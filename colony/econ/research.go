package econ

import "math"

func ResearchValue(good string) float64 {
	tier := TierOf(good)
	if tier == 0 {
		return 0
	}
	return math.Pow(5, float64(tier-1))
}

func ExperimentReward(base float64, repeatsToday int) float64 {
	return base * math.Max(0.1, math.Pow(0.8, float64(max(0, repeatsToday))))
}
