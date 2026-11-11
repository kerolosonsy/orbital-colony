package levels

import (
	"fmt"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			Title: "First Signal",
			Goal:  "Guide the drone to the beacon 📡",
			Lesson: `# Welcome, Colony Commander 🚀

Your drone has landed on **Kepler-G**. It understands only **Go**: everything you write becomes commands, which the drone follows **in order**, from top to bottom.

## The structure of a Go program

~~~go
package main            // An executable program must be in package main

import "orbital/colony" // Import the drone control library

func main() {           // Execution starts here
	colony.Move(colony.East)
}
~~~

- ~package main~: every Go file starts with its package name. ~main~ means this is an executable program rather than a library.
- ~import~: brings in ready-made code from elsewhere. Here we import ~colony~, which provides the drone commands.
- ~func main() { ... }~: the function where Go begins execution. Everything between ~{~ and ~}~ is the function body.

## What is a function call?

~colony.Move(colony.East)~ means: call the ~Move~ function from the ~colony~ library and pass it the value ~colony.East~.

Directions: ~North~ is up, ~South~ is down, ~East~ is right, and ~West~ is left.

Each drone command uses one **tick** of energy. Finish in fewer ticks to earn more stars ⭐.

## Your mission
The beacon 📡 is 3 tiles to the right. Guide the drone there and press **Run** (or Ctrl/⌘ + Enter).`,
			Starter: `package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	// Continue here...
}
`,
			Hints: []string{
				"Each ~colony.Move(colony.East)~ line moves the drone just one tile.",
				"The beacon is 3 tiles to the right. You need 3 Move lines.",
			},
		},
		ID: "1-1-first-signal", Chapter: 1, Title: "إشارة أولى",
		Goal: "وصّل الدرون للمنارة 📡",
		Par:  3, Gold: 3, MaxTicks: 30,
		Setup: func(w *sim.World) {
			w.Init(6, 3, sim.Regolith)
			w.At(3, 0).Entity = sim.Beacon
			w.Check = func(w *sim.World) (bool, string) {
				return w.X == 3 && w.Y == 0, fmt.Sprintf(i18n.T("الدرون عند (%d, %d) — المنارة عند (3, 0)", "Drone at (%d, %d) — beacon at (3, 0)"), w.X, w.Y)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	// كمّل هنا...
}
`,
		Hints: []string{
			"كل سطر ~colony.Move(colony.East)~ بيحرّك الدرون خانة واحدة بس.",
			"المنارة على بعد 3 خانات يمين. يعني محتاج 3 أسطر Move.",
		},
		Lesson: `# أهلاً يا قائد المستعمرة 🚀

الدرون بتاعك نزل على كوكب **كيبلر-ج**. الدرون ده مبيفهمش غير لغة **Go**: كل اللي هتكتبه هيتحوّل لأوامر، والدرون هينفّذها **بالترتيب** من فوق لتحت.

## هيكل أي برنامج Go

~~~go
package main            // البرنامج اللي بيشتغل لازم يبقى في package main

import "orbital/colony" // بنستورد مكتبة التحكم في الدرون

func main() {           // التنفيذ بيبدأ من هنا
	colony.Move(colony.East)
}
~~~

- ~package main~: كل ملف Go بيبدأ باسم الـ package. ~main~ معناها "ده برنامج بيشتغل" مش مكتبة.
- ~import~: بتجيب كود جاهز من مكان تاني. هنا بنجيب ~colony~ اللي فيها أوامر الدرون.
- ~func main() { ... }~: الدالة اللي Go بيبدأ ينفّذ منها. أي حاجة بين ~{~ و ~}~ هي جسم الدالة.

## يعني إيه استدعاء دالة؟

~colony.Move(colony.East)~ معناها: من مكتبة ~colony~ نادي الدالة ~Move~ وابعتلها قيمة ~colony.East~.

الاتجاهات: ~North~ فوق، ~South~ تحت، ~East~ يمين، ~West~ شمال.

كل أمر للدرون بياكل **tick** واحد من الطاقة. كل ما تخلّص بـ ticks أقل، تاخد نجوم أكتر ⭐.

## المهمة
المنارة 📡 على بعد 3 خانات يمين الدرون. وصّله لها ودوس **Run** (أو Ctrl/⌘ + Enter).`,
	})

	add(&Level{
		En: &LevelText{
			Title: "First Harvest",
			Goal:  "Harvest the 3 moss plants 🌿",
			Lesson: `# Food First 🌿

The colony needs food. There are 3 ripe moss plants to the right of the drone.

## A function with no arguments

~~~go
colony.Harvest()
~~~

~Harvest~ needs no values, so we leave the parentheses empty: ~()~. **You still need the parentheses**: they tell Go to run the function now.

## Order is logic

The computer runs your lines one after another, exactly as written. Harvest before moving and you will harvest the wrong tile.

~~~go
colony.Move(colony.East) // Move onto the moss tile
colony.Harvest()         // Harvest it
~~~

## Tip
If something goes wrong, look at the **Console**. The drone tells you when a command fails and why.`,
			Starter: `package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	colony.Harvest()
}
`,
			Hints: []string{
				"Order matters: move onto the moss first, then call ~colony.Harvest()~.",
				"Move, then Harvest. Repeat 3 times.",
			},
		},
		ID: "1-2-first-harvest", Chapter: 1, Title: "أول حصاد",
		Goal: "احصد الـ 3 طحالب 🌿",
		Par:  6, Gold: 6, MaxTicks: 40,
		Setup: func(w *sim.World) {
			w.Init(5, 2, sim.Regolith)
			for x := 1; x <= 3; x++ {
				w.At(x, 0).Ground = sim.Soil
				w.PlantAt(x, 0, sim.Moss, 0)
			}
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Moss]
				return n >= 3, fmt.Sprintf(i18n.T("طحالب: %d / 3", "Moss: %d / 3"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	colony.Harvest()
}
`,
		Hints: []string{
			"الترتيب مهم: الأول تتحرك فوق الطحلب، وبعدين ~colony.Harvest()~.",
			"Move ثم Harvest، وكرّرهم 3 مرات.",
		},
		Lesson: `# الأكل أولاً 🌿

المستعمرة محتاجة أكل. فيه 3 طحالب ناضجة جنب الدرون على اليمين.

## دالة من غير مدخلات

~~~go
colony.Harvest()
~~~

~Harvest~ مش محتاجة أي قيمة، فبنكتب الأقواس فاضية ~()~. **الأقواس لازم تتكتب** حتى لو فاضية، لأنها هي اللي بتقول لـ Go "نفّذ الدالة دي دلوقتي".

## الترتيب = المنطق

الكمبيوتر بينفّذ الأسطر واحد ورا التاني بالظبط. لو حصدت قبل ما تتحرك، هتحصد الخانة الغلط.

~~~go
colony.Move(colony.East) // روح للخانة اللي عليها الطحلب
colony.Harvest()         // احصده
~~~

## نصيحة
لو حصلت مشكلة، بص على الـ **Console**. الدرون بيقولك فيه لو أمر فشل وليه.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "New Ground",
			Goal:  "Plant moss on the 3 tiles to the right",
			Lesson: `# Grow Your Own 🌱

