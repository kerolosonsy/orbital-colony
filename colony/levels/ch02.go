package levels

import (
	"fmt"
	"strconv"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			Title: "Inventory Check",
			Goal:  "Report the total moss + crystal in storage to base",
			Lesson: `# Variables 📦

The inventory numbers change each time you run the program, so you cannot type the answer yourself. Your program must **read** each number and store it in a **variable**.

## Declaring a variable

~~~go
moss := colony.Count(colony.Moss)
~~~

- ~:=~ means "create a variable and give it this value." Go infers its type (here, ~int~, a whole number).
- You can then use ~moss~ below that line.

## Arithmetic
~~~go
total := moss + crystals
half := total / 2 // Integer division: 7 / 2 = 3
~~~

## Print or report?
- ~fmt.Println(...)~ prints to the Console **for you**. It is useful for inspecting values while debugging.
- ~colony.Report(...)~ sends a value **to base**. This is what the mission checks.

## An important Go rule ⚠️
If you declare a variable and never use it, **the program will not build** (declared and not used). Go enforces this to keep code tidy.

## Your mission
Report the total moss and crystal in storage to base.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	moss := colony.Count(colony.Moss)
	fmt.Println("Moss in storage:", moss)

	// Find the crystal count, add the two counts, and send the total with colony.Report
}
`,
			Hints: []string{
				"Just as you created ~moss~, create ~crystals := colony.Count(colony.Crystal)~.",
				"Use ~total := moss + crystals~, then ~colony.Report(total)~.",
			},
		},
		ID: "2-1-inventory", Chapter: 2, Title: "جرد المخزن",
		Goal: "ابعت للقاعدة مجموع الطحالب + الكريستال اللي في المخزن",
		Par:  0, Gold: 0, MaxTicks: 20,
		Setup: func(w *sim.World) {
			w.Init(4, 3, sim.Regolith)
			w.Inv[sim.Moss] = 10 + w.Rand.Intn(60)
			w.Inv[sim.Crystal] = 5 + w.Rand.Intn(40)
			want := strconv.Itoa(w.Inv[sim.Moss] + w.Inv[sim.Crystal])
			w.Check = func(w *sim.World) (bool, string) {
				if reported(w, want) {
					return true, i18n.T("القاعدة استلمت الرقم الصح", "Base received the correct number")
				}
				if len(w.Reports) > 0 {
					return false, i18n.T("القاعدة استلمت ", "Base received ") + w.Reports[len(w.Reports)-1] + i18n.T(" — ده مش المجموع الصح", " — that is not the correct total")
				}
				return false, i18n.T("القاعدة مستنية التقرير", "Base is waiting for your report")
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	moss := colony.Count(colony.Moss)
	fmt.Println("طحالب في المخزن:", moss)

	// اعرف عدد الكريستال، واجمعهم، وابعت المجموع بـ colony.Report
}
`,
		Hints: []string{
			"زي ما عملنا ~moss~، اعمل متغير ~crystals := colony.Count(colony.Crystal)~.",
			"~total := moss + crystals~ وبعدين ~colony.Report(total)~.",
		},
		Lesson: `# المتغيرات 📦

الأرقام في المخزن بتتغير كل مرة تشغّل البرنامج، فمينفعش تكتب الرقم بإيدك. لازم البرنامج **يقرأه** ويحتفظ بيه في **متغير**.

## تعريف متغير

~~~go
moss := colony.Count(colony.Moss)
~~~

- ~:=~ معناها "اعمل متغير جديد واديله القيمة دي". Go بيعرف نوعه لوحده (هنا ~int~، يعني رقم صحيح).
- بعد كده تقدر تستخدم ~moss~ في أي مكان تحت السطر ده.

## الحساب
~~~go
total := moss + crystals
half := total / 2 // قسمة أرقام صحيحة: 7 / 2 = 3
~~~

## نطبع ولا نبعت؟
- ~fmt.Println(...)~ بيطبع في الـ Console **ليك إنت**. مفيد جداً تشوف القيم وإنت بتجرّب (debugging).
- ~colony.Report(...)~ بيبعت القيمة **للقاعدة**. ده اللي بيتحسب في المهمة.

## قاعدة مهمة في Go ⚠️
لو عرّفت متغير ومستخدمتوش، **البرنامج مش هيتبني** (declared and not used). Go صارم في كده عشان الكود يفضل نضيف.

## المهمة
ابعت للقاعدة مجموع الطحالب والكريستال اللي في المخزن.`,
	})

	add(&Level{
		En: &LevelText{
			Title: "Survey the Area",
			Goal:  "Report the map's area (width × height), then its perimeter",
			Lesson: `# Returning More Than One Value 📐

One useful Go feature: a function can **return several values** at once.

~~~go
width, height := colony.WorldSize()
~~~

~WorldSize~ returns two values, so we create two variables at once. If you do not need one of them, use ~_~:

~~~go
width, _ := colony.WorldSize() // We do not need the height
~~~

## Types
Every Go variable has a fixed type:

~~~go
var count int = 5         // A whole number
var name string = "Kepler" // Text
var ratio float64 = 0.75   // A floating-point number
ready := true              // bool: true or false
~~~

~var~ is the full declaration form; ~:=~ is the shorthand used inside functions. You cannot add an ~int~ to a ~string~: Go rejects the program, preventing many mistakes.

## Your mission
The map changes size each time. Report the **area** first, then the **perimeter**.`,
			Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	width, height := colony.WorldSize()
	fmt.Println("Width:", width, "Height:", height)

	// Area = width × height
	// Perimeter = 2 × (width + height)
}
`,
			Hints: []string{
				"Go uses ~*~ for multiplication: ~area := width * height~.",
				"Parentheses are evaluated first: ~2 * (width + height)~.",
				"Send the area with ~colony.Report(area)~ first, then the perimeter.",
			},
		},
		ID: "2-2-survey", Chapter: 2, Title: "مسح المنطقة",
		Goal: "ابعت مساحة الخريطة (العرض × الطول)، وبعدها محيطها",
		Par:  0, Gold: 0, MaxTicks: 20,
		Setup: func(w *sim.World) {
			wd, ht := 3+w.Rand.Intn(6), 2+w.Rand.Intn(5)
			w.Init(wd, ht, sim.Regolith)
			w.Each(func(x, y int, t *sim.Tile) {
				if w.Rand.Intn(5) == 0 {
					t.Ground = sim.Rock
				}
			})
			area, per := strconv.Itoa(wd*ht), strconv.Itoa(2*(wd+ht))
			w.Check = func(w *sim.World) (bool, string) {
				a := -1
				for i, r := range w.Reports {
					if r == area {
						a = i
						break
					}
				}
				if a == -1 {
					return false, i18n.T("القاعدة مستنية المساحة", "Base is waiting for the area")
				}
				for _, r := range w.Reports[a+1:] {
					if r == per {
						return true, i18n.T("المساحة والمحيط اتسجلوا", "Area and perimeter recorded")
					}
				}
				return false, fmt.Sprintf(i18n.T("المساحة وصلت (%s). مستنيين المحيط", "Area received (%s). Waiting for the perimeter"), area)
			}
		},
		Starter: `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	width, height := colony.WorldSize()
	fmt.Println("العرض:", width, "الطول:", height)

	// المساحة = العرض × الطول
	// المحيط = 2 × (العرض + الطول)
}
`,
		Hints: []string{
			"الضرب في Go بـ ~*~. يعني ~area := width * height~.",
			"الأقواس بتتحسب الأول: ~2 * (width + height)~.",
			"ابعت المساحة الأول بـ ~colony.Report(area)~ وبعدين المحيط.",
		},
		Lesson: `# دالة بترجّع أكتر من قيمة 📐

حاجة مميزة في Go: الدالة ممكن **ترجّع أكتر من قيمة** مرة واحدة.

~~~go
width, height := colony.WorldSize()
~~~

هنا ~WorldSize~ رجّعت قيمتين، فعملنا متغيرين مرة واحدة. لو مش محتاج واحدة منهم، استخدم ~_~:

~~~go
width, _ := colony.WorldSize() // مش محتاج الطول
~~~

## الأنواع (types)
كل متغير في Go ليه نوع ثابت مبيتغيرش:

~~~go
var count int = 5         // رقم صحيح
var name string = "Kepler" // نص
var ratio float64 = 0.75   // رقم عشري
ready := true              // bool: true أو false
~~~

~var~ هي الطريقة الكاملة، و~:=~ هي الاختصار جوه الدوال. مينفعش تجمع ~int~ على ~string~: Go هيرفض البرنامج، وده بيحميك من أخطاء كتير.

## المهمة
الخريطة حجمها بيتغير كل مرة. ابعت **المساحة** الأول، وبعدها **المحيط**.`,
	})
}
