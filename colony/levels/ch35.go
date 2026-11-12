package levels

func init() {
	add(&Level{ID: "35-1-task-queue", Chapter: 35, Kind: "lab", Title: "طابور المهام", Goal: "وزّع مهام الورش على عمّال بـ Redis: عقد إيجار وإعادة محاولة وطابور للمهام الميتة", Setup: seededLabRoom,
		Requires: []string{"sel:json.Marshal"}, Deps: []string{"github.com/redis/go-redis/v9", "github.com/alicebob/miniredis/v2"}, Bridge: []int{86}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch351Starter}},
		Lesson: `# طابور مهام الورش

## قبل ما تكتب
- لو عامل أخد مهمة ووقع قبل ما يخلّصها، مين هيرجّعها للطابور، وإمتى؟
- لو نفس المهمة وصلت مرتين، إزاي تعرف إنها اتسجلت قبل كده؟

## التسليم مرة على الأقل
الطابور الموزّع بيوزّع مهام على عمّال في عمليات مختلفة. العامل ممكن يقع في أي لحظة، فالطابور ميقدرش يمسح المهمة لحظة ما يسلّمها. بيستنى تأكيد (ack) من العامل.

لو التأكيد موصلش، المهمة بترجع تتسلّم تاني. ده اسمه at-least-once: كل مهمة بتتنفذ مرة على الأقل، وممكن أكتر. عشان كده شغل العامل لازم يستحمل التكرار.

## عقد الإيجار (lease)
لما عامل يطلب مهمة (claim) بياخد عقد إيجار ليها لمدة محددة. لو العقد خلص من غير ack، المهمة تعتبر ضايعة وترجع للطابور. مثال: عقد ~30*time.Second~ اتاخد الساعة ~10:00:00~ بينتهي بالظبط ~10:00:30~.

الوقت هنا جاي من ساعة محقونة ~now func() time.Time~ مش من ~time.Now~. التست بيمرر ساعة مزيفة ويقدّمها بنفسه:
~~~go
clock := time.Unix(1000, 0)
now := func() time.Time { return clock }
clock = clock.Add(30 * time.Second) // العقد خلص
~~~

## المجموعة المرتبة في Redis
Redis عندها sorted set: أعضاء نصية كل واحد ليه score رقمي. ~ZAdd~ بتضيف عضو، و ~redis.Z~ هي العضو بالـ score بتاعه. خلّي الـ score هو موعد نهاية العقد بالمللي ثانية. الـ score رقم ~float64~، فالنانوثانية بتاعة تاريخ حقيقي أكبر من إنها تتخزن من غير تقريب:
~~~go
err := client.ZAdd(ctx, "dock:leased", redis.Z{Score: float64(deadline.UnixMilli()), Member: "pump#0"}).Err()
~~~
~ZRangeByScore~ بترجع الأعضاء اللي الـ score بتاعهم في مدى. ~redis.ZRangeBy~ فيها ~Min~ و ~Max~ كنصوص، و ~"-inf"~ معناها من غير حد أدنى:
~~~go
ids, err := client.ZRangeByScore(ctx, "dock:leased", &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(now, 10)}).Result()
~~~
~ZRem~ بتشيل عضو وترجع عدد اللي اتشال فعلاً. لو عاملين حاولوا يشيلوا نفس العضو، واحد بس هياخد ~1~ والتاني ~0~. ده اللي بيخلّي استرجاع العقد المنتهي يحصل مرة واحدة حتى مع عمّال متوازيين.

## القوائم والمجموعات
~RPush~ بتحط في آخر قائمة، و ~LPop~ بتسحب من أولها بشكل ذرّي؛ قائمة فاضية بترجع ~redis.Nil~. ~LRange(ctx, key, 0, -1)~ بترجع القائمة كلها من غير ما تغيرها:
~~~go
client.RPush(ctx, "dock:ready", "pump")
id, err := client.LPop(ctx, "dock:ready").Result()
~~~
~SAdd~ بتضيف لمجموعة وترجع ~1~ لو العضو جديد و ~0~ لو موجود. ده أبسط مفتاح idempotency: الـ ~ID~ اللي اتشاف قبل كده ميتضافش تاني:
~~~go
added, err := client.SAdd(ctx, "dock:seen", "pump").Result()
~~~

## تخزين المهمة
خزّن كل مهمة كـ JSON في مفتاح لوحدها، والقوائم تشيل الـ ~ID~ بس. ~json.Marshal(t)~ بترجع البايتات، و ~client.Set(ctx, key, raw, 0)~ صفر معناه من غير انتهاء.

## العامل المتأخر
عامل اتأخر لحد ما عقده خلص، والمهمة راحت لعامل تاني. لو الأول بعت ack بالـ ~ID~ لوحده هيمسح عقد العامل التاني. الحل إن العضو في المجموعة المرتبة يبقى ~ID~ زائد عدد المحاولات، زي ~"pump#1"~. بعد الاسترجاع العدد بيزيد، فالعامل المتأخر ماسك اسم قديم ومش هيلاقيه.

## طابور المهام الميتة (DLQ)
مهمة بتفشل كل مرة متفضلش تلف للأبد. بعد عدد محاولات محدد بتروح لقائمة منفصلة اسمها dead-letter queue، والمشغّل يبص عليها بنفسه.

## أخطاء شائعة
~~~go
client.LPop(ctx, ready)
client.Del(ctx, "dock:task:"+id) // غلط: المهمة اتمسحت قبل الـ ack، ولو العامل وقع ضاعت
if q.now().After(deadline) { retry(id) } // غلط: عاملين ممكن يسترجعوا نفس المهمة؛ استخدم نتيجة ZRem
~~~

## المهمة
نفّذ ~NewQueue(client, name, lease, maxAttempts, now)~ والدوال ~Enqueue~ و ~Claim~ و ~Ack~ و ~Fail~ و ~Dead~ بتوقيعات الـ starter، وسيب ~main~ فاضية. كل مفاتيح Redis تبدأ بـ ~name~ زائد ~:~ عشان طابورين بأسماء مختلفة ميخلطوش. ~maxAttempts~ دايماً ~1~ أو أكتر، والـ ~lease~ موجبة.

~Enqueue~ بترفض ~ID~ فاضي بـ ~ErrTask~ قبل أي I/O. ~ID~ اتشاف قبل كده (حتى لو خلص) يرجع ~false~ و ~nil~ من غير تغيير. المهمة الجديدة تبدأ بـ ~Attempts~ صفر وتروح لآخر الطابور وترجع ~true~.

~Claim~ الأول تسترجع كل عقد نهايته ~<=~ الوقت الحالي: تزوّد ~Attempts~ بواحد وترجعها لآخر الطابور، أو للـ DLQ لو ~Attempts~ وصل ~maxAttempts~. بعدين تسحب أول مهمة جاهزة وتسجل عقدها حتى ~now() + lease~ (بدقة المللي ثانية) وترجعها بـ ~Attempts~ الحالي. مفيش مهمة جاهزة يرجع ~ErrEmpty~.

~Ack(ctx, t)~ و ~Fail(ctx, t)~ بياخدوا المهمة زي ما ~Claim~ رجّعتها. لو العقد ده مش موجود (اتأكد قبل كده، اتسترجع، أو مهمة مجهولة) يرجعوا ~ErrLease~ من غير أي تغيير. ~Ack~ تنهي المهمة، و ~Fail~ تحسبها محاولة فاشلة بنفس قاعدة الاسترجاع. العقد المنتهي اللي لسه محدش استرجعه يقبل ~Ack~.

~Dead~ ترجع مهام الـ DLQ بترتيب وصولها، كل واحدة بالـ ~Body~ والـ ~Attempts~ بتوعها، وشريحة فاضية مش ~nil~ لو مفيش. كذا ~Queue~ على نفس السيرفر وعمّال متوازيين: كل مهمة تتسلّم لعامل واحد في كل مرة.

النجمة التالتة: ~Redrive(ctx)~ ترجع كل مهام الـ DLQ لآخر الطابور بنفس الترتيب وتصفّر ~Attempts~ وترجع عددها.
`,
		Hints: []string{"ابدأ بـ Enqueue: SAdd على مجموعة seen، واحفظ المهمة JSON، وبعدين RPush.", "في Claim: LPop من ready، و ZAdd باسم ID#Attempts والـ score موعد النهاية؛ Ack تشيله بـ ZRem ولو رجع 0 يبقى ErrLease.", "قبل السحب استرجع بـ ZRangeByScore لحد الوقت الحالي، واللي ZRem بتاعه رجع 1 بس هو اللي يزوّد Attempts ويرجّعها أو يبعتها للـ DLQ."},
		En: &LevelText{Title: "Task Queue", Goal: "Distribute workshop tasks to workers through Redis: leases, retries, and a dead-letter queue", Files: []LabFile{{Name: "main.go", Starter: ch351Starter}},
			Lesson: `# The workshop task queue

## Before you write
- If a worker takes a task and crashes before finishing, who returns it to the queue, and when?
- If the same task arrives twice, how do you know it was already recorded?

## At-least-once delivery
A distributed queue hands tasks to workers in different processes. A worker can crash at any moment, so the queue cannot delete a task the moment it hands it out. It waits for an acknowledgement (ack) from the worker.

If the ack never arrives, the task is delivered again. This is at-least-once delivery: each task runs at least once, possibly more. Worker code must therefore tolerate repeats.

## Leases
When a worker claims a task, it receives a lease for a fixed duration. If the lease ends without an ack, the task is considered lost and returns to the queue. Example: a ~30*time.Second~ lease taken at ~10:00:00~ ends exactly at ~10:00:30~.

Time comes from an injected clock ~now func() time.Time~, not ~time.Now~. The test passes a fake clock and advances it itself:
~~~go
clock := time.Unix(1000, 0)
now := func() time.Time { return clock }
clock = clock.Add(30 * time.Second) // the lease ended
~~~

## Redis sorted sets
A Redis sorted set holds string members, each with a numeric score. ~ZAdd~ adds a member, and ~redis.Z~ is a member with its score. Make the score the lease deadline in milliseconds. A score is a ~float64~, so nanoseconds of a real date are too large to store without rounding:
~~~go
err := client.ZAdd(ctx, "dock:leased", redis.Z{Score: float64(deadline.UnixMilli()), Member: "pump#0"}).Err()
~~~
~ZRangeByScore~ returns members whose score falls in a range. ~redis.ZRangeBy~ has string ~Min~ and ~Max~, and ~"-inf"~ means no lower bound:
~~~go
ids, err := client.ZRangeByScore(ctx, "dock:leased", &redis.ZRangeBy{Min: "-inf", Max: strconv.FormatInt(now, 10)}).Result()
~~~
~ZRem~ removes a member and returns how many were actually removed. If two workers remove the same member, exactly one gets ~1~ and the other ~0~. That makes reclaiming an expired lease happen once, even with concurrent workers.

## Lists and sets
~RPush~ appends to a list and ~LPop~ atomically takes from its front; an empty list returns ~redis.Nil~. ~LRange(ctx, key, 0, -1)~ returns the whole list without changing it:
~~~go
client.RPush(ctx, "dock:ready", "pump")
id, err := client.LPop(ctx, "dock:ready").Result()
~~~
~SAdd~ adds to a set and returns ~1~ when the member is new and ~0~ when it exists. This is the simplest idempotency key: an ~ID~ seen before is not added again:
~~~go
added, err := client.SAdd(ctx, "dock:seen", "pump").Result()
~~~

## Storing a task
Store each task as JSON under its own key; the lists hold only the ~ID~. ~json.Marshal(t)~ returns the bytes, and in ~client.Set(ctx, key, raw, 0)~ zero means no expiry.

## The late worker
A worker stalls until its lease ends, and the task goes to another worker. If the first one acks by ~ID~ alone, it deletes the second worker's lease. The fix: the sorted-set member is the ~ID~ plus the attempt count, like ~"pump#1"~. Reclaiming increments the count, so the late worker holds an old name that no longer exists.

## Dead-letter queue (DLQ)
A task that fails every time must not circulate forever. After a fixed number of attempts it moves to a separate list called the dead-letter queue, where an operator inspects it.

## Common mistakes
~~~go
client.LPop(ctx, ready)
client.Del(ctx, "dock:task:"+id) // Wrong: deleted before the ack; a crash loses it
if q.now().After(deadline) { retry(id) } // Wrong: two workers may reclaim the same task; use ZRem's result
~~~

## Your mission
Implement ~NewQueue(client, name, lease, maxAttempts, now)~ and ~Enqueue~, ~Claim~, ~Ack~, ~Fail~, ~Dead~ with the starter signatures; leave ~main~ empty. Every Redis key starts with ~name~ plus ~:~ so queues with different names never mix. ~maxAttempts~ is always ~1~ or more and ~lease~ is positive.

~Enqueue~ rejects an empty ~ID~ with ~ErrTask~ before any I/O. An ~ID~ seen before (even if finished) returns ~false~ and ~nil~ without changes. A new task starts with ~Attempts~ zero, goes to the back of the queue, and returns ~true~.

~Claim~ first reclaims every lease whose deadline is ~<=~ the current time: it increments ~Attempts~ by one and returns the task to the back of the queue, or to the DLQ when ~Attempts~ reaches ~maxAttempts~. Then it takes the first ready task, records its lease until ~now() + lease~ (millisecond precision), and returns it with its current ~Attempts~. No ready task returns ~ErrEmpty~.

~Ack(ctx, t)~ and ~Fail(ctx, t)~ take the task exactly as ~Claim~ returned it. If that lease does not exist (already acked, reclaimed, or an unknown task) they return ~ErrLease~ without changes. ~Ack~ finishes the task; ~Fail~ counts one failed attempt with the same rule as reclaiming. An expired lease nobody has reclaimed yet still accepts ~Ack~.

~Dead~ returns DLQ tasks in arrival order, each with its ~Body~ and ~Attempts~, and an empty non-nil slice when there are none. Several ~Queue~ values on one server and concurrent workers: each task is held by one worker at a time.

Third star: ~Redrive(ctx)~ moves every DLQ task back to the end of the queue in order, resets ~Attempts~ to zero, and returns the count.
`, Hints: []string{"Start with Enqueue: SAdd to a seen set, store the task as JSON, then RPush.", "In Claim: LPop from ready and ZAdd the name ID#Attempts scored by the deadline; Ack removes it with ZRem, and 0 means ErrLease.", "Before popping, reclaim with ZRangeByScore up to now; only the caller whose ZRem returned 1 increments Attempts and requeues or dead-letters it."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch351Hidden},
	})
}

const ch351Starter = `package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Task struct {
	ID       string
	Body     string
	Attempts int
}

var ErrTask = errors.New("invalid task")
var ErrEmpty = errors.New("queue empty")
var ErrLease = errors.New("lease not held")

type Queue struct{}

func NewQueue(client *redis.Client, name string, lease time.Duration, maxAttempts int, now func() time.Time) *Queue {
	return &Queue{}
}
func (q *Queue) Enqueue(ctx context.Context, t Task) (bool, error) { return false, errors.New("TODO") }
func (q *Queue) Claim(ctx context.Context) (Task, error)           { return Task{}, errors.New("TODO") }
func (q *Queue) Ack(ctx context.Context, t Task) error             { return errors.New("TODO") }
func (q *Queue) Fail(ctx context.Context, t Task) error            { return errors.New("TODO") }
func (q *Queue) Dead(ctx context.Context) ([]Task, error)          { return nil, errors.New("TODO") }
func (q *Queue) Redrive(ctx context.Context) (int, error)          { return 0, errors.New("TODO") }
func main()                                                        {}
`

const ch351Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

type hiddenClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *hiddenClock) Now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.now }
func (c *hiddenClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}

func hiddenRedis(t *testing.T) *miniredis.Miniredis {
	t.Helper()
	srv := miniredis.NewMiniRedis()
	conn, err := srv.Dial()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	t.Cleanup(srv.Close)
	return srv
}
func hiddenClient(t *testing.T, srv *miniredis.Miniredis) *redis.Client {
	t.Helper()
	client := redis.NewClient(&redis.Options{Dialer: func(context.Context, string, string) (net.Conn, error) { return srv.Dial() }, MaxRetries: -1})
	t.Cleanup(func() { client.Close() })
	return client
}
func hiddenQueue(t *testing.T, max int) (*miniredis.Miniredis, *hiddenClock, *Queue) {
	srv := hiddenRedis(t)
	clock := &hiddenClock{now: time.Unix(1700000000, 0)}
	return srv, clock, NewQueue(hiddenClient(t, srv), "dock", 30*time.Second, max, clock.Now)
}
func hiddenClaim(t *testing.T, q *Queue) Task {
	t.Helper()
	got, err := q.Claim(context.Background())
	if err != nil {
		t.Fatal("claim", err)
	}
	return got
}

func TestHiddenQueueFlowAndIdempotency(t *testing.T) {
	_, _, q := hiddenQueue(t, 3)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(3501))
	var want []Task
	for i := 0; i < 15; i++ {
		task := Task{ID: fmt.Sprintf("job-%d-%d", i, rng.Intn(9999)), Body: fmt.Sprintf("weld %d", rng.Intn(100)), Attempts: 7}
		if ok, err := q.Enqueue(ctx, task); !ok || err != nil {
			t.Fatal("enqueue", ok, err)
		}
		if ok, err := q.Enqueue(ctx, Task{ID: task.ID, Body: "duplicate"}); ok || err != nil {
			t.Fatal("duplicate accepted", ok, err)
		}
		task.Attempts = 0
		want = append(want, task)
	}
	for _, w := range want {
		got := hiddenClaim(t, q)
		if got != w {
			t.Fatalf("FIFO claim: got %+v want %+v", got, w)
		}
		if err := q.Ack(ctx, got); err != nil {
			t.Fatal("ack", err)
		}
		if err := q.Ack(ctx, got); !errors.Is(err, ErrLease) {
			t.Fatal("double ack", err)
		}
	}
	if _, err := q.Claim(ctx); !errors.Is(err, ErrEmpty) {
		t.Fatal("empty", err)
	}
	if ok, err := q.Enqueue(ctx, want[0]); ok || err != nil {
		t.Fatal("finished task re-enqueued", ok, err)
	}
	if _, err := q.Enqueue(ctx, Task{Body: "x"}); !errors.Is(err, ErrTask) {
		t.Fatal("empty id", err)
	}
	if err := q.Ack(ctx, Task{ID: "ghost"}); !errors.Is(err, ErrLease) {
		t.Fatal("unknown ack", err)
	}
	if err := q.Fail(ctx, Task{ID: "ghost"}); !errors.Is(err, ErrLease) {
		t.Fatal("unknown fail", err)
	}
	if dead, err := q.Dead(ctx); err != nil || dead == nil || len(dead) != 0 {
		t.Fatal("dead should be empty non-nil", dead, err)
	}
}

