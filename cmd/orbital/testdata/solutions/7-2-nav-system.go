package main

import "orbital/colony"

type Nav struct {
	X, Y int
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

func (n *Nav) Collect(p colony.Point) {
	n.GoTo(p)
	colony.Harvest()
}

func main() {
	nav := &Nav{}
	for _, t := range colony.Targets() {
		nav.Collect(t)
	}
}
