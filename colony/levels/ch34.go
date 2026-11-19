package levels

func init() {
	add(&Level{ID: "34-1-kv-store", Chapter: 34, Kind: "lab", Title: "مخزن مفاتيح", Goal: "ابني خزنة مفاتيح متزامنة بصلاحية ولقطة بترجع زي ما هي",
		LessonRevision: 1, Steps: ch34LessonArabic["34-1-kv-store"],
		Requires: []string{"sel:sync.RWMutex"}, Bridge: []int{87}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch341Starter}},
		Lesson: `# خزنة مفاتيح المستعمرة

## قبل ما تكتب
- لو قيمة صلاحيتها خلصت ومحدش مسحها، هل لسه موجودة؟
- لو اللقطة اتحفظت النهارده واتحمّلت بكرة، الصلاحية تتحسب من إمتى؟

## خريطة محمية بقفل قراية وكتابة
~sync.RWMutex~ بيسمح لقرّايين كتير مع بعض، أو كاتب واحد لوحده. ~RLock~ و ~RUnlock~ للقراية، و ~Lock~ و ~Unlock~ لأي تعديل. الـ map نفسها مش آمنة مع goroutines من غير قفل.
~~~go
type shelf struct {
 mu    sync.RWMutex
 items map[string]int
}
func (s *shelf) count(k string) int {
 s.mu.RLock()
 defer s.mu.RUnlock()
 return s.items[k]
}
~~~

## الصلاحية بساعة محقونة
بدل ~time.Now~ المخزن بياخد دالة ~now func() time.Time~ من اللي بينشئه. عند الكتابة احسب موعد انتهاء مطلق ~now().Add(ttl)~ وخزّنه. القيمة منتهية لما ~!now().Before(expires)~، يعني لحظة الانتهاء نفسها منتهية. صلاحية صفر معناها مفيش انتهاء.
~~~go
exp := now().Add(5 * time.Second)
alive := now().Before(exp)
~~~

## نسخ البايتات
~[]byte~ بيشاور على مصفوفة مشتركة. لو خزّنت شريحة المتصل زي ما هي، تعديله بعد ~Put~ هيغيّر المخزن. انسخ عند الدخول وعند الخروج.
~~~go
c := append([]byte(nil), v...)
~~~

## اللقطة (snapshot)
اللقطة صورة كاملة للحالة الحية في ملف أو تيار، تتحمّل بعدين وترجع نفس النتايج. هنا هتكتبها JSON بـ ~json.NewEncoder(w).Encode~ وتقراها بـ ~json.NewDecoder(r).Decode~. ~[]byte~ بيتحوّل في JSON لـ base64 تلقائي. خزّن موعد الانتهاء المطلق كـ ~UnixNano~ عشان اللقطة ما تطوّلش عمر القيمة.
~~~go
type row struct {
 Key     string ` + "`" + `json:"key"` + "`" + `
 Expires int64  ` + "`" + `json:"expires"` + "`" + `
}
err := json.NewEncoder(w).Encode([]row{{"a", 0}})
~~~
رتّب الصفوف بالمفتاح بـ ~sort.Slice~ عشان نفس الحالة تدي نفس البايتات.

## أخطاء شائعة
~~~go
func (v *Vault) Get(k string) ([]byte, bool) {
 e, ok := v.items[k] // غلط: قراية من غير RLock، والانتهاء مش متفحص
 return e.value, ok  // وكمان بيرجع الشريحة الداخلية نفسها
}
~~~

## المهمة
نفّذ ~NewVault(now func() time.Time) *Vault~ والدوال:
- ~Put(key string, value []byte, ttl time.Duration) error~: مفتاح فاضي يرجع ~ErrVaultKey~، و ~ttl~ سالب يرجع ~ErrVaultTTL~ ومفيش تغيير. صفر = من غير انتهاء. الكتابة بتستبدل القيمة والصلاحية.
- ~Get(key string) ([]byte, bool)~: نسخة من القيمة الحية؛ المنتهية أو المجهولة ترجع ~nil, false~.
- ~Delete(key string) bool~: ~true~ لو كان فيه قيمة حية واتمسحت.
- ~Len() int~: عدد القيم الحية دلوقتي.
- ~Snapshot(w io.Writer) error~: القيم الحية بس، JSON، وأخطاء الكاتب ترجع زي ما هي بـ ~errors.Is~.
- ~LoadVault(r io.Reader, now func() time.Time) (*Vault, error)~: JSON بايظ أو مفتاح فاضي أو مكرر يرجع خطأ ملفوف بـ ~ErrVaultSnapshot~. القيم اللي انتهت حسب ~now~ الجديدة تتساب.

كل الدوال آمنة مع goroutines كتير. لقطتين لنفس الحالة لازم يطلعوا نفس البايتات بالظبط مهما كان ترتيب الإدخال (نجمة إضافية). ~main~ فاضية.
`,
		Hints: []string{"الـslice ممكن تشارك backing array؛ وضّح مين يملك البايتات قبل تعديلها.", "قسّم التصميم بين حماية الحالة، وحساب الصلاحية، وحدود نسخ القيم.", "ابدأ من تعريف ~Vault~ وتوقيعات methods في الـstarter، وحدد إمتى القيمة تعتبر حية."},
		En: &LevelText{Title: "Key Vault", Goal: "Build a concurrent key vault with expiry and a faithful snapshot", LessonRevision: 1, Steps: ch34LessonEnglish["34-1-kv-store"], Files: []LabFile{{Name: "main.go", Starter: ch341Starter}},
			Lesson: `# The Colony Key Vault

## Before you write
- If a value has expired but nobody deleted it, does it still exist?
- If a snapshot is saved today and loaded tomorrow, when does expiry count from?

## A map guarded by a read/write lock
~sync.RWMutex~ admits many readers together or one writer alone. ~RLock~ and ~RUnlock~ guard reads; ~Lock~ and ~Unlock~ guard every change. A map is not safe across goroutines without a lock.
~~~go
type shelf struct {
 mu    sync.RWMutex
 items map[string]int
}
func (s *shelf) count(k string) int {
 s.mu.RLock()
 defer s.mu.RUnlock()
 return s.items[k]
}
~~~

## Expiry with an injected clock
Instead of ~time.Now~ the store receives ~now func() time.Time~ from its creator. On write compute an absolute deadline ~now().Add(ttl)~ and store it. A value has expired when ~!now().Before(expires)~, so the deadline instant itself is expired. A zero ttl means no expiry.
~~~go
exp := now().Add(5 * time.Second)
alive := now().Before(exp)
~~~

## Copying bytes
A ~[]byte~ points at a shared array. If you keep the caller's slice, their later edits change the store. Copy on the way in and on the way out.
~~~go
c := append([]byte(nil), v...)
~~~

## Snapshots
A snapshot is a complete image of the live state in a file or stream, loaded later to give the same results. Here you write JSON with ~json.NewEncoder(w).Encode~ and read it with ~json.NewDecoder(r).Decode~. JSON encodes ~[]byte~ as base64 automatically. Store the absolute deadline as ~UnixNano~ so a snapshot never extends a value's life.
~~~go
type row struct {
 Key     string ` + "`" + `json:"key"` + "`" + `
 Expires int64  ` + "`" + `json:"expires"` + "`" + `
}
err := json.NewEncoder(w).Encode([]row{{"a", 0}})
~~~
Sort rows by key with ~sort.Slice~ so equal states produce equal bytes.

## Common mistakes
~~~go
func (v *Vault) Get(k string) ([]byte, bool) {
 e, ok := v.items[k] // Wrong: read without RLock, and expiry is not checked
 return e.value, ok  // it also returns the internal slice itself
}
~~~

## Your mission
Implement ~NewVault(now func() time.Time) *Vault~ and:
- ~Put(key string, value []byte, ttl time.Duration) error~: an empty key returns ~ErrVaultKey~ and a negative ~ttl~ returns ~ErrVaultTTL~, changing nothing. Zero means no expiry. Writing replaces value and expiry.
- ~Get(key string) ([]byte, bool)~: a copy of a live value; expired or unknown keys return ~nil, false~.
- ~Delete(key string) bool~: ~true~ when a live value was removed.
- ~Len() int~: the number of live values now.
- ~Snapshot(w io.Writer) error~: live values only, as JSON; writer errors are preserved for ~errors.Is~.
- ~LoadVault(r io.Reader, now func() time.Time) (*Vault, error)~: invalid JSON or an empty or duplicate key returns an error wrapping ~ErrVaultSnapshot~. Values already expired by the new ~now~ are skipped.

Every method is safe for many goroutines. Two snapshots of the same state must be byte-identical regardless of insertion order (bonus star). Leave ~main~ empty.
`,
			Hints: []string{"A slice may share a backing array; clarify byte ownership before changing it.", "Separate state protection, expiry decisions, and value-copy boundaries.", "Start at the ~Vault~ type and starter methods, and define when a value is live."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch341Hidden},
	})
}

