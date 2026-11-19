package levels

type ch32Pair struct{ ar, en string }

type ch32FrameSpec struct {
	line    int
	caption ch32Pair
	vars    []LessonVariable
	output  string
}

type ch32WalkSpec struct {
	text      ch32Pair
	code      string
	exampleID string
	imports   []string
	frames    []ch32FrameSpec
}

type ch32QuestionSpec struct {
	text     ch32Pair
	options  []ch32Pair
	answer   int
	feedback []ch32Pair
}

type ch32Language uint8

const (
	ch32Arabic ch32Language = iota
	ch32English
)

type ch32LessonSpec struct {
	prefix       string
	mission      ch32Pair
	idea         ch32WalkSpec
	tryPrompt    ch32Pair
	tryCode      string
	tryExampleID string
	tryImports   []string
	tryFrames    []ch32FrameSpec
	tryOptions   []ch32Pair
	tryAnswer    int
	tryFeedback  []ch32Pair
	watch        ch32Pair
	plan         ch32Pair
	check        ch32QuestionSpec
	turn         ch32Pair
	hints        [3]ch32Pair
}

func ch32Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch32Question(id string, spec ch32QuestionSpec, language ch32Language) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: language.text(spec.text), Question: &LessonQuestion{Options: ch32Strings(spec.options, language), Answer: spec.answer, Feedback: ch32Strings(spec.feedback, language)}}
}

func ch32Walkthrough(id string, spec ch32WalkSpec, language ch32Language) LessonBlock {
	frames := make([]LessonFrame, 0, len(spec.frames))
	for _, frame := range spec.frames {
		frames = append(frames, LessonFrame{Line: frame.line, Caption: language.text(frame.caption), Variables: frame.vars, Output: frame.output})
	}
	return LessonBlock{ID: id, Kind: "walkthrough", Text: language.text(spec.text), Code: spec.code, Imports: spec.imports, ExampleID: spec.exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func (language ch32Language) text(value ch32Pair) string {
	if language == ch32English {
		return value.en
	}
	return value.ar
}

func ch32Strings(values []ch32Pair, language ch32Language) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, language.text(value))
	}
	return result
}

func ch32BuildLesson(spec ch32LessonSpec) (arabic, english []LessonStep, arabicHints, englishHints []string) {
	arabic = ch32BuildLocalizedLesson(spec, ch32Arabic)
	english = ch32BuildLocalizedLesson(spec, ch32English)
	arabicHints = []string{spec.hints[0].ar, spec.hints[1].ar, spec.hints[2].ar}
	englishHints = []string{spec.hints[0].en, spec.hints[1].en, spec.hints[2].en}
	return arabic, english, arabicHints, englishHints
}

func ch32BuildLocalizedLesson(spec ch32LessonSpec, language ch32Language) []LessonStep {
	return []LessonStep{
		ch32MissionStep(spec, language),
		ch32IdeaStep(spec, language),
		ch32TryStep(spec, language),
		ch32WatchStep(spec, language),
		ch32CheckStep(spec, language),
		ch32TurnStep(spec, language),
	}
}

func ch32MissionStep(spec ch32LessonSpec, language ch32Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".mission-step", Title: language.text(ch32Pair{ar: "المهمة", en: "Mission"}), Blocks: []LessonBlock{ch32Markdown(spec.prefix+".mission", language.text(spec.mission))}}
}

func ch32IdeaStep(spec ch32LessonSpec, language ch32Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".idea-step", Title: language.text(ch32Pair{ar: "الفكرة", en: "Idea"}), Blocks: []LessonBlock{ch32Walkthrough(spec.prefix+".idea", spec.idea, language)}}
}

func ch32TryStep(spec ch32LessonSpec, language ch32Language) LessonStep {
	question := ch32QuestionSpec{text: ch32Pair{ar: spec.tryPrompt.ar + "\n\n~~~go\n" + spec.tryCode + "\n~~~", en: spec.tryPrompt.en + "\n\n~~~go\n" + spec.tryCode + "\n~~~"}, options: spec.tryOptions, answer: spec.tryAnswer, feedback: spec.tryFeedback}
	reveal := ch32WalkSpec{text: ch32Pair{ar: "تابع التنفيذ سطر بسطر وشوف الناتج.", en: "Follow execution line by line and check the output."}, code: spec.tryCode, exampleID: spec.tryExampleID, imports: spec.tryImports, frames: spec.tryFrames}
	return LessonStep{ID: spec.prefix + ".try-step", Title: language.text(ch32Pair{ar: "جرّب", en: "Try it"}), Blocks: []LessonBlock{ch32Question(spec.prefix+".predict", question, language), ch32Walkthrough(spec.prefix+".reveal", reveal, language)}}
}

func ch32WatchStep(spec ch32LessonSpec, language ch32Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".watch-step", Title: language.text(ch32Pair{ar: "خد بالك", en: "Watch out"}), Blocks: []LessonBlock{ch32Markdown(spec.prefix+".watch", language.text(spec.watch))}}
}

