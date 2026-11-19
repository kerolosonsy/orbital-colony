package levels

func ch37Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch37Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch37Prediction(text, code string) string {
	return text + "\n\n~~~go\n" + code + "\n~~~"
}

func ch37Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

func ch37Pair(ar, en string) map[string]any {
	return map[string]any{"ar": ar, "en": en}
}

func ch37Sandbox(id, text, kind string, config map[string]any) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: kind, Config: config}}
}

func ch37AtomicConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"title":  ch37Pair("تغييرين مع بعض", "Two changes together"),
			"schema": ch37Pair("schema", "schema"), "ledger": ch37Pair("سجل النسخ", "version ledger"),
			"fail": ch37Pair("محاكاة فشل", "Simulate failure"), "commit": ch37Pair("محاكاة نجاح", "Simulate commit"), "reset": ch37Pair("ابدأ تاني", "Reset"),
			"ready": ch37Pair("جاهز لتغيير تجريبي", "Ready for a sample change"), "failed": ch37Pair("فشل؛ الاتنين زي ما هم", "Failed; both values stay"), "committed": ch37Pair("نجاح؛ القيمتين اتحدثوا سوا", "Committed; both values moved"),
		},
		"values": map[string]any{"schema": 4, "ledger": 4, "next": 5},
	}
}

func ch37SessionConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"title":  ch37Pair("التوكن ومخزونه", "Token and stored digest"),
			"client": ch37Pair("مع العميل", "Client holds"), "store": ch37Pair("المخزن", "Store holds"),
			"issue": ch37Pair("إصدار جلسة", "Issue session"), "logout": ch37Pair("خروج", "Log out"), "reset": ch37Pair("ابدأ تاني", "Reset"),
			"empty": ch37Pair("مفيش", "None"), "ready": ch37Pair("مفيش جلسة", "No active session"), "issued": ch37Pair("العميل ماسك التوكن", "Client holds the token"), "ended": ch37Pair("الجلسة اتقفلت", "Session ended"), "cleared": ch37Pair("الحالة اتصفّرت", "Session state cleared"),
		},
		"values": map[string]any{"token": "token-amber", "digest": "digest-amber"},
	}
}

func ch37CrewConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"title": ch37Pair("المفتاح فيه نطاق", "The key includes scope"),
			"key":   ch37Pair("المفتاح المركّب", "Composite key"), "role": ch37Pair("الدور", "Role"), "choose": ch37Pair("اختار نطاق", "Choose a scope"),
			"blue": ch37Pair("blue + mira", "blue + mira"), "gold": ch37Pair("gold + mira", "gold + mira"), "red": ch37Pair("red + mira", "red + mira"), "empty": ch37Pair("—", "—"),
		},
		"lookups": []any{
			map[string]any{"key": ch37Pair("blue + mira", "blue + mira"), "role": ch37Pair("مالك", "owner")},
			map[string]any{"key": ch37Pair("gold + mira", "gold + mira"), "role": ch37Pair("عضو", "member")},
			map[string]any{"key": ch37Pair("red + mira", "red + mira"), "role": ch37Pair("مش موجود", "not found")},
		},
	}
}

func ch37CursorConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"title": ch37Pair("صفحات بمؤشر ثابت", "Pages with a stable marker"),
			"page":  ch37Pair("الصفحة", "Page"), "cursor": ch37Pair("cursor", "cursor"),
			"next": ch37Pair("الصفحة الجاية", "Next page"), "reset": ch37Pair("ارجع للبداية", "Back to start"),
			"ready": ch37Pair("أول صفحة", "First page"), "advanced": ch37Pair("المؤشر اتحرك مع الصفحة", "The marker moved with the page"), "resetStatus": ch37Pair("رجعنا لأول صفحة", "Returned to the first page"),
		},
		"pages":   [][]any{{"task-amber", "task-blue"}, {"task-copper", "task-violet"}},
		"cursors": []string{"start", "after-amber"},
	}
}

func ch37OperationsConfig() map[string]any {
	return map[string]any{
		"labels": map[string]any{
			"title":    ch37Pair("مسار الطلب", "Request path"),
			"previous": ch37Pair("اللي قبلها", "Previous"), "next": ch37Pair("اللي بعدها", "Next"), "reset": ch37Pair("من الأول", "Reset"),
		},
		"frames": []any{
			map[string]any{"stage": ch37Pair("وصول", "Arrival"), "detail": "GET /v1/tasks?sample=1", "caption": ch37Pair("طلب تجريبي وصل", "A sample request arrived")},
			map[string]any{"stage": ch37Pair("هوية الطلب", "Request identity"), "detail": "req-42", "caption": ch37Pair("اتضاف request ID", "A request ID was added")},
			map[string]any{"stage": ch37Pair("الرد", "Response"), "detail": "202 Accepted", "caption": ch37Pair("الرد رجع للعميل", "The response returned to the client")},
			map[string]any{"stage": ch37Pair("السجل", "Log"), "detail": "GET /v1/tasks · 202", "caption": ch37Pair("السجل فيه path من غير query", "The log keeps the path, not query")},
		},
	}
}

