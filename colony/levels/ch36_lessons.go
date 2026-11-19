package levels

type ch36SandboxPair struct {
	Arabic  string `json:"ar"`
	English string `json:"en"`
}

type ch36SandboxRow struct {
	Label ch36SandboxPair `json:"label"`
	Value ch36SandboxPair `json:"value"`
}

type ch36SandboxFrame struct {
	Caption ch36SandboxPair  `json:"caption"`
	Rows    []ch36SandboxRow `json:"rows"`
}

func ch36Pair(arabic, english string) ch36SandboxPair {
	return ch36SandboxPair{Arabic: arabic, English: english}
}

func ch36LanguagePair(pair ch36SandboxPair, language int) ch36SandboxPair {
	if language == 0 {
		pair.English = ""
	} else {
		pair.Arabic = ""
	}
	return pair
}

func ch36Row(labelArabic, labelEnglish, valueArabic, valueEnglish string) ch36SandboxRow {
	return ch36SandboxRow{Label: ch36Pair(labelArabic, labelEnglish), Value: ch36Pair(valueArabic, valueEnglish)}
}

func ch36Frame(captionArabic, captionEnglish string, rows ...ch36SandboxRow) ch36SandboxFrame {
	return ch36SandboxFrame{Caption: ch36Pair(captionArabic, captionEnglish), Rows: rows}
}

func ch36SandboxInstruction(prefix string, language int) string {
	text := map[string][2]string{
		"c361.": {"سجّل قراءة وراقب counter وgauge.", "Record a reading and compare counter and gauge."},
		"c362.": {"بدّل runtime وشغّل نفس method.", "Swap the runtime and call the same method."},
		"c363.": {"بدّل حالة DB واطلب probes.", "Change database state and request both probes."},
		"c364.": {"جرّب فشل delivery وبعده retry.", "Try a failed delivery, then retry."},
	}[prefix]
	return text[language]
}

