package levels

type ch28LessonCopy struct {
	Mission       string
	Idea          string
	IdeaCode      string
	IdeaImports   []string
	IdeaID        string
	IdeaStates    []LessonFrame
	IdeaCaptions  []string
	TryPrompt     string
	TryCode       string
	TryImports    []string
	TryID         string
	TryAnswer     int
	TryOptions    []string
	TryFeedback   []string
	TryStates     []LessonFrame
	TryCaptions   []string
	Reveal        string
	Watch         string
	Plan          string
	CheckQuestion string
	CheckOptions  []string
	CheckAnswer   int
	CheckFeedback []string
	Turn          string
	Hints         []string
}

type ch28SandboxSpec struct {
	id     string
	kind   string
	text   [2]string
	config map[string]any
}

func ch28Pair(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

var ch28SandboxSpecs = map[string]ch28SandboxSpec{
	"c281": {
		id: "ch28-gin-playground", kind: "ch28-gin-flow",
		text: [2]string{"تابع الطلب بين الطبقات.", "Follow a request across layers."},
		config: map[string]any{
			"labels": map[string]any{
				"step":     ch28Pair("الطبقة", "Layer"),
				"value":    ch28Pair("القيمة", "Value"),
				"previous": ch28Pair("السابق", "Previous"),
				"next":     ch28Pair("التالي", "Next"),
				"reset":    ch28Pair("من الأول", "Reset"),
			},
			"frames": []any{
				map[string]any{"step": "request", "value": "  Vega  ", "caption": ch28Pair("الطلب دخل بقيمة خام.", "The request carries raw input.")},
				map[string]any{"step": "handler", "value": "Vega", "caption": ch28Pair("الـhandler نظّف المدخل.", "The handler normalized the input.")},
				map[string]any{"step": "registry", "value": "Vega", "caption": ch28Pair("الـregistry هو صاحب الحالة.", "The registry owns the state.")},
				map[string]any{"step": "response", "value": "Vega", "caption": ch28Pair("الرد رجّع القيمة النظيفة.", "The response returns the clean value.")},
			},
		},
	},
	"c282": {
		id: "ch28-fiber-playground", kind: "ch28-fiber-copy",
		text: [2]string{"غيّر الـbuffer وراقب النسخة.", "Change the ~buffer~ and watch the copy."},
		config: map[string]any{
			"labels": map[string]any{
				"buffer":  ch28Pair("buffer", "Buffer"),
				"saved":   ch28Pair("النص المحفوظ", "Saved string"),
				"reuse":   ch28Pair("أعد استخدام buffer", "Reuse buffer"),
				"caption": ch28Pair("النسخة المحفوظة مستقلة.", "The saved string stays independent."),
			},
			"values": []string{"north", "sorth"},
		},
	},
	"c283": {
		id: "ch28-openapi-playground", kind: "ch28-json-shape",
		text: [2]string{"بدّل نوع القيمة وراقب الفحص.", "Switch the value type and inspect the assertion."},
		config: map[string]any{
			"labels": map[string]any{
				"field":    ch28Pair("قيمة depth", "depth value"),
				"jsonType": ch28Pair("نوع JSON", "JSON type"),
				"number":   ch28Pair("فحص الرقم", "Number assertion"),
				"next":     ch28Pair("قيمة تانية", "Next sample"),
			},
			"samples": []any{
				map[string]any{"value": "3", "type": "number", "number": "true", "caption": ch28Pair("الرقم يطابق النوع.", "A number matches the type.")},
				map[string]any{"value": "\"3\"", "type": "string", "number": "false", "caption": ch28Pair("النص مش رقمًا.", "A string is not a number.")},
				map[string]any{"value": "—", "type": "missing", "number": "false", "caption": ch28Pair("المفتاح الغايب ملوش قيمة.", "A missing key has no value.")},
			},
		},
	},
	"c284": {
		id: "ch28-shutdown-playground", kind: "ch28-shutdown-order",
		text: [2]string{"شغّل تاني؛ ترتيب الطلبات ممكن يتغير.", "Run again; request order can vary."},
		config: map[string]any{
			"labels": map[string]any{
				"timeline": ch28Pair("تسلسل توضيحي", "Illustrative timeline"),
				"runAgain": ch28Pair("بدّل الترتيب / شغّل تاني", "Shuffle / run again"),
				"caption":  ch28Pair("طلبات نشطة بترتيب متغير.", "Active requests may finish in either order."),
			},
			"orders": []any{
				[]any{
					ch28Pair("الإغلاق بدأ", "Shutdown starts"),
					ch28Pair("الطلب A خلص", "Request A completes"),
					ch28Pair("الطلب B خلص", "Request B completes"),
					ch28Pair("التنضيف بدأ", "Cleanup starts"),
				},
				[]any{
					ch28Pair("الإغلاق بدأ", "Shutdown starts"),
					ch28Pair("الطلب B خلص", "Request B completes"),
					ch28Pair("الطلب A خلص", "Request A completes"),
					ch28Pair("التنضيف بدأ", "Cleanup starts"),
				},
			},
		},
	},
	"c28b": {
		id: "ch28-boss-playground", kind: "ch28-response-commit",
		text: [2]string{"تابع الرد لحد ما metadata تثبت.", "Follow the response until metadata commits."},
		config: map[string]any{
			"labels": map[string]any{
				"status":   ch28Pair("status", "Status"),
				"header":   ch28Pair("header", "Header"),
				"body":     ch28Pair("body", "Body"),
				"state":    ch28Pair("حالة الرد", "Response state"),
				"previous": ch28Pair("السابق", "Previous"),
				"next":     ch28Pair("التالي", "Next"),
				"reset":    ch28Pair("من الأول", "Reset"),
			},
			"frames": []any{
				map[string]any{"status": "unset", "header": "unset", "body": "—", "state": "open", "caption": ch28Pair("لسه مفيش بيانات اتكتبت.", "Nothing has been written yet.")},
				map[string]any{"status": "sample", "header": "example/type", "body": "—", "state": "open", "caption": ch28Pair("الـmetadata جاهزة والرد مفتوح.", "Metadata is set; the response is open.")},
				map[string]any{"status": "sample", "header": "example/type", "body": "ready", "state": "committed", "caption": ch28Pair("كتابة body ثبّتت الرد.", "Writing the body committed the response.")},
				map[string]any{"status": "sample", "header": "example/type", "body": "ready", "state": "committed", "caption": ch28Pair("تغيير metadata المتأخر ما أثرش.", "A late metadata change had no effect.")},
			},
		},
	},
}

func ch28SandboxBlock(prefix string, english bool) LessonBlock {
	spec := ch28SandboxSpecs[prefix]
	text := spec.text[0]
	if english {
		text = spec.text[1]
	}
	return LessonBlock{
		ID: spec.id, Kind: "sandbox", Text: text,
		Sandbox: &LessonSandbox{Kind: spec.kind, Config: spec.config},
	}
}

func ch28Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch28Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch28Walkthrough(id, text, code string, imports []string, exampleID string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch28Frames(states []LessonFrame, captions []string) []LessonFrame {
	frames := make([]LessonFrame, len(states))
	for i, frame := range states {
		frame.Caption = captions[i]
		frames[i] = frame
	}
	return frames
}

func ch28BuildSteps(prefix string, english bool, lesson ch28LessonCopy) []LessonStep {
	titles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	tryPrompt := "توقّع الناتج:"
	if english {
		titles = []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
		tryPrompt = "Predict the output:"
	}
	if lesson.TryPrompt == "" {
		lesson.TryPrompt = tryPrompt
	}
	questionText := lesson.TryPrompt + "\n\n~~~go\n" + lesson.TryCode + "\n~~~"
	return []LessonStep{
		{ID: prefix + ".mission-step", Title: titles[0], Blocks: []LessonBlock{ch28Markdown(prefix+".mission", lesson.Mission)}},
		{ID: prefix + ".idea-step", Title: titles[1], Blocks: []LessonBlock{ch28Walkthrough(prefix+".idea-trace", lesson.Idea, lesson.IdeaCode, lesson.IdeaImports, lesson.IdeaID, ch28Frames(lesson.IdeaStates, lesson.IdeaCaptions))}},
		{ID: prefix + ".try-step", Title: titles[2], Blocks: []LessonBlock{
			ch28Question(prefix+".try-predict", questionText, lesson.TryOptions, lesson.TryAnswer, lesson.TryFeedback),
			ch28Walkthrough(prefix+".try-reveal", lesson.Reveal, lesson.TryCode, lesson.TryImports, lesson.TryID, ch28Frames(lesson.TryStates, lesson.TryCaptions)),
			ch28SandboxBlock(prefix, english),
		}},
		{ID: prefix + ".watch-step", Title: titles[3], Blocks: []LessonBlock{ch28Markdown(prefix+".watch", lesson.Watch)}},
		{ID: prefix + ".check-step", Title: titles[4], Blocks: []LessonBlock{
			ch28Markdown(prefix+".plan", lesson.Plan),
			ch28Question(prefix+".check", lesson.CheckQuestion, lesson.CheckOptions, lesson.CheckAnswer, lesson.CheckFeedback),
		}},
		{ID: prefix + ".turn-step", Title: titles[5], Blocks: []LessonBlock{ch28Markdown(prefix+".turn", lesson.Turn)}},
	}
}

var ch28GinIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "label", Value: "Vega"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "label", Value: "Vega"}}, Output: "Vega"},
}

