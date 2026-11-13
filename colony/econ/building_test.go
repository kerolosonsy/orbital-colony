package econ

import (
	"errors"
	"maps"
	"reflect"
	"testing"

	"orbital/colony/i18n"
)

func TestBuildingStopsWithoutChangingBuffers(t *testing.T) {
	for _, tc := range []struct {
		name   string
		recipe string
		iron   float64
		out    float64
		wear   float64
		power  bool
		want   error
	}{
		{"missing input", "plate", 1, 0, 0, true, ErrNoInput},
		{"full output", "plate", 2, 8, 0, true, ErrFull},
		{"less than one output slot", "plate", 2, 7.5, 0, true, ErrFull},
		{"no power", "plate", 2, 0, 0, false, ErrNoPower},
		{"worn out", "plate", 2, 0, 100, true, ErrWorn},
		{"unknown recipe", "missing", 2, 0, 0, true, ErrBadRecipe},
		{"wrong building", "fuel", 2, 0, 0, true, ErrBadRecipe},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := Building{Kind: "smelter", Recipe: tc.recipe, In: map[string]float64{"iron": tc.iron}, Out: map[string]float64{"plate": tc.out}, Wear: tc.wear, Enabled: true}
			before := b
			before.In, before.Out = maps.Clone(b.In), maps.Clone(b.Out)
			if err := b.Step(1, tc.power); !errors.Is(err, tc.want) {
				t.Fatalf("Step error = %v, want %v", err, tc.want)
			}
			if !reflect.DeepEqual(b, before) {
				t.Fatalf("stalled building changed: %+v, before %+v", b, before)
			}
		})
	}
}

func TestInactiveBuildingPreservesProgress(t *testing.T) {
	for _, tc := range []struct {
		ticks   int
		enabled bool
	}{{0, true}, {-1, true}, {100, false}} {
		b := Building{Kind: "smelter", Recipe: "plate", Progress: 10, Enabled: tc.enabled}
		before := b
		if err := b.Step(tc.ticks, false); err != nil || !reflect.DeepEqual(b, before) {
			t.Fatalf("ticks=%d enabled=%v: building %+v, error %v", tc.ticks, tc.enabled, b, err)
		}
	}
}

func TestBuildingSpendsInputsAtStartAndResumesAfterPowerLoss(t *testing.T) {
	b := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 2}, Enabled: true}
	if err := b.Step(1, true); err != nil || b.In["iron"] != 0 || b.Progress != 1 || b.Out["plate"] != 0 || b.Wear != 0 {
		t.Fatalf("cycle start: %+v, error %v", b, err)
	}
	if err := b.Step(100, false); err != ErrNoPower || b.Progress != 1 {
		t.Fatalf("power loss: %+v, error %v", b, err)
	}
	if err := b.Step(59, true); err != nil || b.Out["plate"] != 1 || b.Progress != 0 || b.Wear != 1 {
		t.Fatalf("resumed cycle: %+v, error %v", b, err)
	}
}

func TestBuildingCompletesMultipleCyclesAndKeepsRemainder(t *testing.T) {
	b := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 8}, Enabled: true}
	if err := b.Step(185, true); err != nil || b.Out["plate"] != 3 || b.Progress != 5 || b.In["iron"] != 0 || b.Wear != 3 {
		t.Fatalf("three cycles and a partial cycle: %+v, error %v", b, err)
	}
	if err := b.Step(55, true); err != nil || b.Out["plate"] != 4 || b.Progress != 0 || b.Wear != 4 {
		t.Fatalf("final cycle: %+v, error %v", b, err)
	}
	if err := b.Step(1, true); err != ErrNoInput {
		t.Fatalf("exhausted input: error %v", err)
	}
}

func TestBuildingBatchMatchesSingleTicks(t *testing.T) {
	for _, tc := range []struct {
		name  string
		level int
		wear  float64
		out   float64
	}{
		{"ordinary", 0, 0, 0},
		{"upgraded", 16, 0, 0},
		{"crossing slowdown", 1, 80, 0},
		{"crossing shutdown", 1, 99, 0},
		{"filling mixed output", 0, 0, 6},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, ticks := range []int{1, 19, 20, 39, 40, 47, 48, 59, 60, 61, 145, 241, 5000} {
				batch := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 8}, Out: map[string]float64{"wire": tc.out}, Level: tc.level, Wear: tc.wear, Enabled: true}
				single := batch
				single.In, single.Out = maps.Clone(batch.In), maps.Clone(batch.Out)
				batchErr := batch.Step(ticks, true)
				var singleErr error
				for tick := 0; tick < ticks; tick++ {
					if singleErr = single.Step(1, true); singleErr != nil {
						break
					}
				}
				if batchErr != singleErr || !reflect.DeepEqual(batch, single) {
					t.Fatalf("%d ticks: batch %+v (%v), single %+v (%v)", ticks, batch, batchErr, single, singleErr)
				}
			}
		})
	}
}

