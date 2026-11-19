# Chapter 34 lessons — review draft

Four labs, each with the same six short steps in Arabic and English. Every walkthrough uses illustrative data and a small standard-library Go example; no example implements a lab method.

## 34-1 — Key Vault

Mission names concurrent state, expiry, and faithful snapshots. Idea walks through slice aliasing. Try asks about the exact expiry instant, then steps through the `time.Time.Before` result. Watch out covers map coordination and byte ownership; Think checks why an injected clock helps repeatable reasoning.

**Playground:** Change the first byte in a shared slice and watch the alias change too. The existing Try it question keeps the separate fixed-clock expiry boundary.

<!-- go-example id=e341-byte-alias imports=fmt outcome=ok -->
```go
func main() {
	value := []byte("air")
	alias := value
	value[0] = 'f'
	fmt.Println(string(alias), string(value))
}
```

<!-- go-example id=e341-expiry-boundary imports=fmt,time outcome=ok -->
```go
func main() {
	deadline := time.Unix(100, 0)
	now := deadline
	fmt.Println(now.Before(deadline))
}
```

## 34-2 — Write-Ahead Log

Mission contrasts a short tail with complete corruption. Idea follows `io.ReadFull` on a short input. Try predicts a big-endian round trip. Watch out separates incomplete data from a bad checksum; Think checks how the program should treat a complete invalid record.

**Playground:** Step between an intact record, a short tail, and a complete record with a mismatched checksum.

<!-- go-example id=e342-short-read imports=fmt,io,strings outcome=ok -->
```go
func main() {
	r := strings.NewReader("abc")
	buf := make([]byte, 4)
	n, err := io.ReadFull(r, buf)
	fmt.Println(n, err == io.ErrUnexpectedEOF)
}
```

<!-- go-example id=e342-binary-roundtrip imports=encoding/binary,fmt outcome=ok -->
```go
func main() {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, 9)
	fmt.Println(binary.BigEndian.Uint32(data))
}
```

## 34-3 — Compaction

Mission continues the previous ledger. Idea walks through sorting keys. Try predicts a tombstone field. Watch out explains how old immutable tables can resurrect deleted values; Think checks the role of a delete marker during reads and compaction.

**Playground:** Step through an older value hidden by a tombstone, early marker removal, and compaction before marker removal.

<!-- go-example id=e343-sort-keys imports=fmt,sort outcome=ok -->
```go
func main() {
	keys := []string{"fuel", "air"}
	sort.Strings(keys)
	fmt.Println(keys[0], keys[1])
}
```

<!-- go-example id=e343-tombstone imports=fmt outcome=ok -->
```go
type Slot struct {
	Dead bool
}
func main() {
	slot := Slot{Dead: true}
	fmt.Println(slot.Dead)
}
```

## 34-4 — Distributed Storage

Mission frames chunk storage, replicas, verification, and repair. Idea follows equal content fingerprints. Try predicts the hex identifier length. Watch out distinguishes a successful fetch from intact bytes; Think asks whether a copy with the wrong fingerprint can be used.

**Playground:** Trace a damaged copy, a verified fallback replica, and repair from that intact copy.

<!-- go-example id=e344-content-hash imports=crypto/sha256,fmt outcome=ok -->
```go
func main() {
	a := sha256.Sum256([]byte("ore"))
	b := sha256.Sum256([]byte("ore"))
	fmt.Println(a == b, len(a))
}
```

<!-- go-example id=e344-hex-length imports=crypto/sha256,encoding/hex,fmt outcome=ok -->
```go
func main() {
	sum := sha256.Sum256([]byte("ore"))
	id := hex.EncodeToString(sum[:])
	fmt.Println(len(id))
}
```

## Toolchain transcripts

#### Verified e341-byte-alias — ok

Real stdout:

```text
fir fir
```

#### Verified e341-expiry-boundary — ok

Real stdout:

```text
false
```

#### Verified e342-short-read — ok

Real stdout:

```text
3 true
```

#### Verified e342-binary-roundtrip — ok

Real stdout:

```text
9
```

#### Verified e343-sort-keys — ok

Real stdout:

```text
air fuel
```

#### Verified e343-tombstone — ok

Real stdout:

```text
true
```

#### Verified e344-content-hash — ok

Real stdout:

```text
true 32
```

#### Verified e344-hex-length — ok

Real stdout:

```text
64
```
