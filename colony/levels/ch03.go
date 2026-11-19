package levels

import (
	"fmt"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch03LongRowEnglish,
			Title: "The Long Row",
			Goal:  "Harvest 14 moss plants — use for instead of writing 28 lines!",
			Lesson: `# The for Loop 🔁

Imagine writing the same line 14 times... then the farm grows to 1000 tiles. Programmers **do not repeat themselves**: they write code once and let the computer repeat it.

This example moves the drone north 5 times:

~~~go
for i := 0; i < 5; i++ {
	colony.Move(colony.North)
}
~~~

A for loop has 3 parts, separated by ~;~:
1. ~i := 0~: initialization. Create a counter named ~i~ starting at 0.
2. ~i < 5~: the condition. Run the body while this is ~true~.
3. ~i++~: after each iteration, increase ~i~ by one.

The body runs for ~i~ = 0, 1, 2, 3, 4. That is **5 times** exactly. Want a different number of repeats? Change the number in the condition.

## Why start at 0?
In programming, counting usually starts at 0. You will see this again with slices.

## Mission rule
Mission Control checks that your code really uses **for**. No shortcuts 😉`,
			Starter: `package main

import "orbital/colony"

func main() {
	for i := 0; i < 3; i++ {
		colony.Move(colony.East)
	}
}
`,
			Hints: []string{
				"A ~for~ body repeats while its condition stays true.",
				"Keep the repeated drone work inside the loop body.",
				"Check the starter's loop condition and commands inside its braces.",
			},
		},
		LessonRevision: 1, Steps: ch03LongRowArabic,
		ID: "3-1-long-row", Chapter: 3, Title: "الصف الطويل",
		Goal: "احصد 14 طحلب بـfor loop بدل ما تكتب 28 سطر!",
		Par:  28, Gold: 28, MaxTicks: 80,
		Requires: []string{"for"},
		Setup: func(w *sim.World) {
			w.Init(15, 2, sim.Regolith)
			for x := 1; x < 15; x++ {
				w.At(x, 0).Ground = sim.Soil
				w.PlantAt(x, 0, sim.Moss, 0)
			}
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Moss]
				return n >= 14, fmt.Sprintf(i18n.T("طحالب: %d / 14", "Moss: %d / 14"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	for i := 0; i < 3; i++ {
		colony.Move(colony.East)
	}
}
`,
		Hints: []string{
			"جسم الـfor loop بيتكرر طول ما الشرط true.",
			"خلّي شغل الدرون اللي محتاج تكراره جوه جسم الـloop.",
			"راجع شرط الـfor وأوامر الـstarter اللي بين الأقواس.",
		},
		Lesson: `# الحلقة for 🔁

تخيّل تكتب نفس السطر 14 مرة… وبعدين المزرعة تكبر لـ 1000 خانة. المبرمجين **مبيكرّروش نفسهم**: بيكتبوا الكود مرة واحدة، ويخلّوا الكمبيوتر يكرّره.

المثال ده بيحرّك الدرون لفوق (North) 5 مرات:

~~~go
for i := 0; i < 5; i++ {
	colony.Move(colony.North)
}
~~~

الـ for ليها 3 أجزاء مفصولين بـ ~;~:
1. ~i := 0~: بداية. بنعمل عدّاد اسمه ~i~ ويبدأ من 0.
2. ~i < 5~: الشرط. طول ما هو ~true~، الجسم يتنفّذ.
3. ~i++~: بعد كل لفّة، زوّد ~i~ واحد.

يعني الجسم هيتنفّذ لـ ~i~ = 0، 1، 2، 3، 4. ده **5 مرات** بالظبط. عاوز عدد تاني؟ غيّر الرقم اللي في الشرط.

## ليه بنبدأ من 0؟
في البرمجة العدّ تقريباً دايماً بيبدأ من 0. هتشوف ده تاني لما نوصل للـ slices.

## قاعدة المهمة
الـ Mission Control بيتأكد إن الكود بتاعك فيه **for** فعلاً. مفيش غش 😉`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch03FullFieldEnglish,
			Title: "The Whole Field",
			Goal:  "Harvest the entire 6×6 field (36 moss)",
			Lesson: `# A Loop Inside a Loop 🧩

This field has **rows** and **columns**. A natural solution is a loop over rows with a loop over tiles inside it.

~~~go
for y := 0; y < 6; y++ {
	for x := 0; x < 6; x++ {
		// This code runs 6 × 6 = 36 times
	}
}
~~~

The inner loop completes all 6 iterations **for every iteration** of the outer loop.

## Watch the walls
The map is 6 tiles wide, so the drone cannot move right from x = 5. If you try, ~Move~ returns ~false~ and wastes a tick. How many moves does each row actually need?

## Changing a variable
Use ~=~ to change an existing variable's value (rather than ~:=~):

~~~go
dir := colony.East
// ...
dir = colony.West // Same variable, new value
~~~

## Challenge ⭐⭐⭐
Can you finish in just **71 ticks**? (36 harvests + 35 moves)`,
			Starter: `package main

import "orbital/colony"

func main() {
	// A for inside a for:
	// The outer loop handles rows; the inner loop handles tiles in each row
	colony.Harvest()
}
`,
			Hints: []string{
				"Nested loops give separate counters to rows and items.",
				"Plan one row at a time, then account for the next row.",
				"Compare the two loop bounds and the direction choice in your row logic.",
			},
		},
		LessonRevision: 1, Steps: ch03FullFieldArabic,
		ID: "3-2-full-field", Chapter: 3, Title: "الحقل الكامل",
		Goal: "احصد الحقل كله 6×6 (36 طحلب)",
		Par:  100, Gold: 71, MaxTicks: 300,
		Requires: []string{"for2"},
		Setup: func(w *sim.World) {
			w.Init(6, 6, sim.Soil)
			w.Each(func(x, y int, _ *sim.Tile) { w.PlantAt(x, y, sim.Moss, 0) })
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Moss]
				return n >= 36, fmt.Sprintf(i18n.T("طحالب: %d / 36", "Moss: %d / 36"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	// for جوه for:
	// الحلقة الخارجية للصفوف، والداخلية للخانات في كل صف
	colony.Harvest()
}
`,
		Hints: []string{
			"الـloops المتداخلة بتدي عدّاد منفصل للصفوف والخانات.",
			"خطط صف واحد في كل مرة، وبعده فكّر في الصف اللي بعده.",
			"راجع حدود الحلقتين واختيار الاتجاه في منطق الصفوف.",
		},
		Lesson: `# حلقة جوه حلقة 🧩

الحقل ده **صفوف** و**أعمدة**. الحل الطبيعي: حلقة للصفوف، وجوّاها حلقة للخانات.

~~~go
for y := 0; y < 6; y++ {
	for x := 0; x < 6; x++ {
		// الكود هنا بيتنفّذ 6 × 6 = 36 مرة
	}
}
~~~

الحلقة الداخلية بتخلّص لفّاتها الـ 6 كلها **في كل لفّة** من الخارجية.

## خلّي بالك من الحيطان
الخريطة عرضها 6، فالدرون مقدرش يتحرك يمين من x = 5. لو حاولت، ~Move~ بترجّع ~false~ وبتضيّع tick على الفاضي. فكّر: محتاج كام حركة فعلاً في كل صف؟

## متغير بيتغيّر
تقدر تغيّر قيمة متغير موجود بـ ~=~ (مش ~:=~):

~~~go
dir := colony.East
// ...
dir = colony.West // نفس المتغير، قيمة جديدة
~~~

## تحدّي ⭐⭐⭐
تقدر تخلّص في **71 tick** بس؟ (36 حصاد + 35 حركة)`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch03EndlessEnglish,
			Title: "A Farm That Never Sleeps",
			Goal:  "Collect 30 moss from a 3×3 field",
			Lesson: `# for Without a Counter ♾️

You do not always know how many repetitions you need. Sometimes you want to repeat **until something happens**:

~~~go
for colony.Count(colony.Moss) < 30 {
	// Repeat while the moss count is below 30
}
~~~

This is the same for loop, with just a condition. Other languages call it ~while~, but Go uses **for** for every kind of loop.

## The infinite loop

~~~go
for {
	// Forever
}
~~~

Here that is fine: the mission ends automatically once you reach the goal. But take care: without any drone commands inside the loop, game time never advances and the program hangs (timeout).

## Growth takes time ⏳
Moss ripens **8 ticks** after planting. Harvest it earlier and **it is destroyed**. Plan a route that lets time pass while you work on the rest of the field.`,
			Starter: `package main

import "orbital/colony"

func main() {
	for colony.Count(colony.Moss) < 30 {
		// Tour the field: harvest ripe plants and replant each tile
	}
}
`,
			Hints: []string{
				"A condition-only ~for~ repeats until its condition becomes false.",
				"Make each pass add progress and advance the game state.",
				"Inspect the ~colony.Count~ condition and add a drone action to the body.",
			},
		},
		LessonRevision: 1, Steps: ch03EndlessArabic,
		ID: "3-3-endless", Chapter: 3, Title: "مزرعة لا تنام",
		Goal: "اجمع 30 طحلب من حقل 3×3",
		Par:  200, Gold: 140, MaxTicks: 1200,
		Requires: []string{"for"},
		Setup: func(w *sim.World) {
			w.Init(3, 3, sim.Soil)
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Moss]
				return n >= 30, fmt.Sprintf(i18n.T("طحالب: %d / 30", "Moss: %d / 30"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

func main() {
	for colony.Count(colony.Moss) < 30 {
		// لف على الحقل: احصد اللي نضج، وازرع مكانه
	}
}
`,
		Hints: []string{
			"for loop بشرط بس بتتكرر لحد ما الشرط يبقى false.",
			"خلّي كل دورة تزود التقدم وتحرك حالة اللعبة لقدّام.",
			"راجع شرط ~colony.Count~ وضيف أمر للدرون جوه جسم الـloop.",
		},
		Lesson: `# for من غير عدّاد ♾️

مش دايماً بتعرف هتكرر كام مرة. أحياناً عايز تكرر **لحد ما حاجة تحصل**:

~~~go
for colony.Count(colony.Moss) < 30 {
	// اتكرر طول ما الطحالب أقل من 30
}
~~~

ده نفس الـ for، بس بالشرط بس. لغات تانية بتسميه ~while~، لكن Go معندهاش غير **for** لكل أنواع التكرار.

## الحلقة اللانهائية

~~~go
for {
	// للأبد
}
~~~

في المستعمرة دي ده تمام، لأن المهمة بتقفل لوحدها أول ما الهدف يتحقق. بس خد بالك: لو جوه الحلقة مفيش ولا أمر للدرون، الوقت مش هيتحرك أبداً والبرنامج هيعلّق (timeout).

## النمو بياخد وقت ⏳
الطحلب بينضج بعد **8 ticks** من زرعه. لو حصدته قبل كده، **هيتدمّر**. فكّر في ترتيب يخلّي الوقت يعدّي وإنت بتشتغل في باقي الحقل.`,
	})
}
