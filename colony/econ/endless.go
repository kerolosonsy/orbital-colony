package econ

import (
	"fmt"
	"math"
	"strings"
)

// The former caps are compatibility boundaries, not purchase limits.
func configureEndless(u *Upgrade) {
	growth := 0.0
	switch {
	case strings.HasPrefix(u.ID, "yield."):
		growth = 2.05
	case strings.HasPrefix(u.ID, "grow."):
		growth = 1.95
	case strings.HasPrefix(u.ID, "research."):
		growth = 1.8
	case u.GroupEn == "Buildings", u.ID == "cargo":
		growth = 1.6
	default:
		growth = map[string]float64{"engine": 1.62, "field": 2.25, "veins.rate": 1.9, "veins.size": 2.1, "market.contracts": 2, "market.logistics": 1.95}[u.ID]
	}
	if growth == 0 {
		return
	}
	u.rule = &CostRule{Credits: u.cost(0).Credits, Growth: growth, Tiers: upgradeTiers(u.ID), Qty: 5, QtyGrowth: 1.15}
	if u.ID == "veins.size" {
		u.rule.Qty, u.rule.QtyGrowth = 20, 1.7
	}
	if u.ID == "cargo" {
		u.rule.Qty = 10
	}
	if u.ID != "field" {
		cap := u.Max
		u.Max = 0
		u.Desc += fmt.Sprintf(" بعد المستوى %d العائد بيتناقص.", cap)
		u.DescEn += fmt.Sprintf(" After level %d, gains diminish.", cap)
	}
}

func upgradeTiers(id string) [][]string {
	raw := "moss"
	for _, prefix := range []string{"yield.", "grow.", "research.yield.", "research.grow.", "research.price.", "research.recovery."} {
		if good, ok := strings.CutPrefix(id, prefix); ok {
			raw = good
		}
	}
	kind := strings.TrimPrefix(strings.TrimPrefix(id, "research.speed."), "research.buffer.")
	if good := map[string]string{"engine": "moss", "cargo": "moss", "composter": "moss", "silo": "moss", "smelter": "iron", "kiln": "silica", "refinery": "water", "assembler": "copper", "reactor-lab": "uranium", "generator": "silica", "collector": "silica", "lab": "crystal", "veins.rate": "crystal", "veins.size": "crystal", "research.vein-rate": "crystal", "research.vein-size": "crystal", "research.sift-stock": "silica", "research.sift-refill": "silica", "research.melt": "water"}[kind]; good != "" {
		raw = good
	}
	t1, t2, t3 := "biomass", "fertilizer", "processor"
	switch raw {
	case "iron", "titanium", "uranium":
		t1, t2, t3 = "plate", "alloy", "frame"
	case "copper":
		t1, t2 = "wire", "circuit"
	case "crystal", "geode", "silica", "helium":
		t1, t2 = "glass", "circuit"
	case "water", "sulfur":
		t1, t2, t3 = "fuel", "coolant", "cell"
	}
	if id == "engine" || id == "cargo" || id == "research.engine" {
		t1, t2, t3 = "plate", "battery", "frame"
	}
	return [][]string{{raw}, {t1}, {t2}, {t3}, {"core"}}
}

// Keep purchased effects exact through their old cap, then join the power curve
// at that endpoint. No save migration or retroactive weakening is needed.
func continuedMult(level, cap int, endpoint float64) float64 {
	return EffectMult((endpoint-1)/math.Pow(float64(cap), 0.75), level)
}

func linearMult(step float64, level, cap int) float64 {
	if level <= cap {
		return 1 + step*float64(level)
	}
	return continuedMult(level, cap, 1+step*float64(cap))
}

func timeScale(factor float64, level, cap int) float64 {
	if level <= cap {
		return math.Pow(factor, float64(level))
	}
	return 1 / continuedMult(level, cap, math.Pow(factor, -float64(cap)))
}

func engineMillis(level int) float64 {
	if level <= 18 {
		return 260 * math.Pow(0.92, float64(level))
	}
	return math.Max(20, 260*math.Pow(0.92, 18)*EngineMillis(level)/EngineMillis(18))
}

func placementMult(level int) float64 {
	if level <= 3 {
		return 1 - 0.1*float64(level)
	}
	return 0.7 / EffectMult(0.1, level-3)
}
