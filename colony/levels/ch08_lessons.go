package levels

func ch08Step(id, title string, blocks ...LessonBlock) LessonStep {
	return LessonStep{ID: id, Title: title, Blocks: blocks}
}

func ch08Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch08Frame(line int, caption, output string, variables ...LessonVariable) LessonFrame {
	return LessonFrame{Line: line, Caption: caption, Output: output, Variables: variables}
}

func ch08Variable(name, value string) LessonVariable {
	return LessonVariable{Name: name, Value: value}
}

func ch08Walkthrough(id, text, exampleID, code string, frames ...LessonFrame) LessonBlock {
	return LessonBlock{
		ID: id, Kind: "walkthrough", Text: text, ExampleID: exampleID, Code: code,
		Imports: []string{"fmt"}, Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch08Question(id, text, exampleID, code string, options []string, answer int, feedback []string) LessonBlock {
	block := LessonBlock{ID: id, Kind: "question", Text: text}
	if code != "" {
		block.Code, block.Imports, block.ExampleID = code, []string{"fmt"}, exampleID
	}
	block.Question = &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}
	return block
}

const ch08CopyExample = `x := 3
adjust := func(n int) {
	n = 10
}
adjust(x)
fmt.Println(x)`

const ch08PointerExample = `x := 4
p := &x
fmt.Println(*p)
*p = 9
fmt.Println(x)`

const ch08ValueReceiverExample = `type Gauge struct{ Value int }
func (g Gauge) Add(n int) {
	g.Value += n
}
func main() {
	g := Gauge{Value: 2}
	g.Add(3)
	fmt.Println(g.Value)
}`

const ch08PointerReceiverExample = `type Gauge struct{ Value int }
func (g *Gauge) Add(n int) {
	g.Value += n
}
func main() {
	g := Gauge{Value: 2}
	g.Add(3)
	fmt.Println(g.Value)
}`

const ch08NilGuardExample = `var p *int
if p == nil {
	fmt.Println("none")
	return
}
fmt.Println(*p)`

const ch08NilValueExample = `n := 6
p := &n
if p == nil {
	fmt.Println("none")
} else {
	fmt.Println(*p)
}`

const ch08LinkExample = `type Lamp struct{ Tag string; Next *Lamp }
last := &Lamp{Tag: "blue"}
first := &Lamp{Tag: "red", Next: last}
fmt.Println(first.Next.Tag)`

const ch08WalkLinkExample = `type Lamp struct{ Tag string; Next *Lamp }
tail := &Lamp{Tag: "gold"}
middle := &Lamp{Tag: "green", Next: tail}
head := &Lamp{Tag: "red", Next: middle}
p := head
fmt.Println(p.Tag)
p = p.Next
fmt.Println(p.Tag)`

var ch08CopyArabic = []LessonStep{
	ch08Step("c81.mission-step", "المهمة",
		ch08Markdown("c81.mission", `شغّل الـstarter وقارن كمية الكريستال اللي اتجمعت بالرقم اللي اتطبع.

المفروض كل إضافة تفضل محسوبة بعد رجوع الدالة.`)),
	ch08Step("c81.idea-step", "الفكرة",
		ch08Walkthrough("c81.copy-trace", `Go بيبعت الـargument كنسخة. تغيير parameter نوعه ~int~ ما يغيّرش الأصل.

تابع النسخة، وبعد رجوع الدالة راجع طباعة ~x~.`, "e81-copy", ch08CopyExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(1, "x يبدأ بالرقم 3.", "", ch08Variable("x", "3")),
			ch08Frame(2, "الدالة اتجهزت لاستقبال قيمة.", "", ch08Variable("x", "3")),
			ch08Frame(3, "n نسخة محلية من x.", "", ch08Variable("x", "3"), ch08Variable("n", "3")),
			ch08Frame(4, "التغيير يخص النسخة.", "", ch08Variable("x", "3"), ch08Variable("n", "10")),
			ch08Frame(5, "الدالة رجعت؛ x كما هي.", "", ch08Variable("x", "3")),
			ch08Frame(6, "الخرج: 3.", "3", ch08Variable("x", "3")),
		),
	),
	ch08Step("c81.try-step", "جرّب",
		ch08Question("c81.try-question", "توقّع الخرج بعد التغيير عن طريق p:", "e81-pointer", ch08PointerExample,
			[]string{"4 ثم 9", "4 ثم 4", "9 ثم 4"}, 0,
			[]string{"القراءة الأولى 4، والكتابة عبر العنوان تغيّر x إلى 9.", "الكتابة عبر p توصل للمتغير الأصلي.", "السطر الأول يطبع قبل التغيير؛ الترتيب لا يتبدل."}),
		ch08Walkthrough("c81.pointer-trace", `~p~ شايل عنوان ~x~. القراءة بـ~*p~ تعرض القيمة، والكتابة من خلاله تغيّر الأصل.`, "e81-pointer", ch08PointerExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(1, "x يبدأ بالرقم 4.", "", ch08Variable("x", "4")),
			ch08Frame(2, "p يحمل عنوان x.", "", ch08Variable("x", "4"), ch08Variable("p", "&x")),
			ch08Frame(3, "قراءة العنوان تعرض 4.", "4", ch08Variable("x", "4"), ch08Variable("p", "&x")),
			ch08Frame(4, "الكتابة عبر p تغيّر x.", "4", ch08Variable("x", "9"), ch08Variable("p", "&x")),
			ch08Frame(5, "الخرج النهائي: 4 ثم 9.", "4\n9", ch08Variable("x", "9")),
		),
	),
	ch08Step("c81.watch-step", "خد بالك",
		ch08Markdown("c81.watch", `~int~ قيمة، و~*int~ نوع pointer. ~&x~ يجيب العنوان؛ و~*p~ يقرأ القيمة أو يكتب فيها.

	التعديل على نسخة محلية يختفي بعد رجوع الدالة. حدّد هل المطلوب تغيير الأصل أو حساب قيمة جديدة.`)),
	ch08Step("c81.check-step", "فكّر واختبر",
		ch08Markdown("c81.plan", `امشِ على مراحل: حدّد المتغير اللي يحتفظ بالتغيير، وافحص parameter، وبعدها راجع مكان النداء.

	اختبر بأرقام صغيرة الأول.`),
		ch08Question("c81.check-question", "لو دالة غيّرت parameter نوعه ~int~، هل يتغير متغير المستدعي؟",
			"", "", []string{"دائمًا", "لا، دي نسخة", "بس لو القيمة موجبة"}, 1,
			[]string{"نوع int وحده لا يربط التغيير بمتغير المستدعي.", "صحيح: parameter بياخد نسخة من القيمة.", "إشارة الرقم لا تغيّر قاعدة النسخ."}),
	),
	ch08Step("c81.turn-step", "دورك",
		ch08Markdown("c81.turn", `ابدأ من دالة الـstarter. غيّر حاجة واحدة، شغّل البرنامج، وقارن التقرير بجمع الدرون.`)),
}

var ch08CopyEnglish = []LessonStep{
	ch08Step("c81.mission-step", "Mission",
		ch08Markdown("c81.mission", `Run the starter and compare collected crystal with the printed total.

Each addition should still count after the function returns.`)),
	ch08Step("c81.idea-step", "Idea",
		ch08Walkthrough("c81.copy-trace", `Go passes an argument as a copy. Changing an ~int~ parameter does not change the original.

Follow the copy, then check ~x~ after the function returns.`, "e81-copy", ch08CopyExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(1, "x starts at 3.", "", ch08Variable("x", "3")),
			ch08Frame(2, "Function is ready.", "", ch08Variable("x", "3")),
			ch08Frame(3, "n copies x.", "", ch08Variable("x", "3"), ch08Variable("n", "3")),
			ch08Frame(4, "Only n changes.", "", ch08Variable("x", "3"), ch08Variable("n", "10")),
			ch08Frame(5, "x stays 3 after return.", "", ch08Variable("x", "3")),
			ch08Frame(6, "Output: 3.", "3", ch08Variable("x", "3")),
		),
	),
	ch08Step("c81.try-step", "Try it",
		ch08Question("c81.try-question", "Predict both prints:", "e81-pointer", ch08PointerExample,
			[]string{"4 then 9", "4 then 4", "9 then 4"}, 0,
			[]string{"Writing through p changes x to 9.", "p reaches the original.", "The first print comes before the change."}),
		ch08Walkthrough("c81.pointer-trace", `p holds x's address. ~*p~ reads or changes the value.`, "e81-pointer", ch08PointerExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(1, "x starts at 4.", "", ch08Variable("x", "4")),
			ch08Frame(2, "p holds x's address.", "", ch08Variable("x", "4"), ch08Variable("p", "&x")),
			ch08Frame(3, "Read 4 through p.", "4", ch08Variable("x", "4"), ch08Variable("p", "&x")),
			ch08Frame(4, "Write 9 through p.", "4", ch08Variable("x", "9"), ch08Variable("p", "&x")),
			ch08Frame(5, "Prints: 4 then 9.", "4\n9", ch08Variable("x", "9")),
		),
	),
	ch08Step("c81.watch-step", "Watch out",
		ch08Markdown("c81.watch", `~int~ is a value; ~*int~ is a pointer type. ~&x~ gets an address; ~*p~ reads or writes there.

A local copy disappears on return.

Decide whether you need to change the original or calculate a new value.`)),
	ch08Step("c81.check-step", "Think & check",
		ch08Markdown("c81.plan", `Work in phases: find the variable to preserve, inspect the parameter, then check the call site.

Try small numbers first.`),
		ch08Question("c81.check-question", "If a function changes an ~int~ parameter, does the caller's variable change?",
			"", "", []string{"Always", "No, it is a copy", "Only for positive values"}, 1,
			[]string{"An int alone does not connect changes to the caller.", "Right: a parameter receives a copy.", "A number's sign does not change copying."}),
	),
	ch08Step("c81.turn-step", "Your turn",
		ch08Markdown("c81.turn", `Start with the starter's function. Change one thing and run it.

Compare the report with collected crystal.`)),
}

var ch08CargoArabic = []LessonStep{
	ch08Step("c82.mission-step", "المهمة",
		ch08Markdown("c82.mission", `الدرون يشيل لحد 10 كريستالات في الشحنة، والمخزن محتاج 30.

الـstarter فيه struct وmethods ومسار جاهز، بس عدّاد الشحنة مش بيحتفظ بتغييره.

اقرأ اللي بيحصل بعد أول تحميل، وبعدها راقب التفريغ عند القاعدة.`)),
	ch08Step("c82.idea-step", "الفكرة",
		ch08Walkthrough("c82.value-trace", `الـvalue receiver بياخد نسخة من الـstruct. تغيير الحقل في النسخة ما يغيّرش قيمة ~g~ اللي عند المستدعي.

في المثال، method بترفع Gauge محلية؛ تابع القيمة بعد الرجوع.`, "e82-value", ch08ValueReceiverExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(6, "Gauge تبدأ بقيمة 2.", "", ch08Variable("g.Value", "2")),
			ch08Frame(7, "استدعاء Add يعمل نسخة receiver.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver.Value", "2")),
			ch08Frame(3, "التعديل يقع على النسخة.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver.Value", "5")),
			ch08Frame(4, "نسخة receiver انتهت.", "", ch08Variable("g.Value", "2")),
			ch08Frame(8, "الخرج يفضل 2.", "2", ch08Variable("g.Value", "2")),
		),
	),
	ch08Step("c82.try-step", "جرّب",
		ch08Question("c82.try-question", "الـmethod هنا ليها pointer receiver؛ إيه الخرج؟", "e82-pointer", ch08PointerReceiverExample,
			[]string{"2", "5", "0"}, 1,
			[]string{"دي كانت نتيجة receiver بياخد نسخة عادية.", "صح: الـreceiver بيشاور على g، فتظهر الزيادة.", "قيمة Gauge اتحددت بـ2 قبل نداء method."}),
		ch08Walkthrough("c82.pointer-trace", `الـpointer receiver يغيّر الـstruct الأصلية. لما تنادي method على متغير عادي، Go يمرّر عنوانه تلقائيًا.`, "e82-pointer", ch08PointerReceiverExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(6, "g تبدأ بقيمة 2.", "", ch08Variable("g.Value", "2")),
			ch08Frame(7, "النداء يمرّر عنوان g.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver", "&g")),
			ch08Frame(2, "الـreceiver يوصل للأصل.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver.Value", "2")),
			ch08Frame(3, "الحقل الأصلي يزيد.", "", ch08Variable("g.Value", "5"), ch08Variable("receiver.Value", "5")),
			ch08Frame(8, "الخرج: 5.", "5", ch08Variable("g.Value", "5")),
		),
	),
	ch08Step("c82.watch-step", "خد بالك",
		ch08Markdown("c82.watch", `~func (g Gauge)~ تستقبل نسخة؛ ~func (g *Gauge)~ تستقبل عنوان الـstruct. ده فرق في تعريف receiver.

~colony.Unload()~ تنقل الشحنة للمخزن من القاعدة فقط. والحصاد يفشل لو الكمية مش هتساع في الشحنة.`)),
	ch08Step("c82.check-step", "فكّر واختبر",
		ch08Markdown("c82.plan", `رتّب التفكير كده: راقب حالة Cargo، وافحص method اللي بتغيّرها، وبعدها راجع شرط المساحة ومكان التفريغ.

اختبر تغييراتك مع شحنة صغيرة قبل المسار كله.`),
		ch08Question("c82.check-question", "أنهي receiver يخلّي method تعدّل struct المستدعي؟",
			"", "", []string{"Value receiver", "Pointer receiver", "الاتنين دايمًا"}, 1,
			[]string{"value receiver يعدّل نسخته المحلية.", "صحيح: pointer receiver يقدر يوصل للأصل.", "النتيجة تعتمد على نوع receiver، مش على اسم method."}),
	),
	ch08Step("c82.turn-step", "دورك",
		ch08Markdown("c82.turn", `ارجع للـmethod اللي بتسجّل الشحنة. شغّل الـstarter، عدّلها، وبعدين تابع عدّاد Cargo ورسائل المخزن.`)),
}

