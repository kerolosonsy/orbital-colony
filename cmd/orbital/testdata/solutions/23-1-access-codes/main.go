package main

import (
	"crypto/rand"
	"fmt"
	"io"
)

func accessPick(source io.Reader, n int) (int, error) {
	for {
		var b [1]byte
		var err error
		if source == nil {
			_, err = rand.Read(b[:])
		} else {
			_, err = io.ReadFull(source, b[:])
		}
		if err != nil {
			return 0, err
		}
		if int(b[0]) < 256-256%n {
			return int(b[0]) % n, nil
		}
	}
}
func accessGroups(length int, categories []string) ([]string, string, error) {
	if length < 1 || length > 128 || len(categories) == 0 || length < len(categories) {
		return nil, "", fmt.Errorf("invalid code length or categories")
	}
	alphabets := map[string]string{"lower": "abcdefghijklmnopqrstuvwxyz", "upper": "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "digits": "0123456789", "marks": "!@#$"}
	seen := map[string]bool{}
	var groups []string
	union := ""
	for _, name := range categories {
		alphabet, ok := alphabets[name]
		if !ok || seen[name] {
			return nil, "", fmt.Errorf("invalid category %q", name)
		}
		seen[name] = true
		groups = append(groups, alphabet)
		union += alphabet
	}
	return groups, union, nil
}
func IssueAccess(source io.Reader, length int, categories []string) (string, error) {
	groups, union, err := accessGroups(length, categories)
	if err != nil {
		return "", err
	}
	code := make([]byte, length)
	for i := range code {
		alphabet := union
		if i < len(groups) {
			alphabet = groups[i]
		}
		j, err := accessPick(source, len(alphabet))
		if err != nil {
			return "", err
		}
		code[i] = alphabet[j]
	}
	for i := len(code) - 1; i > 0; i-- {
		j, err := accessPick(source, i+1)
		if err != nil {
			return "", err
		}
		code[i], code[j] = code[j], code[i]
	}
	return string(code), nil
}
func main() {}
