package levels

import (
	"fmt"
	"sort"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 2, Steps: surveyLessonEnglish,
			Title: "Crystal Log",
			Goal:  "Measure every crystal and report their sizes in order in one slice",
			Lesson: `# Mission

1. Measure the 8 crystals from left to right.
2. Keep their sizes in order in one slice.
3. Report the whole slice once with ~colony.Report~.

The starter already saved the first measurement.

Gold: **15 ticks or fewer**.

8 measures + 7 moves = 15 ticks; Report=0, assuming no other tick costs.

## Idea

A slice is a box: ~append~ adds each size at its end.

Illustrative sizes:

~~~go
sizes := []int{}
sizes = append(sizes, 5)
sizes = append(sizes, 2)
sizes = append(sizes, 8)
fmt.Println(sizes, len(sizes), sizes[0])
~~~

## Watch out

- Keep append's result: ~sizes = append(...)~.
- Indexes: 0 through ~len(sizes)-1~; outside means panic.

Use a loop for repeated work.

For each crystal: measure, append, move toward the next.

Report once afterward; the starter saved the first reading.

Write your first attempt, then run it.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	var sizes []int // An empty slice

	sizes = append(sizes, colony.Measure())
	fmt.Println(sizes, len(sizes))

	// Continue to the end of the row, then send all of sizes with colony.Report
}
`,
			Hints: []string{"append grows the slice while keeping existing sizes.", "Loop over the crystals; report once afterward, remembering the first saved reading.", "Check the ~append~ / ~colony.Measure()~ line and move guard: save once; move toward another crystal."},
		},
		LessonRevision: 2, Steps: surveyLessonArabic,
		ID: "6-1-survey", Chapter: 6, Title: "سجل الكريستال",
		Goal: "قيس كل الكريستال وابعت أحجامهم بالترتيب في slice واحدة",
		Par:  15, Gold: 15, MaxTicks: 60,
		Requires: []string{"slice"},
		Setup: func(w *sim.World) {
			w.Init(8, 2, sim.Regolith)
			sizes := make([]int, 8)
			for x := range sizes {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
				sizes[x] = t.Size
			}
			want := fmt.Sprint(sizes)
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("السجل وصل صح", "The log is correct")
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[n-1] + i18n.T(" — مش مطابق", " — does not match")
				}
				return false, i18n.T("القاعدة مستنية slice فيها 8 أحجام", "Base is waiting for a slice of 8 sizes")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	var sizes []int // slice فاضية

	sizes = append(sizes, colony.Measure())
	fmt.Println(sizes, len(sizes))

	// كمّل لحد آخر الصف، وبعدين ابعت sizes كلها بـ colony.Report
}
`,
		Hints: []string{"الـslice سجل بيكبر مع append، مع الاحتفاظ بالأحجام اللي اتسجلت قبله.", "رتّب loop للكريستالات وابعت تقرير واحد بعده؛ افتكر أول قراءة في الـstarter.", "راجع سطر ~append~ اللي فيه ~colony.Measure()~ وشرط الحركة: بتسجل مرة وتتحرك لكريستالة تالية؟"},
		Lesson: `# المهمة

1. قيس الـ8 كريستالات من الشمال لليمين.
2. احفظ أحجامهم بالترتيب في slice واحدة.
3. ابعت الـslice كلها مرة واحدة بـ~colony.Report~.

الـstarter قاس أول كريستالة وحفظها فعلًا.

الذهب: **15 ticks أو أقل**.

8 قياسات + 7 حركات = 15 ticks؛ Report=0 وبافتراض مفيش تكلفة تانية.

## الفكرة

تخيل الـslice صندوق، وكل حجم جديد بيتحط في آخره بـ~append~.

أحجام للتوضيح:

~~~go
sizes := []int{}
sizes = append(sizes, 5)
sizes = append(sizes, 2)
sizes = append(sizes, 8)
fmt.Println(sizes, len(sizes), sizes[0])
~~~

## خد بالك

- احتفظ بنتيجة ~append~: ~sizes = append(...)~.
- أرقام الـindex من 0 لحد ~len(sizes)-1~؛ براها يحصل panic.

التكرار محتاج loop.

لكل كريستالة: قيس، ضيف، واتجه للي بعدها.

التقرير مرة واحدة بعد الـloop؛ أول قراءة موجودة فعلًا.

