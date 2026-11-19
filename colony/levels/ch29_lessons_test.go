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

func TestCh29Lessons(t *testing.T) {
	levelIDs := []string{"29-1-sqlite-crud", "29-2-migrations", "29-3-transfer", "29-4-pool", "29-5-sqlx", "29-6-gorm"}
	prefixes := []string{"c291", "c292", "c293", "c294", "c295", "c296"}
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	globalIDs := map[string]string{}
	exampleIDs := map[string]string{}
	for i, id := range levelIDs {
		level := Get(id)
		if level == nil || level.En == nil {
			t.Fatalf("%s: missing level or English text", id)
		}
		if level.Kind != "lab" {
			t.Errorf("%s: expected lab, got %q", id, level.Kind)
		}
		if level.LessonRevision != 1 || level.En.LessonRevision != 1 {
			t.Errorf("%s: lesson revision must be 1 in both languages", id)
		}
		if len(level.Steps) != 6 || len(level.En.Steps) != 6 {
			t.Errorf("%s: expected six steps in both languages", id)
			continue
		}
		arWords := ch29ValidateLessonCopy(t, id+" Arabic", prefixes[i]+".", level.Steps, level.Hints, arabicTitles)
		enWords := ch29ValidateLessonCopy(t, id+" English", prefixes[i]+".", level.En.Steps, level.En.Hints, englishTitles)
		if arWords < 200 || arWords > 320 || enWords < 200 || enWords > 320 {
			t.Errorf("%s: learner words Arabic=%d English=%d, want 200-320 each", id, arWords, enWords)
		}
		for stepIndex, step := range level.Steps {
			translated := level.En.Steps[stepIndex]
			if step.ID != translated.ID || len(step.Blocks) != len(translated.Blocks) {
				t.Errorf("%s: Arabic and English step structures differ", step.ID)
				continue
			}
			for blockIndex, block := range step.Blocks {
				other := translated.Blocks[blockIndex]
				if block.ID != other.ID || block.Kind != other.Kind {
					t.Errorf("%s: Arabic and English block identities differ", block.ID)
				}
				if block.Sandbox != nil {
					if other.Sandbox == nil || block.Sandbox.Kind != other.Sandbox.Kind || !reflect.DeepEqual(block.Sandbox.Config, other.Sandbox.Config) {
						t.Errorf("%s: sandbox kind or config differs by language", block.ID)
					}
					if !strings.HasPrefix(block.Sandbox.Kind, "ch29-") {
						t.Errorf("%s: sandbox kind %q needs a ch29- prefix", block.ID, block.Sandbox.Kind)
					}
				}
				if block.Question != nil && (other.Question == nil || block.Question.Answer != other.Question.Answer || len(block.Question.Options) != len(other.Question.Options)) {
					t.Errorf("%s: question answer or option count differs", block.ID)
				}
				if block.Walkthrough != nil {
					if other.Walkthrough == nil || block.Code != other.Code || block.ExampleID != other.ExampleID || !reflect.DeepEqual(block.Imports, other.Imports) || len(block.Walkthrough.Frames) != len(other.Walkthrough.Frames) {
						t.Errorf("%s: walkthrough structure differs", block.ID)
						continue
					}
					for frameIndex, frame := range block.Walkthrough.Frames {
						translatedFrame := other.Walkthrough.Frames[frameIndex]
						frame.Caption, translatedFrame.Caption = "", ""
						if !reflect.DeepEqual(frame, translatedFrame) {
							t.Errorf("%s frame %d: execution state differs", block.ID, frameIndex)
						}
					}
					if previous, ok := exampleIDs[block.ExampleID]; ok {
						if previous != id {
							t.Errorf("%s: example ID %q also used by %s", id, block.ExampleID, previous)
						}
					} else {
						exampleIDs[block.ExampleID] = id
					}
				}
			}
		}
		for _, step := range level.Steps {
			if previous, ok := globalIDs[step.ID]; ok {
				t.Errorf("%s: ID %q also used by %s", id, step.ID, previous)
			} else {
				globalIDs[step.ID] = id
			}
			for _, block := range step.Blocks {
				if previous, ok := globalIDs[block.ID]; ok {
					t.Errorf("%s: ID %q also used by %s", id, block.ID, previous)
				} else {
					globalIDs[block.ID] = id
				}
			}
		}
	}
}

