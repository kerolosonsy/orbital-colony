package levels

import (
	"fmt"
	"strconv"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1,
			Steps:          ch12Lesson121English,
			Title:          "The Shared Counter",
			Goal:           "4 drones add to one variable. Report the correct total... without a data race",
			Lesson: `# Data Race: Everyone Writes at Once 🏁

This mission runs with the **race detector** (~go build -race~), a real Go tool that catches dangerous concurrency bugs.

## The problem
~total += n~ looks like one step, but it is actually 3:
1. Read ~total~
2. Add ~n~
3. Write ~total~

If one goroutine reads ~total = 10~ and another reads the same 10 before the first writes, both write 10 plus their own amount, and **an update is lost**. This is a **data race**.

Worse, this bug appears **intermittently**: it might work 1000 times, then fail in production.

## sync.Mutex: a lock

~~~go
var mu sync.Mutex

mu.Lock()     // Wait if someone else holds the lock
total += n    // Only one goroutine here at a time
mu.Unlock()
~~~

The code between Lock and Unlock is the **critical section**. Keep it as **short** as possible.

## defer with Unlock
~~~go
mu.Lock()
defer mu.Unlock() // Unlock even if the function panics
~~~

## An alternative: sync/atomic
For simple counters only:
~~~go
var total atomic.Int64
total.Add(int64(n))
~~~

## ⚠️ Never copy a mutex
A copied mutex is a completely different lock. Always pass a pointer (~*sync.Mutex~), or put it in a struct passed by pointer. ~go vet~ catches this mistake.`,
			Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func main() {
	total := 0
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn(), colony.Spawn()}

	for row, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < row; i++ {
				d.Move(colony.North)
			}
			for x := 0; x < 8; x++ {
				n, _ := d.TryHarvest()
				total += n // ⚠️ 4 goroutines write here at once
				if x < 7 {
					d.Move(colony.East)
				}
			}
		}()
	}
	wg.Wait()
	colony.Report(total)
}
`,
			Hints: []string{
				"A lucky correct result does not prove shared writes are safe.",
				"Separate the total update from movement, harvesting, and waiting.",
				"Inspect the ~total += n~ line inside the goroutine.",
			},
		},
		LessonRevision: 1, Steps: ch12Lesson121Arabic,
		ID: "12-1-shared-total", Chapter: 12, Title: "العدّاد المشترك",
		Goal: "4 درونات بيجمعوا في نفس المتغير. ابعت المجموع الصح… من غير data race",
		Par:  25, Gold: 20, MaxTicks: 300,
		Requires: []string{"mutex"},
		Race:     true,
		Setup: func(w *sim.World) {
			w.Init(8, 4, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 4
			sum := 0
			w.Each(func(_, _ int, t *sim.Tile) {
				t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
				sum += t.Size
			})
			want := strconv.Itoa(sum)
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("المجموع صح", "The total is correct")
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[n-1] + i18n.T(" والمطلوب ", "; expected ") + want
				}
				return false, fmt.Sprintf(i18n.T("كريستال في المخزن: %d", "Crystal in storage: %d"), w.Inv[sim.Crystal])
			}
		},
		Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func main() {
	total := 0
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn(), colony.Spawn()}

	for row, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < row; i++ {
				d.Move(colony.North)
			}
			for x := 0; x < 8; x++ {
				n, _ := d.TryHarvest()
				total += n // ⚠️ 4 goroutines بيكتبوا هنا في نفس الوقت
				if x < 7 {
					d.Move(colony.East)
				}
			}
		}()
	}
	wg.Wait()
	colony.Report(total)
}
`,
		Hints: []string{
			"نتيجة سليمة مرة مش إثبات إن الكتابة المشتركة آمنة.",
			"افصل تحديث الإجمالي عن الحركة والحصاد وانتظار العمال.",
			"راجع ~total += n~ جوه الـgoroutine.",
		},
		Lesson: `# Data race: لما الكل يكتب في نفس الوقت 🏁

المهمة دي شغالة بالـ **race detector** (~go build -race~)، وده أداة حقيقية جوه Go بتمسك أخطر نوع bugs في الـ concurrency.

## المشكلة
~total += n~ شكلها خطوة واحدة، بس هي فعلياً 3 خطوات:
1. اقرا ~total~
2. زوّد ~n~
3. اكتب ~total~

لو goroutine قرت ~total = 10~، وواحدة تانية قرت نفس الـ 10 قبل ما الأولى تكتب، الاتنين هيكتبوا 10 + حاجة، وإضافة **هتضيع**. ده اسمه **data race**.

الأوحش إن الـ bug ده بيظهر **أحياناً**: ممكن يشتغل صح 1000 مرة ويبوظ في الـ production.

## sync.Mutex: قفل

~~~go
var mu sync.Mutex

mu.Lock()     // لو حد تاني ماسك القفل، استنى
total += n    // goroutine واحدة بس هنا في أي لحظة
mu.Unlock()
~~~

الكود اللي بين Lock وUnlock اسمه **critical section**. خلّيه **قصير** على قد ما تقدر.

## defer مع Unlock
~~~go
mu.Lock()
defer mu.Unlock() // هيتفتح حتى لو حصل panic
~~~

## بديل: sync/atomic
للعدّادات البسيطة بس:
~~~go
var total atomic.Int64
total.Add(int64(n))
~~~

## ⚠️ متنسخش الـ mutex
mutex اتنسخ = قفل تاني خالص. ابعته دايماً كـ pointer (~*sync.Mutex~)، أو خلّيه جوه struct بتبعت pointer ليه. و~go vet~ بيمسك الغلطة دي.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1,
			Steps:          ch12Lesson122English,
			Title:          "The Safe Registry",
			Goal:           "3 drones measure and record into one map. Report the largest crystal's location",
			Lesson: `# A Mutex Inside a Struct: A Common Go Pattern 🗃️

