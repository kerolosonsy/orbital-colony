package levels

import "orbital/colony/i18n"

func init() {
	add(&Level{ID: "18-1-mission-log", Chapter: 18, Kind: "lab", Title: "سجل المهام", Goal: "اختبر سجل طلعات بهوية ثابتة وترتيب مستقر",
		Requires: []string{"tabletest"}, Bridge: []int{16}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: missionLogImplementation, ReadOnly: true}, {Name: "main_test.go", Starter: missionLogTestStarter}},
		Lesson: `# سجل مش بيتلخبط

مركز القيادة محتاج سجل طلعات في الذاكرة. الكود موجود وصحيح؛ شغلك هنا تكتب table-driven tests تثبت إن هوية الطلعة مش مكانها الحالي في الـ slice. حذف طلعة مينفعش يخلّي طلعة جديدة تاخد رقمها.

## قبل ما تكتب
- هتختار تسلسل عمليات يفرّق بين عدّاد IDs وعدد السجلات إزاي؟
- هتثبت إزاي إن خطأ أو تعديل نسخة العرض مبيغيّرش السجل الأصلي؟

## جدول سيناريوهات مش نسخ من نفس التست
~~~go
func TestSizes(t *testing.T) {
 cases := []struct { name string; input string; want int }{
  {"empty", "", 0}, {"two bytes", "ab", 2},
 }
 for _,tc := range cases {
  t.Run(tc.name,func(t *testing.T) {
   if len(tc.input)!=tc.want {t.Fatalf("length = %d",len(tc.input))}
  })
 }
}
~~~
~testing.T~ و~t.Run~ و~t.Fatalf~ زي فصل 16. الجدول بيفصل اختلاف المدخلات عن خطوات التحقق المشتركة. في السجل استخدم جدول لحذف أول ووسط وآخر طلعة؛ اعمل ~var log MissionLog~ جديدة جوه كل subtest عشان الحالة متتسرّبش.

## CRUD وسياسة الهوية
CRUD هنا إنشاء، قراية، تحديث، وحذف؛ مش أوامر تيرمنال. ~MissionEntry~ فيها ~ID int~ و~Sector string~ و~Sealed bool~. الـ zero value من ~MissionLog~ جاهزة، والعداد منفصل عن الـ entries.
~~~go
var log MissionLog
id,err:=log.Log(" ice ")
err=log.Seal(id)
entries:=log.Entries()
fmt.Println(entries) // [{1 ice true}]
err=log.Erase(id)
fmt.Println(err) // <nil>
~~~
~Log~ بتشيل المسافات الطرفية بـ ~strings.TrimSpace~، وترجع ID وخطأ؛ ~Seal~ تثبّت ختم الطلعة، و~Erase~ تشيلها. ~Entries~ لقطة مستقلة بترتيب الإنشاء. ده استخدام للعقد، مش تنفيذ للتستات المطلوبة.

## قارن الحالة قبل وبعد
~~~go
before:=log.Entries()
err:=log.Seal(999)
if err==nil || !reflect.DeepEqual(before,log.Entries()) {
 t.Fatal("failed operation changed state")
}
~~~
استورد ~reflect~. ~reflect.DeepEqual~ تقارن حقول structs وعناصر slices بالترتيب؛ nil وslice فاضية مختلفين عندها. العقد بيسمح بالاتنين للقائمة الفاضية، فافحص ~len~ هناك. خُد snapshot قبل العملية نفسها، مش بعدها.

## أخطاء شائعة
~~~go
log.Seal(1) // التست بيعدّي حتى لو الختم مشتغلش: مفيش assertion
if len(log.Entries())==2 { return } // العدد مش بيثبت IDs ولا الترتيب
~~~
الـ mutants نسخ بتتبني بس بتخرق العقد. لازم تستاتك تفشل عليها بتأكيد سلوك غلط، وتعدّي على النسخة الصحيحة. مقارنة الخطأ بـ nil كفاية هنا؛ نصه مش جزء من العقد.

## المهمة
عدّل ~main_test.go~ بس؛ ~main.go~ للقراءة. اكتب table-driven tests بحالات مستقلة لـ ~MissionLog~: ~Log(sector string) (int,error)~ تضيف sector بعد TrimSpace؛ الفاضي يرفض بـ 0,error من غير استهلاك ID. الأرقام تبدأ 1 وتزيد واحد مع كل إضافة ناجحة، ومبتتكررش بعد الحذف.

~Seal(id int) error~ تختم الموجود؛ تكرار الختم نجاح من غير تغيير إضافي. ~Erase(id int) error~ تحذف الموجود بس وتحافظ على ترتيب الباقي. أي ID مش موجود، بما فيه صفر والسالب، يرجّع خطأ من غير تعديل. ~Entries() []MissionEntry~ لقطة مستقلة بترتيب ID تصاعدي؛ الفاضي طوله صفر.

غطّي الفاضي، trimming، الإضافة، الختم المتكرر، حذف الأول والوسط والآخر، ID بعد الحذف، الأخطاء من غير تغيير، واستقلال لقطة العرض. لازم تمسك الأربع mutants: إعادة ID، ختم ضايع، حذف غلط، ولقطة مشتركة. مفيش Report ولا تزامن ولا I/O في العقد.
`, Hints: []string{"اعمل جدول لحذف أول ووسط وآخر سجل، مع سجل جديد لكل حالة.", "بعد حذف ID 3، الإضافة اللي بعدها لازم تطلع 4، مش 3.", "افحص الختم والقيم والترتيب مش العدد بس؛ عدّل لقطة Entries وتأكد إن الأصل ثابت."},
		En: &LevelText{Title: "Mission Log", Goal: "Test a flight log with stable identities and order", Files: []LabFile{{Name: "main.go", Starter: missionLogImplementation, ReadOnly: true}, {Name: "main_test.go", Starter: missionLogTestStarter}},
			Lesson: `# A Log That Keeps Its Identity

Mission control needs an in-memory flight log. Its implementation is supplied and correct. Write table-driven tests proving that a flight's identity is independent of its current slice position. Deleting a flight must never recycle its ID.

## Before you write
- Which operation sequence distinguishes an ID counter from the number of remaining records?
- How will you prove that errors and edits to a returned snapshot leave the log intact?

## A table of scenarios instead of duplicated tests
~~~go
func TestSizes(t *testing.T) {
 cases := []struct { name string; input string; want int }{
  {"empty", "", 0}, {"two bytes", "ab", 2},
 }
 for _,tc := range cases {
  t.Run(tc.name,func(t *testing.T) {
   if len(tc.input)!=tc.want {t.Fatalf("length = %d",len(tc.input))}
  })
 }
}
~~~
~testing.T~, ~t.Run~, and ~t.Fatalf~ come from chapter 16. The table separates varying inputs from shared assertions. For the log, table-drive deletion of the first, middle, and last flights. Declare a fresh ~var log MissionLog~ inside each subtest so state cannot leak.

## CRUD and identity policy
CRUD means create, read, update, and delete here, without a terminal command loop. ~MissionEntry~ contains ~ID int~, ~Sector string~, and ~Sealed bool~. A zero-value ~MissionLog~ is ready to use; its counter is separate from its entries.
~~~go
var log MissionLog
id,err:=log.Log(" ice ")
err=log.Seal(id)
entries:=log.Entries()
fmt.Println(entries) // [{1 ice true}]
err=log.Erase(id)
fmt.Println(err) // <nil>
~~~
~Log~ trims surrounding whitespace using ~strings.TrimSpace~ and returns an ID and error. ~Seal~ seals the flight; ~Erase~ removes it. ~Entries~ returns an independent snapshot in creation order. This demonstrates the contract, rather than implementing the required tests.

## Compare state before and after
~~~go
before:=log.Entries()
err:=log.Seal(999)
if err==nil || !reflect.DeepEqual(before,log.Entries()) {
 t.Fatal("failed operation changed state")
}
~~~
Import ~reflect~. ~reflect.DeepEqual~ compares struct fields and slice elements in order; it distinguishes nil from an empty allocated slice. The contract accepts either empty representation, so check ~len~ there. Take the snapshot before the operation being tested.

## Common mistakes
~~~go
log.Seal(1) // Passes even if sealing does nothing: no assertion
if len(log.Entries())==2 { return } // Count does not prove IDs or order
~~~
Mutants compile but violate the contract. Your assertions must fail on them and pass on the correct implementation. Checking errors against nil is enough here; their wording is outside the contract.

## Your mission
Edit only ~main_test.go~; ~main.go~ is read-only. Write table-driven tests with independent cases for ~MissionLog~. ~Log(sector string) (int,error)~ appends a trimmed sector; blank input returns 0,error without consuming an ID. IDs start at 1 and advance by one per successful addition, never recycling after deletion.

~Seal(id int) error~ seals an existing entry; repeated sealing succeeds without further change. ~Erase(id int) error~ removes only that entry, retaining survivor order. Every absent ID, including zero and negatives, returns an error without mutation. ~Entries() []MissionEntry~ returns an independent snapshot in ascending ID order; an empty log has length zero.

Cover empty state, trimming, addition, repeated sealing, head/middle/tail deletion, IDs after deletion, unchanged state on errors, and independent snapshots. Kill all four mutants: recycled ID, lost seal, wrong deletion, and shared snapshot. There is no Report, concurrency, or I/O in this contract.
`, Hints: []string{"Table-drive head, middle, and tail deletion with a fresh log per case.", "After deleting ID 3, the next addition must issue 4, not 3.", "Assert seals, values, and order rather than count alone; edit an Entries snapshot and verify the log is intact."}},
		Lab: &Lab{ReqFile: "main_test.go", Hidden: missionLogHidden, Mutants: []Mutant{
			{Desc: i18n.T("إعادة استخدام ID", "Recycled ID"), Code: mutate(missionLogImplementation, "m.issued++", "m.issued = len(m.entries)+1")},
			{Desc: i18n.T("الختم ضاع", "Lost seal"), Code: mutate(missionLogImplementation, "m.entries[i].Sealed = true", "m.entries[i].Sealed = false")},
			{Desc: i18n.T("حذف غلط", "Wrong deletion"), Code: mutate(missionLogImplementation, "m.entries = append(m.entries[:i], m.entries[i+1:]...)", "m.entries = m.entries[:i]")},
			{Desc: i18n.T("لقطة مشتركة", "Shared snapshot"), Code: mutate(missionLogImplementation, "return append([]MissionEntry(nil), m.entries...)", "return m.entries")},
		}},
	})
	add(&Level{ID: "18-2-black-box", Chapter: 18, Kind: "lab", Title: "الصندوق الأسود", Goal: "احفظ واقرأ بيانات الصندوق الأسود بصيغة JSON من غير ضياع دليل العطل",
		Requires: []string{"sel:json.Unmarshal"}, Bridge: []int{17}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: blackBoxStarter}},
		Lesson: `# الصندوق الأسود مش مجرد نص

الطلعات خلصت، وعايزين نحفظ أسماء أجهزة النداء مع رقم دورة التشغيل. JSON صيغة بيانات ليها أنواع وأسماء حقول؛ النص اللي شكله قريب منها مش بالضرورة بيانات سليمة.

## قبل ما تكتب
- إزاي تفرّق بين صندوق فاضي صالح وملف ناقص أو schema ناقصة؟
- لو القراية فشلت، هتسيب دليل العطل على القرص من غير تغيير إزاي؟

## struct tags وأسماء السلك
~~~go
type Battery struct {
 Charge int ` + "`" + `json:"charge"` + "`" + `
}
encoded,err:=json.Marshal(Battery{Charge:7})
fmt.Println(string(encoded),err) // {"charge":7} <nil>
~~~
استورد ~encoding/json~. ~json.Marshal~ تحوّل قيمة لـ ~[]byte~ وخطأ. الحقل لازم يبدأ بحرف كبير عشان يتصدّر؛ الـ struct tag، النص بين backticks بعد النوع، يحدد اسم JSON. الـ tag مش اسم المتغير في Go. الصندوق هيستخدم ~cycle~ و~callsigns~ بنفس الفكرة.

## فك الترميز مش تحقق من المعنى
~~~go
var battery Battery
err:=json.Unmarshal([]byte("{\"charge\":7}"),&battery)
fmt.Println(battery.Charge,err) // 7 <nil>
~~~
~json.Unmarshal~ بتحتاج pointer عشان تملا القيمة. بترفض JSON بايظة، نوع غلط، أو مستند تاني بعد الأول. لكنها بتسيب الحقول الناقصة بقيمها الصفرية وبتتجاهل الحقول المجهولة. فراجع المعنى بعد فك الترميز: الموجب مطلوب، وnil غير slice فاضية موجودة فعلاً.

~~~go
var absent []string
present:=[]string{}
fmt.Println(absent==nil,present==nil) // true false
~~~
JSON ~null~ أو حقل callsigns ناقص بيدّي nil؛ ~[]~ بتدي slice موجودة وفاضية. هنرفض الأولى ونقبل التانية. أسماء النداء نفسها ممكن تحتوي عربي وemoji وquotes؛ المكتبة بتعمل escaping، متبنِيش JSON بإيدك.

## ملف محلي وخطأ بسياق
~~~go
err:=os.WriteFile(path,encoded,0600)
raw,err:=os.ReadFile(path)
if err!=nil {return fmt.Errorf("read %s: %w",path,err)}
~~~
استورد ~os~ و~fmt~. ~os.WriteFile~ تكتب أو تستبدل الملف؛ ~0600~ صلاحيات مالك الملف عند الإنشاء. ~os.ReadFile~ تقرأ البايتات كلها وتقفل الملف داخلياً. ~%w~ من 5-3 بيحتفظ بسبب الخطأ؛ ضيف مسار الملف ومرحلة الفشل. متستخدمش القراءة عشان تصلّح أو تمسح ملف بايظ.

## جرّب العقد
~~~go
box:=BlackBox{Cycle:3,Callsigns:[]string{"relay"}}
err:=WriteBlackBox(path,box)
if err!=nil {return err}
loaded,err:=ReadBlackBox(path)
fmt.Println(loaded.Cycle,err) // 3 <nil>
~~~
~path~ مسار ملف جوه مجلد موجود. المثال ده بيستخدم الدالتين بعد تنفيذهم، وبيوضح إن خطأ الحفظ لازم يتفحص قبل القراية.

## أخطاء شائعة
~~~go
json.Unmarshal(raw,box) // من غير &: مش pointer
os.WriteFile(path,[]byte("{}"),0600) // مسح دليل العطل كعلاج لخطأ قراية
~~~
افحص كل خطأ قبل استخدام الناتج. افحص schema قبل الكتابة عشان حفظ مدخل غير صالح ميمسحش النسخة الموجودة.

## المهمة
نفّذ ~WriteBlackBox(path string, box BlackBox) error~ و~ReadBlackBox(path string) (BlackBox,error)~. ~BlackBox~ فيها ~Cycle int~ بـ tag ~json:"cycle"~ و~Callsigns []string~ بـ tag ~json:"callsigns"~. Cycle لازم موجب؛ Callsigns لازم مش nil، ومسموح فاضية. كل اسم لازم TrimSpace بتاعه مش فاضي، بس احفظ الاسم الأصلي بالمسافات والترتيب والتكرار.

الكتابة تتحقق الأول ثم تحفظ JSON بالحقلين الاتنين؛ القراية تفك JSON وتتحقق وترجع نفس القيم. الحقول المجهولة تتجاهلها. ملف ناقص، فاضي، بايظ، نوع غلط، schema ناقصة أو غير صالحة: رجّع ~BlackBox{}~ وخطأ واضح فيه المسار ومرحلة read أو decode أو validate؛ أخطاء نظام الملفات تحتفظ بسببها بـ %w. القراية عمرها ما تكتب. فشل validation قبل الحفظ يسيب الملف زي ما هو؛ فشل الكتابة يرجع خطأ. مش مطلوب atomic rename أو استرجاع بعد انقطاع كهربا. مفيش Report؛ التستات بتستخدم ملفات مؤقتة.
`, Hints: []string{"ابدأ بالـ tags وبعدها Marshal وUnmarshal؛ افصل التحقق عن القراءة والكتابة.", "Cycle>0 وCallsigns!=nil بيفرقوا schema ناقصة عن قائمة فاضية صالحة.", "تحقق قبل WriteFile؛ لف أخطاء القراية وفك الترميز بالمسار و%w، ومتكتبش من ReadBlackBox."},
		En: &LevelText{Title: "Black Box", Goal: "Save and load black-box JSON without destroying failure evidence", Files: []LabFile{{Name: "main.go", Starter: blackBoxStarter}},
			Lesson: `# The Black Box Is More Than Text

The flights have ended. Preserve device callsigns together with the operating cycle. JSON has data types and field names; text that merely resembles JSON is not necessarily valid data.

## Before you write
- How will you distinguish a valid empty box from a missing file or incomplete schema?
- How will a failed read preserve the evidence on disk?

## Struct tags and wire names
~~~go
type Battery struct {
 Charge int ` + "`" + `json:"charge"` + "`" + `
}
encoded,err:=json.Marshal(Battery{Charge:7})
fmt.Println(string(encoded),err) // {"charge":7} <nil>
~~~
Import ~encoding/json~. ~json.Marshal~ converts a value to ~[]byte~ plus an error. Fields must start with an uppercase letter to be exported. A struct tag, the text in backticks after the type, specifies the JSON name. It does not rename the Go field. The box uses ~cycle~ and ~callsigns~ the same way.

## Decoding does not validate meaning
~~~go
var battery Battery
err:=json.Unmarshal([]byte("{\"charge\":7}"),&battery)
fmt.Println(battery.Charge,err) // 7 <nil>
~~~
~json.Unmarshal~ requires a pointer to fill the value. It rejects broken JSON, wrong types, and a second document after the first. Missing fields retain zero values and unknown fields are ignored. Validate meaning after decoding: a positive value is required, and nil differs from a present empty slice.

~~~go
var absent []string
present:=[]string{}
fmt.Println(absent==nil,present==nil) // true false
~~~
JSON ~null~ or an absent callsigns field produces nil; ~[]~ produces a present empty slice. Reject the former and accept the latter. Callsigns may contain Unicode, emoji, or quotes. Let the library escape them instead of assembling JSON by hand.

## A local file and contextual errors
~~~go
err:=os.WriteFile(path,encoded,0600)
raw,err:=os.ReadFile(path)
if err!=nil {return fmt.Errorf("read %s: %w",path,err)}
~~~
Import ~os~ and ~fmt~. ~os.WriteFile~ writes or replaces a file; ~0600~ specifies owner permissions when creating it. ~os.ReadFile~ reads all bytes and closes the file internally. ~%w~ from 5-3 preserves the cause; include the path and failing stage. A read must never repair or erase a corrupt file.

## Calling the contract
~~~go
box:=BlackBox{Cycle:3,Callsigns:[]string{"relay"}}
err:=WriteBlackBox(path,box)
if err!=nil {return err}
loaded,err:=ReadBlackBox(path)
fmt.Println(loaded.Cycle,err) // 3 <nil>
~~~
~path~ is a file path inside an existing directory. This uses the completed functions and shows that save errors must be checked before reading.

## Common mistakes
~~~go
json.Unmarshal(raw,box) // Missing &: not a pointer
os.WriteFile(path,[]byte("{}"),0600) // Erases evidence after a read failure
~~~
Check every error before using its output. Validate the schema before writing so an invalid save cannot erase the existing file.

## Your mission
Implement ~WriteBlackBox(path string, box BlackBox) error~ and ~ReadBlackBox(path string) (BlackBox,error)~. ~BlackBox~ contains ~Cycle int~ tagged ~json:"cycle"~ and ~Callsigns []string~ tagged ~json:"callsigns"~. Cycle must be positive. Callsigns must be non-nil but may be empty. Every callsign must be nonblank after TrimSpace; preserve the original strings, spaces, order, and duplicates.

Validate before writing JSON with both fields. Reading decodes and validates, returning the same values. Ignore unknown fields. Missing, empty, malformed, wrongly typed, incomplete, or invalid files return ~BlackBox{}~ and a clear error containing the path and read, decode, or validate stage. Preserve filesystem error causes with %w. Reading never writes. Failed save validation leaves the file unchanged; write failures return errors. Atomic rename and power-loss recovery are outside this contract. No Report is needed; tests use temporary files.
`, Hints: []string{"Start with tags, then Marshal and Unmarshal; separate validation from reading and writing.", "Cycle>0 and Callsigns!=nil distinguish missing schema from a valid empty list.", "Validate before WriteFile; wrap read/decode failures with path and %w, and never write from ReadBlackBox."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: blackBoxHidden},
	})
	add(&Level{ID: "18-3-cargo-manifest", Chapter: 18, Kind: "lab", Title: "بيان الشحنة", Goal: "اقرأ CSV الشحنة واجمع الكميات من غير بلع صف بايظ",
		Requires: []string{"sel:csv.NewReader"}, Bridge: []int{18}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: cargoManifestStarter}},
		Lesson: `# الصناديق اللي جواها فاصلة

عامل الرصيف بيبعت بيان شحنة CSV. اسم جهاز النداء ممكن نفسه يحتوي فاصلة أو quote، فـ Split على فاصلة مش parser. المكتبة هتفك الصيغة، وإنت مسؤول عن schema والإجماليات.

## قبل ما تكتب
- هتعرف الفرق بين ملف فاضي وبيان فيه header بس إزاي؟
- لو آخر صف بايظ بعد شحنة سليمة، هل ترجع إجمالي ناقص ولا تفشل البيان كله؟

## CSV مش سطور متقسمة بإيدك
~~~go
reader:=csv.NewReader(strings.NewReader("bay,units\n\"west, upper\",4\n"))
header,err:=reader.Read()
row,err:=reader.Read()
fmt.Println(header,row,err) // [bay units] [west, upper 4] <nil>
~~~
استورد ~encoding/csv~. ~csv.NewReader~ بتاخد ~io.Reader~ اللي اتشرح في 9-4؛ ~Read~ ترجع record كـ ~[]string~ وخطأ. أول record بيحدد العدد المتوقع للحقول تلقائياً. quotes بتحافظ على فاصلة جوه الحقل؛ quote جوه quoted field بتتكتب مرتين، والـ field ممكن يحتوي newline. المكتبة كمان بتتعامل مع CRLF وبتتجاهل السطور الفاضية تماماً.

## EOF مش صف بايظ
~~~go
for {
 row,err:=reader.Read()
 if err==io.EOF {break}
 if err!=nil {return err}
 fmt.Println(row)
}
~~~
~io.EOF~ معناها نهاية البيانات الطبيعية بعد header. قبل header هي خطأ لأن schema مش موجودة. متستخدمش ~ReadAll~ هنا: القراءة record بـ record تخليك تجمع الإجماليات من غير الاحتفاظ بكل الشحنة.

## تحقق بعد فك الصيغة
~~~go
units,err:=strconv.Atoi("12")
name:=strings.TrimSpace(" west ")
fmt.Println(units,name,err) // 12 west <nil>
~~~
~strconv.Atoi~ من 4-3 بتقرأ int، مش floating point: NaN وInf و1e3 مش كميات صحيحة في عقدنا. افحص الحدود كمان؛ 0 شحنة صحيحة. trim اسم النداء بس، مش header أو الحقول الرقمية.

~~~go
totals:=map[string]int{}
totals[name]+=units
err=fmt.Errorf("record %d: %w",3,err)
~~~
الـ map بتجمع كل اسم، حتى لو ظهر أكتر من مرة. عدّ record من 1 للـ header؛ quoted newline مش record جديدة. أخطاء parsing والقراية تاخد سياق record وتحافظ على السبب بـ %w. مثال Errorf يتستخدم بس لما err مش nil.

## شكل الناتج
~~~go
totals,err:=ReadManifest(strings.NewReader("id,callsign,cargo\n0,relay,4\n"))
fmt.Println(totals.Rows,totals.Units,totals.ByCallsign["relay"],err) // 1 4 4 <nil>
~~~
ده نداء للدالة المطلوبة؛ ~CargoTotals~ بتجمع عدد records وإجمالي الوحدات وتوزيعها، مش مجرد عدد السطور النصية.

## أخطاء شائعة
~~~go
fields:=strings.Split(line,",") // بيكسر اسم زي "west, upper"
units,_:=strconv.Atoi(fields[2]) // النص البايظ يتحسب صفر بالغلط
~~~
متحدّثش التقرير النهائي لو صف فشل؛ رجّع zero value وخطأ بدل إجمالي شكله ناجح. كرر القراية لحد EOF، مش لحد أول سطر فاضي.

## المهمة
نفّذ ~ReadManifest(in io.Reader) (CargoTotals,error)~. ~CargoTotals~ فيها ~Rows,Units int~ و~ByCallsign map[string]int~. الـ header لازم بالظبط ~id,callsign,cargo~ بالحالة والترتيب دول. الملف الفاضي خطأ فيه كلمة header؛ header بس نجاح بإجماليات صفر وmap متجهزة.

كل record بيانات لازم 3 حقول: id يقبله Atoi وغير سالب؛ callsign بعد TrimSpace مش فاضي؛ cargo يقبله Atoi من 0 لـ1000000 شامل. تكرار IDs مسموح لأنها أرقام محلية لمصدر الشحنة. اجمع عدد records في Rows، وكل cargo في Units، ولكل اسم في ByCallsign. المدخل الصالح أقصى 100000 record بيانات عشان مجموع int يفضل آمن. متخزّنش الصفوف كلها.

أي schema أو quoting أو field count أو قيمة بايظة، وأي خطأ قراية، يرجّع ~CargoTotals{}~ وخطأ. خطأ البيانات يحتوي ~record N~ بالرقم المنطقي بداية من 2؛ احتفظ بأسباب أخطاء المكتبة. اقبل quoted commas وescaped quotes وCRLF وآخر record من غير newline؛ اتجاهل السطور الفاضية زي csv.Reader. النجمة التالتة: quoted field متعدد السطور، حتى لو المصدر بيبعت بايت واحد في كل Read. مفيش Report ولا ملفات قرص.
`, Hints: []string{"اقرأ header بـ csv.Reader قبل الحلقة، وافحصه بالظبط.", "Read بترجع EOF للنهاية؛ أي خطأ تاني يبطل التقرير كله.", "اجمع حسب TrimSpace للاسم؛ خلي رقم record يزيد لكل Read ناجحة مش لكل newline."},
		En: &LevelText{Title: "Cargo Manifest", Goal: "Read cargo CSV and aggregate units without swallowing broken records", Files: []LabFile{{Name: "main.go", Starter: cargoManifestStarter}},
			Lesson: `# Crates With Commas Inside

The dock sends a CSV cargo manifest. A callsign can contain a comma or quote, so splitting on commas is not a parser. The library handles syntax; you own the schema and totals.

## Before you write
- How will you distinguish an empty file from a manifest containing only its header?
- If the final record is broken after valid cargo, should you return an incomplete total or fail the entire manifest?

## CSV is not manually split lines
~~~go
reader:=csv.NewReader(strings.NewReader("bay,units\n\"west, upper\",4\n"))
header,err:=reader.Read()
row,err:=reader.Read()
fmt.Println(header,row,err) // [bay units] [west, upper 4] <nil>
~~~
Import ~encoding/csv~. ~csv.NewReader~ accepts the ~io.Reader~ taught in 9-4. ~Read~ returns a record as ~[]string~ and an error. The first record automatically sets the expected field count. Quotes preserve embedded commas; an embedded quote is doubled inside a quoted field, and a field may contain a newline. The library also handles CRLF and skips completely empty lines.

## EOF is not a broken record
~~~go
for {
 row,err:=reader.Read()
 if err==io.EOF {break}
 if err!=nil {return err}
 fmt.Println(row)
}
~~~
~io.EOF~ means normal completion after the header. Before the header it is an error because the schema is absent. Avoid ~ReadAll~ here: reading one record at a time lets you aggregate without retaining the whole shipment.

## Validate after decoding syntax
~~~go
units,err:=strconv.Atoi("12")
name:=strings.TrimSpace(" west ")
fmt.Println(units,name,err) // 12 west <nil>
~~~
~strconv.Atoi~ from 4-3 parses an int, not floating point. NaN, Inf, and 1e3 are not integer quantities in this contract. Check bounds as well; zero is valid. Trim only the callsign, not the header or numeric fields.

~~~go
totals:=map[string]int{}
totals[name]+=units
err=fmt.Errorf("record %d: %w",3,err)
~~~
The map aggregates each callsign, including repeated appearances. Count records from header number 1: a quoted newline does not start another record. Add record context to parsing and read errors while preserving their causes with %w. Use the Errorf example only when err is non-nil.

## Result shape
~~~go
totals,err:=ReadManifest(strings.NewReader("id,callsign,cargo\n0,relay,4\n"))
fmt.Println(totals.Rows,totals.Units,totals.ByCallsign["relay"],err) // 1 4 4 <nil>
~~~
This calls the required function. ~CargoTotals~ holds record count, total units, and their distribution, rather than counting physical lines.

## Common mistakes
~~~go
fields:=strings.Split(line,",") // Breaks a name such as "west, upper"
units,_:=strconv.Atoi(fields[2]) // Silently counts invalid text as zero
~~~
Do not publish a final report after a failed record. Return the zero value and an error rather than plausible partial totals. Read until EOF, not until the first empty line.

## Your mission
Implement ~ReadManifest(in io.Reader) (CargoTotals,error)~. ~CargoTotals~ contains ~Rows,Units int~ and ~ByCallsign map[string]int~. Require exactly the header ~id,callsign,cargo~, with that case and order. Empty input returns an error mentioning header. Header-only input succeeds with zero totals and an initialized map.

Every data record has three fields: id accepted by Atoi and nonnegative; a nonblank trimmed callsign; and cargo accepted by Atoi in the inclusive range 0 through 1000000. Duplicate IDs are allowed because they are local numbers from shipment sources. Count records in Rows, all cargo in Units, and cargo per normalized callsign in ByCallsign. Valid input contains at most 100000 data records so int totals remain safe. Do not retain all rows.

Any schema, quoting, field-count, value, or read failure returns ~CargoTotals{}~ and an error. Data errors contain ~record N~, counting logical records from 2; preserve library error causes. Accept quoted commas, escaped quotes, CRLF, and a final unterminated record. Skip completely empty lines as csv.Reader does. The third star covers multiline quoted fields, even when input arrives one byte per Read. No Report or disk files are involved.
`, Hints: []string{"Read and check the exact header using csv.Reader before looping.", "Read returns EOF for completion; every other error invalidates the entire report.", "Aggregate by trimmed callsign; increment record number per successful Read, not per newline."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: cargoManifestHidden, Bonus: true},
	})
	add(&Level{ID: "18-4-signal-words", Chapter: 18, Kind: "lab", Title: "كلمات الإشارة", Goal: "رتّب كلمات الراديو بالتكرار مع Unicode وتعادل حتمي",
		Requires: []string{"map", "sel:bufio.Scanner"}, Bridge: []int{19}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: signalWordsStarter}},
		Lesson: `# الإشارة المتكررة تستاهل اهتمام

هوائي المستعمرة استقبل نصوص كتير. عايزين أكتر الكلمات تكراراً، مش أول كلمات وصلت. لازم النتيجة تفضل ثابتة حتى لما map تلف بترتيب مختلف.

## قبل ما تكتب
- إيه اللي يفصل كلمتين: مسافة بس ولا ترقيم وemoji كمان؟
- لو كلمتين ليهم نفس العدد، إيه ترتيب التعادل اللي يخلي التقرير حتمي؟

## تطبيع Unicode وتقسيم الكلمات
~~~go
lower:=strings.ToLower("ÉCHO, Orbit42!")
words:=strings.FieldsFunc(lower,func(r rune) bool {
 return !unicode.IsLetter(r) && !unicode.IsDigit(r)
})
fmt.Println(words) // [écho orbit42]
~~~
استورد ~strings~ و~unicode~. ~strings.ToLower~ توحّد حالة الحروف حسب Unicode. ~strings.FieldsFunc~ بتاخد predicate على rune؛ true معناها فاصل. ~unicode.IsLetter~ و~unicode.IsDigit~ بيميزوا الحروف والأرقام العربية وغيرها، مش ASCII بس. ~abc123~ كلمة واحدة، و~well-known~ كلمتين؛ apostrophe وemoji فواصل. مفيش Unicode normalization أو دمج combining marks مطلوب.

## Scanner بحد معلن
~~~go
var scanner *bufio.Scanner=bufio.NewScanner(in)
scanner.Buffer(make([]byte,4096),(1<<20)+2)
for scanner.Scan() {
 line:=scanner.Text()
 if len(line)>1<<20 {return fmt.Errorf("line too long")}
 fmt.Println(line)
}
return scanner.Err()
~~~
~bufio.NewScanner~ و~Scan~ و~Text~ و~Err~ من 9-4. هنا بنكتب النوع صراحة: ~*bufio.Scanner~ هو pointer للماسح اللي NewScanner بترجّعه؛ استخدم نفس التصريح في حلك. الجديد ~Scanner.Buffer~ قبل أول Scan: buffer مبدئي، ثم أكبر buffer مسموح. ~1<<20~ معناها 1048576 بايت؛ عامل الإزاحة ~<<~ بيضرب في قوة اتنين. البايتين الزيادة يسمحوا بـ CRLF عند الحد، وبعدين افحص طول المحتوى بنفسك. Err غير nil معناها فشل، مش EOF ناجح.

القراءة سطر بسطر بتقبل نص كبير من غير تحميله كله. خزن قاموس الكلمات بس والنتيجة؛ الكلمة مبتعدّيش newline لأنها فاصل. حدود الدرس بالبايت عشان الذاكرة، مش عدد runes.

## عدّ وبعدين رتّب
~~~go
counts:=map[string]int{}
for _,word:=range words {counts[word]++}
fmt.Println(counts["missing"]) // 0
fmt.Println(strings.Compare("a","b")) // -1
~~~
القيمة الصفرية لعنصر map مش موجود بتخلي الزيادة مباشرة. ~strings.Compare~ ترجع سالب أو صفر أو موجب حسب ترتيب string في Go. حوّل العدادات لـ slice، واستخدم ~slices.SortFunc~ من 6-4: العدد تنازلي، وبعدين الكلمة تصاعدي. بعد الترتيب بس اقطع أعلى limit.

## شكل الاستدعاء
~~~go
ranked,err:=TopSignals(strings.NewReader("ice relay ICE"),1)
fmt.Println(ranked,err) // [{ice 2}] <nil>
example:=WordCount{Word:"relay",Count:1}
fmt.Println(example.Word,example.Count) // relay 1
~~~
~WordCount~ عنصر في النتيجة، مش خريطة العد نفسها. limit أقصى عدد نتائج، مش أقصى عدد كلمات تقراها.

## أخطاء شائعة
~~~go
for word,count:=range counts {fmt.Println(word,count)} // ترتيب map مش مضمون
strings.Fields("ice,ice") // كلمة واحدة فيها فاصلة، مش كلمتين
~~~
متعدّش ترقيم جوه الكلمة، ومتقصش قبل ترتيب القاموس كله. لو القراية فشلت ارجع nil,error بدل نتائج النص اللي وصل بس.

## المهمة
نفّذ ~TopSignals(in io.Reader, limit int) ([]WordCount,error)~؛ ~WordCount~ فيها ~Word string~ و~Count int~. الكلمة سلسلة متصلة من Unicode letters أو digits؛ كل حاجة تانية فاصل، بما فيها hyphen وapostrophe وemoji. طبّق ToLower؛ عدّ التكرار، ورتّب Count تنازلي ثم Word تصاعدي بترتيب Go، ورجّع أول limit أو كل المتاح لو أقل.

الفاضي أو الترقيم بس يرجّع طول صفر,nil. limit سالب يرجّع nil,error؛ صفر يرجّع nil,nil من غير قراية المصدر. كل سطر محتواه لحد 1 MiB شامل، من غير LF أو CRLF؛ اقبل عند الحد حتى من غير newline. الأكبر أو أي خطأ قراية يرجّع nil,error، من غير تقرير جزئي أو truncation. استخدم Scanner وmap؛ النص ممكن كبير ومتقسم على Read صغيرة. مفيش Report.
`, Hints: []string{"اقرأ بسطر Scanner؛ FieldsFunc تفصل على أي rune مش حرف ولا رقم.", "ToLower قبل العد؛ حوّل map لـ slice ورتّب العدد ثم الكلمة.", "وسّع Buffer قبل Scan مع مساحة نهاية السطر وافحص الطول وErr؛ فشل المصدر يرجع nil."},
		En: &LevelText{Title: "Signal Words", Goal: "Rank radio words by frequency with Unicode and deterministic ties", Files: []LabFile{{Name: "main.go", Starter: signalWordsStarter}},
			Lesson: `# Repeated Signals Deserve Attention

The colony antenna received a large text stream. Find the most frequent words, rather than the earliest words received. Results must remain stable even when map iteration order changes.

## Before you write
- What separates words: only whitespace, or punctuation and emoji too?
- When counts tie, which rule makes report order deterministic?

## Unicode normalization of case and tokenization
~~~go
lower:=strings.ToLower("ÉCHO, Orbit42!")
words:=strings.FieldsFunc(lower,func(r rune) bool {
 return !unicode.IsLetter(r) && !unicode.IsDigit(r)
})
fmt.Println(words) // [écho orbit42]
~~~
Import ~strings~ and ~unicode~. ~strings.ToLower~ lowercases according to Unicode. ~strings.FieldsFunc~ accepts a rune predicate where true means separator. ~unicode.IsLetter~ and ~unicode.IsDigit~ recognize letters and digits beyond ASCII. ~abc123~ is one word; ~well-known~ is two. Apostrophes and emoji separate words. Unicode normalization and combining-mark composition are outside the contract.

## Scanner with an explicit limit
~~~go
var scanner *bufio.Scanner=bufio.NewScanner(in)
scanner.Buffer(make([]byte,4096),(1<<20)+2)
for scanner.Scan() {
 line:=scanner.Text()
 if len(line)>1<<20 {return fmt.Errorf("line too long")}
 fmt.Println(line)
}
return scanner.Err()
~~~
~bufio.NewScanner~, ~Scan~, ~Text~, and ~Err~ come from 9-4. Here the declaration explicitly names ~*bufio.Scanner~, the pointer type returned by NewScanner; use the same explicit declaration in your solution. The new ~Scanner.Buffer~ call must precede scanning: it supplies an initial buffer and maximum buffer size. ~1<<20~ is 1048576 bytes; the ~<<~ shift multiplies by a power of two. Two extra bytes accommodate CRLF at the limit; check content length separately. Non-nil Err means failure, not successful EOF.

Reading line by line supports large text without loading it all. Retain the word dictionary and result rather than the entire input. Words cannot cross a newline because it is a separator. Limits count bytes for memory control, rather than runes.

## Count, then sort
~~~go
counts:=map[string]int{}
for _,word:=range words {counts[word]++}
fmt.Println(counts["missing"]) // 0
fmt.Println(strings.Compare("a","b")) // -1
~~~
A missing map key has the zero value, making increments straightforward. ~strings.Compare~ returns a negative, zero, or positive value using Go string order. Collect counts into a slice and use ~slices.SortFunc~ from 6-4: descending count, then ascending word. Select the first limit entries only after sorting.

## Calling the contract
~~~go
ranked,err:=TopSignals(strings.NewReader("ice relay ICE"),1)
fmt.Println(ranked,err) // [{ice 2}] <nil>
example:=WordCount{Word:"relay",Count:1}
fmt.Println(example.Word,example.Count) // relay 1
~~~
~WordCount~ is one result entry, rather than the counting map itself. The limit caps result size, not how many input words to read.

## Common mistakes
~~~go
for word,count:=range counts {fmt.Println(word,count)} // Map order is unspecified
strings.Fields("ice,ice") // One chunk containing a comma, not two words
~~~
Do not retain punctuation inside words or truncate before sorting the full dictionary. On read failure return nil,error instead of reporting only the text received so far.

## Your mission
Implement ~TopSignals(in io.Reader, limit int) ([]WordCount,error)~. ~WordCount~ contains ~Word string~ and ~Count int~. Words are maximal runs of Unicode letters or digits. Everything else separates them, including hyphens, apostrophes, and emoji. Apply ToLower, count occurrences, sort by descending Count then ascending Word using Go string order, and return at most limit entries.

Empty or punctuation-only input returns length zero,nil. Negative limit returns nil,error. Zero limit returns nil,nil without reading input. Each line may contain at most 1 MiB inclusive, excluding LF or CRLF; accept the exact boundary with or without a terminator. Longer lines and any read failure return nil,error, without partial results or truncation. Use Scanner and a map; large input may arrive across tiny Reads. No Report is needed.
`, Hints: []string{"Scan lines; FieldsFunc separates on every rune that is neither a letter nor a digit.", "Lowercase before counting, collect the map into a slice, and sort by count then word.", "Increase Buffer before Scan with terminator headroom; check content length and Err, returning nil on source failure."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: signalWordsHidden},
	})
	add(&Level{ID: "18-5-flight-recorder", Chapter: 18, Kind: "lab", Title: "مسجّل الرحلة", Goal: "حلّل سيل الرحلة وميّز السطر البايظ عن فشل النقل",
		Requires: []string{"sel:bufio.Scanner", "struct"}, Bridge: []int{21}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: flightRecorderStarter}},
		Lesson: `# مسجّل الرحلة بيسمع لحد الآخر

بيانات الرحلة بتوصل كسيل سطور، بنفس شكل تليمتري المستعمرة. السطر البايظ محتاج يتسجل مكانه ونكمّل؛ السلك اللي اتقطع محتاج يوقف التقرير كله. دي سياستين مختلفتين للفشل.

## قبل ما تكتب
- أنهي فشل يخلّيك تكمّل مع diagnostic، وأنهي فشل يبطل الملخص كله؟
- هتضمن إن أرقام السطور تشمل الفاضي والبايظ من غير ما تعتمد على رقم الرسالة إزاي؟

## افصل أول حقلين وسيب الرسالة
~~~go
parts:=strings.SplitN("8 WARN low fuel near dock"," ",3)
fmt.Println(parts) // [8 WARN low fuel near dock]
number,err:=strconv.Atoi(parts[0])
fmt.Println(number,err) // 8 <nil>
~~~
~strings.SplitN~ جديدة هنا: بتقسم على separator لعدد أجزاء أقصاه n؛ الجزء الأخير يحتفظ بباقي النص. مش زي Fields اللي بتدمج المسافات. افحص ~len(parts)~ قبل indexing. ~strconv.Atoi~ من 4-3 بتقبل أعداد بإشارة اختيارية؛ هنشترط الناتج غير سالب. ~strings.TrimSpace~ من 9-4 تفحص رسالة فيها محتوى.

## مستويات محصورة وdiagnostics مرتبة
~~~go
type PacketSummary struct { Good int; Rejected []int }
var summary PacketSummary
switch level {
case "WARN": summary.Good++
default: summary.Rejected=append(summary.Rejected,lineNumber)
}
~~~
ده مثال ملخص حزم، مش عقد المهمة. الـ struct بتجمع counters وslice أرقام السطور المرفوضة. append أثناء القراية يحافظ على ترتيب وصول diagnostics؛ مفيش map مطلوب. الـ level حساس لحالة الحروف: ~warn~ مش ~WARN~.

## حدود Scanner وفشل النقل
~~~go
var scanner *bufio.Scanner=bufio.NewScanner(in)
scanner.Buffer(make([]byte,4096),(256<<10)+2)
for n:=1;scanner.Scan();n++ {
 line:=scanner.Text()
 if len(line)>256<<10 {return fmt.Errorf("line too long")}
 fmt.Println(n,line)
}
return scanner.Err()
~~~
زي 18-4، صرّح بالماسح كـ ~*bufio.Scanner~. Buffer توسّع الحد قبل Scan، و~256<<10~ تساوي 262144 بايت. خلي مساحة LF وCRLF زيادة وافحص المحتوى بالبايت؛ Text تشيل نهاية السطر. السطر الأخير من غير newline لسه سطر. افحص Err بعد الحلقة عشان المصدر ممكن يرجع بيانات وخطأ مع بعض.

الكود المرجعي بيحتفظ بالعدادات وأرقام السطور البايظة بس. الملف كله مش مطلوب في الذاكرة. الحد الهيكلي أو فشل المصدر يرجّع zero value، مش summary جزئية شكلها ناجح.

## جرّب الملخص
~~~go
summary,err:=ReadFlight(strings.NewReader("0 INFO launch\nbroken\n"))
fmt.Println(summary.Info,summary.BadLines,err) // 1 [2] <nil>
empty:=FlightSummary{}
fmt.Println(empty.Warn,empty.Error) // 0 0
~~~
BadLines أرقام مواضع في النص؛ INFO الصالحة ترفع Info بس، مش طول slice الأخطاء.

## أخطاء شائعة
~~~go
fields:=strings.Fields(line)
if len(fields)!=3 { /* الرسالة متعددة الكلمات هتترفض بالغلط */ }
for scanner.Scan() {} // ضيّع أرقام السطور وبلع scanner.Err
~~~
متستخدمش رقم التسلسل اللي جوه الرسالة كـ line number: ممكن يتكرر أو يرجع لصفر. رقم السطر مكانه في النص، بداية من 1.

## المهمة
نفّذ ~ReadFlight(in io.Reader) (FlightSummary,error)~. ~FlightSummary~ فيها ~Info,Warn,Error int~ و~BadLines []int~. السطر الصالح له 3 أجزاء باستخدام SplitN على مسافة ASCII واحدة بحد 3: رقم تسلسل يقبله Atoi وغير سالب، level واحد من INFO/WARN/ERROR بالظبط، ورسالة TrimSpace بتاعها مش فاضي. الرسالة ممكن فيها مسافات داخلية أو طرفية؛ متعملش parsing للرسالة نفسها. رقم التسلسل مش لازم فريد أو متزايد.

زوّد counter المستوى لكل سطر صالح. أي سطر تاني، بما فيه الفاضي أو leading whitespace أو level مجهول، ضيف رقمه 1-based لـ BadLines وكمّل. EOF نجاح؛ الفاضي ملخص أصفار. اقبل CRLF وآخر سطر من غير terminator.

أكبر محتوى سطر 256 KiB شامل، من غير LF أو CRLF؛ الأكبر خطأ هيكلي، مش diagnostic. أي خطأ قراية أو حد حجم يرجّع ~FlightSummary{}~,error مع الاحتفاظ بسبب خطأ المصدر. اختبارات النجاح الأساسية تستخدم سطور لحد 4 KiB؛ النجمة التالتة لقبول سطور أطول من حد Scanner الافتراضي، لحد 256 KiB بالظبط. استخدم Scanner وstruct، ومتحتفظش بكل السطور. مفيش Report.
`, Hints: []string{"عدّ السطور من 1 من غير الاعتماد على sequence؛ SplitN بحد 3 تحفظ الرسالة.", "السطر البايظ يضيف رقمه وبس، لكن scanner.Err يبطل الملخص كله.", "للنجمة التالتة وسّع Buffer لحد 256 KiB مع مساحة CRLF وافحص حجم المحتوى بنفسك."},
		En: &LevelText{Title: "Flight Recorder", Goal: "Analyze flight streams and distinguish bad records from transport failure", Files: []LabFile{{Name: "main.go", Starter: flightRecorderStarter}},
			Lesson: `# The Flight Recorder Listens to the End

Flight data arrives as a stream of lines in the colony's telemetry format. A malformed record should retain its location and allow processing to continue. A broken input wire invalidates the whole report. These are two distinct failure policies.

## Before you write
- Which failure allows a diagnostic and continued scanning, and which invalidates the entire summary?
- How will you count blank and malformed lines independently of the sequence number inside each record?

## Split the first two fields and retain the message
~~~go
parts:=strings.SplitN("8 WARN low fuel near dock"," ",3)
fmt.Println(parts) // [8 WARN low fuel near dock]
number,err:=strconv.Atoi(parts[0])
fmt.Println(number,err) // 8 <nil>
~~~
~strings.SplitN~ is new here: it splits on a separator into at most n parts, retaining the remainder in the final part. Unlike Fields, it does not collapse spaces. Check ~len(parts)~ before indexing. ~strconv.Atoi~ from 4-3 accepts an optional sign; require a nonnegative result. ~strings.TrimSpace~ from 9-4 checks that the message has content.

## Fixed levels and ordered diagnostics
~~~go
type PacketSummary struct { Good int; Rejected []int }
var summary PacketSummary
switch level {
case "WARN": summary.Good++
default: summary.Rejected=append(summary.Rejected,lineNumber)
}
~~~
This packet summary is an example, not the mission contract. A struct holds counters and rejected line numbers. Appending during scanning preserves diagnostic arrival order without a map. Levels are case sensitive: ~warn~ is different from ~WARN~.

## Scanner limits and transport failure
~~~go
var scanner *bufio.Scanner=bufio.NewScanner(in)
scanner.Buffer(make([]byte,4096),(256<<10)+2)
for n:=1;scanner.Scan();n++ {
 line:=scanner.Text()
 if len(line)>256<<10 {return fmt.Errorf("line too long")}
 fmt.Println(n,line)
}
return scanner.Err()
~~~
As in 18-4, explicitly declare the scanner as ~*bufio.Scanner~ and call Buffer before Scan to increase the limit. ~256<<10~ equals 262144 bytes. Allow extra room for LF or CRLF and check content length separately; Text strips the line ending. The final unterminated line still counts. Check Err after the loop because a reader may deliver bytes and an error together.

The reference approach retains counters and malformed line numbers, rather than all input. Size-limit and source failures return the zero value instead of a misleading partial success.

## Calling the summary
~~~go
summary,err:=ReadFlight(strings.NewReader("0 INFO launch\nbroken\n"))
fmt.Println(summary.Info,summary.BadLines,err) // 1 [2] <nil>
empty:=FlightSummary{}
fmt.Println(empty.Warn,empty.Error) // 0 0
~~~
BadLines stores positions in the input. Valid INFO records increment Info without adding diagnostics.

## Common mistakes
~~~go
fields:=strings.Fields(line)
if len(fields)!=3 { /* Incorrectly rejects a multiword message */ }
for scanner.Scan() {} // Loses line numbers and swallows scanner.Err
~~~
Do not use the embedded sequence number as the line number: it can repeat or restart at zero. A line number is its position in the input, starting at 1.

## Your mission
Implement ~ReadFlight(in io.Reader) (FlightSummary,error)~. ~FlightSummary~ contains ~Info,Warn,Error int~ and ~BadLines []int~. A valid line has three parts using SplitN on a single ASCII space with limit 3: a nonnegative sequence number accepted by Atoi, exactly INFO/WARN/ERROR, and a message nonblank after TrimSpace. Messages may have internal or surrounding whitespace; do not parse their content. Sequence numbers need not be unique or increasing.

Increment the matching counter for valid lines. For every other line, including blank lines, leading whitespace, and unknown levels, append its 1-based line number to BadLines and continue. EOF succeeds; empty input returns zero counters. Accept CRLF and a final line without a terminator.

Maximum content length is 256 KiB inclusive, excluding LF or CRLF. Longer lines are structural errors, not diagnostics. Any read or size failure returns ~FlightSummary{}~,error, preserving source error causes. Core success tests use lines up to 4 KiB. The third star covers lines exceeding Scanner's default limit, through exactly 256 KiB. Use Scanner and a struct without retaining all lines. No Report is needed.
`, Hints: []string{"Count lines from 1 independently of sequence numbers; SplitN with limit 3 retains the message.", "Malformed records only add their line number; scanner.Err invalidates the entire summary.", "For the third star, increase Buffer through 256 KiB plus CRLF headroom and check content length explicitly."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: flightRecorderHidden, Bonus: true},
	})
	add(&Level{ID: "18-6-crew-directory", Chapter: 18, Kind: "lab", Title: "دليل الطاقم", Goal: "سجّل الطاقم وابحث عنهم واحفظ عن طريق interface صغيرة",
		Requires: []string{"iface"}, Bridge: []int{22}, Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: crewDirectoryStarter}},
		Lesson: `# دليل الطاقم من غير ربطه بالقرص

مركز الاتصالات محتاج يلاقي حد بالنداء أو التخصص ويحفظ الدليل. البحث مش لازم يعرف التخزين ملف ولا ذاكرة؛ خلي الحفظ ورا interface صغيرة عند المستهلك زي مقابس 9-4.

## قبل ما تكتب
- إزاي ترفض ID مكرر من غير ما تستبدل العضو القديم؟
- لو التخزين عدّل الـ slice اللي استلمها، هتحمي حالة الدليل إزاي؟

## بحث substring من غير حساسية للحالة
~~~go
name:=strings.ToLower("ÉCHO Pilot")
query:=strings.ToLower("PIL")
fmt.Println(strings.Contains(name,query)) // true
~~~
~strings.ToLower~ من 18-4؛ ~strings.Contains~ جديدة هنا وبتفحص وجود substring، مش مساواة ولا كلمة كاملة. الـ query الفاضية موجودة في أي نص، فبترجع الكل. لازم تطبّق ToLower على الطرفين؛ مش مطلوب Unicode normalization.

## interface للحفظ عند المستهلك
~~~go
type BadgeSink interface { PutBadges([]string) error }
type MemoryBadges struct { Names []string }
func (m *MemoryBadges) PutBadges(names []string) error {
 m.Names=append([]string(nil),names...)
 return nil
}
var sink BadgeSink=&MemoryBadges{}
err:=sink.PutBadges([]string{"pilot"})
fmt.Println(err) // <nil>
~~~
أي نوع عنده method بنفس التوقيع يحقق interface من غير implements. مثال الشارات بيستعمل pointer receiver عشان يغيّر المخزن؛ الـ spread ~names...~ يضيف عناصر slice زي 8-4. عرّف ~CrewStore~ بنفس الفكرة بعقد ~SaveCrew([]Crew) error~، من غير فرض JSON أو مسار ملف على المستهلك.

## لقطة مستقلة وترتيب ثابت
~~~go
names:=[]string{"pilot"}
copyOfNames:=append([]string(nil),names...)
copyOfNames[0]="changed"
fmt.Println(names[0]) // الاسم القديم
~~~
خُد نسخة قبل ما تدي البيانات للمخزن. في المهمة Crew حقولها strings، فنسخ عناصر slice كفاية؛ مفيش pointers داخلية. استخدم ~slices.SortFunc~ من 6-4 و~strings.Compare~ من 18-4 عشان ID تصاعدي، مش ترتيب map أو ترتيب تسجيل الأعضاء.

خطأ ~SaveCrew~ لازم يوصل للمنادي، ومايغيّرش الدليل. المخزن نفسه مسؤول عن اللي يكتبه؛ المستهلك مش بيعد إنه يقدر يرجّع آثار تخزين خارجي فشل.

## استخدام الدليل
~~~go
var directory CrewDirectory
err:=directory.Register(Crew{ID:"C7",Callsign:"Echo",Role:"Pilot"})
fmt.Println(directory.Search("PIL"),err) // [{C7 Echo Pilot}] <nil>
err=directory.Save(store)
~~~
~store~ أي تنفيذ لعقد CrewStore، زي مثال MemoryBadges بس بـ Crew. Search بتعرض نسخة، وSave بتطلب الحفظ من التنفيذ المحقون من غير ما تغيّر الدليل.

## أخطاء شائعة
~~~go
if strings.Contains(crew.Callsign,strings.ToLower(query)) {} // الطرف الأول لسه case مختلفة
return store.SaveCrew(d.crew) // المخزن يقدر يبوّظ slice الدليل
~~~
افحص التكرار قبل الإضافة. الحفظ مش عملية تسجيل جديدة، ومش سبب لتغيير البيانات في الذاكرة.

## المهمة
نفّذ ~CrewDirectory~ بقيمة صفرية جاهزة. ~Crew~ فيها ~ID,Callsign,Role string~. ~Register(crew Crew) error~ ترفض أي حقل TrimSpace بتاعه فاضي، أو ID موجود بالظبط، من غير تغيير؛ احفظ الحقول الأصلية من غير trimming. الـ ID حساس لحالة الحروف والمسافات، ومش بيتولد تلقائي. كل دليل مستقل.

~Search(query string) []Crew~ ترجع اللي callsign أو role فيهم query بعد ToLower للطرفين. متبحثش في ID، ومتعملش TrimSpace للـ query. الفاضية ترجع الكل، والمفيش تطابق طول صفر. النتائج نسخ مستقلة مرتبة بـ ID تصاعدي بترتيب Go.

عرّف ~CrewStore interface { SaveCrew([]Crew) error }~. ~Save(store CrewStore) error~ تدي المخزن لقطة مستقلة مرتبة لكل الطاقم، حتى لو فاضية. nil interface ترجع error؛ typed nil pointer جوه interface مش مطلوب دعمها. رجّع خطأ المخزن مع الحفاظ على السبب، وسيب الدليل زي ما هو حتى لو المخزن عدّل البيانات أو فشل. العقد تسلسلي؛ مفيش CRUD تاني أو ملف أو Report مطلوب.
`, Hints: []string{"افحص الفراغ والتكرار قبل append، ومتغيّرش الحقول الأصلية.", "Lowercase للطرفين ثم Contains للاسم أو الدور؛ رتّب نسخة النتائج بالـ ID.", "Save تدي CrewStore نسخة مرتبة؛ رجّع خطأ التخزين، واختبر مخزن بيعدّل slice اللي استلمها."},
		En: &LevelText{Title: "Crew Directory", Goal: "Register crew, search them, and save through a small interface", Files: []LabFile{{Name: "main.go", Starter: crewDirectoryStarter}},
			Lesson: `# A Crew Directory Independent of Disk

Communications needs to find crew by callsign or specialty and save the directory. Search need not know whether storage uses a file or memory. Put saving behind a small consumer-owned interface as in 9-4.

## Before you write
- How will you reject a duplicate ID without replacing the existing crew member?
- If storage modifies the slice it receives, how will you protect the directory's state?

## Case-insensitive substring search
~~~go
name:=strings.ToLower("ÉCHO Pilot")
query:=strings.ToLower("PIL")
fmt.Println(strings.Contains(name,query)) // true
~~~
~strings.ToLower~ comes from 18-4. The new ~strings.Contains~ checks for a substring, rather than equality or a whole word. An empty query is contained in every string, so it matches everyone. Lowercase both operands; Unicode normalization is not required.

## A consumer-owned save interface
~~~go
type BadgeSink interface { PutBadges([]string) error }
type MemoryBadges struct { Names []string }
func (m *MemoryBadges) PutBadges(names []string) error {
 m.Names=append([]string(nil),names...)
 return nil
}
var sink BadgeSink=&MemoryBadges{}
err:=sink.PutBadges([]string{"pilot"})
fmt.Println(err) // <nil>
~~~
Any type with the exact method signature satisfies the interface without an implements declaration. This badge example uses a pointer receiver to change storage. The ~names...~ expansion appends slice elements as in 8-4. Define ~CrewStore~ similarly with ~SaveCrew([]Crew) error~, without imposing JSON or a file path on its consumer.

## Independent snapshots and stable order
~~~go
names:=[]string{"pilot"}
copyOfNames:=append([]string(nil),names...)
copyOfNames[0]="changed"
fmt.Println(names[0]) // Original name
~~~
Copy before handing data to storage. Crew fields are strings, so copying slice elements is sufficient; there are no nested pointers. Use ~slices.SortFunc~ from 6-4 and ~strings.Compare~ from 18-4 for ascending IDs, rather than map or registration order.

Propagate ~SaveCrew~ failures without changing the directory. Storage owns its external writes; the consumer does not promise to reverse external effects after a storage failure.

## Using the directory
~~~go
var directory CrewDirectory
err:=directory.Register(Crew{ID:"C7",Callsign:"Echo",Role:"Pilot"})
fmt.Println(directory.Search("PIL"),err) // [{C7 Echo Pilot}] <nil>
err=directory.Save(store)
~~~
~store~ is any CrewStore implementation, like MemoryBadges above but accepting Crew values. Search exposes a copy; Save delegates to injected storage without changing the directory.

## Common mistakes
~~~go
if strings.Contains(crew.Callsign,strings.ToLower(query)) {} // First operand still has mixed case
return store.SaveCrew(d.crew) // Storage can corrupt the directory's slice
~~~
Check duplicates before appending. Saving is not another registration and must not modify in-memory crew records.

## Your mission
Implement a usable zero-value ~CrewDirectory~. ~Crew~ contains ~ID,Callsign,Role string~. ~Register(crew Crew) error~ rejects fields blank after TrimSpace and exact duplicate IDs, without mutation. Preserve original field strings without trimming. IDs are case and whitespace sensitive and are supplied rather than generated. Directories have independent state.

~Search(query string) []Crew~ returns crew whose callsign or role contains the query after ToLower on both sides. Do not search IDs or TrimSpace the query. Empty queries return everyone; misses have length zero. Return independent snapshots sorted by ascending ID using Go string order.

Define ~CrewStore interface { SaveCrew([]Crew) error }~. ~Save(store CrewStore) error~ gives storage an independent sorted snapshot of all crew, including an empty directory. A nil interface returns an error; typed nil pointers inside interfaces are outside the contract. Return storage failures preserving their causes. Leave the directory unchanged even if storage modifies its input or fails. This sequential contract requires no other CRUD operations, file implementation, or Report.
`, Hints: []string{"Validate blank fields and duplicates before append without changing original strings.", "Lowercase both sides and use Contains on callsign or role; sort a fresh result by ID.", "Give CrewStore a sorted copy; propagate storage errors and test a store that mutates its input."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: crewDirectoryHidden},
	})
}