const ch341Starter = `package main

import (
	"errors"
	"io"
	"time"
)

var ErrVaultKey = errors.New("empty vault key")
var ErrVaultTTL = errors.New("negative vault ttl")
var ErrVaultSnapshot = errors.New("invalid vault snapshot")

type Vault struct {
	now func() time.Time
}

func NewVault(now func() time.Time) *Vault { return &Vault{now: now} }

func (v *Vault) Put(key string, value []byte, ttl time.Duration) error {
	return errors.New("TODO")
}
func (v *Vault) Get(key string) ([]byte, bool) { return nil, false }
func (v *Vault) Delete(key string) bool        { return false }
func (v *Vault) Len() int                      { return 0 }
func (v *Vault) Snapshot(w io.Writer) error    { return errors.New("TODO") }

func LoadVault(r io.Reader, now func() time.Time) (*Vault, error) {
	return nil, errors.New("TODO")
}
func main() {}
`

const ch341Hidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

type vaultClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *vaultClock) now() time.Time      { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *vaultClock) add(d time.Duration) { c.mu.Lock(); c.t = c.t.Add(d); c.mu.Unlock() }

func newVaultClock() *vaultClock { return &vaultClock{t: time.Unix(1700000000, 123)} }

type vaultModel struct {
	val []byte
	exp time.Time
}

func modelAlive(m map[string]vaultModel, k string, now time.Time) ([]byte, bool) {
	e, ok := m[k]
	if !ok || (!e.exp.IsZero() && !now.Before(e.exp)) {
		return nil, false
	}
	return e.val, true
}

// Random operations against a model: replace, expiry boundary, delete, Len.
func TestHiddenVaultModel(t *testing.T) {
	rng := rand.New(rand.NewSource(341))
	for round := 0; round < 20; round++ {
		clk := newVaultClock()
		v := NewVault(clk.now)
		model := map[string]vaultModel{}
		for i := 0; i < 200; i++ {
			k := fmt.Sprintf("k%d", rng.Intn(12))
			switch rng.Intn(5) {
			case 0, 1:
				val := []byte(fmt.Sprint(rng.Int()))
				ttl := time.Duration(rng.Intn(4)) * time.Second
				if err := v.Put(k, val, ttl); err != nil {
					t.Fatal(err)
				}
				e := vaultModel{val: append([]byte(nil), val...)}
				if ttl > 0 {
					e.exp = clk.now().Add(ttl)
				}
				model[k] = e
				val[0] = '#'
			case 2:
				_, want := modelAlive(model, k, clk.now())
				if got := v.Delete(k); got != want {
					t.Fatal("delete", k, got, want)
				}
				delete(model, k)
			case 3:
				clk.add(time.Duration(rng.Intn(3)) * time.Second)
			}
			want, wok := modelAlive(model, k, clk.now())
			got, gok := v.Get(k)
			if gok != wok || !bytes.Equal(got, want) || (!gok && got != nil) {
				t.Fatalf("get %s = %q %v, want %q %v", k, got, gok, want, wok)
			}
			if gok {
				got[0] = '!'
				if again, _ := v.Get(k); !bytes.Equal(again, want) {
					t.Fatal("Get returned internal slice")
				}
			}
			n := 0
			for mk := range model {
				if _, ok := modelAlive(model, mk, clk.now()); ok {
					n++
				}
			}
			if v.Len() != n {
				t.Fatal("len", v.Len(), n)
			}
		}
	}
}

func TestHiddenVaultValidationAndBoundary(t *testing.T) {
	clk := newVaultClock()
	v := NewVault(clk.now)
	if err := v.Put("", []byte("x"), 0); !errors.Is(err, ErrVaultKey) {
		t.Fatal("empty key", err)
	}
	if err := v.Put("a", []byte("x"), -time.Nanosecond); !errors.Is(err, ErrVaultTTL) {
		t.Fatal("negative ttl", err)
	}
	if v.Len() != 0 {
		t.Fatal("rejected writes must not change the vault")
	}
	if err := v.Put("a", nil, time.Second); err != nil {
		t.Fatal(err)
	}
	if _, ok := v.Get("a"); !ok {
		t.Fatal("empty value is still a value")
	}
	clk.add(time.Second - 1)
	if _, ok := v.Get("a"); !ok {
		t.Fatal("alive one nanosecond before the deadline")
	}
	clk.add(1)
	if _, ok := v.Get("a"); ok || v.Delete("a") || v.Len() != 0 {
		t.Fatal("deadline instant is expired")
	}
	if err := v.Put("forever", []byte("y"), 0); err != nil {
		t.Fatal(err)
	}
	clk.add(1000 * time.Hour)
	if got, ok := v.Get("forever"); !ok || string(got) != "y" {
		t.Fatal("zero ttl never expires")
	}
}

type failVaultWriter struct{ err error }

func (w failVaultWriter) Write(p []byte) (int, error) { return 0, w.err }

func TestHiddenVaultSnapshotRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(342))
	for round := 0; round < 10; round++ {
		clk := newVaultClock()
		v := NewVault(clk.now)
		type want struct {
			val []byte
			ttl time.Duration
		}
		wants := map[string]want{}
		for i := 0; i < 30; i++ {
			k := fmt.Sprintf("dock-%d", rng.Intn(40))
			val := make([]byte, rng.Intn(20))
			rng.Read(val)
			ttl := time.Duration(rng.Intn(5)) * time.Minute
			if err := v.Put(k, val, ttl); err != nil {
				t.Fatal(err)
			}
			wants[k] = want{val, ttl}
		}
		v.Put("gone", []byte("z"), time.Second)
		clk.add(time.Second)
		var buf bytes.Buffer
		if err := v.Snapshot(&buf); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(buf.String(), "gone") {
			t.Fatal("snapshot holds expired entry")
		}
		clk.add(90 * time.Second)
		w, err := LoadVault(bytes.NewReader(buf.Bytes()), clk.now)
		if err != nil {
			t.Fatal(err)
		}
		for k, e := range wants {
			got, ok := w.Get(k)
			alive := e.ttl == 0 || e.ttl > 91*time.Second
			if ok != alive || (ok && !bytes.Equal(got, e.val)) {
				t.Fatalf("loaded %s ok=%v alive=%v", k, ok, alive)
			}
		}
	}
	fault := errors.New("disk full")
	v := NewVault(newVaultClock().now)
	v.Put("a", []byte("b"), 0)
	if err := v.Snapshot(failVaultWriter{fmt.Errorf("w: %w", fault)}); !errors.Is(err, fault) {
		t.Fatal("writer error", err)
	}
}

func TestHiddenVaultBadSnapshot(t *testing.T) {
	clk := newVaultClock()
	v := NewVault(clk.now)
	v.Put("dup", []byte("1"), 0)
	var buf bytes.Buffer
	if err := v.Snapshot(&buf); err != nil {
		t.Fatal(err)
	}
	good := buf.String()
	cases := []string{"", "{", "not json", good[:len(good)/2], strings.Replace(good, "dup", "", 1)}
	var two bytes.Buffer
	v.Put("dupe", []byte("2"), 0)
	v.Snapshot(&two)
	cases = append(cases, strings.Replace(two.String(), "dupe", "dup", 1))
	for _, c := range cases {
		if _, err := LoadVault(strings.NewReader(c), clk.now); !errors.Is(err, ErrVaultSnapshot) {
			t.Fatalf("snapshot %q: %v", c, err)
		}
	}
}

func TestHiddenVaultConcurrent(t *testing.T) {
	clk := newVaultClock()
	v := NewVault(clk.now)
	var wg sync.WaitGroup
	for g := 0; g < 8; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			for i := 0; i < 300; i++ {
				k := fmt.Sprintf("k%d", i%7)
				v.Put(k, []byte{byte(g)}, time.Duration(i%3)*time.Second)
				v.Get(k)
				v.Len()
				if i%11 == 0 {
					v.Delete(k)
					v.Snapshot(new(bytes.Buffer))
					clk.add(time.Millisecond)
				}
			}
		}(g)
	}
	wg.Wait()
	for g := 0; g < 8; g++ {
		v.Put(fmt.Sprintf("final%d", g), []byte("x"), 0)
	}
	if v.Len() < 8 {
		t.Fatal("lost writes", v.Len())
	}
}

func TestHiddenBonusVaultDeterministicSnapshot(t *testing.T) {
	rng := rand.New(rand.NewSource(343))
	clk := newVaultClock()
	keys := make([]string, 25)
	for i := range keys {
		keys[i] = fmt.Sprintf("bay-%02d", i)
	}
	snap := func() string {
		v := NewVault(clk.now)
		for _, i := range rng.Perm(len(keys)) {
			v.Put(keys[i], []byte(keys[i]), time.Duration(i)*time.Hour)
		}
		var buf bytes.Buffer
		if err := v.Snapshot(&buf); err != nil {
			t.Fatal(err)
		}
		return buf.String()
	}
	first := snap()
	for i := 0; i < 5; i++ {
		if snap() != first {
			t.Fatal("same state gave different snapshot bytes")
		}
	}
}
`

func init() {
	add(&Level{ID: "34-2-wal", Chapter: 34, Kind: "lab", Title: "سجل الكتابة", Goal: "سجّل كل تعديل على القرص قبل الذاكرة وارجع لنفس الحالة بعد انهيار",
		LessonRevision: 1, Steps: ch34LessonArabic["34-2-wal"],
		Requires: []string{"sel:crc32.ChecksumIEEE"}, Bridge: []int{88}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch342Starter}},
		Lesson: `# دفتر عمليات المستعمرة