const ch37TransactionExample = `func stage(applied []int, version int, commit bool) []int {
	if !commit {
		return applied
	}
	return append(applied, version)
}
func main() {
	current := []int{2}
	fmt.Println(stage(current, 3, false))
	fmt.Println(stage(current, 3, true))
}`

const ch37AppendExample = `func main() {
	versions := []int{4}
	fmt.Println(append(versions, 5))
}`

const ch37DigestExample = `func main() {
	token := "session-demo"
	digest := sha256.Sum256([]byte(token))
	fmt.Println(len(token), len(digest))
}`

const ch37DigestTryExample = `func main() {
	digest := sha256.Sum256([]byte("sample"))
	fmt.Println(len(digest))
}`

const ch37MembershipExample = `type MembershipKey struct {
	Crew string
	User string
}
func main() {
	roles := map[MembershipKey]string{
		{Crew: "blue", User: "nova"}: "lead",
		{Crew: "gold", User: "nova"}: "reader",
	}
	fmt.Println(roles[MembershipKey{Crew: "blue", User: "nova"}], roles[MembershipKey{Crew: "gold", User: "nova"}])
}`

const ch37ScopeTryExample = `type Scope struct { Crew, User string }
func main() {
	roles := map[Scope]string{{"blue", "nova"}: "lead", {"gold", "nova"}: "reader"}
	fmt.Println(roles[Scope{"gold", "nova"}])
}`

const ch37CursorExample = `func main() {
	cursor := base64.RawURLEncoding.EncodeToString([]byte("row-amber"))
	fmt.Println(cursor)
}`

const ch37CursorTryExample = `func main() {
	id, _ := base64.RawURLEncoding.DecodeString("cm93LWJsdWU")
	fmt.Println(string(id))
}`

const ch37PathExample = `func main() {
	req, _ := http.NewRequest(http.MethodGet, "https://ops.example.test/v1/tasks?token=demo", nil)
	fmt.Println(req.URL.Path)
}`

const ch37ContextExample = `func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(ctx.Err())
}`

const ch37FoundationArabic = `1. كمّل migrations وrepository وسيرفر صغير في ~main.go~.
2. خلّي سجل النسخ يواكب كل تغيير في schema.
3. حافظ على عقد واحد لـSQLite وPostgres.

الـstarter فيه types والتواقيع.

ده أول جزء في capstone مركز القيادة.`

const ch37FoundationEnglish = `1. Complete the migrations, repository, and small server in ~main.go~.
2. Keep the version ledger in step with each schema change.
3. Preserve one contract for SQLite and Postgres.

The starter has the types and signatures.

This is the first part of the command-center capstone.`

var ch37FoundationArabicSteps = []LessonStep{
	{
		ID: "c371.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch37Markdown("c371.mission", ch37FoundationArabic)},
	},
	{
		ID: "c371.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c371.atomic-step",
			"الـtransaction تخلي تغييرين مرتبطين يا يثبتوا سوا يا يفضلوا زي ما كانوا. المثال model صغير، مش database.",
			ch37TransactionExample, "e371-atomic-step", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "الدالة تستقبل حالة قديمة."},
				{Line: 2, Caption: "المعاملة فشلت."},
				{Line: 3, Caption: "الحالة القديمة تفضل."},
				{Line: 5, Caption: "النجاح يضيف نسخة."},
				{Line: 8, Caption: "جهّز النسخة الحالية."},
				{Line: 9, Caption: "اعرض الحالة بعد الفشل.", Output: "[2]"},
				{Line: 10, Caption: "اعرض الحالة بعد النجاح.", Output: "[2]\n[2 3]"},
			},
		)},
	},
	{
		ID: "c371.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch37Question("c371.predict-append", ch37Prediction("توقّع سجل النسخ بعد إضافة نسخة واحدة.", ch37AppendExample), []string{"[4]", "[4 5]", "[5 4]"}, 1, []string{
				"الإضافة ما بتمسحش النسخة الموجودة.",
				"النسخة الجديدة تيجي في الآخر.",
				"append مش بتعكس ترتيب العناصر.",
			}),
			ch37Walkthrough(
				"c371.append-reveal",
				"تابع سجل صغير قبل وبعد إضافة نسخة.",
				ch37AppendExample, "e371-append-version", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 1, Caption: "ابدأ main."},
					{Line: 2, Caption: "سجل فيه نسخة واحدة."},
					{Line: 3, Caption: "append تضيف في الآخر.", Output: "[4 5]"},
				},
			),
			ch37Sandbox("c371.atomic-playground", "جرّب فشل أو نجاح، وراقب الاتنين.", "ch37-atomic-ledger", ch37AtomicConfig()),
		},
	},
	{
		ID: "c371.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch37Markdown("c371.watch", "- خلي placeholder والأنواع مقبولين من الدرايفرين.\n- صف الـledger جزء من نفس تغيير الـschema.")},
	},
	{
		ID: "c371.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch37Markdown("c371.plan", "ابدأ بقراءة النسخة، وبعدها نفّذ كل ترحيل كوحدة واحدة.\n\nافصل schema عن repository وعن HTTP؛ خلي كل دالة واضحة."),
			ch37Question("c371.ledger-transaction", "سطر النسخة يتسجل فين؟", []string{"بعد الـtransaction", "جوا نفس الـtransaction", "بعد رد HTTP"}, 1, []string{
				"الفشل هنا يسيب schema والـledger مختلفين.",
				"التغيير وسطر نسخته لازم يثبتوا مع بعض.",
				"HTTP مش مالك ترحيل schema.",
			}),
		},
	},
	{
		ID: "c371.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch37Markdown("c371.turn", "ابدأ بجدول النسخ والترحيلات. بعدهم وصّل repository والـhandlers، وراجع أن Down يعكس Up.")},
	},
}

