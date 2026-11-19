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

func TestCh34SandboxFixtures(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch34.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]map[string]any
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	ids := []string{"34-1-kv-store", "34-2-wal", "34-3-compaction", "34-4-chunk-store"}
	if len(fixtures) != len(ids) {
		t.Fatalf("fixture count = %d, want %d", len(fixtures), len(ids))
	}
	for _, id := range ids {
		level := Get(id)
		block := level.Steps[2].Blocks[2]
		if block.Sandbox == nil || block.Sandbox.Kind != "ch34-storage-flow" || !reflect.DeepEqual(block.Sandbox.Config, fixtures[id]) {
			t.Errorf("%s sandbox differs from its JavaScript fixture", id)
		}
	}
}

var ch34ExamplePattern = regexp.MustCompile("(?s)<!-- go-example id=(e34[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n```go\n(.*?)\n```")
var ch34StdoutPattern = regexp.MustCompile("(?s)Real stdout:\n\n```text\n(.*?)\n```")
var ch34InlineCodePattern = regexp.MustCompile("~[^~]*~")

func TestCh34Lessons(t *testing.T) {
	idToPrefix := map[string]string{
		"34-1-kv-store":    "c341.",
		"34-2-wal":         "c342.",
		"34-3-compaction":  "c343.",
		"34-4-chunk-store": "c344.",
	}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	seenLevels := map[string]bool{}
	for _, level := range All {
		if level.Chapter != 34 {
			continue
		}
		prefix, ok := idToPrefix[level.ID]
		if !ok {
			t.Errorf("unexpected chapter 34 level %q", level.ID)
			continue
		}
		seenLevels[level.ID] = true
		ch34ValidateLesson(t, level, prefix, arabicTitles, englishTitles)
	}
	for id := range idToPrefix {
		if !seenLevels[id] {
			t.Errorf("chapter 34 lesson is missing for %s", id)
		}
	}
}

func TestCh34LessonIDsGlobal(t *testing.T) {
	seen := map[string]string{}
	for _, level := range All {
		if level.Chapter != 34 && level.Chapter != 35 {
			continue
		}
		if level.Kind != "lab" && !level.Boss {
			t.Errorf("chapter %d level %s is not covered as a lab or boss", level.Chapter, level.ID)
		}
		if len(level.Steps) == 0 || level.En == nil || len(level.En.Steps) == 0 {
			continue
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

func ch34BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, 0, len(blocks))
	for _, block := range blocks {
		ids = append(ids, block.ID)
	}
	return ids
}

func ch34ValidateLesson(t *testing.T, level *Level, prefix string, arabicTitles, englishTitles []string) {
	t.Helper()
	if level.En == nil || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
		t.Fatalf("%s: missing English copy or revision 1 wiring", level.ID)
	}
	if len(level.Steps) != 6 || len(level.En.Steps) != 6 || len(level.Hints) != 3 || len(level.En.Hints) != 3 {
		t.Fatalf("%s: expected six steps and three hints in both languages", level.ID)
	}
	ch34ValidateCopy(t, level.ID+" Arabic", level.Steps, level.Hints, prefix, arabicTitles)
	ch34ValidateCopy(t, level.ID+" English", level.En.Steps, level.En.Hints, prefix, englishTitles)
	for i, step := range level.Steps {
		translated := level.En.Steps[i]
		if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
			t.Errorf("%s: Arabic and English step structure differs", step.ID)
			continue
		}
		for j, block := range step.Blocks {
			other := translated.Blocks[j]
			if block.ID != other.ID || block.Kind != other.Kind {
				t.Errorf("%s: Arabic and English block identity differs", block.ID)
			}
			if block.Question != nil && (other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options)) {
				t.Errorf("%s: question answer or option count differs", block.ID)
			}
			if block.Sandbox != nil && (other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind || !reflect.DeepEqual(block.Sandbox.Config, other.Sandbox.Config)) {
				t.Errorf("%s: sandbox config differs between languages", block.ID)
			}
			if block.Walkthrough == nil {
				continue
			}
			if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
				t.Errorf("%s: walkthrough structure differs", block.ID)
				continue
			}
			for k, frame := range block.Walkthrough.Frames {
				translatedFrame := other.Walkthrough.Frames[k]
				frame.Caption, translatedFrame.Caption = "", ""
				if !reflect.DeepEqual(frame, translatedFrame) {
					t.Errorf("%s frame %d: execution state differs", block.ID, k)
				}
			}
		}
	}
}

