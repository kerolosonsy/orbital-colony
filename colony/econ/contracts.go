package econ

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"orbital/colony/i18n"
)

// A contract is a code challenge on the colony. Finishing it unlocks
// blueprints (api.* and crop.* upgrades can only be bought with their
// blueprint) and sometimes free levels or credits.
type Contract struct {
	ID       string   `json:"id"`
	Chapter  int      `json:"chapter"` // must be complete before the contract opens
	Title    string   `json:"title"`
	TitleEn  string   `json:"titleEn"`
	Brief    string   `json:"brief"`
	BriefEn  string   `json:"briefEn"`
	Requires []string `json:"requires,omitempty"` // features.go keys; "a|b" means either
	After    []string `json:"after,omitempty"`    // contracts that must be done first
	Window   int64    `json:"window"`             // ticks from the first run; 0 = no limit
	Goal     Goal     `json:"goal"`
	Limits   Limits   `json:"limits"`
	Reward   Reward   `json:"reward"`

	Repeatable bool    `json:"repeatable,omitempty"`
	RewardRP   float64 `json:"rewardRP,omitempty"`
}

type Goal struct {
	Deliver   map[string]float64 `json:"deliver,omitempty"`   // goods sold with Export inside the window
	Any       float64            `json:"any,omitempty"`       // or this many of any goods
	Kinds     int                `json:"kinds,omitempty"`     // or this many different goods...
	Each      float64            `json:"each,omitempty"`      // ...with at least Each of every one
	Earn      float64            `json:"earn,omitempty"`      // credits earned
	Submit    string             `json:"submit,omitempty"`    // colony.Submit key (phase C5)
	MaxLost   int                `json:"maxLost"`             // plants destroyed or rotten; NoLimit = any
	MinDrones int                `json:"minDrones,omitempty"` // drones with at least MinCommands commands
	MinPrice  float64            `json:"minPrice,omitempty"`  // every unit sells for at least this share of its base price
	MinVein   int                `json:"minVein,omitempty"`   // every crystal harvested has at least this size
	MinAvg    map[string]float64 `json:"minAvg,omitempty"`    // average base harvest of a crop
	Exporters int                `json:"exporters,omitempty"` // at most this many drones may call Export
	MaxMoves  int                `json:"maxMoves,omitempty"`  // total successful moves of every drone
}

type Limits struct {
	MaxDrones int      `json:"maxDrones,omitempty"` // 0 = no limit
	Deny      []string `json:"deny,omitempty"`      // upgrades switched off while the contract runs
}

// ContractRun is the colony's record of one contract.
type ContractRun struct {
	Status   string             `json:"status"` // "", "active", "done", "failed"
	Start    int64              `json:"start"`  // tick
	Sold     map[string]float64 `json:"sold,omitempty"`
	Earned   float64            `json:"earned,omitempty"`
	Lost     int                `json:"lost,omitempty"`
	Best     int64              `json:"best,omitempty"` // fewest ticks it was finished in
	Tries    int                `json:"tries,omitempty"`
	Cmds     map[int]int        `json:"cmds,omitempty"`     // commands per drone
	Harvests map[string]int     `json:"harvests,omitempty"` // harvests per crop
	Base     map[string]float64 `json:"base,omitempty"`     // base yield per crop, before upgrades
	Sellers  map[int]bool       `json:"sellers,omitempty"`  // drones that exported
	Moves    int                `json:"moves,omitempty"`
	Broke    string             `json:"broke,omitempty"` // a rule the run broke: price, vein, exporter

	Submission string `json:"submission,omitempty"`

	Stopping  bool    `json:"stopping,omitempty"`
	Stopped   bool    `json:"stopped,omitempty"`
	CargoLost float64 `json:"cargoLost,omitempty"`

	ResearchDay  int64 `json:"researchDay,omitempty"`
	RepeatsToday int   `json:"repeatsToday,omitempty"`
}

// NoLimit marks a goal without a loss limit.
const NoLimit = -1

// MinCommands is how many commands a drone needs to count as working.
const MinCommands = 20

func deliver(kv ...any) map[string]float64 {
	m := map[string]float64{}
	for i := 0; i < len(kv); i += 2 {
		m[kv[i].(string)] = float64(kv[i+1].(int))
	}
	return m
}