func TestBuildingWearSlowsAbove80AndStopsAt100(t *testing.T) {
	b := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 4}, Wear: 80, Enabled: true}
	if err := b.Step(60, true); err != nil || b.Out["plate"] != 1 || b.Wear != 81 {
		t.Fatalf("wear 80 should run at full speed: %+v, error %v", b, err)
	}
	if err := b.Step(119, true); err != nil || b.Out["plate"] != 1 || b.Progress != 119 {
		t.Fatalf("wear 81 should need 120 ticks: %+v, error %v", b, err)
	}
	if err := b.Step(1, true); err != nil || b.Out["plate"] != 2 || b.Wear != 82 {
		t.Fatalf("slower cycle: %+v, error %v", b, err)
	}
	b.Wear, b.In["iron"] = 99, 4
	if err := b.Step(121, true); err != ErrWorn || b.Wear != 100 || b.Out["plate"] != 3 || b.In["iron"] != 2 || b.Progress != 0 {
		t.Fatalf("shutdown must finish one cycle and retain next inputs: %+v, error %v", b, err)
	}
}

func TestBuildingLevelSpeedsUpWholeTickCycles(t *testing.T) {
	for _, tc := range []struct{ level, ticks int }{{0, 60}, {1, 48}, {16, 20}, {10000, 1}} {
		b := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 2}, Level: tc.level, Enabled: true}
		if err := b.Step(tc.ticks-1, true); err != nil || b.Out["plate"] != 0 {
			t.Fatalf("level %d finished early: %+v, error %v", tc.level, b, err)
		}
		if err := b.Step(1, true); err != nil || b.Out["plate"] != 1 {
			t.Fatalf("level %d did not finish: %+v, error %v", tc.level, b, err)
		}
	}
}

func TestBuildingChecksAllInputsBeforeConsumingAny(t *testing.T) {
	b := Building{Kind: "refinery", Recipe: "fuel", In: map[string]float64{"biomass": 2}, Enabled: true}
	if err := b.Step(100, true); err != ErrNoInput || b.In["biomass"] != 2 || b.Progress != 0 {
		t.Fatalf("incomplete inputs were spent: %+v, error %v", b, err)
	}
	if _, err := b.Put("water", 1); err != nil {
		t.Fatal(err)
	}
	if err := b.Step(100, true); err != nil || b.Out["fuel"] != 1 || b.In["biomass"] != 0 || b.In["water"] != 0 {
		t.Fatalf("complete recipe: %+v, error %v", b, err)
	}
}

func TestBuildingResumesAfterOutputIsDrained(t *testing.T) {
	b := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 2}, Out: map[string]float64{"wire": 8}, Enabled: true}
	if err := b.Step(60, true); err != ErrFull || b.In["iron"] != 2 {
		t.Fatalf("full output spent inputs: %+v, error %v", b, err)
	}
	if n, err := b.Take("wire", 1); n != 1 || err != nil {
		t.Fatalf("draining output: took %g, error %v", n, err)
	}
	if err := b.Step(60, true); err != nil || b.Out["plate"] != 1 || b.Out["wire"] != 7 || b.In["iron"] != 0 {
		t.Fatalf("drained building did not resume: %+v, error %v", b, err)
	}
}

func TestBuildingRecipeChangesPreserveBuffersAndActiveCycle(t *testing.T) {
	b := Building{Kind: "smelter", In: map[string]float64{"iron": 2}, Out: map[string]float64{"plate": 1}, Enabled: true}
	for _, out := range []string{"missing", "fuel"} {
		if err := b.SetRecipe(out); err != ErrBadRecipe || b.Recipe != "" {
			t.Fatalf("invalid recipe %q: %+v, error %v", out, b, err)
		}
	}
	if err := b.SetRecipe("plate"); err != nil {
		t.Fatal(err)
	}
	if err := b.Step(1, true); err != nil {
		t.Fatal(err)
	}
	if err := b.SetRecipe("wire"); err != ErrBadRecipe || b.Recipe != "plate" || b.Progress != 1 {
		t.Fatalf("active recipe changed: %+v, error %v", b, err)
	}
	if err := b.SetRecipe("plate"); err != nil || b.Progress != 1 {
		t.Fatalf("same recipe lost progress: %+v, error %v", b, err)
	}
	if err := b.Step(59, true); err != nil {
		t.Fatal(err)
	}
	if err := b.SetRecipe("wire"); err != nil || b.Out["plate"] != 2 || b.In["iron"] != 0 {
		t.Fatalf("idle recipe change lost buffers: %+v, error %v", b, err)
	}
}

