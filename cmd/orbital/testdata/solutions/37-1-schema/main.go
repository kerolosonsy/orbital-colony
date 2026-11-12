package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Migration struct {
	Version  int
	Up, Down string
}

var Migrations = []Migration{
	{1, `CREATE TABLE users (
 id TEXT PRIMARY KEY,
 email TEXT NOT NULL UNIQUE CHECK (length(email) > 0),
 password_hash TEXT NOT NULL
)`, `DROP TABLE users`},
	{2, `CREATE TABLE crews (
 id TEXT PRIMARY KEY,
 name TEXT NOT NULL CHECK (length(trim(name)) > 0)
)`, `DROP TABLE crews`},
	{3, `CREATE TABLE tasks (
 id TEXT PRIMARY KEY,
 crew_id TEXT NOT NULL REFERENCES crews(id) ON DELETE CASCADE,
 title TEXT NOT NULL CHECK (length(trim(title)) > 0),
 done BOOLEAN NOT NULL DEFAULT FALSE,
 version BIGINT NOT NULL CHECK (version >= 1)
)`, `DROP TABLE tasks`},
}

type User struct{ ID, Email, PasswordHash string }
type Crew struct{ ID, Name string }
type Task struct {
	ID      string `json:"id"`
	CrewID  string `json:"crew_id"`
	Title   string `json:"title"`
	Done    bool   `json:"done"`
	Version int64  `json:"version"`
}

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Repository interface {
	CreateUser(ctx context.Context, u User) error
	UserByEmail(ctx context.Context, email string) (User, error)
	CreateCrew(ctx context.Context, c Crew) error
	CreateTask(ctx context.Context, t Task) error
	GetTask(ctx context.Context, crewID, id string) (Task, error)
}

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func SchemaVersion(ctx context.Context, db *sql.DB) (int, error) {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY)`); err != nil {
		return 0, err
	}
	var v int
	err := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&v)
	return v, err
}

func Migrate(ctx context.Context, db *sql.DB, target int) error {
	if target < 0 || target > len(Migrations) {
		return fmt.Errorf("unknown target %d", target)
	}
	for i, m := range Migrations {
		if m.Version != i+1 {
			return fmt.Errorf("migration %d out of order", m.Version)
		}
	}
	current, err := SchemaVersion(ctx, db)
	if err != nil {
		return err
	}
	if current > len(Migrations) {
		return fmt.Errorf("schema %d is newer than this binary", current)
	}
	for current != target {
		var m Migration
		var query, record string
		if current < target {
			m = Migrations[current]
			query, record = m.Up, `INSERT INTO schema_migrations(version) VALUES ($1)`
		} else {
			m = Migrations[current-1]
			query, record = m.Down, `DELETE FROM schema_migrations WHERE version = $1`
		}
		if err := step(ctx, db, query, record, m.Version); err != nil {
			return fmt.Errorf("migration %d: %w", m.Version, err)
		}
		if current < target {
			current++
		} else {
			current--
		}
	}
	return nil
}

func step(ctx context.Context, db *sql.DB, query, record string, version int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, record, version); err != nil {
		return err
	}
	return tx.Commit()
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func insertOnce(ctx context.Context, q execer, query string, args ...any) error {
	res, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrConflict
	}
	return nil
}

func normalEmail(email string) string { return strings.ToLower(strings.TrimSpace(email)) }

func (s *Store) CreateUser(ctx context.Context, u User) error {
	return insertOnce(ctx, s.db, `INSERT INTO users(id, email, password_hash) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		u.ID, normalEmail(u.Email), u.PasswordHash)
}

func (s *Store) UserByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := s.db.QueryRowContext(ctx, `SELECT id, email, password_hash FROM users WHERE email = $1`, normalEmail(email)).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrNotFound
	}
	return u, err
}

func (s *Store) CreateCrew(ctx context.Context, c Crew) error {
	return insertOnce(ctx, s.db, `INSERT INTO crews(id, name) VALUES ($1, $2) ON CONFLICT DO NOTHING`, c.ID, c.Name)
}

func (s *Store) CreateTask(ctx context.Context, t Task) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var one int
	err = tx.QueryRowContext(ctx, `SELECT 1 FROM crews WHERE id = $1`, t.CrewID).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if err := insertOnce(ctx, tx, `INSERT INTO tasks(id, crew_id, title, done, version) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
		t.ID, t.CrewID, t.Title, t.Done, t.Version); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) GetTask(ctx context.Context, crewID, id string) (Task, error) {
	var t Task
	err := s.db.QueryRowContext(ctx, `SELECT id, crew_id, title, done, version FROM tasks WHERE crew_id = $1 AND id = $2`, crewID, id).
		Scan(&t.ID, &t.CrewID, &t.Title, &t.Done, &t.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrNotFound
	}
	return t, err
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func NewServer(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		v, err := SchemaVersion(r.Context(), s.db)
		if err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "unavailable"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]int{"schema": v})
	})
	mux.HandleFunc("GET /v1/crews/{crew}/tasks/{task}", func(w http.ResponseWriter, r *http.Request) {
		t, err := s.GetTask(r.Context(), r.PathValue("crew"), r.PathValue("task"))
		switch {
		case errors.Is(err, ErrNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
		case err != nil:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
		default:
			writeJSON(w, http.StatusOK, t)
		}
	})
	return mux
}

func main() {}