// Contracts lists every contract in chapter order.
var Contracts = []*Contract{
	{ID: "c02-price-watch", Chapter: 2, Requires: []string{"if"}, Window: 1500,
		Title: "مراقبة السعر", TitleEn: "Price Watch",
		Brief:   "اكسب 50 كريدت، وكل وحدة تتباع بـ 90% من سعرها الأساسي على الأقل. السعر بيقل مع كل وحدة بتبيعها وبيرجع مع الوقت: اقرا colony.Price واستنى. شاشة البورصة متاحة طول العقد. الدرس: فصل 2.",
		BriefEn: "Earn 50 credits, and sell every unit for at least 90% of its base price. Each unit you sell lowers the price, and it recovers over time: read colony.Price and wait. The market display works during the contract. Lesson: chapter 2.",
		Goal:    Goal{Earn: 50, MinPrice: 0.9, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"api.prices"}}},
	{ID: "c03-full-sweep", Chapter: 3, Requires: []string{"for2"}, Window: 1750,
		Title: "المسح الكامل", TitleEn: "Full Sweep",
		Brief:   "سلّم 120 طحالب في 1750 tick على أرض 5×5. لف على الأرض كلها بحلقتين متداخلتين ومتضيّعش أوامر. الدرس: فصل 3.",
		BriefEn: "Deliver 120 moss within 1750 ticks on a 5×5 field. Sweep the whole field with two nested loops and waste no commands. Lesson: chapter 3.",
		Goal:    Goal{Deliver: deliver("moss", 120), MaxLost: NoLimit}, Reward: Reward{Grants: []string{"crop.lichen"}}},
	{ID: "c04-right-ground", Chapter: 4, Requires: []string{"switch"}, Window: 1500,
		Title: "الأرض الصح", TitleEn: "The Right Ground",
		Brief:   "سلّم 30 طحالب و20 أشنة من غير ما تدمّر ولا نبات. الطحالب على التربة والأشنة على الجليد: اختار بـ switch. السكانر متاح طول العقد. الدرس: فصل 4.",
		BriefEn: "Deliver 30 moss and 20 lichen without destroying a single plant. Moss goes on soil and lichen on ice: choose with switch. The scanner works during the contract. Lesson: chapter 4.",
		Goal:    Goal{Deliver: deliver("moss", 30, "lichen", 20), MaxLost: 0}, Reward: Reward{Grants: []string{"api.scanner"}}},
	{ID: "c05-careful-harvest", Chapter: 5, Requires: []string{"errcheck", "func"}, Window: 1200, After: []string{"c04-right-ground"},
		Title: "الحصاد الحذر", TitleEn: "Careful Harvest",
		Brief:   "سلّم 80 من أي محصول من غير ما تدمّر ولا نبات. السكانر مقفول في العقد ده: استخدم TryHarvest وعالج الـ error اللي بترجّعه. الدرس: فصل 5.",
		BriefEn: "Deliver 80 of any crop without destroying a single plant. The scanner is off for this contract: use TryHarvest and handle the error it returns. Lesson: chapter 5.",
		Goal:    Goal{Any: 80, MaxLost: 0}, Limits: Limits{Deny: []string{"api.scanner"}}, Reward: Reward{Grants: []string{"api.try"}}},
	{ID: "c06-ripe-ledger", Chapter: 6, Requires: []string{"map", "slice"}, Window: 1000,
		Title: "دفتر النضج", TitleEn: "Ripe Ledger",
		Brief:   "سلّم 120 طحالب في 1000 tick. سجّل إمتى كل خانة هتنضج (map أو slice) وروح للناضج بس، بدل ما تلف على الكل. الدرس: فصل 6.",
		BriefEn: "Deliver 120 moss within 1000 ticks. Record when each tile ripens (a map or a slice) and visit only ripe ones instead of sweeping everything. Lesson: chapter 6.",
		Goal:    Goal{Deliver: deliver("moss", 120), MaxLost: NoLimit}, Reward: Reward{Grants: []string{"crop.crystal", "api.measure"}}},
	{ID: "c07-plot-type", Chapter: 7, Requires: []string{"struct", "method"}, Window: 800,
		Title: "نوع للقطعة", TitleEn: "A Plot Type",
		Brief:   "سلّم 40 فطر مضيء، ومتوسط الحصاد لكل فطرة 3 على الأقل: ازرعهم عناقيد واحصدهم وهما ناضجين جنب بعض. الفطر المضيء متاح طول العقد. الدرس: فصل 7.",
		BriefEn: "Deliver 40 glowcap with an average harvest of at least 3 per cap: plant them in clusters and harvest them while their neighbors are ripe. Glowcap works during the contract. Lesson: chapter 7.",
		Goal:    Goal{Deliver: deliver("glowcap", 40), MinAvg: map[string]float64{"glowcap": 3}, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"crop.glowcap"}}},
	{ID: "c08-big-veins", Chapter: 8, Requires: []string{"pointer", "nilcheck"}, Window: 2000,
		Title: "العروق الكبيرة", TitleEn: "Big Veins",
		Brief:   "سلّم 60 كريستال من عروق حجمها 4 أو أكتر بس. أي عرق أصغر بيلغي العقد. الرادار متاح طول العقد. الدرس: فصل 8.",
		BriefEn: "Deliver 60 crystal from veins of size 4 or more only. Any smaller vein voids the contract. The radar works during the contract. Lesson: chapter 8.",
		Goal:    Goal{Deliver: deliver("crystal", 60), MinVein: 4, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"api.radar"}}},
	{ID: "c09-crop-protocol", Chapter: 9, Requires: []string{"iface"}, Window: 2000,
		Title: "بروتوكول المحاصيل", TitleEn: "Crop Protocol",
		Brief:   "سلّم 3 محاصيل مختلفة، 40 من كل واحد، في نفس العقد. خلّي كل محصول نوع بيطبّق interface واحدة. الدرس: فصل 9.",
		BriefEn: "Deliver 3 different crops, 40 of each, in one contract. Make every crop a type that implements one interface. Lesson: chapter 9.",
		Goal:    Goal{Kinds: 3, Each: 40, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"crop.reed"}}},
	{ID: "c10-fleet", Chapter: 10, Requires: []string{"go", "waitgroup"}, Window: 800,
		Title: "الأسطول", TitleEn: "The Fleet",
		Brief:   "سلّم 200 طحالب في 800 tick، و3 درونات على الأقل شغالين فعلاً. colony.Spawn متاحة طول العقد. الدرس: فصل 10.",
		BriefEn: "Deliver 200 moss within 800 ticks with at least 3 drones actually working. colony.Spawn works during the contract. Lesson: chapter 10.",
		Goal:    Goal{Deliver: deliver("moss", 200), MinDrones: 3, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"api.spawn", "hangar"}}},
	{ID: "c11-dispatch", Chapter: 11, Requires: []string{"chan", "select"}, Window: 1500,
		Title: "مركز التوزيع", TitleEn: "Dispatch",
		Brief:   "سلّم 100 كريستال بـ 3 درونات على الأقل، ودرون واحد بس هو اللي يبيع (Export). الباقيين يبعتوله بالـ channels. الدرس: فصل 11.",
		BriefEn: "Deliver 100 crystal with at least 3 drones, and only one drone may sell (Export). The others send to it over channels. Lesson: chapter 11.",
		Goal:    Goal{Deliver: deliver("crystal", 100), MinDrones: 3, Exporters: 1, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"hangar"}}},
	{ID: "c12-one-ledger", Chapter: 12, Requires: []string{"mutex"}, Window: 600,
		Title: "دفتر واحد", TitleEn: "One Ledger",
		Brief:   "سلّم 300 طحالب بـ 4 درونات شغالين، وكلهم بيكتبوا في دفتر مشترك محمي بـ sync.Mutex. لازم تشغّل فحص السباقات ويعدّي من غير data race. الدرس: فصل 12.",
		BriefEn: "Deliver 300 moss with 4 working drones, all writing to one shared ledger guarded by a sync.Mutex. Enable race mode and finish without a data race. Lesson: chapter 12.",
		Goal:    Goal{Deliver: deliver("moss", 300), MinDrones: 4, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"market.logistics"}}},
	{ID: "c13-reed-clock", Chapter: 13, Requires: []string{"sel:context.WithTimeout"}, Window: 600,
		Title: "ساعة البوص", TitleEn: "Reed Clock",
		Brief:   "سلّم 50 بوص، ومسموحلك تخسر نباتين بس. البوص بيعفّن لو اتأخرت عليه: حط لكل انتظار مهلة بـ context.WithTimeout. الدرس: فصل 13.",
		BriefEn: "Deliver 50 pulse reed, losing at most two plants. Reed rots if you are late: give every wait a deadline with context.WithTimeout. Lesson: chapter 13.",
		Goal:    Goal{Deliver: deliver("reed", 50), MaxLost: 2}, Reward: Reward{Grants: []string{"market.contracts"}}},
	{ID: "c14-generic-queue", Chapter: 14, Requires: []string{"generictype"}, Window: 1500,
		Title: "طابور عام", TitleEn: "Generic Queue",
		Brief:   "سلّم 150 من محصولين مختلفين (75 من كل واحد على الأقل)، وجدول الشغل نوع generic إنت عرّفته. الدرس: فصل 14.",
		BriefEn: "Deliver 150 of two different crops (at least 75 of each), with the work queue a generic type you define. Lesson: chapter 14.",
		Goal:    Goal{Kinds: 2, Each: 75, MaxLost: NoLimit}, Reward: Reward{Credits: 500}},
	{ID: "c15-shortest-run", Chapter: 15, Requires: []string{"recursion|sel:heap.Push"}, Window: 2000,
		Title: "أقصر طريق", TitleEn: "Shortest Run",
		Brief:   "سلّم 80 كريستال وكل الدرونات تتحرك 400 خطوة بس على بعض. رتّب العروق بأقصر مسار (recursion أو heap). الدرس: فصل 15.",
		BriefEn: "Deliver 80 crystal with all drones moving no more than 400 steps in total. Order the veins by the shortest route (recursion or a heap). Lesson: chapter 15.",
		Goal:    Goal{Deliver: deliver("crystal", 80), MaxMoves: 400, MaxLost: NoLimit}, Reward: Reward{Grants: []string{"api.swap"}}},
	testedPlannerContract,
	{ID: "c17-deep-haul", Chapter: 17, Requires: []string{"for2", "slice"}, Window: 200,
		Title: "الحمولة العميقة", TitleEn: "Deep Haul",
		Brief:   "سلّم 20 جيود من الطبقة 2. أوامر الحفر متاحة طول العقد، والحفار لازم يوصل المستوى 2. الدرس: فصل 17.",
		BriefEn: "Deliver 20 geodes from layer 2. The digging commands work during the contract; the drill must reach level 2. Lesson: chapter 17.",
		Goal:    Goal{Deliver: deliver("geode", 20), MaxLost: NoLimit}, Reward: Reward{Grants: []string{"api.dig", "drill"}}},
}

