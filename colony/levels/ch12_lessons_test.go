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

var ch12ExamplePattern = regexp.MustCompile("(?s)<!-- go-example id=(e12[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n\\x60\\x60\\x60go\\n(.*?)\\n\\x60\\x60\\x60")
var ch12StdoutPattern = regexp.MustCompile("(?m)^Real stdout: ([^\\r\\n]+)$")
var ch12InlineCodePattern = regexp.MustCompile("~[^~]*~")

func TestCh12Lessons(t *testing.T) {
	fixtures := ch12SandboxFixtures(t)
	cases := []struct {
		id       string
		prefix   string
		arabic   []LessonStep
		english  []LessonStep
		arHints  []string
		enHints  []string
		arTitles []string
		enTitles []string
	}{
		{
			id: "12-1-shared-total", prefix: "c121.", arabic: ch12Lesson121Arabic, english: ch12Lesson121English,
			arHints: Get("12-1-shared-total").Hints, enHints: Get("12-1-shared-total").En.Hints,
			arTitles: []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
			enTitles: []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
		},
		{
			id: "12-2-registry", prefix: "c122.", arabic: ch12Lesson122Arabic, english: ch12Lesson122English,
			arHints: Get("12-2-registry").Hints, enHints: Get("12-2-registry").En.Hints,
			arTitles: []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
			enTitles: []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
		},
		{
			id: "12-3-deadlock", prefix: "c123.", arabic: ch12Lesson123Arabic, english: ch12Lesson123English,
			arHints: Get("12-3-deadlock").Hints, enHints: Get("12-3-deadlock").En.Hints,
			arTitles: []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
			enTitles: []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
		},
	}

	for _, item := range cases {
		t.Run(item.id, func(t *testing.T) {
			level := Get(item.id)
			if level == nil || level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("missing bilingual revision 1 wiring")
			}
			if len(item.arabic) != 6 || len(item.english) != 6 || len(item.arHints) != 3 || len(item.enHints) != 3 {
				t.Fatal("expected six steps and three hints in each language")
			}
			if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
				t.Fatal("level does not expose six structured steps")
			}
			validateLessonSteps(t, ch12StepsWithoutSandbox(item.arabic))
			validateLessonSteps(t, ch12StepsWithoutSandbox(item.english))
			for _, hints := range [][]string{item.arHints, item.enHints} {
				for _, hint := range hints {
					if !ch12OneSentence(hint) {
						t.Errorf("hint should be one sentence: %q", hint)
					}
				}
			}
			if len(item.arHints) != len(level.Hints) || len(item.enHints) != len(level.En.Hints) {
				t.Fatal("hint wiring differs from tested lesson")
			}
			for i := range item.arabic {
				ar, en := item.arabic[i], item.english[i]
				if ar.ID != en.ID || ar.Title != item.arTitles[i] || en.Title != item.enTitles[i] || ar.ID != level.Steps[i].ID || en.ID != level.En.Steps[i].ID {
					t.Fatalf("step %d: ID or heading mismatch", i)
				}
				if len(ar.Blocks) != len(en.Blocks) || len(ar.Blocks) != len(level.Steps[i].Blocks) || len(en.Blocks) != len(level.En.Steps[i].Blocks) {
					t.Fatalf("step %s: block counts differ", ar.ID)
				}
				if !strings.HasPrefix(ar.ID, item.prefix) {
					t.Fatalf("step ID %q lacks %q prefix", ar.ID, item.prefix)
				}
				for j, block := range ar.Blocks {
					translated := en.Blocks[j]
					if !strings.HasPrefix(block.ID, item.prefix) || block.ID != translated.ID || block.ID != level.Steps[i].Blocks[j].ID || block.ID != level.En.Steps[i].Blocks[j].ID {
						t.Fatalf("block %q has mismatched or missing prefix", block.ID)
					}
					assertLessonParity(t, block, translated)
					if block.Question != nil {
						if block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) || len(block.Question.Options) != len(block.Question.Feedback) {
							t.Fatalf("%s: invalid answer or feedback count", block.ID)
						}
						for _, feedback := range block.Question.Feedback {
							if !ch12OneSentence(feedback) {
								t.Errorf("%s feedback should be one sentence: %q", block.ID, feedback)
							}
						}
					}
					if block.Walkthrough != nil {
						for _, frame := range block.Walkthrough.Frames {
							if len(strings.Fields(frame.Caption)) > 8 {
								t.Errorf("%s caption exceeds eight words: %q", block.ID, frame.Caption)
							}
						}
					}
				}
			}
			for _, steps := range [][]LessonStep{item.arabic, item.english} {
				try := steps[2].Blocks
				if len(try) != 3 || try[0].Kind != "question" || try[1].Kind != "walkthrough" || try[2].Kind != "sandbox" || try[2].Sandbox == nil {
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
				if sandboxCount != 1 || !strings.HasPrefix(try[2].Sandbox.Kind, "ch12-") {
					t.Fatalf("%s must have exactly one ch12 sandbox", item.id)
				}
				fixture, ok := fixtures[try[2].Sandbox.Kind]
				if !ok || !reflect.DeepEqual(ch12JSONValue(t, try[2].Sandbox.Config), fixture) {
					t.Fatalf("%s sandbox config differs from its JS fixture", try[2].Sandbox.Kind)
				}
			}
			for lang, content := range map[string]struct {
				steps []LessonStep
				hints []string
			}{"Arabic": {item.arabic, item.arHints}, "English": {item.english, item.enHints}} {
				words := ch12LessonWordCount(content.steps, content.hints)
				t.Logf("%s learner text: %d words", lang, words)
				if words < 200 || words > 320 {
					t.Errorf("%s learner text has %d words; expected 200–320", lang, words)
				}
			}
		})
	}
}

