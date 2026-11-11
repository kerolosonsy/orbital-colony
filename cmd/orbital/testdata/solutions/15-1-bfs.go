package main

import (
	"orbital/colony"
	"slices"
)

func walk(path []colony.Point) {
	x, y := colony.Position()
	for _, p := range path {
		switch {
		case p.X > x:
			colony.Move(colony.East)
		case p.X < x:
			colony.Move(colony.West)
		case p.Y > y:
			colony.Move(colony.North)
		default:
			colony.Move(colony.South)
		}
		x, y = p.X, p.Y
	}
}

func neighbors(grid [][]colony.Ground, p colony.Point) []colony.Point {
	var out []colony.Point
	for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		q := colony.Point{X: p.X + d[0], Y: p.Y + d[1]}
		if q.X >= 0 && q.Y >= 0 && q.X < len(grid) && q.Y < len(grid[0]) && grid[q.X][q.Y] != colony.Rock {
			out = append(out, q)
		}
	}
	return out
}

func route(parent map[colony.Point]colony.Point, start, goal colony.Point) []colony.Point {
	var path []colony.Point
	for p := goal; p != start; p = parent[p] {
		path = append(path, p)
	}
	slices.Reverse(path)
	return path
}

func main() {
	grid := colony.Grid()
	start, goal := colony.Point{}, colony.Targets()[0]
	parent := map[colony.Point]colony.Point{start: start}
	q := []colony.Point{start}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		for _, n := range neighbors(grid, p) {
			if _, ok := parent[n]; !ok {
				parent[n] = p
				q = append(q, n)
			}
		}
	}
	walk(route(parent, start, goal))
	colony.Harvest()
}