// FindContract returns the contract with that id, or nil.
func FindContract(id string) *Contract {
	for _, c := range Contracts {
		if c.ID == id {
			return c
		}
	}
	return nil
}

// LocalTitle is the localized title of a contract.
func (c *Contract) LocalTitle() string { return i18n.T(c.Title, c.TitleEn) }

// needsBlueprint reports whether an upgrade can only be bought after a contract.
func needsBlueprint(id string) bool {
	return strings.HasPrefix(id, "api.") || strings.HasPrefix(id, "crop.")
}

// BlueprintSource is the contract whose reward is that upgrade's blueprint.
func BlueprintSource(id string) *Contract {
	for _, c := range Contracts {
		if slices.Contains(c.Reward.Grants, id) {
			return c
		}
	}
	return nil
}

func (s *State) run(id string) *ContractRun {
	r := s.Contracts[id]
	if r == nil {
		r = &ContractRun{}
		s.Contracts[id] = r
	}
	return r
}

// Won reports whether a contract was finished at least once.
func (s *State) Won(id string) bool { r := s.Contracts[id]; return r != nil && r.Best > 0 }

// Available lists the contracts the colony can take now: their chapter is
// in done, the contracts before them are won, and none is active.
func (s *State) Available(done map[int]bool) []*Contract {
	var out []*Contract
	if s.Active != "" {
		return out
	}
	for _, c := range Contracts {
		if !done[c.Chapter] {
			continue
		}
		ok := true
		for _, a := range c.After {
			ok = ok && s.Won(a)
		}
		if ok {
			out = append(out, c)
		}
	}
	return out
}

