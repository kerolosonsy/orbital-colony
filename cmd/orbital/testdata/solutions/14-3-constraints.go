package main

import "orbital/colony"

type Number interface {
	~int | ~float64
}

func Sum[T Number](xs []T) T {
	var t T
	for _, x := range xs {
		t += x
	}
	return t
}

func main() {
	colony.Report(Sum(colony.Depths()))
	colony.Report(Sum(colony.Weights()))
	colony.Report(Sum(colony.Temperatures()))
}
