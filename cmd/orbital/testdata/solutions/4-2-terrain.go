package main

import "orbital/colony"

func main() {
	dir := colony.East
	for y := 0; y < 3; y++ {
		for x := 0; x < 6; x++ {
			switch colony.GetGround() {
			case colony.Soil:
				colony.Plant(colony.Moss)
			case colony.Ice:
				colony.Plant(colony.Lichen)
			case colony.Regolith:
				colony.Till()
				colony.Plant(colony.Moss)
			}
			if x < 5 {
				colony.Move(dir)
			}
		}
		colony.Move(colony.North)
		if dir == colony.East {
			dir = colony.West
		} else {
			dir = colony.East
		}
	}
}