func TestHiddenQueueRetryAndDeadLetter(t *testing.T) {
	_, _, q := hiddenQueue(t, 3)
	ctx := context.Background()
	q.Enqueue(ctx, Task{ID: "pump", Body: "fix pump"})
	q.Enqueue(ctx, Task{ID: "fan", Body: "fix fan"})
	for attempt := 0; attempt < 3; attempt++ {
		got := hiddenClaim(t, q)
		if got.ID != "pump" || got.Attempts != attempt {
			t.Fatal("pump attempt", attempt, got)
		}
		if err := q.Fail(ctx, got); err != nil {
			t.Fatal("fail", err)
		}
		if err := q.Fail(ctx, got); !errors.Is(err, ErrLease) {
			t.Fatal("double fail", err)
		}
		fan := hiddenClaim(t, q)
		if fan.ID != "fan" {
			t.Fatal("failed task should go to the back", fan)
		}
		if err := q.Fail(ctx, fan); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := q.Claim(ctx); !errors.Is(err, ErrEmpty) {
		t.Fatal("dead tasks still ready", err)
	}
	dead, err := q.Dead(ctx)
	if err != nil || len(dead) != 2 || dead[0] != (Task{"pump", "fix pump", 3}) || dead[1] != (Task{"fan", "fix fan", 3}) {
		t.Fatal("dead letters", dead, err)
	}
	_, _, one := hiddenQueue(t, 1)
	one.Enqueue(ctx, Task{ID: "once"})
	if err := one.Fail(ctx, hiddenClaim(t, one)); err != nil {
		t.Fatal(err)
	}
	if dead, _ := one.Dead(ctx); len(dead) != 1 || dead[0].Attempts != 1 {
		t.Fatal("maxAttempts 1", dead)
	}
}

func TestHiddenQueueCrashedWorkerLease(t *testing.T) {
	_, clock, q := hiddenQueue(t, 2)
	ctx := context.Background()
	q.Enqueue(ctx, Task{ID: "drill", Body: "core"})
	crashed := hiddenClaim(t, q)
	clock.Advance(30*time.Second - time.Millisecond)
	if _, err := q.Claim(ctx); !errors.Is(err, ErrEmpty) {
		t.Fatal("lease reclaimed early", err)
	}
	clock.Advance(time.Millisecond)
	second := hiddenClaim(t, q)
	if second.ID != "drill" || second.Attempts != 1 {
		t.Fatal("lease boundary reclaim", second)
	}
	if err := q.Ack(ctx, crashed); !errors.Is(err, ErrLease) {
		t.Fatal("late worker acked a reclaimed task", err)
	}
	clock.Advance(time.Hour)
	if _, err := q.Claim(ctx); !errors.Is(err, ErrEmpty) {
		t.Fatal("exhausted lease should dead-letter", err)
	}
	if dead, _ := q.Dead(ctx); len(dead) != 1 || dead[0].Attempts != 2 {
		t.Fatal("dead after expiries", dead)
	}
	if err := q.Ack(ctx, second); !errors.Is(err, ErrLease) {
		t.Fatal("ack after dead-letter", err)
	}
	q.Enqueue(ctx, Task{ID: "slow"})
	slow := hiddenClaim(t, q)
	clock.Advance(time.Minute)
	if err := q.Ack(ctx, slow); err != nil {
		t.Fatal("expired but unreclaimed lease should ack", err)
	}
	if _, err := q.Claim(ctx); !errors.Is(err, ErrEmpty) {
		t.Fatal("acked task came back", err)
	}
}

func TestHiddenQueueNamespacesAndConcurrentWorkers(t *testing.T) {
	srv := hiddenRedis(t)
	clock := &hiddenClock{now: time.Unix(1700000000, 0)}
	ctx := context.Background()
	other := NewQueue(hiddenClient(t, srv), "lab", time.Minute, 3, clock.Now)
	other.Enqueue(ctx, Task{ID: "j-0"})
	producer := NewQueue(hiddenClient(t, srv), "dock", time.Minute, 3, clock.Now)
	const n = 80
	for i := 0; i < n; i++ {
		if ok, err := producer.Enqueue(ctx, Task{ID: fmt.Sprintf("j-%d", i)}); !ok || err != nil {
			t.Fatal("namespaces mixed", ok, err)
		}
	}
	var mu sync.Mutex
	seen := map[string]int{}
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		worker := NewQueue(hiddenClient(t, srv), "dock", time.Minute, 3, clock.Now)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				got, err := worker.Claim(ctx)
				if errors.Is(err, ErrEmpty) {
					return
				}
				if err != nil {
					t.Error(err)
					return
				}
				mu.Lock()
				seen[got.ID]++
				mu.Unlock()
				if err := worker.Ack(ctx, got); err != nil {
					t.Error("ack", err)
				}
			}
		}()
	}
	wg.Wait()
	if len(seen) != n {
		t.Fatal("tasks lost", len(seen))
	}
	for id, c := range seen {
		if c != 1 {
			t.Fatal("claimed twice", id, c)
		}
	}
	if got := hiddenClaim(t, other); got.ID != "j-0" {
		t.Fatal("other namespace", got)
	}
}

func TestHiddenQueueConcurrentReclaimOnce(t *testing.T) {
	srv := hiddenRedis(t)
	clock := &hiddenClock{now: time.Unix(1700000000, 0)}
	ctx := context.Background()
	q := NewQueue(hiddenClient(t, srv), "dock", time.Second, 5, clock.Now)
	const n = 30
	for i := 0; i < n; i++ {
		q.Enqueue(ctx, Task{ID: fmt.Sprintf("r-%d", i)})
	}
	for i := 0; i < n; i++ {
		hiddenClaim(t, q)
	}
	clock.Advance(time.Second)
	var mu sync.Mutex
	seen := map[string]int{}
	var wg sync.WaitGroup
	for w := 0; w < 6; w++ {
		worker := NewQueue(hiddenClient(t, srv), "dock", time.Hour, 5, clock.Now)
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				got, err := worker.Claim(ctx)
				if errors.Is(err, ErrEmpty) {
					return
				}
				if err != nil || got.Attempts != 1 {
					t.Error("reclaimed attempts", got, err)
					return
				}
				mu.Lock()
				seen[got.ID]++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if len(seen) != n {
		t.Fatal("reclaim lost tasks", len(seen))
	}
	for id, c := range seen {
		if c != 1 {
			t.Fatal("reclaimed twice", id, c)
		}
	}
}

func TestHiddenBonusQueueRedrive(t *testing.T) {
	_, _, q := hiddenQueue(t, 1)
	ctx := context.Background()
	for _, id := range []string{"a", "b", "c"} {
		q.Enqueue(ctx, Task{ID: id, Body: "body-" + id})
		q.Fail(ctx, hiddenClaim(t, q))
	}
	n, err := q.Redrive(ctx)
	if err != nil || n != 3 {
		t.Fatal("redrive", n, err)
	}
	if dead, _ := q.Dead(ctx); len(dead) != 0 {
		t.Fatal("dead not drained", dead)
	}
	for _, id := range []string{"a", "b", "c"} {
		got := hiddenClaim(t, q)
		if got != (Task{id, "body-" + id, 0}) {
			t.Fatal("redriven task", got)
		}
	}
	if n, err := q.Redrive(ctx); n != 0 || err != nil {
		t.Fatal("empty redrive", n, err)
	}
}
`

func init() {
	add(&Level{ID: "35-2-scheduler", Chapter: 35, Kind: "lab", Title: "المجدول", Goal: "شغّل مهام الصيانة الدورية بطابور أولويات وساعة مزيفة، من غير تداخل وبسياسة واضحة للمواعيد الفايتة", Setup: seededLabRoom,
		Requires: []string{"sel:heap.Push"}, Deps: []string{}, Bridge: []int{90}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch352Starter}},
		Lesson: `# مجدول الصيانة

## قبل ما تكتب
- لو عندك ألف مهمة، إزاي تعرف أقرب واحدة عليها الدور من غير ما تلف عليهم كلهم كل مرة؟
- لو المجدول وقف ساعة ورجع، المواعيد اللي فاتت تتنفذ كلها، ولا مرة، ولا خالص؟

## طابور الأولويات
المجدول محتاج دايماً «أقرب موعد». الـ heap شجرة مترتبة جوه شريحة: أصغر عنصر دايماً في الخانة ~0~، والإضافة والسحب بياخدوا وقت ~log n~ بدل ما ترتب الشريحة كلها.

الحزمة ~container/heap~ بتشغل الخوارزمية على أي نوع بينفّذ ~heap.Interface~: ~Len~ و ~Less~ و ~Swap~ زي الترتيب، وكمان ~Push(x any)~ و ~Pop() any~ بيكبّروا ويصغّروا الشريحة من آخرها:
~~~go
type minInts []int

func (h minInts) Len() int           { return len(h) }
func (h minInts) Less(i, j int) bool { return h[i] < h[j] }
func (h minInts) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minInts) Push(x any)        { *h = append(*h, x.(int)) }
func (h *minInts) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}
~~~
~Push~ و ~Pop~ بتوعك بيتنادوا من الحزمة بس. إنت بتنادي ~heap.Push(&h, 5)~ و ~heap.Pop(&h).(int)~ اللي بيرتبوا الشجرة. ~any~ اسم تاني لـ ~interface{}~، فلازم type assertion بعد ~Pop~. ~h[0]~ بتبص على الأصغر من غير ما تسحبه.

## ترتيب ثابت
لو مهمتين ليهم نفس الموعد، لازم ~Less~ تقرر بقاعدة ثابتة وإلا ترتيب التشغيل يبقى عشوائي. هنا: الموعد الأبدر، بعدين الأولوية الأعلى، بعدين الاسم الأصغر أبجدياً. ~a.Equal(b)~ و ~a.Before(b)~ بيقارنوا ~time.Time~.

## ساعة محقونة
المجدول مبيقراش ساعة الجهاز ولا بينام. اللي بيشغّله بينادي ~Tick(now)~ بالوقت اللي هو عايزه، والتست بيقدّم الوقت بنفسه. كده نفس المدخلات دايماً بتدّي نفس التشغيلات.

## الشبكة والمواعيد الفايتة
مهمة كل ~10*time.Second~ ابتدت من ~start~ مواعيدها ~start+10s~ و ~start+20s~ وهكذا. لو ~Tick~ اتنادت متأخر عند ~start+45s~ يبقى فات ~10s~ و ~20s~ و ~30s~ و ~40s~. القسمة على المدة بتحسب كام موعد فات بعد الأول:
~~~go
missed := now.Sub(next) / every // 35s / 10s = 3
last := next.Add(missed * every) // start+40s
~~~
سياسة ~FireOnce~ بتشغّل مرة واحدة بس عن أول موعد فايت. سياسة ~Skip~ بتسيب المواعيد الفايتة كلها لو فات أكتر من موعد. في الحالتين الموعد الجاي ~last + every~، يعني ~start+50s~، فالمجدول ميعملش طوفان تشغيلات.

## منع التداخل
مهمة لسه شغالة متبتديش تاني فوق نفسها. ~Tick~ بتعلّمها ~running~، و ~Done(name)~ بتشيل العلامة. موعد وصل والمهمة شغالة بيعدّي من غير تشغيل.

## أخطاء شائعة
~~~go
sort.Slice(jobs, ...) // غلط: ترتيب كامل في كل Tick، ده O(n log n) بدل log n
h.Push(e)             // غلط: Push بتاعتك مبترتبش الشجرة؛ استخدم heap.Push(&h, e)
for next <= now { run(); next += every } // غلط: طوفان تشغيلات بعد التأخير
~~~

## المهمة
نفّذ ~NewScheduler(start, policy)~ و ~Add~ و ~Tick~ و ~Done~ بتوقيعات الـ starter، وسيب ~main~ فاضية. استخدم ~container/heap~ لترتيب المواعيد.

~Add(j)~ ترجع ~ErrJob~ لو ~Name~ فاضي، أو ~Every <= 0~، أو الاسم متضاف قبل كده. أول موعد ~start + Every~.

~Tick(now)~ بتعالج كل مهمة موعدها ~<= now~ وترجع التشغيلات (~Run~ بالاسم والموعد ~At~) مترتبة بالموعد، بعدين الأولوية الأعلى، بعدين الاسم. شريحة فاضية مش ~nil~ لو مفيش. مهمة شغالة متتشغلش. ~FireOnce~ بتشغّل مرة بـ ~At~ أول موعد فايت. ~Skip~ بتشغّل بس لو مفيش ولا موعد فات قبل الموعد ده. بعد المعالجة، شُغّلت أو لأ، الموعد الجاي هو أول موعد على الشبكة بعد ~now~.

~Done(name)~ بتنهي التشغيل الحالي؛ مهمة مجهولة أو مش شغالة ترجع ~ErrJob~. مفيش goroutines ولا ~time.Now~ ولا ~time.Sleep~.
`,
		Hints: []string{"عرّف نوع شريحة من مؤشرات لمدخلات فيها الـ Job والموعد الجاي وعلامة running، ونفّذ عليها heap.Interface.", "في Tick: طول ما h[0] موعده مش بعد now، اسحبه بـ heap.Pop، احسب missed و last، وقرر التشغيل حسب running والسياسة.", "حدّث next = last + Every ورجّع المدخل بـ heap.Push؛ وخلي Less تكسر التعادل بالأولوية بعدين الاسم."},
		En: &LevelText{Title: "Scheduler", Goal: "Run periodic maintenance jobs with a priority queue and a fake clock, without overlap and with a clear misfire policy", Files: []LabFile{{Name: "main.go", Starter: ch352Starter}},
			Lesson: `# The maintenance scheduler

## Before you write
- With a thousand jobs, how do you find the next one due without scanning them all every time?
- If the scheduler stops for an hour and resumes, do missed runs all execute, once, or not at all?

## Priority queues
A scheduler always needs "the earliest deadline". A heap is a tree arranged inside a slice: the smallest element is always at index ~0~, and inserting or removing costs ~log n~ instead of sorting the whole slice.

Package ~container/heap~ runs the algorithm on any type implementing ~heap.Interface~: ~Len~, ~Less~ and ~Swap~ as for sorting, plus ~Push(x any)~ and ~Pop() any~ that grow and shrink the slice at its end:
~~~go
type minInts []int

func (h minInts) Len() int           { return len(h) }
func (h minInts) Less(i, j int) bool { return h[i] < h[j] }
func (h minInts) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minInts) Push(x any)        { *h = append(*h, x.(int)) }
func (h *minInts) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}
~~~
Your ~Push~ and ~Pop~ are called only by the package. You call ~heap.Push(&h, 5)~ and ~heap.Pop(&h).(int)~, which keep the tree ordered. ~any~ is another name for ~interface{}~, so a type assertion follows ~Pop~. ~h[0]~ peeks at the smallest without removing it.

## A stable order
When two jobs share a deadline, ~Less~ must decide by a fixed rule or run order becomes arbitrary. Here: earlier deadline, then higher priority, then smaller name. ~a.Equal(b)~ and ~a.Before(b)~ compare ~time.Time~ values.

## An injected clock
The scheduler never reads the machine clock or sleeps. Its driver calls ~Tick(now)~ with whatever time it wants, and the test advances time itself. The same inputs always produce the same runs.

## The grid and misfires
A job every ~10*time.Second~ from ~start~ is due at ~start+10s~, ~start+20s~, and so on. If ~Tick~ is called late at ~start+45s~, the slots ~10s~, ~20s~, ~30s~ and ~40s~ were missed. Dividing by the period counts how many slots passed after the first:
~~~go
missed := now.Sub(next) / every // 35s / 10s = 3
last := next.Add(missed * every) // start+40s
~~~
Policy ~FireOnce~ runs once, for the first missed slot. Policy ~Skip~ drops the missed slots entirely when more than one passed. Either way the next deadline is ~last + every~, that is ~start+50s~, so the scheduler never floods.

## Preventing overlap
A job still running must not start again on top of itself. ~Tick~ marks it ~running~ and ~Done(name)~ clears the mark. A slot reached while the job runs passes without a run.

## Common mistakes
~~~go
sort.Slice(jobs, ...) // Wrong: a full sort every Tick, O(n log n) instead of log n
h.Push(e)             // Wrong: your Push does not order the tree; use heap.Push(&h, e)
for next <= now { run(); next += every } // Wrong: floods runs after a delay
~~~

## Your mission
Implement ~NewScheduler(start, policy)~, ~Add~, ~Tick~ and ~Done~ with the starter signatures; leave ~main~ empty. Use ~container/heap~ to order deadlines.

~Add(j)~ returns ~ErrJob~ for an empty ~Name~, ~Every <= 0~, or a name already added. The first deadline is ~start + Every~.

~Tick(now)~ handles every job due at ~<= now~ and returns runs (~Run~ with the name and slot ~At~) ordered by slot, then higher priority, then name; an empty non-nil slice when none. A running job does not run. ~FireOnce~ runs once with ~At~ the first missed slot. ~Skip~ runs only if no slot was missed before this one. After handling, run or not, the next deadline is the first grid slot after ~now~.

~Done(name)~ finishes the current run; an unknown or idle job returns ~ErrJob~. No goroutines, ~time.Now~ or ~time.Sleep~.
`, Hints: []string{"Define a slice of pointers to entries holding the Job, its next deadline, and a running flag, and implement heap.Interface on it.", "In Tick: while h[0] is not after now, take it with heap.Pop, compute missed and last, and decide by running and the policy.", "Set next = last + Every and put the entry back with heap.Push; let Less break ties by priority, then name."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch352Hidden},
	})
}

const ch352Starter = `package main

import (
	"errors"
	"time"
)

type Job struct {
	Name     string
	Every    time.Duration
	Priority int
}

type Misfire int

const (
	FireOnce Misfire = iota
	Skip
)

type Run struct {
	Name string
	At   time.Time
}

var ErrJob = errors.New("invalid job")

type Scheduler struct{}

func NewScheduler(start time.Time, policy Misfire) *Scheduler { return &Scheduler{} }
func (s *Scheduler) Add(j Job) error                          { return errors.New("TODO") }
func (s *Scheduler) Tick(now time.Time) []Run                 { return nil }
func (s *Scheduler) Done(name string) error                   { return errors.New("TODO") }
func main()                                                   {}
`

const ch352Hidden = `package main

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"
)

var hiddenStart = time.Date(2091, 3, 4, 5, 6, 7, 0, time.UTC)

type hiddenModelJob struct {
	job     Job
	next    time.Time
	running bool
}

// hiddenModel is a slow linear reference: no heap, same contract.
func hiddenModel(jobs []*hiddenModelJob, now time.Time, policy Misfire) []Run {
	type pending struct {
		run Run
		pri int
	}
	var out []pending
	for _, m := range jobs {
		if m.next.After(now) {
			continue
		}
		missed := now.Sub(m.next) / m.job.Every
		if !m.running && (policy == FireOnce || missed == 0) {
			out = append(out, pending{Run{m.job.Name, m.next}, m.job.Priority})
			m.running = true
		}
		m.next = m.next.Add((missed + 1) * m.job.Every)
	}
	sort.Slice(out, func(i, j int) bool {
		a, b := out[i], out[j]
		if !a.run.At.Equal(b.run.At) {
			return a.run.At.Before(b.run.At)
		}
		if a.pri != b.pri {
			return a.pri > b.pri
		}
		return a.run.Name < b.run.Name
	})
	runs := []Run{}
	for _, p := range out {
		runs = append(runs, p.run)
	}
	return runs
}

