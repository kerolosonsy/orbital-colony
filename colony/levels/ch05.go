package levels

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch05FunctionsEnglish,
			Title: "Write It Once",
			Goal:  "Harvest the 3 rows (24 moss) using a function you write",
			Lesson: `# Build Your Own Tools 🛠️

So far you have called ready-made functions. Now you will write **your own**:

~~~go
func harvestRow(n int, dir colony.Direction) {
	for i := 0; i < n; i++ {
		colony.Harvest()
		if i < n-1 {
			colony.Move(dir)
		}
	}
}
~~~

- Write ~func~, the name, then the **parameters**, each with its type: ~n int~.
- Calling ~harvestRow(8, colony.East)~ gives ~n~ the value 8 and ~dir~ the value East inside the function.

## Returning a value

~~~go
func area(w, h int) int { // The type after the parentheses is the return type
	return w * h
}

func divmod(a, b int) (int, int) { // Multiple return values
	return a / b, a % b
}
~~~

## Why functions matter
1. **Do not repeat yourself**: keep the same logic in one place.
2. **Names explain intent**: ~harvestRow(8, East)~ is clearer than 15 lines.
3. **Fix a bug** in just one place.

## Capital letters in Go
A name starting with a **capital letter** (~Move~, ~Harvest~) is **exported**, so other packages can use it. A name starting with a **lowercase letter** (~harvestRow~) is private to its package. That is why you write ~colony.Move~, not ~colony.move~.`,
			Starter: `package main

import "orbital/colony"

// harvestRow harvests a whole row of length n, moving in direction dir
func harvestRow(n int, dir colony.Direction) {
	// ...
}

func main() {
	harvestRow(8, colony.East)
	// ...
}
`,
			Hints: []string{
				"Parameters carry changing inputs into one reusable function.",
				"Separate the repeated row work from the transitions in main.",
				"Check the ~harvestRow~ signature and match each call to its parameter types.",
			},
		},
		ID: "5-1-functions", Chapter: 5, LessonRevision: 1, Steps: ch05FunctionsArabic, Title: "اكتبها مرة",
		Goal: "احصد التلات صفوف (24 طحلب) باستخدام function إنت كاتبها",
		Par:  65, Gold: 49, MaxTicks: 300,
		Requires: []string{"func"},
		Setup: func(w *sim.World) {
			w.Init(8, 5, sim.Regolith)
			for _, y := range []int{0, 2, 4} {
				for x := 0; x < 8; x++ {
					w.At(x, y).Ground = sim.Soil
					w.PlantAt(x, y, sim.Moss, 0)
				}
			}
			w.Check = func(w *sim.World) (bool, string) {
				n := w.Inv[sim.Moss]
				return n >= 24, fmt.Sprintf(i18n.T("طحالب: %d / 24", "Moss: %d / 24"), n)
			}
		},
		Starter: `package main

import "orbital/colony"

// harvestRow تحصد صف كامل طوله n وهي ماشية في الاتجاه dir
func harvestRow(n int, dir colony.Direction) {
	// ...
}

func main() {
	harvestRow(8, colony.East)
	// ...
}
`,
		Hints: []string{
			"الـparameters بتنقل المدخلات المتغيرة لـfunction واحدة تقدر تستخدمها أكتر من مرة.",
			"افصل شغل الصف المتكرر عن الانتقالات في main.",
			"راجع توقيع ~harvestRow~ وخلي كل نداء يطابق أنواع الـparameters.",
		},
		Lesson: `# اعمل أدواتك بنفسك 🛠️

لحد دلوقتي كنت بتستخدم دوال جاهزة. دلوقتي هتكتب **دالتك**:

~~~go
func harvestRow(n int, dir colony.Direction) {
	for i := 0; i < n; i++ {
		colony.Harvest()
		if i < n-1 {
			colony.Move(dir)
		}
	}
}
~~~

- ~func~ ثم الاسم، وبعدين **المدخلات (parameters)** وكل واحد بنوعه: ~n int~.
- لما تنادي ~harvestRow(8, colony.East)~، ~n~ بتبقى 8 و~dir~ بتبقى East جوه الدالة.

## دالة بترجّع قيمة

~~~go
func area(w, h int) int { // النوع بعد الأقواس = نوع الناتج
	return w * h
}

func divmod(a, b int) (int, int) { // أكتر من ناتج
	return a / b, a % b
}
~~~

## ليه الدوال مهمة؟
1. **متكرّرش نفسك**: نفس المنطق، مكان واحد.
2. **الأسماء بتشرح**: ~harvestRow(8, East)~ مفهومة أكتر من 15 سطر.
3. **لو فيه bug** بتصلّحه في مكان واحد.

## الحروف الكبيرة في Go
الاسم اللي بيبدأ بحرف **كبير** (~Move~، ~Harvest~) بيبقى **exported**، يعني packages تانية تقدر تستخدمه. واللي بحرف **صغير** (~harvestRow~) خاص بالـ package بتاعته. عشان كده إنت بتكتب ~colony.Move~ مش ~colony.move~.`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch05StormEnglish,
			Title: "Crystal Storm",
			Goal:  "Collect all the crystal in the storm and report the total collected",
			Lesson: `# Errors Are Ordinary Values ⛈️

Many languages raise exceptions when something goes wrong. In Go, a function that can fail **returns the error as a value**, and you decide what to do.

~~~go
n, err := colony.TryHarvest()
if err != nil {
	// Something went wrong; err contains the details
	return
}
// Here we know n is valid
~~~

- ~error~ is a type, just like ~int~ and ~string~.
- ~nil~ means "no value." If ~err == nil~, there is no error.
- You will see this pattern (~if err != nil~) throughout **real Go code**.

## How do you identify an error?

~~~go
switch {
case errors.Is(err, colony.ErrStorm):
	// A storm. Try again
case errors.Is(err, colony.ErrEmpty):
	// Nothing here
}
~~~

~errors.Is~ compares an error with known errors called **sentinel errors**.

## Retrying

~~~go
for {
	n, err := colony.TryHarvest()
	if err == nil {
		total += n
		break // Leave the loop
	}
}
~~~

~break~ leaves the nearest ~for~ loop. ~continue~ skips to the next iteration.

## Your mission
Regular ~Harvest~ is disabled here. Use ~TryHarvest~, handle the storm, and report the total crystal collected.`,
			Starter: `package main

import (
	"errors"
	"fmt"
	"orbital/colony"
)

func main() {
	n, err := colony.TryHarvest()
	if errors.Is(err, colony.ErrStorm) {
		fmt.Println("Storm!", err)
	}
	fmt.Println("Collected:", n)
}
`,
			Hints: []string{
				"A returned error tells you whether the operation produced a usable result.",
				"Separate failure handling from the logic that records successful crystal.",
				"Inspect both results from ~TryHarvest~ and compare causes with ~errors.Is~.",
			},
		},
		ID: "5-2-storm", Chapter: 5, LessonRevision: 1, Steps: ch05StormArabic, Title: "عاصفة الكريستال",
		Goal: "اجمع كل الكريستال وسط العاصفة، وابعت مجموع اللي جمعته",
		Par:  35, Gold: 25, MaxTicks: 200,
		Requires: []string{"errcheck"},
		Setup: func(w *sim.World) {
			w.Init(8, 2, sim.Regolith)
			w.StormChance = 0.35
			w.DisableHarvest = true
			sum := 0
			for x := 0; x < 8; x++ {
				t := w.At(x, 0)
				t.Entity, t.Size = sim.Crystal, 1+w.Rand.Intn(9)
				sum += t.Size
			}
			want := strconv.Itoa(sum)
			w.Check = func(w *sim.World) (bool, string) {
				left := countEntity(w, sim.Crystal)
				if left > 0 {
					return false, fmt.Sprintf(i18n.T("كريستال فاضل: %d", "Crystal remaining: %d"), left)
				}
				if !reported(w, want) {
					return false, i18n.T("الكريستال اتجمع. القاعدة مستنية المجموع", "Crystal collected. Base is waiting for the total")
				}
				return true, i18n.T("تمام", "All done")
			}
		},
		Starter: `package main

import (
	"errors"
	"fmt"
	"orbital/colony"
)

func main() {
	n, err := colony.TryHarvest()
	if errors.Is(err, colony.ErrStorm) {
		fmt.Println("عاصفة!", err)
	}
	fmt.Println("جمعت:", n)
}
`,
		Hints: []string{
			"الـerror الراجع بيقولك لو العملية رجّعت value ينفع تستخدمها.",
			"افصل معالجة الفشل عن تسجيل الكريستال اللي اتجمع بنجاح.",
			"راجع النتيجتين من ~TryHarvest~ وقارن الأسباب بـ~errors.Is~.",
		},
		Lesson: `# الأخطاء قيم عادية ⛈️

في لغات كتير الخطأ بيحصل فجأة (exception). في Go، الدالة اللي ممكن تفشل **بترجّع الخطأ كقيمة**، وإنت اللي تقرر تعمل إيه.

~~~go
n, err := colony.TryHarvest()
if err != nil {
	// حصلت مشكلة. err فيها التفاصيل
	return
}
// هنا متأكدين إن n صح
~~~

- ~error~ نوع زي ~int~ و~string~.
- ~nil~ معناها "مفيش قيمة". لو ~err == nil~ يبقى مفيش خطأ.
- النمط ده (~if err != nil~) هتشوفه في **كل** كود Go حقيقي.

## تعرف نوع الخطأ إزاي؟

~~~go
switch {
case errors.Is(err, colony.ErrStorm):
	// عاصفة. جرّب تاني
case errors.Is(err, colony.ErrEmpty):
	// مفيش حاجة هنا
}
~~~

~errors.Is~ بتقارن الخطأ بأخطاء معروفة. دي اسمها **sentinel errors**.

## إعادة المحاولة

~~~go
for {
	n, err := colony.TryHarvest()
	if err == nil {
		total += n
		break // اخرج من الحلقة
	}
}
~~~

~break~ بتخرجك من أقرب ~for~. و~continue~ بتنط للفّة اللي بعدها.

## المهمة
~Harvest~ العادية متعطّلة هنا. استخدم ~TryHarvest~، واتعامل مع العاصفة، وابعت مجموع الكريستال اللي جمعته.`,
	})

	faultLogStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.DiagnosticCount())
	// Build the report described in the lesson.
}
`
	add(&Level{ID: "5-3-fault-log", Chapter: 5, LessonRevision: 1, Steps: ch05FaultArabic, Title: "سجل الأعطال", Goal: "لفّ كل error وابعت كوده، وبعدها ابعت سجل الأسباب",
		Requires: []string{"errwrap"}, Bridge: []int{2, 14},
		Par: 0, Gold: 0, MaxTicks: 100,
		Starter: faultLogStarter, Lesson: `# الخطأ يحتفظ بسببه 🧰

