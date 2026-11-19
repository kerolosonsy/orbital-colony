package levels

func ch23SandboxPair(arabic, english string) map[string]string {
	return map[string]string{"ar": arabic, "en": english}
}

func ch23SandboxRow(arabic, english, displayValue string) map[string]any {
	return map[string]any{"label": ch23SandboxPair(arabic, english), "value": displayValue}
}

func ch23SandboxFrame(arabic, english string, rows ...map[string]any) map[string]any {
	return map[string]any{"caption": ch23SandboxPair(arabic, english), "rows": rows}
}

func ch23SandboxConfig(action map[string]string, frames ...map[string]any) map[string]any {
	return map[string]any{"action": action, "frames": frames}
}

func ch23Playground(id, kind, description string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: description, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

var ch23ReadFullSandboxConfig = ch23SandboxConfig(ch23SandboxPair("الخطوة التالية", "Next step"),
	ch23SandboxFrame("ابدأ بمصدر ثابت.", "Start with a fixed source.",
		ch23SandboxRow("المصدر", "Source", "[9, 4]"), ch23SandboxRow("المخزن", "Buffer", "[_, _]"), ch23SandboxRow("القراءة", "Read", "n=0; waiting")),
	ch23SandboxFrame("أول بايت وصل.", "The first byte arrived.",
		ch23SandboxRow("المصدر", "Source", "[9, 4]"), ch23SandboxRow("المخزن", "Buffer", "[9, _]"), ch23SandboxRow("القراءة", "Read", "n=1; still reading")),
	ch23SandboxFrame("المخزن اكتمل.", "The buffer is full.",
		ch23SandboxRow("المصدر", "Source", "[9, 4]"), ch23SandboxRow("المخزن", "Buffer", "[9, 4]"), ch23SandboxRow("القراءة", "Read", "n=2; err=nil")))

var ch23AADSandboxConfig = ch23SandboxConfig(ch23SandboxPair("غيّر بيانات الفتح", "Change open data"),
	ch23SandboxFrame("البيانات المصاحبة متطابقة.", "Associated data matches.",
		ch23SandboxRow("المحتوى", "Payload", "encrypted"), ch23SandboxRow("وقت الختم", "Seal AAD", "label"), ch23SandboxRow("وقت الفتح", "Open AAD", "label"), ch23SandboxRow("التحقق", "Check", "accepted")),
	ch23SandboxFrame("البيانات المصاحبة مختلفة.", "Associated data differs.",
		ch23SandboxRow("المحتوى", "Payload", "encrypted"), ch23SandboxRow("وقت الختم", "Seal AAD", "label"), ch23SandboxRow("وقت الفتح", "Open AAD", "other"), ch23SandboxRow("التحقق", "Check", "rejected")))

var ch23ReadLimitSandboxConfig = ch23SandboxConfig(ch23SandboxPair("وسّع حد القراءة", "Increase read limit"),
	ch23SandboxFrame("اقرأ جزءًا محدودًا.", "Read a bounded part.",
		ch23SandboxRow("الحد", "Limit", "2 bytes"), ch23SandboxRow("المقروء", "Read", "pl"), ch23SandboxRow("المتبقي", "Remaining", "4 bytes")),
	ch23SandboxFrame("اسمح ببايتات أكتر.", "Allow more bytes.",
		ch23SandboxRow("الحد", "Limit", "4 bytes"), ch23SandboxRow("المقروء", "Read", "plan"), ch23SandboxRow("المتبقي", "Remaining", "2 bytes")),
	ch23SandboxFrame("اقرأ النص كله.", "Read the whole text.",
		ch23SandboxRow("الحد", "Limit", "6 bytes"), ch23SandboxRow("المقروء", "Read", "planet"), ch23SandboxRow("المتبقي", "Remaining", "0 bytes")))

var ch23OrderSandboxConfig = ch23SandboxConfig(ch23SandboxPair("ترتيب إضافة تاني", "Another insertion order"),
	ch23SandboxFrame("ترتيب الإضافة اتغيّر.", "Insertion order changed.",
		ch23SandboxRow("الإضافة", "Inserted", "[bay/z, bay/a]"), ch23SandboxRow("بعد الترتيب", "Sorted", "[bay/a, bay/z]")),
	ch23SandboxFrame("النتيجة المرتبة ثابتة.", "The sorted result stays stable.",
		ch23SandboxRow("الإضافة", "Inserted", "[bay/a, bay/z]"), ch23SandboxRow("بعد الترتيب", "Sorted", "[bay/a, bay/z]")))

var ch23CSVSandboxConfig = ch23SandboxConfig(ch23SandboxPair("السجل التالي", "Next record"),
	ch23SandboxFrame("اقرأ صف العناوين.", "Read the header row.",
		ch23SandboxRow("السطر", "Line", "station,reading"), ch23SandboxRow("الحقول", "Fields", "[station] [reading]")),
	ch23SandboxFrame("كل عمود بقى field.", "Each column is a field.",
		ch23SandboxRow("السطر", "Line", "sensor-2,7"), ch23SandboxRow("الحقول", "Fields", "[sensor-2] [7]")))

func ch23Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch23Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch23QuestionExample(id, prompt, code string, options []string, answer int, feedback []string) LessonBlock {
	return ch23Question(id, prompt+"\n\n~~~go\n"+code+"\n~~~", options, answer, feedback)
}

func ch23Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

func ch23EnglishFrames(frames []LessonFrame) []LessonFrame {
	translations := map[string]string{
		"قبل التنفيذ.":                        "Before run.",
		"مساحة بايتات للتوضيح.":               "Illustrative byte buffer.",
		"المصدر يملأ المساحة.":                "The source fills the buffer.",
		"اطبع الحالة من غير السر.":            "Print status without the secret.",
		"مصدر ثابت للتجربة.":                  "Fixed source for the demo.",
		"جهّز مكان القراءة.":                  "Prepare the read buffer.",
		"املأ البايتين حتى لو القراءة قصيرة.": "Fill both bytes despite short reads.",
		"القيم وصلت كاملة.":                   "Both values arrived.",
		"مفتاح صالح للتوضيح.":                 "Valid demo key.",
		"كوّن AEAD.":                          "Build the AEAD.",
		"طول nonce والـtag.":                  "Nonce and tag lengths.",
		"أنشئ cipher بمفتاح صالح.":            "Create a cipher with a valid key.",
		"أضف تحقق GCM.":                       "Add GCM authentication.",
		"جهّز nonce مناسبة.":                  "Prepare a matching nonce.",
		"الختم يشمل البيانات المصاحبة.":       "The seal covers associated data.",
		"افتح ببيانات مصاحبة مختلفة.":         "Open with different associated data.",
		"فشل التحقق ظاهر.":                    "Authentication failure is visible.",
		"محتوى خام للتوضيح.":                  "Illustrative raw content.",
		"بصمة المحتوى الخام.":                 "Hash the raw content.",
		"ضغط نفس البايتات.":                   "Compress the same bytes.",
		"إقفال zlib يكمل stream.":             "Closing zlib completes the stream.",
		"بصمة البايتات المضغوطة.":             "Hash the compressed bytes.",
		"المضغوط مختلف عن الخام.":             "Compressed differs from raw.",
		"حد القراءة بثلاث بايتات.":            "Limit the read to three bytes.",
		"اقرأ الجزء المحدود.":                 "Read the limited portion.",
		"النتيجة قصيرة ومكتملة.":              "The short result is complete.",
		"ترتيب إضافة للتوضيح.":                "Insertion order for the demo.",
		"رتّب الأسماء.":                       "Sort the names.",
		"اعرض ترتيب ثابت.":                    "Show stable ordering.",
		"الـstruct يحدد الترتيب.":             "The struct defines the order.",
		"حوّل القيمة لـJSON.":                 "Encode the value as JSON.",
		"المفاتيح بنفس ترتيب الحقول.":         "Keys follow the field order.",
		"قراءات للتوضيح.":                     "Illustrative readings.",
		"ثبّت ترتيب العرض.":                   "Make display order stable.",
		"الترتيب بقى ثابت.":                   "Ordering is now stable.",
		"قارئ CSV يحفظ شكل الحقول.":           "CSV reader keeps field structure.",
		"اقرأ العنوان كسجل.":                  "Read the header as a record.",
		"اقرأ السطر اللي بعده.":               "Read the next row.",
		"CSV رجّع الحقول والأخطاء.":           "CSV returned fields and errors.",
	}
	translated := make([]LessonFrame, len(frames))
	for i, frame := range frames {
		translated[i] = frame
		translated[i].Caption = translations[frame.Caption]
	}
	return translated
}

const ch23Example231Idea = `func main() {
	buffer := make([]byte, 2)
	n, err := rand.Read(buffer)
	fmt.Println(n == len(buffer), err == nil)
}`

const ch23Example231Try = `func main() {
	reader := bytes.NewReader([]byte{9, 4})
	var sample [2]byte
	n, err := io.ReadFull(reader, sample[:])
	fmt.Println(n, sample, err)
}`

var ch23Frames231Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "مساحة بايتات للتوضيح."},
	{Line: 3, Caption: "المصدر يملأ المساحة."},
	{Line: 4, Caption: "اطبع الحالة من غير السر.", Output: "true true"},
}

