package sim

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"orbital/colony/econ"
)

func TestLayerDigCollectAndReturn(t *testing.T) {
	var trace bytes.Buffer
	w := New(&trace, 1)
	w.Init(2, 2, Regolith)
	w.InitLayers(3)
	w.AtZ(0, 0, 1).Entity = Iron
	w.AtZ(1, 0, 1).Entity = Iron
	w.AtZ(0, 0, 2).Entity = Geode
	if w.Dig(0, East) || w.Ascend(0) {
		t.Fatal("surface digging or ascending succeeded")
	}
	if !w.DigDown(0) || !w.Descend(0) || w.Layer(0) != 1 {
		t.Fatal("could not open shaft and descend")
	}
	if w.Move(0, East) || w.X != 0 || w.Count(Iron) != 1 {
		t.Fatal("stone was walkable or shaft ore was lost")
	}
	if !w.Dig(0, East) || !w.Move(0, East) || w.Count(Iron) != 2 {
		t.Fatal("tunnel digging did not collect iron or open movement")
	}
	if w.GetGround(0) != Tunnel || w.DroneGrid(0)[0][1] != Stone || w.At(1, 0).Ground != Regolith {
		t.Fatal("sensor/grid changed or read the surface instead of the current layer")
	}
	if !w.Move(0, West) || !w.DigDown(0) || w.Count(Geode) != 3 || !w.Descend(0) {
		t.Fatal("second shaft did not collect three geodes")
	}
	if w.DigDown(0) || !w.Ascend(0) || !w.Ascend(0) || w.Layer(0) != 0 {
		t.Fatal("depth boundary or return to surface failed")
	}
	if w.Tick != 14 || w.LayerDepth() != 3 {
		t.Fatalf("expected 14 action ticks and depth 3, got %d ticks and depth %d", w.Tick, w.LayerDepth())
	}
	var sawOreTile, sawDescend bool
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "tile" && ev.Tile.X == 1 && ev.Tile.Z == 1 && ev.Tile.G == "tunnel" {
			sawOreTile = true
		}
		if ev.Op == "descend" && ev.OK && ev.Z == 2 {
			sawDescend = true
		}
	}
	if !sawOreTile || !sawDescend {
		t.Fatal("missing layer coordinates in tile/action events")
	}
}

func TestLayerInitSnapshot(t *testing.T) {
	for _, depth := range []int{1, 3} {
		var trace bytes.Buffer
		w := New(&trace, 1)
		w.Init(2, 2, Regolith)
		w.InitLayers(depth)
		w.Start()
		var ev struct {
			Layers int `json:"layers"`
			Tiles  []struct {
				X, Y, Z int
				G       string
			} `json:"tiles"`
		}
		if err := json.Unmarshal(bytes.TrimSpace(trace.Bytes()), &ev); err != nil {
			t.Fatal(err)
		}
		if len(ev.Tiles) != 4*depth || depth > 1 && ev.Layers != depth {
			t.Fatalf("depth %d: snapshot has %d tiles, layers=%d", depth, len(ev.Tiles), ev.Layers)
		}
		seen := map[[3]int]bool{}
		for _, tile := range ev.Tiles {
			seen[[3]int{tile.X, tile.Y, tile.Z}] = true
			if tile.Z > 0 && tile.G != "stone" {
				t.Fatal("underground not solid stone")
			}
		}
		if len(seen) != 4*depth {
			t.Fatal("duplicate tile coordinates")
		}
		if depth == 1 && (bytes.Contains(trace.Bytes(), []byte(`"layers"`)) || bytes.Contains(trace.Bytes(), []byte(`"z"`))) {
			t.Fatal("surface-only snapshot changed its JSON format")
		}
	}
}

