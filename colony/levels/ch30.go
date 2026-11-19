package levels

func init() {
	add(&Level{LessonRevision: 1, Steps: ch30Lesson301Arabic, ID: "30-1-cache-aside", Chapter: 30, Kind: "lab", Title: "كاش جانبي", Goal: "سرّع قراية المنارات بـ Redis مع صلاحية وإلغاء آمن للكاش", Setup: seededLabRoom,
		Requires: []string{"iface", "mutex", "sel:redis.Nil"}, Deps: []string{"github.com/redis/go-redis/v9", "github.com/alicebob/miniredis/v2"}, Bridge: []int{65}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch301Starter}},
		Lesson: `# منارات بكاش جانبي

## قبل ما تكتب
- لو الكاش وقع، مين لسه مصدر الحقيقة؟
- لو الحفظ نجح وحذف الكاش فشل، هتقول إيه للمنادي؟

## نمط الكاش الجانبي
الكاش نسخة مؤقتة تسرّع القراية، مش مخزن الحقيقة. في نمط cache-aside، القراية بتسأل الكاش الأول؛ لو لقت قيمة (hit) ترجعها من غير المخزن. لو المفتاح غايب (miss) تقرا المخزن وتحاول تملا الكاش. فشل الملء ميحوّلش قراية ناجحة لفشل.

~BeaconStore~ واجهة المخزن المحقون. مثال استخدام: ~b, err := store.Load(ctx, "west")~ للقراية، و ~err = store.Save(ctx, Beacon{ID: "west", Charge: 0})~ للحفظ. ~Beacon~ فيها ~ID~ نص و ~Charge~ عدد غير سالب، والصفر مسموح.

## أوامر Redis الحقيقية
استورد ~github.com/redis/go-redis/v9~. ~redis.NewClient~ بتبني عميل من ~redis.Options~؛ ~Addr~ عنوان السيرفر. المنادي بيملك العميل وبيقفله بـ ~Close~. المهمة بتستقبل ~*redis.Client~ جاهز، ومبتفتحش سيرفر.
~~~go
client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
defer client.Close()
raw, err := client.Get(ctx, "probe:west").Bytes()
if errors.Is(err, redis.Nil) { /* المفتاح غايب */ }
~~~
~Get~ أمر قراية، و ~Bytes~ بترجع البايتات والخطأ. ~redis.Nil~ معناها غياب مؤكد؛ أي خطأ تاني ممكن يبقى عطل اتصال. افحص ~ctx.Err()~ قبل أي شغل وبعد قراية الكاش عشان إلغاء الطلب ميتعاملش كعطل يسمح بالـ fallback.

~Set~ بتكتب القيمة، وآخر وسيط مدة صلاحية TTL. صفر مدة معناه مفيش انتهاء، فهنا المنادي دايماً بيوفر مدة موجبة. ~Err~ بترجع خطأ الأمر. ~Del~ بتحذف مفتاح؛ حذف مفتاح غايب نجاح.
~~~go
err := client.Set(ctx, "probe:west", []byte("ready"), 9*time.Second).Err()
err = client.Del(ctx, "probe:west").Err()
~~~
Redis هي اللي بتملك الساعة. التست بيبدأ سيرفر من ~github.com/alicebob/miniredis/v2~ بـ ~srv := miniredis.NewMiniRedis()~. ~conn, err := srv.Dial()~ بتبدأ السيرفر وترجع اتصال في الذاكرة من غير بورت؛ المنادي بيقفل الاتصال بـ ~conn.Close()~. ~srv.FastForward(9*time.Second)~ بتنهي الصلاحية من غير نوم، و ~srv.Close()~ بتحاكي العطل. التست بيوقف إعادة الاتصال وقت العطل، وبيرجع يشغله عند اختبار التعافي. التنفيذ مبيقراش ساعة الجهاز.

## الترميز والتحديث
~json.Marshal(Beacon{ID: "west", Charge: 2})~ بترجع JSON بالحقول المصدّرة ~ID~ و ~Charge~. ~json.Unmarshal(raw, &b)~ بتفكها؛ JSON بايظ يعتبر miss، ومتخزّنش نتيجة فاشلة من المخزن.

احفظ في المخزن الأول، وبعد النجاح احذف المفتاح القديم. لو المخزن فشل متلمسش الكاش. لو الحفظ نجح والحذف فشل، رجّع تحذير يطابق ~ErrCacheSync~؛ الحفظ حصل بالفعل، فإعادة المحاولة العمياء مش ضمان ذرّية.
~~~go
return fmt.Errorf("%w: %w", ErrCacheSync, cacheErr)
~~~
~ErrCacheSync~ وسبب الكاش يفضلوا قابلين للكشف بـ ~errors.Is~. علّم المفتاح في خريطة محلية كـ ~dirty~ بعد فشل الحذف. القراية اللاحقة تتجاهل نسخته القديمة، تقرا المخزن، وتحاول الحذف تاني؛ امسح العلامة واملأ الكاش بس بعد نجاح الحذف. الخريطة دي ضمان داخل نفس الكائن فقط، وممكن تكبر وقت عطل طويل؛ مش معاملة موزعة ولا ضمان بين عمليات نظام مختلفة.

## التزامن
قفل ~sync.Mutex~ حوالين القراية والكتابة كلها يمنع قراية قديمة تملأ الكاش بعد تحديث أحدث داخل نفس الكائن. ده اختيار بسيط بيسلسل الشغل، حتى لمفاتيح مختلفة. مثال: ~c.mu.Lock(); defer c.mu.Unlock()~. المخزن والعميل محقونين ومش ~nil~، والـ TTL موجبة؛ مفيش إعدادات جهاز جوه المنطق.

## أخطاء شائعة
~~~go
client.Del(ctx, key) // غلط: الحذف قبل الحفظ ممكن يخبي فشل المصدر
store.Save(ctx, b)
if err != nil { return Beacon{}, nil } // غلط: عطل المصدر اتحوّل لنجاح
~~~

## المهمة
نفّذ ~NewBeaconCache(client, store, ttl) *BeaconCache~ و ~Read(ctx, id) (Beacon, error)~ و ~Write(ctx, b) error~ بتوقيعات الـ starter، وسيب ~main~ فاضية. المفتاح بالظبط ~colony:beacon:~ زائد الـ ~ID~ زي ما هو، من غير trim. الـ ~ID~ الفاضي أو شحنة سالبة في الكتابة يرجعوا ~ErrBeacon~ قبل أي I/O؛ القراية بترفض ~ID~ فاضي. افحص إلغاء السياق قبل I/O ومرّر نفس السياق لكل نداء.

الـ hit يرجع JSON الصالح من غير المصدر، والـ miss يقرا المصدر ويملا بـ TTL المحقونة. JSON بايظ يعتبر miss ويتستبدل بعد نجاح المصدر. خطأ ~Get~ غير ~redis.Nil~ يرجع نتيجة المصدر من غير محاولة ملء، إلا لو سياق الطلب اتلغى فارجع خطأ السياق. فشل ~Set~ بعد قراية مصدر ناجحة بيتجاهل، وأخطاء المصدر بتترجع من غير تخزين نتيجة فاشلة. عند فشل القراية ارجع ~Beacon{}~ والخطأ؛ المخزن بيرجع قيمة صفرية مع خطأ.

~Write~ تحفظ المصدر الأول وتلغي المفتاح بعد النجاح. فشل الحذف يرجع خطأ يطابق ~ErrCacheSync~ وسبب العطل، ويسجل ~dirty~. قراية ~dirty~ تقرا المصدر الأول؛ لو فشل ارجع خطأ من غير Redis، ولو نجح جرّب الحذف: فشله يسيب العلامة ويرجع القيمة من غير ملء، ونجاحه يمسح العلامة ويحاول ملء. نجاح كتابة لاحقة وحذفها يمسح العلامة كمان. كل العمليات آمنة بالتزامن ومفيش قراية قديمة تفضل بعد كتابة مكتملة على نفس الكائن.

النجمة التالتة: اعتبر JSON سليم نحوياً لكن ~ID~ مختلف أو ~Charge~ سالبة فساد، واتعامل معاه زي JSON البايظ. مفيش ~Report~ أو خدمات خارج التست.
`,
		Hints: []string{"الكاش نسخة مؤقتة، والمخزن يفضل المرجع لما القيمة تغيب.", "قسّم قراءة ~Read~ بين hit وmiss، وخلي كتابة ~Write~ تبدأ بالمصدر.", "فرّق بين ~redis.Nil~ وعطل Redis قبل اختيار مسار القراءة."},
		En: &LevelText{LessonRevision: 1, Steps: ch30Lesson301English, Title: "Cache Aside", Goal: "Accelerate beacon reads with Redis, expiry, and safe invalidation", Files: []LabFile{{Name: "main.go", Starter: ch301Starter}},
			Lesson: `# Beacons with a cache beside the store

## Before you write
- When the cache fails, which component remains authoritative?
- What should a caller hear if saving succeeds but cache deletion fails?

## Cache-aside
A cache is a temporary read accelerator, not the source of truth. In cache-aside, consult the cache first. A hit returns without accessing the store. A miss loads the store and attempts to fill the cache. Failed filling must not turn a successful source read into failure.

~BeaconStore~ is the injected store interface. For example, ~b, err := store.Load(ctx, "west")~ reads and ~err = store.Save(ctx, Beacon{ID: "west", Charge: 0})~ saves. ~Beacon~ has a string ~ID~ and a nonnegative ~Charge~; zero is valid.

## Real Redis commands
Import ~github.com/redis/go-redis/v9~. ~redis.NewClient~ builds a client from ~redis.Options~; ~Addr~ names the server. The caller owns the client and its ~Close~. This mission receives an existing ~*redis.Client~ and never starts a server itself.
~~~go
client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
defer client.Close()
raw, err := client.Get(ctx, "probe:west").Bytes()
if errors.Is(err, redis.Nil) { /* confirmed absence */ }
~~~
~Get~ reads; ~Bytes~ returns bytes and an error. ~redis.Nil~ means confirmed absence; other errors may mean transport failure. Check ~ctx.Err()~ before work and after the cache read so cancellation is never mistaken for an outage permitting fallback.

~Set~ writes, with a final time-to-live (TTL) argument. Zero means no expiry; the caller supplies a positive duration here. ~Err~ returns the command error. ~Del~ deletes a key; deleting an absent key succeeds.
~~~go
err := client.Set(ctx, "probe:west", []byte("ready"), 9*time.Second).Err()
err = client.Del(ctx, "probe:west").Err()
~~~
Redis owns expiry time. Tests create a server from ~github.com/alicebob/miniredis/v2~ with ~srv := miniredis.NewMiniRedis()~. ~conn, err := srv.Dial()~ starts it and returns an in-memory connection without a port; the caller releases the connection with ~conn.Close()~. ~srv.FastForward(9*time.Second)~ expires entries without sleeping; ~srv.Close()~ simulates outage. The fixture prevents reconnection during outage and restarts for recovery checks. Your implementation never reads the machine clock.

## Encoding and writes
~json.Marshal(Beacon{ID: "west", Charge: 2})~ encodes exported ~ID~ and ~Charge~ fields. ~json.Unmarshal(raw, &b)~ decodes; malformed JSON is a miss. Never cache failed store results.

Save to the source first, then delete the old cache key. Source failure leaves the cache untouched. When saving succeeds but deletion fails, return a warning matching ~ErrCacheSync~: the source has already committed, so blindly retrying does not provide atomicity.
~~~go
return fmt.Errorf("%w: %w", ErrCacheSync, cacheErr)
~~~
Both ~ErrCacheSync~ and the cache cause remain discoverable with ~errors.Is~. Mark that key dirty in a local map after failed deletion. Later reads ignore its old cached value, load the source, and retry deletion. Clear the marker and refill only after deletion succeeds. This protects one object only; the map can grow during a prolonged outage. It is neither a distributed transaction nor a cross-process guarantee.

## Concurrency
A ~sync.Mutex~ around each entire read/write prevents an old read from refilling after a newer update within this object. This deliberately simple approach serializes even different keys. Example: ~c.mu.Lock(); defer c.mu.Unlock()~. The injected store and client are non-nil and TTL positive; machine configuration does not belong inside this logic.

## Common mistakes
~~~go
client.Del(ctx, key) // Wrong: invalidating before saving hides source failure
store.Save(ctx, b)
if err != nil { return Beacon{}, nil } // Wrong: source failure became success
~~~

## Your mission
Implement ~NewBeaconCache(client, store, ttl) *BeaconCache~, ~Read(ctx, id) (Beacon, error)~, and ~Write(ctx, b) error~ with the starter signatures; leave ~main~ empty. Keys are exactly ~colony:beacon:~ plus the unchanged ID, without trimming. Empty ID or negative write charge returns ~ErrBeacon~ before I/O; reads reject empty IDs. Check cancellation before I/O and pass the original context to every call.

Hits return valid JSON without the source. Misses load the source and fill with the injected TTL. Malformed JSON is a miss and is replaced after source success. ~Get~ errors other than ~redis.Nil~ return the source result without filling, except a canceled caller context returns its context error. Failed ~Set~ after source success is ignored. Source errors propagate without caching failed results. A failed read returns ~Beacon{}~ and the error; the store itself returns a zero value on errors.

~Write~ saves the source first and invalidates after success. Failed deletion returns an error matching ~ErrCacheSync~ and its cause, and marks dirty. Dirty reads load the source first; source failure returns without Redis calls. On success retry deletion: failure keeps dirty and returns the source value without filling, while success clears dirty and attempts filling. A later successful write/deletion also clears dirty. All operations are concurrent-safe, and no stale read may survive a completed write through the same object.

Third star: syntactically valid JSON with mismatched ~ID~ or negative ~Charge~ is also corrupt and follows the malformed-JSON path. No ~Report~ or services outside the test.
`, Hints: []string{"A cache is temporary; the store remains authoritative when a value is missing.", "Separate ~Read~ hit and miss paths, and let ~Write~ begin with the store.", "Distinguish ~redis.Nil~ from a Redis outage before choosing a read path."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch301Hidden},
	})
}

