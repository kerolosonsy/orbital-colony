package main

import "orbital/colony"

// Reference for c09-crop-protocol. Every crop is a type behind one interface
// that says which ground it needs and how to prepare it. Ice grows lichen,
// the upper rows grow glowcap in a cluster, the rest grows moss.
type Crop interface {
	Entity() colony.Entity
	Prepare() bool
}

type moss struct{}
type lichen struct{}
type glowcap struct{}

func (moss) Entity() colony.Entity    { return colony.Moss }
func (moss) Prepare() bool            { return colony.Till() }
func (lichen) Entity() colony.Entity  { return colony.Lichen }
func (lichen) Prepare() bool          { return true }
func (glowcap) Entity() colony.Entity { return colony.Glowcap }
func (glowcap) Prepare() bool         { return colony.Till() }

func cropFor(g colony.Ground, y, h int) Crop {
	switch {
	case g == colony.Ice:
		return lichen{}
	case g == colony.Rock:
		return nil
	case y >= h-3:
		return glowcap{}
	}
	return moss{}
}

func main() {
	w, h := colony.WorldSize()
	grid := colony.Grid()
	for {
		for y := 0; y < h; y++ {
			for i := 0; i < w; i++ {
				x := i
				if y%2 == 1 {
					x = w - 1 - i
				}
				c := cropFor(grid[x][y], y, h)
				if c == nil {
					continue
				}
				goTo(x, y)
				if colony.CanHarvest() {
					colony.Harvest()
					colony.Plant(c.Entity())
				} else if colony.GetEntity() == colony.Empty && c.Prepare() {
					colony.Plant(c.Entity())
				}
			}
		}
		for _, e := range []colony.Entity{colony.Moss, colony.Lichen, colony.Glowcap} {
			colony.Export(e, 0)
		}
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