// Accept makes id the active contract and starts its clock. An empty id
// cancels the active contract; nothing is lost but the attempt.
func (s *State) Accept(id string) error {
	if id == "" {
		if s.Active == "" {
			return fmt.Errorf("%s", i18n.T("مفيش عقد نشط", "No active contract"))
		}
		s.end("")
		return nil
	}
	c := FindContract(id)
	if c == nil {
		return fmt.Errorf(i18n.T("مفيش عقد اسمه %q", "No contract named %q"), id)
	}
	if s.Active != "" {
		return fmt.Errorf("%s", i18n.T("عندك عقد نشط بالفعل. خلّصه أو ألغيه الأول", "You already have an active contract. Finish or cancel it first"))
	}
	for _, a := range c.After {
		if !s.Won(a) {
			return fmt.Errorf("%s", i18n.T("خلّص العقد الأول: ", "Finish the contract first: ")+FindContract(a).LocalTitle())
		}
	}
	old := s.run(id)
	s.Contracts[id] = &ContractRun{Status: "active", Start: s.Tick, Best: old.Best, Tries: old.Tries + 1, ResearchDay: old.ResearchDay, RepeatsToday: old.RepeatsToday}
	s.Active = id
	return nil
}

// end closes the active contract with a status ("" for a cancel).
func (s *State) end(status string) {
	r := s.run(s.Active)
	if status != "done" && r.Best > 0 {
		status = "done" // a failed replay keeps the earlier win
	}
	r.Status = status
	s.Active = ""
}