func ch32CheckStep(spec ch32LessonSpec, language ch32Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".check-step", Title: language.text(ch32Pair{ar: "فكّر واختبر", en: "Think & check"}), Blocks: []LessonBlock{ch32Markdown(spec.prefix+".plan", language.text(spec.plan)), ch32Question(spec.prefix+".check", spec.check, language)}}
}

func ch32TurnStep(spec ch32LessonSpec, language ch32Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".turn-step", Title: language.text(ch32Pair{ar: "دورك", en: "Your turn"}), Blocks: []LessonBlock{ch32Markdown(spec.prefix+".turn", language.text(spec.turn))}}
}

var ch32Lesson321Arabic, ch32Lesson321English, ch32Lesson321HintsArabic, ch32Lesson321HintsEnglish = ch32BuildLesson(ch32LessonSpec{
	prefix:  "c321",
	mission: ch32Pair{ar: "راجع عقد ~RelayTelemetry~ في المهمة، وحدد حدود الرسالة وعمر الاتصال.\n\nالـstarter لسه TODO؛ ركّز على فصل القراية عن الكتابة والإغلاق.", en: "Read the ~RelayTelemetry~ contract and note the message and connection boundaries.\n\nThe starter is still TODO; focus on separating reads, writes, and shutdown."},
	idea: ch32WalkSpec{text: ch32Pair{ar: "الرسالة ظرف له نوع ومحتوى. كل دورة تتعامل مع رسالة كاملة، والـpump تفضل مسؤولة عن شغلها.", en: "A message is an envelope with a type and payload. Each cycle handles a whole message, while each pump keeps one responsibility."}, code: `type Packet struct { kind, text string }
func main() {
	p := Packet{"text", "air"}
	fmt.Println(p.kind, p.text)
}`, exampleID: "e321-packet", imports: []string{"fmt"}, frames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch32Pair{ar: "Packet تجمع النوع والمحتوى.", en: "Packet groups type and payload."}},
		{line: 3, caption: ch32Pair{ar: "الرسالة نص عن الهواء.", en: "This message is text about air."}, vars: []LessonVariable{{Name: "p", Value: "{text air}"}}},
		{line: 4, caption: ch32Pair{ar: "اطبع النوع والمحتوى.", en: "Print type and payload."}, vars: []LessonVariable{{Name: "p", Value: "{text air}"}}, output: "text air"},
	}},
	tryPrompt: ch32Pair{ar: "توقّع الناتج للرسالتين.", en: "Predict the output for both messages."},
	tryCode: `type Packet struct { kind, text string }
func main() {
	queue := []Packet{{"text", "air"}, {"binary", "map"}}
	for _, p := range queue {
		fmt.Println(p.kind, p.text)
	}
}`, tryExampleID: "e321-queue", tryImports: []string{"fmt"}, tryFrames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch32Pair{ar: "Packet تحمل وصف الرسالة.", en: "Packet carries the message label."}},
		{line: 3, caption: ch32Pair{ar: "الطابور فيه رسالتين.", en: "The queue has two messages."}, vars: []LessonVariable{{Name: "queue", Value: "[{text air} {binary map}]"}}},
		{line: 4, caption: ch32Pair{ar: "بدأت أول لفة.", en: "The first iteration starts."}, vars: []LessonVariable{{Name: "p", Value: "{text air}"}}, output: ""},
		{line: 5, caption: ch32Pair{ar: "اطبع الرسالة الأولى.", en: "Print the first message."}, vars: []LessonVariable{{Name: "p", Value: "{text air}"}}, output: "text air"},
		{line: 4, caption: ch32Pair{ar: "وصلت للرسالة التانية.", en: "The next message is selected."}, vars: []LessonVariable{{Name: "p", Value: "{binary map}"}}, output: "text air"},
		{line: 5, caption: ch32Pair{ar: "اطبع الرسالة التانية.", en: "Print the second message."}, vars: []LessonVariable{{Name: "p", Value: "{binary map}"}}, output: "text air\nbinary map"},
	},
	tryOptions: []ch32Pair{{ar: "text air وبعدها binary map", en: "text air, then binary map"}, {ar: "binary map بس", en: "binary map only"}, {ar: "ترتيب عكسي", en: "Reverse order"}}, tryAnswer: 0,
	tryFeedback: []ch32Pair{{ar: "range ماشية بترتيب الـslice.", en: "range follows slice order."}, {ar: "اللفة بتزور العنصرين.", en: "The loop visits both elements."}, {ar: "مفيش عكس للترتيب هنا.", en: "Nothing reverses the order here."}},
	watch:       ch32Pair{ar: "- خليك واضح مين بيملك كل نوع I/O.\n- الإلغاء لازم يفك الطرف اللي مستني.", en: "- Keep one clear owner for each kind of I/O.\n- Cancellation must release whichever side is waiting."},
	plan:        ch32Pair{ar: "قسّم التفكير: قبول آمن، قراية، كتابة، وبعدين خروج.\nاسأل نفسك إيه اللي يحصل لو الطرف التاني اختفى فجأة.", en: "Separate the work: safe acceptance, reading, writing, then shutdown.\nAsk what happens when the peer disappears."},
	check:       ch32QuestionSpec{text: ch32Pair{ar: "إيه اللي يخلي pump ما تفضلش معلّقة بعد الإلغاء؟", en: "What helps a pump avoid getting stranded after cancellation?"}, options: []ch32Pair{{ar: "تختار بين الشغل وسياق الإلغاء", en: "Select between work and context cancellation"}, {ar: "تتجاهل أخطاء القراية", en: "Ignore read errors"}, {ar: "تفتح goroutine لكل بايت", en: "Start a goroutine per byte"}}, answer: 0, feedback: []ch32Pair{{ar: "سياق الإلغاء بيدي طريق خروج.", en: "Cancellation gives the loop an exit."}, {ar: "تجاهل الخطأ يخلي الحلقة تدور بعد الفصل.", en: "Ignoring errors can spin after disconnect."}, {ar: "تجزئة المسؤولية تكفي؛ goroutine لكل بايت زيادة.", en: "Separate responsibilities; a goroutine per byte is excessive."}}},
	turn:        ch32Pair{ar: "ابدأ من عقد الرسالة، وبعدها تتبّع مين بيقرا ومين بيكتب.\nشغّل التستات وراقب كل مسار خروج.", en: "Start from the message contract, then trace who reads and who writes.\nRun the tests and inspect each shutdown path."},
	hints:       [3]ch32Pair{{ar: "رسالة WebSocket لها نوع ومحتوى؛ خليك متعامل معاها كقطعة واحدة.", en: "A WebSocket message has a type and payload; handle it as one unit."}, {ar: "افصل مسار القراية عن مسار الكتابة وارسم نقطة الإلغاء بينهم.", en: "Separate read and write paths, then mark where cancellation connects them."}, {ar: "راجع سلوك ~Read~ و~Write~ وسياق ~Ping~ في توثيق الحزمة.", en: "Review ~Read~, ~Write~, and the ~Ping~ context in the package docs."}},
})

