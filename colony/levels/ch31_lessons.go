package levels

func ch31Markdown(id, text string) LessonBlock {
	return LessonBlock{ID: id, Kind: "markdown", Text: text}
}

func ch31Question(id, text string, options []string, answer int, feedback []string) LessonBlock {
	return LessonBlock{
		ID:       id,
		Kind:     "question",
		Text:     text,
		Question: &LessonQuestion{Options: options, Answer: answer, Feedback: feedback},
	}
}

func ch31Caption(arabic, englishText string, useEnglish bool) string {
	if useEnglish {
		return englishText
	}
	return arabic
}

func ch31Walkthrough(id, text, code, exampleID string, imports []string, frames []LessonFrame) LessonBlock {
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

const ch31ReadChunkCode = `func main() {
	r := strings.NewReader("signal")
	chunk := make([]byte, 3)
	n, _ := r.Read(chunk)
	fmt.Println(string(chunk[:n]))
}`

func ch31ReadChunk(text string, english bool) LessonBlock {
	return ch31Walkthrough("c311.reader-chunk", text, ch31ReadChunkCode, "e311-reader-chunk", []string{"fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("القارئ فيه stream.", "The reader holds a stream.", english)},
		{Line: 3, Caption: ch31Caption("المخزن يطلب ثلاث بايتات.", "The buffer asks for three bytes.", english)},
		{Line: 4, Caption: ch31Caption("القراءة ملأت الجزء.", "The read fills that chunk.", english), Variables: []LessonVariable{{Name: "n", Value: "3"}, {Name: "chunk", Value: "[115 105 103]"}}},
		{Line: 5, Caption: ch31Caption("اطبع الجزء المقروء.", "Print the bytes read.", english), Variables: []LessonVariable{{Name: "n", Value: "3"}}, Output: "sig"},
	})
}

const ch31ByteCountCode = `func main() {
	word := "é"
	fmt.Println(len(word), utf8.RuneCountInString(word))
}`

func ch31ByteCount(text string, english bool) LessonBlock {
	return ch31Walkthrough("c311.byte-count", text, ch31ByteCountCode, "e311-byte-count", []string{"fmt", "unicode/utf8"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("الحرف يحتل بايتين.", "The character uses two bytes.", english), Variables: []LessonVariable{{Name: "word", Value: "é"}}},
		{Line: 3, Caption: ch31Caption("قارن البايتات بالـrunes.", "Compare bytes and runes.", english), Variables: []LessonVariable{{Name: "word", Value: "é"}}, Output: "2 1"},
	})
}