// Fail ends the active contract as failed.
func (s *State) Fail() {
	if s.Active != "" {
		s.end("failed")
	}
}

// ActiveContract is the active contract and its record, or nil.
func (s *State) ActiveContract() (*Contract, *ContractRun) {
	c := FindContract(s.Active)
	if c == nil {
		return nil, nil
	}
	return c, s.run(c.ID)
}

// ---------- what the simulation records ----------

// RecordSale counts n units of g sold by drone for earned credits; low is the
// cheapest unit price of the sale.
func (s *State) RecordSale(drone int, g string, n, earned, low float64) {
	c, r := s.ActiveContract()
	if c == nil {
		return
	}
	if r.Sold == nil {
		r.Sold = map[string]float64{}
	}
	if r.Sellers == nil {
		r.Sellers = map[int]bool{}
	}
	r.Sold[g] += n
	r.Earned += earned
	r.Sellers[drone] = true
	switch {
	case c.Goal.MinPrice > 0 && low < c.Goal.MinPrice*BasePrice[g]-1e-9:
		r.Broke = "price"
	case c.Goal.Exporters > 0 && len(r.Sellers) > c.Goal.Exporters:
		r.Broke = "exporter"
	}
}

// RecordHarvest counts one harvest of g whose yield before upgrades was base.
// size is the crystal's size (0 for crops).
func (s *State) RecordHarvest(g string, base float64, size int) {
	c, r := s.ActiveContract()
	if c == nil {
		return
	}
	if r.Harvests == nil {
		r.Harvests, r.Base = map[string]int{}, map[string]float64{}
	}
	r.Harvests[g]++
	r.Base[g] += base
	if g == "crystal" && c.Goal.MinVein > 0 && size < c.Goal.MinVein {
		r.Broke = "vein"
	}
}

// RecordLost counts a plant destroyed before it ripened or left to rot.
func (s *State) RecordLost() {
	if c, r := s.ActiveContract(); c != nil {
		r.Lost++
	}
}

// RecordCommand counts one drone command; move tells whether it moved the drone.
func (s *State) RecordCommand(drone int, move bool) {
	c, r := s.ActiveContract()
	if c == nil {
		return
	}
	if r.Cmds == nil {
		r.Cmds = map[int]int{}
	}
	r.Cmds[drone]++
	if move {
		r.Moves++
	}
}

// ---------- progress ----------

func (r *ContractRun) working() int {
	n := 0
	for _, c := range r.Cmds {
		if c >= MinCommands {
			n++
		}
	}
	return n
}

func fmtNum(v float64) string { return fmt.Sprint(math.Floor(v)) }

