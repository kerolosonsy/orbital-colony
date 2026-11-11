# خطة توسيع المنهج — Orbital Colony × 100 Go Projects

> **الغرض:** مواصفات تنفيذ كاملة يستلمها منفّذ تاني (إنسان أو agent) من غير ما يحتاج يرجع للمحادثة.
> **المصدر:** مقارنة منهج اللعبة (17 فصل، 51 مهمة) بمنهج `~/Desktop/Projects/Golang/golang-100-projects/plan.md` (100 مشروع).
> **الحالة:** خطة فقط. مفيش كود اتكتب.
> **النتيجة المستهدفة:** كل مشروع من الـ 100 ليه مهمة مقابلة في اللعبة. من 51 مهمة إلى ~168 مهمة، ومن 17 فصل إلى 37.

## 0. قرارات صاحب المشروع (نهائية)

1. اللي مكانه الطبيعي في فصل قديم يتحط فيه كمهمة إضافية؛ اللي محتاج فصل جديد يتعمله فصل.
2. مسموح بـ `net`، `net/http`، `os` وأي مكتبة في الـ labs. اللعبة محلية ولمستخدم واحد.
3. **كل حاجة تتضاف بشكلها الحقيقي**: SQLite، Postgres، Redis، MongoDB، Kafka، gRPC، Gin، Fiber، GORM... مش fakes بس.
4. التنفيذ على يد حد تاني؛ الملف ده هو المرجع.

## 1. قواعد للمنفّذ

- ريبو `~/Desktop/Projects/Golang` **للقراءة فقط** (سياسة `AI_POLICY.md` بتاعته). ممنوع أي كتابة فيه. الـ README بتاع كل مشروع هناك هو مرجع المتطلبات والـ edge cases للمهمة المقابلة.
- المهمات بتعلّم **نفس المفهوم بقصة المستعمرة**، ومش نسخة من المشروع. ممنوع الحل المرجعي بتاع اللعبة يبقى حل جاهز لمشروع من الـ 100 (أسماء، دومين، وشكل API مختلفين).
- اتبع أسلوب الكود الموجود. متعملش refactor لحاجة مش مطلوبة.
- لغة الدروس: عربي مصري + نسخة `En` كاملة (الـ `i18n_test` بيفرض ده). الكود والمصطلحات إنجليزي.
- كل مرحلة تخلص بـ `go test ./...` أخضر في `cmd/orbital` و `colony`.

## 2. خريطة الكود الحالي (اللي هتتعامل معاه)

| ملف | الدور | ملاحظة مهمة |
|---|---|---|
| `colony/levels/levels.go` | `Level`، `Lab`، `LabFile`، `Chapters`، `add()` | ترتيب `All` = ترتيب أسماء الملفات (`ch01.go`…)؛ `TestLevelOrder` بيفرض ترتيب الفصول |
| `colony/levels/chNN.go` | مهمات كل فصل: `Setup`، `Check`، درس، starter، hints، `En` | الدرس markdown مصغّر: `~code~` و `~~~` (مفيش backticks) |
| `cmd/orbital/policy.go` | `playerImports` allowlist واحدة للكل | هتتقسم لمهمة / lab |
| `cmd/orbital/features.go` | `scan()` + `checkRequirements()`؛ `sel:pkg.Name` شغال عام | هنا تتضاف الـ Requires الجديدة |
| `cmd/orbital/runner.go` | `Run()` للمهمات، `goCmd()` | `GOPROXY=off`، `GOFLAGS=-mod=mod`، go.mod بيتولّد بـ `replace orbital/colony` |
| `cmd/orbital/lab.go` | `RunLab()`: تستات اللاعب → Buggy → Mutants → Hidden (`^TestHidden`) → Bench | أسماء الملفات مسطّحة و `.go` بس؛ timeout 90s؛ مفيش `-race` |
| `cmd/orbital/levels_test.go` | `TestLevelsSolvable`: حل مرجعي لكل مهمة ×5 seeds، والـ starter لازم يتبني وميعدّيش | الحلول في `testdata/solutions/<id>.go` أو `<id>/` للـ labs |
| `colony/levels/i18n_test.go` | `En` كامل، نفس tokens الـ starter، رسائل `Check` مترجمة | |
| `cmd/orbital/workspace.go` | go.mod بتاع VS Code | لازم ياخد نفس الـ deps |
| `web/app.js:964` | `unlocked(i)`: السابق أو الحالي عنده نجوم | هيتعدّل |

## 3. المرحلة P0 — البنية (لازم تخلص قبل أي مهمة جديدة)

### P0.1 توسيع `Level` / `Lab` / `Chapter`
- `Chapter.Track string` + `TrackEn`: مسار الفصل (القسم 9).
- `Level.Deps []string`: module paths خارجية مسموحة للمهمة (مثلاً `github.com/gin-gonic/gin`).
- `Level.Service string`: `""` | `postgres` | `mongo` | `kafka`.
- `Level.CarryFrom string`: ID lab سابق؛ ملفات اللاعب المحفوظة منه تبقى الـ starter (للـ capstone).
- `Lab.Fuzz string`: اسم `FuzzXxx` يتشغّل `-fuzz=^Name$ -fuzztime=5s`.
- `Lab.Bonus bool`: لو فيه `TestHiddenBonus*`: من غيرها نجمتين، بيها 3.
- `Level.Race` يشتغل على الـ labs كمان (`go test -race`).
- `LabFile.Name`: يقبل مسار فرعي محلي (`filepath.IsLocal`) وامتدادات `.go .sql .proto .json .yaml .txt .md .html`. غير الـ `.go` ميتفحصش imports.

