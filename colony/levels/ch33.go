package levels

func init() {
	add(&Level{ID: "33-1-retry-client", Chapter: 33, Kind: "lab", LessonRevision: 1, Steps: ch33Lesson331Arabic, Title: "عميل عنيد", Goal: "عيد طلبات رصيف الإمداد اللي فشلت بأمان، بانتظار متزايد ومن غير ما تكرر شحنة مرتين",
		Requires: []string{"iface"}, Bridge: []int{83}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch331Starter}},
		Lesson: `# عميل عنيد لرصيف الإمداد

## قبل ما تكتب
- لو الرصيف رد 503 على طلب "اخصم 40 لتر وقود"، تعيد الطلب ولا لأ؟ وإيه اللي يضمن إن الخصم ميحصلش مرتين؟
- لو 500 عميل فشلوا في نفس اللحظة وكلهم عادوا بعد ثانية بالظبط، إيه اللي هيحصل للرصيف؟

## http.RoundTripper: طبقة تحت الـ Client
~http.Client~ مش هو اللي بيبعت البايتات؛ بيسلّم كل طلب لحاجة اسمها ~Transport~ من نوع ~http.RoundTripper~. ده interface فيه method واحدة:
~~~go
type RoundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}
~~~
أي نوع عنده ~RoundTrip~ بالشكل ده ينفع يتحط في ~Transport~. يعني تقدر تلف الـ transport الحقيقي بطبقة بتاعتك، زي الـ middleware بتاع فصل 26 بس من ناحية العميل.
~~~go
type logged struct{ next http.RoundTripper }

func (l logged) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println(r.Method, r.URL.Path)
	return l.next.RoundTrip(r)
}

client := &http.Client{Transport: logged{next: http.DefaultTransport}}
~~~
~http.DefaultTransport~ هو الـ transport العادي اللي بيفتح اتصالات TCP. قاعدة الـ RoundTripper: رجّع يا رد يا خطأ، ومتعدلش الطلب الأصلي؛ لو محتاج تغيّره خد نسخة بـ ~r.Clone(ctx)~.

## interface بتاعك للانتظار
الانتظار بين المحاولات لو اتكتب ~time.Sleep~ جوه الكود، التست هيستنى بجد. الحل interface صغير بتعرّفه أنت، والتست يحقن نسخة بتسجل المدة ومش بتنام:
~~~go
type Pauser interface {
	Pause(ctx context.Context, d time.Duration) error
}
~~~
النسخة الحقيقية تستنى ~d~ أو لحد ما ~ctx~ يتلغي، وترجع ~ctx.Err()~ ساعتها. الكود اللي بيعيد مش فارق معاه مين اللي بينفّذ الانتظار.

## مين يستاهل إعادة
مش كل فشل يتعاد. الأكواد 429 و 502 و 503 و 504 معناها "الرصيف مشغول أو الوسيط مش واصل"، وغالبًا الطلب متنفذش. خطأ من الـ transport نفسه (اتصال اتقطع) برضه يتعاد. أما 400 و 404 و 500 فبيرجعوا زي ما هم: إعادتهم مش هتغيّر النتيجة، و 500 ممكن تكون العملية اتنفذت نصها.
~~~go
func retryStatus(code int) bool { return code == 503 || code == 504 }
~~~

## idempotency: الإعادة الآمنة
طلب idempotent يعني تنفيذه مرتين زي تنفيذه مرة. ~GET~ و ~HEAD~ و ~OPTIONS~ بيقروا بس، و ~PUT~ و ~DELETE~ بيحطوا حالة نهائية، فإعادتهم آمنة. ~POST~ و ~PATCH~ ممكن يخصموا مرتين.
الحل المعروف هيدر ~Idempotency-Key~: العميل يبعت مفتاح فريد مع الطلب، والسيرفر يفتكر المفتاح ومينفذش نفس المفتاح مرتين. لو المفتاح موجود، إعادة ~POST~ بقت آمنة.
~~~go
if r.Method == "POST" && r.Header.Get("Idempotency-Key") == "" { tries = 1 }
~~~

## الجسم بيتقري مرة واحدة
~r.Body~ مجرى بيخلص بعد أول قراءة؛ لو بعته تاني هيوصل فاضي. ~http.NewRequest~ لما تديله ~strings.Reader~ أو ~bytes.Reader~ أو ~bytes.Buffer~ بيملأ ~r.GetBody~: دالة بترجع نسخة جديدة من الجسم.
~~~go
body, err := r.GetBody()
again := r.Clone(r.Context())
again.Body = body
~~~
لو فيه جسم و ~GetBody~ بـ ~nil~ مفيش طريقة ترجع البايتات، فالطلب ده ميتعادش. ~http.NoBody~ معناها جسم فاضي، وده زي ~nil~.

## الرد الفاشل لازم يتقفل
كل ~resp.Body~ ماسك اتصال. لو هتسيب رد عشان تعيد، اقرا شوية منه واقفله، وإلا الاتصالات بتتسرب والـ transport مبيعرفش يعيد استخدامها.
~~~go
io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
resp.Body.Close()
~~~
الرد الأخير اللي هترجعه للمستدعي متقفلوش؛ ده شغله هو.

## backoff مع jitter
backoff أُسّي يعني الانتظار بيتضاعف: ~Base~ ثم ~2*Base~ ثم ~4*Base~، لحد سقف ~Max~. ده بيدّي الرصيف وقت يفوق.
لو كل العملاء استنوا نفس المدة بالظبط هيرجعوا مع بعض تاني. الـ jitter بيضيف عشوائية: نستنى نص المدة الثابت وعليه رقم عشوائي لحد النص التاني.
~~~go
d := 400 * time.Millisecond
half := d / 2
wait := half + time.Duration(rand.Int63n(int64(d-half)+1)) // من 200ms لـ 400ms
~~~
~rand.Int63n(n)~ بترجع رقم من 0 لحد ~n-1~. هنا بنحقن دالة بنفس الشكل عشان التست يعرف الرقم.

## أخطاء شائعة
~~~go
resp, err := next.RoundTrip(r)
if resp.StatusCode == 503 { continue } // الرد اتساب مفتوح والاتصال اتسرب
next.RoundTrip(r) // في المحاولة التانية: r.Body اتقرا خلاص، الجسم يوصل فاضي
d := base << attempt // مفيش سقف: المحاولة العاشرة تستنى ساعات
time.Sleep(d) // التست يستنى بجد، والإلغاء مش بيوقفه
~~~

## المهمة
عرّف ~Pauser~ بالشكل اللي فوق، و ~Policy{Attempts int; Base, Max time.Duration}~. نفّذ ~NewRetrier(next http.RoundTripper, p Policy, pause Pauser, jitter func(n int64) int64) (*Retrier, error)~: يرجع ~nil,error~ لو ~next~ أو ~pause~ أو ~jitter~ بـ ~nil~، أو ~Attempts~ برة 1 لـ 10، أو ~Base <= 0~، أو ~Max < Base~.

~(r *Retrier) RoundTrip(req)~ يبعت لحد ~Attempts~ محاولة إجمالي. قبل كل محاولة لو ~req.Context()~ اتلغى يرجع ~nil~ وخطأ السياق. رد كوده مش 429 أو 502 أو 503 أو 504 يرجع فورًا. خطأ transport وسياق الطلب اتلغى يرجع فورًا. غير كده يعيد.

الإعادة مسموحة بس لـ ~GET~ و ~HEAD~ و ~OPTIONS~ و ~PUT~ و ~DELETE~، أو ~POST~ و ~PATCH~ ومعاهم ~Idempotency-Key~ مش فاضي. أي طريقة تانية، أو جسم مش فاضي و ~GetBody~ بـ ~nil~: محاولة واحدة بس. كل إعادة تبعت نسخة من الطلب بجسم جديد من ~GetBody~.

قبل الإعادة رقم n (أول إعادة n=1): ~d = Base~ متضاعف n-1 مرة ومحدود بـ ~Max~، و ~half = d/2~، والانتظار ~half + jitter(int64(d-half)+1)~. اقفل الرد الفاشل الأول، وبعدين نادي ~pause.Pause(ctx, wait)~؛ لو رجّع خطأ رجّعه مع ~nil~. بعد آخر محاولة رجّع آخر رد أو آخر خطأ زي ما هو ومن غير ما تقفله. مفيش نوم حقيقي، و ~main~ فاضية.

النجمة التالتة: لو الرد 429 أو 503 وفيه ~Retry-After~ رقم ثواني صحيح غير سالب، الانتظار يبقى الثواني دي محدودة بـ ~Max~، من غير jitter. أي قيمة تانية تمشي على الـ backoff العادي.`,
		Hints: ch33Lesson331HintsArabic,
		En: &LevelText{LessonRevision: 1, Steps: ch33Lesson331English, Title: "Stubborn Client", Goal: "Retry failed supply-dock requests safely, with growing waits, without ever shipping a cargo twice", Files: []LabFile{{Name: "main.go", Starter: ch331Starter}},
			Lesson: `# A Stubborn Client for the Supply Dock

## Before you write
- If the dock answers 503 to "deduct 40 litres of fuel", do you resend? What guarantees the deduction does not happen twice?
- If 500 clients fail at the same instant and all retry exactly one second later, what happens to the dock?

## http.RoundTripper: the layer under the Client
~http.Client~ does not send the bytes itself; it hands every request to its ~Transport~, an ~http.RoundTripper~. That is an interface with one method:
~~~go
type RoundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}
~~~
Any type with that ~RoundTrip~ fits in ~Transport~. So you can wrap the real transport with your own layer, like the chapter 26 middleware but on the client side.
~~~go
type logged struct{ next http.RoundTripper }

func (l logged) RoundTrip(r *http.Request) (*http.Response, error) {
	fmt.Println(r.Method, r.URL.Path)
	return l.next.RoundTrip(r)
}

client := &http.Client{Transport: logged{next: http.DefaultTransport}}
~~~
~http.DefaultTransport~ is the normal transport that opens TCP connections. The RoundTripper rule: return a response or an error, and never modify the original request; to change it, take a copy with ~r.Clone(ctx)~.

## Your own interface for waiting
If the wait between attempts is a ~time.Sleep~ inside the code, the test really waits. The fix is a small interface you define yourself; the test injects a version that records the duration and never sleeps:
~~~go
type Pauser interface {
	Pause(ctx context.Context, d time.Duration) error
}
~~~
The real version waits ~d~ or until ~ctx~ is cancelled, returning ~ctx.Err()~ then. The retrying code does not care who performs the wait.

## Which failures deserve a retry
Not every failure is retried. Codes 429, 502, 503, and 504 mean "the dock is busy or the intermediary cannot reach it", and the request was most likely not executed. An error from the transport itself (a dropped connection) is retried too. But 400, 404, and 500 are returned as they are: retrying will not change the answer, and a 500 may mean the operation half happened.
~~~go
func retryStatus(code int) bool { return code == 503 || code == 504 }
~~~

## Idempotency: safe retries
An idempotent request has the same effect run twice as run once. ~GET~, ~HEAD~, and ~OPTIONS~ only read, and ~PUT~ and ~DELETE~ set a final state, so retrying them is safe. ~POST~ and ~PATCH~ may deduct twice.
The known fix is an ~Idempotency-Key~ header: the client sends a unique key, and the server remembers keys and never executes the same key twice. With a key present, retrying a ~POST~ becomes safe.
~~~go
if r.Method == "POST" && r.Header.Get("Idempotency-Key") == "" { tries = 1 }
~~~

## A body is read once
~r.Body~ is a stream that is exhausted after the first read; sent again it arrives empty. When ~http.NewRequest~ gets a ~strings.Reader~, ~bytes.Reader~, or ~bytes.Buffer~, it fills ~r.GetBody~: a function returning a fresh copy of the body.
~~~go
body, err := r.GetBody()
again := r.Clone(r.Context())
again.Body = body
~~~
If there is a body and ~GetBody~ is ~nil~, the bytes cannot be recovered, so that request is never retried. ~http.NoBody~ means an empty body, the same as ~nil~.

## A failed response must be closed
Every ~resp.Body~ holds a connection. When you drop a response to retry, read a little of it and close it, or connections leak and the transport cannot reuse them.
~~~go
io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
resp.Body.Close()
~~~
Do not close the final response you return to the caller; that is the caller's job.

## Backoff with jitter
Exponential backoff doubles the wait: ~Base~, then ~2*Base~, then ~4*Base~, up to a ceiling ~Max~. It gives the dock time to recover.
If every client waits exactly the same time, they all come back together again. Jitter adds randomness: wait the fixed half of the delay plus a random amount up to the other half.
~~~go
d := 400 * time.Millisecond
half := d / 2
wait := half + time.Duration(rand.Int63n(int64(d-half)+1)) // 200ms to 400ms
~~~
~rand.Int63n(n)~ returns a number from 0 to ~n-1~. Here a function of the same shape is injected so the test knows the number.

## Common mistakes
~~~go
resp, err := next.RoundTrip(r)
if resp.StatusCode == 503 { continue } // Response left open, connection leaked
next.RoundTrip(r) // On the second attempt r.Body is already read and arrives empty
d := base << attempt // No ceiling: the tenth attempt waits for hours
time.Sleep(d) // The test really waits, and cancellation cannot stop it
~~~

## Your mission
Define ~Pauser~ as shown above, and ~Policy{Attempts int; Base, Max time.Duration}~. Implement ~NewRetrier(next http.RoundTripper, p Policy, pause Pauser, jitter func(n int64) int64) (*Retrier, error)~: return ~nil,error~ if ~next~, ~pause~, or ~jitter~ is ~nil~, ~Attempts~ is outside 1 to 10, ~Base <= 0~, or ~Max < Base~.

~(r *Retrier) RoundTrip(req)~ sends up to ~Attempts~ attempts in total. Before each attempt, if ~req.Context()~ is cancelled, return ~nil~ and the context error. A response whose code is not 429, 502, 503, or 504 returns at once. A transport error while the request context is cancelled returns at once. Anything else is retried.

Retries are allowed only for ~GET~, ~HEAD~, ~OPTIONS~, ~PUT~, and ~DELETE~, or for ~POST~ and ~PATCH~ with a nonempty ~Idempotency-Key~. Any other method, or a nonempty body with a ~nil~ ~GetBody~, gets exactly one attempt. Every retry sends a copy of the request with a fresh body from ~GetBody~.

Before retry number n (the first retry is n=1): ~d = Base~ doubled n-1 times and capped at ~Max~, ~half = d/2~, and the wait is ~half + jitter(int64(d-half)+1)~. Close the failed response first, then call ~pause.Pause(ctx, wait)~; if it returns an error, return that error with ~nil~. After the last attempt return the last response or the last error as is, without closing it. No real sleeping, and ~main~ stays empty.

Third star: if the response is 429 or 503 with a ~Retry-After~ that is a valid non-negative whole number of seconds, the wait is those seconds capped at ~Max~, with no jitter. Any other value follows the normal backoff.`,
			Hints: ch33Lesson331HintsEnglish},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch331Hidden},
	})
}

