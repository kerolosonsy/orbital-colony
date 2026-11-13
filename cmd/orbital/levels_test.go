package main

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"orbital/colony/levels"
)

func testRunner(t *testing.T) *runner {
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not on PATH")
	}
	abs, _ := filepath.Abs("../../colony")
	return &runner{goBin: goBin, colonyDir: abs}
}

// Every level must be solvable by its reference solution, and every
// starter program must compile.
func TestLevelsSolvable(t *testing.T) {
	r := testRunner(t)
	for _, l := range levels.All {
		if l.Kind == "lab" {
			t.Run(l.ID, func(t *testing.T) {
				if l.Service != "" && os.Getenv("ORBITAL_SERVICES") != "1" {
					t.Skip("set ORBITAL_SERVICES=1 for service labs")
				}
				testLab(t, r, l)
			})
			continue
		}
		t.Run(l.ID, func(t *testing.T) {
			if l.Service != "" && os.Getenv("ORBITAL_SERVICES") != "1" {
				t.Skip("set ORBITAL_SERVICES=1 for service levels")
			}
			code, err := os.ReadFile(filepath.Join("testdata", "solutions", l.ID+".go"))
			if err != nil {
				t.Fatalf("no reference solution: %v", err)
			}
			// Exercise concurrent drone references even when the lesson itself
			// does not enable the detector. Keep the starter's original mode.
			reference := *l
			if l.Chapter >= 10 && l.Chapter <= 13 {
				reference.Race = true
			}
			for i := 0; i < 5; i++ { // levels are randomized; try several seeds
				resp, err := r.Run(context.Background(), &reference, string(code))
				if err != nil {
					t.Fatal(err)
				}
				if !resp.Result.Success || strings.Contains(resp.Stderr, "WARNING: DATA RACE") {
					t.Fatalf("phase=%s msg=%q missing=%v compile=%s stderr=%s", resp.Phase, resp.Result.Message, resp.Missing, resp.CompileOutput, resp.Stderr)
				}
				t.Logf("ticks=%d stars=%d (par %d, gold %d)", resp.Result.Ticks, resp.Result.Stars, resp.Result.Par, resp.Result.Gold)
			}

			resp, err := r.Run(context.Background(), l, l.Starter)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Phase == "compile" {
				t.Fatalf("starter does not compile: %s", resp.CompileOutput)
			}
		})
	}
}

func TestCompileErrorTips(t *testing.T) {
	r := testRunner(t)
	resp, err := r.Run(context.Background(), levels.All[0], "package main\n\nimport \"orbital/colony\"\n\nfunc main() {\n\tx := 1\n\tcolony.Move(colony.East)\n}\n")
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.CompileErrors) != 1 || resp.CompileErrors[0].Line != 6 || resp.CompileErrors[0].Tip == "" {
		t.Fatalf("got %+v\n%s", resp.CompileErrors, resp.CompileOutput)
	}
}

// Missions with a fail condition must actually fail when it is triggered.
func TestFailConditions(t *testing.T) {
	r := testRunner(t)
	cases := map[string]struct{ code, want string }{
		"7-2-nav-system": {`package main

import "orbital/colony"

type sweeper struct{}

func (sweeper) row(w int) {
	for x := 0; x < w; x++ {
		colony.Harvest()
		colony.Move(colony.East)
	}
}

func main() {
	w, h := colony.WorldSize()
	sweeper{}.row(0)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			colony.Harvest()
			colony.Move(colony.East)
		}
		for x := 0; x < w; x++ {
			colony.Move(colony.West)
		}
		colony.Move(colony.North)
	}
}
`, "ملوّث"},
		"6-2-top-three": {`package main

import "orbital/colony"

func main() {
	for i := 0; i < 3; i++ {
		colony.Harvest()
		colony.Move(colony.East)
	}
	for range []int{1} {
	}
}
`, "المهمة فشلت"},
	}
	cases["10-1-second-drone"] = struct{ code, want string }{`package main

import (
	"orbital/colony"
	"sync"
)

func main() {
	helper := colony.Spawn()
	var wg sync.WaitGroup
	go func() {
		for i := 0; i < 30; i++ {
			helper.Move(colony.East)
		}
	}()
	colony.Harvest()
	_ = &wg
}
`, "main خلصت"}
	cases["11-1-scout"] = struct{ code, want string }{`package main

import "orbital/colony"

func main() {
	ch := make(chan int)
	go func() { ch <- 1 }()
	<-ch
	for x := 0; x < 12; x++ {
		colony.Harvest()
		colony.Move(colony.East)
	}
}
`, "الصغير"}
	cases["11-3-storm"] = struct{ code, want string }{"", "العاصفة"}
	for _, id := range []string{"12-1-shared-total", "12-2-registry", "12-3-deadlock", "13-1-recall", "13-2-timeouts", "13-3-first-wins"} {
		cases[id] = struct{ code, want string }{"", ""} // buggy starters must not pass
	}
	cases["12-1-shared-total"] = struct{ code, want string }{"", "race"}
	cases["12-3-deadlock"] = struct{ code, want string }{"", "deadlock"}
	cases["13-2-timeouts"] = struct{ code, want string }{"", "الوقت خلص"}
	cases["8-1-copy-bug"] = struct{ code, want string }{"", ""} // starter as shipped must not pass
	cases["8-2-cargo-runs"] = struct{ code, want string }{"", ""}
	for id, c := range cases {
		if c.code == "" {
			c.code = levels.Get(id).Starter
		}
		t.Run(id, func(t *testing.T) {
			resp, err := r.Run(context.Background(), levels.Get(id), c.code)
			// Worlds are random: a rare layout lets a wrong program finish
			// before it reaches its mistake, so a success gets two more tries.
			for try := 0; err == nil && resp.Result.Success && try < 2; try++ {
				resp, err = r.Run(context.Background(), levels.Get(id), c.code)
			}
			if err != nil {
				t.Fatal(err)
			}
			if resp.Result.Success || !strings.Contains(resp.Result.Message+resp.Result.Progress, c.want) {
				t.Fatalf("success=%v msg=%q stderr=%s compile=%s", resp.Result.Success, resp.Result.Message, resp.Stderr, resp.CompileOutput)
			}
		})
	}
}

