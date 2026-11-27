package sim

import (
	"io"
	"testing"

	"orbital/colony/econ"
)

func TestColonyRestartSpawnReturnsSavedDrones(t *testing.T) {
	st := econ.New(1)
	st.Levels["api.spawn"] = 1
	st.Fleet = []econ.DroneState{{ID: 0, Battery: 90}, {ID: 1, Battery: 80}, {ID: 2, Battery: 70}}
	w := NewColony(io.Discard, st, 1)
	w.Realtime, w.ActionDur = false, 0
	w.Start()
	for want := 1; want <= 2; want++ {
		if id := w.Spawn(); id != want {
			t.Fatalf("Spawn = %d, want saved drone %d", id, want)
		}
	}
	if len(w.Drones) != 3 {
		t.Fatalf("restart spawned extra drones: %d", len(w.Drones))
	}
}
