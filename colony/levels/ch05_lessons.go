package levels

func ch05Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch05Question(id, text, code string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text + "\n\n~~~go\n" + code + "\n~~~",
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch05Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

func ch05Pair(ar, en string) map[string]any {
	return map[string]any{"ar": ar, "en": en}
}

func ch05Sandbox(id, text, kind string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch05ParamsSandboxConfig() map[string]any {
	return map[string]any{
		"source": ch05Pair("قيم للتوضيح", "Illustrative values"),
		"labels": map[string]any{
			"previous": ch05Pair("السابق", "Previous"), "next": ch05Pair("التالي", "Next"), "reset": ch05Pair("من الأول", "Reset"),
			"call": ch05Pair("النداء", "Call"), "inside": ch05Pair("جوه الدالة", "Inside function"), "result": ch05Pair("النتيجة", "Result"),
		},
		"frames": []any{
			map[string]any{"caption": ch05Pair("النداء يمرر 2.", "The call passes 2."), "call": "triple(2)", "inside": "—", "result": "—"},
			map[string]any{"caption": ch05Pair("n تستقبل القيمة.", "n receives the value."), "call": "triple(2)", "inside": "2", "result": "—"},
			map[string]any{"caption": ch05Pair("الـreturn ترجع 6.", "return carries 6."), "call": "triple(2)", "inside": "2", "result": "6"},
		},
	}
}

func ch05ErrorValueSandboxConfig() map[string]any {
	return map[string]any{
		"source": ch05Pair("قراءة توضيحية", "Illustrative read"),
		"labels": map[string]any{
			"outcome": ch05Pair("حالة error", "Error state"), "present": ch05Pair("موجود", "Present"), "absent": ch05Pair("nil", "nil"),
			"value": ch05Pair("القيمة الراجعة", "Returned value"), "error": ch05Pair("الخطأ", "Error"), "action": ch05Pair("التصرف", "Next step"),
		},
		"success": map[string]any{
			"error": ch05Pair("nil", "nil"), "action": ch05Pair("استخدم القيمة 12.", "Use value 12."), "caption": ch05Pair("مفيش خطأ؛ النتيجة صالحة.", "No error; the result is usable."),
		},
		"failure": map[string]any{
			"error": ch05Pair("sensor", "sensor"), "action": ch05Pair("افحص الخطأ الأول.", "Check the error first."), "caption": ch05Pair("فيه خطأ؛ سيب القيمة.", "An error is present; skip the value."),
		},
		"value": 12,
	}
}

func ch05WrapSandboxConfig() map[string]any {
	return map[string]any{
		"source": ch05Pair("سبب توضيحي", "Illustrative cause"),
		"labels": map[string]any{
			"format": ch05Pair("صيغة التنسيق", "Format"), "wrap": "%w", "plain": "%v",
			"message": ch05Pair("الرسالة", "Message"), "identity": ch05Pair("errors.Is", "errors.Is"),
		},
		"wrapped": map[string]any{
			"message": "bay 2: valve", "identity": "true", "caption": ch05Pair("السبب لسه معروف.", "The cause stays identifiable."),
		},
		"plain": map[string]any{
			"message": "bay 2: valve", "identity": "false", "caption": ch05Pair("النص موجود؛ الهوية ضاعت.", "Text remains; identity is lost."),
		},
	}
}

func ch05DeferSandboxConfig() map[string]any {
	return map[string]any{
		"source": ch05Pair("حروف تدريبية", "Sample markers"),
		"labels": map[string]any{
			"previous": ch05Pair("السابق", "Previous"), "next": ch05Pair("التالي", "Next"), "reset": ch05Pair("من الأول", "Reset"),
			"output": ch05Pair("الناتج", "Output"),
		},
		"frames": []any{
			map[string]any{"caption": ch05Pair("مفيش طباعة لسه.", "No output yet."), "output": "—"},
			map[string]any{"caption": ch05Pair("النداء الفوري طبع K.", "The immediate call prints K."), "output": "K"},
			map[string]any{"caption": ch05Pair("آخر defer طبع J.", "The last defer prints J."), "output": "KJ"},
			map[string]any{"caption": ch05Pair("أول defer طبع H.", "The first defer prints H."), "output": "KJH"},
		},
	}
}

func ch05ClosureSandboxConfig() map[string]any {
	return map[string]any{
		"source": ch05Pair("أمر تدريبي", "Sample command"),
		"labels": map[string]any{
			"store": ch05Pair("خزّن الدالة", "Store function"), "call": ch05Pair("نفّذ النداء", "Call function"), "reset": ch05Pair("من الأول", "Reset"),
			"count": ch05Pair("عدد النداءات", "Call count"),
		},
		"states": map[string]any{
			"idle":   ch05Pair("لسه مفيش function value.", "No function value stored yet."),
			"stored": ch05Pair("اتخزنت؛ لسه ما اتنفذتش.", "Stored; nothing has run yet."),
			"called": ch05Pair("الـclosure حدّثت العداد.", "The closure updated its counter."),
		},
		"initialCount": 0,
	}
}

const (
	ch05ParamsCode = `func twice(n int) int {
	return n * 2
}
func main() {
	value := twice(4)
	fmt.Println(value)
}`
	ch05PredictParamsCode = `func total(a, b int) int {
	return a + b
}
func main() {
	fmt.Println(total(3, 4))
}`
	ch05ErrorCode = `func sample() (int, error) {
	return 6, errors.New("storm")
}
func main() {
	n, err := sample()
	if err != nil {
		fmt.Println("retry")
		return
	}
	fmt.Println(n)
}`
	ch05PredictErrorCode = `func sample() (int, error) {
	return 6, nil
}
func main() {
	n, err := sample()
	if err != nil {
		fmt.Println("retry")
	} else {
		fmt.Println(n)
	}
}`
	ch05WrapCode = `func main() {
	cause := errors.New("valve")
	wrapped := fmt.Errorf("bay %d: %w", 2, cause)
	fmt.Println(errors.Is(wrapped, cause))
	fmt.Println(wrapped)
}`
	ch05PredictWrapCode = `func main() {
	cause := errors.New("sensor")
	wrapped := fmt.Errorf("unit: %w", cause)
	fmt.Println(errors.Is(wrapped, cause))
}`
	ch05DeferCode = `func work() {
	defer fmt.Println("closed")
	fmt.Println("working")
}
func main() {
	work()
}`
	ch05PredictDeferCode = `func work() {
	defer fmt.Print("A")
	defer fmt.Print("B")
	fmt.Print("C")
}
func main() {
	work()
	fmt.Println()
}`
	ch05FunctionValueCode = `func ping() { fmt.Print("P") }
func main() {
	action := ping
	action()
	fmt.Println("!")
}`
	ch05PredictClosureCode = `func main() {
	hits := 0
	count := func() { hits++ }
	count()
	count()
	fmt.Println(hits)
}`
)

var ch05FunctionsArabic = []LessonStep{
	{ID: "c51.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch05Markdown("c51.mission", `احصد الطحالب في الصفوف التلاتة باستخدام function هتكتبها.

الـstarter جهّز ~harvestRow(n, dir)~ واستدعاها مرة.

	الهدف 24 طحلب؛ ركّز على وظيفة الـparameters وتقسيم الشغل بين الـfunctions.`)}},
	{ID: "c51.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch05Walkthrough("c51.params-trace", `الـparameter هو value بتدخل الـfunction. الـfunction ممكن ترجّع نتيجة بنوع واضح.

تابع ~n~ من النداء، وبعدها شوف الـvalue اللي رجعت للـcaller.`, ch05ParamsCode, "e51-params", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 1, Caption: "الـfunction بتستقبل رقم."},
		{Line: 5, Caption: "twice تستقبل 4.", Variables: []LessonVariable{{Name: "n", Value: "4"}}},
		{Line: 2, Caption: "الـfunction بترجّع 8.", Variables: []LessonVariable{{Name: "n", Value: "4"}}},
		{Line: 6, Caption: "اطبع الـvalue اللي رجعت.", Variables: []LessonVariable{{Name: "value", Value: "8"}}, Output: "8"},
	})}},
	{ID: "c51.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch05Question("c51.predict-result", "توقّع الناتج من النداء:", ch05PredictParamsCode, []string{"12", "7", "34"}, 1, []string{"الـfunction مش بتستخدم الضرب.", "الـfunction بتجمع الـvalues: 3 و4.", "النداء ما بيطبعش الأرقام جنب بعض."}),
		ch05Walkthrough("c51.result-reveal", `تابع الـvalues وهي داخلة للـfunction وراجعة.`, ch05PredictParamsCode, "e51-total", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 5, Caption: "النداء يمرر 3 و4."},
			{Line: 2, Caption: "الجمع بيرجّع 7."},
			{Line: 6, Caption: "الناتج المطبوع 7.", Output: "7"},
		}),
		ch05Sandbox("c51.params-sandbox", "تابع value للتوضيح وهي داخلة ~function~ وراجعة.", "ch05-params", ch05ParamsSandboxConfig()),
	}},
	{ID: "c51.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch05Markdown("c51.watch", `الـparameter جوه الـfunction ليه الاسم والنوع المكتوبين في التوقيع.

النداء لازم يمرر النوع المناسب؛ والـreturn لازم يطابق نوع النتيجة.

	خلّي كل function مسؤولة عن جزء مفهوم.`)}},
	{ID: "c51.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch05Markdown("c51.plan", `ابدأ بتحديد شغل الصف الواحد، وبعدها قرر إيه اللي يتكرر وإيه اللي يتغير بين النداءات.

خلّي ترتيب الانتقال واضح في ~main~. راجع اتجاه الدرون بعد كل صف.

		ما تخلطش بين عدد parameters الـfunction وعدد مرات استدعائها.`),
		ch05Question("c51.q-parameter", "لو عايز الـfunction تتعامل مع طول مختلف، أنهي تعديل يعبّر عن ده؟", "func walk(n int)", []string{"إضافة parameter للطول", "طباعة n بدل استخدامها", "تغيير اسم الدالة فقط"}, 0, []string{"الـparameter بينقل الـvalue للـfunction.", "الطباعة مش بتغيّر سلوك الـfunction.", "اسم الـfunction مش بيمرر value."}),
	}},
	{ID: "c51.turn-step", Title: "دورك", Blocks: []LessonBlock{ch05Markdown("c51.turn", `كمّل جسم الـfunction، وبعدها رتّب النداءات في ~main~. شغّل المهمة وراقب الحصاد بعد كل انتقال.`)}},
}

var ch05FunctionsEnglish = []LessonStep{
	{ID: "c51.mission-step", Title: "Mission", Blocks: []LessonBlock{ch05Markdown("c51.mission", `Harvest moss in three rows using a function you write.

The starter declares ~harvestRow(n, dir)~ and calls it once.

The goal is 24 moss; focus on parameters and splitting work across functions.`)}},
	{ID: "c51.idea-step", Title: "Idea", Blocks: []LessonBlock{ch05Walkthrough("c51.params-trace", `A parameter is a value passed into a function. A function can return a result with a declared type.

Follow ~n~ from the call, then watch the returned value reach its caller.`, ch05ParamsCode, "e51-params", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "The function accepts a number."},
		{Line: 5, Caption: "twice receives 4.", Variables: []LessonVariable{{Name: "n", Value: "4"}}},
		{Line: 2, Caption: "The result returns as 8.", Variables: []LessonVariable{{Name: "n", Value: "4"}}},
		{Line: 6, Caption: "Print the returned value.", Variables: []LessonVariable{{Name: "value", Value: "8"}}, Output: "8"},
	})}},
	{ID: "c51.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch05Question("c51.predict-result", "Predict the output of this call:", ch05PredictParamsCode, []string{"12", "7", "34"}, 1, []string{"The function does not multiply.", "The function adds 3 and 4.", "The call does not print the numbers side by side."}),
		ch05Walkthrough("c51.result-reveal", `Follow the values into the function and back.`, ch05PredictParamsCode, "e51-total", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 5, Caption: "The call passes 3 and 4."},
			{Line: 2, Caption: "The sum returns 7."},
			{Line: 6, Caption: "The printed output is 7.", Output: "7"},
		}),
		ch05Sandbox("c51.params-sandbox", "Follow a sample value into a ~function~ and back.", "ch05-params", ch05ParamsSandboxConfig()),
	}},
	{ID: "c51.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch05Markdown("c51.watch", `A parameter inside a function has the name and type from its signature.

Calls must pass matching types; a return must match the result type.

Give each function one clear responsibility.`)}},
	{ID: "c51.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch05Markdown("c51.plan", `First identify the work for one row, then decide what repeats and what changes between calls.

Keep transitions clear in ~main~. Check the drone's direction after each row.

Do not confuse a function's parameter count with its number of calls.`),
		ch05Question("c51.q-parameter", "To support a different length, which change expresses it?", "func walk(n int)", []string{"Add a length parameter", "Print n instead of using it", "Only rename the function"}, 0, []string{"A parameter carries the value inside.", "Printing does not change the function's behavior.", "A function name does not pass a value."}),
	}},
	{ID: "c51.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch05Markdown("c51.turn", `Complete the function body, then order the calls in ~main~. Run the mission and watch each row's harvest.`)}},
}

var ch05StormArabic = []LessonStep{
	{ID: "c52.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch05Markdown("c52.mission", `اجمع كل الكريستال وسط العاصفة، وبعدها ابعت المجموع.

الـstarter جرّب ~TryHarvest~ مرة، لكنه لسه ما فصلش النجاح عن الفشل.

المهمة بتستخدم error كـvalue راجعة.`)}},
	{ID: "c52.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch05Walkthrough("c52.error-trace", `الـfunction ممكن ترجّع نتيجة و~error~ مع بعض. افحص الـerror قبل ما تعتمد على الـvalue.

هنا العطل للتوضيح؛ مفيش حركة أو بيانات مهمة.`, ch05ErrorCode, "e52-error", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 5, Caption: "sample بترجّع value وerror."},
		{Line: 6, Caption: "فيه error.", Variables: []LessonVariable{{Name: "n", Value: "6"}, {Name: "err", Value: "storm"}}},
		{Line: 7, Caption: "اختار مسار إعادة المحاولة."},
		{Line: 8, Caption: "اطبع سبب المحاولة.", Variables: []LessonVariable{{Name: "err", Value: "storm"}}, Output: "retry"},
	})}},
	{ID: "c52.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch05Question("c52.predict-success", "في المثال ده، إيه اللي هيتطبع؟", ch05PredictErrorCode, []string{"retry", "6", "nil"}, 1, []string{"الـerror nil، فشرط if مش هيدخل.", "المسار الناجح بيطبع الـvalue اللي في n.", "nil مش value مطبوعة في الكود."}),
		ch05Walkthrough("c52.success-reveal", `اتبع شرط الـerror للمسار المناسب.`, ch05PredictErrorCode, "e52-success", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 5, Caption: "sample بترجّع 6 وnil."},
			{Line: 6, Caption: "مفيش error."},
			{Line: 9, Caption: "اطبع الـvalue الناجحة.", Variables: []LessonVariable{{Name: "n", Value: "6"}}, Output: "6"},
		}),
		ch05Sandbox("c52.error-value-sandbox", "بدّل حالة القراءة ولاحظ فحص ~error~ قبل الـvalue.", "ch05-error-value", ch05ErrorValueSandboxConfig()),
	}},
	{ID: "c52.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch05Markdown("c52.watch", `ما تستخدمش الـvalue لو العملية رجّعت error.

~errors.Is~ بتقارن سبب الـerror المعروف حتى لو اتلف برسالة.

إعادة المحاولة لازم تكون لسبب ينفع يتعالج.`)}},
	{ID: "c52.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch05Markdown("c52.plan", `افصل بين محاولة العملية ومعالجة نتيجتها. احتفظ بالمجموع خارج المحاولات.

		ما تكررّش تسجيل الـvalue لو العملية فشلت، وراجع شرط اكتمال الجمع قبل التقرير.

فكر: إيه اللي يتغير بعد نجاح واحد؟`),
		ch05Question("c52.q-error", "لو ~err != nil~، أنهي تصرف آمن؟", "value, err := read()", []string{"افحص سبب الخطأ قبل استخدام value", "اعتبر value ناجحة دائمًا", "تجاهل err"}, 0, []string{"الـerror بيحدد لو الـvalue صالحة.", "العملية ممكن تكون فشلت.", "تجاهل الـerror بيخفي سبب الفشل."}),
	}},
	{ID: "c52.turn-step", Title: "دورك", Blocks: []LessonBlock{ch05Markdown("c52.turn", `رتّب فحص الـerror حول ~TryHarvest~. شغّل المهمة وراقب الفرق بين العاصفة والنجاح.`)}},
}

var ch05StormEnglish = []LessonStep{
	{ID: "c52.mission-step", Title: "Mission", Blocks: []LessonBlock{ch05Markdown("c52.mission", `Collect all the crystal during the storm, then report the total.

The starter calls ~TryHarvest~ once, but does not yet separate success from failure.

This mission returns an error as a value.`)}},
	{ID: "c52.idea-step", Title: "Idea", Blocks: []LessonBlock{ch05Walkthrough("c52.error-trace", `A function can return a result and an ~error~ together. Check the error before trusting the result.

This failure is illustrative; it uses no mission state.`, ch05ErrorCode, "e52-error", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 5, Caption: "sample returns a value and error."},
		{Line: 6, Caption: "An error is present.", Variables: []LessonVariable{{Name: "n", Value: "6"}, {Name: "err", Value: "storm"}}},
		{Line: 7, Caption: "Choose the retry path."},
		{Line: 8, Caption: "Print the retry reason.", Variables: []LessonVariable{{Name: "err", Value: "storm"}}, Output: "retry"},
	})}},
	{ID: "c52.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch05Question("c52.predict-success", "What does this example print?", ch05PredictErrorCode, []string{"retry", "6", "nil"}, 1, []string{"The error is nil, so if is skipped.", "The success path prints n.", "The code does not print nil."}),
		ch05Walkthrough("c52.success-reveal", `Follow the error check to its branch.`, ch05PredictErrorCode, "e52-success", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 5, Caption: "sample returns six and nil."},
			{Line: 6, Caption: "There is no error."},
			{Line: 9, Caption: "Print the successful value.", Variables: []LessonVariable{{Name: "n", Value: "6"}}, Output: "6"},
		}),
		ch05Sandbox("c52.error-value-sandbox", "Toggle the read outcome; check ~error~ before value.", "ch05-error-value", ch05ErrorValueSandboxConfig()),
	}},
	{ID: "c52.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch05Markdown("c52.watch", `Do not use a result when the operation returned an error.

~errors.Is~ matches a known cause even after wrapping.

Retry only when the failure can be handled.`)}},
	{ID: "c52.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch05Markdown("c52.plan", `Separate the operation attempt from result handling. Keep the total outside individual attempts.

Do not record a failed value, and check that collection is complete before reporting.

Ask: what changes after one success?`),
		ch05Question("c52.q-error", "If ~err != nil~, what is a safe next step?", "value, err := read()", []string{"Check the error before using value", "Always treat value as successful", "Ignore err"}, 0, []string{"The error tells whether the result is usable.", "The operation may have failed.", "Ignoring the error hides the failure."}),
	}},
	{ID: "c52.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch05Markdown("c52.turn", `Place the error check around ~TryHarvest~. Run the mission and watch the difference between a storm and success.`)}},
}

var ch05FaultArabic = []LessonStep{
	{ID: "c53.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch05Markdown("c53.mission", `افحص وحدات التشخيص، ابعت كود كل عطل، وبعدها سجل الأسباب معًا.

الـstarter بيعرض عدد الوحدات؛ كمّل على نفس المصدر.

فيه وحدات سليمة، فمش كل نتيجة هتكون error.`)}},
	{ID: "c53.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch05Walkthrough("c53.wrap-trace", `~%w~ بيضيف context وبيحافظ على هوية السبب. عشان كده ~errors.Is~ تقدر تلاقيه بعد اللف.

الـvalues هنا إنذار تدريبي، مش تشخيص المهمة.`, ch05WrapCode, "e53-wrap", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 2, Caption: "اعمل سبب للتوضيح."},
		{Line: 3, Caption: "لف السبب بـcontext الوحدة."},
		{Line: 4, Caption: "Is لسه لاقية السبب.", Variables: []LessonVariable{{Name: "wrapped", Value: "bay 2: valve"}}, Output: "true"},
		{Line: 5, Caption: "النص بيحتفظ بالـcontext.", Output: "true\nbay 2: valve"},
	})}},
	{ID: "c53.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch05Question("c53.predict-cause", "توقّع نتيجة فحص هوية السبب:", ch05PredictWrapCode, []string{"true", "false", "sensor"}, 0, []string{"%w حفظ السبب الأصلي.", "الرسالة اتغيرت لكن الهوية فضلت.", "errors.Is بترجّع true أو false."}),
		ch05Walkthrough("c53.cause-reveal", `تابع الـerror الملفوف لحد السبب الأصلي.`, ch05PredictWrapCode, "e53-cause", []string{"errors", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "أنشئ السبب الأصلي."},
			{Line: 3, Caption: "%w يحتفظ بالسبب."},
			{Line: 4, Caption: "السبب موجود جوه الـerror.", Output: "true"},
		}),
		ch05Sandbox("c53.wrap-sandbox", "بدّل ~format~ وشوف النص وهوية السبب.", "ch05-wrap", ch05WrapSandboxConfig()),
	}},
	{ID: "c53.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch05Markdown("c53.watch", `~%v~ بيطبع السبب بس؛ مش بيحافظ على هويته زي ~%w~.

الـvalue اللي بترجّعها ~Diagnose(i)~ ممكن تكون ~nil~؛ افحصها قبل استخراج الكود.

~errors.Join~ بتجمع الأسباب؛ خزّن الـvalue الراجعة عشان السجل يكبر.`)}},
	{ID: "c53.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch05Markdown("c53.plan", `امشِ بخطوات صغيرة: افحص index، وتعامل مع الـerror لو موجود، وسجّل الناتج المناسب.

		خزّن تراكم الأسباب بين اللفات؛ ما تستبدلش اللي اتجمع.

بعد ما الفحص يخلص، قارن ترتيب السجل بالتقارير الفردية.`),
		ch05Question("c53.q-wrap", "أنهي صيغة بتحافظ على سبب error جوه error جديد؟", "fmt.Errorf(\"unit: %w\", err)", []string{"%w", "%v", "%s"}, 0, []string{"%w بيربط السبب الأصلي.", "%v بيعرض النص بس.", "%s للنص، مش للـwrapping."}),
	}},
	{ID: "c53.turn-step", Title: "دورك", Blocks: []LessonBlock{ch05Markdown("c53.turn", `ضيف دورة الفحص ومعالجة الـerrors. راجع كل تقرير بالترتيب بعد التشغيل.`)}},
}

var ch05FaultEnglish = []LessonStep{
	{ID: "c53.mission-step", Title: "Mission", Blocks: []LessonBlock{ch05Markdown("c53.mission", `Inspect the diagnostic bays, report each fault code, then combine the causes.

The starter prints the bay count; continue from that source.

Some bays are healthy, so not every result is an error.`)}},
	{ID: "c53.idea-step", Title: "Idea", Blocks: []LessonBlock{ch05Walkthrough("c53.wrap-trace", `~%w~ adds context while preserving the cause's identity. ~errors.Is~ can still find it after wrapping.

These values are sample alarms, not mission diagnostics.`, ch05WrapCode, "e53-wrap", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create a sample cause."},
		{Line: 3, Caption: "Wrap it with a location."},
		{Line: 4, Caption: "Is still finds the cause.", Variables: []LessonVariable{{Name: "wrapped", Value: "bay 2: valve"}}, Output: "true"},
		{Line: 5, Caption: "The text keeps its context.", Output: "true\nbay 2: valve"},
	})}},
	{ID: "c53.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch05Question("c53.predict-cause", "Predict the cause identity check:", ch05PredictWrapCode, []string{"true", "false", "sensor"}, 0, []string{"%w preserved the original cause.", "The message changed, but identity stayed.", "errors.Is returns true or false."}),
		ch05Walkthrough("c53.cause-reveal", `Follow the wrapped error to its cause.`, ch05PredictWrapCode, "e53-cause", []string{"errors", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Create the original cause."},
			{Line: 3, Caption: "%w preserves the cause."},
			{Line: 4, Caption: "The cause remains inside.", Output: "true"},
		}),
		ch05Sandbox("c53.wrap-sandbox", "Toggle the ~format~; watch text and cause identity.", "ch05-wrap", ch05WrapSandboxConfig()),
	}},
	{ID: "c53.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch05Markdown("c53.watch", `~%v~ prints the cause but does not preserve its identity like ~%w~.

~Diagnose(i)~ may return ~nil~; check before extracting a code.

~errors.Join~ combines causes; keep its returned value to grow the log.`)}},
	{ID: "c53.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch05Markdown("c53.plan", `Follow a small cycle: inspect an index, handle an error if present, and record the appropriate result.

Keep accumulated causes between iterations; do not replace what is already collected.

After inspection, compare the combined log with the individual reports.`),
		ch05Question("c53.q-wrap", "Which format preserves an error cause inside a new error?", "fmt.Errorf(\"unit: %w\", err)", []string{"%w", "%v", "%s"}, 0, []string{"%w wraps the original cause.", "%v displays text only.", "%s formats text; it does not wrap."}),
	}},
	{ID: "c53.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch05Markdown("c53.turn", `Add the inspection cycle and error handling. Check every report in order after running.`)}},
}

var ch05UndockArabic = []LessonStep{
	{ID: "c54.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch05Markdown("c54.mission", `نفّذ محاولات الحفر المطلوبة، وفك ارتباط الدرون بعد كل محاولة حتى لو حصل ~panic~.

الـstarter يعرض عدد المحاولات بس.

بعد ما المحاولات تخلص، فيه حركة وتقرير واحد.`)}},
	{ID: "c54.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch05Walkthrough("c54.defer-trace", `الـdefer بيتنفّذ لما الـfunction تخلص. كده التنضيف بيحصل بعد الشغل العادي.

كل frame يوضح السطر اللي اشتغل وترتيب الخروج.`, ch05DeferCode, "e54-defer", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 5, Caption: "ابدأ استدعاء work."},
		{Line: 2, Caption: "سجّل طباعة الإغلاق."},
		{Line: 3, Caption: "نفّذ الشغل دلوقتي.", Output: "working"},
		{Line: 2, Caption: "defer بيتنفّذ لما الـfunction تخلص.", Output: "working\nclosed"},
	})}},
	{ID: "c54.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch05Question("c54.predict-order", "توقّع النص بعد انتهاء ~work~:", ch05PredictDeferCode, []string{"ABC", "CBA", "CAB"}, 1, []string{"defer مش بيتنفّذ بترتيب تسجيله.", "آخر defer بيتنفّذ الأول بعد C.", "A اتسجلت الأول، فهتتنفّذ في الآخر."}),
		ch05Walkthrough("c54.order-reveal", `راقب التنفيذ الفوري ثم ترتيب النداءات المؤجلة.`, ch05PredictDeferCode, "e54-order", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 7, Caption: "ابدأ استدعاء work."},
			{Line: 4, Caption: "اطبع C فورًا.", Output: "C"},
			{Line: 3, Caption: "آخر defer بيطبع B.", Output: "CB"},
			{Line: 2, Caption: "أول defer بيطبع A.", Output: "CBA"},
			{Line: 8, Caption: "ضيف سطر النهاية.", Output: "CBA"},
		}),
		ch05Sandbox("c54.defer-sandbox", "شغّل الخطوات وشوف ترتيب نداءات ~defer~.", "ch05-defer", ch05DeferSandboxConfig()),
	}},
	{ID: "c54.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch05Markdown("c54.watch", `لازم تنادي ~recover~ مباشرة جوه function معمول لها ~defer~.

بعد التعافي، التنفيذ بيرجع للـcaller؛ مش بيكمل بعد سطر ~panic~.

التنضيف مكانه جوه حدود المحاولة.`)}},
	{ID: "c54.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch05Markdown("c54.plan", `فكّر في كل محاولة كحد مستقل: بداية العمل، خطوة ممكن تفشل، وتنضيف مضمون عند الخروج.

خلّي الـloop في مستوى يقدر يبدأ محاولة جديدة بعد ما الـfunction ترجع.

افصل رسائل التشخيص عن تقرير المهمة.`),
		ch05Question("c54.q-recover", "فين يقدر ~recover~ يمسك panic؟", "defer func() { _ = recover() }()", []string{"داخل function مؤجلة", "في أي سطر من main", "بعد انتهاء البرنامج"}, 0, []string{"recover بيشتغل مباشرة جوه deferred function.", "السطر العادي مش بيلقط panic.", "بعد ما البرنامج يخلص، مفيش تنفيذ."}),
	}},
	{ID: "c54.turn-step", Title: "دورك", Blocks: []LessonBlock{ch05Markdown("c54.turn", `قسّم المحاولة عن الـloop. اختبر التشغيل وشوف إن كل محاولة اتقفلت قبل التقرير.`)}},
}

var ch05UndockEnglish = []LessonStep{
	{ID: "c54.mission-step", Title: "Mission", Blocks: []LessonBlock{ch05Markdown("c54.mission", `Perform the required drill attempts and undock after each one, even if it panics.

The starter only prints the attempt count.

After all attempts, the drone moves and sends one report.`)}},
	{ID: "c54.idea-step", Title: "Idea", Blocks: []LessonBlock{ch05Walkthrough("c54.defer-trace", `A deferred call runs when its function exits. This makes cleanup happen after normal work.

Each frame shows the executed line and exit order.`, ch05DeferCode, "e54-defer", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 5, Caption: "Call work."},
		{Line: 2, Caption: "Register the closing print."},
		{Line: 3, Caption: "Run the work now.", Output: "working"},
		{Line: 2, Caption: "defer runs on exit.", Output: "working\nclosed"},
	})}},
	{ID: "c54.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch05Question("c54.predict-order", "Predict the text after ~work~ returns:", ch05PredictDeferCode, []string{"ABC", "CBA", "CAB"}, 1, []string{"defer does not run in registration order.", "The last defer runs first after C.", "A was registered first, so it runs last."}),
		ch05Walkthrough("c54.order-reveal", `Watch immediate work, then deferred calls.`, ch05PredictDeferCode, "e54-order", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 7, Caption: "Call work."},
			{Line: 4, Caption: "Print C immediately.", Output: "C"},
			{Line: 3, Caption: "The last defer prints B.", Output: "CB"},
			{Line: 2, Caption: "The first defer prints A.", Output: "CBA"},
			{Line: 8, Caption: "Add the final newline.", Output: "CBA"},
		}),
		ch05Sandbox("c54.defer-sandbox", "Step through calls and watch ~defer~ order.", "ch05-defer", ch05DeferSandboxConfig()),
	}},
	{ID: "c54.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch05Markdown("c54.watch", `~recover~ must be called directly inside a deferred function.

After recovery, execution returns to the caller; it does not resume after ~panic~.

Keep cleanup inside the attempt boundary.`)}},
	{ID: "c54.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch05Markdown("c54.plan", `Think of each attempt as its own boundary: start work, reach a step that may fail, and guarantee cleanup on exit.

Keep the loop at a level that can start another attempt after the function returns.

Separate diagnostics from the mission report.`),
		ch05Question("c54.q-recover", "Where can ~recover~ catch a panic?", "defer func() { _ = recover() }()", []string{"Inside a deferred function", "On any line in main", "After the program exits"}, 0, []string{"recover works directly inside a deferred function.", "An ordinary line does not catch a panic.", "Nothing executes after the program exits."}),
	}},
	{ID: "c54.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch05Markdown("c54.turn", `Separate the attempt from the loop. Run the mission and check that every attempt is closed before the report.`)}},
}

var ch05CommandArabic = []LessonStep{
	{ID: "c55.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch05Markdown("c55.mission", `اختار function لكل أمر، نفّذها، وعدّ الأوامر المعروفة بـclosure.

الـstarter بيعرض عدد الأوامر؛ النصوص نفسها تتقري من API.

حافظ على ترتيب النتائج وتعامل مع الاسم المجهول.`)}},
	{ID: "c55.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch05Walkthrough("c55.function-trace", `اسم الـfunction من غير أقواس هو value تقدر تخزّنها. الأقواس هي اللي بتنفّذ النداء.

الـfunction value هنا نوعها ~func()~.`, ch05FunctionValueCode, "e55-function-value", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 1, Caption: "ping بتطبع علامة."},
		{Line: 3, Caption: "خزّن الـfunction نفسها."},
		{Line: 4, Caption: "نفّذ الـfunction اللي خزّنتها."},
		{Line: 5, Caption: "ظهر النداء في الناتج.", Output: "P!"},
	})}},
	{ID: "c55.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch05Question("c55.predict-closure", "كام هتطبع closure بعد النداءين؟", ch05PredictClosureCode, []string{"1", "2", "0"}, 1, []string{"hits اتعرّف برّه closure وبيتشارك.", "كل نداء بيزوّد نفس المتغير.", "hits مش بيرجع لصفر بعد كل نداء."}),
		ch05Walkthrough("c55.closure-reveal", `تابع المتغير اللي closure بتحتفظ بيه.`, ch05PredictClosureCode, "e55-closure", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "العداد بيبدأ بصفر.", Variables: []LessonVariable{{Name: "hits", Value: "0"}}},
			{Line: 3, Caption: "closure تلتقط hits."},
			{Line: 4, Caption: "النداء الأول بيزوّد العداد.", Variables: []LessonVariable{{Name: "hits", Value: "1"}}},
			{Line: 5, Caption: "النداء التاني بيزوّده كمان.", Variables: []LessonVariable{{Name: "hits", Value: "2"}}},
			{Line: 6, Caption: "اطبع الـvalue النهائية.", Variables: []LessonVariable{{Name: "hits", Value: "2"}}, Output: "2"},
		}),
		ch05Sandbox("c55.closure-sandbox", "خزّن function value، وبعدها نادِها وشوف عدّاد الـclosure.", "ch05-closure", ch05ClosureSandboxConfig()),
	}},
	{ID: "c55.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch05Markdown("c55.watch", `~action := ping~ بيخزّن الـfunction؛ ~action()~ بينفّذها.

اتأكد إن نوع الـfunction اللي اختارتها يناسب الـparameter اللي هيستقبلها.

الـclosure لازم تلتقط عدادًا يعيش بين النداءات.`)}},
	{ID: "c55.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch05Markdown("c55.plan", `اقرا القائمة بالترتيب، ونضّف المسافات قبل اختيار الأمر.

		افصل بين اختيار الـfunction وتنفيذها، وبعدها حدّث العداد لو الأمر معروف.

اتأكد إن المجهول ما يوقفش بقية القائمة.`),
		ch05Question("c55.q-value", "أنهي صيغة بتمرر function value من غير ما تشغّلها؟", "run(ping)", []string{"run(ping)", "run(ping())", "run(nil())"}, 0, []string{"ping هنا function value.", "الأقواس بتشغّل ping فورًا.", "nil مش function ينفع تناديها."}),
	}},
	{ID: "c55.turn-step", Title: "دورك", Blocks: []LessonBlock{ch05Markdown("c55.turn", `اربط اختيار الأمر بالـclosure. شغّل المهمة وراجع ترتيب التقارير والعدد النهائي.`)}},
}

var ch05CommandEnglish = []LessonStep{
	{ID: "c55.mission-step", Title: "Mission", Blocks: []LessonBlock{ch05Markdown("c55.mission", `Choose a function for each command, run it, and count known commands with a closure.

The starter prints the order count; read command names from the API.

Keep result order and handle an unknown name.`)}},
	{ID: "c55.idea-step", Title: "Idea", Blocks: []LessonBlock{ch05Walkthrough("c55.function-trace", `A function name without parentheses is a value you can store. Parentheses execute the call.

This function value has type ~func()~.`, ch05FunctionValueCode, "e55-function-value", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "ping prints a marker."},
		{Line: 3, Caption: "Store the function itself."},
		{Line: 4, Caption: "Call the stored value."},
		{Line: 5, Caption: "The call appears in output.", Output: "P!"},
	})}},
	{ID: "c55.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch05Question("c55.predict-closure", "What does the closure print after two calls?", ch05PredictClosureCode, []string{"1", "2", "0"}, 1, []string{"hits is declared outside and shared.", "Each call increments the same variable.", "The variable is not reset at each call."}),
		ch05Walkthrough("c55.closure-reveal", `Follow the variable captured by the closure.`, ch05PredictClosureCode, "e55-closure", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The counter starts at zero.", Variables: []LessonVariable{{Name: "hits", Value: "0"}}},
			{Line: 3, Caption: "The closure captures hits."},
			{Line: 4, Caption: "First call increments it.", Variables: []LessonVariable{{Name: "hits", Value: "1"}}},
			{Line: 5, Caption: "Second call increments it again.", Variables: []LessonVariable{{Name: "hits", Value: "2"}}},
			{Line: 6, Caption: "Print the final value.", Variables: []LessonVariable{{Name: "hits", Value: "2"}}, Output: "2"},
		}),
		ch05Sandbox("c55.closure-sandbox", "Store a ~function~ value, then call it and watch the ~closure~ count.", "ch05-closure", ch05ClosureSandboxConfig()),
	}},
	{ID: "c55.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch05Markdown("c55.watch", `~action := ping~ stores the function; ~action()~ runs it.

Check that the selected function type matches its receiver parameter.

The closure must capture a counter that survives between calls.`)}},
	{ID: "c55.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch05Markdown("c55.plan", `Read the list in order and trim spaces before choosing a command.

Separate function selection from execution, then update the counter for a known command.

Make sure an unknown entry does not stop the rest of the list.`),
		ch05Question("c55.q-value", "Which form passes a function value without running it?", "run(ping)", []string{"run(ping)", "run(ping())", "run(nil())"}, 0, []string{"ping is a function value here.", "Parentheses run ping immediately.", "nil is not a callable function."}),
	}},
	{ID: "c55.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch05Markdown("c55.turn", `Connect command selection to the closure. Run the mission and check report order and the final count.`)}},
}
