package levels

func init() {
	add(&Level{ID: "21-1-sysinfo", Chapter: 21, Kind: "lab", Title: "حالة المحطة", Goal: "اقرأ flags واعرض لقطة نظام محقونة وبيئة مسموحة بس",
		Requires: []string{"sel:flag.NewFlagSet"}, Bridge: []int{9}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: stationInfoStarterAr}},
		Lesson: `# تقرير محطة من غير مفاجآت

الطاقم محتاج تقرير حالة يقدر يقارنه بين تشغيلتين. بيانات الجهاز الحقيقي بتتجمع عند حدود البرنامج، لكن دالة التقرير بتستقبل لقطة جاهزة عشان الاختبار ميتغيرش بين Windows وLinux.

## قبل ما تكتب
- إزاي تفرّق بين متغير بيئة مش موجود ومتغير موجود وقيمته فاضية؟
- إزاي ترفض آخر argument بايظ قبل ما تكتب أول سطر؟

## flags محلية لكل استدعاء
استورد ~flag~. ~flag.NewFlagSet~ بتنشيء parser مستقل؛ ~flag.ContinueOnError~ تخليه يرجّع خطأ بدل ما يقفل العملية. ده مهم عشان نختبر أكتر من تشغيل في نفس العملية.
~~~go
fs:=flag.NewFlagSet("antenna",flag.ContinueOnError)
fs.SetOutput(io.Discard)
verbose:=fs.Bool("verbose",false,"include details")
err:=fs.Parse([]string{"-verbose"})
fmt.Println(*verbose,fs.NArg(),err) // true 0 <nil>
~~~
~FlagSet.Bool~ بترجع pointer لقيمة الـ flag. القيمة الافتراضية بتفضل لو الـ flag مش موجودة. ~-verbose=false~ تطفيها صراحة؛ Boolean flags مبتاخدش القيمة كـ argument منفصل.

~FlagSet.Parse~ تستقبل args من غير اسم البرنامج. ~FlagSet.NArg~ عدد arguments الباقية بعد flags. ~FlagSet.SetOutput~ يحدد مكان diagnostics؛ ~io.Discard~ writer بيرميها. في العقد هنرجّع الخطأ للمنادي من غير diagnostics تلقائية في التقرير.

## البيئة: اسم واحد، وثلاث حالات
استورد ~os~ لو بتوصل بالبيئة الحقيقية. ~os.Getenv~ بترجع string بس؛ ~os.LookupEnv~ بترجع القيمة ومعاها هل الاسم موجود.
~~~go
text:=os.Getenv("COLONY_CREW")
crew,present:=os.LookupEnv("COLONY_CREW")
fmt.Println(text,crew,present)
lookup:=func(name string)(string,bool){return "pilot",name=="COLONY_CREW"}
crew,present=lookup("COLONY_CREW")
~~~
الـ function value من 5-5 بتسمح بحقن lookup بدل قراءة البيئة الحقيقية. متلفش على كل متغيرات البيئة؛ قائمة السماح هنا اسم واحد فقط. غياب الاسم غير وجوده بقيمة فاضية.

## لقطة runtime وحدود البرنامج
استورد ~runtime~. ~runtime.GOOS~ و ~runtime.GOARCH~ ثوابت لهدف بناء البرنامج؛ ~runtime.Version()~ نسخة Go اللي بنته. مفيش داعي تستعلم عنهم جوه منطق التقرير.
~~~go
host:=StationHost{OS:runtime.GOOS,Arch:runtime.GOARCH,Go:runtime.Version(),Lookup:os.LookupEnv}
err:=StationInfo(os.Args[1:],host,os.Stdout)
if err!=nil {fmt.Fprintln(os.Stderr,err)}
~~~
~os.Args~ slice أولها اسم البرنامج؛ ~os.Stdout~ و ~os.Stderr~ writers للخروج العادي والتشخيص. المثال adapter حقيقي، مش مطلوب في الاختبار؛ ~main~ ممكن تفضل فاضية. التست بيدي ~StationHost~ بقيم ثابتة ودالة lookup، ومش بيقرأ نظامك.

## تنسيق ثابت حتى لو القيمة فيها newline
~~~go
_,err:=fmt.Fprintf(out,"label=%q\n","north\nrelay")
~~~
~fmt.Fprintf~ من 9-4 بتكتب للـ writer وترجّع عدد بايتات وخطأ. ~%q~ بتحط string بين quotes وتهرّب newline وquotes. احتفظ بسبب خطأ الكتابة؛ متحوّلوش لنجاح. مفاتيح التقرير بروتوكول ثابت بالإنجليزي، مش رسائل مترجمة.

## أخطاء شائعة
~~~go
if os.Getenv("COLONY_CREW")=="" { /* غلط: الفاضي مش بالضرورة غايب */ }
fmt.Println(runtime.GOOS) // غلط جوه StationInfo: تجاهل اللقطة والـ writer
~~~

## المهمة
نفّذ ~StationInfo(args []string, host StationHost, out io.Writer) error~. ~StationHost~ فيها ~OS,Arch,Go string~ و ~Lookup func(string)(string,bool)~. الـ writer ودالة lookup غير nil. اعمل ~flag.NewFlagSet~ جديدة كل مرة: ~-system~ Boolean افتراضي true، و ~-crew~ Boolean افتراضي false. استخدم ~flag.ContinueOnError~ و ~io.Discard~ للتشخيص.

ارفض أي خطأ parsing، بما فيه طلب help، وأي arguments متبقية، قبل كتابة أي تقرير أو استدعاء lookup. مفيش تعديل لـ args. بعد النجاح، لو system مفعلة اكتب بالترتيب ~os=%q~ ثم ~arch=%q~ ثم ~go=%q~، كل قيمة من host وكل سطر آخره newline.

لو crew مفعلة، بعدها اقرأ ~COLONY_CREW~ مرة واحدة فقط: لو موجود اكتب ~crew=%q~ بقيمته الأصلية حتى لو فاضية؛ لو غايب اكتب ~crew=<missing>~. لو crew مطفية متستدعيش lookup. إطفاء الاتنين نجاح من غير إخراج. فشل الكتابة يوقف فوراً ويرجع خطأ يحتفظ بالسبب؛ الأجزاء اللي اتكتبت ممكن تفضل. مفيش قراءة حقيقية لـ os أو runtime جوه StationInfo، ولا Report.
`, Hints: []string{"اعمل FlagSet جديدة وحدد defaults قبل Parse؛ افحص NArg قبل أي أثر.", "dereference لنتائج Bool؛ حضّر سطور system بترتيب ثابت واستخدم %q.", "افحص bool من Lookup مش القيمة الفاضية؛ رجّع أول خطأ كتابة بنفس السبب."},
		En: &LevelText{Title: "Station Status", Goal: "Parse flags and report injected system data and an allowlisted environment value", Files: []LabFile{{Name: "main.go", Starter: stationInfoStarter}},
			Lesson: `# A Predictable Station Report

The crew needs reports it can compare between runs. Gather real host data at the program boundary; pass a snapshot into the reporting function so tests stay the same on Windows and Linux.

## Before you write
- How will you distinguish a missing environment variable from a present empty one?
- How will you reject a bad final argument before writing the first line?

## Flags local to each invocation
Import ~flag~. ~flag.NewFlagSet~ creates an independent parser. ~flag.ContinueOnError~ returns errors rather than terminating the process, allowing several invocations in one test process.
~~~go
fs:=flag.NewFlagSet("antenna",flag.ContinueOnError)
fs.SetOutput(io.Discard)
verbose:=fs.Bool("verbose",false,"include details")
err:=fs.Parse([]string{"-verbose"})
fmt.Println(*verbose,fs.NArg(),err) // true 0 <nil>
~~~
~FlagSet.Bool~ returns a pointer to the flag's value. The default remains when the flag is absent. ~-verbose=false~ explicitly disables it; Boolean flags do not take a separate value argument.

~FlagSet.Parse~ receives arguments without the executable name. ~FlagSet.NArg~ counts remaining positional arguments. ~FlagSet.SetOutput~ chooses the diagnostics writer; ~io.Discard~ throws writes away. This contract returns errors without automatic diagnostics in the report.

## Environment: one name, three states
Import ~os~ when connecting to the real environment. ~os.Getenv~ returns only a string; ~os.LookupEnv~ returns the value and a presence flag.
~~~go
text:=os.Getenv("COLONY_CREW")
crew,present:=os.LookupEnv("COLONY_CREW")
fmt.Println(text,crew,present)
lookup:=func(name string)(string,bool){return "pilot",name=="COLONY_CREW"}
crew,present=lookup("COLONY_CREW")
~~~
Function values from 5-5 allow an injected lookup instead of real environment access. Never enumerate the environment; the allowlist here contains exactly one name. Missing differs from present with an empty value.

## Runtime snapshot and the program boundary
Import ~runtime~. ~runtime.GOOS~ and ~runtime.GOARCH~ are build-target constants; ~runtime.Version()~ reports the Go version that built the program. The report logic need not query them itself.
~~~go
host:=StationHost{OS:runtime.GOOS,Arch:runtime.GOARCH,Go:runtime.Version(),Lookup:os.LookupEnv}
err:=StationInfo(os.Args[1:],host,os.Stdout)
if err!=nil {fmt.Fprintln(os.Stderr,err)}
~~~
~os.Args~ is a slice starting with the executable name. ~os.Stdout~ and ~os.Stderr~ are ordinary and diagnostic output writers. This is a real adapter example; tests do not require it and ~main~ may remain empty. Tests supply a ~StationHost~ with fixed strings and a lookup function, without reading your host.

## Stable formatting even with embedded newlines
~~~go
_,err:=fmt.Fprintf(out,"label=%q\n","north\nrelay")
~~~
~fmt.Fprintf~ from 9-4 writes to a writer and returns a byte count and error. ~%q~ quotes strings and escapes newlines and quotes. Preserve write error causes rather than returning success. Report keys are a fixed English protocol, rather than translated messages.

## Common mistakes
~~~go
if os.Getenv("COLONY_CREW")=="" { /* Wrong: empty does not imply missing */ }
fmt.Println(runtime.GOOS) // Wrong inside StationInfo: ignores snapshot and writer
~~~

## Your mission
Implement ~StationInfo(args []string, host StationHost, out io.Writer) error~. ~StationHost~ has ~OS,Arch,Go string~ and ~Lookup func(string)(string,bool)~. The writer and lookup are non-nil. Create a fresh ~flag.NewFlagSet~ per call: Boolean ~-system~ defaults to true and Boolean ~-crew~ defaults to false. Use ~flag.ContinueOnError~ with ~io.Discard~ for diagnostics.

Reject any parsing error, including help, and remaining positional arguments before output or lookup calls. Do not mutate args. After success, if system is enabled, write ~os=%q~, ~arch=%q~, then ~go=%q~, taking strings from host and ending every line with a newline.

If crew is enabled, look up only ~COLONY_CREW~, exactly once, after the system lines. When present, write ~crew=%q~ with the original value, including empty; when absent write ~crew=<missing>~. Disabled crew must not call lookup. Disabling both succeeds without output. Stop immediately on write failure and preserve its cause; previously written bytes may remain. StationInfo must not query real os or runtime values. No Report.
`, Hints: []string{"Create a fresh FlagSet with defaults; check NArg before side effects.", "Dereference Bool results and write system lines in fixed order using %q.", "Check Lookup's presence bool, not empty text; return the first write error with its cause."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: stationInfoHidden},
	})
	add(&Level{ID: "21-2-subcommands", Chapter: 21, Kind: "lab", Title: "أوامر فرعية", Goal: "شغّل أوامر ورشة الإصلاح بمدخلات محقونة وأكواد خروج واضحة",
		Requires: []string{"iface"}, Bridge: []int{16}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: repairCommandsStarter}},
		Lesson: `# أوامر ورشة الإصلاح

الورشة بتستقبل طلبات إصلاح: قطاع وعدد خلايا. هنفصل سجل الطلبات عن التيرمنال بـ interface، ونشغّل أمر واحد من args أو دفعة أوامر من stdin. الحالة في الذاكرة لجلسة واحدة.

## قبل ما تكتب
- إزاي حذف طلب ميخلّيش الإضافة الجديدة تعيد رقمه؟
- إمتى تكمل بعد خطأ أمر، وإمتى توقف الجلسة بسبب فشل النقل؟

## دالة run وكود الخروج
~run(args []string, stdin io.Reader, stdout io.Writer) int~ ترجع نتيجة للمنادي. الصفر نجاح، وغير الصفر فشل؛ الأرقام معناها عقد البرنامج نفسه، مش قاعدة عامة لكل أدوات النظام.
~~~go
func preview(args []string,in io.Reader,out io.Writer) int {
 if len(args)!=0 {return 2}
 if _,err:=fmt.Fprintln(out,"ready");err!=nil {return 1}
 return 0
}
func main() {os.Exit(preview(os.Args[1:],os.Stdin,os.Stdout))}
~~~
~os.Stdin~ reader للمدخل القياسي. ~os.Exit(code)~ تنهي العملية فوراً بكود الخروج ومبتشغّلش defer؛ خليها في ~main~ بس. الدالة القابلة للاختبار ترجع int، والتست ميناديش ~main~. في حلك ~main~ ممكن تفضل فاضية.

## interface عند المستهلك
هنرجع لـ interfaces من 9-4، بعقد ~RepairBook~ اللي تفصل تنفيذ الأمر عن الذاكرة. عرّفها في ~main.go~ واستخدمها كنوع parameter لدالة توزيع الأوامر.
~~~go
type LampBank interface { SetLamp(int) error }
func light(bank LampBank) error {return bank.SetLamp(3)}
~~~
أي struct عندها methods المطلوبة تحقق الـ interface تلقائياً. في المهمة ~RepairOrder~ فيها ~ID,Units int~ و ~Sector string~ و ~Done bool~. التنفيذ الملموس يحتفظ بـ slice وعدّاد إصدار منفصل، زي سجل 18-1، لكن المجال هنا طلبات إصلاح بعدد خلايا.
~~~go
id,err:=book.Queue("ice",4)
orders:=book.Orders()
err=book.Complete(id)
err=book.Remove(id)
fmt.Println(orders,err)
~~~
~book~ هنا ~RepairBook~. ~Queue~ تضيف طلب؛ ~Orders~ لقطة مستقلة بالترتيب؛ ~Complete~ تختم الإصلاح؛ ~Remove~ تحذفه. فشل أي عملية ميغيرش الحالة ولا يستهلك رقم. ختم طلب مختوم نجاح.

## أمر واحد أو سيل سطور
لو args موجودة، نفّذ الأمر ده فقط من غير قراءة stdin. لو فاضية، اقرأ سطر بسطر بـ ~bufio.Scanner~ من 9-4. ~strings.Fields~ تقسم السطر على Unicode whitespace وترجّع ~[]string~، فاسم القطاع هنا token واحدة.
~~~go
words:=strings.Fields("add ice 4")
fmt.Println(words) // [add ice 4]
units,err:=strconv.Atoi(words[2])
~~~
~fmt.Fprint~ تكتب string زي ما هي من غير إضافة newline؛ مفيدة لو جهزت رد الأمر كامل الأول.
~~~go
_,err:=fmt.Fprint(out,"queued=7\n")
~~~
افحص عدد الكلمات قبل indexing. ~strconv.Atoi~ من 4-3؛ بعد parsing افحص المدى. الأمر المرفوض يكتب سطر بروتوكول ~error~ ويكمل في وضع السيل، ويخلي كود النهاية 2 حتى لو بعده أوامر ناجحة.

~Scanner.Err~ غير nil أو فشل writer يوقف فوراً بكود 1. ~EOF~ العادية نجاح لو مفيش أوامر مرفوضة. كل سطر أمر أقصاه 4096 بايت؛ استخدم ~Scanner.Buffer~ من 18-4 مع مساحة CRLF وافحص طول المحتوى. الأكبر فشل نقل بكود 1 من غير تنفيذ السطر.

## أخطاء شائعة
~~~go
id:=len(orders)+1 // غلط بعد الحذف: رقم قديم هيتكرر
os.Exit(2) // غلط جوه run: هيقفل سيرفر الاختبار بدل الرجوع
~~~

## المهمة
نفّذ ~run(args []string, stdin io.Reader, stdout io.Writer) int~. كل استدعاء يبدأ بسجل جديد. args غير فاضية = أمر واحد من غير قراءة stdin؛ الفاضية = أوامر سطور لحد EOF، من غير prompts أو quit. السطر الفاضي أمر بايظ. متعدلش args.

عرّف واستخدم ~RepairBook interface{Queue(string,int)(int,error); Orders() []RepairOrder; Complete(int) error; Remove(int) error}~. وفّر ~RepairOrder~ بالحقول المذكورة وتنفيذ ذاكرة مستقل؛ مش مطلوب اسم معين للتنفيذ. ~Queue~ تقبل sector مش فاضية ومفيهاش Unicode whitespace، وunits من 1 لـ 100 شامل. IDs تبدأ 1 وتزيد مع كل إضافة ناجحة ولا تتعاد بعد الحذف. ~Orders~ لقطة مستقلة بترتيب الإضافة. ~Complete~ و ~Remove~ يرفضوا ID غير موجود، بما فيه الصفر والسالب، من غير تغيير.

الأوامر وعدد arguments بالظبط: ~add sector units~ يكتب ~queued=N~؛ ~done id~ يكتب ~done=N~؛ ~rm id~ يكتب ~removed=N~؛ ~list~ يكتب سطر لكل طلب بالترتيب بصيغة ~%d|%q|%d|%t~ للـ ID وSector وUnits وDone. الفاضي مفيهوش سطور. كل سطر بينتهي newline. الـ ID والـ units يقبلهم Atoi مع فحص المدى. أي أمر مجهول أو malformed أو خطأ domain يكتب ~error~ بس، من غير تغيير حالة.

الكود 0 لو كل الأوامر نجحت، 2 لو أي أمر اترفض، و1 لفشل القراءة أو الكتابة حتى بعد خطأ أمر. في وضع السيل كمّل بعد code 2؛ قف فوراً بعد code 1. إخراج سابق أو تعديل ناجح قبل فشل writer ممكن يفضل. الـ readers والـ writers غير nil، والأوامر الصالحة لحد 4096 بايت محتوى مع LF أو CRLF أو آخر سطر من غير newline. الأكبر code 1. مفيش ملفات أو حالة بين استدعاءين ولا Report.
`, Hints: []string{"ابدأ بذاكرة خلف RepairBook: slice وعدّاد منفصل؛ تحقق قبل التعديل.", "افصل تنفيذ أمر واحد عن Scanner؛ استخدم نفس الموزّع لو args موجودة.", "احتفظ بكود 2 بعد أمر بايظ وكمّل؛ خطأ Scanner أو writer يوقف بكود 1 فوراً."},
		En: &LevelText{Title: "Subcommands", Goal: "Run repair-shop commands with injected streams and explicit exit codes", Files: []LabFile{{Name: "main.go", Starter: repairCommandsStarter}},
			Lesson: `# Repair-Shop Commands

The shop receives repair orders: a sector and a cell count. Separate the order book from the terminal through an interface, then execute one command from args or a stream from stdin. State lives in memory for one invocation.

## Before you write
- How will deletion avoid recycling an old ID on the next insertion?
- When should a bad command allow continuation, and when should transport failure stop the session?

## run and exit codes
~run(args []string, stdin io.Reader, stdout io.Writer) int~ returns an outcome to its caller. Zero means success and nonzero means failure; individual numbers belong to this program's contract, rather than all system tools.
~~~go
func preview(args []string,in io.Reader,out io.Writer) int {
 if len(args)!=0 {return 2}
 if _,err:=fmt.Fprintln(out,"ready");err!=nil {return 1}
 return 0
}
func main() {os.Exit(preview(os.Args[1:],os.Stdin,os.Stdout))}
~~~
~os.Stdin~ reads standard input. ~os.Exit(code)~ immediately terminates the process without running defers; keep it at the ~main~ boundary. Testable logic returns an int; tests never call ~main~. Your ~main~ may remain empty.

## An interface at the consumer
Reuse interfaces from 9-4. Define ~RepairBook~ in ~main.go~ and use it as the parameter type of your command dispatcher to separate command execution from memory storage.
~~~go
type LampBank interface { SetLamp(int) error }
func light(bank LampBank) error {return bank.SetLamp(3)}
~~~
Any struct with the required methods satisfies the interface implicitly. Here ~RepairOrder~ contains ~ID,Units int~, ~Sector string~, and ~Done bool~. A concrete implementation retains a slice and a separate issue counter, like 18-1, now for repair orders with cell counts.
~~~go
id,err:=book.Queue("ice",4)
orders:=book.Orders()
err=book.Complete(id)
err=book.Remove(id)
fmt.Println(orders,err)
~~~
Here ~book~ is a ~RepairBook~. ~Queue~ adds an order, ~Orders~ returns an independent ordered snapshot, ~Complete~ seals a repair, and ~Remove~ deletes it. Failed operations never mutate state or consume IDs. Completing an already completed order succeeds.

## One command or a line stream
Nonempty args execute only that command without reading stdin. Empty args scan lines using ~bufio.Scanner~ from 9-4. ~strings.Fields~ splits Unicode whitespace and returns ~[]string~, so sector names are single tokens here.
~~~go
words:=strings.Fields("add ice 4")
fmt.Println(words) // [add ice 4]
units,err:=strconv.Atoi(words[2])
~~~
~fmt.Fprint~ writes a string without appending a newline; use it when the full reply is already formatted.
~~~go
_,err:=fmt.Fprint(out,"queued=7\n")
~~~
Validate arity before indexing. ~strconv.Atoi~ comes from 4-3; validate the numeric bounds after parsing. A rejected command prints the protocol line ~error~, continues in stream mode, and makes the final code 2 even after later successes.

A non-nil ~Scanner.Err~ or writer failure stops immediately with code 1. Ordinary ~EOF~ succeeds if no command was rejected. Command content is limited to 4096 bytes. Use ~Scanner.Buffer~ from 18-4 with CRLF headroom and check content length. Longer lines are transport failures with code 1, without executing that line.

## Common mistakes
~~~go
id:=len(orders)+1 // Recycles an ID after deletion
os.Exit(2) // Wrong inside run: kills the test process instead of returning
~~~

## Your mission
Implement ~run(args []string, stdin io.Reader, stdout io.Writer) int~. Each call starts a fresh book. Nonempty args execute one command without reading stdin. Empty args process lines through EOF, without prompts or quit. Blank lines are malformed commands. Do not mutate args.

Define and use ~RepairBook interface{Queue(string,int)(int,error); Orders() []RepairOrder; Complete(int) error; Remove(int) error}~. Supply ~RepairOrder~ with the stated fields and independent memory storage; its concrete name is up to you. ~Queue~ accepts a nonempty sector without Unicode whitespace and units in the inclusive range 1..100. IDs start at 1, advance only on successful additions, and are never recycled after deletion. ~Orders~ returns independent snapshots in creation order. ~Complete~ and ~Remove~ reject absent IDs, including zero and negatives, without mutation.

Exact commands and arities: ~add sector units~ prints ~queued=N~; ~done id~ prints ~done=N~; ~rm id~ prints ~removed=N~. ~list~ prints one line per order in creation order using ~%d|%q|%d|%t~ for ID, Sector, Units, Done. Empty lists print nothing. End every line with a newline. Parse IDs and units with Atoi and check bounds. Unknown or malformed commands and domain errors print only ~error~ without mutation.

Return 0 when all commands succeed, 2 if any command is rejected, and 1 on read or write failure even after command errors. Stream mode continues after code 2 and stops immediately on code 1. Prior output and successful mutations before writer failure may remain. Readers and writers are non-nil. Valid command content is at most 4096 bytes, with LF, CRLF, or an unterminated final line; longer content yields code 1. No files, cross-invocation state, or Report.
`, Hints: []string{"Start with memory behind RepairBook: a slice and a separate counter; validate before mutation.", "Separate single-command dispatch from Scanner and reuse it for nonempty args.", "Retain code 2 after malformed commands and continue; Scanner and writer failures stop immediately with 1."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: repairCommandsHidden},
	})

	add(&Level{ID: "21-3-guess-depth", Chapter: 21, Kind: "lab", Title: "خمّن العمق", Goal: "وجّه مسبار التعدين بلعبة تخمين بميزانية محاولات وعشوائية محقونة",
		Requires: []string{"funcvalue"}, Bridge: []int{4}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: guessDepthStarterAr}},
		Lesson: `# المسبار بيدوّر على طبقة الخام

طبقة الخام مستخبية بين عمق 1 و 64 شامل. المسبار عنده 4 محاولات صحيحة؛ كلام مش رقم مش محاولة حفر. هنخلي اختيار الطبقة مدخل محقون، عشان نختبر الفوز والخسارة من غير حظ.

## قبل ما تكتب
- مين يحدد السر: منطق اللعبة ولا مصدر محقون؟
- هل المدخل البايظ يستهلك محاولة، وهل EOF معناها خطأ نقل؟

## مصدر عشوائية على شكل function value
الدوال قيم زي 5-5. ~draw func(int) int~ parameter بياخد حد أعلى ويرجّع رقم من صفر لحد أقل من الحد. مثال مختلف لتحديد بوابة:
~~~go
func chooseGate(draw func(int)int) int {return 10+draw(3)}
fixed:=func(n int)int{return n-1}
fmt.Println(chooseGate(fixed)) // 12
~~~
الـ callback ثابتة في التست. في التشغيل الحقيقي تقدر توصّل method من مولّد محلي:
~~~go
rng:=rand.New(rand.NewSource(7))
gate:=chooseGate(rng.Intn)
~~~
استورد ~math/rand~. ~rand.NewSource(seed)~ مصدر pseudo-random؛ نفس seed تنتج نفس السلسلة. ~rand.New~ تعمل مولّد مستقل، و ~Intn(n)~ ترجع رقم في ~[0,n)~؛ لازم n موجبة. ~rng.Intn~ من غير أقواس method value تحقق توقيع callback. دي عشوائية لعب، مش أسرار أمنية.

في المهمة اطلب ~draw(64)~ مرة واحدة بس، ثم زوّد 1. افحص nil callback وحدود الرقم الراجع قبل القراية أو الإخراج؛ المصدر المخالف للعقد خطأ، مش سبب لتصحيح الرقم في السر. متستخدمش مصدر global أو ساعة الجهاز جوه اللعبة.

## سطر كامل ثم رقم ثم محاولة
~~~go
text:=strings.TrimSpace(" 12 ")
depth,err:=strconv.Atoi(text)
if err!=nil { /* تشخيص invalid من غير زيادة المحاولات */ }
~~~
~bufio.Scanner~ تقرأ سطر، و ~TrimSpace~ تشيل الأطراف، و ~Atoi~ تقرأ integer. السطر ~12 13~ مش رقم واحد. الرقم برا المدى تشخيص مختلف؛ لا النوع البايظ ولا المدى البايظ يستهلك محاولة.

زود عداد المحاولات بعد التحقق، وبعدين افحص التطابق قبل الاستنفاد؛ كده آخر محاولة ممكن تكسب. العمق الأصغر من السر محتاج ~deeper~، والأكبر محتاج ~shallower~. التلميح ميفضحش السر.

## نهاية اللعبة وفشل النقل
~Scanner.Scan~ بترجع false عند EOF أو خطأ؛ افحص ~Scanner.Err~ قبل اعتبارها انسحاب. خطأ writer يرجع فوراً مع سببه زي 21-1. كل diagnostics بروتوكول ثابت، مفيش prompts.

للنجمة التالتة، سطر بايظ طويل لحد 128 KiB لسه تشخيص عادي. ~Scanner.Buffer~ من 18-4 محتاجة مساحة CRLF فوق الحد. افحص طول ~Scanner.Text~ كمان؛ فوق 128 KiB خطأ من غير تشخيص السطر.
~~~go
scanner.Buffer(make([]byte,4096),(128<<10)+2)
~~~

## أخطاء شائعة
~~~go
secret:=1+rand.Intn(64) // تجاهل المصدر المحقون
used++ // غلط لو قبل parsing: كلمة مش رقم استهلكت حفر
~~~

## المهمة
نفّذ ~GuessDepth(in io.Reader, out io.Writer, draw func(int) int) error~. الـ reader والـ writer غير nil. ارفض nil draw، وناتجها لو برا ~[0,64)~، بخطأ من غير I/O. استدعِ ~draw(64)~ مرة واحدة لكل لعبة؛ السر ناتجها زائد 1. مفيش args ولا flags جديدة هنا.

اقرأ سطور لحد النهاية. طبّق TrimSpace ثم Atoi: لو parsing فشل اكتب ~invalid~؛ لو الرقم برا 1..64 اكتب ~range~. الاتنين ميستهلكوش محاولة. كل تخمين صالح يستهلك واحدة من 4. الصح يكتب ~found=N attempts=A~ وينهي فوراً؛ N السر وA عدد المحاولات الصالحة. غير كده اكتب ~deeper~ لو أقل من السر أو ~shallower~ لو أكبر. بعد رابع تخمين غلط اكتب ~lost=N~ وانهي، بعد تلميحه. متعالجش سطور زيادة بعد الفوز أو الاستنفاد.

EOF قبل النهاية تكتب ~aborted~ وترجع nil من غير كشف السر؛ حتى لو مفيش تخمينات. خطأ reader أو سطر أكبر من 128 KiB محتوى يرجّع error من غير aborted؛ خطأ writer يوقف ويحافظ على سببه. إخراج سابق ممكن يفضل. كل سطر إخراج آخره newline، ومفيش غير الصيغ المذكورة. اقبل CRLF وآخر سطر من غير newline. الاختبارات الأساسية تستخدم سطور لحد 4096 بايت؛ النجمة التالتة لسطور بايظة أطول، لحد 128 KiB شامل. مفيش Report ولا إعادة لعب تلقائية.
`, Hints: []string{"اختار السر مرة واحدة من draw(64)+1 بعد فحص المصدر؛ العداد يبدأ صفر.", "TrimSpace ثم Atoi وفحص المدى، وبعدهم بس زوّد المحاولة وافحص الفوز قبل حد 4.", "بعد Scan false افحص Err؛ وسّع Buffer للنجمة التالتة وارجع أول خطأ كتابة."},
		En: &LevelText{Title: "Guess the Depth", Goal: "Guide a mining probe with bounded guesses and injected randomness", Files: []LabFile{{Name: "main.go", Starter: guessDepthStarter}},
			Lesson: `# The Probe Searches for Ore

Ore is hidden at a depth from 1 through 64. The probe has four valid guesses; nonnumeric text is not a drilling attempt. Inject selection of the layer to test wins and losses without luck.

## Before you write
- Who chooses the secret: game logic or an injected source?
- Does invalid input consume an attempt, and does EOF mean transport failure?

## Randomness as a function value
Functions are values as in 5-5. A ~draw func(int) int~ parameter accepts an exclusive upper bound and returns a number from zero below that bound. Here is a different example selecting a gate:
~~~go
func chooseGate(draw func(int)int) int {return 10+draw(3)}
fixed:=func(n int)int{return n-1}
fmt.Println(chooseGate(fixed)) // 12
~~~
Tests use a fixed callback. Real execution can supply a method from a local generator:
~~~go
rng:=rand.New(rand.NewSource(7))
gate:=chooseGate(rng.Intn)
~~~
Import ~math/rand~. ~rand.NewSource(seed)~ creates a pseudorandom source; the same seed reproduces the sequence. ~rand.New~ creates an independent generator and ~Intn(n)~ returns a number in ~[0,n)~ with positive n. ~rng.Intn~ without parentheses is a method value matching the callback. This randomness is suitable for games, rather than security secrets.

Request ~draw(64)~ exactly once, then add 1. Check a nil callback and its returned bounds before reading or writing; an invalid source is an error, not something to silently normalize. Do not use global randomness or the host clock inside the game.

## A whole line, then a number, then an attempt
~~~go
text:=strings.TrimSpace(" 12 ")
depth,err:=strconv.Atoi(text)
if err!=nil { /* Diagnose invalid without incrementing attempts */ }
~~~
~bufio.Scanner~ reads a line, ~TrimSpace~ removes surrounding whitespace, and ~Atoi~ parses an integer. ~12 13~ is not one number. Out-of-range integers have a separate diagnostic; neither invalid category consumes an attempt.

Increment after validation and check equality before exhaustion so the last attempt can win. Depth below the secret needs ~deeper~; above needs ~shallower~. Hints must not reveal the secret.

## Game termination and transport failures
~Scanner.Scan~ returns false for either EOF or errors; check ~Scanner.Err~ before treating that as withdrawal. Return writer errors immediately with their causes as in 21-1. Diagnostics use a fixed protocol without prompts.

For the third star, an invalid line through 128 KiB still receives an ordinary diagnostic. ~Scanner.Buffer~ from 18-4 needs CRLF headroom. Check ~Scanner.Text~ length separately; exceeding 128 KiB is an error without a diagnostic for that line.
~~~go
scanner.Buffer(make([]byte,4096),(128<<10)+2)
~~~

## Common mistakes
~~~go
secret:=1+rand.Intn(64) // Ignores the injected source
used++ // Wrong before parsing: nonnumeric text consumed a drilling attempt
~~~

## Your mission
Implement ~GuessDepth(in io.Reader, out io.Writer, draw func(int) int) error~. The reader and writer are non-nil. Reject nil draw and a result outside ~[0,64)~ with an error before I/O. Call ~draw(64)~ exactly once per game; its result plus 1 is the secret. This lab adds no args or flags.

Read lines until termination. TrimSpace, then Atoi: parsing failures print ~invalid~; numbers outside 1..64 print ~range~. Neither consumes an attempt. Each valid guess consumes one of four attempts. Equality prints ~found=N attempts=A~ and ends immediately, with secret N and valid-attempt count A. Otherwise print ~deeper~ below the secret or ~shallower~ above it. After the fourth wrong guess print ~lost=N~ and stop, after that guess's hint. Do not process further lines after winning or exhaustion.

EOF before termination prints ~aborted~ and returns nil without revealing the secret, even with no guesses. Reader failures or content exceeding 128 KiB return an error without aborted. Writer failures stop immediately and preserve their causes. Prior output may remain. End every output line with a newline; print only the stated forms. Accept CRLF and a final unterminated line. Core tests use lines through 4096 bytes; the third star covers longer invalid lines through exactly 128 KiB. No Report or automatic replay.
`, Hints: []string{"Choose the secret once using validated draw(64)+1; start attempts at zero.", "TrimSpace, Atoi, then bounds; increment only after those and check equality before the four-attempt limit.", "After Scan returns false, check Err; expand Buffer for the third star and return the first write failure."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: guessDepthHidden, Bonus: true},
	})

	add(&Level{ID: "21-4-duel", Chapter: 21, Kind: "lab", Title: "مبارزة الدرونات", Goal: "حوّل أوامر المبارزة لأنواع مسماة واحسب كل الجولات بعشوائية محقونة",
		Requires: []string{"iota"}, Bridge: []int{8}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: droneDuelStarterAr}},
		Lesson: `# تدريب ساحة الدرونات

درون التدريب بيختار مناورة: نبضة ~pulse~، درع ~shield~، أو مثقاب ~drill~. النبضة تعطل المثقاب، المثقاب يخترق الدرع، والدرع يصد النبضة. نفس المناورة تعادل. هنختبر التسع أزواج المرتبة، وبعدين نوصلهم بسيل أوامر.

## قبل ما تكتب
- هتثبت إزاي إن تبديل اللاعب والمنافس يقلب الفوز لخسارة؟
- هل سطر مرفوض أو أمر عرض النتيجة يستهلك اختيار عشوائي؟

## نوع مسمى وقيم iota
الـ custom types و ~iota~ من 2-4؛ هنا بنستخدمهم كقاموس محدود للمناورات بدل strings في المنطق. مثال مختلف لحالة البوابة:
~~~go
type GateState int
const (
 GateClosed GateState=iota
 GateOpen
 GateLocked
)
fmt.Println(GateClosed,GateOpen,GateLocked) // 0 1 2
~~~
~iota~ تبدأ صفر وتزيد في كل سطر const؛ السطور اللي بعدها تكرر التعبير والنوع. في المهمة عرّف ~Maneuver int~ وثوابت ~Pulse~ ثم ~Shield~ ثم ~Drill~ بنفس الأسلوب.

النوع المسمى مش بيمنع قيمة باطلة جاية من تحويل؛ لازم تفحص حدودها قبل حساب النتيجة.
~~~go
state:=GateState(99)
if state<GateClosed || state>GateLocked { /* قيمة غير صالحة */ }
~~~

## النص على الحدود، والأنواع جوه المنطق
~ParseManeuver~ تحول المدخل بعد ~strings.TrimSpace~ و ~strings.ToLower~. متستخدمش substring: ~pulse shield~ مش مناورة واحدة. ~ResolveDuel~ تقبل قيمتين مكتوبين بالنوع المسمى.
~~~go
move,err:=ParseManeuver(" SHIELD ")
verdict,err:=ResolveDuel(move,Pulse)
fmt.Println(verdict,err) // win <nil>
~~~
رجّع خطأ لو أي قيمة مش ضمن التلاتة. مفيش default معناها نبضة. لو الاتنين زي بعض النتيجة ~draw~؛ غير كده افحص الثلاث أزواج اللي اللاعب فيها كسبان، والباقي ~lose~. ترتيب الـ arguments معناه اللاعب ثم المنافس.

## جولات بمصدر محقون
استخدم ~draw func(int) int~ زي 21-3، لكن اسأل عن ~draw(3)~ لكل سطر مناورة صالح بس. حوّل الناتج لـ ~Maneuver~ بعد فحصه؛ 0 نبضة، 1 درع، 2 مثقاب. الاختيار المحقون ثابت في التست ومفيش ساعة أو global randomness.
~~~go
fixed:=func(n int)int{return 2}
rival:=Maneuver(fixed(3))
fmt.Println(rival==Drill) // true
score:=DuelScore{Wins:2,Losses:1,Draws:3}
fmt.Println(score.Wins+score.Losses+score.Draws) // 6
~~~
~DuelScore~ حقولها ~Wins,Losses,Draws int~. زوّد العداد بعد نجاح كتابة الجولة بس. لو الكاتب فشل، ارجع نتيجة الجولات اللي اتكتبت بنجاح قبل الفشل، مع الخطأ وسببه. فشل القراية زي 21-3، بس هنا EOF ترجع النتيجة من غير سطر خروج.

## النجمة التالتة: استعلام النتيجة
السطر ~score~ بعد التطبيع مش جولة؛ اكتب ~score=W/L/D~ من العدادات الحالية، من غير سحب عشوائي أو تغييرها. مثال أول الجلسة ~score=0/0/0~. الاستعلام ممكن يتكرر بين الجولات، ومفيش حالة مشتركة بين استدعاءين.

## أخطاء شائعة
~~~go
rival:=Maneuver(draw(3)) // غلط لو قبل التحقق من سطر اللاعب: استهلك العشوائية للبايظ
if player==rival {return "win",nil} // التعادل مش فوز
~~~

## المهمة
عرّف ~type Maneuver int~ وثوابت ~Pulse=0~ و ~Shield=1~ و ~Drill=2~ باستخدام iota. نفّذ ~ParseManeuver(text string) (Maneuver,error)~: TrimSpace وToLower ثم قبول ~pulse~ أو ~shield~ أو ~drill~ بالظبط؛ الباطل يرجع ~Pulse,error~. نفّذ ~ResolveDuel(player,rival Maneuver) (string,error)~: غير الصالح يرجع string فاضية وخطأ؛ التطابق ~draw~؛ Pulse تكسب Drill، وDrill تكسب Shield، وShield تكسب Pulse؛ غير كده ~lose~.

نفّذ ~Duel(in io.Reader, out io.Writer, draw func(int) int) (DuelScore,error)~. الـ reader والـ writer غير nil. nil draw خطأ قبل I/O. ابدأ العدادات أصفار واقرأ سطر بسطر. أي مدخل مش مناورة يكتب ~invalid~ من غير سحب أو تغيير العدادات، ويكمل. كل مناورة صالحة تستدعي draw(3) مرة واحدة؛ الناتج برا [0,3) خطأ يوقف من غير إخراج للجولة الحالية. اكتب ~player=X rival=Y verdict=V~، X وY أسماء lowercase الأصلية للمناورات وV نتيجة ResolveDuel. بعد نجاح الكتابة زوّد العداد المقابل.

كل سطر إخراج ينتهي newline؛ مفيش prompts أو ملخص عند EOF. اقبل CRLF وآخر سطر من غير newline، وحد المحتوى 4096 بايت شامل؛ الأكبر خطأ من غير invalid للسطر الحالي. أي خطأ قراءة أو كتابة أو مصدر يرجّع النتيجة الجزئية للجولات المكتوبة بنجاح وخطأ يحتفظ بسبب النقل. مفيش معالجة بعد أول فشل نقل، ولا Report.

النجمة التالتة: ~score~ بعد TrimSpace وToLower تكتب ~score=W/L/D~ بالعدادات الحالية، من غير اعتبارها invalid أو جولة ومن غير سحب عشوائي. كل جلسة مستقلة؛ مفيش حد جولات أو إعادة ضبط بعد استعلام النتيجة.
`, Hints: []string{"عرّف Maneuver وثوابتها بـ iota؛ جدول 3×3 يوضح كل نتائج ResolveDuel.", "Parse الأول؛ السطر البايظ يكتب invalid من غير draw. الصالح يسحب draw(3) مرة واحدة.", "اكتب الجولة ثم زوّد العداد؛ score استعلام مستقل، وأخطاء النقل ترجع النتيجة السابقة وسبب الفشل."},
		En: &LevelText{Title: "Drone Duel", Goal: "Map duel commands to named types and score rounds with injected randomness", Files: []LabFile{{Name: "main.go", Starter: droneDuelStarter}},
			Lesson: `# Drone Training Arena

A training drone chooses a ~pulse~, ~shield~, or ~drill~ maneuver. Pulse disables drill, drill pierces shield, and shield blocks pulse. Identical maneuvers draw. Test all nine ordered pairs before connecting them to command input.

## Before you write
- How will swapping player and rival prove that wins become losses?
- Should invalid input or a score query consume a random choice?

## A named type and iota values
Custom types and ~iota~ come from 2-4. Use them as a finite maneuver vocabulary instead of strings inside the rules. Here is a different gate-state example:
~~~go
type GateState int
const (
 GateClosed GateState=iota
 GateOpen
 GateLocked
)
fmt.Println(GateClosed,GateOpen,GateLocked) // 0 1 2
~~~
~iota~ starts at zero and increments on each const specification; following lines repeat the expression and type. Define ~Maneuver int~ and constants ~Pulse~, ~Shield~, then ~Drill~ in this style.

A named type still allows invalid values through conversions; validate bounds before resolving outcomes.
~~~go
state:=GateState(99)
if state<GateClosed || state>GateLocked { /* Invalid value */ }
~~~

## Text at boundaries, types inside the rules
~ParseManeuver~ converts input after ~strings.TrimSpace~ and ~strings.ToLower~. Do not match substrings: ~pulse shield~ is not a maneuver. ~ResolveDuel~ takes two values of the named type.
~~~go
move,err:=ParseManeuver(" SHIELD ")
verdict,err:=ResolveDuel(move,Pulse)
fmt.Println(verdict,err) // win <nil>
~~~
Reject either invalid value; there is no fallback pulse. Equal maneuvers yield ~draw~. Otherwise identify the three winning ordered pairs; the rest yield ~lose~. Arguments are player first, rival second.

## Rounds with an injected source
Reuse ~draw func(int) int~ from 21-3, now requesting ~draw(3)~ only for valid maneuver lines. Validate the returned value before converting to ~Maneuver~: 0 means pulse, 1 shield, 2 drill. Tests control the source without clock access or global randomness.
~~~go
fixed:=func(n int)int{return 2}
rival:=Maneuver(fixed(3))
fmt.Println(rival==Drill) // true
score:=DuelScore{Wins:2,Losses:1,Draws:3}
fmt.Println(score.Wins+score.Losses+score.Draws) // 6
~~~
~DuelScore~ contains ~Wins,Losses,Draws int~. Increment a counter only after successfully writing that round. On writer failure return the scores of previously written rounds and the error with its cause. Handle reader failures as in 21-3, but EOF here returns the score without a farewell line.

## Third star: query the score
A normalized ~score~ line is a query rather than a round. Write ~score=W/L/D~ from current counters, without drawing randomness or changing them. At session start this is ~score=0/0/0~. Queries can repeat between rounds; calls never share state.

## Common mistakes
~~~go
rival:=Maneuver(draw(3)) // Wrong before validating the player line: consumes randomness for invalid input
if player==rival {return "win",nil} // A draw is not a win
~~~

## Your mission
Define ~type Maneuver int~ and constants ~Pulse=0~, ~Shield=1~, ~Drill=2~ using iota. Implement ~ParseManeuver(text string) (Maneuver,error)~: TrimSpace and ToLower, then accept exactly ~pulse~, ~shield~, or ~drill~. Invalid text returns ~Pulse,error~. Implement ~ResolveDuel(player,rival Maneuver) (string,error)~: an invalid value returns an empty string and error; equality yields ~draw~; Pulse beats Drill, Drill beats Shield, and Shield beats Pulse; other pairs yield ~lose~.

Implement ~Duel(in io.Reader, out io.Writer, draw func(int) int) (DuelScore,error)~. Reader and writer are non-nil. Reject nil draw before I/O. Start counters at zero and read lines. Invalid maneuver input prints ~invalid~ without drawing or changing counters, then continues. Each valid maneuver draws exactly once using draw(3). A result outside [0,3) stops with an error and no output for the current round. Print ~player=X rival=Y verdict=V~, using canonical lowercase maneuver names for X and Y and the ResolveDuel verdict for V. Increment the corresponding counter only after a successful write.

End each output line with a newline. No prompts or summary at EOF. Accept CRLF and an unterminated final line. Maximum content is 4096 bytes inclusive; longer lines fail without invalid for that line. Any read, write, or source failure returns the partial score for successfully written rounds and an error preserving transport causes. Stop on the first transport failure. No Report.

Third star: a ~score~ line after TrimSpace and ToLower prints ~score=W/L/D~ from current counters. It is neither invalid nor a round and consumes no randomness. Sessions are independent, with no round cap or reset after a score query.
`, Hints: []string{"Declare Maneuver and iota constants; a 3 by 3 table specifies every ResolveDuel result.", "Parse first; invalid lines print invalid without drawing. Valid lines call draw(3) once.", "Write the round before incrementing; score is an independent query and transport failures return previous scores with their causes."}},
		Lab: &Lab{ReqFile: "main.go", Hidden: droneDuelHidden, Bonus: true},
	})

}