Go maps are **not safe** for concurrent writes. The program can crash with ~fatal error: concurrent map writes~, or the race detector can report the problem.

## Keep the lock beside the data it protects

~~~go
type Registry struct {
	mu    sync.Mutex
	sizes map[colony.Point]int
}

func (r *Registry) Record(p colony.Point, size int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sizes[p] = size
}
~~~

This is a very common Go pattern:
- The lock lives **inside** the struct, beside its data.
- Access the data only through **methods**, each taking the lock.
- Users of ~Registry~ need not know about the lock. This is **encapsulation**.

## Why a pointer receiver?
You need ~*Registry~ because a value receiver would **copy the mutex**. The copy is a different lock, making the protection ineffective. (Remember chapter 8?)

## RWMutex: many reads, few writes
~~~go
var mu sync.RWMutex
mu.RLock()   // Multiple readers may read at once
mu.RUnlock()
mu.Lock()    // Writing requires exclusive access
mu.Unlock()
~~~`,
			Starter: `package main

import (
	"orbital/colony"
	"sync"
)

// Registry is shared among the drones
type Registry struct {
	sizes map[colony.Point]int
}

func (r *Registry) Record(p colony.Point, size int) {
	r.sizes[p] = size // ⚠️ Unsafe with multiple goroutines
}

func main() {
	reg := &Registry{sizes: map[colony.Point]int{}}
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}

	for row, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < row; i++ {
				d.Move(colony.North)
			}
			for x := 0; x < 10; x++ {
				reg.Record(colony.Point{X: x, Y: row}, d.Measure())
				if x < 9 {
					d.Move(colony.East)
				}
			}
		}()
	}
	wg.Wait()

	// Search reg.sizes and report the largest crystal's location with colony.Report
}
`,
			Hints: []string{
				"Concurrent writers to one map need one shared mutex.",
				"Let Registry own the map and lock; separate recording from reading.",
				"Inspect ~Record~ and the ~r.sizes[p] = size~ line.",
			},
		},
		LessonRevision: 1, Steps: ch12Lesson122Arabic,
		ID: "12-2-registry", Chapter: 12, Title: "السجل الآمن",
		Goal: "3 درونات بيقيسوا ويسجلوا في نفس الـmap. ابعت مكان أكبر كريستال",
		Par:  30, Gold: 24, MaxTicks: 300,
		Requires: []string{"mutex", "method"},
		Race:     true,
		Setup: func(w *sim.World) {
			w.Init(10, 3, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 3
			perm := w.Rand.Perm(30)
			best, bestP := -1, sim.Point{}
			w.Each(func(x, y int, t *sim.Tile) {
				t.Entity, t.Size = sim.Crystal, perm[y*10+x]+1
				if t.Size > best {
					best, bestP = t.Size, sim.Point{X: x, Y: y}
				}
			})
			want := fmt.Sprint(bestP)
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("لقيت أكبر كريستال", "Largest crystal found")
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[n-1] + i18n.T(" — ده مش مكان أكبر كريستال", " — that is not the largest crystal's location")
				}
				return false, i18n.T("القاعدة مستنية مكان أكبر كريستال (colony.Point)", "Base is waiting for the largest crystal's location (colony.Point)")
			}
		},
		Starter: `package main

import (
	"orbital/colony"
	"sync"
)

// Registry سجل مشترك بين الدرونات
type Registry struct {
	sizes map[colony.Point]int
}

func (r *Registry) Record(p colony.Point, size int) {
	r.sizes[p] = size // ⚠️ مش آمن مع أكتر من goroutine
}