var ch08CargoEnglish = []LessonStep{
	ch08Step("c82.mission-step", "Mission",
		ch08Markdown("c82.mission", `The drone carries at most 10 crystal per load; storage needs 30.

The starter's cargo counter does not keep its change.

Run it and watch the first unload.`)),
	ch08Step("c82.idea-step", "Idea",
		ch08Walkthrough("c82.value-trace", `A value receiver gets a struct copy. Changing it does not change the caller's ~g~.

Follow the Gauge value after the method returns.`, "e82-value", ch08ValueReceiverExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(6, "Gauge starts at 2.", "", ch08Variable("g.Value", "2")),
			ch08Frame(7, "Add receives a receiver copy.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver.Value", "2")),
			ch08Frame(3, "The copy is changed.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver.Value", "5")),
			ch08Frame(4, "The receiver copy ends.", "", ch08Variable("g.Value", "2")),
			ch08Frame(8, "Output stays 2.", "2", ch08Variable("g.Value", "2")),
		),
	),
	ch08Step("c82.try-step", "Try it",
		ch08Question("c82.try-question", "This method has a pointer receiver; what is printed?", "e82-pointer", ch08PointerReceiverExample,
			[]string{"2", "5", "0"}, 1,
			[]string{"That was the result with an ordinary value receiver.", "Right: the receiver points to g, so the increase remains.", "Gauge was set to 2 before the method call."}),
		ch08Walkthrough("c82.pointer-trace", `A pointer receiver changes the original.

Go passes an ordinary variable's address when calling its method.`, "e82-pointer", ch08PointerReceiverExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(6, "g starts at 2.", "", ch08Variable("g.Value", "2")),
			ch08Frame(7, "The call passes g's address.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver", "&g")),
			ch08Frame(2, "The receiver reaches the original.", "", ch08Variable("g.Value", "2"), ch08Variable("receiver.Value", "2")),
			ch08Frame(3, "The original field increases.", "", ch08Variable("g.Value", "5"), ch08Variable("receiver.Value", "5")),
			ch08Frame(8, "Output: 5.", "5", ch08Variable("g.Value", "5")),
		),
	),
	ch08Step("c82.watch-step", "Watch out",
		ch08Markdown("c82.watch", `~func (g Gauge)~ receives a copy; ~func (g *Gauge)~ receives its address.


~colony.Unload()~ transfers cargo at base only. Harvest fails when it will not fit.`)),
	ch08Step("c82.check-step", "Think & check",
		ch08Markdown("c82.plan", `Think in phases: inspect Cargo, its update method, the space check, and unload timing.

Try a small load first.`),
		ch08Question("c82.check-question", "Which receiver lets a method change the caller's struct?",
			"", "", []string{"Value receiver", "Pointer receiver", "Both always"}, 1,
			[]string{"A value receiver changes its copy.", "Right: a pointer receiver reaches the original.", "The receiver type determines the result."}),
	),
	ch08Step("c82.turn-step", "Your turn",
		ch08Markdown("c82.turn", `Return to the cargo method. Run the starter, change it, then watch Cargo and storage messages.`)),
}

