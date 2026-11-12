package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"orbital/colony/labkit"
)

var (
	ErrForced   = errors.New("drain deadline passed")
	ErrGateUsed = errors.New("gate already ran")
)

type Gate struct {
	h       http.Handler
	clock   labkit.Clock
	grace   time.Duration
	release func() error
	used    atomic.Bool
	drain   chan struct{}
	once    sync.Once
}

func NewGate(h http.Handler, clock labkit.Clock, grace time.Duration, release func() error) (*Gate, error) {
	if h == nil || clock == nil || release == nil || grace <= 0 {
		return nil, errors.New("gate needs a handler, a clock, a release, and a positive grace")
	}
	return &Gate{h: h, clock: clock, grace: grace, release: release, drain: make(chan struct{})}, nil
}

func (g *Gate) Draining() <-chan struct{} { return g.drain }

func (g *Gate) Run(ctx context.Context, ln net.Listener) error {
	if !g.used.CompareAndSwap(false, true) {
		return ErrGateUsed
	}
	srv := &http.Server{Handler: g.h, ReadHeaderTimeout: 2 * time.Second, IdleTimeout: 30 * time.Second}
	srv.RegisterOnShutdown(func() { g.once.Do(func() { close(g.drain) }) })
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	select {
	case err := <-served:
		return errors.Join(err, g.release())
	case <-ctx.Done():
	}
	result := g.shutdown(srv)
	if err := <-served; !errors.Is(err, http.ErrServerClosed) {
		result = errors.Join(result, err)
	}
	return errors.Join(result, g.release())
}

// shutdown drains with a fresh context that the injected clock cancels after grace.
func (g *Gate) shutdown(srv *http.Server) error {
	drainCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	finished := make(chan struct{})
	defer close(finished)
	go func() {
		select {
		case <-g.clock.After(g.grace):
			cancel()
		case <-finished:
		}
	}()
	if err := srv.Shutdown(drainCtx); err != nil {
		srv.Close()
		return fmt.Errorf("%w: %v", ErrForced, err)
	}
	return nil
}

func main() {}
