package levels

func ch18Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch18Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch18Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch18Localized(ar, en string) map[string]any {
	return map[string]any{"ar": ar, "en": en}
}

func ch18Sandbox(id, kind, ar, en, mode string, english bool) LessonBlock {
	if english {
		ar = en
	}
	return LessonBlock{
		ID: id, Kind: "sandbox", Text: ar,
		Sandbox: &LessonSandbox{Kind: kind, Config: ch18SandboxConfig(mode)},
	}
}

func ch18SandboxConfig(mode string) map[string]any {
	labels := map[string]any{}
	data := map[string]any{}
	switch mode {
	case "index":
		labels = map[string]any{
			"values": ch18Localized("القيم", "Values"),
			"read":   ch18Localized("القراءة", "Read"),
			"action": ch18Localized("index التالي", "Next index"),
			"caption": []map[string]any{
				ch18Localized("ابدأ من index صفر", "Start at index zero"),
				ch18Localized("كل index له قيمة", "Each index has a value"),
				ch18Localized("الترتيب يحدد القراءة", "Order determines the read"),
			},
		}
		data = map[string]any{"values": []int{7, 4, 9}}
	case "json":
		labels = map[string]any{
			"source": ch18Localized("قيمة Go", "Go value"),
			"result": ch18Localized("JSON bytes", "JSON bytes"),
			"action": ch18Localized("Marshal", "Marshal"),
			"before": ch18Localized("جاهزة للتحويل", "Ready to encode"),
			"after":  ch18Localized("بايتات JSON", "JSON bytes"),
		}
		data = map[string]any{"source": "Beacon{Charge: 7}", "encoded": `{"charge":7}`}
	case "csv":
		labels = map[string]any{
			"source": ch18Localized("سطر CSV", "CSV row"),
			"fields": ch18Localized("الحقول", "Fields"),
			"action": ch18Localized("السطر التالي", "Next row"),
			"captions": []map[string]any{
				ch18Localized("حقلان منفصلان", "Two separate fields"),
				ch18Localized("الفاصلة جوه quotes", "Comma stays inside quotes"),
			},
		}
		data = map[string]any{"rows": []map[string]any{
			{"source": ch18Localized("شمال,4", "north,4"), "fields": []map[string]any{ch18Localized("شمال", "north"), ch18Localized("4", "4")}},
			{"source": ch18Localized(`"غرب، علوي",4`, `"west, upper",4`), "fields": []map[string]any{ch18Localized("غرب، علوي", "west, upper"), ch18Localized("4", "4")}},
		}}
	case "unicode":
		labels = map[string]any{
			"source":       ch18Localized("النص", "Text"),
			"result":       ch18Localized("الكلمات", "Words"),
			"split":        ch18Localized("FieldsFunc", "FieldsFunc"),
			"lower":        ch18Localized("ToLower", "ToLower"),
			"splitCaption": ch18Localized("الفواصل تفصل الكلمات", "Separators split words"),
			"rawCaption":   ch18Localized("النص الأصلي", "Original text"),
			"lowerCaption": ch18Localized("الحروف بقت صغيرة", "Letters are lowercase"),
		}
		data = map[string]any{
			"raw":        ch18Localized("إيكو—Relay 7", "Echo—Relay 7"),
			"lowerRaw":   ch18Localized("إيكو—relay 7", "echo—relay 7"),
			"words":      []map[string]any{ch18Localized("إيكو", "Echo"), ch18Localized("Relay", "Relay"), ch18Localized("7", "7")},
			"lowerWords": []map[string]any{ch18Localized("إيكو", "echo"), ch18Localized("relay", "relay"), ch18Localized("7", "7")},
		}
	case "splitn":
		labels = map[string]any{
			"source":  ch18Localized("النص", "Text"),
			"limit":   ch18Localized("عدد الأجزاء", "Piece limit"),
			"parts":   ch18Localized("الناتج", "Parts"),
			"action":  ch18Localized("زوّد الحد", "Raise limit"),
			"caption": ch18Localized("آخر جزء يحتفظ بالباقي", "Last part keeps the rest"),
		}
		data = map[string]any{"raw": "9|relay|warm|again", "separator": "|", "limits": []int{2, 3}}
	case "search":
		labels = map[string]any{
			"name":    ch18Localized("الاسم", "Name"),
			"query":   ch18Localized("البحث", "Search"),
			"result":  ch18Localized("النتيجة", "Result"),
			"pending": ch18Localized("لسه ما اتبحثش", "Not searched yet"),
			"action":  ch18Localized("Normalize + Contains", "Normalize + Contains"),
			"caption": ch18Localized("وحّد الحالة قبل البحث", "Normalize case before searching"),
		}
		data = map[string]any{"name": "Pilot", "normalized": "pilot", "query": "lot", "match": "true"}
	}
	return map[string]any{"mode": mode, "labels": labels, "data": data}
}

