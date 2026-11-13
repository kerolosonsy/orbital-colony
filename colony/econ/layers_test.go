package econ

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestLayerSaveNormalizationAndExpansion(t *testing.T) {
	// A legacy save has no depth or underground data.
	legacy := State{Seed: 42}
	legacy.Normalize()
	legacy.AtZ(0, 0, 1).G, legacy.AtZ(0, 0, 1).E = "shaft", ""
	legacy.AtZ(0, 0, 2).G, legacy.AtZ(0, 0, 2).E = "tunnel", ""
	encoded, err := json.Marshal(&legacy)
	if err != nil {
		t.Fatal(err)
	}
	var loaded State
	if err := json.Unmarshal(encoded, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.Normalize()
	if !reflect.DeepEqual(legacy.Under, loaded.Under) || legacy.Depth != loaded.Depth {
		t.Fatal("normalization changed saved underground tiles")
	}
	loaded.Levels["field"] = 2
	loaded.Normalize()
	direct := New(42)
	direct.Levels["field"] = 2
	direct.Normalize()
	for z := 1; z < loaded.Depth; z++ {
		if len(loaded.Under[z-1]) != loaded.Size*loaded.Size {
			t.Fatal("underground layer did not grow with field")
		}
		for y := 0; y < loaded.Size; y++ {
			for x := 0; x < loaded.Size; x++ {
				if x >= legacy.Size || y >= legacy.Size {
					if *loaded.AtZ(x, y, z) != *direct.AtZ(x, y, z) {
						t.Fatal("new ore depends on expansion history")
					}
				}
				if z == 1 && loaded.AtZ(x, y, z).E == "geode" {
					t.Fatal("geodes generated above layer 2")
				}
			}
		}
	}
	if loaded.AtZ(0, 0, 1).G != "shaft" || loaded.AtZ(0, 0, 2).G != "tunnel" {
		t.Fatal("expansion replaced existing excavation")
	}
}

func TestColonyDrillPurchaseAndOreMarket(t *testing.T) {
	s := New(1)
	s.Credits, s.RP = 1000000, 1000
	s.Levels["field"] = 1
	s.Blueprints["api.dig"] = true
	if err := s.Buy("api.dig"); err != nil {
		t.Fatal(err)
	}
	if err := s.Buy("drill"); err != nil {
		t.Fatal(err)
	}
	if err := s.Buy("drill"); err == nil {
		t.Fatal("second drill level did not require iron")
	}
	s.Inv["iron"] = 10
	if err := s.Buy("drill"); err != nil || s.Level("drill") != 2 || s.Inv["iron"] != 0 {
		t.Fatalf("drill purchase did not spend iron: %v", err)
	}
	for _, ore := range []string{"iron", "geode"} {
		s.Inv[ore] = 1
		if s.Prices()[ore] == 0 || s.Sell(ore, 1) <= 0 || s.Inv[ore] != 0 {
			t.Fatalf("ore %s cannot be sold", ore)
		}
		s.Inv[ore] = 5
		if err := s.Buy("yield." + ore); err != nil || s.YieldMult(ore) <= 1 {
			t.Fatalf("ore %s yield upgrade failed: %v", ore, err)
		}
	}
}
