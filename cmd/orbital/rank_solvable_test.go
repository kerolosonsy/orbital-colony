package main

import (
	"os"
	"path/filepath"
	"testing"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

func TestRankReferenceOutproducesStarterOnSameClock(t *testing.T) {
	if testing.Short() {
		t.Skip("plays colony programs in real time")
	}
	code, err := os.ReadFile(filepath.Join("testdata", "colony", "e5-rank", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, program := range []struct {
		name string
		code string
		want int
	}{{"reference", string(code), 1}, {"starter", levels.Colony.Starter, 0}} {
		t.Run(program.name, func(t *testing.T) {
			st := econ.New(11)
			st.Levels["field"] = 1
			st.Normalize()
			for i := range st.Tiles {
				st.Tiles[i] = econ.Tile{G: "soil"}
			}
			var final *econ.State
			var promotedAt int64
			playBuildingProgram(t, st, program.code, func(ev buildingProgramEvent) bool {
				if ev.Op == "rank" && ev.State != nil && promotedAt == 0 {
					promotedAt = ev.State.Tick
				}
				if ev.State != nil && ev.State.Tick >= econ.Window+600 {
					final = ev.State
					return true
				}
				return false
			})
			if program.want == 1 && (promotedAt < econ.Window || promotedAt > econ.Window+600) {
				t.Fatalf("promotion outside observation window: %d", promotedAt)
			}
			if final.Rank != program.want {
				t.Fatalf("rank=%d want=%d, progress=%+v", final.Rank, program.want, final.RankProgress())
			}
			t.Logf("%s: rank %d at tick %d, promoted at %d, goal %+v", program.name, final.Rank, final.Tick, promotedAt, final.RankProgress())
		})
	}
}