var ch23Frames231Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "مصدر ثابت للتجربة."},
	{Line: 3, Caption: "جهّز مكان القراءة."},
	{Line: 4, Caption: "املأ البايتين حتى لو القراءة قصيرة."},
	{Line: 5, Caption: "القيم وصلت كاملة.", Output: "2 [9 4] <nil>"},
}

var ch23Lesson231Arabic = []LessonStep{
	{ID: "c231.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch23Markdown("c231.mission", "أكواد الدخول محتاجة عشوائية مناسبة للأسرار، واختبارات تقدر تعيدها.\n\nالـstarter فيه ~IssueAccess~؛ افصل مصدر البايتات عن قرار اختيار الحروف."),
	}},
	{ID: "c231.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch23Walkthrough("c231.random-trace", "~crypto/rand~ للمصدر الحقيقي؛ ~math/rand~ مش مناسب للأسرار.\n\nالـreader المحقون يخلي السلوك قابل للتكرار في الاختبار، وخطأ القراءة لازم يوقف الاستخدام.", ch23Example231Idea, "e231-random-read", []string{"crypto/rand", "fmt"}, ch23Frames231Idea),
	}},
	{ID: "c231.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch23QuestionExample("c231.readfull-predict", "توقّع الناتج لما المصدر يسلّم بايتين.", ch23Example231Try, []string{"2 [9 4] <nil>", "1 [9 0] EOF", "0 [0 0] <nil>"}, 0, []string{
			"الـreader فيه البايتين المطلوبين.",
			"ReadFull كملت العدد المطلوب.",
			"القراءة نجحت ووصلت القيمتين.",
		}),
		ch23Walkthrough("c231.readfull-reveal", "تابع قراءة المصدر الثابت لحد امتلاء المكان.", ch23Example231Try, "e231-readfull", []string{"bytes", "fmt", "io"}, ch23Frames231Try),
		ch23Playground("ch23-read-full", "ch23-read-full", "تابع امتلاء المخزن.", ch23ReadFullSandboxConfig),
	}},
	{ID: "c231.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch23Markdown("c231.watch", "- متطبعش الكود السري في log.\n- فشل المصدر يرجّع كود فاضي، مش جزء منه."),
	}},
	{ID: "c231.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch23Markdown("c231.plan", "راجع الطلب قبل استهلاك المصدر. بعد كده فكّر إزاي تضمن كل الفئات المطلوبة، وإزاي تخلّي نفس المصدر المحقون يتحكم في كل اختيار."),
		ch23Question("c231.reader-question", "إيه اللي يكمل قراءة قصيرة؟", []string{"io.ReadFull", "fmt.Println", "time.Now"}, 0, []string{
			"الأداة دي بتكمل لحد امتلاء الـbuffer أو خطأ.",
			"الطباعة مش قراءة من المصدر.",
			"الوقت مش بديل عن بايتات العشوائية.",
		}),
	}},
	{ID: "c231.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch23Markdown("c231.turn", "ابدأ بمصدر bytes ثابت. راجع حالة EOF قبل تجربة المصدر الحقيقي."),
	}},
}

