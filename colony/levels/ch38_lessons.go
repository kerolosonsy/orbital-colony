package levels

func ch38Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

type ch38QuestionDetails struct {
	Options  []string
	Answer   int
	Feedback []string
}

func ch38Question(id, text string, details ch38QuestionDetails) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: details.Options, Answer: details.Answer, Feedback: details.Feedback}}
}

type ch38WalkthroughDetails struct {
	ExampleID string
	Imports   []string
	Frames    []LessonFrame
}

func ch38Walkthrough(id, text, code string, details ch38WalkthroughDetails) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: code, Imports: details.Imports, ExampleID: details.ExampleID, Walkthrough: &LessonWalkthrough{Frames: details.Frames}}
}

var ch38IteratorsArabic = []LessonStep{
	{ID: "c381.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c381.mission", "هتكتب iterators بتخلي المستهلك يحدد يكمل ولا يوقف.\n\nحافظ على ترتيب القراءات، وخلي التوقف يحصل فور طلبه.\n\nفرّق بين إنشاء الـSeq واستهلاكها.")}},
	{ID: "c381.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c381.range-values", "الـiterator ممكن يبقى function بتسلّم قيمة كل مرة.\n\n~slices.Values~ بتحوّل slice لـSeq؛ و~range~ هو اللي يبدأ المرور عليها.\n\nهنا المثال بيقرأ قيمتين بس، وبعدها الحلقة تخلص.", `func main() {
	values := []int{2, 5}
	for value := range slices.Values(values) {
		fmt.Println(value)
	}
}`, ch38WalkthroughDetails{ExampleID: "e381-values", Imports: []string{"fmt", "slices"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "القيم موجودة بالترتيب.", Variables: []LessonVariable{{Name: "values", Value: "[2 5]"}}}, {Line: 3, Caption: "range يبدأ الاستهلاك."}, {Line: 4, Caption: "أول قيمة تظهر.", Variables: []LessonVariable{{Name: "value", Value: "2"}}, Output: "2"}, {Line: 4, Caption: "القيمة التالية تظهر.", Variables: []LessonVariable{{Name: "value", Value: "5"}}, Output: "2\n5"}}})}},
	{ID: "c381.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c381.predict-break", "توقّع الناتج بعد ~break~:\n\n~~~go\nfunc main() {\n\tvalues := []int{3, 6, 9}\n\tfor n := range slices.Values(values) {\n\t\tfmt.Println(n)\n\t\tif n == 6 { break }\n\t}\n}\n~~~", ch38QuestionDetails{Options: []string{"3 و6", "3 و6 و9", "9 بس"}, Answer: 0, Feedback: []string{"الـbreak يوقف الاستهلاك بعد القيمة دي.", "القيمة 9 مش هتوصل للحلقة.", "الحلقة تبدأ بأول عنصر في الـslice."}}),
		ch38Walkthrough("c381.break-reveal", "تابع كل قيمة لحد شرط التوقف.", `func main() {
	values := []int{3, 6, 9}
	for n := range slices.Values(values) {
		fmt.Println(n)
		if n == 6 { break }
	}
}`, ch38WalkthroughDetails{ExampleID: "e381-break", Imports: []string{"fmt", "slices"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "القيم بالترتيب."}, {Line: 2, Caption: "range يطلب أول قيمة."}, {Line: 4, Caption: "اطبع 3.", Variables: []LessonVariable{{Name: "n", Value: "3"}}, Output: "3"}, {Line: 4, Caption: "اطبع 6.", Variables: []LessonVariable{{Name: "n", Value: "6"}}, Output: "3\n6"}, {Line: 5, Caption: "break يمنع القيمة التالية.", Variables: []LessonVariable{{Name: "n", Value: "6"}}, Output: "3\n6"}}}),
	}},
	{ID: "c381.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c381.watch", "- إنشاء الـSeq مش معناه إن جسمها اشتغل؛ الاستهلاك بيبدأ مع ~range~.\n- لو ~yield~ رجّعت false، أي yield بعدها تكسر العقد.")}},
	{ID: "c381.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c381.plan", "قسّم التفكير: حدّد إمتى المنتج يسلّم قيمة، وإمتى المستهلك يوقف.\n\nبعدها راجع إن حالة المرور بتبدأ من جديد مع استهلاك جديد.\n\nفي ~First~ فكّر في عدد القيم المطلوبة قبل ما تبدأ."),
		ch38Question("c381.stop-contract", "لما المستهلك يعمل ~break~، المنتج يعمل إيه؟", ch38QuestionDetails{Options: []string{"يوقف إرسال القيم", "يكمل كل القيم", "يعيد أول قيمة"}, Answer: 0, Feedback: []string{"التوقف جزء من عقد الـiterator.", "كده المنتج بيتجاهل قرار المستهلك.", "إعادة المرور تحصل مع استهلاك جديد."}}),
	}},
	{ID: "c381.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c381.turn", "ابدأ بعقد ~Readings~، وبعدها جرّب مستهلك بيقف بدري. شغّل التحقق مع كل جزء صغير.")}},
}

