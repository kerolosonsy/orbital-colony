// Package colony controls your drone on planet Kepler-G.
//
// Every drone action (Move, Harvest, Plant, Till, Wait and the sensors
// CanHarvest, GetEntity, GetGround, Measure) costs one tick of energy.
// Colony information (Count, Position, WorldSize, Tick, Report) is free.
package colony

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"orbital/colony/econ"
	"orbital/colony/internal/sim"
	"orbital/colony/levels"
)

// Direction is where the drone moves. North is up, East is right.
type Direction = sim.Direction

// The four Directions for Move: North (up), East (right), South, West.
const (
	North = sim.North
	East  = sim.East
	South = sim.South
	West  = sim.West
)

// Entity is what sits on a tile.
type Entity = sim.Entity

// Entities: crops and resources on tiles, and the goods made from them.
const (
	Empty      = sim.Empty
	Moss       = sim.Moss
	Lichen     = sim.Lichen
	Crystal    = sim.Crystal
	Beacon     = sim.Beacon
	Glowcap    = sim.Glowcap
	Reed       = sim.Reed
	Iron       = sim.Iron
	Geode      = sim.Geode
	Water      = sim.Water
	Silica     = sim.Silica
	Copper     = sim.Copper
	Sulfur     = sim.Sulfur
	Titanium   = sim.Titanium
	Helium     = sim.Helium
	Uranium    = sim.Uranium
	Biomass    = sim.Biomass
	Plate      = sim.Plate
	Wire       = sim.Wire
	Glass      = sim.Glass
	Fuel       = sim.Fuel
	Fertilizer = sim.Fertilizer
	Circuit    = sim.Circuit
	Alloy      = sim.Alloy
	Battery    = sim.Battery
	Coolant    = sim.Coolant
	Processor  = sim.Processor
	Frame      = sim.Frame
	Cell       = sim.Cell
	Core       = sim.Core
)

// Ground is the surface of a tile.
type Ground = sim.Ground

// Grounds a tile can have: Regolith (Till it into Soil, or Sift it), Soil
// (grows Moss), Ice (grows Lichen, or Melt it), Rock, underground Stone (Dig
// it) and the dug-out Tunnel and Shaft.
const (
	Regolith = sim.Regolith
	Soil     = sim.Soil
	Ice      = sim.Ice
	Rock     = sim.Rock
	Stone    = sim.Stone
	Tunnel   = sim.Tunnel
	Shaft    = sim.Shaft
)

// Point is a position on the map. (0, 0) is the bottom-left corner.
type Point = sim.Point

// Errors returned by TryHarvest. Compare with errors.Is.
var (
	ErrNotRipe = sim.ErrNotRipe
	ErrEmpty   = sim.ErrEmpty
	ErrStorm   = sim.ErrStorm
	ErrFull    = sim.ErrFull
)

var world *sim.World

func init() {
	var out io.Writer = io.Discard
	trace := os.Getenv("ORBITAL_TRACE")
	if trace != "" {
		if f, err := os.Create(trace); err == nil {
			out = f
		}
	}
	seed := time.Now().UnixNano()
	if s, err := strconv.ParseInt(os.Getenv("ORBITAL_SEED"), 10, 64); err == nil {
		seed = s
	}
	id := levelID()
	if id == "colony" {
		startColony(out, seed, trace == "" && os.Getenv("ORBITAL_EVENTS_FD") == "")
		return
	}
	world = levels.Start(id, out, seed)
	if trace == "" {
		// Run from a terminal or VS Code: narrate the mission on stderr.
		world.Console = os.Stderr
		levels.Narrate(world)
	}
	world.Launch()
}

// Move moves the drone one tile. It returns false if a wall is in the way.
func Move(d Direction) bool { return world.Move(0, d) }

// Harvest collects what is on the current tile and returns true on success.
// Harvesting an unripe plant destroys it.
func Harvest() bool { return world.Harvest(0) }

// TryHarvest is the careful version of Harvest. It returns how much was
// collected, or an error (ErrNotRipe, ErrEmpty, ErrStorm). It never destroys.
func TryHarvest() (int, error) { return world.TryHarvest(0) }

// Plant plants e on the current tile. Moss needs Soil, Lichen needs Ice.
func Plant(e Entity) bool { return world.Plant(0, e) }

// Till turns Regolith into Soil.
func Till() bool { return world.Till(0) }

// Wait does nothing for one tick.
func Wait() { world.Wait(0) }

// CanHarvest reports whether the current tile is ready to harvest.
func CanHarvest() bool { return world.CanHarvest(0) }

// GetEntity returns what is on the current tile.
func GetEntity() Entity { return world.GetEntity(0) }

// GetGround returns the ground type of the current tile.
func GetGround() Ground { return world.GetGround(0) }

// Measure returns the size of a crystal on the current tile, or 0.
func Measure() int { return world.Measure(0) }

// Count returns how many of e the colony has stored.
func Count(e Entity) int { return world.Count(e) }

// Position returns the drone's x and y. (0, 0) is the bottom-left corner.
func Position() (x, y int) { return world.Position(0) }

// WorldSize returns the width and height of the map.
func WorldSize() (width, height int) { return world.Size() }

// Tick returns how many ticks have been spent so far.
func Tick() int { return world.Now() }

// Report transmits a value to the base.
func Report(v any) { world.Report(v) }

// Waypoints returns the points the drone must visit, in order.
func Waypoints() []Point { return world.WaypointList() }

// Targets returns the positions of this mission's target tiles.
func Targets() []Point { return world.TargetList() }

// levelID picks the mission: ORBITAL_LEVEL when the game runs the program,
// otherwise the folder name (workspace/<level>/), which is also the binary
// name go run uses.
func levelID() string {
	if id := os.Getenv("ORBITAL_LEVEL"); id != "" {
		return id
	}
	candidates := []string{filepath.Base(os.Args[0])}
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Base(wd))
	}
	for _, c := range candidates {
		if levels.Get(c) != nil {
			return c
		}
	}
	return ""
}

// startColony loads the open colony from the save the game passes in (or a
// fresh one) and wires the live event stream and the purchase channel.
func startColony(out io.Writer, seed int64, terminal bool) {
	st := econ.New(seed)
	if p := os.Getenv("ORBITAL_COLONY"); p != "" {
		if b, err := os.ReadFile(p); err == nil {
			loaded := &econ.State{}
			if json.Unmarshal(b, loaded) == nil && loaded.Size > 0 {
				st = loaded
			}
		}
	}
	if os.Getenv("ORBITAL_EVENTS_FD") != "" {
		out = os.NewFile(3, "events")
	}
	world = sim.NewColony(out, st, seed)
	if terminal {
		world.Console = os.Stderr
	}
	world.Start()
	world.Launch()
	var cmds io.Reader
	if os.Getenv("ORBITAL_CMD_FD") != "" {
		cmds = os.NewFile(4, "commands")
	}
	world.RunColony(cmds)
}

// Export sells n of e on the orbital market (n <= 0 sells all of it) and
// returns the credits earned. Prices drop a little with every unit you sell
// and recover over time.
func Export(e Entity, n int) float64 { return world.Export(0, e, n) }

// Price is what one unit of e sells for right now. It needs the market
// screen upgrade; before that it returns 0.
func Price(e Entity) float64 { return world.Price(e) }

// Credits is the colony's money.
func Credits() float64 { return world.Credits() }

// Export is the per-drone version of the package-level Export.
func (d *Drone) Export(e Entity, n int) float64 { return world.Export(d.id, e, n) }
