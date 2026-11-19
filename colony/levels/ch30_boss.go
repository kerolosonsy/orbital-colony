package levels

func init() {
	add(&Level{
		LessonRevision: 1, Steps: ch30Lesson30bArabic,
		ID: "30-boss-nursery", Chapter: 30, Kind: "lab", Boss: true, Setup: seededLabRoom, Race: true,
		Title: "Boss: موسم الزراعة", Goal: "احجز مؤن المشتل في transaction وحدّث الـcache بعد النجاح",
		Requires: []string{"iface", "sel:sql.ErrNoRows"}, Bridge: []int{61, 64, 65, 66, 68, 70},
		Deps: []string{"modernc.org/sqlite"},
		Files: []LabFile{{Name: "main.go", Starter: `package main
import ("context";"database/sql";"errors")
type SeedLine struct { Name string; Units int }
type NurseryCache interface { Invalidate(context.Context,string) error }
var ErrPlanting=errors.New("invalid planting")
var ErrStock=errors.New("insufficient stock")
var ErrReplay=errors.New("planting already booked")
var ErrCacheSync=errors.New("committed but cache invalidation failed")
func InitNursery(ctx context.Context, db *sql.DB) error { return ErrPlanting }
func BookPlanting(ctx context.Context, db *sql.DB, cache NurseryCache, request string, lines []SeedLine) error { return ErrPlanting }
func main() {}
`}},
		Lesson: `# Boss: موسم الزراعة

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 29-1-sqlite-crud, 29-2-migrations, 29-3-transfer, 29-4-pool, 30-1-cache-aside, 30-4-read-write-split.

مثال أصغر ومختلف:
~~~go
rows, err := db.QueryContext(ctx, "SELECT name FROM nursery_stock")
if err != nil { return err }
defer rows.Close()
~~~

## أخطاء شائعة
~~~go
cache.Invalidate(ctx, name) // wrong before commit: rollback must not affect cache
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~InitNursery(ctx, db) error~ و~BookPlanting(ctx, db, cache, request, lines) error~. التست بيحقن SQLite مؤقتة بحوض اتصال واحد و~NurseryCache~ fake في الذاكرة، من غير Docker؛ المنادي مالك db وcache. التهيئة قابلة للتكرار: جدول ~nursery_stock(name TEXT PRIMARY KEY, units INTEGER NOT NULL CHECK(units>=0))~ وجدول ~nursery_bookings(request TEXT PRIMARY KEY)~.

كل ~SeedLine{Name: "moss", Units: 3}~ حجز مورد. request مش فاضي، batch مش فاضية، الأسماء مش فاضية وفريدة، وUnits موجبة؛ غير كده ~ErrPlanting~ قبل I/O. متغيّرش lines. في transaction واحدة افحص الطلب المكرر (~ErrReplay~)، اخصم كل الموارد باستعلامات parameters مع شرط رصيد كافي، وسجّل الطلب ثم commit. مورد غايب أو ناقص يرجع ~ErrStock~؛ أي فشل يرجّع السبب ويلغي كل الخصومات وتسجيل الطلب. مرّر السياق لكل I/O وأقفل المعاملة حتى مع الخطأ.

بعد commit بس نادِ ~cache.Invalidate(ctx, name)~ بترتيب lines. مفيش قراية من الكاش هنا؛ القاعدة مصدر الحقيقة. نجاح الاختبارات العادية بيفترض كاش شغال. سيب ~main~ فاضية. النجمة التالتة: أول فشل invalidation يوقف الباقي ويرجع خطأ يطابق ~ErrCacheSync~ وسبب الكاش معاً؛ الحجز اتثبت ومبيتكررش عند إعادة الطلب. ده تحذير بعد commit، مش transaction موزعة؛ إصلاح الكاش مسؤولية المنادي.
`,
		Hints: []string{"الـtransaction تجمع تغييرات قاعدة البيانات كعملية واحدة.", "افحص الـrequest قبل ~BeginTx~، وبعدها افصل commit عن مزامنة الـcache.", "استخدم rollback المؤجل واحتفظ بسبب الـerror بعد commit."},
		En: &LevelText{LessonRevision: 1, Steps: ch30Lesson30bEnglish, Title: "Boss: Planting Season", Goal: "Reserve nursery supplies transactionally and invalidate cache after commit", Files: []LabFile{{Name: "main.go", Starter: `package main
import ("context";"database/sql";"errors")
type SeedLine struct { Name string; Units int }
type NurseryCache interface { Invalidate(context.Context,string) error }
var ErrPlanting=errors.New("invalid planting")
var ErrStock=errors.New("insufficient stock")
var ErrReplay=errors.New("planting already booked")
var ErrCacheSync=errors.New("committed but cache invalidation failed")
func InitNursery(ctx context.Context, db *sql.DB) error { return ErrPlanting }
func BookPlanting(ctx context.Context, db *sql.DB, cache NurseryCache, request string, lines []SeedLine) error { return ErrPlanting }
func main() {}
`}},
			Lesson: `# Boss: Planting Season

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 29-1-sqlite-crud, 29-2-migrations, 29-3-transfer, 29-4-pool, 30-1-cache-aside, 30-4-read-write-split.

A smaller, different example:
~~~go
rows, err := db.QueryContext(ctx, "SELECT name FROM nursery_stock")
if err != nil { return err }
defer rows.Close()
~~~

## Common mistakes
~~~go
cache.Invalidate(ctx, name) // wrong before commit: rollback must not affect cache
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~InitNursery(ctx, db) error~ and ~BookPlanting(ctx, db, cache, request, lines) error~. Tests inject temporary SQLite with a single-connection pool and an in-memory ~NurseryCache~ fake, without Docker. The caller owns db and cache. Initialization is repeatable: create ~nursery_stock(name TEXT PRIMARY KEY, units INTEGER NOT NULL CHECK(units>=0))~ and ~nursery_bookings(request TEXT PRIMARY KEY)~.

Each ~SeedLine{Name: "moss", Units: 3}~ reserves a resource. Require a nonempty request, nonempty batch, unique nonempty names, and positive Units; otherwise return ~ErrPlanting~ before I/O. Do not mutate lines. In one transaction check duplicate requests (~ErrReplay~), debit all resources using parameterized queries guarded by sufficient balance, insert the request, then commit. Missing or insufficient stock returns ~ErrStock~; any failure preserves its cause and rolls back every debit and request record. Pass context to I/O and always finish the transaction.

Only after commit call ~cache.Invalidate(ctx, name)~ in line order. This API never reads cache; SQL remains authoritative. Required tests assume healthy cache. Leave ~main~ empty. Third star: stop at the first invalidation failure and return an error matching both ~ErrCacheSync~ and the cache cause. The reservation remains committed and replay must never debit again. This is a post-commit warning, not a distributed transaction; the caller repairs cache.
`,
			Hints: []string{"A transaction groups database changes into one unit.", "Validate before ~BeginTx~, then separate commit from cache synchronization.", "Defer rollback and preserve the cause of a post-commit error."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: `package main
import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "math/rand"
    "path/filepath"
    "reflect"
    "sync"
    "testing"
    _ "modernc.org/sqlite"
)
type hiddenNurseryCache struct {mu sync.Mutex; names []string; failure error; db *sql.DB; t *testing.T}
func(c *hiddenNurseryCache)Invalidate(ctx context.Context,name string)error{
    // A single-connection pool proves the transaction has ended before cache I/O.
    if c.db!=nil {var n int;if err:=c.db.QueryRowContext(ctx,"SELECT count(*) FROM nursery_bookings").Scan(&n);err!=nil{c.t.Error(err)}else if n==0{c.t.Error("cache touched before commit")}}
    c.mu.Lock();defer c.mu.Unlock();c.names=append(c.names,name);return c.failure
}
func hiddenNursery(t *testing.T)(*sql.DB,*hiddenNurseryCache){
    t.Helper();db,err:=sql.Open("sqlite",filepath.Join(t.TempDir(),"nursery.db"));if err!=nil{t.Fatal(err)}
    db.SetMaxOpenConns(1);t.Cleanup(func(){db.Close()})
    if err=InitNursery(context.Background(),db);err!=nil{t.Fatal(err)}
    if err=InitNursery(context.Background(),db);err!=nil{t.Fatal("migration not repeatable",err)}
    return db,&hiddenNurseryCache{db:db,t:t}
}
func hiddenStock(t *testing.T,db *sql.DB,name string)int{t.Helper();var n int;if err:=db.QueryRow("SELECT units FROM nursery_stock WHERE name=?",name).Scan(&n);err!=nil{t.Fatal(err)};return n}
func hiddenNurserySeed(t *testing.T,db *sql.DB,a,b int){t.Helper();if _,err:=db.Exec("INSERT INTO nursery_stock(name,units) VALUES('algae',?),('binder',?)",a,b);err!=nil{t.Fatal(err)}}
func TestHiddenNurseryTransactionAndCache(t *testing.T){
    rng:=rand.New(rand.NewSource(30091))
    for trial:=0;trial<5;trial++{
        db,c:=hiddenNursery(t);n:=20+rng.Intn(100);hiddenNurserySeed(t,db,n,n)
        ctx:=context.Background();lines:=[]SeedLine{{"algae",3},{"binder",2}}
        before:=append([]SeedLine(nil),lines...)
        if err:=BookPlanting(ctx,db,c,fmt.Sprint(trial),lines);err!=nil{t.Fatal(err)}
        if hiddenStock(t,db,"algae")!=n-3 || hiddenStock(t,db,"binder")!=n-2 || !reflect.DeepEqual(lines,before){t.Fatal("wrong debit or mutated input")}
        if !reflect.DeepEqual(c.names,[]string{"algae","binder"}){t.Fatal("cache invalidation",c.names)}
        if err:=BookPlanting(ctx,db,c,fmt.Sprint(trial),lines);!errors.Is(err,ErrReplay){t.Fatal(err)}
        if hiddenStock(t,db,"algae")!=n-3{t.Fatal("replayed debit")}
    }
}
func TestHiddenNurseryRollbackAndBoundaries(t *testing.T){
    for _,tc:=range []struct{name string;request string;lines []SeedLine;want error}{
        {"empty request","",[]SeedLine{{"algae",1}},ErrPlanting},
        {"empty batch","r",nil,ErrPlanting},
        {"empty name","r",[]SeedLine{{"",1}},ErrPlanting},
        {"zero","r",[]SeedLine{{"algae",0}},ErrPlanting},
        {"negative","r",[]SeedLine{{"algae",-1}},ErrPlanting},
        {"duplicate","r",[]SeedLine{{"algae",1},{"algae",2}},ErrPlanting},
        {"missing","r",[]SeedLine{{"algae",1},{"absent",1}},ErrStock},
        {"insufficient","r",[]SeedLine{{"algae",1},{"binder",9}},ErrStock},
    }{t.Run(tc.name,func(t *testing.T){
        db,c:=hiddenNursery(t);hiddenNurserySeed(t,db,5,5)
        err:=BookPlanting(context.Background(),db,c,tc.request,tc.lines)
        if !errors.Is(err,tc.want){t.Fatal(err)}
        var booked int;if err=db.QueryRow("SELECT count(*) FROM nursery_bookings").Scan(&booked);err!=nil{t.Fatal(err)}
        if hiddenStock(t,db,"algae")!=5 || hiddenStock(t,db,"binder")!=5 || booked!=0 || len(c.names)!=0{t.Fatal("partial reservation/cache invalidation")}
    })}
    db,c:=hiddenNursery(t);hiddenNurserySeed(t,db,5,5)
    if _,err:=db.Exec("CREATE TRIGGER reject_binder BEFORE UPDATE ON nursery_stock WHEN OLD.name='binder' BEGIN SELECT RAISE(ABORT,'disk simulation'); END");err!=nil{t.Fatal(err)}
    if err:=BookPlanting(context.Background(),db,c,"fault",[]SeedLine{{"algae",1},{"binder",1}});err==nil{t.Fatal("database failure hidden")}
    if hiddenStock(t,db,"algae")!=5 || len(c.names)!=0{t.Fatal("SQL failure did not roll back")}
    ctx,cancel:=context.WithCancel(context.Background());cancel()
    if err:=BookPlanting(ctx,db,c,"cancel",[]SeedLine{{"algae",1}});!errors.Is(err,context.Canceled){t.Fatal(err)}
    // Real constraints are part of the migration contract.
    if _,err:=db.Exec("INSERT INTO nursery_stock VALUES('bad',-1)");err==nil{t.Fatal("negative stock allowed")}
}
func TestHiddenNurseryConcurrentReplay(t *testing.T){
    db,c:=hiddenNursery(t);hiddenNurserySeed(t,db,7,7)
    results:=make(chan error,12);var wg sync.WaitGroup
    for i:=0;i<12;i++{wg.Add(1);go func(){defer wg.Done();results<-BookPlanting(context.Background(),db,c,"same",[]SeedLine{{"algae",7}})}()}
    wg.Wait();close(results);success:=0
    for err:=range results{if err==nil{success++}else if !errors.Is(err,ErrReplay){t.Fatal(err)}}
    if success!=1 || hiddenStock(t,db,"algae")!=0{t.Fatal("atomic request/zero boundary",success)}
}
func TestHiddenBonusNurseryCommittedWarning(t *testing.T){
    db,c:=hiddenNursery(t);hiddenNurserySeed(t,db,5,5)
    failure:=errors.New("cache offline");c.failure=failure
    err:=BookPlanting(context.Background(),db,c,"committed",[]SeedLine{{"algae",2}})
    if !errors.Is(err,ErrCacheSync) || !errors.Is(err,failure) || hiddenStock(t,db,"algae")!=3{t.Fatal("lost committed warning",err)}
    c.failure=nil
    if err=BookPlanting(context.Background(),db,c,"committed",[]SeedLine{{"algae",2}});!errors.Is(err,ErrReplay){t.Fatal(err)}
    if hiddenStock(t,db,"algae")!=3{t.Fatal("warning caused double booking")}
}
`},
	})
}
