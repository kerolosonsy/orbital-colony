package levels

func init() {
	add(&Level{LessonRevision: 1, Steps: ch26Lesson261.arabic, ID: "26-1-first-server", Chapter: 26, Kind: "lab", Title: "أول سيرفر", Goal: "ابني صفحة القاعدة ونبض الجاهزية بحدود سيرفر واضحة",
		Requires: []string{"sel:http.Server"}, Bridge: []int{46}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch261Starter}},
		Lesson: `# أول بوابة HTTP للقاعدة

## قبل ما تكتب
- إيه اللي لازم تحدده قبل أول بايت في الرد؟
- إزاي تختبر الصفحة من غير ما تفتح بورت؟

## الطلب والرد
الـ ~http.Handler~ واجهة ليها ~ServeHTTP~. الـ ~http.HandlerFunc~ بتحوّل دالة لنفس الواجهة. ~http.ResponseWriter~ بيكتب الرد، و ~*http.Request~ فيه ~Method~ و ~URL.Path~؛ الـ query مش جزء من المسار.
~~~go
h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 if r.URL.Path != "/dock" { http.NotFound(w,r); return }
 w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 w.WriteHeader(http.StatusOK)
 _, _ = w.Write([]byte("open\n"))
})
h.ServeHTTP(w,r)
fallback := http.NotFoundHandler() // handler بيرجع 404
~~~
~Header().Set~ لازم قبل ~WriteHeader~. أول ~Write~ من غير حالة بيختار 200. ~http.StatusOK~ معناها 200 نجاح، و ~http.NotFound~ تكتب 404. الحالة 405 معناها المسار معروف بس الطريقة ممنوعة؛ حط ~Allow~ بالطرق المسموحة. الحالة 500 فشل داخلي. ~http.MethodGet~ و ~http.MethodHead~ أسماء للطرق؛ ~HEAD~ هنا مرفوض، بس رده من غير جسم.
~~~go
w.Header().Set("Allow", "GET")
w.WriteHeader(http.StatusMethodNotAllowed) // 405
if r.Method != http.MethodHead { _, _ = w.Write([]byte("wrong method\n")) }
~~~

## قالب HTML موثوق وبيانات غير موثوقة
استورد ~html/template~. ~template.New~ بتسمّي القالب، و ~Parse~ تراجعه مرة قبل بناء الـ handler. ~Execute~ بتدخل البيانات مع escaping تلقائي؛ ~{{.Label}}~ تقرأ حقل متصدّر. متحوّلش بيانات العميل لـ ~template.HTML~ لأنها بتلغي الحماية.
~~~go
page, err := template.New("dock").Parse("<p>{{.Label}}</p>")
if err != nil { return err }
var out bytes.Buffer
err = page.Execute(&out, struct{ Label string }{"<alarm>"})
~~~
~bytes.Buffer~ من الدروس السابقة ينفع كمساحة مؤقتة. افحص خطأ التنفيذ قبل نسخ ~out.Bytes()~ للرد؛ بعد الالتزام بالرد مينفعش تغيّر الحالة أو تضيف مستند خطأ تاني. أخطاء كتابة الرد بعد الالتزام تنهي الكتابة، مش محاولة رد جديد.

## إعداد السيرفر مش تشغيله
~http.Server~ بيربط ~Addr~ و ~Handler~ بحدود الاتصال. ~ReadHeaderTimeout~ لقراءة الهيدر، و ~ReadTimeout~ لقراءة الطلب، و ~WriteTimeout~ لفترة الكتابة، و ~IdleTimeout~ لانتظار اتصال خامل. ~MaxHeaderBytes~ حد الهيدر مش الجسم. الحدود دي مش بتوقف حساب handler بيتجاهل الإلغاء.
~~~go
s := &http.Server{Addr: "127.0.0.1:0", Handler: h,
 ReadHeaderTimeout: time.Second, ReadTimeout: 3*time.Second,
 WriteTimeout: 4*time.Second, IdleTimeout: 20*time.Second, MaxHeaderBytes: 8192}
~~~
في المهمة رجّع الإعداد بس. مفيش نداء تشغيل شبكة أو وقت حقيقي.

## اختبار في الذاكرة
~httptest.NewRequest~ تبني طلب داخل، و ~httptest.NewRecorder~ تسجّل الحالة والهيدر والجسم. ~Code~ و ~Body.String()~ و ~Header().Get~ تخليك تفحص العقد بدقة.
~~~go
w := httptest.NewRecorder()
h.ServeHTTP(w, httptest.NewRequest("GET", "/dock", nil))
if w.Code != 200 { t.Fatal(w.Code) }
~~~
استورد ~net/http/httptest~ في التست. مفيش اتصال شبكة ولا بورت ثابت؛ إعدادات الوقت تتفحص كقيم من غير نوم.

## أخطاء شائعة
~~~go
w.Write([]byte(base))
w.Header().Set("Content-Type", "text/html") // متأخر، وكمان البيانات مش escaped
w.WriteHeader(500) // الحالة اتثبتت خلاص
~~~

## المهمة
نفّذ ~NewGateway(base, markup string) (http.Handler,error)~. القالب موثوق ومحقون، وبيتنفذ ببيانات ~struct{ Base string }~. خطأ ~Parse~ يرجع ~nil,error~. ~GET /base~ يرجع 200 و ~text/html; charset=utf-8~ وجسم القالب ببيانات القاعدة escaped. ~GET /pulse~ يرجع 200 و ~text/plain; charset=utf-8~ و ~ready\n~ بالظبط. المسارات حرفية وحساسة للحالة والـ query متأثرش عليها.

أي مسار تاني يرجع 404 و ~unknown gate\n~؛ افحص المسار قبل الطريقة. أي طريقة غير ~GET~ للمسارين المعروفين ترجع 405 و ~Allow: GET~ و ~use GET\n~. كل الأخطاء نص عادي بنفس content type، وطلبات ~HEAD~ جسمها فاضي حتى مع 404.

نفّذ ~GatewayServer(addr string,h http.Handler) *http.Server~ يحفظ العنوان والـ handler المحقونين؛ حدود الثواني بالترتيب: ~ReadHeaderTimeout=2~ و ~ReadTimeout=5~ و ~WriteTimeout=8~ و ~IdleTimeout=30~، و ~MaxHeaderBytes=16384~. متشغّلوش، وسيب ~main~ فاضية؛ مفيش ~Report~.

النجمة التالتة: خطأ ~Execute~ يرجع 500 و ~page unavailable\n~ بنص عادي من غير أي جزء من الصفحة.`,
		Hints: []string{"افصل اختيار المسار عن اختيار طريقة الطلب.", "جهّز handler صغيرًا واختبره بمسجل داخل الذاكرة.", "راجع ~ResponseRecorder~ قبل إضافة إعدادات السيرفر."},
		En: &LevelText{LessonRevision: 1, Steps: ch26Lesson261.english, Title: "First Server", Goal: "Build the base page and readiness pulse with explicit server limits", Files: []LabFile{{Name: "main.go", Starter: ch261Starter}},
			Lesson: `# The Base's First HTTP Gateway

## Before you write
- What must be decided before the first response byte?
- How can you test a page without opening a port?

## Request and response
~http.Handler~ is an interface with ~ServeHTTP~. ~http.HandlerFunc~ adapts a function to it. ~http.ResponseWriter~ writes the response; ~*http.Request~ supplies ~Method~ and ~URL.Path~. Query parameters are separate from the path.
~~~go
h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 if r.URL.Path != "/dock" { http.NotFound(w,r); return }
 w.Header().Set("Content-Type", "text/plain; charset=utf-8")
 w.WriteHeader(http.StatusOK)
 _, _ = w.Write([]byte("open\n"))
})
h.ServeHTTP(w,r)
fallback := http.NotFoundHandler() // A handler returning 404
~~~
Call ~Header().Set~ before ~WriteHeader~. A first ~Write~ implicitly selects 200. ~http.StatusOK~ means successful 200; ~http.NotFound~ writes 404. Status 405 means a known path rejects the method; provide ~Allow~. Status 500 means internal failure. ~http.MethodGet~ and ~http.MethodHead~ name methods. This mission rejects ~HEAD~ but keeps its response body empty.
~~~go
w.Header().Set("Allow", "GET")
w.WriteHeader(http.StatusMethodNotAllowed) // 405
if r.Method != http.MethodHead { _, _ = w.Write([]byte("wrong method\n")) }
~~~

## Trusted HTML, untrusted data
Import ~html/template~. ~template.New~ names a template; ~Parse~ validates it once before handler construction. ~Execute~ inserts escaped data; ~{{.Label}}~ reads an exported field. Never bypass escaping by declaring untrusted text to be trusted HTML.
~~~go
page, err := template.New("dock").Parse("<p>{{.Label}}</p>")
if err != nil { return err }
var out bytes.Buffer
err = page.Execute(&out, struct{ Label string }{"<alarm>"})
~~~
The previously learned ~bytes.Buffer~ stages output. Check execution errors before copying ~out.Bytes()~ into the response. After commitment, changing status or appending another error document cannot repair it. Stop after a response write failure rather than attempting another response.

## Configure without starting
~http.Server~ connects ~Addr~ and ~Handler~ with connection limits. ~ReadHeaderTimeout~ covers headers, ~ReadTimeout~ request reading, ~WriteTimeout~ the writing period, and ~IdleTimeout~ idle keep-alive waiting. ~MaxHeaderBytes~ caps headers rather than the body. These settings cannot stop arbitrary computation ignoring cancellation.
~~~go
s := &http.Server{Addr: "127.0.0.1:0", Handler: h,
 ReadHeaderTimeout: time.Second, ReadTimeout: 3*time.Second,
 WriteTimeout: 4*time.Second, IdleTimeout: 20*time.Second, MaxHeaderBytes: 8192}
~~~
Return configuration only. Do not start networking or depend on real time.

## Test in memory
~httptest.NewRequest~ creates an incoming request. ~httptest.NewRecorder~ records status, headers, and body. Inspect ~Code~, ~Body.String()~, and ~Header().Get~ for exact assertions.
~~~go
w := httptest.NewRecorder()
h.ServeHTTP(w, httptest.NewRequest("GET", "/dock", nil))
if w.Code != 200 { t.Fatal(w.Code) }
~~~
Import ~net/http/httptest~ in tests. No socket or fixed port is needed; inspect timeout values without sleeping.

## Common mistakes
~~~go
w.Write([]byte(base))
w.Header().Set("Content-Type", "text/html") // Too late; data is also unescaped
w.WriteHeader(500) // Status already committed
~~~

## Your mission
Implement ~NewGateway(base, markup string) (http.Handler,error)~. The trusted template is injected and executed with ~struct{ Base string }~ data. A ~Parse~ failure returns ~nil,error~. ~GET /base~ returns 200, ~text/html; charset=utf-8~, and the template body with escaped base data. ~GET /pulse~ returns 200, ~text/plain; charset=utf-8~, and exactly ~ready\n~. Paths are exact and case-sensitive; ignore query parameters.

Every other path returns 404 and ~unknown gate\n~; select the path before checking methods. Non-~GET~ methods on known paths return 405, ~Allow: GET~, and ~use GET\n~. Errors use the same plain-text content type. ~HEAD~ always has an empty body, including 404.

Implement ~GatewayServer(addr string,h http.Handler) *http.Server~ retaining the injected address and handler. Timeout seconds: ~ReadHeaderTimeout=2~, ~ReadTimeout=5~, ~WriteTimeout=8~, ~IdleTimeout=30~; ~MaxHeaderBytes=16384~. Do not start it. Leave ~main~ empty; no ~Report~.

Third star: an ~Execute~ failure returns 500 and plain-text ~page unavailable\n~ without any partial page.`,
			Hints: []string{"Separate path selection from request-method selection.", "Build a small handler and check it with an in-memory recorder.", "Inspect ~ResponseRecorder~ before adding server settings."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch261Hidden},
	})
}

const ch261Starter = `package main

import (
	"errors"
	"net/http"
)

func NewGateway(base, markup string) (http.Handler, error)   { return nil, errors.New("TODO") }
func GatewayServer(addr string, h http.Handler) *http.Server { return &http.Server{} }
func main()                                                  {}
`

