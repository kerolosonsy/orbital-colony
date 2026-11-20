package colony

import "orbital/colony/internal/sim"

// ErrFault means a command needs Repair before it can succeed.
var ErrFault = sim.ErrFault

// DroneError identifies the drone; errors.Is recognizes ErrFault and ErrStorm.
type DroneError = sim.DroneError

// Fault is nil for a healthy drone, ErrFault when broken, or ErrStorm while stunned.
func Fault() error { return world.Fault(0) }

// Repair spends one command repairing a fault, at any location, without materials.
// A storm-stunned drone must wait until its stun ends before repairing.
func Repair() bool { return world.Repair(0) }

// Fault is Fault for this drone.
func (d *Drone) Fault() error { return world.Fault(d.id) }

// Repair is Repair for this drone.
func (d *Drone) Repair() bool { return world.Repair(d.id) }

// Rush returns the active order's good and remaining amount, or ok=false.
// Exports automatically fill it at twice the current market price.
func Rush() (good Entity, amount int, ok bool) { return world.Rush() }

// Shutdown registers for a shutdown signal. After it closes, return to base
// and unload; the colony saves and exits after at most two seconds.
func Shutdown() <-chan struct{} { return world.Shutdown() }