var ch37FoundationEnglishSteps = []LessonStep{
	{
		ID: "c371.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch37Markdown("c371.mission", ch37FoundationEnglish)},
	},
	{
		ID: "c371.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c371.atomic-step",
			"A transaction makes two related changes commit together or leave the old state. This small example is a model, not a database.",
			ch37TransactionExample, "e371-atomic-step", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "The function receives old state."},
				{Line: 2, Caption: "The transaction failed."},
				{Line: 3, Caption: "Old state remains."},
				{Line: 5, Caption: "Success adds a version."},
				{Line: 8, Caption: "Prepare the current version."},
				{Line: 9, Caption: "Show state after failure.", Output: "[2]"},
				{Line: 10, Caption: "Show state after success.", Output: "[2]\n[2 3]"},
			},
		)},
	},
	{
		ID: "c371.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch37Question("c371.predict-append", ch37Prediction("Predict the version ledger after adding one version.", ch37AppendExample), []string{"[4]", "[4 5]", "[5 4]"}, 1, []string{
				"Appending does not remove the earlier version.",
				"The new version goes at the end.",
				"append does not reverse the elements.",
			}),
			ch37Walkthrough(
				"c371.append-reveal",
				"Follow a small ledger before and after adding a version.",
				ch37AppendExample, "e371-append-version", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Start main."},
					{Line: 2, Caption: "A ledger with one version."},
					{Line: 3, Caption: "append adds at the end.", Output: "[4 5]"},
				},
			),
			ch37Sandbox("c371.atomic-playground", "Try failure or commit; watch both values.", "ch37-atomic-ledger", ch37AtomicConfig()),
		},
	},
	{
		ID: "c371.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch37Markdown("c371.watch", "- Keep placeholders and types accepted by both drivers.\n- The ledger row belongs to the same schema change.")},
	},
	{
		ID: "c371.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch37Markdown("c371.plan", "Start by reading the version, then run each migration as one unit.\n\nSeparate schema, repository, and HTTP work; keep each function clear."),
			ch37Question("c371.ledger-transaction", "Where does the version row belong?", []string{"After the transaction", "Inside the same transaction", "After the HTTP reply"}, 1, []string{
				"A failure here leaves schema and ledger mismatched.",
				"The change and its version must commit together.",
				"HTTP does not own a schema migration.",
			}),
		},
	},
	{
		ID: "c371.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch37Markdown("c371.turn", "Start with the version table and migrations. Then connect the repository and handlers, and check that Down reverses Up.")},
	},
}

const ch37AuthArabic = `1. كمّل التسجيل والدخول والجلسات في ~auth.go~.
2. خزّن digest للجلسة بدل token الخام.
3. اسمح بالخروج وإلغاء الجلسات وقت تغيير كلمة السر.

الملفات السابقة جاية من مرحلة الأساس.

خلي رد بيانات الدخول الغلط موحّد.`

const ch37AuthEnglish = `1. Complete registration, sign-in, and sessions in ~auth.go~.
2. Store a session digest instead of the raw token.
3. Support logout and revoke sessions on password change.

Earlier files carry over from Foundation.

Keep incorrect sign-in replies uniform.`

