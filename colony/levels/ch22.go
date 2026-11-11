package levels

func init() {
	add(&Level{ID: "22-1-tree", Chapter: 22, Kind: "lab", Title: "شجرة المخزن", Goal: "اعرض مخزن المستعمرة بشجرة مرتبة وعمق محدود",
		Requires: []string{"recursion", "sel:fs.FS"}, Bridge: []int{24}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch221StarterAr}},
		Lesson: `# خريطة المخزن

الطاقم محتاج يشوف رفوف المخزن من غير فتح كل صندوق. هنقرأ شجرة افتراضية ونكتب للـ writer اللي المنادي يختاره.

## قبل ما تكتب
- إزاي عمق 0 يختلف عن مخزن فاضي؟
- هتوقف النزول عند الحد من غير محاولة قراءة مجلد أعمق إزاي؟

## نظام ملفات محقون
استورد ~io/fs~. ~fs.FS~ واجهة فيها ~Open~، مش مسار على جهازك. الأسماء جواها نسبية بشرطة مائلة، والجذر اسمه ~"."~. ~fs.Stat~ بترجع معلومات، و ~IsDir~ بتحدد مجلد.
~~~go
func isShelf(store fs.FS) (bool,error) {
 info,err:=fs.Stat(store,".")
 if err!=nil {return false,err}
 return info.IsDir(),nil
}
~~~
استورد ~testing/fstest~ في التست. ~fstest.MapFS~ خريطة أسماء لـ ~fstest.MapFile~، والمجلدات الأب بتتولد تلقائياً. ~Data~ محتوى ملف؛ ~Mode~ نوعه وصلاحياته.
~~~go
store:=fstest.MapFS{"bay/cell":&fstest.MapFile{Data:[]byte("7")},"empty":&fstest.MapFile{Mode:fs.ModeDir|0700}}
entries,err:=fs.ReadDir(store,".")
fmt.Println(len(entries),err) // 2 <nil>
~~~
~fs.ReadDir~ ترجع ~[]fs.DirEntry~ مرتبة بالاسم. ~Name~ الاسم المحلي، و ~IsDir~ هل هو مجلد. استورد ~path~، واستخدم ~path.Join~ لأسماء الواجهة بدل مسارات نظام التشغيل.
~~~go
for _,entry:=range entries {
 name:=path.Join("bay",entry.Name())
 fmt.Println(name,entry.IsDir())
}
~~~

## recursion وحد عمق
زي شجرة الأرشيف، الدالة المساعدة تنادي نفسها للمجلدات بس. افحص الحد قبل ~fs.ReadDir~؛ قراءة فرع مخفي ممكن تفشل بلا داعي. مثال عد تنازلي مختلف:
~~~go
func steps(n int) int {
 if n==0 {return 0}
 return 1+steps(n-1)
}
~~~
~strings.Repeat~ تكرر نص بعدد معين؛ هتعمل مسافتين لكل مستوى. ~fmt.Fprintf~ تكتب للـ writer وتدي خطأ لازم يرجع فوراً.
~~~go
indent:=strings.Repeat("  ",2)
_,err:=fmt.Fprintf(out,"%s%s\n",indent,"cell")
~~~

## النجمة التالتة: نوع الملف
~fs.FileMode~ مجموعة bits. ~0700~ صلاحيات المالك للمجلد. ~fs.ModeDir~ علامة مجلد؛ ~|~ تجمع bits، و ~&~ تختبر علامة. ~entry.Type()~ ترجع bits النوع بس.
~~~go
if entry.Type()&fs.ModeSymlink!=0 {fmt.Println(entry.Name()+"@")}
~~~
الرابط الرمزي اسم بيشير لمكان تاني. اعرضه بـ ~@~ من غير ما تتبعه؛ الأنواع الخاصة التانية بـ ~!~. الملف العادي نوعه 0؛ متفتحش محتوى الملفات عشان تعرض أسماءها.

## نداء الدالة
~~~go
err:=StoreTree(store,2,out)
~~~
~StoreTree~ بتعرض لحد مستويين تحت الجذر؛ افحص الخطأ قبل الاعتماد على الخريطة.

## أخطاء شائعة
~~~go
entries,_:=fs.ReadDir(store,".") // خطأ القراية اتحول لشجرة ناقصة
StoreTree(store,depth,out) // استدعاء نفس المدخلات مش recursion بتتقدم
~~~

## المهمة
نفّذ ~StoreTree(store fs.FS, depth int, out io.Writer) error~. المدخلان غير nil. ~depth~ سالب خطأ قبل أي I/O. تحقق إن الجذر ~"."~ موجود ومجلد، ثم اكتب ~store/~ وسطر جديد. العمق 0 يكتب الجذر بس؛ 1 يضيف أبناءه؛ كل مستوى زيادة يسمح بأبناء أعمق.

اعرض depth-first بترتيب أسماء الإخوة في Go. كل سطر مسافتين لكل مستوى، ثم الاسم المحلي؛ المجلد آخره ~/~. كل سطر ينتهي newline. متقراش مجلد عند حد العمق، ومتغيرش المصدر. استخدم recursion و ~fs.FS~. فشل filesystem أو writer يوقف فوراً بخطأ يحتفظ بالسبب وباسم المسار المتأثر؛ الإخراج السابق ممكن يفضل.

النجمة التالتة: الروابط آخرها ~@~ والملفات الخاصة ~!~، من غير نزول أو فتح. مفيش Report ولا اعتماد على قرص جهازك.
`, Hints: []string{"تحقق من الجذر، وبعدها اكتب store/ حتى لو مفيش أبناء.", "استخدم fs.ReadDir عشان ترتيب الإخوة؛ path.Join لبناء مسار الطفل.", "مرر المستوى للدالة recursive وافحص حد العمق قبل القراءة؛ افحص خطأ كل كتابة."},
		En: &LevelText{Title: "Store Tree", Goal: "Render a sorted, depth-limited colony store", Files: []LabFile{{Name: "main.go", Starter: ch221Starter}}, Lesson: `# A Map of the Store

The crew needs a shelf map without opening every crate. Read an injected filesystem and write to the caller's writer.

## Before you write
- How does depth 0 differ from an empty store?
- How will you stop at the depth boundary without reading a deeper directory?

## An injected filesystem
Import ~io/fs~. ~fs.FS~ is an interface with ~Open~, rather than a host path. Its names are relative slash-separated paths; the root is ~"."~. ~fs.Stat~ returns metadata whose ~IsDir~ method identifies directories.
~~~go
func isShelf(store fs.FS) (bool,error) {
 info,err:=fs.Stat(store,".")
 if err!=nil {return false,err}
 return info.IsDir(),nil
}
~~~
Import ~testing/fstest~ in tests. ~fstest.MapFS~ maps names to ~fstest.MapFile~ pointers and synthesizes parent directories. ~Data~ contains bytes; ~Mode~ specifies type and permissions.
~~~go
store:=fstest.MapFS{"bay/cell":&fstest.MapFile{Data:[]byte("7")},"empty":&fstest.MapFile{Mode:fs.ModeDir|0700}}
entries,err:=fs.ReadDir(store,".")
fmt.Println(len(entries),err) // 2 <nil>
~~~
~fs.ReadDir~ returns ~[]fs.DirEntry~ sorted by name. ~Name~ is the local name; ~IsDir~ identifies directories. Import ~path~ and use ~path.Join~ for interface names rather than operating-system paths.
~~~go
for _,entry:=range entries {
 name:=path.Join("bay",entry.Name())
 fmt.Println(name,entry.IsDir())
}
~~~

## Recursion and depth
As with the archive tree, a helper calls itself for directories only. Check the limit before ~fs.ReadDir~: reading a hidden branch can fail unnecessarily. A different countdown example:
~~~go
func steps(n int) int {
 if n==0 {return 0}
 return 1+steps(n-1)
}
~~~
~strings.Repeat~ repeats text a specified number of times; use two spaces per level. ~fmt.Fprintf~ writes to the injected writer and returns errors that must stop rendering.
~~~go
indent:=strings.Repeat("  ",2)
_,err:=fmt.Fprintf(out,"%s%s\n",indent,"cell")
~~~

## Third star: file types
~fs.FileMode~ is a bit set. ~0700~ grants owner directory permissions. ~fs.ModeDir~ marks directories; ~|~ combines bits and ~&~ tests them. ~entry.Type()~ returns only type bits.
~~~go
if entry.Type()&fs.ModeSymlink!=0 {fmt.Println(entry.Name()+"@")}
~~~
A symbolic link points elsewhere. Mark it with ~@~ without following it. Mark other special types with ~!~. Regular files have type 0; do not open their contents to display names.

## Calling the function
~~~go
err:=StoreTree(store,2,out)
~~~
~StoreTree~ shows two levels below the root; check its error before relying on the map.

## Common mistakes
~~~go
entries,_:=fs.ReadDir(store,".") // Turns read failure into an incomplete tree
StoreTree(store,depth,out) // Identical arguments make no recursive progress
~~~

## Your mission
Implement ~StoreTree(store fs.FS, depth int, out io.Writer) error~. Both injected values are non-nil. Reject negative depth before I/O. Validate that root ~"."~ exists and is a directory, then write ~store/~ and a newline. Depth 0 prints only the root; 1 includes direct children; further levels include deeper children.

Use depth-first traversal with sibling names in Go string order. Each line has two spaces per level, the local name, and ~/~ for directories. End every line with a newline. Never read a directory at the depth boundary or mutate the source. Use recursion and ~fs.FS~. Stop immediately on filesystem or writer errors, preserving the cause and affected path; prior output may remain.

Third star: mark links with ~@~ and other special files with ~!~, without opening or descending into them. No Report or host disk dependency.
`, Hints: []string{"Validate the root, then print store/ even when it has no children.", "Use fs.ReadDir for sorted siblings and path.Join for child paths.", "Pass the level recursively; check the depth before reading and every write error."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch221Hidden, Bonus: true},
	})
}

