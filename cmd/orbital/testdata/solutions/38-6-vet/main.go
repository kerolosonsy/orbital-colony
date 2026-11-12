package main

import (
	"context"
	"fmt"
	"sync"
)

type Gauge struct {
	mu    sync.Mutex
	total int
}

func (g *Gauge) Add(n int)  { g.mu.Lock(); defer g.mu.Unlock(); g.total += n }
func (g *Gauge) Total() int { g.mu.Lock(); defer g.mu.Unlock(); return g.total }
func Status(n int) string   { return fmt.Sprintf("readings=%d", n) }
func ReadWithCancel(parent context.Context, read func(context.Context) error) error {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()
	return read(ctx)
}
func main() {}