const ch261Hidden = `package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Boundaries: empty/Unicode/markup data, exact paths, methods, parse and execution failures.
func TestHiddenGatewayRoutesAndEscaping(t *testing.T) {
	rng := rand.New(rand.NewSource(261))
	for _, base := range []string{"", "orbital", "محطة 🚀", fmt.Sprintf("<script>%d</script>&\"", rng.Intn(999999))} {
		h, e := NewGateway(base, "<h1>{{.Base}}</h1>")
		if e != nil || h == nil {
			t.Fatal("construct", e)
		}
		for _, tc := range []struct {
			method, path    string
			status          int
			body, ct, allow string
		}{
			{"GET", "/base?crew=1", 200, "<h1>" + template.HTMLEscapeString(base) + "</h1>", "text/html; charset=utf-8", ""},
			{"GET", "/pulse", 200, "ready\n", "text/plain; charset=utf-8", ""},
			{"POST", "/pulse", 405, "use GET\n", "text/plain; charset=utf-8", "GET"},
			{"PUT", "/base", 405, "use GET\n", "text/plain; charset=utf-8", "GET"},
			{"HEAD", "/base", 405, "", "text/plain; charset=utf-8", "GET"},
			{"GET", "/", 404, "unknown gate\n", "text/plain; charset=utf-8", ""},
			{"GET", "/pulse/", 404, "unknown gate\n", "text/plain; charset=utf-8", ""},
			{"POST", "/missing", 404, "unknown gate\n", "text/plain; charset=utf-8", ""},
			{"HEAD", "/missing", 404, "", "text/plain; charset=utf-8", ""},
		} {
			w := httptest.NewRecorder()
			h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status || w.Body.String() != tc.body || w.Header().Get("Content-Type") != tc.ct || w.Header().Get("Allow") != tc.allow {
				t.Fatalf("%s %s: %d %q %v", tc.method, tc.path, w.Code, w.Body.String(), w.Header())
			}
		}
	}
	if h, e := NewGateway("x", "{{"); h != nil || e == nil {
		t.Fatal("parse failure")
	}
}
func TestHiddenGatewayServerLimits(t *testing.T) {
	h := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	s := GatewayServer("127.0.0.1:0", h)
	if s == nil {
		t.Fatal("nil server")
	}
	if s.Addr != "127.0.0.1:0" || s.Handler == nil || s.ReadHeaderTimeout != 2*time.Second || s.ReadTimeout != 5*time.Second || s.WriteTimeout != 8*time.Second || s.IdleTimeout != 30*time.Second || s.MaxHeaderBytes != 16<<10 {
		t.Fatalf("server limits: %+v", s)
	}
}
func TestHiddenBonusGatewayRenderBeforeCommit(t *testing.T) {
	h, e := NewGateway("station", "prefix{{.Missing}}")
	if e != nil {
		t.Fatal(e)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/base", nil))
	if w.Code != 500 || w.Body.String() != "page unavailable\n" || strings.Contains(w.Body.String(), "prefix") {
		t.Fatal("partial page", w.Code, w.Body.String())
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch26Lesson262.arabic, ID: "26-2-rest-crud", Chapter: 26, Kind: "lab", Title: "REST", Goal: "دير أرصفة المستعمرة بـ CRUD متزامن وIDs ثابتة",
		Requires: []string{"mutex"}, Bridge: []int{47}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch262Starter}},
		Lesson: `# سجل أرصفة من خلال HTTP

## قبل ما تكتب
- إيه اللي لازم يحصل تحت نفس القفل عشان رقم الرصيف ميتكررش؟
- هترفض الطلب البايظ قبل ما تغيّر المخزن إزاي؟

## JSON على السلك
CRUD هنا بيتحوّل لطرق: ~POST~ ينشئ، و ~GET~ يقرأ، و ~PUT~ يستبدل بالكامل، و ~DELETE~ يحذف. 201 معناها إنشاء، و ~Location~ مسار العنصر الجديد. 204 نجاح من غير جسم، و 400 مدخل بايظ، و 413 جسم أكبر من الحد.
~~~go
raw, err := io.ReadAll(io.LimitReader(r.Body, 101))
if err != nil { return }
if len(raw) > 100 { w.WriteHeader(413); return }
var incoming *Berth
err = json.Unmarshal(raw, &incoming)
~~~
~r.Body~ قارئ. ~io.LimitReader~ توقف القراءة عند العدد المحدد؛ اقرأ الحد زائد 1 عشان تكتشف الزيادة. ~json.Unmarshal~ من فصل الأرشيف بترفض مستندين متتاليين. الـ pointer يخلي ~null~ مميّزة بـ ~nil~؛ افحصه قبل الحقول. الحقول المجهولة مسموحة هنا.
~~~go
raw, err := json.Marshal(Berth{ID: 7, Label: "west", Slots: 2})
if err != nil { return }
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Location", "/berths/7")
w.WriteHeader(201)
_, _ = w.Write(raw)
~~~
~Berth~ فيها ~ID~ و ~Label~ و ~Slots~، بأسماء JSON ~id~ و ~label~ و ~slots~ باستخدام tags اللي اتشرحت قبل كده. مثال النداء بعد التنفيذ:
~~~go
h := NewBerthAPI()
h.ServeHTTP(w, httptest.NewRequest("GET", "/berths", nil))
~~~

## هوية ومسار وحالة مستقرة
~strconv.Atoi~ تقرأ رقم، و ~strconv.Itoa~ ترجع صورته القياسية؛ مقارنة الاتنين ترفض ~01~ و ~+1~. افصل شكل المسار عن صلاحية الرقم. خُد نسخة من القائمة تحت ~sync.Mutex~، وافتح القفل قبل ترميز الرد. رتّب النسخة بـ ~slices.SortFunc~ حسب ~ID~. العداد منفصل عن طول الخريطة عشان الحذف ميعيدش الأرقام.
~~~go
segment := strings.TrimPrefix("/berths/7", "/berths/")
id, err := strconv.Atoi(segment)
if err != nil || id <= 0 || strconv.Itoa(id) != segment { w.WriteHeader(400); return }
~~~

## أخطاء شائعة
~~~go
next = len(berths)+1 // يعيد رقم اتحذف
mu.Unlock()
for _, b := range berths { list = append(list,b) } // قراية خريطة من غير القفل
~~~

## المهمة
نفّذ ~NewBerthAPI() http.Handler~ بمخزن مستقل وآمن بالتزامن. ~Berth~ زي الـ starter. ~GET /berths~ يرجع 200 ومصفوفة JSON مرتبة تصاعدياً، والفاضي ~[]~. ~POST /berths~ يقبل كائن JSON واحد بـ ~label~ بعد trim غير فاضي و ~slots~ بين 1 و 64 شاملة؛ يتجاهل ~id~ والحقول المجهولة. أول ID هو 1 وبعده أرقام متزايدة لا يعاد استخدامها؛ الطلب المرفوض ميستهلكش رقم. يرجع 201 والعنصر و ~Location: /berths/ID~.

~GET /berths/ID~ يرجع 200 والعنصر. ~PUT~ لنفس المسار يستبدل الحقلين بنفس التحقق ويحفظ ID ويرجع 200 والعنصر؛ ~DELETE~ يحذف ويرجع 204 بجسم فاضي. العنصر الغائب يرجع 404 و ~unknown berth\n~. خطأ JSON أو نوع أو تحقق أو رقم غير موجب أو غير قياسي يرجع 400 و ~invalid berth\n~. حدد الشكل والطريقة ثم تحقق الجسم قبل البحث عن عنصر ~PUT~. حد الجسم في الإنشاء والتحديث 4096 بايت شاملة المسافات؛ الأكبر 413 و ~too large\n~. الفاضي و ~null~ والمصفوفة مرفوضين؛ مفيش شرط على content type الداخل.

المسارات الحرفية بس؛ رقم غير صالح في جزء واحد بعد ~/berths/~ هو 400، والجزء الفاضي أو أجزاء زيادة أو مسار غريب 404. الطرق غير المدعومة 405 و ~use route methods\n~ و ~Allow: GET, POST~ للمجموعة أو ~Allow: DELETE, GET, PUT~ للعنصر. أي ~HEAD~ جسمه فاضي. الردود الناجحة ذات الجسم ~application/json~، والأخطاء ~text/plain; charset=utf-8~. الفشل ميغيرش الحالة. الـ query تتجاهلها. مفيش شبكة أو وقت أو ~Report~، و ~main~ فاضية.`,
		Hints: []string{"عامل الهوية كحالة مستقلة عن عدد السجلات الحالية.", "خلي القراءة والتعديل ضمن حدود القفل نفسه.", "راجع نسخة القائمة قبل فرزها أو ترميزها."},
		En: &LevelText{LessonRevision: 1, Steps: ch26Lesson262.english, Title: "REST", Goal: "Manage colony berths with concurrent CRUD and stable IDs", Files: []LabFile{{Name: "main.go", Starter: ch262Starter}},
			Lesson: `# A Berth Registry over HTTP

## Before you write
- What must share one lock to prevent duplicate berth IDs?
- How will malformed requests fail before changing state?

## JSON on the wire
CRUD maps to methods: ~POST~ creates, ~GET~ reads, ~PUT~ fully replaces, and ~DELETE~ removes. Status 201 means creation; ~Location~ names the new resource. Status 204 has no body, 400 means invalid input, and 413 means oversized input.
~~~go
raw, err := io.ReadAll(io.LimitReader(r.Body, 101))
if err != nil { return }
if len(raw) > 100 { w.WriteHeader(413); return }
var incoming *Berth
err = json.Unmarshal(raw, &incoming)
~~~
~r.Body~ is a reader. ~io.LimitReader~ stops at a byte count; read the cap plus one to detect overflow. Previously learned ~json.Unmarshal~ rejects concatenated documents. A pointer distinguishes ~null~ as ~nil~; check before accessing fields. Unknown fields are allowed here.
~~~go
raw, err := json.Marshal(Berth{ID: 7, Label: "west", Slots: 2})
if err != nil { return }
w.Header().Set("Content-Type", "application/json")
w.Header().Set("Location", "/berths/7")
w.WriteHeader(201)
_, _ = w.Write(raw)
~~~
~Berth~ has ~ID~, ~Label~, and ~Slots~, tagged ~id~, ~label~, and ~slots~ using earlier JSON lessons. After implementation:
~~~go
h := NewBerthAPI()
h.ServeHTTP(w, httptest.NewRequest("GET", "/berths", nil))
~~~

## Stable identity, paths, and snapshots
~strconv.Atoi~ parses numbers; ~strconv.Itoa~ produces canonical spelling. Comparing rejects ~01~ and ~+1~. Distinguish path shape from numeric validity. Copy lists under ~sync.Mutex~, release it before encoding, and sort the copy by ~ID~ using ~slices.SortFunc~. Keep the ID counter separate from map length so deletion cannot recycle IDs.
~~~go
segment := strings.TrimPrefix("/berths/7", "/berths/")
id, err := strconv.Atoi(segment)
if err != nil || id <= 0 || strconv.Itoa(id) != segment { w.WriteHeader(400); return }
~~~

## Common mistakes
~~~go
next = len(berths)+1 // Recycles deleted IDs
mu.Unlock()
for _, b := range berths { list = append(list,b) } // Unlocked map read
~~~

## Your mission
Implement ~NewBerthAPI() http.Handler~ with an independent concurrent-safe store. Keep the starter's ~Berth~ fields. ~GET /berths~ returns 200 and an ascending-ID JSON array, ~[]~ when empty. ~POST /berths~ accepts one JSON object: trim a nonblank ~label~, require ~slots~ in inclusive 1..64, and ignore input ~id~ and unknown fields. IDs start at 1, strictly increase, and are never reused; rejected requests consume none. Return 201, the stored object, and ~Location: /berths/ID~.

~GET /berths/ID~ returns 200 and the object. ~PUT~ fully replaces both fields with the same validation, retains ID, and returns 200 and the object. ~DELETE~ removes it and returns bodyless 204. Missing items return 404 and ~unknown berth\n~. Invalid JSON, types, values, or nonpositive/noncanonical IDs return 400 and ~invalid berth\n~. Select shape and method, then validate a ~PUT~ body before item lookup. Creation/replacement bodies are capped at 4096 bytes including whitespace; larger bodies return 413 and ~too large\n~. Empty, ~null~, and array inputs fail; request content type is unrestricted.

Paths are exact. An invalid single ID segment after ~/berths/~ is 400; empty/extra segments and unknown paths are 404. Unsupported methods return 405, ~use route methods\n~, and ~Allow: GET, POST~ for collection or ~Allow: DELETE, GET, PUT~ for items. ~HEAD~ is bodyless. Successful bodies use ~application/json~; errors use ~text/plain; charset=utf-8~. Failures never mutate state. Ignore query parameters. No network, time, or ~Report~; leave ~main~ empty.`,
			Hints: []string{"Treat identity as separate from the number of current records.", "Keep reads and mutations inside the same lock boundary.", "Inspect a list copy before sorting or encoding it."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch262Hidden},
	})
}

const ch262Starter = `package main

