package levels

const ch24HiddenSupport = `package main

import (
	"math/rand"
	"testing"
	"time"
)

func supplyRand() *rand.Rand { return rand.New(rand.NewSource(24017)) }
func supplyReceive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("deadlock watchdog: missing supply event")
		var zero T
		return zero
	}
}
`

const ch24HTTPHidden = `package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

// net.Pipe keeps httptest's real HTTP stack usable in socket-restricted sandboxes.
type supplyListener struct {
	connections chan net.Conn
	closed      chan struct{}
	once        sync.Once
}

func (l *supplyListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.connections:
		return c, nil
	case <-l.closed:
		return nil, net.ErrClosed
	}
}
func (l *supplyListener) Close() error   { l.once.Do(func() { close(l.closed) }); return nil }
func (l *supplyListener) Addr() net.Addr { return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)} }
func supplyServer(t *testing.T, h http.Handler) *httptest.Server {
	t.Helper()
	listener := &supplyListener{connections: make(chan net.Conn), closed: make(chan struct{})}
	server := &httptest.Server{Listener: listener, Config: &http.Server{Handler: h}}
	server.Start()
	transport := server.Client().Transport.(*http.Transport)
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		client, remote := net.Pipe()
		select {
		case listener.connections <- remote:
			return client, nil
		case <-ctx.Done():
			client.Close()
			remote.Close()
			return nil, ctx.Err()
		case <-listener.closed:
			client.Close()
			remote.Close()
			return nil, net.ErrClosed
		}
	}
	t.Cleanup(func() { server.Client().CloseIdleConnections(); server.Close() })
	return server
}
`

func init() {
	add(&Level{ID: "24-1-ping-pong", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24PingArabic, Title: "بينج بونج", Goal: "خلّي درونين يسلّموا تصريح شحنة بالتبادل ويقفوا من غير تسريب",
		Requires: []string{"chan", "close"}, Bridge: []int{32}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch241StarterAr}},
		Lesson: `# تصريح الشحنة بين درونين

## قبل ما تكتب
- مين يملك القناة ومين يقفلها بعد آخر إرسال؟
- آخر رد هيوصل إزاي من غير إرسال زيادة يعلّق الطرفين؟

## تسليم من غير مخزن
القناة من غير سعة بتعمل مقابلة بين المرسل والمستقبل. إرسال التصريح معناه إن الطرف التاني استلمه، مش إنه خلص كل شغله.
~~~go
permit := make(chan int)
go func() { permit <- 9 }()
code := <-permit
fmt.Println(code)
~~~
هنمثّل الدرونين بـ callback محقونة عشان الاختبار يشوف ترتيب التسليم من غير حركة أو وقت حقيقي. كل جولة فيها تسليم للدرون 0 وبعده للدرون 1.

## الملكية والقفل
~close~ معناها مفيش إرسال تاني. صاحب القناة يقفل مرة واحدة بعد انتهاء كل المرسلين، مش لمجرد إنه عايز يوقف مستقبل. مع أكتر من مرسل، استخدم ~sync.WaitGroup~ من الفصل 10 قبل القفل.
~~~go
close(permit)
value, ok := <-permit
fmt.Println(value, ok) // 0 false لما القناة فاضية ومقفولة
~~~
الـ ~ok~ هو اللي بيفرّق نهاية القناة عن قيمة 0 حقيقية. ~range~ على قناة بيخلص بعد استهلاك كل القيم وقفلها.

## طلب ورد
درون الإرسال يبعت رقم الجولة، ودرون الاستقبال يرجع إقرار على قناة تانية. استنى الإقرار قبل الجولة الجديدة. كده الـ callback التانية تخلص قبل البداية اللي بعدها؛ آخر جولة ليها رد واحد بس.
~~~go
select {
case permit <- code:
case <-ctx.Done():
 return ctx.Err()
}
~~~
نفس ~select~ مطلوبة حوالين أي إرسال أو استقبال ممكن يستنى. افحص ~ctx.Err()~ قبل callback؛ اختيار ~select~ بين حالتين جاهزين مش أولوية للإلغاء. عند خطأ callback، ألغي السياق الداخلي واستنى الاتنين قبل الرجوع. الـ callback نفسها بتخلص ومش بتعمل panic.

## اختبار العمر
~labkit.LeakCheck~ بتسجّل فحص في نهاية التست. استخدم إشارات قنوات لتحديد الإلغاء؛ مفيش ~time.Sleep~ في منطق المهمة.
~~~go
func TestRelay(t *testing.T) {
 labkit.LeakCheck(t)
}
~~~

## أخطاء شائعة
~~~go
close(permit)
permit <- 4 // panic: إرسال بعد القفل
// إرسال رد زيادة بعد آخر استقبال هيعمل deadlock
~~~

## المهمة
نفّذ ~RelaySupplies(ctx context.Context, rounds int, visit func(drone, round int) error) error~. السياق غير ~nil~. ارفض ~rounds < 0~ أو callback ~nil~ قبل تشغيل أي شغل. بعدها الإلغاء المسبق يرجع ~ctx.Err()~ حتى مع 0. صفر جولات ناجح من غير callbacks أو goroutines.

في النجاح، شغّل مشاركين للدرونين 0 و 1. استخدم قنوات unbuffered لتسليم الجولة والإقرار. نادِ ~visit(0,r)~ وبعدها ~visit(1,r)~ لكل ~r~ من 0 لحد ~rounds-1~ بالترتيب ومن غير تداخل callbacks. استخدم ~close~ بعد انتهاء المرسلين، واستنى كل goroutine قبل الرجوع. أي خطأ callback يرجع بنفس السبب ويوقف الجولات؛ الإلغاء يرجع خطأ السياق، ومفيش callback جديدة بعد ملاحظة الإلغاء. ممنوع وقت حقيقي أو ~Report~. ~main~ فاضية؛ دي محاكاة لتسليم درونين بحقن نقاط التفاعل عشان نختبر العمر والقفل.`,
		Hints: []string{"فرّق بين تسليم الرسالة وانتهاء شغل المستقبل.", "ارسم اتجاه الإرسال والرد قبل مراجعة goroutines.", "راجع استخدام ~close~ وقراءة ~ok~ عند نهاية القناة."},
		En: &LevelText{LessonRevision: 1, Steps: ch24PingEnglish, Title: "Ping Pong", Goal: "Have two drones hand over a cargo permit in alternating rounds and stop without leaks", Files: []LabFile{{Name: "main.go", Starter: ch241Starter}},
			Lesson: `# A Cargo Permit Between Two Drones

## Before you write
- Who owns each channel and closes it after the last sender?
- How will the final acknowledgement arrive without an extra blocking send?

## An unbuffered handover
An unbuffered channel pairs a sender with a receiver. Sending a permit proves receipt, rather than completion of the receiver's work.
~~~go
permit := make(chan int)
go func() { permit <- 9 }()
code := <-permit
fmt.Println(code)
~~~
Injected callbacks represent drone actions so tests observe handover order without movement or real time. Each round visits drone 0 then drone 1.

## Ownership and closure
~close~ means no further sends. The owner closes exactly once after all senders finish. Closing is not a way to stop a receiver while senders remain. With multiple senders, join using ~sync.WaitGroup~ from chapter 10 before closing.
~~~go
close(permit)
value, ok := <-permit
fmt.Println(value, ok) // 0 false when closed and drained
~~~
~ok~ distinguishes channel termination from a real zero. A channel ~range~ ends after closure and draining.

## Request and acknowledgement
The sending drone passes the round number; the receiving drone acknowledges on a second channel. ~Wait~ for that acknowledgement before the next round, so the second callback finishes first. The final round needs only one reply.
~~~go
select {
case permit <- code:
case <-ctx.Done():
 return ctx.Err()
}
~~~
Wrap every potentially blocking send and receive this way. Check ~ctx.Err()~ before callbacks; ~select~ does not prioritize cancellation when both cases are ready. A callback error cancels an internal context; join both participants before returning. Callbacks terminate and do not panic.

## Testing lifetimes
~labkit.LeakCheck~ registers an end-of-test leak check. Use channel signals to arrange cancellation; the implementation needs no ~time.Sleep~.
~~~go
func TestRelay(t *testing.T) {
 labkit.LeakCheck(t)
}
~~~

## Common mistakes
~~~go
close(permit)
permit <- 4 // panic: send after close
// An extra reply after the final receive deadlocks
~~~

## Your mission
Implement ~RelaySupplies(ctx context.Context, rounds int, visit func(drone, round int) error) error~. Context is non-~nil~. Reject negative ~rounds~ or a ~nil~ callback before work. Then pre-cancellation returns ~ctx.Err()~, even for zero. Zero ~rounds~ succeeds without callbacks or goroutines.

For success, run two participants representing drones 0 and 1. Use unbuffered channels for round handover and acknowledgement. Call ~visit(0,r)~ then ~visit(1,r)~ for every ~r~ from 0 through ~rounds-1~, in order with no overlapping callbacks. Use ~close~ after senders finish and join every goroutine before returning. Callback failure preserves its cause and stops further ~rounds~. Cancellation returns the context error; no new callback starts after cancellation is observed. No real time or ~Report~. Leave ~main~ empty; this lab simulates two drones through injected interaction points so lifetime and closure can be tested.`,
			Hints: []string{"Separate message handoff from completion of receiver work.", "Draw send and reply directions before reviewing goroutines.", "Review ~close~ and the ~ok~ receive result at channel end."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch241Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport}},
	})
}

const ch241Starter = `package main

import (
	"context"
	"errors"
)

func RelaySupplies(ctx context.Context, rounds int, visit func(drone, round int) error) error {
	return errors.New("unfinished")
}
func main() {}
`
const ch241Hidden = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
	"sync/atomic"
	"testing"
)

