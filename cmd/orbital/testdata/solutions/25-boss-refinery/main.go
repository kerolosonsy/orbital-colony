package main

import (
	"context"
	"errors"
	"sync"
)

var ErrCrew = errors.New("crew must be positive")

func Refine(ctx context.Context, ore []int, crew int, assay func(context.Context, int) (int, error)) ([]int, error) {
	if crew <= 0 {
		return nil, ErrCrew
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if crew > len(ore) {
		crew = len(ore)
	}
	out := make([]int, len(ore))
	jobs := make(chan int)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var first error
	for w := 0; w < crew; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				if ctx.Err() != nil {
					return
				}
				v, err := assay(ctx, ore[i])
				if err != nil {
					mu.Lock()
					if first == nil {
						first = err
						cancel()
					}
					mu.Unlock()
					return
				}
				out[i] = v
			}
		}()
	}
	func() {
		for i := range ore {
			select {
			case jobs <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	close(jobs)
	wg.Wait()
	if first != nil {
		return nil, first
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
func main() {}
