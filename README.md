# Orbital Colony 🚀

لعبة بتعلّمك **Go** من الصفر لحد الاحتراف. إنت قائد مستعمرة على كوكب "كيبلر-ج"، وبتكتب كود Go حقيقي يتحكم في الدرونات: تزرع، وتحصد، وتعدّن، وتنجو من العواصف.

الكود اللي بتكتبه بيتبني ويشتغل بالـ `go` الحقيقي على جهازك، فرسائل الأخطاء والسلوك هما نفس اللي هتقابلهم في أي شغل حقيقي.

## التشغيل

```bash
go run ./cmd/orbital
```

وبعدين افتح http://127.0.0.1:7777

أو بالـ Makefile (اكتب `make` لوحدها تشوف كل الأوامر):

```bash
make start     # شغّل في الخلفية وافتح المتصفح
make restart   # بعد أي تحديث للكود
make stop
make backup    # نسخة احتياطية من كل اللاعبين
make players   # مين متسجل
```

محتاج Go 1.24 أو أحدث للمهمات العادية، وGo 1.27.1 أو أحدث للمختبرات اللي بتستخدم مكتبات خارجية (حسب `labdeps/go.mod`). السيرفر بيشتغل على localhost بس، لأنه بيشغّل الكود اللي بتكتبه.

المهمات المُقيّمة بتقبل مكتبات الدروس والحساب ومعالجة البيانات مع `orbital/colony` و`orbital/colony/labkit`. المختبرات بتقبل مكتبات Go القياسية، ومنها `net/http` و`os`، ما عدا `unsafe` و`os/exec` و`syscall` و`plugin`. المكتبات الخارجية مسموحة بس لو مسارها تحت module موجود في `Level.Deps`. ملفات اللعبة الداخلية مش مسموحة. التفاصيل في [policy.go](cmd/orbital/policy.go). تشغيل الملفات من الـ terminal بيستخدم Go مباشرة. ده مش بديل لعزل كود غير موثوق، والسيرفر مخصص للتشغيل المحلي.

شغّل `make deps` مرة واحدة مع الإنترنت لتنزيل المكتبات المثبّتة وعمل warm-up build؛ تشغيل كود اللاعب بعد كده بيستخدم `GOPROXY=off`. مختبرات الخدمات محتاجة `make services-up` (Docker Compose: Postgres وMongoDB وKafka على localhost). `make services-status` يعرض حالتها و`make services-down` يقفلها. اختبارات مهمات الخدمات بتتخطّى تلقائياً إلا لو `ORBITAL_SERVICES=1`.

## الشغل من VS Code (autocomplete وdocs)

دوس **🧩 افتح في VS Code** في أي مهمة. اللعبة بتعمل فولدر `workspace/<المهمة>/` وتفتحه في VS Code، وبيكون فيه `gopls` شغال عادي:

- autocomplete (Ctrl+Space)، وأسماء الدوال والـ parameters بتظهر وإنت بتكتب.
- تقف بالماوس على أي حاجة (hover) فيظهرلك الشرح بتاعها، وتعمل ⌘+Click على `colony.Move` فيوديك للكود الحقيقي.
- الأخطاء بتظهر وإنت بتكتب، قبل ما تشغّل.

كل ما تعمل Save، اللعبة في المتصفح بتشوف التغيير وبتشغّل الكود لوحدها (تقدر تقفلها من "شغّل لما أحفظ"). وكمان تقدر تشغّل من الـ terminal بتاع VS Code:

```bash
cd workspace/3-1-long-row && go run .
```

ساعتها هيطبعلك نتيجة المهمة في الـ terminal من غير رسوم.

## الفصول

