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

var ch14ExamplePattern = regexp.MustCompile("(?s)<!-- go-example id=(e14-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
var ch14StdoutPattern = regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
var ch14InlineCodePattern = regexp.MustCompile("~[^~]*~")
var ch14SentenceBreakPattern = regexp.MustCompile(`[.!?؟]\s+\S`)

func TestCh14Lessons(t *testing.T) {
	sandboxFixtures := ch14ReadSandboxFixtures(t)
	cases := []struct {
		id     string
		prefix string
		kind   string
		pair   ch14LessonPair
	}{
		{"14-1-generic-max", "c141.", "ch14-generic-max", ch14Lesson141},
		{"14-2-stack", "c142.", "ch14-stack", ch14Lesson142},
		{"14-3-constraints", "c143.", "ch14-type-constraints", ch14Lesson143},
		{"14-4-airlock-queue", "c144.", "ch14-queue", ch14Lesson144},
		{"14-5-archive-tree", "c145.", "ch14-recursion", ch14Lesson145},
	}
	arTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	enTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	globalIDs := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch14BlockIDs(step.Blocks)...) {
				if previous, exists := globalIDs[id]; exists {
					t.Errorf("duplicate lesson ID %q in %s and %s", id, previous, level.ID)
				}
				globalIDs[id] = level.ID
			}
		}
	}
	for _, item := range cases {
		t.Run(item.id, func(t *testing.T) {
			level := Get(item.id)
			if level == nil || level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("missing bilingual revision 1 lesson wiring")
			}
			if len(level.Steps) != 6 || len(level.En.Steps) != 6 || len(level.Hints) != 3 || len(level.En.Hints) != 3 {
				t.Fatal("expected six steps and three hints in both languages")
			}
			if !reflect.DeepEqual(level.Steps, item.pair.Arabic) || !reflect.DeepEqual(level.En.Steps, item.pair.English) || !reflect.DeepEqual(level.Hints, item.pair.ArabicHints) || !reflect.DeepEqual(level.En.Hints, item.pair.EnglishHints) {
				t.Fatal("registered level does not reference the authored lesson data")
			}
			ch14ValidateLessonCopy(t, item.id+" Arabic", item.prefix, level.Steps, level.Hints, arTitles)
			ch14ValidateLessonCopy(t, item.id+" English", item.prefix, level.En.Steps, level.En.Hints, enTitles)
			for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
				var sandboxes []LessonBlock
				for _, step := range steps {
					for _, block := range step.Blocks {
						if block.Kind == "sandbox" {
							sandboxes = append(sandboxes, block)
						}
					}
				}
				if len(sandboxes) != 1 || len(steps[2].Blocks) != 3 || steps[2].Blocks[2].Kind != "sandbox" {
					t.Errorf("%s: expected exactly one sandbox after the Try it question and walkthrough", item.id)
					continue
				}
				sandbox := sandboxes[0]
				if sandbox.Sandbox == nil || !strings.HasPrefix(sandbox.Sandbox.Kind, "ch14-") {
					t.Errorf("%s: sandbox kind must use the ch14- prefix", item.id)
					continue
				}
				if sandbox.Sandbox.Kind != item.kind {
					t.Errorf("%s: got sandbox kind %q, want %q", item.id, sandbox.Sandbox.Kind, item.kind)
				}
				fixture, exists := sandboxFixtures[sandbox.Sandbox.Kind]
				if !exists || !reflect.DeepEqual(ch14JSONValue(sandbox.Sandbox.Config), fixture) {
					t.Errorf("%s: Config differs from fixture %q", item.id, sandbox.Sandbox.Kind)
				}
			}
			for i, step := range level.Steps {
				translated := level.En.Steps[i]
				if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
					t.Fatalf("%s: Arabic and English step structures differ", step.ID)
				}
				for j, block := range step.Blocks {
					other := translated.Blocks[j]
					if block.ID != other.ID || block.Kind != other.Kind {
						t.Errorf("%s: Arabic and English block identities differ", block.ID)
						continue
					}
					if block.Question != nil {
						if other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options) || len(block.Question.Feedback) != len(other.Question.Feedback) {
							t.Errorf("%s: question structure differs across languages", block.ID)
						}
					}
					if block.Walkthrough != nil {
						if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
							t.Errorf("%s: walkthrough structure differs across languages", block.ID)
							continue
						}
						for frameIndex, frame := range block.Walkthrough.Frames {
							translatedFrame := other.Walkthrough.Frames[frameIndex]
							frame.Caption, translatedFrame.Caption = "", ""
							if !reflect.DeepEqual(frame, translatedFrame) {
								t.Errorf("%s frame %d: execution state differs", block.ID, frameIndex)
							}
						}
					}
					if block.Sandbox != nil && (other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind || !reflect.DeepEqual(ch14ConfigShape(block.Sandbox.Config), ch14ConfigShape(other.Sandbox.Config))) {
						t.Errorf("%s: sandbox Config structure differs across languages", block.ID)
					}
				}
			}
		})
	}
}

func ch14BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, 0, len(blocks))
	for _, block := range blocks {
		ids = append(ids, block.ID)
	}
	return ids
}

func ch14ReadSandboxFixtures(t *testing.T) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch14.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]any
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	return fixtures
}

func ch14JSONValue(value any) any {
	encoded, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	var normalized any
	if err := json.Unmarshal(encoded, &normalized); err != nil {
		panic(err)
	}
	return normalized
}

func ch14ConfigShape(value any) any {
	switch value := value.(type) {
	case string:
		return "string"
	case map[string]any:
		result := make(map[string]any, len(value))
		for key, child := range value {
			result[key] = ch14ConfigShape(child)
		}
		return result
	case []any:
		result := make([]any, len(value))
		for index, child := range value {
			result[index] = ch14ConfigShape(child)
		}
		return result
	case []string:
		return make([]string, len(value))
	case []int:
		return append([]int(nil), value...)
	case [][]int:
		return append([][]int(nil), value...)
	default:
		return value
	}
}

