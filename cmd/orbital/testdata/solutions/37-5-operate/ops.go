package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

type Ops struct {
	Now    func() time.Time
	Limit  int
	Window time.Duration
	Log    io.Writer
	NewID  func() string

	mu       sync.Mutex
	counts   map[int]int
	inFlight int
	windows  map[string]window
}

type window struct {
	start time.Time
	used  int
}

var requestIDPattern = regexp.MustCompile(`^[A-Za-z0-9-]{1,64}$`)

// recorder passes successful replies straight through and buffers error
// replies so their body can be rewritten into the shared error shape.
type recorder struct {
	http.ResponseWriter
	status int
	buf    bytes.Buffer
}

func (r *recorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}
	r.status = status
	if status < 400 {
		r.ResponseWriter.WriteHeader(status)
	}
}

func (r *recorder) Write(p []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}
	if r.status >= 400 {
		return r.buf.Write(p)
	}
	return r.ResponseWriter.Write(p)
}

func writeError(w http.ResponseWriter, status int, code, id string) {
	writeJSON(w, status, map[string]string{"error": code, "request_id": id})
}

func (o *Ops) allow(key string) (bool, time.Duration) {
	o.mu.Lock()
	defer o.mu.Unlock()
	now := o.Now()
	win := o.windows[key]
	if win.start.IsZero() || now.Sub(win.start) >= o.Window {
		win = window{start: now}
	}
	if win.used >= o.Limit {
		o.windows[key] = win
		return false, o.Window - now.Sub(win.start)
	}
	win.used++
	o.windows[key] = win
	return true, 0
}

func (o *Ops) metrics(w http.ResponseWriter) {
	o.mu.Lock()
	defer o.mu.Unlock()
	codes := []int{}
	for code := range o.counts {
		codes = append(codes, code)
	}
	sort.Ints(codes)
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	for _, code := range codes {
		fmt.Fprintf(w, "command_requests_total{code=\"%d\"} %d\n", code, o.counts[code])
	}
	fmt.Fprintf(w, "command_requests_in_flight %d\n", o.inFlight)
}

func Operate(h http.Handler, o *Ops) http.Handler {
	o.counts, o.windows = map[int]int{}, map[string]window{}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if !requestIDPattern.MatchString(id) {
			id = o.NewID()
		}
		w.Header().Set("X-Request-ID", id)
		if r.Method == http.MethodGet && r.URL.Path == "/metrics" {
			o.metrics(w)
			return
		}
		rec := &recorder{ResponseWriter: w}
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		if ok, wait := o.allow(host); !ok {
			secs := int((wait + time.Second - 1) / time.Second)
			w.Header().Set("Retry-After", fmt.Sprint(secs))
			rec.WriteHeader(http.StatusTooManyRequests)
			rec.buf.WriteString(`{"error":"rate_limited"}`)
		} else {
			o.mu.Lock()
			o.inFlight++
			o.mu.Unlock()
			func() {
				defer func() {
					if v := recover(); v != nil {
						rec.status = 0
						rec.buf.Reset()
						rec.WriteHeader(http.StatusInternalServerError)
						rec.buf.WriteString(`{"error":"internal"}`)
					}
				}()
				h.ServeHTTP(rec, r)
			}()
			o.mu.Lock()
			o.inFlight--
			o.mu.Unlock()
		}
		if rec.status == 0 {
			rec.WriteHeader(http.StatusOK)
		}
		if rec.status >= 400 {
			var body struct {
				Error string `json:"error"`
			}
			if json.Unmarshal(rec.buf.Bytes(), &body) != nil || body.Error == "" {
				body.Error = strings.ToLower(strings.ReplaceAll(http.StatusText(rec.status), " ", "_"))
			}
			w.Header().Del("Content-Length")
			writeError(w, rec.status, body.Error, id)
		}
		o.mu.Lock()
		o.counts[rec.status]++
		line, _ := json.Marshal(map[string]any{"request_id": id, "method": r.Method, "path": r.URL.Path, "status": rec.status})
		o.Log.Write(append(line, '\n'))
		o.mu.Unlock()
	})
}

func Serve(ctx context.Context, ln net.Listener, h http.Handler, drain time.Duration) error {
	srv := &http.Server{Handler: h}
	failed := make(chan error, 1)
	go func() { failed <- srv.Serve(ln) }()
	select {
	case err := <-failed:
		return err
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), drain)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		srv.Close()
		return err
	}
	if err := <-failed; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
