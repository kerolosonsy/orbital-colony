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
			colony.Wait()
		}
		colony.Harvest()
		colony.Export(colony.Moss, 0)
		colony.Plant(colony.Moss)
		ready[next] = colony.Tick() + 80
	}
}