var ch38IteratorsEnglish = []LessonStep{
	{ID: "c381.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c381.mission", "Write iterators that let the consumer decide whether to continue.\n\nKeep readings in order, and stop as soon as requested.\n\nSeparate creating a Seq from consuming it.")}},
	{ID: "c381.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c381.range-values", "An iterator can be a function that hands over one value at a time.\n\n~slices.Values~ turns a slice into a Seq; ~range~ starts consuming it.\n\nThis example reads two values, then the loop ends.", `func main() {
	values := []int{2, 5}
	for value := range slices.Values(values) {
		fmt.Println(value)
	}
}`, ch38WalkthroughDetails{ExampleID: "e381-values", Imports: []string{"fmt", "slices"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Values keep their order.", Variables: []LessonVariable{{Name: "values", Value: "[2 5]"}}}, {Line: 3, Caption: "range starts consumption."}, {Line: 4, Caption: "The first value appears.", Variables: []LessonVariable{{Name: "value", Value: "2"}}, Output: "2"}, {Line: 4, Caption: "The next value appears.", Variables: []LessonVariable{{Name: "value", Value: "5"}}, Output: "2\n5"}}})}},
	{ID: "c381.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c381.predict-break", "Predict the output after ~break~:\n\n~~~go\nfunc main() {\n\tvalues := []int{3, 6, 9}\n\tfor n := range slices.Values(values) {\n\t\tfmt.Println(n)\n\t\tif n == 6 { break }\n\t}\n}\n~~~", ch38QuestionDetails{Options: []string{"3 and 6", "3, 6, and 9", "Only 9"}, Answer: 0, Feedback: []string{"break stops consumption after this value.", "Value 9 never reaches the loop body.", "The loop starts at the first slice element."}}),
		ch38Walkthrough("c381.break-reveal", "Follow each value until the stop condition.", `func main() {
	values := []int{3, 6, 9}
	for n := range slices.Values(values) {
		fmt.Println(n)
		if n == 6 { break }
	}
}`, ch38WalkthroughDetails{ExampleID: "e381-break", Imports: []string{"fmt", "slices"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Values keep their order."}, {Line: 2, Caption: "range asks for the first value."}, {Line: 4, Caption: "Print 3.", Variables: []LessonVariable{{Name: "n", Value: "3"}}, Output: "3"}, {Line: 4, Caption: "Print 6.", Variables: []LessonVariable{{Name: "n", Value: "6"}}, Output: "3\n6"}, {Line: 5, Caption: "break prevents the next value.", Variables: []LessonVariable{{Name: "n", Value: "6"}}, Output: "3\n6"}}}),
	}},
	{ID: "c381.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c381.watch", "- Creating a Seq does not run its body; consumption begins with ~range~.\n- After ~yield~ returns false, another yield breaks the iterator contract.")}},
	{ID: "c381.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c381.plan", "Split the reasoning: decide when the producer hands over a value and when the consumer stops.\n\nThen check that traversal state starts fresh on another consumption.\n\nFor ~First~, consider the requested count before iteration."),
		ch38Question("c381.stop-contract", "When the consumer uses ~break~, what should the producer do?", ch38QuestionDetails{Options: []string{"Stop yielding values", "Continue through all values", "Repeat the first value"}, Answer: 0, Feedback: []string{"Stopping is part of the iterator contract.", "That ignores the consumer's decision.", "A new consumption starts a new traversal."}}),
	}},
	{ID: "c381.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c381.turn", "Start with the ~Readings~ contract, then try a consumer that stops early. Run the checks after each small piece.")}},
}

var ch38SetsArabic = []LessonStep{
	{ID: "c382.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c382.mission", "اختار قيم من صفوف، مع منع المحظور والتكرار.\n\nترتيب الصفوف والمرشحين مهم للنتيجة.\n\nهتحتاج تفرّق بين تخطي عنصر وتخطي الصف كله.")}},
	{ID: "c382.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c382.set-membership", "الـmap ممكن تمثل set لما المفتاح هو العضوية.\n\nالقيمة ~struct{}{}~ فاضية؛ اللي يهم هو وجود المفتاح ونتيجة comma-ok.\n\nالمثال هنا بيفحص أسماء توضيحية.", `func main() {
	labels := map[string]struct{}{}
	labels["north"] = struct{}{}
	_, exists := labels["north"]
	fmt.Println(exists, len(labels))
}`, ch38WalkthroughDetails{ExampleID: "e382-set", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "أنشئ set فاضية.", Variables: []LessonVariable{{Name: "labels", Value: "map[]"}}}, {Line: 3, Caption: "أضف north مرة واحدة.", Variables: []LessonVariable{{Name: "labels", Value: "map[north:{}]"}}}, {Line: 4, Caption: "افحص العضوية."}, {Line: 5, Caption: "المفتاح موجود.", Variables: []LessonVariable{{Name: "exists", Value: "true"}}, Output: "true 1"}}})}},
	{ID: "c382.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c382.predict-label", "توقّع الناتج:\n\n~~~go\nfunc main() {\nDecks:\n\tfor _, deck := range [][]int{{1, 0}, {2, 3}} {\n\t\tfor _, slot := range deck {\n\t\t\tif slot == 0 { continue Decks }\n\t\t\tfmt.Println(slot)\n\t\t}\n\t}\n}\n~~~", ch38QuestionDetails{Options: []string{"1 ثم 2 و3", "1 و0 و2 و3", "2 و3 بس"}, Answer: 0, Feedback: []string{"الصف يبدأ بطباعة 1؛ الصفر يتخطّى باقي صفه.", "الصفر نفسه مش بيتطبع قبل continue.", "القيمة 1 اتطبعت قبل ما الصف يتخطى."}}),
		ch38Walkthrough("c382.label-reveal", "تابع تأثير الـlabel على الحلقة الكبيرة.", `func main() {
Decks:
	for _, deck := range [][]int{{1, 0}, {2, 3}} {
		for _, slot := range deck {
			if slot == 0 { continue Decks }
			fmt.Println(slot)
		}
	}
}`, ch38WalkthroughDetails{ExampleID: "e382-label", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 3, Caption: "ابدأ الصف الأول."}, {Line: 6, Caption: "اطبع أول خانة.", Variables: []LessonVariable{{Name: "slot", Value: "1"}}, Output: "1"}, {Line: 5, Caption: "الـlabel يتخطّى الصف.", Output: "1"}, {Line: 3, Caption: "ابدأ الصف التالي.", Output: "1"}, {Line: 6, Caption: "اطبع خاناته.", Variables: []LessonVariable{{Name: "slot", Value: "2"}}, Output: "1\n2"}, {Line: 6, Caption: "ثم اطبع الخانة التالية.", Variables: []LessonVariable{{Name: "slot", Value: "3"}}, Output: "1\n2\n3"}}}),
	}},
	{ID: "c382.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c382.watch", "- الـmap مش مرتّبة؛ خُد ترتيب النتيجة من المرور على الصفوف.\n- اسم label بيحدد أنهي loop تتخطى؛ اتأكد إنها حلقة محيطة بيك.")}},
	{ID: "c382.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c382.plan", "ارسم الحل كحلقتين: واحدة للصفوف وواحدة للعناصر.\n\nقبل أي اختيار، صنّف القيمة: محظورة، اتاخدت قبل كده، أو جديدة.\n\nبعدها راجع هل القرار يخص العنصر الحالي ولا بقية الصف."),
		ch38Question("c382.break-scope", "أي عبارة بتوقف الحلقة الداخلية بس؟", ch38QuestionDetails{Options: []string{"~break Slots~", "~continue Decks~", "~return~"}, Answer: 0, Feedback: []string{"الـbreak المسمّى ينهي الحلقة دي.", "الـcontinue ينقل المرور للحلقة الكبيرة.", "return يخرج من الدالة كلها."}}),
	}},
	{ID: "c382.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c382.turn", "ابدأ بتمثيل العضوية، وبعدها راجع نطاق كل label على صفوف صغيرة. خليك محافظ على ترتيب المدخل.")}},
}

