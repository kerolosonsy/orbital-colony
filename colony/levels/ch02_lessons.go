package levels

import "strings"

func ch02Paragraphs(text string) string {
	var formatted strings.Builder
	inFence, inInline := false, false
	lines := strings.Split(text, "\n")
	for lineIndex, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
			inFence = !inFence
			formatted.WriteString(line)
		} else {
			runes := []rune(line)
			for i := 0; i < len(runes); i++ {
				r := runes[i]
				if r == '~' {
					inInline = !inInline
				}
				formatted.WriteRune(r)
				if !inFence && !inInline && (r == '.' || r == '?' || r == '!') && i+1 < len(runes) && runes[i+1] == ' ' && (i == 0 || runes[i-1] < '0' || runes[i-1] > '9') {
					j := i + 1
					for j < len(runes) && runes[j] == ' ' {
						j++
					}
					if j < len(runes) {
						formatted.WriteString("\n\n")
						i = j - 1
					}
				}
			}
		}
		if lineIndex < len(lines)-1 {
			formatted.WriteByte('\n')
		}
	}
	return formatted.String()
}

func ch02Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: ch02Paragraphs(text)}
}

func ch02Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: ch02Paragraphs(text), Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch02Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: ch02Paragraphs(text), Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch02Variable(name, value string) LessonVariable {
	return LessonVariable{Name: name, Value: value}
}

var ch02InventoryArabic = []LessonStep{
	{ID: "c21.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch02Markdown("c21.mission", "1. اقرأ رصيد الطحالب والكريستال من المخزن.\n2. ابعت مجموعهم للقاعدة.\n3. القيم بتتغير مع كل تشغيل.\n\nالـstarter قرأ الطحالب وطبعها للمراجعة؛ ده مش تقرير للقاعدة.")}},
	{ID: "c21.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch02Walkthrough("c21.variables", "المتغير اسم لقيمة هنستخدمها بعدين. الرمز ~:=~ بيعرّف متغيرًا داخل الدالة. بعد القراءة، تقدر تحسب قيمة جديدة وتطبعها؛ المثال بأرقام للتوضيح.",
		`func main() {
	moss := 3
	crystal := 2
	total := moss + crystal
	fmt.Println(total)
}`,
		"e21-variables", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 2, Caption: "moss تخزن 3.", Variables: []LessonVariable{ch02Variable("moss", "3")}}, {Line: 3, Caption: "crystal تخزن 2.", Variables: []LessonVariable{ch02Variable("moss", "3"), ch02Variable("crystal", "2")}}, {Line: 4, Caption: "اجمع القيمتين.", Variables: []LessonVariable{ch02Variable("moss", "3"), ch02Variable("crystal", "2"), ch02Variable("total", "5")}}, {Line: 5, Caption: "اطبع المجموع.", Variables: []LessonVariable{ch02Variable("total", "5")}, Output: "5"}})}},
	{ID: "c21.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch02Question("c21.predict-total", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tparts := 4\n\textra := 2\n\tfmt.Println(parts + extra)\n}\n~~~", []string{"42", "6", "2"}, 1, []string{"+ بتجمع أعدادًا، مش نصوصًا.", "أربعة زائد اتنين يساوي ستة.", "الناتج مش قيمة المتغير التاني وحده."}),
		ch02Walkthrough("c21.total-reveal", "تابع القيم لحد ما توصل لعملية الجمع.",
			`func main() {
	parts := 4
	extra := 2
	fmt.Println(parts + extra)
}`,
			"e21-total", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 2, Caption: "parts تساوي 4.", Variables: []LessonVariable{ch02Variable("parts", "4")}}, {Line: 3, Caption: "extra تساوي 2.", Variables: []LessonVariable{ch02Variable("parts", "4"), ch02Variable("extra", "2")}}, {Line: 4, Caption: "اجمع واطبع الناتج.", Variables: []LessonVariable{ch02Variable("parts", "4"), ch02Variable("extra", "2")}, Output: "6"}}),
	}},
	{ID: "c21.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch02Markdown("c21.watch", "- Go مش هتقبل متغير اتعرّف وماتستخدمش.\n- ~fmt.Println~ بيطبع ليك؛ ~colony.Report~ هو اللي يبعت للقاعدة.")}},
	{ID: "c21.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch02Markdown("c21.plan", "اقرأ الفئتين من الـAPI، وخزن كل نتيجة باسم واضح. استخدم القيم المخزنة للحساب، وبعدها ابعت الناتج بالطريقة اللي المهمة بتفحصها."),
		ch02Question("c21.report-check", "أنهي أمر يرسل قيمة للمهمة؟", []string{"fmt.Println", "colony.Report", ":="}, 1, []string{"ده يعرض القيمة في Console.", "ده يرسل القيمة للقاعدة.", "ده يعرّف متغيرًا جديدًا."}),
	}},
	{ID: "c21.turn-step", Title: "دورك", Blocks: []LessonBlock{ch02Markdown("c21.turn", "كمّل قراءة المخزن واحسب المجموع. شغّل البرنامج، وراجع إن القاعدة استلمت القيمة المطلوبة.")}},
}

