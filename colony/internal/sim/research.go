package sim

import "orbital/colony/econ"

func (w *World) Research() float64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return 0
	}
	return w.Econ.RP
}

func (w *World) advanceResearch(b *econ.Building, tick int) bool {
	w.Econ.ConfigureBuilding(b)
	if b.Kind != "lab" {
		return false
	}
	before, wear := b.Progress, b.Wear
	w.Econ.StepResearch(b)
	if before == 0 && b.Progress > 0 || b.Wear != wear {
		w.emitBuilding(b, tick)
	}
	return true
}

func researchSeed(st *econ.State, seed int64) int64 {
	c, run := st.ActiveContract()
	if c == nil || !c.Repeatable {
		return seed
	}
	return st.Seed ^ int64(c.Chapter)*0x165667B19E3779F9 ^ int64(run.Tries)*0x27D4EB2F165667C5
}
