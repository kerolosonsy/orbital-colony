package econ

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"testing"
)

func TestLegacyStorageKeepsAllHoldingsAndDoesNotGrowOnReload(t *testing.T) {
	for _, version := range []int{1, 2, 3} {
		s := New(11)
		s.Version, s.StorageBase = version, nil
		s.Inv = map[string]float64{"moss": 4000.5, "lichen": 2, "iron": 2000, "water": 3.5, "cell": 3000}
		s.Buildings = []*Building{{Kind: "silo", Enabled: true}}
		before := map[string]float64{}
		for good, n := range s.Inv {
			before[good] = n
		}
		s.Normalize()
		for group, st := range s.Storage {
			if st.Used > st.Cap {
				t.Fatalf("v%d %s over capacity: %+v", version, group, st)
			}
		}
		if !reflect.DeepEqual(before, s.Inv) {
			t.Fatal("migration lost goods")
		}
		capacity := s.Storage["crops"].Cap
		s.Inv["moss"] = float64(capacity + 1)
		blob, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		var loaded State
		if err := json.Unmarshal(blob, &loaded); err != nil {
			t.Fatal(err)
		}
		loaded.Normalize()
		if loaded.Storage["crops"].Cap != capacity {
			t.Fatal("reload raised an initialized storage limit")
		}
	}
}

func TestStorageRejectsWholeTransferAndSilosExpandEachGroup(t *testing.T) {
	s := New(1)
	s.StorageBase = map[string]int{"crops": 2, "raw": 2, "crafted": 2}
	s.Inv["moss"], s.Inv["water"] = 1.5, 2
	if err := s.Store(map[string]float64{"moss": .5}); err != nil {
		t.Fatal(err)
	}
	if err := s.Store(map[string]float64{"water": 1, "cell": 1}); !errors.Is(err, ErrFull) || s.Inv["water"] != 2 || s.Inv["cell"] != 0 {
		t.Fatal("overflow partially transferred goods")
	}
	s.Buildings = []*Building{{Kind: "silo"}, {Kind: "silo"}}
	if err := s.Store(map[string]float64{"moss": 400, "water": 200, "cell": 202}); err != nil {
		t.Fatal(err)
	}
	for group, st := range s.Storage {
		if st.Used != st.Cap {
			t.Fatalf("%s silo limit: %+v", group, st)
		}
	}
}

func TestGeneratorFuelExhaustionAndModeChangesKeepBuffers(t *testing.T) {
	for _, mode := range []string{"solar", "burner", "reactor"} {
		s := New(1)
		b := &Building{Kind: "generator", Enabled: true, In: map[string]float64{}}
		s.Buildings = []*Building{b}
		if err := b.SetMode(mode); err != nil {
			t.Fatal(err)
		}
		good, fuel := GeneratorFuel(mode, 1)
		if fuel > 0 {
			if s.ConsumePower(600) != 0 {
				t.Fatal("unfueled generator produced")
			}
			b.In[good] = fuel * 3
			for i := 0; i < 3; i++ {
				if s.ConsumePower(600) != GeneratorOutput(mode, 600) {
					t.Fatal("fuel exhausted too soon")
				}
			}
			if s.ConsumePower(600) != 0 || math.Abs(b.In[good]) > 1e-9 {
				t.Fatal("fuel exhausted too late")
			}
		}
		b.Enabled = false
		if s.ConsumePower(0) != 0 {
			t.Fatal("disabled generator produced")
		}
		b.In["fuel"] = 1
		if err := b.SetMode("invalid"); !errors.Is(err, ErrBadRecipe) || b.Mode != mode || b.In["fuel"] != 1 {
			t.Fatal("invalid mode changed generator")
		}
	}
}

func TestRentChargesOnlyRunningTicksAndStopsAtZero(t *testing.T) {
	s := New(1)
	s.Buildings = []*Building{{Kind: "silo"}, {Kind: "composter"}}
	s.Credits = 10
	for i := 0; i < 600; i++ {
		if !s.PayRent(2) {
			t.Fatal("rent ran out early")
		}
	}
	if math.Abs(s.Credits-8.5) > 1e-9 || s.Rent != 1.5 {
		t.Fatalf("minute rent: %+v", s)
	}
	s.Credits = .001
	if s.PayRent(2) || s.Credits != 0 {
		t.Fatal("unaffordable rent went negative or ran buildings")
	}
	s.Inv["moss"] = 1
	s.Sell("moss", 1)
	if !s.PayRent(2) {
		t.Fatal("credits did not restore service")
	}
}

func TestFractionalStorageAndExactRentDoNotOverflowOnRounding(t *testing.T) {
	s := New(1)
	s.StorageBase["raw"] = 1
	for i := 0; i < 10; i++ {
		if err := s.Store(map[string]float64{"water": .1}); err != nil {
			t.Fatal(err)
		}
	}
	if s.Storage["raw"].Used != 1 {
		t.Fatalf("rounded storage: %+v", s.Storage["raw"])
	}
	s.Buildings = []*Building{{Kind: "silo"}}
	s.Credits = .5
	for i := 0; i < 600; i++ {
		if !s.PayRent(1) {
			t.Fatalf("rent exhausted at tick %d", i)
		}
	}
	if s.PayRent(1) || s.Credits != 0 {
		t.Fatal("exact minute rent was not exhausted")
	}
}
