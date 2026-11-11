package sim

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"orbital/colony/econ"
)

// The open colony: the same world, but persistent, endless, and run by an
// economy. Ticks are 100ms of real time; drone actions take ActionDur.

var groundByName = map[string]Ground{"regolith": Regolith, "soil": Soil, "ice": Ice, "rock": Rock, "stone": Stone, "tunnel": Tunnel, "shaft": Shaft}
var entityByName = map[string]Entity{"moss": Moss, "lichen": Lichen, "crystal": Crystal, "glowcap": Glowcap, "reed": Reed, "iron": Iron, "geode": Geode}

// cropUnlock names the upgrade that a crop needs before it can be planted
// or harvested. Moss is free.
var cropUnlock = map[Entity]string{Lichen: "crop.lichen", Glowcap: "crop.glowcap", Reed: "crop.reed", Crystal: "crop.crystal"}

// NewColony builds a running colony world from saved state.
func NewColony(out io.Writer, st *econ.State, seed int64) *World {
	st.Normalize()
	w := New(out, seed)
	w.Econ = st
	w.Realtime = true
	w.Unit = 100 * time.Millisecond
	w.MaxTicks = math.MaxInt32
	w.TickBase = int(st.Tick)
	w.Tick = w.TickBase
	w.loadTiles()
	w.applyEcon()
	w.Check = func(w *World) (bool, string) {
		return false, sayf("كريدت: %.0f", "Credits: %.0f", w.Econ.Credits)
	}
	return w
}

func (w *World) loadTiles() {
	st := w.Econ
	w.Init(st.Size, st.Size, Regolith)
	w.InitLayers(st.Depth)
	for z := 0; z < st.Depth; z++ {
		for y := 0; y < st.Size; y++ {
			for x := 0; x < st.Size; x++ {
				et := st.AtZ(x, y, z)
				t := w.AtZ(x, y, z)
				t.Ground = groundByName[et.G]
				t.Entity = entityByName[et.E]
				t.Planted, t.Ready, t.Size = int(et.P), int(et.R), et.S
			}
		}
	}
	// keep drones on the map after it changes size
	clamp := func(v int) int { return min(max(v, 0), st.Size-1) }
	w.X, w.Y, w.Z = clamp(w.X), clamp(w.Y), min(max(w.Z, 0), w.Depth-1)
	for _, d := range w.Drones {
		d.X, d.Y, d.Z = clamp(d.X), clamp(d.Y), min(max(d.Z, 0), w.Depth-1)
	}
}

// applyEcon refreshes everything the upgrades control.
func (w *World) applyEcon() {
	w.ActionDur = time.Duration(w.Econ.ActionMillis() * float64(time.Millisecond))
	w.MaxDrones = w.Econ.Drones()
	w.syncInv()
}

func (w *World) syncInv() {
	for name, e := range entityByName {
		w.Inv[e] = int(w.Econ.Inv[name])
	}
}

// syncEcon writes the live world back into the saved state. Caller holds the lock.
func (w *World) syncEcon() {
	st := w.Econ
	st.Tick = int64(w.Tick)
	st.Depth = w.Depth
	for z := 0; z < w.Depth; z++ {
		for y := 0; y < w.H; y++ {
			for x := 0; x < w.W; x++ {
				t := w.AtZ(x, y, z)
				et := econ.Tile{G: t.Ground.String(), P: int64(t.Planted), R: int64(t.Ready), S: t.Size}
				if t.Entity != Empty {
					et.E = t.Entity.String()
				}
				*st.AtZ(x, y, z) = et
			}
		}
	}
}

// gate returns a refusal message when a command has not been bought yet.
func (w *World) gate(id string) string {
	if w.Econ == nil || w.Econ.Has(id) {
		return ""
	}
	if u := econ.Find(id); u != nil {
		return say("🔒 محتاج تشتري «", "🔒 Buy “") + say(u.Name, u.NameEn) + say("» من المتجر", "” from the shop first")
	}
	return say("🔒 مقفول", "🔒 Locked")
}

func (w *World) neighbors4(x, y int, f func(t *Tile)) {
	for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		nx, ny := x+d[0], y+d[1]
		if nx >= 0 && ny >= 0 && nx < w.W && ny < w.H {
			f(w.AtZ(nx, ny, w.Z))
		}
	}
}