const ch301Starter = `package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type Beacon struct {
	ID     string
	Charge int
}
type BeaconStore interface {
	Load(context.Context, string) (Beacon, error)
	Save(context.Context, Beacon) error
}

var ErrBeacon = errors.New("invalid beacon")
var ErrCacheSync = errors.New("saved, cache invalidation failed")

type BeaconCache struct{}

func NewBeaconCache(client *redis.Client, store BeaconStore, ttl time.Duration) *BeaconCache {
	return &BeaconCache{}
}
func (c *BeaconCache) Read(ctx context.Context, id string) (Beacon, error) {
	return Beacon{}, errors.New("TODO")
}
func (c *BeaconCache) Write(ctx context.Context, b Beacon) error { return errors.New("TODO") }
func main()                                                      {}
`
const ch301Hidden = `package main

import (
	"context"
	"encoding/json"
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

type hiddenBeacons struct {
	mu     sync.Mutex
	rows   map[string]Beacon
	loads  int
	fault  error
	seen   context.Context
	onLoad func()
	onSave func()
}

func (s *hiddenBeacons) Load(ctx context.Context, id string) (Beacon, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loads++
	s.seen = ctx
	if s.onLoad != nil {
		s.onLoad()
	}
	if s.fault != nil {
		return Beacon{}, s.fault
	}
	b, ok := s.rows[id]
	if !ok {
		return Beacon{}, ErrBeacon
	}
	return b, nil
}
func (s *hiddenBeacons) Save(ctx context.Context, b Beacon) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen = ctx
	if s.fault != nil {
		return s.fault
	}
	s.rows[b.ID] = b
	if s.onSave != nil {
		s.onSave()
	}
	return nil
}

// hiddenRedisServer guards Dial's automatic startup so Close really simulates an outage.
type hiddenRedisServer struct {
	*miniredis.Miniredis
	mu     sync.Mutex
	online bool
}

func (s *hiddenRedisServer) dial(ctx context.Context, _, _ string) (net.Conn, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if !s.online {
		return nil, net.ErrClosed
	}
	return s.Miniredis.Dial()
}
func (s *hiddenRedisServer) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.online = false
	s.Miniredis.Close()
}
func (s *hiddenRedisServer) Restart() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	conn, err := s.Miniredis.Dial()
	if err != nil {
		return err
	}
	conn.Close()
	s.online = true
	return nil
}
func hiddenCache(t *testing.T) (*hiddenRedisServer, *redis.Client, *hiddenBeacons, *BeaconCache) {
	t.Helper()
	srv := &hiddenRedisServer{Miniredis: miniredis.NewMiniRedis()}
	if err := srv.Restart(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(srv.Close)
	client := redis.NewClient(&redis.Options{Dialer: srv.dial, MaxRetries: -1, ReadTimeout: time.Second, WriteTimeout: time.Second})
	t.Cleanup(func() { client.Close() })
	store := &hiddenBeacons{rows: make(map[string]Beacon)}
	return srv, client, store, NewBeaconCache(client, store, 7*time.Second)
}
func TestHiddenBeaconHitMissExpiryAndInvalidation(t *testing.T) {
	srv, client, store, c := hiddenCache(t)
	ctx := context.WithValue(context.Background(), struct{}{}, "caller")
	rng := rand.New(rand.NewSource(3001))
	for i := 0; i < 20; i++ {
		b := Beacon{fmt.Sprintf("灯-%d-%d", i, rng.Intn(999)), rng.Intn(1000)}
		store.rows[b.ID] = b
		got, e := c.Read(ctx, b.ID)
		if e != nil || got != b {
			t.Fatal("miss", got, e)
		}
		if store.seen != ctx {
			t.Fatal("caller context lost")
		}
		key := "colony:beacon:" + b.ID
		if srv.TTL(key) != 7*time.Second {
			t.Fatal("TTL", srv.TTL(key))
		}
		before := store.loads
		if got, e = c.Read(ctx, b.ID); e != nil || got != b || store.loads != before {
			t.Fatal("hit reached store", got, e)
		}
		srv.FastForward(7 * time.Second)
		if _, e = c.Read(ctx, b.ID); e != nil || store.loads != before+1 {
			t.Fatal("expiry", e)
		}
		b.Charge = 0
		if e = c.Write(ctx, b); e != nil || srv.Exists(key) {
			t.Fatal("invalidation", e)
		}
		if got, e = c.Read(ctx, b.ID); e != nil || got != b {
			t.Fatal("updated read", got, e)
		}
		if e = client.Set(ctx, key, "{broken", time.Minute).Err(); e != nil {
			t.Fatal(e)
		}
		if got, e = c.Read(ctx, b.ID); e != nil || got != b {
			t.Fatal("corrupt fallback", got, e)
		}
	}
}
func TestHiddenBeaconFailuresAndRecovery(t *testing.T) {
	srv, _, store, c := hiddenCache(t)
	ctx := context.Background()
	b := Beacon{"east", 4}
	store.rows[b.ID] = b
	if _, e := c.Read(ctx, b.ID); e != nil {
		t.Fatal(e)
	}
	sentinel := errors.New("source offline")
	store.fault = sentinel
	if e := c.Write(ctx, Beacon{"east", 8}); !errors.Is(e, sentinel) {
		t.Fatal("source error", e)
	}
	if got, e := c.Read(ctx, b.ID); e != nil || got != b {
		t.Fatal("failed write invalidated hit", got, e)
	}
	if _, e := c.Read(ctx, "absent"); !errors.Is(e, sentinel) || srv.Exists("colony:beacon:absent") {
		t.Fatal("failed miss cached", e)
	}
	store.fault = nil
	srv.Close()
	b.Charge = 9
	if e := c.Write(ctx, b); !errors.Is(e, ErrCacheSync) || store.rows[b.ID] != b {
		t.Fatal("committed write warning", e)
	}
	if got, e := c.Read(ctx, b.ID); e != nil || got != b {
		t.Fatal("outage fallback", got, e)
	}
	// Restart preserves the old cache entry: dirty keys must bypass it until deletion succeeds.
	if e := srv.Restart(); e != nil {
		t.Fatal(e)
	}
	if got, e := c.Read(ctx, b.ID); e != nil || got != b {
		t.Fatal("stale cache after recovery", got, e)
	}
	before := store.loads
	if got, e := c.Read(ctx, b.ID); e != nil || got != b || store.loads != before {
		t.Fatal("recovered hit", got, e)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	before = store.loads
	if _, e := c.Read(canceled, b.ID); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if e := c.Write(canceled, b); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if store.loads != before {
		t.Fatal("canceled read accessed source")
	}
	for _, bad := range []Beacon{{"", 0}, {"east", -1}} {
		if e := c.Write(ctx, bad); !errors.Is(e, ErrBeacon) {
			t.Fatal("validation", e)
		}
	}
	if _, e := c.Read(ctx, ""); !errors.Is(e, ErrBeacon) {
		t.Fatal(e)
	}
}
func TestHiddenBeaconConcurrentAccess(t *testing.T) {
	_, _, store, c := hiddenCache(t)
	ctx := context.Background()
	store.rows["shared"] = Beacon{"shared", 0}
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprint(i)
			if e := c.Write(ctx, Beacon{id, i}); e != nil {
				t.Error(e)
			}
			if b, e := c.Read(ctx, id); e != nil || b.Charge != i {
				t.Error("isolated", b, e)
			}
			if e := c.Write(ctx, Beacon{"shared", i}); e != nil {
				t.Error(e)
			}
			if _, e := c.Read(ctx, "shared"); e != nil {
				t.Error(e)
			}
		}(i)
	}
	wg.Wait()
	if e := c.Write(ctx, Beacon{"shared", 77}); e != nil {
		t.Fatal(e)
	}
	if b, e := c.Read(ctx, "shared"); e != nil || b.Charge != 77 {
		t.Fatal("final state", b, e)
	}
}
func TestHiddenBeaconValidationAndSourceFailureBeforeRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1", MaxRetries: -1})
	defer client.Close()
	store := &hiddenBeacons{rows: map[string]Beacon{}, fault: errors.New("source rejected")}
	c := NewBeaconCache(client, store, time.Second)
	ctx := context.Background()
	for _, b := range []Beacon{{"", 0}, {"a", -1}} {
		if e := c.Write(ctx, b); !errors.Is(e, ErrBeacon) {
			t.Fatal(e)
		}
	}
	if _, e := c.Read(ctx, ""); !errors.Is(e, ErrBeacon) {
		t.Fatal(e)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := c.Read(canceled, "a"); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if e := c.Write(canceled, Beacon{"a", 0}); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if store.seen != nil {
		t.Fatal("invalid/canceled input touched source")
	}
	if e := c.Write(ctx, Beacon{"a", 0}); !errors.Is(e, store.fault) {
		t.Fatal("failed source write reached cache", e)
	}
}
func TestHiddenBeaconBestEffortFillAndCommittedCancellation(t *testing.T) {
	srv, _, store, c := hiddenCache(t)
	ctx := context.Background()
	b := Beacon{"faint", 0}
	store.rows[b.ID] = b
	store.onLoad = func() { srv.SetError("ERR cache full") }
	if got, e := c.Read(ctx, b.ID); e != nil || got != b {
		t.Fatal("fill failure replaced source success", got, e)
	}
	store.onLoad = nil
	srv.SetError("")
	if srv.Exists("colony:beacon:" + b.ID) {
		t.Fatal("failed fill cached")
	}
	if _, e := c.Read(ctx, b.ID); e != nil {
		t.Fatal(e)
	}
	canceled, cancel := context.WithCancel(ctx)
	defer cancel()
	store.onSave = cancel
	b.Charge = 5
	if e := c.Write(canceled, b); !errors.Is(e, ErrCacheSync) || !errors.Is(e, context.Canceled) {
		t.Fatal("committed cancellation", e)
	}
	store.onSave = nil
	if got, e := c.Read(ctx, b.ID); e != nil || got != b {
		t.Fatal("bypass canceled invalidation", got, e)
	}
}
func TestHiddenBonusBeaconSemanticCorruption(t *testing.T) {
	_, client, store, c := hiddenCache(t)
	ctx := context.Background()
	want := Beacon{"north", 3}
	store.rows[want.ID] = want
	for _, b := range []Beacon{{"other", 4}, {"north", -1}, {"", 0}} {
		raw, _ := json.Marshal(b)
		if e := client.Set(ctx, "colony:beacon:north", raw, time.Minute).Err(); e != nil {
			t.Fatal(e)
		}
		if got, e := c.Read(ctx, "north"); e != nil || got != want {
			t.Fatal("semantic corruption", got, e)
		}
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch30Lesson302Arabic, ID: "30-2-documents", Chapter: 30, Kind: "lab", Title: "المستندات", Goal: "احفظ عينات الصخور كمستندات بفهرس فريد وصفحات ومهلة", Setup: seededLabRoom,
		Requires: []string{"iface", "sel:mongo.Collection"}, Deps: []string{"go.mongodb.org/mongo-driver/v2"}, Bridge: []int{67}, Race: true, Service: "mongo",
		Files: []LabFile{{Name: "main.go", Starter: ch302Starter}},
		Lesson: `# أرشيف عينات على MongoDB

## قبل ما تكتب
- إزاي القاعدة ترفض نفس علامة العينة حتى لو طلبين وصلوا مع بعض؟
- ليه رقم آخر عينة أوضح من عدد الصفوف اللي اتخطّيناها؟

## المستند وBSON
MongoDB بتخزن مستندات جوه collection، والـ collections جوه database. المستند مجموعة حقول بأنواع، مش نص JSON محفوظ وخلاص. BSON هو ترميز MongoDB الثنائي؛ بيمثل أعداد ونصوص وأنواع تانية. الحقل ~_id~ مفتاح فريد افتراضياً. هنا بنختار ~int64~ محقونة عشان مفيش توليد عشوائي أو ساعة في المنطق.

