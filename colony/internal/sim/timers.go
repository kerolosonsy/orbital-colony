package sim

import (
	"context"
	"time"
)

func (w *World) startTimers() {
	if w.Storm != nil {
		go w.runStorm()
	}
	if w.Recall != nil {
		go w.runRecall()
	}
	go w.watchdog()
}

// watchdog ends a realtime mission when its time is up, even if every
// drone is stuck (a deadlock, or a receive that never completes).
func (w *World) watchdog() {
	w.sleepUntil(w.MaxTicks)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {}
	}
	w.clock()
	done, progress := w.Check(w)
	if done {
		w.finish("goal", say("المهمة اكتملت!", "Mission complete!"), 0)
	}
	w.finish("fail", sayf("الوقت خلص (%d tick). %s — لو الدرونات وقفت فجأة، يمكن فيه deadlock أو حاجة مستنية للأبد.", "Time ran out (%d ticks). %s — If the drones stopped suddenly, check for a deadlock or an endless wait.", w.MaxTicks, progress), 3)
}

func (w *World) pos(id int) Point {
	if id == 0 {
		return Point{w.X, w.Y}
	}
	d := w.Drones[id]
	return Point{d.X, d.Y}
}

// AllAtBase reports whether every drone is at the base, and if not, which
// one is away. Caller holds the lock (setup and checks only).
func (w *World) AllAtBase() (bool, int) {
	for _, d := range w.Drones {
		if w.pos(d.ID) != w.Base {
			return false, d.ID
		}
	}
	return true, -1
}

// Recall cancels the mission context at tick At; every drone must be at the
// base Grace ticks later.
type Recall struct {
	At, Grace int
	Ctx       context.Context
	cancel    context.CancelFunc
	Over      bool
}

func NewRecall(at, grace int) *Recall {
	ctx, cancel := context.WithCancel(context.Background())
	return &Recall{At: at, Grace: grace, Ctx: ctx, cancel: cancel}
}

func (w *World) runRecall() {
	r := w.Recall
	w.sleepUntil(r.At)
	w.mu.Lock()
	if w.done {
		select {}
	}
	w.clock()
	r.cancel()
	w.emit(Event{Op: "recall", T: w.Tick, OK: true, Msg: sayf("📣 أمر رجوع! كل الدرونات لازم تبقى في القاعدة خلال %d tick", "📣 Recall! All drones must return to base within %d ticks", r.Grace)})
	w.mu.Unlock()

	w.sleepUntil(r.At + r.Grace)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {}
	}
	w.clock()
	r.Over = true
	if ok, id := w.AllAtBase(); !ok {
		p := w.pos(id)
		w.finish("fail", sayf("درون %d لسه عند (%d, %d) ومرجعش للقاعدة في الوقت.", "Drone %d is still at (%d, %d) and did not return to base in time.", id, p.X, p.Y), 4)
	}
	done, progress := w.Check(w)
	if done {
		w.finish("goal", say("كل الدرونات رجعت بأمان! ", "All drones returned safely! ")+progress, 0)
	}
	w.finish("fail", say("كل الدرونات رجعت، بس ", "All drones returned, but ")+progress, 4)
}

// ScanInfo is how a crystal behaves under DeepScan. Ticks < 0 never finishes.
type ScanInfo struct {
	Ticks, Purity int
	Scanned       bool
}

// DeepScan analyses the crystal under drone id. It takes ScanInfo.Ticks
// ticks, or forever, unless ctx ends first.
func (w *World) DeepScan(ctx context.Context, id int) (int, error) {
	d := w.drone(id)
	if d == nil {
		panic(sayf("colony: مفيش درون بالمعرّف %d", "colony: no drone with id %d", id))
	}
	d.busy.Lock()
	defer d.busy.Unlock()

	w.mu.Lock()
	p := w.pos(id)
	z := w.layer(id)
	info, ok := w.Scans[p]
	ok = ok && z == 0
	w.mu.Unlock()
	dur := time.Duration(info.Ticks) * w.Unit
	switch {
	case !ok:
		dur = w.Unit
	case info.Ticks < 0:
		dur = time.Hour
	}
	var err error
	select {
	case <-time.After(dur):
	case <-ctx.Done():
		err = ctx.Err()
	}
	d.nextAt = time.Now()

	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {}
	}
	w.clock()
	if w.Econ != nil {
		w.advanceProduction(w.Tick)
		if blocked := w.droneFailure(d, "scan"); blocked != nil {
			err = blocked
		}
	}
	ev := Event{Op: "scan", T: w.Tick, D: id, X: p.X, Y: p.Y, Z: w.layer(id), OK: err == nil && ok}
	purity := 0
	switch {
	case err != nil:
		ev.Msg = say("⏱️ الفحص اتلغى: ", "⏱️ Scan cancelled: ") + err.Error()
	case !ok:
		ev.Msg, err = say("مفيش كريستال هنا يتفحص", "No crystal here to scan"), ErrEmpty
	default:
		purity = info.Purity
		info.Scanned = true
		w.Scans[p] = info
		ev.Msg = sayf("🔬 نقاء %d%%", "🔬 Purity %d%%", purity)
	}
	w.after(ev)
	return purity, err
}
