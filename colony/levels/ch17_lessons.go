package levels

type ch17FrameCopy struct {
	line   int
	ar     string
	en     string
	vars   []LessonVariable
	output string
}

func ch17Copy(ar, en string, english bool) string {
	if english {
		return en
	}
	return ar
}

func ch17Markdown(id, ar, en string, english bool) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: ch17Copy(ar, en, english)}
}

func ch17Question(id, ar, en, code string, optionsAR, optionsEN []string, answer int, feedbackAR, feedbackEN []string, english bool) LessonBlock {
	return LessonBlock{
		ID:   id,
		Kind: "question",
		Text: ch17Copy(ar, en, english) + "\n\n~~~go\n" + code + "\n~~~",
		Question: &LessonQuestion{
			Options:  ch17CopySlice(optionsAR, optionsEN, english),
			Answer:   answer,
			Feedback: ch17CopySlice(feedbackAR, feedbackEN, english),
		},
	}
}

func ch17Choice(id, ar, en string, optionsAR, optionsEN []string, answer int, feedbackAR, feedbackEN []string, english bool) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     ch17Copy(ar, en, english),
		Question: &LessonQuestion{Options: ch17CopySlice(optionsAR, optionsEN, english), Answer: answer, Feedback: ch17CopySlice(feedbackAR, feedbackEN, english)},
	}
}

func ch17CopySlice(ar, en []string, english bool) []string {
	if english {
		return en
	}
	return ar
}

func ch17Frames(copies []ch17FrameCopy, english bool) []LessonFrame {
	frames := make([]LessonFrame, len(copies))
	for i, frame := range copies {
		frames[i] = LessonFrame{Line: frame.line, Caption: ch17Copy(frame.ar, frame.en, english), Variables: frame.vars, Output: frame.output}
	}
	return frames
}

func ch17Walkthrough(id, ar, en, code, exampleID string, imports []string, frames []ch17FrameCopy, english bool) LessonBlock {
	return LessonBlock{
		ID:        id,
		Kind:      "walkthrough",
		Text:      ch17Copy(ar, en, english),
		Code:      code,
		Imports:   imports,
		ExampleID: exampleID,
		Walkthrough: &LessonWalkthrough{
			Frames: ch17Frames(frames, english),
		},
	}
}

