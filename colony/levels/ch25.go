package levels

import "orbital/colony/i18n"

func init() {
	add(&Level{ID: "25-1-fuel-account", Chapter: 25, Kind: "lab", Race: true, LessonRevision: 1, Steps: ch25FuelArabic,
		Title: "حساب الوقود", Goal: "احمي رصيد خزان الوقود وحدوده مع التعبئة والسحب المتوازي",
		Requires: []string{"mutex"}, Bridge: []int{38}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch251Starter}},
		Lesson: `# رصيد الخزان لازم يفضل صحيح

## قبل ما تكتب
- إيه البيانات اللي لازم تتقري وتتغير تحت نفس القفل؟
- هتثبت عدد عمليات السحب الناجحة إزاي من غير اختيار درون بعينه؟

## الشرط الثابت تحت التزامن
الـ invariant شرط لازم يفضل صحيح بعد كل عملية مكتملة. خزان المستعمرة لازم يفضل بين الصفر والسعة. العملية المرفوضة متغيرش الرصيد.

قفل القراءة لوحدها والكتابة لوحدها يمنع data race لكنه يسمح بخطأ منطقي: درونين يوافقوا على نفس الرصيد. القرار كله لازم يبقى مقطع حرج واحد.
~~~go
mu.Lock()
if seats > 0 { seats--; accepted = true }
mu.Unlock()
~~~
المثال بيحجز كرسي، مش وقود. حتى قراءة الرصيد لازم تستخدم نفس القفل عشان تتزامن مع الكتابة.

## ملكية القفل
استخدم حقل ~sync.Mutex~ وطرق بمستقبل pointer. ممنوع نسخ كائن فيه قفل بعد استخدامه.
~~~go
type SeatPool struct { mu sync.Mutex; seats int }
func (p *SeatPool) Count() int {
 p.mu.Lock()
 defer p.mu.Unlock()
 return p.seats
}
~~~

## حدود الأعداد
استخدم ~int64~ لوحدات صحيحة. افحص المساحة المتاحة بالطرح قبل الجمع، عشان الجمع ممكن يلف لعدد سالب.
~~~go
if incoming > capacity-current { return errors.New("full") }
current += incoming
~~~

## واجهة الخزان
~~~go
account, err := NewFuelAccount(80, 20)
if err != nil { return err }
err = account.Load(7)
err = account.Burn(3)
fmt.Println(account.Remaining(), err)
~~~
~NewFuelAccount~ ترجع pointer ثابت للخزان. ~Load~ تعبئة و ~Burn~ سحب و ~Remaining~ لقطة متزامنة؛ متعرضش الحقول الداخلية.

## أخطاء شائعة
~~~go
if account.Remaining() >= n {
 mu.Lock()
 fuel -= n // الرصيد ممكن يكون اتسحب بعد الفحص
 mu.Unlock()
}
~~~

## المهمة
نفّذ ~NewFuelAccount(capacity, initial int64) (*FuelAccount,error)~ وطرق ~Load(n int64) error~ و ~Burn(n int64) error~ و ~Remaining() int64~. ارفض سعة غير موجبة أو بداية خارج ~0..capacity~ بـ ~nil,error~. التعبئة والسحب يقبلوا أعداد موجبة بس. التعبئة فوق السعة والسحب فوق الرصيد أخطاء من غير تغيير. السحب بكل الرصيد والتعبئة لحد السعة نجاح.

كل قرار وفحص وتحديث، وكل قراءة، تحت نفس ~sync.Mutex~. المتزامن يعادل ترتيب متسلسل للعمليات؛ مفيش ضمان لدرون فائز. الحسابات مستقلة وممنوع نسخها بعد الاستخدام. مفيش وقت أو شبكة أو ~Report~، و ~main~ فاضية.

النجمة التالتة: نفس العقد جنب أكبر ~int64~ من غير overflow.`,
		Hints: []string{"المطلوب الحفاظ على invariant للحساب.", "خلّي فحص القيمة والتغيير متزامنين.", "راجع طرق ~Load~ و~Burn~ في التوقيع قبل حماية القراءات."},
		En: &LevelText{LessonRevision: 1, Steps: ch25FuelEnglish, Title: "Fuel Account", Goal: "Protect tank balance and capacity during concurrent loading and burning", Files: []LabFile{{Name: "main.go", Starter: ch251Starter}},
			Lesson: `# Keep the Tank Invariant True

## Before you write
- Which reads and writes must share one critical section?
- How will you count successful burns without choosing a winning drone?

## An invariant under concurrency
An invariant must hold after every completed operation. Colony fuel stays between zero and tank capacity. Failed operations leave it unchanged.

Separately locked reads and writes can eliminate data races while allowing a logical race: two drones approve spending the same fuel. Protect the complete decision.
~~~go
mu.Lock()
if seats > 0 { seats--; accepted = true }
mu.Unlock()
~~~
This example reserves a seat rather than fuel. Balance reads also need the same lock to synchronize with writes.

## Lock ownership
Store a ~sync.Mutex~ beside state and use pointer receivers. Never copy a lock-containing object after use.
~~~go
type SeatPool struct { mu sync.Mutex; seats int }
func (p *SeatPool) Count() int {
 p.mu.Lock()
 defer p.mu.Unlock()
 return p.seats
}
~~~

## Integer boundaries
Use ~int64~ whole units. Check remaining room by subtraction before addition, which could wrap negative.
~~~go
if incoming > capacity-current { return errors.New("full") }
current += incoming
~~~

## Tank API
~~~go
account, err := NewFuelAccount(80, 20)
if err != nil { return err }
err = account.Load(7)
err = account.Burn(3)
fmt.Println(account.Remaining(), err)
~~~
~NewFuelAccount~ returns a stable pointer. ~Load~ adds fuel, ~Burn~ consumes it, and ~Remaining~ reads a synchronized snapshot. Keep fields private.

## Common mistakes
~~~go
if account.Remaining() >= n {
 mu.Lock()
 fuel -= n // Another drone may have spent it after the check
 mu.Unlock()
}
~~~

## Your mission
Implement ~NewFuelAccount(capacity, initial int64) (*FuelAccount,error)~, ~Load(n int64) error~, ~Burn(n int64) error~, and ~Remaining() int64~. Reject nonpositive capacity or initial fuel outside ~0..capacity~ with ~nil,error~. Loading and burning require positive amounts. Exceeding capacity or available fuel fails without mutation. Exact emptying and filling succeed.

Protect every complete check/update and every read with one ~sync.Mutex~. Concurrent behavior must match a serial order; no winner is promised. Accounts are independent and must not be copied after use. No time, network, or ~Report~; leave ~main~ empty.

Third star: preserve the same contract near maximum ~int64~ without overflow.`,
			Hints: []string{"The account must preserve its balance invariant.", "Keep checking and changing state synchronized.", "Review ~Load~ and ~Burn~ in the signature before protecting reads."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch251Hidden},
	})
}

const ch251Starter = `package main

import "errors"

type FuelAccount struct{}

func NewFuelAccount(capacity, initial int64) (*FuelAccount, error) { return nil, errors.New("TODO") }
func (a *FuelAccount) Load(n int64) error                          { return errors.New("TODO") }
func (a *FuelAccount) Burn(n int64) error                          { return errors.New("TODO") }
func (a *FuelAccount) Remaining() int64                            { return 0 }
func main()                                                        {}
`

const ch251Hidden = `package main

import (
	"math"
	"math/rand"
	"sync"
	"testing"
)

// Boundaries: invalid construction/amounts, empty/full/exact, random sequences, concurrent spending.
func TestHiddenFuelSerialModel(t *testing.T) {
	for _, c := range [][2]int64{{0, 0}, {-1, 0}, {10, -1}, {10, 11}} {
		if a, e := NewFuelAccount(c[0], c[1]); a != nil || e == nil {
			t.Fatal("configuration", c)
		}
	}
	rng := rand.New(rand.NewSource(251))
	for trial := 0; trial < 30; trial++ {
		capacity := int64(1 + rng.Intn(500))
		balance := int64(rng.Intn(int(capacity + 1)))
		a, e := NewFuelAccount(capacity, balance)
		if e != nil || a == nil {
			t.Fatal("construct", e)
		}
		for step := 0; step < 100; step++ {
			n := int64(rng.Intn(700) - 20)
			load := rng.Intn(2) == 0
			valid := n > 0 && n <= balance
			if load {
				valid = n > 0 && n <= capacity-balance
				e = a.Load(n)
			} else {
				e = a.Burn(n)
			}
			if (e == nil) != valid {
				t.Fatal("acceptance", capacity, balance, n, load, e)
			}
			if valid {
				if load {
					balance += n
				} else {
					balance -= n
				}
			}
			if a.Remaining() != balance {
				t.Fatal("mutation", a.Remaining(), balance)
			}
		}
	}
	a, _ := NewFuelAccount(9, 0)
	for _, n := range []int64{-1, 0, 10} {
		if a.Load(n) == nil || a.Remaining() != 0 {
			t.Fatal("invalid load")
		}
	}
	if a.Load(9) != nil || a.Burn(10) == nil || a.Remaining() != 9 || a.Burn(9) != nil || a.Remaining() != 0 {
		t.Fatal("exact boundaries")
	}
}
func TestHiddenFuelConcurrentInvariant(t *testing.T) {
	rng := rand.New(rand.NewSource(25101))
	for trial := 0; trial < 8; trial++ {
		units := int64(1 + rng.Intn(9))
		affordable := 1 + rng.Intn(20)
		a, _ := NewFuelAccount(10000, int64(affordable)*units)
		start := make(chan struct{})
		outcomes := make(chan bool, 64)
		var wg sync.WaitGroup
		for i := 0; i < 64; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				e := a.Burn(units)
				if a.Remaining() < 0 {
					t.Error("negative")
				}
				outcomes <- e == nil
			}()
		}
		close(start)
		wg.Wait()
		close(outcomes)
		successes := 0
		for ok := range outcomes {
			if ok {
				successes++
			}
		}
		if successes != affordable || a.Remaining() != 0 {
			t.Fatal("overspent", successes, a.Remaining())
		}
		increments := make([]int64, 64)
		want := int64(0)
		for i := range increments {
			increments[i] = int64(1 + rng.Intn(20))
			want += increments[i]
		}
		start = make(chan struct{})
		for _, n := range increments {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				if e := a.Load(n); e != nil {
					t.Error(e)
				}
				_ = a.Remaining()
			}()
		}
		close(start)
		wg.Wait()
		if a.Remaining() != want {
			t.Fatal("lost deposits")
		}
		// Mixed operations reconcile all successful deltas, regardless of winners.
		deltas := make(chan int64, 128)
		start = make(chan struct{})
		for i := 0; i < 128; i++ {
			n := int64(1 + rng.Intn(30))
			load := i%2 == 0
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				if load {
					if a.Load(n) == nil {
						deltas <- n
					}
				} else if a.Burn(n) == nil {
					deltas <- -n
				}
				if a.Remaining() < 0 {
					t.Error("negative snapshot")
				}
			}()
		}
		close(start)
		wg.Wait()
		close(deltas)
		for n := range deltas {
			want += n
		}
		other, _ := NewFuelAccount(2, 1)
		if a.Remaining() != want || other.Remaining() != 1 {
			t.Fatal("mixed reconciliation/isolation")
		}
	}
}
func TestHiddenBonusFuelInt64Capacity(t *testing.T) {
	a, e := NewFuelAccount(math.MaxInt64, math.MaxInt64-1)
	if e != nil || a.Load(2) == nil || a.Remaining() != math.MaxInt64-1 {
		t.Fatal("overflow mutation")
	}
	if a.Load(1) != nil || a.Load(1) == nil || a.Burn(math.MaxInt64) != nil || a.Remaining() != 0 {
		t.Fatal("maximum boundary")
	}
}
`