var ch37AuthArabicSteps = []LessonStep{
	{
		ID: "c372.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch37Markdown("c372.mission", ch37AuthArabic)},
	},
	{
		ID: "c372.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c372.session-digest",
			"الـtoken بيرجع للعميل مرة واحدة. خزّن بصمة ثابتة منه عشان تقدر تقارنها من غير ما تحفظ سر الجلسة نفسه.",
			ch37DigestExample, "e372-session-digest", []string{"crypto/sha256", "fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "ابدأ main."},
				{Line: 2, Caption: "الـtoken نص توضيحي.", Variables: []LessonVariable{{Name: "token", Value: "session-demo"}}},
				{Line: 3, Caption: "احسب digest ثابت الطول.", Variables: []LessonVariable{{Name: "digest", Value: "32 bytes"}}},
				{Line: 4, Caption: "قارن الطولين.", Variables: []LessonVariable{{Name: "digest", Value: "32 bytes"}}, Output: "12 32"},
			},
		)},
	},
	{
		ID: "c372.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch37Question("c372.predict-digest", ch37Prediction("توقّع طول الـdigest اللي بيرجعه المثال.", ch37DigestTryExample), []string{"6", "32", "64"}, 1, []string{
				"ده طول كلمة sample، مش البصمة.",
				"SHA-256 بيرجع 32 bytes.",
				"64 هو عدد خانات hex بعد التحويل.",
			}),
			ch37Walkthrough(
				"c372.digest-reveal",
				"شوف طول بصمة SHA-256 لقيمة توضيحية.",
				ch37DigestTryExample, "e372-digest-length", []string{"crypto/sha256", "fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "ابدأ main."},
					{Line: 2, Caption: "احسب بصمة ثابتة."},
					{Line: 3, Caption: "اطبع طولها بالـbytes.", Output: "32"},
				},
			),
			ch37Sandbox("c372.session-playground", "بدّل حالة الجلسة وشوف المخزون.", "ch37-session-digest", ch37SessionConfig()),
		},
	},
	{
		ID: "c372.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch37Markdown("c372.watch", "- password hashing غير digest الجلسة؛ ما تستخدمش SHA-256 وحده لكلمات المرور.\n- رد واحد للـemail الغايب وكلمة السر الغلط يقلل كشف الحسابات.")},
	},
	{
		ID: "c372.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch37Markdown("c372.plan", "افصل إنشاء المستخدم عن التحقق من كلمة السر، وبعدها أنشئ جلسة واربطها بالمستخدم.\n\nراجع أثر logout وتغيير كلمة السر على الجلسات السابقة."),
			ch37Question("c372.session-storage", "أنهي قيمة تتخزن للجلسة؟", []string{"الـtoken الخام", "digest للـtoken", "كلمة السر"}, 1, []string{
				"الـtoken سر يحمله العميل.",
				"البصمة تكفي للمقارنة من غير تخزين السر.",
				"كلمة السر مش session credential.",
			}),
		},
	},
	{
		ID: "c372.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch37Markdown("c372.turn", "ابدأ بتدفق التسجيل والدخول، ثم جرّب إنشاء جلسة والخروج وتغيير كلمة السر على ملفات المرحلة الموروثة.")},
	},
}

var ch37AuthEnglishSteps = []LessonStep{
	{
		ID: "c372.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch37Markdown("c372.mission", ch37AuthEnglish)},
	},
	{
		ID: "c372.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c372.session-digest",
			"The token goes back to the client once. Store a stable fingerprint so you can compare it without keeping the session secret.",
			ch37DigestExample, "e372-session-digest", []string{"crypto/sha256", "fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Start main."},
				{Line: 2, Caption: "The token is illustrative.", Variables: []LessonVariable{{Name: "token", Value: "session-demo"}}},
				{Line: 3, Caption: "Compute a fixed-size digest.", Variables: []LessonVariable{{Name: "digest", Value: "32 bytes"}}},
				{Line: 4, Caption: "Compare their lengths.", Variables: []LessonVariable{{Name: "digest", Value: "32 bytes"}}, Output: "12 32"},
			},
		)},
	},
	{
		ID: "c372.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch37Question("c372.predict-digest", ch37Prediction("Predict the length of the digest returned by this example.", ch37DigestTryExample), []string{"6", "32", "64"}, 1, []string{
				"That is the length of sample, not its digest.",
				"SHA-256 returns 32 bytes.",
				"64 is the number of hex digits after encoding.",
			}),
			ch37Walkthrough(
				"c372.digest-reveal",
				"Watch the digest length for an illustrative value.",
				ch37DigestTryExample, "e372-digest-length", []string{"crypto/sha256", "fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Start main."},
					{Line: 2, Caption: "Compute a fixed-size digest."},
					{Line: 3, Caption: "Print its byte length.", Output: "32"},
				},
			),
			ch37Sandbox("c372.session-playground", "Change session state; inspect stored data.", "ch37-session-digest", ch37SessionConfig()),
		},
	},
	{
		ID: "c372.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch37Markdown("c372.watch", "- Password hashing differs from a session digest; do not use plain SHA-256 for passwords.\n- One reply for a missing email and a wrong password reveals less about accounts.")},
	},
	{
		ID: "c372.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch37Markdown("c372.plan", "Separate user creation from password verification, then create a session and bind it to the user.\n\nCheck how logout and password changes affect earlier sessions."),
			ch37Question("c372.session-storage", "Which value belongs in session storage?", []string{"The raw token", "A digest of the token", "The password"}, 1, []string{
				"The client carries the token as a secret.",
				"The digest supports comparison without storing the secret.",
				"The password is not a session credential.",
			}),
		},
	},
	{
		ID: "c372.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch37Markdown("c372.turn", "Start with registration and sign-in, then try session creation, logout, and password change across the carried files.")},
	},
}

