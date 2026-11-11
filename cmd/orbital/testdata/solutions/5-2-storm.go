package main

import (
	"errors"
	"orbital/colony"
)

func main() {
	total := 0
	for x := 0; x < 8; x++ {
		for {
			n, err := colony.TryHarvest()
			if errors.Is(err, colony.ErrStorm) {
				continue
			}
			total += n
			break
		}
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	colony.Report(total)
}