func init() {
	add(&Level{ID: "25-2-race-report", Chapter: 25, Kind: "lab", Race: true, LessonRevision: 1, Steps: ch25RaceArabic, Title: "تقرير السباق", Goal: "اثبت فقدان تحديث بتست حتمي، اشرح السبب، وصلّح جمع تقارير المسح",
		Requires: []string{"mutex", "subtest"}, Bridge: []int{40}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch252Buggy}, {Name: "main_test.go", Starter: ch252TestStarter}},
		Lesson: `# دليل على سباق تقارير المسح

## قبل ما تكتب
- هل النتيجة الصح مرة واحدة دليل إن مفيش سباق؟
- إزاي تفرض قراءة الرصيد القديم قبل أي كتابة من غير نوم؟

## data race و happens-before
الـ data race وصولين متعارضين لنفس الذاكرة من غير ترتيب بالمزامنة، وواحد فيهم كتابة. ~total += n~ قراءة وحساب وكتابة، مش خطوة واحدة.

~happens-before~ علاقة ترتيب: إرسال قناة يسبق اكتمال الاستقبال المقابل، وفتح قفل يسبق القفل التالي، و ~Done~ تسبق رجوع ~Wait~ المرتبطة بيها. بدء goroutines لوحده مش ترتيب للذاكرة.
~~~go
ready := make(chan struct{})
go func() { answer = 9; close(ready) }()
<-ready
fmt.Println(answer) // الكتابة حصلت قبل القراءة
~~~

## نقطة التنسيق المحقونة
~MergeScans~ تجمع حمولات المسح. لكل index فيه callback اسمها ~checkpoint~ قبل كتابة المجموع، وخارج أي قفل. الاختبار يقدر يوقف العامل عندها، والتنفيذ الصحيح يناديها قبل قراية المجموع كمان.
~~~go
arrived := make(chan int, 2)
release := []chan struct{}{make(chan struct{}), make(chan struct{})}
checkpoint := func(i int) { arrived <- i; <-release[i] }
~~~
استنى وصول كل العمال، وبعدها افتح عامل واحد واستنى رجوع callback بتاعته؟ ده مش كفاية: لازم تستنى **اكتمال التحديث**. عشان كده ~checkpoint~ ترجع دالة إقرار، والتنفيذ يناديها مرة بعد التحديث وخارج القفل.
~~~go
checkpoint := func(i int) func() {
 arrived <- i
 <-release[i]
 return func() { finished <- i }
}
~~~
ده seam للاختبار، مش عداد جاهز. افتح العمال واحد واحد واستنى إقرار كل واحد. الأصل البايظ قرأ قبل الحاجز، فيكتب نسخة قديمة؛ الإصلاح يقرأ ويجمع تحت قفل بعد الحاجز. إشارات الاختبار بترتب الوصولات، فالدليل هنا فقدان تحديث حتمي. التشغيل الحر للأصل ممكن يطلع تقرير detector حقيقي؛ مش مطلوب الاحتفاظ بتست بيسبّب race في النسخة النهائية.

## تستات قبل الإصلاح
استخدم ~t.Run~ لسيناريو الدليل، ولو فشل ارجع قبل اختبار التزامن الحر. اجمع عدد زيارات كل index تحت ~sync.Mutex~ في التست. بعد الإصلاح شغّل ضغط متزامن وقراءات فارغة وسالبة.
~~~go
if !t.Run("controlled interleaving", func(t *testing.T) { /* assertions */ }) { return }
~~~
الـ detector بيفحص المسارات اللي اتنفذت بس. جولة نظيفة مش إثبات لكل المسارات، ومش مقياس سرعة بسبب تكلفة instrumentation.

## أخطاء شائعة
~~~go
old := total
mu.Lock()
total = old + n // القفل جه بعد القراءة؛ التحديث ممكن يضيع
mu.Unlock()
~~~

## المهمة
اكتب تستات الأول وبعدين صلّح ~MergeScans(scans []int, checkpoint func(int) func()) (int,error)~. ارفض callback ~nil~ أو أكتر من 64 عنصر قبل أي callback بـ ~0,error~. غير كده الفاضي يرجع ~0,nil~. كل عنصر، حتى الصفر والسالب، له عامل ونداء ~checkpoint(index)~ مرة قبل التحديث، ثم نداء دالة الإقرار الناتجة مرة بعد التحديث. الدوال المحقونة بترجع، وممكن تستنى إشارات، ودالة الإقرار غير ~nil~. الاتنين خارج القفل. استنى كل العمال، متغيرش المدخل، وارجع المجموع الصحيح؛ المدخلات لا تسبب overflow.

في ~main_test.go~ استخدم ~t.Run~ و ~sync.Mutex~ لحماية سجل الزيارات. التست الحتمي يرفض الأصل قبل اختبار الضغط الحر، ويمسك 3 mutants: تجاهل السالب، نسيان آخر عنصر، وجمع مرتين. في تعليق في ~main.go~ اشرح ~data race~ و ~happens-before~ وليه القفل يشمل القراءة والتحديث. مفيش وقت حقيقي أو ~Report~.`,
		Hints: []string{"اجعل فقد التحديث مفهومًا قبل تعديل السلوك.", "فكّر في حدود العملية المركبة المحمية.", "راجع ~checkpoint~ في العقد لتحديد نقطة التزامن."},
		En: &LevelText{LessonRevision: 1, Steps: ch25RaceEnglish, Title: "Race Report", Goal: "Prove a lost update deterministically, explain it, and repair concurrent scan aggregation", Files: []LabFile{{Name: "main.go", Starter: ch252Buggy}, {Name: "main_test.go", Starter: ch252TestStarter}},
			Lesson: `# Evidence of a Scan Report Race

## Before you write
- Does one correct total prove race freedom?
- How can you force all stale reads before any writes without sleeping?

## Data race and happens-before
A data race consists of conflicting unordered memory accesses, at least one writing. ~total += n~ reads, computes, and writes; it is not indivisible.

~happens-before~ describes synchronization ordering: a send precedes completion of its receive, unlocking precedes a later lock, and ~Done~ precedes the corresponding ~Wait~ return. Starting goroutines alone does not order shared memory.
~~~go
ready := make(chan struct{})
go func() { answer = 9; close(ready) }()
<-ready
fmt.Println(answer) // Write precedes read
~~~

## An injected checkpoint
~MergeScans~ aggregates scan payloads. Each index calls ~checkpoint~ before updating, outside any lock. Tests can pause workers there; correct implementations also read the total after this checkpoint.
~~~go
arrived := make(chan int, 2)
release := []chan struct{}{make(chan struct{}), make(chan struct{})}
checkpoint := func(i int) { arrived <- i; <-release[i] }
~~~
Wait for every worker to arrive, then release workers individually. Callback return alone does not prove completion of the update. Therefore the mission checkpoint returns an acknowledgement function, called once after updating, outside the lock.
~~~go
checkpoint := func(i int) func() {
 arrived <- i
 <-release[i]
 return func() { finished <- i }
}
~~~
This is a testing seam. Release one worker, wait for its acknowledgement, then release the next. The buggy version reads before the barrier and writes stale totals; the repair reads and adds under one lock after it. Test signals order accesses, demonstrating a deterministic lost update. Freely running the original can produce an actual detector report; final tests need not retain a deliberately racing execution.

## Tests first
Use ~t.Run~ for the controlled proof and return if it fails before free contention. Protect per-index visit counts with a ~sync.Mutex~ in tests. After repair, exercise free contention, empty input, and signed values.
~~~go
if !t.Run("controlled interleaving", func(t *testing.T) { /* assertions */ }) { return }
~~~
The detector observes executed paths only. A clean run does not prove all paths safe; instrumentation overhead also makes it unsuitable for universal performance claims.

## Common mistakes
~~~go
old := total
mu.Lock()
total = old + n // The read happened before locking; updates can be lost
mu.Unlock()
~~~

## Your mission
Write tests first, then repair ~MergeScans(scans []int, checkpoint func(int) func()) (int,error)~. Reject a ~nil~ callback or more than 64 elements before callbacks with ~0,error~. Otherwise empty input returns ~0,nil~. Each element, including zero and negative values, gets a worker and exactly one ~checkpoint(index)~ before updating, followed by exactly one call of its returned acknowledgement after updating. Injected functions terminate, may await signals, and return non-~nil~ acknowledgements. Both callbacks run outside the lock. Join workers, preserve input, and return the exact sum; inputs do not overflow.

In ~main_test.go~, use ~t.Run~ and ~sync.Mutex~ protecting the visit log. The controlled test rejects the original before free contention and catches three mutants: skipping negatives, omitting the last item, and doubling additions. Add a ~main.go~ comment explaining ~data race~, ~happens-before~, and why the lock includes reading and updating. No real time or ~Report~.`,
			Hints: []string{"Understand the lost update before changing behavior.", "Think about the boundary of the protected compound operation.", "Review ~checkpoint~ in the contract to locate synchronization."}},
		Lab: &Lab{ReqFile: "main_test.go", Buggy: ch252Buggy, Hidden: ch252Hidden, Mutants: []Mutant{
			{i18n.T("تجاهل المسح السالب", "Skip negative scans"), mutate(ch252Correct, "total += n", "if n > 0 { total += n }")},
			{i18n.T("نسيان آخر مسح", "Omit final scan"), mutate(ch252Correct, "total += n", "if i < len(scans)-1 { total += n }")},
			{i18n.T("مسح متحسب مرتين", "Double counted scan"), mutate(ch252Correct, "total += n", "total += 2*n")},
		}},
	})
}