var ch02InventoryEnglish = []LessonStep{
	{ID: "c21.mission-step", Title: "Mission", Blocks: []LessonBlock{ch02Markdown("c21.mission", "1. Read the moss and crystal inventory.\n2. Send their total to base.\n3. The values change on each run.\n\nThe starter reads moss and prints it for inspection; that is not a report to base.")}},
	{ID: "c21.idea-step", Title: "Idea", Blocks: []LessonBlock{ch02Walkthrough("c21.variables", "A variable names a value we will use later. The ~:=~ operator declares a variable inside a function. After reading values, you can calculate a new one and print it; these sample numbers are illustrative.",
		`func main() {
	moss := 3
	crystal := 2
	total := moss + crystal
	fmt.Println(total)
}`,
		"e21-variables", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 2, Caption: "moss stores 3.", Variables: []LessonVariable{ch02Variable("moss", "3")}}, {Line: 3, Caption: "crystal stores 2.", Variables: []LessonVariable{ch02Variable("moss", "3"), ch02Variable("crystal", "2")}}, {Line: 4, Caption: "Add both values.", Variables: []LessonVariable{ch02Variable("moss", "3"), ch02Variable("crystal", "2"), ch02Variable("total", "5")}}, {Line: 5, Caption: "Print the total.", Variables: []LessonVariable{ch02Variable("total", "5")}, Output: "5"}})}},
	{ID: "c21.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch02Question("c21.predict-total", "Predict the output:\n\n~~~go\nfunc main() {\n\tparts := 4\n\textra := 2\n\tfmt.Println(parts + extra)\n}\n~~~", []string{"42", "6", "2"}, 1, []string{"+ adds numbers rather than joining text.", "Four plus two equals six.", "The result is not only the second value."}),
		ch02Walkthrough("c21.total-reveal", "Follow the values until the addition runs.",
			`func main() {
	parts := 4
	extra := 2
	fmt.Println(parts + extra)
}`,
			"e21-total", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 2, Caption: "parts is 4.", Variables: []LessonVariable{ch02Variable("parts", "4")}}, {Line: 3, Caption: "extra is 2.", Variables: []LessonVariable{ch02Variable("parts", "4"), ch02Variable("extra", "2")}}, {Line: 4, Caption: "Add and print the result.", Variables: []LessonVariable{ch02Variable("parts", "4"), ch02Variable("extra", "2")}, Output: "6"}}),
	}},
	{ID: "c21.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch02Markdown("c21.watch", "- Go rejects a variable that is declared but never used.\n- ~fmt.Println~ prints for you; ~colony.Report~ sends to base.")}},
	{ID: "c21.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch02Markdown("c21.plan", "Read both categories from the API and store each result under a clear name. Use the saved values for the calculation, then send it with the operation the mission checks."),
		ch02Question("c21.report-check", "Which call sends a value to the mission?", []string{"fmt.Println", "colony.Report", ":="}, 1, []string{"This displays a value in the Console.", "This sends a value to base.", "This declares a new variable."}),
	}},
	{ID: "c21.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch02Markdown("c21.turn", "Finish reading the inventory and calculate its total. Run the program and check that base received the requested value.")}},
}

