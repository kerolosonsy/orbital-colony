package levels

func init() {
	add(&Level{ID: "36-1-metrics", Chapter: 36, Kind: "lab", Title: "العدادات", Goal: "اعرض مقاييس رصيف الهبوط بـ Prometheus بـ labels محدودة وendpoint وتحديث متوازي",
		Requires: []string{"sel:prometheus.NewCounterVec", "sel:promhttp.HandlerFor"}, Bridge: []int{96}, Setup: seededLabRoom, Race: true,
		Deps:  []string{"github.com/prometheus/client_golang"},
		Files: []LabFile{{Name: "main.go", Starter: ch361Starter}},
		Lesson: `# عدادات رصيف الهبوط

## قبل ما تكتب
- لو كل اسم سفينة بقى label، كام سلسلة هتتخزن بعد شهر؟
- مين اللي بيقرا المقاييس، وبيقراها إمتى؟

## ليه مقاييس؟
السجل بيحكي حكاية طلب واحد. المقياس رقم بيتجمّع على مدار الوقت: كام سفينة نزلت، وكام واحدة مستنية دلوقتي، والهبوط بياخد قد إيه. Prometheus بيزور السيرفر كل شوية ويقرا الأرقام كلها من endpoint واحد.

استورد ~github.com/prometheus/client_golang/prometheus~. كل مقياس اسمه ~Collector~، ولازم يتسجّل في ~Registry~ عشان يظهر.

## السجل الخاص
~prometheus.NewRegistry()~ بيعمل سجل فاضي بتاعك. فيه سجل عام اسمه ~prometheus.DefaultRegisterer~، لكن أي حاجة تتسجّل فيه بتفضل موجودة للبرنامج كله، والتست التاني بيلاقيها. هنا السجل بييجي من برا كـ parameter.
~~~go
reg := prometheus.NewRegistry()
err := reg.Register(c) // نفس الاسم مرتين: خطأ
~~~
~reg.Register~ بيرجع خطأ لو اسم المقياس متسجّل قبل كده في نفس السجل. ~MustRegister~ بيعمل panic بدل الخطأ؛ هنا رجّع الخطأ.

## العداد (counter)
العداد بيزيد بس، وعمره ما بينقص. ~prometheus.NewCounterVec~ بياخد اسم ونص مساعدة وأسماء الـ labels. ~WithLabelValues~ بيرجع السلسلة اللي ليها القيم دي بالظبط، بنفس ترتيب الأسماء.
~~~go
pings := prometheus.NewCounterVec(prometheus.CounterOpts{
 Name: "relay_pings_total", Help: "Pings by relay.",
}, []string{"relay"})
pings.WithLabelValues("r7").Inc()   // 1
pings.WithLabelValues("r7").Add(2)  // 3
~~~
اسم العداد بينتهي بـ ~_total~ عادةً. نداء ~WithLabelValues~ من غير ~Inc~ بيعمل السلسلة بقيمة 0، فتظهر في الـ endpoint من الأول.

## المؤشر (gauge)
المؤشر رقم بيطلع وينزل: طول طابور، أو عدد شغّال دلوقتي. ~prometheus.NewGauge~ مؤشر واحد من غير labels، و ~prometheus.NewGaugeVec~ بـ labels.
~~~go
fuel := prometheus.NewGauge(prometheus.GaugeOpts{Name: "tank_fuel_liters", Help: "Fuel."})
fuel.Set(42)
fuel.Inc()
fuel.Dec()
~~~

## المدرّج (histogram)
المدرّج بيعدّ القيم في صناديق (buckets). ~Observe(v)~ بيزوّد كل صندوق حده الأعلى ~le~ أكبر من أو يساوي v، وكمان ~_sum~ و ~_count~. الحدود في ~Buckets~ لازم تكون تصاعدية، وصندوق ~+Inf~ بيتضاف لوحده.
~~~go
lat := prometheus.NewHistogramVec(prometheus.HistogramOpts{
 Name: "relay_reply_seconds", Help: "Reply time.", Buckets: []float64{0.5, 2},
}, []string{"relay"})
lat.WithLabelValues("r7").Observe(1.2) // في le="2" و le="+Inf" بس
~~~
~d.Seconds()~ بيحوّل ~time.Duration~ لـ ~float64~ بالثواني، والمقاييس دايماً بالثواني.

## الـ labels المحدودة
كل قيمة label جديدة سلسلة جديدة في الذاكرة لحد ما البرنامج يقفل. القيمة اللي جاية من برا، زي اسم أو ID، ممكن تبقى ملايين. الحل إنك تحوّل أي قيمة مش في قائمة معروفة لقيمة واحدة ثابتة.
~~~go
func relayLabel(r string) string {
 if r == "r7" || r == "r9" { return r }
 return "other" // عدد السلاسل ثابت مهما جه
}
~~~

## الـ endpoint
استورد ~github.com/prometheus/client_golang/prometheus/promhttp~. ~promhttp.HandlerFor(reg, promhttp.HandlerOpts{})~ بيرجع ~http.Handler~ بيكتب كل مقاييس السجل ده بصيغة نصية. كل سطر اسم و labels مرتبة بالاسم وقيمة.
~~~
relay_pings_total{relay="r7"} 3
relay_reply_seconds_bucket{relay="r7",le="0.5"} 0
~~~
~promhttp.Handler()~ بيقرا السجل العام، مش بتاعك.

## التحديث المتوازي
العدادات والمؤشرات آمنة لو اتنادت من كذا goroutine في نفس الوقت؛ مش محتاج قفل حواليها. الدالة اللي بترجع «اقفل اللي فتحته» ممكن تتنادى مرتين بالغلط، و ~sync.Once~ من فصل 25 بيخلّيها تشتغل مرة واحدة.

## أخطاء شائعة
~~~go
prometheus.MustRegister(c)                  // السجل العام: التست التاني بيلاقيه
pings.WithLabelValues(shipName).Inc()       // label من غير حد
lat.Observe(d)                              // Duration مش ثواني، ومش بتتبني
return promhttp.Handler()                   // بيعرض السجل الغلط
~~~

## المهمة
نفّذ ~NewDockMetrics(reg *prometheus.Registry) (*DockMetrics, error)~ يسجّل في ~reg~ بس، ويرجع خطأ التسجيل لو حصل (زي نداء تاني على نفس السجل). المقاييس:
- عداد ~dock_landings_total~ بـ labels ~bay~ و ~outcome~. بعد الإنشاء على طول فيه 10 سلاسل بقيمة 0: كل bay من ~north~ و ~south~ و ~east~ و ~west~ و ~other~ مع ~ok~ و ~aborted~.
- مدرّج ~dock_landing_seconds~ بـ label ~bay~ وحدود 1 و 5 و 15 و 60 بالظبط.
- مؤشر ~dock_queue_depth~ من غير labels.

~RecordLanding(bay string, ok bool, d time.Duration)~ يزوّد العداد واحد بـ ~outcome~ قيمته ~ok~ لو ~ok~ صح و ~aborted~ لو غلط، ويضيف ~d~ بالثواني للمدرّج. أي bay مش من الأربعة بالظبط (حتى ~North~ أو ~""~) يتسجّل ~other~. المدة السالبة تتحسب 0. ~SetQueue(n int)~ يحط المؤشر على n، والسالب 0. ~Handler() http.Handler~ يعرض نفس السجل بـ ~promhttp~. كل الدوال آمنة بالتزامن. ~main~ فاضية.

النجمة التالتة: ~Inflight(bay string) func()~ يزوّد مؤشر ~dock_inflight~ بـ label ~bay~ (بنفس قاعدة ~other~) واحد، ويرجع دالة بتنقّصه واحد أول مرة تتنادى بس؛ النداء التاني مبيعملش حاجة.`,
		Hints: []string{"ابدأ بـ prometheus.NewCounterVec وNewHistogramVec وNewGauge، وسجّلهم بـ reg.Register ورجّع أول خطأ.", "اعمل دالة صغيرة بتحوّل الـ bay لقيمة من الخمسة، وبعد التسجيل نادي WithLabelValues لكل bay مع ok وaborted عشان السلاسل تظهر بصفر.", "في RecordLanding استخدم max(d, 0).Seconds()، وHandler يرجع promhttp.HandlerFor(reg, promhttp.HandlerOpts{})؛ للنجمة لفّ Dec في sync.Once."},
		En: &LevelText{Title: "Counters", Goal: "Expose landing dock metrics with Prometheus using bounded labels, an endpoint, and concurrent updates", Files: []LabFile{{Name: "main.go", Starter: ch361Starter}},
			Lesson: `# Landing Dock Counters

## Before you write
- If every ship name became a label, how many series would be stored after a month?
- Who reads the metrics, and when?

## Why metrics?
A log tells the story of one request. A metric is a number aggregated over time: how many ships landed, how many are waiting now, how long a landing takes. Prometheus visits the server every few seconds and reads all the numbers from one endpoint.

Import ~github.com/prometheus/client_golang/prometheus~. Every metric is a ~Collector~, and it must be registered in a ~Registry~ to appear.

## A private registry
~prometheus.NewRegistry()~ makes an empty registry of your own. There is a global one called ~prometheus.DefaultRegisterer~, but anything registered there stays for the whole program, and the next test finds it. Here the registry comes in as a parameter.
~~~go
reg := prometheus.NewRegistry()
err := reg.Register(c) // Same name twice: an error
~~~
~reg.Register~ returns an error when a metric name is already registered in the same registry. ~MustRegister~ panics instead; here return the error.

## The counter
A counter only goes up, never down. ~prometheus.NewCounterVec~ takes a name, a help text, and label names. ~WithLabelValues~ returns the series with exactly those values, in the same order as the names.
~~~go
pings := prometheus.NewCounterVec(prometheus.CounterOpts{
 Name: "relay_pings_total", Help: "Pings by relay.",
}, []string{"relay"})
pings.WithLabelValues("r7").Inc()   // 1
pings.WithLabelValues("r7").Add(2)  // 3
~~~
Counter names usually end in ~_total~. Calling ~WithLabelValues~ without ~Inc~ creates the series at 0, so it shows in the endpoint from the start.

## The gauge
A gauge goes up and down: a queue length, or how many are running now. ~prometheus.NewGauge~ is one gauge without labels, and ~prometheus.NewGaugeVec~ has labels.
~~~go
fuel := prometheus.NewGauge(prometheus.GaugeOpts{Name: "tank_fuel_liters", Help: "Fuel."})
fuel.Set(42)
fuel.Inc()
fuel.Dec()
~~~

## The histogram
A histogram counts values into buckets. ~Observe(v)~ increments every bucket whose upper bound ~le~ is greater than or equal to v, plus ~_sum~ and ~_count~. The bounds in ~Buckets~ must ascend, and the ~+Inf~ bucket is added for you.
~~~go
lat := prometheus.NewHistogramVec(prometheus.HistogramOpts{
 Name: "relay_reply_seconds", Help: "Reply time.", Buckets: []float64{0.5, 2},
}, []string{"relay"})
lat.WithLabelValues("r7").Observe(1.2) // Only in le="2" and le="+Inf"
~~~
~d.Seconds()~ turns a ~time.Duration~ into a ~float64~ in seconds, and metrics are always in seconds.

## Bounded labels
Every new label value is a new series kept in memory until the program exits. A value from outside, like a name or an ID, can be millions. The fix is to map any value not in a known list to one fixed value.
~~~go
func relayLabel(r string) string {
 if r == "r7" || r == "r9" { return r }
 return "other" // The series count stays fixed whatever arrives
}
~~~

## The endpoint
Import ~github.com/prometheus/client_golang/prometheus/promhttp~. ~promhttp.HandlerFor(reg, promhttp.HandlerOpts{})~ returns an ~http.Handler~ that writes every metric of that registry in a text format. Each line is a name, labels sorted by name, and a value.
~~~
relay_pings_total{relay="r7"} 3
relay_reply_seconds_bucket{relay="r7",le="0.5"} 0
~~~
~promhttp.Handler()~ reads the global registry, not yours.

## Concurrent updates
Counters and gauges are safe to call from many goroutines at once; you need no lock around them. A function that returns "close what you opened" may be called twice by mistake, and ~sync.Once~ from chapter 25 makes it run only once.

## Common mistakes
~~~go
prometheus.MustRegister(c)                  // Global registry: the next test finds it
pings.WithLabelValues(shipName).Inc()       // An unbounded label
lat.Observe(d)                              // A Duration is not seconds and does not compile
return promhttp.Handler()                   // Serves the wrong registry
~~~

## Your mission
Implement ~NewDockMetrics(reg *prometheus.Registry) (*DockMetrics, error)~ registering only in ~reg~ and returning any registration error (such as a second call on the same registry). The metrics:
- Counter ~dock_landings_total~ with labels ~bay~ and ~outcome~. Right after creation it has 10 series at 0: each bay of ~north~, ~south~, ~east~, ~west~, and ~other~ with ~ok~ and ~aborted~.
- Histogram ~dock_landing_seconds~ with label ~bay~ and bounds exactly 1, 5, 15, and 60.
- Gauge ~dock_queue_depth~ with no labels.

~RecordLanding(bay string, ok bool, d time.Duration)~ adds one to the counter with ~outcome~ ~ok~ when ~ok~ is true and ~aborted~ when false, and adds ~d~ in seconds to the histogram. Any bay that is not exactly one of the four (even ~North~ or ~""~) is recorded as ~other~. A negative duration counts as 0. ~SetQueue(n int)~ sets the gauge to n, and a negative n to 0. ~Handler() http.Handler~ serves the same registry with ~promhttp~. Every method is safe for concurrent use. Leave ~main~ empty.

Third star: ~Inflight(bay string) func()~ adds one to gauge ~dock_inflight~ with label ~bay~ (same ~other~ rule) and returns a function that subtracts one the first time it is called only; a second call does nothing.`,
			Hints: []string{"Start with prometheus.NewCounterVec, NewHistogramVec, and NewGauge, register them with reg.Register, and return the first error.", "Write a small function mapping the bay to one of the five, and after registering call WithLabelValues for every bay with ok and aborted so the series appear at zero.", "In RecordLanding use max(d, 0).Seconds(), and Handler returns promhttp.HandlerFor(reg, promhttp.HandlerOpts{}); for the star wrap Dec in a sync.Once."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch361Hidden},
	})
}

const ch361Starter = `package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type DockMetrics struct {
	reg *prometheus.Registry
}

func NewDockMetrics(reg *prometheus.Registry) (*DockMetrics, error) {
	return &DockMetrics{reg: reg}, nil // TODO: register the metrics
}

func (m *DockMetrics) RecordLanding(bay string, ok bool, d time.Duration) {}

func (m *DockMetrics) SetQueue(n int) {}

func (m *DockMetrics) Inflight(bay string) func() { return func() {} }

func (m *DockMetrics) Handler() http.Handler { return http.NotFoundHandler() }

func main() {}
`

const ch361Hidden = `package main

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

var dockTestBays = []string{"north", "south", "east", "west"}

func dockFresh(t *testing.T) (*DockMetrics, *prometheus.Registry) {
	t.Helper()
	reg := prometheus.NewRegistry()
	m, err := NewDockMetrics(reg)
	if err != nil || m == nil {
		t.Fatalf("NewDockMetrics: %v %v", m, err)
	}
	return m, reg
}

// dockScrape reads the endpoint and maps "name{labels}" to its value.
func dockScrape(t *testing.T, m *DockMetrics) map[string]float64 {
	t.Helper()
	w := httptest.NewRecorder()
	m.Handler().ServeHTTP(w, httptest.NewRequest("GET", "/metrics", nil))
	if w.Code != 200 || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") {
		t.Fatalf("endpoint: %d %q", w.Code, w.Header().Get("Content-Type"))
	}
	out := map[string]float64{}
	for _, line := range strings.Split(w.Body.String(), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		i := strings.LastIndex(line, " ")
		v, err := strconv.ParseFloat(line[i+1:], 64)
		if err != nil {
			t.Fatalf("bad line %q", line)
		}
		out[line[:i]] = v
	}
	return out
}

func dockKey(bay, outcome string) string {
	return fmt.Sprintf("dock_landings_total{bay=%q,outcome=%q}", bay, outcome)
}

func dockSeries(s map[string]float64, prefix string) int {
	n := 0
	for k := range s {
		if strings.HasPrefix(k, prefix) {
			n++
		}
	}
	return n
}

func TestHiddenDockRegistration(t *testing.T) {
	m, reg := dockFresh(t)
	if _, err := NewDockMetrics(reg); err == nil {
		t.Fatal("second registration on the same registry must fail")
	}
	if n, err := testutil.GatherAndCount(reg, "dock_landings_total"); err != nil || n != 10 {
		t.Fatalf("fresh counter series = %d %v, want 10", n, err)
	}
	s := dockScrape(t, m)
	for _, bay := range append(dockTestBays, "other") {
		for _, o := range []string{"ok", "aborted"} {
			if v, ok := s[dockKey(bay, o)]; !ok || v != 0 {
				t.Fatalf("%s missing or not zero", dockKey(bay, o))
			}
		}
	}
	if v, ok := s["dock_queue_depth"]; !ok || v != 0 {
		t.Fatal("dock_queue_depth missing")
	}
	families, _ := prometheus.DefaultGatherer.Gather()
	for _, f := range families {
		if strings.HasPrefix(f.GetName(), "dock_") {
			t.Fatal("registered in the global registry:", f.GetName())
		}
	}
	if _, err := NewDockMetrics(prometheus.NewRegistry()); err != nil {
		t.Fatal("a second private registry must work:", err)
	}
}

// Boundaries: bucket edges, negative durations, unknown and mis-cased bays, negative queue.
func TestHiddenDockRecording(t *testing.T) {
	m, reg := dockFresh(t)
	rng := rand.New(rand.NewSource(361))
	counts := map[string]float64{}
	type hist struct{ b [5]float64; sum float64; n float64 }
	hists := map[string]*hist{}
	bounds := []float64{1, 5, 15, 60}
	weird := []string{"North", "", "dock-9", "east ", "WEST", "ship-" + strconv.Itoa(rng.Intn(1e6))}
	durs := []time.Duration{0, time.Second, 5 * time.Second, 15 * time.Second, 60 * time.Second, 60*time.Second + 1, -3 * time.Second, 999 * time.Millisecond}
	for i := 0; i < 400; i++ {
		bay := dockTestBays[rng.Intn(4)]
		label := bay
		if rng.Intn(3) == 0 {
			bay = weird[rng.Intn(len(weird))]
			label = "other"
		}
		ok := rng.Intn(2) == 0
		d := durs[rng.Intn(len(durs))]
		if rng.Intn(2) == 0 {
			d = time.Duration(rng.Int63n(int64(90 * time.Second)))
		}
		m.RecordLanding(bay, ok, d)
		o := "aborted"
		if ok {
			o = "ok"
		}
		counts[dockKey(label, o)]++
		h := hists[label]
		if h == nil {
			h = &hist{}
			hists[label] = h
		}
		sec := max(d, 0).Seconds()
		for j, b := range bounds {
			if sec <= b {
				h.b[j]++
			}
		}
		h.b[4]++
		h.sum += sec
		h.n++
	}
	if n, _ := testutil.GatherAndCount(reg, "dock_landings_total"); n != 10 {
		t.Fatalf("unbounded labels: %d counter series", n)
	}
	s := dockScrape(t, m)
	for k, want := range counts {
		if s[k] != want {
			t.Fatalf("%s = %v, want %v", k, s[k], want)
		}
	}
	if dockSeries(s, "dock_landing_seconds_count{") != len(hists) {
		t.Fatal("histogram has unexpected bay series")
	}
	for bay, h := range hists {
		for j, le := range []string{"1", "5", "15", "60", "+Inf"} {
			k := fmt.Sprintf("dock_landing_seconds_bucket{bay=%q,le=%q}", bay, le)
			if got, ok := s[k]; !ok || got != h.b[j] {
				t.Fatalf("%s = %v, want %v", k, got, h.b[j])
			}
		}
		if got := s[fmt.Sprintf("dock_landing_seconds_count{bay=%q}", bay)]; got != h.n {
			t.Fatalf("%s count = %v, want %v", bay, got, h.n)
		}
		if got := s[fmt.Sprintf("dock_landing_seconds_sum{bay=%q}", bay)]; got < h.sum-1e-6 || got > h.sum+1e-6 {
			t.Fatalf("%s sum = %v, want %v", bay, got, h.sum)
		}
	}
	if dockSeries(s, "dock_landing_seconds_bucket{") != 5*len(hists) {
		t.Fatal("histogram buckets must be exactly 1, 5, 15, 60")
	}
	for _, n := range []int{7, 0, -4, 1 + rng.Intn(500)} {
		m.SetQueue(n)
		if got := dockScrape(t, m)["dock_queue_depth"]; got != float64(max(n, 0)) {
			t.Fatalf("SetQueue(%d) shows %v", n, got)
		}
	}
}

func TestHiddenDockParallel(t *testing.T) {
	m, _ := dockFresh(t)
	var wg sync.WaitGroup
	for g := 0; g < 32; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				m.RecordLanding(dockTestBays[(g+i)%4], i%2 == 0, time.Duration(i)*time.Second)
				m.SetQueue(i)
				if i%10 == 0 {
					_ = dockScrape(t, m)
				}
			}
		}(g)
	}
	wg.Wait()
	s := dockScrape(t, m)
	total := 0.0
	for _, bay := range dockTestBays {
		total += s[dockKey(bay, "ok")] + s[dockKey(bay, "aborted")]
		if s[dockKey(bay, "ok")] != 200 || s[dockKey(bay, "aborted")] != 200 {
			t.Fatalf("%s: ok=%v aborted=%v, want 200 each", bay, s[dockKey(bay, "ok")], s[dockKey(bay, "aborted")])
		}
	}
	if total != 1600 {
		t.Fatal("lost updates", total)
	}
}

func TestHiddenBonusDockInflight(t *testing.T) {
	m, _ := dockFresh(t)
	key := func(bay string) string { return fmt.Sprintf("dock_inflight{bay=%q}", bay) }
	a := m.Inflight("north")
	b := m.Inflight("Spoof")
	c := m.Inflight("north")
	s := dockScrape(t, m)
	if s[key("north")] != 2 || s[key("other")] != 1 {
		t.Fatal("inflight up", s[key("north")], s[key("other")])
	}
	a()
	a()
	b()
	s = dockScrape(t, m)
	if s[key("north")] != 1 || s[key("other")] != 0 {
		t.Fatal("done must subtract once only", s[key("north")], s[key("other")])
	}
	c()
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			done := m.Inflight(dockTestBays[i%4])
			var inner sync.WaitGroup
			for j := 0; j < 3; j++ {
				inner.Add(1)
				go func() { defer inner.Done(); done() }()
			}
			inner.Wait()
		}(i)
	}
	wg.Wait()
	s = dockScrape(t, m)
	for _, bay := range append(dockTestBays, "other") {
		if s[key(bay)] != 0 {
			t.Fatalf("%s inflight = %v after all done", bay, s[key(bay)])
		}
	}
}
`

func init() {
	add(&Level{ID: "36-2-containers", Chapter: 36, Kind: "lab", Title: "مدير الحاويات", Goal: "اكتب أداة سطر أوامر بتعرض وتشغّل وتوقف كبسولات المستعمرة عن طريق interface، بتأكيد قبل الإيقاف ومن غير حذف",
		Requires: []string{"iface"}, Bridge: []int{92}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch362Starter}},
		Lesson: `# مدير كبسولات المستعمرة

## قبل ما تكتب
- لو التست شغّل كبسولة حقيقية بالغلط، مين هيوقفها؟
- لو المشغّل ضغط Enter من غير ما يكتب حاجة، الإيقاف يحصل ولا لأ؟

## الـ runtime من ورا interface
كبسولات المستعمرة شغّالة على runtime حقيقي، زي محرّك حاويات. الأداة بتاعتك متعرفش عنه حاجة غير ثلاث عمليات: عرض وتشغيل وإيقاف. اكتبهم كـ interface، والتست يدّيك نسخة مزيفة (fake) في الذاكرة.
~~~go
type Hatch interface {
 Open(ctx context.Context, id string) error
}
type fakeHatch struct{ opened []string }
func (f *fakeHatch) Open(ctx context.Context, id string) error {
 f.opened = append(f.opened, id)
 return nil
}
~~~
الـ interface الصغير اللي فيه بس اللي الأداة محتاجاه بيحميك: لو الـ runtime الحقيقي فيه عملية حذف، الأداة متقدرش توصلها لأنها مش في الـ interface.

## الـ context لحد آخر نداء
كل عملية على الـ runtime ممكن تاخد وقت، فبتاخد ~ctx~ أول parameter. مرّر نفس ~ctx~ اللي جالك؛ متعملش ~context.Background()~ جوّه الدالة، عشان الإلغاء من برا يوصل.
~~~go
func openAll(ctx context.Context, h Hatch, ids []string) error {
 for _, id := range ids {
  if err := h.Open(ctx, id); err != nil { return err }
 }
 return nil
}
~~~

## الأوامر والـ exit code
الأداة دالة بتاخد ~args~ و ~in~ و ~out~ و ~errOut~ وبترجع رقم الخروج، زي فصل 21. النتيجة العادية على ~out~، والأخطاء على ~errOut~. 0 نجاح، و 1 فشل في التنفيذ، و 2 استخدام غلط.
~~~go
if len(args) == 0 {
 fmt.Fprint(errOut, "usage: hatchctl open ID\n")
 return 2
}
~~~
الأمر الغلط بيرجع 2 قبل ما يلمس الـ runtime خالص.

## التأكيد قبل الفعل الخطير
الإيقاف بيقطع شغل ناس. قبله اطبع سؤال واقرا سطر واحد من ~in~ بـ ~bufio.Scanner~. الافتراضي «لأ»: أي حاجة غير ~y~ أو ~yes~، حتى السطر الفاضي أو نهاية الدخل، معناها إلغاء.
~~~go
fmt.Fprintf(out, "open %s? [y/N] ", id)
sc := bufio.NewScanner(in)
sc.Scan() // false عند نهاية الدخل، و sc.Text() بيبقى ""
answer := strings.ToLower(strings.TrimSpace(sc.Text()))
~~~
الـ N الكبيرة في ~[y/N]~ معناها إن Enter لوحده بيختار «لأ». الـ flag ~--yes~ بيتخطى السؤال للسكربتات، ولما يتبعت متقراش ~in~ أصلاً.

## الحالة قبل الفعل
قبل التشغيل أو الإيقاف، اعرض القائمة ودوّر على الاسم. الاسم الغائب خطأ، والكبسولة اللي في الحالة المطلوبة خلاص مش محتاجة نداء. ~slices.IndexFunc~ بيرجع أول index بيحقق الشرط أو -1.
~~~go
i := slices.IndexFunc(list, func(p Probe) bool { return p.ID == id })
~~~

## أخطاء شائعة
~~~go
var rt = realEngine{}                       // التست بيلمس runtime حقيقي
rt.Stop(context.Background(), name)         // الإلغاء من برا مبيوصلش
if answer != "n" { rt.Stop(ctx, name) }     // السطر الفاضي بقى موافقة
case "rm": rt.(interface{ Remove(string) }) // باب خلفي للحذف
~~~

## المهمة
عرّف ~Runtime~ interface فيه بالظبط ~List(ctx context.Context) ([]Pod, error)~ و ~Start(ctx context.Context, name string) error~ و ~Stop(ctx context.Context, name string) error~. نفّذ ~RunPodCtl(ctx context.Context, rt Runtime, args []string, in io.Reader, out, errOut io.Writer) int~ ومرّر نفس ~ctx~ لكل نداء.

- ~list~: سطر لكل كبسولة ~NAME\trunning~ أو ~NAME\tstopped~ مرتبين بالاسم تصاعدياً، والخروج 0. القائمة الفاضية مبتطبعش حاجة.
- ~start NAME~: الاسم مش في القائمة: ~unknown pod NAME\n~ على ~errOut~ والخروج 1. شغّالة خلاص: ~already running NAME\n~ والخروج 0 من غير نداء ~Start~. غير كده نادي ~Start~ واطبع ~started NAME\n~.
- ~stop NAME~: نفس قاعدة الاسم الغائب، والواقفة خلاص ~already stopped NAME\n~ والخروج 0 من غير سؤال. غير كده اطبع ~stop NAME? [y/N] ~ على ~out~ واقرا سطر واحد. ~y~ أو ~yes~ بعد trim ومن غير فرق حروف كبيرة وصغيرة: نادي ~Stop~ واطبع ~stopped NAME\n~. أي إجابة تانية أو مفيش إجابة: ~aborted\n~ والخروج 1 من غير ~Stop~.
- ~stop --yes NAME~: من غير سؤال ومن غير قراية ~in~.
- أي خطأ من الـ runtime: ~runtime: ERR\n~ على ~errOut~ والخروج 1.
- أي حاجة تانية، منها ~rm~ و ~delete~ وعدد args غلط و ~stop NAME --yes~: ~usage: podctl list [--running] | start NAME | stop [--yes] NAME\n~ على ~errOut~ والخروج 2، من غير أي نداء على الـ runtime. مفيش حذف خالص.

النجمة التالتة: ~list --running~ يطبع الشغّالة بس بنفس الشكل والترتيب.`,
		Hints: []string{"اكتب الـ interface بالتلات methods، وابدأ RunPodCtl بـ switch على args[0] وعدد الباقي، وأي حالة مش معروفة ترجع رسالة الاستخدام و 2.", "لـ start و stop نادي List الأول ودوّر بـ slices.IndexFunc؛ لو الحالة هي المطلوبة خلاص اطبع already ورجّع 0.", "في stop من غير --yes اطبع السؤال، واقرا سطر بـ bufio.Scanner، وقارن strings.ToLower(strings.TrimSpace(...)) بـ y و yes بس."},
		En: &LevelText{Title: "Container Manager", Goal: "Write a command-line tool that lists, starts, and stops colony pods through an interface, confirming before a stop and never deleting", Files: []LabFile{{Name: "main.go", Starter: ch362Starter}},
			Lesson: `# The Colony Pod Manager

## Before you write
- If a test started a real pod by mistake, who would stop it?
- If the operator presses Enter without typing anything, does the stop happen?

## The runtime behind an interface
Colony pods run on a real runtime, like a container engine. Your tool knows only three operations about it: list, start, and stop. Write them as an interface, and the test hands you an in-memory fake.
~~~go
type Hatch interface {
 Open(ctx context.Context, id string) error
}
type fakeHatch struct{ opened []string }
func (f *fakeHatch) Open(ctx context.Context, id string) error {
 f.opened = append(f.opened, id)
 return nil
}
~~~
A small interface holding only what the tool needs protects you: if the real runtime has a delete operation, the tool cannot reach it because it is not in the interface.

## The context to the last call
Every runtime operation may take time, so it takes ~ctx~ as the first parameter. Pass the same ~ctx~ you received; never create ~context.Background()~ inside the function, so a cancel from outside reaches it.
~~~go
func openAll(ctx context.Context, h Hatch, ids []string) error {
 for _, id := range ids {
  if err := h.Open(ctx, id); err != nil { return err }
 }
 return nil
}
~~~

## Commands and the exit code
The tool is a function taking ~args~, ~in~, ~out~, and ~errOut~ and returning the exit code, as in chapter 21. Normal output goes to ~out~, errors to ~errOut~. 0 is success, 1 is a failure while running, and 2 is wrong usage.
~~~go
if len(args) == 0 {
 fmt.Fprint(errOut, "usage: hatchctl open ID\n")
 return 2
}
~~~
A wrong command returns 2 before touching the runtime at all.

## Confirming before a dangerous action
A stop cuts people's work. Before it, print a question and read one line from ~in~ with ~bufio.Scanner~. The default is "no": anything other than ~y~ or ~yes~, even an empty line or the end of input, means cancel.
~~~go
fmt.Fprintf(out, "open %s? [y/N] ", id)
sc := bufio.NewScanner(in)
sc.Scan() // false at end of input, and sc.Text() is ""
answer := strings.ToLower(strings.TrimSpace(sc.Text()))
~~~
The capital N in ~[y/N]~ means Enter alone picks "no". The ~--yes~ flag skips the question for scripts, and when it is given, never read ~in~ at all.

## The state before the action
Before a start or stop, list the pods and look up the name. A missing name is an error, and a pod already in the wanted state needs no call. ~slices.IndexFunc~ returns the first index matching the condition, or -1.
~~~go
i := slices.IndexFunc(list, func(p Probe) bool { return p.ID == id })
~~~

## Common mistakes
~~~go
var rt = realEngine{}                       // The test touches a real runtime
rt.Stop(context.Background(), name)         // An outside cancel never arrives
if answer != "n" { rt.Stop(ctx, name) }     // An empty line became a yes
case "rm": rt.(interface{ Remove(string) }) // A back door to deletion
~~~

## Your mission
Define a ~Runtime~ interface with exactly ~List(ctx context.Context) ([]Pod, error)~, ~Start(ctx context.Context, name string) error~, and ~Stop(ctx context.Context, name string) error~. Implement ~RunPodCtl(ctx context.Context, rt Runtime, args []string, in io.Reader, out, errOut io.Writer) int~ and pass the same ~ctx~ to every call.

- ~list~: one line per pod, ~NAME\trunning~ or ~NAME\tstopped~, in ascending name order, exit 0. An empty list prints nothing.
- ~start NAME~: a name not in the list prints ~unknown pod NAME\n~ to ~errOut~ and exits 1. Already running prints ~already running NAME\n~ and exits 0 without calling ~Start~. Otherwise call ~Start~ and print ~started NAME\n~.
- ~stop NAME~: the same missing-name rule, and an already stopped pod prints ~already stopped NAME\n~ and exits 0 with no question. Otherwise print ~stop NAME? [y/N] ~ to ~out~ and read one line. ~y~ or ~yes~ after trimming, ignoring case: call ~Stop~ and print ~stopped NAME\n~. Any other answer or no answer: ~aborted\n~, exit 1, no ~Stop~.
- ~stop --yes NAME~: no question and no reading of ~in~.
- Any runtime error: ~runtime: ERR\n~ to ~errOut~ and exit 1.
- Anything else, including ~rm~, ~delete~, a wrong argument count, and ~stop NAME --yes~: ~usage: podctl list [--running] | start NAME | stop [--yes] NAME\n~ to ~errOut~ and exit 2, with no runtime call at all. There is no deletion.

Third star: ~list --running~ prints only running pods in the same form and order.`,
			Hints: []string{"Write the interface with the three methods, and start RunPodCtl with a switch on args[0] and the count of the rest; any unknown case prints the usage and returns 2.", "For start and stop call List first and search with slices.IndexFunc; when the state is already the wanted one print already and return 0.", "In stop without --yes print the question, read one line with bufio.Scanner, and compare strings.ToLower(strings.TrimSpace(...)) with y and yes only."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch362Hidden},
	})
}

