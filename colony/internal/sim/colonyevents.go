package sim

import (
	"fmt"
	"math"
	"time"

	"orbital/colony/econ"
)

// DroneError identifies a disabled drone and supports errors.Is with ErrFault or ErrStorm.
type DroneError struct {
	Drone int
	Cause error
}

func (e DroneError) Error() string { return fmt.Sprintf("%s (%d)", e.Cause, e.Drone) }
func (e DroneError) Unwrap() error { return e.Cause }

type faultError struct{}

func (faultError) Error() string {
	return say("colony: الدرون محتاج تصليح", "colony: drone needs repair")
}

var ErrFault error = faultError{}

func (w *World) initColonyEvents() {
	w.Storm = &Storm{Warning: make(chan struct{}, 1), Clear: make(chan struct{})}
	w.Storm.active = w.Econ.Weather.Storm == "active"
	if w.Econ.Weather.Storm == "warn" {
		w.Storm.Warning <- struct{}{}
	}
	w.scheduleStorm(int64(w.Tick))
}

func (w *World) scheduleStorm(tick int64) {
	st := w.Econ
	if w.seasonStorm(tick) {
		return
	}
	if !st.Won("c13-reed-clock") || st.Weather.At != 0 {
		return
	}
	const period = 2400
	at := econ.EventAt(st.Seed, tick/period, period, 1)
	if at <= tick {
		at = econ.EventAt(st.Seed, tick/period+1, period, 1)
	}
	st.Weather = econ.Weather{Storm: "none", At: at}
}

func (w *World) advanceColonyEvents(tick int) {
	st := w.Econ
	st.Tick = int64(tick)
	st.ExpireEvents(int64(tick))
	w.recoverSeasonMarket()
	w.scheduleStorm(int64(tick))
	if st.Weather.At > 0 && int64(tick) >= st.Weather.At {
		switch st.Weather.Storm {
		case "none":
			w.colonyStormPhase("warn")
		case "warn":
			w.colonyStormPhase("hit")
		case "active":
			w.colonyStormPhase("clear")
		}
	}
	if period := st.ShockPeriod(); period > 0 && int64(tick) == econ.EventAt(st.Seed, int64(tick)/period, period, 2) {
		h := econ.EventHash(st.Seed, int64(tick), 2)
		factor := 0.4
		if h&1 != 0 {
			factor = 1.6
		}
		st.Shock = &econ.MarketShock{Good: w.eventGood(h >> 1), Factor: factor, Until: int64(tick) + econ.ShockTicks}
		w.emit(Event{Op: "shock", T: tick, OK: true, Msg: sayf("صدمة سوق: %s ×%.1f", "Market shock: %s ×%.1f", st.Shock.Good, factor)})
	}
	if st.Won("c06-ripe-ledger") && int64(tick) == econ.EventAt(st.Seed, int64(tick)/1600, 1600, 3) {
		h := econ.EventHash(st.Seed, int64(tick), 3)
		n := 20 + int(h%21)
		st.Rush = &econ.RushOrder{Good: w.eventGood(h >> 5), Amount: n, Left: float64(n), Until: int64(tick) + econ.RushTicks, Bonus: econ.RushBonus}
		w.emit(Event{Op: "rush", T: tick, OK: true, Msg: sayf("طلب عاجل: %d %s بسعر مضاعف", "Rush order: %d %s at double price", n, st.Rush.Good)})
	}
}

func (w *World) eventGood(hash uint64) string {
	goods := []string{"moss"}
	for _, good := range []string{"lichen", "glowcap", "reed", "crystal"} {
		if w.Econ.Usable("crop." + good) {
			goods = append(goods, good)
		}
	}
	if w.Econ.Usable("api.dig") {
		goods = append(goods, "iron", "geode")
	}
	return goods[hash%uint64(len(goods))]
}

