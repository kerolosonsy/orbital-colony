package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

// contractUnit is the fast colony clock for these tests: one tick is 2ms
// instead of 100ms, and drone actions scale with it.
const contractUnit = 2 * time.Millisecond

// contractFixture is a colony as a player who reached the contract's chapter
// might have it: earlier contracts won and their unlocks bought, plus a few
// upgrades bought with credits.
func contractFixture(c *econ.Contract) *econ.State {
	if c.Repeatable {
		return experimentFixture(c)
	}
	seed := int64(11)
	if c.ID == "c17-deep-haul" {
		seed = 7 // more geodes on layer 2
	}
	st := econ.New(seed)
	for _, prev := range econ.Contracts {
		if prev.Chapter >= c.Chapter && !(c.ID == econ.CleanStopContract && prev.ID == "c13-reed-clock") {
			continue
		}
		st.Contracts[prev.ID] = &econ.ContractRun{Status: "done", Best: 1, Tries: 1}
		for _, id := range prev.Reward.Grants {
			st.Blueprints[id] = true
			st.Levels[id]++
		}
	}
	bought := map[string]map[string]int{
		"c02-price-watch":     {"field": 1, "engine": 2, "yield.moss": 3},
		"c03-full-sweep":      {"field": 2, "engine": 2},
		"c04-right-ground":    {"field": 2, "engine": 3, "yield.lichen": 4},
		"c05-careful-harvest": {"field": 2, "engine": 4, "yield.moss": 3},
		"c06-ripe-ledger":     {"field": 3, "engine": 5, "yield.moss": 4},
		"c07-plot-type":       {"field": 3, "engine": 5},
		"c08-big-veins":       {"field": 3, "engine": 5, "veins.rate": 6, "yield.crystal": 3},
		"c09-crop-protocol":   {"field": 4, "engine": 6, "yield.moss": 4, "yield.lichen": 4, "yield.glowcap": 2},
		"c10-fleet":           {"field": 4, "engine": 6, "yield.moss": 6, "grow.moss": 3},
		"c11-dispatch":        {"field": 4, "engine": 7, "veins.rate": 8, "yield.crystal": 4},
		"c12-one-ledger":      {"field": 5, "engine": 8, "yield.moss": 8, "grow.moss": 3, "hangar": 3},
		"c13-reed-clock":      {"field": 4, "engine": 8, "yield.reed": 2},
		"c14-generic-queue":   {"cargo": 8, "field": 4, "engine": 8, "yield.moss": 6, "yield.glowcap": 4},
		"c15-shortest-run":    {"cargo": 8, "field": 4, "engine": 8, "veins.rate": 8, "yield.crystal": 4},
		"c17-deep-haul":       {"cargo": 8, "field": 4, "engine": 10, "drill": 2, "yield.geode": 3},
	}[c.ID]
	for id, l := range bought {
		st.Levels[id] = l
	}
	st.Normalize()
	if err := st.Accept(c.ID); err != nil {
		panic(err)
	}
	return st
}

type contractResult struct {
	Requires []string
	OK       bool
	Ticks    int64
	Msg      string
}