### P0.2 `policy.go`
- `missionImports` = القائمة الحالية + `bufio`، `io/fs`، `path`، `html`، `text/template`، `log/slog`، `crypto/sha256`، `crypto/rand`، `time/tzdata`.
- `labDenied` = `unsafe`، `os/exec`، `syscall`، `plugin`. أي stdlib تاني مسموح في الـ labs.
- مكتبة خارجية مسموحة بس لو prefix بتاعها في `Level.Deps`.
- `restrictedImports(code, lvl, lang)`؛ عدّل المناديين في `runner.go` و `lab.go`. تست: mission بيرفض `net/http`، lab بيقبله، lab بيرفض `os/exec`، dep مش في `Deps` بيترفض.

### P0.3 المكتبات الخارجية (offline)
- module جديد `labdeps/` (`module orbital/labdeps`) فيه `deps.go` بـ blank imports لكل المكتبات (القسم 8) + `go.sum`.
- `make deps`: `cd labdeps && GOPROXY=https://proxy.golang.org go mod download all` (المرة الوحيدة اللي فيها نت).
- `goTest()` و `Run()`: لو `len(lvl.Deps) > 0` → ضيف `require` بنفس versions من `labdeps/go.mod` وانسخ `labdeps/go.sum`. `GOPROXY=off` يفضل زي ما هو.
- لو الـ build فشل بـ `missing go.sum entry` / `cannot find module` → رسالة: «شغّل `make deps` مرة واحدة».
- `workspace.go`: نفس الـ requires عشان gopls.
- ارفع timeout الـ lab لـ 240s لما `Deps` مش فاضية (أول build لـ grpc/gin بطيء).
- راجع نسخة Go في `go.mod` (حالياً 1.24)؛ ارفعها لو مكتبة محتاجة أحدث، وحدّث README.

### P0.4 الخدمات الحقيقية
- `services/compose.yaml`: `postgres:16`، `mongo:7`، `redpandadata/redpanda` (Kafka API)، كلهم على `127.0.0.1` بس.
- Makefile: `services-up`، `services-down`، `services-status`.
- السيرفر بيبعت للتستات: `ORBITAL_PG_DSN`، `ORBITAL_MONGO_URI`، `ORBITAL_KAFKA_BROKERS`.
- `RunLab`: لو `lvl.Service != ""` اعمل TCP dial الأول؛ لو مش شغالة → step فاشل برسالة «شغّل `make services-up`» من غير ما تبني.
- مهمات الخدمات **متقفلش التقدّم**: `unlocked()` يعدّي مهمة `service` لو مش مكتملة.
- `levels_test.go`: مهمات `Service != ""` تعمل `t.Skip` إلا لو `ORBITAL_SERVICES=1`.
- Redis **مش محتاج Docker**: `miniredis` سيرفر حقيقي in-process جوه التست.
- SQLite: `modernc.org/sqlite` (pure Go، من غير cgo).

### P0.5 `features.go` — Requires جديدة
| key | قاعدة الـ AST |
|---|---|
| `rune` | تحويل `[]rune(...)` أو `range` على string أو selector من `utf8.` |
| `iota` | `ast.Ident` اسمه `iota` جوه `const` |
| `const` | `GenDecl` بـ `token.CONST` |
| `defer` | `ast.DeferStmt` ≥ 1 |
| `recover` | call لـ `recover` جوه `FuncLit` في `defer` |
| `errwrap` | `fmt.Errorf` فيه `%w` **و** (`errors.As` أو `errors.Is`) |
| `errtype` | type عليه method `Error() string` |
| `funcvalue` | `FuncType` كـ field/param/map value، أو `FuncLit` متخزّن في متغير |
| `closure` | `FuncLit` بيستخدم متغير من الـ scope الخارجي (تقريب: FuncLit جوه func بيرجّع func) |
| `embed` | struct field من غير اسم |
| `variadic` | param بـ `Ellipsis` |
| `bufchan` | `make(chan T, n)` بـ arg تاني |
| `close` | call لـ `close` |
| `recursion` | دالة بتنادي نفسها |
| `fuzz` | `FuncDecl` اسمه `Fuzz*` بـ `*testing.F` |
| `subtest` | call لـ `t.Run` |
| `iface` | `TypeSpec` نوعه `InterfaceType` معرّف من اللاعب |

لكل key: نص عربي/إنجليزي في `requirementTexts` + تست موجب وسالب.

### P0.6 `lab.go`
- `-race` لو `lvl.Race`؛ `WARNING: DATA RACE` يفشّل الـ step بنفس tip الموجودة.
- خطوة fuzz بعد الـ Hidden (`Lab.Fuzz`).
- `Lab.Bonus` للنجوم.
- `Checks` الحالية (Buggy/Mutants/Hidden/Bench) تفضل زي ما هي.
- الـ Hidden يقدر يبقى أكتر من ملف: `Lab.HiddenFiles map[string]string` (بجانب `Hidden` القديم).

### P0.7 باكدج مساعد `orbital/colony/labkit` (جديد، مسموح استيراده)
- `Clock` interface (`Now`، `After`، `NewTicker`) + `FakeClock` بـ `Advance(d)`.
- `LeakCheck(t)`: يقارن `runtime.NumGoroutine` قبل/بعد مع retry.
- `TestCA()` / `IssueCert(...)`: شهادات TLS للتستات (للفصل 32).
- `FreeAddr()`: `127.0.0.1:0`.
- مولّدات بيانات بـ seed: logs، CSV، نصوص عربي/emoji، frames binary.

### P0.8 colony API لمهمات البيانات (الفصول 2–9)
دوال جديدة في `colony/` بتقرا من `sim.World.Data` (بيتولّد في `Setup` بـ `w.Rand`)، والتحقق بـ `colony.Report` زي المهمات الحالية:
`Callsigns() []string` · `Readings() []string` · `Diagnose(i int) error` + `*FaultError` · `Dock()/Undock()/Drill()` · `Orders() []string` · `Exam() []Question` + `Answer(i, choice)` · `Invoice() (totalCents, people int)`.

