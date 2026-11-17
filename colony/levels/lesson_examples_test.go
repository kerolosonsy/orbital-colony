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

func TestSurveyStoryboardExamplesMatchRealGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH6.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e61-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 6 {
		t.Fatalf("expected all six approved survey examples, found %d", len(examples))
	}
	for _, example := range examples {
		t.Run(example[1], func(t *testing.T) {
			verifySurveyExample(t, goBin, storyboard, example)
		})
	}
}

func verifySurveyExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	for _, steps := range [][]LessonStep{surveyLessonArabic, surveyLessonEnglish} {
		for _, step := range steps {
			for _, block := range step.Blocks {
				if block.ExampleID == id && block.Code != code {
					t.Fatal("encoded example differs from approved storyboard")
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
	for _, steps := range [][]LessonStep{surveyLessonArabic, surveyLessonEnglish} {
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
	if outcome == "compile-error" {
		message := regexp.MustCompile(`(?m)^\./main.go:\d+:\d+: (.+)$`)
		want, actual := message.FindStringSubmatch(transcript), message.FindStringSubmatch(string(diagnostics))
		if buildErr == nil || len(want) != 2 || len(actual) != 2 || want[1] != actual[1] {
			t.Fatalf("compile error differs from storyboard: %v\n%s", buildErr, diagnostics)
		}
		return
	}
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
	if outcome == "panic" {
		panicPattern := regexp.MustCompile("(?m)^panic: (.+)$")
		wantPanic := panicPattern.FindString(transcript)
		if runErr == nil || wantPanic == "" || !strings.HasPrefix(stderr.String(), wantPanic+"\n") {
			t.Fatalf("expected %q, got %v\n%s", wantPanic, runErr, stderr)
		}
	} else if runErr != nil || stderr.Len() != 0 {
		t.Fatalf("unexpected failure: %v\n%s", runErr, stderr)
	}
}