## قبل ما تكتب
- لو الكهربا قطعت وسط كتابة سجل، إزاي تعرف إن آخر سجل ناقص؟
- إمتى الأمان إنك تكمّل، وإمتى لازم ترفض تفتح الدفتر؟

## سجل الكتابة المسبقة (WAL)
الـ write-ahead log ملف بنضيف في آخره بس. كل تعديل يتكتب في السجل الأول، وبعدها يتطبق على الذاكرة. لو البرنامج وقع، إعادة قراية السجل من أوله (replay) بترجّع نفس الحالة.
~~~go
f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
if err != nil { return err }
_, err = f.Write(record) // سجل كامل في نداء واحد
~~~
~os.O_APPEND~ بيخلي كل كتابة تروح لآخر الملف. ~O_CREATE~ ينشئه لو مش موجود.

## الـ memtable
الـ memtable هي الحالة الحالية في الذاكرة: هنا ~map[string][]byte~. القراية منها بس، والسجل للاسترجاع بس. حذف مفتاح يتسجل كسجل حذف عشان الـ replay يعيده.
~~~go
mem := map[string][]byte{}
mem["dock"] = []byte("open")
delete(mem, "dock")
~~~

## الـ checksum
الـ checksum رقم صغير محسوب من البايتات؛ لو بايت اتغير غالبًا الرقم يتغير. ~crc32.ChecksumIEEE(b)~ من ~hash/crc32~ بيرجع ~uint32~. احسبه عند الكتابة، واحسبه تاني عند القراية وقارن.
~~~go
sum := crc32.ChecksumIEEE([]byte("dock"))
ok := sum == crc32.ChecksumIEEE(read)
~~~

## شكل السجل
كل سجل: رأس 13 بايت ثم المفتاح ثم القيمة. الأرقام ~binary.BigEndian~.
~~~
[0:4]   crc32.ChecksumIEEE(record[4:])
[4]     op: 1 كتابة، 2 حذف
[5:9]   طول المفتاح
[9:13]  طول القيمة (صفر للحذف)
[13:]   المفتاح ثم القيمة
~~~
~~~go
h := make([]byte, 13)
binary.BigEndian.PutUint32(h[5:9], uint32(len(key)))
n := binary.BigEndian.Uint32(h[5:9])
~~~

## ذيل ممزق ولا فساد؟
السجل الناقص في آخر الملف (رأس أقل من 13 بايت، أو أطوال بتعدّي نهاية الملف) معناه كتابة اتقطعت ومحدش اتأكدله إنها نجحت، فتسيبها. ~os.Truncate(path, size)~ بيقص الملف لآخر سجل سليم عشان الإضافات الجاية تبدأ من مكان صح.
~~~go
err := os.Truncate(path, int64(good))
~~~
لكن سجل كامل بـ checksum غلط، أو op مجهول، أو أطوال خارج الحدود، ده فساد حقيقي: ارفض الفتح ومتخمّنش.

## أخطاء شائعة
~~~go
mem[key] = value               // غلط: الذاكرة اتعدلت قبل السجل
f.Write(header); f.Write(key)  // وكمان السجل اتكتب على كذا نداء والخطأ اتبلع
~~~

## المهمة
نفّذ ~OpenLedger(dir string) (*Ledger, error)~ اللي بيفتح أو ينشئ ~dir/ledger.wal~ ويعمل replay، والدوال:
- ~Set(key string, value []byte) error~ و ~Remove(key string) error~: مفتاح فاضي أو أطول من 1024 بايت يرجع ~ErrLedgerKey~، وقيمة أطول من 65536 بايت ترجع ~ErrLedgerValue~، ومفيش كتابة. غير كده اكتب السجل كامل في نداء ~Write~ واحد، وبعد ما ينجح بس عدّل الـ memtable. ~Remove~ بيكتب سجل حذف دايمًا حتى لو المفتاح مش موجود.
- ~Lookup(key string) ([]byte, bool)~: نسخة من القيمة.
- ~Close() error~.

في الفتح: ذيل ناقص يتقص بـ ~os.Truncate~ والفتح ينجح. سجل كامل فاسد زي ما اتشرح يرجع خطأ ملفوف بـ ~ErrLedgerCorrupt~. سجل الحذف لازم طول قيمته صفر. الدفتر مش لازم يبقى آمن مع goroutines كتير. ~main~ فاضية.
`,
		Hints: []string{"السجل بيسمح بإعادة بناء الحالة بعد إعادة التشغيل.", "رتّب تفكيرك من bytes إلى record ثم حالة الذاكرة.", "ابدأ من ~Ledger~ وتوقيعات ~OpenLedger~ و~Lookup~ في الـstarter."},
		En: &LevelText{Title: "Write-Ahead Log", Goal: "Log every change to disk before memory and recover the same state after a crash", LessonRevision: 1, Steps: ch34LessonEnglish["34-2-wal"], Files: []LabFile{{Name: "main.go", Starter: ch342Starter}},
			Lesson: `# The Colony Operations Ledger

## Before you write
- If power fails mid-write, how do you know the last record is incomplete?
- When is it safe to continue, and when must you refuse to open the ledger?

## Write-ahead logging (WAL)
A write-ahead log is a file you only append to. Each change is written to the log first and applied to memory afterwards. After a crash, replaying the log from the start rebuilds the same state.
~~~go
f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
if err != nil { return err }
_, err = f.Write(record) // one whole record in one call
~~~
~os.O_APPEND~ sends every write to the end of the file. ~O_CREATE~ creates it if missing.

## The memtable
The memtable is the current state in memory: here a ~map[string][]byte~. Reads use it alone; the log exists for recovery. Deleting a key is logged as a delete record so replay repeats it.
~~~go
mem := map[string][]byte{}
mem["dock"] = []byte("open")
delete(mem, "dock")
~~~

## Checksums
A checksum is a small number computed from bytes; changing a byte almost always changes it. ~crc32.ChecksumIEEE(b)~ from ~hash/crc32~ returns a ~uint32~. Compute it when writing, recompute it when reading, and compare.
~~~go
sum := crc32.ChecksumIEEE([]byte("dock"))
ok := sum == crc32.ChecksumIEEE(read)
~~~

## Record layout
Each record is a 13-byte header, then the key, then the value. Numbers are ~binary.BigEndian~.
~~~
[0:4]   crc32.ChecksumIEEE(record[4:])
[4]     op: 1 set, 2 delete
[5:9]   key length
[9:13]  value length (zero for delete)
[13:]   key then value
~~~
~~~go
h := make([]byte, 13)
binary.BigEndian.PutUint32(h[5:9], uint32(len(key)))
n := binary.BigEndian.Uint32(h[5:9])
~~~

## Torn tail or corruption?
An incomplete record at the end (a header shorter than 13 bytes, or lengths that run past the end of the file) is an interrupted write that nobody was told succeeded, so drop it. ~os.Truncate(path, size)~ cuts the file back to the last good record so later appends start in the right place.
~~~go
err := os.Truncate(path, int64(good))
~~~
A complete record with a wrong checksum, an unknown op, or lengths outside the limits is real corruption: refuse to open rather than guess.

## Common mistakes
~~~go
mem[key] = value               // Wrong: memory changed before the log
f.Write(header); f.Write(key)  // the record is split across calls and errors are lost
~~~

## Your mission
Implement ~OpenLedger(dir string) (*Ledger, error)~, which opens or creates ~dir/ledger.wal~ and replays it, and:
- ~Set(key string, value []byte) error~ and ~Remove(key string) error~: an empty key or one longer than 1024 bytes returns ~ErrLedgerKey~, and a value longer than 65536 bytes returns ~ErrLedgerValue~, writing nothing. Otherwise write the whole record in one ~Write~ call, and only after it succeeds update the memtable. ~Remove~ always writes a delete record, even for a missing key.
- ~Lookup(key string) ([]byte, bool)~: a copy of the value.
- ~Close() error~.

On open, an incomplete tail is cut with ~os.Truncate~ and opening succeeds. A complete corrupt record as described returns an error wrapping ~ErrLedgerCorrupt~. A delete record must have value length zero. The ledger need not be safe for concurrent use. Leave ~main~ empty.
`,
			Hints: []string{"A log lets the program rebuild state after a restart.", "Reason from bytes to records, then to in-memory state.", "Start at ~Ledger~ and the ~OpenLedger~ and ~Lookup~ signatures in the starter."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch342Hidden},
	})
}

const ch342Starter = `package main

import "errors"

var ErrLedgerKey = errors.New("invalid ledger key")
var ErrLedgerValue = errors.New("ledger value too large")
var ErrLedgerCorrupt = errors.New("corrupt ledger")

type Ledger struct {
	mem map[string][]byte
}

