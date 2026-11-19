# Chapter 15 lesson storyboards

Draft for `lessons-ch14-15`, 2026-11-19. Each level has six short steps; practice predicts a small output before the walkthrough reveals it. Examples use illustrative graphs, values, and routes.

## 15-1 — Cave Maze

Playground: remove the first queued item, append a new discovery, and predict the remaining order. The walkthrough connects FIFO order to BFS layers and shortest routes with equal move costs.

<!-- go-example id=e15-1-idea imports=fmt outcome=ok -->
```go
func main() {
	queue := []string{"A", "B"}
	current := queue[0]
	queue = queue[1:]
	queue = append(queue, "C")
	fmt.Println(current, queue)
}
```

<!-- go-example id=e15-1-try imports=fmt outcome=ok -->
```go
func main() {
	queue := []string{"north", "east"}
	first := queue[0]
	queue = queue[1:]
	fmt.Println(first, queue)
}
```

## 15-2 — Line Up the Crystal

Playground: compare one adjacent pair, update both slice entries, and predict the row after one swap. The discussion links each local change to the physical crystals and tick cost.

<!-- go-example id=e15-2-idea imports=fmt outcome=ok -->
```go
func main() {
	sizes := []int{5, 2, 4}
	if sizes[0] > sizes[1] {
		sizes[0], sizes[1] = sizes[1], sizes[0]
	}
	fmt.Println(sizes)
}
```

<!-- go-example id=e15-2-try imports=fmt outcome=ok -->
```go
func main() {
	sizes := []int{3, 1}
	if sizes[0] > sizes[1] {
		sizes[0], sizes[1] = sizes[1], sizes[0]
	}
	fmt.Println(sizes)
}
```

## 15-3 — Slippery Ice

Playground: add entry costs for a route through ice and one around it. Predict both totals, then connect lowest accumulated cost to Dijkstra's choice.

<!-- go-example id=e15-3-idea imports=fmt outcome=ok -->
```go
func enterCost(ice bool) int {
	if ice {
		return 4
	}
	return 1
}
func main() {
	through := enterCost(true) + enterCost(false)
	detour := enterCost(false) + enterCost(false) + enterCost(false)
	fmt.Println(through, detour)
}
```

<!-- go-example id=e15-3-try imports=fmt outcome=ok -->
```go
func enterCost(ice bool) int {
	if ice {
		return 4
	}
	return 1
}
func main() {
	through := enterCost(true) + enterCost(false) + enterCost(false)
	around := enterCost(false) + enterCost(false) + enterCost(false) + enterCost(false)
	fmt.Println(through, around)
}
```

## 15-4 — Deep Probe

Playground: move a binary-search lower bound after one comparison, then predict the count from a short recursive run. The two examples show range narrowing and a base case without completing either mission function.

<!-- go-example id=e15-4-idea imports=fmt outcome=ok -->
```go
func main() {
	xs := []int{2, 5, 9}
	lo, hi := 0, len(xs)
	mid := lo + (hi-lo)/2
	fmt.Println(xs[mid], lo, hi)
	if xs[mid] < 6 {
		lo = mid + 1
	}
	fmt.Println(lo, hi)
}
```

<!-- go-example id=e15-4-try imports=fmt outcome=ok -->
```go
func run(open []bool, i int) int {
	if i >= len(open) || !open[i] {
		return 0
	}
	return 1 + run(open, i+1)
}
func main() {
	fmt.Println(run([]bool{true, true, false}, 0))
}
```

## 15-5 — Rescue Priority

Playground: push two integers into a toy min-heap and predict the popped value. The comparison intentionally differs from rescue priority; learners connect heap ordering to the request contract.

<!-- go-example id=e15-5-idea imports=container/heap,fmt outcome=ok -->
```go
type Scores []int
func (s Scores) Len() int { return len(s) }
func (s Scores) Less(i, j int) bool { return s[i] < s[j] }
func (s Scores) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s *Scores) Push(v any) { *s = append(*s, v.(int)) }
func (s *Scores) Pop() any {
	old := *s
	last := old[len(old)-1]
	*s = old[:len(old)-1]
	return last
}
func main() {
	q := &Scores{}
	heap.Push(q, 7)
	heap.Push(q, 2)
	fmt.Println(heap.Pop(q).(int))
}
```

<!-- go-example id=e15-5-try imports=container/heap,fmt outcome=ok -->
```go
type Scores []int
func (s Scores) Len() int { return len(s) }
func (s Scores) Less(i, j int) bool { return s[i] < s[j] }
func (s Scores) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s *Scores) Push(v any) { *s = append(*s, v.(int)) }
func (s *Scores) Pop() any {
	old := *s
	last := old[len(old)-1]
	*s = old[:len(old)-1]
	return last
}
func main() {
	q := &Scores{}
	heap.Push(q, 6)
	heap.Push(q, 3)
	fmt.Println(heap.Pop(q).(int))
}
```

## Toolchain transcripts

Each example is a complete program body with imports declared by its marker.

#### Verified e15-1-idea — ok

Real stdout:

```text
A [B C]
```

#### Verified e15-1-try — ok

Real stdout:

```text
north [east]
```

#### Verified e15-2-idea — ok

Real stdout:

```text
[2 5 4]
```

#### Verified e15-2-try — ok

Real stdout:

```text
[1 3]
```

#### Verified e15-3-idea — ok

Real stdout:

```text
5 3
```

#### Verified e15-3-try — ok

Real stdout:

```text
6 4
```

#### Verified e15-4-idea — ok

Real stdout:

```text
5 0 3
2 3
```

#### Verified e15-4-try — ok

Real stdout:

```text
2
```

#### Verified e15-5-idea — ok

Real stdout:

```text
2
```

#### Verified e15-5-try — ok

Real stdout:

```text
3
```
