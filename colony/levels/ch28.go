package levels

func init() {
	add(&Level{ID: "28-1-gin", Chapter: 28, Kind: "lab", Title: "Gin", Goal: "انقل سجل الأرصفة لـ Gin من غير ما العقد يتغير، بربط وmiddleware",
		Requires: []string{"sel:gin.New", "sel:binding.JSON"}, Bridge: []int{54}, Setup: seededLabRoom, Race: true,
		Deps:  []string{"github.com/gin-gonic/gin"},
		Files: []LabFile{{Name: "main.go", Starter: ch281Starter}},
		Lesson: `# نفس سجل الأرصفة، على Gin

## قبل ما تكتب
- إيه الإعدادات الافتراضية في Gin اللي ممكن تغيّر عقد 26-2 من غير ما تاخد بالك؟
- فين بيخلص شغل الإطار، وفين بيبدأ المخزن والقفل بتوعك؟

## ليه إطار أصلاً؟
في 26-2 حلّلت المسار والطرق بإيدك. Gin بيعمل الراوتنج وباراميترات المسار وربط JSON وسلسلة middleware. العقد اللي العميل شايفه لازم يفضل زي ما هو بالظبط؛ الإطار تفصيلة جوّه.

استورد ~github.com/gin-gonic/gin~. ~gin.New()~ بيرجع ~*gin.Engine~ فاضي من غير أي middleware. الـ engine بينفّذ ~http.Handler~، يعني تختبره بـ ~httptest~ زي أي handler من فصل 26.
~~~go
r := gin.New()
r.GET("/probes/:name", func(c *gin.Context) {
 c.String(http.StatusOK, "probe "+c.Param("name")+"\n")
})
w := httptest.NewRecorder()
r.ServeHTTP(w, httptest.NewRequest("GET", "/probes/vega", nil))
~~~
~r.GET~ و ~r.POST~ و ~r.PUT~ و ~r.DELETE~ بيسجّلوا مسار لطريقة واحدة. ~:name~ جزء متغير من المسار، و ~c.Param("name")~ بيرجعه نص. ~c.Request~ هو نفس ~*http.Request~ اللي تعرفه.

## الرد من الـ Context
~c.String~ بيكتب نص بنوع ~text/plain; charset=utf-8~. ~c.JSON~ بيرمّز أي قيمة بـ ~json.Marshal~ ويحط ~application/json; charset=utf-8~. ~c.Header~ بيحط هيدر رد، و ~c.Status~ بيحدد حالة من غير جسم.
~~~go
c.Header("Location", "/probes/7")
c.JSON(http.StatusCreated, probe) // 201 والجسم JSON
c.Status(http.StatusNoContent)    // 204 من غير جسم
~~~
~gin.Default()~ بيضيف logger و recovery من عنده وبيكتب على الـ stdout؛ هنا ابدأ بـ ~gin.New()~ عشان السلسلة كلها تبقى بتاعتك. ~gin.SetMode(gin.TestMode)~ بيسكّت رسايل الـ debug؛ ده إعداد عام للباكدج كله، والتستات المخفية بتعمله مرة واحدة.

## افتراضيات بتكسر العقد
~RedirectTrailingSlash~ في ~gin.New()~ قيمته ~true~: طلب ~/berths/~ بيترد بتحويل 301 لـ ~/berths~ بدل 404. ~HandleMethodNotAllowed~ قيمته ~false~: طريقة غلط على مسار موجود بترجع 404 بدل 405. دول حقول على الـ engine، غيّرهم قبل ما تسجّل المسارات.
~~~go
r.RedirectTrailingSlash = false
r.HandleMethodNotAllowed = true
r.NoRoute(func(c *gin.Context) { c.String(404, "lost\n") })
r.NoMethod(func(c *gin.Context) { c.String(405, "wrong method\n") })
~~~
~NoRoute~ بيشتغل لما مفيش مسار خالص، و ~NoMethod~ لما المسار موجود بطريقة تانية. قبل ~NoMethod~ بيحط Gin هيدر ~Allow~ بنفسه بطرق المسار المسجلة، بترتيب تسجيلها.

## middleware بشكل Gin
الـ middleware في Gin دالة ~gin.HandlerFunc~ عادية بتنادي ~c.Next()~ عشان باقي السلسلة يشتغل. ~r.Use~ بيضيفها لكل المسارات، وكمان لـ ~NoRoute~ و ~NoMethod~. ~c.Abort()~ بيوقف السلسلة لو الـ middleware رد بنفسه.
~~~go
stamp := func(c *gin.Context) {
 c.Header("X-Deck", "7") // قبل أي كتابة
 c.Next()                // الـ handler بيشتغل هنا
}
r.Use(stamp)
~~~
حط هيدرات الرد قبل ~c.Next()~؛ بعد ما الـ handler يكتب الجسم الهيدرات اتبعتت خلاص.

## الربط (binding) والتحقق
استورد ~github.com/gin-gonic/gin/binding~. ~binding.JSON.BindBody(raw, &in)~ بيفك بايتات JSON في struct، وبعدين بيشغّل قواعد التحقق المكتوبة في tag اسمه ~binding~. ~required~ معناها القيمة مش القيمة الصفرية، و ~min~ و ~max~ حدود الأرقام شاملة.
~~~go
type probeInput struct {
 Name  string ` + "`" + `json:"name" binding:"required"` + "`" + `
 Depth int    ` + "`" + `json:"depth" binding:"required,min=1,max=9"` + "`" + `
}
var in probeInput
err := binding.JSON.BindBody(raw, &in) // نوع غلط أو قاعدة مكسورة
~~~
الربط بيقرا أول قيمة JSON بس، فـ ~{"name":"a","depth":1} {}~ بيعدّي. ~json.Valid(raw)~ بترجع ~true~ لو البايتات قيمة JSON واحدة صحيحة بالظبط، فافحصها قبل الربط. ~required~ مش بتعمل trim؛ الاسم ~" "~ بيعدّي وانت اللي ترفضه.
~~~go
json.Valid([]byte("{} {}"))  // false: قيمتين
json.Valid([]byte(" [1] \n")) // true: المسافات حوالين القيمة مسموحة
~~~
~c.ShouldBindJSON~ بيقرا الجسم كله من غير حد. اقرا بـ ~io.LimitReader~ زي 26-2 الأول، وبعدين اربط البايتات.

## الـ Context مش بيعيش بعد الرد
Gin بيعيد استخدام نفس ~*gin.Context~ لطلب تاني بعد ما الـ handler يرجع. متخزّنش ~c~، ومتستخدموش من goroutine بعد الرد؛ انسخ القيم اللي محتاجها. المخزن هو نفس فكرة 26-2: عداد وخريطة تحت ~sync.Mutex~ واحد.

## النجمة التالتة: الربط من الـ query
~c.ShouldBindQuery(&f)~ بيملا struct من الـ query باستخدام tag اسمه ~form~. استخدم pointer عشان تفرّق بين «مش موجود» و «صفر»؛ ~omitempty~ بتتخطى باقي القواعد لما الـ pointer يبقى ~nil~.
~~~go
type depthFilter struct {
 Max *int ` + "`" + `form:"max_depth" binding:"omitempty,min=1,max=9"` + "`" + `
}
var f depthFilter
if err := c.ShouldBindQuery(&f); err != nil { c.String(400, "bad filter\n"); return }
~~~

## أخطاء شائعة
~~~go
r := gin.Default()                   // logger وrecovery مش بتوعك
r.GET("/berths/", list)              // الشرطة الأخيرة بقت جزء من المسار
if c.ShouldBindJSON(&in) == nil { }  // من غير حد للجسم، وبيقبل قيمة تانية بعد الأولى
go func() { use(c.Param("id")) }()   // c ممكن يكون اتعاد استخدامه
~~~

## المهمة
نفّذ ~NewBerthEngine() *gin.Engine~ يرجع engine جديد بمخزن مستقل وآمن بالتزامن، بعقد 26-2 على نفس ~Berth~. ~GET /berths~ يرجع 200 ومصفوفة مرتبة تصاعدياً بالـ ID، والفاضية ~[]~. ~POST /berths~ يقبل كائن JSON واحد: ~label~ بعد trim مش فاضي، و ~slots~ رقم صحيح من 1 لـ 64 شاملة؛ ~id~ والحقول المجهولة بتتجاهل. الـ IDs تبدأ من 1 وتزيد ومتتعادش، والطلب المرفوض مبياخدش رقم. النجاح 201 والعنصر و ~Location: /berths/ID~.

~GET /berths/ID~ يرجع 200 والعنصر. ~PUT~ لنفس المسار يستبدل الحقلين بنفس التحقق، ويحتفظ بالـ ID، ويرجع 200 والعنصر. ~DELETE~ يرجع 204 بجسم فاضي. العنصر الغائب 404 و ~unknown berth\n~. JSON بايظ أو نوع غلط أو قيمة برا الحدود أو ID مش رقم موجب مكتوب بالشكل القياسي: 400 و ~invalid berth\n~. في ~PUT~ افحص الـ ID، وبعده الجسم، وبعدين دوّر على العنصر. جسم الإنشاء والاستبدال حده 4096 بايت شاملة المسافات؛ الأكبر 413 و ~too large\n~ قبل أي تحليل.

أي مسار تاني، منهم ~/berths/~ و ~/Berths~ و ~/berths/2/x~، يرجع 404 و ~no route\n~. طريقة مش مدعومة على مسار موجود ترجع 405 و ~use route methods\n~، و ~Allow~ فيه طرق المسار: ~GET, POST~ للمجموعة و ~DELETE, GET, PUT~ للعنصر، والترتيب مش مهم. كل رد، حتى 404 و 405 و 413، فيه ~Cache-Control: no-store~ من middleware واحد. الردود الناجحة اللي ليها جسم نوعها ~application/json~ والأخطاء ~text/plain~، والـ charset مسموح. الـ query بتتجاهل، والفشل ميغيّرش الحالة. مفيش شبكة ولا ~Report~، و ~main~ فاضية.

النجمة التالتة: ~GET /berths?min_slots=N~ يرجع بس الأرصفة اللي ~slots~ بتاعتها N أو أكتر بنفس الترتيب؛ N لازم رقم صحيح من 1 لـ 64، وإلا 400 و ~invalid berth\n~. من غير ~min_slots~ القائمة كاملة.`,
		Hints: []string{"ابدأ بـ gin.New() وصفّر RedirectTrailingSlash وشغّل HandleMethodNotAllowed، وسجّل NoRoute وNoMethod وmiddleware الـ Cache-Control بـ Use.", "اقرا الجسم بـ io.LimitReader لحد 4097، وبعدين json.Valid، وبعدين binding.JSON.BindBody على struct بقواعد required وmin وmax، وبعدين trim للاسم.", "خلّي المخزن struct فيه sync.Mutex وعداد وخريطة؛ في PUT افحص c.Param والجسم قبل البحث، وللنجمة استخدم ShouldBindQuery بـ *int."},
		En: &LevelText{Title: "Gin", Goal: "Move the berth registry to Gin without changing its contract, using binding and middleware", Files: []LabFile{{Name: "main.go", Starter: ch281Starter}},
			Lesson: `# The Same Berth Registry, on Gin

## Before you write
- Which Gin defaults could silently change the 26-2 contract?
- Where does the framework's work end, and where do your store and lock begin?

## Why a framework at all?
In 26-2 you parsed paths and methods by hand. Gin does routing, path parameters, JSON binding, and middleware chains. The contract the client sees must stay exactly the same; the framework is an inside detail.

Import ~github.com/gin-gonic/gin~. ~gin.New()~ returns an empty ~*gin.Engine~ with no middleware. The engine implements ~http.Handler~, so you test it with ~httptest~ like any handler from chapter 26.
~~~go
r := gin.New()
r.GET("/probes/:name", func(c *gin.Context) {
 c.String(http.StatusOK, "probe "+c.Param("name")+"\n")
})
w := httptest.NewRecorder()
r.ServeHTTP(w, httptest.NewRequest("GET", "/probes/vega", nil))
~~~
~r.GET~, ~r.POST~, ~r.PUT~, and ~r.DELETE~ register a path for one method. ~:name~ is a variable path segment, and ~c.Param("name")~ returns it as text. ~c.Request~ is the same ~*http.Request~ you already know.

## Replying through the Context
~c.String~ writes text as ~text/plain; charset=utf-8~. ~c.JSON~ encodes any value with ~json.Marshal~ and sets ~application/json; charset=utf-8~. ~c.Header~ sets a response header, and ~c.Status~ sets a status with no body.
~~~go
c.Header("Location", "/probes/7")
c.JSON(http.StatusCreated, probe) // 201 with a JSON body
c.Status(http.StatusNoContent)    // 204 with no body
~~~
~gin.Default()~ adds its own logger and recovery and writes to stdout; here start from ~gin.New()~ so the whole chain is yours. ~gin.SetMode(gin.TestMode)~ silences debug messages; it is a package-wide setting, and the hidden tests set it once.

## Defaults that break the contract
~RedirectTrailingSlash~ is ~true~ in ~gin.New()~: a request for ~/berths/~ gets a 301 redirect to ~/berths~ instead of 404. ~HandleMethodNotAllowed~ is ~false~: a wrong method on an existing path returns 404 instead of 405. These are engine fields; change them before registering routes.
~~~go
r.RedirectTrailingSlash = false
r.HandleMethodNotAllowed = true
r.NoRoute(func(c *gin.Context) { c.String(404, "lost\n") })
r.NoMethod(func(c *gin.Context) { c.String(405, "wrong method\n") })
~~~
~NoRoute~ runs when no path matches at all, and ~NoMethod~ when the path exists for other methods. Before ~NoMethod~, Gin sets the ~Allow~ header itself from the path's registered methods, in registration order.

## Middleware, Gin style
Gin middleware is an ordinary ~gin.HandlerFunc~ that calls ~c.Next()~ to run the rest of the chain. ~r.Use~ adds it to every route, and also to ~NoRoute~ and ~NoMethod~. ~c.Abort()~ stops the chain when the middleware has replied itself.
~~~go
stamp := func(c *gin.Context) {
 c.Header("X-Deck", "7") // Before any write
 c.Next()                // The handler runs here
}
r.Use(stamp)
~~~
Set response headers before ~c.Next()~; once the handler writes the body, headers are already sent.

## Binding and validation
Import ~github.com/gin-gonic/gin/binding~. ~binding.JSON.BindBody(raw, &in)~ decodes JSON bytes into a struct, then runs validation rules written in a tag named ~binding~. ~required~ means the value is not its zero value, and ~min~ and ~max~ are inclusive number bounds.
~~~go
type probeInput struct {
 Name  string ` + "`" + `json:"name" binding:"required"` + "`" + `
 Depth int    ` + "`" + `json:"depth" binding:"required,min=1,max=9"` + "`" + `
}
var in probeInput
err := binding.JSON.BindBody(raw, &in) // Wrong type or broken rule
~~~
Binding reads only the first JSON value, so ~{"name":"a","depth":1} {}~ passes. ~json.Valid(raw)~ returns ~true~ when the bytes are exactly one well-formed JSON value, so check it before binding. ~required~ does not trim; the name ~" "~ passes and you must reject it.
~~~go
json.Valid([]byte("{} {}"))  // false: two values
json.Valid([]byte(" [1] \n")) // true: whitespace around the value is fine
~~~
~c.ShouldBindJSON~ reads the whole body with no limit. Read through ~io.LimitReader~ first as in 26-2, then bind the bytes.

## The Context does not outlive the reply
Gin reuses the same ~*gin.Context~ for another request once the handler returns. Never store ~c~ or use it from a goroutine after the reply; copy the values you need. The store is the same idea as 26-2: a counter and a map under one ~sync.Mutex~.

## Third star: binding the query
~c.ShouldBindQuery(&f)~ fills a struct from the query using a tag named ~form~. Use a pointer to tell "absent" from "zero"; ~omitempty~ skips the other rules when the pointer is ~nil~.
~~~go
type depthFilter struct {
 Max *int ` + "`" + `form:"max_depth" binding:"omitempty,min=1,max=9"` + "`" + `
}
var f depthFilter
if err := c.ShouldBindQuery(&f); err != nil { c.String(400, "bad filter\n"); return }
~~~

## Common mistakes
~~~go
r := gin.Default()                   // Logger and recovery you did not choose
r.GET("/berths/", list)              // The trailing slash became part of the path
if c.ShouldBindJSON(&in) == nil { }  // No body cap, and a second value after the first passes
go func() { use(c.Param("id")) }()   // c may already be reused
~~~

## Your mission
Implement ~NewBerthEngine() *gin.Engine~ returning a new engine with an independent concurrent-safe store and the 26-2 contract over the same ~Berth~. ~GET /berths~ returns 200 and an array in ascending ID order, ~[]~ when empty. ~POST /berths~ accepts one JSON object: a ~label~ that is nonblank after trimming and an integer ~slots~ in inclusive 1..64; input ~id~ and unknown fields are ignored. IDs start at 1, increase, and are never reused; rejected requests consume none. Success is 201 with the stored object and ~Location: /berths/ID~.

~GET /berths/ID~ returns 200 and the object. ~PUT~ on that path replaces both fields with the same validation, keeps the ID, and returns 200 and the object. ~DELETE~ returns 204 with an empty body. A missing item is 404 and ~unknown berth\n~. Invalid JSON, a wrong type, an out-of-range value, or an ID that is not a canonical positive integer is 400 and ~invalid berth\n~. For ~PUT~, check the ID, then the body, then look up the item. Create and replace bodies are capped at 4096 bytes including whitespace; larger ones are 413 and ~too large\n~ before any parsing.

Every other path, including ~/berths/~, ~/Berths~, and ~/berths/2/x~, returns 404 and ~no route\n~. An unsupported method on an existing path returns 405, ~use route methods\n~, and an ~Allow~ listing the path's methods: ~GET, POST~ for the collection and ~DELETE, GET, PUT~ for an item, in any order. Every reply, including 404, 405, and 413, carries ~Cache-Control: no-store~ from one middleware. Successful replies with a body are ~application/json~ and errors are ~text/plain~; a charset is allowed. The query is ignored, and failures never change state. No network and no ~Report~; leave ~main~ empty.

Third star: ~GET /berths?min_slots=N~ returns only berths whose ~slots~ is at least N, in the same order; N must be an integer from 1 to 64, otherwise 400 and ~invalid berth\n~. Without ~min_slots~ the list is complete.`,
			Hints: []string{"Start with gin.New(), turn off RedirectTrailingSlash, turn on HandleMethodNotAllowed, and register NoRoute, NoMethod, and the Cache-Control middleware with Use.", "Read the body through io.LimitReader up to 4097, then json.Valid, then binding.JSON.BindBody into a struct with required, min, and max rules, then trim the label.", "Make the store a struct with a sync.Mutex, a counter, and a map; in PUT check c.Param and the body before the lookup, and for the star use ShouldBindQuery with *int."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch281Hidden, HiddenFiles: map[string]string{"berth_contract_test.go": ch28BerthContract}},
	})
}

const ch281Starter = `package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Berth struct {
	ID    int    ` + "`" + `json:"id"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
	Slots int    ` + "`" + `json:"slots"` + "`" + `
}

func NewBerthEngine() *gin.Engine {
	r := gin.New()
	r.NoRoute(func(c *gin.Context) { c.String(http.StatusNotImplemented, "TODO\n") })
	return r
}
func main() {}
`

// ch28BerthContract is the one contract test shared by 28-1 (Gin) and 28-2
// (Fiber): each lab only supplies a contractCaller adapter.
const ch28BerthContract = `package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mime"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
)

type contractReply struct {
	Code   int
	Header http.Header
	Body   string
}

type contractCaller func(method, path, body string) contractReply

func contractMedia(r contractReply) string {
	media, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	return media
}

func contractShort(body string) string {
	if len(body) > 120 {
		return body[:120] + "..."
	}
	return body
}

func contractText(t *testing.T, label string, r contractReply, status int, body string) {
	t.Helper()
	if r.Code != status || r.Body != body || contractMedia(r) != "text/plain" || r.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("%s: got %d %q %v, want %d %q text/plain no-store", label, r.Code, contractShort(r.Body), r.Header, status, body)
	}
}

func contractJSONOK(r contractReply, status int, into any) bool {
	return r.Code == status && contractMedia(r) == "application/json" && r.Header.Get("Cache-Control") == "no-store" && json.Unmarshal([]byte(r.Body), into) == nil
}

func contractJSON(t *testing.T, label string, r contractReply, status int, into any) {
	t.Helper()
	if !contractJSONOK(r, status, into) {
		t.Fatalf("%s: got %d %q %v, want %d application/json no-store", label, r.Code, contractShort(r.Body), r.Header, status)
	}
}

func contractAllow(h http.Header) []string {
	var out []string
	for _, v := range h.Values("Allow") {
		for _, m := range strings.Split(v, ",") {
			if m = strings.TrimSpace(m); m != "" && m != "HEAD" {
				out = append(out, m)
			}
		}
	}
	slices.Sort(out)
	return out
}

func contractBody(label string, slots int) string {
	raw, _ := json.Marshal(map[string]any{"label": label, "slots": slots})
	return string(raw)
}

// Boundaries: empty list, strict JSON, the 4096-byte cap, slot limits, canonical IDs,
// stable IDs after deletes, validation before lookup, framework defaults, isolation.
func contractCRUD(t *testing.T, fresh func() contractCaller, seed int64) {
	call := fresh()
	rng := rand.New(rand.NewSource(seed))
	var list []Berth
	contractJSON(t, "empty list", call("GET", "/berths", ""), 200, &list)
	if list == nil || len(list) != 0 {
		t.Fatal("empty list must be [] not null")
	}
	for _, body := range []string{"", "null", "[]", "{}", "7", ` + "`" + `{"label":7,"slots":2}` + "`" + `, ` + "`" + `{"label":"x","slots":0}` + "`" + `,
		` + "`" + `{"label":"x","slots":65}` + "`" + `, ` + "`" + `{"label":" ","slots":1}` + "`" + `, ` + "`" + `{"label":"x","slots":1} {}` + "`" + `, ` + "`" + `{"label":"x","slots":"2"}` + "`" + `,
		` + "`" + `{"label":"x","slots":1.5}` + "`" + `, ` + "`" + `{"label":"x"}` + "`" + `, ` + "`" + `{"slots":3}` + "`" + `, ` + "`" + `{"label":"x","slots":1` + "`" + `} {
		contractText(t, fmt.Sprintf("POST %q", body), call("POST", "/berths", body), 400, "invalid berth\n")
	}
	contractText(t, "4097 spaces", call("POST", "/berths", strings.Repeat(" ", 4097)), 413, "too large\n")
	padded := contractBody("edge", 1)
	contractText(t, "4097 bytes", call("POST", "/berths", padded+strings.Repeat(" ", 4097-len(padded))), 413, "too large\n")
	var first Berth
	r := call("POST", "/berths", padded+strings.Repeat(" ", 4096-len(padded)))
	contractJSON(t, "4096 bytes", r, 201, &first)
	if first != (Berth{1, "edge", 1}) || r.Header.Get("Location") != "/berths/1" {
		t.Fatal("rejected requests consumed an ID or wrong reply", first, r.Header)
	}
	r = call("DELETE", "/berths/1", "")
	if r.Code != 204 || r.Body != "" || r.Header.Get("Cache-Control") != "no-store" {
		t.Fatal("delete", r.Code, r.Body, r.Header)
	}
	want := []Berth{}
	n := 12 + rng.Intn(8)
	for i := 2; i < 2+n; i++ {
		label := fmt.Sprintf("bay-%d-🚀", rng.Intn(100000))
		slots := 1 + rng.Intn(64)
		var b Berth
		r := call("POST", "/berths", fmt.Sprintf(` + "`" + `{"label":" %s ","slots":%d,"id":99,"crew":"x"}` + "`" + `, label, slots))
		contractJSON(t, "create", r, 201, &b)
		if b != (Berth{i, label, slots}) || r.Header.Get("Location") != fmt.Sprintf("/berths/%d", i) {
			t.Fatal("identity", b, r.Header)
		}
		want = append(want, b)
	}
	for i, b := range want {
		path := fmt.Sprintf("/berths/%d", b.ID)
		var got Berth
		contractJSON(t, "read", call("GET", path+"?crew=1", ""), 200, &got)
		if got != b {
			t.Fatal("read", got, b)
		}
		b.Label += "-updated"
		b.Slots = 1 + rng.Intn(64)
		contractJSON(t, "replace", call("PUT", path, contractBody(b.Label, b.Slots)), 200, &got)
		if got != b {
			t.Fatal("replacement", got, b)
		}
		want[i] = b
		contractText(t, "invalid replace", call("PUT", path, ` + "`" + `{"label":"bad","slots":-1}` + "`" + `), 400, "invalid berth\n")
		contractText(t, "oversized replace", call("PUT", path, strings.Repeat(" ", 5000)), 413, "too large\n")
	}
	contractText(t, "validate before lookup", call("PUT", "/berths/1", "[]"), 400, "invalid berth\n")
	for _, tc := range []struct {
		method, path, body string
		status             int
		text               string
	}{
		{"GET", "/berths/0", "", 400, "invalid berth\n"}, {"GET", "/berths/01", "", 400, "invalid berth\n"}, {"GET", "/berths/-1", "", 400, "invalid berth\n"},
		{"GET", "/berths/+2", "", 400, "invalid berth\n"}, {"GET", "/berths/99999999999999999999999", "", 400, "invalid berth\n"}, {"DELETE", "/berths/no", "", 400, "invalid berth\n"},
		{"PUT", "/berths/x", contractBody("a", 1), 400, "invalid berth\n"},
		{"GET", "/berths/1", "", 404, "unknown berth\n"}, {"DELETE", "/berths/1", "", 404, "unknown berth\n"}, {"PUT", "/berths/1", contractBody("a", 1), 404, "unknown berth\n"},
		{"GET", "/berths/", "", 404, "no route\n"}, {"GET", "/berths/2/x", "", 404, "no route\n"}, {"GET", "/berths/2/", "", 404, "no route\n"},
		{"POST", "/missing", contractBody("a", 1), 404, "no route\n"}, {"GET", "/Berths", "", 404, "no route\n"}, {"GET", "/", "", 404, "no route\n"},
	} {
		contractText(t, tc.method+" "+tc.path, call(tc.method, tc.path, tc.body), tc.status, tc.text)
	}
	for _, tc := range []struct {
		method, path string
		allow        []string
	}{
		{"DELETE", "/berths", []string{"GET", "POST"}}, {"PATCH", "/berths", []string{"GET", "POST"}}, {"PUT", "/berths", []string{"GET", "POST"}},
		{"POST", "/berths/2", []string{"DELETE", "GET", "PUT"}}, {"PATCH", "/berths/2", []string{"DELETE", "GET", "PUT"}},
	} {
		r := call(tc.method, tc.path, contractBody("a", 1))
		contractText(t, tc.method+" "+tc.path, r, 405, "use route methods\n")
		if got := contractAllow(r.Header); !reflect.DeepEqual(got, tc.allow) {
			t.Fatalf("%s %s Allow = %v, want %v", tc.method, tc.path, got, tc.allow)
		}
	}
	var got []Berth
	contractJSON(t, "final list", call("GET", "/berths?crew=1", ""), 200, &got)
	if !reflect.DeepEqual(got, want) {
		t.Fatal("list ordering or failed requests changed state", got, want)
	}
	var other Berth
	contractJSON(t, "isolated", fresh()("POST", "/berths", contractBody("isolated", 1)), 201, &other)
	if other.ID != 1 {
		t.Fatal("two instances share a store")
	}
}

func contractParallel(t *testing.T, fresh func() contractCaller) {
	call := fresh()
	var wg sync.WaitGroup
	start := make(chan struct{})
	ids := make(chan int, 64)
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			var b Berth
			if !contractJSONOK(call("POST", "/berths", contractBody(fmt.Sprintf("dock-%d", i), 1+i%64)), 201, &b) {
				t.Error("parallel create failed")
			}
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
			var b Berth
			if !contractJSONOK(call("GET", path, ""), 200, &b) || !contractJSONOK(call("PUT", path, contractBody("replaced", 2)), 200, &b) {
				t.Error("parallel read/replace")
			}
			if call("DELETE", path, "").Code != 204 {
				t.Error("parallel delete")
			}
			_ = call("GET", "/berths", "")
		}(id)
	}
	wg.Wait()
	var left []Berth
	if !contractJSONOK(call("GET", "/berths", ""), 200, &left) || len(left) != 0 {
		t.Fatal("surviving entries", left)
	}
}
`

const ch281Hidden = `package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

func ginCaller(t *testing.T) contractCaller {
	e := NewBerthEngine()
	if e == nil {
		t.Fatal("NewBerthEngine returned nil")
	}
	var h http.Handler = e
	return func(method, path, body string) contractReply {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
		return contractReply{w.Code, w.Header(), w.Body.String()}
	}
}

func TestHiddenGinBerthContract(t *testing.T) {
	contractCRUD(t, func() contractCaller { return ginCaller(t) }, 281)
}

func TestHiddenGinBerthsParallel(t *testing.T) {
	contractParallel(t, func() contractCaller { return ginCaller(t) })
}

// Boundaries: absent filter, both limits, zero, out of range, non-integers, empty result.
func TestHiddenBonusGinSlotFilter(t *testing.T) {
	call := ginCaller(t)
	rng := rand.New(rand.NewSource(2812))
	var all []Berth
	for i := 0; i < 30; i++ {
		var b Berth
		contractJSON(t, "create", call("POST", "/berths", contractBody(fmt.Sprintf("f-%d", i), 1+rng.Intn(63))), 201, &b)
		all = append(all, b)
	}
	for _, n := range []int{1, 64, 1 + rng.Intn(64), 1 + rng.Intn(64), 32} {
		want := []Berth{}
		for _, b := range all {
			if b.Slots >= n {
				want = append(want, b)
			}
		}
		var got []Berth
		r := call("GET", fmt.Sprintf("/berths?min_slots=%d&crew=1", n), "")
		contractJSON(t, "filter", r, 200, &got)
		if got == nil || !reflect.DeepEqual(got, want) {
			raw, _ := json.Marshal(want)
			t.Fatalf("min_slots=%d: got %s want %s", n, r.Body, raw)
		}
	}
	for _, bad := range []string{"0", "65", "-3", "abc", "1.5", "9999999999999999999999"} {
		contractText(t, "min_slots="+bad, call("GET", "/berths?min_slots="+bad, ""), 400, "invalid berth\n")
	}
	var got []Berth
	contractJSON(t, "no filter", call("GET", "/berths", ""), 200, &got)
	if !reflect.DeepEqual(got, all) {
		t.Fatal("unfiltered list changed")
	}
}
`

func init() {
	add(&Level{ID: "28-2-fiber", Chapter: 28, Kind: "lab", Title: "Fiber", Goal: "شغّل نفس عقد الأرصفة على Fiber واختبره بـ app.Test، واكتب مقارنة قصيرة بين الإطارين",
		Requires: []string{"sel:fiber.New", "sel:fiber.Config"}, Bridge: []int{55}, Setup: seededLabRoom, Race: true,
		Deps:  []string{"github.com/gofiber/fiber/v2"},
		Files: []LabFile{{Name: "main.go", Starter: ch282Starter}, {Name: "compare.go", Starter: ch282CompareAr}},
		Lesson: `# السجل نفسه على Fiber

## قبل ما تكتب
- Fiber مش مبني على ~net/http~؛ إيه اللي بيتغير في الاختبار وفي عمر البيانات؟
- الخطأ اللي بيرجع من handler بيروح فين، ومين بيكتب رده؟

## Fiber فوق fasthttp
استورد ~github.com/gofiber/fiber/v2~؛ اسم الباكدج ~fiber~. Fiber مبني على fasthttp مش ~net/http~، فـ ~*fiber.App~ مش ~http.Handler~ ومفيهوش ~ServeHTTP~. الـ handler شكله ~func(c *fiber.Ctx) error~.
~~~go
app := fiber.New()
app.Get("/probes/:name", func(c *fiber.Ctx) error {
 return c.Status(fiber.StatusOK).SendString("probe " + c.Params("name") + "\n")
})
~~~
~app.Get~ و ~app.Post~ و ~app.Put~ و ~app.Delete~ بيسجّلوا مسار لطريقة. ~c.Params("name")~ بيقرا جزء ~:name~، و ~c.Body()~ بيرجع بايتات الجسم كلها، و ~c.Method()~ الطريقة.

## الرد
~c.Status~ بيحدد الحالة ويرجع نفس الـ ctx عشان تكمّل عليه. ~c.SendString~ بيكتب نص؛ نوع المحتوى الافتراضي ~text/plain; charset=utf-8~. ~c.JSON~ بيرمّز القيمة ويحط ~application/json~. ~c.Set~ بيحط هيدر رد.
~~~go
c.Set("Location", "/probes/7")
return c.Status(fiber.StatusCreated).JSON(probe)
~~~
~c.SendStatus(204)~ بيحط الحالة، ولو الجسم فاضي بيكتب اسم الحالة نص؛ السيرفر بيشيل الجسم ده لوحده مع 204. ~c.Status(204)~ ووراها ~return nil~ أوضح.

## الاختبار بـ app.Test
~app.Test(req, -1)~ بياخد ~*http.Request~ عادي من ~httptest.NewRequest~ ويشغّله في الذاكرة من غير بورت، ويرجع ~*http.Response~ وخطأ. المهلة الافتراضية ثانية واحدة؛ ~-1~ بتلغيها عشان الطلبات المتوازية تحت ~-race~ متفشلش من البطء. اقفل ~resp.Body~ بعد القراية.
~~~go
resp, err := app.Test(httptest.NewRequest("GET", "/probes/vega", nil), -1)
if err != nil { t.Fatal(err) }
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
~~~

## الإعداد اللي بيحمي العقد
~fiber.New~ بياخد ~fiber.Config~ اختياري. ~StrictRouting~ افتراضياً ~false~، يعني ~/berths/~ بيطابق ~/berths~. ~CaseSensitive~ افتراضياً ~false~، يعني ~/Berths~ بيطابق. الاتنين لازم ~true~ هنا. ~BodyLimit~ الافتراضي 4 MiB؛ ده سقف السيرفر مش حد العقد، فافحص ~len(c.Body())~ بنفسك.
~~~go
app := fiber.New(fiber.Config{StrictRouting: true, CaseSensitive: true, ErrorHandler: onFault})
~~~

## الأخطاء بترجع، مش بتتكتب
الـ handler بيرجع ~error~، و Fiber بيبعت أي خطأ لـ ~ErrorHandler~ واحد. المسار المجهول بيوصله ~*fiber.Error~ بـ ~Code~ قيمته 404. الطريقة الغلط على مسار موجود بتوصله بـ 405، ومعاها هيدر ~Allow~ حطه Fiber بالطرق اللي ليها مسار. افحص النوع بـ ~errors.As~.
~~~go
func onFault(c *fiber.Ctx, err error) error {
 var fe *fiber.Error
 if errors.As(err, &fe) && fe.Code == fiber.StatusNotFound {
  return c.Status(404).SendString("lost\n")
 }
 return c.Status(500).SendString("fault\n")
}
~~~
~fiber.NewError(code, msg)~ بيبني نفس النوع لو عايز ترجّعه من handler. ~app.Get~ بيسجّل ~HEAD~ كمان لنفس المسار، فـ ~Allow~ ممكن يبان فيه ~HEAD~ وده مسموح.

## middleware بشكل Fiber
middleware الـ ~net/http~ اللي بتلف ~http.Handler~ مينفعش تتحط هنا. في Fiber هي handler عادي بينادي ~c.Next()~ ويرجّع الخطأ اللي رجع منه. ~app.Use~ بيضيفها قبل كل المسارات، حتى الطلب اللي هيخلص 404.
~~~go
app.Use(func(c *fiber.Ctx) error {
 c.Set("X-Deck", "7")
 return c.Next() // رجّع خطأ الباقي عشان يوصل لـ ErrorHandler
})
~~~

## عمر البيانات
fasthttp بيعيد استخدام بافرات الطلب بعد ما الـ handler يرجع. النص اللي راجع من ~c.Params~ أو ~c.Get~ صالح بس جوه الـ handler؛ لو هتخزنه انسخه بـ ~strings.Clone~. ~json.Unmarshal~ بتنسخ النصوص لوحدها، فالحقول اللي فكّيتها من ~c.Body()~ آمنة.
~~~go
name := strings.Clone(c.Params("name")) // نسخة بتعيش بعد الرد
~~~
~c.BodyParser~ بيختار الـ decoder من ~Content-Type~ وبيرفض لو الهيدر مش موجود. العقد مبيشترطش content type، فافحص ~c.Body()~ بـ ~json.Unmarshal~ زي 26-2؛ هي أصلاً بترفض قيمة تانية بعد الأولى.

## ملف المقارنة
في ~compare.go~ املا ~Comparison~ بأربع ~Contrast~، واحدة لكل موضوع: ~handler~ (شكل الـ handler والواجهة)، و ~errors~ (404 و 405 بيتكتبوا فين)، و ~testing~ (إزاي بتختبر من غير بورت)، و ~lifetime~ (إيه اللي ممنوع يعيش بعد الرد). كل خانة جملة قصيرة عن الإطار ده بالذات.
~~~go
{Topic: "testing", Gin: "...", Fiber: "..."}
~~~

## أخطاء شائعة
~~~go
app := fiber.New()                        // /berths/ و /Berths بيعدّوا
if err := c.BodyParser(&in); err != nil {} // من غير Content-Type بيفشل
seen[c.Params("id")] = true                // البافر بيتعاد استخدامه
resp, _ := app.Test(req)                   // مهلة ثانية، والخطأ اتجاهل
~~~

## المهمة
نفّذ ~NewBerthApp() *fiber.App~ يرجع app جديد بمخزن مستقل وآمن بالتزامن، بنفس عقد 28-1 بالظبط على نفس ~Berth~. ~GET /berths~ يرجع 200 ومصفوفة مرتبة تصاعدياً بالـ ID، والفاضية ~[]~. ~POST /berths~ يقبل كائن JSON واحد: ~label~ بعد trim مش فاضي، و ~slots~ رقم صحيح من 1 لـ 64 شاملة؛ ~id~ والحقول المجهولة بتتجاهل. الـ IDs تبدأ من 1 وتزيد ومتتعادش، والطلب المرفوض مبياخدش رقم. النجاح 201 والعنصر و ~Location: /berths/ID~.

~GET /berths/ID~ يرجع 200 والعنصر. ~PUT~ يستبدل الحقلين بنفس التحقق ويحتفظ بالـ ID ويرجع 200. ~DELETE~ يرجع 204 بجسم فاضي. الغائب 404 و ~unknown berth\n~. JSON بايظ أو نوع غلط أو قيمة برا الحدود أو ID مش رقم موجب قياسي: 400 و ~invalid berth\n~. في ~PUT~ افحص الـ ID، وبعده الجسم، وبعدين دوّر. الجسم حده 4096 بايت شاملة المسافات؛ الأكبر 413 و ~too large\n~ قبل أي تحليل.

أي مسار تاني، منهم ~/berths/~ و ~/Berths~ و ~/berths/2/x~، يرجع 404 و ~no route\n~. طريقة مش مدعومة على مسار موجود ترجع 405 و ~use route methods\n~، و ~Allow~ فيه طرق المسار (~HEAD~ مسموح زيادة). كل رد، حتى 404 و 405 و 413، فيه ~Cache-Control: no-store~ من middleware واحد. النجاح بجسم نوعه ~application/json~ والأخطاء ~text/plain~. الـ query بتتجاهل، والفشل ميغيّرش الحالة. مفيش شبكة ولا ~Report~، و ~main~ فاضية.

في ~compare.go~: ~Comparison~ فيها 4 عناصر بالظبط، كل موضوع من ~handler~ و ~errors~ و ~testing~ و ~lifetime~ مرة واحدة. ~Gin~ و ~Fiber~ في كل عنصر مش فاضيين بعد trim، وكل واحد 160 بايت أو أقل، والاتنين مختلفين.`,
		Hints: []string{"ابدأ بـ fiber.New(fiber.Config{...}) فيه StrictRouting وCaseSensitive وErrorHandler، وبعدين app.Use لـ Cache-Control.", "في ErrorHandler فك الخطأ بـ errors.As لـ *fiber.Error واكتب no route أو use route methods حسب Code.", "افحص len(c.Body()) قبل json.Unmarshal، وحوّل c.Params(\"id\") برقم قياسي، وخلّي المخزن تحت sync.Mutex؛ وبعدين املا Comparison."},
		En: &LevelText{Title: "Fiber", Goal: "Run the same berth contract on Fiber, test it with app.Test, and write a short comparison of the two frameworks", Files: []LabFile{{Name: "main.go", Starter: ch282Starter}, {Name: "compare.go", Starter: ch282CompareEn}},
			Lesson: `# The Same Registry on Fiber

## Before you write
- Fiber is not built on ~net/http~; what changes in testing and in data lifetime?
- Where does an error returned from a handler go, and who writes its reply?

## Fiber on top of fasthttp
Import ~github.com/gofiber/fiber/v2~; its package name is ~fiber~. Fiber is built on fasthttp, not ~net/http~, so ~*fiber.App~ is not an ~http.Handler~ and has no ~ServeHTTP~. A handler has the shape ~func(c *fiber.Ctx) error~.
~~~go
app := fiber.New()
app.Get("/probes/:name", func(c *fiber.Ctx) error {
 return c.Status(fiber.StatusOK).SendString("probe " + c.Params("name") + "\n")
})
~~~
~app.Get~, ~app.Post~, ~app.Put~, and ~app.Delete~ register a path for a method. ~c.Params("name")~ reads the ~:name~ segment, ~c.Body()~ returns the whole body as bytes, and ~c.Method()~ the method.

## Replying
~c.Status~ sets the status and returns the same ctx so you can chain. ~c.SendString~ writes text; the default content type is ~text/plain; charset=utf-8~. ~c.JSON~ encodes the value and sets ~application/json~. ~c.Set~ sets a response header.
~~~go
c.Set("Location", "/probes/7")
return c.Status(fiber.StatusCreated).JSON(probe)
~~~
~c.SendStatus(204)~ sets the status and, when the body is empty, writes the status name as text; the server drops that body by itself for 204. ~c.Status(204)~ followed by ~return nil~ is clearer.

## Testing with app.Test
~app.Test(req, -1)~ takes an ordinary ~*http.Request~ from ~httptest.NewRequest~, runs it in memory without a port, and returns an ~*http.Response~ and an error. The default timeout is one second; ~-1~ disables it so parallel requests under ~-race~ do not fail for being slow. Close ~resp.Body~ after reading.
~~~go
resp, err := app.Test(httptest.NewRequest("GET", "/probes/vega", nil), -1)
if err != nil { t.Fatal(err) }
defer resp.Body.Close()
body, _ := io.ReadAll(resp.Body)
~~~

## Configuration that protects the contract
~fiber.New~ takes an optional ~fiber.Config~. ~StrictRouting~ defaults to ~false~, so ~/berths/~ matches ~/berths~. ~CaseSensitive~ defaults to ~false~, so ~/Berths~ matches. Both must be ~true~ here. The default ~BodyLimit~ is 4 MiB; that is a server ceiling, not your contract's cap, so check ~len(c.Body())~ yourself.
~~~go
app := fiber.New(fiber.Config{StrictRouting: true, CaseSensitive: true, ErrorHandler: onFault})
~~~

## Errors are returned, not written
A handler returns an ~error~, and Fiber sends every error to one ~ErrorHandler~. An unknown path reaches it as a ~*fiber.Error~ whose ~Code~ is 404. A wrong method on an existing path reaches it as 405, together with an ~Allow~ header Fiber filled with the methods that have a route. Check the type with ~errors.As~.
~~~go
func onFault(c *fiber.Ctx, err error) error {
 var fe *fiber.Error
 if errors.As(err, &fe) && fe.Code == fiber.StatusNotFound {
  return c.Status(404).SendString("lost\n")
 }
 return c.Status(500).SendString("fault\n")
}
~~~
~fiber.NewError(code, msg)~ builds the same type if you want to return one from a handler. ~app.Get~ also registers ~HEAD~ for the same path, so ~Allow~ may show ~HEAD~, which is allowed.

## Middleware, Fiber style
A ~net/http~ middleware that wraps an ~http.Handler~ cannot be plugged in here. In Fiber it is an ordinary handler that calls ~c.Next()~ and returns the error that came back. ~app.Use~ adds it in front of every route, even a request that will end in 404.
~~~go
app.Use(func(c *fiber.Ctx) error {
 c.Set("X-Deck", "7")
 return c.Next() // Return the rest's error so it reaches ErrorHandler
})
~~~

## Data lifetime
fasthttp reuses request buffers once the handler returns. A string from ~c.Params~ or ~c.Get~ is valid only inside the handler; copy it with ~strings.Clone~ before storing it. ~json.Unmarshal~ copies strings itself, so fields decoded from ~c.Body()~ are safe.
~~~go
name := strings.Clone(c.Params("name")) // A copy that outlives the reply
~~~
~c.BodyParser~ picks a decoder from ~Content-Type~ and fails when the header is missing. The contract does not require a content type, so check ~c.Body()~ with ~json.Unmarshal~ as in 26-2; it already rejects a second value after the first.

## The comparison file
In ~compare.go~ fill ~Comparison~ with four ~Contrast~ values, one per topic: ~handler~ (handler shape and interface), ~errors~ (where 404 and 405 are written), ~testing~ (how you test without a port), and ~lifetime~ (what must not outlive the reply). Each cell is a short sentence about that framework specifically.
~~~go
{Topic: "testing", Gin: "...", Fiber: "..."}
~~~

## Common mistakes
~~~go
app := fiber.New()                        // /berths/ and /Berths pass
if err := c.BodyParser(&in); err != nil {} // Fails without Content-Type
seen[c.Params("id")] = true                // The buffer is reused
resp, _ := app.Test(req)                   // One-second timeout, error ignored
~~~

## Your mission
Implement ~NewBerthApp() *fiber.App~ returning a new app with an independent concurrent-safe store and exactly the 28-1 contract over the same ~Berth~. ~GET /berths~ returns 200 and an array in ascending ID order, ~[]~ when empty. ~POST /berths~ accepts one JSON object: a ~label~ that is nonblank after trimming and an integer ~slots~ in inclusive 1..64; input ~id~ and unknown fields are ignored. IDs start at 1, increase, and are never reused; rejected requests consume none. Success is 201 with the stored object and ~Location: /berths/ID~.

~GET /berths/ID~ returns 200 and the object. ~PUT~ replaces both fields with the same validation, keeps the ID, and returns 200. ~DELETE~ returns 204 with an empty body. A missing item is 404 and ~unknown berth\n~. Invalid JSON, a wrong type, an out-of-range value, or an ID that is not a canonical positive integer is 400 and ~invalid berth\n~. For ~PUT~, check the ID, then the body, then look up. Bodies are capped at 4096 bytes including whitespace; larger ones are 413 and ~too large\n~ before any parsing.

Every other path, including ~/berths/~, ~/Berths~, and ~/berths/2/x~, returns 404 and ~no route\n~. An unsupported method on an existing path returns 405, ~use route methods\n~, and an ~Allow~ listing the path's methods (an extra ~HEAD~ is allowed). Every reply, including 404, 405, and 413, carries ~Cache-Control: no-store~ from one middleware. Successful bodies are ~application/json~ and errors are ~text/plain~. The query is ignored, and failures never change state. No network and no ~Report~; leave ~main~ empty.

In ~compare.go~: ~Comparison~ holds exactly 4 entries, each topic of ~handler~, ~errors~, ~testing~, and ~lifetime~ once. In every entry ~Gin~ and ~Fiber~ are nonblank after trimming, each at most 160 bytes, and different from each other.`,
			Hints: []string{"Start with fiber.New(fiber.Config{...}) holding StrictRouting, CaseSensitive, and ErrorHandler, then app.Use for Cache-Control.", "In ErrorHandler unwrap the error with errors.As into *fiber.Error and write no route or use route methods based on Code.", "Check len(c.Body()) before json.Unmarshal, convert c.Params(\"id\") as a canonical number, and keep the store under a sync.Mutex; then fill Comparison."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch282Hidden, HiddenFiles: map[string]string{"berth_contract_test.go": ch28BerthContract}},
	})
}

const ch282Starter = `package main

import "github.com/gofiber/fiber/v2"

type Berth struct {
	ID    int    ` + "`" + `json:"id"` + "`" + `
	Label string ` + "`" + `json:"label"` + "`" + `
	Slots int    ` + "`" + `json:"slots"` + "`" + `
}

func NewBerthApp() *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error { return c.Status(fiber.StatusNotImplemented).SendString("TODO\n") })
	return app
}
func main() {}
`

const ch282CompareAr = `package main

// Contrast بيقارن موضوع واحد بين Gin وFiber في جملة قصيرة لكل إطار.
type Contrast struct {
	Topic string
	Gin   string
	Fiber string
}

// اكتب هنا أربع مقارنات: handler وerrors وtesting وlifetime.
var Comparison = []Contrast{}
`

const ch282CompareEn = `package main

// Contrast compares one topic between Gin and Fiber in one short sentence per framework.
type Contrast struct {
	Topic string
	Gin   string
	Fiber string
}

// Write four comparisons here: handler, errors, testing, and lifetime.
var Comparison = []Contrast{}
`

const ch282Hidden = `package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func fiberCaller(t *testing.T) contractCaller {
	app := NewBerthApp()
	if app == nil {
		t.Fatal("NewBerthApp returned nil")
	}
	return func(method, path, body string) contractReply {
		resp, err := app.Test(httptest.NewRequest(method, path, strings.NewReader(body)), -1)
		if err != nil {
			return contractReply{Code: -1, Header: http.Header{}, Body: err.Error()}
		}
		defer resp.Body.Close()
		raw, _ := io.ReadAll(resp.Body)
		return contractReply{resp.StatusCode, resp.Header, string(raw)}
	}
}

func TestHiddenFiberBerthContract(t *testing.T) {
	contractCRUD(t, func() contractCaller { return fiberCaller(t) }, 282)
}

func TestHiddenFiberBerthsParallel(t *testing.T) {
	contractParallel(t, func() contractCaller { return fiberCaller(t) })
}

// Boundaries: exact topic set, duplicates, blank or oversized cells, identical sides.
func TestHiddenFiberComparison(t *testing.T) {
	want := map[string]bool{"handler": true, "errors": true, "testing": true, "lifetime": true}
	if len(Comparison) != 4 {
		t.Fatalf("Comparison has %d entries, want 4", len(Comparison))
	}
	for _, c := range Comparison {
		if !want[c.Topic] {
			t.Fatalf("unknown or repeated topic %q", c.Topic)
		}
		delete(want, c.Topic)
		g, f := strings.TrimSpace(c.Gin), strings.TrimSpace(c.Fiber)
		if g == "" || f == "" || len(c.Gin) > 160 || len(c.Fiber) > 160 || g == f {
			t.Fatalf("topic %q needs two different short sentences: %q / %q", c.Topic, c.Gin, c.Fiber)
		}
	}
}
`

func init() {
	add(&Level{ID: "28-3-openapi", Chapter: 28, Kind: "lab", Title: "عقد الـ API", Goal: "اقرا مستند OpenAPI لسجل المرحّلات، واتحقق إن كل رد من السيرفر مطابق له",
		Requires: []string{"sel:json.Unmarshal"}, Bridge: []int{58}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch283Starter}, {Name: "openapi.json", Starter: ch283Doc, ReadOnly: true}},
		Lesson: `# عقد مكتوب لسجل المرحّلات

## قبل ما تكتب
- لو المستند قال حاجة والسيرفر عمل حاجة تانية، مين يكتشف الفرق قبل العميل؟
- إيه الفرق بين «المستند مكتوب صح» و «السيرفر ماشي على المستند»؟

## OpenAPI في سطور
OpenAPI مستند JSON بيوصف الـ API: المسارات والطرق والحالات وشكل كل جسم. ~openapi~ رقم نسخة المواصفة، و ~paths~ خريطة من قالب مسار لعمليات. كل عملية (~get~ و ~post~ و ~delete~ ...) فيها ~responses~: خريطة من الحالة مكتوبة نص لرد. الرد فيه ~content~: خريطة من نوع المحتوى لـ ~schema~.
~~~json
{"openapi": "3.1.0",
 "paths": {"/probes/{name}": {"get": {"responses": {
   "200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Probe"}}}},
   "204": {"description": "nothing to report"}}}}},
 "components": {"schemas": {"Probe": {"type": "object", "required": ["name"]}}}}
~~~
~{name}~ في المسار قالب بيطابق جزء واحد مش فاضي. الرد اللي مفيهوش ~content~ معناه إن جسمه فاضي، زي 204. ~$ref~ إشارة لـ schema متسمّية تحت ~components.schemas~؛ هنا بنقبل بس الشكل ~#/components/schemas/Name~.

## فك المستند في structs
~json.Unmarshal~ بتملا struct من غير ما تحتاج كل المفاتيح؛ المفتاح اللي ملوش حقل بيتجاهل، زي ~description~ و ~operationId~. ~$ref~ مينفعش يبقى اسم حقل Go، فالـ tag هو اللي بيربطه.
~~~go
type schema struct {
 Ref        string             ` + "`" + `json:"$ref"` + "`" + `
 Type       string             ` + "`" + `json:"type"` + "`" + `
 Properties map[string]*schema ` + "`" + `json:"properties"` + "`" + `
 MinLength  *int               ` + "`" + `json:"minLength"` + "`" + `
}
~~~
الـ pointer في ~MinLength~ بيفرّق بين «مش مكتوب» و «صفر». الـ struct بيشاور على نفسه من خلال pointer أو خريطة، وده اللي بيخلّي الـ schema تتداخل لأي عمق.

## قيمة JSON من غير شكل معروف
الجسم اللي هتتحقق منه شكله مش معروف وقت الكتابة، فافكّه في ~any~. ~json.Unmarshal~ بتحط الكائن ~map[string]any~، والمصفوفة ~[]any~، والنص ~string~، وأي رقم ~float64~، و ~true~ أو ~false~ ~bool~، و ~null~ ~nil~.
~~~go
var v any
err := json.Unmarshal([]byte(` + "`" + `{"depth": 3, "tags": ["a"]}` + "`" + `), &v)
m, ok := v.(map[string]any)    // ok لو كائن
depth, ok := m["depth"].(float64) // 3
~~~
الـ type switch من فصل البروتوكولات بيختار حسب النوع الحقيقي: ~switch x := v.(type) { case string: ... }~. الرقم الصحيح ~float64~ من غير كسر؛ ~math.Trunc~ بتشيل الكسر، فلو ~f == math.Trunc(f)~ يبقى صحيح.
~~~go
f := 2.5
whole := f == math.Trunc(f) // false
~~~

## التحقق بدالة بتنادي نفسها
كل كلمة في الـ schema قاعدة مستقلة. ~type~ النوع، و ~required~ مفاتيح لازم تبقى موجودة، و ~properties~ schema لكل مفتاح موجود، و ~additionalProperties: false~ بتمنع أي مفتاح مش في ~properties~. ~items~ schema لكل عنصر في مصفوفة، و ~enum~ قايمة القيم المسموحة، و ~minimum~ أقل رقم، و ~minLength~ أقل عدد حروف. الدالة بتنادي نفسها على الأبناء.
~~~go
func walk(s *schema, v any) error {
 if s.Ref != "" { return walk(named[strings.TrimPrefix(s.Ref, "#/components/schemas/")], v) }
 if s.Type == "string" {
  text, ok := v.(string)
  if !ok { return errors.New("want string") }
  if s.MinLength != nil && utf8.RuneCountInString(text) < *s.MinLength { return errors.New("too short") }
 }
 return nil
}
~~~
~minLength~ بيعد حروف مش بايتات: ~"🚀"~ حرف واحد. اتأكد من كل ~$ref~ في ~LoadSpec~ مرة واحدة، عشان الإشارة المكسورة تبان وقت التحميل مش في نص طلب.

## مطابقة المسار والنوع
~strings.Split("/relays/7", "/")~ بترجع ~["", "relays", "7"]~. قارن الأجزاء واحد واحد مع القالب: الجزء الحرفي لازم يتساوى، والجزء اللي بين ~{~ و ~}~ بيقبل أي جزء مش فاضي. ~mime.ParseMediaType~ من 26-6 بتفصل ~application/json; charset=utf-8~ لنوع وباراميترات؛ قارن النوع بس.

## اختبار العقد
اختبار العقد بيشغّل الـ API الحقيقي في الذاكرة بـ ~httptest~، وكل رد بيتعرض على المستند: المسار موجود، والطريقة موجودة، والحالة متوثقة، والجسم مطابق للـ schema. أي فرق معناه إن واحد فيهم اتغيّر من غير التاني.
~~~go
w := httptest.NewRecorder()
api.ServeHTTP(w, httptest.NewRequest("GET", "/relays/1", nil))
err := spec.Check("GET", "/relays/1", w.Code, w.Header().Get("Content-Type"), w.Body.Bytes())
~~~
المستند هنا ~openapi.json~، ملف للقراية بس جنب ~main.go~. ~//go:embed~ من 26-5 بيضمّه في البرنامج، والتستات بتقراه بـ ~os.ReadFile("openapi.json")~ من فولدر الباكدج.

## أخطاء شائعة
~~~go
n := v.(float64)            // panic لو القيمة نص؛ استخدم الشكل اللي فيه ok
len(name) < *s.MinLength    // بايتات مش حروف
ops := byPath["/relays/7"]  // المساواة عمرها ما هتطابق "/relays/{id}"
return nil                  // ونسيت additionalProperties: false
~~~

## المهمة
نفّذ ~LoadSpec(raw []byte) (*Spec, error)~: يفك مستند JSON واحد، ويرجع ~nil,error~ لو: JSON بايظ، أو ~openapi~ مش بيبدأ بـ ~3.~، أو ~paths~ فاضية، أو مسار مش بيبدأ بـ ~/~، أو عملية من غير ردود، أو مفتاح رد مش 3 أرقام بين 100 و 599 (زي ~2XX~ و ~default~)، أو ~$ref~ في أي schema مش بالشكل ~#/components/schemas/Name~ أو الاسم مش موجود. العمليات المقروءة ~get~ و ~put~ و ~post~ و ~delete~ و ~patch~؛ أي مفتاح تاني بيتجاهل.

نفّذ ~(s *Spec) Check(method, path string, status int, contentType string, body []byte) error~ يرجع ~nil~ بس لو الرد مطابق. ~path~ مسار من غير query، بيطابق قالب بنفس عدد الأجزاء: الحرفي متساوي، و ~{x}~ أي جزء مش فاضي. الطريقة بأي حالة حروف. مسار أو طريقة أو حالة مش متوثقين: خطأ. رد من غير ~content~: الجسم لازم فاضي والنوع بيتجاهل. رد بـ ~content~: نوع ~contentType~ لازم ~application/json~ وموجود في ~content~، والجسم قيمة JSON واحدة مطابقة للـ schema.

قواعد الـ schema: ~$ref~، و ~type~ (~object~ و ~array~ و ~string~ و ~integer~ و ~number~ و ~boolean~ و ~null~؛ ومن غير ~type~ أي نوع)، و ~required~، و ~properties~، و ~additionalProperties: false~، و ~items~، و ~enum~ لقيم بسيطة، و ~minimum~ للأرقام، و ~minLength~ بالحروف. أي كلمة تانية بتتجاهل. ~integer~ رقم من غير كسر.

نفّذ ~NewRelayAPI() http.Handler~ بمخزن مستقل ماشي على ~openapi.json~. ~GET /relays~ يرجع 200 ومصفوفة مرتبة بالـ ID، والفاضية ~[]~. ~POST /relays~ جسمه حده 4096 بايت، وكائن واحد مطابق لـ ~RelayInput~: ~name~ نص مش فاضي، و ~band~ واحد من ~S~ و ~X~ و ~Ka~، و ~tags~ اختياري فيه نصوص مش فاضية، ومفيش مفاتيح تانية. النجاح 201 و ~Relay~ بـ ~id~ بيبدأ من 1 ويزيد، و ~online~ بـ ~false~، و ~tags~ زي ما اتبعتت أو ~[]~، و ~Location: /relays/ID~. أي غلط 400 و ~{"error":"invalid relay"}~ من غير ما ياخد رقم.

~GET /relays/{id}~ يرجع 200 والعنصر، و ~DELETE~ يرجع 204 بجسم فاضي. الـ ID الغائب أو اللي مش رقم موجب قياسي: 404 و ~{"error":"unknown relay"}~. ردود JSON بنوع ~application/json~. أي طريقة أو مسار برا المستند مش جزء من المهمة. مفيش شبكة ولا ~Report~، و ~main~ فاضية.`,
		Hints: []string{"اعمل structs للمستند: document فيه OpenAPI وPaths وComponents، وpathItem بحقول get وput وpost وdelete وpatch، وschema فيها Ref بـ tag json:\"$ref\".", "في LoadSpec لف على كل schema (properties وitems) واتأكد من كل $ref، وافحص مفاتيح الردود بـ strconv.Atoi وطولها 3.", "في Check فك الجسم في any واكتب validate بـ type switch؛ وفي NewRelayAPI ضمّ openapi.json بـ go:embed واتحقق من الجسم بـ RelayInput نفسها."},
		En: &LevelText{Title: "The API Contract", Goal: "Read the relay registry's OpenAPI document and check that every server reply matches it", Files: []LabFile{{Name: "main.go", Starter: ch283Starter}, {Name: "openapi.json", Starter: ch283Doc, ReadOnly: true}},
			Lesson: `# A Written Contract for the Relay Registry

## Before you write
- If the document says one thing and the server does another, who notices before a client does?
- What is the difference between "the document is well written" and "the server follows the document"?

## OpenAPI in brief
OpenAPI is a JSON document describing an API: paths, methods, statuses, and the shape of every body. ~openapi~ is the specification version, and ~paths~ maps a path template to operations. Each operation (~get~, ~post~, ~delete~, ...) has ~responses~: a map from the status written as text to a response. A response has ~content~: a map from content type to ~schema~.
~~~json
{"openapi": "3.1.0",
 "paths": {"/probes/{name}": {"get": {"responses": {
   "200": {"content": {"application/json": {"schema": {"$ref": "#/components/schemas/Probe"}}}},
   "204": {"description": "nothing to report"}}}}},
 "components": {"schemas": {"Probe": {"type": "object", "required": ["name"]}}}}
~~~
~{name}~ in a path is a template matching one nonempty segment. A response without ~content~ means an empty body, like 204. ~$ref~ points to a named schema under ~components.schemas~; here only the form ~#/components/schemas/Name~ is accepted.

## Decoding the document into structs
~json.Unmarshal~ fills a struct without needing every key; a key with no field is ignored, like ~description~ and ~operationId~. ~$ref~ cannot be a Go field name, so the tag connects it.
~~~go
type schema struct {
 Ref        string             ` + "`" + `json:"$ref"` + "`" + `
 Type       string             ` + "`" + `json:"type"` + "`" + `
 Properties map[string]*schema ` + "`" + `json:"properties"` + "`" + `
 MinLength  *int               ` + "`" + `json:"minLength"` + "`" + `
}
~~~
The pointer in ~MinLength~ tells "not written" from "zero". The struct refers to itself through a pointer or a map, which lets schemas nest to any depth.

## A JSON value of unknown shape
The body you check has no shape known at compile time, so decode it into ~any~. ~json.Unmarshal~ produces ~map[string]any~ for an object, ~[]any~ for an array, ~string~ for text, ~float64~ for every number, ~bool~ for ~true~ or ~false~, and ~nil~ for ~null~.
~~~go
var v any
err := json.Unmarshal([]byte(` + "`" + `{"depth": 3, "tags": ["a"]}` + "`" + `), &v)
m, ok := v.(map[string]any)    // ok when it is an object
depth, ok := m["depth"].(float64) // 3
~~~
The type switch from the protocols chapter picks by the real type: ~switch x := v.(type) { case string: ... }~. An integer is a ~float64~ with no fraction; ~math.Trunc~ drops the fraction, so ~f == math.Trunc(f)~ means whole.
~~~go
f := 2.5
whole := f == math.Trunc(f) // false
~~~

## Validation by a function that calls itself
Each schema keyword is an independent rule. ~type~ is the type, ~required~ lists keys that must exist, ~properties~ gives a schema for each present key, and ~additionalProperties: false~ forbids any key missing from ~properties~. ~items~ is the schema for each array element, ~enum~ lists the allowed values, ~minimum~ is the lowest number, and ~minLength~ the fewest characters. The function calls itself on children.
~~~go
func walk(s *schema, v any) error {
 if s.Ref != "" { return walk(named[strings.TrimPrefix(s.Ref, "#/components/schemas/")], v) }
 if s.Type == "string" {
  text, ok := v.(string)
  if !ok { return errors.New("want string") }
  if s.MinLength != nil && utf8.RuneCountInString(text) < *s.MinLength { return errors.New("too short") }
 }
 return nil
}
~~~
~minLength~ counts characters, not bytes: ~"🚀"~ is one. Check every ~$ref~ once in ~LoadSpec~, so a broken reference shows up at load time instead of in the middle of a request.

## Matching paths and content types
~strings.Split("/relays/7", "/")~ returns ~["", "relays", "7"]~. Compare segments one by one with the template: a literal segment must be equal, and a segment between ~{~ and ~}~ accepts any nonempty segment. ~mime.ParseMediaType~ from 26-6 splits ~application/json; charset=utf-8~ into a type and parameters; compare the type only.

## Contract testing
A contract test runs the real API in memory with ~httptest~ and holds every reply against the document: the path exists, the method exists, the status is documented, and the body matches the schema. Any difference means one of them changed without the other.
~~~go
w := httptest.NewRecorder()
api.ServeHTTP(w, httptest.NewRequest("GET", "/relays/1", nil))
err := spec.Check("GET", "/relays/1", w.Code, w.Header().Get("Content-Type"), w.Body.Bytes())
~~~
The document here is ~openapi.json~, a read-only file next to ~main.go~. ~//go:embed~ from 26-5 builds it into the program, and the tests read it with ~os.ReadFile("openapi.json")~ from the package folder.

## Common mistakes
~~~go
n := v.(float64)            // Panics when the value is text; use the ok form
len(name) < *s.MinLength    // Bytes, not characters
ops := byPath["/relays/7"]  // Equality never matches "/relays/{id}"
return nil                  // And additionalProperties: false was forgotten
~~~

## Your mission
Implement ~LoadSpec(raw []byte) (*Spec, error)~: decode one JSON document and return ~nil,error~ when: the JSON is invalid, ~openapi~ does not start with ~3.~, ~paths~ is empty, a path does not start with ~/~, an operation has no responses, a response key is not 3 digits between 100 and 599 (such as ~2XX~ or ~default~), or a ~$ref~ in any schema is not of the form ~#/components/schemas/Name~ or names a missing schema. The operations read are ~get~, ~put~, ~post~, ~delete~, and ~patch~; any other key is ignored.

Implement ~(s *Spec) Check(method, path string, status int, contentType string, body []byte) error~ returning ~nil~ only when the reply matches. ~path~ has no query and matches a template with the same number of segments: literals are equal, and ~{x}~ is any nonempty segment. The method is case-insensitive. An undocumented path, method, or status is an error. For a response without ~content~ the body must be empty and the type is ignored. For a response with ~content~, the type of ~contentType~ must be ~application/json~ and present in ~content~, and the body must be one JSON value matching the schema.

Schema rules: ~$ref~, ~type~ (~object~, ~array~, ~string~, ~integer~, ~number~, ~boolean~, ~null~; without ~type~ anything goes), ~required~, ~properties~, ~additionalProperties: false~, ~items~, ~enum~ of simple values, ~minimum~ for numbers, and ~minLength~ in characters. Any other keyword is ignored. An ~integer~ is a number with no fraction.

Implement ~NewRelayAPI() http.Handler~ with an independent store that follows ~openapi.json~. ~GET /relays~ returns 200 and an array in ID order, ~[]~ when empty. ~POST /relays~ has a 4096-byte body cap and takes one object matching ~RelayInput~: ~name~ is nonempty text, ~band~ is one of ~S~, ~X~, and ~Ka~, ~tags~ is optional nonempty text, and there are no other keys. Success is 201 and a ~Relay~ whose ~id~ starts at 1 and increases, with ~online~ ~false~, ~tags~ as sent or ~[]~, and ~Location: /relays/ID~. Any error is 400 and ~{"error":"invalid relay"}~ without consuming an ID.

~GET /relays/{id}~ returns 200 and the item, and ~DELETE~ returns 204 with an empty body. A missing ID or one that is not a canonical positive integer is 404 and ~{"error":"unknown relay"}~. JSON replies use ~application/json~. Methods or paths outside the document are not part of the mission. No network and no ~Report~; leave ~main~ empty.`,
			Hints: []string{"Make structs for the document: document with OpenAPI, Paths, and Components, a pathItem with get, put, post, delete, and patch fields, and a schema with Ref tagged json:\"$ref\".", "In LoadSpec walk every schema (properties and items) and verify every $ref, and check response keys with strconv.Atoi and length 3.", "In Check decode the body into any and write validate with a type switch; in NewRelayAPI embed openapi.json with go:embed and validate the body against RelayInput itself."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch283Hidden},
	})
}

const ch283Starter = `package main

import (
	"errors"
	"net/http"
)

type Spec struct{}

func LoadSpec(raw []byte) (*Spec, error) { return nil, errors.New("TODO") }
func (s *Spec) Check(method, path string, status int, contentType string, body []byte) error {
	return errors.New("TODO")
}
func NewRelayAPI() http.Handler { return http.NotFoundHandler() }
func main()                     {}
`

const ch283Doc = `{
  "openapi": "3.1.0",
  "info": {"title": "Relay Registry", "version": "1.0.0"},
  "paths": {
    "/relays": {
      "get": {
        "operationId": "listRelays",
        "summary": "All relays in ascending id order.",
        "responses": {
          "200": {"description": "The relays.", "content": {"application/json": {"schema": {"type": "array", "items": {"$ref": "#/components/schemas/Relay"}}}}}
        }
      },
      "post": {
        "operationId": "createRelay",
        "summary": "Register a relay. Bodies over 4096 bytes are rejected.",
        "requestBody": {"required": true, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/RelayInput"}}}},
        "responses": {
          "201": {"description": "Registered; Location names the new relay.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Relay"}}}},
          "400": {"description": "Not one valid RelayInput.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Error"}}}}
        }
      }
    },
    "/relays/{id}": {
      "parameters": [{"name": "id", "in": "path", "required": true, "description": "Canonical positive integer.", "schema": {"type": "integer", "minimum": 1}}],
      "get": {
        "operationId": "getRelay",
        "responses": {
          "200": {"description": "The relay.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Relay"}}}},
          "404": {"description": "Unknown or malformed id.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Error"}}}}
        }
      },
      "delete": {
        "operationId": "deleteRelay",
        "responses": {
          "204": {"description": "Removed."},
          "404": {"description": "Unknown or malformed id.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/Error"}}}}
        }
      }
    }
  },
  "components": {
    "schemas": {
      "Band": {"type": "string", "enum": ["S", "X", "Ka"]},
      "Relay": {
        "type": "object",
        "required": ["id", "name", "band", "online", "tags"],
        "additionalProperties": false,
        "properties": {
          "id": {"type": "integer", "minimum": 1},
          "name": {"type": "string", "minLength": 1},
          "band": {"$ref": "#/components/schemas/Band"},
          "online": {"type": "boolean"},
          "tags": {"type": "array", "items": {"type": "string", "minLength": 1}}
        }
      },
      "RelayInput": {
        "type": "object",
        "required": ["name", "band"],
        "additionalProperties": false,
        "properties": {
          "name": {"type": "string", "minLength": 1},
          "band": {"$ref": "#/components/schemas/Band"},
          "tags": {"type": "array", "items": {"type": "string", "minLength": 1}}
        }
      },
      "Error": {
        "type": "object",
        "required": ["error"],
        "additionalProperties": false,
        "properties": {"error": {"type": "string", "minLength": 1}}
      }
    }
  }
}
`

const ch283Hidden = `package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"mime"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func hiddenSpecFile(t *testing.T) []byte {
	t.Helper()
	raw, err := os.ReadFile("openapi.json")
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func hiddenSpec(t *testing.T) *Spec {
	t.Helper()
	s, err := LoadSpec(hiddenSpecFile(t))
	if err != nil || s == nil {
		t.Fatal("LoadSpec(openapi.json):", err)
	}
	return s
}

// hiddenMutate decodes the shipped document, applies change, and re-encodes it.
func hiddenMutate(t *testing.T, change func(doc map[string]any)) []byte {
	t.Helper()
	var doc map[string]any
	if err := json.Unmarshal(hiddenSpecFile(t), &doc); err != nil {
		t.Fatal(err)
	}
	change(doc)
	raw, _ := json.Marshal(doc)
	return raw
}

func hiddenDig(doc map[string]any, keys ...string) map[string]any {
	for _, k := range keys {
		doc = doc[k].(map[string]any)
	}
	return doc
}

func hiddenEncode(v any) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

func hiddenName(rng *rand.Rand) string {
	return []string{"vega", "سهيل", "🚀", "deck 7", "Ka-band"}[rng.Intn(5)] + fmt.Sprint(rng.Intn(1000))
}

// Boundaries: version, empty paths, path shape, empty responses, status keys, local/dangling/external refs.
func TestHiddenSpecLoadRejectsBrokenDocuments(t *testing.T) {
	hiddenSpec(t)
	raw := hiddenSpecFile(t)
	extra := hiddenMutate(t, func(d map[string]any) {
		d["x-colony"] = true
		hiddenDig(d, "paths", "/relays", "get")["tags"] = []any{"relays"}
	})
	if s, err := LoadSpec(extra); s == nil || err != nil {
		t.Fatal("unknown keys must be ignored:", err)
	}
	cases := map[string][]byte{
		"truncated":      raw[:len(raw)/2],
		"trailing value": append(append([]byte{}, raw...), []byte(" {}")...),
		"not an object":  []byte("[]"),
		"swagger 2":      hiddenMutate(t, func(d map[string]any) { d["openapi"] = "2.0" }),
		"no version":     hiddenMutate(t, func(d map[string]any) { delete(d, "openapi") }),
		"no paths":       hiddenMutate(t, func(d map[string]any) { d["paths"] = map[string]any{} }),
		"relative path": hiddenMutate(t, func(d map[string]any) {
			p := hiddenDig(d, "paths")
			p["relays"] = p["/relays"]
			delete(p, "/relays")
		}),
		"no responses": hiddenMutate(t, func(d map[string]any) {
			hiddenDig(d, "paths", "/relays/{id}", "delete")["responses"] = map[string]any{}
		}),
		"range status": hiddenMutate(t, func(d map[string]any) {
			r := hiddenDig(d, "paths", "/relays", "get", "responses")
			r["2XX"] = r["200"]
		}),
		"default status": hiddenMutate(t, func(d map[string]any) {
			r := hiddenDig(d, "paths", "/relays", "get", "responses")
			r["default"] = r["200"]
		}),
		"status 600": hiddenMutate(t, func(d map[string]any) {
			r := hiddenDig(d, "paths", "/relays", "get", "responses")
			r["600"] = r["200"]
		}),
		"dangling ref": hiddenMutate(t, func(d map[string]any) { delete(hiddenDig(d, "components", "schemas"), "Error") }),
		"deep dangling ref": hiddenMutate(t, func(d map[string]any) {
			hiddenDig(d, "components", "schemas", "Relay", "properties")["tags"] = map[string]any{"type": "array", "items": map[string]any{"$ref": "#/components/schemas/Ghost"}}
		}),
		"external ref": hiddenMutate(t, func(d map[string]any) {
			hiddenDig(d, "paths", "/relays/{id}", "get", "responses", "200", "content", "application/json")["schema"] = map[string]any{"$ref": "relay.json#/Relay"}
		}),
	}
	for name, doc := range cases {
		if s, err := LoadSpec(doc); s != nil || err == nil {
			t.Errorf("%s: LoadSpec accepted a broken document", name)
		}
	}
}

func hiddenRelay(rng *rand.Rand, id int) map[string]any {
	r := map[string]any{"id": id, "name": hiddenName(rng), "band": []string{"S", "X", "Ka"}[rng.Intn(3)], "online": rng.Intn(2) == 0, "tags": []any{}}
	for i := rng.Intn(4); i > 0; i-- {
		r["tags"] = append(r["tags"].([]any), fmt.Sprintf("t%d", rng.Intn(99)))
	}
	return r
}

// Boundaries: templates, method case, undocumented status/method/path, empty-body responses,
// content types, every schema keyword, refs to refs, rune-counted lengths, integral floats.
func TestHiddenSpecCheckResponses(t *testing.T) {
	s := hiddenSpec(t)
	rng := rand.New(rand.NewSource(2832))
	const js = "application/json"
	ok := func(label, method, path string, status int, ct, body string) {
		t.Helper()
		if err := s.Check(method, path, status, ct, []byte(body)); err != nil {
			t.Errorf("%s: valid reply rejected: %v", label, err)
		}
	}
	bad := func(label, method, path string, status int, ct, body string) {
		t.Helper()
		if s.Check(method, path, status, ct, []byte(body)) == nil {
			t.Errorf("%s: contract violation accepted: %s %s %d %q %s", label, method, path, status, ct, body)
		}
	}
	var list []any
	for i := 0; i < 20; i++ {
		r := hiddenRelay(rng, 1+rng.Intn(5000))
		id := fmt.Sprint(r["id"])
		ok("relay", "GET", "/relays/"+id, 200, js+"; charset=utf-8", hiddenEncode(r))
		ok("created", "POST", "/relays", 201, js, hiddenEncode(r))
		list = append(list, r)
	}
	ok("empty list", "GET", "/relays", 200, js, "[]")
	ok("list", "get", "/relays", 200, js, hiddenEncode(list))
	ok("integral float", "GET", "/relays/3", 200, js, ` + "`" + `{"id":3.0,"name":"🚀","band":"Ka","online":true,"tags":[]}` + "`" + `)
	ok("not found", "GET", "/relays/abc", 404, js, ` + "`" + `{"error":"unknown relay"}` + "`" + `)
	ok("bad input", "POST", "/relays", 400, js, ` + "`" + `{"error":"invalid relay"}` + "`" + `)
	ok("deleted", "DELETE", "/relays/9", 204, "", "")
	ok("deleted with type", "DELETE", "/relays/9", 204, "text/plain", "")
	base := hiddenRelay(rng, 7)
	mutate := func(key string, v any) string {
		m := map[string]any{}
		for k, x := range base {
			m[k] = x
		}
		if v == nil {
			delete(m, key)
		} else {
			m[key] = v
		}
		return hiddenEncode(m)
	}
	for label, body := range map[string]string{
		"missing band":  mutate("band", nil),
		"missing id":    mutate("id", nil),
		"missing tags":  mutate("tags", nil),
		"id as text":    mutate("id", "7"),
		"fractional id": mutate("id", 1.5),
		"id below min":  mutate("id", 0),
		"negative id":   mutate("id", -3),
		"empty name":    mutate("name", ""),
		"numeric name":  mutate("name", 5),
		"unknown band":  mutate("band", "L"),
		"band case":     mutate("band", "ka"),
		"online text":   mutate("online", "yes"),
		"null online":   strings.Replace(mutate("online", false), "false", "null", 1),
		"extra key":     mutate("crew", "x"),
		"tag type":      mutate("tags", []any{"a", 3}),
		"empty tag":     mutate("tags", []any{"a", ""}),
		"tags not list": mutate("tags", "a"),
		"trailing":      hiddenEncode(base) + " {}",
		"truncated":     hiddenEncode(base)[:10],
		"empty":         "",
		"null":          "null",
		"array":         "[" + hiddenEncode(base) + "]",
	} {
		bad(label, "GET", "/relays/7", 200, js, body)
	}
	broken := append(append([]any{}, list[:3]...), map[string]any{"id": 1, "name": "x", "band": "S"})
	bad("bad list item", "GET", "/relays", 200, js, hiddenEncode(broken))
	bad("list as object", "GET", "/relays", 200, js, hiddenEncode(base))
	bad("error extra key", "GET", "/relays/1", 404, js, ` + "`" + `{"error":"x","code":1}` + "`" + `)
	bad("empty error", "POST", "/relays", 400, js, ` + "`" + `{"error":""}` + "`" + `)
	bad("text type", "GET", "/relays/7", 200, "text/plain", hiddenEncode(base))
	bad("no type", "GET", "/relays/7", 200, "", hiddenEncode(base))
	bad("broken type", "GET", "/relays/7", 200, "application/", hiddenEncode(base))
	bad("204 with body", "DELETE", "/relays/7", 204, "", "x")
	bad("undocumented status", "GET", "/relays/7", 500, js, ` + "`" + `{"error":"x"}` + "`" + `)
	bad("status from other op", "GET", "/relays", 201, js, hiddenEncode(base))
	bad("undocumented method", "PUT", "/relays/7", 200, js, hiddenEncode(base))
	bad("undocumented method 2", "PATCH", "/relays", 200, js, "[]")
	for _, p := range []string{"/relays/", "/relays/7/x", "/beacons", "/", "relays", "/Relays", "//relays"} {
		bad("path "+p, "GET", p, 200, js, "[]")
	}

	gauge := ` + "`" + `{"openapi":"3.0.3","paths":{"/gauges/{deck}/{slot}":{"get":{"responses":{"200":{"content":{"application/json":{"schema":{"$ref":"#/components/schemas/Gauge"}}}}}}}},
"components":{"schemas":{"Gauge":{"type":"object","required":["level","note"],"properties":{"level":{"$ref":"#/components/schemas/Level"},"spare":{},"note":{"type":"null"},"ok":{"type":"boolean"},"marks":{"type":"array"}}},
"Level":{"$ref":"#/components/schemas/Reading"},"Reading":{"type":"number","minimum":-2.5}}}}` + "`" + `
	g, err := LoadSpec([]byte(gauge))
	if err != nil || g == nil {
		t.Fatal("gauge document:", err)
	}
	for _, body := range []string{` + "`" + `{"level":-2.5,"note":null}` + "`" + `, ` + "`" + `{"level":7.25,"note":null,"spare":{"any":[1,"x"]},"ok":false,"extra":1,"marks":[1,"a",null]}` + "`" + `} {
		if err := g.Check("GET", "/gauges/7/b", 200, js, []byte(body)); err != nil {
			t.Errorf("gauge %s rejected: %v", body, err)
		}
	}
	for _, body := range []string{` + "`" + `{"level":-2.6,"note":null}` + "`" + `, ` + "`" + `{"level":"1","note":null}` + "`" + `, ` + "`" + `{"level":1,"note":"x"}` + "`" + `, ` + "`" + `{"level":1}` + "`" + `, ` + "`" + `{"level":1,"note":null,"ok":1}` + "`" + `, ` + "`" + `{"level":1,"note":null,"marks":{}}` + "`" + `} {
		if g.Check("GET", "/gauges/7/b", 200, js, []byte(body)) == nil {
			t.Errorf("gauge %s accepted", body)
		}
	}
	if g.Check("GET", "/gauges/7", 200, js, []byte(` + "`" + `{"level":1,"note":null}` + "`" + `)) == nil || g.Check("GET", "/gauges//b", 200, js, []byte(` + "`" + `{"level":1,"note":null}` + "`" + `)) == nil {
		t.Error("template matched the wrong number of segments or an empty segment")
	}
}

type hiddenRelayView struct {
	ID     int      ` + "`" + `json:"id"` + "`" + `
	Name   string   ` + "`" + `json:"name"` + "`" + `
	Band   string   ` + "`" + `json:"band"` + "`" + `
	Online bool     ` + "`" + `json:"online"` + "`" + `
	Tags   []string ` + "`" + `json:"tags"` + "`" + `
}

// Boundaries: empty list, every RelayInput rule, the 4096-byte cap, optional tags, canonical ids,
// deletes, isolation; every reply must also pass the document.
func TestHiddenRelayAPIHonoursContract(t *testing.T) {
	spec := hiddenSpec(t)
	api := NewRelayAPI()
	rng := rand.New(rand.NewSource(2833))
	do := func(method, path, body string) *httptest.ResponseRecorder {
		t.Helper()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(method, path, strings.NewReader(body))
		api.ServeHTTP(w, r)
		if err := spec.Check(method, r.URL.Path, w.Code, w.Header().Get("Content-Type"), w.Body.Bytes()); err != nil {
			t.Fatalf("%s %s -> %d %q breaks the contract: %v", method, path, w.Code, w.Body.String(), err)
		}
		if w.Body.Len() > 0 {
			if media, _, _ := mime.ParseMediaType(w.Header().Get("Content-Type")); media != "application/json" {
				t.Fatalf("%s %s: content type %q", method, path, w.Header().Get("Content-Type"))
			}
		}
		return w
	}
	expectError := func(w *httptest.ResponseRecorder, status int, text string) {
		t.Helper()
		var e map[string]string
		if w.Code != status || json.Unmarshal(w.Body.Bytes(), &e) != nil || e["error"] != text {
			t.Fatalf("want %d %q, got %d %q", status, text, w.Code, w.Body.String())
		}
	}
	var list []hiddenRelayView
	w := do("GET", "/relays", "")
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &list) != nil || list == nil || len(list) != 0 {
		t.Fatal("empty list must be []", w.Code, w.Body.String())
	}
	padded := ` + "`" + `{"name":"edge","band":"X"}` + "`" + `
	for _, body := range []string{"", "null", "[]", "7", ` + "`" + `{"name":"a"}` + "`" + `, ` + "`" + `{"band":"S"}` + "`" + `, ` + "`" + `{"name":"","band":"S"}` + "`" + `, ` + "`" + `{"name":"a","band":"L"}` + "`" + `,
		` + "`" + `{"name":"a","band":"S","crew":1}` + "`" + `, ` + "`" + `{"name":"a","band":"S","id":4}` + "`" + `, ` + "`" + `{"name":"a","band":"S","tags":[""]}` + "`" + `, ` + "`" + `{"name":"a","band":"S","tags":"x"}` + "`" + `,
		` + "`" + `{"name":"a","band":"S","tags":null}` + "`" + `, ` + "`" + `{"name":"a","band":"S"} {}` + "`" + `, ` + "`" + `{"name":5,"band":"S"}` + "`" + `, padded + strings.Repeat(" ", 4097-len(padded))} {
		expectError(do("POST", "/relays", body), 400, "invalid relay")
	}
	var want []hiddenRelayView
	n := 8 + rng.Intn(8)
	for i := 1; i <= n; i++ {
		in := map[string]any{"name": hiddenName(rng), "band": []string{"S", "X", "Ka"}[rng.Intn(3)]}
		tags := []string{}
		if rng.Intn(3) > 0 {
			for j := rng.Intn(3); j >= 0; j-- {
				tags = append(tags, fmt.Sprintf("tag-%d", rng.Intn(50)))
			}
			in["tags"] = tags
		}
		body := hiddenEncode(in)
		if i == 1 {
			body += strings.Repeat(" ", 4096-len(body))
		}
		w := do("POST", "/relays", body)
		var got hiddenRelayView
		wantRelay := hiddenRelayView{i, in["name"].(string), in["band"].(string), false, tags}
		if w.Code != 201 || json.Unmarshal(w.Body.Bytes(), &got) != nil || !reflect.DeepEqual(got, wantRelay) || w.Header().Get("Location") != fmt.Sprintf("/relays/%d", i) {
			t.Fatalf("create %d: %d %q %v", i, w.Code, w.Body.String(), w.Header())
		}
		want = append(want, wantRelay)
	}
	for _, r := range want {
		var got hiddenRelayView
		if w := do("GET", fmt.Sprintf("/relays/%d", r.ID), ""); w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &got) != nil || !reflect.DeepEqual(got, r) {
			t.Fatal("read", r.ID, w.Body.String())
		}
	}
	for _, p := range []string{"/relays/0", "/relays/01", "/relays/+2", "/relays/abc", "/relays/-1", "/relays/999999999999999999999", fmt.Sprintf("/relays/%d", n+1)} {
		expectError(do("GET", p, ""), 404, "unknown relay")
		expectError(do("DELETE", p, ""), 404, "unknown relay")
	}
	gone := 1 + rng.Intn(n)
	if w := do("DELETE", fmt.Sprintf("/relays/%d", gone), ""); w.Code != 204 || w.Body.Len() != 0 {
		t.Fatal("delete", w.Code, w.Body.String())
	}
	expectError(do("GET", fmt.Sprintf("/relays/%d", gone), ""), 404, "unknown relay")
	expectError(do("DELETE", fmt.Sprintf("/relays/%d", gone), ""), 404, "unknown relay")
	want = append(want[:gone-1], want[gone:]...)
	w = do("GET", "/relays?band=S", "")
	if json.Unmarshal(w.Body.Bytes(), &list) != nil || !reflect.DeepEqual(list, want) {
		t.Fatal("list after delete", w.Body.String())
	}
	var created hiddenRelayView
	if w := do("POST", "/relays", ` + "`" + `{"name":"next","band":"S"}` + "`" + `); json.Unmarshal(w.Body.Bytes(), &created) != nil || created.ID != n+1 {
		t.Fatal("ids must keep increasing after a delete", w.Body.String())
	}
	other := httptest.NewRecorder()
	NewRelayAPI().ServeHTTP(other, httptest.NewRequest("POST", "/relays", strings.NewReader(` + "`" + `{"name":"solo","band":"X"}` + "`" + `)))
	if json.Unmarshal(other.Body.Bytes(), &created) != nil || created.ID != 1 {
		t.Fatal("two APIs share a store", other.Body.String())
	}
}
`

func init() {
	add(&Level{ID: "28-4-graceful", Chapter: 28, Kind: "lab", Title: "إغلاق مرتب", Goal: "اقفل بوابة القاعدة من غير ما تقطع طلب شغال: ارفض الجديد، صفّي الجاري جوه مهلة، ونضّف مرة واحدة",
		Requires: []string{"sel:http.Server", "sel:context.Context"}, Bridge: []int{60}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch284Starter}},
		Lesson: `# إغلاق مرتب لبوابة القاعدة

## قبل ما تكتب
- لو جالك أمر إيقاف وطلبات الطاقم لسه شغالة، مين يستنى مين؟
- لو طلب مش راضي يخلص، هتستنى لحد إمتى، وبعدها هتعمل إيه؟

## من الإشارة لـ context
البرنامج الحقيقي بيقف لما نظام التشغيل يبعتله إشارة، زي Ctrl+C. ~signal.NotifyContext~ من ~os/signal~ بترجع context بيتلغي أول ما الإشارة توصل، ودالة ~stop~ بترجّع السلوك العادي. ~os.Interrupt~ هي إشارة Ctrl+C.
~~~go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()
err := gate.Run(ctx, ln) // Run متعرفش حاجة عن الإشارات
~~~
~Run~ بتاخد ~context.Context~ بس، فالتست يقدر يوقفها بـ ~cancel()~ من غير إشارات حقيقية. ~main~ هي الحتة الوحيدة اللي تعرف ~os/signal~.

## Serve على listener محقون
~net.Listen("tcp", "127.0.0.1:0")~ بيفتح بورت فاضي على الجهاز نفسه بس، و ~ln.Addr().String()~ بيقول العنوان اللي اتختار. ~srv.Serve(ln)~ بتقبل اتصالات من الـ listener ده لحد ما السيرفر يتقفل، ومبترجعش قبلها، فشغّلها في goroutine.
~~~go
ln, err := net.Listen("tcp", "127.0.0.1:0")
srv := &http.Server{Handler: h, ReadHeaderTimeout: 2 * time.Second}
served := make(chan error, 1)
go func() { served <- srv.Serve(ln) }()
~~~
لما انت تقفل السيرفر، ~Serve~ بترجع ~http.ErrServerClosed~، ودي مش مشكلة. أي خطأ تاني، زي listener مقفول من الأول، فشل حقيقي.

## Shutdown و Close
~srv.Shutdown(ctx)~ بتقفل الـ listeners فوراً عشان أي اتصال جديد يترفض، وبتقفل الاتصالات الخاملة، وبعدين بتستنى الطلبات الجارية لحد ما تخلص. لو الـ ctx اللي ادّيتهولها اتلغى قبل كده، بترجع خطأه وبتسيب الطلبات شغالة. ~srv.Close()~ بتقطع كل الاتصالات فوراً من غير استنى.
~~~go
if err := srv.Shutdown(drainCtx); err != nil {
 srv.Close() // المهلة خلصت: اقطع الباقي
}
~~~
~Serve~ بترجع ~ErrServerClosed~ في أول لحظة من ~Shutdown~، قبل ما الطلبات الجارية تخلص. فاللي بيقول إن الإغلاق خلص هو رجوع ~Shutdown~، مش رجوع ~Serve~.

## مهلة جديدة من ساعة محقونة
الـ context اللي جالك منه أمر الإيقاف متلغي خلاص؛ لو ادّيته لـ ~Shutdown~ هترجع فوراً كأن المهلة خلصت. اعمل context جديد من ~context.Background()~ والغيه لما ساعة المهلة ترن. ~clock.After(d)~ من ~labkit.Clock~ شغالة زي ~time.After~، والتست بيحرّكها بـ ~Advance~.
~~~go
drainCtx, cancel := context.WithCancel(context.Background())
defer cancel()
finished := make(chan struct{})
defer close(finished)
go func() {
 select {
 case <-clock.After(grace): cancel() // المهلة خلصت
 case <-finished:                    // Shutdown رجعت لوحدها
 }
}()
~~~
الـ ~finished~ بيخلّي goroutine المؤقت يخرج بعد ما ~Run~ ترجع بدل ما يفضل مستني.

## التنضيف مرة واحدة بعد السكوت
الموارد زي قاعدة البيانات أو ملف السجل بتتقفل بعد ما السيرفر يسكت خالص، مش قبل؛ الطلب الجاري ممكن لسه بيستخدمها. ~release~ تتنادى مرة واحدة بالظبط في أي طريق خروج. ~errors.Join~ من فصل الأخطاء بيجمع خطأ الإغلاق وخطأ التنضيف، و ~errors.Is~ بتلاقي أي واحد فيهم؛ ~errors.Join(nil, nil)~ بترجع ~nil~.
~~~go
<-served // Serve خرجت
return errors.Join(result, release())
~~~
~Run~ نفسها تشتغل مرة واحدة. ~atomic.Bool~ زي ~atomic.Int64~ اللي استخدمته قبل كده، بس لقيمة ~true~ أو ~false~ من غير قفل. ~CompareAndSwap(old, new)~ بتحط ~new~ بس لو القيمة الحالية ~old~، وبترجع ~true~ لو غيّرت؛ فأول نداء بس هو اللي بيكسب، حتى لو اتنين ندهوا في نفس اللحظة.
~~~go
var used atomic.Bool
first := used.CompareAndSwap(false, true)  // true
second := used.CompareAndSwap(false, true) // false: اتغيرت خلاص
~~~

## النجمة التالتة: خبّر الـ handlers
الطلب الطويل، زي بث قراءات، يقدر يخلص بدري لو عرف إن البوابة بتقفل. ~srv.RegisterOnShutdown(f)~ بتسجّل دالة بتشتغل أول ما ~Shutdown~ تبدأ. قفل channel بيصحّي كل اللي مستنيين عليه مرة واحدة، و ~sync.Once~ بتضمن إنه يتقفل مرة بس.
~~~go
srv.RegisterOnShutdown(func() { once.Do(func() { close(drain) }) })
// جوه handler طويل:
select { case <-drain: w.Write([]byte("bye")) case <-r.Context().Done(): }
~~~

## أخطاء شائعة
~~~go
<-ctx.Done(); srv.Shutdown(ctx)      // ctx متلغي: مفيش مهلة خالص
go srv.Shutdown(c); return <-served  // رجعت قبل ما الطلبات تخلص
release(); srv.Shutdown(c)           // قفلت المورد والطلب لسه بيستخدمه
<-time.After(grace)                  // ساعة حقيقية: التست ميقدرش يحرّكها
~~~

## المهمة
نفّذ ~NewGate(h http.Handler, clock labkit.Clock, grace time.Duration, release func() error) (*Gate, error)~: handler أو ساعة أو ~release~ بـ ~nil~، أو ~grace~ مش موجبة، ترجع ~nil,error~.

نفّذ ~(g *Gate) Run(ctx context.Context, ln net.Listener) error~: يخدم ~h~ على ~ln~ بـ ~http.Server~ لحد ما ~ctx~ يتلغي. ساعتها يبدأ ~Shutdown~ فوراً، فالاتصالات الجديدة تترفض، ويستنى الطلبات الجارية. مهلة الاستنى بتبدأ بعد الإيقاف بنداء واحد لـ ~clock.After(grace)~؛ لو الطلبات خلصت قبلها ترجع ~nil~. لو المهلة خلصت يقطع الباقي بـ ~Close~ ويرجع خطأ يطابق ~ErrForced~ بـ ~errors.Is~.

لو ~Serve~ فشلت قبل الإيقاف، زي listener مقفول، ~Run~ ترجع خطأها من غير ما تستنى ~ctx~. لو ~ctx~ متلغي قبل النداء، ترجع ~nil~ والـ listener مقفول. ~Run~ مترجعش غير بعد ما goroutine الـ ~Serve~ يخرج. في كل الطرق دي ~release~ تتنادى مرة واحدة بعد ما الخدمة تسكت، وخطأها يتجمع في الخطأ الراجع بـ ~errors.Join~. نداء ~Run~ تاني، حتى في نفس اللحظة، يرجع ~ErrGateUsed~ فوراً من غير ما يلمس ~release~. ~http.ErrServerClosed~ مش خطأ. سيب ~main~ فاضية أو اربطها بـ ~signal.NotifyContext~؛ التستات مش بتناديها.

النجمة التالتة: ~(g *Gate) Draining() <-chan struct{}~ ترجع نفس الـ channel دايماً، مفتوح لحد ما الإيقاف يبدأ وبعدها مقفول، عشان الـ handler الطويل يخلص بدري والإغلاق يرجع ~nil~ من غير ما المهلة تخلص.`,
		Hints: []string{"شغّل srv.Serve(ln) في goroutine بيبعت نتيجته على channel، واستنى بـ select على ctx.Done() وعلى الـ channel ده.", "بعد الإيقاف اعمل context جديد من context.Background() والغيه من goroutine فيه select على clock.After(grace) وعلى channel بيتقفل لما Shutdown ترجع؛ لو Shutdown رجعت خطأ نادِ Close ورجّع ErrForced.", "استنى رجوع Serve، وبعدين errors.Join(نتيجتك، release())؛ احمي Run بـ atomic.Bool.CompareAndSwap، وللنجمة سجّل قفل الـ channel بـ RegisterOnShutdown وsync.Once."},
		En: &LevelText{Title: "Graceful Shutdown", Goal: "Close the base gateway without cutting a running request: refuse new ones, drain the rest within a grace period, and clean up once", Files: []LabFile{{Name: "main.go", Starter: ch284Starter}},
			Lesson: `# A Graceful Shutdown for the Base Gateway

## Before you write
- If a stop order arrives while crew requests are still running, who waits for whom?
- If a request refuses to finish, how long do you wait, and then what?

## From signal to context
A real program stops when the operating system sends it a signal, such as Ctrl+C. ~signal.NotifyContext~ from ~os/signal~ returns a context canceled as soon as the signal arrives, plus a ~stop~ function that restores normal behavior. ~os.Interrupt~ is the Ctrl+C signal.
~~~go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
defer stop()
err := gate.Run(ctx, ln) // Run knows nothing about signals
~~~
~Run~ takes only a ~context.Context~, so a test can stop it with ~cancel()~ and no real signals. ~main~ is the only piece that knows ~os/signal~.

## Serve on an injected listener
~net.Listen("tcp", "127.0.0.1:0")~ opens a free port on this machine only, and ~ln.Addr().String()~ tells you the chosen address. ~srv.Serve(ln)~ accepts connections from that listener until the server is closed and does not return before, so run it in a goroutine.
~~~go
ln, err := net.Listen("tcp", "127.0.0.1:0")
srv := &http.Server{Handler: h, ReadHeaderTimeout: 2 * time.Second}
served := make(chan error, 1)
go func() { served <- srv.Serve(ln) }()
~~~
When you close the server yourself, ~Serve~ returns ~http.ErrServerClosed~, which is fine. Any other error, such as an already closed listener, is a real failure.

## Shutdown and Close
~srv.Shutdown(ctx)~ closes the listeners at once so new connections are refused, closes idle connections, then waits for running requests to finish. If the ctx you gave it is canceled first, it returns that error and leaves the requests running. ~srv.Close()~ cuts every connection immediately without waiting.
~~~go
if err := srv.Shutdown(drainCtx); err != nil {
 srv.Close() // Grace is over: cut the rest
}
~~~
~Serve~ returns ~ErrServerClosed~ the moment ~Shutdown~ starts, before running requests finish. So the end of shutdown is marked by ~Shutdown~ returning, not by ~Serve~ returning.

## A fresh deadline from an injected clock
The context that delivered the stop order is already canceled; give it to ~Shutdown~ and it returns at once as if grace had run out. Make a new context from ~context.Background()~ and cancel it when the grace clock fires. ~clock.After(d)~ from ~labkit.Clock~ works like ~time.After~, and the test moves it with ~Advance~.
~~~go
drainCtx, cancel := context.WithCancel(context.Background())
defer cancel()
finished := make(chan struct{})
defer close(finished)
go func() {
 select {
 case <-clock.After(grace): cancel() // Grace is over
 case <-finished:                    // Shutdown returned by itself
 }
}()
~~~
~finished~ lets the timer goroutine exit after ~Run~ returns instead of waiting forever.

## Clean up once, after silence
Resources such as a database or a log file are closed after the server is completely quiet, not before; a running request may still use them. ~release~ is called exactly once on every exit path. ~errors.Join~ from the errors chapter combines the shutdown error and the cleanup error, and ~errors.Is~ finds either; ~errors.Join(nil, nil)~ returns ~nil~.
~~~go
<-served // Serve has exited
return errors.Join(result, release())
~~~
~Run~ itself runs once. ~atomic.Bool~ is like the ~atomic.Int64~ you used before, but for a lock-free ~true~ or ~false~. ~CompareAndSwap(old, new)~ stores ~new~ only if the current value is ~old~, and returns ~true~ when it changed it; so only the first call wins, even when two calls race.
~~~go
var used atomic.Bool
first := used.CompareAndSwap(false, true)  // true
second := used.CompareAndSwap(false, true) // false: already changed
~~~

## Third star: tell the handlers
A long request, such as a stream of readings, can finish early if it knows the gateway is closing. ~srv.RegisterOnShutdown(f)~ registers a function that runs as soon as ~Shutdown~ starts. Closing a channel wakes everyone waiting on it at once, and ~sync.Once~ makes sure it closes only once.
~~~go
srv.RegisterOnShutdown(func() { once.Do(func() { close(drain) }) })
// Inside a long handler:
select { case <-drain: w.Write([]byte("bye")) case <-r.Context().Done(): }
~~~

## Common mistakes
~~~go
<-ctx.Done(); srv.Shutdown(ctx)      // ctx is canceled: no grace at all
go srv.Shutdown(c); return <-served  // Returned before requests finished
release(); srv.Shutdown(c)           // Closed the resource while a request still uses it
<-time.After(grace)                  // A real clock: the test cannot move it
~~~

## Your mission
Implement ~NewGate(h http.Handler, clock labkit.Clock, grace time.Duration, release func() error) (*Gate, error)~: a ~nil~ handler, clock, or ~release~, or a non-positive ~grace~, returns ~nil,error~.

Implement ~(g *Gate) Run(ctx context.Context, ln net.Listener) error~: serve ~h~ on ~ln~ with an ~http.Server~ until ~ctx~ is canceled. Then start ~Shutdown~ at once, so new connections are refused, and wait for running requests. The wait starts after the stop with one call to ~clock.After(grace)~; if requests finish first, return ~nil~. If grace runs out, cut the rest with ~Close~ and return an error matching ~ErrForced~ with ~errors.Is~.

If ~Serve~ fails before a stop, such as on a closed listener, ~Run~ returns its error without waiting for ~ctx~. If ~ctx~ is canceled before the call, return ~nil~ with the listener closed. ~Run~ returns only after the ~Serve~ goroutine has exited. On every one of these paths ~release~ is called once after serving is quiet, and its error is joined into the returned error with ~errors.Join~. A second ~Run~ call, even a simultaneous one, returns ~ErrGateUsed~ at once without touching ~release~. ~http.ErrServerClosed~ is not an error. Leave ~main~ empty or wire it to ~signal.NotifyContext~; the tests never call it.

Third star: ~(g *Gate) Draining() <-chan struct{}~ always returns the same channel, open until the stop begins and closed after, so a long handler can finish early and shutdown returns ~nil~ without grace running out.`,
			Hints: []string{"Run srv.Serve(ln) in a goroutine that sends its result on a channel, and wait with select on ctx.Done() and that channel.", "After the stop, make a fresh context from context.Background() and cancel it from a goroutine selecting on clock.After(grace) and on a channel closed when Shutdown returns; if Shutdown returns an error, call Close and return ErrForced.", "Wait for Serve to return, then errors.Join(your result, release()); guard Run with atomic.Bool.CompareAndSwap, and for the star register closing the channel with RegisterOnShutdown and sync.Once."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch284Hidden},
	})
}

