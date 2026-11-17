package levels

func lateText(arabic, english string, useEnglish bool) string {
	if useEnglish {
		return english
	}
	return arabic
}

func lateQuestion(id, arabic, english string, optionsArabic, optionsEnglish []string, answer int, feedbackArabic, feedbackEnglish []string, useEnglish bool) LessonBlock {
	return LessonBlock{
		ID: id, Kind: "question", Text: lateText(arabic, english, useEnglish),
		Question: &LessonQuestion{Options: lateTextSlice(optionsArabic, optionsEnglish, useEnglish), Answer: answer, Feedback: lateTextSlice(feedbackArabic, feedbackEnglish, useEnglish)},
	}
}

func lateTextSlice(arabic, english []string, useEnglish bool) []string {
	if useEnglish {
		return english
	}
	return arabic
}

const stockOrderComparatorExample = `type Player struct { name string; age int }
players := []Player{{"Mira", 8}, {"Ari", 12}, {"Ari", 9}}
before := func(a, b Player) int {
	if byName := cmp.Compare(a.name, b.name); byName != 0 {
		return byName
	}
	return cmp.Compare(a.age, b.age)
}
fmt.Println(before(players[1], players[0]))
fmt.Println(before(players[1], players[2]))
slices.SortFunc(players, before)
fmt.Println(players)`