### P0.9 الواجهة
- `unlocked(i)`: كمان `true` لو **أي مهمة بعدها** عندها نجوم (عشان المهمات المحشورة في فصول قديمة متقفلش لاعب متقدّم)، ولو المهمة السابقة `service` ومش مكتملة.
- قائمة الفصول مجمّعة بالـ Track، قابلة للطي.
- الـ lab يعرض الملفات غير الـ `.go` (read-only) وشارة للخدمة المطلوبة + حالتها.
- شارة «جاهز لمشروع 0XX» بعد النجاح (حقل `Level.Bridge []int`).

**P0 تم لما:** كل التستات القديمة خضرا، تست لكل بند فوق، ومهمة lab تجريبية بتستورد `github.com/gin-gonic/gin` بتتبني offline بعد `make deps`.

## 4. تعريف «تمّت» لأي مهمة

1. `Level` كامل: ID، Chapter، Title، Goal، Lesson، Starter/Files، ≥3 Hints، Requires، Par/Gold (للمهمات).
2. `En` كامل ونفس tokens الـ starter.
3. بيانات عشوائية بـ seed (حل ثابت منسوخ ميعدّيش).
4. حل مرجعي في `testdata/solutions/`، والـ starter بيتبني وميعدّيش.
5. للـ labs: Hidden tests بتغطي نجاح/فشل/حدود/تزامن حسب عمود «التحقق» في `plan.md` للمشروع المقابل؛ Mutants ≥3 لو المهمة عن كتابة تستات.
6. الدرس: «قبل ما تكتب» (سؤالين تصميم) + المفهوم + مثال مختلف عن الحل + أخطاء شائعة.
7. `Bridge` بأرقام المشاريع المقابلة.

## 5. مهمات تتضاف في الفصول الموجودة (P1)

النوع: **M** = مهمة درون/بيانات (`Report`)، **L** = lab.

| ID | العنوان | نوع | اللاعب بيعمل إيه / التحقق | Requires | مشروع |
|---|---|---|---|---|---|
| `2-3-callsigns` | أسماء النداء | M | يعكس أسماء عربي/emoji ويبلّغ عدد الحروف الحقيقي؛ `len` بيدّي رقم غلط | `rune` | 006 |
| `2-4-units` | وحدات القياس | M | custom types + `const`/`iota` لأنواع الشحنات وتحويل الوحدات؛ يبلّغ الإجمالي | `iota`,`const` | 003، 008 |
| `4-3-checkpoint` | بوابة الفحص | M | قراءات حسّاسات كنصوص، بعضها بايظ: `strconv` + حدود التصنيف (قيم على الحد بالظبط) | `if`,`errcheck` | 001، 007 |
| `5-3-fault-log` | سجل الأعطال | M | يلف خطأ بـ `%w`، يطلّع `*FaultError` بـ `errors.As`، يجمّع بـ `errors.Join` | `errwrap` | 002، 014 |
| `5-4-undock` | فك الارتباط | M | `Drill()` ممكن تعمل panic؛ لازم `Undock()` يحصل دايماً والدرون يكمّل | `defer`,`recover` | — |
| `5-5-command-table` | لوحة الأوامر | M | `map[string]func()` ينفّذ أوامر جاية كنصوص، وأمر مش معروف = خطأ | `funcvalue`,`map` | 011 |
| `6-4-stock-order` | ترتيب المخزون | M | ترتيب تنازلي مع كسر التعادل بموقع ثابت | `sel:slices.SortFunc` | 019 |
| `6-5-heat-grid` | شبكة الحرارة | L | جمع/ضرب/transpose لـ 2D slices؛ Hidden: أبعاد غلط = error، والمدخل ميتعدّلش | `slice` | 026 |
| `7-4-composite-units` | وحدات مركّبة | M | embedding + تقسيم فاتورة بالـ cents بسياسة تقريب معلنة (الباقي يتوزّع) | `embed`,`method` | 012 |
| `7-5-crew-exam` | اختبار الطاقم | M | `[]Question` structs، يجاوب ويحسب النسبة؛ سؤال من غير اختيارات | `struct`,`range` | 005 |
| `8-4-cart-chain` | سلسلة العربات | L | singly linked list: insert/delete/find في head/middle/tail وقايمة فاضية | `pointer`,`struct` | 029 |
| `9-4-small-plugs` | مقابس صغيرة | L | منطق مفصول عن الـ I/O: دالة بتاخد `io.Reader`/`io.Writer`؛ interface صغيرة بتعرّفها بنفسك | `iface`,`sel:io.Reader` | 011، 022 |
| `13-4-chain-of-command` | سلسلة القيادة | M | context بيعدّي 3 مستويات دوال؛ deadline من الأب بيوقف الكل؛ ممنوع ctx في struct | `sel:context.Context`,`select` | 041 |
| `14-4-airlock-queue` | طابور البوابة | L | `Queue[T]` بـ zero value شغال، FIFO، وفاضي من غير panic | `generictype` | 027 |
| `14-5-archive-tree` | شجرة الأرشيف | L | generic BST بـ comparator: insert/search/in-order، وسياسة التكرار | `generictype`,`recursion` | 028 |
| `15-4-deep-probe` | المسبار العميق | M | binary search على أعماق مرتبة + flood fill recursive | `recursion` | 024 |
| `15-5-rescue-priority` | أولوية الإنقاذ | M | `container/heap`: يخدم الأعلى أولوية أولاً مع وصول طلبات جديدة | `sel:heap.Push` | 090 |
| `16-4-wild-inputs` | مدخلات مجنونة | L | يكتب `FuzzXxx` يكسّر parser فيه bug، وبعدين يصلّحه | `fuzz` | 085 |
| `16-5-fake-clock` | ساعة مزيفة | L | subtests + `t.Helper` + `labkit.FakeClock` بدل `time.Sleep`؛ Mutants بتمسك تستات بتنام | `subtest` | 015، 036 |

