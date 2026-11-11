package main

import "fmt"

func CountPairs(xs []int, target int) int {
	seen := make(map[int]int, len(xs))
	n := 0
	for _, x := range xs {
		n += seen[target-x]
		seen[x]++
	}
	return n
}

func main() {
	fmt.Println(CountPairs([]int{1, 9, 5, 5, 3, 7}, 10))
}
