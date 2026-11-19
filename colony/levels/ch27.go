package levels

func init() {
	add(&Level{LessonRevision: 1, Steps: ch27Lesson271.arabic, ID: "27-1-signed-pass", Chapter: 27, Kind: "lab", Title: "تصريح موقّع", Goal: "اطلع تصاريح دخول قصيرة العمر بعد فحص كود الطاقم بـ bcrypt",
		Requires: []string{"sel:bcrypt.CompareHashAndPassword", "sel:jwt.NewWithClaims"}, Bridge: []int{50}, Setup: seededLabRoom, Race: false,
		Deps:  []string{"golang.org/x/crypto", "github.com/golang-jwt/jwt/v5"},
		Files: []LabFile{{Name: "main.go", Starter: ch271Starter}},
		Lesson: `# تصريح موقّع لبوابة المساكن

## قبل ما تكتب
- لو ملف أكواد الطاقم اتسرق، إيه اللي يمنع الحرامي يدخل بيه؟
- السيرفر يعرف إزاي إن التصريح اللي جاله هو اللي طلّعه بنفسه ولسه صالح؟

## الكود مش بيتخزن
bcrypt دالة hash بطيئة عن قصد وفيها salt عشوائي، فالـ hash المسروق صعب يتعكس. استورد ~golang.org/x/crypto/bcrypt~. ~bcrypt.GenerateFromPassword~ بتاخد الكود وتكلفة وترجع hash؛ ~bcrypt.MinCost~ للتستات بس، و ~bcrypt.DefaultCost~ للتشغيل.
~~~go
hash, err := bcrypt.GenerateFromPassword([]byte("nebula-7"), bcrypt.DefaultCost)
if err != nil { return err }
err = bcrypt.CompareHashAndPassword(hash, []byte("nebula-7")) // nil
~~~
~bcrypt.CompareHashAndPassword~ ترجع ~nil~ لو الكود مطابق. الكود الغلط بيرجع ~bcrypt.ErrMismatchedHashAndPassword~؛ افحصه بـ ~errors.Is~. أي خطأ تاني معناه hash متخزن بايظ، يعني عطل عندك مش غلطة العميل.
~~~go
if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { deny(w); return }
~~~
كل حروف الكود مهمة، حتى المسافات؛ متعملش trim. الاسم المجهول يعدّي على مقارنة مع hash وهمي متولّد مرة واحدة، عشان وقت الرد ميفضحش مين موجود، ويطلع نفس رد الكود الغلط.

## السر من البيئة
~os.Getenv~ من فصل سطر الأوامر بتقرا متغير بيئة. هنا بنحقن دالة بنفس الشكل عشان التست ميغيرش بيئة العملية، و ~main~ الحقيقية تمرر ~os.Getenv~.
~~~go
secret, err := LoadSecret(os.Getenv)
~~~
السر مكتوب base64 لأنه بايتات عشوائية مش نص. ~encoding/base64~ فيها ~base64.StdEncoding~؛ ~EncodeToString~ تحوّل بايتات لنص، و ~DecodeString~ ترجعها أو خطأ لو النص مش base64.
~~~go
text := base64.StdEncoding.EncodeToString([]byte("abc")) // "YWJj"
key, err := base64.StdEncoding.DecodeString(text)
~~~
مفتاح HMAC أقصر من 32 بايت ضعيف؛ ارفضه. ممنوع سر افتراضي مكتوب في الكود.

## JWT: ادعاءات وتوقيع
JWT تلات أجزاء base64url بينهم نقطة: header فيه الخوارزمية، و payload فيه الادعاءات (claims)، وتوقيع. أي حد يقدر يقرا الـ payload؛ التوقيع بس هو اللي بيمنع التعديل. استورد ~github.com/golang-jwt/jwt/v5~؛ اسم الباكدج ~jwt~.

~jwt.RegisteredClaims~ فيها ~Issuer~ (iss) و ~Subject~ (sub) و ~IssuedAt~ (iat) و ~ExpiresAt~ (exp). ~jwt.NewNumericDate~ بتحوّل ~time.Time~ لثواني كاملة. ~jwt.NewWithClaims~ مع ~jwt.SigningMethodHS256~ تبني التوكن، و ~SignedString~ توقّعه بالسر.
~~~go
now := clock.Now()
claims := jwt.RegisteredClaims{Issuer: "hangar", Subject: "rigel",
 IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute))}
pass, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
~~~
العمر القصير بيقلل ضرر التصريح المسروق، لأن مفيش طريقة نلغيه قبل ما يخلص.

## التحقق قبل الثقة
~jwt.ParseWithClaims~ بتفك التوكن في claims بتاعتك، وتنادي keyfunc عشان تجيب المفتاح، وبعدين تتحقق من التوقيع والوقت. الخيارات بتشد التحقق:
~~~go
var c jwt.RegisteredClaims
_, err := jwt.ParseWithClaims(raw, &c, func(*jwt.Token) (any, error) { return secret, nil },
 jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("hangar"),
 jwt.WithTimeFunc(clock.Now), jwt.WithExpirationRequired(), jwt.WithIssuedAt())
~~~
~jwt.WithValidMethods~ بترفض ~none~ و HS512 وأي خوارزمية تانية حتى لو التوقيع سليم. ~jwt.WithIssuer~ بتطابق iss. ~jwt.WithTimeFunc~ بتحقن الساعة؛ ~clock.Now~ هنا method value. ~jwt.WithExpirationRequired~ بترفض توكن من غير exp، و ~jwt.WithIssuedAt~ بترفض iat في المستقبل.

التوكن صالح لحد اللحظة اللي قبل exp؛ عند exp نفسها خلاص انتهى. مفيش سماحية وقت ضمنية.

## middleware المصادقة
الهيدر ~Authorization~ شكله ~Bearer <pass>~ بالظبط، بمسافة واحدة. ~strings.HasPrefix~ و ~strings.TrimPrefix~ يفصلوا التوكن. الفشل 401 مع ~WWW-Authenticate: Bearer~ عشان العميل يعرف نوع التصريح المطلوب، ومن غير ما تنادي ~next~.
~~~go
if !strings.HasPrefix(h, "Bearer ") { deny(w); return }
r = r.WithContext(context.WithValue(r.Context(), crewKey{}, c.Subject))
next.ServeHTTP(w, r)
~~~
النجاح يحط اسم الطاقم في السياق بمفتاح خاص زي 26-3، والـ handler الداخلي يقراه من غير ما يلمس التوكن.

## أخطاء شائعة
~~~go
if string(stored) == code { open() } // الكود متخزن نص عادي
jwt.Parse(raw, keyFunc) // من غير WithValidMethods: توكن HS512 بنفس السر يعدّي
var secret = []byte("dev-secret") // سر مكتوب في الكود
~~~

## المهمة
نفّذ ~LoadSecret(getenv func(string) string) ([]byte,error)~: يقرا ~ORBITAL_DESK_SECRET~ بس، يفكه بـ ~base64.StdEncoding~، ويرفض الفاضي أو غير الصالح أو الأقصر من 32 بايت. نفّذ ~NewPassDesk(secret []byte,crew map[string][]byte,clock labkit.Clock) (*PassDesk,error)~: سر أقصر من 32 بايت أو ساعة ~nil~ يرجع ~nil,error~. ~crew~ اسم ← bcrypt hash. خد نسخة من السر والخريطة؛ تعديلهم بعد الإنشاء ميأثرش.

~(d *PassDesk) Login(w,r)~ handler لأي مسار. غير ~POST~ يرجع 405 و ~Allow: POST~ و ~use POST\n~. الجسم حده 4096 بايت شاملة المسافات، وفيه كائن JSON واحد ~{"crew":"...","code":"..."}~. JSON بايظ أو أكبر من الحد أو اسم أو كود فاضي: 400 و ~bad login\n~. اسم مجهول أو كود غلط: 401 و ~denied\n~ بنفس الرد. النجاح: 200 و ~application/json~ و ~{"pass":"<jwt>"}~؛ HS256 بالسر، و ~iss=orbital-desk~ و ~sub~ الاسم و ~iat~ وقت الساعة و ~exp~ بعد 10 دقايق بالظبط.

~(d *PassDesk) Require(next http.Handler) http.Handler~ يقبل بس ~Authorization: Bearer <pass>~ لتوكن HS256 بنفس السر، و ~iss~ صحيح، و ~sub~ غير فاضي، و ~exp~ موجود وبعد الوقت الحالي، و ~iat~ مش في المستقبل. أي فشل: 401 و ~WWW-Authenticate: Bearer~ و ~pass required\n~ من غير ~next~. النجاح ينادي ~next~ بطلب سياقه فيه الاسم؛ ~CrewName(r *http.Request) string~ ترجعه أو نص فاضي. كل الأخطاء ~text/plain; charset=utf-8~. الوقت من الساعة المحقونة بس؛ مفيش ~Report~، و ~main~ فاضية.

النجمة التالتة: hash متخزن بايظ (خطأ غير ~ErrMismatchedHashAndPassword~) يرجع 500 و ~desk failure\n~ مش 401.`,
		Hints: []string{"فرّق بين hash كود الدخول وتوقيع التصريح.", "افصل تحميل السر وLogin عن middleware ~Require~.", "راجع خيارات ~ParseWithClaims~ الخاصة بالمصدر والوقت."},
		En: &LevelText{LessonRevision: 1, Steps: ch27Lesson271.english, Title: "Signed Pass", Goal: "Issue short-lived entry passes after checking crew codes with bcrypt", Files: []LabFile{{Name: "main.go", Starter: ch271Starter}},
			Lesson: `# A Signed Pass for the Quarters Gate

## Before you write
- If the crew code file is stolen, what stops the thief from using it?
- How does the server know a pass it receives was issued by itself and is still valid?

## Codes are never stored
bcrypt is a deliberately slow hash with a random salt, so a stolen hash is hard to reverse. Import ~golang.org/x/crypto/bcrypt~. ~bcrypt.GenerateFromPassword~ takes the code and a cost and returns a hash; ~bcrypt.MinCost~ is for tests only, ~bcrypt.DefaultCost~ for production.
~~~go
hash, err := bcrypt.GenerateFromPassword([]byte("nebula-7"), bcrypt.DefaultCost)
if err != nil { return err }
err = bcrypt.CompareHashAndPassword(hash, []byte("nebula-7")) // nil
~~~
~bcrypt.CompareHashAndPassword~ returns ~nil~ on a match. A wrong code returns ~bcrypt.ErrMismatchedHashAndPassword~; test it with ~errors.Is~. Any other error means a corrupted stored hash: your fault, not the client's.
~~~go
if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) { deny(w); return }
~~~
Every character of a code matters, including spaces; never trim it. An unknown name still runs a comparison against a decoy hash generated once, so response time does not reveal who exists, and it gets the same response as a wrong code.

## The secret comes from the environment
~os.Getenv~ from the command-line chapter reads an environment variable. Here a function of the same shape is injected so tests never change process environment; the real ~main~ passes ~os.Getenv~.
~~~go
secret, err := LoadSecret(os.Getenv)
~~~
The secret is base64 because it is random bytes, not text. ~encoding/base64~ provides ~base64.StdEncoding~: ~EncodeToString~ turns bytes into text, and ~DecodeString~ turns it back or fails when the text is not base64.
~~~go
text := base64.StdEncoding.EncodeToString([]byte("abc")) // "YWJj"
key, err := base64.StdEncoding.DecodeString(text)
~~~
An HMAC key shorter than 32 bytes is weak; reject it. Never fall back to a default secret written in code.

## JWT: claims and a signature
A JWT is three base64url parts separated by dots: a header naming the algorithm, a payload of claims, and a signature. Anyone can read the payload; only the signature prevents tampering. Import ~github.com/golang-jwt/jwt/v5~; its package name is ~jwt~.

~jwt.RegisteredClaims~ has ~Issuer~ (iss), ~Subject~ (sub), ~IssuedAt~ (iat), and ~ExpiresAt~ (exp). ~jwt.NewNumericDate~ converts a ~time.Time~ to whole seconds. ~jwt.NewWithClaims~ with ~jwt.SigningMethodHS256~ builds the token, and ~SignedString~ signs it with the secret.
~~~go
now := clock.Now()
claims := jwt.RegisteredClaims{Issuer: "hangar", Subject: "rigel",
 IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute))}
pass, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
~~~
A short lifetime limits the damage of a stolen pass, because nothing can revoke it before it expires.

## Verify before trusting
~jwt.ParseWithClaims~ decodes the token into your claims, calls a keyfunc for the key, then checks signature and time. Options tighten validation:
~~~go
var c jwt.RegisteredClaims
_, err := jwt.ParseWithClaims(raw, &c, func(*jwt.Token) (any, error) { return secret, nil },
 jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("hangar"),
 jwt.WithTimeFunc(clock.Now), jwt.WithExpirationRequired(), jwt.WithIssuedAt())
~~~
~jwt.WithValidMethods~ rejects ~none~, HS512, and every other algorithm even with a valid signature. ~jwt.WithIssuer~ matches iss. ~jwt.WithTimeFunc~ injects the clock; ~clock.Now~ is a method value. ~jwt.WithExpirationRequired~ rejects tokens without exp, and ~jwt.WithIssuedAt~ rejects a future iat.

A token is valid until the instant before exp; at exp itself it has expired. There is no implicit leeway.

## Authentication middleware
The ~Authorization~ header is exactly ~Bearer <pass>~ with one space. ~strings.HasPrefix~ and ~strings.TrimPrefix~ separate the token. Failure is 401 with ~WWW-Authenticate: Bearer~ so the client knows which credential is expected, and ~next~ is not called.
~~~go
if !strings.HasPrefix(h, "Bearer ") { deny(w); return }
r = r.WithContext(context.WithValue(r.Context(), crewKey{}, c.Subject))
next.ServeHTTP(w, r)
~~~
Success stores the crew name in context under a private key, as in 26-3, so the inner handler reads it without touching the token.

## Common mistakes
~~~go
if string(stored) == code { open() } // Code stored as plain text
jwt.Parse(raw, keyFunc) // Without WithValidMethods an HS512 token with the same secret passes
var secret = []byte("dev-secret") // Secret written in code
~~~

## Your mission
Implement ~LoadSecret(getenv func(string) string) ([]byte,error)~: read only ~ORBITAL_DESK_SECRET~, decode it with ~base64.StdEncoding~, and reject empty, invalid, or shorter-than-32-byte values. Implement ~NewPassDesk(secret []byte,crew map[string][]byte,clock labkit.Clock) (*PassDesk,error)~: a secret shorter than 32 bytes or a ~nil~ clock returns ~nil,error~. ~crew~ maps a name to a bcrypt hash. Copy both secret and map; later changes to them have no effect.

~(d *PassDesk) Login(w,r)~ is a handler for any path. Non-~POST~ returns 405, ~Allow: POST~, and ~use POST\n~. The body is capped at 4096 bytes including whitespace and holds one JSON object ~{"crew":"...","code":"..."}~. Invalid JSON, an oversized body, or an empty name or code returns 400 and ~bad login\n~. An unknown name or wrong code returns the same 401 and ~denied\n~. Success returns 200, ~application/json~, and ~{"pass":"<jwt>"}~: HS256 with the secret, ~iss=orbital-desk~, ~sub~ the name, ~iat~ the clock time, and ~exp~ exactly 10 minutes later.

~(d *PassDesk) Require(next http.Handler) http.Handler~ accepts only ~Authorization: Bearer <pass>~ for an HS256 token with the same secret, the correct ~iss~, a nonempty ~sub~, a present ~exp~ after the current time, and an ~iat~ not in the future. Any failure returns 401, ~WWW-Authenticate: Bearer~, and ~pass required\n~ without calling ~next~. Success calls ~next~ with a request whose context holds the name; ~CrewName(r *http.Request) string~ returns it or an empty string. All errors use ~text/plain; charset=utf-8~. Time comes only from the injected clock; no ~Report~, and ~main~ stays empty.

Third star: a corrupted stored hash (an error other than ~ErrMismatchedHashAndPassword~) returns 500 and ~desk failure\n~, not 401.`,
			Hints: []string{"Separate the login-code hash from the signed pass.", "Keep secret loading, ~Login~, and ~Require~ as distinct jobs.", "Review ~ParseWithClaims~ options for issuer and time."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch271Hidden},
	})
}

const ch271Starter = `package main

import (
	"errors"
	"net/http"

	"orbital/colony/labkit"
)

type PassDesk struct{}

func LoadSecret(getenv func(string) string) ([]byte, error) { return nil, errors.New("TODO") }
func NewPassDesk(secret []byte, crew map[string][]byte, clock labkit.Clock) (*PassDesk, error) {
	return nil, errors.New("TODO")
}
func (d *PassDesk) Login(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
func (d *PassDesk) Require(next http.Handler) http.Handler       { return next }
func CrewName(r *http.Request) string                            { return "" }
func main()                                                      {}
`

const ch271Hidden = `package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"orbital/colony/labkit"
)

func passBytes(rng *rand.Rand, n int) []byte {
	b := make([]byte, n)
	rng.Read(b)
	return b
}
func passLogin(d *PassDesk, method, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	d.Login(w, httptest.NewRequest(method, "/login", strings.NewReader(body)))
	return w
}
func passVisit(d *PassDesk, auth string, reached *string) *httptest.ResponseRecorder {
	*reached = ""
	h := d.Require(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*reached = "in:" + CrewName(r)
		w.Write([]byte(*reached))
	}))
	r := httptest.NewRequest("GET", "/quarters", nil)
	if auth != "" {
		r.Header.Set("Authorization", auth)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func passSign(t *testing.T, method jwt.SigningMethod, key any, claims jwt.Claims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// Boundaries: missing, malformed and short secrets; exactly 32 bytes; other variables ignored.
func TestHiddenPassSecretFromEnv(t *testing.T) {
	rng := rand.New(rand.NewSource(2711))
	good := passBytes(rng, 32)
	long := passBytes(rng, 48)
	for _, tc := range []struct {
		env  map[string]string
		ok   bool
		want []byte
	}{
		{map[string]string{}, false, nil},
		{map[string]string{"ORBITAL_DESK_SECRET": ""}, false, nil},
		{map[string]string{"SECRET": base64.StdEncoding.EncodeToString(good)}, false, nil},
		{map[string]string{"ORBITAL_DESK_SECRET": "not base64!!"}, false, nil},
		{map[string]string{"ORBITAL_DESK_SECRET": base64.StdEncoding.EncodeToString(good[:31])}, false, nil},
		{map[string]string{"ORBITAL_DESK_SECRET": base64.StdEncoding.EncodeToString(good)}, true, good},
		{map[string]string{"ORBITAL_DESK_SECRET": base64.StdEncoding.EncodeToString(long)}, true, long},
	} {
		got, err := LoadSecret(func(k string) string { return tc.env[k] })
		if (err == nil) != tc.ok || string(got) != string(tc.want) {
			t.Fatalf("env %v: %x %v", tc.env, got, err)
		}
	}
}

// Boundaries: invalid desks, structural login errors, the 4096-byte cap, wrong codes, exact claims.
func TestHiddenPassLoginIssuesShortPass(t *testing.T) {
	rng := rand.New(rand.NewSource(2712))
	secret := passBytes(rng, 32)
	clock := labkit.NewFakeClock(time.Unix(1_800_000_000+rng.Int63n(1_000_000), 0))
	if d, err := NewPassDesk(secret[:31], map[string][]byte{}, clock); d != nil || err == nil {
		t.Fatal("short secret accepted")
	}
	if d, err := NewPassDesk(secret, map[string][]byte{}, nil); d != nil || err == nil {
		t.Fatal("nil clock accepted")
	}
	crew := map[string][]byte{}
	codes := map[string]string{}
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("pilot-%d", rng.Intn(100000))
		code := fmt.Sprintf("code %d 🚀", rng.Intn(1000000))
		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.MinCost)
		if err != nil {
			t.Fatal(err)
		}
		crew[name], codes[name] = hash, code
	}
	key := append([]byte(nil), secret...)
	d, err := NewPassDesk(key, crew, clock)
	if err != nil || d == nil {
		t.Fatal("construct", err)
	}
	key[0] ^= 0xff
	var first string
	for name := range codes {
		first = name
		break
	}
	delete(crew, first)
	if w := passLogin(d, "GET", ""); w.Code != 405 || w.Header().Get("Allow") != "POST" || w.Body.String() != "use POST\n" {
		t.Fatal("method", w.Code, w.Header())
	}
	login := func(name, code string) string {
		raw, _ := json.Marshal(map[string]string{"crew": name, "code": code})
		return string(raw)
	}
	for _, body := range []string{"", "null", "[]", "{}", "{\"crew\":\"" + first + "\"}", "{\"crew\":7,\"code\":\"x\"}", login(first, codes[first]) + " {}", login(first, codes[first]) + strings.Repeat(" ", 4097)} {
		if w := passLogin(d, "POST", body); w.Code != 400 || w.Body.String() != "bad login\n" || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
			t.Fatalf("structural %q: %d %q", body, w.Code, w.Body.String())
		}
	}
	for _, body := range []string{login(first, codes[first]+"x"), login(first, " "+codes[first]), login("ghost", codes[first]), login("ghost", "decoy")} {
		if w := passLogin(d, "POST", body); w.Code != 401 || w.Body.String() != "denied\n" {
			t.Fatalf("credential %q: %d %q", body, w.Code, w.Body.String())
		}
	}
	for name, code := range codes {
		body := login(name, code)
		body += strings.Repeat(" ", 4096-len(body))
		w := passLogin(d, "POST", body)
		var out struct{ Pass string }
		if w.Code != 200 || w.Header().Get("Content-Type") != "application/json" || json.Unmarshal(w.Body.Bytes(), &out) != nil || out.Pass == "" {
			t.Fatalf("login %s: %d %q", name, w.Code, w.Body.String())
		}
		var claims jwt.RegisteredClaims
		token, err := jwt.ParseWithClaims(out.Pass, &claims, func(*jwt.Token) (any, error) { return secret, nil }, jwt.WithTimeFunc(clock.Now))
		now := clock.Now()
		if err != nil || token.Method.Alg() != "HS256" || claims.Issuer != "orbital-desk" || claims.Subject != name ||
			claims.IssuedAt == nil || !claims.IssuedAt.Equal(now) || claims.ExpiresAt == nil || !claims.ExpiresAt.Equal(now.Add(10*time.Minute)) {
			t.Fatalf("claims %+v %v", claims, err)
		}
		clock.Advance(time.Duration(1+rng.Intn(50)) * time.Second)
	}
}

// Boundaries: the exact expiry instant, forged signatures, algorithm confusion, missing and future claims, header shape.
func TestHiddenPassRequireVerifiesEverything(t *testing.T) {
	rng := rand.New(rand.NewSource(2713))
	secret := passBytes(rng, 40)
	clock := labkit.NewFakeClock(time.Unix(1_900_000_000+rng.Int63n(1_000_000), 0))
	hash, _ := bcrypt.GenerateFromPassword([]byte("aurora"), bcrypt.MinCost)
	d, err := NewPassDesk(secret, map[string][]byte{"vega": hash}, clock)
	if err != nil {
		t.Fatal(err)
	}
	w := passLogin(d, "POST", "{\"crew\":\"vega\",\"code\":\"aurora\"}")
	var out struct{ Pass string }
	if json.Unmarshal(w.Body.Bytes(), &out) != nil {
		t.Fatal("login", w.Body.String())
	}
	var reached string
	expectDenied := func(label, auth string) {
		t.Helper()
		w := passVisit(d, auth, &reached)
		if w.Code != 401 || w.Header().Get("WWW-Authenticate") != "Bearer" || w.Body.String() != "pass required\n" || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" || reached != "" {
			t.Fatalf("%s accepted: %d %q %q", label, w.Code, w.Body.String(), reached)
		}
	}
	if w := passVisit(d, "Bearer "+out.Pass, &reached); w.Code != 200 || w.Body.String() != "in:vega" {
		t.Fatal("valid pass", w.Code, w.Body.String())
	}
	now := clock.Now()
	base := func() jwt.RegisteredClaims {
		return jwt.RegisteredClaims{Issuer: "orbital-desk", Subject: "vega", IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute))}
	}
	other := passBytes(rng, 32)
	wrongIss, noExp, future, noSub := base(), base(), base(), base()
	wrongIss.Issuer = "orbital-dock"
	noExp.ExpiresAt = nil
	future.IssuedAt = jwt.NewNumericDate(now.Add(time.Minute))
	noSub.Subject = ""
	none := passSign(t, jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType, base())
	parts := strings.Split(out.Pass, ".")
	forgedBody := base64.RawURLEncoding.EncodeToString([]byte(strings.Replace(string(mustDecode(t, parts[1])), "vega", "evil", 1)))
	for label, auth := range map[string]string{
		"missing":       "",
		"lowercase":     "bearer " + out.Pass,
		"no token":      "Bearer ",
		"double space":  "Bearer  " + out.Pass,
		"basic":         "Basic " + out.Pass,
		"tampered":      "Bearer " + parts[0] + "." + forgedBody + "." + parts[2],
		"other secret":  "Bearer " + passSign(t, jwt.SigningMethodHS256, other, base()),
		"hs512":         "Bearer " + passSign(t, jwt.SigningMethodHS512, secret, base()),
		"none":          "Bearer " + none,
		"wrong issuer":  "Bearer " + passSign(t, jwt.SigningMethodHS256, secret, wrongIss),
		"no expiry":     "Bearer " + passSign(t, jwt.SigningMethodHS256, secret, noExp),
		"future issued": "Bearer " + passSign(t, jwt.SigningMethodHS256, secret, future),
		"no subject":    "Bearer " + passSign(t, jwt.SigningMethodHS256, secret, noSub),
	} {
		expectDenied(label, auth)
	}
	clock.Advance(10*time.Minute - time.Second)
	if w := passVisit(d, "Bearer "+out.Pass, &reached); w.Code != 200 {
		t.Fatal("pass expired early")
	}
	clock.Advance(time.Second)
	expectDenied("expiry instant", "Bearer "+out.Pass)
	if CrewName(httptest.NewRequest("GET", "/", nil)) != "" {
		t.Fatal("crew without pass")
	}
}
func mustDecode(t *testing.T, s string) []byte {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestHiddenBonusPassBrokenHashIsServerFault(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(2_000_000_000, 0))
	good, _ := bcrypt.GenerateFromPassword([]byte("lumen"), bcrypt.MinCost)
	d, err := NewPassDesk(make([]byte, 32), map[string][]byte{"orion": []byte("not-a-bcrypt-hash"), "lyra": good}, clock)
	if err != nil {
		t.Fatal(err)
	}
	if w := passLogin(d, "POST", "{\"crew\":\"orion\",\"code\":\"lumen\"}"); w.Code != 500 || w.Body.String() != "desk failure\n" {
		t.Fatal("broken hash", w.Code, w.Body.String())
	}
	if w := passLogin(d, "POST", "{\"crew\":\"lyra\",\"code\":\"lumen\"}"); w.Code != 200 {
		t.Fatal("healthy crew", w.Code)
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch27Lesson272.arabic, ID: "27-2-sessions", Chapter: 27, Kind: "lab", Title: "الجلسات", Goal: "افتح جلسات للمساكن بكوكي آمن وتوكن CSRF وتدوير عند كل دخول",
		Requires: []string{"sel:http.Cookie"}, Bridge: []int{59}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch272Starter}},
		Lesson: `# جلسة على السيرفر وكوكي في المتصفح

## قبل ما تكتب
- لو حد قدر يخلّي متصفحك يبعت طلب للقاعدة، الكوكي هيتبعت معاه؟ وإزاي ترفض الطلب ده؟
- ليه رقم الجلسة لازم يتغير بعد كل دخول ناجح؟

## الجلسة سجل عندك
الجلسة هنا سجل في خريطة على السيرفر: رقم عشوائي ← اسم الطاقم وتوكن CSRF ووقت الانتهاء. العميل شايل الرقم بس، ومفيهوش أي معلومة يقدر يعدّلها. الخروج يمسح السجل، فالرقم القديم يموت فوراً، عكس JWT.
~~~go
type visit struct { crew, csrf string; expires time.Time }
rooms := map[string]visit{}
rooms[id] = visit{crew: "rigel", csrf: token, expires: now.Add(time.Hour)}
~~~
الخريطة مشتركة بين طلبات متوازية، فاحميها بـ ~sync.Mutex~ زي فصل 25.

## أرقام لا تتخمن
الرقم لازم من مصدر عشوائي آمن. بنحقن ~io.Reader~؛ في التشغيل ~crypto/rand.Reader~ وفي التست مصدر ثابت. ~io.ReadFull~ من فصل الخزنة تملا الـ buffer أو ترجع خطأ، و ~%x~ يحوّل 16 بايت لـ 32 حرف hex.
~~~go
b := make([]byte, 16)
if _, err := io.ReadFull(random, b); err != nil { return "", err }
id := fmt.Sprintf("%x", b)
~~~
لو المصدر فشل متطلّعش جلسة. لو الرقم الجديد موجود أصلاً متكتبش فوق جلسة حد تاني.

## الكوكي وخصائصه
~http.Cookie~ فيها ~Name~ و ~Value~ و ~Path~ و ~MaxAge~ و ~HttpOnly~ و ~Secure~ و ~SameSite~. ~http.SetCookie~ تضيف هيدر ~Set-Cookie~ للرد.
~~~go
http.SetCookie(w, &http.Cookie{Name: "hatch", Value: id, Path: "/", MaxAge: 600,
 HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})
~~~
~HttpOnly~ تمنع JavaScript في الصفحة من قراءة الكوكي. ~Secure~ تخليها تتبعت على HTTPS بس. ~http.SameSiteStrictMode~ تمنع المتصفح يبعتها مع طلب جاي من موقع تاني. ~MaxAge~ بالثواني؛ القيمة السالبة في Go معناها «امسحها دلوقتي» وبتتكتب ~Max-Age=0~. القيمة صفر معناها مفيش ~Max-Age~ خالص. سيب ~Domain~ فاضي عشان الكوكي متتبعتش لدومينات فرعية.

## قراءة الكوكي
~r.Cookie(name)~ ترجع أول كوكي بالاسم ده، أو ~http.ErrNoCookie~ لو مش موجودة.
~~~go
c, err := r.Cookie("hatch")
if errors.Is(err, http.ErrNoCookie) { /* مفيش جلسة */ }
~~~
قيمة الكوكي جاية من العميل؛ دوّر عليها في الخريطة ومتثقش فيها لوحدها. الجلسة اللي وقتها خلص امسحها وقت ما تلاقيها (lazy)، وعاملها كأنها مش موجودة.

## CSRF
موقع خبيث يقدر يخلي المتصفح يبعت ~POST~ للقاعدة، والمتصفح ممكن يضيف الكوكي لوحده. الحل: توكن سري تاني بيترجع في جسم رد الدخول، والعميل الحقيقي بيبعته في هيدر ~X-CSRF-Token~ مع كل طلب بيغيّر حالة. الموقع الخبيث مش شايف الرد، فمش عارف التوكن.

قارن التوكن بـ ~crypto/subtle~: ~subtle.ConstantTimeCompare~ ترجع 1 لو البايتات متساوية، ووقتها مش بيعتمد على مكان أول اختلاف.
~~~go
if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-CSRF-Token")), []byte(v.csrf)) != 1 {
 w.WriteHeader(403); return
}
~~~
الطلبات اللي بتقرا بس، زي ~GET~، مش محتاجة التوكن.

