package levels

import (
	"fmt"
	"slices"
	"time"

	"orbital/colony/i18n"
	"orbital/colony/internal/sim"
)

func init() {
	add(&Level{ID: "19-1-durations", Chapter: 19, LessonRevision: 1, Steps: ch19DurationsArabic, Title: "مدد الرحلات", Goal: "احسب وصول الرحلات من جدول المغادرة والمدة", Requires: []string{"sel:time.Parse"}, Bridge: []int{13},
		Lesson: `# مواعيد الرحلات

## قبل ما تكتب
- هترفض المدة السالبة حتى لو parser قبلها إزاي؟
- إزاي تحافظ على الترتيب لما أمر في النص يفشل؟

## مدة الرحلة غير التاريخ
جدول الرحلات لحظة ومعاها مدة. ~time.Duration~ عدد نانوثواني، مش تاريخ. استورد ~time~. الوحدات وSub من16-5؛ هنا هنقرا المدة ونطبع شكلها الموحّد.
~~~go
d,err:=time.ParseDuration("90m")
fmt.Println(d.String(),err) // 1h30m0s <nil>
fmt.Println(d>2*time.Hour) // false
~~~
~time.ParseDuration~ بتقبل h وm وs وتركيباتهم وكسورهم. وحدة غلط أو overflow يرجّع خطأ. ~Duration.String~ تطبع الشكل الموحّد. الصفر صالح؛ السالب ممكن يتقرا صح بس الرحلة لازم ترفضه صراحة.

## الـ layout تاريخ مرجعي
~time.Parse~ بترجع ~time.Time~ وخطأ. أول parameter ترتيب التاريخ المرجعي 2006-01-02 15:04:05، مش حروف YYYY أو MM. من غير منطقة في الـ layout القراية UTC.
~~~go
t,err:=time.Parse("2006-01-02 15:04:05","2031-04-03 23:40:00")
if err!=nil {return}
next:=t.Add(35*time.Minute)
fmt.Println(next.Format("2006-01-02 15:04:05")) // 2031-04-04 00:15:00
~~~
~Time.Add~ بترجع لحظة جديدة و~Time.Format~ بتعرضها بنفس التاريخ المرجعي. عبور نص الليل والسنة بيحصل تلقائي. المدة وقت منقضي؛ مش افتراض إن كل يوم محلي في كل منطقة24 ساعة.

## اقرأ عقد الجدول
~colony.Orders~ نسخة من المدخلات بالترتيب زي5-5. ~strings.Split~ جديدة هنا: بتقسم النص على كل ظهور للـseparator وبترجع []string، امتداد لـSplitN من18-5. الحقلين مضمونين؛ مش مطلوب parser عام للسجلات.
~~~go
parts:=strings.Split("2031-04-03 23:40:00|35m","|")
fmt.Println(parts[1]) // 35m
colony.Report("invalid")
~~~
تقرير لكل مدخل حتى لو بايظ. متستخدمش ساعة الجهاز ومتستناش الرحلة تخلص.

## أخطاء شائعة
~~~go
time.Parse("YYYY-MM-DD",text) // مش layout بتاع Go
_ = start.Add(duration) // ضيّعت القيمة الجديدة
~~~

## المهمة
اقرأ ~colony.Orders() []string~. كل أمر حقلين بالظبط بينهم ~|~: مغادرة UTC تقبلها ~time.Parse("2006-01-02 15:04:05",...)~ ومدة تقبلها ParseDuration. ابعت تقرير string واحد لكل أمر بنفس الترتيب. لو أي parser فشل أو المدة برا [0,48*time.Hour]، ابعت ~invalid~. غير كده ابعت ~arrival|duration~: arrival هي departure.Add(duration).Format("2006-01-02 15:04:05") والمدة Duration.String(). الثواني والكسور حسب سلوك الدوال دي. مفيش تقارير زيادة ولا حركة ولا time.Now ولا انتظار حقيقي.
`, Hints: []string{"الـduration وقت منقضي، والـinstant بيمثل لحظة على التقويم.", "اقرأ الحقول كل واحد لوحده وحافظ على مكان كل نتيجة.", "راجع نتيجة ~Time.Add~ وformat الوقت بعد ما تفحص أخطاء parsing."},
		En: &LevelText{LessonRevision: 1, Steps: ch19DurationsEnglish, Title: "Flight Durations", Goal: "Calculate arrivals from departure and elapsed duration", Lesson: `# Flight Schedules

## Before you write
- How will you reject negative durations even when parsing succeeds?
- How will you preserve order when an intermediate order fails?

## Elapsed duration and calendar time
A flight schedule is an instant plus elapsed time. ~time.Duration~ counts nanoseconds; it is not a date. Import ~time~. The units and Sub were introduced in 16-5; now parse a duration and display its canonical form.
~~~go
d,err:=time.ParseDuration("90m")
fmt.Println(d.String(),err) // 1h30m0s <nil>
fmt.Println(d>2*time.Hour) // false
~~~
~time.ParseDuration~ accepts compound units such as h, m, s and fractions; bad units and overflow return errors. ~Duration.String~ prints a canonical duration. Zero is valid; negative durations parse successfully, so reject them explicitly for flights.

## Layouts describe a reference date
~time.Parse~ returns a ~time.Time~ and error. Its first argument is the reference date 2006-01-02 15:04:05 arranged like the input, rather than YYYY or MM placeholders. A layout without a zone is parsed in UTC.
~~~go
t,err:=time.Parse("2006-01-02 15:04:05","2031-04-03 23:40:00")
if err!=nil {return}
next:=t.Add(35*time.Minute)
fmt.Println(next.Format("2006-01-02 15:04:05")) // 2031-04-04 00:15:00
~~~
~Time.Add~ returns a new instant; ~Time.Format~ uses the same reference layout to render it. It carries midnight and year boundaries automatically. A duration adds elapsed time, without assuming a calendar day is always 24 hours in every location.

## Read the schedule contract
~colony.Orders~ returns copied input strings in order, as in 5-5. ~strings.Split~ is new here: it separates at every occurrence of a delimiter and returns []string, extending SplitN from 18-5. The input guarantees two fields; this task does not require a general record parser.
~~~go
parts:=strings.Split("2031-04-03 23:40:00|35m","|")
fmt.Println(parts[1]) // 35m
colony.Report("invalid")
~~~
Use one Report for every input, even on failure. Do not use the host clock or wait for flights to finish.

## Common mistakes
~~~go
time.Parse("YYYY-MM-DD",text) // Not a Go layout
_ = start.Add(duration) // Discards the new instant
~~~

## Your mission
Read ~colony.Orders() []string~. Each order has exactly two fields separated by ~|~: a UTC departure accepted by ~time.Parse("2006-01-02 15:04:05",...)~ and a duration accepted by ParseDuration. For each order, in input order, send exactly one string Report. If either parser fails or the duration lies outside [0,48*time.Hour], report ~invalid~. Otherwise report ~arrival|duration~, where arrival is departure.Add(duration).Format("2006-01-02 15:04:05") and duration is Duration.String(). Preserve seconds and fractions as those methods specify. No extra reports or movement; no time.Now or real waiting.
`, Hints: []string{"A duration is elapsed time; an instant is a point on the calendar.", "Read fields separately and preserve each result's position.", "Check ~Time.Add~'s new value and format time after parsing errors."}, Starter: flightDurationsStarter}, Starter: flightDurationsStarterAr, Par: 20, Gold: 10, MaxTicks: 100, Setup: setupFlightDurations})
	add(&Level{ID: "19-2-two-suns", Chapter: 19, LessonRevision: 1, Steps: ch19SunsArabic, Title: "شمسين", Goal: "اعرض نفس اللحظة في مرصدين من غير الاعتماد على الجهاز", Requires: []string{"sel:time.LoadLocation"}, Bridge: []int{13},
		Lesson: `# شمسين ونفس اللحظة

## قبل ما تكتب
- مين اللي بيتغير: اللحظة ولا عرضها؟
- لو تاني منطقة غلط هل ترجع نص النتيجة؟

## لحظة واحدة وساعتين حائط
مرصدين المستعمرة بيقروا نفس اللحظة في منطقتين. Parse وFormat من19-1. ~time.RFC3339~ layout جاهزة فيها فرق UTC؛ ~time.LoadLocation~ بتجيب قواعد اسم IANA وبترجع خطأ للاسم المجهول.
~~~go
t,err:=time.Parse(time.RFC3339,"2030-12-31T23:30:00Z")
loc,err:=time.LoadLocation("Asia/Tokyo")
local:=t.In(loc)
fmt.Println(local.Format("2006-01-02"),err) // 2031-01-01 <nil>
~~~
~Time.In~ تغيّر العرض من غير تحريك اللحظة. ~Time.Equal~ تقارن اللحظات؛ == بتقارن بيانات المنطقة الداخلية كمان.
~~~go
fmt.Println(t.Equal(local)) // true
~~~
افحص كل خطأ قبل الناتج. متبدلش اسم مجهول بـUTC في السر.

## قواعد محلية وعرض صريح
اعمل blank import لقاعدة البيانات المرفقة عشان تشتغل من غير ملفات مناطق على الجهاز أو نت. الـ blank identifier بيشغّل تهيئة الباكدج من غير اسم تستخدمه في الكود.
~~~go
import _ "time/tzdata"
~~~
اعرض offset رقمي بدل اختصار حسب اللغة. ~-0700~ في layout فرق ساعات ودقايق بإشارة، و~15:04~ ساعة بنظام24 ودقايق.
~~~go
fmt.Println(local.Format("2006-01-02 15:04 -0700"))
fmt.Println(t.Format("2006-01-02")!=local.Format("2006-01-02"))
~~~
قارن التاريخ كله بالسنة والشهر. قواعد المنطقة ممكن تختلف حسب الموسم؛ النجمة التالتة بتختبر شتا وصيف بعيد عن لحظات تغيير التوقيت الملتبسة.

## أخطاء شائعة
~~~go
local:=t.Add(9*time.Hour) // غيّر اللحظة نفسها
loc,_:=time.LoadLocation(name) // تجاهل الاسم المجهول
~~~

## المهمة
نفّذ ~ReadSuns(instant,home,remote string) (SunReadings,error)~؛ الحقول ~Home,Remote string~ و~DayChanged bool~. اقرأ instant بـRFC3339، وحل الاسمين بـLoadLocation. اعرض كل لحظة محلية بـ ~2006-01-02 15:04 -0700~، وDayChanged true لو تاريخهم الكامل مختلف. أي خطأ يرجّع SunReadings{} وخطأ فيه instant أو اسم المنطقة اللي فشلت. استخدم time/tzdata؛ مفيش Now ولا Local. الاسم الفاضي وLocal مش ضمن المدخلات. النجمة التالتة لقواعد الموسم.
`, Hints: []string{"الـtime zone بتغيّر العرض المحلي، مش اللحظة نفسها.", "حل كل المدخلات قبل ما ترجع نتيجة كاملة.", "راجع ~LoadLocation~ ثم ~In~ وoffset الرقمي في format."},
		En: &LevelText{LessonRevision: 1, Steps: ch19SunsEnglish, Title: "Two Suns", Goal: "Render one instant at two observatories independently of the host", Lesson: `# Two Suns, One Instant

## Before you write
- What changes: the instant or its representation?
- Should an invalid second location return half a result?

## One instant, two wall clocks
The colony's two observatories read one instant in different locations. Reuse Parse and Format from 19-1. ~time.RFC3339~ is a built-in layout including a UTC offset; ~time.LoadLocation~ resolves an IANA name and returns an error for unknown names.
~~~go
t,err:=time.Parse(time.RFC3339,"2030-12-31T23:30:00Z")
loc,err:=time.LoadLocation("Asia/Tokyo")
local:=t.In(loc)
fmt.Println(local.Format("2006-01-02"),err) // 2031-01-01 <nil>
~~~
~Time.In~ changes the representation without moving the instant. ~Time.Equal~ compares instants; == also compares internal location data.
~~~go
fmt.Println(t.Equal(local)) // true
~~~
Check each error before using its returned value. Never silently replace an unknown location with UTC.

## Offline rules and explicit output
Blank-import the bundled database so this lab works without host zone files or network access. The blank identifier runs a package's initialization without naming it in your code.
~~~go
import _ "time/tzdata"
~~~
Use numeric offsets, not localized abbreviations. In a layout, ~-0700~ means signed hours and minutes; ~15:04~ means 24-hour hours and minutes.
~~~go
fmt.Println(local.Format("2006-01-02 15:04 -0700"))
fmt.Println(t.Format("2006-01-02")!=local.Format("2006-01-02"))
~~~
Compare the entire date, including year and month. A location can have seasonal offset rules; the third star checks winter and summer instants away from transition ambiguity.

## Common mistakes
~~~go
local:=t.Add(9*time.Hour) // Changes the instant
loc,_:=time.LoadLocation(name) // Discards unknown-zone errors
~~~

## Your mission
Implement ~ReadSuns(instant,home,remote string) (SunReadings,error)~ with ~Home,Remote string~ and ~DayChanged bool~. Parse instant using RFC3339 and resolve both names using LoadLocation. Render each local instant with ~2006-01-02 15:04 -0700~; DayChanged reports whether their full dates differ. Every error returns SunReadings{} and an error mentioning instant or the failing zone name. Import time/tzdata; do not use Now or Local. Empty names and Local are outside the input domain. The third star checks seasonal rules.
`, Hints: []string{"A time zone changes local display, not the instant itself.", "Resolve all inputs before returning a complete result.", "Review ~LoadLocation~, then ~In~ and the numeric format offset."}, Files: []LabFile{{Name: "main.go", Starter: twoSunsStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: twoSunsStarterAr}}, Lab: &Lab{ReqFile: "main.go", Hidden: twoSunsHidden, Bonus: true}})
	add(&Level{ID: "19-3-countdown", Chapter: 19, LessonRevision: 1, Steps: ch19CountdownArabic, Title: "العد التنازلي", Goal: "شغّل عداد قابل للإلغاء بساعة محقونة ونظّف الـ ticker", Requires: []string{"sel:context.Context", "select"}, Bridge: []int{15},
		Lesson: `# الإطلاق من غير نوم

## قبل ما تكتب
- مين هيوقف ticker لو الشاشة فشلت؟
- هل الصفر محتاج ticker؟

## للـ ticker مالك
شاشة الإطلاق بتعد ticks، مش الوقت الحقيقي. الساعة المحقونة وAdvance من16-5. هنا ~labkit.Clock.NewTicker~ بترجع ~labkit.Ticker~: ~C()~ قناة ticks و~Stop()~ تحرّر المصدر. Stop مبتقفلش القناة.
~~~go
clock:=labkit.NewFakeClock(time.Unix(0,0))
ticker:=clock.NewTicker(250*time.Millisecond)
defer ticker.Stop()
clock.Advance(250*time.Millisecond)
<-ticker.C()
~~~
اعمل ticker بعد التحقق ونجاح العرض الأول بس. defer Stop بعدها فوراً عشان الإلغاء وخطأ العرض يقفلوها مرة واحدة. FakeClock بتسقط ticks الزيادة لو القناة مليانة زي ticker حقيقي. كل tick تستلمها تنقص واحد حتى بعد Advance كبيرة.

## إلغاء وعرض ممكن يفشل
~context.Context~ وDone وErr وselect من13-4. افحص Err قبل الشغل وبعد استلام tick: select مش بتدي أولوية للإلغاء لما الاتنين جاهزين.
~~~go
select {
case <-ctx.Done(): return ctx.Err()
case <-ticker.C():
 if err:=ctx.Err();err!=nil {return err}
}
~~~
الـ callback function value زي5-5، ولازم ترجع بسرعة. رجّع خطأها الأصلي بدل إخفاء فشل الشاشة.
~~~go
emit:=func(n int)error {fmt.Println(n);return nil}
err:=emit(7)
~~~
~labkit.LeakCheck(t)~ بتسجل cleanup تفحص goroutines؛ ناديها قبل تشغيل الشغل. دي مكملة لعدّ Stop لأن fake ticker ممكن متعملش goroutine أصلاً.
~~~go
labkit.LeakCheck(t)
~~~
التست يستنى أول callback وإنشاء ticker قبل Advance، وبعدها كل قيمة معروضة. مفيش نوم عشان الـ scheduler.

## أخطاء شائعة
~~~go
for range ticker.C() {} // Stop مش بتقفل القناة
time.Sleep(time.Second) // تجاهلت الساعة المحقونة
~~~

## المهمة
نفّذ ~Countdown(ctx context.Context,clock labkit.Clock,seconds int,emit func(int)error) error~. ctx غير nil؛ ارفض nil clock أو nil emit أو seconds<0 قبل أي أثر. بعدها ctx الملغي يرجّع Err من غير عرض أو ticker. اعرض seconds فوراً؛ الصفر ينتهي من غير ticker. غير كده ticker مدتها ثانية واحدة، وكل tick تعرض القيمة الأقل لحد0 شامل. افحص الإلغاء قبل كل عرض ناتج عن tick. رجّع خطأ callback زي ما هو أو ctx.Err؛ بعد الصفر الناجح رجّع nil. كل ticker مملوكة تتوقف مرة واحدة على كل مسار. callback سريعة؛ مفيش goroutine مطلوبة جوا الدالة ولا انتظار حقيقي ولا Report.
`, Hints: []string{"الـclock المحقونة بتخلي الوقت قابل للتحكم من غير انتظار حقيقي.", "فكّر في ملكية ticker ومسار الإلغاء قبل كتابة الـloop.", "راجع ~Ticker.Stop~ و~ctx.Err~ بعد استلام tick."},
		En: &LevelText{LessonRevision: 1, Steps: ch19CountdownEnglish, Title: "Countdown", Goal: "Run a cancelable countdown with injected time and ticker cleanup", Lesson: `# Launch Without Sleeping

## Before you write
- Who stops the ticker if the display fails?
- Does zero need a ticker?

## A ticker has an owner
The launch display counts ticks, not wall-clock time. Clock injection and Advance came from 16-5. Now use ~labkit.Clock.NewTicker~: it returns a ~labkit.Ticker~, whose ~C()~ receives ticks and ~Stop()~ releases the source. Its channel is not closed by Stop.
~~~go
clock:=labkit.NewFakeClock(time.Unix(0,0))
ticker:=clock.NewTicker(250*time.Millisecond)
defer ticker.Stop()
clock.Advance(250*time.Millisecond)
<-ticker.C()
~~~
Create the ticker only after validation and the initial display succeeds. Defer Stop immediately after creation so cancellation and callback errors release it exactly once. FakeClock drops extra ticks when its channel is full, like a real ticker. One received tick means one decrement, even after a large Advance.

## Cancellation and fallible output
~context.Context~, Done, Err, and select come from 13-4. Test Err before starting work and again after receiving a tick: select does not prioritize cancellation when both cases are ready.
~~~go
select {
case <-ctx.Done(): return ctx.Err()
case <-ticker.C():
 if err:=ctx.Err();err!=nil {return err}
}
~~~
The callback is an ordinary function value from 5-5. It must return promptly. Return its original error; do not hide a failed display.
~~~go
emit:=func(n int)error {fmt.Println(n);return nil}
err:=emit(7)
~~~
~labkit.LeakCheck(t)~ registers a test cleanup checking goroutines; call it before launching work. It complements explicit Stop-count checks because a fake ticker need not create a goroutine.
~~~go
labkit.LeakCheck(t)
~~~
Synchronize tests on the initial callback and ticker creation before Advance, then on each emitted value. Scheduler sleeps are unnecessary.

## Common mistakes
~~~go
for range ticker.C() {} // Stop does not close the channel
time.Sleep(time.Second) // Ignores injected time
~~~

## Your mission
Implement ~Countdown(ctx context.Context,clock labkit.Clock,seconds int,emit func(int)error) error~. ctx is non-nil. Reject nil clock, nil emit, or seconds<0 before side effects. Then an already-canceled ctx returns Err without output or a ticker. Emit seconds immediately; zero completes without a ticker. Otherwise create a one-second ticker and emit one lower value per received tick through zero inclusive. Check cancellation before every tick-driven emission. Return the original callback error or ctx.Err; return nil after successful zero emission. Stop every owned ticker exactly once on every exit. Callbacks return promptly. No internal goroutine, real waiting, or Report is needed.
`, Hints: []string{"An injected clock makes time controllable without real waiting.", "Consider ticker ownership and cancellation paths before writing the loop.", "Review ~Ticker.Stop~ and ~ctx.Err~ after receiving a tick."}, Files: []LabFile{{Name: "main.go", Starter: countdownStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: countdownStarterAr}}, Lab: &Lab{ReqFile: "main.go", Hidden: countdownHidden, HiddenFiles: map[string]string{"clock_helpers_test.go": cosmicClockHidden}}})
	add(&Level{ID: "19-4-many-timers", Chapter: 19, LessonRevision: 1, Steps: ch19TimersArabic, Title: "مؤقتات كتير", Goal: "انتظر الأرصفة بالتوازي واعرض النتائج بترتيب الطلبات", Requires: []string{"go", "waitgroup"}, Bridge: []int{31},
		Lesson: `# الأرصفة بتخلص بترتيب مختلف

## قبل ما تكتب
- هتمنع append المتزامنة إزاي؟
- هتكتشف آخر طلب بايظ قبل بدء أول مؤقت إزاي؟

## عامل وخانة وticker لكل طلب
طلبات الأرصفة جاية بترتيب العرض؛ ترتيب الانتهاء ميغيروش. استخدم goroutines و~sync.WaitGroup~ من10-2، وكل عامل يكتب خانة مختلفة مجهزة من الأول. عدّي index والطلب صراحة ومتعملش append متزامن.
~~~go
results:=make([]int,len(jobs))
var wg sync.WaitGroup
for i,job:=range jobs {
 wg.Add(1)
 go func(i int,job int){defer wg.Done();results[i]=job*2}(i,job)
}
wg.Wait()
~~~
المثال jobs فيه []int. عناصر slice المختلفة أماكن مستقلة؛ Wait بتضمن إن الكتابات ظهرت قبل الرجوع. العامل يثبّت نتيجته قبل Done. مفيش داعي لترتيب رسائل الوصول.

## حدث واحد قابل للإلغاء
استخدم ticker من19-3 لحدث واحد وبعده Stop دايماً. Clock.After مش بيلغي fake timer معلقة لما الإلغاء يكسب؛ Clock المشتركة مفيهاش constructor لمؤقت قابل للإيقاف.
~~~go
ticker:=clock.NewTicker(delay)
defer ticker.Stop()
select {
case <-ctx.Done():
case <-ticker.C():
 if ctx.Err()==nil {status="done"}
}
~~~
ابدأ status بـcanceled. كل عامل يفحص الإلغاء قبل الصفر أو إنشاء ticker. لو tick والإلغاء جاهزين مع بعض، canceled تكسب لو العامل شاف Err. نتيجة اتثبتت done تفضل كده بعد إلغاء لاحق.

تحقق من الدفعة كلها قبل تشغيل أي عامل. كشف التكرار map من فصل6.
~~~go
seen:=map[string]bool{}
if seen[id] {return nil,fmt.Errorf("duplicate dock")}
seen[id]=true
~~~
التستات FakeClock وإشارات إنشاء وعدّ Stop وLeakCheck من19-3. المختبر بيشغّل race detector. الوقت بيتحرك بعد تسجيل كل المصادر الموجبة عشان يثبت توازي العمال.

## أخطاء شائعة
~~~go
go func(){results=append(results,reply)}() // race وترتيب وصول
<-clock.After(delay) // مش قابل للإلغاء والتنظيف
~~~

## المهمة
نفّذ ~AwaitDocks(ctx context.Context,clock labkit.Clock,jobs []DockTimer) ([]DockResult,error)~. DockTimer حقولها ID string وDelay time.Duration؛ DockResult حقولها ID,Status string. ctx غير nil. nil clock أو ID فاضي أو مكرر أو Delay سالبة يرفضوا الدفعة كلها بـnil,error من غير جدولة. متعدلش jobs. الفاضي نتيجة طول0,nil. عامل لكل طلب يكتب نفس index مرة واحدة، وWaitGroup تستنى الكل. العامل يشوف ctx.Err أولاً: الملغي نتيجته canceled حتى للصفر؛ غير كده صفر done فوراً من غير ticker. الموجب ticker بمدة Delay؛ أول tick وErr==nil يثبّت done، والإلغاء يثبّت canceled. كل ticker تتوقف مرة واحدة. بعد اكتمال العمال ارجع النتائج,nil حتى لو بعضهم ملغي. نتيجة done متتغيرش بإلغاء لاحق؛ إلغاء متزامن قبل تثبيت النتيجة ممكن يكسب. مفيش انتظار حقيقي ولا Report؛ Race شغالة.
`, Hints: []string{"ترتيب انتهاء العمال مش هو ترتيب الطلبات.", "حدّد ملكية كل نتيجة وطريقة انتظار المجموعة قبل تشغيل العمال.", "راجع إن ~WaitGroup.Done~ بعد الكتابة وإن ~Wait~ قبل القراءة."},
		En: &LevelText{LessonRevision: 1, Steps: ch19TimersEnglish, Title: "Many Timers", Goal: "Await dock timers concurrently and display results in request order", Lesson: `# Docks Finish Out of Order

## Before you write
- How will you avoid concurrent append?
- How will you reject the final invalid job before scheduling the first?

## One worker, one slot, one ticker
Dock requests arrive in display order; completion order must not reorder them. Reuse goroutines and ~sync.WaitGroup~ from 10-2, now assigning each worker a distinct preallocated result slot. Pass index and job explicitly and never append concurrently.
~~~go
results:=make([]int,len(jobs))
var wg sync.WaitGroup
for i,job:=range jobs {
 wg.Add(1)
 go func(i int,job int){defer wg.Done();results[i]=job*2}(i,job)
}
wg.Wait()
~~~
The sample assumes jobs is []int. Distinct slice elements are separate storage; Wait makes all worker writes visible before returning. A worker's result is committed before it calls Done. Sorting completion messages is unnecessary.

## Cancelable one-shot ownership
Use the ticker API from 19-3 for a single event and always Stop it. Clock.After cannot release an outstanding fake timer when cancellation wins; the shared Clock interface has no stoppable timer constructor.
~~~go
ticker:=clock.NewTicker(delay)
defer ticker.Stop()
select {
case <-ctx.Done():
case <-ticker.C():
 if ctx.Err()==nil {status="done"}
}
~~~
Initialize status to canceled. A worker observes cancellation before handling zero or creating a ticker. If tick and cancellation are ready together, canceled wins when the worker observes Err. A completed slot stays done after later cancellation.

Validate the entire batch before starting any worker. Duplicate detection uses a map from chapter 6.
~~~go
seen:=map[string]bool{}
if seen[id] {return nil,fmt.Errorf("duplicate dock")}
seen[id]=true
~~~
Tests use FakeClock, creation acknowledgements, Stop counts and LeakCheck from 19-3. The lab runs the race detector. Time advances only after every positive-delay source is registered, proving workers are parallel.

## Common mistakes
~~~go
go func(){results=append(results,reply)}() // Race and completion order
<-clock.After(delay) // Cannot cancel and release the source
~~~

## Your mission
Implement ~AwaitDocks(ctx context.Context,clock labkit.Clock,jobs []DockTimer) ([]DockResult,error)~. DockTimer has ID string and Delay time.Duration; DockResult has ID,Status string. ctx is non-nil. A nil clock, empty or duplicate ID, or negative Delay rejects the whole batch with nil,error before scheduling. Do not mutate jobs. Empty input returns length zero,nil. Each job has a worker writing its original index exactly once; WaitGroup joins all workers. A worker checks ctx.Err first: canceled yields canceled even for zero; otherwise zero yields done without a ticker. Positive delays use a ticker of Delay; the first tick with nil Err commits done, while cancellation commits canceled. Stop each ticker exactly once. Return results,nil after joining, even for canceled entries. A committed done survives later cancellation; concurrent cancellation before commitment may win. No real waiting or Report; the race detector is enabled.
`, Hints: []string{"Worker finish order is not request order.", "Decide who owns each result and how workers will be joined.", "Check that ~WaitGroup.Done~ follows the write and ~Wait~ precedes reading."}, Files: []LabFile{{Name: "main.go", Starter: manyTimersStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: manyTimersStarterAr}}, Lab: &Lab{ReqFile: "main.go", Hidden: manyTimersHidden, HiddenFiles: map[string]string{"clock_helpers_test.go": cosmicClockHidden}}, Race: true})
}

const flightDurationsStarter = `package main

import "orbital/colony"

func main() {
	for range colony.Orders() {
		colony.Report("unfinished")
	}
}
`
const flightDurationsStarterAr = `package main

import "orbital/colony"

func main() {
	for range colony.Orders() {
		colony.Report("لسه مخلصتش")
	}
}
`
const twoSunsHidden = `package main

import (
	"io"
	"orbital/colony/levels"
	"strings"
	"testing"
	"time"
	_ "time/tzdata"
)

func TestHiddenSunsFixedInstants(t *testing.T) {
	w := levels.Start("19-2-two-suns", io.Discard, 192)
	zones := []string{"UTC", "Asia/Tokyo", "Pacific/Honolulu", "Asia/Dubai"}
	for i := 0; i < 80; i++ {
		instant := time.Unix(1735689600+int64(w.Rand.Intn(86400*60)), 0).UTC()
		home, remote := zones[w.Rand.Intn(len(zones))], zones[w.Rand.Intn(len(zones))]
		a, _ := time.LoadLocation(home)
		b, _ := time.LoadLocation(remote)
		x, y := instant.In(a), instant.In(b)
		want := SunReadings{x.Format("2006-01-02 15:04 -0700"), y.Format("2006-01-02 15:04 -0700"), x.Format("2006-01-02") != y.Format("2006-01-02")}
		got, err := ReadSuns(instant.Format(time.RFC3339), home, remote)
		if err != nil || got != want {
			t.Fatalf("%s %s -> %s: %+v want %+v: %v", instant, home, remote, got, want, err)
		}
	}
	for _, pair := range [][2]string{{"Asia/Tokyo", "Pacific/Honolulu"}, {"Pacific/Honolulu", "Asia/Tokyo"}} {
		got, err := ReadSuns("2025-01-01T00:00:00Z", pair[0], pair[1])
		if err != nil || !got.DayChanged {
			t.Fatal("date boundary", got, err)
		}
	}
	for _, tc := range []struct{ instant, home, remote, part string }{
		{"2025-02-30T00:00:00Z", "UTC", "UTC", "instant"}, {"bad", "UTC", "UTC", "instant"},
		{"2025-01-01T00:00:00Z", "Mars/Missing", "UTC", "Mars/Missing"},
		{"2025-01-01T00:00:00Z", "UTC", "Mars/Missing", "Mars/Missing"},
	} {
		got, err := ReadSuns(tc.instant, tc.home, tc.remote)
		if got != (SunReadings{}) || err == nil || !strings.Contains(err.Error(), tc.part) {
			t.Fatal("invalid input", got, err)
		}
	}
}
func TestHiddenSunsInputOffset(t *testing.T) {
	got, err := ReadSuns("2025-01-01T04:00:00+04:00", "UTC", "Asia/Tokyo")
	want := SunReadings{"2025-01-01 00:00 +0000", "2025-01-01 09:00 +0900", false}
	if err != nil || got != want {
		t.Fatal("input offset must preserve the instant", got, err)
	}
}
func TestHiddenBonusSeasonalSuns(t *testing.T) {
	for _, tc := range []struct{ instant, want string }{{"2025-01-15T12:00:00Z", "2025-01-15 13:00 +0100"}, {"2025-07-15T12:00:00Z", "2025-07-15 14:00 +0200"}} {
		got, err := ReadSuns(tc.instant, "UTC", "Europe/Berlin")
		if err != nil || got.Remote != tc.want {
			t.Fatal("seasonal offset", got, err)
		}
	}
}
`
const twoSunsStarter = `package main

import "errors"

type SunReadings struct {
	Home, Remote string
	DayChanged   bool
}

func ReadSuns(instant, home, remote string) (SunReadings, error) {
	return SunReadings{}, errors.New("suns unfinished")
}
func main() {}
`
const twoSunsStarterAr = `package main

import "errors"

type SunReadings struct {
	Home, Remote string
	DayChanged   bool
}

func ReadSuns(instant, home, remote string) (SunReadings, error) {
	return SunReadings{}, errors.New("شمسين لسه مخلصتش")
}
func main() {}
`
const countdownHidden = `package main

import (
	"context"
	"errors"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"testing"
	"time"
)

func TestHiddenCountdownTicksAndCleanup(t *testing.T) {
	labkit.LeakCheck(t)
	w := levels.Start("19-3-countdown", io.Discard, 193)
	for repeat := 0; repeat < 20; repeat++ {
		n := 1 + w.Rand.Intn(6)
		clock := hiddenNewClock()
		emitted := make(chan int, 20)
		done := make(chan error, 1)
		ctx, cancel := context.WithCancel(context.Background())
		go func() { done <- Countdown(ctx, clock, n, func(n int) error { emitted <- n; return nil }) }()
		if hiddenReceive(t, emitted) != n || hiddenReceive(t, clock.created) != time.Second {
			t.Fatal("initial value or interval")
		}
		clock.Advance(time.Second - time.Nanosecond)
		select {
		case n := <-emitted:
			t.Fatal("early emission", n)
		default:
		}
		clock.Advance(time.Nanosecond)
		if hiddenReceive(t, emitted) != n-1 {
			t.Fatal("exact boundary")
		}
		for want := n - 2; want >= 0; want-- {
			clock.Advance(time.Second)
			if hiddenReceive(t, emitted) != want {
				t.Fatal("count order")
			}
		}
		if err := hiddenReceive(t, done); err != nil {
			t.Fatal(err)
		}
		cancel()
		clock.checkStopped(t, 1)
		select {
		case extra := <-emitted:
			t.Fatal("extra result", extra)
		default:
		}
	}
}
func TestHiddenCountdownFailureAndCancellation(t *testing.T) {
	labkit.LeakCheck(t)
	wire := errors.New("display failed")
	for _, failAt := range []int{3, 2, 0} {
		clock := hiddenNewClock()
		emitted := make(chan int, 10)
		done := make(chan error, 1)
		go func() {
			done <- Countdown(context.Background(), clock, 3, func(n int) error {
				emitted <- n
				if n == failAt {
					return wire
				}
				return nil
			})
		}()
		if hiddenReceive(t, emitted) != 3 {
			t.Fatal("initial value")
		}
		if failAt != 3 {
			hiddenReceive(t, clock.created)
			for n := 2; n >= failAt; n-- {
				clock.Advance(time.Second)
				if hiddenReceive(t, emitted) != n {
					t.Fatal("callback order")
				}
			}
		}
		if !errors.Is(hiddenReceive(t, done), wire) {
			t.Fatal("callback cause lost")
		}
		want := 1
		if failAt == 3 {
			want = 0
		}
		clock.checkStopped(t, want)
	}
	for _, afterTick := range []bool{false, true} {
		clock := hiddenNewClock()
		emitted := make(chan int, 10)
		done := make(chan error, 1)
		ctx, cancel := context.WithCancel(context.Background())
		go func() { done <- Countdown(ctx, clock, 4, func(n int) error { emitted <- n; return nil }) }()
		hiddenReceive(t, emitted)
		hiddenReceive(t, clock.created)
		if afterTick {
			clock.Advance(time.Second)
			if hiddenReceive(t, emitted) != 3 {
				t.Fatal("tick")
			}
		}
		cancel()
		if !errors.Is(hiddenReceive(t, done), context.Canceled) {
			t.Fatal("cancellation cause")
		}
		clock.checkStopped(t, 1)
	}
	clock := hiddenNewClock()
	calls := 0
	emit := func(int) error { calls++; return nil }
	if err := Countdown(context.Background(), clock, 0, emit); err != nil || calls != 1 {
		t.Fatal("zero countdown", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Countdown(ctx, clock, 2, emit); !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatal("pre-canceled callback")
	}
	for _, err := range []error{Countdown(context.Background(), nil, 2, emit), Countdown(context.Background(), clock, -1, emit), Countdown(context.Background(), clock, 2, nil)} {
		if err == nil {
			t.Fatal("invalid config")
		}
	}
	if calls != 1 || len(clock.created) != 0 {
		t.Fatal("no-source path scheduled work")
	}
	clock.checkStopped(t, 0)
}
`
const countdownStarter = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
)

func Countdown(ctx context.Context, clock labkit.Clock, seconds int, emit func(int) error) error {
	return errors.New("countdown unfinished")
}
func main() {}
`
const countdownStarterAr = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
)

func Countdown(ctx context.Context, clock labkit.Clock, seconds int, emit func(int) error) error {
	return errors.New("العداد لسه مخلصش")
}
func main() {}
`
const manyTimersHidden = `package main

import (
	"context"
	"fmt"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"slices"
	"testing"
	"time"
)

type hiddenDockReply struct {
	results []DockResult
	err     error
}

func TestHiddenDocksParallelAndOrdered(t *testing.T) {
	labkit.LeakCheck(t)
	w := levels.Start("19-4-many-timers", io.Discard, 194)
	for repeat := 0; repeat < 15; repeat++ {
		clock := hiddenNewClock()
		var jobs []DockTimer
		var want []DockResult
		for i := 0; i < 10; i++ {
			id := fmt.Sprintf("bay-%d", i)
			jobs = append(jobs, DockTimer{id, time.Duration(1+w.Rand.Intn(5)) * time.Second})
			want = append(want, DockResult{id, "done"})
		}
		before := slices.Clone(jobs)
		done := make(chan hiddenDockReply, 1)
		go func() { r, e := AwaitDocks(context.Background(), clock, jobs); done <- hiddenDockReply{r, e} }()
		for range jobs {
			hiddenReceive(t, clock.created)
		}
		clock.Advance(5 * time.Second)
		reply := hiddenReceive(t, done)
		if reply.err != nil || !slices.Equal(reply.results, want) || !slices.Equal(jobs, before) {
			t.Fatal("ordered independent slots", reply)
		}
		clock.checkStopped(t, len(jobs))
	}
}
func TestHiddenDocksCancelAndPreflight(t *testing.T) {
	labkit.LeakCheck(t)
	clock := hiddenNewClock()
	done := make(chan hiddenDockReply, 1)
	ctx, cancel := context.WithCancel(context.Background())
	jobs := []DockTimer{{"slow", 5 * time.Second}, {"fast", time.Second}, {"zero", 0}}
	go func() { r, e := AwaitDocks(ctx, clock, jobs); done <- hiddenDockReply{r, e} }()
	hiddenReceive(t, clock.created)
	hiddenReceive(t, clock.created)
	clock.Advance(time.Second)
	// Completion is committed before Stop; the clock observes that boundary.
	// A wrapper acknowledgement makes cancellation independent of scheduling.
	hiddenWaitStopped(t, clock, 1)
	cancel()
	reply := hiddenReceive(t, done)
	if reply.err != nil || len(reply.results) != 3 || reply.results[0] != (DockResult{"slow", "canceled"}) || reply.results[1] != (DockResult{"fast", "done"}) {
		t.Fatal("mixed cancellation", reply)
	}
	if reply.results[2].ID != "zero" || (reply.results[2].Status != "done" && reply.results[2].Status != "canceled") {
		t.Fatal("zero result")
	}
	clock.checkStopped(t, 2)
	for _, jobs := range [][]DockTimer{{{"a", time.Second}, {"a", 0}}, {{"a", time.Second}, {"b", -1}}, {{"", 0}}} {
		c := hiddenNewClock()
		r, e := AwaitDocks(context.Background(), c, jobs)
		if e == nil || r != nil || len(c.created) != 0 {
			t.Fatal("preflight must schedule nothing", r, e)
		}
	}
	c := hiddenNewClock()
	if r, e := AwaitDocks(context.Background(), nil, nil); e == nil || r != nil {
		t.Fatal("nil clock")
	}
	if r, e := AwaitDocks(context.Background(), c, nil); e != nil || len(r) != 0 {
		t.Fatal("empty jobs")
	}
	if r, e := AwaitDocks(context.Background(), c, []DockTimer{{"z", 0}}); e != nil || !slices.Equal(r, []DockResult{{"z", "done"}}) {
		t.Fatal("zero timer", r, e)
	}
	canceled, stop := context.WithCancel(context.Background())
	stop()
	if r, e := AwaitDocks(canceled, c, []DockTimer{{"a", time.Second}, {"b", 0}}); e != nil || !slices.Equal(r, []DockResult{{"a", "canceled"}, {"b", "canceled"}}) {
		t.Fatal("pre-canceled", r, e)
	}
	if len(c.created) != 0 {
		t.Fatal("immediate paths created ticker")
	}
}
`
const manyTimersStarter = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
	"time"
)

type DockTimer struct {
	ID    string
	Delay time.Duration
}
type DockResult struct{ ID, Status string }

func AwaitDocks(ctx context.Context, clock labkit.Clock, jobs []DockTimer) ([]DockResult, error) {
	return nil, errors.New("timers unfinished")
}
func main() {}
`
const manyTimersStarterAr = `package main

import (
	"context"
	"errors"
	"orbital/colony/labkit"
	"time"
)

type DockTimer struct {
	ID    string
	Delay time.Duration
}
type DockResult struct{ ID, Status string }

func AwaitDocks(ctx context.Context, clock labkit.Clock, jobs []DockTimer) ([]DockResult, error) {
	return nil, errors.New("المؤقتات لسه مخلصتش")
}
func main() {}
`
const cosmicClockHidden = `package main

import (
	"orbital/colony/labkit"
	"sync"
	"testing"
	"time"
)

type hiddenClock struct {
	*labkit.FakeClock
	mu            sync.Mutex
	created       chan time.Duration
	active, stops int
	stopped       chan struct{}
}
type hiddenTicker struct {
	labkit.Ticker
	owner *hiddenClock
}

func hiddenNewClock() *hiddenClock {
	return &hiddenClock{FakeClock: labkit.NewFakeClock(time.Unix(0, 0)), created: make(chan time.Duration, 100), stopped: make(chan struct{}, 100)}
}
func (c *hiddenClock) NewTicker(d time.Duration) labkit.Ticker {
	ticker := c.FakeClock.NewTicker(d)
	c.mu.Lock()
	c.active++
	c.mu.Unlock()
	c.created <- d
	return &hiddenTicker{ticker, c}
}
func (t *hiddenTicker) Stop() {
	t.Ticker.Stop()
	t.owner.mu.Lock()
	t.owner.active--
	t.owner.stops++
	t.owner.mu.Unlock()
	t.owner.stopped <- struct{}{}
}
func hiddenReceive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case n := <-ch:
		return n
	// This is only a deadlock watchdog; successful tests advance fake time.
	case <-time.After(3 * time.Second):
		t.Fatal("deadlock watchdog; fake event was not consumed")
		var zero T
		return zero
	}
}
func (c *hiddenClock) checkStopped(t *testing.T, want int) {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active != 0 || c.stops != want {
		t.Fatalf("ticker lifecycle: active=%d stops=%d want=%d", c.active, c.stops, want)
	}
}