var ch28GinTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "label", Value: "Luna"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "label", Value: "Luna"}}, Output: "Luna"},
}

var ch28GinArabic = ch28LessonCopy{
	Mission: `انقل سجل الأرصفة لـ Gin مع الحفاظ على عقد الطلب والرد.

الـstarter فيه ~NewBerthEngine~ بيرد بـ~TODO~. خليك مركز على مسؤولية كل طبقة.

راجع كمان الطلب الفاضي.`,
	Idea: `~gin.New()~ بيبدأ router من غير middleware تلقائي.

الـhandler ينسّق الطلب؛ المخزن يملك الحالة. افصل تنظيف المدخل عن الراوتنج.`,
	IdeaCode: `func main() {
	label := strings.TrimSpace("  Vega  ")
	fmt.Println(label)
}`,
	IdeaImports: []string{"fmt", "strings"}, IdeaID: "e281-normalize", IdeaStates: ch28GinIdeaStates,
	IdeaCaptions: []string{"قبل التنفيذ.", "نضّف الاسم.", "اطبع الاسم النضيف."},
	TryPrompt:    "توقّع الناتج بعد تنظيف المسافات.",
	TryCode: `func main() {
	label := strings.TrimSpace("  Luna  ")
	fmt.Println(label)
}`,
	TryImports: []string{"fmt", "strings"}, TryID: "e281-trim", TryAnswer: 0,
	TryOptions:  []string{"Luna", "  Luna", "luna"},
	TryFeedback: []string{"TrimSpace شالت المسافات من الطرفين.", "المسافات الطرفية مش جزء من الناتج.", "TrimSpace ما بتغيّرش حالة الحروف."},
	TryStates:   []LessonFrame{{Line: 0}, {Line: 2, Variables: []LessonVariable{{Name: "label", Value: "Luna"}}}, {Line: 3, Variables: []LessonVariable{{Name: "label", Value: "Luna"}}, Output: "Luna"}},
	TryCaptions: []string{"قبل التنفيذ.", "القيمة بقت Luna.", "ده النص المطبوع."},
	Reveal:      `تابع القيمة بعد التنظيف، وبعدين شوف النص المطبوع.`,
	Watch: `إعدادات الراوتر جزء من العقد؛ الشرطة الأخيرة ممكن تغيّر نتيجة المسار.

خلّي الحالة مشتركة بين handlers، واحمي تحديثها.`,
	Plan: `راجع رحلة الطلب: route، تحقق، store، ثم response.

ارفض المدخل قبل تعديل الحالة. الرفض المفروض ما يستهلكش ID.`,
	CheckQuestion: `لو POST اترفض قبل الحفظ، إيه المتوقع من العداد؟`,
	CheckOptions:  []string{"يفضل زي ما هو", "يزيد مرة", "يرجع للصفر"}, CheckAnswer: 0,
	CheckFeedback: []string{"الطلب المرفوض ما بيغيّرش الحالة.", "الرقم بيتقدم مع نجاح الحفظ.", "العداد بيحتفظ بتقدمه السابق."},
	Turn: `ابدأ بالـrouter، وبعدها شغّل طلبات من الذاكرة.

قارن status وbody والهيدرز بالعقد.

راجع كمان المسار بشرطته الأخيرة.`,
	Hints: []string{
		`افصل مسؤولية الراوتنج عن ملكية بيانات السجل.`,
		`امشِ مع الطلب من handler لحد الرد، ودوّن حدود كل جزء.`,
		`راجع placeholder جوه ~NewBerthEngine~ وواجهات ~gin.Context~.`,
	},
}

