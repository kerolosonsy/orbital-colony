package main

import (
	"context"
	"orbital/colony"
)

func main() {
	ctx := context.Background()
	for x := 0; x < 10; x++ {
		sctx, cancel := context.WithTimeout(ctx, colony.Ticks(5))
		purity, err := colony.DeepScan(sctx)
		cancel()
		if err == nil && purity >= 50 {
			colony.Harvest()
		}
		if x < 9 {
			colony.Move(colony.East)
		}
	}
}
