# Chapter 27 lesson storyboards

Short bilingual lessons for the chapter labs. All data in the examples is illustrative.

## 27-1-signed-pass — Signed Pass

A toy MAC ties bytes to a key; base64url payload text remains readable.

<!-- go-example id=e271-hmac-demo imports=crypto/hmac,crypto/sha256,fmt outcome=ok -->
```go
func main() {
	message := []byte("crew=Mira")
	mac := hmac.New(sha256.New, []byte("demo-only"))
	_, _ = mac.Write(message)
	fmt.Printf("%x\n", mac.Sum(nil)[:4])
}
```

#### Verified e271-hmac-demo — ok

Real stdout:

```text
b3614ce6
```

<!-- go-example id=e271-readable-payload imports=encoding/base64,fmt outcome=ok -->
```go
func main() {
	payload := base64.RawURLEncoding.EncodeToString([]byte("crew=Mira"))
	decoded, _ := base64.RawURLEncoding.DecodeString(payload)
	fmt.Println(string(decoded))
}
```

#### Verified e271-readable-payload — ok

Real stdout:

```text
crew=Mira
```

## 27-2-sessions — Sessions

Cookie flags shape browser behavior, while a server map resolves a session ID.

<!-- go-example id=e272-cookie-flags imports=fmt,net/http outcome=ok -->
```go
func main() {
	cookie := http.Cookie{Name: "quarters", Value: "demo", Path: "/",
		HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode}
	fmt.Println(cookie.HttpOnly, cookie.Secure, cookie.SameSite == http.SameSiteStrictMode)
}
```

#### Verified e272-cookie-flags — ok

Real stdout:

```text
true true true
```

<!-- go-example id=e272-session-lookup imports=fmt outcome=ok -->
```go
func main() {
	sessions := map[string]string{"demo-id": "Mira"}
	crew, ok := sessions["demo-id"]
	fmt.Println(crew, ok)
}
```

#### Verified e272-session-lookup — ok

Real stdout:

```text
Mira true
```

## 27-3-cors — CORS

A toy preflight can respond before the inner handler; header names ignore case.

<!-- go-example id=e273-preflight imports=fmt,net/http,net/http/httptest outcome=ok -->
```go
func main() {
	nextCalls := 0
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		nextCalls++
		w.WriteHeader(http.StatusOK)
	})
	shield := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodOptions, "/data", nil)
	r.Header.Set("Access-Control-Request-Method", "GET")
	shield.ServeHTTP(w, r)
	fmt.Println(w.Code, nextCalls)
}
```

#### Verified e273-preflight — ok

Real stdout:

```text
204 0
```

<!-- go-example id=e273-header-case imports=fmt,strings outcome=ok -->
```go
func main() {
	fmt.Println(strings.EqualFold("X-Crew", "x-crew"))
}
```

#### Verified e273-header-case — ok

Real stdout:

```text
true
```

## 27-4-rate-limit — Rate Limit

Each illustrative client address has a separate balance; positive wait rounds upward.

<!-- go-example id=e274-client-tanks imports=fmt outcome=ok -->
```go
func main() {
	tanks := map[string]int{"192.0.2.4": 2, "192.0.2.9": 1}
	tanks["192.0.2.4"]--
	fmt.Println(tanks["192.0.2.4"], tanks["192.0.2.9"])
}
```

#### Verified e274-client-tanks — ok

Real stdout:

```text
1 1
```

<!-- go-example id=e274-round-up imports=fmt,math,time outcome=ok -->
```go
func main() {
	wait := 1500 * time.Millisecond
	seconds := int(math.Ceil(wait.Seconds()))
	fmt.Println(seconds)
}
```

#### Verified e274-round-up — ok

Real stdout:

```text
2
```

## 27-5-short-links — Short Links

Alphabet order selects symbols, and errors.Is follows wrapped failures.

<!-- go-example id=e275-base62-alphabet imports=fmt outcome=ok -->
```go
func main() {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fmt.Printf("%c %c %c\n", alphabet[0], alphabet[10], alphabet[61])
}
```

#### Verified e275-base62-alphabet — ok

Real stdout:

```text
0 a Z
```

<!-- go-example id=e275-wrapped-collision imports=errors,fmt outcome=ok -->
```go
func main() {
	errTaken := errors.New("code taken")
	wrapped := fmt.Errorf("save: %w", errTaken)
	fmt.Println(errors.Is(wrapped, errTaken))
}
```

#### Verified e275-wrapped-collision — ok

Real stdout:

```text
true
```

## Playground notes

- **Signed Pass:** Change an illustrative message or key and watch its readable payload remain separate from a changing toy tag.
- **Sessions:** Switch between two IDs and compare the browser-held ID with the server lookup.
- **CORS:** Cycle a preflight and a normal request; watch whether the inner handler runs.
- **Rate Limit:** Increment two client counters independently, without showing a limit threshold.
- **Short Links:** Step through a few illustrative Base62 indexes and their symbols without generating a code.
