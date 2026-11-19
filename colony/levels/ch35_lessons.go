package levels

import "strings"

type ch35FrameSpec struct {
	line    int
	caption [2]string
	vars    []LessonVariable
	output  string
}

type ch35ExampleSpec struct {
	id      string
	code    string
	imports []string
	frames  []ch35FrameSpec
}

type ch35LessonSpec struct {
	prefix        string
	mission       [2]string
	ideaText      [2]string
	idea          ch35ExampleSpec
	question      [2]string
	options       [2][]string
	answer        int
	feedback      [2][]string
	revealText    [2]string
	reveal        ch35ExampleSpec
	watch         [2]string
	plan          [2]string
	checkQuestion [2]string
	checkOptions  [2][]string
	checkAnswer   int
	checkFeedback [2][]string
	turn          [2]string
	hints         [2][3]string
}

func ch35Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: ch35ShortParagraphs(text)}
}

func ch35Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: ch35ShortParagraphs(text), Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch35Walkthrough(id, text string, example ch35ExampleSpec, language int) LessonBlock {
	frames := make([]LessonFrame, 0, len(example.frames))
	for _, frame := range example.frames {
		frames = append(frames, LessonFrame{Line: frame.line, Caption: frame.caption[language], Variables: frame.vars, Output: frame.output})
	}
	return LessonBlock{
		ID: id, Kind: "walkthrough", Text: ch35ShortParagraphs(text), Code: example.code,
		Imports: example.imports, ExampleID: example.id,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch35FencedQuestion(text, code string) string {
	return ch35ShortParagraphs(text) + "\n\n~~~go\n" + code + "\n~~~"
}

func ch35ShortParagraphs(text string) string {
	lines := strings.Split(text, "\n")
	paragraphs := make([]string, 0, len(lines))
	inFence := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
			inFence = !inFence
			paragraphs = append(paragraphs, line)
			continue
		}
		if !inFence {
			line = strings.ReplaceAll(line, ". ", ".\n\n")
			line = strings.ReplaceAll(line, "? ", "?\n\n")
			line = strings.ReplaceAll(line, "؟ ", "؟\n\n")
			line = strings.ReplaceAll(line, "! ", "!\n\n")
		}
		paragraphs = append(paragraphs, line)
	}
	return strings.Join(paragraphs, "\n")
}

