package levels

func init() {
	add(&Level{
		ID: "23-boss-evidence", Chapter: 23, Kind: "lab", Boss: true, LessonRevision: 1, Steps: ch23LessonBossArabic, Setup: seededLabRoom, Race: false,
		Title: "Boss: دليل العاصفة", Goal: "حوّل قراءات CSV لملف JSON مختوم بحفظ ذرّي",
		Requires: []string{"map", "sel:csv.NewReader", "sel:json.Marshal", "sel:cipher.NewGCM"}, Bridge: []int{17, 18, 25, 30},
		Files: []LabFile{{Name: "main.go", Starter: `package main
import ("io";"time";"errors")
type EvidenceSample struct { Sensor string; Value int }
type EvidenceBundle struct { Captured string; Samples []EvidenceSample }
var ErrEvidence=errors.New("invalid evidence")
func SealEvidence(in io.Reader, dst string, key []byte, entropy io.Reader, at time.Time) error {
    return ErrEvidence
}
func main() {}
`}},
		Lesson: `# Boss: دليل العاصفة

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 18-2-black-box, 18-3-cargo-manifest, 19-1-durations, 22-3-atomic-save, 23-2-sealed-cargo.

مثال أصغر ومختلف:
~~~go
rows := []string{"beta", "alpha"}
sort.Strings(rows) // deterministic report order
~~~

## أخطاء شائعة
~~~go
os.WriteFile(dst, partial, 0644) // destroys the previous report before validation
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~SealEvidence(in, dst, key, entropy, at) error~ بتوقيع الـ starter. الـ readers مش nil، ومجلد ~dst~ موجود. اقرأ CSV بهيدر بالظبط ~sensor,value~ وصفوف من عمودين، اسم مش فاضي وعدد int غير سالب؛ لا تعمل trim. اجمع القراءات لنفس الاسم ورتّب الأسماء تصاعدي بالبايتات. الإجماليات بتسع int. الهيدر لوحده صحيح، والفاضي أو CSV بايظ خطأ. الاسم/العدد/الهيدر غير الصحيح يطابق ~ErrEvidence~، وأخطاء I/O تحتفظ بالسبب.

حوّل ~EvidenceBundle{Captured: at.UTC().Format(time.RFC3339), Samples: samples}~ لـ JSON؛ ~samples~ فاضية مش nil لو مفيش صفوف. مثال صف: ~EvidenceSample{Sensor: "ice", Value: 7}~. شفّر بـ AES-GCM بمفتاح AES صالح، nonce من ~io.ReadFull(entropy, nonce)~، وAAD بالظبط ~[]byte("colony:evidence:v1")~. الملف بايتات ~nonce || ciphertext~. احفظ بملف مؤقت في نفس المجلد ثم Close وRename؛ أي فشل يسيب الملف القديم زي ما هو وينضّف المؤقت. سيب ~main~ فاضية. النجمة التالتة: الملف النهائي صلاحياته 0600 حتى لو استبدل ملف 0644.
`,
		Hints: []string{"حافظ على المدخلات سليمة قبل أي أثر جانبي.", "قسّم الشغل لمراحل وتأكد من ناتج كل مرحلة.", "راجع ~csv.Reader.Read~ وحدود الموارد عند فشل callback."},
		En: &LevelText{LessonRevision: 1, Steps: ch23LessonBossEnglish, Title: "Boss: Storm Evidence", Goal: "Turn CSV readings into an authenticated JSON file with atomic replacement", Files: []LabFile{{Name: "main.go", Starter: `package main
import ("io";"time";"errors")
type EvidenceSample struct { Sensor string; Value int }
type EvidenceBundle struct { Captured string; Samples []EvidenceSample }
var ErrEvidence=errors.New("invalid evidence")
func SealEvidence(in io.Reader, dst string, key []byte, entropy io.Reader, at time.Time) error {
    return ErrEvidence
}
func main() {}
`}},
			Lesson: `# Boss: Storm Evidence

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 18-2-black-box, 18-3-cargo-manifest, 19-1-durations, 22-3-atomic-save, 23-2-sealed-cargo.

A smaller, different example:
~~~go
rows := []string{"beta", "alpha"}
sort.Strings(rows) // deterministic report order
~~~

## Common mistakes
~~~go
os.WriteFile(dst, partial, 0644) // destroys the previous report before validation
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~SealEvidence(in, dst, key, entropy, at) error~ with the starter signature. Readers are nonnil; the destination parent directory exists. Read CSV with exact header ~sensor,value~ and two-column rows: a nonempty name and nonnegative int, without trimming. Aggregate repeated names and sort names by ascending byte order. Totals fit int. Header-only input succeeds; empty or malformed CSV fails. Invalid header/name/number must match ~ErrEvidence~; preserve I/O error causes.

Marshal ~EvidenceBundle{Captured: at.UTC().Format(time.RFC3339), Samples: samples}~ as JSON; samples is nonnil and empty for no rows. Example row: ~EvidenceSample{Sensor: "ice", Value: 7}~. Encrypt using AES-GCM with a valid AES key, a nonce filled by ~io.ReadFull(entropy, nonce)~, and exact AAD ~[]byte("colony:evidence:v1")~. File bytes are ~nonce || ciphertext~. Write a temporary file in the same directory, Close, then Rename. Every failure preserves the old destination and removes temporary files. Leave ~main~ empty. Third star: the final file has permissions 0600 even when replacing a 0644 file.
`,
			Hints: []string{"Validate input before making side effects.", "Split the work into stages and check each result.", "Review ~csv.Reader.Read~ and resource lifetimes when a callback fails."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: `package main
import (
    "bytes"
    "crypto/aes"
    "crypto/cipher"
    "encoding/csv"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "math/rand"
    "os"
    "path/filepath"
    "reflect"
    "sort"
    "strings"
    "testing"
    "time"
)
func hiddenEvidenceOpen(t *testing.T,path string,key []byte) EvidenceBundle {
    t.Helper()
    raw,err:=os.ReadFile(path);if err!=nil {t.Fatal(err)}
    block,err:=aes.NewCipher(key);if err!=nil {t.Fatal(err)}
    gcm,err:=cipher.NewGCM(block);if err!=nil {t.Fatal(err)}
    if len(raw)<gcm.NonceSize()+gcm.Overhead(){t.Fatal("missing authenticated payload")}
    plain,err:=gcm.Open(nil,raw[:gcm.NonceSize()],raw[gcm.NonceSize():],[]byte("colony:evidence:v1"))
    if err!=nil {t.Fatal("not the documented GCM envelope",err)}
    var bundle EvidenceBundle
    if err=json.Unmarshal(plain,&bundle);err!=nil {t.Fatal(err)}
    raw[len(raw)-1]^=1
    if _,err=gcm.Open(nil,raw[:gcm.NonceSize()],raw[gcm.NonceSize():],[]byte("colony:evidence:v1"));err==nil {t.Fatal("tampering accepted")}
    return bundle
}
func TestHiddenSealedEvidence(t *testing.T) {
    rng:=rand.New(rand.NewSource(23091))
    at:=time.Date(2041,2,3,1,2,3,0,time.FixedZone("station",7200))
    for _,n:=range []int{0,1,50} {
        var input bytes.Buffer
        writer:=csv.NewWriter(&input);writer.Write([]string{"sensor","value"})
        totals:=map[string]int{}
        for i:=0;i<n;i++ {name:=fmt.Sprintf("probe,%d",rng.Intn(7));v:=rng.Intn(100);totals[name]+=v;writer.Write([]string{name,fmt.Sprint(v)})}
        writer.Flush();if err:=writer.Error();err!=nil {t.Fatal(err)}
        key:=make([]byte,32);rng.Read(key)
        nonce:=make([]byte,12);rng.Read(nonce)
        path:=filepath.Join(t.TempDir(),"evidence.bin")
        if err:=SealEvidence(&input,path,key,bytes.NewReader(nonce),at);err!=nil {t.Fatal(err)}
        got:=hiddenEvidenceOpen(t,path,key)
        want:=make([]EvidenceSample,0,len(totals))
        for k,v:=range totals {want=append(want,EvidenceSample{k,v})}
        sort.Slice(want,func(i,j int)bool{return want[i].Sensor<want[j].Sensor})
        if got.Captured!=at.UTC().Format(time.RFC3339) || !reflect.DeepEqual(got.Samples,want){t.Fatal("report contract",got,want)}
    }
}
type hiddenEvidenceReader struct {err error}
func (r hiddenEvidenceReader) Read([]byte)(int,error){return 0,r.err}
func TestHiddenEvidenceFailurePreservesDestination(t *testing.T) {
    failure:=errors.New("sensor link lost")
    for _,tc:=range []struct{name string;in io.Reader; key []byte; entropy io.Reader; cause error}{
        {"empty",strings.NewReader(""),make([]byte,32),bytes.NewReader(make([]byte,12)),nil},
        {"header",strings.NewReader("wrong,value\na,1\n"),make([]byte,32),bytes.NewReader(make([]byte,12)),ErrEvidence},
        {"negative",strings.NewReader("sensor,value\na,-1\n"),make([]byte,32),bytes.NewReader(make([]byte,12)),ErrEvidence},
        {"name",strings.NewReader("sensor,value\n,1\n"),make([]byte,32),bytes.NewReader(make([]byte,12)),ErrEvidence},
        {"number",strings.NewReader("sensor,value\na,no\n"),make([]byte,32),bytes.NewReader(make([]byte,12)),ErrEvidence},
        {"row",strings.NewReader("sensor,value\na,1,extra\n"),make([]byte,32),bytes.NewReader(make([]byte,12)),nil},
        {"reader",hiddenEvidenceReader{failure},make([]byte,32),bytes.NewReader(make([]byte,12)),failure},
        {"key",strings.NewReader("sensor,value\na,1\n"),[]byte{1},bytes.NewReader(make([]byte,12)),nil},
        {"entropy",strings.NewReader("sensor,value\na,1\n"),make([]byte,32),hiddenEvidenceReader{failure},failure},
    } {t.Run(tc.name,func(t *testing.T){
        dir:=t.TempDir();path:=filepath.Join(dir,"evidence.bin")
        if err:=os.WriteFile(path,[]byte("previous"),0600);err!=nil{t.Fatal(err)}
        err:=SealEvidence(tc.in,path,tc.key,tc.entropy,time.Unix(0,0))
        if err==nil || tc.cause!=nil && !errors.Is(err,tc.cause){t.Fatal("lost failure",err)}
        got,e:=os.ReadFile(path);if e!=nil || string(got)!="previous"{t.Fatal("destination damaged",string(got),e)}
        entries,e:=os.ReadDir(dir);if e!=nil || len(entries)!=1{t.Fatal("temporary file left",entries,e)}
    })}
    // A directory destination makes rename fail after a successful write.
    dir:=t.TempDir();dst:=filepath.Join(dir,"occupied");if err:=os.Mkdir(dst,0700);err!=nil{t.Fatal(err)}
    if err:=SealEvidence(strings.NewReader("sensor,value\n"),dst,make([]byte,32),bytes.NewReader(make([]byte,12)),time.Unix(0,0));err==nil{t.Fatal("rename error hidden")}
    entries,err:=os.ReadDir(dir);if err!=nil || len(entries)!=1{t.Fatal("rename failure leaked temp",entries,err)}
}
func TestHiddenBonusEvidencePrivateReplacement(t *testing.T) {
    path:=filepath.Join(t.TempDir(),"sealed")
    if err:=os.WriteFile(path,[]byte("public old report"),0644);err!=nil{t.Fatal(err)}
    if err:=SealEvidence(strings.NewReader("sensor,value\nprobe,0\n"),path,make([]byte,16),bytes.NewReader(make([]byte,12)),time.Unix(0,0));err!=nil{t.Fatal(err)}
    info,err:=os.Stat(path);if err!=nil{t.Fatal(err)}
    if info.Mode().Perm()!=0600{t.Fatal("sealed evidence is not private",info.Mode())}
}
`},
	})
}
