package colony

// Dig opens the adjacent underground Stone tile and collects its ore.
func Dig(dir Direction) bool { return world.Dig(0, dir) }

// DigDown opens a shaft to the next layer and collects ore directly below.
func DigDown() bool { return world.DigDown(0) }

// Descend moves down through the shaft on the current tile.
func Descend() bool { return world.Descend(0) }

// Ascend moves up through a shaft in the tile above.
func Ascend() bool { return world.Ascend(0) }

// Layer returns the drone's depth (0 is the surface), without spending a tick.
func Layer() int { return world.Layer(0) }

// Depth returns the world's number of layers, without spending a tick.
func Depth() int { return world.LayerDepth() }

// Dig is Dig for this drone.
func (d *Drone) Dig(dir Direction) bool { return world.Dig(d.id, dir) }

// DigDown is DigDown for this drone.
func (d *Drone) DigDown() bool { return world.DigDown(d.id) }

// Descend is Descend for this drone.
func (d *Drone) Descend() bool { return world.Descend(d.id) }

// Ascend is Ascend for this drone.
func (d *Drone) Ascend() bool { return world.Ascend(d.id) }

// Layer is Layer for this drone.
func (d *Drone) Layer() int { return world.Layer(d.id) }

// Depth returns the world's number of layers, without spending a tick.
func (d *Drone) Depth() int { return world.LayerDepth() }

// Grid returns the ground of every tile on this drone's layer as grid[x][y].
func (d *Drone) Grid() [][]Ground { return world.DroneGrid(d.id) }

// LayerGrid returns the ground of every tile on layer z as grid[x][y],
// from the survey map. Free: no ticks.
func LayerGrid(z int) [][]Ground { return world.LayerGrid(z) }
