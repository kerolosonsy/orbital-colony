package main

import (
	"errors"
	"orbital/colony/labkit"
	"sync"
	"time"
)

type expiryReading struct {
	text    string
	expires time.Time
}
type ExpiryBoard struct {
	mu         sync.RWMutex
	clock      labkit.Clock
	readings   map[string]expiryReading
	closed     bool
	stop, done chan struct{}
	sweeps     chan int
}

func NewExpiryBoard(clock labkit.Clock, sweep time.Duration) (*ExpiryBoard, error) {
	if clock == nil || sweep <= 0 {
		return nil, errors.New("invalid expiry board")
	}
	b := &ExpiryBoard{clock: clock, readings: make(map[string]expiryReading), stop: make(chan struct{}), done: make(chan struct{}), sweeps: make(chan int, 1)}
	ticker := clock.NewTicker(sweep)
	go b.run(ticker)
	return b, nil
}
func (b *ExpiryBoard) Put(key, text string, ttl time.Duration) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed || key == "" || ttl <= 0 {
		return errors.New("invalid reading")
	}
	b.readings[key] = expiryReading{text, b.clock.Now().Add(ttl)}
	return nil
}
func (b *ExpiryBoard) Lookup(key string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	r, ok := b.readings[key]
	if b.closed || !ok || !b.clock.Now().Before(r.expires) {
		return "", false
	}
	return r.text, true
}
func (b *ExpiryBoard) Stored() int        { b.mu.RLock(); defer b.mu.RUnlock(); return len(b.readings) }
func (b *ExpiryBoard) Sweeps() <-chan int { return b.sweeps }
func (b *ExpiryBoard) sweep() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.clock.Now()
	removed := 0
	for key, r := range b.readings {
		if !now.Before(r.expires) {
			delete(b.readings, key)
			removed++
		}
	}
	return removed
}
func (b *ExpiryBoard) run(ticker labkit.Ticker) {
	defer close(b.done)
	defer close(b.sweeps)
	defer ticker.Stop()
	for {
		select {
		case <-b.stop:
			return
		case <-ticker.C():
			removed := b.sweep()
			select {
			case b.sweeps <- removed:
			default:
			}
		}
	}
}
func (b *ExpiryBoard) Close() {
	b.mu.Lock()
	if !b.closed {
		b.closed = true
		clear(b.readings)
		close(b.stop)
	}
	b.mu.Unlock()
	<-b.done
}
func main() {}
