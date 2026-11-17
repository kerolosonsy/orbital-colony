package levels

func ch07Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch07Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch07Walkthrough(id, text, code string, exampleID string, frames []LessonFrame) LessonBlock {
	return LessonBlock{
		ID:          id,
		Kind:        "walkthrough",
		Text:        text,
		Code:        code,
		Imports:     []string{"fmt"},
		ExampleID:   exampleID,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

var ch07WaypointsArabic = []LessonStep{
	{
		ID: "c71.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch07Markdown("c71.mission", "عدّي على نقاط التفتيش حسب ترتيب ~colony.Waypoints()~.\n\nكل نقطة ~colony.Point~ فيها الإحداثيين ~X~ و~Y~.\n\nالـstarter فيه دالة وصول ناقصة، ووصل لأول نقطة بالفعل.")},
	},
	{
		ID: "c71.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c71.point-example",
			"الـstruct اسم واحد بيجمع fields مرتبطة. كل نقطة قيمة واحدة، و~p.X~ يقرأ الإحداثي X.\n\nالـuppercase في اسم field بيخليه exported للـpackages التانية.\n\nكل عنصر في الـslice نقطة كاملة. في loop، ~i~ هو index القائمة و~p~ قيمة النقطة.\n\nأسماء الحقول لازم تطابق التعريف. الحروف الكبيرة جزء من الاسم، وهي كمان بتخلي field متاحة برّه package.\n\nاتجاه الوصول يتحدد من الإحداثيين.",
			`type Point struct { X, Y int }
func main() {
	p := Point{X: 3, Y: 1}
	fmt.Println(p.X, p.Y)
}`,
			"e71-point-struct",
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "النوع يجمع الإحداثيين."},
				{Line: 3, Caption: "p فيها X=3 وY=1.", Variables: []LessonVariable{{Name: "p", Value: "{3 1}"}}},
				{Line: 4, Caption: "اقرأ field من النقطة.", Variables: []LessonVariable{{Name: "p", Value: "{3 1}"}}, Output: "3 1"},
			},
		)},
	},
	{
		ID: "c71.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch07Question("c71.predict-field", "توقّع الناتج:\n\n~~~go\ntype Point struct { X, Y int }\nfunc main() {\n\tp := Point{X: 2, Y: 1}\n\tfmt.Println(p.X)\n}\n~~~", []string{"1", "2", "{2 1}"}, 1, []string{
				"ده حقل Y، مش الحقل المطلوب.",
				"قيمة الحقل X هي 2.",
				"الطباعة طلبت field واحدة، مش النقطة كلها.",
			}),
			ch07Walkthrough(
				"c71.field-reveal",
				"امشِ مع إنشاء النقطة، وبعدها راقب الحقل اللي اتطبع.",
				`type Point struct { X, Y int }
func main() {
	p := Point{X: 2, Y: 1}
	fmt.Println(p.X)
}`,
				"e71-point-field",
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 1, Caption: "Point فيها حقلان."},
					{Line: 3, Caption: "النقطة تحمل 2 و1.", Variables: []LessonVariable{{Name: "p", Value: "{2 1}"}}},
					{Line: 4, Caption: "الناتج هو X.", Variables: []LessonVariable{{Name: "p", Value: "{2 1}"}}, Output: "2"},
				},
			),
		},
	},
	{
		ID: "c71.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch07Markdown("c71.watch", "- ~Waypoints()~ بترجع slice بالترتيب المطلوب؛ ما تبدّلش النقاط.\n- ~p.X~ و~p.Y~ إحداثيات، مش أرقام مواقع في الـslice.")},
	},
	{
		ID: "c71.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch07Markdown("c71.plan", "فكّر بالمراحل: اقرأ القائمة، حدّد نقطة الوصول، اختر اتجاه الحركة، ثم راجع مكان الدرون.\n\nالإحداثي ممكن يزيد أو يقل حسب مكان النقطة. راجع نقطة البداية كمان."),
			ch07Question("c71.range-index", "في ~for i, p := range points~، ~i~ بيمثل إيه؟", []string{"ترتيب النقطة في الـslice", "إحداثي X", "عدد حركات الدرون"}, 0, []string{
				"range بيدي index العنصر وقيمته.",
				"إحداثي X موجود جوه p.",
				"عدد الحركات بيتحدد من المسافة.",
			}),
		},
	},
	{
		ID: "c71.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch07Markdown("c71.turn", "اكتب أول محاولة وشغّلها. راقب ترتيب النقاط ومكان الدرون بعد كل وصول.")},
	},
}

