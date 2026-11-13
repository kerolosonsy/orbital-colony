package main

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"orbital/colony/econ"
)

func colonyTreeNode(t *testing.T, nodes []treeNodeView, id string) treeNodeView {
	t.Helper()
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
	}
	t.Fatalf("missing tree node %q", id)
	return treeNodeView{}
}

func treeChecksPass(node treeNodeView) bool {
	for _, check := range node.Checks {
		if !check.OK {
			return false
		}
	}
	return true
}

func TestColonyTreeFreshVisibilityAndAffordability(t *testing.T) {
	for _, credits := range []float64{0, 40} {
		st := econ.New(1)
		st.Credits = credits
		nodes := treeForLang(st, "en")
		if len(nodes) != len(econ.Tree) {
			t.Fatalf("got %d nodes, want %d", len(nodes), len(econ.Tree))
		}
		base := colonyTreeNode(t, nodes, "base")
		if base.State != "owned" || base.Level != 1 || treeChecksPass(base) {
			t.Fatalf("base must be owned and not buyable: %+v", base)
		}
		for _, id := range []string{"field", "engine", "yield.moss"} {
			node := colonyTreeNode(t, nodes, id)
			want := "near"
			if credits == 40 {
				want = "available"
			}
			if node.State != want || treeChecksPass(node) != (want == "available") {
				t.Fatalf("credits %.0f: %s = %+v, want %s", credits, id, node, want)
			}
		}
		for id, want := range map[string]string{"crop.lichen": "locked", "crop.glowcap": "unknown", "composter": "locked", "assembler": "unknown"} {
			if node := colonyTreeNode(t, nodes, id); node.State != want {
				t.Errorf("%s = %s, want %s", id, node.State, want)
			}
		}
	}
	st := econ.New(1)
	st.Levels["yield.moss"] = 1
	nodes := treeForLang(st, "en")
	if colonyTreeNode(t, nodes, "composter").State != "near" || colonyTreeNode(t, nodes, "refinery").State != "locked" {
		t.Fatal("owning a node must reveal descendants up to two steps away")
	}
}

func TestColonyTreeBlueprintChecksAndLiveGoodsCosts(t *testing.T) {
	st := econ.New(1)
	st.Levels["crop.lichen"] = 1
	st.Credits, st.RP = 450, econ.FindNode("crop.glowcap").RP
	st.Inv["lichen"] = 39
	node := colonyTreeNode(t, treeForLang(st, "en"), "crop.glowcap")
	wantCost, _, _ := st.Status(econ.Find(node.ID))
	if !reflect.DeepEqual(node.Cost, wantCost) {
		t.Fatalf("tree cost %+v differs from today's cost %+v", node.Cost, wantCost)
	}
	var kinds []string
	for _, check := range node.Checks {
		kinds = append(kinds, check.Kind)
	}
	if !reflect.DeepEqual(kinds, []string{"parent", "contract", "rp", "cost", "cost"}) {
		t.Fatalf("check order: %v", kinds)
	}
	contract := econ.BlueprintSource(node.ID)
	if check := node.Checks[1]; check.OK || check.Ref != contract.ID || check.Label != contract.TitleEn {
		t.Fatalf("blueprint check: %+v", check)
	}
	st.Blueprints[node.ID] = true
	node = colonyTreeNode(t, treeForLang(st, "en"), node.ID)
	if node.State != "near" || node.Checks[len(node.Checks)-1].OK {
		t.Fatalf("goods shortage must block purchase: %+v", node)
	}
	st.Inv["lichen"] = 40
	node = colonyTreeNode(t, treeForLang(st, "en"), node.ID)
	if node.State != "available" || !treeChecksPass(node) {
		t.Fatalf("exact live cost must be buyable: %+v", node)
	}
}

func TestColonyTreeBuildingsUseLiveCostsAndPrerequisiteLevels(t *testing.T) {
	st := econ.New(1)
	st.Levels["generator"] = 1
	st.RP = econ.FindNode("assembler").RP
	cost := econ.Find("assembler").CostAt(0)
	st.Credits = cost.Credits
	for good, n := range cost.Goods {
		st.Inv[good] = n
	}
	for _, level := range []int{2, 3} {
		st.Levels["smelter"] = level
		node := colonyTreeNode(t, treeForLang(st, "en"), "assembler")
		if node.After[0] != (treePrereq{"smelter", 3}) || node.Checks[0].OK != (level >= 3) {
			t.Fatalf("smelter level %d: %+v", level, node)
		}
		if !reflect.DeepEqual(node.Cost, cost) || treeChecksPass(node) != (level >= 3) {
			t.Fatalf("building availability: %+v", node)
		}
		if err := st.Buy("assembler"); (err == nil) != (level >= 3) {
			t.Fatalf("tree and shop disagree at smelter level %d: %v", level, err)
		}
	}
	for _, node := range treeForLang(st, "en") {
		for _, check := range node.Checks {
			if check.Kind == "soon" {
				t.Fatalf("obsolete soon check: %+v", node)
			}
		}
	}
}

