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

func TestCh08LessonsAreBilingualAndValid(t *testing.T) {
	missions := []struct {
		id, prefix string
		arabic     []string
		english    []string
	}{
		{"8-1-copy-bug", "c81.", []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}, []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}},
		{"8-2-cargo-runs", "c82.", []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}, []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}},
		{"8-3-nil-radar", "c83.", []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}, []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}},
		{"8-4-cart-chain", "c84.", []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}, []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}},
	}
	for _, mission := range missions {
		t.Run(mission.id, func(t *testing.T) {
			level := Get(mission.id)
			if level == nil || level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("missing bilingual revision 1 lesson")
			}
			if len(level.Steps) != 6 || len(level.En.Steps) != 6 || len(level.Hints) != 3 || len(level.En.Hints) != 3 {
				t.Fatal("expected six steps and three hints in both languages")
			}
			ch08ValidateHints(t, level.Hints, []string{"الفكرة:", "الترتيب:", "المكان:"})
			ch08ValidateHints(t, level.En.Hints, []string{"Concept:", "Structure:", "Specific:"})
			ch08ValidateSteps(t, level.Steps, mission.prefix, mission.arabic)
			ch08ValidateSteps(t, level.En.Steps, mission.prefix, mission.english)
			for i, step := range level.Steps {
				translated := level.En.Steps[i]
				if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
					t.Fatalf("step %s has mismatched translation structure", step.ID)
				}
				for j, block := range step.Blocks {
					ch08ValidateBlockParity(t, block, translated.Blocks[j])
				}
			}
			counts := map[string]int{"Arabic": ch08LessonWords(level.Steps, level.Hints), "English": ch08LessonWords(level.En.Steps, level.En.Hints)}
			t.Logf("lesson words: Arabic %d, English %d", counts["Arabic"], counts["English"])
			for language, count := range counts {
				if count < 200 || count > 320 {
					t.Errorf("%s lesson word count %d is outside 200-320", language, count)
				}
			}
		})
	}
}

func ch08ValidateHints(t *testing.T, hints, labels []string) {
	t.Helper()
	for i, hint := range hints {
		if strings.HasPrefix(hint, labels[i]) || ch08SentenceCount(hint) != 1 {
			t.Fatalf("hint %d must be position-labelled and one sentence: %q", i, hint)
		}
	}
}

func ch08SentenceCount(text string) int {
	count := 0
	for _, r := range text {
		if strings.ContainsRune(".!?؟。", r) {
			count++
		}
	}
	return count
}

func ch08ValidateSteps(t *testing.T, steps []LessonStep, prefix string, titles []string) {
	t.Helper()
	seen := map[string]bool{}
	sandboxCount := 0
	for i, step := range steps {
		if step.Title != titles[i] || !strings.HasPrefix(step.ID, prefix) || seen[step.ID] || len(step.Blocks) == 0 {
			t.Fatalf("invalid step %d: %#v", i, step)
		}
		seen[step.ID] = true
		for _, block := range step.Blocks {
			if !strings.HasPrefix(block.ID, prefix) || seen[block.ID] || strings.TrimSpace(block.Text) == "" {
				t.Fatalf("invalid or duplicate block ID/text: %#v", block)
			}
			seen[block.ID] = true
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Fatalf("%s: markdown has an unexpected payload", block.ID)
				}
			case "question":
				if block.Question == nil || block.Walkthrough != nil || block.Sandbox != nil {
					t.Fatalf("%s: invalid question payload", block.ID)
				}
				q := block.Question
				if len(q.Options) < 2 || q.Answer < 0 || q.Answer >= len(q.Options) || len(q.Feedback) != len(q.Options) {
					t.Fatalf("%s: answer or feedback is out of range", block.ID)
				}
				for _, feedback := range q.Feedback {
					if ch08SentenceCount(feedback) != 1 {
						t.Fatalf("%s: option feedback must be one sentence", block.ID)
					}
				}
				if block.Code != "" && (len(block.Imports) == 0 || block.ExampleID == "") {
					t.Fatalf("%s: question example lacks imports or ID", block.ID)
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Question != nil || block.Sandbox != nil || block.Code == "" || len(block.Imports) == 0 || block.ExampleID == "" {
					t.Fatalf("%s: invalid walkthrough payload", block.ID)
				}
				ch08ValidateFrames(t, block)
			case "sandbox":
				sandboxCount++
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch08-") || len(block.Sandbox.Config) == 0 {
					t.Fatalf("%s: invalid chapter 8 sandbox payload", block.ID)
				}
			default:
				t.Fatalf("%s: unsupported block kind %q", block.ID, block.Kind)
			}
		}
		if i == 2 && (len(step.Blocks) != 3 || step.Blocks[0].Kind != "question" || step.Blocks[1].Kind != "walkthrough" || step.Blocks[2].Kind != "sandbox") {
			t.Fatalf("%s: Try it must predict, reveal, then show one sandbox", prefix)
		}
	}
	if sandboxCount != 1 {
		t.Fatalf("%s: expected exactly one chapter 8 sandbox, got %d", prefix, sandboxCount)
	}
}

