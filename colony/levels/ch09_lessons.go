package levels

func ch09Pair(arabic, english string) map[string]any {
	return map[string]any{"ar": arabic, "en": english}
}

var ch09StringerSandbox = &LessonSandbox{
	Kind: "ch09-stringer",
	Config: map[string]any{
		"intro":       ch09Pair("بدّل العرض وشوف Stringer.", "Switch the view; inspect Stringer."),
		"outputLabel": ch09Pair("شكل الـvalue", "Value display"),
		"default":     ch09Pair("{4}", "{4}"),
		"string":      ch09Pair("dock-04", "dock-04"),
		"showString":  ch09Pair("اعرض String()", "Show String()"),
		"showDefault": ch09Pair("اعرض الشكل الافتراضي", "Show default form"),
	},
}

var ch09AutopilotSandbox = &LessonSandbox{
	Kind: "ch09-plan-callback",
	Config: map[string]any{
		"intro":       ch09Pair("Plan بتتسأل عن كل tile.", "Plan is called for each tile."),
		"button":      ch09Pair("التيلة التالية", "Next tile"),
		"tileLabel":   ch09Pair("الـtile", "Tile"),
		"resultLabel": ch09Pair("النتيجة", "Return"),
		"tiles": []any{
			map[string]any{"tile": ch09Pair("Tile A", "Tile A"), "result": ch09Pair("0 actions", "0 actions"), "caption": ch09Pair("المكتبة نادت Plan.", "The library called Plan.")},
			map[string]any{"tile": ch09Pair("Tile B", "Tile B"), "result": ch09Pair("2 actions", "2 actions"), "caption": ch09Pair("النتيجة تخص الخانة دي.", "This return belongs to this tile.")},
			map[string]any{"tile": ch09Pair("Tile C", "Tile C"), "result": ch09Pair("1 action", "1 action"), "caption": ch09Pair("بعدها المكتبة تكمل.", "Then the library continues.")},
		},
	},
}

var ch09SignalsSandbox = &LessonSandbox{
	Kind: "ch09-type-switch",
	Config: map[string]any{
		"intro":       ch09Pair("النوع الحقيقي يختار الـcase.", "The concrete type selects the case."),
		"button":      ch09Pair("الإشارة التالية", "Next value"),
		"valueLabel":  ch09Pair("الـvalue", "Value"),
		"typeLabel":   ch09Pair("النوع", "Type"),
		"caseLabel":   ch09Pair("الـcase", "Case"),
		"effectLabel": ch09Pair("النتيجة", "Effect"),
		"values": []any{
			map[string]any{"value": ch09Pair("42", "42"), "type": ch09Pair("int", "int"), "case": ch09Pair("int", "int"), "effect": ch09Pair("اجمع رقمًا", "Add a number")},
			map[string]any{"value": ch09Pair("beacon", "beacon"), "type": ch09Pair("string", "string"), "case": ch09Pair("string", "string"), "effect": ch09Pair("تعامل مع النص", "Handle text")},
			map[string]any{"value": ch09Pair("Point{2, 5}", "Point{2, 5}"), "type": ch09Pair("Point", "Point"), "case": ch09Pair("Point", "Point"), "effect": ch09Pair("اقرا الحقول", "Read its fields")},
		},
	},
}

var ch09TrimSandbox = &LessonSandbox{
	Kind: "ch09-trim-space",
	Config: map[string]any{
		"intro":       ch09Pair("رمز ␠ بيمثل المسافة؛ جرّب TrimSpace.", "␠ marks spaces; try TrimSpace."),
		"inputLabel":  ch09Pair("قبل", "Before"),
		"outputLabel": ch09Pair("بعد", "After"),
		"raw":         ch09Pair("␠␠beacon-7␠␠", "␠␠beacon-7␠␠"),
		"trimmed":     ch09Pair("beacon-7", "beacon-7"),
		"trim":        ch09Pair("طبّق TrimSpace", "Apply TrimSpace"),
		"restore":     ch09Pair("اعرض المسافات", "Show spaces"),
	},
}

