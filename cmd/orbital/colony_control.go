package main

import (
	"context"
	"encoding/json"
	"io"
	"maps"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"orbital/colony/econ"
)

const controlLatency = 200 * time.Millisecond

type controlSample struct {
	at  time.Time
	inv map[string]float64
}

type colonyControl struct {
	probeMu                    sync.Mutex
	ctx                        context.Context
	cancel                     context.CancelFunc
	done                       chan struct{}
	client                     *http.Client
	addr                       string
	contract                   string
	seed                       int64
	warehouse, metrics         bool
	healthy                    bool
	lastMS                     float64
	probes, failed, restarts   int
	history                    []controlSample
	series                     map[string][][2]float64
	labels                     map[string]map[string]bool
	stopping, restart, running bool
	restarted                  chan struct{}
	draining                   bool
	verdictReady               chan struct{}
	verdict                    *bool
}

func controlEnv(env []string) []string {
	clean := make([]string, 0, len(env))
	for _, entry := range env {
		if !strings.HasPrefix(entry, "ORBITAL_LISTEN_ADDR=") && !strings.HasPrefix(entry, "ORBITAL_DB_PATH=") {
			clean = append(clean, entry)
		}
	}
	return clean
}

func (c *colonyGame) startControlProcess(cmd *exec.Cmd) error {
	cmd.Env = controlEnv(cmd.Env)
	c.mu.Lock()
	c.control = nil
	owned := c.state.Has("sys.control")
	contract, seed := c.state.Active, c.state.Seed
	warehouse, metrics := c.state.Has("sys.warehouse"), c.state.Has("sys.metrics")
	c.mu.Unlock()
	if !owned {
		return cmd.Start()
	}
	reservation, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	transport := &http.Transport{Proxy: nil, DisableKeepAlives: true, MaxConnsPerHost: 16, ResponseHeaderTimeout: 2 * time.Second}
	control := &colonyControl{ctx: ctx, cancel: cancel, done: make(chan struct{}), verdictReady: make(chan struct{}), addr: reservation.Addr().String(), contract: contract, seed: seed, warehouse: warehouse, metrics: metrics, running: true, restarted: make(chan struct{}, 1), series: map[string][][2]float64{}, labels: map[string]map[string]bool{}, client: &http.Client{Transport: transport, Timeout: 2 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
	cmd.Env = append(cmd.Env, "ORBITAL_LISTEN_ADDR="+control.addr)
	if warehouse {
		cmd.Env = append(cmd.Env, "ORBITAL_DB_PATH="+filepath.Join(filepath.Dir(c.path), "warehouse.sqlite"))
	}
	// net/http players bind by address; release the reservation immediately before exec.
	if err := reservation.Close(); err != nil {
		cancel()
		return err
	}
	if err := cmd.Start(); err != nil {
		cancel()
		return err
	}
	c.mu.Lock()
	c.control = control
	c.mu.Unlock()
	return nil
}

func (c *colonyGame) startControlProbes() {
	c.mu.Lock()
	control := c.control
	c.mu.Unlock()
	if control != nil {
		go func() { defer close(control.done); c.runControl(control) }()
	}
}

// Caller holds mu, so a concurrent restart cannot escape a manual stop.
func (c *colonyGame) markControlStop() *colonyRun {
	if c.control != nil {
		c.control.stopping = true
	}
	return c.active
}

func (c *colonyGame) waitControlRun(active *colonyRun, dir string, readers *sync.WaitGroup) (error, int) {
	c.mu.Lock()
	control := c.control
	c.mu.Unlock()
	current := active
	for {
		waitErr := current.cmd.Wait()
		_ = killProcessTree(current.cmd)
		readers.Wait()
		c.mu.Lock()
		if control == nil || !control.restart || control.stopping {
			c.mu.Unlock()
			if control != nil {
				c.stopControlProbes(control)
			}
			return waitErr, current.cmd.ProcessState.ExitCode()
		}
		next, err := c.restartControlProcess(current, dir, readers)
		if err != nil {
			control.restart = false
			c.mu.Unlock()
			c.stopControlProbes(control)
			return err, -1
		}
		current, c.active = next, next
		control.restart = false
		control.restarts++
		control.restarted <- struct{}{}
		c.mu.Unlock()
	}
}

// Caller holds mu until the replacement child and its streams are registered.
func (c *colonyGame) restartControlProcess(previous *colonyRun, dir string, readers *sync.WaitGroup) (*colonyRun, error) {
	body, err := json.Marshal(c.state)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "colony.json"), body, 0o600); err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(context.Background(), previous.cmd.Path)
	configureProcess(cmd)
	cmd.Dir, cmd.Env = dir, previous.cmd.Env
	pipes := make([]*os.File, 0, 6)
	for range 3 {
		r, w, err := os.Pipe()
		if err != nil {
			for _, file := range pipes {
				file.Close()
			}
			return nil, err
		}
		pipes = append(pipes, r, w)
	}
	cmd.ExtraFiles = []*os.File{pipes[1], previous.commandsReader}
	cmd.Stdout, cmd.Stderr = pipes[3], pipes[5]
	if err := cmd.Start(); err != nil {
		for _, file := range pipes {
			file.Close()
		}
		return nil, err
	}
	for _, i := range []int{1, 3, 5} {
		pipes[i].Close()
	}
	for _, stream := range []struct {
		file *os.File
		op   string
	}{{pipes[0], ""}, {pipes[2], "stdout"}, {pipes[4], "stderr"}} {
		readers.Add(1)
		go func(file *os.File, op string) { defer readers.Done(); defer file.Close(); c.readLines(file, op) }(stream.file, stream.op)
	}
	return &colonyRun{am: previous.am, cmd: cmd, commands: previous.commands, commandsReader: previous.commandsReader, done: previous.done, details: previous.details}, nil
}

