package levels

func init() {
	add(&Level{ID: "38-1-iterators", Chapter: 38, Kind: "lab", Title: "مسارات من غير تخزين", Goal: "اكتب iterators بتقف فوراً لما المستهلك يوقف", Setup: seededLabRoom,
		Requires: []string{"sel:iter.Seq", "sel:iter.Seq2", "range"}, Files: []LabFile{{Name: "main.go", Starter: ch381Starter}},
		Lesson: `# مسارات قراءات المستعمرة

## قبل ما تكتب
- لو الحلقة عملت ~break~ بعد تاني قراءة، هل المنتج يحق له يجهّز التالتة؟
- هل الـ iterator بيتشغّل ساعة إنشائه ولا ساعة استهلاكه؟

## دالة بدل قناة
من Go 1.23، ~range~ يقدر يستهلك دالة. ~iter.Seq[V]~ اسم لنوع ~func(yield func(V) bool)~؛ ~iter.Seq2[K,V]~ نفس الفكرة بقيمتين، زي index وقيمة. ده مش goroutine: المستهلك والمنتج بيتبادلوا التحكم على نفس مسار التنفيذ.
~~~go
func countdown() iter.Seq[int] {
 return func(yield func(int) bool) {
  for n := 3; n > 0; n-- {
   if !yield(n) { return }
  }
 }
}
for n := range countdown() { fmt.Println(n) }
~~~
~yield~ بيرجع false لما المستهلك يعمل ~break~ أو ~return~. لازم ترجع فوراً؛ محاولة yield تانية ممكن تعمل panic. مفيش channel تتقفل ولا goroutine تتسرّب. إنشاء الدالة ماينفّذش جسمها. كل استهلاك جديد يبدأ من الأول لو حالة المرور جوه الدالة الداخلية.

## قيمتين في كل لفة
~~~go
var pairs iter.Seq2[string, int] = func(yield func(string, int) bool) {
 if !yield("oxygen", 4) { return }
}
for name, count := range pairs { fmt.Println(name, count) }
~~~
المثال بيستخدم مفتاح string، لكن المهمة عايزة index من صفر. خليك محافظ على ترتيب المدخل، وماتعدّلوش.

## أخطاء شائعة
~~~go
for _, v := range values {
 visit(v)
 yield(v) // غلط: تجاهلت طلب التوقف
}
~~~
كمان تجميع كل القيم قبل أول yield بيضيّع الكسل (laziness)، وتحويلها لـ channel محتاج إدارة عمر goroutine من غير لازمة هنا.

## المهمة
نفّذ ~Readings(values []int, visit func(int)) iter.Seq[int]~: قبل كل yield مباشرة نادِ ~visit~ مرة للقيمة دي. ممنوع تزور قيمة بعد توقف المستهلك. ~visit~ غير nil.
نفّذ ~Indexed(values []int) iter.Seq2[int,int]~ يرجّع index والقيمة بالترتيب، ويطيع التوقف فوراً.
نفّذ ~First(values []int, n int) []int~ باستخدام range-over-func على ~Readings~؛ يرجّع أول n قيم أو كل القيم لو أقل. n صفر أو سالب يرجّع نتيجة فاضية. مفيش تعديل للمدخل ولا goroutines. نجمة إضافية: نفس ~Readings~ يتستهلك أكتر من مرة بعد break ويبدأ من الأول. ~main~ فاضية.
`,
		Hints: []string{"ارجع فوراً لما yield ترجع false؛ الزيارة تحصل قبلها مباشرة.", "خلّي index جوه الدالة الداخلية عشان يبدأ من صفر كل مرة.", "في First افحص n قبل الحلقة واستخدم break لما تجمع العدد المطلوب."},
		En: &LevelText{Title: "Paths Without Storage", Goal: "Write iterators that stop as soon as their consumer stops", Files: []LabFile{{Name: "main.go", Starter: ch381Starter}}, Lesson: `# Colony Reading Paths

## Before you write
- If the consumer breaks after its second reading, may the producer prepare a third?
- Does an iterator run when constructed or when consumed?

## A function instead of a channel
Since Go 1.23, ~range~ can consume a function. ~iter.Seq[V]~ names ~func(yield func(V) bool)~; ~iter.Seq2[K,V]~ yields two values, such as an index and a reading. This is synchronous: producer and consumer transfer control on the same execution path without a goroutine.
~~~go
func countdown() iter.Seq[int] {
 return func(yield func(int) bool) {
  for n := 3; n > 0; n-- {
   if !yield(n) { return }
  }
 }
}
for n := range countdown() { fmt.Println(n) }
~~~
~yield~ returns false when the consumer uses ~break~ or ~return~. Return immediately; another yield can panic. There is no channel to close or goroutine to leak. Constructing the function does not execute its body. Put traversal state inside the inner function to restart on each consumption.

## Two values per iteration
~~~go
var pairs iter.Seq2[string, int] = func(yield func(string, int) bool) {
 if !yield("oxygen", 4) { return }
}
for name, count := range pairs { fmt.Println(name, count) }
~~~
This example uses a string key; your task uses an index starting at zero. Preserve input order and leave the input unchanged.

## Common mistakes
~~~go
for _, v := range values {
 visit(v)
 yield(v) // Wrong: ignores the consumer's stop request
}
~~~
Collecting everything before the first yield loses laziness. Introducing a channel adds an unnecessary goroutine lifetime to manage here.

## Mission
Implement ~Readings(values []int, visit func(int)) iter.Seq[int]~. Call the non-nil ~visit~ exactly once immediately before each yield, never for a value after the consumer stops.
Implement ~Indexed(values []int) iter.Seq2[int,int]~, yielding index and value in order and stopping immediately on false.
Implement ~First(values []int, n int) []int~ using range-over-func on ~Readings~: return the first n values, or all if shorter. Nonpositive n gives an empty result. Do not mutate inputs or create goroutines. Bonus star: the same ~Readings~ iterator restarts when consumed again after an early break. Leave ~main~ empty.
`, Hints: []string{"Return immediately when yield returns false; visit immediately before yielding.", "Keep traversal state inside the inner function so each iteration restarts.", "In First check n before iterating and break when enough values have been collected."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch381Hidden, Bonus: true},
	})
	add(&Level{ID: "38-2-sets-and-labels", Chapter: 38, Kind: "lab", Title: "بوابات مميّزة", Goal: "استخدم set وlabels لتختار أول بوابة مسموحة في كل صف", Setup: seededLabRoom,
		Requires: []string{"map", "generic"}, Files: []LabFile{{Name: "main.go", Starter: ch382Starter}},
		Lesson: `# اختيار بوابات من غير تكرار

## قبل ما تكتب
- مقابلة بوابة محظورة تلغي المرشح ده بس ولا بقية الصف كله؟
- أول اختيار في صف يوقف أنهي حلقة؟

## struct{} كـ set
~struct{}~ نوع من غير حقول؛ القيمة مش شايلة بيانات. ~map[T]struct{}~ بتخزّن العضوية في المفتاح. ~T comparable~ من فصل generics شرط إن T ينفع مفتاح map. المجموعات مش مرتّبة؛ ترتيب النتيجة ييجي من المرور على الصفوف.
~~~go
allowed := map[string]struct{}{"north": {}}
allowed["south"] = struct{}{}
_, exists := allowed["north"]
delete(allowed, "south")
~~~
استخدم comma-ok للعضوية؛ قيمة ~struct{}{}~ نفسها ماتفرّقش موجود من مش موجود. خلي مجموعة للمحظور ومجموعة للاختيارات السابقة.

## اسم للحلقة
الـ label اسم قبل جملة، و~break Name~ و~continue Name~ بيستهدفوا حلقة محيطة بالاسم. continue للحلقة الكبيرة بيتخطّى بقية الصف كله؛ break للحلقة الصغيرة يكمّل الصف اللي بعده.
~~~go
Decks:
for _, deck := range decks {
Slots:
 for _, slot := range deck {
  if slot < 0 { continue Decks }
  if slot == 0 { continue Slots }
  fmt.Println(slot)
  break Slots
 }
}
~~~
المثال بيفحص أرقام خانات؛ المهمة بتفحص العضوية في مجموعتين. الـ label مش goto ومش اسم function، ومينفعش تقفز لحلقة مش محيطة بيك.

## أخطاء شائعة
~~~go
for _, row := range rows {
 for _, gate := range row {
  if blocked(gate) { continue } // غلط: كده هتكمّل نفس الصف
 }
}
~~~
وممنوع تلف على map عشان تبني النتيجة؛ ترتيب map مش ثابت.

## المهمة
نفّذ ~SelectGates[T comparable](rows [][]T, blocked []T) []T~ باستخدام ~map[T]struct{}~ وlabeled break وcontinue. افحص الصفوف والمرشحين بترتيبهم. لو قابلت محظور قبل ما تختار، سيب بقية الصف. لو قابلت اختيار قديم، جرّب اللي بعده في نفس الصف. أول قيمة جديدة ومسموحة تتضاف مرة واحدة وتوقف البحث في الصف. المحظور اللي بعد اختيار ناجح ميتفحصش. الصف الفاضي ملوش اختيار، والمدخل ميتعدّلش. ينفع int وstring وأي comparable. ~main~ فاضية.
`,
		Hints: []string{"ابني denied من blocked وseen من الاختيارات اللي نجحت.", "افحص المحظور الأول؛ continue للصف الخارجي يلغي الباقي.", "بعد append سجّل العضوية واعمل break للحلقة الداخلية."},
		En: &LevelText{Title: "Distinct Gates", Goal: "Use a set and labels to choose the first allowed gate in each row", Files: []LabFile{{Name: "main.go", Starter: ch382Starter}}, Lesson: `# Choose Gates Without Duplicates

## Before you write
- Does a blocked gate reject just that candidate or the rest of its row?
- Which loop stops after the first successful selection?

## struct{} as a set
~struct{}~ has no fields and carries no payload. ~map[T]struct{}~ represents membership through its keys. The generic constraint ~T comparable~ makes T a valid map key. Sets are unordered; obtain result order by traversing the rows.
~~~go
allowed := map[string]struct{}{"north": {}}
allowed["south"] = struct{}{}
_, exists := allowed["north"]
delete(allowed, "south")
~~~
Use comma-ok for membership: ~struct{}{}~ itself cannot distinguish absence. Keep separate sets for blocked gates and previous selections.

## Name a loop
A label names a statement. ~break Name~ and ~continue Name~ target an enclosing loop. Continuing the outer loop skips the rest of the row; breaking the inner loop moves to the next row.
~~~go
Decks:
for _, deck := range decks {
Slots:
 for _, slot := range deck {
  if slot < 0 { continue Decks }
  if slot == 0 { continue Slots }
  fmt.Println(slot)
  break Slots
 }
}
~~~
This example tests slot numbers; the mission tests membership in two sets. A label is neither goto nor a function name, and cannot target a loop that does not enclose the branch.

## Common mistakes
~~~go
for _, row := range rows {
 for _, gate := range row {
  if blocked(gate) { continue } // Wrong: continues this same row
 }
}
~~~
Do not build the output by ranging over a map: its iteration order is unspecified.

## Mission
Implement ~SelectGates[T comparable](rows [][]T, blocked []T) []T~ using ~map[T]struct{}~ and labeled break and continue. Inspect rows and candidates in input order. A blocked gate encountered before selection skips the rest of that row. A previously selected gate skips only that candidate. Append the first new allowed gate and stop scanning its row. A blocked gate after that successful choice is never examined. Empty rows select nothing. Do not mutate inputs. Support int, string and any comparable type. Leave ~main~ empty.
`, Hints: []string{"Build denied from blocked and seen from successful selections.", "Check blocked membership first; continue the outer loop to discard its remainder.", "After append, record membership and break the inner loop."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch382Hidden},
	})
	add(&Level{ID: "38-3-once", Chapter: 38, Kind: "lab", Title: "تشغيل مرة واحدة", Goal: "هيّأ إعدادات مشتركة مرة واحدة بأمان تحت التزامن", Setup: seededLabRoom,
		Requires: []string{"sel:sync.Once", "sel:sync.OnceValue"}, Files: []LabFile{{Name: "main.go", Starter: ch383Starter}},
		Race: true,
		Lesson: `# معايرة مرة واحدة

## قبل ما تكتب
- لو 40 goroutine طلبوا المعايرة مع بعض، مين ينشر النتيجة؟
- هل القيمة صفر أو النص الفاضي معناهم إن التهيئة ماحصلتش؟

## sync.Once
~sync.Once~ صفرها جاهز. ~Do(func())~ ينفّذ دالة واحدة مرة واحدة لكل instance. الاستدعاءات المتزامنة بتستنى انتهاء نفس التهيئة، وبعد الرجوع بتشوف كل الكتابات اللي حصلت جواها. خزن النتيجة في نفس الـ struct وخلي methods على pointer؛ ماتنسخش Once بعد الاستخدام.
~~~go
type beacon struct { once sync.Once; code string }
func (b *beacon) Code() string {
 b.once.Do(func(){ b.code = "K-7" })
 return b.code
}
~~~
~Once~ بيحمي التهيئة فقط؛ لو هتعدل القيمة بعد كده، محتاج مزامنة تانية. لو الدالة عملت panic، Once بيعتبرها اتنفّذت؛ مش بيعيد المحاولة. نداء نفس Once من جوه Do بتاعتها يعمل deadlock.

## sync.OnceValue
~sync.OnceValue~ موجودة من Go 1.21. بتاخد ~func() T~ وبترجع ~func() T~ آمنة للتزامن، بتخزّن النتيجة حتى لو صفر. كل closure مستقلة.
~~~go
loadPort := sync.OnceValue(func() int { return 8081 })
port := loadPort() // التهيئة هنا، مش وقت إنشاء loadPort
~~~
لو الدالة الأصلية بتعمل panic، كل استدعاء للـ closure بيرمي نفس قيمة panic من غير ما يعيد الأصلية. دي نجمة إضافية في المهمة.

## أخطاء شائعة
~~~go
if c.value == 0 { c.value = load() } // غلط: race، وصفر ممكن نتيجة صحيحة
var once sync.Once
_ = once // لو حطيتها محلية في كل نداء، كل نداء هيهيّأ تاني
~~~

## المهمة
نفّذ ~(c *Calibration) Value(load func() int) int~: أول load بس يتنفّذ عبر ~sync.Once~، وكل المنادين يشوفوا نفس int. كل Calibration مستقلة وصفرها شغّال. load غير nil، ومش بتعمل panic في الجزء الأساسي.
نفّذ ~LazyName(load func() string) func() string~ بـ ~sync.OnceValue~: متستدعيش load وقت إنشاء closure. أول استهلاك يخزّن النتيجة وكل goroutines تشوفها، حتى النص الفاضي. نجمة إضافية: تكرار نفس panic value من غير إعادة load. الاختبارات بتشتغل بـ ~-race~. ~main~ فاضية.
`,
		Hints: []string{"خزّن sync.Once كحقل، مش متغير محلي في Value.", "خلي الكتابة للنتيجة جوه Do والقراية بعدها؛ pointer receiver يمنع النسخ.", "رجّع closure بتاعة sync.OnceValue؛ ماتستدعيهاش في LazyName."},
		En: &LevelText{Title: "Start Once", Goal: "Initialize shared settings exactly once under concurrency", Files: []LabFile{{Name: "main.go", Starter: ch383Starter}}, Lesson: `# Calibrate Once

## Before you write
- If forty goroutines request calibration together, who publishes the result?
- Does zero or an empty string imply initialization has not happened?

## sync.Once
The zero value of ~sync.Once~ is usable. ~Do(func())~ runs one function exactly once per instance. Concurrent callers wait for that initialization and see all its writes after returning. Store the result in the same struct and use pointer receivers; never copy a Once after use.
~~~go
type beacon struct { once sync.Once; code string }
func (b *beacon) Code() string {
 b.once.Do(func(){ b.code = "K-7" })
 return b.code
}
~~~
Once guards initialization only; later mutations need separate synchronization. If initialization panics, Once considers it done and does not retry. Calling the same Once recursively from inside its Do function deadlocks.

## sync.OnceValue
Available since Go 1.21, ~sync.OnceValue~ takes ~func() T~ and returns a concurrent-safe ~func() T~ caching its result, even zero. Each returned closure is independent.
~~~go
loadPort := sync.OnceValue(func() int { return 8081 })
port := loadPort() // Initialization happens here, not at construction
~~~
If the original function panics, every call to the closure panics with the same value without rerunning it. That behavior earns the bonus star.

## Common mistakes
~~~go
if c.value == 0 { c.value = load() } // Wrong: a race, and zero may be valid
var once sync.Once
_ = once // Making this local to every call repeats initialization
~~~

## Mission
Implement ~(c *Calibration) Value(load func() int) int~ using ~sync.Once~. Only the first load executes and every caller sees the same int. Each Calibration is independent and its zero value works. Load is non-nil and does not panic in the required portion.
Implement ~LazyName(load func() string) func() string~ using ~sync.OnceValue~. Do not call load while constructing the closure. Its first consumption caches the result for all goroutines, including an empty string. Bonus: replay the same panic value without rerunning load. Tests run with ~-race~. Leave ~main~ empty.
`, Hints: []string{"Store sync.Once as a field, not a local variable in Value.", "Write the result inside Do and read after it; use a pointer receiver.", "Return the closure produced by sync.OnceValue without calling it in LazyName."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch383Hidden, Bonus: true},
	})
	add(&Level{ID: "38-4-errgroup", Chapter: 38, Kind: "lab", Title: "طاقم محدود", Goal: "شغّل شغل متوازي بحد واضح وألغِ الباقي عند أول خطأ", Setup: seededLabRoom,
		Requires: []string{"sel:errgroup.WithContext"}, Files: []LabFile{{Name: "main.go", Starter: ch384Starter}},
		Race:   true,
		Deps:   []string{"golang.org/x/sync"},
		Bridge: []int{34, 41},
		Lesson: `# طاقم مسح محدود

## قبل ما تكتب
- حد التوازي بيقيّد عدد الدوال الشغالة فعلاً ولا بس عدد النتايج المتجمّعة؟
- لو عامل فشل وعامل تاني مستني، مين يلغيه ومين يستناه؟

## errgroup بدل WaitGroup لوحدها
فصل 10 علّم ~sync.WaitGroup~ للاستنى، وفصل 13 علّم ~context~ للإلغاء. ~golang.org/x/sync/errgroup~ بتجمع الانتظار مع أول error وإلغاء الباقي. الحزمة موجودة في المكتبات المحلية المثبّتة؛ مفيش تنزيل وقت اللعب.
~~~go
g, ctx := errgroup.WithContext(parent)
g.SetLimit(3)
g.Go(func() error { return inspect(ctx, "antenna") })
err := g.Wait()
~~~
~WithContext~ بيرجّع group وcontext مشتقة؛ أول خطأ غير nil بيلغيها. ~Go~ بتبدأ goroutine وقد تحجز المتصل لما الحد مليان. ~SetLimit~ يتنادى قبل الشغل؛ صفر يمنع أي بداية، والسالب معناه من غير حد في الحزمة الحقيقية، لكن عقد مهمتنا بيرفض الاتنين. ~Wait~ تستنى كل العمال وترجّع أول خطأ؛ وبتلغي الـ context المشتقة حتى لو الكل نجح.

## إلغاء تعاوني وحد للتوازي
الـ callback لازم يحترم ~ctx.Done()~؛ الحزمة ماتقدرش تقتل goroutine بالقوة. افحص ~ctx.Err()~ قبل الجدولة وكمان جوه الدالة: الإلغاء ممكن يحصل وإنت محجوز جوه ~Go~. احتفظ بخطأ ~Wait~ بدل ما تستبدله بخطأ الإلغاء. بعد نجاح المجموعة ارجع خطأ الأب لو اتلغى؛ ماترجعش خطأ المشتقة لأن Wait نفسها بتلغيها.
~~~go
select {
case <-ctx.Done(): return ctx.Err()
case reading := <-ready: return save(reading)
}
~~~

## أخطاء شائعة
~~~go
go work(ctx, job) // غلط: تجاهلت الخطأ وحد التوازي والانتظار
return ctx.Err() // لو دي المشتقة بعد Wait، هتبلغ إلغاء حتى بعد نجاح كامل
~~~

## المهمة
نفّذ ~SurveyAll(ctx context.Context, jobs []int, limit int, work func(context.Context,int) error) error~ باستخدام errgroup. limit لازم موجب حتى لو jobs فاضية؛ غير كده error ومفيش work. في النجاح كل عنصر يتنفّذ مرة، حتى القيم المكررة، بحد أقصى limit workers. عند أول خطأ ارجعه من غير فقد هويته لـ ~errors.Is~، الغِ الباقي، وماتبدأش work جديدة بعد ملاحظة الإلغاء. استنى كل worker بدأت قبل الرجوع. الأب الملغي يرجّع خطأه، حتى لو مفيش jobs؛ empty مع أب سليم = nil. work غير nil وبتحترم الإلغاء ومش بتعمل panic. اختبارات ~-race~ وleak check. ~main~ فاضية.
`,
		Hints: []string{"استعمل WithContext وSetLimit قبل حلقة jobs.", "افحص context قبل Go وجوّاها؛ التوازي بيخلق نافذة بين الفحصين.", "استنى Wait حتى بعد الإلغاء؛ خطأ أول worker أهم من context.Canceled."},
		En: &LevelText{Title: "Bounded Crew", Goal: "Bound parallel work and cancel peers on the first error", Files: []LabFile{{Name: "main.go", Starter: ch384Starter}}, Lesson: `# A Bounded Survey Crew

## Before you write
- Does the limit constrain actual running functions or just collected results?
- When one worker fails and another is waiting, who cancels it and waits for it?

## More than a WaitGroup
Chapter 10 introduced ~sync.WaitGroup~ for waiting and chapter 13 introduced ~context~ for cancellation. ~golang.org/x/sync/errgroup~ combines waiting, first-error propagation and cancellation. The package is already pinned in the local dependency catalog; gameplay needs no downloads.
~~~go
g, ctx := errgroup.WithContext(parent)
g.SetLimit(3)
g.Go(func() error { return inspect(ctx, "antenna") })
err := g.Wait()
~~~
~WithContext~ returns a group and a derived context canceled on the first non-nil error. ~Go~ starts a goroutine and may block the caller while the limit is full. Call ~SetLimit~ before starting work. In the real package zero blocks all starts and a negative limit is unbounded; this mission rejects both. ~Wait~ joins all workers and returns the first error. It also cancels the derived context on success.

## Cooperative cancellation and bounded concurrency
The callback must respect ~ctx.Done()~; errgroup cannot forcibly kill a goroutine. Check ~ctx.Err()~ before scheduling and inside the function: cancellation can happen while Go is blocked. Preserve the Wait error instead of replacing it with cancellation. After successful waiting, return the parent's error if canceled; do not return the derived context's error, since Wait itself cancels it.
~~~go
select {
case <-ctx.Done(): return ctx.Err()
case reading := <-ready: return save(reading)
}
~~~

## Common mistakes
~~~go
go work(ctx, job) // Wrong: ignores errors, bounds, and joining
return ctx.Err() // A derived context after Wait reports cancellation even on success
~~~

## Mission
Implement ~SurveyAll(ctx context.Context, jobs []int, limit int, work func(context.Context,int) error) error~ using errgroup. Limit must be positive even for empty jobs; otherwise return an error without calling work. On success process every element exactly once, including duplicates, with at most limit simultaneous callbacks. Return the first error preserving ~errors.Is~ identity, cancel peers, and do not begin new work after observing cancellation. Join all started workers before returning. A canceled parent returns its error, even with no jobs; empty jobs with a live parent succeed. Work is non-nil, cooperates with cancellation and does not panic. Tests use ~-race~ and leak checks. Leave ~main~ empty.
`, Hints: []string{"Use WithContext and SetLimit before the jobs loop.", "Check context before Go and inside its callback to cover cancellation while blocked.", "Always Wait, even after cancellation; preserve the first worker error."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch384Hidden},
	})
	add(&Level{ID: "38-5-allocs", Chapter: 38, Kind: "lab", Title: "ميزانية الذاكرة", Goal: "قلّل تخصيصات النصوص والشرائح وقِس بـ ReportAllocs", Setup: seededLabRoom,
		Requires: []string{"benchmark"}, Files: []LabFile{{Name: "main.go", Starter: ch385Starter}, {Name: "main_test.go", Starter: ch385TestStarter}},
		Bridge: []int{45},
		Lesson: `# ميزانية تخصيصات الرسائل

## قبل ما تكتب
- حجم النص النهائي معروف قبل الكتابة؟
- هل بتقيس جسم الدالة ولا تجهيز بيانات الـ benchmark كمان؟

## السعة قبل الإضافة
~make([]int, 0, n)~ بيعمل slice طولها صفر وسعتها n. ~append~ جوه السعة بيستخدم نفس المصفوفة بدل نمو متكرر. نسخ العناصر لمصفوفة جديدة يمنع نتيجة الدالة من مشاركة ذاكرة المدخل.
~~~go
ports := make([]int, 0, len(all))
for _, p := range all { if p != 0 { ports = append(ports, p) } }
~~~

## strings.Builder
النصوص immutable، فتكرار ~s += piece~ بينسخ بادئة بتكبر كل مرة. ~strings.Builder~ بيجمّع بايتات في مخزن واحد. ~Grow(n)~ يحجز مساحة، ~WriteString~ يضيف نص، ~WriteByte~ يضيف بايت، و~String~ يرجّع النص. ماتنسخش Builder بعد أول استخدام.
~~~go
var b strings.Builder
b.Grow(len(name) + 2)
b.WriteByte('[')
b.WriteString(name)
b.WriteByte(']')
message := b.String()
~~~
~len~ للنص عدد بايتات؛ ده المناسب للحجز حتى مع عربي وemoji. احسب الفواصل كمان.

## قياس التخصيصات وescape analysis
الـ benchmark من فصل 16 بياخد ~*testing.B~. ~b.ReportAllocs()~ يضيف B/op وallocs/op، و~b.ResetTimer()~ يستبعد التحضير. شغّل ~go test -run '^$' -bench . -benchmem~ محلياً. استخدم متغير package يخزّن النتيجة عشان المترجم ميحذفش الشغل.
~~~go
var sink string
func BenchmarkTag(b *testing.B) {
 b.ReportAllocs()
 for i := 0; i < b.N; i++ { sink = tag("relay") }
}
~~~
~go test -gcflags='-m=2'~ بيعرض escape analysis: المترجم بيقرر لو قيمة لازم تعيش بعد stack frame فتتحط على heap. ~return &local~ ممكن يخلّي local تهرب؛ pointer مش ضمان تخصيص، ونتيجة التحليل بتتغير مع inlining ونسخة Go. قلّل التخصيصات بالقياس، مش بتخمين مكان كل متغير. ~testing.AllocsPerRun~ بيقيس متوسط التخصيصات لنداءات متكررة بعد warmup.

## أخطاء شائعة
~~~go
out := make([]int, len(readings))
for _, n := range readings { if n > 0 { out = append(out, n) } } // غلط: أضفت أصفار قبل النتايج
s += next // في حلقة كبيرة: نسخ متكرر للبادئة
~~~

## المهمة
نفّذ ~Assemble(parts []string) string~: وصل كل العناصر بـ ~|~ بالظبط، من غير فاصل زيادة أو حذف عناصر فاضية؛ nil = نص فاضي. نفّذ ~CopyPositive(readings []int) []int~: نسخة مستقلة من العناصر >0 بنفس الترتيب. استعمل preallocation وBuilder.
اكتب ~BenchmarkAssemble~ في main_test.go مع ~b.ReportAllocs~ ونتيجة محفوظة في sink. للمدخل 512 عنصر، Assemble بحد أقصى تخصيصين وCopyPositive بحد أقصى تخصيص واحد. المختبر يقارن تجميع 512 نص مع النسخة البطيئة: حدود النجوم 1x و3x و8x. النسب بتتأثر بالجهاز؛ الصحة والتخصيصات بتتفحص قبل السرعة. ~main~ فاضية.
`,
		Hints: []string{"احسب مجموع len لكل جزء وزوّد عدد الفواصل قبل Grow.", "CopyPositive تبدأ بطول صفر وسعة len(readings)، مش طول len(readings).", "حط ReportAllocs قبل الحلقة، واحفظ النتيجة في package sink."},
		En: &LevelText{Title: "Allocation Budget", Goal: "Reduce string and slice allocations and measure with ReportAllocs", Files: []LabFile{{Name: "main.go", Starter: ch385Starter}, {Name: "main_test.go", Starter: ch385TestStarter}}, Lesson: `# A Message Allocation Budget

## Before you write
- Can you know the final string size before writing it?
- Are you measuring the function or its benchmark setup too?

## Capacity before append
~make([]int, 0, n)~ creates a slice of length zero and capacity n. Appends within that capacity reuse the same array instead of repeatedly growing. Copying into a new array prevents output from sharing the input's storage.
~~~go
ports := make([]int, 0, len(all))
for _, p := range all { if p != 0 { ports = append(ports, p) } }
~~~

## strings.Builder
Strings are immutable, so repeated ~s += piece~ copies an increasingly large prefix. ~strings.Builder~ accumulates bytes in one buffer. ~Grow(n)~ reserves capacity, ~WriteString~ adds text, ~WriteByte~ adds a byte and ~String~ returns the string. Do not copy a Builder after first use.
~~~go
var b strings.Builder
b.Grow(len(name) + 2)
b.WriteByte('[')
b.WriteString(name)
b.WriteByte(']')
message := b.String()
~~~
String ~len~ counts bytes, appropriate for reservation even with Arabic and emoji. Include separators in your size calculation.

## Measuring allocations and escape analysis
Chapter 16 benchmarks receive ~*testing.B~. ~b.ReportAllocs()~ adds B/op and allocs/op; ~b.ResetTimer()~ excludes setup. Run ~go test -run '^$' -bench . -benchmem~ locally. Store results in a package variable so the compiler cannot remove the work.
~~~go
var sink string
func BenchmarkTag(b *testing.B) {
 b.ReportAllocs()
 for i := 0; i < b.N; i++ { sink = tag("relay") }
}
~~~
~go test -gcflags='-m=2'~ reports escape analysis: the compiler decides whether a value must outlive its stack frame and move to the heap. Returning ~&local~ may make it escape; pointers do not guarantee allocations, and inlining and Go versions affect the analysis. Measure rather than guessing placement. ~testing.AllocsPerRun~ measures average allocations over repeated calls after warming up.

## Common mistakes
~~~go
out := make([]int, len(readings))
for _, n := range readings { if n > 0 { out = append(out, n) } } // Wrong: zeros before the results
s += next // In a long loop: repeated prefix copying
~~~

## Mission
Implement ~Assemble(parts []string) string~ joining every element with exactly ~|~, retaining empty elements without an extra separator; nil gives an empty string. Implement ~CopyPositive(readings []int) []int~ returning an independent copy of elements >0 in input order. Use preallocation and Builder.
Write ~BenchmarkAssemble~ in main_test.go with ~b.ReportAllocs~ and a result sink. For 512 elements, Assemble must use at most two allocations and CopyPositive at most one. The lab compares assembly of 512 strings against repeated concatenation: star tiers are 1x, 3x and 8x. Ratios depend on the machine; correctness and allocation budgets run first. Leave ~main~ empty.
`, Hints: []string{"Sum byte lengths plus separator count before Grow.", "Start CopyPositive with zero length and capacity len(readings).", "Call ReportAllocs before the loop and store results in a package sink."}},
		Lab: &Lab{ReqFile: "main_test.go", Hidden: ch385Hidden, Bench: &Bench{Tiers: [3]float64{1, 3, 8}}},
	})
	add(&Level{ID: "38-6-vet", Chapter: 38, Kind: "lab", Title: "تفتيش قبل الإقلاع", Goal: "صلّح نسخ القفل وتنسيق الرسائل وتسريب cancel وشغّل vet", Setup: seededLabRoom,
		Requires: []string{"sel:context.WithCancel"}, Files: []LabFile{{Name: "main.go", Starter: ch386Starter}},
		Race: true,
		Lesson: `# تفتيش قبل الإقلاع

## قبل ما تكتب
- هل نسخ struct فيها Mutex بيحافظ على نفس القفل ونفس العداد؟
- مين مالك cancel ومين مسؤول يناديها في كل مسار رجوع؟

## go vet مش compiler
~go build~ يتأكد إن البرنامج صالح لغوياً؛ ~go vet~ يكتشف أنماط غالباً غلط رغم إنها بتتبني. شغّل ~go vet ./...~. المختبر بيشغّل نفس analyzers عبر ~go test -vet=all -run '^$'~ على ملفاتك قبل الاختبارات المخفية. مرحلة اختبار ملفاتك الأولى هنا بتستخدم ~-vet=off~ عشان تقدر تفرّق البناء عن نتائج التفتيش.
~~~go
func show(n int) { fmt.Printf("count=%s", n) } // printf: نوع غلط للـ %s
~~~
~%d~ للأعداد الصحيحة؛ ~%s~ للنصوص. vet مش إثبات صحة ولا بديل للتستات و~-race~.

## نسخ الأقفال وملكية الإلغاء
~sync.Mutex~ مينفعش يتنسخ بعد الاستخدام. receiver بالقيمة بينسخ struct كلها قبل تنفيذ method؛ استخدم pointer receiver لحماية نفس الحالة المشتركة.
~~~go
type Counter struct { mu sync.Mutex; n int }
func (c Counter) Reset() { c.mu.Lock(); defer c.mu.Unlock(); c.n = 0 } // copylocks
~~~
~context.WithCancel~ بيرجّع context مشتقة وCancelFunc. استدعي cancel حتى لو العملية خلصت بنجاح عشان تحرّر موارد الأب وتبلغ المستمعين. ~defer cancel()~ بعد الإنشاء بيغطي النجاح والخطأ وpanic.
~~~go
ctx, _ := context.WithCancel(parent) // lostcancel
return inspect(ctx)
~~~
قواعد loop variables اتغيّرت من Go 1.22؛ مينفعش نفترض إن vet هيبلغ نفس تحذيرات الإصدارات القديمة. عشان كده المختبر بيركّز على copylocks وprintf وlostcancel اللي لسه موجودين.

## أخطاء شائعة
~~~go
ctx, cancel := context.WithCancel(parent)
err := inspect(ctx)
if err != nil { return err } // غلط: cancel اتنسيت في الخطأ
cancel()
~~~
إخفاء التحذير بـ تغيير النص أو حذف القفل من غير إصلاح السلوك مش حل؛ التستات بتقيس النتيجة المشتركة وصحة الرسالة وإلغاء child.

## المهمة
صلّح starter من غير تغيير signatures غير تحويل Add لـ pointer receiver. ~Gauge~ صفرها شغّال؛ ~Add(n)~ يضيف n و~Total()~ يقرا الإجمالي بأمان تحت التزامن. ~Status(n)~ يرجّع بالظبط ~readings=<n>~ بصيغة عشرية حتى السالب.
~ReadWithCancel(parent, read)~ ينشئ child بـ WithCancel، ينادي read مرة واحدة بالـ child، يرجّع نفس خطأ read ويلغي child على كل مسارات الخروج بعد read، مش قبلها. read غير nil. لازم vet يعدّي و~-race~ ميبلّغش race. ~main~ فاضية.
`,
		Hints: []string{"حوّل Add لـ receiver من نوع *Gauge عشان متنسخش القفل والعداد.", "استخدم %d للأعداد؛ اختبر رقم سالب كمان.", "بعد WithCancel مباشرة حط defer cancel؛ متستبدلش error بتاع read."},
		En: &LevelText{Title: "Preflight Inspection", Goal: "Fix copied locks, formatting and a lost cancel, then pass vet", Files: []LabFile{{Name: "main.go", Starter: ch386Starter}}, Lesson: `# Preflight Inspection

## Before you write
- Does copying a struct with a Mutex preserve the same lock and counter?
- Who owns cancel and calls it on every exit path?

## go vet is not the compiler
~go build~ checks that a program is valid Go; ~go vet~ detects suspicious patterns that still compile. Run ~go vet ./...~. This lab runs the same analyzers through ~go test -vet=all -run '^$'~ on your files before hidden tests. Its initial player-test step uses ~-vet=off~ so compilation and inspection remain distinguishable.
~~~go
func show(n int) { fmt.Printf("count=%s", n) } // printf: wrong type for %s
~~~
~%d~ formats integers; ~%s~ formats strings. Vet is neither a correctness proof nor a replacement for tests and ~-race~.

## Copied locks and cancellation ownership
A ~sync.Mutex~ must not be copied after use. A value receiver copies the entire struct before the method runs; use a pointer receiver to protect the same shared state.
~~~go
type Counter struct { mu sync.Mutex; n int }
func (c Counter) Reset() { c.mu.Lock(); defer c.mu.Unlock(); c.n = 0 } // copylocks
~~~
~context.WithCancel~ returns a derived context and CancelFunc. Call cancel even on success to release parent resources and notify listeners. Deferring cancellation immediately after construction covers success, error and panic exits.
~~~go
ctx, _ := context.WithCancel(parent) // lostcancel
return inspect(ctx)
~~~
Loop variable rules changed in Go 1.22; do not assume vet still reports historical warnings. This lab focuses on copylocks, printf and lostcancel, which remain applicable.

## Common mistakes
~~~go
ctx, cancel := context.WithCancel(parent)
err := inspect(ctx)
if err != nil { return err } // Wrong: forgets cancellation on error
cancel()
~~~
Silencing diagnostics by changing text or deleting a lock without fixing behavior is insufficient. Hidden checks exercise shared state, message content and child cancellation.

## Mission
Repair the starter without changing signatures except making Add a pointer receiver. The zero Gauge works; ~Add(n)~ adds n and ~Total()~ reads the total safely under concurrency. ~Status(n)~ returns exactly ~readings=<n>~ in decimal, including negative numbers.
~ReadWithCancel(parent, read)~ creates a child using WithCancel, calls non-nil read once with it, returns read's original error, and cancels the child on every exit after read, never before it. Vet must pass and ~-race~ must report no races. Leave ~main~ empty.
`, Hints: []string{"Make Add use a *Gauge receiver to avoid copying the lock and counter.", "Use %d for integers and check a negative number too.", "Defer cancel immediately after WithCancel and preserve the error returned by read."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch386Hidden, Vet: true},
	})
	add(&Level{ID: "38-7-pprof", Chapter: 38, Kind: "lab", Title: "فين الوقت بيروح؟", Goal: "اكتب CPU وheap profiles وصلّح جمع القراءات المميّزة", Setup: seededLabRoom,
		Requires: []string{"sel:pprof.StartCPUProfile", "sel:pprof.WriteHeapProfile"}, Files: []LabFile{{Name: "main.go", Starter: ch387Starter}},
		Lesson: `# فين وقت المسح بيروح؟

## قبل ما تكتب
- هل bottleneck من الخوارزمية ولا من تكلفة كتابة ملف القياس؟
- لو profiler شغّال بالفعل، هل تبدأ الشغل ولا ترجع الخطأ؟

## CPU وheap profiles
~runtime/pprof~ بتكتب profiles ثنائية مضغوطة. CPU profile عينات من الـ stack أثناء التنفيذ؛ الـ heap profile لقطة تخصيصات ذاكرة بأخذ عينات، مش قائمة بكل object. محتاج workload طويلة عشان CPU samples تبقى مفيدة؛ الشغل الصغير ممكن ينتج ملف صالح من غير عينات CPU، وده طبيعي.
~~~go
if err := pprof.StartCPUProfile(cpuWriter); err != nil { return err }
scanAntenna()
pprof.StopCPUProfile()
err := pprof.WriteHeapProfile(heapWriter)
~~~
~StartCPUProfile~ بتبدأ profiler واحدة على مستوى العملية وبتاخد ~io.Writer~، الواجهة الصغيرة من فصل 9 اللي فيها ~Write([]byte) (int,error)~. ~StopCPUProfile~ بتوقف التجميع وتستنى اكتمال الكتابة. ~WriteHeapProfile~ بتكتب لقطة heap وترجّع خطأ الكاتب. المتصل مالك الكتاب؛ متقفلهمش. بعد فشل Start متوقفش profiler تخص حد تاني.

## قراءة التقرير محلياً
في برنامج محلي مرّر ملفات من ~os.Create~ كـ writers، افحص الخطأ واقفل الملفات بعد اكتمال الكتابة. شغّل ~go tool pprof -top cpu.prof~ أو ~go tool pprof -top -alloc_space heap.prof~؛ ~-alloc_space~ يعرض الحجم المتراكم للتخصيصات، والافتراضي heap بيركّز على الذاكرة الحية. للاطلاع على السطور استخدم ~-list=UniqueTotal~ مع executable الموافق. ممكن تجمع profile من benchmark كمان بـ ~go test -run '^$' -bench . -cpuprofile cpu.prof -memprofile heap.prof~. اللعبة مش بتفتح UI للأداة.

## إصلاح hotspot بالقياس
فحص كل قراءة مقابل كل اللي قبلها تكلفته O(n²). set من مهمة 38-2 يحفظ العضوية ويحوّل المرور إلى O(n) في المتوسط مقابل ذاكرة O(n). جمّع كل قيمة مختلفة مرة واحدة؛ التكرارات مش تزود المجموع. الرقم السالب قراءة صحيحة.
~~~go
for i, v := range samples {
 for j := 0; j < i; j++ { _ = samples[j] == v } // المقارنات بتكبر تربيعياً
}
~~~
قِس الدالة الساخنة لوحدها بbenchmark؛ ماتحطش تشغيل profiler وكتابة الملفات جوه القياس اللي بتقارن بيه الخوارزميات.

## أخطاء شائعة
~~~go
pprof.StartCPUProfile(w) // غلط: تجاهلت إن profiler ممكن تكون مشغولة
return calculate()     // غلط: مفيش Stop، والملف ممكن ميتكتبش كامل
~~~
CPU writer errors مش بترجع من StopCPUProfile لأنها مابترجّعش error؛ ماتعدش بحاجة الـ API مش بتدعمها. المهمة بتفحص خطأ heap writer وفشل البداية.

## المهمة
نفّذ ~UniqueTotal(readings []int) int~: مجموع القيم المختلفة بالظبط، بما فيها السالب والصفر؛ nil = صفر والمدخل ثابت. افترض المجموع داخل حدود int.
نفّذ ~ProfileSurvey(cpu, heap io.Writer, readings []int) (int,error)~: ابدأ CPU profile قبل UniqueTotal، اوقفها واستنى الكتابة، وبعدين اكتب heap profile. ارجع المجموع وnil في النجاح. لو Start فشلت ارجع 0 والخطأ فوراً؛ لو heap writer فشل ارجع المجموع والخطأ مع الحفاظ على ~errors.Is~. الكتاب غير nil؛ CPU writer سليم. الناتجين لازم يبقوا profiles حقيقية غير فاضية، حتى مع input فاضي، ويتاح بدء profile تانية بعد الرجوع.
بعد الصحة، Lab.Bench تقارن 4096 قراءة مميّزة ضد الحل التربيعي، من غير profiling داخل القياس: حدود النجوم 2x و5x و10x. ~main~ فاضية.
`,
		Hints: []string{"استعمل set لتجمع كل رقم مرة واحدة من غير حلقات متداخلة.", "افحص StartCPUProfile؛ لو نجحت لازم StopCPUProfile قبل الرجوع أو كتابة heap.", "رجّع خطأ WriteHeapProfile مع المجموع؛ وقِس UniqueTotal لوحدها."},
		En: &LevelText{Title: "Where Does the Time Go?", Goal: "Write CPU and heap profiles and optimize distinct-reading totals", Files: []LabFile{{Name: "main.go", Starter: ch387Starter}}, Lesson: `# Where Does Survey Time Go?

## Before you write
- Is the bottleneck the algorithm or writing the measurement file?
- If another CPU profiler is already running, should you start work or return its error?

## CPU and heap profiles
~runtime/pprof~ writes compressed binary profiles. A CPU profile samples execution stacks; a heap profile samples memory allocations rather than listing every object. Useful CPU samples need a sustained workload. A tiny workload can produce a valid profile with no CPU samples.
~~~go
if err := pprof.StartCPUProfile(cpuWriter); err != nil { return err }
scanAntenna()
pprof.StopCPUProfile()
err := pprof.WriteHeapProfile(heapWriter)
~~~
~StartCPUProfile~ starts a process-wide profiler and takes ~io.Writer~, the small interface from chapter 9 with ~Write([]byte) (int,error)~. ~StopCPUProfile~ stops sampling and waits for writing to finish. ~WriteHeapProfile~ writes a heap snapshot and returns writer errors. Writers belong to the caller; do not close them. If Start fails, do not stop someone else's profiler.

## Inspecting a report locally
In a local program pass files created by ~os.Create~ as writers, check errors and close files after writing finishes. Use ~go tool pprof -top cpu.prof~ or ~go tool pprof -top -alloc_space heap.prof~. Alloc_space shows cumulative allocated bytes; heap defaults focus on live memory. Use ~-list=UniqueTotal~ with the matching executable to inspect lines. Benchmarks can also collect profiles through ~go test -run '^$' -bench . -cpuprofile cpu.prof -memprofile heap.prof~. The game does not open the tool's UI.

## Fixing the hotspot
Comparing each reading against all earlier readings costs O(n²). The set from mission 38-2 makes traversal O(n) on average, using O(n) memory. Add each distinct number once. Negative readings are valid.
~~~go
for i, v := range samples {
 for j := 0; j < i; j++ { _ = samples[j] == v } // Quadratically growing comparisons
}
~~~
Benchmark the hot function alone. Keep profiler startup and file writing outside the algorithm comparison.

## Common mistakes
~~~go
pprof.StartCPUProfile(w) // Wrong: ignores an already-busy profiler
return calculate()     // Wrong: no Stop, so output may remain incomplete
~~~
StopCPUProfile returns no error, so CPU writer failures cannot be promised through that API. This mission checks heap writer failures and startup failures.

## Mission
Implement ~UniqueTotal(readings []int) int~, summing distinct values exactly once including negatives and zero. Nil gives zero; do not change inputs. Assume the total fits int.
Implement ~ProfileSurvey(cpu, heap io.Writer, readings []int) (int,error)~: start CPU profiling before UniqueTotal, stop and finish writing, then write a heap profile. Return the total and nil on success. If Start fails return zero and that error immediately. If heap writing fails return the total and preserve its error for ~errors.Is~. Writers are non-nil and the CPU writer succeeds. Both outputs must be real non-empty profiles, even for empty input, and a new CPU profile must be possible after return.
After correctness checks, Lab.Bench compares 4096 distinct readings against a quadratic implementation, with no profiling inside the measurement. Star tiers are 2x, 5x and 10x. Leave ~main~ empty.
`, Hints: []string{"Use a set to add each distinct number once without nested scans.", "Check StartCPUProfile; after success call StopCPUProfile before returning or writing the heap.", "Return the heap writer error with the total; benchmark UniqueTotal alone."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch387Hidden, Bench: &Bench{Tiers: [3]float64{2, 5, 10}}},
	})
}

const ch381Starter = `package main

import "iter"

func Readings(values []int, visit func(int)) iter.Seq[int] {
	return func(yield func(int) bool) {}
}
func Indexed(values []int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {}
}
func First(values []int, n int) []int { return nil }
func main()                           {}
`

const ch381Hidden = `package main

import (
	"math/rand"
	"orbital/colony/labkit"
	"reflect"
	"testing"
)

func TestHiddenIteratorStops(t *testing.T) {
	labkit.LeakCheck(t)
	for seed := int64(381); seed < 386; seed++ {
		rng := rand.New(rand.NewSource(seed))
		values := make([]int, 10+rng.Intn(80))
		for i := range values {
			values[i] = rng.Intn(1000) - 500
		}
		for _, n := range []int{1, len(values) / 2, len(values)} {
			calls := 0
			var got []int
			for v := range Readings(values, func(int) { calls++ }) {
				got = append(got, v)
				if len(got) == n {
					break
				}
			}
			if calls != n || !reflect.DeepEqual(got, values[:n]) {
				t.Fatalf("stop after %d: visits=%d values=%v", n, calls, got)
			}
			if first := First(values, n); !reflect.DeepEqual(first, values[:n]) {
				t.Fatal("First", n, first)
			}
		}
		count := 0
		for i, v := range Indexed(values) {
			if i != count || v != values[i] {
				t.Fatal("index/value order")
			}
			count++
		}
		if count != len(values) {
			t.Fatal("missing indexed readings")
		}
		calls := 0
		Indexed(values)(func(i, v int) bool { calls++; return false })
		if calls != 1 {
			t.Fatal("Seq2 ignored stop")
		}
		if got := First(values, len(values)+5); !reflect.DeepEqual(got, values) {
			t.Fatal("limit beyond end")
		}
	}
	for _, n := range []int{-3, 0, 1} {
		if len(First(nil, n)) != 0 {
			t.Fatal("empty")
		}
	}
	calls := 0
	for range Readings(nil, func(int) { calls++ }) {
		t.Fatal("empty yield")
	}
	for range Indexed(nil) {
		t.Fatal("empty indexed yield")
	}
	if calls != 0 {
		t.Fatal("empty visit")
	}
	if len(First([]int{1}, 0)) != 0 || len(First([]int{1}, -1)) != 0 {
		t.Fatal("nonpositive limit")
	}
}
func TestHiddenBonusReusableIterator(t *testing.T) {
	values := []int{7, 9, 11}
	calls := 0
	seq := Readings(values, func(int) { calls++ })
	for round := 0; round < 3; round++ {
		count := 0
		for v := range seq {
			if v != values[count] {
				t.Fatal("restart")
			}
			count++
			if count == 2 {
				break
			}
		}
		if count != 2 {
			t.Fatal("iterator not reusable")
		}
	}
	if calls != 6 {
		t.Fatal("producer ran ahead")
	}
}
`

const ch382Starter = `package main

func SelectGates[T comparable](rows [][]T, blocked []T) []T { return nil }
func main()                                                 {}
`

const ch382Hidden = `package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"reflect"
	"testing"
)

func TestHiddenGates(t *testing.T) {
	for seed := int64(382); seed < 387; seed++ {
		rng := rand.New(rand.NewSource(seed))
		blocked := []int{rng.Intn(5), 10}
		rows := [][]int{nil, {blocked[0], 99}, {7, blocked[0]}, {7, 8}, {7, 8, 9}}
		for j := 0; j < 80; j++ {
			row := make([]int, rng.Intn(8))
			for i := range row {
				row[i] = rng.Intn(15)
			}
			rows = append(rows, row)
		}
		before := make([][]int, len(rows))
		for i, row := range rows {
			if row != nil {
				before[i] = append([]int{}, row...)
			}
		}
		var want []int
		for _, row := range rows {
			for _, v := range row {
				denied := false
				for _, b := range blocked {
					if v == b {
						denied = true
					}
				}
				if denied {
					break
				}
				used := false
				for _, s := range want {
					if v == s {
						used = true
					}
				}
				if !used {
					want = append(want, v)
					break
				}
			}
		}
		got := SelectGates(rows, blocked)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("seed %d: %v want %v", seed, got, want)
		}
		if !reflect.DeepEqual(rows, before) {
			t.Fatal("input changed")
		}
	}
	if len(SelectGates[int](nil, nil)) != 0 {
		t.Fatal("empty")
	}
	if got := SelectGates([][]string{{"مريخ", "أرض"}, {"مريخ", "قمر"}, {"حظر", "شمس"}, {"قمر", "شمس"}}, []string{"حظر"}); !reflect.DeepEqual(got, []string{"مريخ", "قمر", "شمس"}) {
		t.Fatal("generic strings", got)
	}
}
func TestHiddenSetAndLabels(t *testing.T) {
	src, err := os.ReadFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "main.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	hasSet, hasBreak, hasContinue := false, false, false
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "SelectGates" {
			continue
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.MapType:
				if empty, ok := n.Value.(*ast.StructType); ok && len(empty.Fields.List) == 0 {
					hasSet = true
				}
			case *ast.BranchStmt:
				if n.Label != nil {
					if n.Tok == token.BREAK {
						hasBreak = true
					}
					if n.Tok == token.CONTINUE {
						hasContinue = true
					}
				}
			}
			return true
		})
	}
	if !hasSet || !hasBreak || !hasContinue {
		t.Fatal("use map[T]struct{} and labeled break/continue inside SelectGates")
	}
}
`

const ch383Starter = `package main

import "sync"

type Calibration struct {
	once  sync.Once
	value int
}

func (c *Calibration) Value(load func() int) int { return 0 }
func LazyName(load func() string) func() string  { return func() string { return "" } }
func main()                                      {}
`

const ch383Hidden = `package main

import (
	"math/rand"
	"sync"
	"testing"
)

func TestHiddenOnceConcurrent(t *testing.T) {
	rng := rand.New(rand.NewSource(383))
	for round := 0; round < 12; round++ {
		want := rng.Intn(100000) - 50000
		var calls int
		var c Calibration
		var wg sync.WaitGroup
		start := make(chan struct{})
		for i := 0; i < 40; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				if got := c.Value(func() int { calls++; return want }); got != want {
					t.Errorf("value %d want %d", got, want)
				}
			}()
		}
		close(start)
		wg.Wait()
		if calls != 1 {
			t.Fatal("load count", calls)
		}
		if c.Value(func() int { t.Error("initialized again"); return 0 }) != want {
			t.Fatal("cache lost")
		}
	}
}
func TestHiddenLazyNameConcurrent(t *testing.T) {
	for _, want := range []string{"", "قمر", "Kepler-38"} {
		calls := 0
		lazy := LazyName(func() string { calls++; return want })
		if calls != 0 {
			t.Fatal("eager initialization")
		}
		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if lazy() != want {
					t.Error("name changed")
				}
			}()
		}
		wg.Wait()
		if calls != 1 {
			t.Fatal("name load count", calls)
		}
	}
	a, b := LazyName(func() string { return "a" }), LazyName(func() string { return "b" })
	if a() != "a" || b() != "b" {
		t.Fatal("independent closures")
	}
}
func TestHiddenBonusOnceValuePanic(t *testing.T) {
	marker := &struct{ N int }{38}
	calls := 0
	lazy := LazyName(func() string { calls++; panic(marker) })
	for i := 0; i < 3; i++ {
		func() {
			defer func() {
				if recover() != marker {
					t.Error("panic not replayed")
				}
			}()
			lazy()
		}()
	}
	if calls != 1 {
		t.Fatal("panic retried")
	}
}
`

const ch384Starter = `package main

import "context"

func SurveyAll(ctx context.Context, jobs []int, limit int, work func(context.Context, int) error) error {
	return nil
}
func main() {}
`

const ch384Hidden = `package main

import (
	"context"
	"errors"
	"math/rand"
	"orbital/colony/labkit"
	"sync"
	"testing"
	"time"
)

func TestHiddenBoundedWork(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(384))
	for _, limit := range []int{1, 2, 5} {
		jobs := make([]int, 30)
		want := map[int]int{}
		for i := range jobs {
			jobs[i] = rng.Intn(1000)
			want[jobs[i]]++
		}
		var mu sync.Mutex
		active, peak := 0, 0
		seen := map[int]int{}
		err := SurveyAll(context.Background(), jobs, limit, func(ctx context.Context, n int) error {
			mu.Lock()
			active++
			if active > peak {
				peak = active
			}
			seen[n]++
			mu.Unlock()
			time.Sleep(time.Millisecond)
			mu.Lock()
			active--
			mu.Unlock()
			return nil
		})
		if err != nil || active != 0 || peak > limit {
			t.Fatalf("err=%v active=%d peak=%d", err, active, peak)
		}
		for n, count := range want {
			if seen[n] != count {
				t.Fatal("missing/duplicate job", n)
			}
		}
	}
}
func TestHiddenFirstErrorCancelsPeers(t *testing.T) {
	labkit.LeakCheck(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	boom := errors.New("sensor offline")
	peerStarted := make(chan struct{})
	peerStopped := make(chan struct{})
	calls := 0
	var mu sync.Mutex
	err := SurveyAll(ctx, []int{0, 1, 2, 3, 4}, 2, func(c context.Context, n int) error {
		mu.Lock()
		calls++
		mu.Unlock()
		if n == 0 {
			select {
			case <-peerStarted:
				return boom
			case <-c.Done():
				return c.Err()
			}
		}
		if n == 1 {
			close(peerStarted)
			<-c.Done()
			close(peerStopped)
			return c.Err()
		}
		return nil
	})
	if !errors.Is(err, boom) {
		t.Fatalf("first error lost: %v", err)
	}
	select {
	case <-peerStopped:
	default:
		t.Fatal("returned before peer stopped")
	}
	if calls != 2 {
		t.Fatal("work started after cancellation", calls)
	}
}
func TestHiddenCanceledAndEmptyCrew(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if !errors.Is(SurveyAll(ctx, []int{1}, 1, func(context.Context, int) error { t.Error("called after cancellation"); return nil }), context.Canceled) {
		t.Fatal("parent cancellation")
	}
	if !errors.Is(SurveyAll(ctx, nil, 1, func(context.Context, int) error { return nil }), context.Canceled) {
		t.Fatal("empty cancellation")
	}
	if SurveyAll(context.Background(), nil, 1, func(context.Context, int) error { t.Error("empty work"); return nil }) != nil {
		t.Fatal("empty")
	}
	for _, limit := range []int{0, -1} {
		if SurveyAll(context.Background(), nil, limit, func(context.Context, int) error { return nil }) == nil {
			t.Fatal("invalid limit")
		}
	}
}
`

const ch385Starter = `package main

func Assemble(parts []string) string    { return "" }
func CopyPositive(readings []int) []int { return nil }
func main()                             {}
`

const ch385TestStarter = `package main

import "testing"

func BenchmarkAssemble(b *testing.B) {
	parts := []string{"orbit", "relay", "dock"}
	for i := 0; i < b.N; i++ {
		_ = Assemble(parts)
	}
}
`

const ch385Hidden = `package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
)

var hiddenAssembly string
var hiddenPositive []int

func TestHiddenAssemblyAndCopy(t *testing.T) {
	rng := rand.New(rand.NewSource(385))
	for round := 0; round < 80; round++ {
		parts := []string{"", "قمر", "🚀"}
		nums := []int{0, -1, 1}
		var want []int
		for i := 0; i < rng.Intn(80); i++ {
			parts = append(parts, fmt.Sprint(rng.Int()))
			nums = append(nums, rng.Intn(100)-50)
		}
		if got := Assemble(parts); got != strings.Join(parts, "|") {
			t.Fatal("assembly differs", got)
		}
		before := append([]int{}, nums...)
		for _, n := range nums {
			if n > 0 {
				want = append(want, n)
			}
		}
		got := CopyPositive(nums)
		if !reflect.DeepEqual(got, want) {
			t.Fatal("positive readings", got, want)
		}
		if len(got) > 0 {
			got[0] = -999
		}
		if !reflect.DeepEqual(nums, before) {
			t.Fatal("copy aliases input")
		}
	}
	if Assemble(nil) != "" || Assemble([]string{""}) != "" || Assemble([]string{"", ""}) != "|" {
		t.Fatal("empty strings")
	}
	if len(CopyPositive(nil)) != 0 || len(CopyPositive([]int{-1, 0})) != 0 {
		t.Fatal("empty positives")
	}
}
func TestHiddenAllocationBudget(t *testing.T) {
	parts := make([]string, 512)
	for i := range parts {
		parts[i] = "station"
	}
	nums := make([]int, 512)
	for i := range nums {
		nums[i] = i + 1
	}
	if n := testing.AllocsPerRun(20, func() { hiddenAssembly = Assemble(parts) }); n > 2 {
		t.Fatal("assembly allocation budget <=2, got", n)
	}
	if n := testing.AllocsPerRun(20, func() { hiddenPositive = CopyPositive(nums) }); n > 1 {
		t.Fatal("copy allocation budget <=1, got", n)
	}
}
func TestHiddenBenchmarkReportsAllocations(t *testing.T) {
	src, err := os.ReadFile("main_test.go")
	if err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "main_test.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	reports := false
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "BenchmarkAssemble" {
			continue
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			if call, ok := node.(*ast.CallExpr); ok {
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "ReportAllocs" {
					reports = true
				}
			}
			return true
		})
	}
	if !reports {
		t.Fatal("BenchmarkAssemble must call b.ReportAllocs")
	}
}
func naiveAssemble(parts []string) string {
	s := ""
	for i, p := range parts {
		if i > 0 {
			s += "|"
		}
		s += p
	}
	return s
}
func BenchmarkHiddenNaive(b *testing.B) {
	parts := make([]string, 512)
	for i := range parts {
		parts[i] = "station"
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hiddenAssembly = naiveAssemble(parts)
	}
}
func BenchmarkHiddenTheirs(b *testing.B) {
	parts := make([]string, 512)
	for i := range parts {
		parts[i] = "station"
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hiddenAssembly = Assemble(parts)
	}
}
`

const ch386Starter = `package main

import (
	"context"
	"fmt"
	"sync"
)

type Gauge struct {
	mu    sync.Mutex
	total int
}

func (g Gauge) Add(n int)   { g.mu.Lock(); defer g.mu.Unlock(); g.total += n }
func (g *Gauge) Total() int { g.mu.Lock(); defer g.mu.Unlock(); return g.total }
func Status(n int) string   { return fmt.Sprintf("readings=%s", n) }
func ReadWithCancel(parent context.Context, read func(context.Context) error) error {
	ctx, _ := context.WithCancel(parent)
	return read(ctx)
}
func main() {}
`

const ch386Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

func TestHiddenGaugeAndStatus(t *testing.T) {
	rng := rand.New(rand.NewSource(386))
	for round := 0; round < 10; round++ {
		var g Gauge
		var wg sync.WaitGroup
		want := 0
		for i := 0; i < 40; i++ {
			n := rng.Intn(200) - 100
			want += n
			wg.Add(1)
			go func() { defer wg.Done(); g.Add(n) }()
		}
		wg.Wait()
		if g.Total() != want {
			t.Fatal("copied lock/value", g.Total(), want)
		}
		if got := Status(want); got != fmt.Sprintf("readings=%d", want) {
			t.Fatal("status", got)
		}
	}
}
func TestHiddenCancelOnEveryExit(t *testing.T) {
	boom := errors.New("read failed")
	for _, readErr := range []error{nil, boom} {
		var child context.Context
		err := ReadWithCancel(context.Background(), func(c context.Context) error {
			child = c
			if c.Err() != nil {
				t.Error("premature cancel")
			}
			return readErr
		})
		if !errors.Is(err, readErr) {
			t.Fatal("error lost", err)
		}
		if child == nil {
			t.Fatal("callback missing")
		}
		select {
		case <-child.Done():
		default:
			t.Fatal("cancel lost")
		}
	}
	parent, cancel := context.WithCancel(context.Background())
	cancel()
	if err := ReadWithCancel(parent, func(c context.Context) error { return c.Err() }); !errors.Is(err, context.Canceled) {
		t.Fatal("parent cancellation", err)
	}
}
`

const ch387Starter = `package main

import "io"

func UniqueTotal(readings []int) int                                 { return 0 }
func ProfileSurvey(cpu, heap io.Writer, readings []int) (int, error) { return 0, nil }
func main()                                                          {}
`

const ch387Hidden = `package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"math/rand"
	"runtime/pprof"
	"testing"
)

var hiddenProfileTotal int

func TestHiddenUniqueTotal(t *testing.T) {
	rng := rand.New(rand.NewSource(387))
	for round := 0; round < 80; round++ {
		nums := make([]int, rng.Intn(300))
		seen := map[int]bool{}
		want := 0
		for i := range nums {
			nums[i] = rng.Intn(200) - 100
			if !seen[nums[i]] {
				seen[nums[i]] = true
				want += nums[i]
			}
		}
		before := append([]int{}, nums...)
		if got := UniqueTotal(nums); got != want {
			t.Fatal("unique total", got, want)
		}
		for i := range nums {
			if nums[i] != before[i] {
				t.Fatal("input changed")
			}
		}
	}
	if UniqueTotal(nil) != 0 || UniqueTotal([]int{5, 5, 5}) != 5 {
		t.Fatal("empty/duplicates")
	}
}
func checkProfile(t *testing.T, blob []byte) {
	t.Helper()
	r, err := gzip.NewReader(bytes.NewReader(blob))
	if err != nil {
		t.Fatal("profile is not gzip", err)
	}
	defer r.Close()
	decoded, err := io.ReadAll(r)
	if err != nil || len(decoded) == 0 {
		t.Fatal("empty/corrupt profile", err)
	}
}
func TestHiddenProfilesAndCleanup(t *testing.T) {
	for _, nums := range [][]int{nil, {5, 5, -2, 9}, {17, 17, 17}} {
		var cpu, heap bytes.Buffer
		got, err := ProfileSurvey(&cpu, &heap, nums)
		if err != nil || got != UniqueTotal(nums) {
			t.Fatal("profile result", got, err)
		}
		checkProfile(t, cpu.Bytes())
		checkProfile(t, heap.Bytes())
	}
}

type brokenProfileWriter struct{ err error }

func (w brokenProfileWriter) Write(p []byte) (int, error) { return 0, w.err }
func TestHiddenProfileFailures(t *testing.T) {
	boom := errors.New("heap disk full")
	var cpu bytes.Buffer
	got, err := ProfileSurvey(&cpu, brokenProfileWriter{boom}, []int{8, 8, 3})
	if got != 11 || !errors.Is(err, boom) {
		t.Fatal("heap write error lost", got, err)
	}
	if err := pprof.StartCPUProfile(io.Discard); err != nil {
		t.Fatal("CPU profiler not released", err)
	}
	defer pprof.StopCPUProfile()
	if _, err := ProfileSurvey(io.Discard, io.Discard, nil); err == nil {
		t.Fatal("busy CPU profiler must fail")
	}
}
func naiveUniqueTotal(nums []int) int {
	total := 0
	for i, v := range nums {
		found := false
		for j := 0; j < i; j++ {
			if nums[j] == v {
				found = true
				break
			}
		}
		if !found {
			total += v
		}
	}
	return total
}
func BenchmarkHiddenNaive(b *testing.B) {
	nums := make([]int, 4096)
	for i := range nums {
		nums[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hiddenProfileTotal = naiveUniqueTotal(nums)
	}
}
func BenchmarkHiddenTheirs(b *testing.B) {
	nums := make([]int, 4096)
	for i := range nums {
		nums[i] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hiddenProfileTotal = UniqueTotal(nums)
	}
}
`