func ch34ValidateCopy(t *testing.T, name string, steps []LessonStep, hints []string, prefix string, titles []string) {
	t.Helper()
	seen := map[string]bool{}
	claim := func(id, idPrefix string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, idPrefix) {
			t.Errorf("%s: ID %q lacks prefix %q", name, id, idPrefix)
		}
		seen[id] = true
	}
	words := 0
	for i, step := range steps {
		claim(step.ID, prefix)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has the wrong title or no blocks", name, i)
		}
		words += ch34TextWords(step.Title)
		for _, block := range step.Blocks {
			blockPrefix := prefix
			if block.Kind == "sandbox" && !strings.HasPrefix(block.ID, "ch34-") {
				t.Errorf("%s: sandbox ID %q lacks chapter prefix", name, block.ID)
			}
			if block.Kind == "sandbox" {
				blockPrefix = "ch34-"
			}
			claim(block.ID, blockPrefix)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch34TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has an invalid answer, option, or feedback count", name, block.ID)
					continue
				}
				words += ch34TextWords(strings.Join(block.Question.Options, " "))
				words += ch34TextWords(strings.Join(block.Question.Feedback, " "))
				for i := range block.Question.Options {
					if strings.TrimSpace(block.Question.Options[i]) == "" || !ch34OneSentence(block.Question.Feedback[i]) {
						t.Errorf("%s: %s has empty option or multi-sentence feedback", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch34ValidateFrames(t, name, block)
			case "sandbox":
				if block.Sandbox == nil || !strings.HasPrefix(block.Sandbox.Kind, "ch34-") || len(block.Sandbox.Config) == 0 || block.Question != nil || block.Walkthrough != nil || block.Code != "" {
					t.Errorf("%s: %s has incomplete or unprefixed sandbox data", name, block.ID)
				}
			default:
				t.Errorf("%s: %s uses disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	for i, hint := range hints {
		if !ch34OneSentence(hint) {
			t.Errorf("%s: hint %d must be one sentence", name, i+1)
		}
		for _, label := range []string{"Concept:", "Structure:", "Location:", "الفكرة:", "الهيكل:", "الموضع:", "المفهوم:"} {
			if strings.HasPrefix(hint, label) {
				t.Errorf("%s: hint %d starts with a client-supplied label", name, i+1)
			}
		}
		words += ch34TextWords(hint)
	}
	try := steps[2]
	chapter34 := strings.HasPrefix(prefix, "c34")
	wantTryBlocks := 2
	if chapter34 {
		wantTryBlocks = 3
	}
	if len(try.Blocks) != wantTryBlocks || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || (chapter34 && try.Blocks[2].Kind != "sandbox") {
		if chapter34 {
			t.Errorf("%s: Try it must predict, reveal, then sandbox", name)
		} else {
			t.Errorf("%s: Try it must predict before a walkthrough", name)
		}
	} else if strings.TrimSpace(ch34QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
		t.Errorf("%s: prediction code differs from its reveal", name)
	}
	sandboxes := 0
	for _, step := range steps {
		for _, block := range step.Blocks {
			if block.Kind == "sandbox" {
				sandboxes++
			}
		}
	}
	wantSandboxes := 0
	if chapter34 {
		wantSandboxes = 1
	}
	if sandboxes != wantSandboxes {
		t.Errorf("%s: expected exactly %d sandboxes, found %d", name, wantSandboxes, sandboxes)
	}
	if words < 200 || words > 320 {
		t.Errorf("%s: learner text has %d words; want 200-320", name, words)
	}
	t.Logf("%s learner words: %d", name, words)
}

func ch34ValidateFrames(t *testing.T, name string, block LessonBlock) int {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Errorf("%s: %s has no initial frame", name, block.ID)
		return 0
	}
	maxLine := len(strings.Split(block.Code, "\n"))
	words := 0
	for i, frame := range frames {
		if frame.Line < 0 || frame.Line > maxLine || (i > 0 && frame.Line == 0) || strings.TrimSpace(frame.Caption) == "" {
			t.Errorf("%s: %s frame %d points outside its code", name, block.ID, i)
		}
		if len(strings.Fields(frame.Caption)) > 8 {
			t.Errorf("%s: %s frame %d caption exceeds eight words", name, block.ID, i)
		}
		if i > 0 && !strings.HasPrefix(frame.Output, frames[i-1].Output) {
			t.Errorf("%s: %s frame %d output is not cumulative", name, block.ID, i)
		}
		if frame.Panic != "" && i != len(frames)-1 {
			t.Errorf("%s: execution continues after a panic", block.ID)
		}
		words += len(strings.Fields(frame.Caption))
	}
	return words
}

func ch34TextWords(text string) int {
	lines := strings.Split(text, "\n")
	visible := make([]string, 0, len(lines))
	inFence := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			continue
		}
		if !inFence {
			visible = append(visible, line)
		}
	}
	plain := ch34InlineCodePattern.ReplaceAllString(strings.Join(visible, "\n"), "")
	return len(strings.Fields(plain))
}

func ch34OneSentence(text string) bool {
	plain := strings.TrimSpace(ch34InlineCodePattern.ReplaceAllString(text, ""))
	for _, end := range []string{".", "?", "!", "؟"} {
		plain = strings.TrimSuffix(plain, end)
	}
	return plain != "" && !strings.ContainsAny(plain, ".?!؟")
}

func ch34QuestionCode(text string) string {
	var code []string
	inFence := false
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "~~~") {
			if inFence {
				break
			}
			inFence = true
			continue
		}
		if inFence {
			code = append(code, line)
		}
	}
	return strings.Join(code, "\n")
}

func TestCh34StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute chapter 34 examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH34.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	examples := ch34ExamplePattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 8 {
		t.Fatalf("expected eight chapter 34 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { ch34VerifyExample(t, goBin, storyboard, example) })
	}
}

func ch34VerifyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := 0
	for _, levelID := range []string{"34-1-kv-store", "34-2-wal", "34-3-compaction", "34-4-chunk-store"} {
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
	match := ch34StdoutPattern.FindStringSubmatch(transcript)
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
	for _, levelID := range []string{"34-1-kv-store", "34-2-wal", "34-3-compaction", "34-4-chunk-store"} {
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
