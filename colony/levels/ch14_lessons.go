package levels

type ch14LessonCopy struct {
	Mission      string
	Idea         string
	TryQuestion  string
	TryReveal    string
	Watch        string
	Think        string
	Check        string
	Turn         string
	TryOptions   []string
	TryFeedback  []string
	CheckOptions []string
	CheckReply   []string
	Hints        []string
}

type ch14LessonFrame struct {
	Line      int
	Arabic    string
	English   string
	Variables []LessonVariable
	Output    string
}

type ch14LessonExample struct {
	ID      string
	Code    string
	Imports []string
	Frames  []ch14LessonFrame
}

type ch14LessonPair struct {
	Arabic       []LessonStep
	English      []LessonStep
	ArabicHints  []string
	EnglishHints []string
}

func ch14Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch14Question(id, text string, options, feedback []string, answer int) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch14Walkthrough(id, text string, example ch14LessonExample, arabic bool) LessonBlock {
	frames := make([]LessonFrame, len(example.Frames))
	for i, frame := range example.Frames {
		caption := frame.English
		if arabic {
			caption = frame.Arabic
		}
		frames[i] = LessonFrame{Line: frame.Line, Caption: caption, Variables: frame.Variables, Output: frame.Output}
	}
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: example.Code, Imports: example.Imports, ExampleID: example.ID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch14Locale(arabic, english string) map[string]any {
	return map[string]any{"ar": arabic, "en": english}
}

func ch14Sandbox(prefix string, arabic bool) LessonBlock {
	var kind, arabicText, englishText string
	var config map[string]any
	switch prefix {
	case "c141":
		kind = "ch14-generic-max"
		arabicText, englishText = "تابع المقارنة خطوة خطوة مع نوعين من الأرقام.", "Step through comparisons with two number types."
		config = map[string]any{
			"frames": []any{
				map[string]any{"type": ch14Locale("int", "int"), "values": []any{4, 9, 2}, "current": -1, "best": -1, "caption": ch14Locale("ابدأ بأرقام صحيحة.", "Start with integers.")},
				map[string]any{"type": ch14Locale("int", "int"), "values": []any{4, 9, 2}, "current": 0, "best": 0, "caption": ch14Locale("أول قيمة تبدأ best.", "The first value starts best.")},
				map[string]any{"type": ch14Locale("int", "int"), "values": []any{4, 9, 2}, "current": 1, "best": 1, "caption": ch14Locale("9 أكبر؛ best اتغيّرت.", "9 is larger; best changes.")},
				map[string]any{"type": ch14Locale("int", "int"), "values": []any{4, 9, 2}, "current": 2, "best": 1, "caption": ch14Locale("2 أصغر؛ best ثابتة.", "2 is smaller; best stays.")},
				map[string]any{"type": ch14Locale("float64", "float64"), "values": []any{1.5, 3.25}, "current": 0, "best": 0, "caption": ch14Locale("نداء تاني يستنتج float64.", "Another call infers float64.")},
				map[string]any{"type": ch14Locale("float64", "float64"), "values": []any{1.5, 3.25}, "current": 1, "best": 1, "caption": ch14Locale("جسم واحد يقارن النوعين.", "One body compares both types.")},
			},
			"labels": map[string]any{
				"previous": ch14Locale("اللي فات", "Previous"), "next": ch14Locale("اللي بعده", "Next"), "reset": ch14Locale("من الأول", "Reset"),
				"type": ch14Locale("T المستنتج", "Inferred T"), "values": ch14Locale("القيم", "Values"), "current": ch14Locale("الحالية", "Current"), "best": ch14Locale("best", "best"), "none": ch14Locale("—", "—"),
			},
		}
	case "c142":
		kind = "ch14-stack"
		arabicText, englishText = "جرّب تضيف وتقرا وتشيل آخر اتجاه في الطريق.", "Push, read, and remove the newest route turn."
		config = map[string]any{
			"generic":  ch14Locale("Stack[Direction]", "Stack[Direction]"),
			"examples": []any{map[string]any{"type": ch14Locale("Cell[int]", "Cell[int]"), "value": ch14Locale("12", "12")}, map[string]any{"type": ch14Locale("Cell[string]", "Cell[string]"), "value": ch14Locale("north", "north")}},
			"initial":  []any{ch14Locale("شمال", "North")},
			"arrivals": []any{ch14Locale("شرق", "East"), ch14Locale("جنوب", "South")},
			"labels": map[string]any{
				"generic": ch14Locale("نوع التخزين", "Generic type"), "cells": ch14Locale("قيم توضيحية", "Illustrative values"), "route": ch14Locale("سجل الطريق", "Route history"), "top": ch14Locale("آخر اتجاه", "Newest turn"), "empty": ch14Locale("مفيش اتجاهات", "No turns"),
				"push": ch14Locale("أضف اتجاه", "Push turn"), "peek": ch14Locale("اقرأ الآخر", "Read top"), "pop": ch14Locale("شيل الآخر", "Pop top"), "reset": ch14Locale("من الأول", "Reset"),
				"pushed": ch14Locale("اتضاف اتجاه", "Turn pushed"), "removed": ch14Locale("اتشال اتجاه", "Turn popped"), "resetDone": ch14Locale("رجعنا للبداية", "Back at start"),
			},
		}
	case "c143":
		kind = "ch14-type-constraints"
		arabicText, englishText = "بدّل النوع وشوف الـunderlying type يدخل المجموعة.", "Change types and see which underlying types fit."
		config = map[string]any{
			"constraint": "~int | ~float64",
			"types": []any{
				map[string]any{"name": "Meter", "underlying": "float64", "accepted": true},
				map[string]any{"name": "Count", "underlying": "int", "accepted": true},
				map[string]any{"name": "Route", "underlying": "string", "accepted": false},
			},
			"labels": map[string]any{
				"constraint": ch14Locale("Type set", "Type set"), "type": ch14Locale("النوع", "Type"), "underlying": ch14Locale("النوع الأساسي", "Underlying type"),
				"accepted": ch14Locale("مقبول", "Accepted"), "rejected": ch14Locale("خارج المجموعة", "Outside the set"), "next": ch14Locale("نوع تاني", "Next type"), "reset": ch14Locale("من الأول", "Reset"),
			},
		}
	case "c144":
		kind = "ch14-queue"
		arabicText, englishText = "ضيف طرد، اقرا أول واحد، وبعدها شيله وشوف الترتيب.", "Add a parcel, read the front, then remove it."
		config = map[string]any{
			"elementType": ch14Locale("Queue[string]", "Queue[string]"),
			"initial":     []any{ch14Locale("دواء", "medic"), ch14Locale("أدوات", "tools")},
			"arrivals":    []any{ch14Locale("وقود", "fuel"), ch14Locale("خريطة", "map")},
			"labels": map[string]any{
				"queue": ch14Locale("الطابور", "Queue"), "length": ch14Locale("العدد", "Length"), "front": ch14Locale("الأمام", "Front"), "empty": ch14Locale("فاضي", "empty"), "noRead": ch14Locale("مفيش قيمة للقراءة", "Nothing to read"),
				"add": ch14Locale("وصل طرد", "Arrive"), "read": ch14Locale("اقرأ الأمام", "Read front"), "depart": ch14Locale("شيل الأمام", "Depart"), "reset": ch14Locale("من الأول", "Reset"),
				"added": ch14Locale("اتضاف طرد", "Parcel added"), "departed": ch14Locale("خرج أول طرد", "Front departed"), "resetDone": ch14Locale("الطابور رجع", "Queue restored"),
			},
		}
	case "c145":
		kind = "ch14-recursion"
		arabicText, englishText = "تابع النداءات لحد حالة التوقف، وشوف الـvalues بتتركب وإنت راجع.", "Follow calls to the stopping case, then combine returns."
		config = map[string]any{
			"frames": []any{
				map[string]any{"calls": []string{"total(0)"}, "returned": "—", "caption": ch14Locale("أول نداء يبدأ هنا.", "The first call starts here.")},
				map[string]any{"calls": []string{"total(0)", "total(1)"}, "returned": "—", "caption": ch14Locale("النداء يتقدم لموضع 1.", "The call advances to index 1.")},
				map[string]any{"calls": []string{"total(0)", "total(1)", "total(2)"}, "returned": "—", "caption": ch14Locale("وصلنا لآخر قيمة.", "The last value is reached.")},
				map[string]any{"calls": []string{"total(0)", "total(1)", "total(2)", "total(3)"}, "returned": "0", "caption": ch14Locale("حالة التوقف ترجع صفر.", "The stopping case returns zero.")},
				map[string]any{"calls": []string{"total(0)", "total(1)", "total(2)"}, "returned": "4", "caption": ch14Locale("4 يضاف بعد الرجوع.", "4 is added on return.")},
				map[string]any{"calls": []string{"total(0)", "total(1)"}, "returned": "7", "caption": ch14Locale("3 يركب على 4.", "3 combines with 4.")},
				map[string]any{"calls": []string{"total(0)"}, "returned": "9", "caption": ch14Locale("القيمة رجعت لأول نداء.", "The value returns to the first call.")},
			},
			"labels": map[string]any{
				"previous": ch14Locale("اللي فات", "Previous"), "next": ch14Locale("اللي بعده", "Next"), "reset": ch14Locale("من الأول", "Reset"), "calls": ch14Locale("Call stack", "Call stack"), "returned": ch14Locale("الناتج الراجع", "Returned value"),
			},
		}
	}
	text := englishText
	if arabic {
		text = arabicText
	}
	return LessonBlock{ID: prefix + ".sandbox", Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch14MakeLesson(prefix string, ar, en ch14LessonCopy, idea, try ch14LessonExample, tryAnswer, checkAnswer int) ch14LessonPair {
	makeSteps := func(copy ch14LessonCopy, arabic bool) []LessonStep {
		steps := []LessonStep{
			{ID: prefix + ".mission-step", Title: "Mission", Blocks: []LessonBlock{ch14Markdown(prefix+".mission", copy.Mission)}},
			{ID: prefix + ".idea-step", Title: "Idea", Blocks: []LessonBlock{ch14Walkthrough(prefix+".idea-walkthrough", copy.Idea, idea, arabic)}},
			{ID: prefix + ".try-step", Title: "Try it", Blocks: []LessonBlock{
				ch14Question(prefix+".try-question", copy.TryQuestion+"\n\n~~~go\n"+try.Code+"\n~~~", copy.TryOptions, copy.TryFeedback, tryAnswer),
				ch14Walkthrough(prefix+".try-reveal", copy.TryReveal, try, arabic),
				ch14Sandbox(prefix, arabic),
			}},
			{ID: prefix + ".watch-step", Title: "Watch out", Blocks: []LessonBlock{ch14Markdown(prefix+".watch", copy.Watch)}},
			{ID: prefix + ".check-step", Title: "Think & check", Blocks: []LessonBlock{
				ch14Markdown(prefix+".think", copy.Think),
				ch14Question(prefix+".check-question", copy.Check, copy.CheckOptions, copy.CheckReply, checkAnswer),
			}},
			{ID: prefix + ".turn-step", Title: "Your turn", Blocks: []LessonBlock{ch14Markdown(prefix+".turn", copy.Turn)}},
		}
		if arabic {
			steps[0].Title, steps[1].Title, steps[2].Title = "المهمة", "الفكرة", "جرّب"
			steps[3].Title, steps[4].Title, steps[5].Title = "خد بالك", "فكّر واختبر", "دورك"
		}
		return steps
	}
	return ch14LessonPair{Arabic: makeSteps(ar, true), English: makeSteps(en, false), ArabicHints: ar.Hints, EnglishHints: en.Hints}
}

var ch14Lesson141 = ch14MakeLesson("c141",
	ch14LessonCopy{
		Mission:     "اجمع شغل ~MaxInt~ و~MaxFloat~ في function واحدة، وبعدها ابعت أكبر عمق وأعلى حرارة.\n\nالـstarter فيه بيانات من نوعين؛ فكّر في النوع اللي ينفع يتقارن بـ~>~.\n\nعدّل تعريف الـfunction واستعملها مع الـsliceين.",
		Idea:        "الـgeneric function بتحط type parameter بين أقواس مربعة. ~cmp.Ordered~ constraint بيسمح بالأنواع اللي ينفع تتقارن بـ~>~.\n\nGo بيستنتج ~T~ من الـslice؛ جسم واحد يشتغل مع ~int~ و~float64~. تابع النوع والـvalue في كل نداء.",
		TryQuestion: "توقّع ناتج النداء:",
		TryReveal:   "شوف إزاي نفس الـfunction بتتعامل مع slice أرقام مختلفة النوع.",
		Watch:       "- ~any~ يقبل أنواع أكتر، بس مش بيسمح بالمقارنة تلقائي.\n- القيد لازم يضمن العملية اللي جسم الـfunction بيستخدمها.\n- المثال بيبدأ بأول عنصر؛ بيانات المهمة فيها عناصر بالفعل.",
		Think:       "قبل التشغيل، حدّد النوع اللي هيستنتجه كل نداء. بعدين بص إن كل value بتتبعت مرة واحدة.\n\nخلي بالك إن نتيجة ~Max~ تحتفظ بنفس نوع الـslice، ومش محتاجة تحويل.",
		Check:       "أنهي constraint يسمح باستخدام ~>~ جوه الـfunction؟",
		Turn:        "عدّل الـfunction في الـstarter، وبعدها شغّل البرنامج. تابع العمق والحرارة اللي وصلت للقاعدة، وجرّب كل slice لوحدها لو ظهر error.",
		TryOptions:  []string{"5", "8", "3"}, TryFeedback: []string{"5 أصغر من 8 في المثال.", "8 هي أكبر value في المثال.", "3 مش أكبر من 8."},
		CheckOptions: []string{"~any~", "~cmp.Ordered~", "~[]int~"}, CheckReply: []string{"~any~ مش بيحدد إن المقارنة متاحة.", "القيد ده بيضمن الأنواع القابلة للمقارنة بالترتيب.", "ده نوع slice، مش قيد للـtype parameter."},
		Hints: []string{"~cmp.Ordered~ constraint بيقيد T بالأنواع اللي بتدعم المقارنة بالترتيب.", "انقل جسم الـloop المتكرر لمكان واحد بدل نسختين حسب النوع.", "بص على نداءات ~Max~ مع ~Depths()~ و~Temperatures()~، كل slice لوحدها."},
	},
	ch14LessonCopy{
		Mission:     "Combine ~MaxInt~ and ~MaxFloat~ into one generic function. Report the greatest depth and temperature.\n\nThe starter supplies two slice types; find their shared operation.",
		Idea:        "A generic function puts a type parameter in brackets. ~cmp.Ordered~ permits ordered values.\n\nGo infers ~T~ from each slice, so one body handles both types.",
		TryQuestion: "Predict the call's output:",
		TryReveal:   "One function handles both element types.",
		Watch:       "- ~any~ does not allow comparison.\n- A constraint must permit body operations.\n- Mission slices are nonempty.",
		Think:       "Go infers each call's type. Report both results once.\n\nThe result type matches the input; no conversion is needed.",
		Check:       "Which constraint allows ~>~ inside the function?",
		Turn:        "Edit the function, run it, and inspect both reports. Try each slice separately if a compiler error appears.",
		TryOptions:  []string{"5", "8", "3"}, TryFeedback: []string{"5 is smaller than 8 in the example.", "8 is the greatest value in the example.", "3 is not greater than 8."},
		CheckOptions: []string{"~any~", "~cmp.Ordered~", "~[]int~"}, CheckReply: []string{"~any~ does not guarantee that comparison is available.", "This constraint guarantees ordered comparisons for its types.", "That is a slice type, not a type parameter constraint."},
		Hints: []string{"~cmp.Ordered~ restricts T to types that support ordered comparison.", "Move the repeated loop body into one place instead of keeping type-specific copies.", "Check the ~Max~ calls with ~Depths()~ and ~Temperatures()~ separately."},
	},
	ch14LessonExample{ID: "e14-1-idea", Imports: []string{"cmp", "fmt"}, Code: `func Max[T cmp.Ordered](xs []T) T {
	best := xs[0]
	for _, x := range xs[1:] {
		if x > best {
			best = x
		}
	}
	return best
}
func main() {
	fmt.Println(Max([]int{4, 9, 2}), Max([]float64{1.5, 3.25}))
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "T يقبل أنواع مرتبة.", English: "T accepts ordered types."}, {Line: 2, Arabic: "ابدأ بأول value.", English: "Start from the first value."}, {Line: 4, Arabic: "9 تكسب المقارنة.", English: "9 wins the comparison."}, {Line: 10, Arabic: "main بتنادي النوعين.", English: "main calls both types."}, {Line: 11, Arabic: "أكبر values ظهرت.", English: "Both maximum values print.", Output: "9 3.25"}}},
	ch14LessonExample{ID: "e14-1-try", Imports: []string{"cmp", "fmt"}, Code: `func Max[T cmp.Ordered](xs []T) T {
	best := xs[0]
	for _, x := range xs[1:] {
		if x > best {
			best = x
		}
	}
	return best
}
func main() {
	fmt.Println(Max([]int{8, 3, 5}))
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 2, Arabic: "الـvalue الأولى 8.", English: "The first value is 8."}, {Line: 4, Arabic: "5 أصغر من 8.", English: "5 is smaller than 8."}, {Line: 10, Arabic: "main بتطبع النتيجة.", English: "main prints the result."}, {Line: 11, Arabic: "الناتج هو 8.", English: "The output is 8.", Output: "8"}}},
	1, 1)

var ch14Lesson142 = ch14MakeLesson("c142",
	ch14LessonCopy{
		Mission:     "الخريطة مخفية والنفق متفرّع. استكشف لحد الكريستال، وبعدها ارجع للقاعدة ومعاك الكريستال.\n\nالـstarter فيه ~Stack[T]~ ناقص و~opposite~ جاهزة. سجّل الطريق اللي مشيته.",
		Idea:        "الـgeneric type بيخلي نفس البنية تخزّن أنواع مختلفة. هنا ~Stack[colony.Direction]~ يخزّن اتجاهات، مش أرقام.\n\nالـstack بيرجع آخر حاجة دخلت الأول. ده مناسب للرجوع على نفس الطريق بالعكس، وDFS بيستخدم نفس الفكرة.",
		TryQuestion: "بعد ما تضيف الاتجاهين، أنهي value هتطلع الأول؟",
		TryReveal:   "تابع ترتيب الإضافات والقراءة من آخر القائمة.",
		Watch:       "- اعمل Push للاتجاه بعد حركة ناجحة بس.\n- الخانات المزارة محتاجة سجل لوحده عشان تتفادى الدايرة.\n- Pop من Stack فاضية محتاجة نتيجة توضّح إن مفيش value.",
		Think:       "اسأل نفسك: لو طريق اتسد، أنهي حركة توصّلك للمكان السابق؟\n\nخلي حالة الاستكشاف منفصلة عن نوع البيانات اللي الـstack بيخزنها. الكريستال علامة إن مرحلة الرجوع بدأت.",
		Check:       "بعد Push لـNorth ثم East، أنهي اتجاه Pop بترجّعه؟",
		Turn:        "كمّل سلوك ~Pop~ الأول، وبعدها جرّب الاستكشاف في البرنامج. تابع مكانك بعد كل Move، واتأكد إن طريق الرجوع بيقصر مع كل خطوة.",
		TryOptions:  []string{"North", "East", "South"}, TryFeedback: []string{"North دخل الأول، فيفضل جوه الـstack.", "آخر اتجاه اتضاف هو أول اتجاه يطلع.", "South مش ضمن الـvalues اللي اتضافت."},
		CheckOptions: []string{"North", "East", "South"}, CheckReply: []string{"North أقدم عنصر، فمش هو اللي يخرج الأول.", "الـLIFO بيرجع آخر value دخلت.", "South ما دخلش في الـstack."},
		Hints: []string{"الـtype parameter بيحدد نوع الـvalue اللي هتتخزن في ~Stack~.", "افصل حالة الزيارة عن سجل اتجاهات الرجوع.", "ابدأ من جسم ~Pop~ الفاضي، وبص على نتيجة ~Move~ قبل تسجيل الاتجاه."},
	},
	ch14LessonCopy{
		Mission:     "Explore the hidden branching tunnel, harvest its crystal, and return to base.\n\nThe starter has an unfinished ~Stack[T]~ and a ready ~opposite~ helper.",
		Idea:        "A generic type can hold different types. ~Stack[colony.Direction]~ stores directions.\n\nA stack returns its newest item first, which helps retrace a route; DFS uses this idea.",
		TryQuestion: "After adding both directions, which value comes out first?",
		TryReveal:   "Read the route history from its newest end.",
		Watch:       "- Push only after a successful move.\n- Track visited tiles to avoid loops.\n- Empty Pop needs a presence signal.",
		Think:       "When a route ends, which move returns to the previous tile?\n\nKeep visited state separate from the stack. Finding the crystal starts the return trip.",
		Check:       "After pushing North, then East, which direction does Pop return?",
		Turn:        "Finish ~Pop~, then explore. Watch your position after each move and check that the return route shrinks.",
		TryOptions:  []string{"North", "East", "South"}, TryFeedback: []string{"North arrived first, so it stays deeper in the stack.", "The most recently added direction comes out first.", "South was not added to the stack."},
		CheckOptions: []string{"North", "East", "South"}, CheckReply: []string{"North is older, so it does not leave first.", "LIFO returns the most recent value.", "South was never pushed onto the stack."},
		Hints: []string{"The type parameter chooses which value type the Stack stores.", "Keep visit state separate from the directions used for backtracking.", "Start at the empty body of ~Pop~ and inspect ~Move~ before recording a direction."},
	},
	ch14LessonExample{ID: "e14-2-idea", Imports: []string{"fmt"}, Code: `type Cell[T any] struct { Value T }
func main() {
	depth := Cell[int]{Value: 12}
	turn := Cell[string]{Value: "north"}
	fmt.Println(depth.Value, turn.Value)
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "T بتحدد نوع الـvalue في Cell.", English: "T chooses Cell's value type."}, {Line: 3, Arabic: "الخلية الأولى فيها رقم.", English: "The first cell holds a number."}, {Line: 4, Arabic: "التانية فيها اتجاه كنص.", English: "The next holds a direction as text."}, {Line: 5, Arabic: "اطبع value كل خلية.", English: "Each cell prints its value.", Output: "12 north"}}},
	ch14LessonExample{ID: "e14-2-try", Imports: []string{"fmt"}, Code: `func main() {
	path := []string{"North", "East"}
	next := path[len(path)-1]
	fmt.Println(next)
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 2, Arabic: "North دخلت الأول.", English: "North was pushed first."}, {Line: 3, Arabic: "East دخلت بعديها.", English: "East was pushed next."}, {Line: 4, Arabic: "اقرا آخر value.", English: "Read the last item."}, {Line: 5, Arabic: "East بتطلع الأول.", English: "East comes out first.", Output: "East"}}},
	1, 1)

var ch14Lesson143 = ch14MakeLesson("c143",
	ch14LessonCopy{
		Mission:     "استعمل function ~Sum~ واحدة عشان تجمع الأعماق والأوزان والحرارة، وبعدين ابعت التلات نتائج.\n\nالـstarter فيه type ~Number~ وfunction جمع شغالة مع نوعين؛ بص على الـconstraint مع كل slice.",
		Idea:        "الـconstraint ممكن يحدد type set. علامة ~|~ معناها اختيار، و~+~ مسموح لما كل الأنواع في المجموعة تدعمه.\n\n~Kg~ نوع جديد أساسه ~float64~؛ هو مش النوع نفسه. علامة tilde بتضم الأنواع اللي ليها نفس underlying type.",
		TryQuestion: "توقّع الـvalue اللي هترجع من ~twice(Meter(2.5))~:",
		TryReveal:   "الـconstraint بيسمح بـMeter عشان مبنية على float64.",
		Watch:       "- المقارنة بين نوعين مش معناها إنهم نفس النوع.\n- ~int | float64~ بيستبعد type متعرّف جديد.\n- ~T~ بيحافظ على نوع الـvalue في ناتج الـfunction.",
		Think:       "في كل نداء، طابق النوع الفعلي مع الـtype set. ما تحوّليش البيانات لمجرد إسكات compiler؛ القيد هو اللي يعبّر عن العمليات المقصودة.\n\nالـconstraint هنا خاص بالجمع، مش بالترتيب.",
		Check:       "إيه معنى tilde قبل float64 في type set؟",
		Turn:        "بص على تعريف ~Number~ في الـstarter، وبعدها شغّل نداءات الجمع. تابع رسالة الـcompiler لما نوع الـslice ما يطابقش الـconstraint.",
		TryOptions:  []string{"2.5", "5", "4"}, TryFeedback: []string{"الـvalue بتتجمع مع نفسها، فالناتج 5.", "2.5 هو الإدخال قبل الجمع.", "4 أقل من الناتج المطلوب."},
		CheckOptions: []string{"float64 فقط بالاسم", "أي نوع أساسه float64", "أي نوع قابل للترتيب"}, CheckReply: []string{"النوع المعرّف بشكل منفصل مش داخل المجموعة من غير tilde.", "علامة tilde بتشمل النوع وأي نوع أساسه float64.", "الترتيب خاص بقيود تانية ومش معنى العلامة."},
		Hints: []string{"المشكلة في type ~colony.Kg~ المعرّف، مش في عملية الجمع.", "استعمل type set يطابق underlying type مع العملية المطلوبة.", "قارن تعريف ~Number~ بنداء ~Sum(colony.Weights())~ في الـstarter."},
	},
	ch14LessonCopy{
		Mission:     "Use one ~Sum~ function for depths, weights, and temperatures, then report each total.\n\nThe starter's ~Number~ constraint currently accepts two types.",
		Idea:        "A constraint defines a type set. ~|~ lists choices; ~+~ must work for every member.\n\n~Kg~ is distinct from ~float64~. A tilde term also admits types with that underlying type.",
		TryQuestion: "Predict the value of ~twice(Meter(2.5))~:",
		TryReveal:   "The tilde term admits Meter through its base type.",
		Watch:       "- Shared representation does not mean identical types.\n- ~int | float64~ excludes defined types.\n- ~T~ keeps the result's input type.",
		Think:       "Match each call's type to the set. Do not convert data to silence the compiler.\n\nThis constraint allows addition, not ordering.",
		Check:       "What does a tilde term before float64 mean in a type set?",
		Turn:        "Inspect ~Number~, then run the sum calls. Notice the compiler message for a type outside its set.",
		TryOptions:  []string{"2.5", "5", "4"}, TryFeedback: []string{"The value is added to itself, producing 5.", "2.5 is the input before addition.", "4 is below the result."},
		CheckOptions: []string{"Only the named float64 type", "Any type with float64 underneath", "Any ordered type"}, CheckReply: []string{"A separately defined type is excluded without a tilde term.", "The tilde includes float64 and types based on it.", "Ordering uses different constraints; a tilde does not mean that."},
		Hints: []string{"The issue is the defined type ~colony.Kg~, not addition itself.", "Use a type set that matches the underlying type and required operation.", "Compare the ~Number~ declaration with the starter's ~Sum(colony.Weights())~ call."},
	},
	ch14LessonExample{ID: "e14-3-idea", Imports: []string{"fmt"}, Code: `type Meter float64
type Numeric interface { ~int | ~float64 }
func twice[T Numeric](v T) T { return v + v }
func main() {
	fmt.Println(twice(Meter(1.5)), twice(3))
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "Meter أساسها float64.", English: "Meter is based on float64."}, {Line: 2, Arabic: "المجموعة فيها نوعان أساسيان.", English: "The set allows two underlying types."}, {Line: 3, Arabic: "الجمع مدعوم لكليهما.", English: "Addition is supported by both."}, {Line: 5, Arabic: "النداءان يحتفظان بنوعهما.", English: "Both calls keep their types."}, {Line: 6, Arabic: "النتيجتان اتطبعتا.", English: "Both results print.", Output: "3 6"}}},
	ch14LessonExample{ID: "e14-3-try", Imports: []string{"fmt"}, Code: `type Meter float64
type Numeric interface { ~int | ~float64 }
func twice[T Numeric](v T) T { return v + v }
func main() {
	fmt.Println(twice(Meter(2.5)))
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 1, Arabic: "Meter نوع مبني على float64.", English: "Meter is based on float64."}, {Line: 2, Arabic: "القيد يقبل الأساس ده.", English: "The constraint accepts that base."}, {Line: 3, Arabic: "الـvalue بتتجمع مع نفسها.", English: "The value is added to itself."}, {Line: 5, Arabic: "اطبع النتيجة.", English: "Print the result."}, {Line: 6, Arabic: "الناتج هو 5.", English: "The output is 5.", Output: "5"}}},
	1, 1)

