package levels

import (
	"strings"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

// labRoom is the decorative map shown for lab levels.
func labRoom(w *sim.World) {
	w.Init(7, 3, sim.Rock)
	for x := 1; x < 6; x++ {
		w.At(x, 1).Ground = sim.Ice
		t := w.At(x, 1)
		t.Entity, t.Size = sim.Crystal, x+2
	}
	w.Check = func(*sim.World) (bool, string) {
		return false, i18n.T("المختبر: الكود بيتختبر بـ go test", "Lab: your code is checked with go test")
	}
}

const fuelImpl = `package main

import "fmt"

// Fuel بترجّع كمية الوقود اللازمة لرحلة.
//   - المسافة 0 = 0 وقود، مهما كانت الشحنة
//   - كل خانة مسافة = 2 وقود
//   - كل 5 كجم شحنة كاملين = 1 وقود (الكسر بيتشال: 12 كجم = 2)
//   - لو الشحنة أكتر من 50 كجم: 10 وقود زيادة
func Fuel(distance, cargo int) int {
	if distance == 0 {
		return 0
	}
	f := distance*2 + cargo/5
	if cargo > 50 {
		f += 10
	}
	return f
}

func main() {
	fmt.Println(Fuel(3, 12))
}
`

func mutate(src, from, to string) string {
	if !strings.Contains(src, from) {
		panic("mutant does not apply: " + from)
	}
	return strings.Replace(src, from, to, 1)
}

const medianBuggy = `package main

import (
	"fmt"
	"sort"
)

// Median بترجّع القيمة اللي في النص بعد الترتيب.
//   - لو العدد زوجي: متوسط الاتنين اللي في النص (قسمة صحيحة)
//   - slice فاضية: 0
//   - لازم متغيّرش الـ slice اللي اتبعتتلها
func Median(xs []int) int {
	sort.Ints(xs)
	return xs[len(xs)/2]
}

func main() {
	readings := []int{7, 3, 9}
	fmt.Println(Median(readings), readings)
}
`

const medianHidden = `package main

import (
	"slices"
	"testing"
)

func TestHiddenMedian(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{7, 3, 9}, 7},
		{[]int{4, 1, 3, 2}, 2},
		{[]int{10, 20}, 15},
		{[]int{-5, -1, -3}, -3},
		{[]int{5}, 5},
		{nil, 0},
		{[]int{}, 0},
	}
	for _, c := range cases {
		in := slices.Clone(c.in)
		if got := Median(in); got != c.want {
			t.Errorf("Median(%v) = %d, want %d", c.in, got, c.want)
		}
		if !slices.Equal(in, c.in) {
			t.Errorf("Median(%v) changed its input to %v", c.in, in)
		}
	}
}
`

const pairsNaive = `package main

import "fmt"

// CountPairs بتعدّ كام زوج (i < j) مجموعهم = target.
func CountPairs(xs []int, target int) int {
	n := 0
	for i := range xs {
		for j := i + 1; j < len(xs); j++ {
			if xs[i]+xs[j] == target {
				n++
			}
		}
	}
	return n
}

func main() {
	fmt.Println(CountPairs([]int{1, 9, 5, 5, 3, 7}, 10)) // 3
}
`

const pairsHidden = `package main

import (
	"math/rand"
	"testing"
)

func hiddenNaive(xs []int, target int) int {
	n := 0
	for i := range xs {
		for j := i + 1; j < len(xs); j++ {
			if xs[i]+xs[j] == target {
				n++
			}
		}
	}
	return n
}

func TestHiddenCountPairs(t *testing.T) {
	fixed := []struct {
		in     []int
		target int
	}{
		{nil, 5}, {[]int{5}, 10}, {[]int{5, 5}, 10}, {[]int{5, 5, 5}, 10},
		{[]int{1, 9, 5, 5, 3, 7}, 10}, {[]int{-2, 2, 0, 0}, 0}, {[]int{3, 3, 3, 3}, 6},
	}
	for _, c := range fixed {
		if got, want := CountPairs(c.in, c.target), hiddenNaive(c.in, c.target); got != want {
			t.Errorf("CountPairs(%v, %d) = %d, want %d", c.in, c.target, got, want)
		}
	}
	r := rand.New(rand.NewSource(1))
	for k := 0; k < 300; k++ {
		xs := make([]int, r.Intn(60))
		for i := range xs {
			xs[i] = r.Intn(41) - 20
		}
		target := r.Intn(41) - 20
		if got, want := CountPairs(xs, target), hiddenNaive(xs, target); got != want {
			t.Fatalf("CountPairs(%v, %d) = %d, want %d", xs, target, got, want)
		}
	}
}

var hiddenData = func() []int {
	r := rand.New(rand.NewSource(7))
	xs := make([]int, 10000)
	for i := range xs {
		xs[i] = r.Intn(1000)
	}
	return xs
}()

func BenchmarkHiddenNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hiddenNaive(hiddenData, 1000)
	}
}

func BenchmarkHiddenTheirs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CountPairs(hiddenData, 1000)
	}
}
`

func init() {
	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch16Mission1(true),
			Title: "Catch the Mutants",
			Goal:  "Write Fuel tests that catch every broken version (6 mutants)",
			Lesson: `# The Lab: Real Tests with go test 🧪

There is no drone here. You are in the lab, and your code is checked by real ~go test~. There are two files:
- ~main.go~: the ~Fuel~ function. Read it carefully (read-only).
- ~main_test.go~: your tests.

## Your first test
~~~go
func TestFuel(t *testing.T) {
	if got := Fuel(3, 12); got != 8 {
		t.Errorf("Fuel(3, 12) = %d, want 8", got)
	}
}
~~~
- The filename must end in ~_test.go~, and the function must start with ~Test~ and take ~*testing.T~.
- ~t.Errorf~ records a failure and continues. ~t.Fatalf~ records a failure and stops this test.

## Table-driven tests (the idiomatic Go approach)
~~~go
cases := []struct {
	name            string
	distance, cargo int
	want            int
}{
	{"ordinary", 3, 12, 8},
	{"zero distance", 0, 40, 0},
}
for _, tc := range cases {
	t.Run(tc.name, func(t *testing.T) {
		if got := Fuel(tc.distance, tc.cargo); got != tc.want {
			t.Errorf("got %d, want %d", got, tc.want)
		}
	})
}
~~~
~t.Run~ creates a **subtest** for each case, so you can see exactly which one failed.

## Mutation testing 🧟
How do you know your tests are **good**? Mission Control breaks ~Fuel~ in 6 different ways (changing ~>~ to ~>=~, removing a condition, and so on). Each broken version is a **mutant**. Good tests must **fail** on every mutant. A surviving mutant means your tests miss a case.

Tip: bugs hide at **boundaries** (0, 50, 51...).`,
			Starter: ``,
			Hints: []string{
				"A test is a claim about behavior; a mutant helps reveal a missing claim.",
				"Group related inputs with their expected behavior, then give each row a useful name.",
				"Start from ~TestFuel~ in ~main_test.go~ and use each surviving mutant's description to find the missing promise.",
			},
			Files: []LabFile{
				{Name: "main_test.go", Starter: `package main

import "testing"

func TestFuel(t *testing.T) {
	got := Fuel(3, 12)
	want := 8 // 3*2 + 12/5
	if got != want {
		t.Errorf("Fuel(3, 12) = %d, want %d", got, want)
	}
}
`, ReadOnly: false},
				{Name: "main.go", Starter: `package main

import "fmt"

// Fuel returns the fuel needed for a trip.
//   - Distance 0 means 0 fuel, regardless of cargo
//   - Each tile of distance costs 2 fuel
//   - Each full 5 kg of cargo costs 1 fuel (discard fractions: 12 kg = 2)
//   - Cargo above 50 kg costs 10 extra fuel
func Fuel(distance, cargo int) int {
	if distance == 0 {
		return 0
	}
	f := distance*2 + cargo/5
	if cargo > 50 {
		f += 10
	}
	return f
}

func main() {
	fmt.Println(Fuel(3, 12))
}
`, ReadOnly: true},
			},
		},
		ID: "16-1-mutants", Chapter: 16, LessonRevision: 1, Steps: ch16Mission1(false), Title: "اصطاد الـ Mutants",
		Goal:     "اكتب tests لـ~Fuel~ تمسك كل نسخة متبوّظة منها (6 mutants)",
		Kind:     "lab",
		Requires: []string{"tabletest"},
		Setup:    labRoom,
		Files: []LabFile{
			{Name: "main_test.go", Starter: `package main

import "testing"

func TestFuel(t *testing.T) {
	got := Fuel(3, 12)
	want := 8 // 3*2 + 12/5
	if got != want {
		t.Errorf("Fuel(3, 12) = %d, want %d", got, want)
	}
}
`},
			{Name: "main.go", Starter: fuelImpl, ReadOnly: true},
		},
		Lab: &Lab{
			ReqFile: "main_test.go",
			Mutants: []Mutant{
				{i18n.T("حد الـ 50 كجم: جرّبت 50 بالظبط؟", "The 50 kg boundary: did you try exactly 50?"), mutate(fuelImpl, "cargo > 50", "cargo >= 50")},
				{i18n.T("المسافة صفر مع شحنة", "Zero distance with cargo"), mutate(fuelImpl, "if distance == 0 {\n\t\treturn 0\n\t}\n\t", "")},
				{i18n.T("تكلفة كل خانة مسافة", "Fuel cost per tile of distance"), mutate(fuelImpl, "distance*2", "distance*3")},
				{i18n.T("الكسور في الشحنة (6 كجم مثلاً)", "Cargo fractions (6 kg, for example)"), mutate(fuelImpl, "cargo/5", "(cargo+4)/5")},
				{i18n.T("رسوم الشحنة التقيلة", "Heavy cargo surcharge"), mutate(fuelImpl, "f += 10", "f += 5")},
				{i18n.T("شحنة تقيلة جداً (أكتر من 100)", "Very heavy cargo (over 100)"), mutate(fuelImpl, "if cargo > 50 {", "if cargo > 50 && cargo <= 100 {")},
			},
		},
		Hints: []string{
			"الـtest بيثبت سلوك؛ والـmutant بيكشف الوعد اللي لسه مش واضح.",
			"جمّع المدخلات المرتبطة مع السلوك المتوقع، وادي كل صف اسم مفهوم.",
			"ابدأ من ~TestFuel~ في ~main_test.go~ واستخدم وصف كل mutant عدى عشان تحدد الوعد الناقص.",
		},
		Lesson: `# المختبر: تستات حقيقية بـ go test 🧪

من هنا مفيش درون. إنت في المختبر، والكود بيتختبر بـ ~go test~ الحقيقي. فيه ملفين:
- ~main.go~: الدالة ~Fuel~. اقراها كويس (للقراية بس).
- ~main_test.go~: التستات بتاعتك.

## أول test
~~~go
func TestFuel(t *testing.T) {
	if got := Fuel(3, 12); got != 8 {
		t.Errorf("Fuel(3, 12) = %d, want 8", got)
	}
}
~~~
- الملف لازم ينتهي بـ ~_test.go~، والدالة تبدأ بـ ~Test~ وتاخد ~*testing.T~.
- ~t.Errorf~: سجّل فشل وكمّل. ~t.Fatalf~: سجّل فشل ووقّف الـ test ده.

## Table-driven tests (الطريقة الرسمية في Go)
~~~go
cases := []struct {
	name            string
	distance, cargo int
	want            int
}{
	{"عادي", 3, 12, 8},
	{"مسافة صفر", 0, 40, 0},
}
for _, tc := range cases {
	t.Run(tc.name, func(t *testing.T) {
		if got := Fuel(tc.distance, tc.cargo); got != tc.want {
			t.Errorf("got %d, want %d", got, tc.want)
		}
	})
}
~~~
~t.Run~ بيعمل **subtest** لكل حالة، فبتعرف أنهي حالة بالظبط اللي فشلت.

## Mutation testing 🧟
إزاي تعرف إن التستات بتاعتك **كويسة**؟ Mission Control هياخد ~Fuel~ ويبوّظها 6 مرات (يغيّر ~>~ لـ ~>=~، يشيل شرط، إلخ). كل نسخة متبوّظة اسمها **mutant**. التستات الكويسة لازم **تفشل** على كل mutant. لو mutant عدّى، يبقى فيه حالة إنت مش بتختبرها.

نصيحة: الـ bugs بتستخبى عند **الحدود** (0، 50، 51…).`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch16Mission2(true),
			Title: "Catch the Bug",
			Goal:  "Median has bugs. Write a failing test first, then fix the code until everything passes",
			Lesson: `# TDD: Red First, Then Green 🔴🟢

The drones reported a bug: "Median returns wrong numbers and scrambles our readings." Before fixing it, **prove** the bug exists:

1. 🔴 **Red**: write a test that fails on the current code. Now you have evidence.
2. 🟢 **Green**: fix the code until the test passes.
3. 🔵 **Refactor**: tidy the code while the tests keep passing.

Mission Control checks all three:
- Your tests **fail** on the original code (they catch the bug).
- Your tests pass on your repaired code.
- Mission Control's hidden tests pass on your code.

## A subtle bug: modifying input
~~~go
func Median(xs []int) int {
	sort.Ints(xs) // Sorts the caller's original slice!
	...
}
~~~
Remember chapter 6? Slices refer to shared memory. ~sort.Ints(xs)~ sorts **the caller's data**. Work on a copy (~slices.Clone~).

## Test that input stays unchanged
~~~go
in := []int{7, 3, 9}
Median(in)
if !slices.Equal(in, []int{7, 3, 9}) {
	t.Errorf("Median changed the input: %v", in)
}
~~~

## Testing for panics
~~~go
defer func() {
	if r := recover(); r != nil {
		t.Fatalf("panic: %v", r)
	}
}()
Median(nil)
~~~`,
			Starter: ``,
			Hints: []string{
				"Separate the returned value from whether the function preserves its caller's input.",
				"Make one failing observation first, then repair behavior while keeping both promises visible.",
				"Trace ~Median~ in ~main.go~ around its sort and return, then check which slice each line changes.",
			},
			Files: []LabFile{
				{Name: "main.go", Starter: `package main

import (
	"fmt"
	"sort"
)

// Median returns the middle value after sorting.
//   - Even length: average the middle two values (integer division)
//   - Empty slice: 0
//   - Must not change the input slice
func Median(xs []int) int {
	sort.Ints(xs)
	return xs[len(xs)/2]
}

func main() {
	readings := []int{7, 3, 9}
	fmt.Println(Median(readings), readings)
}
`, ReadOnly: false},
				{Name: "main_test.go", Starter: `package main

import "testing"

func TestMedian(t *testing.T) {
	if got := Median([]int{7, 3, 9}); got != 7 {
		t.Errorf("Median([7 3 9]) = %d, want 7", got)
	}
}
`, ReadOnly: false},
			},
		},
		ID: "16-2-bug-hunt", Chapter: 16, LessonRevision: 1, Steps: ch16Mission2(false), Title: "اصطاد الـ Bug",
		Goal:     "Median فيها bugs. اكتب test يكشفها الأول (يفشل)، وبعدين صلّح الكود لحد ما كله يعدّي",
		Kind:     "lab",
		Requires: []string{"sel:testing.T"},
		Setup:    labRoom,
		Files: []LabFile{
			{Name: "main.go", Starter: medianBuggy},
			{Name: "main_test.go", Starter: `package main

import "testing"

func TestMedian(t *testing.T) {
	if got := Median([]int{7, 3, 9}); got != 7 {
		t.Errorf("Median([7 3 9]) = %d, want 7", got)
	}
}
`},
		},
		Lab: &Lab{ReqFile: "main_test.go", Buggy: medianBuggy, Hidden: medianHidden},
		Hints: []string{
			"افصل الـvalue اللي الـfunction بترجّعها عن سؤال هل حافظت على input بتاع اللي ناداها.",
			"خلّي ملاحظة واحدة تفشل الأول، وبعدها أصلح السلوك مع الحفاظ على الوعدين.",
			"تتبّع ~Median~ في ~main.go~ حوالين الترتيب والـreturn، وشوف كل سطر بيغيّر أنهي slice.",
		},
		Lesson: `# TDD: الأحمر الأول، بعدين الأخضر 🔴🟢

الـ bug اتبلّغ من الدرونات: "Median بترجّع أرقام غلط، وكمان بتلخبط القراءات". قبل ما تصلّح، **اثبت** إن الـ bug موجود:

1. 🔴 **Red**: اكتب test يفشل على الكود الحالي. كده بقى عندك دليل.
2. 🟢 **Green**: صلّح الكود لحد ما الـ test يعدّي.
3. 🔵 **Refactor**: نضّف الكود والتستات لسه بتعدّي.

Mission Control بيتأكد من التلاتة:
- التستات بتاعتك **بتفشل** على الكود الأصلي (يعني بتمسك الـ bug).
- التستات بتاعتك بتعدّي على الكود اللي إنت صلّحته.
- تستات مخفية من Mission Control بتعدّي على كودك.

## bug خبيث: تعديل الـ input
~~~go
func Median(xs []int) int {
	sort.Ints(xs) // بيرتّب الـ slice الأصلية بتاعة اللي ناداك!
	...
}
~~~
فاكر الفصل 6؟ الـ slice بتشاور على نفس الذاكرة. ~sort.Ints(xs)~ بيرتّب **بيانات اللي ناداك**. الحل: اشتغل على نسخة (~slices.Clone~).

## اختبر إن الـ input متغيرش
~~~go
in := []int{7, 3, 9}
Median(in)
if !slices.Equal(in, []int{7, 3, 9}) {
	t.Errorf("Median غيّرت الـ input: %v", in)
}
~~~

## اختبر الـ panic
~~~go
defer func() {
	if r := recover(); r != nil {
		t.Fatalf("panic: %v", r)
	}
}()
Median(nil)
~~~`,
	})

	add(&Level{
		En: &LevelText{
			LessonRevision: 1, Steps: ch16Mission3(true),
			Title: "100 Times Faster",
			Goal:  "CountPairs is very slow (O(n²)). Speed it up without breaking it, and write your own benchmark",
			Lesson: `# Benchmarks: Measure Before Optimizing ⚡

~~~go
func BenchmarkCountPairs(b *testing.B) {
	xs := make([]int, 3000)
	for i := range xs {
		xs[i] = i % 1000
	}
	for b.Loop() { // Go 1.24+ (older: for i := 0; i < b.N; i++)
		CountPairs(xs, 1000)
	}
}
~~~
~go test -bench=.~ runs benchmarks repeatedly until measurements stabilize, then prints **ns/op**: time per call.

## Big-O: why is the code slow?
~~~go
for i := range xs {
	for j := i + 1; j < len(xs); j++ { // A nested loop
~~~
With 10000 elements (Mission Control's benchmark size), that is about 50 million comparisons. Time grows as **n²**: double n and time grows fourfold.

| Approach | Complexity | For 10000 elements |
|---|---|---|
| Two loops | O(n²) | ~50,000,000 |
| Sorting + two pointers | O(n log n) | ~130,000 |
| map | O(n) | ~10,000 |

## The golden rule 🏆
1. **Correctness first**: Mission Control compares your code with the slow version on 300 random cases (negative numbers, duplicates, and empty slices).
2. **Measure**: do not guess. The benchmark decides.
3. **Improve the algorithm** before small details. O(n) beats micro-optimizations of O(n²).

## A surprise 🤯
On paper, a map (O(n)) beats sorting (O(n log n)). But try both and measure: hashing costs time, and sorting small integers is fast, so sorting may win! Big-O describes growth as n gets very large; benchmarks show what happens at your actual n.

## Stars
⭐ 3× faster · ⭐⭐ 15× faster · ⭐⭐⭐ 50× faster`,
			Starter: ``,
			Hints: []string{
				"A benchmark measures repeated runtime; it does not prove result correctness.",
				"Keep correctness and speed checks separate, then compare runs on the same data.",
				"Use the benchmark stub in ~main_test.go~ and inspect the current ~CountPairs~ loops before changing them.",
			},
			Files: []LabFile{
				{Name: "main.go", Starter: `package main

import "fmt"

// CountPairs counts pairs (i < j) whose sum equals target.
func CountPairs(xs []int, target int) int {
	n := 0
	for i := range xs {
		for j := i + 1; j < len(xs); j++ {
			if xs[i]+xs[j] == target {
				n++
			}
		}
	}
	return n
}

func main() {
	fmt.Println(CountPairs([]int{1, 9, 5, 5, 3, 7}, 10)) // 3
}
`, ReadOnly: false},
				{Name: "main_test.go", Starter: `package main

import "testing"

func TestCountPairs(t *testing.T) {
	if got := CountPairs([]int{1, 9, 5, 5, 3, 7}, 10); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

// func BenchmarkCountPairs(b *testing.B) { ... }
`, ReadOnly: false},
			},
		},
		ID: "16-3-benchmark", Chapter: 16, LessonRevision: 1, Steps: ch16Mission3(false), Title: "أسرع 100 مرة",
		Goal:     "~CountPairs~ بطيئة جدًا (O(n²)). خليها أسرع من غير ما تبوّظها، واكتب benchmark بنفسك",
		Kind:     "lab",
		Requires: []string{"benchmark"},
		Setup:    labRoom,
		Files: []LabFile{
			{Name: "main.go", Starter: pairsNaive},
			{Name: "main_test.go", Starter: `package main

import "testing"

func TestCountPairs(t *testing.T) {
	if got := CountPairs([]int{1, 9, 5, 5, 3, 7}, 10); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

// func BenchmarkCountPairs(b *testing.B) { ... }
`},
		},
		Lab: &Lab{ReqFile: "main_test.go", Hidden: pairsHidden, Bench: &Bench{Tiers: [3]float64{3, 15, 50}}},
		Hints: []string{
			"الـbenchmark بيقيس التشغيل المتكرر؛ لكنه مش دليل على صحة النتيجة.",
			"افصل فحص الصحة عن السرعة، وقارن التشغيل على نفس البيانات.",
			"استعمل الـbenchmark stub في ~main_test.go~ وبص على loops ~CountPairs~ الحالية قبل التغيير.",
		},
		Lesson: `# Benchmarks: قيس قبل ما تحسّن ⚡

~~~go
func BenchmarkCountPairs(b *testing.B) {
	xs := make([]int, 3000)
	for i := range xs {
		xs[i] = i % 1000
	}
	for b.Loop() { // Go 1.24+ (القديم: for i := 0; i < b.N; i++)
		CountPairs(xs, 1000)
	}
}
~~~
~go test -bench=.~ بيشغّلها كذا مرة لحد ما القياس يثبت، ويطبعلك **ns/op**: الوقت لكل استدعاء.

## Big-O: ليه الكود بطيء؟
~~~go
for i := range xs {
	for j := i + 1; j < len(xs); j++ { // حلقة جوه حلقة
~~~
مع 10000 عنصر (اللي Mission Control بيقيس بيهم) ده ~50 مليون مقارنة. الوقت بيكبر زي **n²**: ضاعف الـ n والوقت يتضرب في 4.

| الطريقة | التعقيد | لـ 10000 عنصر |
|---|---|---|
| حلقتين | O(n²) | ~50,000,000 |
| ترتيب + مؤشرين | O(n log n) | ~130,000 |
| map | O(n) | ~10,000 |

## القاعدة الذهبية 🏆
1. **صح الأول**: Mission Control هيقارن كودك بالنسخة البطيئة على 300 حالة عشوائية (أرقام سالبة، ومتكررة، وslices فاضية).
2. **قيس**: متخمّنش. الـ benchmark هو الحَكَم.
3. **حسّن الخوارزمية** قبل التفاصيل الصغيرة. O(n) بتكسب أي micro-optimization على O(n²).

## مفاجأة 🤯
على الورق الـ map (O(n)) أحسن من الترتيب (O(n log n)). بس جرّب الاتنين وقيس: الـ hashing ليه تكلفة، والترتيب على أرقام صغيرة سريع جداً، فممكن الترتيب يكسب! Big-O بيقولك الكود بيتصرف إزاي لما n تكبر أوي، والـ benchmark بيقولك الحقيقة عند الـ n بتاعتك.

## النجوم
⭐ أسرع 3 مرات · ⭐⭐ أسرع 15 مرة · ⭐⭐⭐ أسرع 50 مرة`,
	})
	add(&Level{ID: "16-4-wild-inputs", Chapter: 16, LessonRevision: 1, Steps: ch16Mission4(false), Kind: "lab", Title: "مدخلات مجنونة", Goal: "اكتب fuzz يكشف عيوب كبسولة التحكم وبعدين صلّح parser", Requires: []string{"fuzz"}, Bridge: []int{85},
		Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: capsuleBuggy, ReadOnly: false}, {Name: "main_test.go", Starter: wildInputsTestStarter}},
		Lesson: `# جرّب مدخلات ما تخطرش على بالك

الـ table tests بتجرّب حالات اخترتها. الـ fuzzing بيغيّر بايتات عشان يلاقي panic أو خرق للعقد. parser كبسولة التحكم دالة نقية، فممكن نجرب آلاف المدخلات من غير أجهزة أو شبكة.

## قبل ما تكتب
- إيه الشرط اللي لازم يفضل صحيح لأي مدخل، حتى البايظ؟
- هتفرّق رفض مدخل بايظ عن parser غلط بترفض كل حاجة إزاي؟

## شكل fuzz target
~~~go
func FuzzPrefix(f *testing.F) {
 f.Add("ion")
 f.Add("")
 f.Fuzz(func(t *testing.T,s string) {
  if len("x"+s)!=len(s)+1 {t.Fatal("length invariant")}
 })
}
~~~
~testing.F~ مسؤول عن الـ corpus. ~f.Add~ بتضيف seed بنفس نوع parameter اللي في callback. ~f.Fuzz~ بتاخد دالة أول parameter فيها testing.T وبعده المدخل اللي هيتغيّر؛ هنا string، وفي المهمة []byte. المثال بيثبت خاصية طول، مش بيفك رسالة.

~go test~ العادي بيشغّل seeds؛ مرحلة Lab.Fuzz بتشغّل FuzzCapsule خمس ثواني. الـ bug الأصلي والـ mutants بيتجرّبوا على seeds كمان، فاختار مدخلات تكشفهم حتى قبل التوليد.

## افحص قبل الفهرسة
~~~go
func firstByte(packet []byte) (byte,bool) {
 if len(packet)==0 {return 0,false}
 return packet[0],true
}
~~~
~byte~ اسم لنوع uint8، وقيمة زي ~0xA7~ رقم hexadecimal. الطول المكتوب في header مش ضمان لوجود payload؛ حوّله لـ int وافحص الطول الفعلي قبل slicing. المدخل المرفوض يرجّع error، مش panic.

## property بقبول ورفض
~~~go
before:=slices.Clone(packet)
if !bytes.Equal(packet,before) {t.Fatal("input changed")}
~~~
استورد ~bytes~ للمقارنة و~slices~ لنسخ المدخل بـ ~slices.Clone~ زي درس Median. قارن نسخة قبل وبعد الاستدعاء الحقيقي في التست. احسب صلاحية المدخل من المواصفات، وبعدها راجع الخطأ والنتيجة؛ عدم الـ panic لوحده يسمح لحل يرفض كل حاجة. على النجاح جرّب تغيير الـ payload الناتجة عشان تتأكد إنها مش بتشير للمدخل.

## أخطاء شائعة
~~~go
n:=int(frame[2]) // panic لو أقل من3 بايت
return frame[3:3+n] // الطول المعلن ممكن يبقى أكبر من الموجود
~~~
ومتكتبش recover في fuzz لإخفاء panic: ده العيب اللي محتاج التست يكشفه. لازم رفض البايتات الزيادة كمان، مش الرسائل الناقصة بس.

## المهمة
نفّذ ~ParseCapsule(frame []byte) (byte,[]byte,error)~: header ثلاث بايتات: توقيع0xA7، ثم kind يساوي1 أو2، ثم طول payload من0 لـ16. الطول الكلي لازم يساوي3+الطول المعلن بالظبط. النجاح يرجّع kind ونسخة مستقلة من payload وnil؛ أي مدخل تاني يرجّع0 وnil وخطأ من غير تغيير المدخل أو panic.

اكتب ~FuzzCapsule~ في main_test.go: seeds تشمل قصير، توقيع غلط، kind غلط، بايتات زيادة، وصحيح بطول صفر وبـ payload. اختبر القبول والرفض وعدم تعديل المدخل. التستات لازم تمسك الأصل البايظ والتلات mutants، وتعدّي على إصلاحك ومع التوليد. مفيش Report هنا.
`, Hints: []string{
			"ابدأ بـseeds صغيرة؛ ~go test~ بيشغّلهم قبل fuzz طويل.",
			"خلّي صلاحية input جاية من العقد، مش من parser نفسه.",
			"ابدأ من أول قراءة index في parser، وبص على test القبول والرفض مع input ثابت.",
		},
		En: &LevelText{LessonRevision: 1, Steps: ch16Mission4(true), Title: "Wild Inputs", Goal: "Write a fuzz target that exposes capsule defects, then repair the parser", Files: []LabFile{{Name: "main.go", Starter: capsuleBuggy, ReadOnly: false}, {Name: "main_test.go", Starter: wildInputsTestStarter}}, Lesson: `# Try Inputs You Would Never Invent

Table tests cover chosen cases. Fuzzing mutates bytes to find panics or contract violations. The control-capsule parser is pure, so thousands of inputs need no hardware or network.

## Before you write
- Which property must hold for every input, including malformed ones?
- How will you distinguish rejecting malformed input from incorrectly rejecting everything?

## A fuzz target's shape
~~~go
func FuzzPrefix(f *testing.F) {
 f.Add("ion")
 f.Add("")
 f.Fuzz(func(t *testing.T,s string) {
  if len("x"+s)!=len(s)+1 {t.Fatal("length invariant")}
 })
}
~~~
~testing.F~ manages the corpus. ~f.Add~ supplies seeds matching the callback input type. ~f.Fuzz~ receives a function taking testing.T followed by the mutated input: string here, []byte in the mission. This example checks length rather than decoding a message.

Normal ~go test~ runs the seeds. The Lab.Fuzz stage fuzzes FuzzCapsule for five seconds. The original bug and mutants also run against seeds, so choose revealing inputs before generation begins.

## Check before indexing
~~~go
func firstByte(packet []byte) (byte,bool) {
 if len(packet)==0 {return 0,false}
 return packet[0],true
}
~~~
~byte~ aliases uint8; ~0xA7~ is a hexadecimal number. A header length does not guarantee available payload. Convert it to int and check actual length before slicing. Rejection returns an error rather than panic.

## Properties for acceptance and rejection
~~~go
before:=slices.Clone(packet)
if !bytes.Equal(packet,before) {t.Fatal("input changed")}
~~~
Import ~bytes~ to compare and ~slices~ to copy input with ~slices.Clone~, as in Median. In the real test, compare a snapshot before and after the call. Derive input validity from the contract, then check both errors and results. Merely avoiding panic permits an implementation that rejects everything. On success, modify returned payload to check it does not alias input.

## Common mistakes
~~~go
n:=int(frame[2]) // Panics with fewer than three bytes
return frame[3:3+n] // Declared length can exceed available bytes
~~~
Do not hide panics with recover inside the fuzz target: those are bugs the test must expose. Reject trailing bytes as well as truncated messages.

## Your mission
Implement ~ParseCapsule(frame []byte) (byte,[]byte,error)~. The three-byte header contains signature 0xA7, kind 1 or 2, and payload length 0 through 16. Total length must equal exactly 3 plus the declared length. Success returns kind, an independent payload copy, and nil. Every other input returns zero, nil, and an error without mutation or panic.

Write ~FuzzCapsule~ in main_test.go with seeds for short input, wrong signature, unknown kind, trailing bytes, and valid empty/nonempty payloads. Check acceptance, rejection, and unchanged input. Tests must catch the original bug and all three mutants, pass on your repair, and pass generated inputs. No Report is needed.
	`, Hints: []string{
			"Start with small seeds; ~go test~ runs them before a long fuzz session.",
			"Derive input validity from the contract, not from the parser itself.",
			"Start at the parser's first index read, then review acceptance and rejection with a fixed input.",
		}},
		Lab: &Lab{ReqFile: "main_test.go", Hidden: wildInputsHidden, Fuzz: "FuzzCapsule", Buggy: capsuleBuggy, Mutants: []Mutant{
			{i18n.T("توقيع غلط", "Wrong signature"), mutate(capsuleCorrect, "frame[0] != 0xA7", "false")},
			{i18n.T("نوع مش معروف", "Unknown kind"), mutate(capsuleCorrect, "frame[1] != 1 && frame[1] != 2", "false")},
			{i18n.T("بايتات زيادة بعد الرسالة", "Trailing bytes"), mutate(capsuleCorrect, "len(frame) != 3+n", "len(frame)<3+n")},
		}},
	})
	add(&Level{ID: "16-5-fake-clock", Chapter: 16, LessonRevision: 1, Steps: ch16Mission5(false), Kind: "lab", Race: true, Title: "ساعة مزيفة", Goal: "اختبر بوابة إطلاق بوقت محقون من غير نوم", Requires: []string{"subtest"}, Bridge: []int{15, 36},
		Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: launchGateImpl, ReadOnly: true}, {Name: "main_test.go", Starter: fakeClockTestStarter}},
		Lesson: `# حرّك الساعة بدل ما تستناها

تست بيستنى ثانية حقيقية بطيء وممكن يفشل حسب ضغط الجهاز. بوابة الإطلاق بتاخد ساعة كـ interface، فالتست يقدر ينقل الوقت فوراً ويختبر قبل الحد وعنده بالظبط.

## قبل ما تكتب
- إزاي تثبت إن رصيد جديد متاح عند الحد، مش قبله؟
- لو طلب اترفض، هتثبت إن الرصيد القديم لسه موجود إزاي؟

## الساعة المحقونة
~~~go
clock:=labkit.NewFakeClock(time.Unix(0,0))
start:=clock.Now()
clock.Advance(250*time.Millisecond)
fmt.Println(clock.Now().Sub(start)) // 250ms
~~~
استورد ~orbital/colony/labkit~ و~time~. ~labkit.Clock~ عقد Now وAfter وNewTicker، و~labkit.FakeClock~ تنفيذ قابل للتحريك. ~NewFakeClock~ يحدد بداية؛ ~time.Unix~ يبني وقت ثابت، و~Advance~ يحركه للأمام من غير sleep. ~Now~ تقرأ اللحظة، و~Sub~ ترجع المدة. في المهمة محتاجين Now بس، من غير goroutines للوقت.

## subtests بتفصل السيناريوهات
~~~go
func TestWindow(t *testing.T) {
 t.Run("before boundary",func(t *testing.T){
  clock:=labkit.NewFakeClock(time.Unix(0,0))
  clock.Advance(time.Second-time.Nanosecond)
  // افحص السلوك هنا
 })
}
~~~
~t.Run~ بيسمي سيناريو مستقل. اعمل ساعة وبوابة جديدة لكل سيناريو عشان الرصيد والوقت ميتسربوش من تست لتاني. ~time.Duration~ عدد نانوثواني؛ ~time.Second~ و~time.Millisecond~ و~time.Nanosecond~ وحدات ثابتة.

## helper يوجّه الفشل لمكان الاستدعاء
~~~go
func requirePositive(t *testing.T,n int) {
 t.Helper()
 if n<=0 {t.Fatalf("got %d",n)}
}
~~~
~t.Helper~ بتخلي رسالة الفشل تشير لسطر اللي استدعى helper. استخدم helper لفحص نتيجة Take المتكرر، من غير إخفاء سيناريو التست. المثال بيفحص رقم بس، مش بيختبر بوابة الإطلاق.

## اقرأ عقد البوابة
~~~go
gate,err:=NewLaunchGate(clock,3,time.Second)
if err!=nil {t.Fatal(err)}
ok,err:=gate.Take(2)
fmt.Println(ok,err) // true <nil>: رصيدها بدأ3
~~~
~NewLaunchGate~ بتاخد ساعة، وسعة موجبة، ومدة موجبة بين كل تصريح والتاني. ~Take(n)~ طلب موجب: true لو الرصيد كفاية، false,nil لو مش كفاية، وخطأ لو n غير موجب. البوابة بتبدأ مليانة، وكودها للقراءة بس.

كل every كاملة تضيف تصريح واحد، مع الاحتفاظ بباقي المدة وتحديد الرصيد عند capacity. عند الامتلاء مفيش رصيد زيادة يتخزن للمستقبل. الطلب المرفوض أو غير الصالح ميستهلكش. الوقت اللي يرجع للخلف ميضيفش رصيد؛ التنفيذ يحمي الطلبات المتزامنة بـ mutex.

## أخطاء شائعة
~~~go
time.Sleep(time.Second) // نوم حقيقي بدل تحريك الساعة المحقونة
clock.Advance(time.Second)
// من غير assertion التست مش بيثبت حاجة
~~~
اختبر قبل الحد بنانوثانية، ثم عنده. Mutant بيستخدم وقت حقيقي بدل الساعة المحقونة هيفشل مع Advance؛ النوم بيخبّي الفرق. التحقق بيمنع Sleep وAfter وNewTicker الحقيقيين في ملف تستك.

## المهمة
عدّل main_test.go بس. اكتب subtests وhelper بـ t.Helper باستخدام labkit.NewFakeClock وAdvance. غطّي الرصيد الأولي، رفض طلب أكبر من الرصيد من غير استهلاك، قبل refill وعندها، أجزاء مدة متراكمة، سقف السعة، ومدخلات configuration وTake غير الصالحة.

~NewLaunchGate(nil,...)~ أو capacity<=0 أو every<=0 ترجع nil,error. ~Take(n<=0)~ ترجع false,error من غير تغيير. التستات لازم تعدّي على التنفيذ الصحيح وتمسك الأربع mutants، وتشتغل من غير انتظار حقيقي. مفيش Report؛ المختبر كمان بيختبر طلبات متزامنة بـ race detector.
`, Hints: []string{
			"ساعة جديدة لكل subtest، وخلي helper يفحص الـvalue والـerror مع بعض.",
			"نظّم السيناريوهات حوالين الوقت والرصيد، وسيب كل واحد يبدأ بحالة مستقلة.",
			"ابدأ من signatures ~NewLaunchGate~ و~Take~ قبل استعمال ~FakeClock.Advance~ في الـtest.",
		},
		En: &LevelText{LessonRevision: 1, Steps: ch16Mission5(true), Title: "Fake Clock", Goal: "Test a launch gate using injected time without sleeping", Files: []LabFile{{Name: "main.go", Starter: launchGateImpl, ReadOnly: true}, {Name: "main_test.go", Starter: fakeClockTestStarter}}, Lesson: `# Move Time Instead of Waiting

A test that waits a real second is slow and vulnerable to machine load. The launch gate accepts a clock interface, letting tests move time instantly and inspect exact boundaries.

## Before you write
- How will you prove new credit is available exactly at the boundary, never before?
- After a rejected request, how will you prove the old credit is intact?

## An injected clock
~~~go
clock:=labkit.NewFakeClock(time.Unix(0,0))
start:=clock.Now()
clock.Advance(250*time.Millisecond)
fmt.Println(clock.Now().Sub(start)) // 250ms
~~~
Import ~orbital/colony/labkit~ and ~time~. ~labkit.Clock~ specifies Now, After, and NewTicker; ~labkit.FakeClock~ is a controllable implementation. ~NewFakeClock~ sets its starting instant; ~time.Unix~ creates a fixed time; ~Advance~ moves forward without sleep. ~Now~ reads time and ~Sub~ returns elapsed duration. This mission only needs Now, without timer goroutines.

## Subtests separate scenarios
~~~go
func TestWindow(t *testing.T) {
 t.Run("before boundary",func(t *testing.T){
  clock:=labkit.NewFakeClock(time.Unix(0,0))
  clock.Advance(time.Second-time.Nanosecond)
  // Check behavior here.
 })
}
~~~
~t.Run~ names an independent scenario. Create a fresh clock and gate per scenario so time and credit cannot leak between tests. ~time.Duration~ counts nanoseconds; ~time.Second~, ~time.Millisecond~, and ~time.Nanosecond~ are duration constants.

## A helper locates failure at its caller
~~~go
func requirePositive(t *testing.T,n int) {
 t.Helper()
 if n<=0 {t.Fatalf("got %d",n)}
}
~~~
~t.Helper~ makes failures point to the helper's caller. Use a helper for repeated Take assertions without hiding the scenario. This example checks a number rather than testing the launch gate.

## Read the gate contract
~~~go
gate,err:=NewLaunchGate(clock,3,time.Second)
if err!=nil {t.Fatal(err)}
ok,err:=gate.Take(2)
fmt.Println(ok,err) // true <nil>: starts with three credits
~~~
~NewLaunchGate~ accepts a clock, positive capacity, and positive interval between credits. ~Take(n)~ requires a positive request: true when affordable, false,nil when unaffordable, and an error when n is nonpositive. The gate starts full; its implementation is read-only.

Each complete every interval earns one credit, retaining fractional time and clamping balance at capacity. Full gates cannot bank excess credit for future spending. Rejected and invalid requests consume nothing. Backward time adds no credit; the implementation protects concurrent requests with a mutex.

## Common mistakes
~~~go
time.Sleep(time.Second) // Real waiting instead of advancing injected time
clock.Advance(time.Second)
// Without an assertion this proves nothing.
~~~
Check one nanosecond before refill, then exactly at it. A mutant using real time fails when fake time advances; sleeping can conceal that difference. Hidden checks reject real Sleep, After, and NewTicker calls in your test file.

## Your mission
Edit only main_test.go. Write subtests and a t.Helper assertion using labkit.NewFakeClock and Advance. Cover initial credit, oversized rejection without consumption, just-before/exact refill, accumulated fractional intervals, capacity ceiling, and invalid configuration and Take inputs.

~NewLaunchGate(nil,...)~, capacity<=0, or every<=0 return nil,error. ~Take(n<=0)~ returns false,error without mutation. Tests must pass the correct implementation, kill all four mutants, and perform no real waiting. No Report is needed; the lab also checks concurrent requests with the race detector.
	`, Hints: []string{
			"Give each subtest a fresh clock, and let a helper check value and error together.",
			"Organize scenarios around time and credit, with each one starting from isolated state.",
			"Start with the ~NewLaunchGate~ and ~Take~ signatures before using ~FakeClock.Advance~ in the test.",
		}},
		Lab: &Lab{ReqFile: "main_test.go", Hidden: fakeClockHidden, Mutants: []Mutant{
			{i18n.T("لحظة refill بالظبط", "Exact refill instant"), mutate(launchGateImpl, "now.Sub(g.last) / g.every", "(now.Sub(g.last)-1)/g.every")},
			{i18n.T("تجاهل الساعة المحقونة", "Ignored injected clock"), mutate(launchGateImpl, "now := g.clock.Now()", "now:=time.Now()")},
			{i18n.T("طلب صفر مش صالح", "Zero requests are invalid"), mutate(launchGateImpl, "if n <= 0", "if n<0")},
			{i18n.T("رفض الطلب لازم يحافظ على الرصيد", "Rejection must preserve credit"), mutate(launchGateImpl, "if n > g.credits {\n\t\treturn false, nil\n\t}", "if n>g.credits {g.credits=0;return false,nil}")},
		}},
	})
}

