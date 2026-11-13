package main

import "orbital/colony"

// Reference for c08-big-veins. Patrol the rock tiles and harvest a vein only
// once the gauge says it reached size 4.
func main() {
	pos := colony.Point{}
	for {
		if p := bigVein(&pos); p != nil {
			colony.Harvest()
			colony.Export(colony.Crystal, 0)
		}
	}
}

// bigVein walks the rocks from pos and returns the first big vein, or nil
// after a full round without one. pos follows the drone.
func bigVein(pos *colony.Point) *colony.Point {
	grid := colony.Grid()
	for x := range grid {
		for y := range grid[x] {
			if grid[x][y] != colony.Rock {
				continue
			}
			goTo(x, y)
			pos.X, pos.Y = x, y
			if colony.Measure() >= 4 {
				return pos
			}
		}
	}
	return nil
}

func goTo(x, y int) {
	cx, cy := colony.Position()
	for ; cx < x; cx++ {
		colony.Move(colony.East)
	}
	for ; cx > x; cx-- {
		colony.Move(colony.West)
	}
	for ; cy < y; cy++ {
		colony.Move(colony.North)
	}
	for ; cy > y; cy-- {
		colony.Move(colony.South)
	}
}