// Boundaries: negative, zero, one, odd/even rounds, cancellation at either drone, callback errors.
func TestHiddenRelayOrderedRounds(t *testing.T) {
	labkit.LeakCheck(t)
	rng := supplyRand()
	counts := []int{0, 1, 2, 3}
	for i := 0; i < 25; i++ {
		counts = append(counts, rng.Intn(30))
	}
	for _, n := range counts {
		var busy atomic.Int32
		next := 0
		err := RelaySupplies(context.Background(), n, func(d, r int) error {
			if busy.Add(1) != 1 {
				t.Error("overlapping drone actions")
			}
			if d != next%2 || r != next/2 {
				t.Errorf("handover %d: drone=%d round=%d", next, d, r)
			}
			next++
			busy.Add(-1)
			return nil
		})
		if err != nil || next != 2*n {
			t.Fatal("round count", next, err)
		}
	}
}
func TestHiddenRelayStopsAndJoins(t *testing.T) {
	labkit.LeakCheck(t)
	fault := errors.New("permit jam")
	for _, drone := range []int{0, 1} {
		for _, cancelInstead := range []bool{false, true} {
			ctx, cancel := context.WithCancel(context.Background())
			calls := 0
			err := RelaySupplies(ctx, 10, func(d, r int) error {
				calls++
				if d == drone && r == 2 {
					if cancelInstead {
						cancel()
						return nil
					}
					return fault
				}
				return nil
			})
			cancel()
			want := fault
			if cancelInstead {
				want = context.Canceled
			}
			if !errors.Is(err, want) || calls != 5+drone {
				t.Fatal("stop boundary", calls, err)
			}
		}
	}
	called := false
	fn := func(int, int) error { called = true; return nil }
	if RelaySupplies(context.Background(), -1, fn) == nil || called {
		t.Fatal("negative")
	}
	if RelaySupplies(context.Background(), 0, nil) == nil {
		t.Fatal("nil callback")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := RelaySupplies(ctx, 0, fn); !errors.Is(err, context.Canceled) || called {
		t.Fatal("pre-cancel", err)
	}
	if err := RelaySupplies(context.Background(), 0, fn); err != nil || called {
		t.Fatal("zero", err)
	}
}
`

func init() {
	add(&Level{ID: "24-2-producers", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24DepotArabic, Title: "منتجين ومستهلكين", Goal: "اجمع صناديق مخازن متوازية في مستهلك واحد بطابور محدود",
		Requires: []string{"bufchan", "close"}, Bridge: []int{37}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch242StarterAr}},
		Lesson: `# مخازن بتغذّي رصيف واحد

## قبل ما تكتب
- لو مفيش منتجين، مين هيقفل طابور الصناديق؟
- لو المستهلك فشل، إزاي تفك المنتج اللي واقف في إرسال؟

## سعة محدودة
الـ buffered channel تخزن لحد السعة المحددة. لما تتملي، المرسل يستنى: دي بداية مفهوم الضغط العكسي، بدل ذاكرة بتكبر مع كل صندوق.
~~~go
crates := make(chan int, 2)
crates <- 8
crates <- 5
fmt.Println(<-crates)
~~~
السعة 0 شغالة زي القناة unbuffered. السعة السالبة تعمل panic، فافحصها قبل ~make~. السعة مش عدد المنتجين ولا عدد goroutines.

## كذا مرسل ومالك واحد
كل مخزن يبقى goroutine تقرأ slice بتاعته بس. خلي ~sync.WaitGroup~ تعد المنتجين، وgoroutine منفصلة تستناهم وتقفل القناة. ده بيشتغل حتى مع صفر منتجين.
~~~go
go func() {
 producers.Wait()
 close(crates)
}()
for code := range crates { fmt.Println(code) }
~~~
استعمل نفس نمط ~Add~ قبل ~go~ من الفصل 10. المنتج مبيقفلش القناة المشتركة عشان ممكن مخزن تاني لسه بيبعت.

## الإلغاء مسار خروج
اشتق سياق داخلي من ~ctx~ بـ ~context.WithCancel~. لو callback المستهلك فشلت، احتفظ بخطئها وألغي. كل إرسال يعمل ~select~ مع ~Done~. بعد الخطأ، كمّل تصريف القناة من غير callbacks لحد القفل عشان كل المنتجين ينتهوا قبل الرجوع.
~~~go
run, cancel := context.WithCancel(ctx)
defer cancel()
select {
case crates <- 8:
case <-run.Done():
 return
}
~~~
القيم المتساوية صناديق منفصلة: ممنوع إزالة التكرارات. الترتيب بين مخزنين مش مضمون، لكن ترتيب الصناديق داخل كل مخزن محفوظ.

## أخطاء شائعة
~~~go
go func() { defer close(crates); crates <- 8 }()
go func() { crates <- 5 }() // ممكن يبعت بعد القفل
~~~

## المهمة
نفّذ ~UnloadDepots(ctx context.Context, batches [][]int, capacity int, consume func(int) error) error~. السياق غير ~nil~. ارفض سعة سالبة أو callback ~nil~ قبل أي callback. بعدها الإلغاء المسبق يرجع خطأ السياق. شغّل منتج لكل batch غير فاضية، وقناة واحدة بسعة ~capacity~، ومستهلك واحد ينادي ~consume~ بالتتابع.

النجاح يستهلك كل عنصر مرة واحدة بما فيها القيم السالبة والتكرارات، ويحفظ ترتيب كل batch. صفر ~batches~ أو ~batches~ فاضية نجاح من غير callbacks. متغيرش المدخلات. خطأ المستهلك يوقف callbacks التالية ويرجع سببه. الإلغاء يوقف عند ملاحظته ويرجع ~ctx.Err()~؛ لا ضمان لعدد جزئي ثابت وقت السباق مع الإلغاء. استنى كل المنتجين والقافل قبل الرجوع؛ القناة تتقفل مرة واحدة بعد انتهاء المرسلين. ممنوع وقت حقيقي أو ~Report~.`,
		Hints: []string{"افصل مساحة الانتظار عن عدد الرسائل الحالية.", "حدد مين يملك كتابة نتيجة الاستهلاك.", "راجع مقارنة ~cap~ و~len~ في المثال قبل تتبع queue."},
		En: &LevelText{LessonRevision: 1, Steps: ch24DepotEnglish, Title: "Producers and Consumers", Goal: "Merge concurrent depot crates into one consumer through a bounded queue", Files: []LabFile{{Name: "main.go", Starter: ch242Starter}},
			Lesson: `# Many Depots, One Dock

## Before you write
- With no producers, who closes the crate queue?
- If the consumer fails, how will a blocked producer escape its send?

## Bounded ~capacity~
A buffered channel holds up to its declared ~capacity~. Once full, sends wait. This is backpressure: the queue cannot grow indefinitely.
~~~go
crates := make(chan int, 2)
crates <- 8
crates <- 5
fmt.Println(<-crates)
~~~
Capacity zero behaves as an unbuffered channel. Negative ~capacity~ panics, so validate before ~make~. Capacity is neither producer count nor goroutine count.

## Several senders, one owner
Each depot goroutine only reads its own slice. A ~sync.WaitGroup~ counts producers; a separate goroutine waits for them and closes the channel. This works for zero producers too.
~~~go
go func() {
 producers.Wait()
 close(crates)
}()
for code := range crates { fmt.Println(code) }
~~~
Reuse chapter 10's ~Add~ before ~go~ pattern. An individual producer must not close a shared channel while another depot may send.

## Cancellation is an exit path
Derive an internal context with ~context.WithCancel~. On consumer failure, ~save~ its error and cancel. Every send selects against ~Done~. After failure, drain without callbacks until closure, joining all producers before returning.
~~~go
run, cancel := context.WithCancel(ctx)
defer cancel()
select {
case crates <- 8:
case <-run.Done():
 return
}
~~~
Equal values represent distinct crates: do not deduplicate. Inter-depot order is unspecified; each depot's own crate order is preserved.

## Common mistakes
~~~go
go func() { defer close(crates); crates <- 8 }()
go func() { crates <- 5 }() // May send after close
~~~

## Your mission
Implement ~UnloadDepots(ctx context.Context, batches [][]int, capacity int, consume func(int) error) error~. Context is non-~nil~. Reject negative ~capacity~ or a ~nil~ callback before callbacks. Then pre-cancellation returns the context error. Start one producer per nonempty batch, one channel of ~capacity~, and one consumer calling ~consume~ serially.

Success consumes each element once, including negative values and duplicates, preserving each batch's order. No ~batches~ or only empty ~batches~ succeeds without callbacks. Do not mutate ~inputs~. Consumer failure stops subsequent callbacks and preserves its cause. Cancellation stops when observed and returns ~ctx.Err()~; there is no fixed partial count during a cancellation race. Join all producers and the closer before returning; close once after all senders finish. No real time or ~Report~.`,
			Hints: []string{"Separate waiting space from the current message count.", "Identify who owns the consumed result.", "Review ~cap~ and ~len~ in the example before tracing the queue."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch242Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport}},
	})
}

const ch242Starter = `package main

import (
	"context"
	"errors"
)

func UnloadDepots(ctx context.Context, batches [][]int, capacity int, consume func(int) error) error {
	return errors.New("unfinished")
}
func main() {}
`
const ch242Hidden = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
	"reflect"
	"slices"
	"testing"
)

