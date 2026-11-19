package levels

func ch22Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch22Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch22QuestionExample(id, prompt, code string, options []string, answer int, feedback []string) LessonBlock {
	return ch22Question(id, prompt+"\n\n~~~go\n"+code+"\n~~~", options, answer, feedback)
}

func ch22Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
	return LessonBlock{
		ID:          id,
		Kind:        "walkthrough",
		Text:        text,
		Code:        code,
		Imports:     imports,
		ExampleID:   exampleID,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch22EnglishFrames(frames []LessonFrame) []LessonFrame {
	translations := map[string]string{
		"قبل التنفيذ.":              "Before run.",
		"خريطة الملفات للتوضيح.":    "Illustrative file map.",
		"اقرأ أبناء الجذر.":         "Read the root entries.",
		"أول اسم هو alpha.":         "alpha is the first name.",
		"ابدأ العد من 2.":           "Start counting from 2.",
		"لسه ما وصلناش للصفر.":      "Zero is not reached yet.",
		"اطبع القيمة الحالية.":      "Print the current value.",
		"النداء يبدأ قيمة أصغر.":    "The call starts with less.",
		"النداء التاني يطبع 1.":     "The next call prints 1.",
		"الصفر يوقف النداءات.":      "Zero stops further calls.",
		"خطة للتوضيح.":              "Illustrative plan.",
		"ترتيب ثابت قبل التنفيذ.":   "Stable order before execution.",
		"اطبع الترتيب الجديد.":      "Print the new order.",
		"وحّد حالة الامتداد.":       "Normalize the extension case.",
		"الامتداد بقى صغير.":        "The extension is lowercase.",
		"مساحة مؤقتة للتجهيز.":      "A temporary staging space.",
		"الكتابة ما زالت معزولة.":   "The write stays isolated.",
		"الناتج لسه جاهز.":          "The result is ready.",
		"ابدأ بمساحة منفصلة.":       "Start in a separate space.",
		"اكتب للحزمة المؤقتة.":      "Write into the temporary buffer.",
		"اعرض اللي اتجهز.":          "Show the prepared value.",
		"بصمة جديدة لكل stream.":    "A fresh hash for each stream.",
		"Copy يمرر البايتات.":       "Copy passes the bytes along.",
		"عدد البايتات وطول البصمة.": "Byte count and hash length.",
		"محتوى أول للتوضيح.":        "First illustrative content.",
		"محتوى تاني بنفس الحجم.":    "Second content with equal size.",
		"الحجم وحده مش كفاية.":      "Size alone is not enough.",
	}
	translated := make([]LessonFrame, len(frames))
	for i, frame := range frames {
		translated[i] = frame
		translated[i].Caption = translations[frame.Caption]
	}
	return translated
}

const ch22Example221Idea = `func main() {
	store := fstest.MapFS{
		"zeta":  &fstest.MapFile{Data: []byte("z")},
		"alpha": &fstest.MapFile{Data: []byte("a")},
	}
	entries, err := fs.ReadDir(store, ".")
	fmt.Println(entries[0].Name(), len(entries), err)
}`

const ch22Example221Try = `func count(n int) {
	if n == 0 {
		return
	}
	fmt.Println(n)
	count(n - 1)
}
func main() {
	count(2)
}`

var ch22Frames221Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "خريطة الملفات للتوضيح."},
	{Line: 6, Caption: "اقرأ أبناء الجذر."},
	{Line: 7, Caption: "أول اسم هو alpha.", Output: "alpha 2 <nil>"},
}

var ch22Frames221Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 9, Caption: "ابدأ العد من 2."},
	{Line: 2, Caption: "لسه ما وصلناش للصفر."},
	{Line: 5, Caption: "اطبع القيمة الحالية.", Output: "2"},
	{Line: 6, Caption: "النداء يبدأ قيمة أصغر.", Output: "2"},
	{Line: 5, Caption: "النداء التاني يطبع 1.", Output: "2\n1"},
	{Line: 3, Caption: "الصفر يوقف النداءات.", Output: "2\n1"},
}

