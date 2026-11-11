package main

import (
	"context"
	"fmt"
	"sync"
)

func UnloadDepots(ctx context.Context, batches [][]int, capacity int, consume func(int) error) error {
	if capacity < 0 || consume == nil {
		return fmt.Errorf("invalid depots")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	run, cancel := context.WithCancel(ctx)
	defer cancel()
	crates := make(chan int, capacity)
	var producers sync.WaitGroup
	for _, batch := range batches {
		if len(batch) == 0 {
			continue
		}
		producers.Add(1)
		go func(batch []int) {
			defer producers.Done()
			for _, code := range batch {
				if run.Err() != nil {
					return
				}
				select {
				case crates <- code:
				case <-run.Done():
					return
				}
			}
		}(batch)
	}
	joined := make(chan struct{})
	go func() { producers.Wait(); close(crates); close(joined) }()
	var fault error
	for code := range crates {
		if fault != nil || run.Err() != nil {
			continue
		}
		if err := consume(code); err != nil {
			fault = err
			cancel()
		}
	}
	<-joined
	if fault != nil {
		return fault
	}
	return ctx.Err()
}
func main() {}
