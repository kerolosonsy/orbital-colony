package levels

import "orbital/colony/i18n"

func init() {
	add(&Level{ID: "31-1-echo", Chapter: 31, Kind: "lab", Title: "الصدى", Goal: "رجّع سطور منارة الاستغاثة بحد بايتات وموعد نهائي",
		Requires: []string{"sel:net.Conn"}, Bridge: []int{71}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch311Starter}},
		Lesson: `# صدى منارة الاستغاثة

## قبل ما تكتب
- هل نداء قراية واحد معناه رسالة كاملة؟
- مين يقفل الاتصال لو السطر أطول من المسموح؟

## TCP تيار بايتات
~net.Conn~ اتصال بيطبّق ~io.Reader~ و ~io.Writer~ وبيوفّر ~Close~. TCP بيحفظ ترتيب البايتات، لكن ممكن يقسّم رسالة على كذا قراية أو يجمع رسالتين. حدود الرسالة مسؤولية البروتوكول.
~~~go
func send(c net.Conn) error {
 _, err := io.WriteString(c, "dock\n")
 return err
}
~~~
~net.Listen~ بيعمل ~net.Listener~؛ ~Accept~ يستنى اتصال جديد، و ~Addr~ يجيب العنوان الفعلي لما البورت 0 يختاره النظام. قفل المستمع يفك ~Accept~، لكنه ميقفلش الاتصالات المقبولة. صاحب كل اتصال يقفله. المهمة بتستقبل اتصال جاهز عشان تختبره من غير مستمع.
~~~go
ln, err := net.Listen("tcp", "127.0.0.1:0")
if err != nil { return err }
defer ln.Close()
c, err := ln.Accept()
if err != nil { return err }
defer c.Close()
fmt.Println(ln.Addr())
~~~

## حدود السطور والذاكرة
الفاصل هنا البايت ~'\n'~. ~bufio.NewReader~ بيخزن قرايات صغيرة؛ ~ReadByte~ بيرجع بايت أو خطأ. اجمع لحد الفاصل وافحص الحجم قبل الإضافة؛ متستخدمش قراية سطر غير محدودة. الحد بيحسب بايتات المحتوى، مش حروف Unicode، ومش بيحسب الفاصل. ~'\r'~ بايت عادي هنا.
~~~go
r := bufio.NewReader(source)
b, err := r.ReadByte()
if err != nil { return err }
if b == '\n' { /* انتهت الرسالة */ }
~~~
~io.EOF~ قبل أي محتوى جديد نهاية نظيفة. بعد محتوى من غير فاصل يبقى ~io.ErrUnexpectedEOF~، ومترجعش المحتوى الجزئي. ~errors.Is(err, io.EOF)~ يفحص السبب.

## موعد نهائي محقون
~SetDeadline~ بيحدد موعد مطلق للقراية والكتابة؛ ~time.Time{}~ بيلغيه. مش مدة لكل عملية. هنا الموعد محقون، وممنوع ~time.Now~ أو نوم جوه المنطق. خطأ الموعد أو القراية أو الكتابة لازم يرجع بنفس السبب. خطأ شبكة ممكن يتفحص كـ ~net.Error~ وله ~Timeout()~.
~~~go
err := c.SetDeadline(deadline)
if err != nil { return err }
var ne net.Error
if errors.As(err, &ne) && ne.Timeout() { /* مهلة */ }
~~~
~Write~ ممكن ترجع عدد أقل من المطلوب. كمّل الباقي؛ لو رجعت 0 من غير خطأ ارجع ~io.ErrShortWrite~ عشان متلفش للأبد.
~~~go
n, err := sink.Write(payload)
if err != nil { return err }
if n == 0 { return io.ErrShortWrite }
payload = payload[n:]
~~~

## اختبار من غير شبكة
~net.Pipe~ بيرجع اتصالين في الذاكرة؛ الكتابة بتستنى القراية من الطرف التاني. شغّل طرف في goroutine واقفل الاتنين. ~labkit.LeakCheck(t)~ يتنادى قبل تشغيلهم. التستات بتحقن أخطاء المهلة بدل انتظار ساعة الجهاز.
~~~go
left, right := net.Pipe()
defer left.Close()
defer right.Close()
~~~

## أخطاء شائعة
~~~go
buf := make([]byte, 9000)
n, _ := c.Read(buf)
c.Write(buf[:n]) // غلط: قراية TCP مش بالضرورة سطر كامل، والخطأ اتبلع
~~~

## المهمة
نفّذ ~EchoBeacon(c net.Conn, limit int, deadline time.Time) error~. الاتصال غير ~nil~ وملكيته بتتنقل للدالة: اقفله مرة عند كل خروج، وتجاهل خطأ القفل فقط. ارفض ~limit~ خارج 0 لحد 4096 قبل أي I/O غير القفل. نادِ ~SetDeadline(deadline)~ مرة قبل القراية. كل سطر كامل يرجع زي ما هو بما فيه الفاصل، بالترتيب، حتى السطر الفاضي. بالظبط ~limit~ بايت محتوى مسموح، والبايت الزيادة يرجع ~ErrBeaconLine~ من غير رد على السطر ده. الردود السابقة تفضل مكتوبة.

نهاية نظيفة ترجع ~nil~؛ نهاية وسط سطر ترجع ~io.ErrUnexpectedEOF~. أخطاء النقل والموعد تتحفظ بـ ~errors.Is~. اتعامل مع الكتابة الجزئية ومن غير تقدم زي الشرح. مفيش goroutines داخل الدالة أو ~Report~؛ ~main~ فاضية.
`,
		Hints: []string{"ابدأ بملكية الاتصال والتحقق، وبعدها اضبط الموعد المحقون.", "اقرا بايت بايت؛ افحص الحد قبل إضافة محتوى جديد.", "عند الفاصل اكتب السطر كله بحلقة للكتابة الجزئية، وبعدها صفّر السطر."},
		En: &LevelText{Title: "Echo", Goal: "Echo distress-beacon lines with a byte limit and deadline", Files: []LabFile{{Name: "main.go", Starter: ch311Starter}},
			Lesson: `# Distress Beacon Echo

## Before you write
- Does one read mean one complete message?
- Who closes a connection when a line exceeds its limit?

## TCP is a byte stream
~net.Conn~ implements ~io.Reader~, ~io.Writer~, and ~Close~. TCP preserves byte order but can split one message across reads or combine messages. Framing belongs to your protocol.
~~~go
func send(c net.Conn) error {
 _, err := io.WriteString(c, "dock\n")
 return err
}
~~~
~net.Listen~ creates a ~net.Listener~. ~Accept~ waits for a connection; ~Addr~ reports the assigned address when port zero lets the system choose. Closing a listener unblocks ~Accept~, but does not close accepted connections. Each connection needs an owner. This mission injects an accepted connection.
~~~go
ln, err := net.Listen("tcp", "127.0.0.1:0")
if err != nil { return err }
defer ln.Close()
c, err := ln.Accept()
if err != nil { return err }
defer c.Close()
fmt.Println(ln.Addr())
~~~

## Bounded line framing
The delimiter is byte ~'\n'~. ~bufio.NewReader~ buffers reads; ~ReadByte~ returns a byte or error. Accumulate up to the delimiter, checking size before appending; avoid unbounded line reads. The limit counts content bytes, not Unicode characters, and excludes the delimiter. A carriage return is an ordinary content byte.
~~~go
r := bufio.NewReader(source)
b, err := r.ReadByte()
if err != nil { return err }
if b == '\n' { /* message ends */ }
~~~
~io.EOF~ before new content is clean closure. After unterminated content, return ~io.ErrUnexpectedEOF~ without echoing the partial line. ~errors.Is(err, io.EOF)~ checks the cause.

## Injected absolute deadline
~SetDeadline~ sets one absolute deadline for reading and writing; ~time.Time{}~ clears it. It is not a duration per operation. Receive it from the caller; do not call ~time.Now~ or sleep. Preserve deadline, read, and write failures. Network errors can implement ~net.Error~ with ~Timeout()~.
~~~go
err := c.SetDeadline(deadline)
if err != nil { return err }
var ne net.Error
if errors.As(err, &ne) && ne.Timeout() { /* timeout */ }
~~~
~Write~ may accept fewer bytes. Continue with the remainder; zero progress without error becomes ~io.ErrShortWrite~ to avoid an infinite loop.
~~~go
n, err := sink.Write(payload)
if err != nil { return err }
if n == 0 { return io.ErrShortWrite }
payload = payload[n:]
~~~

## In-memory testing
~net.Pipe~ supplies two in-memory connections. Writes wait for peer reads. Run one end in a goroutine and close both. Call ~labkit.LeakCheck(t)~ before starting them. Tests inject timeout errors rather than wait for wall-clock expiry.
~~~go
left, right := net.Pipe()
defer left.Close()
defer right.Close()
~~~

## Common mistakes
~~~go
buf := make([]byte, 9000)
n, _ := c.Read(buf)
c.Write(buf[:n]) // Wrong: one TCP read need not contain a line; errors are lost
~~~

## Your mission
Implement ~EchoBeacon(c net.Conn, limit int, deadline time.Time) error~. The non-~nil~ connection transfers ownership: close exactly once on every exit, ignoring only its close error. Reject limits outside 0 through 4096 before any I/O except closing. Call ~SetDeadline(deadline)~ once before reading. Echo each complete line unchanged including its delimiter, in order, including empty lines. Exactly ~limit~ content bytes are valid; one more returns ~ErrBeaconLine~ without echoing that line. Earlier responses remain written.

Clean closure returns ~nil~; EOF within a line returns ~io.ErrUnexpectedEOF~. Preserve transport and deadline causes for ~errors.Is~. Handle partial writes and zero progress as explained. No internal goroutines or ~Report~; leave ~main~ empty.
`,
			Hints: []string{"Own and validate the connection, then set the injected deadline.", "Read bytes and check the limit before appending content.", "At a delimiter write the whole line using a partial-write loop, then reset it."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch311Hidden},
	})
}

const ch311Starter = `package main

import (
	"errors"
	"net"
	"time"
)

var ErrBeaconLine = errors.New("beacon line too long")

func EchoBeacon(c net.Conn, limit int, deadline time.Time) error {
	defer c.Close()
	return errors.New("TODO")
}
func main() {}
`
const ch311Hidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"orbital/colony/labkit"
	"testing"
	"time"
)

type beaconConn struct {
	net.Conn
	input                     io.Reader
	output                    bytes.Buffer
	deadline                  time.Time
	sets, closes, chunk       int
	setErr, readErr, writeErr error
	stalled                   bool
}

func (c *beaconConn) Read(p []byte) (int, error) {
	if c.readErr != nil {
		return 0, c.readErr
	}
	return c.input.Read(p)
}
func (c *beaconConn) Write(p []byte) (int, error) {
	if c.writeErr != nil {
		return 0, c.writeErr
	}
	if c.stalled {
		return 0, nil
	}
	if c.chunk > 0 && len(p) > c.chunk {
		p = p[:c.chunk]
	}
	return c.output.Write(p)
}
func (c *beaconConn) SetDeadline(d time.Time) error { c.sets++; c.deadline = d; return c.setErr }
func (c *beaconConn) Close() error                  { c.closes++; return nil }

// Boundaries: empty, exact/over limit, partial EOF, byte lengths, partial/no-progress writes, failures.
func TestHiddenBeaconLines(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(311))
	deadline := time.Unix(900, 17)
	for i := 0; i < 50; i++ {
		limit := rng.Intn(80)
		var wire bytes.Buffer
		for j := 0; j < 1+rng.Intn(8); j++ {
			wire.Write(bytes.Repeat([]byte{byte(32 + rng.Intn(90))}, rng.Intn(limit+1)))
			wire.WriteByte('\n')
		}
		c := &beaconConn{input: bytes.NewReader(wire.Bytes()), chunk: 1 + rng.Intn(5)}
		if err := EchoBeacon(c, limit, deadline); err != nil || !bytes.Equal(c.output.Bytes(), wire.Bytes()) || c.closes != 1 || c.sets != 1 || c.deadline != deadline {
			t.Fatal("echo", err, c)
		}
	}
	for _, tc := range []struct {
		wire  string
		limit int
		want  error
		out   string
	}{{"", 0, nil, ""}, {"\n", 0, nil, "\n"}, {"a\n", 0, ErrBeaconLine, ""}, {"ok\nlong\n", 2, ErrBeaconLine, "ok\n"}, {"x", 4, io.ErrUnexpectedEOF, ""}, {"é\n", 1, ErrBeaconLine, ""}, {"\r\n", 1, nil, "\r\n"}, {"\x00\xff\n", 2, nil, "\x00\xff\n"}, {string(bytes.Repeat([]byte{'z'}, 4096)) + "\n", 4096, nil, string(bytes.Repeat([]byte{'z'}, 4096)) + "\n"}} {
		c := &beaconConn{input: bytes.NewBufferString(tc.wire)}
		err := EchoBeacon(c, tc.limit, deadline)
		if !errors.Is(err, tc.want) || c.output.String() != tc.out || c.closes != 1 {
			t.Fatal("boundary", tc.limit, err)
		}
	}
	for _, limit := range []int{-1, 4097} {
		c := &beaconConn{input: bytes.NewReader(nil)}
		if EchoBeacon(c, limit, deadline) == nil || c.sets != 0 || c.closes != 1 {
			t.Fatal("validation")
		}
	}
}
func TestHiddenBeaconFailures(t *testing.T) {
	fault := fmt.Errorf("wrapped: %w", timeoutBeacon{})
	for _, stage := range []string{"set", "read", "write", "stall"} {
		c := &beaconConn{input: bytes.NewBufferString("a\n")}
		want := error(fault)
		switch stage {
		case "set":
			c.setErr = fault
		case "read":
			c.readErr = fault
		case "write":
			c.writeErr = fault
		case "stall":
			c.stalled = true
			want = io.ErrShortWrite
		}
		if err := EchoBeacon(c, 10, time.Time{}); !errors.Is(err, want) || c.closes != 1 {
			t.Fatal(stage, err)
		}
	}
}