// Building previews must never start mission timers: they would call
// os.Exit inside the server process.
func TestPreviewHasNoTimers(t *testing.T) {
	for _, l := range levels.All {
		var buf strings.Builder
		levels.Start(l.ID, &buf, 1)
	}
	// Outlast the longest realtime budget (500 ticks × 4ms); a running
	// watchdog would kill this test binary.
	time.Sleep(2200 * time.Millisecond)
}

// testLab runs a lab level's reference files, then its starters (which must
// build but must not pass).
func testLab(t *testing.T, r *runner, l *levels.Level) {
	dir := filepath.Join("testdata", "solutions", l.ID)
	files := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(name)] = string(body)
		return nil
	})
	if err != nil {
		t.Fatalf("no reference solution: %v", err)
	}
	resp, err := r.RunLab(context.Background(), l, files)
	if err != nil {
		t.Fatal(err)
	}
	if !resp.Result.Success {
		for _, s := range resp.Steps {
			t.Logf("%v %s\n%s", s.OK, s.Title, s.Output)
		}
		t.Fatalf("phase=%s msg=%q missing=%v compile=%s", resp.Phase, resp.Result.Message, resp.Missing, resp.CompileOutput)
	}
	t.Logf("stars=%d %s", resp.Result.Stars, resp.Result.Message)

	resp, err = r.RunLab(context.Background(), l, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Phase == "compile" || resp.Result.Success {
		t.Fatalf("starters must build and must not pass: phase=%s success=%v %s", resp.Phase, resp.Result.Success, resp.CompileOutput)
	}
}

// Missions must be ordered by chapter: unlocking and "next mission" follow
// this order. (Files register levels in file-name order, so ch10.go once
// sorted before ch2.go.)
func TestLevelOrder(t *testing.T) {
	for i := 1; i < len(levels.All); i++ {
		if a, b := levels.All[i-1], levels.All[i]; a.Chapter > b.Chapter {
			t.Fatalf("%s (chapter %d) comes before %s (chapter %d)", a.ID, a.Chapter, b.ID, b.Chapter)
		}
	}
	for _, c := range levels.Chapters {
		found := false
		for _, l := range levels.All {
			found = found || l.Chapter == c.N
		}
		if !found {
			t.Errorf("chapter %d has no missions", c.N)
		}
	}
}

// Reaching the goal doesn't end the run on its own: the program has to stop
// there (or keep the goal true), so overshooting the beacon fails.
func TestGoalMustHoldWhenProgramEnds(t *testing.T) {
	r := testRunner(t)
	moves := func(n int) string {
		return "package main\n\nimport \"orbital/colony\"\n\nfunc main() {\n" +
			strings.Repeat("\tcolony.Move(colony.East)\n", n) + "}\n"
	}
	for _, c := range []struct {
		n     int
		ok    bool
		ticks int
	}{{3, true, 3}, {6, false, 0}} {
		resp, err := r.Run(context.Background(), levels.Get("1-1-first-signal"), moves(c.n))
		if err != nil {
			t.Fatal(err)
		}
		if resp.Result.Success != c.ok || c.ok && resp.Result.Ticks != c.ticks {
			t.Errorf("%d moves: success=%v ticks=%d (%s), want success=%v ticks=%d",
				c.n, resp.Result.Success, resp.Result.Ticks, resp.Result.Message, c.ok, c.ticks)
		}
		if !c.ok && !strings.Contains(resp.Result.Message, "tick 3") {
			t.Errorf("%d moves: message should say the goal was met at tick 3: %s", c.n, resp.Result.Message)
		}
	}
}