var ch38SetsEnglish = []LessonStep{
	{ID: "c382.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c382.mission", "Choose values from rows while excluding blocked and repeated values.\n\nRow and candidate order matters to the result.\n\nYou will distinguish skipping one item from skipping its whole row.")}},
	{ID: "c382.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c382.set-membership", "A map can represent a set when its keys represent membership.\n\nThe value ~struct{}{}~ is empty; the key and comma-ok result matter.\n\nThis example checks illustrative names.", `func main() {
	labels := map[string]struct{}{}
	labels["north"] = struct{}{}
	_, exists := labels["north"]
	fmt.Println(exists, len(labels))
}`, ch38WalkthroughDetails{ExampleID: "e382-set", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Create an empty set.", Variables: []LessonVariable{{Name: "labels", Value: "map[]"}}}, {Line: 3, Caption: "Add north once.", Variables: []LessonVariable{{Name: "labels", Value: "map[north:{}]"}}}, {Line: 4, Caption: "Check membership."}, {Line: 5, Caption: "The key exists.", Variables: []LessonVariable{{Name: "exists", Value: "true"}}, Output: "true 1"}}})}},
	{ID: "c382.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c382.predict-label", "Predict the output:\n\n~~~go\nfunc main() {\nDecks:\n\tfor _, deck := range [][]int{{1, 0}, {2, 3}} {\n\t\tfor _, slot := range deck {\n\t\t\tif slot == 0 { continue Decks }\n\t\t\tfmt.Println(slot)\n\t\t}\n\t}\n}\n~~~", ch38QuestionDetails{Options: []string{"1, then 2 and 3", "1, 0, 2, and 3", "Only 2 and 3"}, Answer: 0, Feedback: []string{"The row prints 1; zero skips the rest of that row.", "Zero is not printed before continue.", "Value 1 prints before the row is skipped."}}),
		ch38Walkthrough("c382.label-reveal", "Follow the label's effect on the outer loop.", `func main() {
Decks:
	for _, deck := range [][]int{{1, 0}, {2, 3}} {
		for _, slot := range deck {
			if slot == 0 { continue Decks }
			fmt.Println(slot)
		}
	}
}`, ch38WalkthroughDetails{ExampleID: "e382-label", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 3, Caption: "Start the first row."}, {Line: 6, Caption: "Print its first slot.", Variables: []LessonVariable{{Name: "slot", Value: "1"}}, Output: "1"}, {Line: 5, Caption: "The label skips this row.", Output: "1"}, {Line: 3, Caption: "Start the next row.", Output: "1"}, {Line: 6, Caption: "Print its first slot.", Variables: []LessonVariable{{Name: "slot", Value: "2"}}, Output: "1\n2"}, {Line: 6, Caption: "Print its next slot.", Variables: []LessonVariable{{Name: "slot", Value: "3"}}, Output: "1\n2\n3"}}}),
	}},
	{ID: "c382.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c382.watch", "- A map is unordered; get result order from the row traversal.\n- A label names the loop to skip; check that it encloses the branch.")}},
	{ID: "c382.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c382.plan", "Sketch two loops: one for rows and one for candidates.\n\nBefore choosing, classify a value as blocked, already used, or new.\n\nThen check whether each decision affects one candidate or the rest of its row."),
		ch38Question("c382.break-scope", "Which statement ends only the inner loop?", ch38QuestionDetails{Options: []string{"~break Slots~", "~continue Decks~", "~return~"}, Answer: 0, Feedback: []string{"A named break ends that loop.", "continue advances the outer loop.", "return exits the whole function."}}),
	}},
	{ID: "c382.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c382.turn", "Start with membership sets, then check each label's scope on small rows. Preserve the input order.")}},
}

var ch38OnceArabic = []LessonStep{
	{ID: "c383.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c383.mission", "هيّأ كل قيمة مرة واحدة، حتى لو كذا goroutine طلبها مع بعض.\n\nالقيمة نفسها ممكن تكون صفر أو نص فاضي.\n\nافصل بين وقت إنشاء دالة التحميل ووقت أول استخدام.")}},
	{ID: "c383.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c383.once", "~sync.Once~ بتربط حالة التنفيذ بinstance واحدة.\n\nكل ~Do~ بعد أول تشغيل بينتظر لو التهيئة لسه شغالة، وبعدها يكمل.\n\nالمثال بيعدّ استدعاء الـfunction، مش قيمة الإعدادات.", `func main() {
	var once sync.Once
	calls := 0
	for i := 0; i < 2; i++ {
		once.Do(func() { calls++ })
	}
	fmt.Println(calls)
}`, ch38WalkthroughDetails{ExampleID: "e383-once", Imports: []string{"fmt", "sync"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "Once جاهزة بصفرها."}, {Line: 3, Caption: "العداد يبدأ بصفر.", Variables: []LessonVariable{{Name: "calls", Value: "0"}}}, {Line: 4, Caption: "ابدأ اللفة الأولى."}, {Line: 5, Caption: "نفّذ التهيئة مرة واحدة.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}}, {Line: 5, Caption: "النداء التالي لا يعيدها.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}}, {Line: 7, Caption: "اطبع عدد مرات التهيئة.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}, Output: "1"}}})}},
	{ID: "c383.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c383.predict-lazy", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tload := sync.OnceValue(func() string {\n\t\tfmt.Println(\"load\")\n\t\treturn \"east\"\n\t})\n\tfmt.Println(\"created\")\n\tfmt.Println(load())\n\tfmt.Println(load())\n}\n~~~", ch38QuestionDetails{Options: []string{"created, load, east, east", "load, created, east, east", "created, east, east"}, Answer: 0, Feedback: []string{"التحميل يبدأ عند أول نداء؛ النتيجة تتخزن.", "إنشاء الـclosure بيحصل قبل طباعة created.", "الـfunction بتطبع load وقت أول استدعاء."}}),
		ch38Walkthrough("c383.lazy-reveal", "تابع الإنشاء، وبعده النداءين.", `func main() {
	load := sync.OnceValue(func() string {
		fmt.Println("load")
		return "east"
	})
	fmt.Println("created")
	fmt.Println(load())
	fmt.Println(load())
}`, ch38WalkthroughDetails{ExampleID: "e383-lazy", Imports: []string{"fmt", "sync"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "أنشئ closure مؤجلة."}, {Line: 6, Caption: "اطبع بعد الإنشاء.", Output: "created"}, {Line: 7, Caption: "أول نداء يشغّل load.", Output: "created\nload\neast", Variables: []LessonVariable{{Name: "result", Value: "east"}}}, {Line: 8, Caption: "النداء التالي يرجّع المخزّن.", Output: "created\nload\neast\neast", Variables: []LessonVariable{{Name: "result", Value: "east"}}}}}),
	}},
	{ID: "c383.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c383.watch", "- ماتنسخش ~sync.Once~ بعد أول استخدام؛ خليه مع الحالة اللي بيحميها.\n- الصفر مش علامة موثوقة على إن التحميل لسه ماحصلش.")}},
	{ID: "c383.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c383.plan", "حدّد مين مالك حالة التهيئة، وهل حالته مستقلة عن الكائنات التانية.\n\nراجع إن كل caller يشوف نفس النتيجة بعد اكتمال التحميل.\n\nمع ~OnceValue~، إنشاء closure لوحده مش استهلاك."),
		ch38Question("c383.zero-result", "لو النتيجة ~0~، هل ده يثبت إن التهيئة ماحصلتش؟", ch38QuestionDetails{Options: []string{"لأ؛ القيمة وحالة التشغيل حاجتين", "أيوه؛ الصفر معناه مفيش تحميل", "بس لو مفيش goroutines"}, Answer: 0, Feedback: []string{"Once بتتبع التنفيذ منفصل عن قيمة النتيجة.", "التحميل ممكن يرجع القيمة الصفرية بشكل صحيح.", "قيمة النتيجة مش مرتبطة بعدد callers."}}),
	}},
	{ID: "c383.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c383.turn", "ابدأ بحالة مرة واحدة لكل instance، وبعدها افصل الـclosure الكسولة. جرّب القيمة الصفرية قبل اختبار التزامن.")}},
}

