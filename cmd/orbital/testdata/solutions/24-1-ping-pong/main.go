package main

import (
	"context"
	"fmt"
	"sync"
)

func RelaySupplies(ctx context.Context, rounds int, visit func(drone, round int) error) error {
	if rounds < 0 || visit == nil {
		return fmt.Errorf("invalid relay")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if rounds == 0 {
		return nil
	}
	run, cancel := context.WithCancel(ctx)
	defer cancel()
	permit := make(chan int)
	ack := make(chan struct{})
	var wg sync.WaitGroup
	var faults [2]error
	wg.Add(2)
	go func() {
		defer wg.Done()
		for r := 0; r < rounds; r++ {
			if run.Err() != nil {
				return
			}
			if err := visit(0, r); err != nil {
				faults[0] = err
				cancel()
				return
			}
			select {
			case permit <- r:
			case <-run.Done():
				return
			}
			select {
			case <-ack:
			case <-run.Done():
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for r := 0; r < rounds; r++ {
			select {
			case <-run.Done():
				return
			case round := <-permit:
				if run.Err() != nil {
					return
				}
				if err := visit(1, round); err != nil {
					faults[1] = err
					cancel()
					return
				}
			}
			select {
			case ack <- struct{}{}:
			case <-run.Done():
				return
			}
		}
	}()
	wg.Wait()
	close(permit)
	close(ack)
	for _, err := range faults {
		if err != nil {
			return err
		}
	}
	return ctx.Err()
}
func main() {}
