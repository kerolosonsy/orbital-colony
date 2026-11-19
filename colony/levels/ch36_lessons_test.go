package levels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestCh36Lessons(t *testing.T) {
	levelIDs := []string{"36-1-metrics", "36-2-containers", "36-3-ops", "36-boss-storm-control"}
	prefixes := []string{"c361.", "c362.", "c363.", "c364."}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	allArabicIDs := map[string]bool{}
	for mission, id := range levelIDs {
		level := Get(id)
		if level == nil || level.En == nil || level.Kind != "lab" {
			t.Fatalf("%s: missing lab or English text", id)
		}
		if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: lesson revision must be 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		ch36ValidateLessonCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, prefixes[mission], 200, 320)
		ch36ValidateLessonCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, prefixes[mission], 200, 320)
		for _, step := range level.Steps {
			ch36ClaimGlobalID(t, allArabicIDs, step.ID)
			for _, block := range step.Blocks {
				ch36ClaimGlobalID(t, allArabicIDs, block.ID)
			}
		}
		for i, step := range level.Steps {
			translated := level.En.Steps[i]
			if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
				t.Errorf("%s: Arabic and English step structures differ", step.ID)
				continue
			}
			for j, block := range step.Blocks {
				other := translated.Blocks[j]
				if block.ID != other.ID || block.Kind != other.Kind {
					t.Errorf("%s: Arabic and English block identities differ", block.ID)
				}
				if block.Question != nil {
					if other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options) {
						t.Errorf("%s: question answer or option count differs", block.ID)
					}
				}
				if block.Walkthrough != nil {
					if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
						t.Errorf("%s: walkthrough structure differs", block.ID)
						continue
					}
					for k, frame := range block.Walkthrough.Frames {
						translatedFrame := other.Walkthrough.Frames[k]
						frame.Caption, translatedFrame.Caption = "", ""
						if !reflect.DeepEqual(frame, translatedFrame) {
							t.Errorf("%s frame %d: execution state differs", block.ID, k)
						}
					}
				}
				if block.Sandbox != nil && (other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind || !ch36SameConfigShape(t, block.Sandbox.Config, other.Sandbox.Config)) {
					t.Errorf("%s: Arabic and English sandbox Config structure differs", block.ID)
				}
			}
		}
	}
}

func ch36ClaimGlobalID(t *testing.T, seen map[string]bool, id string) {
	t.Helper()
	if id == "" || seen[id] {
		t.Errorf("empty or duplicate chapter 36 ID %q", id)
	}
	seen[id] = true
}