const ch252Buggy = `package main

import (
	"errors"
	"sync"
)

func MergeScans(scans []int, checkpoint func(int) func()) (int, error) {
	if checkpoint == nil || len(scans) > 64 {
		return 0, errors.New("invalid scans")
	}
	total := 0
	var wg sync.WaitGroup
	for i, n := range scans {
		wg.Add(1)
		go func() { defer wg.Done(); old := total; ack := checkpoint(i); total = old + n; ack() }()
	}
	wg.Wait()
	return total, nil
}
func main() {}
`
const ch252Correct = `package main

import (
	"errors"
	"sync"
)

// A data race has unordered conflicting accesses; the mutex establishes happens-before.
// Reading and updating must share the critical section to avoid stale additions.
func MergeScans(scans []int, checkpoint func(int) func()) (int, error) {
	if checkpoint == nil || len(scans) > 64 {
		return 0, errors.New("invalid scans")
	}
	total := 0
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i, n := range scans {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ack := checkpoint(i)
			mu.Lock()
			total += n
			mu.Unlock()
			ack()
		}()
	}
	wg.Wait()
	return total, nil
}
func main() {}
`
const ch252TestStarter = `package main

import "testing"

func TestScans(t *testing.T) {
	got, err := MergeScans([]int{3}, func(int) func() { return func() {} })
	if err != nil || got != 3 {
		t.Fatal("single scan")
	}
}
`
const ch252Hidden = `package main

import (
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"slices"
	"strings"
	"sync"
	"testing"
)

// Cases: empty, signed/zero payloads, maximum batch, invalid callback/size, concurrent independent calls.
func TestHiddenScansExactAndCallbacks(t *testing.T) {
	rng := rand.New(rand.NewSource(252))
	for _, size := range []int{0, 1, 2, 17, 64} {
		scans := make([]int, size)
		want := 0
		for i := range scans {
			scans[i] = rng.Intn(101) - 50
			want += scans[i]
		}
		before := slices.Clone(scans)
		visits := make([]int, size)
		acks := make([]int, size)
		var mu sync.Mutex
		got, e := MergeScans(scans, func(i int) func() {
			mu.Lock()
			defer mu.Unlock()
			if i < 0 || i >= size {
				t.Error("index")
				return func() {}
			}
			visits[i]++
			return func() { mu.Lock(); acks[i]++; mu.Unlock() }
		})
		if e != nil || got != want || !slices.Equal(before, scans) {
			t.Fatal("sum/ownership", got, want, e)
		}
		for i := range scans {
			if visits[i] != 1 || acks[i] != 1 {
				t.Fatal("callback count", i)
			}
		}
	}
	if got, e := MergeScans(nil, nil); got != 0 || e == nil {
		t.Fatal("nil callback")
	}
	if got, e := MergeScans(make([]int, 65), func(int) func() { t.Error("invalid callback"); return func() {} }); got != 0 || e == nil {
		t.Fatal("oversize")
	}
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got, e := MergeScans([]int{9, -3, 2}, func(int) func() { return func() {} })
			if e != nil || got != 8 {
				t.Error("independent calls", got, e)
			}
		}()
	}
	wg.Wait()
}
func TestHiddenRaceExplanation(t *testing.T) {
	src, e := os.ReadFile("main.go")
	if e != nil {
		t.Fatal(e)
	}
	f, e := parser.ParseFile(token.NewFileSet(), "main.go", src, parser.ParseComments)
	if e != nil {
		t.Fatal(e)
	}
	explanation := ""
	for _, c := range f.Comments {
		explanation += c.Text()
	}
	if !strings.Contains(strings.ToLower(explanation), "race") || !strings.Contains(strings.ToLower(explanation), "happens-before") {
		t.Fatal("explain data race and happens-before in a comment")
	}
}
`

func init() {
	add(&Level{ID: "25-3-token-bucket", Chapter: 25, Kind: "lab", Race: true, LessonRevision: 1, Steps: ch25BucketArabic, Title: "خزان التصاريح", Goal: "نظّم دفعات إطلاق الشحنات بخزان تصاريح بيتجدد بساعة محقونة",
		Requires: []string{"mutex"}, Bridge: []int{36}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: ch253Starter}},
		Lesson: `# تصريح لكل فترة

## قبل ما تكتب
- هل رفض طلب كبير يستهلك جزء من الرصيد؟
- إزاي تحتفظ بأجزاء الفترة من غير تخزين رصيد فوق السعة؟

## token bucket
الخزان يبدأ مليان، وده يسمح بدفعة فورية بحجم السعة. بعد استهلاكه، كل فترة كاملة تضيف تصريح واحد. الطلب يستهلك كميته كلها أو يترفض كلها؛ مفيش انتظار أو طابور.

~PermitTank~ مختلف عن semaphore: التصريح المستهلك مش بيرجع لما الشحنة تخلص. الوقت هو اللي بيجدده. هنحسب التجديد عند الطلب، من غير goroutine ولا ticker.
~~~go
clock := labkit.NewFakeClock(time.Unix(0, 0))
tank, err := NewPermitTank(clock, 4, time.Second)
if err != nil { return err }
ok, err := tank.Admit(3)
clock.Advance(time.Second)
fmt.Println(ok, err)
~~~
~NewPermitTank~ بتحدد الساعة والسعة وفترة التصريح. ~Admit~ قرار فوري؛ ~false,nil~ نقص عادي و ~false,error~ مدخل غلط.

## حساب كسور الزمن
~Sub~ ترجع مدة، والقسمة الصحيحة تحسب الفترات الكاملة. حرك العلامة بعدد الفترات المكتملة بس عشان الجزء الباقي يفضل محسوب.
~~~go
elapsed := now.Sub(last)
steps := elapsed / every
last = last.Add(steps * every)
~~~
لو الرصيد وصل السعة، امسح الزمن الزائد وابدأ العلامة عند ~now~: مفيش تصريح كامل أو كسر تصريح يتخزن فوق السعة. لو الساعة رجعت للخلف، تجاهل القراءة ومتغيرش العلامة؛ كده الرجوع لنفس اللحظة مش بيكرر الرصيد.

## قرار واحد بقفل واحد
~sync.Mutex~ تحمي قراءة الساعة وعلامة الزمن والرصيد وفحص الطلب والاستهلاك مع بعض. ارفض الطلب غير الموجب قبل تغيير أي حالة. الرفض لنقص الرصيد يسمح بحساب الزمن، لكنه ميستهلكش.
~~~go
mu.Lock()
defer mu.Unlock()
// تحقق، جدد، افحص، واستهلك هنا
~~~
الساعة المحقونة آمنة للتزامن و ~Now~ بترجع فوراً. مفيش I/O تحت القفل.

## أخطاء شائعة
~~~go
last = now // غلط لو الرصيد لسه ناقص: كل كسر فترة بيضيع
credits -= cost // غلط قبل التأكد إن الرصيد كفاية
~~~

## المهمة
نفّذ ~NewPermitTank(clock labkit.Clock, capacity int, every time.Duration) (*PermitTank,error)~ و ~Admit(cost int) (bool,error)~. ارفض ساعة ~nil~ وسعة أو فترة غير موجبة بـ ~nil,error~. ابدأ برصيد السعة. الطلب غير الموجب ~false,error~ من غير تغيير؛ الطلب الأكبر من الرصيد ~false,nil~ من غير استهلاك؛ غير كده استهلك بالظبط وارجع ~true,nil~.

قبل القرار الصحيح احسب تصريح لكل ~every~ كاملة من آخر علامة، مع الاحتفاظ بباقي الزمن لو الخزان لسه مش مليان. لو التجديد يملأ الخزان، اقف عند السعة وامسح كل الرصيد الزمني الزائد بعلامة ~now~. الزمن الثابت أو الراجع للخلف ميجددش وميرجعش العلامة. كل ده تحت نفس ~sync.Mutex~، والحسابات متعملش overflow حتى بعد قفزة كبيرة. الحسابات مستقلة؛ ممنوع نسخها بعد الاستخدام. مفيش نوم أو وقت حقيقي أو goroutines للساعة أو ~Report~.`,
		Hints: []string{"استخدم مرور الوقت كمدخل للحساب.", "افصل عدد الفترات عن عدد التصاريح المتاحة.", "راجع قسمة ~time.Duration~ على ~every~ في العقد."},
		En: &LevelText{LessonRevision: 1, Steps: ch25BucketEnglish, Title: "Permit Tank", Goal: "Control cargo launch bursts with an injected-clock token bucket", Files: []LabFile{{Name: "main.go", Starter: ch253Starter}},
			Lesson: `# One Permit Per Interval

## Before you write
- Does rejecting a large request consume partial credit?
- How can you retain partial intervals without banking credit above capacity?

## Token bucket
A full initial tank permits an immediate burst up to capacity. After consumption, each complete interval earns one permit. Admission consumes the entire request or rejects it; no waiting or queue is involved.

~PermitTank~ differs from a semaphore: finishing cargo work does not return permits. Time replenishes them. Refill lazily at admission, without ticker goroutines.
~~~go
clock := labkit.NewFakeClock(time.Unix(0, 0))
tank, err := NewPermitTank(clock, 4, time.Second)
if err != nil { return err }
ok, err := tank.Admit(3)
clock.Advance(time.Second)
fmt.Println(ok, err)
~~~
~NewPermitTank~ selects clock, capacity, and interval. ~Admit~ decides immediately: ~false,nil~ means insufficient credit, while ~false,error~ means invalid input.

## Fractional time
~Sub~ returns elapsed duration; integer division counts complete intervals. Move the time marker only by complete intervals to retain the remainder.
~~~go
elapsed := now.Sub(last)
steps := elapsed / every
last = last.Add(steps * every)
~~~
When refill reaches capacity, discard surplus time and set the marker to ~now~. Neither whole nor fractional excess credit can be banked above capacity. Ignore backward readings without moving the marker: returning to the same instant cannot mint duplicate credit.

## One decision under one lock
A ~sync.Mutex~ protects clock observation, marker, balance, affordability, and consumption together. Reject nonpositive cost before any mutation. Insufficient-credit rejection accounts for elapsed time but consumes nothing.
~~~go
mu.Lock()
defer mu.Unlock()
// Validate, refill, check, and consume here
~~~
The injected clock is concurrency-safe and ~Now~ returns promptly. No external I/O belongs under the lock.

## Common mistakes
~~~go
last = now // Loses partial intervals when the tank remains below capacity
credits -= cost // Wrong before checking affordability
~~~

## Your mission
Implement ~NewPermitTank(clock labkit.Clock, capacity int, every time.Duration) (*PermitTank,error)~ and ~Admit(cost int) (bool,error)~. Reject a ~nil~ clock or nonpositive capacity/interval with ~nil,error~. Start full. Nonpositive requests return ~false,error~ without mutation; requests exceeding available credit return ~false,nil~ without consumption; other requests consume exactly their cost and return ~true,nil~.

Before each valid decision, earn one permit per complete ~every~ since the marker, retaining fractional time while below capacity. Refill reaching capacity clamps balance and discards all surplus time by setting the marker to ~now~. Frozen/backward time earns nothing and never moves the marker backward. Protect everything with one ~sync.Mutex~ and avoid arithmetic overflow even after large jumps. Instances are independent and must not be copied after use. No sleeping, real time, clock goroutines, or ~Report~.`,
			Hints: []string{"Treat elapsed time as an input to the calculation.", "Separate interval count from available permits.", "Review division of ~time.Duration~ by ~every~ in the contract."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch253Hidden},
	})
}