var ch02SurveyArabic = []LessonStep{
	{ID: "c22.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch02Markdown("c22.mission", "1. اعرف عرض الخريطة وطولها من ~WorldSize~.\n2. ابعت المساحة، وبعدها المحيط.\n3. الحجم بيتغيّر مع كل تشغيل.\n\nالـstarter بيقرأ البعدين ويطبعهم للمراجعة.")}},
	{ID: "c22.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch02Walkthrough("c22.multiple-values", "دالة Go ممكن ترجع قيمتين، ونستقبلهم بمتغيرين في نفس السطر. كل اسم بياخد القيمة المقابلة له. لو مش محتاج قيمة، استخدم ~_~. المثال يرجّع قياسين للتوضيح ويطبعهم بالترتيب.",
		`func dimensions() (int, int) {
	return 4, 2
}
func main() {
	width, height := dimensions()
	fmt.Println(width, height)
}`,
		"e22-values", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 5, Caption: "اطلب قياسي الدالة."}, {Line: 2, Caption: "الدالة ترجع 4 و2."}, {Line: 5, Caption: "استقبل العرض والطول.", Variables: []LessonVariable{ch02Variable("width", "4"), ch02Variable("height", "2")}}, {Line: 6, Caption: "اطبع القيمتين بالترتيب.", Variables: []LessonVariable{ch02Variable("width", "4"), ch02Variable("height", "2")}, Output: "4 2"}})}},
	{ID: "c22.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch02Question("c22.predict-values", "توقّع الناتج:\n\n~~~go\nfunc pair() (int, int) { return 2, 3 }\nfunc main() {\n\ta, b := pair()\n\tfmt.Println(a, b)\n}\n~~~", []string{"3 2", "2 3", "5"}, 1, []string{"القيم مش بتتبدّل تلقائيًا.", "كل قيمة تروح لمتغيرها بالترتيب.", "الطباعة تعرض قيمتين منفصلتين."}),
		ch02Walkthrough("c22.values-reveal", "تابع القيمة الأولى وبعدها الثانية.",
			`func pair() (int, int) { return 2, 3 }
func main() {
	a, b := pair()
	fmt.Println(a, b)
}`,
			"e22-pair", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 3, Caption: "اطلب قيمتين من pair."}, {Line: 1, Caption: "الدالة ترجع 2 و3."}, {Line: 3, Caption: "استقبل النتيجتين.", Variables: []LessonVariable{ch02Variable("a", "2"), ch02Variable("b", "3")}}, {Line: 4, Caption: "اطبع حسب ترتيب الاستقبال.", Variables: []LessonVariable{ch02Variable("a", "2"), ch02Variable("b", "3")}, Output: "2 3"}}),
	}},
	{ID: "c22.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch02Markdown("c22.watch", "- ~WorldSize~ ترجع العرض قبل الارتفاع.\n- ترتيب التقارير مهم: القاعدة مستنية المساحة قبل المحيط.")}},
	{ID: "c22.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch02Markdown("c22.plan", "سمّي كل بُعد بوضوح، وبعدها استخدم الاتنين في الحسابين. ابعت الناتج الأول، ثم ابعت التاني. راجع إن كل تعبير بيرجع رقمًا صحيحًا."),
		ch02Question("c22.order-check", "إيه التقرير اللي لازم يوصل للقاعدة الأول؟", []string{"المحيط", "المساحة", "عرض الخريطة"}, 1, []string{"القاعدة بتفحص المساحة قبل المحيط.", "المساحة أول تقرير مطلوب.", "العرض لوحده مش التقرير المطلوب."}),
	}},
	{ID: "c22.turn-step", Title: "دورك", Blocks: []LessonBlock{ch02Markdown("c22.turn", "استخدم العرض والارتفاع اللي قرأتهم، واحسب القيم المطلوبة. شغّل البرنامج وتأكد من ترتيب التقارير.")}},
}