في عاصفة الكريستال اتعلّمنا إن الخطأ قيمة، وإن ~errors.Is~ بتعرف سبب معروف. دلوقتي هنضيف مكان العطل للرسالة، من غير ما نضيّع السبب الأصلي.

## قبل ما تكتب
- إزاي تضيف سياق للخطأ وتحافظ على هويته؟
- هتجمع كل الأعطال إزاي من غير ما توقف عند أول واحد؟

## افحص وحدة واحدة
~~~go
count := colony.DiagnosticCount()
err := colony.Diagnose(0)
fmt.Println(count, err)
~~~
~DiagnosticCount~ عدد الوحدات. لف بحلقة عدّ من صفر لحد قبل العدد ده، ونادِ ~Diagnose(i)~ لكل وحدة. النتيجة ~nil~ لو سليمة، أو error لو معطوبة. هنا مش محتاج slice تحفظ فيها الوحدات.

## لف الخطأ بـ %w
~~~go
cause := errors.New("antenna offline")
wrapped := fmt.Errorf("relay %d: %w", 7, cause)
fmt.Println(wrapped) // relay 7: antenna offline
~~~
~errors.New~ بتعمل error برسالة بسيطة. ~fmt.Errorf~ بتعمل error جديد من صيغة: ~%d~ مكان عدد صحيح، و~%w~ مكان السبب اللي هيتحفظ جوه الخطأ الجديد. ضيف ~"errors"~ و~"fmt"~ للـ imports.

