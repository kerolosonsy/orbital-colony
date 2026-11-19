package levels

func ch19Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch19Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch19Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

var ch19DurationsArabic = []LessonStep{
	{ID: "c191.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch19Markdown("c191.mission", "اقرأ جدول الرحلات واحسب وصول كل رحلة.\n\nخلي ترتيب التقارير زي ترتيب الأوامر، حتى لو أمر واحد فيه خطأ.")}},
	{ID: "c191.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch19Walkthrough("c191.idea-duration", "~time.Duration~ مدة منقضية، و~ParseDuration~ بتحوّل النص لقيمة زمنية.", ch19DurationParseCode, "e19-1-duration", []string{"fmt", "time"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "اقرأ مدة توضيحية."},
		{Line: 3, Caption: "اتأكد من خطأ parsing."},
		{Line: 7, Caption: "اعرض المدة بصيغتها.", Output: "35m0s 35m0s"},
	})}},
	{ID: "c191.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch19Question("c191.try-arrival", "توقّع الناتج:\n\n~~~go\n"+ch19DurationAddCode+"\n~~~", []string{"23:40", "00:15", "35m"}, 1, []string{"ده وقت البداية قبل إضافة المدة.", "الـminute بتعدّي لليوم التالي تلقائياً.", "Format هنا بيعرض الساعة، مش المدة."}),
		ch19Walkthrough("c191.try-reveal", "تابع وقت توضيحي وهو يعدّي نص الليل.", ch19DurationAddCode, "e19-1-add", []string{"fmt", "time"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "وقت البداية قيمة مستقلة.", Variables: []LessonVariable{{Name: "start", Value: "2031-04-03 23:40 UTC"}}},
			{Line: 3, Caption: "Add رجّعت وقت وصول جديد.", Variables: []LessonVariable{{Name: "arrival", Value: "2031-04-04 00:15 UTC"}}},
			{Line: 4, Caption: "اعرض الساعة الجديدة.", Variables: []LessonVariable{{Name: "arrival", Value: "2031-04-04 00:15 UTC"}}, Output: "00:15"},
		}),
	}},
	{ID: "c191.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch19Markdown("c191.watch", "~Parse~ بتستخدم layout Go من تاريخ مرجعي، مش رموز YYYY/MM.\n\nالـduration ممكن تتقري صح وتفضل مخالفة للعقد؛ افصل parsing عن التحقق.")}},
	{ID: "c191.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch19Markdown("c191.plan", "اقرأ كل field لوحده، وبعدها قرر صلاحية الأمر.\n\nخلّي نتيجة كل أمر في مكانه وما تستخدمش ساعة الجهاز."),
		ch19Question("c191.check-parse", "لو ~ParseDuration~ نجحت، هل ده يثبت إن الأمر صالح؟", []string{"لأ، لسه فيه تحقق للعقد", "آه، parsing كفاية", "بس لو المدة ساعة"}, 0, []string{"الـparser بتأكد الصيغة، مش قواعد المهمة.", "المدخل ممكن يخالف العقد رغم صحة الصيغة.", "نوع المدة مش معيار التحقق."}),
	}},
	{ID: "c191.turn-step", Title: "دورك", Blocks: []LessonBlock{ch19Markdown("c191.turn", "ابدأ بمدخل واحد للتوضيح، شغّله، وبعدها راجع ترتيب كل تقرير.")}},
}

