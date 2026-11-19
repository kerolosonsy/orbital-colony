# Chapter 9 lesson storyboards

Review draft for missions 9-1 through 9-4. Each mission has six steps in the approved 6-1 order, one concept walkthrough, and a Try step that asks for a prediction before revealing it. Text and frame captions are bilingual in the lesson data; code and execution state are identical. Examples use small illustrative values, not mission solutions.

Word counts use whitespace-separated learner text, including titles, captions, questions, feedback, and hints; code is excluded.

## 9-1 — تقرير بشكل محترم (`9-1-stringer`)

**Focus:** implicit interface satisfaction, `fmt.Stringer`, and formatted text. The starter already gathers the fields; the lesson leaves the formatting method for the player.

**Word count:** Arabic 252; English 313.

1. **المهمة / Mission:** identify the required report and the starter's existing `Status` value.
2. **الفكرة / Idea:** a method signature satisfies an interface; `Sprintf` builds text and `Sprint` uses `String`.
3. **جرّب / Try it:** predict a `time.Duration`'s printed form, then watch the interface value print.
4. **خد بالك / Watch out:** method spelling and return type matter.
5. **فكّر واختبر / Think & check:** identify the interface contract and the formatter that returns text.
6. **دورك / Your turn:** implement the behavior for `Status` and compare the report.

`c9m1.format-trace` · illustrative formatting

<!-- go-example id=e9m1-format imports=fmt outcome=ok -->
```go
label := fmt.Sprintf("dock-%02d", 4)
fmt.Println(label)
```

Frames: line 0 before; line 1 `label = dock-04`; line 2 output `dock-04`.

#### Verified `e9m1-format`

Real stdout:

```text
dock-04
```

`c9m1.stringer-predict` → `c9m1.stringer-trace` · predict, then reveal

<!-- go-example id=e9m1-stringer imports=fmt,time outcome=ok -->
```go
var label fmt.Stringer = 90 * time.Second
fmt.Println(label)
```

Frames: line 0 before; line 1 `label` holds a `time.Duration`; line 2 calls `String` and prints `1m30s`.

#### Verified `e9m1-stringer`

Real stdout:

```text
1m30s
```

**Hints, in order:** concept — implicit satisfaction; structure — put display behavior on the value and leave sending in `main`; specific — inspect `Report(s)` and the starter's `Status` fields.

## 9-2 — الطيار الآلي (`9-2-autopilot`)

**Focus:** `Autopilot` owns the map loop and calls a `Planner`; `Plan` returns ordered actions. The example uses the real standard-library callback `slices.IndexFunc`, not a colony plan.

**Word count:** Arabic 253; English 319.

1. **المهمة / Mission:** respond to tile ground and contents using the supplied `Farmer`.
2. **الفكرة / Idea:** connect the `Planner` signature to library-owned callback flow.
3. **جرّب / Try it:** predict how many callback calls happen before the first match.
4. **خد بالك / Watch out:** the autopilot sweeps once; wrong actions and moves cost ticks.
5. **فكّر واختبر / Think & check:** who calls `Plan`, and in what order are actions run?
6. **دورك / Your turn:** extend the starter's tile decisions and run the autopilot.

`c9m2.callback-trace` · callback example

<!-- go-example id=e9m2-callback imports=fmt,slices outcome=ok -->
```go
i := slices.IndexFunc([]int{2, 5, 8}, func(n int) bool {
	fmt.Println("check", n)
	return n > 4
})
fmt.Println("index", i)
```

Frames: line 0 before; line 1 search begins; lines 2–3 check 2 then 5; line 5 reports index 1.

#### Verified `e9m2-callback`

Real stdout:

```text
check 2
check 5
index 1
```

`c9m2.callback-predict` → `c9m2.callback-reveal` · predict, then reveal

<!-- go-example id=e9m2-predict imports=fmt,slices outcome=ok -->
```go
calls := 0
found := slices.IndexFunc([]int{3, 6}, func(n int) bool {
	calls++
	return n > 4
})
fmt.Println(found, calls)
```

Frames: line 1 starts at zero; lines 3–4 check 3 and then 6; line 6 prints index 1 and two calls.

#### Verified `e9m2-predict`

Real stdout:

```text
1 2
```

**Hints, in order:** concept — satisfy the `Planner` contract; structure — decide per tile and preserve action order; specific — inspect `Plan(t colony.Tile)` and its `Ground`, `Entity`, and `Ripe` fields.

## 9-3 — إشارات الراديو (`9-3-signals`)

**Focus:** type switches over `colony.Signal`; the randomized message order does not identify the type or route. The example uses `any` with ordinary Go values because the colony signal interface is sealed.

**Word count:** Arabic 255; English 316.

1. **المهمة / Mission:** visit every distress location and total supply values; ignore noise.
2. **الفكرة / Idea:** an interface stores a concrete type and value; each switch case gives `v` its case type.
3. **جرّب / Try it:** predict the case selected for the string `"3"`, then view the trace.
4. **خد بالك / Watch out:** unchecked assertions can panic; `Move` costs a tick even on failure.
5. **فكّر واختبر / Think & check:** match `Distress.At` and understand assertion risk.
6. **دورك / Your turn:** replace diagnostic printing with collection, visits, and the report.

`c9m3.switch-trace` · illustrative values

<!-- go-example id=e9m3-switch imports=fmt outcome=ok -->
```go
total := 0
for _, msg := range []any{7, "static"} {
	switch v := msg.(type) {
	case int:
		total += v
	case string:
		fmt.Println("ignore", v)
	}
}
fmt.Println("sum", total)
```

