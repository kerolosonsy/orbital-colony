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
			Title: "A Well-Formatted Report",
			Goal:  "Make Report(s) produce exactly: moss=12 crystal=7",
			Lesson: `# Interfaces: An Implicit Agreement 🤝

An **interface** describes **behavior**: "anything with these methods." What a value does matters, regardless of its concrete type.

~~~go
// Defined in the fmt package:
type Stringer interface {
	String() string
}
~~~

Any type with a ~String() string~ method automatically satisfies **Stringer**. You do not write "implements" as in Java. This is **implicit implementation**, one of Go's most useful features.

~fmt.Println~ and ~colony.Report~ check whether the value is a Stringer. If so, they use your ~String()~ method.

~~~go
func (s Status) String() string {
	return fmt.Sprintf("moss=%d crystal=%d", s.Moss, s.Crystal)
}
~~~

## Sprintf: formatting text
| | |
|---|---|
| ~%d~ | An integer |
| ~%s~ | A string |
| ~%v~ | Any value in its default format |
| ~%+v~ | A struct with field names: ~{Moss:12 Crystal:7}~ |
| ~%.2f~ | A number with two decimal places |

## Defining your own interface

~~~go
type Harvester interface {
	Harvest() int
}
~~~

Any type with ~Harvest() int~ can be passed to a function accepting ~Harvester~. This lets your code work with types **that have not been written yet**.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

type Status struct {
	Moss, Crystal int
}

func main() {
	s := Status{
		Moss:    colony.Count(colony.Moss),
		Crystal: colony.Count(colony.Crystal),
	}
	fmt.Println(s) // Prints {12 7}... not the required format
	colony.Report(s)
}
`,
			Hints: []string{
				"Add a method named exactly ~String~ with no arguments and a ~string~ return value.",
				"Inside it: ~return fmt.Sprintf(\"moss=%d crystal=%d\", s.Moss, s.Crystal)~.",
				"Keep main as it is: ~colony.Report(s)~. Your String method will be called automatically.",
			},
		},
		ID: "9-1-stringer", Chapter: 9, Title: "تقرير بشكل محترم",
		Goal: "خلّي Report(s) يطلع بالشكل ده بالظبط: moss=12 crystal=7",
		Par:  0, Gold: 0, MaxTicks: 20,
		Requires: []string{"stringer"},
		Setup: func(w *sim.World) {
			w.Init(4, 3, sim.Regolith)
			w.Inv[sim.Moss] = 5 + w.Rand.Intn(40)
			w.Inv[sim.Crystal] = 5 + w.Rand.Intn(40)
			want := fmt.Sprintf("moss=%d crystal=%d", w.Inv[sim.Moss], w.Inv[sim.Crystal])
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("التقرير وصل بالشكل الصح", "The report has the correct format")
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت \"", "Base received \"") + w.Reports[n-1] + i18n.T("\" — المطلوب \"", "\" — expected \"") + want + "\""
				}
				return false, i18n.T("القاعدة مستنية التقرير", "Base is waiting for your report")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

type Status struct {
	Moss, Crystal int
}

func main() {
	s := Status{
		Moss:    colony.Count(colony.Moss),
		Crystal: colony.Count(colony.Crystal),
	}
	fmt.Println(s) // بيطبع {12 7}… مش الشكل المطلوب
	colony.Report(s)
}
`,
		Hints: []string{
			"ضيف method اسمها ~String~ بالظبط، مبتاخدش حاجة وبترجّع ~string~.",
			"جوّاها: ~return fmt.Sprintf(\"moss=%d crystal=%d\", s.Moss, s.Crystal)~.",
			"خلّي main زي ما هي: ~colony.Report(s)~. الـ String بتاعتك هي اللي هتشتغل لوحدها.",
		},
		Lesson: `# Interfaces: اتفاق من غير توقيع 🤝

الـ **interface** نوع بيوصف **سلوك**: "أي حاجة عندها الـ methods دي". مش مهم هي إيه، المهم تعمل إيه.

~~~go
// موجود في مكتبة fmt:
type Stringer interface {
	String() string
}
~~~

أي نوع عنده method اسمها ~String() string~ بيبقى **Stringer** تلقائياً. مش محتاج تكتب "implements" زي Java. ده اسمه **implicit implementation**، وهو من أقوى حاجات Go.

~fmt.Println~ و~colony.Report~ بيسألوا: "القيمة دي Stringer؟" لو أيوه، بيستخدموا ~String()~ بتاعتك.

