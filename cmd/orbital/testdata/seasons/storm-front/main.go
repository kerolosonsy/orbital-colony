package main

import (
	"sync"
	"sync/atomic"

	"orbital/colony"
)

var shelter atomic.Bool

func main() {
	go func() {
		for range colony.StormWarning() {
			shelter.Store(true)
			<-colony.StormClear()
			shelter.Store(false)
		}
	}()
	drones := []*colony.Drone{colony.Main()}
	for len(drones) < 3 {
		d := colony.Spawn()
		if d == nil {
			break
		}
		drones = append(drones, d)
	}
	var workers sync.WaitGroup
	for k, d := range drones {
		workers.Add(1)
		go func() {
			defer workers.Done()
			farm(d, k, len(drones))
		}()
	}
	workers.Wait()
}

func farm(d *colony.Drone, k, n int) {
	w, h := colony.WorldSize()
	for {
		for y := k; y < h; y += n {
			for column := 0; column < w; column++ {
				x := column
				if (y/n)%2 == 1 {
					x = w - 1 - column
				}
				if shelter.Load() {
					goTo(d, 0, 0)
					d.Export(colony.Moss, 0)
					for shelter.Load() {
						d.Wait()
					}
				}
				goTo(d, x, y)
				if d.CanHarvest() {
					d.Harvest()
				}
				if d.GetEntity() == colony.Empty && d.Till() {
					d.Plant(colony.Moss)
				}
				if colony.Count(colony.Moss) >= 20 {
					d.Export(colony.Moss, 0)
				}
			}
		}
	}
}

func goTo(d *colony.Drone, x, y int) {
	for {
		cx, cy := d.Position()
		switch {
		case cx < x:
			d.Move(colony.East)
		case cx > x:
			d.Move(colony.West)
		case cy < y:
			d.Move(colony.North)
		case cy > y:
			d.Move(colony.South)
		default:
			return
		}
	}
}
