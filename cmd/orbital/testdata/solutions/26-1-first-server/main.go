package main

import (
	"bytes"
	"html/template"
	"net/http"
	"time"
)

func gatewayText(w http.ResponseWriter, r *http.Request, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = w.Write([]byte(text))
	}
}
func NewGateway(base, markup string) (http.Handler, error) {
	page, err := template.New("base").Parse(markup)
	if err != nil {
		return nil, err
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/base" && r.URL.Path != "/pulse" {
			gatewayText(w, r, 404, "unknown gate\n")
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			gatewayText(w, r, 405, "use GET\n")
			return
		}
		if r.URL.Path == "/pulse" {
			gatewayText(w, r, 200, "ready\n")
			return
		}
		var rendered bytes.Buffer
		if err := page.Execute(&rendered, struct{ Base string }{base}); err != nil {
			gatewayText(w, r, 500, "page unavailable\n")
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(rendered.Bytes())
	}), nil
}
func GatewayServer(addr string, h http.Handler) *http.Server {
	return &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 2 * time.Second,
		ReadTimeout: 5 * time.Second, WriteTimeout: 8 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}
}
func main() {}
