package levels

type ch16FrameCopy struct {
	line   int
	ar     string
	en     string
	vars   []LessonVariable
	output string
}

func ch16Copy(ar, en string, english bool) string {
	if english {
		return en
	}
	return ar
}

type ch16SandboxRowData struct {
	label  map[string]string
	values []string
}

func ch16SandboxText(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

func ch16SandboxRow(ar, en string, values ...string) ch16SandboxRowData {
	return ch16SandboxRowData{label: ch16SandboxText(ar, en), values: values}
}

func ch16SandboxFrame(ar, en string, rows ...ch16SandboxRowData) map[string]any {
	encodedRows := make([]map[string]any, len(rows))
	for i, row := range rows {
		encodedRows[i] = map[string]any{"label": row.label, "values": row.values}
	}
	return map[string]any{"caption": ch16SandboxText(ar, en), "rows": encodedRows}
}

func ch16SandboxConfig(frames ...map[string]any) map[string]any {
	return map[string]any{
		"labels": map[string]map[string]string{
			"previous": ch16SandboxText("السابق", "Back"),
			"next":     ch16SandboxText("اللي بعده", "Next"),
			"reset":    ch16SandboxText("من الأول", "Reset"),
		},
		"frames": frames,
	}
}

func ch16TraceBlock(id string, config map[string]any, english bool) LessonBlock {
	return LessonBlock{
		ID:      id,
		Kind:    "sandbox",
		Text:    ch16Copy("تابع الحالات واحدة واحدة.", "Step through each state.", english),
		Sandbox: &LessonSandbox{Kind: "ch16-lab-trace", Config: config},
	}
}

func ch16Markdown(id, ar, en string, english bool) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: ch16Copy(ar, en, english)}
}

func ch16Question(id, ar, en string, code string, optionsAR, optionsEN []string, answer int, feedbackAR, feedbackEN []string, english bool) LessonBlock {
	return LessonBlock{
		ID:   id,
		Kind: "question",
		Text: ch16Copy(ar, en, english) + "\n\n~~~go\n" + code + "\n~~~",
		Question: &LessonQuestion{
			Options:  ch16CopySlice(optionsAR, optionsEN, english),
			Answer:   answer,
			Feedback: ch16CopySlice(feedbackAR, feedbackEN, english),
		},
	}
}

func ch16Choice(id, ar, en string, optionsAR, optionsEN []string, answer int, feedbackAR, feedbackEN []string, english bool) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     ch16Copy(ar, en, english),
		Question: &LessonQuestion{Options: ch16CopySlice(optionsAR, optionsEN, english), Answer: answer, Feedback: ch16CopySlice(feedbackAR, feedbackEN, english)},
	}
}

func ch16CopySlice(ar, en []string, english bool) []string {
	if english {
		return en
	}
	return ar
}

func ch16Frames(copies []ch16FrameCopy, english bool) []LessonFrame {
	frames := make([]LessonFrame, len(copies))
	for i, frame := range copies {
		frames[i] = LessonFrame{Line: frame.line, Caption: ch16Copy(frame.ar, frame.en, english), Variables: frame.vars, Output: frame.output}
	}
	return frames
}

func ch16Walkthrough(id, ar, en, code, exampleID string, imports []string, frames []ch16FrameCopy, english bool) LessonBlock {
	return LessonBlock{
		ID:        id,
		Kind:      "walkthrough",
		Text:      ch16Copy(ar, en, english),
		Code:      code,
		Imports:   imports,
		ExampleID: exampleID,
		Walkthrough: &LessonWalkthrough{
			Frames: ch16Frames(frames, english),
		},
	}
}

