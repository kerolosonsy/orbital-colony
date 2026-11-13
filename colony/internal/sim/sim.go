// Package sim is the colony world simulation. Player programs cannot import
// it directly (it is internal); they reach it through package colony.
package sim

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"orbital/colony/econ"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

func (d Direction) String() string {
	if d < North || d > West {
		return fmt.Sprintf("Direction(%d)", int(d))
	}
	return [...]string{"North", "East", "South", "West"}[d]
}

type Ground int

const (
	Regolith Ground = iota
	Soil
	Ice
	Rock
	Stone
	Tunnel
	Shaft
)

func (g Ground) String() string {
	if g < Regolith || g > Shaft {
		return fmt.Sprintf("Ground(%d)", int(g))
	}
	return [...]string{"regolith", "soil", "ice", "rock", "stone", "tunnel", "shaft"}[g]
}

type Entity int

const (
	Empty Entity = iota
	Moss
	Lichen
	Crystal
	Beacon
	Glowcap
	Reed
	Iron
	Geode
	Water
	Silica
	Copper
	Sulfur
	Titanium
	Helium
	Uranium
	Biomass
	Plate
	Wire
	Glass
	Fuel
	Fertilizer
	Circuit
	Alloy
	Battery
	Coolant
	Processor
	Frame
	Cell
	Core
)

func (e Entity) String() string {
	if e < Empty || e > Core {
		return fmt.Sprintf("Entity(%d)", int(e))
	}
	return [...]string{"empty", "moss", "lichen", "crystal", "beacon", "glowcap", "reed", "iron", "geode", "water", "silica", "copper", "sulfur", "titanium", "helium", "uranium", "biomass", "plate", "wire", "glass", "fuel", "fertilizer", "circuit", "alloy", "battery", "coolant", "processor", "frame", "cell", "core"}[e]
}

// GrowTime is how many ticks a plant needs after planting.
func GrowTime(e Entity) int {
	switch e {
	case Moss:
		return 8
	case Lichen:
		return 12
	}
	return 0
}

var (
	ErrNotRipe = errors.New(say("colony: النبات لسه مش ناضج", "colony: plant is not ripe yet"))
	ErrEmpty   = errors.New(say("colony: مفيش حاجة هنا تتحصد", "colony: nothing to harvest here"))
	ErrStorm   = errors.New(say("colony: العاصفة الترابية عطّلت الدرون", "colony: dust storm interrupted the drone"))
	ErrFull    = econ.ErrFull
	ErrLocked  = errors.New(say("colony: الأمر ده لسه مقفول", "colony: this command is not unlocked yet"))
)

// Point is a position on the map. (0, 0) is the bottom-left corner.
type Point struct {
	X, Y int
}

// Dist is the Manhattan distance: how many moves it takes to get from p to q.
func Dist(p, q Point) int { return abs(p.X-q.X) + abs(p.Y-q.Y) }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type Tile struct {
	Ground   Ground
	Entity   Entity
	Planted  int
	Ready    int
	Size     int
	MeltAt   int
	SiftUsed int
	SiftAt   int
}

type tileView struct {
	B string `json:"b,omitempty"`
	Z int    `json:"z,omitempty"`
	X int    `json:"x"`
	Y int    `json:"y"`
	G string `json:"g"`
	E string `json:"e,omitempty"`
	P int    `json:"p,omitempty"`
	R int    `json:"r,omitempty"`
	S int    `json:"s,omitempty"`
}