var ch08RadarArabic = []LessonStep{
	ch08Step("c83.mission-step", "المهمة",
		ch08Markdown("c83.mission", `الرادار بيرجّع عنوان أقرب هدف لسه موجود. المطلوب إن الدرون يزور الأهداف ويحصدها كلها.

الـstarter بيقرأ النتيجة مرة واحدة؛ فكر إيه اللي يحصل لما تخلص الأهداف.`)),
	ch08Step("c83.idea-step", "الفكرة",
		ch08Walkthrough("c83.nil-trace", `الـpointer ممكن يكون عنوان قيمة، أو ~nil~ لو مفيش قيمة. افحصه قبل ما تقرأ اللي بيشاور عليه.

هنا ~p~ فاضي؛ الشرط يمنع قراءة عنوان غير موجود.`, "e83-nil", ch08NilGuardExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(1, "p يبدأ nil.", "", ch08Variable("p", "nil")),
			ch08Frame(2, "الشرط يكتشف nil.", "", ch08Variable("p", "nil")),
			ch08Frame(3, "نطبع إن مفيش هدف.", "none", ch08Variable("p", "nil")),
			ch08Frame(4, "return ينهي المثال بأمان.", "none", ch08Variable("p", "nil")),
		),
	),
	ch08Step("c83.try-step", "جرّب",
		ch08Question("c83.try-question", "p عنوانه صالح هنا؛ إيه اللي هيتطبع؟", "e83-value", ch08NilValueExample,
			[]string{"none", "6", "panic"}, 1,
			[]string{"الفرع ده محتاج p تكون nil.", "صح: p مش nil، فبنقرأ القيمة اللي بتشاور عليها.", "الـdereference آمن بعد فحص المؤشر."}),
		ch08Walkthrough("c83.value-trace", `فحص ~p == nil~ يختار مسارًا آمنًا. ولما المؤشر صالح، تقدر تقرأ القيمة اللي وراه.`, "e83-value", ch08NilValueExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(1, "n يساوي 6.", "", ch08Variable("n", "6")),
			ch08Frame(2, "p يشاور على n.", "", ch08Variable("n", "6"), ch08Variable("p", "&n")),
			ch08Frame(3, "المؤشر مش nil.", "", ch08Variable("p", "&n")),
			ch08Frame(6, "نقرأ القيمة 6.", "6", ch08Variable("n", "6"), ch08Variable("p", "&n")),
		),
	),
	ch08Step("c83.watch-step", "خد بالك",
		ch08Markdown("c83.watch", `~colony.Nearest()~ بيرجّع ~*colony.Point~، أو ~nil~ لما مفيش أهداف. متقرأش ~p.X~ قبل فحص ~p~.

كل سؤال للرادار بيكلف tick. كل نداء جديد بيبحث عن أقرب هدف باقي وقتها.`)),
	ch08Step("c83.check-step", "فكّر واختبر",
		ch08Markdown("c83.plan", `رتّب المسار: اسأل الرادار، افحص إن فيه نتيجة، وبعدها استخدم النقطة. كرر التفكير بعد إزالة هدف.

ما تفترضش ترتيبًا ثابتًا بين الأهداف؛ الرادار يحدد الأقرب حسب موقع الدرون الحالي.`),
		ch08Question("c83.check-question", "إيه معنى ~nil~ من ~colony.Nearest()~؟",
			"", "", []string{"مفيش هدف متبقّي", "الهدف عند 0,0", "البرنامج لسه ما اتحركش"}, 0,
			[]string{"صحيح: الرادار لم يعد يجد هدفًا متبقيًا.", "nil لا يحدد موقعًا؛ هو غياب النتيجة.", "النتيجة لا تخبرك هل الدرون تحرك أم لا."}),
	),
	ch08Step("c83.turn-step", "دورك",
		ch08Markdown("c83.turn", `راجع طباعة الـstarter: هي بتفك المؤشر قبل ما تعرف حالته. جرّب حالة وجود هدف وحالة انتهاء الأهداف.`)),
}

