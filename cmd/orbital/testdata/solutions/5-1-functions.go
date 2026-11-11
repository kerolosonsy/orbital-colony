package main

import "orbital/colony"

func harvestRow(n int, dir colony.Direction) {
	for i := 0; i < n; i++ {
		colony.Harvest()
		if i < n-1 {
			colony.Move(dir)
		}
	}
}

func up(n int) {
	for i := 0; i < n; i++ {
		colony.Move(colony.North)
	}
}

func main() {
	harvestRow(8, colony.East)
	up(2)
	harvestRow(8, colony.West)
	up(2)
	harvestRow(8, colony.East)
}
