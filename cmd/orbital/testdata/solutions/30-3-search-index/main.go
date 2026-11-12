package main

import (
	"errors"
	"sort"
	"strings"
	"unicode"
)

type Dispatch struct{ Code, Title, Body string }
type DispatchHit struct {
	Dispatch Dispatch
	Score    int
}

var ErrSearch = errors.New("invalid search input")

type DispatchIndex struct {
	docs     map[string]Dispatch
	postings map[string]map[string]int
}

func dispatchTerms(text string) []string {
	return strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
}
func BuildDispatchIndex(docs []Dispatch) (*DispatchIndex, error) {
	index := &DispatchIndex{docs: make(map[string]Dispatch), postings: make(map[string]map[string]int)}
	for _, d := range docs {
		if d.Code == "" {
			return nil, ErrSearch
		}
		if _, exists := index.docs[d.Code]; exists {
			return nil, ErrSearch
		}
		index.docs[d.Code] = d
		index.add(d.Code, d.Title, 3)
		index.add(d.Code, d.Body, 1)
	}
	return index, nil
}
func (index *DispatchIndex) add(code, text string, weight int) {
	for _, term := range dispatchTerms(text) {
		if index.postings[term] == nil {
			index.postings[term] = make(map[string]int)
		}
		index.postings[term][code] += weight
	}
}
func (index *DispatchIndex) Search(query string, offset, limit int) ([]DispatchHit, error) {
	terms := dispatchTerms(query)
	if len(terms) == 0 || offset < 0 || limit < 1 || limit > 25 {
		return nil, ErrSearch
	}
	scores := make(map[string]int)
	seen := make(map[string]bool)
	for _, term := range terms {
		if seen[term] {
			continue
		}
		seen[term] = true
		for code, score := range index.postings[term] {
			scores[code] += score
		}
	}
	hits := make([]DispatchHit, 0, len(scores))
	for code, score := range scores {
		hits = append(hits, DispatchHit{index.docs[code], score})
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].Score != hits[j].Score {
			return hits[i].Score > hits[j].Score
		}
		return hits[i].Dispatch.Code < hits[j].Dispatch.Code
	})
	if offset >= len(hits) {
		return []DispatchHit{}, nil
	}
	end := offset + min(limit, len(hits)-offset)
	return hits[offset:end], nil
}
func (index *DispatchIndex) Vocabulary(prefix string, limit int) ([]string, error) {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	for _, r := range prefix {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return nil, ErrSearch
		}
	}
	if limit < 1 || limit > 20 {
		return nil, ErrSearch
	}
	terms := []string{}
	for term := range index.postings {
		if strings.HasPrefix(term, prefix) {
			terms = append(terms, term)
		}
	}
	sort.Strings(terms)
	if len(terms) > limit {
		terms = terms[:limit]
	}
	return terms, nil
}
func main() {}
