package main

import (
	"orbital/colony"
	"sync"
)

func main() {
	total := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn(), colony.Spawn()}
	for row, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < row; i++ {
				d.Move(colony.North)
			}
			for x := 0; x < 8; x++ {
				n, _ := d.TryHarvest()
				mu.Lock()
				total += n
				mu.Unlock()
				if x < 7 {
					d.Move(colony.East)
				}
			}
		}()
	}
	wg.Wait()
	colony.Report(total)
}
