package levels

import (
	"context"
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

func TestCh20Lessons(t *testing.T) {
	ids := []string{"20-1-bulletin", "20-2-validators", "20-3-lexer", "20-4-evaluator"}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seenIDs := map[string]bool{}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English copy", id)
		}
		if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: expected lab and lesson revision 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		ch20ValidateLessonCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, mission+1, 200, 320, seenIDs)
		englishIDs := map[string]bool{}
		ch20ValidateLessonCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, mission+1, 200, 320, englishIDs)
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
					if other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options) || len(block.Question.Feedback) != len(other.Question.Feedback) {
						t.Errorf("%s: question answer or option structure differs", block.ID)
					}
				}
				if block.Walkthrough != nil {
					if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
						t.Errorf("%s: walkthrough structure differs", block.ID)
						continue
					}
					for frameIndex, frame := range block.Walkthrough.Frames {
						translatedFrame := other.Walkthrough.Frames[frameIndex]
						frame.Caption, translatedFrame.Caption = "", ""
						if !reflect.DeepEqual(frame, translatedFrame) {
							t.Errorf("%s frame %d: execution state differs", block.ID, frameIndex)
						}
					}
				}
			}
		}
	}
}

func ch20ValidateLessonCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, mission, minWords, maxWords int, seen map[string]bool) {
	t.Helper()
	claim := func(id string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, fmt.Sprintf("c20%d.", mission)) {
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
		words += ch20TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch20TextWords(block.Text)
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
				words += ch20TextWords(strings.Join(block.Question.Options, " "))
				words += ch20TextWords(strings.Join(block.Question.Feedback, " "))
				for k := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[k]) == "" || strings.TrimSpace(block.Question.Feedback[k]) == "" {
						t.Errorf("%s: %s has empty option or feedback", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch20ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints, got %d", name, len(hints))
	} else {
		for i, hint := range hints {
			for _, label := range []string{"Concept:", "Structure:", "Location:", "المفهوم:", "الهيكل:", "الموضع:"} {
				if strings.HasPrefix(hint, label) {
					t.Errorf("%s: hint %d starts with forbidden label %q", name, i+1, label)
				}
			}
			words += ch20TextWords(hint)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 2 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" {
			t.Errorf("%s: Try it must predict first, then reveal with a walkthrough", name)
		} else if strings.TrimSpace(ch20QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction example differs from its walkthrough reveal", name)
		}
	}
	if words < minWords || words > maxWords {
		t.Errorf("%s learner text has %d words; want %d-%d", name, words, minWords, maxWords)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch20QuestionCode(text string) string {
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

func ch20ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
		words += len(strings.Fields(frame.Caption))
	}
	return words
}

func ch20TextWords(text string) int {
	lines := strings.Split(text, "\n")
	inFence := false
	visible := make([]string, 0, len(lines))
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

func TestCh20StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 20 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH20.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e20[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 20 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch20VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch20VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected storyboard outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"20-1-bulletin", "20-2-validators", "20-3-lexer", "20-4-evaluator"} {
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
		t.Fatalf("%s: no chapter 20 walkthrough uses this example", id)
	}
	parts := strings.SplitN(storyboard, "#### Verified "+id+" — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing documented output", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
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
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, goBin, "run", ".")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: go run failed: %v\n%s", id, err, output)
	}
	if strings.TrimSuffix(string(output), "\n") != match[1] {
		t.Fatalf("%s: stdout = %q, want %q", id, output, match[1])
	}
}