func ch35BuildLesson(spec ch35LessonSpec) ([]LessonStep, []LessonStep) {
	titles := [2][]string{
		{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
		{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
	}
	steps := [2][]LessonStep{}
	for language := range steps {
		prefix := spec.prefix
		steps[language] = []LessonStep{
			{ID: prefix + "mission-step", Title: titles[language][0], Blocks: []LessonBlock{ch35Markdown(prefix+"mission", spec.mission[language])}},
			{ID: prefix + "idea-step", Title: titles[language][1], Blocks: []LessonBlock{ch35Walkthrough(prefix+"idea-trace", spec.ideaText[language], spec.idea, language)}},
			{ID: prefix + "try-step", Title: titles[language][2], Blocks: []LessonBlock{
				ch35Question(prefix+"predict", ch35FencedQuestion(spec.question[language], spec.reveal.code), spec.options[language], spec.answer, spec.feedback[language]),
				ch35Walkthrough(prefix+"reveal", spec.revealText[language], spec.reveal, language),
			}},
			{ID: prefix + "watch-step", Title: titles[language][3], Blocks: []LessonBlock{ch35Markdown(prefix+"watch", spec.watch[language])}},
			{ID: prefix + "check-step", Title: titles[language][4], Blocks: []LessonBlock{
				ch35Markdown(prefix+"plan", spec.plan[language]),
				ch35Question(prefix+"check-question", spec.checkQuestion[language], spec.checkOptions[language], spec.checkAnswer, spec.checkFeedback[language]),
			}},
			{ID: prefix + "turn-step", Title: titles[language][5], Blocks: []LessonBlock{ch35Markdown(prefix+"turn", spec.turn[language])}},
		}
	}
	return steps[0], steps[1]
}

var ch35QueueArabic, ch35QueueEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c351.",
	mission: [2]string{
		"ابني طابور Redis يوزّع مهام الورش على عمّال مستقلين. فكّر في المهمة وهي جاهزة، مستأجرة، مؤكدة، أو محتاجة محاولة تانية.",
		"Build a Redis queue that gives workshop tasks to independent workers. Track a task as ready, leased, acknowledged, or waiting for another attempt.",
	},
	ideaText: [2]string{
		"الـlease له نهاية محسوبة من ساعة يقدر الاختبار يتحكم فيها. عند نفس لحظة النهاية، ~Before~ مش بتعتبر الوقت لسه سابقها.",
		"A lease has a deadline based on a clock the caller can control. At the deadline itself, ~Before~ no longer treats the current time as earlier.",
	},
	idea: ch35ExampleSpec{
		id: "e351-lease-boundary", imports: []string{"fmt", "time"},
		code: `func main() {
	start := time.Unix(1000, 0)
	deadline := start.Add(3 * time.Second)
	now := start.Add(3 * time.Second)
	fmt.Println(now.Before(deadline))
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"بدأت الساعة المثال.", "The sample clock starts."}, vars: []LessonVariable{{Name: "start", Value: "1970-01-01 00:16:40 +0000 UTC"}}},
			{line: 3, caption: [2]string{"حسبنا deadline ثابتة.", "A fixed deadline is set."}, vars: []LessonVariable{{Name: "deadline", Value: "1970-01-01 00:16:43 +0000 UTC"}}},
			{line: 4, caption: [2]string{"الساعة وصلت للنهاية.", "The clock reaches the deadline."}, vars: []LessonVariable{{Name: "now", Value: "1970-01-01 00:16:43 +0000 UTC"}}},
			{line: 5, caption: [2]string{"الوقت مش قبل النهاية.", "Now is not before the deadline."}, vars: []LessonVariable{{Name: "now", Value: "1970-01-01 00:16:43 +0000 UTC"}}, output: "false"},
		},
	},
	question: [2]string{
		"المثال بيسأل الخريطة هل ~pump~ اتشاف قبل كده. توقّع الناتج.",
		"The map asks whether ~pump~ was seen before. Predict the output.",
	},
	options: [2][]string{{"true", "false", "مفتاح مش موجود"}, {"true", "false", "A missing key"}}, answer: 0,
	feedback: [2][]string{
		{"القيمة true موجودة للمفتاح.", "الخريطة فيها المفتاح بالفعل.", "comma-ok بيرجع القيمة وحالة الوجود."},
		{"The key has a stored true value.", "The map already contains this key.", "Comma-ok reports the value and whether it exists."},
	},
	revealText: [2]string{"اتبع قراءة الـmap وشوف قيمة الوجود.", "Follow the map lookup and watch its presence result."},
	reveal: ch35ExampleSpec{
		id: "e351-map-presence", imports: []string{"fmt"},
		code: `func main() {
	seen := map[string]bool{"pump": true}
	_, ok := seen["pump"]
	fmt.Println(ok)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"المفتاح اتسجل قبل كده.", "The key is already recorded."}, vars: []LessonVariable{{Name: "seen", Value: "map[pump:true]"}}},
			{line: 3, caption: [2]string{"القراءة رجعت حالة الوجود.", "The lookup returns presence."}, vars: []LessonVariable{{Name: "ok", Value: "true"}}},
			{line: 4, caption: [2]string{"النتيجة true.", "The result is true."}, vars: []LessonVariable{{Name: "ok", Value: "true"}}, output: "true"},
		},
	},
	watch: [2]string{
		"العامل ممكن يستلم المهمة تاني لو وقع قبل ack؛ خلي تكرار التسليم متوقع. وما تعتمدش على ساعة محلية لحسم سباق بين عمّال.",
		"A worker may receive a task again if it crashes before ack; expect repeated delivery. Do not use a local clock to settle a race between workers.",
	},
	plan: [2]string{
		"ارسم دورة المهمة من Enqueue لحد Claim، وبعدها اسأل إيه اللي يغيّره Ack أو Fail. راجع كمان نهاية الـlease ومسار المهام الميتة.",
		"Trace a task from Enqueue through Claim, then ask what Ack or Fail changes. Also review lease expiry and the dead-letter path.",
	},
	checkQuestion: [2]string{"إيه اللي يحسم مين استرجع lease منتهية بين عاملين؟", "What decides which of two workers reclaimed an expired lease?"},
	checkOptions: [2][]string{
		{"نتيجة ذرّية من Redis.", "مقارنة الساعة في كل worker.", "ترتيب بدء goroutines."},
		{"An atomic Redis result.", "Each worker's local clock check.", "The order goroutines start."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"Redis تضمن إن عملية واحدة بس تغيّر العضو.", "الساعات المحلية مش منسقة بين العمّال.", "ترتيب goroutines مش مضمون."},
		{"Redis lets only one operation change the member.", "Workers' local clocks are not coordinated.", "Goroutine start order is not guaranteed."},
	},
	turn: [2]string{
		"اقرأ أنواع ~Task~ و~Queue~ وتوقيعات methods. اتبع كل انتقال حالة لوحده، وسيب سلوك Redis قابل لإعادة التشغيل بالساعة المحقونة.",
		"Read the ~Task~ and ~Queue~ types and method signatures. Follow each state transition separately, and keep Redis behavior replayable with the injected clock.",
	},
	hints: [2][3]string{
		{"التسليم الموزّع ممكن يتكرر؛ صمّم المستهلك على الأساس ده.", "تابع انتقال المهمة بين الجاهز والمستأجر والمنتهي أو المعاد.", "ابدأ من ~Task~ وطرق ~Queue~ في الـstarter، وراجع ملكية كل lease."},
		{"Distributed delivery can repeat; design with that in mind.", "Track task movement between ready, leased, finished, or retried states.", "Start at ~Task~ and ~Queue~ methods in the starter, and review lease ownership."},
	},
})

