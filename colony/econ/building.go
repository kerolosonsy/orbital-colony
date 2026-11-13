package econ

import (
	"math"
	"sort"

	"orbital/colony/i18n"
)

type Building struct {
	speedBonus   float64
	bufferBonus  float64
	Kind         string             `json:"kind"`
	X            int                `json:"x"`
	Y            int                `json:"y"`
	Level        int                `json:"level"`
	Recipe       string             `json:"recipe"`
	In           map[string]float64 `json:"in"`
	Out          map[string]float64 `json:"out"`
	Progress     int                `json:"progress"`
	Wear         float64            `json:"wear"`
	Priority     int                `json:"priority"`
	Enabled      bool               `json:"enabled"`
	ActiveRecipe Recipe             `json:"activeRecipe"`
}

type BuildingError struct{ ar, en string }

func (e BuildingError) Error() string { return i18n.T(e.ar, e.en) }

// Messages lets the simulation carry both languages to live subscribers.
func (e BuildingError) Messages() (string, string) { return e.ar, e.en }

var (
	ErrNotHere           = BuildingError{"الدرون مش واقف على المبنى", "Drone is not standing on the building"}
	ErrUnknownKind       = BuildingError{"نوع المبنى مش معروف", "Unknown building kind"}
	ErrTileNotFree       = BuildingError{"الخانة مش متاحة للبناء", "Tile is not free for building"}
	ErrBlueprintNotOwned = BuildingError{"مخطط المبنى مش مملوك", "Building blueprint is not owned"}
	ErrNoInput           = BuildingError{"مفيش مدخلات كفاية", "Not enough inputs"}
	ErrFull              = BuildingError{"المخزن مليان", "Buffer is full"}
	ErrNoPower           = BuildingError{"مفيش طاقة", "No power"}
	ErrWorn              = BuildingError{"المبنى محتاج تصليح", "Building needs repair"}
	ErrBadRecipe         = BuildingError{"الوصفة مش صالحة للمبنى دلوقتي", "Recipe is not valid for this building right now"}
)

func Capacity(level int) float64 { return 8 * EffectMult(0.25, level) }

func (b *Building) recipes(out string) []Recipe {
	var recipes []Recipe
	for _, r := range RecipeFor(out) {
		if r.Building == b.Kind {
			recipes = append(recipes, r)
		}
	}
	sort.SliceStable(recipes, func(i, j int) bool { return !recipes[i].Alt && recipes[j].Alt })
	return recipes
}

func (b *Building) recipe(out string) (Recipe, error) {
	if b.Progress > 0 && b.ActiveRecipe.Out == out {
		return b.ActiveRecipe, nil
	}
	recipes := b.recipes(out)
	if len(recipes) == 0 {
		return Recipe{}, ErrBadRecipe
	}
	for _, r := range recipes {
		if b.hasInputs(r) {
			return r, nil
		}
	}
	return recipes[0], nil
}

func (b *Building) hasInputs(r Recipe) bool {
	for good, n := range r.In {
		if b.In[good] < n {
			return false
		}
	}
	return true
}

func (b *Building) accepts(good string) bool {
	if b.Kind == "lab" {
		return ResearchValue(good) > 0
	}
	for _, r := range b.recipes(b.Recipe) {
		if _, ok := r.In[good]; ok {
			return true
		}
	}
	return false
}

func (b *Building) SetRecipe(out string) error {
	if _, err := b.recipe(out); err != nil {
		return err
	}
	// Inputs already spent on an active cycle cannot become another product.
	if b.Progress > 0 && out != b.Recipe {
		return ErrBadRecipe
	}
	b.Recipe = out
	return nil
}

func (b *Building) Put(good string, n float64) (float64, error) {
	if !b.accepts(good) {
		return 0, ErrBadRecipe
	}
	if n <= 0 {
		return 0, nil
	}
	accepted := math.Min(n, b.capacity()-b.In[good])
	if accepted <= 0 {
		return 0, ErrFull
	}
	if b.In == nil {
		b.In = map[string]float64{}
	}
	b.In[good] += accepted
	return accepted, nil
}

func (b *Building) Take(good string, n float64) (float64, error) {
	if n <= 0 {
		return 0, nil
	}
	taken := math.Min(n, b.Out[good])
	if taken <= 0 {
		return 0, ErrNoInput
	}
	b.Out[good] -= taken
	return taken, nil
}

