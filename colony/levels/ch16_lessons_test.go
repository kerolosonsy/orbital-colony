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

var ch16ExamplePattern = regexp.MustCompile("(?s)<!-- go-example id=(e16[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
var ch16StdoutPattern = regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
var ch16InlineCodePattern = regexp.MustCompile("~[^~]*~")

func TestCh16Lessons(t *testing.T) {
	cases := []struct {
		id, prefix, titleAR, titleEN string
	}{
		{"16-1-mutants", "c161.", "المهمة|الفكرة|جرّب|خد بالك|فكّر واختبر|دورك", "Mission|Idea|Try it|Watch out|Think & check|Your turn"},
		{"16-2-bug-hunt", "c162.", "المهمة|الفكرة|جرّب|خد بالك|فكّر واختبر|دورك", "Mission|Idea|Try it|Watch out|Think & check|Your turn"},
		{"16-3-benchmark", "c163.", "المهمة|الفكرة|جرّب|خد بالك|فكّر واختبر|دورك", "Mission|Idea|Try it|Watch out|Think & check|Your turn"},
		{"16-4-wild-inputs", "c164.", "المهمة|الفكرة|جرّب|خد بالك|فكّر واختبر|دورك", "Mission|Idea|Try it|Watch out|Think & check|Your turn"},
		{"16-5-fake-clock", "c165.", "المهمة|الفكرة|جرّب|خد بالك|فكّر واختبر|دورك", "Mission|Idea|Try it|Watch out|Think & check|Your turn"},
	}
	chapterIDs := make(map[string]bool, len(cases))
	for _, item := range cases {
		chapterIDs[item.id] = true
	}
	chapterLevels := 0
	for _, level := range All {
		if strings.HasPrefix(level.ID, "16-") {
			chapterLevels++
			if !chapterIDs[level.ID] {
				t.Errorf("chapter 16 level %s has no sandbox assertions", level.ID)
			}
		}
	}
	if chapterLevels != len(cases) {
		t.Fatalf("expected %d chapter 16 levels, found %d", len(cases), chapterLevels)
	}
	for _, item := range cases {
		t.Run(item.id, func(t *testing.T) {
			level := Get(item.id)
			if level == nil || level.En == nil {
				t.Fatal("missing bilingual level")
			}
			if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("lesson revision must be 1 in both languages")
			}
			arTitles, enTitles := strings.Split(item.titleAR, "|"), strings.Split(item.titleEN, "|")
			ch16ValidateLessonCopy(t, item.id+" Arabic", item.prefix, level.Steps, level.Hints, arTitles)
			ch16ValidateLessonCopy(t, item.id+" English", item.prefix, level.En.Steps, level.En.Hints, enTitles)
			if len(level.Steps) != len(level.En.Steps) {
				t.Fatal("Arabic and English step counts differ")
			}
			for i, step := range level.Steps {
				translated := level.En.Steps[i]
				if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
					t.Fatalf("step %d has different IDs or block counts", i)
				}
				for j, block := range step.Blocks {
					ch16AssertParity(t, block, translated.Blocks[j])
				}
			}
		})
	}
}

