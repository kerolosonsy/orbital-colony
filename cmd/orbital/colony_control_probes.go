package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httptrace"
	"sync"
	"syscall"
	"time"

	"orbital/colony/econ"
)

type controlProbe struct {
	path   string
	delay  time.Duration
	amount int
}

func controlProbeAt(seed int64, index int) controlProbe {
	hash := econ.EventHash(seed, int64(index), 26)
	path := []string{"/healthz", "/inventory", "/orders"}[index%3]
	amount := 5
	if index%6 == 5 {
		amount = 1000000000
	}
	return controlProbe{path: path, delay: time.Duration(250+hash%251) * time.Millisecond, amount: amount}
}

func controlWait(ctx context.Context, duration time.Duration) bool {
	timer := time.NewTimer(max(0, duration))
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (control *colonyControl) request(ctx context.Context, path string, payload any) (int, []byte, time.Duration, error) {
	var body io.Reader
	method := http.MethodGet
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return 0, nil, 0, err
		}
		body, method = bytes.NewReader(encoded), http.MethodPost
	}
	req, err := http.NewRequestWithContext(ctx, method, "http://"+control.addr+path, body)
	if err != nil {
		return 0, nil, 0, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	started := time.Now()
	response, err := control.client.Do(req)
	if err != nil {
		return 0, nil, time.Since(started), err
	}
	defer response.Body.Close()
	content, err := io.ReadAll(io.LimitReader(response.Body, 1<<20+1))
	if len(content) > 1<<20 {
		err = io.ErrShortBuffer
	}
	return response.StatusCode, content, time.Since(started), err
}

func (c *colonyGame) runControl(control *colonyControl) {
	var workers sync.WaitGroup
	workers.Add(1)
	go func() { defer workers.Done(); c.controlSnapshots(control) }()
	if control.metrics {
		workers.Add(1)
		go func() { defer workers.Done(); c.runMetrics(control) }()
	}
	defer workers.Wait()
	c.broadcastControl()
	if !c.awaitControlReady(control) {
		c.controlVerdict(control, false)
		<-control.ctx.Done()
		return
	}
	switch control.contract {
	case "s29-warehouse":
		c.probeWarehouse(control)
	case "s29-transfer":
		c.controlVerdict(control, c.probeTransfers(control))
	default:
		c.probeControlLoop(control)
	}
	<-control.ctx.Done()
}

func (c *colonyGame) controlSnapshots(control *colonyControl) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-control.ctx.Done():
			return
		case <-ticker.C:
			_ = c.controlCommand(map[string]bool{"controlSnapshot": true})
		}
	}
}

func (c *colonyGame) awaitControlReady(control *colonyControl) bool {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		status, _, _, err := control.request(control.ctx, "/healthz", nil)
		if err == nil && status == http.StatusOK {
			return true
		}
		if !controlWait(control.ctx, 25*time.Millisecond) {
			return false
		}
	}
	return false
}

func (c *colonyGame) probeControlLoop(control *colonyControl) {
	at := time.Now()
	for index := 0; ; index++ {
		probe := controlProbeAt(control.seed, index)
		at = at.Add(probe.delay)
		if !controlWait(control.ctx, time.Until(at)) {
			return
		}
		control.probeMu.Lock()
		c.mu.Lock()
		draining := control.draining
		c.mu.Unlock()
		if draining {
			control.probeMu.Unlock()
			return
		}
		ok, latency := c.judgeControlProbe(control, probe)
		c.recordControlProbe(control, ok, latency)
		if !ok {
			c.controlVerdict(control, false)
		}
		control.probeMu.Unlock()
		if index == 11 && control.contract == "h26-control-api" {
			c.controlVerdict(control, true)
		}
		if index == 5 && control.contract == "h28-drain" {
			c.controlVerdict(control, c.probeDrain(control))
			return
		}
	}
}

func (c *colonyGame) judgeControlProbe(control *colonyControl, probe controlProbe) (bool, time.Duration) {
	if probe.path == "/orders" {
		return c.probeOrder(control, probe.amount)
	}
	started := time.Now()
	status, body, elapsed, err := control.request(control.ctx, probe.path, nil)
	ok := err == nil && status == http.StatusOK && elapsed <= controlLatency
	if probe.path == "/inventory" {
		ok = ok && c.controlInventoryMatches(control, body, started)
	}
	return ok, elapsed
}

func inventoryEqual(actual, expected map[string]float64) bool {
	for good, amount := range actual {
		if math.IsNaN(amount) || math.IsInf(amount, 0) || math.Abs(amount-expected[good]) > 1e-6 {
			return false
		}
	}
	for good, amount := range expected {
		if math.Abs(amount-actual[good]) > 1e-6 {
			return false
		}
	}
	return true
}

func decodeInventory(body []byte) (map[string]float64, bool) {
	var inventory map[string]float64
	err := json.Unmarshal(body, &inventory)
	return inventory, err == nil && inventory != nil
}

