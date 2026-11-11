package main

import "errors"

type archiveNode[T any] struct {
	record      T
	left, right *archiveNode[T]
}
type Archive[T any] struct {
	root    *archiveNode[T]
	compare func(T, T) int
}

func (a *Archive[T]) Order(compare func(T, T) int) error {
	if compare == nil || a.root != nil {
		return errors.New("order needs an empty archive and comparator")
	}
	a.compare = compare
	return nil
}
func (a *Archive[T]) Store(record T) (bool, error) {
	if a.compare == nil {
		return false, errors.New("configure archive order first")
	}
	var added bool
	a.root, added = storeArchive(a.root, record, a.compare)
	return added, nil
}
func (a *Archive[T]) Contains(record T) bool {
	for node := a.root; node != nil; {
		c := a.compare(record, node.record)
		if c == 0 {
			return true
		}
		if c < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}
	return false
}
func (a *Archive[T]) Records() []T { return archiveRecords(a.root) }
func storeArchive[T any](node *archiveNode[T], record T, compare func(T, T) int) (*archiveNode[T], bool) {
	if node == nil {
		return &archiveNode[T]{record: record}, true
	}
	c := compare(record, node.record)
	if c == 0 {
		return node, false
	}
	var added bool
	if c < 0 {
		node.left, added = storeArchive(node.left, record, compare)
	} else {
		node.right, added = storeArchive(node.right, record, compare)
	}
	return node, added
}
func archiveRecords[T any](node *archiveNode[T]) []T {
	if node == nil {
		return nil
	}
	records := archiveRecords(node.left)
	records = append(records, node.record)
	for _, record := range archiveRecords(node.right) {
		records = append(records, record)
	}
	return records
}
func main() {}