// ProgressPair is Progress with the reason in both languages: Arabic, English.
func (s *State) ProgressPair(tick int64) (done, failed bool, why [2]string) {
	c, r := s.ActiveContract()
	if c == nil {
		return false, false, why
	}
	g := c.Goal
	if g.Submit != "" {
		return s.submissionProgress(c, r, tick)
	}
	if c.ID == CleanStopContract {
		return r.Stopped && r.CargoLost == 0, r.Stopped && r.CargoLost > 0, [2]string{"مستني إغلاق من غير خسارة حمولة", "Waiting for shutdown without cargo loss"}
	}
	switch {
	case r.Broke == "price":
		return false, true, [2]string{fmt.Sprintf("بعت تحت %.0f%% من السعر الأساسي", g.MinPrice*100), fmt.Sprintf("You sold below %.0f%% of the base price", g.MinPrice*100)}
	case r.Broke == "vein":
		return false, true, [2]string{fmt.Sprintf("حصدت عرق أصغر من %d", g.MinVein), fmt.Sprintf("You harvested a vein smaller than %d", g.MinVein)}
	case r.Broke == "exporter":
		return false, true, [2]string{fmt.Sprintf("أكتر من %d درون باع", g.Exporters), fmt.Sprintf("More than %d drone sold goods", g.Exporters)}
	case g.MaxLost >= 0 && r.Lost > g.MaxLost:
		return false, true, [2]string{fmt.Sprintf("خسرت %d نبات (المسموح %d)", r.Lost, g.MaxLost), fmt.Sprintf("You lost %d plants (allowed: %d)", r.Lost, g.MaxLost)}
	case g.MaxMoves > 0 && r.Moves > g.MaxMoves:
		return false, true, [2]string{fmt.Sprintf("الدرونات اتحركت %d خطوة (المسموح %d)", r.Moves, g.MaxMoves), fmt.Sprintf("The drones moved %d steps (allowed: %d)", r.Moves, g.MaxMoves)}
	}
	var ar, en []string
	done = true
	need := func(ok bool, a, e string) {
		done = done && ok
		ar, en = append(ar, a), append(en, e)
	}
	for _, k := range sortedKeys(g.Deliver) {
		need(r.Sold[k] >= g.Deliver[k], fmt.Sprintf("%s %s/%s", cropNames[k], fmtNum(r.Sold[k]), fmtNum(g.Deliver[k])), fmt.Sprintf("%s %s/%s", cropNamesEn[k], fmtNum(r.Sold[k]), fmtNum(g.Deliver[k])))
	}
	if g.Any > 0 {
		total := 0.0
		for _, v := range r.Sold {
			total += v
		}
		need(total >= g.Any, fmt.Sprintf("محاصيل %s/%s", fmtNum(total), fmtNum(g.Any)), fmt.Sprintf("goods %s/%s", fmtNum(total), fmtNum(g.Any)))
	}
	if g.Kinds > 0 {
		n := 0
		for _, v := range r.Sold {
			if v >= g.Each {
				n++
			}
		}
		need(n >= g.Kinds, fmt.Sprintf("أنواع فيها %s: %d/%d", fmtNum(g.Each), n, g.Kinds), fmt.Sprintf("kinds with %s: %d/%d", fmtNum(g.Each), n, g.Kinds))
	}
	if g.Earn > 0 {
		need(r.Earned >= g.Earn, fmt.Sprintf("كريدت %s/%s", fmtNum(r.Earned), fmtNum(g.Earn)), fmt.Sprintf("credits %s/%s", fmtNum(r.Earned), fmtNum(g.Earn)))
	}
	for _, k := range sortedKeys(g.MinAvg) {
		avg := 0.0
		if r.Harvests[k] > 0 {
			avg = r.Base[k] / float64(r.Harvests[k])
		}
		need(r.Harvests[k] > 0 && avg >= g.MinAvg[k], fmt.Sprintf("متوسط %s %.1f/%.0f", cropNames[k], avg, g.MinAvg[k]), fmt.Sprintf("%s average %.1f/%.0f", cropNamesEn[k], avg, g.MinAvg[k]))
	}
	if g.MinDrones > 0 {
		n := r.working()
		need(n >= g.MinDrones, fmt.Sprintf("درونات شغالة %d/%d", n, g.MinDrones), fmt.Sprintf("working drones %d/%d", n, g.MinDrones))
	}
	if g.Submit != "" {
		need(false, "مستني "+g.Submit, "waiting for "+g.Submit)
	}
	why = [2]string{strings.Join(ar, "، "), strings.Join(en, ", ")}
	if !done && c.Window > 0 && tick-r.Start > c.Window {
		return false, true, [2]string{fmt.Sprintf("الوقت خلص (%d tick): %s", c.Window, why[0]), fmt.Sprintf("Time ran out (%d ticks): %s", c.Window, why[1])}
	}
	return done, false, why
}

// Progress checks the active contract at tick: done, failed (and why), or the
// progress so far in why.
func (s *State) Progress(tick int64) (done, failed bool, why string) {
	done, failed, pair := s.ProgressPair(tick)
	return done, failed, i18n.T(pair[0], pair[1])
}

