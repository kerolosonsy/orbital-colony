package levels

func ch21Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch21Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch21Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

func ch21Pair(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

func ch21SandboxBlock(id, kind, text string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch21StationSandboxConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"flag": ch21Pair("Flag -details", "Flag -details"), "off": ch21Pair("false", "false"), "on": ch21Pair("true", "true"),
			"toggle": ch21Pair("بدّل الـflag", "Toggle flag"), "lookup": ch21Pair("مصدر البيئة", "Environment lookup"),
			"environment": ch21Pair("os.LookupEnv", "os.LookupEnv"),
			"name":        ch21Pair("الاسم", "Name"), "presence": ch21Pair("الوجود", "Presence"), "value": ch21Pair("القيمة", "Value"),
			"present": ch21Pair("موجود", "Present"), "missing": ch21Pair("غايب", "Missing"), "next": ch21Pair("الحالة الجاية", "Next state"),
			"caption": ch21Pair("الـflag تبدأ false.", "The flag starts false."),
		},
		"states": []any{
			map[string]any{"label": ch21Pair("اسم غايب", "Missing name"), "value": ch21Pair("مش موجود", "Not set"), "present": false},
			map[string]any{"label": ch21Pair("قيمة فاضية", "Empty value"), "value": ch21Pair("\"\"", "\"\""), "present": true},
			map[string]any{"label": ch21Pair("قيمة موجودة", "Present value"), "value": ch21Pair("crew-7", "crew-7"), "present": true},
		},
	}
}

func ch21CommandsSandboxConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"caption": ch21Pair("كل Scan تكشف سطرًا واحدًا.", "Each Scan exposes one line."),
			"current": ch21Pair("السطر الحالي", "Current line"), "written": ch21Pair("سطور وصلت للـwriter", "Lines sent to writer"),
			"previous": ch21Pair("السابق", "Previous"), "next": ch21Pair("التالي", "Next"), "reset": ch21Pair("من الأول", "Reset"),
		},
		"frames": []any{
			map[string]any{"caption": ch21Pair("الـReader فيه تلات سطور.", "The Reader holds three lines."), "current": ch21Pair("—", "—"), "written": []any{}},
			map[string]any{"caption": ch21Pair("Scan تعرض السطر الأول.", "Scan exposes the first line."), "current": ch21Pair("ping one", "ping one"), "written": []any{ch21Pair("ping one", "ping one")}},
			map[string]any{"caption": ch21Pair("السطر الحالي يتكتب مرة.", "The current line is written once."), "current": ch21Pair("ping two", "ping two"), "written": []any{ch21Pair("ping one", "ping one"), ch21Pair("ping two", "ping two")}},
			map[string]any{"caption": ch21Pair("كل سطر له دورة منفصلة.", "Each line has its own pass."), "current": ch21Pair("ping three", "ping three"), "written": []any{ch21Pair("ping one", "ping one"), ch21Pair("ping two", "ping two"), ch21Pair("ping three", "ping three")}},
		},
	}
}

func ch21DepthSandboxConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"source": ch21Pair("مصدر الاختيار", "Choice source"), "fixed": ch21Pair("callback ثابت", "Fixed callback"),
			"varied": ch21Pair("سحبات توضيحية", "Sample random draws"), "draw": ch21Pair("اسحب index", "Draw index"),
			"bound": ch21Pair("الحد exclusive", "Exclusive bound"), "index": ch21Pair("index", "index"), "selected": ch21Pair("الاسم المختار", "Selected label"),
			"calls": ch21Pair("عدد النداءات", "Calls"), "none": ch21Pair("لسه مفيش اختيار", "No choice yet"),
			"caption": ch21Pair("المصدر يرجع index أصغر من الحد.", "The source returns an index below the bound."),
		},
		"choices":      []any{ch21Pair("amber", "amber"), ch21Pair("cyan", "cyan"), ch21Pair("violet", "violet")},
		"fixedResults": []int{1, 1, 1}, "variedResults": []int{0, 2, 1},
	}
}

