package main

func SelectGates[T comparable](rows [][]T, blocked []T) []T {
	denied := make(map[T]struct{}, len(blocked))
	for _, v := range blocked {
		denied[v] = struct{}{}
	}
	seen := make(map[T]struct{})
	var out []T
Rows:
	for _, row := range rows {
	Candidates:
		for _, v := range row {
			if _, found := denied[v]; found {
				continue Rows
			}
			if _, found := seen[v]; found {
				continue Candidates
			}
			seen[v] = struct{}{}
			out = append(out, v)
			break Candidates
		}
	}
	return out
}
func main() {}
