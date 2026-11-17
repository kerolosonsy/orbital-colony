package levels

const ch11ScoutChannelExample = `found := make(chan int, 2)
found <- 3
found <- 8
close(found)
for x := range found {
	fmt.Println(x)
}`

const ch11ScoutTryExample = `places := make(chan int, 2)
places <- 4
places <- 7
close(places)
fmt.Println(<-places)
fmt.Println(<-places)`

const ch11PoolQueueExample = `jobs := make(chan int, 2)
jobs <- 3
jobs <- 5
close(jobs)
for job := range jobs {
	fmt.Println(job * 2)
}`

const ch11PoolTryExample = `results := make(chan int, 2)
go func() { results <- 2 }()
go func() { results <- 3 }()
a := <-results
b := <-results
fmt.Println(a + b)`

const ch11StormSelectExample = `warning := make(chan struct{}, 1)
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
}`

const ch11StormTryExample = `calm := make(chan struct{})
close(calm)
select {
case <-calm:
	fmt.Println("safe")
default:
	fmt.Println("work")
}`

func ch11PoolTryTrace(captions []string) *LessonWalkthrough {
	return &LessonWalkthrough{Frames: []LessonFrame{
		{Line: 0, Caption: captions[0]},
		{Line: 1, Caption: captions[1], Variables: []LessonVariable{{Name: "len(results)", Value: "0"}}},
		{Line: 2, Caption: captions[2], Variables: []LessonVariable{{Name: "len(results)", Value: "1"}}},
		{Line: 3, Caption: captions[3], Variables: []LessonVariable{{Name: "len(results)", Value: "2"}}},
		{Line: 4, Caption: captions[4], Variables: []LessonVariable{{Name: "a", Value: "2"}, {Name: "len(results)", Value: "1"}}},
		{Line: 5, Caption: captions[5], Variables: []LessonVariable{{Name: "b", Value: "3"}, {Name: "len(results)", Value: "0"}}},
		{Line: 6, Caption: captions[6], Variables: []LessonVariable{{Name: "a + b", Value: "5"}}, Output: "5"},
	}}
}

