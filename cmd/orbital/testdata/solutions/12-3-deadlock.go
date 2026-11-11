package main

import (
	"orbital/colony"
	"sync"
)

var drill, cart sync.Mutex

func harvestWith(d *colony.Drone) {
	drill.Lock()
	cart.Lock()
	d.Wait()
	d.Harvest()
	cart.Unlock()
	drill.Unlock()
}

func work(d *colony.Drone, row int) {
	for i := 0; i < row; i++ {
		d.Move(colony.North)
	}
	for x := 0; x < 6; x++ {
		harvestWith(d)
		if x < 5 {
			d.Move(colony.East)
		}
	}
}

func main() {
	a, b := colony.Main(), colony.Spawn()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); work(a, 0) }()
	go func() { defer wg.Done(); work(b, 1) }()
	wg.Wait()
}
