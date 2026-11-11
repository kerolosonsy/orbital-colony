package levels

func init() {
	add(&Level{ID: "23-1-access-codes", Chapter: 23, Kind: "lab", Title: "أكواد الدخول", Goal: "طلّع أكواد للبوابات بفئات مضمونة وعشوائية محقونة",
		Requires: []string{"sel:rand.Read"}, Bridge: []int{10}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch231Starter}},
		Lesson: `# أكواد بوابات الخزنة

البوابات محتاجة أكواد مؤقتة، مش لعبة تخمين. اختيار مصدر العشوائية جزء من صحة الحل؛ متطبعش الكود في log.

## قبل ما تكتب
- هترفض طلب مستحيل قبل قراءة أي بايت إزاي؟
- إزاي تضمن ظهور كل فئة من غير انتظار الحظ؟

## عشوائية للأسرار
استورد ~crypto/rand~ باسم ~rand~. دي عشوائية مناسبة للأسرار من النظام. ~math/rand~ و ~math/rand/v2~ قابلين للتوقع ومش مناسبين هنا، حتى مع seed بالوقت؛ ممنوع استيرادهم في ملفات الحل كلها.
~~~go
buf:=make([]byte,4)
n,err:=rand.Read(buf)
fmt.Println(n,err) // متطبعش البايتات السرية
~~~
~rand.Read~ يملأ البايتات. لما ~source~ تكون ~nil~ في المهمة استخدمه؛ غير كده استخدم الـ reader المحقون فقط. التست يديك بايتات ثابتة أو خطأ، ومش بيغيّر المصدر العام للمكتبة.

~io.ReadFull~ يكمّل قراءات قصيرة لحد ملء الـ buffer؛ EOF بدري خطأ. افحص الخطأ قبل استخدام أي بايت.
~~~go
var b [1]byte
_,err:=io.ReadFull(source,b[:])
if err!=nil {return "",err}
~~~

## فئات مضمونة مش احتمالية
ابدأ بحرف من كل فئة مطلوبة، وكمّل الطول من اتحاد الفئات. امزج الناتج عشان موضع أول حرف ميبقاش دايماً نفس الفئة. ~categories~ أسماء بروتوكول ثابتة، مش نصوص مترجمة.
~~~go
alphabets:=map[string]string{"signal":"xyz","panel":"47"}
joined:=alphabets["signal"]+alphabets["panel"]
fmt.Println(joined) // xyz47، مثال مختلف عن فئات المهمة
~~~
امزج من آخر عنصر: اختار موضع من 0 لحد الموضع الحالي وبدّل الاتنين. ده Fisher-Yates، وكل اختيار لازم ييجي من نفس المصدر. متعدلش slice الفئات اللي المنادي بعتها.
~~~go
for i:=len(code)-1;i>0;i-- {
 j,err:=pick(i+1)
 if err!=nil {return "",err}
 code[i],code[j]=code[j],code[i]
}
~~~
~pick~ هنا دالة عندك بترجع عدد في ~[0,n)~ وخطأ، مش API جديدة في مكتبة.

## نجمة إضافية: اختيار من غير انحياز
~b % n~ يفضّل قيم لما 256 مش بتتقسم على ~n~. احسب ~limit:=256-256%n~ كـ int؛ ارمي البايتات الأكبر من أو تساوي الحد واقرأ تاني. المقبول يتقسم بالتساوي على الاختيارات. مع 10 اختيارات ارفض 250 لحد 255؛ 249 تدي 9. ده يزيل انحياز الاختيار الواحد، مش ادعاء توزيع موحّد على كل الأكواد الممكنة ولا رقم entropy.
~~~go
limit:=256-256%n
if int(b)<limit {return int(b)%n,nil}
~~~

## أخطاء شائعة
~~~go
code[i]=alphabet[b%byte(len(alphabet))] // انحياز، وكمان تجاهل خطأ المصدر
// اختيار كل الحروف من الاتحاد بس ممكن ينسى فئة كاملة
~~~

## المهمة
نفّذ ~IssueAccess(source io.Reader, length int, categories []string) (string,error)~. الطول من 1 لـ 128 شامل. الفئات المتاحة بالظبط: ~lower~ = ~abcdefghijklmnopqrstuvwxyz~، ~upper~ = ~ABCDEFGHIJKLMNOPQRSTUVWXYZ~، ~digits~ = ~0123456789~، ~marks~ = ~!@#$~. ارفض فئات فاضية أو مجهولة أو مكررة، وطول أقل من عددها، قبل أي قراءة.

كل كود ناجح طوله المطلوب بالبايت (كل الحروف ASCII)، مفيهوش حروف برا الفئات المختارة، وفيه حرف على الأقل من كل فئة ببناء صريح، وبعده مزج. ~nil~ للمصدر يعني ~crypto/rand.Read~؛ غير ~nil~ يعني كل العشوائية منه مع دعم القراءات القصيرة. ممنوع ~math/rand~ بنسختيه، ومفيش ملفات أو وقت أو شبكة أو Report. في أي فشل ارجع string فاضية وخطأ يحافظ على سبب فشل القراية؛ متغيرش المدخلات. ~main~ ممكن تفضل فاضية.

النجمة التالتة: كل اختيار حرف أو موضع مزج يستخدم rejection sampling بالحد المذكور. في حالة طول 1 وفئة ~digits~ فقط، المصدر 250 ثم 255 ثم 9 لازم يطلع ~9~، ومصدر 250 ثم EOF لازم يفشل من غير كود جزئي.
`,
		Hints: []string{"تحقق من الطول والفئات كلها قبل اختيار أول حرف.", "اختار حرف من كل فئة، كمّل من الاتحاد، وبعدين امزج من الآخر.", "اقرأ بايت واحد بـ ~io.ReadFull~؛ للـ bonus ارفض اللي فوق الحد قبل modulo."},
		En: &LevelText{Title: "Access Codes", Goal: "Issue gate codes with guaranteed categories and injected randomness", Files: []LabFile{{Name: "main.go", Starter: ch231Starter}},
			Lesson: `# Vault Gate Codes

The gates need temporary secrets, rather than a guessing game. Choosing a random ~source~ is part of correctness. Never log the generated code.

## Before you write
- How will you reject an impossible request before reading any bytes?
- How will every requested category appear without relying on chance?

## Randomness for secrets
Import ~crypto/rand~ as ~rand~. It provides system-backed randomness suitable for secrets. ~math/rand~ and ~math/rand/v2~ are predictable and unsuitable here, even with a time seed. Neither may be imported anywhere in the solution.
~~~go
buf:=make([]byte,4)
n,err:=rand.Read(buf)
fmt.Println(n,err) // Do not print secret bytes
~~~
~rand.Read~ fills bytes. Use it when this mission's ~source~ is ~nil~; otherwise use only the injected reader. Tests supply deterministic bytes or errors without replacing the library's global source.

~io.ReadFull~ completes short reads until the buffer is filled. Early EOF is an error. Check errors before using bytes.
~~~go
var b [1]byte
_,err:=io.ReadFull(source,b[:])
if err!=nil {return "",err}
~~~

## Structural category coverage
Pick one character from each requested category, fill remaining positions from their union, then shuffle so leading positions do not always reveal categories. Category names are fixed protocol strings.
~~~go
alphabets:=map[string]string{"signal":"xyz","panel":"47"}
joined:=alphabets["signal"]+alphabets["panel"]
fmt.Println(joined) // xyz47; different categories from the mission
~~~
Shuffle backwards: choose a position between zero and the current position inclusive, then swap. This is Fisher-Yates. All choices use the same source. Do not mutate the caller's categories slice.
~~~go
for i:=len(code)-1;i>0;i-- {
 j,err:=pick(i+1)
 if err!=nil {return "",err}
 code[i],code[j]=code[j],code[i]
}
~~~
~pick~ denotes your own function returning a number in ~[0,n)~ and an error, rather than a library API.

## Third star: unbiased individual choices
~b % n~ favours values when 256 is not divisible by ~n~. Compute ~limit:=256-256%n~ as an int. Reject bytes greater than or equal to that limit and read again. Accepted bytes divide evenly among choices. With 10 choices reject 250 through 255; 249 maps to 9. This removes individual selection bias; it does not claim a uniform distribution over all valid codes or a numerical entropy value.
~~~go
limit:=256-256%n
if int(b)<limit {return int(b)%n,nil}
~~~

## Common mistakes
~~~go
code[i]=alphabet[b%byte(len(alphabet))] // Biased and ignores source errors
// Drawing only from the union can omit a required category entirely
~~~

## Your mission
Implement ~IssueAccess(source io.Reader, length int, categories []string) (string,error)~. Length is 1 through 128 inclusive. Exactly four categories exist: ~lower~ = ~abcdefghijklmnopqrstuvwxyz~, ~upper~ = ~ABCDEFGHIJKLMNOPQRSTUVWXYZ~, ~digits~ = ~0123456789~, ~marks~ = ~!@#$~. Reject no categories, unknown or duplicate names, and length smaller than the category count before reading.

Success returns exactly length bytes (all characters are ASCII), only from enabled categories, with each category represented by construction, followed by shuffling. A ~nil~ ~source~ means ~crypto/rand.Read~; otherwise all randomness comes from the injected reader, including short-read support. Both math/rand versions are forbidden. No filesystem, time, network, or Report. Every failure returns an empty string and error preserving reader causes; do not mutate inputs. ~main~ may remain empty.

Third star: use the rejection rule above for every character and shuffle-position choice. For length 1 with only ~digits~, bytes 250, 255, 9 must produce ~9~; 250 followed by EOF must fail without a partial code.
`,
			Hints: []string{"Validate length and all category names before drawing.", "Pick one per category, fill from the union, then shuffle backwards.", "Read one byte with ~io.ReadFull~; for the bonus reject bytes beyond the cutoff before modulo."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch231Hidden},
	})
}

const ch231Starter = `package main

import "io"

func IssueAccess(source io.Reader, length int, categories []string) (string, error) {
	return "", io.EOF
}
func main() {}
`

