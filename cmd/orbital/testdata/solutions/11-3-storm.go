package main

import "orbital/colony"

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

func main() {
	d := colony.Main()
	warning := colony.StormWarning()
	calm := colony.StormClear()
	w, h := colony.WorldSize()
	var path []colony.Point
	for y := 0; y < h; y++ {
		for i := 0; i < w; i++ {
			x := i
			if y%2 == 1 {
				x = w - 1 - i
			}
			if x+y > 0 {
				path = append(path, colony.Point{X: x, Y: y})
			}
		}
	}
	for _, p := range path {
		select {
		case <-warning:
			x, y := d.Position()
			goTo(d, colony.Point{})
			<-calm
			goTo(d, colony.Point{X: x, Y: y})
		default:
		}
		goTo(d, p)
		d.Harvest()
	}
}
