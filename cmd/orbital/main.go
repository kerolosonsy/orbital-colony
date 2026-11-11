// Command orbital serves the Orbital Colony game on localhost.
package main

import (
	"bytes"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"flag"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"orbital"
	"orbital/colony/levels"
)

type server struct {
	run       *runner
	token     string
	mu        sync.Mutex // one player program at a time
	web       fs.FS
	colony    *colonyGame
	profiles  *profileStore
	colonyOps sync.Mutex
	ws        *workspace
}

func main() {
	addr := flag.String("addr", "127.0.0.1:7777", "listen address (keep it on localhost)")
	colonyDir := flag.String("colony", "colony", "path to the colony module")
	wsDir := flag.String("workspace", "workspace", "folder for missions opened in VS Code")
	savesDir := flag.String("saves", "saves", "folder for persistent colony saves")
	flag.Parse()

	abs, err := filepath.Abs(*colonyDir)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(abs, "go.mod")); err != nil {
		log.Fatalf("colony module not found at %s (run from the project root or pass -colony)", abs)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		log.Fatal("the go toolchain must be on PATH: ", err)
	}
	host, _, err := net.SplitHostPort(*addr)
	if err != nil {
		log.Fatal(err)
	}
	if ip := net.ParseIP(host); host != "localhost" && (ip == nil || !ip.IsLoopback()) {
		log.Fatal("refusing to listen on a non-loopback address: player code runs on this machine")
	}

	wsAbs, err := filepath.Abs(*wsDir)
	if err != nil {
		log.Fatal(err)
	}
	rel, err := filepath.Rel(wsAbs, abs)
	if err != nil {
		log.Fatal(err)
	}

	savesAbs, err := filepath.Abs(*savesDir)
	if err != nil {
		log.Fatal(err)
	}
	profiles, err := newProfileStore(savesAbs)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 16)
	rand.Read(b)
	web, _ := fs.Sub(orbital.Web, "web")
	s := &server{
		run:      &runner{goBin: goBin, colonyDir: abs},
		token:    hex.EncodeToString(b),
		profiles: profiles,
		web:      web,
		ws:       &workspace{root: wsAbs, colonyRel: filepath.ToSlash(rel)},
	}

	mux := http.NewServeMux()
	s.registerProfiles(mux)
	mux.HandleFunc("GET /{$}", s.index)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(web)))
	mux.HandleFunc("GET /api/levels", s.levels)
	mux.HandleFunc("GET /api/preview", s.preview)
	mux.HandleFunc("POST /api/run", s.runCode)
	mux.HandleFunc("POST /api/workspace/open", s.workspaceOpen)
	mux.HandleFunc("GET /api/workspace/files", s.workspaceFiles)

	mux.HandleFunc("GET /api/colony", s.colonyStatus)
	mux.HandleFunc("POST /api/colony/run", s.colonyRun)
	mux.HandleFunc("POST /api/colony/stop", s.colonyStop)
	mux.HandleFunc("POST /api/colony/buy", s.colonyBuy)
	mux.HandleFunc("POST /api/colony/claim", s.colonyClaim)
	mux.HandleFunc("GET /api/colony/stream", s.colonyStream)

	log.Printf("🚀 Orbital Colony: http://%s", *addr)
	srv := &http.Server{Addr: *addr, Handler: s.guard(mux), ReadHeaderTimeout: 10 * time.Second}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(shutdown)
	go func() {
		<-shutdown
		s.colonyOps.Lock()
		defer s.colonyOps.Unlock()
		s.stopColonies()
		if err := srv.Close(); err != nil {
			log.Printf("server close: %v", err)
		}
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// guard blocks DNS-rebinding and cross-site requests: only pages served by
// this process (which embed the token) may call the API.
func (s *server) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		if ip := net.ParseIP(host); host != "localhost" && (ip == nil || !ip.IsLoopback()) {
			requestError(w, r, "forbidden host", http.StatusForbidden)
			return
		}
		token := r.Header.Get("X-Orbital-Token")
		if token == "" && r.Method == http.MethodGet && r.URL.Path == "/api/colony/stream" {
			token = r.URL.Query().Get("token")
		}
		if strings.HasPrefix(r.URL.Path, "/api/") &&
			subtle.ConstantTimeCompare([]byte(token), []byte(s.token)) != 1 {
			requestError(w, r, "missing token", http.StatusForbidden)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/") && !profileID.MatchString(requestProfile(r)) {
			profileBadRequest(w, r, "invalid profile id", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r.WithContext(withLang(r.Context(), requestLang(r))))
	})
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	page, err := fs.ReadFile(s.web, "index.html")
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Write(bytes.ReplaceAll(page, []byte("__ORBITAL_TOKEN__"), []byte(s.token)))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (s *server) levels(w http.ResponseWriter, r *http.Request) {
	chapters := append([]levels.Chapter(nil), levels.Chapters...)
	missions := make([]*levels.Level, len(levels.All))
	for i, l := range levels.All {
		missions[i] = levelForLang(l, requestLang(r))
	}
	if requestLang(r) == "en" {
		for i := range chapters {
			chapters[i].Title, chapters[i].Concept = chapters[i].TitleEn, chapters[i].ConceptEn
		}
	}
	writeJSON(w, map[string]any{"chapters": chapters, "levels": missions})
}

func (s *server) preview(w http.ResponseWriter, r *http.Request) {
	l := levels.Get(r.URL.Query().Get("level"))
	if l == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	snapshot, err := previewSnapshot(r.Context(), l.ID, requestLang(r))
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(snapshot)
}

func (s *server) runCode(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(withLang(r.Context(), requestLang(r)))
	var req struct {
		Level, Code string
		Files       map[string]string
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	l := levelForLang(levels.Get(req.Level), requestLang(r))
	if l == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var resp *runResponse
	var err error
	if l.Kind == "lab" {
		resp, err = s.run.RunLab(r.Context(), l, req.Files)
	} else {
		resp, err = s.run.Run(r.Context(), l, req.Code)
	}
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, resp)
}
