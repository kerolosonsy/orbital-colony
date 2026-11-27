package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

func profileServer(t *testing.T) (*server, http.Handler) {
	t.Helper()
	p, err := newProfileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	runner := testRunner(t)
	workspaceRoot := t.TempDir()
	colonyRel, err := filepath.Rel(workspaceRoot, runner.colonyDir)
	if err != nil {
		t.Fatal(err)
	}
	s := &server{profiles: p, run: runner, token: "test-token", ws: &workspace{root: workspaceRoot, colonyRel: colonyRel}}
	mux := http.NewServeMux()
	s.registerProfiles(mux)
	mux.HandleFunc("GET /api/colony", s.colonyStatus)
	mux.HandleFunc("POST /api/colony/run", s.colonyRun)
	mux.HandleFunc("POST /api/colony/stop", s.colonyStop)
	mux.HandleFunc("POST /api/colony/claim", s.colonyClaim)
	mux.HandleFunc("GET /api/colony/stream", s.colonyStream)
	mux.HandleFunc("POST /api/run", s.runCode)
	mux.HandleFunc("GET /api/workspace/files", s.workspaceFiles)
	t.Cleanup(func() { s.colonyOps.Lock(); defer s.colonyOps.Unlock(); s.stopColonies() })
	return s, s.guard(mux)
}

func profileRequest(t *testing.T, h http.Handler, method, path, id, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, "http://localhost"+path, strings.NewReader(body))
	req.Header.Set("X-Orbital-Token", "test-token")
	req.Header.Set("X-Orbital-Profile", id)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func decodeProfile(t *testing.T, rec *httptest.ResponseRecorder) profile {
	t.Helper()
	if rec.Code != 200 {
		t.Fatalf("HTTP %d: %s", rec.Code, rec.Body)
	}
	var pr profile
	if err := json.Unmarshal(rec.Body.Bytes(), &pr); err != nil {
		t.Fatal(err)
	}
	return pr
}

func TestProfileCRUDProgressResetAndDeleteLast(t *testing.T) {
	s, h := profileServer(t)
	listed := profileRequest(t, h, "GET", "/api/profiles", "", "")
	var profiles []profile
	if err := json.Unmarshal(listed.Body.Bytes(), &profiles); err != nil {
		t.Fatal(err)
	}
	if listed.Code != 200 || len(profiles) != 1 || profiles[0].ID != "default" || profiles[0].Name != "Player 1" {
		t.Fatalf("initial profiles: %s", listed.Body)
	}
	pr := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles", "", `{"name":"Alex","lang":"en"}`))
	duplicate := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles", "", `{"name":"Alex","lang":"ar"}`))
	if pr.ID == duplicate.ID || !profileID.MatchString(pr.ID) || pr.Lang != "en" || pr.Created.IsZero() {
		t.Fatalf("invalid profiles: %+v %+v", pr, duplicate)
	}
	renamed := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles/"+pr.ID+"/rename", "", `{"name":"Alex 2"}`))
	changed := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles/"+pr.ID+"/lang", "", `{"lang":"ar"}`))
	if renamed.Name != "Alex 2" || changed.Lang != "ar" || changed.Created != pr.Created {
		t.Fatalf("profile update: %+v %+v", renamed, changed)
	}
	base := "/api/profiles/" + pr.ID
	progress := `{"stars":{"1-1":3},"code":"hello","settings":{"sound":false}}`
	put := profileRequest(t, h, "PUT", base+"/progress", "", progress)
	got := profileRequest(t, h, "GET", base+"/progress", "", "")
	if put.Code != 200 || strings.TrimSpace(got.Body.String()) != progress {
		t.Fatalf("progress: %s / %s", put.Body, got.Body)
	}
	if empty := profileRequest(t, h, "GET", "/api/profiles/"+duplicate.ID+"/progress", "", ""); strings.TrimSpace(empty.Body.String()) != "{}" {
		t.Fatalf("progress leaked: %s", empty.Body)
	}
	if rec := profileRequest(t, h, "POST", "/api/colony/claim", pr.ID, `{"chapter":1}`); rec.Code != 200 {
		t.Fatalf("claim: %s", rec.Body)
	}
	ws := s.ws.forProfile(pr.ID)
	if _, err := ws.open(levels.All[0], nil); err != nil {
		t.Fatal(err)
	}
	reset := profileRequest(t, h, "POST", base+"/reset", "", "")
	if reset.Code != 200 {
		t.Fatalf("reset: %s", reset.Body)
	}
	for _, path := range []string{filepath.Join(s.profiles.dir(pr.ID), "progress.json"), filepath.Join(s.profiles.dir(pr.ID), "colony.json"), ws.root} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatalf("reset left %s: %v", path, err)
		}
	}
	if _, err := s.profiles.load(pr.ID); err != nil {
		t.Fatal(err)
	}
	status := profileRequest(t, h, "GET", "/api/colony", pr.ID, "")
	var snapshot struct{ State econ.State }
	if err := json.Unmarshal(status.Body.Bytes(), &snapshot); err != nil {
		t.Fatal(err)
	}
	if len(snapshot.State.Claimed) != 0 {
		t.Fatalf("reset retained rewards: %s", status.Body)
	}
	for _, id := range []string{pr.ID, duplicate.ID} {
		if rec := profileRequest(t, h, "DELETE", "/api/profiles/"+id, "", ""); rec.Code != 200 {
			t.Fatalf("delete: %s", rec.Body)
		}
	}
	if rec := profileRequest(t, h, "DELETE", "/api/profiles/default", "", ""); rec.Code != 409 {
		t.Fatalf("last delete: %d %s", rec.Code, rec.Body)
	}
}