const ch231Hidden = `package main

import (
	"bytes"
	"errors"
	"go/parser"
	"go/token"
	"io"
	"orbital/colony/levels"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"testing/iotest"
)

func TestHiddenAccessNoMathRandom(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range files {
		f, err := parser.ParseFile(token.NewFileSet(), name, nil, parser.ImportsOnly)
		if err != nil {
			t.Fatal(err)
		}
		for _, imp := range f.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				t.Fatal(err)
			}
			if path == "math/rand" || path == "math/rand/v2" {
				t.Fatal("predictable random source", name)
			}
		}
	}
}
func TestHiddenAccessSeededCategories(t *testing.T) {
	w := levels.Start("23-1-access-codes", io.Discard, 231)
	names := []string{"lower", "upper", "digits", "marks"}
	alphabets := []string{"abcdefghijklmnopqrstuvwxyz", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "0123456789", "!@#$"}
	for mask := 1; mask < 16; mask++ {
		var cats []string
		allowed := ""
		for i, name := range names {
			if mask&(1<<i) != 0 {
				cats = append(cats, name)
				allowed += alphabets[i]
			}
		}
		before := slices.Clone(cats)
		for _, length := range []int{len(cats), 16, 128} {
			entropy := make([]byte, 8192)
			w.Rand.Read(entropy)
			got, err := IssueAccess(iotest.OneByteReader(bytes.NewReader(entropy)), length, cats)
			if err != nil || len(got) != length {
				t.Fatal("length", mask, length, len(got), err)
			}
			for _, b := range []byte(got) {
				if !strings.ContainsRune(allowed, rune(b)) {
					t.Fatal("outside alphabet")
				}
			}
			for i, name := range names {
				if slices.Contains(cats, name) && !strings.ContainsAny(got, alphabets[i]) {
					t.Fatal("missing category", name)
				}
			}
			again, err := IssueAccess(bytes.NewReader(entropy), length, cats)
			if err != nil || got != again || !slices.Equal(cats, before) {
				t.Fatal("injection/replay/input mutation", err)
			}
		}
	}
}
func TestHiddenAccessShufflesCategoryPositions(t *testing.T) {
	w := levels.Start("23-1-access-codes", io.Discard, 2317)
	orders := map[string]bool{}
	for i := 0; i < 32; i++ {
		entropy := make([]byte, 128)
		w.Rand.Read(entropy)
		code, err := IssueAccess(bytes.NewReader(entropy), 4, []string{"lower", "upper", "digits", "marks"})
		if err != nil || len(code) != 4 {
			t.Fatal("shuffle case", err)
		}
		var order strings.Builder
		for _, b := range code {
			switch {
			case b >= 'a' && b <= 'z':
				order.WriteByte('l')
			case b >= 'A' && b <= 'Z':
				order.WriteByte('u')
			case b >= '0' && b <= '9':
				order.WriteByte('d')
			default:
				order.WriteByte('m')
			}
		}
		orders[order.String()] = true
	}
	if len(orders) < 2 {
		t.Fatal("category positions never shuffled")
	}
}
func TestHiddenAccessValidationAndSourceFailure(t *testing.T) {
	wire := errors.New("entropy offline")
	for _, tc := range []struct {
		n    int
		cats []string
	}{
		{0, []string{"digits"}}, {-1, []string{"digits"}}, {129, []string{"lower"}},
		{3, nil}, {2, []string{"lower", "upper", "digits"}}, {4, []string{"digits", "digits"}},
		{4, []string{"other"}}, {4, []string{""}},
	} {
		r := &hiddenAccessForbidden{t: t}
		if got, err := IssueAccess(r, tc.n, tc.cats); err == nil || got != "" {
			t.Fatal("invalid request", tc, got, err)
		}
	}
	for _, r := range []io.Reader{iotest.ErrReader(wire), io.MultiReader(bytes.NewReader([]byte{1}), iotest.ErrReader(wire))} {
		if got, err := IssueAccess(r, 8, []string{"digits"}); !errors.Is(err, wire) || got != "" {
			t.Fatal("partial secret/lost cause", got, err)
		}
	}
	if got, err := IssueAccess(bytes.NewReader(nil), 1, []string{"marks"}); err == nil || got != "" {
		t.Fatal("EOF", got, err)
	}
}

type hiddenAccessForbidden struct{ t *testing.T }

func (r *hiddenAccessForbidden) Read([]byte) (int, error) {
	r.t.Fatal("read before validation")
	return 0, io.EOF
}
func TestHiddenBonusAccessRejectBiasedTail(t *testing.T) {
	for _, tc := range []struct {
		input []byte
		want  string
	}{{[]byte{250, 255, 9}, "9"}, {[]byte{249}, "9"}, {[]byte{250}, ""}} {
		got, err := IssueAccess(bytes.NewReader(tc.input), 1, []string{"digits"})
		if got != tc.want || (err != nil) != (tc.want == "") {
			t.Fatal("rejection sampling", got, err)
		}
	}
}
`