var ch11ScoutLessonArabic = []LessonStep{
	{
		ID: "c111.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c111.mission", Kind: "markdown",
			Text: "1. خلّي scout يبعث أماكن الكريستالات الكبيرة.\n2. استقبل كل مكان ووجّه درونك له.\n3. احصد الأماكن اللي وصلتك بس.\n\nالـstarter بدأ القياس والإرسال، ومستنيك تكمل الاستقبال.",
		}},
	},
	{
		ID: "c111.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c111.channel-trace", Kind: "walkthrough",
			Text: "الـchannel بتنقل قيمة بين goroutines.\n\nمن غير buffer، الإرسال والاستقبال بيستّنوا بعض.\n\nالـscout يقيس؛ الدرون التاني يستقبل المكان ويحصد.\n\nهنا فيه رسالتين في الطابور؛ `range` يقرأهم، وبعد القفل يعرف إن مفيش رسائل تانية.",
			Code: ch11ScoutChannelExample, Imports: []string{"fmt"}, ExampleID: "e11-scout-channel",
			Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل التشغيل."},
				{Line: 1, Caption: "قناة فاضية.", Variables: []LessonVariable{{Name: "queued", Value: "0"}}},
				{Line: 2, Caption: "وصل أول مكان.", Variables: []LessonVariable{{Name: "queued", Value: "1"}}},
				{Line: 3, Caption: "وصل تاني مكان.", Variables: []LessonVariable{{Name: "queued", Value: "2"}}},
				{Line: 4, Caption: "المرسل خلّص.", Variables: []LessonVariable{{Name: "queued", Value: "2"}, {Name: "closed", Value: "true"}}},
				{Line: 5, Caption: "استقبل 3.", Variables: []LessonVariable{{Name: "x", Value: "3"}, {Name: "queued", Value: "1"}}},
				{Line: 6, Caption: "ظهر أول رقم.", Variables: []LessonVariable{{Name: "x", Value: "3"}, {Name: "queued", Value: "1"}}, Output: "3"},
				{Line: 5, Caption: "استقبل 8.", Variables: []LessonVariable{{Name: "x", Value: "8"}, {Name: "queued", Value: "0"}}, Output: "3"},
				{Line: 6, Caption: "ظهر الرقم التاني.", Variables: []LessonVariable{{Name: "x", Value: "8"}, {Name: "queued", Value: "0"}}, Output: "3\n8"},
			}},
		}},
	},
	{
		ID: "c111.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID: "c111.predict", Kind: "question",
				Text: "قبل التتبع، توقّع ترتيب السطرين في الـoutput.",
				Code: ch11ScoutTryExample, Imports: []string{"fmt"}, ExampleID: "e11-scout-order",
				Question: &LessonQuestion{
					Options: []string{"4 وبعده 7", "7 وبعده 4", "0 وبعده 0"}, Answer: 0,
					Feedback: []string{"رسائل الإرسال بتتقرأ بالترتيب.", "الاستقبال مش بيقلب ترتيب الرسائل.", "القناة فيها قيم اتبعتت قبل القفل."},
				},
			},
			{
				ID: "c111.reveal", Kind: "walkthrough",
				Text: "الرسالتين بيوصلوا بالترتيب، وبعد استهلاكهم القناة المقفولة بتنهي القراءة.",
				Code: ch11ScoutTryExample, Imports: []string{"fmt"}, ExampleID: "e11-scout-order",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "قناة فاضية."},
					{Line: 2, Caption: "اتضاف 4."},
					{Line: 3, Caption: "اتضاف 7."},
					{Line: 4, Caption: "القناة اتقفلت."},
					{Line: 5, Caption: "اتطبع 4.", Output: "4"},
					{Line: 6, Caption: "اتطبع 7.", Output: "4\n7"},
				}},
			},
		},
	},
	{
		ID: "c111.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c111.watch", Kind: "markdown",
			Text: "- اللي بيبعت هو اللي يقفل القناة بعد آخر رسالة.\n- الإرسال على channel مقفولة يعمل panic.",
		}},
	},
	{
		ID: "c111.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c111.plan", Kind: "markdown", Text: "خلي scout يقيس، والـchannel تنقل المكان، والدرون التاني يستقبل ويحصد.\n\nاستقبل البلاغات لحد نهاية القناة، وخلي كل درون مسؤول عن دوره."},
			{ID: "c111.q-range", Kind: "question", Text: "إمتى `range found` توقف؟", Question: &LessonQuestion{
				Options: []string{"بعد القفل وتفريغ الرسائل", "بعد أول رسالة", "أول ما الطابور يفضى"}, Answer: 0,
				Feedback: []string{"بتقرأ كل الموجود، وبعدها القفل بينهي الحلقة.", "رسالة واحدة مش كفاية لو فيه غيرها.", "القناة المفتوحة ممكن تستقبل رسائل تانية."},
			}},
			{ID: "c111.q-close", Kind: "question", Text: "مين يقفل `found`؟", Question: &LessonQuestion{
				Options: []string{"المرسل", "المستقبل", "أي طرف وقت ما يحب"}, Answer: 0,
				Feedback: []string{"المرسل عارف إنه خلّص إرسال.", "المستقبل يقرأ، لكنه ما يعلنش النهاية.", "الإغلاق المبكر ممكن يكسر إرسال لسه جاي."},
			}},
		},
	},
	{
		ID: "c111.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c111.turn", Kind: "markdown", Text: "كمّل الاستقبال وشغّل البرنامج؛ راقب كل رسالة وحركة الدرون."}},
	},
}

