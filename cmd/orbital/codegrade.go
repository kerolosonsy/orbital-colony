package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"orbital/colony/levels"
)

// The code grade: after a mission passes, the helper's model reads the code
// and scores its craft out of 100, says why, and says how to reach 100. It
// explains in words only; it never writes the code for the player.

type gradeRequest struct {
	Level string            `json:"level"`
	Files map[string]string `json:"files"`
	Stars int               `json:"stars"`
	Ticks int               `json:"ticks"`
}

type gradeResult struct {
	Score   int      `json:"score"`
	Summary string   `json:"summary"`
	Why     []string `json:"why"`
	How     []string `json:"how"`
}

const gradeMaxPoints = 6

func gradePrompt(l *levels.Level, lang string, req gradeRequest) string {
	var b strings.Builder
	b.WriteString("You grade the craft of a player's Go code in Orbital Colony, a game that teaches Go. The code already passes the mission; judge how well it is written.\n\n")
	b.WriteString(`Score from 0 to 100 using this rubric:
- Correctness and robustness beyond the one test map (edge cases, no luck): 25
- Clarity: naming, structure, small focused pieces, no dead code: 25
- Idiomatic Go for what the lessons have taught so far (gofmt style, range, error handling, standard library): 20
- Simplicity and efficiency (no needless work; for drone missions, ticks against the three-star target): 20
- Readability details: comments where they help, consistent formatting, magic numbers named: 10
Judge against what this mission and earlier lessons teach: never deduct for not using a feature the player has not met yet. Clean, correct code that does exactly the job deserves 90 or more. 100 means you would not change a single thing.

Rules:
- Reply with ONE JSON object and nothing else: {"score": <int>, "summary": "<one sentence>", "why": ["..."], "how": ["..."]}.
- "why": up to 6 short points explaining the score: what earned points and what cost points, each naming the line or identifier.
- "how": up to 6 concrete steps to reach 100, most valuable first. If the score is 100, "how" is empty.
- Never write code: no code blocks, no complete Go lines, no rewritten functions. You may name identifiers, functions or keywords in backticks and describe the change in words.
`)
	if lang != "en" {
		b.WriteString("- Write summary, why and how in friendly Egyptian Arabic; keep Go terms and identifiers in English.\n")
	} else {
		b.WriteString("- Write summary, why and how in clear, friendly English.\n")
	}
	fmt.Fprintf(&b, "\n## Mission %s: %s\nGoal: %s\n", l.ID, l.Title, l.Goal)
	if l.Kind == "lab" {
		b.WriteString("Kind: lab (the player's code passed the tests).\n")
	} else {
		fmt.Fprintf(&b, "Kind: drone mission. This run: %d ticks, %d stars. Three stars at %d ticks or fewer.\n", req.Ticks, req.Stars, l.Gold)
	}
	for _, step := range l.Steps {
		for _, block := range step.Blocks {
			if block.Kind == "markdown" && strings.TrimSpace(block.Text) != "" {
				fmt.Fprintf(&b, "\n### Lesson: %s\n%s\n", step.Title, block.Text)
			}
		}
	}
	if lesson := strings.TrimSpace(l.Lesson); lesson != "" {
		fmt.Fprintf(&b, "\n## Lesson\n%s\n", lesson)
	}
	fmt.Fprintf(&b, "\n## Starter code\n```go\n%s\n```\n", strings.TrimSpace(l.Starter))
	names := make([]string, 0, len(req.Files))
	for name := range req.Files {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(&b, "\n## Player's code: %s\n```go\n%s\n```\n", name, clip(req.Files[name], chatMaxCode))
	}
	return b.String()
}

var jsonObject = regexp.MustCompile(`(?s)\{.*\}`)

// parseGrade reads the model's JSON (tolerating text or fences around it)
// and keeps it inside the contract: a 0–100 score, short lists, no code.
func parseGrade(reply string) (gradeResult, error) {
	var g gradeResult
	raw := jsonObject.FindString(reply)
	if raw == "" {
		return g, errors.New("no JSON in reply")
	}
	if err := json.Unmarshal([]byte(raw), &g); err != nil {
		return g, err
	}
	g.Score = max(0, min(100, g.Score))
	g.Summary = strings.TrimSpace(g.Summary)
	g.Why, g.How = tidyPoints(g.Why), tidyPoints(g.How)
	if g.Score == 100 {
		g.How = nil
	}
	if g.Summary == "" && len(g.Why) == 0 {
		return g, errors.New("empty grade")
	}
	return g, nil
}

func tidyPoints(points []string) []string {
	out := []string{}
	for _, p := range points {
		p = strings.TrimSpace(p)
		if p == "" || strings.Contains(p, "```") {
			continue // a code block slipped in: drop it rather than hand over code
		}
		out = append(out, p)
		if len(out) == gradeMaxPoints {
			break
		}
	}
	return out
}

// grade: POST /api/grade → gradeResult or {error, setup}.
func (s *server) grade(w http.ResponseWriter, r *http.Request) {
	var req gradeRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 256<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	lang := requestLang(r)
	l := levelForLang(levels.Get(req.Level), lang)
	if l == nil {
		requestError(w, r, "unknown level", http.StatusNotFound)
		return
	}
	code := 0
	for _, text := range req.Files {
		code += len(strings.TrimSpace(text))
	}
	if code == 0 {
		requestError(w, r, "no code", http.StatusBadRequest)
		return
	}
	if !s.aiAllowed(w, r) {
		return
	}
	// the model's own time limit starts once it takes the question (see codexAsk)
	reply, err := askModel(r.Context(), gradePrompt(l, lang, req))
	if errors.Is(err, errNoCodex) {
		writeJSON(w, map[string]any{"error": "codex", "setup": true})
		return
	}
	var g gradeResult
	if err == nil {
		g, err = parseGrade(reply)
	}
	if err != nil {
		log.Printf("grade: %v", err)
		writeJSON(w, map[string]any{"error": "failed"})
		return
	}
	writeJSON(w, g)
}
