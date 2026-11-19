# Chapter 33 lesson storyboards

Short bilingual lessons for all four chapter 33 labs and the boss lab. Each lesson has six steps; Try it predicts an output, then reveals the same example line by line. All sample values are illustrative.

## 33-1 — Stubborn Client

Mission: retry only when doing so is safe. Idea: request bodies are streams, and some can be recreated. Try: predict whether a request has a replay function. Watch out: protect operation semantics and close discarded responses. Think: decide replayability before network attempts. Your turn: check body and cancellation paths.

<!-- go-example id=e331-replayable-body imports=fmt,net/http,strings outcome=ok -->
```go
func main() {
	req, err := http.NewRequest("POST", "https://dock.example/load", strings.NewReader("sample"))
	fmt.Println(err == nil, req.GetBody != nil)
}
```

<!-- go-example id=e331-getbody-check imports=fmt,net/http,strings outcome=ok -->
```go
func main() {
	req, _ := http.NewRequest("POST", "https://dock.example/load", strings.NewReader("x"))
	fmt.Println(req.GetBody != nil)
}
```

## 33-2 — Reverse Proxy

Mission: route public requests to trusted backends. Idea: a handler receives a request and response writer. Try: predict the path recorded by a local handler. Watch out: trust configured routes and strip forged headers. Think: separate route checks from proxy rewriting. Your turn: test allowed and rejected paths locally.

<!-- go-example id=e332-handler-path imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	})
	req := httptest.NewRequest("GET", "/hydro/pump", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	fmt.Println(res.Body.String())
}
```

<!-- go-example id=e332-recorded-handler imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, r.URL.Path)
	})
	req := httptest.NewRequest("GET", "/air/valve", nil)
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	fmt.Println(res.Body.String())
}
```

## 33-3 — Load Balancer

Mission: share work across healthy units. Idea: use a mutex to protect shared state. Try: predict the value after one locked update. Watch out: do not keep a lock during network waits. Think: review selection, health checks, and recovery separately. Your turn: inspect state under concurrent requests.

<!-- go-example id=e333-protected-counter imports=fmt,sync outcome=ok -->
```go
type Counter struct { mu sync.Mutex; next int }
func (c *Counter) Take() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := c.next
	c.next++
	return n
}
func main() {
	c := &Counter{}
	fmt.Println(c.Take(), c.Take())
}
```

<!-- go-example id=e333-locked-update imports=fmt,sync outcome=ok -->
```go
func main() {
	mu := sync.Mutex{}
	next := 4
	mu.Lock()
	next++
	mu.Unlock()
	fmt.Println(next)
}
```

## 33-4 — Gateway

Mission: compose routing, identity, limits, and forwarding. Idea: middleware wraps a handler and may prepare a response first. Try: predict a local handler's header and body. Watch out: layer order is policy; remove client-supplied identity headers. Think: trace each request through the chain. Your turn: test rejection before proxying.

<!-- go-example id=e334-middleware-wrap imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	stamp := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Lesson", "ready")
			next.ServeHTTP(w, r)
		})
	}
	base := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "ok") })
	res := httptest.NewRecorder()
	stamp(base).ServeHTTP(res, httptest.NewRequest("GET", "/", nil))
	fmt.Println(res.Header().Get("X-Lesson"), res.Body.String())
}
```

<!-- go-example id=e334-handler-response imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Flag", "on")
		fmt.Fprint(w, "done")
	})
	res := httptest.NewRecorder()
	h.ServeHTTP(res, httptest.NewRequest("GET", "/", nil))
	fmt.Println(res.Header().Get("X-Flag"), res.Body.String())
}
```

## Boss — Survey Uplink

Mission: combine framed input, HTTPS, deadlines, and response ownership. Idea: read a length prefix before its content. Try: predict a big-endian length. Watch out: validate before allocation and close every response body. Think: trace ownership across each error boundary. Your turn: use small cases to check each stage.

<!-- go-example id=e33boss-frame-length imports=encoding/binary,fmt outcome=ok -->
```go
func main() {
	frame := []byte{0, 3, 'i', 'c', 'e'}
	length := binary.BigEndian.Uint16(frame[:2])
	fmt.Println(length, string(frame[2:]))
}
```

<!-- go-example id=e33boss-big-endian imports=encoding/binary,fmt outcome=ok -->
```go
func main() {
	var header [2]byte
	header[0], header[1] = 0, 5
	n := binary.BigEndian.Uint16(header[:])
	fmt.Println(n)
}
```

## Toolchain transcripts

Each example is compiled and run as a standalone Go program.

#### Verified e331-replayable-body — ok

Real stdout:

```text
true true
```

#### Verified e331-getbody-check — ok

Real stdout:

```text
true
```

#### Verified e332-handler-path — ok

Real stdout:

```text
/hydro/pump
```

#### Verified e332-recorded-handler — ok

Real stdout:

```text
/air/valve
```

#### Verified e333-protected-counter — ok

Real stdout:

```text
0 1
```

#### Verified e333-locked-update — ok

Real stdout:

```text
5
```

#### Verified e334-middleware-wrap — ok

Real stdout:

```text
ready ok
```

#### Verified e334-handler-response — ok

Real stdout:

```text
on done
```

#### Verified e33boss-frame-length — ok

Real stdout:

```text
3 ice
```

#### Verified e33boss-big-endian — ok

Real stdout:

```text
5
```
