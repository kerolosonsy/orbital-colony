package levels

const ch12CounterExample = "var mu sync.Mutex\n" +
	"total := 0\n" +
	"var wg sync.WaitGroup\n" +
	"add := func(n int) {\n" +
	"\tdefer wg.Done()\n" +
	"\tmu.Lock()\n" +
	"\tdefer mu.Unlock()\n" +
	"\ttotal += n\n" +
	"}\n" +
	"wg.Add(2)\n" +
	"go add(2)\n" +
	"go add(3)\n" +
	"wg.Wait()\n" +
	"fmt.Println(total)"

const ch12RegistryExample = "sizes := map[string]int{}\n" +
	"var mu sync.Mutex\n" +
	"var wg sync.WaitGroup\n" +
	"record := func(key string, size int) {\n" +
	"\tmu.Lock()\n" +
	"\tdefer mu.Unlock()\n" +
	"\tsizes[key] = size\n" +
	"}\n" +
	"wg.Add(2)\n" +
	"go func() { defer wg.Done(); record(\"amber\", 4) }()\n" +
	"go func() { defer wg.Done(); record(\"violet\", 7) }()\n" +
	"wg.Wait()\n" +
	"fmt.Println(len(sizes), sizes[\"violet\"])"

const ch12DeadlockExample = "var first, second sync.Mutex\n" +
	"var wg sync.WaitGroup\n" +
	"done := 0\n" +
	"work := func() {\n" +
	"\tfirst.Lock()\n" +
	"\tsecond.Lock()\n" +
	"\tdone++\n" +
	"\tsecond.Unlock()\n" +
	"\tfirst.Unlock()\n" +
	"}\n" +
	"wg.Add(2)\n" +
	"go func() { defer wg.Done(); work() }()\n" +
	"go func() { defer wg.Done(); work() }()\n" +
	"wg.Wait()\n" +
	"fmt.Println(done)"

func ch12CounterTrace(captions []string) *LessonWalkthrough {
	return &LessonWalkthrough{Frames: []LessonFrame{
		{Line: 0, Caption: captions[0]},
		{Line: 2, Caption: captions[1], Variables: []LessonVariable{{Name: "total", Value: "0"}}},
		{Line: 6, Caption: captions[2], Variables: []LessonVariable{{Name: "total", Value: "0"}, {Name: "n", Value: "2"}}},
		{Line: 7, Caption: captions[3], Variables: []LessonVariable{{Name: "total", Value: "0"}, {Name: "n", Value: "2"}}},
		{Line: 8, Caption: captions[4], Variables: []LessonVariable{{Name: "total", Value: "2"}, {Name: "n", Value: "2"}}},
		{Line: 9, Caption: captions[5], Variables: []LessonVariable{{Name: "total", Value: "2"}, {Name: "n", Value: "2"}}},
		{Line: 6, Caption: captions[6], Variables: []LessonVariable{{Name: "total", Value: "2"}, {Name: "n", Value: "3"}}},
		{Line: 8, Caption: captions[7], Variables: []LessonVariable{{Name: "total", Value: "5"}, {Name: "n", Value: "3"}}},
		{Line: 13, Caption: captions[8], Variables: []LessonVariable{{Name: "total", Value: "5"}}},
		{Line: 14, Caption: captions[9], Variables: []LessonVariable{{Name: "total", Value: "5"}}, Output: "5"},
	}}
}