func TestBuildingPutUsesPerGoodCapacityAndRejectsOtherGoods(t *testing.T) {
	b := Building{Kind: "refinery", Recipe: "fuel"}
	if n, err := b.Put("biomass", 3); n != 3 || err != nil {
		t.Fatalf("empty buffer: accepted %g, error %v", n, err)
	}
	if n, err := b.Put("water", 9); n != 8 || err != nil {
		t.Fatalf("water slot capacity: accepted %g, error %v", n, err)
	}
	if n, err := b.Put("water", 1); n != 0 || err != ErrFull {
		t.Fatalf("full water slot: accepted %g, error %v", n, err)
	}
	if n, err := b.Put("biomass", 1); n != 1 || err != nil {
		t.Fatalf("independent biomass slot: accepted %g, error %v", n, err)
	}
	if n, err := b.Put("iron", 1); n != 0 || err != ErrBadRecipe || b.In["iron"] != 0 {
		t.Fatalf("foreign good: accepted %g, error %v", n, err)
	}
	for _, n := range []float64{0, -1} {
		if accepted, err := b.Put("water", n); accepted != 0 || err != nil || b.In["water"] != 8 {
			t.Fatalf("nonpositive put %g: accepted %g, error %v", n, accepted, err)
		}
	}
	b.Level = 1
	if n, err := b.Put("water", 3); n != 2 || err != nil {
		t.Fatalf("upgraded capacity: accepted %g, error %v", n, err)
	}
}

func TestComposterUsesFirstAvailableRecipe(t *testing.T) {
	for _, tc := range []struct {
		name          string
		moss, lichen  float64
		remainingMoss float64
		remainingLich float64
	}{
		{"lichen only", 0, 2, 0, 0},
		{"primary first", 4, 2, 0, 2},
		{"primary incomplete", 3, 2, 3, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b := Building{Kind: "composter", Enabled: true}
			if err := b.SetRecipe("biomass"); err != nil {
				t.Fatal(err)
			}
			for good, amount := range map[string]float64{"moss": tc.moss, "lichen": tc.lichen} {
				if n, err := b.Put(good, amount); n != amount || err != nil {
					t.Fatalf("Put(%s, %g) = %g, %v", good, amount, n, err)
				}
			}
			if err := b.Step(40, true); err != nil || b.Out["biomass"] != 1 || b.In["moss"] != tc.remainingMoss || b.In["lichen"] != tc.remainingLich {
				t.Fatalf("recipe choice: %+v, error %v", b, err)
			}
		})
	}
}

func TestComposterBatchSelectsAgainAfterEachCycle(t *testing.T) {
	b := Building{Kind: "composter", Recipe: "biomass", In: map[string]float64{"moss": 4, "lichen": 2}, Enabled: true}
	if err := b.Step(80, true); err != nil || b.Out["biomass"] != 2 || b.In["moss"] != 0 || b.In["lichen"] != 0 || b.Progress != 0 {
		t.Fatalf("primary then alternative: %+v, error %v", b, err)
	}
}

func TestComposterKeepsStartedAlternativeWhenPrimaryInputsArrive(t *testing.T) {
	previous := Recipes
	Recipes = append([]Recipe(nil), Recipes...)
	t.Cleanup(func() { Recipes = previous })
	// Distinct durations make a mid-cycle recipe switch observable.
	for i := range Recipes {
		if Recipes[i].Out == "biomass" && Recipes[i].Alt {
			Recipes[i].Ticks = 80
		}
	}
	b := Building{Kind: "composter", Recipe: "biomass", In: map[string]float64{"lichen": 2}, Enabled: true}
	if err := b.Step(1, true); err != nil || b.In["lichen"] != 0 {
		t.Fatalf("alternative start: %+v, error %v", b, err)
	}
	if n, err := b.Put("moss", 4); n != 4 || err != nil {
		t.Fatalf("primary inputs: accepted %g, error %v", n, err)
	}
	if err := b.Step(39, true); err != nil || b.Progress != 40 || b.Out["biomass"] != 0 || b.In["moss"] != 4 {
		t.Fatalf("active alternative changed: %+v, error %v", b, err)
	}
	if err := b.Step(40, true); err != nil || b.Progress != 0 || b.Out["biomass"] != 1 || b.In["moss"] != 4 {
		t.Fatalf("alternative completion: %+v, error %v", b, err)
	}
	if err := b.Step(40, true); err != nil || b.Out["biomass"] != 2 || b.In["moss"] != 0 {
		t.Fatalf("next cycle should use primary: %+v, error %v", b, err)
	}
}

