# Chapter 14 lesson storyboards

Draft for `lessons-ch14-15`, 2026-11-19. Each level has six short steps; practice predicts a small output before the walkthrough reveals it. All examples use illustrative values.

## 14-1 — One Function for Every Type

Playground: step through one generic maximum over integer and decimal slices; watch type inference and the best value. The walkthrough follows the algorithm.

<!-- go-example id=e14-1-idea imports=cmp,fmt outcome=ok -->
```go
func Max[T cmp.Ordered](xs []T) T {
	best := xs[0]
	for _, x := range xs[1:] {
		if x > best {
			best = x
		}
	}
	return best
}
func main() {
	fmt.Println(Max([]int{4, 9, 2}), Max([]float64{1.5, 3.25}))
}
```

<!-- go-example id=e14-1-try imports=cmp,fmt outcome=ok -->
```go
func Max[T cmp.Ordered](xs []T) T {
	best := xs[0]
	for _, x := range xs[1:] {
		if x > best {
			best = x
		}
	}
	return best
}
func main() {
	fmt.Println(Max([]int{8, 3, 5}))
}
```

## 14-2 — The Tunnel

Playground: instantiate one generic value with two types, then trace the last direction in a tiny route history. Connect the type parameter to the mission's direction stack.

<!-- go-example id=e14-2-idea imports=fmt outcome=ok -->
```go
type Cell[T any] struct { Value T }
func main() {
	depth := Cell[int]{Value: 12}
	turn := Cell[string]{Value: "north"}
	fmt.Println(depth.Value, turn.Value)
}
```

<!-- go-example id=e14-2-try imports=fmt outcome=ok -->
```go
func main() {
	path := []string{"North", "East"}
	next := path[len(path)-1]
	fmt.Println(next)
}
```

## 14-3 — Type Constraints

Playground: cycle through defined and built-in types to see which underlying types match the constraint's type set.

<!-- go-example id=e14-3-idea imports=fmt outcome=ok -->
```go
type Meter float64
type Numeric interface { ~int | ~float64 }
func twice[T Numeric](v T) T { return v + v }
func main() {
	fmt.Println(twice(Meter(1.5)), twice(3))
}
```

<!-- go-example id=e14-3-try imports=fmt outcome=ok -->
```go
type Meter float64
type Numeric interface { ~int | ~float64 }
func twice[T Numeric](v T) T { return v + v }
func main() {
	fmt.Println(twice(Meter(2.5)))
}
```

## 14-4 — Airlock Queue

Playground: add, read, and remove parcels in a small generic FIFO queue. Discuss the presence signal and read-versus-remove behavior without implementing the lab API.

<!-- go-example id=e14-4-idea imports=fmt outcome=ok -->
```go
type Parcel[T any] struct { Value T }
func main() {
	box := Parcel[string]{Value: "medic"}
	fmt.Println(box.Value)
}
```

<!-- go-example id=e14-4-try imports=fmt outcome=ok -->
```go
type Parcel[T any] struct { Value T }
func main() {
	box := Parcel[int]{Value: 3}
	fmt.Println(box.Value)
}
```

## 14-5 — Archive Tree

Playground: follow a recursive total down a short slice and watch returns combine. Relate the stopping case and progress to the archive's recursive contract without giving its implementation.

<!-- go-example id=e14-5-idea imports=fmt outcome=ok -->
```go
func total(xs []int, i int) int {
	if i == len(xs) {
		return 0
	}
	return xs[i] + total(xs, i+1)
}
func main() {
	fmt.Println(total([]int{2, 3, 4}, 0))
}
```

<!-- go-example id=e14-5-try imports=fmt outcome=ok -->
```go
func countDown(n int) int {
	if n == 0 {
		return 0
	}
	return 1 + countDown(n-1)
}
func main() {
	fmt.Println(countDown(3))
}
```

## Toolchain transcripts

Each example is a complete program body with imports declared by its marker.

#### Verified e14-1-idea — ok

Real stdout:

```text
9 3.25
```

#### Verified e14-1-try — ok

Real stdout:

```text
8
```

#### Verified e14-2-idea — ok

Real stdout:

```text
12 north
```

#### Verified e14-2-try — ok

Real stdout:

```text
East
```

#### Verified e14-3-idea — ok

Real stdout:

```text
3 6
```

#### Verified e14-3-try — ok

Real stdout:

```text
5
```

#### Verified e14-4-idea — ok

Real stdout:

```text
medic
```

#### Verified e14-4-try — ok

Real stdout:

```text
3
```

#### Verified e14-5-idea — ok

Real stdout:

```text
9
```

#### Verified e14-5-try — ok

Real stdout:

```text
3
```