استورد ~go.mongodb.org/mongo-driver/v2/mongo~ و ~go.mongodb.org/mongo-driver/v2/bson~ و ~go.mongodb.org/mongo-driver/v2/mongo/options~. ~Sample~ بتمثل ~ID~ و ~Tag~ و ~Mass~. الـ tag المسمى ~bson~ بيربط اسم Go بالحقل المحفوظ.
~~~go
type Probe struct {
 Code int64 ` + "`" + `bson:"_id"` + "`" + `
 Site string ` + "`" + `bson:"site"` + "`" + `
}
filter := bson.D{{Key: "_id", Value: int64(8)}}
~~~
~bson.D~ قائمة حقول مرتبة؛ كل عنصر فيه ~Key~ و ~Value~. استخدمها للفلاتر والترتيب. ~Sample{ID: 8, Tag: "basalt", Mass: 0}~ عينة صالحة؛ الصفر مسموح للكتلة.

## الاتصال ملك المنادي
~mongo.Connect~ بتبني عميل، و ~options.Client().ApplyURI~ بتحدد الاتصال. ~Database~ و ~Collection~ بيختاروا النطاق من غير إنشاء اتصال لكل عملية. ~Ping~ بيختبر الاتصال، و ~Disconnect~ بيقفل العميل. التست بس بيقرا ~ORBITAL_MONGO_URI~؛ التنفيذ بياخد collection جاهزة.
~~~go
client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("ORBITAL_MONGO_URI")))
if err != nil { return err }
defer client.Disconnect(ctx)
if err := client.Ping(ctx, nil); err != nil { return err }
collection := client.Database("owned_test_name").Collection("probes")
~~~
اختبارات التكامل بتتخطى لو المتغير غايب، ولو موجود بتستخدم قاعدة باسم فريد وبتحذف نفس القاعدة بعد التست. تستات الوحدة بتشتغل دايماً على ~SampleStore~ وهمية؛ ده يثبت منطق الأرشيف، مش صحة استعلام MongoDB.

## فهرس فريد
~Indexes().CreateOne~ بتنشئ ~mongo.IndexModel~. ~Keys~ الحقول واتجاهاتها؛ 1 تصاعدي. ~options.Index().SetUnique(true).SetName("unique_probe_site")~ بتسمّي الفهرس وتفرض الفرادة. فحص قبل الكتابة لوحده فيه سباق؛ الفهرس هو الحكم.
~~~go
_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
 Keys: bson.D{{Key: "site", Value: 1}},
 Options: options.Index().SetUnique(true).SetName("unique_probe_site"),
})
~~~
~mongo.IsDuplicateKeyError(err)~ بتكشف مخالفة الفرادة من نوع الخطأ الرسمي. ~mongo.ErrNoDocuments~ بيمثل غياب نتيجة واحدة؛ متفحصش نص رسالة السيرفر.
~~~go
if mongo.IsDuplicateKeyError(err) { return ErrSampleDuplicate }
if errors.Is(err, mongo.ErrNoDocuments) { return ErrSampleMissing }
~~~

## إضافة وقراية واستبدال وحذف
~InsertOne~ بتضيف struct باستخدام BSON tags. ~FindOne~ بتاخد filter، و ~Decode~ بتملا struct أو ترجع خطأ، ومنه الغياب. ~ReplaceOne~ بتستبدل المستند الموجود من غير upsert افتراضياً؛ ~MatchedCount~ صفر معناها غياب حتى لو ~ModifiedCount~ صفر ممكن يبقى نفس المحتوى. ~DeleteOne~ بترجع ~DeletedCount~.
~~~go
_, err := collection.InsertOne(ctx, Probe{8, "west"})
var probe Probe
err = collection.FindOne(ctx, filter).Decode(&probe)
changed, err := collection.ReplaceOne(ctx, filter, Probe{8, "east"})
removed, err := collection.DeleteOne(ctx, filter)
_ = changed.MatchedCount
_ = removed.DeletedCount
~~~

## مؤشر قاعدة ومؤشر صفحة
~Find~ بترجع ~cursor~ موارد للسيرفر. ~Next(ctx)~ تتقدم، و ~Decode~ تقرا الحالي، و ~Err~ تكشف فشل اللف، و ~Close(ctx)~ تحرر الموارد حتى عند فشل الفك. الفلتر ~$gt~ معناه أكبر من، مش أكبر أو يساوي. ~SetSort~ و ~SetLimit~ بيحددوا ترتيب وعدد المستندات في السيرفر.
~~~go
cursor, err := collection.Find(ctx,
 bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: int64(8)}}}},
 options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(4))
if err != nil { return err }
defer cursor.Close(ctx)
for cursor.Next(ctx) {
 var probe Probe
 if err := cursor.Decode(&probe); err != nil { return err }
}
return cursor.Err()
~~~
مؤشر الصفحة عندنا رقم آخر ~ID~ رجّعناه، مش الـ ~cursor~ اللي ماسك موارد القاعدة. اطلب ~limit+1~ عشان تعرف إن فيه صفحة تانية. رجّع أول ~limit~ بس، و ~next~ آخر ~ID~ فيهم لو الإضافي موجود. لو مفيش إضافي رجّع صفر. الـ IDs فريدة فمفيش تعادل. الضمان لبيانات ثابتة؛ ده مش snapshot بين نداءات منفصلة.

## سياق محدود من الأب
~context.WithTimeout(ctx, timeout)~ بيحتفظ بقيم وإلغاء الأب، وبيختار المهلة الأقرب لو الأب أقصر. نادِ ~cancel~ دايماً لتحرير المؤقت. ~NewArchive(store, 2*time.Second)~ بياخد مدة موجبة ومخزن مش ~nil~. كل نداء مخزن من الأرشيف لازم ياخد سياق مشتق محدود، ومش ~context.Background()~.
~~~go
bounded, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()
return store.Insert(bounded, sample)
~~~

## أخطاء شائعة
~~~go
collection.Find(context.Background(), bson.D{}) // سياق المنادي ضاع
if changed.ModifiedCount == 0 { return ErrSampleMissing } // ممكن العنصر موجود ومفيش تغيير
for cursor.Next(ctx) { } // من غير Close ولا Err
~~~

## المهمة
احتفظ بتوقيعات الـ starter وسيب ~main~ فاضية. ~SampleStore~ فيها ~Insert~ و ~Get~ و ~Replace~ و ~Remove~ و ~After~. الأخيرة بترجع بحد أقصى العدد المطلوب، بترتيب ~ID~ تصاعدي، وأكبر من ~after~ حصرياً. الأخطاء بترجع قيمة صفرية أو slice بـ ~nil~. المنادي بيحقن المخزن والمهلة الموجبة في ~NewArchive~. ~Archive~ مفيهاش حالة متغيرة وبتدعم التزامن لو المخزن بيدعمه.

~Add(ctx, sample)~ و ~Update(ctx, sample)~ يطبعوا ~Tag~ بـ ~strings.ToLower(strings.TrimSpace(s.Tag))~ ويرفضوا ~ID <= 0~ أو علامة فاضية أو ~Mass < 0~ بـ ~ErrSample~ قبل المخزن. الإضافة تنادي ~Insert~ مرة، والتعديل ~Replace~ مرة. ~Get(ctx, id)~ و ~Delete(ctx, id)~ يرفضوا ~ID~ غير موجب بنفس الخطأ وينادوا ~Get~ و ~Remove~. كل دالة تفحص إلغاء الأب الأول، وتستخدم سياق مشتق بمهلة الأرشيف وتعمل ~cancel~ عند الرجوع، وتحافظ على خطأ المخزن.

~Page(ctx, after, limit)~ ترفض ~after < 0~ أو ~limit~ برا 1 لـ 50 بـ ~ErrSample~ قبل المخزن. نادي ~After~ مرة بـ ~limit+1~ وبسياق محدود. رجّع ~([]Sample, next int64, error)~؛ القائمة الفاضية مش ~nil~، والفشل ~nil, 0, err~، ومؤشر التالي زي الشرح. مثال: ~rows, next, err := archive.Page(ctx, 0, 2)~ ثم ~archive.Page(ctx, next, 2)~ لو ~next~ مش صفر.

~NewMongoSamples(collection)~ بتغلّف collection محقونة. ~Init(ctx)~ تنشئ فهرس تصاعدي فريد اسمه ~unique_sample_tag~ على ~tag~، بشكل قابل للتكرار. ~Insert~ و ~Get~ و ~Replace~ و ~Remove~ و ~After~ تنفذ العقد الحقيقي في MongoDB على حقول ~_id~ و ~tag~ و ~mass~ بالـ tags الموجودة. الاستبدال ميعملش upsert، والكتلة صفر تتكتب. ~After~ تعمل فلتر حصري وترتيب وحد على السيرفر وتقفل ~cursor~ وتفحص خطأها النهائي؛ فشل الفك يرجع ~nil,error~. الـ adapter ميطبعش ولا يتحقق تاني: المدخل جاي من الأرشيف صالح. مرّر السياق، وحوّل الفرادة لـ ~ErrSampleDuplicate~ والغياب لـ ~ErrSampleMissing~، وحافظ على باقي الأخطاء. ~Init~ والـ adapter بياخدوا سياق محدود من المنادي. مفيش ~Report~، والاتصال والتنضيف في التست بس.
`,
		Hints: []string{"المستند record من حقول، والفهرس الفريد يحمي القاعدة وقت السباق.", "افصل ~Archive~ عن adapter Mongo، وخلي المؤشر يستخدم IDs مرتبة.", "راجع ~MatchedCount~ و~cursor.Err~ عند التعديل والصفحات."},
		En: &LevelText{LessonRevision: 1, Steps: ch30Lesson302English, Title: "Documents", Goal: "Store rock samples as documents with a unique index, pagination, and a timeout", Files: []LabFile{{Name: "main.go", Starter: ch302Starter}},
			Lesson: `# A MongoDB sample archive

## Before you write
- How can the database reject a shared sample tag when two requests race?
- Why use the last sample ID rather than the number of skipped records?

## Documents and BSON
MongoDB stores documents inside collections inside databases. A document contains typed fields, rather than merely storing JSON text. BSON is MongoDB's binary encoding, supporting integers, strings, and other types. ~_id~ is unique by default. We choose injected ~int64~ IDs so domain logic needs neither random generation nor a clock.

Import ~go.mongodb.org/mongo-driver/v2/mongo~, ~go.mongodb.org/mongo-driver/v2/bson~, and ~go.mongodb.org/mongo-driver/v2/mongo/options~. ~Sample~ holds ~ID~, ~Tag~, and ~Mass~. The ~bson~ struct tag maps a Go field to its stored name.
~~~go
type Probe struct {
 Code int64 ` + "`" + `bson:"_id"` + "`" + `
 Site string ` + "`" + `bson:"site"` + "`" + `
}
filter := bson.D{{Key: "_id", Value: int64(8)}}
~~~
~bson.D~ is an ordered field list; each element has ~Key~ and ~Value~. Use it for filters and sort order. ~Sample{ID: 8, Tag: "basalt", Mass: 0}~ is valid; zero mass is allowed.

## The caller owns connections
~mongo.Connect~ creates a client, and ~options.Client().ApplyURI~ supplies connection configuration. ~Database~ and ~Collection~ select scope without creating a connection per operation. ~Ping~ tests connectivity; ~Disconnect~ closes the client. Only tests read ~ORBITAL_MONGO_URI~; your implementation receives a collection.
~~~go
client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("ORBITAL_MONGO_URI")))
if err != nil { return err }
defer client.Disconnect(ctx)
if err := client.Ping(ctx, nil); err != nil { return err }
collection := client.Database("owned_test_name").Collection("probes")
~~~
Integration skips when the variable is absent. When supplied, each run uses a uniquely named database and drops only that database. Unit tests always run against a fake ~SampleStore~: that proves archive logic, not MongoDB query behavior.

## A unique index
~Indexes().CreateOne~ creates a ~mongo.IndexModel~. ~Keys~ specifies fields and direction; 1 means ascending. ~options.Index().SetUnique(true).SetName("unique_probe_site")~ names the index and enforces uniqueness. Checking before inserting still races; the index decides.
~~~go
_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
 Keys: bson.D{{Key: "site", Value: 1}},
 Options: options.Index().SetUnique(true).SetName("unique_probe_site"),
})
~~~
~mongo.IsDuplicateKeyError(err)~ inspects official driver errors for uniqueness failures. ~mongo.ErrNoDocuments~ represents single-result absence. Never classify by server message text.
~~~go
if mongo.IsDuplicateKeyError(err) { return ErrSampleDuplicate }
if errors.Is(err, mongo.ErrNoDocuments) { return ErrSampleMissing }
~~~

## Insert, read, replace, delete
~InsertOne~ stores a struct through BSON tags. ~FindOne~ takes a filter; ~Decode~ fills a struct or returns an error, including absence. ~ReplaceOne~ replaces an existing document without upsert by default. Zero ~MatchedCount~ means absent; zero ~ModifiedCount~ can merely mean unchanged content. ~DeleteOne~ reports ~DeletedCount~.
~~~go
_, err := collection.InsertOne(ctx, Probe{8, "west"})
var probe Probe
err = collection.FindOne(ctx, filter).Decode(&probe)
changed, err := collection.ReplaceOne(ctx, filter, Probe{8, "east"})
removed, err := collection.DeleteOne(ctx, filter)
_ = changed.MatchedCount
_ = removed.DeletedCount
~~~