const ch221Starter = `package main

import (
	"errors"
	"io"
	"io/fs"
)

func StoreTree(store fs.FS, depth int, out io.Writer) error { return errors.New("unfinished") }
func main()                                                 {}
`
const ch221StarterAr = `package main

import (
	"errors"
	"io"
	"io/fs"
)

func StoreTree(store fs.FS, depth int, out io.Writer) error { return errors.New("لسه مخلصتش") }
func main()                                                 {}
`
const ch221Hidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"strings"
	"testing"
	"testing/fstest"
)

func TestHiddenTreeSeededDepthAndOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(221))
	for round := 0; round < 25; round++ {
		store := fstest.MapFS{"empty": &fstest.MapFile{Mode: fs.ModeDir | 0700}}
		names := []string{"A", "z", "قمر"}
		counts := make([]int, len(names))
		for i, name := range names {
			counts[i] = 1 + rng.Intn(8)
			for j := 0; j < counts[i]; j++ {
				store[fmt.Sprintf("%s/c%02d", name, j)] = &fstest.MapFile{Data: []byte{byte(rng.Intn(256))}}
			}
		}
		for depth := 0; depth <= 4; depth++ {
			want := "store/\n"
			if depth > 0 {
				for i, name := range names {
					if name == "z" {
						want += "  empty/\n"
					}
					want += "  " + name + "/\n"
					if depth > 1 {
						for j := 0; j < counts[i]; j++ {
							want += fmt.Sprintf("    c%02d\n", j)
						}
					}
				}
			}
			var out bytes.Buffer
			if err := StoreTree(store, depth, &out); err != nil || out.String() != want {
				t.Fatalf("depth %d: %q want %q err %v", depth, out.String(), want, err)
			}
		}
	}
	var out bytes.Buffer
	if err := StoreTree(fstest.MapFS{}, 3, &out); err != nil || out.String() != "store/\n" {
		t.Fatal("empty", out.String(), err)
	}
}

func TestHiddenTreeNestedBoundaryAndFileRoot(t *testing.T) {
	store := fstest.MapFS{"a/b/c/d": &fstest.MapFile{}}
	lines := []string{"store/", "  a/", "    b/", "      c/", "        d"}
	for depth := 0; depth <= 5; depth++ {
		var out bytes.Buffer
		count := depth + 1
		if count > len(lines) {
			count = len(lines)
		}
		want := strings.Join(lines[:count], "\n") + "\n"
		if err := StoreTree(store, depth, &out); err != nil || out.String() != want {
			t.Fatal("nested depth", depth, out.String(), err)
		}
	}
	file, err := fs.Stat(store, "a/b/c/d")
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := StoreTree(hiddenTreeFileRoot{store, file}, 0, &out); err == nil || out.Len() != 0 {
		t.Fatal("file root accepted", err)
	}
}

type hiddenTreeFileRoot struct {
	fs.FS
	file fs.FileInfo
}

func (s hiddenTreeFileRoot) Stat(string) (fs.FileInfo, error) { return s.file, nil }

type hiddenTreeFS struct {
	fs.FS
	fail  string
	cause error
}

func (s hiddenTreeFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if name == s.fail {
		return nil, s.cause
	}
	return fs.ReadDir(s.FS, name)
}
func TestHiddenTreeErrorsAndPruning(t *testing.T) {
	wire := errors.New("shelf offline")
	base := fstest.MapFS{"bay/cell": &fstest.MapFile{}}
	for _, depth := range []int{0, 1, 2} {
		var out bytes.Buffer
		err := StoreTree(hiddenTreeFS{base, "bay", wire}, depth, &out)
		if depth < 2 {
			if err != nil {
				t.Fatal("read beyond limit", err)
			}
		} else if !errors.Is(err, wire) || !strings.Contains(err.Error(), "bay") {
			t.Fatal("lost read error", err)
		}
	}
	var out bytes.Buffer
	if err := StoreTree(base, -1, &out); err == nil || out.Len() != 0 {
		t.Fatal("negative depth")
	}
	if err := StoreTree(hiddenTreeFS{base, ".", wire}, 1, &out); !errors.Is(err, wire) {
		t.Fatal("root read", err)
	}
	if err := StoreTree(base, 3, hiddenTreeWriter{wire}); !errors.Is(err, wire) {
		t.Fatal("write cause", err)
	}
	if err := StoreTree(hiddenTreeRoot{}, 0, io.Discard); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("root stat", err)
	}
}

type hiddenTreeWriter struct{ err error }

func (w hiddenTreeWriter) Write(p []byte) (int, error) { return 0, w.err }

type hiddenTreeRoot struct{}

func (hiddenTreeRoot) Open(name string) (fs.File, error) { return nil, fs.ErrNotExist }
func TestHiddenBonusTreeSpecialEntries(t *testing.T) {
	store := fstest.MapFS{"link": &fstest.MapFile{Mode: fs.ModeSymlink, Data: []byte(".")}, "pipe": &fstest.MapFile{Mode: fs.ModeNamedPipe}}
	var out bytes.Buffer
	if err := StoreTree(store, 100, &out); err != nil || out.String() != "store/\n  link@\n  pipe!\n" {
		t.Fatal("special entries", out.String(), err)
	}
}
`

func init() {
	add(&Level{ID: "22-2-organizer", Chapter: 22, Kind: "lab", Title: "المنظّم", Goal: "خطط توزيع صناديق الوارد قبل نقل أي ملف",
		Requires: []string{"sel:filepath.WalkDir"}, Bridge: []int{20}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch222StarterAr}},
		Lesson: `# الوارد يروح للرف الصح

صناديق الوارد في ~inbox~، والأرشيف في ~records~ و ~images~ و ~misc~. الخطة بيانات قابلة للمراجعة؛ مفيش نقل أثناء اكتشاف الملفات.

## قبل ما تكتب
- إزاي تمنع أول نقل لو آخر صندوق هيعمل تصادم؟
- لو النقل التالت فشل، هتبلّغ عدد النقلات المكتملة إزاي؟

## المشي على واجهة ملفات
~fs.WalkDir~ بتزور الجذر وبعدين أبناءه بترتيب lexical. الـ callback تاخد المسار و ~fs.DirEntry~ وخطأ. افحص الخطأ الأول؛ المدخل ممكن nil وقت الفشل. رجوع خطأ يوقف المشي.
~~~go
err:=fs.WalkDir(store,"inbox",func(name string,entry fs.DirEntry,err error)error{
 if err!=nil {return err}
 if entry.Type().IsRegular() {fmt.Println(name)}
 return nil
})
~~~
~Type().IsRegular()~ تميّز الملفات العادية عن المجلدات والروابط والأنواع الخاصة؛ ~WalkDir~ مبتتبعش الروابط.

## مسارات الواجهة ومسارات الجهاز
استورد ~path/filepath~ و ~os~. ~filepath.WalkDir~ نفس callback بس بمسارات الجهاز. ~filepath.Join~ تجمع مسارات الجهاز؛ ~filepath.Rel~ تحسب المسار النسبي؛ ~filepath.ToSlash~ تحوله لأسماء ~fs.FS~. ~os.DirFS~ تعرض مجلد حقيقي كواجهة قراية، ومش sandbox ضد روابط خبيثة.
~~~go
err:=filepath.WalkDir(filepath.Join(root,"inbox"),func(name string,entry fs.DirEntry,err error)error{
 if err!=nil {return err}
 relative,err:=filepath.Rel(root,name)
 if err!=nil {return err}
 fmt.Println(filepath.ToSlash(relative))
 return nil
})
store:=os.DirFS(root)
~~~
~filepath.FromSlash~ تعمل العكس عند تحويل اسم الخطة لمسار جهاز. ~slices.Sort~ ترتب slice من strings مباشرة.
~~~go
host:=filepath.Join(root,filepath.FromSlash("records/cell.json"))
names:=[]string{"b","a"}
slices.Sort(names)
fmt.Println(host,names)
~~~

