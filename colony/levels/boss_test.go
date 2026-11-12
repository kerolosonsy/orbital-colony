package levels_test

import (
	"orbital/colony/levels"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestTrackBossesCloseTheirTracks(t *testing.T) {
	want := map[int]string{9: "Fundamentals", 17: "Craft (generics, algorithms, tests)", 23: "Data and Tools", 25: "Concurrency", 28: "Services", 30: "Databases", 33: "Networking", 36: "Systems"}
	last := map[string]*levels.Level{}
	tracks := map[int]string{}
	for _, chapter := range levels.Chapters {
		tracks[chapter.N] = chapter.TrackEn
	}
	for _, level := range levels.All {
		last[tracks[level.Chapter]] = level
	}
	found := map[int]bool{}
	missionID := regexp.MustCompile(`[0-9]+-[0-9]+-[a-z][a-z-]+`)
	for _, level := range levels.All {
		if !level.Boss {
			continue
		}
		track, ok := want[level.Chapter]
		if !ok || found[level.Chapter] {
			t.Fatalf("unexpected/duplicate boss: %s", level.ID)
		}
		found[level.Chapter] = true
		if last[track] != level {
			t.Errorf("%s does not close %s", level.ID, track)
		}
		if !strings.HasPrefix(level.ID, strconv.Itoa(level.Chapter)+"-boss-") || level.Kind != "lab" || level.Service != "" {
			t.Errorf("invalid/offline-inaccessible boss: %s", level.ID)
		}
		if level.Lab == nil || !level.Lab.Bonus || level.Lab.ReqFile == "" || len(level.Bridge) == 0 || len(level.Hints) != 3 || level.En == nil || len(level.En.Hints) != 3 {
			t.Errorf("incomplete boss: %s", level.ID)
		}
		if level.En == nil {
			continue
		}
		for _, lesson := range []string{level.Lesson, level.En.Lesson} {
			for _, id := range missionID.FindAllString(lesson, -1) {
				earlier := levels.Get(id)
				if earlier == nil || earlier.Chapter > level.Chapter {
					t.Errorf("%s references untaught or missing mission %s", level.ID, id)
				}
			}
		}
	}
	for chapter := range want {
		if !found[chapter] {
			t.Errorf("chapter %d missing boss", chapter)
		}
	}
	if last["Concurrency"] != nil && !last["Concurrency"].Race {
		t.Error("concurrency boss must use race detector")
	}
}
