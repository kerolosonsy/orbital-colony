package levels

func init() {
	add(&Level{ID: "32-1-websocket", Chapter: 32, Kind: "lab", Title: "بث حي", Goal: "رحّل قراءات الحساسات بـ WebSocket وحلقات بتقف من غير تسريب",
		Requires: []string{"sel:websocket.Accept"}, Deps: []string{"github.com/coder/websocket"}, Bridge: []int{74}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch321Starter}},
		Lesson: `# بث قراءات المستعمرة

## قبل ما تكتب
- مين هيفضل يقرا والـ ping مستني pong؟
- إزاي الحلقتين يقفوا لو الطرف التاني اختفى؟

## مصافحة HTTP وبعدها رسائل
WebSocket بيبدأ بمصافحة HTTP Upgrade. استورد ~github.com/coder/websocket~. ~websocket.Accept~ بيتحقق من المصافحة ويرجع ~*websocket.Conn~؛ لو فشل بيكتب خطأ HTTP بنفسه. افحص سياسة التطبيق قبل القبول. سيب فحص نفس المضيف الافتراضي جنب فحص الـ origin بالظبط؛ متشغّلش ~InsecureSkipVerify~.
~~~go
c, err := websocket.Accept(w, r, nil)
if err != nil { return err }
defer c.CloseNow()
~~~
~Origin~ عنوان الصفحة اللي فتحت الاتصال، مش إثبات هوية. العنوان الموثوق محقون عشان التست يستخدم بورت loopback عشوائي. قارن النص كله، بالنوع والبورت، من غير wildcards. المستدعي بيدي عنوان HTTP صحيح مطابق لمضيف الطلب.
~~~go
if r.Header.Get("Origin") != trusted { http.Error(w, "origin rejected", 403); return }
~~~

## حلقة قراية وحلقة كتابة
الـ pump مجرد حلقة بمسؤولية واحدة. ~Read(ctx)~ بترجع الرسالة كاملة بعد تجميع أجزائها، ونوعها ~MessageType~، وخطأ. ~MessageText~ و ~MessageBinary~ بيفرّقوا الصيغ. حلقة القراية تبعت الرسائل أو الفشل على channel محدود؛ الإرسال لازم يختار الإلغاء كمان عشان الكاتب لو وقف القارئ ميتحبسش.
~~~go
kind, payload, err := c.Read(ctx)
select { case inbox <- packet: case <-ctx.Done(): return }
~~~
حلقة الكتابة تملك كتابة الرسائل والـ ping. ~Write(ctx, websocket.MessageText, payload)~ بتبعت رسالة نصية واحدة. استخدم حلقة كاتب واحدة بتختار بين الرسائل ونبضات الفحص والإلغاء. الكاتب بس يبعت قفل التطبيق؛ صاحب دورة الحياة يقدر يقفل بالقوة عشان يفك الـ I/O.
~~~go
err := c.Write(ctx, websocket.MessageText, []byte("sensor online"))
~~~

## الحدود والتحقق والقفل
~SetReadLimit~ بتحد الرسالة المجمّعة بالبايت قبل تحليلها؛ الزيادة بتنتج كود قفل 1009 وهو ~StatusMessageTooBig~. رسالة على الحد بالظبط مسموحة. ~utf8.Valid~ بتفحص بايتات النص؛ الطول بالبايت مش بعدد الحروف.
~~~go
c.SetReadLimit(256)
valid := utf8.Valid([]byte("قمر")) // true
~~~
~Close(code, reason)~ بتعمل مصافحة قفل بمهلة محدودة. ~CloseStatus(err)~ بتطلع كود القفل، أو -1 لو خطأ عادي. ~StatusNormalClosure~ هو 1000، و ~StatusUnsupportedData~ هو 1003، و ~StatusInvalidFramePayloadData~ هو 1007، و ~StatusPolicyViolation~ هو 1008. السبب يكون عام. ~CloseNow~ قفل بالقوة من غير مصافحة. دايماً الغي السياق واقفل بالقوة واستنى القارئ قبل الرجوع؛ قفل الطرف التاني الطبيعي نجاح والباقي خطأ.
~~~go
if websocket.CloseStatus(err) == websocket.StatusNormalClosure { return nil }
c.Close(websocket.StatusPolicyViolation, "invalid telemetry")
~~~

## النجمة التالتة: نبضة محقونة
~Ping(ctx)~ بتبعت ping وتستنى الـ pong المطابق. القارئ النشط بيعالج رسائل التحكم ويرد على ping الطرف التاني تلقائياً. من غير قارئ الـ ping ممكن يستنى لحد الإلغاء. ~Pulse~ فيها سياق الفحص وchannel اسمه ~Reply~ بسعة؛ مفيش ساعة ولا نوم جوه المرحّل.
~~~go
err := c.Ping(pulse.Context)
pulse.Reply <- err
~~~
~RelayOptions.Pulses~ ممكن تكون ~nil~؛ كده ذراع الـ select يتعطل. لو channel اتقفلت حوّلها لـ ~nil~ عشان متلفش. المستدعي بيلغي سياق أي نبضة لسه شغالة وقت الإغلاق؛ ده بيحد انتظار الـ ping لو سياق المرحّل اتلغى.

## أخطاء شائعة
~~~go
c.Ping(context.Background()) // غلط قبل تشغيل قارئ: مفيش حد بيستهلك pong
for { c.Read(ctx) }          // غلط: تجاهل الخطأ يعمل دوران بعد الفصل
~~~

## المهمة
نفّذ ~RelayTelemetry(ctx context.Context, w http.ResponseWriter, r *http.Request, opts RelayOptions) error~ لوصلة حساسات المستعمرة. ارفض ~Limit~ برا 1 لحد 4096 أو ~Origin~ فاضي بـ HTTP 400 وخطأ قبل المصافحة؛ بعد كده ارفض Origin الطلب الغايب أو المختلف بـ 403 وخطأ. سيب فحوصات المصافحة الافتراضية. جسم خطأ HTTP مش محدد.

بعد القبول استخدم حلقتين منفصلتين للقراية والكتابة وchannel رسائل محدود، وسياق محقون للـ I/O. كل رسالة نصية مش فاضية وUTF-8 صحيح وطولها لحد ~Limit~ تاخد رد واحد: ~relay:~ وبعدها نفس البايتات بالترتيب. binary جوه الحد يقفل بـ 1003، ونص UTF-8 بايظ بـ 1007، والنص الفاضي بـ 1008، وأي رسالة أكبر من الحد بـ 1009، من غير رد تطبيقي للرسالة المرفوضة. المخالفات وأخطاء النقل ترجع خطأ؛ القفل الطبيعي يرجع ~nil~. الإلغاء يرجع خطأ ويفصل بسرعة. مترجعش غير بعد توقف الحلقتين؛ مفيش ~Report~ و ~main~ فاضية.

النجمة التالتة: استهلك ~Pulses~ في الكاتب، ونادِ ~Ping~ بسياق النبضة، وابعت نفس الخطأ، حتى ~nil~، مرة واحدة لـ ~Reply~، ووقف الحلقتين لو الـ ping فشل. الـ channel ليها مستهلك واحد وكل ~Reply~ سعتها 1، ومتقفلش قنوات المستدعي. قفل قناة النبض يعطّل الفحص. النداءات المتزامنة مستقلة.
`,
		Hints: []string{"افحص الإعدادات والـ Origin قبل websocket.Accept، وبعدها اضبط SetReadLimit.", "خلّي القارئ يبعت رسالة أو خطأ على channel محدود بـ select على الإلغاء؛ الكاتب وحده يرد ويقفل بمبرر.", "عند الخروج الغي السياق واقفل CloseNow واستنى القارئ؛ للنجمة نفّذ Ping من الكاتب وابعت نتيجته لـ Reply."},
		En: &LevelText{Title: "Live Broadcast", Goal: "Relay sensor readings through WebSocket pumps that stop without leaks", Files: []LabFile{{Name: "main.go", Starter: ch321Starter}},
			Lesson: `# Live Telemetry Relay

## Before you write
- Who keeps reading while a ping waits for its pong?
- How do both pumps stop when the peer disappears?

## An HTTP handshake, then messages
A WebSocket starts with an HTTP Upgrade handshake. Import ~github.com/coder/websocket~. ~websocket.Accept~ validates the upgrade and returns a ~*websocket.Conn~; failure already writes the HTTP error. Check your application policy before accepting. Keep the library's same-host origin check as well as our exact origin check; never enable ~InsecureSkipVerify~.
~~~go
c, err := websocket.Accept(w, r, nil)
if err != nil { return err }
defer c.CloseNow()
~~~
~Origin~ is the browser's page origin, not authentication. An injected trusted origin lets the test use its random loopback port. Require exact equality, including scheme and port, without wildcards. The caller supplies a valid HTTP origin matching the request host.
~~~go
if r.Header.Get("Origin") != trusted { http.Error(w, "origin rejected", 403); return }
~~~

## A read pump and a write pump
A pump is just a loop with one job. ~Read(ctx)~ returns a complete reassembled message, its ~MessageType~, and an error. ~MessageText~ and ~MessageBinary~ distinguish application formats. A read pump sends messages or failures to a bounded channel; it must select on cancellation when sending so a stopped writer cannot strand it.
~~~go
kind, payload, err := c.Read(ctx)
select { case inbox <- packet: case <-ctx.Done(): return }
~~~
The write pump owns application writes and heartbeat pings. ~Write(ctx, websocket.MessageText, payload)~ sends one text message. Use one writer loop selecting between the inbox, heartbeat requests, and cancellation. Only the writer sends application close frames; the lifecycle owner may force-close to unblock I/O.
~~~go
err := c.Write(ctx, websocket.MessageText, []byte("sensor online"))
~~~

## Limits, validation, and closure
~SetReadLimit~ caps a reassembled message in bytes before application parsing; exceeding it yields close code 1009 (~StatusMessageTooBig~). A message exactly at the limit is accepted. ~utf8.Valid~ checks text bytes; text length counts bytes, not runes.
~~~go
c.SetReadLimit(256)
valid := utf8.Valid([]byte("moon")) // true
~~~
~Close(code, reason)~ performs the bounded close handshake. ~CloseStatus(err)~ extracts a peer's close code, or -1 for an ordinary error. ~StatusNormalClosure~ is 1000; ~StatusUnsupportedData~ is 1003; ~StatusInvalidFramePayloadData~ is 1007; ~StatusPolicyViolation~ is 1008. Use a generic reason. ~CloseNow~ force-closes without a handshake. Always cancel, force-close, and join the reader before returning; normal peer closure is success, other failures are errors.
~~~go
if websocket.CloseStatus(err) == websocket.StatusNormalClosure { return nil }
c.Close(websocket.StatusPolicyViolation, "invalid telemetry")
~~~

## Third star: an injected heartbeat
~Ping(ctx)~ sends a ping and waits for the matching pong. An active reader processes control frames and automatically replies to peer pings. Without an active reader, ping can hang until cancellation. ~Pulse~ supplies the ping context and a buffered ~Reply~ channel; no clock or sleep belongs in this relay.
~~~go
err := c.Ping(pulse.Context)
pulse.Reply <- err
~~~
~RelayOptions.Pulses~ may be nil; a nil channel disables that select arm. Set a closed pulse channel to nil so it cannot spin. The caller cancels every outstanding pulse context when shutting down; this bounds a ping while the relay context is canceled.

## Common mistakes
~~~go
c.Ping(context.Background()) // Wrong before starting a reader: nobody consumes pong
for { c.Read(ctx) }          // Wrong: errors ignored, spins after disconnect
~~~

## Your mission
Implement ~RelayTelemetry(ctx context.Context, w http.ResponseWriter, r *http.Request, opts RelayOptions) error~ for a colony sensor uplink. Reject ~Limit~ outside 1..4096 or empty ~Origin~ with HTTP 400 and an error before upgrade; otherwise reject missing or unequal request Origin with 403 and an error. Keep the library's default handshake checks. No particular HTTP error body is required.

After acceptance use separate read and write pumps and a bounded inbox, with the injected context on I/O. Each nonempty valid UTF-8 text message at or below ~Limit~ receives exactly ~relay:~ followed by its unchanged bytes, in order. Binary within the limit closes with 1003, invalid UTF-8 text with 1007, empty text with 1008, and any oversize message with 1009, with no application reply for that message. Return an error for these violations and transport failures; normal peer close returns nil. Cancellation returns an error and disconnects promptly. Do not return until both pumps have stopped; no ~Report~, and ~main~ is empty.

Third star: consume ~Pulses~ in the writer, call ~Ping~ with each pulse's context, send its exact error (including nil) once to ~Reply~, and stop both pumps on ping failure. The channel has one receiver, every ~Reply~ has capacity 1, and caller-owned channels must not be closed. A closed pulse channel disables heartbeats. Concurrent relay calls share no state.
`,
			Hints: []string{"Validate options and Origin before websocket.Accept, then set SetReadLimit.", "Have the reader send a message or error through a bounded channel selecting on cancellation; only the writer replies and sends policy closes.", "On exit cancel, CloseNow, and join the reader; for the star call Ping from the writer and acknowledge its result through Reply."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch321Hidden},
	})
}

const ch321Starter = `package main

import (
	"context"
	"errors"
	"net/http"
)

type Pulse struct {
	Context context.Context
	Reply   chan error
}
type RelayOptions struct {
	Origin string
	Limit  int
	Pulses <-chan Pulse
}

func RelayTelemetry(ctx context.Context, w http.ResponseWriter, r *http.Request, opts RelayOptions) error {
	http.Error(w, "TODO", http.StatusNotImplemented)
	return errors.New("TODO")
}
func main() {}
`
const ch321Hidden = `package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"orbital/colony/labkit"
)

func relayWait[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("relay did not stop or reply")
	}
	var zero T
	return zero
}
func relayFixture(t *testing.T, limit int) (*websocket.Conn, chan Pulse, context.CancelFunc, <-chan error) {
	t.Helper()
	labkit.LeakCheck(t)
	ctx, cancel := context.WithCancel(context.Background())
	pulses := make(chan Pulse)
	done := make(chan error, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		done <- RelayTelemetry(ctx, w, r, RelayOptions{Origin: "http://" + r.Host, Limit: limit, Pulses: pulses})
	}))
	t.Cleanup(srv.Close)
	t.Cleanup(cancel)
	dialCtx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()
	conn, _, err := websocket.Dial(dialCtx, "ws"+strings.TrimPrefix(srv.URL, "http"), &websocket.DialOptions{HTTPHeader: http.Header{"Origin": {srv.URL}}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.CloseNow() })
	return conn, pulses, cancel, done
}
func TestHiddenRelayEchoBoundaries(t *testing.T) {
	conn, _, cancel, done := relayFixture(t, 128)
	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()
	rng := rand.New(rand.NewSource(3201))
	for i := 0; i < 24; i++ {
		text := fmt.Sprintf("حساس-%d-🚀", rng.Intn(100000))
		if i == 0 {
			text = strings.Repeat("x", 128)
		}
		if err := conn.Write(ctx, websocket.MessageText, []byte(text)); err != nil {
			t.Fatal(err)
		}
		kind, got, err := conn.Read(ctx)
		if err != nil || kind != websocket.MessageText || string(got) != "relay:"+text {
			t.Fatalf("echo: %q %v", got, err)
		}
	}
	conn.Close(websocket.StatusNormalClosure, "")
	if err := relayWait(t, done); err != nil {
		t.Fatal(err)
	}
	cancel()
}
func TestHiddenRelayRejectsFrames(t *testing.T) {
	for _, tc := range []struct {
		name string
		kind websocket.MessageType
		body []byte
		code websocket.StatusCode
	}{
		{"binary", websocket.MessageBinary, []byte("probe"), websocket.StatusUnsupportedData},
		{"utf8", websocket.MessageText, []byte{0xff}, websocket.StatusInvalidFramePayloadData},
		{"empty", websocket.MessageText, nil, websocket.StatusPolicyViolation},
		{"oversize", websocket.MessageText, []byte(strings.Repeat("z", 33)), websocket.StatusMessageTooBig},
	} {
		t.Run(tc.name, func(t *testing.T) {
			conn, _, _, done := relayFixture(t, 32)
			ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
			defer stop()
			if err := conn.Write(ctx, tc.kind, tc.body); err != nil {
				t.Fatal(err)
			}
			_, _, err := conn.Read(ctx)
			if websocket.CloseStatus(err) != tc.code {
				t.Fatalf("close=%v want %v: %v", websocket.CloseStatus(err), tc.code, err)
			}
			if relayWait(t, done) == nil {
				t.Fatal("bad frame accepted")
			}
		})
	}
}
func TestHiddenRelayRejectsHandshake(t *testing.T) {
	for _, tc := range []struct {
		limit           int
		trusted, origin string
		status          int
	}{
		{0, "http://station", "http://station", 400},
		{-1, "http://station", "http://station", 400},
		{4097, "http://station", "http://station", 400},
		{8, "", "http://station", 400},
		{8, "http://station", "", 403},
		{8, "http://station", "http://other", 403},
		{8, "http://station", "http://station", 426},
	} {
		req := httptest.NewRequest("GET", "http://station/", nil)
		req.Header.Set("Origin", tc.origin)
		w := httptest.NewRecorder()
		if err := RelayTelemetry(context.Background(), w, req, RelayOptions{Origin: tc.trusted, Limit: tc.limit}); err == nil || w.Code != tc.status {
			t.Fatalf("invalid options, origin, or upgrade accepted: status=%d want=%d err=%v", w.Code, tc.status, err)
		}
	}
}
func TestHiddenRelayCancellation(t *testing.T) {
	conn, _, cancel, done := relayFixture(t, 4096)
	cancel()
	if relayWait(t, done) == nil {
		t.Fatal("cancellation swallowed")
	}
	conn.CloseNow()
}
func TestHiddenBonusRelayHeartbeat(t *testing.T) {
	conn, pulses, _, done := relayFixture(t, 64)
	ctx, stop := context.WithTimeout(context.Background(), 3*time.Second)
	defer stop()
	reading := make(chan error, 1)
	go func() { _, _, err := conn.Read(ctx); reading <- err }()
	ack := make(chan error, 1)
	select {
	case pulses <- Pulse{Context: ctx, Reply: ack}:
	case <-ctx.Done():
		t.Fatal("writer absent")
	}
	if err := relayWait(t, ack); err != nil {
		t.Fatal(err)
	}
	failed, cancel := context.WithCancel(context.Background())
	cancel()
	select {
	case pulses <- Pulse{Context: failed, Reply: ack}:
	case <-ctx.Done():
		t.Fatal("writer absent")
	}
	if err := relayWait(t, ack); err == nil {
		t.Fatal("ping failure lost")
	}
	if relayWait(t, done) == nil {
		t.Fatal("failed ping did not stop pumps")
	}
	relayWait(t, reading)
}
func TestHiddenRelayConcurrentDisconnects(t *testing.T) {
	labkit.LeakCheck(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	completed := make(chan error, 8)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		completed <- RelayTelemetry(ctx, w, r, RelayOptions{Origin: "http://" + r.Host, Limit: 64})
	}))
	defer srv.Close()
	finished := make(chan error, 8)
	for i := 0; i < 8; i++ {
		go func(i int) {
			conn, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(srv.URL, "http"), &websocket.DialOptions{HTTPHeader: http.Header{"Origin": {srv.URL}}})
			if err != nil {
				finished <- err
				return
			}
			defer conn.CloseNow()
			text := fmt.Sprintf("sensor-%d", i)
			if err = conn.Write(ctx, websocket.MessageText, []byte(text)); err != nil {
				finished <- err
				return
			}
			_, got, err := conn.Read(ctx)
			if err == nil && string(got) != "relay:"+text {
				err = fmt.Errorf("crossed replies: %q", got)
			}
			conn.CloseNow()
			finished <- err
		}(i)
	}
	for i := 0; i < 8; i++ {
		if err := relayWait(t, finished); err != nil {
			t.Fatal(err)
		}
		if err := relayWait(t, completed); err == nil {
			t.Fatal("abrupt disconnect reported clean close")
		}
	}
}
`

