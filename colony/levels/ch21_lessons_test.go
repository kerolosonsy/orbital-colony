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

func TestCh21Lessons(t *testing.T) {
	ch21CheckGlobalLessonIDs(t)
	ids := []string{"21-1-sysinfo", "21-2-subcommands", "21-3-guess-depth", "21-4-duel"}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seenIDs := map[string]string{}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English text", id)
		}
		if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: expected lab lesson revision 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		ch21ValidateLessonCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, mission+1, seenIDs)
		ch21ValidateLessonCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, mission+1, nil)
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
			}
		}
	}
}

func ch21CheckGlobalLessonIDs(t *testing.T) {
	t.Helper()
	owners := map[string]string{}
	claim := func(id, levelID string) {
		if id == "" {
			t.Errorf("%s has an empty lesson ID", levelID)
			return
		}
		if previous, exists := owners[id]; exists {
			t.Errorf("lesson ID %q is shared by %s and %s", id, previous, levelID)
		}
		owners[id] = levelID
	}
	for _, level := range All {
		for _, step := range level.Steps {
			claim(step.ID, level.ID)
			for _, block := range step.Blocks {
				claim(block.ID, level.ID)
			}
		}
	}
}

func ch21ValidateLessonCopy(t *testing.T, name string, steps []LessonStep, hints []string, titles []string, mission int, seenIDs map[string]string) {
	t.Helper()
	localIDs := map[string]bool{}
	claim := func(id string) {
		prefix := fmt.Sprintf("c21%d.", mission)
		if id == "" || !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q has the wrong mission prefix", name, id)
		}
		if localIDs[id] {
			t.Errorf("%s: duplicate ID %q", name, id)
		}
		localIDs[id] = true
		if seenIDs != nil {
			if previous, exists := seenIDs[id]; exists {
				t.Errorf("%s: ID %q duplicates %s", name, id, previous)
			}
			seenIDs[id] = name
		}
	}
	words := 0
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i+1)
		}
		words += ch21TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch21TextWords(block.Text)
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
				words += ch21TextWords(strings.Join(block.Question.Options, " "))
				for i, feedback := range block.Question.Feedback {
					if strings.TrimSpace(block.Question.Options[i]) == "" || strings.TrimSpace(feedback) == "" {
						t.Errorf("%s: %s has empty option or feedback", name, block.ID)
					}
					words += ch21TextWords(feedback)
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch21ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected exactly three hints", name)
	} else {
		for i, hint := range hints {
			if strings.TrimSpace(hint) == "" || ch21HintHasLabel(hint) {
				t.Errorf("%s: hint %d is empty or starts with a label", name, i+1)
			}
			words += ch21TextWords(hint)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 2 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" {
			t.Errorf("%s: Try it must predict first, then reveal with a walkthrough", name)
		} else if strings.TrimSpace(ch21QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction example differs from its walkthrough reveal", name)
		}
	}
	if words < 200 || words > 320 {
		t.Errorf("%s learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch21HintHasLabel(hint string) bool {
	hint = strings.TrimSpace(hint)
	for _, label := range []string{"Concept:", "Structure:", "Location:", "الفكرة:", "المفهوم:", "الهيكل:", "الموضع:"} {
		if strings.HasPrefix(hint, label) {
			return true
		}
	}
	return false
}

func ch21ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
		words += ch21TextWords(frame.Caption)
	}
	return words
}

func ch21QuestionCode(text string) string {
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

func ch21TextWords(text string) int {
	var visible []string
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

func TestCh21StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 21 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH21.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e21[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n\\x60\\x60\\x60go\\n(.*?)\\n\\x60\\x60\\x60")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 21 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) {
			ch21VerifyExample(t, goBin, storyboard, example)
		})
	}
	for _, levelID := range []string{"21-1-sysinfo", "21-2-subcommands", "21-3-guess-depth", "21-4-duel"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.Walkthrough != nil && !seen[block.ExampleID] {
						t.Errorf("%s: walkthrough example is missing from the storyboard", block.ExampleID)
					}
				}
			}
		}
	}
}

func ch21VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected storyboard outcome %q", id, outcome)
	}
	found := 0
	wantOutput := ""
	for _, levelID := range []string{"21-1-sysinfo", "21-2-subcommands", "21-3-guess-depth", "21-4-duel"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						found++
						if block.Code != code {
							t.Fatalf("%s: lesson code differs from storyboard", id)
						}
						if !reflect.DeepEqual(block.Imports, strings.Split(importNames, ",")) {
							t.Fatalf("%s: lesson imports differ from storyboard", id)
						}
						frames := block.Walkthrough.Frames
						if len(frames) > 0 {
							wantOutput = frames[len(frames)-1].Output
						}
					}
				}
			}
		}
	}
	if found != 2 {
		t.Fatalf("%s: found %d matching bilingual walkthroughs; want 2", id, found)
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
	if match[1] != wantOutput {
		t.Fatalf("%s: final walkthrough output %q differs from storyboard %q", id, wantOutput, match[1])
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
		t.Fatalf("%s stdout = %q, want %q", id, stdout.String(), wantOutput+"\n")
	}
}
