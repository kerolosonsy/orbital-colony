package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"orbital/colony/levels"
)

// Each profile has a Go module so VS Code can resolve the colony API.
type workspace struct {
	root      string // absolute path of workspace/
	colonyRel string // colony module path relative to root
}

func (ws *workspace) ensureModule(l *levels.Level) error {
	if err := os.MkdirAll(ws.root, 0o755); err != nil {
		return err
	}
	catalogLevel := *l
	catalogLevel.Deps = append([]string(nil), l.Deps...)
	for _, level := range levels.All {
		catalogLevel.Deps = append(catalogLevel.Deps, level.Deps...)
	}
	colonyDir := filepath.Clean(filepath.Join(ws.root, ws.colonyRel))
	files, err := moduleFiles(colonyDir, &catalogLevel)
	if err != nil {
		return err
	}
	// VS Code needs a relative replace so the workspace remains movable.
	files["go.mod"] = strings.ReplaceAll(files["go.mod"], filepath.ToSlash(colonyDir), filepath.ToSlash(ws.colonyRel))
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(ws.root, name), []byte(body), 0o644); err != nil {
			return err
		}
	}
	return nil
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
	if err := ws.ensureModule(l); err != nil {
		return "", err
	}
	dir := filepath.Join(ws.root, l.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	first := ""
	for _, f := range editorFiles(l) {
		if !validLabName(f.Name) {
			return "", fmt.Errorf("bad file name %q", f.Name)
		}
		path := filepath.Join(dir, f.Name)
		if first == "" {
			first = path
		}
		body := f.Starter
		if c, ok := current[f.Name]; ok && !f.ReadOnly && strings.HasSuffix(f.Name, ".go") && c != "" {
			body = c
		}
		if _, err := os.Stat(path); err == nil && !f.ReadOnly && strings.HasSuffix(f.Name, ".go") {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", err
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
		if !validLabName(f.Name) {
			return nil, newest, fmt.Errorf("bad file name %q", f.Name)
		}
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