## 6. الفصول الجديدة (18 → 37)

كلها **labs** إلا لو مكتوب M. عمود Deps مختصر؛ الأسماء الكاملة في القسم 8.

### فصل 18 — الأرشيف (JSON، CSV، نصوص) — P2
| ID | العنوان | اللاعب بيعمل إيه / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `18-1-mission-log` | سجل المهام | CRUD في الذاكرة بـ IDs ثابتة وترتيب مستقر + table tests | `tabletest` | 016 |
| `18-2-black-box` | الصندوق الأسود | JSON round trip، struct tags، ملف ناقص/بايظ = error واضح | `sel:json.Unmarshal` | 017 |
| `18-3-cargo-manifest` | بيان الشحنة | CSV من `io.Reader`: header، صف بايظ، quotes، ملف فاضي، إجماليات | `sel:csv.NewReader` | 018 |
| `18-4-signal-words` | كلمات الإشارة | أكتر الكلمات تكراراً بترتيب حتمي؛ ترقيم، case، تعادل، نص كبير | `map`,`sel:bufio.Scanner` | 019 |
| `18-5-flight-recorder` | مسجّل الرحلة | تحليل log بالـ stream: مستويات، سطر بايظ، سطر طويل جداً | `sel:bufio.Scanner`,`struct` | 021 |
| `18-6-crew-directory` | دليل الطاقم | دفتر: رفض ID مكرر، بحث case-insensitive، حفظ عن طريق interface | `iface` | 022 |

### فصل 19 — الساعة الكونية (`time`) — P2
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `19-1-durations` | مدد الرحلات | M: `time.Duration`، parse/format بـ layout | `sel:time.Parse` | 013 |
| `19-2-two-suns` | شمسين | تحويل وقت ثابت بين مناطق؛ منطقة مش موجودة؛ تغيّر اليوم | `sel:time.LoadLocation` | 013 |
| `19-3-countdown` | العد التنازلي | عدّاد بـ ticker محقون وقابل للإلغاء؛ مفيش ticker leak ولا انتظار حقيقي | `sel:context.Context`,`select` | 015 |
| `19-4-many-timers` | مؤقتات كتير | مؤقتات متوازية؛ كل نتيجة مرة واحدة وبترتيب عرض حتمي | `go`,`waitgroup` · race | 031 |

### فصل 20 — المترجم (parsing) — P2
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `20-1-bulletin` | نشرة المستعمرة | markup محدود → HTML بـ state machine و `strings.Builder`؛ escaping لـ HTML خبيث؛ syntax مش مقفول | `sel:strings.Builder` | 023 |
| `20-2-validators` | المدققين | validators مستقلة (إيميل/تليفون/تاريخ) + تجميع الأخطاء؛ تاريخ مستحيل؛ Unicode | `errwrap` | 014 |
| `20-3-lexer` | مقطّع الأوامر | lexer للغة أوامر الدرون؛ token غلط + موقعه | `struct`,`switch` | 093 |
| `20-4-evaluator` | المنفّذ | recursive descent + AST + متغيرات؛ أولويات وأقواس؛ الناتج بيحرّك الدرون فعلاً | `recursion`,`iface` | 093 |

### فصل 21 — سطر الأوامر والنظام (CLI) — P3
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `21-1-sysinfo` | حالة المحطة | `flag.FlagSet` + `os`/`runtime`؛ flags افتراضية، env موجود/ناقص، format ثابت | `sel:flag.NewFlagSet` | 009 |
| `21-2-subcommands` | أوامر فرعية | `run(args, stdin, stdout) int` لـ add/list/done/rm؛ exit codes | `iface` | 016 |
| `21-3-guess-depth` | خمّن العمق | لعبة تخمين بعشوائية محقونة: فوز، استنفاد المحاولات، مدخل مش رقم | `funcvalue` | 004 |
| `21-4-duel` | مبارزة الدرونات | جولات بـ custom type + `iota`؛ كل تركيبات فوز/خسارة/تعادل | `iota` | 008 |

### فصل 22 — نظام الملفات — P3
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `22-1-tree` | شجرة المخزن | طباعة شجرة مرتبة بحد عمق على `fs.FS` (`fstest.MapFS`)؛ عمق 0، مجلد فاضي | `recursion`,`sel:fs.FS` | 024 |
| `22-2-organizer` | المنظّم | خطة بـ dry-run (مفيش تغيير) ثم تنفيذ على `t.TempDir()`؛ تصادم أسماء، امتداد مجهول | `sel:filepath.WalkDir` | 020 |
| `22-3-atomic-save` | حفظ ذرّي | temp file + `os.Rename`؛ مفيش ملف نص-مكتوب لو حصل فشل في النص | `defer`,`sel:os.Rename` | 017 |
| `22-4-duplicates` | الملفات المكررة | تجميع بالحجم ثم SHA-256 بالـ stream؛ ممنوع تحميل الملف كله في الذاكرة | `sel:sha256.New`,`sel:io.Copy` | 025 |