func init() {
	add(&Level{ID: "23-2-sealed-cargo", Chapter: 23, Kind: "lab", Title: "شحنة مختومة", Goal: "شفّر حمولة واربطها ببيان الشحنة مع كشف أي تلاعب",
		Requires: []string{"sel:cipher.NewGCM"}, Bridge: []int{30}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch232Starter}},
		Lesson: `# ختم شحنة الخزنة

البصمة بتكشف اختلاف المحتوى، لكنها مش بتخفيه ومش بتثبت مين كتبه. شحنة الخزنة محتاجة سرية وكشف تلاعب بمفتاح يديه المنادي؛ دي مش منظومة إدارة مفاتيح ولا تحويل password لمفتاح.

## قبل ما تكتب
- هتمنع إعادة استخدام nonce مع نفس المفتاح إزاي؟
- هتمنع ظهور بايت واحد من الحمولة قبل التحقق من الختم إزاي؟

## AES ثم GCM
استورد ~crypto/aes~ و ~crypto/cipher~. ~aes.NewCipher~ تقبل مفتاح خام طوله 16 أو 24 أو 32 بايت؛ غير كده خطأ. ~cipher.NewGCM~ تبني تشفير موثّق فوقه: AEAD بتخفي النص وتضيف tag يكشف التلاعب. متكتبش خوارزمية تشفير بنفسك.
~~~go
block,err:=aes.NewCipher(key)
if err!=nil {return nil,err}
aead,err:=cipher.NewGCM(block)
if err!=nil {return nil,err}
fmt.Println(aead.NonceSize(),aead.Overhead()) // 12 16
~~~
~cipher.AEAD~ هي واجهة الختم والفتح اللي بترجعها ~NewGCM~.
~~~go
var box cipher.AEAD = aead
fmt.Println(box.NonceSize()) // 12
~~~

~NonceSize~ طول nonce و ~Overhead~ طول tag. nonce مش سر، لكنه لازم جديد لكل تشفير تحت نفس المفتاح. إعادة استخدامه خطر حتى لو الحمولة مختلفة. التست يحقن سيل بايتات مختلفة ويتأكد إن كل استدعاء يستهلك nonce جديدة؛ مش بيستدل على الأمان من صدفة اختلاف خرجتين عشوائيتين.

## المصدر المحقون والبيانات المصاحبة
زي المهمة السابقة: ~io.ReadFull~ تقرا nonce كاملة، وفشل المصدر يوقف من غير ~packet~ جزئية. ~nil~ هنا معناها ~rand.Reader~ من ~crypto/rand~، reader آمن جاهز، مش ~math/rand~.
~~~go
if source==nil {source=rand.Reader}
nonce:=make([]byte,aead.NonceSize())
if _,err:=io.ReadFull(source,nonce);err!=nil {return nil,err}
~~~

~Seal(dst,nonce,plaintext,aad)~ ترجع ciphertext وبعدها tag. ~aad~ بيانات موثّقة بس مش مشفرة. اربط بيان الشحنة بالختم؛ بيان مختلف لازم يفشل حتى لو المفتاح والـ ~packet~ صح. مثال هوية جهاز مختلف عن عقد المهمة:
~~~go
sealed:=aead.Seal(nil,nonce,[]byte("sensor=7"),[]byte("probe-A"))
plain,err:=aead.Open(nil,nonce,sealed,[]byte("probe-A"))
~~~
~Open~ تتحقق من tag قبل إرجاع النص. غلط المفتاح، أو nonce، أو ciphertext، أو البيان، كلها فشل تحقق. ارجع خطأ عام واحد ~authentication failed~ لكل فشل ~Open~؛ متقولش للمهاجم إيه اللي غلط.

## تنسيق الباكت وحدود الذاكرة
البروتوكول هنا 4 بايت ثابتة ~OCV1~، بعدها nonce من 12 بايت، وبعدها ciphertext مع tag من 16 بايت. استخدم ~append~ مع slice جديدة عشان متعدلش ~key~ أو ~manifest~ أو ~cargo~ أو packet. AAD هي البادئة ~OCV1~ متبوعة ببايتات ~manifest~ بالظبط.
~~~go
aad:=append([]byte("SENSOR1"),label...)
~~~
المثال بيبني prefix مستقل. افحص طول ~packet~ قبل أي slicing؛ الحمولة الفاضية صالحة وبتطلع tag برضه. الحد 65536 بايت حمولة، والحد الأقصى للباكت 65568 بايت.

## أخطاء شائعة
~~~go
nonce:=make([]byte,12) // غلط لو استخدمتها من غير ملء عشوائي كل مرة
plain,_:=aead.Open(nil,nonce,sealed,nil) // تجاهل الخطأ يضيّع التحقق والبيان
~~~

## المهمة
نفّذ ~SealCargo(key,manifest,cargo []byte,source io.Reader) ([]byte,error)~ و ~OpenCargo(key,manifest,packet []byte) ([]byte,error)~. ادعم مفاتيح AES الثلاثة، وحمولة من صفر لـ 65536 بايت شامل. قبل قراءة المصدر ارفض المفتاح والطول الغلط. استخدم AES-GCM وpacket وAAD بالتنسيق المذكور. ~source~ غير ~nil~ مصدر كل بايتات nonce؛ ~nil~ مصدره ~crypto/rand.Reader~. كل تشفير يقرأ nonce جديدة كاملة، ومفيش cache أو توليد بمصدر تاني.

~OpenCargo~ ترفض المفتاح الغلط طولاً، والبادئة الغلط، والباكت الأقصر من 32 أو الأكبر من 65568 بايت، من غير panic. الباقي يمر على ~Open~ مع البيان. فشل التحقق يرجّع نص خطأ ~authentication failed~ من غير تفاصيل. كل فشل يرجّع ~nil~ للخرج؛ فشل المصدر يحتفظ بالسبب. النجاح لا يغير المدخلات، حتى لو slices ليها مساحة زيادة. لا ملفات أو شبكة أو Report ولا output سرّي.

النجمة التالتة: نفّذ ~OpenCargoTo(key,manifest,packet []byte,out io.Writer) error~. writer غير ~nil~؛ مفيش أي كتابة قبل نجاح الفتح بالكامل. بعدها اكتب النص مرة واحدة، واحتفظ بسبب خطأ الكتابة؛ لو الكتابة قصيرة من غير خطأ ارجع ~io.ErrShortWrite~. الجزء المكتوب قبل فشل النقل ممكن يفضل. ~io.ErrShortWrite~ خطأ قياسي لما writer تقبل بايتات أقل من المطلوب.
~~~go
n,err:=out.Write(plain)
if err==nil && n!=len(plain) {err=io.ErrShortWrite}
~~~
`,
		Hints: []string{"ابنِ AES وبعدها GCM؛ افحص الحدود قبل قراءة nonce.", "اقرأ ~NonceSize~ بايت بـ ~io.ReadFull~ لكل تشفير، وحطهم بعد ~OCV1~.", "ابنِ AAD في slice جديدة؛ افتح بـ ~nil~ destination ومترجعش نص لو التحقق فشل."},
		En: &LevelText{Title: "Sealed Cargo", Goal: "Encrypt cargo and authenticate its manifest and sealed bytes", Files: []LabFile{{Name: "main.go", Starter: ch232Starter}},
			Lesson: `# Seal the Vault Cargo

A digest detects different bytes but does not conceal them or prove who wrote them. Vault ~cargo~ needs confidentiality and tamper detection with a caller-supplied raw key. Key management and password derivation are outside this mission.

## Before you write
- How will you avoid nonce reuse under the same ~key~?
- How will you prevent publishing even one plaintext byte before authentication?

## AES followed by GCM
Import ~crypto/aes~ and ~crypto/cipher~. ~aes.NewCipher~ accepts raw keys of 16, 24, or 32 bytes; other lengths fail. ~cipher.NewGCM~ wraps the block cipher in authenticated encryption: AEAD hides plaintext and adds a tag detecting tampering. Never implement your own cipher.
~~~go
block,err:=aes.NewCipher(key)
if err!=nil {return nil,err}
aead,err:=cipher.NewGCM(block)
if err!=nil {return nil,err}
fmt.Println(aead.NonceSize(),aead.Overhead()) // 12 16
~~~
~cipher.AEAD~ is the sealing/opening interface returned by ~NewGCM~.
~~~go
var box cipher.AEAD = aead
fmt.Println(box.NonceSize()) // 12
~~~

~NonceSize~ gives nonce length; ~Overhead~ gives tag length. A nonce is public but must be fresh for every encryption under a key. Reuse is dangerous even for different payloads. Tests inject distinct byte sequences and verify fresh nonce consumption rather than asserting that two real random outputs happen to differ.

## Injected randomness and associated data
As in the preceding mission, ~io.ReadFull~ reads a complete nonce; ~source~ failure returns no partial packet. Nil here selects ~rand.Reader~ from ~crypto/rand~, a ready-made secure reader, rather than math/rand.
~~~go
if source==nil {source=rand.Reader}
nonce:=make([]byte,aead.NonceSize())
if _,err:=io.ReadFull(source,nonce);err!=nil {return nil,err}
~~~
~Seal(dst,nonce,plaintext,aad)~ returns ciphertext followed by a tag. ~aad~ is authenticated but unencrypted associated data. Bind the ~manifest~ to the seal; changing it must fail even with the correct ~key~ and packet. This example binds a sensor identity, using a different format from the mission:
~~~go
sealed:=aead.Seal(nil,nonce,[]byte("sensor=7"),[]byte("probe-A"))
plain,err:=aead.Open(nil,nonce,sealed,[]byte("probe-A"))
~~~
~Open~ verifies the tag before returning plaintext. Wrong ~key~, nonce, ciphertext, or ~manifest~ all fail authentication. Return one generic ~authentication failed~ error for all Open failures rather than revealing which input was wrong.

## Packet framing and bounded memory
This protocol has four fixed bytes ~OCV1~, a 12-byte nonce, then ciphertext with its 16-byte tag. Use ~append~ with a fresh slice to avoid modifying ~key~, ~manifest~, ~cargo~, or packet. AAD is exactly ~OCV1~ followed by ~manifest~ bytes.
~~~go
aad:=append([]byte("SENSOR1"),label...)
~~~
This example builds an independent prefix. Check ~packet~ length before slicing. Empty ~cargo~ is valid and still produces a tag. Cargo is limited to 65536 bytes; packets to 65568 bytes.

## Common mistakes
~~~go
nonce:=make([]byte,12) // Wrong if not freshly filled for every encryption
plain,_:=aead.Open(nil,nonce,sealed,nil) // Ignores authentication and manifest
~~~

## Your mission
Implement ~SealCargo(key,manifest,cargo []byte,source io.Reader) ([]byte,error)~ and ~OpenCargo(key,manifest,packet []byte) ([]byte,error)~. Support all three AES ~key~ sizes and ~cargo~ sizes zero through 65536 inclusive. Reject invalid keys and oversized ~cargo~ before reading randomness. Use AES-GCM with the ~packet~ and AAD formats above. A non-~nil~ ~source~ supplies every nonce byte; ~nil~ selects ~crypto/rand.Reader~. Read a fresh complete nonce for every encryption, with no cache or alternative random source.

~OpenCargo~ rejects invalid ~key~ lengths, wrong prefix, packets shorter than 32 or larger than 65568 bytes, without panics. Authenticate everything else with Open and the supplied manifest. Authentication failures have exactly the error text ~authentication failed~ without details. All failures return ~nil~ output; reader failures preserve their cause. Do not mutate inputs even when slices have spare capacity. No filesystem, network, Report, or secret output.

Third star: implement ~OpenCargoTo(key,manifest,packet []byte,out io.Writer) error~. The writer is non-nil. Do not write before complete successful authentication. Then write plaintext once, preserving write-error causes; a short write without error returns ~io.ErrShortWrite~. ~Bytes~ written before transport failure may remain. ~io.ErrShortWrite~ is the standard error for a writer accepting fewer bytes than requested.
~~~go
n,err:=out.Write(plain)
if err==nil && n!=len(plain) {err=io.ErrShortWrite}
~~~
`,
			Hints: []string{"Build AES then GCM; validate bounds before reading a nonce.", "Read ~NonceSize~ bytes with ~ReadFull~ on every seal, placing them after OCV1.", "Build AAD in a fresh slice; Open into ~nil~ and return no plaintext on failure."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: true, Hidden: ch232Hidden},
	})
}

const ch232Starter = `package main

import "io"

func SealCargo(key, manifest, cargo []byte, source io.Reader) ([]byte, error) { return nil, io.EOF }
func OpenCargo(key, manifest, packet []byte) ([]byte, error)                  { return nil, io.EOF }
func OpenCargoTo(key, manifest, packet []byte, out io.Writer) error           { return io.EOF }
func main()                                                                   {}
`