var ch18MissionLogArabic = []LessonStep{
	{ID: "c181.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch18Markdown("c181.mission", "الكود في ~main.go~ جاهز؛ شغلك في ~main_test.go~.\n\nاكتب table-driven tests تثبت إن سجل الطلعات بيحافظ على هويتها وترتيبها.\n\nكل حالة تبقى مستقلة، والـtest يثبت السلوك مش عدد العناصر بس.")}},
	{ID: "c181.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch18Walkthrough("c181.idea-table", "الـtable بتجمع بيانات حالات مختلفة؛ loop واحدة تشغّل نفس الفحص على كل حالة.", ch18MissionLogTableCode, "e18-1-table", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "ابدأ بسجل حالات."},
		{Line: 5, Caption: "الحالتان جاهزتان."},
		{Line: 6, Caption: "اقرا الحالة الأولى.", Variables: []LessonVariable{{Name: "tc", Value: "{small 2}"}}},
		{Line: 7, Caption: "اطبع الـvalue بتاعتها.", Variables: []LessonVariable{{Name: "tc", Value: "{small 2}"}}, Output: "small 2"},
		{Line: 7, Caption: "نفس الفحص للحالة التالية.", Variables: []LessonVariable{{Name: "tc", Value: "{large 4}"}}, Output: "small 2\nlarge 4"},
		{Line: 8, Caption: "خلصت كل الحالات.", Output: "small 2\nlarge 4"},
	})}},
	{ID: "c181.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch18Question("c181.try-index", "توقّع الناتج:\n\n~~~go\n"+ch18LogIndexCode+"\n~~~", []string{"7", "4", "2"}, 1, []string{"الـindex صفر فيه 7.", "الـindex واحد فيه 4.", "الطول اتنين، مش الـvalue بتاعة العنصر."}),
		ch18Walkthrough("c181.try-reveal", "تابع الـslice وبعدها قراءة الـindex.", ch18LogIndexCode, "e18-1-index", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "اتنين values بالترتيب.", Variables: []LessonVariable{{Name: "values", Value: "[7 4]"}}},
			{Line: 3, Caption: "اقرا الـindex 1.", Variables: []LessonVariable{{Name: "values", Value: "[7 4]"}}, Output: "4"},
		}),
		ch18Sandbox("c181.sandbox", "ch18-slice-index", "غيّر index.", "Move the index.", "index", false),
	}},
	{ID: "c181.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch18Markdown("c181.watch", "حالات مستقلة بتمنع الـtest يأثر على اللي بعده.\n\nقارن الـvalues والـerrors والترتيب؛ الطول لوحده مش كفاية.")}},
	{ID: "c181.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch18Markdown("c181.plan", "حوّل كل سلوك في العقد لنتيجة تقدر تلاحظها.\n\nخلي البيانات تختلف، وخلي خطوات الفحص مشتركة."),
		ch18Question("c181.check-state", "إزاي تثبت إن العملية الفاشلة ما غيّرتش السجل؟", []string{"قارن الحالة قبلها وبعدها", "افحص عدد العناصر بس", "افحص نص الخطأ بس"}, 0, []string{"المقارنة بتكشف أي تعديل في البيانات.", "نفس العدد ممكن يخفي تغيير المحتوى.", "نص الـerror مش دليل على ثبات السجل."}),
	}},
	{ID: "c181.turn-step", Title: "دورك", Blocks: []LessonBlock{ch18Markdown("c181.turn", "ابدأ بحالة صغيرة، شغّل الـtests، وبعدها زوّد حالة تكشف سلوك مختلف.")}},
}

var ch18MissionLogEnglish = []LessonStep{
	{ID: "c181.mission-step", Title: "Mission", Blocks: []LessonBlock{ch18Markdown("c181.mission", "The code in ~main.go~ is ready; work in ~main_test.go~.\n\nWrite table-driven tests showing that the flight log preserves identity and order.\n\nKeep cases independent, and assert behavior rather than only item counts.")}},
	{ID: "c181.idea-step", Title: "Idea", Blocks: []LessonBlock{ch18Walkthrough("c181.idea-table", "A table groups different case data; one loop runs the same check for each case.", ch18MissionLogTableCode, "e18-1-table", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Start a case table."},
		{Line: 5, Caption: "Both cases are ready."},
		{Line: 6, Caption: "Read the first case.", Variables: []LessonVariable{{Name: "tc", Value: "{small 2}"}}},
		{Line: 7, Caption: "Print its values.", Variables: []LessonVariable{{Name: "tc", Value: "{small 2}"}}, Output: "small 2"},
		{Line: 7, Caption: "Run the same check again.", Variables: []LessonVariable{{Name: "tc", Value: "{large 4}"}}, Output: "small 2\nlarge 4"},
		{Line: 8, Caption: "Every case has run.", Output: "small 2\nlarge 4"},
	})}},
	{ID: "c181.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch18Question("c181.try-index", "Predict the output:\n\n~~~go\n"+ch18LogIndexCode+"\n~~~", []string{"7", "4", "2"}, 1, []string{"Index zero holds 7.", "Index one holds 4.", "The length is two, not an element value."}),
		ch18Walkthrough("c181.try-reveal", "Follow the slice, then read its indexed value.", ch18LogIndexCode, "e18-1-index", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Two ordered values.", Variables: []LessonVariable{{Name: "values", Value: "[7 4]"}}},
			{Line: 3, Caption: "Read index one.", Variables: []LessonVariable{{Name: "values", Value: "[7 4]"}}, Output: "4"},
		}),
		ch18Sandbox("c181.sandbox", "ch18-slice-index", "غيّر index.", "Move the index.", "index", true),
	}},
	{ID: "c181.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch18Markdown("c181.watch", "Independent cases keep one test from affecting the next.\n\nCompare values, errors, and order; length alone is not enough.")}},
	{ID: "c181.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch18Markdown("c181.plan", "Turn each contract behavior into an observable result.\n\nVary the data while keeping the check shared."),
		ch18Question("c181.check-state", "How can you show that a failed operation left the log unchanged?", []string{"Compare state before and after", "Check only the item count", "Check only the error text"}, 0, []string{"A comparison reveals changes to the data.", "The same count can hide changed contents.", "An error message does not prove state stayed fixed."}),
	}},
	{ID: "c181.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch18Markdown("c181.turn", "Start with one small case, run the tests, then add a case that exposes different behavior.")}},
}

