# Orbital Colony 🚀

لعبة بتعلّمك **Go** من أول استدعاء دالة لحد بناء خدمات وأنظمة ومشروع API ختامي. إنت قائد مستعمرة على كوكب "كيبلر-ج"، وبتكتب كود Go حقيقي يتحكم في الدرونات، وبعدها تبني أدوات وخدمات للمستعمرة في المختبرات.

الكود اللي بتكتبه بيتبني ويشتغل بالـ `go` الحقيقي على جهازك، فرسائل الأخطاء والسلوك هما نفس اللي هتقابلهم في أي شغل حقيقي.

## اللي هتلاقيه

- **181 مهمة في 38 فصل**: 174 مهمة في الفصول 1–37، و7 مختبرات في فصل 38 الاختياري. الإجمالي 58 مهمة درون و123 مختبر، شاملين 8 مهمات ختام مسار (Boss).
- الدروس والتلميحات ورسائل اللعبة بالعربي المصري والإنجليزي، والكود Go حقيقي.
- درونات وبيانات بتتغيّر بالـ seed، وفحص للمفاهيم المطلوبة بالـ AST، ونجوم حسب شروط كل مهمة.
- مختبرات بـ `go test` الحقيقي: hidden tests، وmutants وbenchmarks وfuzzing حسب المهمة، و`-race` لما المهمة تطلبه.
- 5 أسئلة فهم في كل فصل، بأفضل نتيجة محفوظة؛ ومراجعة بتتعرض كل 5 مهمات جديدة مكتملة ببيانات جديدة. مراجعة الدرونات من غير starter، والمختبرات بتحتفظ بملفات بدايتها؛ النجوم والحلول المحفوظة مابتتغيّرش.
- ملفات لاعبين وحفظ على السيرفر، ونسخة محلية في المتصفح، وربط VS Code.
- **debugger جوه الكود**: أول Run بيتسجّل والخريطة بتمشي مع كل خطوة؛ Step into/over/out، ورجوع لورا، وbreakpoints (F9)، والمتغيرات والـ call stack تحت.
- **شرح لرسائل الأخطاء** بالعربي: معناها، وسببها غالباً، وتراجع إيه. السطر الغلط بيتعلّم بالأحمر والشرح بيظهر فوقه، وفيه زرار يصلّح الـ import الناقص أو الزيادة بضغطة.
- درس قصير لكل مهمة فيه walkthrough وملعب تفاعلي، وإحصائياتك محفوظة على جهازك (`make stats`).
- editor فيه autocomplete لـ `colony.` بشرح كل دالة، وأقواس بتقفل لوحدها، وTab/Shift+Tab للسطور المختارة، وCtrl/⌘+/ للتعليق. دوس `?` تشوف كل الاختصارات.

### اختصارات

| الاختصار | بيعمل إيه |
| --- | --- |
| Ctrl/⌘ + Enter | شغّل الكود |
| Ctrl/⌘ + Shift + Enter | شغّل خطوة بخطوة |
| Esc | وقّف التشغيل |
| Ctrl/⌘ + / | علّق السطور أو شيل التعليق |
| Shift + Alt + F | نسّق الكود بـ gofmt |
| `colony.` · Ctrl + Space | اقتراحات دوال colony مع أول سطر من شرحها |
| Alt + ↑/↓ · Shift + Alt + ↑/↓ | انقل السطور · كرّرها |
| Ctrl + G | روح لسطر |
| F9 | نقطة توقف على سطر المؤشر |
| F10 · F11 · Shift+F11 | عدّي · ادخل · اخرج من الدالة |
| F5 · Shift+F5 | كمّل · اقفل التتبّع |
| Alt+← · Alt+→ · Alt+Space | خطوة لورا · خطوة · شغّل/وقّف (وقت التتبّع) |
| ← / → | الخطوة الجاية/السابقة في الدرس |

تقدر تفتح مهمة مباشرة بالرابط: `http://127.0.0.1:7777/#3-1`.

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

### مكتبات المختبرات والخدمات

شغّل `make deps` مرة واحدة مع الإنترنت: بينزّل كل المكتبات بإصداراتها المثبّتة في `labdeps/go.mod` و`labdeps/go.sum`، وبعدها بيعمل warm-up build. السيرفر بيولّد module منفصل لكود اللاعب؛ لو المهمة ليها `Level.Deps` بياخد قائمة الإصدارات المثبّتة كلها من `labdeps`، بس سياسة الاستيراد بتسمح بالمكتبات الخارجية المذكورة في `Deps` للمهمة دي بس. البناء والاختبارات من اللعبة بيستخدموا `GOPROXY=off` و`GOTOOLCHAIN=local`، فجهّز المكتبات والـ toolchain الأول.