// Cases: no producers, empty producers, zero/one/many buffer slots, repeated and signed crates.
func TestHiddenDepotsNoLostCrates(t *testing.T) {
	labkit.LeakCheck(t)
	rng := supplyRand()
	for _, capacity := range []int{0, 1, 7} {
		for count := 0; count < 12; count++ {
			batches := make([][]int, count)
			before := make([][]int, count)
			want := map[int]int{}
			for i := range batches {
				for j, size := 0, rng.Intn(30); j < size; j++ {
					v := rng.Intn(9) - 4
					batches[i] = append(batches[i], v)
					want[v]++
				}
				before[i] = slices.Clone(batches[i])
			}
			got := map[int]int{}
			if err := UnloadDepots(context.Background(), batches, capacity, func(v int) error { got[v]++; return nil }); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) || !reflect.DeepEqual(batches, before) {
				t.Fatal("lost/duplicated crates or mutated inputs", got, want)
			}
		}
	}
	for _, capacity := range []int{0, 1, 8} {
		var got []int
		if err := UnloadDepots(context.Background(), [][]int{{8, 3, -1, 8}}, capacity, func(v int) error { got = append(got, v); return nil }); err != nil || !slices.Equal(got, []int{8, 3, -1, 8}) {
			t.Fatal("producer order", got, err)
		}
	}
}
func TestHiddenDepotsStopBlockedProducers(t *testing.T) {
	labkit.LeakCheck(t)
	fault := errors.New("dock unavailable")
	for _, capacity := range []int{0, 1, 5} {
		for _, stop := range []bool{false, true} {
			ctx, cancel := context.WithCancel(context.Background())
			batches := make([][]int, 15)
			for i := range batches {
				batches[i] = make([]int, 100)
			}
			calls := 0
			err := UnloadDepots(ctx, batches, capacity, func(int) error {
				calls++
				if stop {
					cancel()
					return nil
				}
				return fault
			})
			cancel()
			want := fault
			if stop {
				want = context.Canceled
			}
			if calls != 1 || !errors.Is(err, want) {
				t.Fatal("consumer stop", calls, err)
			}
		}
	}
	calls := 0
	fn := func(int) error { calls++; return nil }
	if UnloadDepots(context.Background(), nil, -1, fn) == nil || UnloadDepots(context.Background(), nil, 0, nil) == nil || calls != 0 {
		t.Fatal("preflight")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := UnloadDepots(ctx, [][]int{{1}}, 1, fn); !errors.Is(err, context.Canceled) || calls != 0 {
		t.Fatal("pre-cancel", err)
	}
}
`

func init() {
	add(&Level{ID: "24-3-work-crew", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24CrewArabic, Title: "طاقم الشغل", Goal: "ابني طاقم عام بعدد عمّال ثابت ونتيجة أو خطأ لكل شحنة",
		Requires: []string{"chan", "go", "sel:context.Context"}, Bridge: []int{34}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch243StarterAr}},
		Lesson: `# طاقم الشحن تحت الضغط

## قبل ما تكتب
- مين يكتب نتيجة كل موضع من غير ~append~ متوازي؟
- لو أول عامل واقف، هل هتشغّل آلاف goroutines للباقي؟

## عدد ثابت من العمال
الـ worker pool تشغّل عدد ثابت من goroutines تقرأ طابور مشترك. الضغط العكسي معناه إن إرسال مهمة جديدة يستنى مساحة أو عامل، بدل إنشاء goroutine لكل عنصر.
~~~go
jobs := make(chan int)
for i := 0; i < 3; i++ {
 go func() { for index := range jobs { fmt.Println(index) } }()
}
~~~
القناة unbuffered هنا تمنع تكديس الشغل. في المهمة، العامل ينادي callback بنفسه؛ متعملش ~go~ جديدة لكل callback. خلي المرسل يقفل الطابور، واستنى العمال بـ ~sync.WaitGroup~.

## دالة عامة ونتائج مستقلة
الـ generics من الفصل 14 تنفع للدوال والنتائج مع بعض. ~any~ قيد يسمح بكل نوع، مش تحويل للقيم لـ ~interface{}~.
~~~go
type Reply[R any] struct { Value R; Err error }
func identity[T any](x T) T { return x }
fmt.Println(identity("bay"))
~~~
جهّز slice بطول المدخلات وابعت indexes. كل عامل يكتب موضع مختلف. الـ ~Wait~ قبل قراءة النتائج يضمن إن الكتابة خلصت، حتى لو النهاية بترتيب مختلف.
~~~go
replies := make([]Reply[string], 4)
replies[2] = Reply[string]{Value: "ready"}
~~~
خطأ عنصر مش خطأ batch: خزّن الخطأ وخلي العامل ياخد المهمة التالية. متحوّلش كل النتائج لفشل بسبب شحنة واحدة.

## الإلغاء والتوفيق بين النتائج
مرّر ~context.Context~ نفسه للـ callback. الـ callback لازم تحترم الإلغاء وتنتهي؛ مفيش طريقة آمنة لقتل goroutine بالقوة. افحص الإلغاء قبل سحب شغل وقبل callback، وكل انتظار على قناة يبقى قابل للإلغاء.
~~~go
select {
case jobs <- index:
case <-ctx.Done():
 // وقف التغذية، اقفل jobs، واستنى العمال
}
~~~
بعد انتظار العمال، المواضع اللي متبدأتش تاخد صفر النوع وخطأ السياق. الـ callback اللي بدأت تحتفظ بالقيمة والخطأ اللي رجعتهم حتى لو الإلغاء حصل بعدها.

## قيم عامة في الطاقم
استخدم نفس النوع لمدخلات مختلفة، وخلّي خطأ كل عنصر جنب قيمته.
~~~go
reply := CrewReply[string]{Value: "oxygen", Err: nil}
fmt.Println(reply.Value, reply.Err)
~~~

## أخطاء شائعة
~~~go
for _, cargo := range cargos {
 go inspect(cargo) // عدد غير محدود، مفيش backpressure
}
// append مشترك من العمال هيعمل race
~~~

## المهمة
نفّذ ~RunCrew[T,R any](ctx context.Context, inputs []T, workers int, work func(context.Context,T)(R,error)) ([]CrewReply[R],error)~. ~CrewReply[R]~ فيها ~Value R~ و ~Err error~. السياق غير ~nil~ والـ callback بتنتهي. ارفض ~workers <= 0~ أو ~work == nil~ بـ ~nil~ وخطأ قبل أي callback، حتى مع مدخلات فاضية.

بعد التحقق، جهّز نتيجة لكل مدخل بنفس الترتيب. لو السياق متلغي، املأ النتائج بخطئه من غير callbacks وارجع خطأ السياق للـ batch. غير كده الفاضي نجاح. استخدم طابور indexes محدود و ~workers~ عمال كحد أقصى، مع السماح بعدد أكبر من المدخلات. كل مدخل يبدأ مرة واحدة في النجاح، وكل callback ترجع قيمة وخطأ يتخزنوا كما هم؛ خطأ عنصر ميوقفش الباقي. متغيرش المدخلات. عند الإلغاء، أوقف الجديد عند ملاحظته واستنى الجاري؛ المواضع اللي متبدأتش تاخد صفر النوع وخطأ السياق. ارجع خطأ السياق للـ batch عند الإلغاء، أو ~nil~ في النجاح حتى مع أخطاء عناصر. مفيش تسريب أو وقت حقيقي أو ~Report~.

النجمة التالتة: نفس العقد مع structs ونصوص بدل أرقام، واستدعاءات مستقلة متوازية من غير حالة عامة مشتركة.`,
		Hints: []string{"ميّز الخطأ العام عن خطأ عنصر واحد.", "اربط كل مدخل بخانة نتيجة مستقلة.", "ارجع لتعريف ~CrewReply[R]~ قبل تتبع قيم ~R~."},
		En: &LevelText{LessonRevision: 1, Steps: ch24CrewEnglish, Title: "Work Crew", Goal: "Build a generic fixed-size crew with a result or error for each shipment", Files: []LabFile{{Name: "main.go", Starter: ch243Starter}},
			Lesson: `# A Crew Under Backpressure

## Before you write
- Who writes each result slot without concurrent ~append~?
- If the first ~workers~ block, will thousands of other goroutines start?

## A fixed crew
A worker pool starts a fixed number of goroutines reading a shared queue. Backpressure means submitting a job waits for ~capacity~ or a worker, rather than starting a goroutine per input.
~~~go
jobs := make(chan int)
for i := 0; i < 3; i++ {
 go func() { for index := range jobs { fmt.Println(index) } }()
}
~~~
An unbuffered queue prevents a backlog here. ~Workers~ call the callback directly; do not start another goroutine per callback. The feeder closes the queue; join ~workers~ with ~sync.WaitGroup~.

## Generic functions and independent results
Chapter 14's generics apply to functions and result types together. ~any~ accepts every type without converting values to ~interface{}~.
~~~go
type Reply[R any] struct { Value R; Err error }
func identity[T any](x T) T { return x }
fmt.Println(identity("bay"))
~~~
Allocate a result slice matching input length and queue indexes. Each worker writes a distinct slot. ~Wait~ before reading makes completed writes visible despite out-of-order completion.
~~~go
replies := make([]Reply[string], 4)
replies[2] = Reply[string]{Value: "ready"}
~~~
An item error is not a batch error: store it and take the next job. A single shipment must not fail independent shipments.

## Cancellation and reconciliation
Pass the same ~context.Context~ to the callback. Callbacks must honor cancellation and terminate; goroutines cannot safely be killed. Check cancellation before taking work and before invoking the callback. Every channel wait needs cancellation.
~~~go
select {
case jobs <- index:
case <-ctx.Done():
 // Stop feeding, close jobs, and join workers
}
~~~
After joining, unstarted slots receive the type's zero value and the context error. A started callback retains exactly its returned value and error, even if cancellation follows.

## Generic crew values
Use the same result type with different values, retaining an error next to each value.
~~~go
reply := CrewReply[string]{Value: "oxygen", Err: nil}
fmt.Println(reply.Value, reply.Err)
~~~

## Common mistakes
~~~go
for _, cargo := range cargos {
 go inspect(cargo) // Unbounded goroutines, no backpressure
}
// Concurrent append from workers races
~~~

## Your mission
Implement ~RunCrew[T,R any](ctx context.Context, inputs []T, workers int, work func(context.Context,T)(R,error)) ([]CrewReply[R],error)~. ~CrewReply[R]~ contains ~Value R~ and ~Err error~. Context is non-~nil~ and callbacks terminate. Reject ~workers <= 0~ or ~work == nil~ with ~nil~,error before callbacks, even for empty ~inputs~.

After validation, allocate one result per input in source order. Pre-cancellation fills every result with its context error without callbacks and returns the batch context error. Otherwise empty input succeeds. Use a bounded index queue and at most ~workers~ ~workers~; more ~workers~ than ~inputs~ is valid. Each input starts once on success. Store each callback's value and error unchanged; an item failure does not stop others. Do not mutate ~inputs~. On cancellation, stop starting work when observed and join active callbacks. Unstarted slots get zero values and the context error. Return the batch context error for cancellation, otherwise ~nil~ even with item errors. No leaks, real time, or ~Report~.

Third star: preserve the same contract with structs and strings, and with independent concurrent calls without shared global state.`,
			Hints: []string{"Separate a crew-wide error from an item error.", "Give each input its own result slot.", "Review ~CrewReply[R]~ before tracing values of ~R~."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch243Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport}},
	})
}

const ch243Starter = `package main

import (
	"context"
	"errors"
)

type CrewReply[R any] struct {
	Value R
	Err   error
}

func RunCrew[T, R any](ctx context.Context, inputs []T, workers int, work func(context.Context, T) (R, error)) ([]CrewReply[R], error) {
	return nil, errors.New("unfinished")
}
func main() {}
`
const ch243Hidden = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
	"slices"
	"sync/atomic"
	"testing"
)

// Cases: empty/one/many inputs; invalid, one and surplus workers; item faults and canceled queues.
func TestHiddenCrewOrderedResultsAndItemFaults(t *testing.T) {
	labkit.LeakCheck(t)
	rng := supplyRand()
	fault := errors.New("damaged crate")
	for _, workers := range []int{1, 4, 60} {
		inputs := make([]int, 31)
		for i := range inputs {
			inputs[i] = rng.Intn(1000)
		}
		before := slices.Clone(inputs)
		got, err := RunCrew(context.Background(), inputs, workers, func(ctx context.Context, v int) (int, error) {
			if v%3 == 0 {
				return -v, fault
			}
			return v + 7, nil
		})
		if err != nil || len(got) != len(inputs) || !slices.Equal(before, inputs) {
			t.Fatal("batch", err)
		}
		for i, v := range inputs {
			want := v + 7
			var e error
			if v%3 == 0 {
				want = -v
				e = fault
			}
			if got[i].Value != want || !errors.Is(got[i].Err, e) {
				t.Fatal("source order/cause", i, got[i])
			}
		}
	}
	fn := func(context.Context, int) (int, error) { t.Error("unexpected callback"); return 0, nil }
	for _, n := range []int{0, -1} {
		if r, e := RunCrew(context.Background(), []int{}, n, fn); r != nil || e == nil {
			t.Fatal("invalid workers")
		}
	}
	if r, e := RunCrew[int, int](context.Background(), nil, 1, nil); r != nil || e == nil {
		t.Fatal("nil callback")
	}
	if r, e := RunCrew(context.Background(), []int{}, 1, fn); len(r) != 0 || e != nil {
		t.Fatal("empty", e)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, e := RunCrew(ctx, []int{1, 2}, 2, fn)
	if !errors.Is(e, context.Canceled) || len(r) != 2 {
		t.Fatal("pre-cancel", e)
	}
	for _, v := range r {
		if v.Value != 0 || !errors.Is(v.Err, context.Canceled) {
			t.Fatal("missing cancellation slot")
		}
	}
}
func TestHiddenCrewBoundedParallelismAndCancellation(t *testing.T) {
	labkit.LeakCheck(t)
	for _, stop := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		entered := make(chan int, 40)
		release := make(chan struct{})
		var active, peak atomic.Int32
		type reply struct {
			r []CrewReply[int]
			e error
		}
		done := make(chan reply, 1)
		inputs := make([]int, 30)
		for i := range inputs {
			inputs[i] = i
		}
		go func() {
			r, e := RunCrew(ctx, inputs, 3, func(c context.Context, v int) (int, error) {
				a := active.Add(1)
				for old := peak.Load(); a > old && !peak.CompareAndSwap(old, a); old = peak.Load() {
				}
				defer active.Add(-1)
				entered <- v
				select {
				case <-release:
					return v + 1, nil
				case <-c.Done():
					return 0, c.Err()
				}
			})
			done <- reply{r, e}
		}()
		seen := map[int]bool{}
		for i := 0; i < 3; i++ {
			seen[supplyReceive(t, entered)] = true
		}
		if len(seen) != 3 {
			t.Fatal("duplicate work")
		}
		if stop {
			cancel()
		} else {
			close(release)
		}
		got := supplyReceive(t, done)
		cancel()
		if peak.Load() != 3 || active.Load() != 0 {
			t.Fatal("crew concurrency", peak.Load(), active.Load())
		}
		if stop {
			if !errors.Is(got.e, context.Canceled) {
				t.Fatal(got.e)
			}
			for _, v := range got.r {
				if !errors.Is(v.Err, context.Canceled) {
					t.Fatal("lost canceled slot")
				}
			}
			if len(entered) != 0 {
				t.Fatal("started queued callback after cancellation")
			}
		} else {
			if got.e != nil || len(got.r) != 30 {
				t.Fatal("completion", got.e)
			}
			for i, v := range got.r {
				if v.Value != i+1 || v.Err != nil {
					t.Fatal("order", i, v)
				}
			}
		}
	}
}
func TestHiddenBonusCrewGenericIndependentCalls(t *testing.T) {
	labkit.LeakCheck(t)
	type cargo struct {
		Name string
		Mass int
	}
	done := make(chan error, 5)
	for i := 0; i < 5; i++ {
		go func() {
			got, err := RunCrew(context.Background(), []cargo{{"ice", 2}, {"ore", 3}}, 2, func(_ context.Context, c cargo) (string, error) { return c.Name, nil })
			if err == nil && (len(got) != 2 || got[0].Value != "ice" || got[1].Value != "ore") {
				err = errors.New("generic slots")
			}
			done <- err
		}()
	}
	for i := 0; i < 5; i++ {
		if e := supplyReceive(t, done); e != nil {
			t.Fatal(e)
		}
	}
}
`

func init() {
	add(&Level{ID: "24-4-pipeline", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24PipelineArabic, Title: "خط التكرير", Goal: "قسّم خام الرصيف على مسارات تكرير وادمج نتايجهم مع إلغاء نضيف",
		Requires: []string{"chan", "select"}, Bridge: []int{44}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch244StarterAr}},
		Lesson: `# من مسارات كتير لمجرى واحد

## قبل ما تكتب
- مين يقفل خرج كل مسار، ومين يقفل الخرج النهائي؟
- لو العميل وقف قراءة، الإلغاء هيفك كل مرحلة إزاي؟

## توزيع الشغل: fan-out
أكتر من عامل يقرأوا نفس قناة الدخل. كل رسالة يستلمها عامل واحد بس. ده توزيع للشغل، مش بث نسخة لكل عامل.
~~~go
for i := 0; i < 2; i++ {
 go func() { for sample := range feed { fmt.Println(sample) } }()
}
~~~
في المهمة كل عامل عنده قناة خرج مستقلة. كل عامل يقفل خرجه بنفسه لما ينتهي؛ هو المرسل الوحيد عليه.

## تجميع النتايج: fan-in
اعمل goroutine نقل لكل خرج عامل. كل واحدة تقرأ مسار وتبعت في الخرج النهائي. قافل واحد يستنى كل الناقلين بـ ~sync.WaitGroup~ وبعدها يقفل النهائي.
~~~go
go func() {
 defer forwards.Done()
 for n := range lane { merged <- n }
}()
go func() { forwards.Wait(); close(merged) }()
~~~
الترتيب النهائي هنا ترتيب الجاهزية، مش ترتيب الدخل. احتفظ بـ ~ID~ الرسالة عشان تربط النتيجة بالشحنة من غير الاعتماد على الموضع.

## الإلغاء لكل انتظار
المثال السابق ناجح بس من غير إلغاء. في المهمة، استقبال الدخل، إرسال خرج العامل، استقبال الناقل، وإرساله للخرج النهائي كلهم يحتاجوا ~select~ مع ~ctx.Done()~.
~~~go
select {
case <-ctx.Done():
 return
case sample, ok := <-feed:
 if !ok { return }
 fmt.Println(sample)
}
~~~
لو المستهلك وقف قراءة من غير إلغاء، الضغط العكسي هيوقف كل المراحل بشكل طبيعي. مالك السياق لازم يلغي لما يتخلى عن القراءة. ممنوع قفل قناة الدخل من عندك؛ اللي بيغذيها هو مالكها.

## هوية الشحنة
ابعت ~ID~ مع العينة عشان ترتيب الوصول ميغيّرش معناها.
~~~go
batch := OreBatch{ID: 17, Grade: 8}
reply := OreReply{ID: batch.ID, Grade: 9}
fmt.Println(reply.ID, reply.Grade)
~~~

## أخطاء شائعة
~~~go
for n := range lane { merged <- n } // انتظار إرسال مش قابل للإلغاء
close(feed) // المستهلك مش مالك دخل محقون
~~~

## المهمة
نفّذ ~RefineOre(ctx context.Context, input <-chan OreBatch, lanes int, refine func(context.Context,OreBatch)(int,error)) (<-chan OreReply,error)~. ~OreBatch~ فيها ~ID, Grade int~، و ~OreReply~ فيها ~ID, Grade int~ و ~Err error~. السياق غير ~nil~. ارفض قناة ~nil~ أو callback ~nil~ أو ~lanes <= 0~ بشكل متزامن بـ ~nil~ وخطأ، من غير تشغيل شغل.

غير كده ارجع قناة خرج فوراً وnil. وزّع الدخل على ~lanes~ عمال بحد أقصى، كل عامل له خرج خاص، وادمجهم بقنوات fan-in في قناة نهائية unbuffered. في التشغيل غير الملغي: كل batch لها رد واحد بنفس ~ID~ وقيمة وخطأ callback كما هما، بما فيهم أخطاء العناصر؛ الترتيب غير مضمون. الـ callback بتحترم السياق وبتنتهي. الدخل الفاضي المقفول يقفل الخرج، والخطأ في عنصر ميوقفش الباقي.

الإلغاء يقفل الخرج بعد انتهاء كل العمال والناقلين، حتى مع دخل مفتوح ساكت أو مستهلك مبطل قراءة. بعد ملاحظة الإلغاء متبدأش callback جديدة؛ مسموح ردود جزئية سبقت الإلغاء. الإلغاء المسبق يرجع قناة بتتقفل من غير callbacks. متقفلش الدخل. كل خرج داخلي يتقفل مرة واحدة بمالكه، والنهائي بعد الناقلين. الاختبار يصرف الخرج لحد القفل. ممنوع وقت حقيقي أو ~Report~.`,
		Hints: []string{"اتبع إشارة انتهاء الإدخال عبر كل مرحلة.", "اكتب مالك كل قناة قبل مراجعة الإغلاق.", "راجع سلوك ~range~ بعد ~close~ في المثال الصغير."},
		En: &LevelText{LessonRevision: 1, Steps: ch24PipelineEnglish, Title: "Refining Pipeline", Goal: "Fan dock ore out across refining lanes and merge their results with clean cancellation", Files: []LabFile{{Name: "main.go", Starter: ch244Starter}},
			Lesson: `# Many Lanes, One Stream

## Before you write
- Who closes each lane output, and who closes the merged output?
- If the consumer stops reading, how will cancellation release every ~stage~?

## Fan-out distributes work
Several ~workers~ read the same input channel. Each message reaches one worker. This distributes work rather than broadcasting a copy to everyone.
~~~go
for i := 0; i < 2; i++ {
 go func() { for sample := range feed { fmt.Println(sample) } }()
}
~~~
Each mission worker has its own output channel. It closes that output when finished because it is the only sender.

## Fan-in merges results
Start a forwarding goroutine per worker output. Each reads its lane and sends to the merged channel. One closer waits for all forwarders with ~sync.WaitGroup~, then closes the merged output.
~~~go
go func() {
 defer forwards.Done()
 for n := range lane { merged <- n }
}()
go func() { forwards.Wait(); close(merged) }()
~~~
Final order follows readiness rather than source position. Preserve a message ~ID~ to associate replies with shipments without assuming positional order.

## Cancel every wait
The preceding example handles only successful completion. In the mission, input receives, worker output sends, forwarder receives, and merged sends all need ~select~ with ~ctx.Done()~.
~~~go
select {
case <-ctx.Done():
 return
case sample, ok := <-feed:
 if !ok { return }
 fmt.Println(sample)
}
~~~
A consumer that stops reading without canceling naturally backpressures the whole pipeline. The context owner must cancel when abandoning reads. Never close injected input: its feeder owns it.

## Shipment identity
Carry an ~ID~ with a sample so arrival order cannot change its meaning.
~~~go
batch := OreBatch{ID: 17, Grade: 8}
reply := OreReply{ID: batch.ID, Grade: 9}
fmt.Println(reply.ID, reply.Grade)
~~~

## Common mistakes
~~~go
for n := range lane { merged <- n } // Send wait cannot be canceled
close(feed) // Consumer does not own injected input
~~~

## Your mission
Implement ~RefineOre(ctx context.Context, input <-chan OreBatch, lanes int, refine func(context.Context,OreBatch)(int,error)) (<-chan OreReply,error)~. ~OreBatch~ has ~ID, Grade int~; ~OreReply~ has ~ID, Grade int~ and ~Err error~. Context is non-~nil~. Synchronously reject ~nil~ input, ~nil~ callback, or ~lanes <= 0~ with ~nil~,error and no work.

Otherwise return an output channel immediately with ~nil~ error. Fan input out to at most ~lanes~ ~workers~, each with a separate output, and fan those outputs into one unbuffered merged channel. Without cancellation, each batch produces exactly one reply preserving its ~ID~ and the callback's grade and error, including item errors. Order is unspecified. Callbacks honor context and terminate. Closed empty input closes output; item failure does not stop others.

Cancellation closes output after all ~workers~ and forwarders finish, even with idle open input or an abandoned consumer. Start no new callback after observing cancellation; partial replies committed before cancellation are allowed. Pre-cancellation returns a channel that closes without callbacks. Never close input. Close each internal output exactly once by its owner, and the merged output after forwarders finish. Tests drain to closure. No real time or ~Report~.`,
			Hints: []string{"Follow the input-finished signal through each stage.", "Write down each channel owner before reviewing closure.", "Review how ~range~ behaves after ~close~ in the small example."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch244Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport}},
	})
}

