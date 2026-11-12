package main

import (
	"context"
	"database/sql"
	"errors"
	"math"
)

var ErrFuel = errors.New("unavailable oxygen")
var ErrReplay = errors.New("request already applied")
var ErrTransfer = errors.New("invalid transfer")

type OxygenMove struct {
	Request, From, To string
	Units             int64
}

func InitOxygen(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS oxygen_tanks(name TEXT PRIMARY KEY, units INTEGER NOT NULL CHECK(units>=0)); CREATE TABLE IF NOT EXISTS oxygen_moves(request TEXT PRIMARY KEY)`)
	return err
}
func oxygenUpdate(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return ErrFuel
	}
	return nil
}
func MoveOxygen(ctx context.Context, db *sql.DB, m OxygenMove, afterDebit func() error) error {
	if m.Request == "" || m.From == "" || m.To == "" || m.From == m.To || m.Units <= 0 {
		return ErrTransfer
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var request string
	err = tx.QueryRowContext(ctx, "SELECT request FROM oxygen_moves WHERE request=?", m.Request).Scan(&request)
	if err == nil {
		return ErrReplay
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err := oxygenUpdate(ctx, tx, "UPDATE oxygen_tanks SET units=units-? WHERE name=? AND units>=?", m.Units, m.From, m.Units); err != nil {
		return err
	}
	if afterDebit != nil {
		if err := afterDebit(); err != nil {
			return err
		}
	}
	if err := oxygenUpdate(ctx, tx, "UPDATE oxygen_tanks SET units=units+? WHERE name=? AND units<=?", m.Units, m.To, math.MaxInt64-m.Units); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO oxygen_moves(request) VALUES(?)", m.Request); err != nil {
		return err
	}
	return tx.Commit()
}
func main() {}