const ch362Starter = `package main

import (
	"context"
	"io"
)

type Pod struct {
	Name    string
	Running bool
}

// TODO: list the operations the tool needs.
type Runtime interface{}

func RunPodCtl(ctx context.Context, rt Runtime, args []string, in io.Reader, out, errOut io.Writer) int {
	return 2
}

func main() {}
`

const ch362Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"testing"
)

type podCtxKey struct{}

type fakePods struct {
	t        *testing.T
	pods     map[string]bool
	calls    []string
	failList error
	failDo   error
}

func (f *fakePods) check(ctx context.Context) {
	if ctx.Value(podCtxKey{}) != "deck-7" {
		f.t.Error("runtime called without the caller's context")
	}
}

func (f *fakePods) List(ctx context.Context) ([]Pod, error) {
	f.check(ctx)
	f.calls = append(f.calls, "list")
	if f.failList != nil {
		return nil, f.failList
	}
	var out []Pod
	for name, run := range f.pods { // map order: the tool must sort
		out = append(out, Pod{name, run})
	}
	return out, nil
}

func (f *fakePods) Start(ctx context.Context, name string) error {
	f.check(ctx)
	f.calls = append(f.calls, "start "+name)
	if f.failDo != nil {
		return f.failDo
	}
	f.pods[name] = true
	return nil
}

