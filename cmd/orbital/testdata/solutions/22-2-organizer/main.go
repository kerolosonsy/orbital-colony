package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type ShelfMove struct{ From, To string }

func shelfDestination(name string) string {
	shelf := "misc"
	switch strings.ToLower(path.Ext(name)) {
	case ".json", ".csv":
		shelf = "records"
	case ".png":
		shelf = "images"
	}
	return path.Join(shelf, path.Base(name))
}
func shelfPlan(store fs.FS, names []string) ([]ShelfMove, error) {
	slices.Sort(names)
	seen := map[string]bool{}
	var moves []ShelfMove
	for _, name := range names {
		target := shelfDestination(name)
		if seen[target] {
			return nil, fmt.Errorf("collision %s", target)
		}
		seen[target] = true
		if err := shelfAvailable(store, target); err != nil {
			return nil, err
		}
		moves = append(moves, ShelfMove{name, target})
	}
	return moves, nil
}
func shelfAvailable(store fs.FS, target string) error {
	if _, err := fs.Stat(store, target); err == nil {
		return fmt.Errorf("occupied %s", target)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	parent, err := fs.Stat(store, path.Dir(target))
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	if !parent.IsDir() {
		return fmt.Errorf("shelf is not directory %s", path.Dir(target))
	}
	entries, err := fs.ReadDir(store, path.Dir(target))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.Name() == path.Base(target) {
			return fmt.Errorf("occupied %s", target)
		}
	}
	return nil
}
func PlanShelves(store fs.FS) ([]ShelfMove, error) {
	root, err := fs.Stat(store, "inbox")
	if err != nil {
		return nil, err
	}
	if !root.IsDir() {
		return nil, fmt.Errorf("inbox is not directory")
	}
	var names []string
	err = fs.WalkDir(store, "inbox", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			names = append(names, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return shelfPlan(store, names)
}
func hostShelfNames(root string) ([]string, error) {
	var names []string
	err := filepath.WalkDir(filepath.Join(root, "inbox"), func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		names = append(names, filepath.ToSlash(relative))
		return nil
	})
	return names, err
}
func OrganizeShelves(root string) (int, error) { return ExecuteShelves(root, os.Rename) }
func ExecuteShelves(root string, rename func(string, string) error) (int, error) {
	store := os.DirFS(root)
	inbox, err := fs.Stat(store, "inbox")
	if err != nil {
		return 0, err
	}
	if !inbox.IsDir() {
		return 0, fmt.Errorf("inbox is not directory")
	}
	names, err := hostShelfNames(root)
	if err != nil {
		return 0, err
	}
	moves, err := shelfPlan(store, names)
	if err != nil {
		return 0, err
	}
	for i, move := range moves {
		target := filepath.Join(root, filepath.FromSlash(move.To))
		if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
			return i, err
		}
		if err := rename(filepath.Join(root, filepath.FromSlash(move.From)), target); err != nil {
			return i, err
		}
	}
	return len(moves), nil
}
func main() {}
