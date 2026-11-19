# Chapter 24 lesson storyboards

Walkthrough data is illustrative. Try it predictions use the same small programs as their reveals.

## 24-1 — Ping Pong

Mission: pass rounds between two drones and stop cleanly. Idea: an unbuffered channel is a handoff; Try: predict a buffered receive. Watch out: channel ownership and callback completion. Think: trace cancellation and errors. Your turn: map each wait point.

## 24-2 — Producers and Consumers

Mission: merge depot batches through a bounded queue. Idea: channel capacity differs from current length. Try: predict both values. Watch out: a buffer does not guarantee consumption. Think: trace ownership and failures. Your turn: follow each crate.

## 24-3 — Work Crew

Mission: return one typed reply per input from a fixed crew. Idea: a type parameter keeps a result shape reusable. Try: predict a generic field value. Watch out: item errors differ from call errors. Think: separate inputs, results, and cancellation. Your turn: trace ownership.

## 24-4 — Refining Pipeline

Mission: fan work across lanes and merge replies. Idea: closing a channel lets range drain and finish. Try: predict a receive from a closed channel. Watch out: close ownership and cancellation. Think: draw stages and waits. Your turn: trace channel direction.

## 24-5 — Beacon Check

Mission: inspect addresses with an injected HTTP client and a concurrency limit. Idea: requests carry context; response status differs from transport errors. Try: predict a status. Watch out: close response bodies. Think: separate setup, call, and result. Your turn: track ordered results.

## 24-6 — Chunk Download

Mission: stage chart chunks and publish after completion. Idea: WriteAt uses an offset, independent of completion order. Try: predict its byte count. Watch out: partial responses and early publication. Think: plan, validate, then decide. Your turn: trace failure to commit or abort.

## 24-7 — Map Tiles

Mission: resize tiles with a bounded crew and preserve format. Idea: DecodeConfig reads image metadata before pixels. Try: predict rectangle dimensions. Watch out: source format and validation. Think: process each tile through save. Your turn: follow one tile and its result.

## Walkthrough examples

<!-- go-example id=e241-channel-handoff imports=fmt outcome=ok -->
```go
func main() {
	permit := make(chan int)
	go func() { permit <- 7 }()
	code := <-permit
	fmt.Println(code)
}
```

<!-- go-example id=e241-buffer-read imports=fmt outcome=ok -->
```go
func main() {
	queue := make(chan int, 1)
	queue <- 4
	fmt.Println(<-queue)
}
```

<!-- go-example id=e242-buffer-capacity imports=fmt outcome=ok -->
```go
func main() {
	queue := make(chan int, 2)
	queue <- 4
	queue <- 6
	fmt.Println(cap(queue), len(queue))
}
```

<!-- go-example id=e242-buffer-one imports=fmt outcome=ok -->
```go
func main() {
	queue := make(chan int, 1)
	queue <- 8
	fmt.Println(cap(queue), len(queue))
}
```

<!-- go-example id=e243-generic-reply imports=fmt outcome=ok -->
```go
type Reply[T any] struct { Value T }
func main() {
	result := Reply[int]{Value: 7}
	fmt.Println(result.Value)
}
```

<!-- go-example id=e243-generic-text imports=fmt outcome=ok -->
```go
type Reply[T any] struct { Value T }
func main() {
	result := Reply[string]{Value: "ore"}
	fmt.Println(result.Value)
}
```

<!-- go-example id=e244-drain-closed imports=fmt outcome=ok -->
```go
func main() {
	jobs := make(chan int, 2)
	jobs <- 3
	jobs <- 5
	close(jobs)
	for job := range jobs {
		fmt.Println(job)
	}
}
```

<!-- go-example id=e244-closed-read imports=fmt outcome=ok -->
```go
func main() {
	ch := make(chan int)
	close(ch)
	value, ok := <-ch
	fmt.Println(value, ok)
}
```

<!-- go-example id=e245-context-request imports=context,fmt,net/http outcome=ok -->
```go
func main() {
	request, err := http.NewRequestWithContext(context.Background(), "GET", "https://beacon.example/ping", nil)
	if err != nil {
		fmt.Println("invalid")
		return
	}
	fmt.Println(request.Method, request.URL.Host)
}
```

<!-- go-example id=e245-http-status imports=fmt,net/http outcome=ok -->
```go
func main() {
	response := http.Response{StatusCode: http.StatusAccepted}
	fmt.Println(response.StatusCode)
}
```

<!-- go-example id=e246-write-at-offset imports=fmt,os outcome=ok -->
```go
func main() {
	file, err := os.CreateTemp("", "chart-example")
	if err != nil {
		fmt.Println("create failed")
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()
	n, err := file.WriteAt([]byte("map"), 2)
	fmt.Println(n, err == nil)
}
```

<!-- go-example id=e246-write-count imports=fmt,os outcome=ok -->
```go
func main() {
	file, _ := os.CreateTemp("", "part-example")
	defer os.Remove(file.Name())
	defer file.Close()
	count, _ := file.WriteAt([]byte("xy"), 4)
	fmt.Println(count)
}
```

<!-- go-example id=e247-decode-config imports=bytes,fmt,image,image/png outcome=ok -->
```go
func main() {
	pixels := image.NewRGBA(image.Rect(0, 0, 2, 1))
	var encoded bytes.Buffer
	_ = png.Encode(&encoded, pixels)
	config, format, err := image.DecodeConfig(bytes.NewReader(encoded.Bytes()))
	if err != nil {
		fmt.Println("decode failed")
		return
	}
	fmt.Println(format, config.Width, config.Height)
}
```

<!-- go-example id=e247-rect-size imports=fmt,image outcome=ok -->
```go
func main() {
	bounds := image.Rect(1, 2, 5, 8)
	fmt.Println(bounds.Dx(), bounds.Dy())
}
```

## Toolchain transcripts

#### Verified e241-channel-handoff — ok

```text
7
```

#### Verified e241-buffer-read — ok

```text
4
```

#### Verified e242-buffer-capacity — ok

```text
2 2
```

#### Verified e242-buffer-one — ok

```text
1 1
```

#### Verified e243-generic-reply — ok

```text
7
```

#### Verified e243-generic-text — ok

```text
ore
```

#### Verified e244-drain-closed — ok

```text
3
5
```

#### Verified e244-closed-read — ok

```text
0 false
```

#### Verified e245-context-request — ok

```text
GET beacon.example
```

#### Verified e245-http-status — ok

```text
202
```

#### Verified e246-write-at-offset — ok

```text
3 true
```

#### Verified e246-write-count — ok

```text
2
```

#### Verified e247-decode-config — ok

```text
png 2 1
```

#### Verified e247-rect-size — ok

```text
4 6
```