var ch09StringerArabic = []LessonStep{
	{
		ID: "c9m1.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c9m1.mission", Kind: "markdown",
			Text: "1. خلّي `Report(s)` يطلع سطر الموارد بالتنسيق المطلوب.\n2. الـstarter عامل `Status` وجامع الـvalues؛ ركّز على طريقة عرضه.\n3. التقرير بياخد **0 ticks**.",
		}},
	},
	{
		ID: "c9m1.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{
			{ID: "c9m1.string-format", Kind: "markdown", Text: "Go بتطابق الـinterface بالـmethods. `fmt.Stringer` محتاج `String() string`؛ مفيش كلمة `implements`.\n\n`fmt.Sprintf` بيبني نص، و`fmt.Sprint` بيستخدم `String` لو الـvalue عندها الـmethod.\n\nالـmethod بتخلي التنسيق ييجي من الـvalue. `Report` مش محتاج يعرف نوعها؛ ده سلوكها."},
			{ID: "c9m1.format-trace", Kind: "walkthrough", Text: "تنسيق توضيحي صغير؛ الـvalue دي مش بيانات المهمة.", Code: "label := fmt.Sprintf(\"dock-%02d\", 4)\nfmt.Println(label)", Imports: []string{"fmt"}, ExampleID: "e9m1-format", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل التنسيق."},
				{Line: 1, Caption: "Sprintf جهّز النص.", Variables: []LessonVariable{{Name: "label", Value: "dock-04"}}},
				{Line: 2, Caption: "Println عرض النص.", Variables: []LessonVariable{{Name: "label", Value: "dock-04"}}, Output: "dock-04"},
			}}},
		},
	},
	{
		ID: "c9m1.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{ID: "c9m1.stringer-predict", Kind: "question", Text: "تتوقع `fmt.Println(label)` بيطبع إيه؟", Question: &LessonQuestion{
				Options: []string{"90", "1m30s", "Stringer"}, Answer: 1,
				Feedback: []string{"ده شكل الرقم الخام؛ Duration عندها تنسيق نصي.", "`String()` بتاعة Duration بتعرض دقيقة ونص.", "Stringer اسم interface؛ مش النص الناتج."},
			}},
			{ID: "c9m1.stringer-trace", Kind: "walkthrough", Text: "`time.Duration` عندها `String` جاهزة.", Code: "var label fmt.Stringer = 90 * time.Second\nfmt.Println(label)", Imports: []string{"fmt", "time"}, ExampleID: "e9m1-stringer", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل إسناد الـvalue."},
				{Line: 1, Caption: "Duration حققت العقد.", Variables: []LessonVariable{{Name: "label", Value: "time.Duration: 1m30s"}}},
				{Line: 2, Caption: "الطباعة نادت String.", Variables: []LessonVariable{{Name: "label", Value: "time.Duration: 1m30s"}}, Output: "1m30s"},
			}}},
			{ID: "c9m1.stringer-sandbox", Kind: "sandbox", Text: "بدّل ~String()~ وشوف شكل الـvalue.", Sandbox: ch09StringerSandbox},
		},
	},
	{
		ID: "c9m1.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{ID: "c9m1.watch", Kind: "markdown", Text: "اسم الـmethod وحروفه مهمين. `String` لازم ترجّع `string` من غير parameters.\n\n`fmt.Println(s)` بيعرض شكل افتراضي؛ `fmt.Sprintf` هو اللي بيكوّن نص بالشكل اللي تختاره."}},
	},
	{
		ID: "c9m1.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c9m1.plan", Kind: "markdown", Text: "الفكرة: خلي `Status` يوفّر سلوكه النصي، وبعدها ابعت الـvalue الموجودة.\n\nراجع ترتيب الكلمات والأرقام وعلامات الفصل في النص المطلوب."},
			{ID: "c9m1.q-contract", Kind: "question", Text: "إيه اللي يخلّي نوع يحقق `fmt.Stringer`؟", Question: &LessonQuestion{
				Options: []string{"method `String() string`", "حقل اسمه `String`", "method اسمها `Format`"}, Answer: 0,
				Feedback: []string{"توقيع method هو عقد الـinterface.", "الحقول مش methods ومش بتحقق العقد.", "اسم method لازم يطابق `String`."},
			}},
			{ID: "c9m1.q-format", Kind: "question", Text: "أنهي function بترجّع نص متنسّق؟", Question: &LessonQuestion{
				Options: []string{"`fmt.Sprintf`", "`fmt.Println`", "`colony.Report`"}, Answer: 0,
				Feedback: []string{"`Sprintf` بترجّع `string` عشان تستخدمها.", "`Println` بتكتب الناتج، مش بترجّع `string`.", "`Report` بتبعت الـvalue للقاعدة."},
			}},
		},
	},
	{
		ID: "c9m1.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c9m1.turn", Kind: "markdown", Text: "خلّي `Status` يطابق العقد، شغّل الكود، وقارن سطر التقرير."}},
	},
}