const ch284Starter = `package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"orbital/colony/labkit"
)

var (
	ErrForced   = errors.New("drain deadline passed")
	ErrGateUsed = errors.New("gate already ran")
)

type Gate struct{}

func NewGate(h http.Handler, clock labkit.Clock, grace time.Duration, release func() error) (*Gate, error) {
	return nil, errors.New("TODO")
}
func (g *Gate) Run(ctx context.Context, ln net.Listener) error { return errors.New("TODO") }
func (g *Gate) Draining() <-chan struct{}                      { return nil }
func main()                                                    {}
`

const ch284Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"orbital/colony/labkit"
)

// hiddenGateClock reports every After call so a test advances only after the timer exists.
type hiddenGateClock struct {
	*labkit.FakeClock
	armed chan time.Duration
}

func (c hiddenGateClock) After(d time.Duration) <-chan time.Time {
	ch := c.FakeClock.After(d)
	c.armed <- d
	return ch
}

func newHiddenGateClock() hiddenGateClock {
	return hiddenGateClock{labkit.NewFakeClock(time.Unix(1_700_000_000, 0)), make(chan time.Duration, 16)}
}

func hiddenGateListen(t *testing.T) net.Listener {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	return ln
}

func hiddenGateRefuses(addr string) bool {
	c, err := net.DialTimeout("tcp", addr, 200*time.Millisecond)
	if err != nil {
		return true
	}
	c.Close()
	return false
}

