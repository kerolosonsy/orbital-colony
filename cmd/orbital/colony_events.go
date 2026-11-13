package main

import (
	"encoding/json"
	"syscall"

	"orbital/colony/econ"
)

func (c *colonyGame) signalCleanStop(tick int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, run := c.state.ActiveContract()
	active := c.active
	if c.state.Active != econ.CleanStopContract || run == nil || active == nil || active.stopSent {
		return
	}
	if tick >= econ.CleanStopAt(c.state.Seed, run.Start) {
		active.stopSent = true
		_ = active.cmd.Process.Signal(syscall.SIGTERM)
	}
}

// Caller holds mu; broadcast after releasing it.
func (c *colonyGame) finishColonyEvents(lang string) string {
	done, settled := c.state.FinishCleanStop()
	if !settled {
		return ""
	}
	msg := tr(lang, "الإغلاق خلص من غير خسارة حمولة", "Shutdown finished without cargo loss")
	if !done {
		msg = tr(lang, "الشحنة اللي فضلت في الدرونات ضاعت", "Cargo left in drones was lost")
	}
	frac := 0
	if done {
		frac = 1
	}
	line, _ := json.Marshal(map[string]any{"op": "contract", "t": c.state.Tick, "ok": done, "msg": msg, "state": c.state, "contract": map[string]any{"id": econ.CleanStopContract, "frac": frac, "left": -1}})
	return string(line)
}
