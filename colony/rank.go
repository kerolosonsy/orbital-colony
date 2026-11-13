package colony

// Rank is the colony's earned rank. Querying it costs no ticks.
func Rank() int { return world.Rank() }

// RankGoal returns the next goal ("output" or "efficiency"), its target,
// and the current measurement. Querying it costs no ticks.
func RankGoal() (kind string, target, current float64) { return world.RankGoal() }