المهمة اللي ليها `Level.Service` بتفحص اتصال TCP بالخدمة قبل التشغيل، ولو مش متاحة بتطلب منك تشغّل `make services-up`. ده بيشغّل Docker Compose من `services/compose.yaml`: Postgres 16 وMongoDB 7 وKafka-compatible Redpanda على localhost. `make services-status` يعرض الحالة و`make services-down` يقفلها. SQLite شغّال داخل العملية، ومختبرات Redis بتستخدم `miniredis`، فمش محتاجين خدمة Docker منفصلة.

السيرفر بيمرّر القيم دي لاختبارات المختبرات؛ لو مش متضبطة بيستخدم القيم الافتراضية:

| المتغير | القيمة الافتراضية / الاستخدام |
|---|---|
| `ORBITAL_PG_DSN` | `postgres://orbital:orbital@127.0.0.1:5432/orbital?sslmode=disable` |
| `ORBITAL_MONGO_URI` | `mongodb://127.0.0.1:27017` |
| `ORBITAL_KAFKA_BROKERS` | `127.0.0.1:9092`؛ قائمة مفصولة بفواصل، وفحص التوافر بيجرّب أول broker بس |
| `ORBITAL_SERVICES` | اضبطه على `1` عشان اختبارات الحلول المرجعية تشغّل المهمات المعلَّمة بـ `Service`؛ مش مفتاح لتشغيل الخدمات ولا شرط للعب من المتصفح |

بعد تجهيز المكتبات والخدمات، جولة اختبارات السيرفر بالخدمات:

```bash
make services-up
ORBITAL_SERVICES=1 go test ./cmd/orbital/ -count=1 -timeout 30m
```

الـ Capstone بيختبر SQLite دايماً، وPostgres كمان لو الاتصال متاح؛ مفيهوش `Service` إجباري. `ORBITAL_DESK_SECRET` جزء من تمرين التصاريح `27-1-signed-pass` (سر Base64 يفك لـ 32 بايت على الأقل)، مش إعداد مطلوب لتشغيل اللعبة.

لو بتشغّل برنامج اللاعب من الـ terminal، `ORBITAL_LANG=en` بيختار رسائل إنجليزي (الافتراضي عربي)، و`ORBITAL_SEED` عدد صحيح لتكرار نفس العالم، و`ORBITAL_LEVEL` بيحدد المهمة؛ اللعبة بتضبطهم تلقائياً عند التشغيل منها. `ORBITAL_TRACE` و`ORBITAL_PREVIEW` و`ORBITAL_COLONY` و`ORBITAL_EVENTS_FD` و`ORBITAL_CMD_FD` متغيرات داخلية لتسجيل الأحداث والمعاينة وحفظ المستعمرة وقنوات الاتصال؛ مش محتاج تضبطها للتشغيل العادي.

## المساعد (شات جوه اللعبة)

