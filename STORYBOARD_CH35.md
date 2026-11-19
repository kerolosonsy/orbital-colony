# Chapter 35 lessons — review draft

Seven labs use the six-step bilingual format. Their snippets illustrate one Go idea at a time and do not implement lab methods. The Raft examples show deterministic state only; none claims a particular goroutine interleaving.

## 35-1 — Task Queue

Mission frames a Redis queue as task state transitions. Idea walks through a fixed lease deadline; Try predicts a map presence lookup. Watch out covers at-least-once delivery and coordinated ownership; Think checks how concurrent reclaiming is decided.

**Playground idea:** Show ready, leased, and dead-letter lanes with a controllable clock. Let the learner advance time and observe that an unacknowledged task can return for another delivery.

<!-- go-example id=e351-lease-boundary imports=fmt,time outcome=ok -->
```go
func main() {
	start := time.Unix(1000, 0)
	deadline := start.Add(3 * time.Second)
	now := start.Add(3 * time.Second)
	fmt.Println(now.Before(deadline))
}
```

<!-- go-example id=e351-map-presence imports=fmt outcome=ok -->
```go
func main() {
	seen := map[string]bool{"pump": true}
	_, ok := seen["pump"]
	fmt.Println(ok)
}
```

## 35-2 — Scheduler

Mission describes periodic jobs, missed runs, and an injected clock. Idea follows the heap root. Try predicts the minimum in a sample deadline list. Watch out distinguishes heap operations from direct slice methods; Think checks overlap behavior.

**Playground idea:** Let the learner add labeled deadlines, advance a fake clock, and inspect the next heap root. A separate toggle can mark a job running and reveal skipped overlap.

<!-- go-example id=e352-heap-root imports=container/heap,fmt outcome=ok -->
```go
type minInts []int
func (h minInts) Len() int { return len(h) }
func (h minInts) Less(i, j int) bool { return h[i] < h[j] }
func (h minInts) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *minInts) Push(x any) { *h = append(*h, x.(int)) }
func (h *minInts) Pop() any { old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x }
func main() {
	h := &minInts{7, 3, 5}
	heap.Init(h)
	fmt.Println((*h)[0])
}
```

<!-- go-example id=e352-heap-pop imports=fmt,slices outcome=ok -->
```go
func main() {
	due := []int{9, 2, 6}
	fmt.Println(slices.Min(due))
}
```

## 35-3 — Distributed Lock

Mission covers a Redis lock with an owner and expiry. Idea compares an old token with the current value; Try predicts a matching token. Watch out explains the gap between separate check/delete calls and the limit of advisory locks; Think checks atomic ownership release.

**Playground idea:** Let two workers contend for one labeled lock. Advance a clock past expiry, issue a new token, and show why the old worker cannot release the new owner's lease.

<!-- go-example id=e353-token-owner imports=fmt outcome=ok -->
```go
func main() {
	stored := "T-old"
	myToken := "T-new"
	fmt.Println(stored == myToken)
}
```

<!-- go-example id=e353-token-match imports=fmt outcome=ok -->
```go
func main() {
	stored := "T7"
	myToken := "T7"
	fmt.Println(stored == myToken)
}
```

## 35-4 — Outbox

Mission separates writing an order, relaying its event, and consuming it idempotently. Idea steps through a default preserved when JSON omits a field. Try predicts whether an event ID is already in a map. Watch out focuses on transaction boundaries and retry order.

**Playground idea:** Show the order and outbox row as one transaction, then a relay that can fail between publish and marking. Replay the duplicate event and let the learner observe a single consumer total.

<!-- go-example id=e354-json-default imports=encoding/json,fmt outcome=ok -->
```go
type Payload struct {
	Qty int
}
func main() {
	p := Payload{Qty: 1}
	_ = json.Unmarshal([]byte("{\"Item\":\"valve\"}"), &p)
	fmt.Println(p.Qty)
}
```

