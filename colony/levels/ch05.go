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
				"Inside harvestRow, use ~for i := 0; i < n; i++~ with Harvest, then Move in ~dir~ (except on the last iteration).",
				"There is an empty row between planted rows: move North twice.",
				"Travel ~colony.West~ along the second row, then ~colony.East~ along the third.",
			},
		},
		ID: "5-1-functions", Chapter: 5, Title: "اكتبها مرة",
		Goal: "احصد الـ 3 صفوف (24 طحلب) باستخدام دالة إنت كاتبها",
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
			"جوه harvestRow: حلقة ~for i := 0; i < n; i++~ فيها Harvest، وبعدها Move في ~dir~ (إلا آخر لفّة).",
			"بين الصفوف فيه صف فاضي: اطلع North مرتين.",
			"الصف التاني هتمشي فيه ~colony.West~، والتالت ~colony.East~.",
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
				"The storm can make TryHarvest fail. Use a loop: keep trying until ~err == nil~.",
				"~errors.Is(err, colony.ErrStorm)~ tells you the storm caused the failure, so retrying makes sense.",
				"Add ~n~ to a ~total~ variable outside the loops, then Report it at the end.",
			},
		},
		ID: "5-2-storm", Chapter: 5, Title: "عاصفة الكريستال",
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
			"العاصفة ممكن تفشّل TryHarvest. اعمل حلقة: جرّب لحد ما ~err == nil~.",
			"~errors.Is(err, colony.ErrStorm)~ بتقولك إن الفشل كان بسبب العاصفة، يعني تستاهل تجرّب تاني.",
			"اجمع ~n~ في متغير ~total~ برّه الحلقات، وابعته بـ Report في الآخر.",
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
}