### فصل 23 — الخزنة (crypto) — P3
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `23-1-access-codes` | أكواد الدخول | مولّد بطول وفئات حروف إجبارية؛ اختيارات مستحيلة؛ ممنوع `math/rand` | `sel:rand.Read` | 010 |
| `23-2-sealed-cargo` | شحنة مختومة | AES-GCM بـ nonce عشوائي؛ مفتاح غلط، تلاعب، عدم تكرار nonce | `sel:cipher.NewGCM` | 030 |
| `23-3-content-address` | عنوان بالمحتوى | hash-object + zlib + object store على مجلد مؤقت؛ hash معروف | `sel:sha1.New` | 097 |
| `23-4-snapshot-commit` | لقطة و commit | index + tree + commit فوق 23-3 | `CarryFrom` | 097 |

### فصل 24 — خطوط الإمداد (أنماط التزامن) — P4، كله `Race: true`
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `24-1-ping-pong` | بينج بونج | M: درونين يتبادلوا N مرة على unbuffered channel ثم قفل نضيف؛ N=0 | `chan`,`close` | 032 |
| `24-2-producers` | منتجين ومستهلكين | buffered channel؛ مفيش فقد/تكرار؛ صفر منتجين؛ إلغاء بدري | `bufchan`,`close` | 037 |
| `24-3-work-crew` | طاقم الشغل | worker pool عام: نتائج + أخطاء، backpressure، 1 و N workers، إلغاء | `chan`,`go`,`sel:context.Context` | 034 |
| `24-4-pipeline` | خط التكرير | fan-out/fan-in مع إلغاء؛ كل channel تتقفل مرة واحدة | `chan`,`select` | 044 |
| `24-5-beacon-check` | فحص المنارات | فحص URLs على `httptest` بحد توازي فعلي و timeout | `sel:http.Client` | 033 |
| `24-6-chunk-download` | تحميل بالأجزاء | Range requests + `io.WriterAt`؛ سيرفر من غير range؛ جزء فاشل؛ تنظيف | `sel:context.Context` | 035 |
| `24-7-map-tiles` | بلاطات الخريطة | تصغير صور بـ workers محدودين مع حفظ الصيغة؛ ملف بايظ | `go` · `x/image` | 039 |

### فصل 25 — الموارد المشتركة (sync) — P4، كله `Race: true`
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `25-1-fuel-account` | حساب الوقود | invariant: الرصيد ميبقاش سالب تحت إيداع/سحب متوازي | `mutex` | 038 |
| `25-2-race-report` | تقرير السباق | كود فيه race: يثبته بتست، يشرحه في تعليق، يصلّحه | `mutex`,`subtest` | 040 |
| `25-3-token-bucket` | خزان التصاريح | token bucket بساعة محقونة: capacity، refill، burst؛ تستات من غير sleep | `mutex` | 036 |
| `25-4-event-bus` | ناقل الأحداث | subscribe/unsubscribe؛ سياسة المشترك البطيء؛ Close idempotent | `chan`,`mutex` | 042 |
| `25-5-ttl-cache` | كاش بصلاحية | `RWMutex` + انتهاء + cleanup goroutine بيقف من غير leak | `sel:sync.RWMutex` | 043 |
| `25-6-atomic-bench` | atomic ولا mutex | عدّادين + benchmark؛ صحة تحت التزامن | `benchmark`,`sel:atomic.Int64` | 045 |

### فصل 26 — بوابة القاعدة (HTTP) — P5
| ID | العنوان | اللاعب / Hidden (`httptest`) | Requires · Deps | مشروع |
|---|---|---|---|---|
| `26-1-first-server` | أول سيرفر | health + صفحة `html/template` + server timeouts؛ method ممنوع؛ content type | `sel:http.Server` | 046 |
| `26-2-rest-crud` | REST | CRUD في الذاكرة: كل endpoint، 400/404، ID ثابت، طلبات متوازية | `mutex` · race | 047 |
| `26-3-router` | الراوتر | router صغير + سلسلة middleware + request ID في context؛ ترتيب، params، 404/405 | `funcvalue` | 048 |
| `26-4-envelope` | الظرف | envelope موحّد + ربط typed errors بـ status؛ مفيش تسريب خطأ داخلي | `errtype` | 049 |
| `26-5-static` | ملفات ثابتة | `embed` + `FileServer`؛ من غير directory listing؛ ETag/Cache-Control؛ 404 مخصص | `sel:http.FileServer` | 052 |
| `26-6-upload` | رفع ملفات | multipart بالـ stream؛ `MaxBytesReader`؛ نوع مرفوض؛ path traversal؛ تنظيف الجزئي | `sel:http.MaxBytesReader` | 051 |

### فصل 27 — التصاريح والحماية — P5
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `27-1-signed-pass` | تصريح موقّع | bcrypt + JWT قصير العمر + auth middleware؛ توقيع/انتهاء/claims؛ secret من env | bcrypt، jwt | 050 |
| `27-2-sessions` | الجلسات | session على السيرفر + cookie آمن + CSRF + rotation + logout | `sel:http.Cookie` | 059 |
| `27-3-cors` | CORS | allowlist + security headers؛ preflight؛ credentials؛ `Vary: Origin` | `funcvalue` | 056 |
| `27-4-rate-limit` | حد الطلبات | لكل عميل (يبني على 25-3)؛ 429، عزل العملاء، سياسة الـ proxy، تنظيف الذاكرة | `mutex` · race | 057 |
| `27-5-short-links` | روابط مختصرة | Base62 + store interface؛ تصادم محقون؛ URL بايظ | `iface` | 053 |

