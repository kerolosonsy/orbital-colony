package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"orbital/colony/levels"
)

// The single run must score exactly like the plain program: same result,
// ticks, stars, stdout and world events for the same seed.
func TestSingleRunMatchesPlainRun(t *testing.T) {
	plain := testRunner(t)
	plain.plainOnly = true
	traced := testRunner(t)
	checked := 0
	for _, l := range levels.All {
		code, err := os.ReadFile(filepath.Join("testdata", "solutions", l.ID+".go"))
		if err != nil || !singleRunEligible(l, string(code)) {
			continue
		}
		if l.Chapter%3 != 1 && checked >= 6 { // spread across chapters, keep it quick
			continue
		}
		checked++
		t.Run(l.ID, func(t *testing.T) {
			t.Parallel()
			ctx := context.WithValue(context.Background(), reviewSeedKey{}, int64(424242))
			want, err := plain.Run(ctx, l, string(code))
			if err != nil {
				t.Fatal(err)
			}
			got, err := traced.Run(ctx, l, string(code))
			if err != nil {
				t.Fatal(err)
			}
			if got.Trace == nil || !got.Trace.Available || len(got.Trace.Frames) == 0 {
				t.Fatalf("no single-run trace: %+v", got.Trace)
			}
			if want.Trace != nil {
				t.Fatal("plain run must not carry a trace")
			}
			if got.Result != want.Result || got.Stdout != want.Stdout || got.Phase != want.Phase {
				t.Fatalf("instrumented run differs\nwant %+v %q\ngot  %+v %q", want.Result, want.Stdout, got.Result, got.Stdout)
			}
			if !reflect.DeepEqual(eventsOf(t, want.Events), eventsOf(t, got.Events)) {
				t.Fatalf("world events differ: %d vs %d events", len(want.Events), len(got.Events))
			}
			last := got.Trace.Frames[len(got.Trace.Frames)-1]
			if last.Events < 1 || last.Events > len(got.Events) {
				t.Fatalf("last frame events = %d, run has %d", last.Events, len(got.Events))
			}
			for i := 1; i < len(got.Trace.Frames); i++ {
				if got.Trace.Frames[i].Events < got.Trace.Frames[i-1].Events && got.Trace.Frames[i].Goroutine == got.Trace.Frames[i-1].Goroutine {
					t.Fatalf("frame %d goes back in events", i)
				}
			}
		})
	}
	if checked < 15 {
		t.Fatalf("only %d eligible levels checked", checked)
	}
}

func eventsOf(t *testing.T, raw []json.RawMessage) []map[string]any {
	t.Helper()
	out := make([]map[string]any, 0, len(raw))
	for _, r := range raw {
		var m map[string]any
		if err := json.Unmarshal(r, &m); err != nil {
			t.Fatal(err)
		}
		out = append(out, m)
	}
	return out
}

// Stepping onto the second Move in 1-1 must show the drone after exactly one
// move: frames carry how many events happened before each statement.
func TestSingleRunFramesFollowEvents(t *testing.T) {
	r := testRunner(t)
	code := "package main\n\nimport \"orbital/colony\"\n\nfunc main() {\n\tcolony.Move(colony.East)\n\tcolony.Move(colony.East)\n\tcolony.Move(colony.East)\n}\n"
	ctx := context.WithValue(context.Background(), reviewSeedKey{}, int64(7))
	resp, err := r.Run(ctx, levels.Get("1-1-first-signal"), code)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Trace == nil || !resp.Trace.Synced {
		t.Fatalf("expected synced trace, got %+v", resp.Trace)
	}
	moves := map[int]int{} // line -> events before it
	for _, f := range resp.Trace.Frames {
		if f.Event == "line" {
			moves[f.Line] = f.Events
		}
	}
	countMoves := func(n int) int {
		c := 0
		for _, e := range eventsOf(t, resp.Events)[:n] {
			if e["op"] == "move" {
				c++
			}
		}
		return c
	}
	for line, wantMoves := range map[int]int{6: 0, 7: 1, 8: 2} {
		if got := countMoves(moves[line]); got != wantMoves {
			t.Fatalf("line %d: %d moves before it, want %d", line, got, wantMoves)
		}
	}
}

// Code that does not compile still reports the learner's own errors.
func TestSingleRunCompileErrorFallsBack(t *testing.T) {
	r := testRunner(t)
	resp, err := r.Run(context.Background(), levels.Get("1-1-first-signal"), "package main\n\nimport \"orbital/colony\"\n\nfunc main() {\n\tx := 1\n\tcolony.Move(colony.East)\n}\n")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Phase != "compile" || len(resp.CompileErrors) != 1 || resp.CompileErrors[0].Line != 6 || resp.CompileErrors[0].File != "main.go" {
		t.Fatalf("got %+v\n%s", resp.CompileErrors, resp.CompileOutput)
	}
	if resp.Trace != nil {
		t.Fatal("no trace for code that did not build")
	}
}

func TestSingleRunSkipsRaceLevels(t *testing.T) {
	for _, l := range levels.All {
		if l.Race && singleRunEligible(l, "package main\nfunc main(){}\n") {
			t.Fatalf("%s: race level must not use the single run", l.ID)
		}
	}
}

// A String method is called by the recorder while formatting; that must not
// recurse, and the level must still be scored by the plain run.
func TestSingleRunStringerDoesNotRecurse(t *testing.T) {
	r := testRunner(t)
	code, err := os.ReadFile(filepath.Join("testdata", "solutions", "9-1-stringer.go"))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := r.Run(context.Background(), levels.Get("9-1-stringer"), string(code))
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Result.Success {
		t.Fatalf("level not solved: %+v", resp.Result)
	}
	if resp.Trace == nil || !resp.Trace.Synced || len(resp.Trace.Frames) == 0 {
		t.Fatalf("expected a synced trace, got %+v", resp.Trace)
	}
}