func ch17Mission1(english bool) []LessonStep {
	cell := `type Cell struct { X, Y, Z int }
func main() {
	p := Cell{X: 1, Y: 2, Z: 0}
	fmt.Println(p.X, p.Y, p.Z)
	p.Z = 1
	fmt.Println(p.X, p.Y, p.Z)
}`
	increment := `func main() {
	z := 0
	z++
	fmt.Println(z)
}`
	return []LessonStep{
		{ID: "c171.mission-step", Title: ch17Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch17Markdown("c171.mission",
			"1. اجمع 3 حديد من الطبقة اللي تحت.\n\n2. ارجع للسطح بعد ما تخلص.\n\n3. الـstarter فتح shaft ونزل فعلًا؛ هدف الذهب 12 ticks.",
			"1. Collect 3 iron from the layer below.\n\n2. Return to the surface when finished.\n\n3. The starter opened a shaft and descended; gold target is 12 ticks.", english)}},
		{ID: "c171.idea-step", Title: ch17Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch17Walkthrough("c171.z-coordinate",
			"اعتبر المكان نقطة فيها ~x~ و~y~ و~z~.\n\nتغيير ~z~ يغيّر الطبقة، و~x~ و~y~ يفضلوا نفس المكان الأفقي.",
			"Think of a location as ~x~, ~y~, and ~z~.\n\nChanging ~z~ changes the layer while ~x~ and ~y~ keep the same horizontal tile.",
			cell, "e171-z-coordinate", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 1, ar: "Cell تجمع 3 إحداثيات.", en: "Cell groups three coordinates."},
				{line: 3, ar: "النقطة تبدأ على السطح.", en: "The point starts on the surface.", vars: []LessonVariable{{Name: "p", Value: "{1 2 0}"}}},
				{line: 4, ar: "x وy ثابتين؛ z=0.", en: "x and y stay; z=0.", vars: []LessonVariable{{Name: "p", Value: "{1 2 0}"}}, output: "1 2 0"},
				{line: 5, ar: "غيّرنا الطبقة بس.", en: "Only the layer changes.", vars: []LessonVariable{{Name: "p", Value: "{1 2 1}"}}, output: "1 2 0"},
				{line: 6, ar: "x وy زي ما هما.", en: "x and y remain the same.", vars: []LessonVariable{{Name: "p", Value: "{1 2 1}"}}, output: "1 2 0\n1 2 1"},
			}, english)}},
		{ID: "c171.try-step", Title: ch17Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch17Question("c171.predict-layer", "إيه اللي هيتطبع بعد زيادة z؟", "What prints after z increases?", increment, []string{"1", "0", "2"}, []string{"1", "0", "2"}, 0,
				[]string{"بدأت z بصفر وزادت مرة.", "دي قيمتها قبل الزيادة.", "الزيادة خطوة واحدة بس."},
				[]string{"z starts at zero and increases once.", "That is its value before the increment.", "It increases by one step."}, english),
			ch17Walkthrough("c171.layer-reveal", "اتبع قيمة الطبقة خطوة خطوة.", "Follow the layer value one step at a time.", increment, "e171-layer", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "z تبدأ بصفر.", en: "z starts at zero.", vars: []LessonVariable{{Name: "z", Value: "0"}}},
				{line: 3, ar: "زودناها واحد.", en: "Increase it by one.", vars: []LessonVariable{{Name: "z", Value: "1"}}},
				{line: 4, ar: "يطبع رقم الطبقة.", en: "Prints the layer number.", vars: []LessonVariable{{Name: "z", Value: "1"}}, output: "1"},
			}, english),
		}},
		{ID: "c171.watch-step", Title: ch17Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch17Markdown("c171.watch",
			"~Dig~ بيفتح الخانة اللي جنبك تحت الأرض؛ ~Move~ يعدّي في نفق مفتوح.\n\n~DigDown~ و~Descend~ خطوتين مختلفتين.",
			"~Dig~ opens the adjacent underground tile; ~Move~ travels through an open tunnel.\n\n~DigDown~ and ~Descend~ are separate actions.", english)}},
		{ID: "c171.check-step", Title: ch17Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch17Markdown("c171.plan", "افصل السؤال لمرحلتين: اجمع الكمية، وبعدها راجع مكان الدرون.\n\n~Layer()~ بتقولك العمق الحالي من غير ما تنقلك.", "Separate the task into two checks: collect the amount, then check the drone's location.\n\n~Layer()~ reads current depth without moving you.", english),
			ch17Choice("c171.coord-question", "لما ~p.Z~ تبقى 1، إيه اللي بيتغير؟", "When ~p.Z~ becomes 1, what changes?", []string{"الطبقة فقط", "كل الإحداثيات", "قيمة x فقط"}, []string{"Only the layer", "Every coordinate", "Only x"}, 0,
				[]string{"z هي إحداثية العمق؛ x وy ما اتغيروش.", "إسناد z ما بيغيّرش باقي الحقول.", "x وy بيمثلوا المكان الأفقي."},
				[]string{"z is the depth coordinate; x and y stayed unchanged.", "Assigning z does not change other fields.", "x and y represent the horizontal location."}, english),
		}},
		{ID: "c171.turn-step", Title: ch17Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch17Markdown("c171.turn",
			"كمّل من الطبقة اللي الـstarter نزل لها.\n\nتابع الحديد ومكان الدرون، وبعدين شغّل المحاولة كاملة.",
			"Continue from the layer where the starter descended.\n\nTrack the iron and drone position, then run the full attempt.", english)}},
	}
}

