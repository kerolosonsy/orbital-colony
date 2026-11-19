package levels

type ch25Copy struct {
	mission       string
	idea          string
	tryQuestion   string
	tryOptions    []string
	tryFeedback   []string
	reveal        string
	watch         string
	think         string
	checkQuestion string
	checkOptions  []string
	checkFeedback []string
	turn          string
}

func ch25Pair(arabic, english string) map[string]string {
	return map[string]string{"ar": arabic, "en": english}
}

func ch25Sandbox(prefix, language string) LessonBlock {
	kind, config := ch25SandboxConfig(prefix)
	return LessonBlock{ID: prefix + ".sandbox", Kind: "sandbox", Text: ch25SandboxText(prefix, language), Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch25SandboxText(prefix, language string) string {
	translations := map[string]map[string]string{
		"c251":  ch25Pair("شغّل الاتنين transactions؛ جرّب ترتيب تاني.", "Run both transactions; try another order."),
		"c252":  ch25Pair("بدّل الترتيب وشوف الكتابة الأخيرة.", "Change the order; watch the last write."),
		"c253":  ch25Pair("قدّم الساعة؛ احسب الفترات الكاملة.", "Advance the clock; count complete intervals."),
		"c254":  ch25Pair("انشر واسحب؛ لاحظ امتلاء الـbuffer.", "Publish and receive; watch the full buffer."),
		"c255":  ch25Pair("قدّم clock واقرا صلاحية الـvalue.", "Advance the clock; check entry expiry."),
		"c256":  ch25Pair("شغّل Add؛ بدّل ترتيب العمّال.", "Run Add; change worker order."),
		"c25b1": ch25Pair("اقفل الـchannel؛ اسحب الـvalues بعدها.", "Close the channel; receive queued values."),
	}
	return translations[prefix][language]
}

func ch25SandboxConfig(prefix string) (string, map[string]any) {
	switch prefix {
	case "c251":
		return "ch25-locked-account", ch25LockedAccountConfig()
	case "c252":
		return "ch25-race-report", ch25RaceReportConfig()
	case "c253":
		return "ch25-token-clock", ch25TokenClockConfig()
	case "c254":
		return "ch25-nonblocking-bus", ch25NonblockingBusConfig()
	case "c255":
		return "ch25-ttl-cache", ch25TTLCacheConfig()
	case "c256":
		return "ch25-atomic-counter", ch25AtomicCounterConfig()
	case "c25b1":
		return "ch25-closed-drain", ch25ClosedDrainConfig()
	}
	return "", nil
}

func ch25LockedAccountConfig() map[string]any {
	return map[string]any{
		"balance": 7,
		"transactions": []any{
			map[string]any{"worker": ch25Pair("العامل A", "Worker A"), "amount": 4},
			map[string]any{"worker": ch25Pair("العامل B", "Worker B"), "amount": 5},
		},
		"labels": map[string]any{
			"balance": ch25Pair("الرصيد", "Balance"), "next": ch25Pair("العملية الجاية", "Next operation"),
			"again": ch25Pair("ترتيب تاني", "Try other order"), "reset": ch25Pair("من الأول", "Reset"),
			"spends": ch25Pair("خصم", "spends"), "skips": ch25Pair("ماخصمش", "skips"),
			"because": ch25Pair("الرصيد مايكفيش", "insufficient balance"), "left": ch25Pair("المتبقي", "left"),
			"done": ch25Pair("خلصت الـtransactions", "Transactions complete"), "order": ch25Pair("الترتيب ممكن يتغير.", "Order can vary."),
			"protected": ch25Pair("كل فحص وتعديل جوه Mutex.", "Each check and update holds Mutex."),
			"history":   ch25Pair("العمليات", "Operations"),
		},
	}
}

func ch25RaceReportConfig() map[string]any {
	return map[string]any{
		"initial": 10,
		"workers": []any{
			map[string]any{"id": "a", "name": ch25Pair("العامل A", "Worker A"), "add": 3},
			map[string]any{"id": "b", "name": ch25Pair("العامل B", "Worker B"), "add": 5},
		},
		"schedules": [][]string{{"a-read", "b-read", "a-write", "b-write"}, {"b-read", "a-read", "b-write", "a-write"}},
		"labels": map[string]any{
			"value": ch25Pair("العداد", "Counter"), "next": ch25Pair("الخطوة الجاية", "Next step"),
			"again": ch25Pair("ترتيب تاني", "Try other order"), "reset": ch25Pair("من الأول", "Reset"),
			"read": ch25Pair("قرأ", "reads"), "write": ch25Pair("كتب", "writes"),
			"order": ch25Pair("ده ترتيب محتمل واحد.", "This is one possible order."),
			"trace": ch25Pair("الخطوات", "Trace"),
		},
	}
}

func ch25TokenClockConfig() map[string]any {
	return map[string]any{
		"elapsed": 11, "min": 0, "max": 19, "interval": 5, "tokens": 1, "capacity": 6,
		"labels": map[string]any{
			"elapsed": ch25Pair("الوقت اللي عدى", "Elapsed seconds"), "complete": ch25Pair("فترات كاملة", "Complete intervals"),
			"tokens": ch25Pair("التصاريح", "Permits"), "order": ch25Pair("الوقت الزايد ما يكملش فترة.", "Extra time does not complete an interval."),
		},
	}
}

func ch25NonblockingBusConfig() map[string]any {
	return map[string]any{
		"capacity": 1, "initial": []string{"ready"}, "message": "update",
		"labels": map[string]any{
			"queue": ch25Pair("الـbuffer", "Buffer"), "publish": ch25Pair("انشر", "Publish"),
			"receive": ch25Pair("اسحب", "Receive"), "reset": ch25Pair("من الأول", "Reset"),
			"sent": ch25Pair("اتبعت", "sent"), "full": ch25Pair("ممتلئ؛ اتخطّت", "full; skipped"),
			"empty": ch25Pair("فاضي", "empty"), "order": ch25Pair("الإرسال ما يستناش مساحة.", "A send never waits for space."),
			"attempts": ch25Pair("محاولات الإرسال", "Send attempts"),
		},
	}
}

func ch25TTLCacheConfig() map[string]any {
	return map[string]any{
		"ttl": 5, "age": 0, "maxAge": 8,
		"labels": map[string]any{
			"age": ch25Pair("عمر الـvalue بالثواني", "Entry age in seconds"), "advance": ch25Pair("قدّم الساعة", "Advance clock"),
			"read": ch25Pair("اقرا من الـcache", "Read cache"), "reset": ch25Pair("من الأول", "Reset"),
			"fresh": ch25Pair("لسه صالحة", "Fresh"), "expired": ch25Pair("انتهت", "Expired"),
			"ttl": ch25Pair("مدة الصلاحية", "TTL"),
		},
	}
}

func ch25AtomicCounterConfig() map[string]any {
	return map[string]any{
		"workers": []any{
			map[string]any{"name": ch25Pair("العامل A", "Worker A"), "add": 1},
			map[string]any{"name": ch25Pair("العامل B", "Worker B"), "add": 2},
			map[string]any{"name": ch25Pair("العامل C", "Worker C"), "add": 3},
		},
		"orders": [][]int{{0, 1, 2}, {2, 0, 1}},
		"labels": map[string]any{
			"value": ch25Pair("الـvalue بتاعة العداد", "Counter value"), "next": ch25Pair("شغّل atomic.Add", "Run atomic.Add"),
			"again": ch25Pair("بدّل الترتيب", "Shuffle order"), "reset": ch25Pair("من الأول", "Reset"),
			"load": ch25Pair("Load الحالي", "Current Load"), "order": ch25Pair("ترتيب Add يتغير؛ المجموع ثابت.", "Add order varies; total stays."),
			"trace": ch25Pair("خطوات Add", "Add steps"),
		},
	}
}

func ch25ClosedDrainConfig() map[string]any {
	return map[string]any{
		"values": []int{4, 6},
		"labels": map[string]any{
			"channel": ch25Pair("الـchannel بتاعة الشغل", "Job channel"), "open": ch25Pair("مفتوحة", "Open"),
			"closed": ch25Pair("مقفولة", "Closed"), "close": ch25Pair("اقفل الـchannel", "Close channel"),
			"receive": ch25Pair("استقبل الـvalue", "Receive value"), "reset": ch25Pair("من الأول", "Reset"),
			"drained": ch25Pair("الـqueue فاضي", "Queue drained"), "queued": ch25Pair("لسه فيه values", "Values remain queued"),
			"empty": ch25Pair("مفيش values؛ الاستقبال يستنى", "No values; receive waits"),
		},
	}
}

type ch25FrameInfo struct {
	line      int
	arabic    string
	english   string
	variables []LessonVariable
	output    string
}

type ch25Example struct {
	id      string
	code    string
	imports []string
	frames  []ch25FrameInfo
}

func ch25Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch25Question(id, text string, options, feedback []string, answer int) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch25Walkthrough(id, text string, example ch25Example, english bool) LessonBlock {
	frames := make([]LessonFrame, 0, len(example.frames))
	for _, frame := range example.frames {
		caption := frame.arabic
		if english {
			caption = frame.english
		}
		frames = append(frames, LessonFrame{Line: frame.line, Caption: caption, Variables: frame.variables, Output: frame.output})
	}
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: example.code, Imports: example.imports, ExampleID: example.id, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch25BuildLesson(prefix string, copy ch25Copy, examples [2]ch25Example, english bool) []LessonStep {
	titles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	language := "ar"
	if english {
		titles = []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
		language = "en"
	}
	tryCode := copy.tryQuestion + "\n\n~~~go\n" + examples[1].code + "\n~~~"
	return []LessonStep{
		{ID: prefix + ".mission-step", Title: titles[0], Blocks: []LessonBlock{ch25Markdown(prefix+".mission", copy.mission)}},
		{ID: prefix + ".idea-step", Title: titles[1], Blocks: []LessonBlock{ch25Walkthrough(prefix+".idea-example", copy.idea, examples[0], english)}},
		{ID: prefix + ".try-step", Title: titles[2], Blocks: []LessonBlock{
			ch25Question(prefix+".try-predict", tryCode, copy.tryOptions, copy.tryFeedback, 1),
			ch25Walkthrough(prefix+".try-reveal", copy.reveal, examples[1], english),
			ch25Sandbox(prefix, language),
		}},
		{ID: prefix + ".watch-step", Title: titles[3], Blocks: []LessonBlock{ch25Markdown(prefix+".watch", copy.watch)}},
		{ID: prefix + ".check-step", Title: titles[4], Blocks: []LessonBlock{ch25Markdown(prefix+".think", copy.think), ch25Question(prefix+".check-question", copy.checkQuestion, copy.checkOptions, copy.checkFeedback, 1)}},
		{ID: prefix + ".turn-step", Title: titles[5], Blocks: []LessonBlock{ch25Markdown(prefix+".turn", copy.turn)}},
	}
}

var ch25FuelExamples = [2]ch25Example{
	{
		id: "e251-locked-tank", imports: []string{"fmt", "sync"},
		code: `type Tank struct { mu sync.Mutex; fuel int }
func (t *Tank) Add(n int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fuel += n
}
func (t *Tank) Fuel() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.fuel
}
func main() {
	tank := Tank{fuel: 3}
	tank.Add(2)
	fmt.Println(tank.Fuel())
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 1, arabic: "الـlock جزء من Tank.", english: "The lock belongs to Tank."},
			{line: 13, arabic: "الرصيد يبدأ بـ3.", english: "The balance starts at 3."},
			{line: 14, arabic: "الإضافة محمية بالـlock.", english: "The add is protected by the lock."},
			{line: 15, arabic: "القراءة رجّعت 5.", english: "The read returns 5.", output: "5"},
		},
	},
	{
		id: "e251-mutex-scope", imports: []string{"fmt", "sync"},
		code: `func main() {
	var mu sync.Mutex
	fuel := 3
	mu.Lock()
	fuel += 2
	mu.Unlock()
	fmt.Println(fuel)
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الرصيد يبدأ بـ3.", english: "The balance starts at 3."},
			{line: 4, arabic: "بدأت المنطقة المحمية.", english: "The protected section begins."},
			{line: 5, arabic: "الرصيد بقى 5.", english: "The balance becomes 5."},
			{line: 7, arabic: "اطبع بعد فك الـlock.", english: "Print after unlocking.", output: "5"},
		},
	},
}

var ch25FuelArabic = ch25BuildLesson("c251", ch25Copy{
	mission:       "نفّذ ~FuelAccount~ عشان ~Load~ و~Burn~ المتوازيين يحافظوا على الرصيد والسعة. الـstarter فيه API ناقص؛ العقد يطلب تحديثات وقراءات آمنة.",
	idea:          "~Mutex~ زي الـlock على دفتر الرصيد. اقفل حوالين القراءة والتحديث كعملية واحدة، والقراءة المشتركة كمان محتاجة حماية.",
	tryQuestion:   "توقّع الرصيد بعد التحديث المحمي:",
	tryOptions:    []string{"3", "5", "2"},
	tryFeedback:   []string{"ده الرصيد قبل الإضافة.", "الإضافة رفعت الرصيد من 3 إلى 5.", "2 مقدار الإضافة، مش الرصيد النهائي."},
	reveal:        "تابع نطاق الـlock والـvalue بتاعة الرصيد بعد الجمع.",
	watch:         "- افحص الشرط والتحديث جوه نفس الـlock.\n- القارئ من غير الـlock يقدر يشوف كتابة متزامنة.",
	think:         "حدّد invariants الرصيد والسعة. راجع كل مسار نجاح ورفض تحت التزامن من غير تحديث جزئي.\n\nالـlock ما بيغيّرش الـvalues بنفسه.",
	checkQuestion: "إيه اللي يحميه الـmutex هنا؟",
	checkOptions:  []string{"function واحدة بس", "فحص الشرط مع التعديل", "سرعة التنفيذ"},
	checkFeedback: []string{"الـlock يحمي الحالة المشتركة بين النداءات.", "الاتنين لازم يبقوا جزء من المنطقة نفسها.", "الـlock مش أداة قياس أداء."},
	turn:          "ابدأ من أنواع الحقول وطرق الـAPI. تتبع الرصيد قبل وبعد كل عملية مع التوازي.",
}, ch25FuelExamples, false)

var ch25FuelEnglish = ch25BuildLesson("c251", ch25Copy{
	mission:       "Implement ~FuelAccount~ so concurrent ~Load~ and ~Burn~ preserve balance and capacity.\n\nThe starter has an incomplete API; the contract requires safe updates and reads.",
	idea:          "A ~Mutex~ is a lock on the balance ledger.\n\nProtect the check and update as one operation; shared reads need protection too.",
	tryQuestion:   "Predict the balance after the protected update:",
	tryOptions:    []string{"3", "5", "2"},
	tryFeedback:   []string{"That is the balance before the add.", "The add raises the balance from 3 to 5.", "2 is the add amount, not the final balance."},
	reveal:        "Follow the lock scope and the balance after addition.",
	watch:         "- Check the condition and update under the same lock.\n- An unlocked reader can see a concurrent write.",
	think:         "Identify the balance and capacity invariants. Review success and rejection paths under concurrency without partial updates.\n\nThe lock does not change values by itself.",
	checkQuestion: "What does the mutex protect here?",
	checkOptions:  []string{"Only one function", "The check together with the update", "Execution speed"},
	checkFeedback: []string{"The lock protects shared state across calls.", "Both belong in the same protected section.", "A mutex is not a performance meter."},
	turn:          "Start from the field types and API methods. Trace balance before and after each concurrent operation.",
}, ch25FuelExamples, true)

var ch25RaceExamples = [2]ch25Example{
	{
		id: "e252-lost-update", imports: []string{"fmt"},
		code: `func main() {
	base := 10
	first := base + 3
	second := base + 5
	base = first
	base = second
	fmt.Println(base)
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "القراءة الأصلية هي 10.", english: "The original read is 10."},
			{line: 4, arabic: "كل عامل حسب نسخته.", english: "Each worker used its snapshot."},
			{line: 6, arabic: "الكتابة التانية تمسح الأولى.", english: "The second write overwrites the first."},
			{line: 7, arabic: "ده ترتيب توضيحي واحد.", english: "This is one illustrative order.", output: "15"},
		},
	},
	{
		id: "e252-stale-snapshot", imports: []string{"fmt"},
		code: `func main() {
	base := 8
	first := base + 2
	second := base + 4
	base = first
	base = second
	fmt.Println(base)
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "النسخة تبدأ من 8.", english: "Both snapshots start at 8."},
			{line: 4, arabic: "الـvalue القديمة اتحسبت مرتين.", english: "Two stale values were computed."},
			{line: 6, arabic: "التانية تكتب 12.", english: "The second write stores 12."},
			{line: 7, arabic: "اطبع الـvalue الأخيرة.", english: "Print the last value.", output: "12"},
		},
	},
}

var ch25RaceArabic = ch25BuildLesson("c252", ch25Copy{
	mission:       "اكتب برهانًا حتميًا لفقد تحديث في ~Add~، وبعدها صلّح جمع تقارير المسح.\n\nالـstarter فيه موضع مزامنة محقون؛ المطلوب تفسير السبب وحماية الجمع.",
	idea:          "لو عاملين قرأوا نفس الرصيد قبل ما يكتبوا، آخر كتابة ممكن تمسح إضافة سابقة.\n\nالمثال يحاكي ترتيبًا ممكنًا من غير goroutines.",
	tryQuestion:   "في الترتيب المعروض، إيه الـvalue اللي هتفضل؟",
	tryOptions:    []string{"10", "12", "15"},
	tryFeedback:   []string{"دي القراءة الأصلية قبل الإضافات.", "الكتابة الأخيرة جاية من 10 + 5.", "15 محتاجة جمع الإضافتين معًا."},
	reveal:        "تابع القراءتين والكتابتيْن؛ دي محاكاة لترتيب واحد.",
	watch:         "- تزامن التنفيذ مش معناه جمع ذري.\n- الإقرار لازم ييجي بعد اكتمال التحديث.",
	think:         "فكّر إزاي تثبت ترتيبًا محددًا من غير توقيت. بعدين راجع أي قراءة وجمع لازم يفضلوا داخل نفس الحماية.",
	checkQuestion: "إيه اللي يمنع عاملًا من الكتابة فوق تحديث غيره؟",
	checkOptions:  []string{"الـlock حوالين القراءة والجمع", "تغيير ترتيب التقرير", "تأخير ثابت"},
	checkFeedback: []string{"المنطقة المحمية تضم العملية المركبة كلها.", "التقرير مش مصدر التحديث الضايع.", "التأخير ما يضمنش ترتيبًا."},
	turn:          "ابدأ بتحديد موضع القراءة والتحديث والإقرار. شغّل البرهان، ثم افحص السلوك المحمي تحت race detector.",
}, ch25RaceExamples, false)

var ch25RaceEnglish = ch25BuildLesson("c252", ch25Copy{
	mission:       "Write deterministic evidence of a lost update in ~Add~, then repair scan aggregation.\n\nThe starter has an injected synchronization point; explain the cause and protect the aggregate.",
	idea:          "If two workers read the same balance before writing, the last write can erase an earlier add.\n\nThis simulates one possible order without goroutines.",
	tryQuestion:   "In the shown order, which value remains?",
	tryOptions:    []string{"10", "12", "15"},
	tryFeedback:   []string{"That is the original read before either add.", "The last write comes from 10 + 5.", "15 would require combining both adds."},
	reveal:        "Follow both reads and writes; this simulates one possible order.",
	watch:         "- Concurrent execution does not make addition atomic.\n- Acknowledgement belongs after the update completes.",
	think:         "Think how to force a chosen order without timing.\n\nThen check that any read and add stay inside the same protection.",
	checkQuestion: "What prevents one worker from overwriting another update?",
	checkOptions:  []string{"A lock around read and add", "Changing report order", "A fixed delay"},
	checkFeedback: []string{"The protected section contains the whole compound operation.", "Reporting is not the source of the lost update.", "A delay does not guarantee ordering."},
	turn:          "Start by locating read, update, and acknowledgement. Run the evidence, then inspect protected behavior with the race detector.",
}, ch25RaceExamples, true)

var ch25BucketExamples = [2]ch25Example{
	{
		id: "e253-duration-steps", imports: []string{"fmt", "time"},
		code: `func main() {
	elapsed := 11 * time.Second
	interval := 5 * time.Second
	steps := elapsed / interval
	fmt.Println(int(steps))
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "مر 11 ثانية.", english: "Eleven seconds elapsed."},
			{line: 3, arabic: "الفترة 5 ثواني.", english: "The interval is five seconds."},
			{line: 4, arabic: "اكتملت فترتان.", english: "Two intervals are complete."},
			{line: 5, arabic: "اطبع عدد الفترات.", english: "Print the interval count.", output: "2"},
		},
	},
	{
		id: "e253-whole-intervals", imports: []string{"fmt", "time"},
		code: `func main() {
	elapsed := 16 * time.Second
	interval := 4 * time.Second
	fmt.Println(int(elapsed / interval))
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "مر 16 ثانية.", english: "Sixteen seconds elapsed."},
			{line: 3, arabic: "الفترة 4 ثواني.", english: "The interval is four seconds."},
			{line: 4, arabic: "اقسم لعدد كامل.", english: "Divide into whole intervals.", output: "4"},
		},
	},
}

var ch25BucketArabic = ch25BuildLesson("c253", ch25Copy{
	mission:       "نظّم التصاريح في ~TokenBucket~ مع ساعة محقونة. الـlab يراجع التجديد والحد الأقصى؛ الـstarter بيحدد الـinterface بتاعة الخزان والوقت.",
	idea:          "قسمة ~time.Duration~ على الفترة تعطي عدد الفترات الكاملة. الساعة المحقونة تخلي الحساب يتبع وقتًا تختاره بدل انتظار حقيقي.",
	tryQuestion:   "توقّع عدد الفترات الكاملة:",
	tryOptions:    []string{"3", "4", "16"},
	tryFeedback:   []string{"ثلاث فترات مجموعها 12 ثانية.", "16 فيها أربع فترات من 4 ثواني.", "دي الـvalue بتاعة الثواني، مش عدد الفترات."},
	reveal:        "تابع قسمة المدة على الفترة.",
	watch:         "- الكسور ما تكملش فترة جديدة.\n- تحديث العلامة الزمنية جزء من حالة الخزان.",
	think:         "حدد حالة الخزان والوقت المرجعي. راجع التجديد عند مرور وقت قصير، وعند وقت طويل، من غير افتراض نوم.\n\nالرصيد ما يزيدش عن السعة.\n\nالوقت الزايد بعد الامتلاء ما يتحولش لرصيد مخزون. سجّل الحالة قبل العملية وبعدها. قارن الرصيد بالسعة.",
	checkQuestion: "ليه نحقن clock بدل انتظار الوقت الحقيقي؟",
	checkOptions:  []string{"نختبر الزمن مباشرة", "نسرّع goroutine", "نغيّر نوع المدة"},
	checkFeedback: []string{"الـtest يختار اللحظة من غير انتظار.", "السرعة مش المقصودة هنا.", "المدة تفضل ~time.Duration~."},
	turn:          "ابدأ من الـinterface بتاعة الساعة وحالة bucket. تتبع قراءة الوقت والتحديث تحت نفس الملكية.",
}, ch25BucketExamples, false)

var ch25BucketEnglish = ch25BuildLesson("c253", ch25Copy{
	mission:       "Control permits in ~TokenBucket~ with an injected clock.\n\nThe lab checks refill and capacity; the starter defines the bucket and time interfaces.",
	idea:          "Dividing a ~time.Duration~ by an interval gives complete intervals.\n\nAn injected clock lets you choose time instead of waiting in real time.",
	tryQuestion:   "Predict the number of complete intervals:",
	tryOptions:    []string{"3", "4", "16"},
	tryFeedback:   []string{"Three intervals total twelve seconds.", "Sixteen contains four intervals of four seconds.", "That is the seconds value, not the interval count."},
	reveal:        "Follow the duration divided by the interval.",
	watch:         "- A fraction does not complete another interval.\n- Advancing the time marker is part of bucket state.",
	think:         "Identify the bucket state and reference time. Review refill after a short and a long elapsed period without real sleeps.\n\nBalance stays within capacity; extra time after refill does not become stored credit.\n\nRecord state before and after an operation. Compare balance with capacity.",
	checkQuestion: "Why inject a clock instead of waiting in real time?",
	checkOptions:  []string{"Tests choose time directly", "It speeds up a goroutine", "It changes the duration type"},
	checkFeedback: []string{"A test selects the instant without waiting.", "Speed is not the point here.", "The duration stays a ~time.Duration~."},
	turn:          "Start from the clock interface and bucket state. Trace time reads and updates under the same ownership.",
}, ch25BucketExamples, true)

var ch25BusExamples = [2]ch25Example{
	{
		id: "e254-nonblocking-send", imports: []string{"fmt"},
		code: `func main() {
	updates := make(chan string, 1)
	updates <- "ready"
	select {
	case updates <- "next":
		fmt.Println("sent")
	default:
		fmt.Println("full")
	}
	fmt.Println(<-updates)
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الـqueue خانة واحدة.", english: "The queue has one slot."},
			{line: 3, arabic: "ready شغلت الخانة.", english: "ready filled the slot."},
			{line: 8, arabic: "الـdefault اختار full.", english: "default selects full.", output: "full"},
			{line: 10, arabic: "الرسالة الأصلية لسه موجودة.", english: "The original message remains.", output: "full\nready"},
		},
	},
	{
		id: "e254-select-default", imports: []string{"fmt"},
		code: `func main() {
	queue := make(chan int, 1)
	queue <- 9
	select {
	case queue <- 6:
		fmt.Println("sent")
	default:
		fmt.Println("busy")
	}
	fmt.Println(<-queue)
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 3, arabic: "الخانة فيها 9.", english: "The slot contains 9."},
			{line: 8, arabic: "الإرسال ما استناش.", english: "The send does not wait.", output: "busy"},
			{line: 10, arabic: "اقرا الـvalue القديمة.", english: "Read the older value.", output: "busy\n9"},
		},
	},
}

