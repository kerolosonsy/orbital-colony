package main

import "orbital/colony"

// Reference for c04-right-ground. The ground decides the crop: moss on
// regolith and soil, lichen on ice. A plant is only harvested once the
// scanner says it is ripe, so nothing is ever destroyed.
func main() {
	w, h := colony.WorldSize()
	for {
		for y := 0; y < h; y++ {
			for i := 0; i < w; i++ {
				x := i
				if y%2 == 1 {
					x = w - 1 - i
				}
				goTo(x, y)
				switch colony.GetGround() {
				case colony.Regolith:
					colony.Till()
					colony.Plant(colony.Moss)
				case colony.Soil:
					tend(colony.Moss)
				case colony.Ice:
					tend(colony.Lichen)
				}
			}
		}
		colony.Export(colony.Moss, 0)
		colony.Export(colony.Lichen, 0)
		goTo(0, 0)
	}
}

func tend(crop colony.Entity) {
	if colony.CanHarvest() {
		colony.Harvest()
		colony.Plant(crop)
	} else if colony.GetEntity() == colony.Empty {
		colony.Plant(crop)
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
