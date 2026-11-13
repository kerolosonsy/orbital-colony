package plan

func Next(ready []int) int {
	best := -1
	for i := range ready {
		if best < 0 || ready[i] < ready[best] {
			best = i
		}
	}
	return best
}