var ch23Lesson231English = []LessonStep{
	{ID: "c231.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch23Markdown("c231.mission", "Access codes need randomness suitable for secrets and tests that can be repeated.\n\nThe starter has ~IssueAccess~; separate the byte source from character selection."),
	}},
	{ID: "c231.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch23Walkthrough("c231.random-trace", "Use ~crypto/rand~ for a real source; ~math/rand~ is unsuitable for secrets.\n\nAn injected reader makes behavior repeatable in tests, and read errors must stop use.", ch23Example231Idea, "e231-random-read", []string{"crypto/rand", "fmt"}, ch23EnglishFrames(ch23Frames231Idea)),
	}},
	{ID: "c231.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch23QuestionExample("c231.readfull-predict", "Predict the output when the source supplies two bytes.", ch23Example231Try, []string{"2 [9 4] <nil>", "1 [9 0] EOF", "0 [0 0] <nil>"}, 0, []string{
			"The reader contains the two requested bytes.",
			"ReadFull completed the requested length.",
			"The read succeeded and returned both values.",
		}),
		ch23Walkthrough("c231.readfull-reveal", "Follow the fixed source until the buffer is full.", ch23Example231Try, "e231-readfull", []string{"bytes", "fmt", "io"}, ch23EnglishFrames(ch23Frames231Try)),
		ch23Playground("ch23-read-full", "ch23-read-full", "Watch the buffer fill.", ch23ReadFullSandboxConfig),
	}},
	{ID: "c231.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch23Markdown("c231.watch", "- Never print the secret code to a log.\n- A source failure returns an empty code, not a partial one."),
	}},
	{ID: "c231.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch23Markdown("c231.plan", "Validate the request before consuming the source.\n\nThen reason about how every requested category is guaranteed and how one injected source controls each choice."),
		ch23Question("c231.reader-question", "What completes a short read?", []string{"io.ReadFull", "fmt.Println", "time.Now"}, 0, []string{
			"It reads until the buffer is full or an error occurs.",
			"Printing does not read from the source.",
			"Time is not a replacement for random bytes.",
		}),
	}},
	{ID: "c231.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch23Markdown("c231.turn", "Start with a fixed bytes source. Review EOF behavior before trying the real source."),
	}},
}

const ch23Example232Idea = `func main() {
	block, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	fmt.Println(aead.NonceSize(), aead.Overhead())
}`

const ch23Example232Try = `func main() {
	block, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	nonce := make([]byte, aead.NonceSize())
	sealed := aead.Seal(nil, nonce, []byte("ok"), []byte("label"))
	_, err = aead.Open(nil, nonce, sealed, []byte("other"))
	fmt.Println(err != nil)
}`

