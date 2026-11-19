package levels

func ch30Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch30Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch30Caption(arabic, englishText string, useEnglish bool) string {
	if useEnglish {
		return englishText
	}
	return arabic
}

func ch30TextPair(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

func ch30SandboxBlock(id, kind string, config map[string]any, english bool) LessonBlock {
	intro := config["intro"].(map[string]string)
	text := intro["ar"]
	if english {
		text = intro["en"]
	}
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch30Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

const ch30CacheHitCode = `func main() {
	cache := map[string]string{"ridge": "ready"}
	value, ok := cache["ridge"]
	fmt.Printf("%q %t\n", value, ok)
}`

func ch30CacheHit(text string, english bool) LessonBlock {
	return ch30Walkthrough("c301.cache-hit", text, ch30CacheHitCode, "e301-cache-hit", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("الخريطة فيها قيمة.", "The map has a value.", english)},
		{Line: 3, Caption: ch30Caption("القراءة وجدت المفتاح.", "The lookup found the key.", english), Variables: []LessonVariable{{Name: "value", Value: "ready"}, {Name: "ok", Value: "true"}}},
		{Line: 4, Caption: ch30Caption("اطبع القيمة وحالة الوجود.", "Print the value and presence.", english), Variables: []LessonVariable{{Name: "value", Value: "ready"}, {Name: "ok", Value: "true"}}, Output: "\"ready\" true"},
	})
}

const ch30CacheMissCode = `func main() {
	cache := map[string]string{"ridge": "ready"}
	value, ok := cache["cove"]
	fmt.Printf("%q %t\n", value, ok)
}`

func ch30CacheMiss(text string, english bool) LessonBlock {
	return ch30Walkthrough("c301.cache-miss", text, ch30CacheMissCode, "e301-cache-miss", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("الخريطة فيها مفتاح تاني.", "The map has another key.", english)},
		{Line: 3, Caption: ch30Caption("المفتاح غير موجود.", "The key is missing.", english), Variables: []LessonVariable{{Name: "value", Value: ""}, {Name: "ok", Value: "false"}}},
		{Line: 4, Caption: ch30Caption("الصفرية مش hit.", "A zero value is not a hit.", english), Variables: []LessonVariable{{Name: "value", Value: ""}, {Name: "ok", Value: "false"}}, Output: "\"\" false"},
	})
}