func TestHiddenSchedulerOrderOnTime(t *testing.T) {
	rng := rand.New(rand.NewSource(3502))
	s := NewScheduler(hiddenStart, FireOnce)
	var model []*hiddenModelJob
	for i := 0; i < 40; i++ {
		j := Job{Name: fmt.Sprintf("job-%02d", rng.Intn(60)), Every: time.Duration(1+rng.Intn(9)) * time.Second, Priority: rng.Intn(4)}
		err := s.Add(j)
		dup := false
		for _, m := range model {
			dup = dup || m.job.Name == j.Name
		}
		if dup {
			if !errors.Is(err, ErrJob) {
				t.Fatal("duplicate accepted", j.Name, err)
			}
			continue
		}
		if err != nil {
			t.Fatal("add", err)
		}
		model = append(model, &hiddenModelJob{job: j, next: hiddenStart.Add(j.Every)})
	}
	for sec := 0; sec <= 60; sec++ {
		now := hiddenStart.Add(time.Duration(sec) * time.Second)
		want := hiddenModel(model, now, FireOnce)
		got := s.Tick(now)
		if got == nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("second %d: got %v want %v", sec, got, want)
		}
		for _, r := range got {
			if err := s.Done(r.Name); err != nil {
				t.Fatal("done", err)
			}
		}
		for _, m := range model {
			m.running = false
		}
	}
}

func TestHiddenSchedulerOverlap(t *testing.T) {
	s := NewScheduler(hiddenStart, FireOnce)
	s.Add(Job{Name: "scrub", Every: 10 * time.Second})
	at := func(sec int) time.Time { return hiddenStart.Add(time.Duration(sec) * time.Second) }
	if got := s.Tick(at(9)); got == nil || len(got) != 0 {
		t.Fatal("early", got)
	}
	if got := s.Tick(at(10)); !reflect.DeepEqual(got, []Run{{"scrub", at(10)}}) {
		t.Fatal("boundary", got)
	}
	if got := s.Tick(at(20)); len(got) != 0 {
		t.Fatal("overlapping run", got)
	}
	if err := s.Done("scrub"); err != nil {
		t.Fatal(err)
	}
	if err := s.Done("scrub"); !errors.Is(err, ErrJob) {
		t.Fatal("done twice", err)
	}
	if got := s.Tick(at(25)); len(got) != 0 {
		t.Fatal("skipped slot replayed", got)
	}
	if got := s.Tick(at(30)); !reflect.DeepEqual(got, []Run{{"scrub", at(30)}}) {
		t.Fatal("next slot", got)
	}
	if err := s.Done("ghost"); !errors.Is(err, ErrJob) {
		t.Fatal("unknown done", err)
	}
	for _, bad := range []Job{{Name: "", Every: time.Second}, {Name: "zero"}, {Name: "neg", Every: -time.Second}, {Name: "scrub", Every: time.Second}} {
		if err := s.Add(bad); !errors.Is(err, ErrJob) {
			t.Fatal("validation", bad, err)
		}
	}
}

func TestHiddenSchedulerMisfirePolicies(t *testing.T) {
	at := func(sec int) time.Time { return hiddenStart.Add(time.Duration(sec) * time.Second) }
	once := NewScheduler(hiddenStart, FireOnce)
	skip := NewScheduler(hiddenStart, Skip)
	for _, s := range []*Scheduler{once, skip} {
		if err := s.Add(Job{Name: "vent", Every: 10 * time.Second}); err != nil {
			t.Fatal(err)
		}
	}
	if got := once.Tick(at(45)); !reflect.DeepEqual(got, []Run{{"vent", at(10)}}) {
		t.Fatal("fire once", got)
	}
	once.Done("vent")
	if got := skip.Tick(at(45)); len(got) != 0 {
		t.Fatal("skip fired a misfire", got)
	}
	for _, s := range []*Scheduler{once, skip} {
		if got := s.Tick(at(49)); len(got) != 0 {
			t.Fatal("next deadline must be after now", got)
		}
		if got := s.Tick(at(50)); !reflect.DeepEqual(got, []Run{{"vent", at(50)}}) {
			t.Fatal("resume on grid", got)
		}
		s.Done("vent")
	}
	if got := skip.Tick(at(69)); !reflect.DeepEqual(got, []Run{{"vent", at(60)}}) {
		t.Fatal("skip late within one period still runs", got)
	}
	if got := skip.Tick(at(5)); len(got) != 0 {
		t.Fatal("time going back", got)
	}
}

func TestHiddenSchedulerRandomJumps(t *testing.T) {
	for _, policy := range []Misfire{FireOnce, Skip} {
		rng := rand.New(rand.NewSource(int64(35020 + policy)))
		s := NewScheduler(hiddenStart, policy)
		var model []*hiddenModelJob
		for i := 0; i < 300; i++ {
			j := Job{Name: fmt.Sprintf("j%03d", i), Every: time.Duration(1+rng.Intn(5000)) * time.Millisecond, Priority: rng.Intn(3) - 1}
			if err := s.Add(j); err != nil {
				t.Fatal(err)
			}
			model = append(model, &hiddenModelJob{job: j, next: hiddenStart.Add(j.Every)})
		}
		now := hiddenStart
		for step := 0; step < 400; step++ {
			now = now.Add(time.Duration(rng.Intn(9000)) * time.Millisecond)
			want := hiddenModel(model, now, policy)
			got := s.Tick(now)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("policy %d step %d: got %d runs want %d", policy, step, len(got), len(want))
			}
			for _, m := range model {
				if m.running && rng.Intn(2) == 0 {
					m.running = false
					if err := s.Done(m.job.Name); err != nil {
						t.Fatal("done", err)
					}
				}
			}
		}
	}
}
`

func init() {
	add(&Level{ID: "35-3-redis-lock", Chapter: 35, Kind: "lab", Title: "القفل الموزّع", Goal: "امنع ورشتين يشتغلوا على نفس المفاعل في نفس الوقت بقفل Redis له مالك ومدة، ووضّح حدود أمانه", Setup: seededLabRoom,
		Requires: []string{"sel:hex.EncodeToString", "sel:io.ReadFull"}, Deps: []string{"github.com/redis/go-redis/v9", "github.com/alicebob/miniredis/v2"}, Bridge: []int{94}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch353Starter}},
		Lesson: `# قفل المفاعل

## قبل ما تكتب
- لو صاحب القفل وقع ومرجعش يفكّه، القفل يفضل مقفول للأبد؟
- لو القفل خلصت مدته واتاخد من حد تاني، إزاي تمنع صاحبه القديم إنه يفكّ قفل غيره؟

## ليه ~sync.Mutex~ مش كفاية
~sync.Mutex~ بيحمي goroutines جوه نفس العملية. ورش المستعمرة عمليات منفصلة على أجهزة مختلفة، فمحتاجين مكان مشترك يقول مين ماسك القفل. هنا المكان ده Redis.

## القفل بمدة
القفل مفتاح في Redis. اللي يقدر يكتبه لو مش موجود ياخد القفل. ~SetNX~ بتكتب بس لو المفتاح غايب، وترجع ~true~ لو كتبت:
~~~go
ok, err := client.SetNX(ctx, "colony:lock:pump", token, 5*time.Second).Result()
~~~
آخر وسيط مدة (TTL). لو صاحب القفل وقع، Redis نفسها بتمسح المفتاح بعد المدة، فالقفل ميفضلش مقفول للأبد.

## التوكن: مين المالك
قيمة المفتاح توكن عشوائي خاص بكل عملية أخد. اقرا بايتات من مصدر عشوائي محقون بـ ~io.ReadFull~، اللي بترجع خطأ لو ملقتش البايتات كلها، وحوّلها لنص بـ ~hex.EncodeToString~ من ~encoding/hex~:
~~~go
raw := make([]byte, 16)
if _, err := io.ReadFull(random, raw); err != nil { return err }
token := hex.EncodeToString(raw) // 32 حرف
~~~
في الإنتاج المصدر ~crypto/rand.Reader~؛ التست بيمرر مصدر بـ seed عشان التوكن يبقى متوقع.

## قارن وامسح في خطوة واحدة
الفك لازم يمسح المفتاح بس لو قيمته لسه التوكن بتاعك. ~Get~ وبعدها ~Del~ خطوتين: بينهم المدة ممكن تخلص وحد تاني ياخد القفل، والـ ~Del~ تمسح قفله. Redis بتنفذ سكربت Lua كامل من غير ما أمر تاني يدخل في النص:
~~~go
var unlock = redis.NewScript(` + "`" + `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0` + "`" + `)
n, err := unlock.Run(ctx, client, []string{"colony:lock:pump"}, token).Int()
~~~
~redis.NewScript~ بتجهّز السكربت مرة واحدة. ~Run~ بتاخد أسماء المفاتيح (~KEYS~ في Lua) وبعدها القيم (~ARGV~)، والترقيم في Lua بيبدأ من ~1~. ~redis.call~ بتنفذ أمر Redis جوه السكربت. ~Int()~ بترجع رقم النتيجة: ~1~ يعني اتمسح، و ~0~ يعني مش بتاعك.

## حدود الأمان
القفل ده استشاري. صاحب قفل وقف أكتر من المدة (GC، شبكة بطيئة) هيصحى فاكر إنه لسه ماسكه، مع إن Redis إدّته لحد تاني. التوكن بيمنعه يمسح قفل غيره، لكن ميمنعش كتاباته على المورد نفسه. الحماية الكاملة محتاجة fencing token: رقم بيزيد مع كل أخد، والمخزن نفسه بيرفض أي رقم أقدم من آخر رقم شافه.

## أخطاء شائعة
~~~go
client.Del(ctx, key) // غلط: بيمسح قفل أي حد، حتى لو مش بتاعك
if client.Get(ctx, key).Val() == token { client.Del(ctx, key) } // غلط: خطوتين، والمدة ممكن تخلص بينهم
client.SetNX(ctx, key, "locked", 0) // غلط: من غير مدة ومن غير مالك
~~~

## المهمة
نفّذ ~NewLocker(client, random)~ و ~Acquire~ و ~Release~ بتوقيعات الـ starter، وسيب ~main~ فاضية. المفتاح بالظبط ~colony:lock:~ زائد ~key~.

~Acquire(ctx, key, ttl)~ ترجع ~ErrLockKey~ قبل أي I/O لو ~key~ فاضي أو ~ttl <= 0~. بعدين تقرا ~16~ بايت بالظبط من ~random~ (خطأ القراية يترجع من غير ما تلمس Redis) والتوكن هو الـ hex بتاعهم. لو المفتاح موجود ترجع ~ErrHeld~، وإلا تكتبه بالتوكن والمدة وترجع ~Lease{Key: key, Token: token}~.

~Release(ctx, lease)~ تمسح المفتاح بسكربت واحد بس لو قيمته توكن الـ lease، وإلا ترجع ~ErrNotOwner~ من غير تغيير. ده يشمل قفل خلصت مدته، أو اتاخد من حد تاني، أو اتفك قبل كده.

اكتب تعليق في ~main.go~ يشرح حدود أمان القفل ويذكر كلمة ~fencing~؛ التست بيدور عليها في التعليقات.

النجمة التالتة: ~Extend(ctx, lease, ttl)~ بسكربت مقارنة بيجدد المدة لـ ~ttl~ (~PEXPIRE~ بالمللي ثانية) لو المالك لسه هو؛ غير كده ~ErrNotOwner~، و ~ttl <= 0~ ترجع ~ErrLockKey~.
`,
		Hints: []string{"ابدأ بالتحقق، وبعدين io.ReadFull لـ 16 بايت و hex.EncodeToString، وبعدين SetNX بالمدة؛ false يبقى ErrHeld.", "للفك اعمل redis.NewScript فيه GET ومقارنة بـ ARGV[1] و DEL، وشغله بـ Run(ctx, client, []string{key}, token).Int().", "نتيجة 0 من السكربت تبقى ErrNotOwner؛ ومتنساش تعليق حدود الأمان فيه كلمة fencing."},
		En: &LevelText{Title: "Distributed Lock", Goal: "Stop two workshops from working on the same reactor at once with a Redis lock that has an owner and a duration, and state its safety limits", Files: []LabFile{{Name: "main.go", Starter: ch353Starter}},
			Lesson: `# The reactor lock

## Before you write
- If the lock holder crashes and never unlocks, does the lock stay locked forever?
- If the lock expires and someone else takes it, how do you stop the old holder from unlocking the new one?

## Why ~sync.Mutex~ is not enough
~sync.Mutex~ protects goroutines inside one process. Colony workshops are separate processes on different machines, so they need a shared place that says who holds the lock. Here that place is Redis.

## A lock with a duration
The lock is a Redis key. Whoever manages to write it while absent takes the lock. ~SetNX~ writes only when the key is absent and returns ~true~ when it wrote:
~~~go
ok, err := client.SetNX(ctx, "colony:lock:pump", token, 5*time.Second).Result()
~~~
The last argument is a duration (TTL). If the holder crashes, Redis deletes the key itself after that time, so the lock never stays held forever.

## The token: who owns it
The key's value is a random token unique to each acquisition. Read bytes from an injected random source with ~io.ReadFull~, which fails unless it fills the whole buffer, and turn them into text with ~hex.EncodeToString~ from ~encoding/hex~:
~~~go
raw := make([]byte, 16)
if _, err := io.ReadFull(random, raw); err != nil { return err }
token := hex.EncodeToString(raw) // 32 characters
~~~
In production the source is ~crypto/rand.Reader~; the test passes a seeded source so the token is predictable.

## Compare and delete in one step
Unlocking must delete the key only if its value is still your token. ~Get~ then ~Del~ is two steps: between them the lock can expire, someone else takes it, and the ~Del~ removes their lock. Redis runs a whole Lua script with no other command in between:
~~~go
var unlock = redis.NewScript(` + "`" + `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0` + "`" + `)
n, err := unlock.Run(ctx, client, []string{"colony:lock:pump"}, token).Int()
~~~
~redis.NewScript~ prepares the script once. ~Run~ takes key names (~KEYS~ in Lua) then values (~ARGV~), and Lua indexes start at ~1~. ~redis.call~ runs a Redis command inside the script. ~Int()~ returns the numeric result: ~1~ means deleted, ~0~ means not yours.

## Safety limits
This lock is advisory. A holder paused longer than the TTL (GC, a slow network) wakes up believing it still holds the lock, even though Redis gave it to someone else. The token stops it from deleting another holder's lock, but not from writing to the resource itself. Full protection needs a fencing token: a number that grows with each acquisition, which the storage itself checks, rejecting anything older than the latest it has seen.

## Common mistakes
~~~go
client.Del(ctx, key) // Wrong: deletes anyone's lock, even when it is not yours
if client.Get(ctx, key).Val() == token { client.Del(ctx, key) } // Wrong: two steps; the lock can expire between them
client.SetNX(ctx, key, "locked", 0) // Wrong: no duration and no owner
~~~

## Your mission
Implement ~NewLocker(client, random)~, ~Acquire~ and ~Release~ with the starter signatures; leave ~main~ empty. The key is exactly ~colony:lock:~ plus ~key~.

~Acquire(ctx, key, ttl)~ returns ~ErrLockKey~ before any I/O when ~key~ is empty or ~ttl <= 0~. Then it reads exactly ~16~ bytes from ~random~ (a read error is returned without touching Redis) and the token is their hex. If the key exists it returns ~ErrHeld~; otherwise it writes the token with the duration and returns ~Lease{Key: key, Token: token}~.

~Release(ctx, lease)~ deletes the key with one script only when its value is the lease token; otherwise it returns ~ErrNotOwner~ without changes. That covers a lock that expired, was taken by someone else, or was already released.

Write a comment in ~main.go~ explaining the lock's safety limits that mentions the word ~fencing~; the test looks for it in comments.

Third star: ~Extend(ctx, lease, ttl)~ uses a compare script to reset the duration to ~ttl~ (~PEXPIRE~ in milliseconds) while the owner is unchanged; otherwise ~ErrNotOwner~, and ~ttl <= 0~ returns ~ErrLockKey~.
`, Hints: []string{"Validate first, then io.ReadFull 16 bytes and hex.EncodeToString, then SetNX with the duration; false means ErrHeld.", "To unlock, build a redis.NewScript with GET, a comparison to ARGV[1], and DEL, and run it with Run(ctx, client, []string{key}, token).Int().", "A script result of 0 is ErrNotOwner; and do not forget the safety-limits comment containing the word fencing."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch353Hidden},
	})
}

const ch353Starter = `package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"io"
	"time"
)

var ErrLockKey = errors.New("invalid lock request")
var ErrHeld = errors.New("lock held")
var ErrNotOwner = errors.New("not lock owner")

type Lease struct {
	Key   string
	Token string
}

type Locker struct{}

func NewLocker(client *redis.Client, random io.Reader) *Locker { return &Locker{} }
func (l *Locker) Acquire(ctx context.Context, key string, ttl time.Duration) (Lease, error) {
	return Lease{}, errors.New("TODO")
}
func (l *Locker) Release(ctx context.Context, lease Lease) error { return errors.New("TODO") }
func (l *Locker) Extend(ctx context.Context, lease Lease, ttl time.Duration) error {
	return errors.New("TODO")
}
func main() {}
`

const ch353Hidden = `package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"go/parser"
	"go/token"
	"io"
	"math/rand"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func hiddenLockServer(t *testing.T) *miniredis.Miniredis {
	t.Helper()
	srv := miniredis.NewMiniRedis()
	conn, err := srv.Dial()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	t.Cleanup(srv.Close)
	return srv
}
func hiddenLockClient(t *testing.T, srv *miniredis.Miniredis) *redis.Client {
	t.Helper()
	client := redis.NewClient(&redis.Options{Dialer: func(context.Context, string, string) (net.Conn, error) { return srv.Dial() }, MaxRetries: -1})
	t.Cleanup(func() { client.Close() })
	return client
}
func hiddenToken(seed int64, skip int) string {
	rng := rand.New(rand.NewSource(seed))
	raw := make([]byte, 16)
	for i := 0; i <= skip; i++ {
		io.ReadFull(rng, raw)
	}
	return hex.EncodeToString(raw)
}

type hiddenBrokenReader struct{}

func (hiddenBrokenReader) Read(p []byte) (int, error) { return 3, io.ErrUnexpectedEOF }

