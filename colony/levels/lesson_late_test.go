package levels

import (
	"fmt"
	"strings"
	"testing"
)

func TestLateLessonsStayShortAndVisual(t *testing.T) {
	for _, level := range []*Level{Get("6-4-stock-order"), Get("6-5-heat-grid")} {
		t.Run(level.ID, func(t *testing.T) {
			if level.LessonRevision != 1 || level.En == nil || level.En.LessonRevision != 1 {
				t.Fatal("late lesson revision must be one in both languages")
			}
			if len(level.Steps) != 6 || len(level.En.Steps) != 6 || len(level.Hints) != 3 || len(level.En.Hints) != 3 {
				t.Fatal("late lesson must have six steps in both languages")
			}
			arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
			englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
			for i := range arabicTitles {
				if level.Steps[i].Title != arabicTitles[i] || level.En.Steps[i].Title != englishTitles[i] {
					t.Fatalf("step %d title does not match the approved six-step shape", i)
				}
			}
			for _, copy := range []struct {
				steps   []LessonStep
				hints   []string
				english bool
			}{{level.Steps, level.Hints, false}, {level.En.Steps, level.En.Hints, true}} {
				words, questions, walkthroughs, sandboxes := lessonLateWords(copy.steps, copy.hints, level.ID, copy.english), 0, 0, 0
				for i, step := range copy.steps {
					if i == 0 && len(strings.Fields(step.Blocks[0].Text)) > 45 {
						t.Fatal("mission step exceeds 45 words")
					}
					for _, block := range step.Blocks {
						if block.Note != "" || block.Optional {
							t.Fatal("late lesson must not hide extra reading")
						}
						if block.Walkthrough != nil {
							walkthroughs++
							for _, frame := range block.Walkthrough.Frames {
								if len(strings.Fields(frame.Caption)) > 8 {
									t.Fatalf("%s caption exceeds eight words: %q", block.ID, frame.Caption)
								}
							}
						}
						if block.Question != nil {
							questions++
							for _, feedback := range block.Question.Feedback {
								if !strings.HasSuffix(feedback, ".") {
									t.Fatalf("%s feedback must be one sentence: %q", block.ID, feedback)
								}
							}
						}
						if block.Sandbox != nil {
							sandboxes++
							if strings.Contains(block.Text, "\n") {
								t.Fatal("sandbox instruction must be one line")
							}
						}
					}
				}
				t.Logf("learner-authored words, including hints and sandbox labels: %d", words)
				if words < 250 || words > 320 || questions != 3 || walkthroughs != 1 || sandboxes != 1 {
					t.Fatalf("words=%d questions=%d walkthroughs=%d sandboxes=%d", words, questions, walkthroughs, sandboxes)
				}
			}
		})
	}
}

func lessonLateWords(steps []LessonStep, hints []string, levelID string, english bool) int {
	var text []string
	text = append(text, hints...)
	for _, step := range steps {
		text = append(text, step.Title)
		for _, block := range step.Blocks {
			text = append(text, block.Text)
			if block.Walkthrough != nil {
				for _, frame := range block.Walkthrough.Frames {
					text = append(text, frame.Caption)
				}
			}
			if block.Question != nil {
				text = append(text, block.Question.Options...)
				text = append(text, block.Question.Feedback...)
			}
			if block.Sandbox == nil {
				continue
			}
			spec := block.Sandbox
			text = append(text, spec.SourceLabel, spec.ResetLabel)
			if spec.Kind == "sort-compare" {
				text = append(text, spec.SortDirectionLabel, spec.SortTieLabel, spec.AscendingLabel, spec.DescendingLabel, spec.TieOnLabel, spec.TieOffLabel)
				for position, size := range spec.Sequence {
					if english {
						text = append(text, fmt.Sprintf("Position %d, size %d", position, size))
					} else {
						text = append(text, fmt.Sprintf("الموقع %d، الحجم %d", position, size))
					}
				}
			}
			if spec.Kind == "grid-index" {
				text = append(text, spec.GridOuterCloneLabel, spec.GridRowCloneLabel, spec.GridEditCopyLabel)
				if english {
					text = append(text, "Original", "Copy", "Original grid", "Cloned grid")
				} else {
					text = append(text, "الأصل", "النسخة", "الشبكة الأصلية", "الشبكة المنسوخة")
				}
				for row, values := range spec.Grid {
					for column, value := range values {
						if english {
							text = append(text, fmt.Sprintf("grid row %d column %d: %d", row, column, value))
						} else {
							text = append(text, fmt.Sprintf("grid صف %d عمود %d: %d", row, column, value))
						}
					}
				}
			}
		}
	}
	words := 0
	for _, line := range text {
		words += len(strings.Fields(line))
	}
	return words
}

func TestLateMissionGradingIsUnchanged(t *testing.T) {
	stock := Get("6-4-stock-order")
	if stock.Par != 19 || stock.Gold != 19 || stock.MaxTicks != 50 {
		t.Fatal("6-4 grading changed")
	}
	heat := Get("6-5-heat-grid")
	if heat.Kind != "lab" || heat.Requires[0] != "slice" {
		t.Fatal("6-5 lab contract changed")
	}
}