const ch37TenantsArabic = `1. كمّل إدارة الأطقم والأدوار في ~crews.go~.
2. اعزل بيانات كل طاقم عن اللي مش عضو فيه.
3. سجّل تغييرات العضوية في audit.

الأساس وتسجيل الدخول جايين من المراحل السابقة.

فكّر في هوية المستخدم وسياق الطاقم مع بعض.`

const ch37TenantsEnglish = `1. Complete crew and role management in ~crews.go~.
2. Isolate each crew's data from nonmembers.
3. Record membership changes in the audit log.

Foundation and sign-in carry over from earlier stages.

Consider user identity together with crew scope.`

var ch37TenantsArabicSteps = []LessonStep{
	{
		ID: "c373.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch37Markdown("c373.mission", ch37TenantsArabic)},
	},
	{
		ID: "c373.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c373.tenant-key",
			"العضوية مرتبطة بزوج: الطاقم والمستخدم. المستخدم نفسه ممكن يبقى له دور مختلف في طاقم تاني.",
			ch37MembershipExample, "e373-tenant-key", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "المفتاح فيه طاقم ومستخدم."},
				{Line: 5, Caption: "ابدأ main."},
				{Line: 6, Caption: "خريطة العضويات."},
				{Line: 10, Caption: "اقرأ الدورين.", Output: "lead reader"},
			},
		)},
	},
	{
		ID: "c373.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch37Question("c373.predict-scope", ch37Prediction("توقّع دور nova داخل الطاقم gold.", ch37ScopeTryExample), []string{"lead", "reader", "empty"}, 1, []string{
				"lead يخص الطاقم التاني.",
				"مفتاح المجموعة فيه gold وnova.",
				"العضوية موجودة في الخريطة.",
			}),
			ch37Walkthrough(
				"c373.scope-reveal",
				"تابع lookup للمستخدم نفسه في طاقمين.",
				ch37ScopeTryExample, "e373-scope-reveal", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التنفيذ."},
					{Line: 1, Caption: "Scope يجمع قيمتين."},
					{Line: 3, Caption: "خريطة فيها عضويتان."},
					{Line: 4, Caption: "اطلب المفتاح المركب.", Output: "reader"},
				},
			),
			ch37Sandbox("c373.scope-playground", "بدّل نطاق الطاقم وشوف الدور.", "ch37-scoped-membership", ch37CrewConfig()),
		},
	},
	{
		ID: "c373.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch37Markdown("c373.watch", "- الغريب ياخد 404 حتى لو الطاقم موجود.\n- audit والتغيير لازم يفضلوا متسقين لو خطوة فشلت.")},
	},
	{
		ID: "c373.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch37Markdown("c373.plan", "ابدأ بهوية المنادي، بعدين عضويته ودوره في الطاقم المطلوب.\n\nبعد التحقق، اربط تغيير العضوية بسطر audit وراجع رد الـHTTP."),
			ch37Question("c373.outsider-code", "غريب سأل عن طاقم؛ أنهي رد يخفي وجوده؟", []string{"403 Forbidden", "404 Not Found", "200 OK"}, 1, []string{
				"403 يلمّح إن المورد موجود.",
				"404 مايكشفش وجود الطاقم.",
				"200 يكشف بيانات المورد.",
			}),
		},
	},
	{
		ID: "c373.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch37Markdown("c373.turn", "ابدأ بـRole وCreateCrew، وبعدها امشِ في إضافة عضو وقراءة audit من خلال الـAPI.")},
	},
}

