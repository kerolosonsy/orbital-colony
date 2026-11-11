package main

import "orbital/colony"

func main() {
	colony.DigDown()
	colony.Descend()
	for i := 0; i < 3; i++ {
		colony.Dig(colony.East)
		colony.Move(colony.East)
	}
	for i := 0; i < 3; i++ {
		colony.Move(colony.West)
	}
	colony.Ascend()
}
