# Chapter 8 lesson storyboard

Each mission uses the approved six-step sequence: Mission, Idea, Try it, Watch out, Think & check, Your turn. Arabic and English share IDs, block kinds, code, answers, and walkthrough states. Each Try step asks for a prediction before the following walkthrough reveals the same example. Examples use illustrative values and do not encode a mission solution.

## 8-1 — The Counter That Does Not Count

**Playground:** Change a local copy, then compare a write through a pointer.

**Word count:** Arabic 288; English 303. Counts include step titles, prose, captions, options, feedback, and three hints; excludes code.

**Six steps:** Mission compares collected crystal with the report; Idea traces a copied `int`; Try predicts pointer-based output before its trace; Watch out contrasts values and addresses; Think & check reviews parameter and call-site phases; Your turn returns to the starter function.

**Idea:** A local `int` parameter is a copy; follow it through a function call and compare the caller's value after return. The walkthrough ends at the actual output.

<!-- go-example id=e81-copy imports=fmt outcome=ok -->
```go
x := 3
adjust := func(n int) {
	n = 10
}
adjust(x)
fmt.Println(x)
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 1 | x يبدأ بالرقم 3. · x starts at 3. | |
| 2 | الـfunction اتجهزت لاستقبال الـvalue. · The function is ready for a value. | |
| 3 | n نسخة محلية من x. · n is a local copy of x. | |
| 4 | التغيير يخص النسخة. · Only the copy changes. | |
| 5 | الـfunction رجعت؛ x كما هي. · The function returns; x stays. | |
| 6 | الخرج: 3. · Output: 3. | `3` |

<!-- go-example id=e81-pointer imports=fmt outcome=ok -->
```go
x := 4
p := &x
fmt.Println(*p)
*p = 9
fmt.Println(x)
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 1 | x يبدأ بالرقم 4. · x starts at 4. | |
| 2 | p يحمل عنوان x. · p holds x's address. | |
| 3 | قراءة العنوان تعرض 4. · Reading the address shows 4. | `4` |
| 4 | الكتابة عبر p تغيّر x. · Writing through p changes x. | `4` |
| 5 | الخرج النهائي: 4 ثم 9. · Final output: 4 then 9. | `4`, `9` |

**Try:** Predict the two prints when the value is read and then changed through a pointer; inspect the trace afterward.

## 8-2 — Cargo Runs

**Playground:** Run value and pointer receivers against the same small value.

**Word count:** Arabic 295; English 301. Counts include step titles, prose, captions, options, feedback, and three hints; excludes code.

**Six steps:** Mission states the 10-unit hold and storage goal; Idea traces a value receiver; Try predicts a pointer-receiver result before its trace; Watch out notes unload location and fit; Think & check names cargo, space, and unload phases; Your turn points back to `Load`.

**Idea:** A value receiver mutates its copy. The walkthrough follows a small `Gauge` through the method call and shows the caller still has 2.

<!-- go-example id=e82-value imports=fmt outcome=ok -->
```go
type Gauge struct{ Value int }
func (g Gauge) Add(n int) {
	g.Value += n
}
func main() {
	g := Gauge{Value: 2}
	g.Add(3)
	fmt.Println(g.Value)
}
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 6 | Gauge بتبدأ بالـvalue 2. · Gauge starts at 2. | |
| 7 | استدعاء Add يعمل نسخة receiver. · Add receives a receiver copy. | |
| 3 | التعديل يقع على النسخة. · The copy is changed. | |
| 4 | نسخة receiver انتهت. · The receiver copy ends. | |
| 8 | الخرج يفضل 2. · Output stays 2. | `2` |

<!-- go-example id=e82-pointer imports=fmt outcome=ok -->
```go
type Gauge struct{ Value int }
func (g *Gauge) Add(n int) {
	g.Value += n
}
func main() {
	g := Gauge{Value: 2}
	g.Add(3)
	fmt.Println(g.Value)
}
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 6 | g بتبدأ بالـvalue 2. · g starts at 2. | |
| 7 | النداء يمرّر عنوان g. · The call passes g's address. | |
| 2 | الـreceiver يوصل للأصل. · The receiver reaches the original. | |
| 3 | الحقل الأصلي يزيد. · The original field increases. | |
| 8 | الخرج: 5. · Output: 5. | `5` |

**Try:** Predict the field after a pointer-receiver method call, then watch the receiver reach the original struct.

## 8-3 — Radar

**Playground:** Toggle a target and watch the nil guard choose a safe path.

