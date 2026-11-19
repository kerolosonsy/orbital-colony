package main

import (
	"bytes"
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"

	"orbital/colony/levels"
)

// Single run: for ordinary missions the run the level is scored on is also
// the debugger trace. The learner's code is instrumented once, built once
// and executed once, so side effects happen once and the map can follow the
// debugger frame by frame using the same event stream.
//
// Race-detector levels are excluded (the recorder's lock would add
// happens-before edges and could hide the race the level is about), and so
// are labs and service levels, which keep the separate trace request.

const singleRunTraceFile = "learner-trace.jsonl"

func singleRunEligible(lvl *levels.Level, code string) bool {
	if lvl == nil || lvl.Kind != "" || lvl.Race || lvl.Service != "" {
		return false
	}
	// Chapters 10–13 run realtime worlds (ticks follow the wall clock), so
	// two runs of the same code can differ; they keep the separate trace.
	if lvl.Chapter >= 10 && lvl.Chapter <= 13 {
		return false
	}
	file, err := parser.ParseFile(token.NewFileSet(), "main.go", code, 0)
	if err != nil || file.Name.Name != "main" {
		return false
	}
	// Goroutine order varies from run to run and the recorder's lock would
	// shift it further, so concurrent programs keep the separate trace.
	concurrent := false
	ast.Inspect(file, func(n ast.Node) bool {
		if _, ok := n.(*ast.GoStmt); ok {
			concurrent = true
		}
		return !concurrent
	})
	return !concurrent && traceMainPresent(map[string]string{"main.go": code})
}

// singleRunBuild is the instrumented build running next to the plain one.
type singleRunBuild struct {
	dir    string
	done   chan string // instrumented binary path, or "" if it failed
	stop   context.CancelFunc
	active bool
}

func (b *singleRunBuild) wait() string {
	if !b.active {
		return ""
	}
	return <-b.done
}

func (b *singleRunBuild) cancel() {
	if b.active {
		b.stop()
		<-b.done
		b.active = false
	}
}

func (b *singleRunBuild) cleanup() {
	if b.dir != "" {
		os.RemoveAll(b.dir)
	}
}

// startSingleRunBuild instruments code into its own module directory and
// builds it in the background. Any failure just means "no single run".
func (r *runner) startSingleRunBuild(ctx context.Context, lvl *levels.Level, code string, args []string) *singleRunBuild {
	b := &singleRunBuild{}
	if r.plainOnly || !singleRunEligible(lvl, code) {
		return b
	}
	rewritten, recorder, err := rewriteTraceFiles(map[string]string{"main.go": code})
	if err != nil {
		return b
	}
	dir, err := os.MkdirTemp("", "orbital-single-*")
	if err != nil {
		return b
	}
	b.dir = dir
	if err := r.writeTraceProject(dir, lvl, rewritten, recorder); err != nil {
		return b
	}
	bctx, stop := context.WithCancel(ctx)
	b.stop, b.done, b.active = stop, make(chan string, 1), true
	go func() {
		defer stop()
		build := r.goCmd(bctx, dir, append(args, ".")...)
		out := &limitedBuffer{max: 16 << 10}
		build.Stdout, build.Stderr = out, out
		if runCommand(build) != nil {
			b.done <- ""
			return
		}
		b.done <- filepath.Join(dir, "prog")
	}()
	return b
}

// createSingleRunTrace makes the file the recorder appends to.
func createSingleRunTrace(dir string) (string, bool) {
	path := filepath.Join(dir, singleRunTraceFile)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return "", false
	}
	f.Close()
	return path, true
}

// singleRunTrace reads the recorder output of a finished single run and
// links every frame to the number of world events emitted before it.
func singleRunTrace(path string, events []byte, stdout string, completed bool, lang string) *traceResponse {
	f, err := os.Open(path)
	if err != nil {
		return traceUnavailable(lang)
	}
	defer f.Close()
	resp := readTraceFile(f, lang)
	if !resp.Available {
		return resp
	}
	resp.Completed = completed
	resp.Stdout = stdout
	resp.Synced = true
	linkFrameEvents(resp.Frames, events)
	return resp
}

// linkFrameEvents turns each frame's byte offset in the event file into a
// count of complete event lines (the index into the run's event list).
func linkFrameEvents(frames []traceFrame, events []byte) {
	var ends []int64 // byte offset just past each non-empty line
	offset := int64(0)
	for _, line := range bytes.SplitAfter(events, []byte("\n")) {
		offset += int64(len(line))
		if len(bytes.TrimSpace(line)) > 0 && bytes.HasSuffix(line, []byte("\n")) {
			ends = append(ends, offset)
		}
	}
	for i := range frames {
		// With goroutines a frame's offset can be read just before another
		// goroutine's frame is appended, so search each one on its own.
		sim := frames[i].Sim
		frames[i].Events = sort.Search(len(ends), func(k int) bool { return ends[k] > sim })
		frames[i].Sim = 0
	}
}
