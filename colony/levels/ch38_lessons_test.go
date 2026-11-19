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

func TestCh38Lessons(t *testing.T) {
	ids := []string{
		"38-1-iterators",
		"38-2-sets-and-labels",
		"38-3-once",
		"38-4-errgroup",
		"38-5-allocs",
		"38-6-vet",
		"38-7-pprof",
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	for mission, id := range ids {
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
		ch38ValidateLessonCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, fmt.Sprintf("c38%d.", mission+1))
		ch38ValidateLessonCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, fmt.Sprintf("c38%d.", mission+1))
		ch38CheckGlobalIDs(t, level)
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
						t.Errorf("%s: question answers or option counts differ", block.ID)
					}
				}
				if block.Walkthrough != nil {
					if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
						t.Errorf("%s: walkthrough structures differ", block.ID)
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
}

func ch38CheckGlobalIDs(t *testing.T, level *Level) {
	t.Helper()
	own := map[string]bool{}
	for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
		for _, step := range steps {
			own[step.ID] = true
			for _, block := range step.Blocks {
				own[block.ID] = true
			}
		}
	}
	for _, other := range All {
		if other.Chapter == level.Chapter || other.En == nil {
			continue
		}
		for _, steps := range [][]LessonStep{other.Steps, other.En.Steps} {
			for _, step := range steps {
				if own[step.ID] {
					t.Errorf("%s: lesson ID %q is also used by %s", level.ID, step.ID, other.ID)
				}
				for _, block := range step.Blocks {
					if own[block.ID] {
						t.Errorf("%s: lesson ID %q is also used by %s", level.ID, block.ID, other.ID)
					}
				}
			}
		}
	}
}

func ch38ValidateLessonCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string) {
	t.Helper()
	seen := map[string]bool{}
	claim := func(id string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q has wrong mission prefix; want %q", name, id, prefix)
		}
		seen[id] = true
	}
	words := 0
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i)
		}
		words += ch38TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch38TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown contains another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback or answer", name, block.ID)
					continue
				}
				words += ch38TextWords(strings.Join(block.Question.Options, " "))
				words += ch38TextWords(strings.Join(block.Question.Feedback, " "))
				for i := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[i]) == "" || strings.TrimSpace(block.Question.Feedback[i]) == "" {
						t.Errorf("%s: %s has empty option or feedback", name, block.ID)
					}
					feedback := strings.TrimSpace(block.Question.Feedback[i])
					if strings.Count(feedback, ".") != 1 || !strings.HasSuffix(feedback, ".") {
						t.Errorf("%s: %s feedback %d is not one sentence", name, block.ID, i)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch38ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
	} else {
		for i, hint := range hints {
			if strings.TrimSpace(hint) == "" || ch38HasHintPrefix(hint) || !strings.HasSuffix(strings.TrimSpace(hint), ".") {
				t.Errorf("%s: hint %d is empty, prefixed or not one sentence", name, i+1)
			}
			words += ch38TextWords(hint)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 2 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" {
			t.Errorf("%s: Try it must predict first, then reveal with a walkthrough", name)
		} else if strings.TrimSpace(ch38QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction code differs from its walkthrough reveal\nquestion=%q\nwalkthrough=%q", name, ch38QuestionCode(try.Blocks[0].Text), try.Blocks[1].Code)
		}
	}
	if len(steps) > 1 && (len(steps[1].Blocks) != 1 || steps[1].Blocks[0].Kind != "walkthrough") {
		t.Errorf("%s: Idea must show a walkthrough", name)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch38HasHintPrefix(hint string) bool {
	hint = strings.TrimSpace(hint)
	for _, prefix := range []string{"Concept:", "Structure:", "Location:", "الفكرة:", "الهيكل:", "الموضع:", "المفهوم:"} {
		if strings.HasPrefix(hint, prefix) {
			return true
		}
	}
	return false
}

func ch38ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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

func ch38QuestionCode(text string) string {
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

func ch38TextWords(text string) int {
	visible := make([]string, 0)
	inFence := false
	for _, line := range strings.Split(text, "\n") {
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

func TestCh38StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 38 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH38.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e38[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 14 {
		t.Fatalf("expected fourteen chapter 38 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch38VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch38VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected storyboard outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"38-1-iterators", "38-2-sets-and-labels", "38-3-once", "38-4-errgroup", "38-5-allocs", "38-6-vet", "38-7-pprof"} {
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
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
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
	goMod := "module storyboard/" + id + "\n\ngo 1.24\n"
	if strings.Contains(importNames, "golang.org/x/sync/errgroup") {
		goMod += "\nrequire golang.org/x/sync v0.25.0\n"
	}
	files := map[string]string{
		"go.mod":  goMod,
		"main.go": "package main\nimport (" + strings.Join(imports, ";") + ")\n" + code + "\n",
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, goBin, "build", "-mod=mod", "-trimpath", "-gcflags=-lang=go1.24", "-o", "example", ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s: %v\n%s", id, err, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	var stdout, stderr strings.Builder
	run.Stdout, run.Stderr = &stdout, &stderr
	if err := run.Run(); err != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, err, stderr.String())
	}
	if stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), wantOutput+"\n")
	}
	for _, levelID := range []string{"38-1-iterators", "38-2-sets-and-labels", "38-3-once", "38-4-errgroup", "38-5-allocs", "38-6-vet", "38-7-pprof"} {
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
