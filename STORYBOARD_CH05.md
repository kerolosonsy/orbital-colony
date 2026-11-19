# Chapter 5 lesson storyboards

Short bilingual lesson examples for the five chapter 5 missions. Values are illustrative; each lesson keeps the mission's existing setup and contract.

## 5-1 — Write It Once

Mission: harvest three rows using a helper with parameters. Idea: follow a parameter into a function and its result back to the caller. Try: predict a small function result. Watch out: argument and result types must match. Think: identify repeated work and changing inputs. Your turn: finish the helper and arrange calls.

<!-- go-example id=e51-params imports=fmt outcome=ok -->
```go
func twice(n int) int {
	return n * 2
}
func main() {
	value := twice(4)
	fmt.Println(value)
}
```

#### Verified e51-params — ok

Real stdout:

```text
8
```

<!-- go-example id=e51-total imports=fmt outcome=ok -->
```go
func total(a, b int) int {
	return a + b
}
func main() {
	fmt.Println(total(3, 4))
}
```

#### Verified e51-total — ok

Real stdout:

```text
7
```

## 5-2 — Crystal Storm

Mission: collect crystals during storm failures, then report the total. Idea: check the error before using the value. Try: predict the success branch. Watch out: only retry a failure that can be handled. Think: separate attempts from successful accumulation. Your turn: finish the result handling.

<!-- go-example id=e52-error imports=errors,fmt outcome=ok -->
```go
func sample() (int, error) {
	return 6, errors.New("storm")
}
func main() {
	n, err := sample()
	if err != nil {
		fmt.Println("retry")
		return
	}
	fmt.Println(n)
}
```

#### Verified e52-error — ok

Real stdout:

```text
retry
```

<!-- go-example id=e52-success imports=fmt outcome=ok -->
```go
func sample() (int, error) {
	return 6, nil
}
func main() {
	n, err := sample()
	if err != nil {
		fmt.Println("retry")
	} else {
		fmt.Println(n)
	}
}
```

#### Verified e52-success — ok

Real stdout:

```text
6
```

## 5-3 — Fault Log

Mission: preserve fault identity while adding context and combining causes. Idea: wrap with `%w`; check identity with `errors.Is`. Try: predict a cause check. Watch out: `%v` displays text without wrapping. Think: keep diagnostic results in order. Your turn: inspect and report each fault.

<!-- go-example id=e53-wrap imports=errors,fmt outcome=ok -->
```go
func main() {
	cause := errors.New("valve")
	wrapped := fmt.Errorf("bay %d: %w", 2, cause)
	fmt.Println(errors.Is(wrapped, cause))
	fmt.Println(wrapped)
}
```

#### Verified e53-wrap — ok

Real stdout:

```text
true
bay 2: valve
```

<!-- go-example id=e53-cause imports=errors,fmt outcome=ok -->
```go
func main() {
	cause := errors.New("sensor")
	wrapped := fmt.Errorf("unit: %w", cause)
	fmt.Println(errors.Is(wrapped, cause))
}
```

#### Verified e53-cause — ok

Real stdout:

```text
true
```

## 5-4 — Undocking

Mission: undock after each drill attempt, including panic. Idea: deferred calls run as a function exits. Try: predict defer order. Watch out: recovery belongs directly in a deferred function. Think: isolate one attempt so the caller can continue. Your turn: complete cleanup and recovery.

<!-- go-example id=e54-defer imports=fmt outcome=ok -->
```go
func work() {
	defer fmt.Println("closed")
	fmt.Println("working")
}
func main() {
	work()
}
```

#### Verified e54-defer — ok

Real stdout:

```text
working
closed
```

<!-- go-example id=e54-order imports=fmt outcome=ok -->
```go
func work() {
	defer fmt.Print("A")
	defer fmt.Print("B")
	fmt.Print("C")
}
func main() {
	work()
	fmt.Println()
}
```

#### Verified e54-order — ok

Real stdout:

```text
CBA
```

## 5-5 — Command Table

Mission: choose a function for each command and count known commands with a closure. Idea: distinguish a function value from a call. Try: predict a captured counter. Watch out: pass a matching function type and keep the counter outside the closure. Think: preserve command order. Your turn: connect selection, execution, and counting.

<!-- go-example id=e55-function-value imports=fmt outcome=ok -->
```go
func ping() { fmt.Print("P") }
func main() {
	action := ping
	action()
	fmt.Println("!")
}
```

#### Verified e55-function-value — ok

Real stdout:

```text
P!
```

<!-- go-example id=e55-closure imports=fmt outcome=ok -->
```go
func main() {
	hits := 0
	count := func() { hits++ }
	count()
	count()
	fmt.Println(hits)
}
```

#### Verified e55-closure — ok

Real stdout:

```text
2
```
