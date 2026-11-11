package main

import "testing"

func TestCountPairs(t *testing.T) {
	if got := CountPairs([]int{1, 9, 5, 5, 3, 7}, 10); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func BenchmarkCountPairs(b *testing.B) {
	xs := make([]int, 3000)
	for i := range xs {
		xs[i] = i % 1000
	}
	for b.Loop() {
		CountPairs(xs, 1000)
	}
}