## التدوير
لو الرقم اتعرف قبل الدخول (session fixation)، المهاجم يستنى الضحية تدخل ويستخدمه. عشان كده كل دخول ناجح يطلّع رقم جديد وتوكن جديد، ويمسح الجلسة القديمة اللي جت في الكوكي.

## أخطاء شائعة
~~~go
id := fmt.Sprint(len(rooms)+1) // رقم متسلسل سهل التخمين
http.SetCookie(w, &http.Cookie{Name: "hatch", Value: id}) // من غير HttpOnly وSecure وSameSite
if r.Header.Get("X-CSRF-Token") == v.csrf { } // مقارنة وقتها بيفضح التوكن
~~~

## المهمة
نفّذ ~NewQuarters(check func(crew, code string) bool, clock labkit.Clock, random io.Reader) (*Quarters,error)~؛ أي مدخل ~nil~ يرجع ~nil,error~. ~check~ بتحكم على الاسم والكود، و ~random~ آمن للتزامن. ~(q *Quarters) ServeHTTP~ فيها تلات مسارات حرفية: ~POST /enter~ و ~GET /me~ و ~POST /logout~. طريقة غلط: 405 و ~Allow~ بالطريقة الصح و ~wrong method\n~. مسار تاني: 404 و ~no room\n~. الأخطاء ~text/plain; charset=utf-8~.

