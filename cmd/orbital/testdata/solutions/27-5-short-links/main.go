package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type LinkStore interface {
	Put(code, target string) error
	Get(code string) (string, error)
}

var ErrTaken = errors.New("code taken")
var ErrNoLink = errors.New("no link")

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62(n uint64) string {
	if n == 0 {
		return "0"
	}
	var out []byte
	for n > 0 {
		out = append(out, alphabet[n%62])
		n /= 62
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}

func ParseBase62(s string) (uint64, error) {
	if s == "" || len(s) > 1 && s[0] == '0' {
		return 0, errors.New("not canonical base62")
	}
	var n uint64
	for i := 0; i < len(s); i++ {
		d := strings.IndexByte(alphabet, s[i])
		if d < 0 {
			return 0, fmt.Errorf("bad base62 digit %q", s[i])
		}
		if n > (math.MaxUint64-uint64(d))/62 {
			return 0, errors.New("base62 overflow")
		}
		n = n*62 + uint64(d)
	}
	return n, nil
}

type memoryStore struct {
	mu    sync.Mutex
	links map[string]string
}

func NewMemoryStore() LinkStore { return &memoryStore{links: map[string]string{}} }

func (m *memoryStore) Put(code, target string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.links[code]; ok {
		return ErrTaken
	}
	m.links[code] = target
	return nil
}

func (m *memoryStore) Get(code string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	target, ok := m.links[code]
	if !ok {
		return "", ErrNoLink
	}
	return target, nil
}

func plain(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func validTarget(target string) bool {
	for i := 0; i < len(target); i++ {
		if target[i] <= ' ' || target[i] == 0x7f {
			return false
		}
	}
	u, err := url.Parse(target)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil
}

func validCode(code string) bool {
	if code == "" || len(code) > 11 {
		return false
	}
	for i := 0; i < len(code); i++ {
		if strings.IndexByte(alphabet, code[i]) < 0 {
			return false
		}
	}
	return true
}

type desk struct {
	store LinkStore
	next  func() uint64
}

func NewLinkDesk(store LinkStore, next func() uint64) http.Handler {
	return &desk{store: store, next: next}
}

func (d *desk) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/links":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			plain(w, http.StatusMethodNotAllowed, "wrong method\n")
			return
		}
		d.create(w, r)
	case strings.HasPrefix(r.URL.Path, "/go/") && validCode(strings.TrimPrefix(r.URL.Path, "/go/")):
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			plain(w, http.StatusMethodNotAllowed, "wrong method\n")
			return
		}
		target, err := d.store.Get(strings.TrimPrefix(r.URL.Path, "/go/"))
		if errors.Is(err, ErrNoLink) {
			plain(w, http.StatusNotFound, "no link\n")
			return
		}
		if err != nil {
			plain(w, http.StatusInternalServerError, "store failure\n")
			return
		}
		w.Header().Set("Location", target)
		w.WriteHeader(http.StatusFound)
	default:
		plain(w, http.StatusNotFound, "no link\n")
	}
}

func (d *desk) create(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 2049))
	var in struct {
		Target string `json:"target"`
	}
	if err != nil || len(raw) > 2048 || json.Unmarshal(raw, &in) != nil || !validTarget(in.Target) {
		plain(w, http.StatusBadRequest, "bad target\n")
		return
	}
	for attempt := 0; attempt < 3; attempt++ {
		code := Base62(d.next())
		err := d.store.Put(code, in.Target)
		if errors.Is(err, ErrTaken) {
			continue
		}
		if err != nil {
			plain(w, http.StatusInternalServerError, "store failure\n")
			return
		}
		body, _ := json.Marshal(map[string]string{"code": code, "path": "/go/" + code})
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/go/"+code)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(body)
		return
	}
	plain(w, http.StatusServiceUnavailable, "no free code\n")
}

func main() {}
