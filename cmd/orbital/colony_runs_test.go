package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"orbital/colony/econ"
)

func TestColonyRunsContentsStatsCapAndNewestFirst(t *testing.T) {
	s := testColony(t)
	c := s.colony
	c.state.Tick = 50
	c.state.Earned = 100
	files := map[string]string{"main.go": "package main\nfunc main(){}", "plan/note.txt": "plan"}
	details := &colonyRunDetails{lang: "en", record: colonyRunRecord{CodeHash: colonyCodeHash(files)}}
	c.active = &colonyRun{details: details}
	c.beginRun(details)
	details.record.Started = time.Now().UTC().Add(-2 * time.Minute)
	c.event(`{"op":"econ","stats":{"drones":[{"id":1,"commands":{"move":3}}]}}`)
	c.event(`{"op":"econ"}`)
	c.event(`{"op":"econ","stats":null}`)
	c.state.Tick = 1284
	c.state.Earned = 660.5
	c.finishRun(c.active)
	c.active = nil
	response := httptest.NewRecorder()
	s.colonyRuns(response, httptest.NewRequest("GET", "/api/colony/runs", nil))
	var reply struct {
		Runs []colonyRunRecord `json:"runs"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &reply); err != nil {
		t.Fatal(err)
	}
	if len(reply.Runs) != 1 {
		t.Fatalf("history: %s", response.Body)
	}
	run := reply.Runs[0]
	if run.Ticks != 1234 || run.Earned != 560.5 || run.PerMin < 280 || run.PerMin > 280.3 || run.CodeHash != colonyCodeHash(files) || run.Contract != "" || run.Race || !strings.Contains(string(run.Stats), `"move":3`) {
		t.Fatalf("record: %+v", run)
	}
	var raw struct{ Runs []map[string]json.RawMessage }
	if err := json.Unmarshal(response.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	var started string
	if err := json.Unmarshal(raw.Runs[0]["started"], &started); err != nil {
		t.Fatal(err)
	}
	if _, err := time.Parse(time.RFC3339, started); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 55; i++ {
		if err := c.appendRun(colonyRunRecord{Ticks: int64(i)}); err != nil {
			t.Fatal(err)
		}
	}
	reloaded, err := newColonyGame(filepath.Dir(c.path))
	if err != nil {
		t.Fatal(err)
	}
	runs, err := reloaded.readRuns()
	if err != nil || len(runs) != 50 || runs[0].Ticks != 54 || runs[49].Ticks != 5 {
		t.Fatalf("cap/order: %v %v", runs, err)
	}
	if string(runs[0].Stats) != "null" {
		t.Fatalf("absent stats: %s", runs[0].Stats)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(c.path), "colony-runs.json")); err != nil {
		t.Fatal(err)
	}
}

func TestColonyRunsHashIncludesNamesAndContents(t *testing.T) {
	first := colonyCodeHash(map[string]string{"main.go": "main", "plan/p.go": "plan"})
	reordered := colonyCodeHash(map[string]string{"plan/p.go": "plan", "main.go": "main"})
	if first != reordered || len(first) != 64 {
		t.Fatalf("unstable hash: %s %s", first, reordered)
	}
	for _, files := range []map[string]string{{"main.go": "main!", "plan/p.go": "plan"}, {"main.go": "main", "plan/q.go": "plan"}, {"main.go": "mainplan"}} {
		if first == colonyCodeHash(files) {
			t.Fatalf("hash omitted change: %v", files)
		}
	}
}

func TestColonyRaceReportedOnceAndStored(t *testing.T) {
	s := testColony(t)
	ch, unsubscribe := s.colony.subscribe()
	defer unsubscribe()
	source := `package main
import ("sync";"fmt")
func main(){
 var x int
 var wg sync.WaitGroup
 for i:=0;i<8;i++ {wg.Add(1);go func(){defer wg.Done();for k:=0;k<10000;k++{x++}}()}
 wg.Wait()
 fmt.Println(x)
}`
	response := projectPost(t, s.colonyRun, colonyProject{Code: source, Race: true})
	if !strings.Contains(response.Body.String(), `"phase":"run"`) {
		t.Fatal(response.Body)
	}
	timer := time.NewTimer(30 * time.Second)
	defer timer.Stop()
	races := 0
	for {
		select {
		case line := <-ch:
			var ev struct{ Op, Msg string }
			if err := json.Unmarshal([]byte(line), &ev); err != nil {
				t.Fatal(err)
			}
			if ev.Op == "race" {
				var arabic struct{ Msg string }
				if err := json.Unmarshal([]byte(localizeStreamLine(line, "ar")), &arabic); err != nil {
					t.Fatal(err)
				}
				if arabic.Msg != runtimeTipFor("WARNING: DATA RACE", "ar") {
					t.Fatalf("Arabic race tip: %s", arabic.Msg)
				}
				races++
				if ev.Msg != runtimeTipFor("WARNING: DATA RACE", "en") {
					t.Fatalf("race tip: %s", line)
				}
			}
			if ev.Op == "stopped" {
				if races != 1 {
					t.Fatalf("race events: %d", races)
				}
				// Stop also waits for the history and directory cleanup to finish.
				colonyRequest(t, s.colonyStop, "")
				runs, err := s.colony.readRuns()
				if err != nil || len(runs) != 1 || !runs[0].Race {
					t.Fatalf("race history: %v %v", runs, err)
				}
				return
			}
		case <-timer.C:
			t.Fatal("race run did not finish")
		}
	}
}

func TestColonyContractC12RequiresRaceAndRacedReplayKeepsOwnership(t *testing.T) {
	for _, lang := range []string{"ar", "en"} {
		t.Run(lang, func(t *testing.T) {
			s := testColony(t)
			s.colony.state = contractFixture(econ.FindContract("c12-one-ledger"))

			source := "package main\nimport \"sync\"\nvar mu sync.Mutex\nfunc main(){}"
			blob, _ := json.Marshal(colonyProject{Code: source})
			body := string(blob)
			reply := contractPost(t, s.colonyRun, body, lang)
			if reply.Phase != "requires" || len(reply.Missing) != 1 || !strings.Contains(reply.Missing[0], "race:true") {
				t.Fatalf("race gate: %+v", reply)
			}
			if s.colony.active != nil {
				t.Fatal("started without race")
			}
		})
	}
	for _, best := range []int64{0, 123} {
		t.Run(fmt.Sprint(best), func(t *testing.T) {
			s := testColony(t)
			c := s.colony
			c.state = contractFixture(econ.FindContract("c12-one-ledger"))
			c.state.Contracts[c.state.Active].Best = best
			c.state.Levels["market.logistics"] = 3
			details := &colonyRunDetails{lang: "en"}
			c.beginRun(details)
			c.active = &colonyRun{details: details}
			run := c.state.Contracts[c.state.Active]
			run.Best = 50
			run.Status = "done"
			c.state.Active = ""
			if best == 0 {
				c.state.Levels["market.logistics"]++
			}
			ch, unsubscribe := c.subscribe()
			defer unsubscribe()
			c.detectRace("WARNING: DATA RACE")
			c.detectRace("WARNING: DATA RACE")
			if len(ch) != 1 {
				t.Fatalf("duplicate race events: %d", len(ch))
			}
			event := c.finishRun(c.active)
			c.active = nil
			if c.state.Levels["market.logistics"] != 3 || run.Best != best || best == 0 && run.Status != "failed" || best > 0 && run.Status != "done" || !strings.Contains(event, `"ok":false`) {
				t.Fatalf("raced contract lost ownership: %+v %s", c.state, event)
			}
		})
	}
}

func TestColonyContractC12CompletionWaitsForRaceVerdict(t *testing.T) {
	for _, raced := range []bool{false, true} {
		t.Run(fmt.Sprint(raced), func(t *testing.T) {
			c, err := newColonyGame(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			c.state = contractFixture(econ.FindContract("c12-one-ledger"))
			details := &colonyRunDetails{lang: "en"}
			c.beginRun(details)
			c.active = &colonyRun{details: details}
			ch, unsubscribe := c.subscribe()
			defer unsubscribe()
			run := c.state.Contracts[details.record.Contract]
			run.Status = "done"
			run.Best = 100
			c.state.Active = ""
			c.state.Levels["market.logistics"]++
			event, _ := json.Marshal(map[string]any{"op": "contract", "ok": true, "contract": map[string]string{"id": details.record.Contract}, "state": c.state})
			c.event(string(event))
			if len(ch) != 0 {
				t.Fatal("success emitted before race detector exited")
			}
			if raced {
				c.detectRace("WARNING: DATA RACE")
			}
			c.state.Credits = 1234
			final := c.finishRun(c.active)
			var verdict struct {
				OK    bool
				State *econ.State
			}
			if err := json.Unmarshal([]byte(final), &verdict); err != nil {
				t.Fatal(err)
			}
			if verdict.OK == raced || verdict.State.Credits != 1234 || raced && verdict.State.Contracts[details.record.Contract].Status != "failed" {
				t.Fatalf("final verdict: %s", final)
			}
		})
	}
}
