package main

import "strings"

func Assemble(parts []string) string {
	size := 0
	for _, p := range parts {
		size += len(p)
	}
	if len(parts) > 1 {
		size += len(parts) - 1
	}
	var b strings.Builder
	b.Grow(size)
	for i, p := range parts {
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(p)
	}
	return b.String()
}
func CopyPositive(readings []int) []int {
	out := make([]int, 0, len(readings))
	for _, v := range readings {
		if v > 0 {
			out = append(out, v)
		}
	}
	return out
}
func main() {}
