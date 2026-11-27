package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Account mode adds local sign-in and an optional Claude API helper.

type hostedConfig struct {
	invites                []string
	googleID, googleSecret string
	aiKey                  string
	aiDaily                int
}

// readHostedConfig reads the settings from the environment and removes the
// secrets from it, so player programs (which inherit the environment) never
// see them.
func readHostedConfig() hostedConfig {
	cfg := hostedConfig{
		googleID: os.Getenv("GOOGLE_CLIENT_ID"), googleSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		aiKey: os.Getenv("ANTHROPIC_API_KEY"), aiDaily: 40,
	}
	for _, code := range strings.Split(os.Getenv("ORBITAL_INVITES"), ",") {
		if code = strings.TrimSpace(code); code != "" {
			cfg.invites = append(cfg.invites, code)
		}
	}
	if n, err := strconv.Atoi(os.Getenv("ORBITAL_AI_DAILY")); err == nil && n >= 0 {
		cfg.aiDaily = n
	}
	for _, name := range []string{"GOOGLE_CLIENT_SECRET", "ANTHROPIC_API_KEY", "ORBITAL_INVITES"} {
		os.Unsetenv(name)
	}
	if cfg.aiKey != "" {
		askModel = claudeAsker(cfg.aiKey)
	}
	return cfg
}

func (s *server) setupHosted(publicURL, saves string, cfg hostedConfig) error {
	u, err := url.Parse(publicURL)
	if err != nil || u.Host == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return fmt.Errorf("-public must be a full URL like http://127.0.0.1:7777")
	}
	accounts, err := newAccountStore(saves, cfg.invites, u.Scheme == "https")
	if err != nil {
		return err
	}
	accounts.google = newGoogleAuth(cfg.googleID, cfg.googleSecret, publicURL)
	token, err := secretKey(accounts.dir)
	if err != nil {
		return err
	}
	s.accounts, s.public, s.token = accounts, u.Hostname(), token
	s.local = loopbackHost(s.public)
	if cfg.aiKey != "" { // the host pays for Claude; the player's own Codex is theirs
		s.ai = &aiQuota{daily: cfg.aiDaily, used: map[string]int{}}
	}
	return nil
}

// aiQuota caps helper questions and code scores per account per day, since
// on a hosted server they are paid from the host's API key.
type aiQuota struct {
	mu    sync.Mutex
	daily int
	day   string
	used  map[string]int
}

func (q *aiQuota) take(account string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	if today := time.Now().UTC().Format("2006-01-02"); today != q.day {
		q.day, q.used = today, map[string]int{}
	}
	if q.used[account] >= q.daily {
		return false
	}
	q.used[account]++
	return true
}

// aiAllowed spends one of the account's daily helper uses; when none are
// left it answers {"error":"limit"} and returns false.
func (s *server) aiAllowed(w http.ResponseWriter, r *http.Request) bool {
	acc := requestAccount(r)
	if s.ai == nil || acc == nil || s.ai.take(acc.ID) {
		return true
	}
	writeJSON(w, map[string]any{"error": "limit"})
	return false
}
