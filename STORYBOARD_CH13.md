# Chapter 13 lesson storyboards

Short bilingual lesson examples for the four chapter 13 missions. Examples use small illustrative values. The goroutine trace shows one possible order; Go does not promise a scheduling order.

## 13-1 — Recall Order

Mission: let three drones work until recall, then return them and wait for all workers. Idea: `Done` signals cancellation and `Err` gives its reason. Try: predict the error after manual cancellation. Watch out: cancellation is cooperative and a WaitGroup only waits. Think: stop work, return, then join. Your turn: finish the worker's cancellation path.

<!-- go-example id=e131-context imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	fmt.Println(ctx.Err() == nil)
	cancel()
	<-ctx.Done()
	fmt.Println(ctx.Err() == context.Canceled)
}
```

#### Verified e131-context — ok

Real stdout:

```text
true
true
```

<!-- go-example id=e131-cancel imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done()
	fmt.Println(ctx.Err() == context.Canceled)
}
```

#### Verified e131-cancel — ok

Real stdout:

```text
true
```

## 13-2 — Deep Scan

Mission: scan each crystal and harvest only after a successful scan meets the purity rule. Idea: a timeout context ends a hung operation and exposes `DeadlineExceeded`. Try: predict whether a future deadline is still active. Watch out: cancel child contexts after each operation. Think: distinguish timeout from other errors. Your turn: add a bounded scan and continue after timeouts.

<!-- go-example id=e132-timeout imports=context,errors,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	<-ctx.Done()
	fmt.Println(errors.Is(ctx.Err(), context.DeadlineExceeded))
}
```

#### Verified e132-timeout — ok

Real stdout:

```text
true
```

<!-- go-example id=e132-active imports=context,fmt,time outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()
	fmt.Println(ctx.Err() == nil)
}
```

#### Verified e132-active — ok

Real stdout:

```text
true
```

## 13-3 — The Legendary Crystal

Mission: share a search among three drones; the first finder cancels the others, and all return. Idea: cancellation reaches workers through `Done`. Try: predict `Err` after cancellation. Watch out: context cancellation does not move a drone or forcibly stop a goroutine. Think: have workers check the shared signal and wait for completion. Your turn: connect a successful search to cancellation.

<!-- go-example id=e133-cancel imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	finished := make(chan bool, 1)
	go func() {
		<-ctx.Done()
		finished <- true
	}()
	cancel()
	<-finished
	fmt.Println(ctx.Err() == context.Canceled)
}
```

#### Verified e133-cancel — ok

Real stdout:

```text
true
```

<!-- go-example id=e133-done imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done()
	fmt.Println(ctx.Err())
}
```

#### Verified e133-done — ok

Real stdout:

```text
context canceled
```

## 13-4 — Chain of Command

Mission: carry context through command, relay, and probe, preserving its stop reason. Idea: pass the same context into each call and return `Err` outward. Try: predict whether relay preserves cancellation. Watch out: a new `Background` breaks the parent lifetime. Think: inspect entry and exit traces. Your turn: repair propagation and exercise both stop reasons.

<!-- go-example id=e134-chain imports=context,fmt outcome=ok -->
```go
func probe(ctx context.Context) error { return ctx.Err() }
func relay(ctx context.Context) error { return probe(ctx) }
func command(ctx context.Context) error { return relay(ctx) }
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := command(ctx)
	fmt.Println(err == context.Canceled)
}
```

#### Verified e134-chain — ok

Real stdout:

```text
true
```

<!-- go-example id=e134-relay imports=context,fmt outcome=ok -->
```go
func relay(ctx context.Context) error { return probe(ctx) }
func probe(ctx context.Context) error { return ctx.Err() }
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := relay(ctx)
	fmt.Println(err == context.Canceled)
}
```

#### Verified e134-relay — ok

Real stdout:

```text
true
```