var ch28GinEnglish = ch28LessonCopy{
	Mission: `Move the berth registry to Gin while preserving its request and response contract.

The ~NewBerthEngine~ starter returns ~TODO~. Keep each layer's responsibility clear.

Also inspect an empty request.`,
	Idea: `~gin.New()~ starts a router without automatic middleware.

The handler coordinates a request; the store owns state. Separate input cleanup from routing.`,
	IdeaCode: ch28GinArabic.IdeaCode, IdeaImports: ch28GinArabic.IdeaImports, IdeaID: ch28GinArabic.IdeaID, IdeaStates: ch28GinIdeaStates,
	IdeaCaptions: []string{"Before execution.", "Trim the label.", "Print the clean label."},
	TryPrompt:    "Predict the output after trimming spaces.",
	TryCode:      ch28GinArabic.TryCode, TryImports: ch28GinArabic.TryImports, TryID: ch28GinArabic.TryID, TryAnswer: 0,
	TryOptions:  []string{"Luna", "  Luna", "luna"},
	TryFeedback: []string{"TrimSpace removes spaces at both ends.", "Leading and trailing spaces are not printed.", "TrimSpace does not change letter case."},
	TryStates:   ch28GinTryStates,
	TryCaptions: []string{"Before execution.", "The value is Luna.", "That text is printed."},
	Reveal:      `Follow the cleaned value, then watch the printed text.`,
	Watch: `Router settings are part of the contract; a trailing slash can change route results.

Keep state shared across handlers, and protect its updates.`,
	Plan: `Trace a request through route, validation, store, then response.

Reject invalid input before changing state. A rejection should not consume an ID.`,
	CheckQuestion: `If POST is rejected before saving, what should happen to the counter?`,
	CheckOptions:  []string{"It stays the same", "It increases once", "It resets to zero"}, CheckAnswer: 0,
	CheckFeedback: []string{"A rejected request does not change state.", "The number advances after a successful save.", "The counter keeps its previous progress."},
	Turn: `Start with the router, then send in-memory requests.

Compare status, body, and headers with the contract.

Also inspect the path with a trailing slash.`,
	Hints: []string{
		`Separate route handling from ownership of registry data.`,
		`Follow a request from handler to response and note each boundary.`,
		`Inspect the placeholder in ~NewBerthEngine~ and the ~gin.Context~ methods.`,
	},
}

var ch28FiberIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "saved", Value: "north"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "saved", Value: "north"}, {Name: "raw", Value: "sorth"}}},
	{Line: 5, Variables: []LessonVariable{{Name: "saved", Value: "north"}, {Name: "raw", Value: "sorth"}}, Output: "north sorth"},
}