التست يستخدم ~t.TempDir()~ عشان مجلد مستقل بيتشال تلقائياً. أي ملف حقيقي في الفصل يبقى جواه.
~~~go
root:=t.TempDir()
err:=os.MkdirAll(filepath.Join(root,"inbox"),0700)
~~~
~os.MkdirAll~ تنشئ المجلد وآباءه لو ناقصين. ~0700~ صلاحيات المالك بس.

## تصنيف ثم تحقق ثم تنفيذ
~path.Ext~ ترجع الامتداد بالنقطة، و ~path.Base~ آخر اسم. طبق ~strings.ToLower~ على الامتداد بس؛ متغيرش اسم الصندوق.
~~~go
fmt.Println(path.Ext("bay/NOTE.JSON"),path.Base("bay/NOTE.JSON")) // .JSON NOTE.JSON
~~~
اجمع الوجهات في map عشان تمنع مصدرين لنفس الوجهة. افحص الوجهة بـ ~fs.Stat~؛ ~errors.Is(err,fs.ErrNotExist)~ بس معناها المكان فاضي، مش أي خطأ. افحص كمان مجلد الوجهة لو موجود؛ مينفعش يكون ملف. لو stat رجعت غايب، افحص أسماء ~fs.ReadDir~ في الأب الموجود كمان: رابط مقطوع لسه اسم محجوز.
~~~go
_,err:=fs.Stat(store,"records/cycle.json")
if err!=nil && !errors.Is(err,fs.ErrNotExist) {return err}
~~~
~os.Rename~ تنقل الاسم لنفس نظام الملفات؛ فشلها يرجع خطأ، وممكن تستبدل وجهة موجودة على بعض الأنظمة. عشان كده التحقق الكامل يسبق التنفيذ، والعقد يفترض شجرة ثابتة من غير كاتب خارجي أو روابط في مجلدات الوجهة.
~~~go
err:=os.Rename(filepath.Join(root,"inbox","cell.json"),filepath.Join(root,"records","cell.json"))
~~~
~ExecuteShelves~ بتاخد دالة النقل المحقونة عشان نجرب فشل في النص؛ ~OrganizeShelves~ غلاف بيمرر ~os.Rename~.
~~~go
plan,err:=PlanShelves(store)
example:=ShelfMove{From:"inbox/cell.csv",To:"records/cell.csv"}
n,err:=OrganizeShelves(root)
fmt.Println(plan,example,n,err)
n,err=ExecuteShelves(root,os.Rename)
~~~

مفيش نسخ احتياطي تلقائي ولا rollback. احتفظ بعدد النقلات الناجحة وارجع أول خطأ.

## أخطاء شائعة
~~~go
os.MkdirAll(destination,0700) // غلط في dry-run: غيّرت القرص
if err!=nil { /* افتراض إن الوجهة مش موجودة: ممكن permission error */ }
~~~

## المهمة
نفّذ ~PlanShelves(store fs.FS) ([]ShelfMove,error)~ و ~OrganizeShelves(root string) (int,error)~. ~ShelfMove~ فيها ~From,To string~ بأسماء نسبية slash-separated. ~inbox~ لازم موجود ومجلد؛ ممكن فاضي. امشِ داخله recursively، واتجاهل الروابط والأنواع الخاصة. ~.json~ و ~.csv~ يروحوا ~records~؛ ~.png~ يروح ~images~؛ الباقي، بما فيه بدون امتداد، يروح ~misc~. الامتداد case-insensitive؛ الاسم الأصلي محفوظ، والمجلدات الداخلية بتتسطح. رتّب الخطة بـ ~From~ تصاعدي.

افشل بخطة nil لو مصدرين لنفس الوجهة، أو الوجهة موجودة بأي نوع، أو مجلد الهدف موجود ومش مجلد، أو حصل خطأ مشي أو stat. احتفظ بأسباب أخطاء I/O. ~PlanShelves~ للـ dry-run: مفيش إنشاء أو نقل أو كتابة.

~OrganizeShelves~ غلاف يستدعي ~ExecuteShelves(root string, rename func(string,string)error) (int,error)~ مع ~os.Rename~. دالة النقل غير nil؛ في النجاح تنقل فعلاً، وفي الفشل متغيرش الملفات. ~ExecuteShelves~ تستخدم ~filepath.WalkDir~ لاكتشاف ملفات الوارد الحقيقي، وتتحقق من الخطة كلها قبل أي تغيير. بعدها تنشئ مجلدات الهدف المطلوبة بس بـ 0700 وتنقل بالدالة المحقونة بترتيب الخطة. ترجع عدد النقلات المكتملة؛ أول فشل يوقف من غير rollback للنقلات السابقة. إعادة التشغيل بعد النجاح خطة فاضية. الجذر غير متغير أثناء العملية، ومجلدات الوجهة الموجودة حقيقية من غير روابط؛ مفيش تزامن أو Report. التستات الحقيقية داخل ~t.TempDir()~ بس.
`, Hints: []string{"اجمع أسماء الملفات العادية الأول، وبعدها حوّلها لخطة.", "استخدم map للوجهات وافحص stat لكل وجهة ومجلدها قبل أي نقل.", "شارك منطق التصنيف والتحقق بين المشي الافتراضي والحقيقي؛ زوّد العدد بعد Rename الناجحة بس."},
		En: &LevelText{Title: "The Organizer", Goal: "Plan incoming crate placement before moving files", Files: []LabFile{{Name: "main.go", Starter: ch222Starter}}, Lesson: `# Put Incoming Crates on the Right Shelf

Incoming crates live under ~inbox~. Archive shelves are ~records~, ~images~, and ~misc~. A plan is reviewable data: discovery never moves files.

## Before you write
- How will a collision in the final crate prevent the first move?
- If the third move fails, how will you report completed moves?

## Walking a filesystem interface
~fs.WalkDir~ visits the root and then its children in lexical order. Its callback receives a path, ~fs.DirEntry~, and error. Check the error first because the entry can be nil on failure. Returning an error stops traversal.
~~~go
err:=fs.WalkDir(store,"inbox",func(name string,entry fs.DirEntry,err error)error{
 if err!=nil {return err}
 if entry.Type().IsRegular() {fmt.Println(name)}
 return nil
})
~~~
~Type().IsRegular()~ distinguishes regular files from directories, links, and special types. ~WalkDir~ does not follow symbolic links.

## Interface names and host paths
Import ~path/filepath~ and ~os~. ~filepath.WalkDir~ has the same callback shape using host paths. ~filepath.Join~ joins host paths; ~filepath.Rel~ computes a relative path; ~filepath.ToSlash~ converts it to an ~fs.FS~ name. ~os.DirFS~ exposes a real directory for reading, without providing a sandbox against malicious links.
~~~go
err:=filepath.WalkDir(filepath.Join(root,"inbox"),func(name string,entry fs.DirEntry,err error)error{
 if err!=nil {return err}
 relative,err:=filepath.Rel(root,name)
 if err!=nil {return err}
 fmt.Println(filepath.ToSlash(relative))
 return nil
})
store:=os.DirFS(root)
~~~
~filepath.FromSlash~ converts a plan name back to a host path. ~slices.Sort~ directly sorts a string slice.
~~~go
host:=filepath.Join(root,filepath.FromSlash("records/cell.json"))
names:=[]string{"b","a"}
slices.Sort(names)
fmt.Println(host,names)
~~~

Tests use ~t.TempDir()~ for isolated directories removed automatically. Every real file in this chapter belongs inside one.
~~~go
root:=t.TempDir()
err:=os.MkdirAll(filepath.Join(root,"inbox"),0700)
~~~
~os.MkdirAll~ creates missing directories and parents. ~0700~ grants permissions only to the owner.

## Classify, validate, then execute
~path.Ext~ returns the dotted extension; ~path.Base~ returns the last name. Apply ~strings.ToLower~ only to the extension, preserving crate names.
~~~go
fmt.Println(path.Ext("bay/NOTE.JSON"),path.Base("bay/NOTE.JSON")) // .JSON NOTE.JSON
~~~
Track destinations in a map to reject two sources targeting one name. Check destinations using ~fs.Stat~: only ~errors.Is(err,fs.ErrNotExist)~ means available. Also reject an existing target shelf that is not a directory. If stat reports absence, check ~fs.ReadDir~ names in an existing parent too: a dangling link still occupies a name.
~~~go
_,err:=fs.Stat(store,"records/cycle.json")
if err!=nil && !errors.Is(err,fs.ErrNotExist) {return err}
~~~
~os.Rename~ moves a name within one filesystem. It returns failures and may replace existing destinations on some systems. Complete validation must precede execution. The contract assumes a stable tree without external writers or links in destination directories.
~~~go
err:=os.Rename(filepath.Join(root,"inbox","cell.json"),filepath.Join(root,"records","cell.json"))
~~~
~ExecuteShelves~ accepts an injected rename function for testing partial failure; ~OrganizeShelves~ is a wrapper passing ~os.Rename~.
~~~go
plan,err:=PlanShelves(store)
example:=ShelfMove{From:"inbox/cell.csv",To:"records/cell.csv"}
n,err:=OrganizeShelves(root)
fmt.Println(plan,example,n,err)
n,err=ExecuteShelves(root,os.Rename)
~~~