const ch331Starter = `package main

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Pauser interface {
	Pause(ctx context.Context, d time.Duration) error
}

type Policy struct {
	Attempts int
	Base     time.Duration
	Max      time.Duration
}

type Retrier struct{}

func NewRetrier(next http.RoundTripper, p Policy, pause Pauser, jitter func(n int64) int64) (*Retrier, error) {
	return nil, errors.New("TODO")
}
func (r *Retrier) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("TODO")
}
func main() {}
`

const ch331Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type retryPause struct {
	waits []time.Duration
	fail  error
}

func (p *retryPause) Pause(ctx context.Context, d time.Duration) error {
	p.waits = append(p.waits, d)
	return p.fail
}

type retryBody struct {
	io.ReadCloser
	closed *atomic.Int64
}

func (b retryBody) Close() error { b.closed.Add(1); return b.ReadCloser.Close() }

// retryWire counts transport calls, fails the first len(errs) calls, and tracks closed bodies.
type retryWire struct {
	errs   []error
	calls  atomic.Int64
	closed atomic.Int64
}

func (w *retryWire) RoundTrip(r *http.Request) (*http.Response, error) {
	n := int(w.calls.Add(1))
	if n <= len(w.errs) && w.errs[n-1] != nil {
		return nil, w.errs[n-1]
	}
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	resp.Body = retryBody{resp.Body, &w.closed}
	return resp, nil
}

// retryDock answers the scripted codes in order, then 200 with "ok"; it records every body it saw.
type retryDock struct {
	mu     sync.Mutex
	codes  []int
	header map[string]string
	hits   int
	bodies []string
}

func (d *retryDock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	d.mu.Lock()
	defer d.mu.Unlock()
	d.bodies = append(d.bodies, string(raw))
	d.hits++
	for k, v := range d.header {
		w.Header().Set(k, v)
	}
	if d.hits <= len(d.codes) {
		w.WriteHeader(d.codes[d.hits-1])
		fmt.Fprintf(w, "busy-%d", d.hits)
		return
	}
	fmt.Fprintf(w, "ok-%d", d.hits)
}

type retryJitter struct {
	rng  *rand.Rand
	asks []int64
	outs []int64
}

func (j *retryJitter) pick(n int64) int64 {
	v := j.rng.Int63n(n)
	j.asks, j.outs = append(j.asks, n), append(j.outs, v)
	return v
}

func retryWant(t *testing.T, p Policy, j *retryJitter, waits []time.Duration, retries int) {
	t.Helper()
	if len(waits) != retries || len(j.asks) != retries {
		t.Fatalf("want %d waits, got %v (jitter asked %v)", retries, waits, j.asks)
	}
	for i := 0; i < retries; i++ {
		d := p.Base
		for k := 0; k < i && d < p.Max; k++ {
			d *= 2
		}
		if d > p.Max {
			d = p.Max
		}
		half := d / 2
		if j.asks[i] != int64(d-half)+1 || waits[i] != half+time.Duration(j.outs[i]) {
			t.Fatalf("retry %d: delay %v, jitter asked %d (want %d), waited %v (want %v)", i+1, d, j.asks[i], int64(d-half)+1, waits[i], half+time.Duration(j.outs[i]))
		}
	}
}

func retrySetup(t *testing.T, codes []int) (*retryDock, *httptest.Server) {
	dock := &retryDock{codes: codes}
	srv := httptest.NewServer(dock)
	t.Cleanup(srv.Close)
	return dock, srv
}

func retryRead(t *testing.T, resp *http.Response, err error, code int, body string) {
	t.Helper()
	if err != nil || resp == nil {
		t.Fatalf("want %d %q, got %v", code, body, err)
	}
	raw, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != code || string(raw) != body {
		t.Fatalf("want %d %q, got %d %q", code, body, resp.StatusCode, raw)
	}
}

// Boundaries: every invalid field, Attempts 1 and 10, Max equal to Base.
func TestHiddenRetryRejectsBadPolicy(t *testing.T) {
	pause, jit := &retryPause{}, func(n int64) int64 { return 0 }
	ok := Policy{Attempts: 3, Base: time.Millisecond, Max: time.Second}
	for _, tc := range []struct {
		next   http.RoundTripper
		p      Policy
		pause  Pauser
		jitter func(int64) int64
	}{
		{nil, ok, pause, jit},
		{http.DefaultTransport, ok, nil, jit},
		{http.DefaultTransport, ok, pause, nil},
		{http.DefaultTransport, Policy{Attempts: 0, Base: time.Millisecond, Max: time.Second}, pause, jit},
		{http.DefaultTransport, Policy{Attempts: 11, Base: time.Millisecond, Max: time.Second}, pause, jit},
		{http.DefaultTransport, Policy{Attempts: 3, Base: 0, Max: time.Second}, pause, jit},
		{http.DefaultTransport, Policy{Attempts: 3, Base: -time.Millisecond, Max: time.Second}, pause, jit},
		{http.DefaultTransport, Policy{Attempts: 3, Base: time.Second, Max: time.Second - 1}, pause, jit},
	} {
		if r, err := NewRetrier(tc.next, tc.p, tc.pause, tc.jitter); r != nil || err == nil {
			t.Fatalf("accepted %+v", tc.p)
		}
	}
	for _, p := range []Policy{{Attempts: 1, Base: time.Second, Max: time.Second}, {Attempts: 10, Base: time.Nanosecond, Max: time.Hour}} {
		if r, err := NewRetrier(http.DefaultTransport, p, pause, jit); r == nil || err != nil {
			t.Fatalf("rejected %+v: %v", p, err)
		}
	}
}

// Success after k retryable failures: exact backoff and jitter, every failed response closed.
func TestHiddenRetryRecoversWithBackoff(t *testing.T) {
	rng := rand.New(rand.NewSource(3311))
	retryable := []int{429, 502, 503, 504}
	methods := []string{"GET", "HEAD", "OPTIONS", "PUT", "DELETE"}
	for round := 0; round < 8; round++ {
		p := Policy{Attempts: 2 + rng.Intn(7), Base: time.Duration(1+rng.Intn(300)) * time.Millisecond}
		p.Max = p.Base * time.Duration(1+rng.Intn(9))
		if round == 0 {
			p.Max = p.Base
		}
		k := rng.Intn(p.Attempts)
		codes := make([]int, k)
		for i := range codes {
			codes[i] = retryable[rng.Intn(len(retryable))]
		}
		dock, srv := retrySetup(t, codes)
		wire, pause, jit := &retryWire{}, &retryPause{}, &retryJitter{rng: rand.New(rand.NewSource(int64(round)))}
		r, err := NewRetrier(wire, p, pause, jit.pick)
		if err != nil {
			t.Fatal(err)
		}
		method := methods[rng.Intn(len(methods))]
		req, _ := http.NewRequest(method, srv.URL+"/fuel", nil)
		resp, err := r.RoundTrip(req)
		want := fmt.Sprintf("ok-%d", k+1)
		if method == "HEAD" {
			want = ""
		}
		retryRead(t, resp, err, 200, want)
		if dock.hits != k+1 || wire.calls.Load() != int64(k+1) {
			t.Fatalf("round %d: %d failures, dock saw %d hits", round, k, dock.hits)
		}
		if wire.closed.Load() != int64(k+1) {
			t.Fatalf("round %d: %d of %d failed responses closed", round, wire.closed.Load()-1, k)
		}
		retryWant(t, p, jit, pause.waits, k)
	}
}

// Exhaustion returns the last failing response unread and open; transport errors are retried and returned last.
func TestHiddenRetryGivesUpWithLastAnswer(t *testing.T) {
	rng := rand.New(rand.NewSource(3312))
	for round := 0; round < 4; round++ {
		p := Policy{Attempts: 1 + rng.Intn(5), Base: 50 * time.Millisecond, Max: time.Second}
		codes := make([]int, 20)
		for i := range codes {
			codes[i] = 503
		}
		dock, srv := retrySetup(t, codes)
		wire, pause, jit := &retryWire{}, &retryPause{}, &retryJitter{rng: rng}
		r, _ := NewRetrier(wire, p, pause, jit.pick)
		req, _ := http.NewRequest("GET", srv.URL+"/fuel", nil)
		resp, err := r.RoundTrip(req)
		if wire.closed.Load() != int64(p.Attempts-1) {
			t.Fatalf("closed %d of %d dropped responses (or closed the returned one)", wire.closed.Load(), p.Attempts-1)
		}
		retryRead(t, resp, err, 503, fmt.Sprintf("busy-%d", p.Attempts))
		if dock.hits != p.Attempts {
			t.Fatalf("want %d attempts, got %d", p.Attempts, dock.hits)
		}
		retryWant(t, p, jit, pause.waits, p.Attempts-1)
	}
	broken := errors.New("link lost")
	_, srv := retrySetup(t, nil)
	wire, pause := &retryWire{errs: []error{broken, broken, broken, broken}}, &retryPause{}
	r, _ := NewRetrier(wire, Policy{Attempts: 4, Base: time.Millisecond, Max: time.Millisecond}, pause, func(int64) int64 { return 0 })
	req, _ := http.NewRequest("GET", srv.URL, nil)
	if resp, err := r.RoundTrip(req); resp != nil || !errors.Is(err, broken) || wire.calls.Load() != 4 || len(pause.waits) != 3 {
		t.Fatalf("transport errors: %v %v calls=%d waits=%v", resp, err, wire.calls.Load(), pause.waits)
	}
	wire = &retryWire{errs: []error{broken, broken}}
	r, _ = NewRetrier(wire, Policy{Attempts: 4, Base: time.Millisecond, Max: time.Millisecond}, &retryPause{}, func(int64) int64 { return 0 })
	req, _ = http.NewRequest("GET", srv.URL, nil)
	resp, err := r.RoundTrip(req)
	retryRead(t, resp, err, 200, "ok-1")
}

