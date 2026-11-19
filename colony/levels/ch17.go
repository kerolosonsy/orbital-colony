package levels

import (
	"fmt"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

type p3 struct{ x, y, z int }

// carve digs a random tunnel walk of n steps on layer z starting at (x, y)
// and returns every tile it opened, in order.
func carve(w *sim.World, x, y, z, n int) []p3 {
	path := []p3{{x, y, z}}
	w.AtZ(x, y, z).Ground = sim.Tunnel
	for len(path) < n {
		d := dirs[w.Rand.Intn(4)]
		nx, ny := x+d[0], y+d[1]
		if nx < 0 || ny < 0 || nx >= w.W || ny >= w.H {
			continue
		}
		x, y = nx, ny
		t := w.AtZ(x, y, z)
		if t.Ground != sim.Tunnel {
			t.Ground, t.Entity = sim.Tunnel, sim.Empty
			path = append(path, p3{x, y, z})
		}
	}
	return path
}

// shortest3D is the fewest actions from start until standing next to goal
// (then one Dig), moving through tunnels and shafts on every layer.
func shortest3D(w *sim.World, start, goal p3) int {
	open := func(p p3) bool {
		if p.x < 0 || p.y < 0 || p.x >= w.W || p.y >= w.H || p.z < 0 || p.z >= w.Depth {
			return false
		}
		g := w.AtZ(p.x, p.y, p.z).Ground
		return p.z == 0 || g == sim.Tunnel || g == sim.Shaft
	}
	dist := map[p3]int{start: 0}
	q := []p3{start}
	for len(q) > 0 {
		c := q[0]
		q = q[1:]
		for _, d := range dirs {
			if n := (p3{c.x + d[0], c.y + d[1], c.z}); n == goal {
				return dist[c] + 1
			}
		}
		var next []p3
		for _, d := range dirs {
			next = append(next, p3{c.x + d[0], c.y + d[1], c.z})
		}
		if w.AtZ(c.x, c.y, c.z).Ground == sim.Shaft {
			next = append(next, p3{c.x, c.y, c.z + 1})
		}
		if c.z > 0 && w.AtZ(c.x, c.y, c.z-1).Ground == sim.Shaft {
			next = append(next, p3{c.x, c.y, c.z - 1})
		}
		for _, n := range next {
			if _, seen := dist[n]; !seen && open(n) {
				dist[n] = dist[c] + 1
				q = append(q, n)
			}
		}
	}
	return 999
}

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch17Mission1(true),
			Title: "Your First Shaft",
			Goal:  "Dig down, descend, collect 3 iron from the rock, and return to the surface",
			Lesson: `# Layers: The World Has Depth ⛏️

So far everything has happened on the surface. But beneath Kepler-G lies rock containing **minerals**. Each tile now has **3 coordinates**: ~x~, ~y~, and ~z~ (depth):

- ~z = 0~: the surface.
- ~z = 1~: the layer below it, then ~z = 2~, and so on.

## New commands

| Command | What it does |
|---|---|
| ~colony.DigDown()~ | Digs a shaft from your tile to the layer below |
| ~colony.Descend()~ | Climbs down the shaft (z + 1) |
| ~colony.Ascend()~ | Climbs up if there is a shaft above (z − 1) |
| ~colony.Dig(dir)~ | Excavates the adjacent tile on your layer and collects any ore |
| ~colony.Layer()~ | Your current layer (free) |

## Different rules underground
- Solid stone (~colony.Stone~) **blocks** movement. Dig first.
- An excavated tile becomes a ~colony.Tunnel~ that you can walk through.
- A shaft (~colony.Shaft~) is the only connection between layers.

## Your mission
There are 3 iron veins on layer 1 east of the starting point. Dig to them and return to the surface.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	colony.DigDown() // This shaft connects the surface to the layer below
	colony.Descend() // Climb down
	fmt.Println("I am now on layer", colony.Layer())

	// Iron is buried in the rock east of you. Dig(East) excavates the tile and collects its ore
}
`,
			Hints: []string{
				"Depth adds a z coordinate; ~colony.Layer()~ reads the drone's current layer.",
				"Treat opening a shaft, entering it, collecting ore, and returning as separate phases.",
				"Compare ~colony.Dig~, ~colony.Move~, and ~colony.Ascend~ with the starter's initial shaft.",
			},
		},
		ID: "17-1-first-shaft", Chapter: 17, LessonRevision: 1, Steps: ch17Mission1(false), Title: "أول بير",
		Goal: "احفر لتحت، انزل، اجمع 3 حديد من جوه الصخر، واطلع تاني للسطح",
		Par:  18, Gold: 12, MaxTicks: 120,
		Setup: func(w *sim.World) {
			w.Init(5, 3, sim.Regolith)
			w.Underground(2, 0, 0)
			for x := 1; x <= 3; x++ {
				w.AtZ(x, 0, 1).Entity = sim.Iron
			}
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Iron]
				switch {
				case n < 3:
					return false, fmt.Sprintf(i18n.T("حديد: %d / 3", "Iron: %d / 3"), n)
				case w.Z != 0:
					return false, i18n.T("معاك الحديد! اطلع للسطح بـ Ascend", "You have the iron! Return to the surface with Ascend")
				}
				return true, i18n.T("رجعت بالحديد", "Returned with the iron")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	colony.DigDown() // البير ده بيوصل السطح بالطبقة اللي تحته
	colony.Descend() // انزل فيه
	fmt.Println("أنا دلوقتي في الطبقة", colony.Layer())

	// الحديد مدفون في الصخر شرقك. Dig(East) بيحفر الخانة وبيطلّع اللي فيها
}
`,
		Hints: []string{
			"العمق بيضيف إحداثية z؛ ~colony.Layer()~ بتقرأ الطبقة الحالية للدرون.",
			"قسّم الشغل لمراحل: فتح البير، النزول، جمع الحديد، والرجوع.",
			"قارن ~colony.Dig~ و~colony.Move~ و~colony.Ascend~ بالبير اللي الـstarter فتحه.",
		},
		Lesson: `# الطبقات: العالم ليه عمق ⛏️