func (w *World) colonyPlant(t *Tile, e Entity) (bool, string, bool) {
	if id, ok := cropUnlock[e]; ok && e != Crystal {
		if m := w.gate(id); m != "" {
			return false, m, false
		}
	}
	switch e {
	case Moss, Glowcap:
		if t.Ground != Soil {
			return false, e.String() + say(" محتاج تربة (soil). استخدم Till الأول", " needs soil. Use Till first"), false
		}
	case Lichen:
		if t.Ground != Ice {
			return false, say("الأشنة بتتزرع على الجليد (ice) بس", "Lichen only grows on ice"), false
		}
	case Reed:
		if t.Ground != Soil {
			return false, say("البوص محتاج تربة (soil)", "Reeds need soil"), false
		}
		nearIce := false
		w.neighbors4(w.X, w.Y, func(n *Tile) { nearIce = nearIce || n.Ground == Ice })
		if !nearIce {
			return false, say("البوص لازم يبقى جنب خانة جليد", "Reeds must be next to an ice tile"), false
		}
	default:
		return false, e.String() + say(" مش حاجة بتتزرع", " cannot be planted"), false
	}
	t.Entity, t.Planted = e, w.Tick
	t.Ready = w.Tick + int(w.Econ.GrowTime(e.String()))
	return true, "", true
}

func (w *World) colonyHarvest(destroy bool) (int, error, string, bool) {
	t := w.here()
	if t.Entity == Empty || t.Entity == Beacon {
		return 0, ErrEmpty, say("مفيش حاجة هنا تتحصد", "Nothing here to harvest"), false
	}
	if id, ok := cropUnlock[t.Entity]; ok {
		if m := w.gate(id); m != "" {
			return 0, ErrLocked, m, false
		}
	}
	name := t.Entity.String()
	if t.Entity != Crystal && w.Tick < t.Ready {
		if !destroy {
			return 0, ErrNotRipe, say("لسه مش ناضج", "Not ripe yet"), false
		}
		*t = Tile{Ground: t.Ground}
		return 0, ErrNotRipe, say("النبات لسه مش ناضج… واتدمر!", "The plant was not ripe yet… and was destroyed!"), true
	}
	if t.Entity == Reed && w.Tick > t.Ready+econ.ReedWindow {
		*t = Tile{Ground: t.Ground}
		return 0, ErrNotRipe, say("البوص عفّن! اتأخرت عليه", "The reeds rotted! You waited too long"), true
	}
	base := 1.0
	switch t.Entity {
	case Crystal:
		base = float64(t.Size)
	case Reed:
		base = 5
	case Glowcap:
		w.neighbors4(w.X, w.Y, func(n *Tile) {
			if n.Entity == Glowcap && w.Tick >= n.Ready {
				base++
			}
		})
	}
	got := base * w.Econ.YieldMult(name)
	w.Econ.Inv[name] += got
	w.syncInv()
	*t = Tile{Ground: t.Ground}
	return int(got), nil, fmt.Sprintf("+%s %s", trimFloat(got), name), true
}

func trimFloat(v float64) string {
	if v == math.Trunc(v) {
		return fmt.Sprint(int(v))
	}
	return fmt.Sprintf("%.1f", v)
}

// Export sells n units of e (n <= 0 sells everything) and returns the credits earned.
func (w *World) Export(id int, e Entity, n int) (earned float64) {
	w.act(id, "export", func() (bool, string, bool) {
		if w.Econ == nil {
			return false, say("Export في المستعمرة بس", "Export is only available in colony mode"), false
		}
		name := e.String()
		amount := float64(n)
		if n <= 0 {
			amount = w.Econ.Inv[name]
		}
		if amount <= 0 || w.Econ.Inv[name] < 1 && n > 0 {
			return false, say("مفيش ", "No ") + name + say(" تتباع", " to sell"), false
		}
		earned = w.Econ.Sell(name, amount)
		w.syncInv()
		return true, fmt.Sprintf("+%.0f ⌬", earned), false
	})
	return earned
}

// Price is the current sale price of e (0 until the market screen is bought).
func (w *World) Price(e Entity) float64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil || !w.Econ.Has("api.prices") {
		return 0
	}
	return w.Econ.Price(e.String())
}

// Credits is the colony's money.
func (w *World) Credits() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return 0
	}
	return w.Econ.Credits
}

// ---------- background life of the colony ----------