var ch25BusArabic = ch25BuildLesson("c254", ch25Copy{
	mission:       "وزّع إشارات ~Publish~ على المشتركين من غير ما الـchannel البطيئة توقف الباقي.\n\nالـstarter فيه الـinterface بتاعة ناقل الأحداث؛ الـlab يراجع إلغاء الاشتراك والـlock.",
	idea:          "في ~select~، فرع ~default~ يشتغل فورًا لو مفيش إرسال جاهز. كده محاولة الإرسال ما تستناش مستهلكًا بطيئًا.",
	tryQuestion:   "الـchannel ممتلئة؛ توقّع الفرع والناتج:",
	tryOptions:    []string{"sent ثم 9", "busy ثم 9", "تعليق"},
	tryFeedback:   []string{"مافيش خانة متاحة للإرسال الجديد.", "الـdefault يشتغل والرسالة القديمة تفضل.", "وجود default يمنع انتظار الإرسال."},
	reveal:        "تابع فرع ~default~ في ~select~ ثم اقرا الرسالة المنتظرة.",
	watch:         "- الإرسال غير الحاجب ممكن يتخطى مشتركًا ممتلئًا.\n- إغلاق الـchannel أثناء الإرسال يحتاج تنسيقًا.",
	think:         "افصل موضوع الرسالة عن المشتركين المطابقين. راجع منافسة الإرسال مع إلغاء الاشتراك والإغلاق.\n\nعداد النشر يحسب الإرسالات المقبولة بس؛ المشترك غير المطابق ما يدخلش المحاولة.",
	checkQuestion: "إيه أثر ~default~ في إرسال الـchannel؟",
	checkOptions:  []string{"يستنى خانة", "يرجّع لو الإرسال مش جاهز", "يفرغ الـchannel"},
	checkFeedback: []string{"الانتظار يحصل من غير default.", "select ما يوقفش عند إرسال غير جاهز.", "default ما يقرأش الرسائل."},
	turn:          "ابدأ من عمر subscription وملكية الـchannel. تتبع النشر والإلغاء والإغلاق تحت حالة واحدة.",
}, ch25BusExamples, false)

