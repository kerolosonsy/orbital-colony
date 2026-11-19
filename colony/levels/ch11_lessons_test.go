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

func TestCh11Lessons(t *testing.T) {
	fixtures := ch11SandboxFixtures(t)
	cases := []struct {
		id, prefix string
		level      *Level
		arTitles   []string
		enTitles   []string
	}{
		{
			id: "11-1-scout", prefix: "c111.", level: Get("11-1-scout"),
			arTitles: []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
			enTitles: []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
		},
		{
			id: "11-2-worker-pool", prefix: "c112.", level: Get("11-2-worker-pool"),
			arTitles: []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
			enTitles: []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
		},
		{
			id: "11-3-storm", prefix: "c113.", level: Get("11-3-storm"),
			arTitles: []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
			enTitles: []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			level := tc.level
			if level == nil || level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("missing bilingual revision-1 lesson")
			}
			if len(level.Hints) != 3 || len(level.En.Hints) != 3 || len(level.Steps) != 6 || len(level.En.Steps) != 6 {
				t.Fatal("lesson must have six steps and three hints per language")
			}
			validateLessonSteps(t, ch11StepsWithoutSandbox(level.Steps))
			validateLessonSteps(t, ch11StepsWithoutSandbox(level.En.Steps))
			ch11CheckLessonIDs(t, tc.prefix, level.Steps)
			ch11CheckLessonIDs(t, tc.prefix, level.En.Steps)
			ch11CheckTitles(t, level.Steps, tc.arTitles)
			ch11CheckTitles(t, level.En.Steps, tc.enTitles)
			for i, step := range level.Steps {
				translated := level.En.Steps[i]
				if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
					t.Fatalf("step %s: translated structure differs", step.ID)
				}
				for j, block := range step.Blocks {
					other := translated.Blocks[j]
					assertLessonParity(t, block, other)
					if block.Question != nil && (len(block.Question.Feedback) != len(other.Question.Feedback) || len(block.Question.Options) != len(other.Question.Options) || block.Question.Answer != other.Question.Answer) {
						t.Fatalf("%s: translated question differs", block.ID)
					}
					if block.Kind != "markdown" && block.Kind != "walkthrough" && block.Kind != "question" {
						if block.Kind == "sandbox" && block.Sandbox != nil && strings.HasPrefix(block.Sandbox.Kind, "ch11-") {
							continue
						}
						t.Fatalf("%s: unsupported block kind %q", block.ID, block.Kind)
					}
					if block.Walkthrough != nil {
						for _, frame := range block.Walkthrough.Frames {
							if len(strings.Fields(frame.Caption)) > 8 {
								t.Fatalf("%s: frame caption is longer than eight words", block.ID)
							}
						}
					}
				}
			}
			for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
				try := steps[2].Blocks
				if len(try) != 3 || try[0].Kind != "question" || try[1].Kind != "walkthrough" || try[2].Kind != "sandbox" || try[0].Code == "" || try[0].Code != try[1].Code || try[0].ExampleID != try[1].ExampleID {
					t.Fatal("Try it must keep its question and walkthrough before one sandbox")
				}
				sandboxCount := 0
				for _, step := range steps {
					for _, block := range step.Blocks {
						if block.Kind == "sandbox" {
							sandboxCount++
						}
					}
				}
				if sandboxCount != 1 || try[2].Sandbox == nil || !strings.HasPrefix(try[2].Sandbox.Kind, "ch11-") {
					t.Fatalf("%s must have exactly one ch11 sandbox", tc.id)
				}
				fixture, ok := fixtures[try[2].Sandbox.Kind]
				if !ok || !reflect.DeepEqual(ch11JSONValue(t, try[2].Sandbox.Config), fixture) {
					t.Fatalf("%s sandbox config differs from its JS fixture", try[2].Sandbox.Kind)
				}
			}
			for lang, hints := range map[string][]string{"Arabic": level.Hints, "English": level.En.Hints} {
				labels := []string{"Concept:", "Structure:", "Specific:"}
				if lang == "Arabic" {
					labels = []string{"مفهوم:", "الترتيب:", "سطر محدد:"}
				}
				for i, hint := range hints {
					if strings.HasPrefix(hint, labels[i]) || strings.ContainsAny(hint, "\n\r") {
						t.Fatalf("%s: hint %d lacks its position label or spans lines", lang, i+1)
					}
				}
			}
			for lang, words := range map[string]int{
				"Arabic":  ch11LessonWordCount(level.Steps, level.Hints),
				"English": ch11LessonWordCount(level.En.Steps, level.En.Hints),
			} {
				if words < 200 || words > 320 {
					t.Fatalf("%s learner text has %d words; want 200–320", lang, words)
				}
				t.Logf("%s learner-authored words: %d", lang, words)
			}
		})
	}
}

