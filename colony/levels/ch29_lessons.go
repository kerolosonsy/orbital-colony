package levels

type ch29LessonCopy struct {
	Mission       string
	Idea          string
	IdeaCode      string
	IdeaImports   []string
	IdeaID        string
	IdeaStates    []LessonFrame
	IdeaCaptions  []string
	TryPrompt     string
	TryCode       string
	TryImports    []string
	TryID         string
	TryAnswer     int
	TryOptions    []string
	TryFeedback   []string
	TryStates     []LessonFrame
	TryCaptions   []string
	Reveal        string
	Watch         string
	Plan          string
	CheckQuestion string
	CheckOptions  []string
	CheckAnswer   int
	CheckFeedback []string
	Turn          string
	Hints         []string
}

func ch29Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch29Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{ID: id, Kind: "question", Text: text, Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback}}
}

func ch29Walkthrough(id, text, code string, imports []string, exampleID string, frames []LessonFrame) LessonBlock {
	return LessonBlock{ID: id, Kind: "walkthrough", Text: text, Code: code, Imports: imports, ExampleID: exampleID, Walkthrough: &LessonWalkthrough{Frames: frames}}
}

func ch29Frames(states []LessonFrame, captions []string) []LessonFrame {
	frames := make([]LessonFrame, len(states))
	for i, frame := range states {
		frame.Caption = captions[i]
		frames[i] = frame
	}
	return frames
}

func ch29BuildSteps(prefix string, english bool, lesson ch29LessonCopy) []LessonStep {
	titles := []string{"المهمة", "الفكرة", "جرّب", "خد بالك", "فكّر واختبر", "دورك"}
	tryPrompt := "توقّع الناتج:"
	if english {
		titles = []string{"Mission", "Idea", "Try it", "Watch out", "Think & check", "Your turn"}
		tryPrompt = "Predict the output:"
	}
	if lesson.TryPrompt == "" {
		lesson.TryPrompt = tryPrompt
	}
	questionText := lesson.TryPrompt + "\n\n~~~go\n" + lesson.TryCode + "\n~~~"
	return []LessonStep{
		{ID: prefix + ".mission-step", Title: titles[0], Blocks: []LessonBlock{ch29Markdown(prefix+".mission", lesson.Mission)}},
		{ID: prefix + ".idea-step", Title: titles[1], Blocks: []LessonBlock{ch29Walkthrough(prefix+".idea-trace", lesson.Idea, lesson.IdeaCode, lesson.IdeaImports, lesson.IdeaID, ch29Frames(lesson.IdeaStates, lesson.IdeaCaptions))}},
		{ID: prefix + ".try-step", Title: titles[2], Blocks: []LessonBlock{
			ch29Question(prefix+".try-predict", questionText, lesson.TryOptions, lesson.TryAnswer, lesson.TryFeedback),
			ch29Walkthrough(prefix+".try-reveal", lesson.Reveal, lesson.TryCode, lesson.TryImports, lesson.TryID, ch29Frames(lesson.TryStates, lesson.TryCaptions)),
			ch29SandboxBlock(prefix+".playground", prefix, english),
		}},
		{ID: prefix + ".watch-step", Title: titles[3], Blocks: []LessonBlock{ch29Markdown(prefix+".watch", lesson.Watch)}},
		{ID: prefix + ".check-step", Title: titles[4], Blocks: []LessonBlock{
			ch29Markdown(prefix+".plan", lesson.Plan),
			ch29Question(prefix+".check", lesson.CheckQuestion, lesson.CheckOptions, lesson.CheckAnswer, lesson.CheckFeedback),
		}},
		{ID: prefix + ".turn-step", Title: titles[5], Blocks: []LessonBlock{ch29Markdown(prefix+".turn", lesson.Turn)}},
	}
}