var ch31Lesson311Arabic = []LessonStep{
	{ID: "c311.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch31Markdown("c311.mission", "1. رجّع سطور المنارة كاملة وبالترتيب.\n2. احترم حد المحتوى بالبايت.\n3. خلّي ملكية الاتصال واضحة.")}},
	{ID: "c311.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch31ReadChunk("TCP بيقدم stream بايتات؛ القراية الواحدة ممكن ترجع جزء من الرسالة.\n\nالبروتوكول هو اللي يحدد الفاصل. احتفظ بالبايتات لحد ما تعرف إن الإطار اكتمل.\n\nالـdeadline مطلق ومحقون، و~net.Conn~ محتاج مالك يقفله.\n\nالحد هنا بايتات محتوى، فالحرف المرئي مش دايماً بايت واحد.", false)}},
	{ID: "c311.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch31Question("c311.predict-bytes", "كام بايت وكم rune في النص؟\n\n~~~go\nfunc main() {\n\tword := \"é\"\n\tfmt.Println(len(word), utf8.RuneCountInString(word))\n}\n~~~", []string{"2 1", "1 1", "2 2"}, 0, []string{
			"الحرف UTF-8 له بايتان وrune واحدة.",
			"len للنص بيعد البايتات.",
			"RuneCountInString يعد وحدات Unicode.",
		}),
		ch31ByteCount("راقب الفرق بين طول string وعدّ الحروف.", false),
	}},
	{ID: "c311.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch31Markdown("c311.watch", "- EOF وسط سطر مش إغلاق نظيف؛ ما تبعتش الجزء الناقص.\n- الكتابة ممكن تكون جزئية؛ تابع الباقي أو أبلغ عن التوقف.")}},
	{ID: "c311.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch31Markdown("c311.plan", "افصل التملك عن النقل: تحقق من الحد، اضبط الموعد، وبعدها عالج الإطار.\n\nاسأل إيه اللي اتكتب بالفعل قبل أي خطأ، ومين يقفل الاتصال عند كل خروج."),
		ch31Question("c311.eof", "EOF بعد محتوى من غير newline معناه إيه؟", []string{"إطار ناقص", "سطر كامل", "تجاهل كل الخطأ"}, 0, []string{
			"المحتوى لسه ما وصلش للفاصل.",
			"الفاصل جزء من تعريف الإطار.",
			"سبب النقل لازم يفضل ظاهراً.",
		}),
	}},
	{ID: "c311.turn-step", Title: "دورك", Blocks: []LessonBlock{ch31Markdown("c311.turn", "ابدأ بحالة قراءة واحدة؛ بعدها راجع الإغلاق والردود السابقة.")}},
}

var ch31Lesson311English = []LessonStep{
	{ID: "c311.mission-step", Title: "Mission", Blocks: []LessonBlock{ch31Markdown("c311.mission", "1. Echo complete beacon lines in order.\n2. Respect the content limit in bytes.\n3. Make connection ownership clear.")}},
	{ID: "c311.idea-step", Title: "Idea", Blocks: []LessonBlock{ch31ReadChunk("TCP provides a byte stream; one read may return only part of a message.\n\nThe protocol defines the delimiter. Keep bytes until you know the frame is complete.\n\nA deadline is absolute and injected, and a ~net.Conn~ needs one owner to close it.\n\nThis limit counts content bytes, so one visible character may use several bytes.", true)}},
	{ID: "c311.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch31Question("c311.predict-bytes", "How many bytes and runes are in this string?\n\n~~~go\nfunc main() {\n\tword := \"é\"\n\tfmt.Println(len(word), utf8.RuneCountInString(word))\n}\n~~~", []string{"2 1", "1 1", "2 2"}, 0, []string{
			"This UTF-8 character uses two bytes and one rune.",
			"len on a string counts bytes.",
			"RuneCountInString counts Unicode runes.",
		}),
		ch31ByteCount("Watch the difference between string length and rune count.", true),
	}},
	{ID: "c311.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch31Markdown("c311.watch", "- EOF inside a line is incomplete; do not send its partial bytes.\n- Writes may be partial; continue the remainder or report no progress.")}},
	{ID: "c311.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch31Markdown("c311.plan", "Separate ownership from transfer: validate the limit, set the deadline, then process a frame.\n\nAsk what was already written before an error, and who closes the connection on every exit."),
		ch31Question("c311.eof", "What does EOF before a newline mean?", []string{"An incomplete frame", "A complete line", "Ignore the error"}, 0, []string{
			"The content has not reached its delimiter.",
			"The delimiter is part of the frame definition.",
			"The transport cause must remain visible.",
		}),
	}},
	{ID: "c311.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch31Markdown("c311.turn", "Start with one read state; then check closure and earlier replies.")}},
}

const ch31FullMailboxCode = `func main() {
	mail := make(chan string, 1)
	mail <- "first"
	select {
	case mail <- "later":
		fmt.Println("sent")
	default:
		fmt.Println("busy")
	}
}`

func ch31FullMailbox(text string, english bool) LessonBlock {
	return ch31Walkthrough("c312.full-mailbox", text, ch31FullMailboxCode, "e312-full-mailbox", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("صندوق بسعة واحدة.", "The mailbox holds one item.", english)},
		{Line: 3, Caption: ch31Caption("الصندوق بقى مليان.", "The mailbox is now full.", english), Variables: []LessonVariable{{Name: "mail", Value: "[first]"}}},
		{Line: 4, Caption: ch31Caption("select يجرب الإرسال.", "select tries the send.", english)},
		{Line: 8, Caption: ch31Caption("default يتعامل مع الامتلاء.", "default handles the full buffer.", english), Variables: []LessonVariable{{Name: "mail", Value: "[first]"}}, Output: "busy"},
	})
}

const ch31EmptyMailboxCode = `func main() {
	mail := make(chan string, 1)
	select {
	case mail <- "ping":
		fmt.Println("sent")
	default:
		fmt.Println("full")
	}
}`

func ch31EmptyMailbox(text string, english bool) LessonBlock {
	return ch31Walkthrough("c312.empty-mailbox", text, ch31EmptyMailboxCode, "e312-empty-mailbox", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("الصندوق فاضي.", "The mailbox is empty.", english)},
		{Line: 3, Caption: ch31Caption("select يختار الإرسال.", "select chooses the send.", english)},
		{Line: 5, Caption: ch31Caption("الرسالة تدخل الصندوق.", "The message enters the mailbox.", english), Variables: []LessonVariable{{Name: "mail", Value: "[ping]"}}},
		{Line: 6, Caption: ch31Caption("الإرسال تم فوراً.", "The send succeeds now.", english), Variables: []LessonVariable{{Name: "mail", Value: "[ping]"}}, Output: "sent"},
	})
}

