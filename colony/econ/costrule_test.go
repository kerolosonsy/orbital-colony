package econ

import (
	"math"
	"reflect"
	"testing"
)

func TestCostSlicesKeepGrowingAfterLastTier(t *testing.T) {
	r := CostRule{Credits: 10, Growth: 1.2, Tiers: [][]string{{"moss"}, {"biomass", "plate"}, {"core"}}, Qty: 2, QtyGrowth: 1.1}
	for _, tc := range []struct {
		level int
		goods []string
	}{{0, []string{"moss"}}, {4, []string{"moss"}}, {5, []string{"biomass", "plate"}}, {9, []string{"biomass", "plate"}}, {10, []string{"core"}}, {500, []string{"core"}}} {
		want := map[string]float64{}
		for _, good := range tc.goods {
			want[good] = 2 * math.Pow(1.1, float64(tc.level))
		}
		got := r.At(tc.level)
		if !reflect.DeepEqual(got.Goods, want) || got.Credits != 10*math.Pow(1.2, float64(tc.level)) {
			t.Errorf("level %d: %+v, want goods %v", tc.level, got, want)
		}
	}
	if got := (CostRule{Credits: 10, Growth: 1}).At(10); got.Credits != 10 || len(got.Goods) != 0 {
		t.Errorf("credit-only rule: %+v", got)
	}
}

func TestTreeCostsStayFiniteAndMonotonicThroughLevel500(t *testing.T) {
	for _, n := range Tree {
		previous := 0.0
		for level := 0; level <= 500; level++ {
			cost := n.Cost.At(level)
			if math.IsInf(cost.Credits, 0) || math.IsNaN(cost.Credits) || cost.Credits < previous {
				t.Fatalf("%s level %d invalid credit cost %g after %g", n.ID, level, cost.Credits, previous)
			}
			for good, qty := range cost.Goods {
				if qty <= 0 || math.IsInf(qty, 0) || math.IsNaN(qty) {
					t.Fatalf("%s level %d invalid %s quantity: %g", n.ID, level, good, qty)
				}
			}
			previous = cost.Credits
		}
	}
}

func TestEffectsDiminishAndEngineHasAFloor(t *testing.T) {
	for _, a := range []float64{0, 0.1, 0.35, 1} {
		previous, increment := 1.0, math.Inf(1)
		for level := 1; level <= 500; level++ {
			mult := EffectMult(a, level)
			if mult < previous || mult > 1+a*float64(level) || mult-previous > increment+1e-12 {
				t.Fatalf("a=%g level=%d: invalid multiplier %g", a, level, mult)
			}
			increment, previous = mult-previous, mult
		}
		if EffectMult(a, 0) != 1 || EffectMult(a, 16) != 1+8*a {
			t.Fatalf("wrong effect formula for a=%g", a)
		}
	}
	previous := 260.0
	for level := 0; level <= 10000; level++ {
		ms := EngineMillis(level)
		if math.IsNaN(ms) || ms < 20 || ms > previous {
			t.Fatalf("level %d: invalid engine time %g", level, ms)
		}
		previous = ms
	}
	if EngineMillis(0) != 260 || EngineMillis(10000) != 20 || EngineMillis(1) != 260/1.09 {
		t.Fatal("wrong engine formula or floor")
	}
}