var ch25BusEnglish = ch25BuildLesson("c254", ch25Copy{
	mission:       "Broadcast ~Publish~ signals without letting a slow subscriber block the rest.\n\nThe starter has an event bus interface; the lab checks unsubscribe and locking.",
	idea:          "In a ~select~, ~default~ runs immediately when no send is ready.\n\nThis keeps a send attempt from waiting on a slow consumer.",
	tryQuestion:   "The channel is full; predict the branch and output:",
	tryOptions:    []string{"sent then 9", "busy then 9", "blocked"},
	tryFeedback:   []string{"No slot is available for the new send.", "default runs and the old message remains.", "A default branch prevents waiting on the send."},
	reveal:        "Follow the ~default~ choice, then receive the waiting message.",
	watch:         "- A nonblocking send can skip a full subscriber.\n- Closing a channel during sends needs coordination.",
	think:         "Separate a signal topic from matching subscribers. Review publish racing with unsubscribe and close.\n\nThe publish count includes accepted sends only; unmatched subscribers are not attempted.",
	checkQuestion: "What does ~default~ do for a channel send?",
	checkOptions:  []string{"Waits for a slot", "Runs when the send is not ready", "Drains the channel"},
	checkFeedback: []string{"Waiting happens without default.", "The selection does not block on an unready send.", "default does not read messages."},
	turn:          "Start from subscription lifetime and channel ownership. Trace publish, unsubscribe, and close under one state policy.",
}, ch25BusExamples, true)