const ch253Starter = `package main

import (
	"errors"
	"orbital/colony/labkit"
	"time"
)

type PermitTank struct{}

func NewPermitTank(clock labkit.Clock, capacity int, every time.Duration) (*PermitTank, error) {
	return nil, errors.New("TODO")
}
func (p *PermitTank) Admit(cost int) (bool, error) { return false, errors.New("TODO") }
func main()                                        {}
`
const ch253Hidden = `package main

import (
	"math/rand"
	"orbital/colony/labkit"
	"sync"
	"testing"
	"time"
)

func permitWant(t *testing.T, p *PermitTank, n int, want bool) {
	t.Helper()
	got, e := p.Admit(n)
	if e != nil || got != want {
		t.Fatalf("Admit(%d)=%v,%v want %v", n, got, e, want)
	}
}

// Cases: initial burst, rejected/invalid cost, fractional/exact refill, cap saturation and concurrent bursts.
func TestHiddenPermitBoundaries(t *testing.T) {
	rng := rand.New(rand.NewSource(253))
	for trial := 0; trial < 30; trial++ {
		cap := 2 + rng.Intn(20)
		every := time.Duration(2+rng.Intn(500)) * time.Nanosecond
		clock := labkit.NewFakeClock(time.Unix(100, 0))
		p, e := NewPermitTank(clock, cap, every)
		if e != nil || p == nil {
			t.Fatal(e)
		}
		permitWant(t, p, cap+1, false)
		permitWant(t, p, cap, true)
		permitWant(t, p, 1, false)
		clock.Advance(every - 1)
		permitWant(t, p, 1, false)
		clock.Advance(1)
		permitWant(t, p, 1, true)
		for i := time.Duration(0); i < every; i++ {
			clock.Advance(1)
			if i < every-1 {
				permitWant(t, p, 1, false)
			}
		}
		permitWant(t, p, 1, true)
		clock.Advance(every*time.Duration(cap) + every/2)
		permitWant(t, p, cap, true)
		clock.Advance(every - every/2 - 1)
		permitWant(t, p, 1, false)
		clock.Advance(every/2 + 1)
		permitWant(t, p, 1, true)
		for _, n := range []int{0, -1} {
			if ok, e := p.Admit(n); ok || e == nil {
				t.Fatal("invalid request")
			}
		}
	}
	for _, c := range []struct {
		capacity int
		every    time.Duration
	}{{0, 1}, {-1, 1}, {1, 0}, {1, -1}} {
		if p, e := NewPermitTank(labkit.NewFakeClock(time.Time{}), c.capacity, c.every); p != nil || e == nil {
			t.Fatal("config")
		}
	}
	if p, e := NewPermitTank(nil, 1, 1); p != nil || e == nil {
		t.Fatal("nil clock")
	}
}

type permitClock struct {
	labkit.Clock
	now time.Time
}

func (c *permitClock) Now() time.Time { return c.now }
func TestHiddenPermitBackwardAndLongJump(t *testing.T) {
	clock := &permitClock{now: time.Unix(100, 0)}
	p, e := NewPermitTank(clock, 2, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	permitWant(t, p, 2, true)
	clock.now = clock.now.Add(-time.Hour)
	permitWant(t, p, 1, false)
	clock.now = time.Unix(100, 0)
	permitWant(t, p, 1, false)
	clock.now = clock.now.Add(time.Second)
	permitWant(t, p, 1, true)
	clock.now = clock.now.Add(time.Duration(1<<63 - 1))
	permitWant(t, p, 3, false)
	permitWant(t, p, 2, true)
	permitWant(t, p, 1, false)
}
func TestHiddenPermitConcurrentBudget(t *testing.T) {
	rng := rand.New(rand.NewSource(2531))
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	capacity := 5 + rng.Intn(20)
	p, e := NewPermitTank(clock, capacity, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	for round := 0; round < 4; round++ {
		start := make(chan struct{})
		out := make(chan bool, 100)
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				ok, e := p.Admit(1)
				if e != nil {
					t.Error(e)
				}
				out <- ok
			}()
		}
		close(start)
		wg.Wait()
		close(out)
		got := 0
		for ok := range out {
			if ok {
				got++
			}
		}
		if got != capacity {
			t.Fatal("budget", got, capacity)
		}
		clock.Advance(time.Duration(capacity) * time.Second)
	}
	other, e := NewPermitTank(clock, 1, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	permitWant(t, other, 1, true)
	permitWant(t, other, 1, false)
}
`

func init() {
	add(&Level{ID: "25-4-event-bus", Chapter: 25, Kind: "lab", Race: true, LessonRevision: 1, Steps: ch25BusArabic, Title: "ناقل الأحداث", Goal: "وزّع إشارات الرصيف على المشتركين من غير ما البطيء يعطّلهم",
		Requires: []string{"chan", "mutex"}, Bridge: []int{42}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: ch254Starter}},
		Lesson: `# إشارات لكل مشترك

## قبل ما تكتب
- إيه اللي يحصل لو قناة مشترك مليانة؟
- إزاي تمنع إرسال حدث على قناة اتقفلت بالتزامن؟

## بث مش توزيع شغل
في ناقل الأحداث، كل مشترك في الموضوع يستقبل نسخة. ده مختلف عن عمال بيقرأوا قناة واحدة: هناك عنصر واحد يروح لعامل واحد. ~DockSignal~ فيها موضوع ونص غير قابل للتعديل.
~~~go
signal := DockSignal{Topic: "airlock", Text: "sealed"}
var bus DockBus
id, events, err := bus.Subscribe(signal.Topic, 2)
fmt.Println(id, events, err)
~~~
~DockBus~ بصفره جاهز. ~Subscribe~ تعمل قناة خاصة بموضوع وسعة موجبة وترجع رقم وقناة استقبال وخطأ. الاشتراكات المكررة مستقلة.

## سياسة المشترك البطيء
سعة القناة حد الذاكرة لكل مشترك. السياسة هنا **drop newest**: لما القناة تتملي ارفض الحدث الجديد للمشترك ده فقط، وسيب القديم. ~select~ مع ~default~ تمنع انتظار المستهلك.
~~~go
select {
case events <- signal:
 delivered++
default:
 // الجديد اتساب، القديم لسه موجود
}
~~~
~Publish~ ترجع عدد القنوات اللي قبلت الحدث، مش عدد المشتركين ولا عدد اللي قرأوه بالفعل. ده يسمح بمراجعة الفقد من غير تخمين.
~~~go
count, err := bus.Publish(DockSignal{Topic: "airlock", Text: "open"})
fmt.Println(count, err)
~~~

## القفل والملكية
الناقل هو المرسل والمالك الوحيد للقنوات. نفس ~sync.Mutex~ تحمي خريطة الاشتراكات والإرسال وإزالة الاشتراك والقفل النهائي. الإرسال غير منتظر، فمفيش callback أو I/O تحت القفل.

~Unsubscribe~ تقفل قناة الاشتراك وتشيله؛ القديم المخزن يفضل قابل للقراءة. ~Close~ تقفل الباقي مرة واحدة؛ تكرارها آمن حتى بالتزامن، وده معنى idempotent.
~~~go
removed := bus.Unsubscribe(id)
bus.Close()
bus.Close()
for event := range events { fmt.Println(event.Text) }
fmt.Println(removed)
~~~
بعد القفل، الاشتراكات والنشر ترجع خطأ. الأرقام الموجبة بتزيد من غير إعادة استخدام؛ لو وصلت ~math.MaxInt64~ ارفض أي اشتراك جديد. المثال ده للحد بس:
~~~go
if next == math.MaxInt64 { return errors.New("subscription IDs exhausted") }
next++
~~~

## أخطاء شائعة
~~~go
mu.Unlock()
events <- signal // ممكن Unsubscribe تكون قفلت القناة؛ وكمان ممكن يعلق
~~~
نسخ القنوات خارج القفل مش كفاية هنا؛ لازم ضمان عمرها لحد نهاية الإرسال. مفيش goroutines مطلوبة في الناقل نفسه.

## المهمة
نفّذ ~DockSignal{Topic,Text string}~ و ~DockBus~ بصفر شغال، و ~Subscribe(topic string, capacity int) (int64,<-chan DockSignal,error)~ و ~Unsubscribe(id int64) bool~ و ~Publish(signal DockSignal) (int,error)~ و ~Close()~.

الموضوع الفاضي أو السعة غير الموجبة أو الناقل المقفول يرفض الاشتراك بـ ~0,nil,error~ من غير تغيير. أول رقم 1، ثم أرقام متزايدة لا تتكرر؛ ارفض بعد ~math.MaxInt64~. الموضوع الفاضي أو الناقل المقفول يرفض النشر بـ ~0,error~. النص الفاضي صالح. النشر من غير مشتركين نجاح بعدد 0. كل اشتراك مطابق يقبل الحدث مرة أو يسقط الجديد لو ممتلئ؛ غير المطابق ميستقبلش. ترتيب النشر المتسلسل محفوظ لكل مشترك، ومفيش ترتيب مضمون بين نداءات متوازية.

إلغاء موجود يرجع ~true~ ويقفل قناته مرة؛ الغريب أو المكرر ~false~. القفل النهائي يقفل الكل ويمنع النشر والاشتراك الجديد؛ القيم المخزنة تتصرف قبل نهاية القناة. كل الطرق آمنة بالتزامن، مفيش إرسال بعد القفل ولا انتظار مستهلك ولا نسخ للناقل بعد استخدامه. مفيش goroutines داخلية أو وقت أو شبكة أو ~Report~.

النجمة التالتة: نفس دقة عدد التسليمات والقفل مرة واحدة لما ناشرين ونداءات ~Close~ يتسابقوا.`,
		Hints: []string{"حدد معنى الإرسال المقبول للمشترك.", "افحص حالة القناة قبل التعامل مع إغلاقها.", "راجع فرع ~default~ في ~select~ عند إرسال غير حاجب."},
		En: &LevelText{LessonRevision: 1, Steps: ch25BusEnglish, Title: "Event Bus", Goal: "Broadcast dock signals without allowing slow subscribers to block others", Files: []LabFile{{Name: "main.go", Starter: ch254Starter}},
			Lesson: `# One Signal for Every Subscriber

## Before you write
- What happens when a subscriber channel is full?
- How will you prevent sending concurrently with channel closure?

## Broadcast rather than distribute
Every matching subscriber gets a copy. Workers sharing one input channel instead divide work. ~DockSignal~ carries a topic and immutable text.
~~~go
signal := DockSignal{Topic: "airlock", Text: "sealed"}
var bus DockBus
id, events, err := bus.Subscribe(signal.Topic, 2)
fmt.Println(id, events, err)
~~~
A zero ~DockBus~ works. ~Subscribe~ creates a private channel with positive capacity, returning an ID, receive-only channel, and error. Duplicate subscriptions are independent.

## Slow-consumer policy
Channel capacity bounds per-subscriber memory. Use **drop newest**: a full channel drops only the new signal for that subscriber, preserving buffered signals. ~select~ with ~default~ avoids waiting for consumption.
~~~go
select {
case events <- signal:
 delivered++
default:
 // Preserve buffered signals; discard the new one
}
~~~
~Publish~ returns the number of channels accepting the signal, rather than subscriber count or completed reads. Delivery loss becomes observable.
~~~go
count, err := bus.Publish(DockSignal{Topic: "airlock", Text: "open"})
fmt.Println(count, err)
~~~

## Locking and ownership
The bus exclusively sends and closes its channels. One ~sync.Mutex~ protects the subscription map, sending, unsubscribe, and final closure. Nonblocking sends require no callbacks or external I/O under the lock.

~Unsubscribe~ removes and closes a subscription; its buffered events remain readable. ~Close~ closes remaining channels once. Repeated or concurrent closure is safe: this is idempotence.
~~~go
removed := bus.Unsubscribe(id)
bus.Close()
bus.Close()
for event := range events { fmt.Println(event.Text) }
fmt.Println(removed)
~~~
After closure, subscribe and publish fail. Positive IDs increase without reuse; reject new subscriptions after ~math.MaxInt64~. This tiny example illustrates exhaustion only:
~~~go
if next == math.MaxInt64 { return errors.New("subscription IDs exhausted") }
next++
~~~

## Common mistakes
~~~go
mu.Unlock()
events <- signal // Unsubscribe may have closed it; it can also block
~~~
Snapshotting channels outside the lock is insufficient here: channel lifetime must cover the send. The bus itself needs no goroutines.

## Your mission
Implement ~DockSignal{Topic,Text string}~, zero-ready ~DockBus~, ~Subscribe(topic string, capacity int) (int64,<-chan DockSignal,error)~, ~Unsubscribe(id int64) bool~, ~Publish(signal DockSignal) (int,error)~, and ~Close()~.

Empty topic, nonpositive capacity, or closed bus rejects subscription with ~0,nil,error~ without mutation. IDs start at 1, increase without reuse, and reject after ~math.MaxInt64~. Empty topic or closed bus rejects publish with ~0,error~. Empty text is valid. Publishing without subscribers succeeds with count 0. Every matching subscriber accepts once or drops the newest signal when full; other topics receive nothing. Sequential publish order is preserved per subscriber; concurrent calls have no promised global order.

Removing an existing ID returns ~true~ and closes its channel once; unknown/repeated removal returns ~false~. Final closure closes all channels and prevents new subscriptions/publications; buffered events drain before channel termination. All methods are concurrency-safe, never send after closure, and never wait for consumption. Do not copy the bus after use. No internal goroutines, time, network, or ~Report~.

Third star: preserve exact delivery accounting and single closure when publishers race with ~Close~ calls.`,
			Hints: []string{"Define what counts as an accepted subscriber send.", "Check channel state before handling closure.", "Review the ~default~ branch for a nonblocking send."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch254Hidden, Bonus: true},
	})
}