var ch38OnceEnglish = []LessonStep{
	{ID: "c383.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c383.mission", "Initialize each value once, even when several goroutines request it together.\n\nThe value itself may be zero or an empty string.\n\nSeparate creating the loader from its first use.")}},
	{ID: "c383.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c383.once", "~sync.Once~ ties execution state to one instance.\n\nEach later ~Do~ waits if initialization is still running, then continues.\n\nThis example counts calls to the function, not a setting value.", `func main() {
	var once sync.Once
	calls := 0
	for i := 0; i < 2; i++ {
		once.Do(func() { calls++ })
	}
	fmt.Println(calls)
}`, ch38WalkthroughDetails{ExampleID: "e383-once", Imports: []string{"fmt", "sync"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Once is ready at zero."}, {Line: 3, Caption: "The counter starts at zero.", Variables: []LessonVariable{{Name: "calls", Value: "0"}}}, {Line: 4, Caption: "Start the first loop."}, {Line: 5, Caption: "Initialize one time.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}}, {Line: 5, Caption: "The next call does not repeat.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}}, {Line: 7, Caption: "Print initialization calls.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}, Output: "1"}}})}},
	{ID: "c383.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c383.predict-lazy", "Predict the output:\n\n~~~go\nfunc main() {\n\tload := sync.OnceValue(func() string {\n\t\tfmt.Println(\"load\")\n\t\treturn \"east\"\n\t})\n\tfmt.Println(\"created\")\n\tfmt.Println(load())\n\tfmt.Println(load())\n}\n~~~", ch38QuestionDetails{Options: []string{"created, load, east, east", "load, created, east, east", "created, east, east"}, Answer: 0, Feedback: []string{"Loading starts on the first call; its result is cached.", "The closure is created before created prints.", "The function prints load on its first call."}}),
		ch38Walkthrough("c383.lazy-reveal", "Follow construction and both calls.", `func main() {
	load := sync.OnceValue(func() string {
		fmt.Println("load")
		return "east"
	})
	fmt.Println("created")
	fmt.Println(load())
	fmt.Println(load())
}`, ch38WalkthroughDetails{ExampleID: "e383-lazy", Imports: []string{"fmt", "sync"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Create a deferred closure."}, {Line: 6, Caption: "Print after construction.", Output: "created"}, {Line: 7, Caption: "First call runs load.", Output: "created\nload\neast", Variables: []LessonVariable{{Name: "result", Value: "east"}}}, {Line: 8, Caption: "Next call returns the cache.", Output: "created\nload\neast\neast", Variables: []LessonVariable{{Name: "result", Value: "east"}}}}}),
	}},
	{ID: "c383.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c383.watch", "- Do not copy ~sync.Once~ after first use; keep it with the state it protects.\n- Zero is not a reliable signal that loading has not happened.")}},
	{ID: "c383.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c383.plan", "Decide who owns initialization state and whether each object has its own state.\n\nCheck that callers see the same result after loading completes.\n\nWith ~OnceValue~, creating the closure is not consuming it."),
		ch38Question("c383.zero-result", "If the result is ~0~, does that prove initialization did not run?", ch38QuestionDetails{Options: []string{"No; value and run state differ", "Yes; zero means no load", "Only when there are no goroutines"}, Answer: 0, Feedback: []string{"Once tracks execution separately from the result.", "A load can validly return the zero value.", "The result does not depend on caller count."}}),
	}},
	{ID: "c383.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c383.turn", "Start with one-time state per instance, then separate the lazy closure. Try a zero result before checking concurrency.")}},
}

