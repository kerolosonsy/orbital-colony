package main

import "orbital/colony"

func main() {
	var counts []int
	for z := 1; z < colony.Depth(); z++ {
		grid := colony.LayerGrid(z)
		n := 0
		for x := range grid {
			for y := range grid[x] {
				if g := grid[x][y]; g == colony.Tunnel || g == colony.Shaft {
					n++
				}
			}
		}
		counts = append(counts, n)
	}
	colony.Report(counts)
}