func ch36ValidateLessonCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string, minWords, maxWords int) {
	t.Helper()
	seen := map[string]bool{}
	claim := func(id string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q has the wrong mission prefix", name, id)
		}
		seen[id] = true
	}
	words := 0
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i)
		}
		words += ch36TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch36TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown contains another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", name, block.ID)
					continue
				}
				words += ch36TextWords(strings.Join(block.Question.Options, " "))
				words += ch36TextWords(strings.Join(block.Question.Feedback, " "))
				for k := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[k]) == "" || strings.TrimSpace(block.Question.Feedback[k]) == "" {
						t.Errorf("%s: %s has empty option or feedback", name, block.ID)
					}
					if strings.Contains(block.Question.Feedback[k], "\n") || !ch36EndsSentence(block.Question.Feedback[k]) {
						t.Errorf("%s: %s feedback %d must be one sentence", name, block.ID, k)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch36ValidateFrames(t, name, block)
			case "sandbox":
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch36-") || block.Sandbox.Config == nil {
					t.Errorf("%s: %s has an incomplete chapter 36 sandbox", name, block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
	} else {
		for i, hint := range hints {
			trimmed := strings.TrimSpace(hint)
			lower := strings.ToLower(trimmed)
			for _, label := range []string{"concept:", "structure:", "location:", "الفكرة:", "الهيكل:", "الموضع:"} {
				if strings.HasPrefix(lower, label) {
					t.Errorf("%s: hint %d includes a label prefix", name, i+1)
				}
			}
			if strings.Contains(hint, "\n") || !ch36EndsSentence(hint) {
				t.Errorf("%s: hint %d must be one sentence", name, i+1)
			}
			words += ch36TextWords(hint)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 3 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || try.Blocks[2].Kind != "sandbox" {
			t.Errorf("%s: Try it must predict, reveal, then show one sandbox", name)
		} else if strings.TrimSpace(ch36QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction example differs from its walkthrough reveal", name)
		}
	}
	if words < minWords || words > maxWords {
		t.Errorf("%s: learner text has %d words; want %d-%d", name, words, minWords, maxWords)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func TestCh36Sandboxes(t *testing.T) {
	type fixtureLevel struct {
		ID      string                     `json:"id"`
		Kind    string                     `json:"kind"`
		Configs map[string]json.RawMessage `json:"configs"`
	}
	var fixture struct {
		Levels []fixtureLevel `json:"levels"`
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch36.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &fixture); err != nil {
		t.Fatal(err)
	}
	fixtures := make(map[string]fixtureLevel, len(fixture.Levels))
	for _, item := range fixture.Levels {
		if _, exists := fixtures[item.ID]; exists {
			t.Errorf("duplicate sandbox fixture for %s", item.ID)
		}
		fixtures[item.ID] = item
	}
	if len(fixtures) != 4 {
		t.Fatalf("expected four chapter 36 fixtures, found %d", len(fixtures))
	}
	for _, id := range []string{"36-1-metrics", "36-2-containers", "36-3-ops", "36-boss-storm-control"} {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing bilingual lesson", id)
		}
		ar, arCount, arMisplaced := ch36FindSandbox(level.Steps)
		en, enCount, enMisplaced := ch36FindSandbox(level.En.Steps)
		if arCount != 1 || enCount != 1 || ar == nil || en == nil || arMisplaced || enMisplaced || ar.Sandbox == nil || en.Sandbox == nil {
			t.Errorf("%s: want one sandbox after the question and walkthrough in both languages", id)
			continue
		}
		if ar.ID != en.ID || ar.Kind != en.Kind || ar.Sandbox.Kind != en.Sandbox.Kind {
			t.Errorf("%s: sandbox identities differ across languages", id)
		}
		if ar.Kind != "sandbox" || !strings.HasPrefix(ar.Sandbox.Kind, "ch36-") || ar.Sandbox.Config == nil {
			t.Errorf("%s: sandbox must use a ch36 kind and Config", id)
		}
		if !ch36SameConfigShape(t, ar.Sandbox.Config, en.Sandbox.Config) {
			t.Errorf("%s: Arabic and English sandbox Config structure differs", id)
		}
		item, ok := fixtures[id]
		if !ok {
			t.Errorf("%s: missing JavaScript fixture", id)
			continue
		}
		if item.Kind != ar.Sandbox.Kind || !ch36JSONEqual(t, ar.Sandbox.Config, item.Configs["ar"]) || !ch36JSONEqual(t, en.Sandbox.Config, item.Configs["en"]) {
			t.Errorf("%s: Go sandbox Config differs from JavaScript fixture", id)
		}
		delete(fixtures, id)
	}
	for id := range fixtures {
		t.Errorf("unexpected chapter 36 sandbox fixture %s", id)
	}
}

func ch36FindSandbox(steps []LessonStep) (*LessonBlock, int, bool) {
	var found *LessonBlock
	count := 0
	misplaced := false
	for stepIndex, step := range steps {
		for blockIndex := range step.Blocks {
			block := &step.Blocks[blockIndex]
			if block.Kind != "sandbox" {
				continue
			}
			count++
			if stepIndex != 2 || blockIndex != len(step.Blocks)-1 {
				misplaced = true
			}
			found = block
		}
	}
	return found, count, misplaced
}

func ch36JSONEqual(t *testing.T, got any, wantJSON json.RawMessage) bool {
	t.Helper()
	var want, normalized any
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Errorf("decode fixture Config: %v", err)
		return false
	}
	body, err := json.Marshal(got)
	if err != nil {
		t.Errorf("encode sandbox Config: %v", err)
		return false
	}
	if err := json.Unmarshal(body, &normalized); err != nil {
		t.Errorf("normalize sandbox Config: %v", err)
		return false
	}
	return reflect.DeepEqual(normalized, want)
}

func ch36SameConfigShape(t *testing.T, first, second any) bool {
	t.Helper()
	normalize := func(value any) (any, bool) {
		body, err := json.Marshal(value)
		if err != nil {
			t.Errorf("encode sandbox Config: %v", err)
			return nil, false
		}
		var normalized any
		if err := json.Unmarshal(body, &normalized); err != nil {
			t.Errorf("normalize sandbox Config: %v", err)
			return nil, false
		}
		return normalized, true
	}
	left, ok := normalize(first)
	if !ok {
		return false
	}
	right, ok := normalize(second)
	if !ok {
		return false
	}
	return ch36JSONShape(left, right)
}

func ch36JSONShape(first, second any) bool {
	switch left := first.(type) {
	case map[string]any:
		right, ok := second.(map[string]any)
		if !ok || len(left) != len(right) {
			return false
		}
		for key, value := range left {
			other, ok := right[key]
			if !ok || !ch36JSONShape(value, other) {
				return false
			}
		}
		return true
	case []any:
		right, ok := second.([]any)
		if !ok || len(left) != len(right) {
			return false
		}
		for index, value := range left {
			if !ch36JSONShape(value, right[index]) {
				return false
			}
		}
		return true
	case string:
		_, ok := second.(string)
		return ok
	default:
		return reflect.DeepEqual(first, second)
	}
}

func ch36EndsSentence(text string) bool {
	text = strings.TrimSpace(text)
	return strings.HasSuffix(text, ".") || strings.HasSuffix(text, "?") || strings.HasSuffix(text, "!") || strings.HasSuffix(text, "؟")
}

func ch36QuestionCode(text string) string {
	var code []string
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") {
			if inFence {
				break
			}
			inFence = true
			continue
		}
		if inFence {
			code = append(code, line)
		}
	}
	return strings.Join(code, "\n")
}

