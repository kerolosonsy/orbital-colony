package levels

type ch33Pair struct{ ar, en string }

type ch33FrameSpec struct {
	line    int
	caption ch33Pair
	vars    []LessonVariable
	output  string
}

type ch33WalkSpec struct {
	text      ch33Pair
	code      string
	exampleID string
	imports   []string
	frames    []ch33FrameSpec
}

type ch33QuestionSpec struct {
	text     ch33Pair
	options  []ch33Pair
	answer   int
	feedback []ch33Pair
}

type ch33SandboxFrame struct {
	caption                ch33Pair
	leftLabel, leftValue   ch33Pair
	rightLabel, rightValue ch33Pair
}

type ch33SandboxRun struct {
	label  ch33Pair
	frames []ch33SandboxFrame
}

type ch33SandboxSpec struct {
	id, kind string
	text     ch33Pair
	title    ch33Pair
	advance  ch33Pair
	reset    ch33Pair
	repeat   ch33Pair
	runs     []ch33SandboxRun
}

type ch33Language uint8

const (
	ch33Arabic ch33Language = iota
	ch33English
)

type ch33LessonSpec struct {
	prefix       string
	mission      ch33Pair
	idea         ch33WalkSpec
	tryPrompt    ch33Pair
	tryCode      string
	tryExampleID string
	tryImports   []string
	tryFrames    []ch33FrameSpec
	tryOptions   []ch33Pair
	tryAnswer    int
	tryFeedback  []ch33Pair
	watch        ch33Pair
	plan         ch33Pair
	check        ch33QuestionSpec
	turn         ch33Pair
	hints        [3]ch33Pair
	sandbox      ch33SandboxSpec
}

func ch33Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch33Question(id string, spec ch33QuestionSpec, language ch33Language) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: language.text(spec.text), Question: &LessonQuestion{Options: ch33Strings(spec.options, language), Answer: spec.answer, Feedback: ch33Strings(spec.feedback, language)}}
}

func ch33Walkthrough(id string, spec ch33WalkSpec, language ch33Language) LessonBlock {
	frames := make([]LessonFrame, 0, len(spec.frames))
	for _, frame := range spec.frames {
		frames = append(frames, LessonFrame{Line: frame.line, Caption: language.text(frame.caption), Variables: frame.vars, Output: frame.output})
	}
	return LessonBlock{ID: id, Kind: "walkthrough", Text: language.text(spec.text), Code: spec.code, Imports: spec.imports, ExampleID: spec.exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func (language ch33Language) text(value ch33Pair) string {
	if language == ch33English {
		return value.en
	}
	return value.ar
}

func ch33Strings(values []ch33Pair, language ch33Language) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, language.text(value))
	}
	return result
}

func ch33Locale(value ch33Pair) map[string]any {
	return map[string]any{"ar": value.ar, "en": value.en}
}

func ch33SandboxBlock(spec ch33SandboxSpec, language ch33Language) LessonBlock {
	localRuns := make([]any, 0, len(spec.runs))
	for _, run := range spec.runs {
		frames := make([]any, 0, len(run.frames))
		for _, frame := range run.frames {
			frames = append(frames, map[string]any{
				"caption":    ch33Locale(frame.caption),
				"leftLabel":  ch33Locale(frame.leftLabel),
				"leftValue":  ch33Locale(frame.leftValue),
				"rightLabel": ch33Locale(frame.rightLabel),
				"rightValue": ch33Locale(frame.rightValue),
			})
		}
		localRuns = append(localRuns, map[string]any{"label": ch33Locale(run.label), "frames": frames})
	}
	return LessonBlock{
		ID: spec.id, Kind: "sandbox", Text: language.text(spec.text),
		Sandbox: &LessonSandbox{Kind: spec.kind, Config: map[string]any{
			"title": ch33Locale(spec.title),
			"labels": map[string]any{
				"advance": ch33Locale(spec.advance),
				"repeat":  ch33Locale(spec.repeat),
				"reset":   ch33Locale(spec.reset),
			},
			"runs": localRuns,
		}},
	}
}

func ch33BuildLesson(spec ch33LessonSpec) (arabic, english []LessonStep, arabicHints, englishHints []string) {
	arabic = ch33BuildLocalizedLesson(spec, ch33Arabic)
	english = ch33BuildLocalizedLesson(spec, ch33English)
	arabicHints = []string{spec.hints[0].ar, spec.hints[1].ar, spec.hints[2].ar}
	englishHints = []string{spec.hints[0].en, spec.hints[1].en, spec.hints[2].en}
	return arabic, english, arabicHints, englishHints
}

func ch33BuildLocalizedLesson(spec ch33LessonSpec, language ch33Language) []LessonStep {
	return []LessonStep{
		ch33MissionStep(spec, language),
		ch33IdeaStep(spec, language),
		ch33TryStep(spec, language),
		ch33WatchStep(spec, language),
		ch33CheckStep(spec, language),
		ch33TurnStep(spec, language),
	}
}

func ch33MissionStep(spec ch33LessonSpec, language ch33Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".mission-step", Title: language.text(ch33Pair{ar: "المهمة", en: "Mission"}), Blocks: []LessonBlock{ch33Markdown(spec.prefix+".mission", language.text(spec.mission))}}
}

func ch33IdeaStep(spec ch33LessonSpec, language ch33Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".idea-step", Title: language.text(ch33Pair{ar: "الفكرة", en: "Idea"}), Blocks: []LessonBlock{ch33Walkthrough(spec.prefix+".idea", spec.idea, language)}}
}

