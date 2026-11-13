package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		if prev.Chapter >= c.Chapter {
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
		"c14-generic-queue":   {"field": 4, "engine": 8, "yield.moss": 6, "yield.glowcap": 4},
		"c15-shortest-run":    {"field": 4, "engine": 8, "veins.rate": 8, "yield.crystal": 4},
		"c16-tested-planner":  {"engine": 8, "yield.moss": 10, "grow.moss": 3},
		"c17-deep-haul":       {"field": 4, "engine": 8, "drill": 2, "yield.geode": 3},
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
	return buildColonyProjectProgram(t, r, colonyProject{Files: map[string]string{"main.go": code}})
}

func buildColonyProjectProgram(t *testing.T, r *runner, project colonyProject) string {
	t.Helper()
	dir := t.TempDir()
	if err := r.writeModule(dir, levels.Colony); err != nil {
		t.Fatal(err)
	}
	if err := writeProjectFiles(dir, project.Files); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	args := []string{"build", "-o", "prog"}
	if project.Race {
		args = append(args, "-race")
	}
	build := r.goCmd(ctx, dir, append(args, ".")...)
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	return filepath.Join(dir, "prog")
}

// playContract runs code with c active on its fixture until the contract
// settles or fails, or its window has clearly passed.
func playContract(t *testing.T, r *runner, c *econ.Contract, project colonyProject) contractResult {
	t.Helper()
	if missing := projectRequirements(project.Files, c.Requires, "en"); len(missing) > 0 {
		return contractResult{Requires: missing}
	}
	if c.ID == "c16-tested-planner" {
		response, err := r.checkColonyPlanner(withLang(context.Background(), "en"), project.Files)
		if err != nil {
			t.Fatal(err)
		}
		if !response.Result.Success {
			return contractResult{Requires: []string{response.Result.Message}}
		}
	}
	prog := buildColonyProjectProgram(t, r, project)
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
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		if strings.Contains(output.String(), "WARNING: DATA RACE") {
			t.Errorf("reference has a data race:\n%s", output.String())
		}
	}()
	scanner := bufio.NewScanner(events)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	for scanner.Scan() {
		var ev struct {
			Op    string
			OK    bool
			T     int64
			Msg   string
			State *econ.State
		}
		if json.Unmarshal(scanner.Bytes(), &ev) != nil || ev.Op != "contract" {
			continue
		}
		res := contractResult{OK: ev.OK, Msg: ev.Msg, Ticks: ev.T - st.Contracts[c.ID].Start}
		if ev.OK && ev.State != nil {
			res.Ticks = ev.State.Contracts[c.ID].Best
		}
		return res
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
			files := readColonyFixture(t, c.ID)
			ref := playContract(t, r, c, colonyProject{Files: files, Race: c.ID == "c12-one-ledger"})
			if len(ref.Requires) > 0 || !ref.OK {
				t.Fatalf("reference did not finish: requires %v, %s", ref.Requires, ref.Msg)
			}
			share := float64(ref.Ticks) / float64(c.Window)
			t.Logf("%s reference finished in %d of %d ticks (%.0f%%)", c.ID, ref.Ticks, c.Window, share*100)
			if share > 0.85 {
				t.Errorf("reference used %.0f%% of the window; keep it within 40-70%%", share*100)
			}
			starter := playContract(t, r, c, colonyProject{Files: map[string]string{"main.go": levels.Colony.Starter}, Race: c.ID == "c12-one-ledger"})
			if starter.OK {
				t.Errorf("the colony starter finished the contract in %d ticks", starter.Ticks)
			}
			t.Logf("%s starter: requires %v, %s", c.ID, starter.Requires, starter.Msg)
		})
	}
}