func ch36TextWords(text string) int {
	lines := strings.Split(text, "\n")
	visible := make([]string, 0, len(lines))
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence {
			visible = append(visible, line)
		}
	}
	var plain strings.Builder
	inInlineCode := false
	for _, r := range strings.Join(visible, "\n") {
		if r == '~' {
			inInlineCode = !inInlineCode
			continue
		}
		if !inInlineCode {
			plain.WriteRune(r)
		}
	}
	return len(strings.Fields(plain.String()))
}

func ch36ValidateFrames(t *testing.T, name string, block LessonBlock) int {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Errorf("%s: %s has no initial frame", name, block.ID)
		return 0
	}
	maxLine := strings.Count(block.Code, "\n") + 1
	words := 0
	for i, frame := range frames {
		if frame.Line < 0 || frame.Line > maxLine || (i > 0 && frame.Line == 0) || strings.TrimSpace(frame.Caption) == "" {
			t.Errorf("%s: %s frame %d points outside its code", name, block.ID, i)
		}
		if len(strings.Fields(frame.Caption)) > 8 {
			t.Errorf("%s: %s frame %d caption exceeds eight words", name, block.ID, i)
		}
		if i > 0 && !strings.HasPrefix(frame.Output, frames[i-1].Output) {
			t.Errorf("%s: %s frame %d output is not cumulative", name, block.ID, i)
		}
		if frame.Panic != "" && i != len(frames)-1 {
			t.Errorf("%s: execution continues after a panic", block.ID)
		}
		words += len(strings.Fields(frame.Caption))
	}
	return words
}

func TestCh36StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 36 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH36.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	fence := strings.Repeat(string(rune(96)), 3)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e36[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n" + fence + "go\n(.*?)\n" + fence)
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 36 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch36VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch36VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	fence := strings.Repeat(string(rune(96)), 3)
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected storyboard outcome %q", id, outcome)
	}
	found := false
	levelIDs := []string{"36-1-metrics", "36-2-containers", "36-3-ops", "36-boss-storm-control"}
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						found = true
						if block.Code != code {
							t.Fatalf("%s: lesson code differs from storyboard", id)
						}
					}
				}
			}
		}
	}
	if !found {
		t.Fatalf("%s: no lesson walkthrough uses this example", id)
	}
	parts := strings.SplitN(storyboard, "#### Verified "+id+" — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing documented output", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\n\n" + fence + "text\n(.*?)\n" + fence)
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	wantOutput := match[1]
	dir := t.TempDir()
	imports := make([]string, 0)
	for _, importName := range strings.Split(importNames, ",") {
		imports = append(imports, fmt.Sprintf("%q", importName))
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + strings.Join(imports, ";") + ")\n" + code + "\n",
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, goBin, "build", "-trimpath", "-gcflags=-lang=go1.24", "-o", "example", ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s: %v\n%s", id, err, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	var stdout, stderr strings.Builder
	run.Stdout, run.Stderr = &stdout, &stderr
	err = run.Run()
	if err != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, err, stderr.String())
	}
	if stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), wantOutput+"\n")
	}
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != wantOutput {
						t.Fatalf("%s: final walkthrough output differs from transcript", id)
					}
				}
			}
		}
	}
}