func ch33TryStep(spec ch33LessonSpec, language ch33Language) LessonStep {
	question := ch33QuestionSpec{text: ch33Pair{ar: spec.tryPrompt.ar + "\n\n~~~go\n" + spec.tryCode + "\n~~~", en: spec.tryPrompt.en + "\n\n~~~go\n" + spec.tryCode + "\n~~~"}, options: spec.tryOptions, answer: spec.tryAnswer, feedback: spec.tryFeedback}
	reveal := ch33WalkSpec{text: ch33Pair{ar: "تابع التنفيذ سطر بسطر وشوف الناتج.", en: "Follow execution line by line and check the output."}, code: spec.tryCode, exampleID: spec.tryExampleID, imports: spec.tryImports, frames: spec.tryFrames}
	return LessonStep{ID: spec.prefix + ".try-step", Title: language.text(ch33Pair{ar: "جرّب", en: "Try it"}), Blocks: []LessonBlock{ch33Question(spec.prefix+".predict", question, language), ch33Walkthrough(spec.prefix+".reveal", reveal, language), ch33SandboxBlock(spec.sandbox, language)}}
}

func ch33WatchStep(spec ch33LessonSpec, language ch33Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".watch-step", Title: language.text(ch33Pair{ar: "خد بالك", en: "Watch out"}), Blocks: []LessonBlock{ch33Markdown(spec.prefix+".watch", language.text(spec.watch))}}
}

func ch33CheckStep(spec ch33LessonSpec, language ch33Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".check-step", Title: language.text(ch33Pair{ar: "فكّر واختبر", en: "Think & check"}), Blocks: []LessonBlock{ch33Markdown(spec.prefix+".plan", language.text(spec.plan)), ch33Question(spec.prefix+".check", spec.check, language)}}
}

func ch33TurnStep(spec ch33LessonSpec, language ch33Language) LessonStep {
	return LessonStep{ID: spec.prefix + ".turn-step", Title: language.text(ch33Pair{ar: "دورك", en: "Your turn"}), Blocks: []LessonBlock{ch33Markdown(spec.prefix+".turn", language.text(spec.turn))}}
}

var ch33Lesson331Arabic, ch33Lesson331English, ch33Lesson331HintsArabic, ch33Lesson331HintsEnglish = ch33BuildLesson(ch33LessonSpec{
	prefix:  "c331",
	mission: ch33Pair{ar: "الـclient في المهمة بيلف transport عشان يعالج الفشل بأمان.\nراجع إمتى إعادة الـrequest ممكن تكرر أثر على الرصيف.", en: "The client in this mission wraps a transport to handle failures safely.\nThink about when replaying a request could repeat an effect at the dock."},
	idea: ch33WalkSpec{text: ch33Pair{ar: "الـRequest جسمه stream بيتقري مرة. ~http.NewRequest~ يقدر يحتفظ بطريقة لإعادة جسم بعض أنواع الـReader.", en: "A request body is a stream that is read once. ~http.NewRequest~ can keep a way to recreate bodies from certain Reader types."}, code: `func main() {
	req, err := http.NewRequest("POST", "https://dock.example/load", strings.NewReader("sample"))
	fmt.Println(err == nil, req.GetBody != nil)
}`, exampleID: "e331-replayable-body", imports: []string{"fmt", "net/http", "strings"}, frames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "جهّز request بجسم ينفع يتبعت تاني.", en: "Create a request with a recreatable body."}, vars: []LessonVariable{{Name: "err", Value: "nil"}, {Name: "req.GetBody", Value: "non-nil"}}},
		{line: 3, caption: ch33Pair{ar: "افحص الجسم وطريقة إعادته.", en: "Check the body and replay function."}, output: "true true"},
	}},
	tryPrompt: ch33Pair{ar: "توقّع الـvalue اللي هترجعها GetBody.", en: "Predict the GetBody value."},
	tryCode: `func main() {
	req, _ := http.NewRequest("POST", "https://dock.example/load", strings.NewReader("x"))
	fmt.Println(req.GetBody != nil)
}`, tryExampleID: "e331-getbody-check", tryImports: []string{"fmt", "net/http", "strings"}, tryFrames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "NewRequest جهز طريقة النسخ.", en: "NewRequest supplied a replay function."}, vars: []LessonVariable{{Name: "req.GetBody", Value: "non-nil"}}},
		{line: 3, caption: ch33Pair{ar: "الـfunction مش nil.", en: "The function is not nil."}, output: "true"},
	},
	tryOptions: []ch33Pair{{ar: "true", en: "true"}, {ar: "false", en: "false"}, {ar: "nil body", en: "nil body"}}, tryAnswer: 0,
	tryFeedback: []ch33Pair{{ar: "strings.Reader بيدعم إعادة الجسم.", en: "strings.Reader supports body replay."}, {ar: "NewRequest ملأ GetBody هنا.", en: "NewRequest populated GetBody here."}, {ar: "الـrequest له جسم في المثال.", en: "The example request has a body."}},
	watch:       ch33Pair{ar: "- أثر POST ممكن يتكرر لو اتبعت تاني.\n- اقفل الـresponse اللي هتسيبه، وسيب الـresponse النهائي للمستدعي.\nالـtransport بينفذ محاولة؛ سياسة الإعادة مسؤولة عن قرارها.", en: "- Replaying a POST can repeat its effect.\n- Close a response you discard; leave the final response to its caller.\nThe transport performs an attempt; retry policy owns the decision."},
	plan:        ch33Pair{ar: "حلّل قابلية الإعادة قبل محاولات الشبكة. راجع ملكية الجسم والـresponse والـcontext بتاع الانتظار.\n\nخلي الإلغاء يوقف انتظار المحاولة التالية. اختبر قرار الإعادة قبل تعديل الجسم.", en: "Decide whether replay is safe before network attempts. Track body, response, and wait ownership.\n\nLet cancellation stop the next wait. Check replayability before changing the body."},
	check:       ch33QuestionSpec{text: ch33Pair{ar: "إيه اللي يخلي إعادة POST آمنة من التكرار؟", en: "What can make replaying a POST safe from duplication?"}, options: []ch33Pair{{ar: "idempotency key يفهمه الـserver", en: "An idempotency key the server honors"}, {ar: "زيادة عدد المحاولات", en: "More attempts"}, {ar: "إغلاق الاتصال أسرع", en: "Closing the connection sooner"}}, answer: 0, feedback: []ch33Pair{{ar: "الـserver يقدر يميّز نفس العملية بالـkey.", en: "The server can recognize the same operation by its key."}, {ar: "المحاولات الكتير ممكن تكرر الأثر.", en: "More attempts can repeat the effect."}, {ar: "سرعة الإغلاق ما تمنعش تكرار الـrequest.", en: "Closing faster does not prevent replay."}}},
	turn:        ch33Pair{ar: "ابدأ بتصنيف الـrequest قبل محاولة الإعادة.\nافحص إغلاق الـresponse والإلغاء مع حالات الـerror.", en: "Classify the request before considering replay.\nCheck response closure and cancellation across failures."},
	hints:       [3]ch33Pair{{ar: "الـrequest القابل للإعادة لازم يحافظ على أثره والـbody بتاعه.", en: "A replayable request must preserve its effect and body."}, {ar: "افصل قرار إعادة الـrequest عن سياسة الانتظار.", en: "Separate the replay decision from the wait policy."}, {ar: "راجع ~http.Request.GetBody~ وملكية ~Response.Body~.", en: "Review ~http.Request.GetBody~ and ownership of ~Response.Body~."}},
	sandbox: ch33SandboxSpec{
		id: "ch33-retry-client", kind: "ch33-http-flow",
		text:    ch33Pair{ar: "تابع جسم الـrequest قبل المحاولة وبعدها.", en: "Step through a request body across attempts."},
		title:   ch33Pair{ar: "جسم الـrequest", en: "Request body"},
		advance: ch33Pair{ar: "الخطوة الجاية", en: "Next event"}, reset: ch33Pair{ar: "من الأول", en: "Reset"}, repeat: ch33Pair{ar: "ترتيب تاني", en: "Another order"},
		runs: []ch33SandboxRun{{label: ch33Pair{ar: "request توضيحي", en: "Illustrative request"}, frames: []ch33SandboxFrame{
			{caption: ch33Pair{ar: "جسم الـrequest لسه مقروء.", en: "The request body is ready."}, leftLabel: ch33Pair{ar: "الجسم", en: "Body"}, leftValue: ch33Pair{ar: "stream جاهز", en: "stream: ready"}, rightLabel: ch33Pair{ar: "GetBody", en: "GetBody"}, rightValue: ch33Pair{ar: "موجود", en: "available"}},
			{caption: ch33Pair{ar: "المحاولة استهلكت الـstream.", en: "The attempt consumed its stream."}, leftLabel: ch33Pair{ar: "المحاولة", en: "Attempt"}, leftValue: ch33Pair{ar: "اترسلت", en: "sent"}, rightLabel: ch33Pair{ar: "GetBody", en: "GetBody"}, rightValue: ch33Pair{ar: "موجود", en: "available"}},
			{caption: ch33Pair{ar: "إعادة المحاولة تحتاج جسم جديد.", en: "A retry needs a fresh body."}, leftLabel: ch33Pair{ar: "المحاولة", en: "Attempt"}, leftValue: ch33Pair{ar: "إعادة", en: "replay"}, rightLabel: ch33Pair{ar: "نسخة الجسم", en: "Body copy"}, rightValue: ch33Pair{ar: "جديدة", en: "fresh"}},
		}}},
	},
})