import "net/http"

type Berth struct {
	ID    int    ` + "`" + `json:"id"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
	Slots int    ` + "`" + `json:"slots"` + "`" + `
}

func NewBerthAPI() http.Handler { return http.NotFoundHandler() }
func main()                     {}
`

const ch262Hidden = `package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func berthCall(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
	return w
}
func berthDecode(t *testing.T, w *httptest.ResponseRecorder, status int) Berth {
	t.Helper()
	var b Berth
	if w.Code != status || w.Header().Get("Content-Type") != "application/json" || json.Unmarshal(w.Body.Bytes(), &b) != nil {
		t.Errorf("response %d %s", w.Code, w.Body.String())
	}
	return b
}

// Boundaries: empty lists, invalid JSON and IDs, slot limits, replacements, deleted IDs, isolation.
func TestHiddenBerthsCRUDAndValidation(t *testing.T) {
	h := NewBerthAPI()
	rng := rand.New(rand.NewSource(262))
	want := []Berth{}
	w := berthCall(h, "GET", "/berths", "")
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatal("empty list", w.Body.String())
	}
	for _, body := range []string{"", "null", "[]", "{}", ` + "`" + `{"label":7,"slots":2}` + "`" + `, ` + "`" + `{"label":"x","slots":0}` + "`" + `, ` + "`" + `{"label":"x","slots":65}` + "`" + `, ` + "`" + `{"label":" ","slots":1}` + "`" + `, ` + "`" + `{"label":"x","slots":1} {}` + "`" + `} {
		if w := berthCall(h, "POST", "/berths", body); w.Code != 400 {
			t.Fatal("bad body", body, w.Code)
		}
	}
	if w := berthCall(h, "POST", "/berths", strings.Repeat(" ", 4097)); w.Code != 413 {
		t.Fatal("size cap")
	}
	body := ` + "`" + `{"label":"edge","slots":1}` + "`" + `
	body += strings.Repeat(" ", 4096-len(body))
	first := berthDecode(t, berthCall(h, "POST", "/berths", body), 201)
	if first.ID != 1 {
		t.Fatal("invalid input consumed ID")
	}
	if w := berthCall(h, "DELETE", "/berths/1", ""); w.Code != 204 || w.Body.Len() != 0 {
		t.Fatal("delete")
	}
	for i := 2; i <= 25; i++ {
		label := fmt.Sprintf("bay-%d-🚀", rng.Intn(100000))
		slots := 1 + rng.Intn(64)
		w := berthCall(h, "POST", "/berths", fmt.Sprintf(` + "`" + `{"label":" %s ","slots":%d}` + "`" + `, label, slots))
		b := berthDecode(t, w, 201)
		if b != (Berth{i, label, slots}) || w.Header().Get("Location") != fmt.Sprintf("/berths/%d", i) {
			t.Fatal("identity", b, w.Header())
		}
		want = append(want, b)
	}
	for i, b := range want {
		path := fmt.Sprintf("/berths/%d", b.ID)
		if got := berthDecode(t, berthCall(h, "GET", path+"?crew=1", ""), 200); got != b {
			t.Fatal("read", got, b)
		}
		b.Label += "-updated"
		b.Slots = 64
		raw, _ := json.Marshal(b)
		got := berthDecode(t, berthCall(h, "PUT", path, string(raw)), 200)
		if got != b {
			t.Fatal("replacement")
		}
		want[i] = b
		if berthCall(h, "PUT", path, ` + "`" + `{"label":"bad","slots":-1}` + "`" + `).Code != 400 {
			t.Fatal("invalid update")
		}
	}
	var got []Berth
	w = berthCall(h, "GET", "/berths", "")
	if json.Unmarshal(w.Body.Bytes(), &got) != nil || !reflect.DeepEqual(got, want) {
		t.Fatal("list ordering/state", got)
	}
	for _, tc := range []struct {
		method, path, allow string
		status              int
	}{
		{"GET", "/berths/0", "", 400}, {"GET", "/berths/01", "", 400}, {"GET", "/berths/-1", "", 400}, {"GET", "/berths/+2", "", 400}, {"GET", "/berths/9999999999999999999999", "", 400}, {"GET", "/berths/no", "", 400},
		{"GET", "/berths/", "", 404}, {"GET", "/berths/2/x", "", 404}, {"POST", "/missing", "", 404},
		{"GET", "/berths/1", "", 404}, {"DELETE", "/berths/1", "", 404}, {"PUT", "/berths/1", "", 404},
		{"DELETE", "/berths", "GET, POST", 405}, {"POST", "/berths/2", "DELETE, GET, PUT", 405}, {"HEAD", "/berths/2", "DELETE, GET, PUT", 405},
	} {
		w := berthCall(h, tc.method, tc.path, ` + "`" + `{"label":"new","slots":1}` + "`" + `)
		if w.Code != tc.status || w.Header().Get("Allow") != tc.allow || (tc.method == "HEAD" && w.Body.Len() != 0) {
			t.Fatal(tc, w.Code, w.Header())
		}
	}
	other := NewBerthAPI()
	if berthDecode(t, berthCall(other, "POST", "/berths", ` + "`" + `{"label":"isolated","slots":1}` + "`" + `), 201).ID != 1 {
		t.Fatal("shared stores")
	}
}
func TestHiddenBerthsParallelRequests(t *testing.T) {
	h := NewBerthAPI()
	var wg sync.WaitGroup
	start := make(chan struct{})
	ids := make(chan int, 64)
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			b := berthDecode(t, berthCall(h, "POST", "/berths", fmt.Sprintf(` + "`" + `{"label":"dock-%d","slots":1}` + "`" + `, i)), 201)
			ids <- b.ID
		}(i)
	}
	close(start)
	wg.Wait()
	close(ids)
	seen := map[int]bool{}
	for id := range ids {
		if seen[id] || id < 1 || id > 64 {
			t.Fatal("duplicate/invalid ID", id)
		}
		seen[id] = true
	}
	if len(seen) != 64 {
		t.Fatal("lost creates")
	}
	for id := 1; id <= 64; id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			path := fmt.Sprintf("/berths/%d", id)
			berthDecode(t, berthCall(h, "GET", path, ""), 200)
			berthDecode(t, berthCall(h, "PUT", path, ` + "`" + `{"label":"replaced","slots":2}` + "`" + `), 200)
			if berthCall(h, "DELETE", path, "").Code != 204 {
				t.Error("delete")
			}
			_ = berthCall(h, "GET", "/berths", "")
		}(id)
	}
	wg.Wait()
	if w := berthCall(h, "GET", "/berths", ""); strings.TrimSpace(w.Body.String()) != "[]" {
		t.Fatal("surviving entries")
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch26Lesson263.arabic, ID: "26-3-router", Chapter: 26, Kind: "lab", Title: "الراوتر", Goal: "اربط أرصفة بباراميترات وسلسلة middleware وهوية لكل طلب",
		Requires: []string{"funcvalue"}, Bridge: []int{48}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch263Starter}},
		Lesson: `# راوتر صغير وسلسلة بوابات

## قبل ما تكتب
- إزاي تفرّق مسار مش موجود عن طريقة مش مسموحة؟
- مين لازم يشتغل الأول، أول middleware ولا آخر واحدة؟

## المطابقة مش تنظيف المسار
~Route~ بتجمع ~Method~ و ~Pattern~ و ~Handler~. ~slices.Clone~ تعمل نسخة سطحية من slice التسجيلات عشان تعديل الأصل ميمسش الراوتر. الجزء اللي أوله ~:~ اسم parameter؛ قسم النمط والمسار بنفس ~strings.Split~ وقارن الأجزاء. متستخدمش تنظيف يحوّل مسار غلط لصحيح.
~~~go
route := Route{Method: "GET", Pattern: "/bays/:bay", Handler: h}
snapshot := slices.Clone([]Route{route})
fmt.Println(len(snapshot)) // 1
parts := strings.Split(route.Pattern, "/")
fmt.Println(parts) // [ bays :bay]
~~~
امشِ بترتيب التسجيل؛ أول نمط مطابق **وطريقته مطابقة** يكسب. لو فيه أنماط مطابقة بس مفيش طريقة، اجمع الطرق ورتّبها بـ ~slices.Sort~ وشيل التكرار بـ ~slices.Compact~ عشان ~Allow~ حتمي.

## Middleware دالة بتلف handler
~Middleware~ نوع دالة ~func(http.Handler) http.Handler~. اللف من الآخر يخلي أول عنصر هو الخارجي؛ يقدر يكتب رد ويرجع من غير استدعاء التالي، زي رفض دخول 403.
~~~go
stamp := func(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  w.Header().Set("X-Dock", "west")
  next.ServeHTTP(w,r)
 })
}
h = stamp(h)
~~~

## قيم الطلب في context
استخدم نوع مفتاح خاص عشان تتجنب التصادم. ~context.WithValue~ بتشتق سياق مع قيمة، و ~r.WithContext~ ترجع نسخة طلب وتحافظ على سياق الأب. ~Context().Value~ ترجع ~any~؛ type assertion مع ~ok~ تتعامل مع القيمة الغائبة من غير panic.
~~~go
type traceKey struct{}
ctx := context.WithValue(r.Context(), traceKey{}, "trace-7")
r = r.WithContext(ctx)
trace, ok := r.Context().Value(traceKey{}).(string)
fmt.Println(trace,ok)
~~~
خزّن params في ~map[string]string~ بمفتاح مختلف عن request ID. متحطش حالة كل الطلبات في متغير عام. المصدر ~nextID~ محقون عشان التست يختار IDs؛ مفيش ساعة أو عشوائية مخفية.
~~~go
h := NewDockRouter(routes, WithRequestID(func() string { return "test-9" }))
fmt.Println(Param(r,"bay"), RequestID(r))
~~~

## أخطاء شائعة
~~~go
for _, mw := range chain { h = mw(h) } // أول عنصر بقى الداخلي
ctx := context.WithValue(context.Background(), "id", id) // ضيّع إلغاء الأب ومفتاح قابل للتصادم
~~~

## المهمة
نفّذ ~NewDockRouter(routes []Route,chain ...Middleware) http.Handler~. خد نسخة من slice التسجيلات؛ العقد بيفترض أنماط صحيحة تبدأ بـ ~/~ وparams بأسماء غير فاضية وفريدة في النمط، وhandlers وmiddleware غير ~nil~. طابق ~URL.Path~ حرفياً بأجزاء؛ parameter يقبل جزء واحد غير فاضي. من غير تنظيف أو redirects أو معالجة query. أول تسجيل يطابق النمط والطريقة يفوز؛ طرق التسجيلات uppercase.

نفّذ ~Param(r *http.Request,name string) string~ من سياق الطلب؛ الغائب نص فاضي. لو مفيش نمط مطابق: 404 و ~no route\n~. لو النمط موجود بس الطريقة مختلفة: 405 و ~wrong method\n~ و ~Allow~ مرتّبة أبجدياً من غير تكرار ومفصولة بـ ~, ~. الأخطاء ~text/plain; charset=utf-8~، و ~HEAD~ فيها من غير جسم؛ مفيش GET ضمني لـ HEAD.

لف كل الردود بالسلسلة؛ أول middleware خارجي، والترتيب قبل handler للأمام وبعده بالعكس. نفّذ ~WithRequestID(nextID func() string) Middleware~ بمصدر غير ~nil~ يرجع ID غير فاضي صالح للهيدر؛ نداء واحد لكل طلب حتى 404/405، تجاهل هيدر العميل، حط ~X-Request-ID~، ومرّر نفس القيمة في السياق. ~RequestID(r *http.Request) string~ ترجع القيمة أو فاضي. حافظ على سياق الأب ومتغيّرش الطلب الأصلي؛ params ما تتسربش بين الطلبات. مفيش شبكة أو ~Report~، و ~main~ فاضية.`,
		Hints: []string{"قارن أجزاء المسار مع النمط قبل تقرير الطريقة.", "ركّب السلسلة من الداخل للخارج عشان أول عنصر يبقى خارجيًا.", "راجع مفاتيح context الخاصة وWithContext على نسخة الطلب."},
		En: &LevelText{LessonRevision: 1, Steps: ch26Lesson263.english, Title: "The Router", Goal: "Route parameterized docks through middleware with per-request identity", Files: []LabFile{{Name: "main.go", Starter: ch263Starter}},
			Lesson: `# A Small Router and a Middleware Chain

## Before you write
- How do you distinguish an unknown path from a disallowed method?
- Which middleware should run first?

## Match rather than clean
~Route~ groups ~Method~, ~Pattern~, and ~Handler~. ~slices.Clone~ snapshots the route slice so edits to its original elements cannot change the router. A segment starting with ~:~ names a parameter. Split pattern and path using ~strings.Split~ and compare segments. Never clean invalid paths into valid ones.
~~~go
route := Route{Method: "GET", Pattern: "/bays/:bay", Handler: h}
snapshot := slices.Clone([]Route{route})
fmt.Println(len(snapshot)) // 1
parts := strings.Split(route.Pattern, "/")
fmt.Println(parts) // [ bays :bay]
~~~
Registration order wins among routes matching both path and method. If only paths match, collect their methods, ~slices.Sort~ them, and remove duplicates with ~slices.Compact~ to build deterministic ~Allow~.

## Middleware wraps a handler
~Middleware~ is ~func(http.Handler) http.Handler~. Wrap in reverse so the first argument is outermost. A wrapper may respond and return without calling next, such as denying access with 403.
~~~go
stamp := func(next http.Handler) http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
  w.Header().Set("X-Dock", "west")
  next.ServeHTTP(w,r)
 })
}
h = stamp(h)
~~~

## Request context values
Use a private key type to prevent collisions. ~context.WithValue~ derives a context with a value. ~r.WithContext~ returns a request copy preserving its parent context. ~Context().Value~ returns ~any~; a comma-ok type assertion safely handles absence.
~~~go
type traceKey struct{}
ctx := context.WithValue(r.Context(), traceKey{}, "trace-7")
r = r.WithContext(ctx)
trace, ok := r.Context().Value(traceKey{}).(string)
fmt.Println(trace,ok)
~~~
Store parameter maps as ~map[string]string~ under a different key from request IDs. Never use global request state. Inject ~nextID~ so tests choose identifiers without hidden clocks or randomness.
~~~go
h := NewDockRouter(routes, WithRequestID(func() string { return "test-9" }))
fmt.Println(Param(r,"bay"), RequestID(r))
~~~

## Common mistakes
~~~go
for _, mw := range chain { h = mw(h) } // First middleware becomes innermost
ctx := context.WithValue(context.Background(), "id", id) // Loses cancellation; collision-prone key
~~~

## Your mission
Implement ~NewDockRouter(routes []Route,chain ...Middleware) http.Handler~. Snapshot the route slice. Assume valid patterns starting with ~/~, nonempty unique parameter names per pattern, and non-~nil~ handlers/middleware. Match ~URL.Path~ segment-for-segment; parameters match exactly one nonempty segment. No cleaning, redirects, or query matching. First registration matching both path and method wins. Registered methods are uppercase.

Implement ~Param(r *http.Request,name string) string~ reading request context, empty when absent. No matching pattern returns 404 and ~no route\n~. Matching paths with no matching method return 405, ~wrong method\n~, and alphabetically sorted unique ~Allow~ methods separated by ~, ~. Errors use ~text/plain; charset=utf-8~ and have no body for ~HEAD~. No implicit GET-to-HEAD matching.

Wrap all responses, including errors. First middleware is outermost: before-handler work runs forward, after-handler work in reverse. Implement ~WithRequestID(nextID func() string) Middleware~ with a non-~nil~ source returning a nonempty header-safe ID. Call it once per request including 404/405, ignore client IDs, set ~X-Request-ID~, and pass that value in context. ~RequestID(r *http.Request) string~ returns the value or empty. Preserve parent context and the original request; parameters cannot leak between requests. No network or ~Report~; leave ~main~ empty.`,
			Hints: []string{"Compare path segments with the pattern before deciding the method.", "Wrap from inside out so the first middleware stays outermost.", "Review private context keys and ~WithContext~ on a request copy."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch263Hidden},
	})
}