func main() {
	reg := &Registry{sizes: map[colony.Point]int{}}
	var wg sync.WaitGroup
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}

	for row, d := range drones {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < row; i++ {
				d.Move(colony.North)
			}
			for x := 0; x < 10; x++ {
				reg.Record(colony.Point{X: x, Y: row}, d.Measure())
				if x < 9 {
					d.Move(colony.East)
				}
			}
		}()
	}
	wg.Wait()

	// لف على reg.sizes وابعت مكان أكبر واحد بـ colony.Report
}
`,
		Hints: []string{
			"الكتّاب المتزامنين على نفس الـmap محتاجين الـmutex المشتركة.",
			"خلّي Registry تملك الـmap والـmutex، وافصل التسجيل عن القراية.",
			"راجع ~Record~ وسطر ~r.sizes[p] = size~.",
		},
		Lesson: `# Mutex جوه struct: النمط الاحترافي 🗃️

الـ maps في Go **مش آمنة** لو أكتر من goroutine بيكتبوا فيها في نفس الوقت. ممكن البرنامج يقع برسالة ~fatal error: concurrent map writes~، أو الـ race detector يمسكها.

## خلّي القفل جنب البيانات اللي بيحميها

~~~go
type Registry struct {
	mu    sync.Mutex
	sizes map[colony.Point]int
}

func (r *Registry) Record(p colony.Point, size int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sizes[p] = size
}
~~~

ده أشهر نمط في Go:
- القفل **جوه** الـ struct، جنب البيانات بالظبط.
- البيانات مش بيوصلها حد غير من خلال **methods**، وكل method بتقفل.
- اللي بيستخدم ~Registry~ ميعرفش أصلاً إن فيه قفل. ده اسمه **encapsulation**.

## ليه pointer receiver؟
لازم ~*Registry~، لأن value receiver هيعمل **نسخة من الـ mutex**، والنسخة دي قفل مختلف خالص، فالحماية تبقى وهمية. (فاكر الفصل 8؟)

## RWMutex: قراية كتير، كتابة قليلة
~~~go
var mu sync.RWMutex
mu.RLock()   // كذا واحد يقدروا يقروا في نفس الوقت
mu.RUnlock()
mu.Lock()    // الكتابة لوحدها
mu.Unlock()
~~~`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1,
			Steps:          ch12Lesson123English,
			Title:          "Deadlock",
			Goal:           "The drones have stopped! Find and fix the deadlock, then harvest both rows",
			Lesson: `# Deadlock: Each Waits for the Other 🔒🔒

Imagine two people in a narrow corridor, each waiting for the other to pass first. Neither can move.

~~~text
Drone a: holds drill ✓ ... waiting for cart ⏳
Drone b: holds cart ✓ ... waiting for drill ⏳
~~~

This is **deadlock**: each goroutine holds something and waits for something another holds.

## The four conditions (Coffman)
Deadlock occurs when all four conditions hold together:
1. Only one participant can hold a resource at a time (mutual exclusion).
2. A participant holding one resource waits for another.
3. A lock cannot be forcibly taken away.
4. A **cycle**: a waits for b, and b waits for a.

Break any condition and the deadlock disappears.

## The practical solution: a fixed lock order
If **all** code locks drill before cart, the cycle cannot form.

## How do you detect it?
- Go sometimes reports: ~fatal error: all goroutines are asleep - deadlock!~
- But if another goroutine is active (a timer, for example), the program may **simply hang**. Here Mission Control stops it after 200 ticks.
- In real programs, ~kill -QUIT~ or ~SIGQUIT~ prints a stack trace for every goroutine, showing where each is stuck.

## Better still: avoid holding two locks
If you can design with one lock, or use channels instead, that is better.`,
			Starter: `package main

import (
	"orbital/colony"
	"sync"
)

// The colony has one drill and one cart. Every harvest needs both.
var drill, cart sync.Mutex

func harvestWith(d *colony.Drone, first, second *sync.Mutex, ready *sync.WaitGroup) {
	first.Lock()
	if ready != nil {
		ready.Done()
		ready.Wait()
	}
	d.Wait() // Prepare the equipment
	second.Lock()
	d.Harvest()
	second.Unlock()
	first.Unlock()
}

func work(d *colony.Drone, row int, first, second *sync.Mutex, ready *sync.WaitGroup) {
	for i := 0; i < row; i++ {
		d.Move(colony.North)
	}
	for x := 0; x < 6; x++ {
		harvestWith(d, first, second, ready)
		ready = nil
		if x < 5 {
			d.Move(colony.East)
		}
	}
}

