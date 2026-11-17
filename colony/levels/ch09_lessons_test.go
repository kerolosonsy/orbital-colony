package levels

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

type ch09LessonCopy struct {
	language string
	steps    []LessonStep
	titles   []string
	prefix   string
	hints    []string
	labels   []string
}

func TestCh09Lessons(t *testing.T) {
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	missions := []struct {
		id, prefix string
		arabic     []LessonStep
		english    []LessonStep
	}{
		{"9-1-stringer", "c9m1.", ch09StringerArabic, ch09StringerEnglish},
		{"9-2-autopilot", "c9m2.", ch09AutopilotArabic, ch09AutopilotEnglish},
		{"9-3-signals", "c9m3.", ch09SignalsArabic, ch09SignalsEnglish},
		{"9-4-small-plugs", "c9m4.", ch09PlugsArabic, ch09PlugsEnglish},
	}
	for _, mission := range missions {
		level := Get(mission.id)
		t.Run(mission.id, func(t *testing.T) {
			if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("Arabic and English lesson revisions must both be 1")
			}
			copies := []ch09LessonCopy{
				{"Arabic", mission.arabic, arabicTitles, mission.prefix, level.Hints, []string{"مفهوم", "ترتيب", "محدد"}},
				{"English", mission.english, englishTitles, mission.prefix, level.En.Hints, []string{"Concept", "Structure", "Specific"}},
			}
			for _, copy := range copies {
				t.Run(copy.language, func(t *testing.T) { ch09ValidateLessonCopy(t, copy) })
			}
			for i, step := range mission.arabic {
				if i >= len(mission.english) || len(step.Blocks) != len(mission.english[i].Blocks) {
					t.Fatalf("step %d block count differs across languages", i+1)
				}
				for j, block := range step.Blocks {
					assertLessonParity(t, block, mission.english[i].Blocks[j])
				}
			}
		})
	}
}

func ch09ValidateLessonCopy(t *testing.T, copy ch09LessonCopy) {
	t.Helper()
	if len(copy.steps) != 6 || len(copy.titles) != 6 || len(copy.hints) != 3 || len(copy.labels) != 3 {
		t.Fatalf("got %d steps and %d hints; want six and three", len(copy.steps), len(copy.hints))
	}
	validateLessonSteps(t, copy.steps)
	tryBlocks := copy.steps[2].Blocks
	if len(tryBlocks) != 2 || tryBlocks[0].Kind != "question" || tryBlocks[1].Kind != "walkthrough" {
		t.Errorf("Try it must ask for a prediction before its walkthrough")
	}
	for i, hint := range copy.hints {
		if strings.HasPrefix(hint, copy.labels[i]+": ") || strings.ContainsAny(strings.TrimSpace(hint), "\n!?؟") {
			t.Errorf("hint %d needs one sentence with the %q label", i+1, copy.labels[i])
		}
	}
	words := ch09LessonWordCount(copy.steps, copy.hints)
	if words < 200 || words > 320 {
		t.Errorf("learner word count = %d, want 200–320", words)
	}
	t.Logf("learner word count: %d", words)
	seen := map[string]bool{}
	for i, step := range copy.steps {
		if step.Title != copy.titles[i] || !strings.HasPrefix(step.ID, copy.prefix) || seen[step.ID] {
			t.Errorf("step %d title or ID is invalid: %q", i+1, step.ID)
		}
		seen[step.ID] = true
		for _, block := range step.Blocks {
			if !strings.HasPrefix(block.ID, copy.prefix) || seen[block.ID] {
				t.Errorf("empty-prefix or duplicate ID %q", block.ID)
			}
			seen[block.ID] = true
			if block.Kind == "sandbox" {
				t.Errorf("%s: chapter 9 lessons do not add sandbox blocks", block.ID)
			}
			if block.Kind == "walkthrough" {
				for frameIndex, frame := range block.Walkthrough.Frames {
					if len(strings.Fields(frame.Caption)) > 8 {
						t.Errorf("%s frame %d caption exceeds eight words", block.ID, frameIndex)
					}
				}
			}
			for _, line := range strings.Split(block.Text, "\n") {
				if len([]rune(line)) > 140 {
					t.Errorf("%s: a text line may exceed three panel lines", block.ID)
				}
			}
			if block.Question != nil {
				for _, texts := range [][]string{block.Question.Options, block.Question.Feedback} {
					for _, text := range texts {
						if len([]rune(text)) > 140 {
							t.Errorf("%s: question text may exceed three panel lines", block.ID)
						}
					}
				}
			}
		}
	}
}

func ch09LessonWordCount(steps []LessonStep, hints []string) int {
	words := len(strings.Fields(strings.Join(hints, " ")))
	for _, step := range steps {
		words += len(strings.Fields(step.Title))
		for _, block := range step.Blocks {
			words += len(strings.Fields(block.Text))
			if block.Question != nil {
				words += len(strings.Fields(strings.Join(block.Question.Options, " ") + " " + strings.Join(block.Question.Feedback, " ")))
			}
			if block.Walkthrough != nil {
				for _, frame := range block.Walkthrough.Frames {
					words += len(strings.Fields(frame.Caption))
				}
			}
		}
	}
	return words
}

func TestCh09StoryboardExamplesMatchRealGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH09.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e9m[1-4]-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter-9 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Fatalf("duplicate example ID %q", example[1])
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) {
			ch09VerifyExample(t, goBin, storyboard, example)
		})
	}
}

func ch09VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("unexpected outcome %q", outcome)
	}
	found := 0
	for _, levelID := range []string{"9-1-stringer", "9-2-autopilot", "9-3-signals", "9-4-small-plugs"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID != id {
						continue
					}
					found++
					if block.Code != code || strings.Join(block.Imports, ",") != imports {
						t.Fatalf("%s differs from its structured lesson example", id)
					}
				}
			}
		}
	}
	if found != 2 {
		t.Fatalf("%s occurs %d times in structured lessons, want bilingual pair", id, found)
	}
	parts := strings.SplitN(storyboard, "#### Verified `"+id+"`", 2)
	if len(parts) != 2 {
		t.Fatalf("missing documented output for %s", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("missing real stdout for %s", id)
	}
	wantOutput := match[1] + "\n"
	for _, levelID := range []string{"9-1-stringer", "9-2-autopilot", "9-3-signals", "9-4-small-plugs"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough != nil {
						frames := block.Walkthrough.Frames
						if frames[len(frames)-1].Output != strings.TrimSuffix(wantOutput, "\n") {
							t.Fatalf("%s walkthrough output differs from documented output", id)
						}
					}
				}
			}
		}
	}
	dir := t.TempDir()
	quotedImports := []string{}
	for _, name := range strings.Split(imports, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", name))
	}
	main := "package main\nimport (" + strings.Join(quotedImports, ";") + ")\nfunc main() {\n" + code + "\n}\n"
	for name, contents := range map[string]string{"go.mod": "module storyboard/" + id + "\n\ngo 1.24.0\n", "main.go": main} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, goBin, "build", "-trimpath", "-gcflags=-lang=go1.24", "-o", "example", ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, buildErr := build.CombinedOutput()
	if buildErr != nil {
		t.Fatalf("build: %v\n%s", buildErr, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	runErr := run.Run()
	if stdout.String() != wantOutput || runErr != nil || stderr.Len() != 0 {
		t.Fatalf("run: err=%v stdout=%q want=%q stderr=%q", runErr, stdout.String(), wantOutput, stderr.String())
	}
}
