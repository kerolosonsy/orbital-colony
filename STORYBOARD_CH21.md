# Chapter 21 lesson storyboards

Short bilingual lessons for the four labs. Examples use illustrative data and are executed by the chapter lesson test.

## 21-1 — Station Status

Mission: keep reports stable across machines. Idea: isolate host data and parse flags in a local parser. Try: predict a Boolean flag and remaining argument count. Watch out: distinguish missing environment data from an empty value. Think: identify checks before output. Your turn: use injected values and a writer.

<!-- go-example id=e211-local-flag imports=flag,fmt,io outcome=ok -->
```go
func main() {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	active := fs.Bool("active", false, "show details")
	_ = fs.Parse([]string{"-active"})
	fmt.Println(*active, fs.NArg())
}
```

<!-- go-example id=e211-flag-reveal imports=flag,fmt,io outcome=ok -->
```go
func main() {
	fs := flag.NewFlagSet("probe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	active := fs.Bool("active", false, "show details")
	_ = fs.Parse([]string{"-active"})
	fmt.Println(*active, fs.NArg())
}
```

## 21-2 — Subcommands

Mission: accept one command or a line stream. Idea: use injected readers, writers, and an interface around state. Try: predict what a scanner reads. Watch out: validate tokens before indexing and classify command errors separately from transport failures. Think: separate command source, interpretation, and state change. Your turn: inspect state across inputs.

<!-- go-example id=e212-scan-lines imports=bufio,fmt,strings outcome=ok -->
```go
func main() {
	scanner := bufio.NewScanner(strings.NewReader("check bay"))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
```

<!-- go-example id=e212-scan-reveal imports=bufio,fmt,strings outcome=ok -->
```go
func main() {
	scanner := bufio.NewScanner(strings.NewReader("check bay"))
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
```

## 21-3 — Guess the Depth

Mission: guide a probe with bounded guesses. Idea: pass a choice source into the game so tests can repeat it. Try: predict a selected label. Watch out: validate a source result before using it as an index. Think: identify why fixed sources help tests. Your turn: exercise the game through injected input and randomness.

<!-- go-example id=e213-injected-pick imports=fmt outcome=ok -->
```go
func pick(labels []string, draw func(int) int) string {
	return labels[draw(len(labels))]
}
func main() {
	fixed := func(n int) int { return 1 }
	fmt.Println(pick([]string{"amber", "cyan", "violet"}, fixed))
}
```

<!-- go-example id=e213-pick-reveal imports=fmt outcome=ok -->
```go
func pick(labels []string, draw func(int) int) string {
	return labels[draw(len(labels))]
}
func main() {
	fixed := func(n int) int { return 1 }
	fmt.Println(pick([]string{"amber", "cyan", "violet"}, fixed))
}
```

## 21-4 — Drone Duel

Mission: map text commands to maneuver values. Idea: use a named integer type and constants for a small vocabulary. Try: predict the numeric value of a named constant. Watch out: validate text before conversion. Think: parse first, then reason about ordered pairs. Your turn: separate input handling from round rules.

<!-- go-example id=e214-iota-values imports=fmt outcome=ok -->
```go
type LampMode int
const (
	idle LampMode = iota
	scan
	relay
)
func main() {
	mode := scan
	fmt.Println(int(mode))
}
```

<!-- go-example id=e214-mode-reveal imports=fmt outcome=ok -->
```go
type LampMode int
const (
	idle LampMode = iota
	scan
	relay
)
func main() {
	mode := scan
	fmt.Println(int(mode))
}
```

## Verified example outputs

#### Verified e211-local-flag — ok

Real stdout:

```text
true 0
```

#### Verified e211-flag-reveal — ok

Real stdout:

```text
true 0
```

#### Verified e212-scan-lines — ok

Real stdout:

```text
check bay
```

#### Verified e212-scan-reveal — ok

Real stdout:

```text
check bay
```

#### Verified e213-injected-pick — ok

Real stdout:

```text
cyan
```

#### Verified e213-pick-reveal — ok

Real stdout:

```text
cyan
```

#### Verified e214-iota-values — ok

Real stdout:

```text
1
```

#### Verified e214-mode-reveal — ok

Real stdout:

```text
1
```
