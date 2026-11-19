package levels

var ch09BossArabic = []LessonStep{
	{
		ID: "c9b.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{{
			ID: "c9b.mission", Kind: "markdown",
			Text: "1. كمّل ~Pack~ في ~main.go~ عشان ينظم بيانات الشحنات حسب البوابة.\n2. اقرأ تعريف ~Crate~ والتصريح ~Clearance~ قبل ما تختار خطوات التنفيذ.\n3. ده lab: اشتغل في الدالة وسيب ~main~ فاضية؛ عقدها الكامل موجود في وصف المستوى.",
		}},
	},
	{
		ID: "c9b.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{{
			ID: "c9b.idea", Kind: "markdown",
			Text: "~Crate~ بتجمع حقول مرتبطة في struct واحدة. ~Clearance~ بتوصف method محتاجها الكود؛ أي type عنده التوقيع يحقق الـinterface تلقائيًا.\n\nتوقيع ~Pack~ بيحدد المدخلات والناتجين. اتعامل مع الـslice والنتيجة كحالتين منفصلتين.\n\nاقرأ العقد قبل الكود: شكل البيانات، دور التصريح، ومسارات النجاح والخطأ. الأمثلة هنا للتدريب بس.",
		}},
	},
	{
		ID: "c9b.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			{ID: "c9b.map-predict", Kind: "question", Text: "لو ~counts~ بدأت فاضية، تتوقع المثال يطبع إيه؟\n\n~~~go\ncounts := map[string]int{}\ncounts[\"teal\"] += 2\ncounts[\"teal\"]++\nfmt.Println(counts[\"teal\"])\n~~~", Question: &LessonQuestion{
				Options: []string{"2", "3", "0"}, Answer: 1,
				Feedback: []string{"التحديث التاني لسه ما اتضافش.", "التحديثين غيروا قيمة المفتاح نفسه.", "الخريطة بدأت فاضية، بس الكود غيّر قيمتها."},
			}},
			{ID: "c9b.map-reveal", Kind: "walkthrough", Text: "كل تحديث بيغيّر قيمة المفتاح نفسه.", Code: "counts := map[string]int{}\ncounts[\"teal\"] += 2\ncounts[\"teal\"]++\nfmt.Println(counts[\"teal\"])", Imports: []string{"fmt"}, ExampleID: "e9boss-map-count", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "قبل إنشاء الخريطة."},
				{Line: 1, Caption: "بدأنا بخريطة فاضية.", Variables: []LessonVariable{{Name: "counts", Value: "map[]"}}},
				{Line: 2, Caption: "القيمة بقت 2.", Variables: []LessonVariable{{Name: "counts", Value: "map[teal:2]"}}},
				{Line: 3, Caption: "التحديث زوّد نفس المفتاح.", Variables: []LessonVariable{{Name: "counts", Value: "map[teal:3]"}}},
				{Line: 4, Caption: "اطبع قيمة المفتاح.", Variables: []LessonVariable{{Name: "counts", Value: "map[teal:3]"}}, Output: "3"},
			}}},
		},
	},
	{
		ID: "c9b.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{{
			ID: "c9b.watch", Kind: "markdown",
			Text: "~nil~ map وmap فاضية مش نفس الحالة؛ العقد بيحدد النتيجة المطلوبة.\n\nلفة ~range~ على map مفيش لها ترتيب مضمون؛ ما تعتمدش على ترتيب المفاتيح.\n\nالـ~error~ ناتج مستقل، فراجع قيمته قبل ما تستخدم البيانات.",
		}},
	},
	{
		ID: "c9b.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			{ID: "c9b.plan", Kind: "markdown", Text: "رتّب مراجعتك: افهم المدخل، وحدد معنى كل نتيجة، وبعدها راجع مسارات الخطأ.\n\nاسأل نفسك: هل Go بتطلب إعلان إن النوع يحقق interface؟"},
			{ID: "c9b.q-interface", Kind: "question", Text: "إزاي type يحقق interface في Go؟", Question: &LessonQuestion{
				Options: []string{"بكتابة كلمة implements", "بتوفير الـmethods المطلوبة"}, Answer: 1,
				Feedback: []string{"Go ما فيهاش كلمة implements للعقود دي.", "Go بتطابق الـinterface حسب توقيعات الـmethods."},
			}},
			{ID: "c9b.q-error", Kind: "question", Text: "إيه معنى إن ~err == nil~؟", Question: &LessonQuestion{
				Options: []string{"الدالة ما رجعتش error", "الـmap لازم تكون nil"}, Answer: 0,
				Feedback: []string{"دي إشارة إن مفيش خطأ مُبلّغ عنه.", "الـerror والـmap ناتجين منفصلين."},
			}},
		},
	},
	{
		ID: "c9b.turn-step", Title: "دورك",
		Blocks: []LessonBlock{{
			ID: "c9b.turn", Kind: "markdown",
			Text: "طبّق وصف المستوى على ~Pack~ وسيب ~main~ فاضية. شغّل ~go test~ للـcontract؛ بعد أي failure ارجع للشرط المرتبط بيه وعدّل خطوة واحدة.",
		}},
	},
}