var ch19DurationsEnglish = []LessonStep{
	{ID: "c191.mission-step", Title: "Mission", Blocks: []LessonBlock{ch19Markdown("c191.mission", "Read the flight schedule and calculate each arrival.\n\nKeep reports in order, even when one order is invalid.")}},
	{ID: "c191.idea-step", Title: "Idea", Blocks: []LessonBlock{ch19Walkthrough("c191.idea-duration", "~time.Duration~ represents elapsed time; ~ParseDuration~ converts text to a duration.", ch19DurationParseCode, "e19-1-duration", []string{"fmt", "time"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Parse an illustrative duration."},
		{Line: 3, Caption: "Check the parse error."},
		{Line: 7, Caption: "Print the canonical duration.", Output: "35m0s 35m0s"},
	})}},
	{ID: "c191.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch19Question("c191.try-arrival", "Predict the output:\n\n~~~go\n"+ch19DurationAddCode+"\n~~~", []string{"23:40", "00:15", "35m"}, 1, []string{"That is the start time before adding duration.", "The minutes carry into the next day automatically.", "Format displays the clock time, not the duration."}),
		ch19Walkthrough("c191.try-reveal", "Follow an illustrative time across midnight.", ch19DurationAddCode, "e19-1-add", []string{"fmt", "time"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The start is its own value.", Variables: []LessonVariable{{Name: "start", Value: "2031-04-03 23:40 UTC"}}},
			{Line: 3, Caption: "Add returns a new arrival.", Variables: []LessonVariable{{Name: "arrival", Value: "2031-04-04 00:15 UTC"}}},
			{Line: 4, Caption: "Print the new clock time.", Variables: []LessonVariable{{Name: "arrival", Value: "2031-04-04 00:15 UTC"}}, Output: "00:15"},
		}),
	}},
	{ID: "c191.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch19Markdown("c191.watch", "~Parse~ uses a Go layout built from a reference date, not YYYY/MM tokens.\n\nA duration can parse yet violate the contract; separate parsing from validation.")}},
	{ID: "c191.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch19Markdown("c191.plan", "Read each field separately, then decide whether the order is valid.\n\nKeep each result in its original place and avoid the host clock."),
		ch19Question("c191.check-parse", "If ~ParseDuration~ succeeds, is the order valid?", []string{"Not yet; check the contract", "Yes, parsing is enough", "Only when it is one hour"}, 0, []string{"Parsing checks syntax, not mission rules.", "Input can parse and still violate the contract.", "Duration type does not define validity."}),
	}},
	{ID: "c191.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch19Markdown("c191.turn", "Start with one illustrative order, run it, then review the order of every report.")}},
}

var ch19SunsArabic = []LessonStep{
	{ID: "c192.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch19Markdown("c192.mission", "اعرض نفس اللحظة في منطقتين زمنيتين.\n\nالنتيجة كلها أو خطأ واضح؛ الجهاز اللي شغّل الكود ما يختارش المناطق بدالك.\n\nأسماء المناطق IANA وليها قواعد بتختلف حسب الموسم.")}},
	{ID: "c192.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch19Walkthrough("c192.idea-in", "~Time.In~ بيغيّر طريقة عرض اللحظة؛ ~Time.Equal~ يتأكد إن اللحظة نفسها ما اتحركتش.", ch19SunsInCode, "e19-2-in", []string{"fmt", "time"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "أنشئ لحظة بتوقيت UTC."},
		{Line: 3, Caption: "جهّز منطقة توضيحية."},
		{Line: 4, Caption: "اعرض نفس اللحظة محلياً."},
		{Line: 5, Caption: "اللحظة ثابتة؛ العرض اتغير.", Output: "true 04:00 +0400"},
	})}},
	{ID: "c192.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch19Question("c192.try-zone", "توقّع الناتج:\n\n~~~go\n"+ch19SunsPreviousDayCode+"\n~~~", []string{"false 2031-01-01 00:00 +0000", "true 2030-12-31 19:00 -0500", "true 2031-01-01 05:00 +0500"}, 1, []string{"Equal يقارن اللحظة حتى لو المنطقة مختلفة.", "التحويل للخلف يعدّي للتاريخ السابق.", "المنطقة هنا offset سالب خمس ساعات."}),
		ch19Walkthrough("c192.try-reveal", "تابع تحويل لحظة للتوقيت السالب.", ch19SunsPreviousDayCode, "e19-2-date", []string{"fmt", "time"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "ابدأ بمنتصف الليل UTC."},
			{Line: 3, Caption: "حوّل لنفس اللحظة."},
			{Line: 4, Caption: "اعرض اللحظة والتاريخ المحلي.", Output: "true 2030-12-31 19:00 -0500"},
		}),
	}},
	{ID: "c192.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch19Markdown("c192.watch", "~LoadLocation~ بترجع error للاسم المجهول؛ خلي قواعد ~time/tzdata~ متاحة.\n\nقارن التاريخ كامل بالسنة والشهر واليوم، مش رقم اليوم بس.")}},
	{ID: "c192.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch19Markdown("c192.plan", "حلّ اللحظة والمنطقتين قبل تكوين النتيجة.\n\nاستخدم format واضح فيه offset رقمي، وما تعتمدش على ~Local~."),
		ch19Question("c192.check-in", "إيه اللي بيحافظ على نفس اللحظة في منطقة مختلفة؟", []string{"~In~", "~Add~", "~time.Now~"}, 0, []string{"In بتغيّر الـLocation المرتبطة بالعرض.", "Add بتزود مدة وبتحرّك اللحظة.", "Now بتجيب لحظة الجهاز الحالية."}),
	}},
	{ID: "c192.turn-step", Title: "دورك", Blocks: []LessonBlock{ch19Markdown("c192.turn", "جرّب لحظة توضيحية، شغّلها في منطقتين، وبعدها راجع فرق التاريخ والـoffset.")}},
}

