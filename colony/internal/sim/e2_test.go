package sim

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"reflect"
	"testing"

	"orbital/colony/econ"
)

func e2World() (*World, *bytes.Buffer) {
	trace := &bytes.Buffer{}
	w := NewColony(trace, econ.New(11), 11)
	w.Realtime = false
	return w, trace
}

func reloadColony(t *testing.T, w *World) *World {
	t.Helper()
	w.syncEcon()
	blob, err := json.Marshal(w.Econ)
	if err != nil {
		t.Fatal(err)
	}
	var st econ.State
	if err := json.Unmarshal(blob, &st); err != nil {
		t.Fatal(err)
	}
	restored := NewColony(io.Discard, &st, st.Seed)
	restored.Realtime = false
	return restored
}

func buildingEvents(t *testing.T, trace *bytes.Buffer) []Event {
	t.Helper()
	var events []Event
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "building" {
			events = append(events, ev)
		}
	}
	return events
}

func TestColonyMeltRefreezesAfter600TicksAcrossSave(t *testing.T) {
	w, _ := e2World()
	w.At(0, 0).Ground = Ice
	if !w.Melt(0) || w.Count(Water) != 1 || w.At(0, 0).Ground != Regolith {
		t.Fatal("Melt did not produce water and expose regolith")
	}
	melted := w.Tick
	if w.Melt(0) || w.Count(Water) != 1 {
		t.Fatal("melted tile paid twice")
	}
	w = reloadColony(t, w)
	w.Tick = melted + 598
	w.Wait(0)
	if w.At(0, 0).Ground != Regolith {
		t.Fatal("ice returned before 600 ticks")
	}
	w.Wait(0)
	if w.At(0, 0).Ground != Ice || w.At(0, 0).MeltAt != 0 {
		t.Fatal("ice did not return at 600 ticks")
	}
	if !w.Melt(0) || w.Count(Water) != 2 {
		t.Fatal("regrown ice cannot be melted")
	}
}

func TestColonySiftDepletesAndRefillsAcrossSave(t *testing.T) {
	w, _ := e2World()
	for i := 0; i < 3; i++ {
		if !w.Sift(0) {
			t.Fatal("initial stock missing")
		}
	}
	if w.Sift(0) || w.Count(Silica) != 3 {
		t.Fatal("depleted tile still pays")
	}
	if !w.Move(0, East) || !w.Sift(0) {
		t.Fatal("different tile shares depleted stock")
	}
	if !w.Move(0, West) {
		t.Fatal("return to depleted tile")
	}
	w = reloadColony(t, w)
	w.Tick = 99
	if w.Sift(0) {
		t.Fatal("stock refilled before 100 ticks")
	}
	if !w.Sift(0) || w.Sift(0) || w.Count(Silica) != 5 {
		t.Fatal("expected exactly one unit after 100 ticks")
	}
	w.Tick = 1000
	for i := 0; i < 3; i++ {
		if !w.Sift(0) {
			t.Fatal("idle tile failed to refill")
		}
	}
	if w.Sift(0) {
		t.Fatal("refill accumulated beyond the stock cap")
	}
}

func TestColonyGatheringRequiresSurfaceAndRightGround(t *testing.T) {
	for _, ground := range []Ground{Soil, Rock, Stone, Tunnel, Shaft} {
		w, _ := e2World()
		w.At(0, 0).Ground = ground
		if w.Melt(0) || w.Sift(0) || w.Count(Water) != 0 || w.Count(Silica) != 0 {
			t.Fatalf("gathered on %s", ground)
		}
	}
	w, _ := e2World()
	w.At(0, 0).Ground = Ice
	if w.Sift(0) {
		t.Fatal("sifted ice")
	}
	w.Z = 1
	for _, g := range []Ground{Ice, Regolith} {
		w.here().Ground = g
		if w.Melt(0) || w.Sift(0) {
			t.Fatal("surface gathering worked underground")
		}
	}
	mission := New(io.Discard, 1)
	mission.Init(1, 1, Ice)
	if mission.Melt(0) {
		t.Fatal("Melt worked outside colony mode")
	}
	mission.At(0, 0).Ground = Regolith
	if mission.Sift(0) {
		t.Fatal("Sift worked outside colony mode")
	}
}

func TestColonyDigHarvestsNewOres(t *testing.T) {
	for _, ore := range []Entity{Copper, Sulfur, Titanium, Uranium} {
		w, _ := e2World()
		w.Econ.Levels["api.dig"], w.Econ.Levels["drill"] = 1, 5
		w.AtZ(0, 0, 1).Entity = ore
		if !w.DigDown(0) || w.Count(ore) != 1 || !w.Descend(0) {
			t.Fatalf("DigDown failed for %s", ore)
		}
		w.AtZ(1, 0, 1).Entity = ore
		if !w.Dig(0, East) || w.Count(ore) != 2 || w.Econ.Inv[ore.String()] != 2 {
			t.Fatalf("Dig failed for %s", ore)
		}
		if earned := w.Export(0, ore, 1); earned != econ.PriceV3[ore.String()] || w.Count(ore) != 1 {
			t.Fatalf("Export failed for %s: %g", ore, earned)
		}
	}
}