func ch36SandboxConfig(prefix string, language int) map[string]any {
	advance := ch36Pair("كمّل خطوة", "Next step")
	replay := ch36Pair("جرّب الحالة التانية", "Try the other case")
	reset := ch36Pair("من الأول", "Reset")
	variants := [][]ch36SandboxFrame{}
	switch prefix {
	case "c361.":
		variants = [][]ch36SandboxFrame{
			{
				ch36Frame("الـcounter بيجمع القراءات.", "The counter accumulates readings.", ch36Row("label", "Label", "bay-east", "bay-east"), ch36Row("counter", "Counter", "4", "4"), ch36Row("gauge", "Gauge", "18", "18")),
				ch36Frame("قراءة جديدة تزود counter.", "A new reading raises the counter.", ch36Row("label", "Label", "bay-east", "bay-east"), ch36Row("counter", "Counter", "5", "5"), ch36Row("gauge", "Gauge", "18", "18")),
			},
			{
				ch36Frame("labels مختلفة تتجمع في other.", "Different labels group under other.", ch36Row("label", "Label", "other", "other"), ch36Row("counter", "Counter", "5", "5"), ch36Row("gauge", "Gauge", "18", "18")),
				ch36Frame("الـgauge ممكن يقل.", "A gauge can move down.", ch36Row("label", "Label", "other", "other"), ch36Row("counter", "Counter", "5", "5"), ch36Row("gauge", "Gauge", "16", "16")),
			},
		}
	case "c362.":
		variants = [][]ch36SandboxFrame{
			{
				ch36Frame("الـCLI ينادي contract مش engine.", "The CLI calls a contract, not an engine.", ch36Row("caller", "Caller", "RunPodCtl", "RunPodCtl"), ch36Row("interface", "Interface", "Runtime", "Runtime"), ch36Row("implementation", "Implementation", "fake-a", "fake-a")),
				ch36Frame("نفس method وصلت للـfake.", "The same method reached the fake.", ch36Row("method", "Method", "Start(probe-4)", "Start(probe-4)"), ch36Row("implementation", "Implementation", "fake-a", "fake-a"), ch36Row("result", "Result", "started", "started")),
			},
			{
				ch36Frame("الـinterface تقبل implementation تانية.", "The interface accepts another implementation.", ch36Row("caller", "Caller", "RunPodCtl", "RunPodCtl"), ch36Row("interface", "Interface", "Runtime", "Runtime"), ch36Row("implementation", "Implementation", "fake-b", "fake-b")),
				ch36Frame("النداء نفسه؛ التنفيذ اتبدّل.", "Same call, different implementation.", ch36Row("method", "Method", "Start(probe-4)", "Start(probe-4)"), ch36Row("implementation", "Implementation", "fake-b", "fake-b"), ch36Row("result", "Result", "probe ready", "probe ready")),
			},
		}
	case "c363.":
		variants = [][]ch36SandboxFrame{
			{
				ch36Frame("قاعدة البيانات جاهزة.", "The database is ready.", ch36Row("DB", "Database", "connected", "connected"), ch36Row("/livez", "/livez", "ok", "ok"), ch36Row("/readyz", "/readyz", "ready", "ready")),
				ch36Frame("ترتيب الـrequests ممكن يختلف.", "Request order can vary.", ch36Row("أول response", "First response", "/livez: ok", "/livez: ok"), ch36Row("تاني response", "Second response", "/readyz: ready", "/readyz: ready"), ch36Row("الترتيب", "Order", "live ثم ready", "live then ready")),
			},
			{
				ch36Frame("قاعدة البيانات واقفة.", "The database is offline.", ch36Row("DB", "Database", "offline", "offline"), ch36Row("/livez", "/livez", "ok", "ok"), ch36Row("/readyz", "/readyz", "waiting", "waiting")),
				ch36Frame("الـliveness مستقل عن DB.", "Liveness stays independent of the database.", ch36Row("أول response", "First response", "/readyz: waiting", "/readyz: waiting"), ch36Row("تاني response", "Second response", "/livez: ok", "/livez: ok"), ch36Row("الترتيب", "Order", "ready ثم live", "ready then live")),
			},
		}
	case "c364.":
		variants = [][]ch36SandboxFrame{
			{
				ch36Frame("الـrows اتقفلت قبل broker.", "Rows closed before the broker call.", ch36Row("DB rows", "DB rows", "closed", "closed"), ch36Row("محاولة A", "Attempt A", "failed", "failed"), ch36Row("outbox", "Outbox", "pending", "pending")),
				ch36Frame("retry نجح والحدث اتسجل.", "Retry succeeded and delivery is recorded.", ch36Row("DB rows", "DB rows", "closed", "closed"), ch36Row("محاولة B", "Attempt B", "sent", "sent"), ch36Row("outbox", "Outbox", "delivered", "delivered")),
			},
			{
				ch36Frame("محاولة B وصلت الأول.", "Attempt B arrived first.", ch36Row("DB rows", "DB rows", "closed", "closed"), ch36Row("ترتيب المحاولات", "Attempt order", "B ثم A", "B then A"), ch36Row("outbox", "Outbox", "pending", "pending")),
				ch36Frame("الترتيب اختلف، والحالة محفوظة.", "Order changed; durable state remains.", ch36Row("DB rows", "DB rows", "closed", "closed"), ch36Row("المحاولة الناجحة", "Successful attempt", "A", "A"), ch36Row("outbox", "Outbox", "delivered", "delivered")),
			},
		}
	}
	labels := map[string]ch36SandboxPair{"advance": advance, "again": replay, "reset": reset}
	for key, pair := range labels {
		labels[key] = ch36LanguagePair(pair, language)
	}
	for variant := range variants {
		for frame := range variants[variant] {
			variants[variant][frame].Caption = ch36LanguagePair(variants[variant][frame].Caption, language)
			for row := range variants[variant][frame].Rows {
				variants[variant][frame].Rows[row].Label = ch36LanguagePair(variants[variant][frame].Rows[row].Label, language)
				variants[variant][frame].Rows[row].Value = ch36LanguagePair(variants[variant][frame].Rows[row].Value, language)
			}
		}
	}
	return map[string]any{
		"labels":   labels,
		"variants": variants,
	}
}

func ch36Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch36Sandbox(id, text, prefix string, language int) LessonBlock {
	return LessonBlock{ID: id, Kind: "sandbox", Text: text, Sandbox: &LessonSandbox{Kind: "ch36-service-flow", Config: ch36SandboxConfig(prefix, language)}}
}

func ch36Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch36Prediction(text, code string) string {
	return text + "\n\n~~~go\n" + code + "\n~~~"
}

func ch36Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

const ch36DurationExample = `func main() {
	d := 1500 * time.Millisecond
	fmt.Printf("%.1f\n", d.Seconds())
}`

const ch36LabelExample = `func bayLabel(bay string) string {
	if bay == "amber" {
		return bay
	}
	return "other"
}
func main() {
	fmt.Println(bayLabel("amber"), bayLabel("flight-9"))
}`

const ch36InterfaceExample = `type Starter interface {
	Start(string) string
}
type Demo struct{}
func (Demo) Start(name string) string { return "started " + name }
func main() {
	var engine Starter = Demo{}
	fmt.Println(engine.Start("probe-4"))
}`