var ch33Lesson332Arabic, ch33Lesson332English, ch33Lesson332HintsArabic, ch33Lesson332HintsEnglish = ch33BuildLesson(ch33LessonSpec{
	prefix:  "c332",
	mission: ch33Pair{ar: "الوكيل بيستقبل request عام ويختار خدمة خلف من جدول موثوق.\nالتزم بحدود البوابة والهيدرات قبل إرسال الـrequest.", en: "The proxy receives a public request and chooses a backend from trusted routes.\nRespect the gateway and header boundaries before forwarding."},
	idea: ch33WalkSpec{text: ch33Pair{ar: "أي ~http.Handler~ بياخد writer وrequest، ويقرر يكتب response أو يمرر الـrequest.\nالمسار بيتقاس من الـrequest الداخل.", en: "An ~http.Handler~ receives a writer and request, then writes a response or forwards work.\nRouting begins with the incoming request."}, code: `func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	})
	req := httptest.NewRequest("GET", "/hydro/pump", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	fmt.Println(res.Body.String())
}`, exampleID: "e332-handler-path", imports: []string{"fmt", "net/http", "net/http/httptest"}, frames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "الـhandler يستقبل الـrequest.", en: "The handler receives the request."}},
		{line: 5, caption: ch33Pair{ar: "المسار يبدأ بـ hydro.", en: "The path starts with hydro."}, vars: []LessonVariable{{Name: "req.URL.Path", Value: "/hydro/pump"}}},
		{line: 7, caption: ch33Pair{ar: "نفّذ handler محلي.", en: "Run the handler locally."}},
		{line: 8, caption: ch33Pair{ar: "اقرا الـresponse المسجل.", en: "Read the recorded response."}, output: "/hydro/pump"},
	}},
	tryPrompt: ch33Pair{ar: "توقّع المسار اللي يسجله الـresponse.", en: "Predict the path recorded in the response."},
	tryCode: `func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	})
	req := httptest.NewRequest("GET", "/air/valve", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	fmt.Println(res.Body.String())
}`, tryExampleID: "e332-recorded-handler", tryImports: []string{"fmt", "net/http", "net/http/httptest"}, tryFrames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 5, caption: ch33Pair{ar: "الـrequest جاي بمسار air.", en: "The request carries an air path."}, vars: []LessonVariable{{Name: "req.URL.Path", Value: "/air/valve"}}},
		{line: 7, caption: ch33Pair{ar: "writer بيسجل الـresponse.", en: "The recorder stores the response."}},
		{line: 8, caption: ch33Pair{ar: "اطبع الـbody المسجل.", en: "Print the recorded body."}, output: "/air/valve"},
	},
	tryOptions: []ch33Pair{{ar: "/air/valve", en: "/air/valve"}, {ar: "/hydro/pump", en: "/hydro/pump"}, {ar: "مسار فاضي", en: "An empty path"}}, tryAnswer: 0,
	tryFeedback: []ch33Pair{{ar: "الـrequest التجريبي حدد المسار ده.", en: "The test request sets this path."}, {ar: "ده مسار المثال السابق.", en: "That was the previous example's path."}, {ar: "الـrecorder جمع النص المكتوب.", en: "The recorder captured the written text."}},
	watch:       ch33Pair{ar: "- ما تاخدش عنوان الـbackend من query الـclient.\n- امسح الهيدرات اللي الـclient ممكن يزوّرها.\nالمضيف الداخلي يفضل من الإعداد الموثوق.", en: "- Never take a backend address from the client's query.\n- Remove headers that the client could forge.\nKeep backend hosts in trusted configuration."},
	plan:        ch33Pair{ar: "افصل التحقق من route عن إعادة كتابة الـrequest.\nبعدها راجع المهلة، response الوحدة، ومسار الـerror.\nالتوجيه وحده مش بديل عن صلاحية الـrequest.", en: "Separate route validation from request rewriting.\nThen review the timeout, backend response, and error path.\nRouting alone does not authorize a request."},
	check:       ch33QuestionSpec{text: ch33Pair{ar: "مين يحدد عناوين الوحدات الخلفية؟", en: "Who should determine backend unit addresses?"}, options: []ch33Pair{{ar: "جدول مسارات موثوق", en: "A trusted route table"}, {ar: "value يرسلها الـclient", en: "A value supplied by the client"}, {ar: "هيدر X-Forwarded-For", en: "The X-Forwarded-For header"}}, answer: 0, feedback: []ch33Pair{{ar: "الجدول يحدّد الوجهات المسموحة.", en: "The table defines allowed destinations."}, {ar: "عنوان الـclient ممكن يفتح وصول داخلي.", en: "A client URL could expose internal services."}, {ar: "الهيدر بيوصف المصدر، مش الوجهة.", en: "That header describes the source, not the destination."}}},
	turn:        ch33Pair{ar: "ابدأ بالوجهات الموثوقة وسلوك handler.\nجرّب مسار معروف ومسار مرفوض من غير وحدة حقيقية.", en: "Start with trusted targets and handler behavior.\nTry a known route and a rejected route without a real unit."},
	hints:       [3]ch33Pair{{ar: "Reverse proxy يستقبل request ويرجّع response من backend عبر بوابة واحدة.", en: "A reverse proxy receives a request and returns a backend response through one gateway."}, {ar: "افصل جدول الوجهات عن الـrequest الوارد وراجع قواعد الهيدرات.", en: "Separate trusted targets from incoming requests and review header rules."}, {ar: "راجع ~httputil.ProxyRequest.SetURL~ و~ReverseProxy.Rewrite~.", en: "Review ~httputil.ProxyRequest.SetURL~ and ~ReverseProxy.Rewrite~."}},
	sandbox: ch33SandboxSpec{
		id: "ch33-reverse-proxy", kind: "ch33-http-flow",
		text:    ch33Pair{ar: "شوف المسار اللي بيوصل للـhandler.", en: "Watch the path received by a handler."},
		title:   ch33Pair{ar: "المسار الوارد", en: "Incoming path"},
		advance: ch33Pair{ar: "request تاني", en: "Next request"}, reset: ch33Pair{ar: "من الأول", en: "Reset"}, repeat: ch33Pair{ar: "ترتيب تاني", en: "Another order"},
		runs: []ch33SandboxRun{{label: ch33Pair{ar: "requests محلية للتوضيح", en: "Illustrative local requests"}, frames: []ch33SandboxFrame{
			{caption: ch33Pair{ar: "الـhandler استقبل المسار.", en: "The handler received this path."}, leftLabel: ch33Pair{ar: "الـrequest", en: "Request"}, leftValue: ch33Pair{ar: "/hydro/pump", en: "/hydro/pump"}, rightLabel: ch33Pair{ar: "المسجل", en: "Recorded"}, rightValue: ch33Pair{ar: "/hydro/pump", en: "/hydro/pump"}},
			{caption: ch33Pair{ar: "request تاني له مساره.", en: "Another request carries its path."}, leftLabel: ch33Pair{ar: "الـrequest", en: "Request"}, leftValue: ch33Pair{ar: "/air/valve", en: "/air/valve"}, rightLabel: ch33Pair{ar: "المسجل", en: "Recorded"}, rightValue: ch33Pair{ar: "/air/valve", en: "/air/valve"}},
		}}},
	},
})

