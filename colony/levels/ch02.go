package levels

import (
	"fmt"
	"slices"
	"strconv"
	"unicode/utf8"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{LessonRevision: 1, Steps: ch02InventoryEnglish,
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
				"A variable keeps a changing value available for later calculations.",
				"Read each inventory category, calculate from the saved values, and choose the mission output.",
				"The starter already reads moss; look at ~colony.Count~ for the other category.",
			},
		},
		LessonRevision: 1, Steps: ch02InventoryArabic, ID: "2-1-inventory", Chapter: 2, Title: "جرد المخزن",
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
			"المتغير بيحتفظ بقيمة بتتغير عشان تستخدمها بعدين.",
			"اقرأ كل نوع من المخزن، واحسب من القيم المحفوظة، واختار طريقة إرسال المهمة.",
			"الـstarter قرأ الطحالب بالفعل؛ بص على ~colony.Count~ للنوع التاني.",
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
		En: &LevelText{LessonRevision: 1, Steps: ch02SurveyEnglish,
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
				"A function may return multiple values, each received in matching order.",
				"Use both dimensions to reason about the two requested measurements.",
				"Inspect the existing ~colony.WorldSize()~ line and the required report order.",
			},
		},
		LessonRevision: 1, Steps: ch02SurveyArabic, ID: "2-2-survey", Chapter: 2, Title: "مسح المنطقة",
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
			"الدالة ممكن ترجع أكتر من قيمة، وكل متغير يستقبل قيمته بالترتيب.",
			"استخدم البعدين عشان تفكر في القياسين المطلوبين.",
			"راجع سطر ~colony.WorldSize()~ الموجود وترتيب التقارير المطلوبة.",
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

	callsignsStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.Callsign(0))
	// Build the report described in the lesson.
}
`
	add(&Level{LessonRevision: 1, Steps: ch02CallsignsArabic, ID: "2-3-callsigns", Chapter: 2, Title: "أسماء النداء", Goal: "قارن bytes وrunes وابعت أول وآخر حرف لاسمين",
		Requires: []string{"rune"}, Bridge: []int{6},
		Par: 0, Gold: 0, MaxTicks: 20,
		Starter: callsignsStarter, Lesson: `# النص مش عدد bytes بس 📡

اسم الدرون ممكن يكون عربي أو فيه emoji. النهارده هنشوف ليه مساحة النص في الذاكرة مختلفة عن عدد حروفه، من غير حلقات ولا قوائم.

## قبل ما تكتب
- هل ~len~ بيعد الحروف ولا مساحة تخزينها؟
- أول byte في اسم عربي هو أول حرف كامل؟

## هات اسم واحد
~~~go
name := colony.Callsign(0)
fmt.Println(name)
~~~
~colony.Callsign~ بتاخد رقم اسم وترجّع ~string~ واحدة. عندنا اسمين بس، أرقامهم 0 و1. الاتنين مضمون إنهم مش فاضيين، فمش محتاج تفحص حالة اسم فاضي.

## byte وrune
~~~go
text := "é🚀"
fmt.Println(len(text)) // 6
~~~
~len~ دالة جاهزة بترجّع عدد الـ bytes. UTF-8 ممكن يخزّن حرف واحد في أكتر من byte.

الـ ~rune~ هي Unicode code point: قيمة حرف واحدة. ضيف ~"unicode/utf8"~ للـ imports، وبعدين:
~~~go
fmt.Println(utf8.RuneCountInString(text)) // 2
~~~
~RuneCountInString~ بتعد الـ runes بدل الـ bytes. العلامة المركّبة زي التشكيل بتتعد rune لوحدها؛ إحنا مش بنعد الشكل المرئي المركّب كله كحرف واحد.

## أول وآخر rune: مثال محلول
~~~go
first, size := utf8.DecodeRuneInString(text)
last, _ := utf8.DecodeLastRuneInString(text)
fmt.Println(string(first), size, string(last)) // é 2 🚀
~~~
~DecodeRuneInString~ بترجّع أول rune وعدد الـ bytes اللي خدتها. ~DecodeLastRuneInString~ بتعمل نفس الفكرة لآخر rune. استخدم ~_~ لو مش محتاج العدد، زي ما اتعلّمنا في مسح المنطقة.

