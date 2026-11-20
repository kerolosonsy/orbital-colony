package levels

import (
	"strings"
	"testing"
)

// The client renders ~code~ one line at a time, so a code span must open and
// close on the same line. A lone "~ " (the generics tilde) and table rows
// (approximate numbers like "~10,000") are fine.
func TestInlineCodeSpansCloseOnTheirLine(t *testing.T) {
	check := func(id, where, text string) {
		fenced := false
		for _, line := range strings.Split(text, "\n") {
			if strings.HasPrefix(line, "~~~") || strings.HasPrefix(line, "```") {
				fenced = !fenced
				continue
			}
			if fenced || strings.HasPrefix(line, "|") {
				continue
			}
			// a span closes on the same line; one standalone " ~ " may stand for
			// the generics tilde
			total, lone := strings.Count(line, "~"), 0
			for i := 0; i < len(line); i++ {
				if line[i] == '~' && (i == 0 || line[i-1] == ' ') && (i+1 == len(line) || line[i+1] == ' ') {
					lone++
				}
			}
			opens := total
			if total%2 != 0 {
				opens = total - lone
			}
			if opens%2 != 0 {
				t.Errorf("%s %s: unclosed ~code~ span in %q", id, where, line)
			}
			if strings.Count(line, "**")%2 != 0 {
				t.Errorf("%s %s: unclosed **bold** in %q", id, where, line)
			}
		}
	}
	scan := func(id, lang, goal, lesson string, hints []string, steps []LessonStep) {
		check(id, lang+" goal", goal)
		check(id, lang+" lesson", lesson)
		for i, hint := range hints {
			check(id, lang+" hint "+string(rune('1'+i)), hint)
		}
		for _, step := range steps {
			for _, block := range step.Blocks {
				check(id, lang+" "+block.ID, block.Text)
				check(id, lang+" "+block.ID+" note", block.Note)
			}
		}
	}
	for _, l := range All {
		scan(l.ID, "ar", l.Goal, l.Lesson, l.Hints, l.Steps)
		if l.En != nil {
			scan(l.ID, "en", l.En.Goal, l.En.Lesson, l.En.Hints, l.En.Steps)
		}
	}
}
