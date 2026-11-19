# Chapter 30 lesson storyboards

Short bilingual lessons for the four storage labs and the nursery boss. Values are illustrative. The examples teach small Go concepts; they do not implement a lab contract.

## 30-1 — Cache Aside

Mission: keep the store authoritative while Redis speeds reads. Idea: distinguish a cache hit from a missing key. Try: predict the zero value and presence flag. Watch: separate confirmed absence from Redis outages. Think: trace the read and write routes. Your turn: inspect one operation at a time.

<!-- go-example id=e301-cache-hit imports=fmt outcome=ok -->
```go
func main() {
	cache := map[string]string{"ridge": "ready"}
	value, ok := cache["ridge"]
	fmt.Printf("%q %t\n", value, ok)
}
```

<!-- go-example id=e301-cache-miss imports=fmt outcome=ok -->
```go
func main() {
	cache := map[string]string{"ridge": "ready"}
	value, ok := cache["cove"]
	fmt.Printf("%q %t\n", value, ok)
}
```

## 30-2 — Documents

Mission: model stored records and deterministic pages. Idea: a struct groups typed fields; a unique index protects shared storage. Try: predict a small page with one extra row. Watch: close and check cursors. Think: separate the archive interface from its Mongo adapter. Your turn: start with the storage boundary.

<!-- go-example id=e302-document imports=fmt outcome=ok -->
```go
type Card struct {
	ID  int64
	Tag string
}
func main() {
	c := Card{ID: 4, Tag: "basalt"}
	fmt.Println(c.Tag, c.ID)
}
```

<!-- go-example id=e302-page-extra imports=fmt outcome=ok -->
```go
func main() {
	rows := []string{"mica", "slate", "quartz"}
	limit := 2
	fmt.Println(rows[:limit], len(rows) > limit)
}
```

## 30-3 — Search Index

Mission: find dispatches through an inverted index, then rank and page them. Idea: terms map to posting lists. Try: predict lexical order. Watch: map iteration is unspecified. Think: normalize, gather, rank, and slice. Your turn: build the index before its query path.

<!-- go-example id=e303-postings imports=fmt outcome=ok -->
```go
func main() {
	postings := map[string][]string{
		"amber": {"note-q", "note-m"},
	}
	fmt.Println(postings["amber"])
}
```

<!-- go-example id=e303-lexical-order imports=fmt,sort outcome=ok -->
```go
func main() {
	codes := []string{"note-q", "note-m"}
	sort.Strings(codes)
	fmt.Println(codes)
}
```

## 30-4 — Read and Write

Mission: write once to primary and route reads by session. Idea: an injected clock defines a temporary sticky window. Try: predict whether a wrapped cause is found. Watch: domain errors do not mean replica outage. Think: trace each route and session separately. Your turn: mark the route before each call.

<!-- go-example id=e304-sticky-window imports=fmt,time outcome=ok -->
```go
func main() {
	start := time.Unix(120, 0)
	until := start.Add(time.Minute)
	fmt.Println(start.Before(until))
}
```

<!-- go-example id=e304-wrapped-cause imports=errors,fmt outcome=ok -->
```go
func main() {
	var ErrOffline = errors.New("offline")
	err := fmt.Errorf("replica: %w", ErrOffline)
	fmt.Println(errors.Is(err, ErrOffline))
}
```

## 30-Boss — Planting Season

Mission: reserve supplies transactionally and invalidate cache after commit. Idea: keep proposed state distinct from saved state. Try: predict which causes remain discoverable. Watch: cache failure after commit cannot undo SQL. Think: order validation, transaction, commit, and cache work. Your turn: draw the commit boundary first.

<!-- go-example id=e30b-draft-state imports=fmt outcome=ok -->
```go
type Sketch struct {
	Count int
}
func main() {
	saved := Sketch{Count: 7}
	proposal := Sketch{Count: saved.Count + 2}
	fmt.Println(saved.Count, proposal.Count)
}
```

<!-- go-example id=e30b-joined-errors imports=errors,fmt outcome=ok -->
```go
func main() {
	sentinel := errors.New("sync warning")
	cause := errors.New("cache unavailable")
	err := errors.Join(sentinel, cause)
	fmt.Println(errors.Is(err, sentinel), errors.Is(err, cause))
}
```

## Interactive playgrounds

- 30-1: Read through a cache miss, fill the cache, evict the key, then read again.
- 30-2: Move across small document pages and watch the extra-row flag.
- 30-3: Look up terms and inspect their posting lists.
- 30-4: Write for one session, advance an illustrative clock, and compare read routes.
- 30-Boss: Commit proposed state before simulating a cache failure.

## Toolchain transcripts

#### Verified e301-cache-hit — ok

Real stdout:

```text
"ready" true
```

#### Verified e301-cache-miss — ok

Real stdout:

```text
"" false
```

#### Verified e302-document — ok

Real stdout:

```text
basalt 4
```

#### Verified e302-page-extra — ok

Real stdout:

```text
[mica slate] true
```

#### Verified e303-postings — ok

Real stdout:

```text
[note-q note-m]
```

#### Verified e303-lexical-order — ok

Real stdout:

```text
[note-m note-q]
```

#### Verified e304-sticky-window — ok

Real stdout:

```text
true
```

#### Verified e304-wrapped-cause — ok

Real stdout:

```text
true
```

#### Verified e30b-draft-state — ok

Real stdout:

```text
7 9
```

#### Verified e30b-joined-errors — ok

Real stdout:

```text
true true
```
