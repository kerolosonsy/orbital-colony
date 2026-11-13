package sim

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"orbital/colony/econ"
)

func TestSeasonModifiers(t *testing.T) {
	for _, season := range econ.Seasons {
		t.Run(season.ID, func(t *testing.T) {
			st := season.New()
			w := NewColony(io.Discard, st, 999)
			w.Realtime = false
			switch season.ID {
			case "one-drone":
				if w.Spawn() >= 0 {
					t.Fatal("single-drone limit allowed a spawn")
				}
			case "blind-flight":
				if w.gate("api.radar") == "" || w.gate("api.scanner") == "" {
					t.Fatal("sensors are usable")
				}
			case "volatile-market":
				var output bytes.Buffer
				w.out = &output
				w.Tick = 1000
				w.advanceProduction(w.Tick)
				shocks := 0
				for _, line := range bytes.Split(output.Bytes(), []byte{'\n'}) {
					var ev Event
					json.Unmarshal(line, &ev)
					if ev.Op == "shock" {
						shocks++
					}
				}
				if shocks != 5 {
					t.Fatalf("got %d shocks, want 5", shocks)
				}
			case "storm-front":
				w.Tick = 340
				w.advanceProduction(w.Tick)
				if st.Weather.Storm != "warn" || st.Weather.At != 400 {
					t.Fatalf("warning: %+v", st.Weather)
				}
				w.Tick = 400
				w.advanceProduction(w.Tick)
				if !w.Storm.active {
					t.Fatal("storm not active at 400")
				}
				w.Tick = 740
				w.advanceProduction(w.Tick)
				if st.Weather.Storm != "warn" || st.Weather.At != 800 {
					t.Fatalf("second warning: %+v", st.Weather)
				}
			case "fixed-kit":
				st.Credits = 1e9
				st.RP = 1e9
				st.Rank = 100
				if st.Buy("engine") == nil {
					t.Fatal("purchase allowed")
				}
			case "small-field":
				if w.W != econ.FieldSize(st.Level("field"))/2 {
					t.Fatal("field not halved")
				}
				st.Normalize()
				if st.Size != w.W {
					t.Fatal("normalization expanded field")
				}
			}
		})
	}
}

func TestSeasonStopsBeforeActionBeyondBoundary(t *testing.T) {
	st := econ.FindSeason("one-drone").New()
	st.Tick = 5999
	var output bytes.Buffer
	w := NewColony(&output, st, 1)
	w.Realtime = false
	exited := false
	w.exit = func(code int) {
		if code != 0 {
			t.Fatalf("exit %d", code)
		}
		exited = true
	}
	w.Till(0)
	if !exited || !w.done || st.Tick != 6000 || w.At(0, 0).Ground != Regolith {
		t.Fatalf("boundary mutated: exited=%v tick=%d", exited, st.Tick)
	}
	if !bytes.Contains(output.Bytes(), []byte(`"op":"season-end"`)) {
		t.Fatal("missing final snapshot")
	}
}
