package main

import (
	"orbital/colony/labkit"
	"testing"
	"time"
)

func assertTake(t *testing.T, gate *LaunchGate, n int, want bool) {
	t.Helper()
	got, err := gate.Take(n)
	if err != nil || got != want {
		t.Fatalf("Take(%d)=%v,%v want %v", n, got, err, want)
	}
}
func TestLaunchGate(t *testing.T) {
	t.Run("exact boundary and fractional carry", func(t *testing.T) {
		clock := labkit.NewFakeClock(time.Unix(0, 0))
		gate, err := NewLaunchGate(clock, 2, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		assertTake(t, gate, 2, true)
		assertTake(t, gate, 1, false)
		clock.Advance(999 * time.Millisecond)
		assertTake(t, gate, 1, false)
		clock.Advance(time.Millisecond)
		assertTake(t, gate, 1, true)
		clock.Advance(500 * time.Millisecond)
		assertTake(t, gate, 1, false)
		clock.Advance(500 * time.Millisecond)
		assertTake(t, gate, 1, true)
	})
	t.Run("ceiling and rejection", func(t *testing.T) {
		clock := labkit.NewFakeClock(time.Unix(0, 0))
		gate, err := NewLaunchGate(clock, 2, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		assertTake(t, gate, 3, false)
		assertTake(t, gate, 2, true)
		clock.Advance(10 * time.Second)
		assertTake(t, gate, 3, false)
		assertTake(t, gate, 2, true)
		assertTake(t, gate, 1, false)
	})
	t.Run("invalid input leaves credit intact", func(t *testing.T) {
		clock := labkit.NewFakeClock(time.Unix(0, 0))
		for _, tc := range []struct {
			cap   int
			every time.Duration
		}{{0, time.Second}, {-1, time.Second}, {1, 0}, {1, -time.Second}} {
			if gate, err := NewLaunchGate(clock, tc.cap, tc.every); err == nil || gate != nil {
				t.Fatal("invalid configuration")
			}
		}
		if gate, err := NewLaunchGate(nil, 1, time.Second); err == nil || gate != nil {
			t.Fatal("nil clock")
		}
		gate, err := NewLaunchGate(clock, 2, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		for _, n := range []int{0, -1} {
			if ok, err := gate.Take(n); err == nil || ok {
				t.Fatal("invalid request")
			}
		}
		assertTake(t, gate, 2, true)
	})
}