var ch07WaypointsEnglish = []LessonStep{
	{
		ID: "c71.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch07Markdown("c71.mission", "Visit the checkpoints in the order from ~colony.Waypoints()~.\n\nEach ~colony.Point~ holds coordinates ~X~ and ~Y~.\n\nThe starter has an unfinished travel function and already reaches the first point.")},
	},
	{
		ID: "c71.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c71.point-example",
			"A struct groups related fields under one name. A point is one value, and ~p.X~ reads its X coordinate.\n\nAn uppercase field name is exported to other packages.\n\nEach slice element is one whole point. In a loop, ~i~ is its list index and ~p~ is the point value.\n\nField names must match the type definition exactly. Capitalization is part of the name and exports the field outside its package.\n\nBoth coordinates determine the travel direction.",
			`type Point struct { X, Y int }
func main() {
	p := Point{X: 3, Y: 1}
	fmt.Println(p.X, p.Y)
}`,
			"e71-point-struct",
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "The type groups both coordinates."},
				{Line: 3, Caption: "p has X=3 and Y=1.", Variables: []LessonVariable{{Name: "p", Value: "{3 1}"}}},
				{Line: 4, Caption: "Read fields from the point.", Variables: []LessonVariable{{Name: "p", Value: "{3 1}"}}, Output: "3 1"},
			},
		)},
	},
	{
		ID: "c71.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch07Question("c71.predict-field", "Predict the output:\n\n~~~go\ntype Point struct { X, Y int }\nfunc main() {\n\tp := Point{X: 2, Y: 1}\n\tfmt.Println(p.X)\n}\n~~~", []string{"1", "2", "{2 1}"}, 1, []string{
				"That is Y, not the requested field.",
				"Field X has value 2.",
				"The print asks for one field, not the whole point.",
			}),
			ch07Walkthrough(
				"c71.field-reveal",
				"Follow the point creation, then watch which field is printed.",
				`type Point struct { X, Y int }
func main() {
	p := Point{X: 2, Y: 1}
	fmt.Println(p.X)
}`,
				"e71-point-field",
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Point has two fields."},
					{Line: 3, Caption: "The point holds 2 and 1.", Variables: []LessonVariable{{Name: "p", Value: "{2 1}"}}},
					{Line: 4, Caption: "The output is X.", Variables: []LessonVariable{{Name: "p", Value: "{2 1}"}}, Output: "2"},
				},
			),
		},
	},
	{
		ID: "c71.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch07Markdown("c71.watch", "- ~Waypoints()~ returns the required slice order; do not rearrange its points.\n- ~p.X~ and ~p.Y~ are coordinates, not indexes in the slice.")},
	},
	{
		ID: "c71.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch07Markdown("c71.plan", "Think in phases: read the list, identify the destination, choose a travel direction, then check the drone's position.\n\nA coordinate can increase or decrease. Check the starting point too."),
			ch07Question("c71.range-index", "In ~for i, p := range points~, what does ~i~ mean?", []string{"The point's slice index", "The X coordinate", "The drone's move count"}, 0, []string{
				"range gives the element index and value.",
				"The X coordinate is inside p.",
				"The move count depends on distance.",
			}),
		},
	},
	{
		ID: "c71.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch07Markdown("c71.turn", "Write and run your first attempt. Watch the point order and the drone's position after each arrival.")},
	},
}

var ch07NavigationArabic = []LessonStep{
	{
		ID: "c72.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch07Markdown("c72.mission", "احصد كريستالات ~colony.Targets()~ بس.\n\nفي كريستالات تانية للتشتيت؛ سيبها حتى لو عدّيت عليها.\n\nالـstarter فيه ~Nav~ و~GoTo~ ناقصة، وبيبدأ من قائمة الأهداف.")},
	},
	{
		ID: "c72.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c72.pointer-method",
			"الـmethod بتجمع السلوك مع النوع. الـpointer receiver يغيّر نفس قيمة Nav، فحالة الموقع تفضل معاك.\n\nالـreceiver بييجي قبل اسم الـmethod في التعريف.\n\nالنداء بـ~nav.Step()~ يبيّن إن السلوك تابع للنوع، والـpointer يخلي التعديل على القيمة الأصلية.",
			`type Nav struct { X int }
func (n *Nav) Step() { n.X++ }
func main() {
	nav := Nav{X: 2}
	nav.Step()
	fmt.Println(nav.X)
}`,
			"e72-pointer-method",
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "Nav تحفظ الموضع."},
				{Line: 4, Caption: "nav تبدأ عند X=2.", Variables: []LessonVariable{{Name: "nav.X", Value: "2"}}},
				{Line: 5, Caption: "الـmethod تغيّر الأصل.", Variables: []LessonVariable{{Name: "nav.X", Value: "3"}}},
				{Line: 6, Caption: "اطبع الحالة الجديدة.", Variables: []LessonVariable{{Name: "nav.X", Value: "3"}}, Output: "3"},
			},
		)},
	},
	{
		ID: "c72.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch07Question("c72.predict-receiver", "توقّع الناتج بعد ~g.Add()~:\n\n~~~go\ntype Gauge struct { N int }\nfunc (g *Gauge) Add() { g.N++ }\nfunc main() {\n\tg := Gauge{N: 4}\n\tg.Add()\n\tfmt.Println(g.N)\n}\n~~~", []string{"4", "5", "0"}, 1, []string{
				"الـpointer receiver يغيّر القيمة الأصلية.",
				"الزيادة ترفع N من 4 إلى 5.",
				"Gauge تبدأ بـ4، مش بالقيمة الصفرية.",
			}),
			ch07Walkthrough(
				"c72.receiver-reveal",
				"تابع قيمة Gauge قبل النداء وبعده.",
				`type Gauge struct { N int }
func (g *Gauge) Add() { g.N++ }
func main() {
	g := Gauge{N: 4}
	g.Add()
	fmt.Println(g.N)
}`,
				"e72-gauge-step",
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 4, Caption: "Gauge تبدأ عند 4.", Variables: []LessonVariable{{Name: "g.N", Value: "4"}}},
					{Line: 5, Caption: "Add زوّدت القيمة.", Variables: []LessonVariable{{Name: "g.N", Value: "5"}}},
					{Line: 6, Caption: "الناتج هو 5.", Variables: []LessonVariable{{Name: "g.N", Value: "5"}}, Output: "5"},
				},
			),
		},
	},
	{
		ID: "c72.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch07Markdown("c72.watch", "- كل ~Move~ و~Harvest~ بتاخد tick؛ حتى ~Move~ الفاشلة.\n- قراءة ~Targets()~ و~Position()~ مجانية. حدّث حالة Nav مع كل حركة.")},
	},
	{
		ID: "c72.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch07Markdown("c72.plan", "مراحل التفكير: خُد هدف من القائمة، حدّد الوصول له، نفّذ الحصاد المقصود، ثم تابع الهدف التالي.\n\nخلّي تمثيل المكان مواكب للحركة. قارن الهدف بالموقع الحالي، مش بأول الخريطة كل مرة.\n\nالحركة الفعلية هي المرجع؛ راجع قيم Nav بعد نداءات الحركة. راجع موقع الدرون بعد الحركة مباشرة في كل خطوة."),
			ch07Question("c72.receiver-kind", "method بتغيّر حقول ~Nav~؛ أنهي receiver يناسبها؟", []string{"~Nav~ بالقيمة", "~*Nav~ بالـpointer", "الموضوعان زي بعض"}, 1, []string{
				"value receiver يشتغل على نسخة من القيمة.",
				"pointer receiver يقدر يعدّل Nav نفسها.",
				"الـreceiver بيغيّر سلوك التعديل.",
			}),
		},
	},
	{
		ID: "c72.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch07Markdown("c72.turn", "جرّب أول نسخة. راجع عدّاد ~Nav~ بعد كل حركة، واتأكد إن الحصاد ما لمسش كريستال خارج الأهداف.")},
	},
}

