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
	day, baselineNight, night, brownout := false, false, false, false
	playBuildingProgram(t, st, string(code), func(ev buildingProgramEvent) bool {
		if ev.State == nil {
			return false
		}
		s := ev.State
		day = day || s.Power.Day && s.Power.Supply == 10
		baselineNight = baselineNight || !s.Power.Day && s.Power.Supply == 7
		night = night || !s.Power.Day && s.Power.Supply == 25
		for _, b := range s.Buildings {
			brownout = brownout || b.Kind == "composter" && b.Enabled && !b.Powered && s.Power.Supply == 10 && s.Power.Demand == 11
		}
		if s.Tick >= 600 && s.Inv["fuel"] >= 1 {
			if !day || !baselineNight || !night || !brownout {
				t.Fatalf("missing power phases: day=%v baselineNight=%v night=%v brownout=%v", day, baselineNight, night, brownout)
			}
			return true
		}
		return false
	})
}

func TestPowerSolarPlacementWhileWaitingKeepsProductionRunning(t *testing.T) {
	st := econ.New(11)
	st.Tick, st.Credits = 60, 699
	st.Levels["composter"], st.Levels["refinery"], st.Levels["generator"] = 1, 1, 1
	st.Inv["glass"] = 2
	st.Buildings = []*econ.Building{
		{Kind: "composter", X: 0, Level: 1, Enabled: true, Recipe: "biomass", In: map[string]float64{"moss": 8}, Out: map[string]float64{}},
		{Kind: "refinery", X: 1, Level: 1, Enabled: true, Recipe: "fuel", In: map[string]float64{"biomass": 4, "water": 2}, Out: map[string]float64{}},
	}
	code := `package main
import "orbital/colony"
func main() {
 colony.Move(colony.East)
 colony.Move(colony.East)
 if err := colony.Build(colony.Generator); err != nil { panic(err) }
 for { colony.Wait() }
}`
	playBuildingProgram(t, st, code, func(ev buildingProgramEvent) bool {
		if ev.State == nil || len(ev.State.Buildings) < 3 {
			return false
		}
		s := ev.State
		if s.Power != (econ.PowerStatus{Supply: 10, Demand: 7, Day: true}) {
			t.Fatalf("tick %d: power=%+v", s.Tick, s.Power)
		}
		for _, b := range s.Buildings {
			if !b.Powered {
				t.Fatalf("tick %d: unpowered building: %+v", s.Tick, b)
			}
		}
		if s.Tick < 200 {
			return false
		}
		if s.Buildings[0].Out["biomass"] < 2 || s.Buildings[1].Out["fuel"] < 1 {
			t.Fatalf("tick %d: production stalled: %+v", s.Tick, s.Buildings)
		}
		return true
	})
}
