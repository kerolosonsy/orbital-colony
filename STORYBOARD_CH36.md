# Chapter 36 lesson storyboards

Short bilingual lessons for the three chapter labs and the storm-control boss. Every level follows Mission, Idea, Try it, Watch out, Think & check, and Your turn. Example values are illustrative.

## 36-1 — Counters

Mission: map the metric contract and keep labels bounded. Idea: convert durations to seconds. Try: predict a fixed label group. Watch out: counters accumulate while gauges move both ways. Think: match metric shape to meaning and connect methods to the local registry. Your turn: build in small pieces.

<!-- go-example id=e361-duration imports=fmt,time outcome=ok -->
```go
func main() {
	d := 1500 * time.Millisecond
	fmt.Printf("%.1f\n", d.Seconds())
}
```

<!-- go-example id=e361-bounded-label imports=fmt outcome=ok -->
```go
func bayLabel(bay string) string {
	if bay == "amber" {
		return bay
	}
	return "other"
}
func main() {
	fmt.Println(bayLabel("amber"), bayLabel("flight-9"))
}
```

## 36-2 — Container Manager

Mission: finish the CLI behind Runtime. Idea: a small interface isolates the tool from an engine. Try: predict a method result. Watch out: invalid commands make no runtime call, and empty input does not confirm a stop. Think: separate parsing, state lookup, and action. Your turn: exercise the fake runtime.

<!-- go-example id=e362-interface imports=fmt outcome=ok -->
```go
type Starter interface {
	Start(string) string
}
type Demo struct{}
func (Demo) Start(name string) string { return "started " + name }
func main() {
	var engine Starter = Demo{}
	fmt.Println(engine.Start("probe-4"))
}
```

<!-- go-example id=e362-interface-reveal imports=fmt outcome=ok -->
```go
type Name interface { Label() string }
type Relay struct{}
func (Relay) Label() string { return "relay" }
func main() {
	var item Name = Relay{}
	fmt.Println(item.Label())
}
```

## 36-3 — Readiness

Mission: implement structured logs and database-aware probes. Idea: ReplaceAttr removes time and redacts a sensitive value. Try: predict URL.Path without a query. Watch out: liveness stays independent of the database. Think: separate probe branches and inspect shared state under concurrent requests. Your turn: check JSON and probe behavior.

<!-- go-example id=e363-structured-log imports=log/slog,os,strings outcome=ok -->
```go
func main() {
	opts := &slog.HandlerOptions{ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			return slog.Attr{}
		}
		if strings.EqualFold(a.Key, "token") {
			return slog.String(a.Key, "[REDACTED]")
		}
		return a
	}}
	log := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	log.Info("request", "token", "sample", "status", 503)
}
```

<!-- go-example id=e363-path-only imports=fmt,net/http outcome=ok -->
```go
func main() {
	req, _ := http.NewRequest(http.MethodGet, "https://ops.example.test/readyz?token=demo", nil)
	fmt.Println(req.URL.Path)
}
```

## 36-Boss — Storm Shield Control

Mission: persist commands, retry delivery, and expose readiness and metrics. Idea: an outbox tracks a durable command and its delivery state. Try: predict the marker after success versus failure. Watch out: close database rows before broker calls and keep secrets out of logs. Think: trace resource ownership across failure. Your turn: use a local database and the injected broker.

<!-- go-example id=e364-outbox-state imports=fmt outcome=ok -->
```go
type Receipt struct {
	ID        string
	Published bool
}
func main() {
	receipt := Receipt{ID: "storm-demo"}
	fmt.Println(receipt.ID, receipt.Published)
	receipt.Published = true
	fmt.Println(receipt.Published)
}
```

<!-- go-example id=e364-ack-state imports=fmt outcome=ok -->
```go
func acknowledged(publishOK bool) bool {
	if !publishOK {
		return false
	}
	return true
}
func main() {
	fmt.Println(acknowledged(false), acknowledged(true))
}
```

## Toolchain transcripts

Each example is a complete program body with imports declared by its marker.

#### Verified e361-duration — ok

Real stdout:

```text
1.5
```

#### Verified e361-bounded-label — ok

Real stdout:

```text
amber other
```

#### Verified e362-interface — ok

Real stdout:

```text
started probe-4
```

#### Verified e362-interface-reveal — ok

Real stdout:

```text
relay
```

#### Verified e363-structured-log — ok

Real stdout:

```text
{"level":"INFO","msg":"request","token":"[REDACTED]","status":503}
```

#### Verified e363-path-only — ok

Real stdout:

```text
/readyz
```

#### Verified e364-outbox-state — ok

Real stdout:

```text
storm-demo false
true
```

#### Verified e364-ack-state — ok

Real stdout:

```text
false true
```
