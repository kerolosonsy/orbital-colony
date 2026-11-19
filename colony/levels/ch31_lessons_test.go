package levels

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCh31Lessons(t *testing.T) {
	ids := []string{
		"31-1-echo", "31-2-chat", "31-3-discovery",
		"31-4-frames", "31-5-dns", "31-6-port-check",
	}
	prefixes := []string{"c311.", "c312.", "c313.", "c314.", "c315.", "c316."}
	arTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	enTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	ch30CheckGlobalIDs(t)
	registered := registeredSandboxKinds(t)
	for mission, id := range ids {
		t.Run(id, func(t *testing.T) {
			level := Get(id)
			if level == nil || level.En == nil {
				t.Fatal("missing level or English lesson")
			}
			if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("expected a lab with lesson revision 1 in both languages")
			}
			ch31ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, prefixes[mission], arTitles)
			ch31ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, prefixes[mission], enTitles)
			if len(level.Steps) != len(level.En.Steps) {
				t.Fatalf("%s: Arabic and English step counts differ", id)
			}
			for stepIndex, arabicStep := range level.Steps {
				englishStep := level.En.Steps[stepIndex]
				if arabicStep.ID != englishStep.ID || len(arabicStep.Blocks) != len(englishStep.Blocks) {
					t.Fatalf("%s: step %d identity or block count differs", id, stepIndex)
				}
				for blockIndex, arabicBlock := range arabicStep.Blocks {
					englishBlock := englishStep.Blocks[blockIndex]
					if arabicBlock.ID != englishBlock.ID || arabicBlock.Kind != englishBlock.Kind {
						t.Errorf("%s: Arabic and English block structure differs", arabicBlock.ID)
					}
					if arabicBlock.Question != nil && (englishBlock.Question == nil || arabicBlock.Question.Answer != englishBlock.Question.Answer || len(arabicBlock.Question.Options) != len(englishBlock.Question.Options)) {
						t.Errorf("%s: question structure differs", arabicBlock.ID)
					}
					if arabicBlock.Walkthrough != nil && (englishBlock.Walkthrough == nil || arabicBlock.Code != englishBlock.Code || arabicBlock.ExampleID != englishBlock.ExampleID || !reflect.DeepEqual(arabicBlock.Imports, englishBlock.Imports)) {
						t.Errorf("%s: walkthrough structure differs", arabicBlock.ID)
					}
					if arabicBlock.Sandbox != nil {
						if englishBlock.Sandbox == nil || arabicBlock.Sandbox.Kind != englishBlock.Sandbox.Kind || !reflect.DeepEqual(arabicBlock.Sandbox.Config, englishBlock.Sandbox.Config) {
							t.Errorf("%s: sandbox config differs across languages", arabicBlock.ID)
						}
						if !registered[arabicBlock.Sandbox.Kind] {
							t.Errorf("%s: kind %q has no chapter renderer", arabicBlock.ID, arabicBlock.Sandbox.Kind)
						}
						sandboxConfigParity(t, arabicBlock.Sandbox.Config, englishBlock.Sandbox.Config)
					}
				}
			}
		})
	}
}

func ch31ValidateCopy(t *testing.T, name string, steps []LessonStep, hints []string, prefix string, titles []string) {
	t.Helper()
	if len(steps) != 6 {
		t.Fatalf("%s: expected six lesson steps", name)
	}
	words := 0
	sandboxes := 0
	for index, step := range steps {
		if step.Title != titles[index] || len(step.Blocks) == 0 || !strings.HasPrefix(step.ID, prefix) {
			t.Errorf("%s: invalid step %d", name, index+1)
		}
		words += ch30VisibleWords(step.Title)
		for _, block := range step.Blocks {
			if !strings.HasPrefix(block.ID, prefix) || strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has an invalid ID or empty text", name, block.ID)
			}
			words += ch30VisibleWords(block.Text)
			if block.Question != nil {
				for _, option := range block.Question.Options {
					words += ch30VisibleWords(option)
				}
				for _, feedback := range block.Question.Feedback {
					words += ch30VisibleWords(feedback)
				}
			}
			if block.Walkthrough != nil {
				for _, frame := range block.Walkthrough.Frames {
					words += ch30VisibleWords(frame.Caption)
				}
			}
			if block.Sandbox != nil {
				sandboxes++
				if block.Kind != "sandbox" || !strings.HasPrefix(block.Sandbox.Kind, "ch31-") || len(block.Sandbox.Config) == 0 {
					t.Errorf("%s: %s has an invalid sandbox", name, block.ID)
				}
			}
		}
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints", name)
	}
	for _, hint := range hints {
		words += ch30VisibleWords(hint)
	}
	if sandboxes != 1 {
		t.Errorf("%s: expected exactly one sandbox, found %d", name, sandboxes)
	}
	if len(steps[2].Blocks) != 3 || steps[2].Blocks[0].Kind != "question" || steps[2].Blocks[1].Kind != "walkthrough" || steps[2].Blocks[2].Kind != "sandbox" {
		t.Errorf("%s: Try it must keep prediction and walkthrough before one sandbox", name)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: lesson has %d visible words; want 200–320", name, words)
	}
}

func TestCh31SandboxFixtures(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch31.fixtures.json"))
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
	ids := []string{"31-1-echo", "31-2-chat", "31-3-discovery", "31-4-frames", "31-5-dns", "31-6-port-check"}
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

func TestCh31StoryboardExamples(t *testing.T) {
	ch30RunStoryboardExamples(t, "docs/storyboards/STORYBOARD_CH31.md", "e31", []string{
		"31-1-echo", "31-2-chat", "31-3-discovery",
		"31-4-frames", "31-5-dns", "31-6-port-check",
	}, 12)
}
