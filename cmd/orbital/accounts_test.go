package main

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"orbital"
)

// hostedServer is a profile server in hosted mode with one invite code.
func hostedServer(t *testing.T) (*server, http.Handler) {
	t.Helper()
	passwordIters = 1000
	s, _ := profileServer(t)
	s.web, _ = fs.Sub(orbital.Web, "web")
	saves := t.TempDir()
	if err := s.setupHosted("https://orbital.example", saves, hostedConfig{invites: []string{"kepler"}, aiDaily: 2}); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	s.registerAccounts(mux)
	s.registerProfiles(mux)
	mux.HandleFunc("GET /{$}", s.index)
	mux.HandleFunc("GET /api/workspace/files", s.workspaceFiles)
	return s, s.guard(mux)
}

type browser struct {
	t      *testing.T
	h      http.Handler
	s      *server
	cookie string
}

func (b *browser) do(method, path, body string) *httptest.ResponseRecorder {
	b.t.Helper()
	req := httptest.NewRequest(method, ""+b.base()+path, strings.NewReader(body))
	req.Header.Set("X-Orbital-Token", b.s.token)
	if b.cookie != "" {
		req.Header.Set("Cookie", sessionCookie+"="+b.cookie)
	}
	rec := httptest.NewRecorder()
	b.h.ServeHTTP(rec, req)
	for _, c := range rec.Result().Cookies() {
		if c.Name == sessionCookie {
			b.cookie = c.Value
		}
	}
	return rec
}

func signUp(t *testing.T, b *browser, email, name string) {
	t.Helper()
	rec := b.do("POST", "/api/auth/signup", `{"email":"`+email+`","password":"correct horse","name":"`+name+`","invite":"KEPLER","lang":"en"}`)
	if rec.Code != 200 {
		t.Fatalf("signup: %d %s", rec.Code, rec.Body)
	}
}

func profilesOf(t *testing.T, b *browser) []profile {
	t.Helper()
	rec := b.do("GET", "/api/profiles", "")
	if rec.Code != 200 {
		t.Fatalf("profiles: %d %s", rec.Code, rec.Body)
	}
	var list []profile
	json.Unmarshal(rec.Body.Bytes(), &list)
	return list
}

func TestAccountsSignUpLoginLogout(t *testing.T) {
	s, h := hostedServer(t)
	b := &browser{t: t, h: h, s: s}
	if rec := b.do("GET", "/api/profiles", ""); rec.Code != 401 {
		t.Fatalf("signed-out API answered %d", rec.Code)
	}
	if rec := b.do("GET", "/", ""); !strings.Contains(rec.Body.String(), "login.js") {
		t.Fatal("signed-out page is not the sign-in page")
	}
	if rec := b.do("POST", "/api/auth/signup", `{"email":"a@b.co","password":"correct horse","invite":"wrong"}`); rec.Code != 403 {
		t.Fatalf("bad invite answered %d", rec.Code)
	}
	if rec := b.do("POST", "/api/auth/signup", `{"email":"a@b.co","password":"short","invite":"kepler"}`); rec.Code != 400 {
		t.Fatalf("short password answered %d", rec.Code)
	}
	signUp(t, b, "Mona@Example.com", "Mona")
	list := profilesOf(t, b)
	if len(list) != 1 || list[0].Name != "Mona" || list[0].Lang != "en" {
		t.Fatalf("first profile = %+v", list)
	}
	if rec := b.do("GET", "/", ""); !strings.Contains(rec.Body.String(), `name="orbital-account"`) {
		t.Fatal("signed-in page has no account")
	}
	other := &browser{t: t, h: h, s: s}
	if rec := other.do("POST", "/api/auth/signup", `{"email":"mona@example.com","password":"correct horse","invite":"kepler"}`); rec.Code != 409 {
		t.Fatalf("duplicate email answered %d", rec.Code)
	}
	b.do("POST", "/api/auth/logout", "")
	if rec := b.do("GET", "/api/profiles", ""); rec.Code != 401 {
		t.Fatalf("after logout the API answered %d", rec.Code)
	}
	if rec := b.do("POST", "/api/auth/login", `{"email":"mona@example.com","password":"wrong pass"}`); rec.Code != 401 {
		t.Fatalf("wrong password answered %d", rec.Code)
	}
	if rec := b.do("POST", "/api/auth/login", `{"email":"MONA@example.com","password":"correct horse"}`); rec.Code != 200 {
		t.Fatalf("login: %d %s", rec.Code, rec.Body)
	}
	if got := profilesOf(t, b); len(got) != 1 || got[0].Name != "Mona" {
		t.Fatalf("profiles after login = %+v", got)
	}
}