~/enter~: جسم حده 4096 بايت فيه ~{"crew":"...","code":"..."}~. JSON بايظ أو أكبر أو حقل فاضي: 400 و ~bad entry\n~. ~check~ رفض: 401 و ~denied\n~. بعدها اقرا 16 بايت للرقم و 16 للتوكن بالترتيب، وكل واحد 32 حرف hex صغير. فشل المصدر أو رقم موجود أصلاً: 500 و ~quarters failure\n~ من غير كوكي ولا تغيير. النجاح يمسح جلسة الكوكي القديمة لو موجودة، ويحفظ جلسة تخلص بعد 30 دقيقة، ويحط كوكي ~quarters~ بالرقم و ~Path=/~ و ~MaxAge=1800~ و ~HttpOnly~ و ~Secure~ و ~SameSiteStrictMode~ من غير ~Domain~، ويرجع 200 و ~application/json~ و ~{"csrf":"..."}~.

الجلسة صالحة لو كوكي ~quarters~ موجودة في الخريطة والوقت قبل الانتهاء؛ عند لحظة الانتهاء نفسها اتمسحت. ~/me~ من غير جلسة: 401 و ~no session\n~؛ معاها: 200 و ~crew NAME\n~. ~/logout~ من غير جلسة: 401 و ~no session\n~؛ هيدر ~X-CSRF-Token~ مختلف: 403 و ~bad csrf\n~ والجلسة تفضل. النجاح يمسح الجلسة ويحط كوكي ~quarters~ فاضية بنفس الخصائص و ~MaxAge=-1~، ويرجع 204 بجسم فاضي. الوقت من الساعة بس؛ مفيش ~Report~، و ~main~ فاضية.`,
		Hints: []string{"فرّق بين ID الموجود عند العميل وسجل الجلسة على السيرفر.", "رتّب التحقق والإنشاء والتدوير كعمليات واضحة.", "راجع ~ConstantTimeCompare~ عند مقارنة توكن CSRF."},
		En: &LevelText{LessonRevision: 1, Steps: ch27Lesson272.english, Title: "Sessions", Goal: "Open quarters sessions with a secure cookie, a CSRF token, and rotation on every login", Files: []LabFile{{Name: "main.go", Starter: ch272Starter}},
			Lesson: `# A Server Session and a Browser Cookie

## Before you write
- If another site makes your browser send a request to the base, is the cookie sent too? How do you reject that request?
- Why must the session ID change after every successful login?

## The session is your record
A session here is a record in a server-side map: random ID to crew name, CSRF token, and expiry. The client holds only the ID and nothing it could edit. Logout deletes the record, so the old ID dies immediately, unlike a JWT.
~~~go
type visit struct { crew, csrf string; expires time.Time }
rooms := map[string]visit{}
rooms[id] = visit{crew: "rigel", csrf: token, expires: now.Add(time.Hour)}
~~~
The map is shared by concurrent requests, so guard it with ~sync.Mutex~ as in chapter 25.

## IDs nobody can guess
IDs must come from a secure random source. An ~io.Reader~ is injected: ~crypto/rand.Reader~ in production, a fixed source in tests. ~io.ReadFull~ from the vault chapter fills the buffer or fails, and ~%x~ turns 16 bytes into 32 hex characters.
~~~go
b := make([]byte, 16)
if _, err := io.ReadFull(random, b); err != nil { return "", err }
id := fmt.Sprintf("%x", b)
~~~
If the source fails, issue no session. If the new ID already exists, never overwrite someone else's session.

## The cookie and its attributes
~http.Cookie~ has ~Name~, ~Value~, ~Path~, ~MaxAge~, ~HttpOnly~, ~Secure~, and ~SameSite~. ~http.SetCookie~ adds a ~Set-Cookie~ header to the response.
~~~go
http.SetCookie(w, &http.Cookie{Name: "hatch", Value: id, Path: "/", MaxAge: 600,
 HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode})
~~~
~HttpOnly~ hides the cookie from page JavaScript. ~Secure~ sends it only over HTTPS. ~http.SameSiteStrictMode~ stops the browser from sending it with requests started by another site. ~MaxAge~ is in seconds; a negative Go value means delete now and is written as ~Max-Age=0~. Zero means no ~Max-Age~ at all. Leave ~Domain~ empty so subdomains never receive the cookie.

## Reading the cookie
~r.Cookie(name)~ returns the first cookie with that name, or ~http.ErrNoCookie~ when absent.
~~~go
c, err := r.Cookie("hatch")
if errors.Is(err, http.ErrNoCookie) { /* no session */ }
~~~
The value comes from the client; look it up in the map and never trust it alone. Delete an expired session when you find it (lazily) and treat it as absent.

## CSRF
A hostile site can make the browser send a ~POST~ to the base, and the browser may attach the cookie by itself. The defense is a second secret token returned in the login response body; the real client sends it in an ~X-CSRF-Token~ header with every state-changing request. The hostile site cannot read the response, so it does not know the token.

Compare tokens with ~crypto/subtle~: ~subtle.ConstantTimeCompare~ returns 1 when the bytes are equal, and its timing does not depend on where the first difference is.
~~~go
if subtle.ConstantTimeCompare([]byte(r.Header.Get("X-CSRF-Token")), []byte(v.csrf)) != 1 {
 w.WriteHeader(403); return
}
~~~
Read-only requests such as ~GET~ need no token.

## Rotation
If an ID is known before login (session fixation), an attacker can wait for the victim to log in and then use it. So every successful login issues a fresh ID and token and deletes the old session named by the incoming cookie.

## Common mistakes
~~~go
id := fmt.Sprint(len(rooms)+1) // Sequential, easy to guess
http.SetCookie(w, &http.Cookie{Name: "hatch", Value: id}) // No HttpOnly, Secure, or SameSite
if r.Header.Get("X-CSRF-Token") == v.csrf { } // Timing leaks the token
~~~

## Your mission
Implement ~NewQuarters(check func(crew, code string) bool, clock labkit.Clock, random io.Reader) (*Quarters,error)~; any ~nil~ input returns ~nil,error~. ~check~ judges name and code, and ~random~ is safe for concurrent use. ~(q *Quarters) ServeHTTP~ serves three exact paths: ~POST /enter~, ~GET /me~, and ~POST /logout~. A wrong method returns 405, ~Allow~ with the right method, and ~wrong method\n~. Any other path returns 404 and ~no room\n~. Errors use ~text/plain; charset=utf-8~.

~/enter~: a body of at most 4096 bytes holding ~{"crew":"...","code":"..."}~. Invalid JSON, an oversized body, or an empty field returns 400 and ~bad entry\n~. A rejecting ~check~ returns 401 and ~denied\n~. Then read 16 bytes for the ID and 16 for the token, in that order, each as 32 lowercase hex characters. A source failure or an ID already present returns 500 and ~quarters failure\n~ with no cookie and no change. Success deletes the old cookie's session if any, stores a session expiring in 30 minutes, sets cookie ~quarters~ with the ID, ~Path=/~, ~MaxAge=1800~, ~HttpOnly~, ~Secure~, and ~SameSiteStrictMode~ without ~Domain~, and returns 200, ~application/json~, and ~{"csrf":"..."}~.

A session is valid when the ~quarters~ cookie is in the map and the time is before expiry; at the expiry instant it is gone. ~/me~ without a session returns 401 and ~no session\n~; with one, 200 and ~crew NAME\n~. ~/logout~ without a session returns 401 and ~no session\n~; a different ~X-CSRF-Token~ header returns 403 and ~bad csrf\n~ and the session stays. Success deletes the session, sets an empty ~quarters~ cookie with the same attributes and ~MaxAge=-1~, and returns 204 with an empty body. Time comes only from the clock; no ~Report~, and ~main~ stays empty.`,
			Hints: []string{"Separate the client-held ID from the server-side session record.", "Lay out validation, creation, and rotation as distinct operations.", "Review ~ConstantTimeCompare~ for CSRF token comparison."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch272Hidden},
	})
}

const ch272Starter = `package main

import (
	"errors"
	"io"
	"net/http"

	"orbital/colony/labkit"
)

type Quarters struct{}

func NewQuarters(check func(crew, code string) bool, clock labkit.Clock, random io.Reader) (*Quarters, error) {
	return nil, errors.New("TODO")
}
func (q *Quarters) ServeHTTP(w http.ResponseWriter, r *http.Request) { http.NotFound(w, r) }
func main()                                                          {}
`

const ch272Hidden = `package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
	"time"

	"orbital/colony/labkit"
)

type lockedRand struct {
	mu  sync.Mutex
	rng *rand.Rand
}

func (l *lockedRand) Read(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rng.Read(p)
}

// loopRand repeats the same 32 bytes forever, so every login draws the same ID.
type loopRand struct{ n int }

func (l *loopRand) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(l.n % 32)
		l.n++
	}
	return len(p), nil
}

func roomCall(h http.Handler, method, path, body string, cookie *http.Cookie, csrf string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	if cookie != nil {
		r.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	}
	if csrf != "" {
		r.Header.Set("X-CSRF-Token", csrf)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func roomCookie(t *testing.T, w *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()
	cs := w.Result().Cookies()
	if len(cs) != 1 {
		t.Fatalf("want one cookie, got %v", w.Header().Values("Set-Cookie"))
	}
	return cs[0]
}
func roomEnter(t *testing.T, h http.Handler, crew, code string, old *http.Cookie) (*http.Cookie, string) {
	t.Helper()
	w := roomCall(h, "POST", "/enter", fmt.Sprintf("{\"crew\":%q,\"code\":%q}", crew, code), old, "")
	var out struct{ CSRF string }
	if w.Code != 200 || w.Header().Get("Content-Type") != "application/json" || json.Unmarshal(w.Body.Bytes(), &out) != nil {
		t.Fatalf("enter %s: %d %q", crew, w.Code, w.Body.String())
	}
	c := roomCookie(t, w)
	if c.Name != "quarters" || len(c.Value) != 32 || strings.Trim(c.Value, "0123456789abcdef") != "" || c.Path != "/" || c.MaxAge != 1800 ||
		!c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteStrictMode || c.Domain != "" {
		t.Fatalf("cookie flags %+v", c)
	}
	if len(out.CSRF) != 32 || strings.Trim(out.CSRF, "0123456789abcdef") != "" || out.CSRF == c.Value {
		t.Fatalf("csrf %q", out.CSRF)
	}
	return c, out.CSRF
}
func roomExpect(t *testing.T, w *httptest.ResponseRecorder, status int, body string) {
	t.Helper()
	if w.Code != status || w.Body.String() != body || body != "" && w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("want %d %q, got %d %q", status, body, w.Code, w.Body.String())
	}
}

func roomCheck(codes map[string]string) func(string, string) bool {
	return func(crew, code string) bool { want, ok := codes[crew]; return ok && want == code }
}

// Boundaries: invalid construction, routes and methods, structural input, credentials, cookie flags, injected randomness.
func TestHiddenSessionsEnterAndCookie(t *testing.T) {
	rng := rand.New(rand.NewSource(2721))
	clock := labkit.NewFakeClock(time.Unix(1_700_000_000+rng.Int63n(1000000), 0))
	random := &lockedRand{rng: rand.New(rand.NewSource(rng.Int63()))}
	codes := map[string]string{"altair": fmt.Sprintf("k-%d", rng.Intn(99999)), "deneb": fmt.Sprintf("k-%d", rng.Intn(99999))}
	for _, bad := range []func() (*Quarters, error){
		func() (*Quarters, error) { return NewQuarters(nil, clock, random) },
		func() (*Quarters, error) { return NewQuarters(roomCheck(codes), nil, random) },
		func() (*Quarters, error) { return NewQuarters(roomCheck(codes), clock, nil) },
	} {
		if q, err := bad(); q != nil || err == nil {
			t.Fatal("invalid quarters accepted")
		}
	}
	q, err := NewQuarters(roomCheck(codes), clock, random)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ method, path, allow string }{{"GET", "/enter", "POST"}, {"POST", "/me", "GET"}, {"GET", "/logout", "POST"}} {
		w := roomCall(q, tc.method, tc.path, "", nil, "")
		roomExpect(t, w, 405, "wrong method\n")
		if w.Header().Get("Allow") != tc.allow {
			t.Fatal("allow", tc, w.Header())
		}
	}
	roomExpect(t, roomCall(q, "GET", "/", "", nil, ""), 404, "no room\n")
	roomExpect(t, roomCall(q, "GET", "/me/", "", nil, ""), 404, "no room\n")
	for _, body := range []string{"", "null", "[]", "{\"crew\":\"altair\"}", "{\"crew\":\"altair\",\"code\":5}", "{\"crew\":\"altair\",\"code\":\"x\"} {}", "{\"crew\":\"altair\",\"code\":\"x\"}" + strings.Repeat(" ", 4096)} {
		w := roomCall(q, "POST", "/enter", body, nil, "")
		roomExpect(t, w, 400, "bad entry\n")
		if len(w.Result().Cookies()) != 0 {
			t.Fatal("cookie on bad entry")
		}
	}
	w := roomCall(q, "POST", "/enter", "{\"crew\":\"altair\",\"code\":\"wrong\"}", nil, "")
	roomExpect(t, w, 401, "denied\n")
	if len(w.Result().Cookies()) != 0 {
		t.Fatal("cookie on denial")
	}
	seen := map[string]bool{}
	for i := 0; i < 10; i++ {
		crew := "altair"
		if i%2 == 1 {
			crew = "deneb"
		}
		c, csrf := roomEnter(t, q, crew, codes[crew], nil)
		if seen[c.Value] || seen[csrf] {
			t.Fatal("reused random value")
		}
		seen[c.Value], seen[csrf] = true, true
		roomExpect(t, roomCall(q, "GET", "/me?x=1", "", c, ""), 200, "crew "+crew+"\n")
	}
	roomExpect(t, roomCall(q, "GET", "/me", "", nil, ""), 401, "no session\n")
	roomExpect(t, roomCall(q, "GET", "/me", "", &http.Cookie{Name: "quarters", Value: strings.Repeat("ab", 16)}, ""), 401, "no session\n")
	roomExpect(t, roomCall(q, "GET", "/me", "", &http.Cookie{Name: "other", Value: "x"}, ""), 401, "no session\n")
}

