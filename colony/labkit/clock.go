// Package labkit provides deterministic tools for player lab tests.
package labkit

import (
	"sync"
	"time"
)

// Ticker is the part of *time.Ticker a Clock hands out: a channel and Stop.
type Ticker interface {
	C() <-chan time.Time
	Stop()
}

// Clock is what code under test asks for the time, so tests can swap in a
// FakeClock and never sleep.
type Clock interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
	NewTicker(time.Duration) Ticker
}

// RealClock is the production Clock backed by the time package.
type RealClock struct{}

// Now returns time.Now().
func (RealClock) Now() time.Time { return time.Now() }

// After returns time.After(d).
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// NewTicker wraps time.NewTicker(d).
func (RealClock) NewTicker(d time.Duration) Ticker { return realTicker{time.NewTicker(d)} }

type realTicker struct{ *time.Ticker }

// C returns the ticker channel.
func (t realTicker) C() <-chan time.Time { return t.Ticker.C }

// FakeClock is a test Clock that only moves when Advance is called.
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
func NewFakeClock(start time.Time) *FakeClock { return &FakeClock{now: start} }

// Now returns the fake current time.
func (c *FakeClock) Now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.now }

// After returns a channel that receives once the clock is advanced by d.
func (c *FakeClock) After(d time.Duration) <-chan time.Time { return c.timer(d, 0).ch }

// NewTicker returns a fake ticker firing every d of advanced time; d must be positive.
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

// Advance moves the clock forward by d and delivers every timer and tick now due.
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

// C returns the fake timer channel.
func (t *fakeTimer) C() <-chan time.Time { return t.ch }

// Stop cancels the fake timer or ticker.
func (t *fakeTimer) Stop() { t.clock.mu.Lock(); defer t.clock.mu.Unlock(); delete(t.clock.timers, t) }
