package levels

import (
	"fmt"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

// scatterRocks fills ~density of the map with rock, keeping (0, 0) clear.
func scatterRocks(w *sim.World, density float64) {
	w.Each(func(x, y int, t *sim.Tile) {
		if x+y > 0 && w.Rand.Float64() < density {
			t.Ground = sim.Rock
		}
	})
}

// farthest returns a reachable tile at the largest distance (by dist).
func farthest(w *sim.World, dist map[sim.Point]int) (sim.Point, int) {
	best, bd := sim.Point{}, -1
	w.Each(func(x, y int, _ *sim.Tile) {
		p := sim.Point{X: x, Y: y}
		if d, ok := dist[p]; ok && (d > bd || d == bd && w.Rand.Intn(2) == 0) {
			best, bd = p, d
		}
	})
	return best, bd
}

// bfsPathCost walks one fewest-moves path from a to b and sums its tick cost.
func bfsPathCost(w *sim.World, a, b sim.Point) int {
	parent := map[sim.Point]sim.Point{a: a}
	q := []sim.Point{a}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		for _, n := range w.Neighbors(p) {
			if _, ok := parent[n]; !ok {
				parent[n] = p
				q = append(q, n)
			}
		}
	}
	cost := 0
	for p := b; p != a; p = parent[p] {
		cost += w.StepCost(p)
	}
	return cost
}

// shakerCost is a reference plan: measure every crystal walking east, then
// cocktail-sort the known sizes, walking and swapping as needed.
func shakerCost(sizes []int) int {
	s := append([]int(nil), sizes...)
	n := len(s)
	cost, at := 2*n-1, n-1
	walk := func(to int) {
		cost += max(at-to, to-at)
		at = to
	}
	for swapped := true; swapped; {
		swapped = false
		for i := n - 2; i >= 0; i-- { // westward pass
			if s[i] > s[i+1] {
				walk(i)
				s[i], s[i+1] = s[i+1], s[i]
				cost++
				swapped = true
			}
		}
		for i := 0; i < n-1; i++ { // eastward pass
			if s[i] > s[i+1] {
				walk(i)
				s[i], s[i+1] = s[i+1], s[i]
				cost++
				swapped = true
			}
		}
	}
	return cost
}

