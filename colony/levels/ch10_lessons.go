package levels

func ch10Pair(arabic, english string) map[string]any {
	return map[string]any{"ar": arabic, "en": english}
}

var ch10WaitSandbox = &LessonSandbox{
	Kind: "ch10-waitgroup-order",
	Config: map[string]any{
		"intro":    ch10Pair("Wait بتخلص بعد ما الشغل المسجل يخلص.", "Wait returns after registered work ends."),
		"button":   ch10Pair("الخطوة التالية / تشغيل تاني", "Next step / run again"),
		"logLabel": ch10Pair("ترتيب ممكن", "One possible order"),
		"runs": []any{
			map[string]any{"events": []any{
				ch10Pair("main خلصت شغلها.", "main finishes its work."),
				ch10Pair("helper خلصت المهمة.", "helper finishes its task."),
				ch10Pair("Wait خلصت بعد ما helper خلص.", "Wait returns after helper finishes."),
			}},
			map[string]any{"events": []any{
				ch10Pair("helper خلصت المهمة.", "helper finishes its task."),
				ch10Pair("main خلصت شغلها.", "main finishes its work."),
				ch10Pair("Wait خلصت بعد ما helper خلص.", "Wait returns after helper finishes."),
			}},
		},
	},
}

var ch10SwarmSandbox = &LessonSandbox{
	Kind: "ch10-swarm-order",
	Config: map[string]any{
		"intro":    ch10Pair("كل تشغيل بيوضح ترتيب ممكن.", "Each run shows one possible order."),
		"button":   ch10Pair("الخطوة التالية / تشغيل تاني", "Next step / run again"),
		"logLabel": ch10Pair("سجل الأحداث", "Event log"),
		"runs": []any{
			map[string]any{"events": []any{
				ch10Pair("A تبدأ.", "A starts."), ch10Pair("B تبدأ.", "B starts."),
				ch10Pair("A بتطبع.", "A prints."), ch10Pair("A تنادي Done.", "A calls Done."),
				ch10Pair("B بتطبع.", "B prints."), ch10Pair("B تنادي Done.", "B calls Done."),
				ch10Pair("Wait بتخلص بعد الاتنين.", "Wait returns after both."),
			}},
			map[string]any{"events": []any{
				ch10Pair("B تبدأ.", "B starts."), ch10Pair("A تبدأ.", "A starts."),
				ch10Pair("B بتطبع.", "B prints."), ch10Pair("B تنادي Done.", "B calls Done."),
				ch10Pair("A بتطبع.", "A prints."), ch10Pair("A تنادي Done.", "A calls Done."),
				ch10Pair("Wait بتخلص بعد الاتنين.", "Wait returns after both."),
			}},
		},
	},
}

var ch10DivideSandbox = &LessonSandbox{
	Kind: "ch10-groups-order",
	Config: map[string]any{
		"intro":       ch10Pair("كل مجموعة لها شغلها؛ النهاية ممكن تتغير.", "Each group owns tasks; finish order can vary."),
		"button":      ch10Pair("الخطوة التالية / تشغيل تاني", "Next step / run again"),
		"groupsLabel": ch10Pair("مجموعات توضيحية", "Illustrative groups"),
		"logLabel":    ch10Pair("ترتيب انتهاء ممكن", "One possible finish order"),
		"groups": []any{
			ch10Pair("A: [cup, lamp]", "A: [cup, lamp]"),
			ch10Pair("B: [map, ring]", "B: [map, ring]"),
			ch10Pair("C: [key, bowl]", "C: [key, bowl]"),
		},
		"runs": []any{
			map[string]any{"events": []any{
				ch10Pair("A تخلص.", "A finishes."), ch10Pair("C تخلص.", "C finishes."),
				ch10Pair("B تخلص.", "B finishes."), ch10Pair("كل المهام خلصت.", "All tasks finish."),
			}},
			map[string]any{"events": []any{
				ch10Pair("B تخلص.", "B finishes."), ch10Pair("A تخلص.", "A finishes."),
				ch10Pair("C تخلص.", "C finishes."), ch10Pair("كل المهام خلصت.", "All tasks finish."),
			}},
		},
	},
}