func hiddenWaitStopped(t *testing.T, c *hiddenClock, n int) {
	t.Helper()
	for i := 0; i < n; i++ {
		hiddenReceive(t, c.stopped)
	}
}
`

func setupFlightDurations(w *sim.World) {
	w.Init(4, 2, sim.Regolith)
	var want []string
	for i := 0; i < 18; i++ {
		start := time.Date(2030+w.Rand.Intn(4), time.December, 31, 23, w.Rand.Intn(60), w.Rand.Intn(60), 0, time.UTC)
		duration := time.Duration(w.Rand.Intn(172801)) * time.Second
		if i == 0 {
			duration = 0
		}
		if i == 1 {
			duration = 48 * time.Hour
		}
		if i == 2 {
			duration = 1500 * time.Millisecond
		}
		spelling := duration.String()
		if i%3 == 0 {
			spelling = fmt.Sprintf("%ds", int64(duration/time.Second))
		}
		w.Data.Orders = append(w.Data.Orders, start.Format("2006-01-02 15:04:05")+"|"+spelling)
		want = append(want, start.Add(duration).Format("2006-01-02 15:04:05")+"|"+duration.String())
	}
	for _, bad := range []string{"2031-02-30 12:00:00|1h", "2031-01-01 00:00:00|-1s", "2031-01-01 00:00:00|48h1ns", "2031-01-01 00:00:00|nope", "bad|1h", "2031-01-01 00:00:00|999999999999999999999h"} {
		w.Data.Orders = append(w.Data.Orders, bad)
		want = append(want, "invalid")
	}
	w.Rand.Shuffle(len(want), func(i, j int) {
		want[i], want[j] = want[j], want[i]
		w.Data.Orders[i], w.Data.Orders[j] = w.Data.Orders[j], w.Data.Orders[i]
	})
	w.Check = func(w *sim.World) (bool, string) {
		if slices.Equal(w.Reports, want) {
			return true, i18n.T("القاعدة استلمت كل مواعيد الوصول والمدد صح", "Base received every correct arrival and duration")
		}
		return false, i18n.T("مستني تقرير لكل رحلة بالترتيب؛ راجع المدة والتاريخ البايظ", "Waiting for one report per flight in order; check duration and invalid dates")
	}
}