func TestCh29SandboxFixtures(t *testing.T) {
	type fixture struct {
		ID     string         `json:"id"`
		Kind   string         `json:"kind"`
		Config map[string]any `json:"config"`
	}
	data, err := os.ReadFile(filepath.Join("..", "..", "web", "sandbox", "ch29.fixtures.json"))
	if err != nil {
		t.Fatal(err)
	}
	var fixtures map[string]fixture
	if err := json.Unmarshal(data, &fixtures); err != nil {
		t.Fatal(err)
	}
	ids := []string{"29-1-sqlite-crud", "29-2-migrations", "29-3-transfer", "29-4-pool", "29-5-sqlx", "29-6-gorm"}
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

func ch29ValidateLessonCopy(t *testing.T, name, prefix string, steps []LessonStep, hints, titles []string) int {
	t.Helper()
	seen := map[string]bool{}
	claim := func(id string) {
		if id == "" || seen[id] {
			t.Errorf("%s: empty or duplicate ID %q", name, id)
		}
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("%s: ID %q needs prefix %q", name, id, prefix)
		}
		seen[id] = true
	}
	words := 0
	sandboxCount := 0
	for stepIndex, step := range steps {
		claim(step.ID)
		if step.Title != titles[stepIndex] || len(step.Blocks) == 0 {
			t.Errorf("%s: step %d has wrong title or no blocks", name, stepIndex)
		}
		words += ch29TextWords(step.Title)
		for _, block := range step.Blocks {
			claim(block.ID)
			if strings.TrimSpace(block.Text) == "" {
				t.Errorf("%s: %s has no learner text", name, block.ID)
			}
			words += ch29TextWords(block.Text)
			switch block.Kind {
			case "markdown":
				if block.Question != nil || block.Walkthrough != nil || block.Sandbox != nil || block.Code != "" {
					t.Errorf("%s: %s markdown has another payload", name, block.ID)
				}
			case "sandbox":
				sandboxCount++
				if block.Sandbox == nil || block.Question != nil || block.Walkthrough != nil || block.Code != "" || !strings.HasPrefix(block.Sandbox.Kind, "ch29-") {
					t.Errorf("%s: %s has an invalid chapter sandbox", name, block.ID)
				}
			case "question":
				if block.Question == nil || len(block.Question.Options) < 2 || len(block.Question.Options) != len(block.Question.Feedback) || block.Question.Answer < 0 || block.Question.Answer >= len(block.Question.Options) {
					t.Errorf("%s: %s has invalid options, feedback, or answer", name, block.ID)
					continue
				}
				words += ch29TextWords(strings.Join(block.Question.Options, " "))
				for _, feedback := range block.Question.Feedback {
					if strings.TrimSpace(feedback) == "" || strings.Count(ch29PlainText(feedback), ".") != 1 {
						t.Errorf("%s: %s feedback must be one sentence: %q", name, block.ID, feedback)
					}
					words += ch29TextWords(feedback)
				}
				for _, option := range block.Question.Options {
					if strings.TrimSpace(option) == "" {
						t.Errorf("%s: %s has an empty option", name, block.ID)
					}
				}
			case "walkthrough":
				if block.Walkthrough == nil || block.Code == "" || block.ExampleID == "" || len(block.Imports) == 0 || block.Question != nil || block.Sandbox != nil {
					t.Errorf("%s: %s has an incomplete walkthrough", name, block.ID)
					continue
				}
				words += ch29ValidateFrames(t, name, block)
			default:
				t.Errorf("%s: %s uses disallowed block kind %q", name, block.ID, block.Kind)
			}
		}
	}
	if sandboxCount != 1 {
		t.Errorf("%s: expected exactly one sandbox, got %d", name, sandboxCount)
	}
	if len(hints) != 3 {
		t.Errorf("%s: expected exactly three hints", name)
	} else {
		for hintIndex, hint := range hints {
			plain := strings.ToLower(strings.TrimSpace(hint))
			for _, label := range []string{"concept:", "structure:", "location:", "المفهوم:", "الهيكل:", "الموضع:", "الفكرة:"} {
				if strings.HasPrefix(plain, label) {
					t.Errorf("%s: hint %d starts with forbidden label %q", name, hintIndex+1, label)
				}
			}
			if strings.Count(ch29PlainText(hint), ".") != 1 {
				t.Errorf("%s: hint %d should be one sentence", name, hintIndex+1)
			}
			words += ch29TextWords(hint)
		}
	}
	if len(steps) > 2 {
		try := steps[2]
		if len(try.Blocks) != 3 || try.Blocks[0].Kind != "question" || try.Blocks[1].Kind != "walkthrough" || try.Blocks[2].Kind != "sandbox" {
			t.Errorf("%s: Try it must predict, reveal, then show one sandbox", name)
		} else if strings.TrimSpace(ch29QuestionCode(try.Blocks[0].Text)) != strings.TrimSpace(try.Blocks[1].Code) {
			t.Errorf("%s: prediction and reveal code differ", name)
		}
		if len(try.Blocks) == 3 && (try.Blocks[2].Sandbox == nil || !strings.HasPrefix(try.Blocks[2].Sandbox.Kind, "ch29-")) {
			t.Errorf("%s: Try it needs one sandbox with a ch29- kind", name)
		}
	}
	t.Logf("%s learner-authored words: %d", name, words)
	return words
}

