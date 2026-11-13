package main

import (
	"fmt"
	"testing"

	"orbital/colony/econ"
)

func TestColonyTreeRankCheckAgreesWithShopAndBuy(t *testing.T) {
	for _, lang := range []string{"ar", "en"} {
		for _, rank := range []int{0, 1, 2} {
			st := econ.New(1)
			st.Levels["engine"], st.Credits, st.Rank = 5, 1e6, rank
			node := colonyTreeNode(t, treeForLang(st, lang), "engine")
			var rankCheck *treeCheck
			for i := range node.Checks {
				if node.Checks[i].Kind == "rank" {
					rankCheck = &node.Checks[i]
				}
			}
			wantLabel := tr(lang, "محتاج رتبة مستعمرة 2", "Needs colony rank 2")
			if rankCheck == nil || rankCheck.Ref != "rank" || rankCheck.Label != wantLabel || rankCheck.OK != (rank >= 2) || treeChecksPass(node) != (rank >= 2) || node.State != "owned" {
				t.Fatalf("%s rank %d: tree check %+v, node %+v", lang, rank, rankCheck, node)
			}
			for _, upgrade := range shopForLang(st, lang) {
				if upgrade.ID == "engine" && (upgrade.Can != (rank >= 2) || rank < 2 && upgrade.Reason != wantLabel) {
					t.Fatalf("shop disagrees: %+v", upgrade)
				}
			}
			if err := st.Buy("engine"); (err == nil) != (rank >= 2) {
				t.Fatalf("Buy disagrees with rank check: %v", err)
			}
		}
		for _, need := range []int{2, 4, 8} {
			for _, message := range []string{fmt.Sprintf("محتاج رتبة مستعمرة %d", need), fmt.Sprintf("Needs colony rank %d", need)} {
				if got := economyText(lang, message); got != fmt.Sprintf(tr(lang, "محتاج رتبة مستعمرة %d", "Needs colony rank %d"), need) {
					t.Fatalf("%s: untranslated rank reason %q", lang, got)
				}
			}
		}
	}
}