var ch32Lesson322Arabic, ch32Lesson322English, ch32Lesson322HintsArabic, ch32Lesson322HintsEnglish = ch32BuildLesson(ch32LessonSpec{
	prefix:  "c322",
	mission: ch32Pair{ar: "المهمة بتحوّل طلب قطاع لقراءة واحدة ورد واحد.\nافصل رفض الطلب، فشل المخزن، والنجاح؛ كل واحد له معنى مختلف للعميل.", en: "The mission turns one sector request into one reading and one reply.\nSeparate rejected input, store failure, and success; each means something different to the caller."},
	idea: ch32WalkSpec{text: ch32Pair{ar: "Unary RPC بتاخد message واحدة وترجع message واحدة.\nالـstruct بتوضح إيه اللي داخل وإيه اللي راجع.", en: "A unary RPC takes one message and returns one message.\nThe structs make the input and result explicit."}, code: `type Query struct { Sector string }
type Reply struct { Sector string; Pressure int64 }
func inspect(q Query) Reply { return Reply{q.Sector, -3} }
func main() {
	r := inspect(Query{"ring"})
	fmt.Println(r.Sector, r.Pressure)
}`, exampleID: "e322-unary-shape", imports: []string{"fmt"}, frames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch32Pair{ar: "Query تمثل الطلب.", en: "Query represents the request."}},
		{line: 2, caption: ch32Pair{ar: "Reply تمثل النتيجة.", en: "Reply represents the result."}},
		{line: 5, caption: ch32Pair{ar: "الطلب طلب قطاع ring.", en: "The query names sector ring."}, vars: []LessonVariable{{Name: "r", Value: "{ring -3}"}}},
		{line: 6, caption: ch32Pair{ar: "اطبع القطاع والضغط.", en: "Print sector and pressure."}, vars: []LessonVariable{{Name: "r", Value: "{ring -3}"}}, output: "ring -3"},
	}},
	tryPrompt: ch32Pair{ar: "توقّع قيمة ~reply.Pressure~.", en: "Predict the value of ~reply.Pressure~."},
	tryCode: `type Reply struct { Pressure int64 }
func main() {
	reply := Reply{Pressure: 0}
	fmt.Println(reply.Pressure)
}`, tryExampleID: "e322-zero-reading", tryImports: []string{"fmt"}, tryFrames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch32Pair{ar: "Reply فيها ضغط int64.", en: "Reply stores int64 pressure."}},
		{line: 3, caption: ch32Pair{ar: "القراءة صفر صالحة.", en: "A zero reading is valid."}, vars: []LessonVariable{{Name: "reply.Pressure", Value: "0"}}},
		{line: 4, caption: ch32Pair{ar: "اطبع القراءة زي ما هي.", en: "Print the reading unchanged."}, vars: []LessonVariable{{Name: "reply.Pressure", Value: "0"}}, output: "0"},
	},
	tryOptions: []ch32Pair{{ar: "0", en: "0"}, {ar: "-3", en: "-3"}, {ar: "البرنامج مش هيطبع حاجة", en: "The program prints nothing"}}, tryAnswer: 0,
	tryFeedback: []ch32Pair{{ar: "القيمة الصفرية موجودة في الحقل.", en: "The field contains zero."}, {ar: "-3 كانت في مثال تاني.", en: "-3 belonged to another example."}, {ar: "fmt.Println بتطبع الحقل.", en: "fmt.Println prints the field."}},
	watch:       ch32Pair{ar: "- status code جزء من API العام، مش نص خطأ عشوائي.\n- متطلعش تفاصيل المخزن الداخلية للعميل.", en: "- A status code is part of the public API, not an arbitrary error string.\n- Keep private store details away from callers."},
	plan:        ch32Pair{ar: "ابدأ بعقد الرسائل المسجّلة، وبعدها خدمة مستقلة.\nراجع حدود الطلب قبل البحث وحدد تحويل الأخطاء عند الحافة.", en: "Start with the registered message contract, then the service.\nValidate before lookup and decide how errors map at the boundary."},
	check:       ch32QuestionSpec{text: ch32Pair{ar: "فين تحط فحص المفتاح المشترك لكل unary call؟", en: "Where should shared key checks for unary calls live?"}, options: []ch32Pair{{ar: "Unary interceptor", en: "A unary interceptor"}, {ar: "داخل كل lookup", en: "Inside every lookup"}, {ar: "بعد ما الرد يوصل للعميل", en: "After the reply reaches the client"}}, answer: 0, feedback: []ch32Pair{{ar: "الـinterceptor بيغلف استدعاءات الخدمة.", en: "The interceptor wraps service calls."}, {ar: "المخزن مسؤول عن البيانات، مش بيانات الاعتماد.", en: "The store owns data, not credentials."}, {ar: "الفحص لازم يوقف الطلب قبل الخدمة.", en: "The check must stop the call before the service."}}},
	turn:        ch32Pair{ar: "اكتب الخدمة حول عقد الـproto المقفول.\nاختبر النجاح والرفض وأخطاء السياق لوحدهم.", en: "Build the service around the read-only proto contract.\nCheck success, rejection, and context errors separately."},
	hints:       [3]ch32Pair{{ar: "Unary RPC بيربط طلب typed واحد بنتيجة typed واحدة.", en: "A unary RPC pairs one typed request with one typed result."}, {ar: "خلّي حالة الطلب محلية، وافصل اعتراض التصريح عن منطق الخدمة.", en: "Keep request state local and separate authorization from service logic."}, {ar: "ابدأ من ~BerthControlServer~ و~UnaryServerInterceptor~ المولّدين.", en: "Start with the generated ~BerthControlServer~ and ~UnaryServerInterceptor~ APIs."}},
})

