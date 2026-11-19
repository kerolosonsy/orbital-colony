# Chapter 04 lesson storyboards

Mission content is brief and bilingual in `colony/levels/ch04_lessons.go`. Examples use practice values; the checkpoint lab keeps its private classification details out of the lesson.

## 4-1 — Patience Pays

Follow a false readiness condition into the wait branch, then predict a true condition.

<!-- go-example id=e41-if-branch imports=fmt outcome=ok -->
```go
func main() {
	ready := false
	if ready {
		fmt.Println("harvest")
	} else {
		fmt.Println("wait")
	}
}
```

#### Verified e41-if-branch — ok

Real stdout:

```text
wait
```

The Try it reveal enters the guarded body.

<!-- go-example id=e41-if-reveal imports=fmt outcome=ok -->
```go
func main() {
	ready := true
	if ready {
		fmt.Println("ready")
	}
}
```

#### Verified e41-if-reveal — ok

Real stdout:

```text
ready
```

## 4-2 — Mixed Terrain

Show one matching switch case and how unmatched values reach default.

<!-- go-example id=e42-switch-case imports=fmt outcome=ok -->
```go
func main() {
	signal := 2
	switch signal {
	case 1:
		fmt.Println("hold")
	case 2:
		fmt.Println("go")
	default:
		fmt.Println("check")
	}
}
```

#### Verified e42-switch-case — ok

Real stdout:

```text
go
```

The Try it reveal selects the matching amber case.

<!-- go-example id=e42-switch-reveal imports=fmt outcome=ok -->
```go
func main() {
	color := "amber"
	switch color {
	case "blue":
		fmt.Println("move")
	case "amber":
		fmt.Println("pause")
	default:
		fmt.Println("check")
	}
}
```

#### Verified e42-switch-reveal — ok

Real stdout:

```text
pause
```

## 4-3 — Inspection Checkpoint

Teach ordered classification with unrelated delivery-time examples; the lesson does not provide checkpoint bands or test cases.

<!-- go-example id=e43-ordered-bands imports=fmt outcome=ok -->
```go
func main() {
	hours := 5
	label := "long"
	if hours < 3 {
		label = "short"
	} else if hours < 8 {
		label = "medium"
	}
	fmt.Println(label)
}
```

#### Verified e43-ordered-bands — ok

Real stdout:

```text
medium
```

The Try it reveal checks who owns an exact practice boundary.

<!-- go-example id=e43-boundary-reveal imports=fmt outcome=ok -->
```go
func main() {
	level := 4
	label := "high"
	if level < 4 {
		label = "low"
	} else {
		label = "middle"
	}
	fmt.Println(label)
}
```

#### Verified e43-boundary-reveal — ok

Real stdout:

```text
middle
```