var ch22Lesson221Arabic = []LessonStep{
	{ID: "c221.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch22Markdown("c221.mission", "المخزن هنا ~fs.FS~ محقون، مش فولدر جهازك.\n\nاعرض أسماءه في شجرة مرتبة، وخلي حد العمق يوقف الفروع الأبعد.\n\nالـstarter فيه توقيع ~StoreTree~ ولسه محتاج تنفيذ."),
	}},
	{ID: "c221.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch22Walkthrough("c221.fs-trace", "تخيل كل مجلد محطة صغيرة: اقرأ أبناءه، وبعدها تابع المجلدات اللي يسمح بيها الحد.\n\n~fs.ReadDir~ بترتب الأسماء؛ و~fs.DirEntry~ فيها الاسم المحلي ونوع العنصر.", ch22Example221Idea, "e221-fs-readdir", []string{"fmt", "io/fs", "testing/fstest"}, ch22Frames221Idea),
	}},
	{ID: "c221.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch22QuestionExample("c221.recursion-predict", "توقّع الناتج قبل ما تشوف التنفيذ.", ch22Example221Try, []string{"2 بس", "2 وبعدها 1", "2 وبعدها 1 وبعدها 0"}, 1, []string{
			"القيمة 1 بتتطبع في نداء تاني.",
			"كل نداء يطبع قبل ما يستدعي قيمة أصغر.",
			"لما نوصل للصفر بنرجع من غير طباعة.",
		}),
		ch22Walkthrough("c221.recursion-reveal", "تابع القيمة وهي بتصغر لحد حالة التوقف.", ch22Example221Try, "e221-recursion", []string{"fmt"}, ch22Frames221Try),
	}},
	{ID: "c221.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch22Markdown("c221.watch", "- افحص الخطأ قبل استخدام ناتج القراءة.\n- عند حد العمق، متقراش المجلد الأبعد."),
	}},
	{ID: "c221.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch22Markdown("c221.plan", "افصل بين التحقق من الجذر وعرض محتواه.\n\nبعدها حدّد إمتى الفرع يقف، وإمتى يكمل. خطأ من filesystem أو writer لازم يفضل ظاهر."),
		ch22Question("c221.stop-question", "إيه اللي بيوقف recursion بأمان؟", []string{"حالة توقف واضحة", "اسم متغير جديد", "ترتيب مختلف للأبناء"}, 0, []string{
			"من غير حالة توقف، النداءات مبتنتهيش.",
			"اسم المتغير مش بيغيّر عدد النداءات.",
			"الترتيب مش بديل عن حالة التوقف.",
		}),
	}},
	{ID: "c221.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch22Markdown("c221.turn", "ابدأ بالجذر، وبعدها جرّب فرع صغير. راجع النص المكتوب بعد أول خطأ."),
	}},
}

var ch22Lesson221English = []LessonStep{
	{ID: "c221.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch22Markdown("c221.mission", "This store is an injected ~fs.FS~, not your machine's folder.\n\nRender its names as a sorted tree and let the depth limit stop deeper branches.\n\nThe starter has the ~StoreTree~ signature, ready for implementation."),
	}},
	{ID: "c221.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch22Walkthrough("c221.fs-trace", "Treat each directory as a small stop: read its entries, then follow directories allowed by the limit.\n\n~fs.ReadDir~ sorts names; each ~fs.DirEntry~ has a local name and entry type.", ch22Example221Idea, "e221-fs-readdir", []string{"fmt", "io/fs", "testing/fstest"}, ch22EnglishFrames(ch22Frames221Idea)),
	}},
	{ID: "c221.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch22QuestionExample("c221.recursion-predict", "Predict the output before seeing the trace.", ch22Example221Try, []string{"Only 2", "2 then 1", "2 then 1 then 0"}, 1, []string{
			"The value 1 prints in another call.",
			"Each call prints before calling a smaller value.",
			"At zero, the function returns without printing.",
		}),
		ch22Walkthrough("c221.recursion-reveal", "Follow the value as it shrinks to the stopping case.", ch22Example221Try, "e221-recursion", []string{"fmt"}, ch22EnglishFrames(ch22Frames221Try)),
	}},
	{ID: "c221.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch22Markdown("c221.watch", "- Check an error before using a read result.\n- At the depth limit, do not read a deeper directory."),
	}},
	{ID: "c221.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch22Markdown("c221.plan", "Separate root validation from listing its contents.\n\nThen decide when a branch stops and when it continues. Keep filesystem and writer errors visible."),
		ch22Question("c221.stop-question", "What safely stops recursion?", []string{"A clear stopping case", "A new variable name", "A different sibling order"}, 0, []string{
			"Without a stopping case, calls do not end.",
			"A variable name does not change the calls.",
			"Ordering is not a stopping case.",
		}),
	}},
	{ID: "c221.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch22Markdown("c221.turn", "Start at the root, then try a small branch. Check the text written before the first error."),
	}},
}

