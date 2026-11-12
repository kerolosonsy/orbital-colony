package main

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestQuizScoresSurviveSaveReload(t *testing.T) {
	s, h := profileServer(t)
	const body = `{"progress":{"1-1-first-signal":{"stars":3}},"prefs":{"sound":false},"quizScores":{"1":4,"38":0}}`
	path := "/api/profiles/default/progress"
	put := profileRequest(t, h, "PUT", path, "", body)
	if put.Code != 200 {
		t.Fatalf("save: %d %s", put.Code, put.Body)
	}
	// Reopen the store to prove scores survive a server restart, not just an echo.
	reopened, err := newProfileStore(s.profiles.root)
	if err != nil {
		t.Fatal(err)
	}
	s.profiles = reopened
	got := profileRequest(t, h, "GET", path, "", "")
	if got.Code != 200 || strings.TrimSpace(got.Body.String()) != body {
		t.Fatalf("reload: %d %s", got.Code, got.Body)
	}
}

func TestLevelsExposeLocalizedQuizzesWithoutChangingArabic(t *testing.T) {
	s := &server{}
	original := make([]levels.QuizQuestion, len(levels.Chapters))
	for i, ch := range levels.Chapters {
		original[i] = ch.Quiz[0]
	}
	for _, lang := range []string{"en", "ar"} {
		req := httptest.NewRequest("GET", "/api/levels", nil)
		req.Header.Set("X-Orbital-Lang", lang)
		w := httptest.NewRecorder()
		s.levels(w, req)
		var response struct {
			Chapters []levels.Chapter `json:"chapters"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if len(response.Chapters) != len(levels.Chapters) {
			t.Fatal("missing chapters")
		}
		for i, ch := range response.Chapters {
			if len(ch.Quiz) != 5 {
				t.Fatalf("chapter %d: missing quiz", ch.N)
			}
			want := original[i].Question
			if lang == "en" {
				want = original[i].En.Question
			}
			if ch.Quiz[0].Question != want {
				t.Fatalf("chapter %d %s: wrong language", ch.N, lang)
			}
		}
	}
}
