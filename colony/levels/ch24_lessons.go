package levels

type ch24Copy struct {
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

type ch24FrameInfo struct {
	line      int
	arabic    string
	english   string
	variables []LessonVariable
	output    string
}

type ch24Example struct {
	id      string
	code    string
	imports []string
	frames  []ch24FrameInfo
}

func ch24Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch24Question(id, text string, options, feedback []string, answer int) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Feedback: feedback, Answer: answer}}
}

func ch24Walkthrough(id, text string, example ch24Example, english bool) LessonBlock {
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

func ch24BuildLesson(prefix string, copy ch24Copy, examples [2]ch24Example, english bool) []LessonStep {
	titles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	if english {
		titles = []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
	}
	idea := ch24Walkthrough(prefix+".idea-example", copy.idea, examples[0], english)
	tryCode := copy.tryQuestion + "\n\n~~~go\n" + examples[1].code + "\n~~~"
	try := []LessonBlock{
		ch24Question(prefix+".try-predict", tryCode, copy.tryOptions, copy.tryFeedback, 1),
		ch24Walkthrough(prefix+".try-reveal", copy.reveal, examples[1], english),
	}
	check := ch24Question(prefix+".check-question", copy.checkQuestion, copy.checkOptions, copy.checkFeedback, 1)
	return []LessonStep{
		{ID: prefix + ".mission-step", Title: titles[0], Blocks: []LessonBlock{ch24Markdown(prefix+".mission", copy.mission)}},
		{ID: prefix + ".idea-step", Title: titles[1], Blocks: []LessonBlock{idea}},
		{ID: prefix + ".try-step", Title: titles[2], Blocks: try},
		{ID: prefix + ".watch-step", Title: titles[3], Blocks: []LessonBlock{ch24Markdown(prefix+".watch", copy.watch)}},
		{ID: prefix + ".check-step", Title: titles[4], Blocks: []LessonBlock{ch24Markdown(prefix+".think", copy.think), check}},
		{ID: prefix + ".turn-step", Title: titles[5], Blocks: []LessonBlock{ch24Markdown(prefix+".turn", copy.turn)}},
	}
}

var ch24PingExamples = [2]ch24Example{
	{
		id: "e241-channel-handoff", imports: []string{"fmt"},
		code: `func main() {
	permit := make(chan int)
	go func() { permit <- 7 }()
	code := <-permit
	fmt.Println(code)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "قناة بلا مخزن.", english: "The channel has no buffer."},
			{line: 3, arabic: "مرسل واحد ينتظر.", english: "One sender waits."},
			{line: 4, arabic: "الاستقبال يستلم 7.", english: "The receive gets 7."},
			{line: 5, arabic: "اطبع التصريح.", english: "Print the permit.", output: "7"},
		},
	},
	{
		id: "e241-buffer-read", imports: []string{"fmt"},
		code: `func main() {
	queue := make(chan int, 1)
	queue <- 4
	fmt.Println(<-queue)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "السعة خانة واحدة.", english: "Capacity is one slot."},
			{line: 3, arabic: "الرسالة تدخل الطابور.", english: "The value enters the queue."},
			{line: 4, arabic: "الاستقبال يطبع 4.", english: "The receive prints 4.", output: "4"},
		},
	},
}

var ch24PingArabic = ch24BuildLesson("c241", ch24Copy{
	mission:       "نفّذ ~RelaySupplies~ لتسليم جولات بين درونين. راجع الإلغاء والأخطاء ونهاية القنوات من غير goroutine معلّقة. الـstarter فيه التوقيع بس.",
	idea:          "القناة unbuffered زي مقابلة: الإرسال والاستقبال يلاقوا بعض. ده ترتيب واحد ممكن بين goroutines، مش وعد بترتيب تشغيل ثابت.",
	tryQuestion:   "توقّع اللي هيتطبع بعد إرسال القيمة:",
	tryOptions:    []string{"0", "4", "مافيش ناتج"},
	tryFeedback:   []string{"القناة فيها قيمة قبل الاستقبال.", "الاستقبال بياخد القيمة المرسلة.", "الإرسال يسبق الاستقبال في نفس المسار."},
	reveal:        "تابع دخول القيمة للقناة ثم الاستقبال.",
	watch:         "- إشارة الاستلام مش معناها إن callback خلصت.\n- قافل القناة لازم يملك قرار آخر إرسال.",
	think:         "قسّم العقد: مدخلات صالحة، جولات متبادلة، ثم إيقاف واضح. راجع مين منتظر إيه عند النجاح والخطأ.\n\nاكتب كل إرسال واستقبال كسهم واضح.",
	checkQuestion: "إيه اللي يميّز نهاية قناة مقفولة وفاضية؟",
	checkOptions:  []string{"القيمة 0 بس", "الاستقبال بيرجع ~ok == false~", "إرسال جديد"},
	checkFeedback: []string{"الصفر ممكن يكون قيمة حقيقية.", "~ok~ بيفرق النهاية عن قيمة صفر عادية.", "الإرسال بعد القفل يسبب panic."},
	turn:          "ابدأ من التوقيع والـcontext. شغّل التستات بعد ما تحدد ملكية كل قناة ونقطة انتظارها.",
}, ch24PingExamples, false)