var ch37TenantsEnglishSteps = []LessonStep{
	{
		ID: "c373.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch37Markdown("c373.mission", ch37TenantsEnglish)},
	},
	{
		ID: "c373.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c373.tenant-key",
			"Membership belongs to a pair: crew and user. One user can have a different role in another crew.",
			ch37MembershipExample, "e373-tenant-key", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "The key has crew and user."},
				{Line: 5, Caption: "Start main."},
				{Line: 6, Caption: "The membership map."},
				{Line: 10, Caption: "Read both roles.", Output: "lead reader"},
			},
		)},
	},
	{
		ID: "c373.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch37Question("c373.predict-scope", ch37Prediction("Predict nova's role in the gold crew.", ch37ScopeTryExample), []string{"lead", "reader", "empty"}, 1, []string{
				"lead belongs to the other crew.",
				"The lookup key includes gold and nova.",
				"The membership is present in the map.",
			}),
			ch37Walkthrough(
				"c373.scope-reveal",
				"Follow a lookup for one user in two crews.",
				ch37ScopeTryExample, "e373-scope-reveal", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Scope groups two values."},
					{Line: 3, Caption: "The map has two memberships."},
					{Line: 4, Caption: "Look up the combined key.", Output: "reader"},
				},
			),
			ch37Sandbox("c373.scope-playground", "Change crew scope; inspect the role.", "ch37-scoped-membership", ch37CrewConfig()),
		},
	},
	{
		ID: "c373.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch37Markdown("c373.watch", "- An outsider gets 404 even when the crew exists.\n- Keep the audit and change consistent if one step fails.")},
	},
	{
		ID: "c373.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch37Markdown("c373.plan", "Start with the caller's identity, then membership and role in the requested crew.\n\nAfter authorization, connect the membership change to its audit entry and review the HTTP reply."),
			ch37Question("c373.outsider-code", "Which reply hides whether a crew exists from an outsider?", []string{"403 Forbidden", "404 Not Found", "200 OK"}, 1, []string{
				"403 hints that the resource exists.",
				"404 does not reveal whether the crew exists.",
				"200 exposes resource data.",
			}),
		},
	},
	{
		ID: "c373.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch37Markdown("c373.turn", "Start with Role and CreateCrew, then follow member addition and audit reading through the API.")},
	},
}

const ch37TasksArabic = `1. كمّل إنشاء وتعديل وعرض المهام في ~tasks.go~.
2. عالج إعادة الطلب وتعارض النسخ.
3. خلّي صفحات القائمة تكمل من cursor ثابت.

الملفات السابقة محمولة من مرحلة الأطقم.

ركّز على سلوك كل endpoint لوحده.`

const ch37TasksEnglish = `1. Complete task creation, updates, and listing in ~tasks.go~.
2. Handle request retries and version conflicts.
3. Continue list pages from a stable cursor.

Earlier files carry over from the Crews stage.

Focus on each endpoint's behavior.`

var ch37TasksArabicSteps = []LessonStep{
	{
		ID: "c374.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch37Markdown("c374.mission", ch37TasksArabic)},
	},
	{
		ID: "c374.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c374.cursor-encode",
			"الـcursor علامة متابعة opaque للعميل. ترميز الـID يخليه نص آمن للنقل، والترتيب الثابت يحدد الصفحة التالية.",
			ch37CursorExample, "e374-cursor-encode", []string{"encoding/base64", "fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "ابدأ main."},
				{Line: 2, Caption: "حوّل ID توضيحي لنص."},
				{Line: 3, Caption: "اطبع الـcursor.", Output: "cm93LWFtYmVy"},
			},
		)},
	},
	{
		ID: "c374.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch37Question("c374.predict-cursor", ch37Prediction("توقّع الـID بعد فك الـcursor.", ch37CursorTryExample), []string{"cm93LWJsdWU", "row-blue", "row-amber"}, 1, []string{
				"ده النص المرمّز نفسه.",
				"فك الترميز يرجّع ID الأصلي.",
				"ده مثال تاني مش قيمة الـcursor.",
			}),
			ch37Walkthrough(
				"c374.cursor-reveal",
				"فك cursor توضيحي عشان تشوف الـID.",
				ch37CursorTryExample, "e374-cursor-decode", []string{"encoding/base64", "fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "ابدأ main."},
					{Line: 2, Caption: "فك النص المرمّز."},
					{Line: 3, Caption: "اطبع الـID الأصلي.", Output: "row-blue"},
				},
			),
			ch37Sandbox("c374.cursor-playground", "حرّك الصفحة وشوف ~cursor~.", "ch37-cursor-pages", ch37CursorConfig()),
		},
	},
	{
		ID: "c374.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch37Markdown("c374.watch", "- مفتاح التكرار مرتبط بالمستخدم والطاقم والجسم.\n- version تمنع تعديل قديم يمسح تحديث أحدث؛ OFFSET بيتأثر بإضافات الصفوف.")},
	},
	{
		ID: "c374.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch37Markdown("c374.plan", "راجع POST لوحده مع إعادة الطلب، ثم PATCH مع نسخة قديمة، وبعدها مرّ على صفحات GET.\n\nقارن ترتيب العناصر وحدد إزاي كل cursor بيكمل القائمة."),
			ch37Question("c374.retry-check", "نفس المفتاح ونفس الجسم اتبعتوا تاني؛ إيه المتوقع؟", []string{"مهمة جديدة", "الرد المحفوظ", "تعديل للنسخة"}, 1, []string{
				"إعادة الإرسال المفروض ما تكررّش الأثر.",
				"نفس الطلب يرجع الرد اللي اتحفظ.",
				"النسخة تخص تعديل PATCH.",
			}),
		},
	},
	{
		ID: "c374.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch37Markdown("c374.turn", "ابدأ بعملية واحدة في ~tasks.go~ وشغّل contract tests بتاعتها قبل ما تنتقل للمسار اللي بعده.")},
	},
}