var ch14Lesson144 = ch14MakeLesson("c144",
	ch14LessonCopy{
		Mission:     "المختبر طالب ~Queue[T]~ يحافظ على ترتيب وصول الركاب. كمّل الـAPI في ~main.go~ بحيث يدعم الإضافة، القراءة من الأول، الإزالة، والعدد الحالي.\n\nفكّر في الـzero value ونوع الـvalue قبل ما تبدأ.",
		Idea:        "الـgeneric container بيخزّن نوع يختاره المستخدم، والـnil slice ينفع يبدأ فاضي. الـqueue FIFO: الأقدم يخرج قبل الأحدث.\n\nفي العقد، ~Front~ بيعرض أول value بس، و~Depart~ بيشيلها. وجود bool بيفرق بين غياب الـvalue وvalue بتساوي zero.",
		TryQuestion: "توقّع الـvalue اللي هتظهر:",
		TryReveal:   "النوع الفعلي بيتحدد لما تعمل generic value.",
		Watch:       "- ده عقد متسلسل؛ مفيش goroutines مطلوبة.\n- اقرا أول عنصر كاستعلام من غير ما تغيّر الـqueue.\n- إزالة عنصر لازم ما تسيبش مراجع مستهلكة في التخزين.",
		Think:       "اقرا توقيع كل method واسأل: هل بتقرا ولا بتغيّر؟\n\nلو الـvalue نفسها صفر، وجود العنصر بيتحدد من bool. خلي العدد يعكس العناصر الموجودة دلوقتي، مش كل اللي وصلوا من قبل.",
		Check:       "لو أول value هي zero value، إزاي تعرف إنها موجودة؟",
		Turn:        "ابدأ بربط كل method بوعد واحد من العقد. نفّذ وعد في كل مرة، واكتب بنفسك checks للـqueue الفاضية، والقراءة، والترتيب، وإعادة الاستخدام من غير نسخ مثال جاهز.",
		TryOptions:  []string{"3", "medic", "0"}, TryFeedback: []string{"ده الـvalue اللي اتخزنت جوه Parcel.", "medic مش ضمن المثال ده.", "النوع هنا int والـvalue مش صفر."},
		CheckOptions: []string{"من القيمة فقط", "من bool المصاحبة", "من اسم النوع"}, CheckReply: []string{"الصفر ممكن يبقى عنصر حقيقي.", "الـbool المنفصل بيقول لو فيه عنصر.", "النوع ما يقولش إذا كان العنصر موجود."},
		Hints: []string{"الـgeneric type يخلي ~Queue~ تخزّن أي نوع من غير تغيير الترتيب.", "فرّق بين قراءة العنصر وإزالته، وخلي ~Len~ يعكس العدد الحالي.", "بص على توقيع ~Front~ و~Depart~ وحالة الـzero value في الـstarter."},
	},
	ch14LessonCopy{
		Mission:     "The lab asks for ~Queue[T]~ with FIFO order. Complete the ~main.go~ API for adding, reading, removing, and counting items.\n\nConsider the zero value and element type first.",
		Idea:        "A generic container stores a chosen type; a nil slice can start empty. FIFO removes the oldest item first.\n\n~Front~ reads only, while ~Depart~ removes. The bool distinguishes absence from a zero-valued item.",
		TryQuestion: "Predict the value that prints:",
		TryReveal:   "The concrete type is chosen when the value is created.",
		Watch:       "- No goroutines are needed.\n- Reading the front must not change the queue.\n- Clear consumed references from storage.",
		Think:       "Read each signature: does it read or change state?\n\nThe bool signals presence, even for zero. Len counts current items.",
		Check:       "If the first value is zero, how do you know it exists?",
		Turn:        "Translate each method into a promise. Implement one at a time, then write your own checks for empty behavior, order, and reuse.",
		TryOptions:  []string{"3", "medic", "0"}, TryFeedback: []string{"That is the number stored inside Parcel.", "medic does not appear in this example.", "The int value here is not zero."},
		CheckOptions: []string{"From the value alone", "From the accompanying bool", "From the type name"}, CheckReply: []string{"A zero value can still be a real item.", "The separate signal tells whether an item exists.", "A type does not say whether an item is present."},
		Hints: []string{"A generic type lets Queue hold any element type while preserving its order.", "Separate reading an item from removing it, and let Len show the current count.", "Review the ~Front~ and ~Depart~ signatures and zero-value behavior in the starter."},
	},
	ch14LessonExample{ID: "e14-4-idea", Imports: []string{"fmt"}, Code: `type Parcel[T any] struct { Value T }
func main() {
	box := Parcel[string]{Value: "medic"}
	fmt.Println(box.Value)
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "T نوع الـvalue في الطرد.", English: "T is the parcel's value type."}, {Line: 3, Arabic: "حدد النوع عند الاستخدام.", English: "Choose the type at use."}, {Line: 4, Arabic: "اقرا value الطرد.", English: "Read the parcel value."}, {Line: 5, Arabic: "النص ظهر في الناتج.", English: "The text appears in output.", Output: "medic"}}},
	ch14LessonExample{ID: "e14-4-try", Imports: []string{"fmt"}, Code: `type Parcel[T any] struct { Value T }
func main() {
	box := Parcel[int]{Value: 3}
	fmt.Println(box.Value)
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 1, Arabic: "Parcel generic.", English: "Parcel is generic."}, {Line: 3, Arabic: "النوع هنا int.", English: "The type here is int."}, {Line: 4, Arabic: "اطبع value الطرد.", English: "Print the parcel value."}, {Line: 5, Arabic: "الناتج هو 3.", English: "The output is 3.", Output: "3"}}},
	0, 0)

