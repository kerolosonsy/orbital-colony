package levels

func init() {
	add(&Level{
		ID: "25-boss-refinery", Chapter: 25, Kind: "lab", Boss: true, Setup: seededLabRoom, Race: true,
		Title: "Boss: مصفاة الطوارئ", Goal: "شغّل مصفاة محدودة بتقف نضيف عند الخطأ والإلغاء",
		Requires: []string{"go", "chan", "select", "waitgroup", "mutex"}, Bridge: []int{34, 37, 41, 44},
		Files: []LabFile{{Name: "main.go", Starter: `package main
import (
    "context"
    "errors"
)
var ErrCrew = errors.New("crew must be positive")
func Refine(ctx context.Context, ore []int, crew int, assay func(context.Context,int)(int,error)) ([]int,error) {
    return nil, ErrCrew
}
func main() {}
`}},
		Lesson: `# Boss: مصفاة الطوارئ

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 13-4-chain-of-command, 24-3-work-crew, 24-4-pipeline, 25-1-fuel-account, 25-5-ttl-cache.

مثال أصغر ومختلف:
~~~go
jobs := make(chan int)
go func() { defer close(jobs); jobs <- 7 }()
for n := range jobs { _ = n }
~~~

## أخطاء شائعة
~~~go
go func() { results <- value }() // blocks forever if nobody receives
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~Refine(ctx, ore, crew, assay) ([]int, error)~ بتوقيع الـ starter. ~ore~ قراءات، ~crew~ حد العمال، و~assay~ دالة محقونة بتحترم السياق وبترجع نتيجة أو خطأ. مثال: ~Refine(ctx, []int{2, 4}, 2, assay)~. قبل الشغل ارفض ~crew <= 0~ بـ ~ErrCrew~؛ بعدها افحص السياق. شغّل لحد ~min(crew, len(ore))~ نداءات بالتوازي لما فيه شغل، كل عنصر مرة واحدة والنتائج بنفس ترتيب المدخل، من غير تغييره. الفاضي يرجع slice فاضية مش nil.

أول خطأ يلغي سياق العمال ويرجع نفس السبب بـ ~nil~ بدل النتيجة الجزئية. إلغاء الأب يرجع خطأ السياق. لازم كل العمال يخلصوا قبل الرجوع، وكل channel تتقفل من مالكها مرة واحدة. التستات بتستخدم ~-race~ و~labkit.LeakCheck~. سيب ~main~ فاضية. النجمة التالتة: نداءات مصافي مستقلة بالتوازي من غير حالة مشتركة.
`,
		Hints: []string{"قناة الشغل يملكها المرسل؛ اقفلها بعد الإرسال.", "نتائج بمؤشرات ثابتة تمنع خلط الترتيب؛ احمِ أول خطأ بقفل.", "كل إرسال لازم يلاحظ الإلغاء، وبعده Wait قبل الرجوع."},
		En: &LevelText{Title: "Boss: Emergency Refinery", Goal: "Run a bounded refinery that joins on failure and cancellation", Files: []LabFile{{Name: "main.go", Starter: `package main
import (
    "context"
    "errors"
)
var ErrCrew = errors.New("crew must be positive")
func Refine(ctx context.Context, ore []int, crew int, assay func(context.Context,int)(int,error)) ([]int,error) {
    return nil, ErrCrew
}
func main() {}
`}},
			Lesson: `# Boss: Emergency Refinery

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 13-4-chain-of-command, 24-3-work-crew, 24-4-pipeline, 25-1-fuel-account, 25-5-ttl-cache.

A smaller, different example:
~~~go
jobs := make(chan int)
go func() { defer close(jobs); jobs <- 7 }()
for n := range jobs { _ = n }
~~~

## Common mistakes
~~~go
go func() { results <- value }() // blocks forever if nobody receives
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~Refine(ctx, ore, crew, assay) ([]int, error)~ with the starter signature. ~ore~ contains readings, ~crew~ bounds workers, and the injected ~assay~ honors its context and returns a result or error. Example: ~Refine(ctx, []int{2, 4}, 2, assay)~. Reject ~crew <= 0~ with ~ErrCrew~ before work, then check context. Run up to ~min(crew, len(ore))~ callbacks concurrently when work exists, process each element once, preserve input order in output, and do not mutate input. Empty input returns a nonnil empty slice.

The first failure cancels the workers' context and returns its cause with nil output. Parent cancellation returns its context error. Join every worker before returning; the owner closes each channel once. Tests use ~-race~ and ~labkit.LeakCheck~. Leave ~main~ empty. Third star: independent concurrent refinery calls share no state.
`,
			Hints: []string{"The sender owns and closes the job channel.", "Index result slots to preserve order; lock the first error.", "Every send must observe cancellation, then Wait before returning."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: `package main
import (
    "context"
    "errors"
    "math/rand"
    "reflect"
    "sync"
    "sync/atomic"
    "testing"
    "time"
    "orbital/colony/labkit"
)
func TestHiddenRefineryResults(t *testing.T) {
    labkit.LeakCheck(t)
    rng:=rand.New(rand.NewSource(25091))
    for _, n:=range []int{0,1,2,67} {
        in:=make([]int,n); want:=make([]int,n)
        for i:=range in { in[i]=rng.Intn(10000); want[i]=in[i]*3 }
        before:=append([]int{},in...)
        for _, crew:=range []int{1,4,99} {
            got,err:=Refine(context.Background(),in,crew,func(ctx context.Context,x int)(int,error){return x*3,nil})
            if err!=nil || !reflect.DeepEqual(got,want) { t.Fatal(n,crew,got,err) }
        }
        if !reflect.DeepEqual(in,before) { t.Fatal("input changed") }
    }
    for _, crew:=range []int{0,-1} {
        got,err:=Refine(context.Background(),[]int{1},crew,func(context.Context,int)(int,error){t.Error("invalid crew ran work");return 0,nil})
        if got!=nil || !errors.Is(err,ErrCrew) { t.Fatal(got,err) }
    }
}
func TestHiddenRefineryBoundAndCancellation(t *testing.T) {
    labkit.LeakCheck(t)
    ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
    defer cancel()
    started:=make(chan struct{},20)
    release:=make(chan struct{})
    var active,peak atomic.Int64
    done:=make(chan error,1)
    go func(){
        _,err:=Refine(ctx,make([]int,20),3,func(ctx context.Context,x int)(int,error){
            a:=active.Add(1)
            for p:=peak.Load(); a>p && !peak.CompareAndSwap(p,a); p=peak.Load() {}
            defer active.Add(-1)
            started<-struct{}{}
            select { case <-release: return x,nil; case <-ctx.Done(): return 0,ctx.Err() }
        })
        done<-err
    }()
    for i:=0;i<3;i++ { select {case <-started: case <-ctx.Done():t.Fatal("crew did not run concurrently") } }
    close(release)
    select {case err:=<-done: if err!=nil {t.Fatal(err)};case <-ctx.Done():t.Fatal("workers did not finish")}
    if peak.Load()>3 || active.Load()!=0 {t.Fatal("crew limit/join",peak.Load(),active.Load())}

    stopCtx,stop:=context.WithCancel(context.Background())
    entered:=make(chan struct{},4)
    go func(){
        _,err:=Refine(stopCtx,make([]int,50),4,func(ctx context.Context,_ int)(int,error){entered<-struct{}{}; <-ctx.Done();return 0,ctx.Err()})
        done<-err
    }()
    select {case <-entered:case <-ctx.Done():t.Fatal("no callback")}
    stop()
    select {case err:=<-done:if !errors.Is(err,context.Canceled){t.Fatal(err)};case <-ctx.Done():t.Fatal("cancel did not join")}
}
func TestHiddenRefineryFailure(t *testing.T) {
    labkit.LeakCheck(t)
    failure:=errors.New("contaminated ore")
    ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second); defer cancel()
    got,err:=Refine(ctx,[]int{1,2,3,4},2,func(ctx context.Context,x int)(int,error){
        if x==1 {return 0,failure}; <-ctx.Done();return 0,ctx.Err()
    })
    if got!=nil || !errors.Is(err,failure) {t.Fatal("failure was lost",got,err)}
    canceled,stop:=context.WithCancel(ctx);stop()
    got,err=Refine(canceled,[]int{1},1,func(context.Context,int)(int,error){t.Error("canceled callback");return 0,nil})
    if got!=nil || !errors.Is(err,context.Canceled){t.Fatal(got,err)}
}
func TestHiddenBonusIndependentRefineries(t *testing.T) {
    labkit.LeakCheck(t)
    var wg sync.WaitGroup
    for i:=0;i<12;i++ { wg.Add(1);go func(i int){defer wg.Done()
        got,err:=Refine(context.Background(),[]int{i},2,func(_ context.Context,x int)(int,error){return x+1,nil})
        if err!=nil || !reflect.DeepEqual(got,[]int{i+1}) {t.Error("cross-call state",got,err)}
    }(i) }; wg.Wait()
}
`},
	})
}
