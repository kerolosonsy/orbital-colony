package main

import (
	"context"
	"fmt"
	"orbital/colony/labkit"
	"sync"
	"time"
)

type DockTimer struct {
	ID    string
	Delay time.Duration
}
type DockResult struct{ ID, Status string }

func AwaitDocks(ctx context.Context, clock labkit.Clock, jobs []DockTimer) ([]DockResult, error) {
	if clock == nil {
		return nil, fmt.Errorf("clock required")
	}
	seen := map[string]bool{}
	for _, job := range jobs {
		if job.ID == "" || seen[job.ID] || job.Delay < 0 {
			return nil, fmt.Errorf("invalid dock timer")
		}
		seen[job.ID] = true
	}
	results := make([]DockResult, len(jobs))
	var wg sync.WaitGroup
	for i, job := range jobs {
		wg.Add(1)
		go func(i int, job DockTimer) {
			defer wg.Done()
			results[i] = DockResult{job.ID, "canceled"}
			if ctx.Err() != nil {
				return
			}
			if job.Delay == 0 {
				results[i].Status = "done"
				return
			}
			ticker := clock.NewTicker(job.Delay)
			defer ticker.Stop()
			select {
			case <-ctx.Done():
			case <-ticker.C():
				if ctx.Err() == nil {
					results[i].Status = "done"
				}
			}
		}(i, job)
	}
	wg.Wait()
	return results, nil
}
func main() {}