## Database cursors and page cursors
~Find~ returns a server-resource cursor. ~Next(ctx)~ advances; ~Decode~ reads the current document; ~Err~ checks terminal failure; ~Close(ctx)~ releases resources even on decoding failure. ~$gt~ means strictly greater, not greater-or-equal. ~SetSort~ and ~SetLimit~ apply ordering and a document cap on the server.
~~~go
cursor, err := collection.Find(ctx,
 bson.D{{Key: "_id", Value: bson.D{{Key: "$gt", Value: int64(8)}}}},
 options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}).SetLimit(4))
if err != nil { return err }
defer cursor.Close(ctx)
for cursor.Next(ctx) {
 var probe Probe
 if err := cursor.Decode(&probe); err != nil { return err }
}
return cursor.Err()
~~~
Our page cursor is the last returned ID, not a resource-owning database cursor. Fetch ~limit+1~ to discover whether another page exists. Return at most ~limit~, with ~next~ set to the last returned ID only when an extra row exists; otherwise zero. IDs are unique, so there are no ties. This guarantee assumes stable data, not a snapshot across independent calls.

## A bounded child context
~context.WithTimeout(ctx, timeout)~ preserves the parent's values and cancellation, choosing its earlier deadline when applicable. Always call ~cancel~ to release the timer. ~NewArchive(store, 2*time.Second)~ takes a positive timeout and non-nil store. Every archive store call receives a bounded derived context, never ~context.Background()~.
~~~go
bounded, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()
return store.Insert(bounded, sample)
~~~

## Common mistakes
~~~go
collection.Find(context.Background(), bson.D{}) // Caller cancellation lost
if changed.ModifiedCount == 0 { return ErrSampleMissing } // Existing unchanged record
for cursor.Next(ctx) { } // No Close or terminal Err check
~~~

## Your mission
Keep the starter signatures and empty ~main~. ~SampleStore~ provides ~Insert~, ~Get~, ~Replace~, ~Remove~, and ~After~. The last returns at most the requested count ordered by ascending ID, strictly greater than after. On error return a zero object or nil slice. The caller injects a non-nil store and positive timeout into ~NewArchive~. ~Archive~ has no mutable state and supports concurrency when its store does.

~Add(ctx, sample)~ and ~Update(ctx, sample)~ normalize ~Tag~ with ~strings.ToLower(strings.TrimSpace(s.Tag))~ and reject ~ID <= 0~, blank tag, or ~Mass < 0~ as ~ErrSample~ before calling the store. Add calls ~Insert~ once; update calls ~Replace~ once. ~Get(ctx, id)~ and ~Delete(ctx, id)~ reject nonpositive IDs with the same error and call ~Get~ and ~Remove~ respectively. Every method checks parent cancellation first, derives a context bounded by the archive timeout, cancels it on return, and preserves store errors.

~Page(ctx, after, limit)~ rejects ~after < 0~ or limits outside 1..50 as ~ErrSample~ before the store. Call ~After~ once with ~limit+1~ and a bounded context. Return ~([]Sample, next int64, error)~: non-nil empty slices, ~nil, 0, err~ on failure, and the continuation described above. Example: ~rows, next, err := archive.Page(ctx, 0, 2)~, followed by ~archive.Page(ctx, next, 2)~ when next is nonzero.

~NewMongoSamples(collection)~ wraps an injected collection. ~Init(ctx)~ idempotently creates an ascending unique index named ~unique_sample_tag~ on ~tag~. ~Insert~, ~Get~, ~Replace~, ~Remove~, and ~After~ implement the real MongoDB contract over ~_id~, ~tag~, and ~mass~, using the starter tags. Replacement does not upsert and must write zero mass. ~After~ filters exclusively, sorts and limits on the server, closes the cursor, and checks terminal error; decode failure returns ~nil,error~. The adapter does no further validation or normalization: archive inputs are valid. Forward context, map uniqueness to ~ErrSampleDuplicate~ and absence to ~ErrSampleMissing~, and preserve other errors. ~Init~ and adapter methods receive bounded contexts from their caller. No ~Report~; tests alone connect and clean up.
`, Hints: []string{"A document is a typed record, and a unique index protects concurrent writes.", "Separate ~Archive~ from the Mongo adapter, and page with ordered IDs.", "Check ~MatchedCount~ and ~cursor.Err~ during replacement and paging."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch302Hidden, HiddenFiles: map[string]string{"mongo_integration_test.go": ch302Integration}},
	})
}

const ch302Starter = `package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type Sample struct {
	ID   int64  ` + "`" + `bson:"_id"` + "`" + `
	Tag  string ` + "`" + `bson:"tag"` + "`" + `
	Mass int    ` + "`" + `bson:"mass"` + "`" + `
}

var ErrSample = errors.New("invalid sample")
var ErrSampleMissing = errors.New("sample missing")
var ErrSampleDuplicate = errors.New("duplicate sample")

type SampleStore interface {
	Insert(context.Context, Sample) error
	Get(context.Context, int64) (Sample, error)
	Replace(context.Context, Sample) error
	Remove(context.Context, int64) error
	After(context.Context, int64, int) ([]Sample, error)
}
type Archive struct{}

func NewArchive(store SampleStore, timeout time.Duration) *Archive { return &Archive{} }
func (a *Archive) Add(ctx context.Context, s Sample) error         { return errors.New("TODO") }
func (a *Archive) Update(ctx context.Context, s Sample) error      { return errors.New("TODO") }
func (a *Archive) Get(ctx context.Context, id int64) (Sample, error) {
	return Sample{}, errors.New("TODO")
}
func (a *Archive) Delete(ctx context.Context, id int64) error { return errors.New("TODO") }
func (a *Archive) Page(ctx context.Context, after int64, limit int) ([]Sample, int64, error) {
	return nil, 0, errors.New("TODO")
}

type MongoSamples struct{}

func NewMongoSamples(collection *mongo.Collection) *MongoSamples   { return &MongoSamples{} }
func (s *MongoSamples) Init(ctx context.Context) error             { return errors.New("TODO") }
func (s *MongoSamples) Insert(ctx context.Context, b Sample) error { return errors.New("TODO") }
func (s *MongoSamples) Get(ctx context.Context, id int64) (Sample, error) {
	return Sample{}, errors.New("TODO")
}
func (s *MongoSamples) Replace(ctx context.Context, b Sample) error { return errors.New("TODO") }
func (s *MongoSamples) Remove(ctx context.Context, id int64) error  { return errors.New("TODO") }
func (s *MongoSamples) After(ctx context.Context, after int64, limit int) ([]Sample, error) {
	return nil, errors.New("TODO")
}
func main() {}
`
const ch302Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

type hiddenSamples struct {
	rows  map[int64]Sample
	calls int
	ctx   context.Context
	fault error
	after int64
	limit int
}

func (s *hiddenSamples) record(ctx context.Context) error { s.calls++; s.ctx = ctx; return s.fault }
func (s *hiddenSamples) Insert(ctx context.Context, b Sample) error {
	if e := s.record(ctx); e != nil {
		return e
	}
	for _, old := range s.rows {
		if old.ID == b.ID || old.Tag == b.Tag {
			return ErrSampleDuplicate
		}
	}
	s.rows[b.ID] = b
	return nil
}
func (s *hiddenSamples) Get(ctx context.Context, id int64) (Sample, error) {
	if e := s.record(ctx); e != nil {
		return Sample{}, e
	}
	b, ok := s.rows[id]
	if !ok {
		return Sample{}, ErrSampleMissing
	}
	return b, nil
}
func (s *hiddenSamples) Replace(ctx context.Context, b Sample) error {
	if e := s.record(ctx); e != nil {
		return e
	}
	if _, ok := s.rows[b.ID]; !ok {
		return ErrSampleMissing
	}
	for _, old := range s.rows {
		if old.ID != b.ID && old.Tag == b.Tag {
			return ErrSampleDuplicate
		}
	}
	s.rows[b.ID] = b
	return nil
}
func (s *hiddenSamples) Remove(ctx context.Context, id int64) error {
	if e := s.record(ctx); e != nil {
		return e
	}
	if _, ok := s.rows[id]; !ok {
		return ErrSampleMissing
	}
	delete(s.rows, id)
	return nil
}
func (s *hiddenSamples) After(ctx context.Context, after int64, limit int) ([]Sample, error) {
	if e := s.record(ctx); e != nil {
		return nil, e
	}
	s.after = after
	s.limit = limit
	rows := []Sample{}
	for _, b := range s.rows {
		if b.ID > after {
			rows = append(rows, b)
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	if len(rows) > limit {
		rows = rows[:limit]
	}
	return rows, nil
}
func hiddenArchiveContext(t *testing.T, store *hiddenSamples, parent context.Context) {
	t.Helper()
	if _, ok := store.ctx.Deadline(); !ok {
		t.Fatal("store call without deadline")
	}
	if store.ctx.Value(struct{}{}) != parent.Value(struct{}{}) {
		t.Fatal("parent values lost")
	}
	if store.ctx.Err() != context.Canceled {
		t.Fatal("child context not released")
	}
}
func TestHiddenArchiveUnitCRUDAndPages(t *testing.T) {
	store := &hiddenSamples{rows: map[int64]Sample{}}
	a := NewArchive(store, 3*time.Second)
	ctx := context.WithValue(context.Background(), struct{}{}, "identity")
	rows, next, e := a.Page(ctx, 0, 1)
	if e != nil || rows == nil || len(rows) != 0 || next != 0 {
		t.Fatal("empty", rows, next, e)
	}
	rng := rand.New(rand.NewSource(3002))
	want := []Sample{}
	for i := 1; i <= 23; i++ {
		want = append(want, Sample{int64(i * 3), fmt.Sprintf("rock-%02d", i), rng.Intn(100)})
	}
	for _, i := range rng.Perm(len(want)) {
		b := want[i]
		b.Tag = " \t" + strings.ToUpper(b.Tag) + "\n"
		if e := a.Add(ctx, b); e != nil {
			t.Fatal(e)
		}
		if d, ok := store.ctx.Deadline(); !ok || d.After(time.Now().Add(3*time.Second)) || store.ctx.Value(struct{}{}) != "identity" {
			t.Fatal("bounded derived context")
		}
		if store.ctx.Err() != context.Canceled {
			t.Fatal("derived context not released")
		}
	}
	for _, limit := range []int{1, 2, 7, 23, 50} {
		got := []Sample{}
		var after int64
		for attempts := 0; attempts < 30; attempts++ {
			page, n, e := a.Page(ctx, after, limit)
			if e != nil || page == nil {
				t.Fatal(e)
			}
			if store.after != after || store.limit != limit+1 {
				t.Fatal("fetch extra", store.after, store.limit)
			}
			hiddenArchiveContext(t, store, ctx)
			got = append(got, page...)
			if n == 0 {
				break
			}
			if len(page) != limit || n != page[len(page)-1].ID || n <= after {
				t.Fatal("continuation", page, n)
			}
			after = n
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatal("pagination gaps/repeats", limit, got)
		}
	}
	if e := a.Add(ctx, want[0]); !errors.Is(e, ErrSampleDuplicate) {
		t.Fatal(e)
	}
	b := want[0]
	b.Tag = "changed"
	b.Mass = 0
	if e := a.Update(ctx, b); e != nil {
		t.Fatal(e)
	}
	hiddenArchiveContext(t, store, ctx)
	if got, e := a.Get(ctx, b.ID); e != nil || got != b {
		t.Fatal("update", got, e)
	}
	hiddenArchiveContext(t, store, ctx)
	if e := a.Delete(ctx, b.ID); e != nil {
		t.Fatal(e)
	}
	hiddenArchiveContext(t, store, ctx)
	if _, e := a.Get(ctx, b.ID); !errors.Is(e, ErrSampleMissing) {
		t.Fatal(e)
	}
	if e := a.Update(ctx, b); !errors.Is(e, ErrSampleMissing) {
		t.Fatal(e)
	}
	if e := a.Delete(ctx, b.ID); !errors.Is(e, ErrSampleMissing) {
		t.Fatal(e)
	}
}
func TestHiddenArchiveUnitValidationAndFailures(t *testing.T) {
	store := &hiddenSamples{rows: map[int64]Sample{}}
	a := NewArchive(store, time.Second)
	ctx := context.Background()
	for _, b := range []Sample{{0, "x", 0}, {-1, "x", 0}, {1, " \t", 0}, {1, "x", -1}} {
		if e := a.Add(ctx, b); !errors.Is(e, ErrSample) {
			t.Fatal(e)
		}
		if e := a.Update(ctx, b); !errors.Is(e, ErrSample) {
			t.Fatal(e)
		}
	}
	for _, id := range []int64{0, -1} {
		if _, e := a.Get(ctx, id); !errors.Is(e, ErrSample) {
			t.Fatal(e)
		}
		if e := a.Delete(ctx, id); !errors.Is(e, ErrSample) {
			t.Fatal(e)
		}
	}
	for _, p := range []struct {
		after int64
		limit int
	}{{-1, 1}, {0, 0}, {0, -1}, {0, 51}} {
		if rows, n, e := a.Page(ctx, p.after, p.limit); !errors.Is(e, ErrSample) || rows != nil || n != 0 {
			t.Fatal("invalid page", rows, n, e)
		}
	}
	if store.calls != 0 {
		t.Fatal("invalid inputs reached store")
	}
	for _, cause := range []error{errors.New("offline"), context.DeadlineExceeded} {
		store.fault = cause
		if e := a.Add(ctx, Sample{1, "x", 0}); !errors.Is(e, cause) {
			t.Fatal(e)
		}
		if e := a.Update(ctx, Sample{1, "x", 0}); !errors.Is(e, cause) {
			t.Fatal(e)
		}
		if _, e := a.Get(ctx, 1); !errors.Is(e, cause) {
			t.Fatal(e)
		}
		if e := a.Delete(ctx, 1); !errors.Is(e, cause) {
			t.Fatal(e)
		}
		if rows, n, e := a.Page(ctx, 0, 1); !errors.Is(e, cause) || rows != nil || n != 0 {
			t.Fatal("partial failure", rows, n, e)
		}
	}
	for _, deadline := range []bool{false, true} {
		var bounded context.Context
		var cancel context.CancelFunc
		if deadline {
			bounded, cancel = context.WithDeadline(ctx, time.Unix(1, 0))
		} else {
			bounded, cancel = context.WithCancel(ctx)
			cancel()
		}
		defer cancel()
		before := store.calls
		if e := a.Add(bounded, Sample{1, "x", 0}); !errors.Is(e, bounded.Err()) {
			t.Fatal(e)
		}
		if e := a.Update(bounded, Sample{1, "x", 0}); !errors.Is(e, bounded.Err()) {
			t.Fatal(e)
		}
		if _, e := a.Get(bounded, 1); !errors.Is(e, bounded.Err()) {
			t.Fatal(e)
		}
		if e := a.Delete(bounded, 1); !errors.Is(e, bounded.Err()) {
			t.Fatal(e)
		}
		if _, _, e := a.Page(bounded, 0, 1); !errors.Is(e, bounded.Err()) {
			t.Fatal(e)
		}
		if store.calls != before {
			t.Fatal("canceled call reached store")
		}
	}
}
`

