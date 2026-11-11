package main

import "orbital/colony"

func main() {
	scout := colony.Spawn()
	found := make(chan int)
	go func() {
		for x := 0; x < 12; x++ {
			if scout.Measure() >= 5 {
				found <- x
			}
			if x < 11 {
				scout.Move(colony.East)
			}
		}
		close(found)
	}()
	at := 0
	for x := range found {
		for ; at < x; at++ {
			colony.Move(colony.East)
		}
		colony.Harvest()
	}
}
