package colony

import (
	"context"
	"errors"

	"orbital/colony/internal/sim"
)

// FaultError describes a failed diagnostic: its Index, a short Code and a
// Message. Read the code from any wrapped error with FaultCode.
type FaultError = sim.FaultError

// Question is one exam question: a Prompt and its Choices (answer by index).
type Question = sim.Question

// Dataset accessors return copies, so player code cannot change mission inputs.
func Callsigns() []string { return world.CallsignList() }

// Readings returns a copy of the mission's text readings.
func Readings() []string { return world.ReadingList() }

// Orders returns a copy of the mission's commands.
func Orders() []string { return world.OrderList() }

// Diagnose returns the fault found by diagnostic i: nil when healthy, a
// *FaultError when faulty, or an error for an index out of range.
func Diagnose(i int) error { return world.Diagnose(i) }

// Dock connects the drone to the drill rig; Drill needs it.
func Dock() { world.Dock() }

// Undock disconnects the drone from the drill rig (pair it with Dock, often via defer).
func Undock() { world.Undock() }

// Drill panics on a configured drill fault or if the drone is not docked.
func Drill() { world.Drill() }

// Exam returns a copy of the exam questions; answer them with Answer.
func Exam() []Question { return world.Exam() }

// Answer records a zero-based choice and reports whether the indices are valid.
func Answer(i, choice int) bool { return world.Answer(i, choice) }

// Invoice returns the bill total in cents and the number of people sharing it.
func Invoice() (totalCents, people int) { return world.Invoice() }

// Callsign returns one callsign, or empty text for an invalid index.
func Callsign(i int) string { return world.Callsign(i) }

// CargoAmount returns the amount for kind 0..2, or zero for an invalid kind.
func CargoAmount(kind int) int { return world.CargoAmount(kind) }

// ReadingCount returns how many numeric sensor readings there are (see Reading).
func ReadingCount() int { return world.ReadingCount() }

// Reading returns a numeric sensor reading, or zero for an invalid index.
func Reading(i int) int { return world.Reading(i) }

// DiagnosticCount returns how many diagnostics Diagnose accepts.
func DiagnosticCount() int { return world.DiagnosticCount() }

// DrillAttempts returns how many drill attempts the mission asks for.
func DrillAttempts() int { return world.DrillAttempts() }

// OrderCount returns how many commands Order accepts.
func OrderCount() int { return world.OrderCount() }

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

// TraceCommand records the supplied context at a named command boundary.
func TraceCommand(ctx context.Context, stage string) { world.TraceCommand(ctx, stage) }

// RescueRequest is a pending rescue: its ID and Priority (higher goes first).
type RescueRequest = sim.RescueRequest

// RescueArrivals returns a copy of the current arrivals, initially and after
// every successful rescue. Read each wave once to avoid enqueueing duplicates.
func RescueArrivals() []RescueRequest { return world.RescueArrivals() }

// Rescue serves the highest-priority pending request (smaller ID breaks ties).
// A rejected request permanently fails the mission, and releases no new wave.
func Rescue(id int) bool { return world.Rescue(id) }
