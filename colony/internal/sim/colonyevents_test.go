package sim

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"testing"

	"orbital/colony/econ"
)

var eventGates = []string{"c02-price-watch", "c05-careful-harvest", "c06-ripe-ledger", "c13-reed-clock"}

func eventWorld(seed int64, gates ...string) *World {
	st := econ.New(seed)
	for _, gate := range gates {
		st.Contracts[gate] = &econ.ContractRun{Best: 1}
	}
	w := NewColony(io.Discard, st, seed+99)
	w.Realtime = false
	return w
}

func TestColonyEventsDeterministicAndGated(t *testing.T) {
	for _, gate := range append([]string{""}, eventGates...) {
		t.Run(gate, func(t *testing.T) {
			a, b := eventWorld(71, gate), eventWorld(71, gate)
			var left, right bytes.Buffer
			a.out, b.out = &left, &right
			seen := map[string]bool{}
			for tick := 1; tick < 5000; tick++ {
				a.Wait(0)
				b.Wait(0)
				if !reflect.DeepEqual(a.Econ.Weather, b.Econ.Weather) || !reflect.DeepEqual(a.Econ.Shock, b.Econ.Shock) || !reflect.DeepEqual(a.Econ.Rush, b.Econ.Rush) || a.Drones[0].fault != b.Drones[0].fault {
					t.Fatal("seeded events diverged")
				}
			}
			var streams [2][]Event
			for i, stream := range [][]byte{left.Bytes(), right.Bytes()} {
				for _, line := range bytes.Split(stream, []byte{'\n'}) {
					var ev Event
					if json.Unmarshal(line, &ev) == nil && (ev.Op == "storm" || ev.Op == "fault" || ev.Op == "shock" || ev.Op == "rush") {
						ev.Translations = nil
						streams[i] = append(streams[i], ev)
					}
				}
			}
			if !reflect.DeepEqual(streams[0], streams[1]) {
				t.Fatal("live event streams diverged")
			}
			for _, line := range bytes.Split(left.Bytes(), []byte{'\n'}) {
				var ev Event
				if json.Unmarshal(line, &ev) == nil {
					if ev.Op == "fault" {
						var wire map[string]json.RawMessage
						if err := json.Unmarshal(line, &wire); err != nil {
							t.Fatal(err)
						}
						if _, ok := wire["d"]; !ok {
							t.Fatal("fault event omitted drone 0")
						}
					}
					seen[ev.Op] = true
				}
			}
			for op, required := range map[string]string{"storm": "c13-reed-clock", "fault": "c05-careful-harvest", "shock": "c02-price-watch", "rush": "c06-ripe-ledger"} {
				if seen[op] != (gate == required) {
					t.Fatalf("%s fired=%v with only gate %q", op, seen[op], gate)
				}
			}
		})
	}
}

func TestColonyStormWarningShelterDamageAndClear(t *testing.T) {
	for _, sheltered := range []bool{false, true} {
		t.Run(map[bool]string{false: "ignores warning", true: "returns home"}[sheltered], func(t *testing.T) {
			w := eventWorld(9, "c13-reed-clock")
			warn := int(w.Econ.Weather.At)
			clear := w.StormClear()
			w.X = 1
			w.Drones = append(w.Drones, &Drone{ID: 1, X: 1, Hold: map[string]float64{}})
			w.PlantAt(1, 0, Moss, warn+1000)
			w.PlantAt(2, 0, Moss, warn)
			w.Tick = warn
			w.advanceProduction(w.Tick)
			select {
			case <-w.StormWarning():
			default:
				t.Fatal("no warning")
			}
			if clear != w.StormClear() {
				t.Fatal("first warning replaced the clear channel")
			}
			if w.Econ.Weather.Storm != "warn" || w.Econ.Weather.At != int64(warn+60) {
				t.Fatal("wrong warning interval")
			}
			if sheltered {
				if !w.Move(0, West) || !w.Move(1, West) {
					t.Fatal("could not return home")
				}
			}
			w.Tick = warn + 60
			w.advanceProduction(w.Tick)
			if w.done || w.At(1, 0).Entity != Empty || w.At(2, 0).Entity != Moss {
				t.Fatal("storm damaged wrong plants or ended run")
			}
			for id := 0; id < 2; id++ {
				err := w.Fault(id)
				if errors.Is(err, ErrStorm) == sheltered {
					t.Fatalf("drone %d fault=%v, sheltered=%v", id, err, sheltered)
				}
				if !sheltered {
					if _, err := w.TryHarvest(id); !errors.Is(err, ErrStorm) {
						t.Fatalf("stunned harvest: %v", err)
					}
					if w.Repair(id) {
						t.Fatal("repair bypassed stun")
					}
				}
			}
			w.syncInfrastructure()
			if !sheltered && w.Econ.Fleet[0].Stunned <= 0 {
				t.Fatal("stun missing from fleet")
			}
			w.Tick = warn + 60 + 150
			w.advanceProduction(w.Tick)
			select {
			case <-clear:
			default:
				t.Fatal("clear channel not closed")
			}
			if w.Fault(0) != nil || w.Econ.Weather.Storm != "none" {
				t.Fatal("storm did not clear")
			}
		})
	}
}

