package levels

type ch15LessonCopy struct {
	Mission       string
	Idea          string
	TryQuestion   string
	TryReveal     string
	Watch         string
	Think         string
	Check         string
	Turn          string
	TryOptions    []string
	TryFeedback   []string
	CheckOptions  []string
	CheckFeedback []string
	Hints         []string
}

type ch15LessonFrame struct {
	Line    int
	Arabic  string
	English string
	Output  string
}

type ch15LessonExample struct {
	ID      string
	Code    string
	Imports []string
	Frames  []ch15LessonFrame
	Sandbox LessonSandbox
}

type ch15LessonPair struct {
	Arabic       []LessonStep
	English      []LessonStep
	ArabicHints  []string
	EnglishHints []string
}

type ch15SandboxRow struct {
	label  map[string]string
	values []string
}

func ch15SandboxText(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

func ch15SandboxRowFor(ar, en string, values ...string) ch15SandboxRow {
	return ch15SandboxRow{label: ch15SandboxText(ar, en), values: values}
}

func ch15SandboxFrame(ar, en string, rows ...ch15SandboxRow) map[string]any {
	encodedRows := make([]map[string]any, len(rows))
	for i, row := range rows {
		encodedRows[i] = map[string]any{"label": row.label, "values": row.values}
	}
	return map[string]any{"caption": ch15SandboxText(ar, en), "rows": encodedRows}
}

func ch15LessonSandbox(frames ...map[string]any) LessonSandbox {
	return LessonSandbox{
		Kind: "ch15-algorithm-trace",
		Config: map[string]any{
			"labels": map[string]map[string]string{
				"previous": ch15SandboxText("السابق", "Back"),
				"next":     ch15SandboxText("اللي بعده", "Next"),
				"reset":    ch15SandboxText("من الأول", "Reset"),
			},
			"frames": frames,
		},
	}
}

func ch15Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch15Question(id, text string, options, feedback []string, answer int) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch15Walkthrough(id, text string, example ch15LessonExample, arabic bool) LessonBlock {
	frames := make([]LessonFrame, len(example.Frames))
	for i, frame := range example.Frames {
		caption := frame.English
		if arabic {
			caption = frame.Arabic
		}
		frames[i] = LessonFrame{Line: frame.Line, Caption: caption, Output: frame.Output}
	}
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: example.Code, Imports: example.Imports, ExampleID: example.ID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch15MakeLesson(prefix string, ar, en ch15LessonCopy, idea, try ch15LessonExample, tryAnswer, checkAnswer int) ch15LessonPair {
	makeSteps := func(copy ch15LessonCopy, arabic bool) []LessonStep {
		sandboxText := "حرّك خطوة، وراقب تغيّر الحالة."
		if !arabic {
			sandboxText = "Step through and watch the state change."
		}
		steps := []LessonStep{
			{ID: prefix + ".mission-step", Title: "Mission", Blocks: []LessonBlock{ch15Markdown(prefix+".mission", copy.Mission)}},
			{ID: prefix + ".idea-step", Title: "Idea", Blocks: []LessonBlock{ch15Walkthrough(prefix+".idea-walkthrough", copy.Idea, idea, arabic)}},
			{ID: prefix + ".try-step", Title: "Try it", Blocks: []LessonBlock{
				ch15Question(prefix+".try-question", copy.TryQuestion+"\n\n~~~go\n"+try.Code+"\n~~~", copy.TryOptions, copy.TryFeedback, tryAnswer),
				ch15Walkthrough(prefix+".try-reveal", copy.TryReveal, try, arabic),
				{ID: prefix + ".sandbox", Kind: "sandbox", Text: sandboxText, Sandbox: &try.Sandbox},
			}},
			{ID: prefix + ".watch-step", Title: "Watch out", Blocks: []LessonBlock{ch15Markdown(prefix+".watch", copy.Watch)}},
			{ID: prefix + ".check-step", Title: "Think & check", Blocks: []LessonBlock{
				ch15Markdown(prefix+".think", copy.Think),
				ch15Question(prefix+".check-question", copy.Check, copy.CheckOptions, copy.CheckFeedback, checkAnswer),
			}},
			{ID: prefix + ".turn-step", Title: "Your turn", Blocks: []LessonBlock{ch15Markdown(prefix+".turn", copy.Turn)}},
		}
		if arabic {
			steps[0].Title, steps[1].Title, steps[2].Title = "المهمة", "الفكرة", "جرّب"
			steps[3].Title, steps[4].Title, steps[5].Title = "خد بالك", "فكّر واختبر", "دورك"
		}
		return steps
	}
	return ch15LessonPair{Arabic: makeSteps(ar, true), English: makeSteps(en, false), ArabicHints: ar.Hints, EnglishHints: en.Hints}
}

var ch15Lesson151 = ch15MakeLesson("c151",
	ch15LessonCopy{
		Mission:     "وصل للكريستال بأقصر طريق في متاهة الصخور. الخريطة بتتغير، والـstarter بيعرضها من ~colony.Grid()~ وبيحدد الهدف.\n\nخطط الطريق الأول، وبعدها امشِ عليه واحصد الكريستال.",
		Idea:        "BFS بتستخدم queue من نوع FIFO. كل مرة بتوسع الخانات المكتشفة على نفس المسافة، وبعدها بتروح للمسافة اللي بعدها.\n\nلو كل حركة بنفس التكلفة، أول وصول لخانة بيكون بأقل عدد خطوات. خزّن مين وصّلك للخانة عشان تعرف ترجع للطريق.",
		TryQuestion: "توقّع إيه اللي بيتشال وإيه اللي يفضل:",
		TryReveal:   "الـqueue بتحافظ على ترتيب الدخول، وبعدها بتضيف الاكتشاف الجديد.",
		Watch:       "- الصخر ~colony.Rock~ حاجز؛ ما تضفوش كجار صالح.\n- ~grid[x][y]~ معناها X قبل Y.\n- علّم الخانة وقت اكتشافها عشان ما تدخلش الـqueue أكتر من مرة.",
		Think:       "BFS مناسبة لما الهدف أقل عدد حركات وكل حركة بنفس التكلفة. لو التكاليف اختلفت، أول زيارة مش لازم تكون الأرخص.\n\nبعد ما تبني الطريق، اقرا كل خطوة بالنسبة لمكانك الحالي، مش للإحداثيات لوحدها.",
		Check:       "إيه فايدة تسجيل الخانة عند اكتشافها؟",
		Turn:        "ابدأ برسم قائمة الخانات اللي هتفحصها، وبعدها ضيف تسجيل السابق لكل واحدة. شغّل البرنامج على خريطة جديدة وتابع الطريق قبل التنفيذ.",
		TryOptions:  []string{"north و[east]", "east و[north]", "north و[]"}, TryFeedback: []string{"أول value تخرج، وeast تفضل في الـqueue.", "north وصلت الأول، فبتتسحب قبل east.", "east بتفضل بعد سحب north."},
		CheckOptions: []string{"عشان ما نعيدش اكتشافها", "عشان نحرك الدرون", "عشان نرتب الحروف"}, CheckFeedback: []string{"تسجيل الاكتشاف يمنع التكرار في الجولات اللي بعدها.", "الخريطة بتسجل مواقع، مش بتنقل الدرون.", "الترتيب هنا حسب مسافة الوصول، مش الحروف."},
		Hints: []string{"BFS تضمن أقصر طريق لما كل خطوة لها نفس التكلفة.", "استعمل queue للاكتشافات وmap للـparent عشان تسجل الطريق.", "بص على أبعاد ~Grid()~ وشرط الجدار قبل ما تضيف الجيران."},
	},
	ch15LessonCopy{
		Mission:     "Find the crystal by the shortest route through the rock maze. The starter reads the changing map with ~colony.Grid()~ and identifies the goal.",
		Idea:        "BFS uses a FIFO queue and expands tiles layer by layer. With equal move costs, a tile's first visit uses the fewest moves. Record its predecessor to rebuild the route.",
		TryQuestion: "Predict what is removed and what remains:",
		TryReveal:   "A queue keeps arrival order and appends new discoveries.",
		Watch:       "- ~colony.Rock~ blocks passage.\n- ~grid[x][y]~ puts X first.\n- Mark discoveries to prevent repeats.",
		Think:       "BFS finds fewest moves only when costs match. Unequal costs need another method.\n\nFollow route steps from your current position, not coordinates alone.",
		Check:       "Why record a tile as soon as it is discovered?",
		Turn:        "Trace queue entries and their predecessors. Run on a fresh map and inspect the route before following it.",
		TryOptions:  []string{"north and [east]", "east and [north]", "north and []"}, TryFeedback: []string{"The first value leaves, and east remains.", "north arrived first, so it is removed first.", "east remains after north is removed."},
		CheckOptions: []string{"To avoid discovering it again", "To move the drone", "To sort the letters"}, CheckFeedback: []string{"Recording discovery prevents a later duplicate.", "The map records locations; it does not move the drone.", "The order is by travel distance, not letters."},
		Hints: []string{"BFS guarantees a shortest route when each move has equal cost.", "Use a queue for discoveries and a parent map for the route.", "Review ~Grid()~ dimensions and the wall check before adding neighbors."},
	},
	ch15LessonExample{ID: "e15-1-idea", Imports: []string{"fmt"}, Code: `func main() {
	queue := []string{"A", "B"}
	current := queue[0]
	queue = queue[1:]
	queue = append(queue, "C")
	fmt.Println(current, queue)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 2, Arabic: "A وصلت الأول.", English: "A arrived first."}, {Line: 3, Arabic: "اقرا أول عنصر.", English: "Read the first item."}, {Line: 4, Arabic: "شيل A من الأول.", English: "Remove A from the front."}, {Line: 5, Arabic: "C تدخل في الآخر.", English: "C joins at the end."}, {Line: 6, Arabic: "اطبع الحالة الجديدة.", English: "Print the new state."}, {Line: 7, Arabic: "A خرجت وB فضلت.", English: "A left and B remains.", Output: "A [B C]"}}},
	ch15LessonExample{ID: "e15-1-try", Imports: []string{"fmt"}, Code: `func main() {
	queue := []string{"north", "east"}
	first := queue[0]
	queue = queue[1:]
	fmt.Println(first, queue)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 2, Arabic: "north وصلت الأول.", English: "north arrived first."}, {Line: 3, Arabic: "اقرا أول عنصر.", English: "Read the front."}, {Line: 4, Arabic: "شيل أول عنصر.", English: "Slice off the first item."}, {Line: 5, Arabic: "اطبع الباقي.", English: "Show what remains."}, {Line: 6, Arabic: "east لسه في الـqueue.", English: "east is still queued.", Output: "north [east]"}},
		Sandbox: ch15LessonSandbox(
			ch15SandboxFrame("A في أول الطابور.", "A waits at the front.",
				ch15SandboxRowFor("الطابور", "Queue", "A", "B"),
				ch15SandboxRowFor("اتسحب", "Removed", "—")),
			ch15SandboxFrame("اسحب أول عنصر بس.", "Remove just the first item.",
				ch15SandboxRowFor("الطابور", "Queue", "B"),
				ch15SandboxRowFor("اتسحب", "Removed", "A")),
			ch15SandboxFrame("ضيف الاكتشاف الجديد في الآخر.", "Append the new discovery last.",
				ch15SandboxRowFor("الطابور", "Queue", "B", "C"),
				ch15SandboxRowFor("اتسحب", "Removed", "A")),
		),
	},
	0, 0)

var ch15Lesson152 = ch15MakeLesson("c152",
	ch15LessonCopy{
		Mission:     "رتّب أحجام الكريستال من الأصغر للأكبر، من الشمال لليمين. الـstarter قاس الصف كله بالفعل وخزّن الأحجام في ~sizes~.\n\nلازم ترتيب الذاكرة يفضل مطابق للصف بعد كل تبديل.",
		Idea:        "قارن كل زوج جنب بعض، وبدّله لما الشمال يكون أكبر. التبديل يصلّح inversion واحدة في الصف.\n\nفي المهمة، ~colony.Swap~ بتبدّل الكريستال مع الجار من غير ما تحرك الـdrone، وبتكلف tick واحدة حسب توثيق الـAPI وعلى سرعة المهمة العادية.",
		TryQuestion: "توقّع نتيجة تبديل الزوج:",
		TryReveal:   "قارن أول قيمتين، وبعد التبديل اطبع الـslice الجديدة.",
		Watch:       "- حدّث ~sizes~ مع التبديل على الخريطة.\n- حركة الـdrone بتحدد مكانك بالنسبة للزوج المطلوب.\n- loop من غير تبديل معناها إن الترتيب اكتمل.",
		Think:       "التبديل المحلي بيصلّح زوج واحد؛ كرر الفحص لأن حركة واحدة ممكن تكشف inversion تانية.\n\nالدرجة بتعتمد على ticks، ففكّر كمان في طريقة المرور على الصف بدل ما ترجع لنقطة البداية كل مرة.",
		Check:       "بعد تبديل crystalين متجاورين، إيه اللي لازم يتحدث؟",
		Turn:        "شغّل القياس الموجود، وبعدها اكتب خطوة مقارنة وتبديل صغيرة. بعد كل تبديل، طابق السجل مع الصف الفعلي قبل ما توسع الخوارزمية.",
		TryOptions:  []string{"[3 1]", "[1 3]", "[1 1]"}, TryFeedback: []string{"الزوج لسه تنازلي قبل التبديل.", "الـvalues بتتبدل وبتبقى تصاعدية.", "التبديل بيغير المكان مش الـvalue."},
		CheckOptions: []string{"الأرض فقط", "sizes فقط", "الأرض و~sizes~"}, CheckFeedback: []string{"لوحدها تخلي سجل المقاسات قديم.", "الخريطة مش هتعكس الترتيب الجديد.", "لازم الذاكرة والصف يفضلوا متطابقين."},
		Hints: []string{"ابدأ بمقارنة زوج جنب بعض وحدد هل ترتيبهما صحيح.", "خلي قيم ~sizes~ مطابقة لمواضع الكريستال.", "بص على سطر ~colony.Swap~ واتجاه الجار اللي بيبدّله الـAPI."},
	},
	ch15LessonCopy{
		Mission:     "Sort crystals from smallest to largest, left to right. The starter measured the row into ~sizes~.\n\nKeep memory aligned with the crystals after each swap.",
		Idea:        "Swap an adjacent pair when the left value is larger; this fixes one inversion.\n\n~colony.Swap~ exchanges a crystal with its neighbor without moving the drone. The API documents one tick per swap at normal mission speed.",
		TryQuestion: "Predict the result of swapping the pair:",
		TryReveal:   "Swap the first pair when needed, then print the slice.",
		Watch:       "- Update ~sizes~ with each map swap.\n- Position the drone beside the pair.\n- No swaps in a pass means sorted.",
		Think:       "A swap fixes one pair; another inversion may remain.\n\nTicks affect the score, so plan travel along the row.",
		Check:       "After swapping two neighboring crystals, what must also change?",
		Turn:        "Run the measurements and write one compare-and-swap step. Check your record against the row before extending it.",
		TryOptions:  []string{"[3 1]", "[1 3]", "[1 1]"}, TryFeedback: []string{"The pair is still descending before the swap.", "The two values switch into ascending order.", "A swap changes positions; it does not copy a value."},
		CheckOptions: []string{"Only the ground", "Only sizes", "The ground and ~sizes~"}, CheckFeedback: []string{"That alone leaves the size record stale.", "The map would not reflect the new order.", "Memory and the row must stay aligned."},
		Hints: []string{"Start by comparing an adjacent pair and decide if its order is correct.", "Keep ~sizes~ values aligned with their crystal positions.", "Inspect the ~colony.Swap~ line and the neighbor direction that API exchanges."},
	},
	ch15LessonExample{ID: "e15-2-idea", Imports: []string{"fmt"}, Code: `func main() {
	sizes := []int{5, 2, 4}
	if sizes[0] > sizes[1] {
		sizes[0], sizes[1] = sizes[1], sizes[0]
	}
	fmt.Println(sizes)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 2, Arabic: "الـvalues تبدأ 5 و2.", English: "The values start as 5 and 2."}, {Line: 3, Arabic: "الشمال أكبر من اليمين.", English: "The left value is larger."}, {Line: 4, Arabic: "بدّل الـvalues مع بعض.", English: "Swap both values."}, {Line: 5, Arabic: "الزوج بقى مرتب تصاعدي.", English: "The pair is now ascending."}, {Line: 6, Arabic: "اطبع الـslice بعد التعديل.", English: "Print the changed slice."}, {Line: 7, Arabic: "الـvalue التالتة ما اتحركتش.", English: "The third value did not move.", Output: "[2 5 4]"}}},
	ch15LessonExample{ID: "e15-2-try", Imports: []string{"fmt"}, Code: `func main() {
	sizes := []int{3, 1}
	if sizes[0] > sizes[1] {
		sizes[0], sizes[1] = sizes[1], sizes[0]
	}
	fmt.Println(sizes)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 2, Arabic: "3 أكبر من 1.", English: "3 is greater than 1."}, {Line: 3, Arabic: "شرط التبديل اتحقق.", English: "The swap condition is true."}, {Line: 4, Arabic: "بدّل العنصرين.", English: "Swap the two elements."}, {Line: 5, Arabic: "الـslice بقت مرتبة.", English: "The slice is now ordered."}, {Line: 6, Arabic: "اطبع النتيجة.", English: "Print the result."}, {Line: 7, Arabic: "الناتج هو [1 3].", English: "The output is [1 3].", Output: "[1 3]"}},
		Sandbox: ch15LessonSandbox(
			ch15SandboxFrame("قارن زوجًا متجاورًا.", "Compare one adjacent pair.",
				ch15SandboxRowFor("الأحجام", "Sizes", "4", "1", "3"),
				ch15SandboxRowFor("الزوج", "Pair", "4", "1")),
			ch15SandboxFrame("بدّل القيمتين معًا.", "Swap both slice entries.",
				ch15SandboxRowFor("الأحجام", "Sizes", "1", "4", "3"),
				ch15SandboxRowFor("الزوج", "Pair", "1", "4")),
		),
	},
	1, 2)

