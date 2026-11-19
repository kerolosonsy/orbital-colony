package levels

func ch03Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch03Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch03Walkthrough(id, text, code, exampleID string, frames []LessonFrame) LessonBlock {
	return LessonBlock{
		ID:          id,
		Kind:        "walkthrough",
		Text:        text,
		Code:        code,
		Imports:     []string{"fmt"},
		ExampleID:   exampleID,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch03Locale(ar, en string) map[string]any {
	return map[string]any{"ar": ar, "en": en}
}

func ch03Sandbox(id, text, kind string, config map[string]any, english bool) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: ch03ConfigForLanguage(config, english).(map[string]any)}}
}

func ch03ConfigForLanguage(value any, english bool) any {
	switch current := value.(type) {
	case map[string]any:
		ar, hasArabic := current["ar"].(string)
		en, hasEnglish := current["en"].(string)
		if len(current) == 2 && hasArabic && hasEnglish {
			if english {
				return map[string]any{"ar": "", "en": en}
			}
			return map[string]any{"ar": ar, "en": ""}
		}
		localized := make(map[string]any, len(current))
		for key, nested := range current {
			localized[key] = ch03ConfigForLanguage(nested, english)
		}
		return localized
	case []any:
		localized := make([]any, len(current))
		for i, nested := range current {
			localized[i] = ch03ConfigForLanguage(nested, english)
		}
		return localized
	default:
		return current
	}
}

var ch03RepeatConfig = map[string]any{
	"labels": map[string]any{
		"counter": ch03Locale("العداد i", "Counter i"),
		"prints":  ch03Locale("الناتج", "Output"),
		"next":    ch03Locale("الخطوة الجاية", "Next step"),
		"reset":   ch03Locale("من الأول", "Reset"),
	},
	"frames": []any{
		map[string]any{"i": 1, "prints": []int{}, "caption": ch03Locale("قبل أول طبعة.", "Before the first print.")},
		map[string]any{"i": 1, "prints": []int{1}, "caption": ch03Locale("الجسم يطبع القيمة 1.", "The body prints value 1.")},
		map[string]any{"i": 2, "prints": []int{1}, "caption": ch03Locale("العداد يزيد إلى 2.", "The counter advances to 2.")},
		map[string]any{"i": 2, "prints": []int{1, 2}, "caption": ch03Locale("اللفة التانية تطبع 2.", "The second loop prints 2.")},
		map[string]any{"i": 3, "prints": []int{1, 2}, "caption": ch03Locale("العداد يزيد إلى 3.", "The counter advances to 3.")},
		map[string]any{"i": 3, "prints": []int{1, 2, 3}, "caption": ch03Locale("اللفة الأخيرة تطبع 3.", "The final loop prints 3.")},
		map[string]any{"i": 4, "prints": []int{1, 2, 3}, "caption": ch03Locale("عند 4 الشرط يقف.", "At 4, the condition stops.")},
	},
}

var ch03NestedConfig = map[string]any{
	"rows": 2, "columns": 2,
	"labels": map[string]any{
		"position": ch03Locale("الصف والعمود", "Row and column"),
		"next":     ch03Locale("الخطوة الجاية", "Next step"),
		"reset":    ch03Locale("من الأول", "Reset"),
	},
	"frames": []any{
		map[string]any{"row": 0, "column": 0, "visited": []int{}, "caption": ch03Locale("ابدأ الصف الأول.", "Start the first row.")},
		map[string]any{"row": 0, "column": 0, "visited": []int{0}, "caption": ch03Locale("ابدأ العمود من صفر.", "Start the column at zero.")},
		map[string]any{"row": 0, "column": 1, "visited": []int{0, 1}, "caption": ch03Locale("العمود الداخلي يتحرك.", "The inner column advances.")},
		map[string]any{"row": 1, "column": 0, "visited": []int{0, 1, 2}, "caption": ch03Locale("صف جديد؛ العمود يرجع.", "New row; the column resets.")},
		map[string]any{"row": 1, "column": 1, "visited": []int{0, 1, 2, 3}, "caption": ch03Locale("آخر خانة في الشبكة.", "The last grid cell.")},
	},
}

