package levels

import (
	"fmt"
	"testing"
)

func TestCh13Lessons(t *testing.T) {
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
		ch05ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, arabicTitles, prefix, 200, 320)
		ch05ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, englishTitles, prefix, 200, 320)
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
		ch05ValidateParity(t, id, level.Steps, level.En.Steps)
	}
}

func TestCh13StoryboardExamples(t *testing.T) {
	ch05RunStoryboardExamples(t, "STORYBOARD_CH13.md", `(?s)<!-- go-example id=(e13[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n`+"```go\n"+`(.*?)\n`+"```", 8, []string{"13-1-recall", "13-2-timeouts", "13-3-first-wins", "13-4-chain-of-command"})
}
