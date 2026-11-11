package labkit

import "testing"

type leakRecorder struct {
	cleanup func()
	failed  bool
}

func (*leakRecorder) Helper()                 {}
func (r *leakRecorder) Cleanup(f func())      { r.cleanup = f }
func (r *leakRecorder) Errorf(string, ...any) { r.failed = true }

func TestLeakCheckDetectsUnfinishedGoroutine(t *testing.T) {
	r := &leakRecorder{}
	LeakCheck(r)
	stop, ready, done := make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() { close(ready); <-stop; close(done) }()
	<-ready
	r.cleanup()
	close(stop)
	<-done
	if !r.failed {
		t.Fatal("missed live goroutine")
	}
}
func TestLeakCheckAllowsFinishedGoroutine(t *testing.T) {
	r := &leakRecorder{}
	LeakCheck(r)
	done := make(chan struct{})
	go func() { close(done) }()
	<-done
	r.cleanup()
	if r.failed {
		t.Fatal("reported finished goroutine")
	}
}
