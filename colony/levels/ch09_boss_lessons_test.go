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

func TestCh09BossLessons(t *testing.T) {
	level := Get("9-boss-evacuation")
	if level == nil || level.En == nil {
		t.Fatal("boss lesson or English text is missing")
	}
	if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
		t.Fatal("Arabic and English lesson revisions must both be 1")
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	ch09BossValidateCopy(t, "Arabic", level.Steps, level.Hints, arabicTitles)
	ch09BossValidateCopy(t, "English", level.En.Steps, level.En.Hints, englishTitles)
	if len(level.Steps) != len(level.En.Steps) {
		t.Fatal("Arabic and English lesson step counts differ")
	}
	for i, step := range level.Steps {
		translated := level.En.Steps[i]
		if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
			t.Errorf("step %d: Arabic and English identities differ", i+1)
			continue
		}
		for j, block := range step.Blocks {
			ch09BossAssertParity(t, block, translated.Blocks[j])
		}
	}
	ch09BossAssertGlobalIDs(t, level)
}

func ch09BossValidateCopy(t *testing.T, language string, steps []LessonStep, hints, titles []string) {
	t.Helper()
	if len(steps) != 6 || len(titles) != 6 || len(hints) != 3 {
		t.Fatalf("%s: got %d steps and %d hints; want six and three", language, len(steps), len(hints))
	}
	seen := map[string]bool{}
	sandboxes := 0
	claim := func(id string) {
		if !strings.HasPrefix(id, "c9b.") || seen[id] {
			t.Errorf("%s: empty-prefix or duplicate ID %q", language, id)
		}
		seen[id] = true
	}
	words := 0
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has the wrong title or no blocks", language, i+1)
		}
		words += ch09BossVisibleWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", language, block.ID)
			}
			words += ch09BossVisibleWords(block.Text)
			if block.Kind == "sandbox" {
				sandboxes++
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch09-") {
					t.Errorf("%s: sandbox must use a ch09- kind", block.ID)
				}
			}
			for _, line := range strings.Split(block.Text, "\n") {
				if len([]rune(line)) > 140 {
					t.Errorf("%s: %s has a line over three panel lines", language, block.ID)
				}
			}
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", language, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", language, block.ID)
					continue
				}
				words += ch09BossVisibleWords(strings.Join(block.Question.Options, " "))
				words += ch09BossVisibleWords(strings.Join(block.Question.Feedback, " "))
				for i, option := range block.Question.Options {
					feedback := block.Question.Feedback[i]
					if strings.TrimSpace(option) == "" || strings.TrimSpace(feedback) == "" {
						t.Errorf("%s: %s has an empty option or feedback", language, block.ID)
					}
					if strings.ContainsAny(feedback, "\n!?؟") || !strings.HasSuffix(strings.TrimSpace(feedback), ".") {
						t.Errorf("%s: %s feedback %d must be one sentence", language, block.ID, i+1)
					}
					if len([]rune(feedback)) > 140 || len([]rune(option)) > 140 {
						t.Errorf("%s: %s option or feedback exceeds three panel lines", language, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 || block.Sandbox != nil {
					t.Errorf("%s: %s has an incomplete walkthrough", language, block.ID)
					continue
				}
				words += ch09BossValidateFrames(t, language, block)
			case "sandbox":
				if block.Sandbox == nil || block.Code != "" || block.Question != nil || block.Walkthrough != nil {
					t.Errorf("%s: sandbox has an invalid payload", block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed kind %q", language, block.ID, block.Kind)
			}
		}
	}
	if sandboxes != 1 {
		t.Errorf("%s: got %d sandboxes, want exactly one", language, sandboxes)
	}
	for i, hint := range hints {
		if strings.TrimSpace(hint) == "" || strings.ContainsAny(hint, "\n!?؟") || !strings.HasSuffix(strings.TrimSpace(hint), ".") {
			t.Errorf("%s: hint %d must be one sentence", language, i+1)
		}
		lower := strings.ToLower(strings.TrimSpace(hint))
		for _, prefix := range []string{"concept:", "structure:", "specific:", "الفكرة:", "الهيكل:", "الموضع:"} {
			if strings.HasPrefix(lower, prefix) {
				t.Errorf("%s: hint %d starts with a position label", language, i+1)
			}
		}
		words += ch09BossVisibleWords(hint)
	}
	if len(steps) >= 3 {
		try := steps[2].Blocks
		if len(try) != 3 || try[0].Kind != "question" || try[1].Kind != "walkthrough" || try[2].Kind != "sandbox" {
			t.Errorf("%s: Try it must predict, reveal, then offer one sandbox", language)
		} else if strings.TrimSpace(ch09BossQuestionCode(try[0].Text)) != strings.TrimSpace(try[1].Code) {
			t.Errorf("%s: prediction code differs from the walkthrough", language)
		}
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200–320", language, words)
	}
	t.Logf("%s learner-authored words: %d", language, words)
}

func ch09BossValidateFrames(t *testing.T, language string, block LessonBlock) int {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Errorf("%s: %s has no initial frame", language, block.ID)
		return 0
	}
	maxLine := strings.Count(block.Code, "\n") + 1
	words := 0
	for i, frame := range frames {
		if frame.Line < 0 || frame.Line > maxLine || (i > 0 && frame.Line == 0) || strings.TrimSpace(frame.Caption) == "" {
			t.Errorf("%s: %s frame %d points outside its code", language, block.ID, i)
		}
		if len(strings.Fields(frame.Caption)) > 8 {
			t.Errorf("%s: %s frame %d caption exceeds eight words", language, block.ID, i)
		}
		if i > 0 && !strings.HasPrefix(frame.Output, frames[i-1].Output) {
			t.Errorf("%s: %s frame %d output is not cumulative", language, block.ID, i)
		}
		if frame.Panic != "" && i != len(frames)-1 {
			t.Errorf("%s: %s continues after a panic", language, block.ID)
		}
		words += ch09BossVisibleWords(frame.Caption)
	}
	return words
}

