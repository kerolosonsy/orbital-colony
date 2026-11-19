# Chapter 19 lesson storyboards

Short bilingual lessons for four labs. All examples use illustrative values. The concurrent walkthrough shows one possible schedule; goroutine scheduling is not generally ordered.

## 19-1 — Flight Durations

### Idea · `e19-1-duration`

<!-- go-example id=e19-1-duration imports=fmt,time outcome=ok -->
```go
func main() {
	d, err := time.ParseDuration("35m")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d, d.String())
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | duration parsed | — | اقرأ مدة توضيحية. | Parse an illustrative duration. |
| 2 / 3 | error checked | — | اتأكد من خطأ parsing. | Check the parse error. |
| 3 / 7 | duration | `35m0s 35m0s` | اعرض المدة بصيغتها. | Print the canonical duration. |

#### Verified e19-1-duration — ok

Real stdout:

```text
35m0s 35m0s
```

### Try it · `e19-1-add`

<!-- go-example id=e19-1-add imports=fmt,time outcome=ok -->
```go
func main() {
	start := time.Date(2031, 4, 3, 23, 40, 0, 0, time.UTC)
	arrival := start.Add(35 * time.Minute)
	fmt.Println(arrival.Format("15:04"))
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | start | — | وقت البداية قيمة مستقلة. | The start is its own value. |
| 2 / 3 | arrival next day | — | Add رجّعت وقت وصول جديد. | Add returns a new arrival. |
| 3 / 4 | formatted time | `00:15` | اعرض الساعة الجديدة. | Print the new clock time. |

#### Verified e19-1-add — ok

Real stdout:

```text
00:15
```

## 19-2 — Two Suns

### Idea · `e19-2-in`

<!-- go-example id=e19-2-in imports=fmt,time outcome=ok -->
```go
func main() {
	instant := time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC)
	zone := time.FixedZone("UTC+4", 4*60*60)
	local := instant.In(zone)
	fmt.Println(instant.Equal(local), local.Format("15:04 -0700"))
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | UTC instant | — | أنشئ لحظة بتوقيت UTC. | Create a UTC instant. |
| 2 / 3 | fixed zone | — | جهّز منطقة توضيحية. | Set an illustrative location. |
| 3 / 4 | converted view | — | اعرض نفس اللحظة محلياً. | Render the instant locally. |
| 4 / 5 | equal instants | `true 04:00 +0400` | اللحظة ثابتة؛ العرض اتغير. | Same instant, changed view. |

#### Verified e19-2-in — ok

Real stdout:

```text
true 04:00 +0400
```

### Try it · `e19-2-date`

<!-- go-example id=e19-2-date imports=fmt,time outcome=ok -->
```go
func main() {
	t := time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC)
	local := t.In(time.FixedZone("UTC-5", -5*60*60))
	fmt.Println(t.Equal(local), local.Format("2006-01-02 15:04 -0700"))
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | midnight UTC | — | ابدأ بمنتصف الليل UTC. | Start at UTC midnight. |
| 2 / 3 | same instant | — | حوّل لنفس اللحظة. | Convert the same instant. |
| 3 / 4 | previous local day | `true 2030-12-31 19:00 -0500` | اعرض اللحظة والتاريخ المحلي. | Print instant and local date. |

#### Verified e19-2-date — ok

Real stdout:

```text
true 2030-12-31 19:00 -0500
```

## 19-3 — Countdown

### Idea · `e19-3-context`

<!-- go-example id=e19-3-context imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	default:
		fmt.Println("running")
	}
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | active context | — | جهّز سياق قابل للإلغاء. | Create a cancelable context. |
| 2 / 4 | select waits | — | راقب إشارة الإلغاء. | Watch for cancellation. |
| 3 / 8 | default branch | `running` | قبل الإلغاء، البرنامج شغال. | The program is still running. |

#### Verified e19-3-context — ok

Real stdout:

```text
running
```

### Try it · `e19-3-cancel`

<!-- go-example id=e19-3-cancel imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(ctx.Err())
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | context created | — | أنشئ context جديد. | Create a new context. |
| 2 / 3 | canceled | — | نفّذ الإلغاء. | Cancel the context. |
| 3 / 4 | error value | `context canceled` | اقرأ سبب الإلغاء. | Read the cancellation reason. |

#### Verified e19-3-cancel — ok

Real stdout:

```text
context canceled
```

## 19-4 — Many Timers

### Idea · `e19-4-wait`

<!-- go-example id=e19-4-wait imports=fmt,sync outcome=ok -->
```go
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("worker")
	}()
	wg.Wait()
	fmt.Println("ready")
}
```

These frames show one possible schedule. The worker can run before or after main reaches Wait; Wait only lets main continue after Done.

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | counter ready | — | جهّز عدّاد العمال. | Prepare the worker count. |
| 2 / 4 | goroutine launched | — | ابدأ goroutine واحدة. | Start one goroutine. |
| 3 / 6 | worker ran | `worker` | العامل يطبع قبل Done. | Worker prints before Done. |
| 4 / 8 | Wait returned | `worker` | Wait رجعت بعد اكتمال العامل. | Wait returns after completion. |
| 5 / 9 | main continues | `worker` + `ready` | كمل بعد انتهاء العمال. | Continue after the workers. |

#### Verified e19-4-wait — ok

Real stdout:

```text
worker
ready
```

### Try it · `e19-4-slots`

<!-- go-example id=e19-4-slots imports=fmt outcome=ok -->
```go
func main() {
	results := make([]string, 2)
	results[1] = "fast"
	results[0] = "slow"
	fmt.Println(results)
}
```

| Frame / line | State | Output | AR caption | EN caption |
|---|---|---|---|---|
| 0 / before 1 | — | — | قبل التشغيل. | Before execution. |
| 1 / 2 | two reserved slots | — | احجز خانتين للردود. | Reserve two response slots. |
| 2 / 3 | second reply first | — | الرد التاني وصل الأول. | The second reply arrived first. |
| 3 / 4 | first reply next | — | الأول وصل بعده. | The first reply arrived next. |
| 4 / 5 | request order | `[slow fast]` | اعرض حسب ترتيب الـindex. | Print in index order. |

#### Verified e19-4-slots — ok

Real stdout:

```text
[slow fast]
```
