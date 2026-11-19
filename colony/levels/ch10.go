package levels

import (
	"fmt"
	"sort"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func crystalRows(w *sim.World, rows ...int) []sim.Point {
	var ps []sim.Point
	for _, y := range rows {
		for x := 0; x < w.W; x++ {
			t := w.At(x, y)
			t.Entity, t.Size = sim.Crystal, 1
			ps = append(ps, sim.Point{X: x, Y: y})
		}
	}
	return ps
}

const spawnHint = "لو main خلصت قبل الدرونات التانية، البرنامج كله بيقفل وهما معاه. استنّاهم بـ sync.WaitGroup."

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch10LessonEnglish["10-1-second-drone"],
			Title: "A Second Drone",
			Goal:  "Harvest both rows at once: one for you and one for a new drone",
			Lesson: `# Goroutines: Several Things at Once 🚁🚁

From here on, time is **real**: each move takes actual time. One drone harvests two rows in about 40 ticks. Two drones working **at once** finish in about 21.

## go means "start in the background"

~~~go
go harvestRow(helper, colony.East) // Starts and returns immediately
harvestRow(colony.Main(), colony.East) // Runs at the same time
~~~

Putting ~go~ before a function call runs it in a **goroutine**: a very lightweight execution thread. Go can run millions of them.

## The problem: main does not wait 😱
When ~main~ returns, **the whole program exits**, even if goroutines are still running. Try it and see.

## sync.WaitGroup: a work counter

~~~go
var wg sync.WaitGroup

wg.Add(1)              // One new task
go func() {
	defer wg.Done()    // Decrement the counter when finished
	harvestRow(helper, colony.East)
}()

wg.Wait()              // Wait until the counter reaches zero
~~~

- ~func() { ... }()~ is an **anonymous function** called immediately (a closure).
- ~defer~ means "run this when the function exits," including when it panics.

## The drones
- ~colony.Spawn()~ returns a new ~*colony.Drone~ at (0, 0), or nil if the hangar is empty.
- ~colony.Main()~ is your original drone as a ~*colony.Drone~.
- Each drone has the same commands: ~d.Move~, ~d.Harvest~, ~d.Measure~...`,
			Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func harvestRow(d *colony.Drone, dir colony.Direction) {
	for x := 0; x < 10; x++ {
		d.Harvest()
		if x < 9 {
			d.Move(dir)
		}
	}
}

func main() {
	helper := colony.Spawn() // A new drone at (0, 0)
	var wg sync.WaitGroup

	// Run helper in a goroutine: move North twice and harvest row y = 2
	_ = helper
	_ = &wg

	harvestRow(colony.Main(), colony.East) // You harvest row y = 0
}
`,
			Hints: []string{
				"Give each drone separate work.",
				"Start the helper and wait in main.",
				"Check ~_ = helper~ and the end of ~main~.",
			},
		},
		LessonRevision: 1, Steps: ch10LessonArabic["10-1-second-drone"],
		ID: "10-1-second-drone", Chapter: 10, Title: "درون تاني",
		Goal: "احصد الصفين في نفس الوقت: إنت صف، ودرون جديد الصف التاني",
		Par:  32, Gold: 25, MaxTicks: 200,
		Requires: []string{"go", "waitgroup"},
		Setup: func(w *sim.World) {
			w.Init(10, 3, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 2
			ps := crystalRows(w, 0, 2)
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, ps)
				msg := fmt.Sprintf(i18n.T("كريستال فاضل: %d / %d", "Crystal remaining: %d / %d"), left, len(ps))
				if len(w.Drones) > 1 && left > 0 {
					msg += " · " + i18n.T(spawnHint, "If main finishes before the other drones, the whole program exits with them. Wait for them with sync.WaitGroup.")
				}
				return left == 0, msg
			}
		},
		Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func harvestRow(d *colony.Drone, dir colony.Direction) {
	for x := 0; x < 10; x++ {
		d.Harvest()
		if x < 9 {
			d.Move(dir)
		}
	}
}

func main() {
	helper := colony.Spawn() // درون جديد عند (0, 0)
	var wg sync.WaitGroup

	// شغّل الـ helper في goroutine: يطلع North مرتين ويحصد صف y = 2
	_ = helper
	_ = &wg

	harvestRow(colony.Main(), colony.East) // إنت بتحصد صف y = 0
}
`,
		Hints: []string{
			"كل درون له شغل منفصل.",
			"شغّل helper وخلي main تستناه.",
			"راجع ~_ = helper~ وآخر ~main~.",
		},
		Lesson: `# Goroutines: أكتر من حاجة في نفس الوقت 🚁🚁

من هنا الوقت بقى **حقيقي**: كل حركة بتاخد وقت فعلاً. درون واحد بيحصد صفين في ~40 tick. درونين **في نفس الوقت** يخلّصوا في ~21.

## go = "اشتغل في الخلفية"

~~~go
go harvestRow(helper, colony.East) // بيبدأ وبيرجع فوراً
harvestRow(colony.Main(), colony.East) // ده بيشتغل في نفس الوقت
~~~

كلمة ~go~ قبل استدعاء دالة بتشغّلها في **goroutine**: خيط تنفيذ خفيف جداً. Go يقدر يشغّل ملايين منهم.

## المشكلة: main مبتستناش 😱
لما ~main~ تخلص، **البرنامج كله بيقفل**، حتى لو فيه goroutines لسه شغالة. جرّب وشوف بنفسك.

## sync.WaitGroup: عدّاد للشغل

~~~go
var wg sync.WaitGroup

wg.Add(1)              // فيه شغلانة واحدة جديدة
go func() {
	defer wg.Done()    // لما تخلص، نقّص العدّاد
	harvestRow(helper, colony.East)
}()

wg.Wait()              // استنى لحد ما العدّاد يبقى صفر
~~~

- ~func() { ... }()~ دي **دالة من غير اسم** بتتنادى في نفس اللحظة (closure).
- ~defer~ معناها "نفّذ السطر ده لما الدالة تخلص"، حتى لو حصل panic.

## الدرونات
- ~colony.Spawn()~ بيرجّع ~*colony.Drone~ جديد عند (0, 0)، أو nil لو الهانجر فاضي.
- ~colony.Main()~ هو الدرون الأصلي بتاعك كـ ~*colony.Drone~.
- كل درون عنده نفس الأوامر: ~d.Move~ و~d.Harvest~ و~d.Measure~…`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch10LessonEnglish["10-2-swarm"],
			Title: "The Swarm",
			Goal:  "4 drones, 4 rows. Finish the whole field at once",
			Lesson: `# Goroutines in a Loop 🐝

~~~go
for row := 1; row < 4; row++ {
	d := colony.Spawn()
	wg.Add(1)
	go func() {
		defer wg.Done()
		work(d, row)
	}()
}
wg.Wait()
~~~

## What does the closure see?
The goroutine uses ~d~ and ~row~ from the loop. Since Go 1.22, each iteration has a **new copy** of the loop variable, so each goroutine gets the right number.

(In older Go versions this was a famous bug: all goroutines could see the final value. If you see ~row := row~ in old code, that was the workaround.)

## Pass it as a parameter (clearer)

~~~go
go func(d *colony.Drone, row int) {
	defer wg.Done()
	work(d, row)
}(d, row)
~~~

## Order is not guaranteed
Goroutines run **concurrently**. Which prints first? You cannot rely on an order. If order matters, you need **synchronization**, the subject of the next chapter.

## Go 1.25+: wg.Go
~~~go
wg.Go(func() { work(d, row) }) // Add + go + Done in one line
~~~`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
	"sync"
)

