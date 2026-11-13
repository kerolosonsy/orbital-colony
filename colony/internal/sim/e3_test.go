package sim

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"testing"

	"orbital/colony/econ"
)

func powerBuilding(w *World, kind string, x int) *econ.Building {
	w.Econ.Levels[kind] = 1
	b := &econ.Building{Kind: kind, X: x % 3, Y: x / 3, Level: 1, Enabled: true, In: map[string]float64{}, Out: map[string]float64{}}
	if kind == "generator" {
		b.Mode = "solar"
	}
	if kind == "collector" {
		b.Recipe = "helium"
	}
	w.Econ.Buildings = append(w.Econ.Buildings, b)
	w.Econ.Credits = 1000
	w.applyEcon()
	return b
}

func TestPowerBaselineAndDayNightPriorityBrownout(t *testing.T) {
	w, _ := e2World()
	a := powerBuilding(w, "composter", 0)
	a.Recipe, a.In["moss"] = "biomass", 8
	w.Wait(0)
	if !a.Powered || a.Progress != 1 {
		t.Fatal("baseline did not run early composter")
	}
	g := powerBuilding(w, "generator", 1)
	b := powerBuilding(w, "refinery", 2)
	b.Recipe, b.In["water"], b.In["biomass"] = "fuel", 1, 2
	c := powerBuilding(w, "smelter", 3)
	c.Recipe, c.In["iron"] = "plate", 2
	if err := w.BuildingSetPriority(0, Point{2, 0}, 5); err != nil {
		t.Fatal(err)
	}
	if err := w.BuildingSetPriority(0, Point{}, -1); err != nil {
		t.Fatal(err)
	}
	beforePriority := a.Progress
	w.Wait(0)
	if !b.Powered || !c.Powered || a.Powered || a.Progress != beforePriority {
		t.Fatalf("priority did not select 10kW loads: a=%+v b=%+v c=%+v", a, b, c)
	}
	w.Tick = 599
	w.Wait(0)
	before := c.Progress
	if supply, demand := w.Power(); supply != 0 || demand != 11 || w.Daylight() || c.Powered {
		t.Fatalf("night power=%g demand=%g", supply, demand)
	}
	w.Wait(0)
	if c.Progress != before {
		t.Fatal("solar building progressed at night")
	}
	if err := w.BuildingEnable(0, Point{1, 0}, false); err != nil {
		t.Fatal(err)
	}
	w.Tick = 1199
	w.Wait(0)
	if supply, _ := w.Power(); supply != 0 || g.Enabled {
		t.Fatal("disabled generator supplied baseline or solar")
	}
}

func TestCargoActivatesWithFirstBuildingAndRemainsDroneScoped(t *testing.T) {
	w, _ := e2World()
	w.PlantAt(0, 0, Moss, 0)
	if !w.Harvest(0) || w.Count(Moss) != 1 || len(w.Carry(0)) != 0 {
		t.Fatal("early harvest changed")
	}
	w.Econ.Levels["composter"], w.Econ.Inv["moss"], w.Econ.Credits = 1, 10, 100
	if err := w.Build(0, "composter"); err != nil {
		t.Fatal(err)
	}
	w.Move(0, East)
	w.PlantAt(1, 0, Moss, 0)
	if !w.Harvest(0) || w.Carry(0)[Moss] != 1 || w.Count(Moss) != 6 || w.CargoFree(0) != 9 {
		t.Fatal("building harvest bypassed cargo")
	}
	if w.Unload(0) || w.Carry(0)[Moss] != 1 {
		t.Fatal("unloaded away from base")
	}
	w.Move(0, West)
	w.Econ.Levels["api.spawn"], w.Econ.Levels["hangar"] = 1, 1
	id := w.Spawn()
	if id != 1 || len(w.Carry(id)) != 0 || w.Carry(0)[Moss] != 1 {
		t.Fatal("cargo leaked between drones")
	}
	if !w.Unload(0) || w.Count(Moss) != 7 || len(w.Carry(0)) != 0 {
		t.Fatal("unload lost goods")
	}
	w.Z = 1
	w.Drones[0].Hold["moss"] = 1
	if w.Unload(0) {
		t.Fatal("unloaded underground at base coordinates")
	}
	w.Econ.Levels["cargo"] = 2
	if w.CargoFree(0) != 19 {
		t.Fatal("cargo upgrade not applied")
	}
}

