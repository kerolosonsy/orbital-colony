# Chapter 12 lesson storyboards

Content contract, 2026-11-18. Revision 1, branch `lessons-ch12`. Each mission uses the six approved steps and markdown, question, and walkthrough blocks only. Word counts include titles, copy, question options and feedback, frame captions, and hints; code is excluded. Examples use illustrative values, not mission data.

## 12-1 — العدّاد المشترك (`12-1-shared-total`)

**Words:** Arabic 297; English 292.

**Objective:** explain why a shared increment needs one mutex. Reinforce that a WaitGroup waits for completion but does not protect shared memory. The starter already harvests.

### Mission / المهمة

**AR:** 1. خلّي 4 درونات تزود نفس المتغير.  
2. احمِ التحديث من data race.  
3. ابعت الإجمالي بعد انتهاء الشغل.

الـstarter بيحصد بالفعل؛ ركّز على الكتابة المشتركة.

**EN:** 1. Make 4 drones add to one total.  
2. Protect updates from a data race.  
3. Report the result.

The starter already harvests.

### Idea / الفكرة

**AR:** تخيّل ~total~ لوحة واحدة. كل تحديث يقرأها، يزوّدها، ثم يكتبها.

لو التحديثات اتداخلت، ممكن قيمة تضيع؛ النتيجة الصح بالصدفة مش ضمان.

Mutex مشترك يسمح لواحد بس يدخل الجزء المحمي.  
~WaitGroup~ يستنى الشغل، لكنه مش بيحمي المتغير.

**EN:** Think of ~total~ as one board. An update reads, adds, then writes.

Overlapping updates can lose a value. A correct result by luck proves nothing.

One shared mutex protects this short update. ~WaitGroup~ waits; it does not protect total.

### Try it / جرّب

**Predict:** «المثال يضيف 2 و3 تحت قفل واحد. إيه اللي يطبع؟» / “What does this locked example print?”

| Options | Correct | One-sentence feedback |
|---|---:|---|
| 5 · 2 · النتيجة ممكن تتغير / 5 · 2 · It may vary | 1 | التحديثين بيتجمعوا؛ ترتيبهم مش بيغيّر المجموع. / Both additions remain. |
|  |  | ده أول تحديث لوحده؛ العامل التاني كمان بيضيف. / That is only the first. |
|  |  | الترتيب مش مضمون، لكن كل تحديث محمي بالقفل. / Order may vary; the lock protects both. |

**AR walkthrough text:** ده ترتيب واحد ممكن: 2 يدخل الأول. Go مش بتضمن مين يمسك القفل أولاً.  
**EN walkthrough text:** One possible order: 2 enters first. Go does not promise who gets the lock first.

<!-- go-example id=e121-counter imports=fmt,sync outcome=ok -->
```go
var mu sync.Mutex
total := 0
var wg sync.WaitGroup
add := func(n int) {
    defer wg.Done()
    mu.Lock()
    defer mu.Unlock()
    total += n
}
wg.Add(2)
go add(2)
go add(3)
wg.Wait()
fmt.Println(total)
```

| Line | AR caption | EN caption | State / output |
|---:|---|---|---|
| 0 | قبل التشغيل. | Before execution. | — |
| 2 | الإجمالي يبدأ صفر. | Total starts at zero. | total=0 |
| 6 | ترتيب ممكن: n=2 يمسك القفل. | One order: n=2 gets the lock. | total=0, n=2 |
| 7 | defer هيفتحه عند الخروج. | defer unlocks on return. | total=0, n=2 |
| 8 | الإجمالي صار 2. | Total becomes 2. | total=2, n=2 |
| 9 | الرجوع يفتح القفل. | Return releases the lock. | total=2 |
| 6 | بعده n=3 يدخل. | Then n=3 enters. | total=2, n=3 |
| 8 | الإجمالي صار 5. | Total becomes 5. | total=5 |
| 13 | العمال خلصوا. | Workers have finished. | total=5 |
| 14 | يطبع 5. | Prints 5. | stdout: `5` |

### Watch out / خد بالك

**AR:** القفل لازم يبقى هو نفسه عند كل العاملين.

خلي الجزء بين ~Lock~ و~Unlock~ صغير؛ حرّره حتى لو الدالة خرجت.

Race detector ممكن يكشف السباق حتى لو الإجمالي ظهر صح.

**EN:** All workers must use the same mutex.

Keep the locked section small; defer its release.

The race detector can catch a race despite a correct total.

### Think & check / فكّر واختبر