func (w *World) colonyStormPhase(phase string) {
	st := w.Econ
	msg := ""
	switch phase {
	case "warn":
		st.Weather = econ.Weather{Storm: "warn", At: st.Weather.At + econ.StormWarningTicks}
		select {
		case <-w.Storm.Clear:
			w.Storm.Clear = make(chan struct{})
		default:
		}
		select {
		case w.Storm.Warning <- struct{}{}:
		default:
		}
		msg = say("⚠️ عاصفة بعد 60 tick! ارجع للقاعدة", "⚠️ Storm in 60 ticks! Return to base")
	case "hit":
		w.Storm.active = true
		st.Weather = econ.Weather{Storm: "active", At: st.Weather.At + econ.StormStunTicks}
		for _, d := range w.Drones {
			x, y, z := d.X, d.Y, d.Z
			if d.ID == 0 {
				x, y, z = w.X, w.Y, w.Z
			}
			if x != 0 || y != 0 || z != 0 {
				d.stunnedUntil = int(st.Tick) + econ.StormStunTicks
			}
		}
		w.Each(func(x, y int, t *Tile) {
			if (t.Entity == Moss || t.Entity == Lichen || t.Entity == Glowcap || t.Entity == Reed) && t.Ready > int(st.Tick) {
				t.Entity, t.Planted, t.Ready, t.Size = Empty, 0, 0, 0
				w.Econ.RecordLost()
				v := w.view(x, y)
				w.emit(Event{Op: "tile", T: int(st.Tick), OK: true, Tile: &v})
			}
		})
		msg = say("🌪️ العاصفة ضربت! الدرونات المكشوفة اتعطّلت", "🌪️ Storm hit! Exposed drones are stunned")
	case "clear":
		w.Storm.active = false
		close(w.Storm.Clear)
		st.Weather = econ.Weather{Storm: "none"}
		w.scheduleStorm(st.Tick)
		msg = say("☀️ العاصفة خلصت", "☀️ Storm cleared")
	}
	w.emit(Event{Op: "storm", Phase: phase, T: int(st.Tick), OK: true, Msg: msg})
}

func (w *World) droneFailure(d *Drone, op string) error {
	if w.Econ == nil {
		return nil
	}
	if d.stunnedUntil > w.Tick {
		return DroneError{d.ID, ErrStorm}
	}
	if op == "drone-repair" {
		return nil
	}
	if d.fault {
		return DroneError{d.ID, ErrFault}
	}
	if w.Econ.Won("c05-careful-harvest") && econ.EventHash(w.Econ.Seed, int64(w.Tick), uint64(d.ID)+10)%100 == 0 {
		d.fault = true
		w.emit(Event{Op: "fault", D: d.ID, T: w.Tick, OK: false, Msg: ErrFault.Error()})
		return DroneError{d.ID, ErrFault}
	}
	return nil
}

func (w *World) Fault(id int) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return nil
	}
	d := w.Drones[id]
	if d.stunnedUntil > w.Tick {
		return DroneError{id, ErrStorm}
	}
	if d.fault {
		return DroneError{id, ErrFault}
	}
	return nil
}

func (w *World) Repair(id int) bool {
	return w.act(id, "drone-repair", func() (bool, string, bool) {
		w.actingDrone.fault = false
		return true, say("الدرون اتصلّح", "Drone repaired"), false
	})
}

func (w *World) Rush() (Entity, int, bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil || w.Econ.Rush == nil || int64(w.Tick) >= w.Econ.Rush.Until {
		return Empty, 0, false
	}
	r := w.Econ.Rush
	return entityByName[r.Good], int(math.Ceil(r.Left)), true
}

func (w *World) StormWarning() <-chan struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Storm == nil {
		return nil
	}
	return w.Storm.Warning
}

func (w *World) StormClear() <-chan struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Storm == nil {
		return nil
	}
	return w.Storm.Clear
}

func (w *World) Shutdown() <-chan struct{} {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.shutdown == nil {
		w.shutdown = make(chan struct{})
	}
	if w.shuttingDown && !w.shutdownClosed {
		close(w.shutdown)
		w.shutdownClosed = true
	}
	return w.shutdown
}

func (w *World) colonyShutdown() {
	w.mu.Lock()
	w.clock()
	w.shuttingDown = true
	registered := w.shutdown != nil
	if registered {
		close(w.shutdown)
		w.shutdownClosed = true
	}
	w.saveShutdown(false)
	w.mu.Unlock()
	if registered {
		deadline := time.NewTimer(2 * time.Second)
		defer deadline.Stop()
		<-deadline.C
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.clock()
	w.saveShutdown(true)
	w.done = true
	w.exit(0)
}

func (w *World) saveShutdown(final bool) {
	if !w.shuttingDown {
		return
	}
	lost := 0.0
	for _, d := range w.Drones {
		lost += cargoAmount(d)
	}
	if w.Econ.Active == econ.CleanStopContract {
		_, r := w.Econ.ActiveContract()
		r.Stopping = true
		if final {
			r.Stopped = true
			r.CargoLost += lost
			w.checkContract()
		}
	}
	if final {
		for _, d := range w.Drones {
			clear(d.Hold)
			d.Cargo = 0
		}
		w.Cargo = 0
	}
	w.econSnapshot("state", say("المستعمرة اتحفظت", "Colony saved"), true, true)
}

func colonyEventWire(ev Event) any {
	if ev.Op == "fault" {
		return struct {
			Event
			Drone int `json:"d"`
		}{ev, ev.D}
	}
	return ev
}