func TestProfileRejectsInvalidIDsAndProgress(t *testing.T) {
	_, h := profileServer(t)
	for _, id := range []string{"../escape", "UPPER", "-leading", strings.Repeat("a", 41)} {
		rec := profileRequest(t, h, "GET", "/api/colony", id, "")
		if rec.Code != 400 {
			t.Fatalf("accepted %q: %d", id, rec.Code)
		}
	}
	for _, tc := range []struct {
		body   string
		status int
	}{{`{"broken":`, 400}, {`"` + strings.Repeat("a", 2<<20) + `"`, 413}} {
		rec := profileRequest(t, h, "PUT", "/api/profiles/default/progress", "", tc.body)
		if rec.Code != tc.status {
			t.Fatalf("progress status %d, want %d", rec.Code, tc.status)
		}
	}
	req := httptest.NewRequest("GET", "http://localhost/api/profiles", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 403 {
		t.Fatalf("unguarded profiles: %d", rec.Code)
	}
}

func TestProfileColonyStateIsolation(t *testing.T) {
	s, h := profileServer(t)
	pr := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles", "", `{"name":"Second","lang":"en"}`))
	if _, err := s.profiles.ensureDefault(); err != nil {
		t.Fatal(err)
	}
	finishChapter1(t, s, "default")
	rec := profileRequest(t, h, "POST", "/api/colony/claim", "default", `{"chapter":1}`)
	if rec.Code != 200 {
		t.Fatalf("claim: %s", rec.Body)
	}
	for _, id := range []string{"default", pr.ID} {
		rec := profileRequest(t, h, "GET", "/api/colony", id, "")
		var status struct{ State econ.State }
		if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
			t.Fatal(err)
		}
		if status.State.Claimed["chapter-1"] != (id == "default") {
			t.Fatalf("reward leaked to %s: %s", id, rec.Body)
		}
		reloaded, err := newColonyGame(s.profiles.dir(id))
		if err != nil {
			t.Fatal(err)
		}
		if reloaded.state.Claimed["chapter-1"] != (id == "default") {
			t.Fatalf("saved reward leaked to %s", id)
		}
	}
}