اكتب أول محاولة بنفسك، شغّلها، وراجع السجل اللي طلع.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: topThreeLessonEnglish,
			Title: "The Top 3",
			Goal:  "The drone can carry only 3 crystals. Collect the largest 3",
			Lesson: `# range and the slices Package 🏆

The drone can carry **only 3 loads**. Harvest the wrong crystal and you lose your chance. **Measure everything first**, then **decide**.

## range: visit every element

~~~go
for i, s := range sizes {
	fmt.Println("Position", i, "size", s)
}
~~~

~range~ gives you the index and value. If you do not need the index, replace it with ~_~:

~~~go
total := 0
for _, s := range sizes {
	total += s
}
~~~

## The slices package (part of Go)

~~~go
import "slices"

slices.Sort(nums)              // Sort in ascending order (changes nums itself)
c := slices.Clone(nums)        // An independent copy
m := slices.Max(nums)          // The largest value
i := slices.Index(nums, 15)    // The index of 15, or -1
~~~

## An important trap ⚠️
Assigning a slice to another variable does not copy its elements:

~~~go
a := []int{3, 1, 2}
b := a
slices.Sort(b)
fmt.Println(a) // [1 2 3]  ← a changed too!
~~~

Both refer to **the same underlying memory**. If you need a real copy, use ~slices.Clone~.

## Challenge ⭐⭐⭐
The three-star target is recalculated each run from the crystal positions.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	sizes := make([]int, 10)
	for x := 0; x < 10; x++ {
		sizes[x] = colony.Measure()
		if x < 9 {
			colony.Move(colony.East)
		}
	}
	fmt.Println("Sizes:", sizes)

	// The drone is now at x = 9. Find the largest 3 and return to harvest them
}
`,
			Hints: []string{
				"Size is value; index is position.",
				"Compare all sizes, then plan visits.",
				"Review ~slices.Clone~ and ~slices.Sort~; which order must stay linked to positions?",
			},
		},
		LessonRevision: 1, Steps: topThreeLessonArabic,
		ID: "6-2-top-three", Chapter: 6, Title: "أغلى 3",
		Goal: "الدرون يقدر يشيل 3 كريستالات بس. هات أكبر 3",
		Par:  45, Gold: 31, MaxTicks: 200,
		Requires: []string{"range"},
		Setup: func(w *sim.World) {
			w.Init(10, 2, sim.Regolith)
			w.HarvestLimit = 3
			perm := w.Rand.Perm(20)
			type c struct{ x, size int }
			cs := make([]c, 10)
			for x := 0; x < 10; x++ {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, perm[x]+1
				cs[x] = c{x, t.Size}
			}
			sort.Slice(cs, func(i, j int) bool { return cs[i].size > cs[j].size })
			top, minX := 0, 9
			for _, k := range cs[:3] {
				top += k.size
				minX = min(minX, k.x)
			}
			w.Gold = 19 + (9 - minX) + 3
			w.Par = w.Gold + 14
			w.FailCheck = func(w *sim.World) string {
				if w.Harvests >= 3 && w.Inv[sim.Crystal] != top {
					return fmt.Sprintf(i18n.T("جمعت %d بس. أكبر 3 مجموعهم %d. المهمة فشلت.", "You collected only %d. The largest 3 total %d. Mission failed."), w.Inv[sim.Crystal], top)
				}
				return ""
			}
			w.Check = func(w *sim.World) (bool, string) {
				return w.Inv[sim.Crystal] == top, fmt.Sprintf(i18n.T("شحنات: %d / 3", "Loads: %d / 3"), w.Harvests)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	sizes := make([]int, 10)
	for x := 0; x < 10; x++ {
		sizes[x] = colony.Measure()
		if x < 9 {
			colony.Move(colony.East)
		}
	}
	fmt.Println("الأحجام:", sizes)

	// دلوقتي الدرون عند x = 9. حدّد أكبر 3 ورجّع احصدهم
}
`,
		Hints: []string{
			"الحجم قيمة، والـindex مكان.",
			"قارن الأحجام، ثم خطط الزيارة والحصاد.",
			"راجع ~slices.Clone~ و~slices.Sort~؛ أي ترتيب يفضل مربوط بالمواقع؟",
		},
		Lesson: `# range ومكتبة slices 🏆

الدرون يقدر يشيل **3 شحنات بس**. لو حصدت كريستال غلط، فرصتك راحت. يعني لازم **تقيس الكل الأول**، وبعدين **تقرر**.

## range: لف على كل العناصر

~~~go
for i, s := range sizes {
	fmt.Println("المكان", i, "الحجم", s)
}
~~~

~range~ بتديك الـ index والقيمة. لو مش محتاج الـ index، اكتب ~_~ مكانه:

~~~go
total := 0
for _, s := range sizes {
	total += s
}
~~~

## مكتبة slices (جوه Go نفسها)

~~~go
import "slices"

slices.Sort(nums)              // ترتيب تصاعدي (بيغيّر nums نفسها)
c := slices.Clone(nums)        // نسخة مستقلة
m := slices.Max(nums)          // أكبر قيمة
i := slices.Index(nums, 15)    // مكان 15 أو -1
~~~

## فخ مهم ⚠️
الـ slice مش بتتنسخ لما تديها لمتغير تاني:

~~~go
a := []int{3, 1, 2}
b := a
slices.Sort(b)
fmt.Println(a) // [1 2 3]  ← a اتغيرت كمان!
~~~

الاتنين بيبصّوا على **نفس الذاكرة**. لو محتاج نسخة فعلاً، استخدم ~slices.Clone~.

## تحدّي ⭐⭐⭐
رقم النجوم التلاتة بيتحسب كل مرة على حسب مكان الكريستال.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: censusLessonEnglish,
			Title: "Life Census",
			Goal:  "Create and report a map counting each type (moss, lichen, and crystal)",
			Lesson: `# Maps: The Colony Dictionary 🗂️

A slice looks up elements by number (0, 1, 2...). A **map** looks up values by a **key**; you choose the key type.

~~~go
counts := map[string]int{}   // string keys, int values
counts["moss"] = 3
counts["moss"]++             // 4
fmt.Println(counts["lichen"]) // 0 ← A missing key returns the zero value
~~~

## Check whether a key exists

~~~go
n, ok := counts["crystal"]
if !ok {
	fmt.Println("No crystal at all")
}
~~~

This is called the **comma-ok idiom**, and you will see it often in Go.

## Other operations

~~~go
delete(counts, "moss")       // Delete a key
len(counts)                  // The number of keys
for k, v := range counts {   // Visit every entry
	fmt.Println(k, v)
}
~~~

## An important detail 🎲
When you range over a map, **the order is unspecified** and can vary each time. Go does this deliberately so you do not depend on an order it cannot guarantee. (~fmt.Println~ prints maps in sorted order for readability.)

## Initialize a map before writing to it
~~~go
var m map[string]int // nil map
m["x"] = 1           // panic!
m = map[string]int{} // Or make(map[string]int)
~~~

## Your mission
Visit every tile, count each type in a map, and report it to base.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	counts := map[colony.Entity]int{}

	e := colony.GetEntity()
	counts[e]++
	fmt.Println(counts)

	// Visit the whole map (6×4), count, then report counts
}
`,
			Hints: []string{
				"A map counts each kind; len counts keys.",
				"Start from the starter's count; visit other tiles once.",
				"Check the map and ~colony.GetEntity()~ line: do key types match?",
			},
		},
		LessonRevision: 1, Steps: censusLessonArabic,
		ID: "6-3-census", Chapter: 6, Title: "تعداد الحياة",
		Goal: "اعمل map فيها عدد كل نوع (moss و lichen و crystal) وابعتها",
		Par:  60, Gold: 47, MaxTicks: 300,
		Requires: []string{"map"},
		Setup: func(w *sim.World) {
			w.Init(6, 4, sim.Regolith)
			kinds := []sim.Entity{sim.Moss, sim.Lichen, sim.Crystal, sim.Empty}
			w.Each(func(x, y int, t *sim.Tile) {
				e := kinds[w.Rand.Intn(len(kinds))]
				switch {
				case x == 1 && y == 0:
					e = sim.Moss
				case x == 2 && y == 1:
					e = sim.Lichen
				case x == 3 && y == 2:
					e = sim.Crystal
				}
				switch e {
				case sim.Moss:
					t.Ground = sim.Soil
					w.PlantAt(x, y, e, 0)
				case sim.Lichen:
					t.Ground = sim.Ice
					w.PlantAt(x, y, e, 0)
				case sim.Crystal:
					t.Entity, t.Size = e, 1+w.Rand.Intn(5)
				}
			})
			byEntity, byEntityAll := map[sim.Entity]int{}, map[sim.Entity]int{}
			w.Each(func(_, _ int, t *sim.Tile) {
				byEntityAll[t.Entity]++
				if t.Entity != sim.Empty {
					byEntity[t.Entity]++
				}
			})
			wants := map[string]bool{fmt.Sprint(byEntity): true, fmt.Sprint(byEntityAll): true}
			for _, m := range []map[sim.Entity]int{byEntity, byEntityAll} {
				byName := map[string]int{}
				for e, n := range m {
					byName[e.String()] = n
				}
				wants[fmt.Sprint(byName)] = true
			}
			w.Check = func(w *sim.World) (bool, string) {
				for _, r := range w.Reports {
					if wants[r] {
						return true, i18n.T("التعداد مظبوط", "The census is correct")
					}
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[n-1] + i18n.T(" — الأرقام مش مظبوطة", " — the counts are incorrect")
				}
				return false, i18n.T("القاعدة مستنية map بالأعداد", "Base is waiting for a map of counts")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	counts := map[colony.Entity]int{}

	e := colony.GetEntity()
	counts[e]++
	fmt.Println(counts)

	// لف على الخريطة كلها (6×4) وعدّ، وبعدين ابعت counts
}
`,
		Hints: []string{
			"الـmap بتعد كل نوع؛ ~len~ بتعد المفاتيح.",
			"ابدأ بعدّ الـstarter، وزُر باقي الخانات مرة.",
			"راجع الـmap وسطر ~colony.GetEntity()~: هل المفتاح متوافق؟",
		},
		Lesson: `# Maps: قاموس المستعمرة 🗂️

الـ slice بتوصل للعناصر برقم (0، 1، 2…). الـ **map** بتوصل للقيمة بـ **مفتاح** من أي نوع.

~~~go
counts := map[string]int{}   // مفتاح string، قيمة int
counts["moss"] = 3
counts["moss"]++             // 4
fmt.Println(counts["lichen"]) // 0 ← المفتاح مش موجود؟ بترجع القيمة الصفرية
~~~

## اتأكد إن المفتاح موجود

~~~go
n, ok := counts["crystal"]
if !ok {
	fmt.Println("مفيش كريستال خالص")
}
~~~

ده اسمه **comma ok**، وهتشوفه كتير في Go.

## عمليات تانية

~~~go
delete(counts, "moss")       // امسح مفتاح
len(counts)                  // عدد المفاتيح
for k, v := range counts {   // لف على الكل
	fmt.Println(k, v)
}
~~~

## سرّ مهم 🎲
لما تلف على map بـ range، **الترتيب عشوائي** كل مرة. Go عامل كده بقصد، عشان متعتمدش على ترتيب مش مضمون. (بس ~fmt.Println~ بيطبعها مرتبة عشان تقراها بسهولة.)

## map لازم تتعمل قبل ما تكتب فيها
~~~go
var m map[string]int // nil map
m["x"] = 1           // panic!
m = map[string]int{} // أو make(map[string]int)
~~~

## المهمة
لف على كل خانة، وعدّ كل نوع في map، وابعتها للقاعدة.`,
	})

	stockOrderStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.Measure())
	// Build the report described in the lesson.
}
`
	add(&Level{ID: "6-4-stock-order", Chapter: 6, Title: "ترتيب المخزون", Goal: "رتّب مواقع الكريستال بالحجم تنازلي وبالموقع عند التعادل",
		Requires: []string{"sel:slices.SortFunc"}, Bridge: []int{19},
		Par: 19, Gold: 19, MaxTicks: 50,
		Starter: stockOrderStarter, Lesson: `# ترتيب له قاعدة ثابتة 📦

مهمّة أغلى 3 عرّفتك على ترتيب الأرقام. دلوقتي القاعدة مركّبة: الحجم الأكبر الأول، ولو اتنين نفس الحجم، الموقع الأصغر يسبق. مش محتاج struct عشان تحتفظ بالمعلومتين.

## قبل ما تكتب
- هتحافظ على الربط بين الحجم وموقعه إزاي؟
- لو الحجم اتعادل، إيه اللي يمنع ترتيب النتيجة يتغيّر؟

## قايمتين بينهم نفس الفهرس
~~~go
sizes := make([]int, 3)
positions := make([]int, 3)
sizes[0] = 9
positions[0] = 0
~~~
~sizes[x]~ حجم الكريستال عند x، و~positions~ قايمة المواقع اللي هنرتّبها. ~make~ والفهرسة هم نفس أدوات 6-1. الأحجام تفضل ثابتة؛ اللي يتحرّك في الترتيب هو أرقام المواقع بس.

~~~go
width, _ := colony.WorldSize()
size := colony.Measure()
fmt.Println(width, size)
~~~
~Measure~ بتقرأ حجم كريستال خانة الدرون الحالية، وبتاخد tick. اقرأ عرض الصف بـ ~WorldSize~، وقيس كل خانة ثم اتحرك East إلا بعد آخر واحدة.

## comparator بتقول مين يسبق
~~~go
numbers := []int{8, 2, 5}
slices.SortFunc(numbers, func(a, b int) int {
 return a - b
})
fmt.Println(numbers) // [2 5 8]
~~~
ضيف ~"slices"~ للـ imports. ~SortFunc~ بتعدّل الـ slice نفسها، وبتاخد دالة مقارنة زي قيم الدوال اللي عرفناها في 5-5. هنا الأعداد صغيرة فطرحها آمن.

الدالة لازم ترجع عدد: سالب لو a تسبق b، صفر لو متساويين، وموجب لو a تيجي بعد b. في المثال بنرتّب الأرقام تصاعدي؛ دي مش قاعدة أحجام المهمة.

## كسر تعادل: مثال محلول
~~~go
ranks := []int{4, 4, 2}
order := []int{1, 0, 2}
slices.SortFunc(order, func(a, b int) int {
 if ranks[a] < ranks[b] { return -1 }
 if ranks[a] > ranks[b] { return 1 }
 return a - b
})
fmt.Println(order) // [2 0 1]
~~~
الدالة بتلتقط ranks كـ closure وبتقارن المواقع اللي في order. الرتبة الأصغر الأول في المثال، وعند التعادل الموقع الأصغر يسبق. لاحظ إننا ما رتّبناش ranks ولا ضيّعنا الربط بينها وبين المواقع.

## أخطاء شائعة
~~~go
return sizes[a] >= sizes[b] // مش هيتبني
~~~
ده bool، لكن comparator محتاجة int. ولو تجاهلت التعادل وقلت صفر لكل حجمين متساويين، ما حدّدتش أنهي موقع يسبق. كمان متكتبش المقارنة على ~sizes[a]~ بعد ما تكون رتّبت sizes وغيّرت علاقتها بالموقع.

## المهمة
قيس صف الكريستال كله عند y=0. رتّب المواقع بـ ~slices.SortFunc~: الحجم تنازلي، وبعده الموقع x تصاعدي لو الأحجام متساوية. ابعت ~positions~ كلها بتقرير واحد، مثلاً ~[2 0 1]~. متحصدش ولا تبعت الأحجام بدل المواقع.`, Hints: []string{"اعمل ~sizes~ ثابتة و~positions~ تبدأ بمواقع الصف.", "قيمة a في comparator هي موقع، فحجمه ~sizes[a]~ مش a نفسه.", "الحجم الأكبر يطلع نتيجة سالبة عشان يسبق؛ عند التعادل استخدم فرق الموقعين الصغيرين."},
		En: &LevelText{Title: "Stock Order", Goal: "Rank crystal positions by descending size and ascending position on ties", Starter: stockOrderStarter, Lesson: `# Sorting With a Fixed Tie Rule 📦

Top Three introduced sorting numbers. Now use a two-part rule: larger size first, then smaller position on equal sizes. You do not need a struct to retain both facts.

## Before you write
- How will you preserve the link between a size and its position?
- What makes equal-size results deterministic?

## Two collections linked by index
~~~go
sizes := make([]int, 3)
positions := make([]int, 3)
sizes[0] = 9
positions[0] = 0
~~~
~sizes[x]~ holds the crystal size at x; ~positions~ is the list of positions to sort. ~make~ and indexing are the tools from 6-1. Keep sizes fixed and reorder only position numbers.

~~~go
width, _ := colony.WorldSize()
size := colony.Measure()
fmt.Println(width, size)
~~~
~Measure~ reads the crystal size at the drone's current tile and costs one tick. Read the row width with ~WorldSize~, measure each tile, and move East except after the last tile.

## A comparator says which comes first
~~~go
numbers := []int{8, 2, 5}
slices.SortFunc(numbers, func(a, b int) int {
 return a - b
})
fmt.Println(numbers) // [2 5 8]
~~~
Import ~"slices"~. ~SortFunc~ changes the slice in place and accepts a comparison function, using function values from 5-5. These sample numbers are small enough for subtraction to be safe.

The comparator returns an integer: negative if a comes before b, zero for equality, positive if a follows b. The example sorts ascending numbers, not the mission's sizes.

## Breaking ties: a worked example
~~~go
ranks := []int{4, 4, 2}
order := []int{1, 0, 2}
slices.SortFunc(order, func(a, b int) int {
 if ranks[a] < ranks[b] { return -1 }
 if ranks[a] > ranks[b] { return 1 }
 return a - b
})
fmt.Println(order) // [2 0 1]
~~~
The function captures ranks as a closure and compares positions from order. This example puts smaller ranks first and breaks ties by smaller position. We never sort ranks or lose its link to positions.

## Common mistakes
~~~go
return sizes[a] >= sizes[b] // Will not build
~~~
This returns bool, but the comparator needs int. Returning zero for all equal sizes also fails to specify which position comes first. Do not compare ~sizes[a]~ after sorting sizes and destroying its link to position.

## Your mission
Measure the entire crystal row at y=0. Sort positions with ~slices.SortFunc~: descending size, then ascending x for ties. Report the whole ~positions~ slice once, for example ~[2 0 1]~. Do not harvest or report sizes instead of positions.`, Hints: []string{"Keep ~sizes~ fixed and initialize ~positions~ with the row positions.", "Comparator argument a is a position, so its size is ~sizes[a]~, not a itself.", "A larger size needs a negative comparison result to come first; on ties use the small position difference."}},
		Setup: func(w *sim.World) {
			n := 6 + w.Rand.Intn(5)
			w.Init(n, 2, sim.Regolith)
			positions := make([]int, n)
			sizes := make([]int, n)
			for x := range positions {
				positions[x] = x
				sizes[x] = 1 + w.Rand.Intn(5)
				w.At(x, 0).Entity = sim.Crystal
				w.At(x, 0).Size = sizes[x]
			}
			sizes[n-1] = sizes[0]
			w.At(n-1, 0).Size = sizes[0]
			sort.Slice(positions, func(i, j int) bool {
				a, b := positions[i], positions[j]
				if sizes[a] != sizes[b] {
					return sizes[a] > sizes[b]
				}
				return a < b
			})
			want := fmt.Sprint(positions)
			w.Check = func(w *sim.World) (bool, string) {
				if w.Harvests == 0 && len(w.Reports) == 1 && w.Reports[0] == want {
					return true, i18n.T("المواقع مرتبة صح، بما فيها التعادل", "Positions are correctly ranked, including ties")
				}
				return false, i18n.T("مستني المواقع بالحجم تنازلي وبالموقع تصاعدي عند التعادل، من غير حصاد", "Waiting for positions by descending size and ascending position on ties, without harvesting")
			}
		},
	})

	heatGridStarter := `package main

import "fmt"

func OverlayHeat(a, b [][]float64) ([][]float64, error) {
	return nil, fmt.Errorf("implement OverlayHeat")
}
func ThroughHeat(a, b [][]float64) ([][]float64, error) {
	return nil, fmt.Errorf("implement ThroughHeat")
}
func FlipHeat(a [][]float64) ([][]float64, error) { return nil, fmt.Errorf("implement FlipHeat") }
func main()                                       {}
`
	heatFiles := []LabFile{{Name: "main.go", Starter: heatGridStarter}}
	add(&Level{ID: "6-5-heat-grid", Chapter: 6, Title: "شبكة الحرارة", Goal: "نفّذ عمليات شبكة الحرارة بـ 2D slices من غير تعديل المدخل",
		Requires: []string{"slice"}, Bridge: []int{26},
		Kind: "lab", Files: heatFiles,
		Starter: heatGridStarter, Lesson: `# صفوف حرارة جوه slices 🌡️

الحسّاس الواحد بيدي رقم، لكن لوحة المحطة فيها صفوف وأعمدة. هنستخدم slice من slices عشان نمثل شبكة، ونبني نتائج جديدة من غير ما نلمس قراءات المصدر.

## قبل ما تكتب
- إزاي تتأكد إن كل الصفوف بنفس الطول قبل الفهرسة؟
- هتمنع صف الناتج إنه يشارك ذاكرة صف من المدخل إزاي؟

## أول lab
هنا بتعدّل ملف ~main.go~، وزر Run بيشغّل اختبارات بدلاً من حركة الدرون. مش محتاج تكتب tests بنفسك؛ ركّز على الدوال الثلاثة الموجودة في الملف. عند النجاح الاختبارات المخفية هي اللي تحدد النتيجة، مش Report.

## slice من صفوف
~~~go
sectors := [][]float64{{2, 4, 6}, {1, 3, 5}}
fmt.Println(len(sectors), len(sectors[0])) // 2 3
fmt.Println(sectors[1][2]) // 5
~~~
~[][]float64~ يعني كل عنصر في الـ outer slice صف نوعه ~[]float64~. أول index يختار الصف، والتاني العمود. هنا الخريطة 2×3؛ الأعداد العشرية والسالبة صالحة.

## افحص الشكل قبل الحساب
~~~go
func columns(grid [][]float64) (int, error) {
 if len(grid) == 0 { return 0, fmt.Errorf("no rows") }
 return len(grid[0]), nil
}
~~~
المثال بيفحص وجود صف قبل ~grid[0]~؛ ده جزء من الفحص بس. فحصك الكامل لازم يرفض صفر صفوف، أو كل الصفوف فاضية، أو أي اختلاف في أطوال الصفوف. صف فاضي وسط صفوف غير فاضية معناه شبكة غير مستطيلة.

الدوال بترجّع نتيجة وerror زي 5-2. استخدم ~fmt.Errorf~ برسالة مفهومة للمدخل الغلط؛ متعملش panic. افحص المدخلين قبل أي جمع أو ضرب.

## خصّص صفوف مستقلة: مثال محلول
~~~go
copyGrid := make([][]float64, len(sectors))
for i, row := range sectors {
 copyGrid[i] = make([]float64, len(row))
 for j, reading := range row {
  copyGrid[i][j] = reading * 2
 }
}
fmt.Println(copyGrid[0][1]) // 8
~~~
المثال بيضاعف خريطة قطاعات تجريبية. خصّصنا outer slice وبعدين صف جديد لكل i. كده تعديل الناتج مش هيغيّر المصدر، وصف مش هيغيّر صف تاني. الحلقات المتداخلة نفس فكرة 3-2؛ ~range~ من 6-2.

## ثلاث عمليات مختلفة
~~~go
sum, err := OverlayHeat(sectors, sectors)
fmt.Println(sum, err)
~~~
~OverlayHeat(a,b)~ بتجمع عنصر بعنصر، فالشكلين لازم يبقوا متساويين: نفس عدد الصفوف ونفس عدد الأعمدة.

~~~go
flipped, err := FlipHeat(sectors)
fmt.Println(flipped, err)
~~~
~FlipHeat(a)~ هي transpose: الخلية ~a[i][j]~ تروح للناتج عند ~[j][i]~. المثال 2×3 يبقى 3×2، وأول صف جديد يبقى ~[2 1]~.

~~~go
coupler := [][]float64{{1}, {0}, {2}}
through, err := ThroughHeat(sectors, coupler)
fmt.Println(through, err) // [[14] [11]] <nil>
~~~
~ThroughHeat(a,b)~ ضرب مصفوفات. أعمدة a لازم تساوي صفوف b. النتيجة صفوف a × أعمدة b؛ كل خلية ~i,j~ مجموع ~a[i][k]*b[k][j]~ على كل k في البعد المشترك. المثال وضّح خلية 0,0: ~2*1 + 4*0 + 6*2 = 14~.

## أخطاء شائعة
~~~go
out := sectors
out[0][0] = 99 // غيّرت المصدر كمان
~~~
ده مش نسخ شبكة. حتى outer slice جديدة مش كفاية لو صفوفها نفس صفوف المصدر. وبرضه ضرب ~a[i][j]*b[i][j]~ مش ThroughHeat؛ لازم مجموع على البعد المشترك.

## المهمة
كمّل ~OverlayHeat~ و~ThroughHeat~ و~FlipHeat~ بتواقيع الـ starter. كل نتيجة ~[][]float64~ جديدة، وكل صف مستقل عن المدخلات وعن باقي الصفوف. متعدّلش أي مدخل، حتى في مسار الخطأ.

رجّع error غير nil لو المصفوفة فاضية، غير مستطيلة، أو الأبعاد مش مناسبة للعملية. 1×1 وrow/column vectors صالحين. نص الخطأ مش لازم يكون حرفي، بس يكون مفهوم. مفيش Report مطلوب. الاختبارات بتقارن الأرقام بسماحية ~1e-9~ نسبية وبحد أدنى مطلق ~1e-9~؛ مش مطلوب كشف overflow في float64.`, Hints: []string{"اكتب فحص شكل مشترك لكل الدوال؛ افحص وجود الصفوف قبل ~grid[0]~.", "الناتج محتاج outer slice وصفوف جديدة، مش مجرد إسناد أو نسخ outer slice.", "للضرب خلي i صف الناتج وj عموده، واجمع على k؛ للـ transpose بدّل الفهرسين في الناتج."},
		En: &LevelText{Title: "Heat Grid", Goal: "Implement heat operations with 2D slices without mutating inputs", Starter: heatGridStarter, Lesson: `# Rows of Heat Inside Slices 🌡️

One sensor gives one number, but a station panel has rows and columns. Use a slice of slices to represent a grid and produce fresh results without changing source readings.

## Before you write
- How will you validate every row length before indexing?
- How will you keep output rows from sharing input storage?

## Your first lab
Edit ~main.go~ here. Run executes tests rather than moving a drone. You do not need to write tests yourself; implement the three supplied functions. Hidden tests determine success instead of Report.

## A slice of rows
~~~go
sectors := [][]float64{{2, 4, 6}, {1, 3, 5}}
fmt.Println(len(sectors), len(sectors[0])) // 2 3
fmt.Println(sectors[1][2]) // 5
~~~
~[][]float64~ means each outer-slice element is a row of type ~[]float64~. The first index selects a row, the second a column. This grid is 2×3; negative and decimal values are valid.

## Validate shape before computation
~~~go
func columns(grid [][]float64) (int, error) {
 if len(grid) == 0 { return 0, fmt.Errorf("no rows") }
 return len(grid[0]), nil
}
~~~
This checks that a row exists before ~grid[0]~; it is only part of validation. Complete validation must reject zero rows, all-empty rows, and unequal row lengths. Mixing empty and nonempty rows is ragged input.

Functions return a result and error as in 5-2. Use ~fmt.Errorf~ with a descriptive message for invalid inputs rather than panicking. Validate both operands before addition or multiplication.

## Allocate independent rows: a worked example
~~~go
copyGrid := make([][]float64, len(sectors))
for i, row := range sectors {
 copyGrid[i] = make([]float64, len(row))
 for j, reading := range row {
  copyGrid[i][j] = reading * 2
 }
}
fmt.Println(copyGrid[0][1]) // 8
~~~
This doubles a sample sector map. We allocate the outer slice, then a new row for every i. Changing the result cannot change the source or another row. Nested loops come from 3-2 and ~range~ from 6-2.

## Three different operations
~~~go
sum, err := OverlayHeat(sectors, sectors)
fmt.Println(sum, err)
~~~
~OverlayHeat(a,b)~ adds corresponding elements, so both shapes must match: equal row and column counts.

~~~go
flipped, err := FlipHeat(sectors)
fmt.Println(flipped, err)
~~~
~FlipHeat(a)~ transposes: cell ~a[i][j]~ moves to output ~[j][i]~. The 2×3 example becomes 3×2, whose first row is ~[2 1]~.

~~~go
coupler := [][]float64{{1}, {0}, {2}}
through, err := ThroughHeat(sectors, coupler)
fmt.Println(through, err) // [[14] [11]] <nil>
~~~
~ThroughHeat(a,b)~ performs matrix multiplication. Columns of a must equal rows of b. The result has a's row count and b's column count. Each cell ~i,j~ sums ~a[i][k]*b[k][j]~ across the shared dimension k. Example cell 0,0 is ~2*1 + 4*0 + 6*2 = 14~.

## Common mistakes
~~~go
out := sectors
out[0][0] = 99 // Also changes the source
~~~
This does not copy a grid. Even a fresh outer slice is insufficient if its rows reuse source rows. Also, ~a[i][j]*b[i][j]~ is not ThroughHeat; sum across the shared dimension.

## Your mission
Implement ~OverlayHeat~, ~ThroughHeat~, and ~FlipHeat~ with the starter signatures. Return fresh ~[][]float64~ results with rows independent of inputs and each other. Never mutate inputs, even on errors.

Return a non-nil error for empty or ragged inputs or incompatible operation dimensions. Single cells and row/column vectors are valid. Error wording is flexible but must be descriptive. No Report is required. Tests use relative tolerance ~1e-9~ with an absolute floor of ~1e-9~; float64 overflow detection is not required.`, Hints: []string{"Share shape validation across functions; check row existence before ~grid[0]~.", "Allocate a fresh outer slice and fresh rows, rather than assigning or copying only the outer slice.", "For multiplication i is the output row and j its column; sum over k. For transpose swap the output indices."}, Files: heatFiles},
		Setup: func(w *sim.World) {
			labRoom(w)
			for x := 1; x < 6; x++ {
				w.At(x, 1).Size = 1 + w.Rand.Intn(9)
			}
		},
		Lab: &Lab{ReqFile: "main.go", Hidden: `package main

import (
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func hiddenCopy(grid [][]float64) [][]float64 {
	if grid == nil {
		return nil
	}
	out := make([][]float64, len(grid))
	for i, row := range grid {
		if row != nil {
			out[i] = append([]float64{}, row...)
		}
	}
	return out
}
func hiddenEqual(t *testing.T, got, want [][]float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("rows: got %v want %v", got, want)
	}
	for i, row := range want {
		if len(got[i]) != len(row) {
			t.Fatalf("columns: got %v want %v", got, want)
		}
		for j, n := range row {
			if math.IsNaN(got[i][j]) || math.IsInf(got[i][j], 0) || math.Abs(got[i][j]-n) > 1e-9*math.Max(1, math.Abs(n)) {
				t.Fatalf("cell %d,%d: got %v want %v", i, j, got, want)
			}
		}
	}
}
func hiddenIndependent(t *testing.T, out, a, b [][]float64) {
	t.Helper()
	beforeA, beforeB, snapshot := hiddenCopy(a), hiddenCopy(b), hiddenCopy(out)
	for i, row := range snapshot {
		for j, n := range row {
			out[i][j] = n + 12345
			if !reflect.DeepEqual(a, beforeA) || !reflect.DeepEqual(b, beforeB) {
				t.Fatal("output aliases an operand")
			}
			for r, other := range snapshot {
				for c, want := range other {
					if r == i && c == j {
						continue
					}
					if out[r][c] != want {
						t.Fatal("output rows alias one another")
					}
				}
			}
			out[i][j] = n
		}
	}
}
func TestHiddenHeatExamples(t *testing.T) {
	cases := []struct {
		name       string
		a, b, want [][]float64
		run        func([][]float64, [][]float64) ([][]float64, error)
	}{
		{"overlay", [][]float64{{1, -2}, {3, 4}}, [][]float64{{5, 2}, {-3, 8}}, [][]float64{{6, 0}, {0, 12}}, OverlayHeat},
		{"through", [][]float64{{1, 2, 3}, {4, 5, 6}}, [][]float64{{7, 8}, {9, 10}, {11, 12}}, [][]float64{{58, 64}, {139, 154}}, ThroughHeat},
		{"flip rectangular", [][]float64{{1, 2, 3}, {4, 5, 6}}, nil, [][]float64{{1, 4}, {2, 5}, {3, 6}}, func(a, b [][]float64) ([][]float64, error) { return FlipHeat(a) }},
		{"one cell product", [][]float64{{-0.5}}, [][]float64{{2}}, [][]float64{{-1}}, ThroughHeat},
		{"one cell sum", [][]float64{{-0.5}}, [][]float64{{2}}, [][]float64{{1.5}}, OverlayHeat},
		{"row vector", [][]float64{{1, 2, 3}}, nil, [][]float64{{1}, {2}, {3}}, func(a, b [][]float64) ([][]float64, error) { return FlipHeat(a) }},
		{"column vector", [][]float64{{1}, {2}, {3}}, nil, [][]float64{{1, 2, 3}}, func(a, b [][]float64) ([][]float64, error) { return FlipHeat(a) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			beforeA, beforeB := hiddenCopy(c.a), hiddenCopy(c.b)
			got, err := c.run(c.a, c.b)
			if err != nil {
				t.Fatal(err)
			}
			hiddenEqual(t, got, c.want)
			if !reflect.DeepEqual(c.a, beforeA) || !reflect.DeepEqual(c.b, beforeB) {
				t.Fatal("operation mutated inputs")
			}
			hiddenIndependent(t, got, c.a, c.b)
		})
	}
}
func TestHiddenHeatInvalid(t *testing.T) {
	valid := [][]float64{{1, 2}, {3, 4}}
	invalid := []struct {
		name string
		a    [][]float64
	}{
		{"nil", nil}, {"zero rows", [][]float64{}}, {"zero columns", [][]float64{{}, {}}},
		{"ragged", [][]float64{{1, 2}, {3}}}, {"empty first row", [][]float64{{}, {1}}}, {"empty later row", [][]float64{{1}, {}}},
	}
	for _, c := range invalid {
		t.Run(c.name, func(t *testing.T) {
			before, beforeValid := hiddenCopy(c.a), hiddenCopy(valid)
			checks := []func() ([][]float64, error){func() ([][]float64, error) { return FlipHeat(c.a) }, func() ([][]float64, error) { return OverlayHeat(c.a, valid) }, func() ([][]float64, error) { return OverlayHeat(valid, c.a) }, func() ([][]float64, error) { return ThroughHeat(c.a, valid) }, func() ([][]float64, error) { return ThroughHeat(valid, c.a) }}
			for _, run := range checks {
				_, err := run()
				if err == nil || strings.TrimSpace(err.Error()) == "" {
					t.Fatal("invalid grid needs a descriptive error")
				}
				if !reflect.DeepEqual(c.a, before) || !reflect.DeepEqual(valid, beforeValid) {
					t.Fatal("input mutated on error")
				}
			}
		})
	}
	mismatches := []struct {
		name string
		b    [][]float64
		run  func([][]float64, [][]float64) ([][]float64, error)
	}{
		{"add rows", [][]float64{{1, 2}}, OverlayHeat}, {"add columns", [][]float64{{1}, {2}}, OverlayHeat}, {"product inner dimension", [][]float64{{1, 2}}, ThroughHeat},
	}
	for _, c := range mismatches {
		t.Run(c.name, func(t *testing.T) {
			beforeA, beforeB := hiddenCopy(valid), hiddenCopy(c.b)
			if _, err := c.run(valid, c.b); err == nil {
				t.Fatal("incompatible dimensions accepted")
			}
			if !reflect.DeepEqual(valid, beforeA) || !reflect.DeepEqual(c.b, beforeB) {
				t.Fatal("shape failure mutated input")
			}
		})
	}
}
func hiddenGrid(rng *rand.Rand, rows, cols int) [][]float64 {
	grid := make([][]float64, rows)
	for i := range grid {
		grid[i] = make([]float64, cols)
		for j := range grid[i] {
			grid[i][j] = float64(rng.Intn(41)-20) / 4
		}
	}
	return grid
}
func TestHiddenHeatSeeded(t *testing.T) {
	baseSeed := rand.Int63()
	t.Logf("heat grid seed: %d", baseSeed)
	for trial := int64(0); trial < 30; trial++ {
		rng := rand.New(rand.NewSource(baseSeed + trial))
		rows, inner, cols := 1+rng.Intn(4), 1+rng.Intn(4), 1+rng.Intn(4)
		a, b := hiddenGrid(rng, rows, inner), hiddenGrid(rng, inner, cols)
		beforeA, beforeB := hiddenCopy(a), hiddenCopy(b)
		got, err := ThroughHeat(a, b)
		if err != nil {
			t.Fatal(err)
		}
		want := make([][]float64, rows)
		for i := range want {
			want[i] = make([]float64, cols)
			for j := 0; j < cols; j++ {
				for k := 0; k < inner; k++ {
					want[i][j] += beforeA[i][k] * beforeB[k][j]
				}
			}
		}
		hiddenEqual(t, got, want)
		if !reflect.DeepEqual(a, beforeA) || !reflect.DeepEqual(b, beforeB) {
			t.Fatal("multiply mutated inputs")
		}
		hiddenIndependent(t, got, a, b)
		flipped, err := FlipHeat(a)
		if err != nil {
			t.Fatal(err)
		}
		flipWant := make([][]float64, inner)
		for j := range flipWant {
			flipWant[j] = make([]float64, rows)
			for i := 0; i < rows; i++ {
				flipWant[j][i] = beforeA[i][j]
			}
		}
		hiddenEqual(t, flipped, flipWant)
		twice, err := FlipHeat(flipped)
		if err != nil {
			t.Fatal(err)
		}
		hiddenEqual(t, twice, beforeA)
		hiddenIndependent(t, flipped, a, nil)
		other := hiddenGrid(rng, rows, inner)
		beforeOther := hiddenCopy(other)
		sum, err := OverlayHeat(a, other)
		if err != nil {
			t.Fatal(err)
		}
		sumWant := hiddenCopy(beforeA)
		for i := range sumWant {
			for j := range sumWant[i] {
				sumWant[i][j] += beforeOther[i][j]
			}
		}
		hiddenEqual(t, sum, sumWant)
		if !reflect.DeepEqual(a, beforeA) || !reflect.DeepEqual(other, beforeOther) {
			t.Fatal("overlay or flip mutated inputs")
		}
		hiddenIndependent(t, sum, a, other)
	}
	zero := [][]float64{{0, 0}, {0, 0}}
	a := [][]float64{{-1.25, 2}, {3, 4.5}}
	identity := [][]float64{{1, 0}, {0, 1}}
	for _, pair := range [][2][][]float64{{a, identity}, {identity, a}} {
		got, err := ThroughHeat(pair[0], pair[1])
		if err != nil {
			t.Fatal(err)
		}
		hiddenEqual(t, got, a)
	}
	sum, err := OverlayHeat(a, zero)
	if err != nil {
		t.Fatal(err)
	}
	hiddenEqual(t, sum, a)
	product, err := ThroughHeat(a, zero)
	if err != nil {
		t.Fatal(err)
	}
	hiddenEqual(t, product, zero)
}
`},
	})

}
