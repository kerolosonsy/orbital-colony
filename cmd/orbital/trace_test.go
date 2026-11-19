package main

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"orbital/colony/levels"
)

func TestTraceRewriteCapturesLoopsClosuresGoroutinesAndShadowing(t *testing.T) {
	code := `package main
import "fmt"
func main() {
	left, right := 1, 2
	left, right = right, left
	sum := 0
	for index := 0; index < 2; index++ {
		if index == 0 {
			sum += left
		} else {
			sum += right
		}
	}
	shadow := 10
	{
		shadow := 20
		fmt.Println(shadow)
	}
	closure := func(delta int) { sum += delta; fmt.Println(sum) }
	closure(1)
	for item := range []int{5} { fmt.Println(item) }
	fmt.Println("range done")
	done := make(chan struct{})
	go func() { workerValue := 7; fmt.Println(workerValue); close(done) }()
	<-done
	fmt.Println(left, right, shadow, sum)
}`
	resp := runTraceFixture(t, code)
	if !resp.Available || !resp.Completed || len(resp.Frames) == 0 {
		t.Fatalf("trace response: %+v", resp)
	}
	inner := frameAt(t, resp.Frames, sourceLine(code, "fmt.Println(shadow)"))
	if variable(inner, "shadow") != "20" {
		t.Fatalf("inner shadow frame: %+v", inner)
	}
	outer := frameAt(t, resp.Frames, sourceLine(code, "fmt.Println(left, right, shadow, sum)"))
	if variable(outer, "left") != "2" || variable(outer, "right") != "1" || variable(outer, "shadow") != "10" {
		t.Fatalf("multiple assignment or outer shadow frame: %+v", outer)
	}
	closure := frameWithVariable(t, resp.Frames, sourceLine(code, "sum += delta"), "delta")
	if variable(closure, "delta") != "1" || variable(closure, "sum") != "3" {
		t.Fatalf("closure frame: %+v", closure)
	}
	worker := frameWithVariable(t, resp.Frames, sourceLine(code, "fmt.Println(workerValue)"), "workerValue")
	if worker.Goroutine <= 1 || variable(worker, "workerValue") != "7" {
		t.Fatalf("goroutine frame: %+v", worker)
	}
	if got := countFramesAt(resp.Frames, sourceLine(code, "sum += right")); got != 1 {
		t.Fatalf("else branch frame count = %d, want 1", got)
	}
	if got := countFramesAt(resp.Frames, sourceLine(code, "for index :=")); got != 5 {
		t.Fatalf("loop header/body/post frame count = %d, want 5", got)
	}
	if afterRange := frameAt(t, resp.Frames, sourceLine(code, `fmt.Println("range done")`)); variable(afterRange, "item") != "" {
		t.Fatalf("range variable escaped its scope: %+v", afterRange)
	}
}

func TestTraceCapsFramesAndLetsProgramFinish(t *testing.T) {
	code := `package main
import "fmt"
func main() {
	for i := 0; i < 600; i++ {}
	fmt.Println("finished")
}`
	resp := runTraceFixture(t, code)
	if !resp.Available || !resp.Completed || !resp.Truncated {
		t.Fatalf("trace did not finish with a truncated notice: %+v", resp)
	}
	if len(resp.Frames) > traceFrameLimit {
		t.Fatalf("got %d frames, cap is %d", len(resp.Frames), traceFrameLimit)
	}
}

func TestTraceCapsValuesAndTotalBytes(t *testing.T) {
	valueCode := `package main
import (
	"fmt"
	"strings"
)
func main() {
	huge := strings.Repeat("x", 1000)
	fmt.Println(huge)
}`
	valueTrace := runTraceFixture(t, valueCode)
	frame := frameWithVariable(t, valueTrace.Frames, sourceLine(valueCode, "fmt.Println(huge)"), "huge")
	if len(variable(frame, "huge")) > traceValueLimit || !frame.Truncated || !valueTrace.Truncated {
		t.Fatalf("value limit was not recorded: %+v", frame)
	}

	var byteCode strings.Builder
	byteCode.WriteString("package main\nimport \"strings\"\nfunc main() {\n")
	for i := 0; i < traceVarLimit; i++ {
		byteCode.WriteString("value")
		byteCode.WriteString(strconv.Itoa(i))
		byteCode.WriteString(" := strings.Repeat(\"x\", 160)\n")
	}
	byteCode.WriteString("for i := 0; i < 30; i++ {}\n}\n")
	byteTrace := runTraceFixture(t, byteCode.String())
	if !byteTrace.Available || !byteTrace.Completed || !byteTrace.Truncated || len(byteTrace.Frames) >= traceFrameLimit {
		t.Fatalf("total-byte cap did not stop recording and allow completion: frames=%d response=%+v", len(byteTrace.Frames), byteTrace)
	}
}