const missionLogImplementation = `package main

import (
	"errors"
	"strings"
)

type MissionEntry struct {
	ID     int
	Sector string
	Sealed bool
}
type MissionLog struct {
	entries []MissionEntry
	issued  int
}

func (m *MissionLog) Log(sector string) (int, error) {
	sector = strings.TrimSpace(sector)
	if sector == "" {
		return 0, errors.New("sector required")
	}
	m.issued++
	m.entries = append(m.entries, MissionEntry{ID: m.issued, Sector: sector})
	return m.issued, nil
}
func (m *MissionLog) Seal(id int) error {
	for i := range m.entries {
		if m.entries[i].ID == id {
			m.entries[i].Sealed = true
			return nil
		}
	}
	return errors.New("mission ID not found")
}
func (m *MissionLog) Erase(id int) error {
	for i := range m.entries {
		if m.entries[i].ID == id {
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return nil
		}
	}
	return errors.New("mission ID not found")
}
func (m *MissionLog) Entries() []MissionEntry {
	return append([]MissionEntry(nil), m.entries...)
}
func main() {}
`
const missionLogTestStarter = `package main

import "testing"

func TestMissionLog(t *testing.T) {
 t.Log("add table-driven scenarios and assertions")
}
`
const missionLogHidden = `package main

import (
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"reflect"
	"testing"
)

func TestHiddenMissionHistory(t *testing.T) {
	for _, seed := range []int64{181, 904, 7321} {
		world := levels.Start("18-1-mission-log", io.Discard, seed)
		var log MissionLog
		var want []MissionEntry
		issued := 0
		for step := 0; step < 200; step++ {
			switch world.Rand.Intn(4) {
			case 0:
				sector := labkit.Text(world.Rand.Int63(), 2)
				id, err := log.Log(" " + sector + " ")
				issued++
				if err != nil || id != issued {
					t.Fatalf("stable ID: %d %v", id, err)
				}
				want = append(want, MissionEntry{ID: id, Sector: sector})
			case 1:
				if len(want) > 0 {
					i := world.Rand.Intn(len(want))
					if err := log.Seal(want[i].ID); err != nil {
						t.Fatal(err)
					}
					want[i].Sealed = true
				}
			case 2:
				if len(want) > 0 {
					i := world.Rand.Intn(len(want))
					if err := log.Erase(want[i].ID); err != nil {
						t.Fatal(err)
					}
					want = append(want[:i], want[i+1:]...)
				}
			case 3:
				if id, err := log.Log(" \t"); id != 0 || err == nil {
					t.Fatal("blank sector accepted")
				}
				for _, id := range []int{-1, 0, issued + 100} {
					if log.Seal(id) == nil || log.Erase(id) == nil {
						t.Fatal("unknown ID accepted")
					}
				}
			}
			got := log.Entries()
			if len(got) != len(want) {
				t.Fatal("entry count")
			}
			if len(want) > 0 && !reflect.DeepEqual(got, want) {
				t.Fatalf("history: %v want %v", got, want)
			}
			if len(got) > 0 {
				got[0].Sector = "tampered"
				if !reflect.DeepEqual(log.Entries(), want) {
					t.Fatal("snapshot alias")
				}
			}
		}
	}
}
`
const blackBoxStarter = `package main

import "errors"

type BlackBox struct {
	Cycle     int
	Callsigns []string
}

func WriteBlackBox(path string, box BlackBox) error { return errors.New("black box unfinished") }
func ReadBlackBox(path string) (BlackBox, error) {
	return BlackBox{}, errors.New("black box unfinished")
}
func main() {}
`
const blackBoxHidden = `package main

import (
	"encoding/json"
	"errors"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestHiddenBlackBoxRoundTrip(t *testing.T) {
	world := levels.Start("18-2-black-box", io.Discard, 182)
	path := filepath.Join(t.TempDir(), "capsule.json")
	for i := 0; i < 30; i++ {
		names := strings.Split(labkit.Text(world.Rand.Int63(), 1+world.Rand.Intn(20)), " ")
		names = append(names, "dock, \"north\"\n\\relay", "<ice>&")
		if i == 0 {
			names = []string{}
		}
		want := BlackBox{Cycle: 1 + world.Rand.Intn(10000), Callsigns: names}
		if err := WriteBlackBox(path, want); err != nil {
			t.Fatal(err)
		}
		got, err := ReadBlackBox(path)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("round trip: %v %v", got, err)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var schema map[string]json.RawMessage
		if err := json.Unmarshal(raw, &schema); err != nil {
			t.Fatal(err)
		}
		if len(schema) != 2 || schema["cycle"] == nil || schema["callsigns"] == nil {
			t.Fatal("wire names must use tags", string(raw))
		}
	}
}
func TestHiddenBlackBoxFailureKeepsEvidence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	for _, raw := range []string{"", "{", "null", "[]", "{}", ` + "`" + `{"cycle":1}` + "`" + `, ` + "`" + `{"callsigns":[]}` + "`" + `, ` + "`" + `{"cycle":0,"callsigns":[]}` + "`" + `, ` + "`" + `{"cycle":1,"callsigns":null}` + "`" + `, ` + "`" + `{"cycle":1,"callsigns":[" "]}` + "`" + `, ` + "`" + `{"cycle":"1","callsigns":[]}` + "`" + `, ` + "`" + `{"cycle":1,"callsigns":[]} {}` + "`" + `} {
		if err := os.WriteFile(path, []byte(raw), 0600); err != nil {
			t.Fatal(err)
		}
		got, err := ReadBlackBox(path)
		if err == nil || !strings.Contains(err.Error(), path) || !reflect.DeepEqual(got, BlackBox{}) {
			t.Fatalf("invalid %q: %v %v", raw, got, err)
		}
		after, readErr := os.ReadFile(path)
		if readErr != nil || string(after) != raw {
			t.Fatal("load overwrote evidence")
		}
	}
	missing := filepath.Join(dir, "missing.json")
	if _, err := ReadBlackBox(missing); !errors.Is(err, os.ErrNotExist) || !strings.Contains(err.Error(), missing) {
		t.Fatal("missing file context", err)
	}
	if err := WriteBlackBox(filepath.Join(dir, "absent", "box.json"), BlackBox{1, []string{}}); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("write failure", err)
	}
	for _, box := range []BlackBox{{}, {0, []string{}}, {-1, []string{}}, {1, nil}, {1, []string{"\t"}}} {
		if err := os.WriteFile(path, []byte("evidence"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := WriteBlackBox(path, box); err == nil {
			t.Fatal("invalid save accepted")
		}
		after, err := os.ReadFile(path)
		if err != nil || string(after) != "evidence" {
			t.Fatal("invalid save changed destination")
		}
	}
	if err := os.WriteFile(path, []byte(` + "`" + `{"cycle":7,"callsigns":[],"future":true}` + "`" + `), 0600); err != nil {
		t.Fatal(err)
	}
	if got, err := ReadBlackBox(path); err != nil || got.Cycle != 7 || got.Callsigns == nil {
		t.Fatal("unknown fields should be ignored", err)
	}
}
`
const cargoManifestStarter = `package main

import (
	"errors"
	"io"
)

type CargoTotals struct {
	Rows, Units int
	ByCallsign  map[string]int
}

func ReadManifest(in io.Reader) (CargoTotals, error) {
	return CargoTotals{}, errors.New("manifest unfinished")
}
func main() {}
`
const cargoManifestHidden = `package main

import (
	"encoding/csv"
	"errors"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenManifestSeededTotals(t *testing.T) {
	world := levels.Start("18-3-cargo-manifest", io.Discard, 183)
	for _, n := range []int{0, 1, 127, 100000} {
		text := labkit.CSV(world.Rand.Int63(), n)
		rows, err := csv.NewReader(strings.NewReader(text)).ReadAll()
		if err != nil {
			t.Fatal(err)
		}
		want := CargoTotals{ByCallsign: map[string]int{}}
		for _, row := range rows[1:] {
			units, err := strconv.Atoi(row[2])
			if err != nil {
				t.Fatal(err)
			}
			want.Rows++
			want.Units += units
			want.ByCallsign[row[1]] += units
		}
		var source io.Reader = strings.NewReader(text)
		if n < 1000 {
			source = iotest.OneByteReader(source)
		}
		got, err := ReadManifest(source)
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("n=%d totals %v want %v: %v", n, got, want, err)
		}
	}
}
func TestHiddenManifestRejectsBrokenRecords(t *testing.T) {
	for _, tc := range []struct{ text, context string }{
		{"", "header"}, {"id,callsign\n", "header"}, {"ID,callsign,cargo\n", "header"}, {"id,id,cargo\n", "header"},
		{"id,callsign,cargo\n0,ok,1\n1,short\n", "record 3"},
		{"id,callsign,cargo\n0,ok,1\n1,x,2,extra\n", "record 3"},
		{"id,callsign,cargo\n0,ok,1\n1,\"bad,2\n", "record 3"},
	} {
		got, err := ReadManifest(strings.NewReader(tc.text))
		if err == nil || !strings.Contains(err.Error(), tc.context) || !reflect.DeepEqual(got, CargoTotals{}) {
			t.Fatalf("%q: %v %v", tc.text, got, err)
		}
	}
	for _, fields := range []string{"-1,x,1", "x,x,1", "0, ,1", "0,x,abc", "0,x,NaN", "0,x,Inf", "0,x,-1", "0,x,1000001", "0,x,", "0,x,1e3"} {
		got, err := ReadManifest(strings.NewReader("id,callsign,cargo\n" + fields + "\n"))
		if err == nil || !strings.Contains(err.Error(), "record 2") || !reflect.DeepEqual(got, CargoTotals{}) {
			t.Fatal("invalid fields", fields, got, err)
		}
	}
	wire := errors.New("manifest wire failed")
	if got, err := ReadManifest(iotest.ErrReader(wire)); !errors.Is(err, wire) || !reflect.DeepEqual(got, CargoTotals{}) {
		t.Fatal("header reader failure", got, err)
	}
	source := io.MultiReader(strings.NewReader("id,callsign,cargo\n0,x,1\n"), iotest.ErrReader(wire))
	if got, err := ReadManifest(source); !errors.Is(err, wire) || !reflect.DeepEqual(got, CargoTotals{}) {
		t.Fatal("late reader failure", got, err)
	}
	got, err := ReadManifest(strings.NewReader("id,callsign,cargo\r\n0, x ,0\r\n1,x,1000000\r\n\r\n"))
	if err != nil || got.Rows != 2 || got.Units != 1000000 || got.ByCallsign["x"] != 1000000 {
		t.Fatal("boundaries/trim/CRLF", got, err)
	}
}

func TestHiddenManifestQuotedNames(t *testing.T) {
	input := "id,callsign,cargo\n0,\"he said \"\"dock\"\"\",2\n1,relay,3"
	got, err := ReadManifest(strings.NewReader(input))
	if err != nil || got.Rows != 2 || got.Units != 5 || got.ByCallsign["he said \"dock\""] != 2 || got.ByCallsign["relay"] != 3 {
		t.Fatal("escaped quote and final record", got, err)
	}
}

func TestHiddenBonusMultilineCargo(t *testing.T) {
	text := "id,callsign,cargo\n0,\"south, \"\"ice\"\"\nrelay\",7\n1,\"south, \"\"ice\"\"\nrelay\",8"
	got, err := ReadManifest(iotest.OneByteReader(strings.NewReader(text)))
	if err != nil || got.Rows != 2 || got.Units != 15 || len(got.ByCallsign) != 1 || got.ByCallsign["south, \"ice\"\nrelay"] != 15 {
		t.Fatal("multiline quoted record", got, err)
	}
}

`
const signalWordsStarter = `package main

import (
	"errors"
	"io"
)

type WordCount struct {
	Word  string
	Count int
}

func TopSignals(in io.Reader, limit int) ([]WordCount, error) {
	return nil, errors.New("signals unfinished")
}
func main() {}
`
const signalWordsHidden = `package main

import (
	"errors"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenSignalsTokenRules(t *testing.T) {
	for _, tc := range []struct {
		text  string
		limit int
		want  []WordCount
	}{
		{"", 3, nil}, {" \t\n!🚀--'", 9, nil},
		{"Orbit, ORBIT! ice.ice; well-known don't abc123 ٤ É é", 20, []WordCount{{"ice", 2}, {"orbit", 2}, {"é", 2}, {"abc123", 1}, {"don", 1}, {"known", 1}, {"t", 1}, {"well", 1}, {"٤", 1}}},
		{"z b a z B A", 2, []WordCount{{"a", 2}, {"b", 2}}},
		{"lonely", 5, []WordCount{{"lonely", 1}}},
	} {
		got, err := TopSignals(iotest.OneByteReader(strings.NewReader(tc.text)), tc.limit)
		if err != nil || !slices.Equal(got, tc.want) {
			t.Fatalf("%q: %v %v want %v", tc.text, got, err, tc.want)
		}
	}
}
func TestHiddenSignalsSeededStream(t *testing.T) {
	world := levels.Start("18-4-signal-words", io.Discard, 184)
	text := labkit.Text(world.Rand.Int63(), 60000)
	counts := map[string]int{}
	for _, word := range strings.Fields(text) {
		if word != "🚀" && word != "🛰️" {
			counts[word]++
		}
	}
	var want []WordCount
	for word, count := range counts {
		want = append(want, WordCount{word, count})
	}
	slices.SortFunc(want, func(a, b WordCount) int {
		if a.Count != b.Count {
			return b.Count - a.Count
		}
		return strings.Compare(a.Word, b.Word)
	})
	for _, limit := range []int{1, 2, 99} {
		expected := want
		if len(expected) > limit {
			expected = expected[:limit]
		}
		for repeat := 0; repeat < 2; repeat++ {
			got, err := TopSignals(iotest.HalfReader(strings.NewReader(text)), limit)
			if err != nil || !slices.Equal(got, expected) {
				t.Fatal("seeded stream/order", got, expected, err)
			}
		}
	}
}
func TestHiddenSignalsLimitsAndFailure(t *testing.T) {
	wire := errors.New("signal wire failed")
	if got, err := TopSignals(iotest.ErrReader(wire), 0); err != nil || len(got) != 0 {
		t.Fatal("zero limit must not read", got, err)
	}
	if got, err := TopSignals(strings.NewReader("a"), -1); err == nil || got != nil {
		t.Fatal("negative limit", got, err)
	}
	source := io.MultiReader(strings.NewReader("ice\n"), iotest.ErrReader(wire))
	if got, err := TopSignals(source, 2); !errors.Is(err, wire) || got != nil {
		t.Fatal("reader error must discard report", got, err)
	}
	for _, ending := range []string{"", "\n", "\r\n"} {
		word := strings.Repeat("x", 1<<20)
		got, err := TopSignals(strings.NewReader(word+ending), 1)
		if err != nil || len(got) != 1 || got[0].Word != word || got[0].Count != 1 {
			t.Fatal("exact line limit", err)
		}
	}
	for _, size := range []int{(1 << 20) + 1, 2 << 20} {
		if got, err := TopSignals(strings.NewReader("ok\n"+strings.Repeat("x", size)), 1); err == nil || got != nil {
			t.Fatal("oversized line must fail without partial report", err)
		}
	}
}
`
const flightRecorderStarter = `package main

import (
	"errors"
	"io"
)

type FlightSummary struct {
	Info, Warn, Error int
	BadLines          []int
}

func ReadFlight(in io.Reader) (FlightSummary, error) {
	return FlightSummary{}, errors.New("recorder unfinished")
}
func main() {}
`
const flightRecorderHidden = `package main

import (
	"errors"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"reflect"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenFlightCountsAndDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		text string
		want FlightSummary
	}{
		{"", FlightSummary{}},
		{"0000 INFO launch\r\n1 WARN low\n2 ERROR down", FlightSummary{Info: 1, Warn: 1, Error: 1}},
		{"\n \n0 info x\n-1 INFO x\nx WARN y\n2 INFO\n3 ERROR   \n4  INFO x\n 5 INFO x\n6 TRACE x\n7 INFO ok\n8 WARN hi there  ", FlightSummary{Info: 1, Warn: 1, BadLines: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}},
	} {
		got, err := ReadFlight(iotest.OneByteReader(strings.NewReader(tc.text)))
		if err != nil || !reflect.DeepEqual(got, tc.want) {
			t.Fatalf("%q: %v want %v: %v", tc.text, got, tc.want, err)
		}
	}
}
func TestHiddenFlightSeededStream(t *testing.T) {
	world := levels.Start("18-5-flight-recorder", io.Discard, 185)
	var stream strings.Builder
	var want FlightSummary
	line := 0
	for batch := 0; batch < 200; batch++ {
		logs := labkit.Logs(world.Rand.Int63(), 20)
		for _, row := range strings.Split(strings.TrimSuffix(logs, "\n"), "\n") {
			line++
			switch strings.Fields(row)[1] {
			case "INFO":
				want.Info++
			case "WARN":
				want.Warn++
			case "ERROR":
				want.Error++
			}
			stream.WriteString(row + "\n")
		}
		line++
		want.BadLines = append(want.BadLines, line)
		stream.WriteString("signal lost\n")
	}
	got, err := ReadFlight(iotest.HalfReader(strings.NewReader(stream.String())))
	if err != nil || got.Info != want.Info || got.Warn != want.Warn || got.Error != want.Error || !slices.Equal(got.BadLines, want.BadLines) {
		t.Fatal("stream totals/diagnostic order", got, err)
	}
}
func TestHiddenFlightStructuralFailures(t *testing.T) {
	wire := errors.New("flight wire failed")
	source := io.MultiReader(strings.NewReader("0 INFO good\n"), iotest.ErrReader(wire))
	if got, err := ReadFlight(source); !errors.Is(err, wire) || !reflect.DeepEqual(got, FlightSummary{}) {
		t.Fatal("read failure must discard summary", got, err)
	}
	if got, err := ReadFlight(iotest.DataErrReader(io.MultiReader(strings.NewReader("0 INFO good\n"), iotest.ErrReader(wire)))); !errors.Is(err, wire) || !reflect.DeepEqual(got, FlightSummary{}) {
		t.Fatal("data and error", got, err)
	}
	for _, n := range []int{(256 << 10) + 1, 512 << 10} {
		if got, err := ReadFlight(strings.NewReader("0 INFO fine\n" + strings.Repeat("x", n))); err == nil || !reflect.DeepEqual(got, FlightSummary{}) {
			t.Fatal("over-limit line", got, err)
		}
	}
}
func TestHiddenBonusFlightExtendedBuffer(t *testing.T) {
	for _, n := range []int{64 << 10, 200 << 10, 256 << 10} {
		for _, ending := range []string{"", "\n", "\r\n"} {
			line := "12 WARN " + strings.Repeat("x", n-len("12 WARN "))
			got, err := ReadFlight(iotest.HalfReader(strings.NewReader(line + ending)))
			if err != nil || !reflect.DeepEqual(got, FlightSummary{Warn: 1}) {
				t.Fatalf("line bytes=%d ending=%q: %v %v", n, ending, got, err)
			}
		}
	}
}
`
const crewDirectoryStarter = `package main

import "errors"

type Crew struct{ ID, Callsign, Role string }
type CrewStore interface{ SaveCrew([]Crew) error }
type CrewDirectory struct{ crew []Crew }

func (d *CrewDirectory) Register(crew Crew) error   { return errors.New("directory unfinished") }
func (d *CrewDirectory) Search(query string) []Crew { return nil }
func (d *CrewDirectory) Save(store CrewStore) error { return errors.New("directory unfinished") }
func main()                                         {}
`
const crewDirectoryHidden = `package main

import (
	"errors"
	"fmt"
	"io"
	"orbital/colony/labkit"
	"orbital/colony/levels"
	"reflect"
	"slices"
	"strings"
	"testing"
)

type hiddenCrewStore struct {
	saved   []Crew
	failure error
}

func (s *hiddenCrewStore) SaveCrew(crew []Crew) error {
	s.saved = append([]Crew(nil), crew...)
	if len(crew) > 0 {
		crew[0].Callsign = "store mutation"
	}
	return s.failure
}
func TestHiddenCrewRegistrationAndSearch(t *testing.T) {
	world := levels.Start("18-6-crew-directory", io.Discard, 186)
	var directory CrewDirectory
	var want []Crew
	if len(directory.Search("")) != 0 {
		t.Fatal("zero value")
	}
	for i := 40; i >= 0; i-- {
		crew := Crew{ID: fmt.Sprintf("C%03d", i), Callsign: labkit.Text(world.Rand.Int63(), 2), Role: []string{"Pilot", "ENGINEER", "Medic"}[i%3]}
		if err := directory.Register(crew); err != nil {
			t.Fatal(err)
		}
		want = append(want, crew)
		if err := directory.Register(Crew{crew.ID, "intruder", "visitor"}); err == nil {
			t.Fatal("duplicate ID accepted")
		}
	}
	slices.SortFunc(want, func(a, b Crew) int { return strings.Compare(a.ID, b.ID) })
	for _, query := range []string{"", "PILOT", "engine", "orbit", "قمر", "no-match", "C000"} {
		var expected []Crew
		for _, crew := range want {
			if strings.Contains(strings.ToLower(crew.Callsign), strings.ToLower(query)) || strings.Contains(strings.ToLower(crew.Role), strings.ToLower(query)) {
				expected = append(expected, crew)
			}
		}
		got := directory.Search(query)
		if !slices.Equal(got, expected) {
			t.Fatal("search/order", query, got, expected)
		}
		if len(got) > 0 {
			got[0].Role = "altered"
			if !slices.Equal(directory.Search(""), want) {
				t.Fatal("search alias")
			}
		}
	}
	for _, crew := range []Crew{{}, {" ", "a", "b"}, {"new", "\t", "b"}, {"new", "a", " "}} {
		if err := directory.Register(crew); err == nil || !slices.Equal(directory.Search(""), want) {
			t.Fatal("invalid registration mutated state")
		}
	}
	var other CrewDirectory
	if len(other.Search("")) != 0 {
		t.Fatal("directories share state")
	}
	if err := other.Register(Crew{"a", "Écho", "Science"}); err != nil {
		t.Fatal(err)
	}
	if got := other.Search("ÉCHO"); len(got) != 1 || got[0].ID != "a" {
		t.Fatal("Unicode lowercase", got)
	}
	if err := other.Register(Crew{"A", " echo ", " pilot "}); err != nil {
		t.Fatal("IDs are case sensitive", err)
	}
	if len(other.Search(" echo ")) != 1 {
		t.Fatal("fields/query must retain spaces")
	}
}
func TestHiddenCrewSaveIsolationAndFailure(t *testing.T) {
	var directory CrewDirectory
	store := &hiddenCrewStore{}
	if err := directory.Save(store); err != nil || len(store.saved) != 0 {
		t.Fatal("empty save", err)
	}
	for _, crew := range []Crew{{"b", "Moon", "Pilot"}, {"a", "Ice", "Medic"}} {
		if err := directory.Register(crew); err != nil {
			t.Fatal(err)
		}
	}
	want := directory.Search("")
	if err := directory.Save(store); err != nil || !reflect.DeepEqual(store.saved, want) || !reflect.DeepEqual(directory.Search(""), want) {
		t.Fatal("save sorted snapshot", err)
	}
	wire := errors.New("archive offline")
	store.failure = wire
	if err := directory.Save(store); !errors.Is(err, wire) || !reflect.DeepEqual(directory.Search(""), want) {
		t.Fatal("failed store must not change directory", err)
	}
	if err := directory.Save(nil); err == nil {
		t.Fatal("nil store")
	}
}
`