func ch12RegistryTrace(captions []string) *LessonWalkthrough {
	return &LessonWalkthrough{Frames: []LessonFrame{
		{Line: 0, Caption: captions[0]},
		{Line: 1, Caption: captions[1], Variables: []LessonVariable{{Name: "sizes", Value: "{}"}}},
		{Line: 5, Caption: captions[2], Variables: []LessonVariable{{Name: "sizes", Value: "{}"}}},
		{Line: 7, Caption: captions[3], Variables: []LessonVariable{{Name: "sizes", Value: "{amber: 4}"}}},
		{Line: 5, Caption: captions[4], Variables: []LessonVariable{{Name: "sizes", Value: "{amber: 4}"}}},
		{Line: 7, Caption: captions[5], Variables: []LessonVariable{{Name: "sizes", Value: "{amber: 4, violet: 7}"}}},
		{Line: 12, Caption: captions[6], Variables: []LessonVariable{{Name: "sizes", Value: "{amber: 4, violet: 7}"}}},
		{Line: 13, Caption: captions[7], Variables: []LessonVariable{{Name: "sizes", Value: "{amber: 4, violet: 7}"}}, Output: "2 7"},
	}}
}

func ch12DeadlockTrace(captions []string) *LessonWalkthrough {
	return &LessonWalkthrough{Frames: []LessonFrame{
		{Line: 0, Caption: captions[0]},
		{Line: 3, Caption: captions[1], Variables: []LessonVariable{{Name: "done", Value: "0"}}},
		{Line: 5, Caption: captions[2], Variables: []LessonVariable{{Name: "done", Value: "0"}}},
		{Line: 6, Caption: captions[3], Variables: []LessonVariable{{Name: "done", Value: "0"}}},
		{Line: 7, Caption: captions[4], Variables: []LessonVariable{{Name: "done", Value: "1"}}},
		{Line: 8, Caption: captions[5], Variables: []LessonVariable{{Name: "done", Value: "1"}}},
		{Line: 5, Caption: captions[6], Variables: []LessonVariable{{Name: "done", Value: "1"}}},
		{Line: 7, Caption: captions[7], Variables: []LessonVariable{{Name: "done", Value: "2"}}},
		{Line: 15, Caption: captions[8], Variables: []LessonVariable{{Name: "done", Value: "2"}}, Output: "2"},
	}}
}