var ch09StringerEnglish = []LessonStep{
	{
		ID: "c9m1.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c9m1.mission", Kind: "markdown",
			Text: "1. Make `Report(s)` produce the required resource line.\n2. The starter defines `Status` and gathers the values; focus on how it displays them.\n3. Reporting costs **0 ticks**.",
		}},
	},
	{
		ID: "c9m1.idea-step", Title: "Idea",
		Blocks: []LessonBlock{
			{ID: "c9m1.string-format", Kind: "markdown", Text: "Go matches an interface by methods. `fmt.Stringer` needs `String() string`; there is no `implements` keyword.\n\n`fmt.Sprintf` builds text, and `fmt.Sprint` uses `String` when a value has that method.\n\nThe method lets formatting come from the value. `Report` need not know its concrete type."},
			{ID: "c9m1.format-trace", Kind: "walkthrough", Text: "A small formatting example with illustrative data.", Code: "label := fmt.Sprintf(\"dock-%02d\", 4)\nfmt.Println(label)", Imports: []string{"fmt"}, ExampleID: "e9m1-format", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before formatting."},
				{Line: 1, Caption: "Sprintf builds the text.", Variables: []LessonVariable{{Name: "label", Value: "dock-04"}}},
				{Line: 2, Caption: "Println displays it.", Variables: []LessonVariable{{Name: "label", Value: "dock-04"}}, Output: "dock-04"},
			}}},
		},
	},
	{
		ID: "c9m1.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{ID: "c9m1.stringer-predict", Kind: "question", Text: "What will `fmt.Println(label)` print?", Question: &LessonQuestion{
				Options: []string{"90", "1m30s", "Stringer"}, Answer: 1,
				Feedback: []string{"That is the raw number; Duration has a text format.", "Duration's `String()` displays one minute and thirty seconds.", "Stringer is an interface name, not the output."},
			}},
			{ID: "c9m1.stringer-trace", Kind: "walkthrough", Text: "`time.Duration` already has `String`.", Code: "var label fmt.Stringer = 90 * time.Second\nfmt.Println(label)", Imports: []string{"fmt", "time"}, ExampleID: "e9m1-stringer", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before assigning the value."},
				{Line: 1, Caption: "Duration fits the contract.", Variables: []LessonVariable{{Name: "label", Value: "time.Duration: 1m30s"}}},
				{Line: 2, Caption: "Printing calls String.", Variables: []LessonVariable{{Name: "label", Value: "time.Duration: 1m30s"}}, Output: "1m30s"},
			}}},
			{ID: "c9m1.stringer-sandbox", Kind: "sandbox", Text: "Toggle ~String()~ and inspect the display.", Sandbox: ch09StringerSandbox},
		},
	},
	{
		ID: "c9m1.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{ID: "c9m1.watch", Kind: "markdown", Text: "The method name and capitalization matter. `String` takes no parameters and returns `string`.\n\n`fmt.Println(s)` shows a default form; `fmt.Sprintf` builds text in the format you choose."}},
	},
	{
		ID: "c9m1.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c9m1.plan", Kind: "markdown", Text: "Plan: let `Status` provide its text behavior, then send the value you already have.\n\nCheck the word order, numbers, and separators in the required text."},
			{ID: "c9m1.q-contract", Kind: "question", Text: "What makes a type satisfy `fmt.Stringer`?", Question: &LessonQuestion{
				Options: []string{"method `String() string`", "a field named `String`", "method named `Format`"}, Answer: 0,
				Feedback: []string{"A method signature is the interface contract.", "Fields are not methods and do not satisfy it.", "The method name must match `String`."},
			}},
			{ID: "c9m1.q-format", Kind: "question", Text: "Which function returns formatted text?", Question: &LessonQuestion{
				Options: []string{"`fmt.Sprintf`", "`fmt.Println`", "`colony.Report`"}, Answer: 0,
				Feedback: []string{"`Sprintf` returns text for you to use.", "`Println` writes output; it does not return the text.", "`Report` sends a value to the base."},
			}},
		},
	},
	{
		ID: "c9m1.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c9m1.turn", Kind: "markdown", Text: "Make `Status` satisfy the contract, run the code, and compare the report line."}},
	},
}

var ch09AutopilotArabic = []LessonStep{
	{
		ID: "c9m2.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{ID: "c9m2.mission", Kind: "markdown", Text: "1. خلّي `Autopilot` يتعامل مع كل خانة حسب تربتها ومحتواها.\n2. المطلوب حصاد الناضج وزراعة المحصول المناسب.\n3. الـstarter مجهّز `Farmer` وحالة الجليد."}},
	},
	{
		ID: "c9m2.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{
			{ID: "c9m2.contract", Kind: "markdown", Text: "الطيار بيتولى الـloop على الخريطة، وبينادي `Plan` لكل `Tile`. أي نوع عنده توقيع `Plan(t colony.Tile) []colony.Action` يحقق `Planner`.\n\nالـactions اللي بترجّعها بتتنفذ بالترتيب."},
			{ID: "c9m2.callback-trace", Kind: "walkthrough", Text: "مكتبة Go بتنادي الـcallback لحد ما تلاقي نتيجة.", Code: "i := slices.IndexFunc([]int{2, 5, 8}, func(n int) bool {\n\tfmt.Println(\"check\", n)\n\treturn n > 4\n})\nfmt.Println(\"index\", i)", Imports: []string{"fmt", "slices"}, ExampleID: "e9m2-callback", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل البحث."},
				{Line: 1, Caption: "المكتبة بدأت البحث."},
				{Line: 2, Caption: "أول value: 2.", Variables: []LessonVariable{{Name: "n", Value: "2"}}, Output: "check 2"},
				{Line: 3, Caption: "2 مش مطلوبة؛ كمّل.", Output: "check 2"},
				{Line: 2, Caption: "الـvalue التالية: 5.", Variables: []LessonVariable{{Name: "n", Value: "5"}}, Output: "check 2\ncheck 5"},
				{Line: 3, Caption: "الشرط نجح؛ وقف.", Output: "check 2\ncheck 5"},
				{Line: 5, Caption: "أول index مطابق هو 1.", Variables: []LessonVariable{{Name: "i", Value: "1"}}, Output: "check 2\ncheck 5\nindex 1"},
			}}},
		},
	},
	{
		ID: "c9m2.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{ID: "c9m2.callback-predict", Kind: "question", Text: "تتوقع النتيجة وعدد النداءات كام؟", Question: &LessonQuestion{
				Options: []string{"`1 1`", "`1 2`", "`2 2`"}, Answer: 1,
				Feedback: []string{"اتجربت values قبل ما الشرط ينجح.", "الـindex 1، والـcallback اتنادت مرتين.", "الـindex بيبدأ من صفر، مش عدد الـvalues."},
			}},
			{ID: "c9m2.callback-reveal", Kind: "walkthrough", Text: "الـcallback بتوقف البحث عند أول true.", Code: "calls := 0\nfound := slices.IndexFunc([]int{3, 6}, func(n int) bool {\n\tcalls++\n\treturn n > 4\n})\nfmt.Println(found, calls)", Imports: []string{"fmt", "slices"}, ExampleID: "e9m2-predict", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل أول نداء."},
				{Line: 1, Caption: "العداد بدأ بصفر.", Variables: []LessonVariable{{Name: "calls", Value: "0"}}},
				{Line: 3, Caption: "فحص 3 زوّد العداد.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}},
				{Line: 4, Caption: "3 مش بتحقق الشرط."},
				{Line: 3, Caption: "فحص 6 زوّد العداد.", Variables: []LessonVariable{{Name: "calls", Value: "2"}}},
				{Line: 4, Caption: "6 حققت الشرط؛ توقف."},
				{Line: 6, Caption: "الـindex والنداءات.", Variables: []LessonVariable{{Name: "found", Value: "1"}, {Name: "calls", Value: "2"}}, Output: "1 2"},
			}}},
			{ID: "c9m2.plan-sandbox", Kind: "sandbox", Text: "~Plan~ لكل tile.", Sandbox: ch09AutopilotSandbox},
		},
	},
	{
		ID: "c9m2.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{ID: "c9m2.watch", Kind: "markdown", Text: "لفة واحدة: قراية `Tile` ببلاش، والـactions والحركة كلهم ticks.\n\nالـaction الغلط بيضيّع tick."}},
	},
	{
		ID: "c9m2.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c9m2.plan", Kind: "markdown", Text: "قسّم القرار حسب بيانات الخانة.\n\nراجع هل النبات موجود وناضج، وهل الأرض جاهزة للزراعة أو محتاجة تجهيز."},
			{ID: "c9m2.q-caller", Kind: "question", Text: "مين بينادي `Plan` في المهمة؟", Question: &LessonQuestion{
				Options: []string{"`Autopilot`", "`Farmer`"}, Answer: 0,
				Feedback: []string{"المكتبة بتسأل الـPlanner عن كل خانة.", "Farmer بيرد على النداء؛ مش هو اللي بيلف."},
			}},
			{ID: "c9m2.q-order", Kind: "question", Text: "إمتى الـAutopilot ينفّذ action تانية؟", Question: &LessonQuestion{
				Options: []string{"بعد اللي قبلها في الـslice", "كلهم في نفس اللحظة"}, Answer: 0,
				Feedback: []string{"ترتيب الـslice هو ترتيب التنفيذ.", "الـactions هنا مش متوازية."},
			}},
		},
	},
	{
		ID: "c9m2.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c9m2.turn", Kind: "markdown", Text: "وسّع قرار `Plan` للحالات الناقصة، وبعدها شغّل الطيار."}},
	},
}