**Plan:** شغّل العمال → احمِ التحديث → استناهم → ابعت الإجمالي. / Start workers → protect the update → wait → report.

**Q1:** هل ~wg.Wait()~ بيحمي total أثناء تحديثه؟ / Does Wait protect total during updates?  
Options: أيوه / Yes · لأ / No · بس لو عامل واحد / Only with one worker. **Answer: 2 / No.**  
Feedback: الانتظار مش قفل؛ ممكن العمال يكتبوا في نفس الوقت. / Wait only waits; it does not lock.

**Q2:** أنهي اختيار يحمي كل تحديثات الإجمالي؟ / What protects every total update?  
Options: نفس القفل المشترك / One shared mutex · قفل جديد لكل goroutine / One lock per goroutine · القفل بعد Wait / Lock after Wait. **Answer: 1.**  
Feedback: كل كاتب لازم يتفق على نفس القفل. / All writers need the same lock.

### Your turn / دورك

راجع سطر الجمع في الـstarter، ثم شغّل واختبر race detector.  
Check the starter's addition line, then run it with the race detector.

### Hints

1. **Concept / الفكرة:** نتيجة سليمة مرة مش إثبات إن الكتابة المشتركة آمنة. / A lucky correct result does not prove shared writes are safe.
2. **Structure / الترتيب:** افصل تحديث الإجمالي عن الحركة والحصاد وانتظار العمال. / Separate total updates from drone work.
3. **Specific / السطر:** راجع ~total += n~ جوه الـgoroutine. / Inspect ~total += n~ inside the goroutine.

#### Verified e121-counter

Real stdout: 5

## 12-2 — السجل الآمن (`12-2-registry`)

**Words:** Arabic 287; English 283.

**Objective:** protect concurrent map writes with one mutex owned beside the map. Apply the known pointer-receiver pattern and read after the workers finish. The starter already measures and records locations.

### Mission / المهمة

**AR:** 1. خلّي 3 درونات تسجل القياسات في map واحدة.  
2. احمِ السجل من الكتابة المتزامنة.  
3. ابعت مكان أكبر كريستال.

الـstarter بيقيس الصفوف ويسجل المواقع بالفعل.

**EN:** 1. Have 3 drones record in one map.  
2. Protect concurrent writes.  
3. Report the largest crystal's location.

The starter already measures and records locations.

### Idea / الفكرة

**AR:** الـmap المشتركة مش آمنة مع كتابات متزامنة.

النمط المعتاد: حط ~sync.Mutex~ جنب الـmap في ~Registry~.  
كل method بتوصل للبيانات تقفل نفس القفل؛ receiver من نوع ~*Registry~ ما ينسخش القفل.

بعد ~wg.Wait()~، الكتابة خلصت؛ بعدها تقدر تقرأ السجل.

**EN:** Concurrent writes to one map need protection.

Keep ~sync.Mutex~ beside the map in ~Registry~.  
Each method uses that lock; a ~*Registry~ receiver does not copy it.

After ~wg.Wait()~, writes finish; then read the log.

### Try it / جرّب

**Predict:** «السجل التجريبي يضيف مفتاحين. إيه اللي يطبع؟» / “Two keys are added. What prints?”

| Options | Correct | One-sentence feedback |
|---|---:|---|
| 2 7 · 2 4 · 7 2 | 1 | فيه مفتاحين، وقيمة violet هي 7. / Two keys; violet has value 7. |
|  |  | المفتاحين موجودين؛ راجع القيمة المرتبطة بـviolet. / Check the value stored for violet. |
|  |  | البرنامج يطبع عدد المفاتيح قبل قيمة violet. / The key count prints first. |

**AR walkthrough text:** مفاتيح توضيحية، مش بيانات المهمة. ده ترتيب واحد ممكن للكتابتين.  
**EN walkthrough text:** Illustrative keys, not mission data. This is one possible write order.

<!-- go-example id=e122-registry imports=fmt,sync outcome=ok -->
```go
sizes := map[string]int{}
var mu sync.Mutex
var wg sync.WaitGroup
record := func(key string, size int) {
    mu.Lock()
    defer mu.Unlock()
    sizes[key] = size
}
wg.Add(2)
go func() { defer wg.Done(); record("amber", 4) }()
go func() { defer wg.Done(); record("violet", 7) }()
wg.Wait()
fmt.Println(len(sizes), sizes["violet"])
```

