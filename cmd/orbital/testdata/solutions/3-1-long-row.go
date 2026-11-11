package main

import "orbital/colony"

func main() {
	for i := 0; i < 14; i++ {
		colony.Move(colony.East)
		colony.Harvest()
	}
}
