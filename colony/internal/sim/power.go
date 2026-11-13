package sim

import "orbital/colony/econ"

func (w *World) Power() (float64, float64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return 0, 0
	}
	return w.Econ.Power.Supply, w.Econ.Power.Demand
}

func (w *World) Daylight() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return econ.Daylight(int64(w.Tick))
}

func (w *World) BuildingSetMode(id int, p Point, mode string) error {
	return w.buildingAction(id, p, "mode", func(b *econ.Building) error {
		if err := b.SetMode(mode); err != nil {
			return err
		}
		w.refreshPower(b)
		w.emitBuilding(b, w.Tick)
		return nil
	})
}

func (w *World) BuildingSetPriority(id int, p Point, priority int) error {
	return w.buildingAction(id, p, "priority", func(b *econ.Building) error {
		b.Priority = priority
		w.refreshPower(b)
		w.emitBuilding(b, w.Tick)
		return nil
	})
}

func (w *World) BuildingEnable(id int, p Point, enabled bool) error {
	return w.buildingAction(id, p, "enable", func(b *econ.Building) error {
		b.Enabled = enabled
		w.refreshPower(b)
		w.emitBuilding(b, w.Tick)
		return nil
	})
}

func (w *World) productionStep(tick int) {
	paid := w.Econ.PayRent(len(w.Drones))
	supply := 0.0
	if paid {
		supply = w.Econ.ConsumePower(int64(tick))
	}
	powered, demand := econ.Allocate(supply, w.Econ.Buildings)
	w.Econ.Power = econ.PowerStatus{Supply: supply, Demand: demand, Day: econ.Daylight(int64(tick))}
	for i, b := range w.Econ.Buildings {
		b.Level = w.Econ.Level(b.Kind)
		before, wear, wasPowered := b.Progress, b.Wear, b.Powered
		b.Powered = paid && powered[i]
		w.Econ.ConfigureBuilding(b)
		if b.Kind == "lab" {
			if b.Powered {
				w.advanceResearch(b, tick)
			}
			if wasPowered != b.Powered {
				w.emitBuilding(b, tick)
			}
			continue
		}
		if b.Kind != "collector" || !w.Econ.Power.Day {
			_ = b.Step(1, b.Powered)
		}
		if before == 0 && b.Progress > 0 || b.Wear != wear || wasPowered != b.Powered {
			w.emitBuilding(b, tick)
		}
	}
}

func (w *World) refreshPower(changed *econ.Building) {
	before := make([]bool, len(w.Econ.Buildings))
	for i, b := range w.Econ.Buildings {
		before[i] = b.Powered
	}
	w.Econ.RefreshPower(int64(w.Tick))
	for i, b := range w.Econ.Buildings {
		if b != changed && b.Powered != before[i] {
			w.emitBuilding(b, w.Tick)
		}
	}
}
