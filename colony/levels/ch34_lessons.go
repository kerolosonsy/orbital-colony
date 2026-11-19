package levels

import "strings"

type ch34FrameSpec struct {
	line    int
	caption [2]string
	vars    []LessonVariable
	output  string
}

type ch34ExampleSpec struct {
	id      string
	code    string
	imports []string
	frames  []ch34FrameSpec
}

type ch34SandboxFrame struct {
	caption                [2]string
	leftLabel, leftValue   [2]string
	rightLabel, rightValue [2]string
}

type ch34SandboxSpec struct {
	id             string
	text, title    [2]string
	advance, reset [2]string
	frames         []ch34SandboxFrame
}

type ch34LessonSpec struct {
	prefix        string
	mission       [2]string
	ideaText      [2]string
	idea          ch34ExampleSpec
	tryText       [2]string
	question      [2]string
	options       [2][]string
	answer        int
	feedback      [2][]string
	revealText    [2]string
	reveal        ch34ExampleSpec
	watch         [2]string
	plan          [2]string
	checkQuestion [2]string
	checkOptions  [2][]string
	checkAnswer   int
	checkFeedback [2][]string
	turn          [2]string
	hints         [2][3]string
	sandbox       ch34SandboxSpec
}

func ch34Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: ch34ShortParagraphs(text)}
}

func ch34Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: ch34ShortParagraphs(text), Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch34Walkthrough(id, text string, example ch34ExampleSpec, language int) LessonBlock {
	frames := make([]LessonFrame, 0, len(example.frames))
	for _, frame := range example.frames {
		frames = append(frames, LessonFrame{Line: frame.line, Caption: frame.caption[language], Variables: frame.vars, Output: frame.output})
	}
	return LessonBlock{
		ID: id, Kind: "walkthrough", Text: ch34ShortParagraphs(text), Code: example.code,
		Imports: example.imports, ExampleID: example.id,
		Walkthrough: &LessonWalkthrough{Frames: frames},
	}
}

func ch34Locale(value [2]string) map[string]any {
	return map[string]any{"ar": value[0], "en": value[1]}
}

func ch34SandboxBlock(spec ch34SandboxSpec, language int) LessonBlock {
	frames := make([]any, 0, len(spec.frames))
	for _, frame := range spec.frames {
		frames = append(frames, map[string]any{
			"caption":    ch34Locale(frame.caption),
			"leftLabel":  ch34Locale(frame.leftLabel),
			"leftValue":  ch34Locale(frame.leftValue),
			"rightLabel": ch34Locale(frame.rightLabel),
			"rightValue": ch34Locale(frame.rightValue),
		})
	}
	return LessonBlock{
		ID: spec.id, Kind: "sandbox", Text: spec.text[language],
		Sandbox: &LessonSandbox{Kind: "ch34-storage-flow", Config: map[string]any{
			"title":  ch34Locale(spec.title),
			"labels": map[string]any{"advance": ch34Locale(spec.advance), "reset": ch34Locale(spec.reset)},
			"frames": frames,
		}},
	}
}

func ch34FencedQuestion(text, code string) string {
	return ch34ShortParagraphs(text) + "\n\n~~~go\n" + code + "\n~~~"
}

func ch34ShortParagraphs(text string) string {
	lines := strings.Split(text, "\n")
	paragraphs := make([]string, 0, len(lines))
	inFence := false
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "~~~") {
			inFence = !inFence
			paragraphs = append(paragraphs, line)
			continue
		}
		if !inFence {
			line = strings.ReplaceAll(line, ". ", ".\n\n")
			line = strings.ReplaceAll(line, "? ", "?\n\n")
			line = strings.ReplaceAll(line, "؟ ", "؟\n\n")
			line = strings.ReplaceAll(line, "! ", "!\n\n")
		}
		paragraphs = append(paragraphs, line)
	}
	return strings.Join(paragraphs, "\n")
}