var ch18BlackBoxArabic = []LessonStep{
	{ID: "c182.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch18Markdown("c182.mission", "خزّن بيانات الصندوق واقراها بصيغة JSON.\n\nخلي الـerrors اللي من الملف أو المحتوى واضحة، وما تمسحش دليل مشكلة حصلت قبل كده.\n\nالـtag بتوصل اسم Go بـkey JSON. الاسم في Go بيفضل زي ما هو.")}},
	{ID: "c182.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch18Walkthrough("c182.idea-json", "~Marshal~ بتحوّل Go value لـJSON bytes؛ أسماء الحقول على السلك مش لازم تبقى نفس الأسماء في Go.", ch18BlackBoxJSONCode, "e18-2-json", []string{"encoding/json", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 1, Caption: "نوع ببيانات مصدّرة."},
		{Line: 5, Caption: "الـvalue جاهزة للتحويل.", Variables: []LessonVariable{{Name: "b", Value: "{7}"}}},
		{Line: 6, Caption: "حوّلها لصيغة JSON."},
		{Line: 7, Caption: "اطبع الـbytes والناتج.", Output: "{\"charge\":7} <nil>"},
	})}},
	{ID: "c182.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch18Question("c182.try-json", "توقّع الناتج:\n\n~~~go\n"+ch18BlackBoxMapCode+"\n~~~", []string{"{\"hull\":3} <nil>", "hull=3", "خطأ"}, 0, []string{"Marshal بترجّع JSON والـerror هنا nil.", "المكتبة بترجّع JSON، مش صيغة key=value.", "الـvalue دي صالحة للتحويل."}),
		ch18Walkthrough("c182.try-reveal", "شوف إزاي map واحدة بتطلع object في JSON.", ch18BlackBoxMapCode, "e18-2-map", []string{"encoding/json", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "Marshal بترجّع bytes."},
			{Line: 3, Caption: "الناتج object والـerror nil.", Output: "{\"hull\":3} <nil>"},
		}),
		ch18Sandbox("c182.sandbox", "ch18-json-marshal", "جرّب Marshal.", "Try Marshal.", "json", false),
	}},
	{ID: "c182.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch18Markdown("c182.watch", "~Unmarshal~ محتاجة pointer عشان تملا الـvalue.\n\nفرّق بين JSON سليمة ومعناها، وبص على كل error قبل ما تستخدم الناتج.")}},
	{ID: "c182.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch18Markdown("c182.plan", "افصل التحقق عن التحويل والملف.\n\nراجع الحالة قبل ما تحفظ، وخلي القراية ما تغيّرش المصدر."),
		ch18Question("c182.check-slice", "إيه الفرق بين slice ناقصة وواحدة موجودة وفاضية؟", []string{"الأولى nil والتانية مش nil", "الاتنين nil", "الطول مختلف"}, 0, []string{"JSON null والغياب ممكن ينتجوا nil.", "الـslice الفاضية المهيأة مش nil.", "الاتنين طولهم صفر."}),
	}},
	{ID: "c182.turn-step", Title: "دورك", Blocks: []LessonBlock{ch18Markdown("c182.turn", "ابدأ بالتحويل والتحقق، وبعدها جرّب مسارات القراءة والكتابة كل واحد لوحده.")}},
}