var ch28FiberTryStates = []LessonFrame{
	{Line: 0},
	{Line: 3, Variables: []LessonVariable{{Name: "saved", Value: "east"}}},
	{Line: 4, Variables: []LessonVariable{{Name: "saved", Value: "east"}, {Name: "raw", Value: "west"}}},
	{Line: 5, Variables: []LessonVariable{{Name: "saved", Value: "east"}, {Name: "raw", Value: "west"}}, Output: "east"},
}

var ch28FiberArabic = ch28LessonCopy{
	Mission: `انقل نفس العقد لـ Fiber، واختبر التطبيق من غير فتح بورت.

الـstarter فيه handler مؤقت، وملف المقارنة لسه فاضي.`,
	Idea: `Fiber بيستخدم ~fiber.Ctx~ وبيشغّل الطلب عبر ~app.Test~.

بعض النصوص جاية من buffer بيتعاد استخدامه؛ انسخ اللي هتخزنه بعد الرد.`,
	IdeaCode: `func main() {
	raw := []byte("north")
	saved := strings.Clone(string(raw))
	raw[0] = 's'
	fmt.Println(saved, string(raw))
}`,
	IdeaImports: []string{"fmt", "strings"}, IdeaID: "e282-owned-string", IdeaStates: ch28FiberIdeaStates,
	IdeaCaptions: []string{"قبل التنفيذ.", "احتفظ بنسخة النص.", "البافر اتغير؛ النسخة ثابتة.", "اطبع النسختين."},
	TryPrompt:    `بعد تغيير أول بايت، إيه اللي بيظهر؟`,
	TryCode: `func main() {
	raw := []byte("east")
	saved := string(raw)
	raw[0] = 'w'
	fmt.Println(saved)
}`,
	TryImports: []string{"fmt"}, TryID: "e282-string-copy", TryAnswer: 0,
	TryOptions:  []string{"east", "west", "eest"},
	TryFeedback: []string{"تحويل الـslice لـ string عمل نسخة مستقلة.", "التعديل حصل على raw بعد النسخ.", "الحرف الأول المحفوظ لسه e."},
	TryStates:   ch28FiberTryStates,
	TryCaptions: []string{"قبل التنفيذ.", "saved اتنسخت.", "raw اتغيرت.", "اطبع النص المحفوظ."},
	Reveal:      `راقب إن تعديل الـbuffer ما يغيّرش النص المحفوظ.`,
	Watch: `ما تخزّنش ~c.Params~ بعد ما handler يرجع من غير نسخة.

JSON المفكوك بيملك نصوصه؛ راجع مصدر كل قيمة.`,
	Plan: `افصل شكل handler عن عقد البيانات اللي بتخدمه.

قارن الـresponse، اختبار الذاكرة، وسلامة القيم بعد انتهاء الطلب.

جهّز request جديد لكل محاولة عشان عمر السياق يبقى واضح.`,
	CheckQuestion: `إيه أداة Fiber اللي بتشغّل request في الذاكرة؟`,
	CheckOptions:  []string{"~app.Test~", "~http.ListenAndServe~", "~time.Sleep~"}, CheckAnswer: 0,
	CheckFeedback: []string{"~app.Test~ بيشغّل الطلب من غير port مفتوح.", "~http.ListenAndServe~ بيفتح listener حقيقي.", "~time.Sleep~ ما بيختبرش handler."},
	Turn: `طبّق عقد الأرصفة، وبعدها املأ المقارنة بجمل قصيرة.

جرّب المسارات بحروفها وشرطتها زي ما اتكتبت.`,
	Hints: []string{
		`افصل اختلاف handler عن ثبات عقد البيانات.`,
		`رتّب المقارنة حسب handler وerrors وtesting وlifetime.`,
		`راجع عمر ~c.Params~ و~strings.Clone~ قبل تخزين نص.`,
	},
}

var ch28FiberEnglish = ch28LessonCopy{
	Mission: `Move the same contract to Fiber and test the app without opening a port.

The starter has a temporary handler, and the comparison file is empty.`,
	Idea: `Fiber uses ~fiber.Ctx~ and runs requests through ~app.Test~.

Some strings come from a reused buffer; copy any value that must outlive the response.`,
	IdeaCode: ch28FiberArabic.IdeaCode, IdeaImports: ch28FiberArabic.IdeaImports, IdeaID: ch28FiberArabic.IdeaID, IdeaStates: ch28FiberIdeaStates,
	IdeaCaptions: []string{"Before execution.", "Keep a string copy.", "The buffer changed; the copy stayed.", "Print both values."},
	TryPrompt:    `After changing the first byte, what is printed?`,
	TryCode:      ch28FiberArabic.TryCode, TryImports: ch28FiberArabic.TryImports, TryID: ch28FiberArabic.TryID, TryAnswer: 0,
	TryOptions:  []string{"east", "west", "eest"},
	TryFeedback: []string{"Converting the slice to a string made a separate copy.", "The edit was made to raw after the copy.", "The saved first letter is still e."},
	TryStates:   ch28FiberTryStates,
	TryCaptions: []string{"Before execution.", "saved was copied.", "raw changed.", "Print the saved text."},
	Reveal:      `Watch the buffer edit leave the saved string unchanged.`,
	Watch: `Do not keep ~c.Params~ after a handler returns without copying it.

Decoded JSON owns its strings; check where each value came from.`,
	Plan: `Separate the handler shape from the data contract it serves.

Compare the response, in-memory test, and value lifetime after a request.

Use a fresh request per attempt to keep lifetimes clear.`,
	CheckQuestion: `Which Fiber tool runs a request in memory?`,
	CheckOptions:  []string{"~app.Test~", "~http.ListenAndServe~", "~time.Sleep~"}, CheckAnswer: 0,
	CheckFeedback: []string{"~app.Test~ runs the request without an open port.", "~http.ListenAndServe~ opens a real listener.", "~time.Sleep~ does not test a handler."},
	Turn: `Implement the berth contract, then fill the comparison with short sentences.

Try paths with their exact case and slash.`,
	Hints: []string{
		`Separate handler differences from the stable data contract.`,
		`Organize the comparison around handler, errors, testing, and lifetime.`,
		`Check the lifetime of ~c.Params~ and ~strings.Clone~ before storing text.`,
	},
}