var ch31Lesson312Arabic = []LessonStep{
	{ID: "c312.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch31Markdown("c312.mission", "1. خلّي حلقة واحدة تملك خريطة الأعضاء.\n2. ابعت الرسائل من غير تعطيل الباقيين.\n3. اقفل الصناديق عند الخروج.")}},
	{ID: "c312.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch31FullMailbox("القناة صندوق بريد؛ المالك الواحد يحافظ على حالة الأعضاء من غير قفل للخريطة.\n\n~select~ مع ~default~ يجرب إرسالاً فورياً. لو الصندوق مليان، المسار يكمل بدل ما يستنى.\n\nفي المهمة، الامتلاء معناه فصل العضو البطيء؛ الرسالة ممكن تفضل مقروءة قبل ما القناة تخلص.\n\nالمنتج يملك قناة الأوامر؛ الغرفة تملك صندوق العضو المقبول.", false)}},
	{ID: "c312.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch31Question("c312.predict-send", "والصندوق فاضي، إيه اللي يطبع؟\n\n~~~go\nfunc main() {\n\tmail := make(chan string, 1)\n\tselect {\n\tcase mail <- \"ping\":\n\t\tfmt.Println(\"sent\")\n\tdefault:\n\t\tfmt.Println(\"full\")\n\t}\n}\n~~~", []string{"sent", "full", "لا شيء"}, 0, []string{
			"الإرسال متاح في السعة الفاضية.",
			"default يشتغل لو الإرسال مش جاهز.",
			"المسار الجاهز يطبع نتيجة.",
		}),
		ch31EmptyMailbox("اتبع اختيار select قبل قراءة الرسالة المطبوعة.", false),
	}},
	{ID: "c312.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch31Markdown("c312.watch", "- ما تقفلش قناة عضو من برّه مالكها.\n- رد الأمر نقطة مزامنة؛ ابعته بعد ما تخلص أثره.")}},
	{ID: "c312.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch31Markdown("c312.plan", "رتّب الحالة حول owner loop: أمر واحد، تحقق، تحديث، ثم رد.\n\nوقت النهاية، حدّد مين بيقفل كل قناة. راجع إن العضو البطيء ما يعطلش الإرسال لغيره."),
		ch31Question("c312.channel-owner", "مين يقفل صندوق عضو انضم بنجاح؟", []string{"الغرفة المالكة", "العضو القارئ", "أي مرسل"}, 0, []string{
			"نفس المالك مسؤول عن الإرسال والإغلاق.",
			"المشترك يقرأ فقط حسب العقد.",
			"ملكية الإغلاق ماينفعش تتشارك عشوائياً.",
		}),
	}},
	{ID: "c312.turn-step", Title: "دورك", Blocks: []LessonBlock{ch31Markdown("c312.turn", "ابدأ بحلقة المالك، وبعدها أضف كل انتقال حالة لوحده.")}},
}

var ch31Lesson312English = []LessonStep{
	{ID: "c312.mission-step", Title: "Mission", Blocks: []LessonBlock{ch31Markdown("c312.mission", "1. Let one loop own the member map.\n2. Deliver messages without blocking others.\n3. Close mailboxes on shutdown.")}},
	{ID: "c312.idea-step", Title: "Idea", Blocks: []LessonBlock{ch31FullMailbox("A channel is a mailbox; one owner keeps member state simple without a map lock.\n\n~select~ with ~default~ attempts an immediate send. If the mailbox is full, the loop continues instead of waiting.\n\nHere, a full mailbox means disconnecting a slow member; queued messages can still be read after closure.\n\nThe producer owns the command channel; the room owns an accepted member mailbox.", true)}},
	{ID: "c312.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch31Question("c312.predict-send", "What prints when the mailbox is empty?\n\n~~~go\nfunc main() {\n\tmail := make(chan string, 1)\n\tselect {\n\tcase mail <- \"ping\":\n\t\tfmt.Println(\"sent\")\n\tdefault:\n\t\tfmt.Println(\"full\")\n\t}\n}\n~~~", []string{"sent", "full", "nothing"}, 0, []string{
			"The send can use the empty buffer slot.",
			"default runs when the send is not ready.",
			"The ready case prints a result.",
		}),
		ch31EmptyMailbox("Follow the select choice before reading its printed result.", true),
	}},
	{ID: "c312.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch31Markdown("c312.watch", "- Do not close a member queue outside its owner.\n- Acknowledgement is a sync point; send it after effects finish.")}},
	{ID: "c312.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch31Markdown("c312.plan", "Organize state around the owner loop: receive, validate, update, then reply.\n\nAt shutdown, identify who closes each channel. Check that a slow member cannot block delivery to others."),
		ch31Question("c312.channel-owner", "Who closes a successfully joined member queue?", []string{"The room that owns it", "The reading member", "Any sender"}, 0, []string{
			"The same owner controls sends and closure.",
			"The subscriber only reads under this contract.",
			"Channel close ownership cannot be random.",
		}),
	}},
	{ID: "c312.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch31Markdown("c312.turn", "Start with the owner loop, then add each state transition separately.")}},
}

const ch31PacketFieldsCode = `func main() {
	fields := strings.Split("v2|r7|pier", "|")
	fmt.Println(fields[0], fields[2])
}`

func ch31PacketFields(text string, english bool) LessonBlock {
	return ch31Walkthrough("c313.packet-fields", text, ch31PacketFieldsCode, "e313-packet-fields", []string{"fmt", "strings"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("قسّم datagram واحدة.", "Split one datagram.", english)},
		{Line: 3, Caption: ch31Caption("اقرأ حقول محددة.", "Read selected fields.", english), Variables: []LessonVariable{{Name: "fields", Value: "[v2 r7 pier]"}}},
		{Line: 4, Caption: ch31Caption("اطبع أول وآخر حقل.", "Print the first and last fields.", english), Variables: []LessonVariable{{Name: "fields", Value: "[v2 r7 pier]"}}, Output: "v2 pier"},
	})
}

const ch31OnePacketCode = `func main() {
	packet := []byte("v2|r7|pier")
	fmt.Println(string(packet))
}`

func ch31OnePacket(text string, english bool) LessonBlock {
	return ch31Walkthrough("c313.one-packet", text, ch31OnePacketCode, "e313-one-packet", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("بايتات packet واحدة.", "Bytes of one packet.", english), Variables: []LessonVariable{{Name: "packet", Value: "[118 50 124 114 55 124 112 105 101 114]"}}},
		{Line: 3, Caption: ch31Caption("اعرض النص التوضيحي.", "Show the illustrative text.", english), Variables: []LessonVariable{{Name: "packet", Value: "[118 50 124 114 55 124 112 105 101 114]"}}, Output: "v2|r7|pier"},
	})
}

