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

func (d *Drone) Move(dir Direction) bool { return world.Move(d.id, dir) }
func (d *Drone) Harvest() bool           { return world.Harvest(d.id) }
func (d *Drone) TryHarvest() (int, error) {
	return world.TryHarvest(d.id)
}
func (d *Drone) Plant(e Entity) bool  { return world.Plant(d.id, e) }
func (d *Drone) Till() bool           { return world.Till(d.id) }
func (d *Drone) Wait()                { world.Wait(d.id) }
func (d *Drone) CanHarvest() bool     { return world.CanHarvest(d.id) }
func (d *Drone) GetEntity() Entity    { return world.GetEntity(d.id) }
func (d *Drone) GetGround() Ground    { return world.GetGround(d.id) }
func (d *Drone) Measure() int         { return world.Measure(d.id) }
func (d *Drone) Unload() bool         { return world.Unload(d.id) }
func (d *Drone) Position() (x, y int) { return world.Position(d.id) }

// StormWarning receives a value when a dust storm is about to hit.
func StormWarning() <-chan struct{} {
	if world.Storm == nil {
		return nil
	}
	return world.Storm.Warning
}

// StormClear is closed when the storm is over.
func StormClear() <-chan struct{} {
	if world.Storm == nil {
		return nil
	}
	return world.Storm.Clear
}

// Main returns a handle to drone 0, the one the package-level functions
// control, so it can be passed around like any other *Drone.
func Main() *Drone { return &Drone{id: 0} }
