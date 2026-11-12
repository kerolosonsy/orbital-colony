package main

func RelayPath[T comparable](links map[T][]T, start, goal T) ([]T, bool) {
	queue := []T{start}
	seen := map[T]bool{start: true}
	prev := map[T]T{}
	for head := 0; head < len(queue); head++ {
		node := queue[head]
		if node == goal {
			path := []T{node}
			for node != start {
				node = prev[node]
				path = append(path, node)
			}
			for a, b := 0, len(path)-1; a < b; a, b = a+1, b-1 {
				path[a], path[b] = path[b], path[a]
			}
			return path, true
		}
		for _, next := range links[node] {
			if !seen[next] {
				seen[next] = true
				prev[next] = node
				queue = append(queue, next)
			}
		}
	}
	return nil, false
}
func main() {}
