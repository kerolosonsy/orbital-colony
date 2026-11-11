package levels

import (
	"fmt"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			Title: "Patience Pays",
			Goal:  "Collect 10 moss without destroying a single plant",
			Lesson: `# if: The Drone Makes Decisions 🤔

These moss plants were planted at different times, so they are not all ripe. Harvest even one unripe plant and **the mission fails**.

~~~go
if colony.CanHarvest() {
	colony.Harvest()
}
~~~

Code inside ~{ }~ runs **only** when the condition is ~true~.

## else and else if

~~~go
x, y := colony.Position()
if y == 0 && x < 4 {
	colony.Move(colony.East)
} else if y == 0 {
	colony.Move(colony.North)
} else {
	colony.Move(colony.West)
}
~~~

## Comparison and logical operators
| | Meaning |
|---|---|
| ~==~ | Equal (two = signs, not one!) |
| ~!=~ | Not equal |
| ~<~ ~<=~ ~>~ ~>=~ | Less than / greater than, with or without equality |
| ~&&~ | AND (both must be true) |
| ~\|\|~ | OR (either one is enough) |
| ~!~ | NOT: ~!colony.CanHarvest()~ |

## if with a short statement
Go lets you declare a variable inside the if statement itself:

~~~go
if x, _ := colony.Position(); x == 4 {
	// x exists only within this if
}
~~~`,
			Starter: `package main

import "orbital/colony"

func main() {
	for {
		if colony.CanHarvest() {
			// ...
		}
		// Move to the next tile
	}
}
`,
			Hints: []string{
				"~colony.CanHarvest()~ returns ~true~ when the tile is ready. Inside the if, Harvest then Plant to keep the tile productive.",
				"The field has two rows: head right along the first, go North, return left, then go South. Repeat.",
				"~colony.Position()~ returns x and y. Choose a direction from them: ~if y == 0 && x < 4 { ... } else if ...~",
			},
		},
		ID: "4-1-ripe-only", Chapter: 4, Title: "الصبر مفتاح",
		Goal: "اجمع 10 طحالب من غير ما تدمّر ولا نبتة",
		Par:  80, Gold: 55, MaxTicks: 800,
		Requires: []string{"if"},
		Setup: func(w *sim.World) {
			w.Init(5, 2, sim.Soil)
			w.FailOnDestroy = true
			w.Each(func(x, y int, _ *sim.Tile) { w.PlantAt(x, y, sim.Moss, w.Rand.Intn(30)) })
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Moss]
				return n >= 10, fmt.Sprintf(i18n.T("طحالب: %d / 10", "Moss: %d / 10"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	for {
		if colony.CanHarvest() {
			// ...
		}
		// اتحرّك للخانة اللي بعدها
	}
}
`,
		Hints: []string{
			"~colony.CanHarvest()~ بترجّع ~true~ لو الخانة جاهزة. جوه الـ if: Harvest وبعدين Plant عشان الخانة تفضل شغالة.",
			"الحقل صفين: امشي يمين في الصف الأول، اطلع North، ارجع شمال، وانزل South. وكرّر.",
			"~colony.Position()~ بترجّع x و y. ممكن تقرر الاتجاه على حسبهم: ~if y == 0 && x < 4 { ... } else if ...~",
		},
		Lesson: `# if: الدرون بيفكّر 🤔

الطحالب هنا مزروعة في أوقات مختلفة، فمش كلها ناضجة. ولو حصدت نبتة واحدة قبل ما تنضج، **المهمة تفشل**.

~~~go
if colony.CanHarvest() {
	colony.Harvest()
}
~~~

الكود اللي جوه ~{ }~ بيتنفّذ **بس** لو الشرط ~true~.

## else و else if

~~~go
x, y := colony.Position()
if y == 0 && x < 4 {
	colony.Move(colony.East)
} else if y == 0 {
	colony.Move(colony.North)
} else {
	colony.Move(colony.West)
}
~~~

## علامات المقارنة والمنطق
| | المعنى |
|---|---|
| ~==~ | يساوي (اتنين = مش واحدة!) |
| ~!=~ | لا يساوي |
| ~<~ ~<=~ ~>~ ~>=~ | أصغر / أكبر |
| ~&&~ | و (الاتنين لازم true) |
| ~\|\|~ | أو (واحد منهم كفاية) |
| ~!~ | عكس: ~!colony.CanHarvest()~ |

## if ومعاها سطر قصير
Go بيسمحلك تعرّف متغير جوه الـ if نفسها:

~~~go
if x, _ := colony.Position(); x == 4 {
	// x موجود بس جوه الـ if دي
}
~~~`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Mixed Terrain",
			Goal:  "Plant the right crop on every tile: moss on soil, lichen on ice",
			Lesson: `# switch: Many Choices 🔀

When one value has many possible cases, a chain of ~if/else if~ statements gets long and awkward. ~switch~ is clearer:

~~~go
switch colony.GetGround() {
case colony.Soil:
	colony.Plant(colony.Moss)
case colony.Ice:
	colony.Plant(colony.Lichen)
case colony.Regolith:
	colony.Till()
	colony.Plant(colony.Moss)
default:
	// Anything else (such as Rock)
}
~~~

## How Go differs from other languages
- No ~break~ needed: each ~case~ ends automatically without continuing into the next one.
- You can put several values in one case: ~case colony.Soil, colony.Regolith:~
- A ~switch~ without a value works like if/else if:

~~~go
switch {
case x < 5:
	// ...
case y > 2:
	// ...
}
~~~

## Your mission
Every tile has its own terrain. Moss grows on soil (and regolith after Till); Lichen grows on ice.`,
			Starter: `package main

import "orbital/colony"

func main() {
	switch colony.GetGround() {
	case colony.Soil:
		// ...
	case colony.Ice:
		// ...
	}
}
`,
			Hints: []string{
				"Nothing grows on Rock, so leave it alone. Regolith needs Till first, then Moss.",
				"Traverse the whole field in a snake pattern and switch on ~colony.GetGround()~ for each tile.",
				"Each GetGround costs a tick. Call it just once per tile.",
			},
		},
		ID: "4-2-terrain", Chapter: 4, Title: "تضاريس مختلطة",
		Goal: "ازرع كل خانة باللي يناسبها: طحالب على التربة، أشنة على الجليد",
		Par:  80, Gold: 60, MaxTicks: 400,
		Requires: []string{"switch"},
		Setup: func(w *sim.World) {
			w.Init(6, 3, sim.Soil)
			grounds := []sim.Ground{sim.Soil, sim.Ice, sim.Regolith, sim.Ice, sim.Regolith, sim.Rock}
			w.Each(func(x, y int, t *sim.Tile) { t.Ground = grounds[w.Rand.Intn(len(grounds))] })
			w.At(1, 0).Ground, w.At(2, 1).Ground, w.At(3, 2).Ground = sim.Ice, sim.Regolith, sim.Soil
			w.Check = func(w *sim.World) (bool, string) {
				done, total := 0, 0
				w.Each(func(_, _ int, t *sim.Tile) {
					switch t.Ground {
					case sim.Rock:
						return
					case sim.Ice:
						if t.Entity == sim.Lichen {
							done++
						}
					default:
						if t.Entity == sim.Moss {
							done++
						}
					}
					total++
				})
				return done == total, fmt.Sprintf(i18n.T("خانات مزروعة صح: %d / %d", "Correctly planted tiles: %d / %d"), done, total)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	switch colony.GetGround() {
	case colony.Soil:
		// ...
	case colony.Ice:
		// ...
	}
}
`,
		Hints: []string{
			"Rock مبيتزرعش، سيبه في حاله. الـ regolith محتاج Till الأول، وبعدها Moss.",
			"لف على الحقل كله بنفس طريقة الثعبان، ولكل خانة اعمل switch على ~colony.GetGround()~.",
			"كل GetGround بتاخد tick. نادِها مرة واحدة بس لكل خانة.",
		},
		Lesson: `# switch: اختيارات كتير 🔀

لما يكون عندك قيمة واحدة وحالات كتير، سلسلة ~if/else if~ بتبقى طويلة ووحشة. ~switch~ أنضف:

~~~go
switch colony.GetGround() {
case colony.Soil:
	colony.Plant(colony.Moss)
case colony.Ice:
	colony.Plant(colony.Lichen)
case colony.Regolith:
	colony.Till()
	colony.Plant(colony.Moss)
default:
	// أي حاجة تانية (زي Rock)
}
~~~

## حاجات تفرّق Go عن لغات تانية
- مفيش ~break~: كل ~case~ بيخلص لوحده ومش بيكمّل على اللي بعده.
- ممكن تحط أكتر من قيمة في نفس الحالة: ~case colony.Soil, colony.Regolith:~
- ~switch~ من غير قيمة بتشتغل زي if/else if:

~~~go
switch {
case x < 5:
	// ...
case y > 2:
	// ...
}
~~~

## المهمة
كل خانة لها أرض مختلفة. الطحالب بتتزرع على soil (والـ regolith بعد Till)، والأشنة (Lichen) على ice.`,
	})
}
