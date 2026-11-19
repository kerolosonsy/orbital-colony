# Chapter 37 lesson storyboards

Short bilingual lessons for the five command-center stages. Every level follows Mission, Idea, Try it, Watch out, Think & check, and Your turn. Examples use illustrative values and do not implement a lab contract.

## 37-1 — Foundation

Mission: connect migrations, a repository, and HTTP across two database drivers. Idea: a transaction preserves the old state on failure and commits a version on success. Try: predict the version list after append. Watch out: keep the ledger change with the schema change. Think: separate persistence layers. Your turn: start with migrations.

<!-- go-example id=e371-atomic-step imports=fmt outcome=ok -->
```go
func stage(applied []int, version int, commit bool) []int {
	if !commit {
		return applied
	}
	return append(applied, version)
}
func main() {
	current := []int{2}
	fmt.Println(stage(current, 3, false))
	fmt.Println(stage(current, 3, true))
}
```

<!-- go-example id=e371-append-version imports=fmt outcome=ok -->
```go
func main() {
	versions := []int{4}
	fmt.Println(append(versions, 5))
}
```

## 37-2 — Sign-in Gate

Mission: add registration, sign-in, and opaque sessions. Idea: a fixed-size digest can be compared without storing the bearer token. Try: predict SHA-256's byte length. Watch out: password hashing and session hashing serve different purposes. Think: separate password checks from session lifecycle. Your turn: follow signup, login, logout, and password change.

<!-- go-example id=e372-session-digest imports=crypto/sha256,fmt outcome=ok -->
```go
func main() {
	token := "session-demo"
	digest := sha256.Sum256([]byte(token))
	fmt.Println(len(token), len(digest))
}
```

<!-- go-example id=e372-digest-length imports=crypto/sha256,fmt outcome=ok -->
```go
func main() {
	digest := sha256.Sum256([]byte("sample"))
	fmt.Println(len(digest))
}
```

## 37-3 — Isolated Crews

Mission: add crew roles, isolation, and audit records. Idea: membership is keyed by both crew and user. Try: predict a lookup for one user in a second crew. Watch out: outsiders receive not found, while audit follows the change. Think: verify identity, membership, and role before acting. Your turn: follow one membership request.

<!-- go-example id=e373-tenant-key imports=fmt outcome=ok -->
```go
type MembershipKey struct {
	Crew string
	User string
}
func main() {
	roles := map[MembershipKey]string{
		{Crew: "blue", User: "nova"}: "lead",
		{Crew: "gold", User: "nova"}: "reader",
	}
	fmt.Println(roles[MembershipKey{Crew: "blue", User: "nova"}], roles[MembershipKey{Crew: "gold", User: "nova"}])
}
```

<!-- go-example id=e373-scope-reveal imports=fmt outcome=ok -->
```go
type Scope struct { Crew, User string }
func main() {
	roles := map[Scope]string{{"blue", "nova"}: "lead", {"gold", "nova"}: "reader"}
	fmt.Println(roles[Scope{"gold", "nova"}])
}
```

## 37-4 — Task Board

Mission: combine idempotent creation, versioned updates, and cursor pages. Idea: an encoded stable ID can travel as an opaque cursor. Try: predict the ID after decoding. Watch out: retries, stale versions, and page continuation are distinct cases. Think: review each endpoint separately. Your turn: build and check one operation at a time.

<!-- go-example id=e374-cursor-encode imports=encoding/base64,fmt outcome=ok -->
```go
func main() {
	cursor := base64.RawURLEncoding.EncodeToString([]byte("row-amber"))
	fmt.Println(cursor)
}
```

<!-- go-example id=e374-cursor-decode imports=encoding/base64,fmt outcome=ok -->
```go
func main() {
	id, _ := base64.RawURLEncoding.DecodeString("cm93LWJsdWU")
	fmt.Println(string(id))
}
```

## 37-5 — Operations Room

Mission: wrap the API with request IDs, errors, metrics, limits, safe logs, and graceful shutdown. Idea: URL.Path omits the query string. Try: predict the error after canceling a context. Watch out: use a fresh timed context for drain and protect shared state. Think: trace one request through response and log. Your turn: exercise the wrapper and local listener.

<!-- go-example id=e375-path-log imports=fmt,net/http outcome=ok -->
```go
func main() {
	req, _ := http.NewRequest(http.MethodGet, "https://ops.example.test/v1/tasks?token=demo", nil)
	fmt.Println(req.URL.Path)
}
```

<!-- go-example id=e375-canceled-context imports=context,fmt outcome=ok -->
```go
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(ctx.Err())
}
```

## Toolchain transcripts

Each example is a complete program body with imports declared by its marker.

#### Verified e371-atomic-step — ok

Real stdout:

```text
[2]
[2 3]
```

#### Verified e371-append-version — ok

Real stdout:

```text
[4 5]
```

#### Verified e372-session-digest — ok

Real stdout:

```text
12 32
```

#### Verified e372-digest-length — ok

Real stdout:

```text
32
```

#### Verified e373-tenant-key — ok

Real stdout:

```text
lead reader
```

#### Verified e373-scope-reveal — ok

Real stdout:

```text
reader
```

#### Verified e374-cursor-encode — ok

Real stdout:

```text
cm93LWFtYmVy
```

#### Verified e374-cursor-decode — ok

Real stdout:

```text
row-blue
```

#### Verified e375-path-log — ok

Real stdout:

```text
/v1/tasks
```

#### Verified e375-canceled-context — ok

Real stdout:

```text
context canceled
```