لحد دلوقتي كل حاجة كانت على السطح. بس كيبلر-ج تحت أرضه صخر فيه **معادن**. كل خانة دلوقتي ليها **3 إحداثيات**: ~x~ و~y~ و~z~ (العمق):

- ~z = 0~: السطح.
- ~z = 1~: الطبقة اللي تحته، و~z = 2~ اللي تحتها، وهكذا.

## الأوامر الجديدة

| الأمر | بيعمل إيه |
|---|---|
| ~colony.DigDown()~ | يحفر بير في الخانة اللي إنت عليها لحد الطبقة اللي تحت |
| ~colony.Descend()~ | ينزل في البير (z + 1) |
| ~colony.Ascend()~ | يطلع لو فوقك بير (z − 1) |
| ~colony.Dig(dir)~ | يحفر الخانة اللي جنبك في نفس الطبقة، ويجمع أي معدن فيها |
| ~colony.Layer()~ | الطبقة اللي إنت فيها (ببلاش) |

## تحت الأرض قواعد تانية
- الصخر الصلب (~colony.Stone~) **بيقفل** الحركة. لازم تحفر الأول.
- الخانة المحفورة بتبقى نفق (~colony.Tunnel~) وتقدر تمشي فيها.
- البير (~colony.Shaft~) هو الطريق الوحيد بين الطبقات.

