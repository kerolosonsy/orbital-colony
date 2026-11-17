# Chapter 11 lesson storyboard

Revision 1. Each mission uses the six approved steps. Word counts follow the 6-1 convention: titles, learner text, frame captions, question options and feedback, and hints count; code does not. All examples use illustrative values. No mission threshold or completed game algorithm appears here.

## 11-1 — الكشّاف (`11-1-scout`)

**Word count:** Arabic 285; English 311.

### Arabic learner copy

#### Step `c111.mission-step` — المهمة

`c111.mission` · markdown

1. خلّي scout يبعث أماكن الكريستالات الكبيرة.
2. استقبل كل مكان ووجّه درونك له.
3. احصد الأماكن اللي وصلتك بس.

الـstarter بدأ القياس والإرسال، ومستنيك تكمل الاستقبال.

#### Step `c111.idea-step` — الفكرة

`c111.channel-trace` · walkthrough

الـchannel بتنقل قيمة بين goroutines.

من غير buffer، الإرسال والاستقبال بيستّنوا بعض.

الـscout يقيس؛ الدرون التاني يستقبل المكان ويحصد.

هنا فيه رسالتين في الطابور؛ `range` يقرأهم، وبعد القفل يعرف إن مفيش رسائل تانية.

<!-- go-example id=e11-scout-channel imports=fmt outcome=ok -->
```go
found := make(chan int, 2)
found <- 3
found <- 8
close(found)
for x := range found {
	fmt.Println(x)
}
```

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | قبل التشغيل. | — |
| 1 | قناة فاضية. | — |
| 2 | وصل أول مكان. | — |
| 3 | وصل تاني مكان. | — |
| 4 | المرسل خلّص. | — |
| 5 | استقبل 3. | — |
| 6 | ظهر أول رقم. | `3` |
| 5 | استقبل 8. | `3` |
| 6 | ظهر الرقم التاني. | `3` then `8` |

#### Step `c111.try-step` — جرّب

`c111.predict` · question

قبل التتبع، توقّع ترتيب السطرين في الـoutput.

<!-- go-example id=e11-scout-order imports=fmt outcome=ok -->
```go
places := make(chan int, 2)
places <- 4
places <- 7
close(places)
fmt.Println(<-places)
fmt.Println(<-places)
```

1. 4 وبعده 7 — رسائل الإرسال بتتقرأ بالترتيب. **Correct.**
2. 7 وبعده 4 — الاستقبال مش بيقلب ترتيب الرسائل.
3. 0 وبعده 0 — القناة فيها قيم اتبعتت قبل القفل.

`c111.reveal` · walkthrough

الرسالتين بيوصلوا بالترتيب، وبعد استهلاكهم القناة المقفولة بتنهي القراءة. يعيد التتبع المثال أعلاه.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | قبل التشغيل. | — |
| 1 | قناة فاضية. | — |
| 2 | اتضاف 4. | — |
| 3 | اتضاف 7. | — |
| 4 | القناة اتقفلت. | — |
| 5 | اتطبع 4. | `4` |
| 6 | اتطبع 7. | `4` then `7` |

#### Step `c111.watch-step` — خد بالك

`c111.watch` · markdown

- اللي بيبعت هو اللي يقفل القناة بعد آخر رسالة.
- الإرسال على channel مقفولة يعمل panic.

#### Step `c111.check-step` — فكّر واختبر

`c111.plan` · markdown

خلي scout يقيس، والـchannel تنقل المكان، والدرون التاني يستقبل ويحصد.

استقبل البلاغات لحد نهاية القناة، وخلي كل درون مسؤول عن دوره.

`c111.q-range` · question — إمتى `range found` توقف؟

1. بعد القفل وتفريغ الرسائل — بتقرأ كل الموجود، وبعدها القفل بينهي الحلقة. **Correct.**
2. بعد أول رسالة — رسالة واحدة مش كفاية لو فيه غيرها.
3. أول ما الطابور يفضى — القناة المفتوحة ممكن تستقبل رسائل تانية.

