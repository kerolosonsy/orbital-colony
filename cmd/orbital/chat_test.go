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

func TestChatPromptCarriesMissionCodeAndConversation(t *testing.T) {
	l := levelForLang(levels.Get("6-4-stock-order"), "ar")
	prompt := chatPrompt(l, "ar", chatRequest{
		Code:    "package main\n\nfunc main() { colony.Measure() }",
		Console: "Error: waiting for positions",
		Messages: []chatMessage{
			{Role: "user", Text: "مش فاهم المهمة"},
			{Role: "assistant", Text: "Measure بتقيس الخانة اللي تحتك بس."},
			{Role: "user", Text: "طب أعمل إيه؟"},
		},
	})
	for _, want := range []string{
		"Egyptian Arabic", "Teach; do not solve", "6-4-stock-order", l.Goal,
		"colony.Measure() }", "Error: waiting for positions",
		"Player: مش فاهم المهمة", "Tutor: Measure بتقيس", "Player: طب أعمل إيه؟",
		"Three stars at 19 ticks", l.Hints[0],
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt is missing %q", want)
		}
	}
	if !strings.HasSuffix(prompt, "Tutor:") {
		t.Error("prompt should end where the tutor speaks")
	}
	if en := chatPrompt(levelForLang(levels.Get("6-4-stock-order"), "en"), "en", chatRequest{Messages: []chatMessage{{Role: "user", Text: "hi"}}}); !strings.Contains(en, "friendly English") || strings.Contains(en, "Egyptian") {
		t.Error("English players get English answers")
	}
}

func TestChatKeepsOnlyRecentMessagesAndCode(t *testing.T) {
	var messages []chatMessage
	for i := 0; i < 30; i++ {
		messages = append(messages, chatMessage{Role: "user", Text: "question-" + strings.Repeat("x", i)})
	}
	prompt := chatPrompt(levels.Get("1-1-first-signal"), "en", chatRequest{Code: strings.Repeat("a", chatMaxCode+500) + "TAIL", Messages: messages})
	if strings.Count(prompt, "Player: question-") != chatMaxMessages {
		t.Errorf("want the last %d messages", chatMaxMessages)
	}
	if !strings.Contains(prompt, "TAIL") || strings.Contains(prompt, strings.Repeat("a", chatMaxCode+1)) {
		t.Error("long code keeps its end, clipped")
	}
}

func TestChatEndpoint(t *testing.T) {
	saved := askModel
	defer func() { askModel = saved }()
	var got string
	askModel = func(_ context.Context, prompt string) (string, error) {
		got = prompt
		return "جرّب تمشي شرق الأول.", nil
	}
	s := &server{}
	post := func(body any) map[string]any {
		raw, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/chat", bytes.NewReader(raw))
		rec := httptest.NewRecorder()
		s.chat(rec, req)
		var out map[string]any
		json.Unmarshal(rec.Body.Bytes(), &out)
		if rec.Code != http.StatusOK {
			out = map[string]any{"status": float64(rec.Code)}
		}
		return out
	}
	if out := post(chatRequest{Level: "1-1-first-signal", Messages: []chatMessage{{Role: "user", Text: "help"}}}); out["reply"] != "جرّب تمشي شرق الأول." || !strings.Contains(got, "1-1-first-signal") {
		t.Fatalf("reply = %v", out)
	}
	if out := post(chatRequest{Level: "nope", Messages: []chatMessage{{Role: "user", Text: "help"}}}); out["status"] != float64(404) {
		t.Errorf("unknown level: %v", out)
	}
	if out := post(chatRequest{Level: "1-1-first-signal"}); out["status"] != float64(400) {
		t.Errorf("no question: %v", out)
	}
	askModel = func(context.Context, string) (string, error) { return "", errNoCodex }
	if out := post(chatRequest{Level: "1-1-first-signal", Messages: []chatMessage{{Role: "user", Text: "help"}}}); out["setup"] != true {
		t.Errorf("missing codex should ask for setup: %v", out)
	}
}
