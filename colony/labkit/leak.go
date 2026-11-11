package labkit

import (
	"runtime"
	"time"
)

type leakTest interface {
	Helper()
	Cleanup(func())
	Errorf(string, ...any)
}

// LeakCheck registers a cleanup; call it before starting the lab's goroutines.
// Counts are process-wide, so do not use it in parallel tests.
func LeakCheck(t leakTest) {
	t.Helper()
	before := runtime.NumGoroutine()
	t.Cleanup(func() {
		for i := 0; i < 20; i++ {
			if runtime.NumGoroutine() <= before {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Errorf("goroutine leak: before=%d after=%d", before, runtime.NumGoroutine())
	})
}
