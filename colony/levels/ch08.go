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
			LessonRevision: 1, Steps: ch08CopyEnglish,
			Title: "The Counter That Does Not Count",
			Goal:  "Fix the bug: base must receive the actual crystal total",
			Lesson: `# Go Copies Every Argument 📋

Try the starter: the drone collects crystal, but the total is **0**. Why?

When you pass a variable to a function, Go passes a **copy** of its value (pass by value). The function increments the copy, which disappears when the function returns.

~~~go
func addHarvest(total int, n int) {
	total += n // Changes only the copy
}
~~~

## A pointer: a variable's address

A pointer is a value that refers to a variable's **location** in memory.

~~~go
x := 5
p := &x        // & means "the address of x." p has type *int
fmt.Println(*p) // 5   ← * means "read the value at this address"
*p = 9
fmt.Println(x)  // 9   ← x itself changed!
~~~

| | Meaning |
|---|---|
| ~*int~ | A type: "pointer to int" |
| ~&x~ | The address of x |
| ~*p~ | The value at address p |

## The solution

~~~go
func addHarvest(total *int, n int) {
	*total += n
}

addHarvest(&total, n)
~~~

## When should you use a pointer?
1. When a function **needs to change** a caller's variable.
2. When a value is **large** (a struct with many fields) and you want to avoid copying it every time.

Otherwise, returning a value with ~return~ is usually clearer.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// addHarvest should add n to total... but something is wrong
func addHarvest(total int, n int) {
	total += n
}

func main() {
	total := 0
	for x := 0; x < 8; x++ {
		n, _ := colony.TryHarvest()
		addHarvest(total, n)
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	fmt.Println("Total:", total)
	colony.Report(total)
}
`,
			Hints: []string{
				"An ordinary parameter receives a value copy; a pointer refers to an address.",
				"Inspect the function, its call site, then compare the printed value with storage.",
				"Check the ~addHarvest~ declaration and its call in the starter.",
			},
		},
		LessonRevision: 1, Steps: ch08CopyArabic,
		ID: "8-1-copy-bug", Chapter: 8, Title: "العدّاد اللي مبيعدّش",
		Goal: "صلّح الـ bug: القاعدة لازم تستلم مجموع الكريستال الحقيقي",
		Par:  15, Gold: 15, MaxTicks: 60,
		Requires: []string{"pointer"},
		Setup: func(w *sim.World) {
			w.Init(8, 2, sim.Regolith)
			sum := 0
			for x := 0; x < 8; x++ {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
				sum += t.Size
			}
			want := strconv.Itoa(sum)
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("المجموع صح", "The total is correct")
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[n-1] + i18n.T(" بس الدرون جمع ", " but the drone collected ") + strconv.Itoa(w.Inv[sim.Crystal])
				}
				return false, i18n.T("القاعدة مستنية المجموع", "Base is waiting for the total")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// addHarvest المفروض تزوّد total بـ n... بس فيه حاجة غلط
func addHarvest(total int, n int) {
	total += n
}

