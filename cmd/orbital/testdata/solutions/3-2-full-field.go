package main

import "orbital/colony"

func main() {
	dir := colony.East
	for y := 0; y < 6; y++ {
		for x := 0; x < 6; x++ {
			colony.Harvest()
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
