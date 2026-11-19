package levels

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
		ch34ValidateLesson(t, level, prefix, arabicTitles, englishTitles)
	}
	for id := range idToPrefix {
		if !seenLevels[id] {
			t.Errorf("chapter 35 lesson is missing for %s", id)
		}
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