type timeoutBeacon struct{}

func (timeoutBeacon) Error() string   { return "injected timeout" }
func (timeoutBeacon) Timeout() bool   { return true }
func (timeoutBeacon) Temporary() bool { return false }
func TestHiddenBeaconPipe(t *testing.T) {
	labkit.LeakCheck(t)
	left, right := net.Pipe()
	t.Cleanup(func() { left.Close(); right.Close() })
	done := make(chan error, 1)
	go func() { done <- EchoBeacon(left, 12, time.Time{}) }()
	for _, part := range []string{"air", "lock", "\n"} {
		if _, err := io.WriteString(right, part); err != nil {
			t.Fatal(err)
		}
	}
	reply := make([]byte, 8)
	if _, err := io.ReadFull(right, reply); err != nil || string(reply) != "airlock\n" {
		t.Fatal("pipe", err, string(reply))
	}
	right.Close()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("echo did not stop")
	}
}
`

func init() {
	add(&Level{ID: "31-2-chat", Chapter: 31, Kind: "lab", Title: "غرفة الراديو", Goal: "شغّل قلب غرفة الراديو بمالك واحد وعزل العميل البطيء",
		Requires: []string{"chan", "select"}, Bridge: []int{72}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch312Starter}},
		Lesson: `# غرفة راديو المستعمرة

## قبل ما تكتب
- مين يملك خريطة الأسماء وقفل قنوات الخروج؟
- هل عميل مش بيقرا يقدر يعطّل كل الطاقم؟

## افصل النقل عن الغرفة
المهمة السابقة شرحت السطور والاتصال. هنا هنبني قلب الغرفة بس؛ محوّل الاتصال ممكن يقرا السطر ويبعت أمر، وكاتب مستقل يصرف صندوق العميل. الاختبار بيحقن القنوات عشان التزامن يبقى قابل للتحكم، من غير TCP جديد.

~RadioCommand~ أمر فيه ~Kind~ و ~Name~ و ~Text~، وقناة ~Out~ للانضمام وقناة ~Reply~ لنتيجة كل أمر. ~RunRadio~ حلقة واحدة تملك الخريطة؛ مفيش أي goroutine تانية بتقراها أو تكتبها.
~~~go
ack := make(chan error, 1)
mail := make(chan string, 4)
commands <- RadioCommand{Kind: "join", Name: "medic", Out: mail, Reply: ack}
err := <-ack
~~~
المنادي يملك قناة الأوامر وقفلها، ويضمن ~Reply~ فاضية بسعة 1. بعد قبول الانضمام، الغرفة تملك إرسال وقفل ~Out~، والمنادي يقرا بس. القناة جديدة وفاضية ومش مشتركة بين عضوين؛ متقفلهاش برّه الغرفة.

## ضغط عكسي من غير تعطيل
الإرسال غير المعطّل يستخدم ~select~ مع ~default~. امتلاء صندوق عضو معناه فصله وقفل صندوقه، مش إسقاط رسالة مع إبقائه متصل. الرسائل الموجودة تفضل قابلة للقراية بعد القفل.
~~~go
select {
case mailbox <- "alarm":
 default:
 close(mailbox)
 delete(members, "medic")
}
~~~
البث من عضو معروف يوصل لكل عضو تاني بصيغة ~name + ": " + text~. المرسل مبيستقبلش نسخته. الانضمام والخروج ميبعتوش إعلانات. رد الأمر بعد انتهاء التحديث والبث بيمثل نقطة مزامنة للاختبار.
~~~go
commands <- RadioCommand{Kind: "say", Name: "medic", Text: "ready", Reply: ack}
err := <-ack
~~~

~strings.ContainsAny~ تكشف أي بايت ممنوع في النص؛ افحص الطول بالبايتات قبل البث.
~~~go
if strings.ContainsAny(text, "\r\n") { return ErrRadioCommand }
~~~

## نهاية لها مالك
إغلاق قناة الأوامر نهاية نظيفة، وإلغاء ~ctx~ يرجع ~ctx.Err()~. في الحالتين اقفل كل الصناديق المقبولة الباقية مرة واحدة. افحص الإلغاء قبل الدورة وداخل ~select~. لو أمر اتقبل قبل ملاحظة الإلغاء يجوز يكمل؛ مفيش وعد بأولوية ~select~.
~~~go
select {
case <-ctx.Done():
 return ctx.Err()
case command, ok := <-commands:
 if !ok { return nil }
 _ = command
}
~~~

## أخطاء شائعة
~~~go
for _, mailbox := range members {
 mailbox <- text // غلط: عضو بطيء يوقف الحلقة كلها
}
~~~

## المهمة
نفّذ ~RunRadio(ctx context.Context, commands <-chan RadioCommand) error~؛ السياق غير ~nil~. استخدم حلقة مالك واحدة و ~chan~ و ~select~. الأسماء من 1 لحد 16 بايت من ~a-z~ و ~0-9~ و ~-~ فقط. ~join~ يرفض الاسم البايظ أو المكرر أو ~Out~ من غير سعة بـ ~ErrRadioName~؛ الرفض مبيقفلش ~Out~ ولا يغير الحالة. اسم مفصول ينفع يرجع بصندوق جديد.

~leave~ يحذف العضو ويقفل صندوقه؛ اسم مش موجود يرجع ~ErrRadioMember~. ~say~ يتأكد من العضوية الأول، وبعدها يقبل نص من 1 لحد 128 بايت من غير ~\n~ أو ~\r~؛ النص البايظ أو نوع أمر مجهول يرجع ~ErrRadioCommand~. البث لغير المرسل فقط، وبنفس ترتيب معالجة الأوامر لكل مستقبل. افصل الصندوق الممتلي من غير تعطيل الآخرين. كل أمر مقبول من القناة ياخد رد واحد حتى لو فشل، بعد اكتمال أثره؛ قنوات ~Reply~ مضمونة غير ~nil~ وفاضية بسعة 1، ومش بتتقفل هنا. ~Out~ المقبولة مضمونة منفصلة وملك الغرفة. النهاية تقفل كل الصناديق، ومفيش goroutines بتعيش بعد الرجوع. ~main~ فاضية ومفيش ~Report~.

النجمة التالتة: حافظ على نفس العقد مع منتجين متزامنين للأوامر؛ كل رسالة من كل منتج توصل مرة واحدة، من غير فرض ترتيب بين المنتجين.
`,
		Hints: []string{"خلي الخريطة محلية للحلقة، والقفل في مسار خروج واحد.", "رد join بعد الفحص؛ leave يمسح قبل القفل، وsay يتأكد من العضوية.", "اعمل إرسال select/default لكل مستقبل؛ احذف واقفل البطيء، وبعدها ابعت رد الأمر."},
		En: &LevelText{Title: "Radio Room", Goal: "Run a single-owner radio hub that isolates slow crew members", Files: []LabFile{{Name: "main.go", Starter: ch312Starter}},
			Lesson: `# Colony Radio Room

## Before you write
- Who owns the name map and closes outgoing queues?
- Can a crew member who stops reading block everyone?

## Separate transport from the room
The previous mission taught connections and lines. Here build only the room core. A transport adapter could read a line into a command and drain a member's queue independently. Injected channels make concurrency controllable without new TCP code.

~RadioCommand~ contains ~Kind~, ~Name~, ~Text~, an ~Out~ queue for joining, and a ~Reply~ acknowledgement channel. One ~RunRadio~ loop owns the map; no other goroutine accesses it.
~~~go
ack := make(chan error, 1)
mail := make(chan string, 4)
commands <- RadioCommand{Kind: "join", Name: "medic", Out: mail, Reply: ack}
err := <-ack
~~~
The caller owns and closes the commands channel and guarantees an empty capacity-one ~Reply~. On successful join the room owns sending to and closing ~Out~; the caller only reads it. It is a fresh empty queue, never shared by two members. Never close it outside the room.

## Backpressure without blocking
Use ~select~ with ~default~ for nonblocking delivery. A full mailbox causes disconnection and closure, rather than keeping a member after dropping a message. Already queued messages remain readable after closure.
~~~go
select {
case mailbox <- "alarm":
 default:
 close(mailbox)
 delete(members, "medic")
}
~~~
A known member broadcasts to everyone else as ~name + ": " + text~. The sender receives no copy. Joining and leaving emit no announcements. Acknowledgement happens after all state changes and broadcasts, giving tests a synchronization point.
~~~go
commands <- RadioCommand{Kind: "say", Name: "medic", Text: "ready", Reply: ack}
err := <-ack
~~~

~strings.ContainsAny~ detects forbidden characters in text; check byte length before broadcasting.
~~~go
if strings.ContainsAny(text, "\r\n") { return ErrRadioCommand }
~~~

## Owned termination
Closing commands ends successfully; cancellation returns ~ctx.Err()~. Both close every remaining accepted mailbox exactly once. Check cancellation before the loop's selection and within ~select~. An accepted command may finish before cancellation is observed; ~select~ has no priority guarantee.
~~~go
select {
case <-ctx.Done():
 return ctx.Err()
case command, ok := <-commands:
 if !ok { return nil }
 _ = command
}
~~~

## Common mistakes
~~~go
for _, mailbox := range members {
 mailbox <- text // Wrong: one slow member blocks the entire loop
}
~~~

## Your mission
Implement ~RunRadio(ctx context.Context, commands <-chan RadioCommand) error~ with a non-~nil~ context. Use one owner loop, ~chan~, and ~select~. Names contain 1 through 16 bytes drawn only from ~a-z~, ~0-9~, and ~-~. A ~join~ rejects malformed or duplicate names or an ~Out~ without capacity with ~ErrRadioName~, leaving state and rejected queues unchanged. A disconnected name may rejoin with a fresh queue.

A ~leave~ removes and closes its member queue; absent names return ~ErrRadioMember~. A ~say~ checks membership first, then requires 1 through 128 text bytes without ~\n~ or ~\r~; invalid text or an unknown command kind returns ~ErrRadioCommand~. Broadcast only to other members, preserving command processing order per recipient. Disconnect full mailboxes without blocking others. Every command received gets exactly one reply after its effects, including failures. ~Reply~ is guaranteed non-~nil~, empty, capacity one; do not close it. Accepted ~Out~ queues are distinct and owned by the room. Close all remaining mailboxes on termination; no goroutines survive return. Empty ~main~; no ~Report~.

Third star: preserve the contract under concurrent command producers. Every producer's message arrives exactly once, with no specified order between producers.
`,
			Hints: []string{"Keep the map local to the loop and close remaining queues in one exit path.", "Validate before joining; remove before closing; check membership before speaking.", "Use select/default per recipient; evict slow members, then acknowledge."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch312Hidden},
	})
}

const ch312Starter = `package main

import (
	"context"
	"errors"
)

var ErrRadioName = errors.New("invalid radio name")
var ErrRadioMember = errors.New("unknown radio member")
var ErrRadioCommand = errors.New("invalid radio command")

type RadioCommand struct {
	Kind, Name, Text string
	Out              chan string
	Reply            chan error
}