var ch09AutopilotEnglish = []LessonStep{
	{
		ID: "c9m2.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{ID: "c9m2.mission", Kind: "markdown", Text: "1. Make `Autopilot` respond to each tile's ground and contents.\n2. Harvest ripe plants and choose a suitable crop.\n3. The starter defines `Farmer` and handles ice."}},
	},
	{
		ID: "c9m2.idea-step", Title: "Idea",
		Blocks: []LessonBlock{
			{ID: "c9m2.contract", Kind: "markdown", Text: "The autopilot owns the map loop and calls `Plan` for each `Tile`. Any type with `Plan(t colony.Tile) []colony.Action` satisfies `Planner`.\n\nReturned actions run in slice order."},
			{ID: "c9m2.callback-trace", Kind: "walkthrough", Text: "A Go library calls a callback until it finds a match.", Code: "i := slices.IndexFunc([]int{2, 5, 8}, func(n int) bool {\n\tfmt.Println(\"check\", n)\n\treturn n > 4\n})\nfmt.Println(\"index\", i)", Imports: []string{"fmt", "slices"}, ExampleID: "e9m2-callback", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before searching."},
				{Line: 1, Caption: "The library starts."},
				{Line: 2, Caption: "First value: 2.", Variables: []LessonVariable{{Name: "n", Value: "2"}}, Output: "check 2"},
				{Line: 3, Caption: "2 fails; keep going.", Output: "check 2"},
				{Line: 2, Caption: "Next value: 5.", Variables: []LessonVariable{{Name: "n", Value: "5"}}, Output: "check 2\ncheck 5"},
				{Line: 3, Caption: "The test passes; stop.", Output: "check 2\ncheck 5"},
				{Line: 5, Caption: "First matching index: 1.", Variables: []LessonVariable{{Name: "i", Value: "1"}}, Output: "check 2\ncheck 5\nindex 1"},
			}}},
		},
	},
	{
		ID: "c9m2.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{ID: "c9m2.callback-predict", Kind: "question", Text: "Predict the result and number of calls.", Question: &LessonQuestion{
				Options: []string{"`1 1`", "`1 2`", "`2 2`"}, Answer: 1,
				Feedback: []string{"Two values are checked before the test passes.", "The index is 1, and the callback is called twice.", "The index starts at zero; it is not the number of values."},
			}},
			{ID: "c9m2.callback-reveal", Kind: "walkthrough", Text: "The callback stops the search at the first true.", Code: "calls := 0\nfound := slices.IndexFunc([]int{3, 6}, func(n int) bool {\n\tcalls++\n\treturn n > 4\n})\nfmt.Println(found, calls)", Imports: []string{"fmt", "slices"}, ExampleID: "e9m2-predict", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before the first call."},
				{Line: 1, Caption: "The counter starts at zero.", Variables: []LessonVariable{{Name: "calls", Value: "0"}}},
				{Line: 3, Caption: "Checking 3 adds one.", Variables: []LessonVariable{{Name: "calls", Value: "1"}}},
				{Line: 4, Caption: "3 fails the test."},
				{Line: 3, Caption: "Checking 6 adds one.", Variables: []LessonVariable{{Name: "calls", Value: "2"}}},
				{Line: 4, Caption: "6 passes; stop."},
				{Line: 6, Caption: "The index and call count.", Variables: []LessonVariable{{Name: "found", Value: "1"}, {Name: "calls", Value: "2"}}, Output: "1 2"},
			}}},
			{ID: "c9m2.plan-sandbox", Kind: "sandbox", Text: "~Plan~ per tile.", Sandbox: ch09AutopilotSandbox},
		},
	},
	{
		ID: "c9m2.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{ID: "c9m2.watch", Kind: "markdown", Text: "One pass only: reading a `Tile` is free, but actions and moves each cost a tick. Wrong actions waste ticks."}},
	},
	{
		ID: "c9m2.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c9m2.plan", Kind: "markdown", Text: "Split the decision by tile data.\n\nCheck whether a plant exists and is ripe, and whether the ground is ready or needs preparation."},
			{ID: "c9m2.q-caller", Kind: "question", Text: "Who calls `Plan` in this mission?", Question: &LessonQuestion{
				Options: []string{"`Autopilot`", "`Farmer`"}, Answer: 0,
				Feedback: []string{"The library asks the planner about each tile.", "Farmer answers the call; it does not own the loop."},
			}},
			{ID: "c9m2.q-order", Kind: "question", Text: "When does the autopilot run another action?", Question: &LessonQuestion{
				Options: []string{"After the previous slice action", "All at the same time"}, Answer: 0,
				Feedback: []string{"Slice order is execution order.", "These actions are not concurrent."},
			}},
		},
	},
	{
		ID: "c9m2.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c9m2.turn", Kind: "markdown", Text: "Extend `Plan` for the missing cases, then run the autopilot."}},
	},
}

