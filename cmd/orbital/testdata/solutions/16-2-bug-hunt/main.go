package main

import (
	"fmt"
	"slices"
)

func Median(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	s := slices.Clone(xs)
	slices.Sort(s)
	m := len(s) / 2
	if len(s)%2 == 0 {
		return (s[m-1] + s[m]) / 2
	}
	return s[m]
}

func main() {
	readings := []int{7, 3, 9}
	fmt.Println(Median(readings), readings)
}