var ch25CacheExamples = [2]ch25Example{
	{
		id: "e255-rwmutex-read", imports: []string{"fmt", "sync"},
		code: `type Cache struct {
	mu sync.RWMutex
	values map[string]string
}
func (c *Cache) Lookup(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[key]
}
func main() {
	cache := Cache{values: map[string]string{"dock": "ready"}}
	fmt.Println(cache.Lookup("dock"))
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "القراءة محمية بـRWMutex.", english: "The read uses RWMutex."},
			{line: 6, arabic: "RLock يحمي lookup.", english: "RLock protects lookup."},
			{line: 12, arabic: "الـvalue الموجودة هي ready.", english: "The stored value is ready."},
			{line: 13, arabic: "اطبع الـvalue.", english: "Print the value.", output: "ready"},
		},
	},
	{
		id: "e255-stop-join", imports: []string{"fmt"},
		code: `func main() {
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-stop
		fmt.Println("cleanup stopped")
	}()
	close(stop)
	<-done
	fmt.Println("joined")
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "إشارة الإيقاف تبدأ.", english: "Create the stop signal."},
			{line: 9, arabic: "إشارة الإيقاف اتقفلت.", english: "The stop signal is closed."},
			{line: 10, arabic: "انتظر نهاية العامل.", english: "Wait for the worker to end."},
			{line: 11, arabic: "ده ترتيب ممكن.", english: "This is one possible order.", output: "cleanup stopped\njoined"},
		},
	},
}

