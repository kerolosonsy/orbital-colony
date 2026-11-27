// Command orbital serves the Orbital Colony game on localhost.
package main

import (
	"bytes"
	"compress/gzip"
	"context"
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
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
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
	accounts  *accountStore // optional account mode
	public    string        // canonical host name for account mode
	local     bool          // accounts on this computer: VS Code still works
	ai        *aiQuota      // daily helper uses per account
}

func loopbackHost(host string) bool {
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return strings.EqualFold(host, "localhost") || ip != nil && ip.IsLoopback()
}

func validateListenConfig(addr, publicURL string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	if !loopbackHost(host) {
		return errors.New("Orbital Colony runs player programs on this computer. Use a loopback address such as 127.0.0.1:7777")
	}
	if publicURL == "" {
		return nil
	}
	u, err := url.Parse(publicURL)
	if err != nil || !loopbackHost(u.Hostname()) || (u.Scheme != "http" && u.Scheme != "https") ||
		u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || (u.Path != "" && u.Path != "/") {
		return errors.New("-public must be a loopback origin such as http://127.0.0.1:7777")
	}
	return nil
}

func main() {
	addr := flag.String("addr", "127.0.0.1:7777", "listen address (keep it on localhost)")
	colonyDir := flag.String("colony", "colony", "path to the colony module")
	wsDir := flag.String("workspace", "workspace", "folder for missions opened in VS Code")
	savesDir := flag.String("saves", "saves", "folder for persistent colony saves")
	stats := flag.Bool("stats", false, "print the top learning struggle levels")
	publicURL := flag.String("public", os.Getenv("ORBITAL_PUBLIC_URL"), "loopback origin for optional local accounts (e.g. http://127.0.0.1:7777)")
	flag.Parse()
	cfg := readHostedConfig()
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
	if err := validateListenConfig(*addr, *publicURL); err != nil {
		log.Fatal(err)
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
	var profiles *profileStore // with accounts, each account has its own store
	if *publicURL == "" {
		if profiles, err = newProfileStore(savesAbs); err != nil {
			log.Fatal(err)
		}
	}

	web, _ := fs.Sub(orbital.Web, "web")
	s := &server{
		run:      &runner{goBin: goBin, colonyDir: abs},
		token:    randomHex(16),
		profiles: profiles,
		web:      web,
		ws:       &workspace{root: wsAbs, colonyRel: filepath.ToSlash(rel)},
	}
	if *publicURL != "" {
		if err := s.setupHosted(*publicURL, savesAbs, cfg); err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("🚀 Orbital Colony: http://%s", *addr)
	srv := &http.Server{Addr: *addr, Handler: s.routes(), ReadHeaderTimeout: 10 * time.Second, ReadTimeout: 30 * time.Second}
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

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()
	s.registerAccounts(mux)
	s.registerProfiles(mux)
	mux.HandleFunc("GET /{$}", s.index)
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticHandler(s.web)))
	mux.HandleFunc("GET /api/levels", s.levels)
	mux.HandleFunc("GET /api/services", s.services)
	mux.HandleFunc("GET /api/preview", s.preview)
	mux.HandleFunc("POST /api/run", s.runCode)
	mux.HandleFunc("POST /api/trace", s.traceCode)
	mux.HandleFunc("POST /api/format", s.formatCode)
	mux.HandleFunc("GET /api/symbols", s.symbols)
	mux.HandleFunc("POST /api/chat", s.chat)
	mux.HandleFunc("POST /api/grade", s.grade)
	mux.HandleFunc("POST /api/workspace/open", s.workspaceOpen)
	mux.HandleFunc("GET /api/workspace/files", s.workspaceFiles)

	mux.HandleFunc("GET /api/colony", s.colonyStatus)
	mux.HandleFunc("POST /api/colony/run", s.colonyRun)
	s.registerColonyProject(mux)
	s.registerColonySeasons(mux)
	mux.HandleFunc("POST /api/colony/stop", s.colonyStop)
	mux.HandleFunc("POST /api/colony/debug", s.colonyDebugCommand)
	mux.HandleFunc("POST /api/colony/buy", s.colonyBuy)
	mux.HandleFunc("POST /api/colony/claim", s.colonyClaim)
	mux.HandleFunc("POST /api/colony/contract", s.colonyContract)
	mux.HandleFunc("GET /api/colony/stream", s.colonyStream)

	return s.guard(mux)
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
		h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src https://fonts.gstatic.com; img-src 'self' data:; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cross-Origin-Resource-Policy", "same-origin")
		if r.URL.Path == "/healthz" { // the host's health check, on any address
			w.Write([]byte("ok"))
			return
		}
		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			host = r.Host
		}
		if !loopbackHost(host) && (s.public == "" || host != s.public) {
			requestError(w, r, "forbidden host", http.StatusForbidden)
			return
		}
		if s.accounts != nil && s.accounts.secure {
			h.Set("Strict-Transport-Security", "max-age=31536000")
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
		ctx := withLang(r.Context(), requestLang(r))
		if s.accounts != nil {
			acc := s.accounts.current(r)
			if acc == nil && strings.HasPrefix(r.URL.Path, "/api/") && !strings.HasPrefix(r.URL.Path, "/api/auth/") {
				authError(w, "login", http.StatusUnauthorized)
				return
			}
			if acc != nil {
				ctx = withAccount(ctx, acc)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// staticHandler serves the embedded client so a rebuilt server never leaves
// the browser on old scripts: every request revalidates (no-cache) and a
// content-hash ETag turns unchanged files into a cheap 304.
func staticHandler(web fs.FS) http.Handler {
	etags := map[string]string{}
	fs.WalkDir(web, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || (!strings.HasSuffix(path, ".css") && !strings.HasSuffix(path, ".js")) ||
			strings.HasSuffix(path, ".test.js") || path == "sandbox/harness.js" {
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
		tag, ok := etags[strings.TrimPrefix(r.URL.Path, "/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("ETag", tag)
		files.ServeHTTP(w, r)
	})
}

func (s *server) index(w http.ResponseWriter, r *http.Request) {
	name, acc := "index.html", requestAccount(r)
	if s.accounts != nil && acc == nil {
		name = "login.html"
	}
	page, err := fs.ReadFile(s.web, name)
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	page = bytes.ReplaceAll(page, []byte("__ORBITAL_TOKEN__"), []byte(s.token))
	if acc != nil {
		meta := `<meta name="orbital-account" content="` + acc.ID + `">` + "\n"
		if !s.vscodeAvailable() {
			meta += `<meta name="orbital-vscode" content="off">` + "\n"
		}
		page = bytes.Replace(page, []byte(`<meta name="orbital-token"`), []byte(meta+`<meta name="orbital-token"`), 1)
	}
	w.Write(page)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

type cachedJSON struct {
	body, gz []byte
	etag     string
}

// Encode and compress each language once, including during concurrent startup.
var cachedLevels = map[string]func() (cachedJSON, error){
	"ar": sync.OnceValues(func() (cachedJSON, error) { return encodeLevels("ar") }),
	"en": sync.OnceValues(func() (cachedJSON, error) { return encodeLevels("en") }),
}

func encodeLevels(lang string) (cachedJSON, error) {
	body, err := json.Marshal(levelsPayload(lang))
	if err != nil {
		return cachedJSON{}, err
	}
	body = append(body, '\n')
	sum := sha256.Sum256(body)
	var gz bytes.Buffer
	zw := gzip.NewWriter(&gz)
	zw.Write(body)
	zw.Close()
	return cachedJSON{body: body, gz: gz.Bytes(), etag: `W/"` + hex.EncodeToString(sum[:8]) + `"`}, nil
}

func acceptsGzip(header string) bool {
	wildcard := false
	for entry := range strings.SplitSeq(header, ",") {
		name, params, _ := strings.Cut(entry, ";")
		name = strings.ToLower(strings.TrimSpace(name))
		if name != "gzip" && name != "*" {
			continue
		}
		quality := 1.0
		for param := range strings.SplitSeq(params, ";") {
			key, value, _ := strings.Cut(strings.TrimSpace(param), "=")
			if strings.EqualFold(key, "q") {
				quality, _ = strconv.ParseFloat(strings.TrimSpace(value), 64)
			}
		}
		accepted := quality > 0 && quality <= 1
		if name == "gzip" {
			return accepted
		}
		wildcard = accepted
	}
	return wildcard
}

func (s *server) levels(w http.ResponseWriter, r *http.Request) {
	cached, err := cachedLevels[requestLang(r)]()
	if err != nil {
		requestError(w, r, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("ETag", cached.etag)
	w.Header().Set("Vary", "Accept-Encoding, X-Orbital-Lang")
	if match := r.Header.Get("If-None-Match"); match != "" && match == cached.etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if acceptsGzip(strings.Join(r.Header.Values("Accept-Encoding"), ",")) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(cached.gz)
		return
	}
	w.Write(cached.body)
}

func levelsPayload(lang string) map[string]any {
	chapters := append([]levels.Chapter(nil), levels.Chapters...)
	for i := range chapters {
		chapters[i].Quiz = append([]levels.QuizQuestion(nil), chapters[i].Quiz...)
		if lang == "en" {
			for j := range chapters[i].Quiz {
				chapters[i].Quiz[j].QuizText = chapters[i].Quiz[j].En
			}
		}
	}
	missions := make([]*levels.Level, len(levels.All))
	for i, l := range levels.All {
		m := *levelForLang(l, lang)
		m.En = nil // already localized above; the client needs only one language
		if len(m.Steps) > 0 {
			m.Lesson = "" // the client shows Steps instead of the old lesson text
		}
		missions[i] = &m
	}
	if lang == "en" {
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
		if s.store(r) == nil {
			requestError(w, r, "bad request", 400)
			return
		}
		s.store(r).mu.Lock()
		defer s.store(r).mu.Unlock()
		if !profileID.MatchString(requestProfile(r)) {
			requestError(w, r, "bad request", 400)
			return
		}
		reviewPath = filepath.Join(s.store(r).dir(requestProfile(r)), "progress.json")
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
		s.recordProfileLearningRun(s.store(r), requestProfile(r), l.ID, resp)
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
