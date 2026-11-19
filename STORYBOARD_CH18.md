# Chapter 18 lesson storyboards

Short bilingual lessons for six labs. Examples use illustrative values; the code does not implement a lab contract or provide tests.

## 18-1 — Mission Log

Interactive sandbox (`ch18-slice-index`): move an index and watch the selected illustrative value.

### Idea · `e18-1-table`

<!-- go-example id=e18-1-table imports=fmt outcome=ok -->
```go
func main() {
	cases := []struct{ name string; value int }{
		{"small", 2},
		{"large", 4},
	}
	for _, tc := range cases {
		fmt.Println(tc.name, tc.value)
	}
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | table starts | — | ابدأ بسجل حالات. | Start a case table. |
| 2 / 5 | two cases | — | الحالتان جاهزتان. | Both cases are ready. |
| 3 / 6 | tc=small | — | اقرأ الحالة الأولى. | Read the first case. |
| 4 / 7 | tc=small | `small 2` | اعرض قيمتها. | Print its values. |
| 5 / 7 | tc=large | `small 2` + `large 4` | نفس الفحص للحالة التالية. | Run the same check again. |
| 6 / 8 | loop ends | cumulative | خلصت كل الحالات. | Every case has run. |

#### Verified e18-1-table — ok

Real stdout:

```text
small 2
large 4
```

### Try it · `e18-1-index`

<!-- go-example id=e18-1-index imports=fmt outcome=ok -->
```go
func main() {
	values := []int{7, 4}
	fmt.Println(values[1])
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | values=[7 4] | — | قيمتان بترتيبهم. | Two ordered values. |
| 2 / 3 | index 1 | `4` | اقرأ index واحد. | Read index one. |

#### Verified e18-1-index — ok

Real stdout:

```text
4
```

## 18-2 — Black Box

Interactive sandbox (`ch18-json-marshal`): encode an illustrative Go value and inspect JSON bytes.

### Idea · `e18-2-json`

<!-- go-example id=e18-2-json imports=encoding/json,fmt outcome=ok -->
```go
type Beacon struct {
	Charge int `json:"charge"`
}
func main() {
	b := Beacon{Charge: 7}
	raw, err := json.Marshal(b)
	fmt.Println(string(raw), err)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 1 | Beacon type | — | نوع ببيانات مصدّرة. | A type with an exported field. |
| 2 / 5 | b={7} | — | القيمة جاهزة للتحويل. | The value is ready. |
| 3 / 6 | raw JSON bytes | — | حوّلها لصيغة JSON. | Convert it to JSON. |
| 4 / 7 | output | `{"charge":7} <nil>` | اعرض البايتات والناتج. | Show bytes and error. |

#### Verified e18-2-json — ok

Real stdout:

```text
{"charge":7} <nil>
```

### Try it · `e18-2-map`

<!-- go-example id=e18-2-map imports=encoding/json,fmt outcome=ok -->
```go
func main() {
	raw, err := json.Marshal(map[string]int{"hull": 3})
	fmt.Println(string(raw), err)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | raw bytes | — | Marshal ترجع بايتات. | Marshal returns bytes. |
| 2 / 3 | printed object | `{"hull":3} <nil>` | الناتج object وخطأ nil. | Object output and nil error. |

#### Verified e18-2-map — ok

Real stdout:

```text
{"hull":3} <nil>
```

## 18-3 — Cargo Manifest

Interactive sandbox (`ch18-csv-quoted-field`): switch rows and see a quoted comma stay inside one field.

### Idea · `e18-3-csv`

<!-- go-example id=e18-3-csv imports=encoding/csv,fmt,strings outcome=ok -->
```go
func main() {
	r := csv.NewReader(strings.NewReader("north,4\n"))
	row, err := r.Read()
	fmt.Println(row, err)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | reader ready | — | أنشئ reader من نص. | Create a reader from text. |
| 2 / 3 | row decoded | — | اقرأ record واحدة. | Read one record. |
| 3 / 4 | fields | `[north 4] <nil>` | اطبع الحقول والخطأ. | Print fields and error. |

#### Verified e18-3-csv — ok

Real stdout:

```text
[north 4] <nil>
```

### Try it · `e18-3-quoted`

<!-- go-example id=e18-3-quoted imports=encoding/csv,fmt,strings outcome=ok -->
```go
func main() {
	r := csv.NewReader(strings.NewReader("\"west, upper\",4\n"))
	row, err := r.Read()
	fmt.Println(row, err)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | CSV input | — | الـreader يستقبل نص CSV. | The reader receives CSV text. |
| 2 / 3 | row decoded | — | Read تفك الحقول. | Read decodes the fields. |
| 3 / 4 | one quoted field | `[west, upper 4] <nil>` | الاسم يفضل حقل واحد. | The name stays one field. |

#### Verified e18-3-quoted — ok

Real stdout:

```text
[west, upper 4] <nil>
```

## 18-4 — Signal Words

Interactive sandbox (`ch18-unicode-words`): split illustrative text and toggle lowercase normalization.

### Idea · `e18-4-unicode`

<!-- go-example id=e18-4-unicode imports=fmt,strings,unicode outcome=ok -->
```go
func main() {
	words := strings.FieldsFunc("Echo—Relay 7", func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	fmt.Println(strings.ToLower(words[0]), len(words))
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | words parsed | — | جزّئ حسب حروف Unicode. | Split using Unicode letters. |
| 2 / 3 | normalized count | `echo 3` | حوّل أول كلمة واعرض العدد. | Lowercase and print the count. |

#### Verified e18-4-unicode — ok

Real stdout:

```text
echo 3
```

### Try it · `e18-4-lower`

<!-- go-example id=e18-4-lower imports=fmt,strings outcome=ok -->
```go
func main() {
	word := strings.ToLower("Echo")
	fmt.Println(word)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | lowercase word | — | حوّل الكلمة لحروف صغيرة. | Convert the word to lowercase. |
| 2 / 3 | output | `echo` | اعرض القيمة الجديدة. | Print the new value. |

#### Verified e18-4-lower — ok

Real stdout:

```text
echo
```

## 18-5 — Flight Recorder

Interactive sandbox (`ch18-splitn-limit`): raise the split limit and inspect the retained final part.

### Idea · `e18-5-splitn`

<!-- go-example id=e18-5-splitn imports=fmt,strings outcome=ok -->
```go
func main() {
	parts := strings.SplitN("12|engine|battery|warm", "|", 3)
	fmt.Println(len(parts), parts[2])
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | at most 3 parts | — | افصل لحد تلات أجزاء. | Split into at most three parts. |
| 2 / 3 | final part retained | `3 battery|warm` | اطبع العدد والرسالة. | Print count and message. |

#### Verified e18-5-splitn — ok

Real stdout:

```text
3 battery|warm
```

### Try it · `e18-5-limit`

<!-- go-example id=e18-5-limit imports=fmt,strings outcome=ok -->
```go
func main() {
	parts := strings.SplitN("9|relay|warm|again", "|", 3)
	fmt.Println(len(parts), parts[2])
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | two splits | — | اقسم على أول فاصلين. | Split on the first two separators. |
| 2 / 3 | final message | `3 warm|again` | اعرض العدد والجزء الأخير. | Print count and final piece. |

#### Verified e18-5-limit — ok

Real stdout:

```text
3 warm|again
```

## 18-6 — Crew Directory

Interactive sandbox (`ch18-normalized-search`): normalize an illustrative name before `Contains`.

### Idea · `e18-6-interface`

<!-- go-example id=e18-6-interface imports=bytes,fmt,io outcome=ok -->
```go
func emit(w io.Writer, text string) error {
	_, err := fmt.Fprint(w, text)
	return err
}
func main() {
	var b bytes.Buffer
	err := emit(&b, "ready")
	fmt.Println(b.String(), err)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 1 | function takes io.Writer | — | الدالة تستقبل interface. | The function accepts an interface. |
| 2 / 6 | Buffer value | — | Buffer يطبّق io.Writer. | Buffer implements io.Writer. |
| 3 / 7 | write call | — | اكتب من خلال الحد المشترك. | Write through the shared boundary. |
| 4 / 8 | result | `ready <nil>` | اعرض النص وخطأ الكتابة. | Print text and write error. |

#### Verified e18-6-interface — ok

Real stdout:

```text
ready <nil>
```

### Try it · `e18-6-search`

<!-- go-example id=e18-6-search imports=fmt,strings outcome=ok -->
```go
func main() {
	name := strings.ToLower("Pilot")
	fmt.Println(strings.Contains(name, "lot"))
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | normalized name | — | وحّد حالة الأحرف. | Normalize letter case. |
| 2 / 3 | match | `true` | ابحث واعرض النتيجة. | Search and print the result. |

#### Verified e18-6-search — ok

Real stdout:

```text
true
```
