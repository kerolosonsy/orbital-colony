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

func TestCh05Lessons(t *testing.T) {
	ids := []string{"5-1-functions", "5-2-storm", "5-3-fault-log", "5-4-undock", "5-5-command-table"}
	fixtureBytes, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch05.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]any
	if err := json.Unmarshal(fixtureBytes, &fixtures); err != nil {
		t.Fatal(err)
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	foreignIDs := map[string]string{}
	for _, level := range All {
		if level.Chapter == 5 || level.Chapter == 13 {
			continue
		}
		translations := [][]LessonStep{level.Steps}
		if level.En != nil {
			translations = append(translations, level.En.Steps)
		}
		for _, steps := range translations {
			for _, step := range steps {
				foreignIDs[step.ID] = level.ID
				for _, block := range step.Blocks {
					foreignIDs[block.ID] = level.ID
				}
			}
		}
	}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English text", id)
		}
		if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: lesson revision must be 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		prefix := fmt.Sprintf("c5%d.", mission+1)
		ch05ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, prefix, 200, 320)
		ch05ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, prefix, 200, 320)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				if _, exists := foreignIDs[step.ID]; exists {
					t.Errorf("%s: ID %q collides with another chapter", id, step.ID)
				}
				for _, block := range step.Blocks {
					if _, exists := foreignIDs[block.ID]; exists {
						t.Errorf("%s: ID %q collides with another chapter", id, block.ID)
					}
				}
			}
		}
		ch05ValidateParity(t, id, level.Steps, level.En.Steps)
		ch05ValidateSandboxes(t, id, level.Steps, level.En.Steps, fixtures)
	}
}

func ch05ValidateCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string, minWords, maxWords int) {
	t.Helper()
	seen := map[string]bool{}
	claim := func(id string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q does not start with %q", name, id, prefix)
		}
		seen[id] = true
	}
	words := 0
	for i, step := range steps {
		claim(step.ID)
		if i >= len(titles) || step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has the wrong title or no blocks", name, i)
		}
		words += ch05TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch05TextWords(block.Text)
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
				words += ch05TextWords(strings.Join(block.Question.Options, " "))
				words += ch05TextWords(strings.Join(block.Question.Feedback, " "))
				for i := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[i]) == "" || strings.TrimSpace(block.Question.Feedback[i]) == "" {
						t.Errorf("%s: %s has empty option or feedback", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch05ValidateFrames(t, name, block)
			case "sandbox":
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch05-") || len(block.Sandbox.Config) == 0 {
					t.Errorf("%s: %s has an incomplete or unprefixed sandbox", name, block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected exactly three hints, got %d", name, len(hints))
	}
	for i, hint := range hints {
		for _, label := range []string{"Concept:", "Structure:", "Specific:", "المفهوم:", "الهيكل:", "الموضع:", "الفكرة:"} {
			if strings.HasPrefix(hint, label) {
				t.Errorf("%s: hint %d includes a positional label", name, i+1)
			}
		}
		if strings.TrimSpace(hint) == "" {
			t.Errorf("%s: hint %d is empty", name, i+1)
		}
		words += ch05TextWords(hint)
	}
	if len(steps) == 6 {
		try := steps[2]
		wantTryBlocks := 2
		if strings.HasPrefix(prefix, "c5") {
			wantTryBlocks = 3
		}
		if len(try.Blocks) != wantTryBlocks || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || (wantTryBlocks == 3 && try.Blocks[2].Kind != "sandbox") {
			if wantTryBlocks == 3 {
				t.Errorf("%s: Try it must keep its question and walkthrough before one sandbox", name)
			} else {
				t.Errorf("%s: Try it must predict first, then reveal with a walkthrough", name)
			}
		} else if strings.TrimSpace(ch05QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction code differs from its walkthrough", name)
		}
	}
	if words < minWords || words > maxWords {
		t.Errorf("%s: learner text has %d words; want %d-%d", name, words, minWords, maxWords)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch05ValidateSandboxes(t *testing.T, name string, arabic, english []LessonStep, fixtures map[string]any) {
	t.Helper()
	copies := [][]LessonStep{arabic, english}
	var first *LessonBlock
	for copyIndex, steps := range copies {
		count := 0
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.Kind != "sandbox" {
					continue
				}
				count++
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch05-") {
					t.Errorf("%s: sandbox kind must have a ch05- prefix", name)
					continue
				}
				fixture, ok := fixtures[block.Sandbox.Kind]
				if !ok {
					t.Errorf("%s: no JSON fixture for %q", name, block.Sandbox.Kind)
					continue
				}
				actualJSON, err := json.Marshal(block.Sandbox.Config)
				if err != nil {
					t.Errorf("%s: marshal %q config: %v", name, block.Sandbox.Kind, err)
					continue
				}
				fixtureJSON, err := json.Marshal(fixture)
				if err != nil || string(actualJSON) != string(fixtureJSON) {
					t.Errorf("%s: %q config differs from its JavaScript fixture", name, block.Sandbox.Kind)
				}
				ch05ValidateConfigText(t, name, block.Sandbox.Config)
				if copyIndex == 0 {
					first = &block
				} else if first != nil {
					if first.ID != block.ID || first.Sandbox.Kind != block.Sandbox.Kind || !reflect.DeepEqual(first.Sandbox.Config, block.Sandbox.Config) {
						t.Errorf("%s: Arabic and English sandbox IDs, kinds, or configs differ", name)
					}
				}
			}
		}
		if count != 1 {
			t.Errorf("%s: expected exactly one ch05- sandbox, found %d", name, count)
		}
	}
}

