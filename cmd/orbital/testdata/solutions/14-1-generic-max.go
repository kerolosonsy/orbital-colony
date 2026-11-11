package main

import (
	"cmp"
	"orbital/colony"
)

func Max[T cmp.Ordered](xs []T) T {
	m := xs[0]
	for _, x := range xs[1:] {
		m = max(m, x)
	}
	return m
}

func main() {
	colony.Report(Max(colony.Depths()))
	colony.Report(Max(colony.Temperatures()))
}