var ch23Frames232Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "مفتاح صالح للتوضيح."},
	{Line: 5, Caption: "كوّن AEAD."},
	{Line: 10, Caption: "طول nonce والـtag.", Output: "12 16"},
}

var ch23Frames232Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "أنشئ cipher بمفتاح صالح."},
	{Line: 5, Caption: "أضف تحقق GCM."},
	{Line: 9, Caption: "جهّز nonce مناسبة."},
	{Line: 10, Caption: "الختم يشمل البيانات المصاحبة."},
	{Line: 11, Caption: "افتح ببيانات مصاحبة مختلفة."},
	{Line: 13, Caption: "فشل التحقق ظاهر.", Output: "true"},
}

var ch23Lesson232Arabic = []LessonStep{
	{ID: "c232.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch23Markdown("c232.mission", "الشحنة محتاجة سرية وكشف أي تعديل. الـstarter فيه واجهات الختم والفتح؛ المفتاح والـmanifest والـnonce لازم يتعاملوا كمدخلات مختلفة."),
	}},
	{ID: "c232.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch23Walkthrough("c232.aead-trace", "~AES-GCM~ بيرجع AEAD للتشفير والتحقق. الـnonce مش سر، لكن لازم يبقى جديد مع نفس المفتاح؛ والـAAD بتتراجع من غير ما تتشفّر.", ch23Example232Idea, "e232-gcm-shape", []string{"crypto/aes", "crypto/cipher", "fmt"}, ch23Frames232Idea),
	}},
	{ID: "c232.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch23QuestionExample("c232.aad-predict", "إيه اللي يحصل لو الـAAD اختلفت عند الفتح؟", ch23Example232Try, []string{"يرجع النص عادي", "التحقق يفشل", "الـAAD تتشفر تلقائي"}, 1, []string{
			"البيانات المصاحبة جزء من التحقق.",
			"تغيير AAD يبطل التحقق.",
			"الـAAD بتفضل مرئية.",
		}),
		ch23Walkthrough("c232.aad-reveal", "شوف أثر تغيير البيانات المصاحبة على التحقق.", ch23Example232Try, "e232-aad-check", []string{"crypto/aes", "crypto/cipher", "fmt"}, ch23Frames232Try),
		ch23Playground("ch23-aad-check", "ch23-aad-check", "بدّل بيانات الفتح.", ch23AADSandboxConfig),
	}},
	{ID: "c232.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch23Markdown("c232.watch", "- مثال الـnonce الثابت للتوضيح بس؛ متعيدوش مع المفتاح.\n- افحص أطوال الباكت قبل slicing."),
	}},
	{ID: "c232.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch23Markdown("c232.plan", "فرّق بين بناء cipher، تجهيز nonce، والختم. عند الفتح، الخطأ معناه إن النص مش موثوق؛ مفيش كتابة قبل نجاح التحقق."),
		ch23Question("c232.aad-question", "الـAAD بتعمل إيه؟", []string{"تتشفّر وتفضل مخفية", "تفضل ظاهرة وتدخل في التحقق", "تختار طول المفتاح"}, 1, []string{
			"الـAAD لا تتشفّر.",
			"تغييرها يغيّر نتيجة التحقق.",
			"طول المفتاح يخص AES.",
		}),
	}},
	{ID: "c232.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch23Markdown("c232.turn", "جرّب فتح باكت سليمة، وبعدها غيّر manifest في مثال منفصل."),
	}},
}

var ch23Lesson232English = []LessonStep{
	{ID: "c232.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch23Markdown("c232.mission", "The cargo needs secrecy and tamper detection.\n\nThe starter has seal and open APIs; the key, manifest, and nonce are distinct inputs."),
	}},
	{ID: "c232.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch23Walkthrough("c232.aead-trace", "~AES-GCM~ returns an AEAD for encryption and authentication.\n\nA nonce is not secret, but must be new for the same key; AAD is checked without encryption.", ch23Example232Idea, "e232-gcm-shape", []string{"crypto/aes", "crypto/cipher", "fmt"}, ch23EnglishFrames(ch23Frames232Idea)),
	}},
	{ID: "c232.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch23QuestionExample("c232.aad-predict", "What happens if AAD differs during opening?", ch23Example232Try, []string{"The plaintext returns normally", "Authentication fails", "AAD is encrypted automatically"}, 1, []string{
			"Associated data is part of authentication.",
			"Changing AAD invalidates authentication.",
			"AAD remains visible.",
		}),
		ch23Walkthrough("c232.aad-reveal", "See how changing associated data affects authentication.", ch23Example232Try, "e232-aad-check", []string{"crypto/aes", "crypto/cipher", "fmt"}, ch23EnglishFrames(ch23Frames232Try)),
		ch23Playground("ch23-aad-check", "ch23-aad-check", "Change open data.", ch23AADSandboxConfig),
	}},
	{ID: "c232.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch23Markdown("c232.watch", "- The fixed nonce is illustrative only; never reuse it with a key.\n- Check packet lengths before slicing."),
	}},
	{ID: "c232.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch23Markdown("c232.plan", "Separate cipher creation, nonce setup, and sealing.\n\nOn open, an error means the plaintext is untrusted; do not write before verification succeeds."),
		ch23Question("c232.aad-question", "What does AAD do?", []string{"It is encrypted and hidden", "It stays visible and is authenticated", "It selects key length"}, 1, []string{
			"AAD is not encrypted.",
			"Changing it changes the authentication result.",
			"AES determines the key length.",
		}),
	}},
	{ID: "c232.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch23Markdown("c232.turn", "Try opening a valid packet, then change its manifest in a separate example."),
	}},
}

