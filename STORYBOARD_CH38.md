# Chapter 38 lesson storyboards

Short bilingual lesson outlines. Each Try it question predicts the same fragment revealed by its walkthrough. Examples use illustrative values.

## 38-1 — Paths Without Storage / مسارات من غير تخزين

Mission: build iterators that respect the consumer's stop request. Idea: consume a standard-library Seq with `range`. Try: predict which values remain after `break`. Watch out: constructing a Seq does not consume it. Think: separate producer and consumer behavior. Your turn: implement and check the contracts in small pieces.

<!-- go-example id=e381-values imports=fmt,slices outcome=ok -->
```go
func main() {
	values := []int{2, 5}
	for value := range slices.Values(values) {
		fmt.Println(value)
	}
}
```

<!-- go-example id=e381-break imports=fmt,slices outcome=ok -->
```go
func main() {
	values := []int{3, 6, 9}
	for n := range slices.Values(values) {
		fmt.Println(n)
		if n == 6 { break }
	}
}
```

## 38-2 — Distinct Gates / بوابات مميّزة

Mission: choose from ordered rows while tracking membership. Idea: use a map as a set. Try: predict which values a labeled `continue` skips. Watch out: maps do not provide result order. Think: distinguish item scope from row scope. Your turn: trace membership and labels on small rows.

<!-- go-example id=e382-set imports=fmt outcome=ok -->
```go
func main() {
	labels := map[string]struct{}{}
	labels["north"] = struct{}{}
	_, exists := labels["north"]
	fmt.Println(exists, len(labels))
}
```

<!-- go-example id=e382-label imports=fmt outcome=ok -->
```go
func main() {
Decks:
	for _, deck := range [][]int{{1, 0}, {2, 3}} {
		for _, slot := range deck {
			if slot == 0 { continue Decks }
			fmt.Println(slot)
		}
	}
}
```

## 38-3 — Start Once / تشغيل مرة واحدة

Mission: share one initialization result, including a zero value. Idea: `sync.Once` tracks execution per instance. Try: predict when a `OnceValue` loader runs. Watch out: zero is not an initialization flag. Think: keep state ownership clear. Your turn: test independent instances and deferred loading.

<!-- go-example id=e383-once imports=fmt,sync outcome=ok -->
```go
func main() {
	var once sync.Once
	calls := 0
	for i := 0; i < 2; i++ {
		once.Do(func() { calls++ })
	}
	fmt.Println(calls)
}
```

<!-- go-example id=e383-lazy imports=fmt,sync outcome=ok -->
```go
func main() {
	load := sync.OnceValue(func() string {
		fmt.Println("load")
		return "east"
	})
	fmt.Println("created")
	fmt.Println(load())
	fmt.Println(load())
}
```

## 38-4 — Bounded Crew / طاقم محدود

Mission: limit parallel work and handle errors. Idea: a group joins workers; cancellation is cooperative. The first walkthrough shows one possible worker order; another run may schedule differently. Try: predict the result after `Wait`. Watch out: the limit does not stop a running callback. Think: separate scheduling, cancellation and joining. Your turn: begin with small successful jobs.

<!-- go-example id=e384-group imports=fmt,sync/atomic,golang.org/x/sync/errgroup outcome=ok -->
```go
func main() {
	var total atomic.Int64
	var group errgroup.Group
	group.SetLimit(2)
	group.Go(func() error { total.Add(1); return nil })
	group.Go(func() error { total.Add(1); return nil })
	err := group.Wait()
	fmt.Println(err == nil, total.Load())
}
```

<!-- go-example id=e384-wait imports=fmt,golang.org/x/sync/errgroup outcome=ok -->
```go
func main() {
	var group errgroup.Group
	group.Go(func() error { return nil })
	err := group.Wait()
	fmt.Println(err == nil)
}
```

## 38-5 — Allocation Budget / ميزانية الذاكرة

Mission: preserve results while reducing allocations, then measure. Idea: `strings.Builder` collects text. Try: distinguish slice length from capacity. Watch out: allocation reports do not prove correctness. Think: separate benchmark setup from timed work. Your turn: establish behavior, add a sink, and compare measurements.

<!-- go-example id=e385-builder imports=fmt,strings outcome=ok -->
```go
func main() {
	var b strings.Builder
	b.Grow(8)
	b.WriteString("ion")
	b.WriteByte('-')
	b.WriteString("nav")
	fmt.Println(b.String())
}
```

<!-- go-example id=e385-capacity imports=fmt outcome=ok -->
```go
func main() {
	values := make([]int, 0, 4)
	values = append(values, 7)
	fmt.Println(len(values), cap(values))
}
```

## 38-6 — Preflight Inspection / تفتيش قبل الإقلاع

Mission: use `go vet` alongside tests to review shared state, formatting and cancellation. Idea: format verbs should match value types. Try: predict a formatted integer. Watch out: value receivers copy a struct, and every cancel path matters. Think: read diagnostics at their source, then check behavior. Your turn: rerun inspection and tests after each correction.

<!-- go-example id=e386-format imports=fmt outcome=ok -->
```go
func main() {
	count := -3
	fmt.Printf("count=%d\n", count)
}
```

<!-- go-example id=e386-printf imports=fmt outcome=ok -->
```go
func main() {
	fmt.Printf("items=%d\n", 12)
}
```

## 38-7 — Where Does the Time Go? / فين الوقت بيروح؟

Mission: collect CPU and heap profiles, then improve a distinct-value total. Idea: write heap profile data to an `io.Writer`. Try: use a set lookup to reason about membership. Watch out: stop a successful CPU profile and do not assume fixed output size. Think: separate the algorithm from profiling. Your turn: check writer failures and compare benchmarks without profile collection inside them.

<!-- go-example id=e387-heap imports=bytes,fmt,runtime/pprof outcome=ok -->
```go
func main() {
	var profile bytes.Buffer
	err := pprof.WriteHeapProfile(&profile)
	fmt.Println(err == nil, profile.Len() > 0)
}
```

<!-- go-example id=e387-set imports=fmt outcome=ok -->
```go
func main() {
	seen := map[int]struct{}{4: {}}
	_, exists := seen[4]
	fmt.Println(exists)
}
```

## Toolchain transcripts

Each example runs as a standalone program. The output below is captured from the Go toolchain.

#### Verified e381-values — ok

Real stdout:

```text
2
5
```

#### Verified e381-break — ok

Real stdout:

```text
3
6
```

#### Verified e382-set — ok

Real stdout:

```text
true 1
```

#### Verified e382-label — ok

Real stdout:

```text
1
2
3
```

#### Verified e383-once — ok

Real stdout:

```text
1
```

#### Verified e383-lazy — ok

Real stdout:

```text
created
load
east
east
```

#### Verified e384-group — ok

Real stdout:

```text
true 2
```

#### Verified e384-wait — ok

Real stdout:

```text
true
```

#### Verified e385-builder — ok

Real stdout:

```text
ion-nav
```

#### Verified e385-capacity — ok

Real stdout:

```text
1 4
```

#### Verified e386-format — ok

Real stdout:

```text
count=-3
```

#### Verified e386-printf — ok

Real stdout:

```text
items=12
```

#### Verified e387-heap — ok

Real stdout:

```text
true true
```

#### Verified e387-set — ok

Real stdout:

```text
true
```