func RunRadio(ctx context.Context, commands <-chan RadioCommand) error { return errors.New("TODO") }
func main()                                                            {}
`
const ch312Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"orbital/colony/labkit"
	"strings"
	"sync"
	"testing"
	"time"
)

func radioWait[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("radio deadlock")
		var zero T
		return zero
	}
}
func radioStart(t *testing.T) (chan RadioCommand, context.CancelFunc, <-chan error) {
	t.Helper()
	labkit.LeakCheck(t)
	ctx, cancel := context.WithCancel(context.Background())
	commands := make(chan RadioCommand)
	done := make(chan error, 1)
	go func() { done <- RunRadio(ctx, commands) }()
	t.Cleanup(cancel)
	return commands, cancel, done
}
func radioSend(t *testing.T, ch chan RadioCommand, c RadioCommand) error {
	t.Helper()
	c.Reply = make(chan error, 1)
	select {
	case ch <- c:
	case <-time.After(3 * time.Second):
		t.Fatal("blocked command")
	}
	return radioWait(t, c.Reply)
}
func radioClosed(t *testing.T, ch chan string) {
	t.Helper()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("expected closed empty mailbox")
		}
	default:
		t.Fatal("mailbox left open")
	}
}

// Boundaries: invalid/duplicate/reused names, empty/oversize text, slow peer, cancellation, concurrent producers.
func TestHiddenRadioBroadcastAndEviction(t *testing.T) {
	commands, _, done := radioStart(t)
	fast := make(chan string, 64)
	slow := make(chan string, 1)
	sender := make(chan string, 64)
	for name, out := range map[string]chan string{"pilot": sender, "engineer": fast, "medic": slow} {
		if err := radioSend(t, commands, RadioCommand{Kind: "join", Name: name, Out: out}); err != nil {
			t.Fatal(err)
		}
	}
	rejected := make(chan string, 1)
	if !errors.Is(radioSend(t, commands, RadioCommand{Kind: "join", Name: "pilot", Out: rejected}), ErrRadioName) {
		t.Fatal("duplicate")
	}
	rng := rand.New(rand.NewSource(312))
	for i := 0; i < 30; i++ {
		text := fmt.Sprint(rng.Intn(100000))
		if err := radioSend(t, commands, RadioCommand{Kind: "say", Name: "pilot", Text: text}); err != nil {
			t.Fatal(err)
		}
		if got := radioWait(t, fast); got != "pilot: "+text {
			t.Fatal("broadcast", got)
		}
	}
	if got := radioWait(t, slow); !strings.HasPrefix(got, "pilot: ") {
		t.Fatal(got)
	}
	radioClosed(t, slow)
	if !errors.Is(radioSend(t, commands, RadioCommand{Kind: "say", Name: "medic", Text: "hi"}), ErrRadioMember) {
		t.Fatal("slow member retained")
	}
	if len(sender) != 0 || len(rejected) != 0 {
		t.Fatal("self broadcast or rejected member")
	}
	rejoin := make(chan string, 1)
	if radioSend(t, commands, RadioCommand{Kind: "join", Name: "medic", Out: rejoin}) != nil {
		t.Fatal("reuse")
	}
	if radioSend(t, commands, RadioCommand{Kind: "leave", Name: "medic"}) != nil {
		t.Fatal("leave")
	}
	radioClosed(t, rejoin)
	close(commands)
	if err := radioWait(t, done); err != nil {
		t.Fatal(err)
	}
	radioClosed(t, fast)
	radioClosed(t, sender)
	select {
	case _, ok := <-rejected:
		if !ok {
			t.Fatal("rejected channel closed")
		}
	default:
	}
}
func TestHiddenRadioValidationAndCancel(t *testing.T) {
	commands, cancel, done := radioStart(t)
	for _, name := range []string{"", "A", "white space", "é", strings.Repeat("a", 17)} {
		if !errors.Is(radioSend(t, commands, RadioCommand{Kind: "join", Name: name, Out: make(chan string, 1)}), ErrRadioName) {
			t.Fatal("name", name)
		}
	}
	for _, out := range []chan string{nil, make(chan string)} {
		if !errors.Is(radioSend(t, commands, RadioCommand{Kind: "join", Name: "pilot", Out: out}), ErrRadioName) {
			t.Fatal("capacity")
		}
	}
	mail := make(chan string, 2)
	if radioSend(t, commands, RadioCommand{Kind: "join", Name: strings.Repeat("a", 16), Out: mail}) != nil {
		t.Fatal("exact name")
	}
	for _, text := range []string{"", strings.Repeat("z", 129), "x\ny", "x\ry"} {
		if !errors.Is(radioSend(t, commands, RadioCommand{Kind: "say", Name: strings.Repeat("a", 16), Text: text}), ErrRadioCommand) {
			t.Fatal("text")
		}
	}
	if radioSend(t, commands, RadioCommand{Kind: "say", Name: strings.Repeat("a", 16), Text: strings.Repeat("z", 128)}) != nil {
		t.Fatal("exact text")
	}
	for _, kind := range []string{"say", "leave"} {
		if !errors.Is(radioSend(t, commands, RadioCommand{Kind: kind, Name: "ghost"}), ErrRadioMember) {
			t.Fatal("missing")
		}
	}
	if !errors.Is(radioSend(t, commands, RadioCommand{Kind: "unknown"}), ErrRadioCommand) {
		t.Fatal("kind")
	}
	cancel()
	if !errors.Is(radioWait(t, done), context.Canceled) {
		t.Fatal("cancel")
	}
	radioClosed(t, mail)
}
func TestHiddenBonusRadioConcurrentProducers(t *testing.T) {
	commands, _, done := radioStart(t)
	observer := make(chan string, 256)
	if radioSend(t, commands, RadioCommand{Kind: "join", Name: "observer", Out: observer}) != nil {
		t.Fatal("observer")
	}
	for i := 0; i < 4; i++ {
		if radioSend(t, commands, RadioCommand{Kind: "join", Name: fmt.Sprint("crew-", i), Out: make(chan string, 256)}) != nil {
			t.Fatal("crew")
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 40; j++ {
				if radioSend(t, commands, RadioCommand{Kind: "say", Name: fmt.Sprint("crew-", i), Text: fmt.Sprint(j)}) != nil {
					t.Error("say")
				}
			}
		}(i)
	}
	wg.Wait()
	close(commands)
	if radioWait(t, done) != nil {
		t.Fatal("stop")
	}
	seen := map[string]bool{}
	for msg := range observer {
		if seen[msg] {
			t.Fatal("duplicate")
		}
		seen[msg] = true
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 40; j++ {
			if !seen[fmt.Sprintf("crew-%d: %d", i, j)] {
				t.Fatal("missing")
			}
		}
	}
}
`

func init() {
	add(&Level{ID: "31-3-discovery", Chapter: 31, Kind: "lab", Title: "اكتشاف الجيران", Goal: "رد على حزمة تعريف جار واحدة ببروتوكول UDP محدود",
		Requires: []string{"sel:net.ListenPacket"}, Bridge: []int{73}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch313Starter}},
		Lesson: `# حزمة تعريف من رصيف قريب

## قبل ما تكتب
- هتعرف الحزمة الأطول من الحد إزاي من غير قبول جزء منها؟
- الرد يروح لعنوان مكتوب في الحزمة ولا عنوان المرسل الفعلي؟

## UDP مش تيار
UDP بيحافظ على حدود الحزمة؛ ممكن حزمة تضيع أو تتكرر أو توصل بترتيب مختلف. نجاح الإرسال مش إقرار بوصولها. هنا عملية واحدة لكل حزمة من غير retries.

~net.ListenPacket~ بتفتح ~net.PacketConn~ على ~udp4~؛ ~LocalAddr~ يرجع العنوان الفعلي. افصل الفتح عن الرد عشان نحقن transport في التست.
~~~go
pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
if err != nil { return err }
defer pc.Close()
fmt.Println(pc.LocalAddr())
~~~
~OpenNeighbor~ بتاخد دالة فتح بنفس التوقيع، ولو ~nil~ تستخدم ~net.ListenPacket~. مثال حقن: ~OpenNeighbor(net.ListenPacket)~. متضيفش عنوان من المستخدم؛ المصنع مربوط بالـ loopback.

## قراية حزمة والرد لمصدرها
~ReadFrom~ ترجع عدد البايتات وعنوان ~net.Addr~ للمرسل، و ~WriteTo~ تبعت حزمة واحدة له. ~SetDeadline~ نفس الموعد المطلق اللي اتشرح في الصدى. استخدم مخزن 65 بايت لحد 64: لو اتقرت 65 يبقى أطول من الحد حتى لو باقي الحزمة اترمى. متعيدش إرسال جزء من حزمة عند short write؛ ده مش TCP.
~~~go
buf := make([]byte, 65)
n, peer, err := pc.ReadFrom(buf)
if err != nil { return err }
if n > 64 { return ErrNeighborPacket }
n, err = pc.WriteTo([]byte("notice"), peer)
if err == nil && n != len("notice") { return io.ErrShortWrite }
~~~
~net.PacketConn~ ينفع يتنفّذ في الذاكرة لاختبار المهلة وخطأ القراية والكتابة. في تشغيل محلي، المنادي يقفل الـ socket اللي المصنع رجّعه؛ دالة الرد بتستعيره ومبتقفلوش.

~strings.ContainsRune~ تفحص هل حرف موجود في مجموعة مسموحة؛ فحص طول النص لوحده مش كفاية.
~~~go
allowed := strings.ContainsRune("0123456789abcdef", rune('c'))
fmt.Println(allowed) // true
~~~

## نسخة البروتوكول وربط الرد
الطلب نص ASCII بالظبط ~Q1|nonce~ والرد ~A1|nonce|station~. ~1~ نسخة البروتوكول، والـ ~nonce~ كود ربط محقون مش عشوائية جوه المنطق. الكود من 1 لحد 16 بايت من ~0-9a-f~، والمحطة من 1 لحد 16 بايت من ~a-z0-9-~. مفيش مسافات أو فاصل أخير. ~strings.Split~ يكشف الحقول الزيادة.
~~~go
parts := strings.Split("A1|ab4|medic", "|")
if len(parts) != 3 { return "", ErrNeighborPacket }
~~~
~EncodeNeighborProbe("ab4")~ يرجع طلب، و ~DecodeNeighborReply(packet, "ab4")~ يرجع المحطة بعد فحص النسخة والكود. أي اختلاف كود يرفض الرد، حتى لو باقي الحقول سليمة.

## أخطاء شائعة
~~~go
buf := make([]byte, 64)
n, peer, _ := pc.ReadFrom(buf)
// غلط: حزمة 100 بايت ممكن تظهر كأنها 64 سليمة
pc.WriteTo(buf[:n], peer)
~~~

## المهمة
نفّذ ~OpenNeighbor(listen func(string,string)(net.PacketConn,error)) (net.PacketConn,error)~: اختار ~net.ListenPacket~ لو الدالة ~nil~، ونادِها مرة بـ ~udp4~ و ~127.0.0.1:0~ ورجّع نفس الناتج والخطأ. نفّذ ~EncodeNeighborProbe(nonce string) ([]byte,error)~ و ~DecodeNeighborReply(packet []byte, nonce string) (string,error)~ بالبروتوكول والحدود المذكورة. خطأ ترميز يرجع ~nil,ErrNeighborPacket~، وخطأ فك يرجع ~"",ErrNeighborPacket~. ارفض أي حزمة أكبر من 64 بايت أو نسخة أو عدد حقول غلط؛ المدخلات متتغيرش.

نفّذ ~AnswerNeighbor(pc net.PacketConn, station string, deadline time.Time) error~. الاتصال غير ~nil~ ومملوك للمنادي. افحص اسم المحطة قبل I/O، واضبط الموعد مرة، واقرا حزمة واحدة بمخزن 65. الطلب البايظ يرجع ~ErrNeighborPacket~ من غير إرسال. الطلب السليم يترد عليه مرة بعنوان المصدر من ~ReadFrom~ وبنفس الكود. رجّع أخطاء الموعد والقراية والكتابة زي ما هي، والكتابة القصيرة ~io.ErrShortWrite~. متقفلش الاتصال ولا تعمل goroutines أو تقرا وقت حقيقي. ~main~ فاضية؛ مفيش ~Report~.
`,
		Hints: []string{"اعمل فحص صغير للكود والاسم، وافصل codec عن النقل.", "استخدم مخزن أكبر ببايت من الحد، ومتردش على طلب بايظ.", "احتفظ بعنوان ReadFrom للرد؛ المصنع يقبل دالة فتح عشان الاختبار ميحتاجش بورت."},
		En: &LevelText{Title: "Neighbor Discovery", Goal: "Answer one bounded UDP neighbor-discovery packet", Files: []LabFile{{Name: "main.go", Starter: ch313Starter}},
			Lesson: `# A Neighboring Dock's Discovery Packet

## Before you write
- How can you detect an oversized datagram instead of accepting its truncated prefix?
- Should replies target an address inside the packet or its actual source?

## UDP is not a stream
UDP preserves datagram boundaries but may lose, duplicate, or reorder packets. A successful send is not acknowledgement of delivery. Here one operation handles one packet without retries.

~net.ListenPacket~ opens a ~net.PacketConn~ using ~udp4~; ~LocalAddr~ reports the assigned address. Separate opening from answering so tests can inject a transport.
~~~go
pc, err := net.ListenPacket("udp4", "127.0.0.1:0")
if err != nil { return err }
defer pc.Close()
fmt.Println(pc.LocalAddr())
~~~
~OpenNeighbor~ takes a listener function with the same signature; ~nil~ selects ~net.ListenPacket~. Example injection: ~OpenNeighbor(net.ListenPacket)~. No caller-supplied address: this factory binds loopback.

## Receive and answer the source
~ReadFrom~ returns a byte count and the sender's ~net.Addr~; ~WriteTo~ sends one datagram there. ~SetDeadline~ uses the absolute deadline taught in Echo. Allocate 65 bytes for a limit of 64: reading 65 detects overflow even if more bytes were discarded. Never resend a short datagram's remainder; UDP is not TCP.
~~~go
buf := make([]byte, 65)
n, peer, err := pc.ReadFrom(buf)
if err != nil { return err }
if n > 64 { return ErrNeighborPacket }
n, err = pc.WriteTo([]byte("notice"), peer)
if err == nil && n != len("notice") { return io.ErrShortWrite }
~~~
An in-memory ~net.PacketConn~ can inject deadline, read, and write errors. For local operation, the factory's caller closes the socket; the answering function borrows it without closing.

~strings.ContainsRune~ checks membership in an allowed character set; length validation alone is insufficient.
~~~go
allowed := strings.ContainsRune("0123456789abcdef", rune('c'))
fmt.Println(allowed) // true
~~~

## Versioning and reply correlation
Requests are exactly ASCII ~Q1|nonce~ and replies ~A1|nonce|station~. ~1~ is the protocol version. The ~nonce~ is an injected correlation code, not internally generated randomness. It contains 1 through 16 bytes from ~0-9a-f~; station names contain 1 through 16 bytes from ~a-z0-9-~. No spaces or trailing delimiter. ~strings.Split~ exposes surplus fields.
~~~go
parts := strings.Split("A1|ab4|medic", "|")
if len(parts) != 3 { return "", ErrNeighborPacket }
~~~
~EncodeNeighborProbe("ab4")~ creates a request; ~DecodeNeighborReply(packet, "ab4")~ checks version and correlation before returning a station. A different code rejects the reply even if other fields are valid.

## Common mistakes
~~~go
buf := make([]byte, 64)
n, peer, _ := pc.ReadFrom(buf)
// Wrong: a 100-byte datagram may appear to be a valid 64-byte packet
pc.WriteTo(buf[:n], peer)
~~~

## Your mission
Implement ~OpenNeighbor(listen func(string,string)(net.PacketConn,error)) (net.PacketConn,error)~: choose ~net.ListenPacket~ for ~nil~, invoke once with ~udp4~ and ~127.0.0.1:0~, and return its connection and error unchanged. Implement ~EncodeNeighborProbe(nonce string) ([]byte,error)~ and ~DecodeNeighborReply(packet []byte, nonce string) (string,error)~ with the stated wire format and bounds. Encoding failure returns ~nil,ErrNeighborPacket~; decoding failure returns ~"",ErrNeighborPacket~. Reject packets exceeding 64 bytes, wrong versions, or wrong field counts. Do not mutate inputs.

Implement ~AnswerNeighbor(pc net.PacketConn, station string, deadline time.Time) error~. The non-~nil~ connection remains caller-owned. Validate station before I/O, set the deadline once, then read one datagram into 65 bytes. Invalid requests return ~ErrNeighborPacket~ without sending. Answer valid requests exactly once to the ~ReadFrom~ source with the same code. Preserve deadline/read/write errors; short writes return ~io.ErrShortWrite~. Do not close the connection, start goroutines, or read wall time. Empty ~main~; no ~Report~.
`,
			Hints: []string{"Validate codes and names separately; keep the codec independent of transport.", "Read into a max-plus-one buffer and never answer malformed requests.", "Reply to ReadFrom's source; inject the factory listener for socket-free tests."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch313Hidden},
	})
}