// Boundaries: expiry instant, CSRF mismatch keeps the session, logout clears, rotation kills the old ID, random failure and collision.
func TestHiddenSessionsCSRFRotationAndLogout(t *testing.T) {
	rng := rand.New(rand.NewSource(2722))
	clock := labkit.NewFakeClock(time.Unix(1_750_000_000, 0))
	codes := map[string]string{"mira": "sol-" + fmt.Sprint(rng.Intn(1000))}
	q, _ := NewQuarters(roomCheck(codes), clock, &lockedRand{rng: rand.New(rand.NewSource(rng.Int63()))})
	c, csrf := roomEnter(t, q, "mira", codes["mira"], nil)
	clock.Advance(30*time.Minute - time.Second)
	roomExpect(t, roomCall(q, "GET", "/me", "", c, ""), 200, "crew mira\n")
	clock.Advance(time.Second)
	roomExpect(t, roomCall(q, "GET", "/me", "", c, ""), 401, "no session\n")
	roomExpect(t, roomCall(q, "POST", "/logout", "", c, csrf), 401, "no session\n")

	c, csrf = roomEnter(t, q, "mira", codes["mira"], nil)
	roomExpect(t, roomCall(q, "POST", "/logout", "", c, ""), 403, "bad csrf\n")
	roomExpect(t, roomCall(q, "POST", "/logout", "", c, csrf[:31]+"z"), 403, "bad csrf\n")
	roomExpect(t, roomCall(q, "POST", "/logout", "", c, c.Value), 403, "bad csrf\n")
	roomExpect(t, roomCall(q, "GET", "/me", "", c, ""), 200, "crew mira\n")

	rotated, csrf2 := roomEnter(t, q, "mira", codes["mira"], c)
	if rotated.Value == c.Value || csrf2 == csrf {
		t.Fatal("rotation reused values")
	}
	roomExpect(t, roomCall(q, "GET", "/me", "", c, ""), 401, "no session\n")
	roomExpect(t, roomCall(q, "POST", "/logout", "", rotated, csrf), 403, "bad csrf\n")
	w := roomCall(q, "POST", "/logout", "", rotated, csrf2)
	roomExpect(t, w, 204, "")
	cleared := roomCookie(t, w)
	raw := w.Header().Get("Set-Cookie")
	if cleared.Name != "quarters" || cleared.Value != "" || cleared.MaxAge >= 0 || !strings.Contains(raw, "Max-Age=0") || cleared.Path != "/" || !cleared.HttpOnly || !cleared.Secure || cleared.SameSite != http.SameSiteStrictMode {
		t.Fatalf("clearing cookie %q", raw)
	}
	roomExpect(t, roomCall(q, "GET", "/me", "", rotated, ""), 401, "no session\n")

	failing, _ := NewQuarters(roomCheck(codes), clock, iotest.ErrReader(errors.New("entropy gone")))
	w = roomCall(failing, "POST", "/enter", "{\"crew\":\"mira\",\"code\":\""+codes["mira"]+"\"}", nil, "")
	roomExpect(t, w, 500, "quarters failure\n")
	if len(w.Result().Cookies()) != 0 {
		t.Fatal("cookie after random failure")
	}
	stuck, _ := NewQuarters(roomCheck(codes), clock, &loopRand{})
	first, _ := roomEnter(t, stuck, "mira", codes["mira"], nil)
	w = roomCall(stuck, "POST", "/enter", "{\"crew\":\"mira\",\"code\":\""+codes["mira"]+"\"}", nil, "")
	roomExpect(t, w, 500, "quarters failure\n")
	if len(w.Result().Cookies()) != 0 {
		t.Fatal("cookie after collision")
	}
	roomExpect(t, roomCall(stuck, "GET", "/me", "", first, ""), 200, "crew mira\n")
}

func TestHiddenSessionsParallel(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(1_760_000_000, 0))
	codes := map[string]string{}
	for i := 0; i < 32; i++ {
		codes[fmt.Sprintf("crew-%d", i)] = fmt.Sprintf("c%d", i)
	}
	q, _ := NewQuarters(roomCheck(codes), clock, &lockedRand{rng: rand.New(rand.NewSource(2723))})
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			crew := fmt.Sprintf("crew-%d", i)
			r := httptest.NewRequest("POST", "/enter", strings.NewReader(fmt.Sprintf("{\"crew\":%q,\"code\":\"c%d\"}", crew, i)))
			w := httptest.NewRecorder()
			q.ServeHTTP(w, r)
			var out struct{ CSRF string }
			cs := w.Result().Cookies()
			if w.Code != 200 || len(cs) != 1 || json.Unmarshal(w.Body.Bytes(), &out) != nil {
				t.Error("parallel enter", w.Code)
				return
			}
			for j := 0; j < 5; j++ {
				if w := roomCall(q, "GET", "/me", "", cs[0], ""); w.Body.String() != "crew "+crew+"\n" {
					t.Error("crossed sessions", w.Body.String())
				}
			}
			if w := roomCall(q, "POST", "/logout", "", cs[0], out.CSRF); w.Code != 204 {
				t.Error("parallel logout", w.Code)
			}
		}(i)
	}
	wg.Wait()
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch27Lesson273.arabic, ID: "27-3-cors", Chapter: 27, Kind: "lab", Title: "CORS", Goal: "اسمح لشاشات الأسطول المعروفة بس تقرا بيانات القاعدة، وحط هيدرات الحماية على كل رد",
		Requires: []string{"funcvalue"}, Bridge: []int{56}, Setup: seededLabRoom, Race: false,
		Files: []LabFile{{Name: "main.go", Starter: ch273Starter}},
		Lesson: `# درع CORS للبوابة

## قبل ما تكتب
- لو صفحة من موقع تاني طلبت بيانات القاعدة، مين اللي بيمنعها تقرا الرد: السيرفر ولا المتصفح؟
- لو الرد بيتغير حسب ~Origin~، الكاش اللي في النص يعرف ده إزاي؟

## الأصل (origin)
الأصل هو scheme و host و port مع بعض، زي ~https://deck.orbital.test:8443~ من غير مسار. المتصفح بيمنع صفحة من أصل تقرا رد من أصل تاني إلا لو الرد سمح بده. السيرفر بيستقبل الطلب عادي؛ CORS قرار قراءة في المتصفح، مش بديل للمصادقة.

المتصفح بيحط الأصل في هيدر ~Origin~. لو الأصل في القائمة، الرد يرجع ~Access-Control-Allow-Origin~ بنفس القيمة بالظبط. المطابقة حرفية: ~https://Deck...~ أو ~https://deck.../~ مش نفس الأصل.
~~~go
origin := r.Header.Get("Origin")
if slices.Contains(allowed, origin) {
 w.Header().Set("Access-Control-Allow-Origin", origin)
}
~~~
~slices.Contains~ بترجع ~true~ لو العنصر موجود. الأصل الغريب مش بيترفض بخطأ؛ الرد يكمل من غير هيدرات CORS والمتصفح هو اللي يمنع القراءة.

## الكوكيز مع الطلب
~Access-Control-Allow-Credentials: true~ بتسمح للمتصفح يبعت الكوكيز ويقرا الرد. مع الكوكيز ممنوع ~*~؛ لازم الأصل نفسه. عشان كده قائمة السماح هنا أصول محددة بس.

## الطلب التمهيدي (preflight)
قبل طلب «مش بسيط»، زي ~PATCH~ أو هيدر مخصص، المتصفح بيبعت ~OPTIONS~ فيه ~Access-Control-Request-Method~ وممكن ~Access-Control-Request-Headers~. ده سؤال، فالـ handler الحقيقي مبيشتغلش.
~~~go
preflight := r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != ""
~~~
الموافقة 204 من غير جسم، ومعاها ~Access-Control-Allow-Methods~ و ~Access-Control-Allow-Headers~ بالقوايم مفصولة بـ ~, ~، و ~Access-Control-Max-Age~ بالثواني عشان المتصفح يحفظ الموافقة. ~strings.Join~ تجمع القايمة، و ~strconv.Itoa~ تكتب الرقم.

أسماء الهيدرات مش حساسة لحالة الحروف. ~strings.EqualFold~ تقارن من غير حالة، و ~slices.ContainsFunc~ تدوّر بدالة:
~~~go
ok := slices.ContainsFunc(headers, func(a string) bool { return strings.EqualFold(a, "x-crew") })
~~~
الطريقة نفسها حساسة: ~get~ مش ~GET~.

## Vary
الكاش بيخزن الرد حسب الـ URL. لو الرد بيختلف حسب ~Origin~، لازم ~Vary: Origin~ وإلا الكاش يدّي أصل رد أصل تاني. ~Header().Add~ بتضيف قيمة جنب القديمة، و ~Set~ بتمسحها؛ استخدم ~Add~ عشان ~Vary~ اللي بيحطها الـ handler متضيعش. ~Header().Values~ بترجع كل القيم.
~~~go
w.Header().Add("Vary", "Origin")
w.Header().Add("Vary", "Accept-Encoding")
fmt.Println(w.Header().Values("Vary")) // [Origin Accept-Encoding]
~~~

## هيدرات الحماية
~X-Content-Type-Options: nosniff~ تمنع المتصفح يخمّن نوع المحتوى. ~X-Frame-Options: DENY~ تمنع عرض الصفحة جوه frame. ~Referrer-Policy: no-referrer~ متبعتش العنوان لمواقع تانية. ~Content-Security-Policy: default-src 'none'~ تمنع تحميل أي مصدر في رد API. حطهم قبل ما تنادي ~next~ عشان يطلعوا على كل رد.

## تحليل الأصل
~url.Parse~ من ~net/url~ بتقسم النص لـ ~Scheme~ و ~Host~ و ~User~ و ~Path~ و ~RawQuery~. الأصل الصحيح يرجع زي ما هو لو جمعت ~Scheme + "://" + Host~.
~~~go
u, err := url.Parse("https://deck.orbital.test:8443")
fmt.Println(u.Scheme, u.Host, u.Path == "", err) // https deck.orbital.test:8443 true <nil>
~~~

## أخطاء شائعة
~~~go
w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // بيعكس أي أصل
w.Header().Set("Access-Control-Allow-Origin", "*") // مع Credentials مرفوض
w.Header().Set("Vary", "Origin") // بيمسح Vary اللي قبله
~~~

## المهمة
نفّذ ~Shield(p Policy) (func(http.Handler) http.Handler,error)~. خد نسخة من القوايم؛ تعديل ~p~ بعدها ميأثرش. كل رد، حتى الرفض والـ preflight، فيه ~X-Content-Type-Options: nosniff~ و ~X-Frame-Options: DENY~ و ~Referrer-Policy: no-referrer~ و ~Content-Security-Policy: default-src 'none'~، وأول قيمة ~Vary~ هي ~Origin~ مضافة بـ ~Add~ قبل ~next~.

الطلب العادي، يعني أي طلب مش ~OPTIONS~ فيه ~Access-Control-Request-Method~، دايماً يروح لـ ~next~. لو ~Origin~ مطابق حرفياً لعنصر في ~Origins~ حط ~Access-Control-Allow-Origin~ بنفس القيمة، و ~Access-Control-Allow-Credentials: true~ لو ~Credentials~. غير كده من غير هيدرات CORS.

الـ preflight ميوصلش لـ ~next~، وقيم ~Vary~ فيه بالترتيب ~Origin~ و ~Access-Control-Request-Method~ و ~Access-Control-Request-Headers~. يتقبل لو الأصل مسموح، والطريقة المطلوبة موجودة حرفياً في ~Methods~، وكل اسم في ~Access-Control-Request-Headers~ بعد التقسيم بالفاصلة والـ trim وتجاهل الفاضي موجود في ~Headers~ من غير حالة. القبول 204 بجسم فاضي، و ~Access-Control-Allow-Origin~ بالأصل، و ~Access-Control-Allow-Methods~ بـ ~Methods~ بترتيبها مفصولة بـ ~, ~، و ~Access-Control-Allow-Headers~ بنفس الشكل لو ~Headers~ مش فاضية، و ~Access-Control-Max-Age~ لو ~MaxAge~ أكبر من صفر، و ~Access-Control-Allow-Credentials: true~ لو ~Credentials~. الرفض 403 و ~text/plain; charset=utf-8~ و ~preflight denied\n~ من غير أي هيدر ~Access-Control-Allow-*~. مفيش ~Report~، و ~main~ فاضية.

النجمة التالتة: ~Shield~ ترجع ~nil,error~ لسياسة مرخية: أصل مش بالشكل ~http(s)://host[:port]~ بحروف صغيرة من غير مسار أو query أو مستخدم (زي ~*~ و ~null~)، أو طريقة فاضية أو فيها غير الحروف الكبيرة، أو هيدر فاضي أو فيه مسافة أو فاصلة، أو ~MaxAge~ برة 0 لـ 86400.`,
		Hints: []string{"عامل CORS كسياسة قراءة يطبقها المتصفح.", "افصل preflight عن الطلبات العادية داخل wrapper.", "استخدم ~Header.Add~ لحفظ ~Vary~ الموجود."},
		En: &LevelText{LessonRevision: 1, Steps: ch27Lesson273.english, Title: "CORS", Goal: "Let only known fleet displays read base data, and put protection headers on every response", Files: []LabFile{{Name: "main.go", Starter: ch273Starter}},
			Lesson: `# A CORS Shield for the Gateway

## Before you write
- If a page from another site requests base data, who stops it from reading the response: the server or the browser?
- If a response changes with ~Origin~, how does a cache in between know?

## The origin
An origin is scheme, host, and port together, such as ~https://deck.orbital.test:8443~, without a path. A browser stops a page from one origin reading a response from another unless the response allows it. The server still receives the request; CORS is a browser read decision, not a replacement for authentication.

The browser puts the origin in the ~Origin~ header. If the origin is on the list, the response returns ~Access-Control-Allow-Origin~ with exactly that value. Matching is literal: ~https://Deck...~ or ~https://deck.../~ is not the same origin.
~~~go
origin := r.Header.Get("Origin")
if slices.Contains(allowed, origin) {
 w.Header().Set("Access-Control-Allow-Origin", origin)
}
~~~
~slices.Contains~ reports whether an element is present. A foreign origin is not rejected with an error; the response continues without CORS headers and the browser blocks the read.

## Cookies with the request
~Access-Control-Allow-Credentials: true~ lets the browser send cookies and read the response. With credentials ~*~ is forbidden; the exact origin is required. That is why this allowlist holds specific origins only.

## The preflight
Before a non-simple request, such as ~PATCH~ or a custom header, the browser sends ~OPTIONS~ with ~Access-Control-Request-Method~ and possibly ~Access-Control-Request-Headers~. It is a question, so the real handler does not run.
~~~go
preflight := r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != ""
~~~
Approval is a bodyless 204 with ~Access-Control-Allow-Methods~ and ~Access-Control-Allow-Headers~ listing values separated by ~, ~, plus ~Access-Control-Max-Age~ in seconds so the browser caches the approval. ~strings.Join~ joins a list and ~strconv.Itoa~ writes the number.

Header names are case-insensitive. ~strings.EqualFold~ compares ignoring case, and ~slices.ContainsFunc~ searches with a function:
~~~go
ok := slices.ContainsFunc(headers, func(a string) bool { return strings.EqualFold(a, "x-crew") })
~~~
Methods are case-sensitive: ~get~ is not ~GET~.

## Vary
A cache stores responses by URL. If the response differs by ~Origin~, ~Vary: Origin~ is required, or the cache hands one origin's response to another. ~Header().Add~ appends a value beside existing ones, while ~Set~ replaces them; use ~Add~ so a ~Vary~ from the handler survives. ~Header().Values~ returns every value.
~~~go
w.Header().Add("Vary", "Origin")
w.Header().Add("Vary", "Accept-Encoding")
fmt.Println(w.Header().Values("Vary")) // [Origin Accept-Encoding]
~~~

## Protection headers
~X-Content-Type-Options: nosniff~ stops content-type guessing. ~X-Frame-Options: DENY~ forbids framing the page. ~Referrer-Policy: no-referrer~ sends no address to other sites. ~Content-Security-Policy: default-src 'none'~ loads no resources from an API response. Set them before calling ~next~ so every response carries them.

## Parsing an origin
~url.Parse~ from ~net/url~ splits text into ~Scheme~, ~Host~, ~User~, ~Path~, and ~RawQuery~. A proper origin is reproduced exactly by ~Scheme + "://" + Host~.
~~~go
u, err := url.Parse("https://deck.orbital.test:8443")
fmt.Println(u.Scheme, u.Host, u.Path == "", err) // https deck.orbital.test:8443 true <nil>
~~~

## Common mistakes
~~~go
w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // Reflects any origin
w.Header().Set("Access-Control-Allow-Origin", "*") // Forbidden with credentials
w.Header().Set("Vary", "Origin") // Erases earlier Vary values
~~~

## Your mission
Implement ~Shield(p Policy) (func(http.Handler) http.Handler,error)~. Copy the lists; later changes to ~p~ have no effect. Every response, including denials and preflights, has ~X-Content-Type-Options: nosniff~, ~X-Frame-Options: DENY~, ~Referrer-Policy: no-referrer~, and ~Content-Security-Policy: default-src 'none'~, and its first ~Vary~ value is ~Origin~, added with ~Add~ before ~next~.

An ordinary request, meaning anything other than ~OPTIONS~ with ~Access-Control-Request-Method~, always reaches ~next~. If ~Origin~ literally equals an entry of ~Origins~, set ~Access-Control-Allow-Origin~ to it, plus ~Access-Control-Allow-Credentials: true~ when ~Credentials~. Otherwise add no CORS headers.

A preflight never reaches ~next~, and its ~Vary~ values are, in order, ~Origin~, ~Access-Control-Request-Method~, and ~Access-Control-Request-Headers~. It is accepted when the origin is allowed, the requested method is literally in ~Methods~, and every name in ~Access-Control-Request-Headers~, after splitting on commas, trimming, and skipping empties, is in ~Headers~ ignoring case. Acceptance is a bodyless 204 with ~Access-Control-Allow-Origin~ set to the origin, ~Access-Control-Allow-Methods~ listing ~Methods~ in order separated by ~, ~, ~Access-Control-Allow-Headers~ in the same form when ~Headers~ is nonempty, ~Access-Control-Max-Age~ when ~MaxAge~ is above zero, and ~Access-Control-Allow-Credentials: true~ when ~Credentials~. Denial is 403, ~text/plain; charset=utf-8~, and ~preflight denied\n~ with no ~Access-Control-Allow-*~ header. No ~Report~, and ~main~ stays empty.

Third star: ~Shield~ returns ~nil,error~ for a loose policy: an origin not shaped ~http(s)://host[:port]~ in lowercase without path, query, or user (such as ~*~ and ~null~), a method that is empty or has anything but uppercase letters, a header that is empty or contains a space or comma, or ~MaxAge~ outside 0 to 86400.`,
			Hints: []string{"Treat CORS as a browser-enforced read policy.", "Separate preflight from ordinary requests inside the wrapper.", "Use ~Header.Add~ to preserve existing ~Vary~ values."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch273Hidden},
	})
}