const ch36MetricsArabic = `1. كمّل ~DockMetrics~ في ~main.go~ واربط المقاييس بالـregistry المحلي.
2. فرّق بين counter وGauge وHistogram.
3. خلّي values الـlabels محدودة والتحديثات آمنة.

الـstarter فيه الـtypes والتواقيع.

ده lab؛ العقد هو مرجعك.`

const ch36MetricsEnglish = `1. Complete ~DockMetrics~ in ~main.go~ and use the local registry.
2. Distinguish a counter, gauge, and histogram.
3. Keep label values bounded and updates safe.

The starter has the types and signatures.

This is a lab; use its contract as your guide.`

var ch36MetricsArabicSteps = []LessonStep{
	{
		ID: "c361.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch36Markdown("c361.mission", ch36MetricsArabic)},
	},
	{
		ID: "c361.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c361.duration",
			"الـduration مش رقم ثواني تلقائي. ~Seconds~ بتحولها لرقم بالثواني قبل عرضها كمقياس. المثال توضيحي.",
			ch36DurationExample, "e361-duration", []string{"fmt", "time"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التشغيل."},
				{Line: 1, Caption: "ابدأ تشغيل المثال."},
				{Line: 2, Caption: "مدة توضيحية بالملي ثانية.", Variables: []LessonVariable{{Name: "d", Value: "1.5s"}}},
				{Line: 3, Caption: "الناتج بالثواني.", Variables: []LessonVariable{{Name: "d", Value: "1.5s"}}, Output: "1.5"},
			},
		)},
	},
	{
		ID: "c361.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch36Question("c361.predict-label", ch36Prediction("توقّع الناتج؛ المدخل الغريب بيروح لفئة ثابتة.", ch36LabelExample), []string{"amber flight-9", "amber other", "other other"}, 1, []string{
				"الـvalue الغريبة مش label مستقلة.",
				"الغريب بيتجمع في other.",
				"الـvalue المعروفة بتفضل زي ما هي.",
			}),
			ch36Walkthrough(
				"c361.label-reveal",
				"تابع label مع value معروفة وبعدها value جديدة.",
				ch36LabelExample, "e361-bounded-label", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "الـfunction بتاخد bay."},
					{Line: 2, Caption: "قارن بالـvalue المعروفة."},
					{Line: 5, Caption: "الـvalues التانية تتجمع."},
					{Line: 7, Caption: "اطبع الفئتين."},
					{Line: 8, Caption: "الفئة ثابتة للغريب.", Output: "amber other"},
				},
			),
			ch36Sandbox("c361.playground", ch36SandboxInstruction("c361.", 0), "c361.", 0),
		},
	},
	{
		ID: "c361.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch36Markdown("c361.watch", "- اسم سفينة أو ID عشوائي كـlabel بيكبر عدد السلاسل.\n- الـcounter تراكمي؛ الـgauge ينفع يزيد وينقص.")},
	},
	{
		ID: "c361.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch36Markdown("c361.plan", "اكتب معنى كل value ووحدتها قبل اختيار النوع.\n\nبعدها اربط كل method بالمقياس والـhandler، وراجع errors التسجيل."),
			ch36Question("c361.gauge-check", "طول الـqueue الحالي أقرب لأي نوع؟", []string{"Counter", "Gauge", "Histogram"}, 1, []string{
				"العداد بيسجل الزيادات المتراكمة.",
				"الـGauge بيمثل value بتطلع وتنزل.",
				"المدرّج بيجمع قياسات في buckets.",
			}),
		},
	},
	{
		ID: "c361.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch36Markdown("c361.turn", "ابدأ بالتواقيع والـregistry في starter. شغّل go test بعد كل جزء، وراجع القراءات من غير default registry.")},
	},
}

