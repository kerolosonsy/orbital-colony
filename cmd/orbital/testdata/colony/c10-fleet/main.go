package main

import (
	"sync"

	"orbital/colony"
)

// Reference for c10-fleet. Three drones share the field: drone k farms every
// row y with y%3 == k, each in its own goroutine.
func main() {
	drones := []*colony.Drone{colony.Main()}
	for len(drones) < 3 {
		d := colony.Spawn()
		if d == nil {
			break
		}
		drones = append(drones, d)
	}
	var wg sync.WaitGroup
	for k, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			farm(d, k, len(drones))
		}()
	}
	wg.Wait()
}

func farm(d *colony.Drone, k, n int) {
	w, h := colony.WorldSize()
	for {
		for y := k; y < h; y += n {
			for i := 0; i < w; i++ {
				x := i
				if (y/n)%2 == 1 {
					x = w - 1 - i
				}
				goTo(d, x, y)
				if d.CanHarvest() {
					d.Harvest()
					d.Plant(colony.Moss)
				} else if d.GetEntity() == colony.Empty && d.Till() {
					d.Plant(colony.Moss)
				}
			}
		}
		d.Export(colony.Moss, 0)
	}
}

func goTo(d *colony.Drone, x, y int) {
	cx, cy := d.Position()
	for ; cx < x; cx++ {
		d.Move(colony.East)
	}
	for ; cx > x; cx-- {
		d.Move(colony.West)
	}
	for ; cy < y; cy++ {
		d.Move(colony.North)
	}
	for ; cy > y; cy-- {
		d.Move(colony.South)
	}
}
