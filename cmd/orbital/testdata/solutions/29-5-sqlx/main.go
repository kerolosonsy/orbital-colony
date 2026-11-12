package main

import (
	"context"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Ration struct {
	Bay   string `db:"bay"`
	Units int    `db:"units"`
}
type RationStore interface {
	Put(context.Context, Ration) error
	List(context.Context, int) ([]Ration, error)
}

var ErrRation = errors.New("invalid ration")

func StockRation(ctx context.Context, store RationStore, r Ration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.Bay = strings.TrimSpace(r.Bay)
	if r.Bay == "" || r.Units < 0 {
		return ErrRation
	}
	return store.Put(ctx, r)
}

type SQLRations struct{ db *sqlx.DB }

func NewSQLRations(db *sqlx.DB) *SQLRations { return &SQLRations{db: db} }
func InitRations(ctx context.Context, db *sqlx.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS ration_bins(bay TEXT PRIMARY KEY CHECK(length(trim(bay))>0),units INTEGER NOT NULL CHECK(units>=0))`)
	return err
}
func (s *SQLRations) Put(ctx context.Context, r Ration) error {
	_, err := sqlx.NamedExecContext(ctx, s.db, "INSERT INTO ration_bins(bay,units) VALUES(:bay,:units)", r)
	return err
}
func (s *SQLRations) List(ctx context.Context, min int) ([]Ration, error) {
	out := []Ration{}
	err := sqlx.SelectContext(ctx, s.db, &out, "SELECT bay,units FROM ration_bins WHERE units >= $1 ORDER BY bay", min)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func main() {}