var ch11ScoutLessonEnglish = []LessonStep{
	{
		ID: "c111.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c111.mission", Kind: "markdown",
			Text: "1. Scout sends locations of large crystals.\n2. Your drone visits each report.\n3. Harvest reported locations.\n\nThe starter measures and sends; add the receiver.",
		}},
	},
	{
		ID: "c111.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c111.channel-trace", Kind: "walkthrough",
			Text: "A channel carries values between goroutines.\n\nWithout a buffer, send and receive wait.\n\nThe scout measures; the other drone receives the location and harvests.\n\n`range` reads queued messages; close marks the end.",
			Code: ch11ScoutChannelExample, Imports: []string{"fmt"}, ExampleID: "e11-scout-channel",
			Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Empty channel.", Variables: []LessonVariable{{Name: "queued", Value: "0"}}},
				{Line: 2, Caption: "First place sent.", Variables: []LessonVariable{{Name: "queued", Value: "1"}}},
				{Line: 3, Caption: "Second place sent.", Variables: []LessonVariable{{Name: "queued", Value: "2"}}},
				{Line: 4, Caption: "Sender is done.", Variables: []LessonVariable{{Name: "queued", Value: "2"}, {Name: "closed", Value: "true"}}},
				{Line: 5, Caption: "Receive 3.", Variables: []LessonVariable{{Name: "x", Value: "3"}, {Name: "queued", Value: "1"}}},
				{Line: 6, Caption: "First number prints.", Variables: []LessonVariable{{Name: "x", Value: "3"}, {Name: "queued", Value: "1"}}, Output: "3"},
				{Line: 5, Caption: "Receive 8.", Variables: []LessonVariable{{Name: "x", Value: "8"}, {Name: "queued", Value: "0"}}, Output: "3"},
				{Line: 6, Caption: "Second number prints.", Variables: []LessonVariable{{Name: "x", Value: "8"}, {Name: "queued", Value: "0"}}, Output: "3\n8"},
			}},
		}},
	},
	{
		ID: "c111.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID: "c111.predict", Kind: "question",
				Text: "Before the trace, predict the order of the two output lines.",
				Code: ch11ScoutTryExample, Imports: []string{"fmt"}, ExampleID: "e11-scout-order",
				Question: &LessonQuestion{
					Options: []string{"4, then 7", "7, then 4", "0, then 0"}, Answer: 0,
					Feedback: []string{"Sent messages are received in order.", "Receiving does not reverse the messages.", "The channel contains values sent before close."},
				},
			},
			{
				ID: "c111.reveal", Kind: "walkthrough",
				Text: "Messages arrive in order; close ends the loop after both are read.",
				Code: ch11ScoutTryExample, Imports: []string{"fmt"}, ExampleID: "e11-scout-order",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Empty channel."},
					{Line: 2, Caption: "4 is queued."},
					{Line: 3, Caption: "7 is queued."},
					{Line: 4, Caption: "Channel closes."},
					{Line: 5, Caption: "4 prints.", Output: "4"},
					{Line: 6, Caption: "7 prints.", Output: "4\n7"},
				}},
			},
		},
	},
	{
		ID: "c111.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c111.watch", Kind: "markdown",
			Text: "- The sender closes the channel after its last message.\n- Sending on a closed channel causes a panic.",
		}},
	},
	{
		ID: "c111.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c111.plan", Kind: "markdown", Text: "Scout measures; the channel carries places; the other drone receives and harvests.\n\nRead reports until sending ends."},
			{ID: "c111.q-range", Kind: "question", Text: "When does `range found` stop?", Question: &LessonQuestion{
				Options: []string{"After close and draining messages", "After the first message", "As soon as the queue is empty"}, Answer: 0,
				Feedback: []string{"It reads queued values, then close ends the loop.", "One message is not enough if more remain.", "An open channel can still receive messages."},
			}},
			{ID: "c111.q-close", Kind: "question", Text: "Who closes `found`?", Question: &LessonQuestion{
				Options: []string{"The sender", "The receiver", "Either side, whenever"}, Answer: 0,
				Feedback: []string{"The sender knows when it has finished sending.", "The receiver reads but does not announce the end.", "Early close can break a later send."},
			}},
		},
	},
	{
		ID: "c111.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c111.turn", Kind: "markdown", Text: "Finish the receiver; run it and follow each report."}},
	},
}