const wildInputsTestStarter = `package main

import "testing"

func TestCapsule(t *testing.T) {
	kind, payload, err := ParseCapsule([]byte{0xA7, 1, 1, 42})
	if err != nil || kind != 1 || len(payload) != 1 || payload[0] != 42 {
		t.Fatal("valid capsule")
	}
}

// Add FuzzCapsule with seeds and a contract assertion.
`
const wildInputsHidden = `package main

import (
	"bytes"
	"io"
	"orbital/colony/levels"
	"testing"
)

func TestHiddenCapsuleBoundsAndOwnership(t *testing.T) {
	world := levels.Start("16-4-wild-inputs", io.Discard, 164)
	r := world.Rand
	frames := [][]byte{nil, {0xA7}, {0xA7, 1}}
	for n := 0; n <= 18; n++ {
		for _, kind := range []byte{0, 1, 2, 3, 255} {
			frame := make([]byte, n+3)
			for i := range frame {
				frame[i] = byte(world.Data.SensorReadings[(n+i)%256])
			}
			frame[0] = 0xA7
			frame[1] = kind
			frame[2] = byte(n)
			frames = append(frames, frame, append(append([]byte(nil), frame...), 0), frame[:len(frame)-1])
		}
	}
	for i := 0; i < 300; i++ {
		frame := make([]byte, r.Intn(40))
		r.Read(frame)
		frames = append(frames, frame)
	}
	for _, frame := range frames {
		before := append([]byte(nil), frame...)
		valid := len(frame) >= 3 && frame[0] == 0xA7 && (frame[1] == 1 || frame[1] == 2) && frame[2] <= 16 && len(frame) == 3+int(frame[2])
		kind, payload, err := ParseCapsule(frame)
		if !bytes.Equal(frame, before) {
			t.Fatal("mutated frame")
		}
		if !valid {
			if err == nil || kind != 0 || payload != nil {
				t.Fatalf("accepted %x", frame)
			}
			continue
		}
		if err != nil || kind != frame[1] || !bytes.Equal(payload, frame[3:]) {
			t.Fatalf("decode %x", frame)
		}
		if len(payload) > 0 {
			payload[0] ^= 255
			if !bytes.Equal(frame, before) {
				t.Fatal("aliased input")
			}
		}
	}
}
`