// Wire format: tile.b is the building kind, omitted on tiles without a building.
// Econ snapshots carry rank:{level,kind,target,value,filled}; state.rank is the earned level.
// filled is running time / econ.Window, capped at 1; kind names the next goal.
// Full econ snapshots (all events carrying state) include stats:{drones:[{id,work,move,
// wait,failed,idleMoves}],topFail,ticks}. work counts non-move/non-wait commands;
// move counts cells traveled, wait counts action ticks, ticks excludes earlier runs.
// topFail is empty without failures; ties use lexical op order. Idle moves are
// successful moves followed by another move attempt without successful work.
// Shop and tree upgrades use max:0 for unlimited levels, never a maxed check.
// op:"rank" carries ok:true, a localized msg and rank after each promotion.
// op:"building" carries t, x, y and building:{kind,recipe,in,out,progress,wear,level,
// powered,enabled,priority,mode}; mode is present only on generators.
// state includes power:{supply,demand,day}, storage:{crops,raw,crafted}, where
// each group is {used,cap}, fleet:[{id,battery,cargo,cap}], and rent (credits/minute).
// Buffers map good IDs to amounts; progress is a fraction in [0,1]. Events are
// sent on placement, buffer/recipe changes, repair and completed cycles, not each tick.
// Lab consumption also emits op:"building" and never adds output goods.
// State includes rp (float); tree nodes include rp (first-purchase cost, omitted
// at zero), and checks may have kind:"rp". Contracts include repeatable:true
// and rewardRP (base RP); settlement rewards include rp (actual RP paid).
// state.weather is {storm:"none"|"warn"|"active",at:nextPhaseTick}.
// Optional state.rush is {good,amount,left,until,bonus}; state.shock is {good,factor,until}.
// Fleet entries optionally include fault:true and stunned:ticksRemaining.
// Live op:"storm" has phase:"warn"|"hit"|"clear" and msg; op:"fault" has d,msg;
// op:"shock" and op:"rush" have msg.
type Event struct {
	Phase        string         `json:"phase,omitempty"`
	Building     *BuildingView  `json:"building,omitempty"`
	Translations [][2]string    `json:"translations,omitempty"`
	Z            int            `json:"z,omitempty"`
	Layers       int            `json:"layers,omitempty"`
	Op           string         `json:"op"`
	T            int            `json:"t"`
	OK           bool           `json:"ok"`
	X            int            `json:"x"`
	Y            int            `json:"y"`
	D            int            `json:"d,omitempty"`
	Msg          string         `json:"msg,omitempty"`
	Tile         *tileView      `json:"tile,omitempty"`
	Inv          map[string]int `json:"inv,omitempty"`
	Progress     string         `json:"progress,omitempty"`
	W            int            `json:"w,omitempty"`
	H            int            `json:"h,omitempty"`
	Tiles        []tileView     `json:"tiles,omitempty"`
	MaxTicks     int            `json:"maxTicks,omitempty"`
	Par          int            `json:"par,omitempty"`
	Gold         int            `json:"gold,omitempty"`
	// Reached tells, after each action, whether the goal is true right now.
	// The server reads it when the program ends on its own.
	Reached *bool `json:"reached,omitempty"`
}

// World holds the full state of one mission. All exported methods are safe
// for concurrent use, so later chapters can drive it from goroutines.
type World struct {
	archive    *colonyArchive
	marketKeys map[string]systemReceipt

	mu      sync.Mutex
	W, H    int
	tiles   []Tile
	under   [][]Tile
	Depth   int
	X, Y, Z int
	Tick    int

	MaxTicks       int
	Inv            map[Entity]int
	Reports        []string
	Destroyed      int
	StormChance    float64
	DisableHarvest bool
	FailOnDestroy  bool
	Rand           *rand.Rand
	Check          func(w *World) (done bool, progress string)
	FailCheck      func(w *World) string // non-empty result fails the mission
	// Reaching the goal doesn't end the run: the next action has to keep it
	// (an endless loop keeps working), or the program has to end there.
	reached      bool
	reachedAt    int // tick the goal became true
	HarvestLimit int // 0 means unlimited
	Harvests     int
	Waypoints    []Point
	NextWaypoint int
	Targets      []Point
	Par, Gold    int  // per-run star thresholds; override the level's when set
	AnyOrder     bool // waypoints may be visited in any order
	Capacity     int  // cargo limit; 0 means harvests go straight to storage
	Cargo        int
	Base         Point // where Unload works
	Signals      []Signal
	MaxDrones    int           // drones allowed in total, including the first
	Realtime     bool          // ticks follow wall-clock time, so goroutines run in parallel
	Unit         time.Duration // wall-clock length of one tick in realtime mode
	Drones       []*Drone
	Storm        *Storm
	Recall       *Recall
	Scans        map[Point]ScanInfo // deep-scan behaviour per crystal
	RockBlocks   bool               // drones cannot enter Rock tiles
	HideMap      bool               // Grid() returns nil: the drone must explore
	MoveCost     map[Ground]int     // ticks to enter a tile; default 1
	Data         Datasets
	Console      io.Writer     // when set, a human-readable log for terminal runs
	Econ         *econ.State   // non-nil in the open colony: economy, upgrades, gates
	ActionDur    time.Duration // wall time of one drone action; defaults to Unit
	TickBase     int           // ticks already played in earlier colony runs
	lastVein     int
	ScoutRoles   bool // drone 0 cannot scan; other drones cannot harvest
	start        time.Time
	mainSave     [4]int

	shutdown       chan struct{}
	shuttingDown   bool
	shutdownClosed bool
	stats          map[int]*droneCounters

	lastRankTick int

	lastProductionTick int
	actingDrone        *Drone

	out  io.Writer
	exit func(code int)
	done bool
}