var ch25CacheArabic = ch25BuildLesson("c255", ch25Copy{
	mission:       "نفّذ ~TTLCache~ لقراءات صالحة، وعامل تنظيف قابل للإيقاف من غير تسريب. الـstarter فيه الساعة المحقونة ونوع الـcache؛ الـlab يراجع عمر العامل.",
	idea:          "~RWMutex~ يسمح بقراءات متزامنة، ويمنعها وقت الكتابة. صلاحية الـvalue تعتمد على الساعة، مش وجود الـkey وحده.",
	tryQuestion:   "إيه اللي يضمن انتظار إغلاق العامل قبل الـreturn؟",
	tryOptions:    []string{"إغلاق stop بس", "استقبال done بعد stop", "إنشاء ticker"},
	tryFeedback:   []string{"الإشارة تطلب الإيقاف لكنها مش إقرار النهاية.", "استقبال done يثبت إن العامل خلص.", "ticker مش إشارة انضمام."},
	reveal:        "تابع request الإيقاف ثم إشارة انتهاء العامل. ده ترتيب ممكن للتنفيذ.",
	watch:         "- ~RLock~ للقراءة، و~Lock~ للتعديل.\n- كل عامل يحتاج مسار إيقاف وانتظار واضح.",
	think:         "راجع حد الصلاحية في lookup، ثم مسح المنتهي، ثم إغلاق متكرر آمن. ما تخلطش قراءة مع حذف.",
	checkQuestion: "إيه اللي يميّز ~RLock~؟",
	checkOptions:  []string{"يسمح بقراءات متزامنة", "يقفل الكتابة نهائيًا", "يمسح الـmap"},
	checkFeedback: []string{"قراءات متعددة ممكن تشارك الـlock.", "الكتابة تستخدم الـlock الحصري.", "الـlock ما يعدلش البيانات."},
	turn:          "ابدأ بتحديد ملكية map والعامل. تتبع كل ~Close~ من الـrequest لحد نهاية التنظيف.",
}, ch25CacheExamples, false)