var ch15Lesson153 = ch15MakeLesson("c153",
	ch15LessonCopy{
		Mission:     "وصل للكريستال بأقل وقت، مش بأقل عدد خطوات. كل خانة Ice بتكلف 4 ticks حسب قاعدة المهمة والـstarter.\n\nالمسار الأقصر ممكن يضيّع وقت لو عدى على جليد كتير.",
		Idea:        "BFS بترتب الوصول بعدد الحركات، لكن Dijkstra بتقارن التكلفة التراكمية. كل طريق بيجمع تكلفة الخانات اللي دخلها.\n\nاختار أقل مسافة معروفة بين الخانات اللي لسه ما خلصتش؛ حدّث مسافة الجار لو لقيت طريق أرخص.",
		TryQuestion: "توقّع كلفة الطريقين بالترتيب:",
		TryReveal:   "الجليد يرفع كلفة الدخول، حتى لو الخطوات أقل.",
		Watch:       "- خزّن تكلفة الطريق كله، مش عدد الخانات بس.\n- حدّث الـparent لما المسافة تتحسن.\n- تكلفة الجليد تخص دخول الخانة؛ استعمل ~cost~ الموجودة في الـstarter.",
		Think:       "قارن مجموع التكلفة، مش طول المسار. نقطة أبعد ممكن تبقى أرخص لو أرضها أسهل.\n\nDijkstra هنا بتتعامل مع تكاليف مش سالبة؛ الخريطة صغيرة، وتقدر تبدأ من غير heap.",
		Check:       "لو أول طريق لخانة كلفته 8، وبعدها لقيت طريق بكلفة 5، أنهي value تفضل؟",
		Turn:        "استعمل الـstarter عشان تقرا نوع الأرض وتكلفة الدخول. جرّب ترسم طريقين على خريطة مختلفة، وبعدها تابع ترتيب اختيار الخانات.",
		TryOptions:  []string{"6 و4", "4 و6", "6 و5"}, TryFeedback: []string{"الطريق الأول بيدخل جليد، والتاني أرضه عادية.", "المطلوب تكلفة العبور الأول وبعدها الالتفاف.", "الأرض العادية مجموعها أربعة مش خمسة."},
		CheckOptions: []string{"8", "5", "الاتنين"}, CheckFeedback: []string{"دي مسافة أقدم وأغلى من المسار الجديد.", "خلي أقل تكلفة معروفة للخانة.", "المسافة مش بتجمع المسارين مع بعض."},
		Hints: []string{"Dijkstra بتختار حسب مجموع التكلفة، مش عدد الحركات.", "خلي أفضل مسافة معروفة ومعاها علامة للخانات اللي خلصت.", "بص على function ~cost~ ومكان الهدف في ~Targets()[0]~."},
	},
	ch15LessonCopy{
		Mission:     "Reach the crystal in the least time. Each Ice tile costs 4 ticks under the mission rule and starter.\n\nA short route may be slower across ice.",
		Idea:        "BFS counts moves; Dijkstra adds entry costs. Choose the lowest known distance among unfinished tiles, and update a neighbor when a cheaper route appears.",
		TryQuestion: "Predict the two route costs, in order:",
		TryReveal:   "Ice can cost more than several ordinary steps.",
		Watch:       "- Track total cost, not tile count.\n- Update the predecessor when cost improves.\n- Use the starter's entry ~cost~ for ice.",
		Think:       "Compare cost, not distance. Easier ground can make a longer route cheaper.\n\nCosts are nonnegative; this small map needs no heap.",
		Check:       "If a tile first costs 8 to reach, then a route costs 5, which value should remain?",
		Turn:        "Read ground types and entry costs from the starter. Trace two routes, then inspect tile-selection order.",
		TryOptions:  []string{"6 and 4", "4 and 6", "6 and 5"}, TryFeedback: []string{"The first route enters ice; the second uses ordinary ground.", "The requested order is crossing, then detour.", "The ordinary route totals four, not five."},
		CheckOptions: []string{"8", "5", "Both"}, CheckFeedback: []string{"That is the older, more expensive route.", "Keep the lowest known cost for the tile.", "Distances do not add both routes together."},
		Hints: []string{"Dijkstra chooses by total cost, not by move count.", "Keep the best known distance and a mark for finished tiles.", "Review the starter's ~cost~ function and target at ~Targets()[0]~."},
	},
	ch15LessonExample{ID: "e15-3-idea", Imports: []string{"fmt"}, Code: `func enterCost(ice bool) int {
	if ice {
		return 4
	}
	return 1
}
func main() {
	through := enterCost(true) + enterCost(false)
	detour := enterCost(false) + enterCost(false) + enterCost(false)
	fmt.Println(through, detour)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "احسب تكلفة دخول خانة.", English: "Calculate one entry cost."}, {Line: 2, Arabic: "الجليد تكلفته 4.", English: "Ice costs 4."}, {Line: 4, Arabic: "الأرض العادية تكلفتها 1.", English: "Ordinary ground costs 1."}, {Line: 7, Arabic: "اجمع تكلفة طريق الجليد.", English: "Add the ice route."}, {Line: 8, Arabic: "اجمع تكلفة الالتفاف.", English: "Add the detour."}, {Line: 9, Arabic: "قارن مجموع الطريقين.", English: "Compare both route totals."}, {Line: 10, Arabic: "الالتفاف أرخص هنا.", English: "The detour costs less here.", Output: "5 3"}}},
	ch15LessonExample{ID: "e15-3-try", Imports: []string{"fmt"}, Code: `func enterCost(ice bool) int {
	if ice {
		return 4
	}
	return 1
}
func main() {
	through := enterCost(true) + enterCost(false) + enterCost(false)
	around := enterCost(false) + enterCost(false) + enterCost(false) + enterCost(false)
	fmt.Println(through, around)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 2, Arabic: "دخول الجليد تكلفته 4.", English: "Entering ice costs 4."}, {Line: 4, Arabic: "الخانة العادية تكلفتها 1.", English: "An ordinary tile costs 1."}, {Line: 7, Arabic: "المسار الأول بيدخل جليد.", English: "The first route enters ice."}, {Line: 8, Arabic: "المسار التاني كله عادي.", English: "The second route is ordinary."}, {Line: 9, Arabic: "اطبع المجموعين.", English: "Print both totals."}, {Line: 10, Arabic: "الناتج 6 و4.", English: "The result is 6 and 4.", Output: "6 4"}},
		Sandbox: ch15LessonSandbox(
			ch15SandboxFrame("المجموعين بيبدأوا من صفر.", "Both totals start at zero.",
				ch15SandboxRowFor("عبر الجليد", "Through ice", "0"),
				ch15SandboxRowFor("حول الجليد", "Around ice", "0")),
			ch15SandboxFrame("ضيف كلفة أول خانة.", "Add the first entry cost.",
				ch15SandboxRowFor("عبر الجليد", "Through ice", "4"),
				ch15SandboxRowFor("حول الجليد", "Around ice", "1")),
			ch15SandboxFrame("ضيف كلفة الخانة اللي بعدها.", "Add the next entry cost.",
				ch15SandboxRowFor("عبر الجليد", "Through ice", "5"),
				ch15SandboxRowFor("حول الجليد", "Around ice", "2")),
			ch15SandboxFrame("قارن الكلفتين التراكميتين.", "Compare the accumulated costs.",
				ch15SandboxRowFor("عبر الجليد", "Through ice", "6"),
				ch15SandboxRowFor("حول الجليد", "Around ice", "3")),
		),
	},
	0, 1)