var ch09SignalsArabic = []LessonStep{
	{
		ID: "c9m3.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{ID: "c9m3.mission", Kind: "markdown", Text: "1. راجع كل رسالة راديو وحدد نوعها.\n2. روح لمكان كل `Distress`، واجمع أرقام `Supply` للتقرير.\n3. `Noise` يتساب؛ ترتيب الرسائل عشوائي."}},
	},
	{
		ID: "c9m3.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{
			{ID: "c9m3.type-switch", Kind: "markdown", Text: "الـinterface شايل نوع حقيقي والـvalue بتاعته. `colony.Signals()` بترجّع `[]colony.Signal`؛ كل عنصر `Distress` أو `Supply` أو `Noise`.\n\nفي `switch v := s.(type)`، نوع `v` بيتحدد حسب الـcase. كل case ليها fields مختلفة؛ استخدم اللي يناسب النوع."},
			{ID: "c9m3.switch-trace", Kind: "walkthrough", Text: "مثال فيه values توضيحية: رقم ونص ثابت.", Code: "total := 0\nfor _, msg := range []any{7, \"static\"} {\n\tswitch v := msg.(type) {\n\tcase int:\n\t\ttotal += v\n\tcase string:\n\t\tfmt.Println(\"ignore\", v)\n\t}\n}\nfmt.Println(\"sum\", total)", Imports: []string{"fmt"}, ExampleID: "e9m3-switch", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل قراءة الرسائل."},
				{Line: 1, Caption: "المجموع بدأ بصفر.", Variables: []LessonVariable{{Name: "total", Value: "0"}}},
				{Line: 2, Caption: "أول value رقم 7.", Variables: []LessonVariable{{Name: "msg", Value: "int(7)"}}},
				{Line: 5, Caption: "حالة int زوّدت المجموع.", Variables: []LessonVariable{{Name: "total", Value: "7"}}},
				{Line: 2, Caption: "الـvalue التالية نص ثابت.", Variables: []LessonVariable{{Name: "msg", Value: "string(static)"}}},
				{Line: 8, Caption: "حالة النص تجاهلته.", Variables: []LessonVariable{{Name: "total", Value: "7"}}, Output: "ignore static"},
				{Line: 10, Caption: "اطبع الإجمالي مرة.", Variables: []LessonVariable{{Name: "total", Value: "7"}}, Output: "ignore static\nsum 7"},
			}}},
		},
	},
	{
		ID: "c9m3.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{ID: "c9m3.assert-predict", Kind: "question", Text: "الـvalue نص `\"3\"`؛ تتوقع يطلع إيه؟", Question: &LessonQuestion{
				Options: []string{"`4`", "`text 3`", "panic"}, Answer: 1,
				Feedback: []string{"النوع string؛ الكود مش بيحوّله لرقم.", "الـcase المطابق هو string.", "الـtype switch بيتعامل مع النوع بأمان."},
			}},
			{ID: "c9m3.assert-reveal", Kind: "walkthrough", Text: "الـcase بيتحدد بالنوع الحقيقي.", Code: "signal := any(\"3\")\nswitch v := signal.(type) {\ncase int:\n\tfmt.Println(v + 1)\ncase string:\n\tfmt.Println(\"text\", v)\n}", Imports: []string{"fmt"}, ExampleID: "e9m3-type-check", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل فحص النوع."},
				{Line: 1, Caption: "الـvalue نوعها string.", Variables: []LessonVariable{{Name: "signal", Value: "string(3)"}}},
				{Line: 2, Caption: "الـswitch اختار string."},
				{Line: 6, Caption: "اطبع النوع والنص.", Variables: []LessonVariable{{Name: "v", Value: "3"}}, Output: "text 3"},
			}}},
			{ID: "c9m3.type-switch-sandbox", Kind: "sandbox", Text: "حرّك الـvalues وشوف الـcase المطابق.", Sandbox: ch09SignalsSandbox},
		},
	},
	{
		ID: "c9m3.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{ID: "c9m3.watch", Kind: "markdown", Text: "الـassertion من غير `ok` ممكن تعمل panic. `Move` tick حتى لو اتقفلت؛ `Report` تكلفته صفر."}},
	},
	{
		ID: "c9m3.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c9m3.plan", Kind: "markdown", Text: "اقرا الإشارات الأول؛ وبعدها زور المواقع وابعت المجموع. ترتيبها عشوائي."},
			{ID: "c9m3.q-fields", Kind: "question", Text: "أنهي نوع فيه `At`؟", Question: &LessonQuestion{
				Options: []string{"`Distress`", "`Supply`", "`Noise`"}, Answer: 0,
				Feedback: []string{"موقع الزيارة موجود في `Distress.At`.", "`Supply` شايلة `Amount`.", "`Noise` نص ثابت ملوش موقع."},
			}},
			{ID: "c9m3.q-assert", Kind: "question", Text: "إيه خطر assertion مباشرة على رسالة من نوع تاني؟", Question: &LessonQuestion{
				Options: []string{"panic", "بترجّع nil", "تتحول تلقائيًا"}, Answer: 0,
				Feedback: []string{"من غير `ok`، النوع الغلط يعمل panic.", "الـvalue مش pointer؛ مش هترجّع nil.", "Go مش بتحوّل concrete types لوحدها."},
			}},
		},
	},
	{
		ID: "c9m3.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c9m3.turn", Kind: "markdown", Text: "بدّل طباعة التشخيص بجمع البيانات والزيارات المطلوبة، وبعدها راجع التقرير."}},
	},
}