### فصل 28 — الإطارات والتشغيل — P5
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `28-1-gin` | Gin | نفس عقد 26-2 بـ Gin (binding، middleware) | gin | 054 |
| `28-2-fiber` | Fiber | نفس العقد بـ Fiber (`app.Test`) + ملف مقارنة قصير يكتبه اللاعب | fiber | 055 |
| `28-3-openapi` | عقد الـ API | مستند OpenAPI (JSON) مطابق للـ API؛ contract test للـ endpoints والـ schemas | `sel:json.Unmarshal` | 058 |
| `28-4-graceful` | إغلاق مرتب | طلب جاري يخلص جوه deadline، رفض الجديد، Close مرة واحدة | `sel:http.Server`,`sel:context.Context` | 060 |

### فصل 29 — قواعد البيانات (SQL) — P6
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `29-1-sqlite-crud` | أول قاعدة | `database/sql` + SQLite مؤقتة: CRUD، constraint، not found، context | sqlite | 061 |
| `29-2-migrations` | الترحيلات | up/down + جدول versions؛ فشل في النص؛ منع نسخة مكررة | sqlite | 064 |
| `29-3-transfer` | تحويل ذرّي | transaction: نجاح، رصيد غير كافي، فشل في النص = rollback، request ID مكرر | sqlite | 066 |
| `29-4-pool` | حوض الاتصالات | `SetMaxOpenConns` + `DBStats`؛ الحد ميتعدّاش؛ `WaitCount`؛ rows دايماً بتتقفل | sqlite | 068 |
| `29-5-sqlx` | sqlx | repository بـ named queries و struct scan؛ unit بـ fake store، و integration على Postgres | sqlx، pgx · `Service: postgres` | 062 |
| `29-6-gorm` | GORM | علاقات صريحة، preload، وتوثيق عدد الـ queries (N+1) | gorm، gorm-sqlite | 063 |

### فصل 30 — التخزين المتقدم — P6
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `30-1-cache-aside` | كاش جانبي | Redis حقيقي (miniredis): hit/miss، TTL، invalidation عند التحديث، fallback لو Redis وقع | go-redis، miniredis | 065 |
| `30-2-documents` | المستندات | CRUD + unique index + cursor pagination + context timeout | mongo-driver · `Service: mongo` | 067 |
| `30-3-search-index` | فهرس البحث | inverted index + ranking بسيط + normalization + صفحات من غير تكرار | `map` | 069 |
| `30-4-read-write-split` | قراءة وكتابة | router لـ primary/replica؛ فشل replica؛ sticky read بعد الكتابة | `iface` | 070 |

### فصل 31 — الشبكة الخام — P7
| ID | العنوان | اللاعب / Hidden (loopback، port عشوائي) | Requires · Deps | مشروع |
|---|---|---|---|---|
| `31-1-echo` | الصدى | TCP line protocol بحدود و deadlines؛ رسالة أكبر من الحد؛ timeout؛ العميل قفل | `sel:net.Conn` | 071 |
| `31-2-chat` | غرفة الراديو | hub loop + broadcast؛ اسم مكرر؛ عميل بطيء | `chan`,`select` · race | 072 |
| `31-3-discovery` | اكتشاف الجيران | UDP على loopback: encode/decode، timeout، packet بايظ | `sel:net.ListenPacket` | 073 |
| `31-4-frames` | إطارات binary | parser offline بـ bounds checks + fuzz؛ إطار مقطوع/مجهول | `fuzz`,`sel:binary.BigEndian` | 085 |
| `31-5-dns` | دليل العناوين | lookup عن طريق resolver محقون؛ A/AAAA/CNAME، NXDOMAIN، ترتيب حتمي | `iface` | 080 |
| `31-6-port-check` | فحص المنافذ | **localhost فقط**: مفتوح/مقفول، timeout، حد توازي، رفض أي host تاني | `go`,`sel:net.DialTimeout` | 082 |

### فصل 32 — البروتوكولات — P7
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `32-1-websocket` | بث حي | WebSocket: read/write pumps، ping/pong، حد الرسالة، disconnect من غير leak | coder/websocket | 074 |
| `32-2-grpc-unary` | gRPC | خدمة unary + interceptor + status codes على `bufconn`؛ `.proto` و `pb.go` جاهزين read-only | grpc | 075 |
| `32-3-grpc-stream` | gRPC streaming | bidirectional: رسائل متعددة، إلغاء العميل، خطأ سيرفر، EOF نضيف | grpc | 076 |
| `32-4-tls` | TLS و mTLS | شهادات `labkit`: موثوقة/منتهية/اسم غلط؛ mTLS من غير client cert؛ حد أدنى للـ TLS | `sel:tls.Config` | 084 |
| `32-5-http2` | HTTP/2 | TLS + تفاوض `h2` + طلبات متوازية على اتصال واحد (من غير Server Push) | `sel:http.Server` | 077 |
| `32-6-ssh` | SSH | executor على سيرفر SSH in-process: host key verification، allowlist أوامر، timeout | x/crypto/ssh | 081 |

### فصل 33 — الوسطاء — P7
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `33-1-retry-client` | عميل عنيد | `RoundTripper` بـ retry محدود و backoff/jitter محقون؛ مفيش retry لـ POST من غير key؛ body بيتقفل | `iface` | 083 |
| `33-2-reverse-proxy` | الوكيل العكسي | `httputil.ReverseProxy` بـ allowlist و timeout وسياسة headers؛ backend واقع | `sel:httputil.ReverseProxy` | 078 |
| `33-3-load-balancer` | موزّع الأحمال | round-robin + health؛ backend مريض؛ الكل واقع؛ طلبات متوازية | `mutex` · race | 079 |
| `33-4-gateway` | البوابة | routing + auth + rate limit + request ID (يجمع 27 و 33) | `funcvalue` | 091 |

