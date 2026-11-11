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
	add(&Level{ID: "14-4-airlock-queue", Chapter: 14, Kind: "lab", Title: "طابور البوابة", Goal: "ابني طابور generic يحافظ على ترتيب الوصول", Requires: []string{"generictype"}, Bridge: []int{27},
		Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: airlockQueueStarter, ReadOnly: false}},
		Lesson: `# اللي وصل الأول يخرج الأول

الـ stack اللي فات بيمشي LIFO. بوابة الضغط محتاجة العكس: FIFO، أول راكب يدخل هو أول راكب يخرج، مهما نوع بياناته.

## قبل ما تكتب
- الـ zero value تقدر تمثل طابور فاضي من غير constructor؟
- هتميّز راكب قيمته صفر عن مفيش راكب إزاي؟

## generic type بنفس السلوك
~~~go
type Tray[T any] struct { Slots []T }
func (t *Tray[T]) Add(v T) { t.Slots = append(t.Slots,v) }
var tray Tray[string]
tray.Add("lens")
~~~
~T any~ بيسمح بأي نوع. ~Tray[string]~ نسخة بنوع محدد، والـ receiver لازم يذكر ~[T]~. الـ slice اللي قيمتها nil تقبل append، فمش محتاج تهيئة.

## نتيجة ومعاها وجودها
~~~go
func first[T any](xs []T) (T,bool) {
 if len(xs)>0 { return xs[0],true }
 var zero T
 return zero,false
}
n,ok := first([]int{0}) // n=0, ok=true
~~~
~var zero T~ القيمة الصفرية للنوع: صفر أو نص فاضي أو struct صفرية. bool بتقول هل فيه عنصر فعلاً؛ الصفر لوحده مش علامة للفاضي.

## تحرير المراجع
~~~go
slots := []*int{new(int)}
slots[0] = nil
slots = slots[1:]
~~~
التقطيع لوحده ممكن يسيب العنصر القديم في المصفوفة الخلفية. امسح خانة العنصر قبل استبعاده. لما الطابور يفضى رجّع الـ slice لـ nil، عشان المصفوفة القديمة تتحرر؛ متخزّنش prefix مستهلك في field تاني.

## شكل الاستدعاء
~~~go
var q Queue[string]
q.Admit("medic")
first,ok:=q.Front()
fmt.Println(first,ok,q.Len()) // medic true 1
last,ok:=q.Depart()
fmt.Println(last,ok,q.Len()) // medic true 0
~~~
المثال بيختبر API فقط، مش تنفيذ التخزين. نفس الاستدعاءات تشتغل على Queue[int] أو أي struct.

## أخطاء شائعة
~~~go
last := q.waiting[len(q.waiting)-1] // ده آخر واحد، وممكن panic لو فاضي
~~~
ده ترتيب stack مش بوابة. افحص الطول قبل الفهرسة. ~Front~ استعلام: متشيلش عنصر وإنت بتعرض أول الواصلين.

## المهمة
نفّذ ~Queue[T any]~؛ ~var q Queue[int]~ شغال فوراً. ~Admit(v)~ يضيف آخر الطابور، ~Depart() (T,bool)~ يشيل أول عنصر، و~Front() (T,bool)~ يقرأه بس. الاتنين يرجّعوا zero,false لو فاضي من غير panic. ~Len() int~ العدد الحالي.

بعد Admit(4) وAdmit(9)، Depart ترجع 4,true وFront ترجع 9,true من غير تغيير Len. امسح المراجع المستهلكة، وتخلّص من التخزين لما يفضى. العقد تسلسلي ومش آمن لتزامن؛ مفيش Report في المختبر.
`, Hints: []string{"الـ nil slice جاهزة لـ append.", "اقرأ أول خانة، وامسحها بقيمة صفرية من T قبل التقطيع.", "Front متغيّرش التخزين؛ جرّب فاضي ثم عنصر واحد ثم إعادة استعمال."},
		En: &LevelText{Title: "Airlock Queue", Goal: "Build a generic queue that preserves arrival order", Files: []LabFile{{Name: "main.go", Starter: airlockQueueStarter, ReadOnly: false}}, Lesson: `# First Arrival, First Departure

The previous stack used LIFO. An airlock needs FIFO: the first passenger admitted is the first to depart, regardless of the passenger data type.

## Before you write
- Can the zero value represent an empty queue without a constructor?
- How will you distinguish a zero-valued passenger from no passenger?

## A generic type with one behavior
~~~go
type Tray[T any] struct { Slots []T }
func (t *Tray[T]) Add(v T) { t.Slots = append(t.Slots,v) }
var tray Tray[string]
tray.Add("lens")
~~~
~T any~ allows any type. ~Tray[string]~ supplies a concrete type, and the receiver must include ~[T]~. A nil slice accepts append, so initialization is unnecessary.

## A result and its presence
~~~go
func first[T any](xs []T) (T,bool) {
 if len(xs)>0 { return xs[0],true }
 var zero T
 return zero,false
}
n,ok := first([]int{0}) // n=0, ok=true
~~~
~var zero T~ gives the type's zero value: zero, empty text, or a zero struct. The bool records presence; zero alone cannot represent emptiness.

## Release references
~~~go
slots := []*int{new(int)}
slots[0] = nil
slots = slots[1:]
~~~
Slicing alone can leave an old reference in the backing array. Clear the consumed slot before slicing it away. Reset the slice to nil when empty so the old array can be reclaimed; do not retain a consumed prefix in another field.

## Calling the contract
~~~go
var q Queue[string]
q.Admit("medic")
first,ok:=q.Front()
fmt.Println(first,ok,q.Len()) // medic true 1
last,ok:=q.Depart()
fmt.Println(last,ok,q.Len()) // medic true 0
~~~
This exercises the API rather than implementing storage. The same calls work for Queue[int] or any struct.

## Common mistakes
~~~go
last := q.waiting[len(q.waiting)-1] // Last arrival; panics when empty
~~~
That is stack order. Check length before indexing. ~Front~ is a query and must leave the queue intact.

## Your mission
Implement ~Queue[T any]~ so ~var q Queue[int]~ works immediately. ~Admit(v)~ appends, ~Depart() (T,bool)~ removes the first passenger, and ~Front() (T,bool)~ only reads it. Both return zero,false when empty without panic. ~Len() int~ is the current count.

After Admit(4) and Admit(9), Depart returns 4,true and Front returns 9,true without changing Len. Clear consumed references and release storage when empty. The contract is sequential, with no concurrent safety requirement and no Report.
`, Hints: []string{"A nil slice is ready for append.", "Read the first slot, clear it using a zero T, then slice it away.", "Front must not modify storage; try empty, singleton, and reuse."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: airlockQueueHidden},
	})
	add(&Level{ID: "14-5-archive-tree", Chapter: 14, Kind: "lab", Title: "شجرة الأرشيف", Goal: "خزّن سجلات في شجرة generic بمقارنة وترتيب محددين", Requires: []string{"generictype", "recursion"}, Bridge: []int{28},
		Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: archiveTreeStarter, ReadOnly: false}},
		Lesson: `# أرشيف بيتفرّع

بدل ما تقارن بكل سجل، الشجرة بتختار فرع أصغر أو أكبر. هنا هنتعلّم recursion لأول مرة: دالة بتحل جزء أصغر من نفس المشكلة، وبتنادي نفسها عليه.

## قبل ما تكتب
- أنهي حالة توقف الاستدعاء من غير ما تنادي نفسك تاني؟
- سجلين متساويين بالمقارنة: تحتفظ بالأول ولا تستبدله؟

## comparator بدل افتراض نوع
~~~go
compare := func(a,b int) int {
 if a<b {return -1}; if a>b {return 1}; return 0
}
fmt.Println(compare(3,8)) // -1
~~~
زي SortFunc، نتيجة سالبة معناها a يسبق b، وصفر معناها متساويين. ~func(T,T) int~ تسمح بنفس الشجرة لرقم أو نص أو struct. المقارنة لازم تبقى ثابتة ومتسقة وتحدد ترتيب كلي؛ مش هنصلّح comparator متناقضة.

## recursion بمثال مش شجرة
~~~go
func countdown(n int) int {
 if n<=0 {return 0}
 return 1+countdown(n-1)
}
fmt.Println(countdown(3)) // 3
~~~
~n<=0~ حالة التوقف. كل نداء بيقرّبنا منها. بعد ما النداء الداخلي يرجع، الخارجي بيكمل. من غير حالة توقف هتستهلك الـ stack. هنا العمق صغير ومش محتاج goroutines.

## عقد generic وفروع
~~~go
type Branch[T any] struct { Entry T; Left,Right *Branch[T] }
b := &Branch[string]{Entry:"log"}
fmt.Println(b.Left==nil) // true
~~~
الفرع الفاضي nil. القيم اللي مقارنةً تسبق العقدة تروح Left، واللي بعدها Right. عند الإدخال رجّع جذر الفرع بعد تعديله، عشان الأب يحتفظ بالوصلة حتى لو كان الفرع nil.

## ترتيب الزيارة
~~~text
visit(left) → current → visit(right)
~~~
ده in-order: يعرض السجلات بترتيب المقارنة. تجميعه في slice جديدة يحافظ على استقلال النتيجة. الإدخال المرتّب ممكن يعمل سلسلة طويلة؛ الشجرة مش balanced.

## شكل الاستدعاء
~~~go
var logs Archive[int]
err:=logs.Order(compare)
added,err:=logs.Store(7)
fmt.Println(added,err,logs.Contains(7),logs.Records()) // true <nil> true [7]
~~~
~compare~ هي الدالة اللي فوق. ده استخدام الأرشيف بعد التنفيذ، مش خوارزمية بناء الشجرة.

## أخطاء شائعة
~~~go
func countdown(n int) int { return countdown(n) } // مفيش تقدّم أو توقف
~~~
ومتستخدمش ~record < node.record~ مع T any: النوع ممكن يبقى struct. استخدم المقارنة المحقونة.

## المهمة
نفّذ ~Archive[T any]~. ~Order(func(T,T) int) error~ يحدد المقارنة في أرشيف فاضي؛ nil أو محاولة تغييرها بعد وجود سجلات ترجع error من غير تعديل. ~Store(record) (bool,error)~ يحتاج مقارنة، ويرجّع true لو أضاف، وfalse,nil لو تكرار حسب مقارنة صفر، مع الاحتفاظ بالسجل الأول.

~Contains(record) bool~ يرجّع وجود سجل مساوٍ بالمقارنة. ~Records() []T~ نسخة مستقلة in-order. الـ zero value تقبل البحث والعرض وترجع فاضي؛ Store من غير Order ترجع خطأ من غير panic أو تعديل. استخدم دالة recursive للإدخال أو الزيارة. العقد تسلسلي؛ مفيش Report.
`, Hints: []string{"ابدأ بالفرع nil كحالة توقف.", "نتيجة المقارنة صفر معناها تكرار: متغيّرش السجل القديم.", "في الزيارة اجمع الشمال ثم العقدة ثم اليمين؛ متخزّنش الناتج كـ slice مشتركة."},
		En: &LevelText{Title: "Archive Tree", Goal: "Store records in a generic tree with a defined comparison order", Files: []LabFile{{Name: "main.go", Starter: archiveTreeStarter, ReadOnly: false}}, Lesson: `# A Branching Archive

Instead of comparing with every record, a tree selects a smaller or larger branch. This lesson introduces recursion: a function solves a smaller instance of its problem by calling itself.

## Before you write
- Which case stops without another recursive call?
- When two records compare equal, should you retain the first or replace it?

## A comparator instead of a type assumption
~~~go
compare := func(a,b int) int {
 if a<b {return -1}; if a>b {return 1}; return 0
}
fmt.Println(compare(3,8)) // -1
~~~
As with SortFunc, a negative result places a before b, and zero means equality. ~func(T,T) int~ lets the same tree accept numbers, strings, or structs. The comparator must be stable, consistent, and define a total order; contradictory comparators are outside the contract.

## Recursion with a non-tree example
~~~go
func countdown(n int) int {
 if n<=0 {return 0}
 return 1+countdown(n-1)
}
fmt.Println(countdown(3)) // 3
~~~
~n<=0~ is the base case. Each call approaches it. When the inner call returns, the outer call continues. Without a base case the stack is exhausted. These small depths need no goroutines.

## Generic nodes and branches
~~~go
type Branch[T any] struct { Entry T; Left,Right *Branch[T] }
b := &Branch[string]{Entry:"log"}
fmt.Println(b.Left==nil) // true
~~~
An empty branch is nil. Records comparing before the node go Left, and later ones go Right. Return the branch root after insertion so its parent retains the link even when the branch was nil.

## Visit order
~~~text
visit(left) → current → visit(right)
~~~
This is in-order traversal and yields comparator order. Collect into a fresh slice to keep results independent. Sorted insertions may produce a long chain; this tree is not balanced.

## Calling the contract
~~~go
var logs Archive[int]
err:=logs.Order(compare)
added,err:=logs.Store(7)
fmt.Println(added,err,logs.Contains(7),logs.Records()) // true <nil> true [7]
~~~
~compare~ is the function above. This uses the completed archive rather than implementing its tree algorithm.

## Common mistakes
~~~go
func countdown(n int) int { return countdown(n) } // No progress or base case
~~~
Do not use ~record < node.record~ for T any: it might be a struct. Use the injected comparator.

## Your mission
Implement ~Archive[T any]~. ~Order(func(T,T) int) error~ configures an empty archive. A nil comparator or reordering a populated archive returns an error without mutation. ~Store(record) (bool,error)~ requires a comparator, returns true for insertion, or false,nil for comparator equality, retaining the first record.

~Contains(record) bool~ reports comparator equality. ~Records() []T~ returns an independent in-order snapshot. The zero value supports empty reads; Store before Order returns an error without panic or mutation. Use a recursive function for insertion or traversal. This is a sequential contract with no Report.
`, Hints: []string{"Start with a nil branch as the base case.", "Zero comparison means duplicate: keep the original record.", "Traverse left, node, right; do not expose a shared result slice."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: archiveTreeHidden},
	})
}

