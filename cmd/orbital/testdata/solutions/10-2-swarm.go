package main

import (
	"orbital/colony"
	"sync"
)

func work(d *colony.Drone, row int) {
	for i := 0; i < row; i++ {
		d.Move(colony.North)
	}
	for x := 0; x < 8; x++ {
		d.Harvest()
		if x < 7 {
			d.Move(colony.East)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	for row := 1; row < 4; row++ {
		d := colony.Spawn()
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(d, row)
		}()
	}
	work(colony.Main(), 0)
	wg.Wait()
}