var ch38ErrgroupArabic = []LessonStep{
	{ID: "c384.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c384.mission", "شغّل jobs بالتوازي مع حد واضح، وتعامل مع أول خطأ.\n\nالإلغاء تعاوني: كل worker لازم يراجع الـcontext.\n\nالرجوع يستنى العمال اللي بدأوا بالفعل.")}},
	{ID: "c384.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c384.group", "~errgroup~ بتجمع انتظار workers وأول error.\n\n~SetLimit~ يحدد عدد المهام النشطة؛ ~Wait~ يستنى اكتمال كل اللي بدأ.\n\nالـframes بتعرض ترتيب محتمل واحد؛ تشغيل تاني ممكن يرتّب البداية بشكل مختلف.", `func main() {
	var total atomic.Int64
	var group errgroup.Group
	group.SetLimit(2)
	group.Go(func() error { total.Add(1); return nil })
	group.Go(func() error { total.Add(1); return nil })
	err := group.Wait()
	fmt.Println(err == nil, total.Load())
}`, ch38WalkthroughDetails{ExampleID: "e384-group", Imports: []string{"fmt", "sync/atomic", "golang.org/x/sync/errgroup"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "العداد آمن للتزامن."}, {Line: 3, Caption: "أنشئ group للشغل."}, {Line: 4, Caption: "اسمح باثنين نشطين."}, {Line: 5, Caption: "ترتيب محتمل: يبدأ worker."}, {Line: 6, Caption: "ممكن worker تاني يتداخل."}, {Line: 7, Caption: "Wait يستنى الاتنين."}, {Line: 8, Caption: "اطبع بعد اكتمالهما.", Variables: []LessonVariable{{Name: "total", Value: "2"}}, Output: "true 2"}}})}},
	{ID: "c384.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c384.predict-wait", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tvar group errgroup.Group\n\tgroup.Go(func() error { return nil })\n\terr := group.Wait()\n\tfmt.Println(err == nil)\n}\n~~~", ch38QuestionDetails{Options: []string{"true", "false", "البرنامج يعلق"}, Answer: 0, Feedback: []string{"كل worker خلص من غير error.", "Wait بيرجع nil لما كل الشغل ينجح.", "مفيش حد توازي يمنع المهمة من البدء."}}),
		ch38Walkthrough("c384.wait-reveal", "تابع إرسال worker ثم الانتظار.", `func main() {
	var group errgroup.Group
	group.Go(func() error { return nil })
	err := group.Wait()
	fmt.Println(err == nil)
}`, ch38WalkthroughDetails{ExampleID: "e384-wait", Imports: []string{"fmt", "golang.org/x/sync/errgroup"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "جهّز group فاضية."}, {Line: 3, Caption: "ابدأ worker واحدة."}, {Line: 4, Caption: "Wait يجمع النتيجة."}, {Line: 5, Caption: "مافيش error.", Variables: []LessonVariable{{Name: "err", Value: "nil"}}, Output: "true"}}}),
	}},
	{ID: "c384.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c384.watch", "- ماتعتمدش على ترتيب goroutines؛ اعتمد على النتيجة بعد ~Wait~.\n- الـcontext بتطلب من الشغل يوقف؛ مش بتقتل worker بالقوة.")}},
	{ID: "c384.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c384.plan", "راجع حدود التوازي قبل الجدولة، وبعدها اسأل إزاي كل worker بتلاحظ الإلغاء.\n\nافصل خطأ الأب عن context المشتقة من المجموعة.\n\nبعد الإلغاء، تأكد إن كل worker بدأت خلصت قبل الرجوع."),
		ch38Question("c384.cancel", "مين بيوقف worker بعد إلغاء الـcontext؟", ch38QuestionDetails{Options: []string{"الـworker لما تراجعها", "errgroup بالقوة", "SetLimit لوحدها"}, Answer: 0, Feedback: []string{"الإلغاء إشارة؛ callback لازم يتعاون.", "المكتبة ما تقدرش تقتل goroutine بالقوة.", "الحد بيقيّد العدد ومش بيلغي الشغل."}}),
	}},
	{ID: "c384.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c384.turn", "ابدأ بالنجاح على jobs صغيرة، وبعدها راجع حد التوازي والإلغاء. استخدم ~-race~ عشان تراجع الحالة المشتركة.")}},
}

var ch38ErrgroupEnglish = []LessonStep{
	{ID: "c384.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c384.mission", "Run jobs in parallel with a clear limit, and handle the first error.\n\nCancellation is cooperative: each worker must check its context.\n\nReturn only after workers that already started finish.")}},
	{ID: "c384.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c384.group", "~errgroup~ combines worker waiting and the first error.\n\n~SetLimit~ caps active tasks; ~Wait~ joins everything already started.\n\nFrames show one possible order; another run may schedule starts differently.", `func main() {
	var total atomic.Int64
	var group errgroup.Group
	group.SetLimit(2)
	group.Go(func() error { total.Add(1); return nil })
	group.Go(func() error { total.Add(1); return nil })
	err := group.Wait()
	fmt.Println(err == nil, total.Load())
}`, ch38WalkthroughDetails{ExampleID: "e384-group", Imports: []string{"fmt", "sync/atomic", "golang.org/x/sync/errgroup"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "The counter is concurrent-safe."}, {Line: 3, Caption: "Create a work group."}, {Line: 4, Caption: "Allow two active tasks."}, {Line: 5, Caption: "One possible worker starts."}, {Line: 6, Caption: "Another may overlap."}, {Line: 7, Caption: "Wait joins both."}, {Line: 8, Caption: "Print after both finish.", Variables: []LessonVariable{{Name: "total", Value: "2"}}, Output: "true 2"}}})}},
	{ID: "c384.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c384.predict-wait", "Predict the output:\n\n~~~go\nfunc main() {\n\tvar group errgroup.Group\n\tgroup.Go(func() error { return nil })\n\terr := group.Wait()\n\tfmt.Println(err == nil)\n}\n~~~", ch38QuestionDetails{Options: []string{"true", "false", "The program hangs"}, Answer: 0, Feedback: []string{"The worker finishes without an error.", "Wait returns nil when all work succeeds.", "No concurrency limit prevents the task from starting."}}),
		ch38Walkthrough("c384.wait-reveal", "Follow worker submission, then wait.", `func main() {
	var group errgroup.Group
	group.Go(func() error { return nil })
	err := group.Wait()
	fmt.Println(err == nil)
}`, ch38WalkthroughDetails{ExampleID: "e384-wait", Imports: []string{"fmt", "golang.org/x/sync/errgroup"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Prepare an empty group."}, {Line: 3, Caption: "Start one worker."}, {Line: 4, Caption: "Wait collects its result."}, {Line: 5, Caption: "There is no error.", Variables: []LessonVariable{{Name: "err", Value: "nil"}}, Output: "true"}}}),
	}},
	{ID: "c384.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c384.watch", "- Do not rely on goroutine order; rely on the result after ~Wait~.\n- A context asks work to stop; it cannot kill a worker by force.")}},
	{ID: "c384.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c384.plan", "Check the concurrency limit before scheduling, then ask how each worker notices cancellation.\n\nSeparate the parent error from the group's derived context.\n\nAfter cancellation, make sure every started worker finishes before returning."),
		ch38Question("c384.cancel", "Who stops a worker after its context is canceled?", ch38QuestionDetails{Options: []string{"The worker when it checks", "errgroup by force", "SetLimit by itself"}, Answer: 0, Feedback: []string{"Cancellation is a signal; the callback cooperates.", "The package cannot kill a goroutine by force.", "The limit caps task count; it does not cancel work."}}),
	}},
	{ID: "c384.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c384.turn", "Start with successful small jobs, then check the limit and cancellation. Use ~-race~ to inspect shared state.")}},
}

