package main

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

func StoreTree(store fs.FS, depth int, out io.Writer) error {
	if depth < 0 {
		return fmt.Errorf("negative depth")
	}
	root, err := fs.Stat(store, ".")
	if err != nil {
		return fmt.Errorf("stat .: %w", err)
	}
	if !root.IsDir() {
		return fmt.Errorf("root . is not a directory")
	}
	if _, err := fmt.Fprintln(out, "store/"); err != nil {
		return fmt.Errorf("write .: %w", err)
	}
	return storeChildren(store, ".", depth, treeOutput{out, 1})
}

type treeOutput struct {
	writer io.Writer
	level  int
}

func storeChildren(store fs.FS, parent string, remaining int, output treeOutput) error {
	if remaining == 0 {
		return nil
	}
	entries, err := fs.ReadDir(store, parent)
	if err != nil {
		return fmt.Errorf("read %s: %w", parent, err)
	}
	for _, entry := range entries {
		name := path.Join(parent, entry.Name())
		suffix := ""
		switch {
		case entry.IsDir():
			suffix = "/"
		case entry.Type()&fs.ModeSymlink != 0:
			suffix = "@"
		case entry.Type() != 0:
			suffix = "!"
		}
		if _, err := fmt.Fprintf(output.writer, "%s%s%s\n", strings.Repeat("  ", output.level), entry.Name(), suffix); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
		if entry.IsDir() {
			if err := storeChildren(store, name, remaining-1, treeOutput{output.writer, output.level + 1}); err != nil {
				return err
			}
		}
	}
	return nil
}
func main() {}