~string(first)~ بتحوّل قيمة الـ rune لنص عشان تبعت حرف، بدل ما تبعت رقمه. المثال هنا لنص تجريبي؛ بيانات المهمة لازم تقراها بنفسك.

## أخطاء شائعة
~~~go
first := string(text[0]) // غلط لو هدفك أول حرف
~~~
~text[0]~ بتجيب byte واحدة، مش rune كاملة. مع النص التجريبي دي قطعة من ترميز ~é~، فمش هتاخد الحرف المطلوب. وكمان متستخدمش ~len(name)~ مكان عدد الـ runes.

## المهمة
اقرأ الاسم 0 والاسم 1 بنداءين منفصلين. لكل اسم ابعت أربع تقارير منفصلة بـ ~colony.Report~: عدد الـ bytes، عدد الـ runes، أول rune كنص، وآخر rune كنص، بالترتيب ده.

يعني ثمان تقارير في الآخر: الأربعة بتوع الاسم 0، وبعدهم الأربعة بتوع الاسم 1. مفيش عكس أسماء ولا حلقات في المهمة دي.`, Hints: []string{
			"~len~ بتعد bytes؛ دوال UTF-8 بتعد Unicode code points كـrunes.",
			"اشتغل على اسم واحد كل مرة وخلي تقاريره الأربعة مع بعض.",
			"راجع الـimports ودالتي decode قبل ما تبدأ مجموعة الاسم التاني.",
		},
		En: &LevelText{LessonRevision: 1, Steps: ch02CallsignsEnglish, Title: "Callsigns", Goal: "Report byte length, rune count, and endpoint characters for two callsigns", Starter: callsignsStarter, Lesson: `# Text Is More Than Bytes 📡

A drone callsign can contain Arabic or emoji. Today we will distinguish storage size from character count, without loops or collections.

## Before you write
- Does ~len~ count characters or their storage?
- Is the first byte of an Arabic name a complete character?

## Read one callsign
~~~go
name := colony.Callsign(0)
fmt.Println(name)
~~~
~colony.Callsign~ takes a name index and returns one ~string~. There are exactly two names, numbered 0 and 1. Both are guaranteed nonempty, so you do not need an empty-name condition.

## Byte and rune
~~~go
text := "é🚀"
fmt.Println(len(text)) // 6
~~~
~len~ is a built-in function returning the byte count. UTF-8 can store one character in several bytes.

A ~rune~ is a Unicode code point: one character value. Add ~"unicode/utf8"~ to your imports, then:
~~~go
fmt.Println(utf8.RuneCountInString(text)) // 2
~~~
~RuneCountInString~ counts runes rather than bytes. A combining mark counts as its own rune; we are not counting a combined visual symbol as one character.

## First and last rune: a worked example
~~~go
first, size := utf8.DecodeRuneInString(text)
last, _ := utf8.DecodeLastRuneInString(text)
fmt.Println(string(first), size, string(last)) // é 2 🚀
~~~
~DecodeRuneInString~ returns the first rune and its byte width. ~DecodeLastRuneInString~ does the same for the last rune. Use ~_~ when you do not need the width, just as in Survey the Area.

~string(first)~ converts the rune value into text so you report a character rather than its number. This example uses sample text; read the mission inputs yourself.

## Common mistakes
~~~go
first := string(text[0]) // Wrong for reading the first character
~~~
~text[0]~ reads a single byte, not a whole rune. Here it is only part of the encoding of ~é~, so it does not produce the required character. Also, do not substitute ~len(name)~ for the rune count.

## Your mission
Read callsigns 0 and 1 with separate calls. For each name, send four separate ~colony.Report~ calls: byte length, rune count, first rune as text, and last rune as text, in that order.