func TestStorageAndCargoOverflowPreserveHarvestGatherAndBuffers(t *testing.T) {
	w, _ := e2World()
	w.Econ.StorageBase = map[string]int{"crops": 0, "raw": 0, "crafted": 0}
	w.PlantAt(0, 0, Moss, 0)
	w.Econ.Levels["api.try"] = 1
	if _, err := w.TryHarvest(0); !errors.Is(err, ErrFull) || w.At(0, 0).Entity != Moss {
		t.Fatal("storage harvest overflow destroyed crop")
	}
	w.At(0, 0).Entity = Empty
	if _, err := w.TrySift(0); !errors.Is(err, ErrFull) || w.At(0, 0).SiftUsed != 0 {
		t.Fatal("sift overflow consumed stock")
	}
	w.At(0, 0).Ground = Ice
	if _, err := w.TryMelt(0); !errors.Is(err, ErrFull) || w.At(0, 0).Ground != Ice {
		t.Fatal("melt overflow consumed ice")
	}
	b := powerBuilding(w, "composter", 0)
	b.Recipe, b.Out["biomass"] = "biomass", 2
	w.Drones[0].Hold = map[string]float64{"moss": 10}
	w.Cargo = 10
	if _, err := w.BuildingTake(0, Point{}, Biomass, 2); !errors.Is(err, ErrFull) || b.Out["biomass"] != 2 {
		t.Fatal("Take overflow lost output")
	}
	if w.Unload(0) || w.Carry(0)[Moss] != 10 {
		t.Fatal("storage overflow discarded cargo")
	}
	if _, err := w.BuildingPut(0, Point{}, Moss, 4); err != nil || w.Carry(0)[Moss] != 6 || w.Count(Moss) != 0 {
		t.Fatal("Put did not use cargo first")
	}
	if n, err := w.BuildingTake(0, Point{}, Biomass, 2); err != nil || n != 2 || w.Carry(0)[Biomass] != 2 {
		t.Fatal("Take did not fill cargo")
	}
	w = reloadColony(t, w)
	if w.Carry(0)[Moss] != 6 || w.Carry(0)[Biomass] != 2 {
		t.Fatal("save lost cargo")
	}
}

func TestDroneBatteryDrainHalfSpeedChargeAndPersistence(t *testing.T) {
	w, _ := e2World()
	w.Wait(0)
	if w.Battery(0) != 100 {
		t.Fatal("battery active before refinery")
	}
	powerBuilding(w, "refinery", 0)
	w.Wait(0)
	if w.Battery(0) != 99 {
		t.Fatal("action did not drain battery")
	}
	w.Drones[0].Battery = 1
	before := w.Tick
	w.Wait(0)
	if w.Tick-before != 1 || w.Battery(0) != 0 {
		t.Fatal("last charged action slowed early")
	}
	w.Wait(0)
	if w.Tick-before != 3 {
		t.Fatal("empty battery did not halve speed")
	}
	if w.Charge(0) {
		t.Fatal("charged without fuel")
	}
	w.Econ.Inv["fuel"] = 2
	w.Move(0, East)
	if w.Charge(0) || w.Econ.Inv["fuel"] != 2 {
		t.Fatal("charged away from base")
	}
	w.Move(0, West)
	if !w.Charge(0) || w.Econ.Inv["fuel"] != 1 || w.Battery(0) != 100 {
		t.Fatal("charge did not spend one fuel and refill")
	}
	w.Wait(0)
	w = reloadColony(t, w)
	if w.Battery(0) != 99 {
		t.Fatal("reload refilled battery")
	}
}

func TestRentStopsProductionAtZeroAndResumesAfterExport(t *testing.T) {
	w, _ := e2World()
	b := powerBuilding(w, "composter", 0)
	b.Recipe, b.In["moss"], w.Econ.Credits = "biomass", 4, .0001
	w.Wait(0)
	if w.Econ.Credits != 0 || b.Powered || b.Progress != 0 || w.Econ.Rent != .5 {
		t.Fatal("unpaid rent allowed production or negative balance")
	}
	w.Econ.Inv["iron"] = 1
	w.Export(0, Iron, 1)
	w.Wait(0)
	if !b.Powered || b.Progress != 1 {
		t.Fatal("credits did not restart production")
	}
	w.Econ.Levels["api.spawn"], w.Econ.Levels["hangar"] = 1, 1
	w.Spawn()
	if w.Econ.Rent != 1 {
		t.Fatal("extra drone missing from rent")
	}
}