// buildColonyProgram compiles code against the colony module, like colonyRun.
func buildColonyProgram(t *testing.T, r *runner, code string) string {
	t.Helper()
	dir := t.TempDir()
	gomod := fmt.Sprintf("module player\n\ngo 1.24\n\nrequire orbital/colony v0.0.0\n\nreplace orbital/colony => %s\n", r.colonyDir)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	build := r.goCmd(ctx, dir, "build", "-o", "prog", ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return filepath.Join(dir, "prog")
}

// playContract runs code with c active on its fixture until the contract
// settles or fails, or its window has clearly passed.
func playContract(t *testing.T, r *runner, c *econ.Contract, code string) contractResult {
	t.Helper()
	if missing := contractRequirements(code, c.Requires, "en"); len(missing) > 0 {
		return contractResult{Requires: missing}
	}
	prog := buildColonyProgram(t, r, code)
	st := contractFixture(c)
	save := filepath.Join(t.TempDir(), "colony.json")
	blob, _ := json.Marshal(st)
	if err := os.WriteFile(save, blob, 0o600); err != nil {
		t.Fatal(err)
	}
	events, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer events.Close()
	limit := time.Duration(c.Window)*contractUnit*2 + 10*time.Second
	ctx, cancel := context.WithTimeout(context.Background(), limit)
	defer cancel()
	cmd := exec.CommandContext(ctx, prog)
	cmd.Env = append(os.Environ(), "ORBITAL_LANG=en", "ORBITAL_LEVEL=colony", "ORBITAL_COLONY="+save, "ORBITAL_EVENTS_FD=3", "ORBITAL_COLONY_UNIT="+contractUnit.String())
	cmd.ExtraFiles = []*os.File{writer}
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	writer.Close()
	defer func() { _ = cmd.Process.Kill(); _ = cmd.Wait() }()
	scanner := bufio.NewScanner(events)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	var finalState *econ.State
	sentStop := false
	for scanner.Scan() {
		var ev struct {
			Op    string
			OK    bool
			T     int64
			Msg   string
			State *econ.State
		}
		if json.Unmarshal(scanner.Bytes(), &ev) != nil {
			continue
		}
		if ev.State != nil {
			finalState = ev.State
		}
		if c.ID == econ.CleanStopContract && !sentStop && ev.T >= econ.CleanStopAt(st.Seed, st.Contracts[c.ID].Start) {
			sentStop = true
			if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
				t.Fatal(err)
			}
		}
		if ev.Op != "contract" {
			continue
		}
		res := contractResult{OK: ev.OK, Msg: ev.Msg, Ticks: ev.T - st.Contracts[c.ID].Start}
		if ev.OK && ev.State != nil {
			res.Ticks = ev.State.Contracts[c.ID].Best
		}
		return res
	}
	if c.ID == econ.CleanStopContract && finalState != nil {
		done, settled := finalState.FinishCleanStop()
		if settled {
			return contractResult{OK: done, Ticks: finalState.Tick - st.Contracts[c.ID].Start}
		}
	}
	t.Fatalf("%s: no contract event before the process ended or timed out\n%s", c.ID, strings.TrimSpace(output.String()))
	return contractResult{}
}

// TestContractsSolvable plays every contract's reference solution in
// testdata/colony/<id>/main.go and the colony starter on the fast clock.
// Each reference must finish well inside the window; the starter must not
// finish any contract.
func TestContractsSolvable(t *testing.T) {
	if testing.Short() {
		t.Skip("plays every contract in real time")
	}
	r := testRunner(t)
	sem := make(chan struct{}, 4)
	for _, c := range econ.Contracts {
		t.Run(c.ID, func(t *testing.T) {
			t.Parallel()
			sem <- struct{}{}
			defer func() { <-sem }()
			code, err := os.ReadFile(filepath.Join("testdata", "colony", c.ID, "main.go"))
			if err != nil {
				t.Fatal(err)
			}
			ref := playContract(t, r, c, string(code))
			if len(ref.Requires) > 0 || !ref.OK {
				t.Fatalf("reference did not finish: requires %v, %s", ref.Requires, ref.Msg)
			}
			share := float64(ref.Ticks) / float64(c.Window)
			t.Logf("%s reference finished in %d of %d ticks (%.0f%%)", c.ID, ref.Ticks, c.Window, share*100)
			if share > 0.85 {
				t.Errorf("reference used %.0f%% of the window; keep it within 40-70%%", share*100)
			}
			starter := playContract(t, r, c, levels.Colony.Starter)
			if starter.OK {
				t.Errorf("the colony starter finished the contract in %d ticks", starter.Ticks)
			}
			t.Logf("%s starter: requires %v, %s", c.ID, starter.Requires, starter.Msg)
		})
	}
}