var ch32Lesson323Arabic, ch32Lesson323English, ch32Lesson323HintsArabic, ch32Lesson323HintsEnglish = ch32BuildLesson(ch32LessonSpec{
	prefix:  "c323",
	mission: ch32Pair{ar: "كل stream لها خطوات مرتبة وحالتها الخاصة.\nاتبع عقد الرسائل وخلي الإقرار مربوط بالخطوة اللي استلمتها.", en: "Each stream has ordered steps and its own state.\nFollow the message contract and keep each acknowledgement tied to its step."},
	idea: ch32WalkSpec{text: ch32Pair{ar: "الـstream سلسلة رسائل، مش slice كاملة لازم تتجمع الأول.\nاقرا خطوة، وخلي رقمها ظاهر في التتبّع.", en: "A stream is a sequence of messages, not a whole slice that must be collected first.\nKeep each step number visible as you trace it."}, code: `type Step struct { Sequence, Delta int }
func main() {
	steps := []Step{{1, 4}, {2, -2}}
	for _, step := range steps {
		fmt.Println(step.Sequence, step.Delta)
	}
}`, exampleID: "e323-ordered-steps", imports: []string{"fmt"}, frames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch32Pair{ar: "كل Step لها sequence وdelta.", en: "Each Step has a sequence and delta."}},
		{line: 3, caption: ch32Pair{ar: "وصلت خطوتين للتوضيح.", en: "Two illustrative steps are ready."}, vars: []LessonVariable{{Name: "steps", Value: "[{1 4} {2 -2}]"}}},
		{line: 4, caption: ch32Pair{ar: "الخطوة الأولى رقمها 1.", en: "The first step is numbered 1."}, vars: []LessonVariable{{Name: "step", Value: "{1 4}"}}},
		{line: 5, caption: ch32Pair{ar: "اطبع رقمها وتغييرها.", en: "Print its number and delta."}, vars: []LessonVariable{{Name: "step", Value: "{1 4}"}}, output: "1 4"},
		{line: 4, caption: ch32Pair{ar: "الخطوة اللي بعدها رقم 2.", en: "The next step is numbered 2."}, vars: []LessonVariable{{Name: "step", Value: "{2 -2}"}}, output: "1 4"},
		{line: 5, caption: ch32Pair{ar: "اطبع الزوج التاني.", en: "Print the second pair."}, vars: []LessonVariable{{Name: "step", Value: "{2 -2}"}}, output: "1 4\n2 -2"},
	}},
	tryPrompt: ch32Pair{ar: "توقّع السطرين اللي هيتطبعوا.", en: "Predict the two printed lines."},
	tryCode: `func main() {
	for _, n := range []int{3, 5} {
		fmt.Println(n)
	}
}`, tryExampleID: "e323-range-output", tryImports: []string{"fmt"}, tryFrames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch32Pair{ar: "القيمة الأولى 3.", en: "The first value is 3."}, vars: []LessonVariable{{Name: "n", Value: "3"}}},
		{line: 3, caption: ch32Pair{ar: "اطبع 3.", en: "Print 3."}, vars: []LessonVariable{{Name: "n", Value: "3"}}, output: "3"},
		{line: 2, caption: ch32Pair{ar: "بعدها القيمة 5.", en: "Next, the value is 5."}, vars: []LessonVariable{{Name: "n", Value: "5"}}, output: "3"},
		{line: 3, caption: ch32Pair{ar: "اطبع 5.", en: "Print 5."}, vars: []LessonVariable{{Name: "n", Value: "5"}}, output: "3\n5"},
	},
	tryOptions: []ch32Pair{{ar: "3 وبعدها 5", en: "3, then 5"}, {ar: "5 وبعدها 3", en: "5, then 3"}, {ar: "3 بس", en: "3 only"}}, tryAnswer: 0,
	tryFeedback: []ch32Pair{{ar: "range بتمشي بترتيب العناصر.", en: "range visits elements in order."}, {ar: "مفيش ترتيب عكسي في الكود.", en: "The code does not reverse the order."}, {ar: "الـslice فيها عنصرين.", en: "The slice has two elements."}},
	watch:       ch32Pair{ar: "- حالة كل stream مكانها جوه الاستدعاء.\n- EOF معناه نهاية إرسال طبيعية من العميل.", en: "- Each stream's state belongs to its call.\n- EOF means the client finished sending normally."},
	plan:        ch32Pair{ar: "تابع حدود الرسالة، حالة stream، والرد المقابل.\nخلي الفشل يوقف الخطوات اللي بعده.", en: "Track the message boundary, stream state, and matching reply.\nLet a failure stop later steps."},
	check:       ch32QuestionSpec{text: ch32Pair{ar: "ليه حالة الإجمالي تبقى محلية لكل stream؟", en: "Why keep the running total local to each stream?"}, options: []ch32Pair{{ar: "عشان streams متخلطش بياناتها", en: "So streams do not mix their data"}, {ar: "عشان نقدر نعيد نفس الإجمالي لكل crew", en: "So every crew can reuse one total"}, {ar: "عشان نتجاهل ترتيب الخطوات", en: "So step order can be ignored"}}, answer: 0, feedback: []ch32Pair{{ar: "كل استدعاء له حالة مستقلة.", en: "Each call has independent state."}, {ar: "المشاركة تخلط النتائج بين العملاء.", en: "Sharing can mix results between clients."}, {ar: "الترتيب جزء من معنى الرسائل.", en: "Order is part of the message meaning."}}},
	turn:        ch32Pair{ar: "ابدأ بمسار stream واحد، وبعدها جرّب أكتر من stream.\nراجع الإرسال والإغلاق مع بعض.", en: "Start with one stream path, then check simultaneous streams.\nReview sending and closure together."},
	hints:       [3]ch32Pair{{ar: "Bidirectional RPC بيمرر رسائل في الاتجاهين على stream واحدة.", en: "A bidirectional RPC carries messages both ways on one stream."}, {ar: "افصل حالة كل استدعاء عن حالة السيرفر المشتركة.", en: "Separate per-call state from shared server state."}, {ar: "راجع سلوك ~Recv~ و~Send~ و~Context~ في interface المولّد.", en: "Review ~Recv~, ~Send~, and ~Context~ on the generated stream interface."}},
})