func ch29SandboxBlock(id, prefix string, english bool) LessonBlock {
	kind, config := ch29SandboxFor(prefix)
	intro := config["intro"].(map[string]string)
	text := intro["ar"]
	if english {
		text = intro["en"]
	}
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch29TextPair(ar, en string) map[string]string {
	return map[string]string{"ar": ar, "en": en}
}

func ch29SandboxFor(prefix string) (string, map[string]any) {
	switch prefix {
	case "c291":
		return "ch29-sql-binding", map[string]any{
			"intro": ch29TextPair("بدّل الـkey وشوف نتيجة الربط.", "Change the key and watch binding."),
			"query": "SELECT payload FROM beacon WHERE code = ?",
			"labels": map[string]any{
				"query": ch29TextPair("الـquery الثابت", "Fixed query"), "parameter": ch29TextPair("الـvalue المرتبطة", "Bound value"),
				"result": ch29TextPair("النتيجة", "Result"), "found": ch29TextPair("صف موجود", "Row found"),
				"missing": ch29TextPair("مفيش صف مطابق", "No matching row"), "button": ch29TextPair("بدّل الـkey", "Try another key"),
			},
			"rows": []any{
				map[string]any{"key": ch29TextPair("ridge", "ridge"), "value": ch29TextPair("جاهز", "ready"), "found": true},
				map[string]any{"key": ch29TextPair("cove", "cove"), "value": ch29TextPair("—", "—"), "found": false},
			},
		}
	case "c292":
		return "ch29-migration-boundary", map[string]any{
			"intro": ch29TextPair("جهّز تغيير المخطط وبعدين خزّنه مع النسخة.", "Stage a schema change, then save its version."),
			"labels": map[string]any{
				"saved": ch29TextPair("المخطط المحفوظ", "Saved schema"), "version": ch29TextPair("النسخة المحفوظة", "Saved version"),
				"proposal": ch29TextPair("التغيير المقترح", "Proposed change"), "stage": ch29TextPair("جهّز الترحيل", "Stage migration"),
				"commit": ch29TextPair("خزّن الاتنين", "Commit both"), "reset": ch29TextPair("ابدأ من الأول", "Reset"),
				"pending": ch29TextPair("لسه ما اتحفظش", "Not saved yet"), "committed": ch29TextPair("المخطط والنسخة اتحفظوا سوا", "Schema and version saved together"),
			},
			"schema": ch29TextPair("crates(label)", "crates(label)"), "nextSchema": ch29TextPair("crates(label, zone)", "crates(label, zone)"),
			"version": 2,
		}
	case "c293":
		return "ch29-atomic-transfer", map[string]any{
			"intro": ch29TextPair("جرّب ترتيبين صالحين لـrequests متزامنة.", "Try two valid orders for concurrent requests."),
			"labels": map[string]any{
				"possibleOrder": ch29TextPair("ترتيب ممكن", "One possible order"), "source": ch29TextPair("المصدر", "Source"),
				"north": ch29TextPair("خزان الشمال", "North tank"), "south": ch29TextPair("خزان الجنوب", "South tank"),
				"run": ch29TextPair("جرّب ترتيب تاني", "Try another order"), "accepted": ch29TextPair("اتنقل", "Moved"),
				"skipped": ch29TextPair("الرصيد ما كفّاش", "Not enough left"), "vary": ch29TextPair("الترتيب ممكن يتغير.", "The order can vary."),
			},
			"startingBalance": 4,
			"requests": []any{
				map[string]any{"name": ch29TextPair("الـrequest A", "Request A"), "amount": 3, "target": "north"},
				map[string]any{"name": ch29TextPair("الـrequest B", "Request B"), "amount": 2, "target": "south"},
			},
			"orders": [][]int{{0, 1}, {1, 0}},
		}
	case "c294":
		return "ch29-pool-state", map[string]any{
			"intro": ch29TextPair("حرّك اتصال وشوف الفرق بين الحد والحالة.", "Move a connection; compare the limit with current state."),
			"labels": map[string]any{
				"max": ch29TextPair("الحد المفتوح", "Open limit"), "inUse": ch29TextPair("قيد الاستخدام", "In use"),
				"idle": ch29TextPair("متاح", "Idle"), "acquire": ch29TextPair("استخدم اتصال", "Acquire connection"),
				"release": ch29TextPair("حرّر اتصال", "Release connection"), "reset": ch29TextPair("ابدأ من الأول", "Reset"),
			},
			"maxOpen": 3, "inUse": 1, "idle": 1,
		}
	case "c295":
		return "ch29-repository-boundary", map[string]any{
			"intro": ch29TextPair("بدّل التنفيذ وشوف عقد الـinterface ثابت.", "Switch implementations; the interface contract stays fixed."),
			"labels": map[string]any{
				"contract": ch29TextPair("عقد التخزين", "Store contract"), "implementation": ch29TextPair("التنفيذ الحالي", "Current implementation"),
				"call": ch29TextPair("النداء", "Call"), "switch": ch29TextPair("بدّل التنفيذ", "Switch implementation"),
			},
			"contract": "Put(label, units)", "call": "Put(\"north\", 3)",
			"implementations": []any{
				map[string]any{"name": ch29TextPair("مخزن في الذاكرة", "Memory store"), "effect": ch29TextPair("الـvalue وصلت للـfake.", "The value reached the fake.")},
				map[string]any{"name": ch29TextPair("Repository بـsqlx", "sqlx repository"), "effect": ch29TextPair("الـrepository بيربط الاسم بالـvalue.", "The repository binds name to value.")},
			},
		}
	case "c296":
		return "ch29-association", map[string]any{
			"intro": ch29TextPair("حمّل أبناء القبة وشوف الـassociation.", "Load the dome's children and inspect the association."),
			"labels": map[string]any{
				"parent": ch29TextPair("القبة", "Dome"), "children": ch29TextPair("الأحواض", "Beds"),
				"load": ch29TextPair("حمّل الأحواض", "Load beds"), "hide": ch29TextPair("اخفي الأحواض", "Hide beds"),
				"unloaded": ch29TextPair("العلاقة لسه ما اتحملتش", "Association not loaded"), "loaded": ch29TextPair("اتحمّلت العلاقة", "Association loaded"),
			},
			"parentName": ch29TextPair("Vega", "Vega"),
			"children":   []any{ch29TextPair("داخلي", "inner"), ch29TextPair("خارجي", "outer")},
		}
	default:
		panic("unknown chapter 29 sandbox prefix: " + prefix)
	}
}

var ch29SQLiteIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "err", Value: "wrapped sql.ErrNoRows"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "err", Value: "wrapped sql.ErrNoRows"}}, Output: "true"},
}

var ch29SQLiteTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "arg.Name", Value: "label"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "arg.Name", Value: "label"}}, Output: "label"},
}

var ch29SQLiteArabic = ch29LessonCopy{
	Mission: `خزّن صناديق المؤن في SQLite بقيود والـqueries قابلة للإلغاء.

الـstarter فيه توقيعات CRUD و~*sql.DB~ محقونة؛ سيب main فاضية.`,
	Idea: `~database/sql~ بيفصل الـinterface بتاعة Go عن الدرايفر.

اربط الـvalues كـarguments، واقرا الصفوف في متغيرات بعناوين. غياب الصف يظهر كـ~sql.ErrNoRows~.`,
	IdeaCode: `func main() {
	err := fmt.Errorf("lookup: %w", sql.ErrNoRows)
	fmt.Println(errors.Is(err, sql.ErrNoRows))
}`,
	IdeaImports: []string{"database/sql", "errors", "fmt"}, IdeaID: "e291-error-chain", IdeaStates: ch29SQLiteIdeaStates,
	IdeaCaptions: []string{"قبل التشغيل.", "الـerror محفوظ في err.", "~errors.Is~ لقت السبب."},
	TryPrompt:    `إيه اسم argument اللي بترجّعه ~sql.Named~؟`,
	TryCode: `func main() {
	arg := sql.Named("label", "Luna")
	fmt.Println(arg.Name)
}`,
	TryImports: []string{"database/sql", "fmt"}, TryID: "e291-named-arg", TryAnswer: 1,
	TryOptions:  []string{"Luna", "label", "sql.Named"},
	TryFeedback: []string{"Luna هي الـvalue المربوطة، مش الاسم.", "~sql.NamedArg~ بتحفظ اسم المعامل في ~Name~.", "دي الـfunction اللي بتنشيء الـvalue، مش value ~Name~."},
	TryStates:   ch29SQLiteTryStates,
	TryCaptions: []string{"قبل التشغيل.", "argument باسمه label.", "اطبع اسم المعامل."},
	Reveal:      `شوف الاسم والـvalue منفصلين، وبعدها اطبع الاسم.`,
	Watch: `استخدم placeholders بدل لصق النص في SQL.

اقفل ~Rows~ حتى لو عملت return بدري، وافحص الـerror اللي ظهر من اللف.`,
	Plan: `ابدأ بالمخطط والقيود، وبعدها افصل الكتابة عن القراءة.

راجع نتيجة كل نداء بدل ما تعتبر الـerror صف فاضي. راجع كمان ~Scan~.`,
	CheckQuestion: `إيه اللي يشير له ~sql.ErrNoRows~ بعد Scan؟`,
	CheckOptions:  []string{"مفيش صف مطابق", "اتصال القاعدة اتقفل", "الـquery نجح من غير أعمدة"}, CheckAnswer: 0,
	CheckFeedback: []string{"Scan بيستخدم الـerror ده لما الصف المطلوب غايب.", "إغلاق الاتصال له الـerror التاني.", "الصف ممكن يفضل موجود حتى لو قيمته فاضية."},
	Turn: `شغّل الـfunction واحدة كل مرة، وراجع نوع الـerror بتاعها.

ابدأ بتشغيل المحاولة على قاعدة الـtest المحقونة.`,
	Hints: []string{
		`خلي قيود القاعدة تحمي البيانات حتى لو فات تحقق Go.`,
		`افصل مسارات تنفيذ الأمر عن قراية صف واحد أو مجموعة.`,
		`راجع ~QueryRowContext~ و~RowsAffected~ عند الـfunctions بتاعة CRUD.`,
	},
}