var ch19SunsEnglish = []LessonStep{
	{ID: "c192.mission-step", Title: "Mission", Blocks: []LessonBlock{ch19Markdown("c192.mission", "Render the same instant in two time zones.\n\nReturn a full result or a clear error; the host machine should not choose the zones for you.\n\nIANA locations can have seasonal rules.")}},
	{ID: "c192.idea-step", Title: "Idea", Blocks: []LessonBlock{ch19Walkthrough("c192.idea-in", "~Time.In~ changes how an instant is displayed; ~Time.Equal~ checks that the instant did not move.", ch19SunsInCode, "e19-2-in", []string{"fmt", "time"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create a UTC instant."},
		{Line: 3, Caption: "Set an illustrative location."},
		{Line: 4, Caption: "Render the instant locally."},
		{Line: 5, Caption: "Same instant, changed view.", Output: "true 04:00 +0400"},
	})}},
	{ID: "c192.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch19Question("c192.try-zone", "Predict the output:\n\n~~~go\n"+ch19SunsPreviousDayCode+"\n~~~", []string{"false 2031-01-01 00:00 +0000", "true 2030-12-31 19:00 -0500", "true 2031-01-01 05:00 +0500"}, 1, []string{"Equal compares instants even when locations differ.", "Converting backward crosses into the previous date.", "This zone has an offset of minus five hours."}),
		ch19Walkthrough("c192.try-reveal", "Follow an instant converted to a negative offset.", ch19SunsPreviousDayCode, "e19-2-date", []string{"fmt", "time"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Start at UTC midnight."},
			{Line: 3, Caption: "Convert the same instant."},
			{Line: 4, Caption: "Print instant and local date.", Output: "true 2030-12-31 19:00 -0500"},
		}),
	}},
	{ID: "c192.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch19Markdown("c192.watch", "~LoadLocation~ returns an error for an unknown name; make ~time/tzdata~ available.\n\nCompare the full year, month, and day, not just the day number.")}},
	{ID: "c192.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch19Markdown("c192.plan", "Resolve the instant and both locations before building a result.\n\nUse a format with a numeric offset; do not depend on ~Local~."),
		ch19Question("c192.check-in", "What keeps an instant fixed in another location?", []string{"~In~", "~Add~", "~time.Now~"}, 0, []string{"In changes the Location used for display.", "Add advances the instant by a duration.", "Now reads the host machine's current instant."}),
	}},
	{ID: "c192.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch19Markdown("c192.turn", "Try an illustrative instant in two zones, then review the date and offset difference.")}},
}