const ch322Protocol0 = `syntax = "proto3";
package colony.berth;
option go_package = "orbital.lab/protocol;main";

message BerthQuery { string sector = 1; }
message BerthReading { string sector = 1; int64 pressure = 2; }
service BerthControl { rpc Inspect(BerthQuery) returns (BerthReading); }
`

const ch322Protocol1 = `// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.12
// 	protoc        v7.36.2
// source: berth.proto

package main

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BerthQuery struct {
	state         protoimpl.MessageState ` + "`" + `protogen:"open.v1"` + "`" + `
	Sector        string                 ` + "`" + `protobuf:"bytes,1,opt,name=sector,proto3" json:"sector,omitempty"` + "`" + `
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BerthQuery) Reset() {
	*x = BerthQuery{}
	mi := &file_berth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BerthQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BerthQuery) ProtoMessage() {}

func (x *BerthQuery) ProtoReflect() protoreflect.Message {
	mi := &file_berth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BerthQuery.ProtoReflect.Descriptor instead.
func (*BerthQuery) Descriptor() ([]byte, []int) {
	return file_berth_proto_rawDescGZIP(), []int{0}
}

func (x *BerthQuery) GetSector() string {
	if x != nil {
		return x.Sector
	}
	return ""
}

type BerthReading struct {
	state         protoimpl.MessageState ` + "`" + `protogen:"open.v1"` + "`" + `
	Sector        string                 ` + "`" + `protobuf:"bytes,1,opt,name=sector,proto3" json:"sector,omitempty"` + "`" + `
	Pressure      int64                  ` + "`" + `protobuf:"varint,2,opt,name=pressure,proto3" json:"pressure,omitempty"` + "`" + `
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BerthReading) Reset() {
	*x = BerthReading{}
	mi := &file_berth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BerthReading) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BerthReading) ProtoMessage() {}

func (x *BerthReading) ProtoReflect() protoreflect.Message {
	mi := &file_berth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BerthReading.ProtoReflect.Descriptor instead.
func (*BerthReading) Descriptor() ([]byte, []int) {
	return file_berth_proto_rawDescGZIP(), []int{1}
}

func (x *BerthReading) GetSector() string {
	if x != nil {
		return x.Sector
	}
	return ""
}

func (x *BerthReading) GetPressure() int64 {
	if x != nil {
		return x.Pressure
	}
	return 0
}

var File_berth_proto protoreflect.FileDescriptor

const file_berth_proto_rawDesc = "" +
	"\n" +
	"\vberth.proto\x12\fcolony.berth\"$\n" +
	"\n" +
	"BerthQuery\x12\x16\n" +
	"\x06sector\x18\x01 \x01(\tR\x06sector\"B\n" +
	"\fBerthReading\x12\x16\n" +
	"\x06sector\x18\x01 \x01(\tR\x06sector\x12\x1a\n" +
	"\bpressure\x18\x02 \x01(\x03R\bpressure2O\n" +
	"\fBerthControl\x12?\n" +
	"\aInspect\x12\x18.colony.berth.BerthQuery\x1a\x1a.colony.berth.BerthReadingB\x1bZ\x19orbital.lab/protocol;mainb\x06proto3"

var (
	file_berth_proto_rawDescOnce sync.Once
	file_berth_proto_rawDescData []byte
)

func file_berth_proto_rawDescGZIP() []byte {
	file_berth_proto_rawDescOnce.Do(func() {
		file_berth_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_berth_proto_rawDesc), len(file_berth_proto_rawDesc)))
	})
	return file_berth_proto_rawDescData
}

var file_berth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_berth_proto_goTypes = []any{
	(*BerthQuery)(nil),   // 0: colony.berth.BerthQuery
	(*BerthReading)(nil), // 1: colony.berth.BerthReading
}
var file_berth_proto_depIdxs = []int32{
	0, // 0: colony.berth.BerthControl.Inspect:input_type -> colony.berth.BerthQuery
	1, // 1: colony.berth.BerthControl.Inspect:output_type -> colony.berth.BerthReading
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_berth_proto_init() }
func file_berth_proto_init() {
	if File_berth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_berth_proto_rawDesc), len(file_berth_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_berth_proto_goTypes,
		DependencyIndexes: file_berth_proto_depIdxs,
		MessageInfos:      file_berth_proto_msgTypes,
	}.Build()
	File_berth_proto = out.File
	file_berth_proto_goTypes = nil
	file_berth_proto_depIdxs = nil
}
`

const ch322Protocol2 = `// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v7.36.2
// source: berth.proto

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BerthControl_Inspect_FullMethodName = "/colony.berth.BerthControl/Inspect"
)

// BerthControlClient is the client API for BerthControl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BerthControlClient interface {
	Inspect(ctx context.Context, in *BerthQuery, opts ...grpc.CallOption) (*BerthReading, error)
}

type berthControlClient struct {
	cc grpc.ClientConnInterface
}

func NewBerthControlClient(cc grpc.ClientConnInterface) BerthControlClient {
	return &berthControlClient{cc}
}

func (c *berthControlClient) Inspect(ctx context.Context, in *BerthQuery, opts ...grpc.CallOption) (*BerthReading, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BerthReading)
	err := c.cc.Invoke(ctx, BerthControl_Inspect_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BerthControlServer is the server API for BerthControl service.
// All implementations must embed UnimplementedBerthControlServer
// for forward compatibility.
type BerthControlServer interface {
	Inspect(context.Context, *BerthQuery) (*BerthReading, error)
	mustEmbedUnimplementedBerthControlServer()
}

// UnimplementedBerthControlServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBerthControlServer struct{}

func (UnimplementedBerthControlServer) Inspect(context.Context, *BerthQuery) (*BerthReading, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Inspect not implemented")
}
func (UnimplementedBerthControlServer) mustEmbedUnimplementedBerthControlServer() {}
func (UnimplementedBerthControlServer) testEmbeddedByValue()                      {}

// UnsafeBerthControlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BerthControlServer will
// result in compilation errors.
type UnsafeBerthControlServer interface {
	mustEmbedUnimplementedBerthControlServer()
}

func RegisterBerthControlServer(s grpc.ServiceRegistrar, srv BerthControlServer) {
	// If the following call pancis, it indicates UnimplementedBerthControlServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BerthControl_ServiceDesc, srv)
}

func _BerthControl_Inspect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BerthQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BerthControlServer).Inspect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BerthControl_Inspect_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BerthControlServer).Inspect(ctx, req.(*BerthQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// BerthControl_ServiceDesc is the grpc.ServiceDesc for BerthControl service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BerthControl_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "colony.berth.BerthControl",
	HandlerType: (*BerthControlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Inspect",
			Handler:    _BerthControl_Inspect_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "berth.proto",
}
`

func init() {
	add(&Level{ID: "32-2-grpc-unary", Chapter: 32, Kind: "lab", Title: "gRPC", Goal: "افحص ضغط قطاع بخدمة unary وحارس تصريح وحالات واضحة",
		Requires: []string{"sel:grpc.NewServer"}, Deps: []string{"google.golang.org/grpc", "google.golang.org/protobuf"}, Bridge: []int{75}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch322Starter}, {Name: "berth.proto", Starter: ch322Protocol0, ReadOnly: true}, {Name: "berth.pb.go", Starter: ch322Protocol1, ReadOnly: true}, {Name: "berth_grpc.pb.go", Starter: ch322Protocol2, ReadOnly: true}},
		Lesson: `# فحص ضغط القطاع بـ gRPC

## قبل ما تكتب
- إيه الفشل اللي سببه العميل وإيه اللي سببه المخزن؟
- فين تقدر تمنع الطلب قبل ما يوصل للخدمة؟

## Protobuf هو عقد النقل
Protobuf بيوصف رسائل بأنواع وحقول مرقّمة ثابتة. الـ unary RPC بياخد رسالة ويرجع رسالة. ملف ~berth.proto~ المقفول بيعرّف ~BerthQuery~ و ~BerthReading~ و ~BerthControl.Inspect~. متستخدمش رقم حقل قديم بمعنى جديد. ~berth.pb.go~ و ~berth_grpc.pb.go~ متولّدين بـ protoc؛ اللاعب مش محتاج يولّدهم.
~~~proto
message ProbeQuery { string name = 1; }
message ProbeReply { int64 depth = 1; }
service ProbeCatalog { rpc Lookup(ProbeQuery) returns (ProbeReply); }
~~~
أنواع Go المتولّدة فيها حقول عامة وgetters آمنة مع ~nil~. ~GetSector~ و ~GetPressure~ بيرجعوا القيمة الصفرية لو المستقبل ~nil~. الرسائل pointers؛ متغيّرش الرسالة بعد تسليمها للنقل.
~~~go
q := &BerthQuery{Sector: "ring"}
r := &BerthReading{Sector: q.GetSector(), Pressure: 900}
fmt.Println(r.GetPressure())
~~~

## تنفيذ وتسجيل الخدمة
ضمّن ~UnimplementedBerthControlServer~ بالقيمة عشان التوافق المستقبلي، ونفّذ ~Inspect(context.Context, *BerthQuery) (*BerthReading, error)~. الطلبات ممكن تشتغل بالتزامن؛ حالة الطلب محلية. ~PressureLookup~ المحقونة آمنة بالتزامن.
~~~go
type stationService struct {
 UnimplementedBerthControlServer
 lookup PressureLookup
}
~~~
~grpc.NewServer~ بيعمل سيرفر. ~RegisterBerthControlServer~ بيربط مسارات الخدمة المتولّدة. ~Serve(listener)~ بتستنى لحد القفل. ~Stop~ بتقطع الطلبات، و ~GracefulStop~ بتستنى الجاري؛ التست لازم يفك أي حاجز قبل انتظار القفل المرتب.
~~~go
srv := grpc.NewServer()
RegisterBerthControlServer(srv, service)
go srv.Serve(listener)
defer srv.Stop()
~~~

## أكواد الحالة جزء من العقد
خطأ Go عادي بيتحوّل لـ ~codes.Unknown~ على السلك. استخدم ~status.Error(code, message)~ للفشل المقصود؛ العميل يفحص بـ ~status.Code(err)~. ~InvalidArgument~ مدخل بايظ، و ~NotFound~ قطاع صحيح مش موجود، و ~Unauthenticated~ مفتاح مرفوض، و ~Internal~ بيخفي تفاصيل المخزن. الإلغاء والمهلة يتحوّلوا بـ ~status.FromContextError(err).Err()~ لـ ~Canceled~ أو ~DeadlineExceeded~.
~~~go
return nil, status.Error(codes.NotFound, "sector missing")
// عند حدود السياق:
return nil, status.FromContextError(ctx.Err()).Err()
~~~
افحص السياق قبل البحث وبعده. المخزن ممكن يلف ~ErrSectorMissing~؛ استخدم ~errors.Is~. متطلعش أخطائه التانية للعميل.

## Interceptor حوالين الطلب
~grpc.UnaryServerInterceptor~ بياخد السياق والطلب و ~*grpc.UnaryServerInfo~ و ~grpc.UnaryHandler~. ~next(ctx, req)~ بتشغّل الخدمة؛ الرجوع قبلها بيرفض الطلب. ركّبه مرة بـ ~grpc.UnaryInterceptor~ وقت بناء السيرفر.
~~~go
guard := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
 return next(ctx, req)
}
srv := grpc.NewServer(grpc.UnaryInterceptor(guard))
~~~
~metadata.FromIncomingContext~ بتقرا metadata الطلب؛ ~md.Get~ بترجع كل قيم المفتاح. اطلب قيمة واحدة بالظبط عشان الاعتمادات المكررة متبقاش مبهمة. العميل بيستخدم ~metadata.NewOutgoingContext~ و ~metadata.Pairs~.
~~~go
md, _ := metadata.FromIncomingContext(ctx)
keys := md.Get("x-colony-key")
ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-colony-key", "example"))
~~~

## اختبار من غير بورت
~bufconn.Listen~ مستمع ~net.Listener~ في الذاكرة. ~grpc.NewClient~ بيعمل اتصال كسول؛ ~WithContextDialer~ بتحقن النقل. ~insecure.NewCredentials~ بتلغي تشفير النقل للتست في الذاكرة بس؛ ده منفصل عن تصريح التطبيق. ~NewBerthControlClient~ المتولّدة بتوفّر ~Inspect~. اقفل العميل والسيرفر والمستمع.
~~~go
ln := bufconn.Listen(1 << 20)
cc, err := grpc.NewClient("passthrough:///lab",
 grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return ln.DialContext(ctx) }),
 grpc.WithTransportCredentials(insecure.NewCredentials()))
client := NewBerthControlClient(cc)
reply, err := client.Inspect(ctx, &BerthQuery{Sector: "ring"})
~~~

## النجمة التالتة: استرجاع عند الحدود
لو البحث المحقون عمل panic، مينفعش يوقع العملية. استخدم استرجاع مؤجل في الـ interceptor مع قيم رجوع مسمّاة. الاسترجاع يصفّر الرد ويحوّل الخطأ لـ Internal عام؛ مينفعش يحوّل الـ panic لنجاح.
~~~go
func guarded() (reply any, err error) {
 defer func() { if recover() != nil { reply = nil; err = status.Error(codes.Internal, "unavailable") } }()
 return work()
}
~~~

## أخطاء شائعة
~~~go
return nil, err // غلط لخطأ مخزن خاص: بيسرب تفاصيل مع Unknown
// غلط: استخدام أول قيمة بس بيقبل اعتماد مكرر.
if md.Get("x-colony-key")[0] == key { return next(ctx, req) }
~~~

## المهمة
نفّذ ~NewInspectionServer(key string, lookup PressureLookup) *grpc.Server~. المستدعي بيدي مفتاح مش فاضي وبحث مش ~nil~ وآمن بالتزامن. رجّع سيرفر مسجّل عليه ~BerthControl~ مع unary interceptor للتصريح. لازم قيمة ~x-colony-key~ واحدة بالظبط وتساوي ~key~؛ غير كده ~Unauthenticated~ من غير بحث. كل طلب مستقل.

~Inspect~ بيقبل قطاع من 1 لحد 16 بايت ASCII: أول حرف إنجليزي صغير وبعده حروف صغيرة أو أرقام أو شرطات. متعملش trim ولا تغيير. المدخل الغلط ~InvalidArgument~ من غير بحث. البحث بياخد سياق الطلب ونفس اسم القطاع. النجاح يرجع القطاع والضغط int64 زي ما هما، حتى الصفر والسالب. ~ErrSectorMissing~ الملفوف يرجع ~NotFound~؛ أي خطأ مخزن تاني ~Internal~ ورسالة ~inspection unavailable~. افحص الإلغاء قبل وبعد البحث وحوّل خطأ السياق لحالة gRPC. متسرّبش نص المخزن. مفيش ~Report~ و ~main~ فاضية.

النجمة التالتة: استرجع panic الخدمة في الـ interceptor ورد ~nil~ وخطأ ~Internal~ برسالة ~inspection unavailable~ بالظبط. متعدّلش البروتوكول ولا الملفات المتولّدة المقفولة.
`,
		Hints: []string{"ضمّن UnimplementedBerthControlServer بالقيمة، وسجّل تنفيذ Inspect على grpc.NewServer.", "ركّب UnaryInterceptor بيفحص عدد قيم x-colony-key قبل next؛ الخدمة تتحقق من القطاع قبل البحث.", "اربط errors.Is بالحالات، وافحص ctx.Err قبل وبعد البحث؛ للنجمة استخدم defer وrecover وقيم رجوع مسمّاة في الـ interceptor."},
		En: &LevelText{Title: "gRPC", Goal: "Inspect sector pressure with a unary service, an authentication interceptor, and explicit statuses", Files: []LabFile{{Name: "main.go", Starter: ch322Starter}, {Name: "berth.proto", Starter: ch322Protocol0, ReadOnly: true}, {Name: "berth.pb.go", Starter: ch322Protocol1, ReadOnly: true}, {Name: "berth_grpc.pb.go", Starter: ch322Protocol2, ReadOnly: true}},
			Lesson: `# Inspect a Pressure Sector with gRPC

## Before you write
- Which failures belong to the caller, and which belong to the store?
- Where can authentication stop a call before it reaches the service?

## Protobuf is the wire contract
Protobuf describes typed messages using stable numbered fields. A unary RPC takes one message and returns one message. Our read-only ~berth.proto~ declares ~BerthQuery~, ~BerthReading~, and ~BerthControl.Inspect~. Never reuse a field number for a different meaning. The shipped ~berth.pb.go~ and ~berth_grpc.pb.go~ were generated with protoc; you do not need the generator to play.
~~~proto
message ProbeQuery { string name = 1; }
message ProbeReply { int64 depth = 1; }
service ProbeCatalog { rpc Lookup(ProbeQuery) returns (ProbeReply); }
~~~
The generated Go types have exported fields and nil-safe getters. ~GetSector~ and ~GetPressure~ return zero values for nil receivers. Messages are pointers; do not modify a message after handing it to the transport.
~~~go
q := &BerthQuery{Sector: "ring"}
r := &BerthReading{Sector: q.GetSector(), Pressure: 900}
fmt.Println(r.GetPressure())
~~~

## Implement and register a service
Embed ~UnimplementedBerthControlServer~ by value for forward compatibility, then implement ~Inspect(context.Context, *BerthQuery) (*BerthReading, error)~. Each RPC can run concurrently, so request state belongs in local variables. Our injected ~PressureLookup~ is safe for concurrent use.
~~~go
type stationService struct {
 UnimplementedBerthControlServer
 lookup PressureLookup
}
~~~
~grpc.NewServer~ constructs a server. ~RegisterBerthControlServer~ attaches the generated service routes. ~Serve(listener)~ blocks until shutdown. ~Stop~ closes active calls immediately; ~GracefulStop~ waits for them, so tests must release any barrier before waiting for graceful completion.
~~~go
srv := grpc.NewServer()
RegisterBerthControlServer(srv, service)
go srv.Serve(listener)
defer srv.Stop()
~~~

## Status codes are part of the API
An arbitrary Go error becomes ~codes.Unknown~ on the wire. Use ~status.Error(code, message)~ for deliberate public failures; clients inspect ~status.Code(err)~. ~InvalidArgument~ means malformed input, ~NotFound~ means a valid but absent sector, ~Unauthenticated~ means no valid key, and ~Internal~ hides store details. Context cancellation and deadlines map with ~status.FromContextError(err).Err()~ to ~Canceled~ or ~DeadlineExceeded~.
~~~go
return nil, status.Error(codes.NotFound, "sector missing")
// At a context boundary:
return nil, status.FromContextError(ctx.Err()).Err()
~~~
Check context before lookup and again after lookup. An injected store may wrap ~ErrSectorMissing~; use ~errors.Is~. Do not expose its other errors to the client.

## Interceptors wrap RPCs
~grpc.UnaryServerInterceptor~ receives the context, request, ~*grpc.UnaryServerInfo~, and ~grpc.UnaryHandler~. Calling ~next(ctx, req)~ runs the service. Returning early rejects the RPC. Install it once with ~grpc.UnaryInterceptor~ when constructing the server.
~~~go
guard := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
 return next(ctx, req)
}
srv := grpc.NewServer(grpc.UnaryInterceptor(guard))
~~~
~metadata.FromIncomingContext~ reads request metadata; ~md.Get~ returns all values for a key. Require exactly one value so duplicated credentials cannot be ambiguous. A client uses ~metadata.NewOutgoingContext~ and ~metadata.Pairs~.
~~~go
md, _ := metadata.FromIncomingContext(ctx)
keys := md.Get("x-colony-key")
ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("x-colony-key", "example"))
~~~

## Testing without a port
~bufconn.Listen~ is an in-memory ~net.Listener~. ~grpc.NewClient~ creates a lazy connection; ~WithContextDialer~ supplies its transport. ~insecure.NewCredentials~ disables transport encryption only for this in-memory test; it is unrelated to application authentication. The generated ~NewBerthControlClient~ exposes ~Inspect~. Close client, server, and listener.
~~~go
ln := bufconn.Listen(1 << 20)
cc, err := grpc.NewClient("passthrough:///lab",
 grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return ln.DialContext(ctx) }),
 grpc.WithTransportCredentials(insecure.NewCredentials()))
client := NewBerthControlClient(cc)
reply, err := client.Inspect(ctx, &BerthQuery{Sector: "ring"})
~~~

## Third star: recover at the boundary
A panic in an injected lookup must not crash the process. Use a deferred recovery in the interceptor with named return values. Recover changes the reply to nil and the error to a generic Internal status; it must not turn the panic into success.
~~~go
func guarded() (reply any, err error) {
 defer func() { if recover() != nil { reply = nil; err = status.Error(codes.Internal, "unavailable") } }()
 return work()
}
~~~

## Common mistakes
~~~go
return nil, err // Wrong for a private store error: exposes details with Unknown
// Wrong: using only the first metadata value accepts duplicated credentials.
if md.Get("x-colony-key")[0] == key { return next(ctx, req) }
~~~

## Your mission
Implement ~NewInspectionServer(key string, lookup PressureLookup) *grpc.Server~. The caller supplies a nonempty key and non-nil concurrent-safe lookup. Return a registered ~BerthControl~ server with a unary authentication interceptor. Exactly one incoming ~x-colony-key~ equal to ~key~ is required; otherwise return ~Unauthenticated~ without lookup. Each call is independent.

~Inspect~ accepts a sector of 1..16 ASCII bytes: first a lowercase letter, then lowercase letters, digits, or hyphens. Do not trim or normalize. Invalid input returns ~InvalidArgument~ without lookup. Lookup receives the RPC context and exact sector. Success returns that sector and its int64 pressure unchanged, including zero and negative readings. Wrapped ~ErrSectorMissing~ maps to ~NotFound~; other store errors map to ~Internal~ with message ~inspection unavailable~. Check cancellation before and after lookup, mapping the context error to its gRPC status. Do not leak private store text. No ~Report~; ~main~ is empty.

Third star: recover service panics in the interceptor, returning a nil reply and ~Internal~ with exactly ~inspection unavailable~. Do not edit the read-only protocol or generated files.
`,
			Hints: []string{"Embed UnimplementedBerthControlServer by value and register your Inspect implementation on grpc.NewServer.", "Install a UnaryInterceptor checking the count and value of x-colony-key before next; validate the sector before lookup.", "Map errors.Is to status codes and check ctx.Err before and after lookup; for the star use deferred recover and named returns in the interceptor."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch322Hidden},
	})
}

