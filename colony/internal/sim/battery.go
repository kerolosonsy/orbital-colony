package sim

import "orbital/colony/econ"

func (w *World) batteryActive() bool {
	if w.Econ != nil {
		for _, b := range w.Econ.Buildings {
			if b.Kind == "refinery" {
				return true
			}
		}
	}
	return false
}

func (w *World) actionSpeed(d *Drone) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.batteryActive() && d.Battery == 0 {
		return 2
	}
	return 1
}

func (w *World) Battery(id int) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Drones[id].Battery
}

func (w *World) Charge(id int) bool {
	return w.act(id, "charge", func() (bool, string, bool) {
		if !w.batteryActive() || w.X != 0 || w.Y != 0 || w.Z != 0 {
			return false, say("الشحن محتاج مصفاة ودرون عند القاعدة", "Charging requires a refinery and a drone at base"), false
		}
		if w.Econ.Inv["fuel"] < 1 {
			return false, buildingMessage(econ.ErrNoInput), false
		}
		w.Econ.Inv["fuel"]--
		w.actingDrone.Battery = econ.BatteryCapacity
		w.syncInv()
		return true, "", false
	})
}