func TestHiddenLockAcquireRelease(t *testing.T) {
	srv := hiddenLockServer(t)
	ctx := context.Background()
	for seed := int64(3530); seed < 3540; seed++ {
		l := NewLocker(hiddenLockClient(t, srv), rand.New(rand.NewSource(seed)))
		key := fmt.Sprintf("reactor-%d", seed)
		lease, err := l.Acquire(ctx, key, 4*time.Second)
		want := hiddenToken(seed, 0)
		if err != nil || lease != (Lease{key, want}) {
			t.Fatal("acquire", lease, err, want)
		}
		if got, _ := srv.Get("colony:lock:" + key); got != want {
			t.Fatal("stored token", got)
		}
		if srv.TTL("colony:lock:"+key) != 4*time.Second {
			t.Fatal("ttl", srv.TTL("colony:lock:"+key))
		}
		if _, err := l.Acquire(ctx, key, time.Second); !errors.Is(err, ErrHeld) {
			t.Fatal("double acquire", err)
		}
		if err := l.Release(ctx, Lease{key, hiddenToken(seed, 5)}); !errors.Is(err, ErrNotOwner) || !srv.Exists("colony:lock:"+key) {
			t.Fatal("foreign release", err)
		}
		if err := l.Release(ctx, lease); err != nil || srv.Exists("colony:lock:"+key) {
			t.Fatal("release", err)
		}
		if err := l.Release(ctx, lease); !errors.Is(err, ErrNotOwner) {
			t.Fatal("double release", err)
		}
		again, err := l.Acquire(ctx, key, time.Second)
		if err != nil || again.Token != hiddenToken(seed, 2) {
			t.Fatal("reacquire after release with a fresh token", again, err)
		}
	}
}

func TestHiddenLockExpiryAndStaleOwner(t *testing.T) {
	srv := hiddenLockServer(t)
	ctx := context.Background()
	a := NewLocker(hiddenLockClient(t, srv), rand.New(rand.NewSource(1)))
	b := NewLocker(hiddenLockClient(t, srv), rand.New(rand.NewSource(2)))
	old, err := a.Acquire(ctx, "core", 3*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	srv.FastForward(3*time.Second - time.Millisecond)
	if _, err := b.Acquire(ctx, "core", time.Second); !errors.Is(err, ErrHeld) {
		t.Fatal("expired early", err)
	}
	srv.FastForward(time.Millisecond)
	fresh, err := b.Acquire(ctx, "core", 10*time.Second)
	if err != nil {
		t.Fatal("after expiry", err)
	}
	if err := a.Release(ctx, old); !errors.Is(err, ErrNotOwner) {
		t.Fatal("stale owner released the new lock", err)
	}
	if got, _ := srv.Get("colony:lock:core"); got != fresh.Token {
		t.Fatal("new lock damaged", got)
	}
	if err := b.Release(ctx, fresh); err != nil {
		t.Fatal(err)
	}
	if err := a.Release(ctx, Lease{"never", "x"}); !errors.Is(err, ErrNotOwner) {
		t.Fatal("absent lock", err)
	}
}

func TestHiddenLockValidationAndRandomFailure(t *testing.T) {
	ctx := context.Background()
	dead := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", MaxRetries: -1})
	defer dead.Close()
	l := NewLocker(dead, hiddenBrokenReader{})
	for _, c := range []struct {
		key string
		ttl time.Duration
	}{{"", time.Second}, {"a", 0}, {"a", -time.Second}} {
		if _, err := l.Acquire(ctx, c.key, c.ttl); !errors.Is(err, ErrLockKey) {
			t.Fatal("validation", c, err)
		}
	}
	srv := hiddenLockServer(t)
	broken := NewLocker(hiddenLockClient(t, srv), hiddenBrokenReader{})
	if _, err := broken.Acquire(ctx, "pump", time.Second); err == nil || srv.Exists("colony:lock:pump") {
		t.Fatal("short random read must fail before Redis", err)
	}
}

func TestHiddenLockConcurrentContenders(t *testing.T) {
	srv := hiddenLockServer(t)
	ctx := context.Background()
	for round := 0; round < 5; round++ {
		key := fmt.Sprintf("bay-%d", round)
		var wg sync.WaitGroup
		var mu sync.Mutex
		var winners []Lease
		for w := 0; w < 16; w++ {
			l := NewLocker(hiddenLockClient(t, srv), rand.New(rand.NewSource(int64(round*100+w))))
			wg.Add(1)
			go func() {
				defer wg.Done()
				lease, err := l.Acquire(ctx, key, time.Minute)
				if errors.Is(err, ErrHeld) {
					return
				}
				if err != nil {
					t.Error(err)
					return
				}
				mu.Lock()
				winners = append(winners, lease)
				mu.Unlock()
			}()
		}
		wg.Wait()
		if len(winners) != 1 {
			t.Fatal("winners", len(winners))
		}
	}
}

func TestHiddenLockSafetyComment(t *testing.T) {
	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	for _, g := range f.Comments {
		if strings.Contains(strings.ToLower(g.Text()), "fencing") {
			return
		}
	}
	t.Fatal("main.go needs a comment explaining the lock's safety limits and fencing tokens")
}

func TestHiddenBonusLockExtend(t *testing.T) {
	srv := hiddenLockServer(t)
	ctx := context.Background()
	a := NewLocker(hiddenLockClient(t, srv), rand.New(rand.NewSource(7)))
	b := NewLocker(hiddenLockClient(t, srv), rand.New(rand.NewSource(8)))
	lease, _ := a.Acquire(ctx, "arm", 2*time.Second)
	srv.FastForward(time.Second)
	if err := a.Extend(ctx, lease, 5*time.Second); err != nil || srv.TTL("colony:lock:arm") != 5*time.Second {
		t.Fatal("extend", err, srv.TTL("colony:lock:arm"))
	}
	if err := b.Extend(ctx, Lease{"arm", "intruder"}, time.Hour); !errors.Is(err, ErrNotOwner) || srv.TTL("colony:lock:arm") != 5*time.Second {
		t.Fatal("foreign extend", err)
	}
	if err := a.Extend(ctx, lease, 0); !errors.Is(err, ErrLockKey) {
		t.Fatal("zero ttl", err)
	}
	srv.FastForward(5 * time.Second)
	if err := a.Extend(ctx, lease, time.Second); !errors.Is(err, ErrNotOwner) {
		t.Fatal("extend after expiry", err)
	}
}
`

func init() {
	add(&Level{ID: "35-4-outbox", Chapter: 35, Kind: "lab", Title: "صندوق الصادر", Goal: "سجّل طلب قطع الغيار وحدثه في transaction واحدة، وانشر الأحداث بإعادة محاولة لمستهلك ميعدّش مرتين", Setup: seededLabRoom,
		Requires: []string{"mutex", "sel:json.Unmarshal"}, Deps: []string{"modernc.org/sqlite"}, Bridge: []int{95}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch354Starter}},
		Lesson: `# صندوق صادر المخزن

## قبل ما تكتب
- لو الطلب اتحفظ في القاعدة والسيرفر وقع قبل ما يبعت الحدث، باقي المستعمرة هتعرف إمتى؟
- لو الحدث وصل للمستهلك مرتين، إزاي متعدّش الكمية مرتين؟

## مشكلة الكتابتين
خدمة المخزن بتحفظ الطلب في قاعدتها وبتبلّغ الخدمات التانية بحدث. لو حفظت وبعدين بعتت، الوقوع بين الخطوتين يضيّع الحدث. لو بعتت وبعدين حفظت، الوقوع يبعت حدث عن طلب مش موجود. مفيش transaction واحدة بتغطي القاعدة والوسيط مع بعض.

## نمط الـ outbox
الحل إن الحدث نفسه يتكتب كصف في جدول ~outbox~ جوه نفس القاعدة، في نفس الـ transaction بتاعة الطلب. يا الاتنين يتحفظوا، يا ولا واحد. بعدين عملية منفصلة اسمها relay بتقرا الصفوف اللي لسه متنشرتش، تنشرها، وتعلّمها.
~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback()
tx.ExecContext(ctx, "INSERT INTO orders ...")
tx.ExecContext(ctx, "INSERT INTO outbox ...")
return tx.Commit()
~~~
~defer tx.Rollback()~ آمنة بعد ~Commit~: ساعتها مبتعملش حاجة. لو أي خطوة فشلت ورجعت بدري، الـ rollback بيلغي الاتنين.

## الـ relay والتسليم مرة على الأقل
الـ relay بينشر الحدث وبعدين يعلّمه ~published = 1~. لو وقع بين الاتنين، المرة الجاية هينشره تاني. لو الناشر رجّع خطأ، الحدث يفضل مش منشور ويتعاد بعدين. ده نفس at-least-once من طابور المهام: التكرار وارد.

عشان كده الترتيب مهم: اقرا الصفوف بـ ~ORDER BY id~، ولو نشر حدث فشل وقّف، عشان الحدث اللي بعده ميسبقوش.

## مستهلك idempotent
المستهلك بيفتكر أرقام الأحداث اللي عالجها في ~map[int64]bool~. حدث رقمه اتشاف قبل كده يتجاهل، فالتكرار ميغيرش النتيجة. المستهلك ممكن يتنادى من كذا goroutine، فاحمي الخريطة بـ ~sync.Mutex~.

## نسخة الشكل (schema version)
شكل الحدث بيتغير مع الوقت، والمستهلك ممكن يستلم أحداث قديمة متخزنة. كل حدث شايل ~Version~. النسخة ~1~ كان فيها ~id~ و ~item~ بس، ومعناها قطعة واحدة. النسخة ~2~ زوّدت ~qty~. ابدأ بقيمة افتراضية قبل الفك، والحقل الغايب في JSON يسيبها زي ما هي:
~~~go
p := orderPayload{Qty: 1}
err := json.Unmarshal([]byte(` + "`" + `{"id":"o1","item":"valve"}` + "`" + `), &p) // p.Qty == 1
~~~
نسخة أكبر من اللي المستهلك يعرفها ترجع خطأ بدل ما يخمّن معناها. الـ tags زي ~` + "`" + `json:"qty"` + "`" + `~ بتحدد اسم الحقل في JSON.

## أخطاء شائعة
~~~go
db.ExecContext(ctx, "INSERT INTO orders ...")
db.ExecContext(ctx, "INSERT INTO outbox ...") // غلط: كتابتين منفصلتين؛ الوقوع بينهم بيضيّع الحدث
db.ExecContext(ctx, "UPDATE outbox SET published = 1 ...")
pub.Publish(ctx, e) // غلط: علّمته قبل النشر؛ لو النشر فشل الحدث ضاع للأبد
~~~

## المهمة
نفّذ ~NewOutbox(db)~ و ~Init~ و ~PlaceOrder~ و ~Relay~، و ~NewConsumer~ و ~Handle~ و ~Totals~ بتوقيعات الـ starter، وسيب ~main~ فاضية. التست بيفتح القاعدة ويقفلها.

~Init~ بتنشئ الجداول لو مش موجودة بالشكل ده بالظبط: ~orders(id TEXT PRIMARY KEY, item TEXT, qty INTEGER)~ و ~outbox(id INTEGER PRIMARY KEY AUTOINCREMENT, type TEXT, version INTEGER, payload TEXT, published INTEGER DEFAULT 0)~.

~PlaceOrder(ctx, id, item, qty)~ ترجع ~ErrOrder~ قبل أي I/O لو ~id~ أو ~item~ فاضي أو ~qty <= 0~. في transaction واحدة بتضيف الطلب وصف outbox بـ ~type~ ~"order.placed"~ و ~version~ ~SchemaVersion~ و ~payload~ JSON بالحقول ~id~ و ~item~ و ~qty~. أي فشل (طلب مكرر، أو فشل كتابة الحدث) يرجع خطأ ومفيش ولا صف يتحفظ.

~Relay(ctx, pub, limit)~ بتنشر لحد ~limit~ حدث مش منشور بترتيب ~id~، وتعلّم كل واحد بعد نجاح نشره، وترجع عدد اللي اتعلّموا. أول فشل نشر بيوقف ويرجع العدد لحد هناك والخطأ. ~limit <= 0~ ترجع ~0~ و ~nil~.

~Handle(e)~ ترجع ~ErrVersion~ لو ~Version~ أقل من ~1~ أو أكبر من ~SchemaVersion~، وخطأ الفك لو JSON بايظ. حدث ~ID~ بتاعه اتعالج يرجع ~false~ و ~nil~ من غير تغيير، والجديد يزوّد إجمالي الصنف ويرجع ~true~. ~Totals~ ترجع نسخة من الإجماليات لكل ~item~. كله آمن بالتزامن.
`,
		Hints: []string{"في PlaceOrder: اتحقق، اعمل json.Marshal للـ payload، وبعدين BeginTx و defer Rollback واتنين INSERT و Commit.", "Relay: اقرا الصفوف كلها الأول وقفل rows، وبعدين لكل حدث Publish وبعد النجاح UPDATE published = 1.", "Consumer: افحص الـ Version، ابدأ payload بـ Qty: 1 قبل Unmarshal، وتحت القفل تجاهل أي ID اتشاف."},
		En: &LevelText{Title: "Outbox", Goal: "Record a spare-parts order and its event in one transaction, and relay events with retries to a consumer that never counts twice", Files: []LabFile{{Name: "main.go", Starter: ch354Starter}},
			Lesson: `# The depot outbox

## Before you write
- If the order is saved and the server crashes before sending the event, when does the rest of the colony find out?
- If the event reaches the consumer twice, how do you avoid counting the quantity twice?

## The dual-write problem
The depot service saves an order in its database and notifies other services with an event. Save then send, and a crash between the steps loses the event. Send then save, and a crash sends an event about an order that does not exist. No single transaction covers both the database and the broker.

## The outbox pattern
The fix is to write the event itself as a row in an ~outbox~ table inside the same database, in the same transaction as the order. Either both are saved or neither is. A separate process called the relay then reads unpublished rows, publishes them, and marks them.
~~~go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback()
tx.ExecContext(ctx, "INSERT INTO orders ...")
tx.ExecContext(ctx, "INSERT INTO outbox ...")
return tx.Commit()
~~~
~defer tx.Rollback()~ is safe after ~Commit~: it then does nothing. If any step fails and returns early, the rollback undoes both.

## The relay and at-least-once delivery
The relay publishes the event, then marks it ~published = 1~. If it crashes between the two, it publishes again next time. If the publisher returns an error, the event stays unpublished and is retried later. This is the same at-least-once delivery as the task queue: duplicates happen.

So order matters: read rows with ~ORDER BY id~, and stop when a publish fails so the next event does not overtake it.

## An idempotent consumer
The consumer remembers the event numbers it handled in a ~map[int64]bool~. An event number seen before is ignored, so duplicates do not change the result. The consumer may be called from several goroutines, so protect the map with a ~sync.Mutex~.

## Schema versions
Event shapes change over time, and a consumer may receive old stored events. Each event carries a ~Version~. Version ~1~ had only ~id~ and ~item~, meaning one unit. Version ~2~ added ~qty~. Start with a default before decoding; a field absent from the JSON leaves it unchanged:
~~~go
p := orderPayload{Qty: 1}
err := json.Unmarshal([]byte(` + "`" + `{"id":"o1","item":"valve"}` + "`" + `), &p) // p.Qty == 1
~~~
A version newer than the consumer knows returns an error instead of guessing its meaning. Tags like ~` + "`" + `json:"qty"` + "`" + `~ set the field's JSON name.

## Common mistakes
~~~go
db.ExecContext(ctx, "INSERT INTO orders ...")
db.ExecContext(ctx, "INSERT INTO outbox ...") // Wrong: two separate writes; a crash between them loses the event
db.ExecContext(ctx, "UPDATE outbox SET published = 1 ...")
pub.Publish(ctx, e) // Wrong: marked before publishing; a failed publish loses the event forever
~~~

## Your mission
Implement ~NewOutbox(db)~, ~Init~, ~PlaceOrder~ and ~Relay~, plus ~NewConsumer~, ~Handle~ and ~Totals~ with the starter signatures; leave ~main~ empty. The test opens and closes the database.

~Init~ creates the tables if absent with exactly this shape: ~orders(id TEXT PRIMARY KEY, item TEXT, qty INTEGER)~ and ~outbox(id INTEGER PRIMARY KEY AUTOINCREMENT, type TEXT, version INTEGER, payload TEXT, published INTEGER DEFAULT 0)~.

~PlaceOrder(ctx, id, item, qty)~ returns ~ErrOrder~ before any I/O when ~id~ or ~item~ is empty or ~qty <= 0~. In one transaction it inserts the order and an outbox row with ~type~ ~"order.placed"~, ~version~ ~SchemaVersion~, and a JSON ~payload~ with fields ~id~, ~item~ and ~qty~. Any failure (a duplicate order, or a failed event write) returns an error and saves no rows.

~Relay(ctx, pub, limit)~ publishes up to ~limit~ unpublished events in ~id~ order, marks each after its publish succeeds, and returns how many were marked. The first publish failure stops and returns the count so far and the error. ~limit <= 0~ returns ~0~ and ~nil~.

~Handle(e)~ returns ~ErrVersion~ when ~Version~ is below ~1~ or above ~SchemaVersion~, and the decode error for malformed JSON. An event whose ~ID~ was handled returns ~false~ and ~nil~ without changes; a new one adds to its item's total and returns ~true~. ~Totals~ returns a copy of the totals per ~item~. Everything is safe for concurrent use.
`, Hints: []string{"In PlaceOrder: validate, json.Marshal the payload, then BeginTx, defer Rollback, two INSERTs and Commit.", "Relay: read all rows first and close rows, then for each event Publish and, after success, UPDATE published = 1.", "Consumer: check Version, start the payload with Qty: 1 before Unmarshal, and under the lock ignore any seen ID."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch354Hidden},
	})
}

const ch354Starter = `package main

import (
	"context"
	"database/sql"
	"errors"
	_ "modernc.org/sqlite"
)

const SchemaVersion = 2

var ErrOrder = errors.New("invalid order")
var ErrVersion = errors.New("unknown schema version")

type Event struct {
	ID      int64
	Type    string
	Version int
	Payload string
}

type Publisher interface {
	Publish(ctx context.Context, e Event) error
}

type Outbox struct{}

func NewOutbox(db *sql.DB) *Outbox                     { return &Outbox{} }
func (o *Outbox) Init(ctx context.Context) error       { return errors.New("TODO") }
func (o *Outbox) PlaceOrder(ctx context.Context, id, item string, qty int) error {
	return errors.New("TODO")
}
func (o *Outbox) Relay(ctx context.Context, pub Publisher, limit int) (int, error) {
	return 0, errors.New("TODO")
}

type Consumer struct{}

func NewConsumer() *Consumer                    { return &Consumer{} }
func (c *Consumer) Handle(e Event) (bool, error) { return false, errors.New("TODO") }
func (c *Consumer) Totals() map[string]int       { return nil }
func main()                                      {}
`

const ch354Hidden = `package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

type hiddenPublisher struct {
	sent   []Event
	failAt map[int]bool // call number -> fail
	lost   map[int]bool // call number -> deliver, then report failure
	calls  int
}

func (p *hiddenPublisher) Publish(ctx context.Context, e Event) error {
	p.calls++
	if p.failAt[p.calls] {
		return errors.New("broker down")
	}
	p.sent = append(p.sent, e)
	if p.lost[p.calls] {
		return errors.New("ack lost")
	}
	return nil
}

func hiddenOutbox(t *testing.T) (*sql.DB, *Outbox) {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "depot.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	db.SetMaxOpenConns(1)
	o := NewOutbox(db)
	if err := o.Init(context.Background()); err != nil {
		t.Fatal("init", err)
	}
	if err := o.Init(context.Background()); err != nil {
		t.Fatal("init twice", err)
	}
	return db, o
}
func hiddenCount(t *testing.T, db *sql.DB, q string) int {
	t.Helper()
	var n int
	if err := db.QueryRow(q).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func TestHiddenOutboxAtomicWrite(t *testing.T) {
	db, o := hiddenOutbox(t)
	ctx := context.Background()
	rng := rand.New(rand.NewSource(3504))
	for i := 0; i < 12; i++ {
		id, item, qty := fmt.Sprintf("o-%d", i), fmt.Sprintf("valve-%d", rng.Intn(4)), 1+rng.Intn(20)
		if err := o.PlaceOrder(ctx, id, item, qty); err != nil {
			t.Fatal(err)
		}
		var typ, payload string
		var version, published int
		if err := db.QueryRow("SELECT type, version, payload, published FROM outbox ORDER BY id DESC LIMIT 1").Scan(&typ, &version, &payload, &published); err != nil {
			t.Fatal(err)
		}
		var p map[string]any
		if json.Unmarshal([]byte(payload), &p) != nil || p["id"] != id || p["item"] != item || p["qty"] != float64(qty) {
			t.Fatal("payload", payload)
		}
		if typ != "order.placed" || version != SchemaVersion || published != 0 {
			t.Fatal("event row", typ, version, published)
		}
	}
	if err := o.PlaceOrder(ctx, "o-3", "valve-0", 1); err == nil {
		t.Fatal("duplicate order accepted")
	}
	if _, err := db.Exec("CREATE TRIGGER poison BEFORE INSERT ON outbox WHEN NEW.payload LIKE '%poison%' BEGIN SELECT RAISE(ABORT, 'event rejected'); END"); err != nil {
		t.Fatal(err)
	}
	if err := o.PlaceOrder(ctx, "o-99", "poison", 2); err == nil {
		t.Fatal("failed event write reported success")
	}
	if n := hiddenCount(t, db, "SELECT COUNT(*) FROM orders WHERE id = 'o-99'"); n != 0 {
		t.Fatal("order kept without its event")
	}
	if hiddenCount(t, db, "SELECT COUNT(*) FROM orders") != 12 || hiddenCount(t, db, "SELECT COUNT(*) FROM outbox") != 12 {
		t.Fatal("row counts")
	}
	for _, bad := range []struct {
		id, item string
		qty      int
	}{{"", "a", 1}, {"x", "", 1}, {"x", "a", 0}, {"x", "a", -3}} {
		if err := o.PlaceOrder(ctx, bad.id, bad.item, bad.qty); !errors.Is(err, ErrOrder) {
			t.Fatal("validation", bad, err)
		}
	}
}

func TestHiddenOutboxRelayRetries(t *testing.T) {
	db, o := hiddenOutbox(t)
	ctx := context.Background()
	for i := 1; i <= 7; i++ {
		o.PlaceOrder(ctx, fmt.Sprintf("r-%d", i), "seal", i)
	}
	if n, err := o.Relay(ctx, &hiddenPublisher{}, 0); n != 0 || err != nil {
		t.Fatal("zero limit", n, err)
	}
	pub := &hiddenPublisher{failAt: map[int]bool{3: true}, lost: map[int]bool{6: true}}
	n, err := o.Relay(ctx, pub, 10)
	if n != 2 || err == nil || len(pub.sent) != 2 {
		t.Fatal("stop at first failure", n, err, len(pub.sent))
	}
	n, err = o.Relay(ctx, pub, 2)
	if n != 2 || err != nil {
		t.Fatal("limit", n, err)
	}
	n, err = o.Relay(ctx, pub, 10)
	if n != 0 || err == nil {
		t.Fatal("lost ack must not mark", n, err)
	}
	n, err = o.Relay(ctx, pub, 10)
	if n != 3 || err != nil {
		t.Fatal("final relay", n, err)
	}
	if hiddenCount(t, db, "SELECT COUNT(*) FROM outbox WHERE published = 0") != 0 {
		t.Fatal("unpublished left")
	}
	var ids []int64
	for _, e := range pub.sent {
		ids = append(ids, e.ID)
	}
	if !reflect.DeepEqual(ids, []int64{1, 2, 3, 4, 5, 5, 6, 7}) {
		t.Fatal("delivery order with one duplicate", ids)
	}
	c := NewConsumer()
	for _, e := range pub.sent {
		if _, err := c.Handle(e); err != nil {
			t.Fatal(err)
		}
	}
	if got := c.Totals(); !reflect.DeepEqual(got, map[string]int{"seal": 28}) {
		t.Fatal("duplicate counted", got)
	}
}

func TestHiddenConsumerVersionsAndDuplicates(t *testing.T) {
	c := NewConsumer()
	if ok, err := c.Handle(Event{ID: 1, Type: "order.placed", Version: 1, Payload: ` + "`" + `{"id":"a","item":"filter"}` + "`" + `}); !ok || err != nil {
		t.Fatal("v1", ok, err)
	}
	if ok, err := c.Handle(Event{ID: 2, Type: "order.placed", Version: 2, Payload: ` + "`" + `{"id":"b","item":"filter","qty":4}` + "`" + `}); !ok || err != nil {
		t.Fatal("v2", ok, err)
	}
	if ok, err := c.Handle(Event{ID: 2, Version: 2, Payload: ` + "`" + `{"id":"b","item":"filter","qty":4}` + "`" + `}); ok || err != nil {
		t.Fatal("duplicate", ok, err)
	}
	for _, v := range []int{0, 3, -1} {
		if _, err := c.Handle(Event{ID: 9, Version: v, Payload: "{}"}); !errors.Is(err, ErrVersion) {
			t.Fatal("version", v, err)
		}
	}
	if _, err := c.Handle(Event{ID: 10, Version: 2, Payload: "{broken"}); err == nil {
		t.Fatal("malformed payload")
	}
	got := c.Totals()
	if !reflect.DeepEqual(got, map[string]int{"filter": 5}) {
		t.Fatal("totals", got)
	}
	got["filter"] = 99
	if c.Totals()["filter"] != 5 {
		t.Fatal("Totals must return a copy")
	}
	if ok, _ := c.Handle(Event{ID: 9, Version: 2, Payload: ` + "`" + `{"item":"pipe","qty":1}` + "`" + `}); !ok {
		t.Fatal("rejected events must not be marked seen")
	}
}

func TestHiddenConsumerConcurrentDuplicates(t *testing.T) {
	c := NewConsumer()
	var wg sync.WaitGroup
	var mu sync.Mutex
	accepted := 0
	for w := 0; w < 12; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := int64(1); id <= 40; id++ {
				ok, err := c.Handle(Event{ID: id, Version: 2, Payload: fmt.Sprintf(` + "`" + `{"item":"bolt","qty":%d}` + "`" + `, id)})
				if err != nil {
					t.Error(err)
				}
				if ok {
					mu.Lock()
					accepted++
					mu.Unlock()
				}
			}
		}()
	}
	wg.Wait()
	if accepted != 40 || c.Totals()["bolt"] != 820 {
		t.Fatal("concurrent dedupe", accepted, c.Totals())
	}
}
`

func init() {
	add(&Level{ID: "35-5-kafka", Chapter: 35, Kind: "lab", Title: "الوسيط", Goal: "ابعت قراءات العربات لـ Kafka بمفتاح لكل عربة، واستهلكها بإعادة محاولة و offsets تتأكد بعد النجاح بس وإغلاق مرتب", Setup: seededLabRoom,
		Requires: []string{"sel:context.WithoutCancel", "sel:strconv.Atoi"}, Deps: []string{"github.com/segmentio/kafka-go"}, Bridge: []int{98}, Race: true, Service: "kafka",
		Files: []LabFile{{Name: "main.go", Starter: ch355Starter}},
		Lesson: `# وسيط قراءات العربات