const ch244Starter = `package main

import (
	"context"
	"errors"
)

type OreBatch struct{ ID, Grade int }
type OreReply struct {
	ID, Grade int
	Err       error
}

func RefineOre(ctx context.Context, input <-chan OreBatch, lanes int, refine func(context.Context, OreBatch) (int, error)) (<-chan OreReply, error) {
	return nil, errors.New("unfinished")
}
func main() {}
`
const ch244Hidden = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
	"testing"
)

func supplyDrain(t *testing.T, out <-chan OreReply) []OreReply {
	t.Helper()
	done := make(chan []OreReply, 1)
	go func() {
		var got []OreReply
		for v := range out {
			got = append(got, v)
		}
		done <- got
	}()
	return supplyReceive(t, done)
}

// Cases: closed-empty input, all lane counts, per-item failure, idle input and abandoned output.
func TestHiddenRefineryFanOutAndFanIn(t *testing.T) {
	labkit.LeakCheck(t)
	rng := supplyRand()
	fault := errors.New("low grade")
	for _, lanes := range []int{1, 4, 40} {
		input := make(chan OreBatch, 30)
		want := map[int]OreReply{}
		for i := 0; i < 30; i++ {
			g := rng.Intn(90) - 20
			input <- OreBatch{i, g}
			var e error
			if g < 0 {
				e = fault
			}
			want[i] = OreReply{i, g + 1, e}
		}
		close(input)
		out, err := RefineOre(context.Background(), input, lanes, func(_ context.Context, b OreBatch) (int, error) {
			var e error
			if b.Grade < 0 {
				e = fault
			}
			return b.Grade + 1, e
		})
		if err != nil || out == nil {
			t.Fatal(err)
		}
		got := supplyDrain(t, out)
		if len(got) != len(want) {
			t.Fatal("missing replies", len(got))
		}
		for _, v := range got {
			w, ok := want[v.ID]
			if !ok || v.Grade != w.Grade || !errors.Is(v.Err, w.Err) {
				t.Fatal("duplicate/wrong reply", v)
			}
			delete(want, v.ID)
		}
	}
	empty := make(chan OreBatch)
	close(empty)
	out, e := RefineOre(context.Background(), empty, 2, func(context.Context, OreBatch) (int, error) { t.Error("empty callback"); return 0, nil })
	if e != nil || len(supplyDrain(t, out)) != 0 {
		t.Fatal("empty", e)
	}
}
func TestHiddenRefineryCancellationJoinsEveryStage(t *testing.T) {
	labkit.LeakCheck(t)
	for _, mode := range []string{"idle", "blocked-output", "active", "pre-canceled"} {
		ctx, cancel := context.WithCancel(context.Background())
		input := make(chan OreBatch, 20)
		entered := make(chan struct{}, 20)
		exited := make(chan struct{}, 20)
		if mode != "idle" {
			for i := 0; i < 20; i++ {
				input <- OreBatch{i, i}
			}
		}
		if mode == "pre-canceled" {
			cancel()
		}
		out, e := RefineOre(ctx, input, 3, func(c context.Context, b OreBatch) (int, error) {
			entered <- struct{}{}
			defer func() { exited <- struct{}{} }()
			if mode == "active" {
				<-c.Done()
				return 0, c.Err()
			}
			return b.Grade, nil
		})
		if e != nil {
			t.Fatal(e)
		}
		if mode == "active" {
			for i := 0; i < 3; i++ {
				supplyReceive(t, entered)
			}
		}
		if mode == "blocked-output" {
			supplyReceive(t, entered)
		}
		cancel()
		supplyDrain(t, out)
		if mode == "pre-canceled" && len(entered) != 0 {
			t.Fatal("pre-canceled callback")
		}
		if mode == "active" && len(exited) != 3 {
			t.Fatal("output closed before worker exit")
		}
		// The caller still owns this input after the pipeline exits.
		close(input)
	}
}
func TestHiddenRefineryActuallyUsesSeveralLanes(t *testing.T) {
	labkit.LeakCheck(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	input := make(chan OreBatch, 3)
	for i := 0; i < 3; i++ {
		input <- OreBatch{i, i}
	}
	close(input)
	entered := make(chan int, 3)
	release := make([]chan struct{}, 3)
	for i := range release {
		release[i] = make(chan struct{})
	}
	out, e := RefineOre(ctx, input, 3, func(_ context.Context, b OreBatch) (int, error) {
		entered <- b.ID
		select {
		case <-release[b.ID]:
		case <-ctx.Done():
		}
		return b.Grade, nil
	})
	if e != nil {
		t.Fatal(e)
	}
	seen := map[int]bool{}
	for i := 0; i < 3; i++ {
		seen[supplyReceive(t, entered)] = true
	}
	if len(seen) != 3 {
		t.Fatal("not distributed")
	}
	for i := 2; i >= 0; i-- {
		close(release[i])
		r := supplyReceive(t, out)
		if r.ID != i {
			t.Fatal("fan-in waits for wrong lane", r)
		}
	}
	if len(supplyDrain(t, out)) != 0 {
		t.Fatal("extra reply")
	}
}
func TestHiddenRefineryPreflight(t *testing.T) {
	labkit.LeakCheck(t)
	fn := func(context.Context, OreBatch) (int, error) { t.Error("invalid callback"); return 0, nil }
	in := make(chan OreBatch)
	defer close(in)
	for _, n := range []int{0, -1} {
		if out, e := RefineOre(context.Background(), in, n, fn); out != nil || e == nil {
			t.Fatal("invalid lanes")
		}
	}
	if out, e := RefineOre(context.Background(), nil, 1, fn); out != nil || e == nil {
		t.Fatal("nil input")
	}
	if out, e := RefineOre(context.Background(), in, 1, nil); out != nil || e == nil {
		t.Fatal("nil callback")
	}
}
`

func init() {
	add(&Level{ID: "24-5-beacon-check", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24BeaconArabic, Title: "فحص المنارات", Goal: "افحص منارات HTTP بعميل محقون وحد فعلي للتوازي",
		Requires: []string{"sel:http.Client"}, Bridge: []int{33}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch245StarterAr}},
		Lesson: `# فحص منارات القاعدة

## قبل ما تكتب
- إزاي تثبت الحد الأقصى للطلبات من غير الاعتماد على السرعة؟
- مين يقفل جسم الرد حتى لو حالة HTTP فاشلة؟

## أول عميل HTTP
~http.Client~ بينفّذ طلبات HTTP. بناخده محقون عشان التست يختار النقل والـ timeout بدل شبكة حقيقية. ~Timeout~ حد زمني للطلب كله؛ متستبدلش العميل المحقون بـ ~http.DefaultClient~.
~~~go
client := &http.Client{Timeout: 2 * time.Second}
req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1/status", nil)
if err != nil { return err }
resp, err := client.Do(req)
if err != nil { return err }
defer resp.Body.Close()
fmt.Println(resp.StatusCode)
~~~
~http.NewRequestWithContext~ يربط الطلب بإلغاء الأب، و ~http.MethodGet~ اسم الطريقة. ~Do~ ممكن تفشل في النقل؛ حالة زي 503 مش transport error ولازم تفحص ~StatusCode~ بنفسك. الـ ~Body~ مورد لازم يتقفل حتى لو مش هتقراه.

~req.URL.Scheme~ و ~req.URL.Host~ بيميزوا عنوان مطلق عن مسار نسبي. اسمح بس بـ HTTP وHTTPS وhost غير فاضي قبل ~Do~.
~~~go
if req.URL.Host == "" || (req.URL.Scheme != "http" && req.URL.Scheme != "https") {
 return fmt.Errorf("invalid beacon address")
}
~~~

## Semaphore على قناة
قناة بسعة الحد تمثل التصاريح. امتلك تصريح قبل إنشاء goroutine، وارجعه بـ ~defer~ بعد الطلب. لو بدأت goroutine قبل امتلاك التصريح، عدد goroutines نفسه ممكن يبقى غير محدود.
~~~go
slots := make(chan struct{}, 2)
select {
case slots <- struct{}{}:
case <-ctx.Done():
 return ctx.Err()
}
go func() { defer func(){ <-slots }(); /* request */ }()
~~~
كل طلب يكتب خانة مختلفة في slice، و ~WaitGroup~ تستناهم. الإلغاء لازم يوقف انتظار التصريح وكمان الطلب الجاري. متعدلش إعدادات العميل المشترك.

## سيرفر اختبار محلي
~httptest.NewServer~ يفتح على loopback وبورت عشوائي، و ~server.URL~ عنوانه. ~http.HandlerFunc~ بتحوّل دالة لـ handler. ~http.ResponseWriter~ بيكتب الرد و ~*http.Request~ الطلب الداخل. اقفل السيرفر بعد التست.
~~~go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 w.WriteHeader(http.StatusNoContent)
}))
defer server.Close()
fmt.Println(server.URL)
~~~
في فحص timeout من غير انتظار، ~http.RoundTripper~ واجهة نقل فيها ~RoundTrip(*http.Request)(*http.Response,error)~. حط implementation في ~Client.Transport~ ترجع ~context.DeadlineExceeded~؛ وفي فحص التوازي حط barrier داخل handler. الوقت الحقيقي في التست مجرد watchdog للفشل، مش دليل نجاح.

~~~go
type transportFunc func(*http.Request) (*http.Response, error)
func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
client := &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
 return nil, context.DeadlineExceeded
})}
~~~

تستات اللعبة بتستخدم نفس ~httptest.Server~ فوق اتصالات في الذاكرة؛ المثال المحلي فوق ينفع خارج sandbox. مش مطلوب تبني ناقل أو سيرفر في الحل.

## حالة منارة واحدة
خلي الحالة والخطأ منفصلين عشان تميّز رد HTTP عن فشل النقل.
~~~go
reply := BeaconReply{Status: http.StatusNoContent}
fmt.Println(reply.Status, reply.Err)
~~~

## أخطاء شائعة
~~~go
resp, _ := http.Get(address) // عميل عام، تجاهل الخطأ والسياق
if resp.StatusCode == 200 { return nil } // الجسم متقفلش
~~~

## المهمة
نفّذ ~CheckBeacons(ctx context.Context, client *http.Client, addresses []string, limit int) ([]BeaconReply,error)~. ~BeaconReply~ فيها ~Status int~ و ~Err error~. السياق غير ~nil~. ارفض عميل ~nil~ أو ~limit <= 0~ بـ ~nil~ وخطأ قبل أي طلب، حتى لو القائمة فاضية. باقي الحالات ترجع slice بطول الدخل وبنفس الترتيب وnil كخطأ batch.

لكل عنوان اعمل GET بالعميل المحقون وبنفس السياق، مع حد ~limit~ للطلبات المتزامنة. العناوين المكررة مستقلة. العنوان غير المطلق أو غير HTTP/HTTPS له خطأ عنصر ومش بيتبعت. حالة 200 لحد 299 نجاح؛ غير كده ~Status~ الفعلي وخطأ عنصر. خطأ نقل أو timeout يحافظ على السبب، وStatus تبقى 0. اقفل أي جسم رد مستلم. متغيرش العناوين أو العميل. عند الإلغاء، الطلبات اللي متبدأتش تاخد 0 وخطأ السياق، والجارية تنتهي عن طريقه، والنتائج المكتملة تفضل محفوظة. استنى الكل قبل الرجوع. احترام سياسة redirects والـ timeout مسؤولية العميل المحقون؛ ممنوع شبكة عامة أو وقت ثابت في الحل.`,
		Hints: []string{"فرّق بين خطأ النقل ورد HTTP.", "راجع ترتيب النتائج مقارنة بترتيب العناوين.", "ارجع لـ~NewRequestWithContext~ قبل تتبع طلب المنارة."},
		En: &LevelText{LessonRevision: 1, Steps: ch24BeaconEnglish, Title: "Beacon Check", Goal: "Inspect HTTP beacons with an injected client and an enforced concurrency limit", Files: []LabFile{{Name: "main.go", Starter: ch245Starter}},
			Lesson: `# Checking Colony Beacons

## Before you write
- How will you prove the request ~limit~ without relying on speed?
- Who closes response bodies even for failing HTTP statuses?