const ch322Starter = `package main

import (
	"context"
	"errors"
	"google.golang.org/grpc"
)

var ErrSectorMissing = errors.New("sector missing")

type PressureLookup func(context.Context, string) (int64, error)

func NewInspectionServer(key string, lookup PressureLookup) *grpc.Server {
	return grpc.NewServer()
}
func main() {}
`
const ch322Hidden = `package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"orbital/colony/labkit"
)

func inspectionClient(t *testing.T, lookup PressureLookup) BerthControlClient {
	t.Helper()
	labkit.LeakCheck(t)
	ln := bufconn.Listen(1 << 20)
	srv := NewInspectionServer("orbit-32", lookup)
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	conn, err := grpc.NewClient("passthrough:///inspection", grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return ln.DialContext(ctx) }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Close()
		srv.Stop()
		ln.Close()
		select {
		case <-served:
		case <-time.After(3 * time.Second):
			t.Error("gRPC server leaked")
		}
	})
	return NewBerthControlClient(conn)
}
func inspectionContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)
	return metadata.NewOutgoingContext(ctx, metadata.Pairs("x-colony-key", "orbit-32"))
}
func TestHiddenInspectionStatusAndAuth(t *testing.T) {
	rng := rand.New(rand.NewSource(3202))
	readings := map[string]int64{}
	for i := 0; i < 30; i++ {
		readings[fmt.Sprintf("deck-%d", i)] = rng.Int63n(100000)
	}
	readings["a"] = 0
	readings[strings.Repeat("a", 16)] = -13
	var calls atomic.Int32
	client := inspectionClient(t, func(ctx context.Context, sector string) (int64, error) {
		calls.Add(1)
		if sector == "broken" {
			return 0, errors.New("secret tank details")
		}
		n, ok := readings[sector]
		if !ok {
			return 0, fmt.Errorf("lookup: %w", ErrSectorMissing)
		}
		return n, nil
	})
	ctx := inspectionContext(t)
	for sector, want := range readings {
		got, err := client.Inspect(ctx, &BerthQuery{Sector: sector})
		if err != nil || got.GetSector() != sector || got.GetPressure() != want {
			t.Fatalf("%q got %v %v", sector, got, err)
		}
	}
	for _, tc := range []struct {
		sector string
		code   codes.Code
	}{
		{"", codes.InvalidArgument}, {"A", codes.InvalidArgument}, {"a b", codes.InvalidArgument}, {"1x", codes.InvalidArgument}, {strings.Repeat("a", 17), codes.InvalidArgument}, {"unknown", codes.NotFound}, {"broken", codes.Internal},
	} {
		before := calls.Load()
		_, err := client.Inspect(ctx, &BerthQuery{Sector: tc.sector})
		if status.Code(err) != tc.code || strings.Contains(err.Error(), "secret") {
			t.Fatalf("%q status: %v", tc.sector, err)
		}
		if tc.code == codes.InvalidArgument && calls.Load() != before {
			t.Fatal("invalid sector reached store")
		}
	}
	for _, pairs := range [][]string{nil, {"x-colony-key", "wrong"}, {"x-colony-key", "orbit-32", "x-colony-key", "orbit-32"}} {
		denied := metadata.NewOutgoingContext(ctx, metadata.Pairs(pairs...))
		before := calls.Load()
		_, err := client.Inspect(denied, &BerthQuery{Sector: "a"})
		if status.Code(err) != codes.Unauthenticated || calls.Load() != before {
			t.Fatalf("auth: %v", err)
		}
	}
}
func TestHiddenInspectionConcurrentAndCanceled(t *testing.T) {
	entered := make(chan struct{}, 1)
	client := inspectionClient(t, func(ctx context.Context, sector string) (int64, error) {
		if sector == "wait" {
			entered <- struct{}{}
			<-ctx.Done()
			return 0, ctx.Err()
		}
		return int64(len(sector)), nil
	})
	var wg sync.WaitGroup
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sector := fmt.Sprintf("bay-%d", i)
			got, err := client.Inspect(inspectionContext(t), &BerthQuery{Sector: sector})
			if err != nil || got.GetPressure() != int64(len(sector)) {
				t.Errorf("parallel: %v %v", got, err)
			}
		}(i)
	}
	wg.Wait()
	ctx, cancel := context.WithCancel(inspectionContext(t))
	done := make(chan error, 1)
	go func() { _, err := client.Inspect(ctx, &BerthQuery{Sector: "wait"}); done <- err }()
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("RPC never entered")
	}
	cancel()
	select {
	case err := <-done:
		if status.Code(err) != codes.Canceled {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancellation stuck")
	}
	expired, stop := context.WithDeadline(inspectionContext(t), time.Unix(1, 0))
	defer stop()
	if _, err := client.Inspect(expired, &BerthQuery{Sector: "a"}); status.Code(err) != codes.DeadlineExceeded {
		t.Fatal(err)
	}
}
func TestHiddenBonusInspectionRecovery(t *testing.T) {
	client := inspectionClient(t, func(context.Context, string) (int64, error) { panic("private reactor key") })
	_, err := client.Inspect(inspectionContext(t), &BerthQuery{Sector: "reactor"})
	if status.Code(err) != codes.Internal || status.Convert(err).Message() != "inspection unavailable" {
		t.Fatalf("panic escaped: %v", err)
	}
}
`

const ch323Protocol0 = `syntax = "proto3";
package colony.drift;
option go_package = "orbital.lab/protocol;main";
message DriftStep { uint32 sequence = 1; int64 delta = 2; }
message DriftAck { uint32 sequence = 1; int64 total = 2; }
service DriftControl { rpc Calibrate(stream DriftStep) returns (stream DriftAck); }
`

const ch323Protocol1 = `// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.12
// 	protoc        v7.36.2
// source: drift.proto

package main

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DriftStep struct {
	state         protoimpl.MessageState ` + "`" + `protogen:"open.v1"` + "`" + `
	Sequence      uint32                 ` + "`" + `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"` + "`" + `
	Delta         int64                  ` + "`" + `protobuf:"varint,2,opt,name=delta,proto3" json:"delta,omitempty"` + "`" + `
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DriftStep) Reset() {
	*x = DriftStep{}
	mi := &file_drift_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DriftStep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DriftStep) ProtoMessage() {}

func (x *DriftStep) ProtoReflect() protoreflect.Message {
	mi := &file_drift_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DriftStep.ProtoReflect.Descriptor instead.
func (*DriftStep) Descriptor() ([]byte, []int) {
	return file_drift_proto_rawDescGZIP(), []int{0}
}

func (x *DriftStep) GetSequence() uint32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *DriftStep) GetDelta() int64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

type DriftAck struct {
	state         protoimpl.MessageState ` + "`" + `protogen:"open.v1"` + "`" + `
	Sequence      uint32                 ` + "`" + `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"` + "`" + `
	Total         int64                  ` + "`" + `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"` + "`" + `
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DriftAck) Reset() {
	*x = DriftAck{}
	mi := &file_drift_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DriftAck) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DriftAck) ProtoMessage() {}

func (x *DriftAck) ProtoReflect() protoreflect.Message {
	mi := &file_drift_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DriftAck.ProtoReflect.Descriptor instead.
func (*DriftAck) Descriptor() ([]byte, []int) {
	return file_drift_proto_rawDescGZIP(), []int{1}
}

func (x *DriftAck) GetSequence() uint32 {
	if x != nil {
		return x.Sequence
	}
	return 0
}

func (x *DriftAck) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

var File_drift_proto protoreflect.FileDescriptor

const file_drift_proto_rawDesc = "" +
	"\n" +
	"\vdrift.proto\x12\fcolony.drift\"=\n" +
	"\tDriftStep\x12\x1a\n" +
	"\bsequence\x18\x01 \x01(\rR\bsequence\x12\x14\n" +
	"\x05delta\x18\x02 \x01(\x03R\x05delta\"<\n" +
	"\bDriftAck\x12\x1a\n" +
	"\bsequence\x18\x01 \x01(\rR\bsequence\x12\x14\n" +
	"\x05total\x18\x02 \x01(\x03R\x05total2P\n" +
	"\fDriftControl\x12@\n" +
	"\tCalibrate\x12\x17.colony.drift.DriftStep\x1a\x16.colony.drift.DriftAck(\x010\x01B\x1bZ\x19orbital.lab/protocol;mainb\x06proto3"

var (
	file_drift_proto_rawDescOnce sync.Once
	file_drift_proto_rawDescData []byte
)

func file_drift_proto_rawDescGZIP() []byte {
	file_drift_proto_rawDescOnce.Do(func() {
		file_drift_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_drift_proto_rawDesc), len(file_drift_proto_rawDesc)))
	})
	return file_drift_proto_rawDescData
}

var file_drift_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_drift_proto_goTypes = []any{
	(*DriftStep)(nil), // 0: colony.drift.DriftStep
	(*DriftAck)(nil),  // 1: colony.drift.DriftAck
}
var file_drift_proto_depIdxs = []int32{
	0, // 0: colony.drift.DriftControl.Calibrate:input_type -> colony.drift.DriftStep
	1, // 1: colony.drift.DriftControl.Calibrate:output_type -> colony.drift.DriftAck
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_drift_proto_init() }
func file_drift_proto_init() {
	if File_drift_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_drift_proto_rawDesc), len(file_drift_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_drift_proto_goTypes,
		DependencyIndexes: file_drift_proto_depIdxs,
		MessageInfos:      file_drift_proto_msgTypes,
	}.Build()
	File_drift_proto = out.File
	file_drift_proto_goTypes = nil
	file_drift_proto_depIdxs = nil
}
`

const ch323Protocol2 = `// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v7.36.2
// source: drift.proto

package main

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DriftControl_Calibrate_FullMethodName = "/colony.drift.DriftControl/Calibrate"
)

// DriftControlClient is the client API for DriftControl service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DriftControlClient interface {
	Calibrate(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DriftStep, DriftAck], error)
}

type driftControlClient struct {
	cc grpc.ClientConnInterface
}

func NewDriftControlClient(cc grpc.ClientConnInterface) DriftControlClient {
	return &driftControlClient{cc}
}

func (c *driftControlClient) Calibrate(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[DriftStep, DriftAck], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DriftControl_ServiceDesc.Streams[0], DriftControl_Calibrate_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DriftStep, DriftAck]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DriftControl_CalibrateClient = grpc.BidiStreamingClient[DriftStep, DriftAck]

// DriftControlServer is the server API for DriftControl service.
// All implementations must embed UnimplementedDriftControlServer
// for forward compatibility.
type DriftControlServer interface {
	Calibrate(grpc.BidiStreamingServer[DriftStep, DriftAck]) error
	mustEmbedUnimplementedDriftControlServer()
}

// UnimplementedDriftControlServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDriftControlServer struct{}

func (UnimplementedDriftControlServer) Calibrate(grpc.BidiStreamingServer[DriftStep, DriftAck]) error {
	return status.Errorf(codes.Unimplemented, "method Calibrate not implemented")
}
func (UnimplementedDriftControlServer) mustEmbedUnimplementedDriftControlServer() {}
func (UnimplementedDriftControlServer) testEmbeddedByValue()                      {}

// UnsafeDriftControlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DriftControlServer will
// result in compilation errors.
type UnsafeDriftControlServer interface {
	mustEmbedUnimplementedDriftControlServer()
}

func RegisterDriftControlServer(s grpc.ServiceRegistrar, srv DriftControlServer) {
	// If the following call pancis, it indicates UnimplementedDriftControlServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DriftControl_ServiceDesc, srv)
}

func _DriftControl_Calibrate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DriftControlServer).Calibrate(&grpc.GenericServerStream[DriftStep, DriftAck]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DriftControl_CalibrateServer = grpc.BidiStreamingServer[DriftStep, DriftAck]

// DriftControl_ServiceDesc is the grpc.ServiceDesc for DriftControl service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DriftControl_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "colony.drift.DriftControl",
	HandlerType: (*DriftControlServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Calibrate",
			Handler:       _DriftControl_Calibrate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "drift.proto",
}
`