var ch29SQLiteEnglish = ch29LessonCopy{
	Mission: `Persist supply crates in SQLite with constraints and cancelable queries.

The starter injects CRUD signatures and ~*sql.DB~; keep main empty.`,
	Idea: `~database/sql~ separates Go's API from its driver.

Bind values as arguments and scan rows into pointers. A missing row appears as ~sql.ErrNoRows~.`,
	IdeaCode: ch29SQLiteArabic.IdeaCode, IdeaImports: ch29SQLiteArabic.IdeaImports, IdeaID: ch29SQLiteArabic.IdeaID, IdeaStates: ch29SQLiteIdeaStates,
	IdeaCaptions: []string{"Before execution.", "err holds the wrapped error.", "~errors.Is~ found its cause."},
	TryPrompt:    `What argument name does ~sql.Named~ return?`,
	TryCode:      ch29SQLiteArabic.TryCode, TryImports: ch29SQLiteArabic.TryImports, TryID: ch29SQLiteArabic.TryID, TryAnswer: 1,
	TryOptions:  []string{"Luna", "label", "sql.Named"},
	TryFeedback: []string{"Luna is the bound value, not the name.", "~sql.NamedArg~ stores the parameter name in ~Name~.", "That is the constructor, not the ~Name~ value."},
	TryStates:   ch29SQLiteTryStates,
	TryCaptions: []string{"Before execution.", "An argument named label.", "Print the parameter name."},
	Reveal:      `See that name and value are separate, then print the name.`,
	Watch: `Use placeholders instead of joining text into SQL.

Close ~Rows~ on early returns, and check its final error.`,
	Plan: `Start with schema and constraints, then separate writes from reads.

Inspect each call's result instead of treating errors as empty rows.`,
	CheckQuestion: `What does ~sql.ErrNoRows~ mean after Scan?`,
	CheckOptions:  []string{"No row matched", "The database connection closed", "The query succeeded without columns"}, CheckAnswer: 0,
	CheckFeedback: []string{"Scan uses this error when the requested row is absent.", "A closed connection has a different error.", "A row can exist even when a value is empty."},
	Turn: `Implement one function at a time and check its error type.

Start by running your attempt on the injected test database.`,
	Hints: []string{
		`Let database constraints protect data when Go validation is bypassed.`,
		`Separate commands from reading one row or a collection.`,
		`Review ~QueryRowContext~ and ~RowsAffected~ across the CRUD functions.`,
	},
}

var ch29MigrationIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "applied", Value: "[2]"}}},
	{Line: 4, Variables: []LessonVariable{{Name: "applied", Value: "[2]"}, {Name: "next", Value: "[2 5]"}}},
	{Line: 5, Variables: []LessonVariable{{Name: "applied", Value: "[2]"}, {Name: "next", Value: "[2 5]"}}, Output: "[2] [2 5]"},
}

var ch29MigrationTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "blank", Value: "true"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "blank", Value: "true"}}, Output: "true"},
}

var ch29MigrationArabic = ch29LessonCopy{
	Mission: `طوّر مخطط المخزن بترحيلات up/down وسجل نسخ ذري.

الـstarter يديك نوع Migration والـfunction الواحدة؛ المنادي يملك قاعدة البيانات.`,
	Idea: `سجل النسخ بيحفظ تاريخ المخطط بعد إعادة التشغيل.

المثال ده بيقلّد فكرة التجهيز؛ قاعدة البيانات نفسها تحتاج ~sql.Tx~ لضم التغيير والسجل.`,
	IdeaCode: `func main() {
	applied := []int{2}
	next := slices.Clone(applied)
	next = append(next, 5)
	fmt.Println(applied, next)
}`,
	IdeaImports: []string{"fmt", "slices"}, IdeaID: "e292-staged-list", IdeaStates: ch29MigrationIdeaStates,
	IdeaCaptions: []string{"قبل التشغيل.", "السجل الحالي نسخة واحدة.", "جهّز نسخة للتغيير.", "اطبع الحالي والمجهّز."},
	TryPrompt:    `هل النص ده فاضي بعد ~strings.TrimSpace~؟`,
	TryCode: `func main() {
	blank := strings.TrimSpace("   ") == ""
	fmt.Println(blank)
}`,
	TryImports: []string{"fmt", "strings"}, TryID: "e292-trim-script", TryAnswer: 0,
	TryOptions:  []string{"true", "false", "panic"},
	TryFeedback: []string{"بعد إزالة المسافات مفيش نص باقي.", "المقارنة بتشوف النص بعد trim.", "~strings.TrimSpace~ ما بتعملش panic للنص ده."},
	TryStates:   ch29MigrationTryStates,
	TryCaptions: []string{"قبل التشغيل.", "blank بقت true.", "اطبع النتيجة."},
	Reveal:      `شيل المسافات، وبعدها راقب نتيجة المقارنة.`,
	Watch: `الـDown مش backup؛ ممكن يمسح بيانات.

كل خطوة تحفظ أثر المخطط ورقم النسخة مع بعض.`,
	Plan: `اتأكد من المصدر والتاريخ قبل اختيار الاتجاه.

فكر في حدود كل خطوة: إيه اللي يتراجع وإيه اللي يفضل محفوظ؟`,
	CheckQuestion: `لو ترحيل لاحق فشل، إيه المفروض يحصل للخطوة السابقة المحفوظة؟`,
	CheckOptions:  []string{"تفضل مطبقة", "تتراجع كلها", "تتمسح من السجل بس"}, CheckAnswer: 0,
	CheckFeedback: []string{"كل ترحيل ناجح له حفظ مستقل.", "المهمة تحفظ التقدم اللي سبق الـerror.", "المخطط والسجل لازم يفضلوا متوافقين."},
	Turn: `اقرا ترتيب النسخ قبل لمس قاعدة البيانات.

بعدها جرّب اتجاه up وdown على تاريخ صغير.`,
	Hints: []string{
		`السجل الدائم يربط نسخة المخطط بحالتها الحالية.`,
		`قسّم التفكير إلى تحقق المصدر وقراءة التاريخ وتنفيذ خطوة.`,
		`راجع ~BeginTx~ و~Commit~ في توقيع ~MigrateHabitat~.`,
	},
}

