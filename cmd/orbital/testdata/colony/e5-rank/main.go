package main

import "orbital/colony"

func main() {
	// A closed route through a 4x4 field avoids the return trip and repeated tilling.
	route := []colony.Direction{
		colony.East, colony.East, colony.East, colony.North,
		colony.West, colony.West, colony.North, colony.East,
		colony.East, colony.North, colony.West, colony.West,
		colony.West, colony.South, colony.South, colony.South,
	}
	for {
		for _, direction := range route {
			colony.Harvest()
			colony.Plant(colony.Moss)
			colony.Move(direction)
		}
		colony.Export(colony.Moss, 0)
		kind, target, current := colony.RankGoal()
		colony.Report([]any{colony.Rank(), kind, target, current})
	}
}
