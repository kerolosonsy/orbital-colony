package levels

import (
	ch10context "context"
	ch10json "encoding/json"
	ch10fmt "fmt"
	ch10os "os"
	ch10exec "os/exec"
	ch10filepath "path/filepath"
	ch10reflect "reflect"
	ch10regexp "regexp"
	ch10sort "sort"
	ch10strings "strings"
	ch10testing "testing"
	ch10time "time"
)

func TestCh10SandboxFixturesMatchLessons(t *ch10testing.T) {
	body, err := ch10os.ReadFile(ch10filepath.Join("..", "..", "web", "sandbox", "ch10.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]map[string]any
	if err := ch10json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, id := range []string{"10-1-second-drone", "10-2-swarm", "10-3-divide"} {
		level := Get(id)
		for _, block := range level.Steps[2].Blocks {
			if block.Sandbox == nil {
				continue
			}
			fixture, ok := fixtures[block.Sandbox.Kind]
			if !ok || !ch10reflect.DeepEqual(block.Sandbox.Config, fixture) {
				t.Errorf("%s config differs from its JavaScript fixture", block.Sandbox.Kind)
			}
			seen[block.Sandbox.Kind] = true
		}
	}
	if len(seen) != len(fixtures) {
		t.Errorf("found %d chapter 10 sandbox kinds for %d fixtures", len(seen), len(fixtures))
	}
}

func ch10ValidateStepsWithoutSandbox(t *ch10testing.T, steps []LessonStep) {
	t.Helper()
	filtered := make([]LessonStep, len(steps))
	for i, step := range steps {
		filtered[i] = step
		filtered[i].Blocks = nil
		for _, block := range step.Blocks {
			if block.Sandbox == nil {
				filtered[i].Blocks = append(filtered[i].Blocks, block)
			}
		}
	}
	validateLessonSteps(t, filtered)
}

func TestCh10Lessons(t *ch10testing.T) {
	wantTitlesArabic := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	wantTitlesEnglish := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	prefixes := map[string]string{
		"10-1-second-drone": "c101.",
		"10-2-swarm":        "c102.",
		"10-3-divide":       "c103.",
	}
	seenArabic, seenEnglish := map[string]bool{}, map[string]bool{}
	count := 0
	for _, level := range All {
		prefix, inChapter := prefixes[level.ID]
		if !inChapter {
			continue
		}
		count++
		if level.LessonRevision != 1 || level.En == nil || level.En.LessonRevision != 1 {
			t.Fatalf("%s: lesson revision must be 1 in both languages", level.ID)
		}
		if len(level.Hints) != 3 || len(level.En.Hints) != 3 {
			t.Fatalf("%s: expected three hints in both languages", level.ID)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Fatalf("%s: expected six steps in both languages", level.ID)
		}
		ch10ValidateStepsWithoutSandbox(t, level.Steps)
		ch10ValidateStepsWithoutSandbox(t, level.En.Steps)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			if len(steps[2].Blocks) != 3 || steps[2].Blocks[0].Kind != "question" || steps[2].Blocks[1].Kind != "walkthrough" || steps[2].Blocks[2].Kind != "sandbox" {
				t.Errorf("%s: Try it must keep its prediction and walkthrough before one sandbox", level.ID)
			}
			sandboxes := 0
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.Kind != "sandbox" {
						continue
					}
					sandboxes++
					if block.Sandbox == nil || !ch10strings.HasPrefix(block.Sandbox.Kind, "ch10-") {
						t.Errorf("%s: sandbox must use a ch10- kind", block.ID)
					}
				}
			}
			if sandboxes != 1 {
				t.Errorf("%s: got %d sandboxes, want exactly one", level.ID, sandboxes)
			}
		}
		for i, step := range level.Steps {
			translated := level.En.Steps[i]
			if step.Title != wantTitlesArabic[i] || translated.Title != wantTitlesEnglish[i] {
				t.Errorf("%s: step %d has incorrect bilingual title", level.ID, i+1)
			}
			if step.ID != translated.ID || !ch10strings.HasPrefix(step.ID, prefix) || len(step.Blocks) != len(translated.Blocks) {
				t.Errorf("%s: step %d has mismatched or unprefixed identity", level.ID, i+1)
			}
			if seenArabic[step.ID] || seenEnglish[translated.ID] {
				t.Errorf("%s: duplicate chapter step ID", level.ID)
			}
			seenArabic[step.ID], seenEnglish[translated.ID] = true, true
			for j, block := range step.Blocks {
				other := translated.Blocks[j]
				if !ch10strings.HasPrefix(block.ID, prefix) || seenArabic[block.ID] || seenEnglish[other.ID] {
					t.Errorf("%s: duplicate or unprefixed block ID %q", level.ID, block.ID)
				}
				seenArabic[block.ID], seenEnglish[other.ID] = true, true
				assertLessonParity(t, block, other)
				if block.Question != nil && (block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options)) {
					t.Errorf("%s: answer index is out of range", block.ID)
				}
				if block.Walkthrough != nil {
					for _, frame := range block.Walkthrough.Frames {
						if frame.Line < 0 || frame.Line > ch10strings.Count(block.Code, "\n")+1 {
							t.Errorf("%s: frame line %d is outside the code", block.ID, frame.Line)
						}
						if len(ch10strings.Fields(frame.Caption)) > 8 {
							t.Errorf("%s: frame caption exceeds eight words", block.ID)
						}
					}
				}
			}
		}
		for _, copy := range []struct {
			steps []LessonStep
			hints []string
		}{{level.Steps, level.Hints}, {level.En.Steps, level.En.Hints}} {
			words := ch10LearnerWordCount(copy.steps, copy.hints)
			if words < 200 || words > 320 {
				breakdown := []string{}
				for _, step := range copy.steps {
					breakdown = append(breakdown, step.Title+"="+ch10fmt.Sprint(ch10LearnerWordCount([]LessonStep{step}, nil)))
				}
				t.Errorf("%s: learner text has %d words; want 200 (%s; hints=%d)", level.ID, words, ch10strings.Join(breakdown, ", "), ch10LearnerWordCount(nil, copy.hints))
			}
		}
	}
	if count != len(prefixes) {
		t.Fatalf("found %d chapter 10 missions, want %d", count, len(prefixes))
	}
}