var ch18BlackBoxEnglish = []LessonStep{
	{ID: "c182.mission-step", Title: "Mission", Blocks: []LessonBlock{ch18Markdown("c182.mission", "Implement saving and reading the black-box data as JSON.\n\nMake file and content failures clear, and preserve evidence from earlier failures.\n\nA tag maps a Go field name to a JSON key.")}},
	{ID: "c182.idea-step", Title: "Idea", Blocks: []LessonBlock{ch18Walkthrough("c182.idea-json", "~Marshal~ converts a Go value to JSON bytes; wire field names need not match Go names.", ch18BlackBoxJSONCode, "e18-2-json", []string{"encoding/json", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "A type with an exported field."},
		{Line: 5, Caption: "The value is ready.", Variables: []LessonVariable{{Name: "b", Value: "{7}"}}},
		{Line: 6, Caption: "Convert it to JSON."},
		{Line: 7, Caption: "Show bytes and error.", Output: "{\"charge\":7} <nil>"},
	})}},
	{ID: "c182.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch18Question("c182.try-json", "Predict the output:\n\n~~~go\n"+ch18BlackBoxMapCode+"\n~~~", []string{"{\"hull\":3} <nil>", "hull=3", "an error"}, 0, []string{"Marshal returns JSON and a nil error here.", "The library returns JSON, not key=value text.", "This value is valid to encode."}),
		ch18Walkthrough("c182.try-reveal", "See one map become a JSON object.", ch18BlackBoxMapCode, "e18-2-map", []string{"encoding/json", "fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Marshal returns bytes."},
			{Line: 3, Caption: "Object output and nil error.", Output: "{\"hull\":3} <nil>"},
		}),
		ch18Sandbox("c182.sandbox", "ch18-json-marshal", "جرّب Marshal.", "Try Marshal.", "json", true),
	}},
	{ID: "c182.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch18Markdown("c182.watch", "~Unmarshal~ needs a pointer to fill the value.\n\nSeparate valid JSON from valid meaning, and check each error before using output.")}},
	{ID: "c182.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch18Markdown("c182.plan", "Separate validation from conversion and file access.\n\nReview state before saving; reading should not change its source."),
		ch18Question("c182.check-slice", "What differs between a missing slice and a present empty one?", []string{"The first is nil; the second is not", "Both are nil", "Their lengths differ"}, 0, []string{"JSON null and an absent field can produce nil.", "An initialized empty slice is not nil.", "Both have length zero."}),
	}},
	{ID: "c182.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch18Markdown("c182.turn", "Start with conversion and validation, then try file reads and writes separately.")}},
}

var ch18CargoArabic = []LessonStep{
	{ID: "c183.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch18Markdown("c183.mission", "اقرا manifest CSV من ~io.Reader~ واجمع بياناته.\n\nاستعمل parser يحافظ على الحقول المقتبسة؛ أي مشكلة في المحتوى أو المصدر تمنع تقرير ناقص.")}},
	{ID: "c183.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch18Walkthrough("c183.idea-csv", "كل ~Read~ بترجع record كاملة؛ الـCSV parser بيتولى فك الحقول.", ch18CargoCSVCode, "e18-3-csv", []string{"encoding/csv", "fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "اعمل reader من نص."},
		{Line: 3, Caption: "اقرا record واحدة."},
		{Line: 4, Caption: "اطبع الحقول والـerror.", Output: "[north 4] <nil>"},
	})}},
	{ID: "c183.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch18Question("c183.try-csv", "توقّع الناتج:\n\n~~~go\n"+ch18CargoQuotedCode+"\n~~~", []string{"[west upper 4] <nil>", "[west, upper 4] <nil>", "خطأ parsing"}, 1, []string{"الـcomma جوه quotes ما بتقسمش الحقل.", "الـparser حافظ على الـcomma داخل الاسم.", "السطر صالح بصيغة CSV."}),
		ch18Walkthrough("c183.try-reveal", "تابع قراءة CSV فيها فاصلة مقتبسة.", ch18CargoQuotedCode, "e18-3-quoted", []string{"encoding/csv", "fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "الـreader يستقبل نص CSV."},
			{Line: 3, Caption: "Read تفك الحقول."},
			{Line: 4, Caption: "الاسم يفضل حقل واحد.", Output: "[west, upper 4] <nil>"},
		}),
		ch18Sandbox("c183.sandbox", "ch18-csv-quoted-field", "بدّل الصف.", "Switch rows.", "csv", false),
	}},
	{ID: "c183.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch18Markdown("c183.watch", "~io.EOF~ نهاية طبيعية بعد الـheader؛ أي error تاني يوقف الملف.\n\nعدّ records، مش السطور الظاهرة، وخلي سبب الـerrors اللي بتحصل وقت القراءة.")}},
	{ID: "c183.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch18Markdown("c183.plan", "اقرا schema الأول، وبعدها عالج record record.\n\nالـheader بيحدد شكل البيانات قبل صفوف الشحنة. حدّث إجمالي محلي؛ ما ترجعش نجاح جزئي بعد فشل."),
		ch18Question("c183.check-eof", "إمتى ~io.EOF~ معناها إن الملف خلص طبيعي؟", []string{"بعد قراءة الـheader والـrecords", "قبل وجود header", "كل ما يظهر record غلط"}, 0, []string{"الـEOF بعد البيانات علامة نهاية القراءة.", "قبل الـheader مفيش schema صالحة.", "record غلط بيرجع error لوحده."}),
	}},
	{ID: "c183.turn-step", Title: "دورك", Blocks: []LessonBlock{ch18Markdown("c183.turn", "ابدأ واقرا header، وبعدها جرّب record سليمة وsource يطلع error.")}},
}

