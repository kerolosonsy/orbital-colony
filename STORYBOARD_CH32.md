# Chapter 32 lesson storyboards

Short bilingual lessons for the six chapter 32 labs. Each lesson has six steps; Try it predicts an output, then reveals the same example line by line. All sample values are illustrative.

## 32-1 — Live Broadcast

Mission: relay whole sensor messages and stop cleanly. Idea: a message carries type and payload. Try: predict the order of two packets. Watch out: keep I/O ownership clear and let cancellation release waiters. Think: trace acceptance, reading, writing, and shutdown. Your turn: inspect every exit path.

<!-- go-example id=e321-packet imports=fmt outcome=ok -->
```go
type Packet struct { kind, text string }
func main() {
	p := Packet{"text", "air"}
	fmt.Println(p.kind, p.text)
}
```

<!-- go-example id=e321-queue imports=fmt outcome=ok -->
```go
type Packet struct { kind, text string }
func main() {
	queue := []Packet{{"text", "air"}, {"binary", "map"}}
	for _, p := range queue {
		fmt.Println(p.kind, p.text)
	}
}
```

## 32-2 — gRPC Unary

Mission: turn one sector query into one reading. Idea: unary RPC pairs typed request and reply messages. Try: predict a zero-valued measurement. Watch out: map public statuses carefully and hide private store details. Think: separate shared authorization from service logic. Your turn: follow the read-only proto contract.

<!-- go-example id=e322-unary-shape imports=fmt outcome=ok -->
```go
type Query struct { Sector string }
type Reply struct { Sector string; Pressure int64 }
func inspect(q Query) Reply { return Reply{q.Sector, -3} }
func main() {
	r := inspect(Query{"ring"})
	fmt.Println(r.Sector, r.Pressure)
}
```

<!-- go-example id=e322-zero-reading imports=fmt outcome=ok -->
```go
type Reply struct { Pressure int64 }
func main() {
	reply := Reply{Pressure: 0}
	fmt.Println(reply.Pressure)
}
```

## 32-3 — gRPC Streaming

Mission: keep ordered steps and acknowledgements separate per stream. Idea: a stream carries a sequence of messages. Try: predict a small range output. Watch out: keep state local and treat EOF as a normal half-close. Think: trace each message and its matching reply. Your turn: compare single and simultaneous streams.

<!-- go-example id=e323-ordered-steps imports=fmt outcome=ok -->
```go
type Step struct { Sequence, Delta int }
func main() {
	steps := []Step{{1, 4}, {2, -2}}
	for _, step := range steps {
		fmt.Println(step.Sequence, step.Delta)
	}
}
```

<!-- go-example id=e323-range-output imports=fmt outcome=ok -->
```go
func main() {
	for _, n := range []int{3, 5} {
		fmt.Println(n)
	}
}
```

## 32-4 — TLS and mTLS

Mission: configure the airlock peers and their trust. Idea: TLS settings describe each peer's identity expectations. Try: predict a version check. Watch out: encryption alone is not identity verification. Think: separate server presentation from client trust. Your turn: trace the handshake failures.

<!-- go-example id=e324-client-expectation imports=crypto/tls,fmt outcome=ok -->
```go
func main() {
	client := &tls.Config{ServerName: "probe.colony", MinVersion: tls.VersionTLS13}
	fmt.Println(client.ServerName, client.MinVersion == tls.VersionTLS13)
}
```

<!-- go-example id=e324-version-check imports=crypto/tls,fmt outcome=ok -->
```go
func main() {
	config := &tls.Config{MinVersion: tls.VersionTLS13}
	fmt.Println(config.MinVersion == tls.VersionTLS13)
}
```

## 32-5 — HTTP/2

Mission: read sensors concurrently on one TLS connection. Idea: ALPN offers application protocols; the handshake chooses one. Try: predict the number of offered protocols. Watch out: an offer is not proof of negotiation, and shared locks can block siblings. Think: check ALPN, response version, and connection identity. Your turn: inspect concurrent reads.

<!-- go-example id=e325-alpn-offer imports=crypto/tls,fmt outcome=ok -->
```go
func main() {
	config := &tls.Config{NextProtos: []string{"h2", "http/1.1"}}
	fmt.Println(config.NextProtos[0], len(config.NextProtos))
}
```

<!-- go-example id=e325-alpn-count imports=crypto/tls,fmt outcome=ok -->
```go
func main() {
	config := &tls.Config{NextProtos: []string{"h2", "http/1.1"}}
	fmt.Println(len(config.NextProtos))
}
```

## 32-6 — SSH

Mission: run fixed inspection commands over a trusted connection. Idea: an allowlist maps a known choice to a fixed action. Try: predict a map lookup. Watch out: pin the host identity and follow session ownership through cancellation. Think: separate validation, connection, and command output. Your turn: inspect each command outcome.

<!-- go-example id=e326-command-map imports=fmt outcome=ok -->
```go
type Command struct { label, text string }
func main() {
	commands := map[string]Command{"inspect": {"inspect", "show queue"}}
	c, ok := commands["inspect"]
	fmt.Println(ok, c.text)
}
```

<!-- go-example id=e326-allowlist-lookup imports=fmt outcome=ok -->
```go
func main() {
	commands := map[string]string{"scan": "read status"}
	text, ok := commands["scan"]
	fmt.Println(ok, text)
}
```

## Toolchain transcripts

Each example is compiled and run as a standalone Go program.

#### Verified e321-packet — ok

Real stdout:

```text
text air
```

#### Verified e321-queue — ok

Real stdout:

```text
text air
binary map
```

#### Verified e322-unary-shape — ok

Real stdout:

```text
ring -3
```

#### Verified e322-zero-reading — ok

Real stdout:

```text
0
```

#### Verified e323-ordered-steps — ok

Real stdout:

```text
1 4
2 -2
```

#### Verified e323-range-output — ok

Real stdout:

```text
3
5
```

#### Verified e324-client-expectation — ok

Real stdout:

```text
probe.colony true
```

#### Verified e324-version-check — ok

Real stdout:

```text
true
```

#### Verified e325-alpn-offer — ok

Real stdout:

```text
h2 2
```

#### Verified e325-alpn-count — ok

Real stdout:

```text
2
```

#### Verified e326-command-map — ok

Real stdout:

```text
true show queue
```

#### Verified e326-allowlist-lookup — ok

Real stdout:

```text
true read status
```