func ch08ValidateFrames(t *testing.T, block LessonBlock) {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Fatalf("%s: walkthrough needs an initial frame", block.ID)
	}
	lineCount := strings.Count(block.Code, "\n") + 1
	for i, frame := range frames {
		if frame.Line < 0 || frame.Line > lineCount || (i > 0 && frame.Line == 0) || strings.TrimSpace(frame.Caption) == "" {
			t.Fatalf("%s frame %d has invalid line or caption", block.ID, i)
		}
		if len(strings.Fields(frame.Caption)) > 8 {
			t.Fatalf("%s frame %d caption exceeds eight words", block.ID, i)
		}
		if i > 0 && !strings.HasPrefix(frame.Output, frames[i-1].Output) {
			t.Fatalf("%s frame %d output is not cumulative", block.ID, i)
		}
	}
}

func ch08ValidateBlockParity(t *testing.T, arabic, english LessonBlock) {
	t.Helper()
	if arabic.ID != english.ID || arabic.Kind != english.Kind || arabic.ExampleID != english.ExampleID || arabic.Code != english.Code || strings.Join(arabic.Imports, ",") != strings.Join(english.Imports, ",") {
		t.Fatalf("%s: block identity or code differs across languages", arabic.ID)
	}
	if arabic.Question != nil {
		if english.Question == nil || arabic.Question.Answer != english.Question.Answer || len(arabic.Question.Options) != len(english.Question.Options) || len(arabic.Question.Feedback) != len(english.Question.Feedback) {
			t.Fatalf("%s: question parity differs", arabic.ID)
		}
	}
	if arabic.Walkthrough != nil {
		if english.Walkthrough == nil || len(arabic.Walkthrough.Frames) != len(english.Walkthrough.Frames) {
			t.Fatalf("%s: walkthrough frame count differs", arabic.ID)
		}
		for i, frame := range arabic.Walkthrough.Frames {
			translated := english.Walkthrough.Frames[i]
			frame.Caption, translated.Caption = "", ""
			if !reflect.DeepEqual(frame, translated) {
				t.Fatalf("%s frame %d has mismatched execution", arabic.ID, i)
			}
		}
	}
	if arabic.Sandbox != nil {
		if english.Sandbox == nil || arabic.Sandbox.Kind != english.Sandbox.Kind || !reflect.DeepEqual(arabic.Sandbox.Config, english.Sandbox.Config) {
			t.Fatalf("%s: sandbox config differs across languages", arabic.ID)
		}
	}
}

