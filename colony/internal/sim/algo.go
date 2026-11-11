package sim

import (
	"fmt"
)

// Kg is a weight in kilograms.
type Kg float64

// Datasets are sensor readings for the generics chapter.
type Datasets struct {
	Callsigns  []string
	Readings   []string
	Faults     []*FaultError
	Orders     []string
	Questions  []Question
	Answers    map[int]int
	TotalCents int
	People     int
	Docked     bool
	Undocks    int
	Drills     int
	DrillPanic bool
	Depths     []int
	Temps      []float64
	Weights    []Kg
}

// Swap exchanges what is on the drone's tile with the neighbouring tile in
// direction dir. The drone does not move.
func (w *World) Swap(id int, dir Direction) bool {
	return w.act(id, "swap", func() (bool, string, bool) {
		if m := w.gate("api.swap"); m != "" {
			return false, m, false
		}
		dx, dy := 0, 0
		switch dir {
		case North:
			dy = 1
		case South:
			dy = -1
		case East:
			dx = 1
		case West:
			dx = -1
		}
		nx, ny := w.X+dx, w.Y+dy
		if nx < 0 || ny < 0 || nx >= w.W || ny >= w.H {
			return false, say("مفيش خانة ناحية ", "No tile to the ") + dir.String(), false
		}
		a, b := w.here(), w.AtZ(nx, ny, w.Z)
		a.Entity, b.Entity = b.Entity, a.Entity
		a.Size, b.Size = b.Size, a.Size
		a.Planted, b.Planted = b.Planted, a.Planted
		a.Ready, b.Ready = b.Ready, a.Ready
		w.emit(Event{Op: "tile", T: w.Tick, D: id, OK: true, X: w.X, Y: w.Y, Z: w.Z, Tile: ptr(w.viewZ(nx, ny, w.Z))})
		return true, fmt.Sprintf("⇄ %d ↔ %d", b.Size, a.Size), true
	})
}

func ptr[T any](v T) *T { return &v }

// GroundGrid returns grid[x][y] of every tile's ground.
func (w *World) GroundGrid() [][]Ground { return w.DroneGrid(0) }

func (w *World) DroneGrid(id int) [][]Ground {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.HideMap {
		return nil
	}
	z := w.layer(id)
	g := make([][]Ground, w.W)
	for x := range g {
		g[x] = make([]Ground, w.H)
		for y := range g[x] {
			g[x][y] = w.AtZ(x, y, z).Ground
		}
	}
	return g
}

// Neighbors returns the in-bounds neighbours of p that a drone can enter.
func (w *World) Neighbors(p Point) []Point {
	var out []Point
	for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		q := Point{p.X + d[0], p.Y + d[1]}
		if q.X < 0 || q.Y < 0 || q.X >= w.W || q.Y >= w.H {
			continue
		}
		if w.RockBlocks && w.At(q.X, q.Y).Ground == Rock {
			continue
		}
		out = append(out, q)
	}
	return out
}

// StepCost is the tick cost of entering q.
func (w *World) StepCost(q Point) int {
	if c := w.MoveCost[w.At(q.X, q.Y).Ground]; c > 0 {
		return c
	}
	return 1
}

// BFS returns the fewest-moves distance from a to every reachable tile.
func (w *World) BFS(a Point) map[Point]int {
	dist := map[Point]int{a: 0}
	q := []Point{a}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		for _, n := range w.Neighbors(p) {
			if _, ok := dist[n]; !ok {
				dist[n] = dist[p] + 1
				q = append(q, n)
			}
		}
	}
	return dist
}

// Cheapest returns the minimum tick cost from a to every reachable tile.
func (w *World) Cheapest(a Point) map[Point]int {
	cost := map[Point]int{a: 0}
	done := map[Point]bool{}
	for {
		best, found := Point{}, false
		for p, c := range cost {
			if !done[p] && (!found || c < cost[best]) {
				best, found = p, true
			}
		}
		if !found {
			return cost
		}
		done[best] = true
		for _, n := range w.Neighbors(best) {
			c := cost[best] + w.StepCost(n)
			if old, ok := cost[n]; !ok || c < old {
				cost[n] = c
			}
		}
	}
}