var ch12Lesson121Arabic = []LessonStep{
	{
		ID: "c121.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c121.mission", Kind: "markdown",
			Text: "1. خلّي 4 درونات تزود نفس المتغير.\n2. احمِ التحديث من data race.\n3. ابعت الإجمالي بعد انتهاء الشغل.\n\nالـstarter بيحصد بالفعل؛ ركّز على الكتابة المشتركة.",
		}},
	},
	{
		ID: "c121.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c121.idea", Kind: "markdown",
			Text: "تخيّل ~total~ لوحة واحدة. كل تحديث يقرأها، يزوّدها، ثم يكتبها.\n\nلو التحديثات اتداخلت، ممكن قيمة تضيع؛ النتيجة الصح بالصدفة مش ضمان.\n\nMutex مشترك يسمح لواحد بس يدخل الجزء المحمي.\n~WaitGroup~ يستنى الشغل، لكنه مش بيحمي المتغير.",
		}},
	},
	{
		ID: "c121.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID: "c121.predict", Kind: "question",
				Text: "المثال يضيف 2 و3 تحت قفل واحد. إيه اللي يطبع؟",
				Question: &LessonQuestion{
					Options: []string{"5", "2", "النتيجة ممكن تتغير"}, Answer: 0,
					Feedback: []string{
						"التحديثين بيتجمعوا؛ ترتيبهم مش بيغيّر المجموع.",
						"ده أول تحديث لوحده؛ العامل التاني كمان بيضيف.",
						"الترتيب مش مضمون، لكن كل تحديث محمي بالقفل.",
					},
				},
			},
			{
				ID: "c121.counter-trace", Kind: "walkthrough",
				Text: "ده ترتيب واحد ممكن: 2 يدخل الأول. Go مش بتضمن مين يمسك القفل أولاً.",
				Code: ch12CounterExample, Imports: []string{"fmt", "sync"}, ExampleID: "e121-counter",
				Walkthrough: ch12CounterTrace([]string{
					"قبل التشغيل.", "الإجمالي يبدأ صفر.", "ترتيب ممكن: n=2 يمسك القفل.", "defer هيفتحه عند الخروج.", "الإجمالي صار 2.", "الرجوع يفتح القفل.", "بعده n=3 يدخل.", "الإجمالي صار 5.", "العمال خلصوا.", "يطبع 5.",
				}),
			},
		},
	},
	{
		ID: "c121.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c121.watch", Kind: "markdown",
			Text: "القفل لازم يبقى هو نفسه عند كل العاملين.\n\nخلي الجزء بين ~Lock~ و~Unlock~ صغير؛ حرّره حتى لو الدالة خرجت.\n\nRace detector ممكن يكشف السباق حتى لو الإجمالي ظهر صح.",
		}},
	},
	{
		ID: "c121.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c121.plan", Kind: "markdown", Text: "الخطة: شغّل العمال → احمِ التحديث → استناهم → ابعت الإجمالي."},
			{
				ID: "c121.wait-question", Kind: "question", Text: "هل ~wg.Wait()~ بيحمي total أثناء تحديثه؟",
				Question: &LessonQuestion{Options: []string{"أيوه", "لأ", "بس لو عامل واحد"}, Answer: 1, Feedback: []string{
					"الانتظار مش قفل؛ ممكن العمال يكتبوا في نفس الوقت.",
					"صح، Wait يستنى النهاية بس؛ القفل يحمي الكتابة.",
					"عدد العمال مش بيحوّل الانتظار لقفل.",
				}},
			},
			{
				ID: "c121.lock-question", Kind: "question", Text: "أنهي اختيار يحمي كل تحديثات الإجمالي؟",
				Question: &LessonQuestion{Options: []string{"نفس القفل المشترك", "قفل جديد لكل goroutine", "القفل بعد Wait"}, Answer: 0, Feedback: []string{
					"كل كاتب لازم يتفق على نفس القفل.",
					"أقفال مختلفة مش بتنسّق على نفس المتغير.",
					"القفل بعد انتهاء الكتابة متأخر.",
				}},
			},
		},
	},
	{
		ID: "c121.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c121.turn", Kind: "markdown", Text: "راجع سطر الجمع في الـstarter، ثم شغّل واختبر race detector."}},
	},
}

var ch12Lesson121English = []LessonStep{
	{
		ID: "c121.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c121.mission", Kind: "markdown",
			Text: "1. Make 4 drones add to one total.\n2. Protect updates from a data race.\n3. Report the result.\n\nThe starter already harvests.",
		}},
	},
	{
		ID: "c121.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c121.idea", Kind: "markdown",
			Text: "Think of ~total~ as one board. An update reads, adds, then writes.\n\nOverlapping updates can lose a value. A correct result by luck proves nothing.\n\nOne shared mutex protects this short update. ~WaitGroup~ waits; it does not protect total.",
		}},
	},
	{
		ID: "c121.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID: "c121.predict", Kind: "question",
				Text: "What does this locked example print?",
				Question: &LessonQuestion{
					Options: []string{"5", "2", "It may vary"}, Answer: 0,
					Feedback: []string{
						"Both additions remain.",
						"That is only the first.",
						"Order may vary; the lock protects both.",
					},
				},
			},
			{
				ID: "c121.counter-trace", Kind: "walkthrough",
				Text: "One possible order: 2 enters first. Go does not promise who gets the lock first.",
				Code: ch12CounterExample, Imports: []string{"fmt", "sync"}, ExampleID: "e121-counter",
				Walkthrough: ch12CounterTrace([]string{
					"Before execution.", "Total starts at zero.", "One order: n=2 gets the lock.", "defer unlocks on return.", "Total becomes 2.", "Return releases the lock.", "Then n=3 enters.", "Total becomes 5.", "Workers have finished.", "Prints 5.",
				}),
			},
		},
	},
	{
		ID: "c121.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c121.watch", Kind: "markdown",
			Text: "All workers must use the same mutex.\n\nKeep the locked section small; defer its release.\n\nThe race detector can catch a race despite a correct total.",
		}},
	},
	{
		ID: "c121.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c121.plan", Kind: "markdown", Text: "Plan: start workers → protect the update → wait → report."},
			{
				ID: "c121.wait-question", Kind: "question", Text: "Does Wait protect total during updates?",
				Question: &LessonQuestion{Options: []string{"Yes", "No", "Only with one worker"}, Answer: 1, Feedback: []string{
					"Wait only waits; it does not lock.",
					"Right; the mutex protects writes.",
					"Worker count does not make Wait a lock.",
				}},
			},
			{
				ID: "c121.lock-question", Kind: "question", Text: "What protects every total update?",
				Question: &LessonQuestion{Options: []string{"One shared mutex", "One lock per goroutine", "Lock after Wait"}, Answer: 0, Feedback: []string{
					"All writers need the same lock.",
					"Different locks do not coordinate.",
					"That lock comes too late.",
				}},
			},
		},
	},
	{
		ID: "c121.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c121.turn", Kind: "markdown", Text: "Check the starter's addition line, then run it with the race detector."}},
	},
}

