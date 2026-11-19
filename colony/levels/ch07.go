package levels

import (
	"fmt"
	"slices"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

// randomPoints returns n distinct points on the map, avoiding (0, 0).
func randomPoints(w *sim.World, n int) []sim.Point {
	seen := map[sim.Point]bool{{X: 0, Y: 0}: true}
	var ps []sim.Point
	for len(ps) < n {
		p := sim.Point{X: w.Rand.Intn(w.W), Y: w.Rand.Intn(w.H)}
		if !seen[p] {
			seen[p] = true
			ps = append(ps, p)
		}
	}
	return ps
}

// routeCost is the number of moves to visit ps in order starting at (0, 0).
func routeCost(ps []sim.Point) int {
	cost, at := 0, sim.Point{}
	for _, p := range ps {
		cost += sim.Dist(at, p)
		at = p
	}
	return cost
}

// greedyCost always moves to the nearest unvisited point.
func greedyCost(ps []sim.Point) int {
	left := append([]sim.Point(nil), ps...)
	cost, at := 0, sim.Point{}
	for len(left) > 0 {
		best := 0
		for i, p := range left {
			if sim.Dist(at, p) < sim.Dist(at, left[best]) {
				best = i
			}
		}
		cost += sim.Dist(at, left[best])
		at = left[best]
		left = append(left[:best], left[best+1:]...)
	}
	return cost
}

func placeCrystals(w *sim.World, ps []sim.Point) {
	for _, p := range ps {
		t := w.At(p.X, p.Y)
		t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
	}
}

func remaining(w *sim.World, ps []sim.Point) int {
	n := 0
	for _, p := range ps {
		if w.At(p.X, p.Y).Entity == sim.Crystal {
			n++
		}
	}
	return n
}

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch07WaypointsEnglish,
			Title: "Checkpoints",
			Goal:  "Visit the checkpoints 📡 in the order given by colony.Waypoints()",
			Lesson: `# Structs: Values with Parts 🧱

A map point has **x** and **y** coordinates. Instead of carrying two variables everywhere, Go lets you define a **new type** that groups them:

~~~go
type Point struct {
	X, Y int
}
~~~

This is exactly how ~colony.Point~ is defined. Use it like this:

~~~go
p := colony.Point{X: 3, Y: 2}  // struct literal
fmt.Println(p.X)               // 3
p.Y = 5                        // Change a field
fmt.Println(p)                 // {3 5}
~~~

## A slice of structs

~~~go
points := colony.Waypoints() // []colony.Point
for i, p := range points {
	fmt.Println(i, p.X, p.Y)
}
~~~

## Comparing structs
You can compare points with ~==~ when all their fields are comparable:

~~~go
x, y := colony.Position()
if (colony.Point{X: x, Y: y}) == target {
	// We have arrived
}
~~~

## Your mission
Write ~goTo(p colony.Point)~ to reach any point, then visit ~colony.Waypoints()~ in order. The checkpoints change each run, and the star target is based on the shortest route.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// goTo moves the drone to point p
func goTo(p colony.Point) {
	x, y := colony.Position()
	for x < p.X {
		colony.Move(colony.East)
		x++
	}
	for y < p.Y {
		colony.Move(colony.North)
		y++
	}
	// Complete the other directions...
}

func main() {
	points := colony.Waypoints()
	fmt.Println(points)
	goTo(points[0])
}
`,
			Hints: []string{
				"A point groups coordinates, and the list keeps visit order.",
				"Separate travel to one point from visiting the list.",
				"Finish ~goTo~ after ~colony.Position()~ and check coordinate updates.",
			},
		},
		LessonRevision: 1, Steps: ch07WaypointsArabic,
		ID: "7-1-waypoints", Chapter: 7, Title: "نقاط التفتيش",
		Goal: "عدّي على نقاط التفتيش 📡 بالترتيب اللي في colony.Waypoints()",
		Par:  40, Gold: 30, MaxTicks: 200,
		Requires: []string{"range"},
		Setup: func(w *sim.World) {
			w.Init(8, 6, sim.Regolith)
			w.Waypoints = randomPoints(w, 4)
			for _, p := range w.Waypoints {
				w.At(p.X, p.Y).Entity = sim.Beacon
			}
			w.Gold = routeCost(w.Waypoints)
			w.Par = w.Gold + 8
			w.Check = func(w *sim.World) (bool, string) {
				n := w.NextWaypoint
				return n == len(w.Waypoints), fmt.Sprintf(i18n.T("نقاط: %d / %d — اللي عليها الدور: %v", "Checkpoints: %d / %d — next: %v"), n, len(w.Waypoints), w.Waypoints[min(n, len(w.Waypoints)-1)])
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// goTo تحرّك الدرون لحد النقطة p
func goTo(p colony.Point) {
	x, y := colony.Position()
	for x < p.X {
		colony.Move(colony.East)
		x++
	}
	for y < p.Y {
		colony.Move(colony.North)
		y++
	}
	// كمّل باقي الاتجاهات...
}

func main() {
	points := colony.Waypoints()
	fmt.Println(points)
	goTo(points[0])
}
`,
		Hints: []string{
			"الـstruct تجمع الإحداثيين، والـslice بتحفظ ترتيب المرور.",
			"افصل بين الوصول لنقطة وتكرار المرور على النقاط.",
			"كمّل ~goTo~ بعد ~colony.Position()~ وراجع تحديث الإحداثيات.",
		},
		Lesson: `# Structs: حاجات ليها أجزاء 🧱

النقطة على الخريطة ليها **x** و **y**. بدل ما تشيل متغيرين في كل حتة، Go بيخليك تعمل **نوع جديد** يجمعهم:

~~~go
type Point struct {
	X, Y int
}
~~~

ده بالظبط تعريف ~colony.Point~. استخدامه:

~~~go
p := colony.Point{X: 3, Y: 2}  // struct literal
fmt.Println(p.X)               // 3
p.Y = 5                        // غيّر حقل
fmt.Println(p)                 // {3 5}
~~~

## slice من structs

~~~go
points := colony.Waypoints() // []colony.Point
for i, p := range points {
	fmt.Println(i, p.X, p.Y)
}
~~~

## قارن structs
تقدر تقارن نقطتين بـ ~==~ لو كل الحقول بتاعتهم ينفع تتقارن:

~~~go
x, y := colony.Position()
if (colony.Point{X: x, Y: y}) == target {
	// وصلت
}
~~~

## المهمة
اكتب ~goTo(p colony.Point)~ توصّل الدرون لأي نقطة، وبعدين لف على ~colony.Waypoints()~ بالترتيب. نقاط التفتيش بتتغير كل مرة، ورقم النجوم بيتحسب على أقصر طريق.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch07NavigationEnglish,
			Title: "Navigation System",
			Goal:  "Harvest only the crystal in colony.Targets(). Other crystal is contaminated!",
			Lesson: `# Methods: Functions Attached to a Type ⚙️

A **method** is a function attached to a type. Instead of ~goTo(nav, p)~, you write ~nav.GoTo(p)~:

~~~go
type Nav struct {
	X, Y int
}

func (n *Nav) GoTo(p colony.Point) {
	//  ^^^^^^ This is the receiver
	for n.X < p.X {
		colony.Move(colony.East)
		n.X++
	}
}
~~~

## Why ~*Nav~ rather than ~Nav~?
- ~func (n Nav)~ (**value receiver**): the method receives a **copy**, and its changes are lost.
- ~func (n *Nav)~ (**pointer receiver**): the method works on **the original**.

Since GoTo changes ~n.X~ and ~n.Y~, it needs ~*Nav~. We will explore pointers in chapter 8.

~~~go
nav := Nav{}       // X: 0, Y: 0
nav.GoTo(colony.Point{X: 2, Y: 1})
fmt.Println(nav)   // {2 1}
~~~

## Methods that only read

~~~go
func (n Nav) Here() colony.Point {
	return colony.Point{X: n.X, Y: n.Y}
}
~~~

This changes nothing, so a value receiver is enough.

## Your mission
Define a ~Nav~ type with methods, and harvest **only** the crystal in ~colony.Targets()~. Leave any other crystal you pass untouched!`,
			Starter: `package main

import "orbital/colony"

// Nav keeps track of the drone's position
type Nav struct {
	X, Y int
}

// GoTo moves the drone to p and updates its position in n
func (n *Nav) GoTo(p colony.Point) {
	for n.X < p.X {
		colony.Move(colony.East)
		n.X++
	}
	// ...
}

func main() {
	nav := Nav{}
	for _, t := range colony.Targets() {
		nav.GoTo(t)
		// ...
	}
}
`,
			Hints: []string{
				"Keep navigation state with the movement behavior.",
				"Separate travel from handling each target.",
				"Inspect ~colony.Targets()~ and the code after ~nav.GoTo(t)~.",
			},
		},
		LessonRevision: 1, Steps: ch07NavigationArabic,
		ID: "7-2-nav-system", Chapter: 7, Title: "نظام الملاحة",
		Goal: "احصد كريستال الأهداف بس (colony.Targets()). الكريستال التاني ملوّث!",
		Par:  60, Gold: 45, MaxTicks: 300,
		Requires: []string{"struct", "method"},
		Setup: func(w *sim.World) {
			w.Init(9, 7, sim.Regolith)
			ps := randomPoints(w, 9)
			decoys := ps[5:]
			w.Targets = ps[:5]
			placeCrystals(w, ps)
			w.Gold = routeCost(w.Targets) + len(w.Targets)
			w.Par = w.Gold + 12
			w.FailCheck = func(w *sim.World) string {
				if remaining(w, decoys) < len(decoys) {
					return i18n.T("حصدت كريستال ملوّث (مش في Targets). المهمة فشلت.", "You harvested contaminated crystal (not in Targets). Mission failed.")
				}
				return ""
			}
			w.Check = func(w *sim.World) (bool, string) {
				left := remaining(w, w.Targets)
				return left == 0, fmt.Sprintf(i18n.T("أهداف فاضلة: %d / %d", "Targets remaining: %d / %d"), left, len(w.Targets))
			}
		},
		Starter: `package main

import "orbital/colony"

// Nav بيفتكر مكان الدرون بنفسه
type Nav struct {
	X, Y int
}

// GoTo بتحرّك الدرون لـ p وبتحدّث مكانه في n
func (n *Nav) GoTo(p colony.Point) {
	for n.X < p.X {
		colony.Move(colony.East)
		n.X++
	}
	// ...
}

func main() {
	nav := Nav{}
	for _, t := range colony.Targets() {
		nav.GoTo(t)
		// ...
	}
}
`,
		Hints: []string{
			"خَلّي حالة الملاحة مع سلوك الحركة.",
			"افصل الوصول عن التعامل مع كل هدف.",
			"راجع ~colony.Targets()~ والكود بعد ~nav.GoTo(t)~.",
		},
		Lesson: `# Methods: دوال بتاعة نوع معيّن ⚙️

الـ **method** دالة مربوطة بنوع. بدل ~goTo(nav, p)~ بتكتب ~nav.GoTo(p)~:

~~~go
type Nav struct {
	X, Y int
}

func (n *Nav) GoTo(p colony.Point) {
	//  ^^^^^^ ده الـ receiver
	for n.X < p.X {
		colony.Move(colony.East)
		n.X++
	}
}
~~~

## ليه ~*Nav~ مش ~Nav~؟
- ~func (n Nav)~ (**value receiver**): الدالة بتاخد **نسخة**، وأي تغيير جواها بيضيع.
- ~func (n *Nav)~ (**pointer receiver**): الدالة بتشتغل على **الأصل نفسه**.

بما إن GoTo بتغيّر ~n.X~ و~n.Y~، لازم ~*Nav~. هنفهم الـ pointers بالتفصيل في الفصل 8.

~~~go
nav := Nav{}       // X: 0, Y: 0
nav.GoTo(colony.Point{X: 2, Y: 1})
fmt.Println(nav)   // {2 1}
~~~

## methods للقراءة بس

~~~go
func (n Nav) Here() colony.Point {
	return colony.Point{X: n.X, Y: n.Y}
}
~~~

دي مش بتغيّر حاجة، فـ value receiver كفاية.

## المهمة
اعمل نوع ~Nav~ بـ methods، واحصد **بس** الكريستال اللي في ~colony.Targets()~. لو عدّيت على كريستال تاني متحصدهوش!`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch07ShortestRouteEnglish,
			Title: "The Shortest Route",
			Goal:  "Harvest every target. Choose the order; shorter routes earn more stars",
			Lesson: `# A Simple Algorithm: Visit the Nearest 🧭

~Targets()~ returns targets in a random order. Following it directly may mean lots of travel. The **greedy** idea: choose the **nearest** target at each step.

## A method that returns a value

~~~go
type Spot struct{ X, Y int }

func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}

a, b := Spot{0, 0}, Spot{3, 4}
fmt.Println(a.Dist(b)) // 7
~~~

This is **Manhattan distance**: the number of moves when you cannot travel diagonally.

## Finding the nearest target

~~~go
best := 0
for i, t := range left {
	if here.Dist(t) < here.Dist(left[best]) {
		best = i
	}
}
~~~

## Removing an element from a slice

~~~go
left = append(left[:best], left[best+1:]...)
~~~

~left[:best]~ contains the elements before it; ~left[best+1:]~ contains those after it. ~...~ expands the second slice into individual arguments.

## Think about it 🤔
Greedy is not always **optimal**, but it is often close and easy to write. Finding the true optimum is the "traveling salesperson problem," one of computer science's hardest problems.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

type Spot struct {
	X, Y int
}

// Dist is the number of moves from s to o
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func main() {
	here := Spot{0, 0}
	for _, t := range colony.Targets() {
		s := Spot{t.X, t.Y}
		fmt.Println(s, "distance", here.Dist(s))
	}
}
`,
			Hints: []string{
				"Compare target distances from your current position.",
				"Track remaining targets and refresh your current point.",
				"Inspect the ~Dist~ method and the distance print in ~main~.",
			},
		},
		LessonRevision: 1, Steps: ch07ShortestRouteArabic,
		ID: "7-3-shortest-route", Chapter: 7, Title: "أقصر طريق",
		Goal: "احصد كل الأهداف. الترتيب عليك، وأقصر طريق ياخد النجوم",
		Par:  70, Gold: 50, MaxTicks: 400,
		Requires: []string{"method"},
		Setup: func(w *sim.World) {
			w.Init(10, 8, sim.Regolith)
			w.Targets = randomPoints(w, 7)
			placeCrystals(w, w.Targets)
			given := routeCost(w.Targets) + 7
			w.Gold = min(greedyCost(w.Targets)+7, given)
			w.Par = max(given, w.Gold+6)
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

type Spot struct {
	X, Y int
}

// Dist عدد الحركات من s لـ o
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func main() {
	here := Spot{0, 0}
	for _, t := range colony.Targets() {
		s := Spot{t.X, t.Y}
		fmt.Println(s, "بعيدة", here.Dist(s))
	}
}
`,
		Hints: []string{
			"قارن مسافات الأهداف من مكانك الحالي.",
			"تابع الأهداف الباقية وحدّث موضعك الحالي.",
			"راجع method ~Dist~ وسطر طباعة المسافة في ~main~.",
		},
		Lesson: `# خوارزمية بسيطة: روح للأقرب 🧭

الترتيب اللي ~Targets()~ بترجّعه عشوائي. لو مشيت عليه زي ما هو، هتلف كتير. فكرة **greedy**: كل خطوة اختار **أقرب** هدف.

## method بترجّع قيمة

~~~go
type Spot struct{ X, Y int }

func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}