const ch22Example222Idea = `func main() {
	moves := []string{"shelf/z", "shelf/a"}
	slices.Sort(moves)
	fmt.Println(moves)
}`

const ch22Example222Try = `func main() {
	extension := strings.ToLower(path.Ext("bay/NOTE.JSON"))
	fmt.Println(extension)
}`

var ch22Frames222Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "خطة للتوضيح."},
	{Line: 3, Caption: "ترتيب ثابت قبل التنفيذ."},
	{Line: 4, Caption: "اطبع الترتيب الجديد.", Output: "[shelf/a shelf/z]"},
}

var ch22Frames222Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "وحّد حالة الامتداد."},
	{Line: 3, Caption: "الامتداد بقى صغير.", Output: ".json"},
}

var ch22Lesson222Arabic = []LessonStep{
	{ID: "c222.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch22Markdown("c222.mission", "ملفات الوارد محتاجة رف مناسب، بس اكتشافها ماينقلهاش.\n\nابني خطة تتراجع الأول، وبعدها نفّذها على مجلد حقيقي من غير مفاجآت في النص."),
	}},
	{ID: "c222.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch22Walkthrough("c222.plan-trace", "الخطة بيانات عادية؛ تقدر ترتبها وتفحصها قبل أي تغيير.\n\nخلي أسماء الواجهة slash-separated، ومسارات الجهاز مسؤولية ~filepath~.", ch22Example222Idea, "e222-sort-plan", []string{"fmt", "slices"}, ch22Frames222Idea),
	}},
	{ID: "c222.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch22QuestionExample("c222.extension-predict", "توقّع الامتداد بعد التطبيع.", ch22Example222Try, []string{".JSON", ".json", "NOTE.JSON"}, 1, []string{
			"دي الحالة الأصلية قبل التطبيع.",
			"تطبيع الامتداد يغيّر حالته بس.",
			"ده اسم الملف كله، مش الامتداد.",
		}),
		ch22Walkthrough("c222.extension-reveal", "راقب إزاي تصنيف الامتداد يفضل منفصل عن الاسم.", ch22Example222Try, "e222-extension", []string{"fmt", "path", "strings"}, ch22Frames222Try),
	}},
	{ID: "c222.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch22Markdown("c222.watch", "- اكتشاف المصدر مش معناه إنك تنقله.\n- افشل عند أي خطأ غير «غير موجود»."),
	}},
	{ID: "c222.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch22Markdown("c222.plan", "رتّب شغلك على مراحل: اكتشاف، تصنيف، تحقق، ثم تنفيذ.\n\nاسأل نفسك إيه اللي يحصل لو آخر وجهة متعارضة؟ وإيه اللي يفضل بعد فشل نقل؟"),
		ch22Question("c222.phase-question", "إمتى يبدأ تغيير الملفات؟", []string{"بعد مراجعة الخطة", "وقت اكتشاف أول ملف", "بعد تجاهل أخطاء القراءة"}, 0, []string{
			"التحقق المسبق يحمي أول الملفات من تصادم لاحق.",
			"الاكتشاف المفروض ما يعملش تغييرات.",
			"تجاهل أخطاء القراءة يخفي نقص الخطة.",
		}),
	}},
	{ID: "c222.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch22Markdown("c222.turn", "جرّب خطة على filesystem محقون، وبعدها راجع مجلد مؤقت للتنفيذ الفعلي."),
	}},
}