var ch30Lesson301Arabic = []LessonStep{
	{ID: "c301.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch30Markdown("c301.mission", "1. خلي Redis يسرّع قراية المنارة.\n2. المخزن يفضل مصدر الحقيقة.\n3. احفظ فيه الأول، وبعدها حدّث الكاش.")}},
	{ID: "c301.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch30CacheHit("الكاش زي رف قريب: لو المفتاح موجود، القيمة جاهزة.\n\nلو غايب، المصدر هو الحقيقة؛ فشل ملء الكاش مش معناه فشل القراية.\n\n~ok~ بتميّز المفتاح الغايب عن قيمة فاضية. Redis بيفرق بين ~redis.Nil~ وعطل الاتصال.\n\nلو ~Del~ فشل، علامة dirty تخلي القراية تسيب المفتاح القديم لحد نجاح الحذف.\n\n~Mutex~ حوالين العملية يمنع قراية قديمة تملأ الكاش بعد كتابة.\n\nالـ TTL مدة صلاحية؛ اختبارات المستوى بتحرك ساعة Redis نفسها.", false)}},
	{ID: "c301.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch30Question("c301.predict-miss", "توقّع ناتج قراءة مفتاح غايب:\n\n~~~go\nfunc main() {\n\tcache := map[string]string{\"ridge\": \"ready\"}\n\tvalue, ok := cache[\"cove\"]\n\tfmt.Printf(\"%q %t\\n\", value, ok)\n}\n~~~", []string{"\"\" false", "\"ready\" true", "\"\" true"}, 0, []string{
			"قيمة غايبة بتكون فاضية ومعاها false.",
			"الخريطة ما فيهاش المفتاح المطلوب.",
			"ok بيقول إن المفتاح مش موجود.",
		}),
		ch30CacheMiss("تابع القيمة وراية الوجود لحد الطباعة.", false),
		ch30SandboxBlock("c301.playground", "ch30-cache-aside", ch30CacheConfig, false),
	}},
	{ID: "c301.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch30Markdown("c301.watch", "- ~redis.Nil~ غياب مؤكد؛ باقي أخطاء Redis ممكن تكون عطل.\n- نجاح الحفظ يسبق حذف المفتاح؛ فشل الحذف يفضل محتاج مزامنة.")}},
	{ID: "c301.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch30Markdown("c301.plan", "قسّم القرار: تحقّق من المدخل والسياق، اقرأ الكاش، وبعدها اختار المصدر أو النتيجة.\n\nفي الكتابة، اسأل نفسك: هل المصدر نجح قبل ما تلمس المفتاح؟ وإيه اللي يحصل لو الحذف فشل؟"),
		ch30Question("c301.cache-authority", "مين مرجع الحقيقة؟", []string{"المخزن", "Redis", "TTL"}, 0, []string{
			"المخزن هو المرجع.",
			"Redis كاش مؤقت.",
			"TTL للصلاحية.",
		}),
	}},
	{ID: "c301.turn-step", Title: "دورك", Blocks: []LessonBlock{ch30Markdown("c301.turn", "ابدأ بقراءة واحدة؛ راجع نتيجة كل طبقة قبل الانتقال للي بعدها.")}},
}

var ch30Lesson301English = []LessonStep{
	{ID: "c301.mission-step", Title: "Mission", Blocks: []LessonBlock{ch30Markdown("c301.mission", "1. Let Redis speed up beacon reads.\n2. Keep the store as the source of truth.\n3. Save there first, then refresh the cache.")}},
	{ID: "c301.idea-step", Title: "Idea", Blocks: []LessonBlock{ch30CacheHit("A cache is a nearby shelf; a present key is ready.\n\nOn a miss, the store is authoritative. Cache-fill failure does not fail that read.\n\n~ok~ separates missing from empty; ~redis.Nil~ differs from an outage.\n\nA dirty marker skips stale cache until deletion succeeds. A mutex prevents an old read refilling after a write.\n\nRedis owns expiry time; tests advance its clock.", true)}},
	{ID: "c301.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch30Question("c301.predict-miss", "Predict the lookup for a missing key:\n\n~~~go\nfunc main() {\n\tcache := map[string]string{\"ridge\": \"ready\"}\n\tvalue, ok := cache[\"cove\"]\n\tfmt.Printf(\"%q %t\\n\", value, ok)\n}\n~~~", []string{"\"\" false", "\"ready\" true", "\"\" true"}, 0, []string{
			"A missing value is empty and its flag is false.",
			"The map does not contain the requested key.",
			"The flag says the key was not found.",
		}),
		ch30CacheMiss("Follow the value and presence flag to the print.", true),
		ch30SandboxBlock("c301.playground", "ch30-cache-aside", ch30CacheConfig, true),
	}},
	{ID: "c301.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch30Markdown("c301.watch", "- ~redis.Nil~ means confirmed absence; other Redis errors may be outages.\n- Save before deleting the key; a failed deletion still needs synchronization.")}},
	{ID: "c301.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch30Markdown("c301.plan", "Split the decision: validate input and context, read the cache, then choose the source or result.\n\nFor writes, ask whether the source succeeded before touching the key. What remains if deletion fails?"),
		ch30Question("c301.cache-authority", "Who owns truth?", []string{"Store", "Redis", "TTL"}, 0, []string{
			"The store is authoritative.",
			"Redis is temporary.",
			"TTL controls expiry.",
		}),
	}},
	{ID: "c301.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch30Markdown("c301.turn", "Start with one read; check each layer's result before moving on.")}},
}

const ch30DocumentCode = `type Card struct {
	ID  int64
	Tag string
}
func main() {
	c := Card{ID: 4, Tag: "basalt"}
	fmt.Println(c.Tag, c.ID)
}`

func ch30Document(text string, english bool) LessonBlock {
	return ch30Walkthrough("c302.document", text, ch30DocumentCode, "e302-document", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 1, Caption: ch30Caption("Card تجمع حقول العينة.", "Card groups sample fields.", english)},
		{Line: 5, Caption: ch30Caption("القيمة تحمل tag وID.", "The value holds tag and ID.", english), Variables: []LessonVariable{{Name: "c", Value: "{4 basalt}"}}},
		{Line: 6, Caption: ch30Caption("اطبع الحقول المختارة.", "Print the selected fields.", english), Variables: []LessonVariable{{Name: "c", Value: "{4 basalt}"}}, Output: "basalt 4"},
	})
}

const ch30PageCode = `func main() {
	rows := []string{"mica", "slate", "quartz"}
	limit := 2
	fmt.Println(rows[:limit], len(rows) > limit)
}`

func ch30Page(text string, english bool) LessonBlock {
	return ch30Walkthrough("c302.page-extra", text, ch30PageCode, "e302-page-extra", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("القائمة فيها عنصر زيادة.", "The list has an extra row.", english)},
		{Line: 3, Caption: ch30Caption("الحد يختار أول عنصرين.", "The limit selects two rows.", english), Variables: []LessonVariable{{Name: "limit", Value: "2"}}},
		{Line: 4, Caption: ch30Caption("الزيادة تعني صفحة تانية.", "An extra row means another page.", english), Variables: []LessonVariable{{Name: "limit", Value: "2"}}, Output: "[mica slate] true"},
	})
}

