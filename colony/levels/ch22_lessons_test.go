package levels

import (
	"bytes"
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

func TestCh22Lessons(t *testing.T) {
	ids := []string{"22-1-tree", "22-2-organizer", "22-3-atomic-save", "22-4-duplicates"}
	arTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	enTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seen := map[string]bool{}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English copy", id)
		}
		if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: expected a lab with revision 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		ch22ValidateLessonCopy(t, id+" Arabic", level.Steps, level.Hints, arTitles, fmt.Sprintf("c22%d.", mission+1), seen)
		ch22ValidateLessonCopy(t, id+" English", level.En.Steps, level.En.Hints, enTitles, fmt.Sprintf("c22%d.", mission+1), nil)
		ch22CompareLessonCopies(t, id, level.Steps, level.En.Steps)
	}
}

func ch22ValidateLessonCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string, seen map[string]bool) {
	t.Helper()
	claim := func(id string) {
		if id == "" || !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: invalid ID %q for prefix %q", name, id, prefix)
		}
		if seen != nil && seen[id] {
			t.Errorf("%s: globally duplicated ID %q", name, id)
		}
		if seen != nil {
			seen[id] = true
		}
	}
	words := 0
	sandboxCount := 0
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has the wrong title or no blocks", name, i)
		}
		words += ch22TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			words += ch22TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Text == "" || block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s has an invalid markdown block", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has an invalid question", name, block.ID)
					continue
				}
				words += ch22TextWords(strings.Join(block.Question.Options, " "))
				words += ch22TextWords(strings.Join(block.Question.Feedback, " "))
				for _, feedback := range block.Question.Feedback {
					if strings.TrimSpace(feedback) == "" || !ch22EndsSentence(feedback) {
						t.Errorf("%s: %s feedback must be one complete sentence", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch22ValidateFrames(t, name, block)
			case "sandbox":
				sandboxCount++
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch22-") || len(block.Sandbox.Config) == 0 {
					t.Errorf("%s: %s has an invalid chapter 22 sandbox", name, block.ID)
				}
			default:
				t.Errorf("%s: %s has disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	ch22ValidateHints(t, name, hints, &words)
	if sandboxCount != 1 {
		t.Errorf("%s: expected exactly one sandbox, found %d", name, sandboxCount)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch22ValidateHints(t *testing.T, name string, hints []string, words *int) {
	t.Helper()
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
		return
	}
	for i, hint := range hints {
		trimmed := strings.TrimSpace(hint)
		for _, prefix := range []string{"Concept:", "Structure:", "Location:", "المفهوم:", "الهيكل:", "الموضع:", "الفكرة:", "التركيب:", "المكان:"} {
			if strings.HasPrefix(trimmed, prefix) {
				t.Errorf("%s: hint %d has a position label prefix", name, i+1)
			}
		}
		if !ch22EndsSentence(trimmed) {
			t.Errorf("%s: hint %d must be one complete sentence", name, i+1)
		}
		*words += ch22TextWords(hint)
	}
}

func ch22CompareLessonCopies(t *testing.T, id string, arabic, english []LessonStep) {
	t.Helper()
	for i, step := range arabic {
		other := english[i]
		if step.ID != other.ID || len(step.Blocks) != len(other.Blocks) {
			t.Errorf("%s: step %d differs across languages", id, i)
			continue
		}
		if i == 2 {
			if len(step.Blocks) != 3 || step.Blocks[0].Kind != "question" || step.Blocks[1].Kind != "walkthrough" || step.Blocks[2].Kind != "sandbox" {
				t.Errorf("%s: Try it must predict, reveal, then include one sandbox", id)
			} else if ch22QuestionCode(step.Blocks[0].Text) != step.Blocks[1].Code {
				t.Errorf("%s: prediction code differs from its walkthrough", id)
			}
		}
		for j, block := range step.Blocks {
			translated := other.Blocks[j]
			if block.ID != translated.ID || block.Kind != translated.Kind {
				t.Errorf("%s: block identity differs for %s", id, block.ID)
			}
			if block.Question != nil && (translated.Question == nil || block.Question.Answer != translated.Question.Answer || len(block.Question.Options) != len(translated.Question.Options)) {
				t.Errorf("%s: question answer or option count differs for %s", id, block.ID)
			}
			if block.Sandbox != nil && (translated.Sandbox == nil || block.Sandbox.Kind != translated.Sandbox.Kind || !reflect.DeepEqual(block.Sandbox.Config, translated.Sandbox.Config)) {
				t.Errorf("%s: sandbox config differs for %s", id, block.ID)
			}
			if block.Walkthrough != nil {
				if translated.Walkthrough == nil || block.Code != translated.Code || block.ExampleID != translated.ExampleID || !reflect.DeepEqual(block.Imports, translated.Imports) {
					t.Errorf("%s: walkthrough differs for %s", id, block.ID)
					continue
				}
				if len(block.Walkthrough.Frames) != len(translated.Walkthrough.Frames) {
					t.Errorf("%s: frame count differs for %s", id, block.ID)
					continue
				}
				for k, frame := range block.Walkthrough.Frames {
					otherFrame := translated.Walkthrough.Frames[k]
					frame.Caption, otherFrame.Caption = "", ""
					if !reflect.DeepEqual(frame, otherFrame) {
						t.Errorf("%s: execution state differs at %s frame %d", id, block.ID, k)
					}
				}
			}
		}
	}
}

func ch22QuestionCode(text string) string {
	var code []string
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
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

func ch22ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
			t.Errorf("%s: %s continues after panic", name, block.ID)
		}
		words += len(strings.Fields(frame.Caption))
	}
	return words
}

func ch22TextWords(text string) int {
	var visible []string
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence {
			visible = append(visible, line)
		}
	}
	var plain strings.Builder
	inInline := false
	for _, character := range strings.Join(visible, "\n") {
		if character == '~' {
			inInline = !inInline
			continue
		}
		if !inInline {
			plain.WriteRune(character)
		}
	}
	return len(strings.Fields(plain.String()))
}

