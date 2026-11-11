package main

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode"
)

type WordCount struct {
	Word  string
	Count int
}

func TopSignals(in io.Reader, limit int) ([]WordCount, error) {
	if limit < 0 {
		return nil, fmt.Errorf("negative limit")
	}
	if limit == 0 {
		return nil, nil
	}
	counts := make(map[string]int)
	var scanner *bufio.Scanner = bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 4096), (1<<20)+2)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 1<<20 {
			return nil, fmt.Errorf("signal line too long")
		}
		words := strings.FieldsFunc(strings.ToLower(line), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		})
		for _, word := range words {
			counts[word]++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	var ranked []WordCount
	for word, count := range counts {
		ranked = append(ranked, WordCount{word, count})
	}
	slices.SortFunc(ranked, func(a, b WordCount) int {
		if a.Count > b.Count {
			return -1
		}
		if a.Count < b.Count {
			return 1
		}
		return strings.Compare(a.Word, b.Word)
	})
	if len(ranked) > limit {
		ranked = ranked[:limit]
	}
	return ranked, nil
}
func main() {}
