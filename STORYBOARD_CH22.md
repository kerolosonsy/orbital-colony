# Chapter 22 lesson storyboards

Four short lessons are authored for the chapter's four labs. Every walkthrough is a small, runnable example; all example data is illustrative.

## 22-1 — Store Tree

Mission: render an injected filesystem as a sorted tree, with recursion bounded by depth. Idea: ~fs.ReadDir~ returns sorted entries. Try: predict a recursive counter's output. Watch out: stop before reading beyond the depth boundary and preserve I/O errors. Think: identify the stopping case and each error boundary. Your turn: begin with the root.

<!-- go-example id=e221-fs-readdir imports=fmt,io/fs,testing/fstest outcome=ok -->
```go
func main() {
	store := fstest.MapFS{
		"zeta":  &fstest.MapFile{Data: []byte("z")},
		"alpha": &fstest.MapFile{Data: []byte("a")},
	}
	entries, err := fs.ReadDir(store, ".")
	fmt.Println(entries[0].Name(), len(entries), err)
}
```

<!-- go-example id=e221-recursion imports=fmt outcome=ok -->
```go
func count(n int) {
	if n == 0 {
		return
	}
	fmt.Println(n)
	count(n - 1)
}
func main() {
	count(2)
}
```

#### Verified e221-fs-readdir — ok

Real stdout:

```text
alpha 2 <nil>
```

#### Verified e221-recursion — ok

Real stdout:

```text
2
1
```

## 22-2 — The Organizer

Mission: plan crate destinations before moving files. Idea: a plan is reviewable data; sorting keeps its order stable. Try: predict a normalized extension. Watch out: discovery should not mutate the store. Think: separate discovery, classification, validation, and execution. Your turn: use an injected store and a temporary directory.

<!-- go-example id=e222-sort-plan imports=fmt,slices outcome=ok -->
```go
func main() {
	moves := []string{"shelf/z", "shelf/a"}
	slices.Sort(moves)
	fmt.Println(moves)
}
```

<!-- go-example id=e222-extension imports=fmt,path,strings outcome=ok -->
```go
func main() {
	extension := strings.ToLower(path.Ext("bay/NOTE.JSON"))
	fmt.Println(extension)
}
```

#### Verified e222-sort-plan — ok

Real stdout:

```text
[shelf/a shelf/z]
```

#### Verified e222-extension — ok

Real stdout:

```text
.json
```

## 22-3 — Atomic Save

Mission: preserve the last complete packet when generation or publication fails. Idea: stage bytes behind an ~io.Writer~ before publishing. Try: predict the text written to a builder. Watch out: check callback and close errors. Think: follow the temporary file through cleanup and the rename boundary. Your turn: inspect both target and temporary names after failures.

<!-- go-example id=e223-staging-writer imports=fmt,io,strings outcome=ok -->
```go
func main() {
	var staged strings.Builder
	_, err := io.WriteString(&staged, "ready")
	fmt.Println(staged.String(), err)
}
```

<!-- go-example id=e223-writer imports=fmt,io,strings outcome=ok -->
```go
func main() {
	var stage strings.Builder
	_, err := io.WriteString(&stage, "packet")
	fmt.Println(stage.String(), err)
}
```

#### Verified e223-staging-writer — ok

Real stdout:

```text
ready <nil>
```

#### Verified e223-writer — ok

Real stdout:

```text
packet <nil>
```

## 22-4 — Duplicate Files

Mission: detect duplicate content with bounded-memory reads. Idea: narrow candidates by metadata, then stream each candidate through its own hash. Try: predict equal-size and equal-content comparisons. Watch out: size alone is not proof, and every file needs closing. Think: return no partial report after a read failure. Your turn: inspect stable group order.

<!-- go-example id=e224-stream-hash imports=crypto/sha256,fmt,io,strings outcome=ok -->
```go
func main() {
	digest := sha256.New()
	n, err := io.Copy(digest, strings.NewReader("echo"))
	fmt.Println(n, len(digest.Sum(nil)), err)
}
```

<!-- go-example id=e224-size-filter imports=fmt outcome=ok -->
```go
func main() {
	first := []byte("A")
	second := []byte("B")
	fmt.Println(len(first) == len(second), string(first) == string(second))
}
```

#### Verified e224-stream-hash — ok

Real stdout:

```text
4 32 <nil>
```

#### Verified e224-size-filter — ok

Real stdout:

```text
true false
```
