package colony

import "orbital/colony/econ"

// BuildingKind names a building type for Build, e.g. Smelter or Silo.
type BuildingKind string

// The BuildingKinds you can Build once you own their blueprint.
const (
	Composter  BuildingKind = "composter"
	Smelter    BuildingKind = "smelter"
	Kiln       BuildingKind = "kiln"
	Refinery   BuildingKind = "refinery"
	Assembler  BuildingKind = "assembler"
	ReactorLab BuildingKind = "reactor-lab"
	Silo       BuildingKind = "silo"
	Generator  BuildingKind = "generator"
	Lab        BuildingKind = "lab"
	Collector  BuildingKind = "collector"
)

// BuildingError is the type of the Err* building errors below; compare with errors.Is.
type BuildingError = econ.BuildingError

// Building errors returned by Build and Building methods; compare with errors.Is.
var (
	ErrNoInput           = econ.ErrNoInput
	ErrNoPower           = econ.ErrNoPower
	ErrWorn              = econ.ErrWorn
	ErrBadRecipe         = econ.ErrBadRecipe
	ErrNotHere           = econ.ErrNotHere
	ErrUnknownKind       = econ.ErrUnknownKind
	ErrTileNotFree       = econ.ErrTileNotFree
	ErrBlueprintNotOwned = econ.ErrBlueprintNotOwned
)

// Building handles use the drone that obtained them for all subsequent actions.
type Building struct {
	pos   Point
	drone int
}

// BuildingStatus is a snapshot returned by Building.Status: input/output
// buffers, production progress, wear, level, recipe and power state.
type BuildingStatus struct {
	In, Out  map[Entity]float64
	Progress float64
	Wear     float64
	Level    int
	Recipe   Entity
	Powered  bool
	Enabled  bool
	Priority int
	Mode     string
}

// Build places a building of the given kind on the main drone's tile.
// The error says why it failed, e.g. ErrTileNotFree or ErrBlueprintNotOwned.
func Build(kind BuildingKind) error { return world.Build(0, string(kind)) }

// Build is Build for this drone.
func (d *Drone) Build(kind BuildingKind) error { return world.Build(d.id, string(kind)) }

// Here returns the building under the main drone, or nil. Free: no ticks.
func Here() *Building { return Main().Here() }

// Here returns the building under this drone, or nil. Free: no ticks.
func (d *Drone) Here() *Building {
	p := world.Here(d.id)
	if p == nil {
		return nil
	}
	return &Building{pos: *p, drone: d.id}
}

// Buildings returns every building in the colony, driven by the main drone.
func Buildings() []*Building { return Main().Buildings() }

// Buildings returns every building in the colony, driven by this drone.
func (d *Drone) Buildings() []*Building {
	var buildings []*Building
	for _, p := range world.Buildings() {
		buildings = append(buildings, &Building{pos: p, drone: d.id})
	}
	return buildings
}

// Kind returns what the building is, e.g. Smelter.
func (b *Building) Kind() BuildingKind { return BuildingKind(world.BuildingStatus(b.pos).Kind) }

// Pos returns the building's tile.
func (b *Building) Pos() Point { return b.pos }

// SetRecipe chooses the good the building produces. The drone must stand on it.
func (b *Building) SetRecipe(out Entity) error { return world.BuildingSetRecipe(b.drone, b.pos, out) }

// Put moves up to n of good e from the colony into the building and returns how
// many moved. The drone must stand on it.
func (b *Building) Put(e Entity, n int) (int, error) { return world.BuildingPut(b.drone, b.pos, e, n) }

// Take moves up to n of good e out of the building and returns how many moved.
// The drone must stand on it.
func (b *Building) Take(e Entity, n int) (int, error) {
	return world.BuildingTake(b.drone, b.pos, e, n)
}

// Repair reduces the building's wear using plates from the colony (ErrNoInput
// without any). The drone must stand on it.
func (b *Building) Repair() error { return world.BuildingRepair(b.drone, b.pos) }

// Status returns a snapshot of the building: buffers, progress, wear and power.
func (b *Building) Status() BuildingStatus {
	v := world.BuildingStatus(b.pos)
	s := BuildingStatus{In: map[Entity]float64{}, Out: map[Entity]float64{}, Progress: v.Progress, Wear: v.Wear, Level: v.Level, Powered: v.Powered, Enabled: v.Enabled, Priority: v.Priority, Mode: v.Mode}
	for e := Moss; e <= Core; e++ {
		if n, ok := v.In[e.String()]; ok {
			s.In[e] = n
		}
		if n, ok := v.Out[e.String()]; ok {
			s.Out[e] = n
		}
		if v.Recipe == e.String() {
			s.Recipe = e
		}
	}
	return s
}

// Melt turns the surface ice under the main drone into water. It reports
// whether water was collected; the tile regrows ice later.
func Melt() bool { return world.Melt(0) }

// Sift collects silica from the surface regolith under the main drone.
// Each tile's stock refills over time.
func Sift() bool { return world.Sift(0) }

// Melt is Melt for this drone.
func (d *Drone) Melt() bool { return world.Melt(d.id) }

// Sift is Sift for this drone.
func (d *Drone) Sift() bool { return world.Sift(d.id) }