// Fraction is how far the active contract is, from 0 to 1, by its slowest goal.
func (s *State) Fraction() float64 {
	c, r := s.ActiveContract()
	if c == nil {
		return 0
	}
	g := c.Goal
	if g.Submit != "" && r.Submission != "correct" {
		return 0
	}
	if c.ID == CleanStopContract && !r.Stopped {
		return 0
	}
	f := 1.0
	part := func(have, want float64) {
		if want > 0 {
			f = math.Min(f, math.Min(1, have/want))
		}
	}
	for k, v := range g.Deliver {
		part(r.Sold[k], v)
	}
	total := 0.0
	var each []float64
	for _, v := range r.Sold {
		total += v
		each = append(each, math.Min(v, g.Each))
	}
	part(total, g.Any)
	if g.Kinds > 0 {
		slices.Sort(each)
		slices.Reverse(each)
		sum := 0.0
		for i := 0; i < len(each) && i < g.Kinds; i++ {
			sum += each[i]
		}
		part(sum, float64(g.Kinds)*g.Each)
	}
	part(r.Earned, g.Earn)
	part(float64(r.working()), float64(g.MinDrones))
	return f
}

// Left is how many ticks the active contract has before its window closes
// (-1 when it has no window).
func (s *State) Left(tick int64) int64 {
	c, r := s.ActiveContract()
	if c == nil || c.Window == 0 {
		return -1
	}
	return max(0, c.Window-(tick-r.Start))
}

// Settle pays the active contract when its goal is met at tick. The reward
// is paid the first time only; replays only improve the best time.
func (s *State) Settle(tick int64) (Reward, error) {
	c, r := s.ActiveContract()
	if c == nil {
		return Reward{}, fmt.Errorf("%s", i18n.T("مفيش عقد نشط", "No active contract"))
	}
	if done, _, _ := s.Progress(tick); !done {
		return Reward{}, fmt.Errorf("%s", i18n.T("العقد لسه مخلصش", "The contract is not finished yet"))
	}
	first := r.Best == 0
	took := max(1, tick-r.Start)
	if first || took < r.Best {
		r.Best = took
	}
	s.end("done")
	if c.Repeatable {
		return s.settleExperiment(c, r, tick), nil
	}
	if !first {
		return Reward{}, nil
	}
	s.Credits += c.Reward.Credits
	for _, id := range c.Reward.Grants {
		if needsBlueprint(id) {
			s.Blueprints[id] = true
		} else if u := Find(id); u != nil && (u.Max == 0 || s.Levels[id] < u.Max) {
			s.Levels[id]++
			if id == "field" {
				s.resize(FieldSize(s.Levels[id]))
			}
		}
	}
	return c.Reward, nil
}

// ---------- what a running contract lends and denies ----------

func (s *State) denied(id string) bool {
	c, _ := s.ActiveContract()
	return c != nil && slices.Contains(c.Limits.Deny, id)
}

// Lent reports whether the active contract lends id: a contract's own
// reward works while it runs, so the challenge can use what it unlocks.
func (s *State) Lent(id string) bool {
	c, _ := s.ActiveContract()
	return c != nil && slices.Contains(c.Reward.Grants, id)
}

// Usable reports whether a command or crop works right now: bought (or lent
// by the active contract) and not denied by it.
func (s *State) Usable(id string) bool {
	return (s.Has(id) || s.Lent(id)) && !s.denied(id)
}

// Denied reports whether the active contract switches id off.
func (s *State) Denied(id string) bool { return s.denied(id) }

// DroneCap is how many drones may fly now, including the active contract's
// loan (enough drones for its goal) and limit.
func (s *State) DroneCap() int {
	n := s.Drones()
	c, _ := s.ActiveContract()
	if c == nil {
		return n
	}
	if s.Lent("api.spawn") || s.Lent("hangar") {
		n = max(n, c.Goal.MinDrones)
	}
	if c.Limits.MaxDrones > 0 {
		n = min(n, c.Limits.MaxDrones)
	}
	return n
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// migrateV2 turns licences owned before contracts existed into blueprints.
func (s *State) migrateV2() {
	for id, l := range s.Levels {
		if l > 0 && needsBlueprint(id) {
			s.Blueprints[id] = true
		}
	}
	s.Version = 2
}
