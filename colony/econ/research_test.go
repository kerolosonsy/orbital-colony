package econ

import (
	"math"
	"testing"
)

func TestResearchRewardsGrowSteeplyWithCraftedTier(t *testing.T) {
	for _, good := range Goods {
		if got := ResearchValue(good); got != 0 {
			t.Errorf("raw %s yields %g RP", good, got)
		}
	}
	for _, good := range RawV3 {
		if got := ResearchValue(good.ID); got != 0 {
			t.Errorf("raw %s yields %g RP", good.ID, got)
		}
	}
	for _, tc := range []struct {
		good string
		rp   float64
	}{{"biomass", 1}, {"plate", 1}, {"circuit", 5}, {"processor", 25}, {"core", 125}, {"missing", 0}} {
		if got := ResearchValue(tc.good); got != tc.rp {
			t.Errorf("%s yields %g RP, want %g", tc.good, got, tc.rp)
		}
	}
}

func TestExperimentRepeatsCompoundDiscountAndReachFloor(t *testing.T) {
	for _, tc := range []struct {
		base    float64
		repeats int
		want    float64
	}{
		{100, -1, 100}, {100, 0, 100}, {100, 1, 80}, {100, 2, 64},
		{100, 10, 10.73741824}, {100, 11, 10}, {100, 10000, 10}, {50, 2, 32}, {0, 2, 0},
	} {
		if got := ExperimentReward(tc.base, tc.repeats); math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("base %g repeat %d reward = %g, want %g", tc.base, tc.repeats, got, tc.want)
		}
	}
}
