package main

import (
	"orbital/colony"
	"sync"
)

type Registry struct {
	mu    sync.Mutex
	sizes map[colony.Point]int
}

func (r *Registry) Record(p colony.Point, size int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sizes[p] = size
}

func (r *Registry) Best() colony.Point {
	r.mu.Lock()
	defer r.mu.Unlock()
	best, bp := -1, colony.Point{}
	for p, s := range r.sizes {
		if s > best {
			best, bp = s, p
		}
	}
	return bp
}

func main() {
	reg := &Registry{sizes: map[colony.Point]int{}}
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	for row, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < row; i++ {
				d.Move(colony.North)
			}
			for x := 0; x < 10; x++ {
				reg.Record(colony.Point{X: x, Y: row}, d.Measure())
				if x < 9 {
					d.Move(colony.East)
				}
			}
		}()
	}
	wg.Wait()
	colony.Report(reg.Best())
}
