package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Beacon struct {
	ID     string
	Charge int
}
type BeaconStore interface {
	Load(context.Context, string) (Beacon, error)
	Save(context.Context, Beacon) error
}

var ErrBeacon = errors.New("invalid beacon")
var ErrCacheSync = errors.New("saved, cache invalidation failed")

type BeaconCache struct {
	mu    sync.Mutex
	redis *redis.Client
	store BeaconStore
	ttl   time.Duration
	dirty map[string]bool
}

func NewBeaconCache(client *redis.Client, store BeaconStore, ttl time.Duration) *BeaconCache {
	return &BeaconCache{redis: client, store: store, ttl: ttl, dirty: make(map[string]bool)}
}
func (c *BeaconCache) Read(ctx context.Context, id string) (Beacon, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return Beacon{}, err
	}
	if id == "" {
		return Beacon{}, ErrBeacon
	}
	key := "colony:beacon:" + id
	if !c.dirty[id] {
		raw, err := c.redis.Get(ctx, key).Bytes()
		if ctx.Err() != nil {
			return Beacon{}, ctx.Err()
		}
		if err == nil {
			var b Beacon
			if json.Unmarshal(raw, &b) == nil && b.ID == id && b.Charge >= 0 {
				return b, nil
			}
		} else if !errors.Is(err, redis.Nil) {
			return c.store.Load(ctx, id)
		}
	}
	return c.refill(ctx, id)
}
func (c *BeaconCache) refill(ctx context.Context, id string) (Beacon, error) {
	b, err := c.store.Load(ctx, id)
	if err != nil {
		return Beacon{}, err
	}
	key := "colony:beacon:" + id
	if c.dirty[id] {
		if c.redis.Del(ctx, key).Err() != nil {
			return b, nil
		}
		delete(c.dirty, id)
	}
	raw, err := json.Marshal(b)
	if err != nil {
		return Beacon{}, err
	}
	// The store remains authoritative when the best-effort cache fill fails.
	c.redis.Set(ctx, key, raw, c.ttl).Err()
	return b, nil
}
func (c *BeaconCache) Write(ctx context.Context, b Beacon) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return err
	}
	if b.ID == "" || b.Charge < 0 {
		return ErrBeacon
	}
	if err := c.store.Save(ctx, b); err != nil {
		return err
	}
	if err := c.redis.Del(ctx, "colony:beacon:"+b.ID).Err(); err != nil {
		c.dirty[b.ID] = true
		return fmt.Errorf("%w: %w", ErrCacheSync, err)
	}
	delete(c.dirty, b.ID)
	return nil
}
func main() {}