const ch23Example233Idea = `func main() {
	raw := []byte("sample")
	rawID := sha1.Sum(raw)
	var packed bytes.Buffer
	writer := zlib.NewWriter(&packed)
	_, writeErr := writer.Write(raw)
	closeErr := writer.Close()
	packedID := sha1.Sum(packed.Bytes())
	fmt.Println(rawID == packedID, writeErr == nil && closeErr == nil)
}`

const ch23Example233Try = `func main() {
	limited := io.LimitReader(strings.NewReader("abcdef"), 3)
	part, err := io.ReadAll(limited)
	fmt.Println(string(part), err)
}`

var ch23Frames233Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "محتوى خام للتوضيح."},
	{Line: 3, Caption: "بصمة المحتوى الخام."},
	{Line: 5, Caption: "ضغط نفس البايتات."},
	{Line: 7, Caption: "إقفال zlib يكمل stream."},
	{Line: 8, Caption: "بصمة البايتات المضغوطة."},
	{Line: 9, Caption: "المضغوط مختلف عن الخام.", Output: "false true"},
}

var ch23Frames233Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "حد القراءة بثلاث بايتات."},
	{Line: 3, Caption: "اقرأ الجزء المحدود."},
	{Line: 4, Caption: "النتيجة قصيرة ومكتملة.", Output: "abc <nil>"},
}

var ch23Lesson233Arabic = []LessonStep{
	{ID: "c233.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch23Markdown("c233.mission", "الخزنة تضغط القطع، لكن عنوان المحتوى لازم يعرّف البايتات الأصلية ونوع القطعة.\n\nالـstarter فيه قراءة وكتابة؛ خليك واضح بين الصيغة الخام وصيغة التخزين."),
	}},
	{ID: "c233.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch23Walkthrough("c233.content-trace", "اضغط البيانات للتخزين، لكن البصمة تخص التمثيل الخام. لازم ~zlib.Writer.Close~ يكمّل stream؛ نجاح Write لوحده مش كفاية.", ch23Example233Idea, "e233-compress-hash", []string{"bytes", "compress/zlib", "crypto/sha1", "fmt"}, ch23Frames233Idea),
	}},
	{ID: "c233.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch23QuestionExample("c233.limit-predict", "توقّع الجزء اللي هيوصل للقارئ.", ch23Example233Try, []string{"abcdef", "abc", "خطأ EOF"}, 1, []string{
			"القارئ محدود قبل نهاية النص.",
			"الحد يسمح بأول ثلاثة بايتات.",
			"القراءة المحدودة وصلت لنهايتها من غير خطأ.",
		}),
		ch23Walkthrough("c233.limit-reveal", "شوف حد القراءة وهو بيمنع استهلاك الباقي.", ch23Example233Try, "e233-read-limit", []string{"fmt", "io", "strings"}, ch23Frames233Try),
		ch23Playground("ch23-read-limit", "ch23-read-limit", "وسّع حد القراءة.", ch23ReadLimitSandboxConfig),
	}},
	{ID: "c233.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch23Markdown("c233.watch", "- لازم القراءة توصل EOF عشان checksum يتراجع.\n- افحص طول الفك قبل قبول envelope."),
	}},
	{ID: "c233.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch23Markdown("c233.plan", "افصل التحقق من المدخلات عن التخزين. عند القراءة، فك الضغط ضمن حد، وبعدها راجع التمثيل والهوية قبل إرجاع المحتوى."),
		ch23Question("c233.identity-question", "إيه اللي يثبت هوية القطعة؟", []string{"البايتات المضغوطة بس", "النوع والـpayload الخام", "اسم الملف على الجهاز"}, 1, []string{
			"محتوى مضغوط ممكن يتغير مع نفس البيانات.",
			"العنوان مبني على نوع ومحتوى القطعة.",
			"مسار التخزين مش هو العنوان.",
		}),
	}},
	{ID: "c233.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch23Markdown("c233.turn", "جرّب كتابة قطعة وقراءتها، وبعدها افحص ملف مضغوط بايظ من غير تغيير المصدر."),
	}},
}