var ch03ConditionConfig = map[string]any{
	"limit": 3,
	"labels": map[string]any{
		"count":     ch03Locale("العدد", "Count"),
		"condition": ch03Locale("الشرط", "Condition"),
		"running":   ch03Locale("true · شغّال", "true · running"),
		"stopped":   ch03Locale("false · وقف", "false · stopped"),
		"next":      ch03Locale("الخطوة الجاية", "Next step"),
		"reset":     ch03Locale("من الأول", "Reset"),
	},
	"frames": []any{
		map[string]any{"count": 0, "runs": 0, "running": true, "caption": ch03Locale("ابدأ والعدد صفر.", "Start with a zero count.")},
		map[string]any{"count": 1, "runs": 1, "running": true, "caption": ch03Locale("الجسم زوّد العدد.", "The body increased the count.")},
		map[string]any{"count": 2, "runs": 2, "running": true, "caption": ch03Locale("الشرط لسه true.", "The condition is still true.")},
		map[string]any{"count": 3, "runs": 3, "running": false, "caption": ch03Locale("الشرط بقى false.", "The condition became false.")},
	},
}

var ch03LongRowArabic = []LessonStep{
	{ID: "c31.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch03Markdown("c31.mission", "1. احصد الـ14 نبتة طحلب في الصف الطويل.\n2. الـstarter بيحرّك شرق 3 مرات ومش بيحصد.\n3. هدف الذهب 28 ticks؛ خلّي كل تكرار يعمل شغل مفيد.")}},
	{ID: "c31.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch03Walkthrough("c31.loop-trace", "الـfor loop فيها بداية وشرط وتغيير للعداد.\n\nالجسم بيتكرر طول ما الشرط true. الأرقام هنا للتوضيح.", `func main() {
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
	}
}`, "e31-loop-counter", []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "i تبدأ من 1.", Variables: []LessonVariable{{Name: "i", Value: "1"}}},
		{Line: 3, Caption: "اطبع أول value.", Variables: []LessonVariable{{Name: "i", Value: "1"}}, Output: "1"},
		{Line: 2, Caption: "اللفة التالية عند 2.", Variables: []LessonVariable{{Name: "i", Value: "2"}}, Output: "1"},
		{Line: 3, Caption: "اطبع الـvalue التانية.", Variables: []LessonVariable{{Name: "i", Value: "2"}}, Output: "1\n2"},
		{Line: 2, Caption: "اللفة الأخيرة عند 3.", Variables: []LessonVariable{{Name: "i", Value: "3"}}, Output: "1\n2"},
		{Line: 3, Caption: "اطبع آخر value.", Variables: []LessonVariable{{Name: "i", Value: "3"}}, Output: "1\n2\n3"},
		{Line: 2, Caption: "الشرط يقف عند 4.", Variables: []LessonVariable{{Name: "i", Value: "4"}}, Output: "1\n2\n3"},
	})}},
	{ID: "c31.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch03Question("c31.predict-loop", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tfor n := 0; n < 2; n++ {\n\t\tfmt.Println(\"step\")\n\t}\n}\n~~~", []string{"step مرة", "step مرتين", "step 3 مرات"}, 1, []string{
			"الشرط بيسمح بلفتين قبل ما يقف.",
			"العداد 0 وبعدين 1؛ يعني طباعتين.",
			"لما n تبقى 2، الشرط بيبقى false.",
		}),
		ch03Walkthrough("c31.loop-reveal", "تابع العداد مع كل لفة.", `func main() {
	for n := 0; n < 2; n++ {
		fmt.Println("step")
	}
}`, "e31-loop-reveal", []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "العداد بيبدأ بصفر.", Variables: []LessonVariable{{Name: "n", Value: "0"}}},
			{Line: 3, Caption: "اطبع أول مرة.", Variables: []LessonVariable{{Name: "n", Value: "0"}}, Output: "step"},
			{Line: 2, Caption: "العداد بقى 1.", Variables: []LessonVariable{{Name: "n", Value: "1"}}, Output: "step"},
			{Line: 3, Caption: "اطبع المرة التانية.", Variables: []LessonVariable{{Name: "n", Value: "1"}}, Output: "step\nstep"},
			{Line: 2, Caption: "الشرط وقف عند 2.", Variables: []LessonVariable{{Name: "n", Value: "2"}}, Output: "step\nstep"},
		}),
		ch03Sandbox("ch03-loop-counter", "اضغط وتابع العداد والناتج.", "ch03-repeat-counter", ch03RepeatConfig, false),
	},
	},
	{ID: "c31.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch03Markdown("c31.watch", "- حط الأوامر اللي عايز تكررها جوه ~{ }~.\n- شرط غلط ممكن يكرر أكتر أو أقل من المطلوب.")}},
	{ID: "c31.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch03Markdown("c31.plan", "راجع كل لفة: العداد بيتغير؟ والدرون عمل إيه؟\n\nقبل آخر حركة، بص على مكان الدرون عشان ما تطلبش خطوة زيادة."),
		ch03Question("c31.loop-count", "كام مرة الجسم يشتغل مع ~i := 0; i < 3; i++~؟", []string{"مرتين", "3 مرات", "4 مرات"}, 1, []string{
			"الـvalues 0 و1 و2 كلها أقل من 3.",
			"الـvalues 0 و1 و2 بتحقق الشرط.",
			"عند 3 الشرط يبقى false، فالجسم ما يشتغلش.",
		}),
	}},
	{ID: "c31.turn-step", Title: "دورك", Blocks: []LessonBlock{ch03Markdown("c31.turn", "عدّل الـstarter وخلي عندك لفة واحدة مفهومة. شغّلها، وبعدها زوّد التكرار وراقب مكان الدرون والـinventory.")}},
}

var ch03LongRowEnglish = []LessonStep{
	{ID: "c31.mission-step", Title: "Mission", Blocks: []LessonBlock{ch03Markdown("c31.mission", "1. Harvest all 14 moss plants in the long row.\n2. The starter moves east three times and harvests nothing.\n3. Gold is 28 ticks; make each repetition do useful work.")}},
	{ID: "c31.idea-step", Title: "Idea", Blocks: []LessonBlock{ch03Walkthrough("c31.loop-trace", "A for loop has a start, a condition, and a counter update.\n\nIts body repeats while the condition is true. These values are illustrative.", `func main() {
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
	}
}`, "e31-loop-counter", []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "i starts at 1.", Variables: []LessonVariable{{Name: "i", Value: "1"}}},
		{Line: 3, Caption: "Print the first value.", Variables: []LessonVariable{{Name: "i", Value: "1"}}, Output: "1"},
		{Line: 2, Caption: "The next loop has 2.", Variables: []LessonVariable{{Name: "i", Value: "2"}}, Output: "1"},
		{Line: 3, Caption: "Print the second value.", Variables: []LessonVariable{{Name: "i", Value: "2"}}, Output: "1\n2"},
		{Line: 2, Caption: "The final loop has 3.", Variables: []LessonVariable{{Name: "i", Value: "3"}}, Output: "1\n2"},
		{Line: 3, Caption: "Print the final value.", Variables: []LessonVariable{{Name: "i", Value: "3"}}, Output: "1\n2\n3"},
		{Line: 2, Caption: "The condition stops at 4.", Variables: []LessonVariable{{Name: "i", Value: "4"}}, Output: "1\n2\n3"},
	})}},
	{ID: "c31.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch03Question("c31.predict-loop", "Predict the output:\n\n~~~go\nfunc main() {\n\tfor n := 0; n < 2; n++ {\n\t\tfmt.Println(\"step\")\n\t}\n}\n~~~", []string{"one step", "two steps", "three steps"}, 1, []string{
			"The condition allows two loops before stopping.",
			"The counter is 0 then 1, so it prints twice.",
			"When n becomes 2, the condition is false.",
		}),
		ch03Walkthrough("c31.loop-reveal", "Follow the counter through each loop.", `func main() {
	for n := 0; n < 2; n++ {
		fmt.Println("step")
	}
}`, "e31-loop-reveal", []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The counter starts at zero.", Variables: []LessonVariable{{Name: "n", Value: "0"}}},
			{Line: 3, Caption: "Print the first time.", Variables: []LessonVariable{{Name: "n", Value: "0"}}, Output: "step"},
			{Line: 2, Caption: "The counter is now 1.", Variables: []LessonVariable{{Name: "n", Value: "1"}}, Output: "step"},
			{Line: 3, Caption: "Print the second time.", Variables: []LessonVariable{{Name: "n", Value: "1"}}, Output: "step\nstep"},
			{Line: 2, Caption: "The condition stops at 2.", Variables: []LessonVariable{{Name: "n", Value: "2"}}, Output: "step\nstep"},
		}),
		ch03Sandbox("ch03-loop-counter", "Step through the counter and output.", "ch03-repeat-counter", ch03RepeatConfig, true),
	},
	},
	{ID: "c31.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch03Markdown("c31.watch", "- Put the commands to repeat inside ~{ }~.\n- A wrong condition can repeat too many or too few times.")}},
	{ID: "c31.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch03Markdown("c31.plan", "Review each loop: does the counter change, and what does the drone do?\n\nBefore the last move, check the drone's position so you do not request an extra step."),
		ch03Question("c31.loop-count", "How many times does ~i := 0; i < 3; i++~ run?", []string{"Twice", "3 times", "4 times"}, 1, []string{
			"The values 0 and 1 are below 3; check value 2 too.",
			"Values 0, 1, and 2 satisfy the condition.",
			"At 3 the condition is false, so the body stops.",
		}),
	}},
	{ID: "c31.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch03Markdown("c31.turn", "Edit the starter into one clear loop. Run it, then add repetitions and watch the drone position and inventory.")}},
}

var ch03FullFieldArabic = []LessonStep{
	{ID: "c32.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch03Markdown("c32.mission", "1. احصد الطحالب في الحقل كله: 6 صفوف و6 أعمدة.\n2. الـstarter بيحصد الخانة اللي الدرون واقف عليها بس.\n3. هدف الذهب 71 ticks؛ راقب الحركة عند نهاية كل صف.")}},
	{ID: "c32.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch03Walkthrough("c32.nested-trace", "loop جوه loop بتقسّم الشغل لمستويين.\n\nالـloop الخارجية بتختار صف، والداخلية بتعدّي على خاناته. هنا شبكة تدريب صغيرة.", `func main() {
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			fmt.Println(row, col)
		}
	}
}`, "e32-nested-loops", []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "ابدأ بالصف صفر.", Variables: []LessonVariable{{Name: "row", Value: "0"}}},
		{Line: 3, Caption: "ابدأ بالعمود صفر.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}},
		{Line: 4, Caption: "اطبع الزوج الأول.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}, Output: "0 0"},
		{Line: 3, Caption: "العمود التالي هو 1.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "1"}}, Output: "0 0"},
		{Line: 4, Caption: "اطبع الزوج التالي.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "1"}}, Output: "0 0\n0 1"},
		{Line: 2, Caption: "ابدأ الصف التالي.", Variables: []LessonVariable{{Name: "row", Value: "1"}}, Output: "0 0\n0 1"},
		{Line: 3, Caption: "العمود يبدأ من صفر.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "0 0\n0 1"},
		{Line: 4, Caption: "اطبع زوج الصف الجديد.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "0 0\n0 1\n1 0"},
		{Line: 3, Caption: "ننتقل للعمود 1.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "1"}}, Output: "0 0\n0 1\n1 0"},
		{Line: 4, Caption: "اطبع آخر زوج.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "1"}}, Output: "0 0\n0 1\n1 0\n1 1"},
	})}},
	{ID: "c32.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch03Question("c32.predict-nested", "توقّع السطور المطبوعة:\n\n~~~go\nfunc main() {\n\tfor row := 0; row < 2; row++ {\n\t\tfor col := 0; col < 1; col++ {\n\t\t\tfmt.Println(\"row\", row)\n\t\t}\n\t}\n}\n~~~", []string{"row 0 مرة واحدة", "row 0 ثم row 1", "row 0 وrow 1 مرتين"}, 1, []string{
			"الـloop الخارجية بتلف مرتين حتى لو الداخلية لفة واحدة.",
			"كل صف بيعدّي على العمود الوحيد مرة واحدة.",
			"الـloop الداخلية مش بتكرر العنصر لما حدها يبقى واحد.",
		}),
		ch03Walkthrough("c32.nested-reveal", "تابع كل صف واللفة الداخلية.", `func main() {
	for row := 0; row < 2; row++ {
		for col := 0; col < 1; col++ {
			fmt.Println("row", row)
		}
	}
}`, "e32-nested-reveal", []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "الصف الأول رقمه صفر.", Variables: []LessonVariable{{Name: "row", Value: "0"}}},
			{Line: 3, Caption: "الداخلية تبدأ من صفر.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}},
			{Line: 4, Caption: "اطبع الصف صفر.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}, Output: "row 0"},
			{Line: 2, Caption: "الصف التالي رقمه 1.", Variables: []LessonVariable{{Name: "row", Value: "1"}}, Output: "row 0"},
			{Line: 3, Caption: "ابدأ العمود من جديد.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "row 0"},
			{Line: 4, Caption: "اطبع الصف الأول.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "row 0\nrow 1"},
		}),
		ch03Sandbox("ch03-nested-grid", "اضغط وشوف العداد الداخلي يرجع.", "ch03-nested-grid", ch03NestedConfig, false),
	},
	},
	{ID: "c32.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch03Markdown("c32.watch", "- العداد الداخلي بيرجع للبداية مع كل صف جديد.\n- الشرط الحدّي غلط ممكن يفوّت صف أو يطلب حركة زيادة.")}},
	{ID: "c32.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch03Markdown("c32.plan", "قسّم التفكير: صف واحد، ثم الانتقال للصف اللي بعده.\n\nراجع مكان البداية والنهاية قبل ما تختار اتجاه الحركة في الصف التالي."),
		ch03Question("c32.inner-count", "لو الخارجية بتلف 3 مرات والداخلية مرتين، كام مرة يشتغل جسم الداخلية؟", []string{"5", "6", "8"}, 1, []string{
			"الـloop الداخلية بتكرر لفتاتها بدل جمع الأعداد.",
			"كل لفة من الـloop الخارجية بتشغّل لفتين داخليتين.",
			"الضرب هنا 3 في 2، مش 4 في 2.",
		}),
	}},
	{ID: "c32.turn-step", Title: "دورك", Blocks: []LessonBlock{ch03Markdown("c32.turn", "ابدأ بصف واحد واتأكد إن كل خانة اتحصدت مرة. بعدها ضيف الصفوف، وشغّل الكود وراجع الحركة عند الأطراف.")}},
}

var ch03FullFieldEnglish = []LessonStep{
	{ID: "c32.mission-step", Title: "Mission", Blocks: []LessonBlock{ch03Markdown("c32.mission", "1. Harvest the whole field: 6 rows and 6 columns.\n2. The starter harvests only the tile under the drone.\n3. Gold is 71 ticks; watch movement at each row end.")}},
	{ID: "c32.idea-step", Title: "Idea", Blocks: []LessonBlock{ch03Walkthrough("c32.nested-trace", "A loop inside a loop splits work into two levels.\n\nThe outer loop picks a row; the inner visits its items. Values are illustrative.", `func main() {
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			fmt.Println(row, col)
		}
	}
}`, "e32-nested-loops", []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Row zero.", Variables: []LessonVariable{{Name: "row", Value: "0"}}},
		{Line: 3, Caption: "Column zero.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}},
		{Line: 4, Caption: "Print first pair.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}, Output: "0 0"},
		{Line: 3, Caption: "Column one.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "1"}}, Output: "0 0"},
		{Line: 4, Caption: "Print next pair.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "1"}}, Output: "0 0\n0 1"},
		{Line: 2, Caption: "Next row.", Variables: []LessonVariable{{Name: "row", Value: "1"}}, Output: "0 0\n0 1"},
		{Line: 3, Caption: "Restart column.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "0 0\n0 1"},
		{Line: 4, Caption: "Print row pair.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "0 0\n0 1\n1 0"},
		{Line: 3, Caption: "Column one.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "1"}}, Output: "0 0\n0 1\n1 0"},
		{Line: 4, Caption: "Print final pair.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "1"}}, Output: "0 0\n0 1\n1 0\n1 1"},
	})}},
	{ID: "c32.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch03Question("c32.predict-nested", "Predict the printed lines:\n\n~~~go\nfunc main() {\n\tfor row := 0; row < 2; row++ {\n\t\tfor col := 0; col < 1; col++ {\n\t\t\tfmt.Println(\"row\", row)\n\t\t}\n\t}\n}\n~~~", []string{"row 0 once", "row 0 then row 1", "row 0 and row 1 twice"}, 1, []string{
			"The outer loop has two turns, even with one inner turn.",
			"Each row visits its only column once.",
			"A limit of one means one inner turn.",
		}),
		ch03Walkthrough("c32.nested-reveal", "Follow each row and inner loop.", `func main() {
	for row := 0; row < 2; row++ {
		for col := 0; col < 1; col++ {
			fmt.Println("row", row)
		}
	}
}`, "e32-nested-reveal", []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The first row is zero.", Variables: []LessonVariable{{Name: "row", Value: "0"}}},
			{Line: 3, Caption: "The inner loop starts at zero.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}},
			{Line: 4, Caption: "Print row zero.", Variables: []LessonVariable{{Name: "row", Value: "0"}, {Name: "col", Value: "0"}}, Output: "row 0"},
			{Line: 2, Caption: "The next row is 1.", Variables: []LessonVariable{{Name: "row", Value: "1"}}, Output: "row 0"},
			{Line: 3, Caption: "Restart the column at zero.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "row 0"},
			{Line: 4, Caption: "Print row one.", Variables: []LessonVariable{{Name: "row", Value: "1"}, {Name: "col", Value: "0"}}, Output: "row 0\nrow 1"},
		}),
		ch03Sandbox("ch03-nested-grid", "Step through the inner counter resetting.", "ch03-nested-grid", ch03NestedConfig, true),
	},
	},
	{ID: "c32.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch03Markdown("c32.watch", "- The inner counter restarts for each new row.\n- A bad bound may skip a row or add a move.")}},
	{ID: "c32.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch03Markdown("c32.plan", "Plan one row, then the next.\n\nCheck row endpoints before picking a travel direction."),
		ch03Question("c32.inner-count", "If the outer loop has 3 turns and the inner has 2, how often does its body run?", []string{"5", "6", "8"}, 1, []string{
			"The inner loop repeats for every row.",
			"Each row runs the inner body twice.",
			"The product is 3 times 2.",
		}),
	}},
	{ID: "c32.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch03Markdown("c32.turn", "Start with one row and confirm each tile is harvested once. Add rows next, then run it and review movement at the edges.")}},
}

var ch03EndlessArabic = []LessonStep{
	{ID: "c33.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch03Markdown("c33.mission", "1. اجمع 30 طحلب من حقل 3×3.\n2. الـstarter فيه loop بشرط، بس جسمها فاضي.\n3. هدف الذهب 140 ticks؛ تابع العدد والوقت مع بعض.")}},
	{ID: "c33.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch03Walkthrough("c33.condition-loop", "ممكن تكتب for loop بشرط بس، زي while في لغات تانية.\n\nلازم حاجة جوه الجسم تغيّر الشرط، وإلا الـloop هتفضل شغالة.", `func main() {
	count := 0
	for count < 3 {
		count++
		fmt.Println(count)
	}
}`, "e33-condition-loop", []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "العدد يبدأ بصفر.", Variables: []LessonVariable{{Name: "count", Value: "0"}}},
		{Line: 3, Caption: "الشرط بيسمح باللفة.", Variables: []LessonVariable{{Name: "count", Value: "0"}}},
		{Line: 4, Caption: "زوّد العدد لـ1.", Variables: []LessonVariable{{Name: "count", Value: "1"}}},
		{Line: 5, Caption: "اطبع العدد الحالي.", Variables: []LessonVariable{{Name: "count", Value: "1"}}, Output: "1"},
		{Line: 3, Caption: "العدد 1 بيحقق الشرط.", Variables: []LessonVariable{{Name: "count", Value: "1"}}, Output: "1"},
		{Line: 4, Caption: "زوّد العدد لـ2.", Variables: []LessonVariable{{Name: "count", Value: "2"}}, Output: "1"},
		{Line: 5, Caption: "اطبع العدد التالي.", Variables: []LessonVariable{{Name: "count", Value: "2"}}, Output: "1\n2"},
		{Line: 3, Caption: "العدد 2 بيحقق الشرط.", Variables: []LessonVariable{{Name: "count", Value: "2"}}, Output: "1\n2"},
		{Line: 4, Caption: "زوّد العدد لـ3.", Variables: []LessonVariable{{Name: "count", Value: "3"}}, Output: "1\n2"},
		{Line: 5, Caption: "اطبع آخر عدد.", Variables: []LessonVariable{{Name: "count", Value: "3"}}, Output: "1\n2\n3"},
		{Line: 3, Caption: "الشرط وقف عند 3.", Variables: []LessonVariable{{Name: "count", Value: "3"}}, Output: "1\n2\n3"},
	})}},
	{ID: "c33.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch03Question("c33.predict-condition", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tleft := 2\n\tfor left > 0 {\n\t\tfmt.Println(left)\n\t\tleft--\n\t}\n}\n~~~", []string{"2 بس", "2 وبعدها 1", "2 وبعدها 1 وبعدها 0"}, 1, []string{
			"بعد الطباعة، left بتقل قبل ما نفحص الشرط تاني.",
			"الشرط true عند 2 و1، وبيقف عند صفر.",
			"لما توصل لصفر، جسم الـloop مش بيتنفّذ.",
		}),
		ch03Walkthrough("c33.condition-reveal", "شوف الشرط بعد كل تحديث.", `func main() {
	left := 2
	for left > 0 {
		fmt.Println(left)
		left--
	}
}`, "e33-condition-reveal", []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "left تبدأ عند 2.", Variables: []LessonVariable{{Name: "left", Value: "2"}}},
			{Line: 3, Caption: "الشرط بيسمح باللفة.", Variables: []LessonVariable{{Name: "left", Value: "2"}}},
			{Line: 4, Caption: "اطبع الـvalue 2.", Variables: []LessonVariable{{Name: "left", Value: "2"}}, Output: "2"},
			{Line: 5, Caption: "left بقت 1.", Variables: []LessonVariable{{Name: "left", Value: "1"}}, Output: "2"},
			{Line: 4, Caption: "اطبع الـvalue 1.", Variables: []LessonVariable{{Name: "left", Value: "1"}}, Output: "2\n1"},
			{Line: 5, Caption: "left بقت صفر.", Variables: []LessonVariable{{Name: "left", Value: "0"}}, Output: "2\n1"},
			{Line: 3, Caption: "الشرط وقف عند صفر.", Variables: []LessonVariable{{Name: "left", Value: "0"}}, Output: "2\n1"},
		}),
		ch03Sandbox("ch03-condition-loop", "اضغط وتابع تغيّر الشرط.", "ch03-condition-loop", ch03ConditionConfig, false),
	},
	},
	{ID: "c33.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch03Markdown("c33.watch", "- لو الشرط ما بيتغيرش، loop مش هتخلص.\n- الطحلب محتاج 8 ticks بعد الزرع؛ بدري عليه يتدمر.")}},
	{ID: "c33.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch03Markdown("c33.plan", "اتأكد إن كل دورة بتزود التقدم وبتدي الوقت يعدّي.\n\nفي إعداد المهمة ده، ~colony.Count~ قراءة للمخزون؛ أوامر الدرون بتحرك الوقت."),
		ch03Question("c33.loop-progress", "إيه اللي لازم يحصل عشان شرط الـloop يقف؟", []string{"الشرط يفضل true", "قيمة الشرط تتغير", "نضيف تعليق"}, 1, []string{
			"شرط true يخلي الجسم يدخل لفة تانية.",
			"تغيير الـvalue ممكن يخلي الشرط false.",
			"التعليق ما بيغيّرش حالة البرنامج.",
		}),
	}},
	{ID: "c33.turn-step", Title: "دورك", Blocks: []LessonBlock{ch03Markdown("c33.turn", "ابدأ بتحديث واحد واضح في جسم الـloop. شغّل الكود وراقب عدّاد الطحالب؛ بعدها فكّر إزاي الوقت يعدّي لحد النضج.")}},
}

var ch03EndlessEnglish = []LessonStep{
	{ID: "c33.mission-step", Title: "Mission", Blocks: []LessonBlock{ch03Markdown("c33.mission", "1. Collect 30 moss from a 3×3 field.\n2. The starter has a conditional loop, but its body is empty.\n3. Gold is 140 ticks; track both the count and time.")}},
	{ID: "c33.idea-step", Title: "Idea", Blocks: []LessonBlock{ch03Walkthrough("c33.condition-loop", "A condition-only ~for~ repeats until false.\n\nUpdate the condition in the body, or the loop keeps running.", `func main() {
	count := 0
	for count < 3 {
		count++
		fmt.Println(count)
	}
}`, "e33-condition-loop", []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Count starts at zero.", Variables: []LessonVariable{{Name: "count", Value: "0"}}},
		{Line: 3, Caption: "Condition allows looping.", Variables: []LessonVariable{{Name: "count", Value: "0"}}},
		{Line: 4, Caption: "Increase count to 1.", Variables: []LessonVariable{{Name: "count", Value: "1"}}},
		{Line: 5, Caption: "Print the count.", Variables: []LessonVariable{{Name: "count", Value: "1"}}, Output: "1"},
		{Line: 3, Caption: "Count 1 still qualifies.", Variables: []LessonVariable{{Name: "count", Value: "1"}}, Output: "1"},
		{Line: 4, Caption: "Increase count to 2.", Variables: []LessonVariable{{Name: "count", Value: "2"}}, Output: "1"},
		{Line: 5, Caption: "Print the next count.", Variables: []LessonVariable{{Name: "count", Value: "2"}}, Output: "1\n2"},
		{Line: 3, Caption: "Count 2 still qualifies.", Variables: []LessonVariable{{Name: "count", Value: "2"}}, Output: "1\n2"},
		{Line: 4, Caption: "Increase count to 3.", Variables: []LessonVariable{{Name: "count", Value: "3"}}, Output: "1\n2"},
		{Line: 5, Caption: "Print the final count.", Variables: []LessonVariable{{Name: "count", Value: "3"}}, Output: "1\n2\n3"},
		{Line: 3, Caption: "Condition stops at 3.", Variables: []LessonVariable{{Name: "count", Value: "3"}}, Output: "1\n2\n3"},
	})}},
	{ID: "c33.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch03Question("c33.predict-condition", "Predict the output:\n\n~~~go\nfunc main() {\n\tleft := 2\n\tfor left > 0 {\n\t\tfmt.Println(left)\n\t\tleft--\n\t}\n}\n~~~", []string{"2 only", "2 then 1", "2 then 1 then 0"}, 1, []string{
			"After printing, left decreases before the next check.",
			"The condition is true at 2 and 1, then stops at zero.",
			"At zero, the loop body does not run.",
		}),
		ch03Walkthrough("c33.condition-reveal", "Watch the condition after each update.", `func main() {
	left := 2
	for left > 0 {
		fmt.Println(left)
		left--
	}
}`, "e33-condition-reveal", []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "left starts at 2.", Variables: []LessonVariable{{Name: "left", Value: "2"}}},
			{Line: 3, Caption: "Condition allows a loop.", Variables: []LessonVariable{{Name: "left", Value: "2"}}},
			{Line: 4, Caption: "Print the value 2.", Variables: []LessonVariable{{Name: "left", Value: "2"}}, Output: "2"},
			{Line: 5, Caption: "left becomes 1.", Variables: []LessonVariable{{Name: "left", Value: "1"}}, Output: "2"},
			{Line: 4, Caption: "Print the value 1.", Variables: []LessonVariable{{Name: "left", Value: "1"}}, Output: "2\n1"},
			{Line: 5, Caption: "left becomes zero.", Variables: []LessonVariable{{Name: "left", Value: "0"}}, Output: "2\n1"},
			{Line: 3, Caption: "The condition stops at zero.", Variables: []LessonVariable{{Name: "left", Value: "0"}}, Output: "2\n1"},
		}),
		ch03Sandbox("ch03-condition-loop", "Step through the changing condition.", "ch03-condition-loop", ch03ConditionConfig, true),
	},
	},
	{ID: "c33.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch03Markdown("c33.watch", "- If the condition never changes, the loop will not finish.\n- Moss needs 8 ticks after planting; early harvest destroys it.")}},
	{ID: "c33.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch03Markdown("c33.plan", "Make sure each cycle adds progress and lets time pass.\n\nIn this mission setup, ~colony.Count~ reads inventory; drone commands advance time."),
		ch03Question("c33.loop-progress", "What must happen for a loop condition to stop?", []string{"It stays true", "Its value changes", "Add a comment"}, 1, []string{
			"A true condition sends the body around again.",
			"A value change may make the condition false.",
			"A comment does not change program state.",
		}),
	}},
	{ID: "c33.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch03Markdown("c33.turn", "Start with one clear update in the loop body. Run it and watch the moss count; then consider how time passes until ripening.")}},
}