func (c *colonyGame) observeControlEvent(line string) {
	var event struct {
		State *econ.State
	}
	if json.Unmarshal([]byte(line), &event) != nil || event.State == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.control == nil {
		return
	}
	control := c.control
	control.history = append(control.history, controlSample{at: time.Now(), inv: maps.Clone(event.State.Inv)})
	if len(control.history) > 128 {
		control.history = control.history[len(control.history)-128:]
	}
}

// Caller holds mu; returned containers are detached from the live series.
func (c *colonyGame) controlSummary() map[string]any {
	systems := map[string]any{}
	control := c.control
	if c.state.Has("sys.control") {
		systems["control"] = map[string]any{"on": false}
	}
	if c.state.Has("sys.warehouse") {
		systems["warehouse"] = map[string]any{"on": false, "restarts": 0}
	}
	if c.state.Has("sys.metrics") {
		systems["metrics"] = map[string]any{"on": false, "series": map[string]any{}}
	}
	if control == nil {
		return systems
	}
	systems["control"] = map[string]any{"on": control.running, "addr": control.addr, "healthy": control.healthy, "lastMs": control.lastMS, "probes": control.probes, "failed": control.failed}
	if control.warehouse {
		systems["warehouse"] = map[string]any{"on": control.running, "restarts": control.restarts}
	}
	if control.metrics {
		series := make(map[string][][2]float64, len(control.series))
		for name, points := range control.series {
			series[name] = append([][2]float64{}, points...)
		}
		systems["metrics"] = map[string]any{"on": control.running, "series": series}
	}
	return systems
}

func (c *colonyGame) broadcastControl() {
	c.mu.Lock()
	line, _ := json.Marshal(map[string]any{"op": "systems", "systems": c.controlSummary()})
	c.mu.Unlock()
	c.broadcast(string(line))
}

func (c *colonyGame) controlCommand(command any) error {
	line, err := json.Marshal(command)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active == nil {
		return io.ErrClosedPipe
	}
	if err := c.active.commands.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return err
	}
	_, err = c.active.commands.Write(append(line, '\n'))
	return err
}

func (c *colonyGame) finishControlContract(lang string) string {
	control := c.control
	if control == nil || control.verdict == nil || c.state.Active != control.contract {
		return ""
	}
	id := control.contract
	c.state.VerifyControl(id, *control.verdict)
	ok := *control.verdict
	if ok {
		_, _ = c.state.Settle(c.state.Tick)
	} else {
		c.state.Fail()
	}
	line, _ := json.Marshal(map[string]any{"op": "contract", "ok": ok, "msg": tr(lang, "فحص النظام خلص", "System verification finished"), "state": c.state, "contract": map[string]any{"id": id, "frac": c.state.Fraction(), "left": -1}})
	return string(line)
}

func (c *colonyGame) stopControlProbes(control *colonyControl) {
	c.mu.Lock()
	pending := control.draining && control.verdict == nil
	c.mu.Unlock()
	// A child can exit as its final response reaches the HTTP client.
	if pending && control.verdictReady != nil {
		timer := time.NewTimer(3 * time.Second)
		select {
		case <-control.verdictReady:
		case <-timer.C:
		}
		timer.Stop()
	}
	control.cancel()
	<-control.done
	control.client.CloseIdleConnections()
	c.mu.Lock()
	control.running, control.healthy = false, false
	c.mu.Unlock()
	c.broadcastControl()
}

func (c *colonyGame) broadcastControlEvent(line string) {
	if line != "" {
		c.broadcast(line)
	}
}