func (f *fakePods) Stop(ctx context.Context, name string) error {
	f.check(ctx)
	f.calls = append(f.calls, "stop "+name)
	if f.failDo != nil {
		return f.failDo
	}
	f.pods[name] = false
	return nil
}

// Remove exists on the fake but must never be reachable.
func (f *fakePods) Remove(ctx context.Context, name string) error {
	f.t.Error("Remove must never be called")
	delete(f.pods, name)
	return nil
}

type noRead struct{ t *testing.T }

func (r noRead) Read([]byte) (int, error) {
	r.t.Error("--yes must not read the input")
	return 0, io.EOF
}

type podRun struct {
	code     int
	out, err string
}

func podCall(f *fakePods, in io.Reader, args ...string) podRun {
	var out, errOut strings.Builder
	ctx := context.WithValue(context.Background(), podCtxKey{}, "deck-7")
	code := RunPodCtl(ctx, f, args, in, &out, &errOut)
	return podRun{code, out.String(), errOut.String()}
}

func podFleet(t *testing.T, rng *rand.Rand) *fakePods {
	f := &fakePods{t: t, pods: map[string]bool{}}
	n := 3 + rng.Intn(8)
	for len(f.pods) < n {
		f.pods[fmt.Sprintf("%s-%d", []string{"hydro", "relay", "forge", "Ark", "lab"}[rng.Intn(5)], rng.Intn(100))] = rng.Intn(2) == 0
	}
	return f
}

