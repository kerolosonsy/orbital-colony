package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWeakPassword   = errors.New("weak password")
	ErrBadCredentials = errors.New("bad credentials")
	ErrUnauthorized   = errors.New("unauthorized")
)

const SessionCookie = "session"

type Auth struct {
	Store  *Store
	Now    func() time.Time
	Random io.Reader
	TTL    time.Duration
	Cost   int
}

func (a *Auth) randomText(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(a.Random, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func digest(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (a *Auth) Register(ctx context.Context, email, password string) (User, error) {
	if len(password) < 8 || len(password) > 72 {
		return User{}, ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), a.Cost)
	if err != nil {
		return User{}, err
	}
	id, err := a.randomText(12)
	if err != nil {
		return User{}, err
	}
	u := User{ID: id, Email: normalEmail(email), PasswordHash: string(hash)}
	if err := a.Store.CreateUser(ctx, u); err != nil {
		return User{}, err
	}
	return u, nil
}

func (a *Auth) Login(ctx context.Context, email, password string) (string, error) {
	u, err := a.Store.UserByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return "", ErrBadCredentials
	}
	if err != nil {
		return "", err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", ErrBadCredentials
	}
	token, err := a.randomText(32)
	if err != nil {
		return "", err
	}
	_, err = a.Store.db.ExecContext(ctx, `INSERT INTO sessions(digest, user_id, expires_at) VALUES ($1, $2, $3)`,
		digest(token), u.ID, a.Now().Add(a.TTL).Unix())
	if err != nil {
		return "", err
	}
	return token, nil
}

func (a *Auth) Authenticate(ctx context.Context, token string) (User, error) {
	var u User
	var expires int64
	err := a.Store.db.QueryRowContext(ctx, `SELECT u.id, u.email, u.password_hash, s.expires_at
 FROM sessions s JOIN users u ON u.id = s.user_id WHERE s.digest = $1`, digest(token)).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &expires)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUnauthorized
	}
	if err != nil {
		return User{}, err
	}
	if a.Now().Unix() >= expires {
		if _, err := a.Store.db.ExecContext(ctx, `DELETE FROM sessions WHERE digest = $1`, digest(token)); err != nil {
			return User{}, err
		}
		return User{}, ErrUnauthorized
	}
	return u, nil
}

func (a *Auth) Logout(ctx context.Context, token string) error {
	_, err := a.Store.db.ExecContext(ctx, `DELETE FROM sessions WHERE digest = $1`, digest(token))
	return err
}

func (a *Auth) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	var hash string
	err := a.Store.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBadCredentials
	}
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(oldPassword)) != nil {
		return ErrBadCredentials
	}
	if len(newPassword) < 8 || len(newPassword) > 72 {
		return ErrWeakPassword
	}
	next, err := bcrypt.GenerateFromPassword([]byte(newPassword), a.Cost)
	if err != nil {
		return err
	}
	tx, err := a.Store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `UPDATE users SET password_hash = $1 WHERE id = $2`, string(next), userID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return tx.Commit()
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Old      string `json:"old"`
	New      string `json:"new"`
}

func (a *Auth) current(r *http.Request) (User, string, error) {
	c, err := r.Cookie(SessionCookie)
	if err != nil {
		return User{}, "", ErrUnauthorized
	}
	u, err := a.Authenticate(r.Context(), c.Value)
	return u, c.Value, err
}

func authFailure(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUnauthorized), errors.Is(err, ErrBadCredentials):
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	case errors.Is(err, ErrWeakPassword):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "weak_password"})
	case errors.Is(err, ErrConflict):
		writeJSON(w, http.StatusConflict, map[string]string{"error": "conflict"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal"})
	}
}

func decode(w http.ResponseWriter, r *http.Request, c *credentials) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(c); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad_request"})
		return false
	}
	return true
}

func (a *Auth) setCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{Name: SessionCookie, Value: value, Path: "/", MaxAge: maxAge,
		HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})
}

func NewAPI(a *Auth) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", NewServer(a.Store))
	mux.HandleFunc("POST /v1/register", func(w http.ResponseWriter, r *http.Request) {
		var c credentials
		if !decode(w, r, &c) {
			return
		}
		u, err := a.Register(r.Context(), c.Email, c.Password)
		if err != nil {
			authFailure(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"id": u.ID, "email": u.Email})
	})
	mux.HandleFunc("POST /v1/login", func(w http.ResponseWriter, r *http.Request) {
		var c credentials
		if !decode(w, r, &c) {
			return
		}
		token, err := a.Login(r.Context(), c.Email, c.Password)
		if err != nil {
			authFailure(w, err)
			return
		}
		a.setCookie(w, token, int(a.TTL/time.Second))
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("GET /v1/me", func(w http.ResponseWriter, r *http.Request) {
		u, _, err := a.current(r)
		if err != nil {
			authFailure(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"id": u.ID, "email": u.Email})
	})
	mux.HandleFunc("POST /v1/logout", func(w http.ResponseWriter, r *http.Request) {
		if c, err := r.Cookie(SessionCookie); err == nil {
			if err := a.Logout(r.Context(), c.Value); err != nil {
				authFailure(w, err)
				return
			}
		}
		a.setCookie(w, "", -1)
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /v1/password", func(w http.ResponseWriter, r *http.Request) {
		u, _, err := a.current(r)
		if err != nil {
			authFailure(w, err)
			return
		}
		var c credentials
		if !decode(w, r, &c) {
			return
		}
		if err := a.ChangePassword(r.Context(), u.ID, c.Old, c.New); err != nil {
			authFailure(w, err)
			return
		}
		a.setCookie(w, "", -1)
		w.WriteHeader(http.StatusNoContent)
	})
	return mux
}