func ch16ValidateLessonCopy(t *testing.T, name, prefix string, steps []LessonStep, hints, titles []string) {
	t.Helper()
	sandboxCount := 0
	if len(steps) != 6 || len(titles) != 6 || len(hints) != 3 {
		t.Fatalf("%s: need six steps and three hints", name)
	}
	seen := map[string]bool{}
	claim := func(id string) {
		if id == "" || seen[id] || !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: empty, duplicate, or incorrectly prefixed ID %q", name, id)
		}
		seen[id] = true
	}
	words := 0
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i)
		}
		words += ch16TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch16TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid question options, feedback, or answer", name, block.ID)
					continue
				}
				words += ch16TextWords(strings.Join(block.Question.Options, " "))
				words += ch16TextWords(strings.Join(block.Question.Feedback, " "))
				for i := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[i]) == "" || !ch16OneSentence(block.Question.Feedback[i]) {
						t.Errorf("%s: %s has empty option or non-sentence feedback", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch16ValidateFrames(t, name, block)
			case "sandbox":
				sandboxCount++
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch16-") || block.Question != nil || block.Walkthrough != nil || block.Code != "" {
					t.Errorf("%s: %s has an invalid chapter 16 sandbox", name, block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	for i, hint := range hints {
		if !ch16OneSentence(hint) || ch16HintHasLabel(hint) {
			t.Errorf("%s: hint %d must be one unprefixed sentence: %q", name, i+1, hint)
		}
		words += ch16TextWords(hint)
	}
	if sandboxCount != 1 {
		t.Errorf("%s: expected exactly one sandbox, found %d", name, sandboxCount)
	}
	if len(steps[2].Blocks) != 3 || steps[2].Blocks[0].Kind != "question" || steps[2].Blocks[1].Kind != "walkthrough" || steps[2].Blocks[2].Kind != "sandbox" {
		t.Errorf("%s: Try it must predict, reveal, then show one sandbox", name)
	} else if strings.TrimSpace(ch16QuestionCode(steps[2].Blocks[0].Text)) != strings.TrimSpace(steps[2].Blocks[1].Code) {
		t.Errorf("%s: prediction code differs from its walkthrough", name)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200–320", name, words)
	}
	t.Logf("%s learner text: %d words", name, words)
}

func ch16ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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

func ch16AssertParity(t *testing.T, ar, en LessonBlock) {
	t.Helper()
	if ar.ID != en.ID || ar.Kind != en.Kind || ar.Code != en.Code || ar.ExampleID != en.ExampleID || !reflect.DeepEqual(ar.Imports, en.Imports) {
		t.Errorf("%s: Arabic and English block identity differs", ar.ID)
	}
	if ar.Sandbox != nil && (en.Sandbox == nil || ar.Sandbox.Kind != en.Sandbox.Kind || !reflect.DeepEqual(ar.Sandbox.Config, en.Sandbox.Config)) {
		t.Errorf("%s: sandbox config differs across languages", ar.ID)
	}
	if ar.Question != nil {
		if en.Question == nil || ar.Question.Answer != en.Question.Answer || len(ar.Question.Options) != len(en.Question.Options) {
			t.Errorf("%s: question answer or option count differs", ar.ID)
		}
	}
	if ar.Walkthrough != nil {
		if en.Walkthrough == nil || len(ar.Walkthrough.Frames) != len(en.Walkthrough.Frames) {
			t.Errorf("%s: walkthrough frame counts differ", ar.ID)
			return
		}
		for i, frame := range ar.Walkthrough.Frames {
			translated, translatedCopy := frame, en.Walkthrough.Frames[i]
			translated.Caption, translatedCopy.Caption = "", ""
			if !reflect.DeepEqual(translated, translatedCopy) {
				t.Errorf("%s frame %d: execution state differs", ar.ID, i)
			}
		}
	}
}

func TestCh16SandboxFixture(t *testing.T) {
	var fixture struct {
		ID     string         `json:"id"`
		Kind   string         `json:"kind"`
		Config map[string]any `json:"config"`
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch16.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &fixture); err != nil {
		t.Fatal(err)
	}
	level := Get("16-1-mutants")
	if level == nil || len(level.Steps) < 3 || len(level.Steps[2].Blocks) < 3 {
		t.Fatal("missing chapter 16 fixture level sandbox")
	}
	sandbox := level.Steps[2].Blocks[2]
	if sandbox.ID != fixture.ID || sandbox.Sandbox == nil || sandbox.Sandbox.Kind != fixture.Kind {
		t.Fatalf("chapter 16 fixture does not identify the lesson sandbox")
	}
	actualJSON, err := json.Marshal(sandbox.Sandbox.Config)
	if err != nil {
		t.Fatal(err)
	}
	var actual any
	if err := json.Unmarshal(actualJSON, &actual); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, fixture.Config) {
		t.Fatal("chapter 16 fixture config differs from the lesson config")
	}
}

func ch16OneSentence(sentence string) bool {
	plain := strings.TrimSpace(ch16InlineCodePattern.ReplaceAllString(sentence, ""))
	for _, end := range []string{".", "?", "!", "؟"} {
		plain = strings.TrimSuffix(plain, end)
	}
	return plain != "" && !strings.ContainsAny(plain, ".?!؟")
}

func ch16HintHasLabel(hint string) bool {
	plain := strings.ToLower(strings.TrimSpace(hint))
	for _, prefix := range []string{"concept:", "structure:", "location:", "الفكرة:", "الهيكل:", "الموضع:", "الموقع:"} {
		if strings.HasPrefix(plain, prefix) {
			return true
		}
	}
	return false
}

func ch16TextWords(text string) int {
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

func ch16QuestionCode(text string) string {
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

func TestCh16GlobalLessonIDs(t *testing.T) {
	seen := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch16BlockIDs(step.Blocks)...) {
				if prior, ok := seen[id]; ok {
					t.Errorf("lesson ID %q appears in %s and %s", id, prior, level.ID)
				}
				seen[id] = level.ID
			}
		}
	}
}

func ch16BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, len(blocks))
	for i, block := range blocks {
		ids[i] = block.ID
	}
	return ids
}

func TestCh16StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 16 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH16.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	examples := ch16ExamplePattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 10 {
		t.Fatalf("expected ten chapter 16 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch16VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch16VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected storyboard outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"16-1-mutants", "16-2-bug-hunt", "16-3-benchmark", "16-4-wild-inputs", "16-5-fake-clock"} {
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
		t.Fatalf("%s: missing verified output", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	match := ch16StdoutPattern.FindStringSubmatch(transcript)
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
	if got, want := strings.TrimSuffix(stdout.String(), "\n"), match[1]; got != want {
		t.Fatalf("%s stdout = %q, want %q", id, got, want)
	}
}