var ch28OpenAPIIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2},
	{Line: 6, Variables: []LessonVariable{{Name: "object", Value: "map with depth"}}},
	{Line: 7, Variables: []LessonVariable{{Name: "depth", Value: "3"}, {Name: "number", Value: "true"}}},
	{Line: 8, Variables: []LessonVariable{{Name: "depth", Value: "3"}, {Name: "number", Value: "true"}}, Output: "true 3 true"},
}

var ch28OpenAPITryStates = []LessonFrame{
	{Line: 0},
	{Line: 2},
	{Line: 6, Variables: []LessonVariable{{Name: "object", Value: "map with string depth"}}},
	{Line: 7, Variables: []LessonVariable{{Name: "ok", Value: "false"}}},
	{Line: 8, Variables: []LessonVariable{{Name: "ok", Value: "false"}}, Output: "false"},
}

var ch28OpenAPIArabic = ch28LessonCopy{
	Mission: `اقرأ مستند OpenAPI لسجل المرحّلات، وراجع ردود الـAPI قدامه.

الـstarter مقسّم بين محلّل مستند، checker، وhandler.`,
	Idea: `JSON اللي شكله مش معروف بيتفك في ~any~؛ الأرقام تبقى ~float64~.

استخدم type assertion بشكل آمن. المستند يصف المسارات والردود والـschema.`,
	IdeaCode:    "func main() {\n\tvar value any\n\tif err := json.Unmarshal([]byte(`{\"depth\": 3}`), &value); err != nil {\n\t\tpanic(err)\n\t}\n\tobject, ok := value.(map[string]any)\n\tdepth, number := object[\"depth\"].(float64)\n\tfmt.Println(ok, depth, number)\n}",
	IdeaImports: []string{"encoding/json", "fmt"}, IdeaID: "e283-json-shape", IdeaStates: ch28OpenAPIIdeaStates,
	IdeaCaptions: []string{"قبل التنفيذ.", "ابدأ بقيمة عامة.", "الجذر object.", "الرقم اتعرف بأمان.", "اطبع النوع والقيمة."},
	TryPrompt:    `توقّع قيمة ~ok~ لما depth نص مش رقم.`,
	TryCode:      "func main() {\n\tvar value any\n\tif err := json.Unmarshal([]byte(`{\"depth\": \"3\"}`), &value); err != nil {\n\t\tpanic(err)\n\t}\n\tobject := value.(map[string]any)\n\t_, ok := object[\"depth\"].(float64)\n\tfmt.Println(ok)\n}",
	TryImports:   []string{"encoding/json", "fmt"}, TryID: "e283-type-check", TryAnswer: 1,
	TryOptions:  []string{"true", "false", "3"},
	TryFeedback: []string{"JSON decoded a string, not a number.", "The comma-ok assertion reports a type mismatch.", "The code prints the assertion result, not the field."},
	TryStates:   ch28OpenAPITryStates,
	TryCaptions: []string{"Before execution.", "Decode the object.", "The root is a map.", "The number assertion failed.", "Print the result."},
	Reveal:      `تابع فك الـJSON ونوع القيمة قبل الطباعة.`,
	Watch: `ماتعملش assertion مباشر من غير ~ok~؛ النوع الغلط ممكن يعمل panic.

المسار القالب لازم يطابق عدد الأجزاء.`,
	Plan: `افصل تحميل المستند عن مطابقة رد واحد.

راجع المسار والطريقة والحالة والنوع، وبعدها شكل الجسم.`,
	CheckQuestion: `لو response schema مفيهاش ~content~، الجسم المتوقع إيه؟`,
	CheckOptions:  []string{"جسم فاضي", "أي JSON", "نص افتراضي"}, CheckAnswer: 0,
	CheckFeedback: []string{"غياب content يعني إن الرد بلا body.", "الـschema مش موجودة عشان تقبل JSON.", "العقد ما بيضيفش نص تلقائي."},
	Turn: `ابدأ بمحلّل المستند، وبعدها اختبر مطابقة رد بسيط.

خلي بيانات المثال صغيرة، وامشِ من نوع القيمة للقاعدة.`,
	Hints: []string{
		`المستند والعينة الفعلية مصدرين مختلفين للمقارنة.`,
		`افصل تحميل الـschema عن فحص response واحد.`,
		`راجع ~json.Unmarshal~ ونتيجة type assertion عند ~LoadSpec~.`,
	},
}