func work(d *colony.Drone, row int) {
	fmt.Println("Drone", d.ID(), "heading to row", row)
	// Move up to row and harvest the whole row
}

func main() {
	var wg sync.WaitGroup
	for row := 1; row < 4; row++ {
		d := colony.Spawn()
		// Run work(d, row) in a goroutine
		_ = d
	}
	work(colony.Main(), 0)
	wg.Wait()
}
`,
			Hints: []string{
				"One row is one task.",
				"Set up each task inside the loop.",
				"Check ~_ = d~ and ~work~.",
			},
		},
		LessonRevision: 1, Steps: ch10LessonArabic["10-2-swarm"],
		ID: "10-2-swarm", Chapter: 10, Title: "السرب",
		Goal: "4 درونات، 4 صفوف. خلّصوا الحقل كله في نفس الوقت",
		Par:  30, Gold: 22, MaxTicks: 300,
		Requires: []string{"go", "waitgroup", "for"},
		Setup: func(w *sim.World) {
			w.Init(8, 4, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 4
			ps := crystalRows(w, 0, 1, 2, 3)
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, ps)
				msg := fmt.Sprintf(i18n.T("كريستال فاضل: %d / %d · درونات: %d", "Crystal remaining: %d / %d · Drones: %d"), left, len(ps), len(w.Drones))
				return left == 0, msg
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
	"sync"
)

func work(d *colony.Drone, row int) {
	fmt.Println("درون", d.ID(), "رايح للصف", row)
	// اطلع للصف row، واحصده كله
}

func main() {
	var wg sync.WaitGroup
	for row := 1; row < 4; row++ {
		d := colony.Spawn()
		// شغّل work(d, row) في goroutine
		_ = d
	}
	work(colony.Main(), 0)
	wg.Wait()
}
`,
		Hints: []string{
			"كل صف مهمة لوحدها.",
			"جهّز كل مهمة جوه الـloop.",
			"راجع ~_ = d~ و~work~.",
		},
		Lesson: `# goroutines في حلقة 🐝

~~~go
for row := 1; row < 4; row++ {
	d := colony.Spawn()
	wg.Add(1)
	go func() {
		defer wg.Done()
		work(d, row)
	}()
}
wg.Wait()
~~~

## الـ closure بتشوف إيه؟
الـ goroutine بتستخدم ~d~ و~row~ من الحلقة. من Go 1.22، كل لفّة ليها **نسخة جديدة** من متغير الحلقة، فكل goroutine بتاخد الرقم الصح.

(في Go القديم كان ده bug مشهور جداً: كل الـ goroutines كانت بتشوف آخر قيمة. لو شفت كود قديم فيه ~row := row~، ده كان الحل.)

## ابعتها كـ parameter (أوضح)

~~~go
go func(d *colony.Drone, row int) {
	defer wg.Done()
	work(d, row)
}(d, row)
~~~

## الترتيب مش مضمون
الـ goroutines بتشتغل **في نفس الوقت**. مين يطبع الأول؟ محدش يعرف. لو الترتيب مهم، لازم **تزامن** (synchronization)، وده موضوع الفصل الجاي.

## Go 1.25+: wg.Go
~~~go
wg.Go(func() { work(d, row) }) // Add + go + Done في سطر واحد
~~~`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch10LessonEnglish["10-3-divide"],
			Title: "Divide the Work",
			Goal:  "3 drones and 15 scattered targets. Divide the targets wisely",
			Lesson: `# Dividing the Work ✂️

You have 15 targets and 3 drones. Total time is the time of the **slowest** drone. A good split **balances** the work.

## Give each goroutine a slice

~~~go
parts := [][]colony.Point{targets[:5], targets[5:10], targets[10:]}
for i, d := range drones {
	wg.Add(1)
	go func() {
		defer wg.Done()
		collect(d, parts[i])
	}()
}
wg.Wait()
~~~

## [][]T: a slice of slices
~[][]colony.Point~ means a list of lists. Each drone gets its own list.

## Sorting with a comparison function

~~~go
import "slices"

slices.SortFunc(targets, func(a, b colony.Point) int {
	return a.X - b.X // Negative if a comes before b
})
~~~

## An important concept: load balancing ⚖️
- If one drone gets too many distant targets, the others finish and wait for it.
- Splitting by **region** is often better than splitting by **list order**.
- Next chapter introduces an automatic approach: each drone takes the next job as soon as it is free (a worker pool).`,
			Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

