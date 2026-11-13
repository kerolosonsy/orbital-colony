package econ

import (
	"math"
	"sort"

	"orbital/colony/i18n"
)

type Building struct {
	Kind     string
	X, Y     int
	Level    int
	Recipe   string
	In, Out  map[string]float64
	Progress int
	Wear     float64
	Priority int
	Enabled  bool

	activeRecipe Recipe
}

type buildingError struct{ ar, en string }

func (e buildingError) Error() string { return i18n.T(e.ar, e.en) }

var (
	ErrNoInput   = buildingError{"مفيش مدخلات كفاية", "Not enough inputs"}
	ErrFull      = buildingError{"المخزن مليان", "Buffer is full"}
	ErrNoPower   = buildingError{"مفيش طاقة", "No power"}
	ErrWorn      = buildingError{"المبنى محتاج تصليح", "Building needs repair"}
	ErrBadRecipe = buildingError{"الوصفة مش صالحة للمبنى دلوقتي", "Recipe is not valid for this building right now"}
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
	if b.Progress > 0 && b.activeRecipe.Out == out {
		return b.activeRecipe, nil
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
	accepted := math.Min(n, Capacity(b.Level)-b.In[good])
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
	if Capacity(b.Level)-bufferAmount(b.Out) < 1 {
		return ErrFull
	}
	if b.Progress == 0 && !b.hasInputs(r) {
		return ErrNoInput
	}
	return nil
}

func (b *Building) cycleTicks(r Recipe) int {
	ticks := max(1, int(math.Ceil(float64(r.Ticks)/EffectMult(0.25, b.Level))))
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
		b.activeRecipe = r
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
		b.activeRecipe = Recipe{}
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
