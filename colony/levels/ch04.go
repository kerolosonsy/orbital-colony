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
	add(&Level{ID: "4-3-checkpoint", Chapter: 4, Title: "بوابة الفحص", Goal: "صنّف قراءات صحيحة مع فحص الحدود والقيم خارج المدى",
		Requires: []string{"if"}, Bridge: []int{1, 7},
		Par: 0, Gold: 0, MaxTicks: 100,
		Starter: checkpointStarter, Lesson: `# الحدود ليها صاحب 🚦

بوابة الفحص بتستقبل أرقام صحيحة من حسّاسات الأكسجين. هنا مش محتاج تحوّل نصوص؛ شغلتك تقرر كل رقم ينتمي لأنهي فئة.

## قبل ما تكتب
- الرقم اللي على الحد بالظبط هينتمي للفئة اللي قبله ولا بعده؟
- هتفحص المدى المسموح قبل التصنيف ولا بعده؟

## قراءة رقم برقم
~~~go
count := colony.ReadingCount()
n := colony.Reading(0)
fmt.Println(count, n)
~~~
~ReadingCount~ بتقولك عدد القراءات. ~Reading(i)~ بتجيب قراءة واحدة كـ ~int~، ورقمها يبدأ من صفر. الأرقام جاهزة للحساب؛ مفيش parsing أو error.

استخدم حلقة العدّ اللي عرفناها في الفصل 3 عشان تزور الأرقام من 0 لحد قبل count. مثلاً حلقة عرض عدّاد بسيطة:
~~~go
for i := 0; i < 3; i++ {
 fmt.Println(i)
}
~~~

## رتب شروطك
~~~go
power := 7
label := "blocked"
if power < 2 || power > 8 {
 label = "blocked"
} else if power < 5 {
 label = "slow"
} else {
 label = "fast"
}
fmt.Println(label) // fast
~~~
ده مثال تصريح موتور: من 2 لـ8 مسموح، و5 بداية السرعة العالية. بنفحص القيم المرفوضة الأول، وبعدها الحدود من الأصغر للأكبر. أول شرط صحيح هو اللي بيتنفّذ، ومش بنكمل باقي السلسلة.

~label~ اتعرّفت قبل if عشان نقدر نستخدمها في التقرير بعد نهاية الشروط.

## أخطاء شائعة
~~~go
if power <= 5 {
 label = "slow"
}
~~~
لو 5 بداية fast، فـ ~<=~ حطّت القيمة دي في الفئة الغلط. كمان لو بدأت بفحص ~power < 5~ قبل المدى المسموح، رقم سالب ممكن يتصنّف slow بالغلط.

## المهمة
لف على كل القراءات وابعت تقرير نصّي منفصل لكل قراءة بنفس ترتيبها. المدى المسموح من 0 لـ100 شامل الطرفين؛ أي رقم أقل من 0 أو أكبر من 100 تقريره ~invalid~.

جوه المدى: أقل من20 هو ~low~، من20 لأقل من50 هو ~ready~، من50 لأقل من80 هو ~high~، ومن80 لحد100 هو ~alert~. ابعت الكلمة بس بـ ~colony.Report(label)~، من غير index أو أرقام إضافية. القيم 20 و50 و80 بتبدأ الفئة الجديدة بالظبط.`, Hints: []string{"اقرأ ~ReadingCount~ واستعمل حلقة عدّ، بدل ما تعمل قايمة قراءات.", "ابدأ برفض ~n < 0 || n > 100~ قبل شروط الفئات.", "خلي شروط الحدود تصاعدية وبـ ~<~، وابعت تقرير واحد بعد سلسلة if لكل قراءة."},
		En: &LevelText{Title: "Inspection Checkpoint", Goal: "Classify integer readings, including boundaries and out-of-range values", Starter: checkpointStarter, Lesson: `# Every Boundary Has an Owner 🚦

The checkpoint receives integer oxygen readings. There is no text conversion here; decide which band owns each number.

## Before you write
- Does an exact threshold belong to the preceding or following band?
- Will you validate the allowed range before or after classification?

## Read one number at a time
~~~go
count := colony.ReadingCount()
n := colony.Reading(0)
fmt.Println(count, n)
~~~
~ReadingCount~ gives the number of readings. ~Reading(i)~ returns one ~int~, with indices starting at zero. Values are ready for arithmetic; no parsing or errors are involved.

Use the counting loop from chapter 3 to visit indices from 0 up to, but not including, count. For example, this small display loop prints a counter:
~~~go
for i := 0; i < 3; i++ {
 fmt.Println(i)
}
~~~

## Order your conditions
~~~go
power := 7
label := "blocked"
if power < 2 || power > 8 {
 label = "blocked"
} else if power < 5 {
 label = "slow"
} else {
 label = "fast"
}
fmt.Println(label) // fast
~~~
This worked example grants motor access: 2 through 8 are allowed, and 5 starts the fast band. Reject invalid values first, then test thresholds from smallest to largest. Only the first matching branch runs.

~label~ is declared before the if so it can be used in a report after the branches.

## Common mistakes
~~~go
if power <= 5 {
 label = "slow"
}
~~~
If 5 starts fast, ~<=~ puts that exact value in the wrong band. Also, testing ~power < 5~ before validating the range can wrongly classify a negative number as slow.

## Your mission
Visit every reading and send one separate text report for each, in input order. The allowed range is 0 through 100 inclusive; anything below 0 or above 100 is ~invalid~.

Within that range: below 20 is ~low~, 20 to below 50 is ~ready~, 50 to below 80 is ~high~, and 80 through 100 is ~alert~. Send only the word with ~colony.Report(label)~, with no index or extra numbers. Exactly 20, 50, and 80 start their new bands.`, Hints: []string{"Read ~ReadingCount~ and use a counting loop rather than building a collection.", "Reject ~n < 0 || n > 100~ before testing bands.", "Test ascending thresholds with ~<~ and send one report after the if chain for each reading."}},
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