~~~go
fmt.Println(errors.Is(wrapped, cause)) // true
~~~
دي نفس ~errors.Is~ اللي عرفناها قبل كده. بعد اللف لسه قادرة توصل للسبب، مش بس تقارن شكل النص.

## هات كود العطل
~~~go
code := colony.FaultCode(wrapped)
fmt.Println(code)
~~~
~FaultCode~ بتفتّش داخل الأسباب الملفوفة باستخدام ~errors.As~ جوه المكتبة، وبترجّع كود كـ string. مع عطل جاي من ~Diagnose~ الكود موجود؛ مع خطأ عادي زي مثال الهوائي بترجّع نص فاضي. تفاصيل النوع الداخلي متشالة عن كود اللاعب، فمش محتاج pointers أو structs.

## اجمع من غير قايمة: مثال محلول
~~~go
var all error
all = errors.Join(all, errors.New("relay A"))
all = errors.Join(all, errors.New("relay B"))
fmt.Println(all)
~~~
المتغير ~all~ بيبدأ nil. ~errors.Join~ بتتجاهل nil وبتجمع الأسباب، فكل مرة بنرجّع القيمة الجديدة لنفس المتغير. الطباعة بتعرض رسالة كل سبب في سطر: relay A وبعدها relay B. دي تجربة لجمع إنذارين، مش حل فحوصات المهمة.

## أخطاء شائعة
~~~go
wrapped := fmt.Errorf("relay: %v", cause)
~~~
~%v~ بتعرض نص السبب بس؛ ~errors.Is(wrapped, cause)~ مش هتلاقيه. استخدم ~%w~ لما تقصد تحافظ على السبب. ولو كتبت ~all = wrapped~ هتمسح الأعطال السابقة بدل ما تجمعها.

