package main

import (
	"context"

	"orbital/colony"
)

// Reference for c13-reed-clock. Reed grows only on soil next to ice and rots
// 90 ticks after it ripens. Every wait for a reed has a deadline: if it
// passes, the drone moves on instead of standing there forever.
func main() {
	var beds []colony.Point
	grid := colony.Grid()
	w, h := colony.WorldSize()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if grid[x][y] == colony.Regolith && nearIce(grid, x, y) {
				beds = append(beds, colony.Point{X: x, Y: y})
			}
		}
	}
	ready := map[colony.Point]int{}
	for _, b := range beds {
		goTo(b.X, b.Y)
		if colony.Till() && colony.Plant(colony.Reed) {
			ready[b] = colony.Tick() + 220
		}
	}
	for {
		for _, b := range beds {
			if _, ok := ready[b]; !ok {
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), colony.Ticks(250))
			waitUntil(ctx, ready[b]-dist(b))
			cancel()
			goTo(b.X, b.Y)
			for !colony.CanHarvest() && colony.GetEntity() == colony.Reed && colony.Tick() < ready[b]+20 {
				colony.Wait()
			}
			colony.Harvest()
			colony.Plant(colony.Reed)
			ready[b] = colony.Tick() + 220
		}
		colony.Export(colony.Reed, 0)
	}
}

func waitUntil(ctx context.Context, tick int) {
	for colony.Tick() < tick {
		select {
		case <-ctx.Done():
			return
		default:
			colony.Wait()
		}
	}
}

func dist(p colony.Point) int {
	x, y := colony.Position()
	d := p.X - x + p.Y - y
	if p.X < x {
		d += 2 * (x - p.X)
	}
	if p.Y < y {
		d += 2 * (y - p.Y)
	}
	return d * 2
}

func nearIce(grid [][]colony.Ground, x, y int) bool {
	for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && ny >= 0 && nx < len(grid) && ny < len(grid[nx]) && grid[nx][ny] == colony.Ice {
			return true
		}
	}
	return false
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