func podListing(f *fakePods, onlyRunning bool) string {
	var names []string
	for name := range f.pods {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	for _, name := range names {
		if f.pods[name] {
			b.WriteString(name + "\trunning\n")
		} else if !onlyRunning {
			b.WriteString(name + "\tstopped\n")
		}
	}
	return b.String()
}

func podExpect(t *testing.T, label string, got, want podRun) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %+v, want %+v", label, got, want)
	}
}

const podUsageText = "usage: podctl list [--running] | start NAME | stop [--yes] NAME\n"

func podPick(f *fakePods, running bool) string {
	var names []string
	for name, r := range f.pods {
		if r == running {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		return ""
	}
	return names[0]
}

func TestHiddenPodList(t *testing.T) {
	rng := rand.New(rand.NewSource(362))
	for i := 0; i < 20; i++ {
		f := podFleet(t, rng)
		podExpect(t, "list", podCall(f, strings.NewReader(""), "list"), podRun{0, podListing(f, false), ""})
	}
	empty := &fakePods{t: t, pods: map[string]bool{}}
	podExpect(t, "empty list", podCall(empty, strings.NewReader(""), "list"), podRun{0, "", ""})
	empty.failList = errors.New("engine offline")
	podExpect(t, "list failure", podCall(empty, strings.NewReader(""), "list"), podRun{1, "", "runtime: engine offline\n"})
}

// Boundaries: unknown names, already-in-state pods, runtime failures.
func TestHiddenPodStart(t *testing.T) {
	rng := rand.New(rand.NewSource(3621))
	for i := 0; i < 20; i++ {
		f := podFleet(t, rng)
		podExpect(t, "unknown", podCall(f, nil, "start", "ghost-1"), podRun{1, "", "unknown pod ghost-1\n"})
		if name := podPick(f, true); name != "" {
			f.calls = nil
			podExpect(t, "already running", podCall(f, nil, "start", name), podRun{0, "already running " + name + "\n", ""})
			if len(f.calls) != 1 {
				t.Fatal("start on a running pod called the runtime:", f.calls)
			}
		}
		if name := podPick(f, false); name != "" {
			podExpect(t, "start", podCall(f, nil, "start", name), podRun{0, "started " + name + "\n", ""})
			if !f.pods[name] {
				t.Fatal("Start not called")
			}
		}
	}
	f := &fakePods{t: t, pods: map[string]bool{"hydro": false}, failDo: errors.New("no power")}
	podExpect(t, "start failure", podCall(f, nil, "start", "hydro"), podRun{1, "", "runtime: no power\n"})
	f.failList = errors.New("engine offline")
	podExpect(t, "start list failure", podCall(f, nil, "start", "hydro"), podRun{1, "", "runtime: engine offline\n"})
}

// Boundaries: empty input, EOF, case and spaces, only the first line counts, --yes.
func TestHiddenPodStopConfirm(t *testing.T) {
	yes := []string{"y\n", "Y\n", "yes\n", "YES", "  yEs  \n", "y\nn\n"}
	no := []string{"", "\n", "n\n", "no\n", "yess\n", "ye\n", "maybe\n", "\ny\n", "y es\n", "sure\n"}
	for i, answer := range append(yes, no...) {
		f := &fakePods{t: t, pods: map[string]bool{"forge-1": true, "lab-2": false}}
		got := podCall(f, strings.NewReader(answer), "stop", "forge-1")
		if i < len(yes) {
			podExpect(t, fmt.Sprintf("answer %q", answer), got, podRun{0, "stop forge-1? [y/N] stopped forge-1\n", ""})
			if f.pods["forge-1"] {
				t.Fatal("Stop not called")
			}
		} else {
			podExpect(t, fmt.Sprintf("answer %q", answer), got, podRun{1, "stop forge-1? [y/N] aborted\n", ""})
			if !f.pods["forge-1"] || len(f.calls) != 1 {
				t.Fatal("aborted stop touched the runtime:", f.calls)
			}
		}
	}
	f := &fakePods{t: t, pods: map[string]bool{"forge-1": true, "lab-2": false}}
	podExpect(t, "--yes", podCall(f, noRead{t}, "stop", "--yes", "forge-1"), podRun{0, "stopped forge-1\n", ""})
	podExpect(t, "already stopped", podCall(f, noRead{t}, "stop", "lab-2"), podRun{0, "already stopped lab-2\n", ""})
	podExpect(t, "unknown stop", podCall(f, noRead{t}, "stop", "ghost"), podRun{1, "", "unknown pod ghost\n"})
	f = &fakePods{t: t, pods: map[string]bool{"forge-1": true}, failDo: errors.New("jammed")}
	podExpect(t, "stop failure", podCall(f, strings.NewReader("y\n"), "stop", "forge-1"), podRun{1, "stop forge-1? [y/N] ", "runtime: jammed\n"})
}

func TestHiddenPodUsage(t *testing.T) {
	for _, args := range [][]string{nil, {}, {"rm", "forge-1"}, {"delete", "forge-1"}, {"remove"}, {"LIST"}, {"list", "x"}, {"start"},
		{"start", "a", "b"}, {"stop"}, {"stop", "--yes"}, {"stop", "forge-1", "--yes"}, {"stop", "--force", "forge-1"}, {"stop", "-y", "forge-1"},
		{"start", "--yes", "forge-1"}, {"list", "--running", "x"}, {"prune"}} {
		f := &fakePods{t: t, pods: map[string]bool{"forge-1": true}}
		podExpect(t, fmt.Sprintf("%q", args), podCall(f, noRead{t}, args...), podRun{2, "", podUsageText})
		if len(f.calls) != 0 || !f.pods["forge-1"] {
			t.Fatalf("%q touched the runtime: %v", args, f.calls)
		}
	}
}

func TestHiddenBonusPodListRunning(t *testing.T) {
	rng := rand.New(rand.NewSource(3622))
	for i := 0; i < 20; i++ {
		f := podFleet(t, rng)
		podExpect(t, "list --running", podCall(f, nil, "list", "--running"), podRun{0, podListing(f, true), ""})
	}
	f := &fakePods{t: t, pods: map[string]bool{"a": false}}
	podExpect(t, "none running", podCall(f, nil, "list", "--running"), podRun{0, "", ""})
}
`

func init() {
	add(&Level{ID: "36-3-ops", Chapter: 36, Kind: "lab", Title: "جاهزية", Goal: "اكتب سجلات slog منظّمة بتخفي الأسرار، وفحص /livez و /readyz مربوط بقاعدة بيانات القاعدة",
		Requires: []string{"sel:slog.New", "sel:slog.NewJSONHandler"}, Bridge: []int{100}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch363Starter}},
		Lesson: `# جاهزية القاعدة

## قبل ما تكتب
- لو قاعدة البيانات وقعت، السيرفر نفسه لازم يتعاد تشغيله، ولا بس يبطّل يستقبل طلبات؟
- لو حد بعت ~?token=...~ في الرابط، السر ده هيوصل ملف السجل؟

## سجل منظّم
~fmt.Println~ بيطلع سطر نص حر، وصعب تدوّر فيه. السجل المنظّم كل سطر فيه رسالة ثابتة ومفاتيح وقيم، فتقدر تفلتر بـ ~status=503~ مثلاً. استورد ~log/slog~.

~slog.New(h)~ بيعمل ~*slog.Logger~ فوق ~Handler~. ~slog.NewJSONHandler(w, opts)~ بيكتب كل سجل كسطر JSON واحد على أي ~io.Writer~، فالتست يدّيك buffer ويقرا منه.
~~~go
log := slog.New(slog.NewJSONHandler(&buf, nil))
log.Info("docked", "ship", "kai", "bay", 3)
// {"time":"...","level":"INFO","msg":"docked","ship":"kai","bay":3}
~~~
بعد الرسالة بتيجي أزواج: مفتاح نص وبعده قيمة. الرقم بيفضل رقم في الـ JSON.

## المستويات
~Debug~ و ~Info~ و ~Warn~ و ~Error~ أربع مستويات. ~slog.HandlerOptions{Level: slog.LevelInfo}~ بيرمي أي حاجة أقل من ~Info~.
~~~go
opts := &slog.HandlerOptions{Level: slog.LevelInfo}
log := slog.New(slog.NewJSONHandler(&buf, opts))
log.Debug("noise")         // مش بيتكتب
log.Warn("low oxygen", "pct", 18) // "level":"WARN"
~~~

## تعديل كل مفتاح: ReplaceAttr
~ReplaceAttr~ دالة بتتنادى على كل ~slog.Attr~ قبل الكتابة، و ~groups~ فيها أسماء المجموعات اللي المفتاح جوّاها. ترجّع ~slog.Attr~ معدّل، أو ~slog.Attr{}~ الفاضي عشان المفتاح يتشال.
~~~go
opts := &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
 if len(groups) == 0 && a.Key == slog.TimeKey { return slog.Attr{} } // من غير وقت
 if a.Key == "pin" { return slog.String("pin", "***") }
 return a
}}
~~~
~slog.TimeKey~ هو ~"time"~. الوقت بيتغير كل تشغيل، فشيله عشان السطر يبقى ثابت في التست.

