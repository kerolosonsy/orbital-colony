package levels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

var ch35ExamplePattern = regexp.MustCompile("(?s)<!-- go-example id=(e35[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n```go\n(.*?)\n```")
var ch35StdoutPattern = regexp.MustCompile("(?s)Real stdout:\n\n```text\n(.*?)\n```")

func TestCh35Lessons(t *testing.T) {
	idToPrefix := map[string]string{
		"35-1-task-queue":    "c351.",
		"35-2-scheduler":     "c352.",
		"35-3-redis-lock":    "c353.",
		"35-4-outbox":        "c354.",
		"35-5-kafka":         "c355.",
		"35-6-raft-election": "c356.",
		"35-7-raft-log":      "c357.",
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seenLevels := map[string]bool{}
	for _, level := range All {
		if level.Chapter != 35 {
			continue
		}
		prefix, ok := idToPrefix[level.ID]
		if !ok {
			t.Errorf("unexpected chapter 35 level %q", level.ID)
			continue
		}
		seenLevels[level.ID] = true
		copy := *level
		copy.Steps = ch35WithoutSandboxes(level.Steps)
		english := *level.En
		english.Steps = ch35WithoutSandboxes(level.En.Steps)
		copy.En = &english
		ch34ValidateLesson(t, &copy, prefix, arabicTitles, englishTitles)
	}
	for id := range idToPrefix {
		if !seenLevels[id] {
			t.Errorf("chapter 35 lesson is missing for %s", id)
		}
	}
}

func ch35WithoutSandboxes(steps []LessonStep) []LessonStep {
	copy := make([]LessonStep, len(steps))
	for i, step := range steps {
		copy[i] = step
		copy[i].Blocks = make([]LessonBlock, 0, len(step.Blocks))
		for _, block := range step.Blocks {
			if block.Kind != "sandbox" {
				copy[i].Blocks = append(copy[i].Blocks, block)
			}
		}
	}
	return copy
}

func TestCh35Sandboxes(t *testing.T) {
	type fixtureLevel struct {
		ID      string                     `json:"id"`
		Kind    string                     `json:"kind"`
		Configs map[string]json.RawMessage `json:"configs"`
	}
	var fixture struct {
		Levels []fixtureLevel `json:"levels"`
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch35.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(body, &fixture); err != nil {
		t.Fatal(err)
	}
	fixtures := make(map[string]fixtureLevel, len(fixture.Levels))
	for _, item := range fixture.Levels {
		if _, exists := fixtures[item.ID]; exists {
			t.Errorf("duplicate sandbox fixture for %s", item.ID)
		}
		fixtures[item.ID] = item
	}
	if len(fixtures) != 7 {
		t.Fatalf("expected seven chapter 35 fixtures, found %d", len(fixtures))
	}
	for _, id := range []string{
		"35-1-task-queue", "35-2-scheduler", "35-3-redis-lock", "35-4-outbox",
		"35-5-kafka", "35-6-raft-election", "35-7-raft-log",
	} {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing bilingual lesson", id)
		}
		ar, arCount, arMisplaced := ch35FindSandbox(level.Steps)
		en, enCount, enMisplaced := ch35FindSandbox(level.En.Steps)
		if arCount != 1 || enCount != 1 || ar == nil || en == nil || arMisplaced || enMisplaced || ar.Sandbox == nil || en.Sandbox == nil {
			t.Errorf("%s: want one sandbox after the question and walkthrough in both languages", id)
			continue
		}
		if ar.ID != en.ID || ar.Kind != en.Kind || ar.Sandbox.Kind != en.Sandbox.Kind {
			t.Errorf("%s: sandbox identities differ across languages", id)
		}
		if ar.Kind != "sandbox" || !strings.HasPrefix(ar.Sandbox.Kind, "ch35-") || ar.Sandbox.Config == nil {
			t.Errorf("%s: sandbox must use a ch35 kind and Config", id)
		}
		if !ch35SameConfigShape(t, ar.Sandbox.Config, en.Sandbox.Config) {
			t.Errorf("%s: Arabic and English sandbox Config structure differs", id)
		}
		item, ok := fixtures[id]
		if !ok {
			t.Errorf("%s: missing JavaScript fixture", id)
			continue
		}
		if item.Kind != ar.Sandbox.Kind || !ch35JSONEqual(t, ar.Sandbox.Config, item.Configs["ar"]) || !ch35JSONEqual(t, en.Sandbox.Config, item.Configs["en"]) {
			t.Errorf("%s: Go sandbox Config differs from JavaScript fixture", id)
		}
		delete(fixtures, id)
	}
	for id := range fixtures {
		t.Errorf("unexpected chapter 35 sandbox fixture %s", id)
	}
}

func ch35FindSandbox(steps []LessonStep) (*LessonBlock, int, bool) {
	var found *LessonBlock
	count := 0
	misplaced := false
	for stepIndex, step := range steps {
		for blockIndex := range step.Blocks {
			block := &step.Blocks[blockIndex]
			if block.Kind != "sandbox" {
				continue
			}
			count++
			if stepIndex != 2 || blockIndex != len(step.Blocks)-1 {
				misplaced = true
			}
			found = block
		}
	}
	return found, count, misplaced
}

func ch35JSONEqual(t *testing.T, got any, wantJSON json.RawMessage) bool {
	t.Helper()
	var want, normalized any
	if err := json.Unmarshal(wantJSON, &want); err != nil {
		t.Errorf("decode fixture Config: %v", err)
		return false
	}
	body, err := json.Marshal(got)
	if err != nil {
		t.Errorf("encode sandbox Config: %v", err)
		return false
	}
	if err := json.Unmarshal(body, &normalized); err != nil {
		t.Errorf("normalize sandbox Config: %v", err)
		return false
	}
	return reflect.DeepEqual(normalized, want)
}

func ch35SameConfigShape(t *testing.T, first, second any) bool {
	t.Helper()
	normalize := func(value any) (any, bool) {
		body, err := json.Marshal(value)
		if err != nil {
			t.Errorf("encode sandbox Config: %v", err)
			return nil, false
		}
		var normalized any
		if err := json.Unmarshal(body, &normalized); err != nil {
			t.Errorf("normalize sandbox Config: %v", err)
			return nil, false
		}
		return normalized, true
	}
	left, ok := normalize(first)
	if !ok {
		return false
	}
	right, ok := normalize(second)
	if !ok {
		return false
	}
	return ch35JSONShape(left, right)
}

func ch35JSONShape(first, second any) bool {
	switch left := first.(type) {
	case map[string]any:
		right, ok := second.(map[string]any)
		if !ok || len(left) != len(right) {
			return false
		}
		for key, value := range left {
			other, ok := right[key]
			if !ok || !ch35JSONShape(value, other) {
				return false
			}
		}
		return true
	case []any:
		right, ok := second.([]any)
		if !ok || len(left) != len(right) {
			return false
		}
		for index, value := range left {
			if !ch35JSONShape(value, right[index]) {
				return false
			}
		}
		return true
	case string:
		_, ok := second.(string)
		return ok
	default:
		return reflect.DeepEqual(first, second)
	}
}

func TestCh35LessonIDsGlobal(t *testing.T) {
	seen := map[string]string{}
	for _, level := range All {
		if level.Chapter != 34 && level.Chapter != 35 {
			continue
		}
		if level.Kind != "lab" && !level.Boss {
			t.Errorf("chapter %d level %s is not covered as a lab or boss", level.Chapter, level.ID)
		}
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch34BlockIDs(step.Blocks)...) {
				if owner, exists := seen[id]; exists {
					t.Errorf("lesson ID %q is shared by %s and %s", id, owner, level.ID)
				}
				seen[id] = level.ID
			}
		}
	}
}