`c111.q-close` · question — مين يقفل `found`؟

1. المرسل — المرسل عارف إنه خلّص إرسال. **Correct.**
2. المستقبل — المستقبل يقرأ، لكنه ما يعلنش النهاية.
3. أي طرف وقت ما يحب — الإغلاق المبكر ممكن يكسر إرسال لسه جاي.

#### Step `c111.turn-step` — دورك

`c111.turn` · markdown

كمّل الاستقبال وشغّل البرنامج؛ راقب كل رسالة وحركة الدرون.

### Three position-labelled hints

1. **مفهوم:** الـchannel بتنقل قيمة من goroutine للتانية؛ و`close` بيعلن النهاية.
2. **الترتيب:** جهّز إرسال الأماكن، واستقبلها، واحصد، واستنى القناة تخلص.
3. **سطر محدد:** راجع `x := <-found`؛ إيه اللي يحصل بعد أول استقبال؟

### English learner copy

#### Step `c111.mission-step` — Mission

`c111.mission` · markdown

1. Scout sends locations of large crystals.
2. Your drone visits each report.
3. Harvest reported locations.

The starter measures and sends; add the receiver.

#### Step `c111.idea-step` — Idea

`c111.channel-trace` · walkthrough

A channel carries values between goroutines.

Without a buffer, send and receive wait.

The scout measures; the other drone receives the location and harvests.

`range` reads queued messages; close marks the end.

Code and line references: canonical `e11-scout-channel` above.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | Before execution. | — |
| 1 | Empty channel. | — |
| 2 | First place sent. | — |
| 3 | Second place sent. | — |
| 4 | Sender is done. | — |
| 5 | Receive 3. | — |
| 6 | First number prints. | `3` |
| 5 | Receive 8. | `3` |
| 6 | Second number prints. | `3` then `8` |

#### Step `c111.try-step` — Try it

`c111.predict` · question

Before the trace, predict the order of the two output lines.

Code: canonical `e11-scout-order` above.

1. 4, then 7 — Sent messages are received in order. **Correct.**
2. 7, then 4 — Receiving does not reverse the messages.
3. 0, then 0 — The channel contains values sent before close.

`c111.reveal` · walkthrough

Messages arrive in order; close ends the loop after both are read. The trace repeats the example above.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | Before execution. | — |
| 1 | Empty channel. | — |
| 2 | 4 is queued. | — |
| 3 | 7 is queued. | — |
| 4 | Channel closes. | — |
| 5 | 4 prints. | `4` |
| 6 | 7 prints. | `4` then `7` |

#### Step `c111.watch-step` — Watch out

`c111.watch` · markdown

- The sender closes the channel after its last message.
- Sending on a closed channel causes a panic.

#### Step `c111.check-step` — Think & check

`c111.plan` · markdown

Scout measures; the channel carries places; the other drone receives and harvests.

Read reports until sending ends.

`c111.q-range` · question — When does `range found` stop?

1. After close and draining messages — It reads queued values, then close ends the loop. **Correct.**
2. After the first message — One message is not enough if more remain.
3. As soon as the queue is empty — An open channel can still receive messages.

`c111.q-close` · question — Who closes `found`?

1. The sender — The sender knows when it has finished sending. **Correct.**
2. The receiver — The receiver reads but does not announce the end.
3. Either side, whenever — Early close can break a later send.

#### Step `c111.turn-step` — Your turn

`c111.turn` · markdown

Finish the receiver; run it and follow each report.

### Three position-labelled hints

1. **Concept:** Channels pass values; `close` marks the end.
2. **Structure:** Send places, receive them, harvest, then await closure.
3. **Specific:** Check `x := <-found`; what should happen after the first receive?

## 11-2 — طابور الشغل (`11-2-worker-pool`)

**Word count:** Arabic 302; English 318.

### Arabic learner copy

#### Step `c112.mission-step` — المهمة

`c112.mission` · markdown

