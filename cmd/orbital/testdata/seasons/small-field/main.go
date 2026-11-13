package main

import (
	"errors"
	"orbital/colony"
)

func main() {
	w, h := colony.WorldSize()
	planted := make([][]bool, h)
	for y := range planted {
		planted[y] = make([]bool, w)
	}
	for {
		for y := 0; y < h; y++ {
			for column := 0; column < w; column++ {
				x := column
				if y%2 == 1 {
					x = w - 1 - column
				}
				goTo(x, y)
				if !planted[y][x] {
					planted[y][x] = colony.Till() && colony.Plant(colony.Glowcap)
				} else if _, err := colony.TryHarvest(); err == nil || errors.Is(err, colony.ErrEmpty) {
					colony.Plant(colony.Glowcap)
				}
				if colony.Count(colony.Glowcap) >= 20 {
					colony.Export(colony.Glowcap, 0)
				}
			}
		}
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