const stationInfoStarter = `package main

import (
	"errors"
	"io"
)

type StationHost struct {
	OS, Arch, Go string
	Lookup       func(string) (string, bool)
}

func StationInfo(args []string, host StationHost, out io.Writer) error {
	return errors.New("station report unfinished")
}
func main() {}
`
const stationInfoStarterAr = `package main

import (
	"errors"
	"io"
)

type StationHost struct {
	OS, Arch, Go string
	Lookup       func(string) (string, bool)
}

func StationInfo(args []string, host StationHost, out io.Writer) error {
	return errors.New("تقرير المحطة لسه مخلصش")
}
func main() {}
`
const stationInfoHidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"slices"
	"testing"
)

func TestHiddenStationFlagDefaultsAndEnvironment(t *testing.T) {
	w := levels.Start("21-1-sysinfo", io.Discard, 211)
	for i := 0; i < 60; i++ {
		osName := fmt.Sprintf("colony-os-%d", w.Rand.Intn(10000))
		arch := fmt.Sprintf("cpu-%d", w.Rand.Intn(100))
		version := fmt.Sprintf("go%d\nquoted\"", w.Rand.Intn(100))
		crew := []string{"", "pilot\n\"قمر\"", fmt.Sprint(w.Rand.Intn(10000))}[i%3]
		present := i%2 == 0
		for _, tc := range []struct {
			args         []string
			system, crew bool
		}{
			{nil, true, false}, {[]string{"-crew"}, true, true},
			{[]string{"-system=false", "-crew=true"}, false, true},
			{[]string{"-crew=false", "-system=false"}, false, false},
			{[]string{"-system=false", "-system=true", "--"}, true, false},
		} {
			calls := 0
			host := StationHost{osName, arch, version, func(name string) (string, bool) {
				calls++
				if name != "COLONY_CREW" {
					t.Fatal("environment allowlist", name)
				}
				return crew, present
			}}
			before := slices.Clone(tc.args)
			var out bytes.Buffer
			if err := StationInfo(tc.args, host, &out); err != nil {
				t.Fatal(err)
			}
			want := ""
			if tc.system {
				want = fmt.Sprintf("os=%q\narch=%q\ngo=%q\n", osName, arch, version)
			}
			wantCalls := 0
			if tc.crew {
				wantCalls = 1
				if present {
					want += fmt.Sprintf("crew=%q\n", crew)
				} else {
					want += "crew=<missing>\n"
				}
			}
			if out.String() != want || calls != wantCalls || !slices.Equal(before, tc.args) {
				t.Fatalf("args=%v output=%q want=%q calls=%d", tc.args, out.String(), want, calls)
			}
		}
	}
}
func TestHiddenStationPreflightAndWriteFailure(t *testing.T) {
	for _, args := range [][]string{{"-unknown"}, {"-crew=bad"}, {"leftover"}, {"-crew", "leftover"}, {"--", "-crew"}, {"-help"}} {
		var out bytes.Buffer
		host := StationHost{Lookup: func(string) (string, bool) { t.Fatal("lookup before validation"); return "", false }}
		if err := StationInfo(args, host, &out); err == nil || out.Len() != 0 {
			t.Fatal("invalid args produced report", args, err, out.String())
		}
	}
	wire := errors.New("wire")
	calls := 0
	host := StationHost{"x", "y", "z", func(string) (string, bool) { calls++; return "", false }}
	for _, failAfter := range []int{0, 1, 2, 3} {
		calls = 0
		out := &hiddenStationWriter{left: failAfter, err: wire}
		if err := StationInfo([]string{"-crew"}, host, out); !errors.Is(err, wire) {
			t.Fatal("write cause", err)
		}
		if failAfter < 3 && calls != 0 {
			t.Fatal("lookup after system write failed")
		}
		if out.failedWrites != 1 {
			t.Fatal("continued writing after failure", out.failedWrites)
		}
	}
}

