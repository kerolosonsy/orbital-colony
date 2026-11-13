package colony

import "orbital/colony/internal/sim"

// SystemSnapshot contains the authoritative values for HTTP inventory and metrics.
type SystemSnapshot = sim.SystemSnapshot

// SystemStats returns a consistent copy of the colony's lifetime counters and inventory.
func SystemStats() SystemSnapshot { return world.SystemStats() }

// Inventory returns a copy of inventory, including fractional goods.
func Inventory() map[string]float64 { return world.SystemStats().Inventory }
