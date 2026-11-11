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

func cost(g colony.Ground) int {
	if g == colony.Ice {
		return 4
	}
	return 1
}

func main() {
	grid := colony.Grid()
	start, goal := colony.Point{}, colony.Targets()[0]
	dist := map[colony.Point]int{start: 0}
	parent := map[colony.Point]colony.Point{start: start}
	done := map[colony.Point]bool{}
	for {
		best, found := colony.Point{}, false
		for p, d := range dist {
			if !done[p] && (!found || d < dist[best]) {
				best, found = p, true
			}
		}
		if !found {
			break
		}
		done[best] = true
		for _, n := range neighbors(grid, best) {
			nd := dist[best] + cost(grid[n.X][n.Y])
			if old, ok := dist[n]; !ok || nd < old {
				dist[n], parent[n] = nd, best
			}
		}
	}
	walk(route(parent, start, goal))
	colony.Harvest()
}