## The first HTTP ~client~
~http.Client~ performs HTTP requests. Inject it so tests choose transport and timeout without public networking. ~Timeout~ bounds the whole request; never replace the injected ~client~ with ~http.DefaultClient~.
~~~go
client := &http.Client{Timeout: 2 * time.Second}
req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://127.0.0.1/status", nil)
if err != nil { return err }
resp, err := client.Do(req)
if err != nil { return err }
defer resp.Body.Close()
fmt.Println(resp.StatusCode)
~~~
~http.NewRequestWithContext~ ties requests to parent cancellation; ~http.MethodGet~ names the method. ~Do~ can fail in transport; a status such as 503 is not a transport error, so inspect ~StatusCode~. ~Body~ is a resource requiring closure even if unread.

~req.URL.Scheme~ and ~req.URL.Host~ distinguish absolute addresses from relative paths. Permit only HTTP/HTTPS with a nonempty host before ~Do~.
~~~go
if req.URL.Host == "" || (req.URL.Scheme != "http" && req.URL.Scheme != "https") {
 return fmt.Errorf("invalid beacon address")
}
~~~

## A channel semaphore
A ~capacity~-limited channel represents permits. Acquire before launching a goroutine, and release with ~defer~ after the request. Launching before acquisition can still create unbounded goroutines.
~~~go
slots := make(chan struct{}, 2)
select {
case slots <- struct{}{}:
case <-ctx.Done():
 return ctx.Err()
}
go func() { defer func(){ <-slots }(); /* request */ }()
~~~
Each request owns a distinct result slot; join using ~WaitGroup~. Cancellation must interrupt both permit acquisition and in-flight requests. Never mutate a shared ~client~'s settings.

## A local test server
~httptest.NewServer~ binds loopback on a random port; ~server.URL~ is its address. ~http.HandlerFunc~ adapts a function to a handler, with ~http.ResponseWriter~ for the response and ~*http.Request~ for the incoming request. Close the server after testing.
~~~go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 w.WriteHeader(http.StatusNoContent)
}))
defer server.Close()
fmt.Println(server.URL)
~~~
For timeout tests without waiting, ~http.RoundTripper~ is a transport interface with ~RoundTrip(*http.Request)(*http.Response,error)~. Supply a ~Client.Transport~ implementation returning ~context.DeadlineExceeded~. Use handler barriers to prove concurrency. Real test time serves only as a failure watchdog.

~~~go
type transportFunc func(*http.Request) (*http.Response, error)
func (f transportFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
client := &http.Client{Transport: transportFunc(func(r *http.Request) (*http.Response, error) {
 return nil, context.DeadlineExceeded
})}
~~~

Game tests run the same ~httptest.Server~ over in-memory connections; the loopback example works outside a restricted sandbox. Your solution need not implement a server or transport.

## One beacon result
Keep status and error separate to distinguish an HTTP response from transport failure.
~~~go
reply := BeaconReply{Status: http.StatusNoContent}
fmt.Println(reply.Status, reply.Err)
~~~

## Common mistakes
~~~go
resp, _ := http.Get(address) // Global client, ignored error and context
if resp.StatusCode == 200 { return nil } // Body leaked
~~~

## Your mission
Implement ~CheckBeacons(ctx context.Context, client *http.Client, addresses []string, limit int) ([]BeaconReply,error)~. ~BeaconReply~ contains ~Status int~ and ~Err error~. Context is non-~nil~. Reject ~nil~ ~client~ or ~limit <= 0~ with ~nil~,error before requests, including empty input. All other cases return one source-ordered result per address and ~nil~ batch error.

GET every address using the injected ~client~ and same context, with at most ~limit~ simultaneous requests. Duplicate addresses are independent. Nonabsolute or non-HTTP/HTTPS addresses get item errors without requests. ~Status~ 200 through 299 succeeds; other statuses retain the actual ~Status~ with an item error. Transport errors and timeouts preserve causes and use ~Status~ 0. Close every received response body. Do not mutate addresses or ~client~. Cancellation assigns 0 and the context error to unstarted requests, interrupts active requests, and preserves completed results. Join all work before returning. The injected ~client~ owns redirect and timeout policy. No public network or fixed real-time waits in the solution.`,
			Hints: []string{"Distinguish a transport error from an HTTP response.", "Compare result order with address order.", "Review ~NewRequestWithContext~ before tracing a beacon request."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch245Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport, "supply_http_test.go": ch24HTTPHidden}},
	})
}

const ch245Starter = `package main

import (
	"context"
	"errors"
	"net/http"
)

type BeaconReply struct {
	Status int
	Err    error
}

func CheckBeacons(ctx context.Context, client *http.Client, addresses []string, limit int) ([]BeaconReply, error) {
	return nil, errors.New("unfinished")
}
func main() {}
`
const ch245Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"orbital/colony/labkit"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
)

type beaconTransport func(*http.Request) (*http.Response, error)

func (f beaconTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type beaconBody struct {
	io.Reader
	closed *atomic.Int32
}

func (b beaconBody) Close() error { b.closed.Add(1); return nil }

// Cases: all 2xx boundaries, 3xx/4xx/5xx, duplicate addresses, malformed URLs, timeout and cancellation.
func TestHiddenBeaconsLocalStatusesAndOrder(t *testing.T) {
	labkit.LeakCheck(t)
	server := supplyServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Error("method", r.Method)
		}
		var code int
		fmt.Sscanf(r.URL.Path, "/%d", &code)
		w.WriteHeader(code)
	}))
	defer server.Close()
	client := server.Client()
	defer client.CloseIdleConnections()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	rng := supplyRand()
	codes := []int{200, 204, 299, 300, 404, 503}
	var addresses []string
	var want []int
	for i := 0; i < 30; i++ {
		v := codes[rng.Intn(len(codes))]
		addresses = append(addresses, fmt.Sprintf("%s/%d", server.URL, v))
		want = append(want, v)
	}
	before := slices.Clone(addresses)
	got, e := CheckBeacons(context.Background(), client, addresses, 4)
	if e != nil || len(got) != len(want) || !slices.Equal(before, addresses) {
		t.Fatal("batch", e)
	}
	for i, v := range got {
		if v.Status != want[i] || (v.Err == nil) != (want[i] >= 200 && want[i] < 300) {
			t.Fatal("status/order", i, v)
		}
	}
}
func TestHiddenBeaconsCloseBodiesAndPreserveTimeout(t *testing.T) {
	labkit.LeakCheck(t)
	var closed, calls atomic.Int32
	client := &http.Client{Transport: beaconTransport(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		if r.URL.Path == "/timeout" {
			return nil, context.DeadlineExceeded
		}
		return &http.Response{StatusCode: 503, Header: make(http.Header), Body: beaconBody{strings.NewReader("offline"), &closed}, Request: r}, nil
	})}
	addresses := []string{"http://beacon.local/error", "http://beacon.local/timeout", "://broken", "relative", "ftp://beacon.local/x"}
	got, e := CheckBeacons(context.Background(), client, addresses, 2)
	if e != nil || len(got) != 5 {
		t.Fatal(e)
	}
	if got[0].Status != 503 || got[0].Err == nil || got[1].Status != 0 || !errors.Is(got[1].Err, context.DeadlineExceeded) {
		t.Fatal("failure categories", got)
	}
	for _, v := range got[2:] {
		if v.Status != 0 || v.Err == nil {
			t.Fatal("invalid URL sent", v)
		}
	}
	if calls.Load() != 2 || closed.Load() != 1 {
		t.Fatal("requests/body lifetime", calls.Load(), closed.Load())
	}
}
func TestHiddenBeaconsRealConcurrencyAndCancel(t *testing.T) {
	labkit.LeakCheck(t)
	for _, limit := range []int{1, 3} {
		ctx, cancel := context.WithCancel(context.Background())
		entered := make(chan struct{}, 20)
		release := make(chan struct{})
		var active, peak atomic.Int32
		server := supplyServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			n := active.Add(1)
			defer active.Add(-1)
			for old := peak.Load(); n > old && !peak.CompareAndSwap(old, n); old = peak.Load() {
			}
			entered <- struct{}{}
			select {
			case <-release:
				w.WriteHeader(204)
			case <-r.Context().Done():
			}
		}))
		client := server.Client()
		done := make(chan []BeaconReply, 1)
		addresses := make([]string, 12)
		for i := range addresses {
			addresses[i] = server.URL
		}
		go func() {
			r, e := CheckBeacons(ctx, client, addresses, limit)
			if e != nil {
				t.Error(e)
			}
			done <- r
		}()
		for i := 0; i < limit; i++ {
			supplyReceive(t, entered)
		}
		// All permits are held; queued requests cannot begin before cancellation.
		cancel()
		got := supplyReceive(t, done)
		close(release)
		client.CloseIdleConnections()
		server.Close()
		if peak.Load() != int32(limit) || active.Load() != 0 || len(got) != 12 {
			t.Fatal("limit/lifetime", peak.Load(), active.Load())
		}
		for _, r := range got {
			if r.Status != 0 || !errors.Is(r.Err, context.Canceled) {
				t.Fatal("canceled result", r)
			}
		}
	}
}
func TestHiddenBeaconsPreflight(t *testing.T) {
	labkit.LeakCheck(t)
	var calls atomic.Int32
	client := &http.Client{Transport: beaconTransport(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return nil, errors.New("unexpected network")
	})}
	for _, n := range []int{0, -1} {
		if r, e := CheckBeacons(context.Background(), client, nil, n); r != nil || e == nil {
			t.Fatal("limit")
		}
	}
	if r, e := CheckBeacons(context.Background(), nil, nil, 1); r != nil || e == nil {
		t.Fatal("nil client")
	}
	if r, e := CheckBeacons(context.Background(), client, nil, 1); len(r) != 0 || e != nil {
		t.Fatal("empty")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, e := CheckBeacons(ctx, client, []string{"http://beacon.local"}, 1)
	if e != nil || len(r) != 1 || !errors.Is(r[0].Err, context.Canceled) || calls.Load() != 0 {
		t.Fatal("pre-cancel", e)
	}
}
`

func init() {
	add(&Level{ID: "24-6-chunk-download", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24ChunkArabic, Title: "تحميل بالأجزاء", Goal: "نزّل خريطة ملاحة بأجزاء متوازية وانشرها بس بعد اكتمالها",
		Requires: []string{"sel:context.Context"}, Bridge: []int{35}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch246StarterAr}},
		Lesson: `# خريطة الملاحة على أجزاء

## قبل ما تكتب
- لو السيرفر تجاهل Range وبعت الخريطة كلها، هتكتبها فين؟
- لو جزء فشل، إزاي تمنع نشر نص خريطة وتستنى باقي العمال؟

## طلب نطاق بايتات
~Range~ header بتطلب نطاق شامل البداية والنهاية. آخر جزء ممكن يكون أقصر. النجاح الجزئي حالة ~http.StatusPartialContent~ أو 206، ومعاه ~Content-Range~ يثبت المكان والحجم الكلي.
~~~go
req.Header.Set("Range", "bytes=4-7")
// الرد المتوقع للمثال: Content-Range: bytes 4-7/12
~~~
افحص الـ header والحالة قبل استخدام البايتات. 200 معناها الجسم كله؛ المهمة تسمح بده في أول طلب بس، فتقرأه كخريطة كاملة من offset 0 ومن غير باقي الطلبات. 200 في جزء لاحق خطأ. الحجم جاي من خطة محقونة، مش تخمين من الشبكة.

## كتابة في مكان محدد
~io.WriterAt~ فيها ~WriteAt(p []byte, off int64)(n int,err error)~. الـ offset مش cursor مشترك، فالأجزاء غير المتداخلة ينفع تتكتب بالتوازي لو الـ store بيدعم ده. ~io.Writer~ العادية مش بديل.
~~~go
n, err := stage.WriteAt([]byte("dock"), 4)
if err == nil && n != 4 { err = io.ErrShortWrite }
~~~
متستعملش ~Seek~ مشترك بين العمال. افحص الكتابة القصيرة حتى لو الخطأ ~nil~، وحافظ على سبب أخطاء القراءة والكتابة.

## قراءة محدودة
~io.LimitReader~ يمنع جسم طويل من استهلاك ذاكرة مفتوحة. اقرأ المتوقع زائد بايت واحد عشان تكشف زيادة البيانات، وبعدها لازم الطول يساوي المتوقع بالظبط. ~io.ReadAll~ ممكن يرجع بايتات وخطأ؛ الخطأ يرفض الجزء.
~~~go
payload, err := io.ReadAll(io.LimitReader(resp.Body, expected+1))
if err != nil { return err }
if int64(len(payload)) != expected { return io.ErrUnexpectedEOF }
~~~
اقفل ~resp.Body~ في النجاح والفشل. دي نفس ملكية موارد HTTP في المهمة السابقة.

## مرحلة ثم نشر أو إلغاء
الـ interface ~ChartStage~ تجمع ~io.WriterAt~ مع ~Commit() error~ و ~Abort() error~. ~Commit~ تنشر الخريطة المكتملة. ~Abort~ تنظف الجزء المؤقت؛ المتجر المحقون مسؤول عن معنى النشر والتنظيف، ومش مطلوب ملفات من اللاعب.
~~~go
type Stage interface {
 io.WriterAt
 Commit() error
 Abort() error
}
~~~
بعد التحقق، أي فشل يستدعي ~Abort~ مرة بعد انتظار العمال، حتى لو ~Commit~ نفسها فشلت. لو التنظيف فشل، ~errors.Join~ من الفصل 5 يحافظ على السببين.
~~~go
err = errors.Join(err, stage.Abort())
~~~
أول طلب يكشف دعم النطاق. بعده وزع الأجزاء المتبقية على طاقم محدود، وألغي سياق داخلي عند أول فشل، واقفل طابوره واستنى كل العمال قبل التنظيف.

## خطة من القاعدة
الخطة بتحدد الحجم والسعة؛ المتجر هو اللي يملك مرحلة النشر.
~~~go
plan := ChartPlan{Address: "http://127.0.0.1/chart", Size: 12, Chunk: 4, Workers: 2}
fmt.Println(plan.Size, plan.Chunk)
~~~

## أخطاء شائعة
~~~go
stage.WriteAt(payload, 0) // كل الأجزاء بتدوس على أول الخريطة
stage.Commit() // غلط قبل انتظار العمال وفحص أخطائهم
~~~

## المهمة
نفّذ ~DownloadChart(ctx context.Context, client *http.Client, plan ChartPlan, stage ChartStage) error~. ~ChartPlan~ فيها ~Address string~، و ~Size, Chunk int64~، و ~Workers int~. ~ChartStage~ تضم ~io.WriterAt~ و ~Commit() error~ و ~Abort() error~. السياق غير ~nil~، والمتجر يدعم كتابات متوازية غير متداخلة. ارفض ~client~ أو ~stage~ ~nil~، عنوان غير مطلق HTTP/HTTPS، حجم خارج 0 لحد 1048576 شامل، ~Chunk <= 0~، أو ~Workers <= 0~ قبل أي طلب أو نداء للمتجر. المدخلات وclient متتغيرش.

بعد التحقق تبدأ معاملة: صفر حجم يعمل ~Commit~ من غير طلب. الإلغاء المسبق يعمل ~Abort~ ويرجع سبب السياق. غير كده اعمل GET لأول نطاق من 0 لحد ~min(Chunk,Size)-1~. لو 206، افحص ~Content-Range~ بالنص المطابق ~bytes start-end/Size~ وطول الجسم بالظبط، واكتب عند البداية. نفّذ الباقي بطاقم بحد ~Workers~، من غير تداخل نطاقات، وبنفس فحص الحالة والـ header والطول. لو أول رد 200، اقرأ ~Size~ بالظبط واكتب مرة عند 0 ومتبعتش باقي الأجزاء. أي حالة تانية أو 200 لاحق أو نقص/زيادة/خطأ قراءة أو كتابة يوقف النشر.

كل الطلبات مرتبطة بالسياق المشتق من الأب وتستخدم العميل المحقون فقط. اقفل كل الأجسام. أي فشل يلغي العمال ويستناهم ويعمل ~Abort~ مرة؛ احتفظ بالأسباب، بما فيها خطأ ~Abort~. النجاح يستنى كل العمال ثم ~Commit~ مرة، من غير ~Abort~. فشل ~Commit~ يتبعه ~Abort~. مفيش ملفات أو وقت حقيقي أو طلبات عامة أو ~Report~ في الحل؛ تستات السيرفر محلية في الذاكرة.`,
		Hints: []string{"افصل مكان الجزء عن ترتيب اكتمال طلبه.", "راجع الفرق بين كتابة البيانات ونشر المرحلة.", "ارجع لتوقيع ~WriteAt~ في ~ChartStage~ عند تتبع offset."},
		En: &LevelText{LessonRevision: 1, Steps: ch24ChunkEnglish, Title: "Chunk Download", Goal: "Download a navigation chart in parallel chunks and publish only after completion", Files: []LabFile{{Name: "main.go", Starter: ch246Starter}},
			Lesson: `# A Navigation Chart in Pieces

## Before you write
- If the server ignores Range and sends the whole chart, where will you write it?
- If one part fails, how will you prevent publication and join other ~workers~?

## Byte ranges
The ~Range~ header asks for inclusive start and end offsets. The last part may be shorter. Partial success is ~http.StatusPartialContent~, 206, with ~Content-Range~ identifying offsets and total size.
~~~go
req.Header.Set("Range", "bytes=4-7")
// Expected example response: Content-Range: bytes 4-7/12
~~~
Check headers and status before using bytes. A 200 response means the whole body. This mission permits that only on the first request: read it as a full chart at offset 0, without further requests. A later 200 is an error. ~Size~ comes from an injected plan rather than network guesses.

## Positional writes
~io.WriterAt~ provides ~WriteAt(p []byte, off int64)(n int,err error)~. An offset is not a shared cursor, so disjoint parts can be written concurrently when the store supports that. An ordinary ~io.Writer~ is insufficient.
~~~go
n, err := stage.WriteAt([]byte("dock"), 4)
if err == nil && n != 4 { err = io.ErrShortWrite }
~~~
Never share a ~Seek~ cursor among ~workers~. Detect short writes even without an error, preserving read and write failure causes.

## Bounded reads
~io.LimitReader~ prevents oversized bodies from consuming unbounded memory. Read the expected length plus one byte to detect excess, then require exactly the expected length. ~io.ReadAll~ may return bytes and an error; an error rejects the part.
~~~go
payload, err := io.ReadAll(io.LimitReader(resp.Body, expected+1))
if err != nil { return err }
if int64(len(payload)) != expected { return io.ErrUnexpectedEOF }
~~~
Close ~resp.Body~ on success and failure, following the preceding mission's HTTP resource ownership.

## Stage, then publish or abort
The ~ChartStage~ interface combines ~io.WriterAt~, ~Commit() error~, and ~Abort() error~. ~Commit~ publishes the complete chart. ~Abort~ cleans temporary state. The injected store owns publication and cleanup semantics; player filesystem access is unnecessary.
~~~go
type Stage interface {
 io.WriterAt
 Commit() error
 Abort() error
}
~~~
After validation, every failure calls ~Abort~ once after joining ~workers~, including ~Commit~ failure. If cleanup fails too, chapter 5's ~errors.Join~ preserves both causes.
~~~go
err = errors.Join(err, stage.Abort())
~~~
The first request probes range support. Feed remaining parts to a bounded crew; cancel an internal context on failure, close the queue, and join all ~workers~ before cleanup.

## A plan from base
The plan defines size and worker ~capacity~; the store owns publication.
~~~go
plan := ChartPlan{Address: "http://127.0.0.1/chart", Size: 12, Chunk: 4, Workers: 2}
fmt.Println(plan.Size, plan.Chunk)
~~~

## Common mistakes
~~~go
stage.WriteAt(payload, 0) // Every part overwrites the beginning
stage.Commit() // Wrong before joining workers and checking their errors
~~~

## Your mission
Implement ~DownloadChart(ctx context.Context, client *http.Client, plan ChartPlan, stage ChartStage) error~. ~ChartPlan~ has ~Address string~, ~Size, Chunk int64~, and ~Workers int~. ~ChartStage~ embeds ~io.WriterAt~ plus ~Commit() error~ and ~Abort() error~. Context is non-~nil~. The ~stage~ supports concurrent disjoint writes. Reject ~nil~ ~client~ or ~stage~, nonabsolute or non-HTTP/HTTPS address, size outside 0 through 1048576 inclusive, ~Chunk <= 0~, or ~Workers <= 0~ before requests or store calls. Do not mutate ~inputs~ or ~client~.

After validation, begin a transaction: zero size commits without requests. Pre-cancellation aborts and preserves the context cause. Otherwise GET the first range, 0 through ~min(Chunk,Size)-1~. For 206, require the exact ~Content-Range~ text ~bytes start-end/Size~ and exact body length, then write at the range start. Process remaining disjoint parts with at most ~Workers~ concurrent ~workers~ using the same status, header, and length checks. If the first response is 200, read exactly ~Size~ bytes and write once at 0 without further requests. Any other status, later 200, missing/excess bytes, read error, or write error prevents publication.

All requests use a child of the caller context and only the injected ~client~. Close all bodies. Failure cancels and joins ~workers~, then calls ~Abort~ once, preserving causes including ~Abort~ errors. Success joins ~workers~ then calls ~Commit~ once without ~Abort~. ~Commit~ failure triggers ~Abort~. No filesystem, real time, public requests, or ~Report~ in the solution; test HTTP servers use in-memory connections.`,
			Hints: []string{"Separate a chunk position from its request completion order.", "Review the difference between writing data and publishing a stage.", "Return to the ~WriteAt~ signature in ~ChartStage~ when tracing offset."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch246Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport, "supply_http_test.go": ch24HTTPHidden}},
	})
}

