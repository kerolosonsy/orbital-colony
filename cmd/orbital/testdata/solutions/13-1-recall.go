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

func work(ctx context.Context, d *colony.Drone, rows []int) {
	defer goTo(d, colony.Point{})
	for _, row := range rows {
		for k := 0; k < 9; k++ {
			x := k
			if row%2 == 1 {
				x = 8 - k
			}
			if ctx.Err() != nil {
				return
			}
			goTo(d, colony.Point{X: x, Y: row})
			if ctx.Err() != nil {
				return
			}
			d.Harvest()
		}
	}
}

func main() {
	ctx := colony.Mission()
	var wg sync.WaitGroup
	plan := [][]int{{0, 1}, {2, 3}, {4}}
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for i, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			work(ctx, d, plan[i])
		}()
	}
	wg.Wait()
	<-ctx.Done()
	select {}
}