### فصل 34 — محركات التخزين — P8
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `34-1-kv-store` | مخزن مفاتيح | KV متزامن بـ TTL و snapshot round trip | `sel:sync.RWMutex` · race | 087 |
| `34-2-wal` | سجل الكتابة | WAL + memtable + استرجاع بعد restart + checksums | `sel:crc32.ChecksumIEEE` | 088 |
| `34-3-compaction` | الدمج | SSTables + tombstones + compaction بيحافظ على النتائج؛ ملف بايظ | `CarryFrom` | 088 |
| `34-4-chunk-store` | تخزين موزّع | chunks + replication + فقد node + chunk بايظ + repair + dedup | `sel:sha256.Sum256` | 099 |

### فصل 35 — الأنظمة الموزعة — P8
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `35-1-task-queue` | طابور المهام | retry + lease + DLQ؛ worker وقع؛ تكرار؛ idempotency | go-redis، miniredis | 086 |
| `35-2-scheduler` | المجدول | priority queue + ساعة مزيفة؛ منع التداخل؛ misfire policy | `sel:heap.Push` | 090 |
| `35-3-redis-lock` | القفل الموزّع | token + TTL + compare-and-delete؛ غير المالك ميفكّش؛ انتهاء؛ حدود الأمان في تعليق | go-redis، miniredis | 094 |
| `35-4-outbox` | صندوق الصادر | حدث + كتابة في transaction واحدة؛ إعادة نشر؛ مستهلك idempotent؛ schema version | sqlite | 095 |
| `35-5-kafka` | الوسيط | producer/consumer: keys، offsets، retry، graceful close؛ unit بـ fake + integration حقيقي | kafka-go · `Service: kafka` | 098 |
| `35-6-raft-election` | الانتخابات | محاكاة Raft حتمية في الذاكرة: terms، انتخاب، فشل القائد، term قديم | `struct`,`labkit.FakeClock` | 089 |
| `35-7-raft-log` | نسخ السجل | log replication + commit بالأغلبية + state machine | `CarryFrom` | 089 |

### فصل 36 — التشغيل والمراقبة — P8
| ID | العنوان | اللاعب / Hidden | Requires · Deps | مشروع |
|---|---|---|---|---|
| `36-1-metrics` | العدادات | counters/histograms بـ labels محدودة + endpoint + تحديث متوازي | prometheus | 096 |
| `36-2-containers` | مدير الحاويات | CLI لـ list/start/stop عن طريق interface + fake؛ تأكيد قبل stop؛ مفيش حذف | `iface` | 092 |
| `36-3-ops` | جاهزية | `log/slog` منظّم مع redaction + `/livez` + `/readyz` مرتبط بالقاعدة | `sel:slog.New` | 100 |

### فصل 37 — مركز القيادة (Capstone) — P9، `Service: postgres`
API لأطقم/مشاريع/مهام المستعمرة. كل مرحلة lab بـ `CarryFrom` اللي قبلها و Hidden contract tests.

| ID | المرحلة | التحقق | مشروع |
|---|---|---|---|
| `37-1-schema` | domain + migrations up/down + repository | fresh up/down، constraints | 100 |
| `37-2-auth` | تسجيل/دخول، session token معتم (digest بس متخزّن)، logout، إلغاء عند تغيير الباسورد | cookie flags، expiry، revocation | 100 |
| `37-3-tenants` | أطقم، أدوار owner/admin/member، عزل tenants (الوصول عبر الأطقم = not found) | isolation، منع تصعيد، audit | 100 |
| `37-4-tasks` | مهام بـ `Idempotency-Key`، optimistic version، cursor pagination ثابت | تكرار الطلب، conflict، صفحات من غير تكرار | 100 |
| `37-5-operate` | request ID، error shape ثابت، metrics، rate limit، graceful shutdown | drain window، redaction، `-race` | 100، 091، 096 |

## 7. تغطية الـ 100 مشروع (للمراجعة)

001→4-3 · 002→5-3 · 003→2-4 · 004→21-3 · 005→7-5 · 006→2-3 · 007→4-3 · 008→2-4، 21-4 · 009→21-1 · 010→23-1 · 011→5-5، 9-4 · 012→7-4 · 013→19-1، 19-2 · 014→20-2 · 015→19-3 · 016→18-1، 21-2 · 017→18-2، 22-3 · 018→18-3 · 019→18-4 · 020→22-2 · 021→18-5 · 022→18-6 · 023→20-1 · 024→22-1 · 025→22-4 · 026→6-5 · 027→14-4 · 028→14-5 · 029→8-4 · 030→23-2 · 031→19-4 · 032→24-1 · 033→24-5 · 034→24-3 · 035→24-6 · 036→25-3 · 037→24-2 · 038→25-1 · 039→24-7 · 040→25-2 · 041→13-4 · 042→25-4 · 043→25-5 · 044→24-4 · 045→25-6 · 046→26-1 · 047→26-2 · 048→26-3 · 049→26-4 · 050→27-1 · 051→26-6 · 052→26-5 · 053→27-5 · 054→28-1 · 055→28-2 · 056→27-3 · 057→27-4 · 058→28-3 · 059→27-2 · 060→28-4 · 061→29-1 · 062→29-5 · 063→29-6 · 064→29-2 · 065→30-1 · 066→29-3 · 067→30-2 · 068→29-4 · 069→30-3 · 070→30-4 · 071→31-1 · 072→31-2 · 073→31-3 · 074→32-1 · 075→32-2 · 076→32-3 · 077→32-5 · 078→33-2 · 079→33-3 · 080→31-5 · 081→32-6 · 082→31-6 · 083→33-1 · 084→32-4 · 085→31-4، 16-4 · 086→35-1 · 087→34-1 · 088→34-2، 34-3 · 089→35-6، 35-7 · 090→35-2، 15-5 · 091→33-4 · 092→36-2 · 093→20-3، 20-4 · 094→35-3 · 095→35-4 · 096→36-1 · 097→23-3، 23-4 · 098→35-5 · 099→34-4 · 100→37-1…37-5.

