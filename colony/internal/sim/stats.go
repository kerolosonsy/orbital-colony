package sim

import "time"

type droneStats struct {
	ID        int `json:"id"`
	Work      int `json:"work"`
	Move      int `json:"move"`
	Wait      int `json:"wait"`
	Failed    int `json:"failed"`
	IdleMoves int `json:"idleMoves"`
}

type runStats struct {
	Drones  []droneStats `json:"drones"`
	TopFail string       `json:"topFail"`
	Ticks   int          `json:"ticks"`
}

type droneCounters struct {
	droneStats
	commands map[string]int
	failures map[string]int
	waitTime time.Duration
	pending  bool
}

// This wrapper runs under World.mu and observes position before and after the
// command without adding state or bookkeeping to every drone operation.
func (w *World) trackAction(d *Drone, op string, f func() (bool, string, bool)) func() (bool, string, bool) {
	return func() (bool, string, bool) {
		before := [3]int{w.X, w.Y, w.Z}
		ok, msg, changed := f()
		if w.stats == nil {
			w.stats = map[int]*droneCounters{}
		}
		if w.stats[d.ID] == nil {
			w.stats[d.ID] = &droneCounters{droneStats: droneStats{ID: d.ID}, commands: map[string]int{}, failures: map[string]int{}}
		}
		c := w.stats[d.ID]
		c.commands[op]++
		if !ok {
			c.Failed++
			c.failures[op]++
		}
		w.recordActivity(c, op, ok)
		c.Move += abs(w.X-before[0]) + abs(w.Y-before[1]) + abs(w.Z-before[2])
		return ok, msg, changed
	}
}

func (w *World) recordActivity(c *droneCounters, op string, ok bool) {
	switch op {
	case "move":
		if c.pending {
			c.IdleMoves++
			c.pending = false
		}
		if ok {
			c.pending = true
		}
	case "wait":
		duration := w.Unit
		if w.Realtime && w.ActionDur > 0 {
			duration = w.ActionDur
		}
		if w.batteryActive() && w.Drones[c.ID].Battery == 0 {
			duration *= 2
		}
		c.waitTime += duration
		c.Wait = int(c.waitTime / w.Unit)
	default:
		c.Work++
		if ok {
			c.pending = false
		}
	}
}

func (w *World) statsSnapshot() *runStats {
	snapshot := &runStats{Drones: make([]droneStats, len(w.Drones)), Ticks: max(0, w.Tick-w.TickBase)}
	failures := map[string]int{}
	for id := range w.Drones {
		snapshot.Drones[id].ID = id
		if c := w.stats[id]; c != nil {
			snapshot.Drones[id] = c.droneStats
			for op, count := range c.failures {
				failures[op] += count
			}
		}
	}
	most := 0
	for op, count := range failures {
		if count > most || count == most && op < snapshot.TopFail {
			snapshot.TopFail, most = op, count
		}
	}
	return snapshot
}
