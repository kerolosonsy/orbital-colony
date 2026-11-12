package main

import (
	"context"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/redis/go-redis/v9"
)

// Safety limits: this lock is advisory. A holder paused longer than its TTL
// (GC, swap, a slow network) still believes it holds the lock after Redis has
// expired it and handed it to someone else. The token only stops it from
// deleting or extending the new holder's lock; it cannot stop its writes.
// Protecting the resource itself needs fencing tokens: a number that grows on
// every acquisition and that the storage checks, rejecting older numbers.

var ErrLockKey = errors.New("invalid lock request")
var ErrHeld = errors.New("lock held")
var ErrNotOwner = errors.New("not lock owner")

type Lease struct {
	Key   string
	Token string
}

var releaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0`)

var extendScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0`)

type Locker struct {
	client *redis.Client
	random io.Reader
}

func NewLocker(client *redis.Client, random io.Reader) *Locker {
	return &Locker{client: client, random: random}
}

func lockKey(key string) string { return "colony:lock:" + key }

func (l *Locker) Acquire(ctx context.Context, key string, ttl time.Duration) (Lease, error) {
	if key == "" || ttl <= 0 {
		return Lease{}, ErrLockKey
	}
	raw := make([]byte, 16)
	if _, err := io.ReadFull(l.random, raw); err != nil {
		return Lease{}, err
	}
	token := hex.EncodeToString(raw)
	ok, err := l.client.SetNX(ctx, lockKey(key), token, ttl).Result()
	if err != nil {
		return Lease{}, err
	}
	if !ok {
		return Lease{}, ErrHeld
	}
	return Lease{Key: key, Token: token}, nil
}

func (l *Locker) Release(ctx context.Context, lease Lease) error {
	n, err := releaseScript.Run(ctx, l.client, []string{lockKey(lease.Key)}, lease.Token).Int()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotOwner
	}
	return nil
}

func (l *Locker) Extend(ctx context.Context, lease Lease, ttl time.Duration) error {
	if ttl <= 0 {
		return ErrLockKey
	}
	n, err := extendScript.Run(ctx, l.client, []string{lockKey(lease.Key)}, lease.Token, ttl.Milliseconds()).Int()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotOwner
	}
	return nil
}

func main() {}
