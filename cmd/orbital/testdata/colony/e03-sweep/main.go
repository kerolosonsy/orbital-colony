package main

import "orbital/colony"

// The original 3x3 patch is all regolith. Revisit those plots when ready
// instead of spending commands on the larger field's rocks and ice.
func main() {
	var ready [3][3]int
	for pass := 0; ; pass++ {
		for row := 0; row < 3; row++ {
			for i := 0; i < 3; i++ {
				x := i
				if row%2 == 1 {
					x = 2 - i
				}
				goTo(x, row)
				if pass == 0 {
					colony.Till()
				} else {
					for colony.Tick() < ready[x][row] {
						colony.Wait()
					}
					colony.Harvest()
				}
				colony.Plant(colony.Moss)
				ready[x][row] = colony.Tick() + 80
			}
		}
		colony.Export(colony.Moss, 0)
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