var ch12Lesson122Arabic = []LessonStep{
	{
		ID: "c122.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c122.mission", Kind: "markdown",
			Text: "1. خلّي 3 درونات تسجل القياسات في map واحدة.\n2. احمِ السجل من الكتابة المتزامنة.\n3. ابعت مكان أكبر كريستال.\n\nالـstarter بيقيس الصفوف ويسجل المواقع بالفعل.",
		}},
	},
	{
		ID: "c122.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c122.idea", Kind: "markdown",
			Text: "الـmap المشتركة مش آمنة مع كتابات متزامنة.\n\nالنمط المعتاد: حط ~sync.Mutex~ جنب الـmap في ~Registry~.\nكل method بتوصل للبيانات تقفل نفس القفل؛ receiver من نوع ~*Registry~ ما ينسخش القفل.\n\nبعد ~wg.Wait()~، الكتابة خلصت؛ بعدها تقدر تقرأ السجل.",
		}},
	},
	{
		ID: "c122.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID: "c122.predict", Kind: "question",
				Text: "السجل التجريبي يضيف مفتاحين. إيه اللي يطبع؟",
				Question: &LessonQuestion{Options: []string{"2 7", "2 4", "7 2"}, Answer: 0, Feedback: []string{
					"فيه مفتاحين، وقيمة violet هي 7.",
					"المفتاحين موجودين؛ راجع القيمة المرتبطة بـviolet.",
					"البرنامج يطبع عدد المفاتيح قبل قيمة violet.",
				}},
			},
			{
				ID: "c122.registry-trace", Kind: "walkthrough",
				Text: "مفاتيح توضيحية، مش بيانات المهمة. ده ترتيب واحد ممكن للكتابتين.",
				Code: ch12RegistryExample, Imports: []string{"fmt", "sync"}, ExampleID: "e122-registry",
				Walkthrough: ch12RegistryTrace([]string{
					"قبل التشغيل.", "السجل فاضي.", "عامل amber مسك القفل.", "اتسجلت قيمته.", "violet مسك القفل بعده.", "بعده اتسجل violet.", "كل الكتّاب خلصوا.", "يطبع 2 7.",
				}),
			},
		},
	},
	{
		ID: "c122.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c122.watch", Kind: "markdown",
			Text: "كل وصول متزامن لنفس الـmap يستخدم نفس القفل.\n\n~Record~ يقفل الكتابة؛ method قراءة تحتاج حماية لو العمال لسه شغالين.\n\nبعد ~Wait~، مفيش عامل بيكتب، فاقرأ بعدها.",
		}},
	},
	{
		ID: "c122.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c122.plan", Kind: "markdown", Text: "الخطة: سجل بأمان → استنى الكل → دور على الموقع → ابعته."},
			{
				ID: "c122.lock-question", Kind: "question", Text: "لو كل goroutine معاها قفل لوحدها، هل تحمي map واحدة؟",
				Question: &LessonQuestion{Options: []string{"أيوه", "لأ", "بس وقت القراءة"}, Answer: 1, Feedback: []string{
					"الأقفال المختلفة مش بتنسّق على نفس البيانات.",
					"صح، كل الكتاب لازم يستخدموا mutex مشتركة.",
					"قفل القراءة لوحده مش يحمي كتابات الأقفال التانية.",
				}},
			},
			{
				ID: "c122.wait-question", Kind: "question", Text: "بعد ما ~wg.Wait()~ يرجع، هل القراءة آمنة؟",
				Question: &LessonQuestion{Options: []string{"أيوه، لو مفيش كاتب تاني", "لأ، map لا تُقرأ", "بس بقفل مختلف"}, Answer: 0, Feedback: []string{
					"صح، الانتظار خلّص كل الكتّاب؛ مفيش كتابة متزامنة.",
					"ممكن تقرأ بعد انتهاء الكتابة المتزامنة.",
					"قفل مختلف مش مطلوب لما الكتّاب خلصوا.",
				}},
			},
		},
	},
	{
		ID: "c122.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c122.turn", Kind: "markdown", Text: "أمّن تسجيل الـstarter، ثم أضف قراءة للموقع بعد انتهاء العمال."}},
	},
}

