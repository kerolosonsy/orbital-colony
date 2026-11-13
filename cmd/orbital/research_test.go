package main

import (
	"encoding/json"
	"testing"

	"orbital/colony/econ"
)

func TestResearchTreeAndShopAgreeOnRPAndLevelParents(t *testing.T) {
	st := econ.New(1)
	st.Credits = 1000
	id := "research.engine"
	for _, tc := range []struct {
		level int
		rp    float64
		want  bool
	}{{2, 100, false}, {3, 19.5, false}, {3, 20, true}} {
		st.Levels["engine"], st.RP = tc.level, tc.rp
		node := colonyTreeNode(t, treeForLang(st, "en"), id)
		_, ok, _ := st.Status(econ.Find(id))
		if ok != tc.want || treeChecksPass(node) != ok {
			t.Fatalf("shop/tree mismatch: %+v", node)
		}
	}
	if err := st.Buy(id); err != nil {
		t.Fatal(err)
	}
	st.Levels["engine"] = 0
	node := colonyTreeNode(t, treeForLang(st, "en"), id)
	if !treeChecksPass(node) || len(node.Checks) != 1 || st.RP != 0 {
		t.Fatalf("repeat includes RP or parents: %+v", node)
	}
}

func TestResearchSummaryWireAndLocalization(t *testing.T) {
	st := econ.New(1)
	st.RP = 12.5
	st.Levels["engine"] = 3
	st.Contracts["c03-full-sweep"] = &econ.ContractRun{Best: 1}
	st.Contracts["e03-sweep"] = &econ.ContractRun{Best: 1}
	for _, lang := range []string{"ar", "en"} {
		nodes := treeForLang(st, lang)
		node := colonyTreeNode(t, nodes, "research.engine")
		if node.RP != 20 || node.Checks[1].Kind != "rp" {
			t.Fatalf("missing RP wire check: %+v", node)
		}
		blob, err := json.Marshal(node)
		if err != nil {
			t.Fatal(err)
		}
		var wire map[string]any
		if err := json.Unmarshal(blob, &wire); err != nil {
			t.Fatal(err)
		}
		if wire["rp"] != float64(20) {
			t.Fatalf("missing node rp: %s", blob)
		}
		found := false
		for _, c := range contractsForLang(st, lang, map[int]bool{3: true}) {
			if c.ID == "e03-sweep" {
				found = true
				if !c.Repeatable || c.RewardRP != 8 || c.Status != "available" {
					t.Fatalf("repeatable contract wire: %+v", c)
				}
			}
		}
		if !found {
			t.Fatal("experiment missing from contract view")
		}
		_, _, why := st.Status(econ.Find("research.engine"))
		if got := economyText(lang, why); got != tr(lang, "ناقصك 7.5 RP", "You need 7.5 more RP") {
			t.Fatalf("RP reason not localized: %s", got)
		}
	}
	blob, err := json.Marshal(st)
	if err != nil {
		t.Fatal(err)
	}
	var wire map[string]any
	if err := json.Unmarshal(blob, &wire); err != nil {
		t.Fatal(err)
	}
	if wire["rp"] != 12.5 {
		t.Fatal("state RP missing")
	}
}