func TestTraceBuildFailureReturnsFriendlyUnavailable(t *testing.T) {
	r := testRunner(t)
	dir := t.TempDir()
	if err := r.writeModule(dir, &levels.Level{ID: "trace-build-failure"}); err != nil {
		t.Fatal(err)
	}
	if err := writeProjectFiles(dir, map[string]string{
		"main.go": `package main
func main() { println("normal source") }`,
	}); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dir+"/0_trace.go", []byte("package main\nfunc generatedTrace() { missingTraceName() }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	resp, err := r.buildAndRunTrace(context.Background(), dir, &levels.Level{ID: "trace-build-failure"}, "en")
	if err != nil || resp.Available || !strings.Contains(resp.Message, "Trace is unavailable") {
		t.Fatalf("compile failure response = %+v, err %v", resp, err)
	}
}

func TestTraceKeepsNormalMissionOutputAndTicksUnchanged(t *testing.T) {
	r := testRunner(t)
	lvl := levels.Get("1-1-first-signal")
	code := `package main
import (
	"fmt"
	"orbital/colony"
)
func main() {
	count := 1
	fmt.Println("unchanged", count)
	colony.Move(colony.East)
}`
	before, err := r.Run(context.Background(), lvl, code)
	if err != nil {
		t.Fatal(err)
	}
	traced, err := r.Trace(context.Background(), lvl, code, nil)
	if err != nil || !traced.Available || len(traced.Frames) == 0 {
		t.Fatalf("trace response = %+v, err %v", traced, err)
	}
	after, err := r.Run(context.Background(), lvl, code)
	if err != nil {
		t.Fatal(err)
	}
	if before.Stdout != after.Stdout || before.Stderr != after.Stderr || before.Result.Ticks != after.Result.Ticks || !reflect.DeepEqual(before.Events, after.Events) || !reflect.DeepEqual(before.Result, after.Result) {
		t.Fatalf("trace changed normal run: before=%+v after=%+v", before, after)
	}
}

func TestTraceMainIsRequiredForLabs(t *testing.T) {
	if traceMainPresent(map[string]string{"main.go": "package main\nfunc helper() {}"}) {
		t.Fatal("library-only package was marked traceable")
	}
	if !traceMainPresent(map[string]string{"main.go": "package main\nfunc main() {}"}) {
		t.Fatal("package with main was not marked traceable")
	}
	if traceLevelSupported(&levels.Level{Kind: "colony"}) {
		t.Fatal("non-mission level was marked traceable")
	}
}

func TestTraceRecorderFilenameIsValidAndFirst(t *testing.T) {
	if got := traceRecorderName(map[string]string{"main.go": "package main"}); got != "0.go" {
		t.Fatalf("recorder filename = %q, want earliest Go filename", got)
	}
	if got := traceRecorderName(map[string]string{"0.go": "package main"}); got != "" {
		t.Fatalf("recorder filename = %q when no earlier safe filename exists", got)
	}
}

func TestTraceRewrittenCompileFailureReturnsFriendlyUnavailable(t *testing.T) {
	code := `package main
func main() {
	var count int = "not an int"
	println(count)
}`
	resp := runTraceFixture(t, code)
	if resp.Available || resp.Message == "" {
		t.Fatalf("compile failure response = %+v", resp)
	}
}

func runTraceFixture(t *testing.T, code string) *traceResponse {
	t.Helper()
	resp, err := testRunner(t).Trace(context.Background(), &levels.Level{ID: "trace-fixture"}, code, nil)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func sourceLine(source, snippet string) int {
	return strings.Count(source[:strings.Index(source, snippet)], "\n") + 1
}

func frameAt(t *testing.T, frames []traceFrame, line int) traceFrame {
	t.Helper()
	for _, frame := range frames {
		if frame.Line == line {
			return frame
		}
	}
	t.Fatalf("no frame at line %d in %+v", line, frames)
	return traceFrame{}
}

func frameWithVariable(t *testing.T, frames []traceFrame, line int, name string) traceFrame {
	t.Helper()
	for _, frame := range frames {
		if frame.Line == line && variable(frame, name) != "" {
			return frame
		}
	}
	t.Fatalf("no frame at line %d with variable %q", line, name)
	return traceFrame{}
}

func variable(frame traceFrame, name string) string {
	for _, variable := range frame.Variables {
		if variable.Name == name {
			return variable.Formatted
		}
	}
	return ""
}

func countFramesAt(frames []traceFrame, line int) int {
	count := 0
	for _, frame := range frames {
		if frame.Line == line {
			count++
		}
	}
	return count
}