var ch38AllocsArabic = []LessonStep{
	{ID: "c385.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c385.mission", "قلّل تخصيصات النصوص والـslice، وبعدها قِس benchmark.\n\nالنتيجة لازم تفضل صحيحة ومستقلة عن المدخل.\n\nالقياس يوضح الفرق؛ ماتخمنش عدد التخصيصات من شكل الكود.")}},
	{ID: "c385.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c385.builder", "~strings.Builder~ بيجمع أجزاء النص في buffer.\n\n~Grow~ يحجز سعة متوقعة؛ الكتابة نفسها تحدد النص النهائي.\n\nالسعة وسيلة لتقليل النمو، مش جزء من الناتج.", `func main() {
	var b strings.Builder
	b.Grow(8)
	b.WriteString("ion")
	b.WriteByte('-')
	b.WriteString("nav")
	fmt.Println(b.String())
}`, ch38WalkthroughDetails{ExampleID: "e385-builder", Imports: []string{"fmt", "strings"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "أنشئ Builder فاضي."}, {Line: 3, Caption: "احجز سعة للنص."}, {Line: 4, Caption: "أضف الجزء الأول."}, {Line: 5, Caption: "أضف الفاصل."}, {Line: 6, Caption: "أضف الجزء التالي."}, {Line: 7, Caption: "النص الكامل يظهر.", Variables: []LessonVariable{{Name: "b.String()", Value: "ion-nav"}}, Output: "ion-nav"}}})}},
	{ID: "c385.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c385.predict-cap", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tvalues := make([]int, 0, 4)\n\tvalues = append(values, 7)\n\tfmt.Println(len(values), cap(values))\n}\n~~~", ch38QuestionDetails{Options: []string{"1 4", "4 1", "0 4"}, Answer: 0, Feedback: []string{"اتضاف عنصر واحد ولسه فيه سعة لأربعة.", "الطول هو عدد العناصر، مش السعة.", "append زوّدت الطول لواحد."}}),
		ch38Walkthrough("c385.capacity-reveal", "افصل بين عدد العناصر والسعة.", `func main() {
	values := make([]int, 0, 4)
	values = append(values, 7)
	fmt.Println(len(values), cap(values))
}`, ch38WalkthroughDetails{ExampleID: "e385-capacity", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "الطول صفر والسعة أربعة.", Variables: []LessonVariable{{Name: "values", Value: "[]"}}}, {Line: 3, Caption: "أضف عنصراً واحداً.", Variables: []LessonVariable{{Name: "values", Value: "[7]"}}}, {Line: 4, Caption: "اطبع الطول والسعة.", Variables: []LessonVariable{{Name: "len", Value: "1"}, {Name: "cap", Value: "4"}}, Output: "1 4"}}}),
	}},
	{ID: "c385.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c385.watch", "- Benchmark يقيس التنفيذ المتكرر؛ خَلّي بيانات التجهيز برّه القياس.\n- ~ReportAllocs~ يعرض allocations، لكنه مش ضمان لصحة النتيجة.")}},
	{ID: "c385.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c385.plan", "راجع طول الـslice وسعتها بشكل منفصل، وبعدها تأكد إن النسخة الجديدة مش بتشارك المدخل.\n\nفي النص، عدّ bytes والفواصل قبل الحجز من غير ما تغيّر المحتوى.\n\nخلّي الـbenchmark يقيس الدالة، مش تحضير البيانات."),
		ch38Question("c385.report-allocs", "إيه اللي بيضيفه ~b.ReportAllocs()~؟", ch38QuestionDetails{Options: []string{"أرقام تخصيصات للـbenchmark", "عدد العناصر في النتيجة", "نسخة من المدخل"}, Answer: 0, Feedback: []string{"التقرير يعرض allocations لكل عملية.", "الطول بيتقاس بـlen على القيمة.", "الدالة نفسها مسؤولة عن النسخ."}}),
	}},
	{ID: "c385.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c385.turn", "ثبّت السلوك الأول، وبعدها أضف benchmark مع sink. قارن النتيجة والتخصيصات قبل ما تهتم بالسرعة.")}},
}

var ch38AllocsEnglish = []LessonStep{
	{ID: "c385.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c385.mission", "Reduce string and slice allocations, then measure a benchmark.\n\nThe result must stay correct and independent from its input.\n\nMeasurement shows the difference; do not guess allocation counts from code shape.")}},
	{ID: "c385.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c385.builder", "~strings.Builder~ collects text pieces in a buffer.\n\n~Grow~ reserves expected capacity; writes determine the final text.\n\nCapacity helps reduce growth; it is not part of the output.", `func main() {
	var b strings.Builder
	b.Grow(8)
	b.WriteString("ion")
	b.WriteByte('-')
	b.WriteString("nav")
	fmt.Println(b.String())
}`, ch38WalkthroughDetails{ExampleID: "e385-builder", Imports: []string{"fmt", "strings"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Create an empty Builder."}, {Line: 3, Caption: "Reserve text capacity."}, {Line: 4, Caption: "Write the first piece."}, {Line: 5, Caption: "Write the separator."}, {Line: 6, Caption: "Write the next piece."}, {Line: 7, Caption: "The full text appears.", Variables: []LessonVariable{{Name: "b.String()", Value: "ion-nav"}}, Output: "ion-nav"}}})}},
	{ID: "c385.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c385.predict-cap", "Predict the output:\n\n~~~go\nfunc main() {\n\tvalues := make([]int, 0, 4)\n\tvalues = append(values, 7)\n\tfmt.Println(len(values), cap(values))\n}\n~~~", ch38QuestionDetails{Options: []string{"1 4", "4 1", "0 4"}, Answer: 0, Feedback: []string{"One item was appended, with room for four.", "Length counts elements, not capacity.", "append increased the length to one."}}),
		ch38Walkthrough("c385.capacity-reveal", "Separate element count from capacity.", `func main() {
	values := make([]int, 0, 4)
	values = append(values, 7)
	fmt.Println(len(values), cap(values))
}`, ch38WalkthroughDetails{ExampleID: "e385-capacity", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Length zero, capacity four.", Variables: []LessonVariable{{Name: "values", Value: "[]"}}}, {Line: 3, Caption: "Append one element.", Variables: []LessonVariable{{Name: "values", Value: "[7]"}}}, {Line: 4, Caption: "Print length and capacity.", Variables: []LessonVariable{{Name: "len", Value: "1"}, {Name: "cap", Value: "4"}}, Output: "1 4"}}}),
	}},
	{ID: "c385.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c385.watch", "- A benchmark repeats the measured work; keep setup outside the timed loop.\n- ~ReportAllocs~ reports allocations, but does not prove correctness.")}},
	{ID: "c385.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c385.plan", "Check slice length and capacity separately, then confirm the new copy does not share input storage.\n\nFor text, count bytes and separators before reserving space without changing content.\n\nKeep the benchmark focused on the function, not data setup."),
		ch38Question("c385.report-allocs", "What does ~b.ReportAllocs()~ add?", ch38QuestionDetails{Options: []string{"Allocation numbers for the benchmark", "The result's element count", "A copy of the input"}, Answer: 0, Feedback: []string{"The report shows allocations per operation.", "Length is measured with len on a value.", "The function itself must make the copy."}}),
	}},
	{ID: "c385.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c385.turn", "Lock down behavior first, then add a benchmark with a sink. Compare output and allocations before focusing on speed.")}},
}