func (c *colonyGame) controlInventoryMatches(control *colonyControl, body []byte, started time.Time) bool {
	inventory, ok := decodeInventory(body)
	if !ok {
		return false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, sample := range control.history {
		if !sample.at.Before(started.Add(-time.Second)) && inventoryEqual(inventory, sample.inv) {
			return true
		}
	}
	return false
}

func (c *colonyGame) probeOrder(control *colonyControl, amount int) (bool, time.Duration) {
	if control.metrics {
		status, _, elapsed, err := control.request(control.ctx, "/readyz", nil)
		if err != nil || status != http.StatusOK || elapsed > controlLatency {
			return false, elapsed
		}
	}
	c.mu.Lock()
	available, sold := c.state.Inv["moss"], c.state.Systems.SoldByGood["moss"]
	c.mu.Unlock()
	status, _, elapsed, err := control.request(control.ctx, "/orders", map[string]any{"good": "moss", "amount": amount})
	if err != nil || elapsed > controlLatency {
		return false, elapsed
	}
	if available < float64(amount) {
		return status == http.StatusConflict && c.controlSold("moss") == sold, elapsed
	}
	if status != http.StatusAccepted {
		return false, elapsed
	}
	ctx, cancel := context.WithTimeout(control.ctx, time.Second)
	defer cancel()
	return c.awaitControlSale(ctx, sold+float64(amount)), elapsed
}

func (c *colonyGame) awaitControlSale(ctx context.Context, target float64) bool {
	for {
		if c.controlSold("moss") >= target-1e-6 {
			return true
		}
		if !controlWait(ctx, 10*time.Millisecond) {
			return false
		}
	}
}

func (c *colonyGame) controlSold(good string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state.Systems.SoldByGood[good]
}

func (c *colonyGame) recordControlProbe(control *colonyControl, ok bool, elapsed time.Duration) {
	c.mu.Lock()
	control.probes++
	control.lastMS, control.healthy = float64(elapsed)/float64(time.Millisecond), ok
	if !ok {
		control.failed++
	}
	c.mu.Unlock()
	c.broadcastControl()
}

func (c *colonyGame) controlVerdict(control *colonyControl, ok bool) {
	c.mu.Lock()
	if control.verdict != nil || !econ.IsControlContract(control.contract) {
		c.mu.Unlock()
		return
	}
	control.verdict = &ok
	if control.verdictReady != nil {
		close(control.verdictReady)
	}
	c.mu.Unlock()
	_ = c.controlCommand(map[string]any{"controlResult": map[string]any{"id": control.contract, "ok": ok}})
}

func (c *colonyGame) signalControl(control *colonyControl) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.control != control || c.active == nil || control.stopping {
		return false
	}
	c.active.stopSent = true
	return c.active.cmd.Process.Signal(syscall.SIGTERM) == nil
}

func (c *colonyGame) probeDrain(control *colonyControl) bool {
	control.probeMu.Lock()
	defer control.probeMu.Unlock()
	if control.metrics {
		status, _, elapsed, err := control.request(control.ctx, "/readyz", nil)
		if err != nil || status != http.StatusOK || elapsed > controlLatency {
			return false
		}
	}
	c.mu.Lock()
	control.draining = true
	c.mu.Unlock()
	ctx, cancel := context.WithTimeout(control.ctx, 3*time.Second)
	defer cancel()
	wrote := make(chan struct{}, 1)
	trace := &httptrace.ClientTrace{WroteRequest: func(httptrace.WroteRequestInfo) {
		select {
		case wrote <- struct{}{}:
		default:
		}
	}}
	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(ctx, trace), http.MethodPost, "http://"+control.addr+"/orders", bytes.NewBufferString(`{"good":"moss","amount":5}`))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Orbital-Delay-Ms", "500")
	completed := make(chan controlDrainReply, 1)
	c.mu.Lock()
	available, sold := c.state.Inv["moss"], c.state.Systems.SoldByGood["moss"]
	c.mu.Unlock()
	var worker sync.WaitGroup
	worker.Add(1)
	go func() { defer worker.Done(); completed <- control.drainReply(req) }()
	defer func() { cancel(); worker.Wait() }()
	select {
	case <-wrote:
	case <-ctx.Done():
		return false
	}
	if !controlWait(ctx, 100*time.Millisecond) {
		return false
	}
	select {
	case <-completed:
		return false
	default:
	}
	if !c.signalControl(control) {
		return false
	}
	ready := !control.metrics
	refused := false
	for ctx.Err() == nil && (!refused || !ready) {
		status, _, _, err := control.request(ctx, "/readyz", nil)
		ready = ready || err == nil && status == http.StatusServiceUnavailable
		connection, dialErr := (&net.Dialer{Timeout: 50 * time.Millisecond}).DialContext(ctx, "tcp", control.addr)
		if errors.Is(dialErr, syscall.ECONNREFUSED) {
			refused = true
		}
		if connection != nil {
			connection.Close()
		}
		if !controlWait(ctx, 10*time.Millisecond) {
			break
		}
	}
	select {
	case reply := <-completed:
		expected := http.StatusConflict
		if available >= 5 {
			expected = http.StatusAccepted
		}
		return reply.correct(expected) && refused && ready && (available < 5 || c.awaitControlSale(ctx, sold+5))
	case <-ctx.Done():
		return false
	}
}

type controlDrainReply struct {
	status  int
	elapsed time.Duration
	err     error
}

func (control *colonyControl) drainReply(request *http.Request) controlDrainReply {
	started := time.Now()
	response, err := control.client.Do(request)
	if err != nil {
		return controlDrainReply{err: err}
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20+1))
	if len(body) > 1<<20 {
		err = io.ErrShortBuffer
	}
	return controlDrainReply{status: response.StatusCode, elapsed: time.Since(started), err: err}
}

func (reply controlDrainReply) correct(expected int) bool {
	return reply.err == nil && reply.status == expected && reply.elapsed >= 500*time.Millisecond && reply.elapsed < 3*time.Second
}