## 8. المكتبات الخارجية (`labdeps`)

النسخ تتحدد بـ `go get <pkg>@latest` وقت التنفيذ وتتثبّت في `labdeps/go.mod`؛ متتكتبش من الذاكرة.

| الاختصار | module | مهمات |
|---|---|---|
| sqlite | `modernc.org/sqlite` | 29-1…29-4، 35-4 |
| sqlx | `github.com/jmoiron/sqlx` | 29-5 |
| pgx | `github.com/jackc/pgx/v5` (stdlib adapter) | 29-5، 37 |
| gorm | `gorm.io/gorm` + `github.com/glebarez/sqlite` | 29-6 |
| go-redis | `github.com/redis/go-redis/v9` | 30-1، 35-1، 35-3 |
| miniredis | `github.com/alicebob/miniredis/v2` | نفس اللي فوق |
| mongo-driver | `go.mongodb.org/mongo-driver/v2` | 30-2 |
| kafka-go | `github.com/segmentio/kafka-go` | 35-5 |
| gin | `github.com/gin-gonic/gin` | 28-1 |
| fiber | `github.com/gofiber/fiber/v2` | 28-2 |
| jwt | `github.com/golang-jwt/jwt/v5` | 27-1 |
| bcrypt / ssh | `golang.org/x/crypto` | 27-1، 32-6، 37-2 |
| x/image | `golang.org/x/image` | 24-7 |
| grpc | `google.golang.org/grpc` + `google.golang.org/protobuf` | 32-2، 32-3 |
| websocket | `github.com/coder/websocket` | 32-1 |
| prometheus | `github.com/prometheus/client_golang` | 36-1، 37-5 |

الـ `pb.go` بتاع gRPC بيتولّد مرة وقت التأليف (`make proto`، محتاج `protoc` عند المؤلف بس) ويتحفظ في الريبو كملف lab read-only.

## 9. المسارات (Tracks) في الواجهة

| Track | الفصول |
|---|---|
| الأساسيات | 1–9 |
| التزامن | 10–13، 24، 25 |
| الحِرفة (generics، خوارزميات، تستات) | 14–17 |
| البيانات والأدوات | 18–23 |
| الخدمات (HTTP) | 26–28 |
| قواعد البيانات | 29–30 |
| الشبكات | 31–33 |
| الأنظمة | 34–36 |
| Capstone | 37 |

رقم الفصل هو ترتيب الفتح؛ الـ Track للعرض بس.

## 10. المراحل

| مرحلة | المحتوى | عدد المهمات | يعتمد على | تم لما |
|---|---|---|---|---|
| P0 | البنية (القسم 3) | 0 | — | تستات البنية + lab تجريبي بـ dep خارجي |
| P1 | مهمات الفصول القديمة (القسم 5) | 19 | P0.5، P0.8، P0.9 | `TestLevelsSolvable` + `i18n_test` |
| P2 | فصول 18–20 | 14 | P0 | نفس الشرط |
| P3 | فصول 21–23 | 12 | P0.1، P0.2 | نفس الشرط |
| P4 | فصول 24–25 | 13 | P0.6، P0.7 | نفس الشرط تحت `-race` |
| P5 | فصول 26–28 | 15 | P0.3 | نفس الشرط بعد `make deps` |
| P6 | فصول 29–30 | 10 | P0.3، P0.4 | + `ORBITAL_SERVICES=1` للـ gated |
| P7 | فصول 31–33 | 16 | P0.3، P0.7 | نفس الشرط |
| P8 | فصول 34–36 | 14 | P6 | نفس الشرط |
| P9 | Capstone 37 + الإضافات (القسم 11) + README/PRODUCT | 5 | الكل | مراجعة الواجهة على 37 فصل |

P2، P3، P4 مستقلين عن بعض بعد P0 وينفع يتنفّذوا بالتوازي (كل فصل ملف `chNN.go` لوحده). التعارض الوحيد المتوقع: `Chapters` في `levels.go`، و`requirementTexts` في `features.go`.

## 11. إضافات فوق المنهجين (P9)

- **Boss لكل Track** مبني على الـ Level Gates في `plan.md` (مثلاً: التزامن = خدمة واحدة تعدّي `-race` + `LeakCheck`).
- **أسئلة تقييم ذاتي** آخر كل فصل (5 أسئلة فهم، من غير كود).
- **مراجعة متباعدة**: كل 5 مهمات، مهمة قديمة بـ seed جديد ومن غير starter.
- **فصل اختياري «Go الحديثة»**: `iter.Seq` و range-over-func، `b.ReportAllocs` و escape analysis، `struct{}` كـ set، labeled break، `sync.Once`، `errgroup`، `go vet` lab، pprof.

## 12. مخاطر معروفة

- **أول build للمكتبات التقيلة** (grpc، gorm، mongo) ممكن يعدّي الـ timeout: الـ build cache بيحلّها بعد أول مرة؛ `make deps` يعمل warm-up build كمان.
- **الـ Hidden tests قابلة للقراءة** من كود اللاعب بعد السماح بـ `os`: مقبول (قرار 2).
- **مهمات الخدمات** متتختبرش في `go test` العادي: لازم جولة `ORBITAL_SERVICES=1` قبل إعلان P6/P8/P9 خلصوا.
- **حجم المحتوى**: ~117 مهمة × (درس عربي + إنجليزي + حل + Hidden). الجودة أهم من السرعة: مهمة واحدة مكتملة حسب القسم 4 أحسن من عشرة ناقصين.
- **gRPC `pb.go`** لازم يتولّد بنفس نسخة `protobuf` اللي في `labdeps`.