var ch24PingEnglish = ch24BuildLesson("c241", ch24Copy{
	mission:       "Implement ~RelaySupplies~ to hand rounds between two drones.\n\nReview cancellation, errors, and channel shutdown without a stranded goroutine. The starter has only the signature.",
	idea:          "An unbuffered channel is a meeting point: send and receive meet.\n\nThis shows one possible goroutine order, not a fixed scheduling promise.",
	tryQuestion:   "Predict what prints after the send:",
	tryOptions:    []string{"0", "4", "No output"},
	tryFeedback:   []string{"A value is waiting before the receive.", "The receive takes the sent value.", "The send comes before the receive here."},
	reveal:        "Follow the value into the channel, then receive it.",
	watch:         "- Receipt does not mean the callback has finished.\n- The channel owner decides when the final send is done.",
	think:         "Split the contract into valid inputs, alternating rounds, and a clear stop.\n\nCheck who waits for what on success and failure.\n\nDraw each send and receive as a clear arrow.",
	checkQuestion: "How can a receive detect a closed, empty channel?",
	checkOptions:  []string{"Only the value 0", "The receive returns ~ok == false~", "A new send"},
	checkFeedback: []string{"Zero can also be a real value.", "~ok~ distinguishes closure from an ordinary zero.", "Sending after close causes a panic."},
	turn:          "Start from the signature and context. Run the lab checks after mapping each channel owner and wait point.",
}, ch24PingExamples, true)

var ch24DepotExamples = [2]ch24Example{
	{
		id: "e242-buffer-capacity", imports: []string{"fmt"},
		code: `func main() {
	queue := make(chan int, 2)
	queue <- 4
	queue <- 6
	fmt.Println(cap(queue), len(queue))
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "سعة الطابور اتنين.", english: "Queue capacity is two."},
			{line: 3, arabic: "أول صندوق يدخل.", english: "The first crate enters."},
			{line: 4, arabic: "التاني يفضل بعده.", english: "The next one waits behind it."},
			{line: 5, arabic: "السعة والطول: 2 و2.", english: "Capacity and length: 2 and 2.", output: "2 2"},
		},
	},
	{
		id: "e242-buffer-one", imports: []string{"fmt"},
		code: `func main() {
	queue := make(chan int, 1)
	queue <- 8
	fmt.Println(cap(queue), len(queue))
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "فيه خانة واحدة.", english: "There is one slot."},
			{line: 3, arabic: "الخانة بقت مشغولة.", english: "The slot is now occupied."},
			{line: 4, arabic: "السعة والطول: 1 و1.", english: "Capacity and length: 1 and 1.", output: "1 1"},
		},
	},
}

