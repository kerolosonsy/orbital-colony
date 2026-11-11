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
}

func sortedInts(s []int) bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}