func ch21DuelSandboxConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"caption": ch21Pair("اسم النوع يفصل القيمة عن النص.", "A named type separates value from text."),
			"command": ch21Pair("نص الأمر", "Command text"), "named": ch21Pair("القيمة المسماة", "Named value"),
			"ordinal": ch21Pair("قيمة iota", "iota value"), "next": ch21Pair("أمر توضيحي تاني", "Next example command"),
		},
		"commands": []any{
			map[string]any{"text": ch21Pair("idle", "idle"), "name": ch21Pair("LampMode.idle", "LampMode.idle"), "ordinal": 0},
			map[string]any{"text": ch21Pair("scan", "scan"), "name": ch21Pair("LampMode.scan", "LampMode.scan"), "ordinal": 1},
			map[string]any{"text": ch21Pair("relay", "relay"), "name": ch21Pair("LampMode.relay", "LampMode.relay"), "ordinal": 2},
		},
	}
}

var ch21StationArabic = []LessonStep{
	{ID: "c211.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch21Markdown("c211.mission", "1. جهّز تقرير ثابت بين الأجهزة.\n2. خلي بيانات النظام توصل من برّه.\n3. الـstarter فيه نوع اللقطة؛ ركّز على فصل المصدر عن التقرير.")}},
	{ID: "c211.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch21Markdown("c211.boundary", "الـflag اختيار جاي من المستخدم. ~FlagSet~ بتجمع تعريفاته، و~Bool~ بترجع pointer لقيمته.\n\nالبيئة حاجة تانية: function value تسأل عن اسم، والـbool بيفرق بين الغياب والقيمة الفاضية.\n\nاجمع بيانات الجهاز الحقيقي عند ~main~؛ التقرير يستقبل قيم جاهزة."),
		ch21Walkthrough("c211.flags-example", "تابع parser محلي؛ نتيجة الـflag pointer، و~NArg~ بتعدّ البواقي.", `func main() {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	active := fs.Bool("active", false, "show details")
	_ = fs.Parse([]string{"-active"})
	fmt.Println(*active, fs.NArg())
}`, "e211-local-flag", []string{"flag", "fmt", "io"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "ابدأ parser محلي."},
			{Line: 4, Caption: "سجّل flag بقيمة افتراضية."},
			{Line: 5, Caption: "حلّل اختيار المثال."},
			{Line: 6, Caption: "القيمة true والبواقي صفر.", Variables: []LessonVariable{{Name: "active", Value: "true"}}, Output: "true 0"},
		}),
	}},
	{ID: "c211.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch21Question("c211.predict-flag", "توقّع الناتج بعد parsing:\n\n~~~go\nfunc main() {\n\tfs := flag.NewFlagSet(\"probe\", flag.ContinueOnError)\n\tfs.SetOutput(io.Discard)\n\tactive := fs.Bool(\"active\", false, \"show details\")\n\t_ = fs.Parse([]string{\"-active\"})\n\tfmt.Println(*active, fs.NArg())\n}\n~~~", []string{"true 0", "false 0", "true 1"}, 0, []string{
			"الاختيار اتقري ومفيش argument فاضل.",
			"الـflag اتبعتت بقيمة true.",
			"الـparser استهلك الـflag المعروفة.",
		}),
		ch21Walkthrough("c211.flag-reveal", "شوف القيمة اللي الـflag سجلتها بعد التحليل.", `func main() {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	active := fs.Bool("active", false, "show details")
	_ = fs.Parse([]string{"-active"})
	fmt.Println(*active, fs.NArg())
}`, "e211-flag-reveal", []string{"flag", "fmt", "io"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "parser جديد للاستدعاء."},
			{Line: 4, Caption: "القيمة تبدأ false."},
			{Line: 5, Caption: "active بقت مفعّلة."},
			{Line: 6, Caption: "اطبع القيمة والبواقي.", Variables: []LessonVariable{{Name: "active", Value: "true"}}, Output: "true 0"},
		}),
		ch21SandboxBlock("c211.playground", "ch21-flags-env", "بدّل الـflag وقارن الاسم الغايب بالقيمة الفاضية.", ch21StationSandboxConfig()),
	}},
	{ID: "c211.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch21Markdown("c211.watch", "- pointer من ~Bool~ بيتقري بـ~*value~.\n- القيمة الفاضية مش دليل إن اسم البيئة غايب.")}},
	{ID: "c211.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch21Markdown("c211.plan", "قبل أول output، حدّد إيه اللي ممكن يرفض الطلب.\n\nراجع حالات الـflag: متبعتة، أو متسابِة على default.\n\nإمتى التقرير يحتاج يسأل عن اسم البيئة؟"),
		ch21Question("c211.presence-check", "إيه اللي يفرّق بين اسم غايب وقيمة فاضية؟", []string{"النص بس", "النص ومعاه bool", "طباعة النص"}, 1, []string{
			"النص الفاضي ممكن يبقى قيمة موجودة.",
			"الـbool بيوضح هل الاسم موجود.",
			"الطباعة مش بتكشف وجود الاسم.",
		}),
	}},
	{ID: "c211.turn-step", Title: "دورك", Blocks: []LessonBlock{ch21Markdown("c211.turn", "اكتب التقرير على writer محقون. ابدأ بقيم ثابتة، وراجع إن المصدر الحقيقي متجمع عند حدود البرنامج.")}},
}