func New(out io.Writer, seed int64) *World {
	return &World{
		out:      out,
		Depth:    1,
		exit:     os.Exit,
		Inv:      map[Entity]int{},
		Rand:     rand.New(rand.NewSource(seed)),
		MaxTicks: 1000,
		Check:    func(*World) (bool, string) { return false, "" },
		Unit:     4 * time.Millisecond,
		Drones:   []*Drone{{Battery: econ.BatteryCapacity}},
	}
}

// Drone is one drone's state. Drone 0's position and cargo live in
// World.X, World.Y, World.Z and World.Cargo; others are swapped in while they act.
type Drone struct {
	fault              bool
	stunnedUntil       int
	ID, X, Y, Z, Cargo int
	Battery            int
	Hold               map[string]float64
	busy               sync.Mutex
	nextAt             time.Time // when this drone's previous action finished
}

func (w *World) load(d *Drone) {
	if d.ID == 0 {
		return
	}
	w.mainSave = [4]int{w.X, w.Y, w.Z, w.Cargo}
	w.X, w.Y, w.Z, w.Cargo = d.X, d.Y, d.Z, d.Cargo
}

func (w *World) store(d *Drone) {
	if d.ID == 0 {
		return
	}
	d.X, d.Y, d.Z, d.Cargo = w.X, w.Y, w.Z, w.Cargo
	w.X, w.Y, w.Z, w.Cargo = w.mainSave[0], w.mainSave[1], w.mainSave[2], w.mainSave[3]
}

func (w *World) drone(id int) *Drone {
	w.mu.Lock()
	defer w.mu.Unlock()
	if id < 0 || id >= len(w.Drones) {
		return nil
	}
	return w.Drones[id]
}

func (w *World) clock() {
	if w.Realtime {
		w.Tick = max(w.Tick, w.TickBase+int(time.Since(w.start)/w.Unit))
	}
}

// Init sizes the map and fills it with one ground type. Setup code only.
func (w *World) Init(width, height int, g Ground) {
	w.W, w.H = width, height
	w.Depth, w.under = 1, nil
	w.tiles = make([]Tile, width*height)
	for i := range w.tiles {
		w.tiles[i].Ground = g
	}
}

// At returns the tile at (x, y). Setup and goal checks only.
func (w *World) At(x, y int) *Tile { return &w.tiles[y*w.W+x] }

// Each calls f for every tile. Setup and goal checks only.
func (w *World) Each(f func(x, y int, t *Tile)) {
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			f(x, y, w.At(x, y))
		}
	}
}

// PlantAt places a plant that becomes ripe at tick ready. Setup only.
func (w *World) PlantAt(x, y int, e Entity, ready int) {
	t := w.At(x, y)
	t.Entity, t.Planted, t.Ready = e, ready-GrowTime(e), ready
}

func (w *World) view(x, y int) tileView {
	return w.viewZ(x, y, 0)
}

func (w *World) viewZ(x, y, z int) tileView {
	t := w.AtZ(x, y, z)
	v := tileView{X: x, Y: y, Z: z, G: t.Ground.String()}
	if z == 0 {
		if b := w.buildingAt(x, y); b != nil {
			v.B = b.Kind
		}
	}
	if t.Entity != Empty {
		v.E = t.Entity.String()
		v.P, v.R, v.S = t.Planted, t.Ready, t.Size
	}
	return v
}

func (w *World) invView() map[string]int {
	m := map[string]int{}
	for e, n := range w.Inv {
		if n != 0 {
			m[e.String()] = n
		}
	}
	if w.Capacity > 0 {
		m["cargo"] = w.Cargo
	}
	return m
}