func collect(d *colony.Drone, targets []colony.Point) {
	for _, t := range targets {
		goTo(d, t)
		d.Harvest()
	}
}

func main() {
	targets := colony.Targets()
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	var wg sync.WaitGroup

	// Divide targets among the 3 drones, each working in a goroutine
	_ = drones
	_ = &wg
	collect(colony.Main(), targets)
}
`,
			Hints: []string{
				"Balance routes, not just counts.",
				"Group targets before visiting them.",
				"Start at ~colony.Targets()~ in ~main~.",
			},
		},
		LessonRevision: 1, Steps: ch10LessonArabic["10-3-divide"],
		ID: "10-3-divide", Chapter: 10, Title: "قسّم الشغل",
		Goal: "3 درونات و15 هدف متفرقين. قسّم الأهداف عليهم بذكاء",
		Par:  60, Gold: 40, MaxTicks: 400,
		Requires: []string{"go", "waitgroup"},
		Setup: func(w *sim.World) {
			w.Init(10, 8, sim.Regolith)
			w.Realtime, w.MaxDrones = true, 3
			w.Targets = randomPoints(w, 15)
			placeCrystals(w, w.Targets)
			byX := append([]sim.Point(nil), w.Targets...)
			sort.Slice(byX, func(i, j int) bool { return byX[i].X < byX[j].X })
			worst := 0
			for g := 0; g < 3; g++ {
				worst = max(worst, greedyCost(byX[g*5:g*5+5])+5)
			}
			w.Gold = worst + 2
			w.Par = w.Gold + 14
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, w.Targets)
				return left == 0, fmt.Sprintf(i18n.T("أهداف فاضلة: %d / %d", "Targets remaining: %d / %d"), left, len(w.Targets))
			}
		},
		Starter: `package main

