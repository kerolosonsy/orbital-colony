package main

import "orbital/colony"

func main() {
	for {
		if colony.CanHarvest() {
			colony.Harvest()
			colony.Plant(colony.Moss)
		}
		x, y := colony.Position()
		if y == 0 && x < 4 {
			colony.Move(colony.East)
		} else if y == 0 {
			colony.Move(colony.North)
		} else if x > 0 {
			colony.Move(colony.West)
		} else {
			colony.Move(colony.South)
		}
	}
}
