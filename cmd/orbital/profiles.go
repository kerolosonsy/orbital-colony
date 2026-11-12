package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"orbital/colony/econ"
)

type profile struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Lang    string    `json:"lang"`
}

type profileStore struct {
	mu    sync.Mutex
	root  string
	games map[string]*colonyGame
}

var profileID = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,39}$`)

func newProfileStore(root string) (*profileStore, error) {
	p := &profileStore{root: root, games: map[string]*colonyGame{}}
	p.mu.Lock()
	defer p.mu.Unlock()
	_, err := p.list()
	return p, err
}

func (p *profileStore) dir(id string) string { return filepath.Join(p.root, "profiles", id) }

func atomicJSON(path string, blob []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".json-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	_, writeErr := f.Write(blob)
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(f.Name(), path)
}

func (p *profileStore) save(pr profile) error {
	blob, err := json.Marshal(pr)
	if err != nil {
		return err
	}
	return atomicJSON(filepath.Join(p.dir(pr.ID), "profile.json"), blob)
}

func (p *profileStore) load(id string) (profile, error) {
	var pr profile
	blob, err := os.ReadFile(filepath.Join(p.dir(id), "profile.json"))
	if err != nil {
		return pr, err
	}
	err = json.Unmarshal(blob, &pr)
	return pr, err
}

func (p *profileStore) ensureDefault() (profile, error) {
	pr, err := p.load("default")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return pr, err
	}
	if errors.Is(err, os.ErrNotExist) {
		pr = profile{ID: "default", Name: "Player 1", Created: time.Now().UTC(), Lang: "ar"}
		if err := p.save(pr); err != nil {
			return pr, err
		}
	}
	legacy := filepath.Join(p.root, "colony.json")
	destination := filepath.Join(p.dir("default"), "colony.json")
	if _, err := os.Stat(destination); errors.Is(err, os.ErrNotExist) {
		if err := os.Rename(legacy, destination); err != nil && !errors.Is(err, os.ErrNotExist) {
			return pr, err
		}
	} else if err != nil {
		return pr, err
	}
	return pr, nil
}

// Caller holds mu so profile creation and last-profile deletion stay atomic.
func (p *profileStore) list() ([]profile, error) {
	entries, err := os.ReadDir(filepath.Join(p.root, "profiles"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	profiles := []profile{}
	for _, entry := range entries {
		if !entry.IsDir() || !profileID.MatchString(entry.Name()) {
			continue
		}
		pr, err := p.load(entry.Name())
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, pr)
	}
	if len(profiles) == 0 {
		pr, err := p.ensureDefault()
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, pr)
	}
	sort.Slice(profiles, func(i, j int) bool {
		if profiles[i].Created.Equal(profiles[j].Created) {
			return profiles[i].ID < profiles[j].ID
		}
		return profiles[i].Created.Before(profiles[j].Created)
	})
	return profiles, nil
}

func requestProfile(r *http.Request) string {
	id := r.Header.Get("X-Orbital-Profile")
	if id == "" && r.URL.Path == "/api/colony/stream" {
		id = r.URL.Query().Get("profile")
	}
	if id == "" {
		return "default"
	}
	return id
}

func (s *server) profileGame(r *http.Request) (*colonyGame, error) {
	id := requestProfile(r)
	if !profileID.MatchString(id) {
		return nil, errors.New("invalid profile id")
	}
	if s.profiles == nil {
		return s.colony, nil
	}
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	if id == "default" {
		if _, err := p.ensureDefault(); err != nil {
			return nil, err
		}
	} else if _, err := p.load(id); err != nil {
		return nil, err
	}
	if c := p.games[id]; c != nil {
		return c, nil
	}
	c, err := newColonyGame(p.dir(id))
	if err == nil {
		p.games[id] = c
	}
	return c, err
}

func (s *server) requestGame(w http.ResponseWriter, r *http.Request) *colonyGame {
	c, err := s.profileGame(r)
	if err != nil {
		status := http.StatusInternalServerError
		if !profileID.MatchString(requestProfile(r)) {
			status = http.StatusBadRequest
		}
		if errors.Is(err, os.ErrNotExist) {
			status = http.StatusNotFound
		}
		requestError(w, r, err.Error(), status)
		return nil
	}
	return c
}

func (s *server) stopColonies() {
	if s.profiles == nil {
		s.colony.stop()
		return
	}
	s.profiles.mu.Lock()
	defer s.profiles.mu.Unlock()
	for _, c := range s.profiles.games {
		c.stop()
	}
}

func (s *server) registerProfiles(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/profiles", s.profileList)
	mux.HandleFunc("POST /api/profiles", s.profileCreate)
	mux.HandleFunc("POST /api/profiles/{id}/rename", s.profileUpdate)
	mux.HandleFunc("POST /api/profiles/{id}/lang", s.profileUpdate)
	mux.HandleFunc("POST /api/profiles/{id}/reset", s.profileReset)
	mux.HandleFunc("DELETE /api/profiles/{id}", s.profileDelete)
	mux.HandleFunc("GET /api/profiles/{id}/progress", s.profileProgress)
	mux.HandleFunc("PUT /api/profiles/{id}/progress", s.profileProgress)
	mux.HandleFunc("POST /api/profiles/{id}/review/{action}", s.profileReview)
}

func profileError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, os.ErrNotExist) {
		status = http.StatusNotFound
	}
	profileBadRequest(w, r, err.Error(), status)
}

func profileBadRequest(w http.ResponseWriter, r *http.Request, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errorText(requestLang(r), message)})
}

func (s *server) profileList(w http.ResponseWriter, r *http.Request) {
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	profiles, err := p.list()
	if err != nil {
		profileError(w, r, err)
		return
	}
	writeJSON(w, profiles)
}

func profileInput(w http.ResponseWriter, r *http.Request) (profile, bool) {
	var pr profile
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&pr); err != nil {
		profileBadRequest(w, r, "bad request", 400)
		return pr, false
	}
	pr.Name = strings.TrimSpace(pr.Name)
	return pr, true
}

func (s *server) profileCreate(w http.ResponseWriter, r *http.Request) {
	pr, ok := profileInput(w, r)
	if !ok {
		return
	}
	if pr.Name == "" || (pr.Lang != "" && pr.Lang != "en" && pr.Lang != "ar") {
		profileBadRequest(w, r, "name and valid language required", 400)
		return
	}
	if pr.Lang == "" {
		pr.Lang = "ar"
	}
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	profiles, err := p.list()
	if err != nil {
		profileError(w, r, err)
		return
	}
	pr.ID = freshProfileID(pr.Name, profiles)
	pr.Created = time.Now().UTC()
	if err := p.save(pr); err != nil {
		profileError(w, r, err)
		return
	}
	writeJSON(w, pr)
}

func (s *server) profileByPath(w http.ResponseWriter, r *http.Request) (profile, bool) {
	id := r.PathValue("id")
	if !profileID.MatchString(id) {
		profileBadRequest(w, r, "invalid profile id", 400)
		return profile{}, false
	}
	pr, err := s.profiles.load(id)
	if err != nil {
		profileError(w, r, err)
		return profile{}, false
	}
	return pr, true
}

func (s *server) profileUpdate(w http.ResponseWriter, r *http.Request) {
	input, ok := profileInput(w, r)
	if !ok {
		return
	}
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	pr, ok := s.profileByPath(w, r)
	if !ok {
		return
	}
	if strings.HasSuffix(r.URL.Path, "/rename") {
		if input.Name == "" {
			profileBadRequest(w, r, "name required", 400)
			return
		}
		pr.Name = input.Name
	} else {
		if input.Lang != "en" && input.Lang != "ar" {
			profileBadRequest(w, r, "invalid language", 400)
			return
		}
		pr.Lang = input.Lang
	}
	if err := p.save(pr); err != nil {
		profileError(w, r, err)
		return
	}
	writeJSON(w, pr)
}

func (s *server) removeProfileFiles(id string) error {
	if s.ws != nil {
		if err := os.RemoveAll(filepath.Join(s.ws.root, id)); err != nil {
			return err
		}
	}
	for _, name := range []string{"colony.json", "progress.json"} {
		if err := os.Remove(filepath.Join(s.profiles.dir(id), name)); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func (s *server) profileReset(w http.ResponseWriter, r *http.Request) {
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	pr, ok := s.profileByPath(w, r)
	if !ok {
		return
	}
	if c := p.games[pr.ID]; c != nil {
		c.stop()
		c.mu.Lock()
		c.state = econ.New(time.Now().UnixNano())
		c.lastSave = time.Time{}
		c.mu.Unlock()
	}
	if err := s.removeProfileFiles(pr.ID); err != nil {
		profileError(w, r, err)
		return
	}
	writeJSON(w, pr)
}

func (s *server) profileDelete(w http.ResponseWriter, r *http.Request) {
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	pr, ok := s.profileByPath(w, r)
	if !ok {
		return
	}
	profiles, err := p.list()
	if err != nil {
		profileError(w, r, err)
		return
	}
	if len(profiles) <= 1 {
		profileBadRequest(w, r, "cannot delete the last profile", 409)
		return
	}
	if c := p.games[pr.ID]; c != nil {
		c.stop()
	}
	if err := s.removeProfileFiles(pr.ID); err != nil {
		profileError(w, r, err)
		return
	}
	if err := os.RemoveAll(p.dir(pr.ID)); err != nil {
		profileError(w, r, err)
		return
	}
	delete(p.games, pr.ID)
	writeJSON(w, map[string]bool{"ok": true})
}

func (s *server) profileProgress(w http.ResponseWriter, r *http.Request) {
	p := s.profiles
	p.mu.Lock()
	defer p.mu.Unlock()
	pr, ok := s.profileByPath(w, r)
	if !ok {
		return
	}
	path := filepath.Join(p.dir(pr.ID), "progress.json")
	old, err := readProgress(path)
	if err != nil {
		profileError(w, r, err)
		return
	}
	if r.Method == http.MethodGet {
		if old["progress"] == nil {
			blob, err := os.ReadFile(path)
			if err == nil {
				writeJSON(w, json.RawMessage(blob))
				return
			}
		}
		writeJSON(w, old)
		return
	}
	blob, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 2<<20))
	if err != nil {
		status := 400
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			status = 413
		}
		profileBadRequest(w, r, "invalid or oversized progress", status)
		return
	}
	var incoming progressSave
	if json.Unmarshal(blob, &incoming) != nil || incoming == nil {
		profileBadRequest(w, r, "invalid JSON", 400)
		return
	}
	if incoming["progress"] == nil && old["progress"] == nil {
		if err := atomicJSON(path, blob); err != nil {
			profileError(w, r, err)
			return
		}
		writeJSON(w, json.RawMessage(blob))
		return
	}
	mergeReview(old, incoming)
	if err := saveProgress(path, incoming); err != nil {
		profileError(w, r, err)
		return
	}
	writeJSON(w, incoming)
}

func freshProfileID(name string, profiles []profile) string {
	used := map[string]bool{}
	for _, existing := range profiles {
		used[existing.ID] = true
	}
	slug := strings.Trim(regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(name), "-"), "-")
	if len(slug) > 32 {
		slug = strings.TrimRight(slug[:32], "-")
	}
	id := slug
	for n := 1; id == "" || used[id]; n++ {
		if slug == "" {
			id = fmt.Sprintf("player-%d", n)
		} else {
			id = fmt.Sprintf("%s-%d", slug, n)
		}
	}
	return id
}