var ch23Lesson233English = []LessonStep{
	{ID: "c233.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch23Markdown("c233.mission", "The vault compresses objects, but a content address must identify the original bytes and object kind.\n\nThe starter has read and write APIs; keep raw format distinct from storage format."),
	}},
	{ID: "c233.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch23Walkthrough("c233.content-trace", "Compress bytes for storage, but fingerprint the raw representation.\n\n~zlib.Writer.Close~ finishes the stream; a successful Write alone is not enough.", ch23Example233Idea, "e233-compress-hash", []string{"bytes", "compress/zlib", "crypto/sha1", "fmt"}, ch23EnglishFrames(ch23Frames233Idea)),
	}},
	{ID: "c233.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch23QuestionExample("c233.limit-predict", "Predict which part reaches the reader.", ch23Example233Try, []string{"abcdef", "abc", "EOF error"}, 1, []string{
			"The reader stops before the text ends.",
			"The limit allows the first three bytes.",
			"The limited read reaches its end without an error.",
		}),
		ch23Walkthrough("c233.limit-reveal", "Watch the reader limit stop before consuming the rest.", ch23Example233Try, "e233-read-limit", []string{"fmt", "io", "strings"}, ch23EnglishFrames(ch23Frames233Try)),
		ch23Playground("ch23-read-limit", "ch23-read-limit", "Increase the read limit.", ch23ReadLimitSandboxConfig),
	}},
	{ID: "c233.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch23Markdown("c233.watch", "- Read to EOF so the checksum is checked.\n- Check the decompressed length before accepting an envelope."),
	}},
	{ID: "c233.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch23Markdown("c233.plan", "Separate input validation from storage.\n\nOn read, decompress within a limit, then verify representation and identity before returning content."),
		ch23Question("c233.identity-question", "What defines the object's identity?", []string{"Only compressed bytes", "Kind and raw payload", "Its host filename"}, 1, []string{
			"Compressed bytes may differ for the same content.",
			"The address uses the object's kind and content.",
			"The storage path is not the address.",
		}),
	}},
	{ID: "c233.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch23Markdown("c233.turn", "Try writing and reading an object, then inspect a corrupt compressed file without changing the source."),
	}},
}

const ch23Example234Idea = `func main() {
	names := []string{"bay/z", "bay/a"}
	slices.Sort(names)
	fmt.Println(strings.Join(names, ", "))
}`

const ch23Example234Try = `type Stamp struct {
	Tick    int64
	Message string
}
func main() {
	encoded, err := json.Marshal(Stamp{Tick: 3, Message: "ready"})
	fmt.Println(string(encoded), err)
}`

var ch23Frames234Idea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "ترتيب إضافة للتوضيح."},
	{Line: 3, Caption: "رتّب الأسماء."},
	{Line: 4, Caption: "اعرض ترتيب ثابت.", Output: "bay/a, bay/z"},
}

var ch23Frames234Try = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 1, Caption: "الـstruct يحدد الترتيب."},
	{Line: 6, Caption: "حوّل القيمة لـJSON."},
	{Line: 7, Caption: "المفاتيح بنفس ترتيب الحقول.", Output: "{\"Tick\":3,\"Message\":\"ready\"} <nil>"},
}

var ch23Lesson234Arabic = []LessonStep{
	{ID: "c234.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch23Markdown("c234.mission", "بعد تخزين القطع، جهّز index وcommit قابلين لإعادة البناء.\n\nالـstarter معرّف فيه أنواع الفهرس والوقت؛ البيانات اللي تدخل لازم تفضل تحت سيطرة المنادي."),
	}},
	{ID: "c234.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch23Walkthrough("c234.tree-trace", "الـmap مفيدة للربط، لكن ترتيب iteration مش ثابت. حوّل الأسماء لقائمة مرتبة قبل بناء تمثيل حتمي؛ نفس المحتوى يدي نفس الشجرة.", ch23Example234Idea, "e234-stable-order", []string{"fmt", "slices", "strings"}, ch23Frames234Idea),
	}},
	{ID: "c234.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch23QuestionExample("c234.json-predict", "توقّع شكل JSON اللي بيرجعه Marshal.", ch23Example234Try, []string{`{"Message":"ready","Tick":3} <nil>`, `{"Tick":3,"Message":"ready"} <nil>`, "Stamp{3 ready} <nil>"}, 1, []string{
			"Marshal مش بيبدل ترتيب الحقول في struct.",
			"المفاتيح ماشية حسب ترتيب الحقول.",
			"Marshal بيرجع JSON، مش طباعة Go.",
		}),
		ch23Walkthrough("c234.json-reveal", "تابع ترتيب المفاتيح من تعريف النوع.", ch23Example234Try, "e234-json-order", []string{"encoding/json", "fmt"}, ch23Frames234Try),
		ch23Playground("ch23-stable-order", "ch23-stable-order", "بدّل ترتيب الإضافة.", ch23OrderSandboxConfig),
	}},
	{ID: "c234.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch23Markdown("c234.watch", "- حدّث الـindex بعد نجاح التخزين بس.\n- بيانات الوقت والرسالة جاية من ~SupplyStamp~."),
	}},
	{ID: "c234.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch23Markdown("c234.plan", "افصل staging عن إنشاء snapshot. افحص كل مرجع قبل الكتابة، وخلي ترتيب الشجرة مستقل عن ترتيب إضافة الشحنات."),
		ch23Question("c234.order-question", "ليه نرتب أسماء الـindex؟", []string{"لتثبيت نتيجة الشجرة", "لتغيير IDs القديمة", "لترتيب قراءة الملفات على الجهاز"}, 0, []string{
			"الترتيب الثابت يجعل نفس الفهرس ينتج نفس الشجرة.",
			"الترتيب لا يعيد كتابة القطع القديمة.",
			"الأسماء labels وليست أوامر قراءة للجهاز.",
		}),
	}},
	{ID: "c234.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch23Markdown("c234.turn", "ابدأ بـstaging، وبعدها اعمل commit بوقت ورسالة محقونين وراجع الأب."),
	}},
}