const airlockQueueStarter = `package main

type Queue[T any] struct{ waiting []T }

func (q *Queue[T]) Admit(passenger T) {}
func (q *Queue[T]) Depart() (T, bool) { var zero T; return zero, false }
func (q *Queue[T]) Front() (T, bool)  { var zero T; return zero, false }
func (q *Queue[T]) Len() int          { return 0 }
func main()                           {}
`
const airlockQueueHidden = `package main

import (
	"io"
	"orbital/colony/levels"
	"testing"
)

func hiddenFIFO[T comparable](t *testing.T, passengers []T) {
	t.Helper()
	var q Queue[T]
	var zero T
	if got, ok := q.Depart(); ok || got != zero || q.Len() != 0 {
		t.Fatal("empty departure")
	}
	if got, ok := q.Front(); ok || got != zero {
		t.Fatal("empty front")
	}
	var want []T
	world := levels.Start("14-4-airlock-queue", io.Discard, 144)
	r := world.Rand
	for step := 0; step < 600; step++ {
		if len(want) == 0 || r.Intn(2) == 0 {
			v := passengers[world.Data.SensorReadings[step%256]%len(passengers)]
			q.Admit(v)
			want = append(want, v)
		} else {
			got, ok := q.Front()
			if !ok || got != want[0] || q.Len() != len(want) {
				t.Fatal("front mutated or wrong")
			}
			got, ok = q.Depart()
			if !ok || got != want[0] {
				t.Fatal("FIFO")
			}
			want = want[1:]
		}
		if q.Len() != len(want) {
			t.Fatal("length transition")
		}
	}
	for _, v := range want {
		if got, ok := q.Depart(); !ok || got != v {
			t.Fatal("drain")
		}
	}
	if _, ok := q.Depart(); ok || q.Len() != 0 {
		t.Fatal("empty after drain")
	}
	q.Admit(passengers[0])
	if got, ok := q.Depart(); !ok || got != passengers[0] {
		t.Fatal("reuse")
	}
}
func TestHiddenZeroValueAndFIFO(t *testing.T) {
	hiddenFIFO(t, []int{0, -8, 17})
	hiddenFIFO(t, []string{"", "crew", "🚀"})
	hiddenFIFO(t, []struct {
		Gate int
		Name string
	}{{1, "a"}, {2, "b"}})
}
`