var ch12Lesson122English = []LessonStep{
	{
		ID: "c122.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c122.mission", Kind: "markdown",
			Text: "1. Have 3 drones record in one map.\n2. Protect concurrent writes.\n3. Report the largest crystal's location.\n\nThe starter already measures and records locations.",
		}},
	},
	{
		ID: "c122.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c122.idea", Kind: "markdown",
			Text: "Concurrent writes to one map need protection.\n\nKeep ~sync.Mutex~ beside the map in ~Registry~.\nEach method uses that lock; a ~*Registry~ receiver does not copy it.\n\nAfter ~wg.Wait()~, writes finish; then read the log.",
		}},
	},
	{
		ID: "c122.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID: "c122.predict", Kind: "question",
				Text: "Two keys are added. What prints?",
				Question: &LessonQuestion{Options: []string{"2 7", "2 4", "7 2"}, Answer: 0, Feedback: []string{
					"Two keys; violet has value 7.",
					"Check the value stored for violet.",
					"The key count prints first.",
				}},
			},
			{
				ID: "c122.registry-trace", Kind: "walkthrough",
				Text: "Illustrative keys, not mission data. This is one possible write order.",
				Code: ch12RegistryExample, Imports: []string{"fmt", "sync"}, ExampleID: "e122-registry",
				Walkthrough: ch12RegistryTrace([]string{
					"Before execution.", "The log is empty.", "Amber's worker gets the lock.", "Its value is recorded.", "Violet gets the lock next.", "Violet is recorded next.", "All writers finish.", "Prints 2 7.",
				}),
			},
		},
	},
	{
		ID: "c122.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c122.watch", Kind: "markdown",
			Text: "Concurrent map access uses one shared lock.\n\n~Record~ locks writes; protect reads while workers run.\n\nAfter ~Wait~, the writers have finished.",
		}},
	},
	{
		ID: "c122.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c122.plan", Kind: "markdown", Text: "Plan: record safely → wait → find the location → report."},
			{
				ID: "c122.lock-question", Kind: "question", Text: "Do separate goroutine locks protect one map?",
				Question: &LessonQuestion{Options: []string{"Yes", "No", "Only while reading"}, Answer: 1, Feedback: []string{
					"Different locks do not coordinate.",
					"Right; every writer shares one mutex.",
					"A read lock alone does not protect writes.",
				}},
			},
			{
				ID: "c122.wait-question", Kind: "question", Text: "After Wait returns, can you read safely?",
				Question: &LessonQuestion{Options: []string{"Yes, if no writer remains", "No, never", "Only with a new lock"}, Answer: 0, Feedback: []string{
					"Right; every writer has finished.",
					"Maps are readable after concurrent writes end.",
					"No writer remains to coordinate with.",
				}},
			},
		},
	},
	{
		ID: "c122.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c122.turn", Kind: "markdown", Text: "Protect the starter's recording, then read a location after the workers finish."}},
	},
}

