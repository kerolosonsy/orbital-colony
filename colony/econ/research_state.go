package econ

import (
	"fmt"
	"math"

	"orbital/colony/i18n"
)

func (s *State) researchPurchaseReason(id string) string {
	n := FindNode(id)
	if n == nil || s.Level(id) > 0 {
		return ""
	}
	for _, p := range n.After {
		if p.ID != "base" && !s.meets(p.ID, p.Level) {
			parent := FindNode(p.ID)
			return i18n.T("محتاج الأول: ", "Requires first: ") + i18n.T(parent.Name, parent.NameEn)
		}
	}
	if s.RP < n.RP {
		return fmt.Sprintf(i18n.T("ناقصك %.1f RP", "You need %.1f more RP"), n.RP-s.RP)
	}
	return ""
}

// meets reports whether prerequisite id is at level lvl, counting a first
// level lent by the active contract (deep haul lends api.dig for its drill).
func (s *State) meets(id string, lvl int) bool {
	return s.Level(id) >= lvl || lvl <= 1 && s.Lent(id)
}

func (s *State) spendResearch(id string) {
	if n := FindNode(id); n != nil && s.Level(id) == 0 {
		s.RP -= n.RP
	}
}

func (s *State) settleExperiment(c *Contract, r *ContractRun, tick int64) Reward {
	day := tick / 1200
	if r.ResearchDay != day {
		r.ResearchDay, r.RepeatsToday = day, 0
	}
	reward := Reward{RP: ExperimentReward(c.RewardRP, r.RepeatsToday) * s.ResearchGainMult()}
	r.RepeatsToday++
	s.RP += reward.RP
	return reward
}

func (s *State) ResearchYieldMult(good string) float64 {
	return linearMult(0.15, s.Level("research.yield."+good), 3)
}

func (s *State) GrowthResearchMult(good string) float64 {
	return linearMult(0.1, s.Level("research.grow."+good), 3)
}

func (s *State) EngineResearchMult() float64 {
	return linearMult(0.08, s.Level("research.engine"), 3)
}

func (s *State) VeinResearchMult() float64 {
	return linearMult(0.2, s.Level("research.vein-rate"), 3)
}

func (s *State) MarketResearchMult(good string) float64 {
	return linearMult(0.1, s.Level("research.price."+good), 3)
}

func (s *State) RecoveryResearchMult(good string) float64 {
	return linearMult(0.25, s.Level("research.recovery."+good), 3)
}

func (s *State) BuildingSpeedMult(kind string) float64 {
	return linearMult(0.2, s.Level("research.speed."+kind), 3)
}

func (s *State) BuildingCapacity(kind string, level int) float64 {
	return Capacity(level) + (linearMult(4, s.Level("research.buffer."+kind), 3) - 1)
}

func (s *State) ConfigureBuilding(b *Building) {
	b.speedBonus = s.BuildingSpeedMult(b.Kind) - 1
	b.bufferBonus = s.BuildingCapacity(b.Kind, b.Level) - Capacity(b.Level)
}

func (b *Building) capacity() float64 { return Capacity(b.Level) + b.bufferBonus }

func (s *State) BuildingPlacementCost(kind string) *Cost {
	base := PlacementCost(kind)
	if base == nil {
		return nil
	}
	mult := placementMult(s.Level("research.placement"))
	cost := &Cost{Credits: base.Credits * mult, Goods: map[string]float64{}}
	for good, n := range base.Goods {
		cost.Goods[good] = n * mult
	}
	return cost
}

func (s *State) SiftStock() int { return 3 + int(linearMult(1, s.Level("research.sift-stock"), 3)-1) }
func (s *State) SiftRefillTicks() int {
	return int(math.Ceil(100 / linearMult(0.2, s.Level("research.sift-refill"), 3)))
}
func (s *State) MeltRegrowthTicks() int {
	return int(math.Ceil(600 / linearMult(0.2, s.Level("research.melt"), 3)))
}
func (s *State) ResearchGainMult() float64 {
	return linearMult(0.2, s.Level("research.gain"), 3)
}