func ch16Mission1(english bool) []LessonStep {
	rows := `type Sample struct { name string; value int }
func main() {
	rows := []Sample{{"low", 3}, {"high", 7}}
	for _, row := range rows {
		fmt.Println(row.name, row.value)
	}
}`
	additions := `func main() {
	totals := []int{2, 5}
	sum := 0
	for _, n := range totals {
		sum += n
	}
	fmt.Println(sum)
}`
	return []LessonStep{
		{ID: "c161.mission-step", Title: ch16Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch16Markdown(
			"c161.mission",
			"1. المختبر بيشغّل ~TestFuel~ ضد ست نسخ متبوّظة.\n\n2. وسّع تغطية السلوك عشان كل تغيير غلط يبان.\n\n3. الـstarter فيه حالة واحدة؛ راجع ملف الكود وملف الـtest.",
			"1. The lab runs ~TestFuel~ against six broken versions.\n\n2. Expand behavior coverage so each wrong change shows up.\n\n3. Starter test has one case.", english)}},
		{ID: "c161.idea-step", Title: ch16Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch16Walkthrough(
			"c161.rows-example",
			"الجدول هنا مجرد صفوف بيانات للتوضيح، مش test.\n\nكل صف بيربط اسم بـvalue. مرّ عليهم في loop واحدة عشان المقارنة تبقى منظمة بدل أرقام متناثرة.",
			"This table is illustrative data, not a test.\n\nEach row pairs a name with a value. Looping over rows keeps comparisons organized.",
			rows, "e161-rows", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 1, ar: "النوع يجمع حقول الصف.", en: "The type groups row fields."},
				{line: 3, ar: "صفّين للتوضيح.", en: "Two illustrative rows.", vars: []LessonVariable{{Name: "rows", Value: "[{low 3} {high 7}]"}}},
				{line: 4, ar: "أول loop بتقرا low.", en: "First loop reads low.", vars: []LessonVariable{{Name: "row.name", Value: "low"}, {Name: "row.value", Value: "3"}}},
				{line: 5, ar: "ظهر الصف الأول.", en: "Printed the first row.", output: "low 3"},
				{line: 4, ar: "بعدها بتقرا high.", en: "Next it reads high.", output: "low 3", vars: []LessonVariable{{Name: "row.name", Value: "high"}, {Name: "row.value", Value: "7"}}},
				{line: 5, ar: "الصفين ظهروا بالترتيب.", en: "Both rows printed in order.", output: "low 3\nhigh 7"},
			}, english)}},
		{ID: "c161.try-step", Title: ch16Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch16Question("c161.predict-sum", "توقّع ناتج المثال ده:", "Predict this example's output:", additions, []string{"7", "25", "2"}, []string{"7", "25", "2"}, 0,
				[]string{"الـvalues الاتنين بيتجمعوا في sum.", "ده شكل الـvalues جنب بعض، مش ناتج جمع.", "ده أول عنصر بس، والـloop بتكمل."},
				[]string{"Both values are added to sum.", "That shows the values side by side, not their sum.", "That is only the first value; the loop continues."}, english),
			ch16Walkthrough("c161.sum-reveal", "تابع الـloop لحد ناتجها.", "Follow the loop to its result.", additions, "e161-sum", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "الـslice فيها 2 و5.", en: "The slice holds 2 and 5.", vars: []LessonVariable{{Name: "totals", Value: "[2 5]"}}},
				{line: 3, ar: "المجموع يبدأ بصفر.", en: "The total starts at zero.", vars: []LessonVariable{{Name: "sum", Value: "0"}}},
				{line: 4, ar: "أول value هي 2.", en: "The first value is 2.", vars: []LessonVariable{{Name: "sum", Value: "0"}, {Name: "n", Value: "2"}}},
				{line: 5, ar: "المجموع بقى 2.", en: "The total becomes 2.", vars: []LessonVariable{{Name: "sum", Value: "2"}}},
				{line: 5, ar: "بعدها ضاف 5.", en: "Then it adds 5.", vars: []LessonVariable{{Name: "sum", Value: "7"}}},
				{line: 7, ar: "الناتج هو 7.", en: "The output is 7.", vars: []LessonVariable{{Name: "sum", Value: "7"}}, output: "7"},
			}, english),
			ch16TraceBlock("c161.sandbox", ch16SandboxConfig(
				ch16SandboxFrame("كل وعد محتاج ملاحظة.", "Each promise needs an observation.",
					ch16SandboxRow("الوعد", "Promise", "A"), ch16SandboxRow("الملاحظات", "Observations", "0")),
				ch16SandboxFrame("افحص سلوكًا واحدًا.", "Inspect one behavior at a time.",
					ch16SandboxRow("الوعد", "Promise", "A"), ch16SandboxRow("الملاحظات", "Observations", "1")),
				ch16SandboxFrame("الوعد التالي محتاج فحص.", "The next promise needs a check.",
					ch16SandboxRow("الوعد", "Promise", "B"), ch16SandboxRow("الملاحظات", "Observations", "2")),
			), english),
		}},
		{ID: "c161.watch-step", Title: ch16Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch16Markdown(
			"c161.watch",
			"نجاح النسخة السليمة مش دليل إن test قوي.\n\nmutant بيغيّر وعد واحد؛ حالة ناقصة ممكن تسيبه يعدّي.\n\nخليك مركز على السلوك، مش شكل الكود.",
			"Passing the correct version does not prove a strong test.\n\nA mutant changes one promise; a missing case can let it survive.\n\nFocus on behavior, not code shape.", english)}},
		{ID: "c161.check-step", Title: ch16Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch16Markdown("c161.plan", "اكتب لكل وعد ملاحظة واضحة.\n\nراجع الـtest مع كل mutant، وسجّل أنهي سلوك لسه مش متغطّي.", "Write one clear observation per promise.\n\nReview the test against each mutant and note which behavior is still uncovered.", english),
			ch16Choice("c161.mutant-question", "لو mutant عدى، إيه الاستنتاج الأقرب؟", "If a mutant survives, what is the best conclusion?", []string{"فيه سلوك محتاج فحص أوضح", "الكود الأصلي أكيد غلط", "التست بقى أبطأ"}, []string{"A behavior needs a clearer check", "The original code must be wrong", "The test became slower"}, 0,
				[]string{"مرور النسخة الغلط علامة إن وعدها مش ظاهر في الـtest.", "الـmutant هو اللي اتغير؛ ده مش دليل ضد الأصل.", "السرعة مش بتفسر ليه mutant عدى."},
				[]string{"A wrong version passing means its promise is not visible in the test.", "The mutant changed; this does not disprove the original.", "Speed does not explain why a mutant survived."}, english),
		}},
		{ID: "c161.turn-step", Title: ch16Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch16Markdown(
			"c161.turn",
			"ابدأ من ~TestFuel~ الموجود.\n\nزوّد ملاحظات للسلوك اللي ناقص، وبعدها شغّل الـtests واقرا وصف كل mutant عدى.",
			"Start with the existing ~TestFuel~.\n\nAdd observations for missing behavior, then run the tests and read each surviving mutant's description.", english)}},
	}
}