var ch32Lesson324Arabic, ch32Lesson324English, ch32Lesson324HintsArabic, ch32Lesson324HintsEnglish = ch32BuildLesson(ch32LessonSpec{
	prefix:  "c324",
	mission: ch32Pair{ar: "المهمة بتبني إعدادين: واحد للسيرفر وواحد للعميل.\n\nفكّر مين يعرض الشهادة، ومين يقرر يثق في الطرف التاني.\n\nساعة الاختبار تختار وقت التحقق من غير انتظار حقيقي.", en: "The mission builds two configurations: one for the server and one for the client.\n\nAsk who presents a certificate and who decides to trust the other peer.\n\nThe injected clock selects verification time without real waiting."},
	idea: ch32WalkSpec{text: ch32Pair{ar: "إعداد TLS بيحمل توقعات كل طرف. ~ServerName~ اسم الهوية اللي العميل بيتحقق منها، مش عنوان اتصال.\nكل peer له شهادة أو مجموعة ثقة حسب دوره.", en: "A TLS config carries each peer's expectations. ~ServerName~ is the identity the client verifies, not a dial address.\nEach peer presents or trusts certificates according to its role."}, code: `func main() {
	client := &tls.Config{ServerName: "probe.colony", MinVersion: tls.VersionTLS13}
	fmt.Println(client.ServerName, client.MinVersion == tls.VersionTLS13)
}`, exampleID: "e324-client-expectation", imports: []string{"crypto/tls", "fmt"}, frames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch32Pair{ar: "العميل يتوقع اسم الهوية.", en: "The client expects this identity."}, vars: []LessonVariable{{Name: "client.ServerName", Value: "probe.colony"}, {Name: "client.MinVersion", Value: "772"}}},
		{line: 3, caption: ch32Pair{ar: "اعرض الاسم وفحص النسخة.", en: "Print the name and version check."}, vars: []LessonVariable{{Name: "client.ServerName", Value: "probe.colony"}}, output: "probe.colony true"},
	}},
	tryPrompt: ch32Pair{ar: "توقّع ناتج نسخة TLS.", en: "Predict the TLS version check."},
	tryCode: `func main() {
	config := &tls.Config{MinVersion: tls.VersionTLS13}
	fmt.Println(config.MinVersion == tls.VersionTLS13)
}`, tryExampleID: "e324-version-check", tryImports: []string{"crypto/tls", "fmt"}, tryFrames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch32Pair{ar: "الإعداد يبدأ بـ TLS 1.3.", en: "The config starts at TLS 1.3."}, vars: []LessonVariable{{Name: "config.MinVersion", Value: "772"}}},
		{line: 3, caption: ch32Pair{ar: "قارن النسخة المطلوبة.", en: "Compare the requested version."}, output: "true"},
	},
	tryOptions: []ch32Pair{{ar: "true", en: "true"}, {ar: "false", en: "false"}, {ar: "0", en: "0"}}, tryAnswer: 0,
	tryFeedback: []ch32Pair{{ar: "القيمتين هما TLS 1.3.", en: "Both values are TLS 1.3."}, {ar: "مافيش قيمة تانية اتطلبت.", en: "No other version was requested."}, {ar: "المقارنة بترجع bool.", en: "The comparison returns a bool."}},
	watch:       ch32Pair{ar: "- التشفير لوحده مش إثبات هوية.\n- سيب التحقق القياسي شغال ومتعدّليش مدخلات المستدعي.\nنسخ trust pool يحافظ على ملكية المستدعي.", en: "- Encryption alone does not prove identity.\n- Keep standard verification enabled and leave caller inputs untouched.\nCloning a trust pool protects caller ownership."},
	plan:        ch32Pair{ar: "افصل إعداد السيرفر عن العميل، وبعدها راجع اتجاه الثقة.\nاختبر الإعداد بالمصافحة اللي تناسبه، وقارن الاسم والوقت والجهة الموثوقة.", en: "Separate server and client settings, then check the direction of trust.\nUse a matching handshake and compare identity, time, and trust."},
	check:       ch32QuestionSpec{text: ch32Pair{ar: "~ServerName~ عند العميل بيحدد إيه؟", en: "What does client ~ServerName~ identify?"}, options: []ch32Pair{{ar: "اسم السيرفر المتوقع", en: "The expected server identity"}, {ar: "عنوان IP لازم يتصل به", en: "The IP address to dial"}, {ar: "اسم ملف الشهادة", en: "The certificate file name"}}, answer: 0, feedback: []ch32Pair{{ar: "TLS بيقارن هوية الشهادة بالاسم.", en: "TLS checks the certificate identity against the name."}, {ar: "الـconnection بيتحدد برّه الحقل ده.", en: "The connection target is chosen elsewhere."}, {ar: "الحقل مش بيحمّل شهادة.", en: "This field does not load a certificate."}}},
	turn:        ch32Pair{ar: "ابنِ الإعدادات كقيم جديدة قبل المصافحة.\nخلي كل فشل يوضح جهة الثقة اللي اترفضت.", en: "Build fresh config values before the handshake.\nIdentify which trust check rejected each failure."},
	hints:       [3]ch32Pair{{ar: "TLS بيربط التشفير بهوية شهادة موثوقة.", en: "TLS combines encryption with a trusted certificate identity."}, {ar: "فكّر في إعداد مستقل لكل peer وفي نسخ مجموعات الثقة.", en: "Plan a separate config per peer and clone trust pools."}, {ar: "راجع حقول ~tls.Config~ زي ~RootCAs~ و~ClientCAs~ و~ServerName~.", en: "Review ~tls.Config~ fields such as ~RootCAs~, ~ClientCAs~, and ~ServerName~."}},
})