const ch232Hidden = `package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
	"orbital/colony/levels"
	"slices"
	"testing"
	"testing/iotest"
)

func TestHiddenCargoSeededInteropAndFreshNonces(t *testing.T) {
	w := levels.Start("23-2-sealed-cargo", io.Discard, 232)
	for _, size := range []int{16, 24, 32} {
		key := make([]byte, size)
		w.Rand.Read(key)
		block, _ := aes.NewCipher(key)
		aead, _ := cipher.NewGCM(block)
		entropy := make([]byte, 12*36)
		w.Rand.Read(entropy)
		source := iotest.OneByteReader(bytes.NewReader(entropy))
		seen := map[string]bool{}
		for i := 0; i < 36; i++ {
			n := w.Rand.Intn(1024)
			if i == 0 {
				n = 0
			}
			if i == 1 {
				n = 65536
			}
			cargo := make([]byte, n)
			w.Rand.Read(cargo)
			manifest := []byte("bay-7\nmanifest")
			originals := [][]byte{slices.Clone(key), slices.Clone(manifest), slices.Clone(cargo)}
			got, err := SealCargo(key, manifest, cargo, source)
			if err != nil || len(got) != len(cargo)+32 || string(got[:4]) != "OCV1" {
				t.Fatal("seal framing", len(got), err)
			}
			nonce := got[4:16]
			if !bytes.Equal(nonce, entropy[i*12:(i+1)*12]) || seen[string(nonce)] {
				t.Fatal("nonce not freshly read")
			}
			seen[string(nonce)] = true
			aad := append([]byte("OCV1"), manifest...)
			expected := append(append([]byte("OCV1"), nonce...), aead.Seal(nil, nonce, cargo, aad)...)
			if !bytes.Equal(got, expected) {
				t.Fatal("AES-GCM interoperability")
			}
			before := slices.Clone(expected)
			plain, err := OpenCargo(key, manifest, expected)
			if err != nil || !bytes.Equal(plain, cargo) || !bytes.Equal(expected, before) {
				t.Fatal("roundtrip", err)
			}
			for j, input := range [][]byte{key, manifest, cargo} {
				if !bytes.Equal(input, originals[j]) {
					t.Fatal("input modified")
				}
			}
		}
	}
}
func hiddenCargoBacking(payload []byte) ([]byte, []byte) {
	backing := bytes.Repeat([]byte{0xcc}, len(payload)+64)
	copy(backing, payload)
	return backing[:len(payload)], backing
}
func TestHiddenCargoPreservesInputBackingArrays(t *testing.T) {
	key, keyBacking := hiddenCargoBacking(bytes.Repeat([]byte{9}, 24))
	manifest, manifestBacking := hiddenCargoBacking([]byte("bay identity"))
	cargo, cargoBacking := hiddenCargoBacking([]byte("reserve"))
	before := [][]byte{slices.Clone(keyBacking), slices.Clone(manifestBacking), slices.Clone(cargoBacking)}
	packet, err := SealCargo(key, manifest, cargo, bytes.NewReader(make([]byte, 12)))
	if err != nil {
		t.Fatal(err)
	}
	for i, backing := range [][]byte{keyBacking, manifestBacking, cargoBacking} {
		if !bytes.Equal(backing, before[i]) {
			t.Fatal("seal modified input backing array", i)
		}
	}
	packet, packetBacking := hiddenCargoBacking(packet)
	packetBefore := slices.Clone(packetBacking)
	plain, err := OpenCargo(key, manifest, packet)
	if err != nil || !bytes.Equal(plain, cargo) || !bytes.Equal(packetBacking, packetBefore) {
		t.Fatal("open mutated packet", err)
	}
	for i, backing := range [][]byte{keyBacking, manifestBacking, cargoBacking} {
		if !bytes.Equal(backing, before[i]) {
			t.Fatal("open modified input backing array", i)
		}
	}
}
func TestHiddenCargoTamperingAndBounds(t *testing.T) {
	key := bytes.Repeat([]byte{7}, 32)
	manifest := []byte("cargo manifest")
	packet, err := SealCargo(key, manifest, []byte("ice reserve"), bytes.NewReader(make([]byte, 12)))
	if err != nil {
		t.Fatal(err)
	}
	for i := 4; i < len(packet); i++ {
		bad := slices.Clone(packet)
		bad[i] ^= 1
		got, err := OpenCargo(key, manifest, bad)
		if err == nil || err.Error() != "authentication failed" || got != nil {
			t.Fatal("tamper accepted/leaked", i, err)
		}
	}
	wrong := slices.Clone(key)
	wrong[0] ^= 1
	for _, tc := range []struct{ k, m []byte }{{wrong, manifest}, {key, []byte("different")}} {
		got, err := OpenCargo(tc.k, tc.m, packet)
		if err == nil || err.Error() != "authentication failed" || got != nil {
			t.Fatal("wrong key/manifest", err)
		}
	}
	for _, n := range []int{0, 1, 15, 17, 23, 25, 31, 33} {
		bad := make([]byte, n)
		if got, err := SealCargo(bad, nil, nil, &hiddenCargoNoRead{t}); err == nil || got != nil {
			t.Fatal("invalid seal key")
		}
		if got, err := OpenCargo(bad, nil, packet); err == nil || got != nil {
			t.Fatal("invalid open key")
		}
	}
	if got, err := SealCargo(key, nil, make([]byte, 65537), &hiddenCargoNoRead{t}); err == nil || got != nil {
		t.Fatal("cargo bound")
	}
	for n := 0; n < 32; n++ {
		if got, err := OpenCargo(key, manifest, packet[:n]); err == nil || got != nil {
			t.Fatal("truncation", n)
		}
	}
	bad := slices.Clone(packet)
	bad[0] = 'X'
	for _, p := range [][]byte{bad, make([]byte, 65569)} {
		if got, err := OpenCargo(key, manifest, p); err == nil || got != nil {
			t.Fatal("format bound")
		}
	}
	wire := errors.New("random offline")
	for _, r := range []io.Reader{iotest.ErrReader(wire), io.MultiReader(bytes.NewReader([]byte{1, 2}), iotest.ErrReader(wire)), bytes.NewReader(make([]byte, 11))} {
		got, err := SealCargo(key, nil, nil, r)
		if err == nil || got != nil {
			t.Fatal("short nonce")
		}
		if !errors.Is(err, wire) && !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatal("lost source cause", err)
		}
	}
}

type hiddenCargoNoRead struct{ t *testing.T }

func (r *hiddenCargoNoRead) Read([]byte) (int, error) {
	r.t.Fatal("read before validation")
	return 0, io.EOF
}
func TestHiddenBonusCargoPublishAfterAuthentication(t *testing.T) {
	key := make([]byte, 16)
	plain := []byte("sealed food")
	packet, err := SealCargo(key, nil, plain, bytes.NewReader(make([]byte, 12)))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := OpenCargoTo(key, nil, packet, &out); err != nil || !bytes.Equal(out.Bytes(), plain) {
		t.Fatal("publish", err)
	}
	packet[len(packet)-1] ^= 1
	writer := &hiddenCargoWriter{}
	if err := OpenCargoTo(key, nil, packet, writer); err == nil || writer.calls != 0 {
		t.Fatal("published unauthenticated bytes")
	}
	packet[len(packet)-1] ^= 1
	wire := errors.New("output offline")
	for _, cause := range []error{wire, nil} {
		writer := &hiddenCargoWriter{err: cause}
		err := OpenCargoTo(key, nil, packet, writer)
		want := cause
		if want == nil {
			want = io.ErrShortWrite
		}
		if !errors.Is(err, want) || writer.calls != 1 {
			t.Fatal("writer failure", err, writer.calls)
		}
	}
}

type hiddenCargoWriter struct {
	calls int
	err   error
}

func (w *hiddenCargoWriter) Write(p []byte) (int, error) { w.calls++; return len(p) / 2, w.err }
`