var ch07NavigationEnglish = []LessonStep{
	{
		ID: "c72.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch07Markdown("c72.mission", "Harvest only crystals in ~colony.Targets()~.\n\nOther crystals are decoys; leave them alone even when you pass them.\n\nThe starter has an unfinished ~Nav.GoTo~ and begins with the target list.")},
	},
	{
		ID: "c72.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c72.pointer-method",
			"A method groups behavior with a type. A pointer receiver changes the same Nav value, so its position stays current.\n\nThe receiver comes before the method name in its declaration.\n\nCalling ~nav.Step()~ shows behavior attached to the type, and the pointer updates the original value.",
			`type Nav struct { X int }
func (n *Nav) Step() { n.X++ }
func main() {
	nav := Nav{X: 2}
	nav.Step()
	fmt.Println(nav.X)
}`,
			"e72-pointer-method",
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Nav stores the position."},
				{Line: 4, Caption: "nav starts at X=2.", Variables: []LessonVariable{{Name: "nav.X", Value: "2"}}},
				{Line: 5, Caption: "The method changes the original.", Variables: []LessonVariable{{Name: "nav.X", Value: "3"}}},
				{Line: 6, Caption: "Print the updated state.", Variables: []LessonVariable{{Name: "nav.X", Value: "3"}}, Output: "3"},
			},
		)},
	},
	{
		ID: "c72.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch07Question("c72.predict-receiver", "Predict the output after ~g.Add()~:\n\n~~~go\ntype Gauge struct { N int }\nfunc (g *Gauge) Add() { g.N++ }\nfunc main() {\n\tg := Gauge{N: 4}\n\tg.Add()\n\tfmt.Println(g.N)\n}\n~~~", []string{"4", "5", "0"}, 1, []string{
				"The pointer receiver changes the original value.",
				"The increment moves N from 4 to 5.",
				"Gauge starts at 4, not its zero value.",
			}),
			ch07Walkthrough(
				"c72.receiver-reveal",
				"Follow Gauge before and after the method call.",
				`type Gauge struct { N int }
func (g *Gauge) Add() { g.N++ }
func main() {
	g := Gauge{N: 4}
	g.Add()
	fmt.Println(g.N)
}`,
				"e72-gauge-step",
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 4, Caption: "Gauge starts at 4.", Variables: []LessonVariable{{Name: "g.N", Value: "4"}}},
					{Line: 5, Caption: "Add increments the value.", Variables: []LessonVariable{{Name: "g.N", Value: "5"}}},
					{Line: 6, Caption: "The output is 5.", Variables: []LessonVariable{{Name: "g.N", Value: "5"}}, Output: "5"},
				},
			),
		},
	},
	{
		ID: "c72.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch07Markdown("c72.watch", "- Each ~Move~ and ~Harvest~ costs one tick, even a failed ~Move~.\n- Reading ~Targets()~ and ~Position()~ is free. Keep Nav's state in sync with movement.")},
	},
	{
		ID: "c72.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch07Markdown("c72.plan", "Think in phases: take a target from the list, plan how to reach it, harvest the intended tile, then continue.\n\nKeep the position model aligned with movement. Compare against the current point, not the map origin each time.\n\nActual movement is the reference; check Nav after move calls. Check the drone's position right after each move."),
			ch07Question("c72.receiver-kind", "A method changes ~Nav~ fields. Which receiver fits?", []string{"~Nav~ by value", "~*Nav~ by pointer", "Either is the same"}, 1, []string{
				"A value receiver works on a copy.",
				"A pointer receiver can update Nav itself.",
				"The receiver changes mutation behavior.",
			}),
		},
	},
	{
		ID: "c72.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch07Markdown("c72.turn", "Try a first version. Check Nav's position after each move, and confirm harvesting never touches a crystal outside the targets.")},
	},
}