var ch30Lesson302Arabic = []LessonStep{
	{ID: "c302.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch30Markdown("c302.mission", "1. خزّن عينات الصخور كمستندات.\n2. امنع تكرار العلامة من قاعدة البيانات.\n3. اعرض صفحات بترتيب واضح.")}},
	{ID: "c302.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch30Document("الـstruct شكل السجل جوه Go؛ وBSON tags بتربط الحقول بأسماء التخزين.\n\n~_id~ مميز افتراضياً. الفهرس الفريد هو الحكم وقت الطلبات المتزامنة، مش فحص سابق في الذاكرة.\n\nخلي ~Archive~ يتعامل مع واجهة ~SampleStore~؛ الـadapter وحده يعرف MongoDB.\n\nالصفحات تستخدم ID المؤشر، والـcursor مورد لازم يتقفل ويتفحص خطؤه.", false)}},
	{ID: "c302.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch30Question("c302.predict-page", "توقّع ناتج الصفحة الصغيرة:\n\n~~~go\nfunc main() {\n\trows := []string{\"mica\", \"slate\", \"quartz\"}\n\tlimit := 2\n\tfmt.Println(rows[:limit], len(rows) > limit)\n}\n~~~", []string{"[mica slate] true", "[mica slate quartz] false", "[quartz] true"}, 0, []string{
			"أول عنصرين ظاهرين ولسه فيه عنصر زيادة.",
			"الـslice المعروضة بتنتهي عند الحد.",
			"الزيادة مش أول عنصر في الصفحة.",
		}),
		ch30Page("شوف عدد الصفوف قبل وبعد الحد، وإشارة الزيادة.", false),
		ch30SandboxBlock("c302.playground", "ch30-document-pages", ch30DocumentConfig, false),
	}},
	{ID: "c302.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch30Markdown("c302.watch", "- الفحص قبل insert مش بديل عن unique index.\n- اقفل ~cursor~ حتى مع خطأ، وافحص ~Err~ بعد اللف.")}},
	{ID: "c302.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch30Markdown("c302.plan", "افصل طبقة المجال عن التخزين. راجع التحقق والتطبيع قبل النداء.\n\nبعدها افحص أخطاء Mongo، والغياب، والمؤشر. طلب ~limit+1~ يوضح وجود صفحة تالية من غير إرجاع العنصر الإضافي."),
		ch30Question("c302.unique-index", "إيه اللي يمنع تكرار tag وقت السباق؟", []string{"فهرس unique في القاعدة", "فحص map قبل الكتابة", "ترتيب النتائج"}, 0, []string{
			"قاعدة البيانات تفرض uniqueness بين الطلبات المتزامنة.",
			"فحص الذاكرة ممكن يسبق طلباً متزامناً.",
			"الترتيب لا يمنع إضافة مستند مكرر.",
		}),
	}},
	{ID: "c302.turn-step", Title: "دورك", Blocks: []LessonBlock{ch30Markdown("c302.turn", "ابدأ بواجهة التخزين؛ بعد نجاحها أضف adapter Mongo خطوة خطوة.")}},
}

var ch30Lesson302English = []LessonStep{
	{ID: "c302.mission-step", Title: "Mission", Blocks: []LessonBlock{ch30Markdown("c302.mission", "1. Store rock samples as documents.\n2. Let the database prevent duplicate tags.\n3. Return pages in a clear order.")}},
	{ID: "c302.idea-step", Title: "Idea", Blocks: []LessonBlock{ch30Document("A struct describes a record in Go; BSON tags map its fields to stored names.\n\n~_id~ is unique by default. The unique index decides under concurrent requests, not an earlier in-memory check.\n\nKeep ~Archive~ behind the ~SampleStore~ interface; only the adapter knows MongoDB.\n\nPages use an ID cursor, and each database cursor must be closed and checked for errors.", true)}},
	{ID: "c302.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch30Question("c302.predict-page", "Predict this small page:\n\n~~~go\nfunc main() {\n\trows := []string{\"mica\", \"slate\", \"quartz\"}\n\tlimit := 2\n\tfmt.Println(rows[:limit], len(rows) > limit)\n}\n~~~", []string{"[mica slate] true", "[mica slate quartz] false", "[quartz] true"}, 0, []string{
			"The first two are shown, with one extra row present.",
			"The displayed slice stops at the limit.",
			"The extra row is not the page's first item.",
		}),
		ch30Page("Compare rows before and after the limit, plus the extra flag.", true),
		ch30SandboxBlock("c302.playground", "ch30-document-pages", ch30DocumentConfig, true),
	}},
	{ID: "c302.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch30Markdown("c302.watch", "- A pre-insert check cannot replace a unique index.\n- Close the ~cursor~ on errors, and check ~Err~ after iteration.")}},
	{ID: "c302.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch30Markdown("c302.plan", "Separate the domain layer from storage. Validate and normalize before calling it.\n\nThen classify Mongo errors, missing records, and cursor failures. Fetching ~limit+1~ reveals another page without returning the extra row."),
		ch30Question("c302.unique-index", "What prevents duplicate tags during a race?", []string{"A database unique index", "A map check before writing", "Sorting the results"}, 0, []string{
			"The database enforces uniqueness across concurrent requests.",
			"An in-memory check can race with another request.",
			"Sorting does not prevent a duplicate insert.",
		}),
	}},
	{ID: "c302.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch30Markdown("c302.turn", "Start with the storage interface; add the Mongo adapter one step at a time.")}},
}

const ch30IndexCode = `func main() {
	postings := map[string][]string{
		"amber": {"note-q", "note-m"},
	}
	fmt.Println(postings["amber"])
}`

func ch30Index(text string, english bool) LessonBlock {
	return ch30Walkthrough("c303.postings", text, ch30IndexCode, "e303-postings", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("الكلمة مفتاح الفهرس.", "The term is the index key.", english)},
		{Line: 3, Caption: ch30Caption("القيمة قائمة النشرات.", "The value is a posting list.", english), Variables: []LessonVariable{{Name: "postings", Value: "map[amber:[note-q note-m]]"}}},
		{Line: 5, Caption: ch30Caption("الوصول مباشر للكلمة.", "Lookup goes straight to the term.", english), Variables: []LessonVariable{{Name: "postings", Value: "map[amber:[note-q note-m]]"}}, Output: "[note-q note-m]"},
	})
}

const ch30SortCode = `func main() {
	codes := []string{"note-q", "note-m"}
	sort.Strings(codes)
	fmt.Println(codes)
}`

func ch30Sort(text string, english bool) LessonBlock {
	return ch30Walkthrough("c303.lexical-order", text, ch30SortCode, "e303-lexical-order", []string{"fmt", "sort"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("القيم لسه بترتيب الإدخال.", "Values remain in input order.", english), Variables: []LessonVariable{{Name: "codes", Value: "[note-q note-m]"}}},
		{Line: 3, Caption: ch30Caption("sort ترتب نصياً.", "sort orders strings.", english), Variables: []LessonVariable{{Name: "codes", Value: "[note-m note-q]"}}},
		{Line: 4, Caption: ch30Caption("اطبع الترتيب الثابت.", "Print the stable order.", english), Variables: []LessonVariable{{Name: "codes", Value: "[note-m note-q]"}}, Output: "[note-m note-q]"},
	})
}

