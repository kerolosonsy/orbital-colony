package econ

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestResearchPurchaseSpendsOnceAndPreservesFundsOnFailure(t *testing.T) {
	s := New(1)
	n := FindNode("research.engine")
	s.Credits, s.RP = 1000, n.RP-0.5
	s.Levels["engine"] = 3
	s.Inv["moss"] = 100
	if err := s.Buy(n.ID); err == nil || !strings.Contains(err.Error(), "RP") || s.Credits != 1000 || s.RP != n.RP-0.5 || s.Level(n.ID) != 0 {
		t.Fatalf("RP shortage spent funds or lacked reason: %v", err)
	}
	s.RP = n.RP
	if err := s.Buy(n.ID); err != nil || s.RP != 0 || s.Level(n.ID) != 1 {
		t.Fatalf("exact RP purchase: %v", err)
	}
	s.Levels["engine"] = 0
	if err := s.Buy(n.ID); err != nil || s.RP != 0 || s.Level(n.ID) != 2 {
		t.Fatalf("repeat charged RP or parents: %v", err)
	}
}

func TestExperimentDecayCountsWinsNotAttemptsAndResetsAtDayBoundary(t *testing.T) {
	s := New(1)
	id := "e03-sweep"
	s.Contracts["c03-full-sweep"] = &ContractRun{Best: 1}
	if err := s.Accept(id); err != nil {
		t.Fatal(err)
	}
	s.Fail()
	for _, tc := range []struct {
		tick int64
		rp   float64
	}{{1198, 8}, {1199, 6.4}, {1200, 8}, {1201, 6.4}} {
		s.Tick = tc.tick - 1
		if err := s.Accept(id); err != nil {
			t.Fatal(err)
		}
		if _, err := s.Settle(tc.tick); err == nil {
			t.Fatal("unfinished experiment paid")
		}
		s.RecordSale(0, "moss", 40, 40, 1)
		before := s.RP
		reward, err := s.Settle(tc.tick)
		if err != nil || math.Abs(reward.RP-tc.rp) > 1e-9 || math.Abs(s.RP-before-tc.rp) > 1e-9 {
			t.Fatalf("tick %d reward %+v %v", tc.tick, reward, err)
		}
		blob, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		var loaded State
		if err := json.Unmarshal(blob, &loaded); err != nil {
			t.Fatal(err)
		}
		loaded.Normalize()
		s = &loaded
	}
	if s.Contracts[id].RepeatsToday != 2 || s.Contracts[id].Tries != 5 {
		t.Fatalf("counts lost across save: %+v", s.Contracts[id])
	}
}

func TestResearchOldTreeSaveKeepsEveryOwnedLevelAndPossession(t *testing.T) {
	blob := `{"version":3,"seed":4,"credits":123,"inv":{"moss":7,"plate":3},"levels":{"field":2,"engine":9,"yield.moss":12,"api.try":1,"composter":4,"lab":2},"blueprints":{"api.try":true},"claimed":{"chapter-1":true},"contracts":{"c05-careful-harvest":{"status":"done","best":23}}}`
	var s State
	if err := json.Unmarshal([]byte(blob), &s); err != nil {
		t.Fatal(err)
	}
	levels := map[string]int{}
	for id, l := range s.Levels {
		levels[id] = l
	}
	s.Normalize()
	if s.Version != 3 || s.RP != 0 || !reflect.DeepEqual(s.Levels, levels) || s.Credits != 123 || s.Inv["plate"] != 3 || !s.Blueprints["api.try"] || !s.Claimed["chapter-1"] || !s.Won("c05-careful-harvest") {
		t.Fatal("old save lost possessions")
	}
	s.Credits = 1e12
	s.Inv["moss"] = 1000
	if err := s.Buy("composter"); err != nil || s.Level("composter") != 5 || s.RP != 0 {
		t.Fatalf("owned building retroactively locked: %v", err)
	}
}