const (
	ch101WaitExample = `var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("helper finished")
}()
wg.Wait()
fmt.Println("all done")`

	ch101ResultExample = `var wg sync.WaitGroup
message := "waiting"
wg.Add(1)
go func() {
	defer wg.Done()
	message = "ready"
}()
wg.Wait()
fmt.Println(message)`

	ch102LoopExample = `var wg sync.WaitGroup
for _, label := range []string{"gold", "silver"} {
	wg.Add(1)
	go func(name string) {
		defer wg.Done()
		fmt.Println(name)
	}(label)
}
wg.Wait()`

	ch102OrderExample = `var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("A")
}()
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("B")
}()
wg.Wait()`

	ch103GroupsExample = `jobs := []string{"cup", "lamp", "map", "ring"}
parts := [][]string{
	jobs[:2],
	jobs[2:],
}
fmt.Println(parts[0], parts[1])`

	ch103CountExample = `groups := [][]int{{4, 9}, {3}}
for i, group := range groups {
	fmt.Println(i, len(group))
}`
)

func ch10Trace(lines []int, captions []string, outputs []string) *LessonWalkthrough {
	frames := make([]LessonFrame, len(lines))
	for i := range lines {
		output := outputs[i]
		if len(output) > 0 && output[len(output)-1] == '\n' {
			output = output[:len(output)-1]
		}
		frames[i] = LessonFrame{Line: lines[i], Caption: captions[i], Output: output}
	}
	return &LessonWalkthrough{Frames: frames}
}