const ch313Starter = `package main

import (
	"errors"
	"net"
	"time"
)

var ErrNeighborPacket = errors.New("invalid neighbor packet")

func OpenNeighbor(listen func(string, string) (net.PacketConn, error)) (net.PacketConn, error) {
	return nil, errors.New("TODO")
}
func EncodeNeighborProbe(nonce string) ([]byte, error)                { return nil, ErrNeighborPacket }
func DecodeNeighborReply(packet []byte, nonce string) (string, error) { return "", ErrNeighborPacket }
func AnswerNeighbor(pc net.PacketConn, station string, deadline time.Time) error {
	return ErrNeighborPacket
}
func main() {}
`
const ch313Hidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"orbital/colony/labkit"
	"strings"
	"testing"
	"time"
)

type neighborWire struct {
	net.PacketConn
	packet, reply               []byte
	peer, target                net.Addr
	deadline                    time.Time
	sets, reads, writes, closes int
	setErr, readErr, writeErr   error
	short                       bool
}

func (p *neighborWire) SetDeadline(d time.Time) error { p.sets++; p.deadline = d; return p.setErr }
func (p *neighborWire) ReadFrom(b []byte) (int, net.Addr, error) {
	p.reads++
	return copy(b, p.packet), p.peer, p.readErr
}
func (p *neighborWire) WriteTo(b []byte, a net.Addr) (int, error) {
	p.writes++
	p.target = a
	p.reply = append([]byte(nil), b...)
	if p.short {
		return len(b) - 1, nil
	}
	return len(b), p.writeErr
}
func (p *neighborWire) Close() error { p.closes++; return nil }

type neighborTimeout struct{}

func (neighborTimeout) Error() string   { return "injected packet timeout" }
func (neighborTimeout) Timeout() bool   { return true }
func (neighborTimeout) Temporary() bool { return false }

// Boundaries: min/max/invalid tokens, wrong correlation/version, oversize datagram, timeout, short writes.
func TestHiddenNeighborCodecAndReply(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(313))
	deadline := time.Unix(321, 4)
	for i := 0; i < 60; i++ {
		nonce := fmt.Sprintf("%x", rng.Uint64())
		station := fmt.Sprint("dock-", rng.Intn(1000))
		if i == 0 {
			station = strings.Repeat("a", 16)
		}
		packet, err := EncodeNeighborProbe(nonce)
		if err != nil || string(packet) != "Q1|"+nonce {
			t.Fatal("encode", err)
		}
		before := append([]byte(nil), packet...)
		peer := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1000 + i}
		p := &neighborWire{packet: packet, peer: peer}
		if err = AnswerNeighbor(p, station, deadline); err != nil || p.target != peer || p.writes != 1 || p.sets != 1 || p.deadline != deadline || p.closes != 0 {
			t.Fatal("answer", err)
		}
		if got, err := DecodeNeighborReply(p.reply, nonce); err != nil || got != station {
			t.Fatal("decode", got, err)
		}
		if !bytes.Equal(before, packet) {
			t.Fatal("input mutation")
		}
	}
	for _, code := range []string{"", strings.Repeat("a", 17), "A", "g", "a|b", "a\n", "é"} {
		if b, e := EncodeNeighborProbe(code); b != nil || !errors.Is(e, ErrNeighborPacket) {
			t.Fatal("bad code", code)
		}
	}
	for _, code := range []string{"0", strings.Repeat("a", 16)} {
		if _, e := EncodeNeighborProbe(code); e != nil {
			t.Fatal("exact code")
		}
	}
	for _, packet := range []string{"", "A2|a|dock", "A1|b|dock", "A1|a|", "A1|a|Dock", "A1|a|dock|", "A1|a|" + strings.Repeat("x", 17), strings.Repeat("a", 65)} {
		if got, e := DecodeNeighborReply([]byte(packet), "a"); got != "" || !errors.Is(e, ErrNeighborPacket) {
			t.Fatal("bad reply", packet)
		}
	}
	for _, packet := range []string{"", "Q2|a", "Q1|", "Q1|A", "Q1|a|x", strings.Repeat("x", 65), "Q1|a" + strings.Repeat("x", 100)} {
		p := &neighborWire{packet: []byte(packet)}
		if !errors.Is(AnswerNeighbor(p, "dock", deadline), ErrNeighborPacket) || p.writes != 0 {
			t.Fatal("bad request", packet)
		}
	}
}
func TestHiddenNeighborTransportFailures(t *testing.T) {
	fault := fmt.Errorf("wrapped: %w", neighborTimeout{})
	for _, stage := range []string{"set", "read", "write", "short", "station"} {
		p := &neighborWire{packet: []byte("Q1|ab")}
		want := error(fault)
		station := "dock"
		switch stage {
		case "set":
			p.setErr = fault
		case "read":
			p.readErr = fault
		case "write":
			p.writeErr = fault
		case "short":
			p.short = true
			want = io.ErrShortWrite
		case "station":
			station = "bad name"
			want = ErrNeighborPacket
		}
		if err := AnswerNeighbor(p, station, time.Time{}); !errors.Is(err, want) || p.closes != 0 {
			t.Fatal(stage, err)
		}
		if stage == "station" && (p.sets != 0 || p.reads != 0) {
			t.Fatal("validation I/O")
		}
	}
	p := &neighborWire{}
	calls := 0
	c, e := OpenNeighbor(func(network, addr string) (net.PacketConn, error) {
		calls++
		if network != "udp4" || addr != labkit.FreeAddr() {
			t.Fatal("bind", network, addr)
		}
		return p, fault
	})
	if c != p || !errors.Is(e, fault) || calls != 1 || p.closes != 0 {
		t.Fatal("factory")
	}
}
func TestHiddenNeighborLoopback(t *testing.T) {
	labkit.LeakCheck(t)
	server, err := OpenNeighbor(nil)
	if err != nil {
		var op *net.OpError
		if errors.As(err, &op) && errors.Is(op.Err, errors.ErrUnsupported) {
			t.Skip("UDP sockets unsupported")
		}
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox denies loopback UDP bind")
		}
		t.Fatal(err)
	}
	defer server.Close()
	client, err := net.ListenPacket("udp4", labkit.FreeAddr())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	deadline := time.Now().Add(3 * time.Second)
	if err = client.SetDeadline(deadline); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- AnswerNeighbor(server, "greenhouse", deadline) }()
	probe, _ := EncodeNeighborProbe("f0")
	if _, err = client.WriteTo(probe, server.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 65)
	n, _, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	if station, e := DecodeNeighborReply(buf[:n], "f0"); e != nil || station != "greenhouse" {
		t.Fatal(station, e)
	}
	select {
	case err = <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(4 * time.Second):
		t.Fatal("answer did not stop")
	}
}
`

func init() {
	add(&Level{ID: "31-4-frames", Chapter: 31, Kind: "lab", Title: "إطارات binary", Goal: "فك إطارات تليمتري محفوظة بحدود واضحة واختبار fuzz",
		Requires: []string{"fuzz", "sel:binary.BigEndian"}, Bridge: []int{85}, Race: false, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch314Starter}, {Name: "main_test.go", Starter: ch314TestStarter}},
		Lesson: `# إطارات التليمتري المحفوظة

## قبل ما تكتب
- هتفحص طول الهيدر قبل الوصول لأي index إزاي؟
- هل تعديل نتيجة التحليل ينفع يغير ملف الإدخال في الذاكرة؟

## إطارات ثنائية من غير التقاط شبكة
دي بيانات محقونة من أرشيف المستعمرة، مش Ethernet أو التقاط حزم حقيقي. الهيدر بيقول طول المحتوى، فمش محتاج فاصل سطور. مفيش sockets أو صلاحيات خاصة أو وصول لملفات جوه parser.

~encoding/binary~ بتحوّل أعداد متعددة البايتات بترتيب صريح. ~binary.BigEndian~ معناها البايت الأعلى الأول. ~Uint16~ بتقرا بايتين و ~PutUint16~ بتكتبهم؛ لازم الـ slice طولها 2 على الأقل وإلا panic.
~~~go
buf := make([]byte, 2)
binary.BigEndian.PutUint16(buf, 513) // 2 ثم 1
n := binary.BigEndian.Uint16(buf)
fmt.Println(n)
~~~
الأعداد ~uint16~ من 0 لحد 65535. حوّل الطول لـ ~int~ بعد قرايته وقبل جمعه مع الهيدر، عشان جمع ~uint16~ ممكن يلف.

## عقد الهيدر
~Radiogram~ فيه ~Kind byte~ و ~Sequence uint16~ و ~Payload []byte~. ~ParseRadiogram(wire)~ بيرجع إطار واحد أو قيمة صفرية وخطأ. مثال شكل الإطار: ~Radiogram{Kind: 2, Sequence: 513, Payload: []byte{9}}~.

الهيدر 8 بايت: موضعا 0 و 1 هما ASCII ~O~ و ~C~؛ موضع 2 نسخة قيمتها 1؛ موضع 3 نوع قيمته 1 للتليمتري أو 2 للأوامر؛ موضعا 4 و 5 طول المحتوى big-endian؛ موضعا 6 و 7 رقم التسلسل big-endian. بعدها محتوى من 0 لحد 1024 بايت. الطول الكلي لازم يساوي 8 زائد الطول المعلن، من غير padding أو إطار تاني.
~~~go
if len(wire) < 8 { return Radiogram{}, ErrRadiogram }
length := int(binary.BigEndian.Uint16(wire[4:6]))
if length > 1024 || len(wire) != 8+length { return Radiogram{}, ErrRadiogram }
~~~
افحص الشكل والتوقيع والنسخة والطول الأول. بعد نجاحهم نوع مجهول يرجع ~ErrRadiogramKind~. كل فساد تاني يرجع ~ErrRadiogram~. انسخ المحتوى، متعرضش slice بتشارك الإدخال؛ ~append([]byte{}, payload...)~ نسخة مستقلة.