func (w *World) emit(ev Event) {
	w.archiveEvent(ev)
	pairs := takeMessages()
	if w.Econ != nil {
		ev.Translations = pairs
	}
	b, _ := json.Marshal(colonyEventWire(ev))
	w.out.Write(append(b, '\n'))
	if w.Console != nil {
		w.printConsole(ev)
	}
}

// printConsole narrates the mission for players running go run in a terminal.
func (w *World) printConsole(ev Event) {
	who := ""
	if ev.D > 0 {
		who = sayf("[درون %d] ", "[Drone %d] ", ev.D)
	}
	switch {
	case ev.Op == "goal":
		fmt.Fprintf(w.Console, "✅ %s (%d tick)\n", ev.Msg, ev.T)
	case ev.Op == "fail":
		fmt.Fprintf(w.Console, "✖ %s\n", ev.Msg)
	case ev.Op == "report":
		fmt.Fprintf(w.Console, say("📡 تقرير → %s   (%s)\n", "📡 Report → %s   (%s)\n"), ev.Msg, ev.Progress)
	case !ev.OK && ev.Msg != "":
		fmt.Fprintf(w.Console, "t=%d %s%s: %s\n", ev.T, who, ev.Op, ev.Msg)
	case ev.Op == "spawn" || ev.Op == "recall" || strings.HasPrefix(ev.Op, "storm") || ev.Op == "clear":
		fmt.Fprintf(w.Console, "t=%d %s\n", ev.T, ev.Msg)
	}
}

// Start emits the initial snapshot. Call once after setup.
func (w *World) Start() {
	w.mu.Lock()
	defer w.mu.Unlock()
	ev := Event{Op: "init", OK: true, X: w.X, Y: w.Y, W: w.W, H: w.H, MaxTicks: w.MaxTicks, Inv: w.invView(), Par: w.Par, Gold: w.Gold}
	ev.Z, ev.Layers, ev.Tiles = w.Z, w.layerCount(), w.allTiles()
	_, ev.Progress = w.Check(w)
	w.emit(ev)
	w.start = time.Now()
}

// Launch starts the mission clock's background events (storms, recalls and
// the realtime watchdog). Only the player's program calls it; the server
// builds worlds for previews and must never run timers that call exit.
func (w *World) Launch() {
	w.mu.Lock()
	w.start = time.Now()
	w.mu.Unlock()
	if w.Realtime && w.Econ == nil {
		w.startTimers()
	}
}

func (w *World) finish(op, msg string, code int) {
	w.done = true
	t := w.Tick
	if op == "goal" && w.reached {
		t = w.reachedAt
	}
	w.emit(Event{Op: op, T: t, OK: op == "goal", X: w.X, Y: w.Y, Z: w.Z, Msg: msg, Par: w.Par, Gold: w.Gold})
	w.exit(code)
}

// after runs the end-of-action bookkeeping. Caller holds the lock.
func (w *World) after(ev Event) {
	done, progress := w.Check(w)
	held := done && w.reached // still true one action after it was reached
	if done && !w.reached {
		w.reachedAt = w.Tick
	}
	w.reached = done
	ev.Progress = progress
	ev.Inv = w.invView()
	ev.Reached = &done
	w.emit(ev)
	fail := ""
	if w.FailCheck != nil {
		fail = w.FailCheck(w)
	}
	switch {
	case fail != "":
		w.finish("fail", fail, 4)
	case w.FailOnDestroy && w.Destroyed > 0:
		w.finish("fail", say("دمّرت نبات قبل ما ينضج! المهمة فشلت.", "You destroyed a plant before it ripened! Mission failed."), 4)
	case held || done && w.Tick >= w.MaxTicks:
		w.finish("goal", say("المهمة اكتملت!", "Mission complete!"), 0)
	case done:
		// reached: let the program show whether it stops here
	case w.Tick >= w.MaxTicks:
		w.finish("fail", sayf("طاقة الدرون خلصت بعد %d tick. حاول بخطوات أقل.", "The drone ran out of energy after %d ticks. Try fewer steps.", w.Tick), 3)
	}
}

// act spends one tick on an action of drone id. In realtime mode the drone
// is busy for one tick of wall-clock time, so separate drones overlap.
func (w *World) act(id int, op string, f func() (ok bool, msg string, changed bool)) bool {
	ok, _ := w.actErr(id, op, f)
	return ok
}