var ch10LessonArabic = map[string][]LessonStep{
	"10-1-second-drone": {
		{ID: "c101.mission-step", Title: "المهمة", Blocks: []LessonBlock{
			{ID: "c101.mission", Kind: "markdown", Text: "1. احصد الصفين: إنت صف، ودرون جديد صف تاني.\n2. الـstarter جهّز helper وحصد صفك؛ كمّل شغله.\n3. خلي main تستنى المساعد.\n\nالصفين بعيدين عن بعض.\n\nhelper يبدأ من القاعدة ويوصل لصفه.\n\nلما main تخلص، البرنامج كله بيقفل.\n\n~harvestRow~ جاهزة وتاخد ~*colony.Drone~ واتجاه."},
		}},
		{ID: "c101.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
			{ID: "c101.wait-trace", Kind: "walkthrough", Text: "الـWaitGroup عدّاد للشغل: Add تسجل، Done تخلص، وWait تستنى الصفر.\n\nخلي كل Add يقابلها Done. كلمة ~go~ تشغّل الـclosure وتكمّل فورًا؛ مش بتستنى المهمة.\n\nده ترتيب واحد ممكن؛ بداية goroutine مش مضمونة.", Code: ch101WaitExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-1-wait", Walkthrough: ch10Trace(
				[]int{0, 1, 2, 3, 5, 6, 7, 8},
				[]string{"قبل.", "العداد بيبدأ بصفر.", "شغلانة اتسجلت.", "ترتيب ممكن: بدأت.", "المساعد بيطبع.", "Done تصفّر العداد.", "Wait خلصت.", "main بتطبع."},
				[]string{"", "", "", "", "helper finished\n", "helper finished\n", "helper finished\n", "helper finished\nall done\n"},
			)},
		}},
		{ID: "c101.try-step", Title: "جرّب", Blocks: []LessonBlock{
			{ID: "c101.predict", Kind: "question", Text: "توقع الـvalue اللي هتتطبع بعد Wait.", Code: ch101ResultExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-1-result", Question: &LessonQuestion{
				Options: []string{"waiting", "ready", "مفيش سطر"}, Answer: 1,
				Feedback: []string{"التغيير حصل قبل الطباعة.", "Wait استنت التغيير.", "فيه طباعة بعد Wait."},
			}},
			{ID: "c101.reveal", Kind: "walkthrough", Text: "تابع التغيير، وبعده الطباعة.", Code: ch101ResultExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-1-result", Walkthrough: ch10Trace(
				[]int{0, 2, 3, 4, 6, 8, 9},
				[]string{"قبل.", "message تساوي waiting.", "الشغل اتسجل.", "ترتيب ممكن: بدأت.", "message بقت ready.", "Wait خلصت.", "اطبع الـvalue."},
				[]string{"", "", "", "", "", "", "ready\n"},
			)},
			{ID: "c101.wait-sandbox", Kind: "sandbox", Text: "حرّك الأحداث؛ الترتيب ممكن يتغير.", Sandbox: ch10WaitSandbox},
		}},
		{ID: "c101.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
			{ID: "c101.watch", Kind: "markdown", Text: "- سجّل الشغل قبل go.\n- حط Done في defer؛ هتشتغل عند نهاية الـclosure.\n- Wait تمنع main من الخروج بدري.\n- Move وHarvest كل واحد بياخد tick، حتى لو Move فشلت.\n\nالـWaitGroup تراقب كل الشغل المسجل؛ مش بتشغّله."},
		}},
		{ID: "c101.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
			{ID: "c101.plan", Kind: "markdown", Text: "المراحل: وزّع الشغل، شغّل المساعد، واستنى اكتمال المهام.\n\nWait تستنى الكل حتى لو main خلصت شغلها."},
			{ID: "c101.q-wait", Kind: "question", Text: "إمتى main تقدر تكمل؟", Question: &LessonQuestion{
				Options: []string{"بعد Wait", "بعد go", "بعد Add"}, Answer: 0,
				Feedback: []string{"Wait بتخلص عند الصفر.", "go تبدأ الشغل.", "Add تسجل المهمة."},
			}},
			{ID: "c101.q-done", Kind: "question", Text: "أنهي أمر يقلل العداد؟", Question: &LessonQuestion{
				Options: []string{"Add", "Done", "Wait"}, Answer: 1,
				Feedback: []string{"Add تسجل شغل.", "Done تخلص مهمة.", "Wait تستنى."},
			}},
		}},
		{ID: "c101.turn-step", Title: "دورك", Blocks: []LessonBlock{
			{ID: "c101.turn", Kind: "markdown", Text: "كمّل شغل helper وخلي البرنامج يستناه. شغّل المهمة وراقب الصفين."},
		}},
	},
	"10-2-swarm": {
		{ID: "c102.mission-step", Title: "المهمة", Blocks: []LessonBlock{
			{ID: "c102.mission", Kind: "markdown", Text: "1. وزّع الصفوف على الدرون الأصلي والتلاتة الجداد.\n2. الـstarter فيه loop؛ جسم work فاضي.\n3. شغّل الشغل واستنى الجميع.\n\nلكل درون مكانه. ترتيب الطباعة مش ترتيب البداية أو النهاية."},
		}},
		{ID: "c102.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
			{ID: "c102.loop-trace", Kind: "walkthrough", Text: "الـloop تجهّز مهمة لكل value. الـtrace ترتيب واحد ممكن؛ الرسائل ممكن تتعكس.\n\nتمرير name كـparameter يربط الـvalue بمهمتها. مرّر row عشان كل مهمة تعرف صفها؛ الـvalue الممررة تخصها.\n\nWait تضمن انتهاء الشغل المسجل بس.", Code: ch102LoopExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-2-loop", Walkthrough: ch10Trace(
				[]int{0, 2, 3, 4, 5, 6, 7, 8, 6, 9},
				[]string{"قبل.", "الـvalue الأولى.", "مهمة اتسجلت.", "goroutine بدأت.", "تجهّز Done.", "ترتيب ممكن: gold.", "Done تشتغل.", "الـvalue التانية.", "اطبع silver.", "Wait بتخلص بعد الكل."},
				[]string{"", "", "", "", "", "", "gold\n", "gold\n", "gold\nsilver\n", "gold\nsilver\n"},
			)},
		}},
		{ID: "c102.try-step", Title: "جرّب", Blocks: []LessonBlock{
			{ID: "c102.predict", Kind: "question", Text: "توقع ترتيب طباعة goroutineين.", Code: ch102OrderExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-2-order", Question: &LessonQuestion{
				Options: []string{"A ثم B بس", "B ثم A بس", "A ثم B أو B ثم A"}, Answer: 2,
				Feedback: []string{"A مش مضمونة الأول.", "B مش مضمونة الأول.", "الترتيب متغير؛ Wait تستنى الاتنين."},
			}},
			{ID: "c102.reveal", Kind: "walkthrough", Text: "ده ترتيب واحد: A قبل B. ممكن يتعكس؛ Wait تضمن انتهاء الاتنين بس.", Code: ch102OrderExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-2-order", Walkthrough: ch10Trace(
				[]int{0, 2, 3, 5, 6, 7, 8, 10, 11, 12},
				[]string{"قبل.", "مهمة A اتسجلت.", "goroutine A بدأت.", "ترتيب ممكن: اطبع A.", "A خلصت.", "مهمة B اتسجلت.", "goroutine B بدأت.", "اطبع B هنا.", "B خلصت.", "Wait خلصت."},
				[]string{"", "", "", "", "A\n", "A\n", "A\n", "A\n", "A\nB\n", "A\nB\n"},
			)},
			{ID: "c102.swarm-sandbox", Kind: "sandbox", Text: "تابع التشغيل، وبعدها جرّب ترتيبًا تاني.", Sandbox: ch10SwarmSandbox},
		}},
		{ID: "c102.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
			{ID: "c102.watch", Kind: "markdown", Text: "- زوّد العداد قبل go.\n- كل مهمة تنادي Done مرة.\n- مرّر الـvalue بتاعة الصف للمهمة بوضوح.\n\nابدأ Wait بعد تسجيل المهام؛ ترتيب الطباعة ممكن يتغير."},
		}},
		{ID: "c102.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
			{ID: "c102.plan", Kind: "markdown", Text: "المراحل: جهّز شغل الصفوف، شغّل المهام، واستنى الكل. تابع مكان كل درون."},
			{ID: "c102.q-add", Kind: "question", Text: "إمتى تسجّل المهمة؟", Question: &LessonQuestion{
				Options: []string{"قبل go", "بعد Wait", "بعد ما تخلص"}, Answer: 0,
				Feedback: []string{"سجّلها قبل التشغيل.", "Wait ممكن تخلص بدري.", "Done مش بديل عن Add."},
			}},
		}},
		{ID: "c102.turn-step", Title: "دورك", Blocks: []LessonBlock{
			{ID: "c102.turn", Kind: "markdown", Text: "كمّل حركة وحصاد work. شغّل الكود مرتين ولاحظ تغيّر ترتيب الرسائل."},
		}},
	},
	"10-3-divide": {
		{ID: "c103.mission-step", Title: "المهمة", Blocks: []LessonBlock{
			{ID: "c103.mission", Kind: "markdown", Text: "1. قسّم 15 هدف متفرق على تلاتة درونات.\n2. الـstarter جهّز targets وgoTo وcollect؛ اختار التوزيع.\n3. وازن الرحلات عشان درون مايتأخرش.\n\nالمسافة بتفرق جنب عدد الأهداف. مواقع targets بتتغير كل تشغيل."},
		}},
		{ID: "c103.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
			{ID: "c103.groups-trace", Kind: "walkthrough", Text: "~[][]string~ slice من الـslices؛ كل slice فيها مهام للتوضيح.\n\nفي المهمة، ~colony.Targets()~ بترجّع ~[]colony.Point~. كل نقطة فيها X وY.\n\nالنوع ~[][]colony.Point~ يجمع slices من النقاط؛ استخدم ~len~ لمقارنة حجم المجموعات. المثال يعرض الشكل بس.", Code: ch103GroupsExample, Imports: []string{"fmt"}, ExampleID: "e10-3-groups", Walkthrough: ch10Trace(
				[]int{0, 1, 2, 3, 4, 6},
				[]string{"قبل.", "slice مهام.", "slice مجموعات.", "المجموعة الأولى.", "المجموعة التانية.", "اطبعهم."},
				[]string{"", "", "", "", "", "[cup lamp] [map ring]\n"},
			)},
		}},
		{ID: "c103.try-step", Title: "جرّب", Blocks: []LessonBlock{
			{ID: "c103.predict", Kind: "question", Text: "توقع خرج الـloop؛ الأرقام للتوضيح.", Code: ch103CountExample, Imports: []string{"fmt"}, ExampleID: "e10-3-count", Question: &LessonQuestion{
				Options: []string{"0 2 ثم 1 1", "0 1 ثم 1 2", "1 2 بس"}, Answer: 0,
				Feedback: []string{"الأولى فيها اتنين.", "الأعداد متعكسة.", "كل مجموعة ليها سطر."},
			}},
			{ID: "c103.reveal", Kind: "walkthrough", Text: "range تمر على كل مجموعة؛ len تعد العناصر. العدد وحده مش وقت الرحلة.", Code: ch103CountExample, Imports: []string{"fmt"}, ExampleID: "e10-3-count", Walkthrough: ch10Trace(
				[]int{0, 1, 2, 3, 2, 3},
				[]string{"قبل.", "مجموعتان للتوضيح.", "range تمسك الأولى.", "اطبع رقمها وعددها.", "range تمسك التانية.", "اطبع الرقم والعدد."},
				[]string{"", "", "", "0 2\n", "0 2\n", "0 2\n1 1\n"},
			)},
			{ID: "c103.groups-sandbox", Kind: "sandbox", Text: "شوف المجموعات وبدّل ترتيب النهاية.", Sandbox: ch10DivideSandbox},
		}},
		{ID: "c103.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
			{ID: "c103.watch", Kind: "markdown", Text: "- العدد مش كفاية؛ الطريق يفرق.\n- جمع أهداف بعيدة يطوّل الرحلة.\n- كل حركة وHarvest بياخدوا tick.\n\nقارن القرب قبل التقسيم."},
		}},
		{ID: "c103.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
			{ID: "c103.plan", Kind: "markdown", Text: "المراحل: افحص المواقع، كوّن مجموعات قريبة، وقارن أطول رحلة.\n\nترتيب الزيارات من الفصل 7 يقلل الرجوع.\n\nالوقت الكلي يستنى آخر درون يخلص."},
			{ID: "c103.q-balance", Kind: "question", Text: "إيه اللي يغيّر وقت مجموعتين متساويتين؟", Question: &LessonQuestion{
				Options: []string{"المسافة", "اسم المتغير", "الأقواس"}, Answer: 0,
				Feedback: []string{"الحركة تضيف وقت.", "الاسم مش بيغيّر الطريق.", "الأقواس مش بتغيّر المسافة."},
			}},
			{ID: "c103.q-type", Kind: "question", Text: "إيه شكل ~[][]colony.Point~؟", Question: &LessonQuestion{
				Options: []string{"slice من مجموعات نقاط", "نقطة واحدة", "عداد مهام"}, Answer: 0,
				Feedback: []string{"كل slice داخلية مجموعة نقاط.", "النوع بيضم مجموعات.", "ده نوع بيانات مش عداد."},
			}},
		}},
		{ID: "c103.turn-step", Title: "دورك", Blocks: []LessonBlock{
			{ID: "c103.turn", Kind: "markdown", Text: "جرّب تقسيمة أولى. المواقع بتتغير؛ راقب أبطأ درون وعدّل حسب العدد والمسافة."},
		}},
	},
}

