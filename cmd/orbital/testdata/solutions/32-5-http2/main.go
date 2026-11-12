package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"strings"
	"time"
)

type SectorRead func(context.Context, string) (string, error)

func telemetrySector(path string) (string, bool) {
	if !strings.HasPrefix(path, "/sectors/") {
		return "", false
	}
	name := strings.TrimPrefix(path, "/sectors/")
	if len(name) < 1 || len(name) > 16 {
		return "", false
	}
	for _, b := range []byte(name) {
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-') {
			return "", false
		}
	}
	return name, true
}
func sectorHandler(read SectorRead) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sector, ok := telemetrySector(r.URL.Path)
		if !ok {
			http.Error(w, "no sector", http.StatusNotFound)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "use GET", http.StatusMethodNotAllowed)
			return
		}
		text, err := read(r.Context(), sector)
		if err != nil {
			http.Error(w, "sensor unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write([]byte(text))
	})
}
func NewMultiplexGate(config *tls.Config, read SectorRead) (*http.Server, error) {
	if config == nil || len(config.Certificates) == 0 || read == nil {
		return nil, errors.New("invalid multiplex options")
	}
	secure := config.Clone()
	secure.MinVersion = tls.VersionTLS13
	secure.NextProtos = []string{"h2", "http/1.1"}
	return &http.Server{TLSConfig: secure, Handler: sectorHandler(read), ReadHeaderTimeout: 2 * time.Second, IdleTimeout: 30 * time.Second, MaxHeaderBytes: 16 << 10}, nil
}
func main() {}