زرار **اسأل** فوق بيفتح شات بيشوف المهمة والكود بتاعك وآخر رسايل الكونسول، وبيدّيك تلميحات خطوة بخطوة بدل الحل على طول. بيرد من خلال اشتراك ChatGPT بتاعك عن طريق [Codex CLI](https://github.com/openai/codex)، فمفيش مفتاح API في اللعبة:

```bash
brew install codex
codex login
```

اختار تسجيل الدخول بحساب ChatGPT. السيرفر بيشغّل `codex exec` في فولدر مؤقت فاضي وبصلاحية قراءة بس، فالمساعد بيتكلم بس ومبيعدّلش حاجة. كل رد بياخد حوالي 10 ثواني، وبيتحسب من حدود اشتراكك.

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

الأعداد من `levels.All` في `colony/levels`، شاملة مهمات الـ Boss؛ الأسئلة والمراجعات ووضع المستعمرة الحر مش مهمات إضافية.

| # | الفصل | المفهوم | المهمات | الحالة |
|---|---|---|---|---|
| 1 | الهبوط | استدعاء الدوال وترتيب الأوامر | 3 | ✅ |
| 2 | التليمتري | المتغيرات والأنواع والحساب، وbytes وrunes، وconst وiota | 4 | ✅ |
| 3 | الحلقات | for بكل أشكالها | 3 | ✅ |
| 4 | القرارات | if و switch | 3 | ✅ |
| 5 | الدوال والأخطاء | func وerror، ولف الأخطاء، وdefer وrecover، وclosures | 5 | ✅ |
| 6 | المخازن | slices وmaps، وslices.SortFunc، وشبكات 2D | 5 | ✅ |
| 7 | بناء الأدوات | structs وmethods وembedding | 5 | ✅ |
| 8 | العناوين | pointers وnil وسلاسل مترابطة | 4 | ✅ |
| 9 | البروتوكولات | interfaces وtype switch، وio.Reader وio.Writer | 5 | ✅ |
| 10 | السرب | goroutines و WaitGroup (أكتر من درون في نفس الوقت) | 3 | ✅ |
| 11 | قنوات الاتصال | channels و worker pool و select | 3 | ✅ |
| 12 | السباق | mutex، والـ race detector الحقيقي (`-race`)، وdeadlock | 3 | ✅ |
| 13 | العواصف | context: الإلغاء، والـ timeouts، وWithCancel | 4 | ✅ |
| 14 | القوالب | generics، وconstraints، وStack[T]، وطابور وشجرة عامة | 5 | ✅ |
| 15 | الكهوف | BFS، والترتيب، وDijkstra، وbinary search، وrecursion، وheap | 5 | ✅ |
| 16 | المختبر | go test، وtable tests، وmutation testing، وTDD، وbenchmarks، وfuzzing، وساعة محقونة | 5 | ✅ |
| 17 | الأعماق | طبقات، ومصفوفات 3D، وBFS ثلاثي الأبعاد | 4 | ✅ |
| 18 | الأرشيف | JSON، وCSV، ومعالجة النصوص، وstreaming، وحفظ عن طريق interfaces | 6 | ✅ |
| 19 | الساعة الكونية / The Cosmic Clock | المدد، والمناطق الزمنية، والعد التنازلي، ومؤقتات متوازية قابلة للإلغاء | 4 | ✅ |
| 20 | المترجم / The Interpreter | نشرة HTML آمنة، والتحقق، وlexer، وAST بتحرّك الدرون | 4 | ✅ |
| 21 | سطر الأوامر والنظام / The Command Line | flags، والبيئة، وأوامر فرعية، وتخمين ومبارزة بعشوائية محقونة | 4 | ✅ |
| 22 | نظام الملفات / The File System | شجرة مرتبة، وتنظيم بخطة، وحفظ ذرّي، وكشف التكرار بالـ stream | 4 | ✅ |
| 23 | الخزنة / The Vault | أكواد آمنة، وAES-GCM، وعناوين محتوى مضغوطة، وفهرس ولقطات للشحنات / Secure codes, AES-GCM, compressed content addresses, and supply snapshots | 5 | ✅ |
| 24 | خطوط الإمداد / Supply Lines | ملكية القنوات، وطوابير محدودة، وطاقم عام، وfan-out/fan-in، وHTTP ونطاقات البايت وصور متوازية / Channel ownership, bounded queues, generic crews, fan-out/fan-in, HTTP ranges, and concurrent images | 7 | ✅ |
| 25 | الموارد المشتركة / Shared Resources | شروط الرصيد، وتقرير السباق، وخزان التصاريح، وناقل الأحداث، وكاش بصلاحية، ومقارنة atomic وmutex / Invariants, race reports, token buckets, event buses, TTL caches, and atomic versus mutex benchmarks | 7 | ✅ |
| 26 | بوابة القاعدة / Base Gateway | سيرفر HTTP، وREST، وراوتر وmiddleware، وظرف أخطاء، وملفات مضمّنة ورفع آمن / HTTP servers, REST, routing and middleware, error envelopes, embedded assets, and safe uploads | 6 | ✅ |
| 27 | التصاريح والحماية / Passes and Protection | تصاريح bcrypt وJWT، وجلسات بكوكيز آمنة وCSRF، وCORS، وحد طلبات لكل عميل، وروابط مختصرة Base62 / bcrypt and JWT passes, secure-cookie sessions with CSRF, CORS, per-client rate limits, and Base62 short links | 5 | ✅ |
| 28 | الإطارات والتشغيل / Frameworks and Operations | نفس عقد REST على Gin وFiber، ومستند OpenAPI واختبار عقد بالمكتبة القياسية، وإغلاق مرتب بيصفّي الطلبات الجارية / The same REST contract on Gin and Fiber, an OpenAPI document with a standard-library contract test, and graceful shutdown that drains in-flight requests | 5 | ✅ |
| 29 | قواعد البيانات / SQL Databases | SQL وSQLite، وترحيلات وتحويل ذرّي وحوض اتصالات وsqlx وGORM / SQL and SQLite, migrations, atomic transfers, connection pools, sqlx, and GORM | 6 | ✅ |
| 30 | التخزين المتقدم / Advanced Storage | Redis وكاش جانبي، ومستندات MongoDB، وفهرس بحث، وفصل القراءة والكتابة / Redis cache-aside, MongoDB documents, search indexing, and read/write splitting | 5 | ✅ |
| 31 | الشبكة الخام / Raw Networking | TCP وحدود السطور، وغرفة راديو، وUDP، وإطارات binary، وDNS، وفحص منافذ محلي / TCP and bounded lines, a radio hub, UDP, binary frames, DNS, and loopback port checks | 6 | ✅ |
| 32 | البروتوكولات / Protocols | WebSocket وgRPC والبث وTLS وmTLS وHTTP/2 وSSH / WebSocket, unary and streaming gRPC, TLS, mTLS, HTTP/2, and SSH | 6 | ✅ |
| 33 | الوسطاء / Intermediaries | عميل بيعيد بأمان، ووكيل عكسي، وموزّع أحمال بفحص صحة، وبوابة API / Safe retrying clients, reverse proxies, health-checked load balancing, and an API gateway | 5 | ✅ |
| 34 | محركات التخزين / Storage Engines | مخزن مفاتيح بصلاحية ولقطات وسجل كتابة مسبقة بـ checksums ودمج SSTables وتخزين chunks موزّع بالتكرار والإصلاح / A TTL key store with snapshots, a checksummed write-ahead log, SSTable compaction, and replicated chunk storage with repair | 4 | ✅ |
| 35 | الأنظمة الموزعة / Distributed Systems | طوابير مهام بعقود إيجار وDLQ، ومجدول، وأقفال موزعة، وصندوق الصادر، وKafka، وانتخاب Raft ونسخ السجل / Task queues with leases and a DLQ, a scheduler, distributed locks, the outbox, Kafka, Raft elections, and log replication | 7 | ✅ |
| 36 | التشغيل والمراقبة / Operations and Observability | مقاييس Prometheus وlabels محدودة، ومدير حاويات بـ interface، وسجلات slog منظّمة، وفحوص الجاهزية / Prometheus metrics and bounded labels, an interface-driven container manager, structured slog logs, and readiness probes | 4 | ✅ |
| 37 | مركز القيادة / Command Center (Capstone) | API واحد لأطقم ومهام المستعمرة بيكبر على مراحل: schema وترحيلات وrepository على SQLite وPostgres، وجلسات معتمة بـ bcrypt، وأطقم معزولة بأدوار وaudit، ومهام بمفاتيح تكرار ونسخ متفائلة وصفحات cursor، وتشغيل بـ request ID وmetrics وحد طلبات وإغلاق مرتب / One colony crews-and-tasks API grown in stages: schema, migrations, a repository on SQLite and Postgres, opaque bcrypt-backed sessions, isolated crews with roles and audit, idempotent tasks with optimistic versions and cursor pages, and operations with request IDs, metrics, rate limiting, and graceful drain | 5 | ✅ |
| 38 | Go الحديثة / Modern Go — إضافي / Extras | اختياري بعد فصل 17: iterators وsets وlabels وOnce وerrgroup وReportAllocs وescape analysis وvet وpprof / Optional after chapter 17: iterators, sets, labels, Once, errgroup, ReportAllocs, escape analysis, vet, and pprof | 7 | ✅ |

## المسارات

المسارات بتجمع الفصول في الفهرس، وترتيب الفتح الأساسي حسب المهمات. المهمة اللي بعد مهمة بخدمة بتفضل متاحة حتى لو معملتش مهمة الخدمة. فصل 38 اختياري، بيفتح بعد آخر مهمة في فصل 17، ومش شرط لإكمال المسار الأساسي.

| المسار | الفصول | المهمات | الختام |
|---|---|---|---|
| الأساسيات | 1–9 | 37 | Boss في 9 |
| التزامن | 10–13، 24–25 | 27 | Boss في 25 |
| الحِرفة (generics، خوارزميات، تستات) | 14–17 | 19 | Boss في 17 |
| البيانات والأدوات | 18–23 | 27 | Boss في 23 |
| الخدمات | 26–28 | 16 | Boss في 28 |
| قواعد البيانات | 29–30 | 11 | Boss في 30 |
| الشبكات | 31–33 | 17 | Boss في 33 |
| الأنظمة | 34–36 | 15 | Boss في 36 |
| Capstone | 37 | 5 | API واحد بيتطوّر على مراحل |
| إضافي | 38 | 7 | Go الحديثة، اختياري |

مهمات الـ Boss مختبرات بتجمع مفاهيم المسار، من غير خدمات خارجية مطلوبة. فصل 37 بيكمّل على كودك من المرحلة السابقة (`CarryFrom`): schema وترحيلات، وبعدها جلسات، وأطقم معزولة، ومهام، وتشغيل ومراقبة.

## هيكل المشروع

- `colony/`: الـ module اللي كود اللاعب بيستورده (`orbital/colony`).
  - `internal/sim`: محاكاة العالم. اللاعب ميقدرش يستوردها، عشان ميغشّش.
  - `levels/`: كل مهمة: الخريطة والهدف والدرس والتلميحات.
- `cmd/orbital/`: السيرفر وكود التشغيل، وفيه كمان فحص المفاهيم المطلوبة بالـ AST.
  - `testdata/solutions/`: الحلول المرجعية. `go test ./cmd/orbital` بيختبر الحلول والـ starters؛ المهمات المعلَّمة بخدمة محتاجة `ORBITAL_SERVICES=1` وخدمات شغّالة.
- `web/`: الواجهة (المحرر، والخريطة، والكونسول)، JavaScript عادي من غير build step. اختباراتها بـ `make test-web` (Node 20+).
- `labdeps/`: module منفصل لتثبيت مكتبات المختبرات وتجهيزها، مش module اللاعب.
- `services/`: إعداد Docker Compose للخدمات المحلية.

المختبرات مش مقتصرة على الفصل 16؛ موجودة كمان وسط المنهج وفي فصول التوسعة والـ Bosses: بتعرض خطوات التحقق بدل حركة الدرون، وبتشغّل `go test` الحقيقي على ملفات المهمة في module مؤقت. حسب المهمة، بتكتب تستات تمسك نسخ متبوّظة (mutants)، أو تصلّح bug، أو تنفّذ API عشان تعدّي hidden tests. مهمات الأداء بتستخدم benchmarks، وبعض المهمات فيها fuzzing أو bonus tests أو فحص `go vet` منفصل. `Race: true` بيضيف `-race`، وأي data race بتفشّل المهمة.

في الفصل 10 وما بعده الوقت حقيقي: كل حركة بتاخد 4ms، فالـ goroutines بتشتغل بالتوازي فعلاً، والوقت الكلي بيتقاس من أبطأ درون.

### الاختبارات

```bash
make check      # fmt + vet + اختبارات سريعة + اختبارات الواجهة
make test       # كل حاجة، شاملة كل المهمات (حوالي 8 دقايق)
make test-web   # اختبارات المتصفح بس
```

`TestLevelsSolvable` بيشغّل المهمات بالتوازي (نص عدد الـ CPUs)؛ غيّر العدد بـ `ORBITAL_TEST_PARALLEL=4`. فصول 10–13 (الوقت الحقيقي) بتفضل واحدة ورا التانية.

## Known gaps — فجوات معروفة

- نجاح الاختبارات العادية مش دليل على تكامل الخدمات: `TestLevelsSolvable` بيتخطّى مهمات Postgres وMongoDB وKafka المعلَّمة بـ `Service` إلا مع `ORBITAL_SERVICES=1` وخدمات حقيقية متاحة. فحص TCP في الواجهة مش فحص جاهزية أو صحة بيانات؛ واختبارات Postgres في الـ Capstone ممكن تعمل skip لو الاتصال مش متاح حتى في جولة الخدمات.
- `34-3-compaction` بيتحقق من إعادة الفتح، وفساد الجداول، وإعادة تشغيل WAL بعد flush؛ لكنه مش بيعالج crash في نص compaction، بين تثبيت الجدول الجديد وحذف كل الجداول القديمة.
- المختبرات بتسمح بـ `os`، فالـ hidden tests مش سرّية أو معزولة عن كود اللاعب. سياسة الاستيراد مش sandbox لتشغيل كود غير موثوق؛ اللعبة للتشغيل المحلي.