var ch32Lesson325Arabic, ch32Lesson325English, ch32Lesson325HintsArabic, ch32Lesson325HintsEnglish = ch32BuildLesson(ch32LessonSpec{
	prefix:  "c325",
	mission: ch32Pair{ar: "الخدمة بتقرأ من أكتر من طلب على نفس وصلة TLS.\nركز على تفاوض البروتوكول وعلى إن قراءة واحدة واقفة ما توقفش الباقي.", en: "The service reads multiple requests over one TLS connection.\nFocus on protocol negotiation and keeping one blocked read from stopping its siblings."},
	idea: ch32WalkSpec{text: ch32Pair{ar: "ALPN بيختار بروتوكول التطبيق أثناء TLS. وجود ~h2~ في القائمة عرض، والاتصال نفسه يثبت إيه اللي اتفاوض عليه.", en: "ALPN selects an application protocol during TLS. Offering ~h2~ is an offer; the connection reveals what was negotiated."}, code: `func main() {
	config := &tls.Config{NextProtos: []string{"h2", "http/1.1"}}
	fmt.Println(config.NextProtos[0], len(config.NextProtos))
}`, exampleID: "e325-alpn-offer", imports: []string{"crypto/tls", "fmt"}, frames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch32Pair{ar: "القائمة تفضّل h2.", en: "The list prefers h2."}, vars: []LessonVariable{{Name: "config.NextProtos", Value: "[h2 http/1.1]"}}},
		{line: 3, caption: ch32Pair{ar: "اطبع أول اختيار وعددهم.", en: "Print the first choice and count."}, output: "h2 2"},
	}},
	tryPrompt: ch32Pair{ar: "توقّع عدد البروتوكولات المعروضة.", en: "Predict how many protocols are offered."},
	tryCode: `func main() {
	config := &tls.Config{NextProtos: []string{"h2", "http/1.1"}}
	fmt.Println(len(config.NextProtos))
}`, tryExampleID: "e325-alpn-count", tryImports: []string{"crypto/tls", "fmt"}, tryFrames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch32Pair{ar: "فيه اختيارين.", en: "There are two choices."}, vars: []LessonVariable{{Name: "config.NextProtos", Value: "[h2 http/1.1]"}}},
		{line: 3, caption: ch32Pair{ar: "len تعدّهم.", en: "len counts them."}, output: "2"},
	},
	tryOptions: []ch32Pair{{ar: "2", en: "2"}, {ar: "1", en: "1"}, {ar: "h2", en: "h2"}}, tryAnswer: 0,
	tryFeedback: []ch32Pair{{ar: "القائمة فيها h2 وhttp/1.1.", en: "The list has h2 and http/1.1."}, {ar: "الـslice فيها قيمتين.", en: "The slice has two values."}, {ar: "len بترجع عدد العناصر.", en: "len returns the element count."}},
	watch:       ch32Pair{ar: "- متغيّرش TLS config اللي جالك؛ Clone وبعدين عدّل النسخة.\n- القفل المشترك حوالين قراءة الشبكة يعطّل الطلبات الشقيقة.\nنسخة NextProtos بتحافظ على إعداد المستدعي.", en: "- Do not mutate the supplied TLS config; clone before changing it.\n- A shared lock around network reads blocks sibling requests.\nA fresh NextProtos slice protects caller settings."},
	plan:        ch32Pair{ar: "راجع طبقة TLS، handler، ثم القراءات المتزامنة.\nاستخدم نتيجة ALPN وهوية الاتصال كدليل، مش الوقت.\nالحاجز يوضح التوازي من غير تخمين السرعة.", en: "Review TLS, the handler, then concurrent reads.\nUse ALPN and connection identity as evidence, not elapsed time.\nA barrier proves concurrency without timing guesses."},
	check:       ch32QuestionSpec{text: ch32Pair{ar: "إيه اللي يثبت إن الاتصال اتكلم HTTP/2 فعلًا؟", en: "What proves the connection actually spoke HTTP/2?"}, options: []ch32Pair{{ar: "البروتوكول المتفاوض عليه في TLS والـresponse", en: "The negotiated TLS protocol and response"}, {ar: "وجود h2 في القائمة بس", en: "Only that h2 appears in the offer"}, {ar: "إن الطلب خلص بسرعة", en: "That the request finished quickly"}}, answer: 0, feedback: []ch32Pair{{ar: "التفاوض الفعلي أقوى من عرض البروتوكول.", en: "The actual negotiation is stronger evidence than the offer."}, {ar: "القائمة وحدها ما تثبتش الاختيار.", en: "The offer alone does not prove selection."}, {ar: "الوقت مش دليل على نسخة البروتوكول.", en: "Timing does not establish the protocol version."}}},
	turn:        ch32Pair{ar: "ابنِ handler صغير الأول، وبعدها تأكد من التفاوض والتوازي.\nراجع إن طلب جانبي يكمل لو قراءة اتأخرت.", en: "Build the small handler first, then verify negotiation and concurrency.\nCheck that a sibling request finishes while another read waits."},
	hints:       [3]ch32Pair{{ar: "HTTP/2 بيحمل streams متعددة على اتصال واحد.", en: "HTTP/2 carries multiple streams over one connection."}, {ar: "افصل فحص المسار والطريقة عن قراءة الحساس، وسيب كل طلب مستقل.", en: "Separate route and method checks from sensor reads; keep requests independent."}, {ar: "راجع ~tls.Config.Clone~ و~NextProtos~ وإعدادات ~http.Server~.", en: "Review ~tls.Config.Clone~, ~NextProtos~, and ~http.Server~ settings."}},
})

