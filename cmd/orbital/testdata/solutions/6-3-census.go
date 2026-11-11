package main

import "orbital/colony"

func main() {
	counts := map[colony.Entity]int{}
	dir := colony.East
	for y := 0; y < 4; y++ {
		for x := 0; x < 6; x++ {
			if e := colony.GetEntity(); e != colony.Empty {
				counts[e]++
			}
			if x < 5 {
				colony.Move(dir)
			}
		}
		if y < 3 {
			colony.Move(colony.North)
		}
		if dir == colony.East {
			dir = colony.West
		} else {
			dir = colony.East
		}
	}
	colony.Report(counts)
}
