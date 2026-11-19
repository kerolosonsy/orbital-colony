# Chapter 20 lesson storyboards

All four labs use the approved six-step lesson flow. The examples use illustrative data and demonstrate one concept at a time; none implements a lab contract.

## 20-1 — Colony Bulletin

Mission: render a small bulletin as safe HTML. Idea: recognized delimiters select a state. Try: predict escaping. Watch out: escape user text once and reject an unfinished state. Think: process delimiters, state, and safe text. Your turn: add states gradually.

<!-- go-example id=e201-state imports=fmt,strings outcome=ok -->
```go
func main() {
var source = "[alert]low[/alert]"
mode := "plain"
if strings.HasPrefix(source, "[alert]") {
	mode = "alert"
}
fmt.Println(mode)
}
```

<!-- go-example id=e201-escape imports=fmt,html,strings outcome=ok -->
```go
func main() {
var out strings.Builder
out.WriteString(html.EscapeString("<relay>"))
fmt.Println(out.String())
}
```

## 20-2 — Validators

Mission: check three independent crew fields and aggregate classified causes. Idea: wrapping and joining preserve causes. Try: predict a whole-string regex result. Watch out: anchors check the whole string; a date parser checks calendar validity. Think: check fields separately and report in fixed order. Your turn: try a card with multiple failures.

<!-- go-example id=e202-join imports=errors,fmt outcome=ok -->
```go
func main() {
var errPilot = errors.New("pilot")
var errRadio = errors.New("radio")
combined := errors.Join(fmt.Errorf("card: %w", errPilot), errRadio)
fmt.Println(errors.Is(combined, errPilot), errors.Is(combined, errRadio))
}
```

<!-- go-example id=e202-regex imports=fmt,regexp outcome=ok -->
```go
func main() {
pattern := regexp.MustCompile("^[a-z]+$")
fmt.Println(pattern.MatchString("relay!"))
}
```

## 20-3 — Command Lexer

Mission: turn flight text into tokens with source positions. Idea: retain a token's kind, original spelling, and byte offset. Try: predict byte indexes. Watch out: consume maximal runs and stop on an unknown byte. Think: scan left to right and append EOF. Your turn: check offsets by hand.

<!-- go-example id=e203-token imports=fmt outcome=ok -->
```go
type Token struct {
	Kind, Text string
	Offset int
}
func main() {
	t := Token{Kind: "name", Text: "bay_2", Offset: 4}
	fmt.Println(t.Kind, t.Text, t.Offset)
}
```

<!-- go-example id=e203-offset imports=fmt outcome=ok -->
```go
func main() {
source := "go"
for i := 0; i < len(source); i++ {
fmt.Println(i, string(source[i]))
}
}
```

## 20-4 — Flight Evaluator

Mission: interpret expressions and flight statements. Idea: precedence and parentheses determine grouping. Try: predict left-associative subtraction. Watch out: parse the whole program before movement. Think: parse, evaluate, then call the drive interface. Your turn: start with a small expression.

<!-- go-example id=e204-precedence imports=fmt outcome=ok -->
```go
func main() {
	fmt.Println(2 + 3*4)
}
```

<!-- go-example id=e204-associativity imports=fmt outcome=ok -->
```go
func main() {
fmt.Println(9 - 3 - 1)
}
```

## Toolchain transcripts

#### Verified e201-state — ok

Real stdout:

```text
alert
```

#### Verified e201-escape — ok

Real stdout:

```text
&lt;relay&gt;
```

#### Verified e202-join — ok

Real stdout:

```text
true true
```

#### Verified e202-regex — ok

Real stdout:

```text
false
```

#### Verified e203-token — ok

Real stdout:

```text
name bay_2 4
```

#### Verified e203-offset — ok

Real stdout:

```text
0 g
1 o
```

#### Verified e204-precedence — ok

Real stdout:

```text
14
```

#### Verified e204-associativity — ok

Real stdout:

```text
5
```
