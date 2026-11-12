package main

import (
	"reflect"
	"testing"
)

func TestRelayRoutes(t *testing.T) {
	for _, tc := range []struct {
		name        string
		graph       map[int][]int
		start, goal int
		want        []int
		ok          bool
	}{
		{"same", nil, 7, 7, []int{7}, true},
		{"missing", map[int][]int{1: {2}, 2: {1}}, 1, 9, nil, false},
		{"shortest", map[int][]int{1: {2, 3}, 2: {4}, 4: {5}, 3: {5}}, 1, 5, []int{1, 3, 5}, true},
		{"tie", map[int][]int{1: {3, 2}, 3: {4}, 2: {4}}, 1, 4, []int{1, 3, 4}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := RelayPath(tc.graph, tc.start, tc.goal)
			if ok != tc.ok || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("%v %v", got, ok)
			}
		})
	}
}
func BenchmarkRelayRoutes(b *testing.B) {
	graph := map[int][]int{1: {2, 3}, 2: {4}, 3: {4}}
	for i := 0; i < b.N; i++ {
		RelayPath(graph, 1, 4)
	}
}