var ch12Lesson123Arabic = []LessonStep{
	{
		ID: "c123.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c123.mission", Kind: "markdown",
			Text: "1. رجّع حركة الدرونات اللي وقفت.\n2. حلّ الـdeadlock بين قفلين.\n3. احصد الصفين.\n\nالـstarter فيه عاملين بيطلبوا المثقاب والعربية بترتيبين مختلفين.",
		}},
	},
	{
		ID: "c123.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c123.idea", Kind: "markdown",
			Text: "Deadlock دايرة انتظار: goroutine ماسكة قفل ومستنية قفل تاني.\n\nلو عاملين أخدوا القفلين بترتيب عكسي، ممكن كل واحد يستنى التاني.\n\nGo مش بتضمن مين يبدأ؛ ده سيناريو محتمل، مش نتيجة كل تشغيل.\n\nالقاعدة: اطلب الموارد بترتيب ثابت في كل المسارات.",
		}},
	},
	{
		ID: "c123.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{
				ID: "c123.predict", Kind: "question",
				Text: "العاملين بيطلبوا القفلين بنفس الترتيب. إيه اللي يطبع؟",
				Question: &LessonQuestion{Options: []string{"2", "1", "البرنامج يعلّق"}, Answer: 0, Feedback: []string{
					"العاملين بينفذوا الزيادة؛ القفلين مش بيعملوا دورة انتظار.",
					"كل عامل بيزود done مرة واحدة.",
					"الترتيب موحد، فمش هتتكوّن دورة بين القفلين.",
				}},
			},
			{
				ID: "c123.lock-trace", Kind: "walkthrough",
				Text: "ترتيب واحد ممكن: العامل الأول يخلص ثم الثاني. الناتج ثابت، والجدولة ممكن تختلف.",
				Code: ch12DeadlockExample, Imports: []string{"fmt", "sync"}, ExampleID: "e123-lock-order",
				Walkthrough: ch12DeadlockTrace([]string{
					"قبل التشغيل.", "العداد يبدأ صفر.", "ترتيب ممكن: الأول يأخذ القفل.", "بعده يأخذ القفل التاني.", "العداد صار 1.", "يفتح القفل التاني.", "العامل التاني يدخل بعده.", "العداد صار 2.", "يطبع 2.",
				}),
			},
		},
	},
	{
		ID: "c123.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c123.watch", Kind: "markdown",
			Text: "~WaitGroup~ مش بيفك deadlock؛ هو بس يستنى النهاية.\n\nكل المسارات لازم تطلب الأقفال بنفس الترتيب.\n\nافتح الأقفال بعد الشغل المحمي وبالترتيب العكسي.",
		}},
	},
	{
		ID: "c123.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c123.plan", Kind: "markdown", Text: "الخطة: حدد مين ماسك كل قفل → امنع دورة الانتظار → اختبر حصاد الصفين."},
			{
				ID: "c123.order-question", Kind: "question", Text: "إيه اللي يمنع دورة انتظار بين قفلين؟",
				Question: &LessonQuestion{Options: []string{"ترتيب ثابت واحد", "ترتيب مختلف لكل عامل", "WaitGroup إضافية"}, Answer: 0, Feedback: []string{
					"الترتيب الموحد يمنع اختلاف ترتيب الطلب.",
					"الاختلاف ممكن يعمل دورة انتظار.",
					"WaitGroup تستنى؛ مش بتغيّر ترتيب الأقفال.",
				}},
			},
			{
				ID: "c123.schedule-question", Kind: "question", Text: "هل Go تضمن مين يبدأ في الـwalkthrough؟",
				Question: &LessonQuestion{Options: []string{"لأ، ده ترتيب واحد ممكن", "أيوه، العامل الأول دايماً", "WaitGroup بترتبهم"}, Answer: 0, Feedback: []string{
					"صح، goroutines ممكن تتجدول بترتيبات مختلفة.",
					"ترتيب البدء مش مضمون من Go.",
					"WaitGroup تستنى الاكتمال، مش بتختار مين يبدأ.",
				}},
			},
		},
	},
	{
		ID: "c123.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c123.turn", Kind: "markdown", Text: "عدّل الـstarter وشغّله؛ راقب حركة الدرونين وإتمام الصفين."}},
	},
}