type hiddenStationWriter struct {
	left, failedWrites int
	err                error
}

func (w *hiddenStationWriter) Write(p []byte) (int, error) {
	if w.left == 0 {
		w.failedWrites++
		return 0, w.err
	}
	w.left--
	return len(p), nil
}
`

const repairCommandsStarter = `package main

import "io"

type RepairOrder struct {
	ID, Units int
	Sector    string
	Done      bool
}

func run(args []string, stdin io.Reader, stdout io.Writer) int {
	return 2
}
func main() {}
`
const repairCommandsHidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"slices"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenRepairSessionStableIDsAndErrors(t *testing.T) {
	w := levels.Start("21-2-subcommands", io.Discard, 212)
	for repeat := 0; repeat < 15; repeat++ {
		var input, want strings.Builder
		units := 1 + w.Rand.Intn(100)
		sector := fmt.Sprintf("ice-%d", w.Rand.Intn(100000))
		fmt.Fprintf(&input, "list\nadd %s %d\nadd relay 100\nadd dock 1\ndone 2\ndone 2\n", sector, units)
		want.WriteString("queued=1\nqueued=2\nqueued=3\ndone=2\ndone=2\n")
		// Each rejected command must leave the final ordered list and next ID intact.
		for _, bad := range []string{"", " ", "unknown", "list extra", "add", "add x 0", "add x 101", "add x nope", "add x 1 extra", "done", "done 0", "done -1", "done 999", "rm 999", "rm x", "rm 1 extra"} {
			input.WriteString(bad + "\n")
			want.WriteString("error\n")
		}
		removed := 1 + w.Rand.Intn(3)
		fmt.Fprintf(&input, "rm %d\nadd new 7\nlist", removed)
		fmt.Fprintf(&want, "removed=%d\nqueued=4\n", removed)
		sectors := []string{sector, "relay", "dock", "new"}
		counts := []int{units, 100, 1, 7}
		for id := 1; id <= 4; id++ {
			if id != removed {
				fmt.Fprintf(&want, "%d|%q|%d|%t\n", id, sectors[id-1], counts[id-1], id == 2)
			}
		}
		var out bytes.Buffer
		if code := run(nil, iotest.OneByteReader(strings.NewReader(input.String())), &out); code != 2 || out.String() != want.String() {
			t.Fatalf("session code=%d output=%q want=%q", code, out.String(), want.String())
		}
	}
}
func TestHiddenRepairSingleCommandsAndFreshSessions(t *testing.T) {
	for _, tc := range []struct {
		args []string
		code int
		want string
	}{
		{[]string{"list"}, 0, ""}, {[]string{"add", "قمر", "+1"}, 0, "queued=1\n"},
		{[]string{"add", "", "1"}, 2, "error\n"}, {[]string{"add", "two words", "1"}, 2, "error\n"},
		{[]string{"add", "ice\u00a0relay", "1"}, 2, "error\n"}, {[]string{"add", "x", "999999999999999999999"}, 2, "error\n"},
		{[]string{"done", "1"}, 2, "error\n"}, {[]string{"rm", "1"}, 2, "error\n"},
		{[]string{"unknown"}, 2, "error\n"},
	} {
		before := slices.Clone(tc.args)
		var out bytes.Buffer
		code := run(tc.args, iotest.ErrReader(errors.New("must not read")), &out)
		if code != tc.code || out.String() != tc.want || !slices.Equal(before, tc.args) {
			t.Fatal("single command", tc.args, code, out.String())
		}
	}
	for _, input := range []string{"", "list\r\n", "add ice 1\r\nlist\r\n"} {
		var out bytes.Buffer
		want := ""
		if strings.HasPrefix(input, "add") {
			want = "queued=1\n1|\"ice\"|1|false\n"
		}
		if code := run(nil, strings.NewReader(input), &out); code != 0 || out.String() != want {
			t.Fatal("clean EOF/CRLF", code, out.String())
		}
	}
}
func TestHiddenRepairTransportFailuresAndLength(t *testing.T) {
	wire := errors.New("broken wire")
	var out bytes.Buffer
	source := io.MultiReader(strings.NewReader("add ice 1\n"), iotest.ErrReader(wire))
	if code := run(nil, source, &out); code != 1 || out.String() != "queued=1\n" {
		t.Fatal("late read failure", code, out.String())
	}
	for _, prefix := range []string{"", "bad\n"} {
		writer := &hiddenRepairWriter{}
		if code := run(nil, strings.NewReader(prefix+"add x 1\nadd y 2\n"), writer); code != 1 || writer.calls != 1 {
			t.Fatal("stop on write error", code, writer.calls)
		}
	}
	for _, ending := range []string{"", "\n", "\r\n"} {
		for _, n := range []int{4096, 4097, 100000} {
			var out bytes.Buffer
			line := "add " + strings.Repeat("x", n-len("add  1")) + " 1"
			code := run(nil, strings.NewReader(line+ending), &out)
			if n == 4096 {
				if code != 0 || out.String() != "queued=1\n" {
					t.Fatal("exact line limit", code)
				}
			} else if code != 1 || out.Len() != 0 {
				t.Fatal("oversized line", n, code)
			}
		}
	}
}

type hiddenRepairWriter struct{ calls int }

func (w *hiddenRepairWriter) Write(p []byte) (int, error) {
	w.calls++
	return 0, errors.New("write failed")
}
`

