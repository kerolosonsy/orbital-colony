package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"maps"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/labkit"
)

type controlTransport func(*http.Request) (*http.Response, error)

func (transport controlTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func requireControlSockets(t *testing.T) {
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

func TestControlScheduleSeededAndCoversEndpoints(t *testing.T) {
	var first, replay, other []controlProbe
	for index := range 30 {
		first = append(first, controlProbeAt(11, index))
		replay = append(replay, controlProbeAt(11, index))
		other = append(other, controlProbeAt(12, index))
	}
	if !reflect.DeepEqual(first, replay) || reflect.DeepEqual(first, other) {
		t.Fatal("probe schedule is not seeded")
	}
	seen := map[string]int{}
	for _, probe := range first {
		if probe.delay < 250*time.Millisecond || probe.delay > 500*time.Millisecond {
			t.Fatal(probe)
		}
		seen[probe.path]++
	}
	if len(seen) != 3 || seen["/orders"] != 10 {
		t.Fatal(seen)
	}
}

func TestControlJudgesStatusInventoryDeliveryAndLatency(t *testing.T) {
	for _, tc := range []struct {
		name, path, body string
		status           int
		delay            time.Duration
		stock, delivered float64
		age              time.Duration
		ok               bool
	}{
		{name: "health", path: "/healthz", status: 200, ok: true},
		{name: "unhealthy", path: "/healthz", status: 503},
		{name: "late", path: "/healthz", status: 200, delay: 210 * time.Millisecond},
		{name: "inventory", path: "/inventory", status: 200, body: `{"moss":12.5}`, stock: 12.5, ok: true},
		{name: "wrong inventory", path: "/inventory", status: 200, body: `{"moss":12}`, stock: 12.5},
		{name: "stale inventory", path: "/inventory", status: 200, body: `{"moss":12.5}`, stock: 12.5, age: 1100 * time.Millisecond},
		{name: "null inventory", path: "/inventory", status: 200, body: `null`},
		{name: "extra goods", path: "/inventory", status: 200, body: `{"moss":12.5,"iron":1}`, stock: 12.5},
		{name: "fulfilled order", path: "/orders", status: 202, stock: 10, delivered: 5, ok: true},
		{name: "accepted without export", path: "/orders", status: 202, stock: 10},
		{name: "insufficient stock", path: "/orders", status: 409, stock: 3, ok: true},
		{name: "false refusal", path: "/orders", status: 409, stock: 10},
		{name: "false acceptance", path: "/orders", status: 202, stock: 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st := econ.New(11)
			st.Inv["moss"] = tc.stock
			c := &colonyGame{state: st}
			control := &colonyControl{ctx: context.Background(), addr: "127.0.0.1:12345", history: []controlSample{{at: time.Now().Add(-tc.age), inv: maps.Clone(st.Inv)}}}
			control.client = &http.Client{Transport: controlTransport(func(request *http.Request) (*http.Response, error) {
				if request.URL.Path != tc.path {
					t.Errorf("unexpected path %s", request.URL.Path)
				}
				time.Sleep(tc.delay)
				st.Systems.SoldByGood = map[string]float64{"moss": tc.delivered}
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}
			ok, _ := c.judgeControlProbe(control, controlProbe{path: tc.path, amount: 5})
			if ok != tc.ok {
				t.Fatalf("judged %v want %v", ok, tc.ok)
			}
		})
	}
}

func TestControlThirdPartyImportsRequireOwnedNodes(t *testing.T) {
	s := testColony(t)
	for _, tc := range []struct{ node, pkg string }{{"sys.warehouse", "modernc.org/sqlite"}, {"sys.metrics", "github.com/prometheus/client_golang/prometheus"}} {
		for _, owned := range []bool{false, true} {
			c := s.colony
			c.mu.Lock()
			c.state.Levels[tc.node] = 0
			if owned {
				c.state.Levels[tc.node] = 1
			}
			c.mu.Unlock()
			project := &colonyProject{Files: map[string]string{"main.go": "package main", "nested/import.go": "package nested\nimport _ " + strconv.Quote(tc.pkg)}}
			req := httptest.NewRequest("POST", "/", nil)
			response := httptest.NewRecorder()
			if s.checkControlImports(response, req, project) != owned {
				t.Fatalf("%s owned=%v: %s", tc.pkg, owned, response.Body)
			}
		}
	}
}

func TestControlStopCancelsWorkersAndReleasesAddress(t *testing.T) {
	requireControlSockets(t)
	labkit.LeakCheck(t)
	r := testRunner(t)
	files := readColonyFixture(t, "h26-control-api")
	dir, err := r.colonyProjectDir(files)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	if failure := r.buildColonyProject(context.Background(), dir, &colonyProject{Files: files}); failure != nil {
		t.Fatal(failure)
	}
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Levels["sys.control"] = 1
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en"}); err != nil {
		t.Fatal(err)
	}
	c.mu.Lock()
	control := c.control
	c.mu.Unlock()
	if !c.awaitControlReady(control) {
		c.stop()
		t.Fatal("player never listened")
	}
	c.stop()
	select {
	case <-control.done:
	default:
		t.Fatal("probe workers survived stop")
	}
	listener, err := net.Listen("tcp", control.addr)
	if err != nil {
		t.Fatalf("player listener survived stop: %v", err)
	}
	listener.Close()
	c.mu.Lock()
	on := c.control.running
	c.mu.Unlock()
	if on {
		t.Fatal("stopped system still on")
	}
}

func playControlContract(t *testing.T, r *runner, contract *econ.Contract, project colonyProject) contractResult {
	t.Helper()
	requireControlSockets(t)
	if missing := projectRequirements(project.Files, contract.Requires, "en"); len(missing) > 0 {
		return contractResult{Requires: missing}
	}
	dir, err := r.colonyProjectDir(project.Files)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	if failure := r.buildColonyProject(context.Background(), dir, &project); failure != nil {
		t.Fatal(failure)
	}
	if err := os.Rename(filepath.Join(dir, "prog"), filepath.Join(dir, "player")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "prog"), []byte("#!/bin/sh\nexec env ORBITAL_COLONY_UNIT=2ms GORACE=atexit_sleep_ms=0 ./player\n"), 0700); err != nil {
		t.Fatal(err)
	}
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Inv["moss"], c.state.Inv["lichen"], c.state.Credits = 100, 20, 500
	c.state.Levels["sys.control"], c.state.Levels["engine"] = 1, 8
	if contract.Chapter == 29 {
		c.state.Levels["sys.warehouse"] = 1
	}
	if contract.Chapter == 36 {
		c.state.Levels["sys.metrics"] = 1
	}
	if err := c.state.Accept(contract.ID); err != nil {
		t.Fatal(err)
	}
	events, unsubscribe := c.subscribe()
	defer unsubscribe()
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en"}); err != nil {
		t.Fatal(err)
	}
	defer c.stop()
	c.mu.Lock()
	details := c.active.details
	c.mu.Unlock()
	deadline := time.NewTimer(25 * time.Second)
	defer deadline.Stop()
	var diagnostics []string
	for {
		select {
		case <-deadline.C:
			t.Fatalf("%s timed out: %v", contract.ID, diagnostics)
		case line := <-events:
			var event struct {
				Op       string
				OK       bool
				Msg      string
				Contract struct{ ID string }
				State    *econ.State
			}
			if json.Unmarshal([]byte(line), &event) != nil {
				continue
			}
			if event.Op == "stderr" {
				diagnostics = append(diagnostics, event.Msg)
			}
			if event.Op == "contract" && event.Contract.ID == contract.ID {
				c.stop()
				c.mu.Lock()
				race := details.record.Race
				restarts := c.control.restarts
				c.mu.Unlock()
				if event.OK && contract.ID == "s29-warehouse" && restarts != 3 {
					t.Fatalf("warehouse restarted %d times", restarts)
				}
				return contractResult{OK: event.OK && !race, Ticks: event.State.Tick, Msg: strings.Join(diagnostics, "\n")}
			}
			if event.Op == "stopped" {
				return contractResult{Msg: strings.Join(diagnostics, "\n")}
			}
		}
	}
}

