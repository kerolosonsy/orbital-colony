package main

import "orbital/colony"

func main() {
	var sizes []int
	for x := 0; x < 8; x++ {
		sizes = append(sizes, colony.Measure())
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	colony.Report(sizes)
}
