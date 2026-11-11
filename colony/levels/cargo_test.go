package levels

import (
	"io"
	"testing"
)

func TestCargoMissionAlwaysHasEnoughCrystal(t *testing.T) {
	// Seed 56940 previously produced only 29 crystal for a target of 30.
	for _, seed := range []int64{0, 1, 56940} {
		w := Start("8-2-cargo-runs", io.Discard, seed)
		total := 0
		for _, p := range w.Targets {
			size := w.At(p.X, p.Y).Size
			if size > w.Capacity {
				t.Fatalf("seed %d: crystal of size %d cannot fit in cargo", seed, size)
			}
			total += size
		}
		if total < 30 {
			t.Fatalf("seed %d: only %d crystal for a target of 30", seed, total)
		}
	}
}