var ch36MetricsEnglishSteps = []LessonStep{
	{
		ID: "c361.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch36Markdown("c361.mission", ch36MetricsEnglish)},
	},
	{
		ID: "c361.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c361.duration",
			"A duration is not automatically a number of seconds. ~Seconds~ converts it to seconds before a metric uses it. This example is illustrative.",
			ch36DurationExample, "e361-duration", []string{"fmt", "time"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Start the example."},
				{Line: 2, Caption: "An illustrative millisecond duration.", Variables: []LessonVariable{{Name: "d", Value: "1.5s"}}},
				{Line: 3, Caption: "The result is seconds.", Variables: []LessonVariable{{Name: "d", Value: "1.5s"}}, Output: "1.5"},
			},
		)},
	},
	{
		ID: "c361.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch36Question("c361.predict-label", ch36Prediction("Predict the output; an unknown input goes into one fixed group.", ch36LabelExample), []string{"amber flight-9", "amber other", "other other"}, 1, []string{
				"The unknown value is not its own label.",
				"The unknown value joins other.",
				"The known value stays unchanged.",
			}),
			ch36Walkthrough(
				"c361.label-reveal",
				"Follow the label for a known value, then a new one.",
				ch36LabelExample, "e361-bounded-label", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "The function takes a bay."},
					{Line: 2, Caption: "Compare with the known value."},
					{Line: 5, Caption: "Other values join together."},
					{Line: 7, Caption: "Print both groups."},
					{Line: 8, Caption: "Unknown values share one group.", Output: "amber other"},
				},
			),
			ch36Sandbox("c361.playground", ch36SandboxInstruction("c361.", 1), "c361.", 1),
		},
	},
	{
		ID: "c361.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch36Markdown("c361.watch", "- A ship name or random ID as a label grows the number of series.\n- A counter accumulates; a gauge can go up or down.")},
	},
	{
		ID: "c361.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch36Markdown("c361.plan", "Write down what each value means and its unit before choosing a type.\n\nThen connect each method to its metric and handler, and review registration errors."),
			ch36Question("c361.gauge-check", "Which type best represents the current queue length?", []string{"Counter", "Gauge", "Histogram"}, 1, []string{
				"A counter records accumulating increases.",
				"A gauge represents a value that can move either way.",
				"A histogram groups measured values into buckets.",
			}),
		},
	},
	{
		ID: "c361.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch36Markdown("c361.turn", "Start with the starter's signatures and registry. Run go test after each part, and inspect the readings without the default registry.")},
	},
}

const ch36InterfaceExampleReveal = `type Name interface { Label() string }
type Relay struct{}
func (Relay) Label() string { return "relay" }
func main() {
	var item Name = Relay{}
	fmt.Println(item.Label())
}`

const ch36ReadinessLogExample = `func main() {
	opts := &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		if strings.EqualFold(a.Key, "token") {
			return slog.String(a.Key, "[REDACTED]")
		}
		return a
	}}
	log := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	log.Info("request", "token", "sample", "status", 503)
}`

const ch36PathExample = `func main() {
	req, _ := http.NewRequest(http.MethodGet, "https://ops.example.test/readyz?token=demo", nil)
	fmt.Println(req.URL.Path)
}`

const ch36ReceiptExample = `type Receipt struct {
	ID        string
	Published bool
}
func main() {
	receipt := Receipt{ID: "storm-demo"}
	fmt.Println(receipt.ID, receipt.Published)
	receipt.Published = true
	fmt.Println(receipt.Published)
}`

const ch36AckExample = `func acknowledged(publishOK bool) bool {
	if !publishOK {
		return false
	}
	return true
}
func main() {
	fmt.Println(acknowledged(false), acknowledged(true))
}`

const ch36OpsArabic = `1. كمّل logger وفحوص ~/livez~ و ~/readyz~ في ~main.go~.
2. اخفِ الـsecret keys من السجل.
3. خلي فحص الجاهزية مرتبط بقاعدة البيانات.

الـstarter فيه ~Pinger~ وmethod ~ServeHTTP~.

ده lab؛ ابدأ من العقد والتواقيع.`

const ch36OpsEnglish = `1. Complete the logger and ~/livez~ and ~/readyz~ probes in ~main.go~.
2. Hide secret keys from logs.
3. Tie readiness to the database.

The starter includes ~Pinger~ and a ~ServeHTTP~ surface.

This is a lab; start from its contract and signatures.`

