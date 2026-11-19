package levels

func ch04Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch04Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch04Walkthrough(id, text, code, exampleID string, frames []LessonFrame) LessonBlock {
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

func ch04Locale(ar, en string) map[string]any {
	return map[string]any{"ar": ar, "en": en}
}

func ch04Sandbox(id, text, kind string, config map[string]any, english bool) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: ch04ConfigForLanguage(config, english).(map[string]any)}}
}

func ch04ConfigForLanguage(value any, english bool) any {
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
			localized[key] = ch04ConfigForLanguage(nested, english)
		}
		return localized
	case []any:
		localized := make([]any, len(current))
		for i, nested := range current {
			localized[i] = ch04ConfigForLanguage(nested, english)
		}
		return localized
	default:
		return current
	}
}

var ch04IfConfig = map[string]any{
	"labels": map[string]any{
		"condition": ch04Locale("الشرط", "Condition"),
		"body":      ch04Locale("الجسم", "Body"),
		"choose":    ch04Locale("اختار حالة", "Choose a state"),
		"notReady":  ch04Locale("false", "false"),
		"ready":     ch04Locale("true", "true"),
		"reset":     ch04Locale("من الأول", "Reset"),
	},
	"choices": []any{
		map[string]any{"ready": false, "label": ch04Locale("لسه", "Not yet"), "bodyState": ch04Locale("اتخطّى", "skipped"), "caption": ch04Locale("false تعدّي الجسم.", "False skips the body.")},
		map[string]any{"ready": true, "label": ch04Locale("جاهز", "Ready"), "bodyState": ch04Locale("اشتغل", "runs"), "caption": ch04Locale("true تدخل الجسم.", "True enters the body.")},
	},
}

var ch04SwitchConfig = map[string]any{
	"labels": map[string]any{
		"signal": ch04Locale("الإشارة", "Signal"),
		"case":   ch04Locale("الحالة المختارة", "Selected case"),
		"result": ch04Locale("الناتج", "Result"),
		"choose": ch04Locale("اختار قيمة", "Choose a value"),
	},
	"choices": []any{
		map[string]any{"label": ch04Locale("أزرق", "Blue"), "value": ch04Locale("blue", "blue"), "case": ch04Locale("case blue", "case blue"), "result": ch04Locale("move", "move"), "caption": ch04Locale("القيمة تطابق case واحدة.", "The value matches one case.")},
		map[string]any{"label": ch04Locale("كهرماني", "Amber"), "value": ch04Locale("amber", "amber"), "case": ch04Locale("case amber", "case amber"), "result": ch04Locale("pause", "pause"), "caption": ch04Locale("case واحدة بس بتشتغل.", "Only one case runs.")},
		map[string]any{"label": ch04Locale("قيمة تانية", "Other value"), "value": ch04Locale("other", "other"), "case": ch04Locale("default", "default"), "result": ch04Locale("check", "check"), "caption": ch04Locale("default تستقبل الباقي.", "Default handles the rest.")},
	},
}

var ch04OrderedConfig = map[string]any{
	"labels": map[string]any{
		"hours":       ch04Locale("ساعات التوصيل", "Delivery hours"),
		"firstCheck":  ch04Locale("hours < 3", "hours < 3"),
		"secondCheck": ch04Locale("hours < 8", "hours < 8"),
		"true":        ch04Locale("true", "true"),
		"false":       ch04Locale("false", "false"),
		"skipped":     ch04Locale("متخطّي", "skipped"),
		"result":      ch04Locale("التصنيف", "Label"),
		"next":        ch04Locale("العينة الجاية", "Next sample"),
		"reset":       ch04Locale("من الأول", "Reset"),
	},
	"frames": []any{
		map[string]any{"hours": 2, "first": true, "second": nil, "result": ch04Locale("short", "short"), "caption": ch04Locale("الأول true؛ التاني بيتخطّى.", "First is true; second is skipped.")},
		map[string]any{"hours": 5, "first": false, "second": true, "result": ch04Locale("medium", "medium"), "caption": ch04Locale("الأول false، التاني true.", "First is false; second is true.")},
		map[string]any{"hours": 9, "first": false, "second": false, "result": ch04Locale("long", "long"), "caption": ch04Locale("الاتنين false؛ تفضل القيمة الأصلية.", "Both are false; the original stays.")},
	},
}

