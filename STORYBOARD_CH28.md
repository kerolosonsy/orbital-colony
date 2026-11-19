# Chapter 28 lesson storyboard

Short bilingual practice for the four framework/API labs and the dock-quarantine boss. Each lesson has six steps; practice predicts output before revealing a line-by-line trace. Examples use illustrative data and standard-library Go.

## 28-1 — Gin

AR: افصل الراوتنج عن ملكية السجل، ونظّف المدخل قبل تغيير الحالة. EN: Separate routing from registry state and normalize input before updates.

<!-- go-example id=e281-normalize imports=fmt,strings outcome=ok -->
```go
func main() {
	label := strings.TrimSpace("  Vega  ")
	fmt.Println(label)
}
```

<!-- go-example id=e281-trim imports=fmt,strings outcome=ok -->
```go
func main() {
	label := strings.TrimSpace("  Luna  ")
	fmt.Println(label)
}
```

## 28-2 — Fiber

AR: اختبر التطبيق في الذاكرة، وانسخ القيم اللي لازم تعيش بعد handler. EN: Test in memory and copy values that must outlive a handler.

<!-- go-example id=e282-owned-string imports=fmt,strings outcome=ok -->
```go
func main() {
	raw := []byte("north")
	saved := strings.Clone(string(raw))
	raw[0] = 's'
	fmt.Println(saved, string(raw))
}
```

<!-- go-example id=e282-string-copy imports=fmt outcome=ok -->
```go
func main() {
	raw := []byte("east")
	saved := string(raw)
	raw[0] = 'w'
	fmt.Println(saved)
}
```

## 28-3 — OpenAPI contract

AR: فكّ JSON لشكل غير معروف بأمان، وافصل تحميل العقد عن فحص الرد. EN: Decode unknown JSON safely and keep contract loading separate from response checks.

<!-- go-example id=e283-json-shape imports=encoding/json,fmt outcome=ok -->
```go
func main() {
	var value any
	if err := json.Unmarshal([]byte(`{"depth": 3}`), &value); err != nil {
		panic(err)
	}
	object, ok := value.(map[string]any)
	depth, number := object["depth"].(float64)
	fmt.Println(ok, depth, number)
}
```

<!-- go-example id=e283-type-check imports=encoding/json,fmt outcome=ok -->
```go
func main() {
	var value any
	if err := json.Unmarshal([]byte(`{"depth": "3"}`), &value); err != nil {
		panic(err)
	}
	object := value.(map[string]any)
	_, ok := object["depth"].(float64)
	fmt.Println(ok)
}
```

## 28-4 — Graceful shutdown

AR: ميّز بداية الإغلاق عن انتهاء الطلبات، واجمع خطأ التنضيف مع الخطأ الأصلي. EN: Distinguish shutdown starting from requests finishing, and preserve cleanup errors.

<!-- go-example id=e284-joined-errors imports=errors,fmt outcome=ok -->
```go
func main() {
	shutdownErr := errors.New("drain failed")
	releaseErr := errors.New("close failed")
	combined := errors.Join(shutdownErr, releaseErr)
	fmt.Println(errors.Is(combined, releaseErr))
}
```

<!-- go-example id=e284-once-guard imports=fmt,sync/atomic outcome=ok -->
```go
func main() {
	var used atomic.Bool
	fmt.Println(used.CompareAndSwap(false, true))
	fmt.Println(used.CompareAndSwap(false, true))
}
```

## Boss — Dock quarantine

AR: ثبّت status والهيدر قبل كتابة الرد، وخلّي التحقق يسبق أي تعديل. EN: Set status and headers before writing the response; validate before mutation.

<!-- go-example id=e28b-response imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	w := httptest.NewRecorder()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "{\"state\":\"ready\"}")
	fmt.Println(w.Code, w.Header().Get("Content-Type"))
}
```

<!-- go-example id=e28b-forbidden imports=fmt,net/http outcome=ok -->
```go
func main() {
	status := http.StatusForbidden
	fmt.Println(status)
}
```

## Toolchain transcripts

Each example is compiled and run as a standalone Go program by `TestCh28StoryboardExamples`.

#### Verified e281-normalize — ok

Real stdout:

```text
Vega
```

#### Verified e281-trim — ok

Real stdout:

```text
Luna
```

#### Verified e282-owned-string — ok

Real stdout:

```text
north sorth
```

#### Verified e282-string-copy — ok

Real stdout:

```text
east
```

#### Verified e283-json-shape — ok

Real stdout:

```text
true 3 true
```

#### Verified e283-type-check — ok

Real stdout:

```text
false
```

#### Verified e284-joined-errors — ok

Real stdout:

```text
true
```

#### Verified e284-once-guard — ok

Real stdout:

```text
true
false
```

#### Verified e28b-response — ok

Real stdout:

```text
201 application/json
```

#### Verified e28b-forbidden — ok

Real stdout:

```text
403
```

## Playground notes

- **Gin:** Trace one cleaned label through request, handler, registry, and response.
- **Fiber:** Reuse a byte buffer and watch the saved string stay independent.
- **OpenAPI:** Cycle number, string, and missing JSON values to inspect a safe type assertion.
- **Graceful shutdown:** Run two valid schedules so active requests finish in either order before cleanup.
- **Dock quarantine boss:** Step through generic response metadata and body writing until the response commits.
