package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/labkit"
)

func TestSystemsSurvivePlayerRestarts(t *testing.T) {
	requireControlSockets(t)
	labkit.LeakCheck(t)
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state = econ.New(31)
	c.state.Inv["moss"] = 1000
	for _, name := range []string{"archive", "market", "neighbor", "control", "warehouse", "metrics"} {
		c.state.Levels["sys."+name] = 1
	}
	r := testRunner(t)
	files := readColonyFixture(t, "s29-warehouse")
	dir, err := r.colonyProjectDir(files)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	if failure := r.buildColonyProject(context.Background(), dir, &colonyProject{Files: files, Race: true}); failure != nil {
		t.Fatal(failure)
	}
	events, unsubscribe := c.subscribe()
	defer unsubscribe()
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en", env: []string{"ORBITAL_COLONY_UNIT=2ms", "GORACE=atexit_sleep_ms=0"}}); err != nil {
		t.Fatal(err)
	}
	defer c.stop()
	c.mu.Lock()
	control, neighbor, market := c.control, c.neighbor, c.active.am.market
	c.mu.Unlock()
	client := &http.Client{Timeout: 3 * time.Second}
	defer client.CloseIdleConnections()
	for round := range 4 {
		if round > 0 && !c.restartWarehouse(control) {
			t.Fatal("restart failed")
		}
		if !c.awaitControlReady(control) {
			t.Fatal("restarted player never listened")
		}
		key := fmt.Sprintf("restart-%d", round)
		for attempt := 0; ; attempt++ {
			request, err := http.NewRequest(http.MethodPost, market.url+"/orders", strings.NewReader(`{"good":"moss","amount":1}`))
			if err != nil {
				t.Fatal(err)
			}
			request.Header.Set("Idempotency-Key", key)
			response, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			io.Copy(io.Discard, response.Body)
			response.Body.Close()
			if response.StatusCode == 200 {
				break
			}
			if response.StatusCode != 503 || attempt >= 4 {
				t.Fatalf("market after restart %d: HTTP %d", round, response.StatusCode)
			}
		}
		conn, err := net.DialTimeout("tcp", neighbor.snapshot().Addr, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		conn.SetDeadline(time.Now().Add(time.Second))
		reader := bufio.NewReader(conn)
		offer, err := reader.ReadString('\n')
		fields := strings.Fields(offer)
		if err != nil || len(fields) != 5 || fields[0] != "OFFER" {
			conn.Close()
			t.Fatalf("offer: %q %v", offer, err)
		}
		fmt.Fprintf(conn, "ACCEPT %s\nSHIP %s\n", fields[1], fields[1])
		paid, err := reader.ReadString('\n')
		conn.Close()
		if err != nil || !strings.HasPrefix(paid, "PAID "+fields[1]+" ") {
			t.Fatalf("neighbour after restart %d: %q %v", round, paid, err)
		}
	}
	c.mu.Lock()
	systems := c.summary("en", nil)["systems"].(map[string]any)
	c.mu.Unlock()
	for _, name := range []string{"archive", "market", "neighbor", "control", "warehouse", "metrics"} {
		if systems[name] == nil {
			t.Fatalf("summary lost %s: %+v", name, systems)
		}
	}
	if market.orders.Load() != 4 || neighbor.snapshot().Trades != 4 {
		t.Fatal("sidecar progress lost across restart")
	}
	c.stop()
	seen := map[string]bool{}
	for len(events) > 0 {
		var event struct {
			Op      string
			Systems map[string]json.RawMessage
		}
		if json.Unmarshal([]byte(<-events), &event) == nil && event.Op == "systems" {
			if len(event.Systems) == 6 {
				t.Fatal("expected live system patches")
			}
			for key := range event.Systems {
				seen[key] = true
			}
		}
	}
	for key := range systems {
		if !seen[key] {
			t.Fatalf("no live update for %s", key)
		}
	}
	for _, address := range []string{control.addr, neighbor.snapshot().Addr, strings.TrimPrefix(market.url, "http://")} {
		listener, err := net.Listen("tcp", address)
		if err != nil {
			t.Fatalf("listener survived stop: %s: %v", address, err)
		}
		listener.Close()
	}
}