func TestLayerDroneStateIsIndependent(t *testing.T) {
	w := New(io.Discard, 1)
	w.Init(2, 2, Soil)
	w.InitLayers(2)
	w.MaxDrones = 2
	id := w.Spawn()
	if !w.DigDown(id) || !w.Descend(id) || !w.Dig(id, East) || !w.Move(id, East) {
		t.Fatal("second drone could not explore underground")
	}
	if w.Layer(0) != 0 || w.Layer(id) != 1 || w.X != 0 || w.Drones[id].X != 1 {
		t.Fatal("actions leaked position between drones")
	}
	if w.DroneGrid(id)[0][1] != Stone || w.GroundGrid()[0][1] != Soil {
		t.Fatal("drone grid used another drone's layer")
	}
	w.AtZ(1, 0, 1).Entity, w.AtZ(1, 0, 1).Size = Crystal, 4
	if w.Measure(id) != 4 || !w.Swap(id, West) || w.AtZ(0, 0, 1).Entity != Crystal || w.At(0, 0).Entity != Empty {
		t.Fatal("measure or swap used the surface")
	}
}

func TestColonyLayerGatesPersistenceAndResize(t *testing.T) {
	var trace bytes.Buffer
	st := econ.New(42)
	w := NewColony(&trace, st, 42)
	w.Realtime = false
	st.AtZ(0, 0, 1).E = "iron"
	w.loadTiles()
	for _, action := range []func() bool{
		func() bool { return w.DigDown(0) }, func() bool { return w.Descend(0) },
		func() bool { return w.Ascend(0) }, func() bool { return w.Dig(0, East) },
	} {
		if action() {
			t.Fatal("dig API worked without licence")
		}
	}
	st.Levels["api.dig"] = 1
	if w.DigDown(0) {
		t.Fatal("digging worked without drill")
	}
	st.Levels["drill"], st.Levels["yield.iron"] = 1, 2
	if !w.DigDown(0) || !w.Descend(0) || st.Inv["iron"] != 1.7 || w.DigDown(0) {
		t.Fatal("drill limit or ore yield multiplier failed")
	}
	st.Levels["drill"] = 2
	if !w.DigDown(0) || !w.Descend(0) {
		t.Fatal("drill upgrade did not allow the next layer")
	}
	// Buy through the live command path so resize events and underground drones are exercised.
	st.Credits = 100000
	w.colonyCommands(strings.NewReader(`{"buy":"field"}` + "\n"))
	if w.Layer(0) != 2 || w.AtZ(0, 0, 1).Ground != Shaft || w.AtZ(0, 0, 2).Ground != Tunnel {
		t.Fatal("field expansion lost tunnels or drone depth")
	}
	var resized bool
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev econEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "resize" {
			resized = true
			if ev.Layers != st.Depth || len(ev.Tiles) != st.Size*st.Size*st.Depth {
				t.Fatal("resize did not include all layers")
			}
		}
	}
	if !resized {
		t.Fatal("missing resize event")
	}
	w.syncEcon()
	saved, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var loaded econ.State
	if err := json.Unmarshal(saved, &loaded); err != nil {
		t.Fatal(err)
	}
	restored := NewColony(io.Discard, &loaded, 42)
	if restored.Z != 0 || restored.AtZ(0, 0, 1).Ground != Shaft || restored.AtZ(0, 0, 2).Ground != Tunnel || restored.Inv[Iron] != 1 {
		t.Fatal("saved colony did not restore tunnels, shafts, ore, and surface start")
	}
}

func TestLayerUndergroundSetupUsesSeedAndOreDepth(t *testing.T) {
	newWorld := func() *World {
		w := New(io.Discard, 9)
		w.Init(8, 8, Soil)
		w.Underground(3, 0.5, 0.5)
		return w
	}
	a, b := newWorld(), newWorld()
	iron, geode := 0, 0
	for z := 0; z < a.Depth; z++ {
		for y := 0; y < a.H; y++ {
			for x := 0; x < a.W; x++ {
				tile := a.AtZ(x, y, z)
				if *tile != *b.AtZ(x, y, z) {
					t.Fatal("same seed generated different mission ore")
				}
				if z == 0 && (tile.Ground != Soil || tile.Entity != Empty) || z > 0 && tile.Ground != Stone || z == 1 && tile.Entity == Geode {
					t.Fatal("ore setup changed surface or violated underground rules")
				}
				if tile.Entity == Iron {
					iron++
				}
				if tile.Entity == Geode {
					geode++
				}
			}
		}
	}
	if iron == 0 || geode == 0 {
		t.Fatal("mission helper did not generate both ores")
	}
}
