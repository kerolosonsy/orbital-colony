package main

import (
	"orbital/colony"
	"slices"
)

func main() {
	sizes := make([]int, 10)
	for x := range sizes {
		sizes[x] = colony.Measure()
		if x < 9 {
			colony.Move(colony.East)
		}
	}
	sorted := slices.Clone(sizes)
	slices.Sort(sorted)
	threshold := sorted[len(sorted)-3]
	lowest := 9
	for x, s := range sizes {
		if s >= threshold {
			lowest = min(lowest, x)
		}
	}
	for x := 9; x >= lowest; x-- {
		if sizes[x] >= threshold {
			colony.Harvest()
		}
		if x > lowest {
			colony.Move(colony.West)
		}
	}
}