## قبل ما تكتب
- لو عربة بعتت قراءتين ورا بعض، إزاي تضمن إن المستهلك يشوفهم بنفس الترتيب؟
- لو المستهلك وقع في نص معالجة رسالة، يبدأ منين لما يرجع؟

## Topic و partition
Kafka سجل رسائل بيتكتب في آخره بس. الـ topic اسم سجل، زي ~rover-telemetry~. الـ topic متقسم partitions، وكل partition سجل مرتب لوحده. الترتيب مضمون جوه الـ partition الواحد بس، مش على الـ topic كله.

## المفتاح
كل رسالة ليها ~Key~ و ~Value~ كبايتات. الـ writer بيختار الـ partition من الـ hash بتاع المفتاح، فكل الرسايل اللي ليها نفس المفتاح بتروح لنفس الـ partition بالترتيب. خلّي المفتاح اسم العربة:
~~~go
msg := kafka.Message{Key: []byte("rover-7"), Value: []byte("42")}
~~~
~kafka.Message~ من ~github.com/segmentio/kafka-go~. ~strconv.Itoa(42)~ بترجع ~"42"~، و ~strconv.Atoi("42")~ بترجع الرقم وخطأ لو النص مش رقم.

## الكتابة وإعادة المحاولة
~(*kafka.Writer).WriteMessages(ctx, msgs...)~ بتبعت دفعة رسايل. الشبكة ممكن تفشل مؤقتاً، فجرّب كام مرة محدودة. افحص ~ctx.Err()~ قبل كل محاولة: لو المنادي لغى، ملوش لازمة تكمّل.
~~~go
w := &kafka.Writer{Addr: kafka.TCP("127.0.0.1:9092"), Topic: "rover-telemetry", Balancer: &kafka.Hash{}}
defer w.Close()
err := w.WriteMessages(ctx, msg)
~~~
~kafka.Hash~ هو الـ balancer اللي بيوزّع بالمفتاح. الكود بتاعك بياخد ~Writer~ كـ interface فيها ~WriteMessages~ بس، فالتست يقدر يمرر fake.

## الـ offset ومجموعة المستهلكين
كل رسالة ليها ~Offset~: رقمها جوه الـ partition. مجموعة المستهلكين (consumer group) ليها اسم، و Kafka بتفتكر لكل partition آخر offset المجموعة أكدته (commit). المستهلك اللي يرجع بنفس المجموعة يكمّل من بعد آخر commit.
~~~go
r := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{"127.0.0.1:9092"}, Topic: "rover-telemetry", GroupID: "dashboard"})
defer r.Close()
msg, err := r.FetchMessage(ctx)
err = r.CommitMessages(ctx, msg)
~~~
~FetchMessage~ بتجيب الرسالة الجاية من غير commit. ~CommitMessages~ بتأكدها. لو أكدت قبل المعالجة ووقعت، الرسالة ضاعت. لو أكدت بعد المعالجة ووقعت قبل الـ commit، الرسالة هترجع تاني: ده at-least-once من الدروس اللي فاتت.

## رسالة مسمومة
رسالة قيمتها مش رقم هتفشل كل مرة. لو موقّفة السجل، كل اللي وراها هيستنى للأبد. أكدها وعدّيها من غير ما تنادي المعالج.

## إغلاق مرتب
لما المنادي يلغي ~ctx~، ~FetchMessage~ بترجع خطأ؛ ده مش عطل، ده طلب وقوف. لكن لو المعالج خلص رسالة، الـ commit بتاعها لازم يحصل حتى لو السياق اتلغى. ~context.WithoutCancel(ctx)~ بترجع سياق بنفس القيم بس مبيتلغيش:
~~~go
commitCtx := context.WithoutCancel(ctx) // ctx اتلغى؟ commitCtx لأ
~~~

## أخطاء شائعة
~~~go
r.ReadMessage(ctx) // غلط هنا: بتعمل commit تلقائي قبل ما تعالج
w.WriteMessages(ctx, kafka.Message{Value: v}) // غلط: من غير مفتاح، قراءات نفس العربة ممكن تتلخبط
if err := r.FetchMessage(ctx); err != nil { return 0, err } // غلط: الإلغاء المقصود اتعامل كعطل
~~~

## المهمة
نفّذ ~Send~ و ~Consume~ بتوقيعات الـ starter، وسيب ~main~ فاضية. التستات بتمرر fakes، وتستات التكامل بتمرر ~*kafka.Writer~ و ~*kafka.Reader~ حقيقيين لو ~ORBITAL_KAFKA_BROKERS~ متضبوطة.

~Send(ctx, w, batch)~ ترجع ~ErrTelemetry~ من غير كتابة لو أي ~Rover~ فاضي. دفعة فاضية ترجع ~nil~ من غير كتابة. كل قراءة رسالة مفتاحها ~Rover~ وقيمتها ~strconv.Itoa(Reading)~ بنفس الترتيب، والدفعة كلها في نداء ~WriteMessages~ واحد. لحد ~3~ محاولات؛ قبل كل محاولة لو ~ctx~ اتلغى ارجع ~ctx.Err()~. بعد فشل التلاتة ارجع خطأ بيلف آخر خطأ (~errors.Is~ يلاقيه).

~Consume(ctx, r, h, maxAttempts)~ بتلف: تجيب رسالة، تحوّلها لـ ~Telemetry~ (المفتاح هو ~Rover~)، وتنادي ~h~ لحد ~maxAttempts~ مرة لحد ما ينجح. بعد النجاح بس تعمل commit للرسالة دي بسياق مبيتلغيش وتزوّد العداد. قيمة مش رقم: commit من غير ما تنادي ~h~ ومن غير عد. لو المعالج فشل كل المحاولات ارجع العداد وخطأه من غير commit. إلغاء ~ctx~ (في الجلب أو أثناء المعالجة) يرجع العداد و ~nil~، ورسالة اتلغت معالجتها متتأكدش. خطأ جلب تاني يترجع زي ما هو.
`,
		Hints: []string{"Send: ابني []kafka.Message الأول (المفتاح []byte(Rover))، وبعدين لف 3 مرات وافحص ctx.Err() قبل كل WriteMessages.", "Consume: FetchMessage؛ لو فيه خطأ و ctx.Err() مش nil ارجع العداد و nil. strconv.Atoi على القيمة، ولو فشلت اعمل commit وكمّل.", "commit بعد نجاح h بس، بـ context.WithoutCancel(ctx)؛ فشل كل المحاولات يرجع الخطأ من غير commit إلا لو السياق اتلغى."},
		En: &LevelText{Title: "Broker", Goal: "Send rover readings to Kafka keyed by rover, and consume them with retries, offsets committed only after success, and a graceful close", Files: []LabFile{{Name: "main.go", Starter: ch355Starter}},
			Lesson: `# The rover telemetry broker

## Before you write
- If a rover sends two readings in a row, how do you guarantee the consumer sees them in the same order?
- If the consumer crashes in the middle of a message, where does it resume?

## Topics and partitions
Kafka is a message log that is only appended to. A topic names a log, like ~rover-telemetry~. A topic is split into partitions, each its own ordered log. Order is guaranteed only within one partition, not across the whole topic.

## Keys
Each message has a ~Key~ and a ~Value~ as bytes. The writer chooses the partition from the key's hash, so all messages with the same key go to the same partition, in order. Make the key the rover name:
~~~go
msg := kafka.Message{Key: []byte("rover-7"), Value: []byte("42")}
~~~
~kafka.Message~ comes from ~github.com/segmentio/kafka-go~. ~strconv.Itoa(42)~ returns ~"42"~, and ~strconv.Atoi("42")~ returns the number and an error when the text is not a number.

## Writing and retrying
~(*kafka.Writer).WriteMessages(ctx, msgs...)~ sends a batch. The network can fail briefly, so try a bounded number of times. Check ~ctx.Err()~ before each attempt: once the caller cancels, there is no point continuing.
~~~go
w := &kafka.Writer{Addr: kafka.TCP("127.0.0.1:9092"), Topic: "rover-telemetry", Balancer: &kafka.Hash{}}
defer w.Close()
err := w.WriteMessages(ctx, msg)
~~~
~kafka.Hash~ is the balancer that partitions by key. Your code takes ~Writer~ as an interface with only ~WriteMessages~, so the test can pass a fake.

## Offsets and consumer groups
Each message has an ~Offset~: its number inside the partition. A consumer group has a name, and Kafka remembers, per partition, the last offset the group committed. A consumer returning with the same group resumes after the last commit.
~~~go
r := kafka.NewReader(kafka.ReaderConfig{Brokers: []string{"127.0.0.1:9092"}, Topic: "rover-telemetry", GroupID: "dashboard"})
defer r.Close()
msg, err := r.FetchMessage(ctx)
err = r.CommitMessages(ctx, msg)
~~~
~FetchMessage~ gets the next message without committing. ~CommitMessages~ confirms it. Commit before processing and crash, and the message is lost. Commit after processing and crash before the commit, and the message comes again: at-least-once delivery from the earlier lessons.

## Poison messages
A message whose value is not a number fails every time. If it blocks the log, everything behind it waits forever. Commit it and move on without calling the handler.

## Graceful close
When the caller cancels ~ctx~, ~FetchMessage~ returns an error; that is a stop request, not a failure. But if the handler finished a message, its commit must happen even though the context was canceled. ~context.WithoutCancel(ctx)~ returns a context with the same values that is never canceled:
~~~go
commitCtx := context.WithoutCancel(ctx) // ctx canceled? commitCtx is not
~~~

## Common mistakes
~~~go
r.ReadMessage(ctx) // Wrong here: commits automatically before you process
w.WriteMessages(ctx, kafka.Message{Value: v}) // Wrong: without a key, one rover's readings can reorder
if err := r.FetchMessage(ctx); err != nil { return 0, err } // Wrong: a requested stop treated as a failure
~~~

## Your mission
Implement ~Send~ and ~Consume~ with the starter signatures; leave ~main~ empty. Tests pass fakes, and integration tests pass a real ~*kafka.Writer~ and ~*kafka.Reader~ when ~ORBITAL_KAFKA_BROKERS~ is set.

~Send(ctx, w, batch)~ returns ~ErrTelemetry~ without writing when any ~Rover~ is empty. An empty batch returns ~nil~ without writing. Each reading is a message keyed by ~Rover~ with value ~strconv.Itoa(Reading)~, in order, the whole batch in one ~WriteMessages~ call. Up to ~3~ attempts; before each, if ~ctx~ is canceled return ~ctx.Err()~. After three failures return an error wrapping the last one (found by ~errors.Is~).

~Consume(ctx, r, h, maxAttempts)~ loops: fetch a message, turn it into ~Telemetry~ (the key is ~Rover~), and call ~h~ up to ~maxAttempts~ times until it succeeds. Only after success commit that message with a context that is never canceled and increment the count. A non-numeric value: commit without calling ~h~ and without counting. If the handler fails every attempt, return the count and its error without committing. Canceling ~ctx~ (while fetching or handling) returns the count and ~nil~, and a message whose handling was canceled is not committed. Any other fetch error is returned as is.
`, Hints: []string{"Send: build []kafka.Message first (key []byte(Rover)), then loop 3 times checking ctx.Err() before each WriteMessages.", "Consume: FetchMessage; on error with ctx.Err() non-nil return the count and nil. strconv.Atoi the value, and on failure commit and continue.", "Commit only after h succeeds, with context.WithoutCancel(ctx); failing every attempt returns the error without committing unless the context was canceled."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch355Hidden, HiddenFiles: map[string]string{"kafka_integration_test.go": ch355Integration}},
	})
}

const ch355Starter = `package main

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
)