var ch21StationEnglish = []LessonStep{
	{ID: "c211.mission-step", Title: "Mission", Blocks: []LessonBlock{ch21Markdown("c211.mission", "1. Prepare a stable report across machines.\n2. Supply host data from outside the report.\n3. The starter includes the snapshot type; focus on separating source data from the report.")}},
	{ID: "c211.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch21Markdown("c211.boundary", "A flag is a user choice. ~FlagSet~ groups its definitions, and ~Bool~ returns a pointer to its value.\n\nEnvironment access is separate: a function value looks up a name, and its bool distinguishes missing from empty.\n\nGather real host data at ~main~; the report receives ready values."),
		ch21Walkthrough("c211.flags-example", "Follow a local parser; the flag is a pointer and ~NArg~ counts leftovers.", `func main() {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	active := fs.Bool("active", false, "show details")
	_ = fs.Parse([]string{"-active"})
	fmt.Println(*active, fs.NArg())
}`, "e211-local-flag", []string{"flag", "fmt", "io"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Start a local parser."},
			{Line: 4, Caption: "Register a default flag."},
			{Line: 5, Caption: "Parse the example choice."},
			{Line: 6, Caption: "Value true; zero leftovers.", Variables: []LessonVariable{{Name: "active", Value: "true"}}, Output: "true 0"},
		}),
	}},
	{ID: "c211.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch21Question("c211.predict-flag", "Predict the output after parsing:\n\n~~~go\nfunc main() {\n\tfs := flag.NewFlagSet(\"probe\", flag.ContinueOnError)\n\tfs.SetOutput(io.Discard)\n\tactive := fs.Bool(\"active\", false, \"show details\")\n\t_ = fs.Parse([]string{\"-active\"})\n\tfmt.Println(*active, fs.NArg())\n}\n~~~", []string{"true 0", "false 0", "true 1"}, 0, []string{
			"The choice was read with no argument left.",
			"The flag was supplied as true.",
			"The parser consumed the known flag.",
		}),
		ch21Walkthrough("c211.flag-reveal", "Watch the value recorded after parsing.", `func main() {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	active := fs.Bool("active", false, "show details")
	_ = fs.Parse([]string{"-active"})
	fmt.Println(*active, fs.NArg())
}`, "e211-flag-reveal", []string{"flag", "fmt", "io"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "A fresh parser for this call."},
			{Line: 4, Caption: "The default starts false."},
			{Line: 5, Caption: "active is now enabled."},
			{Line: 6, Caption: "Print value and leftovers.", Variables: []LessonVariable{{Name: "active", Value: "true"}}, Output: "true 0"},
		}),
		ch21SandboxBlock("c211.playground", "ch21-flags-env", "Toggle the flag; compare missing and empty names.", ch21StationSandboxConfig()),
	}},
	{ID: "c211.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch21Markdown("c211.watch", "- Read a ~Bool~ pointer with ~*value~.\n- An empty value does not prove an environment name is missing.")}},
	{ID: "c211.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch21Markdown("c211.plan", "Before the first output, identify what can reject a request.\n\nReview both flag states: supplied, or left at its default.\n\nWhen should the report ask for an environment name?"),
		ch21Question("c211.presence-check", "What distinguishes a missing name from an empty value?", []string{"The text alone", "Text plus a bool", "Printing the text"}, 1, []string{
			"Empty text can still be present.",
			"The bool reports whether the name exists.",
			"Printing cannot reveal presence.",
		}),
	}},
	{ID: "c211.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch21Markdown("c211.turn", "Write the report to an injected writer. Start with fixed values, and keep real host access at the program boundary.")}},
}

