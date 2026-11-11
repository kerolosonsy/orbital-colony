package main

import "orbital/colony"

func depthIndex(depths []int, target int) int {
	lo, hi := 0, len(depths)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if depths[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	if lo < len(depths) && depths[lo] == target {
		return lo
	}
	return -1
}
func flood(grid [][]colony.Ground, x, y int, seen map[colony.Point]bool) int {
	if x < 0 || x >= len(grid) || y < 0 || y >= len(grid[x]) || grid[x][y] == colony.Rock {
		return 0
	}
	p := colony.Point{X: x, Y: y}
	if seen[p] {
		return 0
	}
	seen[p] = true
	return 1 + flood(grid, x+1, y, seen) + flood(grid, x-1, y, seen) + flood(grid, x, y+1, seen) + flood(grid, x, y-1, seen)
}
func main() {
	index := depthIndex(colony.Depths(), colony.Reading(0))
	count := flood(colony.Grid(), 0, 0, map[colony.Point]bool{})
	colony.Report([]int{index, count})
}
