package main

import (
	"fmt"
	"orbital/colony"
	"strings"
)

func reportMoss()    { colony.Report(fmt.Sprintf("moss=%d", colony.Count(colony.Moss))) }
func reportCrystal() { colony.Report(fmt.Sprintf("crystal=%d", colony.Count(colony.Crystal))) }
func waitStep()      { colony.Wait(); colony.Report("wait") }
func commandAction(name string) func() {
	switch name {
	case "moss":
		return reportMoss
	case "crystal":
		return reportCrystal
	case "wait":
		return waitStep
	}
	return nil
}
func execute(action func()) { action() }
func main() {
	count := 0
	countAction := func() { count++ }
	for i := 0; i < colony.OrderCount(); i++ {
		command := strings.TrimSpace(colony.Order(i))
		action := commandAction(command)
		if action == nil {
			colony.Report(fmt.Errorf("unknown:%s", command))
		} else {
			execute(action)
			execute(countAction)
		}
	}
	colony.Report(count)
}