var ch21CommandsArabic = []LessonStep{
	{ID: "c212.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch21Markdown("c212.mission", "1. شغّل أمر واحد أو سيل سطور.\n2. افصل الواجهة عن سجل الطلبات.\n3. الـstarter جاهز لحدود البرنامج؛ كمّل المنطق القابل للاختبار.")}},
	{ID: "c212.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch21Markdown("c212.stream", "~io.Reader~ و~io.Writer~ من 9-4 بيفصلوا المنطق عن التيرمنال. الاختبار يقدر يمرر نص ويلمّ الناتج في الذاكرة.\n\nلما المدخل سطور، ~Scanner~ بيدي السطر الحالي. parsing الأمر ومعالجة حالته مسؤوليتين مختلفتين.\n\nاستخدم ~interface~ عند مكان الاستهلاك عشان التنفيذ يبقى قابل للاستبدال."),
		ch21Walkthrough("c212.scanner-example", "الـScanner يقرأ كل سطر من reader بالتتابع.", `func main() {
	scanner := bufio.NewScanner(strings.NewReader("check bay"))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}`, "e212-scan-lines", []string{"bufio", "fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "جهّز سطر من الذاكرة."},
			{Line: 3, Caption: "اقرأ السطر الأول."},
			{Line: 4, Caption: "اطبع النص المقروء.", Output: "check bay"},
		}),
	}},
	{ID: "c212.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch21Question("c212.predict-scan", "توقّع اللي هيظهر:\n\n~~~go\nfunc main() {\n\tscanner := bufio.NewScanner(strings.NewReader(\"check bay\"))\n\tfor scanner.Scan() {\n\t\tfmt.Println(scanner.Text())\n\t}\n}\n~~~", []string{"check bay", "scanner.Text", "مفيش output"}, 0, []string{
			"ده نص السطر اللي اتقري.",
			"الـmethod دي مش هي اللي بتتطبع.",
			"الـReader فيه سطر واحد.",
		}),
		ch21Walkthrough("c212.scan-reveal", "شوف السطر وهو بيتحول لناتج.", `func main() {
	scanner := bufio.NewScanner(strings.NewReader("check bay"))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}`, "e212-scan-reveal", []string{"bufio", "fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "النص بقى reader."},
			{Line: 3, Caption: "Scan لقت سطر."},
			{Line: 4, Caption: "اطبع السطر مرة واحدة.", Output: "check bay"},
		}),
		ch21SandboxBlock("c212.playground", "ch21-reader-lines", "قدّم Scanner سطرًا؛ راقب ما وصل للـwriter.", ch21CommandsSandboxConfig()),
	}},
	{ID: "c212.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch21Markdown("c212.watch", "- افحص عدد الكلمات قبل ما تستخدم index."), ch21Markdown("c212.watch-detail", "- خطأ الأمر غير خطأ القراءة أو الكتابة؛ كل واحد له تصرف مختلف.")}},
	{ID: "c212.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch21Markdown("c212.plan", "ابدأ بحدود واضحة: مصدر الأمر، تفسيره، ثم تحديث الحالة.\n\nخلي رسالة النتيجة مكتوبة في writer محقون.\n\nلو command بايظ، هل ده فشل نقل ولا رفض منطق؟"),
		ch21Question("c212.failure-kind", "مين يوقف الجلسة فورًا؟", []string{"أمر غير معروف", "فشل writer", "سطر فاضي"}, 1, []string{
			"رفض الأمر بيتسجل ويتعالج كطلب.",
			"من غير writer مفيش قناة نكمل بيها.",
			"السطر الفاضي مدخل مرفوض.",
		}),
	}},
	{ID: "c212.turn-step", Title: "دورك", Blocks: []LessonBlock{ch21Markdown("c212.turn", "نفّذ الـcontract من starter. جرّب مدخل واحد وسيل سطور، وراجع state بين الأوامر.")}},
}

var ch21CommandsEnglish = []LessonStep{
	{ID: "c212.mission-step", Title: "Mission", Blocks: []LessonBlock{ch21Markdown("c212.mission", "1. Run one command or a stream of lines.\n2. Separate the interface from the order book.\n3. The starter has the program boundary; complete the testable logic.")}},
	{ID: "c212.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch21Markdown("c212.stream", "~io.Reader~ and ~io.Writer~ from 9-4 separate logic from the terminal. A test can supply text and collect output in memory.\n\nFor line input, ~Scanner~ provides the current line. Parsing a command and changing its state are separate responsibilities.\n\nUse an ~interface~ where behavior is consumed so the implementation can be replaced."),
		ch21Walkthrough("c212.scanner-example", "A scanner reads each line from a reader in order.", `func main() {
	scanner := bufio.NewScanner(strings.NewReader("check bay"))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}`, "e212-scan-lines", []string{"bufio", "fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Prepare a line in memory."},
			{Line: 3, Caption: "Read the first line."},
			{Line: 4, Caption: "Print the text read.", Output: "check bay"},
		}),
	}},
	{ID: "c212.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch21Question("c212.predict-scan", "Predict what appears:\n\n~~~go\nfunc main() {\n\tscanner := bufio.NewScanner(strings.NewReader(\"check bay\"))\n\tfor scanner.Scan() {\n\t\tfmt.Println(scanner.Text())\n\t}\n}\n~~~", []string{"check bay", "scanner.Text", "no output"}, 0, []string{
			"That is the text from the line read.",
			"The method itself is not printed.",
			"The reader contains one line.",
		}),
		ch21Walkthrough("c212.scan-reveal", "Watch the line become output.", `func main() {
	scanner := bufio.NewScanner(strings.NewReader("check bay"))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}`, "e212-scan-reveal", []string{"bufio", "fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The text becomes a reader."},
			{Line: 3, Caption: "Scan finds a line."},
			{Line: 4, Caption: "Print the line once.", Output: "check bay"},
		}),
		ch21SandboxBlock("c212.playground", "ch21-reader-lines", "Advance the scanner; watch lines reach the writer.", ch21CommandsSandboxConfig()),
	}},
	{ID: "c212.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch21Markdown("c212.watch", "- Check the word count before using an index."), ch21Markdown("c212.watch-detail", "- A bad command differs from a read or write failure; each needs a different response.")}},
	{ID: "c212.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch21Markdown("c212.plan", "Start with clear boundaries: command source, interpretation, then state change.\n\nKeep result messages in an injected writer.\n\nIs a bad command a transport failure or a rejected request?"),
		ch21Question("c212.failure-kind", "Which condition stops the session immediately?", []string{"Unknown command", "Writer failure", "Empty line"}, 1, []string{
			"A rejected command is recorded and handled as input.",
			"Without a writer, the session has no output path.",
			"An empty line is rejected input.",
		}),
	}},
	{ID: "c212.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch21Markdown("c212.turn", "Follow the starter contract. Try one input and a line stream, then inspect state between commands.")}},
}