var ch31Lesson313Arabic = []LessonStep{
	{ID: "c313.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch31Markdown("c313.mission", "1. افتح listener محلياً عند الطلب.\n2. افك packet واحدة بحد واضح.\n3. رد لمصدر الإرسال الفعلي.")}},
	{ID: "c313.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch31PacketFields("UDP ينقل datagrams منفصلة؛ نجاح الإرسال مش تأكيد وصول.\n\nقسّم الحزمة بعد التحقق من شكلها، واحتفظ بعنوان المصدر اللي رجع مع القراءة.\n\nالفاصل بين codec والنقل يخلي الترميز قابل للتجربة من غير socket.\n\nالموعد deadline مطلق. المصنع وحده يفتح الاتصال، واللي ناداه هو اللي يقفله.", false)}},
	{ID: "c313.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch31Question("c313.predict-packet", "إيه اللي يطبعه تحويل packet لstring؟\n\n~~~go\nfunc main() {\n\tpacket := []byte(\"v2|r7|pier\")\n\tfmt.Println(string(packet))\n}\n~~~", []string{"v2|r7|pier", "[118 50 124 114 55 124 112 105 101 114]", "v2 pier"}, 0, []string{
			"string تعرض محتوى البايتات كنص.",
			"دي صيغة slice البايتات قبل التحويل.",
			"التقسيم لحقول محتاج خطوة تانية.",
		}),
		ch31OnePacket("تابع البايتات وهي بتظهر كنص للحزمة.", false),
	}},
	{ID: "c313.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch31Markdown("c313.watch", "- مخزن أكبر ببايت يكشف datagram أطول من الحد.\n- في UDP، short write مش معناه إرسال باقي الحزمة لاحقاً.")}},
	{ID: "c313.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch31Markdown("c313.plan", "افصل التحقق من النص عن التعامل مع socket. ارفض شكل الطلب قبل بناء الرد.\n\nبعدها راجع عنوان الرد والموعد والملكية: الدالة المستعيرة لا تقفل الاتصال."),
		ch31Question("c313.reply-address", "عنوان الرد ييجي منين؟", []string{"عنوان ReadFrom", "حقل داخل النص", "عنوان افتراضي"}, 0, []string{
			"ReadFrom يعطي عنوان المرسل الفعلي.",
			"النص لا يثبت عنوان النقل.",
			"العنوان الافتراضي ممكن يرد لطرف غلط.",
		}),
	}},
	{ID: "c313.turn-step", Title: "دورك", Blocks: []LessonBlock{ch31Markdown("c313.turn", "ابدأ بتوقيع codec؛ بعدها صِل طلباً واحداً بالنقل المحقون.")}},
}

var ch31Lesson313English = []LessonStep{
	{ID: "c313.mission-step", Title: "Mission", Blocks: []LessonBlock{ch31Markdown("c313.mission", "1. Open a local listener on request.\n2. Decode one packet within a clear bound.\n3. Reply to the actual sender.")}},
	{ID: "c313.idea-step", Title: "Idea", Blocks: []LessonBlock{ch31PacketFields("UDP carries separate datagrams; a successful send does not confirm delivery.\n\nValidate a packet's shape before splitting it, and keep the source address returned by the read.\n\nSeparating codec from transport makes encoding testable without a socket.\n\nA deadline is absolute. The factory opens the connection, and its caller closes it.", true)}},
	{ID: "c313.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch31Question("c313.predict-packet", "What does converting the packet to a string print?\n\n~~~go\nfunc main() {\n\tpacket := []byte(\"v2|r7|pier\")\n\tfmt.Println(string(packet))\n}\n~~~", []string{"v2|r7|pier", "[118 50 124 114 55 124 112 105 101 114]", "v2 pier"}, 0, []string{
			"string displays the byte content as text.",
			"That is the byte slice before conversion.",
			"Splitting into fields takes another step.",
		}),
		ch31OnePacket("Follow the bytes as they appear in the packet text.", true),
	}},
	{ID: "c313.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch31Markdown("c313.watch", "- A buffer one byte larger detects an oversized datagram.\n- In UDP, a short write does not send a remainder later.")}},
	{ID: "c313.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch31Markdown("c313.plan", "Separate text validation from socket handling. Reject malformed requests before building a reply.\n\nThen check the reply address, deadline, and ownership: a borrowed connection stays open."),
		ch31Question("c313.reply-address", "Where does the reply address come from?", []string{"The ReadFrom source", "A field in the text", "A default address"}, 0, []string{
			"ReadFrom returns the actual sender address.",
			"Text does not establish the transport address.",
			"A default address could reply to the wrong peer.",
		}),
	}},
	{ID: "c313.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch31Markdown("c313.turn", "Start with the codec signatures; then connect one request to injected transport.")}},
}

const ch31BigEndianCode = `func main() {
	wire := []byte{1, 2}
	n := binary.BigEndian.Uint16(wire)
	fmt.Println(n)
}`

func ch31BigEndian(text string, english bool) LessonBlock {
	return ch31Walkthrough("c314.big-endian", text, ch31BigEndianCode, "e314-big-endian", []string{"encoding/binary", "fmt"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("البايتات محفوظة بالترتيب.", "Bytes stay in order.", english), Variables: []LessonVariable{{Name: "wire", Value: "[1 2]"}}},
		{Line: 3, Caption: ch31Caption("BigEndian يقرأ الرقم.", "BigEndian reads the number.", english), Variables: []LessonVariable{{Name: "n", Value: "258"}}},
		{Line: 4, Caption: ch31Caption("اطبع القيمة العددية.", "Print the numeric value.", english), Variables: []LessonVariable{{Name: "n", Value: "258"}}, Output: "258"},
	})
}

const ch31HeaderCode = `func main() {
	wire := []byte{0, 9, 7}
	n := binary.BigEndian.Uint16(wire[:2])
	fmt.Println(n, len(wire[2:]))
}`

func ch31Header(text string, english bool) LessonBlock {
	return ch31Walkthrough("c314.header-length", text, ch31HeaderCode, "e314-header-length", []string{"encoding/binary", "fmt"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("الإطار فيه بايتات.", "The frame contains bytes.", english)},
		{Line: 3, Caption: ch31Caption("اقرأ الحقل الأول.", "Read the first field.", english), Variables: []LessonVariable{{Name: "n", Value: "9"}}},
		{Line: 4, Caption: ch31Caption("الباقي بايت واحد.", "One byte remains.", english), Variables: []LessonVariable{{Name: "n", Value: "9"}}, Output: "9 1"},
	})
}

var ch31Lesson314Arabic = []LessonStep{
	{ID: "c314.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch31Markdown("c314.mission", "1. فك إطارات التليمتري المحفوظة.\n2. راجع الطول قبل قراءة الحقول.\n3. اكتب fuzz test لعقد parser.")}},
	{ID: "c314.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch31BigEndian("الـbinary frame بيانات مرتبة، مش نص مفصول بعلامة.\n\n~binary.BigEndian~ يقرأ البايت الأعلى أولاً. اعمل حدود slice بعد ما تتأكد إن الهيدر موجود.\n\nالطول المعلن لازم يطابق الشكل كله. والـpayload الناتجة تملك نسخة مستقلة.\n\nراجع كمان ترتيب الأخطاء؛ خطأ النوع له أولوية بعد سلامة البنية.", false)}},
	{ID: "c314.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch31Question("c314.predict-header", "بعد قراءة الرقم، كام بايت باقي؟\n\n~~~go\nfunc main() {\n\twire := []byte{0, 9, 7}\n\tn := binary.BigEndian.Uint16(wire[:2])\n\tfmt.Println(n, len(wire[2:]))\n}\n~~~", []string{"9 1", "9 2", "2304 1"}, 0, []string{
			"أول بايتين يمثلوا 9 وواحد يفضل.",
			"slice المتبقية تبدأ بعد البايتين.",
			"BigEndian يقرأ الصفر قبل التسعة.",
		}),
		ch31Header("اتبع قراءة الحقل، وبعدها طول الجزء المتبقي.", false),
	}},
	{ID: "c314.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch31Markdown("c314.watch", "- ما تعملش slicing قبل فحص طول الهيدر.\n- فك الطول إلى int قبل جمع حجم الهيدر.")}},
	{ID: "c314.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch31Markdown("c314.plan", "قسّم parser لمراحل تحقق صغيرة قبل الوصول للبيانات.\n\nصمّم oracle مستقل في fuzzing؛ تكرار نفس حساب parser ممكن يكرر نفس الغلط."),
		ch31Question("c314.ownership", "تعديل payload الناتجة المفروض يعمل إيه؟", []string{"ما يغيرش wire", "يغير wire", "يغير ترتيب الهيدر"}, 0, []string{
			"النتيجة تملك نسخة منفصلة من البايتات.",
			"المشاركة تكشف تغيير المدخل.",
			"الترتيب جزء من البيانات الأصلية.",
		}),
	}},
	{ID: "c314.turn-step", Title: "دورك", Blocks: []LessonBlock{ch31Markdown("c314.turn", "ابدأ بفحوص الحدود؛ بعد القبول قارن الحقول والملكية.")}},
}