var ch07ShortestRouteArabic = []LessonStep{
	{
		ID: "c73.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch07Markdown("c73.mission", "احصد كل الأهداف في ~colony.Targets()~، واختار ترتيب الزيارة.\n\nالمواقع بتتغير، والترتيب اللي الـAPI ترجّعه مش لازم يكون أقصر مشوار.\n\nالـstarter بيحسب مسافات من موقع الدرون الحالي.")},
	},
	{
		ID: "c73.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c73.manhattan-distance",
			"في شبكة من غير حركة قطرية، Manhattan distance تجمع فرق X وفرق Y بعد إزالة الإشارة.\n\nهنا المسافة بتتحسب كـmethod.\n\nالمسافة بتعد خطوات الخانات، مش طول الخط القطري. استدعاء Dist حساب فقط وما بيحرّكش الدرون؛ الإحداثيات تفضل زي ما هي.\n\nلو عكست ترتيب النقطتين، المسافة تفضل زي ما هي.",
			`type Spot struct { X, Y int }
func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}
func main() {
	a, b := Spot{X: 1, Y: 2}, Spot{X: 4, Y: 3}
	fmt.Println(a.Dist(b))
}`,
			"e73-manhattan-distance",
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "Spot تحفظ X وY."},
				{Line: 8, Caption: "Dist تجمع فرق الإحداثيين."},
				{Line: 12, Caption: "النقطتان (1,2) و(4,3).", Variables: []LessonVariable{{Name: "a", Value: "{1 2}"}, {Name: "b", Value: "{4 3}"}}},
				{Line: 13, Caption: "المسافة تساوي 4.", Variables: []LessonVariable{{Name: "a", Value: "{1 2}"}, {Name: "b", Value: "{4 3}"}}, Output: "4"},
			},
		)},
	},
	{
		ID: "c73.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch07Question("c73.predict-distance", "توقّع ناتج المسافة بين (2,1) و(3,3):\n\n~~~go\ntype Spot struct { X, Y int }\nfunc abs(n int) int {\n\tif n < 0 {\n\t\tn = -n\n\t}\n\treturn n\n}\nfunc (s Spot) Dist(o Spot) int {\n\treturn abs(s.X-o.X) + abs(s.Y-o.Y)\n}\nfunc main() {\n\ta, b := Spot{X: 2, Y: 1}, Spot{X: 3, Y: 3}\n\tfmt.Println(a.Dist(b))\n}\n~~~", []string{"2", "3", "4"}, 1, []string{
				"حركة X واحدة ومعاها فرق Y.",
				"فرق X واحد وفرق Y اتنين؛ المجموع 3.",
				"المسافة مش طول الخط القطري.",
			}),
			ch07Walkthrough(
				"c73.distance-reveal",
				"اتبع حساب المسافة على بيانات تدريب صغيرة.",
				`type Spot struct { X, Y int }
func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}
func main() {
	a, b := Spot{X: 2, Y: 1}, Spot{X: 3, Y: 3}
	fmt.Println(a.Dist(b))
}`,
				"e73-distance-reveal",
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 12, Caption: "الفرق الأفقي واحد.", Variables: []LessonVariable{{Name: "a", Value: "{2 1}"}, {Name: "b", Value: "{3 3}"}}},
					{Line: 13, Caption: "الناتج يجمع 1 و2.", Variables: []LessonVariable{{Name: "a", Value: "{2 1}"}, {Name: "b", Value: "{3 3}"}}, Output: "3"},
				},
			),
		},
	},
	{
		ID: "c73.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch07Markdown("c73.watch", "- الـgreedy يختار الأقرب دلوقتي؛ ده مش ضمان لأقصر route كامل.\n- ~Move~ tick حتى لو فشلت، و~Harvest~ tick؛ الحساب وقراءة ~Targets()~ من غير ticks.")},
	},
	{
		ID: "c73.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch07Markdown("c73.plan", "مراحل التفكير: قارن المسافات من مكانك، اختار هدف تالي، حدّث موضعك، وكرّر.\n\nاعتبر nearest-first طريقة عملية، مش برهان على الحل الأمثل."),
			ch07Question("c73.greedy-guarantee", "هل أقرب هدف دلوقتي يضمن أقصر مسار كامل؟", []string{"أيوه، دايماً", "لأ، مش مضمون", "بس لو الأهداف عشوائية"}, 1, []string{
				"قرار واحد محلي مش كفاية لضمان المسار كله.",
				"الاختيار المحلي قد يترك ترتيباً أطول بعده.",
				"عشوائية ترتيب القائمة لا تعطي ضماناً.",
			}),
		},
	},
	{
		ID: "c73.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch07Markdown("c73.turn", "ابدأ بتجربة ترتيب تختاره. شغّل البرنامج وراقب المسافة الكلية قبل ما تغيّر قرارك.")},
	},
}

