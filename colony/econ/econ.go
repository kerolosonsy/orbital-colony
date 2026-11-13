// Package econ is the open colony's economy: its saved state, the upgrade
// catalog, the market and the rules for buying. It is pure data and pure
// functions, shared by the game server and the simulation.
package econ

import (
	"fmt"
	"math"

	"orbital/colony/i18n"
)

// Crops (and the crystal veins) that can be harvested and sold.
var Goods = []string{"moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode"}

// BasePrice is what one unit sells for when the market is calm.
var BasePrice = map[string]float64{"moss": 1, "lichen": 3, "glowcap": 2, "reed": 14, "crystal": 6, "iron": 9, "geode": 40}

// GrowTicks is how long each crop takes to ripen before upgrades. One tick
// is 100ms of real time while the colony runs.
var GrowTicks = map[string]int{"moss": 80, "lichen": 120, "glowcap": 140, "reed": 220}

// ReedWindow is how long a ripe reed stays good before it rots.
const ReedWindow = 90

type Tile struct {
	G string `json:"g"`
	E string `json:"e,omitempty"`
	P int64  `json:"p,omitempty"` // planted tick
	R int64  `json:"r,omitempty"` // ripe tick
	S int    `json:"s,omitempty"` // crystal size
}

// State is everything that persists between runs.
type State struct {
	Depth    int                `json:"depth,omitempty"`
	Under    [][]Tile           `json:"under,omitempty"`
	Version  int                `json:"version"`
	Seed     int64              `json:"seed"`
	Tick     int64              `json:"tick"`
	Size     int                `json:"size"`
	Tiles    []Tile             `json:"tiles"`
	Inv      map[string]float64 `json:"inv"`
	Credits  float64            `json:"credits"`
	Levels   map[string]int     `json:"levels"`
	Demand   map[string]float64 `json:"demand"` // 1 = calm market; selling pushes it down
	Claimed  map[string]bool    `json:"claimed"`
	Earned   float64            `json:"earned"` // lifetime credits
	RunTicks int64              `json:"runTicks"`

	Contracts  map[string]*ContractRun `json:"contracts,omitempty"`
	Active     string                  `json:"active,omitempty"` // the running contract's id
	Blueprints map[string]bool         `json:"blueprints,omitempty"`
}

func (s *State) At(x, y int) *Tile { return &s.Tiles[y*s.Size+x] }

// New returns a fresh colony: a 3x3 patch of regolith and nothing else.
func New(seed int64) *State {
	s := &State{Version: 2, Seed: seed, Inv: map[string]float64{}, Levels: map[string]int{}, Demand: map[string]float64{}, Claimed: map[string]bool{}}
	s.resize(FieldSize(0))
	s.Normalize()
	return s
}

// Normalize fills maps a decoded save may lack.
func (s *State) Normalize() {
	if s.Inv == nil {
		s.Inv = map[string]float64{}
	}
	if s.Levels == nil {
		s.Levels = map[string]int{}
	}
	if s.Demand == nil {
		s.Demand = map[string]float64{}
	}
	if s.Claimed == nil {
		s.Claimed = map[string]bool{}
	}
	if s.Contracts == nil {
		s.Contracts = map[string]*ContractRun{}
	}
	if s.Blueprints == nil {
		s.Blueprints = map[string]bool{}
	}
	if s.Version < 2 {
		s.migrateV2()
	}
	if s.Active != "" && FindContract(s.Active) == nil {
		s.Active = ""
	}
	if want := FieldSize(s.Levels["field"]); s.Size != want || len(s.Tiles) != want*want {
		s.resize(want)
	}
	s.normalizeLayers()
}

