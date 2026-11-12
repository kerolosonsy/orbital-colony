package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SeedLine struct {
	Name  string
	Units int
}
type NurseryCache interface {
	Invalidate(context.Context, string) error
}

var ErrPlanting = errors.New("invalid planting")
var ErrStock = errors.New("insufficient stock")
var ErrReplay = errors.New("planting already booked")
var ErrCacheSync = errors.New("committed but cache invalidation failed")

func InitNursery(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS nursery_stock(name TEXT PRIMARY KEY, units INTEGER NOT NULL CHECK(units>=0)); CREATE TABLE IF NOT EXISTS nursery_bookings(request TEXT PRIMARY KEY)")
	return err
}
func BookPlanting(ctx context.Context, db *sql.DB, cache NurseryCache, request string, lines []SeedLine) error {
	if request == "" || len(lines) == 0 {
		return ErrPlanting
	}
	seen := map[string]bool{}
	for _, line := range lines {
		if line.Name == "" || line.Units <= 0 || seen[line.Name] {
			return ErrPlanting
		}
		seen[line.Name] = true
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var found string
	err = tx.QueryRowContext(ctx, "SELECT request FROM nursery_bookings WHERE request=?", request).Scan(&found)
	if err == nil {
		return ErrReplay
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	for _, line := range lines {
		result, err := tx.ExecContext(ctx, "UPDATE nursery_stock SET units=units-? WHERE name=? AND units>=?", line.Units, line.Name, line.Units)
		if err != nil {
			return err
		}
		n, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return ErrStock
		}
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO nursery_bookings(request) VALUES(?)", request); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	for _, line := range lines {
		if err = cache.Invalidate(ctx, line.Name); err != nil {
			return fmt.Errorf("%w: %w", ErrCacheSync, err)
		}
	}
	return nil
}
func main() {}