var ch38VetArabic = []LessonStep{
	{ID: "c386.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c386.mission", "راجع starter بأدوات static analysis، وصلّح التحذيرات والسلوك مع بعض.\n\nفيه حالة مشتركة ورسالة formatted وإلغاء لازم يتنضف.\n\nالتست و~go vet~ بيكشفوا أنواع مختلفة من المشاكل.")}},
	{ID: "c386.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c386.format", "~go vet~ يلتقط استخدام format مش مناسب حتى لو الكود بيتبني.\n\n~%d~ بيطبع integer؛ راجع نوع القيمة مع verb.\n\nالإخراج هنا توضيحي، وvet نفسه مش بيشغّل البرنامج.", `func main() {
	count := -3
	fmt.Printf("count=%d\n", count)
}`, ch38WalkthroughDetails{ExampleID: "e386-format", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "القيمة عدد صحيح.", Variables: []LessonVariable{{Name: "count", Value: "-3"}}}, {Line: 3, Caption: "%d مناسب للعدد."}, {Line: 3, Caption: "اطبع القيمة بإشارتها.", Variables: []LessonVariable{{Name: "count", Value: "-3"}}, Output: "count=-3"}}})}},
	{ID: "c386.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c386.predict-format", "توقّع ناتج ~Printf~:\n\n~~~go\nfunc main() {\n\tfmt.Printf(\"items=%d\\n\", 12)\n}\n~~~", ch38QuestionDetails{Options: []string{"items=12", "items=%!d(string=12)", "12 items"}, Answer: 0, Feedback: []string{"%d بيطبع الرقم كعدد عشري.", "القيمة int وverb مناسب لنوعها.", "النص يطبع بالترتيب المكتوب."}}),
		ch38Walkthrough("c386.format-reveal", "راجع النوع والـformat قبل الطباعة.", `func main() {
	fmt.Printf("items=%d\n", 12)
}`, ch38WalkthroughDetails{ExampleID: "e386-printf", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "اطبع int باستخدام %d.", Output: "items=12"}}}),
	}},
	{ID: "c386.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c386.watch", "- Receiver بالقيمة بينسخ الـstruct؛ ده خطر لو فيها Mutex مستخدمة.\n- بعد ~WithCancel~ راجع إن كل مسار رجوع يوصل لـcancel.")}},
	{ID: "c386.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c386.plan", "شغّل الفحص على المصدر، واقرأ كل تشخيص في مكانه.\n\nبعدها راجع ملكية القفل والـCancelFunc، واختبر النجاح والخطأ.\n\n~-race~ بيساعد في الحالة المتزامنة؛ نجاحه مش بديل عن vet."),
		ch38Question("c386.vet-role", "إيه دور ~go vet~ هنا؟", ch38QuestionDetails{Options: []string{"يكشف أنماط مريبة في كود بيتبني", "يثبت إن السلوك كله صحيح", "يشغّل goroutines بترتيب ثابت"}, Answer: 0, Feedback: []string{"تحذيرات vet تكمل البناء والاختبارات.", "التشغيل والاختبارات لسه مهمين.", "التزامن ملوش ترتيب ثابت مضمون."}}),
	}},
	{ID: "c386.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c386.turn", "عالج تشخيص واحد كل مرة، وبعد كل تعديل شغّل الفحص والاختبارات. راجع إن خطأ read الأصلي مايتبدلش.")}},
}

var ch38VetEnglish = []LessonStep{
	{ID: "c386.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c386.mission", "Inspect the starter with static analysis, then fix both warnings and behavior.\n\nIt has shared state, a formatted message and cancellation to clean up.\n\nTests and ~go vet~ catch different kinds of problems.")}},
	{ID: "c386.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c386.format", "~go vet~ catches a mismatched format verb even when code builds.\n\n~%d~ prints an integer; match the value's type to its verb.\n\nThis output is illustrative; vet itself does not run the program.", `func main() {
	count := -3
	fmt.Printf("count=%d\n", count)
}`, ch38WalkthroughDetails{ExampleID: "e386-format", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "The value is an integer.", Variables: []LessonVariable{{Name: "count", Value: "-3"}}}, {Line: 3, Caption: "%d matches the number."}, {Line: 3, Caption: "Print the signed value.", Variables: []LessonVariable{{Name: "count", Value: "-3"}}, Output: "count=-3"}}})}},
	{ID: "c386.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c386.predict-format", "Predict the ~Printf~ output:\n\n~~~go\nfunc main() {\n\tfmt.Printf(\"items=%d\\n\", 12)\n}\n~~~", ch38QuestionDetails{Options: []string{"items=12", "items=%!d(string=12)", "12 items"}, Answer: 0, Feedback: []string{"%d prints an integer in decimal.", "The value is int and the verb matches.", "The text prints in its written order."}}),
		ch38Walkthrough("c386.format-reveal", "Check the type and format before printing.", `func main() {
	fmt.Printf("items=%d\n", 12)
}`, ch38WalkthroughDetails{ExampleID: "e386-printf", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Print int with %d.", Output: "items=12"}}}),
	}},
	{ID: "c386.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c386.watch", "- A value receiver copies its struct; that is unsafe with a used Mutex inside.\n- After ~WithCancel~, check that every return path reaches cancel.")}},
	{ID: "c386.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c386.plan", "Run analysis on the source and read each diagnostic at its location.\n\nThen check lock and CancelFunc ownership, and exercise success and error paths.\n\n~-race~ helps with shared state; passing it does not replace vet."),
		ch38Question("c386.vet-role", "What does ~go vet~ do here?", ch38QuestionDetails{Options: []string{"Find suspicious patterns in buildable code", "Prove all behavior is correct", "Run goroutines in a fixed order"}, Answer: 0, Feedback: []string{"Vet diagnostics complement builds and tests.", "Execution and tests still matter.", "Concurrency has no guaranteed fixed order."}}),
	}},
	{ID: "c386.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c386.turn", "Address one diagnostic at a time, then rerun analysis and tests after each change. Preserve the original error from read.")}},
}

