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

func TestCh30Lessons(t *testing.T) {
	ids := []string{"30-1-cache-aside", "30-2-documents", "30-3-search-index", "30-4-read-write-split", "30-boss-nursery"}
	prefixes := []string{"c301.", "c302.", "c303.", "c304.", "c30b."}
	arTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	enTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	ch30CheckGlobalIDs(t)
	ch30ValidateChapter(t, ids, prefixes, arTitles, enTitles)
}

func TestCh30SandboxFixtures(t *testing.T) {
	type fixture struct {
		ID     string         `json:"id"`
		Kind   string         `json:"kind"`
		Config map[string]any `json:"config"`
	}
	data, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch30.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]fixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	ids := []string{"30-1-cache-aside", "30-2-documents", "30-3-search-index", "30-4-read-write-split", "30-boss-nursery"}
	if len(fixtures) != len(ids) {
		t.Fatalf("got %d sandbox fixtures, want %d", len(fixtures), len(ids))
	}
	for _, id := range ids {
		level := Get(id)
		block := level.Steps[2].Blocks[2]
		fixture, ok := fixtures[id]
		if !ok {
			t.Errorf("missing fixture for %s", id)
			continue
		}
		if block.ID != fixture.ID || block.Sandbox.Kind != fixture.Kind {
			t.Errorf("%s: fixture identity differs", id)
		}
		want, _ := json.Marshal(fixture.Config)
		for language, steps := range map[string][]LessonStep{"ar": level.Steps, "en": level.En.Steps} {
			got, _ := json.Marshal(steps[2].Blocks[2].Sandbox.Config)
			if string(got) != string(want) {
				t.Errorf("%s %s: Config differs from JavaScript fixture", id, language)
			}
		}
	}
}

func TestCh30StoryboardExamples(t *testing.T) {
	ch30RunStoryboardExamples(t, "STORYBOARD_CH30.md", "e30", []string{
		"30-1-cache-aside", "30-2-documents", "30-3-search-index", "30-4-read-write-split", "30-boss-nursery",
	}, 10)
}

func ch30CheckGlobalIDs(t *testing.T) {
	t.Helper()
	seen := map[string]string{}
	for _, level := range All {
		for _, step := range level.Steps {
			if previous, exists := seen[step.ID]; exists {
				t.Errorf("lesson ID %q is shared by %s and %s", step.ID, previous, level.ID)
			}
			seen[step.ID] = level.ID
			for _, block := range step.Blocks {
				if previous, exists := seen[block.ID]; exists {
					t.Errorf("lesson ID %q is shared by %s and %s", block.ID, previous, level.ID)
				}
				seen[block.ID] = level.ID
			}
		}
	}
}

func ch30ValidateChapter(t *testing.T, ids, prefixes, arTitles, enTitles []string) {
	t.Helper()
	for mission, id := range ids {
		t.Run(id, func(t *testing.T) {
			level := Get(id)
			if level == nil || level.En == nil {
				t.Fatal("missing level or English text")
			}
			if level.Kind != "lab" || level.LessonRevision != 1 || level.En.LessonRevision != 1 {
				t.Fatal("expected a lab with lesson revision 1 in both languages")
			}
			if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
				t.Fatal("expected six structured steps in both languages")
				return
			}
			ch30ValidateCopy(t, id+" Arabic", level.Steps, level.Hints, prefixes[mission], arTitles, 200, 320)
			ch30ValidateCopy(t, id+" English", level.En.Steps, level.En.Hints, prefixes[mission], enTitles, 200, 320)
			for i, ar := range level.Steps {
				en := level.En.Steps[i]
				if ar.ID != en.ID || ar.Title != arTitles[i] || en.Title != enTitles[i] || len(ar.Blocks) != len(en.Blocks) {
					t.Fatalf("%s: Arabic and English step structures differ", ar.ID)
				}
				for j, block := range ar.Blocks {
					translated := en.Blocks[j]
					if block.ID != translated.ID || block.Kind != translated.Kind || block.ExampleID != translated.ExampleID ||
						!reflect.DeepEqual(block.Imports, translated.Imports) || block.Code != translated.Code {
						t.Fatalf("%s: language payloads differ", block.ID)
					}
					if block.Sandbox != nil {
						if translated.Sandbox == nil || block.Sandbox.Kind != translated.Sandbox.Kind || !reflect.DeepEqual(block.Sandbox.Config, translated.Sandbox.Config) {
							t.Fatalf("%s: sandbox kind or config differs by language", block.ID)
						}
						if !strings.HasPrefix(block.Sandbox.Kind, "ch30-") {
							t.Fatalf("%s: sandbox kind %q needs a ch30- prefix", block.ID, block.Sandbox.Kind)
						}
					}
					if block.Question != nil {
						if translated.Question == nil || block.Question.Answer != translated.Question.Answer ||
							len(block.Question.Options) != len(translated.Question.Options) {
							t.Fatalf("%s: question answers or option counts differ", block.ID)
						}
					}
					if block.Walkthrough != nil {
						if translated.Walkthrough == nil || len(block.Walkthrough.Frames) != len(translated.Walkthrough.Frames) {
							t.Fatalf("%s: walkthrough frames differ", block.ID)
						}
						for k, frame := range block.Walkthrough.Frames {
							other := translated.Walkthrough.Frames[k]
							frame.Caption, other.Caption = "", ""
							if !reflect.DeepEqual(frame, other) {
								t.Fatalf("%s frame %d: execution state differs", block.ID, k)
							}
						}
					}
				}
			}
		})
	}
}