func hiddenGateWait(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func hiddenGateGet(addr, path string) <-chan string {
	out := make(chan string, 1)
	go func() {
		client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
		resp, err := client.Get("http://" + addr + path)
		if err != nil {
			out <- "error"
			return
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		out <- fmt.Sprintf("%d %s", resp.StatusCode, body)
	}()
	return out
}

func hiddenGateRecv[T any](t *testing.T, what string, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out waiting for %s", what)
	}
	var zero T
	return zero
}

func hiddenGateStill(t *testing.T, what string, done <-chan error) {
	t.Helper()
	select {
	case err := <-done:
		t.Fatalf("Run returned %v %s", err, what)
	case <-time.After(100 * time.Millisecond):
	}
}

// Boundaries: a finished request before the stop, an in-flight request across it,
// new connections refused during the drain, release only after silence.
func TestHiddenGateDrainsInFlight(t *testing.T) {
	rng := rand.New(rand.NewSource(2841))
	clock := newHiddenGateClock()
	grace := time.Duration(1+rng.Intn(30)) * time.Second
	entered, unblock := make(chan struct{}), make(chan struct{})
	var finished atomic.Bool
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/slow" {
			close(entered)
			<-unblock
			finished.Store(true)
		}
		io.WriteString(w, "done "+r.URL.Path)
	})
	var releases atomic.Int32
	var early atomic.Bool
	g, err := NewGate(h, clock, grace, func() error {
		early.Store(!finished.Load())
		releases.Add(1)
		return nil
	})
	if err != nil || g == nil {
		t.Fatal("NewGate:", err)
	}
	ln := hiddenGateListen(t)
	addr := ln.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- g.Run(ctx, ln) }()
	if got := hiddenGateRecv(t, "first reply", hiddenGateGet(addr, "/fast")); got != "200 done /fast" {
		t.Fatal("serving before the stop:", got)
	}
	slow := hiddenGateGet(addr, "/slow")
	hiddenGateRecv(t, "slow request to start", entered)
	cancel()
	hiddenGateWait(t, "new connections to be refused", func() bool { return hiddenGateRefuses(addr) })
	hiddenGateStill(t, "while a request was still running", done)
	if releases.Load() != 0 {
		t.Fatal("release ran while a request was still running")
	}
	close(unblock)
	if got := hiddenGateRecv(t, "slow reply", slow); got != "200 done /slow" {
		t.Fatal("the in-flight request was cut:", got)
	}
	if err := hiddenGateRecv(t, "Run to return", done); err != nil {
		t.Fatal("clean drain returned", err)
	}
	if releases.Load() != 1 || early.Load() {
		t.Fatal("release must run exactly once, after the last request", releases.Load(), early.Load())
	}
}

