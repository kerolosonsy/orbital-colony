package econ

import (
	"encoding/json"
	"math"
	"reflect"
	"strings"
	"testing"
)

func formerCap(u *Upgrade) int {
	switch {
	case strings.HasPrefix(u.ID, "yield."):
		return 12
	case strings.HasPrefix(u.ID, "grow."):
		return 8
	case strings.HasPrefix(u.ID, "research."):
		return 3
	case u.GroupEn == "Buildings", u.ID == "cargo":
		return 8
	}
	return map[string]int{"engine": 18, "veins.rate": 8, "veins.size": 6, "market.contracts": 12, "market.logistics": 8}[u.ID]
}

func TestUpgradeCostsPreserveCreditsAndStayFiniteThrough500(t *testing.T) {
	for _, u := range Catalog {
		previous, qty := 0.0, 0.0
		for level := 0; level <= 500; level++ {
			cost := u.CostAt(level)
			if math.IsInf(cost.Credits, 0) || math.IsNaN(cost.Credits) || cost.Credits < previous {
				t.Fatalf("%s level %d: invalid cost %+v", u.ID, level, cost)
			}
			for good, n := range cost.Goods {
				if math.IsInf(n, 0) || math.IsNaN(n) || n <= 0 || u.rule != nil && n < qty {
					t.Fatalf("%s level %d: invalid %s quantity %g", u.ID, level, good, n)
				}
				qty = n
			}
			if level < max(u.Max, formerCap(u)) && cost.Credits != u.cost(level).Credits {
				t.Fatalf("%s level %d: credit price changed", u.ID, level)
			}
			previous = cost.Credits
		}
	}
}

func TestUpgradeSlicesUseProducibleGoodsAtTheExpectedTier(t *testing.T) {
	raw := map[string]bool{"moss": true, "lichen": true, "glowcap": true, "reed": true, "crystal": true, "iron": true, "geode": true}
	for _, r := range RawV3 {
		raw[r.ID] = true
	}
	var producible func(string, map[string]bool) bool
	producible = func(good string, visiting map[string]bool) bool {
		if raw[good] {
			return true
		}
		if visiting[good] {
			return false
		}
		visiting[good] = true
		defer delete(visiting, good)
		for _, r := range RecipeFor(good) {
			possible := true
			for input := range r.In {
				possible = possible && producible(input, visiting)
			}
			if possible {
				return true
			}
		}
		return false
	}
	for _, u := range Catalog {
		for level := 0; level <= 500; level++ {
			for good := range u.CostAt(level).Goods {
				if !producible(good, map[string]bool{}) {
					t.Fatalf("%s level %d: %s cannot be produced", u.ID, level, good)
				}
				if u.rule != nil && TierOf(good) != min(4, level/5) {
					t.Fatalf("%s level %d asks for tier %d %s", u.ID, level, TierOf(good), good)
				}
			}
		}
	}
	for _, tc := range []struct {
		level int
		good  string
	}{{0, "moss"}, {4, "moss"}, {5, "biomass"}, {9, "biomass"}, {10, "fertilizer"}, {14, "fertilizer"}, {15, "processor"}, {19, "processor"}, {20, "core"}, {500, "core"}} {
		if cost := Find("yield.moss").CostAt(tc.level); len(cost.Goods) != 1 || cost.Goods[tc.good] <= 0 {
			t.Fatalf("moss level %d: %+v", tc.level, cost)
		}
	}
	for _, id := range []string{"yield.moss", "engine", "field"} {
		for level := 0; level <= 2; level++ {
			cost := Find(id).CostAt(level)
			if len(cost.Goods) > 1 {
				t.Fatalf("starter %s needs multiple goods: %+v", id, cost)
			}
			for good := range cost.Goods {
				if !raw[good] {
					t.Fatalf("starter %s needs crafted %s", id, good)
				}
			}
		}
	}
}

func TestOldCapSaveKeepsPossessionsAndCanBuyEveryNextLevel(t *testing.T) {
	s := New(17)
	s.Credits, s.Inv["moss"], s.Inv["core"] = 1234, 45, 2
	s.Blueprints["api.try"], s.Claimed["chapter-1"] = true, true
	s.Contracts["c03-full-sweep"] = &ContractRun{Status: "done", Best: 100}
	s.Buildings = []*Building{{Kind: "composter", Level: 8, In: map[string]float64{"moss": 4}, Out: map[string]float64{"biomass": 2}}}
	for _, u := range Catalog {
		if cap := formerCap(u); cap > 0 {
			s.Levels[u.ID] = cap
		} else {
			s.Levels[u.ID] = u.Max
		}
	}
	s.Normalize()
	blob, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var loaded State
	if err := json.Unmarshal(blob, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.Normalize()
	if !reflect.DeepEqual(s, &loaded) {
		t.Fatal("old save lost possessions on load")
	}
	for _, u := range Catalog {
		if cap := formerCap(u); cap > 0 {
			if u.Max != 0 {
				t.Fatalf("%s still capped at %d", u.ID, u.Max)
			}
			cost := u.CostAt(cap)
			loaded.Credits, loaded.Rank = cost.Credits, 2*(cap/5)
			for good, n := range cost.Goods {
				loaded.Inv[good] = n
			}
			if err := loaded.Buy(u.ID); err != nil || loaded.Level(u.ID) != cap+1 {
				t.Fatalf("%s could not advance past old cap: %v", u.ID, err)
			}
		}
	}
	for _, shop := range loaded.Shop() {
		if formerCap(shop.Upgrade) > 0 && (shop.Max != 0 || strings.Contains(shop.Reason, "All levels")) {
			t.Fatalf("unlimited shop upgrade maxed: %+v", shop)
		}
	}
	for _, id := range []string{"field", "hangar", "drill", "api.scanner", "crop.lichen"} {
		if Find(id).Max <= 0 {
			t.Fatalf("fixed upgrade %s lost cap", id)
		}
	}
	if owned, total := loaded.TotalLevels(); owned <= 0 || total != 0 {
		t.Fatalf("unlimited totals = %d/%d", owned, total)
	}
}