a, b := Spot{0, 0}, Spot{3, 4}
fmt.Println(a.Dist(b)) // 7
~~~

ده اسمه **Manhattan distance**: عدد الحركات لو مش بتتحرك بالقطر.

## إزاي تلاقي الأقرب

~~~go
best := 0
for i, t := range left {
	if here.Dist(t) < here.Dist(left[best]) {
		best = i
	}
}
~~~

## امسح عنصر من slice

~~~go
left = append(left[:best], left[best+1:]...)
~~~

~left[:best]~ هو اللي قبله، و~left[best+1:]~ اللي بعده. و~...~ بتفرد الـ slice التانية كأنها قيم منفصلة.

## فكّر 🤔
greedy مش دايماً **أحسن** حل، بس غالباً قريب منه وسهل يتكتب. أحسن حل فعلاً اسمه "مشكلة البائع المتجوّل"، ودي من أصعب المشاكل في علوم الكمبيوتر.`,
	})

	compositeUnitsStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.Invoice())
	// Build the report described in the lesson.
}
`
	add(&Level{ID: "7-4-composite-units", Chapter: 7, Title: "وحدات مركّبة", Goal: "استخدم embedding وmethod لتوزيع فاتورة الطاقم بالـ cents",
		LessonRevision: 1, Steps: ch07CompositeArabic,
		Requires: []string{"embed", "method"}, Bridge: []int{12},
		Par: 0, Gold: 0, MaxTicks: 100,
		Starter: compositeUnitsStarter, Hints: []string{
			"الـembedding بيرقّي fields النوع الداخلي.",
			"افصل بيانات الفاتورة عن method الأنصبة.",
			"استخدم قيم Invoice وراجع الـreceiver والـfields.",
		}, Lesson: `# وحدة جوّا وحدة 🧩

عرفنا structs في نقاط الطريق وmethods في نظام الملاحة. دلوقتي هنجمع نوعين مع بعض بـ embedding، ونستخدمهم لتقسيم فاتورة الطاقم من غير ما نخلق أو نضيّع cent.

## قبل ما تكتب
- الحقول تبقى متاحة مباشرة على الوحدة المركّبة إزاي؟
- مين يدفع cent الزيادة لما القسمة يبقى ليها باقي؟

## embedding مش حقل باسم تاني
~~~go
type Badge struct { Name string }
type Pilot struct { Badge }
p := Pilot{Badge: Badge{Name: "Nova"}}
fmt.Println(p.Name) // Nova
~~~
الحقل ~Badge~ جوه Pilot مكتوب بنوعه من غير اسم. ده embedding: حقول النوع الداخلي بتبقى متاحة على النوع الخارجي مباشرة. لسه Badge قيمة جوه Pilot؛ دي مش وراثة بين classes.

~~~go
p.Badge.Name = "Vega"
fmt.Println(p.Name) // Vega
~~~
الطريقتين بيوصلوا لنفس الحقل. الأسماء المتكررة لو فيه أكتر من نوع داخلي ممكن تحتاج تحديد المسار؛ إحنا عندنا نوع داخلي واحد بس.

## البيانات المالية
~~~go
total, people := colony.Invoice()
fmt.Println(total, people)
~~~
~Invoice~ بترجّع إجمالي نهائي بالـ cents وعدد الطاقم. الإجمالي مش سالب والعدد موجب؛ مفيش parsing أو tips أو مدخلات ناقصة. الاتنين بيتغيّروا مع التشغيل، وممكن الإجمالي يبقى صفر أو الطاقم فرد واحد.

## قسمة صحيحة وباقي: مثال محلول
~~~go
supplies, teams := 17, 4
base := supplies / teams
extra := supplies % teams
fmt.Println(base, extra) // 4 1
~~~
~%~ بتجيب باقي القسمة، زي مثال divmod في 5-1. لو بنوزّع 17 فلتر على 4 فرق، كل فريق ياخد 4، وأول فريق ياخد واحد زيادة. النتيجة ~[5 4 4 4]~ ومجموعها لسه17.

في الفاتورة هنطبق نفس السياسة على cents: أول remainder أفراد ياخدوا cent زيادة، والباقي ياخدوا base. مفيش تقريب عشري لكل فرد.

## method على الوحدة المركّبة
~~~go
type Pilot struct { Badge }
func (p Pilot) Label() string {
 return p.Name
}
~~~
دي method بقراءة فقط، زي 7-2؛ value receiver كفاية. بتستخدم الحقل المرقّى ~p.Name~ من غير كتابة ~p.Badge.Name~ كل مرة. في مهمتك الـ method هترجع slice بالأنصبة، بأدوات make وrange اللي عرفناها في الفصل6.

## أخطاء شائعة
~~~go
type Pilot struct { badge Badge }
p := Pilot{badge: Badge{Name: "Nova"}}
fmt.Println(p.Name) // الحقل ده مش مترقّي
~~~
~badge Badge~ حقل مسمّى عادي، مش embedding. ولو قسّمت17 على4 وكرّرت الناتج4 بس، هتضيّع واحد؛ لازم توزّع الباقي كمان.

## المهمة
عرّف ~Ledger struct { TotalCents, People int }~، وضمّه بدون اسم حقل جوه ~CrewUnit struct { Ledger }~. اكتب ~Shares() []int~ كـ method على CrewUnit، واستخدم ~u.TotalCents~ و~u.People~ مباشرة.

ابني الوحدة من Invoice، وابعت slice الأنصبة بتقرير واحد. ترتيبها هو ترتيب أفراد الطاقم: أول ~total%people~ أفراد بياخدوا cent زيادة، والباقي النصيب الأساسي. المجموع لازم يساوي الإجمالي بالظبط، حتى مع صفر أو فرد واحد.`,
		En: &LevelText{LessonRevision: 1, Steps: ch07CompositeEnglish, Title: "Composite Units", Goal: "Use embedding and a method to allocate crew invoice cents", Starter: compositeUnitsStarter, Hints: []string{
			"Embedding promotes inner type fields.",
			"Separate invoice data from the shares method.",
			"Use Invoice values and inspect the receiver and fields.",
		}, Lesson: `# One Unit Inside Another 🧩

Waypoints introduced structs and Navigation System introduced methods. Now combine types with embedding and split a crew invoice without inventing or losing a cent.

## Before you write
- How can the outer unit expose its embedded fields directly?
- Who receives the extra cent when division leaves a remainder?

## Embedding is not a renamed field
~~~go
type Badge struct { Name string }
type Pilot struct { Badge }
p := Pilot{Badge: Badge{Name: "Nova"}}
fmt.Println(p.Name) // Nova
~~~
Inside Pilot, ~Badge~ has a type but no field name. This is embedding: inner fields become available directly on the outer type. Badge is still a value inside Pilot; this is not class inheritance.

~~~go
p.Badge.Name = "Vega"
fmt.Println(p.Name) // Vega
~~~
Both paths reach the same field. Conflicting names from several embedded types can require explicit paths; this mission embeds only one type.

## Financial inputs
~~~go
total, people := colony.Invoice()
fmt.Println(total, people)
~~~
~Invoice~ returns a final total in cents and a crew count. The total is nonnegative and the count positive; there is no parsing, tip calculation, or missing input. Both vary across runs. Zero total and a single crew member are valid.

## Integer division and remainder: a worked example
~~~go
supplies, teams := 17, 4
base := supplies / teams
extra := supplies % teams
fmt.Println(base, extra) // 4 1
~~~
~%~ returns the remainder, as in divmod in 5-1. Splitting 17 filters across four teams gives everyone 4 and the first team one extra. The result ~[5 4 4 4]~ still totals 17.

Use the same rule for cents: the first remainder crew members get one extra cent; everyone else gets base. Do not round a decimal amount independently for each person.

## A method on the composite unit
~~~go
type Pilot struct { Badge }
func (p Pilot) Label() string {
 return p.Name
}
~~~
This read-only method uses a value receiver, as in 7-2. It accesses promoted ~p.Name~ without repeating ~p.Badge.Name~. Your mission method returns a shares slice using make and range from chapter 6.

## Common mistakes
~~~go
type Pilot struct { badge Badge }
p := Pilot{badge: Badge{Name: "Nova"}}
fmt.Println(p.Name) // This field was not promoted
~~~
~badge Badge~ is an ordinary named field, not embedding. Also, repeating only the quotient 4 when splitting 17 by 4 loses one unit; distribute the remainder too.

## Your mission
Define ~Ledger struct { TotalCents, People int }~ and embed it without a field name in ~CrewUnit struct { Ledger }~. Write a ~Shares() []int~ method on CrewUnit using ~u.TotalCents~ and ~u.People~ directly.

Build the unit from Invoice and report the shares slice once. Its order is crew order: the first ~total%people~ members receive an extra cent; the rest receive the base share. The sum must equal the total exactly, including zero total and one-person cases.`},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			w.Data.People = 1 + w.Rand.Intn(7)
			w.Data.TotalCents = w.Rand.Intn(5000)
			if w.Rand.Intn(4) == 0 {
				w.Data.TotalCents = (w.Data.TotalCents / w.Data.People) * w.Data.People
			}
			if w.Rand.Intn(8) == 0 {
				w.Data.TotalCents = 0
			}
			shares := make([]int, w.Data.People)
			for i := range shares {
				shares[i] = w.Data.TotalCents / w.Data.People
				if i < w.Data.TotalCents%w.Data.People {
					shares[i]++
				}
			}
			want := []string{fmt.Sprint(shares)}
			w.Check = func(w *sim.World) (bool, string) {
				if slices.Equal(w.Reports, want) {
					return true, i18n.T("كل التقارير مطابقة وبالترتيب", "All reports match in order")
				}
				return false, i18n.T("التقارير لسه ناقصة أو مش مطابقة؛ راجع القيم وترتيبها", "Reports are missing or do not match; check values and order")
			}
		},
	})

	crewExamStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.Exam())
	// Build the report described in the lesson.
}
`
	add(&Level{ID: "7-5-crew-exam", Chapter: 7, Title: "اختبار الطاقم", Goal: "جاوب اختبار أرض النباتات واحسب نسبة الأسئلة الصالحة",
		LessonRevision: 1, Steps: ch07CrewExamArabic,
		Requires: []string{"struct", "range"}, Bridge: []int{5},
		Par: 0, Gold: 0, MaxTicks: 100,
		Starter: crewExamStarter, Hints: []string{
			"كل سؤال يجمع بياناته في struct واحدة.",
			"راجع السؤال واختياراته، وبعدها النتيجة.",
			"خلي معاك الـindex اللي range رجّعه لما تنادي ~colony.Answer~.",
		}, Lesson: `# سؤال كامل بدل متغيرات متفرّقة 🧑‍🚀

