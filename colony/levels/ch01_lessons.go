package levels

import "strings"

func ch01Paragraphs(text string) string {
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

func ch01Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: ch01Paragraphs(text)}
}

func ch01Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: ch01Paragraphs(text), Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch01Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: ch01Paragraphs(text), Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch01Variable(name, value string) LessonVariable {
	return LessonVariable{Name: name, Value: value}
}

func ch01TextPair(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

func ch01Sandbox(id, text, kind string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: ch01Paragraphs(text), Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

var ch01SignalSandboxConfig = map[string]any{
	"labels": map[string]any{
		"position": ch01TextPair("المكان", "Position"),
		"command":  ch01TextPair("الأمر", "Command"),
		"previous": ch01TextPair("السابق", "Previous"),
		"next":     ch01TextPair("التالي", "Next"),
		"reset":    ch01TextPair("من الأول", "Reset"),
	},
	"frames": []map[string]any{
		{"x": 1, "y": 1, "command": ch01TextPair("بداية", "Start"), "caption": ch01TextPair("لسه مفيش أمر.", "No command yet.")},
		{"x": 2, "y": 1, "command": ch01TextPair("East", "East"), "caption": ch01TextPair("East يغيّر x بخانة.", "East changes x by one.")},
		{"x": 2, "y": 2, "command": ch01TextPair("North", "North"), "caption": ch01TextPair("North يغيّر y بخانة.", "North changes y by one.")},
		{"x": 1, "y": 2, "command": ch01TextPair("West", "West"), "caption": ch01TextPair("كل أمر ييجي بالترتيب.", "Each command follows in order.")},
	},
}

var ch01HarvestSandboxConfig = map[string]any{
	"labels": map[string]any{
		"command":   ch01TextPair("الأمر", "Command"),
		"drone":     ch01TextPair("الدرون", "Drone"),
		"plant":     ch01TextPair("طحلب جاهز", "Moss ready"),
		"collected": ch01TextPair("اتجمع", "Collected"),
		"harvested": ch01TextPair("اللي اتجمع", "Harvested"),
		"previous":  ch01TextPair("السابق", "Previous"),
		"next":      ch01TextPair("التالي", "Next"),
		"reset":     ch01TextPair("من الأول", "Reset"),
	},
	"tiles": []map[string]any{
		{"label": ch01TextPair("محطة", "Dock"), "plant": false},
		{"label": ch01TextPair("قطعة طحلب", "Moss plot"), "plant": true},
		{"label": ch01TextPair("صخرة", "Rock"), "plant": false},
	},
	"frames": []map[string]any{
		{"position": 0, "harvested": []int{}, "action": ch01TextPair("استعداد", "Ready"), "caption": ch01TextPair("كل حاجة بتبدأ من المكان الحالي.", "Everything starts at the current tile.")},
		{"position": 1, "harvested": []int{}, "action": ch01TextPair("Move East", "Move East"), "caption": ch01TextPair("الحركة تسبق استدعاء الدالة.", "Movement comes before the function call.")},
		{"position": 1, "harvested": []int{1}, "action": ch01TextPair("Harvest()", "Harvest()"), "caption": ch01TextPair("Harvest تستخدم مكان الدرون الحالي.", "Harvest uses the drone's current tile.")},
	},
}

var ch01PlantingSandboxConfig = map[string]any{
	"labels": map[string]any{
		"ground":   ch01TextPair("الأرض", "Ground"),
		"plant":    ch01TextPair("النبتة", "Plant"),
		"action":   ch01TextPair("الأمر", "Command"),
		"empty":    ch01TextPair("مفيش نبتة", "No plant"),
		"unready":  ch01TextPair("لسه متجهزتش", "Unprepared"),
		"ready":    ch01TextPair("متجهزة", "Prepared"),
		"planted":  ch01TextPair("طحلب", "Moss"),
		"previous": ch01TextPair("السابق", "Previous"),
		"next":     ch01TextPair("التالي", "Next"),
		"reset":    ch01TextPair("من الأول", "Reset"),
	},
	"frames": []map[string]any{
		{"ground": "unready", "plant": "empty", "action": ch01TextPair("البداية", "Start"), "caption": ch01TextPair("الأرض لسه محتاجة تجهيز.", "The ground still needs preparation.")},
		{"ground": "ready", "plant": "empty", "action": ch01TextPair("Till", "Till"), "caption": ch01TextPair("Till يجهز حالة الأرض.", "Till changes the ground state.")},
		{"ground": "ready", "plant": "planted", "action": ch01TextPair("Plant", "Plant"), "caption": ch01TextPair("Plant تيجي بعد تجهيز الأرض.", "Plant follows prepared ground.")},
	},
}

var ch01FirstSignalArabic = []LessonStep{
	{ID: "c11.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch01Markdown("c11.mission", "1. وصّل الدرون للمنارة 📡.\n2. الدرون بيبدأ عند (0,0)، والمنارة شرقًا.\n3. شغّل البرنامج وراقب مكانه بعد كل أمر.\n\nالـstarter بدأ الحركة بالفعل؛ كمل من السطر الموجود.")}},
	{ID: "c11.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch01Walkthrough("c11.sequence", "Go بتنّفذ أوامر ~main~ بالترتيب من فوق لتحت. كل استدعاء دالة أمر مستقل، وقيمة الاتجاه بتحدد وجهته. المثال للتوضيح وبيطبع رسالتين بالترتيب.",
		`func main() {
	fmt.Println("drone ready")
	fmt.Println("beacon ahead")
}`,
		"e11-sequence", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "بدأت دالة main."}, {Line: 2, Caption: "الرسالة الأولى تظهر.", Output: "drone ready"}, {Line: 3, Caption: "الثانية تتبع الأولى.", Output: "drone ready\nbeacon ahead"}})}},
	{ID: "c11.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch01Question("c11.predict-order", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tfmt.Println(\"step\")\n\tfmt.Println(\"signal\")\n}\n~~~", []string{"signal ثم step", "step ثم signal", "السطران في نفس الوقت"}, 1, []string{"Go لا تعكس ترتيب السطور.", "السطر الأعلى يطبع أولًا.", "الدالتان هنا بتتنفذا بالتتابع."}),
		ch01Walkthrough("c11.order-reveal", "تابع السطرين وشوف ترتيب ظهور الرسائل.",
			`func main() {
	fmt.Println("step")
	fmt.Println("signal")
}`,
			"e11-order", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 2, Caption: "اطبع الأمر الأول.", Output: "step"}, {Line: 3, Caption: "اطبع الإشارة بعده.", Output: "step\nsignal"}}),
		ch01Sandbox("ch01-first-signal", "تابع مكان الدرون مع كل أمر.", "ch01-route-stepper", ch01SignalSandboxConfig),
	}},
	{ID: "c11.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch01Markdown("c11.watch", "- ~East~ معناها يمين الخريطة؛ الاتجاه مش عدد الخانات.\n- كل استدعاء ~Move~ بيحاول خطوة واحدة بس.")}},
	{ID: "c11.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch01Markdown("c11.plan", "راجع إحداثي الدرون بعد كل حركة. قارن مكانه بإحداثي الهدف، مش بالرسالة في الـConsole.\n\nابدأ من أثر السطر الموجود قبل إضافة أوامر جديدة."),
		ch01Question("c11.direction-check", "لو الدرون عند (0,0)، أي اتجاه يقرّبه من هدف شرقًا؟", []string{"North", "East", "West"}, 1, []string{"North يحرّك لأعلى.", "East يحرّك لليمين.", "West يحرّك لليسار."}),
	}},
	{ID: "c11.turn-step", Title: "دورك", Blocks: []LessonBlock{ch01Markdown("c11.turn", "كمّل البرنامج وشغّله. راقب حركة الدرون خطوة خطوة، وبعدها قارن مكانه بالمنارة.")}},
}