| # | الفصل | المفهوم | الحالة |
|---|---|---|---|
| 1 | الهبوط | استدعاء الدوال وترتيب الأوامر | ✅ |
| 2 | التليمتري | المتغيرات والأنواع والحساب | ✅ |
| 3 | الحلقات | for بكل أشكالها | ✅ |
| 4 | القرارات | if و switch | ✅ |
| 5 | الدوال والأخطاء | func و error | ✅ |
| 6 | المخازن | slices و maps | ✅ |
| 7 | بناء الأدوات | structs و methods | ✅ |
| 8 | العناوين | pointers و nil | ✅ |
| 9 | البروتوكولات | interfaces و type switch | ✅ |
| 10 | السرب | goroutines و WaitGroup (أكتر من درون في نفس الوقت) | ✅ |
| 11 | قنوات الاتصال | channels و worker pool و select | ✅ |
| 12 | السباق | mutex، والـ race detector الحقيقي (`-race`)، وdeadlock | ✅ |
| 13 | العواصف | context: الإلغاء، والـ timeouts، وWithCancel | ✅ |
| 14 | القوالب | generics، وconstraints، وStack[T] | ✅ |
| 15 | الكهوف | BFS، والترتيب، وDijkstra، وbinary search، وrecursion، وheap | ✅ |
| 16 | المختبر | go test، وtable tests، وmutation testing، وTDD، وbenchmarks، وfuzzing، وساعة محقونة | ✅ |
| 18 | الأرشيف | JSON، وCSV، ومعالجة النصوص، وstreaming، وحفظ عن طريق interfaces | ✅ |
| 19 | الساعة الكونية / The Cosmic Clock | المدد، والمناطق الزمنية، والعد التنازلي، ومؤقتات متوازية قابلة للإلغاء | ✅ |
| 20 | المترجم / The Interpreter | نشرة HTML آمنة، والتحقق، وlexer، وAST بتحرّك الدرون | ✅ |
| 21 | سطر الأوامر والنظام / The Command Line | flags، والبيئة، وأوامر فرعية، وتخمين ومبارزة بعشوائية محقونة | ✅ |
| 22 | نظام الملفات / The File System | شجرة مرتبة، وتنظيم بخطة، وحفظ ذرّي، وكشف التكرار بالـ stream | ✅ |
| 23 | الخزنة / The Vault | أكواد آمنة، وAES-GCM، وعناوين محتوى مضغوطة، وفهرس ولقطات للشحنات / Secure codes, AES-GCM, compressed content addresses, and supply snapshots | ✅ |
| 24 | خطوط الإمداد / Supply Lines | ملكية القنوات، وطوابير محدودة، وطاقم عام، وfan-out/fan-in، وHTTP ونطاقات البايت وصور متوازية / Channel ownership, bounded queues, generic crews, fan-out/fan-in, HTTP ranges, and concurrent images | ✅ |
| 25 | الموارد المشتركة / Shared Resources | شروط الرصيد، وتقرير السباق، وخزان التصاريح، وناقل الأحداث، وكاش بصلاحية، ومقارنة atomic وmutex / Invariants, race reports, token buckets, event buses, TTL caches, and atomic versus mutex benchmarks | ✅ |
| 26 | بوابة القاعدة / Base Gateway | سيرفر HTTP، وREST، وراوتر وmiddleware، وظرف أخطاء، وملفات مضمّنة ورفع آمن / HTTP servers, REST, routing and middleware, error envelopes, embedded assets, and safe uploads | ✅ |
| 27 | التصاريح والحماية / Passes and Protection | تصاريح bcrypt وJWT، وجلسات بكوكيز آمنة وCSRF، وCORS، وحد طلبات لكل عميل، وروابط مختصرة Base62 / bcrypt and JWT passes, secure-cookie sessions with CSRF, CORS, per-client rate limits, and Base62 short links | ✅ |

## هيكل المشروع

- `colony/`: الـ module اللي كود اللاعب بيستورده (`orbital/colony`).
  - `internal/sim`: محاكاة العالم. اللاعب ميقدرش يستوردها، عشان ميغشّش.
  - `levels/`: كل مهمة: الخريطة والهدف والدرس والتلميحات.
- `cmd/orbital/`: السيرفر وكود التشغيل، وفيه كمان فحص المفاهيم المطلوبة بالـ AST.
  - `testdata/solutions/`: حل مرجعي لكل مهمة. `go test ./cmd/orbital` بيتأكد إن كل المهمات تتحل.
- `web/`: الواجهة (المحرر، والخريطة، والكونسول).

مهمات الفصل 16 اسمها "مختبر": مفيش درون. بتكتب تستات، واللعبة بتشغّل `go test` الحقيقي عليها: تستاتك لازم تمسك نسخ متبوّظة من الكود (mutants)، وتكشف الـ bug قبل ما تصلّحه، والـ benchmark بيقيس كودك أسرع كام مرة.

في الفصل 10 وما بعده الوقت حقيقي: كل حركة بتاخد 4ms، فالـ goroutines بتشتغل بالتوازي فعلاً، والوقت الكلي بيتقاس من أبطأ درون.