var ch18CargoEnglish = []LessonStep{
	{ID: "c183.mission-step", Title: "Mission", Blocks: []LessonBlock{ch18Markdown("c183.mission", "Read a CSV manifest from ~io.Reader~ and aggregate its data.\n\nUse a parser that preserves quoted fields; any content or source failure prevents a partial report.")}},
	{ID: "c183.idea-step", Title: "Idea", Blocks: []LessonBlock{ch18Walkthrough("c183.idea-csv", "Each ~Read~ returns a complete record; the CSV parser handles field decoding.", ch18CargoCSVCode, "e18-3-csv", []string{"encoding/csv", "fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create a reader from text."},
		{Line: 3, Caption: "Read one record."},
		{Line: 4, Caption: "Print fields and error.", Output: "[north 4] <nil>"},
	})}},
	{ID: "c183.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch18Question("c183.try-csv", "Predict the output:\n\n~~~go\n"+ch18CargoQuotedCode+"\n~~~", []string{"[west upper 4] <nil>", "[west, upper 4] <nil>", "a parse error"}, 1, []string{"A comma inside quotes does not split a field.", "The parser keeps the comma inside the name.", "This record is valid CSV."}),
		ch18Walkthrough("c183.try-reveal", "Follow one CSV read with a quoted comma.", ch18CargoQuotedCode, "e18-3-quoted", []string{"encoding/csv", "fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The reader receives CSV text."},
			{Line: 3, Caption: "Read decodes the fields."},
			{Line: 4, Caption: "The name stays one field.", Output: "[west, upper 4] <nil>"},
		}),
		ch18Sandbox("c183.sandbox", "ch18-csv-quoted-field", "بدّل الصف.", "Switch rows.", "csv", true),
	}},
	{ID: "c183.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch18Markdown("c183.watch", "~io.EOF~ is normal after the header; other errors stop the file.\n\nCount records, not visible lines, and preserve read error causes.")}},
	{ID: "c183.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch18Markdown("c183.plan", "Read the schema first, then process records one by one.\n\nThe header describes fields before cargo rows. Update local totals; do not return partial success after failure."),
		ch18Question("c183.check-eof", "When does ~io.EOF~ mean the file ended normally?", []string{"After its header and records", "Before any header", "Whenever a record is invalid"}, 0, []string{"EOF after the data marks the end of reading.", "Without a header there is no valid schema.", "An invalid record returns a different error."}),
	}},
	{ID: "c183.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch18Markdown("c183.turn", "Start by reading the header, then try a valid record and a source that fails.")}},
}

var ch18SignalArabic = []LessonStep{
	{ID: "c184.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch18Markdown("c184.mission", "حوّل stream الإشارات لقائمة كلمات مرتبة حسب التكرار.\n\nالتعامل مع Unicode وثبات ترتيب التعادل جزء من النتيجة؛ ما تخزّنش كل النص لمجرد العد.\n\nالـpunctuation والـemoji فواصل بين الكلمات.\n\nوحّد حالة الحروف قبل العد؛ Echo وECHO يبقوا نفس الكلمة.\n\nالـmap ما بتضمنش ترتيب، فحدد قاعدة واضحة للتعادل.")}},
	{ID: "c184.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch18Walkthrough("c184.idea-words", "~FieldsFunc~ بتقسّم حسب شرط على rune؛ وبعدها ~ToLower~ توحّد شكل الكلمة.", ch18SignalUnicodeCode, "e18-4-unicode", []string{"fmt", "strings", "unicode"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "جزّئ حسب حروف Unicode."},
		{Line: 2, Caption: "الكلمات اتجمعت."},
		{Line: 3, Caption: "حوّل أول كلمة واطبع العدد.", Output: "echo 3"},
	})}},
	{ID: "c184.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch18Question("c184.try-lower", "توقّع الناتج:\n\n~~~go\n"+ch18SignalLowerCode+"\n~~~", []string{"Echo", "echo", "ECHO"}, 1, []string{"الحرف الكبير بيتحوّل لصغير.", "ToLower بتوحّد الحروف لحالة صغيرة.", "الـfunction ما بتحوّلش لحروف كبيرة."}),
		ch18Walkthrough("c184.try-reveal", "شوف تطبيق ToLower على كلمة للتوضيح.", ch18SignalLowerCode, "e18-4-lower", []string{"fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "حوّل الكلمة لحروف صغيرة."},
			{Line: 3, Caption: "اطبع الـvalue الجديدة.", Output: "echo"},
		}),
		ch18Sandbox("c184.sandbox", "ch18-unicode-words", "افصل وحوّل.", "Split and lowercase.", "unicode", false),
	}},
	{ID: "c184.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch18Markdown("c184.watch", "استعمل ~unicode.IsLetter~ و~IsDigit~ على rune، مش byte.\n\n~Scanner.Buffer~ لازم تتظبط قبل أول ~Scan~؛ بص على ~Err~ بعد الـloop.")}},
	{ID: "c184.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch18Markdown("c184.plan", "عدّ الكلمات في map، وبعدها حوّل النتائج لـslice عشان ترتّبها.\n\nالتعادل محتاج معيار ثابت غير العدد."),
		ch18Question("c184.check-order", "هل ترتيب ~range~ على map ثابت؟", []string{"مش مضمون", "حسب إدخال المفاتيح", "أبجدي دايماً"}, 0, []string{"ترتيب map ما ينفعش تعتمد عليه.", "Go ما بتحفظش ترتيب إدخال map.", "الترتيب الأبجدي محتاج sorting."}),
	}},
	{ID: "c184.turn-step", Title: "دورك", Blocks: []LessonBlock{ch18Markdown("c184.turn", "ابدأ بتقسيم سطر واحد، شغّل الكود، وبعدها جرّب stream وراجع ترتيب النتيجة.")}},
}

