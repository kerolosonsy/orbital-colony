package main

import (
	"orbital/colony"
	"slices"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

func dist(a, b colony.Point) int {
	dx, dy := a.X-b.X, a.Y-b.Y
	return max(dx, -dx) + max(dy, -dy)
}

func collect(d *colony.Drone, left []colony.Point) {
	for len(left) > 0 {
		x, y := d.Position()
		here := colony.Point{X: x, Y: y}
		best := 0
		for i, t := range left {
			if dist(here, t) < dist(here, left[best]) {
				best = i
			}
		}
		goTo(d, left[best])
		d.Harvest()
		left = append(left[:best], left[best+1:]...)
	}
}

func main() {
	targets := colony.Targets()
	slices.SortFunc(targets, func(a, b colony.Point) int { return a.X - b.X })
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	var wg sync.WaitGroup
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			collect(d, slices.Clone(targets[i*5:i*5+5]))
		}()
	}
	wg.Wait()
}
