package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Google sign-in with the OAuth authorization-code flow. The profile comes
// from Google's userinfo endpoint over TLS with the access token, so no ID
// token signature checking is needed here.

type googleAuth struct {
	clientID, clientSecret string
	redirect               string // <public url>/auth/google/callback
	authURL, tokenURL      string
	userinfoURL            string
	client                 *http.Client
}

func newGoogleAuth(clientID, clientSecret, publicURL string) *googleAuth {
	if clientID == "" || clientSecret == "" {
		return nil
	}
	return &googleAuth{
		clientID: clientID, clientSecret: clientSecret,
		redirect:    strings.TrimRight(publicURL, "/") + "/auth/google/callback",
		authURL:     "https://accounts.google.com/o/oauth2/v2/auth",
		tokenURL:    "https://oauth2.googleapis.com/token",
		userinfoURL: "https://openidconnect.googleapis.com/v1/userinfo",
		client:      &http.Client{Timeout: 15 * time.Second},
	}
}

type googleUser struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Verified bool   `json:"email_verified"`
	Name     string `json:"name"`
}

func (g *googleAuth) exchange(ctx context.Context, code string) (googleUser, error) {
	var user googleUser
	form := url.Values{"code": {code}, "client_id": {g.clientID}, "client_secret": {g.clientSecret},
		"redirect_uri": {g.redirect}, "grant_type": {"authorization_code"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, g.tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return user, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := g.client.Do(req)
	if err != nil {
		return user, err
	}
	defer res.Body.Close()
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&tok); err != nil || res.StatusCode != 200 || tok.AccessToken == "" {
		return user, fmt.Errorf("google token: status %d", res.StatusCode)
	}
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, g.userinfoURL, nil)
	if err != nil {
		return user, err
	}
	req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	info, err := g.client.Do(req)
	if err != nil {
		return user, err
	}
	defer info.Body.Close()
	if err := json.NewDecoder(io.LimitReader(info.Body, 1<<20)).Decode(&user); err != nil || info.StatusCode != 200 || user.Sub == "" {
		return user, fmt.Errorf("google userinfo: status %d", info.StatusCode)
	}
	return user, nil
}

const oauthCookie = "orbital_oauth"

// googleStart sends the browser to Google; the state cookie carries the
// anti-forgery state and the invite code typed on the sign-in page.
func (s *server) googleStart(w http.ResponseWriter, r *http.Request) {
	a := s.accounts
	if a == nil || a.google == nil {
		http.NotFound(w, r)
		return
	}
	g := a.google
	// the state cookie must be set on the host Google returns to (localhost
	// and 127.0.0.1 keep separate cookies), so start over there
	if pub, err := url.Parse(g.redirect); err == nil && pub.Host != r.Host {
		http.Redirect(w, r, pub.Scheme+"://"+pub.Host+r.URL.RequestURI(), http.StatusFound)
		return
	}
	state := randomHex(16)
	invite := clipRunes(strings.TrimSpace(r.URL.Query().Get("invite")), 64)
	lang := "ar"
	if r.URL.Query().Get("lang") == "en" {
		lang = "en"
	}
	http.SetCookie(w, &http.Cookie{Name: oauthCookie, Value: url.Values{"s": {state}, "i": {invite}, "l": {lang}}.Encode(),
		Path: "/auth/google", MaxAge: 600, HttpOnly: true, Secure: a.secure, SameSite: http.SameSiteLaxMode})
	q := url.Values{"client_id": {g.clientID}, "redirect_uri": {g.redirect}, "response_type": {"code"},
		"scope": {"openid email profile"}, "state": {state}, "prompt": {"select_account"}}
	http.Redirect(w, r, g.authURL+"?"+q.Encode(), http.StatusFound)
}

func (s *server) googleCallback(w http.ResponseWriter, r *http.Request) {
	a := s.accounts
	if a == nil || a.google == nil {
		http.NotFound(w, r)
		return
	}
	fail := func(code string) { http.Redirect(w, r, "/?auth="+code, http.StatusFound) }
	cookie, err := r.Cookie(oauthCookie)
	http.SetCookie(w, &http.Cookie{Name: oauthCookie, Value: "", Path: "/auth/google", MaxAge: -1, HttpOnly: true, Secure: a.secure})
	if err != nil {
		fail("expired")
		return
	}
	saved, _ := url.ParseQuery(cookie.Value)
	state := r.URL.Query().Get("state")
	if state == "" || saved.Get("s") != state {
		fail("expired")
		return
	}
	if r.URL.Query().Get("error") != "" {
		fail("cancelled")
		return
	}
	user, err := a.google.exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		fail("failed")
		return
	}
	email, ok := normalizeEmail(user.Email)
	if !ok || !user.Verified {
		fail("unverified")
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	var acc *account
	for _, existing := range a.accounts {
		if existing.Google == user.Sub {
			acc = existing
		}
	}
	if acc == nil {
		// Password sign-up does not verify email ownership. Use the existing
		// password to reach that account; identify Google users by subject.
		if a.byEmail(email) != nil {
			fail("email_taken")
			return
		}
		if !a.inviteOK(saved.Get("i")) {
			fail("invite")
			return
		}
		acc = &account{Email: email, Name: clipRunes(user.Name, 40), Google: user.Sub}
		if err := a.create(acc, saved.Get("l")); err != nil {
			fail("failed")
			return
		}
	}
	if err := a.login(w, acc); err != nil {
		fail("failed")
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
