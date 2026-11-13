package main

import "orbital/colony"

// Reference for c17-deep-haul. Open a shaft down to layer 2, then tunnel
// through the whole layer in a snake. Every stone tile dug collects its ore;
// geodes are sold as they come in.
func main() {
	for colony.Layer() < 2 {
		colony.DigDown()
		colony.Descend()
	}
	w, h := colony.WorldSize()
	dug := make([][]bool, w)
	for x := range dug {
		dug[x] = make([]bool, h)
	}
	dug[0][0] = true
	for y := 0; y < h; y++ {
		for i := 1; i < w; i++ {
			dir := colony.East
			if y%2 == 1 {
				dir = colony.West
			}
			x, _ := colony.Position()
			if dir == colony.East {
				x++
			} else {
				x--
			}
			if !dug[x][y] {
				colony.Dig(dir)
				dug[x][y] = true
			}
			colony.Move(dir)
			if colony.Count(colony.Geode) >= 6 {
				colony.Export(colony.Geode, 0)
			}
		}
		if y < h-1 {
			colony.Dig(colony.North)
			colony.Move(colony.North)
		}
	}
	colony.Export(colony.Geode, 0)
	for {
		colony.Wait()
	}
}
