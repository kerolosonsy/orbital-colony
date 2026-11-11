package colony

import "orbital/colony/internal/sim"

type FaultError = sim.FaultError
type Question = sim.Question

// Dataset accessors return copies, so player code cannot change mission inputs.
func Callsigns() []string  { return world.CallsignList() }
func Readings() []string   { return world.ReadingList() }
func Orders() []string     { return world.OrderList() }
func Diagnose(i int) error { return world.Diagnose(i) }
func Dock()                { world.Dock() }
func Undock()              { world.Undock() }

// Drill panics on a configured drill fault or if the drone is not docked.
func Drill()           { world.Drill() }
func Exam() []Question { return world.Exam() }

// Answer records a zero-based choice and reports whether the indices are valid.
func Answer(i, choice int) bool         { return world.Answer(i, choice) }
func Invoice() (totalCents, people int) { return world.Invoice() }