func ch17Mission2(english bool) []LessonStep {
	grid := `func main() {
	layers := [][][]string{
		{{"A", "B"}, {"C", "D"}},
		{{"E", "F"}, {"G", "H"}},
	}
	fmt.Println(layers[1][0][1])
}`
	index := `func main() {
	layers := [][][]int{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}
	fmt.Println(layers[1][1][0])
}`
	return []LessonStep{
		{ID: "c172.mission-step", Title: ch17Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch17Markdown("c172.mission",
			"1. عدّ الخانات المفتوحة في كل طبقة تحت الأرض.\n\n2. ابعت الأعداد في slice واحدة وبالترتيب.\n\n3. الـstarter بيطبع عدد الطبقات وأبعاد أول grid بس.",
			"1. Count open tiles on every underground layer.\n\n2. Report counts in one ordered slice.\n\n3. The starter prints layer count and the first grid's dimensions only.", english)}},
		{ID: "c172.idea-step", Title: ch17Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch17Walkthrough("c172.grid-index",
			"لو جمعت الطبقات في slice واحدة، الترتيب هنا ~[z][x][y]~.\n\n~LayerGrid(z)~ لوحدها بترجع ~grid[x][y]~؛ ثبّت الترتيب في دماغك.",
			"When layers are collected in one slice, the order here is ~[z][x][y]~.\n\nOne ~LayerGrid(z)~ returns ~grid[x][y]~; keep that order in mind.",
			grid, "e172-grid-index", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "طبقتين للتوضيح.", en: "Two illustrative layers.", vars: []LessonVariable{{Name: "layers", Value: "2 layers; each 2 by 2"}}},
				{line: 6, ar: "نختار z=1 وx=0 وy=1.", en: "Choose z=1, x=0, and y=1.", vars: []LessonVariable{{Name: "layers[1][0][1]", Value: "F"}}},
				{line: 6, ar: "الخانة المختارة فيها F.", en: "The selected tile contains F.", output: "F"},
			}, english)}},
		{ID: "c172.try-step", Title: ch17Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch17Question("c172.predict-index", "أنهي قيمة هتطلع؟", "Which value is printed?", index, []string{"7", "6", "3"}, []string{"7", "6", "3"}, 0,
				[]string{"الطبقة 1، ثم x=1، ثم y=0 توصل لـ7.", "6 عند y=1 في نفس العمود.", "3 في الطبقة 0، العمود 1."},
				[]string{"Layer 1, then x=1, then y=0 reaches 7.", "Six is at y=1 in that column.", "Three is in layer 0, column 1."}, english),
			ch17Walkthrough("c172.index-reveal", "اقرأ كل index بالترتيب.", "Read each index in order.", index, "e172-index", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "قيم مختلفة لكل بُعد.", en: "Distinct values for each dimension.", vars: []LessonVariable{{Name: "layers", Value: "2 layers; each 2 by 2"}}},
				{line: 3, ar: "اختار الطبقة التانية.", en: "Select the second layer.", vars: []LessonVariable{{Name: "z", Value: "1"}}},
				{line: 3, ar: "اختار x=1 ثم y=0.", en: "Select x=1, then y=0.", vars: []LessonVariable{{Name: "layers[1][1][0]", Value: "7"}}},
				{line: 3, ar: "الناتج 7.", en: "The result is 7.", output: "7"},
			}, english),
		}},
		{ID: "c172.watch-step", Title: ch17Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch17Markdown("c172.watch",
			"الخريطة مش لازم تبقى مربعة؛ استخدم طول كل slice في مستواها.\n\nالخانة المفتوحة تشمل ~Tunnel~ و~Shaft~.",
			"The map need not be square; use each slice's own length.\n\nOpen ground includes ~Tunnel~ and ~Shaft~.", english)}},
		{ID: "c172.check-step", Title: ch17Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch17Markdown("c172.plan", "اكتب شكل الأبعاد بتعليق صغير، وافحص قراءة خانة واحدة قبل توسيع اللفّات.\n\nبعد كل طبقة، ضيف عددها لمكانها في slice التقرير.", "Write the dimension order in a short comment and verify one tile before expanding loops.\n\nAdd each layer's count to its position in the report slice.", english),
			ch17Choice("c172.layer-grid-question", "داخل ~LayerGrid(z)~، إيه ترتيب index؟", "Inside ~LayerGrid(z)~, what is the index order?", []string{"[x][y]", "[y][x]", "[z][y]"}, []string{"[x][y]", "[y][x]", "[z][y]"}, 0,
				[]string{"الدالة بترجع grid لطبقة واحدة: x وبعدها y.", "ترتيب المحاور جزء من العقد.", "z اختار الطبقة قبل استدعاء الدالة."},
				[]string{"The function returns one layer grid: x, then y.", "Axis order is part of the contract.", "z selects the layer before the call."}, english),
		}},
		{ID: "c172.turn-step", Title: ch17Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch17Markdown("c172.turn",
			"ابدأ بالـgrid اللي الـstarter عرضها، وبعدها وسّع الحساب لكل layer.\n\nراجع ترتيب الأعداد قبل ما تبعتها بـ~colony.Report~.",
			"Start with the grid shown by the starter, then extend the count to every layer.\n\nCheck count order before sending it with ~colony.Report~.", english)}},
	}
}