const ch254Starter = `package main

import "errors"

type DockSignal struct{ Topic, Text string }
type DockBus struct{}

func (b *DockBus) Subscribe(topic string, capacity int) (int64, <-chan DockSignal, error) {
	return 0, nil, errors.New("TODO")
}
func (b *DockBus) Unsubscribe(id int64) bool              { return false }
func (b *DockBus) Publish(signal DockSignal) (int, error) { return 0, errors.New("TODO") }
func (b *DockBus) Close()                                 {}
func main()                                               {}
`
const ch254Hidden = `package main

import (
	"fmt"
	"math/rand"
	"orbital/colony/labkit"
	"slices"
	"sync"
	"testing"
)

// Cases: empty bus, invalid input, zero text, topics, duplicates, full queues, buffered closure.
func TestHiddenBusTopicsAndSlowConsumer(t *testing.T) {
	var b DockBus
	defer b.Close()
	for _, c := range []struct {
		topic string
		cap   int
	}{{"", 1}, {"x", 0}, {"x", -1}} {
		if id, ch, e := b.Subscribe(c.topic, c.cap); id != 0 || ch != nil || e == nil {
			t.Fatal("invalid subscription")
		}
	}
	if n, e := b.Publish(DockSignal{}); n != 0 || e == nil {
		t.Fatal("empty topic")
	}
	if n, e := b.Publish(DockSignal{"none", ""}); n != 0 || e != nil {
		t.Fatal("empty bus")
	}
	id, slow, e := b.Subscribe("dock", 1)
	if e != nil || id != 1 || slow == nil {
		t.Fatal("first subscription")
	}
	id2, fast, e := b.Subscribe("dock", 3)
	if e != nil || id2 <= id || fast == nil {
		t.Fatal("duplicate subscription")
	}
	otherID, other, e := b.Subscribe("ore", 1)
	if e != nil || otherID <= id2 || other == nil {
		t.Fatal("other topic")
	}
	rng := rand.New(rand.NewSource(254))
	want := []DockSignal{{"dock", ""}, {"dock", fmt.Sprint(rng.Int())}, {"dock", fmt.Sprint(rng.Int())}}
	for i, s := range want {
		n, e := b.Publish(s)
		count := 1
		if i == 0 {
			count = 2
		}
		if e != nil || n != count {
			t.Fatal("drop newest", n, count, e)
		}
	}
	if s := <-slow; s != want[0] {
		t.Fatal("old event overwritten")
	}
	for _, s := range want {
		if got := <-fast; got != s {
			t.Fatal("order", got, s)
		}
	}
	select {
	case <-other:
		t.Fatal("topic isolation")
	default:
	}
	if !b.Unsubscribe(id) || b.Unsubscribe(id) || b.Unsubscribe(-1) {
		t.Fatal("removal")
	}
	if _, ok := <-slow; ok {
		t.Fatal("unsubscribe must close")
	}
	newID, newEvents, e := b.Subscribe("dock", 1)
	if e != nil || newID <= otherID || newEvents == nil {
		t.Fatal("ID reuse")
	}
	if n, e := b.Publish(want[1]); n != 2 || e != nil {
		t.Fatal("remaining subscriptions")
	}
	b.Close()
	b.Close()
	if got := <-newEvents; got != want[1] {
		t.Fatal("buffer lost at close")
	}
	if _, ok := <-newEvents; ok {
		t.Fatal("close")
	}
	if _, ok := <-other; ok {
		t.Fatal("other closure")
	}
	if got := <-fast; got != want[1] {
		t.Fatal("fast buffer")
	}
	if _, ok := <-fast; ok {
		t.Fatal("fast closure")
	}
	if b.Unsubscribe(id2) {
		t.Fatal("closed removal")
	}
	if id, ch, e := b.Subscribe("dock", 1); id != 0 || ch != nil || e == nil {
		t.Fatal("subscribe after close")
	}
	if n, e := b.Publish(want[0]); n != 0 || e == nil {
		t.Fatal("publish after close")
	}
}
func TestHiddenBusConcurrentSubscribeAndPublish(t *testing.T) {
	labkit.LeakCheck(t)
	var b DockBus
	defer b.Close()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			id, ch, e := b.Subscribe(fmt.Sprint(i%3), 2)
			if e != nil {
				t.Error(e)
				return
			}
			for j := 0; j < 10; j++ {
				_, e = b.Publish(DockSignal{fmt.Sprint(j % 3), fmt.Sprint(i, j)})
				if e != nil {
					t.Error(e)
				}
			}
			if !b.Unsubscribe(id) {
				t.Error("missing active subscription")
			}
			for range ch {
			}
		}()
	}
	close(start)
	wg.Wait()
	if n, e := b.Publish(DockSignal{"0", "done"}); n != 0 || e != nil {
		t.Fatal("stale subscribers")
	}
}
func TestHiddenBonusBusConcurrentCloseAccounting(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(2541))
	for trial := 0; trial < 5; trial++ {
		var b DockBus
		_, ch, e := b.Subscribe("dock", 128)
		if e != nil {
			t.Fatal(e)
		}
		// Guarantee a buffered event, then race publishers with several closers.
		if n, e := b.Publish(DockSignal{"dock", "first"}); n != 1 || e != nil {
			t.Fatal("initial publish")
		}
		start := make(chan struct{})
		accepted := make(chan string, 64)
		var wg sync.WaitGroup
		payloads := make([]string, 64)
		for i := range payloads {
			payloads[i] = fmt.Sprint(i, "-", rng.Int())
		}
		for _, text := range payloads {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				n, e := b.Publish(DockSignal{"dock", text})
				if e == nil {
					if n != 1 {
						t.Error("delivery count")
					}
					accepted <- text
				} else if n != 0 {
					t.Error("failed publication count")
				}
			}()
		}
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() { defer wg.Done(); <-start; b.Close() }()
		}
		close(start)
		wg.Wait()
		close(accepted)
		b.Close()
		want := []string{"first"}
		for text := range accepted {
			want = append(want, text)
		}
		var got []string
		for s := range ch {
			got = append(got, s.Text)
		}
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Fatal("close lost accepted events", got, want)
		}
	}
}
`