| Line | AR caption | EN caption | State / output |
|---:|---|---|---|
| 0 | قبل التشغيل. | Before execution. | — |
| 1 | السجل فاضي. | The log is empty. | sizes={} |
| 5 | عامل amber مسك القفل. | Amber's worker gets the lock. | sizes={} |
| 7 | اتسجلت قيمته. | Its value is recorded. | sizes={amber: 4} |
| 5 | violet مسك القفل بعده. | Violet gets the lock next. | sizes={amber: 4} |
| 7 | بعده اتسجل violet. | Violet is recorded next. | sizes={amber: 4, violet: 7} |
| 12 | كل الكتّاب خلصوا. | All writers finish. | sizes unchanged |
| 13 | يطبع 2 7. | Prints 2 7. | stdout: `2 7` |

### Watch out / خد بالك

**AR:** كل وصول متزامن لنفس الـmap يستخدم نفس القفل.

~Record~ يقفل الكتابة؛ method قراءة تحتاج حماية لو العمال لسه شغالين.

بعد ~Wait~، مفيش عامل بيكتب، فاقرأ بعدها.

**EN:** Concurrent map access uses one shared lock.

~Record~ locks writes; protect reads while workers run.

After ~Wait~, the writers have finished.

### Think & check / فكّر واختبر

**Plan:** سجل بأمان → استنى الكل → دور على الموقع → ابعته. / Record safely → wait → find the location → report.

**Q1:** لو كل goroutine معاها قفل لوحدها، هل تحمي map واحدة؟ / Do separate goroutine locks protect one map?  
Options: أيوه / Yes · لأ / No · بس وقت القراءة / Only while reading. **Answer: 2 / No.**  
Feedback: الأقفال المختلفة مش بتنسّق على نفس البيانات. / Different locks do not coordinate.

**Q2:** بعد ما ~wg.Wait()~ يرجع، هل القراءة آمنة؟ / After Wait returns, can you read safely?  
Options: أيوه، لو مفيش كاتب تاني / Yes, if no writer remains · لأ، map لا تُقرأ / No, never · بس بقفل مختلف / Only with a new lock. **Answer: 1.**  
Feedback: صح، الانتظار خلّص كل الكتّاب؛ مفيش كتابة متزامنة. / Right; every writer has finished.

### Your turn / دورك

أمّن تسجيل الـstarter، ثم أضف قراءة للموقع بعد انتهاء العمال.  
Protect the starter's recording, then read a location after the workers finish.

### Hints

1. **Concept / الفكرة:** الكتّاب المتزامنين في map واحدة محتاجين mutex مشتركة. / Concurrent writers to one map need one shared mutex.
2. **Structure / الترتيب:** خلّي Registry تملك الـmap والقفل، وافصل التسجيل عن القراءة. / Let Registry own the map and lock; separate recording from reading.
3. **Specific / السطر:** راجع ~Record~ وسطر ~r.sizes[p] = size~. / Inspect ~Record~ and the ~r.sizes[p] = size~ line.

#### Verified e122-registry

Real stdout: 2 7

## 12-3 — الجمود (`12-3-deadlock`)

**Words:** Arabic 290; English 318.

**Objective:** recognize the possible wait cycle created by inconsistent lock order. Emphasize that Go does not promise a goroutine schedule; a fixed lock order prevents the two-lock cycle. The starter has two workers and two shared locks.

### Mission / المهمة

**AR:** 1. رجّع حركة الدرونات اللي وقفت.  
2. حلّ الـdeadlock بين قفلين.  
3. احصد الصفين.

الـstarter فيه عاملين بيطلبوا المثقاب والعربية بترتيبين مختلفين.

**EN:** 1. Get the stopped drones moving.  
2. Fix a two-lock deadlock.  
3. Harvest both rows.

The workers request locks in different orders.

### Idea / الفكرة

**AR:** Deadlock دايرة انتظار: goroutine ماسكة قفل ومستنية قفل تاني.

لو عاملين أخدوا القفلين بترتيب عكسي، ممكن كل واحد يستنى التاني.

Go مش بتضمن مين يبدأ؛ ده سيناريو محتمل، مش نتيجة كل تشغيل.

القاعدة: اطلب الموارد بترتيب ثابت في كل المسارات.

**EN:** A deadlock is a wait cycle: a goroutine holds one lock and waits for another.

If workers acquire locks in reverse orders, each may wait for the other.

Go does not promise who starts; this is one possible order.

Request resources in one fixed order on every path.

### Try it / جرّب

**Predict:** «العاملين بيطلبوا القفلين بنفس الترتيب. إيه اللي يطبع؟» / “Both workers use one lock order. What prints?”