## fuzz كعقد قابل للتوليد
الفصل 16 شرح ~testing.F~ و ~f.Add~ و ~f.Fuzz~. هنا بنختبر حقول binary وحدودها. استخدم ~binary.BigEndian.PutUint16~ لبناء seed طويلة في ~main_test.go~؛ ده بيحقق شرط المفهوم في ملف التست. مثال مختلف بيختبر طول بادئة:
~~~go
func FuzzTag(f *testing.F) {
 f.Add([]byte{7})
 f.Fuzz(func(t *testing.T, b []byte) {
  if len(b) == 0 { return }
  if len(b[:1]) != 1 { t.Fatal("prefix") }
 })
}
~~~
التست الجيد مش بس «مفيش panic». يتحقق من القبول والرفض والحقول ونسخ المحتوى. احسب oracle الطول بالضرب في 256 والجمع بدل نسخ parser. أضف seeds للصفر و 1 و 256 و 1024 بايت، وهيدر مقطوع ونوع مجهول وطول غلط وبايت زيادة. التوليد يشتغل 5 ثواني عن طريق ~Lab.Fuzz~.

## أخطاء شائعة
~~~go
length := binary.BigEndian.Uint16(wire[4:6]) // غلط لو الهيدر مقطوع
return Radiogram{Payload: wire[8:8+length]}, nil // alias ومن غير فحص الطول
~~~

## المهمة
نفّذ ~ParseRadiogram(wire []byte) (Radiogram,error)~ بالهيدر والأولوية والحدود المذكورة؛ أي خطأ يرجع ~Radiogram{}~ وسببه المحدد. اقبل payload فاضية ورقم تسلسل 0 أو 65535، وانسخ المحتوى من غير تعديل الإدخال. مفيش I/O أو شبكة أو وقت أو ~Report~؛ ~main~ فاضية.

اكتب ~FuzzRadiogram(f *testing.F)~ في ~main_test.go~ مع seeds المذكورة وبناء seed واحدة على الأقل بـ ~binary.BigEndian~. اثبت عدم الـ panic وعدم تعديل المدخل والقبول الصحيح ونسخ المحتوى، وشغّل نفس الاختبارات ضد التلات mutants: endian غلط، نوع مجهول مقبول، و payload بتشارك المدخل. التستات لازم ترفض التلاتة وتعدّي على حلك.
`,
		Hints: []string{"افحص الهيدر قبل slicing؛ افحص الطول قبل النوع عشان أولوية الخطأ.", "حوّل الطول لـ int واطلب مساواة الطول الكلي، وانسخ payload.", "ابني seeds بـ BigEndian، وخلي توقعات fuzz بحساب مستقل، وغيّر payload الناتجة لكشف alias."},
		En: &LevelText{Title: "Binary Frames", Goal: "Parse offline telemetry frames with bounds checks and fuzzing", Files: []LabFile{{Name: "main.go", Starter: ch314Starter}, {Name: "main_test.go", Starter: ch314TestStarter}},
			Lesson: `# Archived Telemetry Frames

## Before you write
- How will you check header length before accessing any index?
- Should mutating a decoded result change the input in memory?

## Binary framing without traffic capture
These injected colony archive bytes are neither Ethernet nor live packet capture. The header declares payload length instead of using line delimiters. The parser needs no sockets, privileges, or file access.

~encoding/binary~ converts multibyte numbers with explicit byte order. ~binary.BigEndian~ puts the high byte first. ~Uint16~ reads two bytes and ~PutUint16~ writes them; slices shorter than two bytes panic.
~~~go
buf := make([]byte, 2)
binary.BigEndian.PutUint16(buf, 513) // 2 then 1
n := binary.BigEndian.Uint16(buf)
fmt.Println(n)
~~~
~uint16~ ranges from 0 through 65535. Convert the decoded length to ~int~ before adding header size; addition in ~uint16~ can wrap.

## Header contract
~Radiogram~ has ~Kind byte~, ~Sequence uint16~, and ~Payload []byte~. ~ParseRadiogram(wire)~ returns one frame or a zero value and error. Example: ~Radiogram{Kind: 2, Sequence: 513, Payload: []byte{9}}~.

The header is 8 bytes: offsets 0 and 1 are ASCII ~O~ and ~C~; offset 2 is version 1; offset 3 is kind 1 for telemetry or 2 for commands; offsets 4 and 5 are big-endian payload length; offsets 6 and 7 are big-endian sequence. The payload has 0 through 1024 bytes. Total size must equal 8 plus declared length, with no padding or second frame.
~~~go
if len(wire) < 8 { return Radiogram{}, ErrRadiogram }
length := int(binary.BigEndian.Uint16(wire[4:6]))
if length > 1024 || len(wire) != 8+length { return Radiogram{}, ErrRadiogram }
~~~
Check shape, signature, version, and length first. Only then an unknown kind returns ~ErrRadiogramKind~. All other corruption returns ~ErrRadiogram~. Copy payload instead of aliasing input; ~append([]byte{}, payload...)~ makes an independent copy.

## Fuzz a generatable contract
Chapter 16 introduced ~testing.F~, ~f.Add~, and ~f.Fuzz~. Here test binary fields and their limits. Build a long seed using ~binary.BigEndian.PutUint16~ in ~main_test.go~, satisfying the concept requirement in the test file. A different example checks prefix length:
~~~go
func FuzzTag(f *testing.F) {
 f.Add([]byte{7})
 f.Fuzz(func(t *testing.T, b []byte) {
  if len(b) == 0 { return }
  if len(b[:1]) != 1 { t.Fatal("prefix") }
 })
}
~~~
A useful test checks more than absence of panic: acceptance, rejection, fields, and independent payload ownership. Compute oracle length by multiplying by 256 and adding instead of copying the parser. Seed payload sizes 0, 1, 256, and 1024, truncated headers, unknown kind, wrong length, and trailing bytes. ~Lab.Fuzz~ generates inputs for 5 seconds.

## Common mistakes
~~~go
length := binary.BigEndian.Uint16(wire[4:6]) // Wrong for truncated headers
return Radiogram{Payload: wire[8:8+length]}, nil // Aliasing without bounds checks
~~~

## Your mission
Implement ~ParseRadiogram(wire []byte) (Radiogram,error)~ with the stated header, validation precedence, and bounds. Every failure returns ~Radiogram{}~ and its specified cause. Accept empty payload and sequence 0 or 65535; copy payload without changing input. No I/O, network, time, or ~Report~. Empty ~main~.

Write ~FuzzRadiogram(f *testing.F)~ in ~main_test.go~ with the specified seeds and at least one seed built using ~binary.BigEndian~. Assert no panic, unchanged input, correct acceptance and fields, and independent payload ownership. Tests must reject three mutants: wrong endianness, accepted unknown kind, and payload aliasing, while passing on your solution.
`,
			Hints: []string{"Check the header before slicing and length before kind to preserve error precedence.", "Convert length to int, require exact total size, and copy payload.", "Build seeds with BigEndian, use an independent fuzz oracle, and mutate returned payload to detect aliasing."}},
		Lab: &Lab{ReqFile: "main_test.go", Bonus: false, Fuzz: "FuzzRadiogram", Hidden: ch314Hidden, Mutants: []Mutant{{Desc: i18n.T("إطار بترتيب بايتات غلط", "Wrong byte order"), Code: ch314EndianMutant}, {Desc: i18n.T("نوع مجهول مقبول", "Unknown kind accepted"), Code: ch314KindMutant}, {Desc: i18n.T("محتوى بيشارك المدخل", "Payload aliases input"), Code: ch314AliasMutant}}},
	})
}

const ch314Starter = `package main

import "errors"

var ErrRadiogram = errors.New("invalid radiogram")
var ErrRadiogramKind = errors.New("unknown radiogram kind")

type Radiogram struct {
	Kind     byte
	Sequence uint16
	Payload  []byte
}

func ParseRadiogram(wire []byte) (Radiogram, error) { return Radiogram{}, ErrRadiogram }
func main()                                         {}
`
const ch314Hidden = `package main

import (
	"bytes"
	"errors"
	"math/rand"
	"testing"
)

// Boundaries: header 0..7, empty/max/oversize payload, all kinds, all sequence bytes, no trailing bytes, no aliasing.
func TestHiddenRadiogramModel(t *testing.T) {
	rng := rand.New(rand.NewSource(314))
	for trial := 0; trial < 200; trial++ {
		length := rng.Intn(1030)
		if trial < 4 {
			length = []int{0, 1, 1024, 1025}[trial]
		}
		sequence := rng.Intn(65536)
		kind := byte(1 + rng.Intn(4))
		wire := make([]byte, 8+length)
		copy(wire, "OC")
		wire[2] = 1
		wire[3] = kind
		wire[4] = byte(length / 256)
		wire[5] = byte(length)
		wire[6] = byte(sequence / 256)
		wire[7] = byte(sequence)
		rng.Read(wire[8:])
		before := append([]byte(nil), wire...)
		got, err := ParseRadiogram(wire)
		var want error
		if length > 1024 {
			want = ErrRadiogram
		} else if kind > 2 {
			want = ErrRadiogramKind
		}
		if !errors.Is(err, want) {
			t.Fatal("error precedence", length, kind, err)
		}
		if want != nil {
			if got.Kind != 0 || got.Sequence != 0 || got.Payload != nil {
				t.Fatal("partial result")
			}
			continue
		}
		if got.Kind != kind || int(got.Sequence) != sequence || !bytes.Equal(got.Payload, wire[8:]) {
			t.Fatal("fields")
		}
		if len(got.Payload) > 0 {
			got.Payload[0] ^= 255
		}
		if !bytes.Equal(before, wire) {
			t.Fatal("alias")
		}
	}
	valid := []byte{'O', 'C', 1, 1, 0, 0, 255, 255}
	if g, e := ParseRadiogram(valid); e != nil || g.Sequence != 65535 {
		t.Fatal("max sequence")
	}
	for n := 0; n < 8; n++ {
		if g, e := ParseRadiogram(valid[:n]); !errors.Is(e, ErrRadiogram) || g.Payload != nil {
			t.Fatal("truncated", n)
		}
	}
	for _, index := range []int{0, 1, 2, 4, 5} {
		bad := append([]byte(nil), valid...)
		bad[index] ^= 255
		if _, e := ParseRadiogram(bad); !errors.Is(e, ErrRadiogram) {
			t.Fatal("header", index, e)
		}
	}
	if _, e := ParseRadiogram(append(valid, 9)); !errors.Is(e, ErrRadiogram) {
		t.Fatal("trailing")
	}
}
`

const ch314TestStarter = `package main

import (
	"encoding/binary"
	"testing"
)

func FuzzRadiogram(f *testing.F) {
	seed := make([]byte, 8)
	copy(seed, "OC")
	seed[2] = 1
	seed[3] = 1
	binary.BigEndian.PutUint16(seed[4:6], 0)
	f.Add(seed)
	f.Fuzz(func(t *testing.T, wire []byte) { /* Add contract assertions and more seeds. */ })
}
`
const ch314EndianMutant = `package main

import (
	"encoding/binary"
	"errors"
)

var ErrRadiogram = errors.New("invalid radiogram")
var ErrRadiogramKind = errors.New("unknown radiogram kind")

type Radiogram struct {
	Kind     byte
	Sequence uint16
	Payload  []byte
}

func ParseRadiogram(wire []byte) (Radiogram, error) {
	if len(wire) < 8 || wire[0] != 'O' || wire[1] != 'C' || wire[2] != 1 {
		return Radiogram{}, ErrRadiogram
	}
	length := int(binary.LittleEndian.Uint16(wire[4:6]))
	if length > 1024 || len(wire) != 8+length {
		return Radiogram{}, ErrRadiogram
	}
	if wire[3] != 1 && wire[3] != 2 {
		return Radiogram{}, ErrRadiogramKind
	}
	return Radiogram{Kind: wire[3], Sequence: binary.LittleEndian.Uint16(wire[6:8]), Payload: append([]byte{}, wire[8:]...)}, nil
}
func main() {}
`
const ch314KindMutant = `package main

import (
	"encoding/binary"
	"errors"
)

var ErrRadiogram = errors.New("invalid radiogram")
var ErrRadiogramKind = errors.New("unknown radiogram kind")

type Radiogram struct {
	Kind     byte
	Sequence uint16
	Payload  []byte
}

func ParseRadiogram(wire []byte) (Radiogram, error) {
	if len(wire) < 8 || wire[0] != 'O' || wire[1] != 'C' || wire[2] != 1 {
		return Radiogram{}, ErrRadiogram
	}
	length := int(binary.BigEndian.Uint16(wire[4:6]))
	if length > 1024 || len(wire) != 8+length {
		return Radiogram{}, ErrRadiogram
	}
	if wire[3] == 0 {
		return Radiogram{}, ErrRadiogramKind
	}
	return Radiogram{Kind: wire[3], Sequence: binary.BigEndian.Uint16(wire[6:8]), Payload: append([]byte{}, wire[8:]...)}, nil
}
func main() {}
`
const ch314AliasMutant = `package main

import (
	"encoding/binary"
	"errors"
)

var ErrRadiogram = errors.New("invalid radiogram")
var ErrRadiogramKind = errors.New("unknown radiogram kind")