الطاقم محتاج يراجع أرض كل محصول. بنك الأسئلة بيتغيّر وترتيب الاختيارات بيتبدّل، لكن كل سؤال بيجمع بياناته في struct واحدة زي ما نقاط الطريق جمعت X وY.

## قبل ما تكتب
- هتحافظ على رقم السؤال الأصلي لو تخطّيت سؤال بايظ إزاي؟
- السؤال اللي ملوش اختيارات يدخل في مقام النسبة؟

## اقرأ الأسئلة كـ structs
~~~go
questions := colony.Exam()
for i, question := range questions {
 fmt.Println(i, question.Prompt, question.Choices)
}
~~~
~Exam~ بترجّع ~[]colony.Question~. كل Question فيها ~Prompt string~ اسم المحصول، و~Choices []string~ أسماء أراضٍ. ~range~ بترجّع index السؤال ونسخة من قيمته، زي slice النقاط في 7-1.

~~~go
if len(question.Choices) == 0 {
 continue
}
~~~
~continue~ بتنقل للّفة اللي بعدها، زي 5-2. افحص الاختيارات قبل الفهرسة؛ فيه سؤال offline من غير اختيارات، وسيبه من غير إجابة ولا زيادة في عدد الأسئلة الصالحة.

## اختار حسب بيانات السؤال
~~~go
switch question.Prompt {
case "moss":
 fmt.Println("soil")
case "lichen":
 fmt.Println("ice")
}
~~~
دي قواعد عرفنا أصلها في تضاريس مختلطة. في الاختبار ده الإجابات بالظبط: moss أرضها soil، وlichen أرضها ice، وcrystal مكانها regolith. مفيش parsing؛ قارن النصوص مباشرة، وبعدين دور في Choices على النص المطلوب.