var ch28OpenAPIEnglish = ch28LessonCopy{
	Mission: `Read the relay registry's OpenAPI document and compare API replies with it.

The starter separates document loading, a checker, and a handler.`,
	Idea: `JSON with an unknown shape decodes into ~any~; numbers become ~float64~.

Use a safe type assertion. The document describes paths, responses, and schemas.`,
	IdeaCode: ch28OpenAPIArabic.IdeaCode, IdeaImports: ch28OpenAPIArabic.IdeaImports, IdeaID: ch28OpenAPIArabic.IdeaID, IdeaStates: ch28OpenAPIIdeaStates,
	IdeaCaptions: []string{"Before execution.", "Start with a general value.", "The root is an object.", "The number was checked safely.", "Print type and value."},
	TryPrompt:    `Predict ~ok~ when depth is text, not a number.`,
	TryCode:      ch28OpenAPIArabic.TryCode, TryImports: ch28OpenAPIArabic.TryImports, TryID: ch28OpenAPIArabic.TryID, TryAnswer: 1,
	TryOptions:  []string{"true", "false", "3"},
	TryFeedback: []string{"JSON decoded a string, not a number.", "The comma-ok assertion reports a type mismatch.", "The code prints the assertion result, not the field."},
	TryStates:   ch28OpenAPITryStates,
	TryCaptions: []string{"Before execution.", "Decode the object.", "The root is a map.", "The number assertion failed.", "Print the result."},
	Reveal:      `Follow JSON decoding and the value type before printing.`,
	Watch: `Do not assert a type without ~ok~; a mismatch can panic.

The path template must match the same number of segments.`,
	Plan: `Separate loading the document from checking one response.

Check path, method, status, and media type, then the body shape.`,
	CheckQuestion: `If a response schema has no ~content~, what body is expected?`,
	CheckOptions:  []string{"An empty body", "Any JSON", "A default text"}, CheckAnswer: 0,
	CheckFeedback: []string{"No content means the response has no body.", "There is no schema there to accept JSON.", "The contract does not add default text."},
	Turn: `Start with the document loader, then check a small response.

Keep examples small and trace each value to its rule.`,
	Hints: []string{
		`The document and actual sample are separate sources to compare.`,
		`Separate schema loading from checking one response.`,
		`Inspect ~json.Unmarshal~ and safe type assertions around ~LoadSpec~.`,
	},
}

var ch28GracefulIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "shutdownErr", Value: "drain failed"}}},
	{Line: 4, Variables: []LessonVariable{{Name: "combined", Value: "joined errors"}}},
	{Line: 5, Variables: []LessonVariable{{Name: "combined", Value: "joined errors"}}, Output: "true"},
}

var ch28GracefulTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2},
	{Line: 3, Output: "true"},
	{Line: 4, Output: "true\nfalse"},
}

var ch28GracefulArabic = ch28LessonCopy{
	Mission: `اقفل البوابة بعد ما ترفض طلبات جديدة وتستنى الجاري بمهلة.

الـstarter فيه Gate ناقص؛ ركّز على عمر الخدمة والتنضيف.`,
	Idea: `~Serve~ ممكن يرجع أول ما الإغلاق يبدأ، قبل ما الطلب الجاري يخلص.

استخدم رجوع ~Shutdown~ لتعرف إن الانتظار انتهى، وبعدها اجمع أخطاء التنضيف.`,
	IdeaCode: `func main() {
	shutdownErr := errors.New("drain failed")
	releaseErr := errors.New("close failed")
	combined := errors.Join(shutdownErr, releaseErr)
	fmt.Println(errors.Is(combined, releaseErr))
}`,
	IdeaImports: []string{"errors", "fmt"}, IdeaID: "e284-joined-errors", IdeaStates: ch28GracefulIdeaStates,
	IdeaCaptions: []string{"قبل التنفيذ.", "خطأ الإغلاق موجود.", "اجمع أخطاء الإغلاق والتنضيف.", "هل خطأ التنضيف لسه معروف؟"},
	TryPrompt:    `توقّع نتيجتي محاولتي ~CompareAndSwap~.`,
	TryCode: `func main() {
	var used atomic.Bool
	fmt.Println(used.CompareAndSwap(false, true))
	fmt.Println(used.CompareAndSwap(false, true))
}`,
	TryImports: []string{"fmt", "sync/atomic"}, TryID: "e284-once-guard", TryAnswer: 1,
	TryOptions:  []string{"true ثم true", "true ثم false", "false ثم false"},
	TryFeedback: []string{"بعد التغيير الأول، القيمة مش false تاني.", "المحاولة الأولى غيّرت القيمة؛ التانية ما لقتش false.", "القيمة الأولية false فالمقارنة الأولى بتنجح."},
	TryStates:   ch28GracefulTryStates,
	TryCaptions: []string{"قبل التنفيذ.", "used تبدأ false.", "المقارنة الأولى تنجح.", "التانية تشوف true."},
	Reveal:      `نفّذ المقارنتين بالترتيب؛ المثال sequential.`,
	Watch: `ما تستخدمش context الإلغاء نفسه كمهلة drain.

المهلة تبدأ بعد الإشارة، والتنضيف ينتظر هدوء الخدمة.`,
	Plan: `افصل توقف listener عن انتهاء الطلبات الجارية.

بعدها راجع إن ~release~ يحصل مرة واحدة في كل طريق خروج.`,
	CheckQuestion: `إمتى يبقى آمن تقفل مورد بيستخدمه handler؟`,
	CheckOptions:  []string{"بعد ما الخدمة تهدى", "فور إلغاء context", "قبل ~Shutdown~"}, CheckAnswer: 0,
	CheckFeedback: []string{"الطلب الجاري لازم يخلص قبل التنضيف.", "الإلغاء بيبدأ الإغلاق، لكنه مش نهايته.", "المورد لازم يفضل متاح أثناء drain."},
	Turn: `ابدأ بمسار نجاح بسيط، وبعده جرّب الإلغاء والمهلة.

راجع ترتيب رجوع ~Serve~ و~Shutdown~ و~release~.`,
	Hints: []string{
		`الإغلاق الهادئ له بداية ووقت انتهاء مختلفين.`,
		`ارسم مسارين: Serve فشل، أو context طلب التوقف.`,
		`راجع نتائج ~Shutdown~ و~clock.After~ حول ~Run~.`,
	},
}

