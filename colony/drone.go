package colony

// Drone is an extra drone launched with Spawn. Each drone can be driven from
// its own goroutine. A single drone does one thing at a time: if two
// goroutines drive the same drone, its actions simply queue up.
type Drone struct {
	id int
}

// Spawn launches a new drone from the base. It returns nil when the hangar
// is empty.
func Spawn() *Drone {
	id := world.Spawn()
	if id < 0 {
		return nil
	}
	return &Drone{id: id}
}

// ID returns the drone's number. The first drone (the package-level
// functions) is drone 0.
func (d *Drone) ID() int { return d.id }

// Move is Move for this drone.
func (d *Drone) Move(dir Direction) bool { return world.Move(d.id, dir) }

// Harvest is Harvest for this drone.
func (d *Drone) Harvest() bool { return world.Harvest(d.id) }

// TryHarvest is TryHarvest for this drone.
func (d *Drone) TryHarvest() (int, error) {
	return world.TryHarvest(d.id)
}

// Plant is Plant for this drone.
func (d *Drone) Plant(e Entity) bool { return world.Plant(d.id, e) }

// Till is Till for this drone.
func (d *Drone) Till() bool { return world.Till(d.id) }

// Wait is Wait for this drone.
func (d *Drone) Wait() { world.Wait(d.id) }

// CanHarvest is CanHarvest for this drone.
func (d *Drone) CanHarvest() bool { return world.CanHarvest(d.id) }

// GetEntity is GetEntity for this drone.
func (d *Drone) GetEntity() Entity { return world.GetEntity(d.id) }

// GetGround is GetGround for this drone.
func (d *Drone) GetGround() Ground { return world.GetGround(d.id) }

// Measure is Measure for this drone.
func (d *Drone) Measure() int { return world.Measure(d.id) }

// Unload is Unload for this drone.
func (d *Drone) Unload() bool { return world.Unload(d.id) }

// Position returns this drone's x and y. Free: no ticks.
func (d *Drone) Position() (x, y int) { return world.Position(d.id) }

// StormWarning receives a value when a dust storm is about to hit.
func StormWarning() <-chan struct{} { return world.StormWarning() }

// StormClear closes when the warned storm clears. Fetch it after each warning.
func StormClear() <-chan struct{} { return world.StormClear() }

// Main returns a handle to drone 0, the one the package-level functions
// control, so it can be passed around like any other *Drone.
func Main() *Drone { return &Drone{id: 0} }
