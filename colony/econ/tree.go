package econ

import (
	"fmt"
	"sort"
	"strings"

	"orbital/colony/i18n"
)

type Prereq struct {
	ID    string
	Level int
}

type Node struct {
	ID, Branch                 string
	Name, NameEn, Desc, DescEn string
	After                      []Prereq
	Contract                   string
	Chapter                    int
	Rank                       int
	RP                         float64
	Cost                       CostRule
}

var Tree []Node

// Catalog is populated by econ.go's init before this file's init runs.
func init() {
	Tree = []Node{{ID: "base", Name: "القاعدة", NameEn: "Base", Desc: "بداية المستعمرة.", DescEn: "The colony's starting point.", Chapter: 1}}
	for _, u := range Catalog {
		if PlacementCost(u.ID) == nil {
			Tree = append(Tree, upgradeNode(u))
		}
	}
	Tree = append(Tree, buildingNodes()...)
	addTreeContent()
	for i := range Tree {
		Tree[i].Chapter = nodeChapter(Tree[i])
		if Tree[i].ID != "base" && (len(Tree[i].After) != 1 || Tree[i].After[0].ID != "base") {
			Tree[i].RP = float64(max(8, Tree[i].Chapter*2))
		}
	}
}

func upgradeNode(u *Upgrade) Node {
	n := Node{ID: u.ID, Branch: upgradeBranch(u.ID), Name: u.Name, NameEn: u.NameEn, Desc: u.Desc, DescEn: u.DescEn, Chapter: max(1, u.Lesson)}
	for _, id := range u.Requires {
		n.After = append(n.After, Prereq{id, 1})
	}
	if len(n.After) == 0 {
		parent := "field"
		switch u.ID {
		case "field", "engine", "yield.moss":
			parent = "base"
		case "grow.moss", "api.prices":
			parent = "yield.moss"
		case "drill":
			parent = "api.dig"
		}
		n.After = []Prereq{{parent, 1}}
	}
	c := u.CostAt(0)
	n.Cost = nodeCost(c.Credits, u.CostAt(1).Credits/c.Credits)
	for good, qty := range c.Goods {
		n.Cost.Tiers[0], n.Cost.Qty = []string{good}, qty
	}
	if good, ok := strings.CutPrefix(u.ID, "yield."); ok {
		n.Cost.Tiers[0] = []string{good}
	}
	if good, ok := strings.CutPrefix(u.ID, "grow."); ok {
		n.Cost.Tiers[0] = []string{good}
	}
	return n
}

func upgradeBranch(id string) string {
	switch {
	case id == "engine", id == "hangar", id == "api.spawn":
		return "الأسطول"
	case id == "api.scanner":
		return "البيانات"
	case id == "drill", id == "api.dig", id == "api.measure", id == "api.radar", id == "api.swap", id == "crop.crystal", strings.HasPrefix(id, "veins."), id == "yield.iron", id == "yield.geode", id == "yield.crystal":
		return "التعدين"
	case id == "api.prices", strings.HasPrefix(id, "market."):
		return "الخدمات"
	default:
		return "الزراعة"
	}
}

func nodeCost(credits, growth float64) CostRule {
	return CostRule{Credits: credits, Growth: growth, Tiers: [][]string{{"moss"}, {"biomass"}, {"fertilizer"}, {"processor"}, {"core"}}, Qty: 5, QtyGrowth: 1.15}
}

func buildingNodes() []Node {
	return []Node{
		{ID: "composter", Branch: "التصنيع", Name: "السمّادة", NameEn: "Composter", Desc: "تحوّل المحاصيل لكتلة حيوية.", DescEn: "Turns crops into biomass.", After: []Prereq{{"yield.moss", 1}}, Chapter: 7, Cost: nodeCost(100, 1.6)},
		{ID: "smelter", Branch: "التصنيع", Name: "المصهر", NameEn: "Smelter", Desc: "يصنع الصفائح والأسلاك والسبائك.", DescEn: "Makes plates, wire, and alloy.", After: []Prereq{{"drill", 1}}, Chapter: 17, Cost: nodeCost(200, 1.6)},
		{ID: "kiln", Branch: "التصنيع", Name: "الفرن", NameEn: "Kiln", Desc: "يحوّل السيليكا لزجاج.", DescEn: "Turns silica into glass.", After: []Prereq{{"crop.crystal", 1}}, Chapter: 7, Cost: nodeCost(180, 1.6)},
		{ID: "refinery", Branch: "التصنيع", Name: "المصفاة", NameEn: "Refinery", Desc: "تصنع الوقود والسماد وسائل التبريد.", DescEn: "Makes fuel, fertilizer, and coolant.", After: []Prereq{{"composter", 1}, {"kiln", 1}}, Chapter: 9, Cost: nodeCost(300, 1.6)},
		{ID: "assembler", Branch: "التصنيع", Name: "المجمّع", NameEn: "Assembler", Desc: "يجمع المكونات لمنتجات متقدمة.", DescEn: "Assembles components into advanced goods.", After: []Prereq{{"smelter", 3}, {"generator", 1}}, Chapter: 19, Cost: nodeCost(600, 1.6)},
		{ID: "reactor-lab", Branch: "التصنيع", Name: "معمل المفاعل", NameEn: "Reactor Lab", Desc: "يصنع خلايا الطاقة من اليورانيوم.", DescEn: "Makes energy cells from uranium.", After: []Prereq{{"assembler", 1}, {"drill", 4}}, Chapter: 25, Cost: nodeCost(2000, 1.6)},
		{ID: "silo", Branch: "التصنيع", Name: "الصومعة", NameEn: "Silo", Desc: "توسّع سعة تخزين الموارد.", DescEn: "Expands resource storage.", After: []Prereq{{"field", 1}, {"composter", 1}}, Chapter: 7, Cost: nodeCost(150, 1.6)},
		{ID: "generator", Branch: "الطاقة", Name: "المولّد", NameEn: "Generator", Desc: "يمد المباني بالطاقة.", DescEn: "Supplies buildings with power.", After: []Prereq{{"engine", 1}, {"kiln", 1}}, Chapter: 19, Cost: nodeCost(400, 1.6)},
		{ID: "lab", Branch: "البيانات", Name: "المعمل", NameEn: "Lab", Desc: "يستهلك المنتجات لإنتاج نقط البحث.", DescEn: "Consumes goods to produce research points.", After: []Prereq{{"assembler", 1}, {"api.scanner", 1}}, Chapter: 19, Cost: nodeCost(800, 1.6)},
	}
}