func TestControlDrainJudgesCompletedSlowResponse(t *testing.T) {
	for _, tc := range []struct {
		reply controlDrainReply
		ok    bool
	}{
		{controlDrainReply{status: 202, elapsed: 600 * time.Millisecond}, true},
		{controlDrainReply{status: 202, elapsed: 100 * time.Millisecond}, false},
		{controlDrainReply{status: 202, elapsed: 3 * time.Second}, false},
		{controlDrainReply{status: 409, elapsed: 600 * time.Millisecond}, false},
		{controlDrainReply{status: 202, elapsed: 600 * time.Millisecond, err: io.ErrUnexpectedEOF}, false},
	} {
		if tc.reply.correct(202) != tc.ok {
			t.Fatalf("reply %+v judged incorrectly", tc.reply)
		}
	}
}

func TestControlReadinessBlocksOrders(t *testing.T) {
	for _, status := range []int{200, 503} {
		st := econ.New(11)
		c := &colonyGame{state: st}
		orders := 0
		control := &colonyControl{ctx: context.Background(), metrics: true, addr: "127.0.0.1:1"}
		control.client = &http.Client{Transport: controlTransport(func(request *http.Request) (*http.Response, error) {
			code := status
			if request.URL.Path == "/orders" {
				orders++
				code = 409
			}
			return &http.Response{StatusCode: code, Body: io.NopCloser(strings.NewReader("")), Header: make(http.Header)}, nil
		})}
		ok, _ := c.probeOrder(control, 5)
		if ok != (status == 200) || (orders > 0) != (status == 200) {
			t.Fatalf("ready=%d ok=%v orders=%d", status, ok, orders)
		}
	}
}