const ch246Starter = `package main

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type ChartPlan struct {
	Address     string
	Size, Chunk int64
	Workers     int
}
type ChartStage interface {
	io.WriterAt
	Commit() error
	Abort() error
}

func DownloadChart(ctx context.Context, client *http.Client, plan ChartPlan, stage ChartStage) error {
	return errors.New("unfinished")
}
func main() {}
`
const ch246Hidden = `package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"orbital/colony/labkit"
	"sync"
	"sync/atomic"
	"testing"
	"testing/iotest"
)

type chartMemory struct {
	mu                            sync.Mutex
	bytes                         []byte
	writes, commits, aborts       int
	writeErr, commitErr, abortErr error
	short                         bool
}

func (s *chartMemory) WriteAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.writes++
	if s.writeErr != nil && off > 0 {
		return 0, s.writeErr
	}
	if s.short {
		return len(p) - 1, nil
	}
	if off < 0 || off+int64(len(p)) > int64(len(s.bytes)) {
		return 0, errors.New("write outside chart")
	}
	return copy(s.bytes[off:], p), nil
}
func (s *chartMemory) Commit() error { s.commits++; return s.commitErr }
func (s *chartMemory) Abort() error  { s.aborts++; clear(s.bytes); return s.abortErr }
func chartHandler(t *testing.T, payload []byte, mode string, calls *atomic.Int32) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		if r.Method != http.MethodGet {
			t.Error("method")
		}
		var start, end int64
		if _, err := fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &start, &end); err != nil || start < 0 || end < start || end >= int64(len(payload)) {
			w.WriteHeader(416)
			return
		}
		if mode == "whole" {
			w.Write(payload)
			return
		}
		if mode == "later200" && start > 0 {
			w.Write(payload)
			return
		}
		if mode == "failure" && start > 0 {
			w.WriteHeader(503)
			return
		}
		header := fmt.Sprintf("bytes %d-%d/%d", start, end, len(payload))
		if mode == "header" {
			header = "bytes 0-0/0"
		}
		w.Header().Set("Content-Range", header)
		w.WriteHeader(206)
		part := payload[start : end+1]
		if mode == "short" {
			part = part[:len(part)-1]
		}
		w.Write(part)
		if mode == "long" {
			w.Write([]byte{99})
		}
	})
}

// Boundaries: zero/one/exact/partial chunks, unsupported ranges, malformed replies, store failures.
func TestHiddenChartSeededBytesAndFallback(t *testing.T) {
	labkit.LeakCheck(t)
	rng := supplyRand()
	for _, size := range []int{0, 1, 16, 17, 129} {
		payload := make([]byte, size)
		rng.Read(payload)
		for _, mode := range []string{"range", "whole"} {
			var calls atomic.Int32
			server := supplyServer(t, chartHandler(t, payload, mode, &calls))
			stage := &chartMemory{bytes: make([]byte, size)}
			plan := ChartPlan{server.URL, int64(size), 16, 3}
			err := DownloadChart(context.Background(), server.Client(), plan, stage)
			server.Close()
			if err != nil || !bytes.Equal(stage.bytes, payload) || stage.commits != 1 || stage.aborts != 0 {
				t.Fatal("download", size, mode, err)
			}
			want := (size + 15) / 16
			if mode == "whole" && size > 0 {
				want = 1
			}
			if calls.Load() != int32(want) {
				t.Fatal("range/fallback requests", calls.Load(), want)
			}
		}
	}
}
func TestHiddenChartFailuresAbortWithoutPublishing(t *testing.T) {
	labkit.LeakCheck(t)
	payload := make([]byte, 65)
	supplyRand().Read(payload)
	for _, mode := range []string{"header", "short", "long", "later200", "failure"} {
		var calls atomic.Int32
		server := supplyServer(t, chartHandler(t, payload, mode, &calls))
		stage := &chartMemory{bytes: make([]byte, len(payload))}
		err := DownloadChart(context.Background(), server.Client(), ChartPlan{server.URL, 65, 8, 3}, stage)
		server.Close()
		if err == nil || stage.commits != 0 || stage.aborts != 1 || !bytes.Equal(stage.bytes, make([]byte, 65)) {
			t.Fatal("unclean failure", mode, err)
		}
	}
	fault := errors.New("stage storage failure")
	cleanup := errors.New("cleanup failure")
	for _, mode := range []string{"write", "short-write", "commit", "abort"} {
		var calls atomic.Int32
		server := supplyServer(t, chartHandler(t, payload, "range", &calls))
		stage := &chartMemory{bytes: make([]byte, 65)}
		switch mode {
		case "write":
			stage.writeErr = fault
		case "short-write":
			stage.short = true
		case "commit":
			stage.commitErr = fault
		case "abort":
			stage.writeErr = fault
			stage.abortErr = cleanup
		}
		err := DownloadChart(context.Background(), server.Client(), ChartPlan{server.URL, 65, 8, 3}, stage)
		server.Close()
		want := fault
		if mode == "short-write" {
			want = io.ErrShortWrite
		}
		if !errors.Is(err, want) || stage.aborts != 1 {
			t.Fatal("store cause", mode, err)
		}
		if mode == "abort" && !errors.Is(err, cleanup) {
			t.Fatal("cleanup cause lost", err)
		}
	}
}
func TestHiddenChartParallelPartsCancelAndJoin(t *testing.T) {
	labkit.LeakCheck(t)
	for _, stop := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		entered := make(chan struct{}, 10)
		release := make(chan struct{})
		var active, peak atomic.Int32
		payload := make([]byte, 40)
		supplyRand().Read(payload)
		server := supplyServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var a, b int64
			fmt.Sscanf(r.Header.Get("Range"), "bytes=%d-%d", &a, &b)
			if a > 0 {
				n := active.Add(1)
				defer active.Add(-1)
				for old := peak.Load(); n > old && !peak.CompareAndSwap(old, n); old = peak.Load() {
				}
				entered <- struct{}{}
				select {
				case <-release:
				case <-r.Context().Done():
					return
				}
			}
			w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/40", a, b))
			w.WriteHeader(206)
			w.Write(payload[a : b+1])
		}))
		stage := &chartMemory{bytes: make([]byte, 40)}
		done := make(chan error, 1)
		go func() { done <- DownloadChart(ctx, server.Client(), ChartPlan{server.URL, 40, 8, 2}, stage) }()
		supplyReceive(t, entered)
		supplyReceive(t, entered)
		if stop {
			cancel()
		} else {
			close(release)
		}
		err := supplyReceive(t, done)
		cancel()
		if stop {
			close(release)
		}
		server.Close()
		if peak.Load() != 2 || active.Load() != 0 {
			t.Fatal("parts not bounded/parallel", peak.Load())
		}
		if stop {
			if !errors.Is(err, context.Canceled) || stage.aborts != 1 || stage.commits != 0 {
				t.Fatal("cancel cleanup", err)
			}
		} else if err != nil || stage.commits != 1 || !bytes.Equal(stage.bytes, payload) {
			t.Fatal("parallel bytes", err)
		}
	}
}

type chartTransport func(*http.Request) (*http.Response, error)

func (f chartTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type chartBody struct {
	io.Reader
	closed *atomic.Int32
}

func (b chartBody) Close() error { b.closed.Add(1); return nil }
func TestHiddenChartReadErrorClosesBodyAndPreflightDoesNothing(t *testing.T) {
	labkit.LeakCheck(t)
	fault := errors.New("radio read failed")
	var calls, closed atomic.Int32
	client := &http.Client{Transport: chartTransport(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		return &http.Response{StatusCode: 206, Header: http.Header{"Content-Range": []string{"bytes 0-7/8"}}, Body: chartBody{iotest.ErrReader(fault), &closed}, Request: r}, nil
	})}
	stage := &chartMemory{bytes: make([]byte, 8)}
	plan := ChartPlan{"http://chart.local", 8, 8, 2}
	if err := DownloadChart(context.Background(), client, plan, stage); !errors.Is(err, fault) || stage.aborts != 1 || closed.Load() != 1 {
		t.Fatal("read cause/body", err)
	}
	for _, bad := range []ChartPlan{{"relative", 8, 1, 1}, {"ftp://x", 8, 1, 1}, {"http://x", -1, 1, 1}, {"http://x", 1048577, 1, 1}, {"http://x", 8, 0, 1}, {"http://x", 8, 1, 0}} {
		s := &chartMemory{}
		before := calls.Load()
		if err := DownloadChart(context.Background(), client, bad, s); err == nil || s.aborts+s.commits+s.writes != 0 || calls.Load() != before {
			t.Fatal("preflight side effects", bad, err)
		}
	}
	if DownloadChart(context.Background(), nil, plan, stage) == nil || DownloadChart(context.Background(), client, plan, nil) == nil {
		t.Fatal("nil dependency")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	s := &chartMemory{}
	if err := DownloadChart(ctx, client, plan, s); !errors.Is(err, context.Canceled) || s.aborts != 1 || calls.Load() != 1 {
		t.Fatal("pre-cancel", err)
	}
}
`