const capsuleCorrect = `package main

import ("errors"; "slices")

func ParseCapsule(frame []byte) (byte, []byte, error) {
	if len(frame) < 3 {
		return 0, nil, errors.New("short capsule")
	}
	if frame[0] != 0xA7 {
		return 0, nil, errors.New("signature")
	}
	if frame[1] != 1 && frame[1] != 2 {
		return 0, nil, errors.New("kind")
	}
	n := int(frame[2])
	if n > 16 || len(frame) != 3+n {
		return 0, nil, errors.New("length")
	}
	return frame[1], slices.Clone(frame[3:]), nil
}
func main() {}
`
const capsuleBuggy = `package main

import "errors"

func ParseCapsule(frame []byte) (byte, []byte, error) {

	if frame[0] != 0xA7 {
		return 0, nil, errors.New("signature")
	}
	if frame[1] != 1 && frame[1] != 2 {
		return 0, nil, errors.New("kind")
	}
	n := int(frame[2])
	if n > 16 || len(frame) < 3+n {
		return 0, nil, errors.New("length")
	}
	return frame[1], append([]byte(nil), frame[3:]...), nil
}
func main() {}
`

const fakeClockTestStarter = `package main

import (
	"orbital/colony/labkit"
	"testing"
	"time"
)

func TestLaunchGate(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	gate, err := NewLaunchGate(clock, 2, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ok, err := gate.Take(1)
	if !ok || err != nil {
		t.Fatal("initial credit")
	}
}
`
const fakeClockHidden = `package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestHiddenInjectedTimeAndBoundaries(t *testing.T) {
	world := levels.Start("16-5-fake-clock", io.Discard, 165)
	for iteration := 0; iteration < 40; iteration++ {
		cap := 1 + world.Data.SensorReadings[iteration]%8
		every := time.Duration(1+world.Data.SensorReadings[iteration+40]) * time.Millisecond
		clock := labkit.NewFakeClock(time.Unix(100, 0))
		gate, err := NewLaunchGate(clock, cap, every)
		if err != nil {
			t.Fatal(err)
		}
		if ok, err := gate.Take(cap + 1); err != nil || ok {
			t.Fatal("oversized request")
		}
		if ok, err := gate.Take(cap); err != nil || !ok {
			t.Fatal("rejection consumed credit")
		}
		clock.Advance(every - time.Nanosecond)
		if ok, _ := gate.Take(1); ok {
			t.Fatal("early refill")
		}
		clock.Advance(time.Nanosecond)
		if ok, _ := gate.Take(1); !ok {
			t.Fatal("exact boundary")
		}
		for i := 0; i < 10; i++ {
			clock.Advance(every / 10)
		}
		// Account for integer division by finishing exactly at the next boundary.
		clock.Advance(every - (every/10)*10)
		if ok, _ := gate.Take(1); !ok {
			t.Fatal("fractional carry")
		}
		clock.Advance(100 * every)
		if ok, _ := gate.Take(cap + 1); ok {
			t.Fatal("capacity ceiling")
		}
		if ok, _ := gate.Take(cap); !ok {
			t.Fatal("full refill")
		}
		for _, n := range []int{-1, 0} {
			if ok, err := gate.Take(n); ok || err == nil {
				t.Fatal("invalid input")
			}
		}
		if ok, _ := gate.Take(1); ok {
			t.Fatal("invalid input changed credit")
		}
	}
}
func TestHiddenConcurrentAdmissions(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	gate, err := NewLaunchGate(clock, 7, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	accepted := 0
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok, err := gate.Take(1)
			if err != nil {
				t.Error(err)
			}
			if ok {
				mu.Lock()
				accepted++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if accepted != 7 {
		t.Fatalf("accepted %d want 7", accepted)
	}
	clock.Advance(time.Second)
	if ok, _ := gate.Take(1); !ok {
		t.Fatal("refill after concurrency")
	}
}
func TestHiddenClockTestDiscipline(t *testing.T) {
	src, err := os.ReadFile("main_test.go")
	if err != nil {
		t.Fatal(err)
	}
	f, err := parser.ParseFile(token.NewFileSet(), "main_test.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	imports := map[string]string{}
	for _, spec := range f.Imports {
		path, _ := strconv.Unquote(spec.Path.Value)
		name := path
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				name = path[i+1:]
				break
			}
		}
		if spec.Name != nil {
			name = spec.Name.Name
		}
		imports[name] = path
	}
	helper, fake, advance := false, false, false
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		if sel.Sel.Name == "Helper" {
			helper = true
		}
		if sel.Sel.Name == "Advance" {
			advance = true
		}
		if pkg, ok := sel.X.(*ast.Ident); ok {
			if imports[pkg.Name] == "time" && (sel.Sel.Name == "Sleep" || sel.Sel.Name == "NewTicker" || sel.Sel.Name == "After") {
				t.Error("tests must advance fake time rather than wait for real time")
			}
			if imports[pkg.Name] == "orbital/colony/labkit" && sel.Sel.Name == "NewFakeClock" {
				fake = true
			}
		}
		return true
	})
	if !helper || !fake || !advance {
		t.Fatal("use a test helper, NewFakeClock, and Advance")
	}
}
`

