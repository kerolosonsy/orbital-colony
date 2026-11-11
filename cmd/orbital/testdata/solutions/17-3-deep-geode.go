package main

import (
	"orbital/colony"
	"slices"
)

type pos struct{ x, y, z int }

var dirs = []struct {
	dx, dy int
	d      colony.Direction
}{{1, 0, colony.East}, {-1, 0, colony.West}, {0, 1, colony.North}, {0, -1, colony.South}}

func main() {
	goal := colony.Targets()[0]
	world := make([][][]colony.Ground, colony.Depth())
	for z := range world {
		world[z] = colony.LayerGrid(z)
	}
	w, h := len(world[0]), len(world[0][0])
	walkable := func(p pos) bool {
		if p.x < 0 || p.y < 0 || p.x >= w || p.y >= h || p.z < 0 || p.z >= len(world) {
			return false
		}
		g := world[p.z][p.x][p.y]
		return p.z == 0 || g == colony.Tunnel || g == colony.Shaft
	}
	start := pos{0, 0, 0}
	parent := map[pos]pos{start: start}
	q := []pos{start}
	var end pos
	var digDir colony.Direction
	found := false
	for len(q) > 0 && !found {
		c := q[0]
		q = q[1:]
		for _, d := range dirs {
			if c.z == 2 && c.x+d.dx == goal.X && c.y+d.dy == goal.Y {
				end, digDir, found = c, d.d, true
				break
			}
		}
		if found {
			break
		}
		var next []pos
		for _, d := range dirs {
			next = append(next, pos{c.x + d.dx, c.y + d.dy, c.z})
		}
		if world[c.z][c.x][c.y] == colony.Shaft {
			next = append(next, pos{c.x, c.y, c.z + 1})
		}
		if c.z > 0 && world[c.z-1][c.x][c.y] == colony.Shaft {
			next = append(next, pos{c.x, c.y, c.z - 1})
		}
		for _, n := range next {
			if _, seen := parent[n]; !seen && walkable(n) {
				parent[n] = c
				q = append(q, n)
			}
		}
	}
	var path []pos
	for p := end; p != start; p = parent[p] {
		path = append(path, p)
	}
	slices.Reverse(path)
	cur := start
	for _, p := range path {
		switch {
		case p.z > cur.z:
			colony.Descend()
		case p.z < cur.z:
			colony.Ascend()
		case p.x > cur.x:
			colony.Move(colony.East)
		case p.x < cur.x:
			colony.Move(colony.West)
		case p.y > cur.y:
			colony.Move(colony.North)
		default:
			colony.Move(colony.South)
		}
		cur = p
	}
	colony.Dig(digDir)
}
