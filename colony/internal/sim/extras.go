package sim

import (
	"time"
)

// Signal is a message picked up by the colony radio. Its concrete type is
// one of Distress, Supply or Noise.
type Signal interface{ isSignal() }

// Distress is a call for help: fly to At.
type Distress struct{ At Point }

// Supply announces a supply drop worth Amount.
type Supply struct{ Amount int }

// Noise is static; it carries no information.
type Noise struct{ Text string }

func (Distress) isSignal() {}
func (Supply) isSignal()   {}
func (Noise) isSignal()    {}

// Unload moves the cargo into storage. It only works at the base.
func (w *World) Unload(id int) bool {
	return w.act(id, "unload", func() (bool, string, bool) {
		if (Point{w.X, w.Y}) != w.Base || w.Econ != nil && w.Z != 0 {
			return false, sayf("التفريغ في القاعدة بس %v", "Unload only at base %v", w.Base), false
		}
		if w.Econ != nil {
			return w.unloadColony()
		}
		if w.Cargo == 0 {
			return false, say("الشحنة فاضية", "Cargo is empty"), false
		}
		n := w.Cargo
		w.Inv[Crystal] += n
		w.Cargo = 0
		return true, sayf("📦 +%d كريستال", "📦 +%d crystal", n), false
	})
}

// Nearest returns the closest remaining target crystal, or nil.
func (w *World) Nearest(id int) (p *Point) {
	w.act(id, "scan", func() (bool, string, bool) {
		if m := w.gate("api.radar"); m != "" {
			return false, m, false
		}
		here, best := Point{w.X, w.Y}, -1
		for i, t := range w.Targets {
			if w.AtZ(t.X, t.Y, w.Z).Entity != Crystal {
				continue
			}
			if best < 0 || Dist(here, t) < Dist(here, w.Targets[best]) {
				best = i
			}
		}
		if best >= 0 {
			q := w.Targets[best]
			p = &q
			return true, "", false
		}
		return true, say("مفيش أهداف فاضلة (nil)", "No targets left (nil)"), false
	})
	return p
}

// Peek reads a tile for the autopilot without spending a tick.
func (w *World) Peek(x, y int) (g Ground, e Entity, ripe bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	t := w.AtZ(x, y, w.Z)
	ripe = t.Entity == Crystal || (t.Entity == Moss || t.Entity == Lichen) && w.Tick >= t.Ready
	return t.Ground, t.Entity, ripe
}

// SignalList returns a copy of the received signals.
func (w *World) SignalList() []Signal {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]Signal(nil), w.Signals...)
}

// Spawn launches a new drone at the base and returns its id, or -1 when the
// hangar is empty.
func (w *World) Spawn() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {}
	}
	if w.Econ != nil {
		w.MaxDrones = w.Econ.DroneCap()
		if !w.Econ.Usable("api.spawn") {
			w.emit(Event{Op: "spawn", T: w.Tick, OK: false, Msg: say("🔒 colony.Spawn محتاجة «رخصة الهانجر» من المتجر", "🔒 colony.Spawn needs a Hangar License from the shop")})
			return -1
		}
	}
	if len(w.Drones) >= max(1, w.MaxDrones) {
		return -1
	}
	w.clock()
	if w.Econ != nil {
		w.advanceProduction(w.Tick)
	}
	d := &Drone{ID: len(w.Drones), X: w.Base.X, Y: w.Base.Y, Battery: 100, Hold: map[string]float64{}}
	w.Drones = append(w.Drones, d)
	if w.Econ != nil {
		w.syncInfrastructure()
	}
	w.emit(Event{Op: "spawn", T: w.Tick, D: d.ID, OK: true, X: d.X, Y: d.Y, Msg: sayf("🚁 درون %d انطلق", "🚁 Drone %d launched", d.ID)})
	return d.ID
}

// Storm is a dust storm with a warning, a hit and a clear time (in ticks).
type Storm struct {
	WarnAt, HitAt, ClearAt int
	Warning, Clear         chan struct{}
	active                 bool
}

func (w *World) sleepUntil(tick int) {
	time.Sleep(time.Until(w.start.Add(time.Duration(tick) * w.Unit)))
}

func (w *World) stormEvent(op, msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {}
	}
	w.clock()
	switch op {
	case "storm":
		w.Storm.active = true
		for _, d := range w.Drones {
			x, y := d.X, d.Y
			if d.ID == 0 {
				x, y = w.X, w.Y
			}
			if (Point{x, y}) != w.Base {
				w.emit(Event{Op: op, T: w.Tick, OK: false, Msg: msg})
				w.finish("fail", sayf("العاصفة ضربت ودرون %d كان عند (%d, %d) برّه القاعدة. المهمة فشلت.", "The storm hit while drone %d was at (%d, %d), outside base. Mission failed.", d.ID, x, y), 4)
			}
		}
	case "clear":
		w.Storm.active = false
	}
	w.emit(Event{Op: op, T: w.Tick, OK: true, Msg: msg})
}

func (w *World) runStorm() {
	s := w.Storm
	w.sleepUntil(s.WarnAt)
	w.stormEvent("storm-warning", sayf("⚠️ إنذار: عاصفة هتضرب بعد %d tick!", "⚠️ Warning: a storm will hit in %d ticks!", s.HitAt-s.WarnAt))
	s.Warning <- struct{}{}
	w.sleepUntil(s.HitAt)
	w.stormEvent("storm", say("🌪️ العاصفة ضربت!", "🌪️ The storm has hit!"))
	w.sleepUntil(s.ClearAt)
	w.stormEvent("clear", say("☀️ الجو صفي. ارجعوا للشغل", "☀️ The storm has cleared. Back to work!"))
	close(s.Clear)
}