## المهمة
3 عروق حديد في الطبقة 1 شرق نقطة البداية. احفر لهم، وارجع للسطح.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch17Mission2(true),
			Title: "Map the Depths",
			Goal:  "Count open tunnels on each underground layer and report the counts in one slice, in order",
			Lesson: `# 3D Arrays: A Slice of Slices of Slices 🧊

The surface map was ~[][]colony.Ground~: rows and columns. Layers add a third dimension:

~~~go
world := make([][][]colony.Ground, colony.Depth())
for z := range world {
	world[z] = colony.LayerGrid(z) // world[z][x][y]
}
~~~

- ~world[z]~ is a whole layer (~[][]colony.Ground~).
- ~world[z][x]~ is a column in that layer.
- ~world[z][x][y]~ is one tile.

## Order matters
Choose an order for dimensions and document it in a comment. Here ~LayerGrid~ returns ~[x][y]~. Accidentally using ~[y][x]~ reads the wrong tile without a compiler warning, a very common bug.

## Three nested loops

~~~go
for z := 1; z < colony.Depth(); z++ {
	for x := range world[z] {
		for y := range world[z][x] {
			// ...
		}
	}
}
~~~

## Fixed-size arrays
If the size is known when writing the code, you can use an array instead of a slice:

~~~go
var counts [3]int // 3 zeros; the length is part of the type
counts[0]++
~~~`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println("Number of layers:", colony.Depth())
	layer1 := colony.LayerGrid(1) // layer1[x][y]
	fmt.Println(len(layer1), "×", len(layer1[0]))

	// For each layer from 1 to 3, count tiles that are Tunnel or Shaft
	// Report the result as []int with colony.Report
}
`,
			Hints: []string{
				"A 3D slice uses one index for each dimension: layer, x, then y.",
				"Inspect one layer's shape before writing loops that visit its columns and rows.",
				"Use the ~grid[x][y]~ comment beside ~LayerGrid~ in the starter to check your indexes.",
			},
		},
		ID: "17-2-layer-census", Chapter: 17, LessonRevision: 1, Steps: ch17Mission2(false), Title: "خريطة الأعماق",
		Goal: "عدّ الأنفاق المفتوحة في كل طبقة تحت الأرض، وابعتهم في slice واحدة بالترتيب",
		Par:  0, Gold: 0, MaxTicks: 60,
		Requires: []string{"slice"},
		Setup: func(w *sim.World) {
			w.Init(7, 5, sim.Regolith)
			w.Underground(4, 0, 0)
			w.At(0, 0).Ground = sim.Shaft
			counts := make([]int, 3)
			for z := 1; z <= 3; z++ {
				carve(w, w.Rand.Intn(7), w.Rand.Intn(5), z, 6+w.Rand.Intn(14))
				w.AtZ(0, 0, z).Ground = sim.Tunnel
				if z < 3 {
					w.AtZ(0, 0, z).Ground = sim.Shaft
				}
				w.Each(func(x, y int, _ *sim.Tile) {
					if g := w.AtZ(x, y, z).Ground; g == sim.Tunnel || g == sim.Shaft {
						counts[z-1]++
					}
				})
			}
			want := fmt.Sprint(counts)
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("خريطة الأعماق صح", "The depth map is correct")
				}
				if n := len(w.Reports); n > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[n-1] + i18n.T(" — العدد مش مظبوط", " — the count is incorrect")
				}
				return false, i18n.T("القاعدة مستنية slice فيها عدد الأنفاق في الطبقات 1 و2 و3", "Base is waiting for a slice counting tunnels on layers 1, 2, and 3")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println("عدد الطبقات:", colony.Depth())
	layer1 := colony.LayerGrid(1) // layer1[x][y]
	fmt.Println(len(layer1), "×", len(layer1[0]))

	// لكل طبقة من 1 لـ 3: عدّ الخانات اللي هي Tunnel أو Shaft
	// وابعت النتيجة كـ []int بـ colony.Report
}
`,
		Hints: []string{
			"الـ3D slice ليها index لكل بُعد: الطبقة، وبعدها x، وبعدها y.",
			"بص على شكل layer واحدة قبل ما تكتب loops تمرّ على أعمدتها وصفوفها.",
			"استعمل تعليق ~grid[x][y]~ جنب ~LayerGrid~ في الـstarter عشان تراجع الـindexes.",
		},
		Lesson: `# مصفوفات 3D: slice جوه slice جوه slice 🧊

الخريطة على السطح كانت ~[][]colony.Ground~: صفوف وأعمدة. دلوقتي عندنا طبقات، يعني بُعد تالت:

~~~go
world := make([][][]colony.Ground, colony.Depth())
for z := range world {
	world[z] = colony.LayerGrid(z) // world[z][x][y]
}
~~~

- ~world[z]~ هي طبقة كاملة (~[][]colony.Ground~).
- ~world[z][x]~ عمود في الطبقة دي.
- ~world[z][x][y]~ خانة واحدة.

## الترتيب مهم
اتفق مع نفسك على ترتيب الأبعاد واكتبه في comment. في المهمة دي ~LayerGrid~ بترجّع ~[x][y]~، فلو كتبت ~[y][x]~ بالغلط، هتقرا خانة غلط من غير ما Go يشتكي، وده من أشهر الـ bugs.

## حلقة متداخلة في 3 مستويات

~~~go
for z := 1; z < colony.Depth(); z++ {
	for x := range world[z] {
		for y := range world[z][x] {
			// ...
		}
	}
}
~~~

## arrays ثابتة الحجم
لو الحجم معروف وقت الكتابة، تقدر تستخدم array بدل slice:

~~~go
var counts [3]int // 3 أصفار، الحجم جزء من النوع
counts[0]++
~~~`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch17Mission3(true),
			Title: "The Deep Geode",
			Goal:  "A geode lies on layer 2. Reach it in the fewest steps through the existing tunnels and shafts",
			Lesson: `# BFS in Three Dimensions 🕳️

The BFS from chapter 15 works on any **graph**: nodes with connections between them. The number of dimensions does not matter. The only difference here is that each node has more neighbors:

~~~go
type pos struct{ x, y, z int }

func neighbors(world [][][]colony.Ground, p pos) []pos {
	var out []pos
	for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		n := pos{p.x + d[0], p.y + d[1], p.z}
		if walkable(world, n) {
			out = append(out, n)
		}
	}
	if world[p.z][p.x][p.y] == colony.Shaft { // Down
		out = append(out, pos{p.x, p.y, p.z + 1})
	}
	if p.z > 0 && world[p.z-1][p.x][p.y] == colony.Shaft { // Up
		out = append(out, pos{p.x, p.y, p.z - 1})
	}
	return out
}
~~~

## Structs as map keys
~pos~ is a struct of 3 numbers. Go accepts it **as a map key** because every field supports ~==~. Both ~map[pos]bool~ and ~map[pos]pos~ work directly.

## Why BFS instead of going straight down?
The shaft leading to layer 2 is not beside the geode. Digging through stone costs Dig + Move per tile. Existing tunnels are much cheaper, but you need an **algorithm** to find their shortest route. ⭐⭐⭐ = exactly the shortest route.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// pos is a tile in the three-dimensional world
type pos struct{ x, y, z int }

func main() {
	goal := colony.Targets()[0] // The geode's location on layer 2
	world := make([][][]colony.Ground, colony.Depth())
	for z := range world {
		world[z] = colony.LayerGrid(z)
	}
	fmt.Println("Geode at", goal, "layers:", len(world))

	// BFS from pos{0, 0, 0}: neighbors in 4 directions on the same layer,
	// + down if this tile is Shaft, + up if the tile above is Shaft.
	// Stop beside the geode and excavate it with Dig.
}
`,
			Hints: []string{
				"BFS still finds shortest equal-cost paths when a node gains a depth coordinate.",
				"Separate map state, allowed transitions, route reconstruction, and actions in your plan.",
				"Start with the ~pos~ definition and the shaft-neighbor comments in ~main.go~.",
			},
		},
		ID: "17-3-deep-geode", Chapter: 17, LessonRevision: 1, Steps: ch17Mission3(false), Title: "الجيود العميق",
		Goal: "فيه جيود في الطبقة 2. وصله بأقل عدد خطوات عن طريق الأنفاق والآبار الموجودة",
		Par:  60, Gold: 30, MaxTicks: 400,
		Setup: func(w *sim.World) {
			for {
				w.Init(8, 6, sim.Regolith)
				w.Underground(3, 0, 0)
				w.At(0, 0).Ground = sim.Shaft
				path1 := carve(w, 0, 0, 1, 22)
				for range 2 { // dead-end branches
					b := path1[w.Rand.Intn(len(path1))]
					carve(w, b.x, b.y, 1, 5)
				}
				mid := path1[len(path1)/2+w.Rand.Intn(len(path1)/2)]
				w.AtZ(mid.x, mid.y, 1).Ground = sim.Shaft
				path2 := carve(w, mid.x, mid.y, 2, 18)
				for range 2 {
					b := path2[w.Rand.Intn(len(path2))]
					carve(w, b.x, b.y, 2, 5)
				}
				end := path2[len(path2)-1]
				var goal p3
				found := false
				for _, d := range dirs { // the geode sits in the rock beside the tunnel's end
					g := p3{end.x + d[0], end.y + d[1], 2}
					if g.x >= 0 && g.y >= 0 && g.x < w.W && g.y < w.H && w.AtZ(g.x, g.y, 2).Ground == sim.Stone {
						goal, found = g, true
						break
					}
				}
				if !found {
					continue
				}
				w.AtZ(goal.x, goal.y, 2).Entity = sim.Geode
				w.Targets = []sim.Point{{X: goal.x, Y: goal.y}}
				best := shortest3D(w, p3{0, 0, 0}, goal)
				if best >= 999 || best < 10 {
					continue
				}
				w.Gold = best // the last step counted is the Dig into the geode
				w.Par = w.Gold + 14
				break
			}
			w.Check = func(w *sim.World) (bool, string) {
				if w.Inv[sim.Geode] > 0 {
					return true, i18n.T("الجيود في إيدك", "You have the geode")
				}
				return false, i18n.T("الجيود لسه في الطبقة 2", "The geode is still on layer 2")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

// pos خانة في العالم الثلاثي الأبعاد
type pos struct{ x, y, z int }

func main() {
	goal := colony.Targets()[0] // مكان الجيود في الطبقة 2
	world := make([][][]colony.Ground, colony.Depth())
	for z := range world {
		world[z] = colony.LayerGrid(z)
	}
	fmt.Println("الجيود عند", goal, "الطبقات:", len(world))

	// BFS من pos{0, 0, 0}: الجيران 4 اتجاهات في نفس الطبقة،
	// + لتحت لو الخانة Shaft، + لفوق لو اللي فوقك Shaft.
	// ووقف أول ما توصل لخانة جنب الجيود، واحفره بـ Dig.
}
`,
		Hints: []string{
			"BFS لسه بيلاقي أقصر طريق بتكلفة واحدة، حتى لو ضفنا العمق لإحداثية النقطة.",
			"افصل حالة الـmap، والانتقالات المسموحة، واسترجاع الطريق، والحركة في خطتك.",
			"ابدأ من تعريف ~pos~ وتعليقات جيران الـshaft في ~main.go~.",
		},
		Lesson: `# BFS في 3 أبعاد 🕳️

الـ BFS اللي كتبته في الفصل 15 بيشتغل على أي **جراف**: نقط وروابط بينها. الأبعاد مش فارقة معاه. الفرق الوحيد هنا إن كل نقطة ليها جيران أكتر:

~~~go
type pos struct{ x, y, z int }

func neighbors(world [][][]colony.Ground, p pos) []pos {
	var out []pos
	for _, d := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		n := pos{p.x + d[0], p.y + d[1], p.z}
		if walkable(world, n) {
			out = append(out, n)
		}
	}
	if world[p.z][p.x][p.y] == colony.Shaft { // لتحت
		out = append(out, pos{p.x, p.y, p.z + 1})
	}
	if p.z > 0 && world[p.z-1][p.x][p.y] == colony.Shaft { // لفوق
		out = append(out, pos{p.x, p.y, p.z - 1})
	}
	return out
}
~~~

## structs كمفاتيح
~pos~ struct فيه 3 أرقام، وGo بيقبله **كمفتاح في map** لأن كل حقوله بتتقارن بـ ~==~. يعني ~map[pos]bool~ و~map[pos]pos~ شغالين على طول.

## ليه BFS ومش "انزل على طول"؟
البير اللي بيوصلك للطبقة 2 مش جنب الجيود. ولو حفرت طريقك في الصخر، كل خانة بتكلّفك Dig + Move. الأنفاق الموجودة أرخص بكتير، بس محتاج **خوارزمية** تلاقي أقصر طريق فيها. ⭐⭐⭐ = أقصر طريق بالظبط.`,
	})
}
