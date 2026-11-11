package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	colony.Harvest()
	colony.Move(colony.East)
	colony.Harvest()
	colony.Move(colony.East)
	colony.Harvest()
}
