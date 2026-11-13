package econ

import "fmt"

func addTreeContent() {
	addArchiveMarketNodes()
	addNeighborTree()
	addResourceResearch()
	addBuildingResearch()
	for _, n := range []Node{
		{ID: "research.engine", Branch: "الأسطول", Name: "ضبط المحرك", NameEn: "Engine Tuning", Desc: "سرعة الأوامر تزيد 8% لكل مستوى.", DescEn: "Command speed increases by 8% per level.", After: []Prereq{{"engine", 3}}, Chapter: 10},
		{ID: "research.vein-rate", Branch: "التعدين", Name: "استكشاف العروق", NameEn: "Vein Prospecting", Desc: "معدل ظهور العروق يزيد 20% لكل مستوى.", DescEn: "Vein spawn rate increases by 20% per level.", After: []Prereq{{"veins.rate", 2}, {"api.radar", 1}}, Chapter: 8},
		{ID: "research.vein-size", Branch: "التعدين", Name: "امتداد العروق", NameEn: "Vein Extension", Desc: "حجم العرق الأقصى يزيد واحد لكل مستوى.", DescEn: "Maximum vein size increases by one per level.", After: []Prereq{{"veins.size", 2}}, Chapter: 8},
		{ID: "research.sift-stock", Branch: "التعدين", Name: "مسح السيليكا", NameEn: "Silica Survey", Desc: "مخزون الغربلة يزيد وحدة لكل خانة لكل مستوى.", DescEn: "Sift stock increases by one unit per tile per level.", After: []Prereq{{"field", 2}}, Chapter: 7},
		{ID: "research.sift-refill", Branch: "التعدين", Name: "تجدد السيليكا", NameEn: "Silica Renewal", Desc: "معدل تجدد الغربلة يزيد 20% لكل مستوى.", DescEn: "Sift refill rate increases by 20% per level.", After: []Prereq{{"research.sift-stock", 1}}, Chapter: 8},
		{ID: "research.melt", Branch: "الزراعة", Name: "تجدد الجليد", NameEn: "Ice Renewal", Desc: "معدل رجوع الجليد بعد التسييح يزيد 20% لكل مستوى.", DescEn: "Ice regrowth rate after melting increases by 20% per level.", After: []Prereq{{"crop.lichen", 1}}, Chapter: 7},
		{ID: "research.placement", Branch: "التصنيع", Name: "بناء اقتصادي", NameEn: "Economical Construction", Desc: "تكلفة وضع المباني بالكريدت والموارد تقل 10% لكل مستوى.", DescEn: "Building placement credits and goods decrease by 10% per level.", After: []Prereq{{"composter", 2}, {"market.contracts", 1}}, Chapter: 13},
		{ID: "research.gain", Branch: "البيانات", Name: "تحليل البحث", NameEn: "Research Analysis", Desc: "نقط البحث من المعمل والتجارب تزيد 20% لكل مستوى.", DescEn: "Lab and experiment research points increase by 20% per level.", After: []Prereq{{"lab", 2}, {"api.scanner", 1}}, Chapter: 19},
	} {
		addResearchNode(n)
	}
	for i := range Tree {
		n := &Tree[i]
		if src := BlueprintSource(n.ID); src != nil && needsBlueprint(n.ID) {
			n.Contract = src.ID
		}
	}
}

func addResearchNode(n Node) {
	var requires []string
	for _, p := range n.After {
		requires = append(requires, p.ID)
	}
	u := &Upgrade{ID: n.ID, Group: n.Branch, GroupEn: "Research", Name: n.Name, NameEn: n.NameEn, Desc: n.Desc, DescEn: n.DescEn, Max: 3, Requires: requires, Lesson: n.Chapter, cost: credits(100, 1.8)}
	add(u)
	n.Cost = CostRule{Credits: 100, Growth: 1.8}
	Tree = append(Tree, n)
}

