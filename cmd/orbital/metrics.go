package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
)

type learningFailures struct {
	Compile int `json:"compile"`
	Runtime int `json:"runtime"`
	Check   int `json:"check"`
}

type levelLearningMetrics struct {
	Runs               int              `json:"runs"`
	FailedRuns         learningFailures `json:"failedRuns"`
	HintCount          int              `json:"hintCount"`
	HintsRevealed      []string         `json:"hintsRevealed"`
	LessonsOpened      int              `json:"lessonsOpened"`
	LessonStepsViewed  map[string]int   `json:"lessonStepsViewed"`
	Interactions       int              `json:"interactions"`
	ActiveTimeMs       int64            `json:"activeTimeMs"`
	Solved             bool             `json:"solved"`
	TimeToFirstSolveMs int64            `json:"timeToFirstSolveMs"`
	BestTicks          int              `json:"bestTicks"`
}

type learningMetrics map[string]levelLearningMetrics

type struggleMetric struct {
	Level  string
	Score  int
	Fails  int
	Hints  int
	Active int64
	Solved bool
}

func learningMetricsFromSave(save progressSave) learningMetrics {
	metrics := learningMetrics{}
	if len(save["metrics"]) == 0 {
		return metrics
	}
	if err := json.Unmarshal(save["metrics"], &metrics); err != nil || metrics == nil {
		return learningMetrics{}
	}
	for id, entry := range metrics {
		metrics[id] = normalizeLevelLearningMetrics(entry)
	}
	return metrics
}

func normalizeLevelLearningMetrics(metric levelLearningMetrics) levelLearningMetrics {
	metric.Runs = max(0, metric.Runs)
	metric.FailedRuns = learningFailures{Compile: max(0, metric.FailedRuns.Compile), Runtime: max(0, metric.FailedRuns.Runtime), Check: max(0, metric.FailedRuns.Check)}
	metric.HintCount = max(0, metric.HintCount)
	metric.HintsRevealed = uniqueMetricStrings(metric.HintsRevealed)
	metric.HintCount = max(metric.HintCount, len(metric.HintsRevealed))
	metric.LessonsOpened = max(0, metric.LessonsOpened)
	metric.Interactions = max(0, metric.Interactions)
	metric.ActiveTimeMs = max(int64(0), metric.ActiveTimeMs)
	metric.TimeToFirstSolveMs = max(int64(0), metric.TimeToFirstSolveMs)
	metric.BestTicks = max(0, metric.BestTicks)
	if metric.LessonStepsViewed == nil {
		metric.LessonStepsViewed = map[string]int{}
	}
	for step, views := range metric.LessonStepsViewed {
		metric.LessonStepsViewed[step] = max(0, views)
	}
	return metric
}

func uniqueMetricStrings(hintIDs []string) []string {
	seen := map[string]bool{}
	unique := make([]string, 0, len(hintIDs))
	for _, hintID := range hintIDs {
		if hintID != "" && !seen[hintID] {
			seen[hintID] = true
			unique = append(unique, hintID)
		}
	}
	return unique
}

func migrateLearningMetrics(save progressSave) {
	if len(save["metrics"]) == 0 {
		return
	}
	blob, err := json.Marshal(learningMetricsFromSave(save))
	if err == nil {
		save["metrics"] = blob
	}
}

func setLearningMetrics(save progressSave, metrics learningMetrics) {
	blob, err := json.Marshal(metrics)
	if err == nil {
		save["metrics"] = blob
	}
}

func learningFailureKind(response *runResponse) string {
	if response.Phase == "compile" {
		return "compile"
	}
	if response.Failure == "runtime" || response.Race {
		return "runtime"
	}
	return "check"
}

func recordLearningRun(save progressSave, levelID string, response *runResponse) {
	metrics := learningMetricsFromSave(save)
	entry := normalizeLevelLearningMetrics(metrics[levelID])
	entry.Runs++
	if !response.Result.Success {
		switch learningFailureKind(response) {
		case "compile":
			entry.FailedRuns.Compile++
		case "runtime":
			entry.FailedRuns.Runtime++
		default:
			entry.FailedRuns.Check++
		}
	} else {
		if !entry.Solved {
			entry.TimeToFirstSolveMs = entry.ActiveTimeMs
		}
		entry.Solved = true
		if entry.BestTicks == 0 || response.Result.Ticks < entry.BestTicks {
			entry.BestTicks = max(0, response.Result.Ticks)
		}
	}
	metrics[levelID] = entry
	setLearningMetrics(save, metrics)
}

func mergeLearningMetrics(old, incoming progressSave) {
	if len(old["metrics"]) == 0 && len(incoming["metrics"]) == 0 {
		return
	}
	merged := learningMetricsFromSave(old)
	for levelID, entry := range learningMetricsFromSave(incoming) {
		merged[levelID] = mergeLevelLearningMetrics(merged[levelID], entry)
	}
	setLearningMetrics(incoming, merged)
}

func mergeLevelLearningMetrics(current, newer levelLearningMetrics) levelLearningMetrics {
	current = normalizeLevelLearningMetrics(current)
	newer = normalizeLevelLearningMetrics(newer)
	current.Runs = max(current.Runs, newer.Runs)
	current.FailedRuns = learningFailures{Compile: max(current.FailedRuns.Compile, newer.FailedRuns.Compile), Runtime: max(current.FailedRuns.Runtime, newer.FailedRuns.Runtime), Check: max(current.FailedRuns.Check, newer.FailedRuns.Check)}
	current.HintCount = max(current.HintCount, newer.HintCount)
	current.HintsRevealed = uniqueMetricStrings(append(current.HintsRevealed, newer.HintsRevealed...))
	current.HintCount = max(current.HintCount, len(current.HintsRevealed))
	current.LessonsOpened = max(current.LessonsOpened, newer.LessonsOpened)
	current.Interactions = max(current.Interactions, newer.Interactions)
	current.ActiveTimeMs = max(current.ActiveTimeMs, newer.ActiveTimeMs)
	current.Solved = current.Solved || newer.Solved
	current.TimeToFirstSolveMs = max(current.TimeToFirstSolveMs, newer.TimeToFirstSolveMs)
	if current.BestTicks == 0 || newer.BestTicks > 0 && newer.BestTicks < current.BestTicks {
		current.BestTicks = newer.BestTicks
	}
	for step, views := range newer.LessonStepsViewed {
		current.LessonStepsViewed[step] = max(current.LessonStepsViewed[step], views)
	}
	return current
}

