package main

import (
	"fmt"
	"orbital/colony"
)

type Status struct {
	Moss, Crystal int
}

func (s Status) String() string {
	return fmt.Sprintf("moss=%d crystal=%d", s.Moss, s.Crystal)
}

func main() {
	colony.Report(Status{colony.Count(colony.Moss), colony.Count(colony.Crystal)})
}
