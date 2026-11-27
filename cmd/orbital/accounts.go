package main

import (
	"context"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Accounts: on a hosted server each person signs in (email and password, or
// Google) and owns a set of player profiles. Account data lives in
// <saves>/auth, readable only by the server; each account's profiles live in
// <saves>/users/<account id>, the same layout the local game uses.

type account struct {
	ID      string    `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Hash    string    `json:"hash,omitempty"`   // pbkdf2-sha256 salt:key, hex
	Google  string    `json:"google,omitempty"` // Google account subject
	Created time.Time `json:"created"`
	Legacy  bool      `json:"legacy,omitempty"` // took over the players this computer had before accounts
}

type session struct {
	Account string    `json:"account"`
	Expires time.Time `json:"expires"`
}

const (
	sessionCookie   = "orbital_session"
	sessionLifetime = 30 * 24 * time.Hour
	minPassword     = 8
	authFailLimit   = 10
	authFailWindow  = 15 * time.Minute
)

// passwordIters is the PBKDF2 work per password (OWASP's figure for
// SHA-256); tests lower it.
var passwordIters = 600_000

type accountStore struct {
	mu       sync.Mutex
	dir      string // <saves>/auth
	users    string // <saves>/users
	legacy   string // <saves>/profiles: players from before accounts
	accounts map[string]*account
	sessions map[string]session // sha256(token) → session
	stores   map[string]*profileStore
	invites  []string
	secure   bool // cookies only over https
	failures map[string][]time.Time
	google   *googleAuth
}

func newAccountStore(saves string, invites []string, secure bool) (*accountStore, error) {
	a := &accountStore{
		dir: filepath.Join(saves, "auth"), users: filepath.Join(saves, "users"), legacy: filepath.Join(saves, "profiles"),
		accounts: map[string]*account{}, sessions: map[string]session{}, stores: map[string]*profileStore{},
		invites: invites, secure: secure, failures: map[string][]time.Time{},
	}
	if err := os.MkdirAll(a.dir, 0o700); err != nil {
		return nil, err
	}
	if err := os.Chmod(a.dir, 0o700); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(a.users, 0o755); err != nil {
		return nil, err
	}
	if err := readJSONFile(filepath.Join(a.dir, "accounts.json"), &a.accounts); err != nil {
		return nil, err
	}
	if err := readJSONFile(filepath.Join(a.dir, "sessions.json"), &a.sessions); err != nil {
		return nil, err
	}
	return a, nil
}

func readJSONFile(path string, v any) error {
	blob, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(blob, v)
}

// Caller holds mu.
func (a *accountStore) saveAccounts() error {
	blob, err := json.Marshal(a.accounts)
	if err != nil {
		return err
	}
	return atomicJSON(filepath.Join(a.dir, "accounts.json"), blob)
}

// Caller holds mu.
func (a *accountStore) saveSessions() error {
	now := time.Now()
	for key, s := range a.sessions {
		if now.After(s.Expires) {
			delete(a.sessions, key)
		}
	}
	blob, err := json.Marshal(a.sessions)
	if err != nil {
		return err
	}
	return atomicJSON(filepath.Join(a.dir, "sessions.json"), blob)
}

// secretKey returns a random value kept across restarts, so pages stay valid
// when the server is redeployed.
func secretKey(dir string) (string, error) {
	path := filepath.Join(dir, "secret")
	if blob, err := os.ReadFile(path); err == nil && len(blob) >= 32 {
		return strings.TrimSpace(string(blob)), nil
	}
	key := randomHex(32)
	return key, atomicJSON(path, []byte(key))
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func tokenKey(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	key, err := pbkdf2.Key(sha256.New, password, salt, passwordIters, 32)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(key), nil
}

func checkPassword(hash, password string) bool {
	saltHex, keyHex, ok := strings.Cut(hash, ":")
	salt, err1 := hex.DecodeString(saltHex)
	want, err2 := hex.DecodeString(keyHex)
	if !ok || err1 != nil || err2 != nil {
		return false
	}
	got, err := pbkdf2.Key(sha256.New, password, salt, passwordIters, len(want))
	return err == nil && subtle.ConstantTimeCompare(got, want) == 1
}

func normalizeEmail(email string) (string, bool) {
	email = strings.ToLower(strings.TrimSpace(email))
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email || len(email) > 254 {
		return "", false
	}
	return email, true
}

// Caller holds mu.
func (a *accountStore) byEmail(email string) *account {
	for _, acc := range a.accounts {
		if acc.Email == email {
			return acc
		}
	}
	return nil
}

func (a *accountStore) inviteOK(code string) bool {
	if len(a.invites) == 0 {
		return true // open sign-up
	}
	code = strings.TrimSpace(code)
	ok := false
	for _, invite := range a.invites {
		if subtle.ConstantTimeCompare([]byte(strings.ToLower(code)), []byte(strings.ToLower(invite))) == 1 {
			ok = true
		}
	}
	return ok
}

// create adds an account with its first player profile, named after the
// person and speaking the language they signed up in. The very first
// account on a computer that already had players takes those players over,
// so turning accounts on never loses anyone's progress. Caller holds mu.
func (a *accountStore) create(acc *account, lang string) error {
	acc.ID = randomHex(8)
	acc.Created = time.Now().UTC()
	adopt := len(a.accounts) == 0 && hasProfiles(a.legacy)
	if adopt {
		if err := os.MkdirAll(filepath.Join(a.users, acc.ID), 0o755); err != nil {
			return err
		}
		if err := os.Rename(a.legacy, filepath.Join(a.users, acc.ID, "profiles")); err != nil {
			return err
		}
		acc.Legacy = true
	}
	a.accounts[acc.ID] = acc
	if err := a.saveAccounts(); err != nil {
		delete(a.accounts, acc.ID)
		return err
	}
	if adopt {
		return nil
	}
	p := &profileStore{root: filepath.Join(a.users, acc.ID), games: map[string]*colonyGame{}}
	name := strings.TrimSpace(acc.Name)
	if name == "" {
		name = strings.Split(acc.Email, "@")[0]
	}
	if lang != "en" {
		lang = "ar"
	}
	return p.save(profile{ID: "default", Name: clipRunes(name, 24), Created: acc.Created, Lang: lang})
}

// hasProfiles reports whether dir holds at least one saved player.
func hasProfiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if _, err := os.Stat(filepath.Join(dir, entry.Name(), "profile.json")); entry.IsDir() && err == nil {
			return true
		}
	}
	return false
}

func clipRunes(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}

// login starts a session for acc and sets its cookie. Caller holds mu.
func (a *accountStore) login(w http.ResponseWriter, acc *account) error {
	token := randomHex(32)
	a.sessions[tokenKey(token)] = session{Account: acc.ID, Expires: time.Now().Add(sessionLifetime)}
	if err := a.saveSessions(); err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: token, Path: "/", MaxAge: int(sessionLifetime.Seconds()),
		HttpOnly: true, Secure: a.secure, SameSite: http.SameSiteLaxMode})
	return nil
}

// current returns the signed-in account for a request, or nil.
func (a *accountStore) current(r *http.Request) *account {
	cookie, err := r.Cookie(sessionCookie)
	if err != nil || cookie.Value == "" {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	s, ok := a.sessions[tokenKey(cookie.Value)]
	if !ok || time.Now().After(s.Expires) {
		return nil
	}
	return a.accounts[s.Account]
}

// profiles returns the account's own profile store.
func (a *accountStore) profiles(id string) *profileStore {
	a.mu.Lock()
	defer a.mu.Unlock()
	if p := a.stores[id]; p != nil {
		return p
	}
	p, err := newProfileStore(filepath.Join(a.users, id))
	if err != nil {
		return nil
	}
	a.stores[id] = p
	return p
}

func (a *accountStore) allStores() []*profileStore {
	a.mu.Lock()
	defer a.mu.Unlock()
	stores := make([]*profileStore, 0, len(a.stores))
	for _, p := range a.stores {
		stores = append(stores, p)
	}
	return stores
}

// throttled counts a failed sign-in for key; too many in the window block it.
func (a *accountStore) throttled(key string, failed bool) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	cutoff := time.Now().Add(-authFailWindow)
	recent := a.failures[key][:0]
	for _, t := range a.failures[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	if failed {
		recent = append(recent, time.Now())
	}
	a.failures[key] = recent
	return len(recent) >= authFailLimit
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ---- request plumbing ----

type accountKey struct{}

func withAccount(ctx context.Context, acc *account) context.Context {
	return context.WithValue(ctx, accountKey{}, acc)
}

func requestAccount(r *http.Request) *account {
	acc, _ := r.Context().Value(accountKey{}).(*account)
	return acc
}

// store is the profile store a request works on: the signed-in account's on
// a hosted server, the one local store otherwise.
func (s *server) store(r *http.Request) *profileStore {
	if s.accounts == nil {
		return s.profiles
	}
	acc := requestAccount(r)
	if acc == nil {
		return nil
	}
	return s.accounts.profiles(acc.ID)
}

// ---- handlers ----

func (s *server) registerAccounts(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/auth/me", s.authMe)
	mux.HandleFunc("POST /api/auth/signup", s.authSignup)
	mux.HandleFunc("POST /api/auth/login", s.authLogin)
	mux.HandleFunc("POST /api/auth/logout", s.authLogout)
	mux.HandleFunc("GET /auth/google", s.googleStart)
	mux.HandleFunc("GET /auth/google/callback", s.googleCallback)
}

type authInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Invite   string `json:"invite"`
	Lang     string `json:"lang"`
}

func authError(w http.ResponseWriter, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": code})
}

func readAuthInput(w http.ResponseWriter, r *http.Request) (authInput, bool) {
	var in authInput
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&in); err != nil {
		authError(w, "bad_request", http.StatusBadRequest)
		return in, false
	}
	return in, true
}

func publicAccount(acc *account) map[string]any {
	return map[string]any{"id": acc.ID, "email": acc.Email, "name": acc.Name, "google": acc.Google != "", "password": acc.Hash != ""}
}

func (s *server) authMe(w http.ResponseWriter, r *http.Request) {
	if s.accounts == nil {
		authError(w, "no_accounts", http.StatusNotFound)
		return
	}
	google := s.accounts.google != nil
	acc := s.accounts.current(r)
	if acc == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{"error": "login", "googleEnabled": google, "inviteRequired": len(s.accounts.invites) > 0})
		return
	}
	info := publicAccount(acc)
	info["googleEnabled"] = google
	writeJSON(w, info)
}

func (s *server) authSignup(w http.ResponseWriter, r *http.Request) {
	a := s.accounts
	if a == nil {
		authError(w, "no_accounts", http.StatusNotFound)
		return
	}
	in, ok := readAuthInput(w, r)
	if !ok {
		return
	}
	ip := clientIP(r)
	if a.throttled(ip, false) {
		authError(w, "too_many", http.StatusTooManyRequests)
		return
	}
	email, valid := normalizeEmail(in.Email)
	if !valid {
		authError(w, "bad_email", http.StatusBadRequest)
		return
	}
	if len(in.Password) < minPassword || len(in.Password) > 200 {
		authError(w, "weak_password", http.StatusBadRequest)
		return
	}
	if !a.inviteOK(in.Invite) {
		a.throttled(ip, true)
		authError(w, "bad_invite", http.StatusForbidden)
		return
	}
	hash, err := hashPassword(in.Password)
	if err != nil {
		authError(w, "failed", http.StatusInternalServerError)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.byEmail(email) != nil {
		authError(w, "email_taken", http.StatusConflict)
		return
	}
	acc := &account{Email: email, Name: clipRunes(strings.TrimSpace(in.Name), 40), Hash: hash}
	if err := a.create(acc, in.Lang); err != nil {
		authError(w, "failed", http.StatusInternalServerError)
		return
	}
	if err := a.login(w, acc); err != nil {
		authError(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, publicAccount(acc))
}

func (s *server) authLogin(w http.ResponseWriter, r *http.Request) {
	a := s.accounts
	if a == nil {
		authError(w, "no_accounts", http.StatusNotFound)
		return
	}
	in, ok := readAuthInput(w, r)
	if !ok {
		return
	}
	ip := clientIP(r)
	email, _ := normalizeEmail(in.Email)
	if a.throttled(ip, false) || a.throttled("email:"+email, false) {
		authError(w, "too_many", http.StatusTooManyRequests)
		return
	}
	a.mu.Lock()
	acc := a.byEmail(email)
	hash := ""
	if acc != nil {
		hash = acc.Hash
	}
	a.mu.Unlock()
	if hash == "" {
		hash = "00:00" // same work either way: don't reveal which emails exist
	}
	if !checkPassword(hash, in.Password) || acc == nil {
		a.throttled(ip, true)
		a.throttled("email:"+email, true)
		authError(w, "wrong_login", http.StatusUnauthorized)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.login(w, acc); err != nil {
		authError(w, "failed", http.StatusInternalServerError)
		return
	}
	writeJSON(w, publicAccount(acc))
}

func (s *server) authLogout(w http.ResponseWriter, r *http.Request) {
	if a := s.accounts; a != nil {
		if cookie, err := r.Cookie(sessionCookie); err == nil {
			a.mu.Lock()
			delete(a.sessions, tokenKey(cookie.Value))
			a.saveSessions()
			a.mu.Unlock()
		}
		http.SetCookie(w, &http.Cookie{Name: sessionCookie, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, Secure: a.secure, SameSite: http.SameSiteLaxMode})
	}
	writeJSON(w, map[string]bool{"ok": true})
}