const guessDepthStarter = `package main

import (
	"errors"
	"io"
)

func GuessDepth(in io.Reader, out io.Writer, draw func(int) int) error {
	return errors.New("depth game unfinished")
}
func main() {}
`
const guessDepthStarterAr = `package main

import (
	"errors"
	"io"
)

func GuessDepth(in io.Reader, out io.Writer, draw func(int) int) error {
	return errors.New("لعبة العمق لسه مخلصتش")
}
func main() {}
`
const guessDepthHidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenDepthSeededWinsLossesAndInvalidGuesses(t *testing.T) {
	w := levels.Start("21-3-guess-depth", io.Discard, 213)
	for i := 0; i < 80; i++ {
		secret := 1 + w.Rand.Intn(64)
		if i == 0 {
			secret = 1
		}
		if i == 1 {
			secret = 64
		}
		wrong := secret%64 + 1
		hint := "shallower\n"
		if wrong < secret {
			hint = "deeper\n"
		}
		for winAt := 1; winAt <= 5; winAt++ {
			input := "bad\n\n12 13\n9999999999999999999999\n0\n65\n"
			want := "invalid\ninvalid\ninvalid\ninvalid\nrange\nrange\n"
			for attempt := 1; attempt <= 4; attempt++ {
				guess := wrong
				if attempt == winAt {
					guess = secret
				}
				input += fmt.Sprintf(" %d \r\n", guess)
				if attempt == winAt {
					want += fmt.Sprintf("found=%d attempts=%d\n", secret, attempt)
					break
				}
				want += hint
				if attempt == 4 {
					want += fmt.Sprintf("lost=%d\n", secret)
				}
			}
			input += "unreachable\n"
			var out bytes.Buffer
			calls := 0
			draw := func(n int) int {
				calls++
				if n != 64 {
					t.Fatal("wrong source bound", n)
				}
				return secret - 1
			}
			if err := GuessDepth(iotest.OneByteReader(strings.NewReader(input)), &out, draw); err != nil || out.String() != want || calls != 1 {
				t.Fatalf("secret=%d winAt=%d got=%q want=%q calls=%d err=%v", secret, winAt, out.String(), want, calls, err)
			}
		}
	}
}
func TestHiddenDepthEOFAndSourceValidation(t *testing.T) {
	for _, tc := range []struct{ input, want string }{
		{"", "aborted\n"}, {"junk", "invalid\naborted\n"},
		{"1", "deeper\naborted\n"}, {"32", "found=32 attempts=1\n"},
		{"33\n33", "shallower\nshallower\naborted\n"},
	} {
		var out bytes.Buffer
		if err := GuessDepth(strings.NewReader(tc.input), &out, func(int) int { return 31 }); err != nil || out.String() != tc.want {
			t.Fatal("EOF", tc, err, out.String())
		}
	}
	for _, draw := range []func(int) int{nil, func(int) int { return -1 }, func(int) int { return 64 }} {
		var out bytes.Buffer
		reader := &hiddenDepthForbiddenReader{t: t}
		if err := GuessDepth(reader, &out, draw); err == nil || out.Len() != 0 {
			t.Fatal("source preflight", err)
		}
	}
}
func TestHiddenDepthReadAndWriteFailures(t *testing.T) {
	wire := errors.New("wire failed")
	source := io.MultiReader(strings.NewReader("1\n"), iotest.ErrReader(wire))
	var out bytes.Buffer
	if err := GuessDepth(source, &out, func(int) int { return 63 }); !errors.Is(err, wire) || out.String() != "deeper\n" {
		t.Fatal("read cause/aborted", err, out.String())
	}
	for _, input := range []string{"", "bad\n", "0\n", "1\n", "64\n", "1\n1\n1\n1\n"} {
		writer := &hiddenDepthWriter{err: wire}
		if err := GuessDepth(strings.NewReader(input), writer, func(int) int { return 63 }); !errors.Is(err, wire) || writer.calls != 1 {
			t.Fatal("stop on writer failure", err, writer.calls)
		}
	}
	for _, size := range []int{(128 << 10) + 1, 256 << 10} {
		var out bytes.Buffer
		if err := GuessDepth(strings.NewReader(strings.Repeat("x", size)), &out, func(int) int { return 0 }); err == nil || out.Len() != 0 {
			t.Fatal("line limit", size, err)
		}
	}
}