func init() {
	add(&Level{ID: "32-3-grpc-stream", Chapter: 32, Kind: "lab", Title: "gRPC streaming", Goal: "عاير انحراف المركبة بخطوات وإقرارات مستقلة لكل stream",
		Requires: []string{"sel:grpc.NewServer"}, Deps: []string{"google.golang.org/grpc", "google.golang.org/protobuf"}, Bridge: []int{76}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch323Starter}, {Name: "drift.proto", Starter: ch323Protocol0, ReadOnly: true}, {Name: "drift.pb.go", Starter: ch323Protocol1, ReadOnly: true}, {Name: "drift_grpc.pb.go", Starter: ch323Protocol2, ReadOnly: true}},
		Lesson: `# وصلة معايرة انحراف في الاتجاهين

## قبل ما تكتب
- الإجمالي المتراكم يعيش فين عشان طاقمين ميشاركوش نفس الحالة؟
- قفل إرسال العميل معناه إن كل الردود وصلت؟

## من unary للبث
الـ unary بيبادل طلب واحد ورد واحد. الـ bidirectional RPC بيبادل سلسلة مرتبة في الاتجاهين. كلمة ~stream~ في الناحيتين في proto بتعلن ده؛ الكود المتولّد بيدي streams بأنواع محددة بدل slices.
~~~proto
rpc Track(stream ProbePoint) returns (stream ProbePosition);
~~~
~drift.proto~ المقفول بيعرّف ~DriftStep~ بحقل ~Sequence~ نوعه uint32 و ~Delta~ نوعه int64، و ~DriftAck~ بحقل ~Sequence~ و ~Total~. الـ getters المتولّدة زي المهمة السابقة.
~~~go
step := &DriftStep{Sequence: 1, Delta: -9}
ack := &DriftAck{Sequence: step.GetSequence(), Total: step.GetDelta()}
fmt.Println(ack.GetTotal())
~~~

## Stream بأنواع واضحة
ضمّن ~UnimplementedDriftControlServer~ بالقيمة. نفّذ ~Calibrate(grpc.BidiStreamingServer[DriftStep, DriftAck]) error~. دي interface عامة متخصصة في النوعين دول. ~Recv~ بتقرا طلب واحد؛ ~Send~ بتكتب رد واحد؛ ~Context~ بتعدّي الإلغاء والمهلة.
~~~go
func (s probeService) Calibrate(link grpc.BidiStreamingServer[DriftStep, DriftAck]) error {
 q, err := link.Recv()
 if err != nil { return err }
 return link.Send(&DriftAck{Sequence: q.GetSequence(), Total: 4})
}
~~~
استخدم ~grpc.NewServer~ و ~RegisterDriftControlServer(srv, service)~ زي تسجيل unary. المتغيرات المحلية جوه ~Calibrate~ تخص الـ stream ده. ~DriftAdjust~ المحقونة بتاخد الإجمالي الحالي والتغيير الجاي؛ استخدم الإجمالي اللي رجعته بدل افتراض الجمع.
~~~go
next, err := adjust(link.Context(), total, q.GetDelta())
~~~

## EOF والفشل والضغط العكسي
~io.EOF~ من ~Recv~ معناه إن العميل قفل ناحية الإرسال. ارجع ~nil~ بعد إرسال الشغل المقبول كله؛ متحوّلش EOF لفشل. خطأ قراية تاني فشل حقيقي. فشل الكتابة يوقف الـ stream فوراً؛ متحاولش تاني ولا تبعت إقرارات بعده.
~~~go
q, err := link.Recv()
if errors.Is(err, io.EOF) { return nil }
if err != nil { return err }
~~~
حلقة قراية ثم تحقق ثم معايرة ثم إرسال بتعمل ضغط عكسي: مش بتقرا شغل زيادة بلا حدود. لسه النقل في الاتجاهين. التبادل خطوة بخطوة كفاية هنا؛ مش محتاج قنوات ولا goroutines إضافية جوه الـ stream. gRPC بيشغّل streams مختلفة بالتزامن.

خطأ التحقق ~codes.InvalidArgument~، وتجاوز الحصة ~codes.ResourceExhausted~. فشل المخزن ~codes.Internal~ عام. زي قبل كده افحص ~ctx.Err()~ قبل وبعد الشغل المحقون وحوّل بـ ~status.FromContextError~. الإقرارات القديمة تفضل صحيحة لو رسالة لاحقة فشلت.
~~~go
return status.Error(codes.ResourceExhausted, "step limit")
~~~

## العميل يقفل الإرسال وبعدها يكمّل قراية
~NewDriftControlClient(cc).Calibrate(ctx)~ بتفتح stream على نفس إعداد bufconn اللي اتعلمته. ~Send~ بتبعت الخطوات؛ ~CloseSend~ بتعلن نهاية الخطوات من غير تجاهل الردود. ~Recv~ تكمل لحد EOF. متعملش ~Send~ متزامن على نفس الـ stream ولا ~Recv~ متزامن؛ مرسل واحد وقارئ واحد يقدروا يشتغلوا مع بعض.
~~~go
link, err := NewDriftControlClient(cc).Calibrate(ctx)
err = link.Send(&DriftStep{Sequence: 1, Delta: 6})
reply, err := link.Recv()
err = link.CloseSend()
_, err = link.Recv() // EOF بعد آخر إقرار
~~~

## أخطاء شائعة
~~~go
type service struct { total int64 } // غلط: مشترك بين streams الأطقم
if err == io.EOF { return status.Error(codes.Internal, "failed") } // نهاية نظيفة بقت فشل
link.Send(reply); continue // غلط: تجاهلت فشل الإرسال
~~~

## المهمة
نفّذ ~NewDriftServer(adjust DriftAdjust) *grpc.Server~ وسجّل ~DriftControl~. دالة المعايرة المحقونة مش ~nil~ وآمنة بالتزامن. كل ~Calibrate~ يبدأ بإجمالي 0 وتسلسل متوقع 1. اقبل لحد 64 خطوة لكل stream. في الخطوات من 1 لحد 64 لازم نفس التسلسل المتوقع وdelta من -1000 لحد 1000 شاملة؛ غير كده انهي بـ ~InvalidArgument~ قبل المعايرة. بعد 64 خطوة مقبولة أي رسالة زيادة تنهي بـ ~ResourceExhausted~ قبل التحقق أو المعايرة.

لكل خطوة مقبولة نادِ ~adjust(stream.Context(), total, delta)~ مرة واحدة، واحفظ الإجمالي الراجع، وابعت ~DriftAck~ واحد بنفس التسلسل والإجمالي. افحص السياق قبل وبعد المعايرة. أي خطأ معايرة غير السياق ينهي بـ ~Internal~ ورسالة ~calibration unavailable~ بالظبط؛ متبعتش إقرار للخطوة الفاشلة ولا تسرّب نص خاص. خطأ السياق يتحوّل لكود gRPC. خطأ القراية غير EOF وكل فشل إرسال يوقفوا الشغل فوراً.

Stream فاضي أو قفل إرسال بعد خطوات مقبولة يدي EOF نظيف للعميل. مفيش معايرة ولا رد بعد رسالة مرفوضة؛ الردود القديمة تفضل. الـ streams المتزامنة حالتها منفصلة. متعدّلش البروتوكول ولا الملفات المتولّدة المقفولة. مفيش ~Report~ و ~main~ فاضية.
`,
		Hints: []string{"ابدأ بإجمالي وتسلسل محليين جوه Calibrate، مش جوه struct الخدمة.", "في كل لفة استقبل، ميّز EOF، افحص الحصة والتسلسل والحدود، وبعدها نادِ المعايرة.", "افحص السياق قبل وبعد المعايرة، وابعت إقرار واحد بالنتيجة؛ أول خطأ Send أو Recv يوقف الحلقة."},
		En: &LevelText{Title: "gRPC Streaming", Goal: "Calibrate vehicle drift with ordered steps and acknowledgements isolated per stream", Files: []LabFile{{Name: "main.go", Starter: ch323Starter}, {Name: "drift.proto", Starter: ch323Protocol0, ReadOnly: true}, {Name: "drift.pb.go", Starter: ch323Protocol1, ReadOnly: true}, {Name: "drift_grpc.pb.go", Starter: ch323Protocol2, ReadOnly: true}},
			Lesson: `# A Bidirectional Drift Calibration Link

## Before you write
- Where does the running total live so two crews cannot share it?
- Does a client half-close mean all responses have already arrived?

## From unary to streaming
A unary call exchanges one request and one reply. A bidirectional RPC exchanges ordered sequences in both directions. ~stream~ on both sides of the proto method declares this; the generated code supplies typed streams rather than slices.
~~~proto
rpc Track(stream ProbePoint) returns (stream ProbePosition);
~~~
Our read-only ~drift.proto~ declares ~DriftStep~ with ~Sequence~ (uint32) and ~Delta~ (int64), and ~DriftAck~ with ~Sequence~ and ~Total~. The generated getters work like the previous mission.
~~~go
step := &DriftStep{Sequence: 1, Delta: -9}
ack := &DriftAck{Sequence: step.GetSequence(), Total: step.GetDelta()}
fmt.Println(ack.GetTotal())
~~~

## A typed server stream
Embed ~UnimplementedDriftControlServer~ by value. Implement ~Calibrate(grpc.BidiStreamingServer[DriftStep, DriftAck]) error~. This is a generic interface generated for a stream of these exact message types. ~Recv~ reads one request; ~Send~ writes one reply; ~Context~ carries cancellation and deadlines.
~~~go
func (s probeService) Calibrate(link grpc.BidiStreamingServer[DriftStep, DriftAck]) error {
 q, err := link.Recv()
 if err != nil { return err }
 return link.Send(&DriftAck{Sequence: q.GetSequence(), Total: 4})
}
~~~
Use ~grpc.NewServer~ and ~RegisterDriftControlServer(srv, service)~ just like registration in the unary mission. Local variables inside ~Calibrate~ belong to that stream. The injected ~DriftAdjust~ gets the current total and next delta; trust its returned total rather than assuming addition.
~~~go
next, err := adjust(link.Context(), total, q.GetDelta())
~~~

## EOF, errors, and backpressure
~io.EOF~ from ~Recv~ means the client has half-closed its sending side. Return nil after all accepted work has been sent; do not turn EOF into an error status. A different receive failure is a real failure. A send failure stops the stream immediately; never retry or send later acknowledgements.
~~~go
q, err := link.Recv()
if errors.Is(err, io.EOF) { return nil }
if err != nil { return err }
~~~
A single receive/validate/adjust/send loop provides backpressure: it does not read unlimited work ahead. It still implements bidirectional transport. For this mission lockstep is sufficient; there is no need for queues or additional goroutines inside a stream. Multiple streams run concurrently through gRPC.

Validation errors use ~codes.InvalidArgument~; a stream quota uses ~codes.ResourceExhausted~. A store failure becomes generic ~codes.Internal~. As before, check ~ctx.Err()~ before and after injected work and translate with ~status.FromContextError~. Earlier acknowledgements remain valid when a later message fails.
~~~go
return status.Error(codes.ResourceExhausted, "step limit")
~~~

## The client half-closes, then drains
~NewDriftControlClient(cc).Calibrate(ctx)~ opens a stream over the same bufconn setup as before. ~Send~ sends steps; ~CloseSend~ announces no more steps without discarding replies. ~Recv~ continues until EOF. Do not call concurrent ~Send~ operations on one stream or concurrent ~Recv~ operations; one sender and one receiver may run concurrently.
~~~go
link, err := NewDriftControlClient(cc).Calibrate(ctx)
err = link.Send(&DriftStep{Sequence: 1, Delta: 6})
reply, err := link.Recv()
err = link.CloseSend()
_, err = link.Recv() // EOF after the final acknowledgement
~~~

## Common mistakes
~~~go
type service struct { total int64 } // Wrong: shared by all crews' streams
if err == io.EOF { return status.Error(codes.Internal, "failed") } // Clean end became failure
link.Send(reply); continue // Wrong: ignored send failure
~~~

## Your mission
Implement ~NewDriftServer(adjust DriftAdjust) *grpc.Server~ registering ~DriftControl~. The non-nil injected adjust function is concurrent-safe. Each ~Calibrate~ starts with total 0 and expected sequence 1. Accept up to 64 steps per stream. For steps 1..64 require the exact next sequence and delta in inclusive -1000..1000; otherwise terminate with ~InvalidArgument~ before adjustment. After 64 accepted steps, any further received message terminates with ~ResourceExhausted~ before validation or adjustment.

For each accepted step call ~adjust(stream.Context(), total, delta)~ once, store its returned total, then send exactly one ~DriftAck~ containing that sequence and total. Check context before and after adjustment. Any non-context adjustment error terminates with ~Internal~ and exactly ~calibration unavailable~; do not acknowledge the failed step or expose private text. Context errors map to their gRPC codes. Non-EOF receive failures and all send failures stop the stream without further work.

An empty stream or a half-close after accepted steps returns clean EOF to the client. No adjustment or reply follows a rejected message; prior replies remain. Simultaneous streams keep separate state. Do not edit the read-only protocol or generated files. No ~Report~; ~main~ stays empty.
`,
			Hints: []string{"Keep the running total and expected sequence local to Calibrate, outside the service struct.", "Receive each step, distinguish EOF, check quota, sequence, and bounds, then call the injected adjustment.", "Check context before and after adjustment, send one acknowledgement with its result, and stop on the first Send or Recv error."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch323Hidden},
	})
}

const ch323Starter = `package main

import (
	"context"
	"google.golang.org/grpc"
)

type DriftAdjust func(context.Context, int64, int64) (int64, error)

func NewDriftServer(adjust DriftAdjust) *grpc.Server { return grpc.NewServer() }
func main()                                          {}
`
const ch323Hidden = `package main

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"orbital/colony/labkit"
)

func driftClient(t *testing.T, adjust DriftAdjust) DriftControlClient {
	t.Helper()
	labkit.LeakCheck(t)
	ln := bufconn.Listen(1 << 20)
	srv := NewDriftServer(adjust)
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()
	cc, err := grpc.NewClient("passthrough:///drift", grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return ln.DialContext(ctx) }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cc.Close()
		srv.Stop()
		ln.Close()
		select {
		case <-served:
		case <-time.After(3 * time.Second):
			t.Error("server leaked")
		}
	})
	return NewDriftControlClient(cc)
}
func driftContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(cancel)
	return ctx
}
func TestHiddenDriftOrderedAndIsolated(t *testing.T) {
	client := driftClient(t, func(_ context.Context, total, delta int64) (int64, error) { return total + 2*delta, nil })
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(seed int64) {
			defer wg.Done()
			stream, err := client.Calibrate(driftContext(t))
			if err != nil {
				t.Error(err)
				return
			}
			rng := rand.New(rand.NewSource(seed))
			var total int64
			count := 64
			if seed == 0 {
				count = 0
			}
			for i := 1; i <= count; i++ {
				delta := int64(rng.Intn(2001) - 1000)
				if i == 1 {
					delta = -1000
				}
				if i == 2 {
					delta = 1000
				}
				if i == 3 {
					delta = 0
				}
				if err := stream.Send(&DriftStep{Sequence: uint32(i), Delta: delta}); err != nil {
					t.Error(err)
					return
				}
				total += 2 * delta
				ack, err := stream.Recv()
				if err != nil || ack.GetSequence() != uint32(i) || ack.GetTotal() != total {
					t.Errorf("ack %v %v want %d", ack, err, total)
					return
				}
			}
			if err := stream.CloseSend(); err != nil {
				t.Error(err)
				return
			}
			if _, err := stream.Recv(); !errors.Is(err, io.EOF) {
				t.Errorf("half-close: %v", err)
			}
		}(int64(worker))
	}
	wg.Wait()
}
func TestHiddenDriftRejectsAndStops(t *testing.T) {
	var calls atomic.Int32
	client := driftClient(t, func(_ context.Context, total, delta int64) (int64, error) { calls.Add(1); return total + delta, nil })
	for _, q := range []*DriftStep{{Sequence: 0}, {Sequence: 2}, {Sequence: 1, Delta: -1001}, {Sequence: 1, Delta: 1001}} {
		before := calls.Load()
		stream, err := client.Calibrate(driftContext(t))
		if err != nil {
			t.Fatal(err)
		}
		if err := stream.Send(q); err != nil {
			t.Fatal(err)
		}
		_, err = stream.Recv()
		if status.Code(err) != codes.InvalidArgument || calls.Load() != before {
			t.Fatalf("invalid step: %v", err)
		}
	}
	stream, err := client.Calibrate(driftContext(t))
	if err != nil {
		t.Fatal(err)
	}
	for i := uint32(1); i <= 64; i++ {
		if err := stream.Send(&DriftStep{Sequence: i}); err != nil {
			t.Fatal(err)
		}
		if _, err := stream.Recv(); err != nil {
			t.Fatal(err)
		}
	}
	before := calls.Load()
	if err := stream.Send(&DriftStep{Sequence: 65}); err != nil {
		t.Fatal(err)
	}
	if _, err := stream.Recv(); status.Code(err) != codes.ResourceExhausted || calls.Load() != before {
		t.Fatalf("limit: %v", err)
	}
	duplicate, err := client.Calibrate(driftContext(t))
	if err != nil {
		t.Fatal(err)
	}
	duplicate.Send(&DriftStep{Sequence: 1})
	if _, err := duplicate.Recv(); err != nil {
		t.Fatal(err)
	}
	duplicate.Send(&DriftStep{Sequence: 1})
	if _, err := duplicate.Recv(); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("duplicate: %v", err)
	}
}
func TestHiddenDriftFailureAndCancellation(t *testing.T) {
	entered := make(chan struct{}, 1)
	exited := make(chan struct{}, 1)
	client := driftClient(t, func(ctx context.Context, total, delta int64) (int64, error) {
		if delta == 7 {
			entered <- struct{}{}
			<-ctx.Done()
			exited <- struct{}{}
			return 0, ctx.Err()
		}
		return 0, errors.New("secret thruster serial")
	})
	stream, err := client.Calibrate(driftContext(t))
	if err != nil {
		t.Fatal(err)
	}
	stream.Send(&DriftStep{Sequence: 1, Delta: 2})
	_, err = stream.Recv()
	if status.Code(err) != codes.Internal || status.Convert(err).Message() != "calibration unavailable" || strings.Contains(err.Error(), "secret") {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(driftContext(t))
	defer cancel()
	blocked, err := client.Calibrate(ctx)
	if err != nil {
		t.Fatal(err)
	}
	blocked.Send(&DriftStep{Sequence: 1, Delta: 7})
	select {
	case <-entered:
	case <-time.After(3 * time.Second):
		t.Fatal("adjust never entered")
	}
	cancel()
	if _, err := blocked.Recv(); status.Code(err) != codes.Canceled {
		t.Fatal(err)
	}
	select {
	case <-exited:
	case <-time.After(3 * time.Second):
		t.Fatal("adjust did not receive cancellation")
	}
	expired, stop := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer stop()
	s, err := client.Calibrate(expired)
	if err == nil {
		_, err = s.Recv()
	}
	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatal(err)
	}
}
`