func (w *World) actErr(id int, op string, f func() (ok bool, msg string, changed bool)) (bool, error) {
	d := w.drone(id)
	if d == nil {
		panic(sayf("colony: مفيش درون بالمعرّف %d", "colony: no drone with id %d", id))
	}
	d.busy.Lock()
	defer d.busy.Unlock()
	// Purchases can change the action duration while other drones are waiting.
	w.mu.Lock()
	speed := w.actionSpeed(d)
	step := w.Unit
	if w.ActionDur > 0 {
		step = w.ActionDur
	}
	start := w.start
	w.mu.Unlock()
	var slot int
	if w.Realtime {
		// Each drone works on a fixed beat: its next action ends one Unit
		// after the previous one, unless it sat idle. Ticks come from that
		// schedule rather than from noisy sleep wake-ups.
		base := d.nextAt
		if now := time.Now(); base.IsZero() || now.Sub(base) > step/2 {
			base = now
		}
		step *= time.Duration(speed)
		target := base.Add(step)
		time.Sleep(time.Until(target))
		d.nextAt = target
		slot = w.TickBase + int((target.Sub(start)+w.Unit/2)/w.Unit)
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {} // mission over; os.Exit is on its way
	}
	if w.Realtime {
		w.Tick = max(w.Tick, slot)
	} else {
		w.Tick += speed
	}
	if w.Econ != nil {
		w.advanceProduction(w.Tick)
	}
	w.load(d)
	w.actingDrone = d
	if op != "wait" && w.batteryActive() {
		d.Battery = max(0, d.Battery-1)
	}
	blocked := w.droneFailure(d, op)
	if blocked != nil {
		blockedMsg := blocked.Error()
		f = func() (bool, string, bool) { return false, blockedMsg, false }
	}
	f = w.trackAction(d, op, f)
	ok, msg, changed := f()
	w.actingDrone = nil
	if w.Econ != nil {
		w.Econ.RecordCommand(id, op == "move" && ok)
		w.recordRankCommand()
	}
	if w.Econ == nil && w.Storm != nil && w.Storm.active && (Point{w.X, w.Y}) != w.Base {
		msg = say("الدرون كان برّه القاعدة وقت العاصفة!", "The drone was outside base during the storm!")
		ev := Event{Op: op, T: w.Tick, D: id, OK: false, X: w.X, Y: w.Y, Z: w.Z, Msg: msg}
		w.store(d)
		w.emit(ev)
		w.finish("fail", sayf("درون %d اتدمر في العاصفة. لازم كل الدرونات تبقى في القاعدة (0, 0) قبل ما تضرب.", "Drone %d was destroyed in the storm. All drones must reach base (0, 0) before it hits.", id), 4)
	}
	ev := Event{Op: op, T: w.Tick, D: id, OK: ok, X: w.X, Y: w.Y, Z: w.Z, Msg: msg}
	if changed {
		v := w.viewZ(w.X, w.Y, w.Z)
		ev.Tile = &v
	}
	w.store(d)
	if w.Econ != nil {
		w.syncInfrastructure()
	}
	w.after(ev)
	if w.Econ != nil {
		w.checkRank()
	}
	if w.Econ != nil {
		w.saveShutdown(false)
	}
	return ok, blocked
}