var ch25CacheEnglish = ch25BuildLesson("c255", ch25Copy{
	mission:       "Implement ~TTLCache~ for live readings with a stoppable cleanup worker and no leaks.\n\nThe starter has the injected clock and cache type; the lab checks worker lifetime.",
	idea:          "~RWMutex~ allows concurrent reads and excludes them during writes. A value's validity depends on the clock, not just key presence.",
	tryQuestion:   "What ensures the worker has ended before return?",
	tryOptions:    []string{"Closing stop only", "Receive done after stop", "Creating a ticker"},
	tryFeedback:   []string{"The signal requests shutdown but does not confirm the end.", "Receiving done confirms the worker finished.", "A ticker is not a join signal."},
	reveal:        "Follow the stop request, then the worker completion signal. This is one possible execution order.",
	watch:         "- Use ~RLock~ to read and ~Lock~ to modify.\n- Every worker needs a clear stop and wait path.",
	think:         "Review the expiry boundary in lookup, then removing stale entries, then safe repeated close. Keep reads separate from deletion.",
	checkQuestion: "What is distinct about ~RLock~?",
	checkOptions:  []string{"It allows concurrent reads", "It blocks writes forever", "It clears a map"},
	checkFeedback: []string{"Multiple readers can share the lock.", "Writes use the exclusive lock.", "A lock does not change data."},
	turn:          "Start by identifying map and worker ownership. Trace every ~Close~ from request through cleanup completion.",
}, ch25CacheExamples, true)