// groundAt is the terrain the colony finds when the field grows onto (x, y).
func groundAt(seed int64, x, y int) string {
	if x < 3 && y < 3 {
		return "regolith"
	}
	h := uint64(seed) ^ uint64(x)*0x9E3779B97F4A7C15 ^ uint64(y)*0xC2B2AE3D27D4EB4F
	h ^= h >> 31
	h *= 0xBF58476D1CE4E5B9
	h ^= h >> 29
	switch r := h % 100; {
	case r < 16:
		return "ice"
	case r < 28:
		return "rock"
	}
	return "regolith"
}

// resize grows the square field, keeping what is already there.
func (s *State) resize(n int) {
	old, on := s.Tiles, s.Size
	s.Tiles = make([]Tile, n*n)
	s.Size = n
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if x < on && y < on && y*on+x < len(old) {
				s.Tiles[y*n+x] = old[y*on+x]
			} else {
				s.Tiles[y*n+x] = Tile{G: groundAt(s.Seed, x, y)}
			}
		}
	}
	s.resizeUnder(n, on)
}

// ---------- the catalog ----------

type Cost struct {
	Credits float64            `json:"credits"`
	Goods   map[string]float64 `json:"goods,omitempty"`
}

type Upgrade struct {
	NameEn   string   `json:"nameEn"`
	DescEn   string   `json:"descEn"`
	GroupEn  string   `json:"groupEn"`
	ID       string   `json:"id"`
	Group    string   `json:"group"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Max      int      `json:"max"`
	Requires []string `json:"requires,omitempty"` // upgrade ids that must be at level >= 1
	Lesson   int      `json:"lesson,omitempty"`   // chapter that teaches what this unlocks
	cost     func(level int) Cost
}

// CostAt is the price of going from level to level+1.
func (u *Upgrade) CostAt(level int) Cost { return u.cost(level) }

func geo(base, growth float64) func(int) float64 {
	return func(l int) float64 { return math.Round(base * math.Pow(growth, float64(l))) }
}

func credits(base, growth float64) func(int) Cost {
	f := geo(base, growth)
	return func(l int) Cost { return Cost{Credits: f(l)} }
}

func mixed(base, growth float64, good string, gBase, gGrowth float64) func(int) Cost {
	f, g := geo(base, growth), geo(gBase, gGrowth)
	return func(l int) Cost { return Cost{Credits: f(l), Goods: map[string]float64{good: g(l)}} }
}

var fieldSizes = []int{3, 4, 5, 6, 7, 8, 10, 12, 14, 16, 19, 22, 26, 30}

// FieldSize is the side of the square field at a field-upgrade level.
func FieldSize(level int) int { return fieldSizes[min(level, len(fieldSizes)-1)] }

var Catalog []*Upgrade

func add(u *Upgrade) {
	switch u.ID {
	case "drill":
		u.NameEn, u.DescEn, u.GroupEn = "Drill", "Each level lets you dig one layer deeper.", "The Depths"
	case "api.dig":
		u.NameEn, u.DescEn, u.GroupEn = "Digging License", "Dig, DigDown, Descend, and Ascend.", "The Depths"
	}
	Catalog = append(Catalog, u)
}

// Find returns the upgrade with that id, or nil.
func Find(id string) *Upgrade {
	for _, u := range Catalog {
		if u.ID == id {
			return u
		}
	}
	return nil
}

var cropNames = map[string]string{"moss": "الطحالب", "lichen": "الأشنة", "glowcap": "الفطر المضيء", "reed": "البوص النابض", "crystal": "الكريستال", "iron": "الحديد", "geode": "الجيود"}

var cropNamesEn = map[string]string{"moss": "Moss", "lichen": "Lichen", "glowcap": "Glowcap", "reed": "Pulse Reed", "crystal": "Crystal", "iron": "Iron", "geode": "Geode"}

// CropName is the localized name of a good.
func CropName(g string) string { return i18n.T(cropNames[g], cropNamesEn[g]) }

func init() {
	addDepthUpgrades()
	// land and machines
	add(&Upgrade{ID: "field", Group: "الأرض", GroupEn: "Land", Name: "توسيع الأرض", NameEn: "Expand the Field", Desc: "المزرعة تكبر. الأرض الجديدة فيها ريجوليث وجليد وصخر.", DescEn: "Grow your farm. New land contains regolith, ice, and rock.", Max: len(fieldSizes) - 1, cost: credits(30, 2.25)})
	add(&Upgrade{ID: "engine", Group: "المحركات", GroupEn: "Engines", Name: "محرك أسرع", NameEn: "Faster Engine", Desc: "كل أوامر الدرون أسرع 8%.", DescEn: "All drone commands run 8% faster.", Max: 18, cost: credits(25, 1.62)})
	add(&Upgrade{ID: "hangar", Group: "الهانجر", GroupEn: "Hangar", Name: "درون زيادة", NameEn: "Extra Drone", Desc: "درون كمان تقدر تطلقه بـ colony.Spawn().", DescEn: "Another drone to launch with colony.Spawn().", Max: 15, Requires: []string{"api.spawn"}, Lesson: 10, cost: mixed(600, 1.9, "crystal", 8, 1.6)})

	// API licences: commands that only work once bought
	apis := []struct {
		id, name, desc string
		nameEn, descEn string
		lesson         int
		price          float64
		req            []string
	}{
		{"api.scanner", "السكانر", "CanHarvest و GetEntity و GetGround.", "Scanner", "CanHarvest, GetEntity, and GetGround.", 4, 25, nil},
		{"api.try", "الحصاد الحذر", "TryHarvest: بترجّع error بدل ما تدمّر النبات.", "Careful Harvesting", "TryHarvest: returns an error instead of destroying the plant.", 5, 60, []string{"api.scanner"}},
		{"api.prices", "شاشة البورصة", "colony.Price(colony.Moss): سعر البيع دلوقتي.", "Market Display", "colony.Price(colony.Moss): the current selling price.", 2, 90, nil},
		{"api.measure", "مقياس الكريستال", "Measure: حجم الكريستال قبل ما تحصده.", "Crystal Gauge", "Measure: check a crystal's size before harvesting.", 6, 250, []string{"crop.crystal"}},
		{"api.radar", "الرادار", "Nearest: أقرب كريستال، وTargets: كل العروق.", "Radar", "Nearest: the nearest crystal; Targets: every vein.", 8, 1500, []string{"crop.crystal"}},
		{"api.swap", "الذراع", "Swap: تبدّل محتوى خانتين.", "Manipulator Arm", "Swap: exchange the contents of two tiles.", 15, 900, []string{"field"}},
		{"api.spawn", "رخصة الهانجر", "colony.Spawn(): درونات تشتغل في goroutines.", "Hangar License", "colony.Spawn(): drones working in goroutines.", 10, 2500, []string{"engine"}},
	}
	for _, a := range apis {
		p := a.price
		add(&Upgrade{ID: a.id, Group: "الأوامر", GroupEn: "Commands", Name: a.name, NameEn: a.nameEn, Desc: a.desc, DescEn: a.descEn, Max: 1, Requires: a.req, Lesson: a.lesson, cost: func(int) Cost { return Cost{Credits: p} }})
	}

	// new crops
	add(&Upgrade{ID: "crop.lichen", Group: "المحاصيل", GroupEn: "Crops", Name: "بذور الأشنة", NameEn: "Lichen Seeds", Desc: "الأشنة بتتزرع على الجليد. أبطأ من الطحالب بس سعرها 3 أضعاف.", DescEn: "Lichen grows on ice. Slower than moss, but worth 3 times as much.", Max: 1, Requires: []string{"field"}, cost: func(int) Cost { return Cost{Credits: 60} }})
	add(&Upgrade{ID: "crop.glowcap", Group: "المحاصيل", GroupEn: "Crops", Name: "الفطر المضيء", NameEn: "Glowcap", Desc: "بيتزرع على التربة. كل فطرة ناضجة جنبه بتزوّد حصاده: ازرعهم عناقيد.", DescEn: "Grows on soil. Each ripe neighbor increases its yield: plant in clusters.", Max: 1, Requires: []string{"crop.lichen"}, cost: func(int) Cost {
		return Cost{Credits: 450, Goods: map[string]float64{"lichen": 40}}
	}})
	add(&Upgrade{ID: "crop.reed", Group: "المحاصيل", GroupEn: "Crops", Name: "البوص النابض", NameEn: "Pulse Reed", Desc: "لازم يبقى جنب جليد. غالي جداً، بس لو اتأخرت عليه بعد ما ينضج بيعفّن.", DescEn: "Must grow beside ice. Very valuable, but rots if left too long after ripening.", Max: 1, Requires: []string{"crop.glowcap"}, cost: func(int) Cost {
		return Cost{Credits: 2500, Goods: map[string]float64{"glowcap": 120}}
	}})
	add(&Upgrade{ID: "crop.crystal", Group: "المحاصيل", GroupEn: "Crops", Name: "تعدين الكريستال", NameEn: "Crystal Mining", Desc: "عروق كريستال بتطلع على الصخر وبتكبر لو سبتها. استنّى ولا احصد؟", DescEn: "Crystal veins appear on rock and grow if left alone. Wait or harvest?", Max: 1, Requires: []string{"field"}, cost: func(int) Cost { return Cost{Credits: 350} }})

	// yields and growth per crop
	yieldBase := map[string]float64{"moss": 40, "lichen": 120, "glowcap": 300, "reed": 1200, "crystal": 500, "iron": 800, "geode": 2000}
	growBase := map[string]float64{"moss": 35, "lichen": 110, "glowcap": 260, "reed": 900}
	for _, g := range Goods {
		var req []string
		if g != "moss" {
			req = []string{"crop." + g}
		}
		group := "الإنتاج"
		groupEn := "Production"
		if g == "iron" || g == "geode" {
			req, group = []string{"api.dig"}, "الأعماق"
			groupEn = "The Depths"
		}
		add(&Upgrade{ID: "yield." + g, Group: group, GroupEn: groupEn, Name: "حصاد " + cropNames[g], NameEn: "Harvest " + cropNamesEn[g], Desc: "+35% من كل حصاد لـ " + cropNames[g] + ".", DescEn: "+35% from each harvest of " + cropNamesEn[g] + ".", Max: 12, Requires: req, cost: credits(yieldBase[g], 2.05)})
		if gb, ok := growBase[g]; ok {
			add(&Upgrade{ID: "grow." + g, Group: "النمو", GroupEn: "Growth", Name: "نمو " + cropNames[g], NameEn: "Grow " + cropNamesEn[g], Desc: cropNames[g] + " بتنضج أسرع 10%.", DescEn: cropNamesEn[g] + " ripens 10% faster.", Max: 8, Requires: req, cost: credits(gb, 1.95)})
		}
	}

	// crystals and the market
	add(&Upgrade{ID: "veins.rate", Group: "الكريستال", GroupEn: "Crystal", Name: "عروق أكتر", NameEn: "More Veins", Desc: "عروق الكريستال بتطلع أسرع 20%.", DescEn: "Crystal veins appear 20% faster.", Max: 8, Requires: []string{"crop.crystal"}, cost: credits(400, 1.9)})
	add(&Upgrade{ID: "veins.size", Group: "الكريستال", GroupEn: "Crystal", Name: "عروق أعمق", NameEn: "Deeper Veins", Desc: "العرق يقدر يكبر لحجم أكبر بـ 1.", DescEn: "Increases the maximum vein size by 1.", Max: 6, Requires: []string{"crop.crystal"}, cost: mixed(700, 2.1, "crystal", 20, 1.7)})
	add(&Upgrade{ID: "market.contracts", Group: "السوق", GroupEn: "Market", Name: "عقود تصدير", NameEn: "Export Contracts", Desc: "كل البيع بسعر أعلى 10%.", DescEn: "Sell everything for 10% more.", Max: 12, Requires: []string{"api.prices"}, cost: credits(150, 2.0)})
	add(&Upgrade{ID: "market.logistics", Group: "السوق", GroupEn: "Market", Name: "شحن أسرع", NameEn: "Faster Shipping", Desc: "السوق بيرجع لسعره الطبيعي أسرع بعد ما تبيع كتير.", DescEn: "Market prices recover faster after you sell large quantities.", Max: 8, Requires: []string{"api.prices"}, cost: credits(200, 1.95)})
}

// ---------- effects ----------

func (s *State) Level(id string) int { return s.Levels[id] }

// Has reports whether an upgrade (or an always-free command) is available.
func (s *State) Has(id string) bool { return s.Levels[id] > 0 }

// ActionMillis is how long one drone action takes in real time.
func (s *State) ActionMillis() float64 { return 260 * math.Pow(0.92, float64(s.Level("engine"))) }

// Drones is how many drones the hangar holds, the first one included.
func (s *State) Drones() int { return 1 + s.Level("hangar") }

func (s *State) YieldMult(g string) float64 { return 1 + 0.35*float64(s.Level("yield."+g)) }

func (s *State) GrowTime(g string) int64 {
	return int64(math.Max(10, float64(GrowTicks[g])*math.Pow(0.9, float64(s.Level("grow."+g)))))
}

func (s *State) VeinEvery() int64 { return int64(160 * math.Pow(0.8, float64(s.Level("veins.rate")))) }
func (s *State) VeinMax() int     { return 4 + s.Level("veins.size") }

// Price is what one unit of g sells for right now.
func (s *State) Price(g string) float64 {
	d := s.Demand[g]
	if d == 0 {
		d = 1
	}
	return BasePrice[g] * (1 + 0.1*float64(s.Level("market.contracts"))) * d
}

// Sell turns n units of g into credits, pushing that good's price down a
// little per unit. It returns the credits earned.
func (s *State) Sell(g string, n float64) float64 {
	n = math.Min(n, s.Inv[g])
	if n <= 0 {
		return 0
	}
	earned := 0.0
	for n > 0 {
		chunk := math.Min(n, 1)
		earned += chunk * s.Price(g)
		d := s.Demand[g]
		if d == 0 {
			d = 1
		}
		s.Demand[g] = math.Max(0.25, d*math.Pow(0.996, chunk))
		n -= chunk
		s.Inv[g] -= chunk
	}
	s.Credits += earned
	s.Earned += earned
	return earned
}

// Recover lets every market drift back toward calm. Call once per tick.
func (s *State) Recover(ticks int64) {
	rate := 0.0025 * (1 + 0.35*float64(s.Level("market.logistics")))
	for g, d := range s.Demand {
		s.Demand[g] = math.Min(1, d+(1-d)*(1-math.Pow(1-rate, float64(ticks))))
	}
}

// ---------- buying ----------

// Status says whether an upgrade can be bought and why not.
func (s *State) Status(u *Upgrade) (cost Cost, ok bool, reason string) {
	l := s.Levels[u.ID]
	if l >= u.Max {
		return Cost{}, false, i18n.T("خلصت كل المستويات", "All levels completed")
	}
	if needsBlueprint(u.ID) && !s.Blueprints[u.ID] {
		why := i18n.T("محتاج مخطط من عقد", "Needs a contract blueprint")
		if c := BlueprintSource(u.ID); c != nil {
			why = i18n.T("خلّص العقد الأول: ", "Finish the contract first: ") + c.LocalTitle()
		}
		return u.CostAt(l), false, why
	}
	for _, r := range u.Requires {
		if !s.Has(r) {
			if ru := Find(r); ru != nil {
				return u.CostAt(l), false, i18n.T("محتاج الأول: ", "Requires first: ") + i18n.T(ru.Name, ru.NameEn)
			}
		}
	}
	cost = u.CostAt(l)
	if s.Credits < cost.Credits {
		return cost, false, fmt.Sprintf(i18n.T("ناقصك %.0f كريدت", "You need %.0f more credits"), cost.Credits-s.Credits)
	}
	for g, n := range cost.Goods {
		if s.Inv[g] < n {
			return cost, false, fmt.Sprintf(i18n.T("ناقصك %.0f %s", "You need %.0f more %s"), n-s.Inv[g], CropName(g))
		}
	}
	return cost, true, ""
}

// Buy raises an upgrade one level if the colony can afford it.
func (s *State) Buy(id string) error {
	u := Find(id)
	if u == nil {
		return fmt.Errorf(i18n.T("مفيش ترقية اسمها %q", "No upgrade named %q"), id)
	}
	cost, ok, reason := s.Status(u)
	if !ok {
		return fmt.Errorf("%s", reason)
	}
	s.Credits -= cost.Credits
	for g, n := range cost.Goods {
		s.Inv[g] -= n
	}
	s.Levels[id]++
	if id == "field" {
		s.resize(FieldSize(s.Levels["field"]))
	}
	return nil
}

// ---------- lessons feed the colony ----------

type Reward struct {
	Credits float64  `json:"credits"`
	Grants  []string `json:"grants,omitempty"` // upgrade ids given for free
}

// ChapterReward is what finishing every mission of a chapter pays out:
// credits only. Commands and crops come from contracts.
func ChapterReward(n int) Reward {
	return Reward{Credits: math.Round(60 * math.Pow(float64(n), 2.2))}
}

// Claim pays a chapter's reward once.
func (s *State) Claim(n int) (Reward, error) {
	key := fmt.Sprint("chapter-", n)
	if s.Claimed[key] {
		return Reward{}, fmt.Errorf("%s", i18n.T("المكافأة دي اتصرفت قبل كده", "This reward has already been claimed"))
	}
	r := ChapterReward(n)
	s.Claimed[key] = true
	s.Credits += r.Credits
	for _, g := range r.Grants {
		if u := Find(g); u != nil && s.Levels[g] < u.Max {
			s.Levels[g]++
		}
	}
	return r, nil
}

// ---------- views for the client ----------

type ShopItem struct {
	*Upgrade
	Level  int    `json:"level"`
	Cost   Cost   `json:"cost"`
	Can    bool   `json:"can"`
	Reason string `json:"reason,omitempty"`
	// Blueprint names the contract that unlocks this upgrade while its
	// blueprint is missing.
	Blueprint string `json:"blueprint,omitempty"`
}

// Shop lists every upgrade with its current level and price.
func (s *State) Shop() []ShopItem {
	out := make([]ShopItem, 0, len(Catalog))
	for _, u := range Catalog {
		c, ok, why := s.Status(u)
		item := ShopItem{Upgrade: u, Level: s.Levels[u.ID], Cost: c, Can: ok, Reason: why}
		if needsBlueprint(u.ID) && !s.Blueprints[u.ID] {
			if src := BlueprintSource(u.ID); src != nil {
				item.Blueprint = src.ID
			}
		}
		out = append(out, item)
	}
	return out
}

// Prices lists the current price of every good.
func (s *State) Prices() map[string]float64 {
	m := map[string]float64{}
	for _, g := range Goods {
		m[g] = math.Round(s.Price(g)*100) / 100
	}
	return m
}

// TotalLevels is how many upgrade levels exist and how many are owned.
func (s *State) TotalLevels() (owned, total int) {
	for _, u := range Catalog {
		total += u.Max
		owned += min(s.Levels[u.ID], u.Max)
	}
	return owned, total
}
