package main

import (
	"math/rand"
	"sync"
	"testing"
)

func TestScans(t *testing.T) {
	if !t.Run("controlled lost update", func(t *testing.T) {
		arrived := make(chan int, 2)
		finished := make(chan int, 2)
		release := []chan struct{}{make(chan struct{}), make(chan struct{})}
		done := make(chan int, 1)
		go func() {
			got, e := MergeScans([]int{7, 11}, func(i int) func() { arrived <- i; <-release[i]; return func() { finished <- i } })
			if e != nil {
				t.Error(e)
			}
			done <- got
		}()
		<-arrived
		<-arrived
		for i := 0; i < 2; i++ {
			close(release[i])
			if got := <-finished; got != i {
				t.Error("ack index")
			}
		}
		if got := <-done; got != 18 {
			t.Fatal("stale read lost update", got)
		}
	}) {
		return
	}
	t.Run("signed exact totals and visits", func(t *testing.T) {
		rng := rand.New(rand.NewSource(2521))
		for _, size := range []int{0, 1, 2, 64} {
			scans := make([]int, size)
			want := 0
			for i := range scans {
				scans[i] = rng.Intn(30) - 20
				want += scans[i]
			}
			var mu sync.Mutex
			visits := make([]int, size)
			got, e := MergeScans(scans, func(i int) func() { mu.Lock(); visits[i]++; mu.Unlock(); return func() {} })
			if e != nil || got != want {
				t.Fatal("sum", got, want, e)
			}
			for _, n := range visits {
				if n != 1 {
					t.Fatal("visits")
				}
			}
		}
	})
	t.Run("preflight", func(t *testing.T) {
		if n, e := MergeScans(nil, nil); n != 0 || e == nil {
			t.Fatal("nil callback")
		}
		if n, e := MergeScans(make([]int, 65), func(int) func() { t.Error("callback"); return func() {} }); n != 0 || e == nil {
			t.Fatal("oversize")
		}
	})
}