func ch05ValidateConfigText(t *testing.T, name string, value any) {
	t.Helper()
	switch current := value.(type) {
	case map[string]any:
		ar, hasArabic := current["ar"].(string)
		en, hasEnglish := current["en"].(string)
		if hasArabic && hasEnglish {
			if strings.TrimSpace(ar) == "" || strings.TrimSpace(en) == "" {
				t.Errorf("%s: sandbox localized text cannot be empty", name)
			}
			if len(strings.Fields(ar)) > 8 || len(strings.Fields(en)) > 8 {
				t.Errorf("%s: sandbox caption exceeds eight words", name)
			}
			return
		}
		for _, child := range current {
			ch05ValidateConfigText(t, name, child)
		}
	case []any:
		for _, child := range current {
			ch05ValidateConfigText(t, name, child)
		}
	}
}

func ch05ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
		words += ch05TextWords(frame.Caption)
	}
	return words
}

func ch05ValidateParity(t *testing.T, name string, arabic, english []LessonStep) {
	t.Helper()
	for i, step := range arabic {
		translated := english[i]
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
			if block.Sandbox != nil && (other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind || !reflect.DeepEqual(block.Sandbox.Config, other.Sandbox.Config)) {
				t.Errorf("%s: sandbox kind or config differs between languages", block.ID)
			}
		}
	}
}

func ch05TextWords(text string) int {
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

func ch05QuestionCode(text string) string {
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

func TestCh05StoryboardExamples(t *testing.T) {
	ch05RunStoryboardExamples(t, "docs/storyboards/STORYBOARD_CH05.md", `(?s)<!-- go-example id=(e5[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n`+"```go\n"+`(.*?)\n`+"```", 10, []string{"5-1-functions", "5-2-storm", "5-3-fault-log", "5-4-undock", "5-5-command-table"})
}

func ch05RunStoryboardExamples(t *testing.T, filename, pattern string, wantCount int, levelIDs []string) {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", filename))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	examples := regexp.MustCompile(pattern).FindAllStringSubmatch(storyboard, -1)
	if len(examples) != wantCount {
		t.Fatalf("expected %d storyboard examples, found %d", wantCount, len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch05VerifyExample(t, goBin, storyboard, example, levelIDs) })
	}
}

func ch05VerifyExample(t *testing.T, goBin, storyboard string, example []string, levelIDs []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						found = true
						if block.Code != code || strings.Join(block.Imports, ",") != importNames {
							t.Fatalf("%s: lesson code or imports differ from storyboard", id)
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
	outputMatch := regexp.MustCompile("(?s)Real stdout:\n\n```text\n(.*?)\n```").FindStringSubmatch(transcript)
	if outputMatch == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	wantOutput := outputMatch[1]
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
	if stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), wantOutput+"\n")
	}
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != wantOutput {
						t.Fatalf("%s: final walkthrough output differs from the transcript", id)
					}
				}
			}
		}
	}
}
