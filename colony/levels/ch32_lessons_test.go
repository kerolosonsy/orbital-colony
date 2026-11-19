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

func TestCh32Lessons(t *testing.T) {
	ids := []string{"32-1-websocket", "32-2-grpc-unary", "32-3-grpc-stream", "32-4-tls", "32-5-http2", "32-6-ssh"}
	registered := registeredSandboxKinds(t)
	globalIDs := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			ch32ClaimGlobalID(t, globalIDs, level.ID, step.ID)
			for _, block := range step.Blocks {
				ch32ClaimGlobalID(t, globalIDs, level.ID, block.ID)
			}
		}
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English lesson", id)
		}
		if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: expected a lab with revision 1 in both languages", id)
		}
		prefix := fmt.Sprintf("c32%d", mission+1)
		ch32ValidateCopy(t, id+" Arabic", prefix, level.Steps, level.Hints, arabicTitles)
		ch32ValidateCopy(t, id+" English", prefix, level.En.Steps, level.En.Hints, englishTitles)
		ch32CompareCopies(t, id, level.Steps, level.En.Steps)
		sandbox := level.Steps[2].Blocks[2].Sandbox
		if sandbox == nil || !strings.HasPrefix(sandbox.Kind, "ch32-") || !registered[sandbox.Kind] {
			t.Errorf("%s: missing registered chapter 32 sandbox", id)
		} else if englishSandbox := level.En.Steps[2].Blocks[2].Sandbox; englishSandbox == nil {
			t.Errorf("%s: missing English sandbox", id)
		} else {
			sandboxConfigParity(t, sandbox.Config, englishSandbox.Config)
		}
	}
}

func ch32ClaimGlobalID(t *testing.T, seen map[string]string, levelID, id string) {
	t.Helper()
	if id == "" {
		t.Errorf("%s: empty lesson ID", levelID)
		return
	}
	if previous, ok := seen[id]; ok {
		t.Errorf("lesson ID %q is shared by %s and %s", id, previous, levelID)
		return
	}
	seen[id] = levelID
}

func ch32ValidateCopy(t *testing.T, name, prefix string, steps []LessonStep, hints, titles []string) {
	t.Helper()
	if len(steps) != 6 || len(titles) != 6 {
		t.Fatalf("%s: expected six lesson steps", name)
	}
	words := 0
	sandboxes := 0
	for i, step := range steps {
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i)
		}
		if !strings.HasPrefix(step.ID, prefix+".") {
			t.Errorf("%s: step ID %q has wrong prefix", name, step.ID)
		}
		words += ch32TextWords(step.Title)
		for _, block := range step.Blocks {
			words += ch32TextWords(block.Text)
			if block.Kind == "sandbox" {
				sandboxes++
			}
			if !strings.HasPrefix(block.ID, prefix+".") || strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has wrong prefix or empty text", name, block.ID)
			}
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid question data", name, block.ID)
					continue
				}
				words += ch32TextWords(strings.Join(block.Question.Options, " "))
				words += ch32TextWords(strings.Join(block.Question.Feedback, " "))
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has incomplete walkthrough data", name, block.ID)
					continue
				}
				frames := block.Walkthrough.Frames
				if len(frames) < 2 || frames[0].Line != 0 {
					t.Errorf("%s: %s needs an initial frame", name, block.ID)
				}
				maxLine := strings.Count(block.Code, "\n") + 1
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
					words += ch32TextWords(frame.Caption)
				}
			case "sandbox":
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch32-") || len(block.Sandbox.Config) == 0 {
					t.Errorf("%s: %s has invalid sandbox data", name, block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected exactly three hints", name)
	} else {
		for i, hint := range hints {
			for _, label := range []string{"Concept:", "Structure:", "Location:", "الفكرة:", "الهيكل:", "الموضع:", "المفهوم:"} {
				if strings.HasPrefix(strings.TrimSpace(hint), label) {
					t.Errorf("%s: hint %d starts with forbidden label %q", name, i+1, label)
				}
			}
			if strings.TrimSpace(hint) == "" {
				t.Errorf("%s: hint %d is empty", name, i+1)
			}
			words += ch32TextWords(hint)
		}
	}
	if sandboxes != 1 {
		t.Errorf("%s: expected exactly one sandbox, found %d", name, sandboxes)
	}
	if len(steps[2].Blocks) != 3 || steps[2].Blocks[0].Kind != "question" || steps[2].Blocks[1].Kind != "walkthrough" || steps[2].Blocks[2].Kind != "sandbox" {
		t.Errorf("%s: Try it must predict, reveal, then sandbox", name)
	} else if strings.TrimSpace(ch32QuestionCode(steps[2].Blocks[0].Text)) != strings.TrimSpace(steps[2].Blocks[1].Code) {
		t.Errorf("%s: prediction and walkthrough code differ", name)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200–320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch32CompareCopies(t *testing.T, name string, arabic, english []LessonStep) {
	t.Helper()
	if len(arabic) != len(english) {
		t.Errorf("%s: Arabic and English step counts differ", name)
		return
	}
	for i, step := range arabic {
		translated := english[i]
		if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
			t.Errorf("%s: step %d identity or block count differs", name, i)
			continue
		}
		for j, block := range step.Blocks {
			other := translated.Blocks[j]
			if block.ID != other.ID || block.Kind != other.Kind {
				t.Errorf("%s: %s Arabic/English block structure differs", name, block.ID)
			}
			if block.Question != nil {
				if other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options) || len(block.Question.Feedback) != len(other.Question.Feedback) {
					t.Errorf("%s: %s question parity differs", name, block.ID)
				}
			}
			if block.Walkthrough != nil {
				if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
					t.Errorf("%s: %s walkthrough parity differs", name, block.ID)
					continue
				}
				for k, frame := range block.Walkthrough.Frames {
					translatedFrame := other.Walkthrough.Frames[k]
					frame.Caption, translatedFrame.Caption = "", ""
					if !reflect.DeepEqual(frame, translatedFrame) {
						t.Errorf("%s: %s frame %d execution state differs", name, block.ID, k)
					}
				}
			}
			if block.Sandbox != nil {
				if other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind {
					t.Errorf("%s: sandbox kind differs across languages", block.ID)
				} else {
					sandboxConfigParity(t, block.Sandbox.Config, other.Sandbox.Config)
				}
			}
		}
	}
}

