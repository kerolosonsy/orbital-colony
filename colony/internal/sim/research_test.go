package sim

import (
	"math"
	"testing"

	"orbital/colony/econ"
)

func TestResearchLabConvertsCraftedTiersWithoutOutputAcrossSave(t *testing.T) {
	for _, tc := range []struct {
		good Entity
		rp   float64
	}{{Biomass, 1}, {Circuit, 5}, {Processor, 25}, {Core, 125}} {
		t.Run(tc.good.String(), func(t *testing.T) {
			w, trace := e2World()
			w.Econ.Levels["lab"] = 1
			w.Econ.Credits = 1000 // enough left after placement to keep paying rent
			w.Econ.Inv["plate"] = 2
			w.Econ.Inv[tc.good.String()] += 1
			if err := w.Build(0, "lab"); err != nil {
				t.Fatal(err)
			}
			if n, err := w.BuildingPut(0, Point{}, tc.good, 1); err != nil || n != 1 {
				t.Fatalf("Put=%d %v", n, err)
			}
			putAt := w.Tick
			w.Tick = putAt + 20
			w.advanceProduction(w.Tick)
			if w.Research() != 0 || w.Econ.Buildings[0].In[tc.good.String()] != 0 || w.BuildingStatus(Point{}).Progress != 0.25 {
				t.Fatal("lab did not consume at start or paid early")
			}
			loaded := reloadColony(t, w)
			for _, world := range []*World{w, loaded} {
				world.Tick = putAt + 79
				world.advanceProduction(world.Tick)
				if world.Research() != 0 {
					t.Fatal("lab paid early")
				}
				world.Tick++
				world.advanceProduction(world.Tick)
				if world.Research() != tc.rp || len(world.BuildingStatus(Point{}).Out) != 0 {
					t.Fatal("wrong RP or lab produced output")
				}
			}
			events := buildingEvents(t, trace)
			if len(events) != 4 || events[2].Building.In[tc.good.String()] != 0 || events[3].Building.Wear != 1 {
				t.Fatalf("lab events missing: %+v", events)
			}
		})
	}
}

func TestResearchLabRejectsRawGoodsAndPausesDisabledOrWorn(t *testing.T) {
	w, _ := e2World()
	w.Econ.Levels["lab"] = 1
	w.Econ.Credits = 1000 // enough left after placement to keep paying rent
	w.Econ.Inv["plate"] = 2
	if err := w.Build(0, "lab"); err != nil {
		t.Fatal(err)
	}
	w.Econ.Inv["moss"] = 3
	if _, err := w.BuildingPut(0, Point{}, Moss, 1); err != econ.ErrBadRecipe {
		t.Fatalf("lab accepted raw moss: %v", err)
	}
	b := w.Econ.Buildings[0]
	b.In["biomass"] = 1
	b.Enabled = false
	w.advanceProduction(100)
	if b.In["biomass"] != 1 || w.Research() != 0 {
		t.Fatal("disabled lab consumed")
	}
	b.Enabled = true
	b.Wear = 100
	w.advanceProduction(200)
	if b.In["biomass"] != 1 {
		t.Fatal("worn lab consumed")
	}
	b.Wear = 0
	w.Econ.Levels["research.gain"] = 1
	w.advanceProduction(280)
	if math.Abs(w.Research()-1.2) > 1e-9 {
		t.Fatalf("research multiplier not applied: %g", w.Research())
	}
}

func TestResearchExperimentSeedChangesOnEachAttempt(t *testing.T) {
	s := econ.New(11)
	s.Contracts["c03-full-sweep"] = &econ.ContractRun{Best: 1}
	var previous int64
	for attempt := 0; attempt < 3; attempt++ {
		if err := s.Accept("e03-sweep"); err != nil {
			t.Fatal(err)
		}
		seed := researchSeed(s, 11)
		if seed == previous || seed == 11 {
			t.Fatal("experiment repeated seed")
		}
		previous = seed
		s.Fail()
	}
}
