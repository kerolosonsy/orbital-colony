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

func TestCh33Lessons(t *testing.T) {
	ids := []string{"33-1-retry-client", "33-2-reverse-proxy", "33-3-load-balancer", "33-4-gateway", "33-boss-survey-uplink"}
	globalIDs := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			ch33ClaimGlobalID(t, globalIDs, level.ID, step.ID)
			for _, block := range step.Blocks {
				ch33ClaimGlobalID(t, globalIDs, level.ID, block.ID)
			}
		}
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	prefixes := []string{"c331", "c332", "c333", "c334", "c33boss"}
	for i, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English lesson", id)
		}
		if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: expected a lab with revision 1 in both languages", id)
		}
		if i == 4 && !level.Boss {
			t.Errorf("%s: boss lesson is not wired to the boss level", id)
		}
		ch33ValidateCopy(t, id+" Arabic", prefixes[i], level.Steps, level.Hints, arabicTitles)
		ch33ValidateCopy(t, id+" English", prefixes[i], level.En.Steps, level.En.Hints, englishTitles)
		ch33CompareCopies(t, id, level.Steps, level.En.Steps)
	}
}

func ch33ClaimGlobalID(t *testing.T, seen map[string]string, levelID, id string) {
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

func ch33ValidateCopy(t *testing.T, name, prefix string, steps []LessonStep, hints, titles []string) {
	t.Helper()
	if len(steps) != 6 || len(titles) != 6 {
		t.Fatalf("%s: expected six lesson steps", name)
	}
	words := 0
	for i, step := range steps {
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, i)
		}
		if !strings.HasPrefix(step.ID, prefix+".") {
			t.Errorf("%s: step ID %q has wrong prefix", name, step.ID)
		}
		words += ch33TextWords(step.Title)
		for _, block := range step.Blocks {
			words += ch33TextWords(block.Text)
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
				words += ch33TextWords(strings.Join(block.Question.Options, " "))
				words += ch33TextWords(strings.Join(block.Question.Feedback, " "))
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
					words += ch33TextWords(frame.Caption)
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
			words += ch33TextWords(hint)
		}
	}
	if len(steps[2].Blocks) != 2 || steps[2].Blocks[0].Kind != "question" || steps[2].Blocks[1].Kind != "walkthrough" {
		t.Errorf("%s: Try it must predict, then reveal", name)
	} else if strings.TrimSpace(ch33QuestionCode(steps[2].Blocks[0].Text)) != strings.TrimSpace(steps[2].Blocks[1].Code) {
		t.Errorf("%s: prediction and walkthrough code differ", name)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200–320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch33CompareCopies(t *testing.T, name string, arabic, english []LessonStep) {
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
		}
	}
}

func ch33TextWords(text string) int {
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

func ch33QuestionCode(text string) string {
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

func TestCh33StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 33 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH33.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e33[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 10 {
		t.Fatalf("expected ten chapter 33 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch33VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch33VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unsupported outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range []string{"33-1-retry-client", "33-2-reverse-proxy", "33-3-load-balancer", "33-4-gateway", "33-boss-survey-uplink"} {
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
	for _, levelID := range []string{"33-1-retry-client", "33-2-reverse-proxy", "33-3-load-balancer", "33-4-gateway", "33-boss-survey-uplink"} {
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