const ch273Starter = `package main

import "net/http"

type Policy struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool
	MaxAge      int
}

func Shield(p Policy) (func(http.Handler) http.Handler, error) {
	return func(next http.Handler) http.Handler { return next }, nil
}
func main() {}
`

const ch273Hidden = `package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type shieldCall struct {
	w       *httptest.ResponseRecorder
	reached bool
}

func shieldDo(t *testing.T, p Policy, method, origin string, extra map[string]string) shieldCall {
	t.Helper()
	mw, err := Shield(p)
	if err != nil || mw == nil {
		t.Fatal("valid policy rejected", err)
	}
	var c shieldCall
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.reached = true
		w.Header().Add("Vary", "Accept-Encoding")
		w.Write([]byte("cargo"))
	}))
	r := httptest.NewRequest(method, "/manifest", nil)
	if origin != "" {
		r.Header.Set("Origin", origin)
	}
	for k, v := range extra {
		r.Header.Set(k, v)
	}
	c.w = httptest.NewRecorder()
	h.ServeHTTP(c.w, r)
	return c
}
func shieldSecure(t *testing.T, c shieldCall) {
	t.Helper()
	h := c.w.Header()
	if h.Get("X-Content-Type-Options") != "nosniff" || h.Get("X-Frame-Options") != "DENY" || h.Get("Referrer-Policy") != "no-referrer" || h.Get("Content-Security-Policy") != "default-src 'none'" {
		t.Fatalf("security headers %v", h)
	}
}

// Boundaries: no Origin, exact-match allowlist, case variants, credentials on and off, Vary merging with the handler.
func TestHiddenShieldActualRequests(t *testing.T) {
	rng := rand.New(rand.NewSource(2731))
	for round := 0; round < 6; round++ {
		allowed := fmt.Sprintf("https://deck-%d.orbital.test:%d", rng.Intn(1000), 1024+rng.Intn(5000))
		other := fmt.Sprintf("https://deck-%d.orbital.test", 1000+rng.Intn(1000))
		p := Policy{Origins: []string{other + ":1", allowed}, Methods: []string{"GET", "PUT"}, Headers: []string{"X-Crew"}, Credentials: round%2 == 0, MaxAge: 60 * round}
		creds := ""
		if p.Credentials {
			creds = "true"
		}
		for _, method := range []string{"GET", "POST", "OPTIONS"} {
			c := shieldDo(t, p, method, "", nil)
			shieldSecure(t, c)
			if !c.reached || c.w.Header().Get("Access-Control-Allow-Origin") != "" || c.w.Body.String() != "cargo" {
				t.Fatal("no-origin request", method, c.w.Header())
			}
			if !reflect.DeepEqual(c.w.Header().Values("Vary"), []string{"Origin", "Accept-Encoding"}) {
				t.Fatal("vary", c.w.Header().Values("Vary"))
			}
			c = shieldDo(t, p, method, allowed, nil)
			shieldSecure(t, c)
			if !c.reached || c.w.Header().Get("Access-Control-Allow-Origin") != allowed || c.w.Header().Get("Access-Control-Allow-Credentials") != creds {
				t.Fatal("allowed origin", method, c.w.Header())
			}
			if !reflect.DeepEqual(c.w.Header().Values("Vary"), []string{"Origin", "Accept-Encoding"}) {
				t.Fatal("vary", c.w.Header().Values("Vary"))
			}
			for _, bad := range []string{other, allowed + "/", "HTTPS" + allowed[5:], "null", "http" + allowed[5:]} {
				c = shieldDo(t, p, method, bad, nil)
				shieldSecure(t, c)
				if !c.reached || c.w.Header().Get("Access-Control-Allow-Origin") != "" || c.w.Header().Get("Access-Control-Allow-Credentials") != "" {
					t.Fatal("foreign origin got CORS", bad, c.w.Header())
				}
			}
		}
	}
}

// Boundaries: allowed and denied preflights by origin, method and header; header case and spacing; Max-Age zero; policy snapshot.
func TestHiddenShieldPreflight(t *testing.T) {
	rng := rand.New(rand.NewSource(2732))
	origin := fmt.Sprintf("https://bridge-%d.orbital.test", rng.Intn(9999))
	p := Policy{Origins: []string{origin}, Methods: []string{"GET", "POST", "PATCH"}, Headers: []string{"Content-Type", "X-Crew"}, Credentials: true, MaxAge: 600}
	mw, err := Shield(p)
	if err != nil {
		t.Fatal(err)
	}
	p.Origins[0] = "https://evil.test"
	p.Methods[0] = "DELETE"
	p.Headers[0] = "X-Evil"
	run := func(o string, extra map[string]string) shieldCall {
		var c shieldCall
		h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { c.reached = true }))
		r := httptest.NewRequest("OPTIONS", "/manifest", nil)
		r.Header.Set("Origin", o)
		for k, v := range extra {
			r.Header.Set(k, v)
		}
		c.w = httptest.NewRecorder()
		h.ServeHTTP(c.w, r)
		return c
	}
	preVary := []string{"Origin", "Access-Control-Request-Method", "Access-Control-Request-Headers"}
	for _, extra := range []map[string]string{
		{"Access-Control-Request-Method": "GET"},
		{"Access-Control-Request-Method": "PATCH", "Access-Control-Request-Headers": "x-crew , CONTENT-TYPE,,"},
		{"Access-Control-Request-Method": "POST", "Access-Control-Request-Headers": "Content-Type"},
	} {
		c := run(origin, extra)
		shieldSecure(t, c)
		h := c.w.Header()
		if c.reached || c.w.Code != 204 || c.w.Body.Len() != 0 || h.Get("Access-Control-Allow-Origin") != origin || h.Get("Access-Control-Allow-Methods") != "GET, POST, PATCH" ||
			h.Get("Access-Control-Allow-Headers") != "Content-Type, X-Crew" || h.Get("Access-Control-Max-Age") != "600" || h.Get("Access-Control-Allow-Credentials") != "true" || !reflect.DeepEqual(h.Values("Vary"), preVary) {
			t.Fatalf("preflight %v: %d %v", extra, c.w.Code, h)
		}
	}
	for label, tc := range map[string]struct {
		origin string
		extra  map[string]string
	}{
		"foreign":     {"https://evil.test", map[string]string{"Access-Control-Request-Method": "GET"}},
		"method":      {origin, map[string]string{"Access-Control-Request-Method": "DELETE"}},
		"lower":       {origin, map[string]string{"Access-Control-Request-Method": "get"}},
		"header":      {origin, map[string]string{"Access-Control-Request-Method": "GET", "Access-Control-Request-Headers": "X-Crew, X-Evil"}},
		"no origin":   {"", map[string]string{"Access-Control-Request-Method": "GET"}},
		"origin case": {"HTTPS" + origin[5:], map[string]string{"Access-Control-Request-Method": "GET"}},
	} {
		c := run(tc.origin, tc.extra)
		shieldSecure(t, c)
		h := c.w.Header()
		if c.reached || c.w.Code != 403 || c.w.Body.String() != "preflight denied\n" || h.Get("Content-Type") != "text/plain; charset=utf-8" ||
			h.Get("Access-Control-Allow-Origin") != "" || h.Get("Access-Control-Allow-Methods") != "" || !reflect.DeepEqual(h.Values("Vary"), preVary) {
			t.Fatalf("%s preflight allowed: %d %v", label, c.w.Code, h)
		}
	}
	if c := run(origin, nil); !c.reached {
		t.Fatal("OPTIONS without Access-Control-Request-Method is an ordinary request")
	}
	quiet := shieldDo(t, Policy{Origins: []string{origin}, Methods: []string{"GET"}}, "OPTIONS", origin, map[string]string{"Access-Control-Request-Method": "GET"})
	h := quiet.w.Header()
	if quiet.w.Code != 204 || h.Get("Access-Control-Max-Age") != "" || h.Get("Access-Control-Allow-Credentials") != "" || h.Get("Access-Control-Allow-Headers") != "" || h.Get("Access-Control-Allow-Methods") != "GET" {
		t.Fatal("optional preflight headers", h)
	}
}

func TestHiddenBonusShieldRejectsLoosePolicies(t *testing.T) {
	for _, p := range []Policy{
		{Origins: []string{"*"}},
		{Origins: []string{"null"}},
		{Origins: []string{"https://deck.orbital.test/"}},
		{Origins: []string{"https://deck.orbital.test/panel"}},
		{Origins: []string{"https://deck.orbital.test?x=1"}},
		{Origins: []string{"deck.orbital.test"}},
		{Origins: []string{"ftp://deck.orbital.test"}},
		{Origins: []string{"https://crew@deck.orbital.test"}},
		{Origins: []string{"https://Deck.orbital.test"}},
		{Origins: []string{""}},
		{Methods: []string{"get"}},
		{Methods: []string{""}},
		{Headers: []string{"X Crew"}},
		{Headers: []string{""}},
		{MaxAge: -1},
		{MaxAge: 86401},
	} {
		if mw, err := Shield(p); mw != nil || err == nil {
			t.Fatalf("loose policy accepted: %+v", p)
		}
	}
	for _, p := range []Policy{{}, {Origins: []string{"http://localhost:8080", "https://deck.orbital.test"}, Methods: []string{"GET"}, Headers: []string{"X-Crew"}, MaxAge: 86400}} {
		if mw, err := Shield(p); mw == nil || err != nil {
			t.Fatalf("strict policy rejected: %+v %v", p, err)
		}
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch27Lesson274.arabic, ID: "27-4-rate-limit", Chapter: 27, Kind: "lab", Title: "حد الطلبات", Goal: "ادّي كل عميل خزان تصاريح خاص بيه، وارفض الزيادة بـ 429، ونضّف العملاء الخاملين",
		Requires: []string{"mutex"}, Bridge: []int{57}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch274Starter}},
		Lesson: `# خزان تصاريح لكل عميل