// Unsafe or unreplayable requests get one attempt; keyed POSTs replay the identical body.
func TestHiddenRetryOnlyRepeatsSafeRequests(t *testing.T) {
	jit := func(int64) int64 { return 0 }
	p := Policy{Attempts: 5, Base: time.Millisecond, Max: time.Second}
	for _, code := range []int{400, 404, 409, 500, 501} {
		dock, srv := retrySetup(t, []int{code, code})
		r, _ := NewRetrier(&retryWire{}, p, &retryPause{}, jit)
		req, _ := http.NewRequest("GET", srv.URL, nil)
		resp, err := r.RoundTrip(req)
		retryRead(t, resp, err, code, "busy-1")
		if dock.hits != 1 {
			t.Fatalf("status %d retried %d times", code, dock.hits)
		}
	}
	once := func(method, key string, body io.Reader) *retryDock {
		dock, srv := retrySetup(t, []int{503, 503})
		r, _ := NewRetrier(&retryWire{}, p, &retryPause{}, jit)
		req, _ := http.NewRequest(method, srv.URL+"/fuel", body)
		if key != "" {
			req.Header.Set("Idempotency-Key", key)
		}
		resp, err := r.RoundTrip(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		return dock
	}
	for _, tc := range []struct {
		method, key string
		body        io.Reader
	}{
		{"POST", "", strings.NewReader("deduct 40")},
		{"POST", "", nil},
		{"PATCH", "", strings.NewReader("deduct 40")},
		{"TRACE", "k-1", nil},
		{"PUT", "", io.NopCloser(strings.NewReader("level 9"))},
		{"POST", "k-2", io.NopCloser(strings.NewReader("deduct 40"))},
	} {
		if dock := once(tc.method, tc.key, tc.body); dock.hits != 1 {
			t.Fatalf("%s key=%q body=%T sent %d times", tc.method, tc.key, tc.body, dock.hits)
		}
	}
	rng := rand.New(rand.NewSource(3313))
	for _, method := range []string{"POST", "PATCH", "PUT"} {
		cargo := fmt.Sprintf("deduct %d litres", rng.Intn(1000))
		dock := once(method, fmt.Sprintf("key-%d", rng.Int63()), strings.NewReader(cargo))
		if dock.hits != 3 {
			t.Fatalf("%s with key: %d attempts, want 3", method, dock.hits)
		}
		for i, b := range dock.bodies {
			if b != cargo {
				t.Fatalf("%s attempt %d body %q, want %q", method, i+1, b, cargo)
			}
		}
	}
}

// Cancellation: before the first attempt, during a pause, and inside the transport.
func TestHiddenRetryStopsOnCancel(t *testing.T) {
	p := Policy{Attempts: 5, Base: time.Millisecond, Max: time.Second}
	jit := func(int64) int64 { return 0 }
	_, srv := retrySetup(t, []int{503, 503, 503})
	wire := &retryWire{}
	r, _ := NewRetrier(wire, p, &retryPause{}, jit)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
	if resp, err := r.RoundTrip(req); resp != nil || !errors.Is(err, context.Canceled) || wire.calls.Load() != 0 {
		t.Fatalf("cancelled before start: %v %v calls=%d", resp, err, wire.calls.Load())
	}
	halt := errors.New("pause interrupted")
	dock, srv2 := retrySetup(t, []int{503, 503, 503})
	wire = &retryWire{}
	r, _ = NewRetrier(wire, p, &retryPause{fail: halt}, jit)
	req, _ = http.NewRequest("GET", srv2.URL, nil)
	if resp, err := r.RoundTrip(req); resp != nil || !errors.Is(err, halt) || dock.hits != 1 || wire.closed.Load() != 1 {
		t.Fatalf("pause error: %v %v hits=%d closed=%d", resp, err, dock.hits, wire.closed.Load())
	}
	ctx, cancel = context.WithCancel(context.Background())
	wire = &retryWire{errs: []error{context.Canceled}}
	pause := &retryPause{}
	r, _ = NewRetrier(roundTripFunc(func(req *http.Request) (*http.Response, error) {
		cancel()
		return wire.RoundTrip(req)
	}), p, pause, jit)
	req, _ = http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
	if resp, err := r.RoundTrip(req); resp != nil || err == nil || wire.calls.Load() != 1 || len(pause.waits) != 0 {
		t.Fatalf("cancelled in transport: %v %v calls=%d waits=%v", resp, err, wire.calls.Load(), pause.waits)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// Retry-After on 429/503 replaces the backoff, capped at Max; anything else keeps the backoff.
func TestHiddenBonusRetryHonorsRetryAfter(t *testing.T) {
	p := Policy{Attempts: 2, Base: 100 * time.Millisecond, Max: 5 * time.Second}
	for _, tc := range []struct {
		code  int
		after string
		want  time.Duration
	}{
		{503, "2", 2 * time.Second},
		{429, "0", 0},
		{429, "100", 5 * time.Second},
		{503, "-1", 50*time.Millisecond + 7},
		{503, "soon", 50*time.Millisecond + 7},
		{502, "3", 50*time.Millisecond + 7},
	} {
		dock, srv := retrySetup(t, []int{tc.code})
		dock.header = map[string]string{"Retry-After": tc.after}
		pause := &retryPause{}
		r, _ := NewRetrier(&retryWire{}, p, pause, func(int64) int64 { return 7 })
		req, _ := http.NewRequest("GET", srv.URL, nil)
		resp, err := r.RoundTrip(req)
		retryRead(t, resp, err, 200, "ok-2")
		if len(pause.waits) != 1 || pause.waits[0] != tc.want {
			t.Fatalf("%d Retry-After %q: waited %v, want %v", tc.code, tc.after, pause.waits, tc.want)
		}
	}
}
`

func init() {
	add(&Level{ID: "33-2-reverse-proxy", Chapter: 33, Kind: "lab", LessonRevision: 1, Steps: ch33Lesson332Arabic, Title: "الوكيل العكسي", Goal: "وصّل طلبات القاعدة لوحدات الخلف المسموحة بس، بمهلة وسياسة هيدرات وأخطاء واضحة لما وحدة تقع",
		Requires: []string{"sel:httputil.ReverseProxy"}, Bridge: []int{78}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch332Starter}},
		Lesson: `# الوكيل العكسي لبوابة القاعدة

## قبل ما تكتب
- العميل بيكلم عنوان واحد، والوحدات اللي ورا (المياه، الهوا، المخازن) ليها عناوين داخلية. مين يقرر الطلب يروح فين، ومين يمنع طلب لوحدة مش معروفة؟
- لو العميل بعت بنفسه ~X-Forwarded-For: 10.0.0.1~، الوحدة اللي ورا تصدّق مين؟

## الوكيل العكسي
الوكيل العكسي (reverse proxy) سيرفر بيستقبل الطلب، ويبعته لسيرفر تاني من ورا، ويرجّع رده للعميل. العميل مش شايف غير الوكيل. ده عكس عميل 33-1: هناك كنا بنلف الـ transport، هنا بنقف قدام السيرفرات.

~net/http/httputil~ فيها ~httputil.ReverseProxy~، وده ~http.Handler~ جاهز. أهم حقل فيه ~Rewrite~: دالة بتاخد ~*httputil.ProxyRequest~ فيه ~In~ (الطلب اللي جه، متلمسوش) و ~Out~ (النسخة اللي هتتبعت، عدّلها براحتك).
~~~go
target, _ := url.Parse("http://127.0.0.1:9001")
p := &httputil.ReverseProxy{Rewrite: func(pr *httputil.ProxyRequest) {
	pr.SetURL(target)
}}
http.ListenAndServe(":8080", p)
~~~
~pr.SetURL(target)~ بيغيّر الـ scheme والـ host والـ Host header لعنوان الوحدة، ويضم مسار الهدف قبل مسار الطلب، ويسيب الـ query.

## المسار بعد البادئة
البوابة بتقسم بالبادئة: ~/hydro/~ للمياه. الوحدة نفسها مش عارفة البادئة، فبنشيلها من ~Out~:
~~~go
pr.Out.URL.Path = "/" + strings.TrimPrefix(pr.In.URL.Path, "/hydro/")
pr.Out.URL.RawPath = ""
~~~
لو فيه بادئتين زي ~/hydro/~ و ~/hydro/deep/~، الأطول هي اللي تكسب، عشان الأدق يتطبق. ~RawPath~ بنفضيه عشان ميفضلش ماسك المسار القديم المشفّر.

## allowlist: المعروف بس
الوكيل اللي بيبعت لأي عنوان في الطلب هو باب مفتوح لأي حد يوصل للشبكة الداخلية. الوكيل بتاعنا عنده جدول ثابت من البادئات والأهداف، وأي مسار مش في الجدول يترفض قبل ما يلمس الـ proxy.
~~~go
if best == "" { http.Error(w, "no route", http.StatusNotFound); return }
~~~
الأهداف نفسها بتتفحص وقت الإنشاء: ~url.Parse~ وبعدين ~Scheme~ و ~Host~؛ الخطأ في الجدول يبان عند التشغيل مش عند أول طلب.

## الهيدرات: إيه اللي يعدّي
فيه هيدرات بتخص الاتصال الواحد بس، اسمها hop-by-hop: ~Connection~ و ~Keep-Alive~ و ~Upgrade~ و ~Transfer-Encoding~ وغيرهم، وأي اسم مكتوب جوه ~Connection~. ~ReverseProxy~ بيشيلهم من ~Out~ لوحده، فمتنسخهمش بإيدك.

في وضع ~Rewrite~، الوكيل بيشيل كمان أي ~X-Forwarded-For~ و ~X-Forwarded-Host~ و ~X-Forwarded-Proto~ جاي من العميل، عشان محدش يزوّر عنوانه. ~pr.SetXForwarded()~ بيحطهم من جديد من الحقيقة: عنوان الـ peer، والـ Host الأصلي، و ~http~ أو ~https~.
~~~go
pr.Out.Header.Del("X-Colony-Role") // العميل ميقدرش يدّعي دور
pr.Out.Header.Set("X-Relay-Route", "/hydro/")
pr.SetXForwarded()
~~~
أي هيدر بيدّي صلاحية جوه الشبكة لازم يتشال من الطلب الخارجي.

## الرد قبل ما يرجع
~ModifyResponse~ بتشوف الرد قبل ما يتنسخ للعميل. نشيل منه ~Server~ عشان منكشفش نوع الوحدة. لو رجّعت خطأ، الرد يتقفل ومحدش يشوفه، والخطأ يروح لـ ~ErrorHandler~.
~~~go
ModifyResponse: func(resp *http.Response) error {
	resp.Header.Del("Server")
	return nil
},
~~~

## لما الوحدة تقع
~ErrorHandler~ بتتنادى لما الطلب للوحدة يفشل: الاتصال اترفض، أو المهلة خلصت، أو ~ModifyResponse~ رجّعت خطأ. من غيرها الوكيل بيكتب 502 ورسالة عامة في الـ log.
~~~go
ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.DeadlineExceeded) { http.Error(w, "too slow", 504); return }
	http.Error(w, "down", 502)
},
~~~
502 (Bad Gateway) معناها الوحدة مردتش رد سليم، و 504 (Gateway Timeout) معناها الوحدة مردتش في الوقت.

## المهلة
المهلة بتتحط على سياق الطلب قبل ما يدخل الـ proxy، والـ transport بيقطع الطلب للوحدة لما السياق يخلص:
~~~go
ctx, cancel := context.WithTimeout(r.Context(), timeout)
defer cancel()
proxy.ServeHTTP(w, r.WithContext(ctx))
~~~
نفس السياق ينفع يشيل قيمة صغيرة زي البادئة اللي اتختارت، بـ ~context.WithValue~ ومفتاح خاص زي 26-3، و ~Rewrite~ تقراها من ~pr.In.Context()~.

## أخطاء شائعة
~~~go
target, _ := url.Parse(r.URL.Query().Get("to")) // أي حد يوصل لأي حاجة جوه
pr.Out.Header = pr.In.Header.Clone() // رجّعت Connection و X-Forwarded-For المزوّر
p := &httputil.ReverseProxy{Rewrite: rw, Director: d} // Rewrite و Director مع بعض: الوكيل يرفض الطلب
~~~

## المهمة
نفّذ ~NewRelay(routes map[string]string, timeout time.Duration, transport http.RoundTripper) (http.Handler, error)~. ~routes~ بادئة ← عنوان وحدة. يرجع ~nil,error~ لو ~routes~ فاضية، أو ~timeout~ مش موجبة أو أكبر من دقيقة، أو ~transport~ بـ ~nil~، أو أي بادئة مش بتبدأ وتخلص بـ ~/~ أو طولها أقل من 3، أو أي هدف ميتفكش بـ ~url.Parse~ أو الـ scheme بتاعه مش ~http~ أو ~https~ أو من غير host أو فيه مسار غير ~/~ أو query. خد نسخة من الجدول، واستخدم ~transport~ في الـ proxy.

الـ handler يختار أطول بادئة يبدأ بيها ~r.URL.Path~؛ مفيش: 404 و ~no route\n~ من غير ما يلمس أي وحدة. غير كده يبعت بـ ~httputil.ReverseProxy~ و ~Rewrite~: الهدف بـ ~SetURL~، والمسار ~"/"~ وبعده الباقي بعد البادئة، والـ query زي ما هي، وشيل ~X-Colony-Role~، وحط ~X-Relay-Route~ بالبادئة، و ~SetXForwarded()~. الرد بيرجع بكوده وجسمه وهيدراته من غير ~Server~.

كل طلب عليه مهلة ~timeout~. الوحدة اتأخرت: 504 و ~upstream timeout\n~. الوحدة مش موجودة أو الاتصال فشل: 502 و ~upstream unavailable\n~. كل الأخطاء بـ ~http.Error~؛ مفيش ~Report~، و ~main~ فاضية.

النجمة التالتة: رد من الوحدة كوده 500 أو أكبر ميوصلش للعميل خالص؛ مكانه 502 و ~upstream failed\n~.`,
		Hints: ch33Lesson332HintsArabic,
		En: &LevelText{LessonRevision: 1, Steps: ch33Lesson332English, Title: "Reverse Proxy", Goal: "Relay base requests only to allowed backend units, with a timeout, a header policy, and clear errors when a unit is down", Files: []LabFile{{Name: "main.go", Starter: ch332Starter}},
			Lesson: `# A Reverse Proxy for the Base Gateway

## Before you write
- The client talks to one address, while the units behind it (water, air, storage) have internal addresses. Who decides where a request goes, and who blocks a request for an unknown unit?
- If the client itself sends ~X-Forwarded-For: 10.0.0.1~, whom does the backend unit believe?

## The reverse proxy
A reverse proxy is a server that receives a request, sends it to another server behind it, and returns that server's response to the client. The client only sees the proxy. It is the opposite of the 33-1 client: there we wrapped the transport, here we stand in front of the servers.

~net/http/httputil~ provides ~httputil.ReverseProxy~, a ready ~http.Handler~. Its key field is ~Rewrite~: a function receiving a ~*httputil.ProxyRequest~ with ~In~ (the incoming request, do not touch it) and ~Out~ (the copy that will be sent, change it freely).
~~~go
target, _ := url.Parse("http://127.0.0.1:9001")
p := &httputil.ReverseProxy{Rewrite: func(pr *httputil.ProxyRequest) {
	pr.SetURL(target)
}}
http.ListenAndServe(":8080", p)
~~~
~pr.SetURL(target)~ changes the scheme, host, and Host header to the unit's address, joins the target path before the request path, and keeps the query.

## The path after the prefix
The gateway splits by prefix: ~/hydro/~ for water. The unit itself does not know the prefix, so it is removed from ~Out~:
~~~go
pr.Out.URL.Path = "/" + strings.TrimPrefix(pr.In.URL.Path, "/hydro/")
pr.Out.URL.RawPath = ""
~~~
With two prefixes like ~/hydro/~ and ~/hydro/deep/~, the longest wins, so the more specific one applies. ~RawPath~ is cleared so it does not keep the old encoded path.

## An allowlist: only the known
A proxy that forwards to any address named in the request is an open door into the internal network. Ours has a fixed table of prefixes and targets, and any path not in the table is refused before it touches the proxy.
~~~go
if best == "" { http.Error(w, "no route", http.StatusNotFound); return }
~~~
The targets themselves are checked at construction: ~url.Parse~, then ~Scheme~ and ~Host~; a bad table fails at startup, not on the first request.

## Headers: what gets through
Some headers belong to a single connection only; they are called hop-by-hop: ~Connection~, ~Keep-Alive~, ~Upgrade~, ~Transfer-Encoding~, and others, plus any name listed inside ~Connection~. ~ReverseProxy~ removes them from ~Out~ by itself, so never copy them by hand.

In ~Rewrite~ mode, the proxy also removes any ~X-Forwarded-For~, ~X-Forwarded-Host~, and ~X-Forwarded-Proto~ sent by the client, so nobody can forge an address. ~pr.SetXForwarded()~ sets them again from the facts: the peer address, the original Host, and ~http~ or ~https~.
~~~go
pr.Out.Header.Del("X-Colony-Role") // The client cannot claim a role
pr.Out.Header.Set("X-Relay-Route", "/hydro/")
pr.SetXForwarded()
~~~
Any header that grants privileges inside the network must be removed from an outside request.

## The response before it returns
~ModifyResponse~ sees the response before it is copied to the client. We remove ~Server~ so the unit's software is not revealed. If it returns an error, the response is closed unseen and the error goes to ~ErrorHandler~.
~~~go
ModifyResponse: func(resp *http.Response) error {
	resp.Header.Del("Server")
	return nil
},
~~~

## When a unit is down
~ErrorHandler~ is called when the request to the unit fails: connection refused, timeout expired, or ~ModifyResponse~ returned an error. Without it the proxy writes a 502 and a generic log line.
~~~go
ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.DeadlineExceeded) { http.Error(w, "too slow", 504); return }
	http.Error(w, "down", 502)
},
~~~
502 (Bad Gateway) means the unit gave no valid answer, and 504 (Gateway Timeout) means it did not answer in time.