var ch04RipeOnlyArabic = []LessonStep{
	{ID: "c41.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch04Markdown("c41.mission", "1. اجمع 10 طحالب من غير ما تدمّر نبتة.\n2. الـstarter بيفحص الجاهزية، وجسم الشرط لسه فاضي.\n3. هدف الذهب 55 ticks؛ خليك حريص مع كل حصاد.")}},
	{ID: "c41.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch04Walkthrough("c41.if-trace", "~if~ بتشغّل جسمها لما الشرط يبقى true.\n\nالقرار هنا على value تدريب، مش فحص المهمة.", `func main() {
	ready := false
	if ready {
		fmt.Println("harvest")
	} else {
		fmt.Println("wait")
	}
}`, "e41-if-branch", []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "الجاهزية false.", Variables: []LessonVariable{{Name: "ready", Value: "false"}}},
		{Line: 3, Caption: "اختبر الـvalue بتاعة الشرط."},
		{Line: 5, Caption: "اختار فرع الانتظار.", Variables: []LessonVariable{{Name: "ready", Value: "false"}}},
		{Line: 6, Caption: "اطبع الانتظار.", Variables: []LessonVariable{{Name: "ready", Value: "false"}}, Output: "wait"},
	})}},
	{ID: "c41.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch04Question("c41.predict-if", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tready := true\n\tif ready {\n\t\tfmt.Println(\"ready\")\n\t}\n}\n~~~", []string{"مفيش output", "ready", "false"}, 1, []string{
			"الجسم بيتنفّذ لأن الشرط true.",
			"الـvalue true بتدخل جسم if وبتطبع الكلمة.",
			"الكود بيطبع string، مش الـvalue بتاعة الشرط.",
		}),
		ch04Walkthrough("c41.if-reveal", "شوف الفرع اللي بيتنفّذ.", `func main() {
	ready := true
	if ready {
		fmt.Println("ready")
	}
}`, "e41-if-reveal", []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "الجاهزية true.", Variables: []LessonVariable{{Name: "ready", Value: "true"}}},
			{Line: 3, Caption: "الشرط يسمح بالتنفيذ.", Variables: []LessonVariable{{Name: "ready", Value: "true"}}},
			{Line: 4, Caption: "اطبع حالة الجاهزية.", Variables: []LessonVariable{{Name: "ready", Value: "true"}}, Output: "ready"},
		}),
		ch04Sandbox("ch04-if-playground", "اختار حالة وشوف فرع if.", "ch04-if-branch", ch04IfConfig, false),
	},
	},
	{ID: "c41.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch04Markdown("c41.watch", "- في إعداد المهمة الحالي، مسح ~CanHarvest~ بياخد tick؛ الحصاد قبل النضج يدمّر النبات.\n- حط الفعل الخطر جوه شرط الحماية.")}},
	{ID: "c41.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch04Markdown("c41.plan", "افصل القرار عن الحركة: افحص الخانة، وبعدها قرر تعمل إيه.\n\nراجع كل فرع عشان تعرف إمتى بيتنفّذ وإمتى بيتخطّى."),
		ch04Question("c41.safe-guard", "لو ~ready~ false، أنهي فرع من ~if ready { ... }~ يشتغل؟", []string{"الجسم بين الأقواس", "ولا حاجة", "كل الفروع"}, 1, []string{
			"الجسم محتاج شرط true عشان يشتغل.",
			"من غير else، الشرط false بيعدّي من غير جسم.",
			"if ما بتشغّلاش كل الفروع مع بعض.",
		}),
	}},
	{ID: "c41.turn-step", Title: "دورك", Blocks: []LessonBlock{ch04Markdown("c41.turn", "كمّل جسم الـif في الـstarter. شغّله وراقب الجاهزية قبل وبعد ما الدرون يلف على الخانات.")}},
}

var ch04RipeOnlyEnglish = []LessonStep{
	{ID: "c41.mission-step", Title: "Mission", Blocks: []LessonBlock{ch04Markdown("c41.mission", "1. Collect 10 moss without destroying a plant.\n2. The starter checks readiness, but the guarded body is empty.\n3. Gold is 55 ticks; be careful with every harvest.")}},
	{ID: "c41.idea-step", Title: "Idea", Blocks: []LessonBlock{ch04Walkthrough("c41.if-trace", "~if~ runs its body when a condition is true.\n\nThis decision uses a practice value, not the mission scan.", `func main() {
	ready := false
	if ready {
		fmt.Println("harvest")
	} else {
		fmt.Println("wait")
	}
}`, "e41-if-branch", []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Readiness is false.", Variables: []LessonVariable{{Name: "ready", Value: "false"}}},
		{Line: 3, Caption: "Check the condition."},
		{Line: 5, Caption: "Choose the wait branch.", Variables: []LessonVariable{{Name: "ready", Value: "false"}}},
		{Line: 6, Caption: "Print the wait message.", Variables: []LessonVariable{{Name: "ready", Value: "false"}}, Output: "wait"},
	})}},
	{ID: "c41.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch04Question("c41.predict-if", "Predict the output:\n\n~~~go\nfunc main() {\n\tready := true\n\tif ready {\n\t\tfmt.Println(\"ready\")\n\t}\n}\n~~~", []string{"No output", "ready", "false"}, 1, []string{
			"The body runs because the condition is true.",
			"The true value enters if and prints the word.",
			"The code prints text, not the condition value.",
		}),
		ch04Walkthrough("c41.if-reveal", "Watch which branch runs.", `func main() {
	ready := true
	if ready {
		fmt.Println("ready")
	}
}`, "e41-if-reveal", []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Readiness is true.", Variables: []LessonVariable{{Name: "ready", Value: "true"}}},
			{Line: 3, Caption: "The condition allows execution.", Variables: []LessonVariable{{Name: "ready", Value: "true"}}},
			{Line: 4, Caption: "Print the ready state.", Variables: []LessonVariable{{Name: "ready", Value: "true"}}, Output: "ready"},
		}),
		ch04Sandbox("ch04-if-playground", "Choose a state and watch the if branch.", "ch04-if-branch", ch04IfConfig, true),
	},
	},
	{ID: "c41.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch04Markdown("c41.watch", "- In this mission's default clock, the ~CanHarvest~ scan costs a tick; early harvest destroys the plant.\n- Put the risky action inside its guard.")}},
	{ID: "c41.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch04Markdown("c41.plan", "Separate the decision from movement: check a tile, then choose an action.\n\nReview each branch to see when it runs and when it is skipped."),
		ch04Question("c41.safe-guard", "If ~ready~ is false, which part of ~if ready { ... }~ runs?", []string{"The body in braces", "Nothing", "Every branch"}, 1, []string{
			"The body needs a true condition to run.",
			"Without else, a false condition skips the body.",
			"An if does not run every branch at once.",
		}),
	}},
	{ID: "c41.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch04Markdown("c41.turn", "Fill in the if body in the starter. Run it and watch readiness as the drone visits tiles.")}},
}

var ch04TerrainArabic = []LessonStep{
	{ID: "c42.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch04Markdown("c42.mission", "1. ازرع كل خانة بنوع النبات اللي يناسب أرضها.\n2. الـstarter عامل ~switch~ لحالتين، ولسه من غير أفعال.\n3. هدف الذهب 60 ticks؛ في خانات محتاجة تعامل مختلف.")}},
	{ID: "c42.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch04Walkthrough("c42.switch-trace", "لما value واحدة يبقى ليها اختيارات كتير، ~switch~ بتوضّح الفروع.\n\nGo بتنفّذ case واحدة؛ مش محتاج ~break~.", `func main() {
	signal := 2
	switch signal {
	case 1:
		fmt.Println("hold")
	case 2:
		fmt.Println("go")
	default:
		fmt.Println("check")
	}
}`, "e42-switch-case", []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "value الإشارة 2.", Variables: []LessonVariable{{Name: "signal", Value: "2"}}},
		{Line: 3, Caption: "قارن الـvalue بالحالات."},
		{Line: 6, Caption: "الحالة المطابقة هي 2.", Variables: []LessonVariable{{Name: "signal", Value: "2"}}},
		{Line: 7, Caption: "اطبع أمر الحركة.", Variables: []LessonVariable{{Name: "signal", Value: "2"}}, Output: "go"},
	})}},
	{ID: "c42.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch04Question("c42.predict-switch", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tcolor := \"amber\"\n\tswitch color {\n\tcase \"blue\":\n\t\tfmt.Println(\"move\")\n\tcase \"amber\":\n\t\tfmt.Println(\"pause\")\n\tdefault:\n\t\tfmt.Println(\"check\")\n\t}\n}\n~~~", []string{"move", "pause", "check"}, 1, []string{
			"الـvalue amber مش هي blue.",
			"الـcase دي مطابقة للـvalue اللي في color.",
			"default بتشتغل لما مفيش case مطابقة.",
		}),
		ch04Walkthrough("c42.switch-reveal", "تابع اختيار الحالة المطابقة.", `func main() {
	color := "amber"
	switch color {
	case "blue":
		fmt.Println("move")
	case "amber":
		fmt.Println("pause")
	default:
		fmt.Println("check")
	}
}`, "e42-switch-reveal", []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "color الـvalue بتاعتها amber.", Variables: []LessonVariable{{Name: "color", Value: "amber"}}},
			{Line: 3, Caption: "افحص الحالات بالترتيب."},
			{Line: 6, Caption: "وصلنا للحالة المطابقة.", Variables: []LessonVariable{{Name: "color", Value: "amber"}}},
			{Line: 7, Caption: "اطبع pause.", Variables: []LessonVariable{{Name: "color", Value: "amber"}}, Output: "pause"},
		}),
		ch04Sandbox("ch04-switch-playground", "اختار value وشوف case المطابقة.", "ch04-switch-case", ch04SwitchConfig, false),
	},
	},
	{ID: "c42.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch04Markdown("c42.watch", "- في إعداد المهمة الحالي، مسح ~GetGround~ بياخد tick. خد الـvalue مرة وإنت عارف هتستخدمها فين.\n- كل ~case~ بتنتهي لوحدها؛ مفيش fallthrough لوحدها.")}},
	{ID: "c42.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch04Markdown("c42.plan", "لف على الخانات بنفس ترتيب traversal اللي اتعلمته.\n\nلكل خانة، اقرا الأرض واختار الفرع المناسب؛ سيب الحالة غير المدعومة واضحة."),
		ch04Question("c42.switch-default", "إمتى ~default~ بتشتغل؟", []string{"دايماً بعد case", "لو مفيش case مطابقة", "قبل أول case"}, 1, []string{
			"Go ما بتكملش لوحدها من case للتانية.",
			"default هو الفرع الاحتياطي لما مفيش تطابق.",
			"الاختيار بيتحدد حسب الـvalue اللي في switch.",
		}),
	}},
	{ID: "c42.turn-step", Title: "دورك", Blocks: []LessonBlock{ch04Markdown("c42.turn", "كمّل فروع الـswitch في الـstarter. جرّب التشغيل وراجع نوع الأرض والفعل اللي حصل عند كل خانة.")}},
}

var ch04TerrainEnglish = []LessonStep{
	{ID: "c42.mission-step", Title: "Mission", Blocks: []LessonBlock{ch04Markdown("c42.mission", "1. Plant each tile with the crop that fits its ground.\n2. The starter switches on two cases, but has no actions yet.\n3. Gold is 60 ticks; some tiles need different handling.")}},
	{ID: "c42.idea-step", Title: "Idea", Blocks: []LessonBlock{ch04Walkthrough("c42.switch-trace", "When one value has several choices, ~switch~ makes branches clear.\n\nGo runs one case; no ~break~ is needed.", `func main() {
	signal := 2
	switch signal {
	case 1:
		fmt.Println("hold")
	case 2:
		fmt.Println("go")
	default:
		fmt.Println("check")
	}
}`, "e42-switch-case", []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "The signal value is 2.", Variables: []LessonVariable{{Name: "signal", Value: "2"}}},
		{Line: 3, Caption: "Compare the cases."},
		{Line: 6, Caption: "Case 2 is a match.", Variables: []LessonVariable{{Name: "signal", Value: "2"}}},
		{Line: 7, Caption: "Print the move command.", Variables: []LessonVariable{{Name: "signal", Value: "2"}}, Output: "go"},
	})}},
	{ID: "c42.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch04Question("c42.predict-switch", "Predict the output:\n\n~~~go\nfunc main() {\n\tcolor := \"amber\"\n\tswitch color {\n\tcase \"blue\":\n\t\tfmt.Println(\"move\")\n\tcase \"amber\":\n\t\tfmt.Println(\"pause\")\n\tdefault:\n\t\tfmt.Println(\"check\")\n\t}\n}\n~~~", []string{"move", "pause", "check"}, 1, []string{
			"The value amber is not blue.",
			"This case matches the color value.",
			"default runs when no case matches.",
		}),
		ch04Walkthrough("c42.switch-reveal", "Follow the matching case.", `func main() {
	color := "amber"
	switch color {
	case "blue":
		fmt.Println("move")
	case "amber":
		fmt.Println("pause")
	default:
		fmt.Println("check")
	}
}`, "e42-switch-reveal", []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "color is amber.", Variables: []LessonVariable{{Name: "color", Value: "amber"}}},
			{Line: 3, Caption: "Check the cases in order."},
			{Line: 6, Caption: "We found the matching case.", Variables: []LessonVariable{{Name: "color", Value: "amber"}}},
			{Line: 7, Caption: "Print pause.", Variables: []LessonVariable{{Name: "color", Value: "amber"}}, Output: "pause"},
		}),
		ch04Sandbox("ch04-switch-playground", "Choose a value and watch its case.", "ch04-switch-case", ch04SwitchConfig, true),
	},
	},
	{ID: "c42.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch04Markdown("c42.watch", "- In this mission's default clock, the ~GetGround~ scan costs a tick. Read its value once and plan its use.\n- Each ~case~ ends by itself; there is no automatic fallthrough.")}},
	{ID: "c42.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch04Markdown("c42.plan", "Visit the tiles using the traversal pattern you learned.\n\nFor each tile, read its ground and choose a branch; keep unsupported ground clear."),
		ch04Question("c42.switch-default", "When does ~default~ run?", []string{"Always after a case", "When no case matches", "Before the first case"}, 1, []string{
			"Go does not continue automatically from one case to another.",
			"default is the fallback branch when none match.",
			"The switch value controls the selection.",
		}),
	}},
	{ID: "c42.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch04Markdown("c42.turn", "Complete the switch branches in the starter. Run it and check the ground type and action at each tile.")}},
}

var ch04CheckpointArabic = []LessonStep{
	{ID: "c43.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch04Markdown("c43.mission", "1. صنّف قراءات الحساسات وابعت label لكل قراءة.\n2. الـstarter بيعرض أول رقم بس؛ البيانات والترتيب بيتغيروا.\n3. مفيش هدف ticks للذهب هنا؛ المهم العقد والحدود.")}},
	{ID: "c43.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch04Walkthrough("c43.ordered-bands", "سلسلة ~if / else if~ بتفحص الشروط بالترتيب.\n\nأول شرط true بيختار فرعه؛ المثال ده عن مدد توصيل للتوضيح.", `func main() {
	hours := 5
	label := "long"
	if hours < 3 {
		label = "short"
	} else if hours < 8 {
		label = "medium"
	}
	fmt.Println(label)
}`, "e43-ordered-bands", []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "المدة 5 ساعات.", Variables: []LessonVariable{{Name: "hours", Value: "5"}}},
		{Line: 3, Caption: "label بتبدأ بـlong.", Variables: []LessonVariable{{Name: "label", Value: "long"}}},
		{Line: 4, Caption: "5 مش أقل من 3."},
		{Line: 6, Caption: "5 أقل من 8.", Variables: []LessonVariable{{Name: "hours", Value: "5"}}},
		{Line: 7, Caption: "اختار تصنيف medium.", Variables: []LessonVariable{{Name: "label", Value: "medium"}}},
		{Line: 9, Caption: "اطبع التصنيف.", Variables: []LessonVariable{{Name: "label", Value: "medium"}}, Output: "medium"},
	})}},
	{ID: "c43.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch04Question("c43.predict-boundary", "عند ~level = 4~، إيه اللي يتطبع؟\n\n~~~go\nfunc main() {\n\tlevel := 4\n\tlabel := \"high\"\n\tif level < 4 {\n\t\tlabel = \"low\"\n\t} else {\n\t\tlabel = \"middle\"\n\t}\n\tfmt.Println(label)\n}\n~~~", []string{"low", "middle", "high"}, 1, []string{
			"4 مش أصغر من 4، فالفرع الأول بيتخطّى.",
			"الشرط false، فـelse تختار middle.",
			"الـlabel القديمة بتتغير في فرع else.",
		}),
		ch04Walkthrough("c43.boundary-reveal", "تابع الـvalue عند الحد.", `func main() {
	level := 4
	label := "high"
	if level < 4 {
		label = "low"
	} else {
		label = "middle"
	}
	fmt.Println(label)
}`, "e43-boundary-reveal", []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "الـvalue في level هي 4.", Variables: []LessonVariable{{Name: "level", Value: "4"}}},
			{Line: 4, Caption: "4 مش أقل من 4.", Variables: []LessonVariable{{Name: "level", Value: "4"}}},
			{Line: 6, Caption: "اختار فرع else.", Variables: []LessonVariable{{Name: "label", Value: "high"}}},
			{Line: 7, Caption: "حدّث label لـmiddle.", Variables: []LessonVariable{{Name: "label", Value: "middle"}}},
			{Line: 9, Caption: "اطبع التصنيف الجديد.", Variables: []LessonVariable{{Name: "label", Value: "middle"}}, Output: "middle"},
		}),
		ch04Sandbox("ch04-ordered-playground", "اضغط وتابع أول شرط مطابق.", "ch04-ordered-branches", ch04OrderedConfig, false),
	},
	},
	{ID: "c43.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch04Markdown("c43.watch", "- الحد لازم يتبع فئة واحدة بس. اختار ~<~ أو ~<=~ بوعي.\n- ~Reading(i)~ بترجع zero value للـindex الغلط؛ خليك جوه العدد.")}},
	{ID: "c43.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch04Markdown("c43.plan", "اتبع contract المهمة: قراءة واحدة، قرار واحد، وlabel واحدة بنفس الترتيب.\n\nافصل رفض الـvalues اللي برا المدى عن اختيار الفئة؛ راجع الحالة الحدّية."),
		ch04Question("c43.validate-first", "إيه الأفضل قبل ما تصنّف رقم؟", []string{"تتأكد إنه جوه المدى", "تطبع قيمته", "تختار أي فئة"}, 0, []string{
			"الفحص بيمنع رقم برا المدى يدخل فئة عادية.",
			"طباعة الرقم مش فحص للـcontract.",
			"الاختيار قبل الفحص ممكن يصنّف value مرفوضة.",
		}),
	}},
	{ID: "c43.turn-step", Title: "دورك", Blocks: []LessonBlock{ch04Markdown("c43.turn", "ابدأ بقراءة مفهرسة واحدة واكتب قرارها بنفسك. بعدين وسّع الشغل لكل القراءات، وشغّل الاختبارات المتاحة في المحرر.")}},
}

var ch04CheckpointEnglish = []LessonStep{
	{ID: "c43.mission-step", Title: "Mission", Blocks: []LessonBlock{ch04Markdown("c43.mission", "1. Classify sensor readings and report one label per reading.\n2. The starter prints one number; the data and order vary.\n3. There is no gold tick target here; the contract and boundaries matter.")}},
	{ID: "c43.idea-step", Title: "Idea", Blocks: []LessonBlock{ch04Walkthrough("c43.ordered-bands", "An ~if / else if~ chain checks conditions in order.\n\nThe first true condition selects its branch; this example uses delivery times.", `func main() {
	hours := 5
	label := "long"
	if hours < 3 {
		label = "short"
	} else if hours < 8 {
		label = "medium"
	}
	fmt.Println(label)
}`, "e43-ordered-bands", []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "The duration is 5 hours.", Variables: []LessonVariable{{Name: "hours", Value: "5"}}},
		{Line: 3, Caption: "label starts as long.", Variables: []LessonVariable{{Name: "label", Value: "long"}}},
		{Line: 4, Caption: "5 is not below 3."},
		{Line: 6, Caption: "5 is below 8.", Variables: []LessonVariable{{Name: "hours", Value: "5"}}},
		{Line: 7, Caption: "Choose the medium label.", Variables: []LessonVariable{{Name: "label", Value: "medium"}}},
		{Line: 9, Caption: "Print the category.", Variables: []LessonVariable{{Name: "label", Value: "medium"}}, Output: "medium"},
	})}},
	{ID: "c43.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch04Question("c43.predict-boundary", "At ~level = 4~, what prints?\n\n~~~go\nfunc main() {\n\tlevel := 4\n\tlabel := \"high\"\n\tif level < 4 {\n\t\tlabel = \"low\"\n\t} else {\n\t\tlabel = \"middle\"\n\t}\n\tfmt.Println(label)\n}\n~~~", []string{"low", "middle", "high"}, 1, []string{
			"4 is not below 4, so the first branch is skipped.",
			"The condition is false, so else chooses middle.",
			"The old label changes in the else branch.",
		}),
		ch04Walkthrough("c43.boundary-reveal", "Follow the boundary value.", `func main() {
	level := 4
	label := "high"
	if level < 4 {
		label = "low"
	} else {
		label = "middle"
	}
	fmt.Println(label)
}`, "e43-boundary-reveal", []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "level equals 4.", Variables: []LessonVariable{{Name: "level", Value: "4"}}},
			{Line: 4, Caption: "4 is not below 4.", Variables: []LessonVariable{{Name: "level", Value: "4"}}},
			{Line: 6, Caption: "Choose the else branch.", Variables: []LessonVariable{{Name: "label", Value: "high"}}},
			{Line: 7, Caption: "Update label to middle.", Variables: []LessonVariable{{Name: "label", Value: "middle"}}},
			{Line: 9, Caption: "Print the new label.", Variables: []LessonVariable{{Name: "label", Value: "middle"}}, Output: "middle"},
		}),
		ch04Sandbox("ch04-ordered-playground", "Step through the first matching condition.", "ch04-ordered-branches", ch04OrderedConfig, true),
	},
	},
	{ID: "c43.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch04Markdown("c43.watch", "- A boundary should belong to one band. Choose ~<~ or ~<=~ deliberately.\n- ~Reading(i)~ returns zero for an invalid index; stay within the count.")}},
	{ID: "c43.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch04Markdown("c43.plan", "Follow the mission contract: one reading, one decision, one label in the same order.\n\nSeparate out-of-range rejection from band choice, then review boundary cases."),
		ch04Question("c43.validate-first", "What should you do before classifying a number?", []string{"Check it is in range", "Print its value", "Choose any band"}, 0, []string{
			"Validation keeps an out-of-range value from entering a normal band.",
			"Printing a value does not check the contract.",
			"Choosing before validation may classify a rejected value.",
		}),
	}},
	{ID: "c43.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch04Markdown("c43.turn", "Start by reading one indexed value and deciding its label yourself. Then cover all readings and run the available editor checks.")}},
}
