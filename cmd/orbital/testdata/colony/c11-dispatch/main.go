package main

import "orbital/colony"

// Reference for c11-dispatch. Two drones patrol half of the rock tiles each
// and report every harvest on a channel. Drone 0 is the only one that sells:
// it waits on the channel with select and exports in batches.
func main() {
	var rocks []colony.Point
	grid := colony.Grid()
	for x := range grid {
		for y := range grid[x] {
			if grid[x][y] == colony.Rock {
				rocks = append(rocks, colony.Point{X: x, Y: y})
			}
		}
	}
	found := make(chan int)
	for k := 0; k < 2; k++ {
		d := colony.Spawn()
		if d == nil {
			break
		}
		go patrol(d, rocks, k, found)
	}
	pending := 0
	for {
		select {
		case n := <-found:
			pending += n
		default:
			colony.Wait()
		}
		if pending >= 8 {
			colony.Export(colony.Crystal, 0)
			pending = 0
		}
	}
}

func patrol(d *colony.Drone, rocks []colony.Point, k int, found chan<- int) {
	for {
		for i := k; i < len(rocks); i += 2 {
			goTo(d, rocks[i].X, rocks[i].Y)
			if d.Measure() >= 3 {
				n, err := d.TryHarvest()
				if err == nil {
					found <- n
				}
			}
		}
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