~~~go
func (s Status) String() string {
	return fmt.Sprintf("moss=%d crystal=%d", s.Moss, s.Crystal)
}
~~~

## Sprintf: تنسيق النصوص
| | |
|---|---|
| ~%d~ | رقم صحيح |
| ~%s~ | نص |
| ~%v~ | أي قيمة بالشكل الافتراضي |
| ~%+v~ | struct بأسماء الحقول: ~{Moss:12 Crystal:7}~ |
| ~%.2f~ | عشري برقمين بعد العلامة |

## تعريف interface بتاعك

~~~go
type Harvester interface {
	Harvest() int
}
~~~

أي نوع عنده ~Harvest() int~ يقدر يتبعت لدالة بتستقبل ~Harvester~. ده بيخلّي كودك يشتغل مع أنواع **لسه متكتبتش**.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Autopilot",
			Goal:  "Write a Planner for colony.Autopilot: harvest ripe plants and choose the right crop for each tile",
			Lesson: `# You Write It, the Library Calls It 🛩️

Until now, **you** have called library functions. Here it works the other way around: ~colony.Autopilot~ visits the map and **calls your code** on each tile to ask, "What should I do?"

~~~go
// In the colony library:
type Planner interface {
	Plan(t Tile) []Action
}

func Autopilot(p Planner)
~~~

You can give the autopilot any type with ~Plan(colony.Tile) []colony.Action~:

~~~go
type Farmer struct{}

func (f Farmer) Plan(t colony.Tile) []colony.Action {
	if t.Ground == colony.Ice {
		return []colony.Action{colony.PlantLichen}
	}
	return nil // Nothing to do
}

colony.Autopilot(Farmer{})
~~~

## This is how real Go programs work
- Servers: you write an ~http.Handler~, and the server calls your ~ServeHTTP~ for each request.
- Sorting: ~sort.Interface~.
- Reading and writing: ~io.Reader~ and ~io.Writer~, two of Go's most important interfaces.

## Available actions
~DoHarvest~, ~DoTill~, ~PlantMoss~, ~PlantLichen~. The autopilot executes them in slice order.

## The Tile
~t.X~, ~t.Y~, ~t.Ground~, ~t.Entity~, ~t.Ripe~. Reading the tile is free; actions cost ticks as usual.`,
			Starter: `package main

import "orbital/colony"

// Farmer will be our Planner
type Farmer struct{}

// Plan tells the autopilot what to do on each tile
func (f Farmer) Plan(t colony.Tile) []colony.Action {
	switch t.Ground {
	case colony.Ice:
		return []colony.Action{colony.PlantLichen}
	}
	return nil
}

func main() {
	colony.Autopilot(Farmer{})
}
`,
			Hints: []string{
				"Regolith needs ~[]colony.Action{colony.DoTill, colony.PlantMoss}~.",
				"On soil: if ~t.Entity == colony.Moss && t.Ripe~, return Harvest then PlantMoss. If empty, return only PlantMoss.",
				"The autopilot makes one pass. A wrong action wastes a tick.",
			},
		},
		ID: "9-2-autopilot", Chapter: 9, Title: "الطيار الآلي",
		Goal: "اكتب Planner لـ colony.Autopilot: احصد الناضج، وازرع كل خانة باللي يناسبها",
		Par:  80, Gold: 60, MaxTicks: 400,
		Requires: []string{"autopilot", "method"},
		Setup: func(w *sim.World) {
			w.Init(6, 4, sim.Soil)
			grounds := []sim.Ground{sim.Soil, sim.Soil, sim.Ice, sim.Regolith, sim.Rock}
			actions, ripe := 0, 0
			w.Each(func(x, y int, t *sim.Tile) {
				t.Ground = grounds[w.Rand.Intn(len(grounds))]
				switch t.Ground {
				case sim.Soil:
					actions++
					if w.Rand.Intn(2) == 0 {
						w.PlantAt(x, y, sim.Moss, 0)
						actions++
						ripe++
					}
				case sim.Ice, sim.Regolith:
					actions += map[sim.Ground]int{sim.Ice: 1, sim.Regolith: 2}[t.Ground]
				}
			})
			w.Gold = w.W*w.H - 1 + actions
			w.Par = w.Gold + 12
			w.Check = func(w *sim.World) (bool, string) {
				done, total := 0, 0
				w.Each(func(_, _ int, t *sim.Tile) {
					switch {
					case t.Ground == sim.Rock:
						return
					case t.Ground == sim.Ice && t.Entity == sim.Lichen,
						t.Ground == sim.Soil && t.Entity == sim.Moss && t.Planted > 0:
						done++
					}
					total++
				})
				ok := done == total && w.Inv[sim.Moss] >= ripe
				return ok, fmt.Sprintf(i18n.T("خانات جاهزة: %d / %d · طحالب محصودة: %d / %d", "Ready tiles: %d / %d · Moss harvested: %d / %d"), done, total, w.Inv[sim.Moss], ripe)
			}
		},
		Starter: `package main

import "orbital/colony"

// Farmer هيبقى الـ Planner بتاعنا
type Farmer struct{}

// Plan بيقول للطيار الآلي يعمل إيه في كل خانة
func (f Farmer) Plan(t colony.Tile) []colony.Action {
	switch t.Ground {
	case colony.Ice:
		return []colony.Action{colony.PlantLichen}
	}
	return nil
}

func main() {
	colony.Autopilot(Farmer{})
}
`,
		Hints: []string{
			"الـ regolith محتاج ~[]colony.Action{colony.DoTill, colony.PlantMoss}~.",
			"على soil: لو ~t.Entity == colony.Moss && t.Ripe~ ارجّع Harvest ثم PlantMoss. ولو فاضية PlantMoss بس.",
			"الطيار الآلي بيلف مرة واحدة. لو رجّعت action غلط، هيضيّع tick.",
		},
		Lesson: `# إنت بتكتب، والمكتبة بتنادي 🛩️

لحد دلوقتي كنت **إنت** اللي بتنادي دوال المكتبة. هنا العكس: ~colony.Autopilot~ هي اللي بتلف على الخريطة، وفي كل خانة **بتنادي كودك** وتسأله: "أعمل إيه؟"

~~~go
// في مكتبة colony:
type Planner interface {
	Plan(t Tile) []Action
}

func Autopilot(p Planner)
~~~

أي نوع عنده ~Plan(colony.Tile) []colony.Action~ تقدر تديه للطيار الآلي:

~~~go
type Farmer struct{}

func (f Farmer) Plan(t colony.Tile) []colony.Action {
	if t.Ground == colony.Ice {
		return []colony.Action{colony.PlantLichen}
	}
	return nil // مفيش حاجة تتعمل
}

colony.Autopilot(Farmer{})
~~~

## ده بالظبط إزاي Go الحقيقي شغال
- السيرفرات: إنت بتكتب ~http.Handler~، والسيرفر بينادي ~ServeHTTP~ بتاعتك مع كل request.
- الترتيب: ~sort.Interface~.
- القراءة والكتابة: ~io.Reader~ و~io.Writer~. دول أهم interfaces في Go كلها.

## الـ Actions المتاحة
~DoHarvest~، ~DoTill~، ~PlantMoss~، ~PlantLichen~. الطيار بينفّذهم بالترتيب اللي في الـ slice.

## الـ Tile
~t.X~، ~t.Y~، ~t.Ground~، ~t.Entity~، ~t.Ripe~. قراية الخانة ببلاش، والـ actions بتاخد ticks عادي.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Radio Signals",
			Goal:  "Respond to every Distress by visiting its location; add the Supply amounts and report the total",
			Lesson: `# Type Switch: What Is Inside the Interface? 📻