func ch30ValidateCopy(t *testing.T, name string, steps []LessonStep, hints []string, prefix string, titles []string, minWords, maxWords int) {
	t.Helper()
	requireSandbox := strings.HasPrefix(prefix, "c30")
	seen := map[string]bool{}
	words := 0
	sandboxCount := 0
	claim := func(id string) {
		if id == "" || seen[id] || !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: empty, duplicate, or incorrectly prefixed ID %q", name, id)
		}
		seen[id] = true
	}
	for i, step := range steps {
		claim(step.ID)
		if step.Title != titles[i] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has the wrong title or no blocks", name, i+1)
		}
		words += ch30VisibleWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch30VisibleWords(block.Text) + ch30VisibleWords(block.Note)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown contains another payload", name, block.ID)
				}
			case "sandbox":
				sandboxCount++
				if block.Sandbox == nil || block.Question != nil || block.Walkthrough != nil || block.Code != "" || (requireSandbox && !strings.HasPrefix(block.Sandbox.Kind, "ch30-")) {
					t.Errorf("%s: %s has an invalid chapter sandbox", name, block.ID)
				}
			case "question":
				if block.Question == nil || block.Walkthrough != nil || block.Sandbox != nil ||
					len(block.Question.Options) != 3 || len(block.Question.Options) != len(block.Question.Feedback) ||
					block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid question data", name, block.ID)
					continue
				}
				words += ch30VisibleWords(strings.Join(block.Question.Options, " "))
				for _, feedback := range block.Question.Feedback {
					if !ch30OneSentence(feedback) {
						t.Errorf("%s: feedback must be one sentence: %q", block.ID, feedback)
					}
					words += ch30VisibleWords(feedback)
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Question != nil || block.Sandbox != nil || block.Code == "" ||
					block.ExampleID == "" || len(block.Imports) == 0 {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch30ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed kind %q", name, block.ID, block.Kind)
			}
		}
	}
	expectedSandboxCount := 0
	if requireSandbox {
		expectedSandboxCount = 1
	}
	if sandboxCount != expectedSandboxCount {
		t.Errorf("%s: expected %d sandbox(es), got %d", name, expectedSandboxCount, sandboxCount)
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected three hints, got %d", name, len(hints))
	}
	for i, hint := range hints {
		trimmed := strings.TrimSpace(hint)
		for _, label := range []string{"Concept:", "Structure:", "Location:", "المفهوم:", "الفكرة:", "الهيكل:", "الموضع:"} {
			if strings.HasPrefix(trimmed, label) {
				t.Errorf("%s: hint %d starts with a label: %q", name, i+1, hint)
			}
		}
		if !ch30OneSentence(hint) {
			t.Errorf("%s: hint %d must be one sentence: %q", name, i+1, hint)
		}
		words += ch30VisibleWords(hint)
	}
	if len(steps) > 2 {
		try := steps[2]
		wantBlocks := 2
		if requireSandbox {
			wantBlocks = 3
		}
		if len(try.Blocks) != wantBlocks || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || (requireSandbox && try.Blocks[2].Kind != "sandbox") {
			t.Errorf("%s: Try it has the wrong prediction, reveal, and sandbox structure", name)
		} else if strings.TrimSpace(ch30QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction code differs from its walkthrough", name)
		}
		if requireSandbox && len(try.Blocks) == 3 && (try.Blocks[2].Sandbox == nil || !strings.HasPrefix(try.Blocks[2].Sandbox.Kind, "ch30-")) {
			t.Errorf("%s: Try it needs one sandbox with a ch30- kind", name)
		}
	}
	if words < minWords || words > maxWords {
		t.Errorf("%s: learner text has %d words; expected %d–%d", name, words, minWords, maxWords)
	}
	t.Logf("%s learner text: %d words", name, words)
}