func TestCh08SandboxConfigsMatchFixtures(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch08.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]any
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}

	kinds := map[string]bool{}
	for _, id := range []string{"8-1-copy-bug", "8-2-cargo-runs", "8-3-nil-radar", "8-4-cart-chain"} {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing bilingual level", id)
		}
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			count := 0
			for stepIndex, step := range steps {
				for blockIndex, block := range step.Blocks {
					if block.Sandbox == nil {
						continue
					}
					count++
					if stepIndex != 2 || blockIndex != len(step.Blocks)-1 {
						t.Errorf("%s: sandbox must follow the Try it question and walkthrough", id)
					}
					kinds[block.Sandbox.Kind] = true
					body, err := json.Marshal(block.Sandbox.Config)
					if err != nil {
						t.Fatal(err)
					}
					var normalized map[string]any
					if err := json.Unmarshal(body, &normalized); err != nil {
						t.Fatal(err)
					}
					if !reflect.DeepEqual(normalized, fixtures[block.Sandbox.Kind]) {
						t.Errorf("%s: config for %s differs from its JS fixture", id, block.Sandbox.Kind)
					}
				}
			}
			if count != 1 {
				t.Errorf("%s: expected one sandbox in each language, got %d", id, count)
			}
		}
	}
	if len(kinds) != len(fixtures) {
		t.Errorf("found %d sandbox kinds but %d fixture entries", len(kinds), len(fixtures))
	}
}

func ch08LessonWords(steps []LessonStep, hints []string) int {
	count := 0
	add := func(text string) { count += len(strings.Fields(text)) }
	for _, step := range steps {
		add(step.Title)
		for _, block := range step.Blocks {
			add(block.Text)
			if block.Question != nil {
				for _, item := range block.Question.Options {
					add(item)
				}
				for _, item := range block.Question.Feedback {
					add(item)
				}
			}
			if block.Walkthrough != nil {
				for _, frame := range block.Walkthrough.Frames {
					add(frame.Caption)
				}
			}
		}
	}
	for _, hint := range hints {
		add(hint)
	}
	return count
}

func TestCh08StoryboardExamplesRunWithGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH08.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e8[1-4]-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 8 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] || example[3] != "ok" {
			t.Fatalf("duplicate or unsupported example %q", example[1])
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch08VerifyExample(t, goBin, storyboard, example) })
	}
	for _, id := range []string{"8-1-copy-bug", "8-2-cargo-runs", "8-3-nil-radar", "8-4-cart-chain"} {
		level := Get(id)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.Code != "" && !seen[block.ExampleID] {
						t.Errorf("%s references example %q missing from storyboard", block.ID, block.ExampleID)
					}
				}
			}
		}
	}
}

func ch08VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	var wantOutput string
	for _, steps := range [][]LessonStep{Get("8-1-copy-bug").Steps, Get("8-1-copy-bug").En.Steps, Get("8-2-cargo-runs").Steps, Get("8-2-cargo-runs").En.Steps, Get("8-3-nil-radar").Steps, Get("8-3-nil-radar").En.Steps, Get("8-4-cart-chain").Steps, Get("8-4-cart-chain").En.Steps} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID != id {
					continue
				}
				if block.Code != code || strings.Join(block.Imports, ",") != imports {
					t.Fatalf("%s: lesson example differs from storyboard", id)
				}
				if block.Walkthrough != nil {
					frames := block.Walkthrough.Frames
					wantOutput = frames[len(frames)-1].Output
				}
			}
		}
	}
	sections := strings.SplitN(storyboard, "#### Verified `"+id+"`", 2)
	if len(sections) != 2 {
		t.Fatalf("%s: missing verified output", id)
	}
	section := strings.Split(sections[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(section)
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	documented := match[1] + "\n"
	if strings.TrimSuffix(documented, "\n") != wantOutput {
		t.Fatalf("%s: walkthrough output %q differs from storyboard %q", id, wantOutput, documented)
	}
	if outcome != "ok" {
		t.Fatalf("%s: unsupported outcome %q", id, outcome)
	}

	dir := t.TempDir()
	quotedImports := []string{}
	for _, name := range strings.Split(imports, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", name))
	}
	program := "package main\nimport (" + strings.Join(quotedImports, ";") + ")\n"
	if strings.Contains(code, "func main()") {
		program += code + "\n"
	} else {
		program += "func main() {\n" + code + "\n}\n"
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": program,
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
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	if err := run.Run(); err != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, err, stderr)
	}
	if stdout.String() != documented {
		t.Fatalf("%s: stdout %q differs from documented %q", id, stdout.String(), documented)
	}
}