func ch17Mission3(english bool) []LessonStep {
	keys := `type Cell struct { X, Y, Z int }
func main() {
	seen := map[Cell]bool{}
	p := Cell{X: 1, Y: 0, Z: 1}
	seen[p] = true
	fmt.Println(seen[Cell{X: 1, Y: 0, Z: 1}])
	fmt.Println(seen[Cell{X: 1, Y: 0, Z: 2}])
}`
	queue := `func main() {
	queue := []string{"dock", "shaft"}
	current := queue[0]
	queue = queue[1:]
	fmt.Println(current, queue)
}`
	return []LessonStep{
		{ID: "c173.mission-step", Title: ch17Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch17Markdown("c173.mission",
			"1. وصّل للجيود في layer 2 بأقل خطوات عبر الأنفاق والآبار.\n\n2. هدف الذهب بيتحدد حسب الخريطة المعروضة.\n\n3. الـstarter حمّل خرائط الطبقات وحدد مكان الهدف.",
			"1. Reach the layer 2 geode in the fewest steps through tunnels and shafts.\n\n2. The gold target depends on the displayed map.\n\n3. The starter loaded layer maps and located the target.", english)}},
		{ID: "c173.idea-step", Title: ch17Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch17Walkthrough("c173.cell-key",
			"BFS من الفصل 15 لسه هو نفسه، بس node بقت فيها z كمان.\n\nstruct حقولها أرقام ينفع تبقى map key؛ نفس x وy على طبقة مختلفة node تانية.",
			"The chapter 15 BFS is still the same, but a node now includes z.\n\nA struct of numbers can be a map key; the same x and y on another layer is a different node.",
			keys, "e173-cell-key", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 1, ar: "Cell فيها 3 إحداثيات.", en: "Cell has three coordinates."},
				{line: 3, ar: "map تبدأ فاضية.", en: "The map starts empty.", vars: []LessonVariable{{Name: "seen", Value: "{}"}}},
				{line: 4, ar: "سجلنا خلية في layer 1.", en: "Record a cell on layer 1.", vars: []LessonVariable{{Name: "p", Value: "{1 0 1}"}, {Name: "seen[p]", Value: "true"}}},
				{line: 5, ar: "نفس الإحداثيات تلاقيها.", en: "The same coordinates are found.", vars: []LessonVariable{{Name: "seen[1,0,1]", Value: "true"}}, output: "true"},
				{line: 6, ar: "طبقة تانية مفتاح مختلف.", en: "Another layer is a different key.", vars: []LessonVariable{{Name: "seen[1,0,2]", Value: "false"}}, output: "true\nfalse"},
			}, english)}},
		{ID: "c173.try-step", Title: ch17Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch17Question("c173.predict-queue", "إيه اللي هيتطبع بعد إزالة أول عنصر؟", "What prints after removing the first item?", queue, []string{"dock [shaft]", "shaft [dock]", "dock [dock shaft]"}, []string{"dock [shaft]", "shaft [dock]", "dock [dock shaft]"}, 0,
				[]string{"الطابور FIFO: خُد الأول وسيب الباقي.", "shaft لسه في الطابور.", "الإزالة قصّت العنصر الأول."},
				[]string{"The queue is FIFO: take the first and keep the rest.", "shaft remains in the queue.", "Removal sliced off the first item."}, english),
			ch17Walkthrough("c173.queue-reveal", "تابع حركة أول عنصر في الطابور.", "Follow the first item leaving the queue.", queue, "e173-queue", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 2, ar: "الدور: dock ثم shaft.", en: "Order: dock, then shaft.", vars: []LessonVariable{{Name: "queue", Value: "[dock shaft]"}}},
				{line: 3, ar: "اختارنا أول محطة.", en: "Select the first station.", vars: []LessonVariable{{Name: "current", Value: "dock"}}},
				{line: 4, ar: "باقي shaft في الطابور.", en: "shaft remains in the queue.", vars: []LessonVariable{{Name: "queue", Value: "[shaft]"}}},
				{line: 5, ar: "يطبع المحطة والباقي.", en: "Print the item and remainder.", vars: []LessonVariable{{Name: "current", Value: "dock"}, {Name: "queue", Value: "[shaft]"}}, output: "dock [shaft]"},
			}, english),
		}},
		{ID: "c173.watch-step", Title: ch17Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch17Markdown("c173.watch",
			"BFS يضمن أقصر طريق لما كل edge تكلفته واحدة.\n\nتحت الأرض الحركة مسموحة في Tunnel أو Shaft؛ والانتقال بين طبقات محتاج shaft.",
			"BFS guarantees a shortest route when every edge has equal cost.\n\nUnderground, movement is allowed on Tunnel or Shaft; layer changes require a shaft.", english)}},
		{ID: "c173.check-step", Title: ch17Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch17Markdown("c173.plan", "ارسم الجيران المسموحين لكل node، وبعدها راجع أي طبقة أو اتجاه اتغير.\n\nاحسب خطوة الحفر الأخيرة ضمن الخطة.", "List allowed neighbors for each node, then check which layer or direction changed.\n\nInclude the final dig action in the plan.", english),
			ch17Choice("c173.node-question", "نفس x وy مع z مختلف؛ نفس node؟", "Same x and y but different z: same node?", []string{"لأ، دي node مختلفة", "أيوه، z مش مهم", "حسب ترتيب الـmap"}, []string{"No, that is a different node", "Yes, z does not matter", "It depends on map order"}, 0,
				[]string{"z جزء من قيمة Cell ومن مفتاح الـmap.", "الأبعاد التلاتة بيميزوا المكان.", "ترتيب الـmap مش بيحدد الهوية."},
				[]string{"z is part of the Cell value and map key.", "All three dimensions identify the location.", "Map order does not define identity."}, english),
		}},
		{ID: "c173.turn-step", Title: ch17Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch17Markdown("c173.turn",
			"ابدأ من BFS اللي تعرفه، وسجّل z مع كل نقطة.\n\nشغّل الكود وراجع المسار مع الخريطة قبل الحفر.",
			"Start from the BFS you know and record z with each point.\n\nRun the code and compare the route with the map before digging.", english)}},
	}
}

