package main

import (
	"context"
	"fmt"
	"sync"
)

type CrewReply[R any] struct {
	Value R
	Err   error
}

func RunCrew[T, R any](ctx context.Context, inputs []T, workers int, work func(context.Context, T) (R, error)) ([]CrewReply[R], error) {
	if workers <= 0 || work == nil {
		return nil, fmt.Errorf("invalid crew")
	}
	replies := make([]CrewReply[R], len(inputs))
	started := make([]bool, len(inputs))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < min(workers, len(inputs)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				if ctx.Err() != nil {
					return
				}
				started[index] = true
				replies[index].Value, replies[index].Err = work(ctx, inputs[index])
			}
		}()
	}
feed:
	for i := range inputs {
		if ctx.Err() != nil {
			break
		}
		select {
		case jobs <- i:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()
	if err := ctx.Err(); err != nil {
		for i := range replies {
			if !started[i] {
				replies[i].Err = err
			}
		}
		return replies, err
	}
	return replies, nil
}
func main() {}