func TestColonyBuildPreconditionsLeavePossessionsUntouched(t *testing.T) {
	for _, tc := range []struct {
		name  string
		kind  string
		setup func(*World)
		want  error
	}{
		{"unknown kind", "missing", func(*World) {}, econ.ErrUnknownKind},
		{"no blueprint", "composter", func(w *World) { w.Econ.Levels["composter"] = 0 }, econ.ErrBlueprintNotOwned},
		{"underground", "composter", func(w *World) { w.Z = 1; w.here().Ground = Regolith }, econ.ErrTileNotFree},
		{"ice", "composter", func(w *World) { w.here().Ground = Ice }, econ.ErrTileNotFree},
		{"rock", "composter", func(w *World) { w.here().Ground = Rock }, econ.ErrTileNotFree},
		{"shaft", "composter", func(w *World) { w.here().Ground = Shaft }, econ.ErrTileNotFree},
		{"entity", "composter", func(w *World) { w.here().Entity = Moss }, econ.ErrTileNotFree},
		{"building", "composter", func(w *World) { w.Econ.Buildings = []*econ.Building{{Kind: "silo"}} }, econ.ErrTileNotFree},
		{"credits", "composter", func(w *World) { w.Econ.Credits = 24 }, econ.ErrNoInput},
		{"goods", "composter", func(w *World) { w.Econ.Inv["moss"] = 3 }, econ.ErrNoInput},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, trace := e2World()
			w.Econ.Levels["composter"], w.Econ.Credits, w.Econ.Inv["moss"] = 1, 25, 4
			tc.setup(w)
			credits, moss, buildings := w.Econ.Credits, w.Econ.Inv["moss"], len(w.Econ.Buildings)
			if err := w.Build(0, tc.kind); !errors.Is(err, tc.want) {
				t.Fatalf("Build = %v, want %v", err, tc.want)
			}
			if w.Tick != 1 || w.Econ.Credits != credits || w.Econ.Inv["moss"] != moss || len(w.Econ.Buildings) != buildings || len(buildingEvents(t, trace)) != 0 {
				t.Fatal("failed build spent possessions or emitted placement")
			}
		})
	}
}

func TestColonyBuildingsPlacePersistBlockFarmingAndAllowMovement(t *testing.T) {
	for _, ground := range []Ground{Regolith, Soil} {
		w, trace := e2World()
		w.At(0, 0).Ground = ground
		w.Econ.Levels["composter"], w.Econ.Credits, w.Econ.Inv["moss"] = 2, 25, 4
		if w.Here(0) != nil || len(w.Buildings()) != 0 {
			t.Fatal("fresh colony contains buildings")
		}
		if err := w.Build(0, "composter"); err != nil {
			t.Fatal(err)
		}
		if w.Econ.Credits != 0 || w.Count(Moss) != 0 || w.Here(0) == nil || len(w.Buildings()) != 1 || w.BuildingStatus(Point{}).Level != 2 {
			t.Fatal("placement lost cost, location or level")
		}
		if w.Plant(0, Moss) || w.Till(0) || !w.Move(0, East) || w.Here(0) != nil || !w.Move(0, West) {
			t.Fatal("building farming or movement rules failed")
		}
		if v := w.view(0, 0); v.B != "composter" || w.viewZ(0, 0, 1).B != "" {
			t.Fatal("tile building view lost surface scope")
		}
		if len(buildingEvents(t, trace)) != 1 {
			t.Fatal("placement should emit one building event")
		}
		w = reloadColony(t, w)
		if w.Here(0) == nil || w.BuildingStatus(Point{}).Kind != "composter" {
			t.Fatal("building lost on reload")
		}
	}
}