func init() {
	add(&Level{ID: "25-5-ttl-cache", Chapter: 25, Kind: "lab", Race: true, LessonRevision: 1, Steps: ch25CacheArabic, Title: "كاش بصلاحية", Goal: "خزّن قراءات الملاحة بصلاحية ونظّفها بعامل قابل للقفل من غير تسريب",
		Requires: []string{"sel:sync.RWMutex"}, Bridge: []int{43}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: ch255Starter}},
		Lesson: `# قراءة قديمة مش قراءة صالحة

## قبل ما تكتب
- هل وجود المفتاح في الخريطة معناه إن قراءته لسه صالحة؟
- مين يوقف عامل التنظيف ويستناه لما اللوحة تتقفل؟

## صلاحية بعقد واضح
~TTL~ مدة من آخر كتابة، مش من آخر قراءة. خزّن القيمة ووقت انتهائها مع بعض. عند ~now == expires~ القيمة منتهية؛ القراءة متجددش المدة.
~~~go
type reading struct { text string; expires time.Time }
r := reading{text: "clear", expires: clock.Now().Add(time.Second)}
live := clock.Now().Before(r.expires)
fmt.Println(live)
~~~
~Lookup~ لازم ترفض المنتهي فوراً حتى لو عامل التنظيف لسه مشتغلش. الكتابة لنفس المفتاح تستبدل القيمة والميعاد. النص الفاضي قيمة صالحة، عشان كده النتيجة فيها ~bool~.

## RWMutex للقراءات
~sync.RWMutex~ تسمح لقراء كتير مع بعض بـ ~RLock~ و ~RUnlock~. الكتابة والمسح والقفل النهائي بـ ~Lock~ و ~Unlock~ حصرياً. متحاولش ترقّي read lock لwrite lock وإنت لسه ماسكها.
~~~go
mu.RLock()
r, ok := readings[key]
mu.RUnlock()
// لو محتاج حذف، خد Lock وافحص تاني؛ نسخة r ممكن تكون قديمة
~~~
في المهمة القراءة مش محتاجة تحذف أصلاً. ~Stored~ تعد العناصر الفعلية، حتى المنتهية اللي لسه متنظفتش؛ دي قياس لذاكرة الخريطة، مش عدد القراءات الصالحة.

## تنظيف بساعة محقونة
~labkit.Clock.NewTicker~ ترجع ~labkit.Ticker~ فيها ~C()~ لقناة النبضات و ~Stop()~ لإيقافها. اعمل ticker قبل رجوع constructor، وعامل واحد ينضف بكل نبضة بناءً على ~clock.Now()~، مش وقت النبضة القديم اللي ممكن يكون متأخر.
~~~go
ticker := clock.NewTicker(time.Second)
go func() {
 defer ticker.Stop()
 for {
  select {
  case <-ticker.C():
   // احذف المنتهي تحت Lock
  case <-stop:
   return
  }
 }
}()
~~~
~Sweeps~ قناة تقارير بعد كل تنظيف، تبعت عدد المحذوفات من غير انتظار، بسعة 1 وسياسة إسقاط الجديد لو ممتلئة. كده التست يستنى اكتمال دورة بعد ~FakeClock.Advance~ بدل تخمين موعد تنفيذ goroutine.
~~~go
clock := labkit.NewFakeClock(time.Unix(0, 0))
board, err := NewExpiryBoard(clock, time.Second)
if err != nil { return err }
err = board.Put("north", "clear", time.Second)
clock.Advance(time.Second)
removed := <-board.Sweeps()
fmt.Println(removed, board.Stored(), err)
board.Close()
~~~
~NewExpiryBoard~ تملك العامل. ~Put~ تكتب، و ~Lookup~ تقرا، و ~Stored~ تعد الذاكرة. ~Sweeps~ للملاحظة بس، المستهلك ميقفلهاش.

## قفل العمر كله
~Close~ تعلّم اللوحة مقفولة وتقفل إشارة الإيقاف مرة واحدة تحت القفل، وبعدها تستنى العامل خارج القفل. العامل يوقف ticker ويقفل قناة التقارير ثم إشارة الانتهاء. كل نداء ~Close~ يستنى نفس الانتهاء.
~~~go
mu.Lock()
if !closed { closed = true; close(stop) }
mu.Unlock()
<-done
~~~

## أخطاء شائعة
~~~go
mu.RLock()
delete(readings, key) // حذف تحت read lock غلط
mu.RUnlock()
// close(stop) مرتين تعمل panic؛ انتظار العامل تحت القفل ممكن يعمل deadlock
~~~

## المهمة
نفّذ ~NewExpiryBoard(clock labkit.Clock, sweep time.Duration) (*ExpiryBoard,error)~ و ~Put(key,text string, ttl time.Duration) error~ و ~Lookup(key string) (string,bool)~ و ~Stored() int~ و ~Sweeps() <-chan int~ و ~Close()~. ساعة ~nil~ أو فترة تنظيف غير موجبة ترجع ~nil,error~ من غير عامل. اللوحة الناتجة تبدأ فاضية؛ هي مش zero-value API وممنوع نسخها.

ارفض مفتاح فاضي أو ~ttl <= 0~ أو كتابة بعد القفل بخطأ من غير تغيير. النص الفاضي صالح. الكتابة تستبدل القيمة وموعد الانتهاء عند ~clock.Now().Add(ttl)~. القراءة المفقودة أو المنتهية أو بعد القفل ترجع ~"",false~؛ غير كده النص و ~true~. القراءة متجددش الصلاحية. ~Stored~ عدد عناصر الخريطة الفعلي، والمنتهي يتحذف في التنظيف. استخدم ~sync.RWMutex~ وقفل قراءة للقرايات وقفل حصري للتغيير.

ابدأ ticker محقون وعامل تنظيف واحد قبل رجوع constructor. كل نبضة تمسح كل ما موعده أقل من أو يساوي ~clock.Now()~ وترسل عدد المحذوفات في قناة ~Sweeps~ بسعة 1 من غير انتظار. القناة الممتلئة تسقط التقرير الجديد. ~Close~ آمنة ومتكررة ومتزامنة: تمنع الكتابة، تفضي الخريطة، توقف العامل وticker، وتقفل ~Sweeps~ بعد آخر إرسال، وتستنى ده كله قبل الرجوع. ممنوع نوم أو ساعة حقيقية أو شبكة أو ~Report~.`,
		Hints: []string{"فصل صلاحية القراءة عن عملية التنظيف.", "فكّر في إشارة الإيقاف وإشارة الانتهاء كل واحدة لوحدها.", "راجع ~RLock~ و~Close~ في العقد قبل تتبع العامل."},
		En: &LevelText{LessonRevision: 1, Steps: ch25CacheEnglish, Title: "TTL Cache", Goal: "Cache navigation readings with expiry and a joined cleanup worker", Files: []LabFile{{Name: "main.go", Starter: ch255Starter}},
			Lesson: `# An Old Reading Is Not a Live Reading

## Before you write
- Does map membership guarantee the reading is still valid?
- Who stops and joins cleanup when the board closes?

## Explicit expiry
~TTL~ is measured from the last write, not the last read. Store text and expiry together. At ~now == expires~, the entry is expired; reads do not renew it.
~~~go
type reading struct { text string; expires time.Time }
r := reading{text: "clear", expires: clock.Now().Add(time.Second)}
live := clock.Now().Before(r.expires)
fmt.Println(live)
~~~
~Lookup~ rejects expired entries immediately, even before cleanup runs. Overwriting replaces both text and expiry. Empty text is valid, so lookup also returns a ~bool~.

## RWMutex for readers
~sync.RWMutex~ permits concurrent readers with ~RLock~ and ~RUnlock~. Writes, deletion, and closure need exclusive ~Lock~ and ~Unlock~. Never upgrade while holding a read lock.
~~~go
mu.RLock()
r, ok := readings[key]
mu.RUnlock()
// Deletion would need Lock and rechecking: r may already be stale
~~~
Lookup need not delete here. ~Stored~ counts physical entries, including expired entries awaiting cleanup: it measures map occupancy, not live readings.

## Cleanup with injected time
~labkit.Clock.NewTicker~ returns a ~labkit.Ticker~, whose ~C()~ provides ticks and ~Stop()~ stops ticking. Create it before construction returns. One worker sweeps per tick using ~clock.Now()~, rather than a potentially stale tick timestamp.
~~~go
ticker := clock.NewTicker(time.Second)
go func() {
 defer ticker.Stop()
 for {
  select {
  case <-ticker.C():
   // Delete expired entries under Lock
  case <-stop:
   return
  }
 }
}()
~~~
~Sweeps~ reports deleted counts after cleanup, nonblocking with capacity 1 and drop-newest when full. Tests await cycle completion after ~FakeClock.Advance~ instead of guessing scheduling.
~~~go
clock := labkit.NewFakeClock(time.Unix(0, 0))
board, err := NewExpiryBoard(clock, time.Second)
if err != nil { return err }
err = board.Put("north", "clear", time.Second)
clock.Advance(time.Second)
removed := <-board.Sweeps()
fmt.Println(removed, board.Stored(), err)
board.Close()
~~~
~NewExpiryBoard~ owns its worker. ~Put~ writes, ~Lookup~ reads, and ~Stored~ counts occupancy. ~Sweeps~ is observation only; consumers never close it.

## Close the whole lifetime
Under the lock, ~Close~ marks closure and closes the stop signal once, then joins outside the lock. The worker stops the ticker, closes reports, and signals completion. Every ~Close~ caller waits for that same completion.
~~~go
mu.Lock()
if !closed { closed = true; close(stop) }
mu.Unlock()
<-done
~~~

## Common mistakes
~~~go
mu.RLock()
delete(readings, key) // Mutation under a read lock is unsafe
mu.RUnlock()
// Closing stop twice panics; joining under the lock can deadlock
~~~

## Your mission
Implement ~NewExpiryBoard(clock labkit.Clock, sweep time.Duration) (*ExpiryBoard,error)~, ~Put(key,text string, ttl time.Duration) error~, ~Lookup(key string) (string,bool)~, ~Stored() int~, ~Sweeps() <-chan int~, and ~Close()~. A ~nil~ clock or nonpositive sweep interval returns ~nil,error~ without a worker. Constructed boards start empty; this is not a zero-value API and must not be copied.

Reject empty keys, ~ttl <= 0~, or post-close writes with an error without mutation. Empty text is valid. Writes replace text and expiry with ~clock.Now().Add(ttl)~. Missing, expired, or post-close lookup returns ~"",false~; otherwise return text,~true~. Reads do not renew expiry. ~Stored~ counts physical map entries; cleanup removes expired entries. Use ~sync.RWMutex~ with read locks for reads and exclusive locks for mutation.

Start an injected ticker and one cleanup worker before construction returns. Each tick removes every entry expiring at or before ~clock.Now()~ and nonblockingly reports the deleted count through capacity-1 ~Sweeps~. A full report channel drops the new report. ~Close~ is safe, repeatable, and concurrent: reject future writes, clear storage, stop worker and ticker, close ~Sweeps~ after its last send, and join all of that before returning. No sleeping, real clock, network, or ~Report~.`,
			Hints: []string{"Separate lookup validity from cleanup work.", "Treat stop and completion as different signals.", "Review ~RLock~ and ~Close~ in the contract before tracing the worker."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch255Hidden},
	})
}