## إخفاء الأسرار (redaction)
السجل بيتنسخ ويتبعت ويتقرا من ناس كتير؛ السر اللي يدخله بيتسرّب. إخفاء المفاتيح الحساسة في ~ReplaceAttr~ بيحمي كل نداء في البرنامج، حتى اللي هيتكتب بعدين.

~slog.Group("auth", "token", t)~ بيعمل مفتاح جوّه مجموعة: ~{"auth":{"token":"..."}}~. ~ReplaceAttr~ بتتنادى على المفاتيح اللي جوّه المجموعة كمان، و ~groups~ قيمتها ~["auth"]~. قارن المفتاح من غير فرق حروف: ~Token~ زي ~token~.
~~~go
log.Info("auth", slog.Group("crew", "Token", "x1")) // لازم يتخفى برضه
~~~

## الحياة والجاهزية
~/livez~ بيجاوب سؤال: البروسيس عايش ولا عالق؟ لو فشل، المنسّق بيعيد تشغيله. عشان كده ميعتمدش على أي حاجة برا؛ لو القاعدة وقعت وكل نسخة اتعاد تشغيلها، المشكلة بتكبر.

~/readyz~ بيجاوب: أبعتلك طلبات دلوقتي؟ ده بيفحص الاعتماديات، زي ~Ping~ على قاعدة البيانات. لو فشل، المنسّق بيبطّل يبعت طلبات بس.
~~~go
type Pinger interface{ Ping(ctx context.Context) error }
~~~
الـ ~*sql.DB~ من فصل 29 عنده ~PingContext~؛ هنا الـ interface صغير والتست بيدّيك fake.