var ch21DepthArabic = []LessonStep{
	{ID: "c213.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch21Markdown("c213.mission", "1. وجّه المسبار بتخمينات محدودة.\n2. استقبل الاختيار العشوائي من برّه.\n3. الـstarter محدد reader وwriter؛ ركّز على حدود اللعبة.")}},
	{ID: "c213.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch21Markdown("c213.random-source", "زي ما اتعلمت في 5-5، function value تقدر تبعتها كـdependency. هنا المصدر بياخد حد exclusive ويرجّع index أصغر منه.\n\nاختيار ثابت في الاختبار يخلّي walkthrough قابل للتكرار. في التشغيل الحقيقي، الـadapter يقدر يوصّل بمصدر عشوائي.\n\nخلي الاختيار عند الحد، وخلي قواعد اللعبة تستهلك نتيجة واضحة."),
		ch21Walkthrough("c213.draw-example", "اختيار ثابت يحدد عنصرًا في قائمة صغيرة.", `func pick(labels []string, draw func(int) int) string {
	return labels[draw(len(labels))]
}
func main() {
	fixed := func(n int) int { return 1 }
	fmt.Println(pick([]string{"amber", "cyan", "violet"}, fixed))
}`, "e213-injected-pick", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 5, Caption: "المصدر الثابت يرجع index 1.", Variables: []LessonVariable{{Name: "fixed", Value: "returns 1"}}},
			{Line: 6, Caption: "العنصر الثاني هو cyan.", Output: "cyan"},
		}),
	}},
	{ID: "c213.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch21Question("c213.predict-pick", "توقّع ناتج المثال:\n\n~~~go\nfunc pick(labels []string, draw func(int) int) string {\n\treturn labels[draw(len(labels))]\n}\nfunc main() {\n\tfixed := func(n int) int { return 1 }\n\tfmt.Println(pick([]string{\"amber\", \"cyan\", \"violet\"}, fixed))\n}\n~~~", []string{"amber", "cyan", "violet"}, 1, []string{
			"الـcallback بيرجع index صفر.",
			"الـcallback بيرجع index واحد.",
			"القائمة فيها عنصرين بعد index واحد.",
		}),
		ch21Walkthrough("c213.pick-reveal", "تابع الحد اللي وصل للمصدر والاختيار الناتج.", `func pick(labels []string, draw func(int) int) string {
	return labels[draw(len(labels))]
}
func main() {
	fixed := func(n int) int { return 1 }
	fmt.Println(pick([]string{"amber", "cyan", "violet"}, fixed))
}`, "e213-pick-reveal", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 5, Caption: "المصدر الثابت يرجع 1.", Variables: []LessonVariable{{Name: "fixed", Value: "returns 1"}}},
			{Line: 6, Caption: "اطبع الاسم المختار.", Output: "cyan"},
		}),
		ch21SandboxBlock("c213.playground", "ch21-injected-draw", "بدّل مصدر الاختيار وراقب الـindex والاسم.", ch21DepthSandboxConfig()),
	}},
	{ID: "c213.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch21Markdown("c213.watch", "- افحص نتيجة المصدر قبل استخدامها كـindex."), ch21Markdown("c213.watch-detail", "- متعتمدش على random global جوه منطق تقدر تختبره.")}},
	{ID: "c213.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch21Markdown("c213.plan", "قسّم التفكير: جهّز مصدر الاختيار، اقرأ التخمين، ثم حدّث حالة الجولة.\n\nخلي I/O واضح من خلال reader وwriter.\n\nإيه فايدة callback ثابت في الاختبار؟"),
		ch21Question("c213.fixed-source", "إيه اللي يضيفه مصدر ثابت؟", []string{"نتيجة قابلة للتكرار", "عشوائية أقوى", "حالة مشتركة"}, 0, []string{
			"نفس الاستدعاء يقدر يكرر نفس المسار.",
			"الثبات مش بيحسن العشوائية.",
			"الحقن بيساعد على تجنب الحالة المشتركة.",
		}),
	}},
	{ID: "c213.turn-step", Title: "دورك", Blocks: []LessonBlock{ch21Markdown("c213.turn", "نفّذ contract اللعبة من starter. جرّب reader صغير ومصدر محقون، وراجع نتيجة كل مرحلة.")}},
}

