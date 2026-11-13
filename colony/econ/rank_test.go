package econ

import (
	"math"
	"testing"
)

func TestRankGoalsAlternateAndIncreaseThroughRank400(t *testing.T) {
	previous := map[string]float64{}
	for rank := 0; rank <= 400; rank++ {
		goal := GoalFor(rank)
		want := "output"
		if rank%2 == 1 {
			want = "efficiency"
		}
		if goal.Kind != want || goal.Target <= previous[goal.Kind] || math.IsInf(goal.Target, 0) || math.IsNaN(goal.Target) {
			t.Fatalf("rank %d: invalid goal %+v, previous %g", rank, goal, previous[goal.Kind])
		}
		previous[goal.Kind] = goal.Target
	}
	for _, tc := range []struct {
		rank   int
		target float64
	}{{0, 60}, {1, 1}, {2, 90}, {3, 1.15}, {4, 135}} {
		if got := GoalFor(tc.rank).Target; got != tc.target {
			t.Errorf("rank %d target = %g, want %g", tc.rank, got, tc.target)
		}
	}
}

func TestRankMetricsUseBaseValuesAndTickDuration(t *testing.T) {
	sold := map[string]float64{"moss": 10, "plate": 2, "water": 4, "missing": 100}
	for _, tc := range []struct {
		count                  int64
		throughput, efficiency float64
	}{{300, 152, 76.0 / 300}, {600, 76, 76.0 / 600}, {4, 11400, 19}, {0, 0, 0}, {-1, 0, 0}} {
		if got := Throughput(sold, tc.count); math.Abs(got-tc.throughput) > 1e-10 {
			t.Errorf("%d ticks throughput = %g, want %g", tc.count, got, tc.throughput)
		}
		if got := Efficiency(sold, tc.count); math.Abs(got-tc.efficiency) > 1e-10 {
			t.Errorf("%d commands efficiency = %g, want %g", tc.count, got, tc.efficiency)
		}
	}
	if Throughput(nil, Window) != 0 || Efficiency(nil, 1) != 0 {
		t.Fatal("empty sales must have zero metrics")
	}
}

func TestTierAccessRequiresTwoRanksPerFiveLevels(t *testing.T) {
	for _, tc := range []struct {
		rank, level int
		allowed     bool
	}{
		{0, 0, true}, {0, 4, true}, {0, 5, false}, {1, 5, false},
		{2, 5, true}, {2, 9, true}, {3, 10, false}, {4, 10, true},
		{7, 20, false}, {8, 20, true}, {199, 500, false}, {200, 500, true},
	} {
		if got := TierAllowed(tc.rank, tc.level); got != tc.allowed {
			t.Errorf("rank %d level %d allowed=%v, want %v", tc.rank, tc.level, got, tc.allowed)
		}
	}
}