const ch255Starter = `package main

import (
	"errors"
	"orbital/colony/labkit"
	"time"
)

type ExpiryBoard struct{}

func NewExpiryBoard(clock labkit.Clock, sweep time.Duration) (*ExpiryBoard, error) {
	return nil, errors.New("TODO")
}
func (b *ExpiryBoard) Put(key, text string, ttl time.Duration) error { return errors.New("TODO") }
func (b *ExpiryBoard) Lookup(key string) (string, bool)              { return "", false }
func (b *ExpiryBoard) Stored() int                                   { return 0 }
func (b *ExpiryBoard) Sweeps() <-chan int                            { return nil }
func (b *ExpiryBoard) Close()                                        {}
func main()                                                          {}
`
const ch255Hidden = `package main

import (
	"fmt"
	"math/rand"
	"orbital/colony/labkit"
	"sync"
	"testing"
	"time"
)

// Cases: constructor/write rejection, exact expiry, overwrite, missing/stored empty, cleanup, close/backpressure.
func TestHiddenExpiryReadBoundariesAndReplacement(t *testing.T) {
	for _, d := range []time.Duration{0, -1} {
		if b, e := NewExpiryBoard(labkit.NewFakeClock(time.Time{}), d); b != nil || e == nil {
			t.Fatal("invalid sweep")
		}
	}
	if b, e := NewExpiryBoard(nil, time.Second); b != nil || e == nil {
		t.Fatal("nil clock")
	}
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	b, e := NewExpiryBoard(clock, time.Hour)
	if e != nil || b == nil {
		t.Fatal(e)
	}
	defer b.Close()
	if v, ok := b.Lookup("missing"); v != "" || ok {
		t.Fatal("missing")
	}
	if b.Stored() != 0 {
		t.Fatal("initial storage")
	}
	if b.Put("key", "old", 10*time.Second) != nil {
		t.Fatal("put")
	}
	for _, c := range []struct {
		k   string
		ttl time.Duration
	}{{"", 1}, {"key", 0}, {"key", -1}} {
		if b.Put(c.k, "bad", c.ttl) == nil {
			t.Fatal("invalid put")
		}
	}
	clock.Advance(9 * time.Second)
	if v, ok := b.Lookup("key"); !ok || v != "old" {
		t.Fatal("early expiry")
	}
	if b.Put("key", "", 2*time.Second) != nil {
		t.Fatal("overwrite")
	}
	clock.Advance(2*time.Second - 1)
	if v, ok := b.Lookup("key"); !ok || v != "" {
		t.Fatal("stored empty")
	}
	clock.Advance(1)
	if v, ok := b.Lookup("key"); ok || v != "" {
		t.Fatal("exact expiry/read renewed TTL")
	}
	if b.Stored() != 1 {
		t.Fatal("lookup should not remove storage")
	}
}

type expiryTicker struct {
	labkit.Ticker
	stopped chan struct{}
	once    sync.Once
}

func (t *expiryTicker) Stop() { t.Ticker.Stop(); t.once.Do(func() { close(t.stopped) }) }

type expiryClock struct {
	*labkit.FakeClock
	stopped chan struct{}
}

func (c *expiryClock) NewTicker(d time.Duration) labkit.Ticker {
	return &expiryTicker{Ticker: c.FakeClock.NewTicker(d), stopped: c.stopped}
}
func TestHiddenExpirySweepAndStop(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(255))
	clock := &expiryClock{labkit.NewFakeClock(time.Unix(100, 0)), make(chan struct{})}
	b, e := NewExpiryBoard(clock, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer b.Close()
	if cap(b.Sweeps()) != 1 {
		t.Fatal("bounded sweep telemetry")
	}
	expired := 0
	total := 40
	for i := 0; i < total; i++ {
		ttl := time.Duration(1+rng.Intn(5)) * time.Second
		if ttl <= 3*time.Second {
			expired++
		}
		if e := b.Put(fmt.Sprint(i), fmt.Sprint(i), ttl); e != nil {
			t.Fatal(e)
		}
	}
	// FakeClock delivers the earliest due timestamp; cleanup must instead observe current Now.
	clock.Advance(3 * time.Second)
	if n := <-b.Sweeps(); n != expired || b.Stored() != total-expired {
		t.Fatal("sweep current time", n, expired, b.Stored())
	}
	clock.Advance(10 * time.Second)
	if n := <-b.Sweeps(); n != total-expired || b.Stored() != 0 {
		t.Fatal("remaining cleanup", n)
	}
	// Leave a report unread; Close must still join the worker and stop the ticker.
	clock.Advance(time.Second)
	b.Close()
	b.Close()
	select {
	case <-clock.stopped:
	default:
		t.Fatal("ticker not stopped")
	}
	for range b.Sweeps() {
	}
	clock.Advance(time.Hour)
	if _, ok := <-b.Sweeps(); ok {
		t.Fatal("telemetry after close")
	}
	if b.Put("after", "x", time.Second) == nil || b.Stored() != 0 {
		t.Fatal("post-close mutation")
	}
}

type expiryObservedClock struct {
	*labkit.FakeClock
	observations chan struct{}
}

func (c *expiryObservedClock) Now() time.Time {
	now := c.FakeClock.Now()
	select {
	case c.observations <- struct{}{}:
	default:
	}
	return now
}
func TestHiddenExpiryFullTelemetryDoesNotBlockClose(t *testing.T) {
	labkit.LeakCheck(t)
	clock := &expiryObservedClock{labkit.NewFakeClock(time.Unix(0, 0)), make(chan struct{}, 1)}
	b, e := NewExpiryBoard(clock, time.Second)
	if e != nil {
		t.Fatal(e)
	}
	defer b.Close()
	if e := b.Put("live", "reading", time.Hour); e != nil {
		t.Fatal(e)
	}
	// Discard constructor/Put observations before driving two distinct sweeps.
	select {
	case <-clock.observations:
	default:
	}
	clock.Advance(time.Second)
	<-clock.observations
	clock.Advance(time.Second)
	<-clock.observations
	// The second sweep began after the first filled the report channel.
	b.Close()
	reports := 0
	for range b.Sweeps() {
		reports++
	}
	if reports != 1 {
		t.Fatal("full telemetry must drop newest", reports)
	}
}
func TestHiddenExpiryConcurrentOwnershipAndClose(t *testing.T) {
	labkit.LeakCheck(t)
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	b, e := NewExpiryBoard(clock, time.Hour)
	if e != nil {
		t.Fatal(e)
	}
	defer b.Close()
	rng := rand.New(rand.NewSource(2551))
	payloads := make([]string, 32)
	for i := range payloads {
		payloads[i] = fmt.Sprint(rng.Int())
		if e := b.Put(fmt.Sprint(i), payloads[i], 2*time.Hour); e != nil {
			t.Fatal(e)
		}
	}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i, text := range payloads {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			k := fmt.Sprint(i)
			for j := 0; j < 30; j++ {
				if e := b.Put(k, text, 2*time.Hour); e != nil {
					t.Error(e)
				}
				if v, ok := b.Lookup(k); !ok || v != text {
					t.Error("disjoint ownership")
				}
				_ = b.Stored()
			}
		}()
	}
	close(start)
	clock.Advance(time.Hour)
	<-b.Sweeps()
	wg.Wait()
	if b.Stored() != 32 {
		t.Fatal("lost writes")
	}
	start = make(chan struct{})
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 20; j++ {
				_ = b.Put("shared", "x", time.Second)
				_, _ = b.Lookup("shared")
				_ = b.Stored()
			}
			b.Close()
		}()
	}
	close(start)
	wg.Wait()
	if b.Stored() != 0 {
		t.Fatal("close did not clear")
	}
	if v, ok := b.Lookup("shared"); ok || v != "" {
		t.Fatal("post-close read")
	}
	for range b.Sweeps() {
	}
}
`

