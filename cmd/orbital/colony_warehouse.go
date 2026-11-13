package main

import (
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func (c *colonyGame) warehouseInventory(control *colonyControl) (map[string]float64, bool) {
	status, body, _, err := control.request(control.ctx, "/inventory", nil)
	inventory, ok := decodeInventory(body)
	return inventory, ok && err == nil && status == http.StatusOK
}

func (c *colonyGame) probeWarehouse(control *colonyControl) {
	if !control.warehouse {
		c.controlVerdict(control, false)
		return
	}
	expected, ok := c.warehouseInventory(control)
	if !ok {
		c.controlVerdict(control, false)
		return
	}
	c.mu.Lock()
	run := c.state.Contracts[control.contract]
	prefix := fmt.Sprintf("%d-%d-%d", control.seed, run.Start, run.Tries)
	c.mu.Unlock()
	for round := range 3 {
		goods := map[string]float64{"moss": float64(7 + round), "lichen": float64(round + 1)}
		batch := map[string]any{"id": fmt.Sprintf("%s-%d", prefix, round), "goods": goods}
		for range 2 {
			status, _, elapsed, err := control.request(control.ctx, "/inventory", batch)
			if err != nil || status != http.StatusNoContent || elapsed > controlLatency {
				c.controlVerdict(control, false)
				return
			}
		}
		for good, amount := range goods {
			expected[good] += amount
		}
		actual, ok := c.warehouseInventory(control)
		if !ok || !inventoryEqual(actual, expected) || !c.warehouseSQLiteExists() {
			c.controlVerdict(control, false)
			return
		}
		if !c.restartWarehouse(control) || !c.awaitControlReady(control) {
			c.controlVerdict(control, false)
			return
		}
		status, _, elapsed, err := control.request(control.ctx, "/inventory", batch)
		if err != nil || status != http.StatusNoContent || elapsed > controlLatency {
			c.controlVerdict(control, false)
			return
		}
		actual, ok = c.warehouseInventory(control)
		if !ok || !inventoryEqual(actual, expected) {
			c.controlVerdict(control, false)
			return
		}
		c.broadcastControl()
	}
	c.controlVerdict(control, true)
}

func (c *colonyGame) warehouseSQLiteExists() bool {
	file, err := os.Open(filepath.Join(filepath.Dir(c.path), "warehouse.sqlite"))
	if err != nil {
		return false
	}
	defer file.Close()
	header := make([]byte, 16)
	n, err := file.Read(header)
	return err == nil && n == 16 && string(header) == "SQLite format 3\x00"
}

func (c *colonyGame) restartWarehouse(control *colonyControl) bool {
	c.mu.Lock()
	if c.active == nil || control.stopping {
		c.mu.Unlock()
		return false
	}
	control.restart = true
	err := c.active.cmd.Process.Signal(syscall.SIGTERM)
	c.mu.Unlock()
	if err != nil {
		return false
	}
	timer := time.NewTimer(4 * time.Second)
	defer timer.Stop()
	select {
	case <-control.ctx.Done():
		return false
	case <-timer.C:
		c.mu.Lock()
		control.restart = false
		if c.active != nil {
			_ = killProcessTree(c.active.cmd)
		}
		c.mu.Unlock()
		return false
	case <-control.restarted:
		return true
	}
}

func transferBalances(body []byte) (map[string]float64, bool) {
	balances, ok := decodeInventory(body)
	if !ok || len(balances) != 2 {
		return nil, false
	}
	_, a := balances["a"]
	_, b := balances["b"]
	return balances, a && b && balances["a"] >= 0 && balances["b"] >= 0 && balances["a"]+balances["b"] == 200
}

func (c *colonyGame) readStores(control *colonyControl) (map[string]float64, bool) {
	status, body, _, err := control.request(control.ctx, "/stores", nil)
	balances, ok := transferBalances(body)
	return balances, ok && err == nil && status == http.StatusOK
}

func (c *colonyGame) probeTransfers(control *colonyControl) bool {
	if !control.warehouse {
		return false
	}
	before, ok := c.readStores(control)
	if !ok || !c.warehouseSQLiteExists() {
		return false
	}
	// An id is deliberately shared by eight concurrent requests.
	c.mu.Lock()
	prefix := fmt.Sprintf("%d-%d-%d", control.seed, c.state.Contracts[control.contract].Start, c.state.Contracts[control.contract].Tries)
	c.mu.Unlock()
	first := map[string]any{"id": prefix + "-transfer-once", "from": "a", "to": "b", "amount": 7}
	var workers sync.WaitGroup
	passed := make(chan bool, 64)
	for range 8 {
		workers.Add(1)
		go func() {
			defer workers.Done()
			status, _, _, err := control.request(control.ctx, "/transfer", first)
			passed <- err == nil && status == http.StatusOK
		}()
	}
	workers.Wait()
	for range 8 {
		if !<-passed {
			return false
		}
	}
	after, ok := c.readStores(control)
	if !ok || after["a"] != before["a"]-7 || after["b"] != before["b"]+7 {
		return false
	}
	for index := range 32 {
		workers.Add(1)
		go func(index int) {
			defer workers.Done()
			from, to := "a", "b"
			if index%2 == 1 {
				from, to = to, from
			}
			request := map[string]any{"id": fmt.Sprintf("%s-transfer-%d", prefix, index), "from": from, "to": to, "amount": 17}
			status, _, _, err := control.request(control.ctx, "/transfer", request)
			_, invariant := c.readStores(control)
			passed <- err == nil && (status == http.StatusOK || status == http.StatusConflict) && invariant
		}(index)
	}
	workers.Wait()
	for range 32 {
		if !<-passed {
			return false
		}
	}
	balances, ok := c.readStores(control)
	if !ok {
		return false
	}
	for index := range 2 {
		status, _, _, err := control.request(control.ctx, "/transfer", map[string]any{"id": fmt.Sprintf("%s-overdraw-%d", prefix, index), "from": "a", "to": "b", "amount": 1000})
		if err != nil || status != http.StatusConflict {
			return false
		}
	}
	final, ok := c.readStores(control)
	return ok && maps.Equal(balances, final)
}