func ch17Boss(english bool) []LessonStep {
	genericMap := `func mark[T comparable](seen map[T]bool, key T) {
	seen[key] = true
}
func main() {
	seen := map[string]bool{}
	mark(seen, "dock")
	fmt.Println(seen["dock"])
}`
	cellMap := `type Cell struct { X, Z int }
func main() {
	known := map[Cell]bool{{X: 2, Z: 1}: true}
	fmt.Println(known[Cell{X: 2, Z: 1}], known[Cell{X: 2, Z: 2}])
}`
	return []LessonStep{
		{ID: "c17b.mission-step", Title: ch17Copy("المهمة", "Mission", english), Blocks: []LessonBlock{ch17Markdown("c17b.mission",
			"1. المختبر طالب ~RelayPath[T comparable]~ لمسارات موجهة.\n\n2. راجع عقد النتيجة والفشل والملكية والدورات.\n\n3. اكتب tests وbenchmark؛ الكودين ظاهرين في الـstarter.",
			"1. The lab asks for ~RelayPath[T comparable]~ over directed routes.\n\n2. Review result, failure, ownership, and cycle promises.\n\n3. Write tests and a benchmark; both files are in the starter.", english)}},
		{ID: "c17b.idea-step", Title: ch17Copy("الفكرة", "Idea", english), Blocks: []LessonBlock{ch17Walkthrough("c17b.comparable-map",
			"~comparable~ بتسمح باستخدام T كمفتاح map.\n\nالمثال يسجل مفتاح نوعه string؛ في العقد النوع بيتغير، لكن خاصية المقارنة تفضل مطلوبة.",
			"~comparable~ allows T to be a map key.\n\nThis example records a string key; the contract can change its type, but comparability is still required.",
			genericMap, "e17b-comparable", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 1, ar: "T لازم تقبل المقارنة.", en: "T must support comparison."},
				{line: 5, ar: "map string تبدأ فاضية.", en: "The string map starts empty.", vars: []LessonVariable{{Name: "seen", Value: "{}"}}},
				{line: 6, ar: "سجلنا dock.", en: "Record dock.", vars: []LessonVariable{{Name: "seen[dock]", Value: "true"}}},
				{line: 7, ar: "المفتاح موجود.", en: "The key is present.", vars: []LessonVariable{{Name: "seen[dock]", Value: "true"}}, output: "true"},
			}, english)}},
		{ID: "c17b.try-step", Title: ch17Copy("جرّب", "Try it", english), Blocks: []LessonBlock{
			ch17Question("c17b.predict-cell", "أنهي lookup يرجّع false؟", "Which lookup returns false?", cellMap, []string{"اللي z فيها 2", "اللي z فيها 1", "الاتنين true"}, []string{"The one with z equal to 2", "The one with z equal to 1", "Both are true"}, 0,
				[]string{"المفتاح المخزن فيه z=1.", "ده نفس المفتاح اللي اتخزن.", "تغيير z يغيّر قيمة المفتاح."},
				[]string{"The stored key has z=1.", "That is the stored key.", "Changing z changes the key value."}, english),
			ch17Walkthrough("c17b.cell-reveal", "قارن مفاتيح الـstruct كاملة.", "Compare the complete struct keys.", cellMap, "e17b-cell", []string{"fmt"}, []ch17FrameCopy{
				{line: 0, ar: "قبل التنفيذ.", en: "Before execution."},
				{line: 1, ar: "Cell فيها X وZ.", en: "Cell has X and Z."},
				{line: 3, ar: "الخريطة فيها z=1.", en: "The map stores z=1.", vars: []LessonVariable{{Name: "known", Value: "{{2 1}:true}"}}},
				{line: 4, ar: "الأول true، والتاني false.", en: "First is true; second is false.", vars: []LessonVariable{{Name: "known[{2 1}]", Value: "true"}, {Name: "known[{2 2}]", Value: "false"}}, output: "true false"},
			}, english),
		}},
		{ID: "c17b.watch-step", Title: ch17Copy("خد بالك", "Watch out", english), Blocks: []LessonBlock{ch17Markdown("c17b.watch",
			"العقد بيوصف مسار موجه؛ الجار مش لازم يرجعلك.\n\nالنجاح الوهمي للوجهة الغايبة يخرق العقد، و~nil~ غير slice فاضية.",
			"The contract describes directed routes; a neighbor need not link back.\n\nInventing success for a missing destination breaks the contract, and ~nil~ differs from an empty slice.", english)}},
		{ID: "c17b.check-step", Title: ch17Copy("فكّر واختبر", "Think & check", english), Blocks: []LessonBlock{
			ch17Markdown("c17b.plan", "حوّل كل وعد لملاحظة مستقلة: قِصر الطريق، الإشارة عند الفشل، عدم تغيير المدخل، وانتهاء البحث مع الدورات.\n\nقِس الأداء لوحده.", "Turn each promise into a separate observation: shortest route, failure signal, unchanged input, and termination with cycles.\n\nMeasure performance separately.", english),
			ch17Choice("c17b.comparable-question", "إيه اللي بتضمنه ~comparable~؟", "What does ~comparable~ guarantee?", []string{"T ينفع مفتاح map", "كل T نوع رقمي", "كل القيم مرتبة"}, []string{"T can be a map key", "Every T is numeric", "All values are ordered"}, 0,
				[]string{"المقارنة مطلوبة للمفتاح؛ مش لازم يكون رقم.", "string وstruct مناسبين كمان لو حقولهم قابلة للمقارنة.", "comparable ما بتفرضش ترتيب."},
				[]string{"Keys need comparison; they need not be numeric.", "Strings and comparable structs can work too.", "comparable does not define an order."}, english),
		}},
		{ID: "c17b.turn-step", Title: ch17Copy("دورك", "Your turn", english), Blocks: []LessonBlock{ch17Markdown("c17b.turn",
			"اكتب قائمة بالعقد جنب placeholder التست، وبعدها طبّق البحث وbenchmark.\n\nراجع النتيجة والملكية قبل ما تقيس السرعة.",
			"Write a contract checklist beside the test placeholder, then implement the search and benchmark.\n\nReview results and ownership before measuring speed.", english)}},
	}
}
