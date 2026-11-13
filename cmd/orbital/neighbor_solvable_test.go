package main

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/labkit"
)

func playNeighborContract(t *testing.T, r *runner, contract *econ.Contract, project colonyProject) contractResult {
	t.Helper()
	if missing := projectRequirements(project.Files, contract.Requires, "en"); len(missing) > 0 {
		return contractResult{Requires: missing}
	}
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state = econ.New(31)
	c.state.Levels["sys.neighbor"] = 1
	c.state.Levels["yield.moss"] = 4
	c.state.Levels["api.try"] = 1
	if contract.ID == "n31-bad-line" {
		c.state.Contracts["n31-trade-route"] = &econ.ContractRun{Best: 1, Status: "done"}
	}
	if err := c.state.Accept(contract.ID); err != nil {
		t.Fatal(err)
	}
	dir, err := r.colonyProjectDir(project.Files)
	if err != nil {
		t.Fatal(err)
	}
	if failure := r.buildColonyProject(withLang(t.Context(), "en"), dir, &project); failure != nil {
		t.Fatalf("build: %+v", failure)
	}
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if errors.Is(err, syscall.EPERM) {
		os.RemoveAll(dir)
		t.Skip("sandbox prohibits loopback listeners")
	}
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	probe.Close()
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en", env: []string{"ORBITAL_COLONY_UNIT=2ms", "GORACE=atexit_sleep_ms=0"}}); err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	defer c.stop()
	deadline := time.NewTimer(30 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case <-deadline.C:
			t.Fatal("neighbour reference did not finish")
		case line := <-ch:
			var event struct {
				Op  string
				OK  bool
				T   int64
				Msg string
			}
			if json.Unmarshal([]byte(line), &event) == nil && event.Op == "contract" {
				return contractResult{OK: event.OK, Ticks: event.T, Msg: event.Msg}
			}
		}
	}
}

func TestNeighborOwnedSidecarRunLifecycle(t *testing.T) {
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if errors.Is(err, syscall.EPERM) {
		t.Skip("sandbox prohibits loopback listeners")
	}
	if err != nil {
		t.Fatal(err)
	}
	probe.Close()
	labkit.LeakCheck(t)
	t.Setenv("ORBITAL_COLONY_UNIT", "2ms")
	r := testRunner(t)
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Levels["sys.neighbor"] = 1
	dir, err := r.colonyProjectDir(map[string]string{"main.go": `package main
import ("os"; "fmt"; "orbital/colony")
func main(){fmt.Println(os.Getenv("ORBITAL_NEIGHBOR_ADDR"));for{colony.Wait()}}`})
	if err != nil {
		t.Fatal(err)
	}
	if failure := r.buildColonyProject(withLang(t.Context(), "en"), dir, &colonyProject{}); failure != nil {
		os.RemoveAll(dir)
		t.Fatal(failure)
	}
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en", env: []string{"ORBITAL_COLONY_UNIT=2ms", "GORACE=atexit_sleep_ms=0"}}); err != nil {
		t.Fatal(err)
	}
	c.mu.Lock()
	n := c.neighbor
	active := c.active
	c.mu.Unlock()
	addr := n.snapshot().Addr
	deadline := time.NewTimer(3 * time.Second)
	received := false
	for !received {
		select {
		case <-deadline.C:
			c.stop()
			t.Fatal("player did not receive neighbour address")
		case line := <-ch:
			var event struct{ Op, Msg string }
			if json.Unmarshal([]byte(line), &event) == nil && event.Op == "stdout" && event.Msg == addr {
				received = true
			}
		}
	}
	deadline.Stop()

	c.stop()
	select {
	case <-active.done:
	default:
		t.Fatal("run still active")
	}
	if conn, err := net.DialTimeout("tcp", addr, 20*time.Millisecond); err == nil {
		conn.Close()
		t.Fatal("listener survived stop")
	}
}
