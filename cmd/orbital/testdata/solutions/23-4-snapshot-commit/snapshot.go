package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/fs"
	"slices"
	"strings"
)

type SupplyIndex map[string]string
type SupplyStamp struct {
	Tick    int64
	Message string
}
type SupplyCommit struct {
	Tree    string
	Parent  string
	Tick    int64
	Message string
}

func supplyName(name string) bool {
	return name != "." && len(name) <= 255 && fs.ValidPath(name) && !strings.ContainsAny(name, "\\\x00\t\r\n")
}
func supplyConflict(index SupplyIndex, name string) bool {
	for old := range index {
		if strings.HasPrefix(name, old+"/") || strings.HasPrefix(old, name+"/") {
			return true
		}
	}
	return false
}
func StageSupply(root string, index SupplyIndex, name string, payload []byte) error {
	if index == nil || !supplyName(name) || supplyConflict(index, name) {
		return fmt.Errorf("invalid staging name or index")
	}
	if _, exists := index[name]; !exists && len(index) >= 128 {
		return fmt.Errorf("index full")
	}
	id, err := WriteVault(root, "cargo", payload)
	if err != nil {
		return err
	}
	index[name] = id
	return nil
}
func supplyStamp(stamp SupplyStamp) bool {
	return stamp.Tick >= 0 && len(stamp.Message) <= 512 && strings.TrimSpace(stamp.Message) != ""
}
func supplyParentID(id string) bool {
	if id == "" {
		return true
	}
	if len(id) != 40 {
		return false
	}
	for _, c := range id {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func supplyParentTree(root, parent string) (string, error) {
	if parent == "" {
		return "", nil
	}
	kind, encoded, err := ReadVault(root, parent)
	if err != nil {
		return "", err
	}
	if kind != "snapshot" {
		return "", fmt.Errorf("parent is not a snapshot")
	}
	var commit SupplyCommit
	if err := json.Unmarshal(encoded, &commit); err != nil {
		return "", err
	}
	canonical, err := json.Marshal(commit)
	if err != nil {
		return "", err
	}
	if !bytes.Equal(canonical, encoded) || !supplyStamp(SupplyStamp{Tick: commit.Tick, Message: commit.Message}) || !supplyParentID(commit.Parent) {
		return "", fmt.Errorf("invalid parent snapshot")
	}
	kind, _, err = ReadVault(root, commit.Tree)
	if err != nil {
		return "", err
	}
	if kind != "tree" {
		return "", fmt.Errorf("parent tree kind")
	}
	return commit.Tree, nil
}
func supplyTree(root string, index SupplyIndex) ([]byte, error) {
	if len(index) == 0 || len(index) > 128 {
		return nil, fmt.Errorf("invalid index size")
	}
	names := make([]string, 0, len(index))
	for name := range index {
		if !supplyName(name) || supplyConflict(index, name) {
			return nil, fmt.Errorf("invalid index name")
		}
		names = append(names, name)
	}
	slices.Sort(names)
	var tree strings.Builder
	for _, name := range names {
		kind, _, err := ReadVault(root, index[name])
		if err != nil {
			return nil, err
		}
		if kind != "cargo" {
			return nil, fmt.Errorf("index points to non-cargo")
		}
		fmt.Fprintf(&tree, "%s\t%s\n", name, index[name])
	}
	return []byte(tree.String()), nil
}
func CommitSupply(root string, index SupplyIndex, parent string, stamp SupplyStamp) (string, error) {
	if !supplyStamp(stamp) {
		return "", fmt.Errorf("invalid snapshot metadata")
	}
	tree, err := supplyTree(root, index)
	if err != nil {
		return "", err
	}
	previous, err := supplyParentTree(root, parent)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	fmt.Fprintf(h, "OC23 tree %d\n", len(tree))
	h.Write(tree)
	if fmt.Sprintf("%x", h.Sum(nil)) == previous {
		return "", fmt.Errorf("unchanged tree")
	}
	treeID, err := WriteVault(root, "tree", tree)
	if err != nil {
		return "", err
	}
	encoded, err := json.Marshal(SupplyCommit{Tree: treeID, Parent: parent, Tick: stamp.Tick, Message: stamp.Message})
	if err != nil {
		return "", err
	}
	return WriteVault(root, "snapshot", encoded)
}