var ch07ShortestRouteEnglish = []LessonStep{
	{
		ID: "c73.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch07Markdown("c73.mission", "Harvest every target in ~colony.Targets()~, and choose the visit order.\n\nLocations vary, and the API's order is not necessarily the shortest route.\n\nThe starter calculates distances from the drone's current position.")},
	},
	{
		ID: "c73.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c73.manhattan-distance",
			"On a grid without diagonal moves, Manhattan distance adds the absolute X and Y differences.\n\nHere, distance is a method.\n\nIt counts tile steps, not diagonal line length. Calling Dist only calculates; it does not move the drone or change its coordinates.\n\nReversing the two points leaves the distance unchanged.",
			`type Spot struct { X, Y int }
func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}
func main() {
	a, b := Spot{X: 1, Y: 2}, Spot{X: 4, Y: 3}
	fmt.Println(a.Dist(b))
}`,
			"e73-manhattan-distance",
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Spot stores X and Y."},
				{Line: 8, Caption: "Dist adds coordinate differences."},
				{Line: 12, Caption: "The points are (1,2) and (4,3).", Variables: []LessonVariable{{Name: "a", Value: "{1 2}"}, {Name: "b", Value: "{4 3}"}}},
				{Line: 13, Caption: "The distance is 4.", Variables: []LessonVariable{{Name: "a", Value: "{1 2}"}, {Name: "b", Value: "{4 3}"}}, Output: "4"},
			},
		)},
	},
	{
		ID: "c73.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch07Question("c73.predict-distance", "Predict the distance between (2,1) and (3,3):\n\n~~~go\ntype Spot struct { X, Y int }\nfunc abs(n int) int {\n\tif n < 0 {\n\t\tn = -n\n\t}\n\treturn n\n}\nfunc (s Spot) Dist(o Spot) int {\n\treturn abs(s.X-o.X) + abs(s.Y-o.Y)\n}\nfunc main() {\n\ta, b := Spot{X: 2, Y: 1}, Spot{X: 3, Y: 3}\n\tfmt.Println(a.Dist(b))\n}\n~~~", []string{"2", "3", "4"}, 1, []string{
				"One X move combines with the Y difference.",
				"X differs by one and Y by two; the sum is 3.",
				"The distance is not the diagonal line length.",
			}),
			ch07Walkthrough(
				"c73.distance-reveal",
				"Follow the distance calculation with small practice data.",
				`type Spot struct { X, Y int }
func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}
func main() {
	a, b := Spot{X: 2, Y: 1}, Spot{X: 3, Y: 3}
	fmt.Println(a.Dist(b))
}`,
				"e73-distance-reveal",
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 12, Caption: "The horizontal difference is one.", Variables: []LessonVariable{{Name: "a", Value: "{2 1}"}, {Name: "b", Value: "{3 3}"}}},
					{Line: 13, Caption: "The output adds 1 and 2.", Variables: []LessonVariable{{Name: "a", Value: "{2 1}"}, {Name: "b", Value: "{3 3}"}}, Output: "3"},
				},
			),
		},
	},
	{
		ID: "c73.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch07Markdown("c73.watch", "- Greedy picks the nearest target now; it does not guarantee a shortest full route.\n- ~Move~ costs a tick even when it fails; ~Harvest~ costs a tick too. Calculation and ~Targets()~ are free.")},
	},
	{
		ID: "c73.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch07Markdown("c73.plan", "Think in phases: compare distances from your position, choose a next target, update your position, and repeat.\n\nTreat nearest-first as practical, not proof of an optimal route."),
			ch07Question("c73.greedy-guarantee", "Does the nearest target now guarantee the shortest full route?", []string{"Yes, always", "No, not guaranteed", "Only for random targets"}, 1, []string{
				"One local decision cannot guarantee the whole route.",
				"A local choice may leave a longer order afterward.",
				"Random list order gives no such guarantee.",
			}),
		},
	},
	{
		ID: "c73.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch07Markdown("c73.turn", "Try an order you choose. Run it and watch the total distance before changing your decision.")},
	},
}