var ch25AtomicExamples = [2]ch25Example{
	{
		id: "e256-atomic-count", imports: []string{"fmt", "sync/atomic"},
		code: `func main() {
	var pulses atomic.Int64
	pulses.Add(2)
	pulses.Add(3)
	fmt.Println(pulses.Load())
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "العداد يبدأ بصفر.", english: "The counter starts at zero."},
			{line: 3, arabic: "أضف نبضتين.", english: "Add two pulses."},
			{line: 4, arabic: "أضف تلاتة كمان.", english: "Add three more."},
			{line: 5, arabic: "اقرا الإجمالي 5.", english: "Read the total, 5.", output: "5"},
		},
	},
	{
		id: "e256-atomic-load", imports: []string{"fmt", "sync/atomic"},
		code: `func main() {
	var count atomic.Int64
	count.Add(4)
	count.Add(1)
	fmt.Println(count.Load())
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "ابدأ بعداد ذري.", english: "Start an atomic counter."},
			{line: 3, arabic: "الـvalue بقت 4.", english: "The value becomes 4."},
			{line: 4, arabic: "زود نبضة واحدة.", english: "Add one pulse."},
			{line: 5, arabic: "اقرا الإجمالي.", english: "Read the total.", output: "5"},
		},
	},
}

var ch25AtomicArabic = ch25BuildLesson("c256", ch25Copy{
	mission:       "قارن عداد ~atomic~ وعداد ~mutex~ بنفس الحمل في ~Bench~. الـstarter فيه التنفيذ والـbenchmark ناقص؛ المطلوب عقد صحيح وقياس عادل.",
	idea:          "~atomic.Int64~ يوفر عمليات ذرية للعداد. استخدم ~Add~ للتحديث و~Load~ للقراءة؛ المقارنة لازم تعمل نفس الشغل.",
	tryQuestion:   "توقّع الـvalue بتاعة العداد بعد العمليتين:",
	tryOptions:    []string{"4", "5", "1"},
	tryFeedback:   []string{"4 هي الـvalue قبل الزيادة الأخيرة.", "الإجمالي بعد الزيادة هو 5.", "1 مقدار الزيادة الأخيرة."},
	reveal:        "تابع عمليتي ~Add~ لحد قراءة ~Load~.",
	watch:         "- عدّل المتغير من خلال atomic API بس.\n- Benchmark يقارن نفس عدد العمليات في النسختين.",
	think:         "اكتب العقد قبل القياس: معنى ~Add~ و~Total~ واحد. بعدين قارن صحة النتيجة والحمل نفسه.\n\natomic مناسب لرقم مستقل. لو invariant بيربط حقول، mutex أوضح؛ والوقت وحده مش دليل صحة.\n\nكرر الحمل قبل مقارنة النتيجة.\n\nتحقق من الإجمالي؛ الصحة تسبق السرعة.",
	checkQuestion: "إيه اللي يخلي benchmark مقارنة عادلة؟",
	checkOptions:  []string{"نفس الشغل على النسختين", "عدادات تبدأ بـvalues مختلفة", "طباعة جوه الـloop"},
	checkFeedback: []string{"غيّر آلية الحماية بس.", "البداية المختلفة تفسد المقارنة.", "الطباعة تغيّر الحمل المقاس."},
	turn:          "ابدأ بعقد العدادين وبنية الـbenchmark. راجع الإجمالي بعد انتهاء المشاركين.",
}, ch25AtomicExamples, false)

var ch25AtomicEnglish = ch25BuildLesson("c256", ch25Copy{
	mission:       "Compare an ~atomic~ counter and a ~mutex~ counter under the same workload in ~Bench~.\n\nThe starter has implementations and an incomplete benchmark; define a correct contract and fair measurement.",
	idea:          "~atomic.Int64~ provides atomic counter operations. Use ~Add~ to update and ~Load~ to read; both versions must do the same work.",
	tryQuestion:   "Predict the counter after both operations:",
	tryOptions:    []string{"4", "5", "1"},
	tryFeedback:   []string{"4 is the value before the last increment.", "The total after the add is 5.", "1 is the last increment amount."},
	reveal:        "Follow both ~Add~ calls through the ~Load~ read.",
	watch:         "- Update through the atomic API only.\n- A benchmark compares the same operation count in both versions.",
	think:         "Write the counter contract before measuring: ~Add~ and ~Total~ mean the same. Then compare correctness and equal work.\n\nAtomic fits an independent number. Use a mutex for a multi-field invariant; timing alone does not prove correctness.\n\nRepeat the workload before comparing results.\n\nVerify the total; correctness comes before speed.",
	checkQuestion: "What makes a benchmark comparison fair?",
	checkOptions:  []string{"The same work in both versions", "Counters starting at different values", "Printing inside the loop"},
	checkFeedback: []string{"Change only the synchronization mechanism.", "Different starting values distort the comparison.", "Printing changes the measured workload."},
	turn:          "Start from both counter contracts and the benchmark shape. Check totals after participants finish.",
}, ch25AtomicExamples, true)