func init() {
	add(&Level{ID: "23-3-content-address", Chapter: 23, Kind: "lab", Title: "عنوان بالمحتوى", Goal: "خزّن قطع الخزنة مضغوطة بعنوان محسوب من نوعها ومحتواها",
		Requires: []string{"sel:sha1.New"}, Bridge: []int{97}, Setup: seededLabRoom,
		Files: []LabFile{{Name: "main.go", Starter: ch233Starter}},
		Lesson: `# عنوان القطعة هو بصمتها

الخزنة مش بتختار أرقام متزايدة للصناديق. عنوان القطعة بيتحسب من نوعها ومحتواها، فنفس القطعة تتخزن مرة واحدة. المخزن هنا بروتوكول مستعمرة صغير، مش repository لـ Git.

## قبل ما تكتب
- هتحسب البصمة قبل الضغط ولا بعده، وليه؟
- لو الاسم صح بس الملف الموجود بايظ، هتمنع استبدال دليل التلف إزاي؟

## البصمة والنوع والطول
اتعلمنا SHA-256 في 22-4. استورد ~crypto/sha1~ هنا عشان تجربة عنونة قديمة؛ ~sha1.New~ ليها نفس واجهة الكتابة و ~Sum~. SHA-1 مكسورة ضد التصادمات المتعمدة، فمش إثبات أصالة ولا اختيار لنظام أمان جديد؛ SHA-256 أنسب لعنونة جديدة. البصمة مش تشفير، ومش بديل عن GCM.
~~~go
h:=sha1.New()
_,err:=io.WriteString(h,"OC23 note 2\nhi")
id:=fmt.Sprintf("%x",h.Sum(nil))
~~~
~io.WriteString~ تكتب string لـ writer وترجع عدد وخطأ. غلاف القطعة في المهمة نص ~OC23 kind N~ ثم newline واحدة ثم بايتات المحتوى. ~N~ طول المحتوى بالبايتات بصيغة decimal من غير leading zeros. النوع جزء من البصمة، عشان tree وcargo بنفس البايتات يبقوا عنوانين مختلفين. المثال ~note~ للتوضيح فقط، مش نوع مسموح في المهمة.

## ضغط zlib
استورد ~compress/zlib~ و ~bytes~. ~zlib.NewWriter~ يرجّع writer ضاغط؛ ~Close~ لازم تتفحص لأنها بتكمّل السيل وبتكتب checksum. ~bytes.Buffer~ بتجمع الخرج في الذاكرة، و ~Bytes~ تعرض بايتاته.
~~~go
var packed bytes.Buffer
zw:=zlib.NewWriter(&packed)
if _,err:=zw.Write([]byte("relay relay relay"));err!=nil {return err}
if err:=zw.Close();err!=nil {return err}
~~~
البصمة تتحسب من الغلاف الخام، مش من ناتج الضغط. اختلاف مستوى الضغط مينفعش يغير عنوان المحتوى.

~zlib.NewReader~ تفتح السيل المضغوط وبتحتاج ~Close~. ~io.ReadAll~ تقرأ لحد EOF أو خطأ؛ استخدم ~io.LimitReader~ عشان الملف الصغير المضغوط ميفجّرش الذاكرة بعد الفك. ~bytes.NewReader~ تحول slice لـ reader.
~~~go
zr,err:=zlib.NewReader(bytes.NewReader(encoded))
if err!=nil {return err}
defer zr.Close()
raw,err:=io.ReadAll(io.LimitReader(zr,1048641))
if err!=nil {return err}
if len(raw)>1048640 {return fmt.Errorf("object too large")}
~~~
حد الغلاف الخام 1048640 بايت. اقرأ بايت زيادة عشان تكتشف تجاوز الحد بدل قبول نسخة مبتورة. لازم تكمل القراءة للـ EOF في الملف الصالح عشان checksum تتراجع؛ فساد آخر السيل خطأ برضه.

## تحقق ثم خزن
زي 22-3، اكتب ملف مؤقت في نفس مجلد الوجهة، اقفله، وبعدين ~os.Rename~. شيل الملف المؤقت في الفشل. خزن داخل ~filepath.Join(root,"vault-objects",id+".z")~ فقط. كل جذور التست تحت ~t.TempDir()~؛ مفيش مجلد جهاز ثابت.
~~~go
file,err:=os.CreateTemp(directory,".vault-*")
if err!=nil {return err}
defer os.Remove(file.Name())
~~~
القراية لازم تفك الضغط وتراجع البصمة والغلاف والطول قبل إرجاع أي محتوى. ~bytes.IndexByte(raw,'\n')~ تلاقي نهاية أول سطر؛ الغياب يرجع -1. ~bytes.Equal~ تقارن المحتوى، مش عناوين الـ slices.
~~~go
end:=bytes.IndexByte(raw,'\n')
if end<0 {return fmt.Errorf("missing envelope")}
header:=string(raw[:end])
fmt.Println(header)
~~~

## أخطاء شائعة
~~~go
id:=fmt.Sprintf("%x",sha1.Sum(compressed)) // غلط: بصمة الضغط، مش الغلاف الخام
os.WriteFile(existing,newBytes,0600) // غلط: ممكن تمسح دليل تلف القطعة
~~~

## المهمة
نفّذ ~WriteVault(root,kind string,payload []byte) (string,error)~ و ~ReadVault(root,id string) (string,[]byte,error)~. الأنواع المسموحة فقط ~cargo~ و ~tree~ و ~snapshot~. المحتوى من صفر لـ 1048576 بايت شامل. قبل أي I/O ارفض نوع أو حجم غلط. الغلاف الخام بالظبط ~fmt.Sprintf("OC23 %s %d\n",kind,len(payload))~ وبعده payload؛ SHA-1 بـ ~sha1.New~ على الغلاف كله، وID lowercase hex من 40 حرف. القطعة الفاضية من نوع ~cargo~ عنوانها ~c9417a3932fe31683d560dd513662e1b3355360d~.

أنشئ مجلد ~vault-objects~ لو ناقص بـ 0700. الملف ~id.z~ فيه الغلاف كله مضغوط zlib؛ كتابة ملف جديد عن طريق مؤقت ثم rename من غير ملفات مؤقتة متروكة. لو القطعة موجودة فكها وراجعها وقارن النوع والمحتوى؛ التطابق نجاح من غير إعادة كتابة، وأي فساد خطأ من غير لمس الملف الموجود. احتفظ بأسباب أخطاء الملفات. المخزن محلي موثوق من غير symlinks أو كاتب موازي؛ لا تزور أي مكان تاني ولا تشغّل أوامر Git.

~ReadVault~ ترفض ID مش 40 حرف من ~0-9a-f~ قبل I/O، وتقرأ القطعة المطلوبة فقط. ارفض zlib بايظة أو ناقصة أو checksum غلط، وغلاف خام فوق 1048640 بايت، وبادئة أو نوع أو طول غير canonical، ومحتوى فوق الحد، وبصمة مختلفة عن الاسم. فشل القراية يرجّع ~"",nil,error~، وفشل الكتابة ~"",error~. لا تعديل للمدخلات ولا للمخزن أثناء القراية، ولا Report. احتفظ بتوقيعات الدوال وملف ~main.go~ للمهمة اللي بعدها.
`,
		Hints: []string{"كوّن الغلاف الخام مرة، واحسب بصمته قبل zlib.", "في القراية حدّد حجم الفك، وافحص الغلاف canonical والبصمة قبل رجوع المحتوى.", "لو الاسم موجود استخدم ~ReadVault~ وقارن؛ لو جديد اكتب مؤقت واقفله قبل rename."},
		En: &LevelText{Title: "Content Address", Goal: "Store compressed vault objects addressed by their kind and content", Files: []LabFile{{Name: "main.go", Starter: ch233Starter}},
			Lesson: `# Address an Object by Its Content

The vault does not allocate increasing crate numbers. An object's address is computed from its kind and content, storing identical objects once. This is a small colony protocol, not a Git repository.

## Before you write
- Should the digest be computed before or after compression, and why?
- If a filename is correct but its bytes are corrupt, how will you preserve the evidence?

## Digest, kind, and length
SHA-256 appeared in 22-4. Import ~crypto/sha1~ here to explore legacy addressing; ~sha1.New~ supports the same writing and ~Sum~ interface. SHA-1 is broken against deliberate collisions. It is neither proof of authenticity nor a choice for a new security design; SHA-256 is preferable for new addressing schemes. Hashing is not encryption or a replacement for GCM.
~~~go
h:=sha1.New()
_,err:=io.WriteString(h,"OC23 note 2\nhi")
id:=fmt.Sprintf("%x",h.Sum(nil))
~~~
~io.WriteString~ writes a string to a writer, returning a count and error. Our envelope is the text ~OC23 kind N~, one newline, then payload bytes. ~N~ is the byte count in decimal without leading zeros. Kind participates in the digest, so tree and ~cargo~ with identical bytes have different addresses. The example's ~note~ is illustrative, not an allowed mission kind.

## zlib compression
Import ~compress/zlib~ and ~bytes~. ~zlib.NewWriter~ creates a compressing writer. Check ~Close~ because it finishes the stream and writes its checksum. ~bytes.Buffer~ collects output in memory; ~Bytes~ exposes its bytes.
~~~go
var packed bytes.Buffer
zw:=zlib.NewWriter(&packed)
if _,err:=zw.Write([]byte("relay relay relay"));err!=nil {return err}
if err:=zw.Close();err!=nil {return err}
~~~
Hash the raw envelope, not compressed output. Changing compression level must not change a content address.

~zlib.NewReader~ opens the compressed stream and must be closed. ~io.ReadAll~ reads to EOF or an error; ~io.LimitReader~ bounds decompression so a small compressed file cannot exhaust memory. ~bytes.NewReader~ exposes a byte slice as a reader.
~~~go
zr,err:=zlib.NewReader(bytes.NewReader(encoded))
if err!=nil {return err}
defer zr.Close()
raw,err:=io.ReadAll(io.LimitReader(zr,1048641))
if err!=nil {return err}
if len(raw)>1048640 {return fmt.Errorf("object too large")}
~~~
The raw envelope limit is 1048640 bytes. Read one extra byte to detect excess rather than accepting a truncated prefix. Valid streams must be read to EOF to check the checksum; corruption at the end is still failure.

## Verify before storing
As in 22-3, write a temporary file in the destination directory, close it, then ~os.Rename~. Remove temporaries on failure. Store only at ~filepath.Join(root,"vault-objects",id+".z")~. Every test root is inside ~t.TempDir()~, with no fixed machine directory.
~~~go
file,err:=os.CreateTemp(directory,".vault-*")
if err!=nil {return err}
defer os.Remove(file.Name())
~~~
Reads decompress and verify digest, envelope, and length before returning content. ~bytes.IndexByte(raw,'\n')~ finds the first header terminator, returning -1 if absent. ~bytes.Equal~ compares content rather than slice addresses.
~~~go
end:=bytes.IndexByte(raw,'\n')
if end<0 {return fmt.Errorf("missing envelope")}
header:=string(raw[:end])
fmt.Println(header)
~~~

## Common mistakes
~~~go
id:=fmt.Sprintf("%x",sha1.Sum(compressed)) // Wrong: hashes compression, not raw envelope
os.WriteFile(existing,newBytes,0600) // Wrong: may destroy evidence of corruption
~~~

## Your mission
Implement ~WriteVault(root,kind string,payload []byte) (string,error)~ and ~ReadVault(root,id string) (string,[]byte,error)~. Allowed kinds are exactly ~cargo~, ~tree~, and ~snapshot~. Payload size is zero through 1048576 bytes inclusive. Reject invalid kinds and sizes before I/O. The raw envelope is exactly ~fmt.Sprintf("OC23 %s %d\n",kind,len(payload))~ followed by payload. Hash the entire envelope using ~sha1.New~; the ID is 40 lowercase hexadecimal characters. Empty ~cargo~ has address ~c9417a3932fe31683d560dd513662e1b3355360d~.

Create ~vault-objects~ if absent with 0700. File ~id.z~ contains the entire zlib-compressed envelope. Publish new files through a temporary and rename, leaving no temporaries behind. If an object already exists, decompress, verify, and compare kind and content. Identical content succeeds without rewriting; corruption fails without touching the existing file. Preserve filesystem error causes. The local store is trusted, without symlinks or concurrent writers. Visit no other locations and never run Git commands.

~ReadVault~ rejects IDs other than 40 characters from ~0-9a-f~ before I/O and reads only the named object. Reject corrupt or truncated zlib or bad checksums, raw envelopes exceeding 1048640 bytes, noncanonical prefix/kind/length, oversized payloads, and digest mismatches. Read failures return ~"",nil,error~; write failures return ~"",error~. Do not mutate inputs or write during reads. No Report. Keep signatures and ~main.go~ for the following mission.
`,
			Hints: []string{"Build the raw envelope once and hash it before zlib.", "Bound decompression and verify the canonical envelope and digest before returning content.", "For existing IDs use ~ReadVault~ and compare; publish new objects through a closed temporary file."}},
		Lab: &Lab{ReqFile: "main.go", Bonus: false, Hidden: ch233Hidden},
	})
}

const ch233Starter = `package main

import "io"

func WriteVault(root, kind string, payload []byte) (string, error) { return "", io.EOF }
func ReadVault(root, id string) (string, []byte, error)            { return "", nil, io.EOF }
func main()                                                        {}
`

const ch233Hidden = `package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"orbital/colony/levels"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func hiddenVaultID(raw []byte) string { return fmt.Sprintf("%x", sha1.Sum(raw)) }
func hiddenVaultPack(t *testing.T, raw []byte) []byte {
	t.Helper()
	var out bytes.Buffer
	z := zlib.NewWriter(&out)
	if _, err := z.Write(raw); err != nil {
		t.Fatal(err)
	}
	if err := z.Close(); err != nil {
		t.Fatal(err)
	}
	return out.Bytes()
}
func TestHiddenVaultKnownVectorAndSeededRoundtrips(t *testing.T) {
	root := t.TempDir()
	w := levels.Start("23-3-content-address", io.Discard, 233)
	id, err := WriteVault(root, "cargo", nil)
	if err != nil || id != "c9417a3932fe31683d560dd513662e1b3355360d" {
		t.Fatal("known empty cargo", id, err)
	}
	for i := 0; i < 28; i++ {
		n := w.Rand.Intn(8192)
		if i == 0 {
			n = 1048576
		}
		payload := make([]byte, n)
		w.Rand.Read(payload)
		before := slices.Clone(payload)
		kind := []string{"cargo", "tree", "snapshot"}[i%3]
		raw := append([]byte(fmt.Sprintf("OC23 %s %d\n", kind, len(payload))), payload...)
		id, err := WriteVault(root, kind, payload)
		if err != nil || id != hiddenVaultID(raw) || !bytes.Equal(payload, before) {
			t.Fatal("write content address", id, err)
		}
		file := filepath.Join(root, "vault-objects", id+".z")
		packed, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		z, err := zlib.NewReader(bytes.NewReader(packed))
		if err != nil {
			t.Fatal(err)
		}
		unpacked, err := io.ReadAll(z)
		z.Close()
		if err != nil || !bytes.Equal(unpacked, raw) {
			t.Fatal("zlib envelope", err)
		}
		k, got, err := ReadVault(root, id)
		if err != nil || k != kind || !bytes.Equal(got, payload) {
			t.Fatal("read", err)
		}
		// Alternate valid compression must be accepted and not rewritten on deduplication.
		var alternate bytes.Buffer
		zw, _ := zlib.NewWriterLevel(&alternate, zlib.NoCompression)
		zw.Write(raw)
		zw.Close()
		if err := os.WriteFile(file, alternate.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
		again, err := WriteVault(root, kind, payload)
		after, _ := os.ReadFile(file)
		if err != nil || again != id || !bytes.Equal(after, alternate.Bytes()) {
			t.Fatal("dedup rewrites", err)
		}
	}
	entries, err := os.ReadDir(filepath.Join(root, "vault-objects"))
	if err != nil || len(entries) != 29 {
		t.Fatal("object count/temp leak", len(entries), err)
	}
}
func TestHiddenVaultCorruptionAndCanonicalEnvelope(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "vault-objects")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	for _, raw := range [][]byte{
		[]byte("OC23 cargo 01\nx"), []byte("OC23 cargo 2\nx"), []byte("OC23 wrong 1\nx"), []byte("OC23 cargo -1\n"),
		[]byte("OTHER cargo 1\nx"), []byte("OC23 cargo 1"), []byte("OC23  cargo 1\nx"),
		append([]byte("OC23 cargo 1048577\n"), make([]byte, 1048577)...),
		bytes.Repeat([]byte{'x'}, 1048641),
	} {
		id := hiddenVaultID(raw)
		file := filepath.Join(dir, id+".z")
		if err := os.WriteFile(file, hiddenVaultPack(t, raw), 0600); err != nil {
			t.Fatal(err)
		}
		k, p, err := ReadVault(root, id)
		if err == nil || k != "" || p != nil {
			t.Fatal("malformed envelope accepted", len(raw), err)
		}
	}
	id, err := WriteVault(root, "cargo", []byte("reserve"))
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(dir, id+".z")
	valid, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	checksum := slices.Clone(valid)
	checksum[len(checksum)-1] ^= 1
	variants := [][]byte{nil, []byte("not zlib"), valid[:len(valid)/2], checksum, hiddenVaultPack(t, []byte("OC23 cargo 7\nchanged"))}
	for _, bad := range variants {
		if err := os.WriteFile(file, bad, 0600); err != nil {
			t.Fatal(err)
		}
		k, p, err := ReadVault(root, id)
		if err == nil || k != "" || p != nil {
			t.Fatal("corruption accepted")
		}
		got, err := WriteVault(root, "cargo", []byte("reserve"))
		after, _ := os.ReadFile(file)
		if err == nil || got != "" || !bytes.Equal(after, bad) {
			t.Fatal("corrupt object overwritten")
		}
	}
}
func TestHiddenVaultValidationAndIOFailures(t *testing.T) {
	root := filepath.Join(t.TempDir(), "absent")
	for _, id := range []string{"", strings.Repeat("a", 39), strings.Repeat("A", 40), strings.Repeat("g", 40), "../" + strings.Repeat("a", 40)} {
		k, p, err := ReadVault(root, id)
		if err == nil || k != "" || p != nil {
			t.Fatal("invalid id")
		}
	}
	for _, tc := range []struct {
		k string
		p []byte
	}{{"blob", nil}, {"cargo\n", nil}, {"cargo", make([]byte, 1048577)}} {
		if id, err := WriteVault(root, tc.k, tc.p); err == nil || id != "" {
			t.Fatal("invalid write")
		}
	}
	if _, err := os.Stat(root); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("invalid request created store")
	}
	k, p, err := ReadVault(root, strings.Repeat("0", 40))
	if !errors.Is(err, fs.ErrNotExist) || k != "" || p != nil {
		t.Fatal("missing cause", err)
	}
	blocked := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocked, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	if id, err := WriteVault(blocked, "cargo", nil); err == nil || id != "" {
		t.Fatal("blocked root")
	}
	after, _ := os.ReadFile(blocked)
	if string(after) != "unchanged" {
		t.Fatal("root modified")
	}
}
`