## الفحص بمهلة
فحص الجاهزية لازم يرد بسرعة حتى لو القاعدة عالقة. ~context.WithTimeout~ من فصل 13 فوق ~r.Context()~ بيدّي الـ Ping مهلة، والـ Pinger الكويس بيرجع لما المهلة تخلص.
~~~go
ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
defer cancel()
err := db.Ping(ctx)
~~~

## سطر لكل طلب
بعد الرد اكتب سطر ~Info~ فيه الطريقة والمسار والحالة. ~r.URL.Path~ هو المسار من غير الـ query؛ ~r.RequestURI~ و ~r.URL.String()~ فيهم الـ query بأسرارها.

## أخطاء شائعة
~~~go
log.Info("login", "token", tok)            // من غير ReplaceAttr: السر في الملف
if a.Key == "token" { ... }                 // Token بالكبير بيعدّي
log.Info("request", "path", r.RequestURI)  // ?token=... اتسجّل
func livez(...) { db.Ping(ctx) }            // القاعدة وقعت: كل النسخ بتتعاد
db.Ping(context.Background())               // القاعدة عالقة: الفحص عالق
~~~

## المهمة
نفّذ ~NewOpsLogger(w io.Writer) *slog.Logger~ بـ ~slog.New~ و ~slog.NewJSONHandler~ على ~w~: مستوى ~Info~ وأعلى بس، ومن غير مفتاح ~time~ في المستوى الأعلى. أي مفتاح اسمه ~password~ أو ~token~ أو ~secret~ أو ~authorization~ من غير فرق حروف، في أي مجموعة، قيمته تبقى النص ~[REDACTED]~ ومفتاحه زي ما هو. باقي المفاتيح زي ما هي.