var ch30Lesson303Arabic = []LessonStep{
	{ID: "c303.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch30Markdown("c303.mission", "1. ابحث في نشرات المستعمرة.\n2. رتّب النتائج والصفحات بثبات.\n3. خلّي القراءة من الفهرس مش من كل النصوص.")}},
	{ID: "c303.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch30Index("الفهرس المقلوب زي دليل: كل كلمة بتوصل لقائمة المستندات اللي ظهرت فيها.\n\nالتطبيع لازم يكون واحد للنشرة والسؤال. ~unicode.IsLetter~ و~unicode.IsDigit~ يغطّوا أكتر من ASCII.\n\nالـmap مش بتضمن ترتيب. بعد جمع النتائج، قاعدة ترتيب صريحة تخلي الصفحات قابلة للتكرار.\n\nانسخ بيانات الإدخال؛ الفهرس الثابت ينفع لقراءات متزامنة.", false)}},
	{ID: "c303.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch30Question("c303.predict-sort", "بعد الترتيب النصي، إيه الناتج؟\n\n~~~go\nfunc main() {\n\tcodes := []string{\"note-q\", \"note-m\"}\n\tsort.Strings(codes)\n\tfmt.Println(codes)\n}\n~~~", []string{"[note-m note-q]", "[note-q note-m]", "[note-m]"}, 0, []string{
			"المقارنة النصية تبدأ بالحروف الأولى.",
			"ده ترتيب الإدخال قبل sort.",
			"العنصرين موجودين بعد الترتيب.",
		}),
		ch30Sort("تابع slice قبل sort وبعدها، ثم راقب الطباعة.", false),
		ch30SandboxBlock("c303.playground", "ch30-search-index", ch30IndexConfig, false),
	}},
	{ID: "c303.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch30Markdown("c303.watch", "- اجمع كلمات السؤال المختلفة مرة واحدة؛ ترتيب map غير مضمون.\n- رتّب كل النتائج قبل تطبيق ~offset~ و~limit~.")}},
	{ID: "c303.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch30Markdown("c303.plan", "قسّم الشغل: طبّع، اجمع postings، احسب الدرجات، رتّب، ثم قص الصفحة.\n\nراجع إن النتيجة قابلة للتعديل من غير ما تغيّر فهرساً أو قراءة لاحقة."),
		ch30Question("c303.stable-tie", "إزاي تثبّت ترتيب نتيجة متعادلة؟", []string{"مفتاح ترتيب ثانوي واضح", "الاعتماد على ترتيب map", "الاحتفاظ بترتيب goroutine"}, 0, []string{
			"مقارنة ثانية تجعل التعادل حتمياً.",
			"ترتيب map غير مضمون بين التشغيلات.",
			"وصول goroutine مش قاعدة ترتيب للمخرجات.",
		}),
	}},
	{ID: "c303.turn-step", Title: "دورك", Blocks: []LessonBlock{ch30Markdown("c303.turn", "ابنِ طبقة الفهرس أولاً، وبعدها افحص صفحة صغيرة من النتائج.")}},
}

var ch30Lesson303English = []LessonStep{
	{ID: "c303.mission-step", Title: "Mission", Blocks: []LessonBlock{ch30Markdown("c303.mission", "1. Search the colony dispatches.\n2. Rank results and pages consistently.\n3. Read the index instead of rescanning every text.")}},
	{ID: "c303.idea-step", Title: "Idea", Blocks: []LessonBlock{ch30Index("An inverted index is a directory: each term points to documents where it appears.\n\nUse the same normalization for documents and queries. ~unicode.IsLetter~ and ~unicode.IsDigit~ cover more than ASCII.\n\nA map does not promise order. An explicit ranking rule makes pages repeatable.\n\nCopy input data; an immutable index supports concurrent reads.", true)}},
	{ID: "c303.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch30Question("c303.predict-sort", "What prints after lexical sorting?\n\n~~~go\nfunc main() {\n\tcodes := []string{\"note-q\", \"note-m\"}\n\tsort.Strings(codes)\n\tfmt.Println(codes)\n}\n~~~", []string{"[note-m note-q]", "[note-q note-m]", "[note-m]"}, 0, []string{
			"String comparison begins with the first letters.",
			"That is the input order before sorting.",
			"Both values remain after sorting.",
		}),
		ch30Sort("Follow the slice before and after sort, then watch the print.", true),
		ch30SandboxBlock("c303.playground", "ch30-search-index", ch30IndexConfig, true),
	}},
	{ID: "c303.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch30Markdown("c303.watch", "- Count distinct query terms once; map order is unspecified.\n- Sort all hits before applying ~offset~ and ~limit~.")}},
	{ID: "c303.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch30Markdown("c303.plan", "Separate the work: normalize, gather postings, score, sort, then slice the page.\n\nCheck that callers can edit a result without changing the index or a later read."),
		ch30Question("c303.stable-tie", "How do you order tied results consistently?", []string{"An explicit secondary key", "Map iteration order", "Goroutine arrival order"}, 0, []string{
			"A second comparison makes ties deterministic.",
			"Map order is unspecified across runs.",
			"Goroutine arrival is not an output-order rule.",
		}),
	}},
	{ID: "c303.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch30Markdown("c303.turn", "Build the index layer first, then inspect a small result page.")}},
}

const ch30StickyCode = `func main() {
	start := time.Unix(120, 0)
	until := start.Add(time.Minute)
	fmt.Println(start.Before(until))
}`

func ch30Sticky(text string, english bool) LessonBlock {
	return ch30Walkthrough("c304.sticky-window", text, ch30StickyCode, "e304-sticky-window", []string{"fmt", "time"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("الوقت المحقون ثابت.", "The injected time is fixed.", english), Variables: []LessonVariable{{Name: "start", Value: "1970-01-01 00:02:00 +0000 UTC"}}},
		{Line: 3, Caption: ch30Caption("النهاية بعد البداية.", "The endpoint follows the start.", english), Variables: []LessonVariable{{Name: "until", Value: "1970-01-01 00:03:00 +0000 UTC"}}},
		{Line: 4, Caption: ch30Caption("الجلسة داخل النافذة.", "The session is inside the window.", english), Variables: []LessonVariable{{Name: "start", Value: "1970-01-01 00:02:00 +0000 UTC"}, {Name: "until", Value: "1970-01-01 00:03:00 +0000 UTC"}}, Output: "true"},
	})
}

const ch30WrappedErrorCode = `func main() {
	var ErrOffline = errors.New("offline")
	err := fmt.Errorf("replica: %w", ErrOffline)
	fmt.Println(errors.Is(err, ErrOffline))
}`

func ch30WrappedError(text string, english bool) LessonBlock {
	return ch30Walkthrough("c304.wrapped-cause", text, ch30WrappedErrorCode, "e304-wrapped-cause", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("أنشئ سبب العطل.", "Create the outage cause.", english)},
		{Line: 3, Caption: ch30Caption("التغليف يحفظ السبب.", "Wrapping preserves the cause.", english), Variables: []LessonVariable{{Name: "err", Value: "replica: offline"}}},
		{Line: 4, Caption: ch30Caption("Is تتعرف عليه.", "Is recognizes it.", english), Variables: []LessonVariable{{Name: "err", Value: "replica: offline"}}, Output: "true"},
	})
}

