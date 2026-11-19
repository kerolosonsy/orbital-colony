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

func TestCh26Lessons(t *testing.T) {
	ids := []string{
		"26-1-first-server",
		"26-2-rest-crud",
		"26-3-router",
		"26-4-envelope",
		"26-5-static",
		"26-6-upload",
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seenArabic, seenEnglish := map[string]bool{}, map[string]bool{}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English text", id)
		}
		if level.Kind != "lab" {
			t.Errorf("%s: expected a lab", id)
		}
		if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: lesson revision must be 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		prefix := fmt.Sprintf("c26%d.", mission+1)
		ch26ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, prefix, seenArabic)
		ch26ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, prefix, seenEnglish)
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
					continue
				}
				if block.Question != nil {
					if other.Question == nil || other.Question.Answer != block.Question.Answer ||
						len(other.Question.Options) != len(block.Question.Options) ||
						len(other.Question.Feedback) != len(block.Question.Feedback) {
						t.Errorf("%s: question structure differs between languages", block.ID)
					}
				}
				if block.Walkthrough != nil {
					if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID ||
						!reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
						t.Errorf("%s: walkthrough structure differs between languages", block.ID)
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
			}
		}
	}
	ch26CheckGlobalIDs(t, seenArabic, seenEnglish)
}

func ch26CheckGlobalIDs(t *testing.T, seenArabic, seenEnglish map[string]bool) {
	t.Helper()
	for _, level := range All {
		if level.Chapter == 26 {
			continue
		}
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch26BlockIDs(step.Blocks)...) {
				if seenArabic[id] {
					t.Errorf("duplicate global Arabic lesson ID %q", id)
				}
				seenArabic[id] = true
			}
		}
		if level.En != nil {
			for _, step := range level.En.Steps {
				for _, id := range append([]string{step.ID}, ch26BlockIDs(step.Blocks)...) {
					if seenEnglish[id] {
						t.Errorf("duplicate global English lesson ID %q", id)
					}
					seenEnglish[id] = true
				}
			}
		}
	}
}

func ch26BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, 0, len(blocks))
	for _, block := range blocks {
		ids = append(ids, block.ID)
	}
	return ids
}

func ch26ValidateCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string, seen map[string]bool) {
	t.Helper()
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
			t.Errorf("%s: step %d has the wrong title or no blocks", name, i)
		}
		words += ch26TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch26TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 ||
					len(block.Question.Options) != len(block.Question.Feedback) ||
					block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", name, block.ID)
					continue
				}
				words += ch26TextWords(strings.Join(block.Question.Options, " "))
				words += ch26TextWords(strings.Join(block.Question.Feedback, " "))
				for option := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[option]) == "" ||
						strings.TrimSpace(block.Question.Feedback[option]) == "" {
						t.Errorf("%s: %s has an empty option or feedback", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch26ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints, got %d", name, len(hints))
	} else {
		for i, hint := range hints {
			lower := strings.ToLower(strings.TrimSpace(hint))
			for _, label := range []string{"concept:", "structure:", "specific:", "location:", "المفهوم:", "الهيكل:", "الموضع:", "الفكرة:"} {
				if strings.HasPrefix(lower, label) {
					t.Errorf("%s: hint %d starts with a position label", name, i+1)
				}
			}
			if strings.TrimSpace(hint) == "" {
				t.Errorf("%s: hint %d is empty", name, i+1)
			}
			words += ch26TextWords(hint)
		}
	}
	if len(steps) == 6 {
		try := steps[2]
		if len(try.Blocks) != 2 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" {
			t.Errorf("%s: Try it must predict before revealing with a walkthrough", name)
		} else if strings.TrimSpace(ch26QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction code differs from its walkthrough", name)
		}
		if len(steps[1].Blocks) != 1 || steps[1].Blocks[0].Kind != "walkthrough" {
			t.Errorf("%s: Idea must be a walkthrough", name)
		}
		if len(steps[4].Blocks) != 2 || steps[4].Blocks[1].Kind != "question" {
			t.Errorf("%s: Think & check needs a question", name)
		}
	}
	if words < 200 || words > 320 {
		t.Errorf("%s learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch26QuestionCode(text string) string {
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

func ch26ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
			t.Errorf("%s: execution continues after panic", block.ID)
		}
		words += len(strings.Fields(frame.Caption))
	}
	return words
}

func ch26TextWords(text string) int {
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

func TestCh26StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 26 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH26.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e26[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n\\x60\\x60\\x60go\\n(.*?)\\n\\x60\\x60\\x60")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 12 {
		t.Fatalf("expected twelve chapter 26 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) {
			ch26VerifyExample(t, goBin, storyboard, example)
		})
	}
}

func ch26VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"26-1-first-server", "26-2-rest-crud", "26-3-router", "26-4-envelope", "26-5-static", "26-6-upload"} {
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
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n\\x60\\x60\\x60text\\n(.*?)\\n\\x60\\x60\\x60")
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
	for _, levelID := range []string{"26-1-first-server", "26-2-rest-crud", "26-3-router", "26-4-envelope", "26-5-static", "26-6-upload"} {
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
