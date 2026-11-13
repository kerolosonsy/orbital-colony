package econ

import "testing"

func TestSolarFollowsDaylightCycle(t *testing.T) {
	for _, tc := range []struct {
		tick int64
		day  bool
	}{
		{0, true}, {599, true}, {600, false}, {1199, false}, {1200, true}, {1799, true}, {1800, false},
		{-1, false}, {-600, false}, {-601, true}, {-1200, true},
	} {
		want := 0.0
		if tc.day {
			want = 10
		}
		if got := Daylight(tc.tick); got != tc.day {
			t.Errorf("Daylight(%d) = %v, want %v", tc.tick, got, tc.day)
		}
		if got := GeneratorOutput("solar", tc.tick); got != want {
			t.Errorf("solar at tick %d = %g, want %g", tc.tick, got, want)
		}
	}
}

func TestFueledGeneratorsProduceAtNightAndScaleFuelWithTicks(t *testing.T) {
	for _, tc := range []struct {
		kind, good string
		ticks      int
		kw, fuel   float64
	}{
		{"burner", "fuel", 100, 25, 1}, {"burner", "fuel", 50, 25, 0.5},
		{"reactor", "cell", 100, 100, 0.1}, {"reactor", "cell", 1000, 100, 1},
		{"burner", "fuel", 0, 25, 0}, {"reactor", "cell", -1, 100, 0},
		{"missing", "", 100, 0, 0},
	} {
		if got := GeneratorOutput(tc.kind, 600); got != tc.kw {
			t.Errorf("%s night output = %g, want %g", tc.kind, got, tc.kw)
		}
		if good, fuel := GeneratorFuel(tc.kind, tc.ticks); good != tc.good || fuel != tc.fuel {
			t.Errorf("%s over %d ticks: %s %g, want %s %g", tc.kind, tc.ticks, good, fuel, tc.good, tc.fuel)
		}
	}
	if good, fuel := GeneratorFuel("solar", 1200); good != "" || fuel != 0 {
		t.Fatalf("solar consumes %g %s", fuel, good)
	}
}

func TestPowerAllocationUsesPrioritySkipsLargeLoadsAndCountsDemand(t *testing.T) {
	for _, tc := range []struct {
		name      string
		supply    float64
		buildings []*Building
		powered   map[int]bool
		demand    float64
	}{
		{"skip oversized", 5, []*Building{
			{Kind: "reactor-lab", Priority: 9, Enabled: true},
			{Kind: "smelter", Priority: 3, Enabled: true},
			{Kind: "composter", Priority: 3, Enabled: true},
		}, map[int]bool{1: true, 2: true}, 25},
		{"ties use index", 4, []*Building{
			{Kind: "smelter", Priority: 3, Enabled: true},
			{Kind: "composter", Priority: 3, Enabled: true},
		}, map[int]bool{0: true}, 5},
		{"later priority wins", 4, []*Building{
			{Kind: "smelter", Priority: 1, Enabled: true},
			{Kind: "composter", Priority: 9, Enabled: true},
		}, map[int]bool{1: true}, 5},
		{"disabled load ignored", 4, []*Building{
			{Kind: "reactor-lab", Priority: 9, Enabled: false},
			{Kind: "smelter", Priority: 1, Enabled: true},
		}, map[int]bool{1: true}, 4},
		{"passive buildings", 0, []*Building{
			{Kind: "silo", Enabled: true},
			{Kind: "generator", Enabled: true},
			{Kind: "lab", Enabled: true},
		}, map[int]bool{0: true, 1: true}, 6},
		{"all consumers", 53, []*Building{
			{Kind: "composter", Enabled: true}, {Kind: "smelter", Enabled: true},
			{Kind: "kiln", Enabled: true}, {Kind: "refinery", Enabled: true},
			{Kind: "assembler", Enabled: true}, {Kind: "reactor-lab", Enabled: true}, {Kind: "lab", Enabled: true},
		}, map[int]bool{0: true, 1: true, 2: true, 3: true, 4: true, 5: true, 6: true}, 49},
		{"empty", 10, nil, map[int]bool{}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			powered, demand := Allocate(tc.supply, tc.buildings)
			if demand != tc.demand {
				t.Errorf("demand = %g, want %g", demand, tc.demand)
			}
			for i := range tc.buildings {
				if powered[i] != tc.powered[i] {
					t.Errorf("building %d powered=%v, want %v", i, powered[i], tc.powered[i])
				}
			}
			for i, on := range powered {
				if on && !tc.powered[i] {
					t.Errorf("unexpected powered building %d", i)
				}
			}
		})
	}
}
