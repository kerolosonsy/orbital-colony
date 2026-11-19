package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLearningMetricsRunRecordingAndStaleSaveMerge(t *testing.T) {
	save := progressSave{"metrics": json.RawMessage(`{"1-1":{"runs":0,"activeTimeMs":94000,"bestTicks":18}}`)}
	for _, response := range []*runResponse{
		{Phase: "compile", Result: result{Success: false}},
		{Phase: "run", Failure: "runtime", Result: result{Success: false}},
		{Phase: "run", Failure: "check", Result: result{Success: false}},
		{Phase: "run", Result: result{Success: true, Ticks: 12}},
	} {
		recordLearningRun(save, "1-1", response)
	}
	metrics := learningMetricsFromSave(save)["1-1"]
	if metrics.Runs != 4 || metrics.FailedRuns != (learningFailures{Compile: 1, Runtime: 1, Check: 1}) {
		t.Fatalf("run totals: %+v", metrics)
	}
	if !metrics.Solved || metrics.TimeToFirstSolveMs != 94000 || metrics.BestTicks != 12 {
		t.Fatalf("solve metrics: %+v", metrics)
	}
	stale := progressSave{"metrics": json.RawMessage(`{"1-1":{"runs":2,"failedRuns":{"compile":1},"hintCount":2}}`)}
	mergeLearningMetrics(save, stale)
	merged := learningMetricsFromSave(stale)["1-1"]
	if merged.Runs != 4 || merged.HintCount != 2 || !merged.Solved || merged.BestTicks != 12 {
		t.Fatalf("stale save lost newer metrics: %+v", merged)
	}
}

func TestLearningMetricsMigrationDefaults(t *testing.T) {
	path := filepath.Join(t.TempDir(), "progress.json")
	if err := os.WriteFile(path, []byte(`{"metrics":{"1-1":{"runs":3,"hintsRevealed":["0","0","1"],"lessonStepsViewed":{"intro":2}}}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	save, err := readProgress(path)
	if err != nil {
		t.Fatal(err)
	}
	entry := learningMetricsFromSave(save)["1-1"]
	if entry.Runs != 3 || entry.FailedRuns != (learningFailures{}) || entry.HintCount != 2 || len(entry.HintsRevealed) != 2 {
		t.Fatalf("old metrics were not normalized: %+v", entry)
	}
	if entry.LessonStepsViewed["intro"] != 2 || entry.LessonStepsViewed["missing"] != 0 || entry.ActiveTimeMs != 0 || entry.BestTicks != 0 || entry.Solved {
		t.Fatalf("missing fields should default to zero: %+v", entry)
	}
}

func TestLearningMetricsEndpointReturnsProfileMetricsOnly(t *testing.T) {
	s, h := profileServer(t)
	if err := saveProgress(filepath.Join(s.profiles.dir("default"), "progress.json"), progressSave{
		"metrics": json.RawMessage(`{"1-1":{"runs":2,"failedRuns":{"check":1},"hintCount":1,"hintsRevealed":["0"]}}`),
	}); err != nil {
		t.Fatal(err)
	}
	rec := profileRequest(t, h, "GET", "/api/profiles/default/metrics", "default", "")
	if rec.Code != 200 {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body)
	}
	var metrics learningMetrics
	if err := json.Unmarshal(rec.Body.Bytes(), &metrics); err != nil {
		t.Fatal(err)
	}
	entry := metrics["1-1"]
	if entry.Runs != 2 || entry.FailedRuns.Check != 1 || entry.HintCount != 1 || len(entry.HintsRevealed) != 1 {
		t.Fatalf("endpoint metrics: %+v", entry)
	}
	if strings.Contains(rec.Body.String(), "Player 1") || strings.Contains(rec.Body.String(), "profile") {
		t.Fatalf("endpoint included profile information: %s", rec.Body)
	}
}

func TestLearningMetricsStatsTableAggregatesProfiles(t *testing.T) {
	root := t.TempDir()
	for id, save := range map[string]progressSave{
		"first":  {"metrics": json.RawMessage(`{"1-1":{"runs":2,"failedRuns":{"check":1},"hintCount":2,"activeTimeMs":1}}`)},
		"second": {"metrics": json.RawMessage(`{"1-1":{"runs":1,"failedRuns":{"runtime":1},"hintCount":1,"activeTimeMs":61000},"1-2":{"runs":1,"failedRuns":{"compile":1}}}`)},
	} {
		if err := saveProgress(filepath.Join(root, "profiles", id, "progress.json"), save); err != nil {
			t.Fatal(err)
		}
	}
	var table strings.Builder
	if err := printLearningStats(root, &table); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(table.String()), "\n")
	if len(lines) != 3 || !strings.Contains(lines[1], "1-1") || !strings.Contains(lines[2], "1-2") {
		t.Fatalf("unexpected top-level stats table: %s", table.String())
	}
}