func ch14ValidateLessonCopy(t *testing.T, name, prefix string, steps []LessonStep, hints, titles []string) {
	t.Helper()
	words := ch14LessonWordCount(steps, hints)
	if words < 200 || words > 320 {
		t.Errorf("%s learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner text: %d words", name, words)
	for i, step := range steps {
		if step.Title != titles[i] || len(step.Blocks) == 0 || !strings.HasPrefix(step.ID, prefix) {
			t.Errorf("%s: step %d has wrong title, prefix, or no blocks", name, i+1)
		}
		for _, block := range step.Blocks {
			if !strings.HasPrefix(block.ID, prefix) || strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has an empty identity or learner text", name, block.ID)
			}
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
				for i, feedback := range block.Question.Feedback {
					if strings.TrimSpace(block.Question.Options[i]) == "" || !ch14OneSentence(feedback) {
						t.Errorf("%s: %s option %d has empty text or multi-sentence feedback", name, block.ID, i)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				ch14ValidateFrames(t, name, block)
			case "sandbox":
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch14-") || len(block.Sandbox.Config) == 0 {
					t.Errorf("%s: %s has incomplete sandbox Config", name, block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
	}
	for i, hint := range hints {
		if !ch14OneSentence(hint) || ch14HasHintPrefix(hint) {
			t.Errorf("%s: hint %d must be one unprefixed sentence", name, i+1)
		}
	}
	if len(steps[2].Blocks) != 3 || steps[2].Blocks[0].Kind != "question" || steps[2].Blocks[1].Kind != "walkthrough" || steps[2].Blocks[2].Kind != "sandbox" {
		t.Errorf("%s: Try it must predict, reveal, then show one sandbox", name)
		return
	}
	if strings.TrimSpace(ch14QuestionCode(steps[2].Blocks[0].Text)) != strings.TrimSpace(steps[2].Blocks[1].Code) {
		t.Errorf("%s: Try it prediction differs from its walkthrough", name)
	}
}

func ch14ValidateFrames(t *testing.T, name string, block LessonBlock) {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Errorf("%s: %s has no initial frame", name, block.ID)
		return
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
		if frame.Panic != "" && i != len(frames)-1 {
			t.Errorf("%s: execution continues after a panic", block.ID)
		}
	}
}

func ch14QuestionCode(text string) string {
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

func ch14LessonWordCount(steps []LessonStep, hints []string) int {
	texts := append([]string(nil), hints...)
	for _, step := range steps {
		texts = append(texts, step.Title)
		for _, block := range step.Blocks {
			texts = append(texts, ch14VisibleText(block.Text))
			if block.Question != nil {
				texts = append(texts, block.Question.Options...)
				texts = append(texts, block.Question.Feedback...)
			}
			if block.Walkthrough != nil {
				for _, frame := range block.Walkthrough.Frames {
					texts = append(texts, frame.Caption)
				}
			}
		}
	}
	return len(strings.Fields(strings.Join(texts, " ")))
}

func ch14VisibleText(text string) string {
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
	return ch14InlineCodePattern.ReplaceAllString(strings.Join(visible, " "), "")
}

func ch14OneSentence(text string) bool {
	plain := strings.TrimSpace(ch14InlineCodePattern.ReplaceAllString(text, ""))
	return plain != "" && !ch14SentenceBreakPattern.MatchString(plain)
}

func ch14HasHintPrefix(hint string) bool {
	for _, prefix := range []string{"Concept:", "Structure:", "Location:", "الفكرة:", "الهيكل:", "الموضع:", "المفهوم:"} {
		if strings.HasPrefix(strings.TrimSpace(hint), prefix) {
			return true
		}
	}
	return false
}

func TestCh14StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 14 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "docs/storyboards/STORYBOARD_CH14.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	examples := ch14ExamplePattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 10 {
		t.Fatalf("expected ten chapter 14 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch14VerifyExample(t, goBin, storyboard, example) })
	}
	for _, item := range []ch14LessonPair{ch14Lesson141, ch14Lesson142, ch14Lesson143, ch14Lesson144, ch14Lesson145} {
		for _, steps := range [][]LessonStep{item.Arabic, item.English} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.Kind == "walkthrough" && !seen[block.ExampleID] {
						t.Errorf("%s references missing storyboard example %q", block.ID, block.ExampleID)
					}
				}
			}
		}
	}
}

func ch14VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unsupported outcome %q", id, outcome)
	}
	found := false
	var finalOutput string
	for _, item := range []ch14LessonPair{ch14Lesson141, ch14Lesson142, ch14Lesson143, ch14Lesson144, ch14Lesson145} {
		for _, steps := range [][]LessonStep{item.Arabic, item.English} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID != id {
						continue
					}
					found = true
					if block.Code != code || strings.Join(block.Imports, ",") != importNames {
						t.Fatalf("%s: lesson code or imports differ from the storyboard", id)
					}
					finalOutput = block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output
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
	match := ch14StdoutPattern.FindStringSubmatch(strings.Split(parts[1], "#### Verified")[0])
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	wantOutput := match[1]
	if finalOutput != wantOutput {
		t.Fatalf("%s: final frame output %q differs from transcript %q", id, finalOutput, wantOutput)
	}
	dir := t.TempDir()
	imports := make([]string, 0)
	for _, name := range strings.Split(importNames, ",") {
		imports = append(imports, fmt.Sprintf("%q", name))
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
	err = run.Run()
	if err != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, err, stderr.String())
	}
	if stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), wantOutput+"\n")
	}
}
