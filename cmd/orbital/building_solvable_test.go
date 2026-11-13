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

type buildingProgramEvent struct {
	Op    string
	OK    bool
	Inv   map[string]float64
	State *econ.State
}

// Uses the contract harness's compiler, saved-state protocol and fast clock.
func playBuildingProgram(t *testing.T, st *econ.State, code string, accept func(buildingProgramEvent) bool) {
	t.Helper()
	prog := buildColonyProjectProgram(t, testRunner(t), colonyProject{Files: map[string]string{"main.go": code}, Race: true})
	save := filepath.Join(t.TempDir(), "colony.json")
	blob, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(save, blob, 0o600); err != nil {
		t.Fatal(err)
	}
	events, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer events.Close()
	defer writer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
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
		var ev buildingProgramEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			t.Fatal(err)
		}
		if accept(ev) {
			return
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	// Wait before reading the shared output buffer, including on timeout or panic.
	_ = cmd.Wait()
	t.Fatalf("program ended before expected colony behavior: %s", output.String())
}

func TestBuildingReferenceProducesAndExportsFuel(t *testing.T) {
	code, err := os.ReadFile(filepath.Join("testdata", "colony", "e2-production", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	st := econ.New(11)
	st.Levels["composter"], st.Levels["refinery"] = 1, 1
	st.Credits, st.Inv["moss"] = 110, 20
	st.At(1, 0).G = "ice"
	fuelCollected, fuelSold, melted, buildings := false, false, false, 0
	playBuildingProgram(t, st, string(code), func(ev buildingProgramEvent) bool {
		if ev.Op == "building" {
			buildings++
		}
		if ev.Op == "melt" && ev.OK {
			melted = true
		}
		if ev.Op == "unload" && ev.OK && ev.Inv["fuel"] >= 1 {
			fuelCollected = true
		}
		if ev.Op == "export" && ev.OK && fuelCollected && ev.Inv["fuel"] == 0 {
			fuelSold = true
		}
		if fuelSold && ev.State != nil {
			if !melted || buildings == 0 || len(ev.State.Buildings) != 2 || ev.State.Earned != econ.PriceV3["fuel"] || ev.State.Inv["fuel"] != 0 {
				t.Fatalf("production/export did not persist: %+v", ev.State)
			}
			return true
		}
		return false
	})
	if !fuelSold {
		t.Fatal("reference never sold fuel")
	}
}

func TestColonyStarterStillBuildsAndRunsWithE2(t *testing.T) {
	planted, harvested := false, false
	playBuildingProgram(t, econ.New(11), levels.Colony.Starter, func(ev buildingProgramEvent) bool {
		if ev.Op == "plant" && ev.OK {
			planted = true
		}
		if ev.Op == "harvest" && ev.OK {
			harvested = true
		}
		return planted && harvested && ev.Op == "export" && ev.OK
	})
}
