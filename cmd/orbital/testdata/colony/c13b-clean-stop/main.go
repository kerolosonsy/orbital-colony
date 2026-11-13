package main

import "orbital/colony"

func home() {
	for {
		if colony.Fault() != nil {
			colony.Repair()
			continue
		}
		x, y := colony.Position()
		if x > 0 {
			colony.Move(colony.West)
		} else if y > 0 {
			colony.Move(colony.South)
		} else {
			break
		}
	}
	for !colony.Unload() {
		if len(colony.Carry()) == 0 {
			break
		}
		colony.Repair()
	}
}

func main() {
	stop := colony.Shutdown()
	colony.Till()
	colony.Plant(colony.Moss)
	for {
		select {
		case <-stop:
			home()
			return
		default:
			if colony.Fault() != nil {
				colony.Repair()
				continue
			}
			if colony.CanHarvest() {
				colony.Harvest()
			}
			colony.Move(colony.East)
			colony.Move(colony.West)
		}
	}
}
