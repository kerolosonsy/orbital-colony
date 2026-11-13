package econ

// SystemCounters survive restarts so counter metrics never reset with a program.
type SystemCounters struct {
	SoldByGood map[string]float64 `json:"soldByGood,omitempty"`
	Sold       float64            `json:"sold"`
	Commands   uint64             `json:"commands"`
}

func addControlNodes() {
	for _, n := range []Node{
		{ID: "sys.control", Branch: "الخدمات", Name: "مركز القيادة", NameEn: "Control API", Desc: "شغّل سيرفر HTTP جنب الدرونات.", DescEn: "Run an HTTP server alongside your drones.", Chapter: 26, After: []Prereq{{"api.prices", 1}}},
		{ID: "sys.warehouse", Branch: "الخدمات", Name: "المخزن", NameEn: "SQLite Warehouse", Desc: "احفظ مخزونك في SQLite بين مرات التشغيل.", DescEn: "Persist inventory in SQLite across runs.", Chapter: 29, After: []Prereq{{"sys.control", 1}}},
		{ID: "sys.metrics", Branch: "الأنظمة", Name: "المراقبة", NameEn: "Metrics", Desc: "اعرض مقاييس Prometheus واستعداد الخدمة.", DescEn: "Expose Prometheus metrics and service readiness.", Chapter: 36, After: []Prereq{{"sys.control", 1}}},
	} {
		u := &Upgrade{ID: n.ID, Group: n.Branch, GroupEn: "Systems", Name: n.Name, NameEn: n.NameEn, Desc: n.Desc, DescEn: n.DescEn, Max: 1, Lesson: n.Chapter, Requires: []string{n.After[0].ID}, cost: credits(1000, 1)}
		add(u)
		n.Cost = CostRule{Credits: 1000, Growth: 1}
		Tree = append(Tree, n)
	}
}

var controlAPIContract = &Contract{ID: "h26-control-api", Chapter: 26, Window: 100000, Requires: []string{"sel:http.ListenAndServe|sel:http.Server|sel:http.Serve"},
	Title: "مركز القيادة", TitleEn: "Control API",
	Brief:   "امتلك sys.control وشغّل HTTP على ORBITAL_LISTEN_ADDR. جاوب 12 فحص في 200ms: GET /healthz = 200، GET /inventory = مخزون حقيقي خلال آخر ثانية، POST /orders {good,amount} = 202 عند القبول مع Export كامل خلال ثانية أو 409 لو المخزون مش كفاية.",
	BriefEn: "Own sys.control and serve HTTP on ORBITAL_LISTEN_ADDR. Answer 12 probes within 200ms each: GET /healthz returns 200; GET /inventory returns real inventory from the last second; POST /orders {good,amount} returns 202 on acceptance and exports the full amount within one second, or 409 for insufficient stock.",
	Goal:    Goal{Submit: "system:h26-control-api", MaxLost: NoLimit}}

var drainContract = &Contract{ID: "h28-drain", Chapter: 28, Window: 100000, Requires: []string{"sel:http.Server", "sel:colony.Shutdown"},
	Title: "إغلاق هادي", TitleEn: "Drain In Flight",
	Brief:   "امتلك sys.control. طلب orders عليه X-Orbital-Delay-Ms: 500 لازم يفضل جاري 500ms ويكمّل وقت SIGTERM. اقفل الاستماع للاتصالات الجديدة خلال 3 ثواني واستخدم colony.Shutdown.",
	BriefEn: "Own sys.control. An orders request with X-Orbital-Delay-Ms: 500 must remain in flight for 500ms and finish across SIGTERM. Refuse new connections within three seconds; register colony.Shutdown.",
	Goal:    Goal{Submit: "system:h28-drain", MaxLost: NoLimit}}

