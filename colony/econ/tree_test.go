package econ

import (
	"reflect"
	"testing"
)

func TestSeedTreeIncludesCatalogAndBuildingPrerequisites(t *testing.T) {
	if err := Validate(Tree); err != nil {
		t.Fatal(err)
	}
	branches := map[string]bool{"الزراعة": true, "التعدين": true, "التصنيع": true, "الأسطول": true, "الطاقة": true, "البيانات": true, "الخدمات": true, "الشبكات": true, "الأنظمة": true}
	for _, n := range Tree {
		if n.Name == "" || n.NameEn == "" || n.Desc == "" || n.DescEn == "" || n.ID != "base" && !branches[n.Branch] {
			t.Errorf("incomplete node: %+v", n)
		}
	}
	for _, u := range Catalog {
		n := FindNode(u.ID)
		if n == nil {
			t.Fatalf("missing catalog node %s", u.ID)
		}
		for _, id := range u.Requires {
			found := false
			for _, p := range n.After {
				found = found || p.ID == id && p.Level == 1
			}
			if !found {
				t.Errorf("%s lost prerequisite %s", u.ID, id)
			}
		}
	}
	for _, id := range []string{"composter", "smelter", "kiln", "refinery", "assembler", "reactor-lab", "silo", "generator", "lab"} {
		if FindNode(id) == nil {
			t.Errorf("missing building %s", id)
		}
	}
	if want := []Prereq{{"smelter", 3}, {"generator", 1}}; !reflect.DeepEqual(FindNode("assembler").After, want) {
		t.Errorf("assembler prerequisites = %v, want %v", FindNode("assembler").After, want)
	}
	if FindNode("missing") != nil {
		t.Fatal("unknown lookup returned a node")
	}
}

func TestInvalidTreesAreRejected(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree []Node
	}{
		{"cycle", []Node{{ID: "base"}, {ID: "a", After: []Prereq{{"b", 1}, {"base", 1}}}, {ID: "b", After: []Prereq{{"a", 1}}}}},
		{"unknown parent", []Node{{ID: "base"}, {ID: "a", After: []Prereq{{"missing", 1}}}}},
		{"unreachable", []Node{{ID: "base"}, {ID: "island"}}},
		{"earlier chapter", []Node{{ID: "base", Chapter: 2}, {ID: "a", Chapter: 1, After: []Prereq{{"base", 1}}}}},
		{"missing root", []Node{{ID: "a"}}},
		{"duplicate", []Node{{ID: "base"}, {ID: "base"}}},
		{"root has parent", []Node{{ID: "base", After: []Prereq{{"a", 1}}}, {ID: "a"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := Validate(tc.tree); err == nil {
				t.Fatal("invalid tree accepted")
			}
		})
	}
}

func TestRequirementsReturnEveryConditionInOrder(t *testing.T) {
	n := Node{ID: "assembler", After: []Prereq{{"smelter", 3}, {"generator", 1}}, Contract: "assembly-contract", Rank: 2, RP: 5,
		Cost: CostRule{Credits: 10, Growth: 1, Tiers: [][]string{{"wire", "plate"}}, Qty: 2, QtyGrowth: 1}}
	want := []Check{{"parent", "smelter", false}, {"parent", "generator", false}, {"contract", "assembly-contract", false}, {"rank", "rank", false}, {"rp", "rp", false}, {"cost", "credits", false}, {"cost", "plate", false}, {"cost", "wire", false}}
	if got := Requirements(n, map[string]int{"smelter": 2}, nil, 0, 0, 0, nil); !reflect.DeepEqual(got, want) {
		t.Fatalf("checks = %+v, want %+v", got, want)
	}
	checks := Requirements(n, map[string]int{"smelter": 3, "generator": 1}, map[string]bool{"assembler": true}, 2, 5, 10, map[string]float64{"plate": 2, "wire": 2})
	if len(checks) != len(want) {
		t.Fatalf("satisfied checks omitted: %+v", checks)
	}
	for _, c := range checks {
		if !c.OK {
			t.Errorf("condition not satisfied at exact threshold: %+v", c)
		}
	}
}

func TestRepeatPurchasesOnlyNeedRankAndCost(t *testing.T) {
	n := Node{ID: "test", After: []Prereq{{"missing", 1}}, Contract: "missing", Rank: 1, RP: 100, Cost: nodeCost(10, 1)}
	for _, rank := range []int{1, 2} {
		checks := Requirements(n, map[string]int{"test": 5}, nil, rank, 0, 10, map[string]float64{"biomass": 100})
		want := []Check{{"rank", "rank", rank == 2}, {"cost", "credits", true}, {"cost", "biomass", true}}
		if !reflect.DeepEqual(checks, want) {
			t.Errorf("rank %d: checks = %+v, want %+v", rank, checks, want)
		}
	}
}

func TestFiveNodeStatesOnOneSave(t *testing.T) {
	previous := Tree
	Tree = []Node{
		{ID: "base"},
		{ID: "owned", After: []Prereq{{"base", 1}}},
		{ID: "available", After: []Prereq{{"owned", 1}}},
		{ID: "near", After: []Prereq{{"owned", 1}}, Contract: "contract", Rank: 1, RP: 5, Cost: nodeCost(10, 1)},
		{ID: "locked", After: []Prereq{{"owned", 2}}},
		{ID: "step", After: []Prereq{{"available", 1}}},
		{ID: "unknown", After: []Prereq{{"step", 1}}},
	}
	t.Cleanup(func() { Tree = previous })
	levels := map[string]int{"owned": 1}
	for _, tc := range []struct{ id, want string }{{"base", "owned"}, {"owned", "owned"}, {"available", "available"}, {"near", "near"}, {"locked", "locked"}, {"step", "locked"}, {"unknown", "unknown"}} {
		if got := NodeState(*FindNode(tc.id), levels, nil, 0, 0, 0, nil); got != tc.want {
			t.Errorf("%s state = %s, want %s", tc.id, got, tc.want)
		}
	}
	if got := Requirements(Node{After: []Prereq{{"base", 1}}}, nil, nil, 0, 0, 0, nil); !got[0].OK {
		t.Fatal("base should be owned in an empty save")
	}
}