func init() {
	add(&Level{ID: "32-4-tls", Chapter: 32, Kind: "lab", Title: "TLS و mTLS", Goal: "جهّز شهادات بوابة الهواء وثقة الطرفين بساعة محقونة",
		Requires: []string{"sel:tls.Config"}, Deps: []string{}, Bridge: []int{84}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch324Starter}},
		Lesson: `# شهادات بوابة الهواء

## قبل ما تكتب
- تشفير الاتصال يثبت إن السيرفر هو المحطة اللي قصدتها؟
- مين بيتحقق من هوية الطاقم لما البوابة تطلب mTLS؟

## هويات X.509 والثقة
TLS بيشفّر النقل وبيثبت الهوية. شهادة X.509 بتربط مفتاح عام بأسماء وتواريخ صلاحية واستخدامات مسموحة، بتوقيع الجهة المصدرة. شهادة CA هي أصل الثقة. شهادة السيرفر محتاجة اسم DNS مطابق في SAN واستخدام server authentication؛ العميل محتاج client authentication. التوقيع لوحده مش كفاية لإثبات الهوية المقصودة.

~tls.Certificate~ فيها سلسلة الشهادات والمفتاح الخاص. التطبيقات الحقيقية بتقرا ملفات PEM بـ ~tls.LoadX509KeyPair~؛ هنا الشهادة محقونة في الذاكرة عشان منحتاجش ملفات الجهاز. متطبعش المفاتيح الخاصة ولا تحفظها كنص.
~~~go
cert, err := tls.LoadX509KeyPair("station.pem", "station-key.pem")
~~~
~x509.NewCertPool~ بتعمل مجموعة ثقة؛ ~AddCert~ بتضيف CA. ~Clone~ بتنسخ المجموعة عشان تغيير الإعدادات ميغيّرش ملك المستدعي.
~~~go
roots := x509.NewCertPool()
roots.AddCert(ca.Certificate)
privateRoots := roots.Clone()
~~~

## إعداد الطرفين
~tls.Config~ بتتجهز قبل الاستخدام وبعدها تعاملها كقيمة ثابتة. ~Certificates~ عند السيرفر فيها السلسلة والمفتاح. ~RootCAs~ عند العميل بتختار الجهات الموثوقة؛ ~ServerName~ بتفحص SAN المطلوب وبتبعت SNI. اسم DNS مش محتاج lookup لما الاتصال محقون جاهز.
~~~go
server := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}
client := &tls.Config{RootCAs: roots, ServerName: "probe.colony", MinVersion: tls.VersionTLS13}
~~~
~MinVersion~ بترفض البروتوكولات الأقدم؛ ~tls.VersionTLS13~ هي TLS 1.3. ~MaxVersion~ بصفر بتسمح بنسخ أحدث مدعومة. متشغّلش ~InsecureSkipVerify~ لأنها بتلغي التحقق اللي بنتعلمه. التحقق بيرفض جهة مجهولة أو SAN غلط أو شهادة منتهية أو لسه مبدأتش أو استخدام موسّع غلط.

~Time~ دالة بترجع ساعة التحقق. الطرفين لازم يستخدموا الساعة المحقونة عشان التست يختار لحظة حتمية للصلاحية من غير نوم.
~~~go
config.Time = func() time.Time { return chosenInstant }
~~~

## TLS متبادل
TLS العادي بيثبت هوية السيرفر بس. في mTLS السيرفر بيطلب شهادة العميل وبيتحقق منها. اضبط ~ClientAuth~ على ~tls.RequireAndVerifyClientCert~ و ~ClientCAs~ على مجموعة ثقة العملاء. العميل بيعرض شهادته في ~Certificates~ بتاعته. ~tls.NoClientCert~ الوضع الافتراضي للسيرفر ومبيطلبش شهادة عميل.
~~~go
server.ClientAuth = tls.RequireAndVerifyClientCert
server.ClientCAs = crewRoots
client.Certificates = []tls.Certificate{crewCert}
~~~
شهادة عميل غايبة أو أجنبية أو منتهية أو للسيرفر بس لازم تفشل. متستبدلش التحقق بـ callback بيقبل أي حاجة.

## مصافحة من غير socket
~tls.Server~ و ~tls.Client~ بيلفوا ~net.Conn~ محقونة. ~HandshakeContext~ بتفاوض وبتتحقق؛ شغّل الطرفين بالتزامن على ~net.Pipe~. اقفل الاتصالات الخام عشان تفك الطرف الفاشل واستنى الاتنين. في TLS 1.3 مصافحة العميل ممكن تخلص قبل رفض السيرفر النهائي لشهادة العميل؛ راقب نتيجة السيرفر كمان.
~~~go
left, right := net.Pipe()
serverConn := tls.Server(left, server)
clientConn := tls.Client(right, client)
go serverConn.HandshakeContext(ctx)
err := clientConn.HandshakeContext(ctx)
state := clientConn.ConnectionState()
fmt.Println(state.Version, state.VerifiedChains)
~~~
~labkit.TestCA~ بتعمل جهة شهادات محلية للتست. ~labkit.IssueCert~ بتاخد ~labkit.CertOptions~ فيها ~DNSNames~ و ~NotBefore~ و ~NotAfter~ و ~Client~؛ بتقدر بالفعل تطلع شهادات منتهية وباستخدام غلط. التست بيختار لحظة من فترة صلاحية CA ويحقنها في كودك. مفيش اتصال بجهة عامة.
~~~go
ca, err := labkit.TestCA()
cert, err := labkit.IssueCert(ca, labkit.CertOptions{
 DNSNames: []string{"probe.colony"}, NotBefore: instant.Add(-time.Hour), NotAfter: instant.Add(time.Hour),
})
~~~

## أخطاء شائعة
~~~go
client.InsecureSkipVerify = true // غلط: تشفير من غير إثبات هوية السيرفر
server.ClientAuth = tls.RequireAnyClientCert // غلط: بيقبل عميل غير موثوق
~~~

## المهمة
نفّذ ~ColonyTLS(opts ColonyTLSOptions) (*tls.Config, *tls.Config, error)~ ترجع إعداد السيرفر وبعده العميل. ارفض سلسلة شهادة سيرفر فاضية، أو مفتاح خاص ~nil~، أو ~Trust~ بـ ~nil~، أو ~ServerName~ فاضي، أو ~Now~ بـ ~nil~ بإعدادين ~nil~ وخطأ. غير كده رجّع إعدادات جديدة مستقلة ومن غير خطأ. كائنات الشهادة والمفاتيح سليمة وثابتة؛ صلاحية الشهادة تتفحص وقت المصافحة مش البناء.

الطرفين يستخدموا ~MinVersion: tls.VersionTLS13~ و ~Now~ المحقونة في ~Time~. السيرفر يعرض ~Server~. العميل يتحقق بنسخة ~Trust~ وبنفس ~ServerName~ والتحقق القياسي شغّال. ~Client~ غير ~nil~ تتعرض عند العميل. ~ClientTrust~ غير ~nil~ تشغّل ~RequireAndVerifyClientCert~ بنسخة المجموعة؛ ~nil~ تسيب ~NoClientCert~. غياب شهادة عميل في وضع الإلزام يفشل. متعدّلش المدخلات ولا تشارك إعدادات قابلة للتغيير بين النداءات. متعملش callbacks لتجاوز التحقق. التستات بتعمل مصافحات TLS حقيقية بشهادات منتهية ومستقبلية واسم وجهة واستخدام غلط. مفيش شبكة ولا ملفات ولا ~Report~ في حلك و ~main~ فاضية.
`,
		Hints: []string{"افحص المدخلات قبل بناء أي إعداد؛ اعمل tls.Config جديد لكل طرف.", "حط الشهادة عند السيرفر وRootCAs وServerName عند العميل، وMinVersion وTime عند الاتنين.", "ClientTrust غير nil معناها RequireAndVerifyClientCert؛ انسخ مجموعات الثقة بـ Clone وسيب التحقق القياسي شغال."},
		En: &LevelText{Title: "TLS and mTLS", Goal: "Configure airlock certificates and peer trust with an injected verification clock", Files: []LabFile{{Name: "main.go", Starter: ch324Starter}},
			Lesson: `# Certificates for the Airlock

## Before you write
- Does encrypting a connection prove the server is the station you intended?
- Who verifies the crew's identity when the airlock requires mutual TLS?

## X.509 identities and trust
TLS encrypts and authenticates a transport. An X.509 certificate binds a public key to names, validity dates, and allowed uses, signed by an issuer. A CA certificate is a trust anchor. A server certificate needs a matching DNS subject alternative name (SAN) and server authentication usage; client authentication needs client usage. A signature alone does not prove the intended identity.

~tls.Certificate~ holds the certificate chain and private key. Real applications load PEM files with ~tls.LoadX509KeyPair~; this lab receives a certificate in memory so no machine files are required. Do not print or serialize private keys.
~~~go
cert, err := tls.LoadX509KeyPair("station.pem", "station-key.pem")
~~~
~x509.NewCertPool~ creates a trust pool; ~AddCert~ adds a CA. ~Clone~ copies the pool so configuration changes do not change caller-owned settings.
~~~go
roots := x509.NewCertPool()
roots.AddCert(ca.Certificate)
privateRoots := roots.Clone()
~~~

## Configure both peers
~tls.Config~ is configured before use and then treated as immutable. The server's ~Certificates~ contains its chain and key. The client's ~RootCAs~ chooses trusted issuers; ~ServerName~ checks the expected SAN and also supplies SNI. A DNS name needs no DNS lookup when a connection is already injected.
~~~go
server := &tls.Config{Certificates: []tls.Certificate{cert}, MinVersion: tls.VersionTLS13}
client := &tls.Config{RootCAs: roots, ServerName: "probe.colony", MinVersion: tls.VersionTLS13}
~~~
~MinVersion~ rejects older protocols; ~tls.VersionTLS13~ is TLS 1.3. Leaving ~MaxVersion~ at zero permits newer supported versions. Never set ~InsecureSkipVerify~: it disables the very verification being taught. Verification rejects an unknown issuer, wrong SAN, expired or not-yet-valid certificate, and incorrect extended key usage.

~Time~ is a function returning the verification clock. Both peers must use the injected clock, so expiry tests can choose a deterministic instant without sleeping.
~~~go
config.Time = func() time.Time { return chosenInstant }
~~~

## Mutual TLS
Ordinary TLS authenticates only the server. In mTLS the server also requires and verifies a client certificate. Set ~ClientAuth~ to ~tls.RequireAndVerifyClientCert~ and ~ClientCAs~ to the client trust pool. The client presents its certificate through its own ~Certificates~. ~tls.NoClientCert~ is the default server mode and does not request client identity.
~~~go
server.ClientAuth = tls.RequireAndVerifyClientCert
server.ClientCAs = crewRoots
client.Certificates = []tls.Certificate{crewCert}
~~~
A missing, foreign, expired, or server-only client certificate must fail. Do not replace verification with a callback that accepts everything.

## Handshake without a socket
~tls.Server~ and ~tls.Client~ wrap an injected ~net.Conn~. ~HandshakeContext~ negotiates and verifies both ends; run peers concurrently on ~net.Pipe~. Close the raw connections to unblock failed peers, and join both goroutines. TLS 1.3 client handshake completion can precede the server's final rejection of a client certificate; observe the server outcome too.
~~~go
left, right := net.Pipe()
serverConn := tls.Server(left, server)
clientConn := tls.Client(right, client)
go serverConn.HandshakeContext(ctx)
err := clientConn.HandshakeContext(ctx)
state := clientConn.ConnectionState()
fmt.Println(state.Version, state.VerifiedChains)
~~~
~labkit.TestCA~ creates a local fixture issuer. ~labkit.IssueCert~ takes ~labkit.CertOptions~ with ~DNSNames~, ~NotBefore~, ~NotAfter~, and ~Client~; it can already issue expired and wrong-purpose fixtures. Tests derive their chosen instant from the CA validity interval, then inject it into your code. No public authority is contacted.
~~~go
ca, err := labkit.TestCA()
cert, err := labkit.IssueCert(ca, labkit.CertOptions{
 DNSNames: []string{"probe.colony"}, NotBefore: instant.Add(-time.Hour), NotAfter: instant.Add(time.Hour),
})
~~~

## Common mistakes
~~~go
client.InsecureSkipVerify = true // Wrong: encryption without verified server identity
server.ClientAuth = tls.RequireAnyClientCert // Wrong: accepts an untrusted client
~~~

## Your mission
Implement ~ColonyTLS(opts ColonyTLSOptions) (*tls.Config, *tls.Config, error)~, returning server config then client config. Reject an empty server certificate chain, nil server private key, nil ~Trust~, empty ~ServerName~, or nil ~Now~ with two nil configurations and an error. Otherwise return fresh independent configurations and no error. Supplied certificate objects and keys are immutable and well-formed; certificate validity is checked during the handshake, not construction.

Both peers use ~MinVersion: tls.VersionTLS13~ and the injected ~Now~ through ~Time~. The server presents ~Server~. The client verifies using a clone of ~Trust~ and exact ~ServerName~, with standard verification enabled. Non-nil ~Client~ is presented by the client. Non-nil ~ClientTrust~ enables ~RequireAndVerifyClientCert~ with a clone of that pool; nil keeps ~NoClientCert~. A missing client certificate in required mode must fail. Do not mutate input options or share mutable config settings across calls. Do not set custom verification bypasses. The tests use real TLS handshakes, including expired, future, wrong-name, foreign, and wrong-purpose certificates. No network, files, or ~Report~ in your solution; ~main~ is empty.
`,
			Hints: []string{"Validate inputs before constructing either configuration; create a new tls.Config for each peer.", "Set the server certificate, client RootCAs and ServerName, and both peers' MinVersion and Time.", "A non-nil ClientTrust selects RequireAndVerifyClientCert; clone trust pools and retain standard verification."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch324Hidden},
	})
}

