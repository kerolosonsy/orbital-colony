# Chapter 2 lesson storyboards

Short bilingual lessons for the four chapter 2 missions. Mission setups, starters, checks, and fallback markdown stay unchanged. All examples use illustrative values, real Go APIs, and verified stdout.

## 2-1 — Inventory Check (`2-1-inventory`)

The learner watches values enter variables, sees addition, and distinguishes Console output from a mission report.

<!-- go-example id=e21-variables imports=fmt outcome=ok -->
```go
func main() {
	moss := 3
	crystal := 2
	total := moss + crystal
	fmt.Println(total)
}
```

#### Verified e21-variables — ok

Real stdout:

```text
5
```

<!-- go-example id=e21-total imports=fmt outcome=ok -->
```go
func main() {
	parts := 4
	extra := 2
	fmt.Println(parts + extra)
}
```

#### Verified e21-total — ok

Real stdout:

```text
6
```

**Playground idea:** Show two changing storage counters feeding named value cards and a total card. Let the learner step through the read and addition before choosing Console or Base output.

## 2-2 — Survey the Area (`2-2-survey`)

The learner observes multiple return values and keeps the two values in their returned order.

<!-- go-example id=e22-values imports=fmt outcome=ok -->
```go
func dimensions() (int, int) {
	return 4, 2
}
func main() {
	width, height := dimensions()
	fmt.Println(width, height)
}
```

#### Verified e22-values — ok

Real stdout:

```text
4 2
```

<!-- go-example id=e22-pair imports=fmt outcome=ok -->
```go
func pair() (int, int) { return 2, 3 }
func main() {
	a, b := pair()
	fmt.Println(a, b)
}
```

#### Verified e22-pair — ok

Real stdout:

```text
2 3
```

**Playground idea:** Let learners resize a small map and watch width and height populate separate slots. A second panel accepts two reports and marks their arrival order.

## 2-3 — Callsigns (`2-3-callsigns`)

The learner compares UTF-8 byte length with rune count and sees first and last rune decoding on sample text.

<!-- go-example id=e23-utf8 imports=fmt,unicode/utf8 outcome=ok -->
```go
func main() {
	text := "é🚀"
	first, _ := utf8.DecodeRuneInString(text)
	last, _ := utf8.DecodeLastRuneInString(text)
	fmt.Println(len(text), utf8.RuneCountInString(text), string(first), string(last))
}
```

#### Verified e23-utf8 — ok

Real stdout:

```text
6 2 é 🚀
```

<!-- go-example id=e23-count imports=fmt,unicode/utf8 outcome=ok -->
```go
func main() {
	text := "é🚀"
	fmt.Println(utf8.RuneCountInString(text))
}
```

#### Verified e23-count — ok

Real stdout:

```text
2
```

**Playground idea:** Render sample text as UTF-8 byte cells and rune groups. Selecting a rune highlights its byte span and updates the first/last markers.

## 2-4 — Cargo Units (`2-4-units`)

The learner sees named constants created with `iota` and explicitly converted for printing.

<!-- go-example id=e24-iota imports=fmt outcome=ok -->
```go
type Role int
const (
	Guest Role = iota
	Pilot
	Chief
)
func main() {
	fmt.Println(int(Guest), int(Pilot), int(Chief))
}
```

#### Verified e24-iota — ok

Real stdout:

```text
0 1 2
```

<!-- go-example id=e24-badge imports=fmt outcome=ok -->
```go
type Badge int
const (
	New Badge = iota
	Ready
	Sent
)
func main() { fmt.Println(int(Ready)) }
```

#### Verified e24-badge — ok

Real stdout:

```text
1
```

**Playground idea:** Display three differently sized cargo containers with their counts, then let learners choose shared-unit conversion cards and watch units align before a total appears.