func ch11StepsWithoutSandbox(steps []LessonStep) []LessonStep {
	filtered := make([]LessonStep, len(steps))
	for i, step := range steps {
		filtered[i].ID, filtered[i].Title = step.ID, step.Title
		for _, block := range step.Blocks {
			if block.Kind != "sandbox" {
				filtered[i].Blocks = append(filtered[i].Blocks, block)
			}
		}
	}
	return filtered
}

func ch11SandboxFixtures(t *testing.T) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch11.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	fixtures := map[string]any{}
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	return fixtures
}

func ch11JSONValue(t *testing.T, value any) any {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var normalized any
	if err := json.Unmarshal(body, &normalized); err != nil {
		t.Fatal(err)
	}
	return normalized
}

func ch11CheckLessonIDs(t *testing.T, prefix string, steps []LessonStep) {
	t.Helper()
	seen := map[string]bool{}
	for _, step := range steps {
		if !strings.HasPrefix(step.ID, prefix) || seen[step.ID] {
			t.Fatalf("invalid or duplicate step ID %q", step.ID)
		}
		seen[step.ID] = true
		for _, block := range step.Blocks {
			if !strings.HasPrefix(block.ID, prefix) || seen[block.ID] {
				t.Fatalf("invalid or duplicate block ID %q", block.ID)
			}
			seen[block.ID] = true
		}
	}
}

func ch11CheckTitles(t *testing.T, steps []LessonStep, want []string) {
	t.Helper()
	if len(steps) != len(want) {
		t.Fatalf("got %d step titles, want %d", len(steps), len(want))
	}
	for i, step := range steps {
		if step.Title != want[i] {
			t.Fatalf("step %d title = %q, want %q", i+1, step.Title, want[i])
		}
	}
}

func ch11LessonWordCount(steps []LessonStep, hints []string) int {
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

func TestCh11StoryboardExamplesMatchRealGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH11.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e11-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 6 {
		t.Fatalf("expected six chapter 11 examples, found %d", len(examples))
	}
	lessonExamples := map[string]string{}
	for _, steps := range [][]LessonStep{ch11ScoutLessonArabic, ch11PoolLessonArabic, ch11StormLessonArabic} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == "" {
					continue
				}
				if old, exists := lessonExamples[block.ExampleID]; exists && old != block.Code {
					t.Fatalf("%s: repeated example ID has different code", block.ExampleID)
				}
				lessonExamples[block.ExampleID] = block.Code
			}
		}
	}
	if len(lessonExamples) != len(examples) {
		t.Fatalf("found %d lesson examples and %d storyboard examples", len(lessonExamples), len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		id, imports, outcome, code := example[1], example[2], example[3], example[4]
		if seen[id] || outcome != "ok" {
			t.Fatalf("duplicate example ID or unsupported outcome for %s", id)
		}
		seen[id] = true
		if want, ok := lessonExamples[id]; !ok || want != code {
			t.Fatalf("%s: storyboard code differs from lesson", id)
		}
		t.Run(id, func(t *testing.T) {
			verifyCh11Example(t, goBin, storyboard, id, imports, code)
		})
	}
}

func verifyCh11Example(t *testing.T, goBin, storyboard, id, imports, code string) {
	t.Helper()
	parts := strings.SplitN(storyboard, "#### Verified `"+id+"`", 2)
	if len(parts) != 2 {
		t.Fatal("missing documented example output")
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatal("missing documented stdout")
	}
	wantOutput := match[1] + "\n"
	for _, steps := range [][]LessonStep{ch11ScoutLessonArabic, ch11PoolLessonArabic, ch11StormLessonArabic} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == id && block.Walkthrough != nil {
					frames := block.Walkthrough.Frames
					if frames[len(frames)-1].Output != strings.TrimSuffix(wantOutput, "\n") {
						t.Fatalf("%s walkthrough output differs from documented execution", id)
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
	diagnostics, buildErr := build.CombinedOutput()
	if buildErr != nil {
		t.Fatalf("build: %v\n%s", buildErr, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	runErr := run.Run()
	if stdout.String() != wantOutput {
		t.Fatalf("stdout = %q, documented %q", stdout.String(), wantOutput)
	}
	if runErr != nil || stderr.Len() != 0 {
		t.Fatalf("unexpected failure: %v\n%s", runErr, stderr)
	}
}