var ch33Lesson333Arabic, ch33Lesson333English, ch33Lesson333HintsArabic, ch33Lesson333HintsEnglish = ch33BuildLesson(ch33LessonSpec{
	prefix:  "c333",
	mission: ch33Pair{ar: "المهمة بتوزّع الـrequests على وحدات متشابهة وتراقب صحتها.\nالحالة المشتركة محتاجة تفضل سليمة تحت requests متزامنة.", en: "The mission distributes requests across similar units and tracks their health.\nShared state must stay correct under concurrent requests."},
	idea: ch33WalkSpec{text: ch33Pair{ar: "Mutex بيحمي قراءة وتحديث الحالة كعملية واحدة.\nالمثال يبيّن العداد من غير goroutines أو افتراض ترتيب متزامن.", en: "A mutex protects reading and updating shared state as one operation.\nThis example shows a counter without goroutines or assumptions about concurrent order."}, code: `type Counter struct { mu sync.Mutex; next int }
func (c *Counter) Take() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := c.next
	c.next++
	return n
}
func main() {
	c := &Counter{}
	fmt.Println(c.Take(), c.Take())
}`, exampleID: "e333-protected-counter", imports: []string{"fmt", "sync"}, frames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 1, caption: ch33Pair{ar: "Counter تحتفظ بالحالة.", en: "Counter holds shared state."}},
		{line: 10, caption: ch33Pair{ar: "العداد يبدأ بصفر.", en: "The counter starts at zero."}, vars: []LessonVariable{{Name: "c.next", Value: "0"}}},
		{line: 11, caption: ch33Pair{ar: "أول نداء يرجع صفر.", en: "The first call returns zero."}, vars: []LessonVariable{{Name: "c.next", Value: "1"}}},
		{line: 11, caption: ch33Pair{ar: "النداء التاني يرجع واحد.", en: "The second call returns one."}, vars: []LessonVariable{{Name: "c.next", Value: "2"}}},
		{line: 11, caption: ch33Pair{ar: "اطبع الـvalues.", en: "Print both values."}, vars: []LessonVariable{{Name: "c.next", Value: "2"}}, output: "0 1"},
	}},
	tryPrompt: ch33Pair{ar: "توقّع أول values للعداد.", en: "Predict the first two counter values."},
	tryCode: `func main() {
	mu := sync.Mutex{}
	next := 4
	mu.Lock()
	next++
	mu.Unlock()
	fmt.Println(next)
}`, tryExampleID: "e333-locked-update", tryImports: []string{"fmt", "sync"}, tryFrames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "الحالة تبدأ عند 4.", en: "The state starts at 4."}, vars: []LessonVariable{{Name: "next", Value: "4"}}},
		{line: 3, caption: ch33Pair{ar: "الـlock اتاخد.", en: "The lock is held."}, vars: []LessonVariable{{Name: "next", Value: "4"}}},
		{line: 4, caption: ch33Pair{ar: "زوّد الـvalue مرة.", en: "Increment the value once."}, vars: []LessonVariable{{Name: "next", Value: "5"}}},
		{line: 6, caption: ch33Pair{ar: "اطبع الـvalue الجديدة.", en: "Print the updated value."}, vars: []LessonVariable{{Name: "next", Value: "5"}}, output: "5"},
	},
	tryOptions: []ch33Pair{{ar: "5", en: "5"}, {ar: "4", en: "4"}, {ar: "0", en: "0"}}, tryAnswer: 0,
	tryFeedback: []ch33Pair{{ar: "التحديث زوّد 4 بواحد.", en: "The update added one to 4."}, {ar: "الطباعة بعد التحديث.", en: "Printing happens after the update."}, {ar: "next بدأت بـvalue 4.", en: "next started at 4."}},
	watch:       ch33Pair{ar: "- ماتمسكش الـlock وانت مستني الشبكة.\n- نتيجة الفحص والـbody محتاجين تنظيف واضح.\nفحص الصحة يقدر يرجّع وحدة سليمة للخدمة.\nحدّث الصحة بعد ما نتيجة الشبكة تبقى معروفة.", en: "- Do not hold a lock while waiting on the network.\n- Health results and response bodies need clear cleanup.\nA health check can bring a recovered unit back.\nUpdate health after the network result is known."},
	plan:        ch33Pair{ar: "راجع تحديث الحالة، اختيار الوحدة، ثم فحص الصحة.\nاستخدم -race كدليل إضافي على الحالة المشتركة.", en: "Review state updates, unit selection, then health checks.\nUse -race as additional evidence about shared state."},
	check:       ch33QuestionSpec{text: ch33Pair{ar: "فين تعمل request health check بطيء؟", en: "Where should a slow health check run?"}, options: []ch33Pair{{ar: "برّه الـlock", en: "Outside the lock"}, {ar: "والـlock ماسك طول الشبكة", en: "While holding the lock across the network call"}, {ar: "في fmt.Println", en: "Inside fmt.Println"}}, answer: 0, feedback: []ch33Pair{{ar: "الشبكة ممكن تتأخر؛ خلي الـrequests التانية تكمل.", en: "Network calls can wait; let other requests proceed."}, {ar: "الـlock الطويل يوقف تحديثات واختيارات تانية.", en: "A long-held lock blocks other updates and picks."}, {ar: "الطباعة مش مكان فحص الصحة.", en: "Printing is not a health check."}}},
	turn:        ch33Pair{ar: "ابدأ بحالة صغيرة وتتبّع lock الحالة.\nبعدها راجع فشل وحدة وتعافيها من غير انتظار دوري.", en: "Start with small state and trace its lock.\nThen inspect unit failure and recovery without a periodic timer."},
	hints:       [3]ch33Pair{{ar: "Mutex بيحمي الحالة المشتركة من التحديثات المتداخلة.", en: "A mutex protects shared state from overlapping updates."}, {ar: "افصل قرار الاختيار عن health check requests على الشبكة.", en: "Separate selection decisions from network health checks."}, {ar: "راجع ~sync.Mutex~ و~http.NewRequestWithContext~ و~RoundTrip~.", en: "Review ~sync.Mutex~, ~http.NewRequestWithContext~, and ~RoundTrip~."}},
	sandbox: ch33SandboxSpec{
		id: "ch33-load-balancer", kind: "ch33-http-flow",
		text:    ch33Pair{ar: "بدّل ترتيب تحديثات goroutines وشوف العداد.", en: "Switch possible goroutine orders and watch the counter."},
		title:   ch33Pair{ar: "ترتيب التحديثات", en: "Update order"},
		advance: ch33Pair{ar: "الحدث الجاي", en: "Next event"}, reset: ch33Pair{ar: "من الأول", en: "Reset"}, repeat: ch33Pair{ar: "ترتيب محتمل تاني", en: "Another possible order"},
		runs: []ch33SandboxRun{
			{label: ch33Pair{ar: "ترتيب محتمل: A ثم B", en: "One possible order: A then B"}, frames: []ch33SandboxFrame{
				{caption: ch33Pair{ar: "العداد محمي بـlock.", en: "The counter is protected by a lock."}, leftLabel: ch33Pair{ar: "الـvalue", en: "Value"}, leftValue: ch33Pair{ar: "0", en: "0"}, rightLabel: ch33Pair{ar: "التالي", en: "Next"}, rightValue: ch33Pair{ar: "A", en: "A"}},
				{caption: ch33Pair{ar: "A حدّثت الحالة.", en: "A updated the shared state."}, leftLabel: ch33Pair{ar: "الـvalue", en: "Value"}, leftValue: ch33Pair{ar: "1", en: "1"}, rightLabel: ch33Pair{ar: "التالي", en: "Next"}, rightValue: ch33Pair{ar: "B", en: "B"}},
				{caption: ch33Pair{ar: "B حدّثت الحالة.", en: "B updated the shared state."}, leftLabel: ch33Pair{ar: "الـvalue", en: "Value"}, leftValue: ch33Pair{ar: "2", en: "2"}, rightLabel: ch33Pair{ar: "الحالة", en: "State"}, rightValue: ch33Pair{ar: "محمية", en: "protected"}},
			}},
			{label: ch33Pair{ar: "ترتيب محتمل: B ثم A", en: "Another possible order: B then A"}, frames: []ch33SandboxFrame{
				{caption: ch33Pair{ar: "العداد محمي بـlock.", en: "The counter is protected by a lock."}, leftLabel: ch33Pair{ar: "الـvalue", en: "Value"}, leftValue: ch33Pair{ar: "0", en: "0"}, rightLabel: ch33Pair{ar: "التالي", en: "Next"}, rightValue: ch33Pair{ar: "B", en: "B"}},
				{caption: ch33Pair{ar: "B حدّثت الحالة.", en: "B updated the shared state."}, leftLabel: ch33Pair{ar: "الـvalue", en: "Value"}, leftValue: ch33Pair{ar: "1", en: "1"}, rightLabel: ch33Pair{ar: "التالي", en: "Next"}, rightValue: ch33Pair{ar: "A", en: "A"}},
				{caption: ch33Pair{ar: "A حدّثت الحالة.", en: "A updated the shared state."}, leftLabel: ch33Pair{ar: "الـvalue", en: "Value"}, leftValue: ch33Pair{ar: "2", en: "2"}, rightLabel: ch33Pair{ar: "الحالة", en: "State"}, rightValue: ch33Pair{ar: "محمية", en: "protected"}},
			}},
		},
	},
})