func init() {
	add(&Level{ID: "25-6-atomic-bench", Chapter: 25, Kind: "lab", Race: true, LessonRevision: 1, Steps: ch25AtomicArabic, Title: "atomic ولا mutex", Goal: "قارن عدّادي نبضات صحيحين تحت التزامن وقيس نفس الحمل على الاتنين",
		Requires: []string{"benchmark", "sel:atomic.Int64"}, Bridge: []int{45}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch256Main, ReadOnly: true}, {Name: "main_test.go", Starter: ch256Starter}},
		Lesson: `# اختار بعد ما تقيس

## قبل ما تكتب
- هل عندك رقم مستقل ولا شرط بيربط أكتر من حقل؟
- إزاي تتأكد إن القياس بيقارن نفس العمل على الاتنين؟

## atomic.Int64
استورد ~sync/atomic~ باسم ~atomic~. ~atomic.Int64~ رقم بصفر جاهز؛ ~Add~ تضيف وترجع القيمة الجديدة، و ~Load~ ترجع قراءة متزامنة. كل الوصولات لنفس الرقم لازم تبقى ذرّية، وممنوع نسخ الكائن بعد استخدامه.
~~~go
var pulses atomic.Int64
next := pulses.Add(3)
current := pulses.Load()
fmt.Println(next, current) // 3 3
~~~
العمليات الذرّية ليها ترتيب متسلسل متسق. لو قراءة لاحظت نتيجة عملية ذرّية، العملية دي بتسبقها بعلاقة ~happens-before~. ده يحمي رقم واحد؛ مش بيحوّل كذا خطوة لمعاملة واحدة.

## متستبدلش invariant بـ Load ثم Add
~~~go
if remaining.Load() >= requested {
 remaining.Add(-requested) // غلط: حد تاني ممكن يكون سحب بين الخطوتين
}
~~~
العدّاد في المهمة ملوش قيد سحب ولا حقول مترابطة، فـ ~Add~ واحدة كفاية. قارن ده بحساب الوقود اللي محتاج قرار كامل بقفل.

## عدّادين بنفس العقد
~LockedPulses~ تستخدم ~sync.Mutex~، و ~AtomicPulses~ تستخدم ~atomic.Int64~. الاتنين بصفر جاهز وطرق pointer، و ~Add(delta int64) int64~ ترجع القيمة الجديدة، و ~Total() int64~ تقرا الرقم. الصفر والسالب صالحين.
~~~go
var locked LockedPulses
var direct AtomicPulses
fmt.Println(locked.Add(2), direct.Add(2))
fmt.Println(locked.Total(), direct.Total())
~~~
التنفيذ والـ benchmark في ~main_test.go~ عشان ده مختبر مقارنة، مش برنامج نشر. ملف ~main.go~ فاضي وللقراءة بس. متحطش التنفيذ في الملف المقفول.

## نفس الحمل المتوازي
~b.RunParallel~ توزع التكرارات على مشاركين، وكل ~pb.Next()~ تصريح لتكرار واحد. العدد الإجمالي عبر المشاركين يساوي ~b.N~، مش ~b.N~ لكل واحد. جهّز عدّاد واحد مشترك خارج callback.
~~~go
func BenchmarkDemo(b *testing.B) {
 var counter AtomicPulses
 b.RunParallel(func(pb *testing.PB) {
  for pb.Next() { counter.Add(1) }
 })
 if counter.Total() != int64(b.N) { b.Fatal("lost work") }
}
~~~
~testing.PB~ حالة تكرار للقياس المتوازي. ~b.Fatal~ يفشل القياس لو الشغل ناقص، بدل قبول سرعة وهمية. اكتب benchmark لكل عدّاد بنفس النمط، من غير إنشاء goroutine جديدة لكل نبضة.

## القياس والنجوم
المختبر يقارن ~BenchmarkHiddenNaive~ للـ mutex مع ~BenchmarkHiddenTheirs~ للـ atomic بنفس الحمل المتوازي، بعد اختبارات الصحة. بيشغّل ~-race~، فالأرقام متأثرة بتكلفة detector والجهاز والمترجم والحمل؛ مفيش وعد إن atomic دايماً أسرع. تقدر تقيس محلياً من غير detector لفهم الأداء الفعلي، لكن متخلطش القياس بإثبات الصحة.

حد النجاح ميزانية تراجع واسعة: نسبة ~mutex/atomic~ على الأقل 0.1، نجمتين عند 0.5، وتلاتة عند 1. النسبة مقارنة محلية، مش ضمان أداء لكل الأجهزة.

## أخطاء شائعة
~~~go
b.RunParallel(func(pb *testing.PB) {
 var counter AtomicPulses // غلط: عدّاد مستقل لكل مشارك يخفي التنافس
 for pb.Next() { counter.Add(1) }
})
~~~

## المهمة
في ~main_test.go~ نفّذ ~LockedPulses~ بـ ~sync.Mutex~ و ~AtomicPulses~ بـ ~atomic.Int64~ وطرق ~Add(delta int64) int64~ و ~Total() int64~. الصفر جاهز؛ القيم ممكن تزيد أو تقل أو تثبت، والمدخلات لا تسبب overflow. كل ~Add~ عملية واحدة وترجع قيمتها الجديدة، والقراءة آمنة مع الكتابة. الكائنات مستقلة وممنوع نسخها بعد الاستخدام.

اكتب ~BenchmarkLockedPulses~ و ~BenchmarkAtomicPulses~ باستخدام ~b.RunParallel~: عدّاد واحد مشترك لكل benchmark وإضافة 1 لكل ~pb.Next~ وفحص المجموع بعد رجوع ~RunParallel~. اكتب تستات للقيم المتسلسلة والمتزامنة. التستات المخفية تفحص الاتنين من غير أحكام زمنية، ثم ~Lab.Bench~ تقارن الحمل نفسه بحدود النجوم الموضحة. مفيش ساعة تطبيق أو نوم أو شبكة أو ~Report~؛ توقيت benchmark مسؤولية أداة Go.`,
		Hints: []string{"حدد سلوك العداد قبل مقارنة آلية الحماية.", "خلي الحمل متساويًا بين النسختين.", "راجع طرق ~Add~ و~Load~ في ~atomic.Int64~."},
		En: &LevelText{LessonRevision: 1, Steps: ch25AtomicEnglish, Title: "Atomic or Mutex", Goal: "Compare two correct concurrent pulse counters under identical benchmark workloads",
			Files: []LabFile{{Name: "main.go", Starter: ch256Main, ReadOnly: true}, {Name: "main_test.go", Starter: ch256Starter}},
			Lesson: `# Choose After Measuring

## Before you write
- Is this an independent number or an invariant spanning multiple fields?
- How will you ensure both benchmarks perform identical work?

## atomic.Int64
Import ~sync/atomic~ as ~atomic~. A zero ~atomic.Int64~ is ready: ~Add~ increments and returns the new value; ~Load~ reads synchronously. Every access to its storage must be atomic. Never copy it after use.
~~~go
var pulses atomic.Int64
next := pulses.Add(3)
current := pulses.Load()
fmt.Println(next, current) // 3 3
~~~
Atomic operations have a consistent sequential order. When one observes another's effect, the observed operation establishes ~happens-before~. This protects one number, without turning several operations into a transaction.

## Do not replace an invariant with Load then Add
~~~go
if remaining.Load() >= requested {
 remaining.Add(-requested) // Another caller may spend between these operations
}
~~~
Our counter has no affordability constraint or related fields, so one ~Add~ suffices. Compare with the fuel account's complete locked decision.

## Two counters, one contract
~LockedPulses~ uses ~sync.Mutex~; ~AtomicPulses~ uses ~atomic.Int64~. Both are zero-ready with pointer methods ~Add(delta int64) int64~ returning the new value and ~Total() int64~ reading it. Zero and negative deltas are valid.
~~~go
var locked LockedPulses
var direct AtomicPulses
fmt.Println(locked.Add(2), direct.Add(2))
fmt.Println(locked.Total(), direct.Total())
~~~
Implementations and benchmarks live in ~main_test.go~ because this is a comparison laboratory. ~main.go~ is empty and read-only; do not move implementation into that file.

## Identical parallel work
~b.RunParallel~ distributes iterations among participants. Each ~pb.Next()~ permits one iteration. Total iterations equal ~b.N~ across all participants, rather than per participant. Create one shared counter outside the callback.
~~~go
func BenchmarkDemo(b *testing.B) {
 var counter AtomicPulses
 b.RunParallel(func(pb *testing.PB) {
  for pb.Next() { counter.Add(1) }
 })
 if counter.Total() != int64(b.N) { b.Fatal("lost work") }
}
~~~
~testing.PB~ holds parallel iteration state. ~b.Fatal~ rejects missing work instead of rewarding false speed. Write the same benchmark for each counter, without launching a goroutine per pulse.

## Measurement and stars
After correctness tests, the lab compares mutex ~BenchmarkHiddenNaive~ with atomic ~BenchmarkHiddenTheirs~ under the same parallel workload. Both run with ~-race~: instrumentation, hardware, compiler, and load affect results. Atomic is not universally faster. Local measurements without instrumentation can explore production performance, but timing is not correctness evidence.

Success uses a broad regression budget: ~mutex/atomic~ ratio at least 0.1; two stars at 0.5 and three at 1. This is a local comparison rather than a performance guarantee across machines.

## Common mistakes
~~~go
b.RunParallel(func(pb *testing.PB) {
 var counter AtomicPulses // Private per participant: hides shared contention
 for pb.Next() { counter.Add(1) }
})
~~~

## Your mission
In ~main_test.go~, implement ~LockedPulses~ with ~sync.Mutex~, ~AtomicPulses~ with ~atomic.Int64~, and methods ~Add(delta int64) int64~ and ~Total() int64~. Zero values work; deltas may increase, decrease, or preserve totals, without overflow. Each ~Add~ is indivisible and returns its new value. Reads synchronize with writes. Instances are independent and must not be copied after use.

Write ~BenchmarkLockedPulses~ and ~BenchmarkAtomicPulses~ using ~b.RunParallel~: one shared counter per benchmark, adding 1 per ~pb.Next~, then checking the total after ~RunParallel~ returns. Write sequential and concurrent correctness tests. Hidden tests check both counters without timing assertions; ~Lab.Bench~ then compares identical workloads using the stated star tiers. No application clock, sleeping, network, or ~Report~; Go's benchmark tool owns measurement time.`,
			Hints: []string{"Define counter behavior before comparing synchronization.", "Keep the workload identical between versions.", "Review ~Add~ and ~Load~ on ~atomic.Int64~."}},
		Lab: &Lab{ReqFile: "main_test.go", Hidden: ch256Hidden, Bench: &Bench{Tiers: [3]float64{0.1, 0.5, 1}}},
	})
}

const ch256Main = `package main

func main() {}
`
const ch256Starter = `package main

import "testing"

type LockedPulses struct{}
type AtomicPulses struct{}

func (p *LockedPulses) Add(delta int64) int64 { return 0 }
func (p *LockedPulses) Total() int64          { return 0 }
func (p *AtomicPulses) Add(delta int64) int64 { return 0 }
func (p *AtomicPulses) Total() int64          { return 0 }
func TestPulses(t *testing.T) {
	var p AtomicPulses
	if p.Total() != 0 {
		t.Fatal("zero")
	}
}
`
const ch256Hidden = `package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"slices"
	"sync"
	"testing"
)

type pulseCounter interface {
	Add(int64) int64
	Total() int64
}

// Cases: zero, signed and zero deltas, independent instances, concurrent Add returns and read/write races.
func TestHiddenPulsesSerialAndConcurrent(t *testing.T) {
	for _, factory := range []func() pulseCounter{func() pulseCounter { return &LockedPulses{} }, func() pulseCounter { return &AtomicPulses{} }} {
		rng := rand.New(rand.NewSource(256))
		p := factory()
		other := factory()
		want := int64(0)
		if p.Total() != 0 {
			t.Fatal("zero")
		}
		for i := 0; i < 200; i++ {
			n := int64(rng.Intn(101) - 50)
			want += n
			if p.Add(n) != want || p.Total() != want {
				t.Fatal("signed delta")
			}
		}
		if other.Total() != 0 {
			t.Fatal("shared instance state")
		}
		p = factory()
		workers := 16 + rng.Intn(16)
		each := 40
		returns := make(chan int64, workers*each)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				for j := 0; j < each; j++ {
					returns <- p.Add(1)
					_ = p.Total()
				}
			}()
		}
		close(start)
		wg.Wait()
		close(returns)
		var got []int64
		for n := range returns {
			got = append(got, n)
		}
		slices.Sort(got)
		if p.Total() != int64(workers*each) {
			t.Fatal("lost increments")
		}
		for i, n := range got {
			if n != int64(i+1) {
				t.Fatal("Add must return its own new total", n, i)
			}
		}
	}
}
func TestHiddenPulsesBenchmarksPresent(t *testing.T) {
	src, e := os.ReadFile("main_test.go")
	if e != nil {
		t.Fatal(e)
	}
	f, e := parser.ParseFile(token.NewFileSet(), "main_test.go", src, 0)
	if e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"BenchmarkLockedPulses", "BenchmarkAtomicPulses"} {
		found := false
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != name {
				continue
			}
			parallel, next := false, false
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if ok {
					parallel = parallel || sel.Sel.Name == "RunParallel"
					next = next || sel.Sel.Name == "Next"
				}
				return true
			})
			found = parallel && next
		}
		if !found {
			t.Error("write the parallel benchmark", name)
		}
	}
}
func BenchmarkHiddenNaive(b *testing.B) {
	var p LockedPulses
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Add(1)
		}
	})
	if p.Total() != int64(b.N) {
		b.Fatal("mutex lost work")
	}
}
func BenchmarkHiddenTheirs(b *testing.B) {
	var p AtomicPulses
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			p.Add(1)
		}
	})
	if p.Total() != int64(b.N) {
		b.Fatal("atomic lost work")
	}
}
`