import (
	"orbital/colony"
	"sync"
)

func goTo(d *colony.Drone, p colony.Point) {
	x, y := d.Position()
	for ; x < p.X; x++ {
		d.Move(colony.East)
	}
	for ; x > p.X; x-- {
		d.Move(colony.West)
	}
	for ; y < p.Y; y++ {
		d.Move(colony.North)
	}
	for ; y > p.Y; y-- {
		d.Move(colony.South)
	}
}

func collect(d *colony.Drone, targets []colony.Point) {
	for _, t := range targets {
		goTo(d, t)
		d.Harvest()
	}
}

func main() {
	targets := colony.Targets()
	drones := []*colony.Drone{colony.Main(), colony.Spawn(), colony.Spawn()}
	var wg sync.WaitGroup

	// قسّم targets على الـ 3 درونات، وكل واحد يشتغل في goroutine
	_ = drones
	_ = &wg
	collect(colony.Main(), targets)
}
`,
		Hints: []string{
			"وازن المسافة مع العدد.",
			"قسّم قبل زيارة الأهداف.",
			"ابدأ عند ~colony.Targets()~ في ~main~.",
		},
		Lesson: `# تقسيم الشغل ✂️

عندك 15 هدف و3 درونات. الوقت الكلي = وقت **أبطأ** درون. يعني التقسيمة الحلوة هي اللي **بتوازن** الشغل.

## ابعت slice لكل goroutine

~~~go
parts := [][]colony.Point{targets[:5], targets[5:10], targets[10:]}
for i, d := range drones {
	wg.Add(1)
	go func() {
		defer wg.Done()
		collect(d, parts[i])
	}()
}
wg.Wait()
~~~

## [][]T: slice من slices
~[][]colony.Point~ معناها قايمة فيها قوايم. كل درون بياخد قايمته.

## ترتيب بـ دالة مقارنة

~~~go
import "slices"

slices.SortFunc(targets, func(a, b colony.Point) int {
	return a.X - b.X // سالب لو a قبل b
})
~~~

## مفهوم مهم: load balancing ⚖️
- لو درون واحد خد أهداف بعيدة وكتير، التانيين هيخلصوا ويقفوا مستنيينه.
- تقسيم حسب **المنطقة** غالباً أحسن من تقسيم حسب **الترتيب**.
- في الفصل الجاي هتشوف طريقة تلقائية: كل درون ياخد الشغلانة الجاية أول ما يفضى (worker pool).`,
	})
}