func ch16Mission2(english bool) []LessonStep {
	views := `func main() {
	readings := []int{8, 3}
	view := readings[:1]
	view[0] = 4
	fmt.Println(readings)
}`
	tail := `func main() {
	data := []int{6, 2}
	tail := data[1:]
	tail[0]++
	fmt.Println(data)
}`
	return []LessonStep{
		{ID: "c162.mission-step", Title: ch16Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch16Markdown("c162.mission",
			"1. ~Median~ فيها سلوك غلط؛ اكتب test يثبته قبل الإصلاح.\n\n2. بعدين خلّي الـtest يعدّي من غير خرق العقد.\n\n3. راجع إن ~main.go~ و~main_test.go~ بيغطّوا النتيجة والـinput.",
			"1. ~Median~ has wrong behavior; write a failing test before repairing it.\n\n2. Then make it pass without breaking the contract.\n\n3. Review whether ~main.go~ and ~main_test.go~ cover the result and input.", english)}},
		{ID: "c162.idea-step", Title: ch16Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch16Walkthrough("c162.slice-view",
			"الـslice زي نافذتين على نفس البيانات.\n\nتغيير عنصر من view ممكن يبان في readings؛ الـslice مش بتتنقل لنسخة تانية لمجرد إنك عملت slice جديدة.",
			"A slice is like a window onto shared data.\n\nChanging an element through view can appear in readings; a new slice value does not automatically copy the elements.",
			views, "e162-slice-view", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "القراءات تبدأ 8 و3.", en: "Readings start as 8 and 3.", vars: []LessonVariable{{Name: "readings", Value: "[8 3]"}}},
				{line: 3, ar: "view شايفة أول عنصر.", en: "view covers the first element.", vars: []LessonVariable{{Name: "view", Value: "[8]"}}},
				{line: 4, ar: "الـvalue اللي بينهم بقت 4.", en: "The shared value becomes 4.", vars: []LessonVariable{{Name: "view[0]", Value: "4"}, {Name: "readings", Value: "[4 3]"}}},
				{line: 5, ar: "الأصل اتأثر بالتغيير.", en: "The original reflects the change.", vars: []LessonVariable{{Name: "readings", Value: "[4 3]"}}, output: "[4 3]"},
			}, english)}},
		{ID: "c162.try-step", Title: ch16Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch16Question("c162.predict-view", "إيه اللي هيظهر بعد تغيير tail؟", "What prints after tail changes?", tail, []string{"[6 3]", "[6 2]", "[3 6]"}, []string{"[6 3]", "[6 2]", "[3 6]"}, 0,
				[]string{"tail بتشارك نفس الـbacking array مع data.", "ده شكل data قبل التعديل.", "التغيير ما بيبدلش ترتيب العناصر."},
				[]string{"tail shares data with data.", "That is data before the change.", "The update does not reorder elements."}, english),
			ch16Walkthrough("c162.tail-reveal", "شوف التغيير في الـslice الأصلية.", "Watch the change in the original slice.", tail, "e162-tail", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "data فيها 6 و2.", en: "data holds 6 and 2.", vars: []LessonVariable{{Name: "data", Value: "[6 2]"}}},
				{line: 3, ar: "tail بتبدأ عند 2.", en: "tail starts at 2.", vars: []LessonVariable{{Name: "tail", Value: "[2]"}}},
				{line: 4, ar: "العنصر بقى 3.", en: "The element becomes 3.", vars: []LessonVariable{{Name: "tail[0]", Value: "3"}, {Name: "data", Value: "[6 3]"}}},
				{line: 5, ar: "data بتطبع بعد التعديل.", en: "data prints after the change.", vars: []LessonVariable{{Name: "data", Value: "[6 3]"}}, output: "[6 3]"},
			}, english),
			ch16TraceBlock("c162.sandbox", ch16SandboxConfig(
				ch16SandboxFrame("الـview بتشارك أول عنصر.", "The view shares its first element.",
					ch16SandboxRow("الأصل", "Original", "8", "3"), ch16SandboxRow("الـview", "View", "8")),
				ch16SandboxFrame("تغيير الـview يغيّر الأصل.", "Changing the view changes its source.",
					ch16SandboxRow("الأصل", "Original", "4", "3"), ch16SandboxRow("الـview", "View", "4")),
			), english),
		}},
		{ID: "c162.watch-step", Title: ch16Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch16Markdown("c162.watch",
			"الـtest على الـvalue لوحدها ممكن يفوّت side effect.\n\nفصل النتيجة عن سلامة input بيوضح الوعد؛ الحالة العادية لوحدها مش كفاية.",
			"Checking the returned value alone can miss a side effect.\n\nSeparate the result from input preservation; an ordinary case alone is not enough.", english)}},
		{ID: "c162.check-step", Title: ch16Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch16Markdown("c162.plan", "امشِ بـ Red ثم Green: خلّي العيب ظاهر، أصلح السلوك، وارجع راجع كل وعد مستقل.", "Work Red then Green: expose the defect, repair behavior, and revisit each independent promise.", english),
			ch16Choice("c162.alias-question", "لو view اتعدّلت، هل readings ممكن تتأثر؟", "Can readings change when view is modified?", []string{"أيوه، الاتنين بيشاركوا العنصر", "لأ، كل slice نسخة تلقائية", "أيوه، بس لو أطوالهم متساوية"}, []string{"Yes, both share the element", "No, every slice is an automatic copy", "Yes, but only when lengths match"}, 0,
				[]string{"الـslice الجديدة ممكن تفضل على نفس الـbacking array.", "تغيير عنوان slice مش معناه نسخ عناصرها.", "الطول مش شرط لمشاركة العناصر."},
				[]string{"The new slice can still use the same backing array.", "A new slice value does not copy its elements.", "Equal lengths are not required for sharing."}, english),
		}},
		{ID: "c162.turn-step", Title: ch16Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch16Markdown("c162.turn",
			"خلّي أول فشل يثبت العيب.\n\nبعد الإصلاح، شغّل الـtests المخفية وبص على إن المدخل فضل زي ما العقد بيقول.",
			"Make the first failure prove the defect.\n\nAfter repair, run the hidden tests and check that the input still matches its contract.", english)}},
	}
}