var ch01FirstSignalEnglish = []LessonStep{
	{ID: "c11.mission-step", Title: "Mission", Blocks: []LessonBlock{ch01Markdown("c11.mission", "1. Guide the drone to the beacon 📡.\n2. It starts at (0,0); the beacon is east.\n3. Run the program and watch its position after each command.\n\nThe starter already begins moving; continue from its existing line.")}},
	{ID: "c11.idea-step", Title: "Idea", Blocks: []LessonBlock{ch01Walkthrough("c11.sequence", "Go runs commands in ~main~ from top to bottom. Each function call is a separate command, and its direction value chooses where it points. This example prints two messages in order.",
		`func main() {
	fmt.Println("drone ready")
	fmt.Println("beacon ahead")
}`,
		"e11-sequence", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "The main function starts."}, {Line: 2, Caption: "The first message appears.", Output: "drone ready"}, {Line: 3, Caption: "The next message follows.", Output: "drone ready\nbeacon ahead"}})}},
	{ID: "c11.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch01Question("c11.predict-order", "Predict the output:\n\n~~~go\nfunc main() {\n\tfmt.Println(\"step\")\n\tfmt.Println(\"signal\")\n}\n~~~", []string{"signal then step", "step then signal", "Both lines at once"}, 1, []string{"Go does not reverse the lines.", "The higher line prints first.", "These calls run one after another."}),
		ch01Walkthrough("c11.order-reveal", "Follow both lines and watch when each message appears.",
			`func main() {
	fmt.Println("step")
	fmt.Println("signal")
}`,
			"e11-order", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 2, Caption: "Print the first command.", Output: "step"}, {Line: 3, Caption: "Print the signal next.", Output: "step\nsignal"}}),
		ch01Sandbox("ch01-first-signal", "Follow the drone's position after each command.", "ch01-route-stepper", ch01SignalSandboxConfig),
	}},
	{ID: "c11.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch01Markdown("c11.watch", "- ~East~ means right on the map; it is not a distance.\n- Each ~Move~ call tries just one step.")}},
	{ID: "c11.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch01Markdown("c11.plan", "Check the drone's coordinates after each move. Compare its position with the target, not with a Console message.\n\nStart by accounting for the line already in the starter."),
		ch01Question("c11.direction-check", "At (0,0), which direction gets closer to a target east?", []string{"North", "East", "West"}, 1, []string{"North moves upward.", "East moves right.", "West moves left."}),
	}},
	{ID: "c11.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch01Markdown("c11.turn", "Continue the program and run it. Watch the drone move one step at a time, then compare its position with the beacon.")}},
}