func TestProfileLegacyMigration(t *testing.T) {
	root := t.TempDir()
	legacy, err := newColonyGame(root)
	if err != nil {
		t.Fatal(err)
	}
	legacy.state.Credits = 321
	if err := legacy.persist(); err != nil {
		t.Fatal(err)
	}
	p, err := newProfileStore(root)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := p.load("default")
	if err != nil || pr.Name != "Player 1" {
		t.Fatalf("default: %+v %v", pr, err)
	}
	migrated := colonySavedState(t, filepath.Join(p.dir("default"), "colony.json"))
	if migrated.Credits != 321 {
		t.Fatalf("lost legacy credits: %v", migrated.Credits)
	}
	if _, err := os.Stat(legacy.path); !os.IsNotExist(err) {
		t.Fatalf("legacy was not moved: %v", err)
	}
	if _, err := newProfileStore(root); err != nil {
		t.Fatal(err)
	}
}

func TestProfileRunSwitchStopsAndPersistsPrevious(t *testing.T) {
	s, h := profileServer(t)
	pr := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles", "", `{"name":"Second"}`))
	body, _ := json.Marshal(map[string]string{"code": levels.Colony.Starter})
	first := profileRequest(t, h, "POST", "/api/colony/run", "default", string(body))
	if first.Code != 200 || !strings.Contains(first.Body.String(), `"running":true`) {
		t.Fatalf("first start: %s", first.Body)
	}
	req := httptest.NewRequest("GET", "http://localhost/", nil)
	c, err := s.profileGame(req)
	if err != nil {
		t.Fatal(err)
	}
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	c.mu.Lock()
	active := c.active
	c.mu.Unlock()
	if active == nil {
		t.Fatal("first run exited early")
	}
	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	advanced := false
	for !advanced {
		select {
		case line := <-ch:
			var event struct{ State *econ.State }
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				t.Fatal(err)
			}
			advanced = event.State != nil && event.State.Tick > 0
		case <-timer.C:
			t.Fatal("first run did not advance")
		}
	}
	second := profileRequest(t, h, "POST", "/api/colony/run", pr.ID, string(body))
	if second.Code != 200 || !strings.Contains(second.Body.String(), `"running":true`) {
		t.Fatalf("second start: %s", second.Body)
	}
	select {
	case <-active.done:
	default:
		t.Fatal("previous profile still running")
	}
	if saved := colonySavedState(t, c.path); saved.Tick == 0 {
		t.Fatal("previous run was not persisted")
	}
	for len(ch) > 0 {
		<-ch
	}
	req.Header.Set("X-Orbital-Profile", pr.ID)
	other, err := s.profileGame(req)
	if err != nil {
		t.Fatal(err)
	}
	other.broadcast(`{"op":"other-profile"}`)
	select {
	case ev := <-ch:
		t.Fatalf("stream crossed profiles: %s", ev)
	default:
	}
	reset := profileRequest(t, h, "POST", "/api/profiles/"+pr.ID+"/reset", "", "")
	if reset.Code != 200 {
		t.Fatalf("running reset: %s", reset.Body)
	}
	other.mu.Lock()
	running := other.active != nil
	other.mu.Unlock()
	if running {
		t.Fatal("reset left process running")
	}
}