## The timeout
The timeout goes on the request context before it enters the proxy, and the transport cuts the request to the unit when the context ends:
~~~go
ctx, cancel := context.WithTimeout(r.Context(), timeout)
defer cancel()
proxy.ServeHTTP(w, r.WithContext(ctx))
~~~
The same context can carry a small value like the chosen prefix, with ~context.WithValue~ and a private key as in 26-3, and ~Rewrite~ reads it from ~pr.In.Context()~.

## Common mistakes
~~~go
target, _ := url.Parse(r.URL.Query().Get("to")) // Anyone reaches anything inside
pr.Out.Header = pr.In.Header.Clone() // Restores Connection and the forged X-Forwarded-For
p := &httputil.ReverseProxy{Rewrite: rw, Director: d} // Rewrite and Director together: the proxy refuses requests
~~~

## Your mission
Implement ~NewRelay(routes map[string]string, timeout time.Duration, transport http.RoundTripper) (http.Handler, error)~. ~routes~ maps a prefix to a unit address. Return ~nil,error~ if ~routes~ is empty, ~timeout~ is not positive or above one minute, ~transport~ is ~nil~, any prefix does not start and end with ~/~ or is shorter than 3, or any target fails ~url.Parse~, has a scheme other than ~http~ or ~https~, has no host, has a path other than ~/~, or has a query. Copy the table, and use ~transport~ in the proxy.

The handler picks the longest prefix that ~r.URL.Path~ starts with; none means 404 and ~no route\n~ without touching any unit. Otherwise it forwards with ~httputil.ReverseProxy~ and ~Rewrite~: the target via ~SetURL~, the path ~"/"~ followed by the rest after the prefix, the query unchanged, ~X-Colony-Role~ removed, ~X-Relay-Route~ set to the prefix, and ~SetXForwarded()~. The response returns with its code, body, and headers, without ~Server~.

Every request has the ~timeout~ deadline. A late unit gives 504 and ~upstream timeout\n~. A missing unit or failed connection gives 502 and ~upstream unavailable\n~. All errors use ~http.Error~; no ~Report~, and ~main~ stays empty.

Third star: a unit response with code 500 or above never reaches the client; it is replaced by 502 and ~upstream failed\n~.`,
			Hints: ch33Lesson332HintsEnglish},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch332Hidden},
	})
}

const ch332Starter = `package main

import (
	"errors"
	"net/http"
	"time"
)

func NewRelay(routes map[string]string, timeout time.Duration, transport http.RoundTripper) (http.Handler, error) {
	return nil, errors.New("TODO")
}
func main() {}
`

const ch332Hidden = `package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type relaySeen struct {
	Unit, Method, Path, Query, Host, Body          string
	For, FwdHost, Proto, Role, Route, Secret, Conn string
}

func relayUnit(t *testing.T, name string) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		w.Header().Set("Server", "unit-os/7")
		w.Header().Set("X-Unit", name)
		json.NewEncoder(w).Encode(relaySeen{Unit: name, Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery, Host: r.Host, Body: string(raw),
			For: r.Header.Get("X-Forwarded-For"), FwdHost: r.Header.Get("X-Forwarded-Host"), Proto: r.Header.Get("X-Forwarded-Proto"),
			Role: r.Header.Get("X-Colony-Role"), Route: r.Header.Get("X-Relay-Route"), Secret: r.Header.Get("X-Hatch-Code"), Conn: r.Header.Get("Connection")})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func relayDo(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func relaySeenBy(t *testing.T, w *httptest.ResponseRecorder) relaySeen {
	t.Helper()
	var s relaySeen
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &s) != nil {
		t.Fatalf("want relayed 200, got %d %q", w.Code, w.Body.String())
	}
	return s
}

func relayFail(t *testing.T, w *httptest.ResponseRecorder, code int, msg string) {
	t.Helper()
	if w.Code != code || w.Body.String() != msg || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("want %d %q, got %d %q %v", code, msg, w.Code, w.Body.String(), w.Header())
	}
}

// Boundaries: every invalid table entry, timeout limits, a nil transport.
func TestHiddenRelayRejectsBadConfig(t *testing.T) {
	ok := "http://127.0.0.1:9"
	tr := http.DefaultTransport
	for _, tc := range []struct {
		routes  map[string]string
		timeout time.Duration
		tr      http.RoundTripper
	}{
		{map[string]string{}, time.Second, tr},
		{nil, time.Second, tr},
		{map[string]string{"/hydro/": ok}, 0, tr},
		{map[string]string{"/hydro/": ok}, -time.Second, tr},
		{map[string]string{"/hydro/": ok}, time.Minute + 1, tr},
		{map[string]string{"/hydro/": ok}, time.Second, nil},
		{map[string]string{"hydro/": ok}, time.Second, tr},
		{map[string]string{"/hydro": ok}, time.Second, tr},
		{map[string]string{"/": ok}, time.Second, tr},
		{map[string]string{"/hydro/": "ftp://127.0.0.1:9"}, time.Second, tr},
		{map[string]string{"/hydro/": "http://"}, time.Second, tr},
		{map[string]string{"/hydro/": "http://127.0.0.1:9/tanks"}, time.Second, tr},
		{map[string]string{"/hydro/": "http://127.0.0.1:9?unit=1"}, time.Second, tr},
		{map[string]string{"/hydro/": "::not a url"}, time.Second, tr},
		{map[string]string{"/hydro/": ok, "/air/": "127.0.0.1:9"}, time.Second, tr},
	} {
		if h, err := NewRelay(tc.routes, tc.timeout, tc.tr); h != nil || err == nil {
			t.Fatalf("accepted %v %v", tc.routes, tc.timeout)
		}
	}
	for _, routes := range []map[string]string{{"/a/": ok}, {"/hydro/": ok + "/", "/air/": "https://127.0.0.1:9"}} {
		if h, err := NewRelay(routes, time.Minute, tr); h == nil || err != nil {
			t.Fatalf("rejected %v: %v", routes, err)
		}
	}
}

// Longest prefix wins, the prefix is stripped, queries and bodies survive, unknown paths never leave.
func TestHiddenRelayRoutesByLongestPrefix(t *testing.T) {
	rng := rand.New(rand.NewSource(3321))
	water, deep, air := relayUnit(t, "water"), relayUnit(t, "deep"), relayUnit(t, "air")
	routes := map[string]string{"/hydro/": water.URL, "/hydro/deep/": deep.URL, "/air/": air.URL + "/"}
	h, err := NewRelay(routes, 5*time.Second, http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}
	routes["/hydro/"] = air.URL
	routes["/vault/"] = air.URL
	for i := 0; i < 12; i++ {
		tail := fmt.Sprintf("tank-%d/level", rng.Intn(1000))
		query := fmt.Sprintf("probe=%d&unit=%c", rng.Intn(100), 'a'+rng.Intn(26))
		prefix, unit := []string{"/hydro/", "/hydro/deep/", "/air/"}[i%3], []string{"water", "deep", "air"}[i%3]
		body := fmt.Sprintf("reading %d", rng.Int63())
		method := []string{"GET", "POST", "PUT", "DELETE"}[rng.Intn(4)]
		s := relaySeenBy(t, relayDo(h, method, prefix+tail+"?"+query, body))
		if s.Unit != unit || s.Path != "/"+tail || s.Query != query || s.Method != method || s.Body != body || s.Route != prefix {
			t.Fatalf("%s %s%s: unit saw %+v", method, prefix, tail, s)
		}
	}
	if s := relaySeenBy(t, relayDo(h, "GET", "/hydro/", "")); s.Unit != "water" || s.Path != "/" {
		t.Fatalf("bare prefix: %+v", s)
	}
	for _, path := range []string{"/", "/hydro", "/hydr/x", "/vault/x", "/airlock/x", "/x/hydro/"} {
		relayFail(t, relayDo(h, "GET", path, ""), 404, "no route\n")
	}
}