func main() {
	total := 0
	for x := 0; x < 8; x++ {
		n, _ := colony.TryHarvest()
		addHarvest(total, n)
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	fmt.Println("المجموع:", total)
	colony.Report(total)
}
`,
		Hints: []string{
			"الـparameter العادي بياخد نسخة من القيمة، والـpointer بيشاور على عنوان.",
			"راجع الدالة، مكان استدعائها، وبعدها قارن الرقم المطبوع بالمخزون.",
			"افحص تعريف ~addHarvest~ وسطر استدعائها في الـstarter.",
		},
		Lesson: `# Go بينسخ كل حاجة 📋

جرّب الكود زي ما هو: الدرون هيجمع الكريستال، بس المجموع هيطلع **0**. ليه؟

لما تبعت متغير لدالة، Go بيبعت **نسخة** من قيمته (pass by value). الدالة بتزوّد النسخة، والنسخة بتختفي أول ما الدالة تخلص.

~~~go
func addHarvest(total int, n int) {
	total += n // بيغيّر النسخة بس
}
~~~

## الـ Pointer: عنوان المتغير

الـ pointer قيمة بتشاور على **مكان** متغير في الذاكرة.

~~~go
x := 5
p := &x        // & = "هات عنوان x". نوع p هو *int
fmt.Println(*p) // 5   ← * = "روح للعنوان وهات القيمة"
*p = 9
fmt.Println(x)  // 9   ← x نفسها اتغيرت!
~~~

| | المعنى |
|---|---|
| ~*int~ | نوع: "pointer لـ int" |
| ~&x~ | هات عنوان x |
| ~*p~ | القيمة اللي في العنوان p |

## الحل

~~~go
func addHarvest(total *int, n int) {
	*total += n
}

addHarvest(&total, n)
~~~

## إمتى تستخدم pointer؟
1. لما الدالة **محتاجة تغيّر** متغير بتاع اللي ناداها.
2. لما القيمة **كبيرة** (struct فيه حاجات كتير) ومش عايز تنسخها كل مرة.

غير كده رجّع القيمة عادي بـ ~return~، لأنه غالباً أوضح.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch08CargoEnglish,
			Title: "Cargo Runs",
			Goal:  "Store 30 crystal at base. The drone carries only 10 at a time",
			Lesson: `# Pointer Receivers: State That Stays ⚓

In chapter 7 we used ~*Nav~ without explaining why. Now you will see for yourself.

The drone carries **only 10**. It must track its load, return to base when full, and call ~colony.Unload()~. The code is ready, but one bug makes the mission fail.

## Value receiver = a copy

~~~go
func (d Drone) Load(n int) {
	d.Cargo += n // d is a copy of the drone; this change is lost
}
~~~

## Pointer receiver = the original

~~~go
func (d *Drone) Load(n int) {
	d.Cargo += n // d points to the actual drone
}
~~~

Note: Go handles the conversion automatically. With an ordinary variable (~d := Drone{}~), you can call ~d.Load(3)~ and Go passes ~&d~ for you.

## A practical rule
- If any method changes a struct, use pointer receivers for **all** its methods to keep the code consistent.
- Value receivers suit small types that do not change (such as ~Point~).

## New command: Unload
- ~colony.Unload()~ transfers cargo into storage. It works **only at base (0, 0)**.
- Harvest fails if the crystal will not fit in the cargo hold.`,
			Starter: `package main

import "orbital/colony"

type Drone struct {
	X, Y  int
	Cargo int
}

// Load records that the drone picked up n
func (d Drone) Load(n int) {
	d.Cargo += n
}

func (d *Drone) GoTo(p colony.Point) {
	for ; d.X < p.X; d.X++ {
		colony.Move(colony.East)
	}
	for ; d.X > p.X; d.X-- {
		colony.Move(colony.West)
	}
	for ; d.Y < p.Y; d.Y++ {
		colony.Move(colony.North)
	}
	for ; d.Y > p.Y; d.Y-- {
		colony.Move(colony.South)
	}
}

func main() {
	d := &Drone{}
	base := colony.Point{X: 0, Y: 0}
	for _, t := range colony.Targets() {
		d.GoTo(t)
		size := colony.Measure()
		if d.Cargo+size > 10 {
			d.GoTo(base)
			colony.Unload()
			d.Cargo = 0
			d.GoTo(t)
		}
		colony.Harvest()
		d.Load(size)
	}
}
`,
			Hints: []string{
				"A value receiver changes its copy; a pointer receiver can change the original struct.",
				"Inspect Cargo updates, the capacity check, and the unload timing as separate phases.",
				"Compare the receiver on ~Load~ with the one on ~GoTo~ in the starter.",
			},
		},
		LessonRevision: 1, Steps: ch08CargoArabic,
		ID: "8-2-cargo-runs", Chapter: 8, Title: "رحلات الشحن",
		Goal: "خزّن 30 كريستال في القاعدة. الدرون بيشيل 10 بس في المرة",
		Par:  120, Gold: 80, MaxTicks: 600,
		Requires: []string{"method", "ptrrecv"},
		Setup: func(w *sim.World) {
			w.Init(8, 4, sim.Regolith)
			w.Capacity = 10
			w.At(0, 0).Entity = sim.Beacon
			w.Targets = randomPoints(w, 12)
			total := 0
			for _, p := range w.Targets {
				t := w.At(p.X, p.Y)
				t.Entity, t.Size = sim.Crystal, 2+w.Rand.Intn(5)
				total += t.Size
			}
			for i := 0; total < 30; i++ {
				p := w.Targets[i%len(w.Targets)]
				w.At(p.X, p.Y).Size++
				total++
			}
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Crystal]
				return n >= 30, fmt.Sprintf(i18n.T("في المخزن: %d / 30 · مع الدرون: %d / 10", "In storage: %d / 30 · On the drone: %d / 10"), n, w.Cargo)
			}
		},
		Starter: `package main

import "orbital/colony"

type Drone struct {
	X, Y  int
	Cargo int
}

// Load بتسجّل إن الدرون شال n
func (d Drone) Load(n int) {
	d.Cargo += n
}

func (d *Drone) GoTo(p colony.Point) {
	for ; d.X < p.X; d.X++ {
		colony.Move(colony.East)
	}
	for ; d.X > p.X; d.X-- {
		colony.Move(colony.West)
	}
	for ; d.Y < p.Y; d.Y++ {
		colony.Move(colony.North)
	}
	for ; d.Y > p.Y; d.Y-- {
		colony.Move(colony.South)
	}
}

func main() {
	d := &Drone{}
	base := colony.Point{X: 0, Y: 0}
	for _, t := range colony.Targets() {
		d.GoTo(t)
		size := colony.Measure()
		if d.Cargo+size > 10 {
			d.GoTo(base)
			colony.Unload()
			d.Cargo = 0
			d.GoTo(t)
		}
		colony.Harvest()
		d.Load(size)
	}
}
`,
		Hints: []string{
			"Value receiver بتغيّر نسختها؛ pointer receiver يقدر يغيّر الـstruct الأصلية.",
			"راجع تحديث Cargo، وفحص السعة، ووقت التفريغ كخطوات منفصلة.",
			"قارن receiver بتاعة ~Load~ باللي على ~GoTo~ في الـstarter.",
		},
		Lesson: `# Pointer receivers: الحالة اللي بتفضل ⚓

في الفصل 7 استخدمنا ~*Nav~ من غير ما نفهم ليه. دلوقتي هتشوف بنفسك.

الدرون بيشيل **10 بس**. لازم يعرف هو شايل كام، ولما يتملى يرجع للقاعدة ويعمل ~colony.Unload()~. الكود جاهز، بس فيه bug واحد بيخلّي المهمة تفشل.

## value receiver = نسخة

~~~go
func (d Drone) Load(n int) {
	d.Cargo += n // d نسخة من الدرون. التغيير بيضيع
}
~~~

## pointer receiver = الأصل

~~~go
func (d *Drone) Load(n int) {
	d.Cargo += n // d بتشاور على الدرون الحقيقي
}
~~~

ملحوظة: Go بيحوّل لوحده. لو ~d~ متغير عادي (~d := Drone{}~)، تقدر تكتب ~d.Load(3)~ وGo هيبعت ~&d~ من غير ما تكتبها.

## القاعدة العملية
- لو أي method بتغيّر الـ struct، خلّي **كل** الـ methods بتاعته pointer receivers. ده بيخلّي الكود متناسق.
- value receivers للأنواع الصغيرة اللي مبتتغيرش (زي ~Point~).

## مهارة جديدة: Unload
- ~colony.Unload()~: بتفرّغ الشحنة في المخزن. بتشتغل **في القاعدة (0, 0) بس**.
- Harvest بتفشل لو الكريستال مش هيساع في الشحنة.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch08RadarEnglish,
			Title: "Radar",
			Goal:  "Use colony.Nearest() to find and harvest every target",
			Lesson: `# nil: A Pointer to Nothing 📡