var ch35SchedulerArabic, ch35SchedulerEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c352.",
	mission: [2]string{
		"اكتب مجدول صيانة يختار أقرب موعد ويقرر التعامل مع المواعيد الفايتة. كل تشغيل يبدأ من ~Tick~، من غير goroutines أو ساعة حقيقية.",
		"Build a maintenance scheduler that selects the next deadline and handles missed runs. Each run begins through ~Tick~, without goroutines or a real clock.",
	},
	ideaText: [2]string{
		"الـheap بيحافظ على أصغر عنصر عند الجذر. مش كل الشريحة مرتبة؛ تابع الرقم اللي تقدر توصله بسرعة.",
		"A heap keeps its smallest item at the root. The whole slice is not sorted; follow the value you can reach quickly.",
	},
	idea: ch35ExampleSpec{
		id: "e352-heap-root", imports: []string{"container/heap", "fmt"},
		code: `type minInts []int
func (h minInts) Len() int { return len(h) }
func (h minInts) Less(i, j int) bool { return h[i] < h[j] }
func (h minInts) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *minInts) Push(x any) { *h = append(*h, x.(int)) }
func (h *minInts) Pop() any { old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x }
func main() {
	h := &minInts{7, 3, 5}
	heap.Init(h)
	fmt.Println((*h)[0])
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 8, caption: [2]string{"القيم لسه slice عادية.", "The values start as a slice."}, vars: []LessonVariable{{Name: "h", Value: "[7 3 5]"}}},
			{line: 9, caption: [2]string{"Init جهزت بنية heap.", "Init establishes the heap."}, vars: []LessonVariable{{Name: "h[0]", Value: "3"}}},
			{line: 10, caption: [2]string{"الجذر هو الأصغر.", "The root is the smallest."}, vars: []LessonVariable{{Name: "h[0]", Value: "3"}}, output: "3"},
		},
	},
	question: [2]string{"في قيم مواعيد للتوضيح، إيه اللي بترجعه ~slices.Min~؟", "For these sample deadlines, what does ~slices.Min~ return?"},
	options:  [2][]string{{"2", "9", "الشريحة كلها"}, {"2", "9", "The whole slice"}}, answer: 0,
	feedback: [2][]string{
		{"Min بترجع أقل قيمة موجودة.", "9 مش أصغر عنصر في المثال.", "Min بترجع قيمة واحدة مش الشريحة."},
		{"Min returns the smallest value present.", "9 is not the smallest item here.", "Min returns one value, not the slice."},
	},
	revealText: [2]string{"اقرأ القيم وبعدين تابع أقل قيمة فيهم.", "Read the values, then follow the smallest one."},
	reveal: ch35ExampleSpec{
		id: "e352-heap-pop", imports: []string{"fmt", "slices"},
		code: `func main() {
	due := []int{9, 2, 6}
	fmt.Println(slices.Min(due))
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"دخلنا مواعيد مثال.", "The sample deadlines are loaded."}, vars: []LessonVariable{{Name: "due", Value: "[9 2 6]"}}},
			{line: 3, caption: [2]string{"اختارنا أقل قيمة.", "The smallest value is selected."}, vars: []LessonVariable{{Name: "due", Value: "[9 2 6]"}}, output: "2"},
		},
	},
	watch: [2]string{
		"~heap.Push~ و~heap.Pop~ بيحافظوا على بنية الأولوية. الاستدعاء المباشر لـPush بتاعت النوع مايرتبش الـheap.",
		"~heap.Push~ and ~heap.Pop~ preserve heap order. Calling the type's Push method directly does not order the heap.",
	},
	plan: [2]string{
		"افصل بيانات المهمة عن حالة تشغيلها. فكّر في التعادل بين المواعيد، وبعدها حدّد أثر تأخير Tick على الموعد القادم.",
		"Separate each job's data from its running state. Consider deadline ties, then decide how a late Tick affects the next slot.",
	},
	checkQuestion: [2]string{"لو موعد المهمة وصل وهي لسه running، إيه السلوك المناسب؟", "If a job is due while it is still running, what should happen?"},
	checkOptions: [2][]string{
		{"ما تبدأش تشغيل متداخل.", "ابدأ goroutine جديدة فورًا.", "امسح المهمة من الجدول."},
		{"Do not start an overlapping run.", "Start another goroutine immediately.", "Delete the job from the schedule."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"حالة running تمنع تداخل نفس المهمة.", "المجدول لا ينشئ goroutines.", "المهمة تفضل موجودة لمواعيد لاحقة."},
		{"The running state prevents overlap for that job.", "The scheduler does not create goroutines.", "The job remains for future deadlines."},
	},
	turn: [2]string{
		"راجع ~Job~ و~Run~ وسياسة المواعيد في الـstarter. جرّب التفكير بوقت افتراضي؛ ما تستخدمش ~Sleep~ أو ساعة النظام.",
		"Review ~Job~, ~Run~, and the misfire policy in the starter. Reason with supplied times; do not use ~Sleep~ or the system clock.",
	},
	hints: [2][3]string{
		{"الـheap بتسهّل متابعة الموعد الأقرب.", "افصل ترتيب المواعيد عن قرار تشغيل كل Job.", "ابدأ من ~Tick~ و~Done~، وراجع سياسة المواعيد في الـstarter."},
		{"A heap makes the next deadline easy to find.", "Separate deadline ordering from each job's run decision.", "Start at ~Tick~ and ~Done~, and review the policy in the starter."},
	},
})

