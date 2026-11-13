package main

import "orbital/colony"

// Reference for c06-ripe-ledger. One pass plants every tile that takes moss
// and writes down when it ripens. After that the drone only visits tiles the
// ledger says are ripe, in the order they were planted.
func main() {
	const grow = 80
	w, h := colony.WorldSize()
	var farm []colony.Point
	ripe := map[colony.Point]int{}
	for y := 0; y < h; y++ {
		for i := 0; i < w; i++ {
			x := i
			if y%2 == 1 {
				x = w - 1 - i
			}
			goTo(x, y)
			if colony.Till() && colony.Plant(colony.Moss) {
				p := colony.Point{X: x, Y: y}
				farm = append(farm, p)
				ripe[p] = colony.Tick() + grow
			}
		}
	}
	for {
		for _, p := range farm {
			for colony.Tick() < ripe[p] {
				colony.Wait()
			}
			goTo(p.X, p.Y)
			colony.Harvest()
			colony.Plant(colony.Moss)
			ripe[p] = colony.Tick() + grow
		}
		colony.Export(colony.Moss, 0)
	}
}

func goTo(x, y int) {
	cx, cy := colony.Position()
	for ; cx < x; cx++ {
		for !colony.Move(colony.East) {
			colony.Repair()
		}
	}
	for ; cx > x; cx-- {
		for !colony.Move(colony.West) {
			colony.Repair()
		}
	}
	for ; cy < y; cy++ {
		for !colony.Move(colony.North) {
			colony.Repair()
		}
	}
	for ; cy > y; cy-- {
		for !colony.Move(colony.South) {
			colony.Repair()
		}
	}
	if colony.Fault() != nil {
		colony.Repair()
	}
}