## سجّل الاختيار
~~~go
accepted := colony.Answer(0, 1)
fmt.Println(accepted)
~~~
~Answer(i, choice)~ بتسجّل رقم الاختيار من صفر للسؤال i، وبترجّع bool لو الرقمين صالحين. true مش معناها الإجابة صح معرفياً؛ إنت بتحدد الصح بمقارنة الاختيار بقاعدة المحصول. استخدم index السؤال الأصلي من range، وindex الاختيار اللي لقيته، مش أرقام ثابتة زي المثال.

## struct للنتيجة: مثال محلول
~~~go
type Progress struct { Done, Total int }
p := Progress{Done: 2, Total: 3}
percent := p.Done * 100 / p.Total
fmt.Println(percent) // 66
~~~
المثال بيحسب تقدّم تدريب، مش درجات اختبار المهمة. جمع العدادات في struct بيخلّي علاقتهم واضحة. القسمة الصحيحة بتشيل الكسر، واضرب في100 قبل القسمة عشان متضيّعش النسبة.

## أخطاء شائعة
~~~go
for i := range questions {
 colony.Answer(i, 0) // أول اختيار مش دايماً الصح
}
~~~
ترتيب الاختيارات بيتبدّل، والسؤال الفاضي أصلاً مفيهوش اختيار0. ولو حسبت كل عناصر بنك الأسئلة في المقام، هتخصم درجة سؤال ملوش إجابة. لو عدد الأسئلة الصالحة صفر، النسبة تبقى صفر من غير قسمة.

