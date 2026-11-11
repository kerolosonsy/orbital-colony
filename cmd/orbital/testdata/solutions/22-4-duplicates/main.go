package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"slices"
	"strings"
)

func packetSizes(store fs.FS) (map[int64][]string, error) {
	root, err := fs.Stat(store, ".")
	if err != nil {
		return nil, fmt.Errorf("stat .: %w", err)
	}
	if !root.IsDir() {
		return nil, fmt.Errorf("root . is not directory")
	}
	buckets := map[int64][]string{}
	err = fs.WalkDir(store, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk %s: %w", name, err)
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return fmt.Errorf("info %s: %w", name, err)
		}
		buckets[info.Size()] = append(buckets[info.Size()], name)
		return nil
	})
	return buckets, err
}
func packetDigest(store fs.FS, name string, size int64) (string, error) {
	file, err := store.Open(name)
	if err != nil {
		return "", fmt.Errorf("open %s: %w", name, err)
	}
	hash := sha256.New()
	n, copyErr := io.Copy(hash, file)
	closeErr := file.Close()
	if copyErr != nil {
		return "", fmt.Errorf("read %s: %w", name, copyErr)
	}
	if closeErr != nil {
		return "", fmt.Errorf("close %s: %w", name, closeErr)
	}
	if n != size {
		return "", fmt.Errorf("size changed %s", name)
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
func PacketEchoes(store fs.FS) ([][]string, error) {
	buckets, err := packetSizes(store)
	if err != nil {
		return nil, err
	}
	var groups [][]string
	for size, names := range buckets {
		if len(names) < 2 {
			continue
		}
		digests := map[string][]string{}
		for _, name := range names {
			digest, err := packetDigest(store, name, size)
			if err != nil {
				return nil, err
			}
			digests[digest] = append(digests[digest], name)
		}
		for _, group := range digests {
			if len(group) > 1 {
				slices.Sort(group)
				groups = append(groups, group)
			}
		}
	}
	slices.SortFunc(groups, func(a, b []string) int { return strings.Compare(a[0], b[0]) })
	return groups, nil
}
func main() {}