var ch19CountdownArabic = []LessonStep{
	{ID: "c193.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch19Markdown("c193.mission", "كمّل countdown بتقبل cancellation وبتستخدم clock محقونة.\n\nخلي كل مسار ينضّف ticker من غير نوم حقيقي أو تسريب.\n\nالإلغاء ممكن يحصل قبل أول tick أو بعده؛ الاختبار هو اللي بيحرّك الوقت.")}},
	{ID: "c193.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch19Walkthrough("c193.idea-context", "~Context.Done~ إشارة للإلغاء؛ ~select~ يقدر يراقبها مع شغل تاني.", ch19CountdownContextCode, "e19-3-context", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "جهّز سياق قابل للإلغاء."},
		{Line: 4, Caption: "راقب إشارة الإلغاء."},
		{Line: 8, Caption: "قبل الإلغاء، البرنامج شغال.", Output: "running"},
	})}},
	{ID: "c193.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch19Question("c193.try-cancel", "توقّع الناتج:\n\n~~~go\n"+ch19CountdownCanceledCode+"\n~~~", []string{"<nil>", "context canceled", "running"}, 1, []string{"السياق اتلغى قبل قراءة Err.", "Err بتوضح إن الإلغاء حصل.", "running كانت هتظهر قبل الإلغاء."}),
		ch19Walkthrough("c193.try-reveal", "شوف قيمة Err بعد إلغاء السياق.", ch19CountdownCanceledCode, "e19-3-cancel", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "أنشئ context جديد."},
			{Line: 3, Caption: "نفّذ الإلغاء."},
			{Line: 4, Caption: "اقرأ سبب الإلغاء.", Output: "context canceled"},
		}),
	}},
	{ID: "c193.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch19Markdown("c193.watch", "اعمل ticker بعد التحقق، و~defer Stop~ بعد إنشائها مباشرة.\n\nبعد tick، راجع ~ctx.Err~ قبل العرض؛ الإلغاء والـtick ممكن يبقوا جاهزين.")}},
	{ID: "c193.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch19Markdown("c193.plan", "الـFakeClock بتحرّك الوقت وقت الاختبار من غير انتظار.\n\nاختبر بداية العملية، tick، الإلغاء، وتنظيف الموارد كل مسار لوحده."),
		ch19Question("c193.check-tick", "بعد ما ~select~ يستقبل tick، إمتى تفحص الإلغاء؟", []string{"قبل الـcallback", "بعد ما ترجع الدالة", "مرة واحدة قبل الـloop"}, 0, []string{"الإلغاء ممكن يكون جاهز مع tick.", "بعد الرجوع يكون العرض حصل بالفعل.", "الفحص الأول مش بيغطي الإلغاء اللاحق."}),
	}},
	{ID: "c193.turn-step", Title: "دورك", Blocks: []LessonBlock{ch19Markdown("c193.turn", "خلّي FakeClock يقدّم الوقت خطوة خطوة، وراقب القيم والتنظيف بعد كل مسار.")}},
}

var ch19CountdownEnglish = []LessonStep{
	{ID: "c193.mission-step", Title: "Mission", Blocks: []LessonBlock{ch19Markdown("c193.mission", "Finish a countdown that accepts cancellation and an injected clock.\n\nClean up the ticker on every path without real sleeps or leaks.\n\nCancellation may happen before or after a tick; tests advance time.")}},
	{ID: "c193.idea-step", Title: "Idea", Blocks: []LessonBlock{ch19Walkthrough("c193.idea-context", "~Context.Done~ signals cancellation; ~select~ can watch it alongside other work.", ch19CountdownContextCode, "e19-3-context", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create a cancelable context."},
		{Line: 4, Caption: "Watch for cancellation."},
		{Line: 8, Caption: "The program is still running.", Output: "running"},
	})}},
	{ID: "c193.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch19Question("c193.try-cancel", "Predict the output:\n\n~~~go\n"+ch19CountdownCanceledCode+"\n~~~", []string{"<nil>", "context canceled", "running"}, 1, []string{"The context was canceled before Err was read.", "Err reports that cancellation happened.", "running would appear before cancellation."}),
		ch19Walkthrough("c193.try-reveal", "See Err after canceling the context.", ch19CountdownCanceledCode, "e19-3-cancel", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Create a new context."},
			{Line: 3, Caption: "Cancel the context."},
			{Line: 4, Caption: "Read the cancellation reason.", Output: "context canceled"},
		}),
	}},
	{ID: "c193.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch19Markdown("c193.watch", "Create the ticker after validation, then ~defer Stop~ immediately.\n\nAfter a tick, check ~ctx.Err~ before emitting; cancellation and a tick can both be ready.")}},
	{ID: "c193.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch19Markdown("c193.plan", "A ~FakeClock~ advances time during tests without waiting.\n\nTest start, tick, cancellation, and cleanup paths independently."),
		ch19Question("c193.check-tick", "After ~select~ receives a tick, when should cancellation be checked?", []string{"Before the callback", "After the function returns", "Only once before the loop"}, 0, []string{"Cancellation can be ready alongside the tick.", "The emission would already have happened.", "The first check cannot catch later cancellation."}),
	}},
	{ID: "c193.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch19Markdown("c193.turn", "Advance FakeClock one step at a time and watch values and cleanup on each path.")}},
}