The radar returns ~*colony.Point~: a pointer to the nearest target. When no targets remain, it returns **nil**, meaning "none."

~~~go
p := colony.Nearest()
if p == nil {
	fmt.Println("No targets left")
	return
}
fmt.Println(p.X, p.Y) // Go automatically uses (*p).X
~~~

## A Common Go Crash 💥

~~~go
var p *colony.Point // nil
fmt.Println(*p)
// panic: runtime error: invalid memory address or nil pointer dereference
~~~

Before using a pointer that might be nil, **check it first**.

## A three-part for loop with a pointer

~~~go
for p := colony.Nearest(); p != nil; p = colony.Nearest() {
	// p is definitely not nil here
}
~~~

This is exactly the same structure as ~for i := 0; i < n; i++~: initialization, condition, and step.

## Why does Go use nil?
nil expresses "no value" without a magic number such as -1. You will also see it with ~error~ (~err == nil~ means no error), maps, slices, and interfaces.

## Stars ⭐
Each radar query costs a tick. The radar chooses the nearest target, so your solution automatically follows the greedy approach.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	p := colony.Nearest() // *colony.Point
	fmt.Println("Address:", p, "value:", *p)

	// Repeat: go to *p, harvest, then ask the radar again.
	// The radar returns nil when no targets remain.
}
`,
			Hints: []string{
				"A nil pointer means the radar found no remaining target; check before reading it.",
				"Plan a scan phase, a safe-result phase, and a visit phase; repeat after a harvest.",
				"Inspect the ~Nearest~ call and the ~*p~ print in the starter.",
			},
		},
		LessonRevision: 1, Steps: ch08RadarArabic,
		ID: "8-3-nil-radar", Chapter: 8, Title: "الرادار",
		Goal: "استخدم colony.Nearest() عشان تلاقي الأهداف وتحصدها كلها",
		Par:  80, Gold: 60, MaxTicks: 400,
		Requires: []string{"nilcheck"},
		Setup: func(w *sim.World) {
			w.Init(9, 7, sim.Regolith)
			w.Targets = randomPoints(w, 7)
			placeCrystals(w, w.Targets)
			w.Gold = greedyCost(w.Targets) + len(w.Targets) + len(w.Targets) + 1
			w.Par = w.Gold + 14
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, w.Targets)
				return left == 0, fmt.Sprintf(i18n.T("أهداف فاضلة: %d / %d", "Targets remaining: %d / %d"), left, len(w.Targets))
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	p := colony.Nearest() // *colony.Point
	fmt.Println("العنوان:", p, "القيمة:", *p)

	// كرّر: روح لـ *p، احصد، واسأل الرادار تاني.
	// الرادار بيرجّع nil لما الأهداف تخلص.
}
`,
		Hints: []string{
			"المؤشر nil معناه إن الرادار ملاقاش هدف فاضل؛ افحصه قبل قراءته.",
			"خطط لمرحلة مسح، وفحص النتيجة، وزيارة؛ وكرر بعد الحصاد.",
			"راجع نداء ~Nearest~ وسطر طباعة ~*p~ في الـstarter.",
		},
		Lesson: `# nil: الـ pointer اللي مبيشاورش على حاجة 📡

الرادار بيرجّع ~*colony.Point~: عنوان لأقرب هدف. ولما الأهداف تخلص بيرجّع **nil**، يعني "مفيش".

~~~go
p := colony.Nearest()
if p == nil {
	fmt.Println("خلاص، مفيش أهداف")
	return
}
fmt.Println(p.X, p.Y) // Go بيعمل (*p).X لوحده
~~~

## أشهر crash في Go 💥

~~~go
var p *colony.Point // nil
fmt.Println(*p)
// panic: runtime error: invalid memory address or nil pointer dereference
~~~

قبل ما تستخدم أي pointer ممكن يكون nil، **اتأكد منه الأول**.

## for بثلاث أجزاء مع pointer

~~~go
for p := colony.Nearest(); p != nil; p = colony.Nearest() {
	// p أكيد مش nil هنا
}
~~~

ده نفس شكل ~for i := 0; i < n; i++~ بالظبط: بداية، وشرط، وخطوة.

## ليه Go بيستخدم nil كده؟
nil بيعبّر عن "مفيش قيمة" من غير ما تحتاج قيمة سحرية زي -1. وهتلاقيه كتير مع ~error~ (لو ~err == nil~ يبقى مفيش خطأ)، ومع الـ maps والـ slices والـ interfaces.

## النجوم ⭐
كل سؤال للرادار بياخد tick. الرادار نفسه بيختار الأقرب، يعني حلّك بيمشي بالـ greedy تلقائياً.`,
	})
	add(&Level{ID: "8-4-cart-chain", Chapter: 8, Kind: "lab", Title: "سلسلة العربات", Goal: "وصّل وافصل عربات من سلسلة من غير ما تضيّع باقي القطر", Requires: []string{"pointer", "struct"}, Bridge: []int{29},
		LessonRevision: 1, Steps: ch08CartArabic,
		Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: cartChainStarter, ReadOnly: false}},
		Hints: []string{
			"كل node بتوصل للي بعدها بعنوان، و~nil~ بتحدد نهاية السلسلة.",
			"افحص الحالة الفاضية، وبعدها التوصيل والفصل والبحث كل واحدة لوحدها.",
			"تابع حقلي ~Head~ و~Next~ في struct قبل تعديل أي وصلة.",
		},
		Lesson: `# عربات ماسكة في بعض

كل عربة شايلة اسمها وعنوان اللي بعدها. القطر مش محتاج generics ولا recursion: هنمشي على الوصلات بحلقة عادية، باستخدام pointers اللي عرفناها في الفصل ده.

## قبل ما تكتب
- مين لازم يتغيّر لما تدخل عربة قبل أول واحدة؟
- هتحتفظ بعنوان باقي السلسلة إزاي وإنت بتغيّر وصلة؟

## عقدة ووصلة
~~~go
type Lamp struct { Label string; Next *Lamp }
last := &Lamp{Label: "blue"}
first := &Lamp{Label: "red", Next: last}
fmt.Println(first.Next.Label) // blue
~~~
~Next *Lamp~ ممكن يبقى nil: ده آخر السلسلة. ~&Lamp{...}~ بيعمل قيمة وياخد عنوانها؛ الـ pointer بيخلّي عقدة تشير لعقدة تانية من نفس النوع من غير حجم لا نهائي.

## امشي من غير ما تغيّر البداية
~~~go
for p := first; p != nil; p = p.Next {
 fmt.Println(p.Label)
}
~~~
~p~ مؤشر مؤقت. تغيير p مش بيغيّر first. افحص nil قبل ما تقرأ field؛ القايمة الفاضية بدايتها nil.

## غيّر الوصلة: مثال مختلف
~~~go
middle := &Lamp{Label: "green", Next: first.Next}
first.Next = middle
~~~
دالة بناء صغيرة بتاخد الوصلة القديمة كـ pointer parameter:
~~~go
func linkedLamp(next *Lamp,label string) *Lamp {
 return &Lamp{Label:label,Next:next}
}
extra:=linkedLamp(last,"white")
fmt.Println(extra.Next.Label) // blue
~~~
استخدم نفس فكرة helper للعربات عند إدخال البداية أو الوسط، عشان دالة البناء تستقبل عنوان الجزء اللي هتكمل عليه.

خزّن الوصلة القديمة في العقدة الجديدة قبل ما توصّلها. للحذف، العقدة السابقة بتشير للي بعد المحذوفة؛ حذف البداية حالة خاصة لأن Head نفسها بتتغيّر.

## شكل الاستدعاء
~~~go
var convoy Convoy
err:=convoy.Couple(0,"water")
at:=convoy.Locate("water")
tags:=convoy.Tags()
removed:=convoy.Release("water")
fmt.Println(err,at,tags,removed) // <nil> 0 [water] true
~~~
دي تجربة للعقد؛ تنفيذ الوصلات هو شغلك. ~ErrCoupling~ sentinel error موجودة في starter؛ رجّعها زي ما هي للموقع البايظ عشان errors.Is تقدر تتعرف عليها.

## أخطاء شائعة
~~~go
first.Next = &Lamp{Label: "green"} // باقي السلسلة ضاع
fmt.Println(last.Next.Label)      // nil: panic
~~~
الـ slice اللي بترجعها لازم تكون جديدة، عشان تعديلها ميغيّرش بيانات العربات. شغلنا هنا تسلسلي؛ الوصول المتزامن مش جزء من العقد.

## المهمة
نفّذ ~Convoy~ من عقد ~Cart{Tag string; Next *Cart}~ وبداية ~Head *Cart~؛ الـ zero value فاضية وشغالة. ~Couple(at, tag)~ تدخل قبل الموقع at من0 لحد الطول شامل الطرفين؛ أي موقع تاني يرجّع ~ErrCoupling~ من غير تعديل.

~Release(tag) bool~ تحذف أول تطابق فقط؛ false لو مفيش. ~Locate(tag) int~ ترجع موقع أول تطابق من صفر، أو -1. ~Tags() []string~ ترجع نسخة مستقلة بالترتيب، وفاضية لو القطر فاضي. مثال: بعد Couple(0,"a") وCouple(1,"b")، Tags نتيجتها [a b]. مفيش Report في المختبر؛ التستات بتراجع العقد مباشرة.
`,
		En: &LevelText{LessonRevision: 1, Steps: ch08CartEnglish, Title: "Cart Chain", Goal: "Couple and release carts without losing the rest of the convoy", Files: []LabFile{{Name: "main.go", Starter: cartChainStarter, ReadOnly: false}}, Hints: []string{
			"Each node links to the next by address, and ~nil~ marks the end.",
			"Inspect the empty case, then linking, release, and search as separate phases.",
			"Follow the ~Head~ and ~Next~ fields before changing any link.",
		}, Lesson: `# Carts Linked Together

Each cart stores its tag and the address of the next cart. No generics or recursion are needed: follow pointers with an ordinary loop, using this chapter's pointer tools.

## Before you write
- What must change when inserting before the first cart?
- How will you preserve the remaining chain while changing a link?

## A node and a link
~~~go
type Lamp struct { Label string; Next *Lamp }
last := &Lamp{Label: "blue"}
first := &Lamp{Label: "red", Next: last}
fmt.Println(first.Next.Label) // blue
~~~
~Next *Lamp~ may be nil, marking the end. ~&Lamp{...}~ creates a value and takes its address. A pointer lets a node refer to its own type without infinite size.

## Walk without changing the beginning
~~~go
for p := first; p != nil; p = p.Next {
 fmt.Println(p.Label)
}
~~~
~p~ is a temporary pointer. Reassigning it does not change first. Check nil before reading a field; an empty list begins with nil.

## Change a link: a different example
~~~go
middle := &Lamp{Label: "green", Next: first.Next}
first.Next = middle
~~~
A small construction function receives the old link as a pointer parameter:
~~~go
func linkedLamp(next *Lamp,label string) *Lamp {
 return &Lamp{Label:label,Next:next}
}
extra:=linkedLamp(last,"white")
fmt.Println(extra.Next.Label) // blue
~~~
Use the same helper idea for cart insertion at the head or middle; the pointer parameter receives the address of the remaining chain.

Save the old link in the new node before attaching it. To delete, the preceding node points past the removed node. Removing the first node is special because Head itself must change.

## Calling the contract
~~~go
var convoy Convoy
err:=convoy.Couple(0,"water")
at:=convoy.Locate("water")
tags:=convoy.Tags()
removed:=convoy.Release("water")
fmt.Println(err,at,tags,removed) // <nil> 0 [water] true
~~~
This exercises the contract; implementing links is your work. ~ErrCoupling~ is the starter's sentinel error. Return it unchanged for invalid positions so errors.Is can recognize it.

## Common mistakes
~~~go
first.Next = &Lamp{Label: "green"} // Loses the remaining chain
fmt.Println(last.Next.Label)      // nil: panic
~~~
Return a fresh slice so changing it cannot change cart data. This contract is sequential; concurrent access is outside its scope.

## Your mission
Implement ~Convoy~ using ~Cart{Tag string; Next *Cart}~ nodes and ~Head *Cart~. Its zero value is usable and empty. ~Couple(at, tag)~ inserts before zero-based position at, accepting 0 through length inclusive. Other positions return ~ErrCoupling~ without mutation.

~Release(tag) bool~ deletes only the first match, returning false for a miss. ~Locate(tag) int~ returns the first zero-based position or -1. ~Tags() []string~ returns an independent ordered snapshot, empty for an empty convoy. For example, Couple(0,"a") then Couple(1,"b") yields [a b]. There is no Report in this lab; tests call the contract directly.
`},
		Lab: &Lab{ReqFile: "main.go", Hidden: cartChainHidden},
	})
}

