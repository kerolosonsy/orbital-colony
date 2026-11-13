package sim

import (
	"bytes"
	"encoding/json"
	"math"
	"slices"
	"testing"

	"orbital/colony/econ"
)

func TestRankOutputWaitsForFullWindowAndPersists(t *testing.T) {
	w, trace := e2World()
	w.Econ.Inv["moss"] = 300
	if w.Export(0, Moss, 0) <= 0 {
		t.Fatal("sale failed")
	}
	w.Tick = econ.Window - 1
	w.econSnapshot("econ", "", true, true)
	if w.Rank() != 0 || w.Econ.RankProgress().Filled >= 1 {
		t.Fatal("rank earned before a full running window")
	}
	w = reloadColony(t, w)
	if w.Rank() != 0 || w.Econ.RankProgress().Filled != float64(econ.Window-1)/econ.Window {
		t.Fatal("reload changed running time")
	}
	w.Wait(0)
	if w.Rank() != 1 {
		t.Fatal("base sales did not meet output goal at full window")
	}
	w = reloadColony(t, w)
	w.Tick += econ.Window + 1
	w.econSnapshot("econ", "", true, true)
	if w.Rank() != 1 || w.Econ.RankProgress().Value != 0 {
		t.Fatal("earned rank fell or expired sales survived")
	}
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev struct{ Op string }
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "rank" {
			t.Fatal("premature promotion emitted a rank event")
		}
	}
}

func TestRankEfficiencyNeedsActivityAndRewardsFewerCommands(t *testing.T) {
	for _, commands := range []int{1, econ.RankMinCommands - 1, econ.RankMinCommands, 201} {
		w, _ := e2World()
		w.Econ.Rank, w.Econ.Inv["moss"] = 1, 200
		for range commands - 1 {
			w.Wait(0)
		}
		w.Tick = econ.Window - 1
		w.Export(0, Moss, 0)
		want := 1
		if commands == econ.RankMinCommands {
			want = 2
		}
		if w.Rank() != want {
			t.Errorf("%d commands with identical sales: rank=%d want=%d", commands, w.Rank(), want)
		}
	}
}

func TestRankSalesIgnoreYieldAndMarketBonusesAndCountAllCommands(t *testing.T) {
	w, _ := e2World()
	w.Econ.Levels["yield.moss"], w.Econ.Levels["market.contracts"] = 4, 3
	w.Econ.Inv["moss"], w.Econ.Demand["moss"] = 24, 0.5
	w.Econ.Levels["api.spawn"], w.Econ.Levels["hangar"] = 1, 1
	w.applyEcon()
	id := w.Spawn()
	if id != 1 {
		t.Fatal("second drone unavailable")
	}
	w.Export(id, Moss, 12)
	w.Export(0, Moss, 0)
	w.Harvest(id) // An unsuccessful command still counts.
	w.CanHarvest(0)
	beforeTick, beforeCommands := w.Tick, w.Econ.RankWindow.Buckets[0].Commands
	kind, target, current := w.RankGoal()
	if kind != "output" || target != 60 || math.Abs(current-10*600/float64(w.Tick)) > 1e-9 {
		t.Fatalf("sales were not normalized: %s %g %g", kind, target, current)
	}
	if beforeCommands != 4 || w.Tick != beforeTick || w.Econ.RankWindow.Buckets[0].Commands != beforeCommands {
		t.Fatal("commands missed a drone, failure or sensor, or queries spent commands")
	}
}

func TestRankWireSnapshotsAndPromotionCarryNextGoalAndTranslations(t *testing.T) {
	w, trace := e2World()
	w.Tick = econ.Window * 4 / 5
	w.econSnapshot("econ", "", true, true)
	var snapshot map[string]json.RawMessage
	if err := json.Unmarshal(bytes.TrimSpace(trace.Bytes()), &snapshot); err != nil {
		t.Fatal(err)
	}
	var progress econ.RankProgress
	if err := json.Unmarshal(snapshot["rank"], &progress); err != nil {
		t.Fatal(err)
	}
	if progress != (econ.RankProgress{Kind: "output", Target: 60, Filled: 0.8}) {
		t.Fatalf("snapshot rank: %+v", progress)
	}
	var state map[string]json.RawMessage
	if err := json.Unmarshal(snapshot["state"], &state); err != nil || string(state["rank"]) != "0" {
		t.Fatalf("state.rank missing: %s, %v", snapshot["state"], err)
	}
	trace.Reset()
	w.Tick, w.Econ.Inv["moss"] = econ.Window-2, 300
	w.Export(0, Moss, 0)
	w.Move(0, South)
	var promotion econEvent
	for _, line := range bytes.Split(bytes.TrimSpace(trace.Bytes()), []byte{'\n'}) {
		var ev econEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		if ev.Op == "rank" {
			promotion = ev
		}
		if ev.Op == "move" && !slices.Contains(ev.Translations, [2]string{"حيطة! مينفعش تطلع برا الخريطة ناحية ", "Wall! You cannot leave the map to the "}) {
			t.Fatal("promotion consumed the action's translations")
		}
	}
	if !promotion.OK || promotion.Rank.Level != 1 || promotion.Rank.Kind != "efficiency" || promotion.Rank.Target != 1 || promotion.Rank.Value != 150 || promotion.Rank.Filled != 1 || promotion.State == nil || promotion.State.Rank != 1 || promotion.Msg == "" || !slices.Contains(promotion.Translations, [2]string{"⬆ المستعمرة وصلت للرتبة 1", "⬆ Colony reached rank 1"}) {
		t.Fatalf("promotion wire: %+v", promotion)
	}
	w.Tick += econ.Window + 100
	w.econSnapshot("econ", "", true, true)
	if w.Rank() != 1 {
		t.Fatal("idling earned efficiency rank from retained sales")
	}
}