func ch10LearnerWordCount(steps []LessonStep, hints []string) int {
	text := append([]string(nil), hints...)
	for _, step := range steps {
		text = append(text, step.Title)
		for _, block := range step.Blocks {
			text = append(text, block.Text)
			if block.Question != nil {
				text = append(text, block.Question.Options...)
				text = append(text, block.Question.Feedback...)
			}
			if block.Walkthrough != nil {
				for _, frame := range block.Walkthrough.Frames {
					text = append(text, frame.Caption)
				}
			}
		}
	}
	return len(ch10strings.Fields(ch10strings.Join(text, " ")))
}

func TestCh10StoryboardExamples(t *ch10testing.T) {
	goBin, err := ch10exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := ch10os.ReadFile(ch10filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH10.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := ch10regexp.MustCompile("(?s)<!-- go-example id=(e10-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 6 {
		t.Fatalf("expected six chapter 10 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Fatalf("duplicate storyboard example ID %q", example[1])
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *ch10testing.T) {
			verifyCh10Example(t, goBin, storyboard, example)
		})
	}
	for _, level := range All {
		if _, ok := map[string]bool{"10-1-second-drone": true, "10-2-swarm": true, "10-3-divide": true}[level.ID]; !ok {
			continue
		}
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.Code != "" && !seen[block.ExampleID] {
						t.Errorf("%s references missing storyboard example %q", block.ID, block.ExampleID)
					}
				}
			}
		}
	}
}

func verifyCh10Example(t *ch10testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" && outcome != "unordered" {
		t.Fatalf("%s: unsupported outcome %q", id, outcome)
	}
	found := false
	for _, level := range All {
		if level.Chapter != 10 {
			continue
		}
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID != id {
						continue
					}
					found = true
					if block.Code != code || ch10strings.Join(block.Imports, ",") != imports {
						t.Fatalf("%s: lesson example differs from storyboard", id)
					}
				}
			}
		}
	}
	if !found {
		t.Fatalf("%s: storyboard example is not used by a lesson", id)
	}

	section := "#### Verified `" + id + "` — " + outcome
	parts := ch10strings.SplitN(storyboard, section, 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing documented output", id)
	}
	transcript := ch10strings.SplitN(parts[1], "#### Verified", 2)[0]
	outputPattern := ch10regexp.MustCompile("(?s)Real stdout(?: \\(one possible order\\))?:\\n\\n```text\\n(.*?)\\n```")
	match := outputPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("%s: missing stdout transcript", id)
	}
	wantOutput := match[1] + "\n"
	for _, level := range All {
		if level.Chapter != 10 {
			continue
		}
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough != nil {
						frames := block.Walkthrough.Frames
						if frames[len(frames)-1].Output != ch10strings.TrimSuffix(wantOutput, "\n") {
							t.Fatalf("%s: walkthrough output differs from storyboard", id)
						}
					}
				}
			}
		}
	}

	dir := t.TempDir()
	quotedImports := []string{}
	for _, name := range ch10strings.Split(imports, ",") {
		quotedImports = append(quotedImports, ch10fmt.Sprintf("%q", name))
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + ch10strings.Join(quotedImports, ";") + ")\nfunc main() {\n" + code + "\n}\n",
	}
	for name, contents := range files {
		if err := ch10os.WriteFile(ch10filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := ch10context.WithTimeout(ch10context.Background(), 90*ch10time.Second)
	defer cancel()
	build := ch10exec.CommandContext(ctx, goBin, "build", "-trimpath", "-gcflags=-lang=go1.24", "-o", "example", ".")
	build.Dir = dir
	build.Env = append(ch10os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, buildErr := build.CombinedOutput()
	if buildErr != nil {
		t.Fatalf("build %s: %v\n%s", id, buildErr, diagnostics)
	}
	run := ch10exec.CommandContext(ctx, ch10filepath.Join(dir, "example"))
	stdout, stderr := &ch10strings.Builder{}, &ch10strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	runErr := run.Run()
	gotOutput := stdout.String()
	if outcome == "unordered" {
		if ch10normalizeLines(gotOutput) != ch10normalizeLines(wantOutput) {
			t.Fatalf("%s: stdout lines = %q, documented lines = %q", id, gotOutput, wantOutput)
		}
	} else if gotOutput != wantOutput {
		t.Fatalf("%s: stdout = %q, documented %q", id, gotOutput, wantOutput)
	}
	if runErr != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, runErr, stderr)
	}
}

func ch10normalizeLines(output string) string {
	lines := ch10strings.Split(ch10strings.TrimSuffix(output, "\n"), "\n")
	ch10sort.Strings(lines)
	return ch10strings.Join(lines, "\n")
}
