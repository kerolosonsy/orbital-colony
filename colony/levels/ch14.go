package levels

import (
	"fmt"
	"slices"
	"strconv"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

// halves returns n random values in [lo, hi) with .5 steps, so float sums are exact.
func halves(w *sim.World, n int, lo, hi float64) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = lo + float64(w.Rand.Intn(int((hi-lo)*2)))/2
	}
	return out
}

func reportsAll(w *sim.World, wants ...string) (int, bool) {
	got := 0
	for _, want := range wants {
		if reported(w, want) {
			got++
		}
	}
	return got, got == len(wants)
}

var dirs = [4][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} // North, East, South, West

// carveTunnel turns a rock map into a tree of tunnels from (0, 0) and
// returns the tunnel's far end and its length.
func carveTunnel(w *sim.World, mainLen int) (sim.Point, int) {
	open := map[sim.Point]bool{{X: 0, Y: 0}: true}
	w.At(0, 0).Ground = sim.Regolith
	canOpen := func(q, from sim.Point) bool {
		if q.X < 0 || q.Y < 0 || q.X >= w.W || q.Y >= w.H || open[q] {
			return false
		}
		for _, d := range dirs {
			n := sim.Point{X: q.X + d[0], Y: q.Y + d[1]}
			if n != from && open[n] {
				return false // would create a loop
			}
		}
		return true
	}
	walk := func(from sim.Point, steps int) (sim.Point, []sim.Point) {
		at, path := from, []sim.Point{}
		for range steps {
			var opts []sim.Point
			for _, d := range dirs {
				if q := (sim.Point{X: at.X + d[0], Y: at.Y + d[1]}); canOpen(q, at) {
					opts = append(opts, q)
				}
			}
			if len(opts) == 0 {
				break
			}
			at = opts[w.Rand.Intn(len(opts))]
			open[at] = true
			w.At(at.X, at.Y).Ground = sim.Regolith
			path = append(path, at)
		}
		return at, path
	}
	end, path := walk(sim.Point{}, mainLen)
	for b := 0; b < 3 && len(path) > 2; b++ {
		walk(path[w.Rand.Intn(len(path)-1)], 2+w.Rand.Intn(3))
	}
	return end, len(path)
}

// dfsCost simulates a depth-first explorer trying North, East, South, West.
func dfsCost(w *sim.World, goal sim.Point) int {
	cost, at := 1, sim.Point{} // scan the start tile
	visited := map[sim.Point]bool{at: true}
	var stack [][2]int
	for at != goal {
		moved := false
		for _, d := range dirs {
			n := sim.Point{X: at.X + d[0], Y: at.Y + d[1]}
			if n.X < 0 || n.Y < 0 || n.X >= w.W || n.Y >= w.H || visited[n] {
				continue
			}
			visited[n] = true
			cost++ // the move attempt
			if w.At(n.X, n.Y).Ground != sim.Rock {
				at, moved = n, true
				stack = append(stack, d)
				cost++ // scan the new tile
				break
			}
		}
		if !moved {
			d := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			at = sim.Point{X: at.X - d[0], Y: at.Y - d[1]}
			cost++
		}
	}
	return cost + 1 + len(stack) // harvest, then walk the stack home
}