var ch37TasksEnglishSteps = []LessonStep{
	{
		ID: "c374.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch37Markdown("c374.mission", ch37TasksEnglish)},
	},
	{
		ID: "c374.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c374.cursor-encode",
			"A cursor is an opaque continuation marker for the client. Encoding an ID makes it safe to carry as text, while stable ordering identifies the next page.",
			ch37CursorExample, "e374-cursor-encode", []string{"encoding/base64", "fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Start main."},
				{Line: 2, Caption: "Encode an illustrative ID."},
				{Line: 3, Caption: "Print the cursor.", Output: "cm93LWFtYmVy"},
			},
		)},
	},
	{
		ID: "c374.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch37Question("c374.predict-cursor", ch37Prediction("Predict the ID after decoding the cursor.", ch37CursorTryExample), []string{"cm93LWJsdWU", "row-blue", "row-amber"}, 1, []string{
				"That is the encoded text itself.",
				"Decoding returns the original ID.",
				"That is a different example value.",
			}),
			ch37Walkthrough(
				"c374.cursor-reveal",
				"Decode an illustrative cursor to see its ID.",
				ch37CursorTryExample, "e374-cursor-decode", []string{"encoding/base64", "fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Start main."},
					{Line: 2, Caption: "Decode the encoded text."},
					{Line: 3, Caption: "Print the original ID.", Output: "row-blue"},
				},
			),
			ch37Sandbox("c374.cursor-playground", "Advance a page and inspect ~cursor~.", "ch37-cursor-pages", ch37CursorConfig()),
		},
	},
	{
		ID: "c374.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch37Markdown("c374.watch", "- An idempotency key belongs to the user, crew, and request body.\n- A version prevents an old edit replacing a newer one; OFFSET shifts when rows are added.")},
	},
	{
		ID: "c374.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch37Markdown("c374.plan", "Review POST on its own with a retry, then PATCH with an old version, then walk through GET pages.\n\nCompare item order and decide how each cursor continues the list."),
			ch37Question("c374.retry-check", "The same key and body arrive again; what should happen?", []string{"Create another task", "Return the saved reply", "Change the version"}, 1, []string{
				"A retry should not repeat its effect.",
				"The same request returns its saved reply.",
				"Versions belong to PATCH updates.",
			}),
		},
	},
	{
		ID: "c374.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch37Markdown("c374.turn", "Start with one operation in ~tasks.go~ and run its contract checks before moving to the next path.")},
	},
}

const ch37OpsArabic = `1. كمّل middleware و~Serve~ في ~ops.go~.
2. ضيف request ID وأخطاء وmetrics وحد طلبات من غير أسرار في السجل.
3. صفّي الطلبات الجارية وقت الإغلاق.

كل ملفات الـAPI السابقة محمولة.

المهمة عن جمعها بأمان تحت حمل متوازي.`

const ch37OpsEnglish = `1. Complete middleware and ~Serve~ in ~ops.go~.
2. Add request IDs, errors, metrics, rate limits, and secret-free logs.
3. Drain active requests during shutdown.

All earlier API files carry over.

The task is to integrate them safely under concurrent load.`

