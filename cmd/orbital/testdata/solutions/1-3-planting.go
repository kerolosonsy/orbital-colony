package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	colony.Till()
	colony.Plant(colony.Moss)
	colony.Move(colony.East)
	colony.Till()
	colony.Plant(colony.Moss)
	colony.Move(colony.East)
	colony.Till()
	colony.Plant(colony.Moss)
}
