package main

import "orbital/colony"

func goTo(p colony.Point) {
	x, y := colony.Position()
	for ; x < p.X; x++ {
		colony.Move(colony.East)
	}
	for ; x > p.X; x-- {
		colony.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		colony.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		colony.Move(colony.South)
	}
}

func main() {
	for p := colony.Nearest(); p != nil; p = colony.Nearest() {
		goTo(*p)
		colony.Harvest()
	}
}
