package main

import (
	"errors"
	"orbital/colony/labkit"
	"sync"
	"time"
)

type LaunchGate struct {
	mu                sync.Mutex
	clock             labkit.Clock
	capacity, credits int
	every             time.Duration
	last              time.Time
}

func NewLaunchGate(clock labkit.Clock, capacity int, every time.Duration) (*LaunchGate, error) {
	if clock == nil || capacity <= 0 || every <= 0 {
		return nil, errors.New("invalid launch gate")
	}
	return &LaunchGate{clock: clock, capacity: capacity, credits: capacity, every: every, last: clock.Now()}, nil
}
func (g *LaunchGate) Take(n int) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if n <= 0 {
		return false, errors.New("positive request required")
	}
	now := g.clock.Now()
	if now.After(g.last) {
		elapsed := now.Sub(g.last) / g.every
		if elapsed > 0 {
			earned := g.capacity - g.credits
			if elapsed < time.Duration(earned) {
				earned = int(elapsed)
			}
			g.credits += earned
			g.last = g.last.Add(elapsed * g.every)
		}
	}
	if n > g.credits {
		return false, nil
	}
	g.credits -= n
	return true, nil
}
func main() {}
