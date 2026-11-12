package main

import "iter"

func Readings(values []int, visit func(int)) iter.Seq[int] {
	return func(yield func(int) bool) {
		for _, v := range values {
			visit(v)
			if !yield(v) {
				return
			}
		}
	}
}
func Indexed(values []int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i, v := range values {
			if !yield(i, v) {
				return
			}
		}
	}
}
func First(values []int, n int) []int {
	if n <= 0 {
		return nil
	}
	var out []int
	for v := range Readings(values, func(int) {}) {
		out = append(out, v)
		if len(out) == n {
			break
		}
	}
	return out
}
func main() {}