var ch07CompositeArabic = []LessonStep{
	{
		ID: "c74.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch07Markdown("c74.mission", "اقرأ إجمالي الفاتورة وعدد الطاقم من ~colony.Invoice()~.\n\nاستخدم ~Ledger~ embedded جوه ~CrewUnit~، وmethod ترجع slice بالأنصبة.\n\nابعت التقرير مرة واحدة، وحافظ على إجمالي الـcents.")},
	},
	{
		ID: "c74.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c74.embedding-method",
			"الـembedding يخلّي field النوع الداخلي ظاهرة بالاسم مباشرة. method القراءة تقدر تستخدم الحقل المرقّى.\n\nالنوع الخارجي لسه شايل قيمة النوع الداخلي؛ ده مش inheritance ولا copy منفصلة.\n\nلو field names اتعارضت، المسار الكامل يوضح أي قيمة تقصد.\n\nالـvalue receiver يستقبل نسخة؛ يناسب method قراءة زي Count.\n\nما تعتمدش عليه لحفظ تعديل على الأصل.",
			`type Cargo struct { Units int }
type Ship struct { Cargo }
func (s Ship) Count() int { return s.Units }
func main() {
	ship := Ship{Cargo: Cargo{Units: 3}}
	fmt.Println(ship.Count())
}`,
			"e74-promoted-method",
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 2, Caption: "Ship يضم Cargo بالـembedding."},
				{Line: 4, Caption: "Count تقرأ القيمة المرقّاة."},
				{Line: 6, Caption: "ship فيها 3 وحدات.", Variables: []LessonVariable{{Name: "ship.Units", Value: "3"}}},
				{Line: 7, Caption: "الـmethod تطبع 3.", Variables: []LessonVariable{{Name: "ship.Units", Value: "3"}}, Output: "3"},
			},
		)},
	},
	{
		ID: "c74.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch07Question("c74.predict-promoted", "توقّع ناتج ~fmt.Println(r.Items)~:\n\n~~~go\ntype Pack struct { Items int }\ntype Rover struct { Pack }\nfunc main() {\n\tr := Rover{Pack: Pack{Items: 4}}\n\tfmt.Println(r.Items)\n}\n~~~", []string{"0", "4", "خطأ compile"}, 1, []string{
				"Rover فيه Pack، وItems مش صفر.",
				"Items field مضمّنة ومقروءة باسم مباشر.",
				"الـembedding يسمح بقراءة field المرقّاة.",
			}),
			ch07Walkthrough(
				"c74.promoted-reveal",
				"راقب الوصول إلى field داخل النوع المضمّن.",
				`type Pack struct { Items int }
type Rover struct { Pack }
func main() {
	r := Rover{Pack: Pack{Items: 4}}
	fmt.Println(r.Items)
}`,
				"e74-promoted-field",
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 2, Caption: "Rover embeds Pack."},
					{Line: 4, Caption: "Pack carries four items.", Variables: []LessonVariable{{Name: "r.Items", Value: "4"}}},
					{Line: 5, Caption: "The promoted field prints.", Variables: []LessonVariable{{Name: "r.Items", Value: "4"}}, Output: "4"},
				},
			),
		},
	},
	{
		ID: "c74.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch07Markdown("c74.watch", "- اكتب النوع المضمّن لوحده؛ ~ledger Ledger~ field عادية مش embedding.\n- القسمة الصحيحة تسيب باقي؛ تجاهله يضيّع جزء من الإجمالي.")},
	},
	{
		ID: "c74.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch07Markdown("c74.plan", "فكّر بالمراحل: اقرأ القيم، مثّلها بأنواع واضحة، اجمع المنطق في method، ثم راجع المجموع النهائي.\n\nجرّب أيضاً حالة إجمالي صفر وفرد واحد."),
			ch07Question("c74.remainder", "ما معنى ~17 % 5~؟", []string{"ناتج القسمة 3", "الباقي 2", "الناتج 3.4"}, 1, []string{
				"القسمة الصحيحة تستخدم / للحصول على 3.",
				"عامل % يعيد المتبقي بعد القسمة.",
				"Go لا يحوّل القسمة الصحيحة لعشرية هنا.",
			}),
		},
	},
	{
		ID: "c74.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch07Markdown("c74.turn", "كوّن الأنواع والـmethod بنفسك، وبعدها قارن slice التقرير بإجمالي الفاتورة. راجع الناتج ببيانات افتراضية.")},
	},
}

