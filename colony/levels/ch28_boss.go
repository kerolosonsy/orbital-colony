package levels

func init() {
	add(&Level{
		ID: "28-boss-quarantine", Chapter: 28, Kind: "lab", Boss: true, LessonRevision: 1, Steps: ch28BossArabicSteps, Setup: seededLabRoom, Race: true,
		Title: "Boss: عزل الرصيف", Goal: "ابني API لعزل الكبسولات بعقد كامل وتصاريح وحدود",
		Requires: []string{"mutex", "sel:http.Server", "sel:http.MaxBytesReader"}, Bridge: []int{46, 47, 48, 49, 50, 60},
		Files: []LabFile{{Name: "main.go", Starter: `package main
import "net/http"
type QuarantinePod struct { ID string; Occupants int }
func NewQuarantine(token string) *http.Server {
    return &http.Server{Handler:http.NotFoundHandler()}
}
func main() {}
`}},
		Lesson: `# Boss: عزل الرصيف

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 26-2-rest-crud, 26-3-router, 26-4-envelope, 26-6-upload, 27-1-signed-pass, 28-4-graceful.

مثال أصغر ومختلف:
~~~go
w := httptest.NewRecorder()
h.ServeHTTP(w, httptest.NewRequest("GET", "/readyz", nil))
~~~

## أخطاء شائعة
~~~go
if token == "" { next.ServeHTTP(w, r) } // missing credentials grant access
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~NewQuarantine(token string) *http.Server~؛ token محقون سري مش فاضي. السيرفر مبيبدأش جوه الدالة: ~srv := NewQuarantine("crew-secret")~ و~srv.Handler~ للاختبار. ~Addr~ بالظبط ~127.0.0.1:0~، و~ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout~ كلهم موجبين وبحد أقصى دقيقة. المخزن جديد لكل سيرفر وآمن بالتزامن.

كل رد ~application/json~. افحص ~Authorization: Bearer <token>~ الأول؛ غيره 401 مع ~{"error":"unauthorized"}~. بعده مسار مش معروف 404 ~not_found~؛ الطرق غير GET أو POST على المجموعة 405 ~method~ مع ~Allow: GET, POST~ (على عنصر GET بس). ~/pods~: GET يرجع array غير nil، وPOST محتاج ~X-Role: medic~ وإلا 403 ~forbidden~. الدور هنا جايلك من حد موثوق سبق وتحقق منه؛ متعتمدش على هيدر دور من عميل في خدمة حقيقية.

POST يقبل JSON واحد ~QuarantinePod{ID: "med-4", Occupants: 0}~ بحجم <=1024 بايت؛ جسم أكبر 413 ~too_large~. JSON بايظ، ID فاضي أو فيه /، أو Occupants سالب: 400 ~invalid~. المكرر 409 ~conflict~، والجديد 201 مع نفس الكائن و~Location: /pods/<ID>~. GET ~/pods/<ID>~ يرجع 200 بالكائن أو 404 ~not_found~. الأخطاء كائن بحقل error بس؛ الفشل ميغيّرش المخزن. سيب ~main~ فاضية. النجمة التالتة: ترتيب قائمة ID تصاعدي، ونسخ ~X-Request-ID~ للرد حتى في الخطأ، وعزل مخازن السيرفرات.
`,
		Hints: ch28BossArabic.Hints,
		En: &LevelText{LessonRevision: 1, Steps: ch28BossEnglishSteps, Title: "Boss: Dock Quarantine", Goal: "Build a bounded quarantine API with a complete authorization contract", Files: []LabFile{{Name: "main.go", Starter: `package main
import "net/http"
type QuarantinePod struct { ID string; Occupants int }
func NewQuarantine(token string) *http.Server {
    return &http.Server{Handler:http.NotFoundHandler()}
}
func main() {}
`}},
			Lesson: `# Boss: Dock Quarantine

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 26-2-rest-crud, 26-3-router, 26-4-envelope, 26-6-upload, 27-1-signed-pass, 28-4-graceful.

A smaller, different example:
~~~go
w := httptest.NewRecorder()
h.ServeHTTP(w, httptest.NewRequest("GET", "/readyz", nil))
~~~

## Common mistakes
~~~go
if token == "" { next.ServeHTTP(w, r) } // missing credentials grant access
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~NewQuarantine(token string) *http.Server~; the injected secret token is nonempty. Do not start the server inside the function: ~srv := NewQuarantine("crew-secret")~ exposes ~srv.Handler~ for tests. Set exact ~Addr~ to ~127.0.0.1:0~; ~ReadHeaderTimeout, ReadTimeout, WriteTimeout, IdleTimeout~ must each be positive and at most one minute. Each server owns fresh concurrency-safe storage.

Every response is ~application/json~. First check ~Authorization: Bearer <token>~; otherwise return 401 ~{"error":"unauthorized"}~. Then unknown routes return 404 ~not_found~; unsupported methods return 405 ~method~ with ~Allow: GET, POST~ on the collection (GET only on an item). At ~/pods~, GET returns a nonnil array; POST requires ~X-Role: medic~, otherwise 403 ~forbidden~. The role is supplied by an already verified trusted boundary here; a real service must not trust a client role header.

POST accepts one JSON ~QuarantinePod{ID: "med-4", Occupants: 0}~, at most 1024 bytes; larger bodies return 413 ~too_large~. Malformed JSON, empty ID, ID containing /, or negative Occupants returns 400 ~invalid~. Duplicate ID returns 409 ~conflict~; creation returns 201 with the object and ~Location: /pods/<ID>~. GET ~/pods/<ID>~ returns 200 with the object or 404 ~not_found~. Errors contain only the error field; failed requests do not mutate storage. Leave ~main~ empty. Third star: sort lists by ID, echo ~X-Request-ID~ even on errors, and isolate server stores.
`,
			Hints: ch28BossEnglish.Hints},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: `package main
import (
    "encoding/json"
    "fmt"
    "math/rand"
    "net/http"
    "net/http/httptest"
    "strings"
    "sync"
    "testing"
    "time"
)
func hiddenPodRequest(h http.Handler,method,path,body,token,role string)*httptest.ResponseRecorder {
    req:=httptest.NewRequest(method,path,strings.NewReader(body))
    req.Header.Set("Authorization","Bearer "+token);req.Header.Set("X-Role",role);req.Header.Set("X-Request-ID","evidence-71")
    w:=httptest.NewRecorder();h.ServeHTTP(w,req);return w
}
func TestHiddenQuarantineContract(t *testing.T) {
    rng:=rand.New(rand.NewSource(28091));token:=fmt.Sprint(rng.Int63())
    srv:=NewQuarantine(token)
    if srv==nil || srv.Handler==nil{t.Fatal("server and handler required")}
    for _,d:=range []time.Duration{srv.ReadHeaderTimeout,srv.ReadTimeout,srv.WriteTimeout,srv.IdleTimeout}{if d<=0 || d>time.Minute{t.Fatal("unbounded timeout",d)}}
    if srv.Addr!="127.0.0.1:0"{t.Fatal("must bind random loopback port")}
    for _,tc:=range []struct{method,path,body,token,role string;status int;code string}{
        {"GET","/pods","","wrong","medic",401,"unauthorized"},
        {"POST","/pods",` + "`" + `{"ID":"a","Occupants":1}` + "`" + `,token,"observer",403,"forbidden"},
        {"GET","/unknown","",token,"medic",404,"not_found"},
        {"GET","/pods/absent","",token,"medic",404,"not_found"},
        {"DELETE","/pods","",token,"medic",405,"method"},
        {"POST","/pods","broken",token,"medic",400,"invalid"},
        {"POST","/pods",` + "`" + `{"ID":"","Occupants":1}` + "`" + `,token,"medic",400,"invalid"},
        {"POST","/pods",` + "`" + `{"ID":"a/b","Occupants":1}` + "`" + `,token,"medic",400,"invalid"},
        {"POST","/pods",` + "`" + `{"ID":"a","Occupants":-1}` + "`" + `,token,"medic",400,"invalid"},
        {"POST","/pods",strings.Repeat("x",1025),token,"medic",413,"too_large"},
    }{
        w:=hiddenPodRequest(srv.Handler,tc.method,tc.path,tc.body,tc.token,tc.role)
        var got map[string]string
        if err:=json.Unmarshal(w.Body.Bytes(),&got);err!=nil{t.Fatal(err,w.Body.String())}
        if w.Code!=tc.status || got["error"]!=tc.code || len(got)!=1 || !strings.HasPrefix(w.Header().Get("Content-Type"),"application/json"){t.Fatal(tc,w.Code,w.Body.String())}
        if tc.status==405 && w.Header().Get("Allow")!="GET, POST"{t.Fatal("missing method contract")}
    }
    for i:=0;i<12;i++{
        id:=fmt.Sprintf("pod-%02d",i);body:=fmt.Sprintf(` + "`" + `{"ID":%q,"Occupants":%d}` + "`" + `,id,i)
        w:=hiddenPodRequest(srv.Handler,"POST","/pods",body,token,"medic")
        if w.Code!=201 || w.Header().Get("Location")!="/pods/"+id{t.Fatal(w.Code,w.Body.String())}
        w=hiddenPodRequest(srv.Handler,"POST","/pods",body,token,"medic");if w.Code!=409{t.Fatal("duplicate",w.Code)}
        w=hiddenPodRequest(srv.Handler,"GET","/pods/"+id,"",token,"observer")
        var pod QuarantinePod;if err:=json.Unmarshal(w.Body.Bytes(),&pod);err!=nil || w.Code!=200 || pod.ID!=id || pod.Occupants!=i{t.Fatal(pod,w.Code,err)}
    }
    var list []QuarantinePod
    w:=hiddenPodRequest(srv.Handler,"GET","/pods","",token,"observer")
    if err:=json.Unmarshal(w.Body.Bytes(),&list);err!=nil || w.Code!=200 || len(list)!=12{t.Fatal(list,w.Code,err)}
}
func TestHiddenQuarantineConcurrentConflict(t *testing.T){
    h:=NewQuarantine("crew").Handler
    var wg sync.WaitGroup
    statuses:=make(chan int,24)
    for i:=0;i<24;i++{wg.Add(1);go func(){defer wg.Done();statuses<-hiddenPodRequest(h,"POST","/pods",` + "`" + `{"ID":"one","Occupants":0}` + "`" + `,"crew","medic").Code}()}
    wg.Wait();close(statuses)
    created:=0
    for code:=range statuses{if code==201{created++}else if code!=409{t.Fatal(code)}}
    if created!=1{t.Fatal("non-atomic duplicate check",created)}
}
func TestHiddenBonusQuarantineCorrelationAndOrder(t *testing.T){
    h:=NewQuarantine("crew").Handler
    for _,id:=range []string{"z","b","a"}{w:=hiddenPodRequest(h,"POST","/pods",fmt.Sprintf(` + "`" + `{"ID":%q,"Occupants":0}` + "`" + `,id),"crew","medic");if w.Code!=201{t.Fatal(w.Code)}}
    w:=hiddenPodRequest(h,"GET","/pods","","crew","observer")
    var list []QuarantinePod
    if err:=json.Unmarshal(w.Body.Bytes(),&list);err!=nil || len(list)!=3{t.Fatal(list,err)}
    if list[0].ID!="a" || list[1].ID!="b" || list[2].ID!="z"{t.Fatal("unstable order",list)}
    for _,token:=range []string{"crew","wrong"}{w:=hiddenPodRequest(h,"GET","/pods","",token,"observer");if w.Header().Get("X-Request-ID")!="evidence-71"{t.Fatal("correlation lost on success/failure")}}
    w=hiddenPodRequest(NewQuarantine("crew").Handler,"GET","/pods","","crew","observer")
    if strings.TrimSpace(w.Body.String())!="[]"{t.Fatal("server instances share state",w.Body.String())}
}
`},
	})
}