var ch35LockArabic, ch35LockEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c353.",
	mission: [2]string{
		"اعمل قفل Redis له مالك ومدة عشان تمنع ورشتين من تشغيل نفس المفاعل معًا. وضّح كمان إيه اللي القفل ما يضمنوش لو مالكه اتأخر.",
		"Build a Redis lock with an owner and expiry so two workshops do not operate one reactor at once. Also explain what the lock cannot guarantee if its owner pauses.",
	},
	ideaText: [2]string{
		"قيمة القفل بتحدد مالكه. مقارنة token قديم بالمالك الحالي بتكشف إن الإيجار اتنقل؛ تابع نتيجة المقارنة.",
		"The lock value identifies its owner. Comparing an old token with the current owner reveals that the lease changed; follow the comparison.",
	},
	idea: ch35ExampleSpec{
		id: "e353-token-owner", imports: []string{"fmt"},
		code: `func main() {
	stored := "T-old"
	myToken := "T-new"
	fmt.Println(stored == myToken)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"Redis فيها مالك حالي.", "Redis holds a current owner."}, vars: []LessonVariable{{Name: "stored", Value: "T-old"}}},
			{line: 3, caption: [2]string{"الـtoken المحلي أقدم.", "The local token is stale."}, vars: []LessonVariable{{Name: "myToken", Value: "T-new"}}},
			{line: 4, caption: [2]string{"المالكين مش واحد.", "The owners do not match."}, vars: []LessonVariable{{Name: "stored == myToken", Value: "false"}}, output: "false"},
		},
	},
	question: [2]string{"هل الـtoken المحلي يطابق قيمة القفل؟", "Does the local token match the lock value?"},
	options:  [2][]string{{"true", "false", "مفيش قفل"}, {"true", "false", "There is no lock"}}, answer: 0,
	feedback: [2][]string{
		{"النصين متساويين حرفيًا.", "القيمة نفسها موجودة في الطرفين.", "المثال بيقارن نصين فقط."},
		{"The two strings are identical.", "The same value appears on both sides.", "The example compares two strings."},
	},
	revealText: [2]string{"قارن token الطلب بالمالك اللي Redis شايفاه.", "Compare the request token with the owner Redis currently sees."},
	reveal: ch35ExampleSpec{
		id: "e353-token-match", imports: []string{"fmt"},
		code: `func main() {
	stored := "T7"
	myToken := "T7"
	fmt.Println(stored == myToken)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"القفل يحمل T7.", "The lock holds T7."}, vars: []LessonVariable{{Name: "stored", Value: "T7"}}},
			{line: 3, caption: [2]string{"الطلب يحمل نفس token.", "The request has the same token."}, vars: []LessonVariable{{Name: "myToken", Value: "T7"}}},
			{line: 4, caption: [2]string{"المقارنة true.", "The comparison is true."}, vars: []LessonVariable{{Name: "stored == myToken", Value: "true"}}, output: "true"},
		},
	},
	watch: [2]string{
		"فحص الملكية ثم الحذف كخطوتين يفتح race. والـtoken يمنع مسح قفل غيرك، لكنه ما يمنعش كتابة مالك قديم على المورد.",
		"Checking ownership and deleting in separate steps creates a race. A token prevents deleting another lock, but not stale writes to the resource.",
	},
	plan: [2]string{
		"افصل التحقق عن القراءة من Redis. اتبع طلب الاكتساب من مدخلاته لنتيجته، وبعدها راجع إن التحرير يعتمد على المالك الحالي.",
		"Separate validation from Redis access. Trace acquisition from its inputs to its result, then check that release depends on the current owner.",
	},
	checkQuestion: [2]string{"إزاي تمنع قفل جديد يتأثر بتحرير قديم؟", "How do you keep an old release from affecting a new lock?"},
	checkOptions: [2][]string{
		{"مقارنة المالك والحذف ذريًا.", "Get وبعده Del منفصلين.", "حذف المفتاح من غير مقارنة."},
		{"Compare the owner and delete atomically.", "Run Get and Del separately.", "Delete the key without a comparison."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"العملية الواحدة تمنع أمر تاني من التداخل.", "المفتاح ممكن يتغير بين الطلبين.", "ده يمسح قفل مالك جديد."},
		{"One atomic operation prevents another command from slipping in.", "The key can change between those requests.", "That can delete a new owner's lock."},
	},
	turn: [2]string{
		"اقرأ ~Lease~ وواجهات Redis في الـstarter. جرّب تشرح حد الأمان بكلامك، وخلي تعليق ~fencing~ جزء من التسليم.",
		"Read ~Lease~ and the Redis-facing starter signatures. Explain the safety limit in your own words, and include the required ~fencing~ comment.",
	},
	hints: [2][3]string{
		{"المدة تحرر المفتاح، والـtoken يحدد المالك الحالي.", "افصل مسار الاكتساب عن التحرير، مع حماية التحديثات القديمة.", "ابدأ من ~Lease~ وتوقيعات ~Acquire~ و~Release~ في الـstarter."},
		{"Expiry frees the key; the token identifies its current owner.", "Separate acquisition from release while guarding against stale updates.", "Start at ~Lease~ and the ~Acquire~ and ~Release~ signatures in the starter."},
	},
})

