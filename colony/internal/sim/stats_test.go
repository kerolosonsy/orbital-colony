package sim

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"orbital/colony/econ"
)

func TestColonyStatsScriptCountsEachDroneAndResetsOnReload(t *testing.T) {
	st := econ.New(1)
	st.Tick = 9000
	st.Levels["api.spawn"], st.Levels["hangar"] = 1, 2
	var trace bytes.Buffer
	w := NewColony(&trace, st, 1)
	w.Realtime = false
	w.Start()
	if w.Spawn() != 1 || w.Spawn() != 2 {
		t.Fatal("could not spawn scripted drones")
	}
	w.Move(0, East)
	w.Wait(0)
	w.Harvest(0) // A failed work attempt does not make the move productive.
	w.Move(0, North)
	w.Till(0)
	w.Move(0, West)
	w.Move(0, South)
	w.Move(0, South) // Failed next move still ends the previous move's interval.
	w.Harvest(0)
	w.Move(1, East)
	w.Move(1, East)
	w.Move(1, East)
	w.Wait(1)
	w.GetEntity(1) // Locked scan ties move failures; lexical order picks harvest.
	w.econSnapshot("econ", "", true, true)
	want := &runStats{Drones: []droneStats{{ID: 0, Work: 3, Move: 4, Wait: 1, Failed: 3, IdleMoves: 3}, {ID: 1, Work: 1, Move: 2, Wait: 1, Failed: 2, IdleMoves: 2}, {ID: 2}}, TopFail: "harvest", Ticks: 14}
	if got := w.statsSnapshot(); !reflect.DeepEqual(got, want) {
		t.Fatalf("script stats %+v, want %+v", got, want)
	}
	if got := w.stats[0].commands; !reflect.DeepEqual(got, map[string]int{"move": 5, "wait": 1, "harvest": 2, "till": 1}) {
		t.Fatalf("operation counts: %v", got)
	}
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev econEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.State != nil && !reflect.DeepEqual(ev.Stats, want) {
			t.Fatalf("full event lost stats: %s", line)
		}
	}
	reloaded := reloadColony(t, w)
	reset := &runStats{Drones: []droneStats{{ID: 0}, {ID: 1}, {ID: 2}}}
	if got := reloaded.statsSnapshot(); !reflect.DeepEqual(got, reset) {
		t.Fatalf("previous run leaked into next run: %+v", got)
	}
}

func TestColonyStatsWaitCountsDurationWithoutIdleGaps(t *testing.T) {
	for _, tc := range []struct {
		name string
		real bool
		slow bool
		want int
	}{{"normal", false, false, 2}, {"flat battery", false, true, 4}, {"fractional realtime ticks", true, false, 3}} {
		t.Run(tc.name, func(t *testing.T) {
			w := NewColony(io.Discard, econ.New(1), 1)
			w.Realtime, w.Unit, w.ActionDur = tc.real, time.Millisecond, 1500*time.Microsecond
			if tc.slow {
				w.Econ.Buildings = []*econ.Building{{Kind: "refinery"}}
				w.Drones[0].Battery = 0
			}
			w.Start()
			w.Wait(0)
			w.Wait(0)
			if got := w.statsSnapshot(); got.Drones[0].Wait != tc.want || got.TopFail != "" {
				t.Fatalf("wait duration: %+v", got)
			}
		})
	}
}

func TestColonyStatsAccompanyEveryStateAndOnlyStateEvents(t *testing.T) {
	var trace bytes.Buffer
	w := NewColony(&trace, econ.New(1), 1)
	w.Realtime = false
	for _, op := range []string{"econ", "bought", "state", "contract", "rank"} {
		for _, full := range []bool{false, true} {
			ev := econEvent{Op: op}
			if full {
				ev.State = w.Econ
			}
			trace.Reset()
			w.emitRaw(ev)
			var wire map[string]json.RawMessage
			if err := json.Unmarshal(trace.Bytes(), &wire); err != nil {
				t.Fatal(err)
			}
			if (wire["stats"] != nil) != full {
				t.Fatalf("%s full=%v: %s", op, full, trace.Bytes())
			}
			if full && string(wire["stats"]) != `{"drones":[{"id":0,"work":0,"move":0,"wait":0,"failed":0,"idleMoves":0}],"topFail":"","ticks":0}` {
				t.Fatalf("stats wire format: %s", wire["stats"])
			}
		}
	}
}
