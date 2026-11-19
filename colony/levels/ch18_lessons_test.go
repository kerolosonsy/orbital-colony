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

func TestCh18Lessons(t *testing.T) {
	ids := []string{"18-1-mission-log", "18-2-black-box", "18-3-cargo-manifest", "18-4-signal-words", "18-5-flight-recorder", "18-6-crew-directory"}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	arabicIDs, englishIDs := map[string]bool{}, map[string]bool{}
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
		ch18ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, mission+1, arabicIDs)
		ch18ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, mission+1, englishIDs)
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			continue
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
				if block.Sandbox != nil || other.Sandbox != nil {
					if block.Sandbox == nil || other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind {
						t.Errorf("%s: sandbox kinds or configs differ between languages", block.ID)
					} else {
						sandboxConfigParity(t, block.Sandbox.Config, other.Sandbox.Config)
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

func TestCh18SandboxConfigsMatchFixtures(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch18.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]any
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, levelID := range []string{"18-1-mission-log", "18-2-black-box", "18-3-cargo-manifest", "18-4-signal-words", "18-5-flight-recorder", "18-6-crew-directory"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			var found *LessonSandbox
			for _, block := range steps[2].Blocks {
				if block.Kind == "sandbox" {
					if found != nil {
						t.Errorf("%s: more than one sandbox", levelID)
					}
					found = block.Sandbox
				}
			}
			if found == nil {
				t.Errorf("%s: Try it has no sandbox", levelID)
				continue
			}
			if !strings.HasPrefix(found.Kind, "ch18-") {
				t.Errorf("%s: sandbox kind %q lacks chapter prefix", levelID, found.Kind)
			}
			fixture, ok := fixtures[found.Kind]
			if !ok {
				t.Errorf("%s: no fixture for %s", levelID, found.Kind)
				continue
			}
			got, err := json.Marshal(found.Config)
			if err != nil {
				t.Fatal(err)
			}
			want, err := json.Marshal(fixture)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(want) {
				t.Errorf("%s: %s config differs from fixture", levelID, found.Kind)
			}
			seen[found.Kind] = true
		}
	}
	if len(seen) != len(fixtures) {
		t.Errorf("chapter 18 has %d fixture kinds but %d lesson kinds", len(fixtures), len(seen))
	}
}

func ch18ValidateCopy(t *testing.T, name string, steps []LessonStep, hints, titles []string, mission int, globalIDs map[string]bool) {
	t.Helper()
	if len(steps) != 6 {
		t.Errorf("%s: expected six steps, got %d", name, len(steps))
		return
	}
	prefix := fmt.Sprintf("c18%d.", mission)
	words := 0
	sandboxes := 0
	claim := func(id string) {
		if id == "" || globalIDs[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q has wrong mission prefix", name, id)
		}
		globalIDs[id] = true
	}
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i)
		}
		words += ch18TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" && block.Kind != "sandbox" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			if block.Kind != "sandbox" {
				words += ch18TextWords(block.Text)
			}
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" || len(block.Imports) != 0 || block.ExampleID != "" {
					t.Errorf("%s: %s markdown contains another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" || block.ExampleID != "" {
					t.Errorf("%s: %s has invalid question payload", name, block.ID)
					continue
				}
				q := block.Question
				if len(q.Options) < 2 || len(q.Options) != len(q.Feedback) || q.Answer < 0 || q.Answer >= len(q.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", name, block.ID)
					continue
				}
				words += ch18TextWords(strings.Join(q.Options, " ")) + ch18TextWords(strings.Join(q.Feedback, " "))
				for k := range q.Options {
					if strings.TrimSpace(q.Options[k]) == "" || !ch18OneSentence(q.Feedback[k]) {
						t.Errorf("%s: %s option %d needs text and one-sentence feedback", name, block.ID, k)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 || block.Question != nil || block.Sandbox != nil {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch18ValidateFrames(t, name, block)
			case "sandbox":
				sandboxes++
				if block.Sandbox == nil || block.Question != nil || block.Walkthrough != nil || block.Code != "" || !strings.HasPrefix(block.Sandbox.Kind, "ch18-") {
					t.Errorf("%s: %s has an invalid chapter 18 sandbox", name, block.ID)
				}
				if i != 2 {
					t.Errorf("%s: sandbox must be in Try it", name)
				}
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if sandboxes != 1 {
		t.Errorf("%s: expected exactly one sandbox, got %d", name, sandboxes)
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three unprefixed hints", name)
	} else {
		labels := []string{"Concept:", "Structure:", "Specific:", "Location:", "الفكرة:", "المفهوم:", "الهيكل:", "الموضع:"}
		for i, hint := range hints {
			for _, label := range labels {
				if strings.HasPrefix(strings.TrimSpace(hint), label) {
					t.Errorf("%s: hint %d includes a positional label", name, i+1)
				}
			}
			if !ch18OneSentence(hint) {
				t.Errorf("%s: hint %d must be one sentence", name, i+1)
			}
			words += ch18TextWords(hint)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 3 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || try.Blocks[2].Kind != "sandbox" {
			t.Errorf("%s: Try it must predict, reveal, then show one sandbox", name)
		} else if strings.TrimSpace(ch18QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction example differs from its reveal", name)
		}
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch18ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
			t.Errorf("%s: %s frame %d points outside code or has no caption", name, block.ID, i)
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
		words += len(strings.Fields(frame.Caption))
	}
	return words
}

func ch18QuestionCode(text string) string {
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

func ch18TextWords(text string) int {
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

func ch18OneSentence(text string) bool {
	var plain strings.Builder
	inInlineCode := false
	for _, r := range text {
		if r == '~' {
			inInlineCode = !inInlineCode
			continue
		}
		if !inInlineCode {
			plain.WriteRune(r)
		}
	}
	count := 0
	for _, r := range plain.String() {
		if strings.ContainsRune(".!?؟", r) {
			count++
		}
	}
	return count == 1
}

func TestCh18StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 18 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH18.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e18[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 12 {
		t.Fatalf("expected twelve chapter 18 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch18VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch18VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"18-1-mission-log", "18-2-black-box", "18-3-cargo-manifest", "18-4-signal-words", "18-5-flight-recorder", "18-6-crew-directory"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID != id {
						continue
					}
					found = true
					if block.Code != code || strings.Join(block.Imports, ",") != importNames {
						t.Fatalf("%s: lesson code or imports differ from storyboard", id)
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
	executable := filepath.Join(dir, "example")
	build := exec.CommandContext(ctx, goBin, "build", "-trimpath", "-gcflags=-lang=go1.24", "-o", executable, ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s: %v\n%s", id, err, diagnostics)
	}
	run := exec.CommandContext(ctx, executable)
	var stdout, stderr strings.Builder
	run.Stdout, run.Stderr = &stdout, &stderr
	if err := run.Run(); err != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, err, stderr.String())
	}
	if stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), wantOutput+"\n")
	}
	for _, levelID := range []string{"18-1-mission-log", "18-2-black-box", "18-3-cargo-manifest", "18-4-signal-words", "18-5-flight-recorder", "18-6-crew-directory"} {
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