var ch30Lesson304Arabic = []LessonStep{
	{ID: "c304.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch30Markdown("c304.mission", "1. وجّه الكتابة للأصل مرة واحدة.\n2. خلّي القراءة العادية من النسخة.\n3. بعد الكتابة الناجحة، راعِ نافذة الجلسة.")}},
	{ID: "c304.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch30Sticky("الـprimary هو مصدر الكتابة؛ والـreplica ممكن تتأخر.\n\nبعد كتابة ناجحة، الجلسة نفسها تقرأ من الأصل وقت قصير. الجلسات التانية تفضل على مسارها.\n\nالساعة المحقونة تخلي الحدث قابل للاختبار من غير نوم أو ~time.Now~.\n\n~Mutex~ حوالين العملية يحافظ على ترتيب القرار والنداء.\n\nالنافذة معناها read-your-writes مؤقتاً، مش إن النسخة لحقت.", false)}},
	{ID: "c304.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch30Question("c304.predict-is", "هل ~errors.Is~ تلاقي السبب الملفوف؟\n\n~~~go\nfunc main() {\n\tvar ErrOffline = errors.New(\"offline\")\n\terr := fmt.Errorf(\"replica: %w\", ErrOffline)\n\tfmt.Println(errors.Is(err, ErrOffline))\n}\n~~~", []string{"true", "false", "offline"}, 0, []string{
			"الـwrapping حافظ على السبب الأصلي.",
			"الفحص بيمر عبر رسالة الخطأ الملفوفة.",
			"Is بترجع true أو false مش نص الرسالة.",
		}),
		ch30WrappedError("امشِ من السبب، للتغليف، وبعدها لفحص errors.Is.", false),
		ch30SandboxBlock("c304.playground", "ch30-read-write-window", ch30ReadWriteConfig, false),
	}},
	{ID: "c304.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch30Markdown("c304.watch", "- خطأ مجال من النسخة مش معناه إنها واقعة؛ fallback له سبب محدد.\n- القراية ما تجددش النافذة، وفشل الكتابة ما يغيّرش علامتها.")}},
	{ID: "c304.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch30Markdown("c304.plan", "راجع المسارات منفصلة: الكتابة، قراءة داخل النافذة، وقراءة عادية.\n\nاسأل مين يناديه كل مسار، وإمتى الإلغاء أو خطأ المجال يمنع fallback. اختبر جلسة أخرى لوحدها."),
		ch30Question("c304.replica-write", "الكتابة الناجحة تروح لمين؟", []string{"للـprimary فقط", "للـreplica فقط", "للاتنين بالتوازي"}, 0, []string{
			"الأصل يستقبل الكتابة مرة واحدة.",
			"النسخة مخصصة لمسار القراءة هنا.",
			"الراوتر لا يكرر الكتابة بين المخزنين.",
		}),
	}},
	{ID: "c304.turn-step", Title: "دورك", Blocks: []LessonBlock{ch30Markdown("c304.turn", "ابدأ بتحديد مسار كل نداء؛ بعده راجع العلامة والساعة.")}},
}

var ch30Lesson304English = []LessonStep{
	{ID: "c304.mission-step", Title: "Mission", Blocks: []LessonBlock{ch30Markdown("c304.mission", "1. Route writes to primary exactly once.\n2. Keep ordinary reads on the replica.\n3. Honor the session window after a successful write.")}},
	{ID: "c304.idea-step", Title: "Idea", Blocks: []LessonBlock{ch30Sticky("The primary is authoritative for writes; the replica may lag.\n\nAfter a successful write, that session reads primary for a short window. Other sessions keep their usual route.\n\nAn injected clock makes the boundary testable without sleeping or calling ~time.Now~.\n\nA mutex keeps route decisions and backend calls in order.\n\nThe window gives temporary read-your-writes, not proof that the replica caught up.", true)}},
	{ID: "c304.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch30Question("c304.predict-is", "Does ~errors.Is~ find a wrapped cause?\n\n~~~go\nfunc main() {\n\tvar ErrOffline = errors.New(\"offline\")\n\terr := fmt.Errorf(\"replica: %w\", ErrOffline)\n\tfmt.Println(errors.Is(err, ErrOffline))\n}\n~~~", []string{"true", "false", "offline"}, 0, []string{
			"Wrapping preserves the original cause.",
			"The check walks through the wrapped error.",
			"Is returns true or false, not the message.",
		}),
		ch30WrappedError("Follow the cause, its wrapper, then the errors.Is check.", true),
		ch30SandboxBlock("c304.playground", "ch30-read-write-window", ch30ReadWriteConfig, true),
	}},
	{ID: "c304.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch30Markdown("c304.watch", "- A replica domain error does not prove an outage; fallback has a specific cause.\n- Reads do not renew the window, and failed writes do not change it.")}},
	{ID: "c304.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch30Markdown("c304.plan", "Review the routes separately: writes, reads inside the window, and ordinary reads.\n\nAsk which store each route calls, and when cancellation or a domain error blocks fallback. Check another session independently."),
		ch30Question("c304.replica-write", "Where does a successful write go?", []string{"Primary only", "Replica only", "Both in parallel"}, 0, []string{
			"The primary receives the write exactly once.",
			"The replica serves the read path here.",
			"The router does not duplicate writes across stores.",
		}),
	}},
	{ID: "c304.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch30Markdown("c304.turn", "Mark the route for each call first; then inspect the marker and clock.")}},
}

const ch30DraftCode = `type Sketch struct {
	Count int
}
func main() {
	saved := Sketch{Count: 7}
	proposal := Sketch{Count: saved.Count + 2}
	fmt.Println(saved.Count, proposal.Count)
}`

func ch30Draft(text string, english bool) LessonBlock {
	return ch30Walkthrough("c30b.draft-state", text, ch30DraftCode, "e30b-draft-state", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 1, Caption: ch30Caption("Sketch تمثل قيمة.", "Sketch represents one value.", english)},
		{Line: 5, Caption: ch30Caption("القيمة المحفوظة ثابتة.", "The saved value stays fixed.", english), Variables: []LessonVariable{{Name: "saved", Value: "{7}"}}},
		{Line: 6, Caption: ch30Caption("الاقتراح قيمة منفصلة.", "The proposal is separate.", english), Variables: []LessonVariable{{Name: "saved", Value: "{7}"}, {Name: "proposal", Value: "{9}"}}},
		{Line: 7, Caption: ch30Caption("اطبع القديم والمقترح.", "Print saved and proposed values.", english), Variables: []LessonVariable{{Name: "saved", Value: "{7}"}, {Name: "proposal", Value: "{9}"}}, Output: "7 9"},
	})
}

const ch30JoinedErrorsCode = `func main() {
	sentinel := errors.New("sync warning")
	cause := errors.New("cache unavailable")
	err := errors.Join(sentinel, cause)
	fmt.Println(errors.Is(err, sentinel), errors.Is(err, cause))
}`

func ch30JoinedErrors(text string, english bool) LessonBlock {
	return ch30Walkthrough("c30b.joined-errors", text, ch30JoinedErrorsCode, "e30b-joined-errors", []string{"errors", "fmt"}, []LessonFrame{
		{Line: 0, Caption: ch30Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch30Caption("علامة التحذير منفصلة.", "The warning marker is separate.", english)},
		{Line: 3, Caption: ch30Caption("سبب الكاش منفصل.", "The cache cause is separate.", english)},
		{Line: 4, Caption: ch30Caption("Join يحفظ السببين.", "Join keeps both causes.", english), Variables: []LessonVariable{{Name: "err", Value: "sync warning\ncache unavailable"}}},
		{Line: 5, Caption: ch30Caption("افحص كل سبب.", "Check both causes.", english), Variables: []LessonVariable{{Name: "err", Value: "sync warning\ncache unavailable"}}, Output: "true true"},
	})
}

var ch30Lesson30bArabic = []LessonStep{
	{ID: "c30b.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch30Markdown("c30b.mission", "1. احجز مؤن المشتل داخل transaction.\n2. خليك قادر ترجع كل الخصومات معاً.\n3. حدّث الكاش بعد commit بس.")}},
	{ID: "c30b.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch30Draft("الـtransaction تجمع تغييرات القاعدة: يا كلها تتثبت، يا كلها تتراجع.\n\nالـcommit حد فاصل. فشل إبطال الكاش بعده تحذير عن المزامنة، مش سبب لإلغاء الحجز.\n\nالنجمة هنا بتطلب الحفاظ على علامة التحذير وسبب الكاش معاً.\n\nاستخدم نفس ~tx~ لكل أوامر القاعدة، وسيب ~main~ فاضية.", false)}},
	{ID: "c30b.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch30Question("c30b.predict-causes", "بعد ~errors.Join~، هل السببان معروفان؟\n\n~~~go\nfunc main() {\n\tsentinel := errors.New(\"sync warning\")\n\tcause := errors.New(\"cache unavailable\")\n\terr := errors.Join(sentinel, cause)\n\tfmt.Println(errors.Is(err, sentinel), errors.Is(err, cause))\n}\n~~~", []string{"true true", "true false", "false false"}, 0, []string{
			"الخطأ المركب يحتفظ بالعلامة وبالسبب.",
			"Is تقدر تلاقي كل سبب داخل Join.",
			"الاثنان جزء من الخطأ المركب.",
		}),
		ch30JoinedErrors("تابع تكوين الخطأ المركب، ثم افحص كل سبب.", false),
		ch30SandboxBlock("c30b.playground", "ch30-commit-cache-boundary", ch30BossConfig, false),
	}},
	{ID: "c30b.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch30Markdown("c30b.watch", "- تحقّق من الطلب والسطور قبل أي كتابة.\n- فشل cache بعد commit ما يرجّعش خصم القاعدة.")}},
	{ID: "c30b.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch30Markdown("c30b.plan", "رتّب المراحل: تحقق، ابدأ transaction، نفّذ تغييراتها، ثم commit.\n\nبعدها فقط فكّر في إبطال الكاش. راجع إن خطأ أي خصم يرجع التغييرات السابقة."),
		ch30Question("c30b.cache-timing", "إمتى يبدأ إبطال الكاش؟", []string{"بعد commit", "قبل التحقق", "داخل كل خصم"}, 0, []string{
			"القاعدة تثبت الحجز قبل لمس الكاش.",
			"المدخلات تتراجع قبل الآثار الجانبية.",
			"إبطال الكاش خطوة بعد transaction كاملة.",
		}),
	}},
	{ID: "c30b.turn-step", Title: "دورك", Blocks: []LessonBlock{ch30Markdown("c30b.turn", "ارسم حد transaction الأول، وبعده أضف أثر الكاش.")}},
}

var ch30Lesson30bEnglish = []LessonStep{
	{ID: "c30b.mission-step", Title: "Mission", Blocks: []LessonBlock{ch30Markdown("c30b.mission", "1. Reserve nursery supplies in a transaction.\n2. Keep all debits reversible together.\n3. Update the cache only after commit.")}},
	{ID: "c30b.idea-step", Title: "Idea", Blocks: []LessonBlock{ch30Draft("A transaction groups database changes: they commit together or roll back together.\n\nCommit is the boundary. Cache invalidation failing afterward warns about synchronization; it does not undo the reservation.\n\nThe bonus asks you to preserve both the warning marker and cache cause.\n\nUse the same ~tx~ for database commands, and leave ~main~ empty.", true)}},
	{ID: "c30b.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch30Question("c30b.predict-causes", "After ~errors.Join~, can both causes be found?\n\n~~~go\nfunc main() {\n\tsentinel := errors.New(\"sync warning\")\n\tcause := errors.New(\"cache unavailable\")\n\terr := errors.Join(sentinel, cause)\n\tfmt.Println(errors.Is(err, sentinel), errors.Is(err, cause))\n}\n~~~", []string{"true true", "true false", "false false"}, 0, []string{
			"The joined error retains its marker and cause.",
			"Is can find each cause inside Join.",
			"Both errors belong to the joined value.",
		}),
		ch30JoinedErrors("Follow the joined error, then check each cause.", true),
		ch30SandboxBlock("c30b.playground", "ch30-commit-cache-boundary", ch30BossConfig, true),
	}},
	{ID: "c30b.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch30Markdown("c30b.watch", "- Validate the request and lines before writing.\n- Cache failure after commit does not undo database debits.")}},
	{ID: "c30b.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch30Markdown("c30b.plan", "Order the phases: validate, begin a transaction, apply its changes, then commit.\n\nOnly afterward consider cache invalidation. Check that any debit failure rolls back earlier changes."),
		ch30Question("c30b.cache-timing", "When does cache invalidation begin?", []string{"After commit", "Before validation", "Inside every debit"}, 0, []string{
			"The reservation commits before cache work begins.",
			"Validate inputs before side effects.",
			"Cache invalidation follows the complete transaction.",
		}),
	}},
	{ID: "c30b.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch30Markdown("c30b.turn", "Draw the transaction boundary first, then add the cache effect.")}},
}

var ch30CacheConfig = map[string]any{
	"intro": ch30TextPair("اقرأ مرة مع cache miss، وبعدها جرّب hit.", "Read through a cache miss, then try a hit."),
	"labels": map[string]any{
		"key": ch30TextPair("المفتاح", "Key"), "cache": ch30TextPair("الكاش", "Cache"), "store": ch30TextPair("المخزن", "Store"),
		"result": ch30TextPair("مصدر القراءة", "Read source"), "read": ch30TextPair("اقرأ", "Read"), "evict": ch30TextPair("امسح من الكاش", "Evict cache"),
		"miss": ch30TextPair("مش موجود", "Miss"), "hit": ch30TextPair("موجود", "Hit"), "cacheSource": ch30TextPair("الكاش", "Cache"),
		"storeSource": ch30TextPair("المخزن", "Store"), "filled": ch30TextPair("الكاش اتملأ", "Cache filled"),
	},
	"key": ch30TextPair("ridge", "ridge"), "value": ch30TextPair("جاهز", "ready"), "cachePresent": false,
}

var ch30DocumentConfig = map[string]any{
	"intro": ch30TextPair("تنقّل بين صفحات صغيرة وشوف وجود صفوف زيادة.", "Move through small pages and watch for extra rows."),
	"labels": map[string]any{
		"page": ch30TextPair("الصفحة", "Page"), "documents": ch30TextPair("المستندات", "Documents"),
		"more": ch30TextPair("فيه صفحة تانية", "More rows"), "next": ch30TextPair("الصفحة التالية", "Next page"),
		"reset": ch30TextPair("ابدأ من الأول", "Reset"), "yes": ch30TextPair("آه", "Yes"), "no": ch30TextPair("لأ", "No"),
	},
	"pageSize": 2,
	"rows":     []any{ch30TextPair("mica · 4", "mica · 4"), ch30TextPair("slate · 7", "slate · 7"), ch30TextPair("quartz · 9", "quartz · 9"), ch30TextPair("basalt · 12", "basalt · 12"), ch30TextPair("flint · 15", "flint · 15")},
}

var ch30IndexConfig = map[string]any{
	"intro": ch30TextPair("اختار term وشوف قائمة النشرات المرتبطة.", "Pick a term and inspect its posting list."),
	"labels": map[string]any{
		"term": ch30TextPair("الكلمة", "Term"), "postings": ch30TextPair("قائمة النشرات", "Posting list"),
		"lookup": ch30TextPair("جرّب كلمة تانية", "Look up another term"), "mapNote": ch30TextPair("ترتيب map مش مضمون.", "Map order is unspecified."),
	},
	"terms": []any{
		map[string]any{"term": ch30TextPair("amber", "amber"), "postings": []any{ch30TextPair("note-q", "note-q"), ch30TextPair("note-m", "note-m")}},
		map[string]any{"term": ch30TextPair("slate", "slate"), "postings": []any{ch30TextPair("note-r", "note-r"), ch30TextPair("note-b", "note-b"), ch30TextPair("note-x", "note-x")}},
	},
}

var ch30ReadWriteConfig = map[string]any{
	"intro": ch30TextPair("اكتب لجلسة، وقدّم الساعة وشوف مسار القراءة.", "Write for one session, advance time, and watch its read route."),
	"labels": map[string]any{
		"clock": ch30TextPair("الوقت", "Time"), "ownSession": ch30TextPair("الجلسة دي", "This session"),
		"otherSession": ch30TextPair("جلسة تانية", "Another session"), "route": ch30TextPair("مسار القراءة", "Read route"),
		"primary": ch30TextPair("Primary", "Primary"), "replica": ch30TextPair("Replica", "Replica"),
		"write": ch30TextPair("كتابة ناجحة", "Successful write"), "advance": ch30TextPair("قدّم الوقت", "Advance time"),
		"reset": ch30TextPair("ابدأ من الأول", "Reset"), "tick": ch30TextPair("دقيقة", "minute"),
	},
	"now": 0, "stickyTicks": 2,
}

var ch30BossConfig = map[string]any{
	"intro": ch30TextPair("احفظ الاقتراح، وبعدها جرّب فشل الكاش.", "Commit the proposal, then simulate a cache failure."),
	"labels": map[string]any{
		"draft": ch30TextPair("اقتراح الخصم", "Proposed debit"), "saved": ch30TextPair("المحفوظ في القاعدة", "Saved in database"),
		"cache": ch30TextPair("الكاش", "Cache"), "notCommitted": ch30TextPair("لسه ما اتحفظش", "Not committed"),
		"committed": ch30TextPair("اتحفظ", "Committed"), "fresh": ch30TextPair("متزامن", "In sync"),
		"warning": ch30TextPair("تحذير مزامنة", "Sync warning"), "commit": ch30TextPair("اعمل commit", "Commit transaction"),
		"failCache": ch30TextPair("حاكي فشل الكاش", "Simulate cache failure"), "reset": ch30TextPair("ابدأ من الأول", "Reset"),
		"savedBoundary": ch30TextPair("القاعدة اتثبتت قبل الكاش.", "The database committed before cache work."),
	},
	"before": 7, "proposed": 9,
}
