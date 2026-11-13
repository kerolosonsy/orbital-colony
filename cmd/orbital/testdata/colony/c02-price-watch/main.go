package main

import "orbital/colony"

// Reference for c02-price-watch. Farm like the full sweep, but after every
// harvest sell one unit at a time and only while the price stays at least
// 92% of moss's base price of 1. The rest waits in storage for the market
// to recover.
func main() {
	w, h := colony.WorldSize()
	var farm [8][8]bool
	for pass := 0; ; pass++ {
		for row := 0; row < h; row++ {
			for i := 0; i < w; i++ {
				x, y := i, row
				if row%2 == 1 {
					x = w - 1 - i
				}
				goTo(x, y)
				if pass == 0 {
					farm[x][y] = colony.Till() && colony.Plant(colony.Moss)
				} else if farm[x][y] {
					colony.Harvest()
					colony.Plant(colony.Moss)
					sell()
				}
			}
		}
	}
}

func sell() {
	for colony.Count(colony.Moss) >= 1 && colony.Price(colony.Moss) >= 0.92 {
		colony.Export(colony.Moss, 1)
	}
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