There is no automatic copy fallback or rollback. Count completed moves and return the first error.

## Common mistakes
~~~go
os.MkdirAll(destination,0700) // Wrong during dry-run: mutates disk
if err!=nil { /* Assuming absence: it could be a permission error */ }
~~~

## Your mission
Implement ~PlanShelves(store fs.FS) ([]ShelfMove,error)~ and ~OrganizeShelves(root string) (int,error)~. ~ShelfMove~ has ~From,To string~ using relative slash-separated names. Require an existing directory ~inbox~, possibly empty. Walk recursively, skipping links and special files. Route ~.json~ and ~.csv~ to ~records~, ~.png~ to ~images~, and everything else including extensionless files to ~misc~. Match extensions case-insensitively, preserve basenames, and flatten nested directories. Sort the plan by ascending From.

Return a nil plan on duplicate destinations, any existing destination, an existing shelf that is not a directory, or traversal/stat failure. Preserve I/O causes. PlanShelves is dry-run: never create, move, or write anything.

~OrganizeShelves~ wraps ~ExecuteShelves(root string, rename func(string,string)error) (int,error)~ with ~os.Rename~. The non-nil callback actually moves files on success and leaves them unchanged on failure. ExecuteShelves discovers real incoming files with ~filepath.WalkDir~ and validates the entire plan before mutation. Then create only required shelf directories with 0700 and move using the injected callback in plan order. Return the number completed; stop at the first failure without rolling back prior moves. Repeating after success yields an empty plan. The tree is stable during execution and existing destination directories are real directories without links. No concurrency or Report. Real tests use only ~t.TempDir()~.
`, Hints: []string{"Collect regular file names before building the plan.", "Use a destination map and stat every target and parent before moving.", "Share classification and validation between interface and host walks; count only successful renames."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch222Hidden, Bonus: false},
	})
}

const ch222Starter = `package main

import (
	"errors"
	"io/fs"
)

type ShelfMove struct{ From, To string }

func PlanShelves(store fs.FS) ([]ShelfMove, error) { return nil, errors.New("unfinished") }
func OrganizeShelves(root string) (int, error)     { return 0, errors.New("unfinished") }
func ExecuteShelves(root string, rename func(string, string) error) (int, error) {
	return 0, errors.New("unfinished")
}
func main() {}
`
const ch222StarterAr = `package main

import (
	"errors"
	"io/fs"
)

type ShelfMove struct{ From, To string }

func PlanShelves(store fs.FS) ([]ShelfMove, error) { return nil, errors.New("لسه مخلصتش") }
func OrganizeShelves(root string) (int, error)     { return 0, errors.New("لسه مخلصتش") }
func ExecuteShelves(root string, rename func(string, string) error) (int, error) {
	return 0, errors.New("لسه مخلصتش")
}
func main() {}
`
const ch222Hidden = `package main

import (
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"testing/fstest"
)

func TestHiddenOrganizerSeededPlanAndExecution(t *testing.T) {
	rng := rand.New(rand.NewSource(222))
	for round := 0; round < 12; round++ {
		store := fstest.MapFS{"inbox": &fstest.MapFile{Mode: fs.ModeDir | 0700}}
		var want []ShelfMove
		for i := 0; i < 20; i++ {
			ext := []string{".JSON", ".csv", ".png", ".alien", ""}[rng.Intn(5)]
			shelf := "misc"
			switch ext {
			case ".JSON", ".csv":
				shelf = "records"
			case ".png":
				shelf = "images"
			}
			base := fmt.Sprintf("crate%02d%s", i, ext)
			from := fmt.Sprintf("inbox/bay%d/%s", rng.Intn(3), base)
			store[from] = &fstest.MapFile{Data: []byte(fmt.Sprint(rng.Int()))}
			want = append(want, ShelfMove{from, shelf + "/" + base})
		}
		slices.SortFunc(want, func(a, b ShelfMove) int {
			if a.From < b.From {
				return -1
			}
			if a.From > b.From {
				return 1
			}
			return 0
		})
		before := fstest.MapFS{}
		for name, file := range store {
			clone := *file
			clone.Data = append([]byte(nil), file.Data...)
			before[name] = &clone
		}
		got, err := PlanShelves(store)
		if err != nil || !reflect.DeepEqual(got, want) || !reflect.DeepEqual(store, before) {
			t.Fatal("plan", got, want, err)
		}
		root := t.TempDir()
		for name, file := range store {
			host := filepath.Join(root, filepath.FromSlash(name))
			if file.Mode.IsDir() {
				if err := os.MkdirAll(host, 0700); err != nil {
					t.Fatal(err)
				}
			} else {
				if err := os.MkdirAll(filepath.Dir(host), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(host, file.Data, 0600); err != nil {
					t.Fatal(err)
				}
			}
		}
		if n, err := OrganizeShelves(root); err != nil || n != len(want) {
			t.Fatal("execute", n, err)
		}
		for _, move := range want {
			raw, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(move.To)))
			if err != nil || string(raw) != string(store[move.From].Data) {
				t.Fatal("bytes", move, err)
			}
			if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(move.From))); !errors.Is(err, fs.ErrNotExist) {
				t.Fatal("source remains", err)
			}
		}
		if n, err := OrganizeShelves(root); err != nil || n != 0 {
			t.Fatal("repeat", n, err)
		}
	}
}
func TestHiddenOrganizerRejectsBeforeMutation(t *testing.T) {
	for _, extra := range []fstest.MapFS{
		{"inbox/b/a.json": &fstest.MapFile{}}, {"records/a.json": &fstest.MapFile{}}, {"records": &fstest.MapFile{}},
	} {
		store := fstest.MapFS{"inbox/a.json": &fstest.MapFile{Data: []byte("keep")}}
		for k, v := range extra {
			store[k] = v
		}
		if got, err := PlanShelves(store); err == nil || got != nil {
			t.Fatal("collision", got, err)
		}
		root := t.TempDir()
		for name, file := range store {
			host := filepath.Join(root, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(host), 0700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(host, file.Data, 0600); err != nil {
				t.Fatal(err)
			}
		}
		if n, err := OrganizeShelves(root); err == nil || n != 0 {
			t.Fatal("preflight", n, err)
		}
		raw, err := os.ReadFile(filepath.Join(root, "inbox", "a.json"))
		if err != nil || string(raw) != "keep" {
			t.Fatal("preflight moved source", err)
		}
	}
	for _, store := range []fstest.MapFS{{}, {"inbox": &fstest.MapFile{}}} {
		if got, err := PlanShelves(store); err == nil || got != nil {
			t.Fatal("invalid inbox", got, err)
		}
	}
	store := fstest.MapFS{"inbox": &fstest.MapFile{Mode: fs.ModeDir | 0700}, "inbox/link": &fstest.MapFile{Mode: fs.ModeSymlink}, "inbox/pipe": &fstest.MapFile{Mode: fs.ModeNamedPipe}}
	if got, err := PlanShelves(store); err != nil || len(got) != 0 {
		t.Fatal("skip special", got, err)
	}
	wire := errors.New("stat offline")
	if got, err := PlanShelves(hiddenShelfFS{store, wire}); !errors.Is(err, wire) || got != nil {
		t.Fatal("walk cause", got, err)
	}
}
func TestHiddenOrganizerMidMoveFailureStopsAndCounts(t *testing.T) {
	wire := errors.New("move failed")
	for failAt := 0; failAt < 3; failAt++ {
		root := t.TempDir()
		if err := os.MkdirAll(filepath.Join(root, "inbox"), 0700); err != nil {
			t.Fatal(err)
		}
		for _, name := range []string{"a.json", "b.json", "c.json"} {
			if err := os.WriteFile(filepath.Join(root, "inbox", name), []byte(name), 0600); err != nil {
				t.Fatal(err)
			}
		}
		calls := 0
		n, err := ExecuteShelves(root, func(from, to string) error {
			calls++
			if calls == failAt+1 {
				return wire
			}
			return os.Rename(from, to)
		})
		if !errors.Is(err, wire) || n != failAt || calls != failAt+1 {
			t.Fatal("partial progress", n, err, calls)
		}
		for i, name := range []string{"a.json", "b.json", "c.json"} {
			shelf := "inbox"
			if i < failAt {
				shelf = "records"
			}
			raw, err := os.ReadFile(filepath.Join(root, shelf, name))
			if err != nil || string(raw) != name {
				t.Fatal("partial state", name, err)
			}
		}
	}
}

func TestHiddenOrganizerDanglingDestinationIsOccupied(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{"inbox", "records"} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "inbox", "a.json"), []byte("keep"), 0600); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "records", "a.json")
	if err := os.Symlink("missing", target); err != nil {
		t.Fatal(err)
	}
	if n, err := OrganizeShelves(root); err == nil || n != 0 {
		t.Fatal("overwrote dangling link", n, err)
	}
	if link, err := os.Readlink(target); err != nil || link != "missing" {
		t.Fatal("link changed", link, err)
	}
}
func TestHiddenOrganizerStatFailureAndEmptyPlan(t *testing.T) {
	wire := errors.New("stat offline")
	base := fstest.MapFS{"inbox/a.json": &fstest.MapFile{}}
	if plan, err := PlanShelves(hiddenShelfStatFS{base, wire}); plan != nil || !errors.Is(err, wire) {
		t.Fatal("stat cause", plan, err)
	}
	empty := fstest.MapFS{"inbox": &fstest.MapFile{Mode: fs.ModeDir | 0700}}
	if plan, err := PlanShelves(empty); err != nil || len(plan) != 0 {
		t.Fatal("empty plan", plan, err)
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "inbox"), 0700); err != nil {
		t.Fatal(err)
	}
	if n, err := OrganizeShelves(root); err != nil || n != 0 {
		t.Fatal("empty execute", n, err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatal("created unneeded shelves", entries, err)
	}
}