var ch09BossEnglish = []LessonStep{
	{
		ID: "c9b.mission-step", Title: "Mission",
		Blocks: []LessonBlock{{
			ID: "c9b.mission", Kind: "markdown",
			Text: "1. Complete ~Pack~ in ~main.go~ to organize cargo data by bay.\n2. Read the ~Crate~ and ~Clearance~ declarations before planning the function.\n3. This is a lab: work in the function and leave ~main~ empty; its full contract is in the level description.",
		}},
	},
	{
		ID: "c9b.idea-step", Title: "Idea",
		Blocks: []LessonBlock{{
			ID: "c9b.idea", Kind: "markdown",
			Text: "~Crate~ groups fields in a struct. ~Clearance~ requires a method; any type with that signature satisfies it automatically.\n\n~Pack~'s signature lists inputs and results. Treat the input slice and result as separate state.\n\nRead the contract first: data, clearance behavior, success and error paths. Examples are illustrative.",
		}},
	},
	{
		ID: "c9b.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			{ID: "c9b.map-predict", Kind: "question", Text: "If ~counts~ starts empty, what will this example print?\n\n~~~go\ncounts := map[string]int{}\ncounts[\"teal\"] += 2\ncounts[\"teal\"]++\nfmt.Println(counts[\"teal\"])\n~~~", Question: &LessonQuestion{
				Options: []string{"2", "3", "0"}, Answer: 1,
				Feedback: []string{"The second update has not been added yet.", "Both updates changed the same key's value.", "The map started empty, but the code changed it."},
			}},
			{ID: "c9b.map-reveal", Kind: "walkthrough", Text: "Each update changes the value for the same key.", Code: "counts := map[string]int{}\ncounts[\"teal\"] += 2\ncounts[\"teal\"]++\nfmt.Println(counts[\"teal\"])", Imports: []string{"fmt"}, ExampleID: "e9boss-map-count", Walkthrough: &LessonWalkthrough{Frames: []LessonFrame{
				{Line: 0, Caption: "Before creating the map."},
				{Line: 1, Caption: "The map starts empty.", Variables: []LessonVariable{{Name: "counts", Value: "map[]"}}},
				{Line: 2, Caption: "The value becomes 2.", Variables: []LessonVariable{{Name: "counts", Value: "map[teal:2]"}}},
				{Line: 3, Caption: "The same key increases.", Variables: []LessonVariable{{Name: "counts", Value: "map[teal:3]"}}},
				{Line: 4, Caption: "Print the key's value.", Variables: []LessonVariable{{Name: "counts", Value: "map[teal:3]"}}, Output: "3"},
			}}},
		},
	},
	{
		ID: "c9b.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{{
			ID: "c9b.watch", Kind: "markdown",
			Text: "A ~nil~ map differs from an empty map; follow the contract.\n\nMap ~range~ order is unspecified; do not rely on key order.\n\nCheck the separate ~error~ before using returned data.",
		}},
	},
	{
		ID: "c9b.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			{ID: "c9b.plan", Kind: "markdown", Text: "Review the input, define each result, then inspect errors.\n\nDoes Go require a declaration that a type satisfies an interface?"},
			{ID: "c9b.q-interface", Kind: "question", Text: "How does a type satisfy an interface in Go?", Question: &LessonQuestion{
				Options: []string{"By writing the word implements", "By providing the required methods"}, Answer: 1,
				Feedback: []string{"Go has no implements keyword for these contracts.", "Go matches interfaces by method signatures."},
			}},
			{ID: "c9b.q-error", Kind: "question", Text: "What does ~err == nil~ mean?", Question: &LessonQuestion{
				Options: []string{"The function reported no error", "The map must be nil"}, Answer: 0,
				Feedback: []string{"It means no error was reported.", "The error and map are separate results."},
			}},
		},
	},
	{
		ID: "c9b.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{{
			ID: "c9b.turn", Kind: "markdown",
			Text: "Implement ~Pack~ from the contract and leave ~main~ empty.\n\nRun ~go test~; after a failure, revisit its condition and change one thing.",
		}},
	},
}