func ch34BuildLesson(spec ch34LessonSpec) ([]LessonStep, []LessonStep) {
	titles := [2][]string{
		{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"},
		{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"},
	}
	steps := [2][]LessonStep{}
	for language := range steps {
		prefix := spec.prefix
		steps[language] = []LessonStep{
			{ID: prefix + "mission-step", Title: titles[language][0], Blocks: []LessonBlock{ch34Markdown(prefix+"mission", spec.mission[language])}},
			{ID: prefix + "idea-step", Title: titles[language][1], Blocks: []LessonBlock{ch34Walkthrough(prefix+"idea-trace", spec.ideaText[language], spec.idea, language)}},
			{ID: prefix + "try-step", Title: titles[language][2], Blocks: []LessonBlock{
				ch34Question(prefix+"predict", ch34FencedQuestion(spec.question[language], spec.reveal.code), spec.options[language], spec.answer, spec.feedback[language]),
				ch34Walkthrough(prefix+"reveal", spec.revealText[language], spec.reveal, language),
				ch34SandboxBlock(spec.sandbox, language),
			}},
			{ID: prefix + "watch-step", Title: titles[language][3], Blocks: []LessonBlock{ch34Markdown(prefix+"watch", spec.watch[language])}},
			{ID: prefix + "check-step", Title: titles[language][4], Blocks: []LessonBlock{
				ch34Markdown(prefix+"plan", spec.plan[language]),
				ch34Question(prefix+"check-question", spec.checkQuestion[language], spec.checkOptions[language], spec.checkAnswer, spec.checkFeedback[language]),
			}},
			{ID: prefix + "turn-step", Title: titles[language][5], Blocks: []LessonBlock{ch34Markdown(prefix+"turn", spec.turn[language])}},
		}
	}
	return steps[0], steps[1]
}

var ch34VaultArabic, ch34VaultEnglish = ch34BuildLesson(ch34LessonSpec{
	prefix: "c341.",
	mission: [2]string{
		"ابني مخزن مفاتيح يستخدمه أكتر من goroutine بأمان. خزّن القيم مع صلاحيتها، وخلي القراءة واللقطات تعكس الحالة الحية. ركّز على ملكية البيانات وثبات اللقطة.",
		"Build a key vault that multiple goroutines can use safely. Store values with expiry, and keep reads and snapshots faithful to live state. Focus on data ownership and repeatable snapshots.",
	},
	ideaText: [2]string{
		"الـslice بتحتفظ بمرجع لبايتات مشتركة. لو غيرت المصدر، أي slice بتشاركها ممكن تتأثر. تابع المثال وشوف القيمة وهي بتتغير.",
		"A slice keeps a reference to shared bytes. Change the source, and another slice may change too. Follow the example to see the value move.",
	},
	idea: ch34ExampleSpec{
		id: "e341-byte-alias", imports: []string{"fmt"},
		code: `func main() {
	value := []byte("air")
	alias := value
	value[0] = 'f'
	fmt.Println(string(alias), string(value))
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"أنشأنا بايتات المصدر.", "The source bytes exist."}, vars: []LessonVariable{{Name: "value", Value: "[97 105 114]"}}},
			{line: 3, caption: [2]string{"الاسمين بيشاوروا لنفس البيانات.", "Both names share the bytes."}, vars: []LessonVariable{{Name: "value", Value: "[97 105 114]"}, {Name: "alias", Value: "[97 105 114]"}}},
			{line: 4, caption: [2]string{"أول بايت اتغير.", "The first byte changed."}, vars: []LessonVariable{{Name: "value", Value: "[102 105 114]"}, {Name: "alias", Value: "[102 105 114]"}}},
			{line: 5, caption: [2]string{"الاسمين بيعرضوا نفس التغيير.", "Both names show the change."}, vars: []LessonVariable{{Name: "alias", Value: "[102 105 114]"}}, output: "fir fir"},
		},
	},
	tryText: [2]string{
		"توقّع الناتج عند لحظة انتهاء الصلاحية بالظبط. الساعة هنا مثال ثابت، مش ساعة الجهاز.",
		"Predict the output at the exact expiry instant. This clock is a fixed example, not the machine clock.",
	},
	question: [2]string{"إيه الناتج؟", "What is printed?"},
	options:  [2][]string{{"true", "false", "وقت جديد"}, {"true", "false", "A new time"}}, answer: 1,
	feedback: [2][]string{
		{"الوقت مش قبل deadline؛ دي مقارنة عكسية.", "Before بيرجع false عند اللحظة نفسها.", "الكود بيقارن قيمتين موجودتين بالفعل."},
		{"That reverses the comparison; now is not before the deadline.", "Before returns false at the deadline itself.", "The code compares two existing time values."},
	},
	revealText: [2]string{"عدّي على تكوين deadline، وبعدها راقب المقارنة عند نفس اللحظة.", "Follow the deadline, then watch the comparison at that exact instant."},
	reveal: ch34ExampleSpec{
		id: "e341-expiry-boundary", imports: []string{"fmt", "time"},
		code: `func main() {
	deadline := time.Unix(100, 0)
	now := deadline
	fmt.Println(now.Before(deadline))
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"حددنا موعد النهاية.", "The deadline is set."}, vars: []LessonVariable{{Name: "deadline", Value: "1970-01-01 00:01:40 +0000 UTC"}}},
			{line: 3, caption: [2]string{"الساعة وصلت لنفس الموعد.", "The clock reaches that instant."}, vars: []LessonVariable{{Name: "now", Value: "1970-01-01 00:01:40 +0000 UTC"}}},
			{line: 4, caption: [2]string{"عند الموعد، المقارنة false.", "At the deadline, the result is false."}, vars: []LessonVariable{{Name: "now", Value: "1970-01-01 00:01:40 +0000 UTC"}}, output: "false"},
		},
	},
	watch: [2]string{
		"الـmap مش آمنة لو goroutines بتقرأ وتكتب من غير تنسيق. والـvalue اللي رجعت من المخزن لازم متكشفش البايتات الداخلية.",
		"A map is unsafe when goroutines read and write without coordination. A returned value must not expose the vault's internal bytes.",
	},
	plan: [2]string{
		"قسّم التفكير لحماية الحالة، وفحص الصلاحية، وملكية النسخ. بعد كده راجع شكل snapshot وإزاي تفضل نفس الحالة قابلة للاسترجاع.",
		"Separate state protection, expiry checks, and copy ownership. Then consider the snapshot shape and how it restores the same state.",
	},
	checkQuestion: [2]string{"ليه الساعة المحقونة مفيدة هنا؟", "Why is an injected clock useful here?"},
	checkOptions: [2][]string{
		{"تخلّي الصلاحية قابلة لإعادة التحقق.", "تمنع كل عمليات الكتابة.", "ترتب مفاتيح الـmap تلقائيًا."},
		{"It makes expiry behavior repeatable.", "It prevents every write.", "It sorts map keys automatically."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"نفس الوقت المدخل يدي نفس قرار الصلاحية.", "الساعة مش قفل للكتابة.", "ترتيب الـmap محتاج معالجة منفصلة."},
		{"The same supplied time gives the same expiry decision.", "A clock does not lock writes.", "Map ordering needs separate handling."},
	},
	turn: [2]string{
		"اقرأ توقيعات ~Vault~ في الـstarter. نفّذ العقد جزء جزء، وشغّل محاولتك بعد ما تراجع الصلاحية والنسخ والتزامن.",
		"Read the ~Vault~ signatures in the starter. Work through the contract in pieces, then run your attempt after reviewing expiry, copies, and concurrency.",
	},
	hints: [2][3]string{
		{"الـslice ممكن تشارك backing array؛ وضّح مين يملك البايتات قبل تعديلها.", "قسّم التصميم بين حماية الحالة، وحساب الصلاحية، وحدود نسخ القيم.", "ابدأ من تعريف ~Vault~ وتوقيعات methods في الـstarter، وحدد إمتى القيمة تعتبر حية."},
		{"A slice may share a backing array; clarify byte ownership before changing it.", "Separate state protection, expiry decisions, and value-copy boundaries.", "Start at the ~Vault~ type and starter methods, and define when a value is live."},
	},
	sandbox: ch34SandboxSpec{
		id:      "ch34-key-vault",
		text:    [2]string{"غيّر بايت وشوف الاسم التاني بيتأثر.", "Change one byte and watch its alias."},
		title:   [2]string{"ملكية البايتات", "Byte ownership"},
		advance: [2]string{"غيّر أول بايت", "Change first byte"}, reset: [2]string{"من الأول", "Reset"},
		frames: []ch34SandboxFrame{
			{caption: [2]string{"الاسمين شايفين نفس البايتات.", "Both names share these bytes."}, leftLabel: [2]string{"المصدر", "Source"}, leftValue: [2]string{"air", "air"}, rightLabel: [2]string{"alias", "Alias"}, rightValue: [2]string{"air", "air"}},
			{caption: [2]string{"التعديل ظهر في الاسمين.", "The edit appears through both names."}, leftLabel: [2]string{"المصدر", "Source"}, leftValue: [2]string{"fir", "fir"}, rightLabel: [2]string{"alias", "Alias"}, rightValue: [2]string{"fir", "fir"}},
		},
	},
})