var ch14Lesson145 = ch14MakeLesson("c145",
	ch14LessonCopy{
		Mission:     "كمّل ~Archive[T]~ عشان تخزّن سجلات بترتيب يحدده comparator. العقد فيه إعداد للترتيب، إدخال، بحث، ونسخة مرتبة من السجلات.\n\nابدأ من الـzero value والـnode الموجودين في الـstarter.",
		Idea:        "الـrecursion بتحل نسخة أصغر من المشكلة. كل نداء محتاج تقدّم وحالة توقف.\n\nفي شجرة، المقارنة بتحدد الاتجاه؛ النوع ~T~ نفسه مش لازم يدعم ~<~. الزيارة in-order تعرض Left ثم العقدة ثم Right حسب المقارنة.",
		TryQuestion: "توقّع ناتج الـfunction:",
		TryReveal:   "كل نداء أصغر بيتبني فوق ناتج النداء اللي بعده.",
		Watch:       "- الـcomparator هو اللي بيحدد ترتيب الـvalues من نوع T.\n- النتيجة صفر معناها الـvalues متساوية حسب الترتيب.\n- اقرا الـarchive الفاضية حتى قبل ما تضبط المقارنة.",
		Think:       "تابع value صغيرة: فين recursion بتقف، وإيه الـvalue اللي بترجّعها؟\n\nافصل قرار المقارنة عن منطق البيانات، وفكّر هل النتيجة المعروضة مستقلة عن البنية الأصلية.",
		Check:       "ليه Archive[T] محتاجة comparator بدل ~record < other~؟",
		Turn:        "قسّم العقد لسلوكيات: اختيار الترتيب، التعامل مع سجل، وبعدها البحث والعرض. اشتغل على function واحدة في كل مرة، وبص إزاي الوصلة بترجع بعد تعديل فرع.",
		TryOptions:  []string{"1", "3", "6"}, TryFeedback: []string{"النداء التاني بيضيف واحد قبل ما يرجع.", "النداء بينزل من 3 لـ0، فالناتج 3.", "الـfunction مش بتجمع values من 1 لـ3."},
		CheckOptions: []string{"T ممكن يكون struct غير مرتب", "كل T أرقام", "الـtree مفيهاش قيم"}, CheckReply: []string{"المقارنة دي بتدعم أنواع مش بتتقارن لوحدها.", "النوع generic مش محصور في الأرقام.", "كل node فيها record."},
		Hints: []string{"ابدأ بتحديد حالة التوقف قبل ما تتابع النداءات recursive.", "استعمل comparator واحد في كل قرارات الترتيب والتساوي.", "بص على حقل ~compare~ وتوقيع ~Records~ في تعريف ~Archive~."},
	},
	ch14LessonCopy{
		Mission:     "Complete ~Archive[T]~ using a supplied comparator. Its API sets order, stores and finds records, then returns an ordered snapshot.\n\nStart with the starter's zero value and node types.",
		Idea:        "Recursion solves smaller problems; each call needs progress and a stopping case.\n\nA comparator chooses tree direction because ~T~ may not support ~<~. In-order visits Left, node, then Right.",
		TryQuestion: "Predict the function's output:",
		TryReveal:   "Each call uses the result from its smaller call.",
		Watch:       "- The comparator defines order for T.\n- Zero means duplicate; keep the first record.\n- Empty reads work, but Store before Order returns an error.",
		Think:       "Trace where recursion stops and what returns next.\n\nKeep comparison separate from data handling. Should displayed records share the tree's storage?",
		Check:       "Why does Archive[T] need a comparator instead of ~record < other~?",
		Turn:        "Split the API into order, store, search, and display. Work one function at a time, tracing how a changed branch returns to its parent.",
		TryOptions:  []string{"1", "3", "6"}, TryFeedback: []string{"The second call adds one before returning.", "The call descends from 3 to 0, so the result is 3.", "The function does not add every value from 1 to 3."},
		CheckOptions: []string{"T may be an unordered struct", "Every T is a number", "The tree holds no values"}, CheckReply: []string{"An injected comparator supports types without built-in order.", "A generic type is not limited to numbers.", "Each node stores a record."},
		Hints: []string{"Identify the stopping case before tracing recursive calls.", "Use one comparator for all ordering and equality decisions.", "Review the ~compare~ field and ~Records~ signature in the ~Archive~ definition."},
	},
	ch14LessonExample{ID: "e14-5-idea", Imports: []string{"fmt"}, Code: `func total(xs []int, i int) int {
	if i == len(xs) {
		return 0
	}
	return xs[i] + total(xs, i+1)
}
func main() {
	fmt.Println(total([]int{2, 3, 4}, 0))
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "أول نداء يبدأ عند 0.", English: "The first call starts at 0."}, {Line: 2, Arabic: "لسه فيه value.", English: "A value remains."}, {Line: 5, Arabic: "النداء بيتقدم خطوة.", English: "The call advances by one."}, {Line: 8, Arabic: "ابدأ الجمع بعد الرجوع.", English: "Addition starts on return."}, {Line: 9, Arabic: "المجموع النهائي ظهر.", English: "The final sum prints.", Output: "9"}}},
	ch14LessonExample{ID: "e14-5-try", Imports: []string{"fmt"}, Code: `func countDown(n int) int {
	if n == 0 {
		return 0
	}
	return 1 + countDown(n-1)
}
func main() {
	fmt.Println(countDown(3))
}`,
		Frames: []ch14LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 1, Arabic: "العدّ يبدأ من 3.", English: "The count starts at 3."}, {Line: 2, Arabic: "الصفر حالة التوقف.", English: "Zero is the stopping case."}, {Line: 5, Arabic: "كل رجوع يضيف واحد.", English: "Each return adds one."}, {Line: 8, Arabic: "اطبع عدد الخطوات.", English: "Print the number of steps."}, {Line: 9, Arabic: "الناتج هو 3.", English: "The output is 3.", Output: "3"}}},
	1, 0)