var ch29MigrationEnglish = ch29LessonCopy{
	Mission: `Evolve the store schema with up/down migrations and an atomic version ledger.

The starter provides Migration and one function; the caller owns the database.`,
	Idea: `The version ledger keeps schema history across restarts.

This example models staging; the database needs ~sql.Tx~ to group the change and ledger.`,
	IdeaCode: ch29MigrationArabic.IdeaCode, IdeaImports: ch29MigrationArabic.IdeaImports, IdeaID: ch29MigrationArabic.IdeaID, IdeaStates: ch29MigrationIdeaStates,
	IdeaCaptions: []string{"Before execution.", "The current ledger is one slice.", "Prepare a copy to change.", "Print current and staged."},
	TryPrompt:    `Is this text blank after ~strings.TrimSpace~?`,
	TryCode:      ch29MigrationArabic.TryCode, TryImports: ch29MigrationArabic.TryImports, TryID: ch29MigrationArabic.TryID, TryAnswer: 0,
	TryOptions:  []string{"true", "false", "panic"},
	TryFeedback: []string{"No text remains after trimming spaces.", "The comparison sees the trimmed text.", "~strings.TrimSpace~ does not panic on this input."},
	TryStates:   ch29MigrationTryStates,
	TryCaptions: []string{"Before execution.", "blank becomes true.", "Print the result."},
	Reveal:      `Remove the spaces, then watch the comparison result.`,
	Watch: `A down migration is not a backup; it can remove data.

Each step saves schema effects and its version together.`,
	Plan: `Validate source and history before choosing a direction.

Think about each step's boundary: what rolls back and what stays saved?`,
	CheckQuestion: `If a later migration fails, what happens to an earlier saved step?`,
	CheckOptions:  []string{"It stays applied", "Everything rolls back", "Only its ledger row is deleted"}, CheckAnswer: 0,
	CheckFeedback: []string{"Each successful migration has its own commit.", "The mission preserves progress before the failure.", "Schema and ledger must stay consistent."},
	Turn: `Read the version order before changing the database.

Then try up and down against a small history.`,
	Hints: []string{
		`A durable ledger ties schema versions to current state.`,
		`Separate source validation, history reading, and step execution.`,
		`Review ~BeginTx~ and ~Commit~ around ~MigrateHabitat~.`,
	},
}

var ch29TransferIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "from", Value: "9"}, {Name: "to", Value: "2"}, {Name: "moved", Value: "3"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "from", Value: "9"}, {Name: "to", Value: "2"}, {Name: "moved", Value: "3"}}, Output: "11 11"},
}

var ch29TransferTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "balance", Value: "5"}, {Name: "requested", Value: "2"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "balance", Value: "5"}, {Name: "requested", Value: "2"}}, Output: "true 3"},
}

var ch29TransferArabic = ch29LessonCopy{
	Mission: `انقل وحدات أكسجين بين خزانات من غير فقد أو تكرار تحت التزامن.

الـstarter يحدد Move وhook اختيارية بعد الخصم.`,
	Idea: `التحويل بيغيّر رصيدين وسجل request في عملية واحدة.

المجموع يفضل ثابت؛ الـtransaction بتخلي أي فشل يلغي التغييرات كلها مع بعض.`,
	IdeaCode: `func main() {
	from, to, moved := int64(9), int64(2), int64(3)
	fmt.Println(from+to, from-moved+to+moved)
}`,
	IdeaImports: []string{"fmt"}, IdeaID: "e293-conservation", IdeaStates: ch29TransferIdeaStates,
	IdeaCaptions: []string{"قبل التشغيل.", "الأرصدة والكمية واضحة.", "المجموع ثابت بعد النقل."},
	TryPrompt:    `توقّع هل الرصيد يكفي، والرصيد المتبقي.`,
	TryCode: `func main() {
	balance, requested := int64(5), int64(2)
	fmt.Println(balance >= requested, balance-requested)
}`,
	TryImports: []string{"fmt"}, TryID: "e293-balance-check", TryAnswer: 0,
	TryOptions:  []string{"true 3", "false 7", "true 7"},
	TryFeedback: []string{"الـ5 تكفي لـrequest من وحدتين؛ يتبقى 3.", "الرصيد الابتدائي أكبر من الـrequest.", "الطرح يقلل رصيد المصدر."},
	TryStates:   ch29TransferTryStates,
	TryCaptions: []string{"قبل التشغيل.", "المتاح أكبر من المطلوب.", "اطبع الكفاية والباقي."},
	Reveal:      `قارن الرصيد بالـrequest، وبعدها احسب الباقي.`,
	Watch: `القراية ثم الكتابة برّه الـtransaction تفتح الـrace.

الـhook ممكن يفشل بعد الخصم؛ لازم التغيير كله يتراجع.`,
	Plan: `راجع المدخلات والـkey بتاع الـrequest قبل تغيير الرصيد.

بعدها افحص أثر الخصم والإضافة، وفكر في الإلغاء والفشل الوسيط.`,
	CheckQuestion: `لو hook رجّعت الـerror بعد الخصم، إيه الحالة المتوقعة؟`,
	CheckOptions:  []string{"الرصيدين والسجل زي الأول", "الخصم محفوظ والإضافة ملغية", "الـrequest يتسجل رغم الـerror"}, CheckAnswer: 0,
	CheckFeedback: []string{"العملية الذرية بتلغي التغييرات كلها مع بعض عند الفشل.", "ده يضيّع وحدات لو ما حصلش rollback.", "السجل ما يعلنش عملية ما اكتملتش."},
	Turn: `جرّب مسار نجاح ومسار hook يفشل.

راجع الأرصدة وسجل الـrequest بعد كل نتيجة.`,
	Hints: []string{
		`التحويل عقده حفظ مجموع الأرصدة وتفرّد الـrequest.`,
		`افصل التحقق والخصم والـhook والإضافة وتسجيل الـrequest.`,
		`راجع ~RowsAffected~ و~afterDebit~ جوه ~MoveOxygen~.`,
	},
}