type hiddenShelfStatFS struct {
	fs.FS
	err error
}

func (s hiddenShelfStatFS) Stat(name string) (fs.FileInfo, error) {
	if name == "records/a.json" {
		return nil, s.err
	}
	return fs.Stat(s.FS, name)
}

type hiddenShelfFS struct {
	fs.FS
	err error
}

func (s hiddenShelfFS) ReadDir(string) ([]fs.DirEntry, error) { return nil, s.err }
`

func init() {
	add(&Level{ID: "22-3-atomic-save", Chapter: 22, Kind: "lab", Title: "حفظ ذرّي", Goal: "احفظ حزمة المنارة من غير ما تضيّع النسخة السليمة عند الفشل",
		Requires: []string{"defer", "sel:os.Rename"}, Bridge: []int{17}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch223StarterAr}},
		Lesson: `# آخر حزمة سليمة

المنارة بتستقبل حزمة ممكن توليدها يفشل بعد أول بايتات. فتح الملف النهائي للكتابة يمسح النسخة السليمة بدري. هنكتب لاسم مؤقت وبعدين ننقل الاسم بعد اكتمال الكتابة والقفل.

## قبل ما تكتب
- إيه آخر خطوة ممكن تفشل قبل نشر الحزمة؟
- إزاي التست يفرض فشل النقل من غير تغيير صلاحيات جهازك؟

## ملف مؤقت جنب الهدف
~os.CreateTemp~ تنشئ اسم فريد وتفتحه للقراءة والكتابة، بصلاحيات 0600 قبل umask. النجمة في النمط بتتبدل بجزء عشوائي. اختار مجلد الهدف بـ ~filepath.Dir~ عشان النقل يفضل على نفس نظام الملفات. ~file.Name~ اسم الملف الحقيقي.
~~~go
file,err:=os.CreateTemp(filepath.Dir(target),".beacon-*")
if err!=nil {return err}
name:=file.Name()
~~~
الاسم العشوائي من نظام التشغيل مش جزء من النتيجة؛ التست يفحص الأثر، مش اسم معين.

## عمر الملف وdefer
~*os.File~ تحقق ~io.Writer~. الـ callback تكتب بس، ومش بتقفل الملف أو تحتفظ بيه بعد الرجوع. ~Close~ ممكن تفشل؛ افحصها قبل النقل. ~os.Remove~ تشيل الاسم المؤقت. ~defer~ من فصل 5 تضمن التنظيف حتى لو رجعت بدري.
~~~go
defer os.Remove(name)
defer file.Close()
err=write(file)
if err!=nil {return err}
if err=file.Close();err!=nil {return err}
~~~
القفل المؤجل يحمي مسار فشل callback؛ القفل الصريح قبل النقل يفحص الفشل. القفل المكرر في التنظيف مش بديل لفحص القفل الصريح. لو التنظيف نفسه فشل، حافظ على خطأ المرحلة الأصلي؛ مفيش ضمان حذف اسم رفض نظام الملفات حذفه.

## حد النشر
~os.Rename~ من المهمة السابقة تنقل الملف المكتمل. متشيلش الهدف الأول، ومتعملش fallback بالنسخ لو النقل فشل. الـ API مش ضمان عام لذرّية الاستبدال على كل نظام تشغيل؛ إحنا بنضمن إننا منكتبش حزمة جزئية للهدف. مفيش وعد بتحمل انقطاع الكهرباء أو ~Sync~ في العقد.
~~~go
err=os.Rename(name,target)
~~~
عشان نجرب فشل النقل بشكل حتمي، مرر function value بنفس التوقيع. ~SaveRelay~ توصل التنفيذ الحقيقي؛ ~CommitRelay~ تستخدم المصدر المحقون.
~~~go
rename:=func(from,to string)error{return os.Rename(from,to)}
err:=SaveRelay(target,func(out io.Writer)error{
 return json.NewEncoder(out).Encode([]int{3,7})
})
err=CommitRelay(target,func(out io.Writer)error{
 _,err:=fmt.Fprint(out,"beacon=ready\n")
 return err
},rename)
~~~
الـ callback ممكن تطبع JSON بـ ~json.NewEncoder(out).Encode(packet)~ لو محتاجة؛ استورد ~encoding/json~. ~json.NewEncoder~ تربط encoder بالـ writer، و ~Encode~ تكتب مستند JSON وسطر جديد وترجع الخطأ. صيغة الحزمة في المهمة اختيار المنادي، مش schema مفروضة.
~~~go
err:=json.NewEncoder(out).Encode([]int{3,7})
~~~

## أخطاء شائعة
~~~go
os.WriteFile(target,partial,0600) // مسح النسخة القديمة قبل اكتمال الحزمة
os.Remove(target) // خلق فجوة وخسارة النسخة لو Rename فشلت
~~~

## المهمة
نفّذ ~SaveRelay(target string, write func(io.Writer)error) error~ كغلاف ينادي ~CommitRelay~ مع ~os.Rename~. نفّذ ~CommitRelay(target string, write func(io.Writer)error, rename func(string,string)error) error~. الدالتان المحقونتان غير nil. مجلد الهدف لازم موجود؛ الهدف ممكن مش موجود أو ملف عادي. التنفيذ تسلسلي، من غير كاتب خارجي.

أنشئ مؤقت فريد جنب الهدف بصلاحيات CreateTemp الافتراضية، وشغّل write مرة واحدة عليه. اكتب للمؤقت بس، واقفله بنجاح قبل استدعاء rename مرة واحدة من اسم المؤقت للهدف. رجّع أول خطأ إنشاء أو callback أو قفل أو نقل مع الحفاظ على سببه. فشل قبل النقل ميناديش rename، والهدف يفضل كما هو أو غايب لو كان غايب. نجاح النقل ينشر البايتات كاملة، بما فيها حزمة فاضية.

استخدم defer لتنظيف الموارد: اقفل الملف على كل مسار وشيل الاسم المؤقت بعد النجاح أو الفشل لما نظام الملفات يسمح. متحذفش الهدف ولا تنسخ للمكان النهائي بعد فشل النقل. المنادي المحقون لـ rename إما ينقل فعلاً أو يرجع خطأ من غير تغيير ملفات. النجمة التالتة: فشل rename المتكرر يحافظ على الهدف ويشيل المؤقت، ومحاولة لاحقة تنجح من غير بقايا. مفيش Report؛ كل اختبار حقيقي داخل ~t.TempDir()~.
`, Hints: []string{"أنشئ المؤقت في filepath.Dir(target)، وسجّل تنظيفه بـ defer فوراً.", "مرر المؤقت للـ callback؛ افحص خطأها وخطأ Close قبل Rename.", "SaveRelay تمرر os.Rename؛ الاختبار المحقون يرفض النقل ويثبت إن الهدف القديم ما اتلمسش."},
		En: &LevelText{Title: "Atomic Save", Goal: "Save a beacon packet without losing the last complete version", Files: []LabFile{{Name: "main.go", Starter: ch223Starter}}, Lesson: `# The Last Complete Packet

A beacon packet can fail generation after its first bytes. Opening the final file would erase the last complete version too early. Write under a temporary name, then publish after writing and closing succeed.