## المهمة
افحص كل الوحدات. لكل error لفّه بصيغة ~fmt.Errorf("bay %d: %w", i, err)~. اتأكد بـ ~errors.Is~ إن السبب محفوظ، وابعت ~colony.FaultCode(wrapped)~ كتقرير فوراً، من غير index أو بادئة.

جمّع بـ ~all = errors.Join(all, wrapped)~. بعد الحلقة ابعت ~all~ نفسه كتقرير أخير لو مش nil؛ Report هيعرضه كنص متعدد السطور. في كل تشغيل فيه وحدات سليمة وعطلين على الأقل، ومفيش تقارير للوحدات السليمة.`, Hints: []string{"~%w~ بيضيف context للـerror وبيحافظ على السبب الأصلي.", "افصل الوحدات السليمة عن الـerrors، وخلي ترتيب المرور زي ما هو.", "راجع ~colony.FaultCode~ علشان تقرأ الكود من error ملفوف."},
		En: &LevelText{LessonRevision: 1, Steps: ch05FaultEnglish, Title: "Fault Log", Goal: "Wrap faults, report each code, then the combined causes", Starter: faultLogStarter, Lesson: `# Keep the Cause Inside the Error 🧰

Crystal Storm taught us that errors are values and ~errors.Is~ recognizes known causes. Now add the fault location without losing the original cause.

## Before you write
- How can you add context while preserving error identity?
- How will you collect every fault without stopping at the first?

## Inspect one bay
~~~go
count := colony.DiagnosticCount()
err := colony.Diagnose(0)
fmt.Println(count, err)
~~~
~DiagnosticCount~ gives the bay count. Use a counting loop from zero up to that count and call ~Diagnose(i)~ for each bay. It returns ~nil~ for a healthy bay or an error for a faulty one. You do not need a slice of bays.

## Wrap with %w
~~~go
cause := errors.New("antenna offline")
wrapped := fmt.Errorf("relay %d: %w", 7, cause)
fmt.Println(wrapped) // relay 7: antenna offline
~~~
~errors.New~ creates an error with a simple message. ~fmt.Errorf~ creates a new error from a format: ~%d~ inserts an integer, while ~%w~ preserves the cause inside the new error. Add ~"errors"~ and ~"fmt"~ to your imports.

~~~go
fmt.Println(errors.Is(wrapped, cause)) // true
~~~
This is the same ~errors.Is~ from the earlier mission. It can still reach the cause after wrapping, rather than merely comparing message text.

## Read the fault code
~~~go
code := colony.FaultCode(wrapped)
fmt.Println(code)
~~~
~FaultCode~ searches wrapped causes using ~errors.As~ inside the library and returns a string code. A fault from ~Diagnose~ has a code; an ordinary error like the antenna example returns empty text. The internal type stays behind the API, so player code needs no pointers or structs.

## Accumulate without a collection: a worked example
~~~go
var all error
all = errors.Join(all, errors.New("relay A"))
all = errors.Join(all, errors.New("relay B"))
fmt.Println(all)
~~~
~all~ starts as nil. ~errors.Join~ ignores nil and combines causes, so assign its new value back to the same variable each time. Printing shows one message per line: relay A, then relay B. This combines two sample alarms rather than solving the mission diagnostics.

## Common mistakes
~~~go
wrapped := fmt.Errorf("relay: %v", cause)
~~~
~%v~ displays the cause's text but ~errors.Is(wrapped, cause)~ will not find it. Use ~%w~ to preserve the cause. Also, ~all = wrapped~ discards earlier faults instead of accumulating them.

## Your mission
Inspect every bay. Wrap each failure with ~fmt.Errorf("bay %d: %w", i, err)~. Check with ~errors.Is~ that its cause remains accessible and immediately report ~colony.FaultCode(wrapped)~, without an index or prefix.

