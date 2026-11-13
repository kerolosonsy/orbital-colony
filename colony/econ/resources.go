package econ

type Resource struct {
	ID, Name, NameEn string
	Tier             int
	BasePrice        float64
}

// New raw prices rise from surface gathering through deeper, hazardous mining and night-only gas collection.
var RawV3 = []Resource{
	{"water", "المياه", "Water", 0, 2},
	{"silica", "السيليكا", "Silica", 0, 5},
	{"copper", "النحاس", "Copper", 0, 12},
	{"sulfur", "الكبريت", "Sulfur", 0, 18},
	{"titanium", "التيتانيوم", "Titanium", 0, 30},
	{"helium", "الهيليوم", "Helium", 0, 45},
	{"uranium", "اليورانيوم", "Uranium", 0, 90},
}

var Crafted = []Resource{
	{"biomass", "الكتلة الحيوية", "Biomass", 1, 6},
	{"plate", "الصفيحة", "Plate", 1, 29},
	{"wire", "السلك", "Wire", 1, 38},
	{"glass", "الزجاج", "Glass", 1, 24},
	{"fuel", "الوقود", "Fuel", 1, 22},
	{"fertilizer", "السماد", "Fertilizer", 2, 48},
	{"circuit", "الدائرة الإلكترونية", "Circuit", 2, 170},
	{"alloy", "السبيكة", "Alloy", 2, 141},
	{"battery", "البطارية", "Battery", 2, 136},
	{"coolant", "سائل التبريد", "Coolant", 2, 78},
	{"processor", "المعالج", "Processor", 3, 880},
	{"frame", "الهيكل", "Frame", 3, 941},
	{"cell", "خلية الطاقة", "Cell", 3, 619},
	{"core", "النواة", "Core", 4, 6818},
}

type Recipe struct {
	Out      string
	In       map[string]float64
	Building string
	Ticks    int
	Alt      bool
}

var Recipes = []Recipe{
	{Out: "biomass", In: map[string]float64{"moss": 4}, Building: "composter", Ticks: 40},
	{Out: "biomass", In: map[string]float64{"lichen": 2}, Building: "composter", Ticks: 40, Alt: true},
	{Out: "plate", In: map[string]float64{"iron": 2}, Building: "smelter", Ticks: 60},
	{Out: "wire", In: map[string]float64{"copper": 2}, Building: "smelter", Ticks: 60},
	{Out: "glass", In: map[string]float64{"silica": 3}, Building: "kiln", Ticks: 80},
	{Out: "fuel", In: map[string]float64{"biomass": 2, "water": 1}, Building: "refinery", Ticks: 100},
	{Out: "fertilizer", In: map[string]float64{"biomass": 2, "sulfur": 1}, Building: "refinery", Ticks: 120},
	{Out: "circuit", In: map[string]float64{"wire": 2, "glass": 1, "crystal": 1}, Building: "assembler", Ticks: 200},
	{Out: "alloy", In: map[string]float64{"plate": 2, "titanium": 1}, Building: "smelter", Ticks: 180},
	{Out: "battery", In: map[string]float64{"wire": 1, "sulfur": 1, "plate": 1}, Building: "assembler", Ticks: 200},
	{Out: "coolant", In: map[string]float64{"water": 2, "helium": 1}, Building: "refinery", Ticks: 160},
	{Out: "processor", In: map[string]float64{"circuit": 3, "geode": 1}, Building: "assembler", Ticks: 500},
	{Out: "frame", In: map[string]float64{"alloy": 2, "circuit": 1, "battery": 1}, Building: "assembler", Ticks: 500},
	{Out: "cell", In: map[string]float64{"uranium": 1, "alloy": 1, "coolant": 2}, Building: "reactor-lab", Ticks: 800},
	{Out: "core", In: map[string]float64{"processor": 2, "frame": 2, "cell": 1}, Building: "assembler", Ticks: 2000},
}

// Biomass uses its primary moss recipe for pricing; the lichen alternative costs more inputs.
var PriceV3 = func() map[string]float64 {
	prices := make(map[string]float64, len(BasePrice)+len(RawV3)+len(Crafted))
	for good, price := range BasePrice {
		prices[good] = price
	}
	for _, goods := range [][]Resource{RawV3, Crafted} {
		for _, good := range goods {
			prices[good.ID] = good.BasePrice
		}
	}
	return prices
}()

func RecipeFor(out string) []Recipe {
	var recipes []Recipe
	for _, r := range Recipes {
		if r.Out == out {
			recipes = append(recipes, r)
		}
	}
	return recipes
}

func TierOf(good string) int {
	for _, g := range Crafted {
		if g.ID == good {
			return g.Tier
		}
	}
	return 0
}
