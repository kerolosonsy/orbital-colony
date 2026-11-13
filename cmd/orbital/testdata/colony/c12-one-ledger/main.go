package main

import (
	"sync"

	"orbital/colony"
)

// Reference for c12-one-ledger. Four drones farm interleaved rows and record
// every harvest in one ledger guarded by a mutex; whoever pushes the ledger
// past a batch sells it.
type ledger struct {
	mu      sync.Mutex
	pending int
}

func (l *ledger) add(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.pending += n
	if l.pending >= 20 {
		l.pending = 0
		return true
	}
	return false
}

func main() {
	book := &ledger{}
	drones := []*colony.Drone{colony.Main()}
	for len(drones) < 4 {
		d := colony.Spawn()
		if d == nil {
			break
		}
		drones = append(drones, d)
	}
	var wg sync.WaitGroup
	for k, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			farm(d, k, len(drones), book)
		}()
	}
	wg.Wait()
}

func farm(d *colony.Drone, k, n int, book *ledger) {
	w, h := colony.WorldSize()
	for {
		for y := k; y < h; y += n {
			for i := 0; i < w; i++ {
				x := i
				if (y/n)%2 == 1 {
					x = w - 1 - i
				}
				goTo(d, x, y)
				if d.CanHarvest() {
					got, err := d.TryHarvest()
					d.Plant(colony.Moss)
					if err == nil && book.add(got) {
						d.Export(colony.Moss, 0)
					}
				} else if d.GetEntity() == colony.Empty && d.Till() {
					d.Plant(colony.Moss)
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
