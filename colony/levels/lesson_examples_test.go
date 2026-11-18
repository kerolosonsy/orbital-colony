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

func TestChapterSixStoryboardExamplesMatchRealGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH6.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e61-(?:crystal-log-v2|sequence-v2)|e62-range-clone-sort|e63-map-count) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	wantIDs := map[string]bool{
		"e61-crystal-log-v2": true, "e61-sequence-v2": true,
		"e62-range-clone-sort": true, "e63-map-count": true,
	}
	if len(examples) != len(wantIDs) {
		t.Fatalf("expected four chapter-six lesson examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] || !wantIDs[example[1]] {
			t.Fatalf("unexpected or duplicate example ID %q", example[1])
		}
		seen[example[1]] = true
		delete(wantIDs, example[1])
		t.Run(example[1], func(t *testing.T) {
			verifyChapterSixExample(t, goBin, storyboard, example)
		})
	}
	if len(wantIDs) != 0 {
		t.Fatalf("missing examples: %v", wantIDs)
	}
}

func verifyChapterSixExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatal("chapter-six lesson example must run successfully")
	}
	for _, steps := range [][]LessonStep{surveyLessonArabic, surveyLessonEnglish, topThreeLessonArabic, topThreeLessonEnglish, censusLessonArabic, censusLessonEnglish} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == id && block.Code != code {
					t.Fatal("encoded example differs from approved storyboard")
				}
				if id == "e61-sequence-v2" && block.Sandbox != nil && block.Sandbox.Kind == "slice-append-sequence" {
					literal := strings.ReplaceAll(fmt.Sprintf("%#v", block.Sandbox.Sequence), " ", "")
					if !strings.Contains(strings.ReplaceAll(code, " ", ""), literal) {
						t.Fatal("sequence fixture differs from playground data")
					}
				}
				if id == "e63-map-count" && block.Sandbox != nil && block.Sandbox.Kind == "map-count" {
					for _, kind := range block.Sandbox.Kinds {
						if !strings.Contains(code, fmt.Sprintf("%q", kind)) {
							t.Fatalf("map fixture omits playground kind %q", kind)
						}
					}
				}
			}
		}
	}
	parts := strings.SplitN(storyboard, "#### Verified `"+id+"`", 2)
	if len(parts) != 2 {
		t.Fatal("missing documented example output")
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	outputPattern := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```")
	wantOutput := ""
	if match := outputPattern.FindStringSubmatch(transcript); match != nil {
		wantOutput = match[1] + "\n"
	}
	for _, steps := range [][]LessonStep{surveyLessonArabic, surveyLessonEnglish, topThreeLessonArabic, topThreeLessonEnglish, censusLessonArabic, censusLessonEnglish} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == id && block.Walkthrough != nil {
					frames := block.Walkthrough.Frames
					if frames[len(frames)-1].Output != strings.TrimSuffix(wantOutput, "\n") {
						t.Fatal("walkthrough terminal output differs from documented execution")
					}
				}
			}
		}
	}
	dir := t.TempDir()
	quotedImports := []string{}
	for _, name := range strings.Split(imports, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", name))
	}
	for name, contents := range map[string]string{"go.mod": "module storyboard/" + id + "\n\ngo 1.24.0\n", "main.go": "package main\nimport (" + strings.Join(quotedImports, ";") + ")\nfunc main() {\n" + code + "\n}\n"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	build := exec.CommandContext(ctx, goBin, "build", "-trimpath", "-gcflags=-lang=go1.24", "-o", "example", ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	diagnostics, buildErr := build.CombinedOutput()
	if buildErr != nil {
		t.Fatalf("build: %v\n%s", buildErr, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	runErr := run.Run()
	if stdout.String() != wantOutput {
		t.Fatalf("stdout = %q, documented %q", stdout.String(), wantOutput)
	}
	if runErr != nil || stderr.Len() != 0 {
		t.Fatalf("unexpected failure: %v\n%s", runErr, stderr)
	}
}