func TestColonyFaultBlocksUntilRepairAndPersists(t *testing.T) {
	w := eventWorld(31, "c05-careful-harvest")
	w.Econ.Levels["api.try"] = 1
	for w.Fault(0) == nil && w.Tick < 2000 {
		w.Wait(0)
	}
	if !errors.Is(w.Fault(0), ErrFault) {
		t.Fatal("fault never fired")
	}
	for i := 0; i < 3; i++ {
		if w.Move(0, East) {
			t.Fatal("faulted drone moved")
		}
		_, err := w.TryHarvest(0)
		var disabled DroneError
		if !errors.Is(err, ErrFault) || !errors.As(err, &disabled) || disabled.Drone != 0 {
			t.Fatalf("wrong typed failure: %v", err)
		}
		if err := w.Build(0, "composter"); !errors.Is(err, ErrFault) {
			t.Fatalf("build failure lost: %v", err)
		}
		if err := w.TryDigDown(0); !errors.Is(err, ErrFault) {
			t.Fatalf("dig failure lost: %v", err)
		}
	}
	w.syncEcon()
	blob, _ := json.Marshal(w.Econ)
	var saved econ.State
	if err := json.Unmarshal(blob, &saved); err != nil {
		t.Fatal(err)
	}
	resumed := NewColony(io.Discard, &saved, 123)
	resumed.Realtime = false
	if !errors.Is(resumed.Fault(0), ErrFault) || !resumed.Repair(0) || resumed.Fault(0) != nil {
		t.Fatal("fault did not survive save or repair")
	}
}

func TestColonyShockPriceRevertsAndRushExpires(t *testing.T) {
	w := eventWorld(12, "c02-price-watch", "c06-ripe-ledger")
	w.Econ.Levels["api.prices"] = 1
	for w.Econ.Shock == nil {
		w.Wait(0)
	}
	s := *w.Econ.Shock
	if got := w.Price(entityByName[s.Good]); got != econ.PriceV3[s.Good]*s.Factor {
		t.Fatalf("shock price=%g", got)
	}
	w.Tick = int(s.Until)
	w.advanceProduction(w.Tick)
	if w.Econ.Shock != nil || w.Price(entityByName[s.Good]) != econ.PriceV3[s.Good] {
		t.Fatal("shock did not revert")
	}
	for w.Econ.Rush == nil {
		w.Wait(0)
	}
	r := *w.Econ.Rush
	good, n, ok := w.Rush()
	if !ok || good.String() != r.Good || n != r.Amount {
		t.Fatal("rush API does not match state")
	}
	w.Tick = int(r.Until)
	w.advanceProduction(w.Tick)
	if _, _, ok := w.Rush(); ok || w.Econ.Rush != nil {
		t.Fatal("rush did not expire")
	}
}

func TestColonyStormSuppressesOnlySolar(t *testing.T) {
	w := eventWorld(11)
	powerBuilding(w, "generator", 0)
	if supply := w.Econ.PowerSupply(1); supply != 10 {
		t.Fatalf("daytime solar=%g", supply)
	}
	w.Econ.Weather = econ.Weather{Storm: "active", At: 151}
	if supply := w.Econ.PowerSupply(1); supply != econ.BaselinePower {
		t.Fatalf("storm solar still generated: %g", supply)
	}
	g := w.Econ.Buildings[0]
	g.Mode = "burner"
	g.In["fuel"] = 10
	if supply := w.Econ.ConsumePower(1); supply <= econ.BaselinePower || g.In["fuel"] >= 10 {
		t.Fatal("storm disabled fueled power")
	}
}
