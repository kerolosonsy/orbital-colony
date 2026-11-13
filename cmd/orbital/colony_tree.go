package main

import (
	"fmt"
	"sort"

	"orbital/colony/econ"
)

type treePrereq struct {
	ID    string `json:"id"`
	Level int    `json:"level"`
}

type treeCheck struct {
	Kind  string `json:"kind"`
	Ref   string `json:"ref"`
	Label string `json:"label"`
	OK    bool   `json:"ok"`
}

type treeNodeView struct {
	RP       float64      `json:"rp,omitempty"`
	ID       string       `json:"id"`
	Branch   string       `json:"branch"`
	BranchEn string       `json:"branchEn"`
	Name     string       `json:"name"`
	NameEn   string       `json:"nameEn"`
	Desc     string       `json:"desc"`
	DescEn   string       `json:"descEn"`
	Effect   string       `json:"effect"`
	EffectEn string       `json:"effectEn"`
	After    []treePrereq `json:"after"`
	Level    int          `json:"level"`
	Cost     econ.Cost    `json:"cost"`
	Checks   []treeCheck  `json:"checks"`
	State    string       `json:"state"`
}

func treeForLang(st *econ.State, lang string) []treeNodeView {
	shop := make(map[string]econ.ShopItem, len(econ.Catalog))
	for _, upgrade := range st.Shop() {
		shop[upgrade.ID] = upgrade
	}
	nodes := make([]treeNodeView, 0, len(econ.Tree))
	for _, node := range econ.Tree {
		view := treeNodeForLang(node, lang)
		view.Level = st.Levels[node.ID]
		view.Checks = treeChecks(st, node, shop[node.ID], lang)
		view.Cost = shop[node.ID].Cost
		if node.ID == "base" {
			view.Level = 1
		}
		view.State = treeState(st, node, view.Checks)
		nodes = append(nodes, view)
	}
	return nodes
}

func treeNodeForLang(node econ.Node, lang string) treeNodeView {
	branchEn := treeBranchEnglish(node.Branch)
	view := treeNodeView{
		RP: node.RP, ID: node.ID, Branch: tr(lang, node.Branch, branchEn), BranchEn: branchEn,
		Name: tr(lang, node.Name, node.NameEn), NameEn: node.NameEn,
		Desc: tr(lang, node.Desc, node.DescEn), DescEn: node.DescEn,
		Effect: tr(lang, node.Desc, node.DescEn), EffectEn: node.DescEn,
		After: make([]treePrereq, 0, len(node.After)),
	}
	for _, parent := range node.After {
		view.After = append(view.After, treePrereq{parent.ID, parent.Level})
	}
	return view
}

func treeChecks(st *econ.State, node econ.Node, upgrade econ.ShopItem, lang string) []treeCheck {
	if node.ID == "base" {
		return []treeCheck{{Kind: "max", Label: tr(lang, "مكتمل", "Maxed")}}
	}
	checks := make([]treeCheck, 0)
	if st.Levels[node.ID] == 0 || upgrade.Upgrade == nil {
		checks = treeParentChecks(st, node, lang)
	}
	if upgrade.Upgrade == nil {
		return checks
	}
	if st.Levels[node.ID] == 0 && upgrade.Blueprint != "" {
		contract := econ.FindContract(upgrade.Blueprint)
		checks = append(checks, treeCheck{Kind: "contract", Ref: contract.ID, Label: tr(lang, contract.Title, contract.TitleEn)})
	}
	if allowed, need := st.RankAllows(node.ID); need > 0 && upgrade.Level < upgrade.Max {
		checks = append(checks, treeCheck{Kind: "rank", Ref: "rank", Label: fmt.Sprintf(tr(lang, "محتاج رتبة مستعمرة %d", "Needs colony rank %d"), need), OK: allowed})
	}
	if st.Level(node.ID) == 0 && node.RP > 0 {
		checks = append(checks, treeCheck{Kind: "rp", Ref: "rp", Label: fmt.Sprintf(tr(lang, "%.0f نقط بحث", "%.0f RP"), node.RP), OK: st.RP >= node.RP})
	}
	checks = append(checks, treeCostChecks(st, upgrade.Cost, lang)...)
	if upgrade.Level >= upgrade.Max {
		checks = append(checks, treeCheck{Kind: "max", Label: tr(lang, "مكتمل", "Maxed")})
	}
	return checks
}

func treeParentChecks(st *econ.State, node econ.Node, lang string) []treeCheck {
	checks := make([]treeCheck, 0, len(node.After))
	for _, parent := range node.After {
		name := econ.FindNode(parent.ID)
		label := fmt.Sprintf(tr(lang, "%s: مستوى %d", "%s: level %d"), tr(lang, name.Name, name.NameEn), parent.Level)
		ok := st.Levels[parent.ID] >= parent.Level || parent.ID == "base" && parent.Level <= 1
		checks = append(checks, treeCheck{Kind: "parent", Ref: parent.ID, Label: label, OK: ok})
	}
	return checks
}

func treeCostChecks(st *econ.State, cost econ.Cost, lang string) []treeCheck {
	checks := []treeCheck{{Kind: "cost", Ref: "credits", Label: fmt.Sprintf(tr(lang, "%.0f كريدت", "%.0f credits"), cost.Credits), OK: st.Credits >= cost.Credits}}
	goods := make([]string, 0, len(cost.Goods))
	for good := range cost.Goods {
		goods = append(goods, good)
	}
	sort.Strings(goods)
	for _, good := range goods {
		name := economyName(lang, econ.CropName(good))
		label := fmt.Sprintf("%.0f %s", cost.Goods[good], name)
		checks = append(checks, treeCheck{Kind: "cost", Ref: good, Label: label, OK: st.Inv[good] >= cost.Goods[good]})
	}
	return checks
}

func treeState(st *econ.State, node econ.Node, checks []treeCheck) string {
	if node.ID == "base" || st.Levels[node.ID] > 0 {
		return "owned"
	}
	if !treeNearOwned(node, st.Levels, 2) {
		return "unknown"
	}
	state := "available"
	for _, check := range checks {
		if !check.OK {
			if check.Kind == "parent" {
				return "locked"
			}
			state = "near"
		}
	}
	return state
}

func treeNearOwned(node econ.Node, levels map[string]int, steps int) bool {
	if node.ID == "base" || levels[node.ID] > 0 {
		return true
	}
	if steps == 0 {
		return false
	}
	for _, parent := range node.After {
		if treeNearOwned(*econ.FindNode(parent.ID), levels, steps-1) {
			return true
		}
	}
	return false
}

func treeBranchEnglish(branch string) string {
	return map[string]string{
		"": "", "الزراعة": "Farming", "الأسطول": "Fleet", "البيانات": "Data",
		"التعدين": "Mining", "الخدمات": "Services", "التصنيع": "Manufacturing", "الطاقة": "Energy",
	}[branch]
}