var ch29TransferEnglish = ch29LessonCopy{
	Mission: `Transfer oxygen between tanks without loss or replay under concurrency.

The starter defines Move and an optional hook after the debit.`,
	Idea: `A transfer changes two balances and one request ledger entry.

The total stays constant; a transaction lets any failure undo the changes together.`,
	IdeaCode: ch29TransferArabic.IdeaCode, IdeaImports: ch29TransferArabic.IdeaImports, IdeaID: ch29TransferArabic.IdeaID, IdeaStates: ch29TransferIdeaStates,
	IdeaCaptions: []string{"Before execution.", "Balances and amount are set.", "The total stays constant."},
	TryPrompt:    `Predict whether the balance is enough and what remains.`,
	TryCode:      ch29TransferArabic.TryCode, TryImports: ch29TransferArabic.TryImports, TryID: ch29TransferArabic.TryID, TryAnswer: 0,
	TryOptions:  []string{"true 3", "false 7", "true 7"},
	TryFeedback: []string{"Five covers a request of two, leaving three.", "The starting balance is greater than the request.", "Subtracting reduces the source balance."},
	TryStates:   ch29TransferTryStates,
	TryCaptions: []string{"Before execution.", "Available is above requested.", "Print enough and remainder."},
	Reveal:      `Compare balance with the request, then calculate the remainder.`,
	Watch: `Reading then writing outside a transaction creates a race.

The hook can fail after debit; the entire change must roll back.`,
	Plan: `Check inputs and the request key before changing balances.

Then inspect debit and credit effects, including cancellation and mid-step failure.`,
	CheckQuestion: `If the hook returns an error after debit, what state is expected?`,
	CheckOptions:  []string{"Both balances and ledger stay as before", "Debit stays, credit is canceled", "The request is recorded anyway"}, CheckAnswer: 0,
	CheckFeedback: []string{"An atomic operation rolls all of its writes back together.", "That would lose units without a rollback.", "The ledger must not mark an incomplete operation."},
	Turn: `Try a successful path and a failing hook path.

Inspect balances and the request ledger after each result.`,
	Hints: []string{
		`The transfer contract preserves total balance and request uniqueness.`,
		`Separate validation, debit, hook, credit, and request recording.`,
		`Review ~RowsAffected~ and ~afterDebit~ inside ~MoveOxygen~.`,
	},
}

var ch29PoolIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "stats.MaxOpenConnections", Value: "3"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "stats.InUse", Value: "2"}, {Name: "stats.Idle", Value: "1"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "stats.MaxOpenConnections", Value: "3"}, {Name: "stats.InUse", Value: "2"}, {Name: "stats.Idle", Value: "1"}}, Output: "3 2 1"},
}

var ch29PoolTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "stats.MaxOpenConnections", Value: "4"}, {Name: "stats.Idle", Value: "2"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "stats.MaxOpenConnections", Value: "4"}, {Name: "stats.Idle", Value: "2"}}, Output: "4 2"},
}

var ch29PoolArabic = ch29LessonCopy{
	Mission: `اضبط حوض قارئ بوابات الترحيل، وحرر الصفوف حتى عند الفشل.

الـstarter فيه دالتين على ~*sql.DB~؛ المنادي يملك عمر القاعدة.`,
	Idea: `~*sql.DB~ حوض آمن للمشاركة، مش اتصال واحد.

~DBStats~ لقطة للحالة؛ أرقام المثال للتوضيح. ~Idle~ عدد متاح دلوقتي، مش حد الاتصالات الفاضية.`,
	IdeaCode: `func main() {
	stats := sql.DBStats{MaxOpenConnections: 3, InUse: 2, Idle: 1}
	fmt.Println(stats.MaxOpenConnections, stats.InUse, stats.Idle)
}`,
	IdeaImports: []string{"database/sql", "fmt"}, IdeaID: "e294-pool-stats", IdeaStates: ch29PoolIdeaStates,
	IdeaCaptions: []string{"قبل التشغيل.", "الحد المضبوط ثلاثة.", "اثنان شغالان وواحد فاضي.", "اطبع اللقطة."},
	TryPrompt:    `إيه الـvalues اللي بيطبعهم المثال؟`,
	TryCode: `func main() {
	stats := sql.DBStats{MaxOpenConnections: 4, Idle: 2}
	fmt.Println(stats.MaxOpenConnections, stats.Idle)
}`,
	TryImports: []string{"database/sql", "fmt"}, TryID: "e294-stats-fields", TryAnswer: 1,
	TryOptions:  []string{"2 4", "4 2", "4 4"},
	TryFeedback: []string{"الـvalues بتتطبع حسب ترتيب الحقول في Println.", "الحد المضبوط أربعة؛ والـidle الحالي اتنين.", "Idle مش نسخة من الحد الأقصى."},
	TryStates:   ch29PoolTryStates,
	TryCaptions: []string{"قبل التشغيل.", "الحد أربعة وIdle اتنين.", "اطبع الحقلين."},
	Reveal:      `اقرا الحقلين بالترتيب اللي اتبعتوا به.`,
	Watch: `ما تعتبرش ~WaitDuration~ نتيجة حتمية؛ ده وقت حقيقي متراكم.

~Rows~ بتحجز سعة لحد ما تتقفل.`,
	Plan: `اتأكد من إعداد الحوض قبل تغييره.

بعد الـquery، تابع ~Scan~ والتنظيف والـerror اللي بيطلع من اللف.

افصل الحد الحالي عن الحد المضبوط.`,
	CheckQuestion: `إيه اللي لازم يحصل لصفوف ~QueryContext~ عند return بدري؟`,
	CheckOptions:  []string{"تتقفل", "تفضل مفتوحة لحد Close للقاعدة", "تتقرأ تاني من البداية"}, CheckAnswer: 0,
	CheckFeedback: []string{"~Close~ يحرر الاتصال حتى لو ~Scan~ فشل.", "ده ممكن يمنع الـquery التاني من استخدام الحوض.", "~Next~ ما بيرجعش للقيم اللي اتقرت."},
	Turn: `ابدأ بضبط الـmax على value صالحة، وبعدها اقرا بيانات صغيرة.

راجع سعة الحوض بعد النجاح وبعد الـerror في القراية.`,
	Hints: []string{
		`حوض الاتصالات يدير سعة مشتركة بين قراء متعددين.`,
		`قسّم الدالتين بين إعداد الحد وقراءة النتائج.`,
		`راجع ~SetMaxOpenConns~ و~Rows.Close~ في starter.`,
	},
}

