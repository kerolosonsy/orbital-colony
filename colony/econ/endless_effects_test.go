package econ

import (
	"math"
	"strings"
	"testing"
)

func TestEveryFormerLevelKeepsItsEffect(t *testing.T) {
	for _, u := range Catalog {
		cap := formerCap(u)
		if cap == 0 {
			continue
		}
		for level := 0; level <= cap; level++ {
			s := New(1)
			s.Levels[u.ID] = level
			var got, want float64
			switch {
			case strings.HasPrefix(u.ID, "yield."):
				got, want = s.YieldMult(strings.TrimPrefix(u.ID, "yield.")), 1+0.35*float64(level)
			case strings.HasPrefix(u.ID, "grow."):
				good := strings.TrimPrefix(u.ID, "grow.")
				got, want = float64(s.GrowTime(good)), math.Trunc(math.Max(10, float64(GrowTicks[good])*math.Pow(0.9, float64(level))))
			case u.GroupEn == "Buildings":
				got, want = s.BuildingCapacity(u.ID, level), 8*EffectMult(0.25, level)
				b := &Building{Kind: u.ID, Level: level}
				if ticks := b.cycleTicks(Recipe{Ticks: 100}); ticks != int(math.Ceil(100/EffectMult(0.25, level))) {
					t.Fatalf("%s level %d: production speed changed", u.ID, level)
				}
			case u.ID == "engine":
				got, want = s.ActionMillis(), 260*math.Pow(0.92, float64(level))
			case u.ID == "cargo":
				got, want = float64(s.CargoCapacity()), 10+5*float64(level)
			case u.ID == "veins.rate":
				got, want = float64(s.VeinEvery()), math.Trunc(160*math.Pow(0.8, float64(level)))
			case u.ID == "veins.size":
				got, want = float64(s.VeinMax()), 4+float64(level)
			case u.ID == "market.contracts":
				got, want = s.Price("moss"), 1+0.1*float64(level)
			case u.ID == "market.logistics":
				s.Demand["moss"] = 0.5
				s.Recover(1)
				got, want = s.Demand["moss"], 0.5+0.5*0.0025*(1+0.35*float64(level))
			default:
				got, want = researchEffect(s, u.ID, level)
			}
			if math.IsNaN(got) || math.Abs(got-want) > 1e-9 {
				t.Fatalf("%s level %d: effect %g, previously %g", u.ID, level, got, want)
			}
		}
	}
}

func researchEffect(s *State, id string, level int) (got, want float64) {
	l := float64(level)
	suffix := strings.TrimPrefix(id, "research.")
	family, good, _ := strings.Cut(suffix, ".")
	switch family {
	case "yield":
		return s.YieldMult(good), 1 + 0.15*l
	case "grow":
		return float64(s.GrowTime(good)), math.Trunc(math.Max(10, float64(GrowTicks[good])/(1+0.1*l)))
	case "price":
		return s.Price(good), PriceV3[good] * (1 + 0.1*l)
	case "recovery":
		s.Demand[good] = 0.5
		s.Recover(1)
		return s.Demand[good], 0.5 + 0.5*0.0025*(1+0.25*l)
	case "speed":
		return s.BuildingSpeedMult(good), 1 + 0.2*l
	case "buffer":
		return s.BuildingCapacity(good, 0), 8 + 4*l
	case "engine":
		return s.ActionMillis(), 260 / (1 + 0.08*l)
	case "vein-rate":
		return float64(s.VeinEvery()), math.Trunc(160 / (1 + 0.2*l))
	case "vein-size":
		return float64(s.VeinMax()), 4 + l
	case "sift-stock":
		return float64(s.SiftStock()), 3 + l
	case "sift-refill":
		return float64(s.SiftRefillTicks()), math.Ceil(100 / (1 + 0.2*l))
	case "melt":
		return float64(s.MeltRegrowthTicks()), math.Ceil(600 / (1 + 0.2*l))
	case "placement":
		return s.BuildingPlacementCost("composter").Credits, 25 * (1 - 0.1*l)
	case "gain":
		return s.ResearchGainMult(), 1 + 0.2*l
	}
	return math.NaN(), 0
}

func TestEndlessEffectsRemainMonotonicAndSafe(t *testing.T) {
	previousEngine, previousGrowth, previousVeins, previousPlacement := 260.0, int64(GrowTicks["moss"]), int64(160), 25.0
	for level := 0; level <= 500; level++ {
		s := New(1)
		for _, u := range Catalog {
			if u.Max == 0 {
				s.Levels[u.ID] = level
			}
		}
		engine, growth, veins := s.ActionMillis(), s.GrowTime("moss"), s.VeinEvery()
		placement := s.BuildingPlacementCost("composter")
		if engine < 20 || engine > previousEngine || growth < 10 || growth > previousGrowth || veins < 1 || veins > previousVeins || placement.Credits <= 0 || placement.Credits > previousPlacement || placement.Goods["moss"] <= 0 {
			t.Fatalf("level %d: engine=%g growth=%d veins=%d placement=%+v", level, engine, growth, veins, placement)
		}
		s.Demand["moss"] = 0.5
		s.Recover(100)
		if price := s.Price("moss"); math.IsNaN(price) || math.IsInf(price, 0) || price < 0 || s.Demand["moss"] > 1 || s.Demand["moss"] < 0.5 {
			t.Fatalf("level %d: unstable market price=%g demand=%g", level, price, s.Demand["moss"])
		}
		previousEngine, previousGrowth, previousVeins, previousPlacement = engine, growth, veins, placement.Credits
	}
}
