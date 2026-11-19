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

func TestCh02Lessons(t *testing.T) {
	ids := []string{"2-1-inventory", "2-2-survey", "2-3-callsigns", "2-4-units"}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seen := map[string]bool{}
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
		prefix := fmt.Sprintf("c2%d.", mission+1)
		ch02ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, prefix, seen)
		ch02ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, prefix, nil)
		ch02ValidateParity(t, level)
		ch02ValidateSandboxes(t, level)
	}
}

func TestCh02SandboxFixtures(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch02.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Kind   string         `json:"kind"`
		Config map[string]any `json:"config"`
	}
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	ids := []string{"2-1-inventory", "2-2-survey", "2-3-callsigns", "2-4-units"}
	if len(fixtures) != len(ids) {
		t.Fatalf("expected %d chapter 2 sandbox fixtures, found %d", len(ids), len(fixtures))
	}
	for _, id := range ids {
		fixture, ok := fixtures[id]
		if !ok {
			t.Errorf("%s: missing sandbox fixture", id)
			continue
		}
		var found *LessonSandbox
		for _, block := range Get(id).Steps[2].Blocks {
			if block.Sandbox != nil {
				found = block.Sandbox
			}
		}
		if found == nil || found.Kind != fixture.Kind {
			t.Errorf("%s: sandbox kind does not match its fixture", id)
			continue
		}
		encoded, err := json.Marshal(found.Config)
		if err != nil {
			t.Fatal(err)
		}
		var actual map[string]any
		if err := json.Unmarshal(encoded, &actual); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, fixture.Config) {
			t.Errorf("%s: sandbox config differs from its JavaScript fixture", id)
		}
	}
}

func TestCh02LessonIDsGloballyUnique(t *testing.T) {
	chapterIDs := []string{
		"1-1-first-signal", "1-2-first-harvest", "1-3-planting",
		"2-1-inventory", "2-2-survey", "2-3-callsigns", "2-4-units",
	}
	owned := map[string]string{}
	for _, levelID := range chapterIDs {
		level := Get(levelID)
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch02BlockIDs(step.Blocks)...) {
				if previous, exists := owned[id]; exists {
					t.Errorf("lesson ID %q is shared by %s and %s", id, previous, levelID)
				}
				owned[id] = levelID
			}
		}
	}
	for _, level := range All {
		if level.Chapter == 1 || level.Chapter == 2 {
			continue
		}
		for _, step := range level.Steps {
			ids := append([]string{step.ID}, ch02BlockIDs(step.Blocks)...)
			for _, id := range ids {
				if previous, exists := owned[id]; exists {
					t.Errorf("lesson ID %q is shared by %s and %s", id, previous, level.ID)
				}
			}
		}
	}
}

func ch02BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, 0, len(blocks))
	for _, block := range blocks {
		ids = append(ids, block.ID)
	}
	return ids
}

func ch02ValidateCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, prefix string, seen map[string]bool) {
	t.Helper()
	words := 0
	claim := func(id string) {
		if id == "" || (!strings.HasPrefix(id, prefix) && !strings.HasPrefix(id, "ch02-")) {
			t.Errorf("%s: ID %q is empty or lacks prefix %q", name, id, prefix)
		}
		if seen != nil && seen[id] {
			t.Errorf("%s: duplicate ID %q", name, id)
		}
		if seen != nil {
			seen[id] = true
		}
	}
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has the wrong title or no blocks", name, i)
		}
		words += ch02TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			maxRunes := 140
			if strings.Contains(name, "Arabic") {
				maxRunes = 110
			}
			ch02ValidateTextLines(t, name, block, maxRunes)
			words += ch02TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", name, block.ID)
					continue
				}
				words += ch02TextWords(strings.Join(block.Question.Options, " "))
				words += ch02TextWords(strings.Join(block.Question.Feedback, " "))
				for k, feedback := range block.Question.Feedback {
					if !strings.HasSuffix(strings.TrimSpace(feedback), ".") {
						t.Errorf("%s: feedback %d must be one sentence", block.ID, k+1)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch02ValidateFrames(t, name, block)
			case "sandbox":
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch02-") || len(block.Sandbox.Config) == 0 || block.Question != nil || block.Walkthrough != nil || block.Code != "" {
					t.Errorf("%s: %s has an invalid chapter 2 sandbox", name, block.ID)
				}
			default:
				t.Errorf("%s: %s has disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
	} else {
		for i, hint := range hints {
			trimmed := strings.TrimSpace(hint)
			for _, label := range []string{"Concept:", "Structure:", "Location:", "المفهوم:", "الهيكل:", "الموضع:"} {
				if strings.HasPrefix(trimmed, label) {
					t.Errorf("%s: hint %d starts with a position label", name, i+1)
				}
			}
			if !strings.HasSuffix(trimmed, ".") {
				t.Errorf("%s: hint %d must be one sentence ending in a period", name, i+1)
			}
			words += ch02TextWords(hint)
		}
	}
	if len(steps) == 6 {
		try := steps[2]
		if len(try.Blocks) != 3 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || try.Blocks[2].Kind != "sandbox" {
			t.Errorf("%s: Try it must predict, reveal, then show one sandbox", name)
		} else if strings.TrimSpace(ch02QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction code differs from its walkthrough", name)
		}
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner copy has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch02ValidateSandboxes(t *testing.T, level *Level) {
	t.Helper()
	for _, copy := range []struct {
		name  string
		steps []LessonStep
	}{{"Arabic", level.Steps}, {"English", level.En.Steps}} {
		count := 0
		for stepIndex, step := range copy.steps {
			for _, block := range step.Blocks {
				if block.Kind != "sandbox" {
					continue
				}
				count++
				if stepIndex != 2 || block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch02-") {
					t.Errorf("%s %s: sandbox must be in Try it with a ch02 kind", level.ID, copy.name)
				}
			}
		}
		if count != 1 {
			t.Errorf("%s %s: expected exactly one sandbox, found %d", level.ID, copy.name, count)
		}
	}
}

func ch02ValidateTextLines(t *testing.T, name string, block LessonBlock, maxRunes int) {
	t.Helper()
	inFence := false
	for i, line := range strings.Split(block.Text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence && len([]rune(strings.TrimSpace(line))) > maxRunes {
			t.Errorf("%s: %s text line %d may wrap beyond three panel lines", name, block.ID, i+1)
		}
	}
}

func ch02ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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

func ch02QuestionCode(text string) string {
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

func ch02TextWords(text string) int {
	var plain strings.Builder
	inFence, inInline := false, false
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
				inInline = !inInline
				continue
			}
			if !inInline {
				plain.WriteRune(r)
			}
		}
		plain.WriteByte(' ')
	}
	return len(strings.Fields(plain.String()))
}