## قبل ما تكتب
- لو كل الطلبات بتيجي من proxy واحد، إزاي تفرّق العملاء من غير ما تصدّق أي هيدر من أي حد؟
- إمتى تقدر تنسى عميل من غير ما تديله تصاريح زيادة؟

## من خزان واحد لخزانات كتير
في 25-3 بنيت ~PermitTank~: يبدأ مليان، وكل ~every~ كاملة تضيف تصريح، والقسمة الصحيحة على المدة بتحتفظ بباقي الفترة. هنا نفس الفكرة لكن خزان لكل عميل في خريطة: ~map[string]*bucket~. أول طلب من عميل جديد يعمل له خزان مليان.
~~~go
b, ok := tanks[key]
if !ok { b = &bucket{credit: burst, mark: now}; tanks[key] = b }
~~~
الخريطة والخزانات كلها تحت ~sync.Mutex~ واحد: البحث والإنشاء والتجديد والاستهلاك قرار واحد. كل طلب بيكلّف تصريح واحد.

## الرد بـ 429
~http.StatusTooManyRequests~ هي 429. هيدر ~Retry-After~ بيقول للعميل يستنى كام ثانية صحيحة. الانتظار هو الباقي لحد التصريح الجاي: ~mark + every - now~. قرّبه لفوق بالقسمة الصحيحة على المدد:
~~~go
wait := mark.Add(every).Sub(now) // 1.2s
secs := (wait + time.Second - 1) / time.Second // 2
w.Header().Set("Retry-After", strconv.Itoa(int(secs)))
~~~
الطلب المرفوض ميوصلش للـ handler، وميستهلكش حاجة. الطلب المقبول يعرف العميل باقي له كام بهيدر ~X-RateLimit-Remaining~.

## مين العميل؟
~r.RemoteAddr~ عنوان الطرف التاني للاتصال، زي ~198.51.100.7:4242~ أو ~[2001:db8::1]:80~. ~net.SplitHostPort~ تفصل المضيف عن البورت وتشيل الأقواس. البورت بيتغير مع كل اتصال، فالمفتاح هو المضيف بس.
~~~go
host, port, err := net.SplitHostPort("[2001:db8::1]:80") // "2001:db8::1" "80" nil
ip := net.ParseIP(host)
fmt.Println(ip.String()) // 2001:db8::1
~~~
~net.ParseIP~ ترجع ~nil~ لو النص مش IP. ~ip.String()~ بتكتبه بشكل قياسي، فـ ~2001:db8:0:0::1~ و ~2001:db8::1~ يبقوا نفس المفتاح.

## X-Forwarded-For
الـ proxy بيحط عنوان العميل اللي شافه في آخر ~X-Forwarded-For~، مفصول بفاصلة عن القيم اللي قبله. أي عميل يقدر يكتب الهيدر ده بنفسه، فلو صدّقته من أي حد، كل طلب يقدر يغيّر هويته ويهرب من الحد. صدّقه بس لو ~RemoteAddr~ نفسه proxy معروف، وخُد آخر عنصر بس، لأنه اللي الـ proxy الموثوق كتبه.
~~~go
values := r.Header.Values("X-Forwarded-For") // ممكن أكتر من سطر
hops := strings.Split(strings.Join(values, ","), ",")
last := strings.TrimSpace(hops[len(hops)-1])
~~~
لو آخر عنصر مش IP، الـ proxy نفسه بايظ؛ ارفض بدل ما تخمّن.

## تنظيف الذاكرة
كل IP جديد بيعمل خزان، فالخريطة ممكن تكبر لحد ما الذاكرة تخلص. احفظ آخر وقت شفت فيه العميل، وامسح اللي خامل من مدة ~Idle~ أو أكتر. العميل اللي اتمسح يرجع بخزان مليان؛ ده آمن بس لو ~Idle~ مش أقل من ~Burst*Every~، لأن الخزان كان هيتملا في المدة دي كده كده.
~~~go
for key, b := range tanks {
 if now.Sub(b.seen) >= idle { delete(tanks, key) }
}
~~~
المسح جوه ~range~ على نفس الخريطة مسموح في Go.

## أخطاء شائعة
~~~go
key := r.RemoteAddr // البورت جزء من المفتاح: كل اتصال خزان جديد
key = r.Header.Get("X-Forwarded-For") // أي عميل يزوّر هويته
mu.Lock(); b := tanks[key]; mu.Unlock(); b.credit-- // استهلاك برة القفل
~~~

## المهمة
نفّذ ~NewLimiter(clock labkit.Clock, l Limits) (*Limiter,error)~. ارفض بـ ~nil,error~: ساعة ~nil~، أو ~Burst~ برة 1 لـ 1000، أو ~Every~ غير موجبة أو أكبر من ساعة، أو ~Idle~ أقل من ~Burst*Every~، أو أي عنصر في ~Proxies~ مش IP.

~(l *Limiter) Wrap(next http.Handler) http.Handler~ يحسب مفتاح العميل: المضيف من ~RemoteAddr~ بشكل ~ip.String()~. لو المضيف ده في ~Proxies~ والطلب فيه ~X-Forwarded-For~، اجمع كل السطور بفاصلة وخُد آخر عنصر بعد trim بشكله القياسي. من غير proxy موثوق تجاهل الهيدر. ~RemoteAddr~ أو آخر عنصر مش صالح: 400 و ~bad client\n~ من غير استهلاك.

كل عميل له خزان بسعة ~Burst~ يبدأ مليان، وتصريح لكل ~Every~ كاملة بنفس قواعد 25-3، وكل طلب بتصريح. القبول يحط ~X-RateLimit-Remaining~ بالباقي بعد الاستهلاك وينادي ~next~. الرفض 429 و ~Retry-After~ بالثواني لحد التصريح الجاي متقرّبة لفوق، و ~slow down\n~، من غير ~next~. الأخطاء ~text/plain; charset=utf-8~.

كل طلب له مفتاح صالح يسجّل وقت آخر ظهور حتى لو اترفض. ~(l *Limiter) Sweep() int~ يمسح العملاء اللي عدّى على آخر ظهور ليهم ~Idle~ أو أكتر ويرجع عددهم، و ~(l *Limiter) Clients() int~ يرجع عدد الخزانات. كله آمن للتزامن تحت قفل واحد؛ مفيش نوم أو وقت حقيقي أو goroutines أو ~Report~، و ~main~ فاضية.`,
		Hints: []string{"خلي لكل عميل معروف رصيدًا وحالة تجديد خاصة.", "اجمع تعديل الخريطة والخزان تحت mutex نفسه.", "ابدأ بتحليل ~SplitHostPort~ ثم ~ParseIP~."},
		En: &LevelText{LessonRevision: 1, Steps: ch27Lesson274.english, Title: "Rate Limit", Goal: "Give each client its own permit tank, refuse excess with 429, and clean up idle clients", Files: []LabFile{{Name: "main.go", Starter: ch274Starter}},
			Lesson: `# A Permit Tank per Client

## Before you write
- If every request arrives through one proxy, how do you tell clients apart without believing any header from anyone?
- When can you forget a client without handing it extra permits?

## From one tank to many
In 25-3 you built ~PermitTank~: it starts full, each complete ~every~ adds a permit, and integer division by the duration keeps the leftover part of an interval. The same idea now runs one tank per client in a map: ~map[string]*bucket~. A new client's first request creates a full tank.
~~~go
b, ok := tanks[key]
if !ok { b = &bucket{credit: burst, mark: now}; tanks[key] = b }
~~~
The map and all tanks share one ~sync.Mutex~: lookup, creation, refill, and spending are one decision. Every request costs one permit.

## Answering with 429
~http.StatusTooManyRequests~ is 429. The ~Retry-After~ header tells the client how many whole seconds to wait. The wait is what remains until the next permit: ~mark + every - now~. Round it up with integer division on durations:
~~~go
wait := mark.Add(every).Sub(now) // 1.2s
secs := (wait + time.Second - 1) / time.Second // 2
w.Header().Set("Retry-After", strconv.Itoa(int(secs)))
~~~
A refused request never reaches the handler and spends nothing. An accepted request tells the client what is left with an ~X-RateLimit-Remaining~ header.

## Who is the client?
~r.RemoteAddr~ is the other end of the connection, such as ~198.51.100.7:4242~ or ~[2001:db8::1]:80~. ~net.SplitHostPort~ separates host from port and strips brackets. The port changes per connection, so the key is the host alone.
~~~go
host, port, err := net.SplitHostPort("[2001:db8::1]:80") // "2001:db8::1" "80" nil
ip := net.ParseIP(host)
fmt.Println(ip.String()) // 2001:db8::1
~~~
~net.ParseIP~ returns ~nil~ for text that is not an IP. ~ip.String()~ writes the canonical form, so ~2001:db8:0:0::1~ and ~2001:db8::1~ become the same key.

## X-Forwarded-For
A proxy appends the client address it saw to the end of ~X-Forwarded-For~, comma-separated from earlier values. Any client can write this header itself, so believing it from anyone lets every request change identity and escape the limit. Believe it only when ~RemoteAddr~ itself is a known proxy, and take only the last item, since that is what the trusted proxy wrote.
~~~go
values := r.Header.Values("X-Forwarded-For") // May be several lines
hops := strings.Split(strings.Join(values, ","), ",")
last := strings.TrimSpace(hops[len(hops)-1])
~~~
If the last item is not an IP, the proxy itself is broken; refuse rather than guess.

## Cleaning up memory
Every new IP creates a tank, so the map could grow until memory runs out. Record when you last saw each client and delete those idle for ~Idle~ or longer. A deleted client returns with a full tank; that is safe only when ~Idle~ is at least ~Burst*Every~, because the tank would have refilled in that time anyway.
~~~go
for key, b := range tanks {
 if now.Sub(b.seen) >= idle { delete(tanks, key) }
}
~~~
Deleting inside a ~range~ over the same map is allowed in Go.

## Common mistakes
~~~go
key := r.RemoteAddr // Port is part of the key: every connection gets a new tank
key = r.Header.Get("X-Forwarded-For") // Any client forges its identity
mu.Lock(); b := tanks[key]; mu.Unlock(); b.credit-- // Spending outside the lock
~~~

## Your mission
Implement ~NewLimiter(clock labkit.Clock, l Limits) (*Limiter,error)~. Return ~nil,error~ for a ~nil~ clock, ~Burst~ outside 1 to 1000, ~Every~ not positive or above one hour, ~Idle~ below ~Burst*Every~, or any ~Proxies~ entry that is not an IP.

~(l *Limiter) Wrap(next http.Handler) http.Handler~ computes the client key: the ~RemoteAddr~ host in ~ip.String()~ form. If that host is in ~Proxies~ and the request has ~X-Forwarded-For~, join all lines with commas and take the last item, trimmed, in canonical form. Without a trusted proxy ignore the header. An invalid ~RemoteAddr~ or last item returns 400 and ~bad client\n~ and spends nothing.

Each client has a tank of capacity ~Burst~ that starts full, earns one permit per complete ~Every~ under the 25-3 rules, and spends one per request. Acceptance sets ~X-RateLimit-Remaining~ to what is left after spending and calls ~next~. Refusal returns 429, ~Retry-After~ in seconds until the next permit rounded up, and ~slow down\n~, without ~next~. Errors use ~text/plain; charset=utf-8~.

Every request with a valid key records its last-seen time, even when refused. ~(l *Limiter) Sweep() int~ deletes clients with ~now - last seen >= Idle~ and returns how many, and ~(l *Limiter) Clients() int~ returns the number of tanks. Everything is concurrency-safe under one lock; no sleeping, real time, goroutines, or ~Report~, and ~main~ stays empty.`,
			Hints: []string{"Give each known client its own balance and refill state.", "Guard map and tank changes with the same mutex.", "Start by parsing ~SplitHostPort~ and then ~ParseIP~."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch274Hidden},
	})
}

const ch274Starter = `package main

import (
	"errors"
	"net/http"
	"time"

	"orbital/colony/labkit"
)

type Limits struct {
	Burst   int
	Every   time.Duration
	Idle    time.Duration
	Proxies []string
}

type Limiter struct{}

func NewLimiter(clock labkit.Clock, l Limits) (*Limiter, error) { return nil, errors.New("TODO") }
func (l *Limiter) Wrap(next http.Handler) http.Handler          { return next }
func (l *Limiter) Sweep() int                                   { return 0 }
func (l *Limiter) Clients() int                                 { return 0 }
func main()                                                     {}
`

const ch274Hidden = `package main

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

func limitCall(h http.Handler, remote string, forwarded ...string) *httptest.ResponseRecorder {
	r := httptest.NewRequest("GET", "/cargo", nil)
	r.RemoteAddr = remote
	for _, f := range forwarded {
		r.Header.Add("X-Forwarded-For", f)
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}
func limitPass(t *testing.T, w *httptest.ResponseRecorder, left int) {
	t.Helper()
	if w.Code != 200 || w.Body.String() != "ok" || w.Header().Get("X-RateLimit-Remaining") != strconv.Itoa(left) {
		t.Fatalf("want pass with %d left, got %d %q %v", left, w.Code, w.Body.String(), w.Header())
	}
}
func limitBlocked(t *testing.T, w *httptest.ResponseRecorder, retry int) {
	t.Helper()
	if w.Code != 429 || w.Body.String() != "slow down\n" || w.Header().Get("Retry-After") != strconv.Itoa(retry) || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("want 429 retry %d, got %d %q %v", retry, w.Code, w.Body.String(), w.Header())
	}
}
func limitHandler(reached *atomic.Int64) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached.Add(1)
		w.Write([]byte("ok"))
	})
}

// Boundaries: invalid limits, Idle exactly Burst*Every, proxy addresses canonicalised.
func TestHiddenLimiterRejectsBadLimits(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	for _, l := range []Limits{
		{Burst: 0, Every: time.Second, Idle: time.Hour},
		{Burst: 1001, Every: time.Second, Idle: 2000 * time.Hour},
		{Burst: 3, Every: 0, Idle: time.Hour},
		{Burst: 3, Every: -time.Second, Idle: time.Hour},
		{Burst: 3, Every: time.Second, Idle: 3*time.Second - 1},
		{Burst: 3, Every: time.Second, Idle: time.Hour, Proxies: []string{"10.0.0"}},
		{Burst: 3, Every: time.Second, Idle: time.Hour, Proxies: []string{"10.0.0.1:80"}},
	} {
		if lim, err := NewLimiter(clock, l); lim != nil || err == nil {
			t.Fatalf("accepted %+v", l)
		}
	}
	if lim, err := NewLimiter(nil, Limits{Burst: 1, Every: time.Second, Idle: time.Second}); lim != nil || err == nil {
		t.Fatal("nil clock")
	}
	if lim, err := NewLimiter(clock, Limits{Burst: 3, Every: time.Second, Idle: 3 * time.Second, Proxies: []string{"10.0.0.1", "2001:db8::7"}}); lim == nil || err != nil {
		t.Fatal("valid limits rejected", err)
	}
}