type Radiogram struct {
	Kind     byte
	Sequence uint16
	Payload  []byte
}

func ParseRadiogram(wire []byte) (Radiogram, error) {
	if len(wire) < 8 || wire[0] != 'O' || wire[1] != 'C' || wire[2] != 1 {
		return Radiogram{}, ErrRadiogram
	}
	length := int(binary.BigEndian.Uint16(wire[4:6]))
	if length > 1024 || len(wire) != 8+length {
		return Radiogram{}, ErrRadiogram
	}
	if wire[3] != 1 && wire[3] != 2 {
		return Radiogram{}, ErrRadiogramKind
	}
	return Radiogram{Kind: wire[3], Sequence: binary.BigEndian.Uint16(wire[6:8]), Payload: wire[8:]}, nil
}
func main() {}
`

func init() {
	add(&Level{ID: "31-5-dns", Chapter: 31, Kind: "lab", Title: "دليل العناوين", Goal: "اجمع عناوين أرصفة المستعمرة من resolver محقون بترتيب ثابت",
		Requires: []string{"iface"}, Bridge: []int{80}, Race: false, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch315Starter}},
		Lesson: `# دليل عناوين الأرصفة

## قبل ما تكتب
- عدم وجود اسم مختلف عن إجابة ناجحة من غير عناوين إزاي؟
- هتختبر DNS من غير سؤال resolver الجهاز إزاي؟

## أنواع سجلات DNS
سجل ~A~ بيربط الاسم بعنوان IPv4، و ~AAAA~ بعنوان IPv6. ~CNAME~ بيربط الاسم باسم canonical؛ مش عنوان IP. ~NXDOMAIN~ معناها الاسم مش موجود، مش مجرد إن نوع سجل معين فاضي. الاختبار هنا عمره ما بينادي DNS حقيقي.

~DockResolver~ واجهة صغيرة يعرّفها اللاعب: ~Addresses(context.Context,string) ([]net.IP,error)~ و ~Alias(context.Context,string) (string,error)~. الأسماء بتتسأل تحت ~.orbital.invalid~؛ نطاق مخصص للشرح والاختبار، مش عنوان جهاز حقيقي. الفاصل الوحيد للشبكة هو الواجهة المحقونة.
~~~go
type DockResolver interface {
 Addresses(context.Context, string) ([]net.IP, error)
 Alias(context.Context, string) (string, error)
}
ips, err := resolver.Addresses(ctx, "medic.orbital.invalid")
canonical, err := resolver.Alias(ctx, "medic.orbital.invalid")
~~~
المنطق يمرّر نفس ~ctx~ ونفس الاسم للعمليتين، بالترتيب: عناوين وبعدها alias. عند خطأ العناوين متسألش alias. افحص إلغاء السياق قبل كل نداء وبعده. تنفيذ الواجهة نفسه مسؤول عن احترام السياق؛ مش هنفتح goroutine تلف حوالين دالة متجاهلاه.

## تمثيل IP وتصنيفه
~net.IP~ بايتات عنوان. ~net.ParseIP~ يفك نص، و ~To4~ يرجع IPv4 أو ~nil~؛ عنوان IPv4-mapped IPv6 بيتحسب IPv4. ~To16~ بيرجع تمثيل IPv6 أو ~nil~ لو العنوان بايظ. ~String~ بيطلع تمثيل canonical.
~~~go
ip := net.ParseIP("2001:db8::9")
if ip.To4() != nil { fmt.Println("A", ip.String()) }
if ip.To16() == nil { return AddressCard{}, ErrDockData }
fmt.Println(ip.String())
~~~
~AddressCard~ نتيجة فيها ~Station~ و ~Canonical~ نصوص، و ~A~ و ~AAAA~ قوائم نصوص. مثال: ~AddressCard{Station: "medic", A: []string{"127.0.0.2"}}~. رتّب كل عائلة ترتيب نصّي تصاعدي، مش ترتيب رقمي. متعدّلش slice أو بايتات resolver. إجابة عناوين فاضية نجاح، حتى لو alias فاضية.

## اسم canonical وخطأ الغياب
ارفض أي بايت غير ASCII في الاسم canonical قبل تحويل الحروف؛ بعدها سوّيه بحروف صغيرة وشيل نقطة نهائية واحدة بس. الاسم الفاضي مسموح كغياب metadata. غيره لازم اسم ASCII صالح: طول كلي لحد 253، labels من 1 لحد 63 من حروف ~a-z~ أو أرقام أو ~-~، ومفيش label تبدأ أو تنتهي بشرطة. مفيش labels فاضية.
~~~go
canonical := strings.TrimSuffix(strings.ToLower("Medic.Orbital.Invalid."), ".")
labels := strings.Split(canonical, ".")
fmt.Println(labels)
~~~
~net.DNSError~ نوع خطأ فيه ~IsNotFound~ للغياب. ~errors.As~ يطلّعه حتى مع wrapping. الغياب يرجع خطأ يطابق ~ErrDockAbsent~ وسبب resolver مع بعض؛ المهلة والأخطاء التانية تتحفظ زي ما هي. الإلغاء بعد I/O له الأولوية على النتيجة.
~~~go
var dns *net.DNSError
if errors.As(err, &dns) && dns.IsNotFound {
 return AddressCard{}, fmt.Errorf("%w: %w", ErrDockAbsent, err)
}
~~~

## أخطاء شائعة
~~~go
ips, _ := net.LookupIP(host) // غلط: اعتماد على جهاز المستخدم وخطأ ضايع
return AddressCard{A: []string{ips[0].String()}}, nil // ترتيب وفحص ناقصين
~~~

## المهمة
نفّذ ~ResolveDock(ctx context.Context, resolver DockResolver, station string) (AddressCard,error)~. السياق والواجهة غير ~nil~. المحطة من 1 لحد 16 بايت من ~a-z0-9-~ ومش بتبدأ أو تنتهي بـ ~-~؛ أي اسم بايظ يرجع ~ErrDockData~ قبل أي نداء. اسأل بالظبط عن ~station + ".orbital.invalid"~. اتبع ترتيب النداءات والإلغاء والتصنيف المذكور. أي خطأ يرجع ~AddressCard{}~ من غير نتيجة جزئية. أي عنوان بايظ أو canonical بايظ يرجع ~ErrDockData~.

في النجاح خزّن اسم المحطة الأصلي والاسم canonical المُسوّى، وصنّف كل IP للعائلة الصحيحة ورتّب كل عائلة نصياً. ممنوع استخدام resolver افتراضي أو نداء شبكة أو وقت حقيقي؛ ~main~ فاضية ومفيش ~Report~.

النجمة التالتة: شيل تكرار العناوين بعد التحويل canonical، بما فيها IPv4-mapped وكتابات IPv6 المختلفة لنفس العنوان، من غير تعديل إجابة resolver.
`,
		Hints: []string{"ابدأ بفحص اسم المحطة والسياق؛ مرّر نفس ctx والاسم للواجهتين.", "استخدم To4 قبل To16 وصنّف الأخطاء بـ errors.As، متخلطش الفراغ بالغياب.", "سوّي canonical ورتّب كل عائلة؛ خريطة للنص canonical تمنع التكرار للنجمة الإضافية."},
		En: &LevelText{Title: "Address Directory", Goal: "Build deterministic dock address cards through an injected resolver", Files: []LabFile{{Name: "main.go", Starter: ch315Starter}},
			Lesson: `# Dock Address Directory

## Before you write
- How does a nonexistent name differ from a successful empty address answer?
- How can you test DNS without consulting the machine's resolver?

## DNS record types
An ~A~ record maps a name to IPv4; ~AAAA~ maps to IPv6. ~CNAME~ supplies a canonical name, not an IP address. ~NXDOMAIN~ means the name does not exist, rather than one record type being empty. Tests never call real DNS.

The player defines ~DockResolver~ with ~Addresses(context.Context,string) ([]net.IP,error)~ and ~Alias(context.Context,string) (string,error)~. Queries use ~.orbital.invalid~, a documentation/test domain rather than a real machine. The injected interface is the only network boundary.
~~~go
type DockResolver interface {
 Addresses(context.Context, string) ([]net.IP, error)
 Alias(context.Context, string) (string, error)
}
ips, err := resolver.Addresses(ctx, "medic.orbital.invalid")
canonical, err := resolver.Alias(ctx, "medic.orbital.invalid")
~~~
Pass the same ~ctx~ and name to both operations, in order: addresses then alias. Do not ask for alias after an address error. Check cancellation before and after each call. The resolver implementation must honor context; do not wrap an uncooperative method in a detached goroutine.

## Represent and classify IPs
~net.IP~ stores address bytes. ~net.ParseIP~ parses text; ~To4~ returns IPv4 or ~nil~. IPv4-mapped IPv6 counts as IPv4. ~To16~ returns IPv6 representation or ~nil~ for invalid bytes. ~String~ produces canonical text.
~~~go
ip := net.ParseIP("2001:db8::9")
if ip.To4() != nil { fmt.Println("A", ip.String()) }
if ip.To16() == nil { return AddressCard{}, ErrDockData }
fmt.Println(ip.String())
~~~
~AddressCard~ has string fields ~Station~ and ~Canonical~, and string slices ~A~ and ~AAAA~. Example: ~AddressCard{Station: "medic", A: []string{"127.0.0.2"}}~. Sort each family lexically ascending, not numerically. Do not mutate resolver slices or IP bytes. Empty addresses succeed even with an empty alias.

## Canonical name and absence
Reject non-ASCII canonical bytes before case conversion, then lowercase and remove at most one trailing dot. Empty metadata is allowed. Otherwise require an ASCII name of at most 253 bytes, with labels of 1 through 63 bytes drawn from ~a-z~, digits, and ~-~. Labels cannot start or end with hyphens; empty labels are invalid.
~~~go
canonical := strings.TrimSuffix(strings.ToLower("Medic.Orbital.Invalid."), ".")
labels := strings.Split(canonical, ".")
fmt.Println(labels)
~~~
~net.DNSError~ has ~IsNotFound~ to identify absence. ~errors.As~ finds it through wrapping. Absence returns an error matching both ~ErrDockAbsent~ and the resolver cause. Preserve timeout and other failures unchanged. Cancellation observed after I/O takes precedence over returned data.
~~~go
var dns *net.DNSError
if errors.As(err, &dns) && dns.IsNotFound {
 return AddressCard{}, fmt.Errorf("%w: %w", ErrDockAbsent, err)
}
~~~

## Common mistakes
~~~go
ips, _ := net.LookupIP(host) // Wrong: machine-dependent and error discarded
return AddressCard{A: []string{ips[0].String()}}, nil // Missing validation and ordering
~~~

## Your mission
Implement ~ResolveDock(ctx context.Context, resolver DockResolver, station string) (AddressCard,error)~. Context and resolver are non-~nil~. Stations have 1 through 16 bytes from ~a-z0-9-~ and cannot start or end with ~-~. Invalid stations return ~ErrDockData~ before any resolver call. Query exactly ~station + ".orbital.invalid"~. Follow the described call order, cancellation, and error classification. All errors return ~AddressCard{}~ without partial data. Invalid IPs or canonical names return ~ErrDockData~.

Success stores the original station, normalized canonical name, and each canonical IP in its correct family, sorted lexically. Never use the default resolver, real networking, or wall time. Empty ~main~; no ~Report~.

Third star: deduplicate after canonicalization, including IPv4-mapped forms and alternative IPv6 spellings, without mutating resolver answers.
`,
			Hints: []string{"Validate station and context first; pass the same context and name to both methods.", "Use To4 before To16 and errors.As for DNS failures; empty answers are not absence.", "Normalize canonical names and sort each family; deduplicate canonical strings for the bonus."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch315Hidden},
	})
}

const ch315Starter = `package main

import (
	"context"
	"errors"
	"net"
)

var ErrDockData = errors.New("invalid dock directory data")
var ErrDockAbsent = errors.New("dock name absent")

type DockResolver interface {
	Addresses(context.Context, string) ([]net.IP, error)
	Alias(context.Context, string) (string, error)
}
type AddressCard struct {
	Station, Canonical string
	A, AAAA            []string
}