type Telemetry struct {
	Rover   string
	Reading int
}

var ErrTelemetry = errors.New("invalid telemetry")

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Reader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Handler func(ctx context.Context, t Telemetry) error

func Send(ctx context.Context, w Writer, batch []Telemetry) error { return errors.New("TODO") }
func Consume(ctx context.Context, r Reader, h Handler, maxAttempts int) (int, error) {
	return 0, errors.New("TODO")
}
func main() {}
`

const ch355Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"math/rand"
	"reflect"
	"strconv"
	"sync"
	"testing"
)

type hiddenWriter struct {
	calls [][]kafka.Message
	fails int
	err   error
}

func (w *hiddenWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	w.calls = append(w.calls, append([]kafka.Message(nil), msgs...))
	if w.fails != 0 {
		w.fails--
		return w.err
	}
	return nil
}

// hiddenLog is an in-memory partition: a reader resumes after the last commit.
type hiddenLog struct {
	mu        sync.Mutex
	msgs      []kafka.Message
	committed int64
}

func (l *hiddenLog) append(key, value string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.msgs = append(l.msgs, kafka.Message{Key: []byte(key), Value: []byte(value), Offset: int64(len(l.msgs))})
}

type hiddenReader struct {
	log      *hiddenLog
	pos      int64
	fetchErr error
}

func (l *hiddenLog) reader() *hiddenReader { return &hiddenReader{log: l, pos: l.committed} }
func (r *hiddenReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	if r.fetchErr != nil {
		return kafka.Message{}, r.fetchErr
	}
	r.log.mu.Lock()
	if r.pos < int64(len(r.log.msgs)) {
		m := r.log.msgs[r.pos]
		r.pos++
		r.log.mu.Unlock()
		return m, nil
	}
	r.log.mu.Unlock()
	<-ctx.Done()
	return kafka.Message{}, ctx.Err()
}
func (r *hiddenReader) CommitMessages(ctx context.Context, msgs ...kafka.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.log.mu.Lock()
	defer r.log.mu.Unlock()
	for _, m := range msgs {
		if m.Offset+1 > r.log.committed {
			r.log.committed = m.Offset + 1
		}
	}
	return nil
}

func TestHiddenKafkaSendKeysAndRetries(t *testing.T) {
	ctx := context.Background()
	rng := rand.New(rand.NewSource(3505))
	var batch []Telemetry
	for i := 0; i < 25; i++ {
		batch = append(batch, Telemetry{fmt.Sprintf("rover-%d", rng.Intn(4)), rng.Intn(2000) - 1000})
	}
	w := &hiddenWriter{}
	if err := Send(ctx, w, batch); err != nil || len(w.calls) != 1 || len(w.calls[0]) != len(batch) {
		t.Fatal("one batch", err, len(w.calls))
	}
	for i, m := range w.calls[0] {
		if string(m.Key) != batch[i].Rover || string(m.Value) != strconv.Itoa(batch[i].Reading) {
			t.Fatal("message", i, string(m.Key), string(m.Value))
		}
	}
	flaky := &hiddenWriter{fails: 2, err: errors.New("leader not available")}
	if err := Send(ctx, flaky, batch[:3]); err != nil || len(flaky.calls) != 3 {
		t.Fatal("retry until success", err, len(flaky.calls))
	}
	sentinel := errors.New("broker down")
	down := &hiddenWriter{fails: -1, err: sentinel}
	if err := Send(ctx, down, batch[:3]); !errors.Is(err, sentinel) || len(down.calls) != 3 {
		t.Fatal("bounded retries", err, len(down.calls))
	}
	empty := &hiddenWriter{}
	if err := Send(ctx, empty, []Telemetry{{"a", 1}, {"", 2}}); !errors.Is(err, ErrTelemetry) || len(empty.calls) != 0 {
		t.Fatal("validation", err)
	}
	if err := Send(ctx, empty, nil); err != nil || len(empty.calls) != 0 {
		t.Fatal("empty batch", err)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if err := Send(canceled, empty, batch); !errors.Is(err, context.Canceled) || len(empty.calls) != 0 {
		t.Fatal("canceled send", err)
	}
}

func TestHiddenKafkaConsumeOrderPoisonAndShutdown(t *testing.T) {
	log := &hiddenLog{}
	rng := rand.New(rand.NewSource(35050))
	var want []Telemetry
	for i := 0; i < 30; i++ {
		if i == 11 {
			log.append("rover-x", "not-a-number")
			continue
		}
		tm := Telemetry{fmt.Sprintf("rover-%d", rng.Intn(3)), rng.Intn(500)}
		want = append(want, tm)
		log.append(tm.Rover, strconv.Itoa(tm.Reading))
	}
	ctx, cancel := context.WithCancel(context.Background())
	var got []Telemetry
	n, err := Consume(ctx, log.reader(), func(_ context.Context, tm Telemetry) error {
		got = append(got, tm)
		if len(got) == len(want) {
			cancel() // shut down right after the last message: its commit must still happen
		}
		return nil
	}, 3)
	if err != nil || n != len(want) || !reflect.DeepEqual(got, want) {
		t.Fatal("consume", n, err)
	}
	if log.committed != 30 {
		t.Fatal("commits incl. poison and last message", log.committed)
	}
}

func TestHiddenKafkaHandlerRetriesAndRedelivery(t *testing.T) {
	log := &hiddenLog{}
	for i := 0; i < 6; i++ {
		log.append("rover-1", strconv.Itoa(i))
	}
	calls := map[int]int{}
	sentinel := errors.New("dashboard offline")
	h := func(_ context.Context, tm Telemetry) error {
		calls[tm.Reading]++
		if tm.Reading == 1 && calls[1] < 3 {
			return errors.New("transient")
		}
		if tm.Reading == 3 {
			return sentinel
		}
		return nil
	}
	ctx := context.Background()
	n, err := Consume(ctx, log.reader(), h, 3)
	if n != 3 || !errors.Is(err, sentinel) || calls[1] != 3 || calls[3] != 3 {
		t.Fatal("retries", n, err, calls)
	}
	if log.committed != 3 {
		t.Fatal("failed message must stay uncommitted", log.committed)
	}
	ctx2, cancel := context.WithCancel(ctx)
	var first []int
	n, err = Consume(ctx2, log.reader(), func(_ context.Context, tm Telemetry) error {
		first = append(first, tm.Reading)
		if len(first) == 3 {
			cancel()
		}
		return nil
	}, 1)
	if err != nil || n != 3 || !reflect.DeepEqual(first, []int{3, 4, 5}) {
		t.Fatal("restart resumes at the failed offset", n, err, first)
	}
	fetchErr := errors.New("rebalance in progress")
	if _, err := Consume(ctx, &hiddenReader{log: log, fetchErr: fetchErr}, h, 1); !errors.Is(err, fetchErr) {
		t.Fatal("fetch error", err)
	}
}

func TestHiddenKafkaCanceledHandlerNotCommitted(t *testing.T) {
	log := &hiddenLog{}
	log.append("rover-2", "7")
	log.append("rover-2", "8")
	ctx, cancel := context.WithCancel(context.Background())
	n, err := Consume(ctx, log.reader(), func(c context.Context, tm Telemetry) error {
		if tm.Reading == 8 {
			cancel()
			return c.Err()
		}
		return nil
	}, 5)
	if err != nil || n != 1 || log.committed != 1 {
		t.Fatal("canceled handling", n, err, log.committed)
	}
}
`

const ch355Integration = `package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func hiddenKafkaTopic(t *testing.T) (string, string) {
	t.Helper()
	brokers := os.Getenv("ORBITAL_KAFKA_BROKERS")
	if brokers == "" {
		t.Skip("set ORBITAL_KAFKA_BROKERS for Kafka integration")
	}
	broker := strings.Split(brokers, ",")[0]
	topic := fmt.Sprintf("orbital-35-5-%d-%d", time.Now().UnixNano(), rand.Int63())
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		t.Fatal(err)
	}
	cc, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		t.Fatal(err)
	}
	defer cc.Close()
	if err := cc.CreateTopics(kafka.TopicConfig{Topic: topic, NumPartitions: 3, ReplicationFactor: 1}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if c, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))); err == nil {
			c.DeleteTopics(topic)
			c.Close()
		}
	})
	return broker, topic
}

func TestHiddenKafkaIntegrationKeyOrderAndOffsets(t *testing.T) {
	broker, topic := hiddenKafkaTopic(t)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	w := &kafka.Writer{Addr: kafka.TCP(broker), Topic: topic, Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll}
	defer w.Close()
	var batch []Telemetry
	for i := 0; i < 30; i++ {
		batch = append(batch, Telemetry{fmt.Sprintf("rover-%d", i%3), i})
	}
	if err := Send(ctx, w, batch); err != nil {
		t.Fatal("send", err)
	}
	group := topic + "-dashboard"
	newReader := func() *kafka.Reader {
		return kafka.NewReader(kafka.ReaderConfig{Brokers: []string{broker}, Topic: topic, GroupID: group, StartOffset: kafka.FirstOffset, MaxWait: 200 * time.Millisecond})
	}
	last := map[string]int{}
	seen := 0
	stopAt := 12
	r := newReader()
	runCtx, stop := context.WithCancel(ctx)
	n, err := Consume(runCtx, r, func(_ context.Context, tm Telemetry) error {
		if prev, ok := last[tm.Rover]; ok && tm.Reading <= prev {
			t.Errorf("rover %s out of order: %d after %d", tm.Rover, tm.Reading, prev)
		}
		last[tm.Rover] = tm.Reading
		seen++
		if seen == stopAt {
			stop()
		}
		return nil
	}, 2)
	stop()
	r.Close()
	if err != nil || n != stopAt {
		t.Fatal("first consumer", n, err)
	}
	r2 := newReader()
	defer r2.Close()
	runCtx2, stop2 := context.WithCancel(ctx)
	defer stop2()
	n2, err := Consume(runCtx2, r2, func(_ context.Context, tm Telemetry) error {
		if prev, ok := last[tm.Rover]; ok && tm.Reading <= prev {
			t.Errorf("redelivered committed message: rover %s %d after %d", tm.Rover, tm.Reading, prev)
		}
		last[tm.Rover] = tm.Reading
		seen++
		if seen == len(batch) {
			stop2()
		}
		return nil
	}, 2)
	if err != nil || n2 != len(batch)-stopAt {
		t.Fatal("second consumer resumes after committed offsets", n2, err)
	}
}
`

func init() {
	add(&Level{ID: "35-6-raft-election", Chapter: 35, Kind: "lab", Title: "الانتخابات", Goal: "خلّي محطات المدار تتفق على قائد واحد لكل فترة بمحاكاة Raft حتمية في الذاكرة، حتى لو القائد وقع", Setup: seededLabRoom,
		Requires: []string{"struct", "switch"}, Deps: []string{}, Bridge: []int{89}, Race: false, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch356Starter}},
		Lesson: `# انتخاب قائد المدار

## قبل ما تكتب
- لو محطتين قرروا في نفس اللحظة إن كل واحدة فيهم القائدة، مين يكسب؟
- محطة كانت مفصولة ورجعت وهي فاكرة نفسها القائدة؛ إزاي الباقيين يعرفوها إنها قديمة؟

## ليه قائد واحد
محطات المدار لازم تتفق على ترتيب الأوامر. أسهل طريقة إن محطة واحدة (القائد) هي اللي ترتب، والباقيين (followers) يمشوا وراها. Raft خوارزمية بتختار القائد ده وبتغيّره لو وقع، من غير ما يبقى فيه قائدين في نفس الوقت.

## الفترات (terms)
الزمن متقسم فترات مترقمة ~1~ و ~2~ و ~3~. كل فترة فيها انتخاب واحد، وفيها قائد واحد على الأكتر. كل رسالة شايلة ~Term~ بتاع اللي باعتها. القاعدة الأساسية: لو شفت ~Term~ أكبر من بتاعك، خده فوراً وارجع follower ونسى صوتك القديم. رسالة ~Term~ بتاعها أصغر من بتاعك جاية من حد قديم: ارفضها ورد عليه بالـ ~Term~ بتاعك عشان يتعلم.

## الأدوار
كل محطة ليها ~Role~ واحد من تلاتة: ~Follower~ و ~Candidate~ و ~Leader~. الـ ~iota~ بيدّي الثوابت أرقام متتالية من صفر:
~~~go
type Role int

const (
	Follower Role = iota // 0
	Candidate            // 1
	Leader               // 2
)
~~~

## مهلة الانتخاب
الـ follower بيستنى نبضات (heartbeats) من القائد. لو عدّت مهلة من غير نبضة، يبقى candidate: يزوّد الـ ~term~، يدّي صوته لنفسه، ويطلب أصوات من الكل. اللي ياخد أغلبية (أكتر من نص عدد المحطات كلها) يبقى القائد ويبعت نبضات على طول.

المهلة عشوائية عشان المحطات متبدأش انتخابات في نفس اللحظة وتفضل تقسّم الأصوات. العشوائية محقونة من ~*rand.Rand~ بـ seed، فالتست بيعرف يعيد نفس السيناريو بالظبط:
~~~go
timeout := ElectionMin + rng.Intn(ElectionSpread) // من 10 لـ 19 نبضة ساعة
~~~

## صوت واحد لكل فترة
كل محطة بتدّي صوت واحد بس في كل فترة، وبتفتكر لمين في ~votedFor~. نفس المرشح لو طلب تاني ياخد نفس الرد. ده اللي بيضمن إن مرشحين مياخدوش أغلبية في نفس الفترة: أي أغلبيتين لازم يشتركوا في محطة، والمحطة دي صوتت مرة واحدة.

## محاكاة حتمية
مفيش goroutines ولا شبكة ولا ساعة حقيقية. ~Tick()~ هي نبضة الساعة المنطقية، و ~Step(m)~ بتسلّم رسالة. الاتنين بيرجعوا الرسايل اللي المحطة عايزة تبعتها، والتست هو الشبكة: بيوصّلها، أو يأخرها، أو يرميها عشان يحاكي الانفصال. نفس المدخلات دايماً بتدّي نفس النتيجة.
~~~go
out := node.Tick()        // ممكن تكون طلبات تصويت
for _, m := range out {
	replies := nodes[m.To].Step(m)
	_ = replies // التست بيحطها في الطابور
}
~~~

## أخطاء شائعة
~~~go
if m.Kind == Vote { reply.Granted = true } // غلط: صوتين في نفس الفترة
if n.votes > len(n.peers)/2 { ... } // غلط: الأغلبية من كل المحطات، والمحطة نفسها منهم
time.Sleep(timeout) // غلط: الوقت هنا نبضات Tick بس
~~~

## المهمة
نفّذ ~NewNode(id, peers, rng)~ و ~State~ و ~Tick~ و ~Step~ بتوقيعات الـ starter، وسيب ~main~ فاضية. ~peers~ أرقام المحطات التانية من غير ~id~. القيم الابتدائية: ~Follower~، ~term~ صفر، مفيش صوت، القائد ~-1~.

المهلة ~ElectionMin + rng.Intn(ElectionSpread)~، بتتسحب مرة في ~NewNode~ ومرة مع بداية كل انتخاب، ومفيش سحب تاني. عدّاد ~elapsed~ بيزيد واحد في كل ~Tick~. لو مش قائد ووصل للمهلة: candidate، ~term~ زائد واحد، صوت لنفسه، القائد ~-1~، ~elapsed~ صفر، مهلة جديدة، وطلب ~Vote~ لكل peer بترتيب ~peers~. لو صوته لوحده أغلبية (محطة واحدة) يبقى قائد على طول.

القائد بيزوّد ~elapsed~ ولما يوصل ~HeartbeatEvery~ يصفّره ويبعت ~Append~ فاضية لكل peer. لما يكسب الانتخاب يبعت النبضات فوراً ويصفّر ~elapsed~. ~Tick~ و ~Step~ بيرجعوا شريحة فاضية مش ~nil~ لو مفيش رسايل.

~Step(m)~: ~Term~ أكبر يخلي المحطة follower بالـ ~term~ الجديد من غير صوت وقائد ~-1~. ~Vote~ يترد عليه دايماً بـ ~VoteReply~ فيه ~term~ المحطة؛ ~Granted~ بس لو نفس الـ ~term~ والصوت فاضي أو لنفس المرشح، والمنح بيصفّر ~elapsed~. ~Append~ يترد عليه بـ ~AppendReply~؛ لو نفس الـ ~term~: follower، القائد هو المرسل، ~elapsed~ صفر، و ~Success~ صح. ~VoteReply~ ممنوح لنفس الـ ~term~ والمحطة candidate يتحسب مرة واحدة لكل مرسل؛ الأغلبية تعمله قائد. ~State()~ بترجع الدور والـ ~term~ والقائد المعروف.
`,
		Hints: []string{"ابدأ بـ Step: افحص m.Term > n.term الأول، وبعدين switch على m.Kind وابني الرد بالـ term الحالي.", "Tick: زوّد elapsed؛ القائد يبعت Append كل HeartbeatEvery، وغيره لما يوصل المهلة يبدأ انتخاب ويسحب مهلة جديدة.", "خزّن الأصوات في map[int]bool عشان التكرار ميتحسبش، والأغلبية len(votes) > (len(peers)+1)/2."},
		En: &LevelText{Title: "Elections", Goal: "Make the orbital stations agree on one leader per term with a deterministic in-memory Raft simulation, even when the leader fails", Files: []LabFile{{Name: "main.go", Starter: ch356Starter}},
			Lesson: `# Electing the orbit leader

## Before you write
- If two stations decide at the same moment that each is the leader, who wins?
- A station that was cut off comes back believing it is the leader; how do the others tell it that it is stale?

## Why one leader
The orbital stations must agree on the order of commands. The simplest way is for one station (the leader) to decide the order and the others (followers) to follow it. Raft is an algorithm that chooses that leader and replaces it when it fails, without two leaders at the same time.

## Terms
Time is split into numbered terms ~1~, ~2~, ~3~. Each term has one election and at most one leader. Every message carries its sender's ~Term~. The core rule: when you see a ~Term~ larger than yours, adopt it at once, become a follower, and forget your old vote. A message with a smaller ~Term~ comes from someone stale: reject it and reply with your ~Term~ so it learns.

## Roles
Each station has one ~Role~ of three: ~Follower~, ~Candidate~ and ~Leader~. ~iota~ gives constants consecutive numbers starting at zero:
~~~go
type Role int

const (
	Follower Role = iota // 0
	Candidate            // 1
	Leader               // 2
)
~~~

## The election timeout
A follower waits for heartbeats from the leader. If a timeout passes without one, it becomes a candidate: it increments its ~term~, votes for itself, and requests votes from everyone. Whoever gets a majority (more than half of all stations) becomes leader and sends heartbeats right away.

The timeout is random so stations do not start elections at the same moment and keep splitting the vote. Randomness is injected as a seeded ~*rand.Rand~, so the test can replay exactly the same scenario:
~~~go
timeout := ElectionMin + rng.Intn(ElectionSpread) // 10 to 19 clock ticks
~~~

## One vote per term
Each station grants one vote per term and remembers whom in ~votedFor~. The same candidate asking again gets the same answer. That guarantees two candidates cannot both get a majority in one term: any two majorities share a station, and that station voted once.

## A deterministic simulation
No goroutines, no network, no real clock. ~Tick()~ is one tick of the logical clock and ~Step(m)~ delivers a message. Both return the messages the station wants to send, and the test is the network: it delivers, delays, or drops them to simulate partitions. The same inputs always give the same result.
~~~go
out := node.Tick()        // possibly vote requests
for _, m := range out {
	replies := nodes[m.To].Step(m)
	_ = replies // the test queues them
}
~~~

## Common mistakes
~~~go
if m.Kind == Vote { reply.Granted = true } // Wrong: two votes in one term
if n.votes > len(n.peers)/2 { ... } // Wrong: a majority of all stations, the station itself included
time.Sleep(timeout) // Wrong: time here is Tick calls only
~~~

## Your mission
Implement ~NewNode(id, peers, rng)~, ~State~, ~Tick~ and ~Step~ with the starter signatures; leave ~main~ empty. ~peers~ are the other station ids, without ~id~. Initial values: ~Follower~, ~term~ zero, no vote, leader ~-1~.

The timeout is ~ElectionMin + rng.Intn(ElectionSpread)~, drawn once in ~NewNode~ and once at the start of each election, with no other draws. An ~elapsed~ counter grows by one per ~Tick~. When not leader and it reaches the timeout: candidate, ~term~ plus one, vote for itself, leader ~-1~, ~elapsed~ zero, a new timeout, and a ~Vote~ request to each peer in ~peers~ order. If its own vote is already a majority (one station) it becomes leader at once.

A leader increments ~elapsed~ and on reaching ~HeartbeatEvery~ resets it and sends an empty ~Append~ to each peer. On winning an election it sends heartbeats immediately and resets ~elapsed~. ~Tick~ and ~Step~ return an empty non-nil slice when there are no messages.

~Step(m)~: a larger ~Term~ makes the station a follower in the new ~term~ with no vote and leader ~-1~. A ~Vote~ always gets a ~VoteReply~ carrying the station's ~term~; ~Granted~ only for the same ~term~ when the vote is empty or already for that candidate, and granting resets ~elapsed~. An ~Append~ gets an ~AppendReply~; for the same ~term~: follower, the sender is leader, ~elapsed~ zero, and ~Success~ true. A granted ~VoteReply~ for the same ~term~ while candidate counts once per sender; a majority makes it leader. ~State()~ returns the role, the ~term~ and the known leader.
`, Hints: []string{"Start with Step: check m.Term > n.term first, then switch on m.Kind and build the reply with the current term.", "Tick: increment elapsed; a leader sends Append every HeartbeatEvery, anyone else starts an election at the timeout and draws a new one.", "Store votes in a map[int]bool so duplicates do not count, and a majority is len(votes) > (len(peers)+1)/2."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch356Hidden},
	})
}

