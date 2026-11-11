package main

import (
	"slices"
	"testing"
)

func TestMedian(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{{[]int{7, 3, 9}, 7}, {[]int{4, 1, 3, 2}, 2}, {nil, 0}}
	for _, c := range cases {
		in := slices.Clone(c.in)
		if got := Median(in); got != c.want {
			t.Errorf("Median(%v) = %d, want %d", c.in, got, c.want)
		}
		if !slices.Equal(in, c.in) {
			t.Errorf("input changed: %v", in)
		}
	}
}