func init() {
	add(&Level{ID: "24-7-map-tiles", Chapter: 24, Kind: "lab", LessonRevision: 1, Steps: ch24TilesArabic, Title: "بلاطات الخريطة", Goal: "صغّر بلاطات الملاحة بطاقم محدود واحفظ صيغة كل صورة",
		Requires: []string{"go"}, Bridge: []int{39}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch247StarterAr}},
		Lesson: `# صور خريطة الملاحة

## قبل ما تكتب
- هتفحص ميزانية البكسلات قبل فك صورة ضخمة إزاي؟
- لو بلاطة بايظة، باقي البلاطات هتكمل ولا هتضيع؟

## الصورة مش امتداد الملف
~image.DecodeConfig~ تقرأ الأبعاد والصيغة من البايتات قبل فك الصورة كلها. ~image.Decode~ بترجع ~image.Image~ والصيغة. استيراد ~image/png~ و ~image/jpeg~ بيسجّل الديكودر المناسب تلقائياً؛ هنستخدم نفس الحزمتين في الترميز.
~~~go
cfg, format, err := image.DecodeConfig(bytes.NewReader(encoded))
if err != nil { return err }
fmt.Println(cfg.Width, cfg.Height, format)
img, format, err := image.Decode(bytes.NewReader(encoded))
~~~
اسم البلاطة ~ID~ للمستعمرة، مش دليل صيغة. اقبل ~png~ و ~jpeg~ من المحتوى بس. قبل ضرب العرض في الارتفاع، افحص ~width > maxPixels/height~ بعد التأكد إن الاتنين موجبين. ده يمنع overflow ويفرض ميزانية الذاكرة قبل فك الصورة.

## مستطيلات وألوان
~image.Image~ بتوفر ~Bounds()~ و ~At(x,y)~. ~image.Rect~ تبني مستطيل نهايته غير شاملة. ~image.NewNRGBA~ تخزن ألوان مع alpha من غير فقد شفافية.
~~~go
canvas := image.NewNRGBA(image.Rect(0, 0, 3, 2))
bounds := img.Bounds()
color := img.At(bounds.Min.X, bounds.Min.Y)
canvas.Set(0, 0, color)
~~~
هنستخدم nearest-neighbor: كل بكسل خرج ياخد عينة من بكسل دخل. ~image.Pt~ تبني نقطة. ~draw.Draw~ تنسخ مستطيل من صورة، و ~draw.Src~ تستبدل اللون بما فيه alpha بدل دمجه فوق خلفية. استورد ~image/draw~ باسم ~draw~.
~~~go
source := image.Pt(bounds.Min.X+sx, bounds.Min.Y+sy)
draw.Draw(canvas, image.Rect(x, y, x+1, y+1), img, source, draw.Src)
~~~
لخرج عرضه ~ow~ وارتفاعه ~oh~، خُد ~sx=x*sw/ow~ و ~sy=y*sh/oh~. النهايات الحصرية تمنع قراءة خارج الدخل. حافظ على النسبة بملاءمة الصورة جوه الصندوق، من غير تكبير، وقرّب البعد التابع لتحت مع حد أدنى 1. استخدم حساب ~int64~ في ضرب الأبعاد.

## حفظ الصيغة
~png.Encode~ و ~jpeg.Encode~ بياخدوا ~io.Writer~ وصورة. ~bytes.Buffer~ يخزّن الناتج في الذاكرة. JPEG ضاغط بفقد ومفيهوش alpha؛ PNG يحفظ العينات والشفافية.
~~~go
var encodedTile bytes.Buffer
err := png.Encode(&encodedTile, canvas)
// في مسار JPEG:
err = jpeg.Encode(&encodedTile, canvas, &jpeg.Options{Quality: 90})
~~~
~jpeg.Options.Quality~ من 1 لحد 100؛ المهمة بتثبّت 90. بعد الترميز، مرّر البايتات لـ callback حفظ محقونة؛ مفيش ملفات أو أسماء مسارات محتاجة للنظام الحقيقي.

## طاقم وإلغاء
استخدم طاقم indexes من 24-3: عامل لكل موضع نتيجة، وعدد ثابت من العمال. افحص السياق قبل فك الصورة، وبين صفوف التصغير، وقبل الحفظ. نداء الحفظ يحترم السياق وينتهي، ومتزامن لحد عدد العمال. خطأ بلاطة يفضل محلي ومش بيلغي الباقي.

## بيانات بلاطة
الاسم هو هوية البلاطة، والصيغة جاية من البايتات.
~~~go
tile := MapTile{Name: "crater", Data: encoded}
spec := TileSpec{Width: 3, Height: 2, Workers: 2, MaxPixels: 4096}
reply := TileReply{Name: tile.Name, Format: "png", Width: spec.Width}
fmt.Println(reply.Name, reply.Width)
~~~

## أخطاء شائعة
~~~go
pixels := cfg.Width * cfg.Height // ممكن overflow قبل الفحص
png.Encode(out, img) // غلط لو محتوى الدخل JPEG ولازم يحافظ على الصيغة
// worker لكل صورة من غير حد هيستهلك ذاكرة كتير
~~~

## المهمة
نفّذ ~BuildMapTiles(ctx context.Context, tiles []MapTile, spec TileSpec, save func(context.Context,string,[]byte)error) ([]TileReply,error)~. ~MapTile~ فيها ~Name string~ و ~Data []byte~. ~TileSpec~ فيها ~Width, Height, Workers, MaxPixels int~. ~TileReply~ فيها ~Name, Format string~ و ~Width, Height int~ و ~Err error~. السياق غير ~nil~، والحفظ ينتهي ويحترم السياق.

ارفض قبل فك أو حفظ أي صورة: callback ~nil~، عدد عمال غير موجب، عرض أو ارتفاع خارج 1 لحد 2048، ميزانية بكسلات خارج 1 لحد 1048576، اسم فاضي أو مكرر. ارجع ~nil~ وخطأ batch. غير كده ارجع نتيجة لكل اسم بنفس الترتيب وnil كخطأ batch، حتى في الإلغاء وأخطاء الصور. الفاضي نجاح. متغيرش المدخلات.

كل بلاطة تفحص ~DecodeConfig~ وتقبل PNG أو JPEG بس، وترفض أبعاد غير موجبة أو أكبر من الميزانية قبل ~Decode~ الكامل. اختار ~scale=min(1,Width/sw,Height/sh)~ بنسب كسرية، وأبعاد الخرج هي ~max(1,floor(sw*scale))~ و ~max(1,floor(sh*scale))~. مفيش تكبير. صغّر بأخذ العينة ~floor(x*sw/ow), floor(y*sh/oh)~؛ احفظ PNG كـ PNG وJPEG كـ JPEG بجودة 90. نادِ ~save~ مرة واحدة بالاسم الأصلي وبايتات الخرج. النجاح يرجع الاسم والصيغة ~png~ أو ~jpeg~ والأبعاد وخطأ ~nil~.

ملف بايظ أو زيادة بكسلات أو خطأ حفظ يرجع الاسم وخطأ مع سبب محفوظ، وباقي حقول النتيجة صفر. كمل باقي الصور. عند الإلغاء أوقف الجديد عند ملاحظته واستنى الجاري؛ غير المكتمل يرجع الاسم وخطأ السياق، والمكتمل قبل الإلغاء يفضل ناجح. العمال والحفظ مع بعض ميتعدوش ~Workers~؛ مفيش goroutine لكل عنصر، أو تسريب، أو ملفات، أو وقت حقيقي، أو ~Report~.

النجمة التالتة: عيّنات PNG الملونة الشفافة تحتفظ بقيم اللون وalpha بالظبط بعد التصغير، باستخدام استبدال ~draw.Src~.`,
		Hints: []string{"افصل قراءة الإعدادات عن فك pixels.", "اربط كل بلاطة بخانة نتيجة مستقلة.", "راجع ~DecodeConfig~ و~Format~ في المثال قبل اختيار encoder."},
		En: &LevelText{LessonRevision: 1, Steps: ch24TilesEnglish, Title: "Map Tiles", Goal: "Resize navigation tiles with a bounded crew while preserving each image format", Files: []LabFile{{Name: "main.go", Starter: ch247Starter}},
			Lesson: `# Navigation Map Images

## Before you write
- How will you check pixel budget before fully decoding a huge image?
- Will a corrupt tile prevent independent tiles from completing?

## Content determines format
~image.DecodeConfig~ reads dimensions and format before full decoding. ~image.Decode~ returns an ~image.Image~ and format. Importing ~image/png~ and ~image/jpeg~ registers their decoders automatically; the same packages encode output.
~~~go
cfg, format, err := image.DecodeConfig(bytes.NewReader(encoded))
if err != nil { return err }
fmt.Println(cfg.Width, cfg.Height, format)
img, format, err := image.Decode(bytes.NewReader(encoded))
~~~
A tile name is a colony ~ID~, not evidence of format. Accept only content identified as ~png~ or ~jpeg~. Before multiplying dimensions, check ~width > maxPixels/height~ after verifying both are positive. This prevents overflow and enforces the memory budget before full decoding.

## Rectangles and colors
~image.Image~ provides ~Bounds()~ and ~At(x,y)~. ~image.Rect~ creates a rectangle with exclusive upper bounds. ~image.NewNRGBA~ stores colors and alpha without losing transparency.
~~~go
canvas := image.NewNRGBA(image.Rect(0, 0, 3, 2))
bounds := img.Bounds()
color := img.At(bounds.Min.X, bounds.Min.Y)
canvas.Set(0, 0, color)
~~~
Use nearest-neighbor sampling: each output pixel takes one source pixel. ~image.Pt~ creates a point. ~draw.Draw~ copies a rectangle from an image; ~draw.Src~ replaces color including alpha rather than compositing over a background. Import ~image/draw~ as ~draw~.
~~~go
source := image.Pt(bounds.Min.X+sx, bounds.Min.Y+sy)
draw.Draw(canvas, image.Rect(x, y, x+1, y+1), img, source, draw.Src)
~~~
For output dimensions ~ow,oh~, use ~sx=x*sw/ow~ and ~sy=y*sh/oh~. Exclusive bounds avoid out-of-source reads. Fit the image inside the target box without upscaling, round the dependent dimension down, and clamp to at least 1. Use ~int64~ for dimension products.

## Preserve encoding
~png.Encode~ and ~jpeg.Encode~ accept an ~io.Writer~ and image. A ~bytes.Buffer~ stores output in memory. JPEG is lossy without alpha; PNG preserves sampled colors and transparency.
~~~go
var encodedTile bytes.Buffer
err := png.Encode(&encodedTile, canvas)
// In the JPEG branch:
err = jpeg.Encode(&encodedTile, canvas, &jpeg.Options{Quality: 90})
~~~
~jpeg.Options.Quality~ ranges from 1 through 100; this mission fixes it at 90. Pass encoded bytes to an injected saving callback, without filesystem paths or real player files.

## Crew and cancellation
Reuse 24-3's index crew: distinct result slots and a fixed worker count. Check context before decoding, between resize rows, and before saving. Save callbacks honor context and terminate; they may run concurrently up to the worker ~limit~. Tile errors stay local without canceling independent tiles.

## Tile values
The name identifies a tile; its format comes from bytes.
~~~go
tile := MapTile{Name: "crater", Data: encoded}
spec := TileSpec{Width: 3, Height: 2, Workers: 2, MaxPixels: 4096}
reply := TileReply{Name: tile.Name, Format: "png", Width: spec.Width}
fmt.Println(reply.Name, reply.Width)
~~~

## Common mistakes
~~~go
pixels := cfg.Width * cfg.Height // May overflow before validation
png.Encode(out, img) // Wrong for JPEG input when preserving format
// One worker per image without a bound consumes excessive memory
~~~

## Your mission
Implement ~BuildMapTiles(ctx context.Context, tiles []MapTile, spec TileSpec, save func(context.Context,string,[]byte)error) ([]TileReply,error)~. ~MapTile~ has ~Name string~ and ~Data []byte~. ~TileSpec~ has ~Width, Height, Workers, MaxPixels int~. ~TileReply~ has ~Name, Format string~, ~Width, Height int~, and ~Err error~. Context is non-~nil~; saving terminates and honors context.

Before any decode or ~save~, reject ~nil~ callback, nonpositive ~workers~, width/height outside 1 through 2048, pixel budget outside 1 through 1048576, or empty/duplicate names, returning ~nil~,batch error. Otherwise return one source-ordered reply per name and ~nil~ batch error, even for cancellation and image errors. Empty input succeeds. Do not mutate ~inputs~.

Each tile uses ~DecodeConfig~, accepts PNG/JPEG only, and rejects nonpositive dimensions or an exceeded pixel budget before full ~Decode~. Use rational ~scale=min(1,Width/sw,Height/sh)~; output dimensions are ~max(1,floor(sw*scale))~ and ~max(1,floor(sh*scale))~. Never upscale. Sample ~floor(x*sw/ow), floor(y*sh/oh)~. Encode PNG as PNG and JPEG as JPEG at quality 90. Call ~save~ once with the original name and output bytes. Success returns name, format ~png~ or ~jpeg~, dimensions, and ~nil~ error.

Corruption, excessive pixels, or ~save~ failure returns name and a cause-preserving error with other fields zero. Continue independent tiles. Cancellation stops new work when observed and joins active work; unfinished tiles return their name and context error, while successes committed before cancellation remain successful. ~Workers~ including saving never exceed ~Workers~. No goroutine per input, leaks, filesystem, real time, or ~Report~.

Third star: colored transparent PNG samples preserve exact color and alpha after resizing, using replacement via ~draw.Src~.`,
			Hints: []string{"Separate config reading from pixel decoding.", "Give every tile its own result slot.", "Review ~DecodeConfig~ and ~Format~ before choosing an encoder."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch247Hidden, HiddenFiles: map[string]string{"supply_helpers_test.go": ch24HiddenSupport}},
	})
}