var ch18SignalEnglish = []LessonStep{
	{ID: "c184.mission-step", Title: "Mission", Blocks: []LessonBlock{ch18Markdown("c184.mission", "Turn a signal stream into words ranked by frequency.\n\nUnicode handling and stable tie order affect the result; do not retain all text just to count it.\n\nPunctuation and emoji separate words.\n\nNormalize case before counting; Echo and ECHO become one word.\n\nMaps do not guarantee order, so define a clear tie rule.")}},
	{ID: "c184.idea-step", Title: "Idea", Blocks: []LessonBlock{ch18Walkthrough("c184.idea-words", "~FieldsFunc~ splits using a predicate on each rune; ~ToLower~ then normalizes the word.", ch18SignalUnicodeCode, "e18-4-unicode", []string{"fmt", "strings", "unicode"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Split using Unicode letters."},
		{Line: 2, Caption: "The words are collected."},
		{Line: 3, Caption: "Lowercase and print the count.", Output: "echo 3"},
	})}},
	{ID: "c184.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch18Question("c184.try-lower", "Predict the output:\n\n~~~go\n"+ch18SignalLowerCode+"\n~~~", []string{"Echo", "echo", "ECHO"}, 1, []string{"The uppercase letter becomes lowercase.", "ToLower normalizes letters to lowercase.", "This function does not uppercase text."}),
		ch18Walkthrough("c184.try-reveal", "Watch ToLower on an illustrative word.", ch18SignalLowerCode, "e18-4-lower", []string{"fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Convert the word to lowercase."},
			{Line: 3, Caption: "Print the new value.", Output: "echo"},
		}),
		ch18Sandbox("c184.sandbox", "ch18-unicode-words", "افصل وحوّل.", "Split and lowercase.", "unicode", true),
	}},
	{ID: "c184.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch18Markdown("c184.watch", "Use ~unicode.IsLetter~ and ~IsDigit~ on runes, not bytes.\n\nSet ~Scanner.Buffer~ before the first ~Scan~; check ~Err~ after the loop.")}},
	{ID: "c184.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch18Markdown("c184.plan", "Count words in a map, then collect results into a slice to sort.\n\nTies need a stable rule beyond frequency."),
		ch18Question("c184.check-order", "Is ~range~ order over a map stable?", []string{"No guarantee", "It follows insertion order", "It is always alphabetical"}, 0, []string{"You cannot rely on map iteration order.", "Go maps do not preserve insertion order.", "Alphabetical order requires sorting."}),
	}},
	{ID: "c184.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch18Markdown("c184.turn", "Start by splitting one line, run it, then try the stream and review result ordering.")}},
}

var ch18RecorderArabic = []LessonStep{
	{ID: "c185.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch18Markdown("c185.mission", "حلّل records من stream وسجّل أرقام اللي مش سليمة.\n\nفرّق بين record بايظة وsource فشل أثناء القراءة؛ التانية تمنع summary ناقصة.\n\nرقم السطر الفعلي غير sequence المكتوبة جواه. خلي التشخيص يشاور للسطر اللي اتقري.")}},
	{ID: "c185.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch18Walkthrough("c185.idea-split", "~SplitN~ بتحد عدد القطع؛ الجزء الأخير يفضل محتفظ بأي separator زيادة.", ch18RecorderSplitCode, "e18-5-splitn", []string{"fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 2, Caption: "افصل لحد تلات أجزاء."},
		{Line: 3, Caption: "اطبع العدد والرسالة.", Output: "3 battery|warm"},
	})}},
	{ID: "c185.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch18Question("c185.try-split", "توقّع الناتج:\n\n~~~go\n"+ch18RecorderLimitCode+"\n~~~", []string{"4 warm", "3 warm|again", "2 relay|warm|again"}, 1, []string{"SplitN وقفت قبل إنشاء جزء رابع.", "الجزء الأخير احتفظ بالـseparator الزيادة.", "العدد المحدد للقطع هو تلاتة."}),
		ch18Walkthrough("c185.try-reveal", "امشِ مع حد القطع وشوف الجزء الأخير.", ch18RecorderLimitCode, "e18-5-limit", []string{"fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "اقسم على أول فاصلين."},
			{Line: 3, Caption: "اطبع العدد والجزء الأخير.", Output: "3 warm|again"},
		}),
		ch18Sandbox("c185.sandbox", "ch18-splitn-limit", "زوّد الحد.", "Raise the limit.", "splitn", false),
	}},
	{ID: "c185.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch18Markdown("c185.watch", "عدّ records بشكل مستقل عن sequence المكتوبة جوه السطر.\n\nالـrecord الغلط له تشخيص؛ ~Scanner.Err~ مشكلة في المصدر كله.")}},
	{ID: "c185.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch18Markdown("c185.plan", "خلي تحليل كل record واضح لوحده.\n\nبعد الـloop، بص على حالة الـscanner قبل اعتبار الملخص صالح."),
		ch18Question("c185.check-error", "إيه اللي يفرّق error في record عن error في قراءة المصدر؟", []string{"مكان ظهور الخطأ", "وجود message بس", "طول السطر"}, 0, []string{"error المصدر بيظهر من ~Scanner.Err~ بعد القراءة.", "الـmessage مش بتحدد مرحلة الفشل.", "طول السطر لوحده مش كفاية."}),
	}},
	{ID: "c185.turn-step", Title: "دورك", Blocks: []LessonBlock{ch18Markdown("c185.turn", "ابدأ بـSplitN لرسالة توضيحية، وبعدها جرّب stream سليمة وأخرى فيها error قراءة.")}},
}