var ch21DepthEnglish = []LessonStep{
	{ID: "c213.mission-step", Title: "Mission", Blocks: []LessonBlock{ch21Markdown("c213.mission", "1. Guide the probe with bounded guesses.\n2. Receive the random choice from outside.\n3. The starter defines a reader and writer; focus on the game boundary.")}},
	{ID: "c213.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch21Markdown("c213.random-source", "From 5-5, a function value can be passed as a dependency. Here the source accepts an exclusive bound and returns a smaller index.\n\nA fixed test source makes a walkthrough repeatable. Real execution can connect an adapter to a random source.\n\nKeep the choice at the boundary and let game rules consume a clear result."),
		ch21Walkthrough("c213.draw-example", "A fixed choice selects an item in a short list.", `func pick(labels []string, draw func(int) int) string {
	return labels[draw(len(labels))]
}
func main() {
	fixed := func(n int) int { return 1 }
	fmt.Println(pick([]string{"amber", "cyan", "violet"}, fixed))
}`, "e213-injected-pick", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 5, Caption: "The fixed source returns index 1.", Variables: []LessonVariable{{Name: "fixed", Value: "returns 1"}}},
			{Line: 6, Caption: "The second item is cyan.", Output: "cyan"},
		}),
	}},
	{ID: "c213.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch21Question("c213.predict-pick", "Predict the example output:\n\n~~~go\nfunc pick(labels []string, draw func(int) int) string {\n\treturn labels[draw(len(labels))]\n}\nfunc main() {\n\tfixed := func(n int) int { return 1 }\n\tfmt.Println(pick([]string{\"amber\", \"cyan\", \"violet\"}, fixed))\n}\n~~~", []string{"amber", "cyan", "violet"}, 1, []string{
			"The callback returns index zero.",
			"The callback returns index one.",
			"The list has two items after index one.",
		}),
		ch21Walkthrough("c213.pick-reveal", "Follow the bound given to the source and the chosen item.", `func pick(labels []string, draw func(int) int) string {
	return labels[draw(len(labels))]
}
func main() {
	fixed := func(n int) int { return 1 }
	fmt.Println(pick([]string{"amber", "cyan", "violet"}, fixed))
}`, "e213-pick-reveal", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 5, Caption: "The fixed source returns 1.", Variables: []LessonVariable{{Name: "fixed", Value: "returns 1"}}},
			{Line: 6, Caption: "Print the chosen name.", Output: "cyan"},
		}),
		ch21SandboxBlock("c213.playground", "ch21-injected-draw", "Switch choice sources; watch the index and label.", ch21DepthSandboxConfig()),
	}},
	{ID: "c213.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch21Markdown("c213.watch", "- Validate the source result before using it as an index."), ch21Markdown("c213.watch-detail", "- Avoid global randomness inside logic you need to test.")}},
	{ID: "c213.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch21Markdown("c213.plan", "Split the work: prepare a choice source, read a guess, then update round state.\n\nKeep I/O visible through the reader and writer.\n\nWhy is a fixed callback useful in a test?"),
		ch21Question("c213.fixed-source", "What does a fixed source provide?", []string{"Repeatable results", "Stronger randomness", "Shared state"}, 0, []string{
			"The same call can repeat the same path.",
			"Fixed behavior does not improve randomness.",
			"Injection helps avoid shared state.",
		}),
	}},
	{ID: "c213.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch21Markdown("c213.turn", "Implement the game contract from the starter. Try a small reader and injected source, then inspect each phase's result.")}},
}

