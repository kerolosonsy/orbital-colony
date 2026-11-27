package sim

import (
	"io"
	"testing"

	"orbital/colony/econ"
)

func TestColonyRadarFindsCrystals(t *testing.T) {
	st := econ.New(1)
	w := NewColony(io.Discard, st, 1)
	w.Realtime, w.ActionDur = false, 0
	for y := 0; y < w.H; y++ {
		for x := 0; x < w.W; x++ {
			if w.At(x, y).Entity == Crystal {
				w.At(x, y).Entity = Empty
			}
		}
	}
	want := Point{w.W - 1, w.H - 1}
	w.At(want.X, want.Y).Entity = Crystal
	if got := w.TargetList(); got != nil {
		t.Fatalf("targets visible without radar: %v", got)
	}
	st.Levels["api.radar"] = 1
	w.Start()
	if got := w.TargetList(); len(got) != 1 || got[0] != want {
		t.Fatalf("TargetList = %v", got)
	}
	if p := w.Nearest(0); p == nil || *p != want {
		t.Fatalf("Nearest = %v", p)
	}
}