var ch24DepotArabic = ch24BuildLesson("c242", ch24Copy{
	mission:       "اجمع صناديق ~batches~ عبر ~UnloadDepots~ من غير فقد. الـlab يفحص استهلاكًا واحدًا و~capacity~ محدودة؛ الـstarter فيه التوقيع بس.",
	idea:          "القناة ذات buffer صف انتظار محدود. ~cap~ هي عدد الخانات، و~len~ عدد الرسائل المنتظرة دلوقتي.",
	tryQuestion:   "توقّع قيمتي ~cap~ و~len~ بعد الإرسال:",
	tryOptions:    []string{"1 و0", "1 و1", "0 و1"},
	tryFeedback:   []string{"السعة ما بتقلش لما الرسالة تدخل.", "خانة واحدة متاحة ومشغولة.", "السعة هنا مش صفر."},
	reveal:        "شوف الفرق بين سعة القناة والرسائل اللي جواها.",
	watch:         "- وجود buffer ما يضمنش وصول الرسالة للمستهلك.\n- حالة التجميع المشتركة تحتاج مالك واضح.",
	think:         "ابدأ بمسار الرسالة من كل batch للمستهلك. راجع مين بيتعامل مع الخطأ وإمتى المنتجين يوقفوا.\n\nكل batch تحتفظ بترتيبها الداخلي؛ تابع القناة حتى النهاية.\n\nعدد العمال مش دليل على صحة القيم. الحفاظ على الترتيب جزء من العقد.",
	checkQuestion: "لو القناة ممتلئة، الإرسال العادي يعمل إيه؟",
	checkOptions:  []string{"ينتظر مستقبلًا", "يمسح أقدم قيمة", "يوسّع السعة"},
	checkFeedback: []string{"الإرسال يكمّل لما تتوفر خانة.", "القنوات ما بتحذفش تلقائيًا.", "السعة ما بتتغيرش بعد الإنشاء."},
	turn:          "جرّب تتبع كل قيمة من المصدر للاستهلاك. راقب مسارات الإلغاء والخطأ كمان.",
}, ch24DepotExamples, false)

var ch24DepotEnglish = ch24BuildLesson("c242", ch24Copy{
	mission:       "Merge the ~batches~ through ~UnloadDepots~ without losing crates.\n\nThe lab checks one consumer and a bounded ~capacity~; the starter has only the signature.",
	idea:          "A buffered channel is a bounded waiting line. ~cap~ is its number of slots; ~len~ is how many messages are waiting now.",
	tryQuestion:   "Predict ~cap~ and ~len~ after the send:",
	tryOptions:    []string{"1 and 0", "1 and 1", "0 and 1"},
	tryFeedback:   []string{"Capacity does not shrink when a value enters.", "One slot is available and occupied.", "This channel capacity is not zero."},
	reveal:        "Compare the channel capacity with the messages inside it.",
	watch:         "- A buffer does not guarantee delivery to a consumer.\n- Shared aggregation state needs a clear owner.",
	think:         "Trace each value from its batch to the consumer. Check who handles an error and when producers stop.\n\nEach batch keeps its internal order; follow the channel to the end.\n\nWorker count does not prove values are correct. Preserving order is part of the contract.",
	checkQuestion: "What does a regular send do when the channel is full?",
	checkOptions:  []string{"Waits for a receiver", "Deletes the oldest value", "Grows capacity"},
	checkFeedback: []string{"The send continues when a slot is free.", "Channels do not delete values automatically.", "Capacity does not change after creation."},
	turn:          "Trace each value from source to consumption. Review cancellation and error paths too.",
}, ch24DepotExamples, true)

var ch24CrewExamples = [2]ch24Example{
	{
		id: "e243-generic-reply", imports: []string{"fmt"},
		code: `type Reply[T any] struct { Value T }
func main() {
	result := Reply[int]{Value: 7}
	fmt.Println(result.Value)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 1, arabic: "T نوع بيتحدد لاحقًا.", english: "T is chosen at use."},
			{line: 3, arabic: "هنا القيمة int.", english: "Here the value is an int."},
			{line: 4, arabic: "اقرأ نتيجة النوع.", english: "Read the typed result.", output: "7"},
		},
	},
	{
		id: "e243-generic-text", imports: []string{"fmt"},
		code: `type Reply[T any] struct { Value T }
func main() {
	result := Reply[string]{Value: "ore"}
	fmt.Println(result.Value)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 1, arabic: "النوع يقبل أنواعًا مختلفة.", english: "The type accepts different types."},
			{line: 3, arabic: "هنا القيمة string.", english: "Here the value is a string."},
			{line: 4, arabic: "الناتج هو ore.", english: "The output is ore.", output: "ore"},
		},
	},
}