const ch324Starter = `package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"time"
)

type ColonyTLSOptions struct {
	Server      tls.Certificate
	Trust       *x509.CertPool
	ClientTrust *x509.CertPool
	Client      *tls.Certificate
	ServerName  string
	Now         func() time.Time
}

func ColonyTLS(opts ColonyTLSOptions) (*tls.Config, *tls.Config, error) {
	return nil, nil, errors.New("TODO")
}
func main() {}
`
const ch324Hidden = `package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"

	"orbital/colony/labkit"
)

func tlsCert(t *testing.T, ca *labkit.CA, opts labkit.CertOptions) tls.Certificate {
	t.Helper()
	cert, err := labkit.IssueCert(ca, opts)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}
func tlsPool(ca *labkit.CA) *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AddCert(ca.Certificate)
	return pool
}
func tlsOptions(t *testing.T, name string) (ColonyTLSOptions, *labkit.CA, time.Time) {
	t.Helper()
	ca, err := labkit.TestCA()
	if err != nil {
		t.Fatal(err)
	}
	now := ca.Certificate.NotBefore.Add(2 * time.Hour)
	cert := tlsCert(t, ca, labkit.CertOptions{DNSNames: []string{name}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour)})
	return ColonyTLSOptions{Server: cert, Trust: tlsPool(ca), ServerName: name, Now: func() time.Time { return now }}, ca, now
}

type tlsOutcome struct {
	state tls.ConnectionState
	err   error
}

func tlsHandshake(t *testing.T, server, client *tls.Config) (tlsOutcome, tlsOutcome) {
	t.Helper()
	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	sc := tls.Server(left, server)
	cc := tls.Client(right, client)
	sd := make(chan tlsOutcome, 1)
	cd := make(chan tlsOutcome, 1)
	drained := make(chan struct{})
	go func() { err := sc.HandshakeContext(ctx); sd <- tlsOutcome{sc.ConnectionState(), err} }()
	go func() {
		defer close(drained)
		err := cc.HandshakeContext(ctx)
		cd <- tlsOutcome{cc.ConnectionState(), err}
		if err == nil {
			cc.Read(make([]byte, 1))
		}
	}()
	var s, c tlsOutcome
	select {
	case s = <-sd:
		if s.err != nil {
			left.Close()
			right.Close()
		}
		c = <-cd
	case c = <-cd:
		if c.err != nil {
			left.Close()
			right.Close()
		}
		s = <-sd
	case <-ctx.Done():
		t.Fatal("TLS handshake stuck")
	}
	left.Close()
	right.Close()
	<-drained
	return s, c
}
func TestHiddenTLSIdentityAndVersions(t *testing.T) {
	labkit.LeakCheck(t)
	rng := rand.New(rand.NewSource(3204))
	for i := 0; i < 8; i++ {
		name := fmt.Sprintf("dome-%d.colony", rng.Intn(100000))
		opts, _, _ := tlsOptions(t, name)
		server, client, err := ColonyTLS(opts)
		if err != nil {
			t.Fatal(err)
		}
		s, c := tlsHandshake(t, server, client)
		if s.err != nil || c.err != nil || c.state.Version != tls.VersionTLS13 || !c.state.HandshakeComplete || len(c.state.VerifiedChains) == 0 {
			t.Fatalf("trusted TLS: %v %v", s.err, c.err)
		}
		if server.MinVersion != tls.VersionTLS13 || client.MinVersion != tls.VersionTLS13 || client.InsecureSkipVerify {
			t.Fatal("unsafe TLS configuration")
		}
		if server.ClientAuth != tls.NoClientCert {
			t.Fatal("default unexpectedly requests client cert")
		}
		old := client.Clone()
		old.MaxVersion = tls.VersionTLS12
		old.MinVersion = tls.VersionTLS12
		s, c = tlsHandshake(t, server, old)
		if s.err == nil && c.err == nil {
			t.Fatal("TLS 1.2 accepted")
		}
	}
}
func TestHiddenTLSRejectsBadServerCertificates(t *testing.T) {
	for _, scenario := range []string{"name", "untrusted", "expired", "future", "wrong-usage"} {
		t.Run(scenario, func(t *testing.T) {
			labkit.LeakCheck(t)
			opts, ca, now := tlsOptions(t, "dome.colony")
			switch scenario {
			case "name":
				opts.ServerName = "other.colony"
			case "untrusted":
				other, err := labkit.TestCA()
				if err != nil {
					t.Fatal(err)
				}
				opts.Trust = tlsPool(other)
			case "expired":
				opts.Server = tlsCert(t, ca, labkit.CertOptions{DNSNames: []string{"dome.colony"}, NotBefore: now.Add(-2 * time.Hour), NotAfter: now.Add(-time.Minute)})
			case "future":
				opts.Server = tlsCert(t, ca, labkit.CertOptions{DNSNames: []string{"dome.colony"}, NotBefore: now.Add(time.Minute), NotAfter: now.Add(time.Hour)})
			case "wrong-usage":
				opts.Server = tlsCert(t, ca, labkit.CertOptions{DNSNames: []string{"dome.colony"}, Client: true, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour)})
			}
			server, client, err := ColonyTLS(opts)
			if err != nil {
				t.Fatal(err)
			}
			_, c := tlsHandshake(t, server, client)
			if c.err == nil {
				t.Fatal("bad server accepted")
			}
		})
	}
}
func TestHiddenTLSMutualIdentity(t *testing.T) {
	for _, scenario := range []string{"valid", "missing", "foreign", "expired", "server-usage"} {
		t.Run(scenario, func(t *testing.T) {
			labkit.LeakCheck(t)
			opts, ca, now := tlsOptions(t, "airlock.colony")
			opts.ClientTrust = tlsPool(ca)
			clientOpts := labkit.CertOptions{Client: true, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour)}
			issuer := ca
			switch scenario {
			case "foreign":
				var err error
				issuer, err = labkit.TestCA()
				if err != nil {
					t.Fatal(err)
				}
			case "expired":
				clientOpts.NotAfter = now.Add(-time.Minute)
			case "server-usage":
				clientOpts.Client = false
			}
			cert := tlsCert(t, issuer, clientOpts)
			if scenario != "missing" {
				opts.Client = &cert
			}
			server, client, err := ColonyTLS(opts)
			if err != nil {
				t.Fatal(err)
			}
			if server.ClientAuth != tls.RequireAndVerifyClientCert {
				t.Fatal("mTLS verification missing")
			}
			s, c := tlsHandshake(t, server, client)
			if scenario == "valid" {
				if s.err != nil || c.err != nil || len(s.state.VerifiedChains) == 0 {
					t.Fatalf("mTLS: %v %v", s.err, c.err)
				}
			} else if s.err == nil {
				t.Fatal("bad client accepted")
			}
		})
	}
}
func TestHiddenTLSOptionsAndIsolation(t *testing.T) {
	opts, _, _ := tlsOptions(t, "ring.colony")
	for _, change := range []func(*ColonyTLSOptions){
		func(o *ColonyTLSOptions) { o.Server = tls.Certificate{} }, func(o *ColonyTLSOptions) { o.Server.PrivateKey = nil }, func(o *ColonyTLSOptions) { o.Trust = nil }, func(o *ColonyTLSOptions) { o.ServerName = "" }, func(o *ColonyTLSOptions) { o.Now = nil },
	} {
		invalid := opts
		change(&invalid)
		s, c, err := ColonyTLS(invalid)
		if err == nil || s != nil || c != nil {
			t.Fatal("invalid options accepted")
		}
	}
	s1, c1, err := ColonyTLS(opts)
	if err != nil {
		t.Fatal(err)
	}
	s2, c2, err := ColonyTLS(opts)
	if err != nil {
		t.Fatal(err)
	}
	s1.MinVersion = tls.VersionTLS10
	c1.ServerName = "changed"
	if s1 == s2 || c1 == c2 || s2.MinVersion != tls.VersionTLS13 || c2.ServerName != opts.ServerName || c1.RootCAs == opts.Trust {
		t.Fatal("configurations share mutable settings")
	}
}
`

func init() {
	add(&Level{ID: "32-5-http2", Chapter: 32, Kind: "lab", Title: "HTTP/2", Goal: "خد قراءات متوازية على اتصال TLS واحد بتفاوض h2",
		Requires: []string{"sel:http.Server"}, Deps: []string{}, Bridge: []int{77}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch325Starter}},
		Lesson: `# بوابة المستعمرة بطلبات متوازية

## قبل ما تكتب
- رد حساس واقف يقدر يمنع رد تاني مستقل على نفس الاتصال؟
- التست يثبت مشاركة الاتصال إزاي من غير قياس وقت؟

## Streams على اتصال واحد
HTTP/1.1 غالباً محتاج اتصالات منفصلة لردود متزامنة مستقلة. HTTP/2 بيقسّم الاتصال لـ streams مرقّمة بإطارات متداخلة؛ طلب واقف مش لازم يعطّل أخوه. التطبيق لسه شايف ~http.Request~ و ~http.ResponseWriter~ عاديين؛ إطارات النقل مسؤولية المكتبة القياسية.

المهمة دي من غير Server Push. ابنِ handler آمن بالتزامن من ~SectorRead~ محقونة ورجّع ~http.Server~. متحطش قفل عام حوالين قراية حساس ممكن تقف؛ الدالة المحقونة بالفعل بتدعم التزامن.
~~~go
srv := &http.Server{Handler: h, ReadHeaderTimeout: 2*time.Second}
~~~

## ALPN بيختار بروتوكول التطبيق
TLS بيثبت هوية الطرفين؛ ALPN بيتفاوض على اللغة بعد المصافحة. ~tls.Config.NextProtos~ قايمة بروتوكولات التطبيق بالترتيب. ~h2~ معناها HTTP/2؛ و ~http/1.1~ بتسمح برجوع لعميل أقدم. لازم الطرفين يعرضوا اختيار مشترك.
~~~go
secure := config.Clone()
secure.NextProtos = []string{"h2", "http/1.1"}
secure.MinVersion = tls.VersionTLS13
~~~
انسخ إعداد المستدعي قبل التعديل. ~tls.Config.Clone~ نسخة سطحية؛ استبدل slice الـ ~NextProtos~ بدل الكتابة جوه مصفوفتها الأصلية. كائنات الشهادات ثابتة زي مهمة TLS.

~srv.ServeTLS(ln, "", "")~ بتستخدم الشهادات الموجودة في ~TLSConfig~ وبتجهز دعم HTTP/2 القياسي. متحطش map فاضية في ~TLSNextProto~ لأنها بتعطّل HTTP/2 التلقائي. التست يملك المستمع وعمر السيرفر؛ المصنع بتاعك ميفتحش بورت.
~~~go
srv := &http.Server{Handler: h, TLSConfig: secure}
err := srv.ServeTLS(listener, "", "")
~~~

## شغّل HTTP/2 عند العميل
~http.Transport~ مخصصة بإعدادات TLS مخصصة لازم تختار ~ForceAttemptHTTP2: true~ بوضوح. ~MaxConnsPerHost: 1~ بتحد الاتصالات مش streams الـ HTTP/2. ثق بالشهادة بالطريقة العادية من ~TLSClientConfig~؛ متتجاوزش التحقق عشان HTTP/2.
~~~go
transport := &http.Transport{
 TLSClientConfig: clientTLS, ForceAttemptHTTP2: true, MaxConnsPerHost: 1,
}
client := &http.Client{Transport: transport}
defer transport.CloseIdleConnections()
~~~
~resp.ProtoMajor~ بتقول نسخة HTTP. ~resp.TLS.NegotiatedProtocol~ بتقول نتيجة ALPN. قايمة البروتوكولات لوحدها مش دليل إن HTTP/2 اشتغل فعلاً.
~~~go
fmt.Println(resp.ProtoMajor, resp.TLS.NegotiatedProtocol) // 2 h2
~~~

## اثبت التزامن بحواجز
~httptrace.WithClientTrace~ بتربط مراقب بالعميل. ~ClientTrace.GotConn~ بتستقبل ~GotConnInfo~ اللي ~Conn~ فيها كائن الاتصال الحقيقي. قارن هوية الكائن بين الطلبات؛ ~Reused~ لوحدها مش دليل إنه نفس الاتصال. التست يجهز الاتصال ويوقف حساس على حاجز channel ويطلب من الباقي يخلصوا على نفس الاتصال قبل فك الحاجز.
~~~go
trace := &httptrace.ClientTrace{GotConn: func(info httptrace.GotConnInfo) {
 observed <- info.Conn
}}
ctx = httptrace.WithClientTrace(ctx, trace)
~~~
~httptest.NewUnstartedServer~ بتسمح بتركيب إعداد السيرفر وشهادات TLS قبل التشغيل. ~EnableHTTP2~ بتسمح للمرفق بـ HTTP/2؛ و ~StartTLS~ بتشغّله على بورت loopback عشوائي. اقفل أجسام الردود والنقل الخامل والسيرفر، وفك الـ handlers الواقفة قبل القفل.
~~~go
ts := httptest.NewUnstartedServer(handler)
ts.EnableHTTP2 = true
ts.StartTLS()
defer ts.Close()
~~~

## أخطاء شائعة
~~~go
config.NextProtos[0] = "h2" // غلط: بتغيّر slice ملك المستدعي
srv.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){} // بتعطّل h2
// قفل مشترك حوالين انتظار الحساس بيمنع streams تانية.
~~~

## المهمة
نفّذ ~NewMultiplexGate(config *tls.Config, read SectorRead) (*http.Server, error)~. ارفض config بـ ~nil~ أو من غير شهادات أو read بـ ~nil~ بسيرفر ~nil~ وخطأ. باقي حقول TLS والشهادات سليمة وثابتة. انسخ config واضبط أقل نسخة TLS 1.3 و ~NextProtos~ بالظبط ~["h2", "http/1.1"]~ ورجّع ~http.Server~ جديد بالإعداد ده و ~ReadHeaderTimeout~ ثانيتين و ~IdleTimeout~ 30 ثانية و ~MaxHeaderBytes~ 16384. حافظ على ساعة التحقق والشهادات المحقونة. متفتحش مستمع ولا تعطّل HTTP/2 التلقائي.

ركّب handler لـ ~GET /sectors/NAME~. الاسم من 1 لحد 16 بايت ASCII حروف صغيرة أو أرقام أو شرطات. أي مسار تاني، منه أجزاء زيادة أو حروف كبيرة، يرجع 404 قبل القراية. طريقة تانية على مسار صحيح ترجع 405 و ~Allow: GET~ من غير قراية. تجاهل query. نادِ ~read(r.Context(), NAME)~ مرة؛ النجاح 200 بنفس البايتات و ~Content-Type: text/plain; charset=utf-8~ و ~Cache-Control: no-store~. أي خطأ قراية يرجع 503 بنص عام من غير محتوى الحساس ولا تفاصيل الخطأ الخاصة. جسم الخطأ مش محدد غير كده.

الطلبات المتزامنة مستقلة؛ قراية واقفة متعطّلش إخواتها على نفس اتصال HTTP/2. التستات بتطلب ALPN حقيقي ~h2~ ونسخة HTTP برقم 2 ونفس هوية الاتصال ورجوع HTTP/1.1 ورفض الشهادات غير الموثوقة. مفيش Server Push ولا ~Report~ و ~main~ فاضية.
`,
		Hints: []string{"انسخ tls.Config واستبدل NextProtos بقائمة جديدة، واضبط MinVersion ومهل http.Server.", "افصل تحقق المسار والطريقة عن قراية الحساس؛ مرّر r.Context من غير قفل مشترك.", "سيب TLSNextProto كما هي عشان ServeTLS تجهّز HTTP/2؛ التست بيستخدم ForceAttemptHTTP2 وهوية GotConn بدل الوقت."},
		En: &LevelText{Title: "HTTP/2", Goal: "Read sensors concurrently over one TLS connection negotiated as h2", Files: []LabFile{{Name: "main.go", Starter: ch325Starter}},
			Lesson: `# Multiplex the Colony Gateway

## Before you write
- Can a blocked sensor response prevent an unrelated response on the same connection?
- How can a test prove connection sharing without measuring elapsed time?

## HTTP/2 streams share one connection
HTTP/1.1 generally needs separate connections for simultaneous independent responses. HTTP/2 splits a connection into numbered streams with interleaved frames, so a held request need not block a sibling. The application still sees ordinary ~http.Request~ and ~http.ResponseWriter~; transport framing belongs to the standard library.

This mission adds no server push. Build a concurrent-safe handler from an injected ~SectorRead~ and return an ~http.Server~. Avoid a global lock around a potentially blocking sensor read; the injected function already supports concurrent use.
~~~go
srv := &http.Server{Handler: h, ReadHeaderTimeout: 2*time.Second}
~~~

## ALPN chooses the application protocol
TLS authenticates the peers; ALPN negotiates what they speak after the handshake. ~tls.Config.NextProtos~ is the ordered list of application protocols. ~h2~ means HTTP/2; ~http/1.1~ allows fallback for older clients. Both ends must offer a common choice.
~~~go
secure := config.Clone()
secure.NextProtos = []string{"h2", "http/1.1"}
secure.MinVersion = tls.VersionTLS13
~~~
Clone the caller's configuration before changes. ~tls.Config.Clone~ is shallow, so replace the ~NextProtos~ slice rather than writing into its existing backing array. Certificate objects stay immutable as in the TLS mission.

~srv.ServeTLS(ln, "", "")~ uses certificates already present in ~TLSConfig~ and configures the standard HTTP/2 server support. Do not populate ~TLSNextProto~ with an empty map: that disables automatic HTTP/2. The test owns the listener and server lifetime; your factory opens no ports.
~~~go
srv := &http.Server{Handler: h, TLSConfig: secure}
err := srv.ServeTLS(listener, "", "")
~~~

## Opt the client in
A custom ~http.Transport~ with custom TLS settings should explicitly use ~ForceAttemptHTTP2: true~. ~MaxConnsPerHost: 1~ bounds connections, not HTTP/2 streams. Trust certificates normally through ~TLSClientConfig~; never bypass verification for HTTP/2.
~~~go
transport := &http.Transport{
 TLSClientConfig: clientTLS, ForceAttemptHTTP2: true, MaxConnsPerHost: 1,
}
client := &http.Client{Transport: transport}
defer transport.CloseIdleConnections()
~~~
~resp.ProtoMajor~ reports the HTTP version. ~resp.TLS.NegotiatedProtocol~ reports ALPN. A protocol list by itself does not prove HTTP/2 was actually used.
~~~go
fmt.Println(resp.ProtoMajor, resp.TLS.NegotiatedProtocol) // 2 h2
~~~

## Prove concurrency with barriers
~httptrace.WithClientTrace~ attaches a client observer. ~ClientTrace.GotConn~ receives ~GotConnInfo~ whose ~Conn~ is the actual connection object. Compare identity across requests; ~Reused~ alone does not prove the same connection. The test warms the connection, holds one sensor behind a channel barrier, and requires siblings to finish on that connection before release.
~~~go
trace := &httptrace.ClientTrace{GotConn: func(info httptrace.GotConnInfo) {
 observed <- info.Conn
}}
ctx = httptrace.WithClientTrace(ctx, trace)
~~~
~httptest.NewUnstartedServer~ lets tests install the returned server configuration and TLS certificates before startup. ~EnableHTTP2~ opts the fixture into HTTP/2; ~StartTLS~ starts it on a random loopback port. Close bodies, idle transports, and the server, and release blocked handlers before closing.
~~~go
ts := httptest.NewUnstartedServer(handler)
ts.EnableHTTP2 = true
ts.StartTLS()
defer ts.Close()
~~~

## Common mistakes
~~~go
config.NextProtos[0] = "h2" // Wrong: mutates a caller-owned slice
srv.TLSNextProto = map[string]func(*http.Server, *tls.Conn, http.Handler){} // Disables h2
// Holding a shared mutex while waiting for a sensor also blocks sibling streams.
~~~

## Your mission
Implement ~NewMultiplexGate(config *tls.Config, read SectorRead) (*http.Server, error)~. Reject nil config, no configured certificates, or nil read with nil server and an error. Other supplied TLS fields and certificates are valid and immutable. Clone config, set minimum TLS 1.3 and ~NextProtos~ exactly ~["h2", "http/1.1"]~, and return a fresh ~http.Server~ with that config, ~ReadHeaderTimeout~ 2 seconds, ~IdleTimeout~ 30 seconds, and ~MaxHeaderBytes~ 16384. Preserve the injected verification clock and certificates. Do not open a listener or disable automatic HTTP/2.

Install a handler for ~GET /sectors/NAME~. NAME is 1..16 ASCII lowercase letters, digits, or hyphens. Any other path, including extra segments or uppercase names, returns 404 before reading. A different method on a valid path returns 405 and ~Allow: GET~, without reading. Ignore query parameters. Call ~read(r.Context(), NAME)~ once; success returns 200 with its exact bytes, ~Content-Type: text/plain; charset=utf-8~, and ~Cache-Control: no-store~. Any read error returns 503 with generic text and no returned sensor contents or private error details. Error bodies are otherwise unspecified.

Concurrent requests must remain independent. A blocked read must not block siblings on the same HTTP/2 connection. The tests require actual ALPN ~h2~, HTTP major version 2, shared connection identity, HTTP/1.1 fallback, and rejection of untrusted certificates. No Server Push, ~Report~, or work in ~main~.
`,
			Hints: []string{"Clone tls.Config, replace NextProtos with a new slice, and set MinVersion and the http.Server timeouts.", "Validate the route and method before sensor reads; pass r.Context without a shared lock.", "Leave TLSNextProto alone so ServeTLS configures HTTP/2; tests use ForceAttemptHTTP2 and GotConn identity rather than timing."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch325Hidden},
	})
}

