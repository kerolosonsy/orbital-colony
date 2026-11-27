package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"orbital/colony/levels"
)

type reviewSeedKey struct{}

// Keep unrelated progress fields intact, including fields from older clients.
type progressSave map[string]json.RawMessage
type reviewState struct {
	Handled int                  `json:"handled,omitempty"`
	Last    string               `json:"last,omitempty"`
	Dates   map[string]time.Time `json:"dates,omitempty"`
	Due     bool                 `json:"due"`
	Level   string               `json:"level,omitempty"`
}

func reviewData(save progressSave) reviewState {
	var state reviewState
	json.Unmarshal(save["review"], &state)
	return state
}
func completed(save progressSave) map[string]bool {
	var records map[string]struct {
		Stars int `json:"stars"`
	}
	json.Unmarshal(save["progress"], &records)
	done := map[string]bool{}
	for _, l := range levels.All {
		if records[l.ID].Stars > 0 {
			done[l.ID] = true
		}
	}
	return done
}

// Curriculum order breaks ties; map iteration never determines the candidate.
func reviewCandidate(save progressSave) (string, bool) {
	state, done := reviewData(save), completed(save)
	if len(done)-state.Handled < 5 {
		return "", false
	}
	candidate := ""
	for _, l := range levels.All {
		if !done[l.ID] || l.ID == state.Last {
			continue
		}
		if candidate == "" || state.Dates[l.ID].Before(state.Dates[candidate]) {
			candidate = l.ID
		}
	}
	return candidate, candidate != ""
}
func setReview(save progressSave, state reviewState) {
	save["review"], _ = json.Marshal(state)
	state.Level, state.Due = reviewCandidate(save)
	save["review"], _ = json.Marshal(state)
}
func readProgress(path string) (progressSave, error) {
	blob, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return progressSave{}, nil
	}
	if err != nil {
		return nil, err
	}
	var save progressSave
	if err := json.Unmarshal(blob, &save); err != nil {
		return nil, err
	}
	if save == nil {
		save = progressSave{}
	}
	if save["progress"] != nil {
		setReview(save, reviewData(save))
	}
	migrateLearningMetrics(save)
	return save, nil
}
func saveProgress(path string, save progressSave) error {
	blob, err := json.Marshal(save)
	if err != nil {
		return err
	}
	return atomicJSON(path, blob)
}
func mergeReview(old, incoming progressSave) {
	state := reviewData(old)
	before, after := completed(old), completed(incoming)
	for _, l := range levels.All {
		if after[l.ID] && !before[l.ID] {
			state.Last = l.ID
		}
	}
	setReview(incoming, state)
}
func finishReview(save progressSave, id string, now time.Time) {
	state := reviewData(save)
	if state.Dates == nil {
		state.Dates = map[string]time.Time{}
	}
	state.Dates[id] = now.UTC()
	state.Handled = len(completed(save)) / 5 * 5
	setReview(save, state)
}

// Labs retain their own starter: their tests depend on the declared API and
// some include trusted read-only implementations. Never carry saved solutions.
func reviewLevel(l *levels.Level) *levels.Level {
	copy := *l
	copy.En = nil // RunLab must retain the starter already selected for this language.
	copy.CarryFrom = ""
	if copy.Kind != "lab" {
		copy.Starter = ""
	}
	return &copy
}
func (s *server) profileReview(w http.ResponseWriter, r *http.Request) {
	p := s.store(r)
	p.mu.Lock()
	defer p.mu.Unlock()
	pr, ok := s.profileByPath(w, r)
	if !ok {
		return
	}
	path := filepath.Join(p.dir(pr.ID), "progress.json")
	save, err := readProgress(path)
	if err != nil {
		profileError(w, r, err)
		return
	}
	id, due := reviewCandidate(save)
	if !due {
		profileBadRequest(w, r, "bad request", 409)
		return
	}
	switch r.PathValue("action") {
	case "start":
		writeJSON(w, reviewLevel(levelForLang(levels.Get(id), requestLang(r))))
	case "skip":
		finishReview(save, id, time.Now())
		if err := saveProgress(path, save); err != nil {
			profileError(w, r, err)
			return
		}
		writeJSON(w, reviewData(save))
	default:
		profileBadRequest(w, r, "bad request", 400)
	}
}