var ch24CrewArabic = ch24BuildLesson("c243", ch24Copy{
	mission:       "نفّذ ~RunCrew~ generic بطاقم ثابت و~CrewReply~ لكل مدخل. لكل نتيجة مكانها، والخطأ العام يوقف الشغل؛ الـstarter فيه الأنواع والتوقيع.",
	idea:          "الـtype parameter زي خانة نوع مؤقتة. عند استخدام ~Reply[int]~ تبقى ~Value~ من نوع int؛ شكل النتيجة يفضل نفسه.",
	tryQuestion:   "توقّع القيمة المطبوعة من ~Reply[string]~:",
	tryOptions:    []string{"7", "ore", "string"},
	tryFeedback:   []string{"دي قيمة المثال الرقمي.", "الـfield فيها النص اللي دخلناه.", "اسم النوع مش قيمة field."},
	reveal:        "تابع اختيار النوع ثم قراءة ~Value~.",
	watch:         "- نتيجة كل عنصر تختلف عن خطأ استدعاء الطاقم كله.\n- العمال ما يكتبوش في نفس خانة النتيجة.",
	think:         "اقرأ العقد وحدد معنى كل return. بعدين تتبع المدخل، نتيجته، وحالة الإلغاء من غير ما تخلطهم.\n\nالـindex يربط الرد بمدخله حتى لو اختلف ترتيب الوصول.\n\nترتيب الردود لازم يطابق ترتيب المدخلات. راجع ترتيبها قبل الاستنتاج.",
	checkQuestion: "في ~Reply[T]~، إيه دور ~T~؟",
	checkOptions:  []string{"عدد العمال", "نوع القيمة داخل Reply", "index المدخل"},
	checkFeedback: []string{"العمال إعداد منفصل عن نوع البيانات.", "كل استخدام يحدد النوع الفعلي لـT.", "الـindex يحدد ترتيب العنصر."},
	turn:          "ابدأ بقراءة أنواع ~T~ و~R~ في التوقيع. افحص ملكية كل نتيجة وعمر كل عامل.",
}, ch24CrewExamples, false)

var ch24CrewEnglish = ch24BuildLesson("c243", ch24Copy{
	mission:       "Implement generic ~RunCrew~ with a fixed crew and one ~CrewReply~ per input.\n\nEach result has a slot, and a crew-level error stops work; the starter has the types and signature.",
	idea:          "A type parameter is a temporary type slot. With ~Reply[int]~, ~Value~ becomes an int; the result shape stays the same.",
	tryQuestion:   "Predict the printed value from ~Reply[string]~:",
	tryOptions:    []string{"7", "ore", "string"},
	tryFeedback:   []string{"That is the numeric example's value.", "The field holds the text we supplied.", "A type name is not a field value."},
	reveal:        "Follow the type choice, then read ~Value~.",
	watch:         "- One item's result differs from a crew-wide error.\n- Workers must not write the same result slot.",
	think:         "Read the contract and identify each return value. Then trace an input, its result, and cancellation without mixing them.\n\nThe index ties each reply to its input if arrival order changes.\n\nReply order must match input order. Review that order before drawing conclusions.",
	checkQuestion: "In ~Reply[T]~, what does ~T~ represent?",
	checkOptions:  []string{"The worker count", "The value type inside Reply", "The input index"},
	checkFeedback: []string{"Worker count is separate from data type.", "Each use chooses the actual type for T.", "The index identifies an item order."},
	turn:          "Start by reading ~T~ and ~R~ in the signature. Check each result's owner and each worker's lifetime.",
}, ch24CrewExamples, true)

var ch24PipelineExamples = [2]ch24Example{
	{
		id: "e244-drain-closed", imports: []string{"fmt"},
		code: `func main() {
	jobs := make(chan int, 2)
	jobs <- 3
	jobs <- 5
	close(jobs)
	for job := range jobs {
		fmt.Println(job)
	}
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الطابور فيه خانتين.", english: "The queue has two slots."},
			{line: 5, arabic: "المرسل أعلن النهاية.", english: "The sender announced the end."},
			{line: 7, arabic: "أول قراءة تطبع 3.", english: "The first receive prints 3.", output: "3"},
			{line: 7, arabic: "التانية تطبع 5.", english: "The next prints 5.", output: "3\n5"},
		},
	},
	{
		id: "e244-closed-read", imports: []string{"fmt"},
		code: `func main() {
	ch := make(chan int)
	close(ch)
	value, ok := <-ch
	fmt.Println(value, ok)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 3, arabic: "القناة اتقفلت فاضية.", english: "The empty channel is closed."},
			{line: 4, arabic: "القيمة صفر وok خطأ.", english: "The value is zero and ok is false."},
			{line: 5, arabic: "دي علامة النهاية.", english: "That is the end signal.", output: "0 false"},
		},
	},
}