var ch02SurveyEnglish = []LessonStep{
	{ID: "c22.mission-step", Title: "Mission", Blocks: []LessonBlock{ch02Markdown("c22.mission", "1. Get the map width and height from ~WorldSize~.\n2. Report area, then perimeter.\n3. The map size changes on each run.\n\nThe starter reads both dimensions and prints them for inspection.")}},
	{ID: "c22.idea-step", Title: "Idea", Blocks: []LessonBlock{ch02Walkthrough("c22.multiple-values", "A Go function can return two values, and we receive them in two variables on one line. Each name gets its matching value. Use ~_~ to ignore one you do not need. This example returns two illustrative measurements and prints them in order.",
		`func dimensions() (int, int) {
	return 4, 2
}
func main() {
	width, height := dimensions()
	fmt.Println(width, height)
}`,
		"e22-values", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 5, Caption: "Request the function values."}, {Line: 2, Caption: "The function returns 4 and 2."}, {Line: 5, Caption: "Receive width and height.", Variables: []LessonVariable{ch02Variable("width", "4"), ch02Variable("height", "2")}}, {Line: 6, Caption: "Print both in order.", Variables: []LessonVariable{ch02Variable("width", "4"), ch02Variable("height", "2")}, Output: "4 2"}})}},
	{ID: "c22.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch02Question("c22.predict-values", "Predict the output:\n\n~~~go\nfunc pair() (int, int) { return 2, 3 }\nfunc main() {\n\ta, b := pair()\n\tfmt.Println(a, b)\n}\n~~~", []string{"3 2", "2 3", "5"}, 1, []string{"The values do not swap themselves.", "Each value goes to its matching variable.", "The print shows two separate values."}),
		ch02Walkthrough("c22.values-reveal", "Follow the first value and then the second.",
			`func pair() (int, int) { return 2, 3 }
func main() {
	a, b := pair()
	fmt.Println(a, b)
}`,
			"e22-pair", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 3, Caption: "Request two values from pair."}, {Line: 1, Caption: "The function returns 2 and 3."}, {Line: 3, Caption: "Receive both results.", Variables: []LessonVariable{ch02Variable("a", "2"), ch02Variable("b", "3")}}, {Line: 4, Caption: "Print in receive order.", Variables: []LessonVariable{ch02Variable("a", "2"), ch02Variable("b", "3")}, Output: "2 3"}}),
	}},
	{ID: "c22.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch02Markdown("c22.watch", "- ~WorldSize~ returns width before height.\n- Report order matters: base expects area before perimeter.")}},
	{ID: "c22.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch02Markdown("c22.plan", "Name each dimension clearly, then use both in the two calculations. Send one result, then the other. Check that each expression produces an integer."),
		ch02Question("c22.order-check", "Which report must reach base first?", []string{"Perimeter", "Area", "Map width"}, 1, []string{"Base checks area before perimeter.", "Area is the first requested report.", "Width alone is not the requested report."}),
	}},
	{ID: "c22.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch02Markdown("c22.turn", "Use the width and height you read, then calculate the requested values. Run the program and check report order.")}},
}