func OpenLedger(dir string) (*Ledger, error) {
	return nil, errors.New("TODO")
}

func (l *Ledger) Set(key string, value []byte) error { return errors.New("TODO") }
func (l *Ledger) Remove(key string) error            { return errors.New("TODO") }
func (l *Ledger) Lookup(key string) ([]byte, bool)   { return nil, false }
func (l *Ledger) Close() error                       { return nil }

func main() {}
`

const ch342Hidden = `package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func walRecord(op byte, key string, value []byte) []byte {
	r := make([]byte, 13, 13+len(key)+len(value))
	r[4] = op
	binary.BigEndian.PutUint32(r[5:9], uint32(len(key)))
	binary.BigEndian.PutUint32(r[9:13], uint32(len(value)))
	r = append(r, key...)
	r = append(r, value...)
	binary.BigEndian.PutUint32(r[0:4], crc32.ChecksumIEEE(r[4:]))
	return r
}

func walOpen(t *testing.T, dir string) *Ledger {
	t.Helper()
	l, err := OpenLedger(dir)
	if err != nil {
		t.Fatal("open", err)
	}
	return l
}

func walMatch(t *testing.T, l *Ledger, model map[string][]byte, keys []string) {
	t.Helper()
	for _, k := range keys {
		want, wok := model[k]
		got, ok := l.Lookup(k)
		if ok != wok || !bytes.Equal(got, want) {
			t.Fatalf("lookup %q = %q %v, want %q %v", k, got, ok, want, wok)
		}
		if ok && len(got) > 0 {
			got[0] ^= 0xff
			if again, _ := l.Lookup(k); !bytes.Equal(again, want) {
				t.Fatal("Lookup returned internal slice")
			}
		}
	}
}

func TestHiddenLedgerExactRecords(t *testing.T) {
	dir := t.TempDir()
	l := walOpen(t, dir)
	val := []byte("pressure=7")
	if err := l.Set("dock", val); err != nil {
		t.Fatal(err)
	}
	val[0] = 'X'
	if err := l.Remove("ghost"); err != nil {
		t.Fatal(err)
	}
	if err := l.Set("e", nil); err != nil {
		t.Fatal(err)
	}
	l.Close()
	want := append(append(walRecord(1, "dock", []byte("pressure=7")), walRecord(2, "ghost", nil)...), walRecord(1, "e", nil)...)
	got, err := os.ReadFile(filepath.Join(dir, "ledger.wal"))
	if err != nil || !bytes.Equal(got, want) {
		t.Fatalf("file bytes differ: %v\n got %x\nwant %x", err, got, want)
	}
	l = walOpen(t, dir)
	defer l.Close()
	if v, ok := l.Lookup("dock"); !ok || string(v) != "pressure=7" {
		t.Fatal("replay", string(v), ok)
	}
	if v, ok := l.Lookup("e"); !ok || len(v) != 0 {
		t.Fatal("empty value lost")
	}
}

func TestHiddenLedgerValidation(t *testing.T) {
	dir := t.TempDir()
	l := walOpen(t, dir)
	defer l.Close()
	if err := l.Set("", []byte("x")); !errors.Is(err, ErrLedgerKey) {
		t.Fatal("empty key", err)
	}
	if err := l.Remove(strings.Repeat("k", 1025)); !errors.Is(err, ErrLedgerKey) {
		t.Fatal("long key", err)
	}
	if err := l.Set("big", make([]byte, 65537)); !errors.Is(err, ErrLedgerValue) {
		t.Fatal("big value", err)
	}
	if fi, err := os.Stat(filepath.Join(dir, "ledger.wal")); err != nil || fi.Size() != 0 {
		t.Fatal("rejected writes must not touch the log", err)
	}
	if l.Set(strings.Repeat("k", 1024), make([]byte, 65536)) != nil {
		t.Fatal("limits are inclusive")
	}
}

// Random operations with reopen after every batch: replay must match a model.
func TestHiddenLedgerReplay(t *testing.T) {
	rng := rand.New(rand.NewSource(342))
	dir := t.TempDir()
	model := map[string][]byte{}
	keys := []string{"a", "b", "air-lock", "bay/7", "κ", "z"}
	for batch := 0; batch < 12; batch++ {
		l := walOpen(t, dir)
		walMatch(t, l, model, keys)
		for i := 0; i < 25; i++ {
			k := keys[rng.Intn(len(keys))]
			if rng.Intn(4) == 0 {
				if err := l.Remove(k); err != nil {
					t.Fatal(err)
				}
				delete(model, k)
			} else {
				v := make([]byte, rng.Intn(40))
				rng.Read(v)
				if err := l.Set(k, v); err != nil {
					t.Fatal(err)
				}
				model[k] = append([]byte{}, v...)
			}
		}
		walMatch(t, l, model, keys)
		if err := l.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

// A crash at any byte inside the last record loses only that record.
func TestHiddenLedgerTornTail(t *testing.T) {
	rng := rand.New(rand.NewSource(343))
	for round := 0; round < 40; round++ {
		dir := t.TempDir()
		path := filepath.Join(dir, "ledger.wal")
		var file []byte
		model := map[string][]byte{}
		for i := 0; i < 1+rng.Intn(6); i++ {
			k := fmt.Sprint("k", rng.Intn(4))
			v := []byte(fmt.Sprint(rng.Int()))
			file = append(file, walRecord(1, k, v)...)
			model[k] = v
		}
		last := walRecord(1, "k0", []byte("torn-value"))
		cut := rng.Intn(len(last))
		os.WriteFile(path, append(append([]byte{}, file...), last[:cut]...), 0o644)
		l := walOpen(t, dir)
		walMatch(t, l, model, []string{"k0", "k1", "k2", "k3"})
		if fi, _ := os.Stat(path); fi.Size() != int64(len(file)) {
			t.Fatal("torn tail must be truncated", fi.Size(), len(file))
		}
		if err := l.Set("after", []byte("crash")); err != nil {
			t.Fatal(err)
		}
		l.Close()
		l = walOpen(t, dir)
		model["after"] = []byte("crash")
		walMatch(t, l, model, []string{"k0", "k1", "k2", "k3", "after"})
		l.Close()
	}
}

func TestHiddenLedgerCorruption(t *testing.T) {
	rng := rand.New(rand.NewSource(344))
	first := walRecord(1, "hull", []byte("sealed"))
	good := append(append([]byte{}, first...), walRecord(1, "fuel", []byte("full-tank"))...)
	for i := 0; i < 20; i++ {
		bad := append([]byte{}, good...)
		start := 0
		if rng.Intn(2) == 1 {
			start = len(first)
		}
		pos := []int{0, 1, 2, 3, 13, 14, 15}[rng.Intn(7)] + start
		bad[pos] ^= byte(1 + rng.Intn(255))
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "ledger.wal"), bad, 0o644)
		if _, err := OpenLedger(dir); !errors.Is(err, ErrLedgerCorrupt) {
			t.Fatal("flipped byte", pos, err)
		}
	}
	for _, rec := range [][]byte{walRecord(7, "k", []byte("v")), walRecord(1, "", []byte("v")), walRecord(2, "k", []byte("v")), walRecord(1, strings.Repeat("k", 1025), nil)} {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "ledger.wal"), append(rec, walRecord(1, "ok", nil)...), 0o644)
		if _, err := OpenLedger(dir); !errors.Is(err, ErrLedgerCorrupt) {
			t.Fatalf("record %x: %v", rec[4:13], err)
		}
	}
}
`

func init() {
	add(&Level{ID: "34-3-compaction", Chapter: 34, Kind: "lab", Title: "الدمج", Goal: "فرّغ الدفتر في جداول مرتبة وادمجها من غير ما نتيجة واحدة تتغير",
		LessonRevision: 1, Steps: ch34LessonArabic["34-3-compaction"],
		Requires: []string{"sel:os.Rename"}, Bridge: []int{88}, Setup: seededLabRoom, CarryFrom: "34-2-wal",
		Files: []LabFile{{Name: "main.go", Starter: ch343Starter}},
		Lesson: `# جداول الأرشيف ودمجها

## قبل ما تكتب
- لو مسحت مفتاح موجود في ملف قديم مبيتعدلش، إزاي القراية تعرف إنه اتمسح؟
- إمتى بقى آمن ترمي علامة المسح نهائي؟

## الملف ده بيكمل من المهمة اللي فاتت
لو حليت ~34-2-wal~، الـ ~main.go~ بيبدأ بالدفتر بتاعك. السجل بيكبر للأبد؛ الحل إنك تفرّغ الـ memtable لملف مرتب وتصفّر السجل.

## الـ SSTable
الـ SSTable (sorted string table) ملف ثابت فيه سجلات مرتبة بالمفتاح، ومبيتعدلش بعد ما يتكتب. هنا كل سجل بنفس شكل سجل الدفتر بالظبط (رأس 13 بايت و checksum)، والمفاتيح تصاعدية من غير تكرار. ترتيب المفاتيح بـ ~sort.Strings~:
~~~go
keys := []string{"fuel", "air"}
sort.Strings(keys) // air, fuel
i := sort.SearchStrings(keys, "fuel") // 1
found := i < len(keys) && keys[i] == "fuel"
~~~
~sort.SearchStrings~ بحث ثنائي في شريحة مرتبة: بيرجع مكان المفتاح أو المكان اللي كان هيتحط فيه.