var ch35OutboxArabic, ch35OutboxEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c354.",
	mission: [2]string{
		"سجّل طلب قطع الغيار وحدثه معًا، وبعدها انشر الحدث وخلّي المستهلك يتحمل تكراره. كل جزء له مسؤولية مختلفة في ضمان عدم فقد الحالة.",
		"Record a spare-parts order and its event together, then publish and let the consumer handle duplicates. Each part has a distinct role in preserving state.",
	},
	ideaText: [2]string{
		"JSON اللي مافيهوش field مش بيمسح قيمة بدأت بيها. المثال يبدأ ~Qty~ افتراضيًا، وبعدين يفك payload قديم.",
		"JSON that omits a field does not erase a value you set first. This example starts with a default ~Qty~, then decodes an older payload.",
	},
	idea: ch35ExampleSpec{
		id: "e354-json-default", imports: []string{"encoding/json", "fmt"},
		code: `type Payload struct {
	Qty int
}
func main() {
	p := Payload{Qty: 1}
	_ = json.Unmarshal([]byte("{\"Item\":\"valve\"}"), &p)
	fmt.Println(p.Qty)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 1, caption: [2]string{"Payload فيها كمية.", "Payload has a quantity field."}},
			{line: 5, caption: [2]string{"بدأنا بكمية افتراضية.", "A default quantity is set."}, vars: []LessonVariable{{Name: "p.Qty", Value: "1"}}},
			{line: 6, caption: [2]string{"الـJSON مافيهوش Qty.", "The JSON has no Qty."}, vars: []LessonVariable{{Name: "p.Qty", Value: "1"}}},
			{line: 7, caption: [2]string{"القيمة الافتراضية فضلت.", "The default value remains."}, vars: []LessonVariable{{Name: "p.Qty", Value: "1"}}, output: "1"},
		},
	},
	question: [2]string{"المستهلك شاف الحدث ~4~ قبل كده. إيه قيمة ~ok~؟", "The consumer has seen event ~4~ before. What is ~ok~?"},
	options:  [2][]string{{"true", "false", "خطأ من map"}, {"true", "false", "A map error"}}, answer: 0,
	feedback: [2][]string{
		{"الخريطة فيها القيمة true للمفتاح.", "القيمة المسجلة هي true.", "القراءة بالشكل comma-ok ما بترجعش error."},
		{"The map stores true for this key.", "The recorded value is true.", "A comma-ok lookup does not return an error."},
	},
	revealText: [2]string{"تابع فحص معرف حدث سبق تسجيله.", "Follow a lookup for an event ID already recorded."},
	reveal: ch35ExampleSpec{
		id: "e354-event-seen", imports: []string{"fmt"},
		code: `func main() {
	seen := map[int64]bool{4: true}
	_, ok := seen[4]
	fmt.Println(ok)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"المعرف موجود في السجل.", "The ID is already recorded."}, vars: []LessonVariable{{Name: "seen", Value: "map[4:true]"}}},
			{line: 3, caption: [2]string{"اتأكدنا من وجوده.", "The lookup confirms it exists."}, vars: []LessonVariable{{Name: "ok", Value: "true"}}},
			{line: 4, caption: [2]string{"النتيجة true.", "The result is true."}, vars: []LessonVariable{{Name: "ok", Value: "true"}}, output: "true"},
		},
	},
	watch: [2]string{
		"الطلب والـoutbox لازم يبقوا في transaction واحدة. والـrelay ممكن يعيد الحدث بعد نجاح النشر؛ المستهلك يتعامل مع duplicate.",
		"The order and outbox row belong in one transaction. A relay may repeat an event after publishing; the consumer handles duplicates.",
	},
	plan: [2]string{
		"افصل كتابة الطلب عن الترحيل والاستهلاك. تتبع ترتيب النشر والتحديث، وفكّر في نسخ schema وفي حالة consumer المشتركة.",
		"Separate order writing, relaying, and consuming. Trace publish and update order, then consider schema versions and shared consumer state.",
	},
	checkQuestion: [2]string{"إمتى يتسجل الحدث على إنه اتنشر؟", "When should an event be marked as published?"},
	checkOptions: [2][]string{
		{"بعد نجاح Publish.", "قبل استدعاء Publisher.", "وقت إنشاء الطلب."},
		{"After Publish succeeds.", "Before calling the publisher.", "When the order is created."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"لو النشر فشل، الحدث يفضل قابل لإعادة المحاولة.", "ده ممكن يفقد الحدث لو النشر فشل.", "إنشاء الطلب مش معناه إنه اتبعت."},
		{"A failed publish leaves the event available to retry.", "That can lose the event when publishing fails.", "Creating an order does not mean it was sent."},
	},
	turn: [2]string{
		"راجع أنواع ~Event~ و~Publisher~ و~Consumer~ في الـstarter. اشتغل من كل عقد لوحده، وسيب ~main~ فاضية.",
		"Review ~Event~, ~Publisher~, and ~Consumer~ in the starter. Work through each contract separately, and leave ~main~ empty.",
	},
	hints: [2][3]string{
		{"الـoutbox بيحفظ نية النشر مع بيانات الطلب.", "قسّم المسار لكتابة ذرّية، relay مرتب، واستهلاك idempotent.", "ابدأ من ~PlaceOrder~ و~Relay~ و~Handle~ في الـstarter."},
		{"The outbox stores publishing intent with the order data.", "Separate atomic writing, ordered relay, and idempotent consumption.", "Start at ~PlaceOrder~, ~Relay~, and ~Handle~ in the starter."},
	},
})