var ch31Lesson314English = []LessonStep{
	{ID: "c314.mission-step", Title: "Mission", Blocks: []LessonBlock{ch31Markdown("c314.mission", "1. Parse archived telemetry frames.\n2. Check length before reading fields.\n3. Write a fuzz test for the parser contract.")}},
	{ID: "c314.idea-step", Title: "Idea", Blocks: []LessonBlock{ch31BigEndian("A binary frame is ordered data, not text split by a delimiter.\n\n~binary.BigEndian~ reads the high byte first. Slice only after confirming the header exists.\n\nThe declared length must match the whole shape. The returned payload owns an independent copy.\n\nAlso preserve error precedence: kind errors come after structural checks.", true)}},
	{ID: "c314.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch31Question("c314.predict-header", "After reading the number, how many bytes remain?\n\n~~~go\nfunc main() {\n\twire := []byte{0, 9, 7}\n\tn := binary.BigEndian.Uint16(wire[:2])\n\tfmt.Println(n, len(wire[2:]))\n}\n~~~", []string{"9 1", "9 2", "2304 1"}, 0, []string{
			"The first two bytes represent 9, leaving one.",
			"The remaining slice starts after two bytes.",
			"BigEndian reads the zero before the nine.",
		}),
		ch31Header("Follow the field read, then count the remaining bytes.", true),
	}},
	{ID: "c314.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch31Markdown("c314.watch", "- Never slice before checking header length.\n- Convert the decoded length to int before adding header size.")}},
	{ID: "c314.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch31Markdown("c314.plan", "Split the parser into small validation stages before accessing data.\n\nBuild an independent fuzz oracle; repeating the parser's calculation can repeat the same mistake."),
		ch31Question("c314.ownership", "What should editing the returned payload do?", []string{"Leave wire unchanged", "Change wire", "Reorder the header"}, 0, []string{
			"The result owns a separate byte copy.",
			"Shared storage would expose input mutation.",
			"Header order belongs to the original data.",
		}),
	}},
	{ID: "c314.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch31Markdown("c314.turn", "Start with bounds checks; after acceptance, compare fields and ownership.")}},
}