var ch34LedgerArabic, ch34LedgerEnglish = ch34BuildLesson(ch34LessonSpec{
	prefix: "c342.",
	mission: [2]string{
		"كمّل دفتر عمليات بيكتب التغيير قبل ما يحدّث الذاكرة. عند الفتح، أعد بناء الحالة من السجل وتعامل مع ذيل ناقص أو بيانات كاملة تالفة بشكل مختلف.",
		"Complete an operations ledger that writes a change before updating memory. On open, rebuild state from the log and treat a short tail differently from complete corrupt data.",
	},
	ideaText: [2]string{
		"السجل الثنائي بيتقري على دفعات كاملة. ~io.ReadFull~ بيوضح إذا الداتا كملت أو انتهت في نص القراءة؛ عدد البايتات يفرق.",
		"A binary log is read in complete chunks. ~io.ReadFull~ distinguishes a filled read from one that ends early; the byte count matters.",
	},
	idea: ch34ExampleSpec{
		id: "e342-short-read", imports: []string{"fmt", "io", "strings"},
		code: `func main() {
	r := strings.NewReader("abc")
	buf := make([]byte, 4)
	n, err := io.ReadFull(r, buf)
	fmt.Println(n, err == io.ErrUnexpectedEOF)
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"المصدر فيه 3 بايتات.", "The source has three bytes."}, vars: []LessonVariable{{Name: "r", Value: "reader(abc)"}}},
			{line: 3, caption: [2]string{"البافر طالب 4 بايتات.", "The buffer asks for four bytes."}, vars: []LessonVariable{{Name: "buf", Value: "[0 0 0 0]"}}},
			{line: 4, caption: [2]string{"اتقرت 3 بايتات بس.", "Only three bytes were read."}, vars: []LessonVariable{{Name: "n", Value: "3"}}},
			{line: 5, caption: [2]string{"النهاية جات في نص القراءة.", "Input ended during the read."}, vars: []LessonVariable{{Name: "n", Value: "3"}}, output: "3 true"},
		},
	},
	tryText: [2]string{
		"الكود بيحفظ رقم في 4 بايتات بترتيب big-endian. توقّع الرقم بعد قراءته.",
		"The code stores a number in four big-endian bytes. Predict the number after reading them.",
	},
	question: [2]string{"إيه الناتج؟", "What is printed?"},
	options:  [2][]string{{"9", "4", "0"}, {"9", "4", "0"}}, answer: 0,
	feedback: [2][]string{
		{"الـencoding والـdecoding استخدموا نفس الترتيب.", "عدد البايتات مش الرقم المخزّن.", "PutUint32 كتب القيمة قبل القراءة."},
		{"Encoding and decoding use the same byte order.", "The byte count is not the stored number.", "PutUint32 wrote a value before the read."},
	},
	revealText: [2]string{"تابع كتابة الرقم، ثم فكّ البايتات بنفس الـbyte order.", "Follow the number being stored, then decode with the same byte order."},
	reveal: ch34ExampleSpec{
		id: "e342-binary-roundtrip", imports: []string{"encoding/binary", "fmt"},
		code: `func main() {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 9)
	fmt.Println(binary.BigEndian.Uint32(data))
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"جهزنا أربع خانات.", "Four byte slots are ready."}, vars: []LessonVariable{{Name: "data", Value: "[0 0 0 0]"}}},
			{line: 3, caption: [2]string{"كتبنا الرقم بالترتيب الكبير.", "The number is stored big-endian."}, vars: []LessonVariable{{Name: "data", Value: "[0 0 0 9]"}}},
			{line: 4, caption: [2]string{"نفس الترتيب رجّع الرقم.", "The same order restores it."}, vars: []LessonVariable{{Name: "data", Value: "[0 0 0 9]"}}, output: "9"},
		},
	},
	watch: [2]string{
		"ذيل ناقص معناه كتابة اتقطعت؛ record كامل وchecksum غلط معناه فساد. ما تخلطش بين الاتنين.",
		"A short tail means an interrupted write; a complete record with a bad checksum is corruption. Keep those cases distinct.",
	},
	plan: [2]string{
		"حدّد حدود كل record قبل ما تطبقه. افصل القراءة والتحقق عن تحديث الذاكرة، وراجع إن الاسترجاع يعيد نفس الحالة.",
		"Identify each record boundary before applying it. Separate reading and validation from memory updates, then check that replay rebuilds the same state.",
	},
	checkQuestion: [2]string{"سجل كامل ببصمة غير مطابقة يتعامل كإيه؟", "How should a complete record with a bad checksum be treated?"},
	checkOptions: [2][]string{
		{"فساد يوقف الفتح.", "ذيل ناقص يتقص.", "تغيير يتتجاهل بصمت."},
		{"Corruption that stops opening.", "A short tail to trim.", "A change to ignore silently."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"سجل كامل غير سليم ماينفعش يتخمّن.", "القص مخصص للبيانات الناقصة في الآخر.", "التجاهل الصامت يضيّع الحالة."},
		{"A complete invalid record cannot be guessed around.", "Trimming is for incomplete trailing data.", "Silently ignoring it loses state."},
	},
	turn: [2]string{
		"ابدأ من بنية ~Ledger~ الفاضية وتوقيعاتها. قسّم الشغل لتسجيل، فتح، واسترجاع، وخلي ~main~ فاضية.",
		"Start with the empty ~Ledger~ type and its signatures. Separate writing, opening, and replay, and leave ~main~ empty.",
	},
	hints: [2][3]string{
		{"السجل بيسمح بإعادة بناء الحالة بعد إعادة التشغيل.", "رتّب تفكيرك من bytes إلى record ثم حالة الذاكرة.", "ابدأ من ~Ledger~ وتوقيعات ~OpenLedger~ و~Lookup~ في الـstarter."},
		{"A log lets the program rebuild state after a restart.", "Reason from bytes to records, then to in-memory state.", "Start at ~Ledger~ and the ~OpenLedger~ and ~Lookup~ signatures in the starter."},
	},
	sandbox: ch34SandboxSpec{
		id:      "ch34-write-ahead-log",
		text:    [2]string{"قارن سجل كامل بذيل ناقص وبايت متغير.", "Compare complete, short, and changed records."},
		title:   [2]string{"قراءة السجل", "Reading a record"},
		advance: [2]string{"الحالة الجاية", "Next case"}, reset: [2]string{"من الأول", "Reset"},
		frames: []ch34SandboxFrame{
			{caption: [2]string{"السجل كامل وفحصه مطابق.", "The record is complete and checksum matches."}, leftLabel: [2]string{"البايتات", "Bytes"}, leftValue: [2]string{"01 02 03 04", "01 02 03 04"}, rightLabel: [2]string{"النتيجة", "Result"}, rightValue: [2]string{"سليم", "intact"}},
			{caption: [2]string{"الذيل ناقص؛ السجل مش كامل.", "The tail is short; the record is incomplete."}, leftLabel: [2]string{"البايتات", "Bytes"}, leftValue: [2]string{"01 02 03", "01 02 03"}, rightLabel: [2]string{"النتيجة", "Result"}, rightValue: [2]string{"غير مكتمل", "incomplete"}},
			{caption: [2]string{"الطول كامل لكن الفحص مختلف.", "The length is complete, but checksum differs."}, leftLabel: [2]string{"البايتات", "Bytes"}, leftValue: [2]string{"01 02 FF 04", "01 02 FF 04"}, rightLabel: [2]string{"النتيجة", "Result"}, rightValue: [2]string{"تلف", "corrupt"}},
		},
	},
})