func TestAccountsKeepProfilesApart(t *testing.T) {
	s, h := hostedServer(t)
	mona, omar := &browser{t: t, h: h, s: s}, &browser{t: t, h: h, s: s}
	signUp(t, mona, "mona@example.com", "Mona")
	signUp(t, omar, "omar@example.com", "Omar")
	mona.do("PUT", "/api/profiles/default/progress", `{"progress":{"1-1":{"stars":3}}}`)
	if rec := omar.do("GET", "/api/profiles/default/progress", ""); strings.Contains(rec.Body.String(), "1-1") {
		t.Fatalf("Omar sees Mona's progress: %s", rec.Body)
	}
	if got := profilesOf(t, omar); len(got) != 1 || got[0].Name != "Omar" {
		t.Fatalf("Omar's profiles = %+v", got)
	}
	if rec := mona.do("GET", "/api/workspace/files?level=1-1", ""); rec.Code != 404 {
		t.Fatalf("hosted VS Code workspace answered %d", rec.Code)
	}
}

func TestProfileDuplicateCopiesProgress(t *testing.T) {
	s, h := hostedServer(t)
	b := &browser{t: t, h: h, s: s}
	signUp(t, b, "mona@example.com", "Mona")
	b.do("PUT", "/api/profiles/default/progress", `{"progress":{"1-1":{"stars":3}}}`)
	rec := b.do("POST", "/api/profiles/default/duplicate", `{"name":"Mona practice"}`)
	copied := decodeProfile(t, rec)
	if copied.ID == "default" || copied.Name != "Mona practice" || copied.Lang != "en" {
		t.Fatalf("copy = %+v", copied)
	}
	if got := b.do("GET", "/api/profiles/"+copied.ID+"/progress", ""); !strings.Contains(got.Body.String(), `"1-1"`) {
		t.Fatalf("copy lost the progress: %s", got.Body)
	}
	b.do("PUT", "/api/profiles/"+copied.ID+"/progress", `{"progress":{"1-2":{"stars":1}}}`)
	if got := b.do("GET", "/api/profiles/default/progress", ""); strings.Contains(got.Body.String(), `"1-2"`) {
		t.Fatal("playing the copy changed the original")
	}
	if got := profilesOf(t, b); len(got) != 2 {
		t.Fatalf("profiles = %+v", got)
	}
	unnamed := decodeProfile(t, b.do("POST", "/api/profiles/default/duplicate", `{}`))
	if unnamed.Name != "Mona 2" {
		t.Fatalf("default copy name = %q", unnamed.Name)
	}
}

func TestAccountsThrottleFailedLogins(t *testing.T) {
	s, h := hostedServer(t)
	b := &browser{t: t, h: h, s: s}
	signUp(t, b, "mona@example.com", "Mona")
	for i := 0; i < authFailLimit; i++ {
		b.do("POST", "/api/auth/login", `{"email":"mona@example.com","password":"guess guess"}`)
	}
	if rec := b.do("POST", "/api/auth/login", `{"email":"mona@example.com","password":"correct horse"}`); rec.Code != 429 {
		t.Fatalf("after %d failures login answered %d", authFailLimit, rec.Code)
	}
}