func (w *World) Move(id int, d Direction) bool {
	return w.act(id, "move", func() (bool, string, bool) {
		dx, dy := 0, 0
		switch d {
		case North:
			dy = 1
		case South:
			dy = -1
		case East:
			dx = 1
		case West:
			dx = -1
		default:
			return false, say("اتجاه مش معروف", "Unknown direction"), false
		}
		nx, ny := w.X+dx, w.Y+dy
		if nx < 0 || ny < 0 || nx >= w.W || ny >= w.H {
			return false, say("حيطة! مينفعش تطلع برا الخريطة ناحية ", "Wall! You cannot leave the map to the ") + d.String(), false
		}
		if w.Z > 0 && !walkable(w.AtZ(nx, ny, w.Z).Ground) {
			return false, say("صخر صلب! احفر الأول بـ Dig", "Solid rock! Use Dig first"), false
		}
		if w.RockBlocks && w.AtZ(nx, ny, w.Z).Ground == Rock {
			return false, say("صخرة! مفيش طريق ناحية ", "Rock! No path to the ") + d.String(), false
		}
		if c := w.MoveCost[w.AtZ(nx, ny, w.Z).Ground]; c > 1 && !w.Realtime {
			w.Tick += c - 1
		}
		w.X, w.Y = nx, ny
		if w.AnyOrder && w.AtZ(nx, ny, w.Z).Entity == Beacon {
			w.NextWaypoint++
			w.AtZ(nx, ny, w.Z).Entity = Empty
			return true, sayf("📍 نقطة %d / %d", "📍 Waypoint %d / %d", w.NextWaypoint, len(w.Waypoints)), true
		}
		if w.Z == 0 && !w.AnyOrder && w.NextWaypoint < len(w.Waypoints) && w.Waypoints[w.NextWaypoint] == (Point{nx, ny}) {
			w.NextWaypoint++
			w.AtZ(nx, ny, w.Z).Entity = Empty
			return true, sayf("📍 نقطة %d / %d", "📍 Waypoint %d / %d", w.NextWaypoint, len(w.Waypoints)), true
		}
		return true, "", false
	})
}

func (w *World) harvest(destroy bool) (int, error, string, bool) {
	if w.Econ != nil {
		return w.colonyHarvest(destroy)
	}
	t := w.here()
	if t.Entity == Empty || t.Entity == Beacon {
		return 0, ErrEmpty, say("مفيش حاجة هنا تتحصد", "Nothing here to harvest"), false
	}
	if w.StormChance > 0 && w.Rand.Float64() < w.StormChance {
		return 0, ErrStorm, say("عاصفة ترابية! الحصاد فشل، حاول تاني", "Dust storm! Harvest failed; try again"), false
	}
	if t.Entity != Crystal && w.Tick < t.Ready {
		if !destroy {
			return 0, ErrNotRipe, say("لسه مش ناضج", "Not ripe yet"), false
		}
		w.Destroyed++
		*t = Tile{Ground: t.Ground}
		return 0, ErrNotRipe, say("النبات لسه مش ناضج… واتدمر!", "The plant was not ripe yet… and was destroyed!"), true
	}
	if w.HarvestLimit > 0 && w.Harvests >= w.HarvestLimit {
		return 0, ErrEmpty, sayf("مخزن الدرون مليان (%d شحنات بس)", "The drone's storage is full (only %d harvests)", w.HarvestLimit), false
	}
	n := 1
	if t.Entity == Crystal {
		n = t.Size
	}
	if w.Capacity > 0 && w.Cargo+n > w.Capacity {
		return 0, ErrFull, sayf("الشحنة مش هتساع: معاك %d والكريستال %d والحد %d", "Not enough cargo space: carrying %d, crystal %d, capacity %d", w.Cargo, n, w.Capacity), false
	}
	w.Harvests++
	name := t.Entity.String()
	if w.Capacity > 0 {
		w.Cargo += n
	} else {
		w.Inv[t.Entity] += n
	}
	*t = Tile{Ground: t.Ground}
	return n, nil, fmt.Sprintf("+%d %s", n, name), true
}

func (w *World) Harvest(id int) bool {
	return w.act(id, "harvest", func() (bool, string, bool) {
		if w.ScoutRoles && id != 0 {
			return false, say("درون الاستكشاف مبيعرفش يحصد", "Scout drones cannot harvest"), false
		}
		if w.DisableHarvest {
			return false, say("في المهمة دي Harvest معطّلة. استخدم TryHarvest", "Harvest is disabled for this mission. Use TryHarvest"), false
		}
		_, err, msg, changed := w.harvest(true)
		return err == nil, msg, changed
	})
}

func (w *World) TryHarvest(id int) (n int, err error) {
	_, blocked := w.actErr(id, "harvest", func() (bool, string, bool) {
		if w.ScoutRoles && id != 0 {
			err = ErrEmpty
			return false, say("درون الاستكشاف مبيعرفش يحصد", "Scout drones cannot harvest"), false
		}
		if m := w.gate("api.try"); m != "" {
			err = ErrLocked
			return false, m, false
		}
		var msg string
		var changed bool
		n, err, msg, changed = w.harvest(false)
		return err == nil, msg, changed
	})
	if blocked != nil {
		err = blocked
	}
	return n, err
}