var ch28GracefulEnglish = ch28LessonCopy{
	Mission: `Close the gateway after refusing new requests and draining active ones within a grace period.

The starter has an unfinished Gate; focus on service lifetime and cleanup.`,
	Idea: `~Serve~ can return as shutdown starts, before an active request finishes.

Use shutdown completion to know when waiting ends, then join cleanup errors.`,
	IdeaCode: ch28GracefulArabic.IdeaCode, IdeaImports: ch28GracefulArabic.IdeaImports, IdeaID: ch28GracefulArabic.IdeaID, IdeaStates: ch28GracefulIdeaStates,
	IdeaCaptions: []string{"Before execution.", "The shutdown error exists.", "Join shutdown and cleanup errors.", "Is the cleanup error still known?"},
	TryPrompt:    `Predict the results of both ~CompareAndSwap~ calls.`,
	TryCode:      ch28GracefulArabic.TryCode, TryImports: ch28GracefulArabic.TryImports, TryID: ch28GracefulArabic.TryID, TryAnswer: 1,
	TryOptions:  []string{"true then true", "true then false", "false then false"},
	TryFeedback: []string{"After the first change, the value is no longer false.", "The first call changed it; the second no longer sees false.", "The initial value is false, so the first comparison succeeds."},
	TryStates:   ch28GracefulTryStates,
	TryCaptions: []string{"Before execution.", "used starts false.", "The first comparison succeeds.", "The second sees true."},
	Reveal:      `Run each comparison in order; this example is sequential.`,
	Watch: `Do not use the canceled request context as the drain deadline.

Grace starts after the stop signal, and cleanup waits for quiet.`,
	Plan: `Separate stopping the listener from finishing active requests.

Then check that ~release~ runs once on every exit path.`,
	CheckQuestion: `When is it safe to close a resource used by a handler?`,
	CheckOptions:  []string{"After serving is quiet", "As soon as context cancels", "Before ~Shutdown~"}, CheckAnswer: 0,
	CheckFeedback: []string{"An active request must finish before cleanup.", "Cancellation starts shutdown but is not its end.", "The resource must remain available while draining."},
	Turn: `Start with a normal exit, then try cancellation and grace expiry.

Check the order of ~Serve~, ~Shutdown~, and ~release~ returning.`,
	Hints: []string{
		`Graceful shutdown has a start and a separate finish.`,
		`Sketch two paths: Serve fails, or context requests a stop.`,
		`Review ~Shutdown~ and ~clock.After~ results around ~Run~.`,
	},
}

var ch28BossIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "w.Code", Value: "200"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "Content-Type", Value: "application/json"}}},
	{Line: 4, Variables: []LessonVariable{{Name: "w.Code", Value: "201"}, {Name: "Content-Type", Value: "application/json"}}},
	{Line: 5, Variables: []LessonVariable{{Name: "w.Code", Value: "201"}, {Name: "Content-Type", Value: "application/json"}}},
	{Line: 6, Variables: []LessonVariable{{Name: "w.Code", Value: "201"}, {Name: "Content-Type", Value: "application/json"}}, Output: "201 application/json"},
}

var ch28BossTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "status", Value: "403"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "status", Value: "403"}}, Output: "403"},
}