var ch24PipelineArabic = ch24BuildLesson("c244", ch24Copy{
	mission:       "نفّذ ~RefineOre~ عشان يقسم شغل الخامات على lanes ويرجع stream ردود. الـlab يراجع الإلغاء والإغلاق؛ الـstarter محدد الأنواع والتوقيع.",
	idea:          "~close~ يقول إن مفيش رسائل جديدة. ~range~ يستهلك الموجود ويقف لوحده بعد ما القناة تتقفل وتفضى.",
	tryQuestion:   "القناة فاضية ومقفولة؛ توقّع نتيجة الاستقبال:",
	tryOptions:    []string{"0 وfalse", "0 وtrue", "panic"},
	tryFeedback:   []string{"القيمة الصفرية ومعاها علامة النهاية.", "true معناها وصلت قيمة فعلية.", "الاستقبال من قناة مقفولة مسموح."},
	reveal:        "اتبع علامة ~ok~ بعد قفل قناة فاضية.",
	watch:         "- اقفل الإرسال بعد آخر مرسل، مش بعد أول عنصر.\n- الإلغاء لازم يفضل ممكن لو مفيش مستقبل.",
	think:         "ارسم مسار كل batch: قناة lane، تحويل، ثم الدمج. حدّد ملكية كل إغلاق ونقطة انتظار.\n\nالخرج النهائي له نقطة إغلاق واحدة.",
	checkQuestion: "إمتى ~range~ على قناة بيخلص؟",
	checkOptions:  []string{"بعد أول قيمة", "لما القناة تتقفل وتفضى", "بعد عدد ثابت"},
	checkFeedback: []string{"range يكمل مع القيم المتاحة.", "القفل مع تفريغ القناة ينهي الحلقة.", "عدد الرسائل مش لازم يبقى معروف."},
	turn:          "ابدأ بتحديد اتجاه القنوات ومن يملك إغلاقها. بعدين تتبع الإلغاء في كل انتظار.",
}, ch24PipelineExamples, false)

var ch24PipelineEnglish = ch24BuildLesson("c244", ch24Copy{
	mission:       "Implement ~RefineOre~ to split ore across lanes and return a reply stream.\n\nThe lab checks cancellation and closure; the starter defines the types and signature.",
	idea:          "~close~ says no more messages will be sent. ~range~ consumes what remains, then stops after the channel is closed and drained.",
	tryQuestion:   "The channel is empty and closed; predict the receive:",
	tryOptions:    []string{"0 and false", "0 and true", "panic"},
	tryFeedback:   []string{"The zero value comes with the end marker.", "true means a real value arrived.", "Receiving from a closed channel is allowed."},
	reveal:        "Follow the ~ok~ result after closing an empty channel.",
	watch:         "- Close sends only after the final sender, not the first item.\n- Cancellation must remain possible when no receiver is ready.",
	think:         "Draw each batch path: lane channel, transform, then merge. Identify every close owner and wait point.\n\nThe final output has one close point.",
	checkQuestion: "When does ~range~ over a channel stop?",
	checkOptions:  []string{"After one value", "When closed and drained", "After a fixed count"},
	checkFeedback: []string{"range continues while values arrive.", "Closure plus an empty channel ends the loop.", "The message count need not be known."},
	turn:          "Start by marking channel directions and close owners. Then trace cancellation at every wait.",
}, ch24PipelineExamples, true)

var ch24BeaconExamples = [2]ch24Example{
	{
		id: "e245-context-request", imports: []string{"context", "fmt", "net/http"},
		code: `func main() {
	request, err := http.NewRequestWithContext(context.Background(), "GET", "https://beacon.example/ping", nil)
	if err != nil {
		fmt.Println("invalid")
		return
	}
	fmt.Println(request.Method, request.URL.Host)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الطلب مربوط بالسياق.", english: "The request carries its context."},
			{line: 3, arabic: "العنوان صالح.", english: "The address is valid."},
			{line: 6, arabic: "اقرأ الطريقة والـhost.", english: "Read the method and host.", output: "GET beacon.example"},
		},
	},
	{
		id: "e245-http-status", imports: []string{"fmt", "net/http"},
		code: `func main() {
	response := http.Response{StatusCode: http.StatusAccepted}
	fmt.Println(response.StatusCode)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الرد يحمل status.", english: "The response carries a status."},
			{line: 3, arabic: "اقرأ قيمة status.", english: "Read the status value.", output: "202"},
		},
	},
}

