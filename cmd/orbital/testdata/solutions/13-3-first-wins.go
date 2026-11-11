package main

import (
	"context"
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
}

func search(ctx context.Context, d *colony.Drone, row int) bool {
	goTo(d, colony.Point{X: 0, Y: row})
	for y := row; y < row+2; y++ {
		for k := 0; k < 10; k++ {
			if ctx.Err() != nil {
				return false
			}
			if d.Measure() == 10 {
				d.Harvest()
				return true
			}
			if k < 9 {
				if y%2 == 0 {
					d.Move(colony.East)
				} else {
					d.Move(colony.West)
				}
			}
		}
		if y == row {
			d.Move(colony.North)
		}
	}
	return false
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if search(ctx, d, i*2) {
				cancel()
			}
			goTo(d, colony.Point{})
		}()
	}
	wg.Wait()
}
