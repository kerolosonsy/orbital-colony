package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"orbital/colony/levels"
)

// The in-game helper answers through the OpenAI Codex CLI (`codex exec`),
// which uses the player's own ChatGPT sign-in: no API key lives in the game.
// It runs read-only in an empty temporary folder, so it can only talk.

type chatMessage struct {
	Role string `json:"role"` // "user" or "assistant"
	Text string `json:"text"`
}

type chatRequest struct {
	Level    string        `json:"level"`
	File     string        `json:"file"`
	Code     string        `json:"code"`
	Console  string        `json:"console"`
	Messages []chatMessage `json:"messages"`
}

const (
	chatMaxMessages = 12
	chatMaxCode     = 24 << 10
	chatMaxConsole  = 4 << 10
	chatTimeout     = 150 * time.Second
)

// askModel sends a prompt to the model and returns its reply. Tests swap it.
var askModel = codexAsk

// chatBusy keeps one codex question in flight: codex is heavy and replies in order.
var chatBusy sync.Mutex

var errNoCodex = errors.New("codex not found")

func codexAsk(ctx context.Context, prompt string) (string, error) {
	bin, err := exec.LookPath("codex")
	if err != nil {
		return "", errNoCodex
	}
	chatBusy.Lock()
	defer chatBusy.Unlock()
	// the time limit starts after waiting for an earlier question, so a
	// queued question isn't killed before codex even starts on it
	ctx, cancel := context.WithTimeout(ctx, chatTimeout)
	defer cancel()
	dir, err := os.MkdirTemp("", "orbital-chat-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	out := filepath.Join(dir, "reply.txt")
	cmd := exec.CommandContext(ctx, bin, "exec",
		"--ephemeral", "--skip-git-repo-check", "--ignore-user-config", "--ignore-rules",
		"-s", "read-only", "-C", dir, "--color", "never",
		"-m", "gpt-6-luna", "-c", "model_reasoning_effort=high",
		"-o", out, "-")
	cmd.Stdin = strings.NewReader(prompt)
	output, runErr := cmd.CombinedOutput()
	reply, _ := os.ReadFile(out)
	if text := strings.TrimSpace(string(reply)); text != "" {
		return text, nil
	}
	if runErr != nil {
		return "", fmt.Errorf("%v: %s", runErr, lastLines(string(output), 6))
	}
	return "", errors.New("empty reply")
}

func lastLines(text string, n int) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}

// chatPrompt is everything the tutor knows: the mission, the player's code,
// the last run, and the conversation so far.
func chatPrompt(l *levels.Level, lang string, req chatRequest) string {
	var b strings.Builder
	arabic := lang != "en"
	b.WriteString("You are the friendly tutor inside Orbital Colony, a game that teaches Go: the player writes real Go to command a drone on planet Kepler-G or to pass a lab's tests.\n\n")
	b.WriteString("Rules:\n")
	if arabic {
		b.WriteString("- Answer in friendly Egyptian Arabic. Keep Go code, identifiers and technical terms in English.\n")
	} else {
		b.WriteString("- Answer in clear, friendly English.\n")
	}
	b.WriteString(`- Be short: at most about 120 words, a few lines or bullets. Use Markdown; wrap code in backticks.
- Teach; do not solve. Start with the smallest hint that gets the player moving: point at the line or idea that is wrong, ask a guiding question, or show a tiny snippet of the one missing piece.
- Give a complete solution only if the player explicitly asks for it twice; even then explain each part briefly.
- Ground every answer in this mission, the player's code and the last run below. If the code has an error, explain that first.
- Never invent game API functions. Only use what appears in the mission, the lesson, the starter or the player's code, plus the Go standard library.
- You cannot run or edit anything; just talk.

`)
	fmt.Fprintf(&b, "## Mission %s: %s\nGoal: %s\n", l.ID, l.Title, l.Goal)
	if l.Kind == "lab" {
		b.WriteString("Kind: lab (Run executes tests; there is no drone).\n")
	} else {
		fmt.Fprintf(&b, "Kind: drone mission. Three stars at %d ticks or fewer.\n", l.Gold)
	}
	if lesson := strings.TrimSpace(l.Lesson); lesson != "" {
		fmt.Fprintf(&b, "\n## Lesson summary\n%s\n", lesson)
	}
	for _, step := range l.Steps {
		for _, block := range step.Blocks {
			if block.Kind == "markdown" && strings.TrimSpace(block.Text) != "" {
				fmt.Fprintf(&b, "\n### Lesson: %s\n%s\n", step.Title, block.Text)
			}
		}
	}
	if len(l.Hints) > 0 {
		b.WriteString("\n## Hints the game offers (in order)\n")
		for i, h := range l.Hints {
			fmt.Fprintf(&b, "%d. %s\n", i+1, h)
		}
	}
	fmt.Fprintf(&b, "\n## Starter code\n```go\n%s\n```\n", strings.TrimSpace(l.Starter))
	file := req.File
	if file == "" {
		file = "main.go"
	}
	fmt.Fprintf(&b, "\n## Player's current code (%s)\n```go\n%s\n```\n", file, clip(req.Code, chatMaxCode))
	if console := strings.TrimSpace(clip(req.Console, chatMaxConsole)); console != "" {
		fmt.Fprintf(&b, "\n## Console after the last run\n```\n%s\n```\n", console)
	}
	b.WriteString("\n## Conversation\n")
	messages := req.Messages
	if len(messages) > chatMaxMessages {
		messages = messages[len(messages)-chatMaxMessages:]
	}
	for _, m := range messages {
		who := "Player"
		if m.Role == "assistant" {
			who = "Tutor"
		}
		fmt.Fprintf(&b, "%s: %s\n", who, clip(strings.TrimSpace(m.Text), 2000))
	}
	b.WriteString("Tutor:")
	return b.String()
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:] // keep the end: the newest lines matter most
}

// chat: POST /api/chat → {reply} or {error, setup}.
func (s *server) chat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 128<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	lang := requestLang(r)
	l := levelForLang(levels.Get(req.Level), lang)
	if l == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	if len(req.Messages) == 0 || req.Messages[len(req.Messages)-1].Role != "user" || strings.TrimSpace(req.Messages[len(req.Messages)-1].Text) == "" {
		requestError(w, r, "empty question", http.StatusBadRequest)
		return
	}
	if !s.aiAllowed(w, r) {
		return
	}
	// the model's own time limit starts once it takes the question (see codexAsk)
	reply, err := askModel(r.Context(), chatPrompt(l, lang, req))
	switch {
	case errors.Is(err, errNoCodex):
		writeJSON(w, map[string]any{"error": "codex", "setup": true})
	case err != nil:
		log.Printf("chat: %v", err)
		writeJSON(w, map[string]any{"error": "failed"})
	default:
		writeJSON(w, map[string]string{"reply": reply})
	}
}
