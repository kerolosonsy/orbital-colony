package main

import (
	"orbital/colony"
	"player/plan"
)

func main() {
	ready := []int{0}
	colony.Till()
	colony.Plant(colony.Moss)
	for {
		next := plan.Next(ready)
		for colony.Tick() < ready[next] || !colony.CanHarvest() {
			if colony.Fault() != nil {
				colony.Repair()
			}
			// A storm can destroy the unripe plant; start the plot again.
			if colony.GetEntity() != colony.Moss {
				colony.Plant(colony.Moss)
				ready[next] = colony.Tick() + 80
			}
			colony.Wait()
		}
		colony.Harvest()
		colony.Unload() // cargo is on after the chapter 13 contracts
		colony.Export(colony.Moss, 0)
		colony.Plant(colony.Moss)
		ready[next] = colony.Tick() + 80
	}
}
