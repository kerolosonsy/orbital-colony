package main

import (
	"errors"
	"orbital/colony/labkit"
	"sync"
	"time"
)

type PermitTank struct {
	mu               sync.Mutex
	clock            labkit.Clock
	capacity, credit int
	every            time.Duration
	last             time.Time
}

func NewPermitTank(clock labkit.Clock, capacity int, every time.Duration) (*PermitTank, error) {
	if clock == nil || capacity <= 0 || every <= 0 {
		return nil, errors.New("invalid permit tank")
	}
	return &PermitTank{clock: clock, capacity: capacity, credit: capacity, every: every, last: clock.Now()}, nil
}
func (p *PermitTank) refill() {
	now := p.clock.Now()
	if !now.After(p.last) {
		return
	}
	steps := now.Sub(p.last) / p.every
	if steps >= time.Duration(p.capacity-p.credit) {
		p.credit = p.capacity
		p.last = now
		return
	}
	p.credit += int(steps)
	p.last = p.last.Add(steps * p.every)
}
func (p *PermitTank) Admit(cost int) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if cost <= 0 {
		return false, errors.New("positive cost required")
	}
	p.refill()
	if cost > p.credit {
		return false, nil
	}
	p.credit -= cost
	return true, nil
}
func main() {}