var ch11PoolLessonArabic = []LessonStep{
	{
		ID: "c112.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c112.mission", Kind: "markdown",
			Text: "1. 3 workers يسحبوا أهداف من channel واحدة.\n2. كل واحد يرجّع ناتج حصاده.\n3. اجمع النواتج وابعت المجموع للقاعدة.\n\nالـstarter مجهّز worker والقنوات وWaitGroup؛ الـmain لسه فاضي.",
		}},
	},
	{
		ID: "c112.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c112.queue-trace", Kind: "walkthrough",
			Text: "الـworker pool طابور مشترك؛ العامل الفاضي ياخد المهمة الجاية.\n\n`range jobs` يخلص بعد القفل وتفريغ الرسائل؛ و`results` بتنقل ناتج الحصاد.\n\n`<-chan T` استقبال بس؛ `chan<- T` إرسال بس.\n\nالتتبع هنا لعامل واحد؛ الـpool يشغّل أكتر من عامل على القناة نفسها.",
			Code: ch11PoolQueueExample, Imports: []string{"fmt"}, ExampleID: "e11-pool-queue",
			Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل التشغيل."},
				{Line: 1, Caption: "طابور فاضي."},
				{Line: 2, Caption: "دخلت مهمة 3."},
				{Line: 3, Caption: "دخلت مهمة 5."},
				{Line: 4, Caption: "المرسل قفل الطابور."},
				{Line: 5, Caption: "العامل استلم 3."},
				{Line: 6, Caption: "الناتج 6.", Output: "6"},
				{Line: 5, Caption: "العامل استلم 5.", Output: "6"},
				{Line: 6, Caption: "الناتج 10.", Output: "6\n10"},
			}},
		}},
	},
	{
		ID: "c112.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID: "c112.predict", Kind: "question",
				Text: "الـworkerين بعتوا 2 و3؛ هتطبع كام لو وصلوا بأي ترتيب؟",
				Code: ch11PoolTryExample, Imports: []string{"fmt"}, ExampleID: "e11-pool-total",
				Question: &LessonQuestion{
					Options: []string{"5", "23", "3"}, Answer: 0,
					Feedback: []string{"الجمع يوصل لنفس الإجمالي بأي ترتيب.", "القيم بتتجمع كأرقام، مش بتتلزق كنص.", "النتيجتين بيتضافوا للإجمالي."},
				},
			},
			{
				ID: "c112.reveal", Kind: "walkthrough",
				Text: "ده ترتيب استقبال واحد ممكن؛ Go مش بيضمن مين يوصل الأول.",
				Code: ch11PoolTryExample, Imports: []string{"fmt"}, ExampleID: "e11-pool-total",
				Walkthrough: ch11PoolTryTrace([]string{
					"قبل التشغيل.", "قناة النتائج فاضية.", "الأول بعت 2.", "التاني بعت 3.",
					"استقبلت a قيمة 2.", "استقبلت b قيمة 3.", "الإجمالي اتطبع.",
				}),
			},
		},
	},
	{
		ID: "c112.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c112.watch", Kind: "markdown",
			Text: "- مرسل `jobs` يقفلها بعد آخر هدف.\n- استخدم WaitGroup قبل قفل `results`؛ من غير القفل، `range` هتفضل مستنية.",
		}},
	},
	{
		ID: "c112.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c112.plan", Kind: "markdown", Text: "رتّب تشغيل العمال، إرسال الأهداف، جمع النواتج، وإغلاق القنوات بعد اكتمال دورها."},
			{ID: "c112.q-close", Kind: "question", Text: "إمتى تقفل `results`؟", Question: &LessonQuestion{
				Options: []string{"بعد ما كل العمال يخلصوا", "بعد أول نتيجة", "قبل تشغيل العمال"}, Answer: 0,
				Feedback: []string{"استنى المجموعة كلها قبل ما تعلن النهاية.", "ممكن لسه فيه عامل هيبعت نتيجة.", "العمال لسه محتاجين القناة المفتوحة."},
			}},
			{ID: "c112.q-range", Kind: "question", Text: "إيه اللي ينهي `range jobs`؟", Question: &LessonQuestion{
				Options: []string{"القناة تتقفل بعد تفريغها", "أول عامل يسحب هدف", "انتهاء أول عامل"}, Answer: 0,
				Feedback: []string{"العمال يكملوا اللي في الطابور قبل ما يقفوا.", "باقي العمال يقدروا ياخدوا شغل كمان.", "عامل واحد خلص مش معناه إن الطابور خلص."},
			}},
		},
	},
	{
		ID: "c112.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c112.turn", Kind: "markdown", Text: "كمّل `main` وشغّلها؛ راقب استقبال الأهداف والنتائج لحد التقرير."}},
	},
}