var ch35KafkaArabic, ch35KafkaEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c355.",
	mission: [2]string{
		"ابعت قراءات العربات لـ Kafka بمفتاح لكل عربة، واستهلك الرسائل مع retry. أكد الـoffset بعد نجاح المعالجة، وتعامل مع الإغلاق من غير ما تخسر الرسالة.",
		"Send rover readings to Kafka with a key per rover, then consume with retries. Commit the offset after processing succeeds, and shut down without losing a message.",
	},
	ideaText: [2]string{
		"مفتاح Kafka وقيمتها bytes. المفتاح يعرّف مجموعة الرسائل اللي لازم تفضل مرتبة جوه partition واحد.",
		"A Kafka key and value are bytes. The key identifies the messages that should stay ordered within one partition.",
	},
	idea: ch35ExampleSpec{
		id: "e355-message-bytes", imports: []string{"fmt"},
		code: `func main() {
	key := []byte("rover-7")
	value := []byte("42")
	fmt.Println(string(key), string(value))
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"المفتاح بايتات العربة.", "The key bytes identify a rover."}, vars: []LessonVariable{{Name: "key", Value: "[114 111 118 101 114 45 55]"}}},
			{line: 3, caption: [2]string{"القيمة بايتات القراءة.", "The value bytes hold a reading."}, vars: []LessonVariable{{Name: "value", Value: "[52 50]"}}},
			{line: 4, caption: [2]string{"رجعنا النصين للعرض.", "Both byte slices are shown as text."}, vars: []LessonVariable{{Name: "key", Value: "rover-7"}, {Name: "value", Value: "42"}}, output: "rover-7 42"},
		},
	},
	question: [2]string{"إيه الناتج عند تحويل نص مش رقمي؟", "What is printed when converting non-numeric text?"},
	options:  [2][]string{{"0 true", "4 true", "4 false"}, {"0 true", "4 true", "4 false"}}, answer: 0,
	feedback: [2][]string{
		{"التحويل الفاشل بيرجع صفر وخطأ.", "النص مش رقم 4.", "فيه error، مش نتيجة ناجحة."},
		{"A failed conversion returns zero and an error.", "The text is not the number 4.", "An error means the conversion did not succeed."},
	},
	revealText: [2]string{"راقب تحويل القيمة وفحص error.", "Watch the value conversion and error check."},
	reveal: ch35ExampleSpec{
		id: "e355-invalid-reading", imports: []string{"fmt", "strconv"},
		code: `func main() {
	reading, err := strconv.Atoi("4x")
	fmt.Println(reading, err != nil)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"حاولنا فك نص القراءة.", "The reading text is parsed."}, vars: []LessonVariable{{Name: "reading", Value: "0"}}},
			{line: 3, caption: [2]string{"التحويل فشل؛ الخطأ موجود.", "Parsing failed; an error exists."}, vars: []LessonVariable{{Name: "reading", Value: "0"}, {Name: "err != nil", Value: "true"}}, output: "0 true"},
		},
	},
	watch: [2]string{
		"الـpartition يضمن ترتيبًا جوه نفسه بس، مش على الـtopic كله. و~FetchMessage~ مش بيأكد الرسالة تلقائيًا.",
		"A partition preserves its own order, not the whole topic's. ~FetchMessage~ does not commit the message automatically.",
	},
	plan: [2]string{
		"راجع مسار الإرسال منفصل عن الاستهلاك. في الاستهلاك، ميّز الرسالة غير الصالحة عن فشل المعالج وعن إلغاء السياق.",
		"Review sending separately from consuming. In the consumer, distinguish invalid messages, handler failures, and context cancellation.",
	},
	checkQuestion: [2]string{"بعد ما ينجح handler، بإيه نعمل commit؟", "After the handler succeeds, what should commit use?"},
	checkOptions: [2][]string{
		{"سياق ما اتلغاش مع الرسالة.", "سياق ملغي دائمًا.", "من غير context."},
		{"A context not canceled with the message.", "A context that is always canceled.", "No context."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"المعالجة خلصت، فلازم التأكيد يقدر يكمل.", "ده يمنع commit بعد معالجة ناجحة.", "واجهة commit بتحتاج context."},
		{"Processing finished, so confirmation must be able to continue.", "That blocks commit after successful processing.", "The commit API requires a context."},
	},
	turn: [2]string{
		"اقرأ واجهتي ~Writer~ و~Reader~ وتوقيعات ~Send~ و~Consume~. تتبع العداد والـoffset مع كل نتيجة، وسيب ~main~ فاضية.",
		"Read the ~Writer~ and ~Reader~ interfaces and the ~Send~ and ~Consume~ signatures. Track the count and offset outcomes, and leave ~main~ empty.",
	},
	hints: [2][3]string{
		{"المفتاح يربط رسائل العربة بنفس مسار الترتيب.", "افصل تحويل البيانات عن الإرسال، والمعالجة عن commit.", "ابدأ من ~Writer~ و~Reader~، وراجع التعامل مع ~ctx.Err~ في الـstarter."},
		{"A key keeps one rover's messages on the same ordering path.", "Separate data conversion from sending, and handling from committing.", "Start at ~Writer~ and ~Reader~, and review ~ctx.Err~ handling in the starter."},
	},
})

