package sim

import (
	"errors"
	"maps"
	"math"

	"orbital/colony/econ"
)

type BuildingView struct {
	Powered  bool               `json:"powered"`
	Enabled  bool               `json:"enabled"`
	Priority int                `json:"priority"`
	Mode     string             `json:"mode,omitempty"`
	Kind     string             `json:"kind"`
	Recipe   string             `json:"recipe"`
	In       map[string]float64 `json:"in"`
	Out      map[string]float64 `json:"out"`
	Progress float64            `json:"progress"`
	Wear     float64            `json:"wear"`
	Level    int                `json:"level"`
}

func (w *World) buildingAt(x, y int) *econ.Building {
	if w.Econ != nil {
		for _, b := range w.Econ.Buildings {
			if b.X == x && b.Y == y {
				return b
			}
		}
	}
	return nil
}

func buildingView(b *econ.Building) BuildingView {
	return BuildingView{Powered: b.Powered, Enabled: b.Enabled, Priority: b.Priority, Mode: b.Mode, Kind: b.Kind, Recipe: b.Recipe, In: maps.Clone(b.In), Out: maps.Clone(b.Out), Progress: b.ProgressFraction(), Wear: b.Wear, Level: b.Level}
}

func (w *World) emitBuilding(b *econ.Building, tick int) {
	v := buildingView(b)
	w.emit(Event{Op: "building", T: tick, OK: true, X: b.X, Y: b.Y, Building: &v})
}

func (w *World) Here(id int) *Point {
	w.mu.Lock()
	defer w.mu.Unlock()
	x, y, z := w.X, w.Y, w.Z
	if id != 0 {
		d := w.Drones[id]
		x, y, z = d.X, d.Y, d.Z
	}
	if z != 0 || w.buildingAt(x, y) == nil {
		return nil
	}
	return &Point{x, y}
}

func (w *World) Buildings() []Point {
	w.mu.Lock()
	defer w.mu.Unlock()
	var positions []Point
	if w.Econ != nil {
		for _, b := range w.Econ.Buildings {
			positions = append(positions, Point{b.X, b.Y})
		}
	}
	return positions
}

func (w *World) BuildingStatus(p Point) BuildingView {
	w.mu.Lock()
	defer w.mu.Unlock()
	if b := w.buildingAt(p.X, p.Y); b != nil {
		return buildingView(b)
	}
	return BuildingView{}
}

func buildingMessage(err error) string {
	var e econ.BuildingError
	if errors.As(err, &e) {
		ar, en := e.Messages()
		return say(ar, en)
	}
	return err.Error()
}

func (w *World) placeBuilding(kind string) error {
	cost := econ.PlacementCost(kind)
	if cost == nil {
		return econ.ErrUnknownKind
	}
	if w.Econ == nil {
		return ErrLocked
	}
	cost = w.Econ.BuildingPlacementCost(kind)
	if !w.Econ.Has(kind) {
		return econ.ErrBlueprintNotOwned
	}
	t := w.here()
	if w.Z != 0 || t.Ground != Regolith && t.Ground != Soil || t.Entity != Empty || w.buildingAt(w.X, w.Y) != nil {
		return econ.ErrTileNotFree
	}
	if w.Econ.Credits < cost.Credits {
		return econ.ErrNoInput
	}
	for good, n := range cost.Goods {
		if w.Econ.Inv[good] < n {
			return econ.ErrNoInput
		}
	}
	w.Econ.Credits -= cost.Credits
	for good, n := range cost.Goods {
		w.Econ.Inv[good] -= n
	}
	b := &econ.Building{Kind: kind, X: w.X, Y: w.Y, Level: w.Econ.Level(kind), Enabled: true, In: map[string]float64{}, Out: map[string]float64{}}
	if kind == "generator" {
		b.Mode = "solar"
	}
	if kind == "collector" {
		b.Recipe = "helium"
	}
	w.Econ.Buildings = append(w.Econ.Buildings, b)
	w.Capacity = w.Econ.CargoCapacity()
	w.refreshPower(b)
	w.syncInv()
	w.emitBuilding(b, w.Tick)
	return nil
}