var ch18RecorderEnglish = []LessonStep{
	{ID: "c185.mission-step", Title: "Mission", Blocks: []LessonBlock{ch18Markdown("c185.mission", "Analyze records from a stream and note the numbers of malformed ones.\n\nSeparate a bad record from a source failure during reading; the latter prevents an incomplete summary.\n\nThe physical line number differs from sequence data inside it. Point diagnostics to the line that was read.")}},
	{ID: "c185.idea-step", Title: "Idea", Blocks: []LessonBlock{ch18Walkthrough("c185.idea-split", "~SplitN~ limits the number of pieces; the last piece keeps any extra separator.", ch18RecorderSplitCode, "e18-5-splitn", []string{"fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Split into at most three parts."},
		{Line: 3, Caption: "Print count and message.", Output: "3 battery|warm"},
	})}},
	{ID: "c185.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch18Question("c185.try-split", "Predict the output:\n\n~~~go\n"+ch18RecorderLimitCode+"\n~~~", []string{"4 warm", "3 warm|again", "2 relay|warm|again"}, 1, []string{"SplitN stops before making a fourth piece.", "The final piece keeps the extra separator.", "The requested limit is three pieces."}),
		ch18Walkthrough("c185.try-reveal", "Follow the split limit and inspect the last piece.", ch18RecorderLimitCode, "e18-5-limit", []string{"fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Split on the first two separators."},
			{Line: 3, Caption: "Print count and final piece.", Output: "3 warm|again"},
		}),
		ch18Sandbox("c185.sandbox", "ch18-splitn-limit", "زوّد الحد.", "Raise the limit.", "splitn", true),
	}},
	{ID: "c185.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch18Markdown("c185.watch", "Count records independently of sequence numbers inside each line.\n\nA malformed record gets a diagnosis; ~Scanner.Err~ indicates a source failure.")}},
	{ID: "c185.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch18Markdown("c185.plan", "Keep parsing for each record explicit.\n\nAfter the loop, inspect scanner state before accepting the summary."),
		ch18Question("c185.check-error", "What separates a record error from a source read error?", []string{"Where the error is reported", "Whether a message exists", "The line length"}, 0, []string{"A source error appears through ~Scanner.Err~ after reading.", "A message does not identify the failing stage.", "Line length alone is not enough."}),
	}},
	{ID: "c185.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch18Markdown("c185.turn", "Start with SplitN on an illustrative message, then try a valid stream and a read failure.")}},
}

var ch18CrewArabic = []LessonStep{
	{ID: "c186.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch18Markdown("c186.mission", "كمّل دليل الطاقم: تسجيل، بحث، وحفظ من خلال interface صغيرة.\n\nحافظ على السجلات الأصلية وخلي النتائج مرتبة من غير ربط الدليل بمخزن واحد.\n\nوتأكد إن التغييرات ما بتتسربش.")}},
	{ID: "c186.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch18Walkthrough("c186.idea-interface", "الـinterface بتوصف السلوك المطلوب؛ أي نوع يطبّق methods المناسبة يقدر يستخدمها caller.", ch18CrewWriterCode, "e18-6-interface", []string{"bytes", "fmt", "io"}, []LessonFrame{
		{Line: 0, Caption: "قبل التشغيل."},
		{Line: 1, Caption: "الـfunction بتستقبل interface."},
		{Line: 6, Caption: "Buffer يطبّق io.Writer."},
		{Line: 7, Caption: "اكتب من خلال الحد المشترك."},
		{Line: 8, Caption: "اطبع النص والـerror اللي رجّعته الكتابة.", Output: "ready <nil>"},
	})}},
	{ID: "c186.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch18Question("c186.try-search", "توقّع الناتج:\n\n~~~go\n"+ch18CrewSearchCode+"\n~~~", []string{"false", "true", "pilot"}, 1, []string{"النص lot موجود جوه pilot.", "Contains بترجع true لما النص موجود.", "الطباعة بتعرض نتيجة مش الاسم."}),
		ch18Walkthrough("c186.try-reveal", "جرّب Contains على نسخة بحروف موحّدة.", ch18CrewSearchCode, "e18-6-search", []string{"fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التشغيل."},
			{Line: 2, Caption: "وحّد حالة الأحرف."},
			{Line: 3, Caption: "ابحث واطبع النتيجة.", Output: "true"},
		}),
		ch18Sandbox("c186.sandbox", "ch18-normalized-search", "جرّب البحث.", "Try the search.", "search", false),
	}},
	{ID: "c186.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch18Markdown("c186.watch", "افحص البيانات قبل تعديل الـslice، وسيب النص الأصلي زي ما اتسجل.\n\nرتّب نسخة للعرض، وتعامل مع error التخزين كجزء من النتيجة.")}},
	{ID: "c186.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch18Markdown("c186.plan", "افصل منطق الدليل عن طريقة التخزين.\n\nراجع إن المخزن يستقبل نسخة مستقلة وإن البحث ما يغيّرش السجلات."),
		ch18Question("c186.check-store", "لو ~Save~ رجّعت error، تعمل إيه؟", []string{"ترجّعه للـcaller", "تتجاهله وتكمل", "تمسح سجل الطاقم"}, 0, []string{"الـcaller محتاج يعرف إن الحفظ فشل.", "تجاهل الـerror بيخفي فشل التخزين.", "مسح السجل مش بيعالج فشل الحفظ."}),
	}},
	{ID: "c186.turn-step", Title: "دورك", Blocks: []LessonBlock{ch18Markdown("c186.turn", "ابدأ بتخزين value واحدة، جرّب بحث مختلف، وبعدها افحص نسخة تصل للمخزن.")}},
}