// Header policy: forged forwarding and role headers dropped, hop-by-hop removed, Server hidden.
func TestHiddenRelayHeaderPolicy(t *testing.T) {
	unit := relayUnit(t, "store")
	h, err := NewRelay(map[string]string{"/store/": unit.URL}, 5*time.Second, http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest("GET", "http://gate.colony/store/crates", nil)
	r.RemoteAddr = "198.51.100.23:40112"
	r.Header.Set("X-Forwarded-For", "10.0.0.1")
	r.Header.Set("X-Forwarded-Host", "evil.example")
	r.Header.Set("X-Colony-Role", "commander")
	r.Header.Set("Connection", "X-Hatch-Code")
	r.Header.Set("X-Hatch-Code", "open-sesame")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	s := relaySeenBy(t, w)
	if s.For != "198.51.100.23" || s.FwdHost != "gate.colony" || s.Proto != "http" {
		t.Fatalf("forwarding headers: %+v", s)
	}
	if s.Role != "" || s.Secret != "" || s.Conn == "X-Hatch-Code" || s.Route != "/store/" {
		t.Fatalf("leaked headers: %+v", s)
	}
	if !strings.HasPrefix(s.Host, "127.0.0.1:") {
		t.Fatalf("unit saw Host %q", s.Host)
	}
	if w.Header().Get("Server") != "" || w.Header().Get("X-Unit") != "store" {
		t.Fatalf("response headers: %v", w.Header())
	}
}

// A dead unit gives 502, a stuck unit gives 504 after the deadline, and 4xx passes through untouched.
func TestHiddenRelayUnitFailures(t *testing.T) {
	dead := httptest.NewServer(http.NotFoundHandler())
	deadURL := dead.URL
	dead.Close()
	release := make(chan struct{})
	stuck := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	t.Cleanup(stuck.Close)
	t.Cleanup(func() { close(release) })
	picky := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Reason", "bad-crate")
		w.WriteHeader(422)
		w.Write([]byte("crate id required"))
	}))
	t.Cleanup(picky.Close)
	h, err := NewRelay(map[string]string{"/dead/": deadURL, "/stuck/": stuck.URL, "/picky/": picky.URL}, 80*time.Millisecond, http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}
	relayFail(t, relayDo(h, "GET", "/dead/x", ""), 502, "upstream unavailable\n")
	start := time.Now()
	relayFail(t, relayDo(h, "GET", "/stuck/x", ""), 504, "upstream timeout\n")
	if time.Since(start) > 3*time.Second {
		t.Fatal("timeout not enforced")
	}
	w := relayDo(h, "POST", "/picky/x", "{}")
	if w.Code != 422 || w.Body.String() != "crate id required" || w.Header().Get("X-Reason") != "bad-crate" {
		t.Fatalf("4xx changed: %d %q %v", w.Code, w.Body.String(), w.Header())
	}
}

// A unit's 5xx body (which may leak internals) is replaced by a clean 502.
func TestHiddenBonusRelayHidesUnitErrors(t *testing.T) {
	for _, code := range []int{500, 503, 599} {
		unit := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(code)
			w.Write([]byte("panic: tank.go:41 nil pointer"))
		}))
		h, err := NewRelay(map[string]string{"/hydro/": unit.URL}, 5*time.Second, http.DefaultTransport)
		if err != nil {
			t.Fatal(err)
		}
		relayFail(t, relayDo(h, "GET", "/hydro/x", ""), 502, "upstream failed\n")
		unit.Close()
	}
	ok := relayUnit(t, "water")
	h, _ := NewRelay(map[string]string{"/hydro/": ok.URL}, 5*time.Second, http.DefaultTransport)
	relaySeenBy(t, relayDo(h, "GET", "/hydro/x", ""))
}
`

func init() {
	add(&Level{ID: "33-3-load-balancer", Chapter: 33, Kind: "lab", LessonRevision: 1, Steps: ch33Lesson333Arabic, Title: "موزّع الأحمال", Goal: "وزّع طلبات المحطة على مولدات متشابهة بالدور، وشيل المولد المريض، وكمّل لو واحد وقع",
		Requires: []string{"mutex"}, Bridge: []int{79}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch333Starter}},
		Lesson: `# موزّع الأحمال على المولدات

## قبل ما تكتب
- عندك تلات مولدات بيعملوا نفس الشغل. لو كل الطلبات راحت لأول واحد، إيه فايدة التانيين؟
- لو مولد وقع في نص الوردية، مين يعرف؟ ومين يرجّعه للخدمة لما يفوق؟

## round-robin: بالدور
موزّع الأحمال (load balancer) وكيل عكسي زي 33-2، بس بدل جدول بادئات عنده مجموعة وحدات متشابهة، وكل طلب يروح لواحدة منهم. أبسط سياسة round-robin: الأول، التاني، التالت، وبعدين الأول تاني.
~~~go
u := units[next%len(units)]
next++
~~~
ده بيوزّع بالتساوي من غير ما يعرف حاجة عن الحمل، وده كفاية لما الطلبات شبه بعض.

## العدّاد مورد مشترك
كل طلب بيتخدم في goroutine لوحده، يعني عشرات الطلبات بتقرا وتزوّد ~next~ في نفس اللحظة. ده نفس سباق فصل 12: من غير قفل، طلبين ياخدوا نفس الوحدة وواحدة تتنسي، و ~-race~ يمسكها.
~~~go
p.mu.Lock()
u := p.units[p.next]
p.next = (p.next + 1) % len(p.units)
p.mu.Unlock()
~~~
الاختيار كله (قراية الحالة، والتخطي، وتحريك العدّاد) لازم يحصل تحت نفس القفل، عشان قرار الطلب يبقى متسق.

## الوحدات المريضة
round-robin بيتخطى الوحدة اللي مش سليمة. التخطي بيبدأ من مكان العدّاد ويلف مرة واحدة على الكل؛ أول سليمة هي اللي تتاخد، والعدّاد يقف بعدها.
~~~go
for i := 0; i < n; i++ {
	u := p.units[(p.next+i)%n]
	if u.healthy { p.next = (p.next + i + 1) % n; return u }
}
return nil // مفيش ولا واحدة سليمة
~~~
لو مفيش ولا وحدة سليمة الرد 503 (Service Unavailable): الخدمة نفسها مش متاحة دلوقتي، مش غلطة العميل.

## فحص الصحة الفعّال
فحص الصحة (health check) طلب صغير بنبعته لكل وحدة على مسار متفق عليه زي ~/healthz~. 200 معناها سليمة؛ أي حاجة تانية أو خطأ معناها مريضة. الفحص نفسه هو اللي بيرجّع الوحدة اللي فاقت.
~~~go
req, _ := http.NewRequestWithContext(ctx, "GET", base+"/healthz", nil)
resp, err := transport.RoundTrip(req)
ok := err == nil && resp.StatusCode == 200
if err == nil { resp.Body.Close() }
~~~
الفحص بيكلم الشبكة وممكن ياخد وقت. متمسكش القفل وانت مستني الرد، وإلا كل الطلبات هتقف ورا الفحص. خد نسخة من القائمة تحت القفل، وافحص من غير قفل، وبعدين اقفل تاني لحظة ما تكتب النتيجة.

في الإنتاج الفحص بيشتغل كل شوية بـ ticker؛ هنا الدالة بتتنادى بإيد التست، عشان مفيش وقت حقيقي.

## فحص الصحة السلبي
مش لازم نستنى الفحص الجاي. لو الطلب الحقيقي فشل في الوصول للوحدة (الـ ~ErrorHandler~ بتاع 33-2 اتنادى)، الوحدة تتعلم مريضة فورًا.
~~~go
ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
	r.Context().Value(pickKey{}).(*attempt).failed = true
},
~~~
لو ~ErrorHandler~ مكتبش حاجة في ~w~، الرد لسه فاضي، فالموزّع يقدر يجرّب وحدة تانية لنفس الطلب. نفس قاعدة 33-1: الإعادة لـ ~GET~ و ~HEAD~ بس؛ أي طريقة تانية ممكن تكون اتنفذت، فترجع 502.

## مين اختار مين
~Rewrite~ لازم تعرف الوحدة اللي اتختارت للمحاولة دي. حطها في السياق قبل ~proxy.ServeHTTP~ بمؤشر لـ struct صغير، و ~ErrorHandler~ يعلّم فيه الفشل. كل محاولة ليها struct جديد، فمفيش مشاركة بين الطلبات.

## أخطاء شائعة
~~~go
u := p.units[p.next%n]; p.next++ // من غير قفل: سباق
p.mu.Lock(); defer p.mu.Unlock(); resp, _ := client.Get(u + "/healthz") // القفل ماسك طول الفحص
http.Error(w, "down", 502); continue // كتبت الرد وبعدين جربت وحدة تانية: الرد اتلخبط
~~~

## المهمة
نفّذ ~NewPool(backends []string, transport http.RoundTripper) (*Pool, error)~: يرجع ~nil,error~ لو القائمة فاضية أو أكتر من 64، أو ~transport~ بـ ~nil~، أو فيه عنوان متكرر، أو أي عنوان مش ~http~ أو ~https~ مع host ومن غير مسار غير ~/~ ومن غير query (نفس قواعد 33-2). كل الوحدات بتبدأ سليمة.

~(p *Pool) ServeHTTP(w,r)~ يختار بالدور أول وحدة سليمة من بعد آخر اختيار (الطلب الأول للوحدة الأولى)، ويبعت بـ ~httputil.ReverseProxy~ بنفس المسار والـ query، و ~SetXForwarded()~. مفيش وحدة سليمة: 503 و ~no healthy backend\n~. فشل الوصول لوحدة: تتعلم مريضة؛ لو الطلب ~GET~ أو ~HEAD~ جرّب الوحدة السليمة الجاية اللي متجرّبتش في الطلب ده، وغير كده 502 و ~backend failed\n~.

~(p *Pool) Check(ctx context.Context) int~ يبعت ~GET /healthz~ لكل وحدة بالـ ~transport~ والسياق، يعلّم 200 سليمة وأي حاجة تانية مريضة، ويقفل الجسم، ويرجّع عدد السليمة. ~(p *Pool) Healthy() []string~ يرجّع العناوين السليمة بترتيب الإدخال وبنفس كتابتها. كله آمن للتوازي تحت قفل واحد ومن غير قفل أثناء الشبكة؛ مفيش ~Report~، و ~main~ فاضية.

النجمة التالتة: رد من وحدة كوده 500 أو أكبر يتحسب فشل زي الاتصال المقطوع: الوحدة تتعلم مريضة، والجسم ده ميوصلش للعميل، و ~GET~ يتجرب على اللي بعدها.`,
		Hints: ch33Lesson333HintsArabic,
		En: &LevelText{LessonRevision: 1, Steps: ch33Lesson333English, Title: "Load Balancer", Goal: "Share station requests across identical generators in turn, drop a sick generator, and keep going when one falls", Files: []LabFile{{Name: "main.go", Starter: ch333Starter}},
			Lesson: `# A Load Balancer for the Generators

## Before you write
- You have three generators doing the same work. If every request goes to the first, what use are the other two?
- If a generator falls mid-shift, who notices? And who brings it back into service once it recovers?

## Round-robin: taking turns
A load balancer is a reverse proxy like 33-2, but instead of a prefix table it has a group of identical units, and each request goes to one of them. The simplest policy is round-robin: first, second, third, then the first again.
~~~go
u := units[next%len(units)]
next++
~~~
It spreads load evenly without knowing anything about it, which is enough when requests are alike.

## The counter is a shared resource
Each request is served in its own goroutine, so dozens of requests read and increment ~next~ at the same instant. It is the chapter 12 race again: without a lock two requests take the same unit and another is skipped, and ~-race~ catches it.
~~~go
p.mu.Lock()
u := p.units[p.next]
p.next = (p.next + 1) % len(p.units)
p.mu.Unlock()
~~~
The whole choice (reading health, skipping, moving the counter) must happen under the same lock, so each request's decision is consistent.

## Sick units
Round-robin skips units that are not healthy. The scan starts at the counter and goes around everyone once; the first healthy unit is taken, and the counter stops just after it.
~~~go
for i := 0; i < n; i++ {
	u := p.units[(p.next+i)%n]
	if u.healthy { p.next = (p.next + i + 1) % n; return u }
}
return nil // Not a single healthy unit
~~~
With no healthy unit the answer is 503 (Service Unavailable): the service itself is unavailable right now, not the client's fault.

