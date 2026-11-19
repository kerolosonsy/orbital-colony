# Chapter 31 lesson storyboards

Short bilingual lessons for the six networking labs. Values are illustrative. Walkthroughs show one deterministic Go example each; no example implements a lab or supplies its tests.

## 31-1 — Echo

Mission: echo complete lines in order under a byte bound. Idea: TCP reads expose stream chunks, not messages. Try: compare UTF-8 bytes and runes. Watch: incomplete EOF and partial writes. Think: track ownership and earlier replies. Your turn: follow one read state.

<!-- go-example id=e311-reader-chunk imports=fmt,strings outcome=ok -->
```go
func main() {
	r := strings.NewReader("signal")
	chunk := make([]byte, 3)
	n, _ := r.Read(chunk)
	fmt.Println(string(chunk[:n]))
}
```

<!-- go-example id=e311-byte-count imports=fmt,unicode/utf8 outcome=ok -->
```go
func main() {
	word := "é"
	fmt.Println(len(word), utf8.RuneCountInString(word))
}
```

## 31-2 — Radio Room

Mission: keep one owner for room state and isolate slow readers. Idea: a nonblocking channel send lets the owner loop continue. Try: predict a send to an empty mailbox. Watch: channel ownership and acknowledgements. Think: process each command before replying. Your turn: add one state transition at a time.

<!-- go-example id=e312-full-mailbox imports=fmt outcome=ok -->
```go
func main() {
	mail := make(chan string, 1)
	mail <- "first"
	select {
	case mail <- "later":
		fmt.Println("sent")
	default:
		fmt.Println("busy")
	}
}
```

<!-- go-example id=e312-empty-mailbox imports=fmt outcome=ok -->
```go
func main() {
	mail := make(chan string, 1)
	select {
	case mail <- "ping":
		fmt.Println("sent")
	default:
		fmt.Println("full")
	}
}
```

## 31-3 — Neighbor Discovery

Mission: handle one bounded UDP request and reply to its sender. Idea: datagrams preserve packet boundaries; the codec can be tested without a socket. Try: predict a packet's string output. Watch: detect oversized packets and do not split short writes. Think: validate before replying. Your turn: connect one request to injected transport.

<!-- go-example id=e313-packet-fields imports=fmt,strings outcome=ok -->
```go
func main() {
	fields := strings.Split("v2|r7|pier", "|")
	fmt.Println(fields[0], fields[2])
}
```

<!-- go-example id=e313-one-packet imports=fmt outcome=ok -->
```go
func main() {
	packet := []byte("v2|r7|pier")
	fmt.Println(string(packet))
}
```

## 31-4 — Binary Frames

Mission: parse bounded archived frames and fuzz the parser contract. Idea: read multibyte fields with explicit byte order. Try: predict a field and remaining length. Watch: check bounds before slicing. Think: preserve error precedence and input ownership. Your turn: start with header checks.

<!-- go-example id=e314-big-endian imports=encoding/binary,fmt outcome=ok -->
```go
func main() {
	wire := []byte{1, 2}
	n := binary.BigEndian.Uint16(wire)
	fmt.Println(n)
}
```

<!-- go-example id=e314-header-length imports=encoding/binary,fmt outcome=ok -->
```go
func main() {
	wire := []byte{0, 9, 7}
	n := binary.BigEndian.Uint16(wire[:2])
	fmt.Println(n, len(wire[2:]))
}
```

## 31-5 — Address Directory

Mission: gather deterministic address cards through an injected resolver. Idea: classify normalized IPs and distinguish DNS absence from an empty answer. Try: predict typed error extraction. Watch: keep resolver data unchanged. Think: validate and check each call in order. Your turn: separate input and resolver errors.

<!-- go-example id=e315-ip-family imports=fmt,net outcome=ok -->
```go
func main() {
	ip := net.ParseIP("192.0.2.8")
	fmt.Println(ip.To4() != nil, ip.String())
}
```

<!-- go-example id=e315-dns-cause imports=errors,fmt,net outcome=ok -->
```go
func main() {
	cause := &net.DNSError{IsNotFound: true}
	var dns *net.DNSError
	fmt.Println(errors.As(cause, &dns), dns.IsNotFound)
}
```

## 31-6 — Port Check

Mission: check only permitted loopback targets with bounded dialing. Idea: store each result at its input index. Try: predict IPv6 host-port formatting. Watch: reject the whole plan before any worker. Think: distinguish configuration failure from port failure. Your turn: validate, then inspect each result slot.

<!-- go-example id=e316-join-host-port imports=fmt,net outcome=ok -->
```go
func main() {
	host := "::1"
	port := "7443"
	address := net.JoinHostPort(host, port)
	fmt.Println(address)
}
```

<!-- go-example id=e316-result-slots imports=fmt outcome=ok -->
```go
func main() {
	ports := []string{"first", "second"}
	states := make([]string, len(ports))
	states[1] = "closed"
	states[0] = "open"
	fmt.Println(states)
}
```

## Toolchain transcripts

#### Verified e311-reader-chunk — ok

Real stdout:

```text
sig
```

#### Verified e311-byte-count — ok

Real stdout:

```text
2 1
```

#### Verified e312-full-mailbox — ok

Real stdout:

```text
busy
```

#### Verified e312-empty-mailbox — ok

Real stdout:

```text
sent
```

#### Verified e313-packet-fields — ok

Real stdout:

```text
v2 pier
```

#### Verified e313-one-packet — ok

Real stdout:

```text
v2|r7|pier
```

#### Verified e314-big-endian — ok

Real stdout:

```text
258
```

#### Verified e314-header-length — ok

Real stdout:

```text
9 1
```

#### Verified e315-ip-family — ok

Real stdout:

```text
true 192.0.2.8
```

#### Verified e315-dns-cause — ok

Real stdout:

```text
true true
```

#### Verified e316-join-host-port — ok

Real stdout:

```text
[::1]:7443
```

#### Verified e316-result-slots — ok

Real stdout:

```text
[open closed]
```