const ch31IPv4Code = `func main() {
	ip := net.ParseIP("192.0.2.8")
	fmt.Println(ip.To4() != nil, ip.String())
}`

func ch31IPv4(text string, english bool) LessonBlock {
	return ch31Walkthrough("c315.ip-family", text, ch31IPv4Code, "e315-ip-family", []string{"fmt", "net"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("ParseIP يفك العنوان.", "ParseIP parses the address.", english), Variables: []LessonVariable{{Name: "ip", Value: "192.0.2.8"}}},
		{Line: 3, Caption: ch31Caption("To4 يصنف العائلة.", "To4 classifies its family.", english), Variables: []LessonVariable{{Name: "ip", Value: "192.0.2.8"}}, Output: "true 192.0.2.8"},
	})
}

const ch31DNSCauseCode = `func main() {
	cause := &net.DNSError{IsNotFound: true}
	var dns *net.DNSError
	fmt.Println(errors.As(cause, &dns), dns.IsNotFound)
}`

func ch31DNSCause(text string, english bool) LessonBlock {
	return ch31Walkthrough("c315.dns-cause", text, ch31DNSCauseCode, "e315-dns-cause", []string{"errors", "fmt", "net"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("الخطأ يحمل نوعه.", "The error carries its type.", english), Variables: []LessonVariable{{Name: "cause.IsNotFound", Value: "true"}}},
		{Line: 3, Caption: ch31Caption("جهّز هدف النوع.", "Prepare a type target.", english)},
		{Line: 4, Caption: ch31Caption("As يستخرج السبب.", "As extracts the cause.", english), Variables: []LessonVariable{{Name: "dns.IsNotFound", Value: "true"}}, Output: "true true"},
	})
}

var ch31Lesson315Arabic = []LessonStep{
	{ID: "c315.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch31Markdown("c315.mission", "1. اجمع عناوين من resolver محقون.\n2. صنّفها ورتّبها بشكل ثابت.\n3. ميّز غياب الاسم عن إجابة فاضية.")}},
	{ID: "c315.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch31IPv4("واجهة resolver هي حد الشبكة؛ الاختبارات ما تسألش DNS الجهاز.\n\n~net.IP~ بتتحول لتمثيل canonical. جرّب ~To4~ قبل ~To16~ عشان العنوان mapped يتحسب IPv4.\n\nإجابة عناوين فاضية نجاح؛ الغياب سبب خطأ مختلف.\n\nنفس ~context~ واسم الاستعلام يمشوا للنداءات بالترتيب.", false)}},
	{ID: "c315.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch31Question("c315.predict-as", "هل ~errors.As~ يستخرج DNS error؟\n\n~~~go\nfunc main() {\n\tcause := &net.DNSError{IsNotFound: true}\n\tvar dns *net.DNSError\n\tfmt.Println(errors.As(cause, &dns), dns.IsNotFound)\n}\n~~~", []string{"true true", "false false", "true false"}, 0, []string{
			"As يملأ المؤشر بالنوع المطابق.",
			"الخطأ من النوع المطلوب بالفعل.",
			"بعد النجاح، IsNotFound تفضل true.",
		}),
		ch31DNSCause("تابع استخراج النوع، وبعده راية الغياب.", false),
	}},
	{ID: "c315.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch31Markdown("c315.watch", "- خطأ الغياب في DNS مش زي قائمة عناوين فاضية.\n- رتّب النص canonical، وما تعدلش بيانات resolver.")}},
	{ID: "c315.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch31Markdown("c315.plan", "ابدأ بالمدخل والسياق، وبعده نداء العناوين ثم alias.\n\nصنّف كل IP، طبّع الاسم، واجمع النتيجة بعد نجاح كل المراحل. أي فشل يرجع نتيجة صفرية."),
		ch31Question("c315.address-order", "العنوان IPv4-mapped يتبع أي عائلة؟", []string{"IPv4", "IPv6", "عائلتين"}, 0, []string{
			"To4 يصنفه كعنوان IPv4.",
			"التمثيل المكتوب وحده لا يحدد الفئة.",
			"كل عنوان يظهر في عائلة واحدة.",
		}),
	}},
	{ID: "c315.turn-step", Title: "دورك", Blocks: []LessonBlock{ch31Markdown("c315.turn", "ابدأ بفصل أخطاء الإدخال عن أخطاء resolver.")}},
}