Frames: line 1 total starts at 0; lines 2 and 5 read and add int 7; the next value matches string and prints `ignore static`; line 10 prints the total.

#### Verified `e9m3-switch`

Real stdout:

```text
ignore static
sum 7
```

`c9m3.assert-predict` → `c9m3.assert-reveal` · predict, then reveal

<!-- go-example id=e9m3-type-check imports=fmt outcome=ok -->
```go
signal := any("3")
switch v := signal.(type) {
case int:
	fmt.Println(v + 1)
case string:
	fmt.Println("text", v)
}
```

Frames: line 1 stores a string; line 2 selects the string case; line 6 prints `text 3`.

#### Verified `e9m3-type-check`

Real stdout:

```text
text 3
```

**Hints, in order:** concept — reveal the concrete signal type; structure — collect data separately from visits and reporting; specific — inspect the `Signals()` loop and the `At` and `Amount` fields.

## 9-4 — مقابس صغيرة (`9-4-small-plugs`)

**Focus:** small `io.Reader`/`io.Writer` interfaces, line scanning, trimming, and distinguishing the `end` command from EOF. The existing fallback lesson stays available.

**Word count:** Arabic 269; English 319.

1. **المهمة / Mission:** make `RelayOrders` work with injected readers and writers, returning errors; print an address or `unknown` for each key.
2. **الفكرة / Idea:** scan one line and remove surrounding spaces.
3. **جرّب / Try it:** predict the result after trimming `" end \n"`, then reveal it.
4. **خد بالك / Watch out:** `end` returns early; after scanning, inspect `Err`; check writes.
5. **فكّر واختبر / Think & check:** choose the small input contract and separate command from EOF.
6. **دورك / Your turn:** implement the contract and try an in-memory reader and buffer.

`c9m4.scan-trace` · one scanner read

<!-- go-example id=e9m4-scan imports=bufio,fmt,strings outcome=ok -->
```go
scanner := bufio.NewScanner(strings.NewReader("  beacon-7 \n"))
if scanner.Scan() {
	key := strings.TrimSpace(scanner.Text())
	fmt.Println(key)
}
fmt.Println(scanner.Err())
```

Frames: line 1 creates the scanner; line 2 reads; line 3 trims; line 4 prints; line 6 shows no read error.

#### Verified `e9m4-scan`

Real stdout:

```text
beacon-7
<nil>
```

`c9m4.trim-predict` → `c9m4.trim-reveal` · predict, then reveal

<!-- go-example id=e9m4-trim imports=fmt,strings outcome=ok -->
```go
key := strings.TrimSpace(" end \n")
fmt.Println(key == "end", key)
```

Frames: line 1 sets `key` to `end`; line 2 prints the comparison and key.

#### Verified `e9m4-trim`

Real stdout:

```text
true end
```

**Hints, in order:** concept — request only the interface behavior needed; structure — separate line reading from command handling and writing; specific — inspect `scanner.Scan()` and `scanner.Err()`.

## Validation notes

Storyboard code is compiled and run as the body of `main` with the listed standard-library imports. `outcome=ok` means it builds and exits successfully; the verified stdout blocks are also checked against the terminal walkthrough frames. No chapter 9 setup, starter, checks, or grading values are changed.

## 9-boss — شحنة الإخلاء (`9-boss-evacuation`)

**Focus:** an integration lab for structs, interfaces, maps, and error results. The full contract stays in the existing fallback lesson.

**Word count:** Arabic 284; English 319, including titles, captions, options, feedback, and hints; excluding code and shared controls.

1. **المهمة / Mission:** complete `Pack` from its contract and keep `main` empty.
2. **الفكرة / Idea:** structs group data; interfaces describe required methods; a signature lists inputs and results.
3. **جرّب / Try it:** predict a tiny map update, then step through it.
4. **خد بالك / Watch out:** distinguish nil and empty maps; map `range` order is unspecified; inspect errors.
5. **فكّر واختبر / Think & check:** review result paths and interface matching.
6. **دورك / Your turn:** implement the contract and use `go test` to guide one change at a time.

`c9b.map-predict` → `c9b.map-reveal` · predict, then reveal

<!-- go-example id=e9boss-map-count imports=fmt outcome=ok -->
```go
counts := map[string]int{}
counts["teal"] += 2
counts["teal"]++
fmt.Println(counts["teal"])
```

| Frame / line | Map | Output | Caption |
|---|---|---|---|
| 0 / before 1 | — | — | قبل إنشاء الـmap. / Before creating the map. |
| 1 / 1 | `map[]` | — | الـmap بدأت فاضية. / The map starts empty. |
| 2 / 2 | `map[teal:2]` | — | الـvalue بقت 2. / The value becomes 2. |
| 3 / 3 | `map[teal:3]` | — | التحديث زوّد نفس الـkey. / The same key increases. |
| 4 / 4 | `map[teal:3]` | `3` | اطبع الـvalue بتاعة الـkey. / Print the key's value. |

#### Verified `e9boss-map-count` — ok

Real stdout:

```text
3
```

**Hints, in order:** concept — a struct groups data and an interface names behavior; structure — keep input, result, and error paths distinct; specific — inspect `Pack`, `ErrManifest`, and return paths in `main.go`.

## Interactive playgrounds

Each sandbox appears after the existing prediction and walkthrough in Try it. Its data is illustrative and does not supply a mission solution.

- **9-1:** Toggle a toy `Status` between its default display and `String()` output.
- **9-2:** Step across sample tiles and watch the library call `Plan` for each one.
- **9-3:** Cycle concrete values and see which type-switch case matches.
- **9-4:** Toggle `TrimSpace` on a sample line with visible spaces.
- **9-boss:** Step through inserting and updating one sample map key.