## Before you write
- What is the last step that can fail before publication?
- How can tests force rename failure without changing host permissions?

## A temporary file beside the target
~os.CreateTemp~ creates a unique name and opens it for reading and writing, with mode 0600 before umask. A star in its pattern is replaced with a random component. Choose the target directory using ~filepath.Dir~ to stay on the same filesystem. ~file.Name~ returns the actual name.
~~~go
file,err:=os.CreateTemp(filepath.Dir(target),".beacon-*")
if err!=nil {return err}
name:=file.Name()
~~~
The operating system's random name is not part of the result. Tests assert effects rather than a particular generated name.

## File lifetime and defer
~*os.File~ satisfies ~io.Writer~. The callback only writes; it neither closes the file nor retains it after returning. ~Close~ can fail: check it before renaming. ~os.Remove~ removes the temporary name. ~defer~ from chapter 5 guarantees cleanup on early returns.
~~~go
defer os.Remove(name)
defer file.Close()
err=write(file)
if err!=nil {return err}
if err=file.Close();err!=nil {return err}
~~~
Deferred closing protects callback failure paths; explicit closing checks failure before publication. A repeated cleanup close cannot replace that check. If cleanup fails, preserve the original stage error; removal cannot be guaranteed when the filesystem refuses it.

## The publication boundary
~os.Rename~ from the preceding mission moves the complete file. Never remove the target first or fall back to copying after rename failure. The API does not guarantee atomic replacement on every operating system; this contract guarantees that partial packets are never written into the target. Power-loss durability and ~Sync~ are outside the contract.
~~~go
err=os.Rename(name,target)
~~~
Inject a function value with the same signature to deterministically test failure. ~SaveRelay~ supplies the real implementation; ~CommitRelay~ uses the injected boundary.
~~~go
rename:=func(from,to string)error{return os.Rename(from,to)}
err:=SaveRelay(target,func(out io.Writer)error{
 return json.NewEncoder(out).Encode([]int{3,7})
})
err=CommitRelay(target,func(out io.Writer)error{
 _,err:=fmt.Fprint(out,"beacon=ready\n")
 return err
},rename)
~~~
A callback can produce JSON using ~json.NewEncoder(out).Encode(packet)~. Import ~encoding/json~: ~json.NewEncoder~ connects an encoder to a writer and ~Encode~ writes a JSON document with a newline, returning errors. The caller chooses packet syntax; no schema is imposed here.
~~~go
err:=json.NewEncoder(out).Encode([]int{3,7})
~~~

## Common mistakes
~~~go
os.WriteFile(target,partial,0600) // Erases the old version before completion
os.Remove(target) // Creates a gap and loses the old version if Rename fails
~~~

## Your mission
Implement ~SaveRelay(target string, write func(io.Writer)error) error~ as a wrapper calling ~CommitRelay~ with ~os.Rename~. Implement ~CommitRelay(target string, write func(io.Writer)error, rename func(string,string)error) error~. Both callbacks are non-nil. The target directory must exist; the target may be absent or a regular file. Calls are sequential without external writers.

Create a unique temporary file beside the target using CreateTemp's default permissions and invoke write exactly once on it. Write only to the temporary file and close it successfully before invoking rename exactly once from its name to target. Return the first create, callback, close, or rename error preserving its cause. Failures before rename must not invoke it. The target stays unchanged, or absent if initially absent. Successful rename publishes all bytes, including an empty packet.

Use defer for resource cleanup: close on every path and remove the temporary name after success or failure when the filesystem permits. Never delete the target or copy into it after rename failure. The injected rename either actually moves the file or returns an error without modifying files. Third star: repeated rename failures preserve the target and clean temporary files; a subsequent attempt succeeds without leftovers. No Report. All real tests live inside ~t.TempDir()~.
`, Hints: []string{"Create the temporary file in filepath.Dir(target) and defer cleanup immediately.", "Give the callback the temporary writer; check its error and Close before Rename.", "SaveRelay passes os.Rename; an injected rejection proves the previous target is untouched."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch223Hidden, Bonus: true},
	})
}

const ch223Starter = `package main

import (
	"errors"
	"io"
)

func SaveRelay(target string, write func(io.Writer) error) error {
	return errors.New("unfinished")
}
func CommitRelay(target string, write func(io.Writer) error, rename func(string, string) error) error {
	return errors.New("unfinished")
}
func main() {}
`
const ch223StarterAr = `package main

import (
	"errors"
	"io"
)

func SaveRelay(target string, write func(io.Writer) error) error {
	return errors.New("لسه مخلصتش")
}
func CommitRelay(target string, write func(io.Writer) error, rename func(string, string) error) error {
	return errors.New("لسه مخلصتش")
}
func main() {}
`
const ch223Hidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func hiddenRelayContents(t *testing.T, dir, target string, want []byte) {
	t.Helper()
	raw, err := os.ReadFile(target)
	if err != nil || !bytes.Equal(raw, want) {
		t.Fatal("target", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 1 {
		t.Fatal("temporary leaked", entries, err)
	}
}
func TestHiddenRelaySeededCompletePackets(t *testing.T) {
	rng := rand.New(rand.NewSource(223))
	dir := t.TempDir()
	target := filepath.Join(dir, "beacon.packet")
	for _, size := range []int{0, 1, 4096, 131073} {
		want := make([]byte, size)
		if _, err := rng.Read(want); err != nil {
			t.Fatal(err)
		}
		calls := 0
		err := SaveRelay(target, func(out io.Writer) error { calls++; _, err := out.Write(want); return err })
		if err != nil || calls != 1 {
			t.Fatal("save", calls, err)
		}
		hiddenRelayContents(t, dir, target, want)
	}
}
func TestHiddenRelayPartialFailureNeverPublishes(t *testing.T) {
	wire := errors.New("generator failed")
	rng := rand.New(rand.NewSource(2231))
	for _, present := range []bool{false, true} {
		for trial := 0; trial < 15; trial++ {
			dir := t.TempDir()
			target := filepath.Join(dir, "relay")
			if present {
				if err := os.WriteFile(target, []byte("old"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			renames := 0
			var captured io.Writer
			err := CommitRelay(target, func(out io.Writer) error {
				captured = out
				_, err := fmt.Fprint(out, string(bytes.Repeat([]byte("x"), 1+rng.Intn(5000))))
				if err != nil {
					return err
				}
				return wire
			}, func(string, string) error { renames++; return nil })
			if !errors.Is(err, wire) || renames != 0 {
				t.Fatal("pre-commit failure", err, renames)
			}
			if _, err := captured.Write([]byte("late")); !errors.Is(err, os.ErrClosed) {
				t.Fatal("file not closed", err)
			}
			if present {
				hiddenRelayContents(t, dir, target, []byte("old"))
			} else {
				if _, err := os.Stat(target); !errors.Is(err, os.ErrNotExist) {
					t.Fatal("created failed target", err)
				}
				entries, err := os.ReadDir(dir)
				if err != nil || len(entries) != 0 {
					t.Fatal("leak", entries, err)
				}
			}
		}
	}
	dir := t.TempDir()
	calls := 0
	if err := SaveRelay(filepath.Join(dir, "absent", "relay"), func(io.Writer) error { calls++; return nil }); !errors.Is(err, os.ErrNotExist) || calls != 0 {
		t.Fatal("create failure", err, calls)
	}
}
func TestHiddenRelayClosedBeforeSameDirectoryRename(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "relay")
	var captured io.Writer
	calls := 0
	err := CommitRelay(target, func(out io.Writer) error { captured = out; _, err := fmt.Fprint(out, "ready"); return err }, func(from, to string) error {
		calls++
		if filepath.Dir(from) != dir || from == target || to != target {
			t.Fatal("rename paths", from, to)
		}
		if _, err := captured.Write([]byte("late")); !errors.Is(err, os.ErrClosed) {
			t.Fatal("rename before close", err)
		}
		raw, err := os.ReadFile(from)
		if err != nil || string(raw) != "ready" {
			t.Fatal("incomplete temporary", err)
		}
		return os.Rename(from, to)
	})
	if err != nil || calls != 1 {
		t.Fatal("commit", err, calls)
	}
	hiddenRelayContents(t, dir, target, []byte("ready"))
}
func TestHiddenBonusRelayFailedRenameCleanupAndRetry(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "relay")
	if err := os.WriteFile(target, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	wire := errors.New("rename rejected")
	for i := 0; i < 5; i++ {
		calls := 0
		err := CommitRelay(target, func(out io.Writer) error { _, err := fmt.Fprintf(out, "new-%d", i); return err }, func(from, to string) error { calls++; return wire })
		if !errors.Is(err, wire) || calls != 1 {
			t.Fatal("rename cause", err, calls)
		}
		hiddenRelayContents(t, dir, target, []byte("old"))
	}
	if err := SaveRelay(target, func(out io.Writer) error { _, err := fmt.Fprint(out, "recovered"); return err }); err != nil {
		t.Fatal(err)
	}
	hiddenRelayContents(t, dir, target, []byte("recovered"))
}
`