const ch302Integration = `package main

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"math/rand"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func hiddenMongoSamples(t *testing.T) (*mongo.Collection, *Archive) {
	t.Helper()
	uri := os.Getenv("ORBITAL_MONGO_URI")
	if uri == "" {
		t.Skip("set ORBITAL_MONGO_URI for isolated MongoDB integration")
	}
	var token [16]byte
	if _, e := cryptorand.Read(token[:]); e != nil {
		t.Fatal(e)
	}
	client, e := mongo.Connect(options.Client().ApplyURI(uri).SetServerSelectionTimeout(5 * time.Second))
	if e != nil {
		t.Fatal("invalid Mongo configuration")
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if e := client.Disconnect(ctx); e != nil {
			t.Error(e)
		}
	})
	db := client.Database("orbital_samples_" + hex.EncodeToString(token[:]))
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if e := db.Drop(ctx); e != nil {
			t.Error("drop owned database", e)
		}
	})
	collection := db.Collection("samples")
	store := NewMongoSamples(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if e := store.Init(ctx); e != nil {
		t.Fatal(e)
	}
	if e := store.Init(ctx); e != nil {
		t.Fatal("repeat init", e)
	}
	return collection, NewArchive(store, 3*time.Second)
}
func TestHiddenArchiveMongoIntegration(t *testing.T) {
	collection, a := hiddenMongoSamples(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	rows, next, e := a.Page(ctx, 0, 50)
	if e != nil || rows == nil || len(rows) != 0 || next != 0 {
		t.Fatal("empty", rows, next, e)
	}
	indexes, e := collection.Indexes().List(ctx)
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for indexes.Next(ctx) {
		var definition struct {
			Name   string ` + "`" + `bson:"name"` + "`" + `
			Unique bool   ` + "`" + `bson:"unique"` + "`" + `
			Key    bson.D ` + "`" + `bson:"key"` + "`" + `
		}
		if e := indexes.Decode(&definition); e != nil {
			indexes.Close(ctx)
			t.Fatal(e)
		}
		if definition.Name == "unique_sample_tag" {
			found = definition.Unique && len(definition.Key) == 1 && definition.Key[0].Key == "tag" && definition.Key[0].Value == int32(1)
		}
	}
	if e := indexes.Err(); e != nil {
		t.Fatal(e)
	}
	indexes.Close(ctx)
	if !found {
		t.Fatal("missing named ascending unique tag index")
	}
	rng := rand.New(rand.NewSource(3020))
	want := []Sample{}
	for i := 1; i <= 18; i++ {
		want = append(want, Sample{int64(i * 7), fmt.Sprintf("stone-%02d", i), rng.Intn(500)})
	}
	for _, i := range rng.Perm(len(want)) {
		b := want[i]
		b.Tag = " " + b.Tag + " "
		if e := a.Add(ctx, b); e != nil {
			t.Fatal(e)
		}
	}
	got := []Sample{}
	var after int64
	for i := 0; i < 10; i++ {
		page, n, e := a.Page(ctx, after, 4)
		if e != nil {
			t.Fatal(e)
		}
		got = append(got, page...)
		if n == 0 {
			break
		}
		after = n
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatal("BSON/order/pages", got)
	}
	for _, b := range want {
		if got, e := a.Get(ctx, b.ID); e != nil || got != b {
			t.Fatal("get", got, e)
		}
	}
	for _, b := range []Sample{want[0], {999, want[0].Tag, 1}} {
		if e := a.Add(ctx, b); !errors.Is(e, ErrSampleDuplicate) {
			t.Fatal("unique index", e)
		}
	}
	b := want[0]
	b.Tag = want[1].Tag
	if e := a.Update(ctx, b); !errors.Is(e, ErrSampleDuplicate) {
		t.Fatal("duplicate update", e)
	}
	b.Tag = "updated"
	b.Mass = 0
	if e := a.Update(ctx, b); e != nil {
		t.Fatal(e)
	}
	if e := a.Update(ctx, b); e != nil {
		t.Fatal("unchanged replacement is not missing", e)
	}
	if got, e := a.Get(ctx, b.ID); e != nil || got != b {
		t.Fatal("zero update", got, e)
	}
	if e := a.Delete(ctx, b.ID); e != nil {
		t.Fatal(e)
	}
	if _, e := a.Get(ctx, b.ID); !errors.Is(e, ErrSampleMissing) {
		t.Fatal(e)
	}
	if e := a.Update(ctx, b); !errors.Is(e, ErrSampleMissing) {
		t.Fatal(e)
	}
	if e := a.Delete(ctx, b.ID); !errors.Is(e, ErrSampleMissing) {
		t.Fatal(e)
	}
	// A colliding tag from concurrent callers is decided by the database index.
	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); errs <- a.Add(ctx, Sample{int64(2000 + i), "shared-tag", i}) }(i)
	}
	wg.Wait()
	close(errs)
	wins := 0
	for e := range errs {
		if e == nil {
			wins++
		} else if !errors.Is(e, ErrSampleDuplicate) {
			t.Fatal(e)
		}
	}
	if wins != 1 {
		t.Fatal("unique race", wins)
	}
	// Corrupt BSON field type must not become a successful partial page.
	if _, e := collection.InsertOne(ctx, bson.D{{Key: "_id", Value: int64(9000)}, {Key: "tag", Value: "bad"}, {Key: "mass", Value: "not-an-int"}}); e != nil {
		t.Fatal(e)
	}
	if rows, n, e := a.Page(ctx, 0, 50); e == nil || rows != nil || n != 0 {
		t.Fatal("decode failure swallowed", rows, n, e)
	}
	expired, stop := context.WithDeadline(ctx, time.Unix(1, 0))
	defer stop()
	store := NewMongoSamples(collection)
	if _, e := store.After(expired, 0, 1); !errors.Is(e, context.DeadlineExceeded) {
		t.Fatal("driver timeout", e)
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch30Lesson303Arabic, ID: "30-3-search-index", Chapter: 30, Kind: "lab", Title: "فهرس البحث", Goal: "افهرس نشرات المستعمرة ورتّب نتائجها وصفّحها من غير تكرار", Setup: seededLabRoom,
		Requires: []string{"map"}, Deps: []string{}, Bridge: []int{69}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch303Starter}},
		Lesson: `# فهرس نشرات المستعمرة

## قبل ما تكتب
- هتوصل للنشرات اللي فيها كلمة معينة من غير ما تفحص النصوص كل مرة إزاي؟
- لو نتيجتين نفس الدرجة، إيه اللي يثبت ترتيب الصفحات؟

## من النص للكلمات
البحث لازم يطبق نفس التطبيع على النشرة والسؤال. الكلمة أطول مجموعة متصلة من حروف Unicode أو أرقامها. أي حاجة تانية فاصل، ومنها الترقيم والإيموجي. حوّل لحروف صغيرة بـ ~strings.ToLower~، واقسم بـ ~strings.FieldsFunc~.
~~~go
terms := strings.FieldsFunc(strings.ToLower("Oxygen،أكسجين 7🚀"), func(r rune) bool {
 return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}) // [oxygen أكسجين 7]
~~~
~unicode.IsLetter~ بتعرف الحروف العربية وغيرها، و ~unicode.IsDigit~ بتعرف الأرقام. التطبيع هنا lowercase بس: مفيش إزالة تشكيل أو اشتقاق أو توحيد لهجات؛ ~é~ مش نفس ~e~ ومعاها combining mark.

## فهرس مقلوب بدل المسح الكامل
الفهرس المقلوب بيربط كل كلمة بالنشرات اللي ظهرت فيها ووزنها. خريطة خارجية للكلمة، وجواها خريطة من كود النشرة للوزن. قائمة الأكواد دي اسمها postings. بدل ما تمسح كل النصوص لكل سؤال، اجمع من postings للكلمات المطلوبة.
~~~go
postings := map[string]map[string]int{}
postings["oxygen"] = map[string]int{"notice-a": 4, "notice-b": 1}
~~~
كل تكرار في العنوان وزنه 3، وكل تكرار في الجسم وزنه 1. في المثال ~Dispatch{Code: "notice-a", Title: "Oxygen", Body: "oxygen"}~ الدرجة 4 لسؤال oxygen. ~Dispatch~ قيم نصية؛ انسخها لخريطة خاصة بالكود بدل الاحتفاظ بـ slice المنادي.

## الجمع والترتيب
السؤال هنا OR: وجود أي كلمة يكفي. ~seen := map[string]bool{}~ تمنع كلمة مكررة في السؤال من زيادة الدرجة. اجمع أوزان كل الكلمات المختلفة الموجودة في النشرة. السؤال اللي مفيهوش كلمات خطأ؛ السؤال بكلمة مش معروفة نتيجة فاضية، مش خطأ.

~DispatchHit{Dispatch: notice, Score: 4}~ بيربط نسخة النشرة بالدرجة. رتّب تنازلي بالدرجة، ثم تصاعدي بالكود بمقارنة Go النصية، عشان تعادل الدرجات ميتحولش لترتيب عشوائي من الخريطة.
~~~go
sort.Slice(hits, func(i, j int) bool {
 if hits[i].Score != hits[j].Score { return hits[i].Score > hits[j].Score }
 return hits[i].Dispatch.Code < hits[j].Dispatch.Code
})
~~~

## الصفحات والملكية
~offset~ عدد النتائج المرتبة اللي هنتخطاها، و ~limit~ أقصى عدد راجع. احسب الدرجة والترتيب قبل قطع الصفحة. لو ~offset~ أكبر من أو يساوي عدد النتائج، رجّع slice فاضية. احسب طول الباقي قبل الجمع عشان ~offset~ كبير ميعملش overflow.
~~~go
end := offset + min(limit, len(hits)-offset)
page := hits[offset:end]
~~~
~min(a,b)~ دالة Go مدمجة بترجع الأصغر. الـ slice دي آمنة للمنادي لو ~hits~ اتبنت جديدة لكل نداء، مش بافر مشترك. الفهرس ثابت بعد البناء، فالقرايات المتوازية مبتحتاجش قفل طالما مبتكتبش حتى cache داخلي. ثبات الفهرس هو سبب إن صفحات متتابعة مبتكررّش أو تفقد نتائج.

## النجمة التالتة: قاموس البادئات
الفهرس عنده كلمات جاهزة، فنقدر نعرض اقتراحات من غير مسح النشرات. ~strings.HasPrefix("oxygen", "oxy")~ بترجع ~true~. طبّع البادئة بـ ~strings.ToLower(strings.TrimSpace(prefix))~، واتأكد إن كل rune حرف أو رقم؛ البادئة الفاضية مسموحة وتعني القاموس كله. ~sort.Strings(terms)~ بترتب الكلمات بمقارنة Go النصية، ثم خد الحد المطلوب.
~~~go
suggestions, err := index.Vocabulary(" OX ", 3)
~~~