func ch16Mission3(english bool) []LessonStep {
	work := `func main() {
	sizes := []int{2, 4}
	for _, n := range sizes {
		fmt.Println(n, n*n)
	}
}`
	square := `func main() {
	n := 3
	fmt.Println(n * n)
}`
	return []LessonStep{
		{ID: "c163.mission-step", Title: ch16Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch16Markdown("c163.mission",
			"1. ~CountPairs~ شغالة، بس النسخة الحالية بطيئة مع مدخلات كبيرة.\n\n2. حسّنها من غير ما تغيّر معناها.\n\n3. ضيف benchmark عشان تقيس بدل التخمين.",
			"1. ~CountPairs~ works, but the current version slows down on large input.\n\n2. Improve it without changing its meaning.\n\n3. Add a benchmark so you measure instead of guessing.", english)}},
		{ID: "c163.idea-step", Title: ch16Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch16Walkthrough("c163.growth-example",
			"المثال بيحسب n² كعدد افتراضي للشغل، مش وقت.\n\nلما n يتضاعف، المربّع يكبر أسرع؛ القياس الحقيقي يفضل محتاج benchmark.",
			"This example calculates n² as an abstract work count, not time.\n\nDoubling n grows the square faster; actual speed still needs a benchmark.",
			work, "e163-growth", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "قيمتان للتوضيح.", en: "Two illustrative sizes.", vars: []LessonVariable{{Name: "sizes", Value: "[2 4]"}}},
				{line: 3, ar: "بدأنا بـ n يساوي 2.", en: "Start with n equal to 2.", vars: []LessonVariable{{Name: "n", Value: "2"}}},
				{line: 4, ar: "الشغل الافتراضي 4.", en: "Abstract work count is 4.", vars: []LessonVariable{{Name: "n*n", Value: "4"}}, output: "2 4"},
				{line: 3, ar: "بعدها n يساوي 4.", en: "Next n equals 4.", output: "2 4", vars: []LessonVariable{{Name: "n", Value: "4"}}},
				{line: 4, ar: "الشغل الافتراضي 16.", en: "Abstract work count is 16.", vars: []LessonVariable{{Name: "n*n", Value: "16"}}, output: "2 4\n4 16"},
			}, english)}},
		{ID: "c163.try-step", Title: ch16Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch16Question("c163.predict-square", "لو n = 3، إيه الناتج؟", "If n = 3, what is printed?", square, []string{"6", "9", "3"}, []string{"6", "9", "3"}, 1,
				[]string{"ده الجمع؛ الكود بيضرب n في نفسها.", "3 مضروبة في 3 تساوي 9.", "3 هي value n قبل الضرب."},
				[]string{"That is addition; the code multiplies n by itself.", "Three times three is nine.", "Three is n before multiplication."}, english),
			ch16Walkthrough("c163.square-reveal", "تابع الضرب لحد الطباعة.", "Follow the multiplication to the print.", square, "e163-square", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "n قيمتها 3.", en: "n has value 3.", vars: []LessonVariable{{Name: "n", Value: "3"}}},
				{line: 3, ar: "3 مضروبة في نفسها.", en: "Three multiplied by itself.", vars: []LessonVariable{{Name: "n*n", Value: "9"}}},
				{line: 3, ar: "يطبع ناتج الضرب.", en: "Prints the product.", vars: []LessonVariable{{Name: "n*n", Value: "9"}}, output: "9"},
			}, english),
			ch16TraceBlock("c163.sandbox", ch16SandboxConfig(
				ch16SandboxFrame("الشغل بيتناسب مع n².", "Work grows with n squared.",
					ch16SandboxRow("حجم الإدخال n", "Input size n", "2"),
					ch16SandboxRow("شغل n²", "n-squared work", "4"),
					ch16SandboxRow("loop واحدة", "Single-pass work", "2")),
				ch16SandboxFrame("المربّع كبر أسرع.", "The square grew faster.",
					ch16SandboxRow("حجم الإدخال n", "Input size n", "4"),
					ch16SandboxRow("شغل n²", "n-squared work", "16"),
					ch16SandboxRow("loop واحدة", "Single-pass work", "4")),
				ch16SandboxFrame("زود الحجم وقارن تاني.", "Grow the input and compare again.",
					ch16SandboxRow("حجم الإدخال n", "Input size n", "6"),
					ch16SandboxRow("شغل n²", "n-squared work", "36"),
					ch16SandboxRow("loop واحدة", "Single-pass work", "6")),
			), english),
		}},
		{ID: "c163.watch-step", Title: ch16Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch16Markdown("c163.watch",
			"Big-O بيشرح نمو الشغل، مش زمن جهازك.\n\nbenchmark ممكن يختلف حسب المدخل والبيئة؛ خليك على نفس البيانات وقت المقارنة.",
			"Big-O describes work growth, not your machine's runtime.\n\nA benchmark varies with input and environment; keep the data consistent when comparing.", english)}},
		{ID: "c163.check-step", Title: ch16Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch16Markdown("c163.plan", "اتأكد من صحة النتيجة الأول.\n\nبعدها شغّل benchmark على نفس الحمل، واقرا الوقت لكل عملية بدل ما تعتمد على إحساسك.", "Check result correctness first.\n\nThen benchmark the same workload and compare time per operation instead of trusting intuition.", english),
			ch16Choice("c163.benchmark-question", "الـbenchmark الأسرع يثبت إيه؟", "What does a faster benchmark establish?", []string{"سرعة القياس في الإعداد ده", "صحة كل النتائج", "التعامل مع كل البيانات"}, []string{"Measured speed in this setup", "Correctness of every result", "Handling every input"}, 0,
				[]string{"القياس بيقارن سرعة التنفيذ بس.", "الصحة محتاجة checks منفصلة.", "حمل واحد مش بيمثل كل inputs."},
				[]string{"A benchmark compares execution speed only.", "Correctness needs separate checks.", "One workload does not represent every input."}, english),
		}},
		{ID: "c163.turn-step", Title: ch16Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch16Markdown("c163.turn",
			"بص على الـbenchmark stub في الـtest، وشغّل التحقق قبل القياس.\n\nقارن نسختك على نفس البيانات واقرا الرقم بهدوء.",
			"Review the benchmark stub in the test file and check correctness before measuring.\n\nCompare your version on the same data and read the result carefully.", english)}},
	}
}