const ch325Starter = `package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
)

type SectorRead func(context.Context, string) (string, error)

func NewMultiplexGate(config *tls.Config, read SectorRead) (*http.Server, error) {
	return nil, errors.New("TODO")
}
func main() {}
`
const ch325Hidden = `package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"orbital/colony/labkit"
)

func multiplexConfig(t *testing.T) (*tls.Config, *tls.Config) {
	t.Helper()
	ca, err := labkit.TestCA()
	if err != nil {
		t.Fatal(err)
	}
	now := ca.Certificate.NotBefore.Add(2 * time.Hour)
	cert, err := labkit.IssueCert(ca, labkit.CertOptions{DNSNames: []string{"gate.colony"}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour)})
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(ca.Certificate)
	return &tls.Config{Certificates: []tls.Certificate{cert}, Time: func() time.Time { return now }}, &tls.Config{RootCAs: roots, ServerName: "gate.colony", MinVersion: tls.VersionTLS13, Time: func() time.Time { return now }}
}
func TestHiddenMultiplexHandlerAndSettings(t *testing.T) {
	config, _ := multiplexConfig(t)
	rng := rand.New(rand.NewSource(3205))
	readings := map[string]string{}
	for i := 0; i < 20; i++ {
		readings[fmt.Sprintf("d-%d", i)] = fmt.Sprintf("pressure=%d\n", rng.Intn(10000))
	}
	var calls int
	srv, err := NewMultiplexGate(config, func(ctx context.Context, sector string) (string, error) {
		calls++
		if sector == "broken" {
			return "private", errors.New("secret")
		}
		return readings[sector], nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if srv.TLSConfig == config || config.MinVersion != 0 || len(config.NextProtos) != 0 || srv.TLSConfig.MinVersion != tls.VersionTLS13 || !reflect.DeepEqual(srv.TLSConfig.NextProtos, []string{"h2", "http/1.1"}) || srv.ReadHeaderTimeout != 2*time.Second || srv.IdleTimeout != 30*time.Second || srv.MaxHeaderBytes != 16<<10 {
		t.Fatal("TLS/HTTP settings or caller ownership violated")
	}
	for name, want := range readings {
		w := httptest.NewRecorder()
		srv.Handler.ServeHTTP(w, httptest.NewRequest("GET", "https://gate/sectors/"+name+"?ignored=1", nil))
		if w.Code != 200 || w.Body.String() != want || w.Header().Get("Cache-Control") != "no-store" || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/plain") {
			t.Fatalf("reading %q: %d %q", name, w.Code, w.Body.String())
		}
	}
	for _, tc := range []struct {
		method, path string
		code         int
	}{
		{"GET", "/", 404}, {"GET", "/sectors/", 404}, {"GET", "/sectors/a/x", 404}, {"GET", "/sectors/UPPER", 404}, {"GET", "/sectors/" + strings.Repeat("a", 17), 404}, {"POST", "/sectors/a", 405}, {"GET", "/sectors/broken", 503},
	} {
		before := calls
		w := httptest.NewRecorder()
		srv.Handler.ServeHTTP(w, httptest.NewRequest(tc.method, "https://gate"+tc.path, nil))
		if w.Code != tc.code || strings.Contains(w.Body.String(), "secret") || strings.Contains(w.Body.String(), "private") {
			t.Fatalf("%s %s: %d %q", tc.method, tc.path, w.Code, w.Body.String())
		}
		if tc.code != 503 && calls != before {
			t.Fatal("rejected route read sensor")
		}
		if tc.code == 405 && w.Header().Get("Allow") != "GET" {
			t.Fatal("missing Allow")
		}
	}
	for _, tc := range []struct {
		cfg  *tls.Config
		read SectorRead
	}{{nil, func(context.Context, string) (string, error) { return "", nil }}, {&tls.Config{}, func(context.Context, string) (string, error) { return "", nil }}, {config, nil}} {
		if s, err := NewMultiplexGate(tc.cfg, tc.read); s != nil || err == nil {
			t.Fatal("invalid options accepted")
		}
	}
}
func multiplexServer(t *testing.T, read SectorRead) (*httptest.Server, *http.Transport) {
	t.Helper()
	labkit.LeakCheck(t)
	config, clientConfig := multiplexConfig(t)
	srv, err := NewMultiplexGate(config, read)
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewUnstartedServer(srv.Handler)
	ts.Config = srv
	ts.TLS = srv.TLSConfig.Clone()
	ts.EnableHTTP2 = true
	ts.StartTLS()
	t.Cleanup(ts.Close)
	tr := &http.Transport{TLSClientConfig: clientConfig, ForceAttemptHTTP2: true, MaxConnsPerHost: 1}
	t.Cleanup(tr.CloseIdleConnections)
	return ts, tr
}

type multiplexReply struct {
	body  string
	major int
	alpn  string
	conn  net.Conn
	err   error
}

func multiplexGet(ctx context.Context, client *http.Client, url string) multiplexReply {
	var reply multiplexReply
	ctx = httptrace.WithClientTrace(ctx, &httptrace.ClientTrace{GotConn: func(c httptrace.GotConnInfo) { reply.conn = c.Conn }})
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		reply.err = err
		return reply
	}
	resp, err := client.Do(req)
	if err != nil {
		reply.err = err
		return reply
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	reply.body = string(body)
	reply.major = resp.ProtoMajor
	reply.alpn = resp.TLS.NegotiatedProtocol
	reply.err = err
	if resp.StatusCode != 200 {
		reply.err = fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return reply
}
func TestHiddenMultiplexOneConnectionDoesNotBlockSiblings(t *testing.T) {
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	defer once.Do(func() { close(release) })
	ts, tr := multiplexServer(t, func(ctx context.Context, sector string) (string, error) {
		if sector == "hold" {
			close(entered)
			select {
			case <-release:
			case <-ctx.Done():
				return "", ctx.Err()
			}
		}
		return "sector:" + sector, nil
	})
	// Register after the server cleanup so failure releases a blocked handler first.
	t.Cleanup(func() { once.Do(func() { close(release) }) })
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := &http.Client{Transport: tr}
	warm := multiplexGet(ctx, client, ts.URL+"/sectors/warm")
	if warm.err != nil || warm.major != 2 || warm.alpn != "h2" || warm.conn == nil {
		t.Fatalf("h2 warmup: %+v", warm)
	}
	held := make(chan multiplexReply, 1)
	go func() { held <- multiplexGet(ctx, client, ts.URL+"/sectors/hold") }()
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("held handler absent")
	}
	replies := make(chan multiplexReply, 12)
	for i := 0; i < 12; i++ {
		go func(i int) { replies <- multiplexGet(ctx, client, fmt.Sprintf("%s/sectors/d-%d", ts.URL, i)) }(i)
	}
	for i := 0; i < 12; i++ {
		select {
		case reply := <-replies:
			if reply.err != nil || reply.major != 2 || reply.conn != warm.conn || !strings.HasPrefix(reply.body, "sector:d-") {
				t.Fatalf("sibling: %+v", reply)
			}
		case <-ctx.Done():
			t.Fatal("blocked sibling or second connection")
		}
	}
	once.Do(func() { close(release) })
	select {
	case reply := <-held:
		if reply.err != nil || reply.conn != warm.conn || reply.body != "sector:hold" {
			t.Fatalf("held: %+v", reply)
		}
	case <-ctx.Done():
		t.Fatal("held request leaked")
	}
}
func TestHiddenMultiplexFallbackAndTrust(t *testing.T) {
	ts, tr := multiplexServer(t, func(context.Context, string) (string, error) { return "ready", nil })
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	fallback := &http.Transport{TLSClientConfig: tr.TLSClientConfig.Clone(), TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{}}
	defer fallback.CloseIdleConnections()
	got := multiplexGet(ctx, &http.Client{Transport: fallback}, ts.URL+"/sectors/a")
	if got.err != nil || got.major != 1 {
		t.Fatalf("HTTP/1 fallback: %+v", got)
	}
	untrusted := &http.Transport{TLSClientConfig: &tls.Config{ServerName: "gate.colony", RootCAs: x509.NewCertPool()}, ForceAttemptHTTP2: true}
	defer untrusted.CloseIdleConnections()
	if got := multiplexGet(ctx, &http.Client{Transport: untrusted}, ts.URL+"/sectors/a"); got.err == nil {
		t.Fatal("untrusted certificate accepted")
	}
}
`

func init() {
	add(&Level{ID: "32-6-ssh", Chapter: 32, Kind: "lab", Title: "SSH", Goal: "نفّذ أوامر فحص الرصيف المحددة بمفتاح مضيف موثوق وإلغاء كامل",
		Requires: []string{"sel:ssh.NewClientConn", "sel:ssh.FixedHostKey"}, Deps: []string{"golang.org/x/crypto"}, Bridge: []int{81}, Race: true, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch326Starter}},
		Lesson: `# وصلة SSH بمفتاح رصيف موثوق

## قبل ما تكتب
- قبول أي host key يثبت هوية الرصيف؟
- الإلغاء هيفك مصافحة أو أمر وقف عن الرد إزاي؟

## هوية المضيف منفصلة عن تصريح المستخدم
SSH بيشفّر الاتصال وبيفحص هويتين مختلفتين. مفتاح المضيف بيثبت السيرفر؛ طرق ~Auth~ بتصرّح للمستخدم. استورد ~golang.org/x/crypto/ssh~. لازم ~ssh.PublicKey~ الموثوق يوصلك من قناة موثوقة، مش تتعلمه من نفس الاتصال اللي لسه متحققناش منه.

~ssh.FixedHostKey~ بتعمل ~HostKeyCallback~ بيقبل المفتاح المثبّت بس. ~ssh.ClientConfig~ فيها الـ callback و ~User~ و ~Auth~. المستدعي بيحقن طرق تصريح زي ~ssh.Password~ أو ~ssh.PublicKeys~؛ مفيش ملف اعتماد في المهمة.
~~~go
config := &ssh.ClientConfig{
 User: "engineer", Auth: []ssh.AuthMethod{ssh.Password("example-only")},
 HostKeyCallback: ssh.FixedHostKey(trustedKey),
}
~~~
متستخدمش ~ssh.InsecureIgnoreHostKey~. اختلاف مفتاح المضيف لازم يفشل قبل تصريح المستخدم أو تنفيذ الجلسة. حد الاتصال المحقون بيرجع ~net.Conn~ متصلة؛ ده يخلي التست يختار مستمع loopback من غير بورت ثابت ولا DNS.

## ابدأ SSH على اتصال محقون
~ssh.NewClientConn~ بتعمل مصافحة SSH على اتصال. نص العنوان بيعرّف الهدف للـ callback؛ فاحص المفتاح المثبّت مش بيستخدم DNS. الاتصال والقنوات الراجعة بيتلفوا بـ ~ssh.NewClient~ اللي بيدير طلبات وقنوات SSH. اقفل العميل واستنى حلقة اتصاله تقف.
~~~go
conn, channels, requests, err := ssh.NewClientConn(raw, "dock", config)
if err != nil { return err }
client := ssh.NewClient(conn, channels, requests)
defer client.Close()
~~~
~NewSession~ بتفتح قناة جلسة مستقلة عن الاتصال المشفّر. ~session.Stdout~ و ~session.Stderr~ مخارج ~io.Writer~ منفصلة. ~session.Run~ بتبعت طلب exec وتستنى نهايته. سيرفر التست بيفهم الطلب كنص محاكاة؛ عمره ما بينادي shell ولا أمر على الجهاز.
~~~go
session, err := client.NewSession()
if err != nil { return err }
defer session.Close()
var stdout, stderr bytes.Buffer
session.Stdout, session.Stderr = &stdout, &stderr
err = session.Run("colony status --summary")
~~~

## الأوامر اختيارات محددة
حوّل اسم أمر التطبيق بالظبط لنص بروتوكول ثابت. متوصلش مدخل المستخدم بنص exec، ولا تعمل trim لأمر مرفوض عشان تعدّيه، ولا تطلب shell أو PTY. افحص الطلب كله قبل الاتصال حتى لو السياق متلغي بالفعل.
~~~go
switch choice {
case "status": command = "colony status --summary"
default: return ErrDockSSHPlan
}
~~~

## الإلغاء لازم يقفل الـ I/O
~ssh.NewClientConn~ و ~NewSession~ و ~Run~ مش بياخدوا context. فالسياق لوحده مش هيوقفهم. شغّل مراقب واحد بعد الاتصال: يختار بين ~ctx.Done()~ وقناة انتهاء؛ عند الإلغاء يقفل الاتصال الخام. القفل بيفك المصافحة والجلسة. في كل خروج وقف المراقب واستناه؛ اقفل الجلسة والعميل والنقل الخام. قفل net connection أكتر من مرة آمن.
~~~go
stopped := make(chan struct{})
go func() {
 select { case <-ctx.Done(): raw.Close(); case <-stopped: }
}()
// عند الرجوع: close(stopped)، وبعدها استنى قناة انتهاء المراقب نفسه.
~~~
المستدعي بيملك المهلة بسياق deadline؛ مفيش ساعة حقيقية مخفية ولا نوم. افحص السياق قبل الاتصال وبعد فشل النقل أو الجلسة عشان الإلغاء يرجع ~ctx.Err()~ بدل خطأ اتصال مقفول بالصدفة. المهلة المنتهية من قبل البداية مش محتاجة شبكة.

## حافظ على حالة خروج الأمر
خروج بعيد مش صفر بيبقى ~*ssh.ExitError~. استخدم ~errors.As~ و ~ExitStatus~ عشان تطلع الحالة وتحافظ على الخطأ الأصلي والمخرجات. أخطاء النقل التانية بترجع من غير مخرجات جزئية في العقد هنا.
~~~go
var exit *ssh.ExitError
if errors.As(err, &exit) { code := exit.ExitStatus(); _ = code }
~~~
السيرفر المخفي بيستخدم ~ssh.ServerConfig~ و ~AddHostKey~ و ~ssh.NewServerConn~ جوه العملية. بيقبل جلسة ويفك طلب exec ويكتب stdout وstderr محاكيين ويبعت حالة الخروج. المخرجات العشوائية بـ seed بتمنع الإجابة الثابتة؛ مفيش مضيف خارجي ولا عملية نظام.

## أخطاء شائعة
~~~go
config.HostKeyCallback = ssh.InsecureIgnoreHostKey() // غلط: بيقبل منتحل هوية
session.Run("colony " + input) // غلط: نص بعيد حر بدل اختيار ثابت
<-ctx.Done() // غلط لوحدها: مش بتفك قراية SSH واقفة
~~~

## المهمة
نفّذ ~ExecuteDockCommand(ctx context.Context, req DockSSHRequest) (DockSSHReply, error)~. قبل الاتصال ارفض ~User~ فاضي أو أكبر من 16 بايت أو فيه غير حروف ASCII صغيرة وأرقام وشرطات، أو ~HostKey~ بـ ~nil~، أو ~Auth~ فاضية، أو ~Dial~ بـ ~nil~، أو ~Command~ مجهول بـ ~ErrDockSSHPlan~ قابل للفحص بـ ~errors.Is~ ورد صفري. الاختيارين بالظبط ~pressure~ و ~seal~ بيتحوّلوا على الترتيب لـ ~colony pressure --brief~ و ~colony seal --check~. متعملش trim ولا تركيب نص ولا shell.

لو الطلب صحيح وسياقه متلغي أو منتهي من البداية، رجّع خطأ السياق من غير اتصال. نادِ ~Dial(ctx)~ المحقونة مرة؛ بترجع اتصال مش ~nil~ ومن غير خطأ أو ~nil~ وخطأ. حافظ على خطأ الاتصال إلا لو السياق فشل، ساعتها خطأه له الأولوية. استلم ملكية الاتصال. صرّح بالمستخدم والطرق المحقونة وثبّت مفتاح المضيف بالظبط. افتح جلسة واحدة ونفّذ النص المقابل مرة، والتقط stdout وstderr منفصلين.

خروج 0 يرجع النصين و ~ExitCode~ بصفر وخطأ ~nil~. خروج بعيد غير صفر يرجع النصين وكود الخروج وخطأ ~*ssh.ExitError~ الأصلي كسبب. أخطاء المصافحة والتصريح والجلسة والـ I/O التانية ترجع رد صفري وخطأ. الإلغاء في أي مرحلة يرجع رد صفري و ~ctx.Err()~؛ لازم يقفل النقل بسرعة ويوقف الجلسة ويستنى مراقب الإلغاء قبل الرجوع. السياقات والاتصال المحقون لازم يتعاونوا وقت الاتصال. النداءات المتزامنة مستقلة. متقراش اعتمادات الجهاز ولا تطبع مفاتيح أو اعتمادات ولا تنفّذ أوامر محلية. مفيش ~Report~ و ~main~ فاضية.
`,
		Hints: []string{"تحقق من الطلب وحوّل enum الأمر لنص ثابت قبل Dial؛ ثبّت HostKeyCallback بـ ssh.FixedHostKey.", "بعد الاتصال شغّل مراقب إلغاء بيقفل النقل، وبعد NewClientConn لف القنوات بـ ssh.NewClient وافتح NewSession.", "التقط stdout وstderr في buffers منفصلة، وافحص ssh.ExitError بـ errors.As؛ أولوية ctx.Err واستنى المراقب عند كل خروج."},
		En: &LevelText{Title: "SSH", Goal: "Execute allowlisted dock checks with a pinned host key and complete cancellation", Files: []LabFile{{Name: "main.go", Starter: ch326Starter}},
			Lesson: `# A Pinned SSH Link to the Dock

## Before you write
- Does accepting any host key authenticate the dock?
- How will cancellation interrupt a handshake or command that has stopped responding?

## Host identity is separate from user authentication
SSH encrypts a connection and verifies two different identities. The host key identifies the server; ~Auth~ methods authenticate the user. Import ~golang.org/x/crypto/ssh~. A trusted ~ssh.PublicKey~ must arrive through a trusted channel, not be learned from the same unverified connection.

~ssh.FixedHostKey~ constructs a ~HostKeyCallback~ that accepts only the pinned key. ~ssh.ClientConfig~ supplies that callback, the ~User~, and ~Auth~. The caller injects auth methods such as ~ssh.Password~ or ~ssh.PublicKeys~; no credential file belongs in this lab.
~~~go
config := &ssh.ClientConfig{
 User: "engineer", Auth: []ssh.AuthMethod{ssh.Password("example-only")},
 HostKeyCallback: ssh.FixedHostKey(trustedKey),
}
~~~
Never use ~ssh.InsecureIgnoreHostKey~. A host mismatch must fail before user authentication or session execution. The injected dial boundary returns a connected ~net.Conn~; this lets tests choose a loopback listener without hard-coded ports or DNS.

## Upgrade an injected connection
~ssh.NewClientConn~ performs the SSH handshake on a connection. Its address string identifies the target to the callback; our pinned-key callback does not use DNS. The returned connection and channels are wrapped with ~ssh.NewClient~, which handles SSH request and channel routing. Close the client and wait for its connection loop to stop.
~~~go
conn, channels, requests, err := ssh.NewClientConn(raw, "dock", config)
if err != nil { return err }
client := ssh.NewClient(conn, channels, requests)
defer client.Close()
~~~
~NewSession~ opens a session channel, independent of the encrypted connection. ~session.Stdout~ and ~session.Stderr~ are separate ~io.Writer~ sinks. ~session.Run~ sends an exec request and waits for completion. Our test server interprets the request as a simulated string; it never invokes a shell or a machine command.
~~~go
session, err := client.NewSession()
if err != nil { return err }
defer session.Close()
var stdout, stderr bytes.Buffer
session.Stdout, session.Stderr = &stdout, &stderr
err = session.Run("colony status --summary")
~~~

## Commands are enums, not user-controlled shell text
Map an exact application command name to a fixed protocol string. Never concatenate user input into exec text, trim unknown names into accepted commands, request a shell, or request a PTY. Validate the whole request before dialing, even when the context is already canceled.
~~~go
switch choice {
case "status": command = "colony status --summary"
default: return ErrDockSSHPlan
}
~~~

## Context cancellation must close I/O
~ssh.NewClientConn~, ~NewSession~, and ~Run~ do not take a context. A context alone therefore cannot stop them. Start one watcher after dialing: select between ~ctx.Done()~ and a completion channel; on cancellation close the raw connection. Closing it unblocks the handshake and session. On every exit stop and join the watcher; close the session, client, and underlying transport. Closing a net connection more than once is safe.
~~~go
stopped := make(chan struct{})
go func() {
 select { case <-ctx.Done(): raw.Close(); case <-stopped: }
}()
// On return: close(stopped), then wait for the watcher's own done channel.
~~~
The caller owns the timeout via a deadline context; no hidden real-time clock or sleep is needed. Check the context before dialing and after a transport or session failure, so cancellation returns ~ctx.Err()~ rather than an incidental closed-connection error. A pre-expired deadline needs no network at all.

## Preserve command exit status
A nonzero remote exit is ~*ssh.ExitError~. Use ~errors.As~ and ~ExitStatus~ to expose that status while preserving the original error and captured output. Other transport failures return no partial output in our contract.
~~~go
var exit *ssh.ExitError
if errors.As(err, &exit) { code := exit.ExitStatus(); _ = code }
~~~
The hidden server uses ~ssh.ServerConfig~, ~AddHostKey~, and ~ssh.NewServerConn~ in-process. It accepts a session, decodes an exec request, writes simulated stdout and stderr, and sends exit status. Seeded output prevents a canned answer from succeeding; no external host or process is used.

## Common mistakes
~~~go
config.HostKeyCallback = ssh.InsecureIgnoreHostKey() // Wrong: accepts an impostor
session.Run("colony " + input) // Wrong: arbitrary remote text instead of a fixed enum
<-ctx.Done() // Wrong by itself: does not interrupt a blocked SSH read
~~~

## Your mission
Implement ~ExecuteDockCommand(ctx context.Context, req DockSSHRequest) (DockSSHReply, error)~. Before dialing reject an empty or over-16-byte ~User~, any username byte other than ASCII lowercase letters, digits, or hyphens, nil ~HostKey~, empty ~Auth~, nil ~Dial~, or unknown ~Command~ with ~ErrDockSSHPlan~ (preservable with ~errors.Is~) and a zero reply. Only exact ~pressure~ and ~seal~ are accepted, mapping respectively to ~colony pressure --brief~ and ~colony seal --check~. Do not trim, concatenate, or run a shell.

For a valid request, a pre-canceled or expired context returns its error with zero dials. Call injected ~Dial(ctx)~ once; it returns either a non-nil connection and nil error, or nil and an error. Preserve dial failures unless the context failed, in which case return its error. Take ownership of the connection. Authenticate with the supplied user and auth methods and verify exactly the pinned host key. Open one session and run the mapped command once, capturing stdout and stderr independently.

Exit 0 returns captured strings, ~ExitCode~ 0, and nil error. Nonzero remote exit returns both captured strings, its exit code, and the original ~*ssh.ExitError~ as a cause. Handshake, authentication, session, and other I/O failures return a zero reply and an error. Cancellation at any stage returns a zero reply and ~ctx.Err()~; it must close the transport promptly, stop the session, and join the cancellation watcher before returning. Contexts and injected dialers must cooperate during dialing. Concurrent calls share no state. Do not read machine credentials, log keys or credentials, or execute local commands. No ~Report~; ~main~ is empty.
`,
			Hints: []string{"Validate the request and map its command enum to fixed text before Dial; pin HostKeyCallback with ssh.FixedHostKey.", "After dialing start a cancellation watcher that closes the transport, wrap NewClientConn channels with ssh.NewClient, then open NewSession.", "Capture stdout and stderr separately and inspect ssh.ExitError with errors.As; prioritize ctx.Err and join the watcher on every exit."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch326Hidden},
	})
}

