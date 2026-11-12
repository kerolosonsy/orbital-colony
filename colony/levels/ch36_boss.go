package levels

func init() {
	add(&Level{
		ID: "36-boss-storm-control", Chapter: 36, Kind: "lab", Boss: true, Setup: seededLabRoom, Race: true,
		Title: "Boss: درع العاصفة", Goal: "استرجع أوامر الدرع وانشرها بأمان واعرض الجاهزية والمقاييس",
		Requires: []string{"mutex", "iface", "sel:sha256.Sum256", "sel:slog.Logger"}, Bridge: []int{87, 88, 95, 96, 98, 100},
		Deps: []string{"modernc.org/sqlite"},
		Files: []LabFile{{Name: "main.go", Starter: `package main
import ("context";"database/sql";"errors";"log/slog";"net/http")
type StormEvent struct { ID string; Delta int; Version int }
type StormBroker interface { Publish(context.Context,StormEvent) error }
var ErrStorm=errors.New("invalid or conflicting storm command")
var ErrStormCorrupt=errors.New("corrupt storm command")
type StormControl struct { db *sql.DB; broker StormBroker; log *slog.Logger }
func NewStormControl(db *sql.DB, broker StormBroker, log *slog.Logger)*StormControl{return &StormControl{db,broker,log}}
func(s *StormControl)Init(ctx context.Context)error{return ErrStorm}
func(s *StormControl)Record(ctx context.Context,id string,delta int)error{return ErrStorm}
func(s *StormControl)Flush(ctx context.Context)error{return ErrStorm}
func(s *StormControl)ServeHTTP(w http.ResponseWriter,r *http.Request){http.Error(w,"unavailable",503)}
func main(){}
`}},
		Lesson: `# Boss: درع العاصفة

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 34-2-wal, 35-4-outbox, 35-5-kafka, 36-1-metrics, 36-3-ops.

مثال أصغر ومختلف:
~~~go
if err := publisher.Publish(ctx, event); err != nil {
 logger.Error("relay unavailable") // no credentials from raw errors
 return err
}
~~~

## أخطاء شائعة
~~~go
logger.Error("broker failed", "error", err) // may expose broker credentials
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~NewStormControl(db, broker, log)~ و~Init/Record/Flush/ServeHTTP~ بتوقيعات الـ starter. المحقونات مش nil؛ المنادي يملك SQLite والـ logger والـ broker. مفيش Kafka أو Docker: broker fake بمستهلك بيمنع تكرار الأثر عن طريق ID. Init قابلة للتكرار: ~storm_commands(seq INTEGER PRIMARY KEY AUTOINCREMENT, id TEXT NOT NULL UNIQUE, delta INTEGER NOT NULL, checksum TEXT NOT NULL, published INTEGER NOT NULL DEFAULT 0)~.

~Record(ctx, id, delta)~ ترفض ID فاضي أو delta صفر بـ ~ErrStorm~. ضيف الأمر الدائم وعلامة pending في نفس الصف بذرّية؛ checksum هي hex صغيرة لـ SHA256 للنص ~fmt.Sprintf("%s:%d", id, delta)~. مثال حدث: ~StormEvent{ID: "west", Delta: -2, Version: 1}~. نفس ID ونفس delta نجاح من غير تغيير؛ delta مختلفة ~ErrStorm~. Record متنشرش. البيانات والـ pending لازم يفضلوا بعد فتح كائن جديد على نفس القاعدة.

Flush متسلسلة مع Record على نفس الكائن: اقرا كل pending بترتيب seq، اقفل Rows، واتحقق من كل checksum وID وdelta قبل أول نشر؛ فساد يرجع ~ErrStormCorrupt~. انشر Version=1 بنفس السياق، وبعد نجاح كل حدث علّم published=1. أول فشل يوقف الباقي ويرجع السبب؛ الحدث الفاشل يفضل pending ويتعاد في Flush لاحقة. ده at-least-once؛ المستهلك لازم يمنع تكرار الأثر. لو broker فشل، سجّل Error برسالة ثابتة من غير raw error أو أسرار.

GET ~/livez~ يرجع 200 مستقل عن القاعدة. ~/readyz~ يرجع 200 لو القاعدة متاحة ومفيش pending، وإلا 503. ~/metrics~ يرجع نص بالظبط ~colony_storm_commands_total N~ ثم ~colony_storm_pending P~، كل سطر newline ومن غير labels؛ N عدد كل الصفوف وP الغير منشورة. فشل القاعدة 503 برسالة عامة من غير تفاصيل، وأي مسار تاني 404. مرّر السياق لكل I/O وامنع السباقات. سيب ~main~ فاضية. النجمة التالتة: فقدان ack ثم إعادة فتح القاعدة وإعادة Record وFlush مينشئش أمر جديد أو يكرر أثر المستهلك.
`,
		Hints: []string{"خزّن الأمر وعلامة النشر في نفس الصف؛ request ID فريد.", "اقفل Rows قبل استدعاء broker خصوصاً مع حوض اتصال واحد.", "علّم النجاح بعد Publish فقط؛ المقاييس ثابتة ومن غير ID كـ label."},
		En: &LevelText{Title: "Boss: Storm Shield Control", Goal: "Recover durable shield commands, retry delivery, and expose readiness and metrics", Files: []LabFile{{Name: "main.go", Starter: `package main
import ("context";"database/sql";"errors";"log/slog";"net/http")
type StormEvent struct { ID string; Delta int; Version int }
type StormBroker interface { Publish(context.Context,StormEvent) error }
var ErrStorm=errors.New("invalid or conflicting storm command")
var ErrStormCorrupt=errors.New("corrupt storm command")
type StormControl struct { db *sql.DB; broker StormBroker; log *slog.Logger }
func NewStormControl(db *sql.DB, broker StormBroker, log *slog.Logger)*StormControl{return &StormControl{db,broker,log}}
func(s *StormControl)Init(ctx context.Context)error{return ErrStorm}
func(s *StormControl)Record(ctx context.Context,id string,delta int)error{return ErrStorm}
func(s *StormControl)Flush(ctx context.Context)error{return ErrStorm}
func(s *StormControl)ServeHTTP(w http.ResponseWriter,r *http.Request){http.Error(w,"unavailable",503)}
func main(){}
`}},
			Lesson: `# Boss: Storm Shield Control

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 34-2-wal, 35-4-outbox, 35-5-kafka, 36-1-metrics, 36-3-ops.

A smaller, different example:
~~~go
if err := publisher.Publish(ctx, event); err != nil {
 logger.Error("relay unavailable") // no credentials from raw errors
 return err
}
~~~

## Common mistakes
~~~go
logger.Error("broker failed", "error", err) // may expose broker credentials
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~NewStormControl(db, broker, log)~ and ~Init/Record/Flush/ServeHTTP~ with the starter signatures. Dependencies are nonnil; the caller owns SQLite, logger, and broker. No Kafka or Docker: the fake broker has a consumer deduplicating effects by ID. Repeatable Init creates ~storm_commands(seq INTEGER PRIMARY KEY AUTOINCREMENT, id TEXT NOT NULL UNIQUE, delta INTEGER NOT NULL, checksum TEXT NOT NULL, published INTEGER NOT NULL DEFAULT 0)~.

~Record(ctx, id, delta)~ rejects empty ID or zero delta with ~ErrStorm~. Atomically insert the durable command and pending marker in one row; checksum is lowercase hex SHA256 of ~fmt.Sprintf("%s:%d", id, delta)~. Example event: ~StormEvent{ID: "west", Delta: -2, Version: 1}~. Same ID and delta succeeds without change; a conflicting delta returns ~ErrStorm~. Record never publishes. Data and pending work survive a new control instance on the same database.

Flush serializes with Record on one instance: read all pending rows by seq, close Rows, and validate every checksum, ID, and delta before publishing anything; corruption returns ~ErrStormCorrupt~. Publish Version=1 using the same context, marking published=1 after each success. Stop on the first failure and preserve its cause; the failed event stays pending for the next Flush. Delivery is at-least-once; consumers deduplicate effects. Log broker failures at Error level with a fixed message, omitting raw errors and secrets.

GET ~/livez~ returns 200 independently of storage. ~/readyz~ returns 200 only when storage is available with no pending work, otherwise 503. ~/metrics~ returns exactly ~colony_storm_commands_total N~ then ~colony_storm_pending P~, each followed by newline, without labels; N counts all rows, P unpublished rows. Database failures return 503 with a generic message; other paths return 404. Pass context to I/O and avoid races. Leave ~main~ empty. Third star: lost acknowledgement followed by reopening, repeated Record, and Flush creates no new command or duplicate consumer effect.
`,
			Hints: []string{"Store the command and delivery marker together with a unique ID.", "Close Rows before broker calls, especially with a one-connection pool.", "Acknowledge only after Publish succeeds; metrics stay fixed without ID labels."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: `package main
import (
    "bytes"
    "context"
    "database/sql"
    "crypto/sha256"
    "errors"
    "fmt"
    "log/slog"
    "math/rand"
    "net/http/httptest"
    "path/filepath"
    "reflect"
    "strings"
    "sync"
    "testing"
    _ "modernc.org/sqlite"
)
type hiddenStormBroker struct { events []StormEvent; effects map[string]int; failAt string; failure error; db *sql.DB }
func(b *hiddenStormBroker)Publish(ctx context.Context,e StormEvent)error{
    // Query through a one-connection pool: Flush must close Rows before publishing.
    var n int
    if err:=b.db.QueryRowContext(ctx,"SELECT count(*) FROM storm_commands WHERE id=?",e.ID).Scan(&n);err!=nil{return err}
    if n!=1{return errors.New("published before durable record")}
    b.events=append(b.events,e)
    // The fake consumer deduplicates effects by event ID even if acknowledgement fails.
    if _,ok:=b.effects[e.ID];!ok{b.effects[e.ID]=e.Delta}
    if e.ID==b.failAt{return b.failure}
    return nil
}
func hiddenStormOpen(t *testing.T,path string)(*StormControl,*hiddenStormBroker,*bytes.Buffer){
    t.Helper();db,err:=sql.Open("sqlite",path);if err!=nil{t.Fatal(err)};db.SetMaxOpenConns(1);t.Cleanup(func(){db.Close()})
    logs:=&bytes.Buffer{};b:=&hiddenStormBroker{effects:map[string]int{},db:db}
    s:=NewStormControl(db,b,slog.New(slog.NewJSONHandler(logs,nil)))
    if err=s.Init(context.Background());err!=nil{t.Fatal(err)}
    if err=s.Init(context.Background());err!=nil{t.Fatal("migration repeated",err)}
    return s,b,logs
}
func hiddenStormStatus(s *StormControl,path string)(int,string){w:=httptest.NewRecorder();s.ServeHTTP(w,httptest.NewRequest("GET",path,nil));return w.Code,w.Body.String()}
func hiddenStormMetrics(t *testing.T,s *StormControl,total,pending int){t.Helper();code,text:=hiddenStormStatus(s,"/metrics");want:=fmt.Sprintf("colony_storm_commands_total %d\ncolony_storm_pending %d\n",total,pending);if code!=200 || text!=want{t.Fatal("bounded metrics",code,text,want)}}
func TestHiddenStormDurabilityAndDelivery(t *testing.T){
    path:=filepath.Join(t.TempDir(),"storm.db");s,b,_:=hiddenStormOpen(t,path);ctx:=context.Background()
    if code,_:=hiddenStormStatus(s,"/readyz");code!=200{t.Fatal("fresh not ready",code)}
    rng:=rand.New(rand.NewSource(36091));var want []StormEvent
    for i:=0;i<15;i++{e:=StormEvent{fmt.Sprintf("shutter-%02d",i),1+rng.Intn(30),1};want=append(want,e);if err:=s.Record(ctx,e.ID,e.Delta);err!=nil{t.Fatal(err)}}
    if len(b.events)!=0{t.Fatal("Record must not publish")}
    hiddenStormMetrics(t,s,15,15)
    if code,_:=hiddenStormStatus(s,"/readyz");code!=503{t.Fatal("pending must block readiness",code)}
    if err:=s.db.Close();err!=nil{t.Fatal(err)}
    s,b,_=hiddenStormOpen(t,path)
    if err:=s.Flush(ctx);err!=nil{t.Fatal(err)}
    if !reflect.DeepEqual(b.events,want){t.Fatal("restart/order/version",b.events,want)}
    hiddenStormMetrics(t,s,15,0)
    if code,_:=hiddenStormStatus(s,"/readyz");code!=200{t.Fatal(code)}
    if err:=s.Flush(ctx);err!=nil || len(b.events)!=15{t.Fatal("repeated flush",err,len(b.events))}
}
func TestHiddenStormFailureRetryAndIdempotency(t *testing.T){
    s,b,logs:=hiddenStormOpen(t,filepath.Join(t.TempDir(),"storm.db"));ctx:=context.Background()
    for _,id:=range []string{"a","b","c"}{if err:=s.Record(ctx,id,4);err!=nil{t.Fatal(err)}}
    if err:=s.Record(ctx,"a",4);err!=nil{t.Fatal("identical replay",err)}
    for _,tc:=range []struct{id string;delta int}{{"a",5},{"",1},{"x",0}}{if err:=s.Record(ctx,tc.id,tc.delta);!errors.Is(err,ErrStorm){t.Fatal("invalid/conflict",err)}}
    failure:=errors.New("BROKER_PASSWORD=do-not-log");b.failure=failure;b.failAt="b"
    if err:=s.Flush(ctx);!errors.Is(err,failure){t.Fatal("broker cause lost",err)}
    hiddenStormMetrics(t,s,3,2)
    if len(b.events)!=2{t.Fatal("did not stop on failed publication",b.events)}
    b.failAt=""
    if err:=s.Flush(ctx);err!=nil{t.Fatal(err)}
    if len(b.events)!=4 || b.events[2].ID!="b" || b.events[3].ID!="c" || len(b.effects)!=3{t.Fatal("retry/dedup contract",b.events,b.effects)}
    if logs.Len()==0 || strings.Contains(logs.String(),"do-not-log") || strings.Contains(logs.String(),"BROKER_PASSWORD"){t.Fatal("missing/redaction log",logs.String())}
    canceled,cancel:=context.WithCancel(ctx);cancel()
    if err:=s.Record(canceled,"cancel",2);!errors.Is(err,context.Canceled){t.Fatal(err)}
    if code,_:=hiddenStormStatus(s,"/livez");code!=200{t.Fatal(code)}
    s.db.Close()
    if code,text:=hiddenStormStatus(s,"/readyz");code!=503 || strings.Contains(text,"sql:"){t.Fatal("database readiness/internal leakage",code,text)}
    if code,_:=hiddenStormStatus(s,"/livez");code!=200{t.Fatal("liveness depended on database",code)}
}
func TestHiddenStormCorruptionAndConcurrency(t *testing.T){
    s,b,_:=hiddenStormOpen(t,filepath.Join(t.TempDir(),"storm.db"));ctx:=context.Background()
    var wg sync.WaitGroup
    for i:=0;i<20;i++{wg.Add(1);go func(i int){defer wg.Done();if err:=s.Record(ctx,fmt.Sprint(i),1);err!=nil{t.Error(err)};hiddenStormStatus(s,"/metrics")}(i)}
    wg.Wait();hiddenStormMetrics(t,s,20,20)
    if _,err:=s.db.Exec("UPDATE storm_commands SET checksum='bad' WHERE id='19'");err!=nil{t.Fatal(err)}
    if err:=s.Flush(ctx);!errors.Is(err,ErrStormCorrupt) || len(b.events)!=0{t.Fatal("batch not validated before publication",err,b.events)}
    var checksum string
    if err:=s.db.QueryRow("SELECT checksum FROM storm_commands WHERE id='0'").Scan(&checksum);err!=nil{t.Fatal(err)}
    if checksum!=fmt.Sprintf("%x",sha256.Sum256([]byte("0:1"))){t.Fatal("checksum contract",checksum)}
}
func TestHiddenBonusStormLostAcknowledgementAfterRestart(t *testing.T){
    path:=filepath.Join(t.TempDir(),"storm.db");s,b,_:=hiddenStormOpen(t,path);ctx:=context.Background()
    if err:=s.Record(ctx,"shield",7);err!=nil{t.Fatal(err)}
    b.failAt="shield";b.failure=errors.New("ack lost")
    if err:=s.Flush(ctx);err==nil{t.Fatal("ack loss hidden")}
    effects:=b.effects;s.db.Close()
    s,b,_=hiddenStormOpen(t,path);b.effects=effects
    if err:=s.Record(ctx,"shield",7);err!=nil{t.Fatal("restart replay",err)}
    if err:=s.Flush(ctx);err!=nil{t.Fatal(err)}
    if len(effects)!=1 || effects["shield"]!=7 || len(b.events)!=1{t.Fatal("delivery/effect duplicate",effects,b.events)}
    hiddenStormMetrics(t,s,1,0)
}
`},
	})
}
