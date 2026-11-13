package main

import (
	"os"
	"path/filepath"
	"testing"

	"orbital/colony/econ"
)

func TestPowerReferenceProducesAcrossDayAndNight(t *testing.T) {
	code, err := os.ReadFile(filepath.Join("testdata", "colony", "e3-power", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	st := econ.New(11)
	for _, kind := range []string{"generator", "composter", "refinery", "kiln"} {
		st.Levels[kind] = 1
	}
	st.Credits = 1000
	st.Inv = map[string]float64{"moss": 20, "glass": 2, "biomass": 2, "silica": 6, "water": 1}
	st.At(1, 0).G = "ice"
	day, night, brownout := false, false, false
	playBuildingProgram(t, st, string(code), func(ev buildingProgramEvent) bool {
		if ev.State == nil {
			return false
		}
		s := ev.State
		day = day || s.Power.Day && s.Power.Supply == 10
		night = night || !s.Power.Day && s.Power.Supply == 25
		for _, b := range s.Buildings {
			brownout = brownout || b.Kind == "composter" && b.Enabled && !b.Powered && s.Power.Supply == 10 && s.Power.Demand == 11
		}
		if s.Tick >= 600 && s.Inv["fuel"] >= 1 {
			if !day || !night || !brownout {
				t.Fatalf("missing power phases: day=%v night=%v brownout=%v", day, night, brownout)
			}
			return true
		}
		return false
	})
}
