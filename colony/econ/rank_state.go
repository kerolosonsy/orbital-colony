package econ

type RankProgress struct {
	Level  int     `json:"level"`
	Kind   string  `json:"kind"`
	Target float64 `json:"target"`
	Value  float64 `json:"value"`
	Filled float64 `json:"filled"`
}

// RankAllows checks the cost slice of the next purchase, leaving owned levels intact.
func (s *State) RankAllows(id string) (ok bool, need int) {
	level := s.Levels[id]
	return TierAllowed(s.Rank, level), 2 * (level / 5)
}

// RecordRankSale removes the sale-time yield multiplier because stored goods
// do not retain the upgrade level at which they were harvested.
func (s *State) RecordRankSale(good string, amount float64) {
	if amount > 0 && PriceV3[good] > 0 {
		s.RankWindow.RecordSale(good, amount/s.YieldMult(good))
	}
}

func (s *State) RankProgress() RankProgress {
	goal := GoalFor(s.Rank)
	sold, commands := s.RankWindow.totals()
	current := Throughput(sold, min(s.RankWindow.Ticks, Window))
	if goal.Kind == "efficiency" {
		current = Efficiency(sold, commands)
	}
	return RankProgress{Level: s.Rank, Kind: goal.Kind, Target: goal.Target, Value: current, Filled: float64(min(s.RankWindow.Ticks, Window)) / Window}
}

func (s *State) PromoteRank() bool {
	progress := s.RankProgress()
	if progress.Filled < 1 || progress.Value < progress.Target {
		return false
	}
	if progress.Kind == "efficiency" {
		_, commands := s.RankWindow.totals()
		if commands < RankMinCommands {
			return false
		}
	}
	s.Rank++
	return true
}
