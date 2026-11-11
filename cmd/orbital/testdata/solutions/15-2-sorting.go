package main

import "orbital/colony"

func main() {
	n := 8
	s := make([]int, n)
	for x := 0; x < n; x++ {
		s[x] = colony.Measure()
		if x < n-1 {
			colony.Move(colony.East)
		}
	}
	at := n - 1
	walk := func(to int) {
		for ; at > to; at-- {
			colony.Move(colony.West)
		}
		for ; at < to; at++ {
			colony.Move(colony.East)
		}
	}
	for swapped := true; swapped; {
		swapped = false
		for i := n - 2; i >= 0; i-- {
			if s[i] > s[i+1] {
				walk(i)
				colony.Swap(colony.East)
				s[i], s[i+1] = s[i+1], s[i]
				swapped = true
			}
		}
		for i := 0; i < n-1; i++ {
			if s[i] > s[i+1] {
				walk(i)
				colony.Swap(colony.East)
				s[i], s[i+1] = s[i+1], s[i]
				swapped = true
			}
		}
	}
}
