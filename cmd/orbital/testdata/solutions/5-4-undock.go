package main

import (
	"fmt"
	"orbital/colony"
)

func drillBay() {
	colony.Dock()
	defer colony.Undock()
	defer func() {
		if fault := recover(); fault != nil {
			fmt.Println("drill recovered:", fault)
		}
	}()
	colony.Drill()
}
func main() {
	attempts := colony.DrillAttempts()
	for i := 0; i < attempts; i++ {
		drillBay()
	}
	colony.Move(colony.East)
	colony.Report(fmt.Sprintf("ready:%d", attempts))
}