var ch10LessonEnglish = map[string][]LessonStep{
	"10-1-second-drone": {
		{ID: "c101.mission-step", Title: "Mission", Blocks: []LessonBlock{
			{ID: "c101.mission", Kind: "markdown", Text: "1. Harvest one row with main and one with the helper.\n2. The starter spawns the helper and harvests your row.\n3. Finish the helper's work and wait for it.\n\nThe rows are apart.\n\nThe helper starts at base and must reach its row.\n\nWhen main returns, the program exits.\n\n~harvestRow~ accepts a ~*colony.Drone~ and a direction, so both drones can use it."},
		}},
		{ID: "c101.idea-step", Title: "Idea", Blocks: []LessonBlock{
			{ID: "c101.wait-trace", Kind: "walkthrough", Text: "A WaitGroup counts tasks: Add registers, Done finishes, and Wait blocks until zero.\n\nBalance each Add with a Done. ~go~ starts the closure and returns; it does not wait.\n\nThis trace shows ONE possible order; goroutine timing is not fixed.", Code: ch101WaitExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-1-wait", Walkthrough: ch10Trace(
				[]int{0, 1, 2, 3, 5, 6, 7, 8},
				[]string{"Before.", "Counter zero.", "Task registered.", "One possible start.", "Helper prints.", "Done reaches zero.", "Wait returns.", "Main prints."},
				[]string{"", "", "", "", "helper finished\n", "helper finished\n", "helper finished\n", "helper finished\nall done\n"},
			)},
		}},
		{ID: "c101.try-step", Title: "Try it", Blocks: []LessonBlock{
			{ID: "c101.predict", Kind: "question", Text: "Predict the value of message printed after Wait.", Code: ch101ResultExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-1-result", Question: &LessonQuestion{
				Options: []string{"waiting", "ready", "No line"}, Answer: 1,
				Feedback: []string{"It changes before printing.", "Wait sees the change.", "Print runs after Wait."},
			}},
			{ID: "c101.reveal", Kind: "walkthrough", Text: "Follow the write, wait, and print.", Code: ch101ResultExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-1-result", Walkthrough: ch10Trace(
				[]int{0, 2, 3, 4, 6, 8, 9},
				[]string{"Before.", "message is waiting.", "Task registered.", "One possible start.", "message becomes ready.", "Wait returns.", "Print the value."},
				[]string{"", "", "", "", "", "", "ready\n"},
			)},
			{ID: "c101.wait-sandbox", Kind: "sandbox", Text: "Step through events; order can vary.", Sandbox: ch10WaitSandbox},
		}},
		{ID: "c101.watch-step", Title: "Watch out", Blocks: []LessonBlock{
			{ID: "c101.watch", Kind: "markdown", Text: "- Add before go.\n- Defer Done; it runs when the closure ends.\n- Wait before main returns.\n- Move and Harvest each cost a tick, even when Move fails.\n\nA WaitGroup tracks all registered work; it does not start it."},
		}},
		{ID: "c101.check-step", Title: "Think & check", Blocks: []LessonBlock{
			{ID: "c101.plan", Kind: "markdown", Text: "Phases: assign work, start the helper, wait for completion.\n\nWait waits for every task, even after main finishes its own work."},
			{ID: "c101.q-wait", Kind: "question", Text: "When can main continue?", Question: &LessonQuestion{
				Options: []string{"After Wait", "After go", "After Add"}, Answer: 0,
				Feedback: []string{"Wait returns at zero.", "go starts work.", "Add registers work."},
			}},
			{ID: "c101.q-done", Kind: "question", Text: "Which call lowers the count?", Question: &LessonQuestion{
				Options: []string{"Add", "Done", "Wait"}, Answer: 1,
				Feedback: []string{"Add registers work.", "Done finishes a task.", "Wait blocks."},
			}},
		}},
		{ID: "c101.turn-step", Title: "Your turn", Blocks: []LessonBlock{
			{ID: "c101.turn", Kind: "markdown", Text: "Finish helper work and wait for it. Run the mission; check both rows."},
		}},
	},
	"10-2-swarm": {
		{ID: "c102.mission-step", Title: "Mission", Blocks: []LessonBlock{
			{ID: "c102.mission", Kind: "markdown", Text: "1. Divide the rows among main and three new drones.\n2. The starter has a loop; work is empty.\n3. Run the tasks and wait for all.\n\nEach drone has its own position. Print order is not start or finish order."},
		}},
		{ID: "c102.idea-step", Title: "Idea", Blocks: []LessonBlock{
			{ID: "c102.loop-trace", Kind: "walkthrough", Text: "The loop prepares one task per value. This trace shows ONE possible order; the prints may reverse.\n\nPassing name as a parameter ties it to the task. Pass row so each task knows its row; the value belongs to that task.\n\nWait only guarantees registered work ends.", Code: ch102LoopExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-2-loop", Walkthrough: ch10Trace(
				[]int{0, 2, 3, 4, 5, 6, 7, 8, 6, 9},
				[]string{"Before.", "First value.", "Task registered.", "Goroutine starts.", "Defer Done.", "One order: gold.", "Gold calls Done.", "Next value.", "Print silver.", "Wait returns after both."},
				[]string{"", "", "", "", "", "", "gold\n", "gold\n", "gold\nsilver\n", "gold\nsilver\n"},
			)},
		}},
		{ID: "c102.try-step", Title: "Try it", Blocks: []LessonBlock{
			{ID: "c102.predict", Kind: "question", Text: "Predict the print order of two goroutines.", Code: ch102OrderExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-2-order", Question: &LessonQuestion{
				Options: []string{"Only A then B", "Only B then A", "A then B or B then A"}, Answer: 2,
				Feedback: []string{"A is not guaranteed first.", "B is not guaranteed first.", "Order varies; Wait waits for both."},
			}},
			{ID: "c102.reveal", Kind: "walkthrough", Text: "This is ONE order: A before B. It may reverse; Wait guarantees only that both finish.", Code: ch102OrderExample, Imports: []string{"fmt", "sync"}, ExampleID: "e10-2-order", Walkthrough: ch10Trace(
				[]int{0, 2, 3, 5, 6, 7, 8, 10, 11, 12},
				[]string{"Before.", "Task A added.", "A starts.", "One order: print A.", "A ends.", "Task B added.", "B starts.", "Print B here.", "B ends.", "Wait returns."},
				[]string{"", "", "", "", "A\n", "A\n", "A\n", "A\n", "A\nB\n", "A\nB\n"},
			)},
			{ID: "c102.swarm-sandbox", Kind: "sandbox", Text: "Step through a run, then compare another order.", Sandbox: ch10SwarmSandbox},
		}},
		{ID: "c102.watch-step", Title: "Watch out", Blocks: []LessonBlock{
			{ID: "c102.watch", Kind: "markdown", Text: "- Add before go.\n- Each task calls Done once.\n- Pass the row value clearly.\n\nWait after the loop registers tasks; print order can change."},
		}},
		{ID: "c102.check-step", Title: "Think & check", Blocks: []LessonBlock{
			{ID: "c102.plan", Kind: "markdown", Text: "Phases: prepare row work, start the tasks, wait for all. Track each drone's position."},
			{ID: "c102.q-add", Kind: "question", Text: "When do you register a task?", Question: &LessonQuestion{
				Options: []string{"Before go", "After Wait", "After it ends"}, Answer: 0,
				Feedback: []string{"Register before launch.", "Wait could finish early.", "Done does not replace Add."},
			}},
		}},
		{ID: "c102.turn-step", Title: "Your turn", Blocks: []LessonBlock{
			{ID: "c102.turn", Kind: "markdown", Text: "Complete movement and harvesting in work. Run again and watch message order change."},
		}},
	},
	"10-3-divide": {
		{ID: "c103.mission-step", Title: "Mission", Blocks: []LessonBlock{
			{ID: "c103.mission", Kind: "markdown", Text: "1. Divide 15 scattered targets among three drones.\n2. The starter sets up targets, goTo, and collect; choose the split.\n3. Balance routes so one drone is not left behind.\n\nDistance matters along with target count. Locations change each run."},
		}},
		{ID: "c103.idea-step", Title: "Idea", Blocks: []LessonBlock{
			{ID: "c103.groups-trace", Kind: "walkthrough", Text: "~[][]string~ is a list of groups; each part is an illustrative task list.\n\nIn the mission, ~colony.Targets()~ returns ~[]colony.Point~. Each point has X and Y.\n\n~[][]colony.Point~ groups point lists; use ~len~ to compare group sizes. This example only shows the shape.", Code: ch103GroupsExample, Imports: []string{"fmt"}, ExampleID: "e10-3-groups", Walkthrough: ch10Trace(
				[]int{0, 1, 2, 3, 4, 6},
				[]string{"Before.", "Task list.", "List of groups.", "First group.", "Second group.", "Print both."},
				[]string{"", "", "", "", "", "[cup lamp] [map ring]\n"},
			)},
		}},
		{ID: "c103.try-step", Title: "Try it", Blocks: []LessonBlock{
			{ID: "c103.predict", Kind: "question", Text: "Predict the loop's output; numbers are illustrative.", Code: ch103CountExample, Imports: []string{"fmt"}, ExampleID: "e10-3-count", Question: &LessonQuestion{
				Options: []string{"0 2, then 1 1", "0 1, then 1 2", "Only 1 2"}, Answer: 0,
				Feedback: []string{"The first group has two.", "The counts are reversed.", "Each group prints a line."},
			}},
			{ID: "c103.reveal", Kind: "walkthrough", Text: "range visits each group; len counts its items. Count alone is not travel time.", Code: ch103CountExample, Imports: []string{"fmt"}, ExampleID: "e10-3-count", Walkthrough: ch10Trace(
				[]int{0, 1, 2, 3, 2, 3},
				[]string{"Before.", "Two example groups.", "range selects first.", "Print number and count.", "range selects next.", "Print next count."},
				[]string{"", "", "", "0 2\n", "0 2\n", "0 2\n1 1\n"},
			)},
			{ID: "c103.groups-sandbox", Kind: "sandbox", Text: "Inspect groups and compare finish orders.", Sandbox: ch10DivideSandbox},
		}},
		{ID: "c103.watch-step", Title: "Watch out", Blocks: []LessonBlock{
			{ID: "c103.watch", Kind: "markdown", Text: "- Item count is not enough; routes matter.\n- Grouping far targets can stretch a trip.\n- Each move and Harvest costs a tick.\n\nCompare distance before splitting."},
		}},
		{ID: "c103.check-step", Title: "Think & check", Blocks: []LessonBlock{
			{ID: "c103.plan", Kind: "markdown", Text: "Phases: inspect locations, group nearby targets, compare the longest route.\n\nChapter 7 showed how visit order reduces backtracking.\n\nTotal time follows the last drone to finish."},
			{ID: "c103.q-balance", Kind: "question", Text: "What changes the time of equal-sized groups?", Question: &LessonQuestion{
				Options: []string{"Distance", "Variable name", "Braces"}, Answer: 0,
				Feedback: []string{"Movement adds travel time.", "Names do not change routes.", "Braces do not change distance."},
			}},
			{ID: "c103.q-type", Kind: "question", Text: "What shape is ~[][]colony.Point~?", Question: &LessonQuestion{
				Options: []string{"A list of point groups", "One point", "A task counter"}, Answer: 0,
				Feedback: []string{"Each inner slice is a point group.", "The type holds groups.", "It is data, not a counter."},
			}},
		}},
		{ID: "c103.turn-step", Title: "Your turn", Blocks: []LessonBlock{
			{ID: "c103.turn", Kind: "markdown", Text: "Try one split. Locations change each run; watch the last drone and adjust for count and distance."},
		}},
	},
}
