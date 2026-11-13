package econ

import (
	"encoding/json"
	"reflect"
	"testing"

	"orbital/colony/i18n"
)

func TestV3MigrationOnlyRerollsUntouchedEmptyStoneOnce(t *testing.T) {
	s := New(42)
	s.Levels["field"] = 6
	s.Normalize()
	s.Version = 2
	for _, layer := range s.Under {
		for i := range layer {
			if layer[i].E != "iron" && layer[i].E != "geode" {
				layer[i].E = ""
			}
		}
	}
	for z := 1; z < s.Depth; z++ {
		*s.AtZ(0, 0, z) = Tile{G: "tunnel"}
		*s.AtZ(1, 0, z) = Tile{G: "shaft"}
		*s.AtZ(2, 0, z) = Tile{G: "stone", E: "iron"}
		*s.AtZ(3, 0, z) = Tile{G: "stone", E: "geode"}
		*s.AtZ(4, 0, z) = Tile{G: "stone", P: 1}
	}
	s.Inv = map[string]float64{"iron": 12.5, "geode": 4, "moss": 2, "water": 3}
	s.Credits, s.Earned = 123, 456
	s.Blueprints["api.dig"] = true
	s.Claimed["chapter-1"] = true
	s.Contracts["c17-deep-haul"] = &ContractRun{Status: "done", Best: 22}
	blob, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var before State
	if err := json.Unmarshal(blob, &before); err != nil {
		t.Fatal(err)
	}
	s.Normalize()
	seen := map[string]int{}
	for z, layer := range s.Under {
		for i, tile := range layer {
			old := before.Under[z][i]
			if old != (Tile{G: "stone"}) && tile != old {
				t.Fatalf("owned tile changed: %v to %v", old, tile)
			}
			if old == (Tile{G: "stone"}) && tile != undergroundAt(s.Seed, i%s.Size, i/s.Size, z+1) {
				t.Fatal("migration did not use coordinate hash")
			}
			seen[tile.E]++
		}
	}
	for _, ore := range []string{"copper", "sulfur", "titanium", "uranium"} {
		if seen[ore] == 0 {
			t.Errorf("migration never added %s", ore)
		}
	}
	before.Version, before.Under = 3, s.Under
	if !reflect.DeepEqual(s, &before) {
		t.Fatal("migration changed possessions or other saved state")
	}
	// A tile emptied after migration must stay empty, including across another load.
	*s.AtZ(5, 5, 5) = Tile{G: "stone"}
	blob, err = json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var loaded State
	if err := json.Unmarshal(blob, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.Normalize()
	if !reflect.DeepEqual(s, &loaded) {
		t.Fatal("version 3 save rerolled a tile on reload")
	}
}

func TestNewOreDepthAndExpansionOrder(t *testing.T) {
	counts := map[string]int{}
	for z := 1; z <= 5; z++ {
		for x := 0; x < 40; x++ {
			for y := 0; y < 40; y++ {
				g := undergroundAt(13, x, y, z).E
				counts[g]++
				if g == "copper" && z < 2 || g == "sulfur" && z < 3 || g == "titanium" && z < 4 || g == "uranium" && z != 5 {
					t.Fatalf("%s at depth %d", g, z)
				}
			}
		}
	}
	for _, ore := range []string{"iron", "geode", "copper", "sulfur", "titanium", "uranium"} {
		if counts[ore] == 0 {
			t.Errorf("missing ore %s", ore)
		}
	}
	staged, direct := New(13), New(13)
	for _, l := range []int{1, 3, 6} {
		staged.Levels["field"] = l
		staged.Normalize()
	}
	direct.Levels["field"] = 6
	direct.Normalize()
	if !reflect.DeepEqual(staged.Under, direct.Under) {
		t.Fatal("new ore depends on expansion order")
	}
}

func TestBuildingSaveResumesSpentAlternativeRecipe(t *testing.T) {
	previous := Recipes
	Recipes = append([]Recipe(nil), Recipes...)
	t.Cleanup(func() { Recipes = previous })
	for i := range Recipes {
		if Recipes[i].Alt {
			Recipes[i].Ticks = 80
		}
	}
	s := New(1)
	b := &Building{Kind: "composter", X: 1, Y: 2, Recipe: "biomass", In: map[string]float64{"lichen": 2}, Out: map[string]float64{}, Enabled: true, Wear: 10}
	s.Buildings = []*Building{b}
	if err := b.Step(20, true); err != nil {
		t.Fatal(err)
	}
	if _, err := b.Put("moss", 4); err != nil {
		t.Fatal(err)
	}
	blob, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var loaded State
	if err := json.Unmarshal(blob, &loaded); err != nil {
		t.Fatal(err)
	}
	loaded.Normalize()
	if !reflect.DeepEqual(s, &loaded) {
		t.Fatal("building or active recipe lost during round trip")
	}
	b = loaded.Buildings[0]
	if err := b.Step(59, true); err != nil || b.Out["biomass"] != 0 || b.Progress != 79 {
		t.Fatalf("alternative finished early after load: %+v, %v", b, err)
	}
	if err := b.Step(1, true); err != nil || b.Out["biomass"] != 1 || b.In["moss"] != 4 || b.In["lichen"] != 0 || b.Wear != 11 {
		t.Fatalf("resumed cycle: %+v, %v", b, err)
	}
}

func TestBuildingCatalogPurchasesMatchTree(t *testing.T) {
	for _, node := range buildingNodes() {
		t.Run(node.ID, func(t *testing.T) {
			s := New(1)
			u := Find(node.ID)
			if u == nil || u.GroupEn != "Buildings" || u.Max <= 1 {
				t.Fatalf("building upgrade missing: %+v", u)
			}
			cost := u.CostAt(0)
			if cost.Credits <= 0 || len(cost.Goods) == 0 {
				t.Fatal("building needs both credits and goods")
			}
			s.Credits = cost.Credits
			for good, n := range cost.Goods {
				s.Inv[good] = n
			}
			for _, p := range node.After {
				s.Levels[p.ID] = p.Level
			}
			for _, p := range node.After {
				s.Levels[p.ID]--
				if err := s.Buy(node.ID); err == nil {
					t.Fatalf("accepted %s below level %d", p.ID, p.Level)
				}
				s.Levels[p.ID]++
			}
			if err := s.Buy(node.ID); err != nil {
				t.Fatal(err)
			}
			if s.Credits != 0 || s.Level(node.ID) != 1 {
				t.Fatal("purchase did not spend exact credits or grant level 1")
			}
			for good := range cost.Goods {
				if s.Inv[good] != 0 {
					t.Fatal("purchase did not spend goods")
				}
			}
			s.Levels[node.ID], s.Credits = u.Max, 1e12
			for good := range cost.Goods {
				s.Inv[good] = 1e12
			}
			if err := s.Buy(node.ID); err == nil {
				t.Fatal("building ignored its level cap")
			}
		})
	}
}

func TestAllGoodsHaveLocalizedNamesAndV3MarketPrices(t *testing.T) {
	previous := i18n.Lang()
	t.Cleanup(func() { i18n.SetLang(previous) })
	s := New(1)
	for _, lang := range []string{"ar", "en"} {
		i18n.SetLang(lang)
		for _, resources := range [][]Resource{RawV3, Crafted} {
			for _, g := range resources {
				want := g.Name
				if lang == "en" {
					want = g.NameEn
				}
				if CropName(g.ID) != want || s.Prices()[g.ID] != PriceV3[g.ID] {
					t.Errorf("missing localized market good %s", g.ID)
				}
			}
		}
	}
	for good, want := range map[string]float64{"moss": 1, "lichen": 3, "glowcap": 2, "reed": 14, "crystal": 6, "iron": 9, "geode": 40} {
		if s.Price(good) != want {
			t.Errorf("old price changed: %s", good)
		}
	}
	if len(s.Prices()) != 28 {
		t.Fatalf("market has %d goods, want 28", len(s.Prices()))
	}
	s.Levels["market.contracts"] = 2
	for _, good := range Goods {
		s.Inv[good], s.Demand[good] = 2, 0.5
		want := PriceV3[good] * 1.2 * 0.5
		if got := s.Sell(good, 1); got != want || s.Inv[good] != 1 {
			t.Errorf("Sell(%s) = %g, want %g", good, got, want)
		}
	}
}
