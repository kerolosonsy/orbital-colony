package main

import "orbital/colony"

// Reference for e07-cluster. Every tile that takes a seed grows glowcap.
// Once all of them are ripe, only the "black" squares of a checkerboard are
// harvested and replanted: their neighbors stay ripe forever, so each black
// harvest counts up to four ripe neighbors.
type plot struct {
	x, y  int
	ready int
}

func (p *plot) black() bool { return (p.x+p.y)%2 == 0 }

func (p *plot) visit() { goTo(p.x, p.y) }

func main() {
	const grow = 140
	w, h := colony.WorldSize()
	var plots []*plot
	for y := 0; y < h; y++ {
		for i := 0; i < w; i++ {
			x := i
			if y%2 == 1 {
				x = w - 1 - i
			}
			goTo(x, y)
			if colony.Till() && colony.Plant(colony.Glowcap) {
				plots = append(plots, &plot{x: x, y: y, ready: colony.Tick() + grow})
			}
		}
	}
	for {
		for _, p := range plots {
			if !p.black() {
				continue
			}
			for colony.Tick() < p.ready {
				colony.Wait()
			}
			p.visit()
			colony.Harvest()
			colony.Plant(colony.Glowcap)
			p.ready = colony.Tick() + grow
		}
		colony.Export(colony.Glowcap, 0)
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