const ch263Starter = `package main

import "net/http"

type Route struct {
	Method, Pattern string
	Handler         http.Handler
}
type Middleware func(http.Handler) http.Handler

func NewDockRouter(routes []Route, chain ...Middleware) http.Handler { return http.NotFoundHandler() }
func Param(r *http.Request, name string) string                      { return "" }
func RequestID(r *http.Request) string                               { return "" }
func WithRequestID(nextID func() string) Middleware {
	return func(next http.Handler) http.Handler { return next }
}
func main() {}
`

const ch263Hidden = `package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

// Boundaries: literal/parameter segments, missing/extra segments, method union, wrapping and context isolation.
func TestHiddenRouterParamsMethodsAndContext(t *testing.T) {
	rng := rand.New(rand.NewSource(263))
	sequence := 0
	rid := WithRequestID(func() string { sequence++; return fmt.Sprintf("req-%d", sequence) })
	h := NewDockRouter([]Route{
		{"PUT", "/sectors/:sector/docks/:dock", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Context().Value(17) != "parent" {
				t.Error("lost parent context")
			}
			fmt.Fprintf(w, "%s|%s|%s", Param(r, "sector"), Param(r, "dock"), RequestID(r))
		})},
		{"GET", "/sectors/:sector/docks/:dock", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })},
	}, rid)
	for i := 1; i <= 30; i++ {
		sector, dock := fmt.Sprintf("s%d", rng.Intn(9999)), fmt.Sprintf("d%d", rng.Intn(9999))
		r := httptest.NewRequest("PUT", "/sectors/"+sector+"/docks/"+dock+"?x=1", nil)
		r.Header.Set("X-Request-ID", "untrusted")
		r = r.WithContext(context.WithValue(r.Context(), 17, "parent"))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		id := fmt.Sprintf("req-%d", i)
		if w.Code != 200 || w.Body.String() != sector+"|"+dock+"|"+id || w.Header().Get("X-Request-ID") != id {
			t.Fatal("params/ID", w.Body.String(), w.Header())
		}
		if Param(r, "sector") != "" || RequestID(r) != "" {
			t.Fatal("mutated original request")
		}
	}
	for _, tc := range []struct {
		method, path string
		status       int
		allow, body  string
	}{
		{"POST", "/sectors/a/docks/b", 405, "GET, PUT", "wrong method\n"},
		{"HEAD", "/sectors/a/docks/b", 405, "GET, PUT", ""},
		{"GET", "/sectors//docks/b", 404, "", "no route\n"},
		{"GET", "/sectors/a/docks/b/", 404, "", "no route\n"},
		{"GET", "/SECTORS/a/docks/b", 404, "", "no route\n"},
		{"HEAD", "/missing", 404, "", ""},
	} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status || w.Header().Get("Allow") != tc.allow || w.Body.String() != tc.body || w.Header().Get("X-Request-ID") == "" {
			t.Fatal(tc, w.Code, w.Body.String(), w.Header())
		}
	}
}
func TestHiddenRouterChainOrderShortCircuitAndSnapshot(t *testing.T) {
	log := []string{}
	wrap := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log = append(log, name+"+")
				next.ServeHTTP(w, r)
				log = append(log, name+"-")
			})
		}
	}
	routes := []Route{{"GET", "/dock/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { log = append(log, "handler:"+Param(r, "id")) })}}
	h := NewDockRouter(routes, wrap("outer"), wrap("inner"))
	routes[0].Pattern = "/changed"
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/dock/7", nil))
	if !reflect.DeepEqual(log, []string{"outer+", "inner+", "handler:7", "inner-", "outer-"}) {
		t.Fatal("wrapping/snapshot", log)
	}
	log = nil
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/missing", nil))
	if !reflect.DeepEqual(log, []string{"outer+", "inner+", "inner-", "outer-"}) {
		t.Fatal("middleware skipped on miss", log)
	}
	block := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(403) })
	}
	h = NewDockRouter([]Route{{"GET", "/", http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Error("blocked handler reached") })}}, block)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 403 {
		t.Fatal("short circuit")
	}
	h = NewDockRouter(nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 404 {
		t.Fatal("empty router")
	}
	routes = []Route{
		{"POST", "/dock/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("post")) })},
		{"GET", "/dock/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("first")) })},
		{"GET", "/dock/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("second")) })},
	}
	h = NewDockRouter(routes)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/dock/9", nil))
	if w.Body.String() != "first" {
		t.Fatal("registration and method precedence", w.Body.String())
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("PUT", "/dock/9", nil))
	if w.Code != 405 || w.Header().Get("Allow") != "GET, POST" {
		t.Fatal("unique methods", w.Header())
	}

}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch26Lesson264.arabic, ID: "26-4-envelope", Chapter: 26, Kind: "lab", Title: "الظرف", Goal: "وحّد ردود المنارات واربط أنواع الأخطاء بحالات آمنة",
		Requires: []string{"errtype"}, Bridge: []int{49}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch264Starter}},
		Lesson: `# ظرف واحد لكل رد منارة

## قبل ما تكتب
- هل رسالة خطأ قاعدة البيانات ينفع يقرأها العميل؟
- إزاي تفرّق خطأ معروف حتى لو ملفوف بـ ~%w~؟

## نوع الخطأ مش نصه
~RelayFault~ فيها ~Kind~ للتصنيف و ~Cause~ للتشخيص الداخلي. طريقة ~Error() string~ تحقق واجهة ~error~ من فصل الأخطاء. متستخدمش ناتجها كرسالة HTTP. ~errors.As~ من 5-3 تلاقي النوع حتى جوه wrapper.
~~~go
fault := &RelayFault{Kind: "missing", Cause: errors.New("private store detail")}
var typed *RelayFault
if errors.As(fmt.Errorf("lookup: %w",fault), &typed) {
 fmt.Println(typed.Kind)
}
~~~
لو ضفت ~Unwrap() error~ تقدر ترجع ~Cause~ عشان سلسلة الأخطاء تكمل، زي الغلاف القياسي؛ الطريقة اختيارية للمهمة.
~~~go
func (e *RelayFault) Unwrap() error { return e.Cause }
~~~

## الشكل الثابت
الـ envelope فيها دائماً ~data~ و ~error~. استخدم pointers عشان ~nil~ تتحول لـ JSON ~null~؛ متستخدمش tag ~omitempty~ لأنها بتخفي الحقل. بيانات ~RelayStatus~ اسم وشحنة، مش رسالة التشخيص.
~~~go
type Problem struct { Code string; Message string }
public := Problem{Code: "missing", Message: "relay missing"}
raw, err := json.Marshal(public)
fmt.Println(string(raw),err)
~~~
المثال له أسماء Go افتراضية؛ في المهمة استخدم tags للحروف الصغيرة زي الفصل 18. رجّع status من تصنيف موثوق، مش من ~Cause.Error()~. الحالة 409 تعارض مع الحالة الحالية، زي منارة محظورة.
~~~go
h := NewRelayEnvelope(func(ctx context.Context,name string)(RelayStatus,error) {
 return RelayStatus{Name:name,Charge:70},nil
})
~~~
الـ callback بتاخد نفس سياق الطلب عشان الإلغاء والقيم متضيعش. الـ handler مش بيخترع بيانات وقت فشل المصدر.

## أخطاء شائعة
~~~go
http.Error(w, err.Error(), 500) // تسريب تفاصيل، وشكل غير موحّد
if err.Error() == "missing" { w.WriteHeader(404) } // النص مش هوية النوع
~~~

## المهمة
نفّذ ~RelayFault~ بحقول ~Kind string~ و ~Cause error~ وطريقة ~Error() string~ للتشخيص الداخلي؛ نصها مش جزء من عقد HTTP. نفّذ ~NewRelayEnvelope(fetch func(context.Context,string)(RelayStatus,error)) http.Handler~؛ المصدر غير ~nil~. ~GET /relays/NAME~ لجزء اسم واحد غير فاضي ينادي المصدر مرة بنفس ~r.Context()~ والاسم، ويتجاهل query.

النجاح: 200 و ~{"data":{"name":"...","charge":70},"error":null}~ بقيم المصدر كما هي. الخطأ: ~data:null~ و ~error~ كائن فيه ~code~ و ~message~ فقط. ~*RelayFault~ حتى لو ملفوف: ~invalid~ يعطي 400 و ~invalid / invalid relay~؛ ~missing~ يعطي 404 و ~missing / relay missing~؛ ~blocked~ يعطي 409 و ~blocked / relay blocked~. أي نوع أو تصنيف آخر يعطي 500 و ~internal / relay unavailable~. تجاهل أي بيانات رجعت جنب خطأ؛ ممنوع تسريب السبب أو wrapper.

