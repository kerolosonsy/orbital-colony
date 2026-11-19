// Command orbital serves the Orbital Colony game on localhost.
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
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
	stats := flag.Bool("stats", false, "print the top learning struggle levels")
	flag.Parse()
	if *stats {
		savesAbs, err := filepath.Abs(*savesDir)
		if err != nil {
			log.Fatal(err)
		}
		if err := printLearningStats(savesAbs, os.Stdout); err != nil {
			log.Fatal(err)
		}
		return
	}

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
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler(web)))
	mux.HandleFunc("GET /api/levels", s.levels)
	mux.HandleFunc("GET /api/services", s.services)
	mux.HandleFunc("GET /api/preview", s.preview)
	mux.HandleFunc("POST /api/run", s.runCode)
	mux.HandleFunc("POST /api/trace", s.traceCode)
	mux.HandleFunc("POST /api/workspace/open", s.workspaceOpen)
	mux.HandleFunc("GET /api/workspace/files", s.workspaceFiles)

	mux.HandleFunc("GET /api/colony", s.colonyStatus)
	mux.HandleFunc("POST /api/colony/run", s.colonyRun)
	s.registerColonyProject(mux)
	s.registerColonySeasons(mux)
	mux.HandleFunc("POST /api/colony/stop", s.colonyStop)
	mux.HandleFunc("POST /api/colony/buy", s.colonyBuy)
	mux.HandleFunc("POST /api/colony/claim", s.colonyClaim)
	mux.HandleFunc("POST /api/colony/contract", s.colonyContract)
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
		if errors.Is(err, syscall.EADDRINUSE) {
			log.Fatalf("%s is already in use (is the game already running? try: make status). Pick another port: -addr 127.0.0.1:7778", *addr)
		}
		log.Fatal(err)
	}
}

// guard blocks DNS-rebinding and cross-site requests: only pages served by
// this process (which embed the token) may call the API.
func (s *server) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The page holds the API token and runs code on this machine: no
		// framing by other sites (clickjacking), no MIME sniffing, and no
		// Referer leaving the game.
		h := w.Header()
		h.Set("X-Frame-Options", "DENY")
		h.Set("Content-Security-Policy", "frame-ancestors 'none'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
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

// staticHandler serves the embedded client so a rebuilt server never leaves
// the browser on old scripts: every request revalidates (no-cache) and a
// content-hash ETag turns unchanged files into a cheap 304.
func staticHandler(web fs.FS) http.Handler {
	etags := map[string]string{}
	fs.WalkDir(web, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if body, err := fs.ReadFile(web, path); err == nil {
			sum := sha256.Sum256(body)
			etags[path] = `"` + hex.EncodeToString(sum[:8]) + `"`
		}
		return nil
	})
	files := http.FileServerFS(web)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if tag, ok := etags[strings.TrimPrefix(r.URL.Path, "/")]; ok {
			w.Header().Set("ETag", tag)
		}
		files.ServeHTTP(w, r)
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

// levelsCache holds the encoded /api/levels payload per language: the level
// data never changes while the server runs, and it is ~3.6 MB of JSON.
var levelsCache sync.Map // lang -> cachedJSON

type cachedJSON struct {
	body []byte
	etag string
}

func (s *server) levels(w http.ResponseWriter, r *http.Request) {
	lang := requestLang(r)
	entry, ok := levelsCache.Load(lang)
	if !ok {
		body, err := json.Marshal(levelsPayload(r))
		if err != nil {
			requestError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}
		sum := sha256.Sum256(body)
		entry, _ = levelsCache.LoadOrStore(lang, cachedJSON{body: append(body, '\n'), etag: `"` + hex.EncodeToString(sum[:8]) + `"`})
	}
	cached := entry.(cachedJSON)
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", cached.etag)
	if match := r.Header.Get("If-None-Match"); match != "" && match == cached.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(cached.body)
}

func levelsPayload(r *http.Request) map[string]any {
	chapters := append([]levels.Chapter(nil), levels.Chapters...)
	for i := range chapters {
		chapters[i].Quiz = append([]levels.QuizQuestion(nil), chapters[i].Quiz...)
		if requestLang(r) == "en" {
			for j := range chapters[i].Quiz {
				chapters[i].Quiz[j].QuizText = chapters[i].Quiz[j].En
			}
		}
	}
	missions := make([]*levels.Level, len(levels.All))
	for i, l := range levels.All {
		missions[i] = levelForLang(l, requestLang(r))
	}
	if requestLang(r) == "en" {
		for i := range chapters {
			chapters[i].Title, chapters[i].Concept = chapters[i].TitleEn, chapters[i].ConceptEn
			chapters[i].Track = chapters[i].TrackEn
		}
	}
	return map[string]any{"chapters": chapters, "levels": missions}
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
		Review      bool
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
	var reviewSave progressSave
	var reviewPath string
	if req.Review {
		if s.profiles == nil {
			requestError(w, r, "bad request", 400)
			return
		}
		s.profiles.mu.Lock()
		defer s.profiles.mu.Unlock()
		if !profileID.MatchString(requestProfile(r)) {
			requestError(w, r, "bad request", 400)
			return
		}
		reviewPath = filepath.Join(s.profiles.dir(requestProfile(r)), "progress.json")
		var loadErr error
		reviewSave, loadErr = readProgress(reviewPath)
		if loadErr != nil {
			requestError(w, r, loadErr.Error(), 500)
			return
		}
		id, due := reviewCandidate(reviewSave)
		if !due || id != l.ID {
			requestError(w, r, "bad request", 409)
			return
		}
		l = reviewLevel(l)
		r = r.WithContext(context.WithValue(r.Context(), reviewSeedKey{}, time.Now().UnixNano()))
	}
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
	if l.Kind == "lab" {
		resp.Traceable = traceMainPresent(traceProjectFiles(l, req.Code, req.Files))
	} else if traceLevelSupported(l) {
		resp.Traceable = resp.Phase == "run" && resp.Trace == nil
	}
	if req.Review {
		resp.Result.Stars = 0
		if resp.Result.Success {
			finishReview(reviewSave, l.ID, time.Now())
		}
		recordLearningRun(reviewSave, l.ID, resp)
		if err := saveProgress(reviewPath, reviewSave); err != nil {
			if resp.Result.Success {
				requestError(w, r, err.Error(), 500)
				return
			}
			log.Printf("save learning metrics: %v", err)
		}
	} else {
		s.recordProfileLearningRun(requestProfile(r), l.ID, resp)
	}
	writeJSON(w, resp)
}

func (s *server) traceCode(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(withLang(r.Context(), requestLang(r)))
	req, ok := decodeTraceRequest(w, r)
	if !ok {
		return
	}
	lvl := levelForLang(levels.Get(req.Level), requestLang(r))
	if lvl == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	if !traceLevelSupported(lvl) {
		writeJSON(w, traceUnavailable(requestLang(r)))
		return
	}
	project := traceProjectFiles(lvl, req.Code, req.Files)
	if !traceImportsAllowed(lvl, req.Code, project, requestLang(r)) {
		writeJSON(w, traceUnavailable(requestLang(r)))
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	resp, err := s.run.Trace(r.Context(), lvl, req.Code, req.Files)
	if err != nil {
		requestError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, resp)
}

func decodeTraceRequest(w http.ResponseWriter, r *http.Request) (traceRequest, bool) {
	var req traceRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return traceRequest{}, false
	}
	return req, true
}