type econEvent struct {
	Translations [][2]string        `json:"translations,omitempty"`
	Layers       int                `json:"layers,omitempty"`
	Op           string             `json:"op"`
	T            int                `json:"t"`
	OK           bool               `json:"ok"`
	Msg          string             `json:"msg,omitempty"`
	Credits      float64            `json:"credits"`
	Prices       map[string]float64 `json:"prices,omitempty"`
	Inv          map[string]float64 `json:"inv,omitempty"`
	State        *econ.State        `json:"state,omitempty"`
	W            int                `json:"w,omitempty"`
	H            int                `json:"h,omitempty"`
	Tiles        []tileView         `json:"tiles,omitempty"`
}

func (w *World) emitRaw(v econEvent) {
	v.Translations = takeMessages()
	b, _ := json.Marshal(v)
	w.out.Write(append(b, '\n'))
}

func (w *World) econSnapshot(op, msg string, ok, full bool) {
	ev := econEvent{Op: op, T: w.Tick, OK: ok, Msg: msg, Credits: w.Econ.Credits, Prices: w.Econ.Prices(), Inv: w.Econ.Inv}
	if full {
		w.syncEcon()
		ev.State = w.Econ
	}
	w.emitRaw(ev)
}

// RunColony starts the colony's background processes: market recovery,
// crystal veins, periodic saves, upgrade purchases from the game, and a
// final save on shutdown.
func (w *World) RunColony(commands io.Reader) {
	go w.colonyTicker()
	if commands != nil {
		go w.colonyCommands(commands)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-sig
		w.mu.Lock()
		w.clock()
		w.econSnapshot("state", say("المستعمرة اتحفظت", "Colony saved"), true, true)
		w.done = true
		w.exit(0)
	}()
}

func (w *World) colonyTicker() {
	const every = 10 // ticks: one second
	n := 0
	for {
		time.Sleep(every * w.Unit)
		w.mu.Lock()
		if w.done {
			w.mu.Unlock()
			return
		}
		w.clock()
		w.Econ.Recover(every)
		w.growVeins()
		n++
		w.econSnapshot("econ", "", true, n%3 == 0)
		w.mu.Unlock()
	}
}

// growVeins spawns and grows crystal veins on rock. Caller holds the lock.
func (w *World) growVeins() {
	if !w.Econ.Has("crop.crystal") {
		return
	}
	maxSize := w.Econ.VeinMax()
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			t := w.At(x, y)
			if t.Entity != Crystal {
				continue
			}
			if s := min(maxSize, 1+(w.Tick-t.Planted)/100); s != t.Size {
				t.Size = s
				v := w.view(x, y)
				w.emit(Event{Op: "tile", T: w.Tick, OK: true, Tile: &v})
			}
		}
	}
	every := int(w.Econ.VeinEvery())
	if w.Tick/max(1, every) == w.lastVein {
		return
	}
	w.lastVein = w.Tick / max(1, every)
	var free []Point
	w.Each(func(x, y int, t *Tile) {
		if t.Ground == Rock && t.Entity == Empty {
			free = append(free, Point{x, y})
		}
	})
	if len(free) == 0 {
		return
	}
	p := free[w.Rand.Intn(len(free))]
	t := w.At(p.X, p.Y)
	t.Entity, t.Planted, t.Size = Crystal, w.Tick, 1
	v := w.view(p.X, p.Y)
	w.emit(Event{Op: "tile", T: w.Tick, OK: true, Msg: say("عرق كريستال جديد", "New crystal vein"), Tile: &v})
}

// colonyCommands applies purchases the player makes in the game while the
// program runs. Each line is {"buy": "<upgrade id>"}.
func (w *World) colonyCommands(r io.Reader) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		var cmd struct{ Buy string }
		if json.Unmarshal(sc.Bytes(), &cmd) != nil || cmd.Buy == "" {
			continue
		}
		w.mu.Lock()
		w.clock()
		w.syncEcon()
		oldSize := w.Econ.Size
		err := w.Econ.Buy(cmd.Buy)
		msg := cmd.Buy
		if err == nil {
			if w.Econ.Size != oldSize {
				w.loadTiles()
				w.emitRaw(econEvent{Op: "resize", T: w.Tick, OK: true, W: w.W, H: w.H, Layers: w.layerCount(), Tiles: w.allTiles(), Credits: w.Econ.Credits})
			}
			w.applyEcon()
		} else {
			msg = err.Error()
		}
		w.econSnapshot("bought", msg, err == nil, true)
		w.mu.Unlock()
	}
}

func (w *World) allTiles() []tileView {
	var out []tileView
	for z := 0; z < w.Depth; z++ {
		for y := 0; y < w.H; y++ {
			for x := 0; x < w.W; x++ {
				out = append(out, w.viewZ(x, y, z))
			}
		}
	}
	return out
}
