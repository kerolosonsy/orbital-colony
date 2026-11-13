package econ

import (
	"encoding/json"
	"strings"
	"testing"

	"orbital/colony/i18n"
)

func TestContractTable(t *testing.T) {
	seen := map[string]bool{}
	granted := map[string]bool{}
	for _, c := range Contracts {
		if seen[c.ID] {
			t.Fatalf("duplicate contract %s", c.ID)
		}
		seen[c.ID] = true
		if c.Title == "" || c.TitleEn == "" || c.Brief == "" || c.BriefEn == "" {
			t.Errorf("%s: missing Arabic or English text", c.ID)
		}
		if c.Chapter < 1 || c.Window <= 0 || len(c.Requires) == 0 {
			t.Errorf("%s: chapter, window and requires must be set", c.ID)
		}
		for _, a := range c.After {
			if p := FindContract(a); p == nil || p.Chapter > c.Chapter {
				t.Errorf("%s: bad After %q", c.ID, a)
			}
		}
		for _, g := range c.Reward.Grants {
			if Find(g) == nil {
				t.Errorf("%s: grants unknown upgrade %q", c.ID, g)
			}
			granted[g] = true
		}
		for _, d := range c.Limits.Deny {
			if Find(d) == nil {
				t.Errorf("%s: denies unknown upgrade %q", c.ID, d)
			}
		}
	}
	for _, u := range Catalog {
		if needsBlueprint(u.ID) && !granted[u.ID] {
			t.Errorf("no contract unlocks the blueprint for %s", u.ID)
		}
	}
}

func TestContractAvailableAcceptCancel(t *testing.T) {
	s := New(1)
	if got := s.Available(map[int]bool{}); len(got) != 0 {
		t.Fatalf("contracts open without chapters: %d", len(got))
	}
	done := map[int]bool{2: true, 4: true, 5: true}
	ids := func() []string {
		var out []string
		for _, c := range s.Available(done) {
			out = append(out, c.ID)
		}
		return out
	}
	if got := strings.Join(ids(), ","); got != "c02-price-watch,c04-right-ground" {
		t.Fatalf("available = %s (c05 needs c04 won first)", got)
	}
	if s.Accept("c05-careful-harvest") == nil {
		t.Fatal("accepted a contract before its After")
	}
	s.Tick = 40
	if err := s.Accept("c02-price-watch"); err != nil {
		t.Fatal(err)
	}
	if s.Accept("c04-right-ground") == nil || len(s.Available(done)) != 0 {
		t.Fatal("two contracts active at once")
	}
	if r := s.Contracts["c02-price-watch"]; r.Status != "active" || r.Start != 40 || r.Tries != 1 {
		t.Fatalf("run = %+v", r)
	}
	if err := s.Accept(""); err != nil || s.Active != "" || s.Contracts["c02-price-watch"].Status != "" {
		t.Fatalf("cancel failed: %v %+v", err, s.Contracts["c02-price-watch"])
	}
	if s.Accept("") == nil || s.Accept("missing") == nil {
		t.Fatal("cancel without contract or unknown id accepted")
	}
}

func TestContractProgressSettleAndReplay(t *testing.T) {
	s := New(1)
	s.Credits = 0
	if err := s.Accept("c04-right-ground"); err != nil {
		t.Fatal(err)
	}
	if !s.Usable("api.scanner") || s.Has("api.scanner") {
		t.Fatal("the contract's reward is not lent while it runs")
	}
	s.RecordSale(0, "moss", 30, 30, 1)
	if done, failed, _ := s.Progress(100); done || failed {
		t.Fatal("finished with lichen missing")
	}
	if f := s.Fraction(); f != 0 {
		t.Fatalf("fraction = %v, want 0 (lichen is the slowest goal)", f)
	}
	s.RecordSale(0, "lichen", 30, 90, 3)
	if s.Left(100) != 1400 {
		t.Fatalf("left = %d", s.Left(100))
	}
	if _, err := s.Settle(50); err != nil {
		t.Fatal(err)
	}
	if !s.Blueprints["api.scanner"] || s.Has("api.scanner") || s.Active != "" || s.Contracts["c04-right-ground"].Best != 50 {
		t.Fatalf("settle: %+v", s.Contracts["c04-right-ground"])
	}
	if _, ok, _ := s.Status(Find("api.scanner")); ok {
		t.Fatal("blueprint should not make the purchase free")
	}
	s.Credits = 1000
	if err := s.Buy("api.scanner"); err != nil {
		t.Fatal(err)
	}
	// a replay never pays twice, keeps the best time, and a failure keeps the win
	s.Tick = 1000
	if err := s.Accept("c04-right-ground"); err != nil {
		t.Fatal(err)
	}
	s.RecordLost()
	if _, failed, why := s.Progress(1001); !failed || why == "" {
		t.Fatal("losing a plant did not fail MaxLost 0")
	}
	s.Fail()
	if r := s.Contracts["c04-right-ground"]; r.Status != "done" || r.Best != 50 || r.Tries != 2 {
		t.Fatalf("failed replay lost the win: %+v", r)
	}
	s.Accept("c04-right-ground")
	s.RecordSale(0, "moss", 30, 30, 1)
	s.RecordSale(0, "lichen", 30, 90, 3)
	credits := s.Credits
	if r, err := s.Settle(1030); err != nil || r.Credits != 0 || s.Credits != credits || s.Contracts["c04-right-ground"].Best != 30 {
		t.Fatalf("replay settle: %v %+v", err, s.Contracts["c04-right-ground"])
	}
}