var ch33Lesson334Arabic, ch33Lesson334English, ch33Lesson334HintsArabic, ch33Lesson334HintsEnglish = ch33BuildLesson(ch33LessonSpec{
	prefix:  "c334",
	mission: ch33Pair{ar: "البوابة تجمع طبقات مشتركة قبل الـrequest اللي رايح للخدمة.\nترتيب الطبقات يحدد إيه اللي يحصل للـrequest المرفوض.", en: "The gateway composes shared layers before forwarding to a service.\nLayer order determines what happens to a rejected request."},
	idea: ch33WalkSpec{text: ch33Pair{ar: "Middleware بتلف handler. تقدر تجهز الـresponse قبل ما تنادي الطبقة اللي بعدها.\n\nالـrequest context بينقل value محتاجاها الطبقة اللي جوه في السلسلة. والطبقة ممكن ترفض الـrequest وتوقف السلسلة.", en: "Middleware wraps a handler. It can set a response before calling the next layer.\n\nRequest context carries values inward. A layer can reject and stop the chain."}, code: `func main() {
	stamp := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Lesson", "ready")
			next.ServeHTTP(w, r)
		})
	}
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "ok") })
	res := httptest.NewRecorder()
	stamp(base).ServeHTTP(res, httptest.NewRequest("GET", "/", nil))
	fmt.Println(res.Header().Get("X-Lesson"), res.Body.String())
}`, exampleID: "e334-middleware-wrap", imports: []string{"fmt", "net/http", "net/http/httptest"}, frames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "stamp بترجع handler جديدة.", en: "stamp returns a wrapped handler."}},
		{line: 4, caption: ch33Pair{ar: "الطبقة بتحط هيدر.", en: "The layer adds a header."}},
		{line: 5, caption: ch33Pair{ar: "بعدها تنادي handler اللي بعدها.", en: "Then it calls the next handler."}},
		{line: 10, caption: ch33Pair{ar: "شغّل السلسلة محليًا.", en: "Run the chain locally."}},
		{line: 11, caption: ch33Pair{ar: "اقرا الـheader والـbody.", en: "Read the header and body."}, output: "ready ok"},
	}},
	tryPrompt: ch33Pair{ar: "توقّع الهيدر والجسم.", en: "Predict the header and body."},
	tryCode: `func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Flag", "on")
		fmt.Fprint(w, "done")
	})
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest("GET", "/", nil))
	fmt.Println(res.Header().Get("X-Flag"), res.Body.String())
}`, tryExampleID: "e334-handler-response", tryImports: []string{"fmt", "net/http", "net/http/httptest"}, tryFrames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 3, caption: ch33Pair{ar: "الـhandler جهز الهيدر.", en: "The handler sets a header."}},
		{line: 4, caption: ch33Pair{ar: "اكتب الجسم.", en: "Write the body."}},
		{line: 7, caption: ch33Pair{ar: "سجّل الـrequest والـresponse.", en: "Record the request and response."}},
		{line: 8, caption: ch33Pair{ar: "اطبع الحقلين.", en: "Print both values."}, output: "on done"},
	},
	tryOptions: []ch33Pair{{ar: "on done", en: "on done"}, {ar: "done on", en: "done on"}, {ar: "فاضي فاضي", en: "empty empty"}}, tryAnswer: 0,
	tryFeedback: []ch33Pair{{ar: "الهيدر والجسم اتحضروا قبل الطباعة.", en: "The header and body were set before printing."}, {ar: "الـfmt.Println بتعرض الهيدر الأول.", en: "fmt.Println prints the header first."}, {ar: "الـhandler كتب الاتنين.", en: "The handler wrote both values."}},
	watch:       ch33Pair{ar: "- ترتيب الطبقات سياسة للبوابة.\n- امسح الأسرار والهيدرات الموثوقة اللي جايه من الـclient.", en: "- Layer order is gateway policy.\n- Strip secrets and trusted headers supplied by the client."},
	plan:        ch33Pair{ar: "تتبّع الـrequest من أول طبقة لحد الـproxy.\nاختبر الرفض قبل الشبكة، وبعدها راجع الهيدرات المرسلة.", en: "Trace a request from the first layer to the proxy.\nCheck rejection before networking, then inspect forwarded headers."},
	check:       ch33QuestionSpec{text: ch33Pair{ar: "إمتى handler ملفوفة بتشغّل اللي بعدها؟", en: "When does a wrapped handler run the next handler?"}, options: []ch33Pair{{ar: "لما تنادي ~next.ServeHTTP~", en: "When it calls ~next.ServeHTTP~"}, {ar: "وقت تعريف الـfunction بس", en: "Only when the function is defined"}, {ar: "بعد ما request يختفي", en: "After the request disappears"}}, answer: 0, feedback: []ch33Pair{{ar: "النداء بيمرر الـrequest للطبقة الداخلية.", en: "The call passes the request to the inner layer."}, {ar: "التعريف لوحده مش بيشغّلها.", en: "Defining a function does not run it."}, {ar: "الـhandler شغالة أثناء الـrequest.", en: "The handler runs during the request."}}},
	turn:        ch33Pair{ar: "ابنِ طبقة واحدة وجربها محليًا.\nبعدها راجع ترتيب الرفض والهيدرات قبل توصيلها للخدمات.", en: "Build one layer and run it locally.\nThen review rejection order and headers before connecting services."},
	hints:       [3]ch33Pair{{ar: "Middleware هي function بتلف ~http.Handler~ تاني.", en: "Middleware is a function that wraps another ~http.Handler~."}, {ar: "ارسم الـvalues اللي تعرفها كل طبقة وإيه اللي بعدها محتاجاه.", en: "Map what each layer learns and what the next layer needs."}, {ar: "راجع ~http.HandlerFunc~ و~Request.WithContext~ و~Chain~.", en: "Review ~http.HandlerFunc~, ~Request.WithContext~, and ~Chain~."}},
	sandbox: ch33SandboxSpec{
		id: "ch33-gateway", kind: "ch33-http-flow",
		text:    ch33Pair{ar: "تابع الـrequest وهو بيعدي من طبقات البوابة.", en: "Follow a request through gateway layers."},
		title:   ch33Pair{ar: "طبقات البوابة", en: "Gateway layers"},
		advance: ch33Pair{ar: "الطبقة الجاية", en: "Next layer"}, reset: ch33Pair{ar: "من الأول", en: "Reset"}, repeat: ch33Pair{ar: "ترتيب تاني", en: "Another order"},
		runs: []ch33SandboxRun{{label: ch33Pair{ar: "request توضيحي", en: "Illustrative request"}, frames: []ch33SandboxFrame{
			{caption: ch33Pair{ar: "الـrequest وصل للبوابة.", en: "The request reached the gateway."}, leftLabel: ch33Pair{ar: "المسار", en: "Path"}, leftValue: ch33Pair{ar: "/hydro/pump", en: "/hydro/pump"}, rightLabel: ch33Pair{ar: "الطبقة", en: "Layer"}, rightValue: ch33Pair{ar: "identity", en: "identity"}},
			{caption: ch33Pair{ar: "طبقة أضافت context.", en: "A layer added request context."}, leftLabel: ch33Pair{ar: "الـcontext", en: "Context"}, leftValue: ch33Pair{ar: "crew key", en: "crew key"}, rightLabel: ch33Pair{ar: "الطبقة", en: "Layer"}, rightValue: ch33Pair{ar: "rate check", en: "rate check"}},
			{caption: ch33Pair{ar: "الـrequest وصل للـhandler الداخلي.", en: "The request reached the inner handler."}, leftLabel: ch33Pair{ar: "الـcontext", en: "Context"}, leftValue: ch33Pair{ar: "ready", en: "ready"}, rightLabel: ch33Pair{ar: "الطبقة", en: "Layer"}, rightValue: ch33Pair{ar: "proxy", en: "proxy"}},
		}}},
	},
})

