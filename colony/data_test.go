package colony

import (
	"errors"
	"io"
	"testing"

	"orbital/colony/internal/sim"
)

func dataWorld(t *testing.T) *sim.World {
	t.Helper()
	previous := world
	world = sim.New(io.Discard, 42)
	t.Cleanup(func() { world = previous })
	return world
}
func TestDataInputsAreCopiesAndFaultsAreTyped(t *testing.T) {
	w := dataWorld(t)
	w.Data.Callsigns = []string{"قمر🚀"}
	w.Data.Readings = []string{"12", "bad"}
	w.Data.Orders = []string{"dock"}
	w.Data.Faults = []*sim.FaultError{nil, {Index: 1, Code: "heat", Message: "hot"}}
	Callsigns()[0] = "changed"
	Readings()[0] = "changed"
	Orders()[0] = "changed"
	if Callsigns()[0] != "قمر🚀" || Readings()[0] != "12" || Orders()[0] != "dock" {
		t.Fatal("mutated source")
	}
	if Diagnose(0) != nil {
		t.Fatal("healthy diagnostic failed")
	}
	var fault *FaultError
	if !errors.As(Diagnose(1), &fault) || fault.Code != "heat" {
		t.Fatal("missing typed fault")
	}
	fault.Code = "changed"
	if w.Data.Faults[1].Code != "heat" {
		t.Fatal("mutated fault")
	}
	for _, i := range []int{-1, 2} {
		if Diagnose(i) == nil {
			t.Fatal("accepted invalid diagnostic")
		}
	}
}
func TestExamAnswersValidateBoundsAndProtectChoices(t *testing.T) {
	w := dataWorld(t)
	w.Data.Questions = []sim.Question{{Prompt: "route?", Choices: []string{"north", "south"}}, {Prompt: "empty"}}
	questions := Exam()
	questions[0].Prompt = "changed"
	questions[0].Choices[0] = "changed"
	if Exam()[0].Prompt != "route?" || Exam()[0].Choices[0] != "north" {
		t.Fatal("mutated exam")
	}
	for _, pair := range [][2]int{{-1, 0}, {2, 0}, {0, -1}, {0, 2}, {1, 0}} {
		if Answer(pair[0], pair[1]) {
			t.Fatal(pair)
		}
	}
	if !Answer(0, 1) || w.Data.Answers[0] != 1 {
		t.Fatal(w.Data.Answers)
	}
	w.Data.TotalCents, w.Data.People = 1001, 3
	if total, people := Invoice(); total != 1001 || people != 3 {
		t.Fatalf("%d %d", total, people)
	}
	Report("checked")
	if len(w.Reports) != 1 || w.Reports[0] != "checked" {
		t.Fatal(w.Reports)
	}
}
func TestDeferredUndockSurvivesDrillPanic(t *testing.T) {
	for _, panics := range []bool{false, true} {
		t.Run(map[bool]string{false: "normal", true: "fault"}[panics], func(t *testing.T) {
			w := dataWorld(t)
			w.Data.DrillPanic = panics
			recovered := false
			func() { Dock(); defer Undock(); defer func() { recovered = recover() != nil }(); Drill() }()
			if recovered != panics || w.Data.Docked || w.Data.Drills != 1 || w.Data.Undocks != 1 {
				t.Fatalf("recovered=%v data=%+v", recovered, w.Data)
			}
			Undock()
			if w.Data.Undocks != 1 {
				t.Fatal("undock not idempotent")
			}
		})
	}
}
