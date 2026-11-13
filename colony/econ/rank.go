package econ

import (
	"math"
	"sort"
)

type RankGoal struct {
	Kind   string
	Target float64
}

const Window = 3000

// Rank 1 starts at 60 value/minute; rank 2 starts at 1 value/command.
func GoalFor(rank int) RankGoal {
	if rank%2 == 0 {
		return RankGoal{"output", 60 * math.Pow(1.5, float64(rank/2))}
	}
	return RankGoal{"efficiency", math.Pow(1.15, float64(rank/2))}
}

func soldValue(sold map[string]float64) float64 {
	goods := make([]string, 0, len(sold))
	for good := range sold {
		goods = append(goods, good)
	}
	sort.Strings(goods)
	total := 0.0
	for _, good := range goods {
		total += sold[good] * PriceV3[good]
	}
	return total
}

func Throughput(sold map[string]float64, ticks int64) float64 {
	if ticks <= 0 {
		return 0
	}
	return soldValue(sold) * 600 / float64(ticks)
}

func Efficiency(sold map[string]float64, commands int64) float64 {
	if commands <= 0 {
		return 0
	}
	return soldValue(sold) / float64(commands)
}

func TierAllowed(rank, level int) bool { return rank >= 2*(level/5) }