## أخطاء شائعة
~~~go
for _, term := range strings.Split(query, " ") { scores[term]++ } // بيضيع الترقيم والعربي وبيكرر السؤال
page := hits[offset:offset+limit] // panic لو النهاية برا النتائج أو حصل overflow
for code := range scores { hits = append(hits, DispatchHit{Dispatch: docs[code]}) } // مفيش ترتيب ولا درجة
~~~

## المهمة
نفّذ ~BuildDispatchIndex(docs []Dispatch) (*DispatchIndex, error)~. ارفض ~Code~ الفاضي أو المكرر بـ ~nil,ErrSearch~. متعملش trim للكود؛ المسافات كود صالح طالما مش فاضي. العنوان والجسم الفاضيين وقائمة النشرات الفاضية صالحين. ابنِ inverted index بخرائط وامتلك نسخة من النشرات: تعديل slice المدخل بعد البناء ميغيّرش البحث. ~Title~ و ~Body~ يرجعوا بالنص الأصلي.

نفّذ ~Search(query string, offset, limit int) ([]DispatchHit, error)~. طبّع زي الشرح، وارفض سؤال من غير كلمات أو ~offset < 0~ أو ~limit~ برا 1 لـ 25 بـ ~nil,ErrSearch~. طابق OR للكلمات المختلفة، ودرجة كل نشرة مجموع 3 لكل ظهور كلمة سؤال في العنوان و1 لكل ظهور في الجسم. كلمة السؤال المكررة تتحسب مرة واحدة. كل نشرة تظهر مرة واحدة، بترتيب درجة تنازلي وبعدها ~Code~ تصاعدي. اقطع الصفحة بعد الترتيب، والنتيجة الفاضية slice مش ~nil~. ~offset~ بعد النهاية صالح، حتى أكبر int. مثال: ~hits, err := index.Search("AIR مؤن air", 0, 5)~، وبعدها ~offset~ يزيد بعدد الراجع.

متغيّرش المدخل ولا الفهرس ولا نتائج نداء سابق؛ نتائج المنادي قابلة للتعديل من غير أثر على أي قراية تانية. البحث آمن بالتزامن على نفس الفهرس، ومش محتاج I/O أو ساعة أو شبكة أو ~Report~. سيب ~main~ فاضية.

النجمة التالتة: نفّذ ~Vocabulary(prefix string, limit int) ([]string, error)~: البادئة بعد trim وlowercase إما فاضية أو حروف وأرقام بس، والحد من 1 لـ 20، وإلا ~nil,ErrSearch~. رجّع كلمات الفهرس المختلفة اللي بتبدأ بالبادئة، مرتبة تصاعدي ومحدودة، وslice فاضية مش ~nil~ لو مفيش. تعديل النتيجة ميغيّرش الفهرس ولا النداءات التالية.
`,
		Hints: []string{"الفهرس المقلوب يربط كل كلمة بالمستندات اللي ظهرت فيها.", "اجمع postings ثم رتّب النتائج قبل ما تقطع الصفحة.", "استخدم ~unicode.IsLetter~ و~sort.Strings~ في مواضعهما."},
		En: &LevelText{LessonRevision: 1, Steps: ch30Lesson303English, Title: "Search Index", Goal: "Index colony dispatches, rank hits, and paginate without duplicates", Files: []LabFile{{Name: "main.go", Starter: ch303Starter}},
			Lesson: `# An index for colony dispatches

## Before you write
- How can you find dispatches containing a word without scanning every text again?
- If two hits tie, what keeps page order stable?

## Text into terms
Use identical normalization for documents and queries. A term is a maximal contiguous run of Unicode letters or digits. Everything else separates terms, including punctuation and emoji. Lowercase with ~strings.ToLower~ and split with ~strings.FieldsFunc~.
~~~go
terms := strings.FieldsFunc(strings.ToLower("Oxygen,Water 7🚀"), func(r rune) bool {
 return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}) // [oxygen water 7]
~~~
~unicode.IsLetter~ recognizes letters across scripts; ~unicode.IsDigit~ recognizes digits. Normalization here is only lowercase: no accent removal, stemming, or language normalization. A precomposed accented letter differs from a base letter followed by a combining mark.

## Invert the lookup
An inverted index maps each term to dispatches containing it and their weights. An outer map names the term; an inner map maps dispatch code to weight. These code lists are called postings. A query gathers from its terms' postings rather than rescanning the original texts.
~~~go
postings := map[string]map[string]int{}
postings["oxygen"] = map[string]int{"notice-a": 4, "notice-b": 1}
~~~
Each title occurrence weighs 3, and each body occurrence weighs 1. ~Dispatch{Code: "notice-a", Title: "Oxygen", Body: "oxygen"}~ scores 4 for oxygen. ~Dispatch~ consists of string values: copy them into a private map by code rather than keeping the caller's slice.

## Accumulate and rank
This query uses OR: any matching term is enough. ~seen := map[string]bool{}~ prevents repeated query terms from increasing scores. Sum the weights for each distinct query term in each dispatch. A query with no terms is an error; an unknown word produces empty results, not an error.

~DispatchHit{Dispatch: notice, Score: 4}~ couples a dispatch copy with its score. Sort by score descending, then code ascending using Go string comparison, so tied scores do not inherit random map iteration order.
~~~go
sort.Slice(hits, func(i, j int) bool {
 if hits[i].Score != hits[j].Score { return hits[i].Score > hits[j].Score }
 return hits[i].Dispatch.Code < hits[j].Dispatch.Code
})
~~~

## Pages and ownership
~offset~ counts sorted hits to skip; ~limit~ caps returned hits. Score and sort before slicing. If offset is at or beyond the total, return an empty slice. Compute remaining length before addition so a huge offset cannot overflow.
~~~go
end := offset + min(limit, len(hits)-offset)
page := hits[offset:end]
~~~
~min(a,b)~ is a Go builtin returning the smaller operand. This slice is safe for callers to mutate when ~hits~ is fresh for each call, rather than shared storage. The index is immutable after building: concurrent reads need no lock provided they never write even an internal cache. Immutability makes consecutive pages free from repeats and omissions.

## Third star: prefix vocabulary
The index already owns a dictionary of terms, so suggestions need no document scan. ~strings.HasPrefix("oxygen", "oxy")~ returns true. Normalize with ~strings.ToLower(strings.TrimSpace(prefix))~ and require every rune to be a letter or digit. Empty prefix is allowed and means the whole dictionary. ~sort.Strings(terms)~ uses ascending Go string order; then apply the requested cap.
~~~go
suggestions, err := index.Vocabulary(" OX ", 3)
~~~

## Common mistakes
~~~go
for _, term := range strings.Split(query, " ") { scores[term]++ } // Punctuation and repetition mishandled
page := hits[offset:offset+limit] // Out-of-bounds or overflow panic
for code := range scores { hits = append(hits, DispatchHit{Dispatch: docs[code]}) } // No order or score
~~~

## Your mission
Implement ~BuildDispatchIndex(docs []Dispatch) (*DispatchIndex, error)~. Empty or duplicate ~Code~ returns ~nil,ErrSearch~. Do not trim codes: whitespace codes are valid when nonempty. Empty titles, bodies, and input collections are valid. Build a map-based inverted index and own a copy of dispatches: later mutation of the input slice must not change search. Return ~Title~ and ~Body~ exactly as originally supplied.

Implement ~Search(query string, offset, limit int) ([]DispatchHit, error)~. Normalize as described and reject no-term queries, ~offset < 0~, or limits outside 1..25 as ~nil,ErrSearch~. Match OR over distinct terms. Score each dispatch as 3 per occurrence of a query term in its title plus 1 per occurrence in its body. Count repeated query terms once. Each dispatch appears once, ordered by descending score then ascending Code. Slice only after sorting; return a non-nil empty slice when empty. Offsets beyond the end are valid, even the maximum int. Example: ~hits, err := index.Search("AIR water air", 0, 5)~, then advance offset by the returned count.

Never mutate inputs, index state, or earlier call results. Mutating a returned hit must not affect another read. Searches on one index are concurrent-safe. No I/O, clock, network, or ~Report~; keep ~main~ empty.

Third star: implement ~Vocabulary(prefix string, limit int) ([]string, error)~. After trimming and lowercasing, prefix must be empty or letters/digits only; limit is 1..20, otherwise ~nil,ErrSearch~. Return distinct indexed terms starting with that prefix in ascending order, capped by limit, with a non-nil empty slice when absent. Mutating the returned slice must not affect the index or future calls.
`, Hints: []string{"An inverted index maps each term to the documents containing it.", "Gather postings and rank results before slicing a page.", "Use ~unicode.IsLetter~ and ~sort.Strings~ where they fit."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch303Hidden},
	})
}

const ch303Starter = `package main

import "errors"

type Dispatch struct{ Code, Title, Body string }
type DispatchHit struct {
	Dispatch Dispatch
	Score    int
}

var ErrSearch = errors.New("invalid search input")

type DispatchIndex struct{}

func BuildDispatchIndex(docs []Dispatch) (*DispatchIndex, error) { return nil, errors.New("TODO") }
func (index *DispatchIndex) Search(query string, offset, limit int) ([]DispatchHit, error) {
	return nil, errors.New("TODO")
}
func (index *DispatchIndex) Vocabulary(prefix string, limit int) ([]string, error) {
	return nil, errors.New("TODO")
}
func main() {}
`
const ch303Hidden = `package main

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"
)

var hiddenWordPattern = regexp.MustCompile(` + "`" + `[\p{L}\p{Nd}]+` + "`" + `)

