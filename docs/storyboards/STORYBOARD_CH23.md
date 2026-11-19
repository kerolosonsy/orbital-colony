# Chapter 23 lesson storyboards

Five short lessons cover the chapter's four labs and evidence boss. The examples use illustrative data and existing Go APIs.

## 23-1 — Access Codes

Mission: issue secret codes with a controllable random source. Idea: use secure randomness for secrets and injected readers for repeatable behavior. Try: predict a full read from a fixed byte source. Watch out: never log secrets or return partial codes. Think: validate before consuming randomness. Your turn: begin with a byte reader.

<!-- go-example id=e231-random-read imports=crypto/rand,fmt outcome=ok -->
```go
func main() {
	buffer := make([]byte, 2)
	n, err := rand.Read(buffer)
	fmt.Println(n == len(buffer), err == nil)
}
```

<!-- go-example id=e231-readfull imports=bytes,fmt,io outcome=ok -->
```go
func main() {
	reader := bytes.NewReader([]byte{9, 4})
	var sample [2]byte
	n, err := io.ReadFull(reader, sample[:])
	fmt.Println(n, sample, err)
}
```

#### Verified e231-random-read — ok

Real stdout:

```text
true true
```

#### Verified e231-readfull — ok

Real stdout:

```text
2 [9 4] <nil>
```

## 23-2 — Sealed Cargo

Mission: encrypt cargo and authenticate its manifest. Idea: AES-GCM provides an AEAD; AAD is visible but authenticated. Try: predict what changing AAD does. Watch out: use a fresh nonce per key and check packet lengths. Think: keep plaintext untrusted until open succeeds. Your turn: compare a valid manifest with a changed one.

<!-- go-example id=e232-gcm-shape imports=crypto/aes,crypto/cipher,fmt outcome=ok -->
```go
func main() {
	block, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	fmt.Println(aead.NonceSize(), aead.Overhead())
}
```

<!-- go-example id=e232-aad-check imports=crypto/aes,crypto/cipher,fmt outcome=ok -->
```go
func main() {
	block, err := aes.NewCipher(make([]byte, 16))
	if err != nil {
		panic(err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}
	nonce := make([]byte, aead.NonceSize())
	sealed := aead.Seal(nil, nonce, []byte("ok"), []byte("label"))
	_, err = aead.Open(nil, nonce, sealed, []byte("other"))
	fmt.Println(err != nil)
}
```

#### Verified e232-gcm-shape — ok

Real stdout:

```text
12 16
```

#### Verified e232-aad-check — ok

Real stdout:

```text
true
```

## 23-3 — Content Address

Mission: store compressed objects with IDs tied to kind and raw content. Idea: compression is storage representation, while the address identifies raw data. Try: predict a limited read. Watch out: finish and verify the compression stream, and bound decompression. Think: validate before returning content. Your turn: compare a valid object with a corrupt stream.

<!-- go-example id=e233-compress-hash imports=bytes,compress/zlib,crypto/sha1,fmt outcome=ok -->
```go
func main() {
	raw := []byte("sample")
	rawID := sha1.Sum(raw)
	var packed bytes.Buffer
	writer := zlib.NewWriter(&packed)
	_, writeErr := writer.Write(raw)
	closeErr := writer.Close()
	packedID := sha1.Sum(packed.Bytes())
	fmt.Println(rawID == packedID, writeErr == nil && closeErr == nil)
}
```

<!-- go-example id=e233-read-limit imports=fmt,io,strings outcome=ok -->
```go
func main() {
	limited := io.LimitReader(strings.NewReader("abcdef"), 3)
	part, err := io.ReadAll(limited)
	fmt.Println(string(part), err)
}
```

#### Verified e233-compress-hash — ok

Real stdout:

```text
false true
```

#### Verified e233-read-limit — ok

Real stdout:

```text
abc <nil>
```

## 23-4 — Snapshot Commit

Mission: stage named supplies and build a deterministic snapshot. Idea: sort names because map iteration order is not stable. Try: predict the JSON key order. Watch out: update the index after successful storage and use the supplied timestamp. Think: validate references before writing new objects. Your turn: stage a supply, then inspect a commit.

<!-- go-example id=e234-stable-order imports=fmt,slices,strings outcome=ok -->
```go
func main() {
	names := []string{"bay/z", "bay/a"}
	slices.Sort(names)
	fmt.Println(strings.Join(names, ", "))
}
```

<!-- go-example id=e234-json-order imports=encoding/json,fmt outcome=ok -->
```go
type Stamp struct {
	Tick    int64
	Message string
}
func main() {
	encoded, err := json.Marshal(Stamp{Tick: 3, Message: "ready"})
	fmt.Println(string(encoded), err)
}
```

#### Verified e234-stable-order — ok

Real stdout:

```text
bay/a, bay/z
```

#### Verified e234-json-order — ok

Real stdout:

```text
{"Tick":3,"Message":"ready"} <nil>
```

## 23-Boss — Storm Evidence

Mission: combine CSV, aggregation, JSON, authenticated encryption, and safe replacement. Idea: keep report order stable and use injected inputs. Try: predict CSV records from a reader. Watch out: do not trim values or touch the target before validation. Think: check each stage and use injected time. Your turn: inspect records before saving.

<!-- go-example id=e23boss-stable-order imports=fmt,slices outcome=ok -->
```go
func main() {
	sensors := []string{"probe-z", "probe-a"}
	slices.Sort(sensors)
	fmt.Println(sensors)
}
```

<!-- go-example id=e23boss-csv-records imports=encoding/csv,fmt,strings outcome=ok -->
```go
func main() {
	reader := csv.NewReader(strings.NewReader("station,reading\nsensor-2,7\n"))
	header, err := reader.Read()
	row, rowErr := reader.Read()
	fmt.Println(header[0], row[0], row[1], err, rowErr)
}
```

#### Verified e23boss-stable-order — ok

Real stdout:

```text
[probe-a probe-z]
```

#### Verified e23boss-csv-records — ok

Real stdout:

```text
station sensor-2 7 <nil> <nil>
```
