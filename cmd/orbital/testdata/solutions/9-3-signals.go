package main

import "orbital/colony"

func goTo(p colony.Point) {
	x, y := colony.Position()
	for ; x < p.X; x++ {
		colony.Move(colony.East)
	}
	for ; x > p.X; x-- {
		colony.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		colony.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		colony.Move(colony.South)
	}
}

func dist(a, b colony.Point) int {
	dx, dy := a.X-b.X, a.Y-b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func main() {
	total := 0
	var calls []colony.Point
	for _, s := range colony.Signals() {
		switch v := s.(type) {
		case colony.Distress:
			calls = append(calls, v.At)
		case colony.Supply:
			total += v.Amount
		}
	}
	colony.Report(total)
	here := colony.Point{}
	for len(calls) > 0 {
		best := 0
		for i, c := range calls {
			if dist(here, c) < dist(here, calls[best]) {
				best = i
			}
		}
		here = calls[best]
		goTo(here)
		calls = append(calls[:best], calls[best+1:]...)
	}
}
