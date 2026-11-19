package levels

type ch27LessonCopy struct {
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

type ch27LessonSpec struct {
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
	copies       [2]ch27LessonCopy
}

type ch27BuiltLesson struct {
	arabic       []LessonStep
	english      []LessonStep
	arabicHints  []string
	englishHints []string
}

func ch27Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch27Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch27Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

type ch27SandboxSpec struct {
	id     string
	kind   string
	text   [2]string
	config map[string]any
}

func ch27Pair(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

var ch27SandboxSpecs = map[string]ch27SandboxSpec{
	"c271": {
		id: "ch27-signed-pass-playground", kind: "ch27-signed-pass",
		text: [2]string{"غيّر الرسالة أو المفتاح وشوف ~tag~.", "Change the message or key and watch the tag."},
		config: map[string]any{
			"labels": map[string]any{
				"payload":       ch27Pair("الـpayload", "Payload"),
				"key":           ch27Pair("المفتاح", "Key"),
				"tag":           ch27Pair("tag توضيحي", "Toy tag"),
				"changeMessage": ch27Pair("غيّر الرسالة", "Change message"),
				"rotateKey":     ch27Pair("بدّل المفتاح", "Rotate key"),
				"caption":       ch27Pair("الرسالة مقروءة والـtag يتغير.", "The payload stays readable; the tag changes."),
			},
			"messages": []string{"crew=Mira", "crew=Vega"},
			"keys":     []string{"key-A", "key-B"},
			"tags":     []string{"tag-01", "tag-02", "tag-03", "tag-04"},
		},
	},
	"c272": {
		id: "ch27-sessions-playground", kind: "ch27-session-lookup",
		text: [2]string{"بدّل ~ID~ وشوف سجل السيرفر.", "Switch IDs and inspect the server record."},
		config: map[string]any{
			"labels": map[string]any{
				"browser": ch27Pair("مع المتصفح", "Browser holds"),
				"server":  ch27Pair("نتيجة السيرفر", "Server lookup"),
				"next":    ch27Pair("بدّل الـID", "Switch ID"),
			},
			"lookups": []any{
				map[string]any{"id": "sess-A7", "result": "Mira", "caption": ch27Pair("السجل ربط الـID بالاسم.", "The record maps the ID to a name.")},
				map[string]any{"id": "sess-X2", "result": "—", "caption": ch27Pair("مفيش سجل بالـID ده.", "No record matches this ID.")},
			},
		},
	},
	"c273": {
		id: "ch27-cors-playground", kind: "ch27-preflight",
		text: [2]string{"بدّل الطلب وتابع استدعاء الـhandler.", "Switch requests and follow the handler call."},
		config: map[string]any{
			"labels": map[string]any{
				"method":     ch27Pair("الطريقة", "Method"),
				"stage":      ch27Pair("الفرع", "Path"),
				"status":     ch27Pair("الرد", "Response"),
				"innerCalls": ch27Pair("نداءات handler", "Handler calls"),
				"next":       ch27Pair("الطلب التالي", "Next request"),
			},
			"requests": []any{
				map[string]any{"method": "OPTIONS", "stage": "preflight", "status": "204", "calls": "0", "caption": ch27Pair("الفرع ده ما ناداش handler.", "This branch skipped the handler.")},
				map[string]any{"method": "GET", "stage": "route", "status": "200", "calls": "1", "caption": ch27Pair("الطلب العادي وصل للـhandler.", "The normal request reached the handler.")},
			},
		},
	},
	"c274": {
		id: "ch27-rate-limit-playground", kind: "ch27-client-counts",
		text: [2]string{"عدّ الطلبات لكل عنوان لوحده.", "Count requests separately for each address."},
		config: map[string]any{
			"labels": map[string]any{
				"client":  ch27Pair("العميل", "Client"),
				"count":   ch27Pair("الطلبات", "Requests"),
				"caption": ch27Pair("كل عنوان له عدّاده.", "Each address has its own count."),
			},
			"clients": []any{
				map[string]any{"address": "198.51.100.8", "button": ch27Pair("طلب من A", "Request A")},
				map[string]any{"address": "198.51.100.9", "button": ch27Pair("طلب من B", "Request B")},
			},
		},
	},
	"c275": {
		id: "ch27-short-links-playground", kind: "ch27-base62-index",
		text: [2]string{"اختار ~index~ وشوف الرمز المقابل.", "Pick an index and see its symbol."},
		config: map[string]any{
			"labels": map[string]any{
				"index":    ch27Pair("index", "Index"),
				"symbol":   ch27Pair("الرمز", "Symbol"),
				"previous": ch27Pair("السابق", "Previous"),
				"next":     ch27Pair("التالي", "Next"),
				"reset":    ch27Pair("من الأول", "Reset"),
			},
			"lookups": []any{
				map[string]any{"index": "0", "symbol": "0", "caption": ch27Pair("أول خانة تبدأ بصفر.", "The first slot starts at zero.")},
				map[string]any{"index": "10", "symbol": "a", "caption": ch27Pair("بعد الأرقام تبدأ الحروف.", "Letters follow the digits.")},
				map[string]any{"index": "61", "symbol": "Z", "caption": ch27Pair("الترتيب يحدد الرمز.", "The order selects the symbol.")},
			},
		},
	},
}

func ch27SandboxBlock(prefix string, english bool) LessonBlock {
	spec := ch27SandboxSpecs[prefix]
	text := spec.text[0]
	if english {
		text = spec.text[1]
	}
	return LessonBlock{
		ID: spec.id, Kind: "sandbox", Text: text,
		Sandbox: &LessonSandbox{Kind: spec.kind, Config: spec.config},
	}
}

func ch27BuildLesson(spec ch27LessonSpec) ch27BuiltLesson {
	arabicTitles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	englishTitles := []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	result := ch27BuiltLesson{}
	stepSets := []*[]LessonStep{&result.arabic, &result.english}
	hintSets := []*[]string{&result.arabicHints, &result.englishHints}
	for language, copy := range spec.copies {
		prefix := spec.prefix
		titles := arabicTitles
		if language == 1 {
			titles = englishTitles
		}
		tryText := copy.tryPrompt + "\n\n~~~go\n" + spec.tryCode + "\n~~~"
		steps := []LessonStep{
			{
				ID: prefix + ".mission-step", Title: titles[0],
				Blocks: []LessonBlock{ch27Markdown(prefix+".mission", copy.mission)},
			},
			{
				ID: prefix + ".idea-step", Title: titles[1],
				Blocks: []LessonBlock{ch27Walkthrough(prefix+".idea-trace", copy.idea, spec.code, spec.exampleID, spec.imports, spec.frames[language])},
			},
			{
				ID: prefix + ".try-step", Title: titles[2],
				Blocks: []LessonBlock{
					ch27Question(prefix+".predict", tryText, copy.options, spec.tryAnswer, copy.feedback),
					ch27Walkthrough(prefix+".try-reveal", copy.tryReveal, spec.tryCode, spec.tryExampleID, spec.tryImports, spec.tryFrames[language]),
					ch27SandboxBlock(prefix, language == 1),
				},
			},
			{
				ID: prefix + ".watch-step", Title: titles[3],
				Blocks: []LessonBlock{ch27Markdown(prefix+".watch", copy.watch)},
			},
			{
				ID: prefix + ".check-step", Title: titles[4],
				Blocks: []LessonBlock{
					ch27Markdown(prefix+".check", copy.check),
					ch27Question(prefix+".check-question", copy.checkQuestion, copy.checkOptions, spec.checkAnswer, copy.checkFeedback),
				},
			},
			{
				ID: prefix + ".turn-step", Title: titles[5],
				Blocks: []LessonBlock{ch27Markdown(prefix+".turn", copy.turn)},
			},
		}
		*stepSets[language] = steps
		*hintSets[language] = append([]string(nil), copy.hints...)
	}
	return result
}

var ch27Lesson271 = ch27BuildLesson(ch27LessonSpec{
	prefix: "c271", exampleID: "e271-hmac-demo",
	imports: []string{"crypto/hmac", "crypto/sha256", "fmt"},
	code: `func main() {
	message := []byte("crew=Mira")
	mac := hmac.New(sha256.New, []byte("demo-only"))
	_, _ = mac.Write(message)
	fmt.Printf("%x\n", mac.Sum(nil)[:4])
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الرسالة مثال توضيحي."},
			{Line: 3, Caption: "جهّز MAC بالمفتاح."},
			{Line: 4, Caption: "اربط الرسالة بالمفتاح."},
			{Line: 5, Caption: "اطبع جزءًا من البصمة.", Output: "b3614ce6"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The message is illustrative."},
			{Line: 3, Caption: "Prepare a MAC with its key."},
			{Line: 4, Caption: "Bind the message to the key."},
			{Line: 5, Caption: "Print part of the tag.", Output: "b3614ce6"},
		},
	},
	tryExampleID: "e271-readable-payload",
	tryImports:   []string{"encoding/base64", "fmt"},
	tryCode: `func main() {
	payload := base64.RawURLEncoding.EncodeToString([]byte("crew=Mira"))
	decoded, _ := base64.RawURLEncoding.DecodeString(payload)
	fmt.Println(string(decoded))
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "حوّل النص لـbase64url."},
			{Line: 3, Caption: "رجّعه لبايتات."},
			{Line: 4, Caption: "اقرأ النص الأصلي.", Output: "crew=Mira"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Encode text as base64url."},
			{Line: 3, Caption: "Decode it back to bytes."},
			{Line: 4, Caption: "Read the original text.", Output: "crew=Mira"},
		},
	},
	tryAnswer: 0, checkAnswer: 2,
	copies: [2]ch27LessonCopy{
		{
			mission:       "بوابة المساكن تفحص code وتطلع تصريح قصير العمر.\n\nالـstarter فيه signatures بس؛ هتربط تخزين آمن بـbcrypt وتحقق ~JWT~.\n\nالسر ييجي من مصدر إعدادات.",
			idea:          "الـHMAC بيربط الرسالة بمفتاح؛ أي تعديل يخلي التحقق يفشل.\n\nده مثال لفكرة توقيع ~JWT~، مش بديل عن ~bcrypt~ لكود الطاقم.\n\nالـpayload ممكن يتقري؛ التوقيع هو اللي يثبت سلامته.",
			tryPrompt:     "توقّع النص بعد فك ترميز payload:",
			tryReveal:     "تابع الترميز وفك الترميز.",
			watch:         "- ~base64url~ ترميز، مش تشفير ولا إثبات هوية.\n- قارن ~bcrypt~ errors: الكود الغلط غير عطل hash.",
			check:         "قسّم الشغل بين تحميل السر، إصدار التصريح، والتحقق منه.\n\nاستخدم الساعة المحقونة عشان الوقت يفضل قابل للفحص.",
			checkQuestion: "إيه اللي يثبت إن claims ما اتعدلتش؟",
			turn:          "ابدأ بفحص بيانات دخول تجريبية. بعدين تتبع claims قبل وبعد التحقق.",
			options:       []string{"crew=Mira", "Y3Jldz1NaXJh", "Mira=crew"},
			feedback:      []string{"فك الترميز يعيد البايتات الأصلية.", "ده النص المشفّر بترميز ~base64url~.", "الترتيب في الرسالة ما اتبدلش."},
			checkOptions:  []string{"إن payload نصه مقروء", "وجود اسم crew", "توقيع صالح بخوارزمية مسموحة"},
			checkFeedback: []string{"القراءة لا تثبت إن البيانات موثوقة.", "الاسم وحده ممكن يتزور.", "التحقق من التوقيع والخوارزمية يثبت السلامة."},
			hints:         []string{"فرّق بين hash كود الدخول وتوقيع التصريح.", "افصل تحميل السر وLogin عن middleware ~Require~.", "راجع خيارات ~ParseWithClaims~ الخاصة بالمصدر والوقت."},
		},
		{
			mission:       "The quarters gate checks a crew code and issues a short-lived pass.\n\nThe starter has signatures; connect safe storage with ~bcrypt~ and ~JWT~ verification.\n\nLoad the secret from configuration.",
			idea:          "~HMAC~ binds a message to a key; changing the message makes verification fail.\n\nThis illustrates a ~JWT~ signature, not a replacement for ~bcrypt~ on crew codes.\n\nA payload can be read; the signature proves its integrity.",
			tryPrompt:     "Predict the text after decoding the payload:",
			tryReveal:     "Follow encoding and decoding.",
			watch:         "- ~base64url~ is encoding, not encryption or proof of identity.\n- Separate a wrong ~bcrypt~ code from a broken stored hash.",
			check:         "Separate secret loading, pass issuance, and pass verification.\n\nUse the injected clock so time stays testable.",
			checkQuestion: "What proves the claims were not changed?",
			turn:          "Start by checking sample login data. Then trace claims before and after verification.",
			options:       []string{"crew=Mira", "Y3Jldz1NaXJh", "Mira=crew"},
			feedback:      []string{"Decoding restores the original bytes.", "That is the ~base64url~-encoded text.", "The message order did not change."},
			checkOptions:  []string{"A readable payload", "A crew name is present", "A valid signature with an allowed method"},
			checkFeedback: []string{"Readable content does not make it trusted.", "A name alone can be forged.", "Checking signature and method proves integrity."},
			hints:         []string{"Separate the login-code hash from the signed pass.", "Keep secret loading, ~Login~, and ~Require~ as distinct jobs.", "Review ~ParseWithClaims~ options for issuer and time."},
		},
	},
})

var ch27Lesson272 = ch27BuildLesson(ch27LessonSpec{
	prefix: "c272", exampleID: "e272-cookie-flags",
	imports: []string{"fmt", "net/http"},
	code: `func main() {
	cookie := http.Cookie{Name: "quarters", Value: "demo", Path: "/",
		HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode}
	fmt.Println(cookie.HttpOnly, cookie.Secure, cookie.SameSite == http.SameSiteStrictMode)
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الكوكي له خصائص صريحة.", Variables: []LessonVariable{{Name: "cookie.Name", Value: "quarters"}, {Name: "cookie.Path", Value: "/"}}},
			{Line: 3, Caption: "افحص خصائص الحماية.", Variables: []LessonVariable{{Name: "cookie.HttpOnly", Value: "true"}, {Name: "cookie.Secure", Value: "true"}}},
			{Line: 3, Caption: "القيم الثلاثة شغالة.", Output: "true true true"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The cookie sets explicit attributes.", Variables: []LessonVariable{{Name: "cookie.Name", Value: "quarters"}, {Name: "cookie.Path", Value: "/"}}},
			{Line: 3, Caption: "Check the protection flags.", Variables: []LessonVariable{{Name: "cookie.HttpOnly", Value: "true"}, {Name: "cookie.Secure", Value: "true"}}},
			{Line: 3, Caption: "All three checks are true.", Output: "true true true"},
		},
	},
	tryExampleID: "e272-session-lookup",
	tryImports:   []string{"fmt"},
	tryCode: `func main() {
	sessions := map[string]string{"demo-id": "Mira"}
	crew, ok := sessions["demo-id"]
	fmt.Println(crew, ok)
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الخريطة فيها جلسة."},
			{Line: 3, Caption: "المفتاح لقى الاسم."},
			{Line: 4, Caption: "اطبع الاسم ووجوده.", Output: "Mira true"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The map holds a session."},
			{Line: 3, Caption: "The key finds its name."},
			{Line: 4, Caption: "Print the name and lookup.", Output: "Mira true"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch27LessonCopy{
		{
			mission:       "الجلسة بتربط متصفح الطاقم بهويته وتمنع طلبات مزيفة.\n\nالـstarter يديك Quarters بس؛ هتدير الجلسات والكوكي وCSRF.\n\nوقت الاختبار ييجي من clock محقون.",
			idea:          "السيرفر يحتفظ بسجل الجلسة؛ المتصفح يحمل ID بس.\n\nخصائص الكوكي تقلل قراءته أو إرساله في سياق غلط.\n\nبعد دخول ناجح، الهوية تتجدد عشان ID قديم مايفضلش صالح.",
			tryPrompt:     "توقّع نتيجة البحث عن الجلسة:",
			tryReveal:     "شوف الاسم وقيمة comma-ok.",
			watch:         "- الجلسة المنتهية تتعامل كأنها مش موجودة.\n- طلب بيغيّر حالة محتاج CSRF token بجانب الكوكي.",
			check:         "افصل التحقق من بيانات الدخول عن البحث في الجلسة.\n\nراجع التدوير والخروج كخطوات تغيّر الحالة، وGET كقراءة.",
			checkQuestion: "إيه اللي يميّز طلبًا شرعيًا بيغيّر الحالة؟",
			turn:          "ابدأ بسجل جلسة toy. جرّب قراءة الاسم، وبعدها تتبع ما يتغير عند الخروج.",
			options:       []string{"Mira false", "Mira true", "demo-id true"},
			feedback:      []string{"المفتاح موجود فـok تساوي true.", "الاسم المخزن Mira والمفتاح اتلاقى.", "الـmap بترجع القيمة مش المفتاح."},
			checkOptions:  []string{"كوكي سليم بس", "CSRF token مطابق في الهيدر", "اسم الطاقم في query"},
			checkFeedback: []string{"المتصفح ممكن يرسل الكوكي تلقائيًا.", "الطلب يثبت امتلاكه للتوكن السري الإضافي.", "الـquery مش إثباتًا سريًا."},
			hints:         []string{"فرّق بين ID الموجود عند العميل وسجل الجلسة على السيرفر.", "رتّب التحقق والإنشاء والتدوير كعمليات واضحة.", "راجع ~ConstantTimeCompare~ عند مقارنة توكن CSRF."},
		},
		{
			mission:       "The session connects a crew browser to its identity and rejects forged requests.\n\nThe starter gives you Quarters; manage sessions, cookies, and CSRF.\n\nTests supply an injected clock.",
			idea:          "The server keeps a session record; the browser carries only its ID.\n\nCookie attributes limit who can read or send it in the wrong ~context~.\n\nAfter login, rotate identity so an old ID does not stay valid.",
			tryPrompt:     "Predict the session lookup result:",
			tryReveal:     "See the name and comma-ok value.",
			watch:         "- Treat an expired session as absent.\n- State-changing requests need a CSRF token beside the cookie.",
			check:         "Separate credential checking from session lookup.\n\nReview rotation and logout as state changes, and ~GET~ as a read.",
			checkQuestion: "What distinguishes a legitimate state-changing request?",
			turn:          "Start with a toy session record. Read its name, then trace what logout changes.",
			options:       []string{"Mira false", "Mira true", "demo-id true"},
			feedback:      []string{"The key exists, so ok is true.", "The stored name is Mira and the key was found.", "A map lookup returns the value, not the key."},
			checkOptions:  []string{"A valid cookie alone", "A matching CSRF token in the header", "A crew name in the query"},
			checkFeedback: []string{"The browser may send its cookie automatically.", "The request proves it has the additional secret token.", "The query is not secret proof."},
			hints:         []string{"Separate the client-held ID from the server-side session record.", "Lay out validation, creation, and rotation as distinct operations.", "Review ~ConstantTimeCompare~ for CSRF token comparison."},
		},
	},
})

var ch27Lesson273 = ch27BuildLesson(ch27LessonSpec{
	prefix: "c273", exampleID: "e273-preflight",
	imports: []string{"fmt", "net/http", "net/http/httptest"},
	code: `func main() {
	nextCalls := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusOK)
	})
	shield := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/data", nil)
	r.Header.Set("Access-Control-Request-Method", "GET")
	shield.ServeHTTP(w, r)
	fmt.Println(w.Code, nextCalls)
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الـhandler الداخلي جاهز."},
			{Line: 8, Caption: "الطلب التمهيدي له فرع."},
			{Line: 14, Caption: "جهّز مسجّل الرد."},
			{Line: 15, Caption: "جهّز الطلب التمهيدي."},
			{Line: 16, Caption: "نوع الطلب اتسجل."},
			{Line: 17, Caption: "الفرع رد قبل next."},
			{Line: 18, Caption: "الحالة وعدد النداءات.", Variables: []LessonVariable{{Name: "nextCalls", Value: "0"}}, Output: "204 0"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The inner handler is ready."},
			{Line: 8, Caption: "Preflight has its own branch."},
			{Line: 14, Caption: "Prepare the response recorder."},
			{Line: 15, Caption: "Prepare the preflight request."},
			{Line: 16, Caption: "The request method is set."},
			{Line: 17, Caption: "The branch responds before next."},
			{Line: 18, Caption: "Status and call count.", Variables: []LessonVariable{{Name: "nextCalls", Value: "0"}}, Output: "204 0"},
		},
	},
	tryExampleID: "e273-header-case",
	tryImports:   []string{"fmt", "strings"},
	tryCode: `func main() {
	fmt.Println(strings.EqualFold("X-Crew", "x-crew"))
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "قارن اسمي الهيدر."},
			{Line: 3, Caption: "حالة الحروف مش فارقة.", Output: "true"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Compare the header names."},
			{Line: 3, Caption: "Letter case does not matter.", Output: "true"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch27LessonCopy{
		{
			mission:       "درع البوابة يسمح لأصول الأسطول المعروفة تقرأ الرد.\n\nالـstarter فيه Policy وShield بس؛ هتضيف validation وسلوك wrapper.\n\nCORS مش تسجيل دخول.",
			idea:          "CORS بيتحكم في قراءة المتصفح لرد من origin تاني.\n\nالطلب العادي يكمل لـnext؛ preflight يسأل عن الإذن قبل الطلب.\n\nضيف ~Vary: Origin~ من غير ما تمسح قيم الـhandler.",
			tryPrompt:     "توقّع هل اسما الهيدر متساويين:",
			tryReveal:     "تابع المقارنة اللي تتجاهل حالة الحروف.",
			watch:         "- أسماء الهيدرات غير حساسة لحالة الحروف؛ methods حساسة.\n- الحماية و~Vary~ لازم يوصلوا لكل رد.",
			check:         "قسّم wrapper لطلب عادي وpreflight.\n\nاجمع شروط الأصل والطريقة والهيدرات، وسيب رد الـhandler الداخلي يمر.",
			checkQuestion: "CORS بيمنع مين من قراءة الرد؟",
			turn:          "جرّب wrapper مع recorder وعدّاد للـhandler. بدّل Origin وشوف الهيدرات.",
			options:       []string{"false", "true", "compile error"},
			feedback:      []string{"~EqualFold~ بتتجاهل فرق الحالة.", "الحروف مختلفة لكن الاسمين متساويين هنا.", "الاستدعاء API صحيح في strings."},
			checkOptions:  []string{"كل clients قبل وصولهم للسيرفر", "المتصفح من إتاحة الرد لصفحة خارجية", "الـhandler من استقبال الطلب"},
			checkFeedback: []string{"CORS مش جدار مصادقة على الاتصال.", "المتصفح يقرر السماح للصفحة بقراءة الرد.", "السيرفر يستقبل الطلب العادي."},
			hints:         []string{"عامل CORS كسياسة قراءة يطبقها المتصفح.", "افصل preflight عن الطلبات العادية داخل wrapper.", "استخدم ~Header.Add~ لحفظ Vary الموجود."},
		},
		{
			mission:       "The gateway shield lets known fleet origins read a response.\n\nThe starter has Policy and Shield; add validation and wrapper behavior.\n\nCORS is not login.",
			idea:          "CORS controls whether a browser page can read a cross-origin response.\n\nAn ordinary request reaches next; preflight asks permission first.\n\nAdd ~Vary: Origin~ without erasing values from the inner handler.",
			tryPrompt:     "Predict whether these header names compare equal:",
			tryReveal:     "Follow the comparison that ignores letter case.",
			watch:         "- Header names ignore case; methods do not.\n- Protection headers and ~Vary~ must reach every response.",
			check:         "Split the wrapper into ordinary and preflight requests.\n\nCombine origin, method, and header checks, while letting inner responses pass through.",
			checkQuestion: "Who does CORS stop from reading a response?",
			turn:          "Try a wrapper with a recorder and handler counter. Change Origin and inspect headers.",
			options:       []string{"false", "true", "compile error"},
			feedback:      []string{"~EqualFold~ ignores case differences.", "The letters differ, but these names compare equal.", "This is a valid strings API call."},
			checkOptions:  []string{"All clients before they reach the server", "The browser from exposing a response to a foreign page", "The handler from receiving a request"},
			checkFeedback: []string{"CORS is not a connection-authentication wall.", "The browser decides whether the page can read the response.", "The server still receives ordinary requests."},
			hints:         []string{"Treat CORS as a browser-enforced read policy.", "Separate preflight from ordinary requests inside the wrapper.", "Use ~Header.Add~ to preserve existing Vary values."},
		},
	},
})

var ch27Lesson274 = ch27BuildLesson(ch27LessonSpec{
	prefix: "c274", exampleID: "e274-client-tanks",
	imports: []string{"fmt"},
	code: `func main() {
	tanks := map[string]int{"192.0.2.4": 2, "192.0.2.9": 1}
	tanks["192.0.2.4"]--
	fmt.Println(tanks["192.0.2.4"], tanks["192.0.2.9"])
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "كل عنوان له رصيد."},
			{Line: 3, Caption: "رصيد عميل واحد قل."},
			{Line: 4, Caption: "الرصيد التاني زي ما هو.", Output: "1 1"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "Each address has a balance."},
			{Line: 3, Caption: "One client's balance drops."},
			{Line: 4, Caption: "The other stays unchanged.", Output: "1 1"},
		},
	},
	tryExampleID: "e274-round-up",
	tryImports:   []string{"fmt", "math", "time"},
	tryCode: `func main() {
	wait := 1500 * time.Millisecond
	seconds := int(math.Ceil(wait.Seconds()))
	fmt.Println(seconds)
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الانتظار ثانية ونص."},
			{Line: 3, Caption: "قرّب الثواني لفوق."},
			{Line: 4, Caption: "العميل ينتظر ثانيتين.", Output: "2"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The wait is one and a half seconds."},
			{Line: 3, Caption: "Round seconds upward."},
			{Line: 4, Caption: "The client waits two seconds.", Output: "2"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch27LessonCopy{
		{
			mission:       "الـgateway يدي كل عميل خزان تصاريح مستقل.\n\nالـstarter فيه حدود الإعداد؛ هتضيف مفاتيح العملاء والحساب والتنظيف.\n\nالوقت جاي من clock، مش نوم فعلي.",
			idea:          "خريطة الخزانات تفصل رصيد كل client.\n\nالـkey يعتمد على IP موثوق؛ بورت الاتصال بيتغير وممكن الهيدر يتزور.\n\nاحمِ البحث والتجديد والاستهلاك تحت قفل واحد.",
			tryPrompt:     "توقّع عدد الثواني بعد تقريب الانتظار لفوق:",
			tryReveal:     "شوف الـduration تتحول لثواني كاملة.",
			watch:         "- تجاهل forwarded header لو الاتصال مش من proxy موثوق.\n- سجل آخر ظهور حتى للطلب المرفوض.",
			check:         "ابدأ بتحديد هوية العميل، ثم خزان مستقل، ثم قرار قبول أو رفض.\n\nفكّر إمتى يبقى حذف الخزان الخامل آمنًا.",
			checkQuestion: "أي عنوان يستخدم لتحديد العميل أولًا؟",
			turn:          "جرّب خزانين بعناوين توضيحية. غيّر واحدًا وتأكد إن التاني مستقل.",
			options:       []string{"1", "2", "1.5"},
			feedback:      []string{"التقريب لفوق ما بينزلش لثانية.", "ثانية ونصف تتقرب للثانية الكاملة التالية.", "Retry-After هنا عدد ثواني صحيح."},
			checkOptions:  []string{"RemoteAddr بعد تحليل IP", "X-Forwarded-For من أي طلب", "رقم المنفذ"},
			checkFeedback: []string{"الـpeer المباشر هو نقطة الثقة الأولى.", "أي عميل يقدر يكتب الهيدر بنفسه.", "المنفذ يتغير بين الاتصالات."},
			hints:         []string{"خلي لكل عميل معروف رصيدًا وحالة تجديد خاصة.", "اجمع تعديل الخريطة والخزان تحت mutex نفسه.", "ابدأ بتحليل ~SplitHostPort~ ثم ~ParseIP~."},
		},
		{
			mission:       "The gateway gives each client an independent permit tank.\n\nThe starter has configuration limits; add client keys, accounting, and cleanup.\n\nTime comes from a clock, not a real sleep.",
			idea:          "A tank map keeps each client's balance separate.\n\nBuild the key from a trusted IP; connection ports change, and headers can be forged.\n\nProtect lookup, refill, and spending under one lock.",
			tryPrompt:     "Predict the seconds after rounding the wait upward:",
			tryReveal:     "See the duration become whole seconds.",
			watch:         "- Ignore forwarded headers unless the peer is a trusted proxy.\n- Record last-seen time even for a refused request.",
			check:         "Start by identifying the client, then its tank, then acceptance or refusal.\n\nDecide when an idle tank can safely be forgotten.",
			checkQuestion: "Which address is the first source for client identity?",
			turn:          "Try two tanks with illustrative addresses. Change one and confirm the other stays separate.",
			options:       []string{"1", "2", "1.5"},
			feedback:      []string{"Rounding upward does not round down to one second.", "One and a half seconds rounds to the next whole second.", "Retry-After here is a whole-second count."},
			checkOptions:  []string{"RemoteAddr after parsing its IP", "X-Forwarded-For from any request", "The port number"},
			checkFeedback: []string{"The direct peer is the first trust boundary.", "Any client can write that header itself.", "Ports change between connections."},
			hints:         []string{"Give each known client its own balance and refill state.", "Guard map and tank changes with the same mutex.", "Start by parsing ~SplitHostPort~ and then ~ParseIP~."},
		},
	},
})

var ch27Lesson275 = ch27BuildLesson(ch27LessonSpec{
	prefix: "c275", exampleID: "e275-base62-alphabet",
	imports: []string{"fmt"},
	code: `func main() {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fmt.Printf("%c %c %c\n", alphabet[0], alphabet[10], alphabet[61])
}`,
	frames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الأبجدية مرتبة."},
			{Line: 3, Caption: "اقرأ أول ووسط وآخر خانة.", Output: "0 a Z"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The alphabet has an order."},
			{Line: 3, Caption: "Read first, middle, and last.", Output: "0 a Z"},
		},
	},
	tryExampleID: "e275-wrapped-collision",
	tryImports:   []string{"errors", "fmt"},
	tryCode: `func main() {
	errTaken := errors.New("code taken")
	wrapped := fmt.Errorf("save: %w", errTaken)
	fmt.Println(errors.Is(wrapped, errTaken))
}`,
	tryFrames: [2][]LessonFrame{
		{
			{Line: 0, Caption: "قبل التنفيذ."},
			{Line: 2, Caption: "الخطأ له هوية."},
			{Line: 3, Caption: "التغليف بيضيف سياق."},
			{Line: 4, Caption: "Is لسه لاقته.", Output: "true"},
		},
		{
			{Line: 0, Caption: "Before execution."},
			{Line: 2, Caption: "The error has an identity."},
			{Line: 3, Caption: "Wrapping adds context."},
			{Line: 4, Caption: "Is still finds it.", Output: "true"},
		},
	},
	tryAnswer: 1, checkAnswer: 0,
	copies: [2]ch27LessonCopy{
		{
			mission:       "روابط الأسطول تتحول لأكواد قصيرة وتعيد توجيه الطلب.\n\nالـstarter يعرّف الواجهة؛ هتضيف التحويل والمخزن والـhandler.\n\nافحص الهدف قبل طلب كود.",
			idea:          "~Base62~ أبجدية أرقام وحروف بترتيب ثابت.\n\nالترتيب نفسه يفضل ثابت في كل تحويل.\n\nالـstore هو صاحب قرار التكرار؛ الواجهة تخلي الـhandler يشتغل مع مخزن قابل للاستبدال.\n\nتحقق من ~URL~ قبل ما تسمح بتحويله.",
			tryPrompt:     "توقّع هل ~errors.Is~ تلاقي الخطأ بعد التغليف:",
			tryReveal:     "راقب هوية الخطأ خلال السلسلة.",
			watch:         "- رابط التحويل لازم يكون ~http~ أو ~https~ ومن غير user.\n- خليك بعيد عن محاولات التصادم اللانهائية.",
			check:         "قسّم المهمة لتحويل الأرقام، تخزين الأكواد، ثم routes.\n\nخطأ تصادم يختلف عن عطل مخزن، حتى لو الاتنين ملفوفين.",
			checkQuestion: "مين يعرف إن الكود اتستخدم بالفعل؟",
			turn:          "ابدأ بأبجدية toy وواجهة تخزين صغيرة. جرّب خطأ ملفوفًا قبل بناء التحويل.",
			options:       []string{"false", "true", "compile error"},
			feedback:      []string{"Is بتتبع الخطأ الأصلي جوه wrapper.", "الخطأ الملفوف لسه مرتبط بالأصل.", "~errors.Is~ و%w استخدام صحيح."},
			checkOptions:  []string{"المخزن", "مصدر الأرقام", "الـredirect response"},
			checkFeedback: []string{"المخزن شايف الأكواد الموجودة كلها.", "المصدر يعطي رقمًا ولا يعرف الحالة.", "الرد يحصل بعد قبول الكود."},
			hints:         []string{"ثبّت ترتيب الأبجدية قبل التفكير في تمثيل الرقم.", "خلّي واجهة التخزين مسؤولة عن حفظ وربط الكود.", "استخدم ~errors.Is~ للتعرف على خطأ تصادم ملفوف."},
		},
		{
			mission:       "Fleet links become short codes and redirect requests.\n\nThe starter defines the interface; add conversion, storage, and the handler.\n\nValidate a target before requesting a code.",
			idea:          "~Base62~ uses a fixed alphabet of digits and letters.\n\nKeep that order fixed for every conversion.\n\nThe store owns duplicate decisions; an interface lets the handler work with replaceable storage.\n\nValidate a ~URL~ before allowing a redirect.",
			tryPrompt:     "Predict whether ~errors.Is~ finds an error after wrapping:",
			tryReveal:     "Watch the error identity through the chain.",
			watch:         "- A redirect target must use ~http~ or ~https~ and have no user.\n- Avoid endless collision attempts.",
			check:         "Split the work into number conversion, code storage, and routes.\n\nA collision differs from a store failure, even when both are wrapped.",
			checkQuestion: "Who knows whether a code is already used?",
			turn:          "Start with a toy alphabet and small store interface. Try a wrapped error before building conversion.",
			options:       []string{"false", "true", "compile error"},
			feedback:      []string{"Is follows the original error inside the wrapper.", "The wrapped error still points to the original.", "~errors.Is~ and %w are valid here."},
			checkOptions:  []string{"The store", "The number source", "The redirect response"},
			checkFeedback: []string{"The store can see every code it holds.", "The source returns a number but knows no state.", "The response happens after a code is accepted."},
			hints:         []string{"Fix the alphabet order before thinking about number representation.", "Let the storage interface own saving and lookup.", "Use ~errors.Is~ to recognize a wrapped collision error."},
		},
	},
})