var ch29PoolEnglish = ch29LessonCopy{
	Mission: `Bound the relay-port reader pool and release rows even on failure.

The starter has two functions over ~*sql.DB~; the caller owns database lifetime.`,
	Idea: `~*sql.DB~ is a shareable pool, not one connection.

~DBStats~ is a snapshot; its sample counts are illustrative. ~Idle~ is free connections now, not the idle limit.`,
	IdeaCode: ch29PoolArabic.IdeaCode, IdeaImports: ch29PoolArabic.IdeaImports, IdeaID: ch29PoolArabic.IdeaID, IdeaStates: ch29PoolIdeaStates,
	IdeaCaptions: []string{"Before execution.", "The configured limit is three.", "Two are in use, one is idle.", "Print the snapshot."},
	TryPrompt:    `Which two values does this example print?`,
	TryCode:      ch29PoolArabic.TryCode, TryImports: ch29PoolArabic.TryImports, TryID: ch29PoolArabic.TryID, TryAnswer: 1,
	TryOptions:  []string{"2 4", "4 2", "4 4"},
	TryFeedback: []string{"Values print in the order passed to Println.", "The limit is four; the current idle count is two.", "Idle is not a copy of the maximum."},
	TryStates:   ch29PoolTryStates,
	TryCaptions: []string{"Before execution.", "Limit is four; Idle is two.", "Print both fields."},
	Reveal:      `Read the fields in the order they were passed.`,
	Watch: `Do not treat ~WaitDuration~ as deterministic; it is real accumulated time.

~Rows~ holds capacity until it is closed.`,
	Plan: `Validate pool settings before changing them.

After a query, trace ~Scan~, cleanup, and the terminal iteration error.

Separate current use from configured capacity.`,
	CheckQuestion: `What should happen to ~QueryContext~ rows on an early return?`,
	CheckOptions:  []string{"They close", "They stay open until DB.Close", "They are reread from the beginning"}, CheckAnswer: 0,
	CheckFeedback: []string{"~Close~ releases the connection even when ~Scan~ fails.", "That can block another query from using the pool.", "~Next~ does not go back to values already read."},
	Turn: `Start with a valid max setting, then read a small result set.

Check pool capacity after success and after a scan error.`,
	Hints: []string{
		`A connection pool manages shared capacity for readers.`,
		`Split the functions into limit setup and result reading.`,
		`Review ~SetMaxOpenConns~ and ~Rows.Close~ in the starter.`,
	},
}

var ch29SQLXIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 10, Variables: []LessonVariable{{Name: "store", Value: "memoryStore"}}},
	{Line: 11, Variables: []LessonVariable{{Name: "bay", Value: "east"}, {Name: "units", Value: "3"}}, Output: "east 3"},
}

var ch29SQLXTryStates = []LessonFrame{
	{Line: 0},
	{Line: 2, Variables: []LessonVariable{{Name: "arg.Name", Value: "bay"}}},
	{Line: 3, Variables: []LessonVariable{{Name: "arg.Name", Value: "bay"}}, Output: "bay"},
}