var ch23Lesson234English = []LessonStep{
	{ID: "c234.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch23Markdown("c234.mission", "After storing objects, build an index and a reproducible commit.\n\nThe starter defines the index and timestamp types; callers keep ownership of the input data."),
	}},
	{ID: "c234.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch23Walkthrough("c234.tree-trace", "A map is useful for lookup, but iteration order is unstable.\n\nSort names before building a deterministic representation; identical content then produces the same tree.", ch23Example234Idea, "e234-stable-order", []string{"fmt", "slices", "strings"}, ch23EnglishFrames(ch23Frames234Idea)),
	}},
	{ID: "c234.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch23QuestionExample("c234.json-predict", "Predict the JSON that Marshal returns.", ch23Example234Try, []string{`{"Message":"ready","Tick":3} <nil>`, `{"Tick":3,"Message":"ready"} <nil>`, "Stamp{3 ready} <nil>"}, 1, []string{
			"Marshal does not reorder struct fields.",
			"Keys follow the field order.",
			"Marshal returns JSON, not Go formatting.",
		}),
		ch23Walkthrough("c234.json-reveal", "Follow the key order from the type definition.", ch23Example234Try, "e234-json-order", []string{"encoding/json", "fmt"}, ch23EnglishFrames(ch23Frames234Try)),
		ch23Playground("ch23-stable-order", "ch23-stable-order", "Change insertion order.", ch23OrderSandboxConfig),
	}},
	{ID: "c234.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch23Markdown("c234.watch", "- Update the index only after storage succeeds.\n- Timestamp and message come from ~SupplyStamp~."),
	}},
	{ID: "c234.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch23Markdown("c234.plan", "Separate staging from snapshot creation.\n\nVerify every reference before writing, and make tree order independent of supply insertion order."),
		ch23Question("c234.order-question", "Why sort index names?", []string{"To make the tree deterministic", "To change old IDs", "To order host file reads"}, 0, []string{
			"Stable ordering makes the same index produce the same tree.",
			"Sorting does not rewrite old objects.",
			"Names are labels, not host read commands.",
		}),
	}},
	{ID: "c234.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch23Markdown("c234.turn", "Start with staging, then commit with an injected stamp and inspect the parent."),
	}},
}

const ch23Example23BossIdea = `func main() {
	sensors := []string{"probe-z", "probe-a"}
	slices.Sort(sensors)
	fmt.Println(sensors)
}`

const ch23Example23BossTry = `func main() {
	reader := csv.NewReader(strings.NewReader("station,reading\nsensor-2,7\n"))
	header, err := reader.Read()
	row, rowErr := reader.Read()
	fmt.Println(header[0], row[0], row[1], err, rowErr)
}`

var ch23Frames23BossIdea = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "قراءات للتوضيح."},
	{Line: 3, Caption: "ثبّت ترتيب العرض."},
	{Line: 4, Caption: "الترتيب بقى ثابت.", Output: "[probe-a probe-z]"},
}

var ch23Frames23BossTry = []LessonFrame{
	{Line: 0, Caption: "قبل التنفيذ."},
	{Line: 2, Caption: "قارئ CSV يحفظ شكل الحقول."},
	{Line: 3, Caption: "اقرأ العنوان كسجل."},
	{Line: 4, Caption: "اقرأ السطر اللي بعده."},
	{Line: 5, Caption: "CSV رجّع الحقول والأخطاء.", Output: "station sensor-2 7 <nil> <nil>"},
}