func ch02ValidateParity(t *testing.T, level *Level) {
	t.Helper()
	for i, step := range level.Steps {
		other := level.En.Steps[i]
		if len(step.Blocks) != len(other.Blocks) {
			t.Errorf("%s: Arabic and English block counts differ", step.ID)
			continue
		}
		if step.ID != other.ID {
			t.Errorf("%s: Arabic and English step IDs differ", step.ID)
		}
		for j, block := range step.Blocks {
			translated := other.Blocks[j]
			if block.ID != translated.ID || block.Kind != translated.Kind {
				t.Errorf("%s: Arabic and English block identities differ", block.ID)
			}
			if block.Sandbox != nil {
				if translated.Sandbox == nil || block.Sandbox.Kind != translated.Sandbox.Kind {
					t.Errorf("%s: Arabic and English sandbox kinds differ", block.ID)
				} else {
					sandboxConfigParity(t, block.Sandbox.Config, translated.Sandbox.Config)
				}
			}
			if block.Question != nil {
				if translated.Question == nil || block.Question.Answer != translated.Question.Answer || len(block.Question.Options) != len(translated.Question.Options) || len(block.Question.Feedback) != len(translated.Question.Feedback) {
					t.Errorf("%s: Arabic and English question structures differ", block.ID)
				}
			}
			if block.Walkthrough != nil {
				if translated.Walkthrough == nil || block.Code != translated.Code || block.ExampleID != translated.ExampleID || !reflect.DeepEqual(block.Imports, translated.Imports) || len(block.Walkthrough.Frames) != len(translated.Walkthrough.Frames) {
					t.Errorf("%s: Arabic and English walkthrough structures differ", block.ID)
					continue
				}
				for k, frame := range block.Walkthrough.Frames {
					translatedFrame := translated.Walkthrough.Frames[k]
					frame.Caption, translatedFrame.Caption = "", ""
					if !reflect.DeepEqual(frame, translatedFrame) {
						t.Errorf("%s frame %d: execution states differ", block.ID, k)
					}
				}
			}
		}
	}
}

func TestCh02StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("Go toolchain unavailable")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH02.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e2[1-4]-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 2 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch02VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch02VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	var found *LessonBlock
	foundCount := 0
	for _, levelID := range []string{"2-1-inventory", "2-2-survey", "2-3-callsigns", "2-4-units"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						if found != nil && (found.Code != block.Code || found.ExampleID != block.ExampleID) {
							t.Fatalf("%s: translated lesson examples differ", id)
						}
						copy := block
						found = &copy
						foundCount++
					}
				}
			}
		}
	}
	if found == nil || foundCount != 2 || found.Code != code || strings.Join(found.Imports, ",") != importNames {
		t.Fatalf("%s: storyboard code must match both translated walkthroughs", id)
	}
	parts := strings.SplitN(storyboard, "#### Verified "+id+" — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing verified output", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	want := match[1]
	frames := found.Walkthrough.Frames
	if len(frames) == 0 || frames[len(frames)-1].Output != want {
		t.Fatalf("%s: final walkthrough frame output differs from storyboard stdout", id)
	}
	dir := t.TempDir()
	imports := make([]string, 0)
	for _, importName := range strings.Split(importNames, ",") {
		imports = append(imports, fmt.Sprintf("%q", importName))
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (\n" + strings.Join(imports, "\n") + "\n)\n" + code + "\n",
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, goBin, "run", ".")
	cmd.Dir = dir
	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: go run failed: %v\n%s", id, err, got)
	}
	if strings.TrimSuffix(string(got), "\n") != want {
		t.Fatalf("%s: stdout mismatch: got %q, want %q", id, got, want)
	}
}