var ch36OpsArabicSteps = []LessonStep{
	{
		ID: "c363.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch36Markdown("c363.mission", ch36OpsArabic)},
	},
	{
		ID: "c363.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c363.structured-log",
			"~slog~ بيطلع keys وvalues منظمة. ~ReplaceAttr~ تقدر تشيل key أو تبدّل الـvalue، حتى جوه group.",
			ch36ReadinessLogExample, "e363-structured-log", []string{"log/slog", "os", "strings"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "ابدأ main."},
				{Line: 2, Caption: "جهّز تعديل السجل."},
				{Line: 3, Caption: "احذف الوقت المتغير."},
				{Line: 6, Caption: "فتّش عن الـkey الحساس."},
				{Line: 7, Caption: "بدّل الـvalue السرية."},
				{Line: 11, Caption: "جهّز JSON handler."},
				{Line: 12, Caption: "اكتب record منظم.", Output: "{\"level\":\"INFO\",\"msg\":\"request\",\"token\":\"[REDACTED]\",\"status\":503}"},
			},
		)},
	},
	{
		ID: "c363.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch36Question("c363.predict-path", ch36Prediction("توقّع المسار اللي يفضل بعد تجاهل query.", ch36PathExample), []string{"/readyz?token=demo", "/readyz", "https://ops.example.test/readyz"}, 1, []string{
				"ده الـURL كامل ومعاه query.",
				"~Path~ فيه المسار بس.",
				"ده فيه scheme وhost كمان.",
			}),
			ch36Walkthrough(
				"c363.path-reveal",
				"راقب الفرق بين مسار الـrequest والـquery.",
				ch36PathExample, "e363-path-only", []string{"fmt", "net/http"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "ابدأ main."},
					{Line: 2, Caption: "جهّز request فيه query."},
					{Line: 3, Caption: "اطبع المسار وحده.", Output: "/readyz"},
				},
			),
			ch36Sandbox("c363.playground", ch36SandboxInstruction("c363.", 0), "c363.", 0),
		},
	},
	{
		ID: "c363.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch36Markdown("c363.watch", "- ~/livez~ يجاوب عن العملية نفسها؛ مش لازم يعمل Ping.\n- ~URL.Path~ من غير query؛ ده يقلل فرصة تسجيل token.")},
	},
	{
		ID: "c363.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch36Markdown("c363.plan", "جهّز logger الأول، وبعده فرّق بين فروع HTTP وسجّل النتيجة بعد تحديدها.\n\nاقرا الـrequests المتوازية بعين على الحالة المشتركة."),
			ch36Question("c363.live-check", "لو قاعدة البيانات وقعت، فحص الحياة يعمل إيه؟", []string{"يفحص القاعدة", "يفضل مستقل عنها", "يرجع نجاح الجاهزية"}, 1, []string{
				"فحص الجاهزية هو اللي يتابع dependency.",
				"حياة العملية مش رهينة قاعدة البيانات.",
				"الحياة والجاهزية سؤالان مختلفان.",
			}),
		},
	},
	{
		ID: "c363.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch36Markdown("c363.turn", "نفّذ logger لوحده الأول. بعده جرّب كل probe، وافحص JSON الناتج وسلوك التنفيذ المتوازي.")},
	},
}

var ch36OpsEnglishSteps = []LessonStep{
	{
		ID: "c363.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch36Markdown("c363.mission", ch36OpsEnglish)},
	},
	{
		ID: "c363.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c363.structured-log",
			"~slog~ writes structured keys and values. ~ReplaceAttr~ can remove a key or change its value, including inside a group.",
			ch36ReadinessLogExample, "e363-structured-log", []string{"log/slog", "os", "strings"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "Start main."},
				{Line: 2, Caption: "Prepare log changes."},
				{Line: 3, Caption: "Remove the changing time."},
				{Line: 6, Caption: "Check for a sensitive key."},
				{Line: 7, Caption: "Replace the secret value."},
				{Line: 11, Caption: "Create a JSON handler."},
				{Line: 12, Caption: "Write a structured record.", Output: "{\"level\":\"INFO\",\"msg\":\"request\",\"token\":\"[REDACTED]\",\"status\":503}"},
			},
		)},
	},
	{
		ID: "c363.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch36Question("c363.predict-path", ch36Prediction("Predict the path that remains after ignoring the query.", ch36PathExample), []string{"/readyz?token=demo", "/readyz", "https://ops.example.test/readyz"}, 1, []string{
				"That is the full URL with its query.",
				"~Path~ contains only the path.",
				"That also includes the scheme and host.",
			}),
			ch36Walkthrough(
				"c363.path-reveal",
				"Watch the difference between path and query.",
				ch36PathExample, "e363-path-only", []string{"fmt", "net/http"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Start main."},
					{Line: 2, Caption: "Build a request with a query."},
					{Line: 3, Caption: "Print only the path.", Output: "/readyz"},
				},
			),
			ch36Sandbox("c363.playground", ch36SandboxInstruction("c363.", 1), "c363.", 1),
		},
	},
	{
		ID: "c363.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch36Markdown("c363.watch", "- ~/livez~ answers about the process itself; it need not ping.\n- ~URL.Path~ has no query, which helps keep tokens out of logs.")},
	},
	{
		ID: "c363.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch36Markdown("c363.plan", "Set up the logger first, then separate HTTP branches and record the result after choosing it.\n\nReview shared state with concurrent requests in mind."),
			ch36Question("c363.live-check", "If the database is down, what should liveness do?", []string{"Check the database", "Stay independent", "Report readiness success"}, 1, []string{
				"The readiness probe checks dependencies.",
				"Process liveness does not depend on the database.",
				"Liveness and readiness ask different questions.",
			}),
		},
	},
	{
		ID: "c363.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch36Markdown("c363.turn", "Implement the logger on its own first. Then try each probe and inspect the JSON output and concurrent behavior.")},
	},
}

