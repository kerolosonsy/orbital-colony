package levels

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

func TestEveryChapterHasFiveBilingualQuizQuestions(t *testing.T) {
	for _, chapter := range Chapters {
		t.Run(fmt.Sprint(chapter.N), func(t *testing.T) {
			if len(chapter.Quiz) != 5 {
				t.Fatalf("got %d questions, want 5", len(chapter.Quiz))
			}
			seen := map[string]bool{}
			for i, q := range chapter.Quiz {
				if seen[q.Question] {
					t.Errorf("question %d repeats a question", i)
				}
				seen[q.Question] = true
				for lang, text := range map[string]QuizText{"ar": q.QuizText, "en": q.En} {
					if strings.TrimSpace(text.Question) == "" || strings.TrimSpace(text.Explanation) == "" {
						t.Errorf("question %d %s: missing question or explanation", i, lang)
					}
					if len(text.Options) < 3 || len(text.Options) > 4 {
						t.Errorf("question %d %s: got %d options", i, lang, len(text.Options))
					}
					if q.Answer < 0 || q.Answer >= len(text.Options) {
						t.Errorf("question %d %s: invalid answer %d", i, lang, q.Answer)
					}
					options := map[string]bool{}
					for _, option := range text.Options {
						if strings.TrimSpace(option) == "" || options[option] {
							t.Errorf("question %d %s: empty or duplicate option", i, lang)
						}
						options[option] = true
					}
				}
				if len(q.Options) != len(q.En.Options) {
					t.Errorf("question %d: option counts differ", i)
				}
				if !strings.ContainsFunc(q.Question, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) {
					t.Errorf("question %d: Arabic question missing", i)
				}
				if strings.ContainsFunc(q.En.Question+q.En.Explanation+strings.Join(q.En.Options, ""), func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) {
					t.Errorf("question %d: Arabic text in English version", i)
				}
			}
		})
	}
}
