package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"orbital/colony/econ"
)

func archiveMarketContract(id string) bool {
	return strings.HasPrefix(id, "d18-") || strings.HasPrefix(id, "d19-") || strings.HasPrefix(id, "d22-") || id == "h24-bulk-seller" || id == "h33-stubborn-client"
}

func requireMarketSocket(t *testing.T) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
		t.Skip("sandbox prohibits loopback listeners")
	}
	if err != nil {
		t.Fatal(err)
	}
	listener.Close()
}

func prepareArchiveContract(t *testing.T, c *colonyGame, id string) {
	t.Helper()
	if strings.HasPrefix(id, "h") {
		c.state.Levels["sys.market"] = 1
		c.state.Inv["moss"] = 1000
		return
	}
	c.state.Levels["sys.archive"] = 1
	dir := filepath.Join(filepath.Dir(c.path), "colony-data")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if id == "d18-audit" {
		ledger := "tick,good,amount,unit_price,credits\n1,moss,2,1.5,3\n2,lichen,3,2,6\n3,moss,bad,1,4\n4,moss,1,1\n5,\"moss,damaged\",2,1,2\n6,moss,2,2,4\n"
		if err := os.WriteFile(filepath.Join(dir, "ledger.csv"), []byte(ledger), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if id == "d19-shift-report" {
		c.state.Tick = 3600
		events := ""
		for day := 0; day < 3; day++ {
			for interval := 0; interval < 2; interval++ {
				events += fmt.Sprintf("{\"op\":\"harvest\",\"t\":%d,\"day\":%d,\"goods\":{\"moss\":%d,\"water\":1.5}}\n", day*1200+interval*120, day, day+1)
			}
		}
		if err := os.WriteFile(filepath.Join(dir, "events.jsonl"), []byte(events), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}

func startFastSystemProgram(t *testing.T, c *colonyGame, prog string) {
	t.Helper()
	dir := filepath.Dir(prog)
	if err := os.Rename(prog, filepath.Join(dir, "player")); err != nil {
		t.Fatal(err)
	}
	wrapper := "#!/bin/sh\nexport ORBITAL_COLONY_UNIT=2ms\nexport GORACE=atexit_sleep_ms=0\nexec ./player\n"
	if err := os.WriteFile(prog, []byte(wrapper), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en"}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.stop)
}

func playArchiveMarketContract(t *testing.T, r *runner, contract *econ.Contract, project colonyProject) contractResult {
	t.Helper()
	if missing := projectRequirements(project.Files, contract.Requires, "en"); len(missing) > 0 {
		return contractResult{Requires: missing}
	}
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	prepareArchiveContract(t, c, contract.ID)
	if err := c.state.Accept(contract.ID); err != nil {
		t.Fatal(err)
	}
	start := c.state.Tick
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	prog := buildColonyProjectProgram(t, r, project)
	if strings.HasPrefix(contract.ID, "h") {
		requireMarketSocket(t)
	}
	startFastSystemProgram(t, c, prog)
	defer c.stop()
	deadline := time.NewTimer(time.Duration(contract.Window)*contractUnit + 10*time.Second)
	defer deadline.Stop()
	var candidate *contractResult
	for {
		select {
		case line := <-ch:
			var ev struct {
				Op       string
				OK       bool
				Msg      string
				Code     int
				State    *econ.State
				Contract struct{ ID string }
			}
			if json.Unmarshal([]byte(line), &ev) != nil {
				continue
			}
			if ev.Op == "contract" && ev.Contract.ID == contract.ID {
				ticks := int64(1)
				if ev.State != nil {
					ticks = max(1, ev.State.Tick-start)
				}
				candidate = &contractResult{OK: ev.OK, Ticks: ticks, Msg: ev.Msg}
			}
			if ev.Op == "stopped" {
				if candidate == nil {
					t.Fatalf("%s exited without contract settlement: %s", contract.ID, ev.Msg)
				}
				if strings.HasPrefix(contract.ID, "h") {
					c.mu.Lock()
					stock := c.state.Inv["moss"]
					c.mu.Unlock()
					if ev.Code != 0 || stock != 980 {
						t.Fatalf("client lost an order or sim duplicated a sale: exit=%d stock=%v", ev.Code, stock)
					}
				}
				return *candidate
			}
		case <-deadline.C:
			t.Fatalf("%s timed out", contract.ID)
		}
	}
}