func TestCh32SandboxFixtures(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch32.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]struct {
		Kind   string          `json:"kind"`
		Config json.RawMessage `json:"config"`
	}
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	ids := []string{"32-1-websocket", "32-2-grpc-unary", "32-3-grpc-stream", "32-4-tls", "32-5-http2", "32-6-ssh"}
	if len(fixtures) != len(ids) {
		t.Fatalf("fixture has %d levels; want %d", len(fixtures), len(ids))
	}
	for _, id := range ids {
		fixture, ok := fixtures[id]
		if !ok {
			t.Errorf("missing fixture for %s", id)
			continue
		}
		block := Get(id).Steps[2].Blocks[2]
		if fixture.Kind != block.Sandbox.Kind {
			t.Errorf("%s: fixture kind %q differs from lesson kind %q", id, fixture.Kind, block.Sandbox.Kind)
		}
		actualJSON, err := json.Marshal(block.Sandbox.Config)
		if err != nil {
			t.Fatal(err)
		}
		var actual, expected any
		if err := json.Unmarshal(actualJSON, &actual); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(fixture.Config, &expected); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("%s: JavaScript fixture differs from Go sandbox config", id)
		}
	}
}

func ch32TextWords(text string) int {
	lines := strings.Split(text, "\n")
	var plain strings.Builder
	inFence, inInline := false, false
	for _, line := range lines {
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

func ch32QuestionCode(text string) string {
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

func TestCh32StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 32 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH32.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e32[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 12 {
		t.Fatalf("expected twelve chapter 32 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch32VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch32VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unsupported outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"32-1-websocket", "32-2-grpc-unary", "32-3-grpc-stream", "32-4-tls", "32-5-http2", "32-6-ssh"} {
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
	match := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```").FindStringSubmatch(strings.Split(parts[1], "#### Verified")[0])
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	want := match[1]
	dir := t.TempDir()
	imports := make([]string, 0)
	for _, name := range strings.Split(importNames, ",") {
		imports = append(imports, fmt.Sprintf("%q", name))
	}
	files := map[string]string{"go.mod": "module storyboard/" + id + "\n\ngo 1.24.0\n", "main.go": "package main\nimport (" + strings.Join(imports, ";") + ")\n" + code + "\n"}
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
	if stdout.String() != want+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), want+"\n")
	}
	for _, levelID := range []string{"32-1-websocket", "32-2-grpc-unary", "32-3-grpc-stream", "32-4-tls", "32-5-http2", "32-6-ssh"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != want {
						t.Fatalf("%s: final walkthrough output differs from transcript", id)
					}
				}
			}
		}
	}
}
