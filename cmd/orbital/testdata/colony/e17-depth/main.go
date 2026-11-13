package main

import "orbital/colony"

var path []colony.Direction

func retry(action func() bool) bool {
	for {
		ok := action()
		if colony.Fault() == nil {
			return ok
		}
		colony.Repair()
	}
}

func move(dir colony.Direction) {
	retry(func() bool { return colony.Move(dir) })
}

func deliver() {
	for i := len(path) - 1; i >= 0; i-- {
		move((path[i] + 2) % 4)
	}
	for colony.Layer() > 0 {
		retry(colony.Ascend)
	}
	retry(colony.Unload)
	for {
		colony.Export(colony.Geode, 0)
		if colony.Fault() == nil {
			break
		}
		colony.Repair()
	}
	for colony.Layer() < 2 {
		retry(colony.Descend)
	}
	for _, dir := range path {
		move(dir)
	}
}

func step(dir colony.Direction) {
	if !retry(func() bool { return colony.Dig(dir) }) {
		deliver()
		retry(func() bool { return colony.Dig(dir) })
	}
	move(dir)
	path = append(path, dir)
	if colony.Carry()[colony.Geode] >= 6 {
		deliver()
	}
}

func main() {
	for colony.Layer() < 2 {
		retry(colony.DigDown)
		retry(colony.Descend)
	}
	w, h := colony.WorldSize()
	dug := make([][]bool, w)
	for x := range dug {
		dug[x] = make([]bool, h)
	}
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
				step(dir)
				dug[x][y] = true
			}
		}
		if y < h-1 {
			step(colony.North)
		}
	}
	deliver()
	for {
		colony.Wait()
	}
}