var ch01FirstHarvestArabic = []LessonStep{
	{ID: "c12.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch01Markdown("c12.mission", "1. اجمع الطحالب 🌿 من الخانات اللي على اليمين.\n2. لازم الدرون يقف فوق النبات وقت الحصاد.\n3. راقب العداد بعد التشغيل.\n\nالـstarter فيه حركة وحصاد بالفعل؛ افهم ترتيبهم قبل ما تكمل.")}},
	{ID: "c12.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch01Walkthrough("c12.function-call", "الدالة ممكن تتنادى من غير arguments، بس الأقواس الفاضية لسه مطلوبة. استدعاء الدالة جوه ~main~ يشغّل جسمها. المثال بيوضح دالة تدريب بسيطة.",
		`func harvest() {
	fmt.Println("harvest")
}
func main() {
	harvest()
}`,
		"e12-call", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 1, Caption: "تعريف دالة الحصاد."}, {Line: 5, Caption: "main تستدعي harvest."}, {Line: 2, Caption: "الدالة تطبع رسالة.", Output: "harvest"}})}},
	{ID: "c12.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch01Question("c12.predict-order", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tfmt.Println(\"move\")\n\tfmt.Println(\"harvest\")\n}\n~~~", []string{"harvest ثم move", "move ثم harvest", "مفيش ناتج"}, 1, []string{"الترتيب مش بيتعكس.", "الحركة تظهر قبل الحصاد.", "كل استدعاء يطبع سطرًا."}),
		ch01Walkthrough("c12.order-reveal", "شوف إزاي ترتيب الاستدعاءين بيظهر في الناتج.",
			`func main() {
	fmt.Println("move")
	fmt.Println("harvest")
}`,
			"e12-order", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 2, Caption: "نفّذ الحركة أولًا.", Output: "move"}, {Line: 3, Caption: "الحصاد يأتي بعدها.", Output: "move\nharvest"}}),
		ch01Sandbox("ch01-first-harvest", "حرّك الدرون ثم تابع أثر ~Harvest()~.", "ch01-harvest-call", ch01HarvestSandboxConfig),
	}},
	{ID: "c12.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch01Markdown("c12.watch", "- ~Harvest()~ مالهاش argument، لكن لازم تكتب ~()~.\n- الحصاد بيحصل في مكان الدرون الحالي وقت الاستدعاء.")}},
	{ID: "c12.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch01Markdown("c12.plan", "افصل بين الحركة والفعل: اسأل نفسك الدرون واقف فين قبل كل ~Harvest~.\n\nالعداد بيقيس النباتات اللي اتجمعت، مش عدد الأسطر."),
		ch01Question("c12.harvest-check", "إيه اللي يحدد الخانة اللي هيحصل عليها الحصاد؟", []string{"مكان الدرون", "اسم الدالة", "عدد الرسائل"}, 0, []string{"الفعل يستخدم المكان الحالي.", "اسم الدالة لا يغيّر المكان.", "الرسائل لا تحدد موقع الدرون."}),
	}},
	{ID: "c12.turn-step", Title: "دورك", Blocks: []LessonBlock{ch01Markdown("c12.turn", "شغّل محاولتك وراقب كل حركة وحصاد في الـConsole. قارن مكان الدرون بعد كل أمر بالطحلب التالي.")}},
}