<!-- go-example id=e354-event-seen imports=fmt outcome=ok -->
```go
func main() {
	seen := map[int64]bool{4: true}
	_, ok := seen[4]
	fmt.Println(ok)
}
```

## 35-5 — Broker

Mission previews keyed rover messages and delayed offset commits. Idea walks through byte keys and values; Try predicts the result of parsing malformed numeric text. Watch out covers partition-local ordering and manual commits; Think checks cancellation-safe confirmation.

**Playground idea:** Put messages from several rovers into separate partition lanes by key. Step through fetch, parse, handler retry, and commit, with a cancel control before each stage.

<!-- go-example id=e355-message-bytes imports=fmt outcome=ok -->
```go
func main() {
	key := []byte("rover-7")
	value := []byte("42")
	fmt.Println(string(key), string(value))
}
```

<!-- go-example id=e355-invalid-reading imports=fmt,strconv outcome=ok -->
```go
func main() {
	reading, err := strconv.Atoi("4x")
	fmt.Println(reading, err != nil)
}
```

## 35-6 — Elections

Mission introduces a deterministic Raft election simulation. Idea steps through `iota` role values; Try predicts a seeded random draw. Watch out covers newer terms and explicit message order; Think checks the response to a higher term.

**Playground idea:** Use a small station map with a manually advanced tick and message queue. Let the learner deliver, delay, or drop one message and inspect role and term changes.

<!-- go-example id=e356-role-iota imports=fmt outcome=ok -->
```go
type Role int
const (
	Follower Role = iota
	Candidate
	Leader
)
func main() {
	fmt.Println(Follower, Candidate, Leader)
}
```

<!-- go-example id=e356-seeded-random imports=fmt,math/rand outcome=ok -->
```go
func main() {
	rng := rand.New(rand.NewSource(7))
	fmt.Println(rng.Intn(10))
}
```

## 35-7 — Log Replication

Mission extends elections with a replicated command log. Idea tracks term and slice index; Try predicts struct equality when terms differ. Watch out covers prefix matching and commit bounds; Think checks what a follower can safely commit.

**Playground idea:** Show leader and follower logs side by side. Let the learner add a conflicting entry, step back to a matching prefix, then advance the committed prefix as acknowledgements arrive.

<!-- go-example id=e357-log-index imports=fmt outcome=ok -->
```go
type Entry struct {
	Term int
	Cmd  string
}
func main() {
	log := []Entry{{Term: 1, Cmd: "scan"}, {Term: 2, Cmd: "dock"}}
	fmt.Println(len(log), log[1].Term)
}
```

<!-- go-example id=e357-entry-term imports=fmt outcome=ok -->
```go
type Entry struct {
	Term int
	Cmd  string
}
func main() {
	local := Entry{Term: 1, Cmd: "scan"}
	incoming := Entry{Term: 2, Cmd: "scan"}
	fmt.Println(local == incoming)
}
```

## Toolchain transcripts

#### Verified e351-lease-boundary — ok

Real stdout:

```text
false
```

#### Verified e351-map-presence — ok

Real stdout:

```text
true
```

#### Verified e352-heap-root — ok

Real stdout:

```text
3
```

#### Verified e352-heap-pop — ok

Real stdout:

```text
2
```

#### Verified e353-token-owner — ok

Real stdout:

```text
false
```

#### Verified e353-token-match — ok

Real stdout:

```text
true
```

#### Verified e354-json-default — ok

Real stdout:

```text
1
```

#### Verified e354-event-seen — ok

Real stdout:

```text
true
```

#### Verified e355-message-bytes — ok

Real stdout:

```text
rover-7 42
```

#### Verified e355-invalid-reading — ok

Real stdout:

```text
0 true
```

#### Verified e356-role-iota — ok

Real stdout:

```text
0 1 2
```

#### Verified e356-seeded-random — ok

Real stdout:

```text
6
```

#### Verified e357-log-index — ok

Real stdout:

```text
2 2
```

#### Verified e357-entry-term — ok

Real stdout:

```text
false
```