var ch15Lesson154 = ch15MakeLesson("c154",
	ch15LessonCopy{
		Mission:     "دور على أول index للعمق أو -1 لو مش موجود، وبعدها عدّ الخانات المفتوحة المتصلة بالقاعدة. الـstarter بيبعت النتيجتين.\n\nالقائمة مرتبة، والخريطة فيها صخور وجيوب معزولة.",
		Idea:        "Binary search بتقلل مجال البحث كل خطوة. المجال ~[lo, hi)~؛ لو الـvalue الوسطى أصغر من الهدف، شيلها وكل اللي قبلها.\n\nفي المثال تابع منتصف المجال وحدوده بعد مقارنة واحدة. دي خطوة توضيحية، مش نتيجة البحث بالكامل.",
		TryQuestion: "توقّع عدد الخانات المفتوحة في الخط المتصل:",
		TryReveal:   "الـrecursive function بتقف عند نهاية الـslice أو أول خانة مقفولة.",
		Watch:       "- عند المساواة، كمّل شمال عشان تلاقي أول تكرار.\n- ~Grid~ مفهرسة ~grid[x][y]~، والصخر حاجز.\n- علّم الخانة قبل الجيران، وخلي المدخلات الفاضية تقف بأمان.",
		Think:       "البحث الثنائي بيتعامل مع بيانات مرتبة؛ بص على الحدود قبل ما تقرا index. الـflood بيعد المنطقة من (0,0) بس، مش كل الأرض المفتوحة.\n\nالخانة المعزولة مش جزء من المنطقة المتصلة حتى لو ما كانتش صخر.",
		Check:       "إمتى تعلّم الخانة visited في recursion؟",
		Turn:        "اشتغل على كل function لوحدها. بص على types الـvalues في ~main~، وبعد التشغيل طابق الـindex والعدد مع الخريطة.",
		TryOptions:  []string{"1", "2", "3"}, TryFeedback: []string{"أول خانة مفتوحة تتحسب مرة واحدة.", "الخانتان المفتوحتان المتتاليتان تدخلان في العدد.", "الخانة المقفولة توقف السلسلة."},
		CheckOptions: []string{"بعد زيارة كل الجيران", "قبل استدعاء الجيران", "مش محتاج علامة"}, CheckFeedback: []string{"جار تاني ممكن يرجع للخانة قبل تعليمها.", "العلامة المبكرة تمنع إعادة الزيارة.", "الشبكة فيها دوائر، فبدونها recursion تكرر نفسها."},
		Hints: []string{"الترتيب التصاعدي يسمح بالتخلص من نص مجال البحث كل مرة.", "استعمل مجال half-open للبحث وmap لوحدها للخانات المزارة.", "بص على ترتيب نداء ~depthIndex~ و~flood~ قبل ~colony.Report~ في ~main~."},
	},
	ch15LessonCopy{
		Mission:     "Find the first depth index or -1 when missing, then count open tiles connected to base. The starter reports both results.\n\nDepths are sorted; the map has rocks and isolated pockets.",
		Idea:        "Binary search shrinks ~[lo, hi)~. If the middle is below target, discard it and everything before it.\n\nThe walkthrough shows one comparison.",
		TryQuestion: "Predict the number of open cells in the connected line:",
		TryReveal:   "Recursion stops at the end or a closed cell.",
		Watch:       "- Search left on equality to find the first match.\n- ~Grid~ uses ~grid[x][y]~; rock blocks passage.\n- Empty inputs must stop safely.",
		Think:       "Check bounds before reading. Flood fill counts from (0,0), marking before neighbors.\n\nIsolated cells are outside the connected region.",
		Check:       "When should recursion mark a cell visited?",
		Turn:        "Work on each function separately. Use ~main~ to confirm result types, then compare the index and count after running.",
		TryOptions:  []string{"1", "2", "3"}, TryFeedback: []string{"The first open cell is counted once.", "Both consecutive open cells belong in the count.", "The closed cell ends the run."},
		CheckOptions: []string{"After visiting every neighbor", "Before calling neighbors", "No mark is needed"}, CheckFeedback: []string{"Another neighbor could return before the cell is marked.", "An early mark prevents a repeat visit.", "The grid has cycles, so recursion would repeat without it."},
		Hints: []string{"Ascending order lets each search step discard half the range.", "Use a half-open search interval and a separate visited map for cells.", "Review the call order for ~depthIndex~ and ~flood~ before ~colony.Report~ in ~main~."},
	},
	ch15LessonExample{ID: "e15-4-idea", Imports: []string{"fmt"}, Code: `func main() {
	xs := []int{2, 5, 9}
	lo, hi := 0, len(xs)
	mid := lo + (hi-lo)/2
	fmt.Println(xs[mid], lo, hi)
	if xs[mid] < 6 {
		lo = mid + 1
	}
	fmt.Println(lo, hi)
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 2, Arabic: "المجال يبدأ كاملًا.", English: "The range starts full."}, {Line: 3, Arabic: "اختار المنتصف بأمان.", English: "Choose the midpoint safely."}, {Line: 4, Arabic: "الـvalue الوسطى هي 5.", English: "The middle value is 5."}, {Line: 5, Arabic: "5 أصغر من الهدف.", English: "5 is below the target."}, {Line: 6, Arabic: "حرّك الحد الأدنى.", English: "Move the lower bound."}, {Line: 8, Arabic: "الحدود الجديدة 2 و3.", English: "The new bounds are 2 and 3."}, {Line: 9, Arabic: "اطبع المجال المتبقي.", English: "Show the remaining range.", Output: "5 0 3\n2 3"}}},
	ch15LessonExample{ID: "e15-4-try", Imports: []string{"fmt"}, Code: `func run(open []bool, i int) int {
	if i >= len(open) || !open[i] {
		return 0
	}
	return 1 + run(open, i+1)
}
func main() {
	fmt.Println(run([]bool{true, true, false}, 0))
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 1, Arabic: "ابدأ عند أول خانة.", English: "Begin at the first cell."}, {Line: 2, Arabic: "الخانة مفتوحة.", English: "The cell is open."}, {Line: 5, Arabic: "كل نداء يتقدم.", English: "Each call moves forward."}, {Line: 2, Arabic: "الخانة التالية مقفولة.", English: "The next cell is closed."}, {Line: 8, Arabic: "اجمع عدد الخانات.", English: "Add the cell counts."}, {Line: 9, Arabic: "العدد هو 2.", English: "The count is 2.", Output: "2"}},
		Sandbox: ch15LessonSandbox(
			ch15SandboxFrame("أول نداء يبدأ من الصفر.", "The first call starts at zero.",
				ch15SandboxRowFor("النداءات", "Call stack", "0"),
				ch15SandboxRowFor("العدد", "Count", "—")),
			ch15SandboxFrame("الخلية المفتوحة تنادي بعدها.", "An open cell calls the next.",
				ch15SandboxRowFor("النداءات", "Call stack", "0", "1"),
				ch15SandboxRowFor("العدد", "Count", "—")),
			ch15SandboxFrame("الخلية المقفولة توقف recursion.", "The closed cell stops recursion.",
				ch15SandboxRowFor("النداءات", "Call stack", "0", "1", "2"),
				ch15SandboxRowFor("العدد", "Count", "0")),
			ch15SandboxFrame("كل رجوع يضيف خلية.", "Each return adds one cell.",
				ch15SandboxRowFor("النداءات", "Call stack", "0"),
				ch15SandboxRowFor("العدد", "Count", "2")),
		),
	},
	0, 1)

