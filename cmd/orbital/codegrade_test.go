package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestParseGrade(t *testing.T) {
	g, err := parseGrade("Here you go:\n```json\n{\"score\": 140, \"summary\": \" Clean. \", \"why\": [\"good names\", \"\", \"```go\\nx := 1\\n```\"], \"how\": [\"name the 3\"]}\n```")
	if err != nil {
		t.Fatal(err)
	}
	if g.Score != 100 || g.Summary != "Clean." {
		t.Errorf("score clamps and summary trims: %+v", g)
	}
	if len(g.Why) != 1 || g.Why[0] != "good names" {
		t.Errorf("empty points and code blocks are dropped: %q", g.Why)
	}
	if len(g.How) != 0 {
		t.Errorf("a perfect score needs no steps: %q", g.How)
	}
	if g, _ := parseGrade(`{"score": -5, "summary": "x", "why": ["1","2","3","4","5","6","7","8"]}`); g.Score != 0 || len(g.Why) != gradeMaxPoints {
		t.Errorf("negative clamps to 0 and lists cap at %d: %+v", gradeMaxPoints, g)
	}
	for _, bad := range []string{"no json here", `{"score": 50}`, `{"score": "high"}`} {
		if _, err := parseGrade(bad); err == nil {
			t.Errorf("%q should fail", bad)
		}
	}
}

func TestGradePrompt(t *testing.T) {
	req := gradeRequest{Level: "1-1-first-signal", Files: map[string]string{"main.go": "package main // MINE"}, Stars: 3, Ticks: 3}
	ar := gradePrompt(levelForLang(levels.Get("1-1-first-signal"), "ar"), "ar", req)
	en := gradePrompt(levelForLang(levels.Get("1-1-first-signal"), "en"), "en", req)
	for _, want := range []string{"MINE", "Never write code", "0 to 100", "3 ticks, 3 stars"} {
		if !strings.Contains(ar, want) {
			t.Errorf("prompt lacks %q", want)
		}
	}
	if !strings.Contains(ar, "Egyptian Arabic") || strings.Contains(en, "Egyptian Arabic") {
		t.Error("the language rule follows the player's language")
	}
}

func TestGradeEndpoint(t *testing.T) {
	saved := askModel
	defer func() { askModel = saved }()
	askModel = func(context.Context, string) (string, error) {
		return `{"score": 82, "summary": "كويس", "why": ["أسماء واضحة"], "how": ["سمّي الرقم 3"]}`, nil
	}
	s := &server{}
	post := func(body any) (int, map[string]any) {
		raw, _ := json.Marshal(body)
		rec := httptest.NewRecorder()
		s.grade(rec, httptest.NewRequest("POST", "/api/grade", bytes.NewReader(raw)))
		var out map[string]any
		json.Unmarshal(rec.Body.Bytes(), &out)
		return rec.Code, out
	}
	ok := gradeRequest{Level: "1-1-first-signal", Files: map[string]string{"main.go": "package main"}}
	if code, out := post(ok); code != http.StatusOK || out["score"] != float64(82) || len(out["how"].([]any)) != 1 {
		t.Fatalf("grade = %d %v", code, out)
	}
	if code, _ := post(gradeRequest{Level: "nope", Files: ok.Files}); code != http.StatusNotFound {
		t.Errorf("unknown level: %d", code)
	}
	if code, _ := post(gradeRequest{Level: ok.Level, Files: map[string]string{"main.go": "  "}}); code != http.StatusBadRequest {
		t.Errorf("empty code: %d", code)
	}
	askModel = func(context.Context, string) (string, error) { return "I think it's fine", nil }
	if _, out := post(ok); out["error"] != "failed" {
		t.Errorf("an unreadable grade is an error, not a score: %v", out)
	}
	askModel = func(context.Context, string) (string, error) { return "", errNoCodex }
	if _, out := post(ok); out["setup"] != true {
		t.Errorf("missing codex asks for setup: %v", out)
	}
}
