package colony

import (
	"errors"

	"orbital/colony/internal/sim"
)

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

// Callsign returns one callsign, or empty text for an invalid index.
func Callsign(i int) string { return world.Callsign(i) }

// CargoAmount returns the amount for kind 0..2, or zero for an invalid kind.
func CargoAmount(kind int) int { return world.CargoAmount(kind) }
func ReadingCount() int        { return world.ReadingCount() }

// Reading returns a numeric sensor reading, or zero for an invalid index.
func Reading(i int) int    { return world.Reading(i) }
func DiagnosticCount() int { return world.DiagnosticCount() }
func DrillAttempts() int   { return world.DrillAttempts() }
func OrderCount() int      { return world.OrderCount() }

// Order returns one command, or empty text for an invalid index.
func Order(i int) string { return world.Order(i) }

// FaultCode follows wrapped and joined causes without exposing pointer syntax
// to the early error lessons. It returns empty text when no fault is present.
func FaultCode(err error) string {
	var fault *FaultError
	if errors.As(err, &fault) {
		return fault.Code
	}
	return ""
}