var warehouseContract = &Contract{ID: "s29-warehouse", Chapter: 29, Window: 100000, Requires: []string{"sel:sql.Open"},
	Title: "مخزن دائم", TitleEn: "Persistent Warehouse",
	Brief:   "امتلك sys.warehouse وافتح SQLite في ORBITAL_DB_PATH. POST /inventory {id,goods} يسجل دفعة مرة واحدة ويرد 204؛ GET /inventory يرجع الإجمالي من القاعدة. هنكرر الدفعات ونوقف ونشغّل البرنامج 3 مرات. بيانات المخزن مستقلة عن مخزون الدرونات وقت العقد.",
	BriefEn: "Own sys.warehouse and open SQLite at ORBITAL_DB_PATH. POST /inventory {id,goods} records a batch once and returns 204; GET /inventory returns database totals. Batches are retried and the program is stopped and restarted three times. Warehouse records are independent of drone inventory during this contract.",
	Goal:    Goal{Submit: "system:s29-warehouse", MaxLost: NoLimit}}

var transferContract = &Contract{ID: "s29-transfer", Chapter: 29, Window: 100000, Requires: []string{"sel:sql.Open"},
	Title: "تحويل ذري", TitleEn: "Atomic Transfer",
	Brief:   "امتلك sys.warehouse. ابدأ مخزنين a وb كل واحد 100. GET /stores يرجع {a:amount,b:amount}. POST /transfer {id,from,to,amount} يحوّل في transaction، 200 للنجاح أو إعادة نفس id، و409 لو الرصيد مش كفاية. المجموع ثابت ومفيش رصيد سالب تحت طلبات متوازية.",
	BriefEn: "Own sys.warehouse. Initialize stores a and b with 100 each. GET /stores returns {a:amount,b:amount}. POST /transfer {id,from,to,amount} transfers in a transaction: 200 for success or replay of an id, 409 for insufficient funds. Concurrent requests must conserve totals and never create negative balances.",
	Goal:    Goal{Submit: "system:s29-transfer", MaxLost: NoLimit}}

var metricsContract = &Contract{ID: "o36-dashboard", Chapter: 36, Window: 100000, Requires: []string{"sel:http.ListenAndServe|sel:http.Server|sel:http.Serve"},
	Title: "لوحة المراقبة", TitleEn: "Colony Dashboard",
	Brief:   "امتلك sys.metrics. اعرض colony_credits وcolony_goods_sold_total وcolony_drone_commands_total وcolony_buildings بقيم colony.SystemStats (±2%). 3 قراءات كل 5 ثواني. حد 16 قيمة labels مختلفة لكل مقياس طول التشغيل. /readyz = 200 قبل orders و503 وقت الإغلاق قبل قفل الاستماع؛ خلّي نافذة 150ms قبل Shutdown.",
	BriefEn: "Own sys.metrics. Expose colony_credits, colony_goods_sold_total, colony_drone_commands_total and colony_buildings matching colony.SystemStats (±2%). Pass three scrapes five seconds apart. At most 16 distinct label values per metric across the run. /readyz returns 200 before orders and 503 during shutdown before closing the listener; leave a 150ms readiness window before Shutdown.",
	Goal:    Goal{Submit: "system:o36-dashboard", MaxLost: NoLimit}}

func IsControlContract(id string) bool {
	switch id {
	case "h26-control-api", "h28-drain", "s29-warehouse", "s29-transfer", "o36-dashboard":
		return true
	}
	return false
}

func (s *State) VerifyControl(id string, ok bool) {
	if s.Active != id || !IsControlContract(id) {
		return
	}
	_, run := s.ActiveContract()
	run.Broke = "system.failed"
	if ok {
		run.Broke = "system.ok"
	}
}

func (s *State) controlProgress() (bool, bool, [2]string) {
	_, run := s.ActiveContract()
	return run.Broke == "system.ok", run.Broke == "system.failed", [2]string{"مستني فحص النظام", "Waiting for system verification"}
}

func (s *State) recordSystemSale(good string, n float64) {
	s.Systems.Sold += n
	if s.Systems.SoldByGood == nil {
		s.Systems.SoldByGood = map[string]float64{}
	}
	s.Systems.SoldByGood[good] += n
}

func (s *State) controlFraction() float64 {
	done, _, _ := s.controlProgress()
	if done {
		return 1
	}
	return 0
}