var ch12Lesson123English = []LessonStep{
	{
		ID: "c123.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c123.mission", Kind: "markdown",
			Text: "1. Get the stopped drones moving.\n2. Fix a two-lock deadlock.\n3. Harvest both rows.\n\nThe workers request locks in different orders.",
		}},
	},
	{
		ID: "c123.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c123.idea", Kind: "markdown",
			Text: "A deadlock is a wait cycle: a goroutine holds one lock and waits for another.\n\nIf workers acquire locks in reverse orders, each may wait for the other.\n\nGo does not promise who starts; this is one possible order.\n\nRequest resources in one fixed order on every path.",
		}},
	},
	{
		ID: "c123.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{
				ID: "c123.predict", Kind: "question",
				Text: "Both workers use one lock order. What prints?",
				Question: &LessonQuestion{Options: []string{"2", "1", "The program hangs"}, Answer: 0, Feedback: []string{
					"Both workers increment; there is no wait cycle.",
					"Each worker increments once.",
					"One order prevents a lock cycle.",
				}},
			},
			{
				ID: "c123.lock-trace", Kind: "walkthrough",
				Text: "One possible order: worker one finishes, then two. Output is fixed; scheduling may vary.",
				Code: ch12DeadlockExample, Imports: []string{"fmt", "sync"}, ExampleID: "e123-lock-order",
				Walkthrough: ch12DeadlockTrace([]string{
					"Before execution.", "Counter starts at zero.", "One order: worker one locks.", "Then it gets the second lock.", "Counter becomes 1.", "It releases the second lock.", "Worker two enters next.", "Counter becomes 2.", "Prints 2.",
				}),
			},
		},
	},
	{
		ID: "c123.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c123.watch", Kind: "markdown",
			Text: "~WaitGroup~ cannot break a deadlock; it only waits.\n\nEvery path needs one lock order.\n\nUnlock after protected work, in reverse order.",
		}},
	},
	{
		ID: "c123.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c123.plan", Kind: "markdown", Text: "Plan: find held locks → prevent the cycle → test both rows."},
			{
				ID: "c123.order-question", Kind: "question", Text: "What prevents a two-lock cycle?",
				Question: &LessonQuestion{Options: []string{"One fixed order", "A different order per worker", "Another WaitGroup"}, Answer: 0, Feedback: []string{
					"One shared order prevents conflicting requests.",
					"Different orders can create a wait cycle.",
					"WaitGroup waits; it does not change lock order.",
				}},
			},
			{
				ID: "c123.schedule-question", Kind: "question", Text: "Does Go promise who starts?",
				Question: &LessonQuestion{Options: []string{"No, this is one possible order", "Yes, worker one always starts", "WaitGroup chooses the order"}, Answer: 0, Feedback: []string{
					"Right, goroutine scheduling can vary.",
					"Go does not promise a start order.",
					"WaitGroup waits; it does not choose who starts.",
				}},
			},
		},
	},
	{
		ID: "c123.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c123.turn", Kind: "markdown", Text: "Edit the starter, then watch both drones finish the rows."}},
	},
}