const ch247Starter = `package main

import (
	"context"
	"errors"
)

type MapTile struct {
	Name string
	Data []byte
}
type TileSpec struct{ Width, Height, Workers, MaxPixels int }
type TileReply struct {
	Name, Format  string
	Width, Height int
	Err           error
}

func BuildMapTiles(ctx context.Context, tiles []MapTile, spec TileSpec, save func(context.Context, string, []byte) error) ([]TileReply, error) {
	return nil, errors.New("unfinished")
}
func main() {}
`
const ch247Hidden = `package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"orbital/colony/labkit"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
)

func tileEncoded(t *testing.T, w, h int, format string, pattern bool) []byte {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	rng := supplyRand()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := color.NRGBA{33, 80, 190, 255}
			if pattern {
				c = color.NRGBA{uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(rng.Intn(256)), uint8(1 + rng.Intn(254))}
			}
			if pattern && x == 0 {
				c.A = 0
			}
			if pattern && x == 1 {
				c.A = 255
			}
			img.SetNRGBA(x, y, c)
		}
	}
	var out bytes.Buffer
	var err error
	if format == "png" {
		err = png.Encode(&out, img)
	} else {
		err = jpeg.Encode(&out, img, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}

// Cases: zero/one/many tiles, no upscale, tall/wide/one-pixel rounding, corrupt and oversized inputs.
func TestHiddenTilesFormatFitAndInputOwnership(t *testing.T) {
	labkit.LeakCheck(t)
	rng := supplyRand()
	var tiles []MapTile
	sizes := [][2]int{{1, 1}, {1, 15}, {15, 1}, {20, 10}, {10, 20}, {9, 7}}
	for i := 0; i < 24; i++ {
		sizes = append(sizes, [2]int{1 + rng.Intn(32), 1 + rng.Intn(32)})
	}
	for i, d := range sizes {
		f := "png"
		if i%2 == 1 {
			f = "jpeg"
		}
		tiles = append(tiles, MapTile{fmt.Sprintf("sector-%d.wrong-extension", i), tileEncoded(t, d[0], d[1], f, false)})
	}
	before := make([]MapTile, len(tiles))
	for i, v := range tiles {
		before[i] = MapTile{v.Name, slices.Clone(v.Data)}
	}
	for _, workers := range []int{1, 4, 40} {
		outputs := map[string][]byte{}
		var mu sync.Mutex
		got, e := BuildMapTiles(context.Background(), tiles, TileSpec{7, 5, workers, 4096}, func(_ context.Context, name string, p []byte) error {
			mu.Lock()
			defer mu.Unlock()
			if outputs[name] != nil {
				t.Error("duplicate save")
			}
			outputs[name] = slices.Clone(p)
			return nil
		})
		if e != nil || len(got) != len(tiles) || !reflect.DeepEqual(before, tiles) {
			t.Fatal("batch/input mutation", e)
		}
		for i, r := range got {
			cfg, format, err := image.DecodeConfig(bytes.NewReader(outputs[r.Name]))
			want := "png"
			if i%2 == 1 {
				want = "jpeg"
			}
			sw, sh := sizes[i][0], sizes[i][1]
			// Independent rational oracle: choose the smallest of 1, 7/sw, and 5/sh.
			num, den := 1, 1
			if 7*den < num*sw {
				num, den = 7, sw
			}
			if 5*den < num*sh {
				num, den = 5, sh
			}
			w, h := max(1, sw*num/den), max(1, sh*num/den)
			if err != nil || r.Err != nil || format != want || r.Format != want || cfg.Width != w || cfg.Height != h || r.Width != w || r.Height != h || r.Name != tiles[i].Name {
				t.Fatal("fit/format/order", i, r, cfg, err)
			}
		}
	}
}
func TestHiddenTilesIndependentFailuresAndPreflight(t *testing.T) {
	labkit.LeakCheck(t)
	fault := errors.New("map store full")
	tiles := []MapTile{{"ok", tileEncoded(t, 4, 4, "png", false)}, {"corrupt", []byte("broken")}, {"large", tileEncoded(t, 20, 20, "png", false)}, {"save-fail", tileEncoded(t, 3, 3, "jpeg", false)}}
	var calls atomic.Int32
	got, e := BuildMapTiles(context.Background(), tiles, TileSpec{2, 2, 3, 16}, func(_ context.Context, name string, p []byte) error {
		calls.Add(1)
		if name == "save-fail" {
			return fault
		}
		return nil
	})
	if e != nil || len(got) != 4 || got[0].Err != nil || calls.Load() != 2 {
		t.Fatal("independent tiles", e, calls.Load())
	}
	for i := 1; i < 4; i++ {
		if got[i].Name != tiles[i].Name || got[i].Err == nil || got[i].Format != "" || got[i].Width+got[i].Height != 0 {
			t.Fatal("failed metadata", got[i])
		}
	}
	if !errors.Is(got[3].Err, fault) {
		t.Fatal("save cause", got[3])
	}
	fn := func(context.Context, string, []byte) error { t.Error("preflight/empty/canceled save"); return nil }
	for _, s := range []TileSpec{{0, 2, 1, 10}, {2, -1, 1, 10}, {2049, 2, 1, 10}, {2, 2, 0, 10}, {2, 2, 1, 0}, {2, 2, 1, 1048577}} {
		if r, e := BuildMapTiles(context.Background(), tiles, s, fn); r != nil || e == nil {
			t.Fatal("settings", s)
		}
	}
	valid := TileSpec{2, 2, 1, 1024}
	for _, ts := range [][]MapTile{{{"", nil}}, {{"same", nil}, {"same", nil}}} {
		if r, e := BuildMapTiles(context.Background(), ts, valid, fn); r != nil || e == nil {
			t.Fatal("names")
		}
	}
	if r, e := BuildMapTiles(context.Background(), tiles, valid, nil); r != nil || e == nil {
		t.Fatal("nil save")
	}
	if r, e := BuildMapTiles(context.Background(), nil, valid, fn); len(r) != 0 || e != nil {
		t.Fatal("empty")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, e := BuildMapTiles(ctx, tiles, valid, fn)
	if e != nil || len(r) != len(tiles) {
		t.Fatal("pre-cancel batch")
	}
	for i, v := range r {
		if v.Name != tiles[i].Name || !errors.Is(v.Err, context.Canceled) {
			t.Fatal("pre-cancel slot", v)
		}
	}
}
func TestHiddenTilesBoundedSavingAndCancellation(t *testing.T) {
	labkit.LeakCheck(t)
	for _, stop := range []bool{false, true} {
		ctx, cancel := context.WithCancel(context.Background())
		payload := tileEncoded(t, 8, 8, "png", false)
		tiles := make([]MapTile, 20)
		for i := range tiles {
			tiles[i] = MapTile{fmt.Sprint(i), payload}
		}
		entered := make(chan string, 20)
		release := make(chan struct{})
		var active, peak atomic.Int32
		type reply struct {
			r []TileReply
			e error
		}
		done := make(chan reply, 1)
		go func() {
			r, e := BuildMapTiles(ctx, tiles, TileSpec{4, 4, 3, 64}, func(c context.Context, name string, p []byte) error {
				n := active.Add(1)
				defer active.Add(-1)
				for old := peak.Load(); n > old && !peak.CompareAndSwap(old, n); old = peak.Load() {
				}
				entered <- name
				select {
				case <-release:
					return nil
				case <-c.Done():
					return c.Err()
				}
			})
			done <- reply{r, e}
		}()
		for i := 0; i < 3; i++ {
			supplyReceive(t, entered)
		}
		if stop {
			cancel()
		} else {
			close(release)
		}
		got := supplyReceive(t, done)
		cancel()
		if got.e != nil || len(got.r) != 20 || active.Load() != 0 || peak.Load() != 3 {
			t.Fatal("worker lifetime/limit", got.e, peak.Load())
		}
		for _, r := range got.r {
			if stop {
				if !errors.Is(r.Err, context.Canceled) {
					t.Fatal("canceled tile", r)
				}
			} else if r.Err != nil {
				t.Fatal(r.Err)
			}
		}
		if stop && len(entered) != 0 {
			t.Fatal("queued saves started after cancellation")
		}
	}
}
func TestHiddenBonusTilesExactTransparentSamples(t *testing.T) {
	labkit.LeakCheck(t)
	source := tileEncoded(t, 9, 7, "png", true)
	original, _, e := image.Decode(bytes.NewReader(source))
	if e != nil {
		t.Fatal(e)
	}
	var encoded []byte
	got, e := BuildMapTiles(context.Background(), []MapTile{{"alpha", source}}, TileSpec{4, 3, 1, 63}, func(_ context.Context, _ string, p []byte) error { encoded = slices.Clone(p); return nil })
	if e != nil || len(got) != 1 || got[0].Err != nil {
		t.Fatal("alpha tile", e)
	}
	out, format, e := image.Decode(bytes.NewReader(encoded))
	if e != nil || format != "png" {
		t.Fatal("alpha format", e)
	}
	w, h := out.Bounds().Dx(), out.Bounds().Dy()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			want := color.NRGBAModel.Convert(original.At(x*9/w, y*7/h)).(color.NRGBA)
			actual := color.NRGBAModel.Convert(out.At(x, y)).(color.NRGBA)
			if actual != want {
				t.Fatal("sample/alpha lost", x, y, actual, want)
			}
		}
	}
}
`

const ch241StarterAr = `package main

import (
	"context"
	"errors"
)

func RelaySupplies(ctx context.Context, rounds int, visit func(drone, round int) error) error {
	return errors.New("لسه مخلصتش")
}
func main() {}
`

const ch242StarterAr = `package main

import (
	"context"
	"errors"
)

func UnloadDepots(ctx context.Context, batches [][]int, capacity int, consume func(int) error) error {
	return errors.New("لسه مخلصتش")
}
func main() {}
`

const ch243StarterAr = `package main

import (
	"context"
	"errors"
)

type CrewReply[R any] struct {
	Value R
	Err   error
}

func RunCrew[T, R any](ctx context.Context, inputs []T, workers int, work func(context.Context, T) (R, error)) ([]CrewReply[R], error) {
	return nil, errors.New("لسه مخلصتش")
}
func main() {}
`

const ch244StarterAr = `package main

import (
	"context"
	"errors"
)

type OreBatch struct{ ID, Grade int }
type OreReply struct {
	ID, Grade int
	Err       error
}

func RefineOre(ctx context.Context, input <-chan OreBatch, lanes int, refine func(context.Context, OreBatch) (int, error)) (<-chan OreReply, error) {
	return nil, errors.New("لسه مخلصتش")
}
func main() {}
`

const ch245StarterAr = `package main

import (
	"context"
	"errors"
	"net/http"
)

type BeaconReply struct {
	Status int
	Err    error
}

func CheckBeacons(ctx context.Context, client *http.Client, addresses []string, limit int) ([]BeaconReply, error) {
	return nil, errors.New("لسه مخلصتش")
}
func main() {}
`

const ch246StarterAr = `package main

import (
	"context"
	"errors"
	"io"
	"net/http"
)

type ChartPlan struct {
	Address     string
	Size, Chunk int64
	Workers     int
}
type ChartStage interface {
	io.WriterAt
	Commit() error
	Abort() error
}

func DownloadChart(ctx context.Context, client *http.Client, plan ChartPlan, stage ChartStage) error {
	return errors.New("لسه مخلصتش")
}
func main() {}
`

const ch247StarterAr = `package main

import (
	"context"
	"errors"
)

type MapTile struct {
	Name string
	Data []byte
}
type TileSpec struct{ Width, Height, Workers, MaxPixels int }
type TileReply struct {
	Name, Format  string
	Width, Height int
	Err           error
}

func BuildMapTiles(ctx context.Context, tiles []MapTile, spec TileSpec, save func(context.Context, string, []byte) error) ([]TileReply, error) {
	return nil, errors.New("لسه مخلصتش")
}
func main() {}
`