var ch11PoolLessonEnglish = []LessonStep{
	{
		ID: "c112.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c112.mission", Kind: "markdown",
			Text: "1. Three workers share the target queue.\n2. Each sends its harvest count.\n3. Sum and report.\n\nThe starter defines `worker`; `main` is empty.",
		}},
	},
	{
		ID: "c112.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c112.queue-trace", Kind: "walkthrough",
			Text: "A pool shares jobs; a free worker takes the next.\n\n`range jobs` ends after close and draining; `results` carries harvest counts.\n\n`<-chan T` receives only; `chan<- T` sends only.\n\nThis trace shows one worker.",
			Code: ch11PoolQueueExample, Imports: []string{"fmt"}, ExampleID: "e11-pool-queue",
			Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Queue is empty."},
				{Line: 2, Caption: "Job 3 arrives."},
				{Line: 3, Caption: "Job 5 arrives."},
				{Line: 4, Caption: "Sender closes queue."},
				{Line: 5, Caption: "Worker gets 3."},
				{Line: 6, Caption: "Result is 6.", Output: "6"},
				{Line: 5, Caption: "Worker gets 5.", Output: "6"},
				{Line: 6, Caption: "Result is 10.", Output: "6\n10"},
			}},
		}},
	},
	{
		ID: "c112.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID: "c112.predict", Kind: "question",
				Text: "Workers send 2 and 3; what total prints in either receive order?",
				Code: ch11PoolTryExample, Imports: []string{"fmt"}, ExampleID: "e11-pool-total",
				Question: &LessonQuestion{
					Options: []string{"5", "23", "3"}, Answer: 0,
					Feedback: []string{"Addition gives the same total in either order.", "The values are added as numbers, not joined as text.", "Both results add to the total."},
				},
			},
			{
				ID: "c112.reveal", Kind: "walkthrough",
				Text: "One possible receive order; Go does not guarantee it.",
				Code: ch11PoolTryExample, Imports: []string{"fmt"}, ExampleID: "e11-pool-total",
				Walkthrough: ch11PoolTryTrace([]string{
					"Before run.", "Results empty.", "First sends 2.", "Second sends 3.",
					"a receives 2.", "b receives 3.", "Sum prints.",
				}),
			},
		},
	},
	{
		ID: "c112.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c112.watch", Kind: "markdown",
			Text: "- The `jobs` sender closes it after the last target.\n- Use WaitGroup before closing `results`; without close, `range` waits.",
		}},
	},
	{
		ID: "c112.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c112.plan", Kind: "markdown", Text: "Start workers, send targets, collect results, then close channels after work ends."},
			{ID: "c112.q-close", Kind: "question", Text: "When should `results` close?", Question: &LessonQuestion{
				Options: []string{"After every worker finishes", "After the first result", "Before workers start"}, Answer: 0,
				Feedback: []string{"Wait for every worker before signaling the end.", "A worker may still send a result.", "Workers need the channel open."},
			}},
			{ID: "c112.q-range", Kind: "question", Text: "What ends `range jobs`?", Question: &LessonQuestion{
				Options: []string{"The channel closes after draining", "The first worker takes a job", "The first worker finishes"}, Answer: 0,
				Feedback: []string{"Workers drain queued jobs before stopping.", "Other workers can take more jobs.", "One worker finishing does not empty it."},
			}},
		},
	},
	{
		ID: "c112.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c112.turn", Kind: "markdown", Text: "Finish `main` and run it; watch jobs and results arrive through the report."}},
	},
}