var ch02CallsignsArabic = []LessonStep{
	{ID: "c23.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch02Markdown("c23.mission", "1. اقرأ اسمي النداء، اسمًا في كل مرة.\n2. قارن byte length بعدد الـrunes، وابعت طرفَي كل اسم.\n3. التقارير لازم تمشي بالترتيب المحدد في الـLesson.\n\nالاسمين للتشغيل الحالي؛ متكتبش قيم ثابتة.")}},
	{ID: "c23.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch02Walkthrough("c23.utf8", "~len~ بتعد bytes في UTF-8، مش الحروف. الـrune قيمة Unicode واحدة؛ استخدم دوال ~unicode/utf8~ للعد وفك أول وآخر rune. المثال للتدريب، والنص مش بيانات المهمة.",
		`func main() {
	text := "é🚀"
	first, _ := utf8.DecodeRuneInString(text)
	last, _ := utf8.DecodeLastRuneInString(text)
	fmt.Println(len(text), utf8.RuneCountInString(text), string(first), string(last))
}`,
		"e23-utf8", []string{"fmt", "unicode/utf8"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 2, Caption: "النص فيه رمزان.", Variables: []LessonVariable{ch02Variable("text", "é🚀")}}, {Line: 3, Caption: "اقرأ أول rune.", Variables: []LessonVariable{ch02Variable("first", "é")}}, {Line: 4, Caption: "اقرأ آخر rune.", Variables: []LessonVariable{ch02Variable("first", "é"), ch02Variable("last", "🚀")}}, {Line: 5, Caption: "اطبع العدّ والطرفين.", Variables: []LessonVariable{ch02Variable("first", "é"), ch02Variable("last", "🚀")}, Output: "6 2 é 🚀"}})}},
	{ID: "c23.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch02Question("c23.predict-runes", "توقّع ناتج عدّ الـrunes:\n\n~~~go\nfunc main() {\n\ttext := \"é🚀\"\n\tfmt.Println(utf8.RuneCountInString(text))\n}\n~~~", []string{"6", "2", "8"}, 1, []string{"6 هو عدد الـbytes في UTF-8.", "في النص رمزان Unicode.", "الرقم ده مش مجموع أحجام التخزين."}),
		ch02Walkthrough("c23.rune-reveal", "تابع النص ثم عدّ الـrunes.",
			`func main() {
	text := "é🚀"
	fmt.Println(utf8.RuneCountInString(text))
}`,
			"e23-count", []string{"fmt", "unicode/utf8"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 2, Caption: "النص محفوظ كما هو.", Variables: []LessonVariable{ch02Variable("text", "é🚀")}}, {Line: 3, Caption: "الدالة تطبع عدد runes.", Variables: []LessonVariable{ch02Variable("text", "é🚀")}, Output: "2"}}),
	}},
	{ID: "c23.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch02Markdown("c23.watch", "- ~name[0]~ byte واحدة، مش أول rune كاملة.\n- حوّل الـrune لـ~string~ قبل تقرير الحرف.")}},
	{ID: "c23.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch02Markdown("c23.plan", "ابدأ باسم واحد: اقرأه، احسب النوعين من العد، وفك أول وآخر rune. راجع إن كل تقرير رايح بالترتيب، وبعدها كرر للاسم التاني."),
		ch02Question("c23.byte-check", "إيه اللي بترجعه ~len(name)~؟", []string{"عدد الـbytes", "عدد الـrunes", "آخر rune"}, 0, []string{"len على string بتعد bytes.", "استخدم utf8 للـrune count.", "len مش بتفك الحروف."}),
	}},
	{ID: "c23.turn-step", Title: "دورك", Blocks: []LessonBlock{ch02Markdown("c23.turn", "نفّذ خطوة خطوة وراجع أسماء النداء والناتج في الـConsole. اتأكد إن كل اسم له تقاريره بالترتيب قبل الانتقال للتاني.")}},
}

