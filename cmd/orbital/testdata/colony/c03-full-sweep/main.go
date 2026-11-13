package main

import "orbital/colony"

// Reference for c03-full-sweep. The first pass tills and plants every tile
// that takes a seed; later passes only harvest and replant those tiles and
// sell once per pass. Tiles that refused a seed (ice, rock) are skipped for
// good.
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
				}
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