1. 3 workers يسحبوا أهداف من channel واحدة.
2. كل واحد يرجّع ناتج حصاده.
3. اجمع النواتج وابعت المجموع للقاعدة.

الـstarter مجهّز worker والقنوات وWaitGroup؛ الـmain لسه فاضي.

#### Step `c112.idea-step` — الفكرة

`c112.queue-trace` · walkthrough

الـworker pool طابور مشترك؛ العامل الفاضي ياخد المهمة الجاية.

`range jobs` يخلص بعد القفل وتفريغ الرسائل؛ و`results` بتنقل ناتج الحصاد.

`<-chan T` استقبال بس؛ `chan<- T` إرسال بس.

التتبع هنا لعامل واحد؛ الـpool يشغّل أكتر من عامل على القناة نفسها.

<!-- go-example id=e11-pool-queue imports=fmt outcome=ok -->
```go
jobs := make(chan int, 2)
jobs <- 3
jobs <- 5
close(jobs)
for job := range jobs {
	fmt.Println(job * 2)
}
```

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | قبل التشغيل. | — |
| 1 | طابور فاضي. | — |
| 2 | دخلت مهمة 3. | — |
| 3 | دخلت مهمة 5. | — |
| 4 | المرسل قفل الطابور. | — |
| 5 | العامل استلم 3. | — |
| 6 | الناتج 6. | `6` |
| 5 | العامل استلم 5. | `6` |
| 6 | الناتج 10. | `6` then `10` |

#### Step `c112.try-step` — جرّب

`c112.predict` · question

الـworkerين بعتوا 2 و3؛ هتطبع كام لو وصلوا بأي ترتيب؟

<!-- go-example id=e11-pool-total imports=fmt outcome=ok -->
```go
results := make(chan int, 2)
go func() { results <- 2 }()
go func() { results <- 3 }()
a := <-results
b := <-results
fmt.Println(a + b)
```

1. 5 — الجمع يوصل لنفس الإجمالي بأي ترتيب. **Correct.**
2. 23 — القيم بتتجمع كأرقام، مش بتتلزق كنص.
3. 3 — النتيجتين بيتضافوا للإجمالي.

`c112.reveal` · walkthrough

ده ترتيب استقبال واحد ممكن؛ Go مش بيضمن مين يوصل الأول. يعيد التتبع المثال أعلاه.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | قبل التشغيل. | — |
| 1 | قناة النتائج فاضية. | — |
| 2 | الأول بعت 2. | — |
| 3 | التاني بعت 3. | — |
| 4 | استقبلت a قيمة 2. | — |
| 5 | استقبلت b قيمة 3. | — |
| 6 | الإجمالي اتطبع. | `5` |

#### Step `c112.watch-step` — خد بالك

`c112.watch` · markdown

- مرسل `jobs` يقفلها بعد آخر هدف.
- استخدم WaitGroup قبل قفل `results`؛ من غير القفل، `range` هتفضل مستنية.

#### Step `c112.check-step` — فكّر واختبر

`c112.plan` · markdown

رتّب تشغيل العمال، إرسال الأهداف، جمع النواتج، وإغلاق القنوات بعد اكتمال دورها.

`c112.q-close` · question — إمتى تقفل `results`؟

1. بعد ما كل العمال يخلصوا — استنى المجموعة كلها قبل ما تعلن النهاية. **Correct.**
2. بعد أول نتيجة — ممكن لسه فيه عامل هيبعت نتيجة.
3. قبل تشغيل العمال — العمال لسه محتاجين القناة المفتوحة.

`c112.q-range` · question — إيه اللي ينهي `range jobs`؟

1. القناة تتقفل بعد تفريغها — العمال يكملوا اللي في الطابور قبل ما يقفوا. **Correct.**
2. أول عامل يسحب هدف — باقي العمال يقدروا ياخدوا شغل كمان.
3. انتهاء أول عامل — عامل واحد خلص مش معناه إن الطابور خلص.

#### Step `c112.turn-step` — دورك

`c112.turn` · markdown