func hiddenSearchOracle(docs []Dispatch, query string) []DispatchHit {
	unique := map[string]bool{}
	for _, term := range hiddenWordPattern.FindAllString(strings.ToLower(query), -1) {
		unique[term] = true
	}
	hits := []DispatchHit{}
	for _, d := range docs {
		score := 0
		for _, term := range hiddenWordPattern.FindAllString(strings.ToLower(d.Title), -1) {
			if unique[term] {
				score += 3
			}
		}
		for _, term := range hiddenWordPattern.FindAllString(strings.ToLower(d.Body), -1) {
			if unique[term] {
				score++
			}
		}
		if score > 0 {
			hits = append(hits, DispatchHit{d, score})
		}
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		return hits[i].Dispatch.Code < hits[j].Dispatch.Code
	})
	return hits
}
func hiddenDispatches() []Dispatch {
	rng := rand.New(rand.NewSource(3003))
	words := []string{"AIR", "air", "مؤن", "灯", "7", "café", "dock", "unknown", "a½b", "Ⅳ", "²x", "٢"}
	docs := []Dispatch{}
	for i := 0; i < 90; i++ {
		d := Dispatch{Code: fmt.Sprintf("deck-%03d", i)}
		for j := 0; j < rng.Intn(8)+1; j++ {
			d.Title += words[rng.Intn(len(words))] + ",🚀"
		}
		for j := 0; j < rng.Intn(30)+1; j++ {
			d.Body += words[rng.Intn(len(words))] + " / "
		}
		docs = append(docs, d)
	}
	docs = append(docs, Dispatch{"empty", "", ""}, Dispatch{"tie-a", "same", "same"}, Dispatch{"tie-b", "same", "same"})
	return docs
}
func TestHiddenDispatchRankingNormalizationAndPages(t *testing.T) {
	docs := hiddenDispatches()
	original := append([]Dispatch(nil), docs...)
	rng := rand.New(rand.NewSource(3030))
	rng.Shuffle(len(docs), func(i, j int) { docs[i], docs[j] = docs[j], docs[i] })
	index, e := BuildDispatchIndex(docs)
	if e != nil || index == nil {
		t.Fatal(e)
	}
	// Mutating caller-owned storage must not mutate the immutable index.
	for i := range docs {
		docs[i] = Dispatch{"changed", "changed", "changed"}
	}
	queries := []string{"air", "AIR air", "air مؤن", "7 café", "灯", "missing", "same", "  مؤن،灯! ", "e\u0301", "🚀Air🚀", "a½b", "٢", "²x"}
	for i := 0; i < 30; i++ {
		queries = append(queries, fmt.Sprintf("%s / %s", []string{"air", "مؤن", "7"}[rng.Intn(3)], []string{"dock", "灯", "café"}[rng.Intn(3)]))
	}
	for _, query := range queries {
		want := hiddenSearchOracle(original, query)
		for _, limit := range []int{1, 7, 25} {
			all := []DispatchHit{}
			for offset := 0; offset <= len(want)+limit; offset += limit {
				got, e := index.Search(query, offset, limit)
				if e != nil || got == nil || len(got) > limit {
					t.Fatal("page", got, e)
				}
				start := min(offset, len(want))
				end := min(start+limit, len(want))
				if !reflect.DeepEqual(got, want[start:end]) {
					t.Fatalf("query=%q offset=%d limit=%d got=%v want=%v", query, offset, limit, got, want[start:end])
				}
				all = append(all, got...)
			}
			if !reflect.DeepEqual(all, want) {
				t.Fatal("duplicates or omissions")
			}
		}
	}
	// Returned objects and slices must not provide write access to index state.
	hit, e := index.Search("same", 0, 1)
	if e != nil || len(hit) != 1 {
		t.Fatal(hit, e)
	}
	hit[0].Dispatch.Title = "corrupted"
	hit[0].Score = 999
	again, e := index.Search("same", 0, 1)
	if e != nil || again[0].Score != 4 || again[0].Dispatch.Code != "tie-a" || again[0].Dispatch.Title != "same" {
		t.Fatal("result alias", again, e)
	}
}
func TestHiddenDispatchBoundariesAndParallelReads(t *testing.T) {
	for _, docs := range [][]Dispatch{{{"", "x", ""}}, {{"a", "x", ""}, {"a", "y", ""}}} {
		if index, e := BuildDispatchIndex(docs); index != nil || !errors.Is(e, ErrSearch) {
			t.Fatal("invalid build", index, e)
		}
	}
	empty, e := BuildDispatchIndex(nil)
	if e != nil {
		t.Fatal(e)
	}
	if got, e := empty.Search("air", 0, 25); e != nil || got == nil || len(got) != 0 {
		t.Fatal("empty index", got, e)
	}
	index, e := BuildDispatchIndex(hiddenDispatches())
	if e != nil {
		t.Fatal(e)
	}
	for _, query := range []string{"", " \t", "🚀!?", "Ⅳ", "½"} {
		if got, e := index.Search(query, 0, 1); got != nil || !errors.Is(e, ErrSearch) {
			t.Fatal("empty query", got, e)
		}
	}
	for _, p := range [][2]int{{-1, 1}, {0, 0}, {0, -1}, {0, 26}} {
		if got, e := index.Search("air", p[0], p[1]); got != nil || !errors.Is(e, ErrSearch) {
			t.Fatal("invalid page", got, e)
		}
	}
	if got, e := index.Search("air", int(^uint(0)>>1), 25); e != nil || got == nil || len(got) != 0 {
		t.Fatal("large offset", got, e)
	}
	want := hiddenSearchOracle(hiddenDispatches(), "AIR مؤن")[:25]
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				if got, e := index.Search("AIR مؤن", 0, 25); e != nil || !reflect.DeepEqual(got, want) {
					t.Error("parallel search", e)
					return
				}
			}
		}()
	}
	wg.Wait()
}
func TestHiddenBonusDispatchVocabulary(t *testing.T) {
	docs := hiddenDispatches()
	index, e := BuildDispatchIndex(docs)
	if e != nil {
		t.Fatal(e)
	}
	dictionary := map[string]bool{}
	for _, d := range docs {
		for _, word := range hiddenWordPattern.FindAllString(strings.ToLower(d.Title+" "+d.Body), -1) {
			dictionary[word] = true
		}
	}
	for _, prefix := range []string{"", " A ", "م", "ca", "7", "missing"} {
		want := []string{}
		for term := range dictionary {
			if strings.HasPrefix(term, strings.ToLower(strings.TrimSpace(prefix))) {
				want = append(want, term)
			}
		}
		sort.Strings(want)
		for _, limit := range []int{1, 3, 20} {
			expected := want[:min(limit, len(want))]
			got, e := index.Vocabulary(prefix, limit)
			if e != nil || got == nil || !reflect.DeepEqual(got, expected) {
				t.Fatal("dictionary", prefix, got, expected, e)
			}
			if len(got) > 0 {
				got[0] = "modified"
			}
		}
	}
	for _, prefix := range []string{"a b", "a!", "🚀"} {
		if got, e := index.Vocabulary(prefix, 2); got != nil || !errors.Is(e, ErrSearch) {
			t.Fatal("invalid prefix", got, e)
		}
	}
	for _, limit := range []int{-1, 0, 21} {
		if got, e := index.Vocabulary("", limit); got != nil || !errors.Is(e, ErrSearch) {
			t.Fatal("invalid limit", got, e)
		}
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch30Lesson304Arabic, ID: "30-4-read-write-split", Chapter: 30, Kind: "lab", Title: "قراءة وكتابة", Goal: "وجّه الكتابة للأصل والقراءة للنسخة مع نافذة قراية لزقة", Setup: seededLabRoom,
		Requires: []string{"iface"}, Deps: []string{}, Bridge: []int{70}, Race: true, Service: "",
		Files: []LabFile{{Name: "main.go", Starter: ch304Starter}},
		Lesson: `# مسارين لسجل المرحّلات

## قبل ما تكتب
- لو القراءة رجعت «مش موجود»، هل ده دليل إن النسخة وقعت؟
- بعد كتابة ناجحة، إزاي جلسة تشوف كتابتها من غير ما تغيّر جلسة تانية؟

## الأصل والنسخة
الـ ~primary~ هو مخزن الكتابة المعتمد. الـ ~replica~ نسخة بتستقبل تحديثاته بتأخير؛ قراءة عادية منها ممكن تشوف قيمة قديمة. فصل القراءة والكتابة يخفف الحمل، لكنه مش معناه إن الاتنين متطابقين في كل لحظة. هنا بنحاكي قرار التوجيه بس، مش بننفّذ بروتوكول النسخ نفسه.

~RelayStore~ واجهة صغيرة بتعرّفها بنفسك: قراءة مفتاح ترجع نص أو خطأ، وكتابة مفتاح ونص ترجع خطأ. نفس الواجهة ليها تنفيذين محقونين: ~primary~ وreplica. التستات بتناديهم في الذاكرة، من غير شبكة.
~~~go
content, err := replica.Read(ctx, "weather")
err = primary.Write(ctx, "weather", "clear")
~~~
الكتابة تروح للأصل مرة واحدة بس. متعملش retry تلقائي: ممكن الخطأ حصل بعد الحفظ فعلاً. الـ ~replica~ عمرها ما تستقبل كتابة من الراوتر.

## عطل إتاحة مش خطأ مجال
~ErrReplicaUnavailable~ علامة إن النسخة غير متاحة. اكتشفها بـ ~errors.Is~ حتى لو ملفوفة. خطأ غياب أو رفض صلاحية مش عطل نسخة، ولازم يرجع من غير fallback. مثال: ~errors.Is(fmt.Errorf("link: %w", ErrReplicaUnavailable), ErrReplicaUnavailable)~ نتيجته ~true~.

الـ fallback هنا محاولة واحدة على ~primary~ بعد عطل ~replica~ في قراءة عادية. لو الأصل نفسه فشل، رجّع خطأه من غير حلقة إعادة. أي خطأ يطابق ~context.Canceled~ أو ~context.DeadlineExceeded~ يرجع فوراً حتى لو بيطابق علامة العطل كمان.
~~~go
if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
 return "", err
}
if errors.Is(err, ErrReplicaUnavailable) { return primary.Read(ctx, key) }
~~~
افحص ~ctx.Err()~ كمان بعد نداء القراءة: لو الطلب اتلغى أثناء خطأ اتصال، ارجع خطأ السياق بدل محاولة جديدة على الأصل. السياق الأصلي لازم يعدّي لكل backend.

## قراية لزقة بعد الكتابة
بعد كتابة ناجحة، اربط الجلسة بالـ ~primary~ لمدة ~sticky~. كل مفاتيح نفس الجلسة تقرا من الأصل خلال المدة، مش بس المفتاح المكتوب. جلسة تانية تفضل على ~replica~. ده read-your-writes داخل النافذة، مش ضمان دائم إن النسخة لحقت.

الساعة دالة محقونة ~now func() time.Time~. كل كتابة ناجحة تسجل نهايتها ~now().Add(sticky)~، وقراية الجلسة بتستخدم ~now().Before(until)~. عند النهاية بالظبط القراءة عادية. التست بيحرّك الوقت بإيده؛ ممنوع ~time.Now~ أو نوم جوه الراوتر.
~~~go
instant := time.Unix(1700000000, 0)
now := func() time.Time { return instant }
until := now().Add(5*time.Second)
stillSticky := now().Before(until)
instant = instant.Add(5*time.Second) // عند النهاية: مش sticky
_ = stillSticky
~~~
خريطة ~map[string]time.Time~ من الجلسة لنهاية النافذة تكفي. القراية مبتجددش النافذة؛ كتابة ناجحة جديدة بس بتجددها. الكتابة الفاشلة متغيرش أي علامة. احذف علامة الجلسة المنتهية لما تتقري؛ من غير cleanup شامل، الجلسات اللي مبتقراش تاني ممكن تفضل في الخريطة. ده حد واضح للمحاكاة، مش خدمة جلسات للإنتاج.

## التزامن والحدود
~sync.Mutex~ حوالين العملية كلها اختيار بسيط يحافظ على ترتيب الكتابة وعلامتها والقراءة التالية. مثال: ~r.mu.Lock(); defer r.mu.Unlock()~. القفل بيسلسل الـ backends والساعة المحقونة كمان؛ متندهش الراوتر من جوه backend عشان القفل مش قابل للدخول تاني.

بعد انتهاء النافذة، قراءة ~replica~ ممكن تفضل قديمة. ضمان أقوى يحتاج إثبات تقدم النسخ، زي مقارنة رقم نسخة محفوظ، وده برا عقد المحاكاة دي. متوصفش مجرد مرور الوقت إنه إثبات إن النسخة لحقت.

## أخطاء شائعة
~~~go
if err != nil { return primary.Read(ctx, key) } // بيخبي أخطاء المجال والإلغاء
replica.Write(ctx, key, content) // النسخة مش مكان الكتابة
ununtil = time.Now().Add(sticky) // ساعة حقيقية وتجديد غلط لو في مسار القراءة
~~~

## المهمة
نفّذ ~NewRelayRouter(primary, replica RelayStore, now func() time.Time, sticky time.Duration) *RelayRouter~. المنادي بيوفر مخزنين وساعة مش ~nil~ ومدة موجبة؛ الساعة غير متراجعة، والـ backends مبتناديش الراوتر. كل كائن مستقل وبيبدأ من غير جلسات لزقة. سيب ~main~ فاضية.

~Write(ctx, session, key, content) error~ تفحص إلغاء السياق الأول، وبعده ترفض ~session~ أو ~key~ الفاضي بـ ~ErrRoute~ قبل أي نداء. النص الفاضي صالح، والمسافات مبتتشالش من أي نص. تكتب على ~primary~ مرة واحدة بنفس السياق؛ الخطأ يرجع زي ما هو من غير retry ولا تغيير النافذة. النجاح يسجل نهاية نافذة الجلسة من الساعة بعد رجوع الكتابة. لو الـ backend رجّع نجاح، اعتبره حفظ حتى لو السياق اتلغى في اللحظة دي، وسجّل النافذة وارجع ~nil~.

~Read(ctx, session, key) (string, error)~ بنفس تحقق الإلغاء والمدخل، وبنفس السياق للـ backend. لو الجلسة جوه نافذتها اقرا ~primary~ بس. غير كده اقرا ~replica~؛ نجاحها يرجع النص، وخطأ مجال يرجع فوراً، وعطلها اللي يطابق ~ErrReplicaUnavailable~ يعمل محاولة واحدة على ~primary~. الإلغاء أو المهلة سواء من السياق الحالي بعد القراءة أو من الخطأ يمنعوا fallback؛ خطأ السياق الحالي له أولوية. عند أي خطأ قراءة ارجع نص فاضي والخطأ. فشل ~primary~ في مسار ~sticky~ ميرجعش للـ ~replica~. النهاية بالظبط خارج النافذة، والقراية متجددهاش. مثال: ~router.Write(ctx, "crew-a", "weather", "clear")~ وبعدها ~router.Read(ctx, "crew-a", "weather")~ تقرا الأصل خلال النافذة.

ادعم نداءات متزامنة على نفس الكائن، ومن غير تسريب حالة جلسة لجلسة تانية. التستات seeded بساعة ومخازن محقونين؛ مفيش شبكة أو ساعة جهاز أو ~Report~.
`,
		Hints: []string{"الـreplica ممكن تتأخر عن الـprimary حتى بعد كتابة ناجحة.", "افصل مسار الكتابة عن القراية اللزقة والقراية العادية.", "راجع ~errors.Is~ والساعة المحقونة عند قرار fallback."},
		En: &LevelText{LessonRevision: 1, Steps: ch30Lesson304English, Title: "Read and Write", Goal: "Route writes to primary and reads to replica with a sticky read window", Files: []LabFile{{Name: "main.go", Starter: ch304Starter}},
			Lesson: `# Two routes for relay records

## Before you write
- Does a read returning "missing" prove that the replica is unavailable?
- How can one session read its own writes without changing another session's routing?

## Primary and replica
The primary is the authoritative write store. A replica receives its updates with delay, so an ordinary read can see older content. Splitting reads and writes distributes load; it does not make both stores identical at every instant. This lab models routing decisions, not the replication protocol itself.

Define a small ~RelayStore~ interface: reading a key returns text or an error, and writing a key and content returns an error. Two injected implementations play primary and replica. Tests call them in memory without networking.
~~~go
content, err := replica.Read(ctx, "weather")
err = primary.Write(ctx, "weather", "clear")
~~~
Every write goes to primary exactly once. Never retry automatically: an error might have occurred after commitment. The router never writes to replica.

## Availability versus domain failure
~ErrReplicaUnavailable~ identifies replica unavailability. Recognize it through ~errors.Is~, including wrapped errors. Missing records and permission failures are domain outcomes and return without fallback. Example: ~errors.Is(fmt.Errorf("link: %w", ErrReplicaUnavailable), ErrReplicaUnavailable)~ is true.

Fallback means one primary attempt after replica availability failure during an ordinary read. Primary failure returns without any retry loop. Errors matching ~context.Canceled~ or ~context.DeadlineExceeded~ return immediately, even if they also match the availability sentinel.
~~~go
if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
 return "", err
}
if errors.Is(err, ErrReplicaUnavailable) { return primary.Read(ctx, key) }
~~~
Also check ~ctx.Err()~ after the backend read: if cancellation happened during a transport failure, return context error instead of starting another primary call. Forward the original context to every backend.

## Sticky reads after writes
After a successful write, pin the session to primary for ~sticky~ duration. Every key in that session reads primary during the window, not only the written key. Other sessions keep reading replica. This provides read-your-writes within the window, not permanent proof that replica has caught up.

Time comes from an injected ~now func() time.Time~. Each successful write records ~now().Add(sticky)~; a session read checks ~now().Before(until)~. At the exact endpoint it is ordinary again. Tests move time explicitly; do not call ~time.Now~ or sleep inside the router.
~~~go
instant := time.Unix(1700000000, 0)
now := func() time.Time { return instant }
until := now().Add(5*time.Second)
stillSticky := now().Before(until)
instant = instant.Add(5*time.Second) // Exact endpoint: no longer sticky
_ = stillSticky
~~~
A ~map[string]time.Time~ from session to endpoint is enough. Reads never renew the window; only another successful write renews it. Failed writes leave markers unchanged. Remove expired markers when that session reads. Without global cleanup, sessions that never return may remain in memory: a limitation of this simulation rather than a production session service.

## Concurrency and limits
A ~sync.Mutex~ around each whole operation is a simple way to preserve write, marker, and subsequent read ordering. Example: ~r.mu.Lock(); defer r.mu.Unlock()~. The lock also serializes backends and the injected clock. Backends must not call back into the router: this lock is not reentrant.

Once the window ends, replica may still be stale. A stronger guarantee needs replication progress evidence, such as comparing an applied version. That is outside this simulation's contract. Elapsed time alone never proves catch-up.

## Common mistakes
~~~go
if err != nil { return primary.Read(ctx, key) } // Hides domain errors and cancellation
replica.Write(ctx, key, content) // Replica is not the write target
until = time.Now().Add(sticky) // Real clock, and wrong renewal if placed in the read path
~~~

## Your mission
Implement ~NewRelayRouter(primary, replica RelayStore, now func() time.Time, sticky time.Duration) *RelayRouter~. The caller supplies non-nil stores and clock, a positive duration, a nondecreasing clock, and backends that never reenter the router. Each instance is independent and starts with no sticky sessions. Keep ~main~ empty.

~Write(ctx, session, key, content) error~ checks context cancellation first, then rejects empty session or key as ~ErrRoute~ before calls. Empty content is valid; no text is trimmed. Write primary once with the original context. Return its error unchanged without retry or marker changes. Success records the session endpoint using the clock after the write returns. If the backend reports success, consider it committed even if context has just canceled: record stickiness and return ~nil~.

~Read(ctx, session, key) (string, error)~ uses the same cancellation/input validation and original context for backends. During the session window read primary only. Otherwise read replica: return its success, return domain errors immediately, and make exactly one primary attempt for errors matching ~ErrReplicaUnavailable~. Cancellation/deadline in the current context after a read, or in the returned error, prevents fallback; the current context's error takes precedence. Every failed read returns empty text and its error. Sticky primary failure never falls back to replica. The exact endpoint is outside the window, and reads never renew it. Example: ~router.Write(ctx, "crew-a", "weather", "clear")~ followed by ~router.Read(ctx, "crew-a", "weather")~ reads primary within the window.

Support concurrent calls on one object and keep sessions isolated. Tests use seeded values, injected time and stores; no network, machine clock, or ~Report~.
`, Hints: []string{"A replica can lag behind primary after a successful write.", "Separate writes, sticky reads, and ordinary reads into distinct routes.", "Use ~errors.Is~ and the injected clock when deciding on fallback."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch304Hidden},
	})
}

