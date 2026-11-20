package colony

// Power returns the colony grid's generated supply and the demand of its
// powered buildings. Free: no ticks.
func Power() (supply, demand float64) { return world.Power() }

// Daylight reports whether it is currently day on the colony clock. Free: no ticks.
func Daylight() bool { return world.Daylight() }

// Storage returns, per storage group, how much is stored and the capacity. Free: no ticks.
func Storage() (used, capacity map[string]int) { return world.Storage() }

// Carry returns what the main drone is holding, by good. Free: no ticks.
func Carry() map[Entity]int { return world.Carry(0) }

// CargoFree returns how many more units the main drone can carry. Free: no ticks.
func CargoFree() int { return world.CargoFree(0) }

// BatteryLevel returns the main drone's battery charge. Free: no ticks.
func BatteryLevel() int { return world.Battery(0) }

// Charge refills the main drone's battery at the base (0, 0) using one fuel
// (needs a refinery). It reports whether charging worked.
func Charge() bool { return world.Charge(0) }

// Carry is Carry for this drone.
func (d *Drone) Carry() map[Entity]int { return world.Carry(d.id) }

// CargoFree is CargoFree for this drone.
func (d *Drone) CargoFree() int { return world.CargoFree(d.id) }

// Battery returns this drone's battery charge (BatteryLevel for this drone).
func (d *Drone) Battery() int { return world.Battery(d.id) }

// Charge is Charge for this drone.
func (d *Drone) Charge() bool { return world.Charge(d.id) }

// SetMode switches the building's operating mode. The drone must stand on it.
func (b *Building) SetMode(mode string) error { return world.BuildingSetMode(b.drone, b.pos, mode) }

// SetPriority sets the building's power priority when supply runs short.
// The drone must stand on it.
func (b *Building) SetPriority(n int) error { return world.BuildingSetPriority(b.drone, b.pos, n) }

// Enable turns the building on or off. A disabled building draws no power.
// The drone must stand on it.
func (b *Building) Enable(on bool) error { return world.BuildingEnable(b.drone, b.pos, on) }

// TryMelt is Melt with a reason: it returns the units collected and an error
// such as ErrNotHere when the drone is not on surface ice.
func TryMelt() (int, error) { return world.TryMelt(0) }

// TrySift is Sift with a reason: it returns the units collected and an error
// such as ErrNotHere (not on regolith) or ErrNoInput (stock depleted).
func TrySift() (int, error) { return world.TrySift(0) }

// TryMelt is TryMelt for this drone.
func (d *Drone) TryMelt() (int, error) { return world.TryMelt(d.id) }

// TrySift is TrySift for this drone.
func (d *Drone) TrySift() (int, error) { return world.TrySift(d.id) }

// TryDig is Dig with a reason: it returns nil on success or the error that stopped it.
func TryDig(dir Direction) error { return world.TryDig(0, dir) }

// TryDigDown is DigDown with a reason: it returns nil on success or the error that stopped it.
func TryDigDown() error { return world.TryDigDown(0) }

// TryDig is TryDig for this drone.
func (d *Drone) TryDig(dir Direction) error { return world.TryDig(d.id, dir) }

// TryDigDown is TryDigDown for this drone.
func (d *Drone) TryDigDown() error { return world.TryDigDown(d.id) }