اسم كل جدول ~sst-000001.sst~ ورقمه بيزيد. الأحدث بيكسب: القراية تبص في الـ memtable، وبعدها الجداول من الأحدث للأقدم، وأول إجابة هي النتيجة. ~fmt.Sprintf("sst-%06d.sst", n)~ بيكوّن الاسم و ~fmt.Sscanf~ بيقراه:
~~~go
var n int
_, err := fmt.Sscanf("sst-000042.sst", "sst-%06d.sst", &n)
~~~

## علامة المسح (tombstone)
مقدرش أشيل مفتاح من جدول قديم. فالمسح بيتسجل كـ tombstone: سجل op 2 معناه «المفتاح ده اتمسح». الـ memtable لازم تحتفظ بيه كمان، لأن القراية لو ملقتش حاجة هتنزل للجداول القديمة وتلاقي القيمة الميتة.
~~~go
type slot struct {
 value []byte
 dead  bool // tombstone
}
~~~

## الكتابة الذرّية للجدول
اكتب الجدول في ملف مؤقت ~.tmp~، وبعدين ~os.Rename~ للاسم النهائي. الـ rename على نفس الفولدر ذرّي: الملف يا موجود كامل يا مش موجود. عشان كده جدول ~.sst~ ناقص أو checksum غلط يبقى فساد، مش ذيل ممزق.
~~~go
err := os.WriteFile(tmp, data, 0o644)
if err == nil { err = os.Rename(tmp, final) }
~~~
بعد الجدول بس صفّر السجل بـ ~Truncate(0)~. لو وقعنا في النص، السجل بيتعاد على نفس القيم ومفيش ضرر.

## الدمج (compaction)
الدمج بيقرا كل الجداول ويطلّع جدول واحد بآخر قيمة لكل مفتاح. لما الجدول الناتج هو الأقدم في الدنيا، الـ tombstones ملهاش لازمة: اترميها هي والمفاتيح الميتة. بعد ما الجديد يتعمله rename، امسح القديم بـ ~os.Remove~. ~os.ReadDir(dir)~ بيرجع أسماء الملفات مرتبة.
~~~go
entries, err := os.ReadDir(dir)
for _, e := range entries { fmt.Println(e.Name()) }
~~~

## أخطاء شائعة
~~~go
delete(mem, key) // غلط: الجدول القديم هيرجّع القيمة الممسوحة
os.WriteFile(filepath.Join(dir, "sst-000002.sst"), data, 0o644) // ومن غير tmp و rename
~~~

## المهمة
كمّل ~Ledger~ بنفس ~OpenLedger~ و ~Set~ و ~Remove~ و ~Lookup~ و ~Close~ وأخطاء ~34-2-wal~ وشكل السجل، وزوّد:
- ~Flush() error~: لو الـ memtable فاضية ارجع ~nil~ من غير ملفات. غير كده اكتب كل مدخلاتها، tombstones كمان، مرتبة بالمفتاح في ~sst-N.sst~ جديد (N أكبر من أي جدول موجود، بستة أرقام) عن طريق ~.tmp~ و ~os.Rename~، وبعدين صفّر ~ledger.wal~ وفضّي الـ memtable.
- ~Compact() error~: لو مفيش جداول ارجع ~nil~. غير كده ادمج كل الجداول في جدول واحد جديد رقمه أكبر من الكل، من غير tombstones، وامسح القديمة. الـ memtable متتلمسش.

القراية: الـ memtable، ثم الجداول من الأحدث. في الفتح: الجداول تتقري بالترتيب؛ أي جدول ناقص أو checksum غلط أو مفاتيح مش تصاعدية يرجع ~ErrLedgerCorrupt~. ملفات ~.tmp~ القديمة تتمسح وتتجاهل (نجمة إضافية). ~main~ فاضية.
`,
		Hints: []string{"العلامة الصريحة بتخلي الحذف حالة قابلة للقراءة.", "افصل بين تفريغ الحالة الحالية ودمج التاريخ القديم.", "راجع حقول ~Ledger~ الموروثة وتوقيعات ~Flush~ و~Compact~ في الـstarter."},
		En: &LevelText{Title: "Compaction", Goal: "Flush the ledger into sorted tables and merge them without changing a single result", LessonRevision: 1, Steps: ch34LessonEnglish["34-3-compaction"], Files: []LabFile{{Name: "main.go", Starter: ch343Starter}},
			Lesson: `# Archive Tables and Compaction

## Before you write
- If you delete a key that lives in an old read-only file, how does a read learn it was deleted?
- When is it safe to throw a delete marker away for good?

## This file continues the previous mission
If you solved ~34-2-wal~, ~main.go~ starts with your ledger. The log grows forever; the fix is to flush the memtable to a sorted file and reset the log.

## SSTables
An SSTable (sorted string table) is an immutable file of records sorted by key, never changed after it is written. Here every record has exactly the ledger record layout (13-byte header and checksum), and keys are strictly increasing. Sort keys with ~sort.Strings~:
~~~go
keys := []string{"fuel", "air"}
sort.Strings(keys) // air, fuel
i := sort.SearchStrings(keys, "fuel") // 1
found := i < len(keys) && keys[i] == "fuel"
~~~
~sort.SearchStrings~ is a binary search over a sorted slice: it returns the key's index or where it would be inserted.

Each table is named ~sst-000001.sst~ with an increasing number. The newest wins: a read checks the memtable, then tables from newest to oldest, and the first answer is the result. ~fmt.Sprintf("sst-%06d.sst", n)~ builds the name and ~fmt.Sscanf~ parses it:
~~~go
var n int
_, err := fmt.Sscanf("sst-000042.sst", "sst-%06d.sst", &n)
~~~

## Tombstones
You cannot remove a key from an old table. So a delete is recorded as a tombstone: an op 2 record meaning "this key was deleted". The memtable must keep it too, or a read that finds nothing will fall through to old tables and find the dead value.
~~~go
type slot struct {
 value []byte
 dead  bool // tombstone
}
~~~

## Writing a table atomically
Write the table to a temporary ~.tmp~ file, then ~os.Rename~ it to its final name. A rename within one directory is atomic: the file either exists complete or not at all. That is why a short ~.sst~ or a bad checksum is corruption, not a torn tail.
~~~go
err := os.WriteFile(tmp, data, 0o644)
if err == nil { err = os.Rename(tmp, final) }
~~~
Only after the table exists reset the log with ~Truncate(0)~. A crash in between replays the log onto the same values, which is harmless.

## Compaction
Compaction reads every table and writes one table holding the latest value per key. When the output is the oldest table in existence, tombstones are no longer needed: drop them along with dead keys. After the new table is renamed into place, delete the old ones with ~os.Remove~. ~os.ReadDir(dir)~ returns file entries sorted by name.
~~~go
entries, err := os.ReadDir(dir)
for _, e := range entries { fmt.Println(e.Name()) }
~~~

## Common mistakes
~~~go
delete(mem, key) // Wrong: the old table resurrects the deleted value
os.WriteFile(filepath.Join(dir, "sst-000002.sst"), data, 0o644) // and no tmp plus rename
~~~

## Your mission
Extend ~Ledger~, keeping ~OpenLedger~, ~Set~, ~Remove~, ~Lookup~, ~Close~, the errors and the record layout of ~34-2-wal~, and add:
- ~Flush() error~: with an empty memtable return ~nil~ and create no files. Otherwise write all its entries, tombstones included, sorted by key into a new ~sst-N.sst~ (N greater than any existing table, six digits) via ~.tmp~ and ~os.Rename~, then reset ~ledger.wal~ to empty and clear the memtable.
- ~Compact() error~: with no tables return ~nil~. Otherwise merge all tables into one new table numbered above them all, without tombstones, and delete the old ones. The memtable is untouched.

Reads: the memtable, then tables newest first. On open, tables are read in order; any short table, bad checksum, or keys not strictly increasing returns ~ErrLedgerCorrupt~. Leftover ~.tmp~ files are removed and ignored (bonus star). Leave ~main~ empty.
`,
			Hints: []string{"An explicit marker makes deletion a readable state.", "Separate flushing the current state from merging old history.", "Review the inherited ~Ledger~ fields and ~Flush~ and ~Compact~ signatures in the starter."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch343Hidden},
	})
}

const ch343Starter = `package main

import "errors"

var ErrLedgerKey = errors.New("invalid ledger key")
var ErrLedgerValue = errors.New("ledger value too large")
var ErrLedgerCorrupt = errors.New("corrupt ledger")

type Ledger struct {
	mem map[string][]byte
}

func OpenLedger(dir string) (*Ledger, error) {
	return nil, errors.New("TODO")
}

