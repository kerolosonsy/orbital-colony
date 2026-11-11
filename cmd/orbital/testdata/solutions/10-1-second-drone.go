package main

import (
	"orbital/colony"
	"sync"
)

func harvestRow(d *colony.Drone, dir colony.Direction) {
	for x := 0; x < 10; x++ {
		d.Harvest()
		if x < 9 {
			d.Move(dir)
		}
	}
}

func main() {
	helper := colony.Spawn()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		helper.Move(colony.North)
		helper.Move(colony.North)
		harvestRow(helper, colony.East)
	}()
	harvestRow(colony.Main(), colony.East)
	wg.Wait()
}
