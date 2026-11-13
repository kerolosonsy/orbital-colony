package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

func TestSeasonEndpointsIsolationResumeAndBestPersistence(t *testing.T) {
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Inv["geode"] = 47
	c.state.Levels["engine"] = 8
	if err := c.persist(); err != nil {
		t.Fatal(err)
	}
	mainPath := c.path
	before, _ := os.ReadFile(mainPath)
	s := &server{colony: c}
	post := func(id string) *httptest.ResponseRecorder {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/colony/season", strings.NewReader(`{"id":"`+id+`"}`))
		s.colonySeason(rec, req)
		return rec
	}
	if rec := post("../bad"); rec.Code != 400 {
		t.Fatalf("invalid: %d", rec.Code)
	}
	for _, season := range econ.Seasons {
		if rec := post(season.ID); rec.Code != 200 {
			t.Fatalf("select: %d %s", rec.Code, rec.Body)
		}
		if c.state.Seed != season.Seed || c.path != filepath.Join(filepath.Dir(mainPath), "season-"+season.ID+".json") {
			t.Fatal("wrong season save")
		}
		c.state.Tick = 400
		c.state.Earned = 123
		if rec := post(""); rec.Code != 200 {
			t.Fatal(rec.Body)
		}
		if rec := post(season.ID); rec.Code != 200 || c.state.Tick != 400 {
			t.Fatal("attempt did not resume")
		}
		c.state.Tick = season.Ticks
		for _, score := range []float64{123, 50, 200} {
			c.state.Earned = score
			var ev struct {
				OK   bool
				Best float64
			}
			json.Unmarshal([]byte(c.finishSeason("en")), &ev)
			if !ev.OK || ev.Best < score {
				t.Fatalf("completion: %+v", ev)
			}
		}
		if rec := post(""); rec.Code != 200 {
			t.Fatal(rec.Body)
		}
	}
	after, _ := os.ReadFile(mainPath)
	if !bytes.Equal(before, after) || c.state.Inv["geode"] != 47 {
		t.Fatal("main colony changed")
	}
	fresh, err := newColonyGame(filepath.Dir(mainPath))
	if err != nil {
		t.Fatal(err)
	}
	bests, err := fresh.seasonBests()
	if err != nil {
		t.Fatal(err)
	}
	for _, season := range econ.Seasons {
		if bests[seasonBestKey(&season)] != 200 {
			t.Fatal("best not persisted")
		}
	}
	for _, lang := range []string{"ar", "en"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/colony/seasons", nil)
		req.Header.Set("X-Orbital-Lang", lang)
		req.URL.RawQuery = "lang=" + lang
		s.colonySeasons(rec, req)
		var view struct {
			Seasons []struct {
				ID, Title, Desc string
				Best, Reference float64
				Ticks           int64
			}
			Active any
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &view); err != nil {
			t.Fatal(err)
		}
		if len(view.Seasons) != 6 || view.Active != nil {
			t.Fatal(rec.Body)
		}
	}
	c.active = &colonyRun{}
	if rec := post("one-drone"); rec.Code != 409 {
		t.Fatal("running colony switched season")
	}
	c.active = nil
}

func playSeason(t *testing.T, prog string, season *econ.Season) float64 {
	t.Helper()
	save := filepath.Join(t.TempDir(), "season.json")
	body, _ := json.Marshal(season.New())
	if err := os.WriteFile(save, body, 0600); err != nil {
		t.Fatal(err)
	}
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, prog)
	cmd.Env = append(os.Environ(), "GORACE=atexit_sleep_ms=0", "ORBITAL_LANG=en", "ORBITAL_LEVEL=colony", "ORBITAL_COLONY="+save, "ORBITAL_EVENTS_FD=3", "ORBITAL_COLONY_UNIT=2ms")
	cmd.ExtraFiles = []*os.File{writer}
	var diagnostics bytes.Buffer
	cmd.Stderr = &diagnostics
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	writer.Close()
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	score := -1.0
	for scanner.Scan() {
		var event struct {
			Op    string
			State *econ.State
		}
		if json.Unmarshal(scanner.Bytes(), &event) == nil && event.Op == "season-end" && event.State != nil {
			if event.State.Tick != season.Ticks {
				t.Errorf("ended at %d ticks", event.State.Tick)
			}
			score = event.State.Earned
		}
	}
	err = cmd.Wait()
	if err != nil || score < 0 || strings.Contains(diagnostics.String(), "DATA RACE") {
		t.Fatalf("season %s: score=%f, %v\n%s", season.ID, score, err, diagnostics.String())
	}
	return score
}

func TestSeasonsSolvableAndSingleDroneDeterministic(t *testing.T) {
	r := testRunner(t)
	starter := buildColonyProgram(t, r, levels.Colony.Starter)
	for i := range econ.Seasons {
		season := &econ.Seasons[i]
		t.Run(season.ID, func(t *testing.T) {
			t.Parallel()
			body, err := os.ReadFile(filepath.Join("testdata", "seasons", season.ID, "main.go"))
			if err != nil {
				t.Fatal(err)
			}
			ref := buildColonyProjectProgram(t, r, colonyProject{Files: map[string]string{"main.go": string(body)}, Race: true})
			score := playSeason(t, ref, season)
			baseline := playSeason(t, starter, season)
			t.Logf("%s reference=%.2f stored=%.2f starter=%.2f", season.ID, score, season.Reference, baseline)
			if score < baseline*1.15+10 {
				t.Errorf("reference needs a clear margin over starter")
			}
			if math.Abs(score-season.Reference) > math.Max(10, season.Reference*0.08) {
				t.Errorf("reference score drift: measured %.2f stored %.2f", score, season.Reference)
			}
			if season.ID == "one-drone" {
				if again := playSeason(t, ref, season); again != score {
					t.Errorf("not deterministic: %.10f vs %.10f", score, again)
				}
			}
		})
	}
}

func TestSeasonRunCompletionLeavesMainSaveUntouched(t *testing.T) {
	r := testRunner(t)
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mainPath := c.path
	before, _ := os.ReadFile(mainPath)
	if err := c.selectSeason("one-drone"); err != nil {
		t.Fatal(err)
	}
	c.state.Tick = 5980
	dir, err := r.colonyProjectDir(map[string]string{"main.go": `package main;import "orbital/colony";func main(){for{colony.Wait()}}`})
	if err != nil {
		t.Fatal(err)
	}
	if failure := r.buildColonyProject(withLang(t.Context(), "en"), dir, &colonyProject{Race: true}); failure != nil {
		t.Fatal(failure)
	}
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en", env: []string{"ORBITAL_COLONY_UNIT=2ms", "GORACE=atexit_sleep_ms=0"}}); err != nil {
		t.Fatal(err)
	}
	c.mu.Lock()
	done := c.active.done
	c.mu.Unlock()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		c.stop()
		t.Fatal("season didn't stop")
	}
	found := false
	for len(ch) > 0 {
		var event struct {
			Op          string
			OK          bool
			Score, Best float64
		}
		if err := json.Unmarshal([]byte(<-ch), &event); err != nil {
			t.Fatal(err)
		}
		if event.Op == "season" {
			found = true
			if !event.OK || event.Score != 0 || event.Best != 0 {
				t.Fatalf("completion %+v", event)
			}
		}
	}
	if !found {
		t.Fatal("missing season completion event")
	}
	if c.state.Tick != 6000 {
		t.Fatalf("final tick %d", c.state.Tick)
	}
	c.mu.Lock()
	err = c.selectSeason("")
	c.mu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(mainPath)
	if !bytes.Equal(before, after) {
		t.Fatal("season run changed main save")
	}
}
