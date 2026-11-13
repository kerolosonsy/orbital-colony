package sim

import (
	"maps"
	"math"
	"slices"

	"orbital/colony/econ"
)

func (w *World) cargoActive() bool {
	return w.Econ != nil && (len(w.Econ.Buildings) > 0 || w.Econ.Won("c13-reed-clock"))
}

func cargoAmount(d *Drone) float64 {
	total := 0.0
	for _, good := range slices.Sorted(maps.Keys(d.Hold)) {
		total += d.Hold[good]
	}
	return total
}

func cargoUsed(d *Drone) int { return int(math.Ceil(cargoAmount(d) - 1e-9)) }

func (w *World) collectGood(e Entity, amount float64) error {
	if !w.cargoActive() {
		return w.Econ.Store(map[string]float64{e.String(): amount})
	}
	d := w.actingDrone
	if cargoAmount(d)+amount > float64(w.Econ.CargoCapacity())+1e-9 {
		return ErrFull
	}
	if d.Hold == nil {
		d.Hold = map[string]float64{}
	}
	d.Hold[e.String()] += amount
	w.Cargo = cargoUsed(d)
	return nil
}

func (w *World) unloadColony() (bool, string, bool) {
	d := w.actingDrone
	if cargoUsed(d) == 0 {
		return false, say("الشحنة فاضية", "Cargo is empty"), false
	}
	if err := w.Econ.Store(d.Hold); err != nil {
		return false, buildingMessage(err), false
	}
	clear(d.Hold)
	w.Cargo = 0
	w.syncInv()
	return true, "", false
}

func (w *World) Carry(id int) map[Entity]int {
	w.mu.Lock()
	defer w.mu.Unlock()
	cargo := map[Entity]int{}
	if !w.cargoActive() {
		return cargo
	}
	for good, n := range w.Drones[id].Hold {
		cargo[entityByName[good]] = int(n)
	}
	return cargo
}

func (w *World) CargoFree(id int) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.cargoActive() {
		return 0
	}
	return max(0, w.Econ.CargoCapacity()-cargoUsed(w.Drones[id]))
}

func (w *World) Storage() (map[string]int, map[string]int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	used, capacity := map[string]int{}, map[string]int{}
	if w.Econ != nil {
		w.Econ.RefreshStorage()
		for group, st := range w.Econ.Storage {
			used[group], capacity[group] = st.Used, st.Cap
		}
	}
	return used, capacity
}

func (w *World) loadFleet() {
	w.Drones = nil
	for _, saved := range w.Econ.Fleet {
		d := &Drone{ID: len(w.Drones), Battery: saved.Battery, Hold: maps.Clone(saved.Cargo)}
		d.fault, d.stunnedUntil = saved.Fault, w.Tick+saved.Stunned
		d.Cargo = cargoUsed(d)
		w.Drones = append(w.Drones, d)
	}
	d := w.Drones[0]
	w.X, w.Y, w.Z, w.Cargo = d.X, d.Y, d.Z, d.Cargo
}

func (w *World) syncInfrastructure() {
	st := w.Econ
	st.Fleet = nil
	for _, d := range w.Drones {
		hold := maps.Clone(d.Hold)
		if hold == nil {
			hold = map[string]float64{}
		}
		st.Fleet = append(st.Fleet, econ.DroneState{Fault: d.fault, Stunned: max(0, d.stunnedUntil-w.Tick), ID: d.ID, Battery: d.Battery, Cargo: hold, Cap: st.CargoCapacity()})
	}
	st.RefreshStorage()
	st.Rent = st.RentRate(len(w.Drones))
}