var ch21DuelArabic = []LessonStep{
	{ID: "c214.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch21Markdown("c214.mission", "1. حوّل أوامر النص لقيم مناورة مسماة.\n2. احسب الجولة بمصدر عشوائية محقون.\n3. الـstarter فيه أنواع ناقصة؛ كمّل قواعد سهلة المراجعة.")}},
	{ID: "c214.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch21Markdown("c214.named-values", "الـnamed type و~iota~ من 2-4 بيرتبوا قيم ليها معنى بدل strings متفرقة. مجموعة constants بتعمل قاموس صغير.\n\nحوّل النص لقيمة بعد التحقق عند الحد. جوه القاعدة، استخدم النوع المسمى وخلي input بعيد عن حساب النتيجة.\n\nالفصل ده يسهل مراجعة الحالات واحدة واحدة."),
		ch21Walkthrough("c214.enum-example", "~iota~ بترقم قيم النوع بالترتيب.", `type LampMode int
const (
	idle LampMode = iota
	scan
	relay
)
func main() {
	mode := scan
	fmt.Println(int(mode))
}`, "e214-iota-values", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "ابدأ مجموعة constants."},
			{Line: 4, Caption: "القيمة التانية رقمها 1."},
			{Line: 8, Caption: "mode تحفظ scan."},
			{Line: 9, Caption: "اطبع القيمة الرقمية.", Output: "1"},
		}),
	}},
	{ID: "c214.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch21Question("c214.predict-mode", "توقّع الناتج:\n\n~~~go\ntype LampMode int\nconst (\n\tidle LampMode = iota\n\tscan\n\trelay\n)\nfunc main() {\n\tmode := scan\n\tfmt.Println(int(mode))\n}\n~~~", []string{"0", "1", "2"}, 1, []string{
			"idle هي أول قيمة ورقمها صفر.",
			"scan تيجي بعد idle مباشرة.",
			"relay هي القيمة التالتة.",
		}),
		ch21Walkthrough("c214.mode-reveal", "امشِ من أول constant لحد الطباعة.", `type LampMode int
const (
	idle LampMode = iota
	scan
	relay
)
func main() {
	mode := scan
	fmt.Println(int(mode))
}`, "e214-mode-reveal", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الثوابت بتبدأ من صفر."},
			{Line: 4, Caption: "scan قيمتها واحد."},
			{Line: 8, Caption: "اختار القيمة المسماة."},
			{Line: 9, Caption: "الناتج الرقمي هو 1.", Output: "1"},
		}),
		ch21SandboxBlock("c214.playground", "ch21-maneuver-values", "مرّر أوامر توضيحية وشوف القيمة المسماة.", ch21DuelSandboxConfig()),
	}},
	{ID: "c214.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch21Markdown("c214.watch", "- متحوّلش نص غير موثوق للنوع قبل ما تتحقق منه."), ch21Markdown("c214.watch-detail", "- parsing وقرار الجولة مرحلتين؛ خليك واضح في الحد بينهم.")}},
	{ID: "c214.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch21Markdown("c214.plan", "امشِ في خطين: طبّع المدخل وتحقق منه، وبعدها طبّق قاعدة الجولة.\n\nفكّر في كل زوج قيم بالترتيب، مش في كل قيمة لوحدها.\n\nلو عكست اللاعب والمنافس، إيه اللي لازم تراجعه؟"),
		ch21Question("c214.parse-boundary", "فين يتحول نص الأمر لنوع مسمى؟", []string{"عند حد الإدخال", "داخل الطباعة", "بعد تحديث النتيجة"}, 0, []string{
			"تحقق من النص قبل وصوله لقواعد اللعب.",
			"الطباعة تعرض النتيجة فقط.",
			"الحالة تتحدث بعد القرار.",
		}),
	}},
	{ID: "c214.turn-step", Title: "دورك", Blocks: []LessonBlock{ch21Markdown("c214.turn", "ابدأ من النوع المسمى في starter. افصل parsing عن القرار، وبعدها راجع الجولات المكتوبة.")}},
}