func TestNightHeliumCompletesRawToCellChain(t *testing.T) {
	w, _ := e2World()
	w.Tick, w.lastProductionTick = 600, 600
	g := powerBuilding(w, "generator", 0)
	g.Mode, g.In["fuel"] = "burner", 10
	collector := powerBuilding(w, "collector", 1)
	smelter := powerBuilding(w, "smelter", 2)
	refinery := powerBuilding(w, "refinery", 3)
	lab := powerBuilding(w, "reactor-lab", 4)
	lab.Priority = 10
	// Two generators cover the 32 kW chain; the second has ample fuel for the run.
	g2 := powerBuilding(w, "generator", 5)
	g2.Mode, g2.In["fuel"] = "burner", 10
	w.Econ.Inv = map[string]float64{"iron": 4, "titanium": 1, "water": 4, "uranium": 1, "fuel": 10}
	w.syncInv()
	put := func(b *econ.Building, e Entity, n int) {
		t.Helper()
		w.X, w.Y = b.X, b.Y
		if moved, err := w.BuildingPut(0, Point{b.X, b.Y}, e, n); err != nil || moved != n {
			t.Fatalf("Put %s: %d %v", e, moved, err)
		}
	}
	produce := func(b *econ.Building, e Entity, n int) {
		t.Helper()
		p := Point{b.X, b.Y}
		if err := w.BuildingSetRecipe(0, p, e); err != nil {
			t.Fatal(err)
		}
		deadline := w.Tick + 2000
		for b.Out[e.String()] < float64(n) && w.Tick < deadline {
			if g.In["fuel"] < 1 {
				put(g, Fuel, 5)
				put(g2, Fuel, 5)
			}
			w.Wait(0)
		}
		if b.Out[e.String()] < float64(n) {
			t.Fatalf("%s stalled: %+v", e, b)
		}
		w.X, w.Y = b.X, b.Y
		if _, err := w.BuildingTake(0, p, e, n); err != nil {
			t.Fatal(err)
		}
	}
	produce(collector, Helium, 2)
	if err := w.BuildingSetRecipe(0, Point{smelter.X, smelter.Y}, Plate); err != nil {
		t.Fatal(err)
	}
	put(smelter, Iron, 4)
	produce(smelter, Plate, 2)
	if err := w.BuildingSetRecipe(0, Point{smelter.X, smelter.Y}, Alloy); err != nil {
		t.Fatal(err)
	}
	put(smelter, Plate, 2)
	put(smelter, Titanium, 1)
	produce(smelter, Alloy, 1)
	if err := w.BuildingSetRecipe(0, Point{refinery.X, refinery.Y}, Coolant); err != nil {
		t.Fatal(err)
	}
	put(refinery, Water, 4)
	put(refinery, Helium, 2)
	produce(refinery, Coolant, 2)
	if err := w.BuildingSetRecipe(0, Point{lab.X, lab.Y}, Cell); err != nil {
		t.Fatal(err)
	}
	put(lab, Uranium, 1)
	put(lab, Alloy, 1)
	put(lab, Coolant, 2)
	produce(lab, Cell, 1)
	w.X, w.Y = 0, 0
	if !w.Unload(0) || w.Count(Cell) != 1 {
		t.Fatal("full cell chain did not reach inventory")
	}
}

func TestInfrastructureWireAndFleetRoundTrip(t *testing.T) {
	w, trace := e2World()
	g := powerBuilding(w, "generator", 0)
	g.Priority = 3
	w.Drones[0].Hold = map[string]float64{"moss": 4}
	w.Drones[0].Battery = 80
	w.Wait(0)
	w.econSnapshot("state", "", true, true)
	lines := bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'})
	var ev struct {
		State struct {
			Power     econ.PowerStatus
			Storage   map[string]econ.StorageStatus
			Fleet     []econ.DroneState
			Rent      float64
			Buildings []map[string]any
		}
	}
	if err := json.Unmarshal(lines[len(lines)-1], &ev); err != nil {
		t.Fatal(err)
	}
	if ev.State.Power.Supply != 10 || !ev.State.Power.Day || ev.State.Storage["crops"].Cap != 2000 || len(ev.State.Fleet) != 1 || ev.State.Fleet[0].Battery != 80 || ev.State.Fleet[0].Cargo["moss"] != 4 || ev.State.Fleet[0].Cap != 10 || ev.State.Rent != .5 {
		t.Fatalf("state wire: %s", lines[len(lines)-1])
	}
	b := ev.State.Buildings[0]
	if b["powered"] != true || b["enabled"] != true || b["priority"] != float64(3) || b["mode"] != "solar" {
		t.Fatalf("building wire: %+v", b)
	}
	loaded := reloadColony(t, w)
	if !reflect.DeepEqual(loaded.Carry(0), w.Carry(0)) || math.Abs(loaded.Econ.Credits-w.Econ.Credits) > 1e-9 {
		t.Fatal("round trip changed holdings")
	}
}

