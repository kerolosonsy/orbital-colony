package main

import "orbital/colony"

// Reference for c15-shortest-run. The rock tiles are ordered once into a
// short tour by a recursive nearest-neighbor search. The drone walks the
// tour and, instead of walking on, waits at a vein until it reaches size 4:
// waiting costs no moves.
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
	route := tour(colony.Point{}, rocks)
	for {
		for _, p := range route {
			goTo(p.X, p.Y)
			for i := 0; i < 400; i++ {
				size := colony.Measure()
				if size >= 4 {
					colony.Harvest()
					unload()
					colony.Export(colony.Crystal, 0)
					break
				}
				if size == 0 && i > 30 {
					break // no vein here yet; try the next rock
				}
				colony.Wait()
			}
		}
	}
}

// tour orders left starting from the nearest point to from, recursively.
func tour(from colony.Point, left []colony.Point) []colony.Point {
	if len(left) == 0 {
		return nil
	}
	best := 0
	for i, p := range left {
		if dist(from, p) < dist(from, left[best]) {
			best = i
		}
	}
	next := left[best]
	rest := append(append([]colony.Point{}, left[:best]...), left[best+1:]...)
	return append([]colony.Point{next}, tour(next, rest)...)
}

func dist(a, b colony.Point) int {
	d := 0
	for _, v := range []int{a.X - b.X, a.Y - b.Y} {
		if v < 0 {
			v = -v
		}
		d += v
	}
	return d
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

func unload() {
	x, y := colony.Position()
	goTo(0, 0)
	for !colony.Unload() {
		if len(colony.Carry()) == 0 {
			break
		}
		colony.Repair()
	}
	colony.Export(colony.Moss, 0)
	colony.Export(colony.Glowcap, 0)
	colony.Export(colony.Crystal, 0)
	goTo(x, y)
}