func ch16Mission4(english bool) []LessonStep {
	guard := `func main() {
	packet := []byte{11, 22}
	if len(packet) > 0 {
		fmt.Println(packet[0])
	}
}`
	read := `func main() {
	packet := []byte{9, 4}
	fmt.Println(len(packet), packet[1])
}`
	return []LessonStep{
		{ID: "c164.mission-step", Title: ch16Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch16Markdown("c164.mission",
			"1. مدخلات كبسولة التحكم ممكن تبقى ناقصة أو غريبة.\n\n2. اكتب fuzz target يكشف panic وخرق العقد، وبعدها أصلح parser.\n\n3. الملفين قابلين للتعديل؛ ابدأ بعينات صغيرة مفهومة.",
			"1. Control capsule input may be truncated or unusual.\n\n2. Write a fuzz target that reveals panics and contract violations, then repair the parser.\n\n3. Both files are editable; start with understandable seeds.", english)}},
		{ID: "c164.idea-step", Title: ch16Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch16Walkthrough("c164.length-guard",
			"قبل ما تقرا index، اتأكد إن العنصر موجود.\n\nالطول الفعلي هو اللي يحمي القراءة؛ معلومة جاية من input مش ضمان لوحدها.",
			"Before reading an index, check that the element exists.\n\nActual length protects the read; information supplied by input is not a guarantee by itself.",
			guard, "e164-length-guard", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "وصل 2 bytes.", en: "Two bytes arrived.", vars: []LessonVariable{{Name: "packet", Value: "[11 22]"}}},
				{line: 3, ar: "الطول أكبر من صفر.", en: "Length is above zero.", vars: []LessonVariable{{Name: "len(packet)", Value: "2"}}},
				{line: 4, ar: "أول byte قيمته 11.", en: "The first byte is 11.", vars: []LessonVariable{{Name: "packet[0]", Value: "11"}}, output: "11"},
			}, english)}},
		{ID: "c164.try-step", Title: ch16Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch16Question("c164.predict-byte", "إيه اللي هيتطبع؟", "What will this print?", read, []string{"2 4", "9 4", "panic"}, []string{"2 4", "9 4", "panic"}, 0,
				[]string{"في عنصرين، والعنصر عند index 1 هو 4.", "9 موجودة عند index 0.", "الـindex 1 جوه الـslice."},
				[]string{"There are two elements; index 1 holds 4.", "Nine is at index 0.", "Index 1 is inside the slice."}, english),
			ch16Walkthrough("c164.byte-reveal", "راجع الطول والـindex سوا.", "Check the length and index together.", read, "e164-byte", []string{"fmt"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "الـslice فيها 9 و4.", en: "The slice holds 9 and 4.", vars: []LessonVariable{{Name: "packet", Value: "[9 4]"}}},
				{line: 3, ar: "length=2 وindex 1 قيمته 4.", en: "Length is 2; index 1 is 4.", vars: []LessonVariable{{Name: "len(packet)", Value: "2"}, {Name: "packet[1]", Value: "4"}}},
				{line: 3, ar: "يطبع الاتنين بالترتيب.", en: "Prints both values in order.", output: "2 4"},
			}, english),
			ch16TraceBlock("c164.sandbox", ch16SandboxConfig(
				ch16SandboxFrame("وصل 2 bytes.", "Two bytes arrive.",
					ch16SandboxRow("bytes", "Bytes", "9", "4"),
					ch16SandboxRow("الطول", "Length", "2"),
					ch16SandboxRow("أول byte", "First byte", "9")),
				ch16SandboxFrame("وصل byte واحد.", "One byte arrives.",
					ch16SandboxRow("bytes", "Bytes", "11"),
					ch16SandboxRow("الطول", "Length", "1"),
					ch16SandboxRow("أول byte", "First byte", "11")),
				ch16SandboxFrame("مفيش bytes وصلت.", "No bytes arrive.",
					ch16SandboxRow("bytes", "Bytes", "—"),
					ch16SandboxRow("الطول", "Length", "0"),
					ch16SandboxRow("أول byte", "First byte", "—")),
			), english),
		}},
		{ID: "c164.watch-step", Title: ch16Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch16Markdown("c164.watch",
			"عدم الـpanic مش معناه إن كل حاجة صح؛ parser ممكن يرفض كل input.\n\nكمان بلاش recover يخفي الـpanic اللي fuzz المفروض يكشفه.",
			"Avoiding a panic does not mean the behavior is correct; a parser could reject every input.\n\nDo not use recover to hide the panic fuzzing should reveal.", english)}},
		{ID: "c164.check-step", Title: ch16Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch16Markdown("c164.plan", "خلّي الـtest يفرّق بين قبول صحيح ورفض آمن.\n\nseeds بتدي fuzz بداية ثابتة، والـgenerated inputs توسّع البحث.", "Make the test distinguish valid acceptance from safe rejection.\n\nSeeds give fuzzing a stable start; generated inputs widen the search.", english),
			ch16Choice("c164.panic-question", "لو fuzz input سبّب panic، الـtest يعمل إيه؟", "If a fuzz input panics, what should the target do?", []string{"يسيب الفشل يبان", "يعمل recover ويتجاهله", "يمسح الـinput"}, []string{"Let the failure surface", "Recover and ignore it", "Delete the input"}, 0,
				[]string{"الـpanic bug ظاهر للـtest؛ ما تخبيهوش.", "recover هنا بيشيل الدليل.", "الـinput هو اللي كشف المشكلة."},
				[]string{"The panic is a bug for the test to expose.", "recover would hide the evidence.", "The input revealed the problem."}, english),
		}},
		{ID: "c164.turn-step", Title: ch16Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch16Markdown("c164.turn",
			"ابدأ بالـseeds اللي في العقد، وشغّلها قبل fuzz طويل.\n\nكل crash أو نتيجة غريبة حوّلها لملاحظة ثابتة.",
			"Start with the contract's seeds and run them before a long fuzz session.\n\nTurn every crash or strange result into a stable observation.", english)}},
	}
}