func TestAssemblerFullWireSlotStillAcceptsOtherCircuitInputs(t *testing.T) {
	b := Building{Kind: "assembler", Enabled: true}
	if err := b.SetRecipe("circuit"); err != nil {
		t.Fatal(err)
	}
	for _, input := range []struct {
		good   string
		amount float64
	}{{"wire", 8}, {"glass", 1}, {"crystal", 1}} {
		if n, err := b.Put(input.good, input.amount); n != input.amount || err != nil {
			t.Fatalf("Put(%s, %g) = %g, %v", input.good, input.amount, n, err)
		}
	}
	if n, err := b.Put("wire", 1); n != 0 || err != ErrFull {
		t.Fatalf("full wire slot: accepted %g, error %v", n, err)
	}
	if err := b.Step(200, true); err != nil || b.Out["circuit"] != 1 || b.In["wire"] != 6 || b.In["glass"] != 0 || b.In["crystal"] != 0 {
		t.Fatalf("circuit production: %+v, error %v", b, err)
	}
}

func TestBuildingTakeOnlyRemovesAvailableOutput(t *testing.T) {
	b := Building{In: map[string]float64{"iron": 2}, Out: map[string]float64{"plate": 2.5}}
	for _, tc := range []struct {
		good string
		n    float64
		want float64
		err  error
	}{{"plate", 2, 2, nil}, {"plate", 2, 0.5, nil}, {"plate", 1, 0, ErrNoInput}, {"iron", 1, 0, ErrNoInput}, {"plate", 0, 0, nil}, {"plate", -1, 0, nil}} {
		if n, err := b.Take(tc.good, tc.n); n != tc.want || err != tc.err || b.In["iron"] != 2 {
			t.Fatalf("Take(%s, %g) = %g, %v; want %g, %v", tc.good, tc.n, n, err, tc.want, tc.err)
		}
	}
}

func TestBuildingCapacityGrowsWithLevel(t *testing.T) {
	if Capacity(0) != 8 || Capacity(1) != 10 || Capacity(16) != 24 {
		t.Fatal("capacity does not use the chosen diminishing returns")
	}
	for level := 1; level <= 500; level++ {
		if Capacity(level) <= Capacity(level-1) {
			t.Fatalf("capacity did not grow at level %d", level)
		}
	}
}

func TestBuildingRepairUsesOnlyNeededPlates(t *testing.T) {
	for _, tc := range []struct{ wear, plates, used, remaining float64 }{
		{100, 1, 1, 80}, {10, 2, 0.5, 0}, {80, 0.25, 0.25, 75}, {0, 1, 0, 0}, {50, 0, 0, 50}, {50, -1, 0, 50},
	} {
		b := Building{Wear: tc.wear}
		if used := b.Repair(tc.plates); used != tc.used || b.Wear != tc.remaining {
			t.Fatalf("wear %g, plates %g: used %g, wear %g; want %g, %g", tc.wear, tc.plates, used, b.Wear, tc.used, tc.remaining)
		}
	}
}

func TestBuildingRepairRestartsWornProductionAtFullSpeed(t *testing.T) {
	b := Building{Kind: "smelter", Recipe: "plate", In: map[string]float64{"iron": 2}, Wear: 100, Enabled: true}
	if used := b.Repair(1); used != 1 || b.Wear != 80 {
		t.Fatalf("repair: used %g, wear %g", used, b.Wear)
	}
	if err := b.Step(60, true); err != nil || b.Out["plate"] != 1 || b.Wear != 81 {
		t.Fatalf("repaired production: %+v, error %v", b, err)
	}
}

func TestBuildingErrorsFollowCurrentLanguage(t *testing.T) {
	previous := i18n.Lang()
	t.Cleanup(func() { i18n.SetLang(previous) })
	for _, tc := range []struct {
		err    error
		ar, en string
	}{
		{ErrNoInput, "مفيش مدخلات كفاية", "Not enough inputs"},
		{ErrFull, "المخزن مليان", "Buffer is full"},
		{ErrNoPower, "مفيش طاقة", "No power"},
		{ErrWorn, "المبنى محتاج تصليح", "Building needs repair"},
		{ErrBadRecipe, "الوصفة مش صالحة للمبنى دلوقتي", "Recipe is not valid for this building right now"},
	} {
		i18n.SetLang("en")
		if got := tc.err.Error(); got != tc.en {
			t.Errorf("English error = %q, want %q", got, tc.en)
		}
		i18n.SetLang("ar")
		if got := tc.err.Error(); got != tc.ar {
			t.Errorf("Arabic error = %q, want %q", got, tc.ar)
		}
	}
}
