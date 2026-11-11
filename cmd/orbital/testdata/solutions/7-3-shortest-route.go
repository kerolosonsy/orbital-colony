package main

import "orbital/colony"

type Nav struct {
	X, Y int
}

func (n *Nav) Dist(p colony.Point) int {
	return abs(n.X-p.X) + abs(n.Y-p.Y)
}

func (n *Nav) GoTo(p colony.Point) {
	for ; n.X < p.X; n.X++ {
		colony.Move(colony.East)
	}
	for ; n.X > p.X; n.X-- {
		colony.Move(colony.West)
	}
	for ; n.Y < p.Y; n.Y++ {
		colony.Move(colony.North)
	}
	for ; n.Y > p.Y; n.Y-- {
		colony.Move(colony.South)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func main() {
	nav := &Nav{}
	left := colony.Targets()
	for len(left) > 0 {
		best := 0
		for i, t := range left {
			if nav.Dist(t) < nav.Dist(left[best]) {
				best = i
			}
		}
		nav.GoTo(left[best])
		colony.Harvest()
		left = append(left[:best], left[best+1:]...)
	}
}
