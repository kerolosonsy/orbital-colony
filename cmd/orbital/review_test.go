package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"orbital/colony/levels"
)

func reviewSaveWithCompleted(n int) progressSave {
	records := map[string]any{}
	for _, l := range levels.All[:n] {
		records[l.ID] = map[string]any{"stars": 1, "best": 99, "code": "saved solution"}
	}
	blob, _ := json.Marshal(records)
	return progressSave{"progress": blob}
}

func TestReviewCandidateThresholdAndRotation(t *testing.T) {
	for _, n := range []int{0, 4, 5} {
		save := reviewSaveWithCompleted(n)
		state := reviewState{Last: levels.All[4].ID}
		setReview(save, state)
		id, due := reviewCandidate(save)
		if due != (n == 5) || (due && id != levels.All[0].ID) {
			t.Fatalf("%d completed: %q, %v", n, id, due)
		}
	}
	save := reviewSaveWithCompleted(5)
	setReview(save, reviewState{Last: levels.All[0].ID})
	if id, _ := reviewCandidate(save); id == levels.All[0].ID {
		t.Fatal("selected just-finished mission")
	}
	setReview(save, reviewState{Last: levels.All[4].ID})
	for i := 0; i < 3; i++ {
		id, due := reviewCandidate(save)
		if !due || id != levels.All[i].ID {
			t.Fatalf("rotation %d: %q, %v", i, id, due)
		}
		finishReview(save, id, time.Date(2026, 1, i+1, 0, 0, 0, 0, time.UTC))
		if _, due := reviewCandidate(save); due {
			t.Fatal("review still due after skip/pass")
		}
		next := reviewSaveWithCompleted(10 + i*5)
		mergeReview(save, next)
		save = next
	}
}

func TestReviewSaveRoundTripAndStaleClient(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	old := reviewSaveWithCompleted(5)
	old["prefs"] = json.RawMessage(`{"current":"1-1"}`)
	if err := saveProgress(path, old); err != nil {
		t.Fatal(err)
	}
	loaded, err := readProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, due := reviewCandidate(loaded); !due {
		t.Fatal("old save did not offer review")
	}
	finishReview(loaded, levels.All[0].ID, time.Now())
	stale := reviewSaveWithCompleted(5)
	mergeReview(loaded, stale)
	if err := saveProgress(path, stale); err != nil {
		t.Fatal(err)
	}
	loaded, err = readProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	if reviewData(loaded).Dates[levels.All[0].ID].IsZero() || reviewData(loaded).Due {
		t.Fatal("stale client overwrote review history")
	}
	if string(loaded["progress"]) != string(old["progress"]) {
		t.Fatal("review changed mission records")
	}
}

func TestReviewRunPreservesStarsAndSolutions(t *testing.T) {
	_, h := profileServer(t)
	save := reviewSaveWithCompleted(5)
	blob, _ := json.Marshal(save)
	put := profileRequest(t, h, "PUT", "/api/profiles/default/progress", "default", string(blob))
	if put.Code != 200 {
		t.Fatal(put.Body)
	}
	start := profileRequest(t, h, "POST", "/api/profiles/default/review/start", "default", "{}")
	var level levels.Level
	if err := json.Unmarshal(start.Body.Bytes(), &level); err != nil {
		t.Fatal(err)
	}
	if level.ID != levels.All[0].ID || level.Starter != "" {
		t.Fatalf("review starter: %s", start.Body)
	}
	code, err := os.ReadFile(filepath.Join("testdata", "solutions", level.ID+".go"))
	if err != nil {
		t.Fatal(err)
	}
	payload, _ := json.Marshal(map[string]any{"level": level.ID, "code": string(code), "review": true})
	run := profileRequest(t, h, "POST", "/api/run", "default", string(payload))
	var response runResponse
	if err := json.Unmarshal(run.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if run.Code != 200 || !response.Result.Success || response.Result.Stars != 0 {
		t.Fatalf("review run: %s", run.Body)
	}
	got := profileRequest(t, h, "GET", "/api/profiles/default/progress", "default", "")
	var after progressSave
	json.Unmarshal(got.Body.Bytes(), &after)
	if string(after["progress"]) != string(save["progress"]) || reviewData(after).Due || reviewData(after).Dates[level.ID].IsZero() {
		t.Fatalf("review changed progress or stayed due: %s", got.Body)
	}
}

func TestReviewLabsKeepOwnStarters(t *testing.T) {
	for _, l := range levels.All {
		if l.Kind != "lab" {
			continue
		}
		for _, lang := range []string{"ar", "en"} {
			localized := levelForLang(l, lang)
			review := reviewLevel(localized)
			if review.CarryFrom != "" || review.En != nil || review.Starter != localized.Starter || !reflect.DeepEqual(review.Files, localized.Files) {
				t.Fatalf("%s (%s): review must retain only its own starter files", l.ID, lang)
			}
		}
	}
}

func TestReviewSkipKeepsFiveMissionCadence(t *testing.T) {
	_, h := profileServer(t)
	save := reviewSaveWithCompleted(7)
	blob, _ := json.Marshal(save)
	put := profileRequest(t, h, "PUT", "/api/profiles/default/progress", "default", string(blob))
	if put.Code != 200 {
		t.Fatal(put.Body)
	}
	skipped := profileRequest(t, h, "POST", "/api/profiles/default/review/skip", "default", "{}")
	var state reviewState
	if err := json.Unmarshal(skipped.Body.Bytes(), &state); err != nil {
		t.Fatal(err)
	}
	if skipped.Code != 200 || state.Due || state.Dates[levels.All[0].ID].IsZero() {
		t.Fatalf("skip: %s", skipped.Body)
	}
	// Re-saving existing completions must not bring the offer back.
	for _, n := range []int{7, 9, 10} {
		next, _ := json.Marshal(reviewSaveWithCompleted(n))
		response := profileRequest(t, h, "PUT", "/api/profiles/default/progress", "default", string(next))
		var stored progressSave
		json.Unmarshal(response.Body.Bytes(), &stored)
		if reviewData(stored).Due != (n == 10) {
			t.Fatalf("%d completed: %s", n, response.Body)
		}
	}
}
