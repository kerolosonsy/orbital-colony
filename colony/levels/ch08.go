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
				"Run the starter first and inspect the printed total. Why is it zero?",
				"Change the parameter to ~total *int~ and use ~*total += n~ inside the function.",
				"Pass the address when calling it: ~addHarvest(&total, n)~.",
			},
		},
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
			"شغّل الكود زي ما هو الأول وشوف المجموع اللي اتطبع. ليه صفر؟",
			"غيّر الـ parameter لـ ~total *int~، وجوه الدالة ~*total += n~.",
			"وإنت بتنادي الدالة ابعت العنوان: ~addHarvest(&total, n)~.",
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
				"Run it and inspect the Console: \"Not enough cargo space.\" The drone does not realize it is full. Why does ~d.Cargo~ stay zero?",
				"Compare ~func (d Drone) Load~ with ~func (d *Drone) GoTo~. One works on a copy.",
				"For more stars: return to base only when full, start with crystal near base, and remember the final Unload!",
			},
		},
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
			"شغّله وبص على الـ Console: \"الشحنة مش هتساع\". يعني الدرون مبيعرفش إنه اتملى. ليه ~d.Cargo~ فاضلة صفر؟",
			"قارن ~func (d Drone) Load~ بـ ~func (d *Drone) GoTo~. واحدة فيهم بتشتغل على نسخة.",
			"للنجوم: مترجعش للقاعدة غير لما تتملى، وابدأ بالكريستال القريب من القاعدة. وفي الآخر متنساش Unload!",
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
				"Use: ~for p := colony.Nearest(); p != nil; p = colony.Nearest() { ... }~",
				"Inside the loop: ~goTo(*p)~, then ~colony.Harvest()~. You need goTo from chapter 7.",
				"Try dereferencing ~*p~ when it is nil to see the panic for yourself. Once is enough 😅",
			},
		},
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
			"الشكل: ~for p := colony.Nearest(); p != nil; p = colony.Nearest() { ... }~",
			"جوه الحلقة: ~goTo(*p)~ ثم ~colony.Harvest()~. محتاج goTo من الفصل 7.",
			"جرّب تعمل ~*p~ وهي nil عشان تشوف الـ panic بعينك. مرة واحدة كفاية 😅",
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
}