## Active health checks
A health check is a small request sent to every unit on an agreed path like ~/healthz~. 200 means healthy; anything else or an error means sick. The check is also what brings a recovered unit back.
~~~go
req, _ := http.NewRequestWithContext(ctx, "GET", base+"/healthz", nil)
resp, err := transport.RoundTrip(req)
ok := err == nil && resp.StatusCode == 200
if err == nil { resp.Body.Close() }
~~~
The check talks to the network and may take time. Do not hold the lock while waiting for the answer, or every request stalls behind the check. Copy the list under the lock, probe without it, then lock again just to write the result.

In production the check runs periodically from a ticker; here the test calls the method directly, so no real time is involved.

## Passive health checks
There is no need to wait for the next check. If a real request fails to reach the unit (the 33-2 ~ErrorHandler~ was called), the unit is marked sick at once.
~~~go
ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
	r.Context().Value(pickKey{}).(*attempt).failed = true
},
~~~
If ~ErrorHandler~ writes nothing to ~w~, the response is still empty, so the balancer can try another unit for the same request. The 33-1 rule holds: retry only ~GET~ and ~HEAD~; any other method may have executed, so it gets a 502.

## Who picked whom
~Rewrite~ must know which unit was chosen for this attempt. Put it in the context before ~proxy.ServeHTTP~ as a pointer to a small struct, and ~ErrorHandler~ marks the failure in it. Each attempt gets a fresh struct, so nothing is shared between requests.

## Common mistakes
~~~go
u := p.units[p.next%n]; p.next++ // No lock: a race
p.mu.Lock(); defer p.mu.Unlock(); resp, _ := client.Get(u + "/healthz") // Lock held for the whole probe
http.Error(w, "down", 502); continue // Wrote the response, then tried another unit: garbled response
~~~

## Your mission
Implement ~NewPool(backends []string, transport http.RoundTripper) (*Pool, error)~: return ~nil,error~ if the list is empty or longer than 64, ~transport~ is ~nil~, an address repeats, or any address is not ~http~ or ~https~ with a host, no path other than ~/~, and no query (the 33-2 rules). All units start healthy.

~(p *Pool) ServeHTTP(w,r)~ picks, in turn, the first healthy unit after the last pick (the first request goes to the first unit), and forwards with ~httputil.ReverseProxy~ keeping the path and query, plus ~SetXForwarded()~. No healthy unit: 503 and ~no healthy backend\n~. A failure to reach a unit marks it sick; for ~GET~ or ~HEAD~ try the next healthy unit not yet tried for this request, otherwise return 502 and ~backend failed\n~.

~(p *Pool) Check(ctx context.Context) int~ sends ~GET /healthz~ to every unit through ~transport~ with the context, marks 200 as healthy and anything else as sick, closes the body, and returns the healthy count. ~(p *Pool) Healthy() []string~ returns the healthy addresses in input order, spelled as given. Everything is concurrency-safe under one lock, never held during network calls; no ~Report~, and ~main~ stays empty.

Third star: a unit response with code 500 or above counts as a failure just like a dropped connection: the unit is marked sick, that body never reaches the client, and a ~GET~ is tried on the next one.`,
			Hints: ch33Lesson333HintsEnglish},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch333Hidden},
	})
}

const ch333Starter = `package main

import (
	"context"
	"errors"
	"net/http"
)

type Pool struct{}

func NewPool(backends []string, transport http.RoundTripper) (*Pool, error) {
	return nil, errors.New("TODO")
}
func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
func (p *Pool) Check(ctx context.Context) int                    { return 0 }
func (p *Pool) Healthy() []string                                { return nil }
func main()                                                      {}
`

const ch333Hidden = `package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type genUnit struct {
	srv    *httptest.Server
	name   string
	health atomic.Int64 // status code for /healthz
	answer atomic.Int64 // status code for everything else
	hits   atomic.Int64
	probes atomic.Int64
}

func genStart(t *testing.T, name string) *genUnit {
	g := &genUnit{name: name}
	g.health.Store(200)
	g.answer.Store(200)
	g.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			g.probes.Add(1)
			w.WriteHeader(int(g.health.Load()))
			return
		}
		g.hits.Add(1)
		w.WriteHeader(int(g.answer.Load()))
		fmt.Fprintf(w, "%s %s?%s %s", g.name, r.URL.Path, r.URL.RawQuery, r.Header.Get("X-Forwarded-For"))
	}))
	t.Cleanup(g.srv.Close)
	return g
}

func genFleet(t *testing.T, n int) ([]*genUnit, []string) {
	var units []*genUnit
	var urls []string
	for i := 0; i < n; i++ {
		g := genStart(t, fmt.Sprintf("gen-%d", i))
		units, urls = append(units, g), append(urls, g.srv.URL)
	}
	return units, urls
}

func genCall(h http.Handler, method, target string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, target, nil)
	r.RemoteAddr = "203.0.113.9:5000"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func genWho(t *testing.T, w *httptest.ResponseRecorder) string {
	t.Helper()
	if w.Code != 200 {
		t.Fatalf("want 200, got %d %q", w.Code, w.Body.String())
	}
	return strings.Fields(w.Body.String())[0]
}

func genFail(t *testing.T, w *httptest.ResponseRecorder, code int, msg string) {
	t.Helper()
	if w.Code != code || w.Body.String() != msg || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("want %d %q, got %d %q", code, msg, w.Code, w.Body.String())
	}
}

func genDeadURL() string {
	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()
	return srv.URL
}

// Boundaries: empty, 65 entries, duplicates, bad addresses, nil transport; 1 and 64 entries accepted.
func TestHiddenPoolRejectsBadFleet(t *testing.T) {
	tr := http.DefaultTransport
	many := make([]string, 65)
	for i := range many {
		many[i] = fmt.Sprintf("http://127.0.0.1:%d", 10000+i)
	}
	for _, list := range [][]string{nil, {}, many, {"http://127.0.0.1:1", "http://127.0.0.1:1"}, {"ftp://127.0.0.1:1"}, {"http://"}, {"http://127.0.0.1:1/x"}, {"http://127.0.0.1:1?a=1"}, {"127.0.0.1:1"}, {"::x"}} {
		if p, err := NewPool(list, tr); p != nil || err == nil {
			t.Fatalf("accepted %v", list)
		}
	}
	if p, err := NewPool([]string{"http://127.0.0.1:1"}, nil); p != nil || err == nil {
		t.Fatal("nil transport accepted")
	}
	for _, list := range [][]string{{"http://127.0.0.1:1/"}, many[:64]} {
		p, err := NewPool(list, tr)
		if p == nil || err != nil || !reflect.DeepEqual(p.Healthy(), list) {
			t.Fatalf("rejected %d units: %v", len(list), err)
		}
	}
}

// Exact turn order, path and query kept, forwarding header set, input list copied.
func TestHiddenPoolTakesTurns(t *testing.T) {
	rng := rand.New(rand.NewSource(3331))
	for round := 0; round < 4; round++ {
		n := 2 + rng.Intn(4)
		_, urls := genFleet(t, n)
		p, err := NewPool(urls, http.DefaultTransport)
		if err != nil {
			t.Fatal(err)
		}
		given := append([]string(nil), urls...)
		urls[0] = "http://127.0.0.1:1"
		for i := 0; i < 3*n+rng.Intn(n); i++ {
			path := fmt.Sprintf("/load/%d", rng.Intn(1000))
			w := genCall(p, []string{"GET", "POST", "PUT"}[rng.Intn(3)], path+"?shift=night")
			want := fmt.Sprintf("gen-%d %s?shift=night 203.0.113.9", i%n, path)
			if w.Code != 200 || w.Body.String() != want {
				t.Fatalf("request %d: got %d %q, want %q", i, w.Code, w.Body.String(), want)
			}
		}
		if !reflect.DeepEqual(p.Healthy(), given) {
			t.Fatalf("Healthy %v, want %v", p.Healthy(), given)
		}
	}
}

// Active checks remove sick units, keep the turn order among the rest, and revive recovered ones.
func TestHiddenPoolActiveHealth(t *testing.T) {
	units, urls := genFleet(t, 4)
	p, _ := NewPool(urls, http.DefaultTransport)
	units[1].health.Store(503)
	units[3].health.Store(404)
	if got := p.Check(context.Background()); got != 2 {
		t.Fatalf("Check = %d, want 2", got)
	}
	for _, g := range units {
		if g.probes.Load() != 1 {
			t.Fatalf("%s probed %d times", g.name, g.probes.Load())
		}
	}
	if !reflect.DeepEqual(p.Healthy(), []string{urls[0], urls[2]}) {
		t.Fatalf("Healthy %v", p.Healthy())
	}
	var order []string
	for i := 0; i < 4; i++ {
		order = append(order, genWho(t, genCall(p, "GET", "/x")))
	}
	if strings.Join(order, ",") != "gen-0,gen-2,gen-0,gen-2" {
		t.Fatalf("order %v", order)
	}
	units[3].health.Store(200)
	if got := p.Check(context.Background()); got != 3 {
		t.Fatalf("revive: Check = %d", got)
	}
	order = order[:0]
	for i := 0; i < 3; i++ {
		order = append(order, genWho(t, genCall(p, "GET", "/x")))
	}
	if strings.Join(order, ",") != "gen-3,gen-0,gen-2" {
		t.Fatalf("order after revive %v", order)
	}
	for _, g := range units {
		g.health.Store(500)
	}
	if got := p.Check(context.Background()); got != 0 || len(p.Healthy()) != 0 {
		t.Fatalf("all sick: Check = %d, Healthy %v", got, p.Healthy())
	}
	before := units[0].hits.Load() + units[2].hits.Load()
	genFail(t, genCall(p, "GET", "/x"), 503, "no healthy backend\n")
	if units[0].hits.Load()+units[2].hits.Load() != before {
		t.Fatal("request reached a sick unit")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	units[0].health.Store(200)
	if got := p.Check(ctx); got != 0 {
		t.Fatalf("cancelled check found %d healthy", got)
	}
}

// Passive failure: a dead unit is marked sick; GET moves on, POST gets 502; all dead gives 503.
func TestHiddenPoolPassiveFailure(t *testing.T) {
	units, urls := genFleet(t, 2)
	dead := genDeadURL()
	p, _ := NewPool([]string{dead, urls[0], urls[1]}, http.DefaultTransport)
	if who := genWho(t, genCall(p, "GET", "/x")); who != "gen-0" {
		t.Fatalf("GET after dead unit went to %s", who)
	}
	if !reflect.DeepEqual(p.Healthy(), urls) {
		t.Fatalf("dead unit still healthy: %v", p.Healthy())
	}
	if w := genCall(p, "HEAD", "/x"); w.Code != 200 || units[1].hits.Load() != 1 {
		t.Fatalf("HEAD: %d, gen-1 hits %d", w.Code, units[1].hits.Load())
	}
	p2, _ := NewPool([]string{dead, urls[0]}, http.DefaultTransport)
	genFail(t, genCall(p2, "POST", "/x"), 502, "backend failed\n")
	if units[0].hits.Load() != 1 || !reflect.DeepEqual(p2.Healthy(), []string{urls[0]}) {
		t.Fatalf("POST retried or dead unit kept: hits=%d healthy=%v", units[0].hits.Load(), p2.Healthy())
	}
	p3, _ := NewPool([]string{dead, genDeadURL()}, http.DefaultTransport)
	genFail(t, genCall(p3, "GET", "/x"), 503, "no healthy backend\n")
	if len(p3.Healthy()) != 0 {
		t.Fatalf("dead units healthy: %v", p3.Healthy())
	}
}

// Parallel requests with concurrent checks: every unit gets exactly its share.
func TestHiddenPoolParallelShare(t *testing.T) {
	units, urls := genFleet(t, 3)
	p, _ := NewPool(urls, http.DefaultTransport)
	var wg sync.WaitGroup
	var bad atomic.Int64
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if w := genCall(p, "GET", "/x"); w.Code != 200 {
				bad.Add(1)
			}
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if p.Check(context.Background()) != 3 || len(p.Healthy()) != 3 {
				bad.Add(1)
			}
		}()
	}
	wg.Wait()
	if bad.Load() != 0 {
		t.Fatalf("%d failed calls", bad.Load())
	}
	for _, g := range units {
		if g.hits.Load() != 50 {
			t.Fatalf("%s served %d of 150", g.name, g.hits.Load())
		}
	}
}