var ch31Lesson315English = []LessonStep{
	{ID: "c315.mission-step", Title: "Mission", Blocks: []LessonBlock{ch31Markdown("c315.mission", "1. Gather addresses through an injected resolver.\n2. Classify and order them consistently.\n3. Distinguish a missing name from an empty answer.")}},
	{ID: "c315.idea-step", Title: "Idea", Blocks: []LessonBlock{ch31IPv4("The resolver interface is the network boundary; tests never consult machine DNS.\n\nNormalize ~net.IP~ to canonical text. Try ~To4~ before ~To16~ so mapped addresses count as IPv4.\n\nAn empty address answer is success; name absence is a different error.\n\nPass the same context and query name to calls in order.", true)}},
	{ID: "c315.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch31Question("c315.predict-as", "Does ~errors.As~ extract this DNS error?\n\n~~~go\nfunc main() {\n\tcause := &net.DNSError{IsNotFound: true}\n\tvar dns *net.DNSError\n\tfmt.Println(errors.As(cause, &dns), dns.IsNotFound)\n}\n~~~", []string{"true true", "false false", "true false"}, 0, []string{
			"As fills the pointer with the matching type.",
			"The error already has the requested type.",
			"After extraction, IsNotFound remains true.",
		}),
		ch31DNSCause("Follow the type extraction, then its absence flag.", true),
	}},
	{ID: "c315.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch31Markdown("c315.watch", "- A DNS not-found error differs from an empty address list.\n- Sort canonical text, and do not mutate resolver data.")}},
	{ID: "c315.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch31Markdown("c315.plan", "Validate input and context first, then call addresses followed by alias.\n\nClassify each IP, normalize the name, and build the result only after every phase succeeds. Any failure returns a zero result."),
		ch31Question("c315.address-order", "Which family gets an IPv4-mapped address?", []string{"IPv4", "IPv6", "Both families"}, 0, []string{
			"To4 classifies it as an IPv4 address.",
			"Textual spelling alone does not decide the family.",
			"Each address belongs to one family.",
		}),
	}},
	{ID: "c315.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch31Markdown("c315.turn", "Start by separating input errors from resolver errors.")}},
}

const ch31JoinHostPortCode = `func main() {
	host := "::1"
	port := "7443"
	address := net.JoinHostPort(host, port)
	fmt.Println(address)
}`

func ch31JoinHostPort(text string, english bool) LessonBlock {
	return ch31Walkthrough("c316.join-host-port", text, ch31JoinHostPortCode, "e316-join-host-port", []string{"fmt", "net"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("العنوان IPv6 محلي.", "The IPv6 host is loopback.", english), Variables: []LessonVariable{{Name: "host", Value: "::1"}}},
		{Line: 3, Caption: ch31Caption("احفظ المنفذ كنص.", "Keep the port as text.", english), Variables: []LessonVariable{{Name: "host", Value: "::1"}, {Name: "port", Value: "7443"}}},
		{Line: 4, Caption: ch31Caption("JoinHostPort يضيف الأقواس.", "JoinHostPort adds brackets.", english), Variables: []LessonVariable{{Name: "address", Value: "[::1]:7443"}}},
		{Line: 5, Caption: ch31Caption("اطبع عنوان الاتصال.", "Print the network address.", english), Variables: []LessonVariable{{Name: "address", Value: "[::1]:7443"}}, Output: "[::1]:7443"},
	})
}

const ch31ResultSlotsCode = `func main() {
	ports := []string{"first", "second"}
	states := make([]string, len(ports))
	states[1] = "closed"
	states[0] = "open"
	fmt.Println(states)
}`

func ch31ResultSlots(text string, english bool) LessonBlock {
	return ch31Walkthrough("c316.result-slots", text, ch31ResultSlotsCode, "e316-result-slots", []string{"fmt"}, []LessonFrame{
		{Line: 0, Caption: ch31Caption("قبل التشغيل.", "Before execution.", english)},
		{Line: 2, Caption: ch31Caption("الفهرس يحفظ الترتيب.", "The index preserves order.", english)},
		{Line: 3, Caption: ch31Caption("خانة لكل عنصر.", "A slot for each item.", english), Variables: []LessonVariable{{Name: "states", Value: "[ ]"}}},
		{Line: 4, Caption: ch31Caption("اكتب نتيجة العنصر الثاني.", "Store the second result.", english), Variables: []LessonVariable{{Name: "states", Value: "[ closed]"}}},
		{Line: 5, Caption: ch31Caption("اكتب نتيجة الأول.", "Store the first result.", english), Variables: []LessonVariable{{Name: "states", Value: "[open closed]"}}},
		{Line: 6, Caption: ch31Caption("الطباعة تتبع القائمة.", "Output follows the input list.", english), Variables: []LessonVariable{{Name: "states", Value: "[open closed]"}}, Output: "[open closed]"},
	})
}