var ch07CompositeEnglish = []LessonStep{
	{
		ID: "c74.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch07Markdown("c74.mission", "Read the invoice total and crew count from ~colony.Invoice()~.\n\nEmbed ~Ledger~ inside ~CrewUnit~, with a method that returns a shares slice.\n\nSend one report and preserve the total cents.")},
	},
	{
		ID: "c74.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c74.embedding-method",
			"Embedding makes an inner type's fields available by direct name. A read method can use the promoted field.\n\nThe outer type still contains the inner value; this is not class inheritance or a separate copy.\n\nIf field names conflict, the full path shows which value you mean.\n\nA value receiver gets a copy and fits a read method like Count.\n\nDo not rely on it to retain changes to the original.",
			`type Cargo struct { Units int }
type Ship struct { Cargo }
func (s Ship) Count() int { return s.Units }
func main() {
	ship := Ship{Cargo: Cargo{Units: 3}}
	fmt.Println(ship.Count())
}`,
			"e74-promoted-method",
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 2, Caption: "Ship embeds Cargo."},
				{Line: 4, Caption: "Count reads the promoted value."},
				{Line: 6, Caption: "ship holds 3 units.", Variables: []LessonVariable{{Name: "ship.Units", Value: "3"}}},
				{Line: 7, Caption: "The method prints 3.", Variables: []LessonVariable{{Name: "ship.Units", Value: "3"}}, Output: "3"},
			},
		)},
	},
	{
		ID: "c74.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch07Question("c74.predict-promoted", "Predict the output of ~fmt.Println(r.Items)~:\n\n~~~go\ntype Pack struct { Items int }\ntype Rover struct { Pack }\nfunc main() {\n\tr := Rover{Pack: Pack{Items: 4}}\n\tfmt.Println(r.Items)\n}\n~~~", []string{"0", "4", "Compile error"}, 1, []string{
				"Rover contains Pack, and Items is not zero.",
				"Items is embedded and readable by its direct name.",
				"Embedding permits reading the promoted field.",
			}),
			ch07Walkthrough(
				"c74.promoted-reveal",
				"Watch access to a field inside the embedded type.",
				`type Pack struct { Items int }
type Rover struct { Pack }
func main() {
	r := Rover{Pack: Pack{Items: 4}}
	fmt.Println(r.Items)
}`,
				"e74-promoted-field",
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 2, Caption: "Rover embeds Pack."},
					{Line: 4, Caption: "Pack contains four items.", Variables: []LessonVariable{{Name: "r.Items", Value: "4"}}},
					{Line: 5, Caption: "The promoted field prints.", Variables: []LessonVariable{{Name: "r.Items", Value: "4"}}, Output: "4"},
				},
			),
		},
	},
	{
		ID: "c74.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch07Markdown("c74.watch", "- Write the embedded type alone; ~ledger Ledger~ is an ordinary field, not embedding.\n- Integer division leaves a remainder; ignoring it loses part of the total.")},
	},
	{
		ID: "c74.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch07Markdown("c74.plan", "Think in phases: read the values, represent them with clear types, group the logic in a method, then check the final sum.\n\nAlso consider zero total and one crew member."),
			ch07Question("c74.remainder", "What does ~17 % 5~ mean?", []string{"The quotient is 3", "The remainder is 2", "The result is 3.4"}, 1, []string{
				"Integer division uses / to get 3.",
				"The % operator returns what remains after division.",
				"Go does not turn this integer division into a decimal.",
			}),
		},
	},
	{
		ID: "c74.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch07Markdown("c74.turn", "Build the types and method yourself, then compare the report slice with the invoice total. Check using sample data.")},
	},
}

var ch07CrewExamArabic = []LessonStep{
	{
		ID: "c75.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch07Markdown("c75.mission", "جاوب أسئلة أرض المحاصيل، وبعدها ابعت نسبة الأسئلة الصح اللي ينفع تتجاوب.\n\nبنك الأسئلة وترتيب الاختيارات بيتغيروا.\n\nفي سؤال offline من غير اختيارات؛ سيبه من غير إجابة.")},
	},
	{
		ID: "c75.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c75.question-struct",
			"كل Question struct فيها prompt وslice اختيارات. range بيدي index السؤال وقيمة منه.\n\nالـindex بيحدد رقم السؤال الأصلي، والـChoices لها indexes خاصة بيها.\n\nترتيب الاختيارات بيتغير؛ اقرأ النصوص بدل ما تفترض رقم ثابت. رقم الاختيار يبدأ من صفر.",
			`type Question struct {
	Prompt string
	Choices []string
}
func main() {
	questions := []Question{{Prompt: "ground?", Choices: []string{"soil", "ice"}}}
	for i, q := range questions {
		fmt.Println(i, q.Prompt, len(q.Choices))
	}
}`,
			"e75-question-struct",
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "Question تجمع prompt واختيارات."},
				{Line: 6, Caption: "القائمة فيها سؤال واحد."},
				{Line: 7, Caption: "i هو index السؤال.", Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "q.Prompt", Value: "ground?"}}},
				{Line: 8, Caption: "للسؤال اختياران.", Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "q.Prompt", Value: "ground?"}}, Output: "0 ground? 2"},
			},
		)},
	},
	{
		ID: "c75.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch07Question("c75.predict-range", "توقّع سطور الطباعة:\n\n~~~go\nfunc main() {\n\tlabels := []string{\"A\", \"B\"}\n\tfor i, label := range labels {\n\t\tfmt.Println(i, label)\n\t}\n}\n~~~", []string{"0 A ثم 1 B", "1 A ثم 2 B", "A ثم B بلا أرقام"}, 0, []string{
				"index يبدأ من صفر، وبيرافق قيمة كل عنصر.",
				"index مش بيبدأ من واحد.",
				"range تقدر ترجع index والقيمة معاً.",
			}),
			ch07Walkthrough(
				"c75.range-reveal",
				"تابع index والقيمة في كل لفة.",
				`func main() {
	labels := []string{"A", "B"}
	for i, label := range labels {
		fmt.Println(i, label)
	}
}`,
				"e75-range-output",
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 2, Caption: "labels فيها A وB.", Slice: &LessonSlice{Name: "labels", Values: []int{0, 1}, Length: 2, ElementsOnly: true}},
					{Line: 4, Caption: "أول index هو صفر.", Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "label", Value: "A"}}, Output: "0 A"},
					{Line: 4, Caption: "التالي index واحد.", Variables: []LessonVariable{{Name: "i", Value: "1"}, {Name: "label", Value: "B"}}, Output: "0 A\n1 B"},
				},
			),
		},
	},
	{
		ID: "c75.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch07Markdown("c75.watch", "- افحص ~len(q.Choices)~ قبل أي index؛ القائمة ممكن تكون فاضية.\n- ~Answer(i, choice)~ ترجع صلاحية الأرقام، مش صحة المعلومة.")},
	},
	{
		ID: "c75.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch07Markdown("c75.plan", "مراحل التفكير: اقرأ بيانات السؤال، قرر هل هو قابل للإجابة، اختار من بياناته، ثم سجّل النتيجة واحسب النسبة.\n\nاحتفظ بالـindex الأصلي من range."),
			ch07Question("c75.empty-choices", "سؤال اختياراته فاضية؛ إيه التصرف الآمن؟", []string{"اقرأ ~Choices[0]~", "تخطّى السؤال", "اعتبر الإجابة صحيحة"}, 1, []string{
				"مفيش عنصر عند index صفر في slice فاضية.",
				"التخطي يمنع الفهرسة وإرسال إجابة غير موجودة.",
				"السؤال الفاضي ماقدّمش إجابة نقارنها.",
			}),
		},
	},
	{
		ID: "c75.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch07Markdown("c75.turn", "ابدأ بنسخة صغيرة من loop الأسئلة. شغّلها وراجع index كل سؤال والاختيار اللي سجلته.")},
	},
}