// Boundaries: the timer starts after the stop with exactly grace, one tick short of grace keeps waiting,
// grace itself forces Close, the stuck client is cut, release still runs once.
func TestHiddenGateForcesAfterGrace(t *testing.T) {
	rng := rand.New(rand.NewSource(2842))
	clock := newHiddenGateClock()
	grace := time.Duration(1+rng.Intn(60)) * time.Second
	entered, unblock := make(chan struct{}), make(chan struct{})
	defer close(unblock)
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-unblock
	})
	var releases atomic.Int32
	g, err := NewGate(h, clock, grace, func() error { releases.Add(1); return nil })
	if err != nil {
		t.Fatal(err)
	}
	ln := hiddenGateListen(t)
	addr := ln.Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- g.Run(ctx, ln) }()
	stuck := hiddenGateGet(addr, "/stuck")
	hiddenGateRecv(t, "stuck request to start", entered)
	select {
	case d := <-clock.armed:
		t.Fatal("grace timer started before the stop:", d)
	case <-time.After(50 * time.Millisecond):
	}
	cancel()
	if d := hiddenGateRecv(t, "clock.After(grace) after the stop", (<-chan time.Duration)(clock.armed)); d != grace {
		t.Fatalf("grace timer %v, want %v", d, grace)
	}
	hiddenGateWait(t, "new connections to be refused", func() bool { return hiddenGateRefuses(addr) })
	clock.Advance(grace - time.Nanosecond)
	hiddenGateStill(t, "before the grace period ended", done)
	clock.Advance(time.Nanosecond)
	if err := hiddenGateRecv(t, "Run to force the close", done); !errors.Is(err, ErrForced) {
		t.Fatal("want ErrForced, got", err)
	}
	if got := hiddenGateRecv(t, "stuck client", stuck); got != "error" {
		t.Fatal("forced close must cut the stuck connection:", got)
	}
	if releases.Load() != 1 {
		t.Fatal("release calls:", releases.Load())
	}
}

