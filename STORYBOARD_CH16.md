# Chapter 16 lesson storyboards

Short bilingual lessons for the labs. Each Try it predicts a small illustrative result, then reveals it line by line. The examples teach concepts without supplying lab tests or implementations.

## 16-1 — Catch the Mutants

Mission: strengthen behavior coverage against broken variants. Idea: organize named data rows. Try: add two values. Watch out: a passing correct version may leave a mutant uncovered. Think: map each promise to an observation. Your turn: extend the existing test.

<!-- go-example id=e161-rows imports=fmt outcome=ok -->
```go
type Sample struct { name string; value int }
func main() {
	rows := []Sample{{"low", 3}, {"high", 7}}
	for _, row := range rows {
		fmt.Println(row.name, row.value)
	}
}
```

<!-- go-example id=e161-sum imports=fmt outcome=ok -->
```go
func main() {
	totals := []int{2, 5}
	sum := 0
	for _, n := range totals {
		sum += n
	}
	fmt.Println(sum)
}
```

#### Verified e161-rows — ok

Real stdout:

```text
low 3
high 7
```

#### Verified e161-sum — ok

Real stdout:

```text
7
```

## 16-2 — Catch the Bug

Mission: make a failing observation before repairing Median. Idea: slices can share elements. Try: mutate a view. Watch out: check input effects as well as the returned value. Think: keep Red and Green separate. Your turn: repair, then recheck the contract.

<!-- go-example id=e162-slice-view imports=fmt outcome=ok -->
```go
func main() {
	readings := []int{8, 3}
	view := readings[:1]
	view[0] = 4
	fmt.Println(readings)
}
```

<!-- go-example id=e162-tail imports=fmt outcome=ok -->
```go
func main() {
	data := []int{6, 2}
	tail := data[1:]
	tail[0]++
	fmt.Println(data)
}
```

#### Verified e162-slice-view — ok

Real stdout:

```text
[4 3]
```

#### Verified e162-tail — ok

Real stdout:

```text
[6 3]
```

## 16-3 — 100 Times Faster

Mission: preserve CountPairs behavior while improving runtime and adding a benchmark. Idea: n-squared work grows quickly. Try: predict a square. Watch out: Big-O is not elapsed time. Think: verify correctness before measuring the same workload. Your turn: compare benchmark results.

<!-- go-example id=e163-growth imports=fmt outcome=ok -->
```go
func main() {
	sizes := []int{2, 4}
	for _, n := range sizes {
		fmt.Println(n, n*n)
	}
}
```

<!-- go-example id=e163-square imports=fmt outcome=ok -->
```go
func main() {
	n := 3
	fmt.Println(n * n)
}
```

#### Verified e163-growth — ok

Real stdout:

```text
2 4
4 16
```

#### Verified e163-square — ok

Real stdout:

```text
9
```

## 16-4 — Wild Inputs

Mission: use fuzzing to expose parser defects. Idea: check actual length before reading. Try: read a byte by index. Watch out: a parser that rejects everything is still wrong. Think: separate valid acceptance from safe rejection. Your turn: start with seeds, then fuzz.

<!-- go-example id=e164-length-guard imports=fmt outcome=ok -->
```go
func main() {
	packet := []byte{11, 22}
	if len(packet) > 0 {
		fmt.Println(packet[0])
	}
}
```

<!-- go-example id=e164-byte imports=fmt outcome=ok -->
```go
func main() {
	packet := []byte{9, 4}
	fmt.Println(len(packet), packet[1])
}
```

#### Verified e164-length-guard — ok

Real stdout:

```text
11
```

#### Verified e164-byte — ok

Real stdout:

```text
2 4
```

## 16-5 — Fake Clock

Mission: test the launch gate with injected time and no sleep. Idea: time can be moved as input. Try: compare two instants. Watch out: isolate each subtest and inspect both result and error. Think: name contract scenarios. Your turn: use a new fake clock each time.

<!-- go-example id=e165-time-move imports=fmt,time outcome=ok -->
```go
func main() {
	start := time.Unix(0, 0)
	later := start.Add(2 * time.Second)
	fmt.Println(later.Sub(start))
}
```

<!-- go-example id=e165-after imports=fmt,time outcome=ok -->
```go
func main() {
	start := time.Unix(0, 0)
	later := start.Add(500 * time.Millisecond)
	fmt.Println(later.After(start))
}
```

#### Verified e165-time-move — ok

Real stdout:

```text
2s
```

#### Verified e165-after — ok

Real stdout:

```text
true
```
