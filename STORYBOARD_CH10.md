# Chapter 10 lessons — review draft

Short bilingual lessons for the three missions. Each uses Mission → Idea → Try it → Watch out → Think & check → Your turn. The Idea step has a line-by-line walkthrough; Try it asks for a prediction before showing a second walkthrough. Values in the examples are illustrative, not mission data. All concurrency traces show one possible order.

## 10-1 — A Second Drone / درون تاني

- **Mission / المهمة:** Finish the helper drone's row and wait for its work.
- **Idea / الفكرة:** Walk through `sync.WaitGroup`: register work, launch it, signal completion, then wait.
- **Try it / جرّب:** Predict a shared message printed after `Wait`; trace the write and the wait.
- **Watch out / خد بالك:** Register before launch; defer `Done`; `main` returning ends the program.
- **Think & check / فكّر واختبر:** Recall what `Wait` and `Done` guarantee.
- **Your turn / دورك:** Finish the helper's task in the starter and observe both rows.

### Idea example — WaitGroup

<!-- go-example id=e10-1-wait imports=fmt,sync outcome=ok -->
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("helper finished")
}()
wg.Wait()
fmt.Println("all done")
```

### Try example — shared message after Wait

<!-- go-example id=e10-1-result imports=fmt,sync outcome=ok -->
```go
var wg sync.WaitGroup
message := "waiting"
wg.Add(1)
go func() {
	defer wg.Done()
	message = "ready"
}()
wg.Wait()
fmt.Println(message)
```

## 10-2 — The Swarm / السرب

- **Mission / المهمة:** Give each drone independent row work, run it concurrently, and wait for all tasks.
- **Idea / الفكرة:** A loop registers tasks; the closure receives its own label. The trace is one possible order.
- **Try it / جرّب:** Predict whether two goroutine prints have a fixed order.
- **Watch out / خد بالك:** Add before launch and call `Done` once per task; print order can change.
- **Think & check / فكّر واختبر:** Separate task completion from output order.
- **Your turn / دورك:** Complete movement and harvesting in `work` and run the swarm again.

### Idea example — tasks from a loop

<!-- go-example id=e10-2-loop imports=fmt,sync outcome=unordered -->
```go
var wg sync.WaitGroup
for _, label := range []string{"gold", "silver"} {
	wg.Add(1)
	go func(name string) {
		defer wg.Done()
		fmt.Println(name)
	}(label)
}
wg.Wait()
```

### Try example — two independent prints

<!-- go-example id=e10-2-order imports=fmt,sync outcome=unordered -->
```go
var wg sync.WaitGroup
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("A")
}()
wg.Add(1)
go func() {
	defer wg.Done()
	fmt.Println("B")
}()
wg.Wait()
```

## 10-3 — Divide the Work / قسّم الشغل

- **Mission / المهمة:** Split scattered targets among three drones with route length in mind.
- **Idea / الفكرة:** See a nested slice as a list of task groups.
- **Try it / جرّب:** Predict how `range` and `len` report each group's size.
- **Watch out / خد بالك:** Equal item counts can still mean different travel times.
- **Think & check / فكّر واختبر:** Compare nearby targets and the longest route.
- **Your turn / دورك:** Run a first split, see which drone finishes last, then adjust.

### Idea example — a list of groups

<!-- go-example id=e10-3-groups imports=fmt outcome=ok -->
```go
jobs := []string{"cup", "lamp", "map", "ring"}
parts := [][]string{
	jobs[:2],
	jobs[2:],
}
fmt.Println(parts[0], parts[1])
```

### Try example — count each group

<!-- go-example id=e10-3-count imports=fmt outcome=ok -->
```go
groups := [][]int{{4, 9}, {3}}
for i, group := range groups {
	fmt.Println(i, len(group))
}
```

## Verified examples

Toolchain check: build with the local Go toolchain and language version `go1.24`, with network dependency resolution disabled. For `unordered` examples, output lines are compared as a multiset; the transcript records one observed order.

#### Verified `e10-1-wait` — ok

Real stdout:

```text
helper finished
all done
```

Real stderr: empty.

#### Verified `e10-1-result` — ok

Real stdout:

```text
ready
```

Real stderr: empty.

#### Verified `e10-2-loop` — unordered

Real stdout (one possible order):

```text
gold
silver
```

Real stderr: empty.

#### Verified `e10-2-order` — unordered

Real stdout (one possible order):

```text
A
B
```

Real stderr: empty.

#### Verified `e10-3-groups` — ok

Real stdout:

```text
[cup lamp] [map ring]
```

Real stderr: empty.

#### Verified `e10-3-count` — ok

Real stdout:

```text
0 2
1 1
```

Real stderr: empty.
