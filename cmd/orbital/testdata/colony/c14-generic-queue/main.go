package main

import "orbital/colony"

// Reference for c14-generic-queue. The work order is a generic FIFO queue of
// tiles. The upper rows grow glowcap in a cluster, the rest grows moss.
type Queue[T any] struct{ items []T }

func (q *Queue[T]) Push(v T) { q.items = append(q.items, v) }

func (q *Queue[T]) Pop() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	v := q.items[0]
	q.items = q.items[1:]
	return v, true
}

func (q *Queue[T]) Len() int { return len(q.items) }

func main() {
	w, h := colony.WorldSize()
	grid := colony.Grid()
	for {
		var work Queue[colony.Point]
		for y := 0; y < h; y++ {
			for i := 0; i < w; i++ {
				x := i
				if y%2 == 1 {
					x = w - 1 - i
				}
				if grid[x][y] == colony.Regolith || grid[x][y] == colony.Soil {
					work.Push(colony.Point{X: x, Y: y})
				}
			}
		}
		for work.Len() > 0 {
			p, _ := work.Pop()
			crop := colony.Moss
			if p.Y >= h-3 {
				crop = colony.Glowcap
			}
			goTo(p.X, p.Y)
			if colony.CanHarvest() {
				colony.Harvest()
				colony.Plant(crop)
			} else if colony.GetEntity() == colony.Empty && colony.Till() {
				colony.Plant(crop)
			}
		}
		colony.Export(colony.Moss, 0)
		colony.Export(colony.Glowcap, 0)
		goTo(0, 0)
	}
}

func goTo(x, y int) {
	cx, cy := colony.Position()
	for ; cx < x; cx++ {
		colony.Move(colony.East)
	}
	for ; cx > x; cx-- {
		colony.Move(colony.West)
	}
	for ; cy < y; cy++ {
		colony.Move(colony.North)
	}
	for ; cy > y; cy-- {
		colony.Move(colony.South)
	}
}