func TestColonyTreeOwnedUpgradeCostsAndMax(t *testing.T) {
	st := econ.New(1)
	upgrade := econ.Find("engine")
	st.Levels[upgrade.ID] = 1
	cost, _, _ := st.Status(upgrade)
	for _, credits := range []float64{cost.Credits - 1, cost.Credits} {
		st.Credits = credits
		node := colonyTreeNode(t, treeForLang(st, "en"), upgrade.ID)
		if node.State != "owned" || len(node.Checks) != 1 || node.Checks[0].Kind != "cost" || treeChecksPass(node) != (credits >= cost.Credits) {
			t.Fatalf("repeat upgrade: %+v", node)
		}
	}
	st.Levels[upgrade.ID] = upgrade.Max
	st.Credits = 1e12
	node := colonyTreeNode(t, treeForLang(st, "en"), upgrade.ID)
	last := node.Checks[len(node.Checks)-1]
	if node.State != "owned" || treeChecksPass(node) || last.Kind != "max" || last.OK {
		t.Fatalf("maxed upgrade must be disabled: %+v", node)
	}
}

func TestColonyTreeLocalization(t *testing.T) {
	st := econ.New(1)
	for _, lang := range []string{"ar", "en"} {
		nodes := treeForLang(st, lang)
		for _, source := range econ.Tree {
			node := colonyTreeNode(t, nodes, source.ID)
			if node.Name != tr(lang, source.Name, source.NameEn) || node.Desc != tr(lang, source.Desc, source.DescEn) || node.Effect != node.Desc {
				t.Errorf("%s: %s text not localized: %+v", lang, source.ID, node)
			}
			if source.Branch != "" && node.Branch == "" {
				t.Errorf("missing branch translation for %s", source.Branch)
			}
		}
		field := colonyTreeNode(t, nodes, "field")
		if field.Branch != tr(lang, "الزراعة", "Farming") || field.Checks[0].Label != tr(lang, "القاعدة: مستوى 1", "Base: level 1") {
			t.Errorf("%s: branch/parent localization: %+v", lang, field)
		}
		composter := colonyTreeNode(t, nodes, "composter")
		if composter.Checks[2].Label != tr(lang, "100 كريدت", "100 credits") || composter.Checks[3].Label != tr(lang, "5 الطحالب", "5 Moss") {
			t.Errorf("%s: building cost localization: %+v", lang, composter)
		}
		for _, good := range []struct{ id, ar, en string }{{"kiln", "السيليكا", "Silica"}, {"refinery", "الكتلة الحيوية", "Biomass"}, {"lab", "الدائرة الإلكترونية", "Circuit"}} {
			checks := treeCostChecks(st, econ.Find(good.id).CostAt(0), lang)
			if got := checks[len(checks)-1].Label; got != "5 "+tr(lang, good.ar, good.en) {
				t.Errorf("%s: %s cost = %s", lang, good.id, got)
			}
		}
		scanner := colonyTreeNode(t, nodes, "api.scanner")
		contract := econ.BlueprintSource(scanner.ID)
		if scanner.Checks[1].Label != tr(lang, contract.Title, contract.TitleEn) {
			t.Errorf("%s: contract localization: %+v", lang, scanner)
		}
	}
}

func TestColonyTreeSummaryAndPurchaseResponse(t *testing.T) {
	s := testColony(t)
	s.colony.state.Credits = 100
	for _, lang := range []string{"ar", "en"} {
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/api/colony", nil)
		request.Header.Set("X-Orbital-Lang", lang)
		s.colonyStatus(response, request)
		var summary struct {
			Tree []treeNodeView `json:"tree"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &summary); err != nil {
			t.Fatal(err)
		}
		if node := colonyTreeNode(t, summary.Tree, "engine"); node.State != "available" || node.Name != tr(lang, "محرك أسرع", "Faster Engine") {
			t.Fatalf("summary tree: %+v", node)
		}
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/colony/buy", strings.NewReader(`{"id":"engine"}`))
	request.Header.Set("X-Orbital-Lang", "en")
	s.colonyBuy(response, request)
	var purchase struct {
		OK   bool           `json:"ok"`
		Tree []treeNodeView `json:"tree"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &purchase); err != nil {
		t.Fatal(err)
	}
	if node := colonyTreeNode(t, purchase.Tree, "engine"); !purchase.OK || node.Level != 1 || node.State != "owned" || node.Cost.Credits != 41 {
		t.Fatalf("purchase tree not refreshed: %s", response.Body)
	}
}