var ch09SignalsEnglish = []LessonStep{
	{
		ID: "c9m3.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{ID: "c9m3.mission", Kind: "markdown", Text: "1. Read each radio message and identify its type.\n2. Visit every `Distress` location and add `Supply` amounts for the report.\n3. Ignore `Noise`; message order is shuffled."}},
	},
	{
		ID: "c9m3.idea-step", Title: "Idea",
		Blocks: []LessonBlock{
			{ID: "c9m3.type-switch", Kind: "markdown", Text: "An interface holds a concrete type and value. `colony.Signals()` returns `[]colony.Signal`; each item is `Distress`, `Supply`, or `Noise`.\n\nIn `switch v := s.(type)`, `v` gets the type of its matching case. Each case exposes different fields."},
			{ID: "c9m3.switch-trace", Kind: "walkthrough", Text: "An example with illustrative values: a number and static.", Code: "total := 0\nfor _, msg := range []any{7, \"static\"} {\n\tswitch v := msg.(type) {\n\tcase int:\n\t\ttotal += v\n\tcase string:\n\t\tfmt.Println(\"ignore\", v)\n\t}\n}\nfmt.Println(\"sum\", total)", Imports: []string{"fmt"}, ExampleID: "e9m3-switch", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before reading messages."},
				{Line: 1, Caption: "The total starts at zero.", Variables: []LessonVariable{{Name: "total", Value: "0"}}},
				{Line: 2, Caption: "First value is int 7.", Variables: []LessonVariable{{Name: "msg", Value: "int(7)"}}},
				{Line: 5, Caption: "The int case adds it.", Variables: []LessonVariable{{Name: "total", Value: "7"}}},
				{Line: 2, Caption: "Next value is static text.", Variables: []LessonVariable{{Name: "msg", Value: "string(static)"}}},
				{Line: 8, Caption: "The text case ignores it.", Variables: []LessonVariable{{Name: "total", Value: "7"}}, Output: "ignore static"},
				{Line: 10, Caption: "Print the total once.", Variables: []LessonVariable{{Name: "total", Value: "7"}}, Output: "ignore static\nsum 7"},
			}}},
		},
	},
	{
		ID: "c9m3.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{ID: "c9m3.assert-predict", Kind: "question", Text: "The value is the string `\"3\"`. What prints?", Question: &LessonQuestion{
				Options: []string{"`4`", "`text 3`", "panic"}, Answer: 1,
				Feedback: []string{"Its type is string; this code does not convert it to a number.", "The matching case is string.", "A type switch checks the type safely."},
			}},
			{ID: "c9m3.assert-reveal", Kind: "walkthrough", Text: "The real type chooses the case.", Code: "signal := any(\"3\")\nswitch v := signal.(type) {\ncase int:\n\tfmt.Println(v + 1)\ncase string:\n\tfmt.Println(\"text\", v)\n}", Imports: []string{"fmt"}, ExampleID: "e9m3-type-check", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before checking the type."},
				{Line: 1, Caption: "The signal is a string.", Variables: []LessonVariable{{Name: "signal", Value: "string(3)"}}},
				{Line: 2, Caption: "The switch chooses string."},
				{Line: 6, Caption: "Print the type and text.", Variables: []LessonVariable{{Name: "v", Value: "3"}}, Output: "text 3"},
			}}},
			{ID: "c9m3.type-switch-sandbox", Kind: "sandbox", Text: "Cycle values and inspect each matching case.", Sandbox: ch09SignalsSandbox},
		},
	},
	{
		ID: "c9m3.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{ID: "c9m3.watch", Kind: "markdown", Text: "An assertion without `ok` can panic. `Move` costs a tick even when blocked; `Report` costs zero."}},
	},
	{
		ID: "c9m3.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c9m3.plan", Kind: "markdown", Text: "Read signals first; visit their sites and report afterward. Their order is shuffled."},
			{ID: "c9m3.q-fields", Kind: "question", Text: "Which type has `At`?", Question: &LessonQuestion{
				Options: []string{"`Distress`", "`Supply`", "`Noise`"}, Answer: 0,
				Feedback: []string{"The visit location is `Distress.At`.", "`Supply` carries `Amount`.", "`Noise` is static text with no location."},
			}},
			{ID: "c9m3.q-assert", Kind: "question", Text: "What can a direct assertion to another message type do?", Question: &LessonQuestion{
				Options: []string{"panic", "return nil", "convert automatically"}, Answer: 0,
				Feedback: []string{"Without `ok`, the wrong type panics.", "The value is not a pointer that returns nil.", "Go does not convert concrete types automatically."},
			}},
		},
	},
	{
		ID: "c9m3.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c9m3.turn", Kind: "markdown", Text: "Replace the diagnostic print with the required collection and visits, then check the report."}},
	},
}

