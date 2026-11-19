package levels

func init() {
	add(&Level{
		ID: "9-boss-evacuation", Chapter: 9, Kind: "lab", Boss: true, Setup: seededLabRoom, Race: false,
		LessonRevision: 1, Steps: ch09BossArabic,
		Title: "Boss: شحنة الإخلاء", Goal: "اجمع شحنة إخلاء باستخدام التصاريح، ومن غير ما ترجّع results جزئية",
		Requires: []string{"map", "struct", "iface", "range"}, Bridge: []int{5, 11, 22},
		Files: []LabFile{{Name: "main.go", Starter: `package main

import "errors"

type Crate struct { Bay string; Units int }
type Clearance interface { Allow(string) bool }
var ErrManifest = errors.New("invalid evacuation manifest")

func Pack(crates []Crate, permit Clearance) (map[string]int, error) {
    return nil, ErrManifest
}
func main() {}
`}},
		Lesson: `# Boss: شحنة الإخلاء

## قبل ما تكتب
- إيه المدخلات اللي لازم تتراجع قبل أي تأثير جانبي؟
- مين مالك النتيجة وكل مورد بعد الخطأ؟

## اختبار التجميع
طوارئ المستعمرة بتجمع أدوات إنت عارفها. مفيش مفهوم جديد. راجع: 2-3-callsigns, 5-3-fault-log, 6-3-census, 7-5-crew-exam, 9-4-small-plugs.

مثال أصغر ومختلف:
~~~go
counts := map[string]int{}
for _, name := range []string{"antenna", "antenna"} { counts[name]++ }
~~~

## أخطاء شائعة
~~~go
if err != nil { return map[string]int{}, nil }
~~~
ده بيضيّع فشل أو بيغيّر الملكية؛ التزم بالعقد المكتوب.

## المهمة
نفّذ ~Pack(crates []Crate, permit Clearance) (map[string]int, error)~. كل ~Crate~ فيها ~Bay~ اسم البوابة و~Units~ عدد الوحدات. ~Clearance.Allow(bay)~ بتحدد السماح؛ مثال نداء: ~totals, err := Pack([]Crate{{"north", 3}}, permit)~. التصريح محقون ومش nil. ارفض الاسم الفاضي، العدد السالب أو بوابة مرفوضة بـ ~nil, ErrManifest~. اجمع التكرار، اقبل الصفر، وارجع map جديدة حتى للمدخل الفاضي. متغيّرش المدخل. الإجمالي بيسع int. سيب ~main~ فاضية؛ مفيش Report.

النجمة التالتة: افحص صلاحية كل الصفوف قبل أول نداء للتصريح، وحافظ على أسماء Unicode زي ما هي.
`,
		Hints: []string{"الـstruct بتجمع البيانات، والـinterface بتحدد السلوك المطلوب بالـmethods.", "افصل بين الـinput والـresult ومسارات الـerror.", "راجع توقيع ~Pack~ وتعريف ~ErrManifest~ في ~main.go~ وانت بتراجع كل مسارات الـreturn."},
		En: &LevelText{LessonRevision: 1, Steps: ch09BossEnglish, Title: "Boss: Evacuation Cargo", Goal: "Assemble cleared evacuation cargo without partial results", Files: []LabFile{{Name: "main.go", Starter: `package main

import "errors"

type Crate struct { Bay string; Units int }
type Clearance interface { Allow(string) bool }
var ErrManifest = errors.New("invalid evacuation manifest")

func Pack(crates []Crate, permit Clearance) (map[string]int, error) {
    return nil, ErrManifest
}
func main() {}
`}},
			Lesson: `# Boss: Evacuation Cargo

## Before you write
- Which inputs must be checked before any side effect?
- Who owns the result and each resource after an error?

## Integration checkpoint
A colony emergency combines tools you already know. No new concept is needed. Revisit: 2-3-callsigns, 5-3-fault-log, 6-3-census, 7-5-crew-exam, 9-4-small-plugs.

A smaller, different example:
~~~go
counts := map[string]int{}
for _, name := range []string{"antenna", "antenna"} { counts[name]++ }
~~~

## Common mistakes
~~~go
if err != nil { return map[string]int{}, nil }
~~~
This drops a failure or changes ownership; preserve the documented contract.

## Mission
Implement ~Pack(crates []Crate, permit Clearance) (map[string]int, error)~. Each ~Crate~ has a bay identifier ~Bay~ and count ~Units~. ~Clearance.Allow(bay)~ decides permission; for example ~totals, err := Pack([]Crate{{"north", 3}}, permit)~. The injected permit is nonnil. Reject empty names, negative counts, or denied bays with ~nil, ErrManifest~. Sum repeated bays, accept zero, and return a fresh nonnil map even for empty input. Do not mutate input. Totals fit int. Leave ~main~ empty; no Report.

Third star: validate every row before the first clearance call, preserving Unicode bay names exactly.
`,
			Hints: []string{"A struct groups data, while an interface describes required method behavior.", "Keep input, result, and error paths clear as separate responsibilities.", "Check ~Pack~ and ~ErrManifest~ in ~main.go~ while reviewing each return path."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: `package main
import (
    "errors"
    "fmt"
    "math/rand"
    "reflect"
    "testing"
)
type hiddenClearance struct { denied string; calls int }
func (p *hiddenClearance) Allow(bay string) bool { p.calls++; return bay != p.denied }
func TestHiddenEvacuation(t *testing.T) {
    // Empty, zero, repeated, Unicode, rejected, negative, and empty-name boundaries.
    for _, seed := range []int64{901,902,903,904,905} {
        rng := rand.New(rand.NewSource(seed))
        input := []Crate{{"east",0}}
        want := map[string]int{"east":0}
        for i:=0; i<40; i++ {
            c := Crate{fmt.Sprintf("bay-%d",rng.Intn(6)),rng.Intn(50)}
            input=append(input,c); want[c.Bay]+=c.Units
        }
        before := append([]Crate(nil),input...)
        got,err := Pack(input,&hiddenClearance{})
        if err!=nil || !reflect.DeepEqual(got,want) { t.Fatalf("seed %d: %v %v",seed,got,err) }
        if !reflect.DeepEqual(input,before) { t.Fatal("mutated crates") }
        got["east"]=99
        again,err:=Pack(input,&hiddenClearance{})
        if err!=nil || again["east"]!=0 { t.Fatal("result aliases later call") }
    }
    for _, input := range [][]Crate{nil,{},{{"dock",0}}} {
        got,err:=Pack(input,&hiddenClearance{})
        if err!=nil || got==nil { t.Fatal("empty/zero must succeed",got,err) }
    }
    for _, input := range [][]Crate{{{"",1}},{{"east",-1}},{{"east",4},{"closed",2}}} {
        got,err:=Pack(input,&hiddenClearance{denied:"closed"})
        if got!=nil || !errors.Is(err,ErrManifest) { t.Fatal("rejection must return nil and ErrManifest",got,err) }
    }
}
func TestHiddenBonusValidateBeforeClearance(t *testing.T) {
    p:=&hiddenClearance{}
    got,err:=Pack([]Crate{{"valid",2},{"",3}},p)
    if got!=nil || !errors.Is(err,ErrManifest) || p.calls!=0 { t.Fatal("invalid manifest consulted clearance") }
    got,err=Pack([]Crate{{"\u0628\u0627\u0628-🚀",3},{"\u0628\u0627\u0628-🚀",4}},p)
    if err!=nil || got["\u0628\u0627\u0628-🚀"]!=7 { t.Fatal("bay identifiers changed",got,err) }
}
`},
	})
}