شكل مسار غير مطابق: 404 بنفس ظرف ~missing~ من غير callback. طريقة غير GET لمسار مطابق: 405 و ~Allow: GET~ وظرف ~method / use GET~ من غير callback. كل الردود ~application/json~، والحقلين موجودين دائماً، و ~HEAD~ من غير جسم. اختبر المسار قبل الطريقة. مفيش شبكة أو وقت أو ~Report~؛ ~main~ فاضية.`,
		Hints: []string{"فرّق بين نوع الخطأ ونص التشخيص اللي يحمله.", "خلّي بيانات الرد والخطأ اختيارية من غير إخفاء الحقول.", "ابدأ بـerrors.As قبل اختيار التصنيف العام."},
		En: &LevelText{LessonRevision: 1, Steps: ch26Lesson264.english, Title: "The Envelope", Goal: "Unify relay responses and map typed errors to safe HTTP statuses", Files: []LabFile{{Name: "main.go", Starter: ch264Starter}},
			Lesson: `# One Envelope for Every Relay Response

## Before you write
- May a client read a database error message?
- How can a known error survive a ~%w~ wrapper?

## Classify types, not text
~RelayFault~ has ~Kind~ for classification and ~Cause~ for internal diagnostics. ~Error() string~ implements the earlier ~error~ interface. Never expose its output as the HTTP message. Earlier ~errors.As~ finds a type inside a wrapper.
~~~go
fault := &RelayFault{Kind: "missing", Cause: errors.New("private store detail")}
var typed *RelayFault
if errors.As(fmt.Errorf("lookup: %w",fault), &typed) {
 fmt.Println(typed.Kind)
}
~~~
Optionally implement ~Unwrap() error~ returning ~Cause~ so traversal continues through the fault, like standard wrappers. It is not required here.
~~~go
func (e *RelayFault) Unwrap() error { return e.Cause }
~~~

## A stable shape
The envelope always contains ~data~ and ~error~. Pointers serialize ~nil~ as JSON ~null~. Do not use ~omitempty~, which removes fields. ~RelayStatus~ contains a name and charge rather than diagnostic text.
~~~go
type Problem struct { Code string; Message string }
public := Problem{Code: "missing", Message: "relay missing"}
raw, err := json.Marshal(public)
fmt.Println(string(raw),err)
~~~
This example uses default Go field names; use lowercase JSON tags in the mission as taught in chapter 18. Select statuses from trusted classification rather than ~Cause.Error()~. Status 409 represents a state conflict, such as a blocked relay.
~~~go
h := NewRelayEnvelope(func(ctx context.Context,name string)(RelayStatus,error) {
 return RelayStatus{Name:name,Charge:70},nil
})
~~~
Pass the same request context to preserve cancellation and values. Never fabricate successful data after source failure.

## Common mistakes
~~~go
http.Error(w, err.Error(), 500) // Leaks internals and breaks the envelope
if err.Error() == "missing" { w.WriteHeader(404) } // Text is not type identity
~~~

## Your mission
Implement ~RelayFault~ with ~Kind string~, ~Cause error~, and diagnostic ~Error() string~. Its text is outside the HTTP contract. Implement ~NewRelayEnvelope(fetch func(context.Context,string)(RelayStatus,error)) http.Handler~ with a non-~nil~ source. ~GET /relays/NAME~ with exactly one nonempty name segment calls it once with the same ~r.Context()~ and name. Ignore query parameters.

Success returns 200 and ~{"data":{"name":"...","charge":70},"error":null}~ preserving source values. Failure returns ~data:null~ and an ~error~ object containing only ~code~ and ~message~. A ~*RelayFault~, even wrapped, maps ~invalid~ to 400 and ~invalid / invalid relay~; ~missing~ to 404 and ~missing / relay missing~; ~blocked~ to 409 and ~blocked / relay blocked~. Any other type or kind returns 500 and ~internal / relay unavailable~. Ignore data returned alongside an error; never leak causes or wrappers.

Unmatched paths return 404 with the ~missing~ envelope without fetching. Non-GET methods on matching paths return 405, ~Allow: GET~, and ~method / use GET~ without fetching. Every response uses ~application/json~ and both envelope fields. ~HEAD~ has no body. Check paths before methods. No network, time, or ~Report~; leave ~main~ empty.`,
			Hints: []string{"Separate an error's type from its diagnostic text.", "Keep response data and error optional without omitting their fields.", "Start with ~errors.As~ before selecting the public category."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch264Hidden},
	})
}

const ch264Starter = `package main

import (
	"context"
	"net/http"
)

type RelayStatus struct {
	Name   string ` + "`" + `json:"name"` + "`" + `
	Charge int    ` + "`" + `json:"charge"` + "`" + `
}
type RelayFault struct {
	Kind  string
	Cause error
}

func (e *RelayFault) Error() string { return "TODO" }
func NewRelayEnvelope(fetch func(context.Context, string) (RelayStatus, error)) http.Handler {
	return http.NotFoundHandler()
}
func main() {}
`

const ch264Hidden = `package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func envelopeCall(h http.Handler, method, path string, ctx context.Context) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(method, path, nil).WithContext(ctx)
	h.ServeHTTP(w, r)
	return w
}
func envelopeAssert(t *testing.T, w *httptest.ResponseRecorder, status int, want string) {
	t.Helper()
	var got, expected any
	if w.Code != status || w.Header().Get("Content-Type") != "application/json" || json.Unmarshal(w.Body.Bytes(), &got) != nil || json.Unmarshal([]byte(want), &expected) != nil || !reflect.DeepEqual(got, expected) {
		t.Fatalf("envelope: %d %s, want %d %s", w.Code, w.Body.String(), status, want)
	}
}