var ch09PlugsArabic = []LessonStep{
	{
		ID: "c9m4.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{ID: "c9m4.mission", Kind: "markdown", Text: "1. خلّي `RelayOrders` بتقرا أوامر من أي `io.Reader` وتكتب في أي `io.Writer`.\n2. كل سطر key؛ العنوان أو `unknown` حتى للسطر الفاضي. `end` يقفل، والـEOF العادي `exhausted`.\n3. رجّع الـerrors اللي بتحصل في القراية والكتابة."}},
	},
	{
		ID: "c9m4.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{
			{ID: "c9m4.scanner", Kind: "markdown", Text: "الـinterface الصغير يطلب السلوك اللي محتاجه بس. `io.Reader` يوفّر `Read`، و`io.Writer` يوفّر `Write`.\n\n`Scanner` يجيب سطر؛ `Text` يجيبه من غير newline."},
			{ID: "c9m4.scan-trace", Kind: "walkthrough", Text: "اقرا السطر ونضّف المسافات الطرفية.", Code: "scanner := bufio.NewScanner(strings.NewReader(\"  beacon-7 \\n\"))\nif scanner.Scan() {\n\tkey := strings.TrimSpace(scanner.Text())\n\tfmt.Println(key)\n}\nfmt.Println(scanner.Err())", Imports: []string{"bufio", "fmt", "strings"}, ExampleID: "e9m4-scan", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل إنشاء القارئ."},
				{Line: 1, Caption: "Reader يزوّد Scanner بسطر."},
				{Line: 2, Caption: "Scan جابت السطر."},
				{Line: 3, Caption: "TrimSpace شالت الطرفيات.", Variables: []LessonVariable{{Name: "key", Value: "beacon-7"}}},
				{Line: 4, Caption: "اطبع الـkey.", Variables: []LessonVariable{{Name: "key", Value: "beacon-7"}}, Output: "beacon-7"},
				{Line: 6, Caption: "مفيش error في القراية.", Output: "beacon-7\n<nil>"},
			}}},
		},
	},
	{
		ID: "c9m4.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{ID: "c9m4.trim-predict", Kind: "question", Text: "بعد `TrimSpace`، تتوقع المقارنة والطباعة؟", Question: &LessonQuestion{
				Options: []string{"`false end`", "`true end`", "`true  end`"}, Answer: 1,
				Feedback: []string{"المسافات الطرفية اتشالت قبل المقارنة.", "الـkey بقى `end` من غير المسافات.", "المسافة اللي قبل end مش جزء من الـvalue."},
			}},
			{ID: "c9m4.trim-reveal", Kind: "walkthrough", Text: "القراية مش بتغيّر النص؛ TrimSpace هي اللي بتشيله.", Code: "key := strings.TrimSpace(\" end \\n\")\nfmt.Println(key == \"end\", key)", Imports: []string{"fmt", "strings"}, ExampleID: "e9m4-trim", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل التنضيف."},
				{Line: 1, Caption: "الـkey بقى end.", Variables: []LessonVariable{{Name: "key", Value: "end"}}},
				{Line: 2, Caption: "المقارنة والنص.", Variables: []LessonVariable{{Name: "key", Value: "end"}}, Output: "true end"},
			}}},
			{ID: "c9m4.trim-sandbox", Kind: "sandbox", Text: "بدّل ~TrimSpace~.", Sandbox: ch09TrimSandbox},
		},
	},
	{
		ID: "c9m4.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{ID: "c9m4.watch", Kind: "markdown", Text: "`end` أمر صريح: اكتب `closed` واعمل `return` قبل lookup. أما `Scan` اللي وقفت، رجّع `Err` لو مش `nil`؛ `nil` معناها EOF عادي.\n\nافحص الـerror بعد كل كتابة، ومتستخدمش `fmt.Println` بدل الـwriter اللي اتبعتلك."}},
	},
	{
		ID: "c9m4.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c9m4.plan", Kind: "markdown", Text: "الخطة: افصل المصدر والوجهة عن منطق الأمر.\n\nلكل سطر، نظّفه وافحصه؛ بعد توقف Scanner ميّز EOF عن error."},
			{ID: "c9m4.q-input", Kind: "question", Text: "إيه اللي محتاجه المنطق من المصدر؟", Question: &LessonQuestion{
				Options: []string{"`io.Reader`", "ملف مفتوح بعينه", "`os.Stdin`"}, Answer: 0,
				Feedback: []string{"أي مصدر يحقق Read ينفع.", "ربط الـfunction بملف بعينه يقلل إعادة استخدامها.", "المدخل ييجي كـparameter؛ مش global."},
			}},
			{ID: "c9m4.q-end", Kind: "question", Text: "إيه الفرق بين `end` وEOF؟", Question: &LessonQuestion{
				Options: []string{"`end` أمر؛ EOF نهاية إدخال", "الاتنين نفس المعنى"}, Answer: 0,
				Feedback: []string{"`end` key مقروء؛ EOF مفيش سطر تاني.", "نهاية الإدخال مش أمر إغلاق."},
			}},
		},
	},
	{
		ID: "c9m4.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{ID: "c9m4.turn", Kind: "markdown", Text: "نفّذ العقد جوه `RelayOrders`، وجرّبه بـreader وbuffer في الذاكرة."}},
	},
}