var ch07CrewExamEnglish = []LessonStep{
	{
		ID: "c75.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch07Markdown("c75.mission", "Answer the crop-ground questions, then report the share of usable questions answered correctly.\n\nThe question bank and choice order vary.\n\nOne offline question has no choices; leave it unanswered.")},
	},
	{
		ID: "c75.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch07Walkthrough(
			"c75.question-struct",
			"Each Question struct holds a prompt and a choices slice. range gives the question index and its value.\n\nThe index identifies the original question, while Choices has its own indexes.\n\nChoice order varies; read the text instead of assuming a fixed number. Choice indexes start at zero.",
			`type Question struct {
	Prompt string
	Choices []string
}
func main() {
	questions := []Question{{Prompt: "ground?", Choices: []string{"soil", "ice"}}}
	for i, q := range questions {
		fmt.Println(i, q.Prompt, len(q.Choices))
	}
}`,
			"e75-question-struct",
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Question groups prompt and choices."},
				{Line: 6, Caption: "The list has one question."},
				{Line: 7, Caption: "i is the question index.", Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "q.Prompt", Value: "ground?"}}},
				{Line: 8, Caption: "The question has two choices.", Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "q.Prompt", Value: "ground?"}}, Output: "0 ground? 2"},
			},
		)},
	},
	{
		ID: "c75.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch07Question("c75.predict-range", "Predict the printed lines:\n\n~~~go\nfunc main() {\n\tlabels := []string{\"A\", \"B\"}\n\tfor i, label := range labels {\n\t\tfmt.Println(i, label)\n\t}\n}\n~~~", []string{"0 A then 1 B", "1 A then 2 B", "A then B without numbers"}, 0, []string{
				"Indexes start at zero and pair with each value.",
				"An index does not start at one.",
				"range can return both the index and value.",
			}),
			ch07Walkthrough(
				"c75.range-reveal",
				"Follow the index and value in each loop.",
				`func main() {
	labels := []string{"A", "B"}
	for i, label := range labels {
		fmt.Println(i, label)
	}
}`,
				"e75-range-output",
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 2, Caption: "labels contains A and B.", Slice: &LessonSlice{Name: "labels", Values: []int{0, 1}, Length: 2, ElementsOnly: true}},
					{Line: 4, Caption: "The first index is zero.", Variables: []LessonVariable{{Name: "i", Value: "0"}, {Name: "label", Value: "A"}}, Output: "0 A"},
					{Line: 4, Caption: "The next index is one.", Variables: []LessonVariable{{Name: "i", Value: "1"}, {Name: "label", Value: "B"}}, Output: "0 A\n1 B"},
				},
			),
		},
	},
	{
		ID: "c75.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch07Markdown("c75.watch", "- Check ~len(q.Choices)~ before indexing; the slice may be empty.\n- ~Answer(i, choice)~ validates the indexes, not the factual answer.")},
	},
	{
		ID: "c75.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch07Markdown("c75.plan", "Think in phases: read the question data, decide if it is answerable, choose from its data, then record and score.\n\nKeep the original index from range."),
			ch07Question("c75.empty-choices", "A question has no choices; what is safe?", []string{"Read ~Choices[0]~", "Skip the question", "Treat it as correct"}, 1, []string{
				"A slice with no elements has no index zero.",
				"Skipping avoids indexing and submitting a missing choice.",
				"An empty question provides no answer to compare.",
			}),
		},
	},
	{
		ID: "c75.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch07Markdown("c75.turn", "Start with a small version of the question loop. Run it and review each question index and the choice you recorded.")},
	},
}