func TestAccountsSurviveRestartAndHideAuthFiles(t *testing.T) {
	passwordIters = 1000
	saves := t.TempDir()
	a, err := newAccountStore(saves, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	a.mu.Lock()
	acc := &account{Email: "a@b.co", Hash: "x"}
	if err := a.create(acc, "ar"); err != nil {
		t.Fatal(err)
	}
	if err := a.login(rec, acc); err != nil {
		t.Fatal(err)
	}
	a.mu.Unlock()
	again, err := newAccountStore(saves, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(rec.Result().Cookies()[0])
	if got := again.current(req); got == nil || got.ID != acc.ID {
		t.Fatalf("session lost across restart: %+v", got)
	}
	info, err := os.Stat(again.dir)
	if err != nil || info.Mode().Perm() != 0o700 {
		t.Fatalf("auth folder mode = %v, %v", info.Mode().Perm(), err)
	}
	if c := rec.Result().Cookies()[0]; !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode {
		t.Fatalf("session cookie flags: %+v", c)
	}
}

func TestPasswordHashing(t *testing.T) {
	passwordIters = 1000
	hash, err := hashPassword("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	if !checkPassword(hash, "correct horse") || checkPassword(hash, "correct hors") || checkPassword("junk", "x") {
		t.Fatal("password check is wrong")
	}
	if other, _ := hashPassword("correct horse"); other == hash {
		t.Fatal("hashes are not salted")
	}
}

func TestAIQuotaPerAccountPerDay(t *testing.T) {
	q := &aiQuota{daily: 2, used: map[string]int{}}
	if !q.take("a") || !q.take("a") || q.take("a") {
		t.Fatal("quota did not stop the third use")
	}
	if !q.take("b") {
		t.Fatal("one account's quota spent another's")
	}
}

func TestHostedConfigHidesSecretsFromPlayerPrograms(t *testing.T) {
	saved := askModel
	t.Cleanup(func() { askModel = saved })
	t.Setenv("ANTHROPIC_API_KEY", "sk-test")
	t.Setenv("GOOGLE_CLIENT_SECRET", "g-secret")
	t.Setenv("ORBITAL_INVITES", "a, b,,")
	cfg := readHostedConfig()
	if cfg.aiKey != "sk-test" || cfg.googleSecret != "g-secret" || len(cfg.invites) != 2 {
		t.Fatalf("config = %+v", cfg)
	}
	for _, name := range []string{"ANTHROPIC_API_KEY", "GOOGLE_CLIENT_SECRET", "ORBITAL_INVITES"} {
		if _, set := os.LookupEnv(name); set {
			t.Fatalf("%s still in the environment player programs inherit", name)
		}
	}
}

func TestFirstAccountTakesOverThisComputersPlayers(t *testing.T) {
	passwordIters = 1000
	saves := t.TempDir()
	local, err := newProfileStore(saves)
	if err != nil {
		t.Fatal(err)
	}
	local.mu.Lock()
	local.save(profile{ID: "kero", Name: "Kero", Lang: "ar"})
	local.mu.Unlock()
	os.WriteFile(local.dir("kero")+"/progress.json", []byte(`{"progress":{"1-1":{"stars":3}}}`), 0o600)

	s, _ := profileServer(t)
	s.profiles = nil
	if err := s.setupHosted("http://localhost:7777", saves, hostedConfig{}); err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	s.registerAccounts(mux)
	s.registerProfiles(mux)
	mux.HandleFunc("GET /api/workspace/files", s.workspaceFiles)
	h := s.guard(mux)
	first := &browser{t: t, h: h, s: s}
	signUp(t, first, "kero@example.com", "Kero")
	names := []string{}
	for _, p := range profilesOf(t, first) {
		names = append(names, p.Name)
	}
	if !strings.Contains(","+strings.Join(names, ",")+",", ",Kero,") {
		t.Fatalf("first account's players = %v", names)
	}
	if got := first.do("GET", "/api/profiles/kero/progress", ""); !strings.Contains(got.Body.String(), `"1-1"`) {
		t.Fatalf("the old player's progress was lost: %s", got.Body)
	}
	second := &browser{t: t, h: h, s: s}
	signUp(t, second, "omar@example.com", "Omar")
	if got := profilesOf(t, second); len(got) != 1 || got[0].Name != "Omar" {
		t.Fatalf("second account got = %+v", got)
	}
	if !s.vscodeAvailable() {
		t.Fatal("VS Code should stay available when signed in on localhost")
	}
}

func (b *browser) base() string {
	if b.s.public == "orbital.example" {
		return "https://orbital.example"
	}
	return "http://localhost"
}

func TestGoogleMakesAccountThenSignsIn(t *testing.T) {
	s, h := hostedServer(t)
	sub := "g-123"
	fake := controlTransport(func(r *http.Request) (*http.Response, error) {
		w := httptest.NewRecorder()
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/token" {
			w.Write([]byte(`{"access_token":"tok"}`))
		} else {
			json.NewEncoder(w).Encode(googleUser{Sub: sub, Email: "Sara@Gmail.com", Verified: true, Name: "Sara"})
		}
		return w.Result(), nil
	})
	g := newGoogleAuth("id", "secret", "https://orbital.example")
	g.tokenURL = "https://google.example/token"
	g.client = &http.Client{Transport: fake}
	s.accounts.google = g

	signIn := func(invite string) (*httptest.ResponseRecorder, string) {
		t.Helper()
		start := httptest.NewRecorder()
		h.ServeHTTP(start, httptest.NewRequest("GET", "https://orbital.example/auth/google?lang=en&invite="+invite, nil))
		to, _ := start.Result().Location()
		if start.Code != 302 || to == nil || to.Query().Get("state") == "" {
			t.Fatalf("start: %d %v", start.Code, to)
		}
		req := httptest.NewRequest("GET", "https://orbital.example/auth/google/callback?code=c&state="+to.Query().Get("state"), nil)
		for _, c := range start.Result().Cookies() {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		for _, c := range rec.Result().Cookies() {
			if c.Name == sessionCookie && c.Value != "" {
				return rec, c.Value
			}
		}
		return rec, ""
	}

	if rec, session := signIn("wrong"); session != "" || !strings.Contains(rec.Header().Get("Location"), "auth=invite") {
		t.Fatalf("new Google account without an invite: %s", rec.Header().Get("Location"))
	}
	rec, session := signIn("kepler")
	if session == "" || rec.Header().Get("Location") != "/" {
		t.Fatalf("Google sign-up: %d %s", rec.Code, rec.Header().Get("Location"))
	}
	b := &browser{t: t, h: h, s: s, cookie: session}
	if got := profilesOf(t, b); len(got) != 1 || got[0].Name != "Sara" {
		t.Fatalf("first profile = %+v", got)
	}
	// later: the same Google account signs straight in, no invite needed
	if _, again := signIn(""); again == "" {
		t.Fatal("returning Google account was not signed in")
	}
	if n := len(s.accounts.accounts); n != 1 {
		t.Fatalf("%d accounts, want 1", n)
	}
	sub = "g-other"
	if rec, session := signIn("kepler"); session != "" || rec.Header().Get("Location") != "/?auth=email_taken" {
		t.Fatalf("different Google subject with the same email: %d %s", rec.Code, rec.Header().Get("Location"))
	}

	// started on another name for this host: go to the public one first
	other := httptest.NewRecorder()
	s.googleStart(other, httptest.NewRequest("GET", "http://localhost/auth/google?lang=en", nil))
	if loc := other.Header().Get("Location"); loc != "https://orbital.example/auth/google?lang=en" {
		t.Fatalf("other host went to %q", loc)
	}
}