var ch09PlugsEnglish = []LessonStep{
	{
		ID: "c9m4.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{ID: "c9m4.mission", Kind: "markdown", Text: "1. Make `RelayOrders` read from any `io.Reader` and write to any `io.Writer`.\n2. Each trimmed key prints its address or `unknown`, even blank; `end` closes; normal EOF prints `exhausted`.\n3. Return read and write errors."}},
	},
	{
		ID: "c9m4.idea-step", Title: "Idea",
		Blocks: []LessonBlock{
			{ID: "c9m4.scanner", Kind: "markdown", Text: "A small interface asks only for needed behavior. `io.Reader` provides `Read`; `io.Writer` provides `Write`.\n\n`Scanner` gets a line; `Text` returns it without the newline."},
			{ID: "c9m4.scan-trace", Kind: "walkthrough", Text: "Read a line and trim its surrounding spaces.", Code: "scanner := bufio.NewScanner(strings.NewReader(\"  beacon-7 \\n\"))\nif scanner.Scan() {\n\tkey := strings.TrimSpace(scanner.Text())\n\tfmt.Println(key)\n}\nfmt.Println(scanner.Err())", Imports: []string{"bufio", "fmt", "strings"}, ExampleID: "e9m4-scan", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before creating the reader."},
				{Line: 1, Caption: "Reader gives Scanner a line."},
				{Line: 2, Caption: "Scan gets the line."},
				{Line: 3, Caption: "TrimSpace removes the edges.", Variables: []LessonVariable{{Name: "key", Value: "beacon-7"}}},
				{Line: 4, Caption: "Print the key.", Variables: []LessonVariable{{Name: "key", Value: "beacon-7"}}, Output: "beacon-7"},
				{Line: 6, Caption: "There was no read error.", Output: "beacon-7\n<nil>"},
			}}},
		},
	},
	{
		ID: "c9m4.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{ID: "c9m4.trim-predict", Kind: "question", Text: "After `TrimSpace`, predict the comparison and output.", Question: &LessonQuestion{
				Options: []string{"`false end`", "`true end`", "`true  end`"}, Answer: 1,
				Feedback: []string{"Surrounding spaces are removed before comparison.", "The key becomes `end` without spaces.", "The space before end is not part of the value."},
			}},
			{ID: "c9m4.trim-reveal", Kind: "walkthrough", Text: "Reading does not trim text; TrimSpace does.", Code: "key := strings.TrimSpace(\" end \\n\")\nfmt.Println(key == \"end\", key)", Imports: []string{"fmt", "strings"}, ExampleID: "e9m4-trim", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before trimming."},
				{Line: 1, Caption: "The key becomes end.", Variables: []LessonVariable{{Name: "key", Value: "end"}}},
				{Line: 2, Caption: "The comparison and text.", Variables: []LessonVariable{{Name: "key", Value: "end"}}, Output: "true end"},
			}}},
			{ID: "c9m4.trim-sandbox", Kind: "sandbox", Text: "Try ~TrimSpace~.", Sandbox: ch09TrimSandbox},
		},
	},
	{
		ID: "c9m4.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{ID: "c9m4.watch", Kind: "markdown", Text: "`end` is explicit: write `closed` and return before lookup. When `Scan` stops, return a nonnil `Err`; `nil` means normal EOF.\n\nCheck every write error, and do not use `fmt.Println` instead of the supplied writer."}},
	},
	{
		ID: "c9m4.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c9m4.plan", Kind: "markdown", Text: "Plan: separate the source and destination from command handling.\n\nFor each line, trim and inspect it; after Scanner stops, distinguish EOF from an error."},
			{ID: "c9m4.q-input", Kind: "question", Text: "What does the logic need from its source?", Question: &LessonQuestion{
				Options: []string{"`io.Reader`", "one specific open file", "`os.Stdin`"}, Answer: 0,
				Feedback: []string{"Any source that provides Read works.", "Binding the function to one file reduces reuse.", "Input arrives as a parameter, not a global."},
			}},
			{ID: "c9m4.q-end", Kind: "question", Text: "How do `end` and EOF differ?", Question: &LessonQuestion{
				Options: []string{"`end` is a command; EOF ends input", "They mean the same thing"}, Answer: 0,
				Feedback: []string{"`end` is a read key; EOF means no next line.", "Input exhaustion is not a close command."},
			}},
		},
	},
	{
		ID: "c9m4.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{ID: "c9m4.turn", Kind: "markdown", Text: "Implement `RelayOrders`, then try a reader and buffer."}},
	},
}
