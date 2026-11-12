package main

import (
	"context"
	"database/sql"
	"errors"
)

type Crate struct {
	ID    int64
	Label string
	Units int
}

var ErrCrateMissing = errors.New("crate missing")

func CrateError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrCrateMissing
	}
	return err
}
func InitCrates(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS crates(id INTEGER PRIMARY KEY CHECK(id>0),label TEXT NOT NULL UNIQUE CHECK(length(trim(label))>0),units INTEGER NOT NULL CHECK(units>=0))`)
	return err
}
func PutCrate(ctx context.Context, db *sql.DB, c Crate) error {
	_, err := db.ExecContext(ctx, "INSERT INTO crates(id,label,units) VALUES(?,?,?)", c.ID, c.Label, c.Units)
	return err
}
func GetCrate(ctx context.Context, db *sql.DB, id int64) (Crate, error) {
	var c Crate
	err := db.QueryRowContext(ctx, "SELECT id,label,units FROM crates WHERE id=?", id).Scan(&c.ID, &c.Label, &c.Units)
	return c, CrateError(err)
}
func ListCrates(ctx context.Context, db *sql.DB) ([]Crate, error) {
	rows, err := db.QueryContext(ctx, "SELECT id,label,units FROM crates ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Crate{}
	for rows.Next() {
		var c Crate
		if err := rows.Scan(&c.ID, &c.Label, &c.Units); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}
func crateChanged(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrCrateMissing
	}
	return nil
}
func ResizeCrate(ctx context.Context, db *sql.DB, id int64, units int) error {
	return crateChanged(db.ExecContext(ctx, "UPDATE crates SET units=? WHERE id=?", units, id))
}
func DeleteCrate(ctx context.Context, db *sql.DB, id int64) error {
	return crateChanged(db.ExecContext(ctx, "DELETE FROM crates WHERE id=?", id))
}
func main() {}