const ch36BossArabic = `1. كمّل ~StormControl~ في ~main.go~ باستخدام SQLite والـbroker المحقون.
2. خزّن الأوامر وكمّل تسليم الـpending بأمان.
3. اعرض الجاهزية والمقاييس الثابتة.

الـstarter فيه الـtypes والتواقيع.

ده boss lab؛ راجع ملكية كل مورد.`

const ch36BossEnglish = `1. Complete ~StormControl~ in ~main.go~ with SQLite and the injected broker.
2. Keep commands durable and deliver pending work safely.
3. Expose readiness and fixed metrics.

The starter contains the types and signatures.

This is a boss lab; track ownership of every resource.`

var ch36BossArabicSteps = []LessonStep{
	{
		ID: "c364.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch36Markdown("c364.mission", ch36BossArabic)},
	},
	{
		ID: "c364.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c364.outbox-state",
			"سجل الـoutbox يربط الأمر بحالة نشره. المثال sketch للـstate بس؛ الذاكرة لوحدها مش تخزين دائم.",
			ch36ReceiptExample, "e364-outbox-state", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 6, Caption: "ابدأ بإيصال pending.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "false"}}},
				{Line: 7, Caption: "اطبع الحالة الأولى.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "false"}}, Output: "storm-demo false"},
				{Line: 8, Caption: "الحالة تتغير بعد التسليم.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "true"}}, Output: "storm-demo false"},
				{Line: 9, Caption: "اطبع الحالة الجديدة.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "true"}}, Output: "storm-demo false\ntrue"},
			},
		)},
	},
	{
		ID: "c364.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch36Question("c364.predict-ack", ch36Prediction("توقّع النتيجة لما النشر ينجح ولما يفشل.", ch36AckExample), []string{"true false", "false true", "false false"}, 1, []string{
				"الحالة الفاشلة لازم تفضل pending.",
				"النجاح وحده يثبت التسليم.",
				"الحالتان مش بنفس النتيجة.",
			}),
			ch36Walkthrough(
				"c364.ack-reveal",
				"قارن علامة النشر للفشل والنجاح.",
				ch36AckExample, "e364-ack-state", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "الـfunction بتستقبل نتيجة النشر."},
					{Line: 2, Caption: "افحص الفشل أولاً."},
					{Line: 3, Caption: "الفشل مش بيأكد التسليم."},
					{Line: 5, Caption: "النجاح يؤكد التسليم."},
					{Line: 8, Caption: "العلامة تتبع النجاح.", Output: "false true"},
				},
			),
			ch36Sandbox("c364.playground", ch36SandboxInstruction("c364.", 0), "c364.", 0),
		},
	},
	{
		ID: "c364.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch36Markdown("c364.watch", "- ما تسيبش Rows مفتوحة وقت نداء broker.\n- سجّل الـerror برسالة ثابتة، من غير raw error ممكن يحتوي سر.")},
	},
	{
		ID: "c364.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch36Markdown("c364.plan", "ارسم انتقال الأمر من التسجيل للتسليم، وحدد مين يملك كل نتيجة بعد الفشل.\n\nقارن بيانات التخزين والـpending قبل وبعد Flush."),
			ch36Question("c364.ack-check", "إمتى علامة pending بتتغير؟", []string{"قبل Publish", "بعد نجاح Publish", "عند إنشاء handler"}, 1, []string{
				"الفشل قبل النشر لازم يترك شغل قابل للإعادة.",
				"النجاح هو اللي يثبت التسليم.",
				"إنشاء handler مش تسليم.",
			}),
		},
	},
	{
		ID: "c364.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch36Markdown("c364.turn", "ابدأ بفتح وإعادة فتح قاعدة محلية، وبعدها راجع النشر والـpending والـHTTP من غير broker حقيقي.")},
	},
}