// Boundaries: nil errors, wrapped typed errors, unknown kinds, private causes and route failures.
func TestHiddenEnvelopeSuccessAndErrorMapping(t *testing.T) {
	rng := rand.New(rand.NewSource(264))
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	for i := 0; i < 25; i++ {
		name := fmt.Sprintf("relay-%d", rng.Intn(999999))
		charge := rng.Intn(101)
		calls := 0
		h := NewRelayEnvelope(func(c context.Context, n string) (RelayStatus, error) {
			calls++
			if c != ctx || n != name {
				t.Error("injection lost")
			}
			return RelayStatus{name, charge}, nil
		})
		w := envelopeCall(h, "GET", "/relays/"+name+"?q=1", ctx)
		envelopeAssert(t, w, 200, fmt.Sprintf(` + "`" + `{"data":{"name":%q,"charge":%d},"error":null}` + "`" + `, name, charge))
		if calls != 1 {
			t.Fatal("fetch count")
		}
	}
	for _, tc := range []struct {
		kind          string
		status        int
		code, message string
	}{
		{"invalid", 400, "invalid", "invalid relay"}, {"missing", 404, "missing", "relay missing"}, {"blocked", 409, "blocked", "relay blocked"}, {"database", 500, "internal", "relay unavailable"},
	} {
		for _, wrapped := range []bool{false, true} {
			secret := fmt.Sprintf("private-dsn-%d", rng.Intn(999999))
			var fault error = &RelayFault{tc.kind, errors.New(secret)}
			if wrapped {
				fault = fmt.Errorf("backend: %w", fault)
			}
			h := NewRelayEnvelope(func(context.Context, string) (RelayStatus, error) { return RelayStatus{"must-not-leak", 99}, fault })
			w := envelopeCall(h, "GET", "/relays/x", context.Background())
			envelopeAssert(t, w, tc.status, fmt.Sprintf(` + "`" + `{"data":null,"error":{"code":%q,"message":%q}}` + "`" + `, tc.code, tc.message))
			if strings.Contains(w.Body.String(), secret) {
				t.Fatal("internal leak")
			}
		}
	}
	h := NewRelayEnvelope(func(context.Context, string) (RelayStatus, error) { return RelayStatus{}, errors.New("secret") })
	envelopeAssert(t, envelopeCall(h, "GET", "/relays/x", context.Background()), 500, ` + "`" + `{"data":null,"error":{"code":"internal","message":"relay unavailable"}}` + "`" + `)
}
func TestHiddenEnvelopeRejectsRoutesWithoutFetching(t *testing.T) {
	h := NewRelayEnvelope(func(context.Context, string) (RelayStatus, error) {
		t.Error("unexpected fetch")
		return RelayStatus{}, nil
	})
	for _, path := range []string{"/relays", "/relays/", "/relays/a/b", "/missing"} {
		envelopeAssert(t, envelopeCall(h, "GET", path, context.Background()), 404, ` + "`" + `{"data":null,"error":{"code":"missing","message":"relay missing"}}` + "`" + `)
	}
	for _, method := range []string{"POST", "PUT", "DELETE", "HEAD"} {
		w := envelopeCall(h, method, "/relays/x", context.Background())
		if w.Header().Get("Allow") != "GET" || w.Code != 405 {
			t.Fatal("method", w.Code, w.Header())
		}
		if method == "HEAD" {
			if w.Body.Len() != 0 {
				t.Fatal("HEAD body")
			}
		} else {
			envelopeAssert(t, w, 405, ` + "`" + `{"data":null,"error":{"code":"method","message":"use GET"}}` + "`" + `)
		}
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch26Lesson265.arabic, ID: "26-5-static", Chapter: 26, Kind: "lab", Title: "ملفات ثابتة", Goal: "اخدم خرائط مضمّنة من غير كشف مجلدات وبكاش حتمي",
		Requires: []string{"sel:http.FileServer"}, Bridge: []int{52}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch265Starter}, {Name: "site/panel.html", Starter: `<!doctype html><title>Orbital Console</title><h1>Dock charts</h1>
`, ReadOnly: true}, {Name: "site/charts/notice.txt", Starter: `Airlock inspections: deck 7.
`, ReadOnly: true}, {Name: "site/charts/sector.json", Starter: `{"sector":"Kepler","docks":7}
`, ReadOnly: true}, {Name: "site/missing.html", Starter: `<!doctype html><title>Lost chart</title><p>Chart unavailable</p>
`, ReadOnly: true}},
		Lesson: `# خرائط ثابتة جوه البرنامج

## قبل ما تكتب
- هتمنع عرض أسماء الملفات في مجلد إزاي؟
- إيه اللي يضمن إن بصمة الكاش تخص نفس البايتات المعروضة؟

## Embed مش embedding بتاع structs
استورد ~embed~. التوجيه ~//go:embed~ بيضم الملفات وقت البناء، والمتغير ~embed.FS~ نظام ملفات للقراءة بس؛ اسم التوجيه لازم فوق المتغير. ده مختلف عن fields المضمّنة من فصل 7. الملفات موجودة في الـ starter للقراءة فقط.
~~~go
//go:embed site
var chartSite embed.FS
root, err := fs.Sub(chartSite, "site")
body, err := fs.ReadFile(root, "charts/notice.txt")
~~~
~fs.Sub~ من ~io/fs~ تخلي مجلد نقطة البداية، و ~fs.ReadFile~ تقرأ بايتات ملف من ~fs.FS~ بدل قرص الجهاز. ملفات الحزمة المضمّنة معروفة وقت البناء.

## FileServer وطبقة السماح
~http.FS~ توصل نظام الملفات بواجهة HTTP، و ~http.FileServer~ ترجع handler يخدم الملفات. من غير غلاف ممكن تعرض قوائم المجلدات أو تعمل redirects. استخدم خريطة مسارات مسموحة بالظبط قبل التفويض؛ متعملش تنظيف لمسار مهاجم.
~~~go
files := http.FileServer(http.FS(root))
copyRequest := r.Clone(r.Context())
copyRequest.URL.Path = "/charts/notice.txt"
files.ServeHTTP(w,copyRequest)
~~~
~r.Clone~ تنسخ الطلب والهيدر والـ URL، فالتغيير الداخلي ما يلمسش الأصل. ~URL.RawPath~ هو شكل المسار escaped لما يبقى مختلف عن ~URL.Path~؛ ارفض الشكل البديل هنا بدل تطبيعه.
~~~go
if r.URL.RawPath != "" && r.URL.RawPath != r.URL.Path { w.WriteHeader(404); return }
~~~

## سياسة الكاش
~ETag~ بصمة قوية محاطة بعلامتي quote. استخدم ~sha256.Sum256~ من فصل الخزنة لنفس بايتات الملف، و ~Cache-Control~ لتحديد صلاحية 600 ثانية. ~Content-Length~ طول البايتات. ~HEAD~ زي GET في الهيدر بس من غير جسم.
~~~go
tag := fmt.Sprintf("\"%x\"", sha256.Sum256(body))
w.Header().Set("ETag",tag)
w.Header().Set("Cache-Control","public, max-age=600")
~~~
النجمة الإضافية بتستخدم ~If-None-Match~: تطابق تام يرجع 304 من غير جسم؛ السياسة هنا متعمّدة إنها أضيق من HTTP العامة. ~Header.Del~ تشيل شرط من نسخة الطلب قبل ~FileServer~ عشان ميطبقش سياسة أوسع من الغلاف.
~~~go
copyRequest.Header.Del("If-None-Match")
copyRequest.Header.Del("If-Modified-Since")
~~~

## أخطاء شائعة
~~~go
return http.FileServer(http.Dir(".")) // يكشف ملفات الجهاز وقوائم المجلدات
w.Header().Set("ETag", fmt.Sprint(time.Now())) // مش بصمة المحتوى
~~~

## المهمة
نفّذ ~NewChartFiles() http.Handler~ باستخدام ~embed~ و ~http.FileServer~ والملفات الثابتة للقراءة. المسارات الوحيدة: ~/console~ يخدم ~site/panel.html~ بنوع ~text/html; charset=utf-8~؛ ~/charts/notice.txt~ يخدم ~site/charts/notice.txt~ بنوع ~text/plain; charset=utf-8~؛ ~/charts/sector.json~ يخدم ~site/charts/sector.json~ بنوع ~application/json~. ~GET~ يرجع 200 والبايتات بالظبط، و ~HEAD~ نفس الحالة والهيدر بجسم فاضي. حط ~Content-Length~ الصحيح و ~ETag~ كـ SHA-256 lowercase hex بين quotes و ~Cache-Control: public, max-age=600~.

أي مسار تاني، بما فيه المجلدات و ~..~ والشرطات encoded وأي ~RawPath~ مختلف، يرجع 404 بجسم ~site/missing.html~ ونوع HTML و ~Cache-Control: no-store~؛ ~HEAD~ من غير جسم. افحص المسار قبل الطريقة. طريقة غير GET/HEAD لمسار معروف ترجع 405 و ~Allow: GET, HEAD~ و ~read only\n~ بنوع نص عادي و ~no-store~. تجاهل query، ومتعملش listing أو redirect أو قراءة من قرص الجهاز. الطلب الأصلي يفضل ثابت. كل القراءات آمنة بالتزامن؛ مفيش شبكة أو ~Report~، و ~main~ فاضية.

النجمة التالتة: ~If-None-Match~ لو مساوي البصمة بالظبط على مسار صالح مع GET/HEAD يرجع 304 و ~ETag~ و ~Cache-Control~ بس من غير جسم أو ~Content-Type~ أو ~Content-Length~. النجمة وweak tag وقايمة tags وأي قيمة مختلفة ترجع الرد الكامل. تجاهل ~If-Modified-Since~؛ صفحة 404 عمرها ما تتحول لـ 304.`,
		Hints: []string{"عامل كل مسار عام كاختيار صريح من قائمة قصيرة.", "اقرأ المحتوى من filesystem بدل مسار الجهاز.", "قارن ~RawPath~ قبل إعادة كتابة ~URL~."},
		En: &LevelText{LessonRevision: 1, Steps: ch26Lesson265.english, Title: "Static Files", Goal: "Serve embedded charts without directory disclosure and with deterministic caching", Files: []LabFile{{Name: "main.go", Starter: ch265Starter}, {Name: "site/panel.html", Starter: `<!doctype html><title>Orbital Console</title><h1>Dock charts</h1>
`, ReadOnly: true}, {Name: "site/charts/notice.txt", Starter: `Airlock inspections: deck 7.
`, ReadOnly: true}, {Name: "site/charts/sector.json", Starter: `{"sector":"Kepler","docks":7}
`, ReadOnly: true}, {Name: "site/missing.html", Starter: `<!doctype html><title>Lost chart</title><p>Chart unavailable</p>
`, ReadOnly: true}},
			Lesson: `# Charts Embedded in the Program

## Before you write
- How will you prevent directory listings?
- What ensures the cache tag describes the exact served bytes?

## Embed is not struct embedding
Import ~embed~. ~//go:embed~ includes files at build time; ~embed.FS~ is read-only. Place the directive above its variable. This differs from chapter 7's embedded struct fields. Starter assets are read-only.
~~~go
//go:embed site
var chartSite embed.FS
root, err := fs.Sub(chartSite, "site")
body, err := fs.ReadFile(root, "charts/notice.txt")
~~~
~fs.Sub~ from ~io/fs~ selects a filesystem root. ~fs.ReadFile~ reads from ~fs.FS~ rather than the machine's disk. Embedded package files are known at build time.

## FileServer behind an allowlist
~http.FS~ adapts a filesystem to HTTP; ~http.FileServer~ returns a file-serving handler. Without a wrapper it may list directories or redirect. Allow exact public routes before delegation; never clean attacker paths.
~~~go
files := http.FileServer(http.FS(root))
copyRequest := r.Clone(r.Context())
copyRequest.URL.Path = "/charts/notice.txt"
files.ServeHTTP(w,copyRequest)
~~~
~r.Clone~ copies the request, headers, and URL so internal rewriting leaves the original untouched. ~URL.RawPath~ preserves an escaped spelling when different from ~URL.Path~; reject alternate encodings here.
~~~go
if r.URL.RawPath != "" && r.URL.RawPath != r.URL.Path { w.WriteHeader(404); return }
~~~

## Cache policy
~ETag~ is a strong quoted content fingerprint. Use chapter 23's ~sha256.Sum256~ on exact bytes, and ~Cache-Control~ for a 600-second lifetime. ~Content-Length~ counts bytes. ~HEAD~ has GET headers without its body.
~~~go
tag := fmt.Sprintf("\"%x\"", sha256.Sum256(body))
w.Header().Set("ETag",tag)
w.Header().Set("Cache-Control","public, max-age=600")
~~~
The bonus uses ~If-None-Match~: exact equality produces bodyless 304. This deliberately narrow policy differs from general HTTP matching. ~Header.Del~ removes conditions from the cloned request before ~FileServer~ applies its broader policy.
~~~go
copyRequest.Header.Del("If-None-Match")
copyRequest.Header.Del("If-Modified-Since")
~~~

## Common mistakes
~~~go
return http.FileServer(http.Dir(".")) // Exposes machine files and listings
w.Header().Set("ETag", fmt.Sprint(time.Now())) // Not a content fingerprint
~~~

## Your mission
Implement ~NewChartFiles() http.Handler~ using ~embed~, ~http.FileServer~, and read-only assets. Only these public routes exist: ~/console~ serves ~site/panel.html~ as ~text/html; charset=utf-8~; ~/charts/notice.txt~ serves ~site/charts/notice.txt~ as ~text/plain; charset=utf-8~; ~/charts/sector.json~ serves ~site/charts/sector.json~ as ~application/json~. ~GET~ returns 200 and exact bytes; ~HEAD~ returns the same status and headers without a body. Include exact ~Content-Length~, a quoted lowercase SHA-256 hex ~ETag~, and ~Cache-Control: public, max-age=600~.

Every other path, including directories, ~..~, encoded separators, and differing ~RawPath~, returns 404 with ~site/missing.html~, HTML type, and ~Cache-Control: no-store~; ~HEAD~ is bodyless. Check paths before methods. Non-GET/HEAD on known routes returns 405, ~Allow: GET, HEAD~, plain-text ~read only\n~, and ~no-store~. Ignore queries; no listings, redirects, or machine-disk reads. Preserve the original request. All reads are concurrent-safe. No network or ~Report~; leave ~main~ empty.

Third star: an exactly matching ~If-None-Match~ on valid GET/HEAD returns 304 with only ~ETag~ and ~Cache-Control~, without body, ~Content-Type~, or ~Content-Length~. Wildcard, weak, list, and differing tags return the full response. Ignore ~If-Modified-Since~; 404 never becomes 304.`,
			Hints: []string{"Treat each public path as an explicit choice from a short list.", "Read assets from a filesystem instead of a machine path.", "Compare ~RawPath~ before rewriting the ~URL~."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch265Hidden},
	})
}

const ch265Starter = `package main

import (
	"embed"
	"net/http"
)

//go:embed site
var chartSite embed.FS

func NewChartFiles() http.Handler { return http.NotFoundHandler() }
func main()                       {}
`

const ch265Hidden = `package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
)

var chartCases = []struct{ path, file, media string }{
	{"/console", "site/panel.html", "text/html; charset=utf-8"},
	{"/charts/notice.txt", "site/charts/notice.txt", "text/plain; charset=utf-8"},
	{"/charts/sector.json", "site/charts/sector.json", "application/json"},
}

func chartRequest(h http.Handler, method, path, tag string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, nil)
	r.Header.Set("If-None-Match", tag)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

// Boundaries: exact resources, directory/traversal/encoded paths, HEAD, immutable parallel reads.
func TestHiddenChartsEmbeddedBytesAndHeaders(t *testing.T) {
	h := NewChartFiles()
	rng := rand.New(rand.NewSource(265))
	for i := 0; i < 30; i++ {
		tc := chartCases[rng.Intn(len(chartCases))]
		body, e := os.ReadFile(tc.file)
		if e != nil {
			t.Fatal(e)
		}
		tag := fmt.Sprintf("\"%x\"", sha256.Sum256(body))
		for _, method := range []string{"GET", "HEAD"} {
			w := chartRequest(h, method, tc.path+"?cycle=3", "")
			if w.Code != 200 || w.Header().Get("Content-Type") != tc.media || w.Header().Get("Content-Length") != strconv.Itoa(len(body)) || w.Header().Get("ETag") != tag || w.Header().Get("Cache-Control") != "public, max-age=600" {
				t.Fatal("asset headers", tc, w.Code, w.Header())
			}
			if (method == "GET" && w.Body.String() != string(body)) || (method == "HEAD" && w.Body.Len() != 0) {
				t.Fatal("asset body")
			}
		}
	}
	missing, e := os.ReadFile("site/missing.html")
	if e != nil {
		t.Fatal(e)
	}
	for _, path := range []string{"/", "/charts", "/charts/", "/panel.html", "/missing.html", "/charts/../panel.html", "/charts/%2e%2e/panel.html", "/charts%2fnotice.txt", "/charts/%252e%252e/panel.html", "/charts/absent.txt"} {
		for _, method := range []string{"GET", "HEAD", "POST"} {
			w := chartRequest(h, method, path, "")
			if w.Code != 404 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Content-Type") != "text/html; charset=utf-8" {
				t.Fatal("custom missing", path, w.Code, w.Header())
			}
			if method == "HEAD" {
				if w.Body.Len() != 0 {
					t.Fatal("head missing body")
				}
			} else if w.Body.String() != string(missing) {
				t.Fatal("custom page", path, w.Body.String())
			}
		}
	}
	w := chartRequest(h, "POST", "/console", "")
	if w.Code != 405 || w.Header().Get("Allow") != "GET, HEAD" || w.Header().Get("Cache-Control") != "no-store" || w.Body.String() != "read only\n" {
		t.Fatal("method", w.Code, w.Header())
	}
}
func TestHiddenChartsRequestOwnership(t *testing.T) {
	h := NewChartFiles()
	r := httptest.NewRequest("GET", "/console?crew=1", nil)
	r.Header.Set("If-None-Match", "different")
	r.Header.Set("If-Modified-Since", "Wed, 01 Jan 2098 00:00:00 GMT")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 || r.URL.Path != "/console" || r.URL.RawQuery != "crew=1" || r.Header.Get("If-None-Match") != "different" || r.Header.Get("If-Modified-Since") == "" {
		t.Fatal("request ownership/conditions", w.Code, r)
	}
}
func TestHiddenChartsConcurrentReads(t *testing.T) {
	h := NewChartFiles()
	var wg sync.WaitGroup
	for i := 0; i < 60; i++ {
		tc := chartCases[i%len(chartCases)]
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := chartRequest(h, "GET", tc.path, "")
			if w.Code != 200 || w.Body.Len() == 0 {
				t.Error("parallel asset")
			}
		}()
	}
	wg.Wait()
}
func TestHiddenBonusChartsExactConditionalPolicy(t *testing.T) {
	h := NewChartFiles()
	for _, tc := range chartCases {
		tag := chartRequest(h, "GET", tc.path, "").Header().Get("ETag")
		for _, method := range []string{"GET", "HEAD"} {
			w := chartRequest(h, method, tc.path, tag)
			if w.Code != 304 || w.Body.Len() != 0 || w.Header().Get("ETag") != tag || w.Header().Get("Cache-Control") != "public, max-age=600" || w.Header().Get("Content-Type") != "" || w.Header().Get("Content-Length") != "" {
				t.Fatal("304 metadata", w.Code, w.Header())
			}
		}
		for _, other := range []string{"*", "W/" + tag, tag + ", \"other\"", "\"wrong\""} {
			w := chartRequest(h, "GET", tc.path, other)
			if w.Code != 200 || w.Body.Len() == 0 {
				t.Fatal("nonexact conditional", other, w.Code)
			}
		}
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch26Lesson266.arabic, ID: "26-6-upload", Chapter: 26, Kind: "lab", Title: "رفع ملفات", Goal: "ارفع لقطات شحنات بحدود stream وتخزين آمن وتنظيف عند الفشل",
		Requires: []string{"sel:http.MaxBytesReader"}, Bridge: []int{51}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch266Starter}},
		Lesson: `# سكانر الشحنات يرفع لقطة بأمان