var ch08RadarEnglish = []LessonStep{
	ch08Step("c83.mission-step", "Mission",
		ch08Markdown("c83.mission", `The radar returns the nearest remaining target as a pointer. Harvest every target.

The starter reads once. What if none remain?`)),
	ch08Step("c83.idea-step", "Idea",
		ch08Walkthrough("c83.nil-trace", `A pointer may hold an address or be ~nil~. Check first.

The guard skips a missing value.`, "e83-nil", ch08NilGuardExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(1, "p starts as nil.", "", ch08Variable("p", "nil")),
			ch08Frame(2, "The check detects nil.", "", ch08Variable("p", "nil")),
			ch08Frame(3, "Print that no target exists.", "none", ch08Variable("p", "nil")),
			ch08Frame(4, "return ends safely.", "none", ch08Variable("p", "nil")),
		),
	),
	ch08Step("c83.try-step", "Try it",
		ch08Question("c83.try-question", "Predict the output when p is valid:", "e83-value", ch08NilValueExample,
			[]string{"none", "6", "panic"}, 1,
			[]string{"That branch runs only when p is nil.", "Right: p points to n, whose value is 6.", "The checked dereference is safe."}),
		ch08Walkthrough("c83.value-trace", `Check ~p~ first. When valid, it holds a readable value.`, "e83-value", ch08NilValueExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(1, "n is 6.", "", ch08Variable("n", "6")),
			ch08Frame(2, "p points to n.", "", ch08Variable("n", "6"), ch08Variable("p", "&n")),
			ch08Frame(3, "The pointer is not nil.", "", ch08Variable("p", "&n")),
			ch08Frame(6, "Read the value 6.", "6", ch08Variable("n", "6"), ch08Variable("p", "&n")),
		),
	),
	ch08Step("c83.watch-step", "Watch out",
		ch08Markdown("c83.watch", `~colony.Nearest()~ returns a Point pointer or ~nil~ when targets end. Check ~p~ before reading ~p.X~.

Each radar query costs one tick; the next scan finds the nearest remaining target.`)),
	ch08Step("c83.check-step", "Think & check",
		ch08Markdown("c83.plan", `Phases: scan, check, use the point, then repeat after harvest.

	Nearest depends on the drone's position; do not assume an order.`),
		ch08Question("c83.check-question", "What does ~nil~ from ~colony.Nearest()~ mean?",
			"", "", []string{"No target remains", "The target is at 0,0", "The program has not moved"}, 0,
			[]string{"Right: no target remains.", "nil means no result, not a location.", "It says nothing about movement."}),
	),
	ch08Step("c83.turn-step", "Your turn",
		ch08Markdown("c83.turn", `Check the starter's dereference. Try with targets, then after they run out.`)),
}