func TestContractRules(t *testing.T) {
	cases := []struct {
		id     string
		play   func(s *State)
		failed bool
	}{
		{"c02-price-watch", func(s *State) { s.RecordSale(0, "moss", 50, 50, 0.95) }, false},
		{"c02-price-watch", func(s *State) { s.RecordSale(0, "moss", 50, 50, 0.85) }, true},
		{"c08-big-veins", func(s *State) { s.RecordHarvest("crystal", 3, 3) }, true},
		{"c11-dispatch", func(s *State) { s.RecordSale(0, "crystal", 1, 6, 6); s.RecordSale(1, "crystal", 1, 6, 6) }, true},
		{"c15-shortest-run", func(s *State) {
			for range 401 {
				s.RecordCommand(0, true)
			}
		}, true},
		{"c03-full-sweep", func(s *State) { s.RecordSale(0, "moss", 10, 10, 1) }, true}, // window
	}
	for _, tc := range cases {
		s := New(1)
		s.Accept(tc.id)
		tc.play(s)
		tick := int64(10)
		if tc.id == "c03-full-sweep" {
			tick = 901
		}
		if _, failed, why := s.Progress(tick); failed != tc.failed {
			t.Errorf("%s: failed = %v (%s)", tc.id, failed, why)
		}
	}
	// goal shapes: kinds, average harvest and working drones
	s := New(1)
	s.Accept("c09-crop-protocol")
	s.RecordSale(0, "moss", 40, 40, 1)
	s.RecordSale(0, "lichen", 40, 120, 3)
	if done, _, _ := s.Progress(1); done {
		t.Fatal("two kinds finished a three-kind goal")
	}
	s.RecordSale(0, "glowcap", 40, 80, 2)
	if done, _, _ := s.Progress(1); !done {
		t.Fatal("three kinds did not finish")
	}
	s = New(1)
	s.Accept("c07-plot-type")
	s.RecordSale(0, "glowcap", 40, 80, 2)
	s.RecordHarvest("glowcap", 2, 0)
	if done, _, _ := s.Progress(1); done {
		t.Fatal("average 2 met a minimum of 3")
	}
	s.RecordHarvest("glowcap", 5, 0)
	if done, _, _ := s.Progress(1); !done {
		t.Fatal("average 3.5 did not meet 3")
	}
	s = New(1)
	s.Accept("c10-fleet")
	if s.DroneCap() != 3 {
		t.Fatalf("fleet contract lends %d drones", s.DroneCap())
	}
	s.RecordSale(0, "moss", 200, 200, 1)
	for d := range 3 {
		for range MinCommands - 1 {
			s.RecordCommand(d, false)
		}
	}
	if done, _, _ := s.Progress(1); done {
		t.Fatal("idle drones counted as working")
	}
	for d := range 3 {
		s.RecordCommand(d, false)
	}
	if done, _, _ := s.Progress(1); !done {
		t.Fatal("three working drones did not finish")
	}
	s = New(1)
	s.Levels["api.scanner"] = 1
	s.Contracts["c04-right-ground"] = &ContractRun{Status: "done", Best: 10}
	if err := s.Accept("c05-careful-harvest"); err != nil {
		t.Fatal(err)
	}
	if s.Usable("api.scanner") || !s.Usable("api.try") {
		t.Fatal("Deny or lending failed")
	}
}

func TestContractMigrationV1(t *testing.T) {
	v1 := `{"version":1,"seed":7,"size":3,"tiles":[],"inv":{"moss":4},"credits":12,"levels":{"yield.moss":12,"api.radar":1,"crop.crystal":1,"field":0},"claimed":{"chapter-8":true}}`
	var s State
	if err := json.Unmarshal([]byte(v1), &s); err != nil {
		t.Fatal(err)
	}
	s.Normalize()
	if s.Version != 2 || s.Levels["yield.moss"] != 12 || s.Levels["api.radar"] != 1 || s.Levels["crop.crystal"] != 1 || s.Credits != 12 || !s.Claimed["chapter-8"] {
		t.Fatalf("migration lost owned state: %+v", s)
	}
	if !s.Blueprints["api.radar"] || !s.Blueprints["crop.crystal"] || s.Blueprints["api.scanner"] || s.Blueprints["yield.moss"] {
		t.Fatalf("blueprints = %v", s.Blueprints)
	}
	if !s.Usable("api.radar") {
		t.Fatal("owned radar stopped working")
	}
	// a v2 save is not migrated again: a level without a blueprint stays locked to buy
	s2 := New(1)
	s2.Levels["api.prices"] = 1
	s2.Normalize()
	if s2.Blueprints["api.prices"] {
		t.Fatal("v2 save migrated twice")
	}
}

func TestContractTextsBothLanguages(t *testing.T) {
	previous := i18n.Lang()
	t.Cleanup(func() { i18n.SetLang(previous) })
	for _, lang := range []string{"ar", "en"} {
		i18n.SetLang(lang)
		s := New(1)
		s.Accept("c03-full-sweep")
		s.RecordSale(0, "moss", 12, 12, 1)
		_, _, why := s.Progress(5)
		_, _, pair := s.ProgressPair(5)
		want := map[string]string{"ar": "الطحالب 12/60", "en": "Moss 12/60"}[lang]
		if why != want || pair != [2]string{"الطحالب 12/60", "Moss 12/60"} {
			t.Errorf("%s: progress = %q, pair = %q", lang, why, pair)
		}
		_, _, why = New(1).Status(Find("crop.lichen"))
		want = map[string]string{"ar": "خلّص العقد الأول: المسح الكامل", "en": "Finish the contract first: Full Sweep"}[lang]
		if why != want {
			t.Errorf("%s: blueprint reason = %q", lang, why)
		}
	}
}