func ch12StepsWithoutSandbox(steps []LessonStep) []LessonStep {
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

func ch12SandboxFixtures(t *testing.T) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch12.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	fixtures := map[string]any{}
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	return fixtures
}

func ch12JSONValue(t *testing.T, value any) any {
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

func ch12OneSentence(sentence string) bool {
	plain := ch12InlineCodePattern.ReplaceAllString(sentence, "")
	plain = strings.TrimSpace(plain)
	for _, end := range []string{".", "?", "!", "؟"} {
		plain = strings.TrimSuffix(plain, end)
	}
	return plain != "" && !strings.ContainsAny(plain, ".?!؟")
}

func ch12LessonWordCount(steps []LessonStep, hints []string) int {
	words := len(strings.Fields(strings.Join(hints, " ")))
	for _, step := range steps {
		words += len(strings.Fields(step.Title))
		for _, block := range step.Blocks {
			words += len(strings.Fields(block.Text))
			words += len(strings.Fields(block.Note))
			if block.Question != nil {
				words += len(strings.Fields(strings.Join(block.Question.Options, " ")))
				words += len(strings.Fields(strings.Join(block.Question.Feedback, " ")))
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

func TestCh12StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 12 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH12.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	examples := ch12ExamplePattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 3 {
		t.Fatalf("expected one example per mission, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Fatalf("duplicate example ID %q", example[1])
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) {
			verifyCh12Example(t, goBin, storyboard, example)
		})
	}
}

func verifyCh12Example(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("unexpected example outcome %q", outcome)
	}
	levels := []string{"12-1-shared-total", "12-2-registry", "12-3-deadlock"}
	found := false
	for _, levelID := range levels {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						found = true
						if strings.ReplaceAll(block.Code, "\t", "    ") != strings.ReplaceAll(code, "\t", "    ") || strings.Join(block.Imports, ",") != imports {
							t.Fatalf("%s differs from its structured walkthrough", id)
						}
					}
				}
			}
		}
	}
	if !found {
		t.Fatalf("storyboard example %s is not used in a lesson", id)
	}

	verified := strings.SplitN(storyboard, "#### Verified "+id, 2)
	if len(verified) != 2 {
		t.Fatalf("missing verified output for %s", id)
	}
	transcript := strings.SplitN(verified[1], "#### Verified ", 2)[0]
	stdoutMatch := ch12StdoutPattern.FindStringSubmatch(transcript)
	if stdoutMatch == nil {
		t.Fatalf("missing real stdout for %s", id)
	}
	wantOutput := stdoutMatch[1] + "\n"
	for _, levelID := range levels {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						frames := block.Walkthrough.Frames
						if frames[len(frames)-1].Output+"\n" != wantOutput {
							t.Fatalf("%s walkthrough output differs from verified stdout", id)
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
	for name, contents := range map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + strings.Join(quotedImports, ";") + ")\nfunc main() {\n" + code + "\n}\n",
	} {
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
		t.Fatalf("build: %v\n%s", err, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	err = run.Run()
	if stdout.String() != wantOutput || err != nil || stderr.Len() != 0 {
		t.Fatalf("run: stdout=%q want=%q err=%v stderr=%q", stdout.String(), wantOutput, err, stderr.String())
	}
}