func TestCh35StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 35 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH35.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	examples := ch35ExamplePattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 14 {
		t.Fatalf("expected fourteen chapter 35 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch35VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch35VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := 0
	for _, levelID := range []string{
		"35-1-task-queue", "35-2-scheduler", "35-3-redis-lock", "35-4-outbox",
		"35-5-kafka", "35-6-raft-election", "35-7-raft-log",
	} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID != id {
						continue
					}
					found++
					if block.Code != code {
						t.Fatalf("%s: lesson code differs from storyboard", id)
					}
					if strings.Join(block.Imports, ",") != importNames {
						t.Fatalf("%s: lesson imports differ from storyboard", id)
					}
				}
			}
		}
	}
	if found != 2 {
		t.Fatalf("%s: expected both bilingual walkthroughs, found %d", id, found)
	}
	parts := strings.SplitN(storyboard, "#### Verified "+id+" — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing documented output", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	match := ch35StdoutPattern.FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("%s: missing stdout block", id)
	}
	wantOutput := match[1]
	dir := t.TempDir()
	imports := make([]string, 0)
	for _, importName := range strings.Split(importNames, ",") {
		imports = append(imports, fmt.Sprintf("%q", importName))
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + strings.Join(imports, ";") + ")\n" + code + "\n",
	}
	for name, contents := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, goBin, "build", "-trimpath", "-gcflags=-lang=go1.24", "-o", "example", ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build %s: %v\n%s", id, err, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	var stdout, stderr strings.Builder
	run.Stdout, run.Stderr = &stdout, &stderr
	err = run.Run()
	if err != nil || stderr.Len() != 0 {
		t.Fatalf("run %s: %v\n%s", id, err, stderr.String())
	}
	if stdout.String() != wantOutput+"\n" {
		t.Fatalf("%s stdout = %q, documented %q", id, stdout.String(), wantOutput+"\n")
	}
	for _, levelID := range []string{
		"35-1-task-queue", "35-2-scheduler", "35-3-redis-lock", "35-4-outbox",
		"35-5-kafka", "35-6-raft-election", "35-7-raft-log",
	} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != wantOutput {
						t.Errorf("%s: final walkthrough output differs from transcript", id)
					}
				}
			}
		}
	}
}