We have used up the moss we harvested. Time to plant more.

## A function that takes an argument

~~~go
colony.Plant(colony.Moss)
~~~

The value inside the parentheses is an **argument**. The same function can plant different crops depending on the value you pass.

## Prepare the ground first
The planet's regolith is hard. Moss needs soil:

~~~go
colony.Till()             // Turn regolith into soil
colony.Plant(colony.Moss) // Planting will now succeed
~~~

If you plant before tilling, ~Plant~ fails and the Console explains why. Try it and see!

## Your mission
Till and plant moss at (1,0), (2,0), and (3,0).`,
			Starter: `package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	colony.Till()
	// Plant here...
}
`,
			Hints: []string{
				"Moss only grows on soil. This ground is regolith, so call ~colony.Till()~ first.",
				"For each tile: ~Move~, then ~Till~, then ~colony.Plant(colony.Moss)~.",
			},
		},
		ID: "1-3-planting", Chapter: 1, Title: "أرض جديدة",
		Goal: "ازرع طحالب في الـ 3 خانات اللي على اليمين",
		Par:  9, Gold: 9, MaxTicks: 50,
		Setup: func(w *sim.World) {
			w.Init(5, 2, sim.Regolith)
			w.At(0, 0).Ground = sim.Rock
			w.Check = func(w *sim.World) (bool, string) {
				n := 0
				for x := 1; x <= 3; x++ {
					if w.At(x, 0).Entity == sim.Moss {
						n++
					}
				}
				return n == 3, fmt.Sprintf(i18n.T("خانات مزروعة: %d / 3", "Planted tiles: %d / 3"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	colony.Move(colony.East)
	colony.Till()
	// ازرع هنا...
}
`,
		Hints: []string{
			"الطحالب مبتتزرعش غير على تربة (soil). الأرض هنا regolith، فلازم ~colony.Till()~ الأول.",
			"لكل خانة: ~Move~ ثم ~Till~ ثم ~colony.Plant(colony.Moss)~.",
		},
		Lesson: `# نزرع بنفسنا 🌱

الطحالب اللي حصدناها خلصت. محتاجين نزرع تاني.

## دالة بتاخد مُدخل (argument)

~~~go
colony.Plant(colony.Moss)
~~~

القيمة اللي جوه الأقواس اسمها **argument**. نفس الدالة ممكن تزرع حاجات مختلفة على حسب القيمة اللي بتبعتها.

## الأرض لازم تتجهّز
تربة الكوكب (regolith) قاسية. الطحالب محتاجة تربة (soil):

~~~go
colony.Till()             // حوّل regolith لـ soil
colony.Plant(colony.Moss) // دلوقتي الزرع هينجح
~~~

لو زرعت قبل الحرث، ~Plant~ هتفشل والـ Console هيقولك السبب. جرّبها عشان تشوف بنفسك!

## المهمة
احرث وازرع طحالب في الخانات (1,0) و (2,0) و (3,0).`,
	})
}