var ch36BossEnglishSteps = []LessonStep{
	{
		ID: "c364.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch36Markdown("c364.mission", ch36BossEnglish)},
	},
	{
		ID: "c364.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c364.outbox-state",
			"An outbox row ties a command to its delivery state. This example sketches state only; memory alone is not durable storage.",
			ch36ReceiptExample, "e364-outbox-state", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 6, Caption: "Start with a pending receipt.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "false"}}},
				{Line: 7, Caption: "Print its first state.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "false"}}, Output: "storm-demo false"},
				{Line: 8, Caption: "State changes after delivery.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "true"}}, Output: "storm-demo false"},
				{Line: 9, Caption: "Print the new state.", Variables: []LessonVariable{{Name: "receipt.Published", Value: "true"}}, Output: "storm-demo false\ntrue"},
			},
		)},
	},
	{
		ID: "c364.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch36Question("c364.predict-ack", ch36Prediction("Predict the result when publishing succeeds and when it fails.", ch36AckExample), []string{"true false", "false true", "false false"}, 1, []string{
				"A failed delivery must remain pending.",
				"Only success confirms delivery.",
				"Both results are not the same.",
			}),
			ch36Walkthrough(
				"c364.ack-reveal",
				"Compare the delivery marker for failure and success.",
				ch36AckExample, "e364-ack-state", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "The function receives a publish result."},
					{Line: 2, Caption: "Check failure first."},
					{Line: 3, Caption: "Failure does not confirm delivery."},
					{Line: 5, Caption: "Success confirms delivery."},
					{Line: 8, Caption: "The marker follows success.", Output: "false true"},
				},
			),
			ch36Sandbox("c364.playground", ch36SandboxInstruction("c364.", 1), "c364.", 1),
		},
	},
	{
		ID: "c364.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch36Markdown("c364.watch", "- Do not keep Rows open while calling the broker.\n- Log a fixed message, without an error that might contain a secret.")},
	},
	{
		ID: "c364.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch36Markdown("c364.plan", "Draw the command's path from recording to delivery, and decide who owns each result after failure.\n\nCompare stored rows and pending work before and after Flush."),
			ch36Question("c364.ack-check", "When should the pending marker change?", []string{"Before Publish", "After Publish succeeds", "When a handler is created"}, 1, []string{
				"A failure before publishing must leave work retryable.",
				"Success confirms delivery.",
				"Creating a handler is not delivery.",
			}),
		},
	},
	{
		ID: "c364.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch36Markdown("c364.turn", "Start by opening and reopening a local database, then inspect delivery, pending work, and HTTP without a real broker.")},
	},
}

const ch36ContainersArabic = `1. كمّل ~RunPodCtl~ في ~main.go~ باستخدام ~Runtime~.
2. اعرض حالة الكبسولات ونفّذ الأوامر المسموحة.
3. خلّي الإيقاف محتاج تأكيد واضح.

الـstarter بيفصل الـCLI عن interface الـruntime.

ده lab؛ مفيش runtime حقيقي مطلوب.`

const ch36ContainersEnglish = `1. Complete ~RunPodCtl~ in ~main.go~ using ~Runtime~.
2. Show pod state and perform the allowed commands.
3. Make a stop require clear confirmation.

The starter separates the CLI from the runtime interface.

This is a lab; no real runtime is needed.`

var ch36ContainersArabicSteps = []LessonStep{
	{
		ID: "c362.mission-step", Title: "المهمة",
		Blocks: []LessonBlock{ch36Markdown("c362.mission", ch36ContainersArabic)},
	},
	{
		ID: "c362.idea-step", Title: "الفكرة",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c362.interface",
			"الـinterface بتحدد السلوك اللي الأداة محتاجاه. نوع بسيط يقدر يحققها من غير ما نربط الكود بمحرك حقيقي.",
			ch36InterfaceExample, "e362-interface", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "قبل التنفيذ."},
				{Line: 1, Caption: "الـinterface تصف Start."},
				{Line: 4, Caption: "Demo يحقق الـinterface."},
				{Line: 7, Caption: "المتغير يحمل تطبيقها."},
				{Line: 8, Caption: "العملية بترجّع string.", Output: "started probe-4"},
			},
		)},
	},
	{
		ID: "c362.try-step", Title: "جرّب",
		Blocks: []LessonBlock{
			ch36Question("c362.predict-interface", ch36Prediction("توقّع الـvalue اللي بترجعها الـinterface.", ch36InterfaceExampleReveal), []string{"relay", "Label", "nil"}, 0, []string{
				"ده اسم method، مش ناتجها.",
				"الـfunction بترجّع string value.",
				"Relay موجودة كـvalue.",
			}),
			ch36Walkthrough(
				"c362.interface-reveal",
				"شوف إزاي الـinterface بتوصل للـmethod.",
				ch36InterfaceExampleReveal, "e362-interface-reveal", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "قبل التشغيل."},
					{Line: 1, Caption: "Name تطلب Label."},
					{Line: 3, Caption: "Relay يحقق السلوك."},
					{Line: 5, Caption: "الـvalue تدخل في الـinterface."},
					{Line: 6, Caption: "الـmethod بترجّع string.", Output: "relay"},
				},
			),
			ch36Sandbox("c362.playground", ch36SandboxInstruction("c362.", 0), "c362.", 0),
		},
	},
	{
		ID: "c362.watch-step", Title: "خد بالك",
		Blocks: []LessonBlock{ch36Markdown("c362.watch", "- الأمر غير المفهوم لازم يخرج قبل أي نداء للـruntime.\n- Enter لوحده أو نهاية الدخل مش موافقة على الإيقاف.")},
	},
	{
		ID: "c362.check-step", Title: "فكّر واختبر",
		Blocks: []LessonBlock{
			ch36Markdown("c362.plan", "افصل تحليل الأمر عن فحص الحالة، وبعدها قرر هل الفعل مسموح أو محتاج تأكيد.\n\nمرّر نفس الـcontext لكل عملية."),
			ch36Question("c362.invalid-call", "أمر استخدامه غلط؛ إيه اللي يحصل للـruntime؟", []string{"ينادي List", "مايحصلش أي نداء", "ينادي Stop"}, 1, []string{
				"تحليل الأمر لازم يسبق لمس runtime.",
				"الأمر المرفوض مايناديش runtime.",
				"الإيقاف مش جزء من أمر غلط.",
			}),
		},
	},
	{
		ID: "c362.turn-step", Title: "دورك",
		Blocks: []LessonBlock{ch36Markdown("c362.turn", "ابدأ من توقيع ~RunPodCtl~ وجرّب الحالات واحدة واحدة مع fake runtime. راجع مخرجات stdout وstderr لوحدهم.")},
	},
}

