package main

import (
	"orbital/colony"
	"slices"
)

func main() {
	width, _ := colony.WorldSize()
	sizes := make([]int, width)
	positions := make([]int, width)
	for x := 0; x < width; x++ {
		sizes[x] = colony.Measure()
		positions[x] = x
		if x+1 < width {
			colony.Move(colony.East)
		}
	}
	slices.SortFunc(positions, func(a, b int) int {
		if sizes[a] > sizes[b] {
			return -1
		}
		if sizes[a] < sizes[b] {
			return 1
		}
		return a - b
	})
	colony.Report(positions)
}