var ch25BossExamples = [2]ch25Example{
	{
		id: "e25b1-close-and-drain", imports: []string{"fmt"},
		code: `func main() {
	jobs := make(chan int, 2)
	jobs <- 4
	jobs <- 6
	close(jobs)
	for ore := range jobs {
		fmt.Println(ore)
	}
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الـqueue جاهز.", english: "The queue is ready."},
			{line: 5, arabic: "المرسل أنهى الإرسال.", english: "The sender finished sending."},
			{line: 7, arabic: "اقرا أول خامة.", english: "Read the first ore value.", output: "4"},
			{line: 7, arabic: "range تفضي الباقي.", english: "range drains the rest.", output: "4\n6"},
		},
	},
	{
		id: "e25b1-channel-end", imports: []string{"fmt"},
		code: `func main() {
	jobs := make(chan int, 1)
	jobs <- 8
	close(jobs)
	value, ok := <-jobs
	fmt.Println(value, ok)
}`,
		frames: []ch25FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 4, arabic: "الـchannel اتقفلت بعد الإرسال.", english: "The channel closed after sending."},
			{line: 5, arabic: "الـvalue لسه متاحة.", english: "The value is still available."},
			{line: 6, arabic: "وصلت الـvalue فعلية.", english: "A real value arrived.", output: "8 true"},
		},
	},
}

var ch25BossArabic = ch25BuildLesson("c25b1", ch25Copy{
	mission:       "نفّذ ~Refine~ بعامل محدود، نتيجة بترتيب الخام، وإيقاف منضم عند الـerror أو الإلغاء.\n\nالـboss يجمع أنماط channels وcontext ونتائج من دروس سابقة.",
	idea:          "المرسل اللي بيقفل الـchannel بيعلن انتهاء المصدر. المستهلك يقدر يفضي الرسائل الموجودة؛ الإغلاق مش بديل عن انتظار العمال.",
	tryQuestion:   "بعد إغلاق الـchannel، هل الـvalue اللي اتبعتت لسه بتتقري؟",
	tryOptions:    []string{"لأ، الإغلاق بيمسحها", "آه، لحد ما تتسحب", "الاستقبال يعمل panic"},
	tryFeedback:   []string{"الإغلاق ما بيمسحش الرسائل الموجودة.", "الـvalues بتتستهلك قبل علامة النهاية.", "panic يحصل عند الإرسال بعد الإغلاق."},
	reveal:        "تابع قراءة الرسالة ثم علامة نهاية الـchannel.",
	watch:         "- إغلاق الـchannel مسؤولية المرسل المالك.\n- النتيجة المرتبة مش معناها إن العمال خلصوا بالترتيب.",
	think:         "راجع العقد على مراحل: اتأكد من المدخلات، وزّع العمل، اجمع النتائج، ثم انضم قبل الـreturn. حدد مصدر أول error.",
	checkQuestion: "إيه اللي لازم يحصل قبل ما الـfunction اللي فيها workers تعمل return؟",
	checkOptions:  []string{"إغلاق كل المتغيرات", "انتظار انتهاء العمال", "رجّع أول نتيجة"},
	checkFeedback: []string{"مش كل متغير بيتقفل.", "الانتظار يمنع عمالًا شاردين بعد الـreturn.", "نتيجة واحدة ما تكملش العقد."},
	turn:          "ابدأ بتقسيم ملكية الـchannels والنتائج والـcontext. استخدم الـtests كعقد، وسيب ~main~ فاضية.",
}, ch25BossExamples, false)

var ch25BossEnglish = ch25BuildLesson("c25b1", ch25Copy{
	mission:       "Implement ~Refine~ with a bounded crew, input-ordered results, and joined shutdown on failure or cancellation.\n\nThe boss combines channel, context, and result patterns from earlier lessons.",
	idea:          "The job-channel sender declares that no more work is coming.\n\nA consumer can drain queued messages; closure does not replace joining workers.",
	tryQuestion:   "After the channel closes, can the sent value still be received?",
	tryOptions:    []string{"No, close deletes it", "Yes, until it is consumed", "The receive panics"},
	tryFeedback:   []string{"Closing does not erase queued values.", "Values are consumed before the end signal.", "Sending after close causes a panic."},
	reveal:        "Follow the message receive, then the channel end signal.",
	watch:         "- The sending owner closes the channel.\n- Ordered results do not imply ordered worker completion.",
	think:         "Review the contract in phases: validate inputs, distribute work, collect results, then join before return. Identify the first error source.",
	checkQuestion: "What must happen before a function with workers returns?",
	checkOptions:  []string{"Close every variable", "Wait for workers to finish", "Return the first result"},
	checkFeedback: []string{"Not every variable is closed.", "Joining prevents workers from outliving the call.", "One result does not complete the contract."},
	turn:          "Start by assigning channel, result, and context ownership. Use the lab checks as a contract and leave ~main~ empty.",
}, ch25BossExamples, true)