const ch326Starter = `package main

import (
	"context"
	"errors"
	"golang.org/x/crypto/ssh"
	"net"
)

var ErrDockSSHPlan = errors.New("invalid dock SSH request")

type DockSSHRequest struct {
	User    string
	HostKey ssh.PublicKey
	Auth    []ssh.AuthMethod
	Command string
	Dial    func(context.Context) (net.Conn, error)
}
type DockSSHReply struct {
	Stdout, Stderr string
	ExitCode       int
}

func ExecuteDockCommand(ctx context.Context, req DockSSHRequest) (DockSSHReply, error) {
	return DockSSHReply{}, errors.New("TODO")
}
func main() {}
`
const ch326Hidden = `package main

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
	"orbital/colony/labkit"
)

func dockSigner(t *testing.T, seed int64) ssh.Signer {
	t.Helper()
	rng := rand.New(rand.NewSource(seed))
	keySeed := make([]byte, ed25519.SeedSize)
	rng.Read(keySeed)
	signer, err := ssh.NewSignerFromKey(ed25519.NewKeyFromSeed(keySeed))
	if err != nil {
		t.Fatal(err)
	}
	return signer
}
func dockWait[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(3 * time.Second):
		t.Fatal("SSH lifecycle stuck")
	}
	var zero T
	return zero
}
func TestHiddenDockRejectsBeforeDial(t *testing.T) {
	key := dockSigner(t, 3206).PublicKey()
	var calls int
	failed := errors.New("injected dial failure")
	base := DockSSHRequest{User: "crew-1", HostKey: key, Auth: []ssh.AuthMethod{ssh.Password("badge")}, Command: "pressure", Dial: func(context.Context) (net.Conn, error) { calls++; return nil, failed }}
	for _, change := range []func(*DockSSHRequest){
		func(r *DockSSHRequest) { r.User = "" }, func(r *DockSSHRequest) { r.User = "Crew" }, func(r *DockSSHRequest) { r.User = strings.Repeat("a", 17) }, func(r *DockSSHRequest) { r.HostKey = nil }, func(r *DockSSHRequest) { r.Auth = nil }, func(r *DockSSHRequest) { r.Dial = nil },
		func(r *DockSSHRequest) { r.Command = "" }, func(r *DockSSHRequest) { r.Command = "pressure; reboot" }, func(r *DockSSHRequest) { r.Command = " pressure" }, func(r *DockSSHRequest) { r.Command = "shell" },
	} {
		req := base
		change(&req)
		reply, err := ExecuteDockCommand(context.Background(), req)
		if !errors.Is(err, ErrDockSSHPlan) || calls != 0 || !reflect.DeepEqual(reply, DockSSHReply{}) {
			t.Fatalf("rejected: %+v %v calls=%d", reply, err, calls)
		}
	}
	for _, ctx := range []context.Context{func() context.Context { c, cancel := context.WithCancel(context.Background()); cancel(); return c }(), func() context.Context {
		c, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
		t.Cleanup(cancel)
		return c
	}()} {
		reply, err := ExecuteDockCommand(ctx, base)
		if !errors.Is(err, ctx.Err()) || calls != 0 || reply != (DockSSHReply{}) {
			t.Fatalf("pre-cancel: %+v %v", reply, err)
		}
	}
	if _, err := ExecuteDockCommand(context.Background(), base); !errors.Is(err, failed) || calls != 1 {
		t.Fatalf("dial cause: %v", err)
	}
}
func TestHiddenDockHandshakeCancellationClosesTransport(t *testing.T) {
	labkit.LeakCheck(t)
	left, right := net.Pipe()
	defer right.Close()
	defer left.Close()
	entered := make(chan struct{})
	req := DockSSHRequest{User: "crew", HostKey: dockSigner(t, 2).PublicKey(), Auth: []ssh.AuthMethod{ssh.Password("badge")}, Command: "seal", Dial: func(context.Context) (net.Conn, error) { close(entered); return left, nil }}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		reply, err := ExecuteDockCommand(ctx, req)
		if reply != (DockSSHReply{}) {
			t.Error("partial reply on cancellation")
		}
		done <- err
	}()
	dockWait(t, entered)
	cancel()
	if err := dockWait(t, done); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

type dockSimulation struct {
	stdout, stderr string
	exit           uint32
	hold           bool
}

func dockSSHFixture(t *testing.T, sim dockSimulation) (DockSSHRequest, <-chan string, *atomic.Int32) {
	t.Helper()
	labkit.LeakCheck(t)
	signer := dockSigner(t, 32060)
	var authenticated atomic.Int32
	config := &ssh.ServerConfig{PasswordCallback: func(c ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		authenticated.Add(1)
		if c.User() != "crew" || string(password) != "badge" {
			return nil, errors.New("denied")
		}
		return nil, nil
	}}
	config.AddHostKey(signer)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	command := make(chan string, 1)
	done := make(chan struct{})
	stop := make(chan struct{})
	var mu sync.Mutex
	var accepted net.Conn
	go func() {
		defer close(done)
		raw, err := ln.Accept()
		if err != nil {
			return
		}
		mu.Lock()
		accepted = raw
		select {
		case <-stop:
			raw.Close()
		default:
		}
		mu.Unlock()
		defer raw.Close()
		conn, channels, requests, err := ssh.NewServerConn(raw, config)
		if err != nil {
			return
		}
		defer conn.Close()
		discarded := make(chan struct{})
		go func() { ssh.DiscardRequests(requests); close(discarded) }()
		defer func() { conn.Close(); <-discarded }()
		incoming, ok := <-channels
		if !ok {
			return
		}
		if incoming.ChannelType() != "session" {
			incoming.Reject(ssh.UnknownChannelType, "session only")
			return
		}
		channel, reqs, err := incoming.Accept()
		if err != nil {
			return
		}
		defer channel.Close()
		for r := range reqs {
			if r.Type != "exec" {
				r.Reply(false, nil)
				continue
			}
			var payload struct{ Command string }
			if ssh.Unmarshal(r.Payload, &payload) != nil {
				r.Reply(false, nil)
				return
			}
			command <- payload.Command
			r.Reply(true, nil)
			if sim.hold {
				<-stop
				return
			}
			channel.Write([]byte(sim.stdout))
			channel.Stderr().Write([]byte(sim.stderr))
			channel.SendRequest("exit-status", false, ssh.Marshal(struct{ Status uint32 }{sim.exit}))
			return
		}
	}()
	t.Cleanup(func() {
		close(stop)
		ln.Close()
		mu.Lock()
		if accepted != nil {
			accepted.Close()
		}
		mu.Unlock()
		dockWait(t, done)
	})
	req := DockSSHRequest{User: "crew", HostKey: signer.PublicKey(), Auth: []ssh.AuthMethod{ssh.Password("badge")}, Command: "pressure", Dial: func(ctx context.Context) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp", ln.Addr().String())
	}}
	return req, command, &authenticated
}
func TestHiddenDockExecCapturesExitAndPinsHost(t *testing.T) {
	rng := rand.New(rand.NewSource(32061))
	for _, exit := range []uint32{0, 7} {
		t.Run(fmt.Sprint(exit), func(t *testing.T) {
			sim := dockSimulation{stdout: fmt.Sprintf("ضغط=%d\n", rng.Intn(100000)), stderr: "تنبيه 🚀\n", exit: exit}
			req, command, _ := dockSSHFixture(t, sim)
			if exit != 0 {
				req.Command = "seal"
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			reply, err := ExecuteDockCommand(ctx, req)
			if reply.Stdout != sim.stdout || reply.Stderr != sim.stderr || reply.ExitCode != int(exit) {
				t.Fatalf("reply %+v %v", reply, err)
			}
			if exit == 0 && err != nil {
				t.Fatal(err)
			}
			if exit != 0 {
				var ee *ssh.ExitError
				if !errors.As(err, &ee) || ee.ExitStatus() != int(exit) {
					t.Fatalf("exit error: %v", err)
				}
			}
			want := "colony pressure --brief"
			if exit != 0 {
				want = "colony seal --check"
			}
			if got := dockWait(t, command); got != want {
				t.Fatalf("executed %q want %q", got, want)
			}
		})
	}
	t.Run("host-mismatch", func(t *testing.T) {
		req, _, auth := dockSSHFixture(t, dockSimulation{})
		req.HostKey = dockSigner(t, 99).PublicKey()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		reply, err := ExecuteDockCommand(ctx, req)
		if err == nil || reply != (DockSSHReply{}) || auth.Load() != 0 {
			t.Fatalf("host pin bypassed: %+v %v auth=%d", reply, err, auth.Load())
		}
	})
	t.Run("auth-denied", func(t *testing.T) {
		req, _, _ := dockSSHFixture(t, dockSimulation{})
		req.Auth = []ssh.AuthMethod{ssh.Password("wrong")}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if reply, err := ExecuteDockCommand(ctx, req); err == nil || reply != (DockSSHReply{}) {
			t.Fatalf("auth bypassed: %+v %v", reply, err)
		}
	})
}
func TestHiddenDockCommandCancellation(t *testing.T) {
	req, entered, _ := dockSSHFixture(t, dockSimulation{hold: true})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() {
		reply, err := ExecuteDockCommand(ctx, req)
		if reply != (DockSSHReply{}) {
			t.Error("partial cancellation output")
		}
		done <- err
	}()
	dockWait(t, entered)
	cancel()
	if err := dockWait(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("command cancellation: %v", err)
	}
}
func TestHiddenDockConcurrentCallsAreIndependent(t *testing.T) {
	type pending struct {
		req DockSSHRequest
		sim dockSimulation
	}
	var batch []pending
	rng := rand.New(rand.NewSource(32062))
	for i := 0; i < 4; i++ {
		sim := dockSimulation{stdout: fmt.Sprintf("oxygen-%d", rng.Intn(100000)), stderr: fmt.Sprintf("bay-%d", i)}
		req, _, _ := dockSSHFixture(t, sim)
		batch = append(batch, pending{req, sim})
	}
	var wg sync.WaitGroup
	for _, p := range batch {
		wg.Add(1)
		go func(p pending) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			reply, err := ExecuteDockCommand(ctx, p.req)
			if err != nil || reply.Stdout != p.sim.stdout || reply.Stderr != p.sim.stderr {
				t.Errorf("crossed output: %+v %v", reply, err)
			}
		}(p)
	}
	wg.Wait()
}
`