const archiveTreeStarter = `package main

import "errors"

type archiveNode[T any] struct {
	record      T
	left, right *archiveNode[T]
}
type Archive[T any] struct {
	root    *archiveNode[T]
	compare func(T, T) int
}

func (a *Archive[T]) Order(compare func(T, T) int) error { return errors.New("order unfinished") }
func (a *Archive[T]) Store(record T) (bool, error)       { return false, errors.New("store unfinished") }
func (a *Archive[T]) Contains(record T) bool             { return false }
func (a *Archive[T]) Records() []T                       { return nil }
func main()                                              {}
`
const archiveTreeHidden = `package main

import (
	"cmp"
	"io"
	"orbital/colony/levels"
	"slices"
	"testing"
)

func hiddenArchive[T any](t *testing.T, inputs []T, compare func(T, T) int) {
	t.Helper()
	var archive Archive[T]
	if archive.Contains(inputs[0]) || len(archive.Records()) != 0 {
		t.Fatal("zero value read")
	}
	if _, err := archive.Store(inputs[0]); err == nil {
		t.Fatal("missing comparator")
	}
	if err := archive.Order(nil); err == nil {
		t.Fatal("nil order")
	}
	if err := archive.Order(compare); err != nil {
		t.Fatal(err)
	}
	var want []T
	for _, record := range inputs {
		duplicate := slices.ContainsFunc(want, func(old T) bool { return compare(old, record) == 0 })
		added, err := archive.Store(record)
		if err != nil || added == duplicate {
			t.Fatal("duplicate policy")
		}
		if !duplicate {
			want = append(want, record)
		}
		if !archive.Contains(record) {
			t.Fatal("search miss after insertion")
		}
	}
	slices.SortFunc(want, compare)
	got := archive.Records()
	if !slices.EqualFunc(got, want, func(a, b T) bool { return compare(a, b) == 0 }) {
		t.Fatalf("traversal %v want %v", got, want)
	}
	got[0] = inputs[len(inputs)-1]
	if !slices.EqualFunc(archive.Records(), want, func(a, b T) bool { return compare(a, b) == 0 }) {
		t.Fatal("snapshot aliases tree")
	}
	if archive.Order(compare) == nil {
		t.Fatal("cannot reorder a populated tree")
	}
}
func TestHiddenArchiveOrdering(t *testing.T) {
	world := levels.Start("14-5-archive-tree", io.Discard, 145)
	xs := []int{-5, 0, 5}
	for i := 0; i < 200; i++ {
		xs = append(xs, world.Data.SensorReadings[i]%300-150)
	}
	hiddenArchive(t, xs, cmp.Compare[int])
	hiddenArchive(t, xs, func(a, b int) int { return cmp.Compare(b, a) })
	for _, xs := range [][]int{{3}, {1, 2, 3, 4, 5}, {5, 4, 3, 2, 1}, {2, 2, 2}} {
		hiddenArchive(t, xs, cmp.Compare[int])
	}
	hiddenArchive(t, []string{"z", "a", "z", "🚀"}, cmp.Compare[string])
	type memo struct {
		Rank int
		Text string
	}
	hiddenArchive(t, []memo{{2, "first"}, {1, "low"}, {2, "ignored"}}, func(a, b memo) int { return cmp.Compare(a.Rank, b.Rank) })
	var a Archive[int]
	_ = a.Order(cmp.Compare[int])
	_, _ = a.Store(4)
	if a.Contains(3) || a.Contains(5) {
		t.Fatal("missing search")
	}
	var b Archive[memo]
	_ = b.Order(func(a, b memo) int { return cmp.Compare(a.Rank, b.Rank) })
	_, _ = b.Store(memo{2, "first"})
	_, _ = b.Store(memo{2, "ignored"})
	if b.Records()[0].Text != "first" {
		t.Fatal("duplicate replaced original")
	}
}
`
