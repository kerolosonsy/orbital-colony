package main

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"orbital/colony/labkit"
)

const cookieName = "quarters"
const sessionLife = 30 * time.Minute

type session struct {
	crew, csrf string
	expires    time.Time
}

type Quarters struct {
	check  func(crew, code string) bool
	clock  labkit.Clock
	random io.Reader
	mu     sync.Mutex
	open   map[string]session
}

func NewQuarters(check func(crew, code string) bool, clock labkit.Clock, random io.Reader) (*Quarters, error) {
	if check == nil || clock == nil || random == nil {
		return nil, errors.New("invalid quarters")
	}
	return &Quarters{check: check, clock: clock, random: random, open: map[string]session{}}, nil
}

func plain(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func setCookie(w http.ResponseWriter, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: value, Path: "/", MaxAge: maxAge,
		HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})
}

// current returns the live session ID for the request; expired sessions are removed.
func (q *Quarters) current(r *http.Request) (string, session, bool) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return "", session{}, false
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	s, ok := q.open[c.Value]
	if !ok {
		return "", session{}, false
	}
	if !q.clock.Now().Before(s.expires) {
		delete(q.open, c.Value)
		return "", session{}, false
	}
	return c.Value, s, true
}

func (q *Quarters) token() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(q.random, b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func (q *Quarters) enter(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 4097))
	var in struct {
		Crew string `json:"crew"`
		Code string `json:"code"`
	}
	if err != nil || len(raw) > 4096 || json.Unmarshal(raw, &in) != nil || in.Crew == "" || in.Code == "" {
		plain(w, http.StatusBadRequest, "bad entry\n")
		return
	}
	if !q.check(in.Crew, in.Code) {
		plain(w, http.StatusUnauthorized, "denied\n")
		return
	}
	id, err := q.token()
	csrf, err2 := q.token()
	if err != nil || err2 != nil {
		plain(w, http.StatusInternalServerError, "quarters failure\n")
		return
	}
	q.mu.Lock()
	if _, taken := q.open[id]; taken {
		q.mu.Unlock()
		plain(w, http.StatusInternalServerError, "quarters failure\n")
		return
	}
	if old, err := r.Cookie(cookieName); err == nil {
		delete(q.open, old.Value)
	}
	q.open[id] = session{crew: in.Crew, csrf: csrf, expires: q.clock.Now().Add(sessionLife)}
	q.mu.Unlock()
	setCookie(w, id, int(sessionLife/time.Second))
	body, _ := json.Marshal(map[string]string{"csrf": csrf})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (q *Quarters) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routes := map[string]string{"/enter": http.MethodPost, "/me": http.MethodGet, "/logout": http.MethodPost}
	method, known := routes[r.URL.Path]
	if !known {
		plain(w, http.StatusNotFound, "no room\n")
		return
	}
	if r.Method != method {
		w.Header().Set("Allow", method)
		plain(w, http.StatusMethodNotAllowed, "wrong method\n")
		return
	}
	if r.URL.Path == "/enter" {
		q.enter(w, r)
		return
	}
	id, s, ok := q.current(r)
	if !ok {
		plain(w, http.StatusUnauthorized, "no session\n")
		return
	}
	if r.URL.Path == "/me" {
		plain(w, http.StatusOK, "crew "+s.crew+"\n")
		return
	}
	if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-CSRF-Token")), []byte(s.csrf)) != 1 {
		plain(w, http.StatusForbidden, "bad csrf\n")
		return
	}
	q.mu.Lock()
	delete(q.open, id)
	q.mu.Unlock()
	setCookie(w, "", -1)
	w.WriteHeader(http.StatusNoContent)
}

func main() {}
