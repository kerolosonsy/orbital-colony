package main

import (
	"errors"
	"sync"
)

// A data race has unordered conflicting accesses; the mutex establishes happens-before.
// Reading and updating must share the critical section to avoid stale additions.
func MergeScans(scans []int, checkpoint func(int) func()) (int, error) {
	if checkpoint == nil || len(scans) > 64 {
		return 0, errors.New("invalid scans")
	}
	total := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i, n := range scans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ack := checkpoint(i)
			mu.Lock()
			total += n
			mu.Unlock()
			ack()
		}()
	}
	wg.Wait()
	return total, nil
}
func main() {}