var ch36ContainersEnglishSteps = []LessonStep{
	{
		ID: "c362.mission-step", Title: "Mission",
		Blocks: []LessonBlock{ch36Markdown("c362.mission", ch36ContainersEnglish)},
	},
	{
		ID: "c362.idea-step", Title: "Idea",
		Blocks: []LessonBlock{ch36Walkthrough(
			"c362.interface",
			"An interface describes the behavior the tool needs. A small type can satisfy it without tying the code to a real engine.",
			ch36InterfaceExample, "e362-interface", []string{"fmt"},
			[]LessonFrame{
				{Line: 0, Caption: "Before execution."},
				{Line: 1, Caption: "The interface describes Start."},
				{Line: 4, Caption: "Demo satisfies the interface."},
				{Line: 7, Caption: "The variable holds its implementation."},
				{Line: 8, Caption: "The method returns text.", Output: "started probe-4"},
			},
		)},
	},
	{
		ID: "c362.try-step", Title: "Try it",
		Blocks: []LessonBlock{
			ch36Question("c362.predict-interface", ch36Prediction("Predict the result returned by the interface.", ch36InterfaceExampleReveal), []string{"relay", "Label", "nil"}, 0, []string{
				"That is the method name, not its result.",
				"The method returns a string value.",
				"Relay is present as a value.",
			}),
			ch36Walkthrough(
				"c362.interface-reveal",
				"Watch how the interface reaches the method.",
				ch36InterfaceExampleReveal, "e362-interface-reveal", []string{"fmt"},
				[]LessonFrame{
					{Line: 0, Caption: "Before execution."},
					{Line: 1, Caption: "Name requires Label."},
					{Line: 3, Caption: "Relay provides the behavior."},
					{Line: 5, Caption: "The value enters the interface."},
					{Line: 6, Caption: "The method returns text.", Output: "relay"},
				},
			),
			ch36Sandbox("c362.playground", ch36SandboxInstruction("c362.", 1), "c362.", 1),
		},
	},
	{
		ID: "c362.watch-step", Title: "Watch out",
		Blocks: []LessonBlock{ch36Markdown("c362.watch", "- An unknown command must exit before any runtime call.\n- Enter alone or end of input is not confirmation to stop.")},
	},
	{
		ID: "c362.check-step", Title: "Think & check",
		Blocks: []LessonBlock{
			ch36Markdown("c362.plan", "Separate command parsing from state lookup, then decide whether an action is allowed or needs confirmation.\n\nPass the same context to each operation."),
			ch36Question("c362.invalid-call", "A command has invalid usage; what happens to the runtime?", []string{"List is called", "No runtime call happens", "Stop is called"}, 1, []string{
				"Parse the command before touching the runtime.",
				"A rejected command does not call the runtime.",
				"Stopping is not part of an invalid command.",
			}),
		},
	},
	{
		ID: "c362.turn-step", Title: "Your turn",
		Blocks: []LessonBlock{ch36Markdown("c362.turn", "Start from the ~RunPodCtl~ signature and try cases one by one with a fake runtime. Check stdout and stderr separately.")},
	},
}
