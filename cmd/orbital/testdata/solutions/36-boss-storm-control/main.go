package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

type StormEvent struct {
	ID      string
	Delta   int
	Version int
}
type StormBroker interface {
	Publish(context.Context, StormEvent) error
}

var ErrStorm = errors.New("invalid or conflicting storm command")
var ErrStormCorrupt = errors.New("corrupt storm command")

type StormControl struct {
	db     *sql.DB
	broker StormBroker
	log    *slog.Logger
	mu     sync.Mutex
}

func NewStormControl(db *sql.DB, broker StormBroker, log *slog.Logger) *StormControl {
	return &StormControl{db: db, broker: broker, log: log}
}
func (s *StormControl) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS storm_commands(seq INTEGER PRIMARY KEY AUTOINCREMENT, id TEXT NOT NULL UNIQUE, delta INTEGER NOT NULL, checksum TEXT NOT NULL, published INTEGER NOT NULL DEFAULT 0)")
	return err
}
func stormChecksum(id string, delta int) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s:%d", id, delta))))
}
func (s *StormControl) Record(ctx context.Context, id string, delta int) error {
	if id == "" || delta == 0 {
		return ErrStorm
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var old int
	err = tx.QueryRowContext(ctx, "SELECT delta FROM storm_commands WHERE id=?", id).Scan(&old)
	if err == nil {
		if old != delta {
			return ErrStorm
		}
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO storm_commands(id,delta,checksum) VALUES(?,?,?)", id, delta, stormChecksum(id, delta)); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *StormControl) Flush(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	rows, err := s.db.QueryContext(ctx, "SELECT id,delta,checksum FROM storm_commands WHERE published=0 ORDER BY seq")
	if err != nil {
		return err
	}
	var events []StormEvent
	for rows.Next() {
		var e StormEvent
		var checksum string
		if err = rows.Scan(&e.ID, &e.Delta, &checksum); err != nil {
			rows.Close()
			return err
		}
		if e.ID == "" || e.Delta == 0 || checksum != stormChecksum(e.ID, e.Delta) {
			rows.Close()
			return ErrStormCorrupt
		}
		e.Version = 1
		events = append(events, e)
	}
	err = rows.Err()
	closeErr := rows.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	for _, e := range events {
		if err = ctx.Err(); err != nil {
			return err
		}
		if err = s.broker.Publish(ctx, e); err != nil {
			s.log.Error("storm publish failed")
			return err
		}
		if _, err = s.db.ExecContext(ctx, "UPDATE storm_commands SET published=1 WHERE id=?", e.ID); err != nil {
			return err
		}
	}
	return nil
}
func (s *StormControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if r.URL.Path == "/livez" {
		fmt.Fprintln(w, "live")
		return
	}
	if r.URL.Path != "/readyz" && r.URL.Path != "/metrics" {
		http.NotFound(w, r)
		return
	}
	var total, pending int
	err := s.db.QueryRowContext(r.Context(), "SELECT count(*),coalesce(sum(CASE WHEN published=0 THEN 1 ELSE 0 END),0) FROM storm_commands").Scan(&total, &pending)
	if err != nil {
		s.log.Error("storm store unavailable")
		http.Error(w, "unavailable", 503)
		return
	}
	if r.URL.Path == "/readyz" {
		if pending != 0 {
			http.Error(w, "pending", 503)
			return
		}
		fmt.Fprintln(w, "ready")
		return
	}
	fmt.Fprintf(w, "colony_storm_commands_total %d\ncolony_storm_pending %d\n", total, pending)
}
func main() {}