var ch01FirstHarvestEnglish = []LessonStep{
	{ID: "c12.mission-step", Title: "Mission", Blocks: []LessonBlock{ch01Markdown("c12.mission", "1. Collect the moss 🌿 on the tiles to the right.\n2. The drone must stand on a plant to harvest it.\n3. Check the counter after running.\n\nThe starter already has a move and harvest; understand their order before continuing.")}},
	{ID: "c12.idea-step", Title: "Idea", Blocks: []LessonBlock{ch01Walkthrough("c12.function-call", "A function can be called without arguments, but its empty parentheses are still required. Calling it inside ~main~ runs its body. This example shows a small practice function.",
		`func harvest() {
	fmt.Println("harvest")
}
func main() {
	harvest()
}`,
		"e12-call", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 1, Caption: "Define the harvest function."}, {Line: 5, Caption: "main calls harvest."}, {Line: 2, Caption: "The function prints a message.", Output: "harvest"}})}},
	{ID: "c12.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch01Question("c12.predict-order", "Predict the output:\n\n~~~go\nfunc main() {\n\tfmt.Println(\"move\")\n\tfmt.Println(\"harvest\")\n}\n~~~", []string{"harvest then move", "move then harvest", "No output"}, 1, []string{"The order is not reversed.", "Move appears before harvest.", "Each call prints one line."}),
		ch01Walkthrough("c12.order-reveal", "Watch how the two calls appear in the output.",
			`func main() {
	fmt.Println("move")
	fmt.Println("harvest")
}`,
			"e12-order", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 2, Caption: "Run the move first.", Output: "move"}, {Line: 3, Caption: "Harvest comes next.", Output: "move\nharvest"}}),
		ch01Sandbox("ch01-first-harvest", "Move the drone, then watch ~Harvest()~.", "ch01-harvest-call", ch01HarvestSandboxConfig),
	}},
	{ID: "c12.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch01Markdown("c12.watch", "- ~Harvest()~ has no argument, but it still needs ~()~.\n- Harvest happens where the drone is when the call runs.")}},
	{ID: "c12.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch01Markdown("c12.plan", "Separate movement from the action: ask where the drone stands before each ~Harvest~.\n\nThe counter measures collected plants, not the number of code lines."),
		ch01Question("c12.harvest-check", "What determines the tile where harvesting happens?", []string{"The drone's position", "The function name", "The message count"}, 0, []string{"The action uses the current tile.", "The name does not change position.", "Messages do not set the drone's location."}),
	}},
	{ID: "c12.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch01Markdown("c12.turn", "Run your attempt and watch each move and harvest in the Console. Compare the drone's position after each command with the next moss plant.")}},
}