func ch30ValidateFrames(t *testing.T, name string, block LessonBlock) int {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Errorf("%s: %s has no before-execution frame", name, block.ID)
		return 0
	}
	maxLine := strings.Count(block.Code, "\n") + 1
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
		words += ch30VisibleWords(frame.Caption)
	}
	return words
}

func ch30VisibleWords(text string) int {
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
	inlineCode := regexp.MustCompile("~[^~]*~")
	return len(strings.Fields(inlineCode.ReplaceAllString(strings.Join(visible, "\n"), "")))
}

func ch30OneSentence(sentence string) bool {
	plain := regexp.MustCompile("~[^~]*~").ReplaceAllString(sentence, "")
	plain = strings.TrimSpace(strings.TrimRight(plain, ".!?؟"))
	return plain != "" && !strings.ContainsAny(plain, ".!?؟")
}

func ch30QuestionCode(text string) string {
	lines := strings.Split(text, "\n")
	inFence := false
	var code []string
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
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

func ch30RunStoryboardExamples(t *testing.T, filename, examplePrefix string, levelIDs []string, expectedCount int) {
	t.Helper()
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("Go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", filename))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	fence := strings.Repeat(string(rune(96)), 3)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(" + regexp.QuoteMeta(examplePrefix) + "[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n" + fence + "go\n(.*?)\n" + fence)
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != expectedCount {
		t.Fatalf("expected %d storyboard examples, found %d", expectedCount, len(examples))
	}
	blocks := map[string]LessonBlock{}
	for _, levelID := range levelIDs {
		level := Get(levelID)
		if level == nil || level.En == nil {
			t.Fatalf("missing level %s", levelID)
		}
		for _, step := range level.Steps {
			for _, block := range step.Blocks {
				if block.ExampleID != "" {
					if _, exists := blocks[block.ExampleID]; exists {
						t.Fatalf("duplicate lesson example %q", block.ExampleID)
					}
					blocks[block.ExampleID] = block
				}
			}
		}
	}
	seen := map[string]bool{}
	for _, example := range examples {
		id := example[1]
		if seen[id] || example[3] != "ok" {
			t.Fatalf("duplicate example or unexpected outcome for %q", id)
		}
		seen[id] = true
		block, exists := blocks[id]
		if !exists || block.Walkthrough == nil {
			t.Fatalf("storyboard example %q is not used by a lesson", id)
		}
		if block.Code != example[4] || strings.Join(block.Imports, ",") != example[2] {
			t.Fatalf("%s differs from its structured walkthrough", id)
		}
		wantOutput := ch30StoryboardOutput(t, storyboard, fence, id)
		frames := block.Walkthrough.Frames
		if len(frames) == 0 || frames[len(frames)-1].Output+"\n" != wantOutput {
			t.Fatalf("%s walkthrough output differs from its transcript", id)
		}
		t.Run(id, func(t *testing.T) {
			ch30RunOneExample(t, goBin, id, example[2], example[4], wantOutput)
		})
	}
	if len(seen) != len(blocks) {
		t.Fatalf("storyboard has %d examples but lessons use %d", len(seen), len(blocks))
	}
}

func ch30StoryboardOutput(t *testing.T, storyboard, fence, id string) string {
	t.Helper()
	verified := strings.SplitN(storyboard, "#### Verified "+id, 2)
	if len(verified) != 2 {
		t.Fatalf("missing verified output for %s", id)
	}
	after := verified[1]
	marker := "Real stdout:\n\n" + fence + "text\n"
	start := strings.Index(after, marker)
	if start < 0 {
		t.Fatalf("missing real stdout for %s", id)
	}
	after = after[start+len(marker):]
	end := strings.Index(after, "\n"+fence)
	if end < 0 {
		t.Fatalf("unterminated stdout for %s", id)
	}
	return after[:end] + "\n"
}

func ch30RunOneExample(t *testing.T, goBin, id, importList, code string, wantOutput string) {
	t.Helper()
	dir := t.TempDir()
	imports := make([]string, 0, len(strings.Split(importList, ",")))
	for _, name := range strings.Split(importList, ",") {
		imports = append(imports, fmt.Sprintf("%q", name))
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
		t.Fatalf("build failed: %v\n%s", err, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	var stdout, stderr strings.Builder
	run.Stdout, run.Stderr = &stdout, &stderr
	runErr := run.Run()
	if stdout.String() != wantOutput || runErr != nil || stderr.Len() != 0 {
		t.Fatalf("run mismatch: stdout=%q want=%q err=%v stderr=%q", stdout.String(), wantOutput, runErr, stderr.String())
	}
}
