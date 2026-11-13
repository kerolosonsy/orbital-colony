package main

import (
	"sync"

	"orbital/colony"
)

// The workers finish before the final export so the ledger has one complete
// batch, and main waits for contract settlement without issuing more commands.
type ledger struct {
	mu      sync.Mutex
	pending int
}

func (l *ledger) add(n int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.pending += n
	return l.pending >= 300
}

func (l *ledger) full() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.pending >= 300
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
	colony.Export(colony.Moss, 0)
	select {}
}

func farm(d *colony.Drone, k, n int, book *ledger) {
	w, h := colony.WorldSize()
	for {
		for y := k; y < h; y += n {
			for i := 0; i < w; i++ {
				if book.full() {
					return
				}
				x := i
				if (y/n)%2 == 1 {
					x = w - 1 - i
				}
				goTo(d, x, y)
				if d.CanHarvest() {
					got, err := d.TryHarvest()
					d.Plant(colony.Moss)
					if err == nil && book.add(got) {
						return
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