func (w *World) Plant(id int, e Entity) bool {
	return w.act(id, "plant", func() (bool, string, bool) {
		t := w.here()
		if t.Entity != Empty || w.Z == 0 && w.buildingAt(w.X, w.Y) != nil {
			return false, say("الخانة مش فاضية", "The tile is not empty"), false
		}
		if w.Econ != nil {
			return w.colonyPlant(t, e)
		}
		switch {
		case e == Moss && t.Ground == Soil, e == Lichen && t.Ground == Ice:
		case e == Moss:
			return false, say("الطحالب محتاجة تربة (soil). استخدم Till الأول", "Moss needs soil. Use Till first"), false
		case e == Lichen:
			return false, say("الأشنة بتتزرع على الجليد (ice) بس", "Lichen only grows on ice"), false
		default:
			return false, e.String() + say(" مش حاجة بتتزرع", " cannot be planted"), false
		}
		t.Entity, t.Planted, t.Ready = e, w.Tick, w.Tick+GrowTime(e)
		return true, "", true
	})
}

func (w *World) Till(id int) bool {
	return w.act(id, "till", func() (bool, string, bool) {
		if w.Z == 0 && w.buildingAt(w.X, w.Y) != nil {
			return false, buildingMessage(econ.ErrTileNotFree), false
		}
		t := w.here()
		switch t.Ground {
		case Regolith:
			t.Ground = Soil
			return true, "", true
		case Soil:
			return true, "", false
		}
		return false, say("مينفعش تحرث ", "Cannot till ") + t.Ground.String(), false
	})
}

func (w *World) Wait(id int) {
	w.act(id, "wait", func() (bool, string, bool) { return true, "", false })
}

func (w *World) CanHarvest(id int) (r bool) {
	w.act(id, "scan", func() (bool, string, bool) {
		if m := w.gate("api.scanner"); m != "" {
			return false, m, false
		}
		t := w.here()
		r = t.Entity == Crystal || (t.Entity == Moss || t.Entity == Lichen || t.Entity == Glowcap || t.Entity == Reed) && w.Tick >= t.Ready
		if r && t.Entity == Reed && w.Tick > t.Ready+econ.ReedWindow {
			r = false // rotten
		}
		return true, "", false
	})
	return r
}

func (w *World) GetEntity(id int) (e Entity) {
	w.act(id, "scan", func() (bool, string, bool) {
		if m := w.gate("api.scanner"); m != "" {
			return false, m, false
		}
		e = w.here().Entity
		return true, "", false
	})
	return e
}

func (w *World) GetGround(id int) (g Ground) {
	w.act(id, "scan", func() (bool, string, bool) {
		if m := w.gate("api.scanner"); m != "" {
			return false, m, false
		}
		g = w.here().Ground
		return true, "", false
	})
	return g
}

func (w *World) Measure(id int) (n int) {
	w.act(id, "scan", func() (bool, string, bool) {
		if w.ScoutRoles && id == 0 {
			n = -1
			return false, say("الدرون ده مفيهوش سكانر (-1). خلّي الـ scout يقيس", "This drone has no scanner (-1). Let the scout measure"), false
		}
		if m := w.gate("api.measure"); m != "" {
			return false, m, false
		}
		if t := w.here(); t.Entity == Crystal {
			n = t.Size
		}
		return true, "", false
	})
	return n
}

// Report sends a value to base. It costs no ticks.
func (w *World) Report(v any) {
	// String methods belong to player code and may read colony state.
	s := fmt.Sprint(v)
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.done {
		select {}
	}
	w.clock()
	w.Reports = append(w.Reports, s)
	w.after(Event{Op: "report", T: w.Tick, OK: true, X: w.X, Y: w.Y, Z: w.Z, Msg: s})
}

func (w *World) Count(e Entity) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Inv[e]
}

func (w *World) Position(id int) (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if id == 0 {
		return w.X, w.Y
	}
	d := w.Drones[id]
	return d.X, d.Y
}

func (w *World) Size() (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.W, w.H
}

func (w *World) Now() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.clock()
	return w.Tick
}

// WaypointList returns a copy of the mission's waypoints, in order.
func (w *World) WaypointList() []Point {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]Point(nil), w.Waypoints...)
}

// TargetList returns a copy of the mission's targets.
func (w *World) TargetList() []Point {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]Point(nil), w.Targets...)
}
