package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrForbidden   = errors.New("forbidden")
	ErrInvalidRole = errors.New("invalid role")
)

type AuditEntry struct {
	Actor  string `json:"actor"`
	Action string `json:"action"`
	Target string `json:"target"`
}

type Crews struct{ Auth *Auth }

var rank = map[string]int{"member": 1, "admin": 2, "owner": 3}

func (c *Crews) db() *sql.DB { return c.Auth.Store.db }

func (c *Crews) Role(ctx context.Context, crewID, userID string) (string, error) {
	var role string
	err := c.db().QueryRowContext(ctx, `SELECT role FROM memberships WHERE crew_id = $1 AND user_id = $2`, crewID, userID).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	}
	return role, err
}

func audit(ctx context.Context, tx *sql.Tx, now int64, crewID, actor, action, target string) error {
	var seq int64
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(seq), 0) + 1 FROM audit WHERE crew_id = $1`, crewID).Scan(&seq); err != nil {
		return err
	}
	_, err := tx.ExecContext(ctx, `INSERT INTO audit(crew_id, seq, actor_id, action, target, at) VALUES ($1, $2, $3, $4, $5, $6)`,
		crewID, seq, actor, action, target, now)
	return err
}

func (c *Crews) CreateCrew(ctx context.Context, actor User, name string) (Crew, error) {
	id, err := c.Auth.randomText(12)
	if err != nil {
		return Crew{}, err
	}
	tx, err := c.db().BeginTx(ctx, nil)
	if err != nil {
		return Crew{}, err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `INSERT INTO crews(id, name) VALUES ($1, $2)`, id, name); err != nil {
		return Crew{}, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO memberships(crew_id, user_id, role) VALUES ($1, $2, 'owner')`, id, actor.ID); err != nil {
		return Crew{}, err
	}
	if err := audit(ctx, tx, c.Auth.Now().Unix(), id, actor.ID, "crew.create", id); err != nil {
		return Crew{}, err
	}
	return Crew{ID: id, Name: name}, tx.Commit()
}

func (c *Crews) AddMember(ctx context.Context, actor User, crewID, email, role string) error {
	mine, err := c.Role(ctx, crewID, actor.ID)
	if err != nil {
		return err
	}
	if rank[mine] < rank["admin"] {
		return ErrForbidden
	}
	if role != "admin" && role != "member" {
		return ErrInvalidRole
	}
	if rank[role] >= rank[mine] {
		return ErrForbidden
	}
	target, err := c.Auth.Store.UserByEmail(ctx, email)
	if err != nil {
		return err
	}
	tx, err := c.db().BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := insertOnce(ctx, tx, `INSERT INTO memberships(crew_id, user_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		crewID, target.ID, role); err != nil {
		return err
	}
	if err := audit(ctx, tx, c.Auth.Now().Unix(), crewID, actor.ID, "member.add", target.ID); err != nil {
		return err
	}
	return tx.Commit()
}

func (c *Crews) Audit(ctx context.Context, actor User, crewID string) ([]AuditEntry, error) {
	mine, err := c.Role(ctx, crewID, actor.ID)
	if err != nil {
		return nil, err
	}
	if rank[mine] < rank["admin"] {
		return nil, ErrForbidden
	}
	rows, err := c.db().QueryContext(ctx, `SELECT actor_id, action, target FROM audit WHERE crew_id = $1 ORDER BY seq`, crewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AuditEntry{}
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.Actor, &e.Action, &e.Target); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func crewFailure(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not_found"})
	case errors.Is(err, ErrForbidden):
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
	case errors.Is(err, ErrInvalidRole):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_role"})
	default:
		authFailure(w, err)
	}
}

func (c *Crews) signedIn(next func(w http.ResponseWriter, r *http.Request, u User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _, err := c.Auth.current(r)
		if err != nil {
			authFailure(w, err)
			return
		}
		next(w, r, u)
	}
}

func NewCommandAPI(c *Crews) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", NewAPI(c.Auth))
	mux.HandleFunc("POST /v1/crews", c.signedIn(func(w http.ResponseWriter, r *http.Request, u User) {
		var body struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		crew, err := c.CreateCrew(r.Context(), u, body.Name)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"id": crew.ID, "name": crew.Name})
	}))
	mux.HandleFunc("POST /v1/crews/{crew}/members", c.signedIn(func(w http.ResponseWriter, r *http.Request, u User) {
		var body struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
			return
		}
		if err := c.AddMember(r.Context(), u, r.PathValue("crew"), body.Email, body.Role); err != nil {
			crewFailure(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"role": body.Role})
	}))
	mux.HandleFunc("GET /v1/crews/{crew}/audit", c.signedIn(func(w http.ResponseWriter, r *http.Request, u User) {
		entries, err := c.Audit(r.Context(), u, r.PathValue("crew"))
		if err != nil {
			crewFailure(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"entries": entries})
	}))
	mux.HandleFunc("GET /v1/crews/{crew}/tasks/{task}", c.signedIn(func(w http.ResponseWriter, r *http.Request, u User) {
		if _, err := c.Role(r.Context(), r.PathValue("crew"), u.ID); err != nil {
			crewFailure(w, err)
			return
		}
		t, err := c.Auth.Store.GetTask(r.Context(), r.PathValue("crew"), r.PathValue("task"))
		if err != nil {
			crewFailure(w, err)
			return
		}
		writeJSON(w, http.StatusOK, t)
	}))
	return mux
}
