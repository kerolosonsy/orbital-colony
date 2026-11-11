package main

import "orbital/colony"

func addHarvest(total *int, n int) {
	*total += n
}

func main() {
	total := 0
	for x := 0; x < 8; x++ {
		n, _ := colony.TryHarvest()
		addHarvest(&total, n)
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	colony.Report(total)
}