func main() {
	a, b := colony.Main(), colony.Spawn()
	// The first harvest lines up both drones when they start with different locks.
	aFirst, aSecond := &drill, &cart
	bFirst, bSecond := &cart, &drill

	var ready *sync.WaitGroup
	if aFirst != bFirst {
		ready = &sync.WaitGroup{}
		ready.Add(2)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); work(a, 0, aFirst, aSecond, ready) }()
	go func() { defer wg.Done(); work(b, 1, bFirst, bSecond, ready) }()
	wg.Wait()
}
`,
			Hints: []string{
				"A deadlock can form when each worker holds one resource and waits for another.",
				"Compare the two lock orders on every path before tracing Wait.",
				"Compare ~aFirst, aSecond~ with ~bFirst, bSecond~ inside ~main~.",
			},
		},
		LessonRevision: 1, Steps: ch12Lesson123Arabic,
		ID: "12-3-deadlock", Chapter: 12, Title: "الجمود",
		Goal: "الدرونات وقفت ومحدش بيتحرك! لاقي الـdeadlock وصلّحه، واحصد الصفين",
		Par:  36, Gold: 28, MaxTicks: 200,
		Requires: []string{"mutex"},
		Setup: func(w *sim.World) {
			w.Init(6, 2, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 2
			ps := crystalRows(w, 0, 1)
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, ps)
				return left == 0, fmt.Sprintf(i18n.T("كريستال فاضل: %d / %d", "Crystal remaining: %d / %d"), left, len(ps))
			}
		},
		Starter: `package main

import (
	"orbital/colony"
	"sync"
)

// في المستعمرة مثقاب واحد وعربية واحدة. كل حصاد محتاج الاتنين.
var drill, cart sync.Mutex

func harvestWith(d *colony.Drone, first, second *sync.Mutex, ready *sync.WaitGroup) {
	first.Lock()
	if ready != nil {
		ready.Done()
		ready.Wait()
	}
	d.Wait() // بيجهّز العدة
	second.Lock()
	d.Harvest()
	second.Unlock()
	first.Unlock()
}

func work(d *colony.Drone, row int, first, second *sync.Mutex, ready *sync.WaitGroup) {
	for i := 0; i < row; i++ {
		d.Move(colony.North)
	}
	for x := 0; x < 6; x++ {
		harvestWith(d, first, second, ready)
		ready = nil
		if x < 5 {
			d.Move(colony.East)
		}
	}
}

func main() {
	a, b := colony.Main(), colony.Spawn()
	// أول حصاد بينسّق بين الدرونين لما كل واحد يبدأ بـlock مختلف.
	aFirst, aSecond := &drill, &cart
	bFirst, bSecond := &cart, &drill

	var ready *sync.WaitGroup
	if aFirst != bFirst {
		ready = &sync.WaitGroup{}
		ready.Add(2)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); work(a, 0, aFirst, aSecond, ready) }()
	go func() { defer wg.Done(); work(b, 1, bFirst, bSecond, ready) }()
	wg.Wait()
}
`,
		Hints: []string{
			"الـdeadlock ممكن يحصل لما كل عامل ماسك مورد ومستني التاني.",
			"قارن ترتيب الـlocks في كل مسار قبل تتبّع Wait.",
			"قارن ~aFirst, aSecond~ مع ~bFirst, bSecond~ جوه ~main~.",
		},
		Lesson: `# Deadlock: كل واحد مستني التاني 🔒🔒

تخيّل اتنين في ممر ضيق، وكل واحد مستني التاني يعدّي الأول. محدش هيتحرك أبداً.

~~~text
درون a: مسك drill ✓ … مستني cart ⏳
درون b: مسك cart ✓ … مستني drill ⏳
~~~

ده **deadlock**: كل goroutine ماسكة حاجة ومستنية حاجة ماسكها غيرها.

## الشروط الأربعة (Coffman)
الـ deadlock بيحصل لما الأربعة يتحققوا مع بعض:
1. المورد بيتمسك من واحد بس في نفس الوقت (mutex).
2. اللي ماسك حاجة بيستنى حاجة تانية.
3. محدش يقدر ياخد القفل غصب.
4. **دايرة**: a مستني b، وb مستني a.

اكسر أي شرط منهم، والـ deadlock يختفي.

## الحل العملي: ترتيب ثابت للأقفال
لو **كل** الكود بيقفل drill قبل cart، الدايرة مستحيل تحصل.

## إزاي تكتشفه؟
- Go أحياناً بيقولك: ~fatal error: all goroutines are asleep - deadlock!~
- بس لو فيه goroutine تانية شغالة (timer مثلاً)، البرنامج **بيعلّق وخلاص**. هنا الـ Mission Control بيقطع بعد 200 tick.
- في البرامج الحقيقية: ~kill -QUIT~ أو ~SIGQUIT~ بيطبع stack trace لكل goroutine، فتشوف كل واحدة واقفة فين.

## أحسن: متمسكش قفلين أصلاً
لو تقدر تصمّم الكود بقفل واحد، أو بـ channels بدل الأقفال، يبقى أحسن.`,
	})
}