var ch24BeaconArabic = ch24BuildLesson("c245", ch24Copy{
	mission:       "نفّذ ~CheckBeacons~ بعميل HTTP محقون وحد فعلي للتوازي. النتائج تفضل بترتيب العناوين؛ الـstarter فيه ~BeaconReply~ والتوقيع.",
	idea:          "~NewRequestWithContext~ يربط الطلب بسياق الإلغاء. إنشاء الطلب لوحده مش اتصال، و~StatusCode~ يخص رد HTTP وصل فعلًا.",
	tryQuestion:   "توقّع status اللي يمثله الرد:",
	tryOptions:    []string{"200", "202", "404"},
	tryFeedback:   []string{"ده status نجاح شائع لكن مش الموجود هنا.", "~StatusAccepted~ قيمته 202.", "ده status لمورد مش موجود."},
	reveal:        "تابع الـconstant لحد قيمة status المطبوعة.",
	watch:         "- خطأ ~Do~ مش هو status HTTP.\n- الرد اللي له body يحتاج إغلاقه.",
	think:         "افصل إعداد الطلب عن تنفيذ الاتصال وعن تفسير الرد. راجع كمان مين يحتفظ بنتيجة كل عنوان.\n\nالعميل المحقون يحدد transport والـtimeout. status HTTP تفضل ردًا، مش خطأ تنفيذ.\n\nحدد العنوان غير الصالح قبل الاتصال، وسجل الرد في موضع عنوانه الأصلي.\n\nراجع إن كل response يتقفل، وافصل خطأ النقل. اكتب status بوضوح.",
	checkQuestion: "لو ~Do~ رجّع response من غير error، ~StatusCode~ بيوصف إيه؟",
	checkOptions:  []string{"نتيجة HTTP", "عدد الطلبات", "طول العنوان"},
	checkFeedback: []string{"status جزء من الرد المستلم.", "الطلبات ما تتحسبش من status.", "طول العنوان معلومة منفصلة."},
	turn:          "ابدأ من عقد العميل المحقون. تتبع الإلغاء والحد المتزامن قبل قراءة النتائج.",
}, ch24BeaconExamples, false)

var ch24BeaconEnglish = ch24BuildLesson("c245", ch24Copy{
	mission:       "Implement ~CheckBeacons~ with an injected HTTP client and a real concurrency limit.\n\nResults keep address order; the starter includes ~BeaconReply~ and the signature.",
	idea:          "~NewRequestWithContext~ ties a request to cancellation. Creating a request does not connect; ~StatusCode~ belongs to an HTTP response that arrived.",
	tryQuestion:   "Predict the status represented by this response:",
	tryOptions:    []string{"200", "202", "404"},
	tryFeedback:   []string{"That is a common success status, but not this one.", "~StatusAccepted~ has value 202.", "That status means the resource was not found."},
	reveal:        "Follow the constant to the printed status value.",
	watch:         "- A ~Do~ error is not an HTTP status.\n- A response with a body needs closing.",
	think:         "Separate request setup, performing the call, and interpreting the response. Also track who owns each address result.\n\nThe injected client controls transport and timeout. An HTTP status is a response, not an execution error.\n\nIdentify an invalid address before calling, and store each reply at its original address position.\n\nCheck that every response is closed, distinguish transport errors, and record status clearly.",
	checkQuestion: "If ~Do~ returns a response without an error, what does ~StatusCode~ describe?",
	checkOptions:  []string{"The HTTP result", "The request count", "The address length"},
	checkFeedback: []string{"Status is part of the received response.", "Status does not count requests.", "Address length is separate information."},
	turn:          "Start from the injected client contract. Trace cancellation and the concurrency limit before reading results.",
}, ch24BeaconExamples, true)

var ch24ChunkExamples = [2]ch24Example{
	{
		id: "e246-write-at-offset", imports: []string{"fmt", "os"},
		code: `func main() {
	file, err := os.CreateTemp("", "chart-example")
	if err != nil {
		fmt.Println("create failed")
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()
	n, err := file.WriteAt([]byte("map"), 2)
	fmt.Println(n, err == nil)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "الملف المؤقت اتفتح.", english: "The temporary file opened."},
			{line: 7, arabic: "التنظيف مؤجل.", english: "Cleanup is deferred."},
			{line: 8, arabic: "اكتب عند offset 2.", english: "Write at offset 2."},
			{line: 9, arabic: "ثلاثة bytes اتكتبت.", english: "Three bytes were written.", output: "3 true"},
		},
	},
	{
		id: "e246-write-count", imports: []string{"fmt", "os"},
		code: `func main() {
	file, _ := os.CreateTemp("", "part-example")
	defer os.Remove(file.Name())
	defer file.Close()
	count, _ := file.WriteAt([]byte("xy"), 4)
	fmt.Println(count)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "اتفتح ملف للتوضيح.", english: "A demo file was opened."},
			{line: 5, arabic: "اكتب حرفين عند offset 4.", english: "Write two bytes at offset 4."},
			{line: 6, arabic: "الناتج عدد bytes.", english: "The output is the byte count.", output: "2"},
		},
	},
}