func TestFullCargoLeavesCropsIceSilicaAndOreUntouched(t *testing.T) {
	for _, source := range []string{"crop", "melt", "sift", "digdown", "dig"} {
		t.Run(source, func(t *testing.T) {
			w, _ := e2World()
			powerBuilding(w, "silo", 2)
			w.Drones[0].Hold = map[string]float64{"moss": 10}
			w.Cargo = 10
			w.Econ.Levels["api.try"], w.Econ.Levels["api.dig"], w.Econ.Levels["drill"] = 1, 1, 1
			var err error
			switch source {
			case "crop":
				w.PlantAt(0, 0, Moss, 0)
				_, err = w.TryHarvest(0)
				if w.At(0, 0).Entity != Moss {
					t.Fatal("full cargo destroyed crop")
				}
			case "melt":
				w.At(0, 0).Ground = Ice
				_, err = w.TryMelt(0)
				if w.At(0, 0).Ground != Ice {
					t.Fatal("full cargo melted ice")
				}
			case "sift":
				_, err = w.TrySift(0)
				if w.At(0, 0).SiftUsed != 0 {
					t.Fatal("full cargo depleted silica")
				}
			case "digdown":
				w.AtZ(0, 0, 1).Entity = Iron
				err = w.TryDigDown(0)
				if w.At(0, 0).Ground != Regolith || w.AtZ(0, 0, 1).Ground != Stone || w.AtZ(0, 0, 1).Entity != Iron {
					t.Fatal("full cargo dug a shaft or lost ore")
				}
			case "dig":
				w.Z = 1
				w.AtZ(1, 0, 1).Entity = Iron
				err = w.TryDig(0, East)
				if w.AtZ(1, 0, 1).Ground != Stone || w.AtZ(1, 0, 1).Entity != Iron {
					t.Fatal("full cargo dug a tunnel or lost ore")
				}
			}
			if !errors.Is(err, ErrFull) || w.Carry(0)[Moss] != 10 {
				t.Fatalf("overflow = %v, cargo=%v", err, w.Carry(0))
			}
		})
	}
}

func TestCollectorPausesDuringDayAndNeedsPowerAtNight(t *testing.T) {
	w, _ := e2World()
	b := powerBuilding(w, "collector", 0)
	w.Tick = 598
	w.Wait(0)
	if b.Progress != 0 || b.Out["helium"] != 0 {
		t.Fatal("collector produced during daylight")
	}
	w.Wait(0)
	if b.Progress != 1 {
		t.Fatal("collector did not begin at night")
	}
	powerBuilding(w, "generator", 1)
	w.Wait(0)
	if b.Powered || b.Progress != 1 {
		t.Fatal("night solar powered collector")
	}
}

func TestPriorityChangeEmitsOtherBuildingsPowerTransition(t *testing.T) {
	w, trace := e2World()
	powerBuilding(w, "generator", 0)
	a := powerBuilding(w, "refinery", 1)
	b := powerBuilding(w, "refinery", 2)
	w.Wait(0)
	trace.Reset()
	if err := w.BuildingSetPriority(0, Point{2, 0}, 1); err != nil {
		t.Fatal(err)
	}
	seen := false
	for _, ev := range buildingEvents(t, trace) {
		seen = seen || ev.X == a.X && !ev.Building.Powered
	}
	if !seen || !b.Powered {
		t.Fatal("priority change hid another building's power transition")
	}
}

func TestDigRejectsAnExistingTunnelWithoutReportingSuccess(t *testing.T) {
	w, _ := e2World()
	w.Econ.Levels["api.dig"], w.Econ.Levels["drill"] = 1, 1
	w.Z = 1
	w.AtZ(1, 0, 1).Ground = Tunnel
	if err := w.TryDig(0, East); !errors.Is(err, econ.ErrTileNotFree) {
		t.Fatalf("digging tunnel: %v", err)
	}
	if w.Dig(0, East) {
		t.Fatal("bool Dig reported success on an existing tunnel")
	}
}

func TestNewDroneRentDoesNotApplyBeforeItsLaunch(t *testing.T) {
	w, _ := e2World()
	w.Econ.Levels["api.spawn"], w.Econ.Levels["hangar"], w.Econ.Credits = 1, 1, 10
	w.Tick = 600
	if w.Spawn() != 1 || w.Econ.Credits != 10 {
		t.Fatal("new drone charged rent for time before launch")
	}
	w.Tick = 1200
	w.econSnapshot("econ", "", true, true)
	if math.Abs(w.Econ.Credits-9.5) > 1e-9 {
		t.Fatal("launched drone did not pay a minute of rent")
	}
}
