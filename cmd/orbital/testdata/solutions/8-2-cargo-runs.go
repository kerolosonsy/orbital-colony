package main

import "orbital/colony"

type Drone struct {
	X, Y  int
	Cargo int
}

func (d *Drone) Dist(p colony.Point) int { return abs(d.X-p.X) + abs(d.Y-p.Y) }

func (d *Drone) GoTo(p colony.Point) {
	for ; d.X < p.X; d.X++ {
		colony.Move(colony.East)
	}
	for ; d.X > p.X; d.X-- {
		colony.Move(colony.West)
	}
	for ; d.Y < p.Y; d.Y++ {
		colony.Move(colony.North)
	}
	for ; d.Y > p.Y; d.Y-- {
		colony.Move(colony.South)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func main() {
	d := &Drone{}
	stored := 0
	left := colony.Targets()
	for stored < 30 && len(left) > 0 {
		best := 0
		for i, t := range left {
			if d.Dist(t) < d.Dist(left[best]) {
				best = i
			}
		}
		t := left[best]
		d.GoTo(t)
		size := colony.Measure()
		if d.Cargo+size > 10 {
			d.GoTo(colony.Point{})
			colony.Unload()
			stored += d.Cargo
			d.Cargo = 0
			if stored >= 30 {
				return
			}
			continue
		}
		colony.Harvest()
		d.Cargo += size
		left = append(left[:best], left[best+1:]...)
		if stored+d.Cargo >= 30 {
			d.GoTo(colony.Point{})
			colony.Unload()
			return
		}
	}
	d.GoTo(colony.Point{})
	colony.Unload()
}
