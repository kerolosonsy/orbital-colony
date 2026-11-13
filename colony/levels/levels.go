// Package levels defines every mission: its map, goal and lesson.
//
// Lesson text is a small markdown dialect: "~~~" fences code blocks and
// ~text~ marks inline code (backticks cannot appear in Go raw strings).
package levels

import (
	"io"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

type Chapter struct {
	Quiz      []QuizQuestion `json:"quiz,omitempty"`
	Track     string         `json:"track"`
	TrackEn   string         `json:"trackEn"`
	N         int            `json:"n"`
	Title     string         `json:"title"`
	Concept   string         `json:"concept"`
	Ready     bool           `json:"ready"`
	TitleEn   string         `json:"titleEn"`
	ConceptEn string         `json:"conceptEn"`
}

type LevelText struct {
	Title   string    `json:"title"`
	Goal    string    `json:"goal"`
	Lesson  string    `json:"lesson"`
	Starter string    `json:"starter"`
	Hints   []string  `json:"hints"`
	Files   []LabFile `json:"files,omitempty"`
}

type Level struct {
	Boss      bool       `json:"boss,omitempty"`
	Deps      []string   `json:"deps,omitempty"`
	Service   string     `json:"service,omitempty"`
	CarryFrom string     `json:"carryFrom,omitempty"`
	Bridge    []int      `json:"bridge,omitempty"`
	En        *LevelText `json:"en,omitempty"`
	ID        string     `json:"id"`
	Chapter   int        `json:"chapter"`
	Title     string     `json:"title"`
	Goal      string     `json:"goal"`
	Lesson    string     `json:"lesson"`
	Starter   string     `json:"starter"`
	Hints     []string   `json:"hints"`
	Par       int        `json:"par"`  // ticks for 2 stars
	Gold      int        `json:"gold"` // ticks for 3 stars
	MaxTicks  int        `json:"maxTicks"`
	Requires  []string   `json:"requires"` // source features, checked by the server
	Race      bool       `json:"race"`     // build with -race; a detected race fails the mission

	Setup func(w *sim.World) `json:"-"`

	// Lab levels are about go test: the player edits Files and the server
	// runs test suites instead of a mission.
	Kind  string    `json:"kind,omitempty"` // "" for missions, "lab"
	Files []LabFile `json:"files,omitempty"`
	Lab   *Lab      `json:"-"`
}

type LabFile struct {
	Name     string `json:"name"`
	Starter  string `json:"starter"`
	ReadOnly bool   `json:"readOnly"`
}

type Mutant struct {
	Desc string // shown when the player's tests miss it
	Code string // a broken main.go
}

type Lab struct {
	Vet         bool // run all go vet analyzers as a separate check
	Fuzz        string
	Bonus       bool
	HiddenFiles map[string]string
	ReqFile     string   // file whose source is checked against Requires
	Mutants     []Mutant // the player's tests must fail on each
	Buggy       string   // original main.go the player's tests must fail on
	Hidden      string   // extra test file run against the player's main.go
	Bench       *Bench

	MutantFile string // defaults to main.go; colony planners use a sub-package
}

// Bench compares BenchmarkHiddenNaive with BenchmarkHiddenTheirs from the
// Hidden file; Tiers are the speedups for 1, 2 and 3 stars.
type Bench struct {
	Tiers [3]float64
}

var Chapters = []Chapter{
	{chapterQuizzes[1], "الأساسيات", "Fundamentals", 1, "الهبوط", "استدعاء الدوال وترتيب الأوامر", true, "Landing", "Function calls and command order"},
	{chapterQuizzes[2], "الأساسيات", "Fundamentals", 2, "التليمتري", "المتغيرات والأنواع والحساب", true, "Telemetry", "Variables, types, and arithmetic"},
	{chapterQuizzes[3], "الأساسيات", "Fundamentals", 3, "الحلقات", "for بكل أشكالها", true, "Loops", "Every form of for"},
	{chapterQuizzes[4], "الأساسيات", "Fundamentals", 4, "القرارات", "if و switch", true, "Decisions", "if and switch"},
	{chapterQuizzes[5], "الأساسيات", "Fundamentals", 5, "الدوال والأخطاء", "func و error", true, "Functions and Errors", "func and error"},
	{chapterQuizzes[6], "الأساسيات", "Fundamentals", 6, "المخازن", "slices و maps", true, "Storage", "slices and maps"},
	{chapterQuizzes[7], "الأساسيات", "Fundamentals", 7, "بناء الأدوات", "structs و methods", true, "Building Tools", "structs and methods"},
	{chapterQuizzes[8], "الأساسيات", "Fundamentals", 8, "العناوين", "pointers", true, "Addresses", "pointers"},
	{chapterQuizzes[9], "الأساسيات", "Fundamentals", 9, "البروتوكولات", "interfaces", true, "Protocols", "interfaces"},
	{chapterQuizzes[10], "التزامن", "Concurrency", 10, "السرب", "goroutines و WaitGroup", true, "The Swarm", "goroutines and WaitGroup"},
	{chapterQuizzes[11], "التزامن", "Concurrency", 11, "قنوات الاتصال", "channels و select", true, "Communication Channels", "channels and select"},
	{chapterQuizzes[12], "التزامن", "Concurrency", 12, "السباق", "mutex و race detector", true, "The Race", "mutex and the race detector"},
	{chapterQuizzes[13], "التزامن", "Concurrency", 13, "العواصف", "context و timeouts", true, "Storms", "context and timeouts"},
	{chapterQuizzes[14], "الحِرفة (generics، خوارزميات، تستات)", "Craft (generics, algorithms, tests)", 14, "القوالب", "generics", true, "Templates", "generics"},
	{chapterQuizzes[15], "الحِرفة (generics، خوارزميات، تستات)", "Craft (generics, algorithms, tests)", 15, "الكهوف", "خوارزميات: BFS وsorting وbinary search وrecursion وheap", true, "Caves", "Algorithms: BFS, sorting, binary search, recursion, and heaps"},
	{chapterQuizzes[16], "الحِرفة (generics، خوارزميات، تستات)", "Craft (generics, algorithms, tests)", 16, "المختبر", "go test وbenchmarks وfuzzing وساعة محقونة", true, "The Lab", "go test, benchmarks, fuzzing, and injected time"},
	{chapterQuizzes[17], "الحِرفة (generics، خوارزميات، تستات)", "Craft (generics, algorithms, tests)", 17, "الأعماق", "طبقات، ومصفوفات 3D، وBFS ثلاثي الأبعاد", true, "The Depths", "Layers, 3D arrays, and three-dimensional BFS"},
	{chapterQuizzes[18], "البيانات والأدوات", "Data and Tools", 18, "الأرشيف", "JSON وCSV ومعالجة النصوص والاختبار", true, "The Archive", "JSON, CSV, text processing, and testing"},
	{chapterQuizzes[19], "البيانات والأدوات", "Data and Tools", 19, "الساعة الكونية", "المدد والمناطق الزمنية ومؤقتات قابلة للإلغاء", true, "The Cosmic Clock", "Durations, time zones, and cancelable concurrent timers"},
	{chapterQuizzes[20], "البيانات والأدوات", "Data and Tools", 20, "المترجم", "تحليل النصوص والتحقق وlexer وشجرة التعبير", true, "The Interpreter", "Markup parsing, validation, lexing, and expression trees"},
	{chapterQuizzes[21], "البيانات والأدوات", "Data and Tools", 21, "سطر الأوامر والنظام", "flags والبيئة وأوامر فرعية وألعاب بعشوائية محقونة", true, "The Command Line", "Flags, environment, subcommands, and games with injected randomness"},
	{chapterQuizzes[22], "البيانات والأدوات", "Data and Tools", 22, "نظام الملفات", "شجرة الملفات وخطط النقل والحفظ الذرّي وبصمات الملفات", true, "The File System", "Filesystem trees, move plans, atomic saves, and streaming file hashes"},
	{chapterQuizzes[23], "البيانات والأدوات", "Data and Tools", 23, "الخزنة", "أكواد آمنة وAES-GCM وعنونة المحتوى ولقطات الشحنات", true, "The Vault", "Secure access codes, AES-GCM, content addressing, and supply snapshots"},
	{chapterQuizzes[24], "التزامن", "Concurrency", 24, "خطوط الإمداد", "ملكية القنوات والضغط العكسي وخطوط المعالجة وHTTP والصور المتوازية", true, "Supply Lines", "Channel ownership, backpressure, pipelines, and concurrent HTTP and images"},
	{chapterQuizzes[25], "التزامن", "Concurrency", 25, "الموارد المشتركة", "شروط الرصيد والسباقات وخزان التصاريح والأحداث وكاش بصلاحية والعدّادات الذرّية", true, "Shared Resources", "Invariants, races, token buckets, event buses, TTL caches, and atomic counters"},
	{chapterQuizzes[26], "الخدمات", "Services", 26, "بوابة القاعدة", "سيرفر HTTP وREST وراوتر وmiddleware وأخطاء موحّدة وملفات ثابتة ورفع آمن", true, "Base Gateway", "HTTP servers, REST, routing, middleware, error envelopes, static files, and safe uploads"},
	{chapterQuizzes[27], "الخدمات", "Services", 27, "التصاريح والحماية", "bcrypt وJWT والجلسات والكوكيز وCSRF وCORS وحد الطلبات وروابط مختصرة", true, "Passes and Protection", "bcrypt, JWT, sessions, cookies, CSRF, CORS, rate limiting, and short links"},
	{chapterQuizzes[28], "الخدمات", "Services", 28, "الإطارات والتشغيل", "Gin وFiber وعقد OpenAPI واختبار العقد وإغلاق مرتب بيصفّي الطلبات الجارية", true, "Frameworks and Operations", "Gin, Fiber, OpenAPI contracts, contract testing, and graceful shutdown that drains in-flight requests"},
	{chapterQuizzes[29], "قواعد البيانات", "Databases", 29, "قواعد البيانات", "SQL وSQLite والترحيلات والتحويلات الذرّية وحوض الاتصالات وsqlx وGORM", true, "SQL Databases", "SQL, SQLite, migrations, atomic transfers, connection pools, sqlx, and GORM"},
	{chapterQuizzes[30], "قواعد البيانات", "Databases", 30, "التخزين المتقدم", "Redis والكاش الجانبي ومستندات MongoDB وفهارس البحث وفصل القراءة والكتابة", true, "Advanced Storage", "Redis cache-aside, MongoDB documents, search indexes, and read/write splitting"},
	{chapterQuizzes[31], "الشبكات", "Networking", 31, "الشبكة الخام", "TCP وحدود السطور وغرفة راديو وUDP وإطارات binary وDNS وفحص منافذ محلي", true, "Raw Networking", "TCP, bounded lines, a radio hub, UDP, binary frames, DNS, and loopback port checks"},
	{chapterQuizzes[32], "الشبكات", "Networking", 32, "البروتوكولات", "WebSocket وgRPC والبث وTLS وmTLS وHTTP/2 وSSH", true, "Protocols", "WebSocket, unary and streaming gRPC, TLS, mTLS, HTTP/2, and SSH"},
	{chapterQuizzes[33], "الشبكات", "Networking", 33, "الوسطاء", "عميل بيعيد بأمان وbackoff وjitter ووكيل عكسي وموزّع أحمال بفحص صحة وبوابة API", true, "Intermediaries", "Safe retries with backoff and jitter, reverse proxies, health-checked load balancing, and an API gateway"},
	{chapterQuizzes[34], "الأنظمة", "Systems", 34, "محركات التخزين", "مخزن مفاتيح بصلاحية ولقطات وسجل كتابة مسبقة بـ checksums ودمج SSTables وتخزين chunks موزّع بالتكرار والإصلاح", true, "Storage Engines", "A TTL key store with snapshots, a checksummed write-ahead log, SSTable compaction, and replicated chunk storage with repair"},
	{chapterQuizzes[35], "الأنظمة", "Systems", 35, "الأنظمة الموزعة", "طوابير المهام والمجدول والأقفال الموزعة وصندوق الصادر وKafka وانتخاب Raft ونسخ السجل", true, "Distributed Systems", "Task queues, schedulers, distributed locks, the outbox, Kafka, Raft elections, and log replication"},
	{chapterQuizzes[36], "الأنظمة", "Systems", 36, "التشغيل والمراقبة", "مقاييس Prometheus وlabels محدودة ومدير حاويات بـ interface وسجلات slog منظّمة وفحوص الجاهزية", true, "Operations and Observability", "Prometheus metrics, bounded labels, an interface-driven container manager, structured slog logs, and readiness probes"},
	{chapterQuizzes[37], "Capstone", "Capstone", 37, "مركز القيادة", "مشروع ختامي: schema وترحيلات وجلسات وأطقم معزولة ومهام بمفاتيح تكرار وتشغيل آمن", true, "Command Center (Capstone)", "Capstone: schema, migrations, sessions, isolated crews, idempotent tasks, and safe operations"},
	{chapterQuizzes[38], "إضافي", "Extras", 38, "Go الحديثة", "فصل اختياري بعد 17: iterators وsets وOnce وerrgroup والتخصيصات وvet وpprof", true, "Modern Go", "Optional after chapter 17: iterators, sets, Once, errgroup, allocations, vet, and pprof"},
}

var All []*Level

func add(l *Level) { All = append(All, l) }

func Get(id string) *Level {
	if id == Colony.ID {
		return Colony
	}
	for _, l := range All {
		if l.ID == id {
			return l
		}
	}
	return nil
}

var sandbox = &Level{
	ID: "sandbox", MaxTicks: 5000,
	Setup: func(w *sim.World) {
		w.Init(8, 8, sim.Regolith)
		w.Check = func(*sim.World) (bool, string) { return false, i18n.T("وضع حر", "Free play") }
	},
}

// Start builds the world for level id (or a free sandbox) and emits its
// initial snapshot to out.
func Start(id string, out io.Writer, seed int64) *sim.World {
	l := Get(id)
	if l == nil {
		l = sandbox
	}
	w := sim.New(out, seed)
	w.MaxTicks = l.MaxTicks
	l.Setup(w)
	w.Start()
	return w
}

func countEntity(w *sim.World, e sim.Entity) int {
	n := 0
	w.Each(func(_, _ int, t *sim.Tile) {
		if t.Entity == e {
			n++
		}
	})
	return n
}

func reported(w *sim.World, want string) bool {
	for _, r := range w.Reports {
		if r == want {
			return true
		}
	}
	return false
}

// Narrate prints the start of the mission for terminal runs. Start already
// emitted the init event before the console was attached.
func Narrate(w *sim.World) {
	_, progress := w.Check(w)
	w.Console.Write([]byte("🛰️  " + progress + "\n"))
}

// Colony is the open, endless colony. It is not a mission: the server runs
// it through its own endpoints and the program loads the saved economy.
var Colony = &Level{
	En: &LevelText{
		Title: "The Colony",
		Goal:  "Grow your colony: plant, harvest, sell, and buy upgrades",
		Lesson: `# Your Colony 🌱

Grow your colony by planting, harvesting, selling goods, and buying upgrades.

## Your starter program
The drone visits the field in a snake pattern. At each tile it tills, harvests, and plants moss. Once it has 5 or more moss, it exports them to the market.

## Keep improving
The program runs until you press Stop, and your progress is saved. Improve the route and buy upgrades to earn faster. Return to (0, 0) before starting the next pass.`,
		Starter: `package main

import "orbital/colony"

// The colony keeps running until you press Stop, and everything is saved.
// To start, the drone visits the whole field in a snake pattern, tilling, harvesting, and planting moss.
// Once it has 5 or more moss, it sells them at the market.
// Improve this code and buy upgrades to earn faster.

func visit() {
	colony.Till()
	colony.Harvest()
	colony.Plant(colony.Moss)
	if colony.Count(colony.Moss) >= 5 {
		colony.Export(colony.Moss, 0)
	}
}

func main() {
	for {
		w, h := colony.WorldSize()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				visit()
				if x < w-1 {
					if y%2 == 0 {
						colony.Move(colony.East)
					} else {
						colony.Move(colony.West)
					}
				}
			}
			if y < h-1 {
				colony.Move(colony.North)
			}
		}
		// Return to (0, 0) for the next pass
		x, y := colony.Position()
		for ; y > 0; y-- {
			colony.Move(colony.South)
		}
		for ; x > 0; x-- {
			colony.Move(colony.West)
		}
	}
}
`,
		Hints: []string{},
	},
	ID: "colony", Kind: "colony", Title: "المستعمرة",
	Goal: "كبّر المستعمرة: ازرع، احصد، بيع، واشتري ترقيات",
	Setup: func(w *sim.World) {
		w.Init(3, 3, sim.Regolith)
		w.Check = func(*sim.World) (bool, string) { return false, i18n.T("المستعمرة", "The Colony") }
	},
	Starter: `package main

import "orbital/colony"

// المستعمرة شغالة على طول لحد ما تدوس Stop، وكل حاجة بتتحفظ.
// البداية: الدرون بيلف على الأرض كلها بالثعبان، يحرث ويحصد ويزرع طحالب،
// ولما يتجمع 5 طحالب أو أكتر يبيعهم في السوق.
// حسّن الكود ده واشتري ترقيات عشان تكسب أسرع.

func visit() {
	colony.Till()
	colony.Harvest()
	colony.Plant(colony.Moss)
	if colony.Count(colony.Moss) >= 5 {
		colony.Export(colony.Moss, 0)
	}
}

func main() {
	for {
		w, h := colony.WorldSize()
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				visit()
				if x < w-1 {
					if y%2 == 0 {
						colony.Move(colony.East)
					} else {
						colony.Move(colony.West)
					}
				}
			}
			if y < h-1 {
				colony.Move(colony.North)
			}
		}
		// ارجع لـ (0, 0) عشان اللفة الجاية
		x, y := colony.Position()
		for ; y > 0; y-- {
			colony.Move(colony.South)
		}
		for ; x > 0; x-- {
			colony.Move(colony.West)
		}
	}
}
`,
}