var ch24ChunkArabic = ch24BuildLesson("c246", ch24Copy{
	mission:       "نزّل خريطة مقسمة في ~DownloadChart~، واكتب كل جزء في مرحلته، وانشرها بس بعد الاكتمال. الـstarter محدد ~ChartPlan~ و~ChartStage~.",
	idea:          "~WriteAt~ يكتب bytes عند offset محدد. كده الأجزاء تقدر تروح لمواضعها من غير ما ترتيب وصولها يحدد مكانها.",
	tryQuestion:   "توقّع قيمة ~count~ بعد كتابة حرفين:",
	tryOptions:    []string{"0", "2", "4"},
	tryFeedback:   []string{"الكتابة نجحت وفيها bytes.", "القيمة هي عدد bytes المكتوبة.", "offset مش عدد bytes المكتوبة."},
	reveal:        "اتبع قيمة الإرجاع من ~WriteAt~.",
	watch:         "- الرد الجزئي مش نفس حجم الجزء المتوقع.\n- المرحلة ما تتنشرش قبل نجاح كل الأجزاء.",
	think:         "فكّر في التخطيط، التحقق من كل رد، ثم قرار واحد للنشر أو الإلغاء. راجع حدود كل جزء وملكيته.\n\nكل نطاق جزء من ملف واحد؛ اكتب في موضعه ثم استنى القرار قبل النشر.",
	checkQuestion: "في ~WriteAt(data, offset)~، الـoffset يحدد إيه؟",
	checkOptions:  []string{"موضع بداية الكتابة", "عدد الـworkers", "حجم الملف النهائي"},
	checkFeedback: []string{"هو مكان بداية bytes في الوجهة.", "العمال إعداد مستقل.", "الحجم مش هو موضع الكتابة."},
	turn:          "ابدأ بعقد ~ChartStage~ وعميل HTTP. تتبع كل فشل حتى ~Commit~ أو الإلغاء.",
}, ch24ChunkExamples, false)

var ch24ChunkEnglish = ch24BuildLesson("c246", ch24Copy{
	mission:       "Download a split chart in ~DownloadChart~, write chunks to their stage positions, and publish only after completion.\n\nThe starter defines ~ChartPlan~ and ~ChartStage~.",
	idea:          "~WriteAt~ writes bytes at a chosen offset. Chunks can reach their positions without arrival order deciding where they go.",
	tryQuestion:   "Predict the value of ~count~ after writing two bytes:",
	tryOptions:    []string{"0", "2", "4"},
	tryFeedback:   []string{"The write succeeded and included bytes.", "The value counts bytes written.", "The offset is not the byte count."},
	reveal:        "Follow the return value from ~WriteAt~.",
	watch:         "- A partial response is not the expected chunk size.\n- Do not publish the stage before every chunk succeeds.",
	think:         "Think through planning, validating each response, then one publish-or-cancel decision. Check every range and its ownership.\n\nEach range belongs to one file; write at its position, then wait for the final decision before publishing.",
	checkQuestion: "In ~WriteAt(data, offset)~, what does offset identify?",
	checkOptions:  []string{"The starting write position", "The worker count", "The final file size"},
	checkFeedback: []string{"It is where the destination write starts.", "Worker count is a separate setting.", "File size is not the write position."},
	turn:          "Start from the ~ChartStage~ contract and HTTP client. Trace each failure through ~Commit~ or cancellation.",
}, ch24ChunkExamples, true)