func ch09BossVisibleWords(text string) int {
	var visible strings.Builder
	inFence, inInlineCode := false, false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		for _, r := range line {
			if r == '~' {
				inInlineCode = !inInlineCode
				continue
			}
			if !inInlineCode {
				visible.WriteRune(r)
			}
		}
		visible.WriteByte(' ')
	}
	return len(strings.Fields(visible.String()))
}

func ch09BossQuestionCode(text string) string {
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

func ch09BossAssertParity(t *testing.T, arabic, english LessonBlock) {
	t.Helper()
	if arabic.ID != english.ID || arabic.Kind != english.Kind {
		t.Errorf("%s: block identities differ across languages", arabic.ID)
	}
	if arabic.Question != nil {
		if english.Question == nil || arabic.Question.Answer != english.Question.Answer || len(arabic.Question.Options) != len(english.Question.Options) || len(arabic.Question.Feedback) != len(english.Question.Feedback) {
			t.Errorf("%s: question answer or option count differs across languages", arabic.ID)
		}
	}
	if arabic.Walkthrough != nil {
		if english.Walkthrough == nil || arabic.Code != english.Code || arabic.ExampleID != english.ExampleID || !reflect.DeepEqual(arabic.Imports, english.Imports) || len(arabic.Walkthrough.Frames) != len(english.Walkthrough.Frames) {
			t.Errorf("%s: walkthrough structures differ across languages", arabic.ID)
			return
		}
		for i, frame := range arabic.Walkthrough.Frames {
			other := english.Walkthrough.Frames[i]
			frame.Caption, other.Caption = "", ""
			if !reflect.DeepEqual(frame, other) {
				t.Errorf("%s frame %d: execution states differ", arabic.ID, i)
			}
		}
	}
	if arabic.Sandbox != nil {
		if english.Sandbox == nil || arabic.Sandbox.Kind != english.Sandbox.Kind {
			t.Errorf("%s: sandbox identity differs across languages", arabic.ID)
		} else {
			sandboxConfigParity(t, arabic.Sandbox.Config, english.Sandbox.Config)
		}
	}
}

func ch09BossAssertGlobalIDs(t *testing.T, target *Level) {
	t.Helper()
	targetIDs := map[string]bool{}
	for _, step := range target.Steps {
		targetIDs[step.ID] = true
		for _, block := range step.Blocks {
			targetIDs[block.ID] = true
		}
	}
	for _, level := range All {
		if level.ID == target.ID {
			continue
		}
		copies := [][]LessonStep{level.Steps}
		if level.En != nil {
			copies = append(copies, level.En.Steps)
		}
		for _, steps := range copies {
			for _, step := range steps {
				if targetIDs[step.ID] {
					t.Errorf("lesson ID %q collides with level %s", step.ID, level.ID)
				}
				for _, block := range step.Blocks {
					if targetIDs[block.ID] {
						t.Errorf("lesson ID %q collides with level %s", block.ID, level.ID)
					}
				}
			}
		}
	}
}

func TestCh09BossStoryboardExampleMatchesRealGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute the boss example")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH09.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e9boss-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 1 {
		t.Fatalf("expected one boss example, found %d", len(examples))
	}
	example := examples[0]
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if id != "e9boss-map-count" || outcome != "ok" {
		t.Fatalf("unexpected example metadata: %s %s", id, outcome)
	}
	level := Get("9-boss-evacuation")
	found := 0
	for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == id {
					found++
					if block.Code != code || strings.Join(block.Imports, ",") != importNames {
						t.Fatalf("%s differs from the structured lesson example", id)
					}
				}
			}
		}
	}
	if found != 2 {
		t.Fatalf("%s occurs %d times in structured lessons, want bilingual pair", id, found)
	}
	parts := strings.SplitN(storyboard, "#### Verified `"+id+"` — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("missing verified output for %s", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("missing stdout for %s", id)
	}
	wantOutput := match[1] + "\n"
	for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != strings.TrimSuffix(wantOutput, "\n") {
					t.Fatalf("%s final walkthrough output differs from documented output", id)
				}
			}
		}
	}
	quotedImports := []string{}
	for _, name := range strings.Split(importNames, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", name))
	}
	dir := t.TempDir()
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + strings.Join(quotedImports, ";") + ")\nfunc main() {\n" + code + "\n}\n",
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
	if err := run.Run(); err != nil || stderr.Len() != 0 || stdout.String() != wantOutput {
		t.Fatalf("run %s: err=%v stdout=%q want=%q stderr=%q", id, err, stdout.String(), wantOutput, stderr.String())
	}
}