func addResourceResearch() {
	for _, good := range []string{"moss", "lichen", "glowcap", "reed", "copper", "sulfur", "titanium", "uranium", "water", "silica"} {
		branch, parent, chapter := "التعدين", "drill", 17
		parentLevel := 2
		switch good {
		case "moss", "lichen", "glowcap", "reed":
			branch, parent, chapter = "الزراعة", "yield."+good, 6
		case "water":
			branch, parent, chapter = "الزراعة", "crop.lichen", 7
		case "silica":
			parent, chapter = "field", 7
		case "sulfur":
			parentLevel = 3
		case "titanium":
			parentLevel = 4
		case "uranium":
			parentLevel = 5
		}
		if u := Find(parent); u.Max > 0 {
			parentLevel = min(parentLevel, u.Max)
		}
		addResearchNode(Node{ID: "research.yield." + good, Branch: branch, Name: "إنتاج متقدم: " + cropNames[good], NameEn: "Advanced Yield: " + cropNamesEn[good], Desc: "حصاد المورد ده يزيد 15% لكل مستوى.", DescEn: "This resource's gathered yield increases by 15% per level.", After: []Prereq{{parent, parentLevel}}, Chapter: chapter})
	}
	for _, good := range []string{"moss", "lichen", "glowcap", "reed"} {
		addResearchNode(Node{ID: "research.grow." + good, Branch: "الزراعة", Name: "نمو متقدم: " + cropNames[good], NameEn: "Advanced Growth: " + cropNamesEn[good], Desc: "معدل نمو المحصول ده يزيد 10% لكل مستوى.", DescEn: "This crop's growth rate increases by 10% per level.", After: []Prereq{{"grow." + good, 2}}, Chapter: 6})
		addResearchNode(Node{ID: "research.price." + good, Branch: "الخدمات", Name: "تسعير: " + cropNames[good], NameEn: "Pricing: " + cropNamesEn[good], Desc: "سعر المحصول ده يزيد 10% لكل مستوى.", DescEn: "This crop's selling price increases by 10% per level.", After: []Prereq{{"market.contracts", 2}, {"yield." + good, 2}}, Chapter: 13})
		addResearchNode(Node{ID: "research.recovery." + good, Branch: "الخدمات", Name: "طلب: " + cropNames[good], NameEn: "Demand: " + cropNamesEn[good], Desc: "معدل تعافي سوق المحصول ده يزيد 25% لكل مستوى.", DescEn: "This crop's market recovery rate increases by 25% per level.", After: []Prereq{{"market.logistics", 2}, {"yield." + good, 1}}, Chapter: 12})
	}
}

func addBuildingResearch() {
	for _, kind := range []string{"composter", "smelter", "kiln", "refinery", "assembler", "lab"} {
		building := FindNode(kind)
		for _, effect := range []struct{ id, ar, en, desc, descEn string }{
			{"speed", "سرعة", "Speed", "معدل تشغيل المبنى ده يزيد 20% لكل مستوى.", "This building kind's production speed increases by 20% per level."},
			{"buffer", "سعة", "Buffer", "سعة كل مدخل وإجمالي المخرجات للمبنى ده تزيد 4 لكل مستوى.", "Each input and total output buffer for this building kind increases by 4 per level."},
		} {
			if kind == "lab" && effect.id == "buffer" {
				effect.desc, effect.descEn = "سعة كل مدخل للمعمل تزيد 4 لكل مستوى.", "Each lab input buffer increases by 4 per level."
			}
			addResearchNode(Node{ID: "research." + effect.id + "." + kind, Branch: "التصنيع", Name: fmt.Sprintf("%s: %s", effect.ar, building.Name), NameEn: fmt.Sprintf("%s: %s", effect.en, building.NameEn), Desc: effect.desc, DescEn: effect.descEn, After: []Prereq{{kind, 2}}, Chapter: building.Chapter})
		}
	}
}
