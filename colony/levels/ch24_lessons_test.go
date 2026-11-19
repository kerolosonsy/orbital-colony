package levels

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestCh24Lessons(t *testing.T) {
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	cases := []struct {
		id, prefix string
	}{
		{"24-1-ping-pong", "c241."},
		{"24-2-producers", "c242."},
		{"24-3-work-crew", "c243."},
		{"24-4-pipeline", "c244."},
		{"24-5-beacon-check", "c245."},
		{"24-6-chunk-download", "c246."},
		{"24-7-map-tiles", "c247."},
	}
	for _, tc := range cases {
		t.Run(tc.id, func(t *testing.T) {
			level := Get(tc.id)
			if level == nil || level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("missing bilingual revision-1 lesson")
			}
			if len(level.Hints) != 3 || len(level.En.Hints) != 3 {
				t.Fatal("each language needs exactly three hints")
			}
			ch24ValidateCopy(t, tc.id+" Arabic", tc.prefix, level.Steps, level.Hints, arabicTitles)
			ch24ValidateCopy(t, tc.id+" English", tc.prefix, level.En.Steps, level.En.Hints, englishTitles)
			if len(level.Steps) != len(level.En.Steps) {
				t.Fatal("Arabic and English lesson lengths differ")
			}
			for i, step := range level.Steps {
				translated := level.En.Steps[i]
				if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
					t.Fatalf("step %s: translated structure differs", step.ID)
				}
				for j, block := range step.Blocks {
					other := translated.Blocks[j]
					if block.ID != other.ID || block.Kind != other.Kind {
						t.Errorf("%s: translated block identity differs", block.ID)
					}
					if block.Question != nil && (other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options) || len(block.Question.Feedback) != len(other.Question.Feedback)) {
						t.Errorf("%s: translated question shape differs", block.ID)
					}
					if block.Walkthrough != nil {
						if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
							t.Errorf("%s: translated walkthrough shape differs", block.ID)
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
				}
			}
			try := level.Steps[2].Blocks
			if len(try) != 2 || try[0].Kind != "question" || try[1].Kind != "walkthrough" || ch24QuestionCode(try[0].Text) != try[1].Code {
				t.Fatal("Try it must predict from code, then reveal the same code")
			}
		})
	}
}

func ch24ValidateCopy(t *testing.T, name, prefix string, steps []LessonStep, hints []string, titles []string) {
	t.Helper()
	if len(steps) != 6 {
		t.Fatalf("%s: expected six steps, got %d", name, len(steps))
	}
	seen := map[string]bool{}
	words := 0
	claim := func(id string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or repeated lesson ID %q", name, id)
		}
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q does not begin with %q", name, id, prefix)
		}
		seen[id] = true
	}
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has an unexpected title or no blocks", name, i)
		}
		words += ch24TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no text", name, block.ID)
			}
			words += ch24TextWords(block.Text)
			for _, paragraph := range strings.Split(block.Text, "\n\n") {
				for _, line := range strings.Split(paragraph, "\n") {
					if ch24TextWords(line) > 20 {
						t.Errorf("%s: %s has a paragraph longer than twenty words", name, block.ID)
					}
				}
			}
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid question data", name, block.ID)
					continue
				}
				words += ch24TextWords(strings.Join(block.Question.Options, " "))
				words += ch24TextWords(strings.Join(block.Question.Feedback, " "))
				for i := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[i]) == "" || strings.TrimSpace(block.Question.Feedback[i]) == "" || strings.ContainsAny(block.Question.Feedback[i], "\n\r") {
						t.Errorf("%s: %s has empty or multiline option feedback", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Question != nil || block.Sandbox != nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch24ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected exactly three hints", name)
	} else {
		labels := []string{"Concept:", "Structure:", "Specific:", "المفهوم:", "الهيكل:", "الموضع:", "الفكرة:", "الترتيب:", "سطر محدد:", "مفهوم:"}
		for i, hint := range hints {
			trimmed := strings.TrimSpace(hint)
			for _, label := range labels {
				if strings.HasPrefix(trimmed, label) {
					t.Errorf("%s: hint %d starts with a position label", name, i+1)
				}
			}
			if strings.ContainsAny(hint, "\n\r") || !strings.HasSuffix(trimmed, ".") {
				t.Errorf("%s: hint %d must be one sentence", name, i+1)
			}
			words += ch24TextWords(hint)
		}
	}
	if words < 200 || words > 320 {
		t.Errorf("%s learner text has %d words; want 200–320", name, words)
	}
	t.Logf("%s learner-authored words: %d", name, words)
}

func ch24ValidateFrames(t *testing.T, name string, block LessonBlock) int {
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
		words += ch24TextWords(frame.Caption)
	}
	return words
}

func ch24TextWords(text string) int {
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

func ch24QuestionCode(text string) string {
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

func TestCh24LessonIDsAreGloballyUnique(t *testing.T) {
	seen := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch24BlockIDs(step.Blocks)...) {
				if previous, ok := seen[id]; ok {
					t.Errorf("lesson ID %q is shared by %s and %s", id, previous, level.ID)
				}
				seen[id] = level.ID
			}
		}
	}
}

func ch24BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, len(blocks))
	for i, block := range blocks {
		ids[i] = block.ID
	}
	return ids
}

func TestCh24StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 24 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH24.md"))
	if err != nil {
		t.Fatal(err)
	}
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e24[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n```go\n(.*?)\n```")
	examples := pattern.FindAllStringSubmatch(string(body), -1)
	if len(examples) != 14 {
		t.Fatalf("expected fourteen chapter 24 examples, found %d", len(examples))
	}
	walkthroughs := map[string]LessonBlock{}
	for _, id := range []string{"24-1-ping-pong", "24-2-producers", "24-3-work-crew", "24-4-pipeline", "24-5-beacon-check", "24-6-chunk-download", "24-7-map-tiles"} {
		level := Get(id)
		for _, step := range level.Steps {
			for _, block := range step.Blocks {
				if block.Walkthrough != nil {
					walkthroughs[block.ExampleID] = block
				}
			}
		}
	}
	seen := map[string]bool{}
	for _, example := range examples {
		id, imports, code := example[1], example[2], example[4]
		block, ok := walkthroughs[id]
		if !ok || seen[id] || block.Code != code || strings.Join(block.Imports, ",") != imports {
			t.Fatalf("storyboard example %s is missing, duplicated, or differs from its walkthrough", id)
		}
		seen[id] = true
		t.Run(id, func(t *testing.T) {
			dir := t.TempDir()
			var source strings.Builder
			source.WriteString("package main\nimport (\n")
			for _, importPath := range strings.Split(imports, ",") {
				fmt.Fprintf(&source, "%q\n", importPath)
			}
			source.WriteString(")\n")
			source.WriteString(code)
			path := filepath.Join(dir, "main.go")
			if err := os.WriteFile(path, []byte(source.String()), 0o600); err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(goBin, "run", "main.go")
			cmd.Dir = dir
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("go run failed: %v\n%s", err, output)
			}
			got := strings.TrimSuffix(string(output), "\n")
			want := block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output
			if got != want {
				t.Fatalf("output mismatch: got %q, want %q", got, want)
			}
		})
	}
}
