# Chapter 7 lesson storyboards

Draft for `lessons-ch07`, 2026-11-18. The structured lessons use six steps per mission. Practice predicts a small output, then reveals it in a walkthrough. All example data is illustrative.

## 7-1 — Checkpoints

Mission: visit the ordered point list. Idea: a struct groups exported coordinates. Try: predict a field read. Watch out: keep the supplied order and distinguish fields from slice indexes. Think: plan the phases of reaching and checking each point. Your turn: run a first attempt.

<!-- go-example id=e71-point-struct imports=fmt outcome=ok -->
```go
type Point struct { X, Y int }
func main() {
	p := Point{X: 3, Y: 1}
	fmt.Println(p.X, p.Y)
}
```

<!-- go-example id=e71-point-field imports=fmt outcome=ok -->
```go
type Point struct { X, Y int }
func main() {
	p := Point{X: 2, Y: 1}
	fmt.Println(p.X)
}
```

## 7-2 — Navigation System

Mission: harvest targets and leave decoys. Idea: a pointer receiver updates navigation state. Try: predict a method's effect. Watch out: movement and harvest cost ticks, including failed moves. Think: separate navigation from target handling. Your turn: compare tracked position with movement.

<!-- go-example id=e72-pointer-method imports=fmt outcome=ok -->
```go
type Nav struct { X int }
func (n *Nav) Step() { n.X++ }
func main() {
	nav := Nav{X: 2}
	nav.Step()
	fmt.Println(nav.X)
}
```

<!-- go-example id=e72-gauge-step imports=fmt outcome=ok -->
```go
type Gauge struct { N int }
func (g *Gauge) Add() { g.N++ }
func main() {
	g := Gauge{N: 4}
	g.Add()
	fmt.Println(g.N)
}
```

## 7-3 — The Shortest Route

Mission: choose a visit order for all targets. Idea: Manhattan distance counts grid moves. Try: predict a toy distance. Watch out: nearest-first is a heuristic, not an optimality guarantee. Think: compare, choose, update, repeat. Your turn: inspect total travel.

<!-- go-example id=e73-manhattan-distance imports=fmt outcome=ok -->
```go
type Spot struct { X, Y int }
func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}
func main() {
	a, b := Spot{X: 1, Y: 2}, Spot{X: 4, Y: 3}
	fmt.Println(a.Dist(b))
}
```

<!-- go-example id=e73-distance-reveal imports=fmt outcome=ok -->
```go
type Spot struct { X, Y int }
func abs(n int) int {
	if n < 0 {
		n = -n
	}
	return n
}
func (s Spot) Dist(o Spot) int {
	return abs(s.X-o.X) + abs(s.Y-o.Y)
}
func main() {
	a, b := Spot{X: 2, Y: 1}, Spot{X: 3, Y: 3}
	fmt.Println(a.Dist(b))
}
```

## 7-4 — Composite Units

Mission: embed invoice data and return a shares slice. Idea: embedding promotes a field; a method can read it. Try: predict direct field access. Watch out: preserve any remainder. Think: represent, calculate, and check the total. Your turn: compare the report with the input total.

<!-- go-example id=e74-promoted-method imports=fmt outcome=ok -->
```go
type Cargo struct { Units int }
type Ship struct { Cargo }
func (s Ship) Count() int { return s.Units }
func main() {
	ship := Ship{Cargo: Cargo{Units: 3}}
	fmt.Println(ship.Count())
}
```

<!-- go-example id=e74-promoted-field imports=fmt outcome=ok -->
```go
type Pack struct { Items int }
type Rover struct { Pack }
func main() {
	r := Rover{Pack: Pack{Items: 4}}
	fmt.Println(r.Items)
}
```

## 7-5 — Crew Exam

Mission: answer usable crop questions and report the correct share. Idea: questions are structs; range preserves an index. Try: predict index/value output. Watch out: skip an empty choices slice; Answer validates indices only. Think: inspect, choose, record, score. Your turn: run a small loop.

<!-- go-example id=e75-question-struct imports=fmt outcome=ok -->
```go
type Question struct {
	Prompt string
	Choices []string
}
func main() {
	questions := []Question{{Prompt: "ground?", Choices: []string{"soil", "ice"}}}
	for i, q := range questions {
		fmt.Println(i, q.Prompt, len(q.Choices))
	}
}
```

<!-- go-example id=e75-range-output imports=fmt outcome=ok -->
```go
func main() {
	labels := []string{"A", "B"}
	for i, label := range labels {
		fmt.Println(i, label)
	}
}
```

## Toolchain transcripts

Each example is a complete program body with imports declared by its marker.

#### Verified e71-point-struct — ok

Real stdout:

```text
3 1
```

#### Verified e71-point-field — ok

Real stdout:

```text
2
```

#### Verified e72-pointer-method — ok

Real stdout:

```text
3
```

#### Verified e72-gauge-step — ok

Real stdout:

```text
5
```

#### Verified e73-manhattan-distance — ok

Real stdout:

```text
4
```

#### Verified e73-distance-reveal — ok

Real stdout:

```text
3
```

#### Verified e74-promoted-method — ok

Real stdout:

```text
3
```

#### Verified e74-promoted-field — ok

Real stdout:

```text
4
```

#### Verified e75-question-struct — ok

Real stdout:

```text
0 ground? 2
```

#### Verified e75-range-output — ok

Real stdout:

```text
0 A
1 B
```
