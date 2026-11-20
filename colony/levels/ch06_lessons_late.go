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

// stockOrderExample is the mission in miniature: four readings the drone
// took at x = 0..3, ranked the way the report wants them.
const stockOrderExample = `sizes := []int{3, 5, 3, 5}
positions := []int{0, 1, 2, 3}
slices.SortFunc(positions, func(a, b int) int {
	if sizes[a] != sizes[b] {
		return cmp.Compare(sizes[b], sizes[a])
	}
	return cmp.Compare(a, b)
})
fmt.Println(positions)`

func stockOrderLesson(useEnglish bool) []LessonStep {
	return []LessonStep{
		{
			ID: "c64.late.mission-step", Title: lateText("المهمة", "Mission", useEnglish),
			Blocks: []LessonBlock{{ID: "c64.late.mission", Kind: "markdown", Text: lateText(
				"1. ~Measure~ بتقيس الخانة اللي تحت الدرون بس: امشي شرق وقيس كل خانة.\n2. رتّب المواقع: الحجم الأكبر الأول، ولو اتعادلوا الـx الأصغر.\n3. ابعت المواقع بـ~Report~ مرة واحدة. ممنوع الحصاد.\n\nالذهب: **19 tick**.",
				"1. ~Measure~ reads only the tile under the drone: walk east and measure every tile.\n2. Rank positions: bigger size first, smaller x on ties.\n3. ~Report~ the positions once. No harvesting.\n\nGold: **19 ticks**.", useEnglish)}},
		},
		{
			ID: "c64.late.idea-step", Title: lateText("الفكرة", "Idea", useEnglish),
			Blocks: []LessonBlock{{
				ID: "c64.late.trace", Kind: "walkthrough",
				Text: lateText("زي طابور: ~SortFunc~ بتسأل الـcomparator عن كل اتنين «مين يقف الأول؟». السالب يقدّم a، والموجب يقدّم b، والصفر تعادل.\n\nده نفس شكل المهمة بأربع خانات.", "Like a queue: ~SortFunc~ asks the comparator \"who stands first?\" Negative = a, positive = b, zero = tie.\n\nThe mission in four tiles.", useEnglish),
				Code: stockOrderExample, Imports: []string{"cmp", "fmt", "slices"}, ExampleID: "e64-rank-crystals",
				Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
					{Line: 0, Caption: lateText("الدرون لسه مقاسش.", "Not measured yet.", useEnglish)},
					{Line: 1, Caption: lateText("القراءات: الحجم عند كل x.", "Size at each x.", useEnglish), Variables: []LessonVariable{{Name: "sizes", Value: "[3 5 3 5]"}}},
					{Line: 2, Caption: lateText("دي المواقع اللي هنرتّبها.", "Positions to rank.", useEnglish), Variables: []LessonVariable{{Name: "positions", Value: "[0 1 2 3]"}}},
					{Line: 5, Caption: lateText("1 قبل 0: الحجم 5 أكبر.", "1 before 0: 5 wins.", useEnglish), Variables: []LessonVariable{{Name: "a, b", Value: "1, 0"}, {Name: "sizes", Value: "5 vs 3"}, {Name: "result", Value: "-1"}}},
					{Line: 7, Caption: lateText("1 و3 متعادلين: الـx الأصغر الأول.", "1 and 3 tie: smaller x first.", useEnglish), Variables: []LessonVariable{{Name: "a, b", Value: "1, 3"}, {Name: "sizes", Value: "5 vs 5"}, {Name: "result", Value: "-1"}}},
					{Line: 8, Caption: lateText("الترتيب خلص.", "Ranking done.", useEnglish), Variables: []LessonVariable{{Name: "positions", Value: "[1 3 0 2]"}}},
					{Line: 9, Caption: lateText("ده اللي هتبعته بـReport.", "What you Report.", useEnglish), Variables: []LessonVariable{{Name: "positions", Value: "[1 3 0 2]"}}, Output: "[1 3 0 2]"},
				}},
			}},
		},
		{
			ID: "c64.late.try-step", Title: lateText("جرّب", "Try it", useEnglish),
			Blocks: []LessonBlock{{
				ID: "c64.late.sort-sandbox", Kind: "sandbox",
				Text: lateText("اقلب اتجاه الحجم وكسر التعادل لحد ما توصل للترتيب اللي المهمة عايزاه.", "Flip size direction and tie-break until you reach the order the mission wants.", useEnglish),
				Sandbox: &LessonSandbox{
					Kind: "sort-compare", Sequence: []int{3, 5, 3, 5, 2},
					SourceLabel:        lateText("أحجام زي اللي هتقيسها.", "Sizes like the ones you will measure.", useEnglish),
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
				"- سطر ~fmt.Println(colony.Measure())~ في الـstarter بيضيّع tick: امسحه.\n- كل ~Measure~ وكل ~Move~ بـtick: لو العرض 10، يبقى 10 قياسات و9 حركات، يعني 19 tick.\n- متتحركش بعد آخر خانة.\n- رتّب ~positions~ مش ~sizes~، وإلا المواقع تضيع.",
				"- Delete the starter's ~fmt.Println(colony.Measure())~: it wastes a tick.\n- Each ~Measure~ and ~Move~ costs a tick: width 10 = 10 + 9 = 19.\n- Do not move past the last tile.\n- Sort ~positions~, not ~sizes~.", useEnglish)}},
		},
		{
			ID: "c64.late.check-step", Title: lateText("فكّر واختبر", "Think & check", useEnglish),
			Blocks: []LessonBlock{
				{ID: "c64.late.plan", Kind: "markdown", Text: lateText("امشي وقيس → خزّن في ~sizes~ → رتّب ~positions~ → ~Report~ مرة واحدة.", "Walk and measure → keep ~sizes~ → rank ~positions~ → one ~Report~.", useEnglish)},
				lateQuestion("c64.late.q-bigger", "الموقع 2 حجمه 4، والموقع 0 حجمه 1. الـcomparator يرجّع إيه لـ(2, 0)؟", "Size 4 at x=2, size 1 at x=0. Comparator result for (2, 0)?",
					[]string{"سالب", "صفر", "موجب"}, []string{"Negative", "Zero", "Positive"}, 0,
					[]string{"صح؛ الأكبر يقف الأول.", "الصفر للتعادل بس.", "الموجب كان هيأخّر الأكبر."}, []string{"Right; bigger first.", "Zero means a tie.", "Positive puts the bigger one later."}, useEnglish),
				lateQuestion("c64.late.q-tie", "الموقع 1 والموقع 4 حجمهم 5. مين الأول؟", "x=1 and x=4 both have size 5. Which first?",
					[]string{"1", "4", "مش فارقة"}, []string{"1", "4", "Either"}, 0,
					[]string{"صح؛ عند التعادل الـx الأصغر.", "عند التعادل الأكبر بيتأخر.", "المهمة بتفرق: الأصغر الأول."}, []string{"Right; smaller x wins ties.", "Larger x goes later on ties.", "The mission says smaller first."}, useEnglish),
				lateQuestion("c64.late.q-report", "تبعت إيه في ~Report~؟", "What goes in the ~Report~?",
					[]string{"الأحجام", "المواقع", "الاتنين"}, []string{"Sizes", "Positions", "Both"}, 1,
					[]string{"المهمة عايزة المواقع.", "صح؛ المواقع بعد الترتيب.", "المواقع بس، مرة واحدة."}, []string{"The mission wants positions.", "Right; the ranked positions.", "Positions only, once."}, useEnglish),
			},
		},
		{
			ID: "c64.late.turn-step", Title: lateText("دورك", "Your turn", useEnglish),
			Blocks: []LessonBlock{{ID: "c64.late.turn", Kind: "markdown", Text: lateText("امسح سطر الـPrintln، اكتب الـloop اللي بتقيس، وبعدها الترتيب. Run!", "Delete the Println, write the measuring loop, then the ranking. Run!", useEnglish)}},
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
				"دي lab: زر Run بيشغّل tests للـfunctions، مش حركة للدرون. كمّل ~OverlayHeat~ و~ThroughHeat~ و~FlipHeat~؛ مفيش ~Report~ ولا ticks.",
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
				lateQuestion("c65.late.q-clone", "نسخت الـouter slice وعدّلت خلية؛ إيه اللي يحصل؟", "Clone the outer slice and edit a cell; what happens?",
					[]string{"المصدر يتغير", "النسخة بس", "Compile error"}, []string{"Source changes", "Copy only", "Compile error"}, 0,
					[]string{"الصفوف لسه مشتركة.", "ده محتاج نسخ كل row.", "Go بتسمح بالنسخ السطحي."}, []string{"Rows are still shared.", "That needs each row cloned.", "Go allows shallow clones."}, useEnglish),
			},
		},
		{
			ID: "c65.late.turn-step", Title: lateText("دورك", "Your turn", useEnglish),
			Blocks: []LessonBlock{{ID: "c65.late.turn", Kind: "markdown", Text: lateText("شغّل الـstarter، وحدد أول contract فشل، واختار function.", "Run; inspect the first failure; choose a function.", useEnglish)}},
		},
	}
}

var (
	stockOrderLessonArabic  = stockOrderLesson(false)
	stockOrderLessonEnglish = stockOrderLesson(true)
	heatGridLessonArabic    = heatGridLesson(false)
	heatGridLessonEnglish   = heatGridLesson(true)
)