Accumulate with ~all = errors.Join(all, wrapped)~. After the loop, report ~all~ itself once if it is non-nil; Report renders it as multiline text. Every run includes healthy bays and at least two faults. Healthy bays produce no report.`, Hints: []string{"~%w~ adds context while preserving the original cause.", "Separate healthy bays from faults and keep inspection order.", "Use ~colony.FaultCode~ to read a code from a wrapped error."}},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			n := 5 + w.Rand.Intn(4)
			w.Data.Faults = make([]*sim.FaultError, n)
			var want []string
			var all error
			for i := 0; i < n; i++ {
				if i == 1 || i == n-1 || (i > 1 && w.Rand.Intn(2) == 0) {
					fault := &sim.FaultError{Index: i, Code: fmt.Sprintf("F%d", 100+w.Rand.Intn(900)), Message: "sensor offline"}
					w.Data.Faults[i] = fault
					want = append(want, fault.Code)
					all = errors.Join(all, fmt.Errorf("bay %d: %w", i, fault))
				}
			}
			want = append(want, all.Error())
			w.Check = func(w *sim.World) (bool, string) {
				if slices.Equal(w.Reports, want) {
					return true, i18n.T("كل التقارير مطابقة وبالترتيب", "All reports match in order")
				}
				return false, i18n.T("التقارير لسه ناقصة أو مش مطابقة؛ راجع القيم وترتيبها", "Reports are missing or do not match; check values and order")
			}
		},
	})

	undockStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.DrillAttempts())
	// Build the report described in the lesson.
}
`
	add(&Level{ID: "5-4-undock", Chapter: 5, LessonRevision: 1, Steps: ch05UndockArabic, Title: "فك الارتباط", Goal: "افصل الدرون بعد كل محاولة حفر وكمّل بعد panic",
		Requires: []string{"defer", "recover"}, Bridge: []int{},
		Par: 1, Gold: 1, MaxTicks: 20,
		Starter: undockStarter, Lesson: `# حتى لو الحفر اتعطّل 🔧

بعض الأعطال مش بترجع error: ~Drill~ ممكن تعمل ~panic~ فتوقف مسار التنفيذ العادي. الدرون محتاج حد أمان صغير يحافظ على فك الارتباط ويخلّي المهمة تكمل.

## قبل ما تكتب
- تحط التنضيف في دالة المحاولة ولا في main؟
- بعد التعافي، هتكمل جوه المحاولة ولا عند المنادي؟

## عدد المحاولات وأوامر الربط
~~~go
attempts := colony.DrillAttempts()
colony.Dock()
~~~
~DrillAttempts~ بترجّع عدد صحيح للمحاولات المطلوبة. ~Dock~ بتربط الدرون عشان يبدأ الشغل.

~~~go
colony.Drill()
~~~
~Drill~ بتسجّل محاولة؛ ممكن تكمل عادي أو تعمل panic. لازم تكون مربوط قبل ندائها. حالة المثقاب ممكن تختلف بين تشغيل والتاني.

~~~go
colony.Undock()
~~~
~Undock~ بتفك الربط. مهم إنها تحصل مرة لكل محاولة حتى لو Drill ما رجعتش عادي.

## defer بتضمن خطوة عند الخروج
~~~go
func announce() {
 defer fmt.Println("finished")
 fmt.Println("working")
}
~~~
~defer~ بتسجّل نداء يتنفّذ لما الدالة تخرج، حتى أثناء panic. المثال بيطبع working وبعدها finished. لو سجلت أكتر من defer، آخر واحدة اتسجّلت بتتنفّذ الأول.

## recover وحدّ الأمان: مثال محلول
~~~go
func trainingSignal() {
 defer fmt.Println("signal closed")
 defer func() {
  issue := recover()
  if issue != nil {
   fmt.Println("recovered:", issue)
  }
 }()
 panic("training signal")
}
~~~
~panic~ في المثال بتطلق عطل تدريب. ~func() { ... }~ دالة من غير اسم، والأقواس ~()~ في آخرها بتجهّز نداءها المؤجل. ~recover~ لازم تتنادى مباشرة جوه الدالة المؤجلة دي؛ بترجّع سبب العطل أو nil لو مفيش panic.

المثال بيطبع السبب وبعدين signal closed. التعافي بيرجّع ~trainingSignal~ للمنادي؛ مش بيكمّل بعد سطر panic. عشان كده حلقة المحاولات لازم تبقى في main، وكل محاولة في دالة لوحدها. احنا بنتعافى من عطل المثقاب المتوقع عشان ننضف ونكمّل.

## تجهيز نص التقرير
~~~go
text := fmt.Sprintf("training:%d", 3)
fmt.Println(text) // training:3
~~~
~Sprintf~ بترجّع string جاهزة بدل ما تطبعها. ~%d~ مكان العدد، زي صيغ Errorf في سجل الأعطال.

## أخطاء شائعة
~~~go
func attempt() {
 recover() // مش جوه defer
 panic("training")
}
~~~
النداء ده مش هيمسك panic اللي بعدها. وكمان لو أخّرت ~Undock~ لآخر main، المحاولات مش هتتفك واحدة واحدة. سجّل التنضيف جوه دالة المحاولة بعد Dock مباشرة.

## المهمة
كرّر عدد ~DrillAttempts()~ من المحاولات. كل واحدة تعمل Dock وDrill، وتضمن Undock مرة واحدة مع النجاح أو panic. بعد ما الكل يخلص، اتحرك East مرة واحدة، وابعت تقرير واحد ~ready:N~، حيث N عدد المحاولات اللي قريته. متبعتش تقارير تانية؛ رسائل التشخيص ينفع تطبعها بـ Println.`, Hints: []string{"~defer~ بينفّذ التنضيف لما الـfunction تخلص، حتى لو حصل panic.", "خلّي كل محاولة جوه حدود واضحة علشان اللي بيناديها يقدر يكمل.", "راجع مكان ~recover~ جوه الـdeferred function وسطر Dock في الـstarter."},
		En: &LevelText{LessonRevision: 1, Steps: ch05UndockEnglish, Title: "Undocking", Goal: "Undock after every drill attempt and continue after panic", Starter: undockStarter, Lesson: `# Even When Drilling Fails 🔧

Some faults do not return errors: ~Drill~ can ~panic~, interrupting normal execution. The drone needs a small safety boundary that guarantees cleanup and lets the mission continue.

## Before you write
- Should cleanup belong to each attempt or to main?
- After recovery, does execution resume inside the attempt or at its caller?

## Attempt count and docking commands
~~~go
attempts := colony.DrillAttempts()
colony.Dock()
~~~
~DrillAttempts~ returns the required integer attempt count. ~Dock~ attaches the drone before work begins.

~~~go
colony.Drill()
~~~
~Drill~ records an attempt and can either return normally or panic. Dock before calling it. Drill status can differ between runs.

~~~go
colony.Undock()
~~~
~Undock~ detaches the drone. It must happen once per attempt even if Drill does not return normally.

## defer guarantees a step on exit
~~~go
func announce() {
 defer fmt.Println("finished")
 fmt.Println("working")
}
~~~
~defer~ registers a call for function exit, including panic unwinding. This prints working, then finished. When several defers are registered, the last registered runs first.

## recover and the safety boundary: a worked example
~~~go
func trainingSignal() {
 defer fmt.Println("signal closed")
 defer func() {
  issue := recover()
  if issue != nil {
   fmt.Println("recovered:", issue)
  }
 }()
 panic("training signal")
}
~~~
~panic~ raises a training fault here. ~func() { ... }~ is an unnamed function, and its final ~()~ prepares the deferred call. ~recover~ must be called directly inside that deferred function; it returns the panic cause or nil when no panic occurred.

The example prints the cause and then signal closed. Recovery returns ~trainingSignal~ to its caller; execution does not resume after the panic line. Keep the attempt loop in main and each attempt in its own function. We recover from the expected drill fault to clean up and continue.

## Format the report text
~~~go
text := fmt.Sprintf("training:%d", 3)
fmt.Println(text) // training:3
~~~
~Sprintf~ returns a formatted string instead of printing it. ~%d~ inserts the integer, as in Errorf formats from Fault Log.

## Common mistakes
~~~go
func attempt() {
 recover() // Not inside a defer
 panic("training")
}
~~~
This call does not catch the later panic. Also, postponing ~Undock~ until main exits fails to clean up attempts individually. Register cleanup inside the attempt immediately after Dock.

## Your mission
Perform ~DrillAttempts()~ attempts. Each calls Dock and Drill and guarantees exactly one Undock on success or panic. After all attempts finish, move East once and send one report ~ready:N~, where N is the attempt count you read. Send no other reports; diagnostic messages can use Println.`, Hints: []string{"~defer~ runs cleanup as a function exits, including during panic unwinding.", "Keep one attempt inside a clear boundary so its caller can continue.", "Check where ~recover~ sits inside the deferred function and where the starter docks."}},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			n := 2 + w.Rand.Intn(5)
			w.Data.DrillAttempts = n
			w.Data.DrillPanic = w.Rand.Intn(2) == 0
			want := fmt.Sprintf("ready:%d", n)
			w.Check = func(w *sim.World) (bool, string) {
				if w.Data.Drills == n && w.Data.Undocks == n && !w.Data.Docked && w.X == 1 && w.Y == 0 && len(w.Reports) == 1 && w.Reports[0] == want {
					return true, i18n.T("كل محاولة اتفكت والدرون كمل رحلته", "Every attempt was cleaned up and the drone continued")
				}
				return false, i18n.T("مستني حفر وفك ارتباط لكل محاولة، وبعدهم حركة وتقرير الجاهزية", "Waiting for drilling and cleanup per attempt, then movement and readiness report")
			}
		},
	})

	commandTableStarter := `package main

import (
	"fmt"
	"orbital/colony"
)

func main() {
	fmt.Println(colony.OrderCount())
	// Build the report described in the lesson.
}
`
	add(&Level{ID: "5-5-command-table", Chapter: 5, LessonRevision: 1, Steps: ch05CommandArabic, Title: "لوحة الأوامر", Goal: "اختار command functions ومرّرها وعدّ المعروف منها بـclosure",
		Requires: []string{"funcvalue"}, Bridge: []int{11},
		Par: 10, Gold: 10, MaxTicks: 30,
		Starter: commandTableStarter, Lesson: `# الدالة ممكن تبقى قيمة 🎛️

الدالة مش لازم تتنفّذ لحظة ما تختارها. ممكن تحتفظ بيها، تبعتها لدالة تانية، أو ترجعها كنتيجة. ده هيخلّي لوحة الأوامر تختار الشغل الأول وتنفّذه بعدين.

## قبل ما تكتب
- إيه الفرق بين اسم الدالة وندائها بالأقواس؟
- العدّاد اللي جوه closure لازم يعيش فين عشان يفتكر النداءات السابقة؟

## خزّن الدالة نفسها
~~~go
func beep() { fmt.Println("ping") }
~~~
دي دالة عادية زي اللي كتبناها في 5-1. جوه main نقدر نعمل:
~~~go
action := beep
action()
~~~
~beep~ من غير أقواس قيمة نوعها ~func()~: دالة من غير مدخلات ولا ناتج. ~action()~ هو اللي بينفّذها. الدالة المجهولة اللي استخدمناها مع defer في المهمة اللي فاتت ينفع كمان نخزنها كقيمة.

## رجّع دالة ومرّر دالة: مثال محلول
~~~go
func choose(name string) func() {
 switch name {
 case "beep":
  return beep
 }
 return nil
}
func run(action func()) {
 action()
}
~~~
~func()~ بعد مدخلات choose هي نوع الناتج. ~run~ بتاخد دالة كـ parameter وتنفّذها. قيمة nil معناها مفيش دالة للاسم ده.

~~~go
action := choose("beep")
if action != nil {
 run(action) // ping
}
~~~
المثال اختار صوت تجريبي ومرّره للتنفيذ؛ مش لوحة أوامر المهمة.

## closure بتفتكر متغير برّه جسمها
~~~go
hits := 0
countHit := func() { hits++ }
run(countHit)
run(countHit)
fmt.Println(hits) // 2
~~~
~countHit~ بتستخدم hits اللي اتعرّفت برّه الدالة. دي closure: نفس المتغير بيتحدّث مع كل نداء. لو عرّفته جوّاها هيتصفّر كل مرة.

## اقرأ ونضّف اسم الأمر
~~~go
n := colony.OrderCount()
name := colony.Order(0)
command := strings.TrimSpace(name)
fmt.Println(n, command)
~~~
~OrderCount~ عدد الأوامر، و~Order(i)~ بترجّع نص أمر واحد. ضيف ~"strings"~؛ ~TrimSpace~ بتشيل المسافات من الأطراف بس، ومش بتغيّر الحروف الكبيرة. لف بحلقة عدّ عادية، من غير map أو slice.

صيغة التقرير ممكن تتجهّز بـ ~fmt.Sprintf("id=%d", 7)~ فتطلع ~id=7~. للمجهول، ~fmt.Errorf("unknown:%s", command)~ بتعمل error و~%s~ بتحط النص مكانها.

## أخطاء شائعة
~~~go
run(beep()) // غلط
~~~
ده بينادي beep بدري ويحاول يمرّر ناتج مش موجود. الصح تمرّر قيمة الدالة. وكمان متندهش على action وهي nil؛ افحصها بـ if الأول.

## المهمة
اكتب دالة بتاخد اسم الأمر وترجّع ~func()~، واختار جواها بـ switch: ~moss~ تبعت ~moss=N~ من ~Count(Moss)~، و~crystal~ تبعت ~crystal=N~، و~wait~ تعمل ~colony.Wait()~ وبعدين تبعت ~wait~.

مرّر الـfunctions لـfunction تنفيذ. استخدم closure واحدة بتلتقط عدّاد الأوامر المعروفة، ومرّرها للتنفيذ بعد كل أمر معروف. للأمر المجهول، حتى الفاضي، ابعت error بصيغة ~unknown:command~ وكمّل من غير ما تزود العدّاد. النتيجة تقرير لكل أمر بالترتيب، وبعدهم تقرير أخير بالعدّاد كعدد صحيح. كل التكرارات بتتنفّذ؛ شيل مسافات الأطراف، لكن ~WAIT~ مش ~wait~.`, Hints: []string{"اسم الـfunction من غير أقواس هو value تقدر تخزّنها وتمرّرها.", "افصل اختيار الأمر عن وقت التنفيذ، وخلي النوع متوافق.", "راجع ~Order(i)~ ونطاق المتغير اللي الـclosure محتاجة تفتكره."},
		En: &LevelText{LessonRevision: 1, Steps: ch05CommandEnglish, Title: "Command Table", Goal: "Choose and pass command functions, counting valid commands with a closure", Starter: commandTableStarter, Lesson: `# A Function Can Be a Value 🎛️

A function does not have to run when you choose it. You can store it, pass it to another function, or return it. The command panel will choose work first and execute it later.

## Before you write
- How does a function name differ from calling it with parentheses?
- Where must a closure's counter live to remember earlier calls?

## Store the function itself
~~~go
func beep() { fmt.Println("ping") }
~~~
This is an ordinary function like those from 5-1. Inside main, write:
~~~go
action := beep
action()
~~~
~beep~ without parentheses is a value of type ~func()~: a function with no inputs or result. ~action()~ executes it. The unnamed function used with defer in the previous mission can also be stored as a value.

## Return and pass functions: a worked example
~~~go
func choose(name string) func() {
 switch name {
 case "beep":
  return beep
 }
 return nil
}
func run(action func()) {
 action()
}
~~~
~func()~ after choose's parameters is its return type. ~run~ accepts a function parameter and executes it. A nil value means there is no function for that name.

~~~go
action := choose("beep")
if action != nil {
 run(action) // ping
}
~~~
This chooses a sample sound and passes it for execution; it does not implement the mission's command panel.

## A closure remembers an outer variable
~~~go
hits := 0
countHit := func() { hits++ }
run(countHit)
run(countHit)
fmt.Println(hits) // 2
~~~
~countHit~ uses hits, declared outside its body. This is a closure: each call updates the same variable. Declaring the counter inside would reset it on every call.

## Read and trim a command
~~~go
n := colony.OrderCount()
name := colony.Order(0)
command := strings.TrimSpace(name)
fmt.Println(n, command)
~~~
~OrderCount~ gives the number of commands, and ~Order(i)~ returns one command string. Import ~"strings"~; ~TrimSpace~ removes surrounding whitespace without changing case. Use an ordinary counting loop, with no map or slice.

Prepare report text with ~fmt.Sprintf("id=%d", 7)~, producing ~id=7~. For unknown commands, ~fmt.Errorf("unknown:%s", command)~ makes an error; ~%s~ inserts the command text.

## Common mistakes
~~~go
run(beep()) // Wrong
~~~
This calls beep early and tries to pass a nonexistent result. Pass the function value instead. Also, never call a nil action; check it with if first.

## Your mission
Write a function taking a command name and returning ~func()~. Choose with switch: ~moss~ reports ~moss=N~ from ~Count(Moss)~; ~crystal~ reports ~crystal=N~; ~wait~ calls ~colony.Wait()~ and then reports ~wait~.

Pass functions to an executor function. Use one closure capturing a valid-command counter, passing it for execution after each known command. For unknown commands, including empty text, report an error formatted as ~unknown:command~ and continue without incrementing the counter. Send one report per command in order, followed by one final integer counter report. Execute repetitions every time; trim surrounding spaces, but ~WAIT~ is not ~wait~.`, Hints: []string{"A function name without parentheses is a value you can store and pass.", "Separate choosing a command from executing it, and keep the types aligned.", "Inspect ~Order(i)~ and the scope of the value your closure must remember."}},
		Setup: func(w *sim.World) {
			w.Init(4, 2, sim.Regolith)
			w.Inv[sim.Moss] = 10 + w.Rand.Intn(80)
			w.Inv[sim.Crystal] = 1 + w.Rand.Intn(40)
			w.Data.Orders = []string{"moss", " wait ", "crystal", "moss", "", fmt.Sprintf("bad%d", w.Rand.Intn(1000)), "WAIT"}
			valid := []string{"moss", "crystal", "wait"}
			for i := 0; i < 1+w.Rand.Intn(3); i++ {
				w.Data.Orders = append(w.Data.Orders, valid[w.Rand.Intn(len(valid))])
			}
			w.Rand.Shuffle(len(w.Data.Orders), func(i, j int) { w.Data.Orders[i], w.Data.Orders[j] = w.Data.Orders[j], w.Data.Orders[i] })
			var want []string
			count, waits := 0, 0
			for _, order := range w.Data.Orders {
				command := strings.TrimSpace(order)
				switch command {
				case "moss":
					want = append(want, fmt.Sprintf("moss=%d", w.Inv[sim.Moss]))
					count++
				case "crystal":
					want = append(want, fmt.Sprintf("crystal=%d", w.Inv[sim.Crystal]))
					count++
				case "wait":
					want = append(want, "wait")
					waits++
					count++
				default:
					want = append(want, "unknown:"+command)
				}
			}
			want = append(want, strconv.Itoa(count))
			w.Check = func(w *sim.World) (bool, string) {
				if slices.Equal(w.Reports, want) && w.Tick == waits {
					return true, i18n.T("الأوامر والعداد وصلوا صح", "Commands and counter are correct")
				}
				return false, i18n.T("راجع تنفيذ كل أمر، تقرير المجهول، وعدّاد الأوامر المعروفة", "Check command actions, unknown reports, and the valid-command counter")
			}
		},
	})

}