var ch08CartArabic = []LessonStep{
	ch08Step("c84.mission-step", "المهمة",
		ch08Markdown("c84.mission", `كل عربة فيها tag ووصلة للعربة اللي بعدها. المطلوب توصيل عربات وفصلها من غير ما باقي القطر يختفي.

الـstarter معرّف الأنواع، لكن methods لسه فاضية. اختبر العقد على سلسلة صغيرة.`)),
	ch08Step("c84.idea-step", "الفكرة",
		ch08Walkthrough("c84.link-trace", `الـstruct تقدر تخزن pointer لنوعها نفسه. هنا ~Next~ بيربط عربتين، و~nil~ معناها إن الوصلة انتهت.

اقرأ tag اللي ورا أول وصلة.`, "e84-link", ch08LinkExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(2, "آخر لمبة اسمها blue.", "", ch08Variable("last.Tag", "blue"), ch08Variable("last.Next", "nil")),
			ch08Frame(3, "first تتصل بآخر لمبة.", "", ch08Variable("first.Tag", "red"), ch08Variable("first.Next", "&last")),
			ch08Frame(4, "نقرأ tag بعد الوصلة.", "blue", ch08Variable("first.Next.Tag", "blue")),
		),
	),
	ch08Step("c84.try-step", "جرّب",
		ch08Question("c84.try-question", "المؤشر اتحرك وصلة واحدة؛ إيه الخرج؟", "e84-walk", ch08WalkLinkExample,
			[]string{"red ثم green", "green ثم gold", "red ثم gold"}, 0,
			[]string{"صحيح: p يبدأ عند head، وبعد الوصلة يصل للوسط.", "ده يبدأ من الوسط ويتحرك للآخر.", "الوصلة الواحدة توصل للوسط قبل الطرف."}),
		ch08Walkthrough("c84.walk-trace", `المتغير ~p~ مؤشر مؤقت؛ تحريكه للسطر التالي ما يغيّرش ~head~. تابع tag قبل الحركة وبعدها.`, "e84-walk", ch08WalkLinkExample,
			ch08Frame(0, "قبل التنفيذ.", ""),
			ch08Frame(2, "tail هو gold.", "", ch08Variable("tail.Tag", "gold")),
			ch08Frame(3, "middle يشاور على tail.", "", ch08Variable("middle.Tag", "green"), ch08Variable("middle.Next", "&tail")),
			ch08Frame(4, "head يشاور على الوسط.", "", ch08Variable("head.Tag", "red"), ch08Variable("head.Next", "&middle")),
			ch08Frame(5, "p يبدأ عند head.", "", ch08Variable("p.Tag", "red")),
			ch08Frame(6, "أول tag هو red.", "red", ch08Variable("p.Tag", "red")),
			ch08Frame(7, "p يتبع Next للوسط.", "red", ch08Variable("p.Tag", "green")),
			ch08Frame(8, "tag التالي هو green.", "red\ngreen", ch08Variable("p.Tag", "green")),
		),
	),
	ch08Step("c84.watch-step", "خد بالك",
		ch08Markdown("c84.watch", `تغيير المتغير ~p~ يحرك المؤشر المؤقت؛ تغيير ~Head~ أو ~Next~ يغيّر وصلة في القطر.

قبل التعديل، اعرف مين لسه بيوصل لباقي السلسلة.

تقدر تقرأ field عبر pointer، لكن dereference لـ~nil~ يعمل panic.`)),
	ch08Step("c84.check-step", "فكّر واختبر",
		ch08Markdown("c84.plan", `قسّم الشغل: راجع الحالات الفاضية، ثم التوصيل، ثم الفصل والبحث، وأخيرًا الـsnapshot اللي بترجعه Tags.

اختبر البداية والوسط والنهاية. ما فيش concurrency في عقد المختبر.`),
		ch08Question("c84.check-question", "إمتى تكون ~Next~ نهاية السلسلة؟",
			"", "", []string{"لما تساوي nil", "لما tag تساوي فراغ", "بعد أول عنصر"}, 0,
			[]string{"صحيح: nil تعني عدم وجود عقدة بعدها.", "نص tag لا يحدد نهاية الوصلة.", "السلسلة ممكن يكون فيها عقد كتير."}),
	),
	ch08Step("c84.turn-step", "دورك",
		ch08Markdown("c84.turn", `ابدأ بأبسط حالة في العقد، وشغّل التستات بعد كل method. لو اختفى عنصر، ارجع لتغييرات Head وNext.`)),
}

