package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

func testColony(t *testing.T) *server {
	t.Helper()
	game, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { game.ops.Lock(); defer game.ops.Unlock(); game.stop() })
	return &server{run: testRunner(t), colony: game, token: "colony-test-token"}
}

func colonyRequest(t *testing.T, handler http.HandlerFunc, body string) *httptest.ResponseRecorder {
	t.Helper()
	recorder := httptest.NewRecorder()
	handler(recorder, httptest.NewRequest(http.MethodPost, "http://localhost/", strings.NewReader(body)))
	if recorder.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", recorder.Code, recorder.Body)
	}
	return recorder
}

func colonySavedState(t *testing.T, path string) *econ.State {
	t.Helper()
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var state econ.State
	if err := json.Unmarshal(body, &state); err != nil {
		t.Fatal(err)
	}
	return &state
}

func TestColonyStoppedSummaryDrawsSavedBuildings(t *testing.T) {
	s := testColony(t)
	c := s.colony
	c.state.Buildings = []*econ.Building{
		{Kind: "composter", X: 2, Y: 1, Level: 1, Enabled: true},
		{Kind: "generator", X: 0, Y: 2, Level: 1, Enabled: true, Mode: "solar"},
	}
	c.state.At(2, 1).G = "soil"
	if err := c.persist(); err != nil {
		t.Fatal(err)
	}
	reloaded, err := newColonyGame(filepath.Dir(c.path))
	if err != nil {
		t.Fatal(err)
	}
	s.colony = reloaded
	response := httptest.NewRecorder()
	s.colonyStatus(response, httptest.NewRequest(http.MethodGet, "/api/colony", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", response.Code, response.Body)
	}
	var summary struct {
		Running bool
		State   struct {
			Size      int
			Buildings []*econ.Building
			Tiles     []struct {
				econ.Tile
				B string `json:"b"`
			}
		}
	}
	if err := json.Unmarshal(response.Body.Bytes(), &summary); err != nil {
		t.Fatal(err)
	}
	if summary.Running || summary.State.Size != 3 || len(summary.State.Tiles) != 9 || len(summary.State.Buildings) != 2 {
		t.Fatalf("stopped summary lost saved state: %s", response.Body)
	}
	for i, tile := range summary.State.Tiles {
		want := ""
		switch i {
		case 1*3 + 2:
			want = "composter"
		case 2 * 3:
			want = "generator"
		}
		if tile.B != want || tile.Tile != reloaded.state.Tiles[i] {
			t.Errorf("tile %d = %+v, want saved tile with b=%q", i, tile, want)
		}
	}
}

func TestColonyStoppedPurchasesAndClaimsPersist(t *testing.T) {
	s := testColony(t)
	c := s.colony
	fresh := colonySavedState(t, c.path)
	if fresh.Size != 3 || fresh.Tick != 0 {
		t.Fatalf("unexpected fresh save: %+v", fresh)
	}
	c.state.Credits = 100
	c.state.Inv["moss"] = 5
	response := colonyRequest(t, s.colonyBuy, `{"id":"engine"}`)
	var purchase struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &purchase); err != nil {
		t.Fatal(err)
	}
	saved := colonySavedState(t, c.path)
	if !purchase.OK || saved.Credits != 75 || saved.Levels["engine"] != 1 {
		t.Fatalf("purchase: %s; save: %+v", response.Body, saved)
	}
	response = colonyRequest(t, s.colonyClaim, `{"chapter":1}`)
	var claim struct {
		OK     bool        `json:"ok"`
		Error  string      `json:"error"`
		Reward econ.Reward `json:"reward"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &claim); err != nil {
		t.Fatal(err)
	}
	saved = colonySavedState(t, c.path)
	if !claim.OK || saved.Credits != 75+claim.Reward.Credits || !saved.Claimed["chapter-1"] {
		t.Fatalf("claim: %s; save: %+v", response.Body, saved)
	}
	credits := saved.Credits
	response = colonyRequest(t, s.colonyClaim, `{"chapter":1}`)
	if err := json.Unmarshal(response.Body.Bytes(), &claim); err != nil {
		t.Fatal(err)
	}
	if claim.OK || claim.Error == "" || colonySavedState(t, c.path).Credits != credits {
		t.Fatalf("repeat claim: %s", response.Body)
	}
	reloaded, err := newColonyGame(filepath.Dir(c.path))
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.state.Credits != credits || reloaded.state.Levels["engine"] != 1 {
		t.Fatalf("reload lost changes: %+v", reloaded.state)
	}
}

func TestColonyRealRunStreamsAndSavesOnStop(t *testing.T) {
	s := testColony(t)
	c := s.colony
	c.state.Credits = 100
	c.state.Inv["moss"] = 5
	startingTick := c.state.Tick
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	body, err := json.Marshal(map[string]string{"code": levels.Colony.Starter})
	if err != nil {
		t.Fatal(err)
	}
	response := colonyRequest(t, s.colonyRun, string(body))
	var run struct {
		Phase   string
		Running bool
	}
	if err := json.Unmarshal(response.Body.Bytes(), &run); err != nil {
		t.Fatal(err)
	}
	if run.Phase != "run" || !run.Running {
		t.Fatalf("start: %s", response.Body)
	}
	c.mu.Lock()
	active := c.active
	c.mu.Unlock()
	if active == nil {
		t.Fatal("program exited before streaming")
	}
	timer := time.NewTimer(20 * time.Second)
	defer timer.Stop()
	initSeen, econSeen := false, false
	for !initSeen || !econSeen {
		select {
		case line := <-ch:
			var event struct{ Op string }
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				t.Fatal(err)
			}
			initSeen = initSeen || event.Op == "init"
			econSeen = econSeen || event.Op == "econ"
		case <-timer.C:
			t.Fatalf("stream timed out: init=%v econ=%v", initSeen, econSeen)
		}
	}
	purchase := colonyRequest(t, s.colonyBuy, `{"id":"engine"}`)
	var pending struct{ OK, Pending bool }
	if err := json.Unmarshal(purchase.Body.Bytes(), &pending); err != nil {
		t.Fatal(err)
	}
	if !pending.OK || !pending.Pending {
		t.Fatalf("live purchase: %s", purchase.Body)
	}
	bought := false
	for !bought {
		select {
		case line := <-ch:
			var event struct {
				Op    string
				State *econ.State
			}
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				t.Fatal(err)
			}
			if event.Op == "bought" {
				if event.State == nil || event.State.Levels["engine"] != 1 {
					t.Fatalf("live purchase not applied: %s", line)
				}
				bought = true
			}
		case <-timer.C:
			t.Fatal("live purchase did not emit bought")
		}
	}
	// Rewards must not race the simulation's state snapshots.
	conflict := httptest.NewRecorder()
	s.colonyClaim(conflict, httptest.NewRequest("POST", "/", strings.NewReader(`{"chapter":1}`)))
	if conflict.Code != http.StatusConflict || !strings.Contains(conflict.Body.String(), "وقّف المستعمرة الأول") {
		t.Fatalf("claim during run: %d %s", conflict.Code, conflict.Body)
	}
	colonyRequest(t, s.colonyStop, "")
	select {
	case <-active.done:
	default:
		t.Fatal("stop returned before exit")
	}
	if !active.cmd.ProcessState.Exited() {
		t.Fatalf("unclean exit: %s", active.cmd.ProcessState)
	}
	stopped := false
	for len(ch) > 0 {
		var event struct {
			Op   string
			Code int
		}
		if err := json.Unmarshal([]byte(<-ch), &event); err != nil {
			t.Fatal(err)
		}
		if event.Op == "stopped" {
			stopped = true
			if event.Code != 0 {
				t.Fatalf("exit code %d", event.Code)
			}
		}
	}
	if !stopped {
		t.Fatal("missing stopped event")
	}
	if saved := colonySavedState(t, c.path); saved.Tick <= startingTick || saved.Levels["engine"] != 1 {
		t.Fatalf("tick did not advance: %d -> %d", startingTick, saved.Tick)
	}
	if _, err := os.Stat(active.cmd.Dir); !os.IsNotExist(err) {
		t.Fatalf("temporary run directory still exists: %v", err)
	}
}

func TestColonyCompileFailure(t *testing.T) {
	s := testColony(t)
	response := colonyRequest(t, s.colonyRun, `{"code":"package main\nfunc main() { missing() }"}`)
	var failure struct {
		Phase         string
		CompileErrors []compileError
		CompileOutput string
	}
	if err := json.Unmarshal(response.Body.Bytes(), &failure); err != nil {
		t.Fatal(err)
	}
	if failure.Phase != "compile" || len(failure.CompileErrors) == 0 || !strings.Contains(failure.CompileOutput, "undefined: missing") {
		t.Fatalf("compile response: %s", response.Body)
	}
}

func TestColonyStreamQueryToken(t *testing.T) {
	s := testColony(t)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/colony/stream", s.colonyStream)
	mux.HandleFunc("GET /api/colony", s.colonyStatus)
	guarded := s.guard(mux)
	for _, path := range []string{"/api/colony/stream?token=wrong", "/api/colony?token=" + s.token} {
		response := httptest.NewRecorder()
		guarded.ServeHTTP(response, httptest.NewRequest("GET", "http://localhost"+path, nil))
		if response.Code != http.StatusForbidden {
			t.Fatalf("query bypassed guard on %s: %d", path, response.Code)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := httptest.NewRequest("GET", "http://localhost/api/colony/stream?token="+s.token, nil).WithContext(ctx)
	response := &colonyStreamRecorder{ResponseRecorder: httptest.NewRecorder(), flushed: make(chan struct{}, 4)}
	done := make(chan struct{})
	go func() { defer close(done); guarded.ServeHTTP(response, req) }()
	defer func() { cancel(); <-done }()
	for frame := 0; frame < 2; frame++ {
		select {
		case <-response.flushed:
		case <-ctx.Done():
			t.Fatal("SSE did not flush before timeout")
		}
		if frame == 0 {
			s.colony.broadcast(`{"op":"test"}`)
		}
	}
	cancel()
	<-done
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "text/event-stream" || response.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("stream headers: %d %v", response.Code, response.Header())
	}
	if response.Body.String() != "data: {\"op\":\"test\"}\n\n" {
		t.Fatalf("SSE frame: %q", response.Body.String())
	}
}

// Observe Flush without opening sockets, which are unavailable in the sandbox.
type colonyStreamRecorder struct {
	*httptest.ResponseRecorder
	flushed chan struct{}
}

func (r *colonyStreamRecorder) Flush() {
	r.ResponseRecorder.Flush()
	r.flushed <- struct{}{}
}

func TestGuardSetsSecurityHeaders(t *testing.T) {
	s := testColony(t)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {})
	for _, path := range []string{"/", "/api/colony"} { // allowed and refused
		response := httptest.NewRecorder()
		s.guard(mux).ServeHTTP(response, httptest.NewRequest("GET", "http://localhost"+path, nil))
		h := response.Header()
		if h.Get("X-Frame-Options") != "DENY" || h.Get("X-Content-Type-Options") != "nosniff" ||
			!strings.Contains(h.Get("Content-Security-Policy"), "frame-ancestors 'none'") || h.Get("Referrer-Policy") != "no-referrer" {
			t.Fatalf("%s headers: %v", path, h)
		}
	}
}

func TestStaticFilesRevalidate(t *testing.T) {
	web := fstest.MapFS{"app.js": {Data: []byte("console.log(1)")}}
	handler := staticHandler(web)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest("GET", "/app.js", nil))
	tag := response.Header().Get("ETag")
	if response.Code != 200 || response.Header().Get("Cache-Control") != "no-cache" || tag == "" {
		t.Fatalf("first load: %d %v", response.Code, response.Header())
	}
	again := httptest.NewRequest("GET", "/app.js", nil)
	again.Header.Set("If-None-Match", tag)
	response = httptest.NewRecorder()
	handler.ServeHTTP(response, again)
	if response.Code != http.StatusNotModified {
		t.Fatalf("unchanged file should be 304, got %d", response.Code)
	}
}

func TestStripDirRemovesResolvedTempPath(t *testing.T) {
	dir := t.TempDir()
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	in := "\t" + real + "/main.go:8 +0x1c\n" + dir + "/util.go:3: oops"
	if got := stripDir(in, dir); got != "\tmain.go:8 +0x1c\nutil.go:3: oops" {
		t.Fatalf("stripDir = %q", got)
	}
}

func TestLevelsPayloadIsCachedWithETag(t *testing.T) {
	s := &server{}
	get := func(lang, etag string) *httptest.ResponseRecorder {
		req := httptest.NewRequest("GET", "/api/levels", nil)
		req.Header.Set("X-Orbital-Lang", lang)
		if etag != "" {
			req.Header.Set("If-None-Match", etag)
		}
		response := httptest.NewRecorder()
		s.levels(response, req)
		return response
	}
	ar, en := get("ar", ""), get("en", "")
	if ar.Code != 200 || en.Code != 200 || ar.Header().Get("ETag") == "" || ar.Header().Get("ETag") == en.Header().Get("ETag") {
		t.Fatalf("ar %d %q, en %d %q", ar.Code, ar.Header().Get("ETag"), en.Code, en.Header().Get("ETag"))
	}
	var payload struct{ Levels []json.RawMessage }
	if err := json.Unmarshal(en.Body.Bytes(), &payload); err != nil || len(payload.Levels) != len(levels.All) {
		t.Fatalf("payload: %v, %d levels", err, len(payload.Levels))
	}
	if again := get("ar", ar.Header().Get("ETag")); again.Code != http.StatusNotModified || again.Body.Len() != 0 {
		t.Fatalf("unchanged payload should be 304, got %d", again.Code)
	}
}