## المهمة
عرّف struct عندك للعدادين ~Correct~ و~Attempted~. جاوب كل سؤال صالح بالقاعدة: moss→soil، lichen→ice، crystal→regolith. سيب السؤال اللي من غير اختيارات من غير Answer.

احسب ~correct*100/attempted~ كعدد صحيح، وابعت النسبة في تقرير واحد. المطلوب إجابات صحيحة لكل الأسئلة الصالحة، فالنتيجة النهائية100؛ التحقق بيراجع الإجابات المسجّلة نفسها كمان، مش التقرير بس.`,
		En: &LevelText{LessonRevision: 1, Steps: ch07CrewExamEnglish, Title: "Crew Exam", Goal: "Answer the crop-ground exam and score usable questions", Starter: crewExamStarter, Hints: []string{
			"Each question groups its data in one struct.",
			"Inspect the question, its choices, then the result.",
			"Keep range's index when calling ~colony.Answer~.",
		}, Lesson: `# A Whole Question Instead of Scattered Variables 🧑‍🚀

The crew is reviewing crop terrain. The bank and choice order vary, but every question groups its data in one struct, just as waypoints group X and Y.

## Before you write
- How will you retain the original question index when skipping a broken question?
- Should a question without choices count in the percentage denominator?

## Read questions as structs
~~~go
questions := colony.Exam()
for i, question := range questions {
 fmt.Println(i, question.Prompt, question.Choices)
}
~~~
~Exam~ returns ~[]colony.Question~. Each Question has ~Prompt string~, the crop name, and ~Choices []string~, terrain names. ~range~ gives the question index and a copy of its value, like the points slice in 7-1.