func ResolveDock(ctx context.Context, resolver DockResolver, station string) (AddressCard, error) {
	return AddressCard{}, errors.New("TODO")
}
func main() {}
`
const ch315Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"orbital/colony/labkit"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type dockResolverFake struct {
	t               *testing.T
	ctx             context.Context
	host            string
	ips             []net.IP
	alias           string
	ipErr, aliasErr error
	calls           []string
	cancel          context.CancelFunc
	cancelAt        string
}

func (r *dockResolverFake) check(ctx context.Context, host, stage string) {
	r.t.Helper()
	if ctx != r.ctx || host != r.host {
		r.t.Error("resolver arguments", host)
	}
	r.calls = append(r.calls, stage)
	if r.cancelAt == stage {
		r.cancel()
	}
}
func (r *dockResolverFake) Addresses(ctx context.Context, host string) ([]net.IP, error) {
	r.check(ctx, host, "ips")
	return r.ips, r.ipErr
}
func (r *dockResolverFake) Alias(ctx context.Context, host string) (string, error) {
	r.check(ctx, host, "alias")
	return r.alias, r.aliasErr
}
func newDockFake(t *testing.T, station string) *dockResolverFake {
	return &dockResolverFake{t: t, ctx: context.Background(), host: station + ".orbital.invalid"}
}
func emptyDockCard(g AddressCard) bool {
	return g.Station == "" && g.Canonical == "" && g.A == nil && g.AAAA == nil
}

// Boundaries: empty/single/many families, lexical ordering, malformed names/IPs, absence vs empty, cancellation, canonical duplicates.
func TestHiddenDockFamiliesAndOrdering(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(315))
	for trial := 0; trial < 40; trial++ {
		station := fmt.Sprint("dock-", trial)
		r := newDockFake(t, station)
		r.alias = "Relay.Orbital.Invalid."
		wantA, wantAAAA := []string{}, []string{}
		for j := 0; j < trial; j++ {
			a := fmt.Sprintf("127.0.%d.%d", trial, j)
			b := fmt.Sprintf("2001:db8:%x::%x", trial, j+1)
			wantA = append(wantA, a)
			wantAAAA = append(wantAAAA, b)
			r.ips = append(r.ips, net.ParseIP(a), net.ParseIP(b))
		}
		rng.Shuffle(len(r.ips), func(i, j int) { r.ips[i], r.ips[j] = r.ips[j], r.ips[i] })
		before := make([]net.IP, len(r.ips))
		for i, ip := range r.ips {
			before[i] = append(net.IP(nil), ip...)
		}
		slices.Sort(wantA)
		slices.Sort(wantAAAA)
		g, e := ResolveDock(r.ctx, r, station)
		if e != nil || g.Station != station || g.Canonical != "relay.orbital.invalid" || !slices.Equal(g.A, wantA) || !slices.Equal(g.AAAA, wantAAAA) || !slices.Equal(r.calls, []string{"ips", "alias"}) {
			t.Fatal("card", g, e)
		}
		if len(r.ips) > 0 && !reflect.DeepEqual(before, r.ips) {
			t.Fatal("input mutation")
		}
	}
	for _, alias := range []string{strings.Repeat("a", 63) + ".invalid", strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 61)} {
		station := strings.Repeat("s", 16)
		r := newDockFake(t, station)
		r.alias = alias + "."
		g, e := ResolveDock(r.ctx, r, station)
		if e != nil || g.Canonical != alias || g.Station != station {
			t.Fatal("exact DNS name bounds", g, e)
		}
	}
	r := newDockFake(t, "empty")
	g, e := ResolveDock(r.ctx, r, "empty")
	if e != nil || g.Station != "empty" || g.Canonical != "" || len(g.A)+len(g.AAAA) != 0 {
		t.Fatal("empty is success")
	}
}
func TestHiddenDockFailuresAndValidation(t *testing.T) {
	for _, name := range []string{"", "-a", "a-", "A", "a.b", "é", strings.Repeat("a", 17), "host:53"} {
		r := newDockFake(t, name)
		g, e := ResolveDock(r.ctx, r, name)
		if !errors.Is(e, ErrDockData) || !emptyDockCard(g) || len(r.calls) != 0 {
			t.Fatal("invalid station", name, e)
		}
	}
	for _, alias := range []string{".", "a..b", "a..", "-a.invalid", "a-.invalid", "a/b", "x:53", "é.invalid", "K.invalid", strings.Repeat("a", 64) + ".invalid", strings.Repeat("a.", 127) + "a"} {
		r := newDockFake(t, "dock")
		r.alias = alias
		g, e := ResolveDock(r.ctx, r, "dock")
		if !errors.Is(e, ErrDockData) || !emptyDockCard(g) {
			t.Fatal("invalid alias", alias, e)
		}
	}
	r := newDockFake(t, "dock")
	r.ips = []net.IP{net.ParseIP("127.0.0.1"), {1, 2, 3}}
	if g, e := ResolveDock(r.ctx, r, "dock"); !errors.Is(e, ErrDockData) || !emptyDockCard(g) {
		t.Fatal("invalid IP")
	}
	for _, stage := range []string{"ips", "alias"} {
		for _, cause := range []error{&net.DNSError{Err: "absent", IsNotFound: true}, &net.DNSError{Err: "slow", IsTimeout: true}, context.DeadlineExceeded, errors.New("resolver broken")} {
			r := newDockFake(t, "dock")
			wrapped := fmt.Errorf("cause: %w", cause)
			if stage == "ips" {
				r.ipErr = wrapped
			} else {
				r.aliasErr = wrapped
			}
			g, e := ResolveDock(r.ctx, r, "dock")
			if !errors.Is(e, cause) || !emptyDockCard(g) {
				t.Fatal("cause", e)
			}
			var dns *net.DNSError
			if errors.As(cause, &dns) && dns.IsNotFound && !errors.Is(e, ErrDockAbsent) {
				t.Fatal("NXDOMAIN")
			}
			if stage == "ips" && len(r.calls) != 1 {
				t.Fatal("called alias after error")
			}
		}
	}
	for _, stage := range []string{"pre", "ips", "alias"} {
		r := newDockFake(t, "dock")
		ctx, cancel := context.WithCancel(context.Background())
		r.ctx = ctx
		r.cancel = cancel
		r.cancelAt = stage
		if stage == "pre" {
			cancel()
		}
		g, e := ResolveDock(ctx, r, "dock")
		cancel()
		if !errors.Is(e, context.Canceled) || !emptyDockCard(g) {
			t.Fatal("cancel", stage, e)
		}
		if stage == "pre" && len(r.calls) != 0 || stage == "ips" && len(r.calls) != 1 {
			t.Fatal("work after cancel")
		}
	}
}
func TestHiddenBonusDockCanonicalDeduplication(t *testing.T) {
	r := newDockFake(t, "dock")
	for _, s := range []string{"127.0.0.2", "::ffff:127.0.0.2", "2001:db8:0:0::1", "2001:db8::1", "127.0.0.2"} {
		r.ips = append(r.ips, net.ParseIP(s))
	}
	g, e := ResolveDock(r.ctx, r, "dock")
	if e != nil || !slices.Equal(g.A, []string{"127.0.0.2"}) || !slices.Equal(g.AAAA, []string{"2001:db8::1"}) {
		t.Fatal("dedup", g, e)
	}
}
`

func init() {
	add(&Level{ID: "31-6-port-check", Chapter: 31, Kind: "lab", Title: "فحص المنافذ", Goal: "افحص قائمة منافذ أرصفة محلية بحد توازي وdial محقون",
		Requires: []string{"go", "sel:net.DialTimeout"}, Bridge: []int{82}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch316Starter}},
		Lesson: `# فحص أرصفة على جهازك فقط

## قبل ما تكتب
- هتمنع أي اتصال خارج جهازك قبل تشغيل أول worker إزاي؟
- هتثبت إن حد التوازي فعلي من غير قياس مدة التنفيذ إزاي؟

## حدود الاستخدام
فحص أجهزة مش ملكك مش مسموح. المهمة دي لجهازك المحلي فقط، ومش بتسمح بأي hostname حتى ~localhost~ عشان متعتمدش على DNS أو إعدادات الجهاز. اقبل بس عنوان IP رقمي بيعدّي ~IsLoopback~، زي ~127.0.0.1~ أو ~::1~. عنوان زي ~127.0.0.2~ محلي كمان؛ عنوان عام أو wildcard مرفوض قبل الـ dial.
~~~go
ip := net.ParseIP(host)
if ip == nil || !ip.IsLoopback() { return nil, ErrPortPlan }
~~~
~net.ParseIP~ و ~net.IP~ اتشرحوا في دليل العناوين. ~IsLoopback~ بتفحص نطاق العنوان نفسه، مش تشابه النص. متستخدمش prefix نصّي زي ~127~ أو تحل hostname ثم تفحصه.

## فتح اتصال بمهلة
~net.DialTimeout~ بتفتح اتصال بالـ network والعنوان والمدة، وترجع ~net.Conn~ أو خطأ. المهلة لكل محاولة، مش لكل القائمة. هنا ~Dial~ callback بنفس التوقيع؛ لو ~nil~ اختار الدالة القياسية. الاختبارات بتحقن نتيجة فورية للفتح والقفل والمهلة بدل انتظار وقت حقيقي.
~~~go
dial := net.DialTimeout
conn, err := dial("tcp", "127.0.0.1:8500", 200*time.Millisecond)
if conn != nil { defer conn.Close() }
if err != nil { return err }
~~~
~net.JoinHostPort~ بتضيف أقواس IPv6 صح؛ ~strconv.Itoa~ بتحوّل البورت لنص. استخدم تمثيل ~ip.String()~ المحقون للعنوان، ومتحاولش تركّب IPv6 بالجمع بس.
~~~go
address := net.JoinHostPort(ip.String(), strconv.Itoa(8500))
fmt.Println(address)
~~~
~PortOptions~ فيها ~Workers int~ و ~Timeout time.Duration~ و ~Dial func(string,string,time.Duration)(net.Conn,error)~. مثال حقن الإنتاج: ~PortOptions{Workers: 3, Timeout: time.Second, Dial: net.DialTimeout}~. المنطق مبيقراش ساعة الجهاز.

## توازي محدود بجد
بدل goroutine لكل بورت وبعدين الانتظار، شغّل على الأكثر ~Workers~ عامل بياخدوا indices من قناة شغل. كل عامل يكتب خانته في slice نتائج بنفس ترتيب قائمة المنافذ، وبعدها ~WaitGroup.Wait~ قبل الرجوع. خلي مالك واحد يقفل قناة الشغل بعد إرسال كل indices. ممكن تملأ قناة بسعة طول القائمة قبل تشغيل العمال لأن عدد المنافذ محدود.
~~~go
jobs := make(chan int, len(ports))
for i := range ports { jobs <- i }
close(jobs)
wg.Add(1)
go func() { defer wg.Done(); for i := range jobs { replies[i] = check(ports[i]) } }()
wg.Wait()
~~~
المثال عامل واحد؛ المهمة محتاجة عدد عوامل يساوي الأصغر بين ~Workers~ وطول القائمة. التست يحبس الـ dial المحقونة على قناة، ويستنى وصول العدد المتوقع قبل فتح الحاجز. كده نثبت الحد والقدرة على التوازي بأحداث، مش ~time.Sleep~.

## حالة منفذ وسبب الفشل
~PortStatus~ فيها ~Port int~ و ~State string~ و ~Err error~. اتصال غير ~nil~ من غير خطأ يعني ~open~. خطأ يحقق ~net.Error.Timeout()~ أو يطابق ~context.DeadlineExceeded~ يعني ~timeout~. باقي الأخطاء ~closed~، ودي حالة فشل فتح في العقد المبسّط، مش إثبات تقني لسبب الرفض. حافظ على سبب الخطأ. ~nil,nil~ من dial غير صالح يتصنّف ~closed~ مع ~ErrPortDial~.
~~~go
var timeout net.Error
if errors.As(err, &timeout) && timeout.Timeout() { state = "timeout" }
~~~
كل اتصال راجع من الـ dial يتقفل مرة قبل رجوع العامل، حتى لو معاه خطأ. خطأ القفل مش بيغير نتيجة المحاولة. مفيش قراية banner أو إرسال بيانات أو إعادة محاولة.

## أخطاء شائعة
~~~go
if strings.HasPrefix(host, "127") {
 go dial("tcp", host+":"+port, timeout) // غلط: prefix مش تحقق عنوان، ومفيش حد توازي
}
~~~

## المهمة
نفّذ ~CheckDockPorts(host string, ports []int, options PortOptions) ([]PortStatus,error)~. ارفض أي host غير IP loopback رقمي، أو أكتر من 64 منفذ، أو بورت خارج 1 لحد 65535، أو بورت مكرر، أو ~Workers~ خارج 1 لحد 16، أو مدة غير موجبة بـ ~nil,ErrPortPlan~ قبل أي dial أو goroutine. متعدّلش قائمة المدخل. القائمة الفاضية نجاح من غير dial أو goroutine بعد فحص باقي الإعدادات.

شغّل عدد عوامل يساوي الأصغر بين عدد المنافذ و ~Workers~، وكل منفذ يتحاول مرة واحدة بـ ~tcp~ والعنوان والمدة المحددين. استخدم ~go~ و ~net.DialTimeout~ كافتراضي قابل للحقن. النتائج بنفس ترتيب ~ports~، وكلها تظهر بعد اكتمال العمال وقفل كل اتصال. فشل منفذ يرجع في خانة ~Err~ وميرجعش خطأ عام؛ الإعداد البايظ وحده خطأ عام. حافظ على causes بـ ~errors.Is~، وصنّف زي الشرح. ممنوع DNS أو فحص أي جهاز خارجي أو ~Report~؛ ~main~ فاضية.
`,
		Hints: []string{"افحص الخطة كلها قبل العمال، واستخدم ParseIP وIsLoopback بدل فحص النص.", "وزّع indices على عدد عوامل محدود، واكتب نتيجة كل index في مكانها.", "اختار DialTimeout لو callback غايبة، واقفل أي conn راجعة، وصنّف المهلة بـ errors.As."},
		En: &LevelText{Title: "Port Check", Goal: "Check registered local dock ports with bounded concurrent injected dialing", Files: []LabFile{{Name: "main.go", Starter: ch316Starter}},
			Lesson: `# Check Docks on Your Own Machine

## Before you write
- How will you reject external targets before starting any worker?
- How can you prove a real concurrency bound without timing execution?

## Usage boundary
Scanning machines you do not own is not allowed. This mission permits only your local machine and rejects every hostname, including ~localhost~, to avoid DNS and machine configuration. Accept only numeric IPs satisfying ~IsLoopback~, such as ~127.0.0.1~ or ~::1~. ~127.0.0.2~ is also loopback. Public and wildcard addresses fail before dialing.
~~~go
ip := net.ParseIP(host)
if ip == nil || !ip.IsLoopback() { return nil, ErrPortPlan }
~~~
The directory mission introduced ~net.ParseIP~ and ~net.IP~. ~IsLoopback~ tests the actual address range, not a similar-looking string. Do not use a textual ~127~ prefix or resolve hostnames before checking them.

## Dial with a per-attempt timeout
~net.DialTimeout~ takes network, address, and duration, returning a ~net.Conn~ or error. The timeout covers one attempt, not the whole list. ~Dial~ is an injected callback with the same signature; ~nil~ selects the standard function. Tests inject immediate open, closed, and timeout outcomes rather than wait for wall time.
~~~go
dial := net.DialTimeout
conn, err := dial("tcp", "127.0.0.1:8500", 200*time.Millisecond)
if conn != nil { defer conn.Close() }
if err != nil { return err }
~~~
~net.JoinHostPort~ adds IPv6 brackets correctly; ~strconv.Itoa~ formats a port. Use ~ip.String()~ for the numeric address instead of concatenating IPv6 manually.
~~~go
address := net.JoinHostPort(ip.String(), strconv.Itoa(8500))
fmt.Println(address)
~~~
~PortOptions~ has ~Workers int~, ~Timeout time.Duration~, and ~Dial func(string,string,time.Duration)(net.Conn,error)~. Production injection example: ~PortOptions{Workers: 3, Timeout: time.Second, Dial: net.DialTimeout}~. The logic never reads wall time.

## Actually bounded concurrency
Start at most ~Workers~ goroutines consuming indices from a jobs channel. Each writes a distinct result slot in input order; ~WaitGroup.Wait~ joins them before returning. One owner closes jobs after sending all indices. Since the list is bounded, it can be buffered and filled before workers start.
~~~go
jobs := make(chan int, len(ports))
for i := range ports { jobs <- i }
close(jobs)
wg.Add(1)
go func() { defer wg.Done(); for i := range jobs { replies[i] = check(ports[i]) } }()
wg.Wait()
~~~
This example has one worker; the mission uses the minimum of ~Workers~ and list length. Tests block injected dialing on a channel until the expected number enters, then release a barrier. This proves both the bound and actual parallelism through events rather than ~time.Sleep~.

## Port outcome and failure cause
~PortStatus~ has ~Port int~, ~State string~, and ~Err error~. A non-~nil~ connection without error is ~open~. An error satisfying ~net.Error.Timeout()~ or matching ~context.DeadlineExceeded~ is ~timeout~. Other failures are ~closed~, meaning dial failure in this simplified contract rather than proof of a specific refusal cause. Preserve the cause. Invalid ~nil,nil~ dialing becomes ~closed~ with ~ErrPortDial~.
~~~go
var timeout net.Error
if errors.As(err, &timeout) && timeout.Timeout() { state = "timeout" }
~~~
Close every returned connection exactly once before its worker completes, even when returned alongside an error. A close error does not reclassify the attempt. Do not read banners, send data, or retry.

## Common mistakes
~~~go
if strings.HasPrefix(host, "127") {
 go dial("tcp", host+":"+port, timeout) // Wrong: no address validation or concurrency bound
}
~~~

## Your mission
Implement ~CheckDockPorts(host string, ports []int, options PortOptions) ([]PortStatus,error)~. Reject nonnumeric or non-loopback hosts, more than 64 ports, ports outside 1 through 65535, duplicate ports, ~Workers~ outside 1 through 16, or nonpositive timeout with ~nil,ErrPortPlan~ before any dial or goroutine. Never mutate input. Empty lists succeed without dialing or goroutines after validating the other settings.

Start the minimum of port count and ~Workers~ workers. Try each port exactly once using ~tcp~, the specified address, and timeout. Use ~go~ and injectable ~net.DialTimeout~ as the default. Results preserve input order and return after joining workers and closing every connection. Port failures belong in their ~Err~ slots, with no overall error; only invalid configuration causes an overall error. Preserve causes for ~errors.Is~ and classify as explained. No DNS, external scanning, or ~Report~; empty ~main~.
`,
			Hints: []string{"Validate the whole plan before workers; use ParseIP and IsLoopback instead of text matching.", "Distribute indices to bounded workers and write each result into its own slot.", "Default to DialTimeout, close every returned connection, and classify timeouts with errors.As."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch316Hidden},
	})
}