نفّذ ~NewOps(log *slog.Logger, db Pinger) *Ops~ و ~(*Ops).ServeHTTP~. الجسم ~text/plain~:
- ~GET /livez~: 200 و ~ok\n~ ومن غير ~Ping~ أبداً.
- ~GET /readyz~: ~Ping~ واحد بـ context مشتق من ~r.Context()~ بمهلة 200ms. النجاح 200 و ~ready\n~. الفشل 503 و ~not ready\n~، وقبله سجل ~Warn~ رسالته ~readiness failed~ ومفتاح ~error~ نص الخطأ.
- طريقة غير ~GET~ على المسارين: 405 و ~method not allowed\n~. أي مسار تاني: 404 و ~not found\n~.
- بعد كل رد، سجل ~Info~ رسالته ~request~ ومفاتيح ~method~ و ~path~ (من غير query) و ~status~ رقم.

آمن بالتزامن، ومن غير شبكة. ~main~ فاضية.

النجمة التالتة: ~(*Ops).Drain()~ آمنة بالتزامن؛ بعدها ~/readyz~ يرجع 503 و ~draining\n~ من غير ~Ping~ ومن غير سجل ~Warn~، و ~/livez~ يفضل 200.`,
		Hints: []string{"ابدأ بـ slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: ...})) وشيل time لما groups فاضية.", "في ReplaceAttr قارن strings.ToLower(a.Key) بالأربع مفاتيح ورجّع slog.String(a.Key, \"[REDACTED]\").", "في ServeHTTP احسب الحالة والجسم الأول، واكتبهم، وبعدين log.Info(\"request\", ...) بـ r.URL.Path؛ للجاهزية context.WithTimeout(r.Context(), 200*time.Millisecond)، وللنجمة atomic.Bool."},
		En: &LevelText{Title: "Readiness", Goal: "Write structured slog logs that hide secrets, and /livez and /readyz probes tied to the base database", Files: []LabFile{{Name: "main.go", Starter: ch363Starter}},
			Lesson: `# Base Readiness

## Before you write
- If the database goes down, should the server itself be restarted, or just stop receiving requests?
- If someone sends ~?token=...~ in the URL, does that secret reach the log file?

## Structured logs
~fmt.Println~ prints a line of free text that is hard to search. In a structured log every line has a fixed message plus keys and values, so you can filter on ~status=503~, for example. Import ~log/slog~.

~slog.New(h)~ makes a ~*slog.Logger~ on top of a ~Handler~. ~slog.NewJSONHandler(w, opts)~ writes each record as one JSON line to any ~io.Writer~, so the test hands you a buffer and reads from it.
~~~go
log := slog.New(slog.NewJSONHandler(&buf, nil))
log.Info("docked", "ship", "kai", "bay", 3)
// {"time":"...","level":"INFO","msg":"docked","ship":"kai","bay":3}
~~~
After the message come pairs: a string key, then a value. A number stays a number in the JSON.

## Levels
~Debug~, ~Info~, ~Warn~, and ~Error~ are four levels. ~slog.HandlerOptions{Level: slog.LevelInfo}~ drops anything below ~Info~.
~~~go
opts := &slog.HandlerOptions{Level: slog.LevelInfo}
log := slog.New(slog.NewJSONHandler(&buf, opts))
log.Debug("noise")         // Not written
log.Warn("low oxygen", "pct", 18) // "level":"WARN"
~~~

## Editing every key: ReplaceAttr
~ReplaceAttr~ is a function called on every ~slog.Attr~ before writing, and ~groups~ holds the names of the groups the key sits in. Return a changed ~slog.Attr~, or the empty ~slog.Attr{}~ to remove the key.
~~~go
opts := &slog.HandlerOptions{ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
 if len(groups) == 0 && a.Key == slog.TimeKey { return slog.Attr{} } // No time
 if a.Key == "pin" { return slog.String("pin", "***") }
 return a
}}
~~~
~slog.TimeKey~ is ~"time"~. The time changes on every run, so remove it to keep the line stable in tests.

## Hiding secrets (redaction)
Logs are copied, shipped, and read by many people; a secret that enters one leaks. Hiding sensitive keys in ~ReplaceAttr~ protects every call in the program, even ones written later.

~slog.Group("auth", "token", t)~ makes a key inside a group: ~{"auth":{"token":"..."}}~. ~ReplaceAttr~ is called on the keys inside the group too, with ~groups~ set to ~["auth"]~. Compare the key ignoring case: ~Token~ is the same as ~token~.
~~~go
log.Info("auth", slog.Group("crew", "Token", "x1")) // Must be hidden too
~~~

## Liveness and readiness
~/livez~ answers: is the process alive or stuck? If it fails, the orchestrator restarts it. That is why it depends on nothing outside; if the database goes down and every copy gets restarted, the problem grows.

~/readyz~ answers: should I send you requests now? It checks dependencies, like a ~Ping~ to the database. If it fails, the orchestrator only stops sending requests.
~~~go
type Pinger interface{ Ping(ctx context.Context) error }
~~~
The ~*sql.DB~ from chapter 29 has ~PingContext~; here the interface is small and the test gives you a fake.

## A probe with a deadline
A readiness probe must answer fast even when the database hangs. ~context.WithTimeout~ from chapter 13 over ~r.Context()~ gives the Ping a deadline, and a good Pinger returns when the deadline passes.
~~~go
ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
defer cancel()
err := db.Ping(ctx)
~~~

## One line per request
After the reply, write an ~Info~ line with the method, the path, and the status. ~r.URL.Path~ is the path without the query; ~r.RequestURI~ and ~r.URL.String()~ include the query and its secrets.

## Common mistakes
~~~go
log.Info("login", "token", tok)            // No ReplaceAttr: the secret is in the file
if a.Key == "token" { ... }                 // A capital Token slips through
log.Info("request", "path", r.RequestURI)  // ?token=... got logged
func livez(...) { db.Ping(ctx) }            // Database down: every copy restarts
db.Ping(context.Background())               // Database hangs: the probe hangs
~~~

## Your mission
Implement ~NewOpsLogger(w io.Writer) *slog.Logger~ with ~slog.New~ and ~slog.NewJSONHandler~ on ~w~: only ~Info~ and above, and no top-level ~time~ key. Any key named ~password~, ~token~, ~secret~, or ~authorization~, ignoring case, in any group, gets the text value ~[REDACTED]~ and keeps its key. Every other key is unchanged.

Implement ~NewOps(log *slog.Logger, db Pinger) *Ops~ and ~(*Ops).ServeHTTP~. Bodies are ~text/plain~:
- ~GET /livez~: 200 and ~ok\n~, never calling ~Ping~.
- ~GET /readyz~: one ~Ping~ with a context derived from ~r.Context()~ with a 200ms deadline. Success is 200 and ~ready\n~. Failure is 503 and ~not ready\n~, preceded by a ~Warn~ record with message ~readiness failed~ and key ~error~ holding the error text.
- A method other than ~GET~ on either path: 405 and ~method not allowed\n~. Any other path: 404 and ~not found\n~.
- After every reply, an ~Info~ record with message ~request~ and keys ~method~, ~path~ (no query), and ~status~ as a number.

Safe for concurrent use, no network. Leave ~main~ empty.

Third star: ~(*Ops).Drain()~ is safe for concurrent use; after it, ~/readyz~ returns 503 and ~draining\n~ with no ~Ping~ and no ~Warn~ record, and ~/livez~ stays 200.`,
			Hints: []string{"Start with slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: ...})) and remove time when groups is empty.", "In ReplaceAttr compare strings.ToLower(a.Key) with the four keys and return slog.String(a.Key, \"[REDACTED]\").", "In ServeHTTP compute the status and body first, write them, then log.Info(\"request\", ...) with r.URL.Path; for readiness use context.WithTimeout(r.Context(), 200*time.Millisecond), and for the star an atomic.Bool."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch363Hidden},
	})
}

const ch363Starter = `package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

func NewOpsLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil)) // TODO: JSON on w with redaction
}

type Ops struct{}

func NewOps(log *slog.Logger, db Pinger) *Ops { return &Ops{} }

func (o *Ops) Drain() {}

func (o *Ops) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "TODO", http.StatusNotImplemented)
}

func main() {}
`

const ch363Hidden = `package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"mime"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type opsBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (o *opsBuf) Write(p []byte) (int, error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.b.Write(p)
}