var ch35RaftArabic, ch35RaftEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c356.",
	mission: [2]string{
		"كمّل محاكاة Raft حتمية في الذاكرة عشان محطات المدار تختار قائدًا. تابع الفترة والدور والرسائل، وخلي الانفصال جزءًا من تفكيرك.",
		"Complete a deterministic in-memory Raft simulation so orbital stations can elect a leader. Track term, role, and messages, and account for partitions.",
	},
	ideaText: [2]string{
		"~iota~ بتدي قيم متتابعة للـRole. تابع الرقم اللي اتخزن لكل اسم في الثوابت.",
		"~iota~ gives Role constants consecutive values. Follow the number assigned to each constant name.",
	},
	idea: ch35ExampleSpec{
		id: "e356-role-iota", imports: []string{"fmt"},
		code: `type Role int
const (
	Follower Role = iota
	Candidate
	Leader
)
func main() {
	fmt.Println(Follower, Candidate, Leader)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 1, caption: [2]string{"Role نوع رقمي جديد.", "Role is a new numeric type."}},
			{line: 3, caption: [2]string{"أول قيمة تبدأ من صفر.", "The first value starts at zero."}, vars: []LessonVariable{{Name: "Follower", Value: "0"}}},
			{line: 4, caption: [2]string{"القيمة التالية تزيد واحد.", "The next constant increments by one."}, vars: []LessonVariable{{Name: "Candidate", Value: "1"}}},
			{line: 5, caption: [2]string{"كل Role لها رقم.", "Each role has a number."}, vars: []LessonVariable{{Name: "Leader", Value: "2"}}},
			{line: 8, caption: [2]string{"اطبع القيم بالترتيب.", "Print the values in order."}, output: "0 1 2"},
		},
	},
	question: [2]string{"بـseed ثابت، إيه قيمة السحب ده؟", "With a fixed seed, what value does this draw?"},
	options:  [2][]string{{"6", "تتغير كل مرة", "10"}, {"6", "It changes each run", "10"}}, answer: 0,
	feedback: [2][]string{
		{"المصدر نفسه والـseed نفسه يعيدوا التسلسل.", "المثال لا يستخدم عشوائية النظام.", "Intn يرجع من المجال المحدد في المثال."},
		{"The same source and seed reproduce the sequence.", "This example does not use system randomness.", "Intn returns from the range in this example."},
	},
	revealText: [2]string{"نفّذ سحبًا من مولّد محلي ببذرة ثابتة.", "Run a draw from a local generator with a fixed seed."},
	reveal: ch35ExampleSpec{
		id: "e356-seeded-random", imports: []string{"fmt", "math/rand"},
		code: `func main() {
	rng := rand.New(rand.NewSource(7))
	fmt.Println(rng.Intn(10))
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"أنشأنا مولد ببذرة 7.", "A generator starts with seed 7."}, vars: []LessonVariable{{Name: "rng", Value: "*rand.Rand"}}},
			{line: 3, caption: [2]string{"السحب الأول ثابت هنا.", "The first draw is fixed here."}, vars: []LessonVariable{{Name: "rng.Intn(10)", Value: "6"}}, output: "6"},
		},
	},
	watch: [2]string{
		"الرسالة ذات Term أقدم ماينفعش تغيّر حالة أحدث. والـsimulation بترتب الرسائل من مدخلاتها؛ مافيش goroutine scheduling هنا.",
		"A message with an older term must not replace newer state. The simulation orders messages from its inputs; there is no goroutine scheduling here.",
	},
	plan: [2]string{
		"افصل تحديث الحالة عن بناء الردود. تابع إمتى المحطة follower أو candidate أو leader، وخلي مصدر العشوائية وتوقيته واضحين.",
		"Separate state updates from reply construction. Track when a station is follower, candidate, or leader, and keep randomness and its timing explicit.",
	},
	checkQuestion: [2]string{"وصلت رسالة Term بتاعها أعلى؛ إيه الأولوية؟", "A message arrives with a higher term; what takes priority?"},
	checkOptions: [2][]string{
		{"تبنّي الفترة الأحدث والرجوع follower.", "تجاهلها عشان الدور الحالي.", "اعتبار المرسل قائدًا قديمًا."},
		{"Adopt the newer term and step down to follower.", "Ignore it because of the current role.", "Treat the sender as an old leader."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"الفترة الأكبر هي الحالة الأحدث المعروفة.", "الدور الحالي مايلغيش term أعلى.", "المرسل مش قديم في المثال."},
		{"The larger term is the newest known state.", "The current role does not override a higher term.", "The sender is not stale in this example."},
	},
	turn: [2]string{
		"ابدأ من حالة ~NewNode~ ورسائل ~Tick~ و~Step~. تتبع رسالة واحدة في كل مرة، واكتب المحاكاة من غير وقت حقيقي أو نوم.",
		"Start with ~NewNode~ state and the ~Tick~ and ~Step~ messages. Trace one message at a time, without real time or sleeps.",
	},
	hints: [2][3]string{
		{"الـterm بيميز الحالة الجديدة من رسائل قائد قديم.", "افصل clock tick عن تسليم الرسالة وردها.", "ابدأ من ~State~ و~Tick~ و~Step~، وراجع توقيت سحب timeout في الـstarter."},
		{"A term distinguishes new state from messages by an old leader.", "Separate clock ticks from message delivery and replies.", "Start at ~State~, ~Tick~, and ~Step~, and review timeout draws in the starter."},
	},
})