const ch356Starter = `package main

import "math/rand"

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

type Kind int

const (
	Vote Kind = iota
	VoteReply
	Append
	AppendReply
)

type Msg struct {
	Kind    Kind
	From    int
	To      int
	Term    int
	Granted bool
	Success bool
}

const (
	ElectionMin    = 10
	ElectionSpread = 10
	HeartbeatEvery = 3
)

type Node struct{}

func NewNode(id int, peers []int, rng *rand.Rand) *Node { return &Node{} }
func (n *Node) State() (Role, int, int)                { return Follower, 0, -1 }
func (n *Node) Tick() []Msg                            { return nil }
func (n *Node) Step(m Msg) []Msg                       { return nil }
func main()                                            {}
`

const ch356Hidden = `package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

// hiddenNet is the test-controlled network: FIFO delivery, with cut links.
type hiddenNet struct {
	t     *testing.T
	nodes []*Node
	down  map[int]bool // stopped stations neither tick nor receive
	cut   map[[2]int]bool
	queue []Msg
	terms map[int]int // term -> leader seen
}

func hiddenCluster(t *testing.T, size int, seed int64) *hiddenNet {
	net := &hiddenNet{t: t, down: map[int]bool{}, cut: map[[2]int]bool{}, terms: map[int]int{}}
	for id := 0; id < size; id++ {
		var peers []int
		for p := 0; p < size; p++ {
			if p != id {
				peers = append(peers, p)
			}
		}
		net.nodes = append(net.nodes, NewNode(id, peers, rand.New(rand.NewSource(seed+int64(id)))))
	}
	return net
}
func (net *hiddenNet) partition(groups ...[]int) {
	net.cut = map[[2]int]bool{}
	group := map[int]int{}
	for g, ids := range groups {
		for _, id := range ids {
			group[id] = g
		}
	}
	for a := range net.nodes {
		for b := range net.nodes {
			if group[a] != group[b] {
				net.cut[[2]int{a, b}] = true
			}
		}
	}
}
func (net *hiddenNet) check() {
	for id, n := range net.nodes {
		role, term, _ := n.State()
		if role == Leader {
			if other, ok := net.terms[term]; ok && other != id {
				net.t.Fatalf("two leaders in term %d: %d and %d", term, other, id)
			}
			net.terms[term] = id
		}
	}
}
func (net *hiddenNet) round() {
	for id, n := range net.nodes {
		if !net.down[id] {
			out := n.Tick()
			if out == nil {
				net.t.Fatal("Tick returned nil")
			}
			net.queue = append(net.queue, out...)
		}
	}
	for len(net.queue) > 0 {
		m := net.queue[0]
		net.queue = net.queue[1:]
		if net.down[m.To] || net.down[m.From] || net.cut[[2]int{m.From, m.To}] {
			continue
		}
		net.queue = append(net.queue, net.nodes[m.To].Step(m)...)
		net.check()
	}
	net.check()
}
func (net *hiddenNet) leader(within int, among []int) int {
	for i := 0; i < within; i++ {
		net.round()
		for _, id := range among {
			if role, _, _ := net.nodes[id].State(); role == Leader {
				leader, agree := id, true
				_, term, _ := net.nodes[id].State()
				for _, other := range among {
					_, t2, l2 := net.nodes[other].State()
					agree = agree && t2 == term && l2 == leader
				}
				if agree {
					return leader
				}
			}
		}
	}
	net.t.Fatal("no agreed leader")
	return -1
}

func TestHiddenRaftNodeRules(t *testing.T) {
	n := NewNode(2, []int{0, 1, 3, 4}, rand.New(rand.NewSource(77)))
	if role, term, leader := n.State(); role != Follower || term != 0 || leader != -1 {
		t.Fatal("initial", role, term, leader)
	}
	rng := rand.New(rand.NewSource(77))
	timeout := ElectionMin + rng.Intn(ElectionSpread)
	for i := 1; i < timeout; i++ {
		if out := n.Tick(); out == nil || len(out) != 0 {
			t.Fatal("early election", i, out)
		}
	}
	out := n.Tick()
	want := []Msg{{Kind: Vote, From: 2, To: 0, Term: 1}, {Kind: Vote, From: 2, To: 1, Term: 1}, {Kind: Vote, From: 2, To: 3, Term: 1}, {Kind: Vote, From: 2, To: 4, Term: 1}}
	if !reflect.DeepEqual(out, want) {
		t.Fatal("vote requests", out)
	}
	if role, term, leader := n.State(); role != Candidate || term != 1 || leader != -1 {
		t.Fatal("candidate", role, term, leader)
	}
	if got := n.Step(Msg{Kind: Vote, From: 3, To: 2, Term: 1}); !reflect.DeepEqual(got, []Msg{{Kind: VoteReply, From: 2, To: 3, Term: 1}}) {
		t.Fatal("already voted for itself", got)
	}
	n.Step(Msg{Kind: VoteReply, From: 0, To: 2, Term: 1, Granted: true})
	n.Step(Msg{Kind: VoteReply, From: 0, To: 2, Term: 1, Granted: true})
	if role, _, _ := n.State(); role != Candidate {
		t.Fatal("duplicate vote counted twice")
	}
	n.Step(Msg{Kind: VoteReply, From: 1, To: 2, Term: 1, Granted: false})
	beats := n.Step(Msg{Kind: VoteReply, From: 4, To: 2, Term: 1, Granted: true})
	if role, term, leader := n.State(); role != Leader || term != 1 || leader != 2 || len(beats) != 4 || beats[0] != (Msg{Kind: Append, From: 2, To: 0, Term: 1}) {
		t.Fatal("majority", role, term, leader, beats)
	}
	if got := n.Tick(); len(got) != 0 {
		t.Fatal("heartbeat too early", got)
	}
	n.Tick()
	if got := n.Tick(); len(got) != 4 || got[3] != (Msg{Kind: Append, From: 2, To: 4, Term: 1}) {
		t.Fatal("heartbeat", got)
	}
	if got := n.Step(Msg{Kind: Append, From: 0, To: 2, Term: 0}); !reflect.DeepEqual(got, []Msg{{Kind: AppendReply, From: 2, To: 0, Term: 1}}) {
		t.Fatal("stale append", got)
	}
	if got := n.Step(Msg{Kind: Vote, From: 1, To: 2, Term: 0}); !reflect.DeepEqual(got, []Msg{{Kind: VoteReply, From: 2, To: 1, Term: 1}}) {
		t.Fatal("stale vote", got)
	}
	if got := n.Step(Msg{Kind: Vote, From: 3, To: 2, Term: 5}); !reflect.DeepEqual(got, []Msg{{Kind: VoteReply, From: 2, To: 3, Term: 5, Granted: true}}) {
		t.Fatal("higher term vote", got)
	}
	if role, term, leader := n.State(); role != Follower || term != 5 || leader != -1 {
		t.Fatal("step down", role, term, leader)
	}
	if got := n.Step(Msg{Kind: Vote, From: 4, To: 2, Term: 5}); got[0].Granted {
		t.Fatal("second vote in term")
	}
	if got := n.Step(Msg{Kind: Vote, From: 3, To: 2, Term: 5}); !got[0].Granted {
		t.Fatal("repeat request from same candidate")
	}
	if got := n.Step(Msg{Kind: Append, From: 3, To: 2, Term: 5}); !reflect.DeepEqual(got, []Msg{{Kind: AppendReply, From: 2, To: 3, Term: 5, Success: true}}) {
		t.Fatal("append", got)
	}
	if _, _, leader := n.State(); leader != 3 {
		t.Fatal("leader learned", leader)
	}
	solo := NewNode(9, nil, rand.New(rand.NewSource(1)))
	for i := 0; i < ElectionMin+ElectionSpread; i++ {
		solo.Tick()
	}
	if role, term, _ := solo.State(); role != Leader || term != 1 {
		t.Fatal("single station", role, term)
	}
}

func TestHiddenRaftClusterElectsOneLeader(t *testing.T) {
	for seed := int64(3560); seed < 3570; seed++ {
		net := hiddenCluster(t, 5, seed*10)
		all := []int{0, 1, 2, 3, 4}
		leader := net.leader(200, all)
		for i := 0; i < 100; i++ {
			net.round()
		}
		_, term, _ := net.nodes[leader].State()
		for _, n := range net.nodes {
			if r, t2, l := n.State(); t2 != term || l != leader || (r == Leader) != (n == net.nodes[leader]) {
				t.Fatal(fmt.Sprintf("seed %d: stable leader lost", seed), r, t2, l)
			}
		}
	}
}

func TestHiddenRaftLeaderFailureAndStaleReturn(t *testing.T) {
	for seed := int64(1); seed <= 6; seed++ {
		net := hiddenCluster(t, 5, seed*1000)
		all := []int{0, 1, 2, 3, 4}
		old := net.leader(200, all)
		_, oldTerm, _ := net.nodes[old].State()
		var rest []int
		for _, id := range all {
			if id != old {
				rest = append(rest, id)
			}
		}
		net.partition([]int{old}, rest)
		next := net.leader(200, rest)
		_, newTerm, _ := net.nodes[next].State()
		if next == old || newTerm <= oldTerm {
			t.Fatal("new leader", next, newTerm, oldTerm)
		}
		if role, term, _ := net.nodes[old].State(); role != Leader || term != oldTerm {
			t.Fatal("isolated leader keeps its stale view", role, term)
		}
		net.partition(all)
		net.leader(200, all)
		if role, term, _ := net.nodes[old].State(); role == Leader || term < newTerm {
			t.Fatal("stale leader must step down", role, term)
		}
	}
}

func TestHiddenRaftMinorityCannotElect(t *testing.T) {
	net := hiddenCluster(t, 5, 4242)
	all := []int{0, 1, 2, 3, 4}
	net.leader(200, all)
	_, before, _ := net.nodes[0].State()
	net.partition([]int{0, 1}, []int{2, 3, 4})
	for i := 0; i < 300; i++ {
		net.round()
	}
	for _, id := range []int{0, 1} {
		if role, term, _ := net.nodes[id].State(); role == Leader && term > before {
			t.Fatal("minority of two elected a leader", id, term)
		}
	}
	_, alone, _ := net.nodes[4].State()
	net.down[2] = true
	net.down[3] = true
	for i := 0; i < 300; i++ {
		net.round()
	}
	if role, term, _ := net.nodes[4].State(); role == Leader || term <= alone {
		t.Fatal("a lone station keeps campaigning but never leads", role, term, alone)
	}
}
`

