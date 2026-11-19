package levels

type ch26LessonCopy struct {
	mission       string
	idea          string
	tryPrompt     string
	tryReveal     string
	watch         string
	check         string
	checkQuestion string
	turn          string
	options       []string
	feedback      []string
	checkOptions  []string
	checkFeedback []string
	hints         []string
}

type ch26LessonSpec struct {
	prefix       string
	exampleID    string
	code         string
	imports      []string
	frames       [2][]LessonFrame
	tryExampleID string
	tryCode      string
	tryImports   []string
	tryFrames    [2][]LessonFrame
	tryAnswer    int
	checkAnswer  int
	copies       [2]ch26LessonCopy
}

type ch26BuiltLesson struct {
	arabic       []LessonStep
	english      []LessonStep
	arabicHints  []string
	englishHints []string
}

func ch26Pair(arabic, english string) map[string]string {
	return map[string]string{"ar": arabic, "en": english}
}

func ch26Sandbox(prefix, language string) LessonBlock {
	kind, config := ch26SandboxConfig(prefix)
	return LessonBlock{ID: prefix + ".sandbox", Kind: "sandbox", Text: ch26SandboxText(prefix, language), Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch26SandboxText(prefix, language string) string {
	translations := map[string]map[string]string{
		"c261": ch26Pair("قدّم الخطوات؛ تابع header وstatus وbody.", "Advance the steps; watch header, status, and body."),
		"c262": ch26Pair("أنشئ سجلًا؛ جرّب ترتيب قفل مختلفًا.", "Create records; try another lock order."),
		"c263": ch26Pair("قدّم الطلب؛ شوف دخول وخروج الـmiddleware.", "Advance the request through middleware entry and return."),
		"c264": ch26Pair("قدّم الـerror؛ تابع الوصول للنوع.", "Advance the error; follow its type."),
		"c265": ch26Pair("اختار المسار؛ راقب المطابقة الدقيقة.", "Choose a route; check exact matching."),
		"c266": ch26Pair("انسخ البايتات واحدة واحدة؛ وصل لـEOF.", "Copy bytes one at a time; reach EOF."),
	}
	return translations[prefix][language]
}

func ch26SandboxConfig(prefix string) (string, map[string]any) {
	switch prefix {
	case "c261":
		return "ch26-handler-flow", ch26HandlerConfig()
	case "c262":
		return "ch26-rest-atomic", ch26RestConfig()
	case "c263":
		return "ch26-middleware-order", ch26MiddlewareConfig()
	case "c264":
		return "ch26-error-envelope", ch26ErrorConfig()
	case "c265":
		return "ch26-static-route", ch26StaticConfig()
	case "c266":
		return "ch26-upload-stream", ch26UploadConfig()
	}
	return "", nil
}

func ch26HandlerConfig() map[string]any {
	return map[string]any{
		"frames": []any{
			map[string]any{"caption": ch26Pair("استجابة لسه فاضية.", "Response starts empty."), "header": "—", "status": "—", "body": "—"},
			map[string]any{"caption": ch26Pair("اتحدد Content-Type.", "Content-Type is set."), "header": "text/plain", "status": "—", "body": "—"},
			map[string]any{"caption": ch26Pair("اتحدد status قبل الكتابة.", "Status is set before writing."), "header": "application/demo", "status": 207, "body": "—"},
			map[string]any{"caption": ch26Pair("الـbody اتكتب أخيرًا.", "The body is written last."), "header": "application/demo", "status": 207, "body": "sample"},
		},
		"labels": map[string]any{
			"header": ch26Pair("Header", "Header"), "status": ch26Pair("Status", "Status"),
			"body": ch26Pair("Body", "Body"), "previous": ch26Pair("اللي قبلها", "Previous"),
			"next": ch26Pair("الخطوة الجاية", "Next step"), "reset": ch26Pair("من الأول", "Reset"),
		},
	}
}

func ch26RestConfig() map[string]any {
	return map[string]any{
		"nextID": 0,
		"records": []any{
			map[string]any{"name": ch26Pair("رصيف ألفا", "Quay Alpha")},
			map[string]any{"name": ch26Pair("رصيف بيتا", "Quay Beta")},
		},
		"orders": [][]int{{0, 1}, {1, 0}},
		"labels": map[string]any{
			"nextID": ch26Pair("آخر ID", "Last ID"), "records": ch26Pair("السجلات", "Records"),
			"next": ch26Pair("أنشئ سجلًا", "Create record"), "again": ch26Pair("بدّل ترتيب القفل", "Change lock order"),
			"reset": ch26Pair("من الأول", "Reset"), "order": ch26Pair("ترتيب طلبات القفل ممكن يتغير.", "Lock request order can vary."),
		},
	}
}

func ch26MiddlewareConfig() map[string]any {
	outer := ch26Pair("outer →", "outer →")
	inner := ch26Pair("inner →", "inner →")
	handler := ch26Pair("handler", "handler")
	innerReturn := ch26Pair("← inner", "← inner")
	outerReturn := ch26Pair("← outer", "← outer")
	return map[string]any{
		"frames": []any{
			map[string]any{"caption": ch26Pair("الطلب دخل outer.", "Request enters outer."), "trace": []any{outer}},
			map[string]any{"caption": ch26Pair("الطلب دخل inner.", "Request enters inner."), "trace": []any{outer, inner}},
			map[string]any{"caption": ch26Pair("الـhandler نفّذ الطلب.", "Handler runs the request."), "trace": []any{outer, inner, handler}},
			map[string]any{"caption": ch26Pair("inner رجع الأول.", "Inner returns first."), "trace": []any{outer, inner, handler, innerReturn}},
			map[string]any{"caption": ch26Pair("outer رجع في الآخر.", "Outer returns last."), "trace": []any{outer, inner, handler, innerReturn, outerReturn}},
		},
		"labels": map[string]any{
			"next": ch26Pair("الخطوة الجاية", "Next step"), "previous": ch26Pair("اللي قبلها", "Previous"),
			"reset": ch26Pair("من الأول", "Reset"),
		},
	}
}

func ch26ErrorConfig() map[string]any {
	fault := ch26Pair("*Fault", "*Fault")
	kind := ch26Pair("Kind: sample", "Kind: sample")
	load := ch26Pair("load", "load")
	return map[string]any{
		"frames": []any{
			map[string]any{"caption": ch26Pair("الخطأ typed.", "The error has a type."), "chain": []any{fault, kind}, "result": ch26Pair("قبل التغليف", "Before wrapping")},
			map[string]any{"caption": ch26Pair("اتغلف مع %w.", "Wrapped with %w."), "chain": []any{load, fault, kind}, "result": ch26Pair("النوع لسه جوه السلسلة", "Type remains in the chain")},
			map[string]any{"caption": ch26Pair("errors.As وصل للنوع.", "errors.As finds the type."), "chain": []any{load, fault, kind}, "result": ch26Pair("لقينا *Fault", "Found *Fault")},
		},
		"labels": map[string]any{
			"chain": ch26Pair("سلسلة الخطأ", "Error chain"), "result": ch26Pair("النتيجة", "Result"),
			"next": ch26Pair("الخطوة الجاية", "Next step"), "previous": ch26Pair("اللي قبلها", "Previous"),
			"reset": ch26Pair("من الأول", "Reset"),
		},
	}
}

func ch26StaticConfig() map[string]any {
	return map[string]any{
		"allowedPath": "/sample", "nearPath": "/sample/",
		"labels": map[string]any{
			"route": ch26Pair("المسار", "Route"), "allowed": ch26Pair("مسموح", "Allowed"),
			"denied": ch26Pair("مرفوض", "Denied"), "choose": ch26Pair("اختار مسارًا", "Choose a route"),
		},
	}
}

func ch26UploadConfig() map[string]any {
	return map[string]any{
		"payload": "sun", "count": 0,
		"labels": map[string]any{
			"source": ch26Pair("الملف التوضيحي", "Sample file"), "copied": ch26Pair("بايتات اتنسخت", "Bytes copied"),
			"staging": ch26Pair("منطقة التجهيز", "Staging area"), "eof": ch26Pair("وصلنا EOF", "Reached EOF"),
			"next": ch26Pair("اقرأ التالي", "Read next"), "empty": ch26Pair("فاضي", "empty"), "reset": ch26Pair("من الأول", "Reset"),
		},
	}
}

func ch26Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch26Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch26Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{
		ID:          id,
		Kind:        "walkthrough",
		Text:        text,
		Code:        code,
		Imports:     append([]string(nil), imports...),
		ExampleID:   exampleID,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch26BuildLesson(spec ch26LessonSpec) ch26BuiltLesson {
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	result := ch26BuiltLesson{}
	stepSets := []*[]LessonStep{&result.arabic, &result.english}
	hintSets := []*[]string{&result.arabicHints, &result.englishHints}
	for language, copy := range spec.copies {
		prefix := spec.prefix
		titles := arabicTitles
		languageCode := "ar"
		if language == 1 {
			titles = englishTitles
			languageCode = "en"
		}
		tryText := copy.tryPrompt + "\n\n~~~go\n" + spec.tryCode + "\n~~~"
		steps := []LessonStep{
			{
				ID: prefix + ".mission-step", Title: titles[0],
				Blocks: []LessonBlock{ch26Markdown(prefix+".mission", copy.mission)},
			},
			{
				ID: prefix + ".idea-step", Title: titles[1],
				Blocks: []LessonBlock{ch26Walkthrough(prefix+".idea-trace", copy.idea, spec.code, spec.exampleID, spec.imports, spec.frames[language])},
			},
			{
				ID: prefix + ".try-step", Title: titles[2],
				Blocks: []LessonBlock{
					ch26Question(prefix+".predict", tryText, copy.options, spec.tryAnswer, copy.feedback),
					ch26Walkthrough(prefix+".try-reveal", copy.tryReveal, spec.tryCode, spec.tryExampleID, spec.tryImports, spec.tryFrames[language]),
					ch26Sandbox(prefix, languageCode),
				},
			},
			{
				ID: prefix + ".watch-step", Title: titles[3],
				Blocks: []LessonBlock{ch26Markdown(prefix+".watch", copy.watch)},
			},
			{
				ID: prefix + ".check-step", Title: titles[4],
				Blocks: []LessonBlock{
					ch26Markdown(prefix+".check", copy.check),
					ch26Question(prefix+".check-question", copy.checkQuestion, copy.checkOptions, spec.checkAnswer, copy.checkFeedback),
				},
			},
			{
				ID: prefix + ".turn-step", Title: titles[5],
				Blocks: []LessonBlock{ch26Markdown(prefix+".turn", copy.turn)},
			},
		}
		*stepSets[language] = steps
		*hintSets[language] = append([]string(nil), copy.hints...)
	}
	return result
}

var ch26Lesson261 = ch26BuildLesson(ch26LessonSpec{
	prefix: "c261", exampleID: "e261-handler",
	imports: []string{"fmt", "net/http", "net/http/httptest", "strings"},
	code: `func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("queued"))
	})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/queue", nil))
	fmt.Println(w.Code, w.Header().Get("Content-Type"), strings.TrimSpace(w.Body.String()))
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الـhandler اتجهز."},
			{Line: 3, Caption: "الهيدر اتحدد بدري."},
			{Line: 4, Caption: "اختار الحالة قبل الجسم."},
			{Line: 5, Caption: "اكتب الجسم بعد الحالة."},
			{Line: 9, Caption: "سجّل الرد كله.", Output: "202 text/plain queued"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The handler is ready."},
			{Line: 3, Caption: "Set the header first."},
			{Line: 4, Caption: "Choose status before body."},
			{Line: 5, Caption: "Write the body afterward."},
			{Line: 9, Caption: "Record the full response.", Output: "202 text/plain queued"},
		},
	},
	tryExampleID: "e261-head-check",
	tryImports:   []string{"fmt", "net/http", "net/http/httptest"},
	tryCode: `func main() {
	w := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.ServeHTTP(w, httptest.NewRequest("GET", "/ready", nil))
	fmt.Println(w.Code, w.Body.Len())
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "جهّز مسجّل الرد."},
			{Line: 4, Caption: "الحالة من غير جسم."},
			{Line: 6, Caption: "نفّذ الطلب داخل الذاكرة."},
			{Line: 7, Caption: "الرد فاضي فعلًا.", Output: "204 0"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Prepare a response recorder."},
			{Line: 4, Caption: "Status without a body."},
			{Line: 6, Caption: "Run the request in memory."},
			{Line: 7, Caption: "The response is empty.", Output: "204 0"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch26LessonCopy{
		{
			mission:       "بوابة القاعدة محتاجة صفحة ونبض جاهزية.\n\nالـstarter مديك signatures بس؛ انت هتحدد الردود وحدود السيرفر.\n\nابدأ بالـhandler، وخلي الاختبار في الذاكرة.",
			idea:          "الـhandler بياخد request ويكتب response.\n\nرتّب الهيدر، الحالة، وبعدها الجسم. أول كتابة ممكن تثبّت الحالة تلقائيًا.\n\n~httptest.NewRecorder~ يسجّل الرد من غير بورت.",
			tryPrompt:     "توقّع ناتج مسجّل رد فاضي:",
			tryReveal:     "شوف الحالة وطول الجسم بعد التنفيذ.",
			watch:         "- بعد ~WriteHeader~ الحالة اتثبتت؛ رتّب كل الهيدرز قبله.\n- القالب يهرب البيانات، و~Server~ هنا إعداد بس.",
			check:         "قسّم الشغل: جهّز handler، جرّب طلبات في الذاكرة، وبعدها افحص إعداد السيرفر.\n\nخلي أخطاء التنفيذ برّه الرد قبل ما تكتب أي بايت.",
			checkQuestion: "فين تضبط ~Content-Type~ عشان يطلع في الرد؟",
			turn:          "اكتب جزء صغير وشغّله بـ~httptest~. راجع الحالة والهيدر والجسم قبل ما تكمل.",
			options:       []string{"200 0", "204 0", "204 1"},
			feedback:      []string{"مفيش كتابة عادية تختار 200 هنا.", "الحالة 204 والجسم طوله صفر.", "مافيش بايت اتكتب في الجسم."},
			checkOptions:  []string{"قبل أول كتابة أو حالة", "بعد الجسم", "بعد انتهاء handler"},
			checkFeedback: []string{"الهيدر لازم يتجهز قبل تثبيت الرد.", "الكتابة بدأت إرسال الرد بالفعل.", "الرد خلص، مش هتعدّل هيدره."},
			hints:         []string{"افصل اختيار المسار عن اختيار طريقة الطلب.", "جهّز handler صغيرًا واختبره بمسجل داخل الذاكرة.", "راجع ~ResponseRecorder~ قبل إضافة إعدادات السيرفر."},
		},
		{
			mission:       "The base gateway needs a page and a readiness pulse.\n\nThe starter gives you signatures; you shape responses and server limits.\n\nStart with a handler and keep checks in memory.",
			idea:          "A handler receives a request and writes a response.\n\nSet headers, then status, then body. The first write can choose status implicitly.\n\n~httptest.NewRecorder~ records a response without a port.",
			tryPrompt:     "Predict the empty recorder result:",
			tryReveal:     "Watch the status and body length after execution.",
			watch:         "- After ~WriteHeader~, status is committed; set headers first.\n- Templates escape data; ~Server~ is configuration here.",
			check:         "Split the work: build a handler, try in-memory requests, then inspect server settings.\n\nHandle render errors before writing any response bytes.",
			checkQuestion: "When should ~Content-Type~ be set?",
			turn:          "Write a small piece and run it with ~httptest~. Check status, header, and body before continuing.",
			options:       []string{"200 0", "204 0", "204 1"},
			feedback:      []string{"No ordinary write selects 200 here.", "The status is 204 and its body is empty.", "No body byte was written."},
			checkOptions:  []string{"Before the first write or status", "After the body", "After the handler returns"},
			checkFeedback: []string{"Prepare the header before committing the response.", "Writing already started the response.", "The response is finished; its header is too late."},
			hints:         []string{"Separate path selection from request-method selection.", "Build a small handler and check it with an in-memory recorder.", "Inspect ~ResponseRecorder~ before adding server settings."},
		},
	},
})

var ch26Lesson262 = ch26BuildLesson(ch26LessonSpec{
	prefix: "c262", exampleID: "e262-locked-id",
	imports: []string{"fmt", "sync"},
	code: `type Berth struct { Label string }
func main() {
	var mu sync.Mutex
	nextID := 0
	berths := map[int]Berth{}
	mu.Lock()
	nextID++
	berths[nextID] = Berth{Label: "west"}
	id := nextID
	label := berths[id].Label
	mu.Unlock()
	fmt.Println(id, label)
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 4, Caption: "المخزن لسه فاضي.", Variables: []LessonVariable{{Name: "nextID", Value: "0"}}},
			{Line: 5, Caption: "القرار جوه القفل."},
			{Line: 9, Caption: "خد نسخة الهوية.", Variables: []LessonVariable{{Name: "id", Value: "1"}}},
			{Line: 10, Caption: "اقرأ الاسم تحت القفل.", Variables: []LessonVariable{{Name: "label", Value: "west"}}},
			{Line: 11, Caption: "افتح القفل بعد النسخ."},
			{Line: 12, Caption: "اطبع الهوية والاسم.", Output: "1 west"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 4, Caption: "The store is empty.", Variables: []LessonVariable{{Name: "nextID", Value: "0"}}},
			{Line: 5, Caption: "The decision is locked."},
			{Line: 9, Caption: "Keep a local ID copy.", Variables: []LessonVariable{{Name: "id", Value: "1"}}},
			{Line: 10, Caption: "Read the label while locked.", Variables: []LessonVariable{{Name: "label", Value: "west"}}},
			{Line: 11, Caption: "Unlock after copying values."},
			{Line: 12, Caption: "Print the ID and label.", Output: "1 west"},
		},
	},
	tryExampleID: "e262-sort-copy",
	tryImports:   []string{"fmt", "slices"},
	tryCode: `func main() {
	ids := []int{8, 3, 5}
	slices.Sort(ids)
	fmt.Println(ids)
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "القائمة بالترتيب الأصلي."},
			{Line: 3, Caption: "Sort رتّبت الأرقام."},
			{Line: 4, Caption: "اطبع الترتيب الجديد.", Output: "[3 5 8]"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The original order is shown."},
			{Line: 3, Caption: "Sort orders the numbers."},
			{Line: 4, Caption: "Print the new order.", Output: "[3 5 8]"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch26LessonCopy{
		{
			mission:       "سجل الأرصفة محتاج إنشاء وقراءة وتعديل وحذف.\n\nالـstarter فيه نوع ~Berth~ بس؛ الـAPI والمخزن لسه عليك.\n\nخلي الـIDs ثابتة، واحمِ الحالة المشتركة.",
			idea:          "العداد والخريطة قرار واحد وقت الإنشاء.\n\nالقفل يخلي التحديث المركب يحصل كخطوة محمية من طلب تاني.\n\nللقوائم، انسخ البيانات قبل فك القفل ورتّب النسخة.",
			tryPrompt:     "توقّع القائمة بعد ~slices.Sort~:",
			tryReveal:     "تابع العناصر قبل وبعد الترتيب.",
			watch:         "- طول الخريطة مش هوية جديدة؛ الحذف ممكن يعيد نفس العدد.\n- متقرأش map مشتركة بعد فك القفل.",
			check:         "فرّق بين تحليل المسار، فحص الجسم، وتغيير المخزن.\n\nالطلب المرفوض المفروض يسيب الحالة زي ما هي.",
			checkQuestion: "إيه اللي يخلي ID مايتكرّرش بعد حذف رصيف؟",
			turn:          "ابدأ بمخزن صغير واختبر فكرته. بعدين اربطه بطرق HTTP واحدة واحدة.",
			options:       []string{"[8 3 5]", "[3 5 8]", "[5 3 8]"},
			feedback:      []string{"Sort غيّرت ترتيب العناصر.", "الترتيب التصاعدي هو 3 ثم 5 ثم 8.", "النتيجة مش ترتيبًا تصاعديًا."},
			checkOptions:  []string{"عداد مستقل يزيد مع الإنشاء", "استخدام len للخريطة", "اشتقاق الرقم من ترتيب JSON"},
			checkFeedback: []string{"العداد يحتفظ بالتاريخ حتى بعد الحذف.", "عدد العناصر ممكن يقل ويرجع يتكرر.", "ترتيب العرض مش هوية مستقرة."},
			hints:         []string{"عامل الهوية كحالة مستقلة عن عدد السجلات الحالية.", "خلي القراءة والتعديل ضمن حدود القفل نفسه.", "راجع نسخة القائمة قبل فرزها أو ترميزها."},
		},
		{
			mission:       "The berth registry needs create, read, update, and delete.\n\nThe starter provides ~Berth~ only; the API and store are yours.\n\nKeep IDs stable and protect shared state.",
			idea:          "The counter and map form one decision during creation.\n\nA lock makes the combined update a protected step against another request.\n\nFor lists, copy data before unlocking and sort the copy.",
			tryPrompt:     "Predict the list after ~slices.Sort~:",
			tryReveal:     "Follow the elements before and after sorting.",
			watch:         "- Map length is not a new identity; deletion can repeat the count.\n- Do not read a shared map after unlocking.",
			check:         "Separate path parsing, body checks, and store changes.\n\nA rejected request should leave state untouched.",
			checkQuestion: "What keeps an ID from repeating after deletion?",
			turn:          "Start with a small store and check its behavior. Then connect HTTP methods one at a time.",
			options:       []string{"[8 3 5]", "[3 5 8]", "[5 3 8]"},
			feedback:      []string{"Sort changes the element order.", "Ascending order is 3, then 5, then 8.", "That is not ascending order."},
			checkOptions:  []string{"A separate counter that advances on creation", "Using the map length", "Deriving it from JSON order"},
			checkFeedback: []string{"The counter remembers history after deletion.", "The item count can shrink and repeat.", "Display order is not stable identity."},
			hints:         []string{"Treat identity as separate from the number of current records.", "Keep reads and mutations inside the same lock boundary.", "Inspect a list copy before sorting or encoding it."},
		},
	},
})

var ch26Lesson263 = ch26BuildLesson(ch26LessonSpec{
	prefix: "c263", exampleID: "e263-middleware-order",
	imports: []string{"fmt", "net/http", "net/http/httptest", "strings"},
	code: `func main() {
	var trace []string
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { trace = append(trace, "handler") })
	wrap := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				trace = append(trace, name+">")
				next.ServeHTTP(w, r)
				trace = append(trace, "<"+name)
			})
		}
	}
	h := wrap("outer")(wrap("inner")(handler))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/bay/7", nil))
	fmt.Println(strings.Join(trace, " "))
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 3, Caption: "الـhandler يسجل أثره."},
			{Line: 13, Caption: "الأولى تفضل برّه."},
			{Line: 14, Caption: "الطلب دخل السلسلة.", Variables: []LessonVariable{{Name: "trace", Value: "[outer> inner> handler <inner <outer]"}}},
			{Line: 15, Caption: "اطبع الترتيب الفعلي.", Output: "outer> inner> handler <inner <outer"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 3, Caption: "The handler records its turn."},
			{Line: 13, Caption: "The first wrapper stays outside."},
			{Line: 14, Caption: "The request enters the chain.", Variables: []LessonVariable{{Name: "trace", Value: "[outer> inner> handler <inner <outer]"}}},
			{Line: 15, Caption: "Print the observed order.", Output: "outer> inner> handler <inner <outer"},
		},
	},
	tryExampleID: "e263-path-parts",
	tryImports:   []string{"fmt", "strings"},
	tryCode: `func main() {
	parts := strings.Split("/bays/7", "/")
	fmt.Println(len(parts), parts[2])
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "المسار بقى أجزاء."},
			{Line: 3, Caption: "الجزء الأخير هو 7.", Output: "3 7"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The path became segments."},
			{Line: 3, Caption: "The last segment is 7.", Output: "3 7"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch26LessonCopy{
		{
			mission:       "الراوتر بيوصل طلبات لأرصفة حسب المسار والطريقة.\n\nالـstarter فاضي؛ هتضيف matching وmiddleware وبيانات لكل طلب.\n\nالمسار الغلط يفضل غلط، من غير تنظيف.",
			idea:          "كل middleware بتلف handler تاني.\n\nالأولى تبقى برّه، فتشتغل قبل اللي بعدها وبعدها بالعكس.\n\nحط params وrequest ID في ~context~ مشتق للطلب، مش في متغير عام.",
			tryPrompt:     "توقّع عدد الأجزاء وقيمة آخر جزء:",
			tryReveal:     "اقرأ ناتج تقسيم المسار.",
			watch:         "- parameter بياخد جزء واحد غير فاضي.\n- ~WithContext~ يرجّع نسخة ويحافظ على سياق الأصل.",
			check:         "فرّق بين مطابقة النمط ومطابقة الطريقة.\n\nاستخدم ~context~ بقيمة مستقلة لكل request عشان الطلبات ما تتسربش لبعض.",
			checkQuestion: "فين تخزن params الخاصة بطلب واحد؟",
			turn:          "جرّب route toy وميدلوير يسجلوا ترتيبهم. بعدها افحص الحالة في الطلب الأصلي.",
			options:       []string{"2 و7", "3 و7", "3 وbays"},
			feedback:      []string{"التقسيم بيشمل الجزء الفاضي قبل أول slash.", "فيه 3 أجزاء وآخر قيمة هي 7.", "آخر جزء هو 7، مش bays."},
			checkOptions:  []string{"في ~context~ مشتق للطلب", "في global map", "في ~URL~ query"},
			checkFeedback: []string{"النسخة تحتفظ بقيم الطلب من غير تسريب.", "طلبات متزامنة ممكن تخلط القيم.", "الـquery مش مكان حالة الراوتر هنا."},
			hints:         []string{"قارن أجزاء المسار مع النمط قبل تقرير الطريقة.", "ركّب السلسلة من الداخل للخارج عشان أول عنصر يبقى خارجيًا.", "راجع مفاتيح ~context~ الخاصة وWithContext على نسخة الطلب."},
		},
		{
			mission:       "The router sends requests to berths by path and method.\n\nThe starter is empty; add matching, middleware, and per-request data.\n\nAn invalid path should stay invalid instead of being cleaned.",
			idea:          "Each middleware wraps another handler.\n\nThe first middleware is outermost, so it runs before the rest and returns last.\n\nPut parameters and request IDs in request ~context~, not a global variable.",
			tryPrompt:     "Predict the segment count and last value:",
			tryReveal:     "Read the split path result.",
			watch:         "- A parameter matches one nonempty segment.\n- ~WithContext~ returns a copy and keeps the parent ~context~.",
			check:         "Separate pattern matching from method matching.\n\nUse request ~context~ for values that belong to one request, so requests cannot leak into each other.",
			checkQuestion: "Where should one request's parameters live?",
			turn:          "Try a toy route and middleware that records order. Then inspect the original request.",
			options:       []string{"2 and 7", "3 and 7", "3 and bays"},
			feedback:      []string{"The split includes the empty part before the first slash.", "There are 3 parts and the last value is 7.", "The last part is 7, not bays."},
			checkOptions:  []string{"In that request's derived ~context~", "In a global map", "In the ~URL~ query"},
			checkFeedback: []string{"A derived request carries its value without leaking it.", "Concurrent requests can mix values.", "The router state does not belong in the query."},
			hints:         []string{"Compare path segments with the pattern before deciding the method.", "Wrap from inside out so the first middleware stays outermost.", "Review private ~context~ keys and ~WithContext~ on a request copy."},
		},
	},
})

var ch26Lesson264 = ch26BuildLesson(ch26LessonSpec{
	prefix: "c264", exampleID: "e264-wrapped-fault",
	imports: []string{"errors", "fmt"},
	code: `type Fault struct { Kind string }
func (f *Fault) Error() string { return "private detail" }
func main() {
	wrapped := fmt.Errorf("fetch: %w", &Fault{Kind: "missing"})
	var fault *Fault
	found := errors.As(wrapped, &fault)
	fmt.Println(found, fault.Kind)
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 1, Caption: "النوع شايل التصنيف."},
			{Line: 4, Caption: "الخطأ اتلفّ."},
			{Line: 6, Caption: "As لقت النوع.", Variables: []LessonVariable{{Name: "found", Value: "true"}, {Name: "fault.Kind", Value: "missing"}}},
			{Line: 7, Caption: "اطبع التصنيف بس.", Output: "true missing"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 1, Caption: "The type carries the kind."},
			{Line: 4, Caption: "The error is wrapped."},
			{Line: 6, Caption: "As finds the type.", Variables: []LessonVariable{{Name: "found", Value: "true"}, {Name: "fault.Kind", Value: "missing"}}},
			{Line: 7, Caption: "Print the category only.", Output: "true missing"},
		},
	},
	tryExampleID: "e264-json-null",
	tryImports:   []string{"encoding/json", "fmt"},
	tryCode: `func main() {
	reply := map[string]any{"data": nil, "error": nil}
	raw, _ := json.Marshal(reply)
	fmt.Println(string(raw))
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الحقلين موجودين."},
			{Line: 3, Caption: "JSON احتفظ بـnull.", Output: "{\"data\":null,\"error\":null}"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Both fields are present."},
			{Line: 3, Caption: "JSON keeps both nulls.", Output: "{\"data\":null,\"error\":null}"},
		},
	},
	tryAnswer: 1, checkAnswer: 1,
	copies: [2]ch26LessonCopy{
		{
			mission:       "ردود المنارات محتاجة شكل ثابت وتصنيف آمن للأخطاء.\n\nالـstarter سايب تعريفات المصدر؛ هتكتب نوع الخطأ والـhandler.\n\nرسالة التشخيص تفضل داخلية. أمثلة الدرس للتوضيح مش ردود المهمة.",
			idea:          "الـerror ممكن يتلفّ مع سياق إضافي.\n\n~errors.As~ تدور على النوع جوه السلسلة، و~Kind~ يحدد تصنيفًا معروفًا.\n\nاستخدم رسالة عامة للعميل، وسيب السبب للتشخيص. التصنيف مش نص الرسالة.",
			tryPrompt:     "توقّع JSON لخريطة بالقيمتين ~nil~:",
			tryReveal:     "شوف شكل الحقلين بعد الترميز.",
			watch:         "- ~nil~ مع pointer يظهر ~null~؛ ~omitempty~ ممكن يخفي الحقل.\n- متحوّلش ~err.Error()~ لنص رد عام.",
			check:         "حدد نوع الخطأ أولًا، وبعدها اختر الحالة والرسالة العامة.\n\nفي النجاح والفشل، شكل الظرف يفضل ثابت.",
			checkQuestion: "إيه مصدر الرسالة العامة في الرد؟",
			turn:          "جرّب خطأ ملفوفًا. تأكد إن التصنيف ينجو وإن التفاصيل الخاصة ماتظهرش.",
			options:       []string{"{\"data\":null,\"error\":null}", "{\"data\":{},\"error\":[]}", "{}"},
			feedback:      []string{"القيمتان الفارغتان اتكتبتا كـnull.", "نوع nil interface مش بيتحوّل للشكل ده.", "الخريطة فيها حقلين، مش object فاضي."},
			checkOptions:  []string{"تصنيف موثوق ورسالة ثابتة", "نص Cause", "نص wrapper كامل"},
			checkFeedback: []string{"الرد يشرح الفئة من غير تفاصيل داخلية.", "السبب ممكن يكشف بيانات خاصة.", "الـwrapper ممكن يسرّب سياق داخلي."},
			hints:         []string{"فرّق بين نوع الخطأ ونص التشخيص اللي يحمله.", "خلّي بيانات الرد والخطأ اختيارية من غير إخفاء الحقول.", "ابدأ بـerrors.As قبل اختيار التصنيف العام."},
		},
		{
			mission:       "Relay responses need one stable shape and safe error classification.\n\nThe starter leaves the source types; you build the fault and handler.\n\nDiagnostic text stays internal. Lesson values are illustrative, not response data.",
			idea:          "An error can be wrapped with extra ~context~.\n\n~errors.As~ searches the chain for a type, and ~Kind~ selects a known category.\n\nGive clients a fixed public message; keep the cause for diagnosis. Classification is separate from message text.",
			tryPrompt:     "Predict JSON for a map with two ~nil~ values:",
			tryReveal:     "See how both fields look after encoding.",
			watch:         "- A nil pointer becomes ~null~; ~omitempty~ can hide a field.\n- Do not turn ~err.Error()~ into a public response.",
			check:         "Identify the error type first, then choose a status and fixed message.\n\nKeep the envelope shape stable on success and failure.",
			checkQuestion: "What should supply the public message?",
			turn:          "Try a wrapped error. Check that its category survives and private details stay hidden.",
			options:       []string{"{\"data\":null,\"error\":null}", "{\"data\":{},\"error\":[]}", "{}"},
			feedback:      []string{"Both nil values encode as null.", "A nil interface does not encode to those values.", "The map has two fields, not an empty object."},
			checkOptions:  []string{"A trusted category and fixed message", "The cause text", "The full wrapper text"},
			checkFeedback: []string{"The response describes a category without internals.", "The cause may expose private data.", "The wrapper may leak internal ~context~."},
			hints:         []string{"Separate an error's type from its diagnostic text.", "Keep response data and error optional without omitting their fields.", "Start with ~errors.As~ before selecting the public category."},
		},
	},
})

var ch26Lesson265 = ch26BuildLesson(ch26LessonSpec{
	prefix: "c265", exampleID: "e265-embedded-fs",
	imports: []string{"fmt", "io/fs", "testing/fstest"},
	code: `func main() {
	site := fstest.MapFS{"charts/notice.txt": &fstest.MapFile{Data: []byte("sector")}}
	charts, err := fs.Sub(site, "charts")
	body, readErr := fs.ReadFile(charts, "notice.txt")
	fmt.Println(string(body), err == nil, readErr == nil)
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الملف جوه نظام افتراضي."},
			{Line: 3, Caption: "Sub اختارت جذر charts."},
			{Line: 4, Caption: "اقرأ بايتات الملف."},
			{Line: 5, Caption: "اطبع المحتوى والحالتين.", Output: "sector true true"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The file is in a virtual filesystem."},
			{Line: 3, Caption: "Sub selects the charts root."},
			{Line: 4, Caption: "Read the file bytes."},
			{Line: 5, Caption: "Print content and both checks.", Output: "sector true true"},
		},
	},
	tryExampleID: "e265-exact-route",
	tryImports:   []string{"fmt"},
	tryCode: `func main() {
	allowed := map[string]bool{"/console": true}
	_, exact := allowed["/console"]
	_, slash := allowed["/console/"]
	fmt.Println(exact, slash)
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "المسار الدقيق مسموح."},
			{Line: 3, Caption: "الشرطة الزيادة مش موجودة."},
			{Line: 4, Caption: "اطبع النتيجتين.", Output: "true false"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The exact path is allowed."},
			{Line: 3, Caption: "The extra slash is absent."},
			{Line: 4, Caption: "Print both results.", Output: "true false"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch26LessonCopy{
		{
			mission:       "بوابة الخرائط تخدم ملفات ثابتة من داخل البرنامج.\n\nالـstarter معاه الأصول؛ انت هتحدد اللي يتشاف وسياسة الكاش.\n\nماتفتحش مجلد للزوار.",
			idea:          "الـfilesystem المضمّن بيتقرأ من غير قرص وقت التشغيل.\n\n~fs.Sub~ يختار جذرًا، و~fs.ReadFile~ يجيب البايتات.\n\nقائمة المسارات الدقيقة تحط الحراسة قبل ~FileServer~.",
			tryPrompt:     "توقّع نتيجة المسار الدقيق والمسار بشرطة زيادة:",
			tryReveal:     "قارن المفتاحين في خريطة السماح.",
			watch:         "- احسب البصمة من نفس البايتات اللي هتخدمها.\n- انسخ request قبل تغيير ~URL~ أو الهيدرز.",
			check:         "ابدأ بقائمة مسارات واضحة، وبعدها وصل كل مسار بالملف ونوعه.\n\nراجع ~HEAD~ والكاش من غير ما تغيّر الطلب الأصلي.",
			checkQuestion: "إيه اللي تمنعه قائمة المسارات الدقيقة؟",
			turn:          "ابدأ بقراءة ملف toy من filesystem داخل الذاكرة. بعدها جرّب مسار مش مسموح.",
			options:       []string{"true true", "true false", "false true"},
			feedback:      []string{"المفتاح التاني مختلف بسبب الشرطة الزيادة.", "المطابقة حرفية؛ المفتاح الأول بس موجود.", "المسار الأصلي موجود في الخريطة."},
			checkOptions:  []string{"عرض مجلد أو مسار غير معلن", "تغيير محتوى الملف", "حساب بصمة مختلفة"},
			checkFeedback: []string{"الـallowlist توقف التفويض قبل العرض.", "القائمة ما بتغيّرش بايتات الملف.", "المطابقة مش مسؤولة عن حساب البصمة."},
			hints:         []string{"عامل كل مسار عام كاختيار صريح من قائمة قصيرة.", "اقرأ المحتوى من filesystem بدل مسار الجهاز.", "قارن ~RawPath~ قبل إعادة كتابة ~URL~."},
		},
		{
			mission:       "The chart gateway serves static files from inside the program.\n\nThe starter includes assets; you choose what is visible and the cache policy.\n\nDo not expose a directory to visitors.",
			idea:          "An embedded filesystem reads files without runtime disk access.\n\n~fs.Sub~ selects a root, and ~fs.ReadFile~ returns the bytes.\n\nAn exact route list guards the call to ~FileServer~.",
			tryPrompt:     "Predict the exact path and extra-slash results:",
			tryReveal:     "Compare both keys in the allowlist.",
			watch:         "- Hash the same bytes you are about to serve.\n- Clone the request before changing its ~URL~ or headers.",
			check:         "Start with an explicit route list, then connect each path to a file and content type.\n\nCheck ~HEAD~ and caching without changing the original request.",
			checkQuestion: "What does an exact route list prevent?",
			turn:          "Read a toy file from an in-memory filesystem. Then try a path that is not allowed.",
			options:       []string{"true true", "true false", "false true"},
			feedback:      []string{"The second key differs because of the extra slash.", "Matching is exact; only the first key exists.", "The original path is in the map."},
			checkOptions:  []string{"Listing a directory or serving an undeclared path", "Changing the file contents", "Calculating a different hash"},
			checkFeedback: []string{"The allowlist stops delegation before serving.", "The list does not change file bytes.", "Path matching does not calculate the hash."},
			hints:         []string{"Treat each public path as an explicit choice from a short list.", "Read assets from a filesystem instead of a machine path.", "Compare ~RawPath~ before rewriting the ~URL~."},
		},
	},
})

var ch26Lesson266 = ch26BuildLesson(ch26LessonSpec{
	prefix: "c266", exampleID: "e266-stream-copy",
	imports: []string{"fmt", "io", "strings"},
	code: `func main() {
	source := strings.NewReader("packet")
	n, err := io.Copy(io.Discard, source)
	fmt.Println(n, err)
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "مصدر بيقرا على دفعات."},
			{Line: 3, Caption: "انسخ لحد نهاية المصدر.", Variables: []LessonVariable{{Name: "n", Value: "6"}, {Name: "err", Value: "nil"}}},
			{Line: 4, Caption: "اعرض عدد البايتات والخطأ.", Output: "6 <nil>"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The source can be read in chunks."},
			{Line: 3, Caption: "Copy through source end.", Variables: []LessonVariable{{Name: "n", Value: "6"}, {Name: "err", Value: "nil"}}},
			{Line: 4, Caption: "Show the byte count and error.", Output: "6 <nil>"},
		},
	},
	tryExampleID: "e266-stream-eof",
	tryImports:   []string{"fmt", "io", "strings"},
	tryCode: `func main() {
	reader := strings.NewReader("x")
	_, _ = reader.Read(make([]byte, 1))
	_, err := reader.Read(make([]byte, 1))
	fmt.Println(err == io.EOF)
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "القارئ فيه بايت واحد."},
			{Line: 3, Caption: "اقرأ البايت الموجود."},
			{Line: 4, Caption: "القراءة الجاية وصلت للنهاية.", Output: "true"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The reader holds one byte."},
			{Line: 3, Caption: "Read the available byte."},
			{Line: 4, Caption: "The next read reaches the end.", Output: "true"},
		},
	},
	tryAnswer: 0, checkAnswer: 1,
	copies: [2]ch26LessonCopy{
		{
			mission:       "بوابة الشحن تستقبل لقطة multipart وتخزنها بأمان.\n\nالـstarter فيه الـtypes بس؛ انت هتكتب المعالج والتخزين.\n\nالتعامل مع stream يحمي الذاكرة والتنظيف يحمي الملفات.",
			idea:          "الـpart قارئ؛ انسخه تدريجيًا بدل ما تحمّله كله.\n\nجهّز ملف مؤقت، افحص البايتات، وبعدين انشره بالاسم الداخلي.\n\nفي كل فشل، نظّف الأثر المؤقت.",
			tryPrompt:     "توقّع هل القراءة التانية وصلت لنهاية المصدر:",
			tryReveal:     "تابع حالة القارئ بعد استهلاك البايت.",
			watch:         "- اسم العميل مش مسار تخزين؛ افحص القيمة الخام.\n- جسم الطلب والملف ليهم حدين منفصلين.",
			check:         "افصل قراءة الطلب عن التحقق والنشر.\n\nخلي الوصول للوجهة وخطوة النشر محميين معًا، واربط التنظيف بكل خروج.",
			checkQuestion: "ليه نقرأ زيادة بايت عن حد الملف؟",
			turn:          "جرّب stream صغيرًا في الذاكرة. لاحظ وقت EOF، وبعدها خطط لمرحلة الملف المؤقت.",
			options:       []string{"لا، لسه فيه بيانات", "أيوه، وصلت EOF", "القراءة بتعيد نفس البايت"},
			feedback:      []string{"البايت الوحيد اتقري بالفعل.", "المصدر خلص بعد أول قراءة.", "القارئ بيتقدم بعد القراءة."},
			checkOptions:  []string{"عشان نكتشف إن البيانات عدّت الحد", "عشان نضيف بايت للملف", "عشان نعرف نوع الصورة"},
			checkFeedback: []string{"البايت الإضافي يفرّق الحد من الزيادة.", "مش بنحفظ بايت زيادة مقصودًا.", "التوقيع هو اللي يحدد النوع المقبول."},
			hints:         []string{"فكّر في كل جزء كقارئ له نهاية وخطأ محتمل.", "افصل ملف مؤقت عن الاسم النهائي المنشور.", "ابدأ تطبيق ~MaxBytesReader~ قبل إنشاء ~MultipartReader~."},
		},
		{
			mission:       "The cargo gateway receives a multipart scan and stores it safely.\n\nThe starter has the types; you build the handler and storage flow.\n\nStreaming protects memory, and cleanup protects files.",
			idea:          "A part is a reader; copy it in chunks instead of loading it all.\n\nStage a temporary file, inspect bytes, then publish it under an internal name.\n\nClean up the temporary artifact on every failure.",
			tryPrompt:     "Predict whether the second read reached source end:",
			tryReveal:     "Follow the reader after its byte is consumed.",
			watch:         "- A client filename is not a storage path; inspect the raw value.\n- The request body and file have separate limits.",
			check:         "Separate request reading, validation, and publication.\n\nProtect destination checks and publication together, and connect cleanup to every exit.",
			checkQuestion: "Why read one byte beyond the file limit?",
			turn:          "Try a small in-memory stream. Watch for EOF, then plan the temporary-file stage.",
			options:       []string{"No, data remains", "Yes, it reached EOF", "The reader repeats the same byte"},
			feedback:      []string{"The only byte was already consumed.", "The source ended after the first read.", "The reader advances when read."},
			checkOptions:  []string{"To detect data beyond the limit", "To add an extra byte to the file", "To identify the image type"},
			checkFeedback: []string{"One extra byte distinguishes the limit from overflow.", "The extra byte is not meant to be saved.", "The signature decides which type is accepted."},
			hints:         []string{"Treat each part as a reader with an end and possible error.", "Separate a temporary file from the published name.", "Apply ~MaxBytesReader~ before creating ~MultipartReader~."},
		},
	},
})