func TestResearchFullTreeContentHasReachableAffordablePrerequisites(t *testing.T) {
	if err := Validate(Tree); err != nil {
		t.Fatal(err)
	}
	if len(Tree) < 60 {
		t.Fatalf("only %d nodes", len(Tree))
	}
	counts := map[string]int{}
	for _, n := range Tree {
		counts[n.Branch]++
		if n.ID == "base" {
			continue
		}
		if len(n.After) != 1 || n.After[0].ID != "base" {
			if n.RP <= 0 {
				t.Errorf("%s lacks RP", n.ID)
			}
		}
		for _, p := range n.After {
			if p.ID != "base" && Find(p.ID).Max > 0 && p.Level > Find(p.ID).Max {
				t.Errorf("%s impossible parent %+v", n.ID, p)
			}
		}
		if src := BlueprintSource(n.ID); needsBlueprint(n.ID) && src != nil && n.Contract != src.ID {
			t.Errorf("%s lacks contract", n.ID)
		}
	}
	t.Logf("nodes=%d branches=%v", len(Tree), counts)
}

func TestResearchEffectsChangeTheirNamedNumbers(t *testing.T) {
	for _, n := range Tree {
		if !strings.HasPrefix(n.ID, "research.") {
			continue
		}
		t.Run(n.ID, func(t *testing.T) {
			s := New(1)
			s.Levels[n.ID] = 2
			var got, want float64
			suffix := strings.TrimPrefix(n.ID, "research.")
			switch {
			case strings.HasPrefix(suffix, "yield."):
				got, want = s.YieldMult(strings.TrimPrefix(suffix, "yield.")), 1.3
			case strings.HasPrefix(suffix, "grow."):
				g := strings.TrimPrefix(suffix, "grow.")
				got, want = float64(s.GrowTime(g)), math.Trunc(float64(GrowTicks[g])/1.2)
			case strings.HasPrefix(suffix, "price."):
				g := strings.TrimPrefix(suffix, "price.")
				got, want = s.Price(g), PriceV3[g]*1.2
			case strings.HasPrefix(suffix, "recovery."):
				g := strings.TrimPrefix(suffix, "recovery.")
				s.Demand[g] = 0.5
				s.Recover(1)
				got, want = s.Demand[g], 0.501875
			case strings.HasPrefix(suffix, "speed."):
				got, want = s.BuildingSpeedMult(strings.TrimPrefix(suffix, "speed.")), 1.4
			case strings.HasPrefix(suffix, "buffer."):
				got, want = s.BuildingCapacity(strings.TrimPrefix(suffix, "buffer."), 0), 16
			case suffix == "engine":
				got, want = s.ActionMillis(), 260/1.16
			case suffix == "vein-rate":
				got, want = float64(s.VeinEvery()), 114
			case suffix == "vein-size":
				got, want = float64(s.VeinMax()), 6
			case suffix == "sift-stock":
				got, want = float64(s.SiftStock()), 5
			case suffix == "sift-refill":
				got, want = float64(s.SiftRefillTicks()), 72
			case suffix == "melt":
				got, want = float64(s.MeltRegrowthTicks()), 429
			case suffix == "placement":
				cost := s.BuildingPlacementCost("composter")
				got, want = cost.Credits, 20
				if cost.Goods["moss"] != 3.2 || PlacementCost("composter").Goods["moss"] != 4 {
					t.Fatal("discount not applied to goods or mutated catalog")
				}
			case suffix == "gain":
				got, want = s.ResearchGainMult(), 1.4
			default:
				t.Fatal("node has no tested effect")
			}
			if math.Abs(got-want) > 1e-9 {
				t.Fatalf("effect=%g want %g", got, want)
			}
		})
	}
}

func TestResearchBuildingSpeedAndBuffersApplyToProduction(t *testing.T) {
	s := New(1)
	s.Levels["research.speed.composter"] = 1
	s.Levels["research.buffer.composter"] = 1
	b := &Building{Kind: "composter", Recipe: "biomass", Enabled: true, In: map[string]float64{}, Out: map[string]float64{"biomass": 8}}
	s.ConfigureBuilding(b)
	if n, err := b.Put("moss", 12); err != nil || n != 12 {
		t.Fatalf("input buffer: %g %v", n, err)
	}
	if err := b.Step(33, true); err != nil || b.Out["biomass"] != 8 {
		t.Fatalf("finished before 34 ticks: %v", err)
	}
	if err := b.Step(1, true); err != nil || b.Out["biomass"] != 9 {
		t.Fatalf("speed or output capacity did not apply: %v", err)
	}
}
