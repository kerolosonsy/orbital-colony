package colony

import "orbital/colony/econ"

type BuildingKind string

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
)

type BuildingError = econ.BuildingError

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

type BuildingStatus struct {
	In, Out  map[Entity]float64
	Progress float64
	Wear     float64
	Level    int
	Recipe   Entity
	Powered  bool
}

func Build(kind BuildingKind) error            { return world.Build(0, string(kind)) }
func (d *Drone) Build(kind BuildingKind) error { return world.Build(d.id, string(kind)) }

func Here() *Building { return Main().Here() }
func (d *Drone) Here() *Building {
	p := world.Here(d.id)
	if p == nil {
		return nil
	}
	return &Building{pos: *p, drone: d.id}
}

func Buildings() []*Building { return Main().Buildings() }
func (d *Drone) Buildings() []*Building {
	var buildings []*Building
	for _, p := range world.Buildings() {
		buildings = append(buildings, &Building{pos: p, drone: d.id})
	}
	return buildings
}

func (b *Building) Kind() BuildingKind               { return BuildingKind(world.BuildingStatus(b.pos).Kind) }
func (b *Building) Pos() Point                       { return b.pos }
func (b *Building) SetRecipe(out Entity) error       { return world.BuildingSetRecipe(b.drone, b.pos, out) }
func (b *Building) Put(e Entity, n int) (int, error) { return world.BuildingPut(b.drone, b.pos, e, n) }
func (b *Building) Take(e Entity, n int) (int, error) {
	return world.BuildingTake(b.drone, b.pos, e, n)
}
func (b *Building) Repair() error { return world.BuildingRepair(b.drone, b.pos) }

func (b *Building) Status() BuildingStatus {
	v := world.BuildingStatus(b.pos)
	s := BuildingStatus{In: map[Entity]float64{}, Out: map[Entity]float64{}, Progress: v.Progress, Wear: v.Wear, Level: v.Level, Powered: true}
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

func Melt() bool            { return world.Melt(0) }
func Sift() bool            { return world.Sift(0) }
func (d *Drone) Melt() bool { return world.Melt(d.id) }
func (d *Drone) Sift() bool { return world.Sift(d.id) }
