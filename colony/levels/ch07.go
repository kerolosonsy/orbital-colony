package levels

import (
	"fmt"

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
				"You need 4 loops in goTo: x too small (East), x too large (West), y too small (North), y too large (South).",
				"In main: ~for _, p := range points { goTo(p) }~.",
			},
		},
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
			"محتاج 4 حلقات في goTo: x أصغر (East)، x أكبر (West)، y أصغر (North)، y أكبر (South).",
			"في main: ~for _, p := range points { goTo(p) }~.",
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
				"Complete GoTo for the other directions, updating ~n.X~ or ~n.Y~ after each move.",
				"After ~nav.GoTo(t)~, call ~colony.Harvest()~.",
				"You can add another method: ~func (n *Nav) Collect(p colony.Point) { n.GoTo(p); colony.Harvest() }~.",
			},
		},
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
			"كمّل GoTo بباقي الاتجاهات، وكل مرة تتحرك حدّث ~n.X~ أو ~n.Y~.",
			"بعد ~nav.GoTo(t)~ اعمل ~colony.Harvest()~.",
			"ممكن تضيف method تانية: ~func (n *Nav) Collect(p colony.Point) { n.GoTo(p); colony.Harvest() }~.",
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
				"The greedy approach: from your current position, visit the nearest unvisited target. Repeat.",
				"Keep targets in a slice and remove each one you visit: ~left = append(left[:i], left[i+1:]...)~.",
				"You need goTo from the previous mission; the same Nav type can help here too.",
			},
		},
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
			"الطريقة الطماعة (greedy): من مكانك، روح لأقرب هدف لسه ماروحتلوش. وكرّر.",
			"احفظ الأهداف في slice، ولما تروح لواحد امسحه: ~left = append(left[:i], left[i+1:]...)~.",
			"محتاج goTo من المهمة اللي فاتت، ونفس Nav ممكن ينفع هنا.",
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
}
