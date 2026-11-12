package main

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

type Migration struct {
	Version  int
	Up, Down string
}

func migrationVersions(ctx context.Context, db *sql.DB) ([]int, error) {
	rows, err := db.QueryContext(ctx, "SELECT version FROM habitat_versions ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []int
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
func migrationSQL(ctx context.Context, tx *sql.Tx, m Migration, down bool) error {
	body := m.Up
	if down {
		body = m.Down
	}
	if _, err := tx.ExecContext(ctx, body); err != nil {
		return err
	}
	if down {
		_, err := tx.ExecContext(ctx, "DELETE FROM habitat_versions WHERE version=?", m.Version)
		return err
	}
	_, err := tx.ExecContext(ctx, "INSERT INTO habitat_versions(version) VALUES(?)", m.Version)
	return err
}
func migrationStep(ctx context.Context, db *sql.DB, m Migration, down bool) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := migrationSQL(ctx, tx, m, down); err != nil {
		return err
	}
	return tx.Commit()
}
func MigrateHabitat(ctx context.Context, db *sql.DB, source []Migration, down bool) error {
	last := 0
	for _, m := range source {
		if m.Version <= last || strings.TrimSpace(m.Up) == "" || strings.TrimSpace(m.Down) == "" {
			return errors.New("invalid migration source")
		}
		last = m.Version
	}
	if _, err := db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS habitat_versions(version INTEGER PRIMARY KEY)"); err != nil {
		return err
	}
	applied, err := migrationVersions(ctx, db)
	if err != nil {
		return err
	}
	if len(applied) > len(source) {
		return errors.New("unknown history")
	}
	for i, v := range applied {
		if source[i].Version != v {
			return errors.New("history is not a source prefix")
		}
	}
	if down {
		if len(applied) == 0 {
			return nil
		}
		return migrationStep(ctx, db, source[len(applied)-1], true)
	}
	for _, m := range source[len(applied):] {
		if err := migrationStep(ctx, db, m, false); err != nil {
			return err
		}
	}
	return nil
}
func main() {}