func TestControlReferencesCompileOffline(t *testing.T) {
	r := testRunner(t)
	for _, id := range []string{"h26-control-api", "h28-drain", "s29-warehouse", "s29-transfer", "o36-dashboard"} {
		t.Run(id, func(t *testing.T) {
			files := readColonyFixture(t, id)
			if missing := projectRequirements(files, econ.FindContract(id).Requires, "en"); len(missing) > 0 {
				t.Fatal(missing)
			}
			passed, output, err := r.goTest(context.Background(), colonyControlLevel(), files, "-run", "^$", "./...")
			if err != nil || !passed {
				t.Fatalf("%v\n%s", err, output)
			}
		})
	}
}

func TestControlContractsRejectBrokenHTTPServices(t *testing.T) {
	requireControlSockets(t)
	for _, tc := range []struct{ id, old, replacement string }{
		{"h28-drain", "time.Sleep(time.Duration(delay) * time.Millisecond)", "time.Sleep(0)"},
		{"s29-warehouse", `sql.Open("sqlite", os.Getenv("ORBITAL_DB_PATH"))`, `sql.Open("sqlite", ":memory:"+os.Getenv("ORBITAL_DB_PATH")[:0])`},
		{"s29-transfer", "if err == nil {", "if false {"},
		{"o36-dashboard", "return colony.SystemStats().Credits", "return colony.SystemStats().Credits+100"},
		{"o36-dashboard", "ready.Store(false)", "ready.Store(true)"},
	} {
		t.Run(tc.id+tc.old, func(t *testing.T) {
			files := readColonyFixture(t, tc.id)
			changed := false
			for name, code := range files {
				if strings.Contains(code, tc.old) {
					files[name] = strings.ReplaceAll(code, tc.old, tc.replacement)
					changed = true
				}
			}
			if !changed {
				t.Fatal("negative fixture mutation did not apply")
			}
			played := playControlContract(t, testRunner(t), econ.FindContract(tc.id), colonyProject{Files: files, Race: true})
			if played.OK {
				t.Fatal("broken service passed", played)
			}
		})
	}
}

func TestWarehouseRestartWorkersStopWithoutLeaks(t *testing.T) {
	requireControlSockets(t)
	labkit.LeakCheck(t)
	played := playControlContract(t, testRunner(t), econ.FindContract("s29-warehouse"), colonyProject{Files: readColonyFixture(t, "s29-warehouse"), Race: true})
	if !played.OK {
		t.Fatal(played)
	}
}

func TestControlEndpointsRejectUnownedDependenciesBeforeBuild(t *testing.T) {
	s := testColony(t)
	s.run.goBin = "/no/such/go"
	for _, endpoint := range []http.HandlerFunc{s.colonyRun, s.colonyTest, s.colonyVet} {
		for _, pkg := range []string{"modernc.org/sqlite", "github.com/prometheus/client_golang/prometheus"} {
			response := projectPost(t, endpoint, colonyProject{Files: map[string]string{"main.go": "package main\nimport _ " + strconv.Quote(pkg) + "\nfunc main(){}"}})
			var reply struct{ Phase string }
			if json.Unmarshal(response.Body.Bytes(), &reply) != nil || reply.Phase != "require" {
				t.Fatalf("dependency reached build: HTTP %d %s", response.Code, response.Body)
			}
		}
	}
}

func TestWarehouseManualStopDuringRestart(t *testing.T) {
	requireControlSockets(t)
	labkit.LeakCheck(t)
	r := testRunner(t)
	files := readColonyFixture(t, "s29-warehouse")
	dir, err := r.colonyProjectDir(files)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	if failure := r.buildColonyProject(context.Background(), dir, &colonyProject{Files: files}); failure != nil {
		t.Fatal(failure)
	}
	c, err := newColonyGame(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	c.state.Levels["sys.control"], c.state.Levels["sys.warehouse"] = 1, 1
	if err := c.state.Accept("s29-warehouse"); err != nil {
		t.Fatal(err)
	}
	if err := c.start(dir, "en", &colonyRunDetails{lang: "en"}); err != nil {
		t.Fatal(err)
	}
	defer c.stop()
	deadline := time.Now().Add(5 * time.Second)
	for {
		c.mu.Lock()
		restarting := c.control.restart
		control := c.control
		c.mu.Unlock()
		if restarting {
			c.stop()
			c.mu.Lock()
			running, restarts, active := control.running, control.restarts, c.active
			c.mu.Unlock()
			if running || active != nil || restarts != 0 {
				t.Fatalf("manual stop escaped: running=%v restarts=%d", running, restarts)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("warehouse never requested restart")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestControlExitWaitsForFinalDrainVerdict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	control := &colonyControl{ctx: ctx, cancel: cancel, done: make(chan struct{}), verdictReady: make(chan struct{}), draining: true, contract: "h28-drain", client: &http.Client{}}
	c := &colonyGame{state: econ.New(11), control: control}
	go func() { <-ctx.Done(); close(control.done) }()
	finished := make(chan struct{})
	go func() { c.stopControlProbes(control); close(finished) }()
	select {
	case <-ctx.Done():
		t.Fatal("exit cancelled final response before verification")
	case <-time.After(20 * time.Millisecond):
	}
	c.controlVerdict(control, true)
	select {
	case <-finished:
	case <-time.After(time.Second):
		t.Fatal("workers were not joined after verdict")
	}
}
