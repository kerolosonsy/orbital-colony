package main

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Task struct {
	ID       string
	Body     string
	Attempts int
}

var ErrTask = errors.New("invalid task")
var ErrEmpty = errors.New("queue empty")
var ErrLease = errors.New("lease not held")

type Queue struct {
	client *redis.Client
	name   string
	lease  time.Duration
	max    int
	now    func() time.Time
}

func NewQueue(client *redis.Client, name string, lease time.Duration, maxAttempts int, now func() time.Time) *Queue {
	return &Queue{client: client, name: name, lease: lease, max: maxAttempts, now: now}
}

func (q *Queue) key(part string) string { return q.name + ":" + part }

// lease names one claim: a stale worker holds an older attempt count.
func lease(t Task) string { return t.ID + "#" + strconv.Itoa(t.Attempts) }

func (q *Queue) Enqueue(ctx context.Context, t Task) (bool, error) {
	if t.ID == "" {
		return false, ErrTask
	}
	added, err := q.client.SAdd(ctx, q.key("seen"), t.ID).Result()
	if err != nil || added == 0 {
		return false, err
	}
	t.Attempts = 0
	if err := q.save(ctx, t); err != nil {
		return false, err
	}
	return true, q.client.RPush(ctx, q.key("ready"), t.ID).Err()
}

func (q *Queue) save(ctx context.Context, t Task) error {
	raw, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return q.client.Set(ctx, q.key("task:"+t.ID), raw, 0).Err()
}

func (q *Queue) load(ctx context.Context, id string) (Task, error) {
	raw, err := q.client.Get(ctx, q.key("task:"+id)).Bytes()
	if err != nil {
		return Task{}, err
	}
	var t Task
	err = json.Unmarshal(raw, &t)
	return t, err
}

// retry records one failed attempt, then sends the task back or to the dead list.
func (q *Queue) retry(ctx context.Context, id string) error {
	t, err := q.load(ctx, id)
	if err != nil {
		return err
	}
	t.Attempts++
	if err := q.save(ctx, t); err != nil {
		return err
	}
	if t.Attempts >= q.max {
		return q.client.RPush(ctx, q.key("dead"), id).Err()
	}
	return q.client.RPush(ctx, q.key("ready"), id).Err()
}

func (q *Queue) Claim(ctx context.Context) (Task, error) {
	now := q.now().UnixMilli()
	expired, err := q.client.ZRangeByScore(ctx, q.key("leased"), &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(now, 10)}).Result()
	if err != nil {
		return Task{}, err
	}
	for _, member := range expired {
		removed, err := q.client.ZRem(ctx, q.key("leased"), member).Result()
		if err != nil {
			return Task{}, err
		}
		if removed == 1 {
			id := member[:strings.LastIndex(member, "#")]
			if err := q.retry(ctx, id); err != nil {
				return Task{}, err
			}
		}
	}
	id, err := q.client.LPop(ctx, q.key("ready")).Result()
	if errors.Is(err, redis.Nil) {
		return Task{}, ErrEmpty
	}
	if err != nil {
		return Task{}, err
	}
	t, err := q.load(ctx, id)
	if err != nil {
		return Task{}, err
	}
	deadline := q.now().Add(q.lease).UnixMilli()
	if err := q.client.ZAdd(ctx, q.key("leased"), redis.Z{Score: float64(deadline), Member: lease(t)}).Err(); err != nil {
		return Task{}, err
	}
	return t, nil
}

func (q *Queue) release(ctx context.Context, t Task) error {
	removed, err := q.client.ZRem(ctx, q.key("leased"), lease(t)).Result()
	if err != nil {
		return err
	}
	if removed == 0 {
		return ErrLease
	}
	return nil
}

func (q *Queue) Ack(ctx context.Context, t Task) error {
	if err := q.release(ctx, t); err != nil {
		return err
	}
	return q.client.Del(ctx, q.key("task:"+t.ID)).Err()
}

func (q *Queue) Fail(ctx context.Context, t Task) error {
	if err := q.release(ctx, t); err != nil {
		return err
	}
	return q.retry(ctx, t.ID)
}

func (q *Queue) Dead(ctx context.Context) ([]Task, error) {
	ids, err := q.client.LRange(ctx, q.key("dead"), 0, -1).Result()
	if err != nil {
		return nil, err
	}
	out := []Task{}
	for _, id := range ids {
		t, err := q.load(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (q *Queue) Redrive(ctx context.Context) (int, error) {
	n := 0
	for {
		id, err := q.client.LPop(ctx, q.key("dead")).Result()
		if errors.Is(err, redis.Nil) {
			return n, nil
		}
		if err != nil {
			return n, err
		}
		t, err := q.load(ctx, id)
		if err != nil {
			return n, err
		}
		t.Attempts = 0
		if err := q.save(ctx, t); err != nil {
			return n, err
		}
		if err := q.client.RPush(ctx, q.key("ready"), id).Err(); err != nil {
			return n, err
		}
		n++
	}
}

func main() {}
