package main

import (
	"context"
	"database/sql"
	"errors"
)

func ConfigureRelayPool(db *sql.DB, max int) error {
	if max <= 0 {
		return errors.New("positive pool limit required")
	}
	db.SetMaxOpenConns(max)
	db.SetMaxIdleConns(max)
	return nil
}
func ReadRelayPorts(ctx context.Context, db *sql.DB, min int) ([]int, error) {
	rows, err := db.QueryContext(ctx, "SELECT code FROM relay_ports WHERE code>=? ORDER BY code", min)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []int{}
	for rows.Next() {
		var code int
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		out = append(out, code)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
func main() {}
