package main

import (
	"encoding/json"
	"testing"

	"orbital/colony/econ"
)

func TestColonyUnlimitedShopAndTreeStayBuyablePastOldCaps(t *testing.T) {
	st := econ.New(1)
	st.Rank, st.Credits, st.RP = 1000, 1e100, 1000
	for _, u := range econ.Catalog {
		if u.Max == 0 {
			st.Levels[u.ID] = 25
		}
	}
	for _, good := range econ.Goods {
		st.Inv[good] = 1e100
	}
	for _, lang := range []string{"ar", "en"} {
		for _, shop := range shopForLang(st, lang) {
			if shop.Max != 0 {
				continue
			}
			node := colonyTreeNode(t, treeForLang(st, lang), shop.ID)
			if !shop.Can || shop.Reason != "" || node.Max != 0 || !treeChecksPass(node) || node.State != "owned" {
				t.Fatalf("%s unlimited %s is blocked: shop=%+v tree=%+v", lang, shop.ID, shop, node)
			}
			for _, check := range node.Checks {
				if check.Kind == "max" {
					t.Fatalf("unlimited %s has maxed check", node.ID)
				}
			}
			for _, view := range []any{shop, node} {
				blob, err := json.Marshal(view)
				if err != nil {
					t.Fatal(err)
				}
				var wire map[string]json.RawMessage
				if err := json.Unmarshal(blob, &wire); err != nil {
					t.Fatal(err)
				}
				if string(wire["max"]) != "0" {
					t.Fatalf("unlimited max must be explicit: %s", blob)
				}
			}
		}
	}
}