~~~go
if len(question.Choices) == 0 {
 continue
}
~~~
~continue~ advances to the next iteration, as in 5-2. Check choices before indexing. One offline question has no choices; leave it unanswered and do not increase the usable-question count.

## Choose using the question data
~~~go
switch question.Prompt {
case "moss":
 fmt.Println("soil")
case "lichen":
 fmt.Println("ice")
}
~~~
These rules build on Mixed Terrain. For this exam, the exact answers are moss→soil, lichen→ice, crystal→regolith. No parsing is involved: compare strings directly, then find the required text in Choices.

## Record a choice
~~~go
accepted := colony.Answer(0, 1)
fmt.Println(accepted)
~~~
~Answer(i, choice)~ records a zero-based choice for question i and returns a bool indicating valid indices. True does not mean the answer is factually correct; you determine correctness by matching the crop rule. Use the original question index from range and the choice index you found, rather than the sample's fixed numbers.

## A score struct: a worked example
~~~go
type Progress struct { Done, Total int }
p := Progress{Done: 2, Total: 3}
percent := p.Done * 100 / p.Total
fmt.Println(percent) // 66
~~~
This computes training progress rather than the mission score. Grouping counters in a struct makes their relationship clear. Integer division discards the fraction; multiply by 100 before division to retain the percentage.