var ch22Lesson222English = []LessonStep{
	{ID: "c222.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch22Markdown("c222.mission", "Incoming files need the right shelf, but discovery should not move them.\n\nBuild a plan that can be reviewed first, then execute it on a real directory without mid-run surprises."),
	}},
	{ID: "c222.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch22Walkthrough("c222.plan-trace", "A plan is ordinary data; sort and inspect it before making changes.\n\nKeep interface names slash-separated, and leave host paths to ~filepath~.", ch22Example222Idea, "e222-sort-plan", []string{"fmt", "slices"}, ch22EnglishFrames(ch22Frames222Idea)),
	}},
	{ID: "c222.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch22QuestionExample("c222.extension-predict", "Predict the extension after normalization.", ch22Example222Try, []string{".JSON", ".json", "NOTE.JSON"}, 1, []string{
			"This is the original case before normalization.",
			"Normalization changes the extension's case only.",
			"This is the whole filename, not its extension.",
		}),
		ch22Walkthrough("c222.extension-reveal", "See how extension classification stays separate from the name.", ch22Example222Try, "e222-extension", []string{"fmt", "path", "strings"}, ch22EnglishFrames(ch22Frames222Try)),
	}},
	{ID: "c222.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch22Markdown("c222.watch", "- Discovering a source does not mean moving it.\n- Fail on errors other than “not found.”"),
	}},
	{ID: "c222.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch22Markdown("c222.plan", "Work in phases: discover, classify, validate, then execute.\n\nWhat if the last destination conflicts? What remains after a failed move?"),
		ch22Question("c222.phase-question", "When should file changes begin?", []string{"After reviewing the plan", "When the first file is found", "After ignoring read errors"}, 0, []string{
			"Up-front validation protects early files from later conflicts.",
			"Discovery should not make changes.",
			"Ignoring read errors hides an incomplete plan.",
		}),
	}},
	{ID: "c222.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch22Markdown("c222.turn", "Try a plan on an injected filesystem, then review a temporary directory for real execution."),
	}},
}

const ch22Example223Idea = `func main() {
	var staged strings.Builder
	_, err := io.WriteString(&staged, "ready")
	fmt.Println(staged.String(), err)
}`

const ch22Example223Try = `func main() {
	var stage strings.Builder
	_, err := io.WriteString(&stage, "packet")
	fmt.Println(stage.String(), err)
}`

var ch22Frames223Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "مساحة مؤقتة للتجهيز."},
	{Line: 3, Caption: "الكتابة ما زالت معزولة."},
	{Line: 4, Caption: "الناتج لسه جاهز.", Output: "ready <nil>"},
}

var ch22Frames223Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "ابدأ بمساحة منفصلة."},
	{Line: 3, Caption: "اكتب للحزمة المؤقتة."},
	{Line: 4, Caption: "اعرض اللي اتجهز.", Output: "packet <nil>"},
}

var ch22Lesson223Arabic = []LessonStep{
	{ID: "c223.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch22Markdown("c223.mission", "لو تجهيز الحزمة فشل بعد أول بايت، النسخة القديمة لازم تفضل سليمة.\n\nالمطلوب يفصل الكتابة عن النشر، ويخلّي فشل النقل واضح."),
	}},
	{ID: "c223.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch22Walkthrough("c223.staging-trace", "الـwriter عقدة للكتابة، مش الملف النهائي.\n\nجهّز المحتوى في مساحة مستقلة؛ النشر ييجي بعد نجاح الكتابة والقفل بالتسلسل.", ch22Example223Idea, "e223-staging-writer", []string{"fmt", "io", "strings"}, ch22Frames223Idea),
	}},
	{ID: "c223.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch22QuestionExample("c223.writer-predict", "توقّع النص والخطأ اللي هيظهروا.", ch22Example223Try, []string{"packet <nil>", "فاضي <nil>", "packet EOF"}, 0, []string{
			"الـBuilder احتفظ بالنص والكتابة نجحت.",
			"الـBuilder استقبل البايتات.",
			"الـwriter قبل النص من غير خطأ.",
		}),
		ch22Walkthrough("c223.writer-reveal", "اتبع الكتابة قبل ما تتخيل أي نشر للملف.", ch22Example223Try, "e223-writer", []string{"fmt", "io", "strings"}, ch22Frames223Try),
	}},
	{ID: "c223.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch22Markdown("c223.watch", "- متفتحش الهدف للكتابة قبل اكتمال المحتوى.\n- افحص خطأ ~Close~ قبل النقل."),
	}},
	{ID: "c223.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch22Markdown("c223.plan", "ارسم عمر الملف: إنشاء، callback، قفل، ثم قرار النقل.\n\nلو مرحلة فشلت، فكّر مين بينضّف الاسم المؤقت وإيه اللي لازم يفضل عند الهدف."),
		ch22Question("c223.publish-question", "إيه اللي ينفع يتنشر؟", []string{"ملف اكتمل وقفله نجح", "أول بايت اتكتب", "اسم مؤقت اتعمل"}, 0, []string{
			"النشر ييجي بعد اكتمال الملف وقفل ناجح.",
			"جزء من المحتوى مش نسخة كاملة.",
			"وجود الاسم مش معناه إن المحتوى جاهز.",
		}),
	}},
	{ID: "c223.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch22Markdown("c223.turn", "اختبر callback ناجح وفاشل، وبعدها راقب الهدف والأسماء المؤقتة."),
	}},
}