var ch21DuelEnglish = []LessonStep{
	{ID: "c214.mission-step", Title: "Mission", Blocks: []LessonBlock{ch21Markdown("c214.mission", "1. Turn text commands into named maneuver values.\n2. Score a round using injected randomness.\n3. The starter has unfinished types; complete clear, reviewable rules.")}},
	{ID: "c214.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch21Markdown("c214.named-values", "Named types and ~iota~ from 2-4 organize values with meaning beyond scattered strings. Constants form a small vocabulary.\n\nConvert text after validation at the boundary. Inside the rule, use the named type and keep input apart from outcome logic.\n\nThis separation makes cases easier to review one by one."),
		ch21Walkthrough("c214.enum-example", "~iota~ numbers the type's values in order.", `type LampMode int
const (
	idle LampMode = iota
	scan
	relay
)
func main() {
	mode := scan
	fmt.Println(int(mode))
}`, "e214-iota-values", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Start a constant group."},
			{Line: 4, Caption: "The second value is 1."},
			{Line: 8, Caption: "mode holds scan."},
			{Line: 9, Caption: "Print the numeric value.", Output: "1"},
		}),
	}},
	{ID: "c214.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch21Question("c214.predict-mode", "Predict the output:\n\n~~~go\ntype LampMode int\nconst (\n\tidle LampMode = iota\n\tscan\n\trelay\n)\nfunc main() {\n\tmode := scan\n\tfmt.Println(int(mode))\n}\n~~~", []string{"0", "1", "2"}, 1, []string{
			"idle is the first value, numbered zero.",
			"scan follows idle directly.",
			"relay is the third value.",
		}),
		ch21Walkthrough("c214.mode-reveal", "Follow the constants from their start to the print.", `type LampMode int
const (
	idle LampMode = iota
	scan
	relay
)
func main() {
	mode := scan
	fmt.Println(int(mode))
}`, "e214-mode-reveal", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Constants start at zero."},
			{Line: 4, Caption: "scan has value one."},
			{Line: 8, Caption: "Choose the named value."},
			{Line: 9, Caption: "The numeric output is 1.", Output: "1"},
		}),
		ch21SandboxBlock("c214.playground", "ch21-maneuver-values", "Cycle example commands and inspect each named value.", ch21DuelSandboxConfig()),
	}},
	{ID: "c214.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch21Markdown("c214.watch", "- Do not convert untrusted text before validating it."), ch21Markdown("c214.watch-detail", "- Parsing and deciding a round are separate stages; keep their boundary clear.")}},
	{ID: "c214.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch21Markdown("c214.plan", "Use two passes: normalize and validate input, then apply the round rule.\n\nThink through each ordered pair, not each value alone.\n\nIf player and rival swap places, what should you review?"),
		ch21Question("c214.parse-boundary", "Where should command text become a named type?", []string{"At the input boundary", "Inside printing", "After updating the score"}, 0, []string{
			"Validate text before it reaches game rules.",
			"Printing only shows the result.",
			"State changes after the decision.",
		}),
	}},
	{ID: "c214.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch21Markdown("c214.turn", "Start with the named type in the starter. Separate parsing from deciding, then inspect the rounds written.")}},
}