## قبل ما تكتب
- هتحسب حد الطلب كله ولا بايتات الملف بس؟
- لو القراءة قطعت في النص، مين يمسح الملف الجزئي؟

## Multipart بالـ stream
~r.MultipartReader~ تقرأ جسم ~multipart/form-data~ من غير تحميل الملفات كلها للذاكرة. ~NextPart~ ترجع جزء قارئ أو ~io.EOF~ بعد آخر جزء. اقرا جزء واحد بـ ~io.Copy~ في ملف مؤقت، وبعدين اطلب التالي عشان ترفض الأجزاء الزيادة. بلاش ~ParseMultipartForm~ في المهمة.
~~~go
r.Body = http.MaxBytesReader(w,r.Body,4096)
defer r.Body.Close()
reader, err := r.MultipartReader()
if err != nil { return }
part, err := reader.NextPart()
if err != nil { return }
defer part.Close()
~~~
~http.MaxBytesReader~ تحد الجسم الكلي وترجع ~*http.MaxBytesError~ لما القراءة تتعداه؛ اكتشفه بـ ~errors.As~ واختار 413. ~os.PathError~ بتميّز أخطاء نظام الملفات عشان ترجع 500 بدل خطأ طلب. ~io.LimitReader~ من مهمة REST تحد الملف لوحده؛ حد زائد 1 يكشف الزيادة. بعد نهاية multipart صرّف باقي ~r.Body~ بـ ~io.Copy(io.Discard,r.Body)~ عشان epilogue يتحسب كمان.
~~~go
var oversized *http.MaxBytesError
if errors.As(err,&oversized) { w.WriteHeader(413) }
var disk *os.PathError
if errors.As(err,&disk) { w.WriteHeader(500) }
~~~

## الاسم المعلن مش مسار
~part.FileName()~ ممكن تنظّف المسار وتخفي محاولة traversal. اقرا ~part.Header.Get("Content-Disposition")~ وفكّه بـ ~mime.ParseMediaType~، ثم افحص ~params["filename"]~ الخام و ~params["name"]~ قبل أي إنشاء ملف. ~filename~ هنا ASCII محدود؛ مش اسم التخزين النهائي.
~~~go
_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
fmt.Println(params["name"],params["filename"],err)
~~~
متثقش في الامتداد أو content type المعلن. اللقطة لازم تبدأ ببايتات PNG الثمانية؛ المهمة بتفحص التوقيع بس، مش سلامة صورة كاملة.
~~~go
signature := []byte{137,80,78,71,13,10,26,10}
fmt.Println(len(signature)) // 8
~~~

## تجهيز ثم نشر
~os.CreateTemp~ من فصل الملفات تنشئ ملف فريد داخل الجذر المحقون. ~defer os.Remove~ تمسح الجزئي في كل مسار. بعد النسخ، ~Seek(0,io.SeekStart)~ ترجع للبداية، و ~io.ReadFull~ تقرأ التوقيع، و ~bytes.Equal~ تقارنه. اقفل الملف قبل ~os.Rename~.
~~~go
stage, err := os.CreateTemp(root,".scan-*")
if err != nil { return err }
defer os.Remove(stage.Name())
defer stage.Close()
_, err = stage.Seek(0,io.SeekStart)
_, err = os.Lstat(filepath.Join(root,"0123456789abcdef.png"))
~~~
المصدر ~NextID~ محقون، وبيتنده تحت ~sync.Mutex~ بعد التحقق الكامل. ~os.Lstat~ تفحص وجود الاسم نفسه من غير اتباع symlink. نفس القفل يشمل ~os.Lstat~ للوجهة و ~os.Rename~؛ وجود أي ملف أو symlink تعارض 409 من غير overwrite. الجذر خاص بالتطبيق ومفيش تعديل خارجي عليه؛ القفل يحمي طلبات نفس الـ handler فقط. الاسم النهائي token صالح، مش اسم العميل.
~~~go
cfg := ScanUploadConfig{Root:root, BodyLimit:4096, FileLimit:1024,
 NextID:func()(string,error){ return "0123456789abcdef",nil }}
h, err := NewScanUpload(cfg)
~~~

## أخطاء شائعة
~~~go
name := part.FileName() // دليل المسار الخام ممكن يكون اتشال
os.Create(filepath.Join(root,name)) // اسم العميل ممكن يخرج من الجذر ويمسح ملف
io.ReadAll(r.Body) // ذاكرة غير محدودة قبل تطبيق الحد
~~~

## المهمة
نفّذ ~NewScanUpload(cfg ScanUploadConfig) (http.Handler,error)~. ~Root~ مجلد موجود خاص بالتطبيق؛ ~BodyLimit~ بين 1 و 2097152، و ~FileLimit~ بين 8 و 1048576 شاملة، و ~NextID~ غير ~nil~. إعداد غلط يرجع ~nil,error~. مفيش كتابة خارج الجذر؛ الاختبارات كلها في ~t.TempDir()~.

~POST /cargo-scan~ يقبل multipart بجزء واحد فقط اسمه ~scan~ و ~filename~ خام طوله 1..80 بايت: أول حرف إنجليزي أو رقم، والباقي حروف إنجليزي أو أرقام أو ~.~ أو ~_~ أو ~-~، وممنوع ~..~ في أي موضع. ارفض اسم أو شكل بايظ بـ 400. تجاهل امتداد وcontent type الجزء؛ اقبل فقط توقيع PNG الثماني، حتى لو الملف اسمه ~camera.bin~. الحدين شاملين؛ زيادة الكلي أو الملف 413. أخطاء القراءة أو القطع أو الإلغاء 400. آخر جزء لازم ينتهي صح؛ صرّف epilogue عشان الحد الكلي يشمله.

بعد التحقق وقبل النشر، افحص إلغاء السياق. نادِ ~NextID~ مرة تحت قفل النشر؛ لازم يرجع 16 حرف lowercase hex. خطأ المصدر أو token غلط أو خطأ نظام الملفات 500؛ ملف نهائي موجود 409. احفظ بايتات الملف بالظبط باسم ~TOKEN.png~ داخل الجذر من خلال ملف مؤقت ثم rename؛ امسح المؤقت عند أي فشل. المصدر مش بيتنده لطلب مرفوض قبل مرحلة النشر، والملفات الموجودة لا تتغير. طلبات نفس الـ handler متزامنة بأمان، والمصدر مش لازم يحمي نفسه؛ ممنوع كائنين يكتبوا نفس الجذر.

النجاح 201 و ~TOKEN.png\n~. رسائل الأخطاء: 400 ~invalid scan\n~؛ 409 ~scan exists\n~؛ 413 ~scan too large\n~؛ 415 للتوقيع المرفوض ~PNG required\n~؛ 500 ~scan unavailable\n~. كلها ~text/plain; charset=utf-8~ من غير أسباب داخلية. مسار مختلف 404 و ~unknown upload\n~، وطريقة غير POST للمسار المعروف 405 و ~Allow: POST~ و ~use POST\n~. ~HEAD~ جسمه فاضي، والـ query تتجاهلها. مفيش serving للرفع أو شبكة أو وقت أو ~Report~؛ ~main~ فاضية.`,
		Hints: []string{"فكّر في كل جزء كقارئ له نهاية وخطأ محتمل.", "افصل ملف مؤقت عن الاسم النهائي المنشور.", "ابدأ تطبيق ~MaxBytesReader~ قبل إنشاء ~MultipartReader~."},
		En: &LevelText{LessonRevision: 1, Steps: ch26Lesson266.english, Title: "File Upload", Goal: "Stream cargo scans with size limits, safe storage, and failure cleanup", Files: []LabFile{{Name: "main.go", Starter: ch266Starter}},
			Lesson: `# Safely Upload a Cargo Scan

## Before you write
- Will you cap the complete request or only file bytes?
- Who removes partial files after an interrupted read?

## Streaming multipart
~r.MultipartReader~ streams ~multipart/form-data~ without loading complete files into memory. ~NextPart~ returns a readable part or ~io.EOF~ after the last one. Stream one part with ~io.Copy~, then request another to reject extra parts. Do not use ~ParseMultipartForm~ here.
~~~go
r.Body = http.MaxBytesReader(w,r.Body,4096)
defer r.Body.Close()
reader, err := r.MultipartReader()
if err != nil { return }
part, err := reader.NextPart()
if err != nil { return }
defer part.Close()
~~~
~http.MaxBytesReader~ bounds the entire body and returns ~*http.MaxBytesError~ beyond the cap. Detect it with ~errors.As~ for 413. ~os.PathError~ identifies filesystem failures for 500 rather than a request error. REST's ~io.LimitReader~ bounds file bytes separately; cap plus one detects overflow. After multipart completion, drain ~r.Body~ using ~io.Copy(io.Discard,r.Body)~ so the epilogue also counts.
~~~go
var oversized *http.MaxBytesError
if errors.As(err,&oversized) { w.WriteHeader(413) }
var disk *os.PathError
if errors.As(err,&disk) { w.WriteHeader(500) }
~~~

## A declared name is not a path
~part.FileName()~ can normalize away traversal evidence. Read ~part.Header.Get("Content-Disposition")~, parse using ~mime.ParseMediaType~, then validate raw ~params["filename"]~ and ~params["name"]~ before creating files. The filename here follows a narrow ASCII policy and never becomes the storage name.
~~~go
_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
fmt.Println(params["name"],params["filename"],err)
~~~
Do not trust extensions or declared content types. Require the eight PNG signature bytes. This mission checks the signature, not full image validity.
~~~go
signature := []byte{137,80,78,71,13,10,26,10}
fmt.Println(len(signature)) // 8
~~~

## Stage, then publish
Earlier ~os.CreateTemp~ creates a unique file inside the injected root. ~defer os.Remove~ cleans partial files on every exit. After copying, ~Seek(0,io.SeekStart)~ rewinds, ~io.ReadFull~ reads the signature, and ~bytes.Equal~ compares it. Close before ~os.Rename~.
~~~go
stage, err := os.CreateTemp(root,".scan-*")
if err != nil { return err }
defer os.Remove(stage.Name())
defer stage.Close()
_, err = stage.Seek(0,io.SeekStart)
_, err = os.Lstat(filepath.Join(root,"0123456789abcdef.png"))
~~~
Inject ~NextID~ and call it under ~sync.Mutex~ after complete validation. ~os.Lstat~ inspects the name itself without following symlinks. The same lock covers destination ~os.Lstat~ and ~os.Rename~. Any existing file or symlink means 409 without overwriting. The application privately owns the root, with no external mutation; the mutex protects requests to this handler only. Store by a validated token, never the client filename.
~~~go
cfg := ScanUploadConfig{Root:root, BodyLimit:4096, FileLimit:1024,
 NextID:func()(string,error){ return "0123456789abcdef",nil }}
h, err := NewScanUpload(cfg)
~~~

## Common mistakes
~~~go
name := part.FileName() // Raw traversal evidence may already be removed
os.Create(filepath.Join(root,name)) // Client path may escape or overwrite
io.ReadAll(r.Body) // Unbounded memory before applying limits
~~~

## Your mission
Implement ~NewScanUpload(cfg ScanUploadConfig) (http.Handler,error)~. ~Root~ is an existing privately owned directory. ~BodyLimit~ is inclusive 1..2097152; ~FileLimit~ inclusive 8..1048576; ~NextID~ is non-~nil~. Invalid configuration returns ~nil,error~. Write only inside the root; tests use ~t.TempDir()~.

~POST /cargo-scan~ accepts exactly one multipart part named ~scan~. Its raw filename is 1..80 bytes: first an English letter or digit, then English letters, digits, ~.~, ~_~, or ~-~, with no ~..~ anywhere. Invalid names or structure return 400. Ignore part extension and content type; accept only the eight PNG signature bytes, even for ~camera.bin~. Both caps are inclusive; exceeding either returns 413. Read failures, truncation, and cancellation return 400. Require a complete closing boundary and drain the epilogue to enforce the whole-body cap.

