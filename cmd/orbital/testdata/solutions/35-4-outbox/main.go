package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"

	_ "modernc.org/sqlite"
)

const SchemaVersion = 2

var ErrOrder = errors.New("invalid order")
var ErrVersion = errors.New("unknown schema version")

type Event struct {
	ID      int64
	Type    string
	Version int
	Payload string
}

type Publisher interface {
	Publish(ctx context.Context, e Event) error
}

type orderPayload struct {
	ID   string `json:"id"`
	Item string `json:"item"`
	Qty  int    `json:"qty"`
}

type Outbox struct{ db *sql.DB }

func NewOutbox(db *sql.DB) *Outbox { return &Outbox{db: db} }

func (o *Outbox) Init(ctx context.Context) error {
	_, err := o.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS orders (id TEXT PRIMARY KEY, item TEXT NOT NULL, qty INTEGER NOT NULL);
CREATE TABLE IF NOT EXISTS outbox (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL,
	version INTEGER NOT NULL,
	payload TEXT NOT NULL,
	published INTEGER NOT NULL DEFAULT 0
);`)
	return err
}

func (o *Outbox) PlaceOrder(ctx context.Context, id, item string, qty int) error {
	if id == "" || item == "" || qty <= 0 {
		return ErrOrder
	}
	payload, err := json.Marshal(orderPayload{ID: id, Item: item, Qty: qty})
	if err != nil {
		return err
	}
	tx, err := o.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `INSERT INTO orders (id, item, qty) VALUES (?, ?, ?)`, id, item, qty); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO outbox (type, version, payload) VALUES (?, ?, ?)`, "order.placed", SchemaVersion, string(payload)); err != nil {
		return err
	}
	return tx.Commit()
}

func (o *Outbox) Relay(ctx context.Context, pub Publisher, limit int) (int, error) {
	if limit <= 0 {
		return 0, nil
	}
	rows, err := o.db.QueryContext(ctx, `SELECT id, type, version, payload FROM outbox WHERE published = 0 ORDER BY id LIMIT ?`, limit)
	if err != nil {
		return 0, err
	}
	var pending []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Type, &e.Version, &e.Payload); err != nil {
			rows.Close()
			return 0, err
		}
		pending = append(pending, e)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	sent := 0
	for _, e := range pending {
		if err := pub.Publish(ctx, e); err != nil {
			return sent, err
		}
		if _, err := o.db.ExecContext(ctx, `UPDATE outbox SET published = 1 WHERE id = ?`, e.ID); err != nil {
			return sent, err
		}
		sent++
	}
	return sent, nil
}

type Consumer struct {
	mu     sync.Mutex
	seen   map[int64]bool
	totals map[string]int
}

func NewConsumer() *Consumer {
	return &Consumer{seen: map[int64]bool{}, totals: map[string]int{}}
}

func (c *Consumer) Handle(e Event) (bool, error) {
	if e.Version < 1 || e.Version > SchemaVersion {
		return false, ErrVersion
	}
	p := orderPayload{Qty: 1} // version 1 had no qty: one unit
	if err := json.Unmarshal([]byte(e.Payload), &p); err != nil {
		return false, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.seen[e.ID] {
		return false, nil
	}
	c.seen[e.ID] = true
	c.totals[p.Item] += p.Qty
	return true, nil
}

func (c *Consumer) Totals() map[string]int {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int, len(c.totals))
	for k, v := range c.totals {
		out[k] = v
	}
	return out
}

func main() {}
