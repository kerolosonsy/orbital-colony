package levels

import "testing"

func TestCh31Lessons(t *testing.T) {
	ids := []string{
		"31-1-echo", "31-2-chat", "31-3-discovery",
		"31-4-frames", "31-5-dns", "31-6-port-check",
	}
	prefixes := []string{"c311.", "c312.", "c313.", "c314.", "c315.", "c316."}
	arTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	enTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	ch30CheckGlobalIDs(t)
	ch30ValidateChapter(t, ids, prefixes, arTitles, enTitles)
}

func TestCh31StoryboardExamples(t *testing.T) {
	ch30RunStoryboardExamples(t, "STORYBOARD_CH31.md", "e31", []string{
		"31-1-echo", "31-2-chat", "31-3-discovery",
		"31-4-frames", "31-5-dns", "31-6-port-check",
	}, 12)
}