Check context cancellation after validation and before publishing. Call ~NextID~ once under the publication lock; require 16 lowercase hex characters. Source errors, invalid tokens, and filesystem errors return 500. Existing destinations return 409. Store exact bytes as ~TOKEN.png~ inside root using a temporary file then rename; remove temporary files on failure. Requests rejected before publication never call the source. Existing files stay unchanged. Concurrent requests to one handler are safe; the source need not synchronize itself. Multiple handlers must not share a root.

Success returns 201 and ~TOKEN.png\n~. Error bodies: 400 ~invalid scan\n~; 409 ~scan exists\n~; 413 ~scan too large\n~; rejected signature 415 ~PNG required\n~; 500 ~scan unavailable\n~. All use ~text/plain; charset=utf-8~ without private causes. Unknown paths return 404 and ~unknown upload\n~; non-POST on the known path returns 405, ~Allow: POST~, and ~use POST\n~. ~HEAD~ is bodyless; ignore queries. Do not serve uploads or use networking, time, or ~Report~. Leave ~main~ empty.`,
			Hints: []string{"Treat each part as a reader with an end and possible error.", "Separate a temporary file from the published name.", "Apply ~MaxBytesReader~ before creating ~MultipartReader~."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch266Hidden},
	})
}

const ch266Starter = `package main

import (
	"errors"
	"net/http"
)

type ScanUploadConfig struct {
	Root                 string
	BodyLimit, FileLimit int64
	NextID               func() (string, error)
}

func NewScanUpload(cfg ScanUploadConfig) (http.Handler, error) { return nil, errors.New("TODO") }
func main()                                                    {}
`

const ch266Hidden = `package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var scanSignature = []byte{137, 80, 78, 71, 13, 10, 26, 10}

func scanBody(t *testing.T, filename, field string, payload []byte, extra bool) ([]byte, string) {
	t.Helper()
	var raw bytes.Buffer
	mw := multipart.NewWriter(&raw)
	if err := mw.SetBoundary("orbital-scan-boundary"); err != nil {
		t.Fatal(err)
	}
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(` + "`" + `form-data; name=%q; filename=%q` + "`" + `, field, filename))
	header.Set("Content-Type", "application/octet-stream")
	part, e := mw.CreatePart(header)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = part.Write(payload); e != nil {
		t.Fatal(e)
	}
	if extra {
		if e = mw.WriteField("extra", "no"); e != nil {
			t.Fatal(e)
		}
	}
	if e = mw.Close(); e != nil {
		t.Fatal(e)
	}
	return raw.Bytes(), mw.FormDataContentType()
}
func scanCall(h http.Handler, raw []byte, media string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", "/cargo-scan", bytes.NewReader(raw))
	r.Header.Set("Content-Type", media)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func scanEmpty(t *testing.T, root string) {
	t.Helper()
	entries, e := os.ReadDir(root)
	if e != nil || len(entries) != 0 {
		t.Fatal("partial files", entries, e)
	}
}
func scanConstruct(t *testing.T, cfg ScanUploadConfig) http.Handler {
	t.Helper()
	h, e := NewScanUpload(cfg)
	if e != nil || h == nil {
		t.Fatal("construct", e)
	}
	return h
}

// Boundaries: exact file/body caps, raw filename traversal, malformed/truncated parts, type and source failures.
func TestHiddenScanStreamLimitsAndContents(t *testing.T) {
	rng := rand.New(rand.NewSource(266))
	for _, size := range []int{8, 9, 127, 1024} {
		root := t.TempDir()
		payload := make([]byte, size)
		rng.Read(payload)
		copy(payload, scanSignature)
		raw, media := scanBody(t, "camera.bin", "scan", payload, false)
		calls := 0
		cfg := ScanUploadConfig{root, int64(len(raw)), int64(size), func() (string, error) { calls++; return "abcdef0123456789", nil }}
		h := scanConstruct(t, cfg)
		w := scanCall(h, raw, media)
		if w.Code != 201 || w.Body.String() != "abcdef0123456789.png\n" || calls != 1 {
			t.Fatal("exact upload", w.Code, w.Body.String(), calls)
		}
		stored, e := os.ReadFile(filepath.Join(root, "abcdef0123456789.png"))
		if e != nil || !bytes.Equal(stored, payload) {
			t.Fatal("stored bytes", e)
		}
		entries, _ := os.ReadDir(root)
		if len(entries) != 1 {
			t.Fatal("temporary residue")
		}
		// A repeated token must never overwrite the first upload.
		w = scanCall(h, raw, media)
		if w.Code != 409 {
			t.Fatal("collision", w.Code)
		}
		stored, _ = os.ReadFile(filepath.Join(root, "abcdef0123456789.png"))
		if !bytes.Equal(stored, payload) {
			t.Fatal("overwritten")
		}
		for _, limits := range [][2]int64{{int64(len(raw) - 1), int64(size)}, {int64(len(raw) + 100), int64(size - 1)}} {
			if limits[1] < 8 {
				continue
			}
			dir := t.TempDir()
			cfg.Root = dir
			cfg.BodyLimit = limits[0]
			cfg.FileLimit = limits[1]
			if w := scanCall(scanConstruct(t, cfg), raw, media); w.Code != 413 {
				t.Fatal("cap overflow", limits, w.Code)
			}
			scanEmpty(t, dir)
		}
	}
}
func TestHiddenScanRejectsTraversalTypeAndPartialBodies(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "uploads")
	if e := os.Mkdir(root, 0700); e != nil {
		t.Fatal(e)
	}
	sentinel := filepath.Join(parent, "escape.png")
	if e := os.WriteFile(sentinel, []byte("keep"), 0600); e != nil {
		t.Fatal(e)
	}
	calls := 0
	h := scanConstruct(t, ScanUploadConfig{root, 4096, 1024, func() (string, error) { calls++; return "0123456789abcdef", nil }})
	for _, name := range []string{"../escape.png", "x/escape.png", ` + "`" + `x\escape.png` + "`" + `, "/absolute.png", "C:evil.png", "..", "", ".hidden.png", "two..dots.png"} {
		raw, media := scanBody(t, name, "scan", scanSignature, false)
		if w := scanCall(h, raw, media); w.Code != 400 {
			t.Fatal("unsafe raw filename", name, w.Code)
		}
		scanEmpty(t, root)
	}
	for _, tc := range []struct {
		field   string
		payload []byte
		extra   bool
		status  int
	}{
		{"scan", []byte("not a PNG file"), false, 415}, {"scan", scanSignature[:7], false, 415}, {"wrong", scanSignature, false, 400}, {"scan", scanSignature, true, 400},
	} {
		raw, media := scanBody(t, "frame.png", tc.field, tc.payload, tc.extra)
		if w := scanCall(h, raw, media); w.Code != tc.status {
			t.Fatal("invalid multipart", tc, w.Code)
		}
		scanEmpty(t, root)
	}
	raw, media := scanBody(t, "frame.png", "scan", scanSignature, false)
	for _, cut := range []int{0, 20, len(raw) - 30} {
		if w := scanCall(h, raw[:cut], media); w.Code != 400 {
			t.Fatal("truncated", cut, w.Code)
		}
		scanEmpty(t, root)
	}
	if w := scanCall(h, raw, "text/plain"); w.Code != 400 {
		t.Fatal("not multipart")
	}
	// The whole request cap includes bytes after the closing boundary.
	if w := scanCall(h, append(raw, bytes.Repeat([]byte("x"), 4096)...), media); w.Code != 413 {
		t.Fatal("epilogue cap", w.Code)
	}
	scanEmpty(t, root)
	if calls != 0 {
		t.Fatal("source called before validation", calls)
	}
	if kept, e := os.ReadFile(sentinel); e != nil || string(kept) != "keep" {
		t.Fatal("traversal changed sentinel", e)
	}
	for _, source := range []func() (string, error){func() (string, error) { return "", errors.New("private source") }, func() (string, error) { return "../../escape", nil }} {
		cfg := ScanUploadConfig{root, 4096, 1024, source}
		w := scanCall(scanConstruct(t, cfg), raw, media)
		if w.Code != 500 || w.Body.String() != "scan unavailable\n" {
			t.Fatal("source failure", w.Code, w.Body.String())
		}
		scanEmpty(t, root)
	}
	r := httptest.NewRequest("POST", "/cargo-scan", bytes.NewReader(raw))
	r.Header.Set("Content-Type", media)
	ctx, cancel := context.WithCancel(r.Context())
	cancel()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r.WithContext(ctx))
	if w.Code != 400 {
		t.Fatal("cancellation")
	}
	scanEmpty(t, root)
}

type scanBrokenReader struct{ remaining io.Reader }

func (b scanBrokenReader) Read(p []byte) (int, error) {
	n, e := b.remaining.Read(p)
	if e == io.EOF {
		e = io.ErrUnexpectedEOF
	}
	return n, e
}
func TestHiddenScanReaderFailureAndConfiguration(t *testing.T) {
	root := t.TempDir()
	cfg := ScanUploadConfig{root, 4096, 1024, func() (string, error) { return "0123456789abcdef", nil }}
	raw, media := scanBody(t, "frame.png", "scan", append(scanSignature, bytes.Repeat([]byte("x"), 500)...), false)
	h := scanConstruct(t, cfg)
	r := httptest.NewRequest("POST", "/cargo-scan", scanBrokenReader{bytes.NewReader(raw[:len(raw)-50])})
	r.Header.Set("Content-Type", media)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatal("reader failure", w.Code)
	}
	scanEmpty(t, root)
	for _, change := range []func(*ScanUploadConfig){func(c *ScanUploadConfig) { c.NextID = nil }, func(c *ScanUploadConfig) { c.BodyLimit = 0 }, func(c *ScanUploadConfig) { c.FileLimit = 7 }, func(c *ScanUploadConfig) { c.Root = filepath.Join(root, "missing") }} {
		bad := cfg
		change(&bad)
		if h, e := NewScanUpload(bad); h != nil || e == nil {
			t.Fatal("configuration")
		}
	}
	for _, tc := range []struct {
		method, path string
		status       int
	}{{"GET", "/cargo-scan", 405}, {"HEAD", "/cargo-scan", 405}, {"POST", "/cargo-scan/", 404}} {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
		if w.Code != tc.status || (tc.method == "HEAD" && w.Body.Len() != 0) {
			t.Fatal("route", tc, w.Code)
		}
	}
}
func TestHiddenScanExistingSymlinkAndUnavailableRoot(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "uploads")
	if err := os.Mkdir(root, 0700); err != nil {
		t.Fatal(err)
	}
	cfg := ScanUploadConfig{root, 4096, 1024, func() (string, error) { return "0123456789abcdef", nil }}
	h := scanConstruct(t, cfg)
	raw, media := scanBody(t, "frame.png", "scan", scanSignature, false)
	destination := filepath.Join(root, "0123456789abcdef.png")
	if err := os.Symlink(filepath.Join(parent, "absent.png"), destination); err != nil {
		t.Fatal(err)
	}
	w := scanCall(h, raw, media)
	if w.Code != 409 || w.Body.String() != "scan exists\n" {
		t.Fatal("dangling symlink overwritten", w.Code, w.Body.String())
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatal("collision left stage", entries, err)
	}
	if err := os.Remove(destination); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(root); err != nil {
		t.Fatal(err)
	}
	w = scanCall(h, raw, media)
	if w.Code != 500 || w.Body.String() != "scan unavailable\n" {
		t.Fatal("unavailable filesystem", w.Code, w.Body.String())
	}
}
func TestHiddenScanConcurrentPublication(t *testing.T) {
	root := t.TempDir()
	sequence := 0
	h := scanConstruct(t, ScanUploadConfig{root, 4096, 1024, func() (string, error) { sequence++; return fmt.Sprintf("%016x", sequence), nil }})
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		payload := append(append([]byte{}, scanSignature...), byte(i))
		raw, media := scanBody(t, "frame.png", "scan", payload, false)
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := scanCall(h, raw, media)
			if w.Code != 201 {
				t.Error("parallel", w.Code)
				return
			}
			name := strings.TrimSpace(w.Body.String())
			stored, e := os.ReadFile(filepath.Join(root, name))
			if e != nil || !bytes.Equal(stored, payload) {
				t.Error("parallel contents", e)
			}
		}()
	}
	wg.Wait()
	entries, e := os.ReadDir(root)
	if e != nil || len(entries) != 24 || sequence != 24 {
		t.Fatal("publication count", len(entries), e)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".scan-") {
			t.Fatal("partial file")
		}
	}
}
`