var ch34CompactionArabic, ch34CompactionEnglish = ch34BuildLesson(ch34LessonSpec{
	prefix: "c343.",
	mission: [2]string{
		"كمّل على دفتر المهمة السابقة. فرّغ حالته لجداول ثابتة مرتبة، وبعدها ادمج الجداول من غير ما قراءة مفتاح ترجع نتيجة قديمة أو ممسوحة.",
		"Continue the previous ledger. Flush its state into sorted immutable tables, then merge tables without letting a read return an old or deleted value.",
	},
	ideaText: [2]string{
		"الـSSTable بتكون مرتبة عشان البحث الثنائي. ترتيب المفاتيح قبل البحث مهم؛ الترتيب اللي في slice هو اللي بيتغير هنا.",
		"An SSTable is sorted so it can be searched efficiently. Sort keys before searching; the slice order changes here.",
	},
	idea: ch34ExampleSpec{
		id: "e343-sort-keys", imports: []string{"fmt", "sort"},
		code: `func main() {
	keys := []string{"fuel", "air"}
	sort.Strings(keys)
	fmt.Println(keys[0], keys[1])
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"المفاتيح لسه مش مرتبة.", "The keys are not sorted yet."}, vars: []LessonVariable{{Name: "keys", Value: "[fuel air]"}}},
			{line: 3, caption: [2]string{"sort رتبت الـslice.", "sort reordered the slice."}, vars: []LessonVariable{{Name: "keys", Value: "[air fuel]"}}},
			{line: 4, caption: [2]string{"اطبع الترتيب الجديد.", "Print the new order."}, vars: []LessonVariable{{Name: "keys", Value: "[air fuel]"}}, output: "air fuel"},
		},
	},
	tryText: [2]string{
		"علامة ~Dead~ بتفضل موجودة حتى لو مفيش قيمة حية في المدخل. توقّع قيمة الحقل المطبوع.",
		"A ~Dead~ marker stays present even when the entry has no live value. Predict the printed field.",
	},
	question: [2]string{"إيه الناتج؟", "What is printed?"},
	options:  [2][]string{{"true", "false", "القيمة القديمة"}, {"true", "false", "The old value"}}, answer: 0,
	feedback: [2][]string{
		{"المدخل اتنشأ وعلامة الحذف true.", "القيمة الافتراضية لـbool هي false.", "المثال بيطبع الحقل مش قيمة قديمة."},
		{"The entry is created with its delete marker set.", "The zero value of bool is false.", "The example prints a field, not an old value."},
	},
	revealText: [2]string{"شوف قيمة العلامة في entry تدريب صغيرة.", "Watch the marker value in one small sample entry."},
	reveal: ch34ExampleSpec{
		id: "e343-tombstone", imports: []string{"fmt"},
		code: `type Slot struct {
	Dead bool
}
func main() {
	slot := Slot{Dead: true}
	fmt.Println(slot.Dead)
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 1, caption: [2]string{"النوع بيحفظ حالة المدخل.", "The type stores entry state."}},
			{line: 4, caption: [2]string{"العلامة موجودة في الحالة.", "The marker is part of the state."}, vars: []LessonVariable{{Name: "slot.Dead", Value: "true"}}},
			{line: 5, caption: [2]string{"القيمة true بتظهر.", "The value true is printed."}, vars: []LessonVariable{{Name: "slot.Dead", Value: "true"}}, output: "true"},
		},
	},
	watch: [2]string{
		"الملفات القديمة ثابتة؛ حذف قيمة من الذاكرة لوحده ممكن يرجعها في القراءة. والـtombstone مايتشالش قبل ما التاريخ الأقدم يبقى مغطى.",
		"Old tables are immutable; deleting only from memory can resurrect a value. Keep a tombstone until older history is covered.",
	},
	plan: [2]string{
		"فكّر في أولوية الحالة: الذاكرة أولًا ثم الجداول من الأحدث. في الدمج احتفظ بآخر حالة لكل مفتاح، وراجع أثر علامات الحذف.",
		"Think about state precedence: memory first, then tables from newest to oldest. During compaction, preserve each key's latest state and account for deletes.",
	},
	checkQuestion: [2]string{"ليه tombstone مهم مع جدول قديم ثابت؟", "Why does a tombstone matter when an old table is immutable?"},
	checkOptions: [2][]string{
		{"عشان يمنع رجوع قيمة محذوفة.", "عشان يرتب أسماء الملفات.", "عشان يضغط حجم كل سجل."},
		{"It prevents a deleted value from returning.", "It sorts file names.", "It compresses every record."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"العلامة بتوقف القراءة عند المفتاح المحذوف.", "ترتيب الملفات له خطوة منفصلة.", "علامة الحذف مش خوارزمية ضغط."},
		{"The marker stops a read at a deleted key.", "File ordering is a separate concern.", "A delete marker is not compression."},
	},
	turn: [2]string{
		"راجع ما ورثته من ~34-2-wal~، وبعدها نفّذ امتداد التخزين. قارن نتائج القراءة قبل وبعد كل عملية، وسيب ~main~ فاضية.",
		"Review what carries over from ~34-2-wal~, then extend the storage. Compare reads before and after each operation, and leave ~main~ empty.",
	},
	hints: [2][3]string{
		{"العلامة الصريحة بتخلي الحذف حالة قابلة للقراءة.", "افصل بين تفريغ الحالة الحالية ودمج التاريخ القديم.", "راجع حقول ~Ledger~ الموروثة وتوقيعات ~Flush~ و~Compact~ في الـstarter."},
		{"An explicit marker makes deletion a readable state.", "Separate flushing the current state from merging old history.", "Review the inherited ~Ledger~ fields and ~Flush~ and ~Compact~ signatures in the starter."},
	},
	sandbox: ch34SandboxSpec{
		id:      "ch34-compaction",
		text:    [2]string{"تابع علامة الحذف قدام قيمة أقدم.", "Trace a tombstone beside an older value."},
		title:   [2]string{"دمج الأجيال", "Merging generations"},
		advance: [2]string{"الخطوة الجاية", "Next event"}, reset: [2]string{"من الأول", "Reset"},
		frames: []ch34SandboxFrame{
			{caption: [2]string{"العلامة الجديدة بتخفي القديم.", "The newer marker hides the older value."}, leftLabel: [2]string{"القديم", "Older table"}, leftValue: [2]string{"fuel", "fuel"}, rightLabel: [2]string{"القراءة", "Read"}, rightValue: [2]string{"محذوف", "deleted"}},
			{caption: [2]string{"شيل العلامة بدري؛ القديم يظهر.", "Drop the marker early; the old value returns."}, leftLabel: [2]string{"القديم", "Older table"}, leftValue: [2]string{"fuel", "fuel"}, rightLabel: [2]string{"القراءة", "Read"}, rightValue: [2]string{"fuel", "fuel"}},
			{caption: [2]string{"بعد دمج القديم، الحذف ثابت.", "After merging old tables, deletion remains."}, leftLabel: [2]string{"القديم", "Older table"}, leftValue: [2]string{"اتشال", "removed"}, rightLabel: [2]string{"القراءة", "Read"}, rightValue: [2]string{"محذوف", "deleted"}},
			{caption: [2]string{"دلوقتي اختفاء العلامة آمن.", "Now removing the marker is safe."}, leftLabel: [2]string{"القديم", "Older table"}, leftValue: [2]string{"اتشال", "removed"}, rightLabel: [2]string{"القراءة", "Read"}, rightValue: [2]string{"مش موجود", "absent"}},
		},
	},
})

