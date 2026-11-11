package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

type LockedPulses struct {
	mu    sync.Mutex
	total int64
}
type AtomicPulses struct{ total atomic.Int64 }

func (p *LockedPulses) Add(delta int64) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.total += delta
	return p.total
}
func (p *LockedPulses) Total() int64          { p.mu.Lock(); defer p.mu.Unlock(); return p.total }
func (p *AtomicPulses) Add(delta int64) int64 { return p.total.Add(delta) }
func (p *AtomicPulses) Total() int64          { return p.total.Load() }
func TestPulses(t *testing.T) {
	for _, p := range []interface {
		Add(int64) int64
		Total() int64
	}{&LockedPulses{}, &AtomicPulses{}} {
		if p.Total() != 0 || p.Add(3) != 3 || p.Add(-4) != -1 || p.Add(0) != -1 || p.Total() != -1 {
			t.Fatal("sequential contract")
		}
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				for j := 0; j < 100; j++ {
					p.Add(1)
					_ = p.Total()
				}
			}()
		}
		close(start)
		wg.Wait()
		if p.Total() != 1599 {
			t.Fatal("concurrent total")
		}
	}
}
func BenchmarkLockedPulses(b *testing.B) {
	var p LockedPulses
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Add(1)
		}
	})
	if p.Total() != int64(b.N) {
		b.Fatal("lost work")
	}
}
func BenchmarkAtomicPulses(b *testing.B) {
	var p AtomicPulses
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Add(1)
		}
	})
	if p.Total() != int64(b.N) {
		b.Fatal("lost work")
	}
}
