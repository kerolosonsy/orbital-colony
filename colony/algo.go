package colony

import "orbital/colony/internal/sim"

// Swap exchanges what is on the drone's tile with the neighbouring tile in
// direction dir. The drone stays where it is. Costs one tick.
func Swap(dir Direction) bool { return world.Swap(0, dir) }

// Grid returns the ground of every tile as grid[x][y], from the orbital
// map. It is free: no ticks.
func Grid() [][]Ground { return world.GroundGrid() }

// Kg is a weight in kilograms.
type Kg = sim.Kg

// Depths returns the drill depths of the mission (in metres).
func Depths() []int { return append([]int(nil), world.Data.Depths...) }

// Temperatures returns the temperature readings of the mission.
func Temperatures() []float64 { return append([]float64(nil), world.Data.Temps...) }

// Weights returns the cargo weights of the mission.
func Weights() []Kg { return append([]Kg(nil), world.Data.Weights...) }