**Word count:** Arabic 289; English 271. Counts include step titles, prose, captions, options, feedback, and three hints; excludes code.

**Six steps:** Mission asks for every target; Idea traces a nil guard; Try predicts a valid-pointer result before its trace; Watch out states the radar's pointer and tick behavior; Think & check names scan, check, use, repeat; Your turn tests the starter with and without targets.

**Idea:** A pointer may be `nil`; check before dereferencing. This short guard ends safely when there is no value.

<!-- go-example id=e83-nil imports=fmt outcome=ok -->
```go
var p *int
if p == nil {
	fmt.Println("none")
	return
}
fmt.Println(*p)
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 1 | p يبدأ nil. · p starts as nil. | |
| 2 | الشرط يكتشف nil. · The check detects nil. | |
| 3 | نطبع إن مفيش هدف. · Print that no target exists. | `none` |
| 4 | return ينهي المثال بأمان. · return ends safely. | `none` |

<!-- go-example id=e83-value imports=fmt outcome=ok -->
```go
n := 6
p := &n
if p == nil {
	fmt.Println("none")
} else {
	fmt.Println(*p)
}
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 1 | n يساوي 6. · n is 6. | |
| 2 | p يشاور على n. · p points to n. | |
| 3 | الـpointer مش nil. · The pointer is not nil. | |
| 6 | اقرا الـvalue 6. · Read the value 6. | `6` |

**Try:** Predict the guarded result when the pointer is valid. Mission notes identify `colony.Nearest()` as returning `*colony.Point`, or `nil`, and each radar query costs one tick.

## 8-4 — Cart Chain

**Playground:** Follow a temporary pointer while the links stay fixed.

**Word count:** Arabic 296; English 304. Counts include step titles, prose, captions, options, feedback, and three hints; excludes code.

**Six steps:** Mission introduces preserving the convoy; Idea traces one `Next` link; Try predicts tags while a temporary pointer advances; Watch out separates moving a pointer from changing a link; Think & check names empty, link, release, search, snapshot phases; Your turn tests each method.

**Idea:** Nodes can point to another node of the same type. Follow one `Next` link and read the tag on the next lamp.

<!-- go-example id=e84-link imports=fmt outcome=ok -->
```go
type Lamp struct{ Tag string; Next *Lamp }
last := &Lamp{Tag: "blue"}
first := &Lamp{Tag: "red", Next: last}
fmt.Println(first.Next.Tag)
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 2 | آخر لمبة اسمها blue. · The last lamp is blue. | |
| 3 | first تتصل بآخر لمبة. · first links to the last lamp. | |
| 4 | اقرا tag بعد الوصلة. · Read the tag beyond the link. | `blue` |

<!-- go-example id=e84-walk imports=fmt outcome=ok -->
```go
type Lamp struct{ Tag string; Next *Lamp }
tail := &Lamp{Tag: "gold"}
middle := &Lamp{Tag: "green", Next: tail}
head := &Lamp{Tag: "red", Next: middle}
p := head
fmt.Println(p.Tag)
p = p.Next
fmt.Println(p.Tag)
```

| Line | State / caption (AR · EN) | Output |
|---:|---|---|
| 0 | قبل التنفيذ. · Before execution. | |
| 2 | tail هو gold. · tail is gold. | |
| 3 | middle يشاور على tail. · middle points to tail. | |
| 4 | head يشاور على الوسط. · head points to middle. | |
| 5 | p يبدأ عند head. · p starts at head. | |
| 6 | أول tag هو red. · The first tag is red. | `red` |
| 7 | p يتبع Next للوسط. · p follows Next to middle. | `red` |
| 8 | tag التالي هو green. · The next tag is green. | `red`, `green` |

**Try:** Predict which tags appear after moving a temporary pointer once. The chain itself does not change.

## Verified Go examples

Each example below is built and run by `ch08_lessons_test.go`; stdout must match exactly.

#### Verified `e81-copy`

Real stdout:

```text
3
```

#### Verified `e81-pointer`

Real stdout:

```text
4
9
```

#### Verified `e82-value`

Real stdout:

```text
2
```

#### Verified `e82-pointer`

Real stdout:

```text
5
```

#### Verified `e83-nil`

Real stdout:

```text
none
```

#### Verified `e83-value`

Real stdout:

```text
6
```

#### Verified `e84-link`

Real stdout:

```text
blue
```

#### Verified `e84-walk`

Real stdout:

```text
red
green
```