كمّل `main` وشغّلها؛ راقب استقبال الأهداف والنتائج لحد التقرير.

### Three position-labelled hints

1. **مفهوم:** كل worker يقدر يسحب الشغل التالي من channel مشتركة.
2. **الترتيب:** شغّل العمال، وابعت الأهداف، واجمع النتائج، واقفل القنوات في وقتها.
3. **سطر محدد:** راجع توقيع `worker`؛ أسهم القنوات بتوضح مين يستقبل ومين يبعت.

### English learner copy

#### Step `c112.mission-step` — Mission

`c112.mission` · markdown

1. Three workers share the target queue.
2. Each sends its harvest count.
3. Sum and report.

The starter defines `worker`; `main` is empty.

#### Step `c112.idea-step` — Idea

`c112.queue-trace` · walkthrough

A pool shares jobs; a free worker takes the next.

`range jobs` ends after close and draining; `results` carries harvest counts.

`<-chan T` receives only; `chan<- T` sends only.

This trace shows one worker.

Code and line references: canonical `e11-pool-queue` above.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | Before execution. | — |
| 1 | Queue is empty. | — |
| 2 | Job 3 arrives. | — |
| 3 | Job 5 arrives. | — |
| 4 | Sender closes queue. | — |
| 5 | Worker gets 3. | — |
| 6 | Result is 6. | `6` |
| 5 | Worker gets 5. | `6` |
| 6 | Result is 10. | `6` then `10` |

#### Step `c112.try-step` — Try it

`c112.predict` · question

Workers send 2 and 3; what total prints in either receive order?

Code: canonical `e11-pool-total` above.

1. 5 — Addition gives the same total in either order. **Correct.**
2. 23 — The values are added as numbers, not joined as text.
3. 3 — Both results add to the total.

`c112.reveal` · walkthrough

One possible receive order; Go does not guarantee it. The trace repeats the example above.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | Before run. | — |
| 1 | Results empty. | — |
| 2 | First sends 2. | — |
| 3 | Second sends 3. | — |
| 4 | `a` receives 2. | — |
| 5 | `b` receives 3. | — |
| 6 | Sum prints. | `5` |

#### Step `c112.watch-step` — Watch out

`c112.watch` · markdown

- The `jobs` sender closes it after the last target.
- Use WaitGroup before closing `results`; without close, `range` waits.

#### Step `c112.check-step` — Think & check

`c112.plan` · markdown

Start workers, send targets, collect results, then close channels after work ends.

`c112.q-close` · question — When should `results` close?

1. After every worker finishes — Wait for every worker before signaling the end. **Correct.**
2. After the first result — A worker may still send a result.
3. Before workers start — Workers need the channel open.

`c112.q-range` · question — What ends `range jobs`?

1. The channel closes after draining — Workers drain queued jobs before stopping. **Correct.**
2. The first worker takes a job — Other workers can take more jobs.
3. The first worker finishes — One worker finishing does not empty it.

#### Step `c112.turn-step` — Your turn

`c112.turn` · markdown

Finish `main` and run it; watch jobs and results arrive through the report.

### Three position-labelled hints

1. **Concept:** Each worker can take the next job from one shared channel.
2. **Structure:** Start workers, send targets, collect results, and close channels at the right time.
3. **Specific:** Check the `worker` signature; channel arrows show which side receives or sends.

## 11-3 — العاصفة جاية (`11-3-storm`)

**Word count:** Arabic 261; English 315.

### Arabic learner copy

#### Step `c113.mission-step` — المهمة

`c113.mission` · markdown

1. اجمع الكريستال المطلوب.
2. عند الإنذار، ارجع للقاعدة قبل الضربة.
3. بعد هدوء العاصفة، كمّل الشغل.

الـstarter جهّز قناتي التحذير والهدوء؛ دورك تدير الرجوع والاستئناف.

#### Step `c113.idea-step` — الفكرة

`c113.select-trace` · walkthrough

`select` شبه `switch` للقنوات: يستنى عملية جاهزة.