var ch08CartEnglish = []LessonStep{
	ch08Step("c84.mission-step", "Mission",
		ch08Markdown("c84.mission", `Each cart links to the next. Couple and release carts without losing the convoy.

	The starter has empty methods. Test the contract on a small chain.`)),
	ch08Step("c84.idea-step", "Idea",
		ch08Walkthrough("c84.link-trace", `A struct can point to its own type. ~Next~ links nodes; ~nil~ ends the chain.

Read the next tag.`, "e84-link", ch08LinkExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(2, "The last lamp is blue.", "", ch08Variable("last.Tag", "blue"), ch08Variable("last.Next", "nil")),
			ch08Frame(3, "first links to the last lamp.", "", ch08Variable("first.Tag", "red"), ch08Variable("first.Next", "&last")),
			ch08Frame(4, "Read the tag beyond the link.", "blue", ch08Variable("first.Next.Tag", "blue")),
		),
	),
	ch08Step("c84.try-step", "Try it",
		ch08Question("c84.try-question", "After one link, what is printed?", "e84-walk", ch08WalkLinkExample,
			[]string{"red then green", "green then gold", "red then gold"}, 0,
			[]string{"Right: p moves from head to the middle.", "That starts at the middle, then moves to tail.", "One link reaches the middle first."}),
		ch08Walkthrough("c84.walk-trace", `~p~ is temporary. Moving it leaves ~head~ unchanged. Follow its tag.`, "e84-walk", ch08WalkLinkExample,
			ch08Frame(0, "Before execution.", ""),
			ch08Frame(2, "tail is gold.", "", ch08Variable("tail.Tag", "gold")),
			ch08Frame(3, "middle points to tail.", "", ch08Variable("middle.Tag", "green"), ch08Variable("middle.Next", "&tail")),
			ch08Frame(4, "head points to middle.", "", ch08Variable("head.Tag", "red"), ch08Variable("head.Next", "&middle")),
			ch08Frame(5, "p starts at head.", "", ch08Variable("p.Tag", "red")),
			ch08Frame(6, "The first tag is red.", "red", ch08Variable("p.Tag", "red")),
			ch08Frame(7, "p follows Next to middle.", "red", ch08Variable("p.Tag", "green")),
			ch08Frame(8, "The next tag is green.", "red\ngreen", ch08Variable("p.Tag", "green")),
		),
	),
	ch08Step("c84.watch-step", "Watch out",
		ch08Markdown("c84.watch", `Moving ~p~ leaves links alone; changing ~Head~ or ~Next~ changes the chain.

Know which pointer reaches the rest before changing a link.

Pointer fields are readable; dereferencing ~nil~ panics.`)),
	ch08Step("c84.check-step", "Think & check",
		ch08Markdown("c84.plan", `Phases: empty case, linking, release and search, then Tags snapshot.

	Try the beginning, middle, and end. This lab is sequential.`),
		ch08Question("c84.check-question", "When is ~Next~ the end of the chain?",
			"", "", []string{"When it is nil", "When the tag is empty", "After the first node"}, 0,
			[]string{"Right: nil means no next node.", "An empty tag does not end a link.", "A chain can have many nodes."}),
	),
	ch08Step("c84.turn-step", "Your turn",
		ch08Markdown("c84.turn", `Start with the empty case. Test each method; if a cart disappears, review Head and Next.`)),
}