func TestColonyBuildingTransfersRequireActingDroneOnSurfaceTile(t *testing.T) {
	w, _ := e2World()
	w.Econ.Levels["api.spawn"], w.Econ.Levels["hangar"] = 1, 1
	w.applyEcon()
	id := w.Spawn()
	if id != 1 {
		t.Fatal("could not spawn second drone")
	}
	w.Econ.Levels["composter"], w.Econ.Credits, w.Econ.Inv["moss"] = 1, 25, 20
	if err := w.Build(id, "composter"); err != nil {
		t.Fatal(err)
	}
	if err := w.BuildingSetRecipe(id, Point{}, Biomass); err != nil {
		t.Fatal(err)
	}
	b := w.Econ.Buildings[0]
	b.Out["biomass"] = 2
	if !w.Move(id, East) {
		t.Fatal("second drone could not move")
	}
	for _, drone := range []int{id, 0} {
		if drone == 0 {
			w.Z = 1
		}
		before := w.Tick
		if _, err := w.BuildingPut(drone, Point{}, Moss, 4); !errors.Is(err, econ.ErrNotHere) {
			t.Fatalf("Put from other tile/layer: %v", err)
		}
		if _, err := w.BuildingTake(drone, Point{}, Biomass, 1); !errors.Is(err, econ.ErrNotHere) {
			t.Fatalf("Take from other tile/layer: %v", err)
		}
		if w.Tick != before+2 || b.In["moss"] != 0 || b.Out["biomass"] != 2 {
			t.Fatal("failed transfers changed buffers or did not take time")
		}
	}
	if !w.Move(id, West) || w.Here(id) == nil || w.Here(0) != nil {
		t.Fatal("Here read another drone's position")
	}
	if n, err := w.BuildingPut(id, Point{}, Moss, 100); n != 10 || err != nil || w.Count(Moss) != 6 {
		t.Fatalf("Put capacity/inventory: %d, %v", n, err)
	}
	b.Enabled = false
	if n, err := w.BuildingTake(id, Point{}, Biomass, 100); n != 2 || err != nil || w.Count(Biomass) != 2 {
		t.Fatalf("Take output/inventory: %d, %v", n, err)
	}
	v := w.BuildingStatus(Point{})
	v.In["moss"] = 999
	if b.In["moss"] == 999 {
		t.Fatal("status exposed mutable buffers")
	}
}

func TestColonyBuildingTransferErrorsAndRepair(t *testing.T) {
	w, _ := e2World()
	w.Econ.Levels["composter"], w.Econ.Credits, w.Econ.Inv["moss"] = 1, 25, 4
	if err := w.Build(0, "composter"); err != nil {
		t.Fatal(err)
	}
	if err := w.BuildingSetRecipe(0, Point{}, Plate); !errors.Is(err, econ.ErrBadRecipe) {
		t.Fatalf("bad recipe: %v", err)
	}
	if err := w.BuildingSetRecipe(0, Point{}, Biomass); err != nil {
		t.Fatal(err)
	}
	if _, err := w.BuildingPut(0, Point{}, Moss, 1); !errors.Is(err, econ.ErrNoInput) {
		t.Fatalf("no inventory: %v", err)
	}
	w.Econ.Inv["iron"] = 1
	if _, err := w.BuildingPut(0, Point{}, Iron, 1); !errors.Is(err, econ.ErrBadRecipe) {
		t.Fatalf("wrong input: %v", err)
	}
	if _, err := w.BuildingTake(0, Point{}, Biomass, 1); !errors.Is(err, econ.ErrNoInput) {
		t.Fatalf("no output: %v", err)
	}
	b := w.Econ.Buildings[0]
	b.Enabled, b.In["moss"], b.Wear = false, 10, 30
	w.Econ.Inv["moss"] = 4
	if _, err := w.BuildingPut(0, Point{}, Moss, 1); !errors.Is(err, ErrFull) {
		t.Fatalf("full input: %v", err)
	}
	if err := w.BuildingRepair(0, Point{}); !errors.Is(err, econ.ErrNoInput) || b.Wear != 30 {
		t.Fatal("repair without plates changed wear")
	}
	w.Econ.Inv["plate"] = 2
	if err := w.BuildingRepair(0, Point{}); err != nil || b.Wear != 0 || w.Econ.Inv["plate"] != 0.5 {
		t.Fatalf("repair did not spend needed plates: %v, %+v", err, b)
	}
	if err := w.BuildingRepair(0, Point{}); err != nil || w.Econ.Inv["plate"] != 0.5 {
		t.Fatal("healthy repair spent plates")
	}
	for _, n := range []int{0, -1} {
		if moved, err := w.BuildingPut(0, Point{}, Moss, n); moved != 0 || err != nil {
			t.Fatal("nonpositive Put moved goods")
		}
		if moved, err := w.BuildingTake(0, Point{}, Biomass, n); moved != 0 || err != nil {
			t.Fatal("nonpositive Take moved goods")
		}
	}
}

