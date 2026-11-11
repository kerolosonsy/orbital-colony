// Package labkit provides deterministic tools for player lab tests.
package labkit

import (
	"sync"
	"time"
)

type Ticker interface {
	C() <-chan time.Time
	Stop()
}

type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
	NewTicker(time.Duration) Ticker
}

type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }
func (RealClock) NewTicker(d time.Duration) Ticker       { return realTicker{time.NewTicker(d)} }

type realTicker struct{ *time.Ticker }

func (t realTicker) C() <-chan time.Time { return t.Ticker.C }

type FakeClock struct {
	mu     sync.Mutex
	now    time.Time
	timers map[*fakeTimer]bool
}

type fakeTimer struct {
	clock  *FakeClock
	ch     chan time.Time
	next   time.Time
	period time.Duration
}

// NewFakeClock starts at start. Advance delivers due events without sleeping.
// Like time.Ticker, fake tickers drop ticks when their channel is full.
func NewFakeClock(start time.Time) *FakeClock               { return &FakeClock{now: start} }
func (c *FakeClock) Now() time.Time                         { c.mu.Lock(); defer c.mu.Unlock(); return c.now }
func (c *FakeClock) After(d time.Duration) <-chan time.Time { return c.timer(d, 0).ch }
func (c *FakeClock) NewTicker(d time.Duration) Ticker {
	if d <= 0 {
		panic("non-positive ticker interval")
	}
	return c.timer(d, d)
}
func (c *FakeClock) timer(d, period time.Duration) *fakeTimer {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.timers == nil {
		c.timers = map[*fakeTimer]bool{}
	}
	timer := &fakeTimer{clock: c, ch: make(chan time.Time, 1), next: c.now.Add(d), period: period}
	if d <= 0 {
		timer.ch <- c.now
	} else {
		c.timers[timer] = true
	}
	return timer
}
func (c *FakeClock) Advance(d time.Duration) {
	if d < 0 {
		panic("negative clock advance")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
	for timer := range c.timers {
		if timer.next.After(c.now) {
			continue
		}
		select {
		case timer.ch <- timer.next:
		default:
		}
		if timer.period == 0 {
			delete(c.timers, timer)
		} else {
			timer.next = timer.next.Add((c.now.Sub(timer.next)/timer.period + 1) * timer.period)
		}
	}
}
func (t *fakeTimer) C() <-chan time.Time { return t.ch }
func (t *fakeTimer) Stop()               { t.clock.mu.Lock(); defer t.clock.mu.Unlock(); delete(t.clock.timers, t) }