func hiddenGateRun(t *testing.T, g *Gate, ctx context.Context, ln net.Listener) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- g.Run(ctx, ln) }()
	return hiddenGateRecv(t, "Run to return", done)
}

// Boundaries: invalid construction, a pre-canceled context, a failing listener,
// release errors, and a second or simultaneous Run.
func TestHiddenGateLifecycleEdges(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	clock := newHiddenGateClock()
	ok := func() error { return nil }
	for name, build := range map[string]func() (*Gate, error){
		"nil handler":    func() (*Gate, error) { return NewGate(nil, clock, time.Second, ok) },
		"nil clock":      func() (*Gate, error) { return NewGate(h, nil, time.Second, ok) },
		"zero grace":     func() (*Gate, error) { return NewGate(h, clock, 0, ok) },
		"negative grace": func() (*Gate, error) { return NewGate(h, clock, -time.Second, ok) },
		"nil release":    func() (*Gate, error) { return NewGate(h, clock, time.Second, nil) },
	} {
		if g, err := build(); g != nil || err == nil {
			t.Errorf("%s accepted", name)
		}
	}
	var releases atomic.Int32
	count := func() error { releases.Add(1); return nil }
	g, err := NewGate(h, clock, time.Second, count)
	if err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	ln := hiddenGateListen(t)
	if err := hiddenGateRun(t, g, canceled, ln); err != nil {
		t.Fatal("pre-canceled Run:", err)
	}
	if !hiddenGateRefuses(ln.Addr().String()) || releases.Load() != 1 {
		t.Fatal("pre-canceled Run must close the listener and release once", releases.Load())
	}
	if err := hiddenGateRun(t, g, context.Background(), hiddenGateListen(t)); !errors.Is(err, ErrGateUsed) || releases.Load() != 1 {
		t.Fatal("second Run:", err, releases.Load())
	}

	sentinel := errors.New("store close failed")
	broken := hiddenGateListen(t)
	broken.Close()
	var failed atomic.Int32
	g2, _ := NewGate(h, clock, time.Second, func() error { failed.Add(1); return sentinel })
	err = hiddenGateRun(t, g2, context.Background(), broken)
	if err == nil || errors.Is(err, ErrForced) || !errors.Is(err, net.ErrClosed) || !errors.Is(err, sentinel) || failed.Load() != 1 {
		t.Fatal("failing listener:", err, failed.Load())
	}
	g3, _ := NewGate(h, clock, time.Second, func() error { return sentinel })
	if err := hiddenGateRun(t, g3, canceled, hiddenGateListen(t)); !errors.Is(err, sentinel) || errors.Is(err, ErrForced) {
		t.Fatal("release error after a clean stop:", err)
	}

	var raced atomic.Int32
	g4, _ := NewGate(h, clock, time.Second, func() error { raced.Add(1); return nil })
	ctx, stop := context.WithCancel(context.Background())
	results := make(chan error, 2)
	start := make(chan struct{})
	for _, l := range []net.Listener{hiddenGateListen(t), hiddenGateListen(t)} {
		go func(l net.Listener) {
			<-start
			results <- g4.Run(ctx, l)
		}(l)
	}
	close(start)
	first := hiddenGateRecv(t, "the losing Run", (<-chan error)(results))
	if !errors.Is(first, ErrGateUsed) {
		t.Fatal("simultaneous Run must return ErrGateUsed at once, got", first)
	}
	stop()
	if err := hiddenGateRecv(t, "the winning Run", (<-chan error)(results)); err != nil || raced.Load() != 1 {
		t.Fatal("winning Run:", err, raced.Load())
	}
}

func TestHiddenBonusGateDrainingSignal(t *testing.T) {
	clock := newHiddenGateClock()
	var g *Gate
	entered := make(chan struct{})
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-g.Draining()
		io.WriteString(w, "bye")
	})
	g, err := NewGate(h, clock, time.Minute, func() error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if g.Draining() == nil || g.Draining() != g.Draining() {
		t.Fatal("Draining must always return the same non-nil channel")
	}
	select {
	case <-g.Draining():
		t.Fatal("Draining closed before the stop")
	default:
	}
	ln := hiddenGateListen(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- g.Run(ctx, ln) }()
	reply := hiddenGateGet(ln.Addr().String(), "/stream")
	hiddenGateRecv(t, "stream to start", entered)
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal("drain with an early finish returned", err)
		}
	case <-time.After(3 * time.Second):
		clock.Advance(time.Minute)
		<-done
		t.Fatal("the long handler was never told to wrap up")
	}
	if got := hiddenGateRecv(t, "stream reply", reply); got != "200 bye" {
		t.Fatal("stream reply:", got)
	}
}
`