func TestColonyBuildingProducesOnClockAndEmitsOnlyChanges(t *testing.T) {
	w, trace := e2World()
	w.Econ.Levels["composter"], w.Econ.Credits, w.Econ.Inv["moss"] = 1, 25, 12
	if err := w.Build(0, "composter"); err != nil {
		t.Fatal(err)
	}
	if err := w.BuildingSetRecipe(0, Point{}, Biomass); err != nil {
		t.Fatal(err)
	}
	if n, err := w.BuildingPut(0, Point{}, Moss, 8); n != 8 || err != nil {
		t.Fatalf("Put: %d, %v", n, err)
	}
	putAt := w.Tick
	w.Wait(0)
	v := w.BuildingStatus(Point{})
	if math.Abs(v.Progress-1.0/32) > 1e-9 || v.In["moss"] != 4 || v.Out["biomass"] != 0 {
		t.Fatalf("cycle start: %+v", v)
	}
	if err := w.BuildingSetRecipe(0, Point{}, Plate); !errors.Is(err, econ.ErrBadRecipe) {
		t.Fatal("active cycle changed recipe")
	}
	w.Tick = putAt + 30
	w.econSnapshot("econ", "", true, true)
	if w.BuildingStatus(Point{}).Out["biomass"] != 0 {
		t.Fatal("cycle finished before 32 production ticks")
	}
	if got := len(buildingEvents(t, trace)); got != 4 {
		t.Fatalf("progress emitted per-tick events: %d", got)
	}
	w.Tick = putAt + 32
	w.econSnapshot("econ", "", true, true)
	if w.BuildingStatus(Point{}).Out["biomass"] != 1 {
		t.Fatal("ticker did not complete cycle at level-1 duration")
	}
	events := buildingEvents(t, trace)
	if len(events) != 5 || events[0].Building.Kind != "composter" || events[2].Building.In["moss"] != 8 || events[3].Building.In["moss"] != 4 || events[4].Building.Out["biomass"] != 1 || events[4].T != putAt+32 {
		t.Fatalf("missing placement, transfer, consumption or completion events: %+v", events)
	}
	if n, err := w.BuildingTake(0, Point{}, Biomass, 1); n != 1 || err != nil {
		t.Fatal("cannot take produced biomass")
	}
	w.Econ.Levels["composter"] = 2
	w.applyEcon()
	if w.BuildingStatus(Point{}).Level != 2 {
		t.Fatal("building level did not follow owned tree level")
	}
	loaded := reloadColony(t, w)
	if !reflect.DeepEqual(loaded.BuildingStatus(Point{}), w.BuildingStatus(Point{})) {
		t.Fatal("mid-cycle state changed on reload")
	}
	w.Tick += 200
	loaded.Tick = w.Tick
	w.econSnapshot("econ", "", true, true)
	loaded.econSnapshot("econ", "", true, true)
	if !reflect.DeepEqual(loaded.BuildingStatus(Point{}), w.BuildingStatus(Point{})) || loaded.BuildingStatus(Point{}).Out["biomass"] != 1 {
		t.Fatal("mid-cycle save did not resume production")
	}
}

func TestColonyBuildingWireFormat(t *testing.T) {
	w, trace := e2World()
	w.Econ.Levels["composter"], w.Econ.Credits, w.Econ.Inv["moss"] = 1, 25, 4
	if err := w.Build(0, "composter"); err != nil {
		t.Fatal(err)
	}
	var event map[string]any
	line := bytes.Split(trace.Bytes(), []byte{'\n'})[0]
	if err := json.Unmarshal(line, &event); err != nil {
		t.Fatal(err)
	}
	var want map[string]any
	if err := json.Unmarshal([]byte(`{"kind":"composter","recipe":"","in":{},"out":{},"progress":0,"wear":0,"level":1}`), &want); err != nil {
		t.Fatal(err)
	}
	if event["op"] != "building" || event["t"] != float64(1) || event["x"] != float64(0) || event["y"] != float64(0) || !reflect.DeepEqual(event["building"], want) {
		t.Fatalf("unexpected building wire format: %s", line)
	}
	for _, x := range []int{0, 1} {
		blob, err := json.Marshal(w.view(x, 0))
		if err != nil {
			t.Fatal(err)
		}
		var tile map[string]any
		if err := json.Unmarshal(blob, &tile); err != nil {
			t.Fatal(err)
		}
		b, exists := tile["b"]
		if x == 0 && b != "composter" || x == 1 && exists {
			t.Fatalf("tile building field: %s", blob)
		}
	}
	trace.Reset()
	if err := w.Build(0, "missing"); !errors.Is(err, econ.ErrUnknownKind) {
		t.Fatal(err)
	}
	var narration Event
	if err := json.Unmarshal(bytes.TrimSpace(trace.Bytes()), &narration); err != nil {
		t.Fatal(err)
	}
	ar, en := econ.ErrUnknownKind.Messages()
	found := false
	for _, pair := range narration.Translations {
		found = found || pair == [2]string{ar, en}
	}
	if !found {
		t.Fatal("building error did not carry both languages to live subscribers")
	}
}
