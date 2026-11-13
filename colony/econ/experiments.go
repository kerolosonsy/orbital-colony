package econ

import "sort"

func init() {
	Contracts = append(Contracts, experiments()...)
	sort.SliceStable(Contracts, func(i, j int) bool { return Contracts[i].Chapter < Contracts[j].Chapter })
}

func experiments() []*Contract {
	return []*Contract{
		{ID: "e03-sweep", Chapter: 3, After: []string{"c03-full-sweep"}, Requires: []string{"for2"}, Window: 750, Repeatable: true, RewardRP: 8,
			Title: "تجربة المسح", TitleEn: "Sweep Experiment", Brief: "سلّم 40 طحالب في 750 tick بحلقتين متداخلتين. مكافأة البحث بتقل مع تكرار التجربة في نفس اليوم.", BriefEn: "Deliver 40 moss within 750 ticks using nested loops. Research rewards decrease for repeats in the same day.", Goal: Goal{Deliver: deliver("moss", 40), MaxLost: NoLimit}},
		{ID: "e05-care", Chapter: 5, After: []string{"c05-careful-harvest"}, Requires: []string{"func", "errcheck"}, Window: 700, Repeatable: true, RewardRP: 10,
			Title: "تجربة الحصاد الحذر", TitleEn: "Careful Harvest Experiment", Brief: "سلّم 30 من أي محصول في 700 tick من غير خسارة نبات. السكانر مقفول؛ استخدم TryHarvest وعالج الخطأ.", BriefEn: "Deliver 30 of any crop within 700 ticks without losing plants. The scanner is disabled; use TryHarvest and handle errors.", Goal: Goal{Any: 30, MaxLost: 0}, Limits: Limits{Deny: []string{"api.scanner"}}},
		{ID: "e07-cluster", Chapter: 7, After: []string{"c07-plot-type"}, Requires: []string{"struct", "method"}, Window: 650, Repeatable: true, RewardRP: 12,
			Title: "تجربة العناقيد", TitleEn: "Cluster Experiment", Brief: "سلّم 20 فطر مضيء في 650 tick بمتوسط حصاد 3 على الأقل. مثّل القطعة بـ struct وmethod.", BriefEn: "Deliver 20 glowcap within 650 ticks with average base yield at least 3. Model a plot with a struct and method.", Goal: Goal{Deliver: deliver("glowcap", 20), MinAvg: map[string]float64{"glowcap": 3}, MaxLost: NoLimit}},
		{ID: "e08-veins", Chapter: 8, After: []string{"c08-big-veins"}, Requires: []string{"pointer", "nilcheck"}, Window: 1000, Repeatable: true, RewardRP: 14,
			Title: "تجربة العروق", TitleEn: "Vein Experiment", Brief: "سلّم 16 كريستال في 1000 tick من عروق حجمها 4 أو أكبر. استخدم مؤشر وافحص nil.", BriefEn: "Deliver 16 crystal within 1000 ticks, harvesting only veins of size 4 or greater. Use a pointer and check nil.", Goal: Goal{Deliver: deliver("crystal", 16), MinVein: 4, MaxLost: NoLimit}},
		{ID: "e15-route", Chapter: 15, After: []string{"c15-shortest-run"}, Requires: []string{"recursion|sel:heap.Push"}, Window: 1000, Repeatable: true, RewardRP: 20,
			Title: "تجربة المسار", TitleEn: "Route Experiment", Brief: "سلّم 24 كريستال في 1000 tick وفي 120 حركة بحد أقصى. خطّط المسار بـ recursion أو heap.", BriefEn: "Deliver 24 crystal within 1000 ticks and at most 120 moves. Plan the route with recursion or a heap.", Goal: Goal{Deliver: deliver("crystal", 24), MaxMoves: 120, MaxLost: NoLimit}},
		{ID: "e17-depth", Chapter: 17, After: []string{"c17-deep-haul"}, Requires: []string{"for2", "slice"}, Window: 100, Repeatable: true, RewardRP: 24,
			Title: "تجربة الأعماق", TitleEn: "Depth Experiment", Brief: "سلّم 6 جيود في 100 tick. استخدم slice لتتبع الحفر وحلقتين لمسح الطبقة.", BriefEn: "Deliver 6 geodes within 100 ticks. Track excavation with a slice and sweep the layer with nested loops.", Goal: Goal{Deliver: deliver("geode", 6), MaxLost: NoLimit}},
	}
}