var ch34DepotArabic, ch34DepotEnglish = ch34BuildLesson(ch34LessonSpec{
	prefix: "c344.",
	mission: [2]string{
		"قسّم الملفات لقطع ووزّع منها نسخ على كبسولات تخزين. لازم تقدر تقرأ نسخة سليمة لو غيرها وقع أو اتلف، وتحاول تصلّح النسخ الناقصة.",
		"Split files into chunks and place replicas across storage pods. Reads should survive a missing or damaged copy, and repair should restore missing replicas.",
	},
	ideaText: [2]string{
		"بصمة المحتوى بتدي نفس الاسم لنفس البايتات. تابع hash ثابت مرتين؛ المساواة هنا بتعبر عن بصمة متطابقة.",
		"A content fingerprint gives the same bytes the same name. Follow one hash twice; equality means matching fingerprints.",
	},
	idea: ch34ExampleSpec{
		id: "e344-content-hash", imports: []string{"crypto/sha256", "fmt"},
		code: `func main() {
	a := sha256.Sum256([]byte("ore"))
	b := sha256.Sum256([]byte("ore"))
	fmt.Println(a == b, len(a))
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"البصمة الأولى اتعملت.", "The first fingerprint is ready."}},
			{line: 3, caption: [2]string{"نفس البايتات لها نفس البصمة.", "Same bytes produce the same hash."}},
			{line: 4, caption: [2]string{"البصمة 32 بايت ومتطابقة.", "The 32-byte hashes match."}, vars: []LessonVariable{{Name: "a == b", Value: "true"}, {Name: "len(a)", Value: "32"}}, output: "true 32"},
		},
	},
	tryText: [2]string{
		"بصمة SHA-256 بتتحول لنص hex. توقّع طول النص اللي هيتطبع.",
		"A SHA-256 digest converts to hex text. Predict the length printed.",
	},
	question: [2]string{"كام حرف هيظهر؟", "How many characters appear?"},
	options:  [2][]string{{"32", "64", "حسب اسم الكبسولة"}, {"32", "64", "It depends on the pod name"}}, answer: 1,
	feedback: [2][]string{
		{"كل بايت بيتكتب بحرفين hex.", "النص أطول من عدد البايتات.", "الاسم هنا طالع من المحتوى نفسه."},
		{"Each byte uses two hex characters.", "The text is longer than the byte count.", "This name comes from the content itself."},
	},
	revealText: [2]string{"حوّل البصمة لتنسيق hex، وبعدها اطبع طول الاسم.", "Convert the digest to hex, then print the identifier length."},
	reveal: ch34ExampleSpec{
		id: "e344-hex-length", imports: []string{"crypto/sha256", "encoding/hex", "fmt"},
		code: `func main() {
	sum := sha256.Sum256([]byte("ore"))
	id := hex.EncodeToString(sum[:])
	fmt.Println(len(id))
}`,
		frames: []ch34FrameSpec{
			{line: 0, caption: [2]string{"قبل التنفيذ.", "Before execution."}},
			{line: 2, caption: [2]string{"حسبنا البصمة.", "The digest was computed."}},
			{line: 3, caption: [2]string{"حوّلناها لنص hex.", "The digest is now hex text."}},
			{line: 4, caption: [2]string{"طول الاسم 64.", "The identifier length is 64."}, vars: []LessonVariable{{Name: "len(id)", Value: "64"}}, output: "64"},
		},
	},
	watch: [2]string{
		"وجود نسخة مش كفاية؛ لازم تتأكد إن بايتاتها تطابق بصمتها. فشل كبسولة واحدة مايعنيش إن كل النسخ فشلت.",
		"A returned copy is not enough; verify its bytes against the fingerprint. One failed pod does not mean every replica failed.",
	},
	plan: [2]string{
		"افصل القراءة من المصدر عن تقسيم القطع وتوزيعها. بعد التخزين، فكّر في التحقق عند التحميل وفي إصلاح النسخ من مصدر سليم.",
		"Separate reading the input from chunking and placement. After storage, consider verification on load and repair from an intact source.",
	},
	checkQuestion: [2]string{"نسخة رجعت من ~Fetch~ من غير error، بس بصمتها غلط؛ تعمل إيه؟", "A ~Fetch~ copy has no error, but its fingerprint is wrong; what then?"},
	checkOptions: [2][]string{
		{"اعتبرها تالفة وجرّب نسخة تانية.", "اكتبها للمخرج مباشرة.", "غيّر الاسم عشان يطابقها."},
		{"Treat it as damaged and try another copy.", "Write it to the output immediately.", "Change the name to match it."},
	}, checkAnswer: 0,
	checkFeedback: [2][]string{
		{"نجاح النقل مش دليل سلامة المحتوى.", "التحقق لازم يسبق الكتابة.", "الاسم هو اللي بيثبت البايتات."},
		{"A successful fetch does not prove intact content.", "Verification comes before writing it out.", "The name is the fingerprint of the bytes."},
	},
	turn: [2]string{
		"اقرأ interface ~Pod~ والعقد اللي على ~Depot~. ابدأ بحركة بيانات واحدة ووسّع التفكير للنسخ والإصلاح؛ ~main~ تفضل فاضية.",
		"Read the ~Pod~ interface and the ~Depot~ contract. Start with one data flow, then reason about replicas and repair; keep ~main~ empty.",
	},
	hints: [2][3]string{
		{"الاسم المبني من المحتوى يسهّل اكتشاف التكرار والتلف.", "رتّب الخطوات من chunk إلى placement ثم verification وrepair.", "ابدأ من ~Pod.Fetch~ و~Pod.Store~ وتعريف ~Depot~ في الـstarter."},
		{"A content-derived name helps detect duplicates and damage.", "Think from chunking to placement, then verification and repair.", "Start at ~Pod.Fetch~, ~Pod.Store~, and the ~Depot~ type in the starter."},
	},
	sandbox: ch34SandboxSpec{
		id:      "ch34-distributed-storage",
		text:    [2]string{"اتبع نسخة تالفة وقراءة نسخة بديلة وإصلاحها.", "Trace a damaged copy, fallback, and repair."},
		title:   [2]string{"نسخ موزعة", "Distributed copies"},
		advance: [2]string{"الحدث الجاي", "Next event"}, reset: [2]string{"من الأول", "Reset"},
		frames: []ch34SandboxFrame{
			{caption: [2]string{"النسخ الثلاث سليمة.", "All three copies are intact."}, leftLabel: [2]string{"A / B / C", "A / B / C"}, leftValue: [2]string{"سليمة", "intact"}, rightLabel: [2]string{"قراءة", "Read"}, rightValue: [2]string{"تحقق", "verified"}},
			{caption: [2]string{"النسخة A بقت تالفة.", "Copy A is now damaged."}, leftLabel: [2]string{"A / B / C", "A / B / C"}, leftValue: [2]string{"تالف / سليم / سليم", "bad / good / good"}, rightLabel: [2]string{"قراءة", "Read"}, rightValue: [2]string{"جرّب نسخة تانية", "try another copy"}},
			{caption: [2]string{"النسخة B اجتازت التحقق.", "Copy B passed verification."}, leftLabel: [2]string{"القراءة", "Read"}, leftValue: [2]string{"B", "B"}, rightLabel: [2]string{"البصمة", "Fingerprint"}, rightValue: [2]string{"مطابقة", "matches"}},
			{caption: [2]string{"الإصلاح نسخ بايتات سليمة.", "Repair copied intact bytes."}, leftLabel: [2]string{"A / B / C", "A / B / C"}, leftValue: [2]string{"سليم / سليم / سليم", "good / good / good"}, rightLabel: [2]string{"المصدر", "Source"}, rightValue: [2]string{"B", "B"}},
		},
	},
})

var ch34LessonArabic, ch34LessonEnglish = func() (map[string][]LessonStep, map[string][]LessonStep) {
	return map[string][]LessonStep{
		"34-1-kv-store":    ch34VaultArabic,
		"34-2-wal":         ch34LedgerArabic,
		"34-3-compaction":  ch34CompactionArabic,
		"34-4-chunk-store": ch34DepotArabic,
	}, map[string][]LessonStep{
		"34-1-kv-store":    ch34VaultEnglish,
		"34-2-wal":         ch34LedgerEnglish,
		"34-3-compaction":  ch34CompactionEnglish,
		"34-4-chunk-store": ch34DepotEnglish,
	}
}()