var ch38PprofArabic = []LessonStep{
	{ID: "c387.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch38Markdown("c387.mission", "قِس وقت التنفيذ والذاكرة بprofile حقيقي، وحسّن جمع القيم المميزة.\n\nCPU وheap بيجاوبوا على أسئلة مختلفة.\n\nقارن الخوارزمية من غير تكلفة تشغيل profiler جوه benchmark.")}},
	{ID: "c387.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch38Walkthrough("c387.heap-profile", "~pprof.WriteHeapProfile~ بيكتب بيانات profile في ~io.Writer~.\n\nالمحتوى مضغوط ومش نص للقراءة المباشرة.\n\nهنا بنفحص وجود بيانات بس، مش حجم ثابت.", `func main() {
	var profile bytes.Buffer
	err := pprof.WriteHeapProfile(&profile)
	fmt.Println(err == nil, profile.Len() > 0)
}`, ch38WalkthroughDetails{ExampleID: "e387-heap", Imports: []string{"bytes", "fmt", "runtime/pprof"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "جهّز writer في الذاكرة."}, {Line: 3, Caption: "اكتب heap profile."}, {Line: 4, Caption: "اتأكد من النجاح ووجود بيانات.", Variables: []LessonVariable{{Name: "err", Value: "nil"}, {Name: "profile.Len() > 0", Value: "true"}}, Output: "true true"}}})}},
	{ID: "c387.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch38Question("c387.predict-set", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tseen := map[int]struct{}{4: {}}\n\t_, exists := seen[4]\n\tfmt.Println(exists)\n}\n~~~", ch38QuestionDetails{Options: []string{"true", "false", "4"}, Answer: 0, Feedback: []string{"المفتاح 4 موجود في الـset.", "القراءة comma-ok بترجع true للمفتاح الموجود.", "fmt هنا بيطبع نتيجة الفحص، مش المفتاح."}}),
		ch38Walkthrough("c387.set-reveal", "راجع العضوية بمفتاح موجود.", `func main() {
	seen := map[int]struct{}{4: {}}
	_, exists := seen[4]
	fmt.Println(exists)
}`, ch38WalkthroughDetails{ExampleID: "e387-set", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "ابدأ main."}, {Line: 2, Caption: "set فيها القيمة 4."}, {Line: 3, Caption: "افحص المفتاح."}, {Line: 4, Caption: "النتيجة موجودة.", Variables: []LessonVariable{{Name: "exists", Value: "true"}}, Output: "true"}}}),
	}},
	{ID: "c387.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch38Markdown("c387.watch", "- بعد نجاح ~StartCPUProfile~ لازم توقفها؛ فيه profiler واحدة للـprocess.\n- CPU profile أخذ عينات، والـheap profile كمان له دلالته؛ مافيش حجم ملف مضمون.")}},
	{ID: "c387.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch38Markdown("c387.plan", "ابدأ بفصل حساب المجموع عن جمع القياسات.\n\nفي المجموع، تابع العضوية من غير تعديل القراءات.\n\nراجع ترتيب بدء CPU وإيقافه وكتابة heap، وافحص كل error متاح من الـAPI.\n\nبعدها قارن الأداء على input كبير من غير profiling جوه القياس."),
		ch38Question("c387.profile-types", "أي profile بيأخذ عينات من stacks وقت التشغيل؟", ch38QuestionDetails{Options: []string{"CPU profile", "Heap profile", "Benchmark report بس"}, Answer: 0, Feedback: []string{"CPU profile بيربط العينات بمسارات التنفيذ.", "heap بيركز على تخصيصات الذاكرة.", "التقرير مش profile CPU لوحده."}}),
	}},
	{ID: "c387.turn-step", Title: "دورك", Blocks: []LessonBlock{ch38Markdown("c387.turn", "ابدأ بصحة المجموع، وبعدها اختبر writers والـcleanup. قارن الـbenchmark منفصل عن دورة كتابة profile.")}},
}

var ch38PprofEnglish = []LessonStep{
	{ID: "c387.mission-step", Title: "Mission", Blocks: []LessonBlock{ch38Markdown("c387.mission", "Measure execution and memory with real profiles, and improve distinct-value totals.\n\nCPU and heap answer different questions.\n\nCompare the algorithm without profiler startup inside the benchmark.")}},
	{ID: "c387.idea-step", Title: "Idea", Blocks: []LessonBlock{ch38Walkthrough("c387.heap-profile", "~pprof.WriteHeapProfile~ writes profile data to an ~io.Writer~.\n\nThe contents are compressed, not plain text for reading.\n\nThis checks that data exists, not for a fixed size.", `func main() {
	var profile bytes.Buffer
	err := pprof.WriteHeapProfile(&profile)
	fmt.Println(err == nil, profile.Len() > 0)
}`, ch38WalkthroughDetails{ExampleID: "e387-heap", Imports: []string{"bytes", "fmt", "runtime/pprof"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "Prepare an in-memory writer."}, {Line: 3, Caption: "Write a heap profile."}, {Line: 4, Caption: "Check success and data.", Variables: []LessonVariable{{Name: "err", Value: "nil"}, {Name: "profile.Len() > 0", Value: "true"}}, Output: "true true"}}})}},
	{ID: "c387.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch38Question("c387.predict-set", "Predict the output:\n\n~~~go\nfunc main() {\n\tseen := map[int]struct{}{4: {}}\n\t_, exists := seen[4]\n\tfmt.Println(exists)\n}\n~~~", ch38QuestionDetails{Options: []string{"true", "false", "4"}, Answer: 0, Feedback: []string{"Key 4 exists in the set.", "A comma-ok lookup is true for an existing key.", "fmt prints the lookup result, not the key."}}),
		ch38Walkthrough("c387.set-reveal", "Check membership for an existing key.", `func main() {
	seen := map[int]struct{}{4: {}}
	_, exists := seen[4]
	fmt.Println(exists)
}`, ch38WalkthroughDetails{ExampleID: "e387-set", Imports: []string{"fmt"}, Frames: []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Enter main."}, {Line: 2, Caption: "The set contains 4."}, {Line: 3, Caption: "Check the key."}, {Line: 4, Caption: "The result exists.", Variables: []LessonVariable{{Name: "exists", Value: "true"}}, Output: "true"}}}),
	}},
	{ID: "c387.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch38Markdown("c387.watch", "- After ~StartCPUProfile~ succeeds, stop it; a process can have only one CPU profiler.\n- CPU samples and heap samples mean different things; there is no fixed profile size.")}},
	{ID: "c387.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch38Markdown("c387.plan", "First separate the total calculation from profile collection.\n\nFor the total, track membership without changing readings.\n\nCheck CPU start, stop and heap write order, and inspect each API error you can receive.\n\nThen compare on larger input without profiling inside the measurement."),
		ch38Question("c387.profile-types", "Which profile samples execution stacks while the program runs?", ch38QuestionDetails{Options: []string{"CPU profile", "Heap profile", "Benchmark report only"}, Answer: 0, Feedback: []string{"CPU samples are tied to execution paths.", "Heap focuses on memory allocations.", "The report alone is not a CPU profile."}}),
	}},
	{ID: "c387.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch38Markdown("c387.turn", "Start with total correctness, then check writers and cleanup. Benchmark separately from profile writing.")}},
}
