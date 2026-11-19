# Chapter 17 lesson storyboards

Short bilingual lessons for layers, 3D maps, and graph search. Each Try it predicts a small illustrative result, then reveals it line by line. Examples use toy data and do not provide a mission route or lab implementation.

## 17-1 — Your First Shaft

Mission: collect iron below the surface and return. Idea: z records depth while x and y locate the horizontal tile. Try: increment a depth value. Watch out: digging and moving have different effects. Think: check inventory and current layer separately. Your turn: continue from the starter's shaft.

Interactive sandbox (`ch17-layer-coordinate`): change z and watch x and y stay fixed.

<!-- go-example id=e171-z-coordinate imports=fmt outcome=ok -->
```go
type Cell struct { X, Y, Z int }
func main() {
	p := Cell{X: 1, Y: 2, Z: 0}
	fmt.Println(p.X, p.Y, p.Z)
	p.Z = 1
	fmt.Println(p.X, p.Y, p.Z)
}
```

<!-- go-example id=e171-layer imports=fmt outcome=ok -->
```go
func main() {
	z := 0
	z++
	fmt.Println(z)
}
```

#### Verified e171-z-coordinate — ok

Real stdout:

```text
1 2 0
1 2 1
```

#### Verified e171-layer — ok

Real stdout:

```text
1
```

## 17-2 — Map the Depths

Mission: count open cells by underground layer and report ordered counts. Idea: the full volume is indexed [z][x][y]. Try: read one illustrative cell. Watch out: use each dimension's own length; Tunnel and Shaft are open. Think: verify the axis order first. Your turn: extend the starter's one-layer view.

Interactive sandbox (`ch17-volume-index`): select illustrative cells and inspect the [z][x][y] order.

<!-- go-example id=e172-grid-index imports=fmt outcome=ok -->
```go
func main() {
	layers := [][][]string{
		{{"A", "B"}, {"C", "D"}},
		{{"E", "F"}, {"G", "H"}},
	}
	fmt.Println(layers[1][0][1])
}
```

<!-- go-example id=e172-index imports=fmt outcome=ok -->
```go
func main() {
	layers := [][][]int{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}
	fmt.Println(layers[1][1][0])
}
```

#### Verified e172-grid-index — ok

Real stdout:

```text
F
```

#### Verified e172-index — ok

Real stdout:

```text
7
```

## 17-3 — The Deep Geode

Mission: find the fewest-step route to a layer 2 geode; its gold target depends on the generated map. Idea: BFS still works when a node includes z. Try: remove the first queue item. Watch out: underground movement needs tunnels or shafts. Think: list allowed neighbors and include the final dig. Your turn: add depth to the known BFS state.

Interactive sandbox (`ch17-fifo-queue`): take the front item and watch the illustrative queue shrink.

<!-- go-example id=e173-cell-key imports=fmt outcome=ok -->
```go
type Cell struct { X, Y, Z int }
func main() {
	seen := map[Cell]bool{}
	p := Cell{X: 1, Y: 0, Z: 1}
	seen[p] = true
	fmt.Println(seen[Cell{X: 1, Y: 0, Z: 1}])
	fmt.Println(seen[Cell{X: 1, Y: 0, Z: 2}])
}
```

<!-- go-example id=e173-queue imports=fmt outcome=ok -->
```go
func main() {
	queue := []string{"dock", "shaft"}
	current := queue[0]
	queue = queue[1:]
	fmt.Println(current, queue)
}
```

#### Verified e173-cell-key — ok

Real stdout:

```text
true
false
```

#### Verified e173-queue — ok

Real stdout:

```text
dock [shaft]
```

## 17-Boss — Rescue Capsule Route

Mission: implement a generic directed route helper and prove its contract. Idea: comparable types can be map keys. Try: compare complete coordinate keys. Watch out: directed links need not return, and failure must stay failure. Think: observe each promise separately and benchmark apart from correctness. Your turn: use the contract as a checklist.

Interactive sandbox (`ch17-comparable-key`): change z and compare complete illustrative `Cell` keys.

<!-- go-example id=e17b-comparable imports=fmt outcome=ok -->
```go
func mark[T comparable](seen map[T]bool, key T) {
	seen[key] = true
}
func main() {
	seen := map[string]bool{}
	mark(seen, "dock")
	fmt.Println(seen["dock"])
}
```

<!-- go-example id=e17b-cell imports=fmt outcome=ok -->
```go
type Cell struct { X, Z int }
func main() {
	known := map[Cell]bool{{X: 2, Z: 1}: true}
	fmt.Println(known[Cell{X: 2, Z: 1}], known[Cell{X: 2, Z: 2}])
}
```

#### Verified e17b-comparable — ok

Real stdout:

```text
true
```

#### Verified e17b-cell — ok

Real stdout:

```text
true false
```