مع `default`، لو مفيش تحذير جاهز، يكمل فورًا.

قناة `colony.StormClear()` بتتقفل عند نهاية العاصفة؛ الاستقبال منها ساعتها يرجع فورًا.

<!-- go-example id=e11-storm-select imports=fmt outcome=ok -->
```go
warning := make(chan struct{}, 1)
select {
case <-warning:
	fmt.Println("return")
default:
	fmt.Println("work")
}
warning <- struct{}{}
select {
case <-warning:
	fmt.Println("return")
default:
	fmt.Println("work")
}
```

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | قبل التشغيل. | — |
| 1 | قناة التحذير فاضية. | — |
| 2 | أول select يفحص. | — |
| 5 | مفيش تحذير؛ default جاهزة. | — |
| 6 | الشغل مستمر. | `work` |
| 8 | وصل تحذير للقناة. | `work` |
| 9 | select يفحص تاني. | `work` |
| 10 | وصلت قيمة التحذير. | `work` |
| 11 | الرجوع اتطبع. | `work` then `return` |

#### Step `c113.try-step` — جرّب

`c113.predict` · question

قناة `calm` اتقفلت؛ توقّع أي فرع يطبع.

<!-- go-example id=e11-storm-clear imports=fmt outcome=ok -->
```go
calm := make(chan struct{})
close(calm)
select {
case <-calm:
	fmt.Println("safe")
default:
	fmt.Println("work")
}
```

1. safe — الاستقبال من قناة مقفولة جاهز فورًا. **Correct.**
2. work — `default` تعمل بس لو مفيش استقبال جاهز.
3. مفيش output — القناة المقفولة مش بتعلّق الاستقبال.

`c113.reveal` · walkthrough

القفل بيخلّي الاستقبال جاهز، فـ`select` يدخل الحالة دي. يعيد التتبع المثال أعلاه.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | قبل التشغيل. | — |
| 1 | اتعملت قناة فاضية. | — |
| 2 | القناة اتقفلت. | — |
| 3 | select يفحص. | — |
| 4 | الاستقبال جاهز. | — |
| 5 | اتطبع safe. | `safe` |

#### Step `c113.watch-step` — خد بالك

`c113.watch` · markdown

- `default` ما بتستناش؛ خلّي الحلقة تعمل شغل بعد الفحص.
- حدّث channel الهدوء بعد كل إنذار عشان تستنى العاصفة الجاية.

#### Step `c113.check-step` — فكّر واختبر

`c113.plan` · markdown

اشتغل بين الإنذارات، ارجع وقت التحذير، استنى الهدوء، وبعدها استأنف.

`c113.q-default` · question — لو `warning` لسه فاضية، `default` تعمل إيه؟

1. تتنفذ فورًا — مفيش عملية جاهزة، فـ`default` تكمل فورًا. **Correct.**
2. تستنى التحذير — من غير `default`، `select` هيستنى.
3. تقفل القناة — `default` مش بتغيّر حالة القناة.

`c113.q-clear` · question — إمتى `<-calm` يرجع بعد القفل؟

1. فورًا — الاستقبال من قناة مقفولة جاهز دايمًا. **Correct.**
2. بعد رسالة جديدة — القفل نفسه علامة انتهاء الإشارة.
3. بعد ما يتقفل تاني — القناة المقفولة مش محتاجة قفل تاني.

#### Step `c113.turn-step` — دورك

`c113.turn` · markdown

كمّل مسار الرجوع والاستئناف وشغّل المهمة؛ تابع كل إنذار.

### Three position-labelled hints

1. **مفهوم:** `select` تتعامل مع حالة channel الجاهزة، و`default` تكمّل فورًا.
2. **الترتيب:** اشتغل بين الإنذارات، ارجع، استنى الهدوء، وبعدها استأنف.
3. **سطر محدد:** راجع `colony.StormWarning()` و`colony.StormClear()` لإشارات التحذير والهدوء.

### English learner copy

#### Step `c113.mission-step` — Mission

`c113.mission` · markdown