func (l *Ledger) Set(key string, value []byte) error { return errors.New("TODO") }
func (l *Ledger) Remove(key string) error            { return errors.New("TODO") }
func (l *Ledger) Lookup(key string) ([]byte, bool)   { return nil, false }
func (l *Ledger) Flush() error                       { return errors.New("TODO") }
func (l *Ledger) Compact() error                     { return errors.New("TODO") }
func (l *Ledger) Close() error                       { return nil }

func main() {}
`

const ch343Hidden = `package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func sstRecord(op byte, key string, value []byte) []byte {
	r := make([]byte, 13, 13+len(key)+len(value))
	r[4] = op
	binary.BigEndian.PutUint32(r[5:9], uint32(len(key)))
	binary.BigEndian.PutUint32(r[9:13], uint32(len(value)))
	r = append(r, key...)
	r = append(r, value...)
	binary.BigEndian.PutUint32(r[0:4], crc32.ChecksumIEEE(r[4:]))
	return r
}

func sstOpen(t *testing.T, dir string) *Ledger {
	t.Helper()
	l, err := OpenLedger(dir)
	if err != nil {
		t.Fatal("open", err)
	}
	return l
}

func sstFiles(t *testing.T, dir string) []string {
	t.Helper()
	var out []string
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".sst") {
			out = append(out, e.Name())
		}
	}
	return out
}

func sstCheck(t *testing.T, l *Ledger, model map[string][]byte, keys []string) {
	t.Helper()
	for _, k := range keys {
		want, wok := model[k]
		got, ok := l.Lookup(k)
		if ok != wok || !bytes.Equal(got, want) {
			t.Fatalf("lookup %q = %q %v, want %q %v", k, got, ok, want, wok)
		}
	}
}

func TestHiddenCompactionFlushFormat(t *testing.T) {
	dir := t.TempDir()
	l := sstOpen(t, dir)
	if err := l.Flush(); err != nil || len(sstFiles(t, dir)) != 0 {
		t.Fatal("empty flush must create nothing", err)
	}
	if err := l.Compact(); err != nil {
		t.Fatal("compact without tables", err)
	}
	l.Set("oxygen", []byte("92"))
	l.Set("airlock", []byte("shut"))
	l.Remove("beacon")
	l.Set("oxygen", []byte("93"))
	if err := l.Flush(); err != nil {
		t.Fatal(err)
	}
	files := sstFiles(t, dir)
	if len(files) != 1 || files[0] != "sst-000001.sst" {
		t.Fatal("tables", files)
	}
	got, _ := os.ReadFile(filepath.Join(dir, files[0]))
	want := append(append(sstRecord(1, "airlock", []byte("shut")), sstRecord(2, "beacon", nil)...), sstRecord(1, "oxygen", []byte("93"))...)
	if !bytes.Equal(got, want) {
		t.Fatalf("table bytes\n got %x\nwant %x", got, want)
	}
	if fi, err := os.Stat(filepath.Join(dir, "ledger.wal")); err != nil || fi.Size() != 0 {
		t.Fatal("log must be empty after flush", err)
	}
	l.Remove("oxygen")
	l.Flush()
	if files := sstFiles(t, dir); len(files) != 2 || files[1] != "sst-000002.sst" {
		t.Fatal("second table", files)
	}
	if _, ok := l.Lookup("oxygen"); ok {
		t.Fatal("tombstone in newer table must hide older value")
	}
	if err := l.Compact(); err != nil {
		t.Fatal(err)
	}
	files = sstFiles(t, dir)
	if len(files) != 1 || files[0] != "sst-000003.sst" {
		t.Fatal("compacted tables", files)
	}
	got, _ = os.ReadFile(filepath.Join(dir, files[0]))
	if !bytes.Equal(got, sstRecord(1, "airlock", []byte("shut"))) {
		t.Fatalf("compacted table must drop tombstones and dead keys: %x", got)
	}
	l.Close()
	l = sstOpen(t, dir)
	defer l.Close()
	sstCheck(t, l, map[string][]byte{"airlock": []byte("shut")}, []string{"airlock", "beacon", "oxygen"})
}

// Random operations with flush, compaction and reopen: every lookup matches a model.
func TestHiddenCompactionModel(t *testing.T) {
	rng := rand.New(rand.NewSource(345))
	for round := 0; round < 4; round++ {
		dir := t.TempDir()
		l := sstOpen(t, dir)
		model := map[string][]byte{}
		var keys []string
		for i := 0; i < 15; i++ {
			keys = append(keys, fmt.Sprintf("bay-%02d", i))
		}
		for step := 0; step < 300; step++ {
			k := keys[rng.Intn(len(keys))]
			switch r := rng.Intn(20); {
			case r < 10:
				v := []byte(fmt.Sprint(step))
				if err := l.Set(k, v); err != nil {
					t.Fatal(err)
				}
				model[k] = v
			case r < 15:
				if err := l.Remove(k); err != nil {
					t.Fatal(err)
				}
				delete(model, k)
			case r < 17:
				if err := l.Flush(); err != nil {
					t.Fatal(err)
				}
			case r < 18:
				if err := l.Compact(); err != nil {
					t.Fatal(err)
				}
				if n := len(sstFiles(t, dir)); n > 1 {
					t.Fatal("compact left", n, "tables")
				}
			default:
				l.Close()
				l = sstOpen(t, dir)
			}
			sstCheck(t, l, model, keys)
		}
		l.Close()
	}
}

// A crash after the table rename but before the log reset replays harmlessly.
func TestHiddenCompactionCrashAfterFlush(t *testing.T) {
	dir := t.TempDir()
	l := sstOpen(t, dir)
	l.Set("a", []byte("1"))
	l.Flush()
	l.Set("a", []byte("2"))
	l.Remove("b")
	l.Set("c", []byte("3"))
	l.Close()
	wal, _ := os.ReadFile(filepath.Join(dir, "ledger.wal"))
	l = sstOpen(t, dir)
	l.Flush()
	l.Close()
	os.WriteFile(filepath.Join(dir, "ledger.wal"), wal, 0o644)
	l = sstOpen(t, dir)
	defer l.Close()
	sstCheck(t, l, map[string][]byte{"a": []byte("2"), "c": []byte("3")}, []string{"a", "b", "c"})
}

func TestHiddenCompactionCorruptTable(t *testing.T) {
	rng := rand.New(rand.NewSource(346))
	table := append(append(sstRecord(1, "alpha", []byte("one")), sstRecord(2, "beta", nil)...), sstRecord(1, "gamma", []byte("three"))...)
	var cases [][]byte
	for i := 0; i < 8; i++ {
		bad := append([]byte{}, table...)
		bad[[]int{0, 2, 13, 15}[rng.Intn(4)]+[]int{0, 21}[rng.Intn(2)]] ^= byte(1 + rng.Intn(255))
		cases = append(cases, bad)
		if n := 1 + rng.Intn(len(table)-1); n != 21 && n != 38 {
			cases = append(cases, table[:n])
		}
	}
	cases = append(cases, append(sstRecord(1, "zulu", nil), sstRecord(1, "alpha", nil)...))
	cases = append(cases, append(sstRecord(1, "same", nil), sstRecord(1, "same", nil)...))
	for _, c := range cases {
		dir := t.TempDir()
		os.WriteFile(filepath.Join(dir, "sst-000001.sst"), c, 0o644)
		if _, err := OpenLedger(dir); !errors.Is(err, ErrLedgerCorrupt) {
			t.Fatalf("table %x: %v", c, err)
		}
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "sst-000004.sst"), table, 0o644)
	l := sstOpen(t, dir)
	sstCheck(t, l, map[string][]byte{"alpha": []byte("one"), "gamma": []byte("three")}, []string{"alpha", "beta", "gamma"})
	l.Set("delta", []byte("4"))
	l.Flush()
	l.Close()
	if files := sstFiles(t, dir); len(files) != 2 || files[1] != "sst-000005.sst" {
		t.Fatal("next table number", files)
	}
}

func TestHiddenBonusCompactionStaleTemp(t *testing.T) {
	dir := t.TempDir()
	l := sstOpen(t, dir)
	l.Set("hull", []byte("ok"))
	l.Flush()
	l.Close()
	stale := filepath.Join(dir, "sst-000002.sst.tmp")
	os.WriteFile(stale, sstRecord(1, "hull", []byte("half"))[:9], 0o644)
	l = sstOpen(t, dir)
	defer l.Close()
	if _, err := os.Stat(stale); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("stale temp file must be removed", err)
	}
	sstCheck(t, l, map[string][]byte{"hull": []byte("ok")}, []string{"hull"})
	var names []string
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		names = append(names, e.Name())
	}
	sort.Strings(names)
	if strings.Join(names, ",") != "ledger.wal,sst-000001.sst" {
		t.Fatal("directory", names)
	}
}
`

func init() {
	add(&Level{ID: "34-4-chunk-store", Chapter: 34, Kind: "lab", Title: "تخزين موزّع", Goal: "وزّع الشحنات على كبسولات تخزين بنسخ متكررة واصلح اللي يبوظ أو يضيع",
		LessonRevision: 1, Steps: ch34LessonArabic["34-4-chunk-store"],
		Requires: []string{"sel:sha256.Sum256"}, Bridge: []int{99}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch344Starter}},
		Lesson: `# مستودع الكبسولات