func ch16Mission5(english bool) []LessonStep {
	clock := `func main() {
	start := time.Unix(0, 0)
	later := start.Add(2 * time.Second)
	fmt.Println(later.Sub(start))
}`
	compare := `func main() {
	start := time.Unix(0, 0)
	later := start.Add(500 * time.Millisecond)
	fmt.Println(later.After(start))
}`
	return []LessonStep{
		{ID: "c165.mission-step", Title: ch16Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch16Markdown("c165.mission",
			"1. اختبر بوابة الإطلاق باستخدام وقت قابل للتحريك.\n\n2. عدّل ~main_test.go~ بس؛ التنفيذ في ~main.go~ للقراءة.\n\n3. خلّي السيناريوهات واضحة ومن غير انتظار حقيقي.",
			"1. Test the launch gate using controllable time.\n\n2. Edit only ~main_test.go~; read the implementation in ~main.go~.\n\n3. Keep scenarios clear and avoid real waiting.", english)}},
		{ID: "c165.idea-step", Title: ch16Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch16Walkthrough("c165.time-move",
			"الساعة المزيفة بتخلي الوقت input للـtest.\n\nحرّكها لقدّام واقرا الفرق فورًا؛ مفيش sleep ولا اعتماد على سرعة الجهاز.",
			"A fake clock makes time an input to the test.\n\nAdvance it and read the difference immediately; no sleep or machine-speed dependency is needed.",
			clock, "e165-time-move", []string{"fmt", "time"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "بدأنا من وقت ثابت.", en: "Start from a fixed instant.", vars: []LessonVariable{{Name: "start", Value: "Unix epoch"}}},
				{line: 3, ar: "اتحركنا ثانيتين.", en: "Advance by two seconds.", vars: []LessonVariable{{Name: "later", Value: "start + 2s"}}},
				{line: 4, ar: "الفرق يساوي 2s.", en: "The elapsed time is 2s.", output: "2s"},
			}, english)}},
		{ID: "c165.try-step", Title: ch16Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch16Question("c165.predict-after", "هل later بعد start؟", "Is later after start?", compare, []string{"true", "false", "مفيش output"}, []string{"true", "false", "No output"}, 0,
				[]string{"~Add~ بتحرك الوقت لقدّام.", "~later~ اتبنت بإضافة مدة موجبة.", "~fmt.Println~ بيطبع نتيجة المقارنة."},
				[]string{"~Add~ moves the time forward.", "~later~ adds a positive duration.", "~fmt.Println~ prints the comparison."}, english),
			ch16Walkthrough("c165.after-reveal", "تابع مقارنة اللحظتين.", "Follow the comparison between instants.", compare, "e165-after", []string{"fmt", "time"}, []ch16FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "اللحظة الأولى ثابتة.", en: "The first instant is fixed.", vars: []LessonVariable{{Name: "start", Value: "Unix epoch"}}},
				{line: 3, ar: "later بعدها بنص ثانية.", en: "later is half a second later.", vars: []LessonVariable{{Name: "later", Value: "start + 500ms"}}},
				{line: 4, ar: "المقارنة نتيجتها true.", en: "The comparison is true.", output: "true"},
			}, english),
			ch16TraceBlock("c165.sandbox", ch16SandboxConfig(
				ch16SandboxFrame("ابدأ من لحظة ثابتة.", "Start from a fixed instant.",
					ch16SandboxRow("وقت البداية", "Start time", "0 ms"),
					ch16SandboxRow("المدة", "Elapsed", "0 ms")),
				ch16SandboxFrame("حرّك الساعة من غير انتظار.", "Advance the clock without waiting.",
					ch16SandboxRow("وقت البداية", "Start time", "0 ms"),
					ch16SandboxRow("المدة", "Elapsed", "250 ms")),
				ch16SandboxFrame("الوقت يتقدم كـinput.", "Time advances as an input.",
					ch16SandboxRow("وقت البداية", "Start time", "0 ms"),
					ch16SandboxRow("المدة", "Elapsed", "500 ms")),
			), english),
		}},
		{ID: "c165.watch-step", Title: ch16Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch16Markdown("c165.watch",
			"ساعة جديدة لكل subtest تمنع الوقت والرصيد من التسريب.\n\nبص على الـbool والـerror؛ واحدة منهم لوحدها ممكن تسيب سلوك ناقص.",
			"A fresh clock per subtest keeps time and credit isolated.\n\nCheck both bool and error; either one alone can miss behavior.", english)}},
		{ID: "c165.check-step", Title: ch16Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch16Markdown("c165.plan", "قسّم العقد لسيناريوهات مسماة.\n\nفي كل واحد: حرّك fake clock، اطلب من البوابة، وبص على النتيجة والـerror مع بعض.", "Split the contract into named scenarios.\n\nIn each one, advance the fake clock, ask the gate, and check result and error together.", english),
			ch16Choice("c165.clock-question", "ليه كل subtest محتاج ساعة جديدة؟", "Why should each subtest get a fresh clock?", []string{"تغيير A ما يغيّرش B", "عشان نعمل sleep", "عشان benchmark يبقى أسرع"}, []string{"Changing A does not change B", "So we can sleep", "To make a benchmark faster"}, 0,
				[]string{"قيم time مستقلة؛ كل سيناريو يبدأ من حالته.", "الـtest بيتجنب الانتظار الحقيقي.", "السؤال عن العزل، مش السرعة."},
				[]string{"Time values are independent; each scenario starts from its own state.", "The test avoids real waiting.", "This is about isolation, not speed."}, english),
		}},
		{ID: "c165.turn-step", Title: ch16Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch16Markdown("c165.turn",
			"اقرا signatures الأول، وبعدها اكتب subtests بعناوين واضحة.\n\nشغّل race detector مع الـtests من غير أي sleep.",
			"Read the signatures first, then write clearly named subtests.\n\nRun the race detector with no sleeps in the tests.", english)}},
	}
}