There must be eight reports: the four for name 0, followed by the four for name 1. This mission requires no reversal or loops.`, Hints: []string{
			"~len~ counts bytes; UTF-8 rune functions count Unicode code points.",
			"Work through one name at a time and keep its four reports together.",
			"Check the imports and the two decode functions before building the second group.",
		}},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			// Both names are nonempty UTF-8; Arabic, emoji, and combining marks vary by seed.
			names := []string{"قمر", "نور", "مسبار", "نُور"}
			w.Data.Callsigns = []string{fmt.Sprintf("%s%d🚀", names[w.Rand.Intn(len(names))], w.Rand.Intn(10000)), fmt.Sprintf("🌙%d%s", w.Rand.Intn(10000), names[w.Rand.Intn(len(names))])}
			var want []string
			for _, name := range w.Data.Callsigns {
				first, _ := utf8.DecodeRuneInString(name)
				last, _ := utf8.DecodeLastRuneInString(name)
				want = append(want, strconv.Itoa(len(name)), strconv.Itoa(utf8.RuneCountInString(name)), string(first), string(last))
			}
			w.Check = func(w *sim.World) (bool, string) {
				if slices.Equal(w.Reports, want) {
					return true, i18n.T("كل التقارير مطابقة وبالترتيب", "All reports match in order")
				}
				return false, i18n.T("التقارير لسه ناقصة أو مش مطابقة؛ راجع القيم وترتيبها", "Reports are missing or do not match; check values and order")
			}
		},
	})

	unitsStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.CargoAmount(0))
	// Build the report described in the lesson.
}
`
	add(&Level{LessonRevision: 1, Steps: ch02UnitsArabic, ID: "2-4-units", Chapter: 2, Title: "وحدات القياس", Goal: "عرّف أنواع الشحنات وحوّل ثلاث كميات لـ cells",
		Requires: []string{"iota", "const"}, Bridge: []int{3, 8},
		Par: 0, Gold: 0, MaxTicks: 20,
		Starter: unitsStarter, Lesson: `# أسماء للأنواع ووحدات للشحنة 📦

المخزن فيه cells مفردة وpacks وcrates. جمع أعدادهم مباشرة مش هيقولنا كمية الشحنة الحقيقية؛ لازم نوحّد الوحدة الأول.

## قبل ما تكتب
- هتميّز نوع العبوة عن عددها إزاي؟
- هتجمع الكميات قبل التحويل ولا بعده؟

## نوع باسم إنت مختاره
~~~go
type Minutes int
~~~
اكتب ~type~ برّه ~main~ عشان تعرّف نوع جديد أساسه ~int~. الاسم هنا بيوضح إن العدد مدّة، مش أي رقم وخلاص.

~~~go
visits := 3
period := Minutes(visits)
~~~
~Minutes(visits)~ تحويل صريح. النوع الجديد مش بيتخلط تلقائي مع متغير ~int~، حتى لو الاتنين بيخزّنوا أعداد صحيحة.

## const وiota
~~~go
type AccessLevel int
const (
 Guest AccessLevel = iota
 Operator
 Captain
)
~~~
~const~ بتعرّف قيمة ثابتة متقدرش تغيّرها بعد كده. جوه مجموعة الثوابت، ~iota~ بتبدأ من صفر وتزيد واحد في كل سطر: Guest=0 وOperator=1 وCaptain=2.

السطرين اللي بعد الأول بيكرروا نفس شكل التعريف تلقائياً. الثوابت التلاتة نوعهم ~AccessLevel~. ولو دالة مستنية ~int~، تقدر تحوّل الاسم بـ ~int(Operator)~.

## معامل تحويل بنوع واضح: مثال محلول
~~~go
const VisitLength Minutes = 5
visits := Minutes(3)
duration := visits * VisitLength
fmt.Println(duration) // 15
~~~
المعامل ~VisitLength~ ثابت نوعه Minutes. حوّلنا المتغير لنفس النوع، فالحساب واضح ومفيش خلط أنواع. دي زيارة تدريب، مش كميات مخزن المهمة.

## اقرأ الكميات كأرقام
~~~go
amount := colony.CargoAmount(0)
fmt.Println(amount)
~~~
~CargoAmount~ بترجّع ~int~ جاهز؛ مفيش نص محتاج يتحلّل. الأرقام 0 و1 و2 بتحدد العبوات: Cell وPack وCrate. كل كمية موجبة وبتتغيّر مع التشغيل.

## أخطاء شائعة
~~~go
var minutes Minutes = 5
var count int = 2
fmt.Println(minutes + count) // مش هيتبني
~~~
النوعين مختلفين. لازم تحويل صريح قبل الحساب. وفي الشحنات، جمع عدد crates على عدد cells قبل تحويل الوحدة بيدي مجموع ملوش معنى.

## المهمة
عرّف ~CargoKind~ أساسه int، وثوابت ~Cell~ و~Pack~ و~Crate~ بـ iota بالترتيب 0 و1 و2. عرّف ~Cells~ أساسه int ومعاملات ثابتة من نوع Cells: ~CellFactor=1~ و~PackFactor=6~ و~CrateFactor=24~.

اقرأ الكميات بنداءات منفصلة، باستخدام ~int(Cell)~ وباقي أسماء الأنواع. حوّلها لـ Cells واضرب كل كمية في معاملها. ابعت تقرير واحد: إجمالي cells كعدد صحيح. مفيش حلقات ولا parsing ولا أخطاء هنا.`, Hints: []string{
			"الأنواع المسماة بتدي الأعداد والوحدات معاني مختلفة في Go.",
			"افصل تعريف الأنواع والثوابت عن قراءة الكميات اللي بتتغير.",
			"راجع نوع argument في ~CargoAmount~ وحوّل القيم قبل جمعها.",
		},
		En: &LevelText{LessonRevision: 1, Steps: ch02UnitsEnglish, Title: "Cargo Units", Goal: "Define cargo kinds and convert three quantities to cells", Starter: unitsStarter, Lesson: `# Named Types and Cargo Units 📦

Storage contains individual cells, packs, and crates. Adding their counts directly does not tell us the real cargo quantity; first use one common unit.

## Before you write
- How will you distinguish a package kind from its quantity?
- Will you add quantities before or after conversion?

## A type with your own name
~~~go
type Minutes int
~~~
Write ~type~ outside ~main~ to define a new type based on ~int~. Its name tells us this number is a duration rather than an arbitrary count.

~~~go
visits := 3
period := Minutes(visits)
~~~
~Minutes(visits)~ is an explicit conversion. A named type does not automatically mix with an ~int~ variable, even though both store integers.

## const and iota
~~~go
type AccessLevel int
const (
 Guest AccessLevel = iota
 Operator
 Captain
)
~~~
~const~ defines a value you cannot change later. Inside a constant group, ~iota~ starts at zero and increases on each line: Guest=0, Operator=1, Captain=2.

The later lines automatically repeat the first declaration's form. All three constants have type ~AccessLevel~. If a function expects ~int~, convert with ~int(Operator)~.

## A typed conversion factor: a worked example
~~~go
const VisitLength Minutes = 5
visits := Minutes(3)
duration := visits * VisitLength
fmt.Println(duration) // 15
~~~
~VisitLength~ is a constant of type Minutes. We converted the variable to the same type so the arithmetic does not mix types. This is a training visit, not the mission's cargo.

## Read quantities as numbers
~~~go
amount := colony.CargoAmount(0)
fmt.Println(amount)
~~~
~CargoAmount~ returns a ready-to-use ~int~; there is no text to parse. Indices 0, 1, and 2 identify Cell, Pack, and Crate. Each quantity is positive and varies between runs.

## Common mistakes
~~~go
var minutes Minutes = 5
var count int = 2
fmt.Println(minutes + count) // Will not build
~~~
The types differ, so convert explicitly before arithmetic. For cargo, adding crate counts to cell counts before conversion produces a meaningless total.

## Your mission
Define ~CargoKind~ based on int, and ~Cell~, ~Pack~, ~Crate~ using iota in the order 0, 1, 2. Define ~Cells~ based on int and typed Cells constants ~CellFactor=1~, ~PackFactor=6~, ~CrateFactor=24~.

Read the three quantities with separate calls using ~int(Cell)~ and the other kind names. Convert them to Cells and multiply each by its factor. Send one report containing the integer total in cells. No loops, parsing, or error handling are needed.`, Hints: []string{
			"Named types give counts and units different meanings in Go.",
			"Separate the type and constant declarations from reading changing quantities.",
			"Check the ~CargoAmount~ argument type and convert values before combining them.",
		}},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			for kind := range w.Data.CargoAmounts {
				w.Data.CargoAmounts[kind] = 1 + w.Rand.Intn(40)
			}
			want := []string{strconv.Itoa(w.Data.CargoAmounts[0] + 6*w.Data.CargoAmounts[1] + 24*w.Data.CargoAmounts[2])}
			w.Check = func(w *sim.World) (bool, string) {
				if slices.Equal(w.Reports, want) {
					return true, i18n.T("كل التقارير مطابقة وبالترتيب", "All reports match in order")
				}
				return false, i18n.T("التقارير لسه ناقصة أو مش مطابقة؛ راجع القيم وترتيبها", "Reports are missing or do not match; check values and order")
			}
		},
	})

}