## قبل ما تكتب
- لو كبسولة تخزين وقعت، الملف لسه ممكن يتقري؟
- إزاي تعرف إن النسخة اللي رجعت سليمة من غير ما تقارنها بنسخة تانية؟

## التقطيع (chunks)
الملف الكبير بيتقسم لقطع ثابتة الحجم؛ آخر قطعة ممكن تبقى أقصر. ~io.ReadFull(r, buf)~ بيملا البافر كله، ولو الداتا خلصت في النص يرجع ~io.ErrUnexpectedEOF~ مع العدد اللي اتقرا، ولو مفيش ولا بايت يرجع ~io.EOF~.
~~~go
buf := make([]byte, 4)
n, err := io.ReadFull(r, buf)
if err == io.EOF { /* خلصنا */ }
part := buf[:n] // صالح كمان مع ErrUnexpectedEOF
~~~

## عنونة المحتوى
اسم القطعة هو بصمة محتواها. ~sha256.Sum256(b)~ بيرجع ~[32]byte~، و ~hex.EncodeToString(sum[:])~ بيحوّلها لـ 64 حرف. نفس المحتوى = نفس الاسم، فالقطعة المكررة تتخزن مرة واحدة (dedup). والقراية بتتأكد بنفسها: احسب البصمة للي رجع وقارنها بالاسم.
~~~go
sum := sha256.Sum256([]byte("ore"))
id := hex.EncodeToString(sum[:])
~~~

## التكرار (replication)
كل قطعة تتخزن على أكتر من كبسولة. المكان بيتحسب من الاسم نفسه، فمش محتاج فهرس: البداية ~binary.BigEndian.Uint32(sum[:4]) % uint32(len(pods))~، وبعدها الكبسولات اللي بعدها بالترتيب مع اللف.
~~~go
start := int(binary.BigEndian.Uint32(sum[:4]) % uint32(len(pods)))
for j := 0; j < replicas; j++ { p := pods[(start+j)%len(pods)]; _ = p }
~~~

## الكبسولة interface
الكبسولة ممكن تبقى ذاكرة أو فولدر أو جهاز تاني؛ المستودع يعرف بس الـ interface. التستات بتحقن كبسولات بتقع أو بترجع بايتات بايظة.
~~~go
type Pod interface {
 Store(id string, data []byte) error
 Fetch(id string) ([]byte, error)
}
~~~

## الإصلاح (repair)
الإصلاح بيلف على كل قطعة اتخزنت، يلاقي نسخة سليمة، ويكتبها تاني على أي كبسولة من أماكنها نسختها ناقصة أو بايظة. لو مفيش ولا نسخة سليمة، القطعة ضاعت: سجّل ده وكمّل الباقي.

## أخطاء شائعة
~~~go
data, err := pods[start].Fetch(id)
if err != nil { return err } // غلط: فيه نسخ تانية
w.Write(data)                // ومن غير ما تتأكد من البصمة
~~~

## المهمة
نفّذ ~NewDepot(pods []Pod, chunkSize, replicas int) (*Depot, error)~: ارجع ~ErrDepotConfig~ لو مفيش كبسولات أو فيه ~nil~، أو ~chunkSize~ خارج 1 لحد 1048576، أو ~replicas~ خارج 1 لحد عدد الكبسولات. والدوال:
- ~Save(r io.Reader) ([]string, error)~: قطّع، وارجع الأسماء بالترتيب (المكرر يتكرر في القايمة). القطعة الجديدة تتخزن بـ ~Store~ مرة على كل كبسولة من أماكنها، والقطعة اللي المستودع خزّنها قبل كده بنجاح مبتتخزنش تاني. أي خطأ قراية أو ~Store~ يرجع بـ ~errors.Is~ والقطعة دي متتحسبش متخزنة. مدخل فاضي يرجع قايمة فاضية.
- ~Load(ids []string, w io.Writer) error~: لكل اسم جرّب أماكنه بالترتيب؛ أول نسخة من غير خطأ وبصمتها مطابقة تتكتب. لو مفيش، أو الاسم مش 64 حرف hex صغير، ارجع خطأ ملفوف بـ ~ErrChunkLost~. أخطاء الكاتب ترجع زي ما هي.
- ~Repair() (int, error)~: لكل قطعة متخزنة اعمل الإصلاح، وارجع عدد مرات ~Store~ اللي نجحت. القطع الضايعة ترجع خطأ ملفوف بـ ~ErrChunkLost~ بعد ما الباقي يتصلح.

مش لازم أمان مع goroutines كتير. ~main~ فاضية.
`,
		Hints: []string{"الاسم المبني من المحتوى يسهّل اكتشاف التكرار والتلف.", "رتّب الخطوات من chunk إلى placement ثم verification وrepair.", "ابدأ من ~Pod.Fetch~ و~Pod.Store~ وتعريف ~Depot~ في الـstarter."},
		En: &LevelText{Title: "Distributed Storage", Goal: "Spread cargo across storage pods with replicas and repair what breaks or vanishes", LessonRevision: 1, Steps: ch34LessonEnglish["34-4-chunk-store"], Files: []LabFile{{Name: "main.go", Starter: ch344Starter}},
			Lesson: `# The Pod Depot

## Before you write
- If one storage pod fails, can the file still be read?
- How do you know a returned copy is intact without comparing it to another copy?

## Chunks
A large file is split into fixed-size chunks; the last may be shorter. ~io.ReadFull(r, buf)~ fills the whole buffer; if data ends midway it returns ~io.ErrUnexpectedEOF~ with the count read, and with no bytes at all it returns ~io.EOF~.
~~~go
buf := make([]byte, 4)
n, err := io.ReadFull(r, buf)
if err == io.EOF { /* done */ }
part := buf[:n] // also valid with ErrUnexpectedEOF
~~~

## Content addressing
A chunk's name is the fingerprint of its content. ~sha256.Sum256(b)~ returns a ~[32]byte~, and ~hex.EncodeToString(sum[:])~ turns it into 64 characters. Same content means same name, so a duplicate chunk is stored once (dedup). Reads verify themselves: hash what came back and compare it with the name.
~~~go
sum := sha256.Sum256([]byte("ore"))
id := hex.EncodeToString(sum[:])
~~~

## Replication
Each chunk is stored on several pods. Placement is computed from the name itself, so no index is needed: start at ~binary.BigEndian.Uint32(sum[:4]) % uint32(len(pods))~, then the following pods in order, wrapping around.
~~~go
start := int(binary.BigEndian.Uint32(sum[:4]) % uint32(len(pods)))
for j := 0; j < replicas; j++ { p := pods[(start+j)%len(pods)]; _ = p }
~~~

## The pod interface
A pod might be memory, a directory, or another machine; the depot only knows the interface. Tests inject pods that fail or return damaged bytes.
~~~go
type Pod interface {
 Store(id string, data []byte) error
 Fetch(id string) ([]byte, error)
}
~~~

## Repair
Repair walks every stored chunk, finds an intact copy, and rewrites it to every placement pod whose copy is missing or damaged. If no intact copy exists the chunk is lost: record that and carry on with the rest.

## Common mistakes
~~~go
data, err := pods[start].Fetch(id)
if err != nil { return err } // Wrong: other replicas exist
w.Write(data)                // and the fingerprint is never checked
~~~

## Your mission
Implement ~NewDepot(pods []Pod, chunkSize, replicas int) (*Depot, error)~: return ~ErrDepotConfig~ for no pods or a ~nil~ pod, ~chunkSize~ outside 1 through 1048576, or ~replicas~ outside 1 through the pod count. And:
- ~Save(r io.Reader) ([]string, error)~: chunk the input and return names in order (duplicates repeat in the list). A new chunk is stored with ~Store~ once on each of its placement pods; a chunk this depot already stored successfully is not stored again. Any read or ~Store~ error is returned for ~errors.Is~ and that chunk does not count as stored. Empty input returns an empty list.
- ~Load(ids []string, w io.Writer) error~: for each name try its placement pods in order; the first copy without error whose fingerprint matches is written. If none, or the name is not 64 lowercase hex characters, return an error wrapping ~ErrChunkLost~. Writer errors are returned as is.
- ~Repair() (int, error)~: repair every stored chunk and return the number of successful ~Store~ calls. Lost chunks return an error wrapping ~ErrChunkLost~ after the rest are repaired.

Concurrent use need not be safe. Leave ~main~ empty.
`,
			Hints: []string{"A content-derived name helps detect duplicates and damage.", "Think from chunking to placement, then verification and repair.", "Start at ~Pod.Fetch~, ~Pod.Store~, and the ~Depot~ type in the starter."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch344Hidden},
	})
}

const ch344Starter = `package main

import (
	"errors"
	"io"
)

var ErrDepotConfig = errors.New("invalid depot config")
var ErrChunkLost = errors.New("chunk lost")

type Pod interface {
	Store(id string, data []byte) error
	Fetch(id string) ([]byte, error)
}

type Depot struct {
	pods []Pod
}

