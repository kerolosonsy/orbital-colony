package levels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestCh13Lessons(t *testing.T) {
	sandboxFixtures := ch13ReadSandboxFixtures(t)
	ids := []string{"13-1-recall", "13-2-timeouts", "13-3-first-wins", "13-4-chain-of-command"}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	foreignIDs := map[string]string{}
	for _, level := range All {
		if level.Chapter == 13 {
			continue
		}
		translations := [][]LessonStep{level.Steps}
		if level.En != nil {
			translations = append(translations, level.En.Steps)
		}
		for _, steps := range translations {
			for _, step := range steps {
				foreignIDs[step.ID] = level.ID
				for _, block := range step.Blocks {
					foreignIDs[block.ID] = level.ID
				}
			}
		}
	}
	for mission, id := range ids {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English text", id)
		}
		if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: lesson revision must be 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		prefix := fmt.Sprintf("c13%d.", mission+1)
		ch05ValidateCopy(t, id+" Arabic", ch13WithoutSandboxes(level.Steps), level.Hints, arabicTitles, prefix, 200, 320)
		ch05ValidateCopy(t, id+" English", ch13WithoutSandboxes(level.En.Steps), level.En.Hints, englishTitles, prefix, 200, 320)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				if _, exists := foreignIDs[step.ID]; exists {
					t.Errorf("%s: ID %q collides with another chapter", id, step.ID)
				}
				for _, block := range step.Blocks {
					if _, exists := foreignIDs[block.ID]; exists {
						t.Errorf("%s: ID %q collides with another chapter", id, block.ID)
					}
				}
			}
		}
		ch13ValidateSandboxPair(t, id, level.Steps, level.En.Steps, sandboxFixtures)
		ch05ValidateParity(t, id, level.Steps, level.En.Steps)
	}
}

func ch13WithoutSandboxes(steps []LessonStep) []LessonStep {
	copy := make([]LessonStep, len(steps))
	for index, step := range steps {
		copy[index] = step
		copy[index].Blocks = make([]LessonBlock, 0, len(step.Blocks))
		for _, block := range step.Blocks {
			if block.Kind != "sandbox" {
				copy[index].Blocks = append(copy[index].Blocks, block)
			}
		}
	}
	return copy
}

func ch13ValidateSandboxPair(t *testing.T, levelID string, arabic, english []LessonStep, fixtures map[string]any) {
	t.Helper()
	wantKinds := map[string]string{
		"13-1-recall": "ch13-recall", "13-2-timeouts": "ch13-timeout",
		"13-3-first-wins": "ch13-first-wins", "13-4-chain-of-command": "ch13-context-chain",
	}
	var arabicBlocks, englishBlocks []LessonBlock
	for _, step := range arabic {
		for _, block := range step.Blocks {
			if block.Kind == "sandbox" {
				arabicBlocks = append(arabicBlocks, block)
			}
		}
	}
	for _, step := range english {
		for _, block := range step.Blocks {
			if block.Kind == "sandbox" {
				englishBlocks = append(englishBlocks, block)
			}
		}
	}
	if len(arabicBlocks) != 1 || len(englishBlocks) != 1 || len(arabic) < 3 || len(english) < 3 || len(arabic[2].Blocks) != 3 || len(english[2].Blocks) != 3 || arabic[2].Blocks[2].Kind != "sandbox" || english[2].Blocks[2].Kind != "sandbox" {
		t.Errorf("%s: expected exactly one sandbox after the Try it question and walkthrough", levelID)
		return
	}
	ar, en := arabicBlocks[0], englishBlocks[0]
	if ar.ID != en.ID || ar.Sandbox == nil || en.Sandbox == nil || ar.Sandbox.Kind != en.Sandbox.Kind {
		t.Errorf("%s: Arabic and English sandbox identities differ", levelID)
		return
	}
	if !strings.HasPrefix(ar.Sandbox.Kind, "ch13-") || !strings.HasPrefix(en.Sandbox.Kind, "ch13-") {
		t.Errorf("%s: sandbox kind must use the ch13- prefix", levelID)
	}
	if ar.Sandbox.Kind != wantKinds[levelID] {
		t.Errorf("%s: got sandbox kind %q, want %q", levelID, ar.Sandbox.Kind, wantKinds[levelID])
	}
	if !reflect.DeepEqual(ch13ConfigShape(ar.Sandbox.Config), ch13ConfigShape(en.Sandbox.Config)) {
		t.Errorf("%s: sandbox Config structure differs across languages", levelID)
	}
	fixture, exists := fixtures[ar.Sandbox.Kind]
	if !exists || !reflect.DeepEqual(ch13JSONValue(ar.Sandbox.Config), fixture) || !reflect.DeepEqual(ch13JSONValue(en.Sandbox.Config), fixture) {
		t.Errorf("%s: Config differs from fixture %q", levelID, ar.Sandbox.Kind)
	}
}

func ch13ReadSandboxFixtures(t *testing.T) map[string]any {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch13.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]any
	if err := json.Unmarshal(body, &fixtures); err != nil {
		t.Fatal(err)
	}
	return fixtures
}

func ch13JSONValue(value any) any {
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

func ch13ConfigShape(value any) any {
	switch value := value.(type) {
	case string:
		return "string"
	case map[string]any:
		result := make(map[string]any, len(value))
		for key, child := range value {
			result[key] = ch13ConfigShape(child)
		}
		return result
	case []any:
		result := make([]any, len(value))
		for index, child := range value {
			result[index] = ch13ConfigShape(child)
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

func TestCh13StoryboardExamples(t *testing.T) {
	ch05RunStoryboardExamples(t, "docs/storyboards/STORYBOARD_CH13.md", `(?s)<!-- go-example id=(e13[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n`+"```go\n"+`(.*?)\n`+"```", 8, []string{"13-1-recall", "13-2-timeouts", "13-3-first-wins", "13-4-chain-of-command"})
}