var ch23LessonBossArabic = []LessonStep{
	{ID: "c23boss.mission-step", Title: "المهمة", Blocks: []LessonBlock{
		ch23Markdown("c23boss.mission", "دليل العاصفة يجمع parsing وJSON وGCM وحفظ آمن. استخدم الأدوات اللي اتعلمتها، وخلي الملف القديم سليم لو أي مرحلة قبل النشر فشلت."),
	}},
	{ID: "c23boss.idea-step", Title: "الفكرة", Blocks: []LessonBlock{
		ch23Walkthrough("c23boss.order-trace", "المهمة تركيبة من مهارات سابقة: سجلات واضحة، ترتيب ثابت، بيانات وقت محقونة، وموارد ليها عمر.\n\nماتبدأش النشر وإنت لسه بتعالج المدخلات.", ch23Example23BossIdea, "e23boss-stable-order", []string{"fmt", "slices"}, ch23Frames23BossIdea),
	}},
	{ID: "c23boss.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch23QuestionExample("c23boss.csv-predict", "توقّع الحقول اللي بيرجعها قارئ CSV.", ch23Example23BossTry, []string{"station sensor-2 7 <nil> <nil>", "station,reading sensor-2,7", "sensor-2 7 EOF"}, 0, []string{
			"كل سجل بيتحول لقيم منفصلة.",
			"الفواصل مش جزء من الحقول الراجعة.",
			"فيه سطرين صالحين قبل نهاية المصدر.",
		}),
		ch23Walkthrough("c23boss.csv-reveal", "شوف كل قراءة وهي بترجع سجل مستقل.", ch23Example23BossTry, "e23boss-csv-records", []string{"encoding/csv", "fmt", "strings"}, ch23Frames23BossTry),
		ch23Playground("ch23-csv-records", "ch23-csv-records", "حوّل السطر لحقول.", ch23CSVSandboxConfig),
	}},
	{ID: "c23boss.watch-step", Title: "خد بالك", Blocks: []LessonBlock{
		ch23Markdown("c23boss.watch", "- ما تعملش trim للقيم قبل التحقق.\n- مفيش كتابة للهدف قبل اكتمال التحقق."),
	}},
	{ID: "c23boss.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch23Markdown("c23boss.plan", "قسّم الشغل لمراحل صغيرة وافحص نتيجة كل مرحلة. خلّي ترتيب التقرير ثابت، واستخدم القيم المحقونة بدل الساعة أو العشوائية العامة."),
		ch23Question("c23boss.time-question", "إيه اللي يثبت وقت التقرير؟", []string{"الوقت الممرر للدالة", "وقت تشغيل الجهاز", "آخر تعديل للملف"}, 0, []string{
			"القيمة المحقونة تخلي النتيجة قابلة للتكرار.",
			"وقت الجهاز ممكن يختلف بين التشغيلات.",
			"وقت الملف مش metadata التقرير المطلوبة.",
		}),
	}},
	{ID: "c23boss.turn-step", Title: "دورك", Blocks: []LessonBlock{
		ch23Markdown("c23boss.turn", "ابدأ بقراءة CSV، وبعدها راجع كل مرحلة قبل فتح الملف النهائي."),
	}},
}

var ch23LessonBossEnglish = []LessonStep{
	{ID: "c23boss.mission-step", Title: "Mission", Blocks: []LessonBlock{
		ch23Markdown("c23boss.mission", "Storm evidence combines parsing, JSON, GCM, and safe saving.\n\nUse tools from earlier chapters, and keep the old file intact if any pre-publication stage fails."),
	}},
	{ID: "c23boss.idea-step", Title: "Idea", Blocks: []LessonBlock{
		ch23Walkthrough("c23boss.order-trace", "This mission combines earlier skills: clear records, stable ordering, injected time, and resources with lifetimes.\n\nDo not publish while input is still being processed.", ch23Example23BossIdea, "e23boss-stable-order", []string{"fmt", "slices"}, ch23EnglishFrames(ch23Frames23BossIdea)),
	}},
	{ID: "c23boss.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch23QuestionExample("c23boss.csv-predict", "Predict the fields returned by the CSV reader.", ch23Example23BossTry, []string{"station sensor-2 7 <nil> <nil>", "station,reading sensor-2,7", "sensor-2 7 EOF"}, 0, []string{
			"Each record becomes separate values.",
			"Commas are not part of the returned fields.",
			"Two valid rows appear before the source ends.",
		}),
		ch23Walkthrough("c23boss.csv-reveal", "Watch each read return a separate record.", ch23Example23BossTry, "e23boss-csv-records", []string{"encoding/csv", "fmt", "strings"}, ch23EnglishFrames(ch23Frames23BossTry)),
		ch23Playground("ch23-csv-records", "ch23-csv-records", "Split one CSV record.", ch23CSVSandboxConfig),
	}},
	{ID: "c23boss.watch-step", Title: "Watch out", Blocks: []LessonBlock{
		ch23Markdown("c23boss.watch", "- Do not trim values before validation.\n- Do not write the target before validation completes."),
	}},
	{ID: "c23boss.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch23Markdown("c23boss.plan", "Split the work into small stages and check each result.\n\nKeep report order stable, and use injected values instead of the system clock or global randomness."),
		ch23Question("c23boss.time-question", "What keeps the report time stable?", []string{"The time passed to the function", "The machine clock", "The file's modification time"}, 0, []string{
			"An injected value makes the result repeatable.",
			"The machine clock can change between runs.",
			"File time is not the requested report metadata.",
		}),
	}},
	{ID: "c23boss.turn-step", Title: "Your turn", Blocks: []LessonBlock{
		ch23Markdown("c23boss.turn", "Start with CSV reading, then inspect every stage before opening the final file."),
	}},
}