var ch11StormLessonArabic = []LessonStep{
	{
		ID: "c113.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c113.mission", Kind: "markdown",
			Text: "1. اجمع الكريستال المطلوب.\n2. عند الإنذار، ارجع للقاعدة قبل الضربة.\n3. بعد هدوء العاصفة، كمّل الشغل.\n\nالـstarter جهّز قناتي التحذير والهدوء؛ دورك تدير الرجوع والاستئناف.",
		}},
	},
	{
		ID: "c113.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c113.select-trace", Kind: "walkthrough",
			Text: "`select` شبه `switch` للقنوات: يستنى عملية جاهزة.\n\nمع `default`، لو مفيش تحذير جاهز، يكمل فورًا.\n\nقناة `colony.StormClear()` بتتقفل عند نهاية العاصفة؛ الاستقبال منها ساعتها يرجع فورًا.",
			Code: ch11StormSelectExample, Imports: []string{"fmt"}, ExampleID: "e11-storm-select",
			Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل التشغيل."},
				{Line: 1, Caption: "قناة التحذير فاضية."},
				{Line: 2, Caption: "أول select يفحص."},
				{Line: 5, Caption: "مفيش تحذير؛ default جاهزة."},
				{Line: 6, Caption: "الشغل مستمر.", Output: "work"},
				{Line: 8, Caption: "وصل تحذير للقناة.", Output: "work"},
				{Line: 9, Caption: "select يفحص تاني.", Output: "work"},
				{Line: 10, Caption: "وصلت قيمة التحذير.", Output: "work"},
				{Line: 11, Caption: "الرجوع اتطبع.", Output: "work\nreturn"},
			}},
		}},
	},
	{
		ID: "c113.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID: "c113.predict", Kind: "question",
				Text: "قناة `calm` اتقفلت؛ توقّع أي فرع يطبع.",
				Code: ch11StormTryExample, Imports: []string{"fmt"}, ExampleID: "e11-storm-clear",
				Question: &LessonQuestion{
					Options: []string{"safe", "work", "مفيش output"}, Answer: 0,
					Feedback: []string{"الاستقبال من قناة مقفولة جاهز فورًا.", "`default` تعمل بس لو مفيش استقبال جاهز.", "القناة المقفولة مش بتعلّق الاستقبال."},
				},
			},
			{
				ID: "c113.reveal", Kind: "walkthrough",
				Text: "القفل بيخلّي الاستقبال جاهز، فـ`select` يدخل الحالة دي.",
				Code: ch11StormTryExample, Imports: []string{"fmt"}, ExampleID: "e11-storm-clear",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "اتعملت قناة فاضية."},
					{Line: 2, Caption: "القناة اتقفلت."},
					{Line: 3, Caption: "select يفحص."},
					{Line: 4, Caption: "الاستقبال جاهز."},
					{Line: 5, Caption: "اتطبع safe.", Output: "safe"},
				}},
			},
		},
	},
	{
		ID: "c113.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c113.watch", Kind: "markdown",
			Text: "- `default` ما بتستناش؛ خلّي الحلقة تعمل شغل بعد الفحص.\n- حدّث channel الهدوء بعد كل إنذار عشان تستنى العاصفة الجاية.",
		}},
	},
	{
		ID: "c113.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c113.plan", Kind: "markdown", Text: "اشتغل بين الإنذارات، ارجع وقت التحذير، استنى الهدوء، وبعدها استأنف."},
			{ID: "c113.q-default", Kind: "question", Text: "لو `warning` لسه فاضية، `default` تعمل إيه؟", Question: &LessonQuestion{
				Options: []string{"تتنفذ فورًا", "تستنى التحذير", "تقفل القناة"}, Answer: 0,
				Feedback: []string{"مفيش عملية جاهزة، فـ`default` تكمل فورًا.", "من غير `default`، `select` هيستنى.", "`default` مش بتغيّر حالة القناة."},
			}},
			{ID: "c113.q-clear", Kind: "question", Text: "إمتى `<-calm` يرجع بعد القفل؟", Question: &LessonQuestion{
				Options: []string{"فورًا", "بعد رسالة جديدة", "بعد ما يتقفل تاني"}, Answer: 0,
				Feedback: []string{"الاستقبال من قناة مقفولة جاهز دايمًا.", "القفل نفسه علامة انتهاء الإشارة.", "القناة المقفولة مش محتاجة قفل تاني."},
			}},
		},
	},
	{
		ID: "c113.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c113.turn", Kind: "markdown", Text: "كمّل مسار الرجوع والاستئناف وشغّل المهمة؛ تابع كل إنذار."}},
	},
}