var ch31Lesson316Arabic = []LessonStep{
	{ID: "c316.mission-step", Title: "المهمة", Blocks: []LessonBlock{ch31Markdown("c316.mission", "1. افحص منافذ loopback المسموح بها.\n2. حدّ التوازي قبل تشغيل العمال.\n3. رجّع النتائج بترتيب المدخل.")}},
	{ID: "c316.idea-step", Title: "الفكرة", Blocks: []LessonBlock{ch31ResultSlots("فحص المنفذ محاولة فتح بس؛ ما تبعتش بيانات أو تعمل retry.\n\nتحقق من العنوان الرقمي والخطة كلها قبل أي dial. استخدم callback محقون عشان الاختبار ما يفتحش شبكة حقيقية.\n\nالعامل يكتب في خانة index تخصه؛ ترتيب انتهاء المحاولات ما يغيّرش ترتيب النتائج.\n\nاقفل كل connection راجعة قبل اكتمال العامل.", false)}},
	{ID: "c316.try-step", Title: "جرّب", Blocks: []LessonBlock{
		ch31Question("c316.predict-address", "إيه صيغة عنوان IPv6 مع المنفذ؟\n\n~~~go\nfunc main() {\n\thost := \"::1\"\n\tport := \"7443\"\n\taddress := net.JoinHostPort(host, port)\n\tfmt.Println(address)\n}\n~~~", []string{"[::1]:7443", "::1:7443", "127.0.0.1:7443"}, 0, []string{
			"الدالة تضيف الأقواس حول IPv6.",
			"الصيغة دي ملتبسة من غير الأقواس.",
			"ده عنوان IPv4 مختلف.",
		}),
		ch31JoinHostPort("تابع تركيب عنوان الشبكة قبل طباعته.", false),
	}},
	{ID: "c316.watch-step", Title: "خد بالك", Blocks: []LessonBlock{ch31Markdown("c316.watch", "- ارفض الخطة كلها قبل العمال أو أي اتصال.\n- استخدم ~IsLoopback~؛ prefix نصي مش تحقق عنوان.")}},
	{ID: "c316.check-step", Title: "فكّر واختبر", Blocks: []LessonBlock{
		ch31Markdown("c316.plan", "افصل فشل الإعداد عن فشل منفذ: الأول يوقف العملية، والتاني يسجل في خانته.\n\nراجع عدد العمال النشطين، وقفل الاتصالات، ثم ترتيب النتيجة مقارنة بالمدخل."),
		ch31Question("c316.worker-bound", "إيه اللي يحدد أقصى عمل متزامن؟", []string{"عدد workers المسموح", "عدد كل المنافذ", "سرعة الساعة"}, 0, []string{
			"حد العمال يقيّد المحاولات المتزامنة.",
			"إطلاق عامل لكل منفذ قد يتجاوز الحد.",
			"الساعة لا تتحكم في التوازي.",
		}),
	}},
	{ID: "c316.turn-step", Title: "دورك", Blocks: []LessonBlock{ch31Markdown("c316.turn", "ابدأ بالتحقق قبل العمال، ثم راجع كل نتيجة في مكانها.")}},
}

var ch31Lesson316English = []LessonStep{
	{ID: "c316.mission-step", Title: "Mission", Blocks: []LessonBlock{ch31Markdown("c316.mission", "1. Check permitted loopback ports.\n2. Bound concurrency before starting workers.\n3. Return results in input order.")}},
	{ID: "c316.idea-step", Title: "Idea", Blocks: []LessonBlock{ch31ResultSlots("A port check is only a dial attempt; send no data and do not retry.\n\nValidate the numeric address and whole plan before dialing. Inject a callback so tests never open a real network connection.\n\nEach worker writes to its own index slot; completion order cannot change result order.\n\nClose every returned connection before its worker finishes.", true)}},
	{ID: "c316.try-step", Title: "Try it", Blocks: []LessonBlock{
		ch31Question("c316.predict-address", "What is the IPv6 host-port form?\n\n~~~go\nfunc main() {\n\thost := \"::1\"\n\tport := \"7443\"\n\taddress := net.JoinHostPort(host, port)\n\tfmt.Println(address)\n}\n~~~", []string{"[::1]:7443", "::1:7443", "127.0.0.1:7443"}, 0, []string{
			"The function brackets the IPv6 address.",
			"That form is ambiguous without brackets.",
			"That is a different IPv4 address.",
		}),
		ch31JoinHostPort("Follow address construction before it is printed.", true),
	}},
	{ID: "c316.watch-step", Title: "Watch out", Blocks: []LessonBlock{ch31Markdown("c316.watch", "- Reject the whole plan before workers or any connection.\n- Use ~IsLoopback~; a text prefix is not address validation.")}},
	{ID: "c316.check-step", Title: "Think & check", Blocks: []LessonBlock{
		ch31Markdown("c316.plan", "Separate invalid configuration from a failed port: one stops the operation, the other fills its result slot.\n\nReview active worker count, connection closure, and result order against the input."),
		ch31Question("c316.worker-bound", "What sets the maximum concurrent work?", []string{"The allowed worker count", "The total port count", "Clock speed"}, 0, []string{
			"The worker limit bounds simultaneous attempts.",
			"One worker per port can exceed the bound.",
			"The clock does not control concurrency.",
		}),
	}},
	{ID: "c316.turn-step", Title: "Your turn", Blocks: []LessonBlock{ch31Markdown("c316.turn", "Start with validation before workers, then check each result in its slot.")}},
}