var ch19TimersArabic = []LessonStep{
	{ID: "c194.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch19Markdown("c194.mission", "انتظر timers الأرصفة بالتوازي، واعرض الردود بترتيب الطلبات.\n\nعامل ممكن يخلص قبل غيره؛ ما تخليش وقت الوصول يغيّر مكان النتيجة.")}},
	{ID: "c194.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch19Walkthrough("c194.idea-wait", "الـWaitGroup بتجمع اكتمال العمال. الـframes بتعرض ترتيب محتمل واحد؛ توقيت الـgoroutines الفعلي ممكن يختلف.", ch19TimersWaitCode, "e19-4-wait", []string{"fmt", "sync"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "جهّز عدّاد العمال."},
		{Line: 4, Caption: "ابدأ goroutine واحدة."},
		{Line: 6, Caption: "العامل يطبع قبل Done.", Output: "worker"},
		{Line: 8, Caption: "Wait رجعت بعد اكتمال العامل.", Output: "worker"},
		{Line: 9, Caption: "كمل بعد انتهاء العمال.", Output: "worker\nready"},
	})}},
	{ID: "c194.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch19Question("c194.try-slots", "تخيّل إن slot التانية خلصت الأول؛ توقّع الناتج:\n\n~~~go\n"+ch19TimersSlotCode+"\n~~~", []string{"[fast slow]", "[slow fast]", "[slow]"}, 1, []string{"الطباعة تتبع مواقع الـslice مش وقت الإكمال.", "كل نتيجة فضلت في slot بتاعة طلبها.", "الـslice اتعملت بطول اتنين."}),
		ch19Walkthrough("c194.try-reveal", "الكتابة في index محدد تحفظ ترتيب الطلب.", ch19TimersSlotCode, "e19-4-slots", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "احجز خانتين للردود."},
			{Line: 3, Caption: "الرد التاني وصل الأول."},
			{Line: 4, Caption: "الأول وصل بعده."},
			{Line: 5, Caption: "اعرض حسب ترتيب الـindex.", Output: "[slow fast]"},
		}),
	}},
	{ID: "c194.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch19Markdown("c194.watch", "الـworkers ممكن يخلصوا بأي ترتيب؛ خلي لكل كتابة مكان مستقل.\n\nاستنى اكتمالهم قبل قراءة النتيجة، وما تعملش append مشترك.")}},
	{ID: "c194.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch19Markdown("c194.plan", "افصل صلاحية الدفعة عن تشغيل العمال.\n\nحدّد مين يملك كل نتيجة، وإمتى بقى آمن تقراها."),
		ch19Question("c194.check-wait", "إيه اللي تضمنه ~WaitGroup.Wait~؟", []string{"كل worker أبلغ عن اكتماله", "ترتيب الإكمال", "ترتيب إنشاء الـgoroutines"}, 0, []string{"Wait تستنى كل Done المقابلة.", "الانتظار ما بيفرضش ترتيب تنفيذ.", "التشغيل مش معناه ترتيب ثابت."}),
	}},
	{ID: "c194.turn-step", Title: "دورك", Blocks: []LessonBlock{ch19Markdown("c194.turn", "ابدأ بطلبين توضيحيين، شغّل fake clock، وبعدها راجع ترتيب الردود والتنظيف.")}},
}

