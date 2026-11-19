# Chapter 03 lesson storyboards

Mission content is brief and bilingual in `colony/levels/ch03_lessons.go`. Walkthroughs use small practice counters and grids; none contains the mission's full route.

## 3-1 — The Long Row

Show the loop counter on each turn, then predict how many repeated prints occur.

<!-- go-example id=e31-loop-counter imports=fmt outcome=ok -->
```go
func main() {
	for i := 1; i <= 3; i++ {
		fmt.Println(i)
	}
}
```

#### Verified e31-loop-counter — ok

Real stdout:

```text
1
2
3
```

The Try it reveal stops when `n` reaches 2.

<!-- go-example id=e31-loop-reveal imports=fmt outcome=ok -->
```go
func main() {
	for n := 0; n < 2; n++ {
		fmt.Println("step")
	}
}
```

#### Verified e31-loop-reveal — ok

Real stdout:

```text
step
step
```

## 3-2 — The Whole Field

Trace the inner counter restarting for each outer row; the practice grid is smaller than the mission.

<!-- go-example id=e32-nested-loops imports=fmt outcome=ok -->
```go
func main() {
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			fmt.Println(row, col)
		}
	}
}
```

#### Verified e32-nested-loops — ok

Real stdout:

```text
0 0
0 1
1 0
1 1
```

The Try it reveal prints one item from each practice row.

<!-- go-example id=e32-nested-reveal imports=fmt outcome=ok -->
```go
func main() {
	for row := 0; row < 2; row++ {
		for col := 0; col < 1; col++ {
			fmt.Println("row", row)
		}
	}
}
```

#### Verified e32-nested-reveal — ok

Real stdout:

```text
row 0
row 1
```

## 3-3 — A Farm That Never Sleeps

Show how the loop condition changes after each body pass, without supplying a field tour.

<!-- go-example id=e33-condition-loop imports=fmt outcome=ok -->
```go
func main() {
	count := 0
	for count < 3 {
		count++
		fmt.Println(count)
	}
}
```

#### Verified e33-condition-loop — ok

Real stdout:

```text
1
2
3
```

The Try it reveal prints the decreasing practice counter until zero.

<!-- go-example id=e33-condition-reveal imports=fmt outcome=ok -->
```go
func main() {
	left := 2
	for left > 0 {
		fmt.Println(left)
		left--
	}
}
```

#### Verified e33-condition-reveal — ok

Real stdout:

```text
2
1
```