var ch22Lesson223English = []LessonStep{
	{ID: "c223.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch22Markdown("c223.mission", "If packet generation fails after its first bytes, the old version must stay intact.\n\nSeparate writing from publication and keep a failed move visible."),
	}},
	{ID: "c223.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch22Walkthrough("c223.staging-trace", "A ~writer~ is a writing contract, not the final file.\n\nPrepare content in a separate space; publish only after writing and closing succeed.", ch22Example223Idea, "e223-staging-writer", []string{"fmt", "io", "strings"}, ch22EnglishFrames(ch22Frames223Idea)),
	}},
	{ID: "c223.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch22QuestionExample("c223.writer-predict", "Predict the text and error that appear.", ch22Example223Try, []string{"packet <nil>", "empty <nil>", "packet EOF"}, 0, []string{
			"The builder kept the text and the write succeeded.",
			"The builder received the bytes.",
			"The writer accepted the text without an error.",
		}),
		ch22Walkthrough("c223.writer-reveal", "Follow the write before imagining any file publication.", ch22Example223Try, "e223-writer", []string{"fmt", "io", "strings"}, ch22EnglishFrames(ch22Frames223Try)),
	}},
	{ID: "c223.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch22Markdown("c223.watch", "- Do not open the target for writing before content is complete.\n- Check ~Close~ before renaming."),
	}},
	{ID: "c223.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch22Markdown("c223.plan", "Trace the file lifetime: create, callback, close, then move decision.\n\nIf a stage fails, who cleans the temporary name, and what must remain at the target?"),
		ch22Question("c223.publish-question", "What is ready to publish?", []string{"A complete file that closed successfully", "The first written byte", "A created temporary name"}, 0, []string{
			"Publish only after the complete file closes successfully.",
			"Partial content is not a complete version.",
			"An existing name does not mean content is ready.",
		}),
	}},
	{ID: "c223.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch22Markdown("c223.turn", "Test a successful and failed callback, then inspect the target and temporary names."),
	}},
}

const ch22Example224Idea = `func main() {
	digest := sha256.New()
	n, err := io.Copy(digest, strings.NewReader("echo"))
	fmt.Println(n, len(digest.Sum(nil)), err)
}`

const ch22Example224Try = `func main() {
	first := []byte("A")
	second := []byte("B")
	fmt.Println(len(first) == len(second), string(first) == string(second))
}`

var ch22Frames224Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "بصمة جديدة لكل stream."},
	{Line: 3, Caption: "Copy يمرر البايتات."},
	{Line: 4, Caption: "عدد البايتات وطول البصمة.", Output: "4 32 <nil>"},
}

var ch22Frames224Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "محتوى أول للتوضيح."},
	{Line: 3, Caption: "محتوى تاني بنفس الحجم."},
	{Line: 4, Caption: "الحجم وحده مش كفاية.", Output: "true false"},
}