| Options | Correct | One-sentence feedback |
|---|---:|---|
| 2 · 1 · البرنامج يعلّق / 2 · 1 · The program hangs | 1 | العاملين بينفذوا الزيادة؛ القفلين مش بيعملوا دورة انتظار. / Both workers increment; there is no wait cycle. |
|  |  | كل عامل بيزود done مرة واحدة. / Each worker increments once. |
|  |  | الترتيب موحد، فمش هتتكوّن دورة بين القفلين. / One order prevents a lock cycle. |

**AR walkthrough text:** ترتيب واحد ممكن: العامل الأول يخلص ثم الثاني. الناتج ثابت، والجدولة ممكن تختلف.  
**EN walkthrough text:** One possible order: worker one finishes, then two. Output is fixed; scheduling may vary.

<!-- go-example id=e123-lock-order imports=fmt,sync outcome=ok -->
```go
var first, second sync.Mutex
var wg sync.WaitGroup
done := 0
work := func() {
    first.Lock()
    second.Lock()
    done++
    second.Unlock()
    first.Unlock()
}
wg.Add(2)
go func() { defer wg.Done(); work() }()
go func() { defer wg.Done(); work() }()
wg.Wait()
fmt.Println(done)
```

| Line | AR caption | EN caption | State / output |
|---:|---|---|---|
| 0 | قبل التشغيل. | Before execution. | — |
| 3 | العداد يبدأ صفر. | Counter starts at zero. | done=0 |
| 5 | ترتيب ممكن: الأول يأخذ القفل. | One order: worker one locks. | done=0 |
| 6 | بعده يأخذ القفل التاني. | Then it gets the second lock. | done=0 |
| 7 | العداد صار 1. | Counter becomes 1. | done=1 |
| 8 | يفتح القفل التاني. | It releases the second lock. | done=1 |
| 5 | العامل التاني يدخل بعده. | Worker two enters next. | done=1 |
| 7 | العداد صار 2. | Counter becomes 2. | done=2 |
| 15 | يطبع 2. | Prints 2. | stdout: `2` |

### Watch out / خد بالك

**AR:** ~WaitGroup~ مش بيفك deadlock؛ هو بس يستنى النهاية.

كل المسارات لازم تطلب الأقفال بنفس الترتيب.

افتح الأقفال بعد الشغل المحمي وبالترتيب العكسي.

**EN:** ~WaitGroup~ cannot break a deadlock; it only waits.

Every path needs one lock order.

Unlock after protected work, in reverse order.

### Think & check / فكّر واختبر

**Plan:** حدد مين ماسك كل قفل → امنع دورة الانتظار → اختبر حصاد الصفين. / Find held locks → prevent the cycle → test both rows.

**Q1:** إيه اللي يمنع دورة انتظار بين قفلين؟ / What prevents a two-lock cycle?  
Options: ترتيب ثابت واحد / One fixed order · ترتيب مختلف لكل عامل / A different order per worker · WaitGroup إضافية / Another WaitGroup. **Answer: 1.**  
Feedback: الترتيب الموحد يمنع اختلاف ترتيب الطلب. / One shared order prevents conflicting requests.

**Q2:** هل Go تضمن مين يبدأ في الـwalkthrough؟ / Does Go promise who starts?  
Options: لأ، ده ترتيب واحد ممكن / No, this is one possible order · أيوه، العامل الأول دايماً / Yes, worker one always starts · WaitGroup بترتبهم / WaitGroup chooses the order. **Answer: 1.**  
Feedback: صح، goroutines ممكن تتجدول بترتيبات مختلفة. / Right, goroutine scheduling can vary.

### Your turn / دورك

عدّل الـstarter وشغّله؛ راقب حركة الدرونين وإتمام الصفين.  
Edit the starter, then watch both drones finish the rows.

### Hints

1. **Concept / الفكرة:** الـdeadlock ممكن يحصل لما كل عامل ماسك مورد ومستني التاني. / A deadlock can form when each worker holds one resource and waits for another.
2. **Structure / الترتيب:** قارن ترتيب القفلين في كل مسار قبل تتبّع Wait. / Compare the two lock orders on every path before tracing Wait.
3. **Specific / السطر:** راجع arguments بتوع ~harvestWith~ في استدعاءي ~work~ جوه ~main~. / Inspect the arguments to ~harvestWith~ in both ~work~ calls inside ~main~.

#### Verified e123-lock-order

Real stdout: 2