type hiddenDepthForbiddenReader struct{ t *testing.T }

func (r *hiddenDepthForbiddenReader) Read(p []byte) (int, error) {
	r.t.Fatal("read before source validation")
	return 0, io.EOF
}

type hiddenDepthWriter struct {
	calls int
	err   error
}

func (w *hiddenDepthWriter) Write(p []byte) (int, error) { w.calls++; return 0, w.err }
func TestHiddenBonusDepthLongInvalidLines(t *testing.T) {
	w := levels.Start("21-3-guess-depth", io.Discard, 2139)
	for _, size := range []int{(64 << 10) + w.Rand.Intn(1000), 128 << 10} {
		for _, ending := range []string{"\n", "\r\n"} {
			var out bytes.Buffer
			input := strings.Repeat("x", size) + ending + "1\n1\n1\n64"
			want := "invalid\ndeeper\ndeeper\ndeeper\nfound=64 attempts=4\n"
			if err := GuessDepth(iotest.HalfReader(strings.NewReader(input)), &out, func(int) int { return 63 }); err != nil || out.String() != want {
				t.Fatal("long invalid line consumed budget or truncated", size, err, out.String())
			}
		}
	}
}
`

const droneDuelStarter = `package main

import (
	"errors"
	"io"
)

type Maneuver int