func nodeChapter(n Node) int {
	chapter := n.Chapter
	for _, p := range n.After {
		chapter = max(chapter, nodeChapter(*FindNode(p.ID)))
	}
	return chapter
}

func FindNode(id string) *Node {
	for i := range Tree {
		if Tree[i].ID == id {
			return &Tree[i]
		}
	}
	return nil
}

func Validate(tree []Node) error {
	nodes := make(map[string]Node, len(tree))
	for _, n := range tree {
		if _, exists := nodes[n.ID]; exists {
			return fmt.Errorf(i18n.T("عقدة مكررة: %s", "Duplicate node: %s"), n.ID)
		}
		nodes[n.ID] = n
	}
	if _, ok := nodes["base"]; !ok {
		return fmt.Errorf("%s", i18n.T("جذر القاعدة مش موجود", "Missing base root"))
	}
	if len(nodes["base"].After) != 0 {
		return fmt.Errorf("%s", i18n.T("القاعدة مينفعش يكون ليها أب", "Base cannot have parents"))
	}
	visited := map[string]int{}
	for _, n := range tree {
		if err := validatePath(n.ID, nodes, visited); err != nil {
			return err
		}
	}
	return nil
}

func validatePath(id string, nodes map[string]Node, visited map[string]int) error {
	if visited[id] == 1 {
		return fmt.Errorf(i18n.T("دورة عند العقدة: %s", "Cycle at node: %s"), id)
	}
	if visited[id] == 2 {
		return nil
	}
	visited[id] = 1
	n := nodes[id]
	if id != "base" && len(n.After) == 0 {
		return fmt.Errorf(i18n.T("عقدة ملهاش طريق للقاعدة: %s", "Node unreachable from base: %s"), id)
	}
	for _, p := range n.After {
		parent, ok := nodes[p.ID]
		if !ok {
			return fmt.Errorf(i18n.T("أب مش معروف: %s", "Unknown prerequisite: %s"), p.ID)
		}
		if n.Chapter < parent.Chapter {
			return fmt.Errorf(i18n.T("فصل العقدة %s قبل الأب %s", "Node %s has an earlier chapter than parent %s"), id, p.ID)
		}
		if err := validatePath(p.ID, nodes, visited); err != nil {
			return err
		}
	}
	visited[id] = 2
	return nil
}

type Check struct {
	Kind, Ref string
	OK        bool
}

func Requirements(n Node, levels map[string]int, blueprints map[string]bool, rank int, rp float64, credits float64, inv map[string]float64) []Check {
	level := levels[n.ID]
	var checks []Check
	if level == 0 {
		for _, p := range n.After {
			checks = append(checks, Check{"parent", p.ID, levels[p.ID] >= p.Level || p.ID == "base" && p.Level <= 1})
		}
		if n.Contract != "" {
			checks = append(checks, Check{"contract", n.Contract, blueprints[n.ID]})
		}
	}
	checks = append(checks, Check{"rank", "rank", rank >= max(n.Rank, 2*(level/5))})
	if level == 0 {
		checks = append(checks, Check{"rp", "rp", rp >= n.RP})
	}
	cost := n.Cost.At(level)
	checks = append(checks, Check{"cost", "credits", credits >= cost.Credits})
	goods := make([]string, 0, len(cost.Goods))
	for good := range cost.Goods {
		goods = append(goods, good)
	}
	sort.Strings(goods)
	for _, good := range goods {
		checks = append(checks, Check{"cost", good, inv[good] >= cost.Goods[good]})
	}
	return checks
}

// NodeState uses Tree for graph distance and treats base as owned even in an empty save.
func NodeState(n Node, levels map[string]int, blueprints map[string]bool, rank int, rp float64, credits float64, inv map[string]float64) string {
	if n.ID == "base" || levels[n.ID] > 0 {
		return "owned"
	}
	if !nearOwned(n, levels, 2) {
		return "unknown"
	}
	state := "available"
	for _, c := range Requirements(n, levels, blueprints, rank, rp, credits, inv) {
		if !c.OK {
			if c.Kind == "parent" {
				return "locked"
			}
			state = "near"
		}
	}
	return state
}

func nearOwned(n Node, levels map[string]int, steps int) bool {
	if n.ID == "base" || levels[n.ID] > 0 {
		return true
	}
	if steps == 0 {
		return false
	}
	for _, p := range n.After {
		if parent := FindNode(p.ID); parent != nil && nearOwned(*parent, levels, steps-1) {
			return true
		}
	}
	return false
}