// A 5xx answer is a failure too: the unit is marked sick and its body never reaches the client.
func TestHiddenBonusPoolTreats5xxAsFailure(t *testing.T) {
	units, urls := genFleet(t, 3)
	p, _ := NewPool(urls, http.DefaultTransport)
	units[0].answer.Store(503)
	units[1].answer.Store(500)
	w := genCall(p, "GET", "/x")
	if genWho(t, w) != "gen-2" {
		t.Fatalf("got %q", w.Body.String())
	}
	if !reflect.DeepEqual(p.Healthy(), []string{urls[2]}) {
		t.Fatalf("Healthy %v", p.Healthy())
	}
	units[2].answer.Store(502)
	genFail(t, genCall(p, "POST", "/x"), 502, "backend failed\n")
	genFail(t, genCall(p, "GET", "/x"), 503, "no healthy backend\n")
	if strings.Contains(w.Body.String(), "gen-0") {
		t.Fatal("sick body leaked")
	}
}
`

func init() {
	add(&Level{ID: "33-4-gateway", Chapter: 33, Kind: "lab", LessonRevision: 1, Steps: ch33Lesson334Arabic, Title: "البوابة", Goal: "ابني بوابة واحدة لخدمات القاعدة: توجيه ومفاتيح طاقم وحد طلبات ورقم لكل طلب، كطبقات فوق بعض",
		Requires: []string{"funcvalue"}, Bridge: []int{91}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch334Starter}},
		Lesson: `# بوابة القاعدة

## قبل ما تكتب
- كل خدمة في القاعدة كانت بتفحص المفاتيح وتعدّ الطلبات بنفسها. لو عايز تغيّر قاعدة الحد، هتعدّل كام خدمة؟
- الطلب اللي اترفض بـ 401 يتحسب من حد الطلبات ولا لأ؟ والترتيب ده بيتقرر فين؟

## البوابة: الوكيل اللي بيقرر
بوابة API (API gateway) وكيل عكسي زي 33-2، بس قبل ما يبعت بيعمل الحاجات المشتركة لكل الخدمات مرة واحدة: رقم للطلب، واختيار الخدمة، والتحقق من المفتاح، وحد الطلبات. الخدمات اللي ورا تستلم طلب نضيف ومتحقق منه.

المهمة دي بتجمع حاجات اتعملت قبل كده: رقم الطلب والـ middleware من فصل 26، والتصريح وحد الطلبات بخزان التصاريح من فصل 27 (27-1 و 27-4)، والوكيل من 33-2. الجديد هو تركيبهم في سلسلة واحدة.

## الـ middleware كقيمة
الدالة في Go قيمة: تتخزن في متغير، وتتبعت كمعامل، وتترجع. نوع ~Middleware~ هو دالة بتاخد handler وترجّع handler بيلفه:
~~~go
type Middleware func(http.Handler) http.Handler

var shout Middleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Shout", "yes")
		next.ServeHTTP(w, r)
	})
}
~~~
لأنها قيم، تقدر تحطها في slice وتلف عليها. ~Chain~ بتلف من الآخر للأول، عشان أول middleware في القائمة يبقى الطبقة الخارجية:
~~~go
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }
	return h
}
~~~
الـ method value زي ~g.auth~ برضه قيمة من نوع الدالة دي، وماسكة الـ receiver جواها.

## الترتيب هو السياسة
ترتيب الطبقات مش تفصيلة. في البوابة دي: رقم الطلب الأول عشان كل رد حتى الـ 404 يتعلّم بيه؛ وبعده التوجيه عشان نعرف الخدمة عامة ولا لأ؛ وبعده المفتاح؛ وبعده الحد عشان نعدّ باسم الطاقم مش بالعنوان؛ وفي الآخر الوكيل.
~~~go
h := Chain(proxy, g.requestID, g.route, g.auth, g.limit)
~~~
الطلب بمفتاح غلط بيقف عند ~auth~، فمبيصرفش تصريح من خزان حد. ولو الحد كان قبل المفتاح، أي حد كان يقدر يخلّص تصاريح طاقم تاني بطلبات من غير مفتاح.

## السياق بين الطبقات
كل طبقة بتعرف حاجة، واللي بعدها محتاجاها: التوجيه عرف الخدمة، والمفتاح عرف اسم الطاقم. ~context.WithValue~ بمفتاح خاص (زي 26-3) بيعدّيها لقدام من غير متغيرات عامة:
~~~go
type ctxKey int
const crewCtx ctxKey = 1
next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), crewCtx, "vega")))
crew, _ := r.Context().Value(crewCtx).(string)
~~~

## الهيدرات اللي بتعدّي للخدمة
البوابة هي اللي بتقول للخدمة مين الطاقم، بهيدر ~X-Crew~. يبقى لازم تشيل أي ~X-Crew~ جاي من برة قبل ما تحط بتاعها، وإلا العميل يدّعي إنه أي حد. وتشيل ~X-Api-Key~ كمان: الخدمات مش محتاجة السر، وكل مكان بيوصله السر مكان ممكن يتسرّب منه.
~~~go
pr.Out.Header.Del("X-Crew")
pr.Out.Header.Set("X-Crew", crew)
pr.Out.Header.Set("X-Request-ID", id)
~~~
رقم الطلب بيروح للخدمة كمان، فالـ log بتاع البوابة والخدمة يتربطوا ببعض.

## أخطاء شائعة
~~~go
h := Chain(proxy, g.limit, g.auth) // الحد قبل المفتاح: طلبات من غير مفتاح بتخلّص تصاريح الطاقم
pr.Out.Header.Set("X-Crew", crew) // ولو مفيش طاقم؟ X-Crew المزوّر جاي من العميل فضل زي ما هو
id := r.Header.Get("X-Request-ID") // صدّقت رقم العميل: ممكن يكرر رقم طلب تاني
~~~

## المهمة
الـ starter فيه ~Service{Prefix, Target string; Public bool}~ و ~Config{Services []Service; Keys map[string]string; Burst int; Every time.Duration; Clock labkit.Clock; NextID func() string}~. ~Keys~ مفتاح ← اسم طاقم. نفّذ ~Middleware~ و ~Chain~ زي فوق، و ~NewGateway(c Config) (http.Handler, error)~: يرجع ~nil,error~ لو مفيش خدمات، أو ~Clock~ أو ~NextID~ بـ ~nil~، أو ~Burst~ برة 1 لـ 1000، أو ~Every~ مش موجبة أو أكبر من ساعة، أو بادئة متكررة، أو بادئة أو هدف مخالف لقواعد 33-2، أو مفتاح أو اسم فاضي. خد نسخة من الخدمات والمفاتيح.

الطبقات بالترتيب ده بالظبط: (1) ~NextID()~ مرة واحدة لكل طلب، والقيمة في ~X-Request-ID~ على كل رد، وتجاهل رقم العميل. (2) أطول بادئة؛ مفيش: 404 و ~no route\n~. (3) الخدمة مش ~Public~: ~X-Api-Key~ لازم يكون في ~Keys~، وإلا 401 و ~bad key\n~. (4) خزان تصاريح بقواعد 27-4 (سعة ~Burst~ بيبدأ مليان، تصريح لكل ~Every~ كاملة بساعة ~Clock~)، واحد لكل طاقم، وللخدمات العامة واحد لكل host في ~RemoteAddr~؛ الخزان فاضي: 429 و ~slow down\n~. (5) ~httputil.ReverseProxy~: المسار بعد البادئة زي 33-2، من غير ~X-Api-Key~، و ~X-Crew~ اسم الطاقم أو مشال للعامة، و ~X-Request-ID~ الرقم، و ~SetXForwarded()~. الوحدة مش موجودة: 502 و ~upstream unavailable\n~.

كل الأخطاء بـ ~http.Error~. آمن للتوازي؛ مفيش نوم أو وقت حقيقي، ومفيش ~Report~، و ~main~ فاضية.

النجمة التالتة: رد 429 فيه ~Retry-After~ بالثواني لحد التصريح الجاي، مقرّبة لفوق، زي 27-4.`,
		Hints: ch33Lesson334HintsArabic,
		En: &LevelText{LessonRevision: 1, Steps: ch33Lesson334English, Title: "Gateway", Goal: "Build one gateway for the base services: routing, crew keys, rate limits, and a number for every request, as stacked layers", Files: []LabFile{{Name: "main.go", Starter: ch334Starter}},
			Lesson: `# The Base Gateway

## Before you write
- Every base service used to check keys and count requests by itself. To change the limit rule, how many services do you edit?
- Does a request refused with 401 count against the rate limit? And where is that order decided?

## The gateway: the proxy that decides
An API gateway is a reverse proxy like 33-2, but before forwarding it does the shared work for every service once: a request number, choosing the service, checking the key, and rate limiting. The services behind it receive a clean, verified request.

This mission combines earlier work: request IDs and middleware from chapter 26, passes and token-bucket rate limiting from chapter 27 (27-1 and 27-4), and the proxy from 33-2. What is new is assembling them into one chain.

## Middleware as a value
A function in Go is a value: it is stored in a variable, passed as an argument, and returned. The ~Middleware~ type is a function taking a handler and returning a handler that wraps it:
~~~go
type Middleware func(http.Handler) http.Handler

var shout Middleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Shout", "yes")
		next.ServeHTTP(w, r)
	})
}
~~~
Because they are values, you can put them in a slice and loop over them. ~Chain~ wraps from last to first, so the first middleware in the list becomes the outermost layer:
~~~go
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- { h = mws[i](h) }
	return h
}
~~~
A method value like ~g.auth~ is also a value of that function type, carrying its receiver inside.

## The order is the policy
Layer order is not a detail. In this gateway: the request ID first so every response, even a 404, carries it; then routing to learn whether the service is public; then the key; then the limit so counting is per crew name, not per address; and the proxy last.
~~~go
h := Chain(proxy, g.requestID, g.route, g.auth, g.limit)
~~~
A request with a bad key stops at ~auth~, so it spends no permit from any tank. If the limit came before the key, anyone could drain another crew's permits with keyless requests.

## Context between layers
Each layer learns something the next one needs: routing learns the service, the key check learns the crew name. ~context.WithValue~ with a private key (as in 26-3) carries it forward without globals:
~~~go
type ctxKey int
const crewCtx ctxKey = 1
next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), crewCtx, "vega")))
crew, _ := r.Context().Value(crewCtx).(string)
~~~

## Headers passed to the service
The gateway is what tells the service who the crew is, with an ~X-Crew~ header. So it must remove any ~X-Crew~ from outside before setting its own, or the client claims to be anyone. It also removes ~X-Api-Key~: services do not need the secret, and every place the secret reaches is a place it can leak from.
~~~go
pr.Out.Header.Del("X-Crew")
pr.Out.Header.Set("X-Crew", crew)
pr.Out.Header.Set("X-Request-ID", id)
~~~
The request number goes to the service as well, so gateway and service logs can be joined.

## Common mistakes
~~~go
h := Chain(proxy, g.limit, g.auth) // Limit before key: keyless requests drain the crew's permits
pr.Out.Header.Set("X-Crew", crew) // And with no crew? The client's forged X-Crew stays as is
id := r.Header.Get("X-Request-ID") // Trusted the client's number: it can repeat another request's ID
~~~

## Your mission
The starter has ~Service{Prefix, Target string; Public bool}~ and ~Config{Services []Service; Keys map[string]string; Burst int; Every time.Duration; Clock labkit.Clock; NextID func() string}~. ~Keys~ maps a key to a crew name. Implement ~Middleware~ and ~Chain~ as above, and ~NewGateway(c Config) (http.Handler, error)~: return ~nil,error~ with no services, a ~nil~ ~Clock~ or ~NextID~, ~Burst~ outside 1 to 1000, ~Every~ not positive or above one hour, a repeated prefix, a prefix or target breaking the 33-2 rules, or an empty key or name. Copy services and keys.

The layers, in exactly this order: (1) ~NextID()~ once per request, its value in ~X-Request-ID~ on every response, ignoring the client's number. (2) The longest prefix; none gives 404 and ~no route\n~. (3) A non-~Public~ service needs an ~X-Api-Key~ present in ~Keys~, otherwise 401 and ~bad key\n~. (4) A token tank under the 27-4 rules (capacity ~Burst~ starting full, one permit per complete ~Every~ on ~Clock~), one per crew, and for public services one per ~RemoteAddr~ host; an empty tank gives 429 and ~slow down\n~. (5) ~httputil.ReverseProxy~: the path after the prefix as in 33-2, without ~X-Api-Key~, ~X-Crew~ set to the crew name or removed for public services, ~X-Request-ID~ set to the number, and ~SetXForwarded()~. A missing unit gives 502 and ~upstream unavailable\n~.

All errors use ~http.Error~. Concurrency-safe; no sleeping or real time, no ~Report~, and ~main~ stays empty.

Third star: a 429 response carries ~Retry-After~ in seconds until the next permit, rounded up, as in 27-4.`,
			Hints: ch33Lesson334HintsEnglish},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch334Hidden},
	})
}

const ch334Starter = `package main