func TestLangMissionPerRequest(t *testing.T) {
	_, h := profileServer(t)
	code, err := os.ReadFile(filepath.Join("testdata", "solutions", "1-1-first-signal.go"))
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(map[string]string{"level": "1-1-first-signal", "code": string(code)})
	for _, lang := range []string{"en", "ar", "en"} {
		req := httptest.NewRequest("POST", "http://localhost/api/run", strings.NewReader(string(body)))
		req.Header.Set("X-Orbital-Token", "test-token")
		req.Header.Set("X-Orbital-Lang", lang)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		var response runResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if !response.Result.Success {
			t.Fatalf("%s mission failed: %s", lang, rec.Body)
		}
		var goal struct{ Op, Msg string }
		for _, event := range response.Events {
			var ev struct{ Op, Msg string }
			if err := json.Unmarshal(event, &ev); err != nil {
				t.Fatal(err)
			}
			if ev.Op == "goal" {
				goal = ev
			}
		}
		for _, message := range []string{goal.Msg, response.Result.Message} {
			arabic := strings.IndexFunc(message, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
			if message == "" || arabic != (lang == "ar") {
				t.Fatalf("%s message: %q", lang, message)
			}
		}
	}
}

func TestProfileStreamQuerySelectsOnlyItsProfile(t *testing.T) {
	s, h := profileServer(t)
	pr := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles", "", `{"name":"Stream"}`))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := httptest.NewRequest("GET", "http://localhost/api/colony/stream?token=test-token&profile="+pr.ID+"&lang=en", nil).WithContext(ctx)
	rec := &colonyStreamRecorder{ResponseRecorder: httptest.NewRecorder(), flushed: make(chan struct{}, 4)}
	done := make(chan struct{})
	go func() { defer close(done); h.ServeHTTP(rec, req) }()
	defer func() { cancel(); <-done }()
	select {
	case <-rec.flushed:
	case <-ctx.Done():
		t.Fatal("stream did not open")
	}
	c, err := s.profileGame(req)
	if err != nil {
		t.Fatal(err)
	}
	c.broadcast(`{"op":"profile-event"}`)
	select {
	case <-rec.flushed:
	case <-ctx.Done():
		t.Fatal("stream did not receive own event")
	}
	cancel()
	<-done
	if !strings.Contains(rec.Body.String(), `"op":"profile-event"`) {
		t.Fatalf("stream: %s", rec.Body)
	}
}

func TestProfileWorkspacesKeepIndependentCode(t *testing.T) {
	s, h := profileServer(t)
	pr := decodeProfile(t, profileRequest(t, h, "POST", "/api/profiles", "", `{"name":"Coder"}`))
	for _, id := range []string{"default", pr.ID} {
		ws := s.ws.forProfile(id)
		code := "package main\n// " + id + "\nfunc main() {}\n"
		path, err := ws.open(levels.All[0], map[string]string{"main.go": code})
		if err != nil {
			t.Fatal(err)
		}
		if path != filepath.Join(s.ws.root, id, levels.All[0].ID, "main.go") {
			t.Fatalf("workspace path: %s", path)
		}
		gomod, err := os.ReadFile(filepath.Join(ws.root, "go.mod"))
		if err != nil {
			t.Fatal(err)
		}
		_, replacement, ok := strings.Cut(string(gomod), "replace orbital/colony => ")
		if !ok || !strings.Contains(string(gomod), "module play") || filepath.Clean(filepath.Join(ws.root, strings.TrimSpace(replacement))) != s.run.colonyDir {
			t.Fatalf("workspace module: %s", gomod)
		}
	}
	for _, id := range []string{"default", pr.ID} {
		rec := profileRequest(t, h, "GET", "/api/workspace/files?level="+levels.All[0].ID, id, "")
		var response struct{ Files map[string]string }
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if rec.Code != 200 || !strings.Contains(response.Files["main.go"], "// "+id+"\n") {
			t.Fatalf("workspace leaked for %s: %s", id, rec.Body)
		}
	}
}

func TestLangLabsAndCompilerTips(t *testing.T) {
	runner := testRunner(t)
	lvl := &levels.Level{Kind: "lab", Files: []levels.LabFile{{Name: "main.go", Starter: "package main\nfunc main() {}"}, {Name: "main_test.go", Starter: `package main
import ("fmt"; "testing"; "orbital/colony")
func TestMessage(t *testing.T) { fmt.Println(colony.ErrEmpty) }
`}}, Lab: &levels.Lab{ReqFile: "main.go"}}
	for _, lang := range []string{"en", "ar"} {
		ctx := withLang(context.Background(), lang)
		lab, err := runner.RunLab(ctx, lvl, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !lab.Result.Success || len(lab.Steps) != 1 {
			t.Fatalf("%s lab: %+v", lang, lab)
		}
		for _, message := range []string{lab.Result.Message, lab.Steps[0].Title, lab.Steps[0].Output} {
			arabic := strings.IndexFunc(message, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
			if arabic != (lang == "ar") {
				t.Fatalf("%s lab message: %q", lang, message)
			}
		}
		compiled, err := runner.Run(ctx, levels.All[0], "package main\nfunc main(){ missing() }")
		if err != nil {
			t.Fatal(err)
		}
		if len(compiled.CompileErrors) == 0 {
			t.Fatalf("compile: %+v", compiled)
		}
		tip := compiled.CompileErrors[0].Tip
		arabic := strings.IndexFunc(tip, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
		if tip == "" || arabic != (lang == "ar") {
			t.Fatalf("%s compiler tip: %q", lang, tip)
		}
	}
}

func TestLangPreviewAndMetadataStayRequestScoped(t *testing.T) {
	s, _ := profileServer(t)
	for _, lang := range []string{"en", "ar", "en"} {
		req := httptest.NewRequest("GET", "http://localhost/api/preview?level=1-1-first-signal", nil)
		req.Header.Set("X-Orbital-Lang", lang)
		rec := httptest.NewRecorder()
		s.preview(rec, req)
		if rec.Code != 200 {
			t.Fatalf("preview: %d %s", rec.Code, rec.Body)
		}
		var event struct{ Progress string }
		if err := json.Unmarshal(rec.Body.Bytes(), &event); err != nil {
			t.Fatal(err)
		}
		arabic := strings.IndexFunc(event.Progress, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
		if event.Progress == "" || arabic != (lang == "ar") {
			t.Fatalf("%s preview: %s", lang, rec.Body)
		}
		rec = httptest.NewRecorder()
		s.levels(rec, req)
		var catalog struct{ Levels []struct{ Title string } }
		if err := json.Unmarshal(rec.Body.Bytes(), &catalog); err != nil {
			t.Fatal(err)
		}
		if len(catalog.Levels) == 0 {
			t.Fatal("missing level metadata")
		}
		arabic = strings.IndexFunc(catalog.Levels[0].Title, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
		if arabic != (lang == "ar") {
			t.Fatalf("%s metadata title: %s", lang, catalog.Levels[0].Title)
		}
	}
}

func TestLangColonyShopAndErrors(t *testing.T) {
	s, h := profileServer(t)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/colony/buy", s.colonyBuy)
	buy := s.guard(mux)
	for _, lang := range []string{"en", "ar"} {
		for _, tc := range []struct {
			handler            http.Handler
			method, path, body string
		}{{h, "GET", "/api/colony", ""}, {buy, "POST", "/api/colony/buy", `{"id":"engine"}`}, {buy, "POST", "/api/colony/buy", `{"id":"missing"}`}} {
			req := httptest.NewRequest(tc.method, "http://localhost"+tc.path, strings.NewReader(tc.body))
			req.Header.Set("X-Orbital-Token", "test-token")
			req.Header.Set("X-Orbital-Lang", lang)
			rec := httptest.NewRecorder()
			tc.handler.ServeHTTP(rec, req)
			var response struct {
				Error string
				Shop  []econ.ShopItem
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if rec.Code != 200 {
				t.Fatalf("%s: %s", tc.path, rec.Body)
			}
			messages := []string{}
			if response.Error != "" {
				messages = append(messages, response.Error)
			}
			for _, u := range response.Shop {
				messages = append(messages, u.Name, u.Desc, u.Group)
				if u.Reason != "" {
					messages = append(messages, u.Reason)
				}
			}
			if len(messages) == 0 {
				t.Fatalf("empty %s response: %s", tc.path, rec.Body)
			}
			for _, message := range messages {
				arabic := strings.IndexFunc(message, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
				if lang == "en" && arabic {
					t.Fatalf("English colony text contains Arabic: %q", message)
				}
			}
		}
	}
}

func TestLangConcurrentMetadataRequests(t *testing.T) {
	s, _ := profileServer(t)
	var requests sync.WaitGroup
	for i := 0; i < 8; i++ {
		lang := []string{"en", "ar"}[i%2]
		requests.Add(1)
		go func() {
			defer requests.Done()
			req := httptest.NewRequest("GET", "http://localhost/api/levels", nil)
			req.Header.Set("X-Orbital-Lang", lang)
			rec := httptest.NewRecorder()
			s.levels(rec, req)
			var response struct{ Levels []struct{ Title string } }
			if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
				t.Error(err)
				return
			}
			if len(response.Levels) == 0 {
				t.Error("missing metadata")
				return
			}
			title := response.Levels[0].Title
			arabic := strings.IndexFunc(title, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
			if arabic != (lang == "ar") {
				t.Errorf("concurrent %s title: %q", lang, title)
			}
		}()
	}
	requests.Wait()
}

func TestLangColonyStreamUsesSubscriberLanguage(t *testing.T) {
	for _, runLang := range []string{"ar", "en"} {
		t.Run(runLang, func(t *testing.T) {
			s, h := profileServer(t)
			streamLang := map[string]string{"ar": "en", "en": "ar"}[runLang]
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			req := httptest.NewRequest("GET", "http://localhost/api/colony/stream?token=test-token&profile=default&lang="+streamLang, nil).WithContext(ctx)
			rec := &colonyStreamRecorder{ResponseRecorder: httptest.NewRecorder(), flushed: make(chan struct{}, 8)}
			done := make(chan struct{})
			go func() { defer close(done); h.ServeHTTP(rec, req) }()
			defer func() { cancel(); <-done }()
			select {
			case <-rec.flushed:
			case <-ctx.Done():
				t.Fatal("stream did not open")
			}
			code := `package main
import "orbital/colony"
func main(){for {colony.Measure();colony.Wait()}}
`
			body, _ := json.Marshal(map[string]string{"code": code})
			runReq := httptest.NewRequest("POST", "http://localhost/api/colony/run", strings.NewReader(string(body)))
			runReq.Header.Set("X-Orbital-Token", "test-token")
			runReq.Header.Set("X-Orbital-Lang", runLang)
			started := httptest.NewRecorder()
			h.ServeHTTP(started, runReq)
			if started.Code != 200 || !strings.Contains(started.Body.String(), `"running":true`) {
				t.Fatalf("start: %s", started.Body)
			}
			// Live debugger frames share the stream, so let several flushes
			// through before reading what the colony itself said.
			for frame := 0; frame < 12; frame++ {
				quiet := time.NewTimer(1500 * time.Millisecond)
				select {
				case <-rec.flushed:
					quiet.Stop()
				case <-quiet.C:
					frame = 12
				case <-ctx.Done():
					t.Fatal("missing colony events")
				}
			}
			cancel()
			<-done
			sawGate, sawProgress := false, false
			for _, line := range strings.Split(rec.Body.String(), "\n") {
				if !strings.HasPrefix(line, "data: ") {
					continue
				}
				var event struct {
					Msg, Progress string
					Translations  json.RawMessage
				}
				if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &event); err != nil {
					t.Fatal(err)
				}
				if len(event.Translations) > 0 {
					t.Fatal("private translation metadata leaked into SSE")
				}
				for _, text := range []string{event.Msg, event.Progress} {
					if text == "" {
						continue
					}
					arabic := strings.IndexFunc(text, func(r rune) bool { return unicode.Is(unicode.Arabic, r) }) >= 0
					if arabic != (streamLang == "ar") {
						t.Fatalf("run=%s stream=%s text=%q", runLang, streamLang, text)
					}
				}
				sawGate = sawGate || strings.Contains(event.Msg, "🔒")
				sawProgress = sawProgress || event.Progress != ""
			}
			if !sawGate || !sawProgress {
				t.Fatalf("missing gate or progress: %s", rec.Body)
			}
			// A language change on the stream leaves the same colony process running.
			c, err := s.profileGame(req)
			if err != nil {
				t.Fatal(err)
			}
			c.mu.Lock()
			running := c.active != nil
			c.mu.Unlock()
			if !running {
				t.Fatal("changing stream language stopped the run")
			}
		})
	}
}