1. Collect the required crystals.
2. Return to base when warned.
3. Continue after the storm clears.

The starter has both channels; manage the return and resume.

#### Step `c113.idea-step` — Idea

`c113.select-trace` · walkthrough

`select` is like `switch` for channels: it waits for a ready operation.

With `default`, it continues at once when no warning is ready.

The `colony.StormClear()` channel closes when the storm ends; receiving then returns immediately.

Code and line references: canonical `e11-storm-select` above.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | Before execution. | — |
| 1 | Warning channel is empty. | — |
| 2 | First select checks. | — |
| 5 | No warning; default is ready. | — |
| 6 | Work continues. | `work` |
| 8 | Warning enters channel. | `work` |
| 9 | Select checks again. | `work` |
| 10 | Warning is received. | `work` |
| 11 | Return prints. | `work` then `return` |

#### Step `c113.try-step` — Try it

`c113.predict` · question

The `calm` channel is closed; predict which branch prints.

Code: canonical `e11-storm-clear` above.

1. safe — A receive from a closed channel is ready immediately. **Correct.**
2. work — `default` runs only when no receive is ready.
3. No output — A closed channel does not block the receive.

`c113.reveal` · walkthrough

Closing makes the receive ready, so `select` takes that case. The trace repeats the example above.

| Line | Caption | Cumulative output |
|---:|---|---|
| 0 | Before execution. | — |
| 1 | Empty channel created. | — |
| 2 | Channel closes. | — |
| 3 | Select checks. | — |
| 4 | Receive is ready. | — |
| 5 | safe prints. | `safe` |

#### Step `c113.watch-step` — Watch out

`c113.watch` · markdown

- `default` does not wait; make the loop do work after checking.
- Refresh the calm channel after each warning to await another storm.

#### Step `c113.check-step` — Think & check

`c113.plan` · markdown

Work between warnings, return when warned, wait for calm, then resume.

`c113.q-default` · question — If `warning` is empty, what does `default` do?

1. Runs immediately — No operation is ready, so `default` runs now. **Correct.**
2. Waits for warning — Without `default`, `select` would wait.
3. Closes the channel — `default` does not change channel state.

`c113.q-clear` · question — After close, when does `<-calm` return?

1. Immediately — A receive from a closed channel is always ready. **Correct.**
2. After a new message — Close itself signals the end.
3. After another close — A closed channel does not need another close.

#### Step `c113.turn-step` — Your turn

`c113.turn` · markdown

Finish the return and resume path and run the mission; follow each warning.

### Three position-labelled hints

1. **Concept:** `select` handles ready channel cases; `default` lets the loop continue now.
2. **Structure:** Work between warnings, return, wait for calm, then resume.
3. **Specific:** Check `colony.StormWarning()` and `colony.StormClear()` for both signals.

## Later playground ideas

- **Scout:** Let learners place a small set of location messages on a channel, then step through receiving, closing, and draining it. Keep the values illustrative and make the sender/receiver roles visible.
- **Worker pool:** Show a few generic jobs being pulled by multiple workers, with a results counter and a WaitGroup gate that controls when results close. Display one possible assignment at a time and label the worker order as nondeterministic.
- **Storm:** Provide a warning toggle and a closed/open calm signal; let learners choose between `select` with `default` and a blocking receive, then show whether work continues or pauses.

## Example verification

Each marker above is a canonical fragment referenced by the English copy. The chapter test builds every fragment with the Go toolchain and compares real stdout with the walkthrough's terminal output.

#### Verified `e11-scout-channel` — ok

Real stdout:

```text
3
8
```

#### Verified `e11-scout-order` — ok

Real stdout:

```text
4
7
```

#### Verified `e11-pool-queue` — ok

Real stdout:

```text
6
10
```

#### Verified `e11-pool-total` — ok

Real stdout:

```text
5
```

#### Verified `e11-storm-select` — ok

Real stdout:

```text
work
return
```

#### Verified `e11-storm-clear` — ok

Real stdout:

```text
safe
```
