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
			Hints:   []string{"Put cases in a slice: ~cases := []struct{ d, c, want int }{ {3, 12, 8}, ... }~, then loop with ~for _, tc := range cases~.", "Think about **boundaries**: 50 and 51, distance 0 with cargo, 6 kg, and a very large load such as 200.", "Each surviving mutant means a case is missing. Read its description in the Console."},
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
		ID: "16-1-mutants", Chapter: 16, Title: "اصطاد الـ Mutants",
		Goal:     "اكتب تستات لـ Fuel تمسك كل نسخة متبوّظة منها (6 mutants)",
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
			"حط الحالات في slice: ~cases := []struct{ d, c, want int }{ {3, 12, 8}, ... }~، ولف عليها بـ ~for _, tc := range cases~.",
			"فكّر في **الحدود**: 50 و51، والمسافة 0 مع شحنة، و6 كجم، وشحنة كبيرة أوي زي 200.",
			"كل mutant بيعدّي يبقى فيه حالة مش مغطّية. اقرا الوصف بتاعه في الـ Console.",
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
			Hints:   []string{"Trace main.go mentally: after ~Median~, ~readings~ is sorted! That is a bug: the function changed data belonging to its caller.", "Revealing cases: even length ~[4 1 3 2]~ should give 2; an empty slice should give 0 (it currently panics); input must remain unchanged (compare a copy with ~slices.Equal~).", "The fix: ~s := slices.Clone(xs)~, then ~slices.Sort(s)~. Handle ~len(s) == 0~ and even lengths."},
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
		ID: "16-2-bug-hunt", Chapter: 16, Title: "اصطاد الـ Bug",
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
			"شغّل main.go في دماغك: ~readings~ بعد ~Median~ بقت مترتبة! ده bug: الدالة غيّرت حاجة مش بتاعتها.",
			"حالات تكشف: عدد زوجي ~[4 1 3 2]~ المفروض 2، وslice فاضية المفروض 0 (دلوقتي بتعمل panic)، وإن الـ input ميتغيرش (قارن بـ ~slices.Equal~ مع نسخة).",
			"الإصلاح: ~s := slices.Clone(xs)~ ثم ~slices.Sort(s)~، واتعامل مع ~len(s) == 0~ والعدد الزوجي.",
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
			Hints:   []string{"The benchmark: ~func BenchmarkCountPairs(b *testing.B) { xs := ...; for b.Loop() { CountPairs(xs, 1000) } }~.", "Idea 1: one pass with a map. For each x, new pairs = ~seen[target-x]~, then ~seen[x]++~. This is O(n).", "Idea 2: sort a copy and use two pointers, one at each end. This is O(n log n). Watch out for duplicates! Try both and measure which is actually faster."},
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
		ID: "16-3-benchmark", Chapter: 16, Title: "أسرع 100 مرة",
		Goal:     "CountPairs بطيئة جداً (O(n²)). خليها أسرع من غير ما تبوّظها، واكتب benchmark بنفسك",
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
			"الـ benchmark: ~func BenchmarkCountPairs(b *testing.B) { xs := ...; for b.Loop() { CountPairs(xs, 1000) } }~.",
			"فكرة 1: لفّة واحدة بـ map: لكل x، الأزواج الجديدة = ~seen[target-x]~، وبعدين ~seen[x]++~. ده O(n).",
			"فكرة 2: رتّب نسخة، وامشي بمؤشرين واحد من الأول وواحد من الآخر (two pointers). ده O(n log n). خلّي بالك من الأرقام المتكررة! جرّب الاتنين وقيس: مين أسرع فعلاً؟",
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
}
