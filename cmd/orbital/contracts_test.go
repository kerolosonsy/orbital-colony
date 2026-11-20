package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"orbital/colony/levels"
)

type contractReply struct {
	OK        bool           `json:"ok"`
	Error     string         `json:"error"`
	Phase     string         `json:"phase"`
	Missing   []string       `json:"missing"`
	Contracts []contractView `json:"contracts"`
}

func contractPost(t *testing.T, handler http.HandlerFunc, body, lang string) contractReply {
	t.Helper()
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "http://localhost/", strings.NewReader(body))
	req.Header.Set("X-Orbital-Lang", lang)
	handler(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("HTTP %d: %s", recorder.Code, recorder.Body)
	}
	var reply contractReply
	if err := json.Unmarshal(recorder.Body.Bytes(), &reply); err != nil {
		t.Fatal(err)
	}
	return reply
}

func contractStatus(list []contractView, id string) string {
	for _, c := range list {
		if c.ID == id {
			return c.Status
		}
	}
	return ""
}

func TestColonyContractEndpointAndRequires(t *testing.T) {
	s, _ := profileServer(t)
	if _, err := s.profiles.ensureDefault(); err != nil {
		t.Fatal(err)
	}
	reply := contractPost(t, s.colonyContract, `{"id":"c04-right-ground"}`, "en")
	if reply.OK || reply.Error != "Finish chapter 4 first" || contractStatus(reply.Contracts, "c04-right-ground") != "locked" {
		t.Fatalf("locked contract accepted: %+v", reply)
	}
	// finish chapters 1-4 in the profile's progress
	records := map[string]map[string]int{}
	for _, l := range levels.All {
		if l.Chapter <= 4 {
			records[l.ID] = map[string]int{"stars": 1}
		}
	}
	blob, _ := json.Marshal(map[string]any{"progress": records})
	if err := os.WriteFile(filepath.Join(s.profiles.dir("default"), "progress.json"), blob, 0o644); err != nil {
		t.Fatal(err)
	}
	reply = contractPost(t, s.colonyContract, `{"id":"c04-right-ground"}`, "ar")
	if !reply.OK || contractStatus(reply.Contracts, "c04-right-ground") != "active" || contractStatus(reply.Contracts, "c05-careful-harvest") != "locked" {
		t.Fatalf("accept: %+v", reply)
	}
	reply = contractPost(t, s.colonyContract, `{"id":"c03-full-sweep"}`, "en")
	if reply.OK || !strings.Contains(reply.Error, "already have an active contract") {
		t.Fatalf("second contract: %+v", reply)
	}
	// the starter has no switch, so the run stops before the build
	body, _ := json.Marshal(map[string]string{"code": levels.Colony.Starter})
	reply = contractPost(t, s.colonyRun, string(body), "en")
	if reply.Phase != "requires" || len(reply.Missing) != 1 {
		t.Fatalf("requires: %+v", reply)
	}
	game, err := s.profileGame(httptest.NewRequest("GET", "/", nil))
	if err != nil {
		t.Fatal(err)
	}
	if game.active != nil {
		t.Fatal("a run started despite missing requirements")
	}
	reply = contractPost(t, s.colonyContract, `{"id":""}`, "en")
	if !reply.OK || game.state.Active != "" || contractStatus(reply.Contracts, "c04-right-ground") != "available" {
		t.Fatalf("cancel: %+v", reply)
	}
	if saved := colonySavedState(t, game.path); saved.Active != "" || saved.Contracts["c04-right-ground"].Tries != 1 {
		t.Fatalf("cancel not saved: %+v", saved.Contracts)
	}
	// the shop explains the missing blueprint in the request language
	summary := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Orbital-Lang", "en")
	s.colonyStatus(summary, req)
	if !strings.Contains(summary.Body.String(), `"reason":"Finish the contract first: The Right Ground"`) || !strings.Contains(summary.Body.String(), `"title":"Price Watch"`) {
		t.Fatalf("summary not localized: %.400s", summary.Body)
	}
}

func TestContractRequirementAlternatives(t *testing.T) {
	recursive := "package main\nfunc f(n int) int { if n == 0 { return 0 }; return f(n-1) }\nfunc main() { f(3) }"
	if m := contractRequirements(recursive, []string{"recursion|sel:heap.Push"}, "en"); len(m) != 0 {
		t.Fatalf("recursion did not satisfy the alternative: %v", m)
	}
	if m := contractRequirements("package main\nfunc main() {}", []string{"recursion|sel:heap.Push"}, "en"); len(m) != 1 || !strings.Contains(m[0], " or ") {
		t.Fatalf("missing alternatives = %v", m)
	}
}

func TestColonyClaimNeedsFinishedChapter(t *testing.T) {
	s, _ := profileServer(t)
	if _, err := s.profiles.ensureDefault(); err != nil {
		t.Fatal(err)
	}
	claim := func() (ok bool, errText string) {
		t.Helper()
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "http://localhost/", strings.NewReader(`{"chapter":1}`))
		req.Header.Set("X-Orbital-Lang", "en")
		s.colonyClaim(recorder, req)
		if recorder.Code != http.StatusOK {
			t.Fatalf("HTTP %d: %s", recorder.Code, recorder.Body)
		}
		var reply struct {
			OK    bool   `json:"ok"`
			Error string `json:"error"`
		}
		if err := json.Unmarshal(recorder.Body.Bytes(), &reply); err != nil {
			t.Fatal(err)
		}
		return reply.OK, reply.Error
	}
	if ok, msg := claim(); ok || msg != "Finish chapter 1 first" {
		t.Fatalf("unfinished chapter claimed: %v %q", ok, msg)
	}
	records := map[string]map[string]int{}
	for _, l := range levels.All {
		if l.Chapter == 1 {
			records[l.ID] = map[string]int{"stars": 1}
		}
	}
	blob, _ := json.Marshal(map[string]any{"progress": records})
	if err := os.WriteFile(filepath.Join(s.profiles.dir("default"), "progress.json"), blob, 0o644); err != nil {
		t.Fatal(err)
	}
	if ok, msg := claim(); !ok {
		t.Fatalf("finished chapter refused: %q", msg)
	}
	if ok, _ := claim(); ok {
		t.Fatal("chapter claimed twice")
	}
}
