package main

import (
	"errors"

	"orbital/colony"
)

// Reference for e05-care. The scanner is switched off, so every
// tile is probed with TryHarvest: a ripe plant is collected and replanted, an
// unripe one is left alone, and an empty tile gets a seed.
func main() {
	w, h := colony.WorldSize()
	for {
		for y := 0; y < h; y++ {
			for i := 0; i < w; i++ {
				x := i
				if y%2 == 1 {
					x = w - 1 - i
				}
				goTo(x, y)
				if err := tend(); err != nil && !errors.Is(err, colony.ErrNotRipe) {
					return
				}
			}
		}
		colony.Export(colony.Moss, 0)
		colony.Export(colony.Lichen, 0)
		goTo(0, 0)
	}
}

// tend harvests a ripe plant and replants; it reports only unexpected errors.
func tend() error {
	_, err := colony.TryHarvest()
	if errors.Is(err, colony.ErrFault) {
		colony.Repair()
		return tend()
	}
	switch {
	case errors.Is(err, colony.ErrNotRipe):
		return err
	case err != nil && !errors.Is(err, colony.ErrEmpty):
		return err
	}
	if !colony.Plant(colony.Moss) && !colony.Plant(colony.Lichen) && colony.Till() {
		colony.Plant(colony.Moss)
	}
	return nil
}

func goTo(x, y int) {
	cx, cy := colony.Position()
	for ; cx < x; cx++ {
		for !colony.Move(colony.East) {
			colony.Repair()
		}
	}
	for ; cx > x; cx-- {
		for !colony.Move(colony.West) {
			colony.Repair()
		}
	}
	for ; cy < y; cy++ {
		for !colony.Move(colony.North) {
			colony.Repair()
		}
	}
	for ; cy > y; cy-- {
		for !colony.Move(colony.South) {
			colony.Repair()
		}
	}
	if colony.Fault() != nil {
		colony.Repair()
	}
}