// Boundaries: full burst, the exact refill instant, retained remainders, Retry-After rounding, isolation, ports and IPv6 spelling.
func TestHiddenLimiterBurstRefillAndIsolation(t *testing.T) {
	rng := rand.New(rand.NewSource(2741))
	for round := 0; round < 5; round++ {
		burst := 2 + rng.Intn(5)
		every := time.Duration(700+rng.Intn(2300)) * time.Millisecond
		clock := labkit.NewFakeClock(time.Unix(1_700_000_000+int64(rng.Intn(1000)), 0))
		lim, err := NewLimiter(clock, Limits{Burst: burst, Every: every, Idle: time.Duration(burst) * every})
		if err != nil {
			t.Fatal(err)
		}
		var reached atomic.Int64
		h := lim.Wrap(limitHandler(&reached))
		a := fmt.Sprintf("198.51.100.%d", 1+rng.Intn(200))
		for i := burst - 1; i >= 0; i-- {
			limitPass(t, limitCall(h, a+":"+strconv.Itoa(1000+i)), i)
		}
		retry := int((every + time.Second - 1) / time.Second)
		limitBlocked(t, limitCall(h, a+":5"), retry)
		if reached.Load() != int64(burst) {
			t.Fatal("blocked request reached handler")
		}
		limitPass(t, limitCall(h, "203.0.113.9:7"), burst-1)
		clock.Advance(every - time.Nanosecond)
		limitBlocked(t, limitCall(h, a+":5"), 1)
		clock.Advance(time.Nanosecond)
		limitPass(t, limitCall(h, a+":5"), 0)
		clock.Advance(every + every/2)
		limitPass(t, limitCall(h, a+":5"), 0)
		limitBlocked(t, limitCall(h, a+":5"), int((every/2+time.Second-1)/time.Second))
		clock.Advance(every / 2)
		limitPass(t, limitCall(h, a+":5"), 0)
		clock.Advance(time.Duration(burst+5) * every)
		limitPass(t, limitCall(h, a+":5"), burst-1)
	}
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	lim, _ := NewLimiter(clock, Limits{Burst: 1, Every: time.Minute, Idle: time.Minute})
	var reached atomic.Int64
	h := lim.Wrap(limitHandler(&reached))
	limitPass(t, limitCall(h, "[2001:db8::1]:80"), 0)
	limitBlocked(t, limitCall(h, "[2001:db8:0:0::1]:9"), 60)
	for _, bad := range []string{"", "garbage", "198.51.100.1", "host.test:80"} {
		if w := limitCall(h, bad); w.Code != 400 || w.Body.String() != "bad client\n" {
			t.Fatalf("remote %q: %d", bad, w.Code)
		}
	}
	if reached.Load() != 1 {
		t.Fatal("handler calls", reached.Load())
	}
}

// Boundaries: spoofed headers from untrusted peers, last hop through a trusted proxy, several header lines, broken hops.
func TestHiddenLimiterProxyPolicy(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	lim, err := NewLimiter(clock, Limits{Burst: 1, Every: time.Hour, Idle: time.Hour, Proxies: []string{"10.0.0.5"}})
	if err != nil {
		t.Fatal(err)
	}
	var reached atomic.Int64
	h := lim.Wrap(limitHandler(&reached))
	limitPass(t, limitCall(h, "192.0.2.50:1", "1.1.1.1"), 0)
	limitBlocked(t, limitCall(h, "192.0.2.50:2", "2.2.2.2"), 3600)
	limitBlocked(t, limitCall(h, "192.0.2.50:3"), 3600)
	limitPass(t, limitCall(h, "10.0.0.5:9", "192.0.2.60"), 0)
	limitPass(t, limitCall(h, "10.0.0.5:9", "192.0.2.50, 203.0.113.70"), 0)
	limitBlocked(t, limitCall(h, "10.0.0.5:9", "203.0.113.70"), 3600)
	limitPass(t, limitCall(h, "10.0.0.5:9", "8.8.8.8", " 203.0.113.71 "), 0)
	limitBlocked(t, limitCall(h, "10.0.0.5:9", "203.0.113.71"), 3600)
	limitPass(t, limitCall(h, "10.0.0.5:9"), 0)
	limitBlocked(t, limitCall(h, "10.0.0.5:10"), 3600)
	for _, hop := range []string{"nonsense", "203.0.113.70,", "203.0.113.70:80"} {
		if w := limitCall(h, "10.0.0.5:9", hop); w.Code != 400 || w.Body.String() != "bad client\n" {
			t.Fatalf("hop %q: %d", hop, w.Code)
		}
	}
	if lim.Clients() != 5 {
		t.Fatal("tracked clients", lim.Clients())
	}
}

// Boundaries: Sweep at Idle minus one tick and exactly Idle, recent clients survive, counts returned.
func TestHiddenLimiterSweepForgetsIdleClients(t *testing.T) {
	rng := rand.New(rand.NewSource(2744))
	clock := labkit.NewFakeClock(time.Unix(1_800_000_000, 0))
	idle := time.Duration(10+rng.Intn(20)) * time.Second
	lim, _ := NewLimiter(clock, Limits{Burst: 2, Every: time.Second, Idle: idle})
	var reached atomic.Int64
	h := lim.Wrap(limitHandler(&reached))
	n := 5 + rng.Intn(10)
	for i := 0; i < n; i++ {
		limitCall(h, fmt.Sprintf("192.0.2.%d:1", i+1))
	}
	if lim.Clients() != n || lim.Sweep() != 0 {
		t.Fatal("fresh clients swept")
	}
	clock.Advance(idle - time.Second)
	limitCall(h, "192.0.2.1:1")
	clock.Advance(time.Second - time.Nanosecond)
	if lim.Sweep() != 0 {
		t.Fatal("swept before Idle")
	}
	clock.Advance(time.Nanosecond)
	if got := lim.Sweep(); got != n-1 || lim.Clients() != 1 {
		t.Fatal("sweep at Idle", got, lim.Clients())
	}
	clock.Advance(idle)
	if lim.Sweep() != 1 || lim.Clients() != 0 {
		t.Fatal("last client")
	}
	limitPass(t, limitCall(h, "192.0.2.1:1"), 1)
}

func TestHiddenLimiterParallelClients(t *testing.T) {
	clock := labkit.NewFakeClock(time.Unix(0, 0))
	lim, _ := NewLimiter(clock, Limits{Burst: 10, Every: time.Second, Idle: time.Minute})
	var reached atomic.Int64
	h := lim.Wrap(limitHandler(&reached))
	var wg sync.WaitGroup
	var passed, blocked [4]atomic.Int64
	start := make(chan struct{})
	for g := 0; g < 64; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			<-start
			for i := 0; i < 4; i++ {
				w := limitCall(h, fmt.Sprintf("192.0.2.%d:%d", 10+i, g))
				if w.Code == 200 {
					passed[i].Add(1)
				} else if w.Code == 429 {
					blocked[i].Add(1)
				}
			}
			lim.Sweep()
			lim.Clients()
		}(g)
	}
	close(start)
	wg.Wait()
	for i := range passed {
		if passed[i].Load() != 10 || blocked[i].Load() != 54 {
			t.Fatal("client", i, passed[i].Load(), blocked[i].Load())
		}
	}
}
`

func init() {
	add(&Level{LessonRevision: 1, Steps: ch27Lesson275.arabic, ID: "27-5-short-links", Chapter: 27, Kind: "lab", Title: "روابط مختصرة", Goal: "اختصر روابط الأسطول بأكواد Base62 فوق مخزن بواجهة، واتعامل مع التصادم والروابط البايظة",
		Requires: []string{"iface"}, Bridge: []int{53}, Setup: seededLabRoom, Race: true,
		Files: []LabFile{{Name: "main.go", Starter: ch275Starter}},
		Lesson: `# أكواد قصيرة لروابط الأسطول

## قبل ما تكتب
- لو مولّد الأرقام طلّع رقم اتستخدم قبل كده، مين يكتشف ده: الـ handler ولا المخزن؟
- رابط ~javascript:~ أو رابط فيه اسم مستخدم، إيه ضرره لو اتحوّل له حد؟

## Base62
أي رقم ينكتب بأي أساس. في الأساس 10 كل خانة من 0 لـ 9؛ في Base62 كل خانة من 62 حرف: الأرقام ثم الحروف الصغيرة ثم الكبيرة. النتيجة كود قصير يتكتب في URL من غير ترميز.
~~~go
const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
fmt.Println(digits[61 % 62]) // 90، يعني 'Z'
~~~
التحويل: باقي القسمة على 62 هو آخر خانة، والقسمة الصحيحة تشيلها. كرر لحد ما الرقم يبقى صفر، وبعدين اعكس الخانات. الصفر نفسه خانة واحدة ~0~.
~~~go
n := uint64(125) // 125 = 2*62 + 1
fmt.Println(n%62, n/62) // 1 2 → الكود "21"
~~~
أكبر ~uint64~ بياخد 11 خانة. ~uint64~ نوع صحيح من غير إشارة، و ~math.MaxUint64~ أكبر قيمة فيه.

## المخزن ورا واجهة
~LinkStore~ واجهة زي فصل البروتوكولات: ~Put~ و ~Get~ بس. الـ handler ميعرفش إذا المخزن خريطة في الذاكرة ولا قاعدة بيانات، والتست يقدر يحقن مخزن مزيّف يعمل تصادم أو عطل عند الطلب.
~~~go
type Shelf interface { Put(code, target string) error }
var s Shelf = &memShelf{items: map[string]string{}}
~~~
المخزن هو اللي يحكم على التكرار، لأنه الوحيد اللي شايف كل الأكواد تحت قفله. ~ErrTaken~ و ~ErrNoLink~ أخطاء ثابتة (sentinel)؛ ممكن ترجع ملفوفة بـ ~%w~، فافحصها بـ ~errors.Is~ مش بـ ~==~.
~~~go
if errors.Is(err, ErrTaken) { continue } // جرّب رقم تاني
~~~

## التصادم والمحاولات
الأرقام جاية من دالة محقونة ~next~، في التشغيل ممكن تكون عشوائية. التصادم وارد، فجرّب عدد محدود من المرات. المحاولات اللانهائية ممكن تعلّق الطلب لو المخزن اتملى. أي خطأ تاني غير التصادم عطل حقيقي، متعيدش عليه.

متسحبش رقم قبل ما تتأكد إن الطلب سليم؛ الطلب البايظ ميستهلكش حاجة.

## الرابط الهدف
~url.Parse~ من درس CORS بتفك الرابط. اقبل ~http~ و ~https~ بس: ~javascript:~ و ~data:~ بيشغّلوا كود في متصفح الضحية. لازم ~Host~ مش فاضي، و ~User~ يكون ~nil~ لأن ~https://crew@evil.test~ بيضلل القارئ.
~~~go
u, err := url.Parse("https://dock.orbital.test/a?b=1")
ok := err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil
~~~
ارفض كمان أي بايت تحكم أو مسافة (أي بايت ~<= ' '~ أو 127)، لأنها ممكن تكسر الهيدر أو تخدع العين.

## التحويل
302 (~http.StatusFound~) معناها «الحاجة دي في مكان تاني»، وهيدر ~Location~ بالعنوان. ~http.Redirect~ بتضيف جسم HTML صغير لطلبات GET، فهنا حط الهيدر والحالة بنفسك عشان الجسم يفضل فاضي. 503 (~http.StatusServiceUnavailable~) معناها السيرفر مش قادر دلوقتي.
~~~go
w.Header().Set("Location", target)
w.WriteHeader(http.StatusFound)
~~~

## أخطاء شائعة
~~~go
code := strconv.Itoa(len(links)+1) // أكواد بتكشف عدد الروابط وبتتكرر بعد الحذف
for { if store.Put(code, t) == nil { break } } // لف لا نهائي، وبيعيد حتى على العطل
if err == ErrTaken { } // الخطأ الملفوف بيعدّي
~~~

## المهمة
نفّذ ~Base62(n uint64) string~ بالأبجدية ~0-9~ ثم ~a-z~ ثم ~A-Z~ من غير أصفار في الأول، و ~Base62(0)~ هي ~0~. نفّذ ~NewMemoryStore() LinkStore~ آمن للتزامن ومستقل: ~Put~ لكود موجود ترجع ~ErrTaken~ ومتكتبش فوقه، و ~Get~ لكود غايب ترجع ~ErrNoLink~.

نفّذ ~NewLinkDesk(store LinkStore, next func() uint64) http.Handler~ (المدخلين مش ~nil~). ~POST /links~ بجسم حده 2048 بايت فيه ~{"target":"..."}~. JSON بايظ أو أكبر من الحد أو هدف مش ~http~/~https~ أو من غير host أو فيه مستخدم أو فيه بايت ~<= ' '~ أو 127: 400 و ~bad target\n~ من غير نداء ~next~. بعدها حتى 3 محاولات: الكود ~Base62(next())~ و ~store.Put~. ~ErrTaken~ (حتى ملفوف) جرّب تاني؛ بعد 3 تصادمات 503 و ~no free code\n~. أي خطأ تاني 500 و ~store failure\n~ من غير إعادة. النجاح 201 و ~application/json~ و ~{"code":"C","path":"/go/C"}~ و ~Location: /go/C~.

~GET /go/C~ حيث C من 1 لـ 11 حرف من الأبجدية: ~store.Get~؛ ~ErrNoLink~ (حتى ملفوف) 404 و ~no link\n~، وخطأ تاني 500 و ~store failure\n~، والنجاح 302 و ~Location~ بالهدف وجسم فاضي. الـ query متأثرش. كود بشكل غلط أو أي مسار تاني: 404 و ~no link\n~ من غير ما تسأل المخزن. طريقة غلط على ~/links~ أو ~/go/C~: 405 و ~Allow: POST~ أو ~Allow: GET~ و ~wrong method\n~. الأخطاء ~text/plain; charset=utf-8~. مفيش ~Report~، و ~main~ فاضية.

النجمة التالتة: ~ParseBase62(s string) (uint64,error)~ عكس ~Base62~ بالظبط: ترفض النص الفاضي، وأي حرف برة الأبجدية، والصفر في الأول لأكتر من خانة، وأي قيمة أكبر من ~math.MaxUint64~.`,
		Hints: []string{"ثبّت ترتيب الأبجدية قبل التفكير في تمثيل الرقم.", "خلّي واجهة التخزين مسؤولة عن حفظ وربط الكود.", "استخدم ~errors.Is~ للتعرف على خطأ تصادم ملفوف."},
		En: &LevelText{LessonRevision: 1, Steps: ch27Lesson275.english, Title: "Short Links", Goal: "Shorten fleet links into Base62 codes over an interface store, handling collisions and broken links", Files: []LabFile{{Name: "main.go", Starter: ch275Starter}},
			Lesson: `# Short Codes for Fleet Links

## Before you write
- If the number source returns a number already used, who detects it: the handler or the store?
- What harm does a ~javascript:~ link, or a link carrying a user name, do to someone redirected to it?