func ch29ValidateFrames(t *testing.T, name string, block LessonBlock) int {
	t.Helper()
	frames := block.Walkthrough.Frames
	if len(frames) < 2 || frames[0].Line != 0 {
		t.Errorf("%s: %s has no initial frame", name, block.ID)
		return 0
	}
	maxLine := strings.Count(block.Code, "\n") + 1
	words := 0
	for i, frame := range frames {
		if frame.Line < 0 || frame.Line > maxLine || (i > 0 && frame.Line == 0) || strings.TrimSpace(frame.Caption) == "" {
			t.Errorf("%s: %s frame %d points outside code or has no caption", name, block.ID, i)
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

func ch29QuestionCode(text string) string {
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

func ch29PlainText(text string) string {
	var plain strings.Builder
	inCode := false
	for _, r := range text {
		if r == '~' {
			inCode = !inCode
			continue
		}
		if !inCode {
			plain.WriteRune(r)
		}
	}
	return plain.String()
}

func ch29TextWords(text string) int {
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
	return len(strings.Fields(ch29PlainText(strings.Join(visible, "\n"))))
}

func TestCh29GlobalLessonIDs(t *testing.T) {
	levelIDs := []string{"28-1-gin", "28-2-fiber", "28-3-openapi", "28-4-graceful", "28-boss-quarantine", "29-1-sqlite-crud", "29-2-migrations", "29-3-transfer", "29-4-pool", "29-5-sqlx", "29-6-gorm"}
	seen := map[string]bool{}
	for _, levelID := range levelIDs {
		level := Get(levelID)
		if level == nil {
			t.Fatalf("missing level %s", levelID)
		}
		chapterPrefix := fmt.Sprintf("c%d", level.Chapter)
		for _, step := range level.Steps {
			for _, id := range append([]string{step.ID}, ch29BlockIDs(step.Blocks)...) {
				if !strings.HasPrefix(id, chapterPrefix) && !strings.HasPrefix(id, fmt.Sprintf("ch%02d-", level.Chapter)) {
					t.Errorf("%s: ID %q has wrong chapter prefix", levelID, id)
				}
				if seen[id] {
					t.Errorf("duplicate global lesson ID %q", id)
				}
				seen[id] = true
			}
		}
	}
}

func ch29BlockIDs(blocks []LessonBlock) []string {
	ids := make([]string, 0, len(blocks))
	for _, block := range blocks {
		ids = append(ids, block.ID)
	}
	return ids
}

func TestCh29StoryboardExamples(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("Go toolchain unavailable")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH29.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e29[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\n```go\n(.*?)\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 12 {
		t.Fatalf("expected twelve chapter 29 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	cacheDir := filepath.Join(t.TempDir(), "gocache")
	levelIDs := []string{"29-1-sqlite-crud", "29-2-migrations", "29-3-transfer", "29-4-pool", "29-5-sqlx", "29-6-gorm"}
	for _, example := range examples {
		if seen[example[1]] {
			t.Errorf("duplicate storyboard example %q", example[1])
			continue
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) {
			ch29VerifyExample(t, goBin, storyboard, example, levelIDs, cacheDir)
		})
	}
}

func ch29VerifyExample(t *testing.T, goBin, storyboard string, example []string, levelIDs []string, cacheDir string) {
	t.Helper()
	id, importNames, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("%s: unexpected outcome %q", id, outcome)
	}
	found := false
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id {
						found = true
						if block.Code != code {
							t.Fatalf("%s: lesson code differs from storyboard", id)
						}
					}
				}
			}
		}
	}
	if !found {
		t.Fatalf("%s: no lesson walkthrough uses this example", id)
	}
	parts := strings.SplitN(storyboard, "#### Verified "+id+" — ok", 2)
	if len(parts) != 2 {
		t.Fatalf("%s: missing output transcript", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\n\n```text\n(.*?)\n```")
	match := outputPattern.FindStringSubmatch(transcript)
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
	build.Env = ch29SetEnv(os.Environ(), "GOCACHE", cacheDir)
	build.Env = ch29SetEnv(build.Env, "GOWORK", "off")
	build.Env = ch29SetEnv(build.Env, "GOTOOLCHAIN", "local")
	build.Env = ch29SetEnv(build.Env, "GOPROXY", "off")
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
	for _, levelID := range levelIDs {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID == id && block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != wantOutput {
						t.Fatalf("%s: final walkthrough output differs from transcript", id)
					}
				}
			}
		}
	}
}

func ch29SetEnv(env []string, key, value string) []string {
	for i, entry := range env {
		if strings.HasPrefix(entry, key+"=") {
			env[i] = key + "=" + value
			return env
		}
	}
	return append(env, key+"="+value)
}
