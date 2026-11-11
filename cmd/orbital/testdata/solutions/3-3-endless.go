package main

import "orbital/colony"

func main() {
	dir := colony.East
	for {
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				colony.Harvest()
				colony.Plant(colony.Moss)
				if y < 2 {
					if x%2 == 0 {
						colony.Move(colony.North)
					} else {
						colony.Move(colony.South)
					}
				}
			}
			if !colony.Move(dir) {
				if dir == colony.East {
					dir = colony.West
				} else {
					dir = colony.East
				}
				colony.Move(dir)
			}
		}
	}
}