const cartChainStarter = `package main

import "errors"

var ErrCoupling = errors.New("invalid coupling position")

type Cart struct {
	Tag  string
	Next *Cart
}
type Convoy struct{ Head *Cart }

func (c *Convoy) Couple(at int, tag string) error { return ErrCoupling }
func (c *Convoy) Release(tag string) bool         { return false }
func (c *Convoy) Locate(tag string) int           { return -1 }
func (c *Convoy) Tags() []string                  { return nil }
func main()                                       {}
`
const cartChainHidden = `package main

import (
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"slices"
	"testing"
)

func TestHiddenCouplings(t *testing.T) {
	var convoy Convoy
	if convoy.Head != nil || len(convoy.Tags()) != 0 || convoy.Release("none") || convoy.Locate("none") != -1 {
		t.Fatal("empty convoy")
	}
	world := levels.Start("8-4-cart-chain", io.Discard, 829)
	r := world.Rand
	var want []string
	for step := 0; step < 250; step++ {
		at := r.Intn(len(want) + 1)
		tag := fmt.Sprint(world.Data.SensorReadings[step] % 12)
		if err := convoy.Couple(at, tag); err != nil {
			t.Fatal(err)
		}
		want = append(want, "")
		copy(want[at+1:], want[at:])
		want[at] = tag
		for _, bad := range []int{-1, len(want) + 1} {
			if !errors.Is(convoy.Couple(bad, "bad"), ErrCoupling) {
				t.Fatal("invalid position must match ErrCoupling")
			}
		}
		if !slices.Equal(convoy.Tags(), want) || convoy.Locate(tag) != slices.Index(want, tag) {
			t.Fatal("order or first match")
		}
		snapshot := convoy.Tags()
		snapshot[0] = "changed"
		if !slices.Equal(convoy.Tags(), want) {
			t.Fatal("Tags aliases nodes")
		}
		if step%2 == 0 {
			if !convoy.Release(tag) {
				t.Fatal("release existing")
			}
			i := slices.Index(want, tag)
			want = append(want[:i], want[i+1:]...)
		}
		if convoy.Release("missing") || !slices.Equal(convoy.Tags(), want) {
			t.Fatal("missing release mutated convoy")
		}
	}
	for len(want) > 0 {
		tag := want[len(want)-1]
		i := slices.Index(want, tag)
		if !convoy.Release(tag) {
			t.Fatal("drain")
		}
		want = append(want[:i], want[i+1:]...)
		if !slices.Equal(convoy.Tags(), want) {
			t.Fatal("drain order")
		}
	}
	if convoy.Head != nil {
		t.Fatal("last cart must clear head")
	}
}
`