const (
	Pulse  Maneuver = 0
	Shield Maneuver = 1
	Drill  Maneuver = 2
)

type DuelScore struct{ Wins, Losses, Draws int }

func ParseManeuver(text string) (Maneuver, error) {
	return Pulse, errors.New("duel unfinished")
}
func ResolveDuel(player, rival Maneuver) (string, error) {
	return "", errors.New("duel unfinished")
}
func Duel(in io.Reader, out io.Writer, draw func(int) int) (DuelScore, error) {
	return DuelScore{}, errors.New("duel unfinished")
}
func main() {}
`
const droneDuelStarterAr = `package main

import (
	"errors"
	"io"
)

type Maneuver int

const (
	Pulse  Maneuver = 0
	Shield Maneuver = 1
	Drill  Maneuver = 2
)

type DuelScore struct{ Wins, Losses, Draws int }

func ParseManeuver(text string) (Maneuver, error) {
	return Pulse, errors.New("المبارزة لسه مخلصتش")
}
func ResolveDuel(player, rival Maneuver) (string, error) {
	return "", errors.New("المبارزة لسه مخلصتش")
}
func Duel(in io.Reader, out io.Writer, draw func(int) int) (DuelScore, error) {
	return DuelScore{}, errors.New("المبارزة لسه مخلصتش")
}
func main() {}
`
const droneDuelHidden = `package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenDuelAllOrderedPairsAndParsing(t *testing.T) {
	names := []string{"pulse", "shield", "drill"}
	// Rows are player, columns rival: explicit oracle prevents reversed win rules.
	verdicts := [][]string{{"draw", "lose", "win"}, {"win", "draw", "lose"}, {"lose", "win", "draw"}}
	for player := 0; player < 3; player++ {
		for rival := 0; rival < 3; rival++ {
			verdict, err := ResolveDuel(Maneuver(player), Maneuver(rival))
			if err != nil || verdict != verdicts[player][rival] {
				t.Fatal("ordered pair", player, rival, verdict, err)
			}
			move, err := ParseManeuver(" \t" + strings.ToUpper(names[player]) + "\r\n")
			if err != nil || move != Maneuver(player) {
				t.Fatal("mapping", move, err)
			}
		}
	}
	for _, bad := range []string{"", " ", "rock", "p", "pulse shield", "drilling", "0"} {
		if move, err := ParseManeuver(bad); err == nil || move != Pulse {
			t.Fatal("invalid token default", bad, move, err)
		}
	}
	for _, pair := range [][2]Maneuver{{-1, Pulse}, {3, Shield}, {Drill, -1}, {Pulse, 3}} {
		if verdict, err := ResolveDuel(pair[0], pair[1]); err == nil || verdict != "" {
			t.Fatal("invalid typed value", pair, verdict, err)
		}
	}
}
func TestHiddenDuelSeededStreamAndIndependentScores(t *testing.T) {
	w := levels.Start("21-4-duel", io.Discard, 214)
	names := []string{"pulse", "shield", "drill"}
	verdicts := [][]string{{"draw", "lose", "win"}, {"win", "draw", "lose"}, {"lose", "win", "draw"}}
	for repeat := 0; repeat < 10; repeat++ {
		var input, want strings.Builder
		var samples []int
		var expected DuelScore
		for round := 0; round < 80; round++ {
			if round%7 == 0 {
				input.WriteString("bad\n\n")
				want.WriteString("invalid\ninvalid\n")
			}
			player, rival := w.Rand.Intn(3), w.Rand.Intn(3)
			samples = append(samples, rival)
			input.WriteString(" " + strings.ToUpper(names[player]) + " \r\n")
			verdict := verdicts[player][rival]
			fmt.Fprintf(&want, "player=%s rival=%s verdict=%s\n", names[player], names[rival], verdict)
			switch verdict {
			case "win":
				expected.Wins++
			case "lose":
				expected.Losses++
			case "draw":
				expected.Draws++
			}
		}
		for replay := 0; replay < 2; replay++ {
			calls := 0
			draw := func(n int) int {
				if n != 3 || calls >= len(samples) {
					t.Fatal("draw consumed for invalid input", n, calls)
				}
				sample := samples[calls]
				calls++
				return sample
			}
			var out bytes.Buffer
			score, err := Duel(iotest.OneByteReader(strings.NewReader(input.String())), &out, draw)
			if err != nil || score != expected || out.String() != want.String() || calls != len(samples) {
				t.Fatalf("rounds got=%+v want=%+v calls=%d err=%v", score, expected, calls, err)
			}
		}
	}
}
func TestHiddenDuelEOFAndFailuresKeepWrittenScores(t *testing.T) {
	wire := errors.New("broken wire")
	draw := func(int) int { return 2 }
	var out bytes.Buffer
	score, err := Duel(strings.NewReader(""), &out, draw)
	if err != nil || score != (DuelScore{}) || out.Len() != 0 {
		t.Fatal("empty EOF", score, err)
	}
	score, err = Duel(strings.NewReader("pulse"), &out, draw)
	if err != nil || score != (DuelScore{Wins: 1}) || out.String() != "player=pulse rival=drill verdict=win\n" {
		t.Fatal("final line", score, err)
	}
	out.Reset()
	source := io.MultiReader(strings.NewReader("pulse\n"), iotest.ErrReader(wire))
	score, err = Duel(source, &out, draw)
	if !errors.Is(err, wire) || score != (DuelScore{Wins: 1}) || out.String() != "player=pulse rival=drill verdict=win\n" {
		t.Fatal("late read failure", score, err)
	}
	for _, sample := range []int{-1, 3} {
		var out bytes.Buffer
		calls := 0
		score, err := Duel(strings.NewReader("pulse\npulse\npulse\n"), &out, func(int) int {
			calls++
			if calls == 1 {
				return 2
			}
			return sample
		})
		if err == nil || score != (DuelScore{Wins: 1}) || calls != 2 || out.String() != "player=pulse rival=drill verdict=win\n" {
			t.Fatal("source failure", score, err, calls)
		}
	}
	out.Reset()
	if score, err := Duel(&hiddenDuelForbiddenReader{t: t}, &out, nil); err == nil || score != (DuelScore{}) || out.Len() != 0 {
		t.Fatal("nil source")
	}
	for _, input := range []string{"invalid\n", "pulse\npulse\npulse\n"} {
		writer := &hiddenDuelWriter{err: wire}
		if score, err := Duel(strings.NewReader(input), writer, draw); !errors.Is(err, wire) || score != (DuelScore{}) || writer.calls != 1 {
			t.Fatal("write failure", score, err, writer.calls)
		}
	}
	partialWriter := &hiddenDuelPartialWriter{err: wire}
	score, err = Duel(strings.NewReader("pulse\npulse\npulse\n"), partialWriter, draw)
	if !errors.Is(err, wire) || score != (DuelScore{Wins: 1}) || partialWriter.calls != 2 {
		t.Fatal("failed round counted or processing continued", score, err, partialWriter.calls)
	}
	for _, n := range []int{4096, 4097, 100000} {
		var out bytes.Buffer
		score, err := Duel(strings.NewReader(strings.Repeat("x", n)), &out, draw)
		if score != (DuelScore{}) {
			t.Fatal("invalid line changed score")
		}
		if n == 4096 {
			if err != nil || out.String() != "invalid\n" {
				t.Fatal("exact line limit", err)
			}
		} else if err == nil || out.Len() != 0 {
			t.Fatal("oversize", n, err)
		}
	}
}

