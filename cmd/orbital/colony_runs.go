package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type colonyRunRecord struct {
	Started  time.Time       `json:"started"`
	Ticks    int64           `json:"ticks"`
	Earned   float64         `json:"earned"`
	PerMin   float64         `json:"perMin"`
	CodeHash string          `json:"codeHash"`
	Contract string          `json:"contract"`
	Race     bool            `json:"race"`
	Stats    json.RawMessage `json:"stats"`
}

type colonyRunDetails struct {
	env              []string
	record           colonyRunRecord
	lang             string
	startTick        int64
	startEarned      float64
	ledgerBest       int64
	ledgerLevel      int
	ledgerCompletion map[string]json.RawMessage
}

func colonyCodeHash(files map[string]string) string {
	hash := sha256.New()
	for _, name := range projectNames(files) {
		// Lengths make boundaries unambiguous even when contents contain file names.
		fmt.Fprintf(hash, "%d:%s%d:%s", len(name), name, len(files[name]), files[name])
	}
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Caller holds mu; this baseline must follow the previous run's final snapshot.
func (c *colonyGame) beginRun(details *colonyRunDetails) {
	details.record.Started = time.Now().UTC()
	details.record.Contract = c.state.Active
	details.startTick, details.startEarned = c.state.Tick, c.state.Earned
	if details.record.Contract == "c12-one-ledger" {
		details.ledgerBest = c.state.Contracts[details.record.Contract].Best
		details.ledgerLevel = c.state.Levels["market.logistics"]
	}
}

func (c *colonyGame) observeRunEvent(line string) {
	var event struct {
		Stats json.RawMessage `json:"stats"`
	}
	if json.Unmarshal([]byte(line), &event) != nil || len(event.Stats) == 0 || string(event.Stats) == "null" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active != nil {
		c.active.details.record.Stats = event.Stats
	}
}

func (c *colonyGame) deferLedgerCompletion(line string) bool {
	var event struct {
		Op       string
		OK       bool
		Contract struct{ ID string }
	}
	if json.Unmarshal([]byte(line), &event) != nil || event.Op != "contract" || !event.OK || event.Contract.ID != "c12-one-ledger" {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active == nil || c.active.details.record.Contract != event.Contract.ID {
		return false
	}
	json.Unmarshal([]byte(line), &c.active.details.ledgerCompletion)
	return true
}

func (c *colonyGame) detectRace(line string) {
	if !strings.Contains(line, "WARNING: DATA RACE") {
		return
	}
	c.mu.Lock()
	active := c.active
	if active == nil || active.details.record.Race {
		c.mu.Unlock()
		return
	}
	active.details.record.Race = true
	lang := active.details.lang
	c.mu.Unlock()
	ar, en := runtimeTipFor(line, "ar"), runtimeTipFor(line, "en")
	encoded, _ := json.Marshal(map[string]any{"op": "race", "msg": tr(lang, ar, en), "translations": [][2]string{{ar, en}}})
	c.broadcast(string(encoded))
}

// The race detector may report during process teardown, after the sim settled
// c12. Only that run's new award is revoked; an earlier win stays owned.
func (c *colonyGame) failRacedLedger(details *colonyRunDetails) string {
	if details.record.Contract != "c12-one-ledger" || !details.record.Race {
		return ""
	}
	run := c.state.Contracts[details.record.Contract]
	if run.Best > 0 && details.ledgerBest == 0 {
		c.state.Levels["market.logistics"] = max(details.ledgerLevel, c.state.Levels["market.logistics"]-1)
	}
	run.Best = details.ledgerBest
	run.Status = "failed"
	if run.Best > 0 {
		run.Status = "done"
	}
	if c.state.Active == details.record.Contract {
		c.state.Active = ""
	}
	ar, en := "العقد فشل: اتكشف data race.", "Contract failed: a data race was detected."
	encoded, _ := json.Marshal(map[string]any{"op": "contract", "ok": false, "msg": tr(details.lang, ar, en), "state": c.state, "contract": map[string]string{"id": details.record.Contract}, "translations": [][2]string{{ar, en}}})
	return string(encoded)
}

func (c *colonyGame) ledgerCompletionEvent(details *colonyRunDetails) string {
	if len(details.ledgerCompletion) == 0 {
		return ""
	}
	event := details.ledgerCompletion
	event["state"], _ = json.Marshal(c.state)
	event["credits"], _ = json.Marshal(c.state.Credits)
	encoded, _ := json.Marshal(event)
	return string(encoded)
}

// Caller holds mu after every stream has drained.
func (c *colonyGame) finishRun(active *colonyRun) string {
	details := active.details
	finalEvent := c.ledgerCompletionEvent(details)
	if details.record.Race {
		finalEvent = c.failRacedLedger(details)
	}
	record := details.record
	record.Ticks = c.state.Tick - details.startTick
	record.Earned = c.state.Earned - details.startEarned
	if minutes := time.Since(record.Started).Minutes(); minutes > 0 {
		record.PerMin = record.Earned / minutes
	}
	if err := c.appendRun(record); err != nil {
		log.Printf("colony run history: %v", err)
	}
	return finalEvent
}

func (c *colonyGame) runsPath() string {
	return filepath.Join(filepath.Dir(c.path), "colony-runs.json")
}

func (c *colonyGame) readRuns() ([]colonyRunRecord, error) {
	runs := []colonyRunRecord{}
	body, err := os.ReadFile(c.runsPath())
	if errors.Is(err, os.ErrNotExist) {
		return runs, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(body, &runs); err != nil {
		return nil, err
	}
	if runs == nil {
		runs = []colonyRunRecord{}
	}
	return runs, nil
}

func (c *colonyGame) appendRun(record colonyRunRecord) error {
	runs, err := c.readRuns()
	if err != nil {
		return err
	}
	runs = append([]colonyRunRecord{record}, runs...)
	if len(runs) > 50 {
		runs = runs[:50]
	}
	f, err := os.CreateTemp(filepath.Dir(c.path), ".colony-runs-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	err = json.NewEncoder(f).Encode(runs)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	return os.Rename(f.Name(), c.runsPath())
}

func (s *server) colonyRuns(w http.ResponseWriter, r *http.Request) {
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.mu.Lock()
	runs, err := c.readRuns()
	c.mu.Unlock()
	if err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"runs": runs})
}