var ch02CallsignsEnglish = []LessonStep{
	{ID: "c23.mission-step", Title: "Mission", Blocks: []LessonBlock{ch02Markdown("c23.mission", "1. Read each callsign separately.\n2. Compare byte length with rune count, then report both endpoints.\n3. Keep reports in the order described by the Lesson.\n\nThe names belong to this run; do not hardcode values.")}},
	{ID: "c23.idea-step", Title: "Idea", Blocks: []LessonBlock{ch02Walkthrough("c23.utf8", "~len~ counts bytes in UTF-8, not characters. A ~rune~ is one Unicode value; use ~unicode/utf8~ functions to count and decode the first and last rune. This practice text is not mission data.",
		`func main() {
	text := "é🚀"
	first, _ := utf8.DecodeRuneInString(text)
	last, _ := utf8.DecodeLastRuneInString(text)
	fmt.Println(len(text), utf8.RuneCountInString(text), string(first), string(last))
}`,
		"e23-utf8", []string{"fmt", "unicode/utf8"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 2, Caption: "The text has two symbols.", Variables: []LessonVariable{ch02Variable("text", "é🚀")}}, {Line: 3, Caption: "Read the first rune.", Variables: []LessonVariable{ch02Variable("first", "é")}}, {Line: 4, Caption: "Read the last rune.", Variables: []LessonVariable{ch02Variable("first", "é"), ch02Variable("last", "🚀")}}, {Line: 5, Caption: "Print counts and endpoints.", Variables: []LessonVariable{ch02Variable("first", "é"), ch02Variable("last", "🚀")}, Output: "6 2 é 🚀"}})}},
	{ID: "c23.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch02Question("c23.predict-runes", "Predict the rune count:\n\n~~~go\nfunc main() {\n\ttext := \"é🚀\"\n\tfmt.Println(utf8.RuneCountInString(text))\n}\n~~~", []string{"6", "2", "8"}, 1, []string{"6 is the UTF-8 byte count.", "The text has two Unicode values.", "This is not the sum of storage sizes."}),
		ch02Walkthrough("c23.rune-reveal", "Follow the text, then count its runes.",
			`func main() {
	text := "é🚀"
	fmt.Println(utf8.RuneCountInString(text))
}`,
			"e23-count", []string{"fmt", "unicode/utf8"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 2, Caption: "Keep the text intact.", Variables: []LessonVariable{ch02Variable("text", "é🚀")}}, {Line: 3, Caption: "The function prints rune count.", Variables: []LessonVariable{ch02Variable("text", "é🚀")}, Output: "2"}}),
	}},
	{ID: "c23.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch02Markdown("c23.watch", "- ~name[0]~ is one byte, not a whole first rune.\n- Convert a rune to ~string~ before reporting the character.")}},
	{ID: "c23.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch02Markdown("c23.plan", "Start with one name: read it, calculate both counts, and decode its first and last rune. Check report order, then repeat for the second name."),
		ch02Question("c23.byte-check", "What does ~len(name)~ return?", []string{"The byte count", "The rune count", "The last rune"}, 0, []string{"len on a string counts bytes.", "Use utf8 for the rune count.", "len does not decode characters."}),
	}},
	{ID: "c23.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch02Markdown("c23.turn", "Run one step at a time and inspect the callsigns and output in the Console. Confirm each name's reports are ordered before moving to the second.")}},
}

var ch02UnitsArabic = []LessonStep{
	{ID: "c24.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch02Markdown("c24.mission", "1. عرّف نوعًا للشحنات وأنواعًا للوحدة الموحدة.\n2. اقرأ كميات Cell وPack وCrate وحوّلها.\n3. ابعت إجماليًا واحدًا بوحدة cells.\n\nالـstarter بيعرض كمية واحدة بس للمراجعة؛ باقي العقد في وصف المهمة.")}},
	{ID: "c24.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch02Walkthrough("c24.iota", "النوع المسمّى بيدي الرقم معنى واضح. في مجموعة ~const~، ~iota~ تبدأ من صفر وتزيد مع كل سطر. التحويل الصريح يخلّي النوع مناسبًا لدالة أو حساب؛ دي أدوار تدريبية للتوضيح.",
		`type Role int
const (
	Guest Role = iota
	Pilot
	Chief
)
func main() {
	fmt.Println(int(Guest), int(Pilot), int(Chief))
}`,
		"e24-iota", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 2, Caption: "ابدأ مجموعة الثوابت."}, {Line: 3, Caption: "Guest قيمتها صفر."}, {Line: 4, Caption: "الثابت التالي يزيد واحدًا."}, {Line: 5, Caption: "الثابت الثالث يزيد كمان."}, {Line: 8, Caption: "حوّل واطبع القيم.", Output: "0 1 2"}})}},
	{ID: "c24.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch02Question("c24.predict-iota", "توقّع الناتج:\n\n~~~go\ntype Badge int\nconst (\n\tNew Badge = iota\n\tReady\n\tSent\n)\nfunc main() { fmt.Println(int(Ready)) }\n~~~", []string{"0", "1", "2"}, 1, []string{"New بياخد أول قيمة.", "Ready هو الثابت التاني.", "Sent هو الثابت التالت."}),
		ch02Walkthrough("c24.iota-reveal", "تابع بداية العد ومكان ~Ready~.",
			`type Badge int
const (
	New Badge = iota
	Ready
	Sent
)
func main() { fmt.Println(int(Ready)) }`,
			"e24-badge", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 2, Caption: "ابدأ تعريف الثوابت."}, {Line: 3, Caption: "New تبدأ من صفر."}, {Line: 4, Caption: "Ready بعدها بواحد."}, {Line: 7, Caption: "اطبع القيمة بعد التحويل.", Output: "1"}}),
	}},
	{ID: "c24.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch02Markdown("c24.watch", "- الأنواع المسماة مش بتختلط تلقائيًا مع ~int~.\n- وحّد الوحدات قبل جمع كميات الشحنات.")}},
	{ID: "c24.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch02Markdown("c24.plan", "فكّر في العقد على أجزاء: تعريف الأنواع والثوابت، قراءة كل كمية، ثم تحويلها لوحدة واحدة قبل الجمع. خلي تعريف الأنواع برّه ~main~."),
		ch02Question("c24.unit-check", "ليه نحول كل الكميات لنفس الوحدة؟", []string{"عشان المقارنة والجمع يبقوا بمعنى", "عشان iota تقرأ المخزن", "عشان نلغي أنواع Go"}, 0, []string{"الوحدة الموحدة تخلي المجموع مفهومًا.", "iota تخص قيم الثوابت.", "الأنواع تفضل موجودة."}),
	}},
	{ID: "c24.turn-step", Title: "دورك", Blocks: []LessonBlock{ch02Markdown("c24.turn", "ابدأ بتعريفاتك، وبعدها اقرأ البيانات وحوّل كل كمية. شغّل برنامجك وراجع نتيجة المهمة بعد كل تعديل.")}},
}