const ch316Starter = `package main

import (
	"errors"
	"net"
	"time"
)

var ErrPortPlan = errors.New("invalid local port plan")
var ErrPortDial = errors.New("dial returned no connection")

type PortOptions struct {
	Workers int
	Timeout time.Duration
	Dial    func(string, string, time.Duration) (net.Conn, error)
}
type PortStatus struct {
	Port  int
	State string
	Err   error
}

func CheckDockPorts(host string, ports []int, options PortOptions) ([]PortStatus, error) {
	return nil, errors.New("TODO")
}
func main() {}
`
const ch316Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"orbital/colony/labkit"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type dockPortConn struct {
	net.Conn
	closes atomic.Int32
}

func (c *dockPortConn) Close() error { c.closes.Add(1); return errors.New("close observation only") }

type dockPortTimeout struct{}

func (dockPortTimeout) Error() string   { return "injected dial deadline" }
func (dockPortTimeout) Timeout() bool   { return true }
func (dockPortTimeout) Temporary() bool { return false }

// Boundaries: loopback v4/v6 only, min/max/duplicates, zero jobs, exact worker peak, timeout vs ordinary errors, close ownership.
func TestHiddenPortsValidationBeforeDial(t *testing.T) {
	labkit.LeakCheck(t)
	var calls atomic.Int32
	dial := func(string, string, time.Duration) (net.Conn, error) {
		calls.Add(1)
		return nil, errors.New("must not dial")
	}
	base := PortOptions{Workers: 2, Timeout: time.Second, Dial: dial}
	for _, host := range []string{"", "localhost", "dock.invalid", "127.evil.invalid", "127.0.0.1:80", "127.0.0.999", "0.0.0.0", "192.0.2.1", "::", "::2", "[::1]", "::1%lo0"} {
		if g, e := CheckDockPorts(host, []int{80}, base); g != nil || !errors.Is(e, ErrPortPlan) {
			t.Fatal("host accepted", host, e)
		}
	}
	for _, ports := range [][]int{{0}, {-1}, {65536}, {1, 1}, make([]int, 65)} {
		if g, e := CheckDockPorts("127.0.0.1", ports, base); g != nil || !errors.Is(e, ErrPortPlan) {
			t.Fatal("ports", ports, e)
		}
	}
	for _, workers := range []int{-1, 0, 17} {
		o := base
		o.Workers = workers
		if _, e := CheckDockPorts("127.0.0.1", nil, o); !errors.Is(e, ErrPortPlan) {
			t.Fatal("workers")
		}
	}
	for _, timeout := range []time.Duration{-1, 0} {
		o := base
		o.Timeout = timeout
		if _, e := CheckDockPorts("127.0.0.1", nil, o); !errors.Is(e, ErrPortPlan) {
			t.Fatal("timeout")
		}
	}
	for _, host := range []string{"127.0.0.1", "127.0.0.2", "::1", "::ffff:127.0.0.1"} {
		g, e := CheckDockPorts(host, nil, base)
		if e != nil || len(g) != 0 {
			t.Fatal("empty", host, e)
		}
	}
	if calls.Load() != 0 {
		t.Fatal("dial before validation")
	}
}
func TestHiddenPortsSeededOutcomes(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(316))
	ordinary := errors.New("text says timeout but is ordinary")
	timed := fmt.Errorf("wrapped: %w", dockPortTimeout{})
	for _, host := range []string{"127.0.0.2", "::1"} {
		for trial := 0; trial < 12; trial++ {
			ports := rng.Perm(64)[:1+rng.Intn(64)]
			for i := range ports {
				ports[i]++
			}
			ports[0] = 65535
			before := slices.Clone(ports)
			conns := make([]dockPortConn, len(ports))
			index := map[int]int{}
			for i, p := range ports {
				index[p] = i
			}
			counts := make([]atomic.Int32, len(ports))
			o := PortOptions{Workers: 1 + rng.Intn(16), Timeout: 27 * time.Millisecond}
			o.Dial = func(network, address string, timeout time.Duration) (net.Conn, error) {
				h, p, e := net.SplitHostPort(address)
				if e != nil || h != host || network != "tcp" || timeout != o.Timeout {
					t.Error("dial arguments", address)
				}
				port, _ := strconv.Atoi(p)
				i, ok := index[port]
				if !ok {
					t.Error("unexpected port")
					return nil, ordinary
				}
				counts[i].Add(1)
				switch i % 5 {
				case 0:
					return &conns[i], nil
				case 1:
					return nil, ordinary
				case 2:
					return &conns[i], timed
				case 3:
					return nil, fmt.Errorf("wrapped: %w", context.DeadlineExceeded)
				default:
					return nil, nil
				}
			}
			g, e := CheckDockPorts(host, ports, o)
			if e != nil || len(g) != len(ports) || !slices.Equal(ports, before) {
				t.Fatal("results", e)
			}
			for i, s := range g {
				state := "closed"
				var cause error
				closed := int32(0)
				switch i % 5 {
				case 0:
					state = "open"
					closed = 1
				case 1:
					cause = ordinary
				case 2:
					state = "timeout"
					cause = timed
					closed = 1
				case 3:
					state = "timeout"
					cause = context.DeadlineExceeded
				case 4:
					cause = ErrPortDial
				}
				if s.Port != ports[i] || s.State != state || !errors.Is(s.Err, cause) || counts[i].Load() != 1 || conns[i].closes.Load() != closed {
					t.Fatal("outcome", i, s, counts[i].Load(), conns[i].closes.Load())
				}
			}
		}
	}
}
func TestHiddenPortsActualConcurrencyBound(t *testing.T) {
	labkit.LeakCheck(t)
	for _, tc := range []struct{ jobs, workers int }{{1, 16}, {16, 1}, {32, 7}, {64, 16}} {
		ports := make([]int, tc.jobs)
		for i := range ports {
			ports[i] = i + 1
		}
		entered := make(chan struct{}, tc.jobs)
		release := make(chan struct{})
		var once sync.Once
		t.Cleanup(func() { once.Do(func() { close(release) }) })
		var active, peak atomic.Int32
		dial := func(string, string, time.Duration) (net.Conn, error) {
			n := active.Add(1)
			for old := peak.Load(); n > old; old = peak.Load() {
				if peak.CompareAndSwap(old, n) {
					break
				}
			}
			entered <- struct{}{}
			<-release
			active.Add(-1)
			return nil, errors.New("closed")
		}
		done := make(chan error, 1)
		go func() {
			_, e := CheckDockPorts("127.0.0.1", ports, PortOptions{Workers: tc.workers, Timeout: time.Second, Dial: dial})
			done <- e
		}()
		expected := min(tc.jobs, tc.workers)
		for i := 0; i < expected; i++ {
			select {
			case <-entered:
			case <-time.After(3 * time.Second):
				t.Fatal("failed to reach worker peak")
			}
		}
		if peak.Load() != int32(expected) {
			t.Fatal("unbounded workers", peak.Load(), expected)
		}
		once.Do(func() { close(release) })
		select {
		case e := <-done:
			if e != nil {
				t.Fatal(e)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("workers not joined")
		}
		if active.Load() != 0 || peak.Load() > int32(expected) {
			t.Fatal("bound/join")
		}
	}
}
func TestHiddenPortsDefaultDialLoopback(t *testing.T) {
	labkit.LeakCheck(t)
	listener, err := net.Listen("tcp", labkit.FreeAddr())
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			t.Skip("sandbox denies loopback TCP bind")
		}
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	done := make(chan error, 1)
	go func() {
		c, e := listener.Accept()
		if c != nil {
			c.Close()
		}
		done <- e
	}()
	g, e := CheckDockPorts("127.0.0.1", []int{port}, PortOptions{Workers: 1, Timeout: time.Second})
	if e != nil || len(g) != 1 || g[0].State != "open" || g[0].Err != nil {
		t.Fatal("default dial", g, e)
	}
	select {
	case e = <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("accept did not stop")
	}
}
`
