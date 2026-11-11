package main

import (
	"context"
	"fmt"
	"sync"
)

type OreBatch struct{ ID, Grade int }
type OreReply struct {
	ID, Grade int
	Err       error
}

func oreLane(ctx context.Context, input <-chan OreBatch, refine func(context.Context, OreBatch) (int, error)) <-chan OreReply {
	out := make(chan OreReply)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case batch, ok := <-input:
				if !ok || ctx.Err() != nil {
					return
				}
				grade, err := refine(ctx, batch)
				select {
				case out <- OreReply{batch.ID, grade, err}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
func RefineOre(ctx context.Context, input <-chan OreBatch, lanes int, refine func(context.Context, OreBatch) (int, error)) (<-chan OreReply, error) {
	if input == nil || lanes <= 0 || refine == nil {
		return nil, fmt.Errorf("invalid refinery")
	}
	merged := make(chan OreReply)
	var forwards sync.WaitGroup
	for i := 0; i < lanes; i++ {
		lane := oreLane(ctx, input, refine)
		forwards.Add(1)
		go func() {
			defer forwards.Done()
			for {
				select {
				case <-ctx.Done():
					// Join the worker by draining its output before this forwarder exits.
					for range lane {
					}
					return
				case reply, ok := <-lane:
					if !ok {
						return
					}
					select {
					case merged <- reply:
					case <-ctx.Done():
						for range lane {
						}
						return
					}
				}
			}
		}()
	}
	go func() { forwards.Wait(); close(merged) }()
	return merged, nil
}
func main() {}