var ch18CrewEnglish = []LessonStep{
	{ID: "c186.mission-step", Title: "Mission", Blocks: []LessonBlock{ch18Markdown("c186.mission", "Finish the crew directory: register, search, and save through a small interface.\n\nPreserve the original records and keep results ordered without tying the directory to one store.\n\nAlso check that changes do not leak back.")}},
	{ID: "c186.idea-step", Title: "Idea", Blocks: []LessonBlock{ch18Walkthrough("c186.idea-interface", "An interface describes required behavior; any type with matching methods can serve its caller.", ch18CrewWriterCode, "e18-6-interface", []string{"bytes", "fmt", "io"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "The function accepts an interface."},
		{Line: 6, Caption: "Buffer implements io.Writer."},
		{Line: 7, Caption: "Write through the shared boundary."},
		{Line: 8, Caption: "Print text and write error.", Output: "ready <nil>"},
	})}},
	{ID: "c186.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch18Question("c186.try-search", "Predict the output:\n\n~~~go\n"+ch18CrewSearchCode+"\n~~~", []string{"false", "true", "pilot"}, 1, []string{"The substring lot appears inside pilot.", "Contains returns true when the substring exists.", "The print shows a search result, not the name."}),
		ch18Walkthrough("c186.try-reveal", "Try Contains on a normalized copy.", ch18CrewSearchCode, "e18-6-search", []string{"fmt", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Normalize letter case."},
			{Line: 3, Caption: "Search and print the result.", Output: "true"},
		}),
		ch18Sandbox("c186.sandbox", "ch18-normalized-search", "جرّب البحث.", "Try the search.", "search", true),
	}},
	{ID: "c186.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch18Markdown("c186.watch", "Validate before changing the slice, and keep original text as registered.\n\nSort a copy for display, and treat storage errors as part of the result.")}},
	{ID: "c186.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch18Markdown("c186.plan", "Separate directory logic from storage.\n\nCheck that the store receives an independent copy and search leaves records unchanged."),
		ch18Question("c186.check-store", "What should happen when ~Save~ returns an error?", []string{"Return it to the caller", "Ignore it and continue", "Delete the crew record"}, 0, []string{"The caller needs to know saving failed.", "Ignoring the error hides a storage failure.", "Deleting a record does not fix a failed save."}),
	}},
	{ID: "c186.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch18Markdown("c186.turn", "Start by registering one value, try another search, then inspect the copy given to storage.")}},
}

const ch18MissionLogTableCode = `func main() {
	cases := []struct{ name string; value int }{
		{"small", 2},
		{"large", 4},
	}
	for _, tc := range cases {
		fmt.Println(tc.name, tc.value)
	}
}`

const ch18LogIndexCode = `func main() {
	values := []int{7, 4}
	fmt.Println(values[1])
}`

const ch18BlackBoxJSONCode = "type Beacon struct {\n\tCharge int `json:\"charge\"`\n}\nfunc main() {\n\tb := Beacon{Charge: 7}\n\traw, err := json.Marshal(b)\n\tfmt.Println(string(raw), err)\n}"

const ch18BlackBoxMapCode = `func main() {
	raw, err := json.Marshal(map[string]int{"hull": 3})
	fmt.Println(string(raw), err)
}`

const ch18CargoCSVCode = `func main() {
	r := csv.NewReader(strings.NewReader("north,4\n"))
	row, err := r.Read()
	fmt.Println(row, err)
}`

const ch18CargoQuotedCode = `func main() {
	r := csv.NewReader(strings.NewReader("\"west, upper\",4\n"))
	row, err := r.Read()
	fmt.Println(row, err)
}`

const ch18SignalUnicodeCode = `func main() {
	words := strings.FieldsFunc("Echo—Relay 7", func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	fmt.Println(strings.ToLower(words[0]), len(words))
}`

const ch18SignalLowerCode = `func main() {
	word := strings.ToLower("Echo")
	fmt.Println(word)
}`

const ch18RecorderSplitCode = `func main() {
	parts := strings.SplitN("12|engine|battery|warm", "|", 3)
	fmt.Println(len(parts), parts[2])
}`

const ch18RecorderLimitCode = `func main() {
	parts := strings.SplitN("9|relay|warm|again", "|", 3)
	fmt.Println(len(parts), parts[2])
}`

const ch18CrewWriterCode = `func emit(w io.Writer, text string) error {
	_, err := fmt.Fprint(w, text)
	return err
}
func main() {
	var b bytes.Buffer
	err := emit(&b, "ready")
	fmt.Println(b.String(), err)
}`

const ch18CrewSearchCode = `func main() {
	name := strings.ToLower("Pilot")
	fmt.Println(strings.Contains(name, "lot"))
}`
