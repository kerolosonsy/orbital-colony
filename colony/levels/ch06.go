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
			Title: "Crystal Log",
			Goal:  "Measure every crystal and report their sizes in order in one slice",
			Lesson: `# Slices: A List That Grows 📋

So far each variable has held one value. What if you want to remember 8 sizes? Or 1000? Use a **slice**.

~~~go
var sizes []int            // An empty int slice (length 0)
sizes = append(sizes, 5)   // [5]
sizes = append(sizes, 2)   // [5 2]
fmt.Println(len(sizes))    // 2
fmt.Println(sizes[0])      // 5  ← Counting starts at 0
~~~

- ~[]int~ means "a slice of int."
- ~append~ returns a **new slice value**, so assign it back: ~sizes = append(sizes, x)~.
- ~len(s)~ tells you how many elements it holds.
- ~s[i]~ retrieves element i. If i is outside the bounds, the program **panics** (panic: index out of range).

## Other ways to create a slice

~~~go
nums := []int{4, 8, 15}     // With initial values
grid := make([]int, 10)     // 10 zeros, ready to use: grid[3] = 7
part := nums[1:3]           // [8 15]: from 1 up to, but not including, 3
~~~

## Your mission
Measure the 8 crystals from left to right, store their sizes in a slice, and report it to base.`,
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
			Hints: []string{
				"Loop 8 times: Measure and append, then Move (except on the last iteration).",
				"~colony.Report(sizes)~ sends the whole slice at once, looking like ~[3 7 1 ...]~.",
			},
		},
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
		Hints: []string{
			"حلقة 8 لفّات: Measure وappend، وبعدين Move (إلا آخر لفّة).",
			"~colony.Report(sizes)~ بتبعت الـ slice كلها مرة واحدة، شكلها هيبقى ~[3 7 1 ...]~.",
		},
		Lesson: `# Slices: قايمة بتكبر 📋

لحد دلوقتي كل متغير كان فيه قيمة واحدة. بس لو عايز تفتكر 8 أحجام؟ أو 1000؟ هنا بنستخدم **slice**.

~~~go
var sizes []int            // slice من int، فاضية (طولها 0)
sizes = append(sizes, 5)   // [5]
sizes = append(sizes, 2)   // [5 2]
fmt.Println(len(sizes))    // 2
fmt.Println(sizes[0])      // 5  ← العدّ من 0
~~~

- ~[]int~ معناها "slice من int".
- ~append~ بترجّع slice **جديدة**، فلازم تحطها في نفس المتغير: ~sizes = append(sizes, x)~.
- ~len(s)~ بيقولك فيها كام عنصر.
- ~s[i]~ بيجيبلك العنصر رقم i. ولو i برا الحدود، البرنامج **بيقع** (panic: index out of range).

## طرق تانية تعمل بيها slice

~~~go
nums := []int{4, 8, 15}     // بقيم جاهزة
grid := make([]int, 10)     // 10 أصفار جاهزين: grid[3] = 7
part := nums[1:3]           // [8 15]: من 1 لحد قبل 3
~~~

## المهمة
قيس الكريستال الـ 8 بالترتيب من الشمال لليمين، واحفظهم في slice، وابعتها للقاعدة.`,
	})

	add(&Level{
		En: &LevelText{
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
				"Copy and sort the slice: ~sorted := slices.Clone(sizes)~, then ~slices.Sort(sorted)~ (requires ~import \"slices\"~). The largest 3 are the last 3.",
				"The threshold is ~sorted[len(sorted)-3]~. Any crystal at least that large belongs to the top 3.",
				"Return west from x = 9 to 0, harvesting only where ~sizes[x] >= threshold~.",
			},
		},
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
			"انسخ الـ slice وارتّبها: ~sorted := slices.Clone(sizes)~ ثم ~slices.Sort(sorted)~ (محتاج ~import \"slices\"~). أكبر 3 هيبقوا آخر 3.",
			"الحد الأدنى = ~sorted[len(sorted)-3]~. أي كريستال حجمه أكبر منه أو بيساويه يبقى من الـ 3.",
			"ارجع غرب من x = 9 لـ 0، واحصد بس لو ~sizes[x] >= threshold~.",
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
				"Use the chapter 3 snake route, but replace Harvest with ~counts[colony.GetEntity()]++~.",
				"You do not need to count empty tiles, but including them is accepted too.",
			},
		},
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
			"نفس لفّة الثعبان بتاعة الفصل 3، بس بدل Harvest هتعمل ~counts[colony.GetEntity()]++~.",
			"الخانات الفاضية (empty) مش لازم تتعد، بس لو عدّيتها برضه مقبولة.",
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
}
