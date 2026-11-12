package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func taskFailure(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

func (c *Crews) member(w http.ResponseWriter, r *http.Request, u User) bool {
	if _, err := c.Role(r.Context(), r.PathValue("crew"), u.ID); err != nil {
		crewFailure(w, err)
		return false
	}
	return true
}

func (c *Crews) createTask(w http.ResponseWriter, r *http.Request, u User) {
	if !c.member(w, r, u) {
		return
	}
	key := r.Header.Get("Idempotency-Key")
	if key == "" || len(key) > 128 {
		taskFailure(w, http.StatusBadRequest, "missing_idempotency_key")
		return
	}
	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<16))
	var body struct {
		Title string `json:"title"`
	}
	if err != nil || json.Unmarshal(raw, &body) != nil || strings.TrimSpace(body.Title) == "" {
		taskFailure(w, http.StatusBadRequest, "bad_request")
		return
	}
	sum := sha256.Sum256(raw)
	fingerprint := hex.EncodeToString(sum[:])
	ctx, crew, db := r.Context(), r.PathValue("crew"), c.db()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	defer tx.Rollback()
	// The claim is the transaction's first statement, so it takes the write
	// lock before any read and a concurrent duplicate waits for our commit.
	err = insertOnce(ctx, tx, `INSERT INTO idempotency(crew_id, user_id, key, fingerprint, status, body) VALUES ($1, $2, $3, $4, 0, '') ON CONFLICT DO NOTHING`,
		crew, u.ID, key, fingerprint)
	if errors.Is(err, ErrConflict) {
		tx.Rollback()
		var stored, saved string
		var status int
		err := db.QueryRowContext(ctx, `SELECT fingerprint, status, body FROM idempotency WHERE crew_id = $1 AND user_id = $2 AND key = $3`,
			crew, u.ID, key).Scan(&stored, &status, &saved)
		if err != nil {
			taskFailure(w, http.StatusInternalServerError, "internal")
			return
		}
		if stored != fingerprint {
			taskFailure(w, http.StatusUnprocessableEntity, "idempotency_mismatch")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Idempotent-Replayed", "true")
		w.WriteHeader(status)
		io.WriteString(w, saved)
		return
	}
	if err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	id, err := c.Auth.randomText(12)
	if err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	task := Task{ID: id, CrewID: crew, Title: body.Title, Version: 1}
	if _, err := tx.ExecContext(ctx, `INSERT INTO tasks(id, crew_id, title, done, version) VALUES ($1, $2, $3, $4, $5)`,
		task.ID, task.CrewID, task.Title, task.Done, task.Version); err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	out, _ := json.Marshal(task)
	if _, err := tx.ExecContext(ctx, `UPDATE idempotency SET status = $1, body = $2 WHERE crew_id = $3 AND user_id = $4 AND key = $5`,
		http.StatusCreated, string(out), crew, u.ID, key); err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	if err := tx.Commit(); err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (c *Crews) updateTask(w http.ResponseWriter, r *http.Request, u User) {
	if !c.member(w, r, u) {
		return
	}
	var body struct {
		Title   *string `json:"title"`
		Done    *bool   `json:"done"`
		Version int64   `json:"version"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&body); err != nil || body.Version < 1 ||
		(body.Title != nil && strings.TrimSpace(*body.Title) == "") {
		taskFailure(w, http.StatusBadRequest, "bad_request")
		return
	}
	ctx, crew, id := r.Context(), r.PathValue("crew"), r.PathValue("task")
	current, err := c.Auth.Store.GetTask(ctx, crew, id)
	if err != nil {
		crewFailure(w, err)
		return
	}
	next := current
	if body.Title != nil {
		next.Title = *body.Title
	}
	if body.Done != nil {
		next.Done = *body.Done
	}
	next.Version = body.Version + 1
	res, err := c.db().ExecContext(ctx, `UPDATE tasks SET title = $1, done = $2, version = version + 1 WHERE crew_id = $3 AND id = $4 AND version = $5`,
		next.Title, next.Done, crew, id, body.Version)
	if err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		latest, err := c.Auth.Store.GetTask(ctx, crew, id)
		if err != nil {
			crewFailure(w, err)
			return
		}
		writeJSON(w, http.StatusConflict, map[string]any{"error": "version_conflict", "version": latest.Version})
		return
	}
	writeJSON(w, http.StatusOK, next)
}

func (c *Crews) listTasks(w http.ResponseWriter, r *http.Request, u User) {
	if !c.member(w, r, u) {
		return
	}
	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 100 {
			taskFailure(w, http.StatusBadRequest, "bad_request")
			return
		}
		limit = n
	}
	after := ""
	if v := r.URL.Query().Get("cursor"); v != "" {
		raw, err := base64.RawURLEncoding.DecodeString(v)
		if err != nil || len(raw) == 0 {
			taskFailure(w, http.StatusBadRequest, "bad_request")
			return
		}
		after = string(raw)
	}
	rows, err := c.db().QueryContext(r.Context(), `SELECT id, crew_id, title, done, version FROM tasks WHERE crew_id = $1 AND id > $2 ORDER BY id LIMIT $3`,
		r.PathValue("crew"), after, limit+1)
	if err != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.CrewID, &t.Title, &t.Done, &t.Version); err != nil {
			taskFailure(w, http.StatusInternalServerError, "internal")
			return
		}
		tasks = append(tasks, t)
	}
	if rows.Err() != nil {
		taskFailure(w, http.StatusInternalServerError, "internal")
		return
	}
	next := ""
	if len(tasks) > limit {
		tasks = tasks[:limit]
		next = base64.RawURLEncoding.EncodeToString([]byte(tasks[limit-1].ID))
	}
	writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks, "next": next})
}

func NewTaskAPI(c *Crews) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", NewCommandAPI(c))
	mux.HandleFunc("POST /v1/crews/{crew}/tasks", c.signedIn(c.createTask))
	mux.HandleFunc("PATCH /v1/crews/{crew}/tasks/{task}", c.signedIn(c.updateTask))
	mux.HandleFunc("GET /v1/crews/{crew}/tasks", c.signedIn(c.listTasks))
	return mux
}
