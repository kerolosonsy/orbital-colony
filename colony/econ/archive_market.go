package econ

const ArchiveDay = 1200

func init() { Contracts = append(Contracts, archiveMarketContracts...) }

var archiveMarketContracts = []*Contract{
	{ID: "d18-audit", Chapter: 18, Window: 4800, Title: "مراجعة الأرشيف", TitleEn: "Archive Audit", Brief: "اجمع كريدت كل سلعة من ledger.csv الحالي وتجاهل الصفوف الناقصة والأرقام البايظة والسلع المجهولة. ابعت map باسم audit بعد بيع واحد سليم على الأقل.", BriefEn: "Sum credits per good in the current ledger.csv, skipping missing fields, invalid numbers and unknown goods. Submit a map as audit after at least one valid sale.", Requires: []string{"sel:csv.NewReader"}, Goal: Goal{Submit: "audit", MaxLost: NoLimit}, Reward: Reward{Credits: 1800}},
	{ID: "d19-shift-report", Chapter: 19, Window: 7200, Title: "تقرير الورديات", TitleEn: "Shift Report", Brief: "اليوم 1200 tick. اجمع goods في أحداث harvest من events.jsonl ونسختيه .1 و.2 لآخر 3 أيام مكتملة. ابعت shift-report كخريطة رقم اليوم النصي لخريطة السلع؛ اليوم الفاضي خريطة فاضية.", BriefEn: "A day is 1200 ticks. Sum goods in harvest events from events.jsonl and its .1 and .2 rotations for the last three completed days. Submit shift-report as a map from decimal day strings to goods maps; empty days use empty maps.", Requires: []string{"sel:json.NewDecoder"}, Goal: Goal{Submit: "shift-report", MaxLost: NoLimit}, Reward: Reward{Credits: 1900}},
	{ID: "d22-atomic-journal", Chapter: 22, Window: 1800, Title: "السجل الذرّي", TitleEn: "Atomic Journal", Brief: "انسخ plan-old.txt إلى plan.txt ثم بدّل بينه وبين plan-new.txt باستمرار بحفظ ذرّي. السيرفر بيقتل البرنامج في وقت متحدد بالـ seed، وبيتحقق من المحتوى وسطر sha256 الأخير.", BriefEn: "Initialize plan.txt from plan-old.txt, then continuously alternate it with plan-new.txt using atomic saves. The server kills the program at a seeded tick and verifies the contents and final sha256 line.", Requires: []string{"sel:os.Rename"}, Goal: Goal{Submit: "atomic-journal", MaxLost: NoLimit}, Reward: Reward{Credits: 2200}},
	{ID: "h24-bulk-seller", Chapter: 24, Window: 4800, Title: "بائع الجملة", TitleEn: "Bulk Seller", Brief: "نفّذ 20 طلب سوق بمفاتيح مختلفة من غير فقد أو تكرار، بعمال محدودين وtimeouts ومحاولات محدودة وbackoff. مخزونك هو مصدر البضاعة.", BriefEn: "Fill 20 market orders with distinct keys and no loss or duplicate sales, using bounded workers, timeouts, bounded retries and backoff. Supply goods from your inventory.", Requires: []string{"sel:http.NewRequest", "waitgroup", "chan"}, Goal: Goal{Submit: "bulk-seller", MaxLost: NoLimit}, Reward: Reward{Credits: 2400}},
	{ID: "h33-stubborn-client", Chapter: 33, Window: 4800, Title: "العميل العنيد", TitleEn: "Stubborn Client", Brief: "نفّذ 20 طلب مختلف؛ 30% من محاولات السوق بتفشل ونجاح واحد رده بطيء. استخدم نفس المفتاح عند إعادة المحاولة وحدود للتزامن والانتظار والمحاولات.", BriefEn: "Fill 20 distinct orders; 30% of market attempts fail and one successful response is slow. Reuse keys on retry and bound concurrency, timeouts and retries.", Requires: []string{"sel:http.NewRequest", "waitgroup", "chan"}, Goal: Goal{Submit: "stubborn-client", MaxLost: NoLimit}, Reward: Reward{Credits: 3300}},
}

func addArchiveMarketNodes() {
	for _, n := range []Node{
		{ID: "sys.archive", Branch: "البيانات", Name: "الأرشيف", NameEn: "Archive", Desc: "ملفات بيع وأحداث دائمة للتحليل.", DescEn: "Persistent sale and event files for analysis.", After: []Prereq{{"api.scanner", 1}}, Chapter: 18},
		{ID: "sys.market", Branch: "الخدمات", Name: "السوق HTTP", NameEn: "HTTP Market", Desc: "بيع بالجملة عبر HTTP بسعر أعلى 20% ومفاتيح تمنع التكرار.", DescEn: "Bulk HTTP sales at a 20% premium with idempotency keys.", After: []Prereq{{"api.prices", 1}}, Chapter: 26},
	} {
		addResearchNode(n)
	}
}

// RecordSubmission accepts only the active contract's answer key.
func (s *State) RecordSubmission(key string, correct bool) bool {
	c, r := s.ActiveContract()
	if c == nil || c.Goal.Submit != key || key == "" || r.Submission != "" {
		return false
	}
	r.Submission = "wrong"
	if correct {
		r.Submission = "correct"
	}
	return true
}

func (s *State) submissionProgress(c *Contract, r *ContractRun, tick int64) (bool, bool, [2]string) {
	switch r.Submission {
	case "correct":
		return true, false, [2]string{}
	case "wrong":
		return false, true, [2]string{"الإجابة غلط", "Incorrect answer"}
	}
	if c.Window > 0 && tick-r.Start > c.Window {
		return false, true, [2]string{"الوقت خلص", "Time ran out"}
	}
	return false, false, [2]string{"مستني " + c.Goal.Submit, "Waiting for " + c.Goal.Submit}
}

func AtomicJournalAt(seed, start int64) int64 {
	return start + 160 + int64(EventHash(seed, start, 22)%120)
}
