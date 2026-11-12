package main

import "errors"

type Crate struct {
	Bay   string
	Units int
}
type Clearance interface{ Allow(string) bool }

var ErrManifest = errors.New("invalid evacuation manifest")

func Pack(crates []Crate, permit Clearance) (map[string]int, error) {
	for _, c := range crates {
		if c.Bay == "" || c.Units < 0 {
			return nil, ErrManifest
		}
	}
	totals := make(map[string]int)
	for _, c := range crates {
		if !permit.Allow(c.Bay) {
			return nil, ErrManifest
		}
		totals[c.Bay] += c.Units
	}
	return totals, nil
}
func main() {}