func init() {
	add(&Level{
		En: &LevelText{
			Title: "Cave Maze",
			Goal:  "Reach the crystal by the shortest route through the maze. colony.Grid() gives you the map for free",
			Lesson: `# BFS: The Shortest Route Through a Maze 🗺️

**Breadth-First Search** explores the map **layer by layer**: all tiles at distance 1, then 2, then 3... The first time it reaches a tile, it has found a **shortest route**.

~~~go
start := colony.Point{X: 0, Y: 0}
parent := map[colony.Point]colony.Point{start: start}
queue := []colony.Point{start}

for len(queue) > 0 {
	p := queue[0]
	queue = queue[1:]        // FIFO: first in, first out
	for _, n := range neighbors(p) {
		if _, seen := parent[n]; seen {
			continue
		}
		parent[n] = p        // We reached n from p
		queue = append(queue, n)
	}
}
~~~

## Reconstruct the route
~~~go
var path []colony.Point
for p := goal; p != start; p = parent[p] {
	path = append(path, p)
}
slices.Reverse(path) // Now ordered from start to goal
~~~

## neighbors
The four neighboring tiles that lie inside the map and are not ~colony.Rock~.

## DFS or BFS?
- **DFS** (the stack, chapter 14): follows a path to its end first. Good for exploring an unknown map.
- **BFS** (the queue): finds a **shortest** route when every step costs the same.

## Your mission
Rock is impassable. The map changes each run. ⭐⭐⭐ = exactly the shortest route.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	grid := colony.Grid() // grid[x][y] = colony.Rock or another terrain
	w, h := len(grid), len(grid[0])
	fmt.Println("Map", w, "×", h)

	goal := colony.Targets()[0] // The crystal's location
	fmt.Println("Goal", goal)

	// 1. BFS from (0, 0): a queue + a map recording each tile's parent
	// 2. Follow parents from goal back to start, then reverse the route
	// 3. Follow the route step by step and harvest
	_ = w
	_ = h
}
`,
			Hints: []string{
				"A Go queue: ~q := []colony.Point{start}~; remove with ~p := q[0]; q = q[1:]~ and add with ~q = append(q, n)~.",
				"~parent := map[colony.Point]colony.Point{}~: the first time you reach n from p, record ~parent[n] = p~. Do not add n again if it is already in parent.",
				"To follow the route, compare each next tile with your position: larger x means East, smaller x West, larger y North, smaller y South.",
			},
		},
		ID: "15-1-bfs", Chapter: 15, Title: "متاهة الكهوف",
		Goal: "وصل للكريستال بأقصر طريق في المتاهة. colony.Grid() بتديك الخريطة ببلاش",
		Par:  40, Gold: 25, MaxTicks: 400,
		Setup: func(w *sim.World) {
			var goal sim.Point
			w.RockBlocks = true // BFS below must respect the walls
			for {
				w.Init(12, 8, sim.Regolith)
				scatterRocks(w, 0.32)
				var d int
				goal, d = farthest(w, w.BFS(sim.Point{}))
				if d >= 14 {
					w.Gold = d + 1
					break
				}
			}
			w.Targets = []sim.Point{goal}
			t := w.At(goal.X, goal.Y)
			t.Entity, t.Size = sim.Crystal, 7
			w.Par = w.Gold + 12
			w.Check = func(w *sim.World) (bool, string) {
				if w.At(goal.X, goal.Y).Entity != sim.Crystal {
					return true, i18n.T("وصلت للكريستال", "You reached the crystal")
				}
				return false, fmt.Sprintf(i18n.T("الكريستال عند %v", "Crystal at %v"), goal)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	grid := colony.Grid() // grid[x][y] = colony.Rock أو غيره
	w, h := len(grid), len(grid[0])
	fmt.Println("الخريطة", w, "×", h)

	goal := colony.Targets()[0] // مكان الكريستال
	fmt.Println("الهدف", goal)

	// 1. BFS من (0, 0): queue + map بتفتكر جيت لكل خانة منين
	// 2. ارجع من goal للبداية عن طريق الـ parents، واعكس الطريق
	// 3. امشي فيه خطوة خطوة، واحصد
	_ = w
	_ = h
}
`,
		Hints: []string{
			"queue في Go: ~q := []colony.Point{start}~، تطلّع منها ~p := q[0]; q = q[1:]~، وتضيف ~q = append(q, n)~.",
			"~parent := map[colony.Point]colony.Point{}~: أول مرة توصل لخانة n من p سجّل ~parent[n] = p~. لو n موجودة في parent متضيفهاش تاني.",
			"عشان تمشي في الطريق: لكل خانة جاية قارنها بمكانك: x أكبر يبقى East، وأصغر West، وy أكبر North، وأصغر South.",
		},
		Lesson: `# BFS: أقصر طريق في متاهة 🗺️

**Breadth-First Search** بيستكشف الخريطة **طبقة طبقة**: الأول كل الخانات اللي على بعد 1، وبعدين 2، وبعدين 3… فأول مرة يوصل لخانة، يكون وصلها **بأقصر طريق**.

~~~go
start := colony.Point{X: 0, Y: 0}
parent := map[colony.Point]colony.Point{start: start}
queue := []colony.Point{start}

for len(queue) > 0 {
	p := queue[0]
	queue = queue[1:]        // FIFO: أول داخل أول خارج
	for _, n := range neighbors(p) {
		if _, seen := parent[n]; seen {
			continue
		}
		parent[n] = p        // جينا n من p
		queue = append(queue, n)
	}
}
~~~

## ارجع بالطريق
~~~go
var path []colony.Point
for p := goal; p != start; p = parent[p] {
	path = append(path, p)
}
slices.Reverse(path) // دلوقتي من البداية للهدف
~~~

## neighbors
الجيران الأربعة اللي جوه الخريطة ومش ~colony.Rock~.

## DFS ولا BFS؟
- **DFS** (الـ Stack، الفصل 14): بيروح لآخر طريق الأول. كويس للاستكشاف لما متعرفش الخريطة.
- **BFS** (الـ Queue): بيلاقي **أقصر** طريق لما كل خطوة بنفس التكلفة.

## المهمة
الصخور مقفولة. الخريطة بتتغير كل مرة. ⭐⭐⭐ = أقصر طريق بالظبط.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Line Up the Crystal",
			Goal:  "Sort crystals from smallest to largest (left to right) using colony.Swap",
			Lesson: `# Sorting When Every Step Costs 🔢

~slices.Sort~ sorts in memory almost instantly. Here **comparisons** are free if you remember sizes, but each **move** and **swap** costs a tick. An algorithm that **reduces swaps and walking** wins.

## Bubble sort

~~~go
for swapped := true; swapped; {
	swapped = false
	for i := 0; i < n-1; i++ {
		if s[i] > s[i+1] {
			s[i], s[i+1] = s[i+1], s[i] // Go swaps values in one line
			swapped = true
		}
	}
}
~~~

## How many swaps do you need?
Each adjacent swap fixes exactly one **inversion** (a pair in the wrong order). The minimum number of swaps is therefore the inversion count; no adjacent-swap method can do better.

## Cocktail sort 🍸
The same idea as bubble sort, but **back and forth**: one pass right, one pass left. This saves lots of walking because you need not return to the start each time.

## In real programs
- ~slices.Sort(s)~: fast (pdqsort, O(n log n)).
- ~slices.SortFunc(s, func(a, b T) int { ... })~: for a custom order.
- ~slices.SortStableFunc~: when equal elements must retain their relative order.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	sizes := make([]int, 8)
	for x := 0; x < 8; x++ {
		sizes[x] = colony.Measure()
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	fmt.Println(sizes)

	// colony.Swap(colony.West) swaps your crystal with the one to your left
	// Sort them, updating sizes too so you remember the order
}
`,
			Hints: []string{
				"Bubble sort: inspect adjacent pairs and swap when the left is larger than the right. Repeat until a whole pass makes no swaps.",
				"Keep ~sizes~ in memory consistent with the field: after every Swap, also swap in the slice: ~sizes[i], sizes[i+1] = sizes[i+1], sizes[i]~.",
				"For more stars, travel back and forth (cocktail sort) instead of returning to the start on each pass.",
			},
		},
		ID: "15-2-sorting", Chapter: 15, Title: "رصّ الكريستال",
		Goal: "رتّب الكريستال من الأصغر للأكبر (شمال ← يمين) بـ colony.Swap",
		Par:  70, Gold: 45, MaxTicks: 500,
		Setup: func(w *sim.World) {
			w.Init(8, 2, sim.Regolith)
			var sizes []int
			for {
				sizes = w.Rand.Perm(8)
				if !sortedInts(sizes) {
					break
				}
			}
			for x, s := range sizes {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, s+1
			}
			w.Gold = shakerCost(sizes)
			w.Par = w.Gold + 25
			w.Check = func(w *sim.World) (bool, string) {
				row := make([]int, 8)
				wrong := 0
				for x := range row {
					row[x] = w.At(x, 0).Size
					if row[x] != x+1 {
						wrong++
					}
				}
				return wrong == 0, fmt.Sprintf(i18n.T("الصف: %v · في مكانهم الصح: %d / 8", "Row: %v · In the right place: %d / 8"), row, 8-wrong)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	sizes := make([]int, 8)
	for x := 0; x < 8; x++ {
		sizes[x] = colony.Measure()
		if x < 7 {
			colony.Move(colony.East)
		}
	}
	fmt.Println(sizes)

	// colony.Swap(colony.West) بتبدّل الكريستال اللي تحتك مع اللي على شمالك
	// رتّبهم، وحدّث sizes إنت كمان عشان تفتكر الترتيب
}
`,
		Hints: []string{
			"Bubble sort: لف على الأزواج المتجاورة، ولو الشمال أكبر من اليمين بدّلهم. كرّر لحد ما لفّة كاملة تعدي من غير ولا تبديل.",
			"خلّي ~sizes~ في الذاكرة مطابقة للأرض: كل ما تعمل Swap، بدّل في الـ slice كمان: ~sizes[i], sizes[i+1] = sizes[i+1], sizes[i]~.",
			"للنجوم: الحركة غالية. امشي ذهاب وعودة (cocktail sort) بدل ما ترجع للأول كل لفّة.",
		},
		Lesson: `# الترتيب: لما كل خطوة بتكلّف 🔢

~slices.Sort~ بترتب في الذاكرة في لحظة. هنا كل **مقارنة** ببلاش (لو حافظ الأحجام)، بس كل **حركة** وكل **تبديل** بتاخد tick. يعني الخوارزمية اللي **بتقلل التبديلات والمشي** هي اللي تكسب.

## Bubble sort

~~~go
for swapped := true; swapped; {
	swapped = false
	for i := 0; i < n-1; i++ {
		if s[i] > s[i+1] {
			s[i], s[i+1] = s[i+1], s[i] // Go بيبدّل في سطر واحد
			swapped = true
		}
	}
}
~~~

## كام تبديل محتاج؟
كل تبديل لجارين بيصلّح **inversion** واحدة بالظبط (زوج مترتب غلط). يعني عدد التبديلات الأقل = عدد الـ inversions، ومحدش يقدر يعمل أقل منه بتبديل الجيران.

## Cocktail sort 🍸
نفس bubble بس **ذهاب وعودة**: لفّة يمين ولفّة شمال. هنا ده بيوفّر مشي كتير، لأنك مش محتاج ترجع لأول الصف كل مرة.

## في الشغل الحقيقي
- ~slices.Sort(s)~: سريعة (pdqsort، يعني O(n log n)).
- ~slices.SortFunc(s, func(a, b T) int { ... })~: لو محتاج ترتيب مخصص.
- ~slices.SortStableFunc~: لو العناصر المتساوية لازم تحافظ على ترتيبها.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Slippery Ice",
			Goal:  "Reach the crystal in the least time. Each ice tile costs 4 ticks!",
			Lesson: `# Dijkstra: The Shortest **Time**, Not the Fewest **Steps** ⛸️

In 15-1, each step cost one tick, so fewer steps meant less time. Here ice costs **4 ticks**. Sometimes a 3-step detour around ice is **faster** than crossing it.

## The idea
Just like BFS, except you take the tile with **the lowest cost** next (rather than the first one added):

~~~go
dist := map[colony.Point]int{start: 0}
done := map[colony.Point]bool{}
for {
	// Find the smallest dist not yet done
	best, found := colony.Point{}, false
	for p, d := range dist {
		if !done[p] && (!found || d < dist[best]) {
			best, found = p, true
		}
	}
	if !found {
		break
	}
	done[best] = true
	for _, n := range neighbors(best) {
		nd := dist[best] + cost(ground(n))
		if old, ok := dist[n]; !ok || nd < old {
			dist[n] = nd
			parent[n] = best
		}
	}
}
~~~

## Faster: container/heap
Instead of searching for the minimum each time (O(n²)), a **priority queue** extracts it in O(log n). Go provides ~container/heap~: define a type implementing ~heap.Interface~ (Len, Less, Swap, Push, and Pop). A good exercise in interfaces!

## ⭐ Stars
- Ordinary BFS: ⭐⭐
- Dijkstra: ⭐⭐⭐`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// cost is the cost of entering a tile with ground g
func cost(g colony.Ground) int {
	if g == colony.Ice {
		return 4
	}
	return 1
}

func main() {
	grid := colony.Grid()
	fmt.Println(len(grid), len(grid[0]), cost(grid[1][0]))
	// BFS minimizes steps... but not time. You need Dijkstra
}
`,
			Hints: []string{
				"Dijkstra is like BFS, but each time you choose the unfinished tile with **the lowest cost so far**.",
				"The simplest version: ~dist map[colony.Point]int~ and ~done map[colony.Point]bool~. Each iteration, find the lowest dist not in done. That is enough for this small map.",
				"The crystal is at ~colony.Targets()[0]~. If you have a solution for 15-1, reuse its neighbors function and route-following logic.",
			},
		},
		ID: "15-3-dijkstra", Chapter: 15, Title: "الجليد الزلق",
		Goal: "وصل للكريستال بأقل وقت. الجليد بياخد 4 tick لكل خانة!",
		Par:  40, Gold: 30, MaxTicks: 500,
		Setup: func(w *sim.World) {
			var goal sim.Point
			for {
				w.Init(10, 7, sim.Regolith)
				w.RockBlocks = true
				w.MoveCost = map[sim.Ground]int{sim.Ice: 4}
				w.Each(func(x, y int, t *sim.Tile) {
					switch r := w.Rand.Float64(); {
					case x+y == 0:
					case r < 0.12:
						t.Ground = sim.Rock
					case r < 0.5:
						t.Ground = sim.Ice
					}
				})
				var d int
				goal, d = farthest(w, w.BFS(sim.Point{}))
				if d < 10 {
					continue
				}
				cheap := w.Cheapest(sim.Point{})[goal]
				if naive := bfsPathCost(w, sim.Point{}, goal); naive-cheap >= 5 {
					w.Gold = cheap + 1
					w.Par = naive + 1
					break
				}
			}
			w.Targets = []sim.Point{goal}
			t := w.At(goal.X, goal.Y)
			t.Entity, t.Size = sim.Crystal, 7
			w.Check = func(w *sim.World) (bool, string) {
				if w.At(goal.X, goal.Y).Entity != sim.Crystal {
					return true, i18n.T("وصلت للكريستال", "You reached the crystal")
				}
				return false, fmt.Sprintf(i18n.T("الكريستال عند %v", "Crystal at %v"), goal)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// cost تكلفة إنك تدخل خانة أرضها g
func cost(g colony.Ground) int {
	if g == colony.Ice {
		return 4
	}
	return 1
}

func main() {
	grid := colony.Grid()
	fmt.Println(len(grid), len(grid[0]), cost(grid[1][0]))
	// BFS بيلاقي أقل عدد خطوات… بس مش أقل وقت. محتاج Dijkstra
}
`,
		Hints: []string{
			"Dijkstra: زي BFS، بس بدل الـ queue العادية، كل مرة خد الخانة اللي **تكلفتها لحد دلوقتي أقل** من كل اللي لسه متخلصتش.",
			"أسهل نسخة: ~dist map[colony.Point]int~ و~done map[colony.Point]bool~، وكل لفّة دوّر على أقل dist مش done. ده كفاية للخريطة الصغيرة دي.",
			"مكان الكريستال في ~colony.Targets()[0]~. ولو عندك حل 15-1، خد منه دالة neighbors ومنطق المشي في الطريق.",
		},
		Lesson: `# Dijkstra: أقصر **وقت** مش أقصر **طريق** ⛸️

في 15-1 كل خطوة كانت بـ tick واحد، فأقل خطوات = أقل وقت. هنا الجليد بياخد **4 tick**. ساعات تلف 3 خطوات حوالين الجليد وتطلع **أسرع** من إنك تعدّي من فوقه.

## الفكرة
زي BFS بالظبط، بس كل مرة بتطلّع الخانة اللي **تكلفتها أقل** (مش اللي دخلت الأول):

~~~go
dist := map[colony.Point]int{start: 0}
done := map[colony.Point]bool{}
for {
	// لاقي أقل dist لسه مش done
	best, found := colony.Point{}, false
	for p, d := range dist {
		if !done[p] && (!found || d < dist[best]) {
			best, found = p, true
		}
	}
	if !found {
		break
	}
	done[best] = true
	for _, n := range neighbors(best) {
		nd := dist[best] + cost(ground(n))
		if old, ok := dist[n]; !ok || nd < old {
			dist[n] = nd
			parent[n] = best
		}
	}
}
~~~

## أسرع: container/heap
بدل ما تدوّر على الأقل كل مرة (O(n²))، الـ **priority queue** بتطلّعهولك في O(log n). Go فيها ~container/heap~: بتعمل type يطبّق ~heap.Interface~ (Len وLess وSwap وPush وPop)، وده تمرين حلو على الـ interfaces.

## ⭐ النجوم
- BFS عادي: ⭐⭐
- Dijkstra: ⭐⭐⭐`,
	})
	add(&Level{ID: "15-4-deep-probe", Chapter: 15, Title: "المسبار العميق", Goal: "ابحث في الأعماق المرتبة وعدّ المنطقة المتصلة بالقاعدة", Requires: []string{"recursion"}, Bridge: []int{24}, Par: 0, Gold: 0, MaxTicks: 100,
		Starter: deepProbeStarter, Lesson: `# نص البحث وزيارة الفروع

الأرشيف علّمنا recursion. هنستخدمها دلوقتي للانتشار على خريطة، ومعاها binary search لقائمة أعماق مرتبة. الاتنين بيقلّلوا شغل مكرر بقواعد واضحة.

## قبل ما تكتب
- لو العمق متكرر، أنهي index مطلوب؟
- إيه اللي يمنع زيارة خانتين بيرجعوا لبعض للأبد؟

## binary search بحد يمين مستبعد
~~~go
xs:=[]int{2,5,5,9}
lo,hi:=0,len(xs)
for lo<hi {
 mid:=lo+(hi-lo)/2
 if xs[mid]<6 {lo=mid+1} else {hi=mid}
}
fmt.Println(lo) // 3: أول قيمة >= 6، مش تطابق
~~~
مجال البحث ~[lo,hi)~. لو القيمة أصغر من الهدف، استبعدها وكل اللي قبلها. غير كده احتفظ بالمنتصف كمرشح. عند التوقف lo أول قيمة مش أقل من الهدف، فافحص الحدود والمساواة عشان تفرّق التطابق عن عدم الوجود. ده يختار أول تكرار، ويشتغل للفاضي كمان.

## خريطة وذاكرة زيارة
~~~go
grid:=colony.Grid()
seen:=map[colony.Point]bool{}
p:=colony.Point{X:0,Y:0}
seen[p]=true
fmt.Println(grid[0][0],seen[p])
~~~
~Grid~ زي BFS: الأبعاد grid[x][y] مش العكس، والصخر ~colony.Rock~ حاجز. الـ map بتسجل الخانات اللي حسبتها، ومش محتاج تغيّر العالم نفسه.

## انتشار recursive: مثال خط واحد
~~~go
func openRun(open []bool,i int) int {
 if i<0 || i>=len(open) || !open[i] {return 0}
 return 1+openRun(open,i+1)
}
fmt.Println(openRun([]bool{true,true,false},0)) // 2
~~~
ده خط واحد بلا رجوع، مش حل الخريطة. في شبكة فيها دوائر، علّم الخانة visited قبل النداء على الجيران. التوقف: خارج الحدود أو صخر أو visited. النداء يرجّع عدد الخانات الجديدة، وجمع الأعداد مع1 يعدّ الخانة الحالية.

## أخطاء شائعة
~~~go
count := flood(grid,x+1,y,seen)
seen[colony.Point{X:x,Y:y}]=true // متأخر: الجار ممكن يرجعلك قبل العلامة
~~~
ومتعدّش كل الأرض المفتوحة: جيب معزول مش متصل بالقاعدة مش جزء من المنطقة المطلوبة. ومينفعش أول تطابق في المنتصف يوقف البحث لو فيه تكرار قبله.

## المهمة
اقرأ ~colony.Depths()~ المرتبة تصاعدياً، والهدف ~colony.Reading(0)~. نفّذ binary search في ~depthIndex~ ترجع أول index للهدف أو -1 لو مش موجود.

نفّذ ~flood~ recursive لحساب كل خانة غير Rock متصلة بـ(0,0) شمال/شرق/جنوب/غرب، وعدّ البداية مرة واحدة. ابعت تقرير واحد ~[]int{index,count}~، مثلاً [3 12]. متحصدش ومتعدّش الجيوب المعزولة. مدخلات المهمة فيها تكرارات وهدف ممكن يغيب؛ الدالتين يتعاملوا مع قائمة أو خريطة فاضية بأمان.
`, Hints: []string{"استخدم مجال [lo,hi) وكمّل ناحية الشمال عند المساواة.", "في flood افحص الحدود والصخر والزيارة قبل عدّ الخانة.", "علّم seen قبل الجيران الأربعة، واجمع نتايجهم مع1."},
		En: &LevelText{Title: "Deep Probe", Goal: "Search sorted depths and count the region connected to base", Starter: deepProbeStarter, Lesson: `# Halving the Search and Visiting Branches

The archive introduced recursion. Apply it to spreading across a map, alongside binary search over sorted depths. Both avoid repeated work through explicit rules.

## Before you write
- If the depth repeats, which index is required?
- What prevents neighboring cells from visiting each other forever?

## Binary search with an exclusive upper bound
~~~go
xs:=[]int{2,5,5,9}
lo,hi:=0,len(xs)
for lo<hi {
 mid:=lo+(hi-lo)/2
 if xs[mid]<6 {lo=mid+1} else {hi=mid}
}
fmt.Println(lo) // 3: first value >= 6, not an exact match
~~~
Search in ~[lo,hi)~. If the middle value is too small, discard it and everything before it. Otherwise retain the middle as a candidate. At termination, lo is the first value not below the target. Check bounds and equality to distinguish a match from absence. This selects the first duplicate and handles empty slices.

## A map and visit memory
~~~go
grid:=colony.Grid()
seen:=map[colony.Point]bool{}
p:=colony.Point{X:0,Y:0}
seen[p]=true
fmt.Println(grid[0][0],seen[p])
~~~
~Grid~ uses grid[x][y], as in BFS, and ~colony.Rock~ blocks passage. The map tracks counted cells without changing the world.

## Recursive spreading: a line example
~~~go
func openRun(open []bool,i int) int {
 if i<0 || i>=len(open) || !open[i] {return 0}
 return 1+openRun(open,i+1)
}
fmt.Println(openRun([]bool{true,true,false},0)) // 2
~~~
This is a one-way line rather than the map solution. A grid has cycles, so mark a cell visited before visiting neighbors. Stop for out-of-bounds, rock, or already-visited cells. Each call returns a new-cell count; add those counts and one for the current cell.

## Common mistakes
~~~go
count := flood(grid,x+1,y,seen)
seen[colony.Point{X:x,Y:y}]=true // Too late: a neighbor can return before the mark
~~~
Do not count all open terrain: isolated pockets are outside the connected region. Finding a match in the middle also cannot end binary search when earlier duplicates may exist.

## Your mission
Read sorted ascending ~colony.Depths()~ and target ~colony.Reading(0)~. Implement binary search in ~depthIndex~, returning the first matching index or -1 for a miss.

Implement recursive ~flood~ to count non-Rock cells connected to (0,0) by north/east/south/west steps, counting the start once. Send one ~[]int{index,count}~ report, for example [3 12]. Do not harvest or count isolated pockets. Inputs contain duplicates and may omit the target; both functions must safely handle empty depths or maps.
`, Hints: []string{"Use [lo,hi) and continue left on equality.", "In flood, check bounds, rock, and visit state before counting.", "Mark seen before the four neighbors, then add their counts plus one."}}, Setup: func(w *sim.World) {
			w.Init(7+w.Rand.Intn(4), 5+w.Rand.Intn(3), sim.Regolith)
			w.RockBlocks = true
			scatterRocks(w, 0.3)
			// Isolated open pocket, unreachable from the origin.
			w.At(w.W-1, w.H-1).Ground = sim.Regolith
			w.At(w.W-2, w.H-1).Ground = sim.Rock
			w.At(w.W-1, w.H-2).Ground = sim.Rock
			depth := 10 + w.Rand.Intn(30)
			n := 12 + w.Rand.Intn(8)
			for i := 0; i < n; i++ {
				depth += 1 + w.Rand.Intn(9)
				w.Data.Depths = append(w.Data.Depths, depth, depth)
			}
			target := w.Data.Depths[w.Rand.Intn(len(w.Data.Depths))]
			if w.Rand.Intn(3) == 0 {
				target = w.Data.Depths[len(w.Data.Depths)-1] + 1
			}
			w.Data.SensorReadings = []int{target}
			index := -1
			for i, d := range w.Data.Depths {
				if d == target {
					index = i
					break
				}
			}
			want := fmt.Sprint([]int{index, len(w.BFS(sim.Point{}))})
			w.Check = func(w *sim.World) (bool, string) {
				if w.Harvests == 0 && len(w.Reports) == 1 && w.Reports[0] == want {
					return true, i18n.T("أول عمق وعدد الخانات المتصلة صح", "First depth index and connected tile count are correct")
				}
				return false, i18n.T("مستني [أول موقع للعمق عدد الخانات المتصلة] من غير حصاد", "Waiting for [first-depth-index connected-tile-count] without harvesting")
			}
		}})
	add(&Level{ID: "15-5-rescue-priority", Chapter: 15, Title: "أولوية الإنقاذ", Goal: "اخدم أعلى أولوية حتى مع وصول موجات استغاثة جديدة", Requires: []string{"sel:heap.Push"}, Bridge: []int{90}, Par: 0, Gold: 0, MaxTicks: 100,
		Starter: rescuePriorityStarter, Lesson: `# طابور بيعيد ترتيب نفسه

الـ FIFO بيحفظ الوصول. الإنقاذ محتاج كل مرة يختار أعلى أولوية من الموجود دلوقتي؛ طلب جديد ممكن يسبق اللي مستني. الـ heap يحافظ على أفضل عنصر في الجذر من غير ترتيب الكل كل مرة.

## قبل ما تكتب
- إيه قاعدة التعادل عشان النتيجة تبقى حتمية؟
- هتدخل الوافدين الجدد قبل اختيار الطلب اللي بعده ولا بعده؟

## عقد container/heap
~~~go
type Scores []int
func (s Scores) Len() int { return len(s) }
func (s Scores) Less(i,j int) bool { return s[i]<s[j] }
func (s Scores) Swap(i,j int) { s[i],s[j]=s[j],s[i] }
func (s *Scores) Push(v any) { *s=append(*s,v.(int)) }
func (s *Scores) Pop() any {
 old:=*s; n:=len(old); last:=old[n-1]; *s=old[:n-1]; return last
}
~~~
استورد ~container/heap~. الـ interface محتاج Len وLess وSwap وPush وPop. Less بتقول مين يسبق؛ المثال أصغر رقم الأول، مش أولوية المهمة. pointer receivers لـPush/Pop عشان يغيّروا طول الـ slice. ~v.(int)~ type assertion اتعلمناها مع الإشارات؛ heap هياخد النوع اللي إحنا بندخله بس.

## نادِ الحزمة مش الـ method
~~~go
s:=&Scores{}
heap.Push(s,8)
heap.Push(s,3)
fmt.Println(heap.Pop(s).(int)) // 3
~~~
~heap.Push~ تضيف وتصلّح ترتيب الشجرة؛ ~heap.Pop~ تنقل أفضل عنصر لآخر الـ slice وبعدين تنادي Pop بتاعتك. عشان كده الـ method تشيل الآخر، مش الأول. ~heap.Init~ تجهّز slice موجودة؛ الفاضية جاهزة من غيرها. افحص Len قبل Pop.

## موجات بتوصل
~~~go
arrivals:=colony.RescueArrivals()
for _,request:=range arrivals {
 fmt.Println(request.ID,request.Priority)
}
ok:=colony.Rescue(arrivals[0].ID) // مثال فقط، مش قاعدة الاختيار
fmt.Println(ok)
~~~
~RescueRequest~ فيها ID فريد وPriority int. ~RescueArrivals~ نسخة من الموجة الحالية؛ اقراها مرة في البداية، ومرة بعد كل Rescue ناجحة، وضيفها كلها. ~Rescue(id)~ تقبل أعلى أولوية حالياً فقط، وبعد النجاح تفتح موجة جديدة؛ استعلام Arrivals نفسه مبيستهلكش الموجة.

## أخطاء شائعة
~~~go
q.Push(request) // أضاف من غير إصلاح heap
request := heap.Pop(q) // ممكن panic لو فاضي
~~~
لو رتّبت الدفعة الأولى مرة وبس، هتتجاهل أولوية الطلب الجديد. الـ ID الأصغر يكسر تعادل Priority، ومش ترتيب الوصول.

## المهمة
كمّل ~RescueHeap~: Priority الأكبر يسبق، ثم ID الأصغر. استخدم heap.Push وheap.Pop. ابدأ بالموجة الأولى؛ بعد كل Rescue ناجحة أدخل الموجة الجديدة قبل الاختيار التالي. كرّر لحد ما heap تفضى.

اخدم كل طلب مرة واحدة وابعت عدد اللي اتخدموا بتقرير واحد. طلب غلط أو مكرر يفشّل المهمة حتى لو كملت؛ مفيش تعديل أولوية أو إلغاء في العقد ده، ومفيش goroutines أو انتظار حقيقي.
`, Hints: []string{"في Less افحص Priority الأول، وID عند التعادل.", "Pop اللي أنت بتكتبها تشيل الآخر؛ heap.Pop هي اللي تختار الأفضل.", "بعد نجاح Rescue اقرأ Arrivals وضيفهم قبل اللفة التالية؛ متقراش الموجة مرتين."},
		En: &LevelText{Title: "Rescue Priority", Goal: "Serve highest priority as new distress waves arrive", Starter: rescuePriorityStarter, Lesson: `# A Queue That Reorders Itself

FIFO preserves arrival order. Rescue must select the highest priority currently waiting; new requests can overtake old ones. A heap retains the best element at its root without sorting everything each time.

## Before you write
- Which tie rule makes the result deterministic?
- Will you enqueue new arrivals before or after selecting the next request?

## The container/heap contract
~~~go
type Scores []int
func (s Scores) Len() int { return len(s) }
func (s Scores) Less(i,j int) bool { return s[i]<s[j] }
func (s Scores) Swap(i,j int) { s[i],s[j]=s[j],s[i] }
func (s *Scores) Push(v any) { *s=append(*s,v.(int)) }
func (s *Scores) Pop() any {
 old:=*s; n:=len(old); last:=old[n-1]; *s=old[:n-1]; return last
}
~~~
Import ~container/heap~. The interface requires Len, Less, Swap, Push, and Pop. Less defines who precedes whom; the example puts smaller numbers first rather than mission priorities. Push/Pop need pointer receivers to change slice length. ~v.(int)~ is the type assertion from signals; we only insert the agreed type.

## Call the package, not the method
~~~go
s:=&Scores{}
heap.Push(s,8)
heap.Push(s,3)
fmt.Println(heap.Pop(s).(int)) // 3
~~~
~heap.Push~ inserts and repairs heap order. ~heap.Pop~ moves the best element to the end before calling your Pop method, which therefore removes the last element. ~heap.Init~ prepares an existing slice; an empty heap needs no initialization. Check Len before Pop.

## Arriving waves
~~~go
arrivals:=colony.RescueArrivals()
for _,request:=range arrivals {
 fmt.Println(request.ID,request.Priority)
}
ok:=colony.Rescue(arrivals[0].ID) // Example only, not the selection rule
fmt.Println(ok)
~~~
~RescueRequest~ has unique ID and integer Priority fields. ~RescueArrivals~ returns a copy of the current wave. Read it once initially and once after each successful Rescue, enqueueing every request. ~Rescue(id)~ accepts only the current highest priority and releases another wave after success; reading Arrivals does not consume a wave.

## Common mistakes
~~~go
q.Push(request) // Inserts without repairing heap order
request := heap.Pop(q) // Panics if empty
~~~
Sorting only the initial batch ignores newly arriving priorities. Smaller ID breaks equal Priority; arrival order does not.

## Your mission
Complete ~RescueHeap~: larger Priority first, then smaller ID. Use heap.Push and heap.Pop. Enqueue the initial wave; after every successful Rescue enqueue the new wave before selecting again. Continue until the heap is empty.

Serve each request once and report the served count exactly once. An incorrect or duplicate request permanently fails the mission. Priority changes and cancellation are outside this contract; no goroutines or real waiting are needed.
`, Hints: []string{"In Less, compare Priority first and ID on ties.", "Your Pop method removes the last element; heap.Pop selects the best.", "After successful Rescue, enqueue Arrivals before the next iteration; do not read a wave twice."}}, Setup: func(w *sim.World) {
			w.Init(5, 3, sim.Regolith)
			ids := w.Rand.Perm(40)
			next := 0
			for wave := 0; wave < 6; wave++ {
				n := 1 + w.Rand.Intn(3)
				if wave == 0 {
					n = 4
				}
				requests := []sim.RescueRequest{}
				for j := 0; j < n; j++ {
					requests = append(requests, sim.RescueRequest{ID: ids[next] + 1, Priority: 1 + w.Rand.Intn(4)})
					next++
				}
				if wave == 0 {
					for i := range requests {
						requests[i].Priority = 2
					}
				}
				if wave == 1 {
					requests[0].Priority = 9
				}
				w.Data.RescueWaves = append(w.Data.RescueWaves, requests)
			}
			w.Data.RescuePending = append([]sim.RescueRequest(nil), w.Data.RescueWaves[0]...)
			want := fmt.Sprint(next)
			w.Check = func(w *sim.World) (bool, string) {
				if w.Data.RescueFault {
					return false, i18n.T("طلب مش معروف أو مش الأعلى أولوية؛ ابدأ من جديد", "Unknown or out-of-priority request; restart the mission")
				}
				if len(w.Data.Rescued) == next && len(w.Data.RescuePending) == 0 && len(w.Reports) == 1 && w.Reports[0] == want {
					return true, i18n.T("كل طلب اتخدم مرة وبالأولوية الصح", "Every request was served once in priority order")
				}
				return false, fmt.Sprintf(i18n.T("إنقاذ: %d / %d؛ ضيف الموجة الجديدة بعد كل خدمة", "Rescues: %d / %d; enqueue the new wave after each service"), len(w.Data.Rescued), next)
			}
		}})
}

func sortedInts(s []int) bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}

