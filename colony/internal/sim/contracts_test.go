package sim

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"orbital/colony/econ"
)

func contractEvents(t *testing.T, trace *bytes.Buffer) []econEvent {
	t.Helper()
	var out []econEvent
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev econEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "contract" || ev.Op == "econ" {
			out = append(out, ev)
		}
	}
	return out
}

func TestColonyContractCountsDenyAndFail(t *testing.T) {
	t.Setenv("ORBITAL_COLONY_UNIT", "1ms")
	st := econ.New(3)
	st.Levels["api.scanner"] = 1
	st.Contracts["c04-right-ground"] = &econ.ContractRun{Status: "done", Best: 9}
	if err := st.Accept("c05-careful-harvest"); err != nil {
		t.Fatal(err)
	}
	var trace bytes.Buffer
	w := NewColony(&trace, st, 3)
	w.Start()
	if w.ActionDur >= 10*1000*1000 {
		t.Fatalf("action duration %v ignores the fast clock", w.ActionDur)
	}
	if w.CanHarvest(0) {
		t.Fatal("denied scanner still worked")
	}
	if _, err := w.TryHarvest(0); err != ErrEmpty {
		t.Fatalf("lent TryHarvest: %v", err)
	}
	w.Till(0)
	w.Plant(0, Moss)
	if w.Harvest(0) {
		t.Fatal("harvested an unripe plant")
	}
	w.mu.Lock()
	r := st.Contracts["c05-careful-harvest"]
	if r.Lost != 1 || r.Cmds[0] != 5 {
		t.Fatalf("run = %+v", r)
	}
	w.econSnapshot("econ", "", true, false)
	w.checkContract()
	w.mu.Unlock()
	if st.Active != "" || r.Status != "failed" {
		t.Fatalf("contract not failed: %+v", r)
	}
	evs := contractEvents(t, &trace)
	last := evs[len(evs)-1]
	if last.Op != "contract" || last.OK || !strings.Contains(last.Msg, "Careful Harvest") && !strings.Contains(last.Msg, "الحصاد الحذر") || last.State == nil {
		t.Fatalf("fail event = %+v", last)
	}
	if prev := evs[len(evs)-2]; prev.Contract == nil || prev.Contract.ID != "c05-careful-harvest" || prev.Progress == "" {
		t.Fatalf("econ event lacks contract progress: %+v", prev)
	}
}

func TestColonyContractSettlesFromSales(t *testing.T) {
	t.Setenv("ORBITAL_COLONY_UNIT", "1ms")
	st := econ.New(3)
	st.Inv["moss"] = 60
	if err := st.Accept("c03-full-sweep"); err != nil {
		t.Fatal(err)
	}
	var trace bytes.Buffer
	w := NewColony(&trace, st, 3)
	w.Start()
	if earned := w.Export(0, Moss, 0); earned <= 0 {
		t.Fatal("export failed")
	}
	w.mu.Lock()
	w.checkContract()
	w.mu.Unlock()
	if !st.Blueprints["crop.lichen"] || st.Contracts["c03-full-sweep"].Status != "done" {
		t.Fatalf("contract not settled: %+v", st.Contracts["c03-full-sweep"])
	}
	if evs := contractEvents(t, &trace); !evs[len(evs)-1].OK {
		t.Fatal("missing success event")
	}
}

func TestColonyContractLendsDrones(t *testing.T) {
	t.Setenv("ORBITAL_COLONY_UNIT", "1ms")
	st := econ.New(3)
	st.Accept("c10-fleet")
	w := NewColony(&bytes.Buffer{}, st, 3)
	w.Start()
	if a, b, c := w.Spawn(), w.Spawn(), w.Spawn(); a != 1 || b != 2 || c != -1 {
		t.Fatalf("spawned %d %d %d with three lent drones", a, b, c)
	}
}