func stockOrderLesson(useEnglish bool) []LessonStep {
	return []LessonStep{
		{
			ID: "c64.late.mission-step", Title: lateText("المهمة", "Mission", useEnglish),
			Blocks: []LessonBlock{{ID: "c64.late.mission", Kind: "markdown", Text: lateText(
				"1. رتّب المواقع: الحجم الأكبر أولًا، ثم x الأصغر عند التعادل.\n2. اقرأ العرض بـ~colony.WorldSize~ وابعث المواقع بتقرير واحد.\n3. الـstarter طبع قياسًا من غير حفظ؛ كلف tick. ممنوع الحصاد.\n\nالذهب: **19 ticks أو أقل**.",
				"1. Rank positions: larger sizes first, then smaller x on ties.\n2. Read width with ~colony.WorldSize~; report positions once.\n3. The starter printed an unsaved measurement (one tick); do not harvest.\n\nGold: **19 ticks or fewer**.", useEnglish)}},
		},
		{
			ID: "c64.late.idea-step", Title: lateText("الفكرة", "Idea", useEnglish),
			Blocks: []LessonBlock{{
				ID: "c64.late.trace", Kind: "walkthrough",
				Text: lateText("المقارنة بتقول مين يسبق: السالب يقدّم a، والصفر تعادل، والموجب يأخّره.\n\nالمثال أسماء ثم أعمار؛ قاعدة المهمة مختلفة.", "A comparator returns a sign: negative puts a first, zero ties, positive puts it later.\n\nThis toy sorts names, then ages; the mission differs.", useEnglish),
				Code: stockOrderComparatorExample, Imports: []string{"cmp", "fmt", "slices"}, ExampleID: "e64-sort-sign-late",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: lateText("قبل الكروت.", "Before the cards.", useEnglish)},
					{Line: 1, Caption: lateText("نوع بيانات للتجربة.", "An illustrative type.", useEnglish), Variables: []LessonVariable{{Name: "Player", Value: "name, age"}}},
					{Line: 2, Caption: lateText("أسماء وأعمار للتوضيح.", "Illustrative names and ages.", useEnglish), Variables: []LessonVariable{{Name: "players", Value: "[{Mira 8} {Ari 12} {Ari 9}]"}}},
					{Line: 3, Caption: lateText("قاعدة المقارنة محفوظة.", "Comparator saved.", useEnglish), Variables: []LessonVariable{{Name: "before", Value: "function value"}}},
					{Line: 9, Caption: lateText("Ari يسبق Mira؛ النتيجة -1.", "Ari precedes Mira; result -1.", useEnglish), Variables: []LessonVariable{{Name: "byName", Value: "-1"}}, Output: "-1"},
					{Line: 10, Caption: lateText("الاسم اتعادل؛ العمر حسمها.", "Names tie; age decides.", useEnglish), Variables: []LessonVariable{{Name: "byName", Value: "0"}, {Name: "byAge", Value: "1"}}, Output: "-1\n1"},
					{Line: 11, Caption: lateText("SortFunc رتّبت القائمة.", "SortFunc reorders the slice.", useEnglish), Variables: []LessonVariable{{Name: "players", Value: "[{Ari 9} {Ari 12} {Mira 8}]"}}, Output: "-1\n1"},
					{Line: 12, Caption: lateText("الأسماء ثم الأعمار.", "Names, then ages.", useEnglish), Variables: []LessonVariable{{Name: "players", Value: "[{Ari 9} {Ari 12} {Mira 8}]"}}, Output: "-1\n1\n[{Ari 9} {Ari 12} {Mira 8}]"},
				}},
			}},
		},
		{
			ID: "c64.late.try-step", Title: lateText("جرّب", "Try it", useEnglish),
			Blocks: []LessonBlock{{
				ID: "c64.late.sort-sandbox", Kind: "sandbox",
				Text: lateText("بدّل اتجاه الحجم وكسر التعادل؛ الكروت توضيحية والتعادلات مميزة.", "Toggle size direction and tie-breaking; cards are illustrative and ties are highlighted.", useEnglish),
				Sandbox: &LessonSandbox{
					Kind: "sort-compare", Sequence: []int{3, 5, 3, 5, 2},
					SourceLabel:        lateText("أحجام للتوضيح، مش بيانات المهمة.", "Illustrative sizes, not mission data.", useEnglish),
					SortDirectionLabel: lateText("اتجاه الحجم", "Size direction", useEnglish), SortTieLabel: lateText("كسر التعادل", "Tie-break", useEnglish),
					AscendingLabel: lateText("تصاعدي", "Ascending", useEnglish), DescendingLabel: lateText("تنازلي", "Descending", useEnglish),
					TieOnLabel: lateText("شغّال", "On", useEnglish), TieOffLabel: lateText("مقفول", "Off", useEnglish),
					ResetLabel: lateText("من الأول", "Reset", useEnglish),
				},
			}},
		},
		{
			ID: "c64.late.watch-step", Title: lateText("خد بالك", "Watch out", useEnglish),
			Blocks: []LessonBlock{{ID: "c64.late.watch", Kind: "markdown", Text: lateText(
				"- خلي نتيجة المقارنة متسقة: سالب أو صفر أو موجب.\n- ~slices.SortFunc~ مش stable؛ ~slices.SortStableFunc~ بتحافظ على ترتيب العناصر المتعادلة.\n\n~Measure~ و~Move~ (حتى لو فشلت) و~Harvest~: كل واحدة tick.\n\n~WorldSize~ و~Report~ والحساب بلا ticks.",
				"- Keep comparator results consistent: negative, zero, or positive.\n- ~slices.SortFunc~ is not stable; ~slices.SortStableFunc~ keeps equal items in input order.\n\n~Measure~, ~Move~ (even on failure), and ~Harvest~ each cost one tick.\n\n~WorldSize~, ~Report~, and computation cost none.", useEnglish)}},
		},
		{
			ID: "c64.late.check-step", Title: lateText("فكّر واختبر", "Think & check", useEnglish),
			Blocks: []LessonBlock{
				{ID: "c64.late.plan", Kind: "markdown", Text: lateText("قيس الكل → احتفظ بالقراءات → رتّب المواقع → ابعت تقريرًا واحدًا.", "Measure all → keep the readings → rank positions → send one report.", useEnglish)},
				lateQuestion("c64.late.q-sign", "-1 معناها إيه؟", "What does -1 mean?",
					[]string{"a يسبق", "تعادل", "a يتأخر"}, []string{"a comes first", "Tie", "a comes later"}, 0,
					[]string{"السالب يقدّم a.", "الصفر يعني تعادل.", "الموجب يأخر a."}, []string{"Negative puts a first.", "Zero means a tie.", "Positive puts a later."}, useEnglish),
				lateQuestion("c64.late.q-stable", "هل ~SortFunc~ تحفظ ترتيب المتعادلين؟", "Does ~SortFunc~ keep tied items in input order?",
					[]string{"دايمًا", "مش مضمون", "بس مع أرقام"}, []string{"Always", "Not guaranteed", "Only for numbers"}, 1,
					[]string{"SortFunc مش stable.", "صح؛ SortStableFunc تحفظ ترتيبهم.", "نوع العنصر مش فارق هنا."}, []string{"SortFunc is not stable.", "Right; SortStableFunc keeps their order.", "The element type does not matter here."}, useEnglish),
				lateQuestion("c64.late.q-report", "تقرير المهمة يبعت إيه؟", "What does the mission report?",
					[]string{"أحجام", "مواقع", "طباعة"}, []string{"Sizes", "Positions", "Printed output"}, 1,
					[]string{"المهمة عايزة المواقع.", "صح؛ ابعت المواقع المرتبة.", "الطباعة للتصحيح."}, []string{"The mission needs positions.", "Right; report ranked positions.", "Printing is for debugging."}, useEnglish),
			},
		},
		{
			ID: "c64.late.turn-step", Title: lateText("دورك", "Your turn", useEnglish),
			Blocks: []LessonBlock{{ID: "c64.late.turn", Kind: "markdown", Text: lateText("افتح الـstarter، عدّل حاجة صغيرة، وشغّل المهمة وراجع التقرير.", "Open the starter; make a small change and inspect its report.", useEnglish)}},
		},
	}
}

