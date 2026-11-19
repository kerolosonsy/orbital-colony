# Chapter 1 lesson storyboards

Short bilingual lessons for the three chapter 1 missions. Existing mission setups, starters, checks, and fallback lessons stay unchanged. Examples use illustrative Go snippets and real stdout.

## 1-1 — First Signal (`1-1-first-signal`)

The learner follows function calls from top to bottom and distinguishes direction from distance.

<!-- go-example id=e11-sequence imports=fmt outcome=ok -->
```go
func main() {
	fmt.Println("drone ready")
	fmt.Println("beacon ahead")
}
```

#### Verified e11-sequence — ok

Real stdout:

```text
drone ready
beacon ahead
```

<!-- go-example id=e11-order imports=fmt outcome=ok -->
```go
func main() {
	fmt.Println("step")
	fmt.Println("signal")
}
```

#### Verified e11-order — ok

Real stdout:

```text
step
signal
```

**Playground idea:** Show the drone's position and a short command queue. Let the learner advance one command at a time and predict its next tile.

## 1-2 — First Harvest (`1-2-first-harvest`)

The lesson introduces an empty-argument function call and connects action order with the drone's current tile.

<!-- go-example id=e12-call imports=fmt outcome=ok -->
```go
func harvest() {
	fmt.Println("harvest")
}
func main() {
	harvest()
}
```

#### Verified e12-call — ok

Real stdout:

```text
harvest
```

<!-- go-example id=e12-order imports=fmt outcome=ok -->
```go
func main() {
	fmt.Println("move")
	fmt.Println("harvest")
}
```

#### Verified e12-order — ok

Real stdout:

```text
move
harvest
```

**Playground idea:** Put the drone beside three labeled plants and expose a two-command queue. Highlight the tile under the drone when Harvest runs.

## 1-3 — New Ground (`1-3-planting`)

The lesson focuses on state and prerequisite order: prepare ground before planting, while accounting for the starter's first tile.

<!-- go-example id=e13-order imports=fmt outcome=ok -->
```go
func main() {
	fmt.Println("till soil")
	fmt.Println("plant moss")
}
```

#### Verified e13-order — ok

Real stdout:

```text
till soil
plant moss
```

<!-- go-example id=e13-try imports=fmt outcome=ok -->
```go
func main() {
	fmt.Println("prepare")
	fmt.Println("plant")
}
```

#### Verified e13-try — ok

Real stdout:

```text
prepare
plant
```

**Playground idea:** Use a row of tiles with ground-state and plant-state badges. Each step in the queue reveals which tile changed and whether its ground is ready.
