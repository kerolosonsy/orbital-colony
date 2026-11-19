package levels

func ch13Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch13Question(id, text, code string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text + "\n\n~~~go\n" + code + "\n~~~",
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch13Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{
		ID:          id,
		Kind:        "walkthrough",
		Text:        text,
		Code:        code,
		Imports:     imports,
		ExampleID:   exampleID,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch13Locale(arabic, english string) map[string]any {
	return map[string]any{"ar": arabic, "en": english}
}

func ch13Sandbox(id, kind, text string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch131RecallSandbox(text string) LessonBlock {
	return ch13Sandbox("c131.sandbox", "ch13-recall", text, map[string]any{
		"workers": []string{"D1", "D2", "D3"},
		"orders":  [][]int{{0, 1, 2}, {2, 0, 1}, {1, 2, 0}},
		"labels": map[string]any{
			"signal":       ch13Locale("السياق", "Context"),
			"active":       ch13Locale("شغال", "active"),
			"canceled":     ch13Locale("اتلغى", "canceled"),
			"working":      ch13Locale("بيشتغل", "working"),
			"atBase":       ch13Locale("في القاعدة", "at base"),
			"joined":       ch13Locale("الكل خلص", "all joined"),
			"notJoined":    ch13Locale("لسه مستنيين", "still waiting"),
			"order":        ch13Locale("ترتيب محتمل", "Possible order"),
			"workersLabel": ch13Locale("الدرونات", "Drones"),
			"recall":       ch13Locale("ابعت الرجوع", "Send recall"),
			"advance":      ch13Locale("خطوة عامل", "Worker step"),
			"shuffleReset": ch13Locale("ترتيب تاني", "Shuffle again"),
			"wait":         ch13Locale("استنى الكل", "Wait for all"),
			"explanation":  ch13Locale("الإشارة ما بترجعش الدرون لوحدها.", "The signal does not return drones by itself."),
		},
	})
}

func ch132TimeoutSandbox(text string) LessonBlock {
	return ch13Sandbox("c132.sandbox", "ch13-timeout", text, map[string]any{
		"deadlineTicks": 3,
		"labels": map[string]any{
			"context":         ch13Locale("السياق", "Context"),
			"active":          ch13Locale("شغال", "active"),
			"deadlineReached": ch13Locale("المهلة انتهت", "deadline reached"),
			"manuallyStopped": ch13Locale("اتلغى يدويًا", "canceled manually"),
			"elapsed":         ch13Locale("الوقت", "Time"),
			"deadline":        ch13Locale("المهلة", "Deadline"),
			"done":            ch13Locale("Done", "Done"),
			"open":            ch13Locale("مفتوحة", "open"),
			"closed":          ch13Locale("مقفولة", "closed"),
			"err":             ch13Locale("Err", "Err"),
			"nil":             ch13Locale("nil", "nil"),
			"deadlineError":   ch13Locale("context.DeadlineExceeded", "context.DeadlineExceeded"),
			"canceledError":   ch13Locale("context.Canceled", "context.Canceled"),
			"tick":            ch13Locale("عدّي وقت", "Advance time"),
			"cancel":          ch13Locale("إلغاء يدوي", "Cancel manually"),
			"reset":           ch13Locale("من الأول", "Reset"),
		},
	})
}

func ch133RaceSandbox(text string) LessonBlock {
	return ch13Sandbox("c133.sandbox", "ch13-first-wins", text, map[string]any{
		"workers": []string{"D1", "D2", "D3"},
		"orders":  [][]int{{0, 1, 2}, {2, 0, 1}, {1, 2, 0}},
		"labels": map[string]any{
			"order":        ch13Locale("ترتيب البحث بيتغير", "Possible search order"),
			"signal":       ch13Locale("Done", "Done"),
			"open":         ch13Locale("لسه مفتوحة", "still open"),
			"canceled":     ch13Locale("اتقفلت", "closed"),
			"workersLabel": ch13Locale("العمال", "Workers"),
			"searching":    ch13Locale("بيدوّر", "searching"),
			"found":        ch13Locale("لقى الهدف", "found target"),
			"stopped":      ch13Locale("وقف بعد الإشارة", "stopped after signal"),
			"raceStep":     ch13Locale("شغّل محاولة", "Run search step"),
			"respond":      ch13Locale("عامل يفحص Done", "Worker checks Done"),
			"shuffleAgain": ch13Locale("رتّب وجرب تاني", "Shuffle and try again"),
			"explanation":  ch13Locale("الإلغاء إشارة؛ كل عامل يستجيب بنفسه.", "Cancellation signals; each worker responds."),
		},
	})
}

func ch134ChainSandbox(text string) LessonBlock {
	return ch13Sandbox("c134.sandbox", "ch13-context-chain", text, map[string]any{
		"code": []string{"command(ctx)", "relay(ctx)", "probe(ctx)", "return err"},
		"frames": []any{
			map[string]any{"line": 0, "caption": ch13Locale("القيادة تمرر السياق.", "Command passes context.")},
			map[string]any{"line": 1, "caption": ch13Locale("الوسيط يمرر نفس السياق.", "Relay passes the same context.")},
			map[string]any{"line": 2, "caption": ch13Locale("المسبار يقرأ السبب.", "Probe reads the reason.")},
			map[string]any{"line": 3, "caption": ch13Locale("السبب يرجع للقيادة.", "The reason returns to command.")},
		},
		"labels": map[string]any{
			"previous": ch13Locale("اللي فات", "Previous"),
			"next":     ch13Locale("اللي بعده", "Next"),
			"reset":    ch13Locale("من الأول", "Reset"),
			"reason":   ch13Locale("سبب الإلغاء", "Stop reason"),
			"canceled": ch13Locale("context.Canceled", "context.Canceled"),
			"deadline": ch13Locale("context.DeadlineExceeded", "context.DeadlineExceeded"),
			"context":  ch13Locale("السياق", "Context"),
			"same":     ch13Locale("نفس ctx", "same ctx"),
			"err":      ch13Locale("Err", "Err"),
			"nil":      ch13Locale("nil", "nil"),
		},
	})
}

const (
	ch13ContextCode = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println(ctx.Err() == nil)
	cancel()
	<-ctx.Done()
	fmt.Println(ctx.Err() == context.Canceled)
}`
	ch13PredictCancelCode = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done()
	fmt.Println(ctx.Err() == context.Canceled)
}`
	ch13TimeoutCode = `func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	<-ctx.Done()
	fmt.Println(errors.Is(ctx.Err(), context.DeadlineExceeded))
}`
	ch13PredictTimeoutCode = `func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	fmt.Println(ctx.Err() == nil)
}`
	ch13CancelCode = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	finished := make(chan bool, 1)
	go func() {
		<-ctx.Done()
		finished <- true
	}()
	cancel()
	<-finished
	fmt.Println(ctx.Err() == context.Canceled)
}`
	ch13PredictDoneCode = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done()
	fmt.Println(ctx.Err())
}`
	ch13ChainCode = `func probe(ctx context.Context) error { return ctx.Err() }
func relay(ctx context.Context) error { return probe(ctx) }
func command(ctx context.Context) error { return relay(ctx) }
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := command(ctx)
	fmt.Println(err == context.Canceled)
}`
	ch13PredictChainCode = `func relay(ctx context.Context) error { return probe(ctx) }
func probe(ctx context.Context) error { return ctx.Err() }
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := relay(ctx)
	fmt.Println(err == context.Canceled)
}`
)

var ch131RecallArabic = []LessonStep{
	{ID: "c131.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch13Markdown("c131.mission", `شغّل 3 drones لحد ما القاعدة تبعت إشارة الرجوع.

الـstarter جهّز ~colony.Mission()~ والـWaitGroup؛ function ~work~ لازم تلاحظ الإلغاء وترجّع الـdrones.

تابع كل drone لوحده، وبعدها استنى المجموعة.`)}},
	{ID: "c131.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch13Walkthrough("c131.context-trace", `الـcontext إشارة إلغاء مشتركة. الـchannel ~Done~ بيتقفل لما الـcontext يتلغي، و~Err~ بيوضح السبب.

الإلغاء مش بيوقف function بالقوة؛ الشغل لازم يلاحظه.`, ch13ContextCode, "e131-context", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 2, Caption: "اعمل context قابل للإلغاء."},
		{Line: 3, Caption: "الـcontext لسه شغال.", Output: "true"},
		{Line: 4, Caption: "ابعت إشارة الإلغاء.", Output: "true"},
		{Line: 5, Caption: "استنى إغلاق Done.", Output: "true"},
		{Line: 6, Caption: "الحالة بقت canceled.", Output: "true\ntrue"},
	})}},
	{ID: "c131.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch13Question("c131.predict-cancel", "بعد ما تلغي وتستنى ~Done~، إيه نتيجة المقارنة؟", ch13PredictCancelCode, []string{"true", "false", "nil"}, 0, []string{"السبب هنا هو context.Canceled.", "الـcontext اتلغى بالفعل.", "المقارنة بترجع bool."}),
		ch13Walkthrough("c131.cancel-reveal", `تابع الإلغاء ثم فحص السبب.`, ch13PredictCancelCode, "e131-cancel", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "جهّز الإشارة."},
			{Line: 3, Caption: "الغِ الـcontext."},
			{Line: 4, Caption: "Done اتقفلت."},
			{Line: 5, Caption: "السبب يطابق Canceled.", Output: "true"},
		}),
		ch131RecallSandbox("بدّل الترتيب، وراقب استجابة الدرونات لإشارة الرجوع."),
	}},
	{ID: "c131.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch13Markdown("c131.watch", `كل goroutine لازم تفحص الإلغاء بنفسها.

~WaitGroup~ تستنى المهام المسجلة؛ ما بتلغيهاش.

أي ترتيب بدء بين goroutines ممكن؛ متفترضش ترتيب ثابت.`)}},
	{ID: "c131.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch13Markdown("c131.plan", `بص على مكان كل drone، وخلي الشغل يقف بعد ما الإشارة توصل.

حدد رجوع كل واحد للقاعدة، وبعدها خلي ~main~ تستنى الكل.

التقاط الإلغاء لوحده مش كفاية لو الـdrone ما رجعش.`),
		ch13Question("c131.q-done", "إيه اللي يحصل لما الـcontext يتلغي؟", "<-ctx.Done()", []string{"القناة Done تتقفل", "كل goroutine تنتهي بالقوة", "WaitGroup تزيد العداد"}, 0, []string{"إغلاق الـchannel بيوصل الإشارة للمستمعين.", "كل مهمة لازم تتصرف بنفسها.", "WaitGroup عدّاد انتظار لوحده."}),
	}},
	{ID: "c131.turn-step", Title: "دورك", Blocks: []LessonBlock{ch13Markdown("c131.turn", `كمّل استجابة ~work~ للإلغاء والرجوع. شغّل المحاكاة وراجع أماكن الـdrones بعد الانتظار.`)}},
}

var ch131RecallEnglish = []LessonStep{
	{ID: "c131.mission-step", Title: "Mission", Blocks: []LessonBlock{ch13Markdown("c131.mission", `Run three drones until base sends a recall order.

The starter sets up ~colony.Mission()~ and the WaitGroup; ~work~ must react to cancellation and return its drone.

Track each drone, then wait for the group.`)}},
	{ID: "c131.idea-step", Title: "Idea", Blocks: []LessonBlock{ch13Walkthrough("c131.context-trace", `A context is a shared cancellation signal. Its ~Done~ channel closes on cancellation, and ~Err~ gives the reason.

Cancellation does not force-stop a function; the work must notice it.`, ch13ContextCode, "e131-context", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create a cancelable context."},
		{Line: 3, Caption: "The context is active.", Output: "true"},
		{Line: 4, Caption: "Send the cancel signal.", Output: "true"},
		{Line: 5, Caption: "Wait for Done to close.", Output: "true"},
		{Line: 6, Caption: "The state is canceled.", Output: "true\ntrue"},
	})}},
	{ID: "c131.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch13Question("c131.predict-cancel", "After cancel and waiting on ~Done~, what is the comparison?", ch13PredictCancelCode, []string{"true", "false", "nil"}, 0, []string{"The reason is context.Canceled.", "The context has already been canceled.", "The comparison returns a bool."}),
		ch13Walkthrough("c131.cancel-reveal", `Follow cancellation, then inspect its reason.`, ch13PredictCancelCode, "e131-cancel", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Prepare the signal."},
			{Line: 3, Caption: "Cancel the context."},
			{Line: 4, Caption: "Done is closed."},
			{Line: 5, Caption: "The reason matches Canceled.", Output: "true"},
		}),
		ch131RecallSandbox("Shuffle workers; watch each respond to recall."),
	}},
	{ID: "c131.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch13Markdown("c131.watch", `Each goroutine must check for cancellation itself.

~WaitGroup~ waits for registered work; it does not cancel it.

Goroutines may start in any order; do not assume a fixed order.`)}},
	{ID: "c131.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch13Markdown("c131.plan", `Check each drone's position and stop its work after the signal arrives.

Plan how each one returns to base, then let ~main~ wait for everyone.

Noticing cancellation is not enough if a drone never returns.`),
		ch13Question("c131.q-done", "What happens when a context is canceled?", "<-ctx.Done()", []string{"The Done channel closes", "Every goroutine is force-stopped", "The WaitGroup count increases"}, 0, []string{"Closing the channel signals its listeners.", "Each task must respond on its own.", "WaitGroup is a separate wait counter."}),
	}},
	{ID: "c131.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch13Markdown("c131.turn", `Complete ~work~ cancellation and return handling. Run the simulation and check drone positions after waiting.`)}},
}

var ch132TimeoutArabic = []LessonStep{
	{ID: "c132.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch13Markdown("c132.mission", `بص على كل كريستالة بـ~DeepScan~، واحصد بس اللي نقاوتها تحقق الشرط.

بعض الفحوصات ممكن تفضل معلقة؛ الـstarter بيمرر لها context من غير timeout.

هدفك تفرّق بين الفحص اللي خلص واللي وقته انتهى.`)}},
	{ID: "c132.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch13Walkthrough("c132.timeout-trace", `~WithTimeout~ بتعمل child context بموعد. لو الوقت خلص، ~Err~ بيرجّع ~DeadlineExceeded~.

استعمل ~colony.Ticks(n)~ عشان تحدد المهلة بوحدة mission ticks.

هنا المدة صفر عشان نوضح إن الموعد خلص فورًا؛ دي مش إعداد للمهمة.`, ch13TimeoutCode, "e132-timeout", []string{"context", "errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 2, Caption: "اعمل context موعده فوري."},
		{Line: 3, Caption: "سجّل تنظيف الـcontext."},
		{Line: 4, Caption: "استنى الـcontext يخلص."},
		{Line: 5, Caption: "السبب هو انتهاء المهلة.", Output: "true"},
	})}},
	{ID: "c132.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch13Question("c132.predict-active", "مع timeout لساعة ما خلصتش، إيه نتيجة الشرط؟", ch13PredictTimeoutCode, []string{"true", "false", "context canceled"}, 0, []string{"الـcontext لسه مفيهوش error.", "المهلة لسه قدام.", "cancel المؤجل لسه ما اشتغلش."}),
		ch13Walkthrough("c132.active-reveal", `بص على الـcontext قبل ما موعده يخلص.`, ch13PredictTimeoutCode, "e132-active", []string{"context", "fmt", "time"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "حدد موعدًا في المستقبل."},
			{Line: 3, Caption: "التنظيف مؤجل للخروج."},
			{Line: 4, Caption: "Err ما زال nil.", Output: "true"},
		}),
		ch132TimeoutSandbox("قدّم الوقت أو ألغِ يدوي؛ قارن أخطاء الـcontext."),
	}},
	{ID: "c132.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch13Markdown("c132.watch", `كل عملية طويلة محتاجة child context مناسب.

استعمل ~cancel~ بعد العملية عشان تحرر موارد الـcontext.

الـparent يقدر يلغي الـchild قبل موعده.`)}},
	{ID: "c132.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch13Markdown("c132.plan", `افصل scan عن قرار الحصاد. النقاوة تهم بس بعد ما scan ترجع بنجاح.

لو الوقت خلص، سيب البلورة وانتقل للي بعدها؛ ما تخليش فحص واحد يوقف الجولة.

استعمل سبب الـerror عشان تفرّق بين انتهاء المهلة وأي فشل تاني.`),
		ch13Question("c132.q-timeout", "أنهي error معناه إن وقت الـcontext خلص؟", "ctx.Err()", []string{"context.DeadlineExceeded", "context.Canceled فقط", "nil دائمًا"}, 0, []string{"DeadlineExceeded بيميز انتهاء الموعد.", "الإلغاء اليدوي ليه سبب تاني.", "nil معناها إن الـcontext لسه شغال."}),
	}},
	{ID: "c132.turn-step", Title: "دورك", Blocks: []LessonBlock{ch13Markdown("c132.turn", `ضيف timeout لكل scan وتعامل مع النتيجة. شغّل المهمة وراقب إن الفحوص المعلقة ما توقفش المسح كله.`)}},
}

var ch132TimeoutEnglish = []LessonStep{
	{ID: "c132.mission-step", Title: "Mission", Blocks: []LessonBlock{ch13Markdown("c132.mission", `Run ~DeepScan~ on each crystal and harvest only those meeting the purity rule.

Some scans can hang; the starter passes a context without a deadline.

Distinguish a completed scan from one that reached its deadline.`)}},
	{ID: "c132.idea-step", Title: "Idea", Blocks: []LessonBlock{ch13Walkthrough("c132.timeout-trace", `~WithTimeout~ creates a child context with a deadline. If time expires, ~Err~ returns ~DeadlineExceeded~.

Use ~colony.Ticks(n)~ to express a timeout in mission ticks.

Zero duration shows an already-expired deadline; it is not a mission setting.`, ch13TimeoutCode, "e132-timeout", []string{"context", "errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create an immediate deadline."},
		{Line: 3, Caption: "Register context cleanup."},
		{Line: 4, Caption: "Wait for the context to end."},
		{Line: 5, Caption: "The reason is deadline exceeded.", Output: "true"},
	})}},
	{ID: "c132.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch13Question("c132.predict-active", "With an hour left, what does the condition print?", ch13PredictTimeoutCode, []string{"true", "false", "context canceled"}, 0, []string{"The context has no error yet.", "The future deadline has not expired.", "The deferred cancel has not run."}),
		ch13Walkthrough("c132.active-reveal", `Inspect the context before its deadline.`, ch13PredictTimeoutCode, "e132-active", []string{"context", "fmt", "time"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Set a future deadline."},
			{Line: 3, Caption: "Cleanup is deferred."},
			{Line: 4, Caption: "Err is still nil.", Output: "true"},
		}),
		ch132TimeoutSandbox("Advance or cancel; compare context errors."),
	}},
	{ID: "c132.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch13Markdown("c132.watch", `Give each long operation an appropriate child context.

Call ~cancel~ when it ends to release the context's resources.

The parent can cancel a child before its deadline.`)}},
	{ID: "c132.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch13Markdown("c132.plan", `Separate scanning from the harvest decision. Purity matters only after a scan succeeds.

If its deadline expires, leave that crystal and continue; one scan should not stop the whole pass.

Use the error reason to distinguish a timeout from another failure.`),
		ch13Question("c132.q-timeout", "Which error marks an expired context deadline?", "ctx.Err()", []string{"context.DeadlineExceeded", "Only context.Canceled", "Always nil"}, 0, []string{"DeadlineExceeded marks the expired deadline.", "Manual cancellation can have another reason.", "nil means the context has not ended."}),
	}},
	{ID: "c132.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch13Markdown("c132.turn", `Add a deadline to each scan and handle its result. Run the mission and check that a hanging scan does not stop the full pass.`)}},
}

var ch133FirstWinsArabic = []LessonStep{
	{ID: "c133.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch13Markdown("c133.mission", `3 drones بيدوروا على كريستالة أسطورية. أول drone يلاقيها يلغي بحث الباقي.

الـstarter وزّع اتنين drones كمان وحط هيكل ~search~.

بعد ما البحث يخلص، الكل يرجع للقاعدة.`)}},
	{ID: "c133.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch13Walkthrough("c133.cancel-trace", `الـchannel Done بتوصل الإلغاء لكل العمال. المثال بيستنى رسالة من goroutine بعد الإلغاء.

الـframes بتعرض ترتيب ممكن؛ تشغيل goroutines مالوش ترتيب مضمون.`, ch13CancelCode, "e133-cancel", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 2, Caption: "اعمل context وchannel للنتيجة."},
		{Line: 3, Caption: "الـchannel بتستقبل إشارة الـworker."},
		{Line: 4, Caption: "العامل ينتظر Done."},
		{Line: 8, Caption: "الغِ الـcontext للجميع."},
		{Line: 9, Caption: "استلم رسالة إن الـworker خلص."},
		{Line: 10, Caption: "سبب الإلغاء صحيح.", Output: "true"},
	})}},
	{ID: "c133.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch13Question("c133.predict-done", "بعد الإلغاء، إيه اللي بيرجّعه ~ctx.Err()~؟", ch13PredictDoneCode, []string{"context canceled", "nil", "deadline exceeded"}, 0, []string{"الإلغاء يدوي في المثال.", "nil معناها إن الـcontext لسه شغال.", "مفيش deadline في الكود."}),
		ch13Walkthrough("c133.done-reveal", `تابع الإلغاء وبعده اقرا السبب.`, ch13PredictDoneCode, "e133-done", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "اعمل الـcontext."},
			{Line: 3, Caption: "ابعت الإلغاء."},
			{Line: 4, Caption: "Done تقفل فورًا."},
			{Line: 5, Caption: "اقرا سبب الإلغاء.", Output: "context canceled"},
		}),
		ch133RaceSandbox("بدّل البحث؛ تابع الفائز واستجابة الباقي."),
	}},
	{ID: "c133.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch13Markdown("c133.watch", `الـworker لازم يفحص الـcontext خلال البحث.

إلغاء الـcontext ما بيرجّعش الـdrone للقاعدة لوحده.

كل goroutine لازم تخلص قبل ما البرنامج يقفل.`)}},
	{ID: "c133.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch13Markdown("c133.plan", `خلي كل worker يشارك نفس الـcontext ويعرف إمتى يقف.

اربط نجاح البحث بإشارة الإلغاء، وبعدها افصل الرجوع للقاعدة عن الشغل.

استعمل ~WaitGroup~ عشان تتأكد إن كل العمال خلصوا.`),
		ch13Question("c133.q-cancel", "~cancel()~ بتعمل إيه في بحث مشترك؟", "cancel()", []string{"تقفل Done وتبلّغ العمال", "توقف كل goroutine بالقوة", "تنتظر كل العمال"}, 0, []string{"إغلاق Done بيصحّي اللي مستني الإشارة.", "العمال بيخرجوا بالتعاون.", "الانتظار مسؤولية WaitGroup."}),
	}},
	{ID: "c133.turn-step", Title: "دورك", Blocks: []LessonBlock{ch13Markdown("c133.turn", `أكمل فحص الإلغاء في ~search~. شغّل المهمة وراقب استجابة كل درون بعد العثور على الهدف.`)}},
}

var ch133FirstWinsEnglish = []LessonStep{
	{ID: "c133.mission-step", Title: "Mission", Blocks: []LessonBlock{ch13Markdown("c133.mission", `Three drones search for a legendary crystal. The first finder cancels the others' search.

The starter spawns two drones and provides the ~search~ outline.

After searching ends, everyone returns to base.`)}},
	{ID: "c133.idea-step", Title: "Idea", Blocks: []LessonBlock{ch13Walkthrough("c133.cancel-trace", `The Done channel sends cancellation to every worker. This example waits for a goroutine message after cancellation.

Frames show one possible order; goroutine scheduling has no guaranteed order.`, ch13CancelCode, "e133-cancel", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create context and result channel."},
		{Line: 3, Caption: "The channel receives a worker signal."},
		{Line: 4, Caption: "The worker waits on Done."},
		{Line: 8, Caption: "Cancel the shared context."},
		{Line: 9, Caption: "Receive worker completion."},
		{Line: 10, Caption: "The cancellation reason matches.", Output: "true"},
	})}},
	{ID: "c133.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch13Question("c133.predict-done", "After cancellation, what does ~ctx.Err()~ return?", ch13PredictDoneCode, []string{"context canceled", "nil", "deadline exceeded"}, 0, []string{"The example uses manual cancellation.", "nil means the context is still active.", "The code has no deadline."}),
		ch13Walkthrough("c133.done-reveal", `Follow cancellation and read the reason.`, ch13PredictDoneCode, "e133-done", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Create the context."},
			{Line: 3, Caption: "Send cancellation."},
			{Line: 4, Caption: "Done closes immediately."},
			{Line: 5, Caption: "Read the cancellation reason.", Output: "context canceled"},
		}),
		ch133RaceSandbox("Shuffle searches; watch the winner and worker response."),
	}},
	{ID: "c133.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch13Markdown("c133.watch", `Each worker must check the context during its search.

Canceling a context does not return a drone to base automatically.

Every goroutine must finish before the program exits.`)}},
	{ID: "c133.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch13Markdown("c133.plan", `Have each worker share the same context and notice when to stop.

Connect a successful search to cancellation, then keep returning to base separate from searching.

Use ~WaitGroup~ to ensure every worker has finished.`),
		ch13Question("c133.q-cancel", "What does ~cancel()~ do in a shared search?", "cancel()", []string{"Closes Done and signals workers", "Force-stops every goroutine", "Waits for every worker"}, 0, []string{"Closing Done wakes listeners.", "Workers exit cooperatively.", "WaitGroup handles waiting."}),
	}},
	{ID: "c133.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch13Markdown("c133.turn", `Complete the cancellation check in ~search~. Run the mission and watch each drone respond after the target is found.`)}},
}

var ch134ChainArabic = []LessonStep{
	{ID: "c134.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch13Markdown("c134.mission", `مرّر الـcontext من القيادة للـrelay للـprobe، ورجّع سبب التوقف للقيادة.

الـstarter فيه سلسلة النداءات؛ الـrelay بيمرر context جديد بالغلط.

سجّل دخول وخروج كل مستوى في المرتين.`)}},
	{ID: "c134.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch13Walkthrough("c134.chain-trace", `كل function بتستقبل نفس الـcontext كأول parameter وتمرره للداخل.

الـprobe بيقرا ~Err~، والنداءات الخارجية بترجّع السبب زي ما هو.

~colony.TraceCommand(ctx, stage)~ بيسجل حالة الـcontext عند حد الـfunction.`, ch13ChainCode, "e134-chain", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 1, Caption: "الـprobe بيقرا حالة الـcontext."},
		{Line: 2, Caption: "الـrelay بيمرر الـcontext."},
		{Line: 3, Caption: "القيادة بتمرر الـcontext."},
		{Line: 5, Caption: "اعمل context القيادة."},
		{Line: 6, Caption: "الغِ الـcontext قبل النداء."},
		{Line: 7, Caption: "نادِ السلسلة."},
		{Line: 8, Caption: "السبب وصل للقيادة.", Output: "true"},
	})}},
	{ID: "c134.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch13Question("c134.predict-chain", "حالة الإلغاء بتوصل عبر ~relay~؟", ch13PredictChainCode, []string{"true", "false", "nil"}, 0, []string{"الـrelay مرّرت نفس ctx للـprobe.", "الـprobe بيرجّع Err من ctx الملغي.", "المقارنة بترجّع bool."}),
		ch13Walkthrough("c134.chain-reveal", `تابع الـcontext من الـrelay للـprobe.`, ch13PredictChainCode, "e134-relay", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 4, Caption: "اعمل context."},
			{Line: 5, Caption: "الغيه قبل السلسلة."},
			{Line: 6, Caption: "الـrelay بتنادي الـprobe."},
			{Line: 2, Caption: "الـprobe بترجّع Err."},
			{Line: 7, Caption: "السبب حافظ على هويته.", Output: "true"},
		}),
		ch134ChainSandbox("بدّل السبب؛ تتبّع الـcontext عبر السلسلة."),
	}},
	{ID: "c134.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch13Markdown("c134.watch", `إنشاء ~context.Background()~ في النص بيقطع إلغاء الـparent.

الـdeadline بينتقل للـchildren، ومينفعش تمدده من مستوى داخلي.

لو الـcontext خلص، ما ترجّعش ~nil~ بدل السبب.`)}},
	{ID: "c134.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch13Markdown("c134.plan", `راجع كل حد للـcall: الدخول، نداء المستوى اللي بعده، رجوع النتيجة، وبعدها الخروج.

جرّب الـcontext الملغي والـcontext بموعد كحالتين منفصلتين.

رتّب السجل قبل ما تعلن النتيجة النهائية.`),
		ch13Question("c134.q-propagate", "إيه اللي بيحافظ على إلغاء القيادة؟", "next(ctx)", []string{"تمرير نفس ctx", "إنشاء Background جديد", "حفظ ctx في struct"}, 0, []string{"نفس الـcontext شايل الإلغاء والموعد.", "Background بيبدأ عمر منفصل.", "مرّر الـcontext لكل نداء."}),
	}},
	{ID: "c134.turn-step", Title: "دورك", Blocks: []LessonBlock{ch13Markdown("c134.turn", `صلّح تمرير الـcontext وسجّل حدود الـfunctions. شغّل الحالتين وراجع سبب التوقف في السجل.`)}},
}

var ch134ChainEnglish = []LessonStep{
	{ID: "c134.mission-step", Title: "Mission", Blocks: []LessonBlock{ch13Markdown("c134.mission", `Pass context from command to relay to probe, and return the stop reason to command.

The starter has the call chain; relay incorrectly passes a new context.

Trace every layer's entry and exit in both runs.`)}},
	{ID: "c134.idea-step", Title: "Idea", Blocks: []LessonBlock{ch13Walkthrough("c134.chain-trace", `Each function receives the same context as its first parameter and passes it inward.

The probe reads ~Err~; outer calls return the reason without replacing it.

~colony.TraceCommand(ctx, stage)~ records context state at a function boundary.`, ch13ChainCode, "e134-chain", []string{"context", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "The probe reads context state."},
		{Line: 2, Caption: "The relay passes context."},
		{Line: 3, Caption: "The command passes context."},
		{Line: 5, Caption: "Create command context."},
		{Line: 6, Caption: "Cancel before the call."},
		{Line: 7, Caption: "Call through the chain."},
		{Line: 8, Caption: "The reason reaches command.", Output: "true"},
	})}},
	{ID: "c134.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch13Question("c134.predict-chain", "Does cancellation pass through ~relay~?", ch13PredictChainCode, []string{"true", "false", "nil"}, 0, []string{"relay passes the same ctx to probe.", "The probe returns Err from canceled ctx.", "The comparison returns a bool."}),
		ch13Walkthrough("c134.chain-reveal", `Follow context from relay to probe.`, ch13PredictChainCode, "e134-relay", []string{"context", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 4, Caption: "Create a context."},
			{Line: 5, Caption: "Cancel before the chain."},
			{Line: 6, Caption: "relay calls probe."},
			{Line: 2, Caption: "probe returns Err."},
			{Line: 7, Caption: "The reason keeps its identity.", Output: "true"},
		}),
		ch134ChainSandbox("Switch reasons; trace context through the chain."),
	}},
	{ID: "c134.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch13Markdown("c134.watch", `Creating a new ~context.Background()~ breaks parent cancellation.

A deadline flows to children and cannot be extended by an inner layer.

Do not return ~nil~ after the context ends.`)}},
	{ID: "c134.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch13Markdown("c134.plan", `Check each call boundary: entry, call the next layer, return its result, then exit.

Exercise the already-canceled context and the deadline context separately.

Review trace order before the final report.`),
		ch13Question("c134.q-propagate", "What preserves command's cancellation?", "next(ctx)", []string{"Pass the same ctx", "Create a new Background", "Store ctx in a struct"}, 0, []string{"The same context carries cancellation and deadline.", "Background starts a separate lifetime.", "Pass context to each call."}),
	}},
	{ID: "c134.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch13Markdown("c134.turn", `Fix context propagation and trace function boundaries. Run both cases and check the stop reason in the log.`)}},
}
