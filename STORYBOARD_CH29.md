# Chapter 29 lesson storyboard

Short bilingual practice for the six database labs. Each lesson has six steps; practice predicts output before revealing a line-by-line trace. Data is illustrative, and no lab implementation or test case is supplied.

## 29-1 — SQLite CRUD

AR: افصل قيم SQL عن نص الاستعلام، وتعامل مع غياب الصف كخطأ واضح. EN: Bind SQL values separately and treat a missing row as a clear error.

<!-- go-example id=e291-error-chain imports=database/sql,errors,fmt outcome=ok -->
```go
func main() {
	err := fmt.Errorf("lookup: %w", sql.ErrNoRows)
	fmt.Println(errors.Is(err, sql.ErrNoRows))
}
```

<!-- go-example id=e291-named-arg imports=database/sql,fmt outcome=ok -->
```go
func main() {
	arg := sql.Named("label", "Luna")
	fmt.Println(arg.Name)
}
```

## 29-2 — Migrations

AR: خليك واعي بحدود كل ترحيل؛ المخطط ورقم النسخة بيتحفظوا سوا. EN: Track each migration boundary; schema and version commit together.

<!-- go-example id=e292-staged-list imports=fmt,slices outcome=ok -->
```go
func main() {
	applied := []int{2}
	next := slices.Clone(applied)
	next = append(next, 5)
	fmt.Println(applied, next)
}
```

<!-- go-example id=e292-trim-script imports=fmt,strings outcome=ok -->
```go
func main() {
	blank := strings.TrimSpace("   ") == ""
	fmt.Println(blank)
}
```

## 29-3 — Atomic transfer

AR: حافظ على مجموع الأرصدة، وخلي الخصم والإضافة وتسجيل الطلب عملية واحدة. EN: Preserve the balance total; debit, credit, and request recording form one operation.

<!-- go-example id=e293-conservation imports=fmt outcome=ok -->
```go
func main() {
	from, to, moved := int64(9), int64(2), int64(3)
	fmt.Println(from+to, from-moved+to+moved)
}
```

<!-- go-example id=e293-balance-check imports=fmt outcome=ok -->
```go
func main() {
	balance, requested := int64(5), int64(2)
	fmt.Println(balance >= requested, balance-requested)
}
```

## 29-4 — Connection pool

AR: فرّق بين حد الحوض وحالته الحالية، واقفل الصفوف عند كل رجوع. EN: Distinguish a pool limit from its current state, and close rows on every return.

<!-- go-example id=e294-pool-stats imports=database/sql,fmt outcome=ok -->
```go
func main() {
	stats := sql.DBStats{MaxOpenConnections: 3, InUse: 2, Idle: 1}
	fmt.Println(stats.MaxOpenConnections, stats.InUse, stats.Idle)
}
```

<!-- go-example id=e294-stats-fields imports=database/sql,fmt outcome=ok -->
```go
func main() {
	stats := sql.DBStats{MaxOpenConnections: 4, Idle: 2}
	fmt.Println(stats.MaxOpenConnections, stats.Idle)
}
```

## 29-5 — sqlx repository

AR: الواجهة بتفصل المجال عن قاعدة البيانات، والـnamed argument بيفصل الاسم عن القيمة. EN: An interface separates domain logic from storage; a named argument separates name from value.

<!-- go-example id=e295-interface-store imports=fmt outcome=ok -->
```go
type Store interface {
	Put(string, int) error
}
type memoryStore struct{}
func (memoryStore) Put(bay string, units int) error {
	fmt.Println(bay, units)
	return nil
}
func main() {
	var store Store = memoryStore{}
	_ = store.Put("east", 3)
}
```

<!-- go-example id=e295-named-arg imports=database/sql,fmt outcome=ok -->
```go
func main() {
	arg := sql.Named("bay", "north")
	fmt.Println(arg.Name)
}
```

## 29-6 — GORM associations

AR: الـassociation بتجمع model أب وslice أبناء؛ فكّر في التحميل والترتيب كجزء من العقد. EN: An association groups a parent model and child slice; loading and order are part of the contract.

<!-- go-example id=e296-association imports=fmt outcome=ok -->
```go
type Bed struct{ Label string }
type Dome struct {
	Name string
	Beds []Bed
}
func main() {
	dome := Dome{Name: "Vega", Beds: []Bed{{Label: "north"}, {Label: "south"}}}
	fmt.Println(dome.Name, len(dome.Beds))
}
```

<!-- go-example id=e296-association-index imports=fmt outcome=ok -->
```go
type Bed struct{ Label string }
type Dome struct{ Beds []Bed }
func main() {
	dome := Dome{Beds: []Bed{{Label: "inner"}, {Label: "outer"}}}
	fmt.Println(dome.Beds[1].Label)
}
```

## Interactive playgrounds

- 29-1: Switch a bound key and compare a found row with an absent one.
- 29-2: Stage a schema change, then commit it with its version.
- 29-3: Replay two concurrent transfers in alternate valid orders; the order can vary.
- 29-4: Acquire and release a connection while the configured cap stays fixed.
- 29-5: Switch a storage implementation behind one interface call.
- 29-6: Load and hide the child slice associated with a parent.

## Toolchain transcripts

Each example is compiled and run as a standalone Go program by `TestCh29StoryboardExamples`.

#### Verified e291-error-chain — ok

Real stdout:

```text
true
```

#### Verified e291-named-arg — ok

Real stdout:

```text
label
```

#### Verified e292-staged-list — ok

Real stdout:

```text
[2] [2 5]
```

#### Verified e292-trim-script — ok

Real stdout:

```text
true
```

#### Verified e293-conservation — ok

Real stdout:

```text
11 11
```

#### Verified e293-balance-check — ok

Real stdout:

```text
true 3
```

#### Verified e294-pool-stats — ok

Real stdout:

```text
3 2 1
```

#### Verified e294-stats-fields — ok

Real stdout:

```text
4 2
```

#### Verified e295-interface-store — ok

Real stdout:

```text
east 3
```

#### Verified e295-named-arg — ok

Real stdout:

```text
bay
```

#### Verified e296-association — ok

Real stdout:

```text
Vega 2
```

#### Verified e296-association-index — ok

Real stdout:

```text
outer
```