func init() {
	add(&Level{ID: "23-4-snapshot-commit", Chapter: 23, Kind: "lab", Title: "لقطة و commit", Goal: "ابني فهرس شحنات وشجرة مرتبة وcommit ببيانات وقت محقونة",
		Requires: []string{}, Bridge: []int{97}, Setup: seededLabRoom,
		CarryFrom: "23-3-content-address",
		Files:     []LabFile{{Name: "main.go", Starter: ch233Starter}, {Name: "snapshot.go", Starter: ch234Starter}},
		Lesson: `# لقطة تجهيز طلعة

الطاقم بيجهّز شحنات بأسماء نسبية، والفهرس بيربط كل اسم بعنوان محتواه. الـ commit يثبّت شجرة الأسماء ويربطها باللقطة السابقة من غير تغيير القطع القديمة. دي لقطة شحنات، مش أوامر Git ولا نسخ ملفات جهازك.

## قبل ما تكتب
- إزاي نفس الشحنات بترتيبين إضافة مختلفين تطلع نفس الشجرة؟
- إيه اللي لازم يتراجع قبل كتابة أي tree أو snapshot؟

## الفهرس حالة صريحة
~SupplyIndex~ نوع map من اسم شحنة لـ ID قطعة cargo. المنادي يملكه ويقدر يحفظه أو يعرضه؛ مفيش index مخفية على جهاز اللاعب. ~StageSupply~ تكتب المحتوى في مخزن 23-3 وبعدين تحدّث الربط لو الكتابة نجحت بس.
~~~go
index:=make(SupplyIndex)
err:=StageSupply(root,index,"bay/ice",[]byte("cold"))
fmt.Println(index["bay/ice"],err)
~~~
~SupplyIndex~ لازم تتعمل بـ ~make~ قبل staging؛ ~nil~ للـ staging خطأ، لكن ~nil~ للـ commit معناها فهرس فاضي ومرفوض. لو اسم موجود، staging تانية تستبدل ربطه بعد نجاح التخزين؛ القطعة القديمة بتفضل في الخزنة.

## أسماء سليمة وشجرة حتمية
~fs.ValidPath~ من ~io/fs~ تفحص أسماء slash-separated من غير مكونات فاضية أو نقطة أو رجوع للأب؛ لكنها بتقبل ~"."~ كجذر. المهمة تمنع الجذر نفسه، وbackslash وNUL وtab وCR وLF.
~~~go
fmt.Println(fs.ValidPath("bay/cell"),fs.ValidPath("../cell")) // true false
~~~
~strings.ContainsAny~ تفحص وجود أي حرف من مجموعة المنع.
~~~go
fmt.Println(strings.ContainsAny("bay\tcell","\t\r\n")) // true
~~~

الاسم هنا label للشحنة، مش مسار بنفتحه على الجهاز. ~bay~ و ~bay/cell~ مينفعش يبقوا شحنتين مع بعض: الأول هيبقى ملف ومجلد في نفس الوقت. ~strings.HasPrefix~ تكشف الحد لو أضفت slash، مش prefix نصي عادي؛ ~bay~ و ~bay2/cell~ صالحين.
~~~go
conflict:=strings.HasPrefix("bay/cell","bay"+"/")
fmt.Println(conflict) // true
~~~

اجمع الأسماء ورتبها بـ ~slices.Sort~ من الفصل السابق. لكل اسم اكتب سطر ~name\tID\n~. دي tree مسطحة بأسماء كاملة، مش تمثيل Git الثنائي؛ الترتيب Go string order. محتوى الشجرة يتخزن بالنوع ~tree~ في ~WriteVault~، ومحتوى كل شحنة بالنوع ~cargo~.
~~~go
var tree strings.Builder
fmt.Fprintf(&tree,"%s\t%s\n","bay/cell",cargoID)
~~~
~\t~ tab واحدة و ~\n~ newline واحدة في string Go، مش backslash حرفية. افحص كل ID بـ ~ReadVault~ وتأكد من النوع قبل بناء tree؛ الفهرس public وممكن المنادي يكون حط فيه قيمة بايظة.

## commit بوقت محقون
~SupplyStamp~ فيها ~Tick int64~ و ~Message string~. دي بيانات جاهزة؛ مينفعش تنادي الوقت الحقيقي أو البيئة أو هوية الجهاز. ~SupplyCommit~ فيها بالترتيب ~Tree string~, ~Parent string~, ~Tick int64~, ~Message string~. استخدم ~json.Marshal~ من 18-2 على الـ struct من غير tags؛ المفاتيح بنفس أسماء الحقول الكبيرة وبنفس الترتيب ومن غير newline نهائية.
~~~go
stamp:=SupplyStamp{Tick:73,Message:"food ready"}
record:=SupplyCommit{Tree:treeID,Parent:previousID,Tick:stamp.Tick,Message:stamp.Message}
encoded,err:=json.Marshal(record)
~~~
خزّن JSON بالنوع ~snapshot~. ~Parent~ فاضية للجذر؛ غير كده ID لـ snapshot موجودة وصالحة. ارجع ID جديدة؛ المنادي هو اللي يحتفظ برأس السلسلة. مفيش HEAD أو index files ضمنية، ولا قراءة worktree. نفس الفهرس والأب وstamp يطلعوا نفس commit.

عند قراءة أب فك JSON بـ ~json.Unmarshal~ لنفس الـ struct، وارجع Marshal وقارن البايتات عشان ترفض حقول ناقصة أو زيادة أو ترتيب وتنسيق مختلف. افحص metadata وشجرة الأب كمان. ده schema لعقدنا، مش قاعدة إن كل JSON لازم canonical.
~~~go
var old SupplyCommit
if err:=json.Unmarshal(encoded,&old);err!=nil {return "",err}
canonical,err:=json.Marshal(old)
if err!=nil || !bytes.Equal(encoded,canonical) {return "",fmt.Errorf("invalid snapshot")}
~~~

## أخطاء شائعة
~~~go
index[name]=id // غلط قبل نجاح WriteVault
for name,id:=range index {fmt.Fprintf(&tree,"%s\t%s\n",name,id)} // ترتيب map مش ثابت
stamp.Tick=time.Now().Unix() // غلط: تجاهلت الوقت المحقون
~~~

## المهمة
كمّل ~main.go~ من 23-3، واكتب الدوال الجديدة في ~snapshot.go~. النقل التلقائي يحافظ على اسم الملف. لا تغيّر توقيعات ~WriteVault~ و ~ReadVault~؛ النسخة الابتدائية البديلة بتتبني بس لسه مش محلولة. نفّذ ~StageSupply(root string,index SupplyIndex,name string,payload []byte) error~ و ~CommitSupply(root string,index SupplyIndex,parent string,stamp SupplyStamp) (string,error)~.

الفهرس حدّه 128 اسم، كل اسم 1 لـ 255 بايت، يحقق ~fs.ValidPath~ وغير ~"."~ وخالي من ~\\~ وNUL وtab وCR وLF. ارفض أي تعارض ملف/مجلد. في staging تحقق من الاسم والفهرس غير ~nil~ والسعة والتعارض قبل I/O؛ المخزن ياخد payload بالنوع ~cargo~، والفهرس يتغير بعد نجاح التخزين بس. الاستبدال مسموح حتى عند امتلاء الفهرس. أسماء الفهرس الموجودة صالحة وقت staging. فشل staging لا يغير الفهرس.

الـ commit ترفض الفهرس الفاضي أو فوق الحد، وأي اسم أو تعارض غلط، وID غايبة أو تالفة أو نوعها مش cargo. ~Tick~ لازم غير سالب؛ ~Message~ مش whitespace وطولها لحد 512 بايت؛ احتفظ بالنص الأصلي بما فيه مسافات وnewlines. الأب الفاضي جذر؛ غير كده لازم snapshot بصيغة JSON canonical الموصوفة وبـ metadata سليمة. ~Parent~ جوه الأب فاضية أو 40 حرف ~0-9a-f~؛ مش مطلوب تمشي الأجداد. ~Tree~ جوه الأب لازم تشير لقطعة tree موجودة وسليمة، مش مطلوب تحليل سطور شجرة الأب.

جهّز tree بالسطور المرتبة المذكورة واحسب عنوانها قبل الكتابة. ارفض commit لو شجرتها نفس شجرة الأب، حتى لو الوقت أو الرسالة اختلفوا. كمّل كل التحققات قبل كتابة أي قطعة جديدة. بعدها اكتب tree ثم snapshot بالـ JSON المذكورة؛ ارجع ID snapshot. فشل الكتابة يرجّع string فاضية ويحافظ على سبب خطأ I/O؛ tree نجحت قبل فشل snapshot ممكن تفضل من غير مرجع. الفهرس والمدخلات وكل القطع القديمة ثابتين وقت commit، حتى في الفشل. لا تزامن ولا Git ولا وقت حقيقي ولا Report. كل المخازن الحقيقية داخل ~t.TempDir()~ فقط.
`,
		Hints: []string{"اتحقق من الاسم والسعة قبل staging، وعدّل map بعد نجاح التخزين بس.", "رتّب الأسماء وافحص كل ~cargo~ بـ ~ReadVault~؛ حضّر سطور tree من غير كتابة.", "راجع الأب والـ metadata والشجرة غير المتغيرة قبل ~WriteVault~؛ استخدم struct لترتيب JSON."},
		En: &LevelText{Title: "Snapshot Commit", Goal: "Build a cargo index, canonical tree, and commit with injected metadata", Files: []LabFile{{Name: "main.go", Starter: ch233Starter}, {Name: "snapshot.go", Starter: ch234Starter}},
			Lesson: `# Snapshot a Flight's Supplies

The crew stages supplies under relative names; an index maps each name to ~cargo~ content. A commit freezes the name tree and links to the previous snapshot without changing old objects. This is a ~cargo~ snapshot, not Git commands or a copy of machine files.

## Before you write
- How will different staging orders produce the same tree?
- What must be verified before writing any tree or snapshot?

## The index is explicit state
~SupplyIndex~ is a map type from supply names to ~cargo~ object IDs. The caller owns it and may save or display it; there is no hidden index on the player's machine. ~StageSupply~ writes content using 23-3, then updates the mapping only on success.
~~~go
index:=make(SupplyIndex)
err:=StageSupply(root,index,"bay/ice",[]byte("cold"))
fmt.Println(index["bay/ice"],err)
~~~
Allocate ~SupplyIndex~ with ~make~ before staging. Nil is a staging error; for commits ~nil~ is an empty, rejected index. Staging an existing name replaces its mapping after storage succeeds; the old object remains addressable.

## Valid names and a deterministic tree
~fs.ValidPath~ from ~io/fs~ checks slash-separated names without empty, dot, or parent segments, but accepts ~"."~ for a root. This mission rejects that root, backslash, NUL, tab, CR, and LF.
~~~go
fmt.Println(fs.ValidPath("bay/cell"),fs.ValidPath("../cell")) // true false
~~~
~strings.ContainsAny~ tests whether any character belongs to a forbidden set.
~~~go
fmt.Println(strings.ContainsAny("bay\tcell","\t\r\n")) // true
~~~

A name is a supply label, not a machine path to open. ~bay~ and ~bay/cell~ cannot both be supplies because bay would be a file and a directory. ~strings.HasPrefix~ checks this boundary when the slash is included; bay and bay2/cell are compatible.
~~~go
conflict:=strings.HasPrefix("bay/cell","bay"+"/")
fmt.Println(conflict) // true
~~~
Collect names and sort with ~slices.Sort~ from the preceding chapter. Write one ~name\tID\n~ line per name. This is a flat tree of full names, not Git's binary representation. Use Go string order. Write the tree with kind ~tree~ through ~WriteVault~; each supply has kind ~cargo~.
~~~go
var tree strings.Builder
fmt.Fprintf(&tree,"%s\t%s\n","bay/cell",cargoID)
~~~
~\t~ is one tab and ~\n~ one newline in a Go string, not literal backslashes. Verify every ID through ~ReadVault~ and check its kind before constructing a tree; the public index may have been altered by its caller.

## Commits with injected time
~SupplyStamp~ contains ~Tick int64~ and ~Message string~. These values are supplied, never queried from real time, environment, or machine identity. ~SupplyCommit~ has fields in this order: ~Tree string~, ~Parent string~, ~Tick int64~, ~Message string~. Apply ~json.Marshal~ from 18-2 to the struct without tags. Keys retain their capitalized field names and order, with no final newline.
~~~go
stamp:=SupplyStamp{Tick:73,Message:"food ready"}
record:=SupplyCommit{Tree:treeID,Parent:previousID,Tick:stamp.Tick,Message:stamp.Message}
encoded,err:=json.Marshal(record)
~~~
Store JSON with kind ~snapshot~. Empty ~Parent~ denotes a root; otherwise it names an existing valid snapshot. Return the new ID and let the caller retain its chain head. There are no implicit HEAD or index files and no worktree reads. Identical index, parent, and stamp produce identical commits.

When reading a parent, ~json.Unmarshal~ into the same struct, marshal again, and compare bytes to reject missing or extra fields, different order, or different whitespace. Validate metadata and the parent's tree too. This is our protocol's schema rule, rather than a rule for all JSON.
~~~go
var old SupplyCommit
if err:=json.Unmarshal(encoded,&old);err!=nil {return "",err}
canonical,err:=json.Marshal(old)
if err!=nil || !bytes.Equal(encoded,canonical) {return "",fmt.Errorf("invalid snapshot")}
~~~

## Common mistakes
~~~go
index[name]=id // Wrong before WriteVault succeeds
for name,id:=range index {fmt.Fprintf(&tree,"%s\t%s\n",name,id)} // Map order is unstable
stamp.Tick=time.Now().Unix() // Wrong: ignores injected time
~~~

## Your mission
Continue ~main.go~ from 23-3 and implement new functions in ~snapshot.go~. Automatic carry preserves filenames. Retain ~WriteVault~ and ~ReadVault~ signatures; fallback starters compile but remain unsolved. Implement ~StageSupply(root string,index SupplyIndex,name string,payload []byte) error~ and ~CommitSupply(root string,index SupplyIndex,parent string,stamp SupplyStamp) (string,error)~.

An index holds at most 128 names. Names are 1 through 255 bytes, satisfy ~fs.ValidPath~, differ from ~"."~, and contain no backslash, NUL, tab, CR, or LF. Reject file/directory prefix conflicts. Staging checks the name, non-~nil~ index, capacity, and conflicts before I/O. Store payload as ~cargo~ and update the index only after success. Replacements are allowed at capacity. Existing index names are valid when staging. Failed staging leaves the index unchanged.

Commits reject empty or oversized indexes, invalid names or conflicts, and missing, corrupt, or non-~cargo~ IDs. ~Tick~ must be nonnegative; ~Message~ must contain non-whitespace and be at most 512 bytes. Preserve the original message including spaces and newlines. Empty parent means a root; otherwise it must be a snapshot with the canonical JSON above and valid metadata. Its ~Parent~ must be empty or 40 characters from ~0-9a-f~; ancestry traversal is unnecessary. Its ~Tree~ must name an existing verified tree object; parsing the parent's tree lines is unnecessary.

Prepare the sorted tree lines and compute their address before writing. Reject a commit whose tree equals its parent's, even with changed metadata. Complete all validation before writing any new objects. Then write tree followed by canonical-JSON snapshot and return the snapshot ID. Write failure returns an empty string and preserves I/O causes; a tree successfully written before snapshot failure may remain unreferenced. Commit never changes the index, inputs, or old objects, including on failure. No concurrency, Git, real time, or Report. Real test stores live only inside ~t.TempDir()~.
`,
			Hints: []string{"Validate name and capacity before staging; update the map only after storage succeeds.", "Sort names and verify every ~cargo~ through ~ReadVault~; prepare tree bytes without writing.", "Validate parent, metadata, and changed tree before ~WriteVault~; use a struct for JSON order."}},
		Lab: &Lab{ReqFile: "snapshot.go", Bonus: false, Hidden: ch234Hidden},
	})
}

