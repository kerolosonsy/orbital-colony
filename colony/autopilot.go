package colony

import "orbital/colony/internal/sim"

// Unload moves the drone's cargo into storage. It only works at the base.
func Unload() bool { return world.Unload(0) }

// Nearest scans for the closest remaining target crystal. It returns nil
// when there is none left. Costs one tick.
func Nearest() *Point { return world.Nearest(0) }

// Signal is a message from the colony radio. Use a type switch to find out
// whether it is a Distress, a Supply or a Noise.
type Signal = sim.Signal

// The three kinds of Signal.
type (
	Distress = sim.Distress
	Supply   = sim.Supply
	Noise    = sim.Noise
)

// Signals returns every signal the radio picked up.
func Signals() []Signal { return world.SignalList() }

// Tile is what the autopilot tells your Planner about one tile.
type Tile struct {
	X, Y   int
	Ground Ground
	Entity Entity
	Ripe   bool
}

// Action is one thing the autopilot can do on a tile.
type Action int

const (
	DoHarvest Action = iota + 1
	DoTill
	PlantMoss
	PlantLichen
)

// Planner decides what the autopilot does on each tile.
type Planner interface {
	Plan(t Tile) []Action
}

// Autopilot flies to (0, 0) and sweeps the whole map row by row, asking p
// what to do on every tile. Reading tiles is free; actions and moves cost
// ticks as usual.
func Autopilot(p Planner) {
	w, h := world.Size()
	x, y := world.Position(0)
	for ; x > 0; x-- {
		world.Move(0, West)
	}
	for ; y > 0; y-- {
		world.Move(0, South)
	}
	for y := 0; y < h; y++ {
		for i := 0; i < w; i++ {
			x, dir := i, East
			if y%2 == 1 {
				x, dir = w-1-i, West
			}
			visit(p, x, y)
			if i < w-1 {
				world.Move(0, dir)
			}
		}
		if y < h-1 {
			world.Move(0, North)
		}
	}
}

func visit(p Planner, x, y int) {
	g, e, ripe := world.Peek(x, y)
	for _, a := range p.Plan(Tile{X: x, Y: y, Ground: g, Entity: e, Ripe: ripe}) {
		switch a {
		case DoHarvest:
			world.Harvest(0)
		case DoTill:
			world.Till(0)
		case PlantMoss:
			world.Plant(0, Moss)
		case PlantLichen:
			world.Plant(0, Lichen)
		}
	}
}