func init() {
	add(&Level{ID: "22-4-duplicates", Chapter: 22, Kind: "lab", Title: "الملفات المكررة", Goal: "اكتشف نسخ حزم المستعمرة بالحجم وبصمة stream",
		Requires: []string{"sel:sha256.New", "sel:io.Copy"}, Bridge: []int{25}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch224StarterAr}},
		Lesson: `# نفس الحزمة في أكتر من رف

الأرشيف فيه نسخ متكررة من حزم المنارات. الحجم فلتر رخيص، بس ملفين بنفس الحجم ممكن مختلفين. هنحسب بصمة للمجموعات المرشحة من غير تحميل كل ملف في الذاكرة.

## قبل ما تكتب
- إزاي تثبت إن ملف حجمه فريد متفتحش أصلاً؟
- لو آخر ملف فشل بعد اكتشاف مجموعة مكررة، هل ترجع تقرير ناقص؟

## الحجم قبل المحتوى
امشِ بـ ~fs.WalkDir~ من ~"."~. ~entry.Info()~ ترجع ~fs.FileInfo~ وخطأ؛ ~Size()~ ترجع الحجم كـ ~int64~. افحص النوع قبل Info عشان تتجاهل الروابط والأنواع الخاصة. اجمع المسارات في ~map[int64][]string~.
~~~go
info,err:=entry.Info()
if err!=nil {return err}
size:=info.Size()
buckets[size]=append(buckets[size],name)
~~~
مجموعة فيها ملف واحد مش ممكن تكرر غيرها بالحجم؛ متفتحهاش. ملفين فاضيين حجمهم 0 ولسه مجموعة مرشحة.

## بصمة بتستقبل stream
استورد ~crypto/sha256~ و ~io~. ~sha256.New~ ترجع hash بيحقق ~io.Writer~. ~io.Copy~ تنقل من reader لـ writer بذاكرة محدودة؛ بترجع عدد البايتات ~int64~ والخطأ. ~Sum(nil)~ ترجع بايتات البصمة الحالية من غير تغيير الحالة. دي بصمة محتوى، مش تشفير ولا كلمة سر.
~~~go
hash:=sha256.New()
n,err:=io.Copy(hash,strings.NewReader("antenna"))
digest:=fmt.Sprintf("%x",hash.Sum(nil))
fmt.Println(n,len(digest),err) // 7 64 <nil>
~~~
~%x~ تعرض البايتات كـ hexadecimal. اعمل hash جديدة لكل ملف؛ استخدام واحدة للمجموعة يخلط المحتوى. متستخدمش ~io.ReadAll~ أو ~fs.ReadFile~ عشان الملفات ممكن كبيرة.

## افتح واقفل وراجع الحجم
~store.Open(name)~ ترجع ~fs.File~ بتدعم ~Read~ و ~Close~ و ~Stat~. اقفل كل ملف حتى لو القراية فشلت. افحص خطأ القفل لو القراية نجحت؛ لو الاتنين فشلوا، رجع خطأ القراية الأصلي.
~~~go
file,err:=store.Open(name)
if err!=nil {return err}
hash:=sha256.New()
n,copyErr:=io.Copy(hash,file)
closeErr:=file.Close()
if copyErr!=nil {return copyErr}
if closeErr!=nil {return closeErr}
~~~
قارن n بالحجم المسجل. اختلافه دليل إن الملف اتغير أو المصدر غير ثابت؛ ارفض التقرير. التست يقدر يحقن reader بترفض طلب buffer أكبر من 32 KiB عشان يمنع تحميل ملف كبير دفعة واحدة.

## ترتيب التقرير
بعد الحجم، اجمع حسب البصمة؛ نفس الحجم لوحده مش دليل تكرار. رجّع بس المجموعات اللي فيها مسارين أو أكتر. رتّب أسماء كل مجموعة، وبعدين المجموعات حسب أول اسم بـ ~slices.SortFunc~ و ~strings.Compare~ زي الفصول السابقة.
~~~go
group:=[]string{"relay/b","relay/a"}
slices.Sort(group)
fmt.Println(group[0]) // relay/a
~~~

## نداء التقرير
~~~go
groups,err:=PacketEchoes(store)
fmt.Println(groups,err)
~~~
~PacketEchoes~ ترجع بيانات التقرير؛ العرض أو الحذف قرار المنادي.

## أخطاء شائعة
~~~go
raw,_:=fs.ReadFile(store,name) // تحميل كامل وبلع خطأ
if len(buckets[size])>1 { /* غلط: نفس الحجم مش بالضرورة نفس المحتوى */ }
~~~

## المهمة
نفّذ ~PacketEchoes(store fs.FS) ([][]string,error)~. store غير nil؛ الجذر ~"."~ لازم مجلد. امشِ recursively على الملفات العادية بس، واتجاهل الروابط والأنواع الخاصة من غير فتحها. اجمع بالحجم أولاً؛ ممنوع فتح محتوى ملف حجمه فريد. للمرشحين احسب SHA-256 بـ ~sha256.New~ و ~io.Copy~، بقراءة stream وذاكرة محدودة لكل ملف، مش ~ReadAll~ أو تحميل الملف كله.

رجّع مجموعات بنفس الحجم والبصمة، طول كل مجموعة 2 على الأقل؛ الاسم نسبي slash-separated. رتّب داخل المجموعة تصاعدي بترتيب Go، ورتّب المجموعات حسب أول اسم. الفاضي أو بدون تكرار نجاح بطول صفر. الملفات الفاضية ممكن تتكرر. مفيش حذف أو كتابة أو تغيير للمصدر.

أي فشل walk أو Info أو Open أو Read أو Close يرجع nil,error يحتفظ بالسبب وبالمسار. لو عدد البايتات المقروءة اختلف عن الحجم المسجل، رجّع nil,error. اقفل كل ملف اتفتح، حتى وقت الفشل، ومتفتحش ملفات زيادة بعد أول فشل. مفيش تزامن أو Report؛ التستات على نظام محقون، من غير قرص الجهاز.
`, Hints: []string{"اجمع map من الحجم لقائمة المسارات قبل فتح المحتوى.", "لكل حجم متكرر، hash جديدة لكل ملف وio.Copy، وبعدها Close وفحص الحجم.", "اجمع حسب البصمة داخل مجموعة الحجم؛ صفّي المجموعات الفردية ورتّب النتيجة كلها."},
		En: &LevelText{Title: "Duplicate Files", Goal: "Find repeated colony packets using size buckets and streaming hashes", Files: []LabFile{{Name: "main.go", Starter: ch224Starter}}, Lesson: `# One Packet on Several Shelves

The archive contains repeated beacon packets. Size is a cheap filter, but equal-size files can differ. Hash candidate groups without loading whole files into memory.

## Before you write
- How will you prove that a unique-size file was never opened?
- If the last file fails after finding duplicates, should an incomplete report escape?

## Size before content
Use ~fs.WalkDir~ from ~"."~. ~entry.Info()~ returns ~fs.FileInfo~ and an error; ~Size()~ returns an ~int64~ byte count. Check type before Info to skip links and special files. Collect paths in ~map[int64][]string~.
~~~go
info,err:=entry.Info()
if err!=nil {return err}
size:=info.Size()
buckets[size]=append(buckets[size],name)
~~~
A bucket with one file cannot match any other size. Never open its content. Two empty files both have size 0 and remain candidates.

## A hash consumes a stream
Import ~crypto/sha256~ and ~io~. ~sha256.New~ returns a hash implementing ~io.Writer~. ~io.Copy~ transfers from a reader to a writer using bounded memory and returns an ~int64~ byte count plus error. ~Sum(nil)~ returns the current digest bytes without changing state. This is a content fingerprint, not encryption or a password.
~~~go
hash:=sha256.New()
n,err:=io.Copy(hash,strings.NewReader("antenna"))
digest:=fmt.Sprintf("%x",hash.Sum(nil))
fmt.Println(n,len(digest),err) // 7 64 <nil>
~~~
~%x~ renders bytes in hexadecimal. Create a fresh hash for each file; sharing one across files combines their content. Avoid ~io.ReadAll~ and ~fs.ReadFile~ because files may be large.

## Open, close, and verify size
~store.Open(name)~ returns an ~fs.File~ supporting ~Read~, ~Close~, and ~Stat~. Close every opened file even when reading fails. Return close errors after successful reads; when both fail, retain the original read error.
~~~go
file,err:=store.Open(name)
if err!=nil {return err}
hash:=sha256.New()
n,copyErr:=io.Copy(hash,file)
closeErr:=file.Close()
if copyErr!=nil {return copyErr}
if closeErr!=nil {return closeErr}
~~~
Compare n with the recorded size. A mismatch means a changed file or inconsistent source, so reject the report. Tests can inject a reader refusing buffers larger than 32 KiB to detect whole-file loading.

## Report order
After size bucketing, group by digest: equal size alone does not prove duplication. Keep groups containing at least two paths. Sort each group, then order groups by their first path using ~slices.SortFunc~ and ~strings.Compare~ from previous chapters.
~~~go
group:=[]string{"relay/b","relay/a"}
slices.Sort(group)
fmt.Println(group[0]) // relay/a
~~~

## Calling the report
~~~go
groups,err:=PacketEchoes(store)
fmt.Println(groups,err)
~~~
~PacketEchoes~ returns report data; presentation or removal belongs to the caller.

## Common mistakes
~~~go
raw,_:=fs.ReadFile(store,name) // Loads everything and discards errors
if len(buckets[size])>1 { /* Equal size need not mean equal content */ }
~~~

## Your mission
Implement ~PacketEchoes(store fs.FS) ([][]string,error)~. The non-nil store must have a directory root ~"."~. Walk recursively over regular files only, ignoring links and special files without opening them. Bucket by size first and never open unique-size content. Hash candidates with ~sha256.New~ and ~io.Copy~ using streaming reads and bounded memory per file, without ReadAll or whole-file loading.

Return groups sharing both size and digest, with at least two paths per group. Names are relative slash-separated paths. Sort each group in Go string order and sort groups by their first path. Empty input or no duplicates succeeds with length zero. Empty files may be duplicates. Never delete, write, or mutate the source.

Any walk, Info, Open, Read, or Close failure returns nil,error preserving cause and path. A read byte count different from recorded size also returns nil,error. Close every opened file even on failure and open no further files after the first failure. No concurrency or Report. Tests use injected filesystems without host disk access.
`, Hints: []string{"Build size-to-path buckets before opening content.", "For repeated sizes, use a fresh hash and io.Copy per file, then Close and verify the byte count.", "Group digests within each size bucket; discard singletons and sort the entire report."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: ch224Hidden, Bonus: false},
	})
}