var ch32Lesson326Arabic, ch32Lesson326English, ch32Lesson326HintsArabic, ch32Lesson326HintsEnglish = ch32BuildLesson(ch32LessonSpec{
	prefix:  "c326",
	mission: ch32Pair{ar: "الوصلة بتنفّذ أوامر فحص محددة بعد اتصال موثوق.\nاربط الاختيار بهوية السيرفر، وخلي الإلغاء ينهي الموارد التابعة للطلب.", en: "The link runs a small set of inspection commands over a trusted connection.\nTie command selection to server identity, and let cancellation release request resources."},
	idea: ch32WalkSpec{text: ch32Pair{ar: "جدول allowlist بيربط اختيار معروف بفعل ثابت.\nده يمنع إن اسم جاي من المستخدم يبقى أمر حر.", en: "An allowlist maps a known choice to a fixed action.\nThat keeps user input from becoming an arbitrary command."}, code: `type Command struct { label, text string }
func main() {
	commands := map[string]Command{"inspect": {"inspect", "show queue"}}
	c, ok := commands["inspect"]
	fmt.Println(ok, c.text)
}`, exampleID: "e326-command-map", imports: []string{"fmt"}, frames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch32Pair{ar: "Command تجمع الاختيار والفعل.", en: "Command groups a choice and action."}},
		{line: 3, caption: ch32Pair{ar: "الخريطة فيها أمر توضيحي.", en: "The map has an illustrative command."}, vars: []LessonVariable{{Name: "commands", Value: "map[inspect:{inspect show queue}]"}}},
		{line: 4, caption: ch32Pair{ar: "الاختيار موجود.", en: "The choice exists."}, vars: []LessonVariable{{Name: "ok", Value: "true"}, {Name: "c.text", Value: "show queue"}}},
		{line: 5, caption: ch32Pair{ar: "اطبع النتيجة الثابتة.", en: "Print the fixed action."}, output: "true show queue"},
	}},
	tryPrompt: ch32Pair{ar: "توقّع نتيجة البحث عن المفتاح.", en: "Predict the lookup result."},
	tryCode: `func main() {
	commands := map[string]string{"scan": "read status"}
	text, ok := commands["scan"]
	fmt.Println(ok, text)
}`, tryExampleID: "e326-allowlist-lookup", tryImports: []string{"fmt"}, tryFrames: []ch32FrameSpec{
		{line: 0, caption: ch32Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch32Pair{ar: "الخريطة فيها scan.", en: "The map contains scan."}, vars: []LessonVariable{{Name: "commands", Value: "map[scan:read status]"}}},
		{line: 3, caption: ch32Pair{ar: "المفتاح معروف.", en: "The key is known."}, vars: []LessonVariable{{Name: "ok", Value: "true"}, {Name: "text", Value: "read status"}}},
		{line: 4, caption: ch32Pair{ar: "اطبع وجوده والنص.", en: "Print presence and action."}, output: "true read status"},
	},
	tryOptions: []ch32Pair{{ar: "true read status", en: "true read status"}, {ar: "false <فاضي>", en: "false <empty>"}, {ar: "panic", en: "panic"}}, tryAnswer: 0,
	tryFeedback: []ch32Pair{{ar: "الـlookup رجّع القيمة وtrue.", en: "The lookup returned the value and true."}, {ar: "المفتاح موجود في الخريطة.", en: "The key is present in the map."}, {ar: "قراءة المفتاح الموجود آمنة.", en: "Reading an existing key is safe."}},
	watch:       ch32Pair{ar: "- متبعتش نص أمر حر من الطلب.\n- راقب الإلغاء وأخطاء الجلسة والاتصال.\nخلي stdout وstderr منفصلين عشان تراجع النتيجة.", en: "- Never pass free-form command text from the request.\n- Track cancellation and session or connection errors.\nKeep stdout and stderr separate when reviewing the result."},
	plan:        ch32Pair{ar: "افصل التحقق قبل الاتصال عن تنفيذ الجلسة.\nراجع ملكية الاتصال والـsession وكل مخرج لوحده.", en: "Separate validation before connecting from session execution.\nReview ownership of the connection, session, and each output."},
	check:       ch32QuestionSpec{text: ch32Pair{ar: "ليه نستخدم host key موثوق؟", en: "Why use a trusted host key?"}, options: []ch32Pair{{ar: "عشان نعرف هوية السيرفر", en: "To verify the server identity"}, {ar: "عشان نختار نص الأمر", en: "To choose command text"}, {ar: "عشان نتجاهل الإلغاء", en: "To ignore cancellation"}}, answer: 0, feedback: []ch32Pair{{ar: "التحقق يمنع قبول endpoint متنكّر.", en: "Verification rejects an impersonating endpoint."}, {ar: "الأمر مصدره اختيار مسموح، مش host key.", en: "The command comes from an allowed choice, not the host key."}, {ar: "الإلغاء يفضل جزء من دورة الاتصال.", en: "Cancellation remains part of connection lifetime."}}},
	turn:        ch32Pair{ar: "راجع الـenum والعقد قبل فتح اتصال.\nجرّب حالات الإلغاء ومخارج الأمر كل حالة لوحدها.", en: "Review the enum and contract before opening a connection.\nCheck cancellation and command outcomes separately."},
	hints:       [3]ch32Pair{{ar: "SSH host key بيثبت هوية الطرف البعيد؛ التشفير وحده مش كفاية.", en: "An SSH host key verifies the remote peer; encryption alone is not enough."}, {ar: "ارسم عمر الاتصال والجلسة وراقب مين يقفل كل مورد.", en: "Trace connection and session lifetime, including who closes each resource."}, {ar: "راجع ~ssh.FixedHostKey~ و~ssh.NewClientConn~ وسياق الطلب.", en: "Review ~ssh.FixedHostKey~, ~ssh.NewClientConn~, and the request context."}},
})
