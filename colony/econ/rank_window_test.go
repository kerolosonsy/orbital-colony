package econ

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"orbital/colony/i18n"
)

func TestRankWindowBucketsAndEvictsSalesAndCommands(t *testing.T) {
	var window RankWindow
	window.Advance(1)
	window.RecordSale("moss", 10)
	window.RecordCommand()
	window.Advance(99)
	window.RecordSale("moss", 2)
	window.RecordCommand()
	window.Advance(1)
	window.RecordSale("plate", 1)
	window.RecordCommand()
	window.Advance(Window - 101)
	sold, commands := window.totals()
	if Throughput(sold, Window) != 8.2 || commands != 3 {
		t.Fatalf("full window lost a bucket: sold=%v commands=%d", sold, commands)
	}
	window.Advance(1)
	sold, commands = window.totals()
	if sold["moss"] != 0 || sold["plate"] != 1 || commands != 1 {
		t.Fatalf("oldest bucket did not expire: sold=%v commands=%d", sold, commands)
	}
	window.Advance(Window)
	sold, commands = window.totals()
	if len(sold) != 0 || commands != 0 {
		t.Fatalf("idle gap retained activity: sold=%v commands=%d", sold, commands)
	}
	for range 10 * Window {
		window.Advance(1)
		window.RecordSale("moss", 1)
		window.RecordCommand()
	}
	sold, commands = window.totals()
	if len(window.Buckets) > 30 || sold["moss"] != float64(commands) || commands < Window-99 || commands > Window {
		t.Fatalf("rolling storage not bounded: buckets=%d sold=%v commands=%d", len(window.Buckets), sold, commands)
	}
}

func TestRankAllowsGatesTheNextCostSlice(t *testing.T) {
	for _, tc := range []struct {
		rank, level, need int
		allowed           bool
	}{
		{0, 0, 0, true}, {0, 4, 0, true}, {0, 5, 2, false},
		{1, 9, 2, false}, {2, 5, 2, true}, {2, 9, 2, true},
		{3, 10, 4, false}, {4, 10, 4, true}, {7, 20, 8, false}, {8, 20, 8, true},
	} {
		s := New(1)
		s.Rank, s.Levels["yield.moss"] = tc.rank, tc.level
		ok, need := s.RankAllows("yield.moss")
		if ok != tc.allowed || need != tc.need || s.Levels["yield.moss"] != tc.level {
			t.Errorf("rank %d level %d: allowed=%v need=%d", tc.rank, tc.level, ok, need)
		}
	}
}

func TestRankBlockedPurchasePreservesPossessionsAndExplainsBothLanguages(t *testing.T) {
	previous := i18n.Lang()
	t.Cleanup(func() { i18n.SetLang(previous) })
	for _, lang := range []string{"ar", "en"} {
		i18n.SetLang(lang)
		s := New(1)
		s.Levels["yield.moss"], s.Credits, s.Inv["moss"] = 5, 1e6, 100
		u := Find("yield.moss")
		cost, ok, why := s.Status(u)
		want := i18n.T("محتاج رتبة مستعمرة 2", "Needs colony rank 2")
		if ok || why != want || !reflect.DeepEqual(cost, u.CostAt(5)) {
			t.Fatalf("%s: status cost=%+v ok=%v why=%q", lang, cost, ok, why)
		}
		if err := s.Buy(u.ID); err == nil || err.Error() != want || s.Credits != 1e6 || s.Inv["moss"] != 100 || s.Level(u.ID) != 5 {
			t.Fatalf("%s: blocked purchase changed possessions: %+v, %v", lang, s, err)
		}
		s.Rank = 2
		if err := s.Buy(u.ID); err != nil || s.Level(u.ID) != 6 || math.Abs(s.Credits-(1e6-cost.Credits)) > 1e-9 {
			t.Fatalf("%s: eligible purchase failed: %v", lang, err)
		}
	}
}

func TestRankSaveRoundTripAndLegacyOwnership(t *testing.T) {
	s := New(11)
	s.Rank, s.Tick = 3, 9000
	s.RankWindow.Advance(2400)
	s.RankWindow.RecordSale("plate", 7)
	s.RankWindow.RecordCommand()
	blob, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var loaded State
	if err := json.Unmarshal(blob, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.Normalize()
	if !reflect.DeepEqual(s, &loaded) || loaded.RankProgress().Filled != 0.8 {
		t.Fatal("rank or partial window lost on load")
	}
	for _, version := range []int{2, 3} {
		legacy := `{"version":2,"seed":11,"size":3,"tick":9000,"levels":{"engine":8,"yield.moss":12},"credits":123,"inv":{"moss":45},"blueprints":{"api.try":true},"claimed":{"chapter-1":true}}`
		loaded = State{}
		if err := json.Unmarshal([]byte(legacy), &loaded); err != nil {
			t.Fatal(err)
		}
		loaded.Version = version
		loaded.Normalize()
		if loaded.Rank != 0 || loaded.RankProgress().Filled != 0 || loaded.Version != 3 || loaded.Level("engine") != 8 || loaded.Level("yield.moss") != 12 || loaded.Credits != 123 || loaded.Inv["moss"] != 45 || !loaded.Blueprints["api.try"] || !loaded.Claimed["chapter-1"] {
			t.Fatalf("v%d save lost possessions or inherited running time: %+v", version, loaded)
		}
	}
}
