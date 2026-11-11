package labkit

import (
	"sync"
	"testing"
	"time"
)

func TestFakeClockTimersAndStoppedTicker(t *testing.T) {
	start := time.Unix(100, 0)
	clock := NewFakeClock(start)
	timer := clock.After(2 * time.Second)
	ticker := clock.NewTicker(time.Second)
	clock.Advance(time.Second)
	if got := <-ticker.C(); !got.Equal(start.Add(time.Second)) {
		t.Fatal(got)
	}
	select {
	case <-timer:
		t.Fatal("timer fired early")
	default:
	}
	clock.Advance(time.Second)
	if got := <-timer; !got.Equal(start.Add(2 * time.Second)) {
		t.Fatal(got)
	}
	<-ticker.C()
	ticker.Stop()
	clock.Advance(time.Hour)
	select {
	case <-ticker.C():
		t.Fatal("stopped ticker fired")
	default:
	}
	select {
	case <-timer:
		t.Fatal("one-shot timer fired twice")
	default:
	}
	if !clock.Now().Equal(start.Add(time.Hour + 2*time.Second)) {
		t.Fatal(clock.Now())
	}
	if got := <-clock.After(0); !got.Equal(clock.Now()) {
		t.Fatal(got)
	}
}

func TestFakeClockDropsTicksWithoutBlocking(t *testing.T) {
	var clock FakeClock
	ticker := clock.NewTicker(time.Second)
	clock.Advance(10 * time.Second)
	if got := <-ticker.C(); !got.Equal(time.Time{}.Add(time.Second)) {
		t.Fatal(got)
	}
	clock.Advance(time.Second)
	if got := <-ticker.C(); !got.Equal(time.Time{}.Add(11 * time.Second)) {
		t.Fatal(got)
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				clock.Now()
				clock.Advance(time.Second)
				clock.After(time.Second)
			}
		}()
	}
	wg.Wait()
}

func TestFakeClockRejectsInvalidDurations(t *testing.T) {
	for _, call := range []func(){func() { new(FakeClock).NewTicker(0) }, func() { new(FakeClock).Advance(-1) }} {
		func() {
			defer func() {
				if recover() == nil {
					t.Error("expected panic")
				}
			}()
			call()
		}()
	}
}