func (s *server) recordProfileLearningRun(p *profileStore, profileID, levelID string, response *runResponse) {
	if p == nil {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, err := p.load(profileID); err != nil {
		log.Printf("load profile for learning metrics: %v", err)
		return
	}
	path := filepath.Join(p.dir(profileID), "progress.json")
	save, err := readProgress(path)
	if err != nil {
		log.Printf("read profile learning metrics: %v", err)
		return
	}
	recordLearningRun(save, levelID, response)
	if err := saveProgress(path, save); err != nil {
		log.Printf("save learning metrics: %v", err)
	}
}

func (s *server) profileMetrics(w http.ResponseWriter, r *http.Request) {
	p := s.store(r)
	p.mu.Lock()
	defer p.mu.Unlock()
	pr, ok := s.profileByPath(w, r)
	if !ok {
		return
	}
	save, err := readProgress(filepath.Join(p.dir(pr.ID), "progress.json"))
	if err != nil {
		profileError(w, r, err)
		return
	}
	writeJSON(w, learningMetricsFromSave(save))
}

func struggleScore(metrics levelLearningMetrics) int {
	fails := metrics.FailedRuns.Compile + metrics.FailedRuns.Runtime + metrics.FailedRuns.Check
	minutes := int(activeMinutes(metrics.ActiveTimeMs))
	return fails*4 + metrics.HintCount*2 + minutes
}

func activeMinutes(milliseconds int64) int64 {
	minutes := milliseconds / 60000
	if milliseconds%60000 > 0 {
		minutes++
	}
	return minutes
}

func topStruggleMetrics(metrics learningMetrics, limit int) []struggleMetric {
	rows := make([]struggleMetric, 0, len(metrics))
	for id, entry := range metrics {
		entry = normalizeLevelLearningMetrics(entry)
		score := struggleScore(entry)
		if score == 0 {
			continue
		}
		rows = append(rows, struggleMetric{Level: id, Score: score, Fails: entry.FailedRuns.Compile + entry.FailedRuns.Runtime + entry.FailedRuns.Check,
			Hints: entry.HintCount, Active: entry.ActiveTimeMs, Solved: entry.Solved})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Score != rows[j].Score {
			return rows[i].Score > rows[j].Score
		}
		return rows[i].Level < rows[j].Level
	})
	if len(rows) > limit {
		rows = rows[:limit]
	}
	return rows
}

func printLearningStats(savesDir string, out io.Writer) error {
	metrics, err := aggregateProfileMetrics(savesDir)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(out, "#    LEVEL                SCORE  FAILS  HINTS  ACTIVE(min)  SOLVED"); err != nil {
		return err
	}
	for i, row := range topStruggleMetrics(metrics, 20) {
		_, err := fmt.Fprintf(out, "%-4d %-20s %-6d %-6d %-6d %-12d %t\n", i+1, row.Level, row.Score, row.Fails, row.Hints, activeMinutes(row.Active), row.Solved)
		if err != nil {
			return err
		}
	}
	return nil
}

func aggregateProfileMetrics(savesDir string) (learningMetrics, error) {
	paths, err := filepath.Glob(filepath.Join(savesDir, "profiles", "*", "progress.json"))
	if err != nil {
		return nil, err
	}
	accountPaths, err := filepath.Glob(filepath.Join(savesDir, "users", "*", "profiles", "*", "progress.json"))
	if err != nil {
		return nil, err
	}
	paths = append(paths, accountPaths...)
	all := learningMetrics{}
	for _, path := range paths {
		blob, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var save progressSave
		if err := json.Unmarshal(blob, &save); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		for levelID, entry := range learningMetricsFromSave(save) {
			all[levelID] = addLevelLearningMetrics(all[levelID], entry)
		}
	}
	return all, nil
}

func addLevelLearningMetrics(aggregate, profileMetrics levelLearningMetrics) levelLearningMetrics {
	aggregate = normalizeLevelLearningMetrics(aggregate)
	profileMetrics = normalizeLevelLearningMetrics(profileMetrics)
	aggregate.Runs += profileMetrics.Runs
	aggregate.FailedRuns = learningFailures{Compile: aggregate.FailedRuns.Compile + profileMetrics.FailedRuns.Compile, Runtime: aggregate.FailedRuns.Runtime + profileMetrics.FailedRuns.Runtime, Check: aggregate.FailedRuns.Check + profileMetrics.FailedRuns.Check}
	aggregate.HintCount += profileMetrics.HintCount
	aggregate.HintsRevealed = uniqueMetricStrings(append(aggregate.HintsRevealed, profileMetrics.HintsRevealed...))
	aggregate.LessonsOpened += profileMetrics.LessonsOpened
	aggregate.Interactions += profileMetrics.Interactions
	aggregate.ActiveTimeMs += profileMetrics.ActiveTimeMs
	aggregate.Solved = aggregate.Solved || profileMetrics.Solved
	for step, views := range profileMetrics.LessonStepsViewed {
		aggregate.LessonStepsViewed[step] += views
	}
	return aggregate
}