const ch304Starter = `package main

import (
	"context"
	"errors"
	"time"
)

var ErrRoute = errors.New("invalid route")
var ErrReplicaUnavailable = errors.New("replica unavailable")

type RelayStore interface {
	Read(context.Context, string) (string, error)
	Write(context.Context, string, string) error
}
type RelayRouter struct{}

func NewRelayRouter(primary, replica RelayStore, now func() time.Time, sticky time.Duration) *RelayRouter {
	return &RelayRouter{}
}
func (r *RelayRouter) Write(ctx context.Context, session, key, content string) error {
	return errors.New("TODO")
}
func (r *RelayRouter) Read(ctx context.Context, session, key string) (string, error) {
	return "", errors.New("TODO")
}
func main() {}
`
const ch304Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

type hiddenRelayLog struct {
	mu    sync.Mutex
	calls []string
}

func (l *hiddenRelayLog) add(call string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.calls = append(l.calls, call)
}

type hiddenRelayBackend struct {
	name              string
	log               *hiddenRelayLog
	rows              map[string]string
	readErr, writeErr error
	readHook          func()
	ctx               context.Context
}

func (s *hiddenRelayBackend) Read(ctx context.Context, key string) (string, error) {
	s.ctx = ctx
	s.log.add(s.name + ":read:" + key)
	if s.readHook != nil {
		s.readHook()
	}
	return s.rows[key], s.readErr
}
func (s *hiddenRelayBackend) Write(ctx context.Context, key, content string) error {
	s.ctx = ctx
	s.log.add(s.name + ":write:" + key)
	if s.writeErr != nil {
		return s.writeErr
	}
	s.rows[key] = content
	return nil
}
func hiddenRouterFixture() (*RelayRouter, *hiddenRelayBackend, *hiddenRelayBackend, *time.Time, *hiddenRelayLog) {
	now := time.Unix(1700000000, 0)
	log := &hiddenRelayLog{}
	p := &hiddenRelayBackend{name: "primary", log: log, rows: map[string]string{"k": "new", "other": "authoritative"}}
	replica := &hiddenRelayBackend{name: "replica", log: log, rows: map[string]string{"k": "old", "other": "stale"}}
	return NewRelayRouter(p, replica, func() time.Time { return now }, 5*time.Second), p, replica, &now, log
}
func hiddenRouterCalls(t *testing.T, log *hiddenRelayLog, want ...string) {
	t.Helper()
	if !reflect.DeepEqual(log.calls, want) {
		t.Fatal("routing order", log.calls, want)
	}
	log.calls = nil
}
func TestHiddenRelayStickySessionIsolationAndBoundary(t *testing.T) {
	r, p, replica, now, log := hiddenRouterFixture()
	ctx := context.WithValue(context.Background(), struct{}{}, "request")
	if got, e := r.Read(ctx, "A", "k"); e != nil || got != "old" {
		t.Fatal(got, e)
	}
	hiddenRouterCalls(t, log, "replica:read:k")
	if e := r.Write(ctx, "A", "k", "written"); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:write:k")
	if p.ctx != ctx {
		t.Fatal("write context lost")
	}
	p.readErr = ErrReplicaUnavailable
	if got, e := r.Read(ctx, "A", "k"); got != "" || !errors.Is(e, ErrReplicaUnavailable) {
		t.Fatal("sticky primary failure", got, e)
	}
	hiddenRouterCalls(t, log, "primary:read:k")
	p.readErr = nil
	if got, e := r.Read(ctx, "A", "k"); e != nil || got != "written" {
		t.Fatal(got, e)
	}
	hiddenRouterCalls(t, log, "primary:read:k")
	if got, e := r.Read(ctx, "A", "other"); e != nil || got != "authoritative" {
		t.Fatal(got, e)
	}
	hiddenRouterCalls(t, log, "primary:read:other")
	if got, e := r.Read(ctx, "B", "k"); e != nil || got != "old" {
		t.Fatal("session isolation", got, e)
	}
	hiddenRouterCalls(t, log, "replica:read:k")
	if replica.ctx != ctx {
		t.Fatal("read context lost")
	}
	*now = now.Add(5*time.Second - time.Nanosecond)
	if _, e := r.Read(ctx, "A", "k"); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:read:k")
	*now = now.Add(time.Nanosecond)
	if got, e := r.Read(ctx, "A", "k"); e != nil || got != "old" {
		t.Fatal("exact expiry", got, e)
	}
	hiddenRouterCalls(t, log, "replica:read:k")
	// Renew only on a successful write, not on reads or failed writes.
	if e := r.Write(ctx, "A", "k", ""); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:write:k")
	*now = now.Add(4 * time.Second)
	if e := r.Write(ctx, "A", "k", "renewed"); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:write:k")
	*now = now.Add(2 * time.Second)
	if _, e := r.Read(ctx, "A", "k"); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:read:k")
	p.writeErr = errors.New("write failed")
	if e := r.Write(ctx, "A", "k", "bad"); e != p.writeErr {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:write:k")
	*now = now.Add(3 * time.Second)
	if _, e := r.Read(ctx, "A", "k"); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "replica:read:k")
	if e := r.Write(ctx, "C", "k", "bad"); e != p.writeErr {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "primary:write:k")
	if _, e := r.Read(ctx, "C", "k"); e != nil {
		t.Fatal(e)
	}
	hiddenRouterCalls(t, log, "replica:read:k")
}
func TestHiddenRelayFallbackAndContextPrecedence(t *testing.T) {
	r, p, replica, _, log := hiddenRouterFixture()
	ctx := context.Background()
	domain := errors.New("missing manifest")
	for _, fault := range []error{ErrReplicaUnavailable, fmt.Errorf("wrapped: %w", ErrReplicaUnavailable)} {
		replica.readErr = fault
		if got, e := r.Read(ctx, "A", "k"); e != nil || got != "new" {
			t.Fatal("fallback", got, e)
		}
		hiddenRouterCalls(t, log, "replica:read:k", "primary:read:k")
		for _, primaryErr := range []error{domain, ErrReplicaUnavailable, context.DeadlineExceeded} {
			p.readErr = primaryErr
			if got, e := r.Read(ctx, "A", "k"); got != "" || !errors.Is(e, primaryErr) {
				t.Fatal("primary error", got, e)
			}
			hiddenRouterCalls(t, log, "replica:read:k", "primary:read:k")
		}
		p.readErr = nil
	}
	for _, fault := range []error{domain, context.Canceled, context.DeadlineExceeded, errors.Join(ErrReplicaUnavailable, context.Canceled), errors.Join(ErrReplicaUnavailable, context.DeadlineExceeded)} {
		replica.readErr = fault
		if got, e := r.Read(ctx, "A", "k"); got != "" || e != fault {
			t.Fatal("domain/cancel error", got, e)
		}
		hiddenRouterCalls(t, log, "replica:read:k")
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, e := r.Read(canceled, "A", "k"); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if e := r.Write(canceled, "A", "k", "x"); !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
	if len(log.calls) != 0 {
		t.Fatal("canceled request made calls")
	}
	mid, stop := context.WithCancel(ctx)
	defer stop()
	replica.readHook = stop
	replica.readErr = ErrReplicaUnavailable
	if got, e := r.Read(mid, "A", "k"); got != "" || !errors.Is(e, context.Canceled) {
		t.Fatal("cancellation became fallback", got, e)
	}
	hiddenRouterCalls(t, log, "replica:read:k")
	replica.readHook = nil
	for _, bad := range [][2]string{{"", "k"}, {"A", ""}} {
		if _, e := r.Read(ctx, bad[0], bad[1]); !errors.Is(e, ErrRoute) {
			t.Fatal(e)
		}
		if e := r.Write(ctx, bad[0], bad[1], "x"); !errors.Is(e, ErrRoute) {
			t.Fatal(e)
		}
	}
	if len(log.calls) != 0 {
		t.Fatal("invalid request made calls")
	}
}
func TestHiddenRelayConcurrentSessions(t *testing.T) {
	r, _, _, _, _ := hiddenRouterFixture()
	ctx := context.Background()
	rng := rand.New(rand.NewSource(3004))
	contents := make([]string, 32)
	for i := range contents {
		contents[i] = fmt.Sprint(rng.Int63())
	}
	var wg sync.WaitGroup
	for i, content := range contents {
		wg.Add(1)
		go func(i int, content string) {
			defer wg.Done()
			session := fmt.Sprint(i)
			key := "key-" + session
			if e := r.Write(ctx, session, key, content); e != nil {
				t.Error(e)
				return
			}
			if got, e := r.Read(ctx, session, key); e != nil || got != content {
				t.Error("read your writes", got, e)
			}
			if got, e := r.Read(ctx, "observer", "k"); e != nil || got != "old" {
				t.Error("sticky state leaked", got, e)
			}
		}(i, content)
	}
	wg.Wait()
}
`