func init() {
	add(&Level{ID: "35-7-raft-log", Chapter: 35, Kind: "lab", Title: "نسخ السجل", Goal: "خلّي القائد ينسخ أوامر المدار لكل المحطات، ويأكدها بالأغلبية، وكل محطة تطبقها بنفس الترتيب", Setup: seededLabRoom,
		Requires: []string{"struct", "switch"}, Deps: []string{}, Bridge: []int{89}, Race: false, Service: "", CarryFrom: "35-6-raft-election",
		Files: []LabFile{{Name: "main.go", Starter: ch357Starter}},
		Lesson: `# سجل أوامر المدار

## قبل ما تكتب
- القائد كتب أمر عنده ووقع قبل ما حد ياخده؛ الأمر ده يعتبر حصل ولا لأ؟
- محطة رجعت بعد انفصال وعندها أوامر محدش تاني شافها؛ تعمل فيها إيه؟

## السجل المنسوخ
كود الانتخابات بتاعك من المهمة اللي فاتت هو نقطة البداية؛ هتكمّل في نفس الملف. القائد دلوقتي مش بس بيبعت نبضات: بيستقبل أوامر (~Propose~)، يكتبها في سجل، وينسخ السجل للـ followers. كل محطة بتطبق الأوامر المؤكدة بنفس الترتيب، فكلهم يوصلوا لنفس الحالة. ده اسمه state machine replication.

كل خانة في السجل ~Entry~ فيها الأمر والـ ~Term~ اللي اتكتب فيه. الترقيم بيبدأ من ~1~، و ~0~ معناها «قبل أول خانة». لو السجل شريحة، الخانة ~i~ في ~log[i-1]~:
~~~go
log := []Entry{{Term: 1, Cmd: "dock"}, {Term: 2, Cmd: "vent"}}
last := len(log)            // 2
lastTerm := log[last-1].Term // 2
~~~

## فحص التطابق
رسالة ~Append~ بقت شايلة ~PrevIndex~ و ~PrevTerm~: الخانة اللي قبل الخانات الجديدة على طول. الـ follower بيقبل بس لو عنده خانة رقمها ~PrevIndex~ بنفس ~PrevTerm~. لو مفيش، يرد ~Success~ غلط، والقائد يرجع خطوة لورا في ~next~ ويحاول تاني. كده السجلين بيتطابقوا لحد نقطة، وبعدها القائد يبعت الباقي.

لو خانة جديدة بتتعارض مع خانة موجودة (نفس الرقم، ~Term~ مختلف)، الـ follower بيقص سجله من عند الخانة دي ويكمّل بخانات القائد. الخانات اللي اتقصت عمرها ما كانت مؤكدة.

## الأغلبية بتأكد
القائد بيفتكر لكل peer ~match~: أعلى خانة اتأكد إنها عنده، و ~next~: الخانة الجاية اللي هيبعتها. الخانة تبقى committed لما تبقى متخزنة على أغلبية المحطات (القائد منهم). القائد بيبعت ~Commit~ في كل ~Append~، والـ follower يؤكد لحد ~min(Commit, آخر خانة اتطابقت)~.

قاعدة مهمة: القائد بيأكد بالعدّ بس خانات من الـ ~term~ بتاعه. خانة قديمة على أغلبية ممكن لسه تتمسح لو قائد تاني كسب؛ بتتأكد ضمنياً لما خانة من الـ ~term~ الحالي فوقها تتأكد.

## مين ينفع يبقى قائد
محطة سجلها ناقص لو بقت قائد هتمسح أوامر مؤكدة. عشان كده ~Vote~ بقت شايلة ~LastIndex~ و ~LastTerm~ بتوع المرشح. المحطة تدّي صوتها بس لو سجل المرشح مش أقدم من سجلها: ~LastTerm~ أكبر، أو نفس الـ ~LastTerm~ و ~LastIndex~ أكبر أو مساوي.
~~~go
upToDate := m.LastTerm > myLastTerm || (m.LastTerm == myLastTerm && m.LastIndex >= myLastIndex)
~~~

## أخطاء شائعة
~~~go
n.log = append(n.log, m.Entries...) // غلط: من غير فحص PrevIndex/PrevTerm ومن غير قص التعارض
n.commit = m.Commit // غلط: ممكن يأكد خانات لسه موصلتش للمحطة دي
if count > half { n.commit = index } // غلط لو log[index-1].Term مش الـ term الحالي
~~~

## المهمة
كمّل على ~main.go~ بتاعك من الانتخابات. زوّد على ~Msg~ الحقول ~LastIndex~ و ~LastTerm~ و ~PrevIndex~ و ~PrevTerm~ و ~Entries []Entry~ و ~Commit~ و ~Match~ كلها ~int~ ما عدا ~Entries~، وعرّف ~Entry~ بـ ~Term int~ و ~Cmd string~، ونفّذ ~Propose(cmd) (int, bool)~ و ~Committed() []string~. كل قواعد الانتخابات من المهمة اللي فاتت لسه زي ما هي.

~Vote~ بيشيل ~LastIndex~ و ~LastTerm~ (صفر لسجل فاضي)، والصوت يتمنح بس لمرشح سجله مش أقدم. لما محطة تكسب: ~next~ لكل peer آخر خانة زائد واحد، و ~match~ صفر. كل ~Append~ (عند الفوز وكل نبضة) لكل peer فيها ~PrevIndex = next-1~ و ~PrevTerm~ بتاعتها، و ~Entries~ نسخة من كل الخانات من ~next~ لآخر السجل، و ~Commit~ القائد.

الـ follower في نفس الـ ~term~ بيعمل كل اللي كان بيعمله، وبعدين: لو ~PrevIndex~ أكبر من سجله أو ~Term~ الخانة دي مش ~PrevTerm~ يرد ~Success~ غلط. غير كده يقص عند أول تعارض، يضيف الناقص، يرد ~Success~ صح و ~Match = PrevIndex + len(Entries)~، ويزوّد ~commit~ لـ ~min(Commit, Match)~ لو أكبر. ~commit~ عمره ما يقل.

القائد لما ياخد ~AppendReply~ لنفس الـ ~term~: النجاح يخلي ~match~ أكبر قيمة شافها و ~next = match + 1~ ويحاول يزوّد ~commit~؛ الفشل ينقص ~next~ واحد بحد أدنى ~1~. ~commit~ يبقى أعلى خانة من الـ ~term~ الحالي على أغلبية.

~Propose~ لغير القائد ترجع ~0~ و ~false~؛ للقائد تضيف ~Entry~ بالـ ~term~ الحالي وترجع رقمها و ~true~، والإرسال بيحصل في النبضة الجاية. لو المحطة لوحدها في المجموعة، الخانة تتأكد على طول. ~Committed~ ترجع أوامر الخانات المؤكدة بالترتيب، وشريحة فاضية مش ~nil~ لو مفيش.
`,
		Hints: []string{"زوّد الحقول والأنواع الأول عشان التست يتبني، وخلي Vote يبعت LastIndex/LastTerm ويشترط upToDate قبل المنح.", "عند الفوز جهّز next و match، وخلي broadcast للـ Append يحسب PrevIndex و PrevTerm و Entries من next لكل peer.", "في الـ follower افحص التطابق، اقص التعارض، أضف الناقص، ورد بـ Match؛ في القائد حدّث match/next وعدّ الأغلبية لخانات الـ term الحالي بس."},
		En: &LevelText{Title: "Log Replication", Goal: "Have the leader replicate orbit commands to every station, commit them by majority, and have each station apply them in the same order", Files: []LabFile{{Name: "main.go", Starter: ch357Starter}},
			Lesson: `# The orbit command log

## Before you write
- The leader wrote a command locally and crashed before anyone received it; did that command happen?
- A station returns from a partition with commands nobody else has seen; what does it do with them?

## A replicated log
Your election code from the previous mission is the starting point; you continue in the same file. The leader now does more than send heartbeats: it accepts commands (~Propose~), writes them to a log, and replicates the log to followers. Every station applies committed commands in the same order, so they all reach the same state. This is state machine replication.

Each log slot ~Entry~ holds the command and the ~Term~ it was written in. Numbering starts at ~1~, and ~0~ means "before the first slot". With the log as a slice, slot ~i~ lives at ~log[i-1]~:
~~~go
log := []Entry{{Term: 1, Cmd: "dock"}, {Term: 2, Cmd: "vent"}}
last := len(log)            // 2
lastTerm := log[last-1].Term // 2
~~~

## The consistency check
An ~Append~ now carries ~PrevIndex~ and ~PrevTerm~: the slot right before the new ones. A follower accepts only if it has slot ~PrevIndex~ with that ~PrevTerm~. If not, it replies ~Success~ false and the leader steps ~next~ back one and tries again. That way the two logs match up to a point, and the leader then sends the rest.

When a new entry conflicts with an existing one (same index, different ~Term~), the follower truncates its log from that slot and continues with the leader's entries. The truncated entries were never committed.

## A majority commits
The leader remembers, per peer, ~match~: the highest slot known to be stored there, and ~next~: the next slot it will send. A slot is committed once it is stored on a majority of stations (the leader included). The leader sends ~Commit~ in each ~Append~, and the follower commits up to ~min(Commit, last matched slot)~.

An important rule: the leader commits by counting only entries from its own ~term~. An older entry on a majority could still be overwritten if another leader wins; it becomes committed implicitly once a current-term entry above it commits.

## Who may lead
A station with a missing log would erase committed commands if it became leader. So ~Vote~ now carries the candidate's ~LastIndex~ and ~LastTerm~. A station grants its vote only if the candidate's log is not older than its own: a larger ~LastTerm~, or the same ~LastTerm~ and a ~LastIndex~ at least as large.
~~~go
upToDate := m.LastTerm > myLastTerm || (m.LastTerm == myLastTerm && m.LastIndex >= myLastIndex)
~~~

## Common mistakes
~~~go
n.log = append(n.log, m.Entries...) // Wrong: no PrevIndex/PrevTerm check and no conflict truncation
n.commit = m.Commit // Wrong: may commit slots this station does not have yet
if count > half { n.commit = index } // Wrong when log[index-1].Term is not the current term
~~~

## Your mission
Continue in your ~main.go~ from the elections. Add to ~Msg~ the fields ~LastIndex~, ~LastTerm~, ~PrevIndex~, ~PrevTerm~, ~Entries []Entry~, ~Commit~ and ~Match~, all ~int~ except ~Entries~; define ~Entry~ with ~Term int~ and ~Cmd string~; and implement ~Propose(cmd) (int, bool)~ and ~Committed() []string~. Every election rule from the previous mission still holds.

~Vote~ carries ~LastIndex~ and ~LastTerm~ (zero for an empty log), and a vote is granted only to a candidate whose log is not older. On winning: ~next~ for each peer is the last index plus one, and ~match~ is zero. Each ~Append~ (on winning and on each heartbeat) to each peer has ~PrevIndex = next-1~ with its ~PrevTerm~, ~Entries~ a copy of every slot from ~next~ to the end, and the leader's ~Commit~.

A follower in the same ~term~ does everything it did before, and then: if ~PrevIndex~ is beyond its log or that slot's ~Term~ is not ~PrevTerm~, it replies ~Success~ false. Otherwise it truncates at the first conflict, appends what is missing, replies ~Success~ true with ~Match = PrevIndex + len(Entries)~, and raises ~commit~ to ~min(Commit, Match)~ when larger. ~commit~ never decreases.

A leader receiving an ~AppendReply~ for its ~term~: success sets ~match~ to the largest value seen and ~next = match + 1~, then tries to advance ~commit~; failure decrements ~next~ by one, never below ~1~. ~commit~ is the highest current-term slot stored on a majority.

~Propose~ on a non-leader returns ~0~ and ~false~; on the leader it appends an ~Entry~ with the current ~term~ and returns its index and ~true~, and sending happens on the next heartbeat. A station alone in its cluster commits it at once. ~Committed~ returns the commands of committed slots in order, an empty non-nil slice when none.
`, Hints: []string{"Add the fields and types first so the test builds, and make Vote send LastIndex/LastTerm and require upToDate before granting.", "On winning set up next and match, and have the Append broadcast compute PrevIndex, PrevTerm and Entries from next per peer.", "In the follower check the match, truncate conflicts, append what is missing, and reply with Match; in the leader update match/next and count a majority for current-term slots only."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch357Hidden},
	})
}

const ch357Starter = `package main

import "math/rand"

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

type Kind int

const (
	Vote Kind = iota
	VoteReply
	Append
	AppendReply
)

type Entry struct {
	Term int
	Cmd  string
}

type Msg struct {
	Kind      Kind
	From      int
	To        int
	Term      int
	Granted   bool
	Success   bool
	LastIndex int
	LastTerm  int
	PrevIndex int
	PrevTerm  int
	Entries   []Entry
	Commit    int
	Match     int
}

const (
	ElectionMin    = 10
	ElectionSpread = 10
	HeartbeatEvery = 3
)

// Bring your election code from the previous mission here.
type Node struct{}

func NewNode(id int, peers []int, rng *rand.Rand) *Node { return &Node{} }
func (n *Node) State() (Role, int, int)                { return Follower, 0, -1 }
func (n *Node) Tick() []Msg                            { return nil }
func (n *Node) Step(m Msg) []Msg                       { return nil }
func (n *Node) Propose(cmd string) (int, bool)         { return 0, false }
func (n *Node) Committed() []string                    { return nil }
func main()                                            {}
`

const ch357Hidden = `package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

type hiddenNet struct {
	t     *testing.T
	nodes []*Node
	down  map[int]bool
	cut   map[[2]int]bool
	queue []Msg
}

func hiddenCluster(t *testing.T, size int, seed int64) *hiddenNet {
	net := &hiddenNet{t: t, down: map[int]bool{}, cut: map[[2]int]bool{}}
	for id := 0; id < size; id++ {
		var peers []int
		for p := 0; p < size; p++ {
			if p != id {
				peers = append(peers, p)
			}
		}
		net.nodes = append(net.nodes, NewNode(id, peers, rand.New(rand.NewSource(seed+int64(id)))))
	}
	return net
}
func (net *hiddenNet) partition(groups ...[]int) {
	net.cut = map[[2]int]bool{}
	group := map[int]int{}
	for g, ids := range groups {
		for _, id := range ids {
			group[id] = g
		}
	}
	for a := range net.nodes {
		for b := range net.nodes {
			if group[a] != group[b] {
				net.cut[[2]int{a, b}] = true
			}
		}
	}
}
func hiddenPrefix(a, b []string) bool {
	if len(a) > len(b) {
		a, b = b, a
	}
	return reflect.DeepEqual(a, b[:len(a)])
}
func (net *hiddenNet) safety() {
	for i, a := range net.nodes {
		for j, b := range net.nodes {
			if i < j && !hiddenPrefix(a.Committed(), b.Committed()) {
				net.t.Fatalf("committed logs diverge: %d=%v %d=%v", i, a.Committed(), j, b.Committed())
			}
		}
	}
}
func (net *hiddenNet) round() {
	for id, n := range net.nodes {
		if !net.down[id] {
			net.queue = append(net.queue, n.Tick()...)
		}
	}
	for len(net.queue) > 0 {
		m := net.queue[0]
		net.queue = net.queue[1:]
		if net.down[m.To] || net.down[m.From] || net.cut[[2]int{m.From, m.To}] {
			continue
		}
		net.queue = append(net.queue, net.nodes[m.To].Step(m)...)
	}
	net.safety()
}
func (net *hiddenNet) leader(within int, among []int) int {
	for i := 0; i < within; i++ {
		net.round()
		for _, id := range among {
			role, term, _ := net.nodes[id].State()
			if role != Leader {
				continue
			}
			agree := true
			for _, other := range among {
				_, t2, l2 := net.nodes[other].State()
				agree = agree && t2 == term && l2 == id
			}
			if agree {
				return id
			}
		}
	}
	net.t.Fatal("no agreed leader")
	return -1
}

func TestHiddenRaftFollowerAppendRules(t *testing.T) {
	f := NewNode(1, []int{0, 2}, rand.New(rand.NewSource(9)))
	if got := f.Committed(); got == nil || len(got) != 0 {
		t.Fatal("empty committed must be non-nil", got)
	}
	r := f.Step(Msg{Kind: Append, From: 0, To: 1, Term: 1, PrevIndex: 2, PrevTerm: 1, Entries: []Entry{{1, "x"}}})
	if len(r) != 1 || r[0].Success || r[0].Kind != AppendReply {
		t.Fatal("gap accepted", r)
	}
	r = f.Step(Msg{Kind: Append, From: 0, To: 1, Term: 1, Entries: []Entry{{1, "a"}, {1, "b"}, {1, "c"}}, Commit: 1})
	if !r[0].Success || r[0].Match != 3 || !reflect.DeepEqual(f.Committed(), []string{"a"}) {
		t.Fatal("append", r, f.Committed())
	}
	r = f.Step(Msg{Kind: Append, From: 0, To: 1, Term: 1, PrevIndex: 3, PrevTerm: 2})
	if r[0].Success {
		t.Fatal("term mismatch accepted")
	}
	r = f.Step(Msg{Kind: Append, From: 0, To: 1, Term: 1, PrevIndex: 1, PrevTerm: 1, Commit: 9})
	if !r[0].Success || r[0].Match != 1 || !reflect.DeepEqual(f.Committed(), []string{"a"}) {
		t.Fatal("commit limited by match", r, f.Committed())
	}
	r = f.Step(Msg{Kind: Append, From: 2, To: 1, Term: 2, PrevIndex: 1, PrevTerm: 1, Entries: []Entry{{1, "b"}, {2, "z"}}, Commit: 3})
	if !r[0].Success || r[0].Match != 3 || !reflect.DeepEqual(f.Committed(), []string{"a", "b", "z"}) {
		t.Fatal("conflict truncation", r, f.Committed())
	}
	r = f.Step(Msg{Kind: Append, From: 2, To: 1, Term: 2, PrevIndex: 0, PrevTerm: 0, Commit: 1})
	if !r[0].Success || len(f.Committed()) != 3 {
		t.Fatal("commit must never decrease", f.Committed())
	}
	v := f.Step(Msg{Kind: Vote, From: 0, To: 1, Term: 3, LastIndex: 5, LastTerm: 1})
	if v[0].Granted {
		t.Fatal("vote for an older log")
	}
	v = f.Step(Msg{Kind: Vote, From: 0, To: 1, Term: 4, LastIndex: 2, LastTerm: 2})
	if v[0].Granted {
		t.Fatal("vote for a shorter log in the same last term")
	}
	v = f.Step(Msg{Kind: Vote, From: 2, To: 1, Term: 5, LastIndex: 3, LastTerm: 2})
	if !v[0].Granted {
		t.Fatal("equal log is up to date")
	}
}

func TestHiddenRaftLeaderReplicationAndCommit(t *testing.T) {
	n := NewNode(0, []int{1, 2}, rand.New(rand.NewSource(5)))
	if i, ok := n.Propose("x"); ok || i != 0 {
		t.Fatal("follower accepted proposal")
	}
	n.Step(Msg{Kind: Append, From: 2, To: 0, Term: 1, Entries: []Entry{{1, "old"}}})
	for i := 0; i < 3*(ElectionMin+ElectionSpread); i++ {
		if role, _, _ := n.State(); role == Candidate {
			break
		}
		n.Tick()
	}
	role, term, _ := n.State()
	if role != Candidate || term != 2 {
		t.Fatal("candidate", role, term)
	}
	out := n.Step(Msg{Kind: VoteReply, From: 1, To: 0, Term: 2, Granted: true})
	if len(out) != 2 || out[0].PrevIndex != 1 || out[0].PrevTerm != 1 || len(out[0].Entries) != 0 {
		t.Fatal("first heartbeat", out)
	}
	n.Step(Msg{Kind: AppendReply, From: 1, To: 0, Term: 2, Success: true, Match: 1})
	if len(n.Committed()) != 0 {
		t.Fatal("older-term entry committed by counting")
	}
	idx, ok := n.Propose("vent")
	if !ok || idx != 2 {
		t.Fatal("propose", idx, ok)
	}
	n.Tick()
	n.Tick()
	beats := n.Tick()
	if len(beats) != 2 || beats[0].PrevIndex != 1 || !reflect.DeepEqual(beats[0].Entries, []Entry{{2, "vent"}}) || beats[1].PrevIndex != 1 {
		t.Fatal("heartbeat carries entries", beats)
	}
	n.Step(Msg{Kind: AppendReply, From: 2, To: 0, Term: 2, Success: false})
	n.Tick()
	n.Tick()
	beats = n.Tick()
	if beats[1].PrevIndex != 0 || len(beats[1].Entries) != 2 {
		t.Fatal("next backs off after a rejection", beats[1])
	}
	n.Step(Msg{Kind: AppendReply, From: 1, To: 0, Term: 2, Success: true, Match: 2})
	if !reflect.DeepEqual(n.Committed(), []string{"old", "vent"}) {
		t.Fatal("majority commit", n.Committed())
	}
	n.Step(Msg{Kind: AppendReply, From: 1, To: 0, Term: 2, Success: true, Match: 1})
	n.Tick()
	n.Tick()
	if beats = n.Tick(); beats[0].PrevIndex != 2 || beats[0].Commit != 2 {
		t.Fatal("match never goes back", beats[0])
	}
	solo := NewNode(7, nil, rand.New(rand.NewSource(3)))
	for i := 0; i < ElectionMin+ElectionSpread; i++ {
		solo.Tick()
	}
	if i, ok := solo.Propose("alone"); !ok || i != 1 || !reflect.DeepEqual(solo.Committed(), []string{"alone"}) {
		t.Fatal("single station commit", i, ok, solo.Committed())
	}
}

func TestHiddenRaftIsolatedLeaderEntriesOverwritten(t *testing.T) {
	net := hiddenCluster(t, 5, 3570)
	all := []int{0, 1, 2, 3, 4}
	old := net.leader(200, all)
	net.nodes[old].Propose("base")
	for i := 0; i < 10; i++ {
		net.round()
	}
	var rest []int
	for _, id := range all {
		if id != old {
			rest = append(rest, id)
		}
	}
	net.partition([]int{old}, rest)
	if _, ok := net.nodes[old].Propose("lost"); !ok {
		t.Fatal("isolated leader still believes it leads")
	}
	next := net.leader(200, rest)
	net.nodes[next].Propose("kept")
	for i := 0; i < 20; i++ {
		net.round()
	}
	net.partition(all)
	for i := 0; i < 60; i++ {
		net.round()
	}
	for id, n := range net.nodes {
		if got := n.Committed(); !reflect.DeepEqual(got, []string{"base", "kept"}) {
			t.Fatal("station", id, got)
		}
	}
}

func TestHiddenRaftRandomPartitions(t *testing.T) {
	for seed := int64(1); seed <= 4; seed++ {
		rng := rand.New(rand.NewSource(35700 + seed))
		net := hiddenCluster(t, 5, seed*7919)
		all := []int{0, 1, 2, 3, 4}
		ever := map[string]bool{}
		cmd := 0
		for r := 0; r < 600; r++ {
			if r%60 == 0 {
				if rng.Intn(3) == 0 {
					net.partition(all)
				} else {
					perm := rng.Perm(5)
					cut := 1 + rng.Intn(2)
					net.partition(perm[:cut], perm[cut:])
				}
			}
			if rng.Intn(3) == 0 {
				for _, n := range net.nodes {
					if role, _, _ := n.State(); role == Leader {
						cmd++
						n.Propose(fmt.Sprintf("c%d", cmd))
					}
				}
			}
			net.round()
			for _, n := range net.nodes {
				for _, c := range n.Committed() {
					ever[c] = true
				}
			}
		}
		net.partition(all)
		leader := net.leader(300, all)
		net.nodes[leader].Propose("final")
		converged := func() bool {
			want := net.nodes[leader].Committed()
			for _, n := range net.nodes {
				if !reflect.DeepEqual(n.Committed(), want) {
					return false
				}
			}
			return len(want) > 0 && want[len(want)-1] == "final"
		}
		// A diverged follower is repaired one slot per heartbeat, so allow time.
		for i := 0; i < 3000 && !converged(); i++ {
			net.round()
		}
		want := net.nodes[leader].Committed()
		for id, n := range net.nodes {
			if !reflect.DeepEqual(n.Committed(), want) {
				t.Fatal("seed", seed, "station", id, "did not converge")
			}
		}
		seen := map[string]bool{}
		for _, c := range want {
			if seen[c] {
				t.Fatal("command applied twice", c)
			}
			seen[c] = true
		}
		for c := range ever {
			if !seen[c] {
				t.Fatal("committed command lost", c)
			}
		}
		if !seen["final"] || len(want) < 5 {
			t.Fatal("progress", len(want))
		}
	}
}
`
