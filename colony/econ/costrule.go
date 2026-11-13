package econ

import "math"

func (r CostRule) At(level int) Cost {
	cost := Cost{Credits: r.Credits * math.Pow(r.Growth, float64(level))}
	if len(r.Tiers) == 0 {
		return cost
	}
	tier := min(level/5, len(r.Tiers)-1)
	cost.Goods = make(map[string]float64, len(r.Tiers[tier]))
	for _, good := range r.Tiers[tier] {
		cost.Goods[good] = r.Qty * math.Pow(r.QtyGrowth, float64(level))
	}
	return cost
}

type CostRule struct {
	Credits, Growth float64
	Tiers           [][]string
	Qty, QtyGrowth  float64
}

func EffectMult(a float64, level int) float64 {
	return 1 + a*math.Pow(float64(level), 0.75)
}

func EngineMillis(level int) float64 {
	return math.Max(20, 260/(1+0.09*math.Pow(float64(level), 0.7)))
}
