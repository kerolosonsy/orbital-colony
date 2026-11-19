# Chapter 26 lesson storyboards

Short bilingual lessons for the chapter labs. All data in the examples is illustrative.

## 26-1-first-server — First Server

A handler sets response metadata and status; an in-memory recorder checks the result.

<!-- go-example id=e261-handler imports=fmt,net/http,net/http/httptest,strings outcome=ok -->
```go
func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("queued"))
	})
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/queue", nil))
	fmt.Println(w.Code, w.Header().Get("Content-Type"), strings.TrimSpace(w.Body.String()))
}
```

#### Verified e261-handler — ok

Real stdout:

```text
202 text/plain queued
```

<!-- go-example id=e261-head-check imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	w := httptest.NewRecorder()
	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h.ServeHTTP(w, httptest.NewRequest("GET", "/ready", nil))
	fmt.Println(w.Code, w.Body.Len())
}
```

#### Verified e261-head-check — ok

Real stdout:

```text
204 0
```

## 26-2-rest-crud — REST

A lock protects an ID and its record together; sorting operates on a list view.

<!-- go-example id=e262-locked-id imports=fmt,sync outcome=ok -->
```go
type Berth struct { Label string }
func main() {
	var mu sync.Mutex
	nextID := 0
	berths := map[int]Berth{}
	mu.Lock()
	nextID++
	berths[nextID] = Berth{Label: "west"}
	id := nextID
	label := berths[id].Label
	mu.Unlock()
	fmt.Println(id, label)
}
```

#### Verified e262-locked-id — ok

Real stdout:

```text
1 west
```

<!-- go-example id=e262-sort-copy imports=fmt,slices outcome=ok -->
```go
func main() {
	ids := []int{8, 3, 5}
	slices.Sort(ids)
	fmt.Println(ids)
}
```

#### Verified e262-sort-copy — ok

Real stdout:

```text
[3 5 8]
```

## 26-3-router — The Router

Middleware wraps in order, while a toy path splits into distinct segments.

<!-- go-example id=e263-middleware-order imports=fmt,net/http,net/http/httptest,strings outcome=ok -->
```go
func main() {
	var trace []string
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { trace = append(trace, "handler") })
	wrap := func(name string) func(http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				trace = append(trace, name+">")
				next.ServeHTTP(w, r)
				trace = append(trace, "<"+name)
			})
		}
	}
	h := wrap("outer")(wrap("inner")(handler))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/bay/7", nil))
	fmt.Println(strings.Join(trace, " "))
}
```

#### Verified e263-middleware-order — ok

Real stdout:

```text
outer> inner> handler <inner <outer
```

<!-- go-example id=e263-path-parts imports=fmt,strings outcome=ok -->
```go
func main() {
	parts := strings.Split("/bays/7", "/")
	fmt.Println(len(parts), parts[2])
}
```

#### Verified e263-path-parts — ok

Real stdout:

```text
3 7
```

## 26-4-envelope — The Envelope

errors.As finds a typed fault through a wrapper; nil fields encode as null.

<!-- go-example id=e264-wrapped-fault imports=errors,fmt outcome=ok -->
```go
type Fault struct { Kind string }
func (f *Fault) Error() string { return "private detail" }
func main() {
	wrapped := fmt.Errorf("fetch: %w", &Fault{Kind: "missing"})
	var fault *Fault
	found := errors.As(wrapped, &fault)
	fmt.Println(found, fault.Kind)
}
```

#### Verified e264-wrapped-fault — ok

Real stdout:

```text
true missing
```

<!-- go-example id=e264-json-null imports=encoding/json,fmt outcome=ok -->
```go
func main() {
	reply := map[string]any{"data": nil, "error": nil}
	raw, _ := json.Marshal(reply)
	fmt.Println(string(raw))
}
```

#### Verified e264-json-null — ok

Real stdout:

```text
{"data":null,"error":null}
```

## 26-5-static — Static Files

An in-memory fs demonstrates root selection and exact route allowlisting.

<!-- go-example id=e265-embedded-fs imports=fmt,io/fs,testing/fstest outcome=ok -->
```go
func main() {
	site := fstest.MapFS{"charts/notice.txt": &fstest.MapFile{Data: []byte("sector")}}
	charts, err := fs.Sub(site, "charts")
	body, readErr := fs.ReadFile(charts, "notice.txt")
	fmt.Println(string(body), err == nil, readErr == nil)
}
```

#### Verified e265-embedded-fs — ok

Real stdout:

```text
sector true true
```

<!-- go-example id=e265-exact-route imports=fmt outcome=ok -->
```go
func main() {
	allowed := map[string]bool{"/console": true}
	_, exact := allowed["/console"]
	_, slash := allowed["/console/"]
	fmt.Println(exact, slash)
}
```

#### Verified e265-exact-route — ok

Real stdout:

```text
true false
```

## 26-6-upload — File Upload

A stream copies into a staging area and reports EOF after its final byte.

<!-- go-example id=e266-stream-copy imports=fmt,io,strings outcome=ok -->
```go
func main() {
	source := strings.NewReader("packet")
	n, err := io.Copy(io.Discard, source)
	fmt.Println(n, err)
}
```

#### Verified e266-stream-copy — ok

Real stdout:

```text
6 <nil>
```

<!-- go-example id=e266-stream-eof imports=fmt,io,strings outcome=ok -->
```go
func main() {
	reader := strings.NewReader("x")
	_, _ = reader.Read(make([]byte, 1))
	_, err := reader.Read(make([]byte, 1))
	fmt.Println(err == io.EOF)
}
```

#### Verified e266-stream-eof — ok

Real stdout:

```text
true
```