var ch28BossArabic = ch28LessonCopy{
	Mission: `ابني API لعزل الكبسولات بعقد واضح وتصاريح وحدود.

ده boss بيجمع راوتنج وHTTP وملكية حالة اتعلمتهم قبل كده.

كل server له مخزن مستقل، والـhandler بيتعامل معاه.`,
	Idea: `كل رد HTTP له status وheaders وbody.

اكتب headers قبل الجسم، وخلي التحقق يسبق أي تغيير في المخزن.`,
	IdeaCode: `func main() {
	w := httptest.NewRecorder()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "{\"state\":\"ready\"}")
	fmt.Println(w.Code, w.Header().Get("Content-Type"))
}`,
	IdeaImports: []string{"fmt", "net/http", "net/http/httptest"}, IdeaID: "e28b-response", IdeaStates: ch28BossIdeaStates,
	IdeaCaptions: []string{"قبل التنفيذ.", "جهّز recorder.", "الهيدر جاهز.", "اختار status للإنشاء.", "اكتب الجسم.", "اطبع بيانات الرد."},
	TryPrompt:    `توقّع رقم ~http.StatusForbidden~.`,
	TryCode: `func main() {
	status := http.StatusForbidden
	fmt.Println(status)
}`,
	TryImports: []string{"fmt", "net/http"}, TryID: "e28b-forbidden", TryAnswer: 1,
	TryOptions:  []string{"401", "403", "404"},
	TryFeedback: []string{"ده status فشل المصادقة.", "Forbidden رقمه 403.", "ده status للمسار الغايب."},
	TryStates:   ch28BossTryStates,
	TryCaptions: []string{"قبل التنفيذ.", "القيمة رقم 403.", "اطبع status."},
	Reveal:      `تابع الثابت من net/http لحد الخرج.`,
	Watch: `فرّق بين مصادقة فاشلة ودور غير مسموح.

حد الجسم قبل فك JSON؛ فشل الطلب ما يغيّرش المخزن.`,
	Plan: `رتّب العقد: مصادقة، route وطريقة، دور، body، ثم تخزين.

راجع كل رد خطأ وشوف إنه ما يسرّبش بيانات زيادة.`,
	CheckQuestion: `فين يتحط فحص credential بالنسبة لقراءة المخزن؟`,
	CheckOptions:  []string{"قبلها", "بعد ما يغيّر البيانات", "جوه main فقط"}, CheckAnswer: 0,
	CheckFeedback: []string{"الطلب غير المصرّح يتوقف قبل الوصول للحالة.", "التغيير الجانبي لازم يستنى التحقق.", "الـhandler هو اللي يعالج كل request."},
	Turn: `ابدأ بـ ~NewQuarantine~ من غير تشغيل سيرفر.

راجع كل boundary، وبعدها اتأكد إن كل خطأ يحافظ على الحالة.

السيرفر يبدأ بس من عند المنادي.`,
	Hints: []string{
		`خلّي المصادقة وملكية الكبسولات مسؤوليتين منفصلتين.`,
		`رتّب تحقق المدخلات قبل أي تعديل مشترك.`,
		`راجع placeholder في ~NewQuarantine~ وهيدر ~Authorization~.`,
	},
}

var ch28BossEnglish = ch28LessonCopy{
	Mission: `Build a capsule quarantine API with a clear contract, authorization, and bounds.

This boss combines routing, HTTP, and state ownership from earlier work.

Each server has an independent store used by its handler.`,
	Idea: `Every HTTP reply has a status, headers, and body.

Write headers before the body, and validate before changing stored state.`,
	IdeaCode: ch28BossArabic.IdeaCode, IdeaImports: ch28BossArabic.IdeaImports, IdeaID: ch28BossArabic.IdeaID, IdeaStates: ch28BossIdeaStates,
	IdeaCaptions: []string{"Before execution.", "Prepare a recorder.", "The header is ready.", "Choose a creation status.", "Write the body.", "Print response details."},
	TryPrompt:    `Predict the numeric value of ~http.StatusForbidden~.`,
	TryCode:      ch28BossArabic.TryCode, TryImports: ch28BossArabic.TryImports, TryID: ch28BossArabic.TryID, TryAnswer: 1,
	TryOptions:  []string{"401", "403", "404"},
	TryFeedback: []string{"That status means authentication failed.", "Forbidden has the value 403.", "That status is for a missing route."},
	TryStates:   ch28BossTryStates,
	TryCaptions: []string{"Before execution.", "The value is 403.", "Print the status."},
	Reveal:      `Follow the net/http constant to the output.`,
	Watch: `Distinguish failed authentication from a disallowed role.

Bound a body before decoding JSON; rejected requests leave storage unchanged.`,
	Plan: `Order the contract: authentication, route and method, role, body, then storage.

Review each error response and check that it reveals no extra data.`,
	CheckQuestion: `Where should credential checking happen relative to reading storage?`,
	CheckOptions:  []string{"Before it", "After changing data", "Only inside main"}, CheckAnswer: 0,
	CheckFeedback: []string{"An unauthorized request stops before reaching state.", "Side effects must wait until validation passes.", "The handler processes each request."},
	Turn: `Start with ~NewQuarantine~ without starting a server.

Check every boundary, then confirm errors preserve state.

The caller starts the server separately.`,
	Hints: []string{
		`Keep authentication and capsule ownership as separate responsibilities.`,
		`Order input checks before any shared-state mutation.`,
		`Inspect the placeholder in ~NewQuarantine~ and the ~Authorization~ header.`,
	},
}

var ch28GinArabicSteps = ch28BuildSteps("c281", false, ch28GinArabic)
var ch28GinEnglishSteps = ch28BuildSteps("c281", true, ch28GinEnglish)
var ch28FiberArabicSteps = ch28BuildSteps("c282", false, ch28FiberArabic)
var ch28FiberEnglishSteps = ch28BuildSteps("c282", true, ch28FiberEnglish)
var ch28OpenAPIArabicSteps = ch28BuildSteps("c283", false, ch28OpenAPIArabic)
var ch28OpenAPIEnglishSteps = ch28BuildSteps("c283", true, ch28OpenAPIEnglish)
var ch28GracefulArabicSteps = ch28BuildSteps("c284", false, ch28GracefulArabic)
var ch28GracefulEnglishSteps = ch28BuildSteps("c284", true, ch28GracefulEnglish)
var ch28BossArabicSteps = ch28BuildSteps("c28b", false, ch28BossArabic)
var ch28BossEnglishSteps = ch28BuildSteps("c28b", true, ch28BossEnglish)
