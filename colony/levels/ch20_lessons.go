package levels

func ch20Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch20Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch20Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch20Localized(ar, en string) map[string]any {
	return map[string]any{"ar": ar, "en": en}
}

func ch20SandboxBlock(id, kind, text string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch20EnglishSandboxBlock(id, kind, text string, config map[string]any) LessonBlock {
	return ch20SandboxBlock(id, kind, text, ch20EnglishConfig(config).(map[string]any))
}

func ch20EnglishConfig(value any) any {
	switch value := value.(type) {
	case map[string]any:
		if english, ok := value["en"].(string); ok {
			if _, localized := value["ar"].(string); localized {
				return map[string]any{"ar": english, "en": english}
			}
		}
		copy := make(map[string]any, len(value))
		for key, item := range value {
			copy[key] = ch20EnglishConfig(item)
		}
		return copy
	case []any:
		copy := make([]any, len(value))
		for index, item := range value {
			copy[index] = ch20EnglishConfig(item)
		}
		return copy
	default:
		return value
	}
}

var ch20BulletinArabic = []LessonStep{
	{ID: "c201.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch20Markdown("c201.mission", "1. حوّل نص النشرة لـ HTML آمن بـ~RenderBulletin~.\n2. اتعامل مع tagين معروفين بس؛ الأقواس التانية تفضل نص عادي.\n3. الـstarter بيرجّع error مؤقت، والمهمة من غير Report.")}},
	{ID: "c201.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch20Walkthrough("c201.state-example", "اعتبر كل tag معروفة key لحالة. ~HasPrefix~ بتتعرف على بداية النص، وبعدها اقرا الباقي حسب الحالة.\n\nده مثال توضيحي صغير، مش parser كامل. خليك عارف هل النص عادي ولا جوه tag.", `func main() {
var source = "[alert]low[/alert]"
mode := "plain"
if strings.HasPrefix(source, "[alert]") {
	mode = "alert"
}
fmt.Println(mode)
}`, "e201-state", []string{"fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 2, Caption: "النص يبدأ بوسم معروف.", Variables: []LessonVariable{{Name: "source", Value: "[alert]low[/alert]"}}},
		{Line: 3, Caption: "الحالة تبدأ plain.", Variables: []LessonVariable{{Name: "mode", Value: "plain"}}},
		{Line: 5, Caption: "الحالة بقت alert.", Variables: []LessonVariable{{Name: "mode", Value: "alert"}}},
		{Line: 7, Caption: "اطبع اسم الحالة.", Variables: []LessonVariable{{Name: "mode", Value: "alert"}}, Output: "alert"},
	})}},
	{ID: "c201.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch20Question("c201.predict-escape", "توقّع الناتج بعد escaping النص:\n\n~~~go\nfunc main() {\nvar out strings.Builder\nout.WriteString(html.EscapeString(\"<relay>\"))\nfmt.Println(out.String())\n}\n~~~", []string{"<relay>", "&lt;relay&gt;", "&amp;lt;relay&amp;gt;"}, 1, []string{"دي العلامات قبل التأمين.", "علامتا الزاوية اتحولوا لـ entities.", "ده escaping مرتين، وده مش اللي حصل."}),
		ch20Walkthrough("c201.escape-reveal", "تابع النص وهو بيتحوّل قبل ما يضاف للـBuilder.", `func main() {
var out strings.Builder
out.WriteString(html.EscapeString("<relay>"))
fmt.Println(out.String())
}`, "e201-escape", []string{"fmt", "html", "strings"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "Builder جاهزة فاضية.", Variables: []LessonVariable{{Name: "out", Value: ""}}},
			{Line: 3, Caption: "اتضاف نص آمن.", Variables: []LessonVariable{{Name: "out", Value: "&lt;relay&gt;"}}},
			{Line: 4, Caption: "الناتج فيه entities.", Variables: []LessonVariable{{Name: "out", Value: "&lt;relay&gt;"}}, Output: "&lt;relay&gt;"},
		}),
		ch20SandboxBlock("c201.try-sandbox", "ch20-bulletin-state", "تابع state والـoutput خطوة خطوة.", map[string]any{
			"source": "[alert]<relay>[/alert]",
			"frames": []any{
				map[string]any{"piece": ch20Localized("بداية النص", "Start of input"), "state": ch20Localized("نص عادي", "plain"), "output": "", "caption": ch20Localized("ابدأ في حالة النص العادي.", "Start in plain text.")},
				map[string]any{"piece": ch20Localized("[alert]", "[alert]"), "state": ch20Localized("تنبيه", "alert"), "output": "<strong>", "caption": ch20Localized("اتعرّف على الوسم المعروف.", "Recognize the known opening tag.")},
				map[string]any{"piece": ch20Localized("<relay>", "<relay>"), "state": ch20Localized("تنبيه", "alert"), "output": "<strong>&lt;relay&gt;", "caption": ch20Localized("اعمِل escaping لنص المستخدم.", "Escape the user text once.")},
				map[string]any{"piece": ch20Localized("[/alert]", "[/alert]"), "state": ch20Localized("نص عادي", "plain"), "output": "<strong>&lt;relay&gt;</strong>", "caption": ch20Localized("اقفل الوسم وارجع للنص العادي.", "Close the tag; return to plain text.")},
			},
			"labels": map[string]any{
				"source":   ch20Localized("النص", "Input"),
				"piece":    ch20Localized("الجزء", "Piece"),
				"state":    ch20Localized("الحالة", "State"),
				"output":   ch20Localized("ناتج آمن", "Safe output"),
				"previous": ch20Localized("السابق", "Previous"),
				"next":     ch20Localized("التالي", "Next"),
				"reset":    ch20Localized("من الأول", "Reset"),
			},
		}),
	}},
	{ID: "c201.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch20Markdown("c201.watch", "- اكتب tags بتاعتك مباشرة، واعمِل escaping لنص المستخدم مرة واحدة.\n- القفل الغلط أو نهاية النص وtag لسه مفتوحة معناهم error من غير output جزئي.")}},
	{ID: "c201.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch20Markdown("c201.plan", "قسّم الشغل: ميّز delimiter، حدّث الحالة، وبعدها أضف النص الآمن.\n\nراجع newline حسب الحالة، وخلي أي syntax مش معروفة حرفية."),
		ch20Question("c201.escape-place", "لو النص جوه ~[code]~ فيه ~<~، تعمل إيه؟", []string{"تسيبه حرفي", "تعمله escaping", "تعتبره tag جديدة"}, 1, []string{"ده بيخلّي HTML يفسّره كعنصر.", "نص المستخدم يفضل غير موثوق جوه كل الحالات.", "المحتوى مش delimiter معروفة."}),
	}},
	{ID: "c201.turn-step", Title: "دورك", Blocks: []LessonBlock{ch20Markdown("c201.turn", "ابدأ بحالة واحدة، وشغّلها على نص صغير. زوّد الحالات واحدة واحدة وراجع الناتج والـerrors.")}},
}

var ch20BulletinEnglish = []LessonStep{
	{ID: "c201.mission-step", Title: "Mission", Blocks: []LessonBlock{ch20Markdown("c201.mission", "1. Render bulletin text as safe HTML with ~RenderBulletin~.\n2. Handle only two recognized tags; other bracket forms stay ordinary text.\n3. The starter returns a temporary error, and this lab uses no Report.")}},
	{ID: "c201.idea-step", Title: "Idea", Blocks: []LessonBlock{ch20Walkthrough("c201.state-example", "Treat each recognized tag as a state switch. ~HasPrefix~ recognizes the start of the input; then the rest is read according to that state.\n\nThis is a small illustration, not a complete parser. Track whether text is plain or inside a tag.", `func main() {
var source = "[alert]low[/alert]"
mode := "plain"
if strings.HasPrefix(source, "[alert]") {
	mode = "alert"
}
fmt.Println(mode)
}`, "e201-state", []string{"fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Input starts with a known tag.", Variables: []LessonVariable{{Name: "source", Value: "[alert]low[/alert]"}}},
		{Line: 3, Caption: "State starts plain.", Variables: []LessonVariable{{Name: "mode", Value: "plain"}}},
		{Line: 5, Caption: "State changes to alert.", Variables: []LessonVariable{{Name: "mode", Value: "alert"}}},
		{Line: 7, Caption: "Print the state name.", Variables: []LessonVariable{{Name: "mode", Value: "alert"}}, Output: "alert"},
	})}},
	{ID: "c201.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch20Question("c201.predict-escape", "Predict the output after escaping the text:\n\n~~~go\nfunc main() {\nvar out strings.Builder\nout.WriteString(html.EscapeString(\"<relay>\"))\nfmt.Println(out.String())\n}\n~~~", []string{"<relay>", "&lt;relay&gt;", "&amp;lt;relay&amp;gt;"}, 1, []string{"Those are the marks before escaping.", "The angle brackets become entities.", "That would escape twice, which did not happen."}),
		ch20Walkthrough("c201.escape-reveal", "Follow the text as it is escaped before it enters the Builder.", `func main() {
var out strings.Builder
out.WriteString(html.EscapeString("<relay>"))
fmt.Println(out.String())
}`, "e201-escape", []string{"fmt", "html", "strings"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Builder starts empty.", Variables: []LessonVariable{{Name: "out", Value: ""}}},
			{Line: 3, Caption: "Safe text is appended.", Variables: []LessonVariable{{Name: "out", Value: "&lt;relay&gt;"}}},
			{Line: 4, Caption: "Output contains entities.", Variables: []LessonVariable{{Name: "out", Value: "&lt;relay&gt;"}}, Output: "&lt;relay&gt;"},
		}),
		ch20EnglishSandboxBlock("c201.try-sandbox", "ch20-bulletin-state", "Track state and escaped output.", map[string]any{
			"source": "[alert]<relay>[/alert]",
			"frames": []any{
				map[string]any{"piece": ch20Localized("بداية النص", "Start of input"), "state": ch20Localized("نص عادي", "plain"), "output": "", "caption": ch20Localized("ابدأ في حالة النص العادي.", "Start in plain text.")},
				map[string]any{"piece": ch20Localized("[alert]", "[alert]"), "state": ch20Localized("تنبيه", "alert"), "output": "<strong>", "caption": ch20Localized("اتعرّف على الوسم المعروف.", "Recognize the known opening tag.")},
				map[string]any{"piece": ch20Localized("<relay>", "<relay>"), "state": ch20Localized("تنبيه", "alert"), "output": "<strong>&lt;relay&gt;", "caption": ch20Localized("اعمِل escaping لنص المستخدم.", "Escape the user text once.")},
				map[string]any{"piece": ch20Localized("[/alert]", "[/alert]"), "state": ch20Localized("نص عادي", "plain"), "output": "<strong>&lt;relay&gt;</strong>", "caption": ch20Localized("اقفل الوسم وارجع للنص العادي.", "Close the tag; return to plain text.")},
			},
			"labels": map[string]any{
				"source":   ch20Localized("النص", "Input"),
				"piece":    ch20Localized("الجزء", "Piece"),
				"state":    ch20Localized("الحالة", "State"),
				"output":   ch20Localized("ناتج آمن", "Safe output"),
				"previous": ch20Localized("السابق", "Previous"),
				"next":     ch20Localized("التالي", "Next"),
				"reset":    ch20Localized("من الأول", "Reset"),
			},
		}),
	}},
	{ID: "c201.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch20Markdown("c201.watch", "- Write your own tags directly, and escape user text exactly once.\n- A mismatched close or EOF inside a tag is an error with no partial output.")}},
	{ID: "c201.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch20Markdown("c201.plan", "Split the work: recognize a delimiter, update state, then append safe text.\n\nCheck newlines by state, and leave unknown syntax literal."),
		ch20Question("c201.escape-place", "If text inside ~[code]~ contains ~<~, what should happen?", []string{"Keep it literal", "Escape it", "Treat it as a new tag"}, 1, []string{"That lets HTML interpret it as an element.", "User text stays untrusted in every state.", "The content is not a recognized delimiter."}),
	}},
	{ID: "c201.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch20Markdown("c201.turn", "Start with one state and run it on a small string. Add states one at a time, checking output and errors.")}},
}

var ch20ValidatorArabic = []LessonStep{
	{ID: "c202.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch20Markdown("c202.mission", "1. راجع بطاقة الطاقم بثلاث فحوص مستقلة.\n2. اجمع الأسباب لو أكتر من حقل فشل، بترتيب ثابت.\n3. الـstarter فيه الأنواع والـsentinels؛ كمّل السلوك بس.")}},
	{ID: "c202.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch20Walkthrough("c202.join-example", "الـerror ممكن يفضل شايل سبب محدد حتى لو اتضاف له شرح. ~%w~ بيحافظ على السبب، و~errors.Join~ بتجمع errors مستقلة.\n\n~errors.Is~ بتدور على السبب نفسه وسط التغليف والتجميع. استخدم ~regexp~ للأنماط الصغيرة، و~^~ و~$~ عشان الفحص يشمل الـvalue كلها.\n\nللتاريخ استخدم ~time.Parse~ من 19-1 وبعدها ~Format~ للتنسيق. ولفحص local name استخدم ~unicode.IsSpace~ على كل rune.", `func main() {
var errPilot = errors.New("pilot")
var errRadio = errors.New("radio")
combined := errors.Join(fmt.Errorf("card: %w", errPilot), errRadio)
fmt.Println(errors.Is(combined, errPilot), errors.Is(combined, errRadio))
}`, "e202-join", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create the first cause."},
		{Line: 3, Caption: "Create a second cause."},
		{Line: 4, Caption: "Keep both causes together."},
		{Line: 5, Caption: "Both causes are found.", Output: "true true"},
	})}},
	{ID: "c202.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch20Question("c202.predict-regex", "توقّع نتيجة الـpattern مع ~\"relay!\"~:\n\n~~~go\nfunc main() {\npattern := regexp.MustCompile(\"^[a-z]+$\")\nfmt.Println(pattern.MatchString(\"relay!\"))\n}\n~~~", []string{"true", "false", "panic"}, 1, []string{"الـ! مش ضمن الحروف المسموحة.", "الـanchors بتطلب تطابق النص كله.", "الـpattern اتعملت compile قبل الفحص."}),
		ch20Walkthrough("c202.regex-reveal", "تابع فحص string صغيرة بالـanchors.", `func main() {
pattern := regexp.MustCompile("^[a-z]+$")
fmt.Println(pattern.MatchString("relay!"))
}`, "e202-regex", []string{"fmt", "regexp"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الـpattern بتطلب حروف صغيرة."},
			{Line: 3, Caption: "علامة التعجب مرفوضة.", Output: "false"},
		}),
		ch20SandboxBlock("c202.try-sandbox", "ch20-error-join", "بدّل الحقول وشوف تجميع الأسباب.", map[string]any{
			"fields": []any{
				map[string]any{"id": "pilot", "label": ch20Localized("اسم الطيار", "Pilot name"), "cause": ch20Localized("اسم الطيار ناقص", "Pilot name missing"), "invalid": true},
				map[string]any{"id": "radio", "label": ch20Localized("نداء الراديو", "Radio call"), "cause": ch20Localized("نداء الراديو غير صالح", "Radio call invalid"), "invalid": false},
				map[string]any{"id": "date", "label": ch20Localized("تاريخ المهمة", "Mission date"), "cause": ch20Localized("التاريخ غير صالح", "Date invalid"), "invalid": true},
			},
			"labels": map[string]any{
				"valid":   ch20Localized("سليم", "valid"),
				"invalid": ch20Localized("فيه error", "invalid"),
				"causes":  ch20Localized("الأسباب بالترتيب", "Causes in order"),
			},
		}),
	}},
	{ID: "c202.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch20Markdown("c202.watch", "- regex من غير anchors ممكن تقبل جزء من value غلط.\n- التاريخ محتاج parser حقيقي؛ regex لوحدها مش بتعرف إن اليوم موجود.")}},
	{ID: "c202.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch20Markdown("c202.plan", "خلي كل validator مسؤولة عن حقل واحد وسبب واحد.\n\nبعدها اجمع النتائج، وراجع التصنيفات بالترتيب المعلن بدل ترتيب map."),
		ch20Question("c202.classify", "إيه اللي يخلّي سبب الـerror قابل للفحص بعد إضافة شرح؟", []string{"نقارن نص الرسالة", "نلف بـ~%w~ ونستخدم ~errors.Is~", "نخزّن النص في map"}, 1, []string{"نص الرسالة ممكن يتغير ومش هوية السبب.", "الـwrapping يحفظ السبب للتصنيف.", "ترتيب map مش contract ثابت."}),
	}},
	{ID: "c202.turn-step", Title: "دورك", Blocks: []LessonBlock{ch20Markdown("c202.turn", "ابدأ بفحص واحد، وبعدها جرّب بطاقة فيها أكتر من مشكلة. اتأكد إن كل سبب ظاهر مرة واحدة وبالترتيب.")}},
}

var ch20ValidatorEnglish = []LessonStep{
	{ID: "c202.mission-step", Title: "Mission", Blocks: []LessonBlock{ch20Markdown("c202.mission", "1. Review the crew card with three independent checks.\n2. Join causes from failed fields in a fixed order.\n3. The starter provides the types and sentinels; complete the behavior.")}},
	{ID: "c202.idea-step", Title: "Idea", Blocks: []LessonBlock{ch20Walkthrough("c202.join-example", "An error can keep a specific cause even when it gains context. ~%w~ preserves that cause, and ~errors.Join~ combines independent errors.\n\n~errors.Is~ searches through wrapping and joining. Use ~regexp~ for small patterns, with ~^~ and ~$~ to check the whole value.\n\nReuse ~time.Parse~ from 19-1 and then ~Format~ for dates. Check each local-part rune with ~unicode.IsSpace~.", `func main() {
var errPilot = errors.New("pilot")
var errRadio = errors.New("radio")
combined := errors.Join(fmt.Errorf("card: %w", errPilot), errRadio)
fmt.Println(errors.Is(combined, errPilot), errors.Is(combined, errRadio))
}`, "e202-join", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Create the first cause."},
		{Line: 3, Caption: "Create a second cause."},
		{Line: 4, Caption: "Keep both causes together."},
		{Line: 5, Caption: "Both causes are found.", Output: "true true"},
	})}},
	{ID: "c202.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch20Question("c202.predict-regex", "Predict the pattern result for ~\"relay!\"~:\n\n~~~go\nfunc main() {\npattern := regexp.MustCompile(\"^[a-z]+$\")\nfmt.Println(pattern.MatchString(\"relay!\"))\n}\n~~~", []string{"true", "false", "panic"}, 1, []string{"The ! is outside the allowed letters.", "Anchors require the whole string to match.", "The pattern compiled before the check."}),
		ch20Walkthrough("c202.regex-reveal", "Follow a small string check with anchors.", `func main() {
pattern := regexp.MustCompile("^[a-z]+$")
fmt.Println(pattern.MatchString("relay!"))
}`, "e202-regex", []string{"fmt", "regexp"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Pattern requires lowercase letters."},
			{Line: 3, Caption: "The exclamation is rejected.", Output: "false"},
		}),
		ch20EnglishSandboxBlock("c202.try-sandbox", "ch20-error-join", "Toggle fields; inspect joined causes.", map[string]any{
			"fields": []any{
				map[string]any{"id": "pilot", "label": ch20Localized("اسم الطيار", "Pilot name"), "cause": ch20Localized("اسم الطيار ناقص", "Pilot name missing"), "invalid": true},
				map[string]any{"id": "radio", "label": ch20Localized("نداء الراديو", "Radio call"), "cause": ch20Localized("نداء الراديو غير صالح", "Radio call invalid"), "invalid": false},
				map[string]any{"id": "date", "label": ch20Localized("تاريخ المهمة", "Mission date"), "cause": ch20Localized("التاريخ غير صالح", "Date invalid"), "invalid": true},
			},
			"labels": map[string]any{
				"valid":   ch20Localized("سليم", "valid"),
				"invalid": ch20Localized("فيه error", "invalid"),
				"causes":  ch20Localized("الأسباب بالترتيب", "Causes in order"),
			},
		}),
	}},
	{ID: "c202.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch20Markdown("c202.watch", "- A regex without anchors can accept part of an invalid value.\n- Dates need a real parser; a regex cannot know whether a day exists.")}},
	{ID: "c202.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch20Markdown("c202.plan", "Give each validator one field and one cause.\n\nThen join the results and check classifications in the declared order instead of map order."),
		ch20Question("c202.classify", "How can a cause stay inspectable after adding context?", []string{"Compare the message text", "Wrap with ~%w~ and use ~errors.Is~", "Store the text in a map"}, 1, []string{"Message wording can change and is not cause identity.", "Wrapping preserves the cause for classification.", "Map order is not a fixed contract."}),
	}},
	{ID: "c202.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch20Markdown("c202.turn", "Start with one check, then try a card with several problems. Confirm every cause appears once and in order.")}},
}

var ch20LexerArabic = []LessonStep{
	{ID: "c203.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch20Markdown("c203.mission", "1. حوّل أوامر الطيران لـtokens في ~LexFlight~.\n2. كل token تحتفظ بنوعها وتهجئتها ومكانها بالـbyte.\n3. الـstarter عرّف ~CommandToken~ و~CommandError~؛ مفيش تنفيذ أو Report.")}},
	{ID: "c203.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch20Walkthrough("c203.token-example", "الـtoken بتسجل دليل من المصدر: النوع، النص الأصلي، والـoffset. الـoffset هنا رقم byte من صفر، مش رقم حرف ظاهر.\n\nامشِ byte byte واستهلك أطول run. ~strconv.ParseInt~ بيرجّع error لو الرقم مش قابل للقراءة؛ ما تتجاهلش الـerror.", `type Token struct {
	Kind, Text string
	Offset int
}
func main() {
	t := Token{Kind: "name", Text: "bay_2", Offset: 4}
	fmt.Println(t.Kind, t.Text, t.Offset)
}`, "e203-token", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "Token keeps its kind."},
		{Line: 2, Caption: "Source spelling stays available."},
		{Line: 6, Caption: "Example token starts at byte four.", Variables: []LessonVariable{{Name: "t", Value: "{name bay_2 4}"}}},
		{Line: 7, Caption: "Print the stored fields.", Variables: []LessonVariable{{Name: "t", Value: "{name bay_2 4}"}}, Output: "name bay_2 4"},
	})}},
	{ID: "c203.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch20Question("c203.predict-offsets", "توقّع السطرين اللي هيظهروا:\n\n~~~go\nfunc main() {\nsource := \"go\"\nfor i := 0; i < len(source); i++ {\nfmt.Println(i, string(source[i]))\n}\n}\n~~~", []string{"0 g ثم 1 o", "1 g ثم 2 o", "g ثم o من غير أرقام"}, 0, []string{"الـindex بيبدأ من صفر وبيعد bytes النص.", "الـindex مش بيبدأ من واحد.", "الكود بيطبع index مع كل byte."}),
		ch20Walkthrough("c203.offset-reveal", "تابع index والـbyte في مثال ASCII صغير.", `func main() {
source := "go"
for i := 0; i < len(source); i++ {
fmt.Println(i, string(source[i]))
}
}`, "e203-offset", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "source فيها 2 bytes."},
			{Line: 4, Caption: "أول byte عند صفر.", Output: "0 g"},
			{Line: 4, Caption: "الـbyte التانية عند واحد.", Output: "0 g\n1 o"},
		}),
		ch20SandboxBlock("c203.try-sandbox", "ch20-lexer-offset", "حرّك المسح وراقب token والـoffset.", map[string]any{
			"source": "go_2 > 7",
			"frames": []any{
				map[string]any{"cursor": "0", "token": ch20Localized("لسه مفيش", "not yet"), "offset": "—", "caption": ch20Localized("ابدأ من أول byte.", "Start at the first byte.")},
				map[string]any{"cursor": "4", "token": "name go_2", "offset": "0", "caption": ch20Localized("استهلك الاسم بالكامل.", "Consume the whole name.")},
				map[string]any{"cursor": "6", "token": ">", "offset": "5", "caption": ch20Localized("تخطّى المسافة وخزّن الرمز.", "Skip space; keep the symbol.")},
				map[string]any{"cursor": "8", "token": "number 7", "offset": "7", "caption": ch20Localized("سجّل بداية الرقم بالـbyte.", "Record the number's byte offset.")},
				map[string]any{"cursor": "8", "token": "eof", "offset": "8", "caption": ch20Localized("EOF بعد آخر byte.", "EOF follows the final byte.")},
			},
			"labels": map[string]any{
				"source":   ch20Localized("النص", "Input"),
				"cursor":   ch20Localized("الـindex", "Cursor"),
				"token":    ch20Localized("آخر token", "Current token"),
				"offset":   ch20Localized("offset", "Offset"),
				"previous": ch20Localized("السابق", "Previous"),
				"next":     ch20Localized("التالي", "Next"),
				"reset":    ch20Localized("من الأول", "Reset"),
			},
		}),
	}},
	{ID: "c203.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch20Markdown("c203.watch", "- اقرا أطول run من الاسم أو الأرقام؛ ما تطلعش token لكل حرف.\n- أي byte مجهولة توقف التحليل بـerror عند مكانها، ومفيش tokens جزئية.")}},
	{ID: "c203.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch20Markdown("c203.plan", "امشِ على المصدر من اليسار لليمين: خزّن البداية، حدّد النوع، واستهلك الـstring كاملة.\n\nسيب علامات الجمل منفصلة، وأضف ~eof~ في النهاية."),
		ch20Question("c203.number-sign", "علامة السالب في ~-8~ تتعامل معاها إزاي؟", []string{"جزء من رقم واحد", "token مستقلة عن الرقم", "مسافة تتشال"}, 1, []string{"المطلوب إن الرقم يحتفظ بتهجئته كأرقام.", "الـlexer تفصل الرمز عن run الأرقام.", "السالب رمز معروف في اللغة."}),
	}},
	{ID: "c203.turn-step", Title: "دورك", Blocks: []LessonBlock{ch20Markdown("c203.turn", "ابدأ بأسماء وأرقام ورموز قليلة. راجع offsets يدويًا، وبعدها جرّب whitespace وbyte مرفوضة.")}},
}

var ch20LexerEnglish = []LessonStep{
	{ID: "c203.mission-step", Title: "Mission", Blocks: []LessonBlock{ch20Markdown("c203.mission", "1. Turn flight commands into tokens with ~LexFlight~.\n2. Keep each token's kind, spelling, and byte position.\n3. The starter defines ~CommandToken~ and ~CommandError~; do not execute or Report.")}},
	{ID: "c203.idea-step", Title: "Idea", Blocks: []LessonBlock{ch20Walkthrough("c203.token-example", "A token records source evidence: kind, original spelling, and offset. Offsets count bytes from zero, not displayed characters.\n\nConsume each maximal run. Check the error from ~strconv.ParseInt~ before accepting a number.", `type Token struct {
	Kind, Text string
	Offset int
}
func main() {
	t := Token{Kind: "name", Text: "bay_2", Offset: 4}
	fmt.Println(t.Kind, t.Text, t.Offset)
}`, "e203-token", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 1, Caption: "Token keeps its kind."},
		{Line: 2, Caption: "Source spelling stays available."},
		{Line: 6, Caption: "Example token starts at byte four.", Variables: []LessonVariable{{Name: "t", Value: "{name bay_2 4}"}}},
		{Line: 7, Caption: "Print the stored fields.", Variables: []LessonVariable{{Name: "t", Value: "{name bay_2 4}"}}, Output: "name bay_2 4"},
	})}},
	{ID: "c203.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch20Question("c203.predict-offsets", "Predict the two output lines:\n\n~~~go\nfunc main() {\nsource := \"go\"\nfor i := 0; i < len(source); i++ {\nfmt.Println(i, string(source[i]))\n}\n}\n~~~", []string{"0 g then 1 o", "1 g then 2 o", "g then o without numbers"}, 0, []string{"The index starts at zero and counts bytes.", "The index does not start at one.", "The code prints an index for each byte."}),
		ch20Walkthrough("c203.offset-reveal", "Follow the index and byte in a small ASCII example.", `func main() {
source := "go"
for i := 0; i < len(source); i++ {
fmt.Println(i, string(source[i]))
}
}`, "e203-offset", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "source has two bytes."},
			{Line: 4, Caption: "First byte is at zero.", Output: "0 g"},
			{Line: 4, Caption: "Second byte is at one.", Output: "0 g\n1 o"},
		}),
		ch20EnglishSandboxBlock("c203.try-sandbox", "ch20-lexer-offset", "Advance scanning; watch token offsets.", map[string]any{
			"source": "go_2 > 7",
			"frames": []any{
				map[string]any{"cursor": "0", "token": ch20Localized("لسه مفيش", "not yet"), "offset": "—", "caption": ch20Localized("ابدأ من أول byte.", "Start at the first byte.")},
				map[string]any{"cursor": "4", "token": "name go_2", "offset": "0", "caption": ch20Localized("استهلك الاسم بالكامل.", "Consume the whole name.")},
				map[string]any{"cursor": "6", "token": ">", "offset": "5", "caption": ch20Localized("تخطّى المسافة وخزّن الرمز.", "Skip space; keep the symbol.")},
				map[string]any{"cursor": "8", "token": "number 7", "offset": "7", "caption": ch20Localized("سجّل بداية الرقم بالـbyte.", "Record the number's byte offset.")},
				map[string]any{"cursor": "8", "token": "eof", "offset": "8", "caption": ch20Localized("EOF بعد آخر byte.", "EOF follows the final byte.")},
			},
			"labels": map[string]any{
				"source":   ch20Localized("النص", "Input"),
				"cursor":   ch20Localized("الـindex", "Cursor"),
				"token":    ch20Localized("آخر token", "Current token"),
				"offset":   ch20Localized("offset", "Offset"),
				"previous": ch20Localized("السابق", "Previous"),
				"next":     ch20Localized("التالي", "Next"),
				"reset":    ch20Localized("من الأول", "Reset"),
			},
		}),
	}},
	{ID: "c203.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch20Markdown("c203.watch", "- Read the longest name or digit run; do not emit one token per character.\n- An unknown byte stops scanning with an error at its position, without partial tokens.")}},
	{ID: "c203.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch20Markdown("c203.plan", "Walk left to right: save the start, identify the kind, and consume the whole run.\n\nKeep statement symbols separate, then append one ~eof~."),
		ch20Question("c203.number-sign", "How should the minus in ~-8~ be handled?", []string{"Part of one number", "A separate token from the number", "Whitespace to skip"}, 1, []string{"The number should keep its digit spelling.", "The lexer separates the symbol from the digit run.", "Minus is a recognized symbol in this language."}),
	}},
	{ID: "c203.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch20Markdown("c203.turn", "Start with a few names, numbers, and symbols. Check offsets by hand, then try whitespace and a rejected byte.")}},
}

var ch20EvaluatorArabic = []LessonStep{
	{ID: "c204.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch20Markdown("c204.mission", "1. نفّذ ~ExecuteFlight~ باستخدام lexer اللي خلصتها في 20-3.\n2. افهم الأولوية والمتغيرات قبل ما الدرون يتحرك.\n3. الـstarter مجهز interfaces؛ كمّل المفسر في ~flight.go~.")}},
	{ID: "c204.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch20Walkthrough("c204.precedence-example", "ترتيب العمليات بيفصل مستويات التعبير: الضرب أعمق من الجمع. الأقواس بتفرض تجميع مختلف.\n\nالـparser يبني AST، وبعدها الـnodes تقيّم نفسها عن طريق ~FlightExpr~. Recursive descent بتقسم القواعد حسب الأولوية، والأقواس ترجع لـexpression. المثال ده يوضح Go بس، مش parser المهمة.", `func main() {
	fmt.Println(2 + 3*4)
}`, "e204-precedence", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "قبل التنفيذ."},
		{Line: 2, Caption: "الضرب يحصل قبل الجمع.", Output: "14"},
	})}},
	{ID: "c204.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch20Question("c204.predict-associativity", "توقّع ناتج العملية من اليسار لليمين:\n\n~~~go\nfunc main() {\nfmt.Println(9 - 3 - 1)\n}\n~~~", []string{"7", "5", "9"}, 1, []string{"ده هيكون الناتج لو اتجمعت 9-(3-1).", "الطرح على نفس المستوى بيتحسب left associative.", "الـoperator بتطرح الـvalues اللي بعدها من 9."}),
		ch20Walkthrough("c204.associativity-reveal", "تابع طرح متسلسل بنفس أولوية العمليات.", `func main() {
fmt.Println(9 - 3 - 1)
}`, "e204-associativity", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "احسب من الشمال لليمين.", Output: "5"},
		}),
		ch20SandboxBlock("c204.try-sandbox", "ch20-expression-tree", "تابع تجميع التعبير خطوة خطوة.", map[string]any{
			"frames": []any{
				map[string]any{"expression": "4 + 2 * 3", "value": ch20Localized("لسه ما اتحسبتش", "not evaluated"), "caption": ch20Localized("ابدأ بالتعبير كامل.", "Start with the full expression.")},
				map[string]any{"expression": "2 * 3", "value": "6", "caption": ch20Localized("الضرب يرتبط أقوى.", "Multiplication binds more tightly.")},
				map[string]any{"expression": "4 + 6", "value": "10", "caption": ch20Localized("اجمع بعد تقييم الضرب.", "Add after evaluating multiplication.")},
				map[string]any{"expression": "10", "value": "10", "caption": ch20Localized("الـexpression رجّع الـvalue.", "The expression returns a value.")},
			},
			"labels": map[string]any{
				"expression": ch20Localized("التعبير", "Expression"),
				"value":      ch20Localized("الـvalue", "Value"),
				"previous":   ch20Localized("السابق", "Previous"),
				"next":       ch20Localized("التالي", "Next"),
				"reset":      ch20Localized("من الأول", "Reset"),
			},
		}),
	}},
	{ID: "c204.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch20Markdown("c204.watch", "- حلّل البرنامج كله قبل أي حركة؛ syntax error ما يتحولش لأثر في العالم.\n- runtime error يوقف الجمل الجاية، والحركة اللي تمت تفضل.")}},
	{ID: "c204.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch20Markdown("c204.plan", "افصل المراحل: parse للجمل والتعبيرات، ثم evaluation للـAST، ثم استدعاء ~FlightDrive.Step~.\n\nراجع مكان الـerror في الـtoken أو الـoperator المناسبة قبل التجربة."),
		ch20Question("c204.parse-first", "إمتى ينفع أول استدعاء حركة؟", []string{"بعد كل statement مباشرة أثناء parsing", "بعد نجاح parsing للبرنامج كله", "قبل قراءة tokens"}, 1, []string{"ده ممكن يحرّك الدرون قبل اكتشاف syntax error لاحق.", "الـAST الكامل يفصل syntax عن الآثار.", "لازم القواعد تتقري قبل التنفيذ."}),
	}},
	{ID: "c204.turn-step", Title: "دورك", Blocks: []LessonBlock{ch20Markdown("c204.turn", "ابدأ بتعبير صغير، وبعدها جرّب let وthrust. راقب الـerrors والحركة المكتملة كل واحدة لوحدها.")}},
}

var ch20EvaluatorEnglish = []LessonStep{
	{ID: "c204.mission-step", Title: "Mission", Blocks: []LessonBlock{ch20Markdown("c204.mission", "1. Implement ~ExecuteFlight~ using the lexer from 20-3.\n2. Apply precedence and variables before the drone moves.\n3. The starter provides the interfaces; finish the interpreter in ~flight.go~.")}},
	{ID: "c204.idea-step", Title: "Idea", Blocks: []LessonBlock{ch20Walkthrough("c204.precedence-example", "Expression precedence separates levels: multiplication binds more tightly than addition. Parentheses force different grouping.\n\nA parser builds an AST, then nodes evaluate through ~FlightExpr~. Recursive descent uses one function per precedence level, and parentheses re-enter expression parsing. This example shows Go's rules, not the mission parser.", `func main() {
	fmt.Println(2 + 3*4)
}`, "e204-precedence", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: "Before execution."},
		{Line: 2, Caption: "Multiplication comes first.", Output: "14"},
	})}},
	{ID: "c204.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch20Question("c204.predict-associativity", "Predict the result when operations group left to right:\n\n~~~go\nfunc main() {\nfmt.Println(9 - 3 - 1)\n}\n~~~", []string{"7", "5", "9"}, 1, []string{"That would result from grouping 9-(3-1).", "Subtraction at one level is left associative.", "The expression subtracts two values from 9."}),
		ch20Walkthrough("c204.associativity-reveal", "Follow repeated subtraction at one precedence level.", `func main() {
fmt.Println(9 - 3 - 1)
}`, "e204-associativity", []string{"fmt"}, []LessonFrame{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Calculate from left to right.", Output: "5"},
		}),
		ch20EnglishSandboxBlock("c204.try-sandbox", "ch20-expression-tree", "Step through expression grouping.", map[string]any{
			"frames": []any{
				map[string]any{"expression": "4 + 2 * 3", "value": ch20Localized("لسه ما اتحسبتش", "not evaluated"), "caption": ch20Localized("ابدأ بالتعبير كامل.", "Start with the full expression.")},
				map[string]any{"expression": "2 * 3", "value": "6", "caption": ch20Localized("الضرب يرتبط أقوى.", "Multiplication binds more tightly.")},
				map[string]any{"expression": "4 + 6", "value": "10", "caption": ch20Localized("اجمع بعد تقييم الضرب.", "Add after evaluating multiplication.")},
				map[string]any{"expression": "10", "value": "10", "caption": ch20Localized("الـexpression رجّع الـvalue.", "The expression returns a value.")},
			},
			"labels": map[string]any{
				"expression": ch20Localized("التعبير", "Expression"),
				"value":      ch20Localized("الـvalue", "Value"),
				"previous":   ch20Localized("السابق", "Previous"),
				"next":       ch20Localized("التالي", "Next"),
				"reset":      ch20Localized("من الأول", "Reset"),
			},
		}),
	}},
	{ID: "c204.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch20Markdown("c204.watch", "- Parse the whole program before moving; a later syntax error must not affect the world.\n- A runtime error stops later statements, while completed movement remains.")}},
	{ID: "c204.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch20Markdown("c204.plan", "Separate the phases: parse statements and expressions, evaluate the AST, then call ~FlightDrive.Step~.\n\nCheck that an error points to the right token or operator before you run it."),
		ch20Question("c204.parse-first", "When can the first movement call happen?", []string{"During parsing after each statement", "After the whole program parses", "Before reading tokens"}, 1, []string{"That can move the drone before a later syntax error is found.", "A full AST separates syntax from side effects.", "The grammar must be read before execution."}),
	}},
	{ID: "c204.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch20Markdown("c204.turn", "Start with a small expression, then try let and thrust. Check errors and completed movement separately.")}},
}