var ch37OpsArabicSteps = []LessonStep{
	{
		ID: "c375.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch37Markdown("c375.mission", ch37OpsArabic)},
	},
	{
		ID: "c375.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c375.path-log",
			"~URL.Path~ فيه المسار من غير query. استخدمه لما تسجل الطلب عشان قيمة زي token ماتدخلش في السجل.",
			ch37PathExample, "e375-path-log", []string{"fmt", "net/http"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "ابدأ main."},
				{Line: 2, Caption: "الطلب فيه query."},
				{Line: 3, Caption: "اطبع المسار وحده.", Output: "/v1/tasks"},
			},
		)},
	},
	{
		ID: "c375.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch37Question("c375.predict-context", ch37Prediction("توقّع نتيجة الـcontext بعد cancel.", ch37ContextExample), []string{"<nil>", "context canceled", "deadline exceeded"}, 1, []string{
				"الـcontext اتلغى بالفعل.",
				"Err بيرجع سبب الإلغاء.",
				"مافيش deadline اتحدد.",
			}),
			ch37Walkthrough(
				"c375.context-reveal",
				"تابع الـcontext بعد إلغاء الأب.",
				ch37ContextExample, "e375-canceled-context", []string{"context", "fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "ابدأ main."},
					{Line: 2, Caption: "أنشئ context قابل للإلغاء."},
					{Line: 3, Caption: "ألغِه صراحة."},
					{Line: 4, Caption: "اقرأ سبب الإلغاء.", Output: "context canceled"},
				},
			),
			ch37Sandbox("c375.request-playground", "تتبّع طلبًا تجريبيًا مرحلة مرحلة.", "ch37-request-path", ch37OperationsConfig()),
		},
	},
	{
		ID: "c375.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch37Markdown("c375.watch", "- استخدم context جديد بمهلة للتصفية بعد إلغاء context السيرفر.\n- احمِ العدادات والنوافذ والسجل من race؛ ما تسجلش headers.")},
	},
	{
		ID: "c375.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch37Markdown("c375.plan", "تتبّع الطلب من ID للـhandler ثم للـresponse والـlog، مع قياس الحالة مرة واحدة.\n\nافصل طلب metrics عن العد، وراجع سلوك drain بعد الإلغاء."),
			ch37Question("c375.drain-context", "بعد إلغاء context الأصلي، إيه اللي يسمح بوقت للتصفية؟", []string{"نفس context الملغي", "context جديد", "Background بلا مهلة"}, 1, []string{
				"السياق الملغي ينهي العمل فوراً.",
				"سياق جديد يقدر يحمل مهلة التصفية.",
				"من غير مهلة مفيش حد للانتظار.",
			}),
		},
	},
	{
		ID: "c375.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch37Markdown("c375.turn", "ابدأ بالـmiddleware على handler بسيط. بعده جرّب الخطأ والحد والتزامن، ثم شغّل Serve على listener محلي.")},
	},
}

var ch37OpsEnglishSteps = []LessonStep{
	{
		ID: "c375.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch37Markdown("c375.mission", ch37OpsEnglish)},
	},
	{
		ID: "c375.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch37Walkthrough(
			"c375.path-log",
			"~URL.Path~ contains the path without its query. Use it when logging a request so a value such as a token stays out of the log.",
			ch37PathExample, "e375-path-log", []string{"fmt", "net/http"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Start main."},
				{Line: 2, Caption: "The request has a query."},
				{Line: 3, Caption: "Print only the path.", Output: "/v1/tasks"},
			},
		)},
	},
	{
		ID: "c375.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch37Question("c375.predict-context", ch37Prediction("Predict the context result after cancel.", ch37ContextExample), []string{"<nil>", "context canceled", "deadline exceeded"}, 1, []string{
				"The context has already been canceled.",
				"Err reports why it ended.",
				"No deadline was set.",
			}),
			ch37Walkthrough(
				"c375.context-reveal",
				"Follow the context after its parent is canceled.",
				ch37ContextExample, "e375-canceled-context", []string{"context", "fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Start main."},
					{Line: 2, Caption: "Create a cancelable context."},
					{Line: 3, Caption: "Cancel it explicitly."},
					{Line: 4, Caption: "Read the cancellation reason.", Output: "context canceled"},
				},
			),
			ch37Sandbox("c375.request-playground", "Trace a sample request through its stages.", "ch37-request-path", ch37OperationsConfig()),
		},
	},
	{
		ID: "c375.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch37Markdown("c375.watch", "- Use a fresh timed context to drain after the server context is canceled.\n- Protect counters, windows, and logs from races; do not log headers.")},
	},
	{
		ID: "c375.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch37Markdown("c375.plan", "Trace a request from ID to handler, response, and log while counting its status once.\n\nKeep metrics requests out of the count, and check drain behavior after cancellation."),
			ch37Question("c375.drain-context", "After the original context is canceled, what allows time to drain?", []string{"The same canceled context", "A fresh context", "Background with no deadline"}, 1, []string{
				"A canceled context ends work immediately.",
				"A fresh context can carry the drain timeout.",
				"Without a timeout there is no wait limit.",
			}),
		},
	},
	{
		ID: "c375.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch37Markdown("c375.turn", "Start with middleware around a simple handler. Then try errors, limits, and concurrency, and run Serve on a local listener.")},
	},
}