func NewDepot(pods []Pod, chunkSize, replicas int) (*Depot, error) {
	return nil, errors.New("TODO")
}

func (d *Depot) Save(r io.Reader) ([]string, error)    { return nil, errors.New("TODO") }
func (d *Depot) Load(ids []string, w io.Writer) error { return errors.New("TODO") }
func (d *Depot) Repair() (int, error)                 { return 0, errors.New("TODO") }

func main() {}
`

const ch344Hidden = `package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"testing/iotest"
)

var errPodDown = errors.New("pod down")

type memPod struct {
	data           map[string][]byte
	stores         int
	down, failSave bool
}

func newMemPod() *memPod { return &memPod{data: map[string][]byte{}} }

func (p *memPod) Store(id string, b []byte) error {
	if p.down || p.failSave {
		return fmt.Errorf("store: %w", errPodDown)
	}
	p.stores++
	p.data[id] = append([]byte{}, b...)
	return nil
}
func (p *memPod) Fetch(id string) ([]byte, error) {
	if p.down {
		return nil, errPodDown
	}
	b, ok := p.data[id]
	if !ok {
		return nil, errors.New("no such chunk")
	}
	return append([]byte{}, b...), nil
}

func depotPods(n int) ([]*memPod, []Pod) {
	var mem []*memPod
	var pods []Pod
	for i := 0; i < n; i++ {
		p := newMemPod()
		mem = append(mem, p)
		pods = append(pods, p)
	}
	return mem, pods
}

func depotPlacement(id string, n, replicas int) []int {
	b, _ := hex.DecodeString(id)
	start := int(binary.BigEndian.Uint32(b[:4]) % uint32(n))
	var out []int
	for j := 0; j < replicas; j++ {
		out = append(out, (start+j)%n)
	}
	return out
}

func newTestDepot(t *testing.T, pods []Pod, size, replicas int) *Depot {
	t.Helper()
	d, err := NewDepot(pods, size, replicas)
	if err != nil {
		t.Fatal(err)
	}
	return d
}

func TestHiddenDepotConfig(t *testing.T) {
	_, pods := depotPods(3)
	for _, c := range []struct {
		pods      []Pod
		size, rep int
	}{{nil, 4, 1}, {[]Pod{pods[0], nil}, 4, 1}, {pods, 0, 1}, {pods, 1048577, 1}, {pods, 4, 0}, {pods, 4, 4}} {
		if _, err := NewDepot(c.pods, c.size, c.rep); !errors.Is(err, ErrDepotConfig) {
			t.Fatal("config", c.size, c.rep, err)
		}
	}
	if _, err := NewDepot(pods, 1048576, 3); err != nil {
		t.Fatal("limits are inclusive", err)
	}
}

// Chunking, ids, placement, round trip and dedup against random cargo.
func TestHiddenDepotSaveLoad(t *testing.T) {
	rng := rand.New(rand.NewSource(347))
	for round := 0; round < 30; round++ {
		n := 1 + rng.Intn(5)
		rep := 1 + rng.Intn(n)
		size := 1 + rng.Intn(9)
		mem, pods := depotPods(n)
		d := newTestDepot(t, pods, size, rep)
		blocks := [][]byte{[]byte("ore"), []byte("ice"), []byte("gas")}
		var cargo []byte
		for i := rng.Intn(30); i > 0; i-- {
			cargo = append(cargo, blocks[rng.Intn(3)]...)
		}
		ids, err := d.Save(iotest.OneByteReader(bytes.NewReader(cargo)))
		if err != nil {
			t.Fatal(err)
		}
		if want := (len(cargo) + size - 1) / size; len(ids) != want {
			t.Fatal("chunk count", len(ids), want)
		}
		unique := map[string]bool{}
		for i, id := range ids {
			end := (i + 1) * size
			if end > len(cargo) {
				end = len(cargo)
			}
			sum := sha256.Sum256(cargo[i*size : end])
			if id != hex.EncodeToString(sum[:]) {
				t.Fatal("id is not the sha256 of the chunk", i)
			}
			unique[id] = true
			for _, p := range depotPlacement(id, n, rep) {
				if _, ok := mem[p].data[id]; !ok {
					t.Fatal("chunk not on its placement pod", p)
				}
			}
		}
		stores := 0
		for _, p := range mem {
			stores += p.stores
		}
		if stores != len(unique)*rep {
			t.Fatal("dedup: stores", stores, len(unique)*rep)
		}
		again, err := d.Save(bytes.NewReader(cargo))
		if err != nil || strings.Join(again, ",") != strings.Join(ids, ",") {
			t.Fatal("second save", err)
		}
		total := 0
		for _, p := range mem {
			total += p.stores
		}
		if total != stores {
			t.Fatal("known chunks must not be stored again")
		}
		var out bytes.Buffer
		if err := d.Load(ids, &out); err != nil || !bytes.Equal(out.Bytes(), cargo) {
			t.Fatal("load", err)
		}
	}
}

func TestHiddenDepotFailover(t *testing.T) {
	mem, pods := depotPods(4)
	d := newTestDepot(t, pods, 3, 3)
	cargo := []byte("reactor-core-schematics-v7")
	ids, err := d.Save(bytes.NewReader(cargo))
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range ids {
		pl := depotPlacement(id, 4, 3)
		mem[pl[0]].data[id][0] ^= 0xff // damaged copy on the first replica
		delete(mem[pl[1]].data, id)      // missing copy on the second
	}
	var out bytes.Buffer
	if err := d.Load(ids, &out); err != nil || !bytes.Equal(out.Bytes(), cargo) {
		t.Fatal("load must skip damaged and missing copies", err)
	}
	for _, p := range mem {
		p.down = true
	}
	if err := d.Load(ids, new(bytes.Buffer)); !errors.Is(err, ErrChunkLost) {
		t.Fatal("all pods down", err)
	}
	for _, bad := range []string{"", "xyz", strings.ToUpper(ids[0]), ids[0][:63]} {
		if err := d.Load([]string{bad}, new(bytes.Buffer)); !errors.Is(err, ErrChunkLost) {
			t.Fatalf("id %q: %v", bad, err)
		}
	}
}

type failDepotWriter struct{}

func (failDepotWriter) Write(p []byte) (int, error) { return 0, errPodDown }

func TestHiddenDepotErrors(t *testing.T) {
	mem, pods := depotPods(2)
	d := newTestDepot(t, pods, 4, 2)
	readErr := errors.New("antenna glitch")
	if _, err := d.Save(iotest.ErrReader(readErr)); !errors.Is(err, readErr) {
		t.Fatal("read error", err)
	}
	if ids, err := d.Save(bytes.NewReader(nil)); err != nil || len(ids) != 0 {
		t.Fatal("empty input", ids, err)
	}
	mem[1].failSave = true
	if _, err := d.Save(bytes.NewReader([]byte("abcd"))); !errors.Is(err, errPodDown) {
		t.Fatal("store error", err)
	}
	mem[1].failSave = false
	ids, err := d.Save(bytes.NewReader([]byte("abcd")))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := mem[1].data[ids[0]]; !ok {
		t.Fatal("a failed chunk must be stored again on the next save")
	}
	if err := d.Load(ids, failDepotWriter{}); !errors.Is(err, errPodDown) {
		t.Fatal("writer error", err)
	}
}

func TestHiddenDepotRepair(t *testing.T) {
	rng := rand.New(rand.NewSource(348))
	mem, pods := depotPods(5)
	d := newTestDepot(t, pods, 5, 3)
	cargo := make([]byte, 200)
	rng.Read(cargo)
	ids, err := d.Save(bytes.NewReader(cargo))
	if err != nil {
		t.Fatal(err)
	}
	if n, err := d.Repair(); n != 0 || err != nil {
		t.Fatal("healthy depot needs no repair", n, err)
	}
	wiped := rng.Intn(5)
	broken := 0
	for id := range mem[wiped].data {
		delete(mem[wiped].data, id) // node lost and replaced empty
		broken++
	}
	for _, id := range ids {
		pl := depotPlacement(id, 5, 3)
		if p := pl[2]; p != wiped {
			mem[p].data[id][1] ^= 0x5a // damaged replica
			broken++
		}
	}
	n, err := d.Repair()
	if err != nil || n != broken {
		t.Fatal("repair count", n, broken, err)
	}
	for _, id := range ids {
		for _, p := range depotPlacement(id, 5, 3) {
			b := mem[p].data[id]
			if sum := sha256.Sum256(b); hex.EncodeToString(sum[:]) != id {
				t.Fatal("replica still damaged after repair")
			}
		}
	}
	lost := ids[0]
	for _, p := range depotPlacement(lost, 5, 3) {
		mem[p].data[lost] = []byte("garbage")
	}
	n, err = d.Repair()
	if !errors.Is(err, ErrChunkLost) || n != 0 {
		t.Fatal("lost chunk", n, err)
	}
	delete(mem[depotPlacement(ids[1], 5, 3)[0]].data, ids[1])
	if n, err := d.Repair(); n != 1 || !errors.Is(err, ErrChunkLost) {
		t.Fatal("repair must continue past a lost chunk", n, err)
	}
}
`