const deepProbeStarter = `package main

import "orbital/colony"

func depthIndex(depths []int, target int) int                                { return -1 }
func flood(grid [][]colony.Ground, x, y int, seen map[colony.Point]bool) int { return 0 }
func main() {
	index := depthIndex(colony.Depths(), colony.Reading(0))
	count := flood(colony.Grid(), 0, 0, map[colony.Point]bool{})
	colony.Report([]int{index, count})
}
`
const rescuePriorityStarter = `package main

import (
	"container/heap"
	"orbital/colony"
)

type RescueHeap []colony.RescueRequest

func (q RescueHeap) Len() int           { return len(q) }
func (q RescueHeap) Less(i, j int) bool { return q[i].ID < q[j].ID }
func (q RescueHeap) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
func (q *RescueHeap) Push(v any)        { *q = append(*q, v.(colony.RescueRequest)) }
func (q *RescueHeap) Pop() any {
	old := *q
	last := old[len(old)-1]
	*q = old[:len(old)-1]
	return last
}
func main() {
	q := &RescueHeap{}
	for _, request := range colony.RescueArrivals() {
		heap.Push(q, request)
	}
	served := 0
	for q.Len() > 0 {
		request := heap.Pop(q).(colony.RescueRequest)
		if !colony.Rescue(request.ID) {
			return
		}
		served++
		// Enqueue the newly released wave here.
	}
	colony.Report(served)
}
`