var ch33Lesson33BossArabic, ch33Lesson33BossEnglish, ch33Lesson33BossHintsArabic, ch33Lesson33BossHintsEnglish = ch33BuildLesson(ch33LessonSpec{
	prefix:  "c33boss",
	mission: ch33Pair{ar: "ده تجميع لأدوات اتعلمتها قبل كده في وصلة واحدة.\nابدأ من العقد المكتوب، وحدد مين يملك الاتصال والـresponse في كل خروج.", en: "This checkpoint combines tools you have already learned in one relay.\nStart from the written contract and track who owns the connection and response on every exit."},
	idea: ch33WalkSpec{text: ch33Pair{ar: "~io.ReadFull~ بتملا الـbuffer المطلوب أو بترجع error.\nالـlength اللي في header لازم يتقري قبل المحتوى.", en: "~io.ReadFull~ fills the requested buffer or returns an error.\nRead a header length before its content."}, code: `func main() {
	frame := []byte{0, 3, 'i', 'c', 'e'}
	length := binary.BigEndian.Uint16(frame[:2])
	fmt.Println(length, string(frame[2:]))
}`, exampleID: "e33boss-frame-length", imports: []string{"encoding/binary", "fmt"}, frames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "الإطار بيبدأ بطول 2 bytes.", en: "The frame begins with two length bytes."}, vars: []LessonVariable{{Name: "frame", Value: "[0 3 105 99 101]"}}},
		{line: 3, caption: ch33Pair{ar: "BigEndian تقرا الرقم.", en: "BigEndian reads the number."}, vars: []LessonVariable{{Name: "length", Value: "3"}}},
		{line: 4, caption: ch33Pair{ar: "اعرض الطول والاسم.", en: "Print length and station."}, vars: []LessonVariable{{Name: "length", Value: "3"}}, output: "3 ice"},
	}},
	tryPrompt: ch33Pair{ar: "توقّع الرقم اللي تقراه BigEndian.", en: "Predict the number read by BigEndian."},
	tryCode: `func main() {
	var header [2]byte
	header[0], header[1] = 0, 5
	n := binary.BigEndian.Uint16(header[:])
	fmt.Println(n)
}`, tryExampleID: "e33boss-big-endian", tryImports: []string{"encoding/binary", "fmt"}, tryFrames: []ch33FrameSpec{
		{line: 0, caption: ch33Pair{ar: "قبل التنفيذ.", en: "Before execution."}},
		{line: 2, caption: ch33Pair{ar: "الـheader مساحته 2 bytes.", en: "header has two bytes."}, vars: []LessonVariable{{Name: "header", Value: "[0 0]"}}},
		{line: 3, caption: ch33Pair{ar: "اكتب الطول للتوضيح.", en: "Set an illustrative length."}, vars: []LessonVariable{{Name: "header", Value: "[0 5]"}}},
		{line: 4, caption: ch33Pair{ar: "اقراه كـuint16.", en: "Read it as uint16."}, vars: []LessonVariable{{Name: "n", Value: "5"}}},
		{line: 5, caption: ch33Pair{ar: "اطبع الرقم.", en: "Print the number."}, vars: []LessonVariable{{Name: "n", Value: "5"}}, output: "5"},
	},
	tryOptions: []ch33Pair{{ar: "5", en: "5"}, {ar: "1280", en: "1280"}, {ar: "0", en: "0"}}, tryAnswer: 0,
	tryFeedback: []ch33Pair{{ar: "الـbyte الأعلى جه قبل الأقل.", en: "The high byte comes before the low byte."}, {ar: "ده ترتيب little-endian.", en: "That is little-endian order."}, {ar: "الـheader اتملأ بـvalue 5.", en: "The header contains the value 5."}},
	watch:       ch33Pair{ar: "- ما تحجزش محتوى قبل التحقق من طوله.\n- كل response HTTP له body لازم يتقفل في مساره.", en: "- Do not allocate content before validating its length.\n- Every HTTP response body needs closure on its path."},
	plan:        ch33Pair{ar: "قسّم العمل لقراءة، تحقق، request HTTP، ثم response.\nبعد كل خطوة اسأل مين يملك المورد لو حصل error.", en: "Split the work into reading, validation, an HTTP request, then a reply.\nAt each step ask who owns the resource if an error occurs."},
	check:       ch33QuestionSpec{text: ch33Pair{ar: "إمتى ~io.ReadFull~ بترجع من غير ما تكمل؟", en: "When does ~io.ReadFull~ return before filling the buffer?"}, options: []ch33Pair{{ar: "لما يحصل error أو نهاية مبكرة", en: "When an error or early EOF occurs"}, {ar: "بعد أول byte دايمًا", en: "Always after the first byte"}, {ar: "لما المحتوى يبقى UTF-8", en: "When the content is UTF-8"}}, answer: 0, feedback: []ch33Pair{{ar: "المصدر ممكن يقف قبل المطلوب.", en: "The source can stop before the requested bytes arrive."}, {ar: "هي بتحاول تملا الـbuffer كله.", en: "It tries to fill the entire buffer."}, {ar: "UTF-8 مش شرط لقراية الـbytes.", en: "UTF-8 is unrelated to reading bytes."}}},
	turn:        ch33Pair{ar: "اعتبر العقد خريطة ملكية للمدخلات والمخرجات.\nاستخدم errors صغيرة لتتبع كل حد من غير ما تبني النظام كله مرة واحدة.", en: "Treat the contract as an ownership map for inputs and outputs.\nUse small error cases to trace each boundary before assembling the relay."},
	hints:       [3]ch33Pair{{ar: "اقرا أطوال إطارات الشبكة المحددة قبل محتواها.", en: "Read network frames as bounded lengths before their contents."}, {ar: "ارسم سلسلة الملكية من net.Conn لحد HTTP body والـresponse.", en: "Trace ownership from net.Conn through the HTTP body and reply."}, {ar: "راجع ~binary.BigEndian~ و~io.ReadFull~ و~NewRequestWithContext~.", en: "Review ~binary.BigEndian~, ~io.ReadFull~, and ~NewRequestWithContext~."}},
	sandbox: ch33SandboxSpec{
		id: "ch33-survey-frame", kind: "ch33-http-flow",
		text:    ch33Pair{ar: "اقرا الطول قبل بيانات الإطار.", en: "Read a frame length before its payload."},
		title:   ch33Pair{ar: "إطار توضيحي", en: "Illustrative frame"},
		advance: ch33Pair{ar: "الجزء الجاي", en: "Next part"}, reset: ch33Pair{ar: "من الأول", en: "Reset"}, repeat: ch33Pair{ar: "ترتيب تاني", en: "Another order"},
		runs: []ch33SandboxRun{{label: ch33Pair{ar: "bytes للتوضيح", en: "Illustrative bytes"}, frames: []ch33SandboxFrame{
			{caption: ch33Pair{ar: "الـheader فيه 2 bytes.", en: "The header has two bytes."}, leftLabel: ch33Pair{ar: "الـheader", en: "Header"}, leftValue: ch33Pair{ar: "00 03", en: "00 03"}, rightLabel: ch33Pair{ar: "المحتوى", en: "Payload"}, rightValue: ch33Pair{ar: "لسه", en: "waiting"}},
			{caption: ch33Pair{ar: "الطول اتقرا بترتيب big-endian.", en: "The big-endian length was read."}, leftLabel: ch33Pair{ar: "الطول", en: "Length"}, leftValue: ch33Pair{ar: "3", en: "3"}, rightLabel: ch33Pair{ar: "المحتوى", en: "Payload"}, rightValue: ch33Pair{ar: "بعده", en: "next"}},
			{caption: ch33Pair{ar: "المحتوى اتقرا بعد الطول.", en: "The payload follows its length."}, leftLabel: ch33Pair{ar: "الطول", en: "Length"}, leftValue: ch33Pair{ar: "3", en: "3"}, rightLabel: ch33Pair{ar: "المحتوى", en: "Payload"}, rightValue: ch33Pair{ar: "ice", en: "ice"}},
		}}},
	},
})