var ch15Lesson155 = ch15MakeLesson("c155",
	ch15LessonCopy{
		Mission:     "اخدم طلب الإنقاذ الأعلى أولوية حتى مع وصول موجات جديدة. لو الأولوية اتعادلت، الـID الأصغر يسبق.\n\nابدأ بالدفعة الحالية، واستعمل ~Rescue~، وابعث عدد الطلبات اللي اتخدمت مرة واحدة.",
		Idea:        "الـheap بتحافظ على أفضل عنصر عند الجذر حسب ~Less~. الـpackage بتصلّح ترتيبها عند ~heap.Push~ و~heap.Pop~.\n\nالمثال heap لأرقام صغيرة، فالأقل بيخرج الأول. مهمة الإنقاذ محتاجة قاعدة ترتيب مختلفة، مع كسر التعادل بالـID.",
		TryQuestion: "في المثال، إيه الرقم اللي ~heap.Pop~ هترجعه؟",
		TryReveal:   "الـpackage بتختار الجذر حسب ~Less~ قبل ما تخرج العنصر.",
		Watch:       "- نادِ ~heap.Push~ و~heap.Pop~ من الـpackage؛ الـmethods لوحدها ما بتصلحش الترتيب.\n- بص على ~Len~ قبل الإزالة.\n- موجة جديدة ممكن تسبق طلبات قديمة حسب الأولوية.",
		Think:       "المقارنة تحدد ترتيب كل عنصر، والـheap تحافظ عليه مع الإضافة والإزالة. ما تخلطش ترتيب الوصول بأولوية الخدمة.\n\nبعد إنقاذ ناجح، بيانات وصول جديدة تبقى جزءًا من الاختيار التالي.",
		Check:       "مين يقرر أنهي value تخرج الأول من heap؟",
		Turn:        "بص على الـmethods الموجودة، وبعدها ظبّط قاعدة المقارنة حسب عقد الإنقاذ. شغّل البرنامج وتابع كل ~Rescue~ ودفعة الوصول اللي بعدها قبل التقرير النهائي.",
		TryOptions:  []string{"7", "2", "4"}, TryFeedback: []string{"7 مش الأقل في heap المثال.", "~Less~ بترتب الأصغر عشان يبقى أول عنصر.", "4 ما دخلش في heap المثال."},
		CheckOptions: []string{"ترتيب الوصول", "Less", "آخر عنصر في الـslice"}, CheckFeedback: []string{"الـheap بتعيد الترتيب حسب قاعدتها.", "~Less~ بتحدد أنهي عنصر أولى.", "~Pop~ بتشيل آخر الـslice بعد ترتيب الـpackage."},
		Hints: []string{"~Less~ في heap بتحدد أولوية كل عنصر مقارنة بالباقي.", "خلي الـpackage تدير الإضافة والإزالة عشان تحافظ على ترتيب الـheap.", "بص على ~RescueArrivals()~ بعد نجاح ~Rescue~ وقبل الاختيار اللي بعده."},
	},
	ch15LessonCopy{
		Mission:     "Serve the highest-priority request as new waves arrive; smaller ID breaks ties.\n\nStart with current arrivals and report the number served once.",
		Idea:        "A heap keeps its best element at the root by ~Less~. ~heap.Push~ and ~heap.Pop~ repair its order.\n\nThis toy heap puts smaller integers first; rescue needs its own rule and ID tie-break.",
		TryQuestion: "In the example, which number does ~heap.Pop~ return?",
		TryReveal:   "The package selects the root before removal.",
		Watch:       "- Use package ~heap.Push~ and ~heap.Pop~ to repair order.\n- Check Len before removal.\n- New requests can outrank old ones.",
		Think:       "Comparison defines service order; arrival time does not.\n\nAfter a successful rescue, add new arrivals before the next choice.",
		Check:       "What decides which value leaves the heap first?",
		Turn:        "Review the methods and set comparison from the contract. Run the program; watch each rescue and new batch before reporting.",
		TryOptions:  []string{"7", "2", "4"}, TryFeedback: []string{"7 is not the smallest value in this heap example.", "Less orders the smallest value first.", "4 was not added to the heap."},
		CheckOptions: []string{"Arrival order", "Less", "The last slice element"}, CheckFeedback: []string{"The heap reorders items by its rule.", "Less defines which element comes first.", "Pop removes a slice end after the package orders it."},
		Hints: []string{"~Less~ on a heap expresses each element's priority against another.", "Let the package manage insertion and removal so it preserves heap order.", "Check ~RescueArrivals()~ after a successful ~Rescue~ and before the next choice."},
	},
	ch15LessonExample{ID: "e15-5-idea", Imports: []string{"container/heap", "fmt"}, Code: `type Scores []int
func (s Scores) Len() int { return len(s) }
func (s Scores) Less(i, j int) bool { return s[i] < s[j] }
func (s Scores) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s *Scores) Push(v any) { *s = append(*s, v.(int)) }
func (s *Scores) Pop() any {
	old := *s
	last := old[len(old)-1]
	*s = old[:len(old)-1]
	return last
}
func main() {
	q := &Scores{}
	heap.Push(q, 7)
	heap.Push(q, 2)
	fmt.Println(heap.Pop(q).(int))
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "قبل التنفيذ.", English: "Before execution."}, {Line: 1, Arabic: "Scores بتخزّن أرقام.", English: "Scores stores integers."}, {Line: 3, Arabic: "الأصغر أولى.", English: "The smaller value has priority."}, {Line: 5, Arabic: "Push بتغيّر طول الـslice.", English: "Push changes the slice length."}, {Line: 8, Arabic: "Pop بتشيل الطرف.", English: "Pop removes the end."}, {Line: 13, Arabic: "ابدأ بـheap فاضية.", English: "Start with an empty heap."}, {Line: 14, Arabic: "ضيف أول رقم.", English: "Add the first number."}, {Line: 15, Arabic: "ضيف تاني رقم.", English: "Add the second number."}, {Line: 16, Arabic: "الـpackage بتختار الأصغر.", English: "The package selects the smaller value.", Output: "2"}, {Line: 17, Arabic: "الناتج النهائي 2.", English: "The final output is 2.", Output: "2"}}},
	ch15LessonExample{ID: "e15-5-try", Imports: []string{"container/heap", "fmt"}, Code: `type Scores []int
func (s Scores) Len() int { return len(s) }
func (s Scores) Less(i, j int) bool { return s[i] < s[j] }
func (s Scores) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s *Scores) Push(v any) { *s = append(*s, v.(int)) }
func (s *Scores) Pop() any {
	old := *s
	last := old[len(old)-1]
	*s = old[:len(old)-1]
	return last
}
func main() {
	q := &Scores{}
	heap.Push(q, 6)
	heap.Push(q, 3)
	fmt.Println(heap.Pop(q).(int))
}`,
		Frames: []ch15LessonFrame{{Line: 0, Arabic: "ابدأ قبل أول سطر.", English: "Start before the first line."}, {Line: 3, Arabic: "الأصغر بييجي الأول في المثال.", English: "The smaller value comes first."}, {Line: 13, Arabic: "ابدأ بـheap فاضية.", English: "Start with an empty heap."}, {Line: 14, Arabic: "ضيف 6 الأول.", English: "Add 6 first."}, {Line: 15, Arabic: "ضيف 3 بعده.", English: "Add 3 next."}, {Line: 16, Arabic: "اطبع العنصر الأصغر.", English: "Print the smaller item.", Output: "3"}, {Line: 17, Arabic: "الناتج النهائي 3.", English: "The final output is 3.", Output: "3"}},
		Sandbox: ch15LessonSandbox(
			ch15SandboxFrame("الـheap لسه فاضية.", "The heap starts empty.",
				ch15SandboxRowFor("العناصر", "Items", "—"),
				ch15SandboxRowFor("اتشال", "Popped", "—")),
			ch15SandboxFrame("ضيف 7 للـheap.", "Push 7 into the heap.",
				ch15SandboxRowFor("العناصر", "Items", "7"),
				ch15SandboxRowFor("اتشال", "Popped", "—")),
			ch15SandboxFrame("الأصغر بقى في الجذر.", "The smaller item reaches the root.",
				ch15SandboxRowFor("العناصر", "Items", "2", "7"),
				ch15SandboxRowFor("اتشال", "Popped", "—")),
			ch15SandboxFrame("Pop شالت عنصر الجذر.", "Pop removes the root item.",
				ch15SandboxRowFor("العناصر", "Items", "7"),
				ch15SandboxRowFor("اتشال", "Popped", "2")),
		),
	},
	1, 1)