const launchGateImpl = `package main

import (
	"errors"
	"orbital/colony/labkit"
	"sync"
	"time"
)

type LaunchGate struct {
	mu                sync.Mutex
	clock             labkit.Clock
	capacity, credits int
	every             time.Duration
	last              time.Time
}

func NewLaunchGate(clock labkit.Clock, capacity int, every time.Duration) (*LaunchGate, error) {
	if clock == nil || capacity <= 0 || every <= 0 {
		return nil, errors.New("invalid launch gate")
	}
	return &LaunchGate{clock: clock, capacity: capacity, credits: capacity, every: every, last: clock.Now()}, nil
}
func (g *LaunchGate) Take(n int) (bool, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if n <= 0 {
		return false, errors.New("positive request required")
	}
	now := g.clock.Now()
	if now.After(g.last) {
		elapsed := now.Sub(g.last) / g.every
		if elapsed > 0 {
			earned := g.capacity - g.credits
			if elapsed < time.Duration(earned) {
				earned = int(elapsed)
			}
			g.credits += earned
			g.last = g.last.Add(elapsed * g.every)
		}
	}
	if n > g.credits {
		return false, nil
	}
	g.credits -= n
	return true, nil
}
func main() {}
`

// seededLabRoom supplies a repeatable stimulus corpus to the new lab suites.
// Hidden tests create their own world because RunLab does not invoke Setup.
func seededLabRoom(w *sim.World) {
	labRoom(w)
	for i := 0; i < 256; i++ {
		w.Data.SensorReadings = append(w.Data.SensorReadings, w.Rand.Intn(1000))
	}
}
