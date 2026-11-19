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
	"unicode"
)

func TestCh04Lessons(t *testing.T) {
	ids := []string{"4-1-ripe-only", "4-2-terrain", "4-3-checkpoint"}
	prefixes := []string{"c41.", "c42.", "c43."}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seen := map[string]bool{}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Fatalf("%s: missing bilingual revision 1 lesson", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Fatalf("%s: expected six steps in both languages", id)
		}
		ch04ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, prefixes[mission], seen)
		ch04ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, prefixes[mission], nil)
		englishJSON := fmt.Sprint(level.En.Steps)
		if strings.ContainsFunc(englishJSON, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) {
			t.Errorf("%s: Arabic text appears in English steps", id)
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
				if block.Question != nil && (other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options)) {
					t.Errorf("%s: question answers or option counts differ", block.ID)
				}
				if block.Walkthrough != nil {
					if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
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

func TestCh04LessonIDsAreGloballyUnique(t *testing.T) {
	seen := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch04BlockIDs(step.Blocks)...) {
				if previous, exists := seen[id]; exists {
					t.Errorf("lesson ID %q appears in both %s and %s", id, previous, level.ID)
				}
				seen[id] = level.ID
			}
		}
	}
}

func ch04BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, 0, len(blocks))
	for _, block := range blocks {
		ids = append(ids, block.ID)
	}
	return ids
}

func ch04ValidateCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string, seen map[string]bool) {
	t.Helper()
	localIDs := map[string]bool{}
	claim := func(id string) {
		if id == "" || localIDs[id] || !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: invalid, duplicate, or wrongly prefixed ID %q", name, id)
		}
		localIDs[id] = true
		if seen != nil {
			if seen[id] {
				t.Errorf("%s: ID %q is already used by another chapter 4 level", name, id)
			}
			seen[id] = true
		}
	}
	words, walkthroughs, questions := 0, 0, 0
	for i, step := range steps {
		claim(step.ID)
		if i >= len(titles) || step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has a wrong title or no blocks", name, i+1)
		}
		words += ch04TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch04TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown contains another payload", name, block.ID)
				}
			case "question":
				questions++
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", name, block.ID)
					continue
				}
				for i := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[i]) == "" || strings.TrimSpace(block.Question.Feedback[i]) == "" {
						t.Errorf("%s: %s has empty option or feedback", name, block.ID)
					}
					if ch04SentenceCount(block.Question.Feedback[i]) != 1 {
						t.Errorf("%s: %s feedback %d must be one sentence", name, block.ID, i+1)
					}
				}
				words += ch04TextWords(strings.Join(block.Question.Options, " "))
				words += ch04TextWords(strings.Join(block.Question.Feedback, " "))
			case "walkthrough":
				walkthroughs++
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch04ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(steps) > 0 {
		mission := steps[0]
		if len(mission.Blocks) != 1 || !strings.HasPrefix(strings.TrimSpace(mission.Blocks[0].Text), "1.") || !strings.Contains(mission.Blocks[0].Text, "\n2.") || !strings.Contains(mission.Blocks[0].Text, "\n3.") || ch04TextWords(mission.Blocks[0].Text) > 45 {
			t.Errorf("%s: mission must use three brief numbered points", name)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 2 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" {
			t.Errorf("%s: Try it must predict first, then reveal with a walkthrough", name)
		} else if strings.TrimSpace(ch04QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction code differs from its reveal", name)
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
	} else {
		for i, hint := range hints {
			for _, label := range []string{"Concept:", "Structure:", "Location:", "المفهوم:", "الهيكل:", "الموضع:"} {
				if strings.HasPrefix(strings.TrimSpace(hint), label) {
					t.Errorf("%s: hint %d has a manual role prefix", name, i+1)
				}
			}
			if ch04SentenceCount(hint) != 1 {
				t.Errorf("%s: hint %d must be one sentence", name, i+1)
			}
			words += ch04TextWords(hint)
		}
	}
	if walkthroughs != 2 || questions != 2 {
		t.Errorf("%s: expected two walkthroughs and two questions, got %d and %d", name, walkthroughs, questions)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch04ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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

func ch04QuestionCode(text string) string {
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

func ch04TextWords(text string) int {
	lines := strings.Split(text, "\n")
	var visible strings.Builder
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence {
			visible.WriteString(line)
			visible.WriteByte('\n')
		}
	}
	return len(strings.Fields(ch04StripInline(visible.String())))
}

func ch04StripInline(text string) string {
	var plain strings.Builder
	inCode := false
	for _, r := range text {
		if r == '~' {
			inCode = !inCode
			continue
		}
		if !inCode {
			plain.WriteRune(r)
		}
	}
	return plain.String()
}

func ch04SentenceCount(text string) int {
	plain := ch04StripInline(text)
	return strings.Count(plain, ".") + strings.Count(plain, "!") + strings.Count(plain, "?") + strings.Count(plain, "؟")
}

func TestCh04StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 4 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH04.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e4[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 6 {
		t.Fatalf("expected six chapter 4 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch04VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch04VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"4-1-ripe-only", "4-2-terrain", "4-3-checkpoint"} {
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
		t.Fatalf("%s: no walkthrough uses this example", id)
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
	for _, levelID := range []string{"4-1-ripe-only", "4-2-terrain", "4-3-checkpoint"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != wantOutput {
						t.Fatalf("%s: terminal frame differs from transcript", id)
					}
				}
			}
		}
	}
	dir := t.TempDir()
	quotedImports := make([]string, 0)
	for _, name := range strings.Split(imports, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", name))
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
}