import (
	"errors"
	"net/http"
	"time"

	"orbital/colony/labkit"
)

type Service struct {
	Prefix string
	Target string
	Public bool
}

type Config struct {
	Services []Service
	Keys     map[string]string
	Burst    int
	Every    time.Duration
	Clock    labkit.Clock
	NextID   func() string
}

func NewGateway(c Config) (http.Handler, error) { return nil, errors.New("TODO") }
func main()                                     {}
`

const ch334Hidden = `package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"orbital/colony/labkit"
)

type gateIDs struct{ n atomic.Int64 }

func (g *gateIDs) next() string { return fmt.Sprintf("req-%d", g.n.Add(1)) }

func gateUnit(t *testing.T, name string, hits *atomic.Int64) string {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits != nil {
			hits.Add(1)
		}
		crew, has := r.Header["X-Crew"]
		fmt.Fprintf(w, "%s|%s|%s|%v%v|%s|%s", name, r.URL.RequestURI(), r.Header.Get("X-Request-ID"), has, crew, r.Header.Get("X-Api-Key"), r.Header.Get("X-Forwarded-For"))
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func gateCall(h http.Handler, target, key, remote string, extra ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("GET", target, nil)
	r.RemoteAddr = remote
	if key != "" {
		r.Header.Set("X-Api-Key", key)
	}
	for i := 0; i+1 < len(extra); i += 2 {
		r.Header.Set(extra[i], extra[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func gateFail(t *testing.T, w *httptest.ResponseRecorder, code int, msg, id string) {
	t.Helper()
	if w.Code != code || w.Body.String() != msg || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" || w.Header().Get("X-Request-ID") != id {
		t.Fatalf("want %d %q id %s, got %d %q %v", code, msg, id, w.Code, w.Body.String(), w.Header())
	}
}

func gateOK(t *testing.T, w *httptest.ResponseRecorder, body, id string) {
	t.Helper()
	if w.Code != 200 || w.Body.String() != body || w.Header().Values("X-Request-ID")[0] != id || len(w.Header().Values("X-Request-ID")) != 1 {
		t.Fatalf("want 200 %q id %s, got %d %q %v", body, id, w.Code, w.Body.String(), w.Header())
	}
}

func gateConfig(services []Service, keys map[string]string, burst int, every time.Duration, clock labkit.Clock, ids *gateIDs) Config {
	return Config{Services: services, Keys: keys, Burst: burst, Every: every, Clock: clock, NextID: ids.next}
}

// Boundaries: every invalid field; Burst 1 and 1000, Every one hour accepted.
func TestHiddenGatewayRejectsBadConfig(t *testing.T) {
	clock, ids := labkit.NewFakeClock(time.Unix(0, 0)), &gateIDs{}
	ok := []Service{{Prefix: "/hydro/", Target: "http://127.0.0.1:9"}}
	keys := map[string]string{"k-1": "vega"}
	for i, c := range []Config{
		gateConfig(nil, keys, 5, time.Second, clock, ids),
		gateConfig(ok, keys, 5, time.Second, nil, ids),
		{Services: ok, Keys: keys, Burst: 5, Every: time.Second, Clock: clock},
		gateConfig(ok, keys, 0, time.Second, clock, ids),
		gateConfig(ok, keys, 1001, time.Second, clock, ids),
		gateConfig(ok, keys, 5, 0, clock, ids),
		gateConfig(ok, keys, 5, time.Hour+1, clock, ids),
		gateConfig([]Service{ok[0], ok[0]}, keys, 5, time.Second, clock, ids),
		gateConfig([]Service{{Prefix: "/hydro", Target: "http://127.0.0.1:9"}}, keys, 5, time.Second, clock, ids),
		gateConfig([]Service{{Prefix: "/hydro/", Target: "ftp://127.0.0.1:9"}}, keys, 5, time.Second, clock, ids),
		gateConfig([]Service{{Prefix: "/hydro/", Target: "http://127.0.0.1:9/x"}}, keys, 5, time.Second, clock, ids),
		gateConfig(ok, map[string]string{"": "vega"}, 5, time.Second, clock, ids),
		gateConfig(ok, map[string]string{"k-2": ""}, 5, time.Second, clock, ids),
	} {
		if h, err := NewGateway(c); h != nil || err == nil {
			t.Fatalf("case %d accepted", i)
		}
	}
	for _, c := range []Config{gateConfig(ok, nil, 1, time.Hour, clock, ids), gateConfig(ok, keys, 1000, time.Nanosecond, clock, ids)} {
		if h, err := NewGateway(c); h == nil || err != nil {
			t.Fatalf("rejected %+v: %v", c, err)
		}
	}
	if ids.n.Load() != 0 {
		t.Fatal("NextID called while building")
	}
}

// Request IDs on every outcome, routing, keys checked and never forwarded, forged crew dropped.
func TestHiddenGatewayLayers(t *testing.T) {
	rng := rand.New(rand.NewSource(3341))
	var hits atomic.Int64
	water, deep, news := gateUnit(t, "water", &hits), gateUnit(t, "deep", &hits), gateUnit(t, "news", &hits)
	dead := httptest.NewServer(http.NotFoundHandler())
	deadURL := dead.URL
	dead.Close()
	key := fmt.Sprintf("key-%d", rng.Int63())
	services := []Service{{Prefix: "/hydro/", Target: water}, {Prefix: "/hydro/deep/", Target: deep}, {Prefix: "/news/", Target: news, Public: true}, {Prefix: "/dead/", Target: deadURL}}
	keys := map[string]string{key: "vega", "other": "rigel"}
	ids := &gateIDs{}
	h, err := NewGateway(gateConfig(services, keys, 100, time.Second, labkit.NewFakeClock(time.Unix(1_700_000_000, 0)), ids))
	if err != nil {
		t.Fatal(err)
	}
	services[0].Target = news
	keys["late"] = "ghost"
	delete(keys, key)
	remote := "198.51.100.4:7000"
	gateOK(t, gateCall(h, "/hydro/tank?lvl=3", key, remote, "X-Request-ID", "spoof", "X-Crew", "commander"), "water|/tank?lvl=3|req-1|true[vega]||198.51.100.4", "req-1")
	gateOK(t, gateCall(h, "/hydro/deep/pump", "other", remote), "deep|/pump|req-2|true[rigel]||198.51.100.4", "req-2")
	gateOK(t, gateCall(h, "/news/today", "", remote, "X-Crew", "commander"), "news|/today|req-3|false[]||198.51.100.4", "req-3")
	gateOK(t, gateCall(h, "/news/today", key, remote), "news|/today|req-4|false[]||198.51.100.4", "req-4")
	before := hits.Load()
	gateFail(t, gateCall(h, "/vault/x", key, remote), 404, "no route\n", "req-5")
	gateFail(t, gateCall(h, "/hydro/tank", "", remote), 401, "bad key\n", "req-6")
	gateFail(t, gateCall(h, "/hydro/tank", "late", remote), 401, "bad key\n", "req-7")
	gateFail(t, gateCall(h, "/hydro/tank", key+" ", remote), 401, "bad key\n", "req-8")
	if hits.Load() != before {
		t.Fatal("refused request reached a unit")
	}
	gateFail(t, gateCall(h, "/dead/x", key, remote), 502, "upstream unavailable\n", "req-9")
	if ids.n.Load() != 9 {
		t.Fatalf("NextID called %d times for 9 requests", ids.n.Load())
	}
}

// Tanks per crew and per public host; refused and unauthenticated requests spend nothing; exact refill.
func TestHiddenGatewayLimits(t *testing.T) {
	rng := rand.New(rand.NewSource(3342))
	for round := 0; round < 4; round++ {
		burst := 2 + rng.Intn(5)
		every := time.Duration(500+rng.Intn(3000)) * time.Millisecond
		clock := labkit.NewFakeClock(time.Unix(1_700_000_000+int64(rng.Intn(1000)), 0))
		var hits atomic.Int64
		unit := gateUnit(t, "store", &hits)
		h, err := NewGateway(gateConfig([]Service{{Prefix: "/store/", Target: unit}, {Prefix: "/pub/", Target: unit, Public: true}},
			map[string]string{"ka": "altair", "kb": "deneb", "ka2": "altair"}, burst, every, clock, &gateIDs{}))
		if err != nil {
			t.Fatal(err)
		}
		pass := func(w *httptest.ResponseRecorder) {
			t.Helper()
			if w.Code != 200 {
				t.Fatalf("round %d: want pass, got %d %q", round, w.Code, w.Body.String())
			}
		}
		block := func(w *httptest.ResponseRecorder) {
			t.Helper()
			if w.Code != 429 || w.Body.String() != "slow down\n" {
				t.Fatalf("round %d: want 429, got %d %q", round, w.Code, w.Body.String())
			}
		}
		for i := 0; i < burst; i++ {
			gateCall(h, "/store/x", "wrong", "10.0.0.1:1")
			gateCall(h, "/nowhere", "ka", "10.0.0.1:1")
			key := []string{"ka", "ka2"}[i%2]
			pass(gateCall(h, "/store/x", key, fmt.Sprintf("10.0.0.%d:1", i+1)))
		}
		block(gateCall(h, "/store/x", "ka", "10.0.0.9:1"))
		block(gateCall(h, "/store/x", "ka2", "10.0.0.9:1"))
		pass(gateCall(h, "/store/x", "kb", "10.0.0.1:1"))
		for i := 0; i < burst; i++ {
			pass(gateCall(h, "/pub/x", "", "192.0.2.7:"+strconv.Itoa(1000+i)))
		}
		block(gateCall(h, "/pub/x", "", "192.0.2.7:9"))
		pass(gateCall(h, "/pub/x", "", "192.0.2.8:9"))
		clock.Advance(every - 1)
		block(gateCall(h, "/store/x", "ka", "10.0.0.9:1"))
		clock.Advance(1)
		pass(gateCall(h, "/store/x", "ka", "10.0.0.9:1"))
		block(gateCall(h, "/store/x", "ka", "10.0.0.9:1"))
		clock.Advance(time.Duration(burst+3) * every)
		for i := 0; i < burst; i++ {
			pass(gateCall(h, "/store/x", "ka", "10.0.0.9:1"))
		}
		block(gateCall(h, "/store/x", "ka", "10.0.0.9:1"))
		if want := int64(3*burst + 3); hits.Load() != want {
			t.Fatalf("round %d: unit saw %d requests, want %d", round, hits.Load(), want)
		}
	}
}

// Parallel requests from one crew: exactly Burst pass, every ID unique.
func TestHiddenGatewayParallel(t *testing.T) {
	unit := gateUnit(t, "store", nil)
	ids := &gateIDs{}
	h, err := NewGateway(gateConfig([]Service{{Prefix: "/store/", Target: unit}}, map[string]string{"ka": "altair"}, 40, time.Hour, labkit.NewFakeClock(time.Unix(0, 0)), ids))
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	codes := map[int]int{}
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := gateCall(h, "/store/x", "ka", "10.0.0.1:1")
			mu.Lock()
			codes[w.Code]++
			seen[w.Header().Get("X-Request-ID")] = true
			mu.Unlock()
		}()
	}
	wg.Wait()
	if codes[200] != 40 || codes[429] != 60 || len(seen) != 100 {
		t.Fatalf("codes %v, %d distinct ids", codes, len(seen))
	}
}

// Retry-After counts whole seconds until the next permit, rounded up.
func TestHiddenBonusGatewayRetryAfter(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(1_700_000_000, 0))
	unit := gateUnit(t, "store", nil)
	h, _ := NewGateway(gateConfig([]Service{{Prefix: "/store/", Target: unit}}, map[string]string{"ka": "altair"}, 1, 2500*time.Millisecond, clock, &gateIDs{}))
	gateCall(h, "/store/x", "ka", "10.0.0.1:1")
	for _, tc := range []struct {
		advance time.Duration
		want    string
	}{{0, "3"}, {500 * time.Millisecond, "2"}, {time.Second, "1"}, {999 * time.Millisecond, "1"}} {
		clock.Advance(tc.advance)
		w := gateCall(h, "/store/x", "ka", "10.0.0.1:1")
		if w.Code != 429 || w.Header().Get("Retry-After") != tc.want {
			t.Fatalf("after %v: %d Retry-After %q, want %s", tc.advance, w.Code, w.Header().Get("Retry-After"), tc.want)
		}
	}
	clock.Advance(time.Millisecond)
	if w := gateCall(h, "/store/x", "ka", "10.0.0.1:1"); w.Code != 200 || w.Header().Get("Retry-After") != "" {
		t.Fatalf("refill: %d", w.Code)
	}
}
`