func (w *World) Build(id int, kind string) (err error) {
	_, blocked := w.actErr(id, "build", func() (bool, string, bool) {
		err = w.placeBuilding(kind)
		if err != nil {
			return false, buildingMessage(err), false
		}
		return true, "", true
	})
	if blocked != nil {
		err = blocked
	}
	return err
}

func (w *World) buildingAction(id int, p Point, op string, f func(*econ.Building) error) (err error) {
	_, blocked := w.actErr(id, op, func() (bool, string, bool) {
		b := w.buildingAt(p.X, p.Y)
		if b == nil {
			err = econ.ErrNotHere
		} else {
			err = f(b)
		}
		if err != nil {
			return false, buildingMessage(err), false
		}
		w.syncInv()
		return true, "", false
	})
	if blocked != nil {
		err = blocked
	}
	return err
}

func (w *World) BuildingSetRecipe(id int, p Point, out Entity) error {
	return w.buildingAction(id, p, "recipe", func(b *econ.Building) error {
		before := b.Recipe
		if err := b.SetRecipe(out.String()); err != nil {
			return err
		}
		if before != b.Recipe {
			w.emitBuilding(b, w.Tick)
		}
		return nil
	})
}

func (w *World) standingOn(b *econ.Building) bool {
	return w.Z == 0 && w.X == b.X && w.Y == b.Y
}

func (w *World) BuildingPut(id int, p Point, e Entity, n int) (int, error) {
	accepted := 0
	err := w.buildingAction(id, p, "put", func(b *econ.Building) error {
		if !w.standingOn(b) {
			return econ.ErrNotHere
		}
		if n <= 0 {
			return nil
		}
		amount := math.Min(float64(n), math.Floor(w.Econ.Inv[e.String()]+w.actingDrone.Hold[e.String()]))
		if amount <= 0 {
			return econ.ErrNoInput
		}
		// The player API moves whole goods; fractional capacity stays available for later levels.
		w.Econ.ConfigureBuilding(b)
		space := math.Floor(w.Econ.BuildingCapacity(b.Kind, b.Level) - b.In[e.String()])
		if space <= 0 {
			return econ.ErrFull
		}
		moved, err := b.Put(e.String(), math.Min(amount, space))
		if err == nil && moved > 0 {
			accepted = int(moved)
			cargo := math.Min(moved, w.actingDrone.Hold[e.String()])
			w.actingDrone.Hold[e.String()] -= cargo
			w.Cargo = cargoUsed(w.actingDrone)
			w.Econ.Inv[e.String()] -= moved - cargo
			w.emitBuilding(b, w.Tick)
		}
		return err
	})
	return accepted, err
}

func (w *World) BuildingTake(id int, p Point, e Entity, n int) (int, error) {
	taken := 0
	err := w.buildingAction(id, p, "take", func(b *econ.Building) error {
		if !w.standingOn(b) {
			return econ.ErrNotHere
		}
		if n <= 0 {
			return nil
		}
		amount := math.Min(float64(n), math.Floor(b.Out[e.String()]))
		if amount <= 0 {
			return econ.ErrNoInput
		}
		if err := w.collectGood(e, amount); err != nil {
			return err
		}
		moved, err := b.Take(e.String(), amount)
		if err == nil {
			taken = int(moved)
			w.emitBuilding(b, w.Tick)
		}
		return err
	})
	return taken, err
}

func (w *World) BuildingRepair(id int, p Point) error {
	return w.buildingAction(id, p, "repair", func(b *econ.Building) error {
		if b.Wear <= 0 {
			return nil
		}
		used := b.Repair(w.Econ.Inv["plate"])
		if used == 0 {
			return econ.ErrNoInput
		}
		w.Econ.Inv["plate"] -= used
		w.emitBuilding(b, w.Tick)
		return nil
	})
}

func (w *World) advanceProduction(target int) {
	for tick := w.lastProductionTick + 1; tick <= target; tick++ {
		w.archiveTick(tick)
		w.advanceColonyEvents(tick)
		w.productionStep(tick)
	}
	w.lastProductionTick = target
	w.restoreIce(target)
}