type hiddenDuelForbiddenReader struct{ t *testing.T }

func (r *hiddenDuelForbiddenReader) Read(p []byte) (int, error) {
	r.t.Fatal("read before source validation")
	return 0, io.EOF
}

type hiddenDuelWriter struct {
	calls int
	err   error
}

func (w *hiddenDuelWriter) Write(p []byte) (int, error) { w.calls++; return 0, w.err }

type hiddenDuelPartialWriter struct {
	calls int
	err   error
}

func (w *hiddenDuelPartialWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls == 1 {
		return len(p), nil
	}
	return 0, w.err
}

func TestHiddenBonusDuelScoreQueries(t *testing.T) {
	w := levels.Start("21-4-duel", io.Discard, 2149)
	var input, want strings.Builder
	var expected DuelScore
	calls := 0
	names := []string{"pulse", "shield", "drill"}
	for i := 0; i < 120; i++ {
		input.WriteString(" SCORE \r\nscore\n")
		line := fmt.Sprintf("score=%d/%d/%d\n", expected.Wins, expected.Losses, expected.Draws)
		want.WriteString(line + line)
		player := w.Rand.Intn(3)
		input.WriteString(names[player] + "\n")
		verdict := []string{"win", "lose", "draw"}[player]
		fmt.Fprintf(&want, "player=%s rival=drill verdict=%s\n", names[player], verdict)
		switch verdict {
		case "win":
			expected.Wins++
		case "lose":
			expected.Losses++
		case "draw":
			expected.Draws++
		}
	}
	input.WriteString("score")
	fmt.Fprintf(&want, "score=%d/%d/%d\n", expected.Wins, expected.Losses, expected.Draws)
	var out bytes.Buffer
	score, err := Duel(strings.NewReader(input.String()), &out, func(n int) int {
		calls++
		if n != 3 {
			t.Fatal("bound")
		}
		return 2
	})
	if err != nil || score != expected || out.String() != want.String() || calls != 120 {
		t.Fatal("score queries must not draw or reset", score, expected, calls, err)
	}
	out.Reset()
	score, err = Duel(strings.NewReader("score"), &out, func(int) int { t.Fatal("score query drew"); return 0 })
	if err != nil || score != (DuelScore{}) || out.String() != "score=0/0/0\n" {
		t.Fatal("fresh query", score, err)
	}
}
`
