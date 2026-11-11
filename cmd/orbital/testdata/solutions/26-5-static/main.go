package main

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed site
var chartSite embed.FS

type chartAsset struct{ file, media, tag string }

func NewChartFiles() http.Handler {
	root, _ := fs.Sub(chartSite, "site")
	files := http.FileServer(http.FS(root))
	assets := map[string]chartAsset{
		"/console":            {"panel.html", "text/html; charset=utf-8", ""},
		"/charts/notice.txt":  {"charts/notice.txt", "text/plain; charset=utf-8", ""},
		"/charts/sector.json": {"charts/sector.json", "application/json", ""},
	}
	for route, asset := range assets {
		body, _ := fs.ReadFile(root, asset.file)
		asset.tag = fmt.Sprintf("\"%x\"", sha256.Sum256(body))
		assets[route] = asset
	}
	missing, _ := fs.ReadFile(root, "missing.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		asset, ok := assets[r.URL.Path]
		if !ok || (r.URL.RawPath != "" && r.URL.RawPath != r.URL.Path) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(404)
			if r.Method != "HEAD" {
				_, _ = w.Write(missing)
			}
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			w.Header().Set("Allow", "GET, HEAD")
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(405)
			_, _ = w.Write([]byte("read only\n"))
			return
		}
		w.Header().Set("ETag", asset.tag)
		w.Header().Set("Cache-Control", "public, max-age=600")
		if r.Header.Get("If-None-Match") == asset.tag {
			w.WriteHeader(304)
			return
		}
		w.Header().Set("Content-Type", asset.media)
		copyRequest := r.Clone(r.Context())
		copyRequest.URL.Path = "/" + asset.file
		// The wrapper owns the deliberately narrow conditional policy.
		copyRequest.Header.Del("If-None-Match")
		copyRequest.Header.Del("If-Modified-Since")
		files.ServeHTTP(w, copyRequest)
	})
}
func main() {}
