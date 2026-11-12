package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func NewOpsLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			switch strings.ToLower(a.Key) {
			case "password", "token", "secret", "authorization":
				return slog.String(a.Key, "[REDACTED]")
			}
			return a
		},
	}))
}

type Ops struct {
	log      *slog.Logger
	db       Pinger
	draining atomic.Bool
}

func NewOps(log *slog.Logger, db Pinger) *Ops {
	return &Ops{log: log, db: db}
}

func (o *Ops) Drain() { o.draining.Store(true) }

func (o *Ops) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, body := o.route(r)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprint(w, body)
	o.log.Info("request", "method", r.Method, "path", r.URL.Path, "status", status)
}

func (o *Ops) route(r *http.Request) (int, string) {
	if r.URL.Path != "/livez" && r.URL.Path != "/readyz" {
		return http.StatusNotFound, "not found\n"
	}
	if r.Method != http.MethodGet {
		return http.StatusMethodNotAllowed, "method not allowed\n"
	}
	if r.URL.Path == "/livez" {
		return http.StatusOK, "ok\n"
	}
	if o.draining.Load() {
		return http.StatusServiceUnavailable, "draining\n"
	}
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()
	if err := o.db.Ping(ctx); err != nil {
		o.log.Warn("readiness failed", "error", err.Error())
		return http.StatusServiceUnavailable, "not ready\n"
	}
	return http.StatusOK, "ready\n"
}

func main() {}