func (o *opsBuf) records(t *testing.T) []map[string]any {
	t.Helper()
	o.mu.Lock()
	defer o.mu.Unlock()
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(o.b.String()), "\n") {
		if line == "" {
			continue
		}
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("log line is not one JSON object: %q", line)
		}
		out = append(out, m)
	}
	return out
}

func (o *opsBuf) text() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.b.String()
}

type fakeDB struct {
	calls atomic.Int32
	fail  atomic.Value // error or nil
	hang  bool
	ctxOK atomic.Bool
}

func (f *fakeDB) Ping(ctx context.Context) error {
	f.calls.Add(1)
	if dl, ok := ctx.Deadline(); ok && time.Until(dl) <= 210*time.Millisecond && ctx.Value(opsKey{}) == "base" {
		f.ctxOK.Store(true)
	} else {
		f.ctxOK.Store(false)
	}
	if f.hang {
		<-ctx.Done()
		return ctx.Err()
	}
	if e, _ := f.fail.Load().(error); e != nil {
		return e
	}
	return nil
}

type opsKey struct{}

func opsDo(h *Ops, method, target string) (int, string, string) {
	r := httptest.NewRequest(method, target, nil)
	r = r.WithContext(context.WithValue(r.Context(), opsKey{}, "base"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	media, _, _ := mime.ParseMediaType(w.Header().Get("Content-Type"))
	return w.Code, w.Body.String(), media
}

func TestHiddenOpsLoggerRedaction(t *testing.T) {
	rng := rand.New(rand.NewSource(363))
	buf := &opsBuf{}
	log := NewOpsLogger(buf)
	secret := fmt.Sprintf("s3cr3t-%d", rng.Intn(1e9))
	log.Debug("hidden", "x", 1)
	log.Info("login", "user", "kai", "token", secret, "Password", secret, "attempt", 2,
		slog.Group("auth", "SECRET", secret, "Authorization", "Bearer "+secret, "scope", "dock"))
	log.Warn("low oxygen", "pct", 18, "tokens", 5)
	log.Error("breach", "secret", secret)
	if strings.Contains(buf.text(), secret) {
		t.Fatalf("secret leaked: %s", buf.text())
	}
	recs := buf.records(t)
	if len(recs) != 3 {
		t.Fatalf("want 3 records (Debug dropped), got %d: %s", len(recs), buf.text())
	}
	for _, r := range recs {
		if _, ok := r["time"]; ok {
			t.Fatal("time key must be removed")
		}
	}
	r := recs[0]
	auth, _ := r["auth"].(map[string]any)
	if r["level"] != "INFO" || r["msg"] != "login" || r["user"] != "kai" || r["attempt"] != 2.0 || r["token"] != "[REDACTED]" || r["Password"] != "[REDACTED]" ||
		auth == nil || auth["SECRET"] != "[REDACTED]" || auth["Authorization"] != "[REDACTED]" || auth["scope"] != "dock" {
		t.Fatalf("login record: %v", r)
	}
	if recs[1]["level"] != "WARN" || recs[1]["pct"] != 18.0 || recs[1]["tokens"] != 5.0 {
		t.Fatalf("only exact key names are redacted: %v", recs[1])
	}
	if recs[2]["level"] != "ERROR" || recs[2]["secret"] != "[REDACTED]" {
		t.Fatalf("error record: %v", recs[2])
	}
}

func opsRequestLine(t *testing.T, r map[string]any, method, path string, status int) {
	t.Helper()
	if r["level"] != "INFO" || r["msg"] != "request" || r["method"] != method || r["path"] != path || r["status"] != float64(status) {
		t.Fatalf("request record %v, want %s %s %d", r, method, path, status)
	}
}

// Boundaries: livez never pings, query secrets, wrong methods and paths, failing and hanging database.
func TestHiddenOpsProbes(t *testing.T) {
	buf := &opsBuf{}
	db := &fakeDB{}
	h := NewOps(NewOpsLogger(buf), db)
	cases := []struct {
		method, target, path string
		status               int
		body                 string
	}{
		{"GET", "/livez", "/livez", 200, "ok\n"},
		{"GET", "/readyz?token=q-s3cret", "/readyz", 200, "ready\n"},
		{"POST", "/readyz", "/readyz", 405, "method not allowed\n"},
		{"DELETE", "/livez", "/livez", 405, "method not allowed\n"},
		{"GET", "/healthz", "/healthz", 404, "not found\n"},
		{"GET", "/readyz/", "/readyz/", 404, "not found\n"},
		{"GET", "/", "/", 404, "not found\n"},
	}
	for _, c := range cases {
		code, body, media := opsDo(h, c.method, c.target)
		if code != c.status || body != c.body || media != "text/plain" {
			t.Fatalf("%s %s: got %d %q %q", c.method, c.target, code, body, media)
		}
	}
	if db.calls.Load() != 1 || !db.ctxOK.Load() {
		t.Fatalf("Ping calls = %d (want 1, only GET /readyz), context from the request with a 200ms deadline: %v", db.calls.Load(), db.ctxOK.Load())
	}
	recs := buf.records(t)
	if len(recs) != len(cases) {
		t.Fatalf("want one record per request, got %d", len(recs))
	}
	for i, c := range cases {
		opsRequestLine(t, recs[i], c.method, c.path, c.status)
	}
	if strings.Contains(buf.text(), "s3cret") {
		t.Fatal("query leaked into the log")
	}

	db.fail.Store(errors.New("base db offline"))
	code, body, _ := opsDo(h, "GET", "/readyz")
	if code != 503 || body != "not ready\n" {
		t.Fatalf("failing db: %d %q", code, body)
	}
	if code, body, _ := opsDo(h, "GET", "/livez"); code != 200 || body != "ok\n" {
		t.Fatal("livez must not depend on the database", code, body)
	}
	recs = buf.records(t)[len(cases):]
	if len(recs) != 3 || recs[0]["level"] != "WARN" || recs[0]["msg"] != "readiness failed" || recs[0]["error"] != "base db offline" {
		t.Fatalf("want Warn then request records: %v", recs)
	}
	opsRequestLine(t, recs[1], "GET", "/readyz", 503)
	opsRequestLine(t, recs[2], "GET", "/livez", 200)

	hung := NewOps(NewOpsLogger(&opsBuf{}), &fakeDB{hang: true})
	start := time.Now()
	code, body, _ = opsDo(hung, "GET", "/readyz")
	if took := time.Since(start); code != 503 || body != "not ready\n" || took > 1500*time.Millisecond {
		t.Fatalf("hanging db: %d %q after %v", code, body, took)
	}
}

func TestHiddenOpsParallel(t *testing.T) {
	buf := &opsBuf{}
	db := &fakeDB{}
	h := NewOps(NewOpsLogger(buf), db)
	var wg sync.WaitGroup
	for i := 0; i < 60; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			path := []string{"/livez", "/readyz", "/nope"}[i%3]
			want := []int{200, 200, 404}[i%3]
			if code, _, _ := opsDo(h, "GET", path); code != want {
				t.Errorf("%s: %d", path, code)
			}
		}(i)
	}
	wg.Wait()
	if n := len(buf.records(t)); n != 60 || db.calls.Load() != 20 {
		t.Fatalf("records=%d pings=%d, want 60 and 20", n, db.calls.Load())
	}
}

func TestHiddenBonusOpsDrain(t *testing.T) {
	buf := &opsBuf{}
	db := &fakeDB{}
	h := NewOps(NewOpsLogger(buf), db)
	if code, _, _ := opsDo(h, "GET", "/readyz"); code != 200 {
		t.Fatal("ready before drain")
	}
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() { defer wg.Done(); h.Drain() }()
		go func() { defer wg.Done(); opsDo(h, "GET", "/readyz"); opsDo(h, "GET", "/livez") }()
	}
	wg.Wait()
	before := db.calls.Load()
	for i := 0; i < 3; i++ {
		code, body, _ := opsDo(h, "GET", "/readyz")
		if code != 503 || body != "draining\n" {
			t.Fatalf("draining readyz: %d %q", code, body)
		}
		if code, body, _ := opsDo(h, "GET", "/livez"); code != 200 || body != "ok\n" {
			t.Fatal("livez while draining", code, body)
		}
	}
	if db.calls.Load() != before {
		t.Fatal("draining readyz must not ping")
	}
	recs := buf.records(t)
	for _, r := range recs[len(recs)-6:] {
		if r["level"] == "WARN" {
			t.Fatal("draining is not a readiness failure:", r)
		}
	}
	opsRequestLine(t, recs[len(recs)-2], "GET", "/readyz", 503)
}
`