var ch29SQLXArabic = ch29LessonCopy{
	Mission: `افصل شحن المؤن عن repository حقيقي بـ sqlx وPostgreSQL.

الـstarter فيه interface و~SQLRations~؛ الاتصال بييجي من المنادي.`,
	Idea: `الـinterface بتفصل منادي المخزن عن تنفيذه.

تقدر تستخدم fake بسيط لفهم العقد؛ بس SQL الحقيقي محتاج ربط أسماء وstruct tags.

الـinterface بتاعة المجال صغيرة؛ تفاصيل الاتصال تفضل جوه الـrepository.

الـdb tag بيربط اسم Go باسم العمود.`,
	IdeaCode: `type Store interface {
	Put(string, int) error
}
type memoryStore struct{}
func (memoryStore) Put(bay string, units int) error {
	fmt.Println(bay, units)
	return nil
}
func main() {
	var store Store = memoryStore{}
	_ = store.Put("east", 3)
}`,
	IdeaImports: []string{"fmt"}, IdeaID: "e295-interface-store", IdeaStates: ch29SQLXIdeaStates,
	IdeaCaptions: []string{"قبل التشغيل.", "store جاهز بالـvalue.", "Put وصلت للـfake."},
	TryPrompt:    `إيه اللي يظهر كاسم المعامل؟`,
	TryCode: `func main() {
	arg := sql.Named("bay", "north")
	fmt.Println(arg.Name)
}`,
	TryImports: []string{"database/sql", "fmt"}, TryID: "e295-named-arg", TryAnswer: 0,
	TryOptions:  []string{"bay", "north", "$1"},
	TryFeedback: []string{"~Name~ هو bay؛ north هي الـvalue.", "north بتتخزن في ~Value~.", "~$1~ صيغة موضعية مش ~Name~."},
	TryStates:   ch29SQLXTryStates,
	TryCaptions: []string{"قبل التشغيل.", "الاسم bay والـvalue north.", "اطبع الاسم."},
	Reveal:      `فرّق بين اسم الـargument والـvalue المرتبطة به.`,
	Watch: `الـfake يشرح السلوك، لكنه مش دليل على سلوك PostgreSQL.

مرّر context المنادي، وسيب عمر الحوض للمنادي.

القائمة الفاضية والـerrors جزء من العقد.

كل query يستقبل الـcontext بتاع المنادي.

خلي الحدود بين المجال والتخزين باينة.`,
	Plan: `اتأكد من الـrequest بتاع المجال، وبعدها راجع الـrepository.

خلي سلوك القائمة الفاضية والـerrors واضح في العقد.`,
	CheckQuestion: `إيه اللي بيعرفه مستهلك ~StockRation~ عن التخزين؟`,
	CheckOptions:  []string{"الـinterface ~RationStore~", "تفاصيل SQL الداخلية", "اسم جدول PostgreSQL"}, CheckAnswer: 0,
	CheckFeedback: []string{"الـinterface تحدد عمليات المخزن المتاحة للمستهلك.", "التنفيذ يملك تفاصيل الـquery.", "اسم الجدول يخص repository."},
	Turn: `نفّذ ~StockRation~ على الحدود اللي في التوقيع.

بعدها راجع ~SQLRations~ على قاعدة يجهزها المنادي.`,
	Hints: []string{
		`الـrepository يخفي تفاصيل التخزين وراء عقد صغير.`,
		`افصل التحقق عن المجال وربط قاعدة البيانات.`,
		`راجع ~sqlx.NamedExecContext~ و~SelectContext~ عند ~SQLRations~.`,
	},
}

var ch29SQLXEnglish = ch29LessonCopy{
	Mission: `Separate ration stocking from a real sqlx and PostgreSQL repository.

The starter has an interface and ~SQLRations~; the caller supplies the connection.`,
	Idea: `An interface separates the store's consumer from its implementation.

A small fake can explain the contract; real SQL still needs named binding and struct tags.

Keep the domain interface small; connection details stay in the repository.

A db tag maps a Go field to its column name.`,
	IdeaCode: ch29SQLXArabic.IdeaCode, IdeaImports: ch29SQLXArabic.IdeaImports, IdeaID: ch29SQLXArabic.IdeaID, IdeaStates: ch29SQLXIdeaStates,
	IdeaCaptions: []string{"Before execution.", "store holds the value.", "Put reached the fake."},
	TryPrompt:    `Which value is printed as the argument name?`,
	TryCode:      ch29SQLXArabic.TryCode, TryImports: ch29SQLXArabic.TryImports, TryID: ch29SQLXArabic.TryID, TryAnswer: 0,
	TryOptions:  []string{"bay", "north", "$1"},
	TryFeedback: []string{"~Name~ is bay; north is the value.", "north is stored in ~Value~.", "~$1~ is positional syntax, not ~Name~."},
	TryStates:   ch29SQLXTryStates,
	TryCaptions: []string{"Before execution.", "Name is bay; value is north.", "Print the name."},
	Reveal:      `Separate the argument name from its bound value.`,
	Watch: `A fake explains behavior but does not prove PostgreSQL behavior.

Pass the caller's context and leave pool lifetime to the caller.

Empty results and errors remain part of the contract.

Pass the caller's context to each query.

Keep the domain and storage boundary visible.`,
	Plan: `Start with domain request validation, then inspect the repository.

Make empty-list and error behavior clear in the contract.`,
	CheckQuestion: `What does a ~StockRation~ consumer know about storage?`,
	CheckOptions:  []string{"The ~RationStore~ interface", "Internal SQL details", "The PostgreSQL table name"}, CheckAnswer: 0,
	CheckFeedback: []string{"The interface defines operations the consumer can use.", "The implementation owns query details.", "The table name belongs to the repository."},
	Turn: `Implement ~StockRation~ against the signature's boundary.

Then inspect ~SQLRations~ with a database supplied by the caller.`,
	Hints: []string{
		`The repository hides storage details behind a small contract.`,
		`Separate domain validation from database binding.`,
		`Review ~sqlx.NamedExecContext~ and ~SelectContext~ in ~SQLRations~.`,
	},
}

var ch29GORMIdeaStates = []LessonFrame{
	{Line: 0},
	{Line: 7, Variables: []LessonVariable{{Name: "dome.Name", Value: "Vega"}, {Name: "dome.Beds", Value: "2 beds"}}},
	{Line: 8, Variables: []LessonVariable{{Name: "dome.Name", Value: "Vega"}, {Name: "dome.Beds", Value: "2 beds"}}, Output: "Vega 2"},
}

var ch29GORMTryStates = []LessonFrame{
	{Line: 0},
	{Line: 4, Variables: []LessonVariable{{Name: "dome.Beds[1].Label", Value: "outer"}}},
	{Line: 5, Variables: []LessonVariable{{Name: "dome.Beds[1].Label", Value: "outer"}}, Output: "outer"},
}