func init() {
	add(&Level{
		En: &LevelText{
			Title: "One Function for Every Type",
			Goal:  "Report the greatest depth (int) and highest temperature (float64) using one generic Max function",
			Lesson: `# Generics: Write It Once for Multiple Types 🧬

~MaxInt~ and ~MaxFloat~ have identical code; only the **type** differs. Before Go 1.18, you had to duplicate it or use ~any~ and lose type safety. Now there are **type parameters**:

~~~go
import "cmp"

func Max[T cmp.Ordered](xs []T) T {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

Max([]int{3, 9, 2})        // T = int
Max([]float64{1.5, -2})    // T = float64
Max([]string{"b", "a"})    // T = string ✓
~~~

## Read it this way
- ~[T cmp.Ordered]~: "there is a type named T, and it must support comparisons with ~<~ and ~>~."
- ~cmp.Ordered~ is a **constraint**: it restricts the types T can represent.
- Go infers T from the argument (type inference).

## any as a constraint
~~~go
func First[T any](xs []T) T { return xs[0] }
~~~
~any~ allows any type, but then you cannot use ~>~, because some types do not support it.

## Go's own library is generic
~slices.Max~, ~slices.Sort~, and ~maps.Keys~ use generics in exactly this way.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func MaxInt(xs []int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

// MaxFloat is identical... except for float64 😩

func main() {
	depths := colony.Depths()         // []int
	temps := colony.Temperatures()    // []float64
	fmt.Println(depths, temps)
	colony.Report(MaxInt(depths))
}
`,
			Hints: []string{
				"Change MaxInt to ~func Max[T cmp.Ordered](xs []T) T~ and ~import \"cmp\"~.",
				"Keep the same body, then use ~colony.Report(Max(depths))~ and ~colony.Report(Max(temps))~.",
			},
		},
		ID: "14-1-generic-max", Chapter: 14, Title: "دالة لكل الأنواع",
		Goal: "ابعت أكبر عمق (int) وأعلى حرارة (float64) بدالة Max واحدة generic",
		Par:  0, Gold: 0, MaxTicks: 20,
		Requires: []string{"generic"},
		Setup: func(w *sim.World) {
			w.Init(5, 3, sim.Regolith)
			for range 8 {
				w.Data.Depths = append(w.Data.Depths, 10+w.Rand.Intn(490))
			}
			w.Data.Temps = halves(w, 8, -90, 40)
			wants := []string{strconv.Itoa(slices.Max(w.Data.Depths)), fmt.Sprint(slices.Max(w.Data.Temps))}
			w.Check = func(w *sim.World) (bool, string) {
				got, ok := reportsAll(w, wants...)
				return ok, fmt.Sprintf(i18n.T("قيم صح وصلت: %d / 2", "Correct values received: %d / 2"), got)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func MaxInt(xs []int) int {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

// MaxFloat نسخة طبق الأصل… بس للـ float64 😩

func main() {
	depths := colony.Depths()         // []int
	temps := colony.Temperatures()    // []float64
	fmt.Println(depths, temps)
	colony.Report(MaxInt(depths))
}
`,
		Hints: []string{
			"غيّر MaxInt لـ ~func Max[T cmp.Ordered](xs []T) T~، واعمل ~import \"cmp\"~.",
			"جوّاها نفس الكود بالظبط، وبعدين ~colony.Report(Max(depths))~ و~colony.Report(Max(temps))~.",
		},
		Lesson: `# Generics: اكتبها مرة لكل الأنواع 🧬

~MaxInt~ و~MaxFloat~ نفس الكود بالحرف، والفرق الوحيد هو **النوع**. قبل Go 1.18 كنت لازم تكرّر الكود أو تستخدم ~any~ وتخسر الـ type safety. دلوقتي فيه **type parameters**:

~~~go
import "cmp"

func Max[T cmp.Ordered](xs []T) T {
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

Max([]int{3, 9, 2})        // T = int
Max([]float64{1.5, -2})    // T = float64
Max([]string{"b", "a"})    // T = string ✓
~~~

## اقراها كده
- ~[T cmp.Ordered]~: "فيه نوع اسمه T، ولازم يكون من الأنواع اللي ينفع تتقارن بـ ~<~ و~>~".
- ~cmp.Ordered~ اسمه **constraint**: بيحدد T ممكن يكون إيه.
- Go بيستنتج T لوحده من الـ argument (type inference).

## any كـ constraint
~~~go
func First[T any](xs []T) T { return xs[0] }
~~~
~any~ معناها أي نوع، بس ساعتها مينفعش تستخدم ~>~ لأن مش كل الأنواع بتتقارن.

## مكتبة Go نفسها generic
~slices.Max~ و~slices.Sort~ و~maps.Keys~ اللي استخدمتهم قبل كده مكتوبين بـ generics بالظبط زي كده.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "The Tunnel",
			Goal:  "The map is hidden! Explore the tunnel, harvest the crystal at its end, and return to (0, 0) using your Stack[T]",
			Lesson: `# Generic types: Stack[T] 📚

Functions are not the only things that can be generic. **Types** can too:

~~~go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T // The zero value of any type
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

var moves Stack[colony.Direction] // A stack of directions
var nums Stack[int]               // A stack of numbers
~~~

## var zero T
You do not know what T will be, so you cannot write 0 or "". ~var zero T~ gives the right zero value for any type.

## Why a stack here? (Depth-First Search)
The tunnel branches, and the map is hidden. The approach:
1. Try a new direction. If successful, **Push** the direction.
2. If every direction is blocked (a dead end), **Pop** and backtrack one step.
3. When you find the crystal, the stack holds **the exact return route**: Pop until empty, moving in the opposite direction each time.

This is **DFS (depth-first search)**, one of computing's fundamental algorithms.

## The opposite direction
~North=0, East=1, South=2, West=3~, so ~(d + 2) % 4~ gives the opposite.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// Stack: last in, first out
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop removes and returns the last element. ok = false when empty
func (s *Stack[T]) Pop() (v T, ok bool) {
	// ...
	return v, false
}

func opposite(d colony.Direction) colony.Direction {
	return (d + 2) % 4
}

func main() {
	var path Stack[colony.Direction]
	path.Push(colony.North)
	fmt.Println(path.items, opposite(colony.North))
	// Explore: try each direction. If Move succeeds, Push. At a dead end, Pop and backtrack
}
`,
			Hints: []string{
				"Pop: if ~len(s.items) == 0~, return ok = false. Otherwise take ~s.items[len(s.items)-1]~ and shorten the slice.",
				"Remember visited tiles in ~map[colony.Point]bool~ to avoid circles. On each new tile, check ~colony.GetEntity() == colony.Crystal~.",
				"To return: while ~Pop~ returns ok, move ~opposite(d)~.",
			},
		},
		ID: "14-2-stack", Chapter: 14, Title: "النفق",
		Goal: "الخريطة مخفية! استكشف النفق، احصد الكريستال اللي في آخره، وارجع (0, 0) بـ Stack[T] بتاعك",
		Par:  80, Gold: 50, MaxTicks: 600,
		Requires: []string{"generictype"},
		Setup: func(w *sim.World) {
			w.Init(11, 7, sim.Rock)
			w.RockBlocks, w.HideMap = true, true
			end, _ := carveTunnel(w, 16)
			t := w.At(end.X, end.Y)
			t.Entity, t.Size = sim.Crystal, 5
			explore := dfsCost(w, end)
			w.Gold = explore*11/10 + 6 // explorers differ: some rescan while backtracking
			w.Par = explore*3/2 + 10
			w.Check = func(w *sim.World) (bool, string) {
				got := w.At(end.X, end.Y).Entity != sim.Crystal
				switch {
				case !got:
					return false, i18n.T("الكريستال لسه في آخر النفق", "The crystal is still at the end of the tunnel")
				case w.X != 0 || w.Y != 0:
					return false, fmt.Sprintf(i18n.T("معاك الكريستال! ارجع (0, 0). إنت عند (%d, %d)", "You have the crystal! Return to (0, 0). You are at (%d, %d)"), w.X, w.Y)
				}
				return true, i18n.T("رجعت بالكريستال", "Returned with the crystal")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// Stack: آخر حاجة تدخل، أول حاجة تطلع
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) {
	s.items = append(s.items, v)
}

// Pop بترجّع آخر عنصر وبتشيله. ok = false لو فاضية
func (s *Stack[T]) Pop() (v T, ok bool) {
	// ...
	return v, false
}

func opposite(d colony.Direction) colony.Direction {
	return (d + 2) % 4
}

func main() {
	var path Stack[colony.Direction]
	path.Push(colony.North)
	fmt.Println(path.items, opposite(colony.North))
	// استكشف: جرّب كل اتجاه. لو Move نجحت، اعمل Push. لو كله مقفول، Pop وارجع
}
`,
		Hints: []string{
			"Pop: لو ~len(s.items) == 0~ ارجّع ok = false. غير كده خد ~s.items[len(s.items)-1]~ وقصّر الـ slice.",
			"افتكر الخانات اللي زرتها في ~map[colony.Point]bool~ عشان متلفّش في دواير. وكل خانة جديدة: ~colony.GetEntity() == colony.Crystal~؟",
			"الرجوع: طول ما ~Pop~ بترجّع ok، اتحرك ~opposite(d)~.",
		},
		Lesson: `# Generic types: Stack[T] 📚

مش الدوال بس اللي ينفع تبقى generic. الـ **types** كمان:

~~~go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	var zero T // القيمة الصفرية لأي نوع
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

var moves Stack[colony.Direction] // Stack من الاتجاهات
var nums Stack[int]               // Stack من الأرقام
~~~

## var zero T
مش عارف T هيبقى إيه، فمتقدرش تكتب 0 أو "". ~var zero T~ بيدّيك القيمة الصفرية الصح لأي نوع.

## ليه Stack هنا؟ (Depth-First Search)
النفق متفرّع والخريطة مخفية. الطريقة:
1. جرّب اتجاه جديد. لو نجح، **Push** الاتجاه.
2. لو كل الاتجاهات مقفولة (طريق مسدود)، **Pop** وارجع خطوة.
3. لما تلاقي الكريستال، الـ stack فيه **طريق الرجوع بالظبط**: Pop لحد ما يفضى، وفي كل مرة امشي عكس الاتجاه.

ده اسمه **DFS (depth-first search)**، وهو من أهم الخوارزميات في الكمبيوتر.

## عكس الاتجاه
~North=0، East=1، South=2، West=3~، يعني ~(d + 2) % 4~ هو العكس.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Type Constraints",
			Goal:  "Report the sums of depths, weights (Kg), and temperatures with one Sum function",
			Lesson: `# Write Your Own Constraints 🔧

A constraint is an **interface** describing allowed types. It can contain a **union of types**:

~~~go
type Number interface {
	int | float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x // Allowed because every type in Number supports addition
	}
	return total
}
~~~

## The problem with defined types

~~~go
type Kg float64 // A new type with underlying type float64

Sum([]Kg{1.5, 2}) // ❌ Kg does not satisfy Number
~~~

~Kg~ is **distinct** from float64, even though it is based on it. This helps prevent accidentally adding kilograms to temperatures.

## The solution: the tilde

~~~go
type Number interface {
	~int | ~float64 // Any type with underlying type int or float64
}
~~~

A tilde before a type includes that type **and every type with that underlying type**.

## Ready-made in the library
~~~go
import "cmp"
cmp.Ordered // All ordered types, defined using tilde terms in this same way
~~~

## When should you use generics?
- When you write **the same code** for multiple types: containers and algorithms over slices or maps.
- **Not** for everything. Keep functions ordinary when they work with just one type. Unnecessary generics make code harder to read.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

type Number interface {
	int | float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

func main() {
	fmt.Println(colony.Depths(), colony.Weights(), colony.Temperatures())
	colony.Report(Sum(colony.Depths()))
	colony.Report(Sum(colony.Temperatures()))
	// colony.Report(Sum(colony.Weights())) // ← Uncomment and see what the compiler says
}
`,
			Hints: []string{
				"~colony.Kg~ is a new type based on float64 (~type Kg float64~). It is distinct from float64, so ~int | float64~ does not accept it.",
				"Put a tilde (the wave symbol near 1 on many keyboards) before int and float64 in Number. It means any type with one of those underlying types.",
			},
		},
		ID: "14-3-constraints", Chapter: 14, Title: "قيود الأنواع",
		Goal: "ابعت مجموع الأعماق والأوزان (Kg) ودرجات الحرارة، بدالة Sum واحدة",
		Par:  0, Gold: 0, MaxTicks: 20,
		Requires: []string{"generic"},
		Setup: func(w *sim.World) {
			w.Init(5, 3, sim.Regolith)
			sumD := 0
			for range 6 {
				d := 10 + w.Rand.Intn(200)
				w.Data.Depths = append(w.Data.Depths, d)
				sumD += d
			}
			var sumW sim.Kg
			for _, f := range halves(w, 6, 1, 90) {
				w.Data.Weights = append(w.Data.Weights, sim.Kg(f))
				sumW += sim.Kg(f)
			}
			w.Data.Temps = halves(w, 6, -60, 30)
			sumT := 0.0
			for _, t := range w.Data.Temps {
				sumT += t
			}
			wants := []string{strconv.Itoa(sumD), fmt.Sprint(sumW), fmt.Sprint(sumT)}
			w.Check = func(w *sim.World) (bool, string) {
				got, ok := reportsAll(w, wants...)
				return ok, fmt.Sprintf(i18n.T("مجاميع صح وصلت: %d / 3", "Correct sums received: %d / 3"), got)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

type Number interface {
	int | float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x
	}
	return total
}

func main() {
	fmt.Println(colony.Depths(), colony.Weights(), colony.Temperatures())
	colony.Report(Sum(colony.Depths()))
	colony.Report(Sum(colony.Temperatures()))
	// colony.Report(Sum(colony.Weights())) // ← شيل الـ comment وشوف الـ compiler هيقول إيه
}
`,
		Hints: []string{
			"~colony.Kg~ نوع جديد معمول من float64 (~type Kg float64~). هو مش float64 نفسه، فـ ~int | float64~ مش بتقبله.",
			"اكتب علامة التلدة (الموجة اللي جنب الرقم 1 في الكيبورد) قبل int وقبل float64 جوه Number. معناها: أي نوع أصله الأنواع دي.",
		},
		Lesson: `# Constraints بإيدك 🔧

الـ constraint هو **interface** بيوصف الأنواع المسموحة. وممكن يكون عبارة عن **قايمة أنواع**:

~~~go
type Number interface {
	int | float64
}

func Sum[T Number](xs []T) T {
	var total T
	for _, x := range xs {
		total += x // مسموح، لأن كل أنواع Number بتتجمع
	}
	return total
}
~~~

## مشكلة الأنواع المشتقة

~~~go
type Kg float64 // نوع جديد أساسه float64

Sum([]Kg{1.5, 2}) // ❌ Kg does not satisfy Number
~~~

~Kg~ نوع **مختلف** عن float64، حتى لو هو مبني عليه. ده بيحميك من إنك تجمع كيلوجرامات على درجات حرارة بالغلط.

## الحل: علامة التلدة

~~~go
type Number interface {
	~int | ~float64 // أي نوع "أساسه" int أو float64
}
~~~

علامة ~ قبل النوع معناها: النوع ده نفسه، **وأي نوع أساسه (underlying type) هو**.

## جاهزة في المكتبة
~~~go
import "cmp"
cmp.Ordered // كل الأنواع اللي بتتقارن، ومعمولة بالـ ~ دي بالظبط
~~~

## إمتى تستخدم generics؟
- لما تلاقي نفسك بتكتب **نفس الكود** لأكتر من نوع: containers، وخوارزميات على slices وmaps.
- **مش** لكل حاجة. لو دالة بتشتغل مع نوع واحد، خلّيها عادية. الـ generics بتصعّب القراية لو اتستخدمت من غير داعي.`,
	})
}