var ch19TimersEnglish = []LessonStep{
	{ID: "c194.mission-step", Title: "Mission", Blocks: []LessonBlock{ch19Markdown("c194.mission", "Await dock timers concurrently and display replies in request order.\n\nOne worker can finish first; arrival time should not move its result.")}},
	{ID: "c194.idea-step", Title: "Idea", Blocks: []LessonBlock{ch19Walkthrough("c194.idea-wait", "A ~WaitGroup~ joins worker completion. These frames show one possible schedule; actual goroutine timing can vary.", ch19TimersWaitCode, "e19-4-wait", []string{"fmt", "sync"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Prepare the worker count."},
		{Line: 4, Caption: "Start one goroutine."},
		{Line: 6, Caption: "Worker prints before Done.", Output: "worker"},
		{Line: 8, Caption: "Wait returns after completion.", Output: "worker"},
		{Line: 9, Caption: "Continue after the workers.", Output: "worker\nready"},
	})}},
	{ID: "c194.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch19Question("c194.try-slots", "If the second slot finishes first, predict the output:\n\n~~~go\n"+ch19TimersSlotCode+"\n~~~", []string{"[fast slow]", "[slow fast]", "[slow]"}, 1, []string{"Printing follows slice positions, not finish time.", "Each result stays in its request slot.", "The slice was created with length two."}),
		ch19Walkthrough("c194.try-reveal", "Writing at an index preserves request order.", ch19TimersSlotCode, "e19-4-slots", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Reserve two response slots."},
			{Line: 3, Caption: "The second reply arrived first."},
			{Line: 4, Caption: "The first reply arrived next."},
			{Line: 5, Caption: "Print in index order.", Output: "[slow fast]"},
		}),
	}},
	{ID: "c194.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch19Markdown("c194.watch", "Workers can finish in any order; give each write its own location.\n\nWait before reading results, and avoid a shared append.")}},
	{ID: "c194.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch19Markdown("c194.plan", "Separate batch validation from starting workers.\n\nDecide who owns each result and when it is safe to read."),
		ch19Question("c194.check-wait", "What does ~WaitGroup.Wait~ guarantee?", []string{"Every worker reported completion", "Completion order", "Goroutine creation order"}, 0, []string{"Wait returns after the matching Done calls.", "Waiting does not force an execution order.", "Starting workers does not define their order."}),
	}},
	{ID: "c194.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch19Markdown("c194.turn", "Start with two illustrative jobs, run the fake clock, then review response order and cleanup.")}},
}

const ch19DurationParseCode = `func main() {
	d, err := time.ParseDuration("35m")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d, d.String())
}`

const ch19DurationAddCode = `func main() {
	start := time.Date(2031, 4, 3, 23, 40, 0, 0, time.UTC)
	arrival := start.Add(35 * time.Minute)
	fmt.Println(arrival.Format("15:04"))
}`

const ch19SunsInCode = `func main() {
	instant := time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC)
	zone := time.FixedZone("UTC+4", 4*60*60)
	local := instant.In(zone)
	fmt.Println(instant.Equal(local), local.Format("15:04 -0700"))
}`

const ch19SunsPreviousDayCode = `func main() {
	t := time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC)
	local := t.In(time.FixedZone("UTC-5", -5*60*60))
	fmt.Println(t.Equal(local), local.Format("2006-01-02 15:04 -0700"))
}`

const ch19CountdownContextCode = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	default:
		fmt.Println("running")
	}
}`

const ch19CountdownCanceledCode = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(ctx.Err())
}`

const ch19TimersWaitCode = `func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("worker")
	}()
	wg.Wait()
	fmt.Println("ready")
}`

const ch19TimersSlotCode = `func main() {
	results := make([]string, 2)
	results[1] = "fast"
	results[0] = "slow"
	fmt.Println(results)
}`
