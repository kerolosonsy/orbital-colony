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
		if p.ID != "base" && s.Level(p.ID) < p.Level {
			parent := FindNode(p.ID)
			return i18n.T("محتاج الأول: ", "Requires first: ") + i18n.T(parent.Name, parent.NameEn)
		}
	}
	if s.RP < n.RP {
		return fmt.Sprintf(i18n.T("ناقصك %.1f RP", "You need %.1f more RP"), n.RP-s.RP)
	}
	return ""
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
	return 1 + 0.15*float64(s.Level("research.yield."+good))
}

func (s *State) GrowthResearchMult(good string) float64 {
	return 1 + 0.1*float64(s.Level("research.grow."+good))
}

func (s *State) EngineResearchMult() float64 {
	return 1 + 0.08*float64(s.Level("research.engine"))
}

func (s *State) VeinResearchMult() float64 {
	return 1 + 0.2*float64(s.Level("research.vein-rate"))
}

func (s *State) MarketResearchMult(good string) float64 {
	return 1 + 0.1*float64(s.Level("research.price."+good))
}

func (s *State) RecoveryResearchMult(good string) float64 {
	return 1 + 0.25*float64(s.Level("research.recovery."+good))
}

func (s *State) BuildingSpeedMult(kind string) float64 {
	return 1 + 0.2*float64(s.Level("research.speed."+kind))
}

func (s *State) BuildingCapacity(kind string, level int) float64 {
	return Capacity(level) + 4*float64(s.Level("research.buffer."+kind))
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
	mult := 1 - 0.1*float64(s.Level("research.placement"))
	cost := &Cost{Credits: base.Credits * mult, Goods: map[string]float64{}}
	for good, n := range base.Goods {
		cost.Goods[good] = n * mult
	}
	return cost
}

func (s *State) SiftStock() int { return 3 + s.Level("research.sift-stock") }
func (s *State) SiftRefillTicks() int {
	return int(math.Ceil(100 / (1 + 0.2*float64(s.Level("research.sift-refill")))))
}
func (s *State) MeltRegrowthTicks() int {
	return int(math.Ceil(600 / (1 + 0.2*float64(s.Level("research.melt")))))
}
func (s *State) ResearchGainMult() float64 {
	return 1 + 0.2*float64(s.Level("research.gain"))
}