// Sorting keeps fractional buffer totals independent of Go's map iteration order.
func bufferAmount(buffer map[string]float64) float64 {
	goods := make([]string, 0, len(buffer))
	for good := range buffer {
		goods = append(goods, good)
	}
	sort.Strings(goods)
	total := 0.0
	for _, good := range goods {
		total += buffer[good]
	}
	return total
}

func (b *Building) ready(r Recipe, powered bool) error {
	if b.Wear >= 100 {
		return ErrWorn
	}
	if !powered {
		return ErrNoPower
	}
	if b.capacity()-bufferAmount(b.Out) < 1 {
		return ErrFull
	}
	if b.Progress == 0 && !b.hasInputs(r) {
		return ErrNoInput
	}
	return nil
}

func (b *Building) cycleTicks(r Recipe) int {
	ticks := max(1, int(math.Ceil(float64(r.Ticks)/(EffectMult(0.25, b.Level)*(1+b.speedBonus)))))
	if b.Wear > 80 {
		ticks *= 2
	}
	return ticks
}

func (b *Building) Step(ticks int, powered bool) error {
	if ticks <= 0 || !b.Enabled {
		return nil
	}
	for ticks > 0 {
		r, err := b.recipe(b.Recipe)
		if err != nil {
			return err
		}
		if err := b.ready(r, powered); err != nil {
			return err
		}
		ticks -= b.advanceCycle(r, ticks)
	}
	return nil
}

func (b *Building) advanceCycle(r Recipe, ticks int) int {
	if b.Progress == 0 {
		b.ActiveRecipe = r
		for good, n := range r.In {
			b.In[good] -= n
		}
	}
	duration := b.cycleTicks(r)
	elapsed := min(ticks, max(0, duration-b.Progress))
	b.Progress += elapsed
	if b.Progress >= duration {
		if b.Out == nil {
			b.Out = map[string]float64{}
		}
		b.Out[r.Out]++
		b.Progress = 0
		b.ActiveRecipe = Recipe{}
		b.Wear = math.Min(100, b.Wear+1)
	}
	return elapsed
}

func (b *Building) Repair(plates float64) float64 {
	if plates <= 0 || b.Wear <= 0 {
		return 0
	}
	used := math.Min(plates, b.Wear/20)
	b.Wear = math.Max(0, b.Wear-20*used)
	return used
}

// ProgressFraction is shared by status queries and live building events.
func (b *Building) ProgressFraction() float64 {
	if b.Progress == 0 {
		return 0
	}
	if b.Kind == "lab" {
		return math.Min(1, float64(b.Progress)/float64(b.cycleTicks(Recipe{Ticks: ResearchTicks})))
	}
	r, err := b.recipe(b.Recipe)
	if err != nil {
		return 0
	}
	return math.Min(1, float64(b.Progress)/float64(b.cycleTicks(r)))
}

// PlacementCost is independent of the price of upgrading the blueprint.
func PlacementCost(kind string) *Cost {
	cost, ok := placementCosts[kind]
	if !ok {
		return nil
	}
	return &cost
}

var placementCosts = map[string]Cost{
	"composter":   {Credits: 25, Goods: map[string]float64{"moss": 4}},
	"smelter":     {Credits: 50, Goods: map[string]float64{"iron": 2}},
	"kiln":        {Credits: 40, Goods: map[string]float64{"silica": 3}},
	"refinery":    {Credits: 75, Goods: map[string]float64{"biomass": 2}},
	"assembler":   {Credits: 150, Goods: map[string]float64{"plate": 2}},
	"reactor-lab": {Credits: 500, Goods: map[string]float64{"plate": 4}},
	"silo":        {Credits: 35, Goods: map[string]float64{"moss": 4}},
	"generator":   {Credits: 100, Goods: map[string]float64{"glass": 2}},
	"lab":         {Credits: 200, Goods: map[string]float64{"plate": 2}},
}

func addBuildingUpgrades() {
	goods := map[string]string{"composter": "moss", "smelter": "iron", "kiln": "silica", "refinery": "biomass", "assembler": "plate", "reactor-lab": "alloy", "silo": "moss", "generator": "glass", "lab": "circuit"}
	for _, n := range buildingNodes() {
		var requires []string
		for _, p := range n.After {
			requires = append(requires, p.ID)
		}
		add(&Upgrade{ID: n.ID, Group: "المباني", GroupEn: "Buildings", Name: n.Name, NameEn: n.NameEn, Desc: n.Desc, DescEn: n.DescEn, Max: 8, Requires: requires, Lesson: n.Chapter, cost: mixed(n.Cost.Credits, 1.6, goods[n.ID], 5, 1.15)})
	}
}