var ch22Lesson224Arabic = []LessonStep{
	{ID: "c224.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch22Markdown("c224.mission", "الأرشيف فيه ملفات كبيرة ونسخ متكررة.\n\nاكتشف التكرار من غير تحميل كل المحتوى للذاكرة، وخلي التقرير كامل أو فاشل بوضوح."),
	}},
	{ID: "c224.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch22Walkthrough("c224.stream-trace", "ابدأ بمعلومة رخيصة زي الحجم، وبعدها مرّر المرشحين على hash.\n\n~io.Copy~ يوصل reader بـwriter من غير ما تجمع الملف كله.", ch22Example224Idea, "e224-stream-hash", []string{"crypto/sha256", "fmt", "io", "strings"}, ch22Frames224Idea),
	}},
	{ID: "c224.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch22QuestionExample("c224.size-predict", "توقّع قيم المقارنة.", ch22Example224Try, []string{"false false", "true false", "true true"}, 1, []string{
			"الحرفين طولهم واحد.",
			"الحجم متساوي لكن المحتوى مختلف.",
			"المحتوى مش متساوي.",
		}),
		ch22Walkthrough("c224.size-reveal", "قارن الحجم بالمحتوى في مثال صغير.", ch22Example224Try, "e224-size-filter", []string{"fmt"}, ch22Frames224Try),
	}},
	{ID: "c224.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch22Markdown("c224.watch", "- نفس الحجم مش دليل إن الملفين متطابقين.\n- اقفل كل ملف وافحص خطأ القراءة والقفل."),
	}},
	{ID: "c224.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch22Markdown("c224.plan", "رتّب التفكير: اجمع metadata، اختار الأحجام المشتركة، احسب بصمة مستقلة لكل ملف، ثم جهّز التقرير.\n\nخطأ متأخر ماينفعش يسيب نتيجة ناقصة."),
		ch22Question("c224.hash-question", "ليه hash جديدة لكل ملف؟", []string{"عشان كل بصمة تمثل محتوى واحد", "عشان الحجم يبقى أكبر", "عشان نقدر نكتب في المصدر"}, 0, []string{
			"إعادة استخدام الحالة تخلط محتوى الملفات.",
			"الحساب ما بيغيرش حجم المصدر.",
			"الفحص هنا للقراءة فقط.",
		}),
	}},
	{ID: "c224.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch22Markdown("c224.turn", "جرّب مصدر فاضي ومصدر فيه ملفين مرشحين، وراجع ترتيب التقرير."),
	}},
}

var ch22Lesson224English = []LessonStep{
	{ID: "c224.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch22Markdown("c224.mission", "The archive contains large files and repeated copies.\n\nFind duplicates without loading all content into memory, and make the report complete or clearly failed."),
	}},
	{ID: "c224.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch22Walkthrough("c224.stream-trace", "Start with cheap metadata such as size, then hash candidates.\n\n~io.Copy~ connects a reader to a writer without collecting the whole file.", ch22Example224Idea, "e224-stream-hash", []string{"crypto/sha256", "fmt", "io", "strings"}, ch22EnglishFrames(ch22Frames224Idea)),
	}},
	{ID: "c224.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch22QuestionExample("c224.size-predict", "Predict the comparison values.", ch22Example224Try, []string{"false false", "true false", "true true"}, 1, []string{
			"The two letters have equal length.",
			"Sizes match, but the contents differ.",
			"The contents are not equal.",
		}),
		ch22Walkthrough("c224.size-reveal", "Compare size with content in a small example.", ch22Example224Try, "e224-size-filter", []string{"fmt"}, ch22EnglishFrames(ch22Frames224Try)),
	}},
	{ID: "c224.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch22Markdown("c224.watch", "- Equal size does not prove two files match.\n- Close every file and check read and close errors."),
	}},
	{ID: "c224.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch22Markdown("c224.plan", "Order the reasoning: collect metadata, select shared sizes, hash each file independently, then prepare the report.\n\nA late error must not leave a partial result."),
		ch22Question("c224.hash-question", "Why create a hash per file?", []string{"Each digest should represent one file", "To increase the file size", "To write to the source"}, 0, []string{
			"Reusing state mixes file contents.",
			"Hashing does not change source size.",
			"This scan is read-only.",
		}),
	}},
	{ID: "c224.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch22Markdown("c224.turn", "Try an empty source and one with two candidates, then inspect report ordering."),
	}},
}