var ch02UnitsEnglish = []LessonStep{
	{ID: "c24.mission-step", Title: "Mission", Blocks: []LessonBlock{ch02Markdown("c24.mission", "1. Define a cargo kind and types for a shared unit.\n2. Read Cell, Pack, and Crate quantities and convert them.\n3. Send one total measured in cells.\n\nThe starter displays one quantity for inspection; the full contract is in the mission description.")}},
	{ID: "c24.idea-step", Title: "Idea", Blocks: []LessonBlock{ch02Walkthrough("c24.iota", "A named type gives a number a clear meaning. In a ~const~ group, ~iota~ starts at zero and increases on each line. Explicit conversion makes a value fit another type or calculation; these practice roles are illustrative.",
		`type Role int
const (
	Guest Role = iota
	Pilot
	Chief
)
func main() {
	fmt.Println(int(Guest), int(Pilot), int(Chief))
}`,
		"e24-iota", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 2, Caption: "Start the constant group."}, {Line: 3, Caption: "Guest has value zero."}, {Line: 4, Caption: "The next constant adds one."}, {Line: 5, Caption: "The third adds one again."}, {Line: 8, Caption: "Convert and print values.", Output: "0 1 2"}})}},
	{ID: "c24.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch02Question("c24.predict-iota", "Predict the output:\n\n~~~go\ntype Badge int\nconst (\n\tNew Badge = iota\n\tReady\n\tSent\n)\nfunc main() { fmt.Println(int(Ready)) }\n~~~", []string{"0", "1", "2"}, 1, []string{"New gets the first value.", "Ready is the second constant.", "Sent is the third constant."}),
		ch02Walkthrough("c24.iota-reveal", "Follow where counting starts and where ~Ready~ sits.",
			`type Badge int
const (
	New Badge = iota
	Ready
	Sent
)
func main() { fmt.Println(int(Ready)) }`,
			"e24-badge", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 2, Caption: "Begin the constants."}, {Line: 3, Caption: "New starts at zero."}, {Line: 4, Caption: "Ready follows by one."}, {Line: 7, Caption: "Print after conversion.", Output: "1"}}),
	}},
	{ID: "c24.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch02Markdown("c24.watch", "- Named types do not mix with ~int~ automatically.\n- Convert to a shared unit before adding cargo quantities.")}},
	{ID: "c24.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch02Markdown("c24.plan", "Think through the contract in parts: define types and constants, read each quantity, then convert everything to one unit before adding. Keep type declarations outside ~main~."),
		ch02Question("c24.unit-check", "Why convert every quantity to one unit?", []string{"So addition and comparison have meaning", "So iota can read storage", "So Go types disappear"}, 0, []string{"One shared unit makes the total meaningful.", "iota applies to constant values.", "The types remain in the program."}),
	}},
	{ID: "c24.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch02Markdown("c24.turn", "Start with your declarations, then read and convert each quantity. Run your program and inspect the mission result after each adjustment.")}},
}
