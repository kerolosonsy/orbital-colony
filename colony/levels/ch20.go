package levels

func init() {
	add(&Level{ID: "20-1-bulletin", Chapter: 20, LessonRevision: 1, Steps: ch20BulletinArabic, Title: "نشرة المستعمرة", Goal: "حوّل نشرة محدودة لـ HTML آمن بماكينة حالات", Requires: []string{"sel:strings.Builder"}, Bridge: []int{23},
		Lesson: `# نشرة آمنة على شاشة القاعدة

## قبل ما تكتب
- إزاي تفرّق نص المستخدم عن tags اللي إنت بتولدها؟
- هتعمل إيه لو الملف خلص ولسه tag مفتوحة؟

## ماكينة حالات صغيرة
نشرة المستعمرة ليها tagين صريحين، مش Markdown عام. الحالة نص عادي أو alert أو code. الفتح يغيّر الحالة، والقفل المطابق بس يرجّعها للنص. ارفض tag معروفة متداخلة أو قفل من غير فتح.
~~~go
mode:="plain"
if strings.HasPrefix(rest,"[notice]") {mode="notice"}
fmt.Println(mode)
~~~
~strings.HasPrefix~ بتفحص بداية النص؛ slice من index بتدي باقي المدخل. استهلك delimiter المعروفة كاملة قبل النص العادي. الأقواس المجهولة تفضل حرفية. notice في المثال مش من قواعد المهمة.

## ابني الناتج من غير ثقة في المدخل
استورد ~strings~ و~html~. ~strings.Builder~ بتجمع النص من غير نسخ البادئة كلها كل مرة. قيمتها الصفرية جاهزة. ~WriteString~ تضيف و~String~ تقرا الناتج؛ متنسخش Builder غير فاضية.
~~~go
var out strings.Builder
out.WriteString("<em>")
out.WriteString(html.EscapeString("<relay>&\"'"))
out.WriteString("</em>")
fmt.Println(out.String())
~~~
~html.EscapeString~ بتستبدل & و< و> والـquotes المزدوجة والمفردة بـentities آمنة. طبقها مرة واحدة على نص المستخدم حتى جوه code، واكتب tags بتاعتك مباشرة. نص شكله entity زي &lt; لازم يبقى &amp;lt;.

newline العادي يتحوّل HTML break؛ جوه tag يفضل newline حرفي. قراءة bytes مناسبة لقواعد ASCII دي: انقل كل byte تانية زي ما هي مع escaping. متعتبرش byte حرف Unicode كامل.
~~~go
piece:=source[i:i+1]
out.WriteString(html.EscapeString(piece))
~~~
لو tag مفتوحة عند EOF رجّع ناتج فاضي وخطأ حتى لو Builder فيها بادئة آمنة. النجمة التالتة نص Unicode كبير متعدد السطور ونفس سياسة الفشل.

## أخطاء شائعة
~~~go
out.WriteString(userText) // HTML خبيث شغّال
return out.String(),nil // تجاهلت حالة مفتوحة عند EOF
~~~

## المهمة
نفّذ ~RenderBulletin(source string) (string,error)~ باستخدام strings.Builder. القواعد الوحيدة ~[alert]text[/alert]~ تتحوّل <strong>text</strong> و~[code]text[/code]~ تتحوّل <code>text</code>. المحتوى ممكن فاضي؛ tags المعروفة ممنوع تتداخل حتى جوه code. أي قفل غير مطابق أو من غير فتح أو EOF مفتوحة يرجّع "",error من غير HTML جزئي. أي syntax تانية حرفية. طبّق html.EscapeString مرة واحدة على كل نص المستخدم؛ newline برا tags تبقى ~<br>
~ وجواها تفضل newline. احفظ باقي المسافات والبايتات، من غير newline إضافية. الفاضي يرجّع "",nil. النجمة التالتة نص كبير Unicode وmultiline، مفيش حد Scanner. مفيش Report.
`, Hints: []string{"خلي الحالة واضحة: نص عادي أو جوه tag.", "اكتب HTML بتاعك واعمِل escaping مرة واحدة لكل نص.", "عند EOF، راجع هل فيه tag مفتوحة قبل إرجاع الناتج."},
		En: &LevelText{LessonRevision: 1, Steps: ch20BulletinEnglish, Title: "Colony Bulletin", Goal: "Render a small bulletin grammar as safe HTML with a state machine", Lesson: `# A Safe Bulletin on the Base Screen

## Before you write
- How will you distinguish user text from tags you generate?
- What happens when EOF arrives inside an open tag?

## A small state machine
The colony bulletin uses two explicit tags rather than general Markdown. State is either plain text, alert, or code. Opening a tag changes state; only its matching close returns to plain. Reject nested recognized tags and unexpected closes.
~~~go
mode:="plain"
if strings.HasPrefix(rest,"[notice]") {mode="notice"}
fmt.Println(mode)
~~~
~strings.HasPrefix~ tests the start of a string; slicing from an index gives the remaining input. Consume an entire recognized delimiter before examining ordinary text. Unknown bracket syntax stays literal. This example's notice tag is not part of the mission grammar.

## Build output without trusting input
Import ~strings~ and ~html~. ~strings.Builder~ accumulates text without repeatedly copying the entire prefix. Its zero value is ready. ~WriteString~ appends text and ~String~ reads the accumulated result; do not copy a nonempty Builder.
~~~go
var out strings.Builder
out.WriteString("<em>")
out.WriteString(html.EscapeString("<relay>&\"'"))
out.WriteString("</em>")
fmt.Println(out.String())
~~~
~html.EscapeString~ replaces &, <, >, double quotes and single quotes with safe entities. Apply it once to user text, including text inside code blocks; emit your own tags directly. Existing entity-looking text such as &lt; must become &amp;lt;.

Plain newlines become HTML breaks; inside a tag they remain literal newlines. Scanning bytes is safe for this ASCII grammar: copy every other byte unchanged through escaping. Do not interpret a byte as a complete Unicode character.
~~~go
piece:=source[i:i+1]
out.WriteString(html.EscapeString(piece))
~~~
On an unclosed tag at EOF return empty output and an error, even if the builder already contains a safe prefix. The third star exercises large multiline Unicode content and the same failure policy.

## Common mistakes
~~~go
out.WriteString(userText) // Executes hostile HTML
return out.String(),nil // Ignores an open state at EOF
~~~

## Your mission
Implement ~RenderBulletin(source string) (string,error)~ using strings.Builder. Only ~[alert]text[/alert]~ becomes <strong>text</strong> and ~[code]text[/code]~ becomes <code>text</code>. Content may be empty; recognized tags cannot nest, including inside code. Mismatched or unexpected closes and EOF inside a tag return "",error without partial HTML. All other syntax is literal. Apply html.EscapeString exactly once to every user-text fragment. Newlines outside tags become ~<br>
~; inside tags retain literal newlines. Preserve all other whitespace and bytes; add no trailing newline. Empty input returns "",nil. The third star covers large Unicode multiline text without a Scanner limit. No Report.
`, Hints: []string{"Track plain text or tag content as state.", "Write your HTML and escape each user-text fragment once.", "At EOF, check for an open tag before returning output."}, Files: []LabFile{{Name: "main.go", Starter: bulletinStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: bulletinStarterAr}}, Lab: &Lab{ReqFile: "main.go", Hidden: bulletinHidden, Bonus: true}})
	add(&Level{ID: "20-2-validators", Chapter: 20, LessonRevision: 1, Steps: ch20ValidatorArabic, Title: "المدققين", Goal: "افحص بطاقة الطاقم واجمع الـerrors مع الاحتفاظ بأسبابها", Requires: []string{"errwrap"}, Bridge: []int{14},
		Lesson: `# بطاقة الطاقم تتراجع كلها

## قبل ما تكتب
- إزاي تعرض كل الأخطاء مش أول خطأ بس؟
- هل regex للتاريخ تثبت إن اليوم موجود؟

## فحوص مستقلة وحدود regex
تسجيل الطاقم3 فحوص مستقلة. القواعد الصغيرة دي مش معيار إيميل أو تليفون دولي كامل. استورد ~regexp~: ~regexp.MustCompile~ بتجهز pattern ثابتة من المبرمج مرة واحدة وتعمل panic لو هي غلط؛ ~Regexp.MatchString~ تفحص مدخل المستخدم من غير panic.
~~~go
var badge=regexp.MustCompile("^[A-Z]{2}[0-9]{3}$")
fmt.Println(badge.MatchString("AB123")) // true
~~~
~^~ و~$~ أول وآخر النص، brackets مجموعة حروف، ~{3}~ تكرار3 مرات، ~+~ مرة أو أكتر، والأقواس بتجمع التكرار. النقطة الحرفية regex هي ~\.~، وفي Go quoted string اكتب ~\\.~. ~{7,12}~ مدى شامل.
~~~go
var domain=regexp.MustCompile("^[A-Za-z0-9]+(\\.[A-Za-z0-9]+)+$")
var radio=regexp.MustCompile("^\\+[0-9]{7,12}$")
fmt.Println(domain.MatchString("relay.colony"),radio.MatchString("+1234567"))
~~~
العنوان يتقسم على @ لازم جزئين بالظبط. الجزء المحلي يقبل Unicode من غير أي Unicode whitespace. فحص whitespace ASCII مش كفاية؛ استخدم ~unicode.IsSpace~ على runes، امتداد لفحوص Unicode من18-4.
~~~go
fmt.Println(unicode.IsSpace('\u00a0')) // true
~~~

## تاريخ حقيقي وأخطاء مجمعة بتصنيف
استخدم time.Parse من19-1. parser حقيقي يرفض تاريخ مستحيل؛ مقارنة Format بالمدخل تفرض الأصفار في اليوم والشهر.
~~~go
t,err:=time.Parse("2006/01/02","2032/02/29")
fmt.Println(t.Format("2006/01/02"),err)
~~~
كل validator تلف sentinel بتاعتها بـ%w من5-3. ~errors.Join~ تجمع كل الأسباب غيرnil؛ لو كلهم nil ترجع nil. ~errors.Is~ تلاقي sentinel جوه اللف والتجميع، بدل مقارنة نص الخطأ.
~~~go
err:=errors.Join(fmt.Errorf("badge: %w",ErrRadio),ErrAddress)
fmt.Println(errors.Is(err,ErrRadio)) // true
~~~
رجّع أسماء الحقول بترتيب معلن عن طريق فحص3 sentinels بالترتيب؛ متعتمدش على map. InvalidFields تقبل nil وخطأ مش تابع للفحوص.

## أخطاء شائعة
~~~go
if err!=nil {return err} // بيوقف باقي فحوص البطاقة
regexp.MustCompile("[0-9]+").MatchString(phone) // substring مش عقد كامل
~~~

## المهمة
نفّذ CheckAddress وCheckRadio وCheckLaunchDate: كل واحدة (string) error ومستقلة؛ الصالح nil والباطل خطأ ملفوف حوالين ErrAddress أو ErrRadio أو ErrLaunchDate. العنوان @ واحدة، local غير فاضي ومسموح Unicode من غير unicode.IsSpace؛ domain مطابق ~[A-Za-z0-9]+(\.[A-Za-z0-9]+)+~ بالكامل. مفيش trim ولا قواعد RFC إضافية. الراديو + وبعدها7 لـ12 رقم ASCII بالظبط من غير فواصل. التاريخ صالح ومطابق تماماً ~2006/01/02~. ValidateCrewCard(CrewCard) error تجمع التلات فحوص بـerrors.Join حتى لو أكتر من حقل فشل. InvalidFields(error) []string تستخدم errors.Is وترجع كل التصنيفات المعروفة مرة واحدة بترتيب ~address,radio,launchDate~؛ nil أو خطأ غريب يرجع طول0. مدخلات العنوان UTF-8 صالحة؛ نص الخطأ حر، السبب جزء من العقد. مفيش Report.
`, Hints: []string{"خلي كل validator تفحص حقل واحد وترجّع sentinel السبب.", "استخدم %w وJoin، وبعدها errors.Is للتصنيف.", "افحص الحقول بالترتيب الثابت المعلن."},
		En: &LevelText{LessonRevision: 1, Steps: ch20ValidatorEnglish, Title: "Validators", Goal: "Validate a crew enrollment card and aggregate classified causes", Lesson: `# Review the Whole Crew Card

## Before you write
- How will you expose every failure rather than only the first?
- Does a date regex prove the day exists?

## Independent validators and regex limits
Crew enrollment has three independent checks. These deliberately small rules are not a complete email or international phone standard. Import ~regexp~: ~regexp.MustCompile~ compiles a programmer-owned pattern once, panicking for an invalid pattern; ~Regexp.MatchString~ checks user input without panicking.
~~~go
var badge=regexp.MustCompile("^[A-Z]{2}[0-9]{3}$")
fmt.Println(badge.MatchString("AB123")) // true
~~~
~^~ and ~$~ anchor the whole string, brackets select a character class, ~{3}~ requires three repetitions, ~+~ means one or more, and parentheses group repetitions. A literal dot in a regex is ~\.~; in a Go quoted string write ~\\.~. ~{7,12}~ is an inclusive repetition range.
~~~go
var domain=regexp.MustCompile("^[A-Za-z0-9]+(\\.[A-Za-z0-9]+)+$")
var radio=regexp.MustCompile("^\\+[0-9]{7,12}$")
fmt.Println(domain.MatchString("relay.colony"),radio.MatchString("+1234567"))
~~~
For the address, split on @ and require exactly two parts. The local part may contain Unicode but no Unicode whitespace. ASCII regex whitespace is insufficient; use ~unicode.IsSpace~ on each rune, extending the Unicode predicates from 18-4.
~~~go
fmt.Println(unicode.IsSpace('\u00a0')) // true
~~~

## Real dates and classified aggregate errors
Reuse time.Parse from 19-1. A real parser rejects impossible dates, while a format round trip requires exact zero padding.
~~~go
t,err:=time.Parse("2006/01/02","2032/02/29")
fmt.Println(t.Format("2006/01/02"),err)
~~~
Each validator wraps its sentinel with %w from 5-3. ~errors.Join~ combines every non-nil cause; all nil produces nil. ~errors.Is~ finds sentinels through wrapping and joins, rather than comparing error wording.
~~~go
err:=errors.Join(fmt.Errorf("badge: %w",ErrRadio),ErrAddress)
fmt.Println(errors.Is(err,ErrRadio)) // true
~~~
Return field names in declared order by checking the three sentinels in order; never depend on a map's iteration order. InvalidFields also accepts nil and unrelated errors.

## Common mistakes
~~~go
if err!=nil {return err} // Stops the remaining card checks
regexp.MustCompile("[0-9]+").MatchString(phone) // Substring, not whole contract
~~~

## Your mission
Implement independent ~CheckAddress~, ~CheckRadio~, and ~CheckLaunchDate~, each (string) error: valid input returns nil; invalid input wraps ErrAddress, ErrRadio, or ErrLaunchDate respectively. Address has exactly one @, a nonempty Unicode local part without unicode.IsSpace, and a domain matching ~[A-Za-z0-9]+(\.[A-Za-z0-9]+)+~ in full. No trimming or additional RFC rules. Radio is + followed by exactly 7 through 12 ASCII digits without separators. Date is real and exactly matches ~2006/01/02~. ~ValidateCrewCard(CrewCard) error~ joins all three checks with errors.Join even when several fail. ~InvalidFields(error) []string~ uses errors.Is to return every known classification once in order ~address,radio,launchDate~; nil and unrelated errors yield length zero. Address inputs are valid UTF-8. Error wording is free; causes are part of the contract. No Report.
`, Hints: []string{"Give each validator one field and its sentinel cause.", "Use %w and Join, then errors.Is for classification.", "Check fields in their declared fixed order."}, Files: []LabFile{{Name: "main.go", Starter: crewCardStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: crewCardStarterAr}}, Lab: &Lab{ReqFile: "main.go", Hidden: crewCardHidden}})
	add(&Level{ID: "20-3-lexer", Chapter: 20, LessonRevision: 1, Steps: ch20LexerArabic, Title: "مقطّع الأوامر", Goal: "قطّع لغة طيران صغيرة مع offsets دقيقة للـerrors", Requires: []string{"struct", "switch"}, Bridge: []int{93},
		Lesson: `# قبل ما الدرون يفهم الأمر

## قبل ما تكتب
- هتقيس الموقع بالـbytes ولاrunes؟
- لو رقم طويل عمل overflow هترجع token غلط؟

## الـ token تحفظ دليل المصدر
قبل تنفيذ الطيران حوّل النص لـtokens. كل token نوع وتهجئة من المصدر وoffset بالبايت من صفر. Offset مش عمود runes. اللغة ASCII وبترفض byte غيرASCII بدل تجاهلها.
~~~go
token:=CommandToken{Kind:"name",Text:"cargo",Offset:4}
err:=&CommandError{Offset:10,Message:"unexpected byte"}
fmt.Println(token.Text,err.Error())
~~~
~error~ هي interface جاهزة ~interface{Error() string}~. pointer receiver من فصل8 تخلي *CommandError تحققها؛ رجوع pointer بيحفظ الحقول. لف الأخطاء من5-3، لكن ~errors.As~ جديدة هنا وبتطلع سبب بنوع محدد. استورد ~errors~، عرّف pointer target وابعت عنوانها عشان As تملاها. الـstarter فيها method دي جاهزة، وتنسيق fmt.Sprintf من9-1.
~~~go
func(e *CommandError)Error() string {
 return fmt.Sprintf("byte %d: %s",e.Offset,e.Message)
}
var failure error=&CommandError{Offset:3,Message:"bad byte"}
fmt.Println(failure.Error())
var diagnostic *CommandError
if errors.As(failure,&diagnostic) {fmt.Println(diagnostic.Offset)}
~~~ الأنواع الجاهزة مشتركة مع المهمة اللي بعدها، فحافظ على أسمائها وحقولها المصدّرة. Offset بطول المصدر معناها EOF.

## سلاسل كاملة وأرقام محدودة
افحص الـbyte الجاية بـswitch من فصل4. استهلك أطول سلسلة اسم أو أرقام؛ متطلعش token لكل حرف. تجاهل space وtab وCR وLF بس. أي byte مجهولة تفشل في مكانها.
~~~go
start:=i
for i<len(source) && source[i]>='0' && source[i]<='9' {i++}
text:=source[start:i]
n,err:=strconv.ParseInt(text,10,64)
fmt.Println(n,err)
~~~
~strconv.ParseInt~ جديدة: غير Atoi بترجع int64 وبتاخد base10 وحجم64 صراحة. افحص خطأ overflow وبعدين حد المهمة1000. السالب token منفصلة، مش جزء من الرقم.
~~~go
symbol:=CommandToken{Kind:"-",Text:"-",Offset:0}
end:=CommandToken{Kind:"eof",Text:"",Offset:len(source)}
fmt.Println(symbol,end)
~~~
احفظ التهجئة حتى الأصفار في الأول. eof واحدة حتى للنص الفاضي. أي فشل يرمي كل tokens؛ مفيش برنامج جزئي شكله صالح.

## أخطاء شائعة
~~~go
n,_:=strconv.Atoi(text) // تجاهل overflow
i++ // تجاهل byte مجهولة بدل تشخيصها
~~~

## المهمة
نفّذ ~LexFlight(source string) ([]CommandToken,error)~. المصدر أكبر من4096 بايت يرجّع nil,*CommandError عند4096 قبل القراية. الأسماء ~[a-z][a-z0-9_]*~؛ كل الأسماء حتى let وthrust والاتجاهات Kind ~name~. سلسلة أرقام قيمة عشرية0..1000 ونوع ~number~ مع Text الأصلية بما فيها أصفار البداية. الرموز ~+ - * / ( ) = ;~ كل واحدة Kind وText نفس الرمز. تجاهل space وtab وCR وLF ASCII بس؛ مفيش comments. ضيف eof واحدة ~CommandToken{Kind:"eof",Text:"",Offset:len(source)}~. كل offsets بايت من صفر. byte مجهولة حتى uppercase أوUnicode أوUTF-8 بايظة ترجّع nil,*CommandError في أول byte؛ الرقم الباطل أوoverflow عند بدايته. Message مش فاضية. الـlexer مش بتفرض قواعد الجمل؛ مسموح number جنب name. مفيش تنفيذ ولا Report.
`, Hints: []string{"امشِ على المصدر بـindex للـbyte.", "خزّن بداية ونص كل run بالكامل.", "استخدم ParseInt وافحص error قبل قبول الرقم."},
		En: &LevelText{LessonRevision: 1, Steps: ch20LexerEnglish, Title: "Command Lexer", Goal: "Tokenize a tiny flight language with precise error offsets", Lesson: `# Before the Drone Understands

## Before you write
- Will positions count bytes or runes?
- Will an overflowing numeral become a misleading token?

## Tokens retain source evidence
Before executing a flight program, turn text into tokens. A token has a kind, the exact source spelling, and a zero-based byte offset. Offset is a byte position, not a rune column. This ASCII language rejects non-ASCII bytes instead of silently skipping them.
~~~go
token:=CommandToken{Kind:"name",Text:"cargo",Offset:4}
err:=&CommandError{Offset:10,Message:"unexpected byte"}
fmt.Println(token.Text,err.Error())
~~~
~error~ is the built-in interface ~interface{Error() string}~. A pointer receiver from chapter 8 lets *CommandError satisfy it; returning that pointer preserves its fields. Chapter 5-3 introduced error wrapping; ~errors.As~ is new here and extracts a typed cause. Import ~errors~, declare a target pointer, and pass its address so As can fill it. The starter supplies this method. Its fmt.Sprintf formatting was taught in 9-1.
~~~go
func(e *CommandError)Error() string {
 return fmt.Sprintf("byte %d: %s",e.Offset,e.Message)
}
var failure error=&CommandError{Offset:3,Message:"bad byte"}
fmt.Println(failure.Error())
var diagnostic *CommandError
if errors.As(failure,&diagnostic) {fmt.Println(diagnostic.Offset)}
~~~ The supplied types are shared with the next mission, so retain their exported names and fields. Offset len(source) means EOF.

## Maximal runs and bounded integers
Inspect the next byte using switch from chapter 4. Consume maximal name or digit runs; do not emit a token per character. Skip only space, tab, CR, and LF. All other unknown bytes fail at their own offset.
~~~go
start:=i
for i<len(source) && source[i]>='0' && source[i]<='9' {i++}
text:=source[start:i]
n,err:=strconv.ParseInt(text,10,64)
fmt.Println(n,err)
~~~
~strconv.ParseInt~ is new: unlike Atoi, it returns int64 and accepts explicit base 10 and bit size 64. Check the parser error for overflow, then the mission's literal bound of 1000. Unary minus is a separate token, never part of a number.
~~~go
symbol:=CommandToken{Kind:"-",Text:"-",Offset:0}
end:=CommandToken{Kind:"eof",Text:"",Offset:len(source)}
fmt.Println(symbol,end)
~~~
Keep spelling, including leading zeros. Emit exactly one eof token even for empty input. Any failure discards all tokens; it never returns a plausible partial program.

## Common mistakes
~~~go
n,_:=strconv.Atoi(text) // Discards overflow
i++ // Skips an unknown byte without a diagnostic
~~~

## Your mission
Implement ~LexFlight(source string) ([]CommandToken,error)~. Source longer than 4096 bytes returns nil,*CommandError at offset 4096 before scanning. Names match ~[a-z][a-z0-9_]*~; every name including let, thrust and directions has Kind ~name~. Digit runs are decimal magnitudes 0..1000, Kind ~number~, with original Text (leading zeros allowed). Symbols ~+ - * / ( ) = ;~ each have Kind and Text equal to that symbol. Skip only ASCII space, tab, CR and LF. No comments. Append exactly ~CommandToken{Kind:"eof",Text:"",Offset:len(source)}~. All offsets are zero-based bytes. An unknown byte (including uppercase, Unicode or invalid UTF-8) returns nil,*CommandError at its first byte; invalid/overflowing numbers point at their start. Message must be nonempty. Lexing does not enforce statement grammar: adjacent number/name tokens are allowed. Do not execute or Report.
`, Hints: []string{"Scan the source with a byte index.", "Keep the start and full text of each run.", "Use ParseInt and check its error before accepting a number."}, Files: []LabFile{{Name: "main.go", Starter: flightLexerStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: flightLexerStarterAr}}, Lab: &Lab{ReqFile: "main.go", Hidden: flightLexerHidden}})
	add(&Level{ID: "20-4-evaluator", Chapter: 20, LessonRevision: 1, Steps: ch20EvaluatorArabic, Title: "المنفّذ", Goal: "فسّر برنامج الطيران بأولويات ومتغيرات وحرّك الدرون فعلاً", Requires: []string{"recursion", "iface"}, Bridge: []int{93},
		Lesson: `# من النص لحركة حقيقية

## قبل ما تكتب
- إزاي تمنع برنامج فيه syntax غلط من تنفيذ أول حركة؟
- الخطأ عند العملية ولا عند بداية الجملة؟

## القواعد بتحدد الأولوية
كمّل بالـ lexer بتاعتك في ~main.go~ وضيف ~flight.go~ للمفسر. recursive descent دالة لكل مستوى أولوية، بتستهلك tokens بمؤشر.

AST شجرة syntax مجردة بتحفظ البنية لتقييم لاحق.
~~~text
program   = { statement }
statement = "let" name "=" expression ";" | "thrust" direction expression ";"
expression = term { ("+" | "-") term }
term       = primary { ("*" | "/") primary }
primary    = number | name | "(" expression ")" | ("+" | "-") primary
~~~
الضرب أعمق من الجمع، فـ ~2+3*4~ بتتجمع ~2+(3*4)~. الحلقة تخلي نفس المستوى left associative: ~9-3-1~ هي ~(9-3)-1~.

العلامة الأحادية تنادي ~primary~ تاني مباشرة. خليها دالة حرة ~primary(p *parser)~، وبعد استهلاك العلامة نادِ ~primary(p)~؛ فاحص المفاهيم بيتعرف على recursion المباشرة للدوال.

الأقواس تنادي ~expression~ تاني، بنفس recursion من 14-5.
~~~go
func primary(p *parser)(FlightExpr,error) {
 if p.peek().Kind=="(" {
  p.next++
  expr,err:=p.expression()
  if err!=nil {return nil,err}
  if _,err=p.take(")");err!=nil {return nil,err}
  return expr,nil
 }
 if p.peek().Kind=="+" {
  p.next++
  return primary(p)
 }
 return p.atom()
}
~~~
ده sketch بيفترض ~parser~ فيها ~tokens~ و~next~ ودوال ~peek~ و~take~ و~expression~ و~atom~؛ اختار الأسماء الداخلية بنفسك.

~take~ ترفض token الحالية من غير تعدي ~EOF~. delimiter الناقصة موقعها token الحالية حتى ~eof~. primaries جنب بعض خطأ؛ مفيش ضرب ضمني.

## interface لشجرة التعبير
~FlightExpr~ بتشيل أي node عندها ~Eval(map[string]int64)(int64,error)~. interfaces من 9-4 وmethods من فصل 7: ~literal~ و~variable~ و~binary~ كلهم ينفّذوا نفس العقد.

~map~ هي environment لاستدعاء واحد.
~~~go
type literal struct {n int64}
func(n literal)Eval(env map[string]int64)(int64,error){return n.n,nil}
var expr FlightExpr=literal{7}
env:=map[string]int64{"fuel":3}
n,err:=expr.Eval(env)
fmt.Println(n,err) // 7 <nil>
~~~
احفظ token العملية في node. ~binary~ تقيّم الشمال وبعدين اليمين recursively وبعدها العملية.

~variable~ تستخدم ~bool~ الوجود في ~map~؛ الغايبة خطأ عند اسمها. السالب الأحادي ممكن node طرح شمالها صفر ومكانها علامة السالب.
~~~go
n,ok:=env[name]
if !ok {return 0,&CommandError{Offset:offset,Message:"undefined variable"}}
~~~
القيم ~int64~ من 20-3. كل ناتج وسيط لازم بين -1000000 و 1000000. ضرب قيمتين صالحين أقصاه ~10^12~ فحسابه قبل الفحص آمن في ~int64~.

افحص القسمة على صفر قبل القسمة؛ قسمة الأعداد الصحيحة بتحذف الكسر ناحية الصفر. خطأ المدى أو القسمة عند العملية.

## حلّل الأول وبعدها حرّك
ابني قايمة الجمل كلها قبل أي تقييم عشان syntax غلط متحركش الدرون.

قيّم يمين ~let~ بنجاح قبل تحديث environment. كل استدعاء environment جديدة. runtime error توقف الجمل اللي بعدها بس الحركات اللي خلصت تفضل.
~~~go
type baseDrive struct{}
func(baseDrive)Step(dir colony.Direction)bool{return colony.Move(dir)}
err:=ExecuteFlight("let cells=2; thrust north cells;",baseDrive{})
~~~
~FlightDrive~ interface عند المستهلك بدالة ~Step(colony.Direction) bool~. ~false~ معناها مسدود.

adapter المثال بتحرك الدرون الحقيقي؛ integration بتفحص مكانه، وباقي التستات boundary محقونة. ~colony.Direction~ و~colony.Move~ والاتجاهات الأربع من فصل 1. اربط أسماء الاتجاهات بالثوابت.

~thrust~ تتحقق من العدد كله وبعدين ~Step~ لكل خلية؛ الصفر مفيش نداء.

لو ~Step~ فشلت قف فوراً. الخلايا الناجحة السابقة اتحركت فعلاً ومفيش rollback.

~nil~ interface خطأ عادي؛ typed nil pointer برا العقد.

## أخطاء شائعة
~~~go
env[name]=n // قبل التأكد إن تقييم اليمين نجح
left=left+right // تطبيق الجمع فوراً ضيّع أولوية الضرب
~~~

## المهمة
خلي ~LexFlight~ وأنواعها زي ما هي في ~main.go~.

نفّذ ~ExecuteFlight(source string,drive FlightDrive) error~ في ~flight.go~، مع ~FlightDrive interface{Step(colony.Direction) bool}~ وAST باستخدام ~FlightExpr interface{Eval(map[string]int64)(int64,error)}~.

اتبع القواعد فوق؛ الاتجاه ~east~ أو ~north~ أو ~west~ أو ~south~. الأسماء وحد الحجم وأخطاء lexer زي 20-3. كل جملة لازم semicolon؛ newline مسافة؛ البرنامج الفاضي نجاح.

الأسماء حتى keywords مسموحة كمتغيرات في مكان الاسم. ~let~ تضيف أو تستبدل بعد نجاح تقييم اليمين؛ كل استدعاء يبدأ جديد.

حلّل البرنامج كله لـ AST قبل أي حركة. ~+~ و~-~ و~*~ و~/~ left associative بالأولوية اللي فوق مع unary signs وأقواس.

كل قيمة وناتج وسيط بين -1000000 و 1000000؛ القسمة بتحذف الكسر ناحية الصفر وصفر divisor خطأ.

عدد ~thrust~ من 0 لـ 32 قبل أي ~Step~. ~Step~ بـ ~false~ توقف فوراً وتسيب الخلايا السابقة اللي نجحت.

أول خطأ يوقف؛ الحركات الناجحة السابقة في runtime تفضل.

أخطاء اللغة lexical وsyntax وruntime كلها ~*CommandError~ مع ~Message~ مش فاضية و~Offset~ بايت من صفر: token الحالية للـ syntax (~EOF~ بطول المصدر)، الاسم للمتغير الغايب، العملية للحساب، وكلمة ~thrust~ لعدد باطل أو حركة مسدودة.

~nil~ drive خطأ قبل parsing. مفيش ~Report~؛ المختبر يحرك الدرون الحقيقي بـadapter.

النجمة التالتة شجر unary وأقواس عميقة وأخطاء المدى الوسيط.
`, Hints: []string{"خلي قواعد الـgrammar تفصل مستويات الأولوية.", "ابني AST قبل التقييم، وكل node تحقق FlightExpr.", "خلي primary تنادي نفسها للعلامة الأحادية، وخزّن الـoperator في AST."},
		En: &LevelText{LessonRevision: 1, Steps: ch20EvaluatorEnglish, Title: "Flight Evaluator", Goal: "Interpret flight programs with precedence and variables and move the real drone", Lesson: `# From Text to Real Movement

## Before you write
- How will you prevent a syntax error from executing earlier movement?
- Does an arithmetic error point at the operator or the statement?

## Grammar encodes precedence
Continue with your lexer in ~main.go~. Add ~flight.go~ for the interpreter. A recursive-descent ~parser~ has one function per precedence level, and consumes ~tokens~ using a cursor.

An AST (abstract syntax tree) stores structure for a later evaluation pass.
~~~text
program   = { statement }
statement = "let" name "=" expression ";" | "thrust" direction expression ";"
expression = term { ("+" | "-") term }
term       = primary { ("*" | "/") primary }
primary    = number | name | "(" expression ")" | ("+" | "-") primary
~~~
Multiplication is deeper than addition, so ~2+3*4~ groups as ~2+(3*4)~. Loops make operators at one level left associative: ~9-3-1~ means ~(9-3)-1.

Unary signs call the ~primary~ function directly, so it is also directly recursive. Use a free function ~primary(p *parser)~ and call ~primary(p)~ again after consuming a unary sign; the requirement scanner recognizes direct function recursion. Parentheses call ~expression~ again, reusing recursion from 14-5.
~~~go
func primary(p *parser)(FlightExpr,error) {
 if p.peek().Kind=="(" {
  p.next++
  expr,err:=p.expression()
  if err!=nil {return nil,err}
  if _,err=p.take(")");err!=nil {return nil,err}
  return expr,nil
 }
 if p.peek().Kind=="+" {
  p.next++
  return primary(p)
 }
 return p.atom()
}
~~~
This sketch assumes a ~parser~ with ~tokens~, ~next~, ~peek~, ~take~, ~expression~ and ~atom~ methods; define your own internal names.

~take~ must reject the current token without advancing past ~EOF~. Missing delimiters point at the current token, including ~eof~. Adjacent primaries are errors; implicit multiplication is not allowed.

## An interface for expression trees
~FlightExpr~ holds any AST node with ~Eval(map[string]int64)(int64,error)~. Use interfaces from 9-4 and methods from chapter 7: ~literal~, ~variable~ and ~binary~ nodes can implement this contract independently.

A map is the environment for one execution.
~~~go
type literal struct {n int64}
func(n literal)Eval(env map[string]int64)(int64,error){return n.n,nil}
var expr FlightExpr=literal{7}
env:=map[string]int64{"fuel":3}
n,err:=expr.Eval(env)
fmt.Println(n,err) // 7 <nil>
~~~
Store each operator's token in its node. Binary evaluation recursively evaluates left then right, then applies the operator.

A variable uses map lookup with its presence bool; missing is an error at that name. Unary minus can be a subtraction node with a zero left operand at the sign's position.
~~~go
n,ok:=env[name]
if !ok {return 0,&CommandError{Offset:offset,Message:"undefined variable"}}
~~~
All values use ~int64~ from 20-3. Require each intermediate result in [-1000000,1000000]. Two valid operands have a product at most 10^12, so computing before checking cannot overflow ~int64~.

Check division by zero before dividing; integer division truncates toward zero. Range and division errors point at the operator.

## Parse first, then drive
Build the entire statement list before evaluating anything, so syntax errors cannot move the drone.

Evaluate ~let~'s right side successfully before updating the environment.

Each call starts with a fresh environment. Runtime errors stop subsequent statements but retain completed movements.
~~~go
type baseDrive struct{}
func(baseDrive)Step(dir colony.Direction)bool{return colony.Move(dir)}
err:=ExecuteFlight("let cells=2; thrust north cells;",baseDrive{})
~~~
~FlightDrive~ is a consumer-owned interface with ~Step(colony.Direction) bool~. False means blocked.

Its adapter above really moves the colony drone; hidden integration checks position, while other tests use an injected boundary. ~colony.Direction~, ~Move~ and the four constants come from chapter 1. Map direction names to those constants.

A ~thrust~ first validates its full count, then calls ~Step~ once per cell; zero calls nothing.

If ~Step~ fails, stop immediately. Previously successful cells remain moved, rather than claiming rollback.

A ~nil~ interface returns an ordinary error; typed nil pointers are outside the domain.

## Common mistakes
~~~go
env[name]=n // Before verifying RHS evaluation succeeded
left=left+right // Applying addition immediately loses multiplication precedence
~~~

## Your mission
Keep ~LexFlight~ and its types unchanged in ~main.go~. Implement ~ExecuteFlight(source string,drive FlightDrive) error~ in ~flight.go~ with ~FlightDrive interface{Step(colony.Direction) bool}~ and an AST using ~FlightExpr interface{Eval(map[string]int64)(int64,error)}~.

Follow the grammar above; direction is ~east~, ~north~, ~west~, or ~south~. Names, size limit and lexical errors follow 20-3. Every statement needs a semicolon; newlines are whitespace; an empty program succeeds.

Names including keywords may be variables in name positions. ~let~ creates or replaces a variable after evaluating its RHS; each call starts fresh.

Parse the entire program into AST nodes before any movement. ~+~, ~-~, ~*~, ~/~ are left associative with the stated precedence; unary signs and parentheses are supported.

Values and all intermediates must stay in [-1000000,1000000]; division truncates toward zero and zero divisors fail.

~thrust~ count must be 0..32 before any ~Step~. ~Step~ returning ~false~ fails immediately after retaining prior successful cells. Stop on the first error; successful prior runtime movements stay.

All lexical, syntax and runtime language errors are ~*CommandError~ with nonempty ~Message~ and zero-based byte ~Offset~: bad/current token for syntax (~EOF~ at ~len(source)~), name for undefined variable, operator for arithmetic, ~thrust~ keyword for bad count or blocked movement.

A ~nil~ drive returns an error before parsing.

No ~Report~; the lab drives the actual drone through an adapter.

The third star checks deep unary/parenthesized trees and intermediate range errors.
`, Hints: []string{"Let grammar rules separate precedence levels.", "Build the AST before evaluation; each node implements FlightExpr.", "Have primary call itself for unary signs and keep operators in the AST."}, Files: []LabFile{{Name: "main.go", Starter: flightEvaluatorStarter}, {Name: "flight.go", Starter: flightEvaluatorFlightStarter}}}, Kind: "lab", Setup: seededLabRoom, Files: []LabFile{{Name: "main.go", Starter: flightEvaluatorStarterAr}, {Name: "flight.go", Starter: flightEvaluatorFlightStarterAr}}, Lab: &Lab{ReqFile: "flight.go", Hidden: flightEvaluatorHidden, Bonus: true}, CarryFrom: "20-3-lexer", MaxTicks: 200})
}

const bulletinHidden = `package main

import (
	"fmt"
	"html"
	"io"
	"orbital/colony/levels"
	"strings"
	"testing"
)

func TestHiddenBulletinStatesAndEscaping(t *testing.T) {
	w := levels.Start("20-1-bulletin", io.Discard, 201)
	for _, tc := range []struct{ source, want string }{{"", ""}, {"\n", "<br>\n"}, {"a\nb", "a<br>\nb"}, {"[alert][/alert][code][/code]", "<strong></strong><code></code>"}, {"[other]<x>", "[other]&lt;x&gt;"}} {
		got, err := RenderBulletin(tc.source)
		if err != nil || got != tc.want {
			t.Fatalf("%q => %q %v", tc.source, got, err)
		}
	}
	for i := 0; i < 60; i++ {
		text := fmt.Sprintf("relay%d <script>\"&'</script> &lt;", w.Rand.Intn(100000))
		source := "status\n[alert]" + text + "[/alert][code]" + text + "[/code]"
		want := "status<br>\n<strong>" + html.EscapeString(text) + "</strong><code>" + html.EscapeString(text) + "</code>"
		if got, err := RenderBulletin(source); err != nil || got != want {
			t.Fatal("escape text exactly once", got, err)
		}
	}
	for _, bad := range []string{"[alert]x", "prefix[code]", "[/alert]", "[alert]x[/code]", "[code][alert]x[/alert][/code]", "[alert][alert][/alert][/alert]"} {
		if got, err := RenderBulletin(bad); err == nil || got != "" {
			t.Fatal("bad syntax returned partial HTML", bad, got, err)
		}
	}
}
func TestHiddenBonusMultilineBulletin(t *testing.T) {
	text := strings.Repeat("قمر🚀<&\n", 20000)
	got, err := RenderBulletin("[code]" + text + "[/code]\n[alert]" + text + "[/alert]")
	want := "<code>" + html.EscapeString(text) + "</code><br>\n<strong>" + html.EscapeString(text) + "</strong>"
	if err != nil || got != want {
		t.Fatal("multiline large Unicode block", err)
	}
	if got, err := RenderBulletin("[code]" + text); err == nil || got != "" {
		t.Fatal("large unclosed block")
	}
}
`
const bulletinStarter = `package main

import "errors"

func RenderBulletin(source string) (string, error) {
	return "", errors.New("bulletin unfinished")
}
func main() {}
`
const bulletinStarterAr = `package main

import "errors"

func RenderBulletin(source string) (string, error) {
	return "", errors.New("النشرة لسه مخلصتش")
}
func main() {}
`
const crewCardHidden = `package main

import (
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"slices"
	"testing"
)

func TestHiddenCrewCardIndependentRules(t *testing.T) {
	for _, tc := range []struct {
		s  string
		ok bool
	}{{"echo@relay.colony", true}, {"قمر@relay.colony", true}, {"é@a.b", true}, {"", false}, {"echo", false}, {"a@@b.c", false}, {"a\u00a0b@x.y", false}, {"a@b..c", false}, {"a@قمر.c", false}, {" a@b.c", false}, {"a@b", false}} {
		err := CheckAddress(tc.s)
		if (err == nil) != tc.ok || (!tc.ok && !errors.Is(err, ErrAddress)) {
			t.Fatal("address", tc.s, err)
		}
	}
	for _, tc := range []struct {
		s  string
		ok bool
	}{{"+1234567", true}, {"+123456789012", true}, {"+123456", false}, {"+1234567890123", false}, {"1234567", false}, {"+١٢٣٤٥٦٧", false}, {"+123 4567", false}, {"", false}} {
		err := CheckRadio(tc.s)
		if (err == nil) != tc.ok || (!tc.ok && !errors.Is(err, ErrRadio)) {
			t.Fatal("radio", tc.s, err)
		}
	}
	for _, tc := range []struct {
		s  string
		ok bool
	}{{"2024/02/29", true}, {"2000/02/29", true}, {"1900/02/29", false}, {"2025/02/30", false}, {"2025/13/01", false}, {"2025/1/01", false}, {"2025-01-01", false}, {"", false}} {
		err := CheckLaunchDate(tc.s)
		if (err == nil) != tc.ok || (!tc.ok && !errors.Is(err, ErrLaunchDate)) {
			t.Fatal("date", tc.s, err)
		}
	}
}
func TestHiddenCrewCardAggregateOrder(t *testing.T) {
	w := levels.Start("20-2-validators", io.Discard, 202)
	for i := 0; i < 80; i++ {
		mask := w.Rand.Intn(8)
		card := CrewCard{fmt.Sprintf("pilot%d@relay.colony", w.Rand.Intn(100000)), "+12345678", "2024/02/29"}
		var want []string
		if mask&1 != 0 {
			card.Address = "bad"
			want = append(want, "address")
		}
		if mask&2 != 0 {
			card.Radio = "bad"
			want = append(want, "radio")
		}
		if mask&4 != 0 {
			card.LaunchDate = "2025/02/29"
			want = append(want, "launchDate")
		}
		err := ValidateCrewCard(card)
		if (err == nil) != (mask == 0) || !slices.Equal(InvalidFields(fmt.Errorf("enrollment: %w", err)), want) {
			t.Fatal("all failures in fixed order", mask, err)
		}
	}
	if len(InvalidFields(nil)) != 0 || len(InvalidFields(errors.New("unrelated"))) != 0 {
		t.Fatal("unknown error classified")
	}
}
`
const crewCardStarter = `package main

import "errors"

var ErrAddress = errors.New("address")
var ErrRadio = errors.New("radio")
var ErrLaunchDate = errors.New("launch date")

type CrewCard struct{ Address, Radio, LaunchDate string }

func CheckAddress(s string) error          { return ErrAddress }
func CheckRadio(s string) error            { return ErrRadio }
func CheckLaunchDate(s string) error       { return ErrLaunchDate }
func ValidateCrewCard(card CrewCard) error { return ErrAddress }
func InvalidFields(err error) []string     { return nil }
func main()                                {}
`
const crewCardStarterAr = `package main

import "errors"

var ErrAddress = errors.New("العنوان")
var ErrRadio = errors.New("الراديو")
var ErrLaunchDate = errors.New("تاريخ الإطلاق")

type CrewCard struct{ Address, Radio, LaunchDate string }

func CheckAddress(s string) error          { return ErrAddress }
func CheckRadio(s string) error            { return ErrRadio }
func CheckLaunchDate(s string) error       { return ErrLaunchDate }
func ValidateCrewCard(card CrewCard) error { return ErrAddress }
func InvalidFields(err error) []string     { return nil }
func main()                                {}
`
const flightLexerHidden = `package main

import (
	"errors"
	"fmt"
	"io"
	"orbital/colony/levels"
	"slices"
	"strings"
	"testing"
)

func TestHiddenFlightTokenKindsAndOffsets(t *testing.T) {
	source := "let n=12;\r\nthrust east (n+2)*3/-1;"
	want := []CommandToken{{"name", "let", 0}, {"name", "n", 4}, {"=", "=", 5}, {"number", "12", 6}, {";", ";", 8}, {"name", "thrust", 11}, {"name", "east", 18}, {"(", "(", 23}, {"name", "n", 24}, {"+", "+", 25}, {"number", "2", 26}, {")", ")", 27}, {"*", "*", 28}, {"number", "3", 29}, {"/", "/", 30}, {"-", "-", 31}, {"number", "1", 32}, {";", ";", 33}, {"eof", "", 34}}
	got, err := LexFlight(source)
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("tokens=%v want %v err=%v", got, want, err)
	}
	w := levels.Start("20-3-lexer", io.Discard, 203)
	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("bay_%d", w.Rand.Intn(1000))
		number := fmt.Sprint(w.Rand.Intn(1001))
		source := name + " = " + number + ";"
		got, err := LexFlight(source)
		want := []CommandToken{{"name", name, 0}, {"=", "=", len(name) + 1}, {"number", number, len(name) + 3}, {";", ";", len(source) - 1}, {"eof", "", len(source)}}
		if err != nil || !slices.Equal(got, want) {
			t.Fatal("seeded tokens", got, err)
		}
	}
	for _, source := range []string{"", " \t\r\n"} {
		got, err := LexFlight(source)
		if err != nil || !slices.Equal(got, []CommandToken{{"eof", "", len(source)}}) {
			t.Fatal("empty", got, err)
		}
	}
	for _, tc := range []struct {
		s      string
		offset int
	}{{"abc @", 4}, {"x ق", 2}, {"x \xff", 2}, {"1001", 0}, {"9999999999999999999999999", 0}, {"A", 0}, {"_x", 0}, {strings.Repeat(" ", 4097), 4096}} {
		got, err := LexFlight(tc.s)
		var diagnostic *CommandError
		if got != nil || !errors.As(err, &diagnostic) || diagnostic.Offset != tc.offset || diagnostic.Message == "" {
			t.Fatal("diagnostic", tc.offset, got, err)
		}
	}
	if got, err := LexFlight(strings.Repeat(" ", 4096)); err != nil || len(got) != 1 {
		t.Fatal("exact length boundary")
	}
}
`
const flightLexerStarter = `package main

import "fmt"

type CommandToken struct {
	Kind, Text string
	Offset     int
}
type CommandError struct {
	Offset  int
	Message string
}

func (e *CommandError) Error() string { return fmt.Sprintf("byte %d: %s", e.Offset, e.Message) }
func LexFlight(source string) ([]CommandToken, error) {
	return nil, &CommandError{0, "lexer unfinished"}
}
func main() {}
`
const flightLexerStarterAr = `package main

import "fmt"

type CommandToken struct {
	Kind, Text string
	Offset     int
}
type CommandError struct {
	Offset  int
	Message string
}

func (e *CommandError) Error() string { return fmt.Sprintf("بايت %d: %s", e.Offset, e.Message) }
func LexFlight(source string) ([]CommandToken, error) {
	return nil, &CommandError{0, "مقطع الأوامر لسه مخلصش"}
}
func main() {}
`
const flightEvaluatorFlightStarter = `package main

import (
	"errors"
	"orbital/colony"
)

type FlightDrive interface{ Step(colony.Direction) bool }
type FlightExpr interface {
	Eval(map[string]int64) (int64, error)
}

func ExecuteFlight(source string, drive FlightDrive) error {
	return errors.New("evaluator unfinished")
}
`
const flightEvaluatorFlightStarterAr = `package main

import (
	"errors"
	"orbital/colony"
)

type FlightDrive interface{ Step(colony.Direction) bool }
type FlightExpr interface {
	Eval(map[string]int64) (int64, error)
}

func ExecuteFlight(source string, drive FlightDrive) error {
	return errors.New("المنفذ لسه مخلصش")
}
`
const flightEvaluatorHidden = `package main

import (
	"errors"
	"fmt"
	"io"
	"orbital/colony"
	"orbital/colony/levels"
	"slices"
	"strings"
	"testing"
)

type hiddenFlightDrive struct {
	directions []colony.Direction
	limit      int
}

func (d *hiddenFlightDrive) Step(dir colony.Direction) bool {
	if len(d.directions) == d.limit {
		return false
	}
	d.directions = append(d.directions, dir)
	return true
}
func TestHiddenFlightPrecedenceAndVariables(t *testing.T) {
	w := levels.Start("20-4-evaluator", io.Discard, 204)
	for i := 0; i < 70; i++ {
		a, b, c := w.Rand.Intn(4), w.Rand.Intn(4), w.Rand.Intn(4)
		source := fmt.Sprintf("let fuel=%d; let fuel=fuel+%d*%d; thrust east fuel; thrust north (%d+%d)*%d;", a, b, c, a, b, c)
		drive := &hiddenFlightDrive{limit: 100}
		err := ExecuteFlight(source, drive)
		want := append(slices.Repeat([]colony.Direction{colony.East}, a+b*c), slices.Repeat([]colony.Direction{colony.North}, (a+b)*c)...)
		if err != nil || !slices.Equal(drive.directions, want) {
			t.Fatal("precedence/parentheses/reassignment", source, drive, err)
		}
	}
	for _, tc := range []struct {
		source string
		steps  int
	}{{"", 0}, {" \n", 0}, {"thrust east 10-3-2;", 5}, {"thrust east 20/4/5;", 1}, {"let x=4; let y=x; thrust east y;", 4}, {"thrust east 0;", 0}, {"thrust east 32;", 32}} {
		drive := &hiddenFlightDrive{limit: 100}
		if err := ExecuteFlight(tc.source, drive); err != nil || len(drive.directions) != tc.steps {
			t.Fatal("boundaries/associativity", tc, drive, err)
		}
	}
}
func TestHiddenFlightDiagnosticsAndEffects(t *testing.T) {
	for _, tc := range []struct {
		source        string
		offset, steps int
	}{
		{"thrust east missing;", 12, 0}, {"thrust east 1/0;", 13, 0}, {"thrust east (1+2;", 16, 0},
		{"thrust east 2(3);", 13, 0}, {"thrust up 1;", 7, 0}, {"fly east 1;", 0, 0},
		{"thrust east -1;", 0, 0}, {"thrust east 33;", 0, 0}, {"thrust east 1", 13, 0},
		{"thrust east 1; let x=;", 21, 0}, {"thrust east 1; thrust north 1/0;", 29, 1},
	} {
		drive := &hiddenFlightDrive{limit: 100}
		err := ExecuteFlight(tc.source, drive)
		var diagnostic *CommandError
		if !errors.As(err, &diagnostic) || diagnostic.Offset != tc.offset || diagnostic.Message == "" || len(drive.directions) != tc.steps {
			t.Fatalf("%q: %v steps=%d want offset=%d steps=%d", tc.source, err, len(drive.directions), tc.offset, tc.steps)
		}
	}
	drive := &hiddenFlightDrive{limit: 1}
	var diagnostic *CommandError
	if err := ExecuteFlight("thrust east 3;", drive); !errors.As(err, &diagnostic) || diagnostic.Offset != 0 || len(drive.directions) != 1 {
		t.Fatal("blocked drive", err)
	}
	if err := ExecuteFlight("", nil); err == nil {
		t.Fatal("nil drive")
	}
	drive = &hiddenFlightDrive{limit: 100}
	_ = ExecuteFlight("let x=2;", drive)
	if err := ExecuteFlight("thrust east x;", drive); err == nil {
		t.Fatal("environment leaked across programs")
	}
}

type hiddenColonyDrive struct{}

func (hiddenColonyDrive) Step(dir colony.Direction) bool { return colony.Move(dir) }
func TestHiddenFlightMovesRealDrone(t *testing.T) {
	x, y := colony.Position()
	if err := ExecuteFlight("let cells=1+1; thrust east cells; thrust north 1;", hiddenColonyDrive{}); err != nil {
		t.Fatal(err)
	}
	nx, ny := colony.Position()
	if nx != x+2 || ny != y+1 {
		t.Fatalf("real drone at %d,%d from %d,%d", nx, ny, x, y)
	}
	if err := ExecuteFlight("thrust south 1; thrust west 2;", hiddenColonyDrive{}); err != nil {
		t.Fatal(err)
	}
}
func TestHiddenBonusFlightDeepTreesAndBounds(t *testing.T) {
	drive := &hiddenFlightDrive{limit: 100}
	source := "thrust east " + strings.Repeat("(", 40) + "---(-2)" + strings.Repeat(")", 40) + ";"
	if err := ExecuteFlight(source, drive); err != nil || len(drive.directions) != 2 {
		t.Fatal("nested unary AST", err)
	}
	for _, source := range []string{"let x=1000*1000+1;", "let x=(1000*1000)*1000;", "let x=-(1000*1000)-1;"} {
		var diagnostic *CommandError
		if err := ExecuteFlight(source, drive); !errors.As(err, &diagnostic) {
			t.Fatal("arithmetic bound", source, err)
		}
	}
}
`
const flightEvaluatorStarter = `package main

import "fmt"

type CommandToken struct {
	Kind, Text string
	Offset     int
}
type CommandError struct {
	Offset  int
	Message string
}

func (e *CommandError) Error() string { return fmt.Sprintf("byte %d: %s", e.Offset, e.Message) }
func LexFlight(source string) ([]CommandToken, error) {
	return nil, &CommandError{0, "lexer unfinished"}
}
func main() {}
`
const flightEvaluatorStarterAr = `package main

import "fmt"

type CommandToken struct {
	Kind, Text string
	Offset     int
}
type CommandError struct {
	Offset  int
	Message string
}

func (e *CommandError) Error() string { return fmt.Sprintf("بايت %d: %s", e.Offset, e.Message) }
func LexFlight(source string) ([]CommandToken, error) {
	return nil, &CommandError{0, "مقطع الأوامر لسه مخلصش"}
}
func main() {}
`