var ch11StormLessonEnglish = []LessonStep{
	{
		ID: "c113.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c113.mission", Kind: "markdown",
			Text: "1. Collect the required crystals.\n2. Return to base when warned.\n3. Continue after the storm clears.\n\nThe starter has both channels; manage the return and resume.",
		}},
	},
	{
		ID: "c113.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c113.select-trace", Kind: "walkthrough",
			Text: "`select` is like `switch` for channels: it waits for a ready operation.\n\nWith `default`, it continues at once when no warning is ready.\n\nThe `colony.StormClear()` channel closes when the storm ends; receiving then returns immediately.",
			Code: ch11StormSelectExample, Imports: []string{"fmt"}, ExampleID: "e11-storm-select",
			Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Warning channel is empty."},
				{Line: 2, Caption: "First select checks."},
				{Line: 5, Caption: "No warning; default is ready."},
				{Line: 6, Caption: "Work continues.", Output: "work"},
				{Line: 8, Caption: "Warning enters channel.", Output: "work"},
				{Line: 9, Caption: "Select checks again.", Output: "work"},
				{Line: 10, Caption: "Warning is received.", Output: "work"},
				{Line: 11, Caption: "Return prints.", Output: "work\nreturn"},
			}},
		}},
	},
	{
		ID: "c113.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID: "c113.predict", Kind: "question",
				Text: "The `calm` channel is closed; predict which branch prints.",
				Code: ch11StormTryExample, Imports: []string{"fmt"}, ExampleID: "e11-storm-clear",
				Question: &LessonQuestion{
					Options: []string{"safe", "work", "No output"}, Answer: 0,
					Feedback: []string{"A receive from a closed channel is ready immediately.", "`default` runs only when no receive is ready.", "A closed channel does not block the receive."},
				},
			},
			{
				ID: "c113.reveal", Kind: "walkthrough",
				Text: "Closing makes the receive ready, so `select` takes that case.",
				Code: ch11StormTryExample, Imports: []string{"fmt"}, ExampleID: "e11-storm-clear",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Empty channel created."},
					{Line: 2, Caption: "Channel closes."},
					{Line: 3, Caption: "Select checks."},
					{Line: 4, Caption: "Receive is ready."},
					{Line: 5, Caption: "safe prints.", Output: "safe"},
				}},
			},
		},
	},
	{
		ID: "c113.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c113.watch", Kind: "markdown",
			Text: "- `default` does not wait; make the loop do work after checking.\n- Refresh the calm channel after each warning to await another storm.",
		}},
	},
	{
		ID: "c113.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c113.plan", Kind: "markdown", Text: "Work between warnings, return when warned, wait for calm, then resume."},
			{ID: "c113.q-default", Kind: "question", Text: "If `warning` is empty, what does `default` do?", Question: &LessonQuestion{
				Options: []string{"Runs immediately", "Waits for warning", "Closes the channel"}, Answer: 0,
				Feedback: []string{"No operation is ready, so `default` runs now.", "Without `default`, `select` would wait.", "`default` does not change channel state."},
			}},
			{ID: "c113.q-clear", Kind: "question", Text: "After close, when does `<-calm` return?", Question: &LessonQuestion{
				Options: []string{"Immediately", "After a new message", "After another close"}, Answer: 0,
				Feedback: []string{"A receive from a closed channel is always ready.", "Close itself signals the end.", "A closed channel does not need another close."},
			}},
		},
	},
	{
		ID: "c113.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c113.turn", Kind: "markdown", Text: "Finish the return and resume path and run the mission; follow each warning."}},
	},
}