var ch01PlantingArabic = []LessonStep{
	{ID: "c13.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch01Markdown("c13.mission", "1. ازرع طحالب في الخانات التلاتة شرق البداية.\n2. كل خانة محتاجة أرض جاهزة ونبتة.\n3. تابع عداد الخانات المزروعة.\n\nالـstarter حرث أول خانة بالفعل؛ كمّل من الحالة دي.")}},
	{ID: "c13.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch01Walkthrough("c13.order", "ترتيب الأوامر مهم لما خطوة تجهز اللي بعدها. في المثال بنطبع التحضير قبل الزراعة عشان نتابع التسلسل؛ البيانات هنا للتوضيح.",
		`func main() {
	fmt.Println("till soil")
	fmt.Println("plant moss")
}`,
		"e13-order", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التنفيذ."}, {Line: 2, Caption: "جهّز الأرض أولًا.", Output: "till soil"}, {Line: 3, Caption: "ازرع بعد التجهيز.", Output: "till soil\nplant moss"}})}},
	{ID: "c13.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch01Question("c13.predict-order", "توقّع الناتج:\n\n~~~go\nfunc main() {\n\tfmt.Println(\"prepare\")\n\tfmt.Println(\"plant\")\n}\n~~~", []string{"plant ثم prepare", "prepare ثم plant", "plant بس"}, 1, []string{"السطر الأول يشتغل أولًا.", "التحضير يظهر قبل الزراعة.", "الاستدعاء الأول كمان يطبع."}),
		ch01Walkthrough("c13.reveal", "تابع التحضير وبعده أمر الزراعة.",
			`func main() {
	fmt.Println("prepare")
	fmt.Println("plant")
}`,
			"e13-try", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "قبل التشغيل."}, {Line: 2, Caption: "التحضير في الأول.", Output: "prepare"}, {Line: 3, Caption: "الزراعة تأتي بعدها.", Output: "prepare\nplant"}}),
		ch01Sandbox("ch01-planting", "تابع حالة الأرض قبل الزراعة.", "ch01-planting-state", ch01PlantingSandboxConfig),
	}},
	{ID: "c13.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch01Markdown("c13.watch", "- الطحلب محتاج ~soil~؛ ~regolith~ لسه محتاج تجهيز.\n- الـstarter حرث خانة البداية، فحالتها مختلفة عن الخانات التالية.")}},
	{ID: "c13.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch01Markdown("c13.plan", "لكل خانة، راجع الأرض والنبتة بعد ما يوصلها الدرون. لو وقفت عند خانة جديدة، اسأل إيه اللي اتغيّر وإيه اللي لسه محتاج خطوة."),
		ch01Question("c13.ground-check", "إيه اللي لازم يكون جاهز قبل ما الطحلب ينمو؟", []string{"Soil", "رسالة في Console", "إحداثي جديد"}, 0, []string{"الطحلب يحتاج أرضًا من نوع soil.", "الرسالة لا تجهز الأرض.", "تغيير الإحداثي وحده مش كفاية."}),
	}},
	{ID: "c13.turn-step", Title: "دورك", Blocks: []LessonBlock{ch01Markdown("c13.turn", "شغّل محاولتك وراقب حالة الأرض والنبات في كل خانة. استخدم عداد المهمة لتعرف التقدم.")}},
}

var ch01PlantingEnglish = []LessonStep{
	{ID: "c13.mission-step", Title: "Mission", Blocks: []LessonBlock{ch01Markdown("c13.mission", "1. Plant moss on the three tiles east of the start.\n2. Each tile needs prepared ground and a plant.\n3. Follow the planted-tile counter.\n\nThe starter already tilled the first tile; continue from that state.")}},
	{ID: "c13.idea-step", Title: "Idea", Blocks: []LessonBlock{ch01Walkthrough("c13.order", "Command order matters when one step prepares the next. This example prints preparation before planting so we can follow the sequence; the data is illustrative.",
		`func main() {
	fmt.Println("till soil")
	fmt.Println("plant moss")
}`,
		"e13-order", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before execution."}, {Line: 2, Caption: "Prepare the ground first.", Output: "till soil"}, {Line: 3, Caption: "Plant after preparing.", Output: "till soil\nplant moss"}})}},
	{ID: "c13.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch01Question("c13.predict-order", "Predict the output:\n\n~~~go\nfunc main() {\n\tfmt.Println(\"prepare\")\n\tfmt.Println(\"plant\")\n}\n~~~", []string{"plant then prepare", "prepare then plant", "Only plant"}, 1, []string{"The first line runs first.", "Preparation appears before planting.", "The first call also prints."}),
		ch01Walkthrough("c13.reveal", "Follow preparation and then the planting command.",
			`func main() {
	fmt.Println("prepare")
	fmt.Println("plant")
}`,
			"e13-try", []string{"fmt"}, []LessonFrame{{Line: 0, Caption: "Before running."}, {Line: 2, Caption: "Prepare comes first.", Output: "prepare"}, {Line: 3, Caption: "Planting comes next.", Output: "prepare\nplant"}}),
		ch01Sandbox("ch01-planting", "Watch the ground before planting.", "ch01-planting-state", ch01PlantingSandboxConfig),
	}},
	{ID: "c13.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch01Markdown("c13.watch", "- Moss needs ~soil~; ~regolith~ still needs preparation.\n- The starter tilled the starting tile, so it differs from later tiles.")}},
	{ID: "c13.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch01Markdown("c13.plan", "For each tile, check the ground and plant after the drone reaches it. When you arrive at a new tile, ask what changed and what still needs a step."),
		ch01Question("c13.ground-check", "What must be ready before moss can grow?", []string{"Soil", "A Console message", "A new coordinate"}, 0, []string{"Moss needs soil ground.", "A message does not prepare ground.", "Changing position alone is not enough."}),
	}},
	{ID: "c13.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch01Markdown("c13.turn", "Run your attempt and watch the ground and plant on each tile. Use the mission counter to track progress.")}},
}
