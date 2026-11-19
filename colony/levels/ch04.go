package levels

import (
	"fmt"
	"slices"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch04RipeOnlyEnglish,
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
				"~if~ protects an action by running it only when ready is true.",
				"Keep the readiness decision separate from your route to the next tile.",
				"Use the starter's ~colony.CanHarvest()~ condition to guide its body.",
			},
		},
		LessonRevision: 1, Steps: ch04RipeOnlyArabic,
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
			"~if~ بتحمي فعل وبتشغّله بس لما الجاهزية تبقى true.",
			"افصل قرار الجاهزية عن طريقك للخانة اللي بعدها.",
			"استخدم شرط ~colony.CanHarvest()~ الموجود في الـstarter عشان توجّه جسمه.",
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
			LessonRevision: 1, Steps: ch04TerrainEnglish,
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
				"~switch~ sends one value to one matching case.",
				"Plan a ground choice for each tile, including unsupported terrain.",
				"Inspect the starter's ~colony.GetGround()~ value and its existing cases.",
			},
		},
		LessonRevision: 1, Steps: ch04TerrainArabic,
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
			"~switch~ بتبعت قيمة لفرع case واحد مطابق.",
			"خطط لاختيار حسب الأرض في كل خانة، بما فيها الأرض اللي مش مدعومة.",
			"راجع قيمة ~colony.GetGround()~ وحالات الـswitch الموجودة في الـstarter.",
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

	checkpointStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.Reading(0))
	// Build the report described in the lesson.
}
`
	add(&Level{LessonRevision: 1, Steps: ch04CheckpointArabic, ID: "4-3-checkpoint", Chapter: 4, Title: "بوابة الفحص", Goal: "صنّف قراءات صحيحة مع فحص الحدود والقيم خارج المدى",
		Requires: []string{"if"}, Bridge: []int{1, 7},
		Par: 0, Gold: 0, MaxTicks: 100,
		Starter: checkpointStarter, Lesson: `# بوابة الفحص 🚦

كل قراءة رقم صحيح وليها تصنيف نصّي. استخدم ~ReadingCount~ لمعرفة عدد القراءات و~Reading(i)~ لقراءة واحدة بالـindex.

امشِ على القراءات بترتيبها. افحص القيم الخارجة عن المجال قبل اختيار الفئة، وخلي كل حد تابع لفئة واحدة.

ابعت label واحدة لكل قراءة بـ~colony.Report~. راجع العقد والحدود من غير ما تفترض ترتيب ثابت للبيانات.`, Hints: []string{"الحد لازم يتبع فئة واحدة بس.", "افصل فحص المجال عن اختيار الفئة، وحافظ على ترتيب القراءات.", "استخدم ~colony.ReadingCount()~ لحد الحلقة و~colony.Reading(i)~ للقيمة الحالية."},
		En: &LevelText{LessonRevision: 1, Steps: ch04CheckpointEnglish, Title: "Inspection Checkpoint", Goal: "Classify integer readings, including boundaries and out-of-range values", Starter: checkpointStarter, Lesson: `# Inspection Checkpoint 🚦

Each reading is an integer with a text category. Use ~ReadingCount~ for the number of readings and ~Reading(i)~ for one indexed value.

Visit readings in order. Check values outside the allowed range before choosing a category, and give each boundary one owner.

Send one label per reading with ~colony.Report~. Follow the contract without assuming a fixed data order.`, Hints: []string{"A boundary should belong to exactly one category.", "Separate range validation from band selection and preserve reading order.", "Use ~colony.ReadingCount()~ for the loop bound and ~colony.Reading(i)~ for the current value."}},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			w.Data.SensorReadings = []int{-1, 0, 19, 20, 49, 50, 79, 80, 100, 101}
			for i := 0; i < 5; i++ {
				w.Data.SensorReadings = append(w.Data.SensorReadings, w.Rand.Intn(141)-20)
			}
			w.Rand.Shuffle(len(w.Data.SensorReadings), func(i, j int) {
				w.Data.SensorReadings[i], w.Data.SensorReadings[j] = w.Data.SensorReadings[j], w.Data.SensorReadings[i]
			})
			var want []string
			for _, n := range w.Data.SensorReadings {
				band := "invalid"
				switch {
				case n < 0 || n > 100:
				case n < 20:
					band = "low"
				case n < 50:
					band = "ready"
				case n < 80:
					band = "high"
				default:
					band = "alert"
				}
				want = append(want, band)
			}
			w.Check = func(w *sim.World) (bool, string) {
				if slices.Equal(w.Reports, want) {
					return true, i18n.T("كل التقارير مطابقة وبالترتيب", "All reports match in order")
				}
				return false, i18n.T("التقارير لسه ناقصة أو مش مطابقة؛ راجع القيم وترتيبها", "Reports are missing or do not match; check values and order")
			}
		},
	})

}
