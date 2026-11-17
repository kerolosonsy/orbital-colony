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

func TestLateStoryboardExamplesMatchRealGo(t *testing.T) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain unavailable: cannot execute storyboard examples")
	}
	body, err := os.ReadFile(filepath.Join("..", "..", "STORYBOARD_CH6.md"))
	if err != nil {
		t.Fatal(err)
	}
	storyboard := string(body)
	pattern := regexp.MustCompile("(?s)<!-- go-example id=(e(?:64|65)-[^ ]+) imports=([^ ]+) outcome=([^ ]+) -->\\n```go\\n(.*?)\\n```")
	examples := pattern.FindAllStringSubmatch(storyboard, -1)
	if len(examples) != 2 {
		t.Fatalf("expected two late chapter-6 examples, found %d", len(examples))
	}
	seen := map[string]bool{}
	for _, example := range examples {
		if seen[example[1]] {
			t.Fatalf("duplicate example %s", example[1])
		}
		seen[example[1]] = true
		t.Run(example[1], func(t *testing.T) { verifyLateExample(t, goBin, storyboard, example) })
	}
}

func verifyLateExample(t *testing.T, goBin, storyboard string, example []string) {
	t.Helper()
	id, imports, outcome, code := example[1], example[2], example[3], example[4]
	if outcome != "ok" {
		t.Fatalf("unexpected outcome %q", outcome)
	}
	var found bool
	for _, levelID := range []string{"6-4-stock-order", "6-5-heat-grid"} {
		level := Get(levelID)
		for _, steps := range [][]LessonStep{level.Steps, level.En.Steps} {
			for _, step := range steps {
				for _, block := range step.Blocks {
					if block.ExampleID != id {
						continue
					}
					found = true
					if block.Code != code || strings.Join(block.Imports, ",") != imports || block.Walkthrough == nil || block.Walkthrough.Frames[len(block.Walkthrough.Frames)-1].Output != storyboardExampleOutput(t, storyboard, id) {
						t.Fatalf("%s differs from the Go level or documented output", id)
					}
				}
			}
		}
	}
	if !found {
		t.Fatalf("%s is not encoded in a late level", id)
	}
	wantOutput := storyboardExampleOutput(t, storyboard, id) + "\n"
	dir := t.TempDir()
	quotedImports := []string{}
	for _, name := range strings.Split(imports, ",") {
		quotedImports = append(quotedImports, fmt.Sprintf("%q", name))
	}
	files := map[string]string{
		"go.mod":  "module storyboard/" + id + "\n\ngo 1.24.0\n",
		"main.go": "package main\nimport (" + strings.Join(quotedImports, ";") + ")\nfunc main() {\n" + code + "\n}\n",
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
	diagnostics, buildErr := build.CombinedOutput()
	if buildErr != nil {
		t.Fatalf("build: %v\n%s", buildErr, diagnostics)
	}
	run := exec.CommandContext(ctx, filepath.Join(dir, "example"))
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}
	run.Stdout, run.Stderr = stdout, stderr
	runErr := run.Run()
	if stdout.String() != wantOutput || runErr != nil || stderr.Len() != 0 {
		t.Fatalf("execution mismatch: stdout=%q want=%q error=%v stderr=%q", stdout.String(), wantOutput, runErr, stderr.String())
	}
}

func storyboardExampleOutput(t *testing.T, storyboard, id string) string {
	t.Helper()
	parts := strings.SplitN(storyboard, "#### Verified "+id, 2)
	if len(parts) != 2 {
		t.Fatalf("missing verified output for %s", id)
	}
	transcript := strings.Split(parts[1], "#### Verified")[0]
	match := regexp.MustCompile("(?s)Real stdout:\\n\\n```text\\n(.*?)\\n```").FindStringSubmatch(transcript)
	if match == nil {
		t.Fatalf("missing real stdout for %s", id)
	}
	return match[1]
}