## Base62
Any number can be written in any base. In base 10 every digit is 0 to 9; in Base62 every digit is one of 62 characters: digits, then lowercase, then uppercase letters. The result is a short code that fits in a URL without escaping.
~~~go
const digits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
fmt.Println(digits[61 % 62]) // 90, that is 'Z'
~~~
Conversion: the remainder modulo 62 is the last digit, and integer division removes it. Repeat until the number is zero, then reverse the digits. Zero itself is the single digit ~0~.
~~~go
n := uint64(125) // 125 = 2*62 + 1
fmt.Println(n%62, n/62) // 1 2 -> code "21"
~~~
The largest ~uint64~ takes 11 digits. ~uint64~ is an unsigned integer type, and ~math.MaxUint64~ is its largest value.

## The store behind an interface
~LinkStore~ is an interface as in the protocols chapter: only ~Put~ and ~Get~. The handler does not know whether the store is an in-memory map or a database, and a test can inject a fake store that collides or fails on demand.
~~~go
type Shelf interface { Put(code, target string) error }
var s Shelf = &memShelf{items: map[string]string{}}
~~~
The store judges duplicates, because only it sees every code under its lock. ~ErrTaken~ and ~ErrNoLink~ are sentinel errors; they may come back wrapped with ~%w~, so test them with ~errors.Is~, not ~==~.
~~~go
if errors.Is(err, ErrTaken) { continue } // Try another number
~~~

## Collisions and attempts
Numbers come from an injected ~next~ function, which may be random in production. Collisions can happen, so try a limited number of times. Endless attempts could hang the request once the store is full. Any error other than a collision is a real failure; do not retry it.

Do not draw a number before the request is known to be valid; a bad request consumes nothing.

## The target link
~url.Parse~ from the CORS lesson decodes the link. Accept only ~http~ and ~https~: ~javascript:~ and ~data:~ run code in the victim's browser. ~Host~ must be nonempty and ~User~ must be ~nil~, because ~https://crew@evil.test~ misleads the reader.
~~~go
u, err := url.Parse("https://dock.orbital.test/a?b=1")
ok := err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil
~~~
Also reject any control byte or space (any byte ~<= ' '~ or 127), which can break headers or fool the eye.

## The redirect
302 (~http.StatusFound~) means the thing lives elsewhere, with the address in a ~Location~ header. ~http.Redirect~ adds a small HTML body for GET requests, so set the header and status yourself here to keep the body empty. 503 (~http.StatusServiceUnavailable~) means the server cannot serve right now.
~~~go
w.Header().Set("Location", target)
w.WriteHeader(http.StatusFound)
~~~

## Common mistakes
~~~go
code := strconv.Itoa(len(links)+1) // Codes reveal link counts and repeat after deletion
for { if store.Put(code, t) == nil { break } } // Endless loop that even retries failures
if err == ErrTaken { } // A wrapped error slips through
~~~

## Your mission
Implement ~Base62(n uint64) string~ with the alphabet ~0-9~, then ~a-z~, then ~A-Z~, no leading zeros, and ~Base62(0)~ equal to ~0~. Implement ~NewMemoryStore() LinkStore~, concurrency-safe and independent: ~Put~ on an existing code returns ~ErrTaken~ without overwriting, and ~Get~ on a missing code returns ~ErrNoLink~.

Implement ~NewLinkDesk(store LinkStore, next func() uint64) http.Handler~ (both inputs non-~nil~). ~POST /links~ takes a body of at most 2048 bytes holding ~{"target":"..."}~. Invalid JSON, an oversized body, or a target that is not ~http~/~https~, lacks a host, carries a user, or contains a byte ~<= ' '~ or 127 returns 400 and ~bad target\n~ without calling ~next~. Then make up to 3 attempts: code ~Base62(next())~ and ~store.Put~. On ~ErrTaken~ (even wrapped) try again; after 3 collisions return 503 and ~no free code\n~. Any other error returns 500 and ~store failure\n~ without retrying. Success returns 201, ~application/json~, ~{"code":"C","path":"/go/C"}~, and ~Location: /go/C~.

~GET /go/C~, where C is 1 to 11 alphabet characters, calls ~store.Get~: ~ErrNoLink~ (even wrapped) returns 404 and ~no link\n~, another error 500 and ~store failure\n~, and success 302 with ~Location~ set to the target and an empty body. Ignore the query. A malformed code or any other path returns 404 and ~no link\n~ without asking the store. A wrong method on ~/links~ or ~/go/C~ returns 405, ~Allow: POST~ or ~Allow: GET~, and ~wrong method\n~. Errors use ~text/plain; charset=utf-8~. No ~Report~, and ~main~ stays empty.

Third star: ~ParseBase62(s string) (uint64,error)~ is the exact inverse of ~Base62~: it rejects empty text, any character outside the alphabet, a leading zero on more than one digit, and any value above ~math.MaxUint64~.`,
			Hints: []string{"Fix the alphabet order before thinking about number representation.", "Let the storage interface own saving and lookup.", "Use ~errors.Is~ to recognize a wrapped collision error."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch275Hidden},
	})
}

const ch275Starter = `package main

import (
	"errors"
	"net/http"
)

type LinkStore interface {
	Put(code, target string) error
	Get(code string) (string, error)
}

var ErrTaken = errors.New("code taken")
var ErrNoLink = errors.New("no link")

func Base62(n uint64) string                                       { return "" }
func ParseBase62(s string) (uint64, error)                         { return 0, errors.New("TODO") }
func NewMemoryStore() LinkStore                                    { return nil }
func NewLinkDesk(store LinkStore, next func() uint64) http.Handler { return http.NotFoundHandler() }
func main()                                                        {}
`

const ch275Hidden = `package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

const linkDigits = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func linkRef(n uint64) string {
	s := ""
	for {
		s = string(linkDigits[n%62]) + s
		n /= 62
		if n == 0 {
			return s
		}
	}
}

type fakeStore struct {
	mu         sync.Mutex
	puts, gets int
	put        func(n int, code, target string) error
	get        func(code string) (string, error)
}

func (f *fakeStore) Put(code, target string) error {
	f.mu.Lock()
	f.puts++
	n := f.puts
	f.mu.Unlock()
	return f.put(n, code, target)
}
func (f *fakeStore) Get(code string) (string, error) {
	f.mu.Lock()
	f.gets++
	f.mu.Unlock()
	return f.get(code)
}

func linkCall(h http.Handler, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(method, path, strings.NewReader(body)))
	return w
}
func linkBody(target string) string {
	raw, _ := json.Marshal(map[string]string{"target": target})
	return string(raw)
}
func linkExpect(t *testing.T, w *httptest.ResponseRecorder, status int, body string) {
	t.Helper()
	if w.Code != status || w.Body.String() != body || w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatalf("want %d %q, got %d %q", status, body, w.Code, w.Body.String())
	}
}

// Boundaries: zero, single and double digits, the largest uint64, random values.
func TestHiddenLinksBase62(t *testing.T) {
	for n, want := range map[uint64]string{0: "0", 9: "9", 10: "a", 35: "z", 36: "A", 61: "Z", 62: "10", 3843: "ZZ", 3844: "100", math.MaxUint64: linkRef(math.MaxUint64)} {
		if got := Base62(n); got != want {
			t.Fatalf("Base62(%d) = %q, want %q", n, got, want)
		}
	}
	rng := rand.New(rand.NewSource(2751))
	for i := 0; i < 500; i++ {
		n := rng.Uint64() >> uint(rng.Intn(64))
		if got := Base62(n); got != linkRef(n) {
			t.Fatalf("Base62(%d) = %q", n, got)
		}
	}
}

// Boundaries: duplicate codes, missing codes, separate stores, concurrent claims on one code.
func TestHiddenLinksMemoryStore(t *testing.T) {
	s := NewMemoryStore()
	if _, err := s.Get("x"); !errors.Is(err, ErrNoLink) {
		t.Fatal("missing link", err)
	}
	if s.Put("x", "https://a.test") != nil {
		t.Fatal("put")
	}
	if err := s.Put("x", "https://b.test"); !errors.Is(err, ErrTaken) {
		t.Fatal("duplicate", err)
	}
	if got, err := s.Get("x"); got != "https://a.test" || err != nil {
		t.Fatal("overwritten", got, err)
	}
	if _, err := NewMemoryStore().Get("x"); !errors.Is(err, ErrNoLink) {
		t.Fatal("shared stores")
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	won := 0
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := s.Put("race", fmt.Sprintf("https://r%d.test", i))
			_, _ = s.Get("race")
			mu.Lock()
			defer mu.Unlock()
			if err == nil {
				won++
			} else if !errors.Is(err, ErrTaken) {
				t.Error("unexpected", err)
			}
		}(i)
	}
	wg.Wait()
	if won != 1 {
		t.Fatal("winners", won)
	}
}

// Boundaries: rejected targets never draw a number, the 2048-byte cap, exact redirects, path shapes, methods.
func TestHiddenLinksDeskCreateAndFollow(t *testing.T) {
	rng := rand.New(rand.NewSource(2752))
	draws := 0
	var last uint64
	h := NewLinkDesk(NewMemoryStore(), func() uint64 { draws++; last = rng.Uint64() >> uint(rng.Intn(40)); return last })
	long := "https://relay.orbital.test/" + strings.Repeat("a", 2048)
	for _, target := range []string{"", "relay.orbital.test", "ftp://relay.orbital.test", "javascript:alert(1)", "//relay.orbital.test/x", "https://", "https:///x", "https://crew:pw@relay.orbital.test", "https://relay.orbital.test/a b", "https://relay.orbital.test/\t", "http://relay.orbital.test/\x7f", long} {
		linkExpect(t, linkCall(h, "POST", "/links", linkBody(target)), 400, "bad target\n")
	}
	for _, body := range []string{"", "null", "[]", "{\"target\":5}", linkBody("https://a.test") + " {}"} {
		linkExpect(t, linkCall(h, "POST", "/links", body), 400, "bad target\n")
	}
	if draws != 0 {
		t.Fatal("rejected input drew codes", draws)
	}
	edge := linkBody("https://a.test/")
	edge = strings.TrimSuffix(edge, "\"}") + strings.Repeat("x", 2048-len(edge)) + "\"}"
	if w := linkCall(h, "POST", "/links", edge); w.Code != 201 || len(edge) != 2048 {
		t.Fatal("2048-byte body", w.Code, len(edge))
	}
	for i := 0; i < 20; i++ {
		target := fmt.Sprintf("https://dock-%d.orbital.test/manifest?crate=%d#lid", rng.Intn(999), rng.Intn(99999))
		if i%3 == 0 {
			target = "http://beacon.orbital.test:8080/" + fmt.Sprint(i)
		}
		w := linkCall(h, "POST", "/links", linkBody(target))
		var out struct{ Code, Path string }
		code := Base62(last)
		if w.Code != 201 || w.Header().Get("Content-Type") != "application/json" || json.Unmarshal(w.Body.Bytes(), &out) != nil || out.Code != code || out.Path != "/go/"+code || w.Header().Get("Location") != "/go/"+code {
			t.Fatalf("create %d %q %v", w.Code, w.Body.String(), w.Header())
		}
		w = linkCall(h, "GET", "/go/"+code+"?utm=1", "")
		if w.Code != 302 || w.Header().Get("Location") != target || w.Body.Len() != 0 {
			t.Fatalf("follow %d %v %q", w.Code, w.Header(), w.Body.String())
		}
		w = linkCall(h, "POST", "/go/"+code, "")
		linkExpect(t, w, 405, "wrong method\n")
		if w.Header().Get("Allow") != "GET" {
			t.Fatal("allow GET")
		}
	}
	w := linkCall(h, "GET", "/links", "")
	linkExpect(t, w, 405, "wrong method\n")
	if w.Header().Get("Allow") != "POST" {
		t.Fatal("allow POST")
	}
	for _, path := range []string{"/", "/go/", "/go/zzzzzzzzzzzz", "/go/a-b", "/go/abc/", "/links/", "/GO/abc"} {
		linkExpect(t, linkCall(h, "GET", path, ""), 404, "no link\n")
	}
}

// Boundaries: injected collisions (wrapped), exhausted retries, store failures, and lookups that never reach the store.
func TestHiddenLinksDeskCollisionsAndFailures(t *testing.T) {
	seq := uint64(1000)
	next := func() uint64 { seq++; return seq }
	taken := 0
	f := &fakeStore{put: func(n int, code, target string) error {
		if n <= 2 {
			taken++
			return fmt.Errorf("claim %s: %w", code, ErrTaken)
		}
		if code != Base62(seq) || target != "https://a.test" {
			t.Error("stored", code, target)
		}
		return nil
	}}
	w := linkCall(NewLinkDesk(f, next), "POST", "/links", linkBody("https://a.test"))
	if w.Code != 201 || f.puts != 3 || !strings.Contains(w.Body.String(), Base62(1003)) {
		t.Fatal("retry after collisions", w.Code, f.puts, w.Body.String())
	}
	f = &fakeStore{put: func(int, string, string) error { return ErrTaken }}
	linkExpect(t, linkCall(NewLinkDesk(f, next), "POST", "/links", linkBody("https://a.test")), 503, "no free code\n")
	if f.puts != 3 {
		t.Fatal("attempts", f.puts)
	}
	f = &fakeStore{put: func(int, string, string) error { return errors.New("disk on fire") }}
	linkExpect(t, linkCall(NewLinkDesk(f, next), "POST", "/links", linkBody("https://a.test")), 500, "store failure\n")
	if f.puts != 1 {
		t.Fatal("retried a real failure", f.puts)
	}
	f = &fakeStore{get: func(code string) (string, error) {
		switch code {
		case "gone":
			return "", fmt.Errorf("lookup: %w", ErrNoLink)
		case "boom":
			return "", errors.New("disk on fire")
		}
		return "https://ok.test", nil
	}}
	h := NewLinkDesk(f, next)
	linkExpect(t, linkCall(h, "GET", "/go/gone", ""), 404, "no link\n")
	linkExpect(t, linkCall(h, "GET", "/go/boom", ""), 500, "store failure\n")
	if f.gets != 2 {
		t.Fatal("gets", f.gets)
	}
	for _, path := range []string{"/go/", "/go/a_b", "/go/" + strings.Repeat("Z", 12), "/go/ok/x"} {
		linkExpect(t, linkCall(h, "GET", path, ""), 404, "no link\n")
	}
	if f.gets != 2 {
		t.Fatal("malformed codes reached the store", f.gets)
	}
}

func TestHiddenBonusLinksParseBase62(t *testing.T) {
	rng := rand.New(rand.NewSource(2753))
	for i := 0; i < 300; i++ {
		n := rng.Uint64() >> uint(rng.Intn(64))
		if got, err := ParseBase62(linkRef(n)); got != n || err != nil {
			t.Fatalf("ParseBase62(%q) = %d %v", linkRef(n), got, err)
		}
	}
	if got, err := ParseBase62(linkRef(math.MaxUint64)); got != math.MaxUint64 || err != nil {
		t.Fatal("max", got, err)
	}
	if got, err := ParseBase62("0"); got != 0 || err != nil {
		t.Fatal("zero")
	}
	for _, bad := range []string{"", "00", "0a", "a-b", " a", "é", "lYGhA16ahyg", "zzzzzzzzzzzz", "ZZZZZZZZZZZ"} {
		if _, err := ParseBase62(bad); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
}
`
