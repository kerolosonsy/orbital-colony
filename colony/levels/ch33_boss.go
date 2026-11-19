package levels

func init() {
	add(&Level{
		ID: "33-boss-survey-uplink", Chapter: 33, Kind: "lab", Boss: true, Setup: seededLabRoom, Race: true,
		LessonRevision: 1, Steps: ch33Lesson33BossArabic, Title: "Boss: وصلة المسح", Goal: "وصّل إطار المنارة بخدمة HTTPS بموعد وإعادة محدودة",
		Requires: []string{"sel:binary.BigEndian", "sel:io.ReadFull", "sel:http.NewRequestWithContext"}, Bridge: []int{71, 74, 77, 78, 83, 84},
		Files: []LabFile{{Name: "main.go", Starter: `package main
import ("net";"net/http";"time";"errors")
type SurveyUplink struct { Client *http.Client; URL,Token string }
var ErrSurveyFrame=errors.New("invalid survey frame")
var ErrSurveyUpstream=errors.New("uplink rejected survey")
func RelaySurvey(conn net.Conn, uplink SurveyUplink, deadline time.Time) error { return ErrSurveyFrame }
func main() {}
`}},
		Lesson: `# Boss: وصلة المسح

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 31-1-echo, 31-4-frames, 32-4-tls, 33-1-retry-client, 33-2-reverse-proxy.

مثال أصغر ومختلف:
~~~go
var header [2]byte
if _, err := io.ReadFull(conn, header[:]); err != nil { return err }
~~~

## أخطاء شائعة
~~~go
resp, _ := client.Get(target) // drops errors and leaves the response body open
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~RelaySurvey(conn, uplink, deadline) error~. ~SurveyUplink~ فيها ~Client~ غير nil مجهّز للتحقق من TLS ومن غير redirects، و~URL~ موثوق و~Token~ سري. الدالة تملك conn وتقفلها عند كل خروج. ارفض URL مش https أو من غير host أو فيها userinfo بـ ~ErrSurveyUpstream~. متعدّلش عميل HTTP أو إعدادات شهاداته.

طبّق ~conn.SetDeadline(deadline)~ قبل القراءة. اقرأ رسالة واحدة: طول uint16 big-endian (بايتين) ثم اسم محطة بطول 1..64 بايت؛ صفر أو أكبر يرجّع ~ErrSurveyFrame~ من غير upstream. استخدم ReadFull عشان تجزئة الشبكة. المثال ~[]byte{0, 3, 'i', 'c', 'e'}~ معناه ice. الطلب GET لنفس URL مع query ~station~ بالاسم الأصلي وباقي query زي ما هي، و~Authorization: Bearer <Token>~، وسياق بموعد deadline.

503 بس له إعادة فورية واحدة؛ اقفل جسم كل رد قبل الإعادة أو الرجوع. أي status نهائي غير 200 يرجع ~ErrSurveyUpstream~. أخطاء النقل/القراءة/الكتابة تحتفظ بالسبب. رد 200 بحد 1024 بايت؛ الزيادة ~ErrSurveyFrame~ من غير رد جزئي. ابعت نفس الإطار بطول محتوى الرد (الصفر مسموح). سيب ~main~ فاضية. النجمة التالتة: شهادة غير موثوقة تفضل مرفوضة، والاتصال اللي مش بيبعت ينتهي بموعده. التستات TLS حقيقي فوق net.Pipe من غير بورت أو Docker.
`,
		Hints: ch33Lesson33BossHintsArabic,
		En: &LevelText{LessonRevision: 1, Steps: ch33Lesson33BossEnglish, Title: "Boss: Survey Uplink", Goal: "Bridge a beacon frame to HTTPS with a deadline and bounded retry", Files: []LabFile{{Name: "main.go", Starter: `package main
import ("net";"net/http";"time";"errors")
type SurveyUplink struct { Client *http.Client; URL,Token string }
var ErrSurveyFrame=errors.New("invalid survey frame")
var ErrSurveyUpstream=errors.New("uplink rejected survey")
func RelaySurvey(conn net.Conn, uplink SurveyUplink, deadline time.Time) error { return ErrSurveyFrame }
func main() {}
`}},
			Lesson: `# Boss: Survey Uplink

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 31-1-echo, 31-4-frames, 32-4-tls, 33-1-retry-client, 33-2-reverse-proxy.

A smaller, different example:
~~~go
var header [2]byte
if _, err := io.ReadFull(conn, header[:]); err != nil { return err }
~~~

## Common mistakes
~~~go
resp, _ := client.Get(target) // drops errors and leaves the response body open
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~RelaySurvey(conn, uplink, deadline) error~. ~SurveyUplink~ has a nonnil ~Client~ configured for certificate verification and no redirects, a trusted ~URL~, and a secret ~Token~. The function owns conn and closes it on every exit. Reject non-https URLs, missing hosts, or userinfo with ~ErrSurveyUpstream~. Never mutate the HTTP client or its certificate settings.

Apply ~conn.SetDeadline(deadline)~ before reading. Read one frame: a two-byte big-endian uint16 length followed by a 1..64-byte station name. Zero or oversized names return ~ErrSurveyFrame~ without upstream I/O. Use ReadFull to handle fragmentation. Example ~[]byte{0, 3, 'i', 'c', 'e'}~ names ice. Send GET to the same URL with query ~station~ set to the exact name, retaining other query values, ~Authorization: Bearer <Token>~, and a context bounded by deadline.

Only 503 gets one immediate retry; close every response body before retry or return. Any final non-200 status returns ~ErrSurveyUpstream~. Preserve transport/read/write causes. Bound 200 response content to 1024 bytes; overflow returns ~ErrSurveyFrame~ with no partial reply. Reply using the same frame format and response length (zero allowed). Leave ~main~ empty. Third star: untrusted certificates remain rejected and idle incoming connections expire at their deadline. Tests use real TLS over net.Pipe without ports or Docker.
`,
			Hints: ch33Lesson33BossHintsEnglish},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, HiddenFiles: map[string]string{"zz_transport_test.go": ch24HTTPHidden}, Hidden: `package main
import (
    "bytes"
    "context"
    "crypto/tls"
    "encoding/binary"
    "errors"
    "io"
    "math/rand"
    "net"
    "net/http"
    "net/http/httptest"
    "strings"
    "sync/atomic"
    "testing"
    "time"
    "orbital/colony/labkit"
)
func hiddenSurveyTLS(t *testing.T,h http.Handler)*httptest.Server{
    t.Helper()
    listener:=&supplyListener{connections:make(chan net.Conn),closed:make(chan struct{})}
    srv:=&httptest.Server{Listener:listener,Config:&http.Server{Handler:h}}
    srv.StartTLS();t.Cleanup(srv.Close)
    transport:=srv.Client().Transport.(*http.Transport)
    transport.DialContext=func(ctx context.Context,_,_ string)(net.Conn,error){
        client,remote:=net.Pipe()
        select{case listener.connections<-remote:return client,nil;case <-ctx.Done():client.Close();remote.Close();return nil,ctx.Err()}
    }
    return srv
}
func hiddenSurveyExchange(t *testing.T,u SurveyUplink,name string)([]byte,error){
    t.Helper();client,server:=net.Pipe();defer client.Close()
    deadline:=time.Now().Add(3*time.Second);client.SetDeadline(deadline)
    done:=make(chan error,1);go func(){done<-RelaySurvey(server,u,deadline)}()
    request:=make([]byte,2+len(name));binary.BigEndian.PutUint16(request,uint16(len(name)));copy(request[2:],name)
    // Fragmented writes exercise ReadFull rather than assuming TCP messages.
    for _,part:=range [][]byte{request[:1],request[1:2],request[2:]}{if len(part)>0{if _,err:=client.Write(part);err!=nil{break}}}
    response,readErr:=io.ReadAll(client)
    select{case err:=<-done:if err!=nil{return response,err};return response,readErr;case <-time.After(4*time.Second):t.Fatal("relay did not close");return nil,nil}
}
func TestHiddenSurveyTLSFramingAndRetry(t *testing.T){
    labkit.LeakCheck(t)
    rng:=rand.New(rand.NewSource(33091))
    for _,size:=range []int{0,1,1024}{
        payload:=make([]byte,size);rng.Read(payload)
        name:="probe & +?";var calls atomic.Int64
        srv:=hiddenSurveyTLS(t,http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
            if r.Method!="GET" || r.URL.Query().Get("station")!=name || r.URL.Query().Get("sector")!="ice" || r.Header.Get("Authorization")!="Bearer secret"{t.Error("request contract",r.URL,r.Header)}
            if calls.Add(1)==1{w.WriteHeader(503);return};w.Write(payload)
        }))
        got,err:=hiddenSurveyExchange(t,SurveyUplink{srv.Client(),srv.URL+"/survey?sector=ice","secret"},name)
        if err!=nil || len(got)!=size+2 || int(binary.BigEndian.Uint16(got))!=size || !bytes.Equal(got[2:],payload) || calls.Load()!=2{t.Fatal("frame/retry",len(got),err,calls.Load())}
    }
}
type hiddenSurveyTransport func(*http.Request)(*http.Response,error)
func(f hiddenSurveyTransport)RoundTrip(r *http.Request)(*http.Response,error){return f(r)}
type hiddenSurveyBody struct{io.Reader;closed *int}
func(b hiddenSurveyBody)Close()error{*b.closed++;return nil}
func TestHiddenSurveyRejectionsAndBodyOwnership(t *testing.T){
    labkit.LeakCheck(t)
    for _,tc:=range []struct{status,size,calls int;want error}{
        {503,0,2,ErrSurveyUpstream},{401,0,1,ErrSurveyUpstream},{200,1025,1,ErrSurveyFrame},
    }{
        calls,closed:=0,0
        client:=&http.Client{Transport:hiddenSurveyTransport(func(r *http.Request)(*http.Response,error){calls++;return &http.Response{StatusCode:tc.status,Body:hiddenSurveyBody{strings.NewReader(strings.Repeat("x",tc.size)),&closed},Header:make(http.Header)},nil})}
        got,err:=hiddenSurveyExchange(t,SurveyUplink{client,"https://station.invalid/survey","s"},"beacon")
        if !errors.Is(err,tc.want) || len(got)!=0 || calls!=tc.calls || closed!=calls{t.Fatal("rejection/body ownership",err,len(got),calls,closed)}
    }
    for _,name:=range []string{"",strings.Repeat("a",65)}{
        called:=false;client:=&http.Client{Transport:hiddenSurveyTransport(func(*http.Request)(*http.Response,error){called=true;return nil,errors.New("unexpected")})}
        _,err:=hiddenSurveyExchange(t,SurveyUplink{client,"https://station.invalid","s"},name)
        if !errors.Is(err,ErrSurveyFrame) || called{t.Fatal("bad request reached upstream",err,called)}
    }
    failure:=errors.New("transport broken")
    client:=&http.Client{Transport:hiddenSurveyTransport(func(*http.Request)(*http.Response,error){return nil,failure})}
    _,err:=hiddenSurveyExchange(t,SurveyUplink{client,"https://station.invalid","s"},"beacon")
    if !errors.Is(err,failure){t.Fatal("lost transport cause",err)}
    _,err=hiddenSurveyExchange(t,SurveyUplink{client,"http://station.invalid","s"},"beacon")
    if !errors.Is(err,ErrSurveyUpstream){t.Fatal("plaintext accepted",err)}
}
func TestHiddenBonusSurveyUntrustedTLSAndDeadline(t *testing.T){
    labkit.LeakCheck(t)
    srv:=hiddenSurveyTLS(t,http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.Write([]byte("ok"))}))
    trusted:=srv.Client().Transport.(*http.Transport)
    bad:=trusted.Clone();bad.TLSClientConfig=&tls.Config{MinVersion:tls.VersionTLS12};defer bad.CloseIdleConnections()
    _,err:=hiddenSurveyExchange(t,SurveyUplink{&http.Client{Transport:bad},srv.URL,"s"},"beacon")
    if err==nil{t.Fatal("untrusted certificate accepted")}
    a,b:=net.Pipe();defer a.Close()
    done:=make(chan error,1)
    go func(){done<-RelaySurvey(b,SurveyUplink{srv.Client(),srv.URL,"s"},time.Now().Add(-time.Second))}()
    select{case err=<-done:var ne net.Error;if !errors.As(err,&ne) || !ne.Timeout(){t.Fatal("deadline lost",err)};case <-time.After(3*time.Second):t.Fatal("read deadline not applied")}
}
`},
	})
}