const ch234Starter = `package main

import "io"

type SupplyIndex map[string]string
type SupplyStamp struct {
	Tick    int64
	Message string
}
type SupplyCommit struct {
	Tree    string
	Parent  string
	Tick    int64
	Message string
}

func StageSupply(root string, index SupplyIndex, name string, payload []byte) error { return io.EOF }
func CommitSupply(root string, index SupplyIndex, parent string, stamp SupplyStamp) (string, error) {
	return "", io.EOF
}
`

const ch234Hidden = `package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"orbital/colony/levels"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func hiddenSupplyTree(index SupplyIndex) []byte {
	names := make([]string, 0, len(index))
	for name := range index {
		names = append(names, name)
	}
	slices.Sort(names)
	var out strings.Builder
	for _, name := range names {
		fmt.Fprintf(&out, "%s\t%s\n", name, index[name])
	}
	return []byte(out.String())
}
func hiddenSupplyID(kind string, payload []byte) string {
	raw := append([]byte(fmt.Sprintf("OC23 %s %d\n", kind, len(payload))), payload...)
	return fmt.Sprintf("%x", sha1.Sum(raw))
}
func hiddenSupplyFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		p, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		out[name] = string(p)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func TestHiddenSupplySeededStagingOrderingAndHistory(t *testing.T) {
	root := t.TempDir()
	w := levels.Start("23-4-snapshot-commit", io.Discard, 234)
	index := make(SupplyIndex)
	reverse := make(SupplyIndex)
	for i := 0; i < 30; i++ {
		name := fmt.Sprintf("bay%02d/cell", i)
		payload := make([]byte, w.Rand.Intn(128))
		w.Rand.Read(payload)
		before := slices.Clone(payload)
		if err := StageSupply(root, index, name, payload); err != nil || !bytes.Equal(payload, before) {
			t.Fatal("stage", err)
		}
		if index[name] != hiddenSupplyID("cargo", payload) {
			t.Fatal("stage stores wrong kind")
		}
	}
	names := make([]string, 0, len(index))
	for name := range index {
		names = append(names, name)
	}
	slices.Sort(names)
	for i := len(names) - 1; i >= 0; i-- {
		reverse[names[i]] = index[names[i]]
	}
	stamp := SupplyStamp{Tick: 773, Message: " supplies\nجاهزة 🚀 "}
	before := maps.Clone(index)
	first, err := CommitSupply(root, index, "", stamp)
	if err != nil {
		t.Fatal(err)
	}
	repeat, err := CommitSupply(root, reverse, "", stamp)
	if err != nil || repeat != first {
		t.Fatal("noncanonical order", err)
	}
	tree := hiddenSupplyTree(index)
	treeID := hiddenSupplyID("tree", tree)
	expected := SupplyCommit{Tree: treeID, Tick: stamp.Tick, Message: stamp.Message}
	encoded, _ := json.Marshal(expected)
	if first != hiddenSupplyID("snapshot", encoded) {
		t.Fatal("commit envelope/JSON")
	}
	kind, actual, err := ReadVault(root, first)
	if err != nil || kind != "snapshot" || !bytes.Equal(actual, encoded) {
		t.Fatal("commit document", err)
	}
	kind, actual, err = ReadVault(root, treeID)
	if err != nil || kind != "tree" || !bytes.Equal(actual, tree) {
		t.Fatal("tree document", err)
	}
	for _, s := range []SupplyStamp{{Tick: 774, Message: stamp.Message}, {Tick: 773, Message: "changed"}} {
		other, err := CommitSupply(root, index, "", s)
		if err != nil || other == first {
			t.Fatal("metadata ignored", err)
		}
	}
	files := hiddenSupplyFiles(t, root)
	if id, err := CommitSupply(root, index, first, SupplyStamp{Tick: 900, Message: "unchanged"}); err == nil || id != "" {
		t.Fatal("empty child accepted")
	}
	if !reflect.DeepEqual(files, hiddenSupplyFiles(t, root)) || !maps.Equal(index, before) {
		t.Fatal("rejected child wrote state")
	}
	oldCargo := index[names[0]]
	if err := StageSupply(root, index, names[0], []byte("replacement")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := ReadVault(root, oldCargo); err != nil {
		t.Fatal("old object erased", err)
	}
	child, err := CommitSupply(root, index, first, SupplyStamp{Tick: 901, Message: "changed cargo"})
	if err != nil {
		t.Fatal(err)
	}
	_, encoded, err = ReadVault(root, child)
	if err != nil {
		t.Fatal(err)
	}
	var commit SupplyCommit
	if err := json.Unmarshal(encoded, &commit); err != nil || commit.Parent != first || commit.Tree == treeID {
		t.Fatal("history", commit, err)
	}
	_, original, err := ReadVault(root, first)
	originalExpected, _ := json.Marshal(expected)
	if err != nil || !bytes.Equal(original, originalExpected) {
		t.Fatal("parent modified")
	}
}
func TestHiddenSupplyPreflightAndCapacity(t *testing.T) {
	root := t.TempDir()
	index := make(SupplyIndex)
	if err := StageSupply(root, index, "bay", []byte("one")); err != nil {
		t.Fatal(err)
	}
	invalid := []string{"", ".", "/absolute", "../x", "a/../b", "a//b", "a/./b", "a/", "a\\b", "a\x00b", "a\tb", "a\nb", "a\rb", strings.Repeat("x", 256), "bay/cell"}
	for _, name := range invalid {
		before := maps.Clone(index)
		files := hiddenSupplyFiles(t, root)
		if err := StageSupply(root, index, name, nil); err == nil {
			t.Fatal("invalid staged name", name)
		}
		if !maps.Equal(index, before) || !reflect.DeepEqual(files, hiddenSupplyFiles(t, root)) {
			t.Fatal("invalid stage mutated state")
		}
	}
	for _, name := range []string{strings.Repeat("x", 255), "bay2/cell", "قمر/food"} {
		if err := StageSupply(root, index, name, nil); err != nil {
			t.Fatal("valid name", err)
		}
	}
	if err := StageSupply(root, nil, "cell", nil); err == nil {
		t.Fatal("nil index")
	}
	nested := make(SupplyIndex)
	if err := StageSupply(root, nested, "a/b", nil); err != nil {
		t.Fatal(err)
	}
	if err := StageSupply(root, nested, "a", nil); err == nil {
		t.Fatal("reverse prefix")
	}
	full := make(SupplyIndex)
	for i := 0; i < 128; i++ {
		if err := StageSupply(root, full, fmt.Sprintf("s%03d", i), nil); err != nil {
			t.Fatal(err)
		}
	}
	files := hiddenSupplyFiles(t, root)
	if err := StageSupply(root, full, "extra", []byte("new")); err == nil || len(full) != 128 || !reflect.DeepEqual(files, hiddenSupplyFiles(t, root)) {
		t.Fatal("capacity")
	}
	if err := StageSupply(root, full, "s000", []byte("replace")); err != nil || len(full) != 128 {
		t.Fatal("replace at capacity", err)
	}
	if id, err := CommitSupply(root, full, "", SupplyStamp{Tick: 0, Message: "boundary"}); err != nil || id == "" {
		t.Fatal("full index commit", err)
	}
	blocked := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocked, []byte("unchanged"), 0600); err != nil {
		t.Fatal(err)
	}
	before := maps.Clone(index)
	if err := StageSupply(blocked, index, "new", nil); err == nil || !maps.Equal(before, index) {
		t.Fatal("stage I/O failure")
	}
	if err := StageSupply(root, index, "large", make([]byte, 1048577)); err == nil || !maps.Equal(before, index) {
		t.Fatal("oversized stage mutated index")
	}
}
func TestHiddenSupplyCommitValidationPreservesStore(t *testing.T) {
	root := t.TempDir()
	id, err := WriteVault(root, "cargo", []byte("reserve"))
	if err != nil {
		t.Fatal(err)
	}
	tree, err := WriteVault(root, "tree", []byte("tree"))
	if err != nil {
		t.Fatal(err)
	}
	invalidDoc, err := WriteVault(root, "snapshot", []byte("{}"))
	if err != nil {
		t.Fatal(err)
	}
	invalidJSON, err := WriteVault(root, "snapshot", []byte("{"))
	if err != nil {
		t.Fatal(err)
	}
	badMeta, _ := json.Marshal(SupplyCommit{Tree: tree, Tick: -1, Message: "bad"})
	badParent, err := WriteVault(root, "snapshot", badMeta)
	if err != nil {
		t.Fatal(err)
	}
	wrongTree, _ := json.Marshal(SupplyCommit{Tree: id, Message: "wrong kind"})
	wrongParent, err := WriteVault(root, "snapshot", wrongTree)
	if err != nil {
		t.Fatal(err)
	}
	tooMany := make(SupplyIndex)
	for i := 0; i < 129; i++ {
		tooMany[fmt.Sprintf("c%03d", i)] = id
	}
	for _, tc := range []struct {
		index  SupplyIndex
		parent string
		stamp  SupplyStamp
	}{
		{nil, "", SupplyStamp{Message: "empty"}}, {tooMany, "", SupplyStamp{Message: "large"}},
		{SupplyIndex{"../x": id}, "", SupplyStamp{Message: "path"}},
		{SupplyIndex{"a": id, "a/b": id}, "", SupplyStamp{Message: "conflict"}},
		{SupplyIndex{"cell": tree}, "", SupplyStamp{Message: "kind"}},
		{SupplyIndex{"cell": strings.Repeat("0", 40)}, "", SupplyStamp{Message: "missing"}},
		{SupplyIndex{"cell": "../bad"}, "", SupplyStamp{Message: "bad ID"}},
		{SupplyIndex{"cell": id}, id, SupplyStamp{Message: "parent kind"}},
		{SupplyIndex{"cell": id}, strings.Repeat("0", 40), SupplyStamp{Message: "parent missing"}},
		{SupplyIndex{"cell": id}, invalidDoc, SupplyStamp{Message: "schema"}},
		{SupplyIndex{"cell": id}, invalidJSON, SupplyStamp{Message: "json"}},
		{SupplyIndex{"cell": id}, badParent, SupplyStamp{Message: "metadata"}},
		{SupplyIndex{"cell": id}, wrongParent, SupplyStamp{Message: "tree kind"}},
		{SupplyIndex{"cell": id}, "", SupplyStamp{Tick: -1, Message: "negative"}},
		{SupplyIndex{"cell": id}, "", SupplyStamp{Message: " \t\n"}},
		{SupplyIndex{"cell": id}, "", SupplyStamp{Message: strings.Repeat("x", 513)}},
	} {
		before := maps.Clone(tc.index)
		files := hiddenSupplyFiles(t, root)
		got, err := CommitSupply(root, tc.index, tc.parent, tc.stamp)
		if err == nil || got != "" || !maps.Equal(before, tc.index) || !reflect.DeepEqual(files, hiddenSupplyFiles(t, root)) {
			t.Fatal("commit validation wrote state", tc.parent, got, err)
		}
	}
	if _, err := CommitSupply(root, SupplyIndex{"cell": id}, "", SupplyStamp{Message: strings.Repeat("x", 512)}); err != nil {
		t.Fatal("message boundary", err)
	}
	// A valid tree and a pre-existing corrupt snapshot force failure after tree publication.
	index := SupplyIndex{"new": id}
	stamp := SupplyStamp{Tick: 991, Message: "publish failure"}
	rawTree := hiddenSupplyTree(index)
	treeID := hiddenSupplyID("tree", rawTree)
	document, _ := json.Marshal(SupplyCommit{Tree: treeID, Tick: stamp.Tick, Message: stamp.Message})
	snapshotID := hiddenSupplyID("snapshot", document)
	file := filepath.Join(root, "vault-objects", snapshotID+".z")
	if err := os.WriteFile(file, []byte("corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	before := maps.Clone(index)
	if got, err := CommitSupply(root, index, "", stamp); err == nil || got != "" || !maps.Equal(before, index) {
		t.Fatal("snapshot write failure")
	}
	after, _ := os.ReadFile(file)
	if string(after) != "corrupt" {
		t.Fatal("corrupt snapshot overwritten")
	}
	if _, _, err := ReadVault(root, treeID); err != nil {
		t.Fatal("successful tree should remain", err)
	}
	// Missing cargo preserves the original filesystem error cause.
	if _, err := CommitSupply(root, SupplyIndex{"cell": strings.Repeat("1", 40)}, "", SupplyStamp{Message: "missing"}); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("lost cause", err)
	}
}
`
