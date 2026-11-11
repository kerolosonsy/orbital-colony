package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"orbital/colony/levels"
)

// Each profile has a Go module so VS Code can resolve the colony API.
type workspace struct {
	root      string // absolute path of workspace/
	colonyRel string // colony module path relative to root
}

func (ws *workspace) ensureModule() error {
	if err := os.MkdirAll(ws.root, 0o755); err != nil {
		return err
	}
	gomod := filepath.Join(ws.root, "go.mod")
	if _, err := os.Stat(gomod); err == nil {
		return nil
	}
	body := fmt.Sprintf("module play\n\ngo 1.24\n\nrequire orbital/colony v0.0.0\n\nreplace orbital/colony => %s\n", ws.colonyRel)
	return os.WriteFile(gomod, []byte(body), 0o644)
}

func editorFiles(l *levels.Level) []levels.LabFile {
	if len(l.Files) > 0 {
		return l.Files
	}
	return []levels.LabFile{{Name: "main.go", Starter: l.Starter}}
}

// open writes the mission folder (never overwriting files the player
// already edited there) and returns the main file's path.
func (ws *workspace) open(l *levels.Level, current map[string]string) (string, error) {
	if err := ws.ensureModule(); err != nil {
		return "", err
	}
	dir := filepath.Join(ws.root, l.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	first := ""
	for _, f := range editorFiles(l) {
		path := filepath.Join(dir, f.Name)
		if first == "" {
			first = path
		}
		body := f.Starter
		if c, ok := current[f.Name]; ok && !f.ReadOnly && c != "" {
			body = c
		}
		if _, err := os.Stat(path); err == nil && !f.ReadOnly {
			continue
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return "", err
		}
	}
	return first, nil
}

func (ws *workspace) read(l *levels.Level) (map[string]string, time.Time, error) {
	dir := filepath.Join(ws.root, l.ID)
	files := map[string]string{}
	var newest time.Time
	for _, f := range editorFiles(l) {
		path := filepath.Join(dir, f.Name)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, newest, err
		}
		files[f.Name] = string(b)
		if st, err := os.Stat(path); err == nil && st.ModTime().After(newest) {
			newest = st.ModTime()
		}
	}
	return files, newest, nil
}

func (s *server) workspaceOpen(w http.ResponseWriter, r *http.Request) {
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	if s.requestGame(w, r) == nil {
		return
	}
	ws := s.ws.forProfile(requestProfile(r))
	var req struct {
		Level string
		Files map[string]string
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 512<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	l := levelForLang(levels.Get(req.Level), requestLang(r))
	if l == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	path, err := ws.open(l, req.Files)
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	launched := false
	if code, err := exec.LookPath("code"); err == nil {
		// Open the whole workspace (so gopls loads the module) and jump to the file.
		launched = exec.Command(code, ws.root, "-g", path).Start() == nil
	}
	writeJSON(w, map[string]any{"path": path, "root": ws.root, "launched": launched})
}

func (s *server) workspaceFiles(w http.ResponseWriter, r *http.Request) {
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	if s.requestGame(w, r) == nil {
		return
	}
	ws := s.ws.forProfile(requestProfile(r))
	l := levelForLang(levels.Get(r.URL.Query().Get("level")), requestLang(r))
	if l == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	files, mtime, err := ws.read(l)
	if errors.Is(err, os.ErrNotExist) {
		requestError(w, r, "not in workspace", http.StatusNotFound)
		return
	}
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"files": files, "mtime": mtime.UnixMilli()})
}

func (ws *workspace) forProfile(id string) *workspace {
	return &workspace{root: filepath.Join(ws.root, id), colonyRel: filepath.ToSlash(filepath.Join("..", ws.colonyRel))}
}