~colony.Signals()~ returns ~[]colony.Signal~. ~Signal~ is an interface; each element holds a concrete value of type ~Distress~, ~Supply~, or ~Noise~.

## Interface value = type + value
An interface value holds two things: the **concrete type** and the **value**. ~%T~ prints the type:

~~~go
fmt.Printf("%T\n", s) // sim.Supply
~~~

## Type switch

~~~go
switch v := s.(type) {
case colony.Distress:
	fmt.Println("Distress at", v.At)    // v has type Distress here
case colony.Supply:
	total += v.Amount                   // v has type Supply here
case colony.Noise:
	// Ignore
default:
	fmt.Println("Unknown type")
}
~~~

## Type assertion: ask one question

~~~go
sup, ok := s.(colony.Supply)
if ok {
	fmt.Println(sup.Amount)
}
~~~

Without ~ok~ (~sup := s.(colony.Supply)~), a wrong type makes the program **panic**.

## any
~any~ is an alias for ~interface{}~: an interface with no methods, satisfied by every Go value. That is why ~colony.Report(v any)~ accepts anything. But each use of ~any~ gives up some compiler help with types, so use it thoughtfully.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	for _, s := range colony.Signals() {
		fmt.Printf("%T %v\n", s, s)
	}
}
`,
			Hints: []string{
				"~switch v := s.(type) { case colony.Distress: ... case colony.Supply: ... }~. Inside each case, ~v~ has that case's type.",
				"Add ~v.Amount~ to a variable and save ~v.At~ in a slice. After the loop, visit the points and report the total.",
				"For more stars, visit distress calls nearest first (greedy), rather than in signal order.",
			},
		},
		ID: "9-3-signals", Chapter: 9, Title: "إشارات الراديو",
		Goal: "رد على كل Distress بإنك تروح مكانها، واجمع قيم الـ Supply وابعت المجموع",
		Par:  60, Gold: 40, MaxTicks: 400,
		Requires: []string{"typeswitch"},
		Setup: func(w *sim.World) {
			w.Init(9, 7, sim.Regolith)
			pts := randomPoints(w, 4)
			w.Waypoints, w.AnyOrder = pts, true
			sum := 0
			for _, p := range pts {
				w.At(p.X, p.Y).Entity = sim.Beacon
				w.Signals = append(w.Signals, sim.Distress{At: p})
			}
			for i := 0; i < 3; i++ {
				a := 5 + w.Rand.Intn(46)
				sum += a
				w.Signals = append(w.Signals, sim.Supply{Amount: a}, sim.Noise{Text: "kssshhh"})
			}
			w.Rand.Shuffle(len(w.Signals), func(i, j int) { w.Signals[i], w.Signals[j] = w.Signals[j], w.Signals[i] })
			var ordered []sim.Point
			for _, s := range w.Signals {
				if d, ok := s.(sim.Distress); ok {
					ordered = append(ordered, d.At)
				}
			}
			w.Gold = greedyCost(pts)
			w.Par = max(routeCost(ordered), w.Gold+6)
			want := strconv.Itoa(sum)
			w.Check = func(w *sim.World) (bool, string) {
				visited, rep := w.NextWaypoint, reported(w, want)
				msg := fmt.Sprintf(i18n.T("استغاثات: %d / %d", "Distress calls: %d / %d"), visited, len(w.Waypoints))
				if rep {
					msg += i18n.T(" · المجموع وصل ✓", " · Total received ✓")
				} else {
					msg += i18n.T(" · القاعدة مستنية مجموع الـ Supply", " · Base is waiting for the Supply total")
				}
				return visited == len(w.Waypoints) && rep, msg
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	for _, s := range colony.Signals() {
		fmt.Printf("%T %v\n", s, s)
	}
}
`,
		Hints: []string{
			"~switch v := s.(type) { case colony.Distress: ... case colony.Supply: ... }~. جوه كل case، ~v~ نوعها هو نوع الحالة دي.",
			"اجمع ~v.Amount~ في متغير، واحفظ ~v.At~ في slice. وبعد الحلقة روح للنقط وابعت المجموع.",
			"للنجوم: روح للاستغاثات بترتيب الأقرب فالأقرب (greedy)، مش بترتيب الإشارات.",
		},
		Lesson: `# Type switch: إيه اللي جوه الـ interface؟ 📻

~colony.Signals()~ بترجّع ~[]colony.Signal~. ~Signal~ ده interface، وكل عنصر جوّاه نوع حقيقي مختلف: ~Distress~ أو ~Supply~ أو ~Noise~.

## interface value = نوع + قيمة
متغير من نوع interface شايل حاجتين: **النوع الحقيقي** و**القيمة**. ~%T~ بتطبع النوع:

~~~go
fmt.Printf("%T\n", s) // sim.Supply
~~~

## Type switch

~~~go
switch v := s.(type) {
case colony.Distress:
	fmt.Println("استغاثة عند", v.At)    // v هنا نوعه Distress
case colony.Supply:
	total += v.Amount                   // v هنا نوعه Supply
case colony.Noise:
	// تجاهل
default:
	fmt.Println("نوع غريب")
}
~~~

## Type assertion: سؤال واحد

~~~go
sup, ok := s.(colony.Supply)
if ok {
	fmt.Println(sup.Amount)
}
~~~

من غير ~ok~ (~sup := s.(colony.Supply)~)، لو النوع مطلعش صح البرنامج **بيقع**.

## any
~any~ هو ~interface{}~: interface مفيهوش ولا method، فأي قيمة في Go بتحققه. عشان كده ~colony.Report(v any)~ بتقبل أي حاجة. بس كل مرة تستخدم ~any~، بتخسر مساعدة الـ compiler في الأنواع، فاستخدمه بحساب.`,
	})
}