var ch29GORMArabic = ch29LessonCopy{
	Mission: `اربط قباب الزراعة بأحواضها وحمّل العلاقات من غير N+1.

الـstarter يديك model structs وfunctions ORM ناقصة.

الاسم والـlabels لازم يبقوا صالحين قبل الحفظ.`,
	Idea: `الـassociation بتمثلها slice جوه الـmodel الأب.

في GORM، الـtags بتوضح الـkeys والعلاقة؛ و~Preload~ بيحمّل الأبناء قبل قراية النتيجة.

الـtransaction بتمنع تخزين جزء من المجموعة لو ابن فشل.`,
	IdeaCode: `type Bed struct{ Label string }
type Dome struct {
	Name string
	Beds []Bed
}
func main() {
	dome := Dome{Name: "Vega", Beds: []Bed{{Label: "north"}, {Label: "south"}}}
	fmt.Println(dome.Name, len(dome.Beds))
}`,
	IdeaImports: []string{"fmt"}, IdeaID: "e296-association", IdeaStates: ch29GORMIdeaStates,
	IdeaCaptions: []string{"قبل التشغيل.", "القبة لها اسم وحوضين.", "اطبع الاسم وعدد الأحواض."},
	TryPrompt:    `إيه label الموجودة في index 1؟`,
	TryCode: `type Bed struct{ Label string }
type Dome struct{ Beds []Bed }
func main() {
	dome := Dome{Beds: []Bed{{Label: "inner"}, {Label: "outer"}}}
	fmt.Println(dome.Beds[1].Label)
}`,
	TryImports: []string{"fmt"}, TryID: "e296-association-index", TryAnswer: 1,
	TryOptions:  []string{"inner", "outer", "empty"},
	TryFeedback: []string{"index 0 هو inner.", "index 1 هو العنصر التاني outer.", "الـslice فيها عنصرين."},
	TryStates:   ch29GORMTryStates,
	TryCaptions: []string{"قبل التشغيل.", "اقرا العنصر الثاني.", "اطبع الـlabel."},
	Reveal:      `الـindex صفر يبدأ بأول عنصر؛ واحد يجيب اللي بعده.`,
	Watch: `تحميل كل أبناء كل قبة جوه الـloop يعمل queries كتير.

slice فاضية مش معناها إن القبة نفسها غايبة.

خلي ترتيب الحقول واضح في النتيجة.`,
	Plan: `راجع الـmodel tags والقيود قبل إنشاء العلاقة.

بعدها فكّر في التحميل والترتيب وتحديث الـzero values.`,
	CheckQuestion: `لو القبة مفيهاش أحواض، شكل association الأنسب إيه؟`,
	CheckOptions:  []string{"slice فاضية", "قبة غايبة", "nil دايماً"}, CheckAnswer: 0,
	CheckFeedback: []string{"القبة موجودة وعلاقتها بس فاضية.", "غياب الأب حالة مختلفة.", "العقد ممكن يطلب نتيجة empty واضحة."},
	Turn: `ابدأ بتعريف العلاقة، وبعدها راجع القراءة والتحديث.

راقب ترتيب الآباء والأبناء في النتيجة.`,
	Hints: []string{
		`العلاقة تجمع owner وأبناء مع قواعد مستقلة.`,
		`افصل تعريف الموديل عن إنشاء المجموعة وقراءتها.`,
		`راجع tags و~Preload~ في ~GardenDome~ و~ListGarden~.`,
	},
}

var ch29GORMEnglish = ch29LessonCopy{
	Mission: `Relate garden domes to beds and load associations without N+1 queries.

The starter provides model structs and unfinished ORM functions.

Names and labels must be valid before saving.`,
	Idea: `A parent model can represent its association with a slice.

In GORM, tags describe keys and relations; ~Preload~ loads children before reading results.

A transaction prevents a partial aggregate if a child insert fails.`,
	IdeaCode: ch29GORMArabic.IdeaCode, IdeaImports: ch29GORMArabic.IdeaImports, IdeaID: ch29GORMArabic.IdeaID, IdeaStates: ch29GORMIdeaStates,
	IdeaCaptions: []string{"Before execution.", "The dome has a name and two beds.", "Print the name and bed count."},
	TryPrompt:    `Which label is at index 1?`,
	TryCode:      ch29GORMArabic.TryCode, TryImports: ch29GORMArabic.TryImports, TryID: ch29GORMArabic.TryID, TryAnswer: 1,
	TryOptions:  []string{"inner", "outer", "empty"},
	TryFeedback: []string{"Index 0 is inner.", "Index 1 is the second item, outer.", "The slice contains two items."},
	TryStates:   ch29GORMTryStates,
	TryCaptions: []string{"Before execution.", "Read the second item.", "Print the label."},
	Reveal:      `Index zero starts at the first item; one reads the next.`,
	Watch: `Loading each dome's children inside a loop makes many queries.

An empty slice does not mean the parent is missing.

Keep field order clear in the returned result.`,
	Plan: `Review model tags and constraints before creating the relation.

Then consider loading, ordering, and zero-value updates.`,
	CheckQuestion: `If a dome has no beds, what should its association look like?`,
	CheckOptions:  []string{"An empty slice", "A missing dome", "Always nil"}, CheckAnswer: 0,
	CheckFeedback: []string{"The dome exists; only its relation is empty.", "A missing parent is a different case.", "The contract can require a clear empty result."},
	Turn: `Start by defining the relation, then review reads and updates.

Watch the order of parents and children in the result.`,
	Hints: []string{
		`A relation joins an owner and children with separate rules.`,
		`Separate model definition from aggregate creation and reading.`,
		`Review tags and ~Preload~ on ~GardenDome~ and ~ListGarden~.`,
	},
}

var ch29SQLiteArabicSteps = ch29BuildSteps("c291", false, ch29SQLiteArabic)
var ch29SQLiteEnglishSteps = ch29BuildSteps("c291", true, ch29SQLiteEnglish)
var ch29MigrationArabicSteps = ch29BuildSteps("c292", false, ch29MigrationArabic)
var ch29MigrationEnglishSteps = ch29BuildSteps("c292", true, ch29MigrationEnglish)
var ch29TransferArabicSteps = ch29BuildSteps("c293", false, ch29TransferArabic)
var ch29TransferEnglishSteps = ch29BuildSteps("c293", true, ch29TransferEnglish)
var ch29PoolArabicSteps = ch29BuildSteps("c294", false, ch29PoolArabic)
var ch29PoolEnglishSteps = ch29BuildSteps("c294", true, ch29PoolEnglish)
var ch29SQLXArabicSteps = ch29BuildSteps("c295", false, ch29SQLXArabic)
var ch29SQLXEnglishSteps = ch29BuildSteps("c295", true, ch29SQLXEnglish)
var ch29GORMArabicSteps = ch29BuildSteps("c296", false, ch29GORMArabic)
var ch29GORMEnglishSteps = ch29BuildSteps("c296", true, ch29GORMEnglish)