const heatGridExample = `grid := make([][]int, 2)
grid[0] = make([]int, 3)
grid[1] = make([]int, 3)
grid[0][1] = 4
grid[1][2] = 5
outer := slices.Clone(grid)
outer[0][1] = 9
rows := [][]int{slices.Clone(grid[0]), slices.Clone(grid[1])}
rows[0][1] = 7
fmt.Println(len(grid), len(grid[0]), grid[1][2])
fmt.Println(grid[0][1], outer[0][1], rows[0][1])`

func heatGridLesson(useEnglish bool) []LessonStep {
	return []LessonStep{
		{
			ID: "c65.late.mission-step", Title: lateText("المهمة", "Mission", useEnglish),
			Blocks: []LessonBlock{{ID: "c65.late.mission", Kind: "markdown", Text: lateText(
				"دي lab: زر Run بيشغّل tests للدوال، مش حركة للدرون. كمّل ~OverlayHeat~ و~ThroughHeat~ و~FlipHeat~؛ مفيش ~Report~ ولا ticks.",
				"Lab: Run tests, not drones. Implement ~OverlayHeat~, ~ThroughHeat~, and ~FlipHeat~; no ~Report~ or ticks.", useEnglish)}},
		},
		{
			ID: "c65.late.idea-step", Title: lateText("الفكرة", "Idea", useEnglish),
			Blocks: []LessonBlock{{
				ID: "c65.late.trace", Kind: "walkthrough",
				Text: lateText("الـ2D slice صفوف جوه slice. في ~grid[r][c]~، الأول row والتاني column؛ ~len(grid)~ للصفوف و~len(grid[0])~ للأعمدة.", "2D slices hold rows. In ~grid[r][c]~, row comes first, then column; ~len(grid)~ counts rows and ~len(grid[0])~ counts columns.", useEnglish),
				Code: heatGridExample, Imports: []string{"fmt", "slices"}, ExampleID: "e65-grid-copy-late",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: lateText("قبل الشبكة.", "Before.", useEnglish)},
					{Line: 1, Caption: lateText("اتعمل مكان لصفين.", "Two rows.", useEnglish), Variables: []LessonVariable{{Name: "grid", Value: "[nil nil]"}}},
					{Line: 2, Caption: lateText("الصف صفر بقى 3 خانات.", "Row 0: three cells.", useEnglish), Variables: []LessonVariable{{Name: "grid", Value: "[[0 0 0] nil]"}}},
					{Line: 3, Caption: lateText("الصف الأول بقى 3 خانات.", "Row 1: three cells.", useEnglish), Variables: []LessonVariable{{Name: "grid", Value: "[[0 0 0] [0 0 0]]"}}},
					{Line: 4, Caption: lateText("[0][1] بقت 4.", "Cell [0][1] becomes 4.", useEnglish), Variables: []LessonVariable{{Name: "grid", Value: "[[0 4 0] [0 0 0]]"}}},
					{Line: 5, Caption: lateText("[1][2] بقت 5.", "Cell [1][2] becomes 5.", useEnglish), Variables: []LessonVariable{{Name: "grid", Value: "[[0 4 0] [0 0 5]]"}}},
					{Line: 6, Caption: lateText("نسخنا outer بس.", "Outer slice cloned.", useEnglish), Variables: []LessonVariable{{Name: "outer", Value: "rows still shared"}}},
					{Line: 7, Caption: lateText("الشبكتين شافوا 9.", "Both show 9.", useEnglish), Variables: []LessonVariable{{Name: "grid[0][1]", Value: "9"}, {Name: "outer[0][1]", Value: "9"}}},
					{Line: 8, Caption: lateText("نسخنا كل صف لوحده.", "Clone each row.", useEnglish), Variables: []LessonVariable{{Name: "rows", Value: "independent row slices"}}},
					{Line: 9, Caption: lateText("النسخة بقت 7 لوحدها.", "Copy row now has 7.", useEnglish), Variables: []LessonVariable{{Name: "grid[0][1]", Value: "9"}, {Name: "rows[0][1]", Value: "7"}}},
					{Line: 10, Caption: lateText("عدد الصفوف والأعمدة، ثم خلية.", "Rows, columns, cell.", useEnglish), Variables: []LessonVariable{{Name: "len(grid)", Value: "2"}, {Name: "len(grid[0])", Value: "3"}, {Name: "grid[1][2]", Value: "5"}}, Output: "2 3 5"},
					{Line: 11, Caption: lateText("الخرج يوضح الصف المشترك والمنسوخ.", "Shared 9; copy 7.", useEnglish), Variables: []LessonVariable{{Name: "grid", Value: "[[0 9 0] [0 0 5]]"}, {Name: "outer", Value: "[[0 9 0] [0 0 5]]"}, {Name: "rows", Value: "[[0 7 0] [0 0 5]]"}}, Output: "2 3 5\n9 9 7"},
				}},
			}},
		},
		{
			ID: "c65.late.try-step", Title: lateText("جرّب", "Try it", useEnglish),
			Blocks: []LessonBlock{{
				ID: "c65.late.grid-sandbox", Kind: "sandbox",
				Text: lateText("اختار خلية وشوف ~grid[r][c]~؛ بدّل النسخ وعدّل النسخة.", "Select a cell for ~grid[r][c]~; change clone mode and edit the copy.", useEnglish),
				Sandbox: &LessonSandbox{
					Kind: "grid-index", Grid: [][]int{{3, 7, 2}, {8, 6, 4}},
					SourceLabel:         lateText("شبكة توضيحية، مش بيانات المهمة.", "Illustrative grid, not mission data.", useEnglish),
					GridOuterCloneLabel: lateText("نسخة outer بس", "Clone outer only", useEnglish),
					GridRowCloneLabel:   lateText("نسخ كل row", "Clone each row", useEnglish),
					GridEditCopyLabel:   lateText("عدّل النسخة", "Edit copy", useEnglish),
					ResetLabel:          lateText("من الأول", "Reset", useEnglish),
				},
			}},
		},
		{
			ID: "c65.late.watch-step", Title: lateText("خد بالك", "Watch out", useEnglish),
			Blocks: []LessonBlock{{ID: "c65.late.watch", Kind: "markdown", Text: lateText(
				"- ~grid[r][c]~: row ثم column.\n- ~slices.Clone(grid)~ سطحي؛ الصفوف مشتركة.",
				"- ~grid[r][c]~: row, then column.\n- ~slices.Clone(grid)~ is shallow; rows stay shared.", useEnglish)}},
		},
		{
			ID: "c65.late.check-step", Title: lateText("فكّر واختبر", "Think & check", useEnglish),
			Blocks: []LessonBlock{
				{ID: "c65.late.plan", Kind: "markdown", Text: lateText("الشكل → العملية → نتيجة مستقلة أو error.", "Shape → operation → independent result or error.", useEnglish)},
				{ID: "c65.late.contract", Kind: "markdown", Text: lateText(
					"OverlayHeat بتجمع خانات لشبكتين نفس الشكل.\n\nكل خانة في ThroughHeat مجموع ضرب صف في عمود؛ أعمدة A = صفوف B.\n\nFlipHeat بتبدّل المحاور.\n\nالمدخل الفاضي أو المتعرج أو أبعاده غلط يرجّع error؛ وممنوع تعدّل المدخلات.\n\nصفوف الناتج مستقلة؛ حتى 1×1 والصف أو العمود الواحد صالحين.",
					"OverlayHeat adds matching cells.\n\nEach ThroughHeat cell sums row-by-column products; columns(A)=rows(B).\n\nFlipHeat swaps axes.\n\nEmpty, ragged, or incompatible input returns an error; leave inputs unchanged.\n\nOutput rows are independent; 1×1 and row or column vectors are valid.", useEnglish)},
				lateQuestion("c65.late.q-index", "في شبكة جرّب، ~grid[1][0]~ كام؟", "In the playground grid, what is ~grid[1][0]~?",
					[]string{"7", "8", "6"}, []string{"7", "8", "6"}, 1,
					[]string{"7 في [0][1].", "صح؛ row 1 وcolumn 0.", "6 في [1][2]."}, []string{"7 is at [0][1].", "Right; row 1, column 0.", "6 is at [1][2]."}, useEnglish),
				lateQuestion("c65.late.q-shape", "A 2×3 وB 3×4؛ ThroughHeat صالحة؟", "A is 2×3 and B is 3×4; is ThroughHeat valid?",
					[]string{"آه؛ 2×4", "لأ", "آه؛ 3×3"}, []string{"Yes; 2×4", "No", "Yes; 3×3"}, 0,
					[]string{"صح؛ 3 هو البعد المشترك.", "ThroughHeat محتاجة أعمدة A تساوي صفوف B.", "3 هو البعد المشترك، مش الناتج."}, []string{"Right; 3 is the shared dimension.", "ThroughHeat needs A's columns to equal B's rows.", "3 is shared, not the output shape."}, useEnglish),
				lateQuestion("c65.late.q-clone", "نسخت outer وعدّلت خلية؛ إيه اللي يحصل؟", "Clone the outer slice and edit a cell; what happens?",
					[]string{"المصدر يتغير", "النسخة بس", "Compile error"}, []string{"Source changes", "Copy only", "Compile error"}, 0,
					[]string{"الصفوف لسه مشتركة.", "ده محتاج نسخ كل row.", "Go تسمح بالنسخ السطحي."}, []string{"Rows are still shared.", "That needs each row cloned.", "Go allows shallow clones."}, useEnglish),
			},
		},
		{
			ID: "c65.late.turn-step", Title: lateText("دورك", "Your turn", useEnglish),
			Blocks: []LessonBlock{{ID: "c65.late.turn", Kind: "markdown", Text: lateText("شغّل الـstarter، وحدد أول contract فشل، واختار دالة.", "Run; inspect the first failure; choose a function.", useEnglish)}},
		},
	}
}

var (
	stockOrderLessonArabic  = stockOrderLesson(false)
	stockOrderLessonEnglish = stockOrderLesson(true)
	heatGridLessonArabic    = heatGridLesson(false)
	heatGridLessonEnglish   = heatGridLesson(true)
)
