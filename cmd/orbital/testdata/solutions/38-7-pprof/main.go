package main

import (
	"io"
	"runtime/pprof"
)

func UniqueTotal(readings []int) int {
	seen := make(map[int]struct{}, len(readings))
	total := 0
	for _, v := range readings {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			total += v
		}
	}
	return total
}
func ProfileSurvey(cpu, heap io.Writer, readings []int) (int, error) {
	if err := pprof.StartCPUProfile(cpu); err != nil {
		return 0, err
	}
	total := UniqueTotal(readings)
	pprof.StopCPUProfile()
	if err := pprof.WriteHeapProfile(heap); err != nil {
		return total, err
	}
	return total, nil
}
func main() {}
