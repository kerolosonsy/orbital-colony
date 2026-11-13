package sim

import (
	"encoding/json"

	"orbital/colony/econ"
)

type SystemSnapshot struct {
	Credits   float64            `json:"credits"`
	Sold      float64            `json:"sold"`
	Commands  uint64             `json:"commands"`
	Buildings int                `json:"buildings"`
	Inventory map[string]float64 `json:"inventory"`
}

func (w *World) SystemStats() SystemSnapshot {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Econ == nil {
		return SystemSnapshot{}
	}
	return systemSnapshot(w.Econ)
}

func systemSnapshot(st *econ.State) SystemSnapshot {
	inv := make(map[string]float64, len(st.Inv))
	for good, amount := range st.Inv {
		inv[good] = amount
	}
	return SystemSnapshot{Credits: st.Credits, Sold: st.Systems.Sold, Commands: st.Systems.Commands, Buildings: len(st.Buildings), Inventory: inv}
}

func (w *World) controlSystemCommand(line []byte) bool {
	var cmd struct {
		ControlResult *struct {
			ID string `json:"id"`
			OK bool   `json:"ok"`
		} `json:"controlResult"`
		ControlSnapshot bool `json:"controlSnapshot"`
	}
	if json.Unmarshal(line, &cmd) != nil || cmd.ControlResult == nil && !cmd.ControlSnapshot {
		return false
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if cmd.ControlResult != nil {
		w.Econ.VerifyControl(cmd.ControlResult.ID, cmd.ControlResult.OK)
		w.checkContract()
	}
	w.econSnapshot("econ", "", true, true)
	return true
}