const ch224Starter = `package main

import (
	"errors"
	"io/fs"
)

func PacketEchoes(store fs.FS) ([][]string, error) { return nil, errors.New("unfinished") }
func main()                                        {}
`
const ch224StarterAr = `package main

import (
	"errors"
	"io/fs"
)

func PacketEchoes(store fs.FS) ([][]string, error) { return nil, errors.New("لسه مخلصتش") }
func main()                                        {}
`
const ch224Hidden = `package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

func TestHiddenEchoesSeededSizesDigestsAndOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(224))
	for round := 0; round < 20; round++ {
		base := fstest.MapFS{}
		var want [][]string
		for family := 0; family < 8; family++ {
			content := fmt.Sprintf("%08d/%02d", rng.Intn(10000000), family)
			copies := 1 + rng.Intn(4)
			var group []string
			for copy := 0; copy < copies; copy++ {
				name := fmt.Sprintf("bay%d/packet-%02d-%02d", copy, family, copy)
				base[name] = &fstest.MapFile{Data: []byte(content)}
				group = append(group, name)
			}
			if copies > 1 {
				slices.Sort(group)
				want = append(want, group)
			}
		}
		base["empty/a"] = &fstest.MapFile{}
		base["empty/b"] = &fstest.MapFile{}
		want = append(want, []string{"empty/a", "empty/b"})
		base["unique"] = &fstest.MapFile{Data: []byte(strings.Repeat("u", 37))}
		base["link"] = &fstest.MapFile{Mode: fs.ModeSymlink}
		base["pipe"] = &fstest.MapFile{Mode: fs.ModeNamedPipe}
		slices.SortFunc(want, func(a, b []string) int { return strings.Compare(a[0], b[0]) })
		store := &hiddenEchoStore{FS: base, opened: map[string]int{}, closed: map[string]int{}}
		got, err := PacketEchoes(store)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatal("grouping", got, want, err)
		}
		if store.opened["unique"] != 0 || store.opened["link"] != 0 || store.opened["pipe"] != 0 {
			t.Fatal("opened excluded content", store.opened)
		}
		for name, n := range store.opened {
			if store.closed[name] != n {
				t.Fatal("leaked file", name)
			}
		}
	}
	for _, base := range []fstest.MapFS{{}, {"one": &fstest.MapFile{Data: []byte("x")}}, {"a": &fstest.MapFile{Data: []byte("a")}, "b": &fstest.MapFile{Data: []byte("b")}}} {
		if got, err := PacketEchoes(base); err != nil || len(got) != 0 {
			t.Fatal("no duplicates", got, err)
		}
	}
}
func TestHiddenEchoesLargeFilesUseBoundedReads(t *testing.T) {
	content := []byte(strings.Repeat("orbit", 300000))
	base := fstest.MapFS{"a": &fstest.MapFile{Data: content}, "b": &fstest.MapFile{Data: content}}
	store := &hiddenEchoStore{FS: base, opened: map[string]int{}, closed: map[string]int{}, bounded: true}
	got, err := PacketEchoes(store)
	if err != nil || !reflect.DeepEqual(got, [][]string{{"a", "b"}}) {
		t.Fatal("streaming", got, err)
	}
	for _, name := range []string{"a", "b"} {
		if store.closed[name] != 1 {
			t.Fatal("close", name)
		}
	}
}
func TestHiddenEchoesFailuresDiscardReportAndClose(t *testing.T) {
	wire := errors.New("archive failed")
	for _, stage := range []string{"walk", "info", "open", "read", "close", "size"} {
		base := fstest.MapFS{"a": &fstest.MapFile{Data: []byte("same")}, "b": &fstest.MapFile{Data: []byte("same")}}
		store := &hiddenEchoStore{FS: base, opened: map[string]int{}, closed: map[string]int{}, stage: stage, wire: wire}
		got, err := PacketEchoes(store)
		if err == nil || got != nil || !strings.Contains(err.Error(), "a") {
			t.Fatal("failure discarded/context", stage, got, err)
		}
		if stage != "size" && !errors.Is(err, wire) {
			t.Fatal("cause", stage, err)
		}
		for name, n := range store.opened {
			if store.closed[name] != n {
				t.Fatal("error leaked file", stage, name)
			}
		}
	}
}

func TestHiddenEchoesMissingRootReturnsContext(t *testing.T) {
	groups, err := PacketEchoes(hiddenEchoMissing{})
	if groups != nil || !errors.Is(err, fs.ErrNotExist) || !strings.Contains(err.Error(), ".") {
		t.Fatal("root context", groups, err)
	}
}

type hiddenEchoMissing struct{}

func (hiddenEchoMissing) Open(string) (fs.File, error) { return nil, fs.ErrNotExist }

type hiddenEchoStore struct {
	fs.FS
	opened, closed map[string]int
	bounded        bool
	stage          string
	wire           error
}

func (s *hiddenEchoStore) ReadDir(name string) ([]fs.DirEntry, error) {
	if s.stage == "walk" {
		return nil, &fs.PathError{Op: "read", Path: "a", Err: s.wire}
	}
	entries, err := fs.ReadDir(s.FS, name)
	if s.stage == "info" {
		for i, e := range entries {
			if e.Name() == "a" {
				entries[i] = hiddenEchoEntry{e, s.wire}
			}
		}
	}
	return entries, err
}
func (s *hiddenEchoStore) Open(name string) (fs.File, error) {
	if name == "a" && s.stage == "open" {
		return nil, s.wire
	}
	file, err := s.FS.Open(name)
	if err != nil {
		return nil, err
	}
	if name == "." {
		return file, nil
	}
	s.opened[name]++
	return &hiddenEchoFile{File: file, store: s, name: name}, nil
}

type hiddenEchoEntry struct {
	fs.DirEntry
	err error
}

func (e hiddenEchoEntry) Info() (fs.FileInfo, error) { return nil, e.err }

type hiddenEchoFile struct {
	fs.File
	store *hiddenEchoStore
	name  string
	reads int
}

func (f *hiddenEchoFile) Read(p []byte) (int, error) {
	f.reads++
	if f.store.bounded && len(p) > 32<<10 {
		return 0, fmt.Errorf("whole-file read: %d", len(p))
	}
	if f.name == "a" {
		switch f.store.stage {
		case "read":
			if f.reads > 1 {
				return 0, f.store.wire
			}
			if len(p) > 1 {
				p = p[:1]
			}
		case "size":
			return 0, io.EOF
		}
	}
	return f.File.Read(p)
}
func (f *hiddenEchoFile) Close() error {
	f.store.closed[f.name]++
	err := f.File.Close()
	if f.name == "a" && f.store.stage == "close" {
		return f.store.wire
	}
	return err
}
`