func ch22EndsSentence(text string) bool {
	text = strings.TrimSpace(text)
	return strings.HasSuffix(text, ".") || strings.HasSuffix(text, "؟") || strings.HasSuffix(text, "!")
}

func TestCh22SandboxFixtures(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch22.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]map[string]any
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	ids := []string{"22-1-tree", "22-2-organizer", "22-3-atomic-save", "22-4-duplicates"}
	for _, id := range ids {
		level := Get(id)
		block := level.Steps[2].Blocks[2]
		want, ok := fixtures[block.Sandbox.Kind]
		if !ok {
			t.Errorf("%s: missing fixture for %s", id, block.Sandbox.Kind)
			continue
		}
		actualJSON, err := json.Marshal(block.Sandbox.Config)
		if err != nil {
			t.Fatal(err)
		}
		fixtureJSON, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(actualJSON, fixtureJSON) {
			t.Errorf("%s: sandbox config differs from the JavaScript fixture", id)
		}
	}
	if len(fixtures) != len(ids) {
		t.Errorf("chapter 22 fixture has %d kinds; want %d", len(fixtures), len(ids))
	}
}

func TestCh22StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH22.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e22[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n```go\n(.*?)\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 22 examples, got %d", len(examples))
	}
	ch22RunStoryboardExamples(t, goBin, storyboard, examples, []string{"22-1-tree", "22-2-organizer", "22-3-atomic-save", "22-4-duplicates"})
}

func ch22RunStoryboardExamples(t *testing.T, goBin, storyboard string, examples [][]string, levelIDs []string) {
	t.Helper()
	seen := map[string]bool{}
	for _, example := range examples {
		id := example[1]
		if seen[id] {
			t.Fatalf("duplicate example ID %q", id)
		}
		seen[id] = true
		ch22VerifyExample(t, goBin, storyboard, example, levelIDs)
	}
	expected := map[string]bool{}
	for _, levelID := range levelIDs {
		for _, step := range Get(levelID).Steps {
			for _, block := range step.Blocks {
				if block.ExampleID != "" {
					expected[block.ExampleID] = true
				}
			}
		}
	}
	if len(expected) != len(seen) {
		t.Fatalf("storyboard has %d examples; lessons use %d", len(seen), len(expected))
	}
	for id := range expected {
		if !seen[id] {
			t.Errorf("lesson example %s is missing from storyboard", id)
		}
	}
}

func ch22VerifyExample(t *testing.T, goBin, storyboard string, example []string, levelIDs []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unsupported outcome %q", id, outcome)
	}
	var matches int
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						matches++
						if block.Code != code || strings.Join(block.Imports, ",") != imports {
							t.Fatalf("%s: storyboard differs from lesson", id)
						}
						if block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != ch22StoryboardOutput(t, storyboard, id) {
							t.Fatalf("%s: terminal frame differs from transcript", id)
						}
					}
				}
			}
		}
	}
	if matches != 2 {
		t.Fatalf("%s: expected matching Arabic and English walkthroughs, got %d", id, matches)
	}
	ch22ExecuteExample(t, goBin, id, imports, code, ch22StoryboardOutput(t, storyboard, id))
}

func ch22StoryboardOutput(t *testing.T, storyboard, id string) string {
	t.Helper()
	parts := strings.SplitN(storyboard, "#### Verified "+id+" — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing verified output", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	match := regexp.MustCompile("(?s)Real stdout:\n\n```text\n(.*?)\n```").FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("%s: missing real stdout", id)
	}
	return match[1]
}

func ch22ExecuteExample(t *testing.T, goBin, id, imports, code, wantOutput string) {
	t.Helper()
	dir := t.TempDir()
	var quotedImports []string
	for _, importName := range strings.Split(imports, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", importName))
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + strings.Join(quotedImports, ";") + ")\n" + code + "\n",
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
		t.Fatalf("%s build: %v\n%s", id, err, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	var stdout, stderr strings.Builder
	run.Stdout, run.Stderr = &stdout, &stderr
	err = run.Run()
	if err != nil || stderr.Len() != 0 || stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s output mismatch: got %q, want %q, err=%v stderr=%q", id, stdout.String(), wantOutput+"\n", err, stderr.String())
	}
}
