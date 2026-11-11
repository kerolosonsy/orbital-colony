package main

import "orbital/colony"

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

var delta = map[colony.Direction][2]int{colony.North: {0, 1}, colony.East: {1, 0}, colony.South: {0, -1}, colony.West: {-1, 0}}

func main() {
	w, h := colony.WorldSize()
	var path Stack[colony.Direction]
	at := colony.Point{}
	seen := map[colony.Point]bool{at: true}
	for colony.GetEntity() != colony.Crystal {
		moved := false
		for _, d := range []colony.Direction{colony.North, colony.East, colony.South, colony.West} {
			n := colony.Point{X: at.X + delta[d][0], Y: at.Y + delta[d][1]}
			if n.X < 0 || n.Y < 0 || n.X >= w || n.Y >= h || seen[n] {
				continue
			}
			seen[n] = true
			if colony.Move(d) {
				at, moved = n, true
				path.Push(d)
				break
			}
		}
		if !moved {
			d, _ := path.Pop()
			back := (d + 2) % 4
			colony.Move(back)
			at = colony.Point{X: at.X + delta[back][0], Y: at.Y + delta[back][1]}
			continue
		}
	}
	colony.Harvest()
	for d, ok := path.Pop(); ok; d, ok = path.Pop() {
		colony.Move((d + 2) % 4)
	}
}