## Common mistakes
~~~go
for i := range questions {
 colony.Answer(i, 0) // The first choice is not always correct
}
~~~
Choice order varies, and the empty question has no choice 0 at all. Counting all bank entries in the denominator also penalizes an unanswerable question. If there are no usable questions, report zero without dividing.

## Your mission
Define your own struct with ~Correct~ and ~Attempted~ counters. Answer every usable question using moss→soil, lichen→ice, crystal→regolith. Do not call Answer for the question without choices.

Compute integer ~correct*100/attempted~ and send one percentage report. All usable answers must be correct, so the final score is 100. Verification also inspects recorded answers, not just your report.`},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			n := 4 + w.Rand.Intn(4)
			crops := []string{"moss", "lichen", "crystal"}
			correct := map[int]int{}
			empty := w.Rand.Intn(n + 1)
			for i := 0; i <= n; i++ {
				if i == empty {
					w.Data.Questions = append(w.Data.Questions, sim.Question{Prompt: "offline"})
					continue
				}
				crop := crops[w.Rand.Intn(len(crops))]
				choices := []string{"soil", "ice", "regolith"}
				w.Rand.Shuffle(len(choices), func(i, j int) { choices[i], choices[j] = choices[j], choices[i] })
				ground := "soil"
				switch crop {
				case "lichen":
					ground = "ice"
				case "crystal":
					ground = "regolith"
				}
				for j, choice := range choices {
					if choice == ground {
						correct[i] = j
					}
				}
				w.Data.Questions = append(w.Data.Questions, sim.Question{Prompt: crop, Choices: choices})
			}
			w.Check = func(w *sim.World) (bool, string) {
				done := len(w.Data.Answers) == len(correct)
				for i, choice := range correct {
					got, ok := w.Data.Answers[i]
					done = done && ok && got == choice
				}
				if done && len(w.Reports) == 1 && w.Reports[0] == "100" {
					return true, i18n.T("كل الإجابات صح والسؤال الفاضي خارج النسبة", "All answers are correct and the empty question is excluded")
				}
				return false, i18n.T("مستني إجابة كل سؤال صالح والنسبة؛ سيب السؤال اللي من غير اختيارات", "Waiting for every usable answer and the score; skip the question without choices")
			}
		},
	})

}