var ch35ReplicationArabic, ch35ReplicationEnglish = ch35BuildLesson(ch35LessonSpec{
	prefix: "c357.",
	mission: [2]string{
		"وسّع محاكاة الانتخابات بسجل أوامر يتنسخ بين المحطات. القائد يقترح، المتابع يراجع التطابق، والمحطات تطبق الجزء المؤكد بنفس الترتيب.",
		"Extend the election simulation with a command log copied between stations. The leader proposes, followers check consistency, and stations apply the committed prefix in order.",
	},
	ideaText: [2]string{
		"كل Entry فيها command وterm. رقم الخانة يبدأ من واحد، لكن index الشريحة يبدأ من صفر؛ تابع الاتنين في المثال.",
		"Each Entry holds a command and term. Slot numbers start at one, while slice indexes start at zero; follow both in the example.",
	},
	idea: ch35ExampleSpec{
		id: "e357-log-index", imports: []string{"fmt"},
		code: `type Entry struct {
	Term int
	Cmd  string
}
func main() {
	log := []Entry{{Term: 1, Cmd: "scan"}, {Term: 2, Cmd: "dock"}}
	fmt.Println(len(log), log[1].Term)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 1, caption: [2]string{"Entry بتجمع term والأمر.", "Entry groups a term and command."}},
			{line: 6, caption: [2]string{"السجل فيه خانتين.", "The log has two entries."}, vars: []LessonVariable{{Name: "log", Value: "[{1 scan} {2 dock}]"}}},
			{line: 7, caption: [2]string{"index 1 يحمل term 2.", "Index 1 carries term 2."}, vars: []LessonVariable{{Name: "len(log)", Value: "2"}, {Name: "log[1].Term", Value: "2"}}, output: "2 2"},
		},
	},
	question: [2]string{"الأمر واحد، لكن الـTerm مختلف. هل الـEntry قيمتين متساويتين؟", "The command is the same, but the term differs. Are the two entries equal?"},
	options:  [2][]string{{"true", "false", "مش قابلين للمقارنة"}, {"true", "false", "They cannot be compared"}}, answer: 1,
	feedback: [2][]string{
		{"كل الحقول تدخل في مقارنة الـstruct.", "اختلاف Term يخلي القيمتين مختلفتين.", "حقول النوع قابلة للمقارنة."},
		{"Every field participates in struct comparison.", "Different terms make the values different.", "The type's fields are comparable."},
	},
	revealText: [2]string{"قارن سجلين في نفس الخانة لهما term مختلف.", "Compare two entries at one slot with different terms."},
	reveal: ch35ExampleSpec{
		id: "e357-entry-term", imports: []string{"fmt"},
		code: `type Entry struct {
	Term int
	Cmd  string
}
func main() {
	local := Entry{Term: 1, Cmd: "scan"}
	incoming := Entry{Term: 2, Cmd: "scan"}
	fmt.Println(local == incoming)
}`,
		frames: []ch35FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 1, caption: [2]string{"السجل فيه term وأمر.", "An entry holds a term and command."}},
			{line: 6, caption: [2]string{"الخانة المحلية term 1.", "The local slot has term 1."}, vars: []LessonVariable{{Name: "local", Value: "{1 scan}"}}},
			{line: 7, caption: [2]string{"الواردة لها term 2.", "The incoming slot has term 2."}, vars: []LessonVariable{{Name: "incoming", Value: "{2 scan}"}}},
			{line: 8, caption: [2]string{"الاختلاف يجعلها false.", "That difference makes it false."}, vars: []LessonVariable{{Name: "local == incoming", Value: "false"}}, output: "false"},
		},
	},
	watch: [2]string{
		"تطابق command لوحده مش كفاية؛ راجع PrevIndex وPrevTerm قبل الإضافة. وما تأكدش خانة لسه مش موجودة عند المتابع.",
		"Matching commands alone is not enough; check PrevIndex and PrevTerm before appending. Do not commit a slot a follower does not have.",
	},
	plan: [2]string{
		"افصل الاقتراح، مزامنة سجل المتابع، وتقدم commit. راجع كمان قواعد انتخاب القائد، لأن صلاحية السجل بتأثر على التصويت.",
		"Separate proposals, follower log matching, and commit progress. Keep election rules too, because log freshness affects voting.",
	},
	checkQuestion: [2]string{"المتبع وصلته قيمة Commit أكبر من سجله المطابق. يعمل إيه؟", "A follower receives Commit beyond its matched log. What should it do?"},
	checkOptions: [2][]string{
		{"يحد التقدم عند آخر خانة مطابقة.", "يضيف خانات وهمية.", "يقلل commit الحالي."},
		{"Limit progress to the last matched slot.", "Create placeholder entries.", "Decrease its current commit."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"المتابع ماعندوش محتوى يؤكد بعد الخانة دي.", "الفهرس مش بديل عن سجل موجود.", "حالة commit ماينفعش ترجع لورا."},
		{"The follower has no content to commit beyond that slot.", "An index cannot replace an actual log entry.", "Commit state must not move backward."},
	},
	turn: [2]string{
		"كمّل على كود الانتخابات في الـstarter. تتبع رسائل Append وردودها، وسيب تحديد الأوامر المؤكدة ناتجًا عن حالة السجل.",
		"Continue from the election code in the starter. Trace Append messages and replies, and derive committed commands from log state.",
	},
	hints: [2][3]string{
		{"الـEntry بتجمع command مع الفترة اللي اتكتب فيها.", "افصل فحص التطابق عن تحديث السجل وتقدم commit.", "ابدأ من ~Msg~ و~Node.Step~ في الـstarter، وراجع معنى ~PrevIndex~ و~PrevTerm~."},
		{"An ~Entry~ groups a command with the term when it was written.", "Separate consistency checks from log updates and commit progress.", "Start at ~Msg~ and ~Node.Step~ in the starter, and review ~PrevIndex~ and ~PrevTerm~."},
	},
})

var ch35LessonArabic, ch35LessonEnglish = func() (map[string][]LessonStep, map[string][]LessonStep) {
	return map[string][]LessonStep{
		"35-1-task-queue":    ch35QueueArabic,
		"35-2-scheduler":     ch35SchedulerArabic,
		"35-3-redis-lock":    ch35LockArabic,
		"35-4-outbox":        ch35OutboxArabic,
		"35-5-kafka":         ch35KafkaArabic,
		"35-6-raft-election": ch35RaftArabic,
		"35-7-raft-log":      ch35ReplicationArabic,
	}, map[string][]LessonStep{
		"35-1-task-queue":    ch35QueueEnglish,
		"35-2-scheduler":     ch35SchedulerEnglish,
		"35-3-redis-lock":    ch35LockEnglish,
		"35-4-outbox":        ch35OutboxEnglish,
		"35-5-kafka":         ch35KafkaEnglish,
		"35-6-raft-election": ch35RaftEnglish,
		"35-7-raft-log":      ch35ReplicationEnglish,
	}
}()
