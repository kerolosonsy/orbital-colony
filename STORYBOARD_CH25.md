# Chapter 25 lesson storyboards

Walkthrough data is illustrative. Try it predictions use the same small programs as their reveals.

## 25-1 — Fuel Account

Mission: keep concurrent loading and burning within account invariants. Idea: a mutex protects a complete check and update, including reads. Try: predict a protected balance. Watch out: unlocked reads and partial updates. Think: identify invariants. Your turn: trace each API call.

## 25-2 — Race Report

Mission: demonstrate a lost update, explain it, and repair aggregation. Idea: stale reads can overwrite one another; the example illustrates one possible order. Try: predict the final write. Watch out: scheduling is not guaranteed. Think: use synchronization evidence, not timing. Your turn: locate the protected operation.

## 25-3 — Token Bucket

Mission: admit permits using an injected clock. Idea: duration division counts complete intervals. Try: predict interval count. Watch out: fractions and the time marker. Think: reason from elapsed time. Your turn: trace state under the mutex.

## 25-4 — Event Bus

Mission: publish to subscribers without waiting on a slow channel. Idea: select default makes a send nonblocking. Try: predict the chosen branch. Watch out: channel closure racing with sends. Think: separate topics and subscriptions. Your turn: trace ownership.

## 25-5 — TTL Cache

Mission: expire readings and join a cleanup worker. Idea: RWMutex supports shared reads; expiry comes from the clock. Try: predict stop and done behavior. Watch out: read versus write locks and worker lifetime. Think: inspect lookup, sweep, and close. Your turn: trace shutdown.

## 25-6 — Atomic or Mutex

Mission: compare correct counters with identical benchmark work. Idea: atomic Add and Load handle independent counter operations. Try: predict the total. Watch out: ordinary access and unequal workloads. Think: define a shared contract. Your turn: check totals after workers finish.

## 25-Boss — Emergency Refinery

Mission: combine bounded work, ordered results, cancellation, and joined shutdown. Idea: channel closure marks the end while queued values remain readable. Try: predict the receive. Watch out: owner and worker completion order. Think: validate, distribute, collect, and join. Your turn: map ownership across the contract.

## Walkthrough examples

<!-- go-example id=e251-locked-tank imports=fmt,sync outcome=ok -->
```go
type Tank struct { mu sync.Mutex; fuel int }
func (t *Tank) Add(n int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.fuel += n
}
func (t *Tank) Fuel() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.fuel
}
func main() {
	tank := Tank{fuel: 3}
	tank.Add(2)
	fmt.Println(tank.Fuel())
}
```

<!-- go-example id=e251-mutex-scope imports=fmt,sync outcome=ok -->
```go
func main() {
	var mu sync.Mutex
	fuel := 3
	mu.Lock()
	fuel += 2
	mu.Unlock()
	fmt.Println(fuel)
}
```

<!-- go-example id=e252-lost-update imports=fmt outcome=ok -->
```go
func main() {
	base := 10
	first := base + 3
	second := base + 5
	base = first
	base = second
	fmt.Println(base)
}
```

<!-- go-example id=e252-stale-snapshot imports=fmt outcome=ok -->
```go
func main() {
	base := 8
	first := base + 2
	second := base + 4
	base = first
	base = second
	fmt.Println(base)
}
```

<!-- go-example id=e253-duration-steps imports=fmt,time outcome=ok -->
```go
func main() {
	elapsed := 11 * time.Second
	interval := 5 * time.Second
	steps := elapsed / interval
	fmt.Println(int(steps))
}
```

<!-- go-example id=e253-whole-intervals imports=fmt,time outcome=ok -->
```go
func main() {
	elapsed := 16 * time.Second
	interval := 4 * time.Second
	fmt.Println(int(elapsed / interval))
}
```

<!-- go-example id=e254-nonblocking-send imports=fmt outcome=ok -->
```go
func main() {
	updates := make(chan string, 1)
	updates <- "ready"
	select {
	case updates <- "next":
		fmt.Println("sent")
	default:
		fmt.Println("full")
	}
	fmt.Println(<-updates)
}
```

<!-- go-example id=e254-select-default imports=fmt outcome=ok -->
```go
func main() {
	queue := make(chan int, 1)
	queue <- 9
	select {
	case queue <- 6:
		fmt.Println("sent")
	default:
		fmt.Println("busy")
	}
	fmt.Println(<-queue)
}
```

<!-- go-example id=e255-rwmutex-read imports=fmt,sync outcome=ok -->
```go
type Cache struct {
	mu sync.RWMutex
	values map[string]string
}
func (c *Cache) Lookup(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[key]
}
func main() {
	cache := Cache{values: map[string]string{"dock": "ready"}}
	fmt.Println(cache.Lookup("dock"))
}
```

<!-- go-example id=e255-stop-join imports=fmt outcome=ok -->
```go
func main() {
	stop := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-stop
		fmt.Println("cleanup stopped")
	}()
	close(stop)
	<-done
	fmt.Println("joined")
}
```

<!-- go-example id=e256-atomic-count imports=fmt,sync/atomic outcome=ok -->
```go
func main() {
	var pulses atomic.Int64
	pulses.Add(2)
	pulses.Add(3)
	fmt.Println(pulses.Load())
}
```

<!-- go-example id=e256-atomic-load imports=fmt,sync/atomic outcome=ok -->
```go
func main() {
	var count atomic.Int64
	count.Add(4)
	count.Add(1)
	fmt.Println(count.Load())
}
```

<!-- go-example id=e25b1-close-and-drain imports=fmt outcome=ok -->
```go
func main() {
	jobs := make(chan int, 2)
	jobs <- 4
	jobs <- 6
	close(jobs)
	for ore := range jobs {
		fmt.Println(ore)
	}
}
```

<!-- go-example id=e25b1-channel-end imports=fmt outcome=ok -->
```go
func main() {
	jobs := make(chan int, 1)
	jobs <- 8
	close(jobs)
	value, ok := <-jobs
	fmt.Println(value, ok)
}
```

## Toolchain transcripts

#### Verified e251-locked-tank — ok

```text
5
```

#### Verified e251-mutex-scope — ok

```text
5
```

#### Verified e252-lost-update — ok

```text
15
```

#### Verified e252-stale-snapshot — ok

```text
12
```

#### Verified e253-duration-steps — ok

```text
2
```

#### Verified e253-whole-intervals — ok

```text
4
```

#### Verified e254-nonblocking-send — ok

```text
full
ready
```

#### Verified e254-select-default — ok

```text
busy
9
```

#### Verified e255-rwmutex-read — ok

```text
ready
```

#### Verified e255-stop-join — ok

```text
cleanup stopped
joined
```

#### Verified e256-atomic-count — ok

```text
5
```

#### Verified e256-atomic-load — ok

```text
5
```

#### Verified e25b1-close-and-drain — ok

```text
4
6
```

#### Verified e25b1-channel-end — ok

```text
8 true
```