var ch24TilesExamples = [2]ch24Example{
	{
		id: "e247-decode-config", imports: []string{"bytes", "fmt", "image", "image/png"},
		code: `func main() {
	pixels := image.NewRGBA(image.Rect(0, 0, 2, 1))
	var encoded bytes.Buffer
	_ = png.Encode(&encoded, pixels)
	config, format, err := image.DecodeConfig(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		fmt.Println("decode failed")
		return
	}
	fmt.Println(format, config.Width, config.Height)
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "صورة أبعادها 2×1.", english: "The image is 2 by 1."},
			{line: 4, arabic: "حوّل الصورة لـPNG.", english: "Encode the image as PNG."},
			{line: 5, arabic: "اقرأ بيانات الصورة.", english: "Read image configuration."},
			{line: 9, arabic: "الصيغة والعرض والارتفاع.", english: "Format, width, and height.", output: "png 2 1"},
		},
	},
	{
		id: "e247-rect-size", imports: []string{"fmt", "image"},
		code: `func main() {
	bounds := image.Rect(1, 2, 5, 8)
	fmt.Println(bounds.Dx(), bounds.Dy())
}`,
		frames: []ch24FrameInfo{
			{line: 0, arabic: "قبل التشغيل.", english: "Before execution."},
			{line: 2, arabic: "المستطيل 4×6.", english: "The rectangle is 4 by 6."},
			{line: 3, arabic: "اطبع العرض والارتفاع.", english: "Print width and height.", output: "4 6"},
		},
	},
}

var ch24TilesArabic = ch24BuildLesson("c247", ch24Copy{
	mission:       "ابنِ بلاطات أصغر بـ~BuildMapTiles~ مع طاقم محدود، وحافظ على format كل صورة. الـstarter يحدد ~MapTile~ و~TileSpec~ ونتيجة الحفظ.",
	idea:          "~DecodeConfig~ يقرأ أبعاد وصيغة الصورة من غير فك كل pixels. افحص الشكل قبل اختيار طريقة المعالجة.",
	tryQuestion:   "توقّع عرض وارتفاع ~image.Rect(1, 2, 5, 8)~:",
	tryOptions:    []string{"5 و8", "4 و6", "6 و4"},
	tryFeedback:   []string{"دي إحداثيات الحافة، مش الأبعاد.", "الفرق بين الحافتين هو 4×6.", "العرض والارتفاع مش مقلوبين."},
	reveal:        "تابع حساب الحجم من حدود المستطيل.",
	watch:         "- صورة source ممكن تكون PNG أو JPEG.\n- افحص الأسماء والإعدادات قبل تشغيل العمال.",
	think:         "قسّم المهمة لقراءة الإعدادات، معالجة كل بلاطة، ثم حفظ ردها. راجع ملكية كل نتيجة وسياق الحفظ.\n\nأبعاد الخرج لازم تكون موجبة. راجع الصيغة قبل الحفظ، وخلي كل صورة مستقلة.\n\nكل نتيجة تخص بلاطة واحدة.\n\nراجع اسم كل صورة.",
	checkQuestion: "إيه المعلومة اللي ترجع من ~DecodeConfig~؟",
	checkOptions:  []string{"أبعاد وصيغة الصورة", "كل pixels المفكوكة", "اسم ملف الحفظ"},
	checkFeedback: []string{"تقدر تفحص الإعدادات قبل فك pixels.", "دي خطوة تانية بعد قراءة config.", "اسم الحفظ يأتي من بيانات البلاطة."},
	turn:          "ابدأ من ~TileSpec~ وأنواع الرد. تتبع صورة واحدة ذهنيًا قبل مراجعة مسارات التوازي.",
}, ch24TilesExamples, false)

var ch24TilesEnglish = ch24BuildLesson("c247", ch24Copy{
	mission:       "Build smaller tiles with ~BuildMapTiles~ and a bounded crew, preserving each image format.\n\nThe starter defines ~MapTile~, ~TileSpec~, and the save result.",
	idea:          "~DecodeConfig~ reads an image's dimensions and format without decoding every pixel. Inspect the shape before choosing how to process it.",
	tryQuestion:   "Predict the width and height of ~image.Rect(1, 2, 5, 8)~:",
	tryOptions:    []string{"5 and 8", "4 and 6", "6 and 4"},
	tryFeedback:   []string{"Those are edge coordinates, not dimensions.", "The difference between edges is 4 by 6.", "Width and height are not reversed."},
	reveal:        "Follow the size calculation from the rectangle bounds.",
	watch:         "- A source image may be PNG or JPEG.\n- Validate names and settings before starting workers.",
	think:         "Split the work into reading settings, processing each tile, and saving its reply. Check result ownership and save context.\n\nOutput dimensions must be positive. Check the format before saving and keep tiles independent.\n\nEach result belongs to one tile. Check every image name.",
	checkQuestion: "What information does ~DecodeConfig~ return?",
	checkOptions:  []string{"Image dimensions and format", "Every decoded pixel", "The saved file name"},
	checkFeedback: []string{"Inspect settings before decoding pixels.", "That is a later step after reading config.", "The tile data supplies the save name."},
	turn:          "Start from ~TileSpec~ and the reply types. Trace one tile mentally before reviewing concurrency paths.",
}, ch24TilesExamples, true)
