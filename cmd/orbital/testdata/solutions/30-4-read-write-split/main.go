package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrRoute = errors.New("invalid route")
var ErrReplicaUnavailable = errors.New("replica unavailable")

type RelayStore interface {
	Read(context.Context, string) (string, error)
	Write(context.Context, string, string) error
}
type RelayRouter struct {
	mu               sync.Mutex
	primary, replica RelayStore
	now              func() time.Time
	sticky           time.Duration
	until            map[string]time.Time
}

func NewRelayRouter(primary, replica RelayStore, now func() time.Time, sticky time.Duration) *RelayRouter {
	return &RelayRouter{primary: primary, replica: replica, now: now, sticky: sticky, until: make(map[string]time.Time)}
}
func (r *RelayRouter) Write(ctx context.Context, session, key, content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if session == "" || key == "" {
		return ErrRoute
	}
	if err := r.primary.Write(ctx, key, content); err != nil {
		return err
	}
	r.until[session] = r.now().Add(r.sticky)
	return nil
}
func relayRead(ctx context.Context, store RelayStore, key string) (string, error) {
	content, err := store.Read(ctx, key)
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", err
	}
	return content, nil
}
func (r *RelayRouter) Read(ctx context.Context, session, key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if session == "" || key == "" {
		return "", ErrRoute
	}
	if until, ok := r.until[session]; ok {
		if r.now().Before(until) {
			return relayRead(ctx, r.primary, key)
		}
		delete(r.until, session)
	}
	content, err := relayRead(ctx, r.replica, key)
	if err == nil {
		return content, nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || !errors.Is(err, ErrReplicaUnavailable) {
		return "", err
	}
	return relayRead(ctx, r.primary, key)
}
func main() {}
