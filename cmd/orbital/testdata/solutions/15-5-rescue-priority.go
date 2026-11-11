package main

import (
	"container/heap"
	"orbital/colony"
)

type RescueHeap []colony.RescueRequest

func (q RescueHeap) Len() int { return len(q) }
func (q RescueHeap) Less(i, j int) bool {
	if q[i].Priority != q[j].Priority {
		return q[i].Priority > q[j].Priority
	}
	return q[i].ID < q[j].ID
}
func (q RescueHeap) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *RescueHeap) Push(v any)   { *q = append(*q, v.(colony.RescueRequest)) }
func (q *RescueHeap) Pop() any {
	old := *q
	last := old[len(old)-1]
	*q = old[:len(old)-1]
	return last
}
func main() {
	q := &RescueHeap{}
	for _, request := range colony.RescueArrivals() {
		heap.Push(q, request)
	}
	served := 0
	for q.Len() > 0 {
		request := heap.Pop(q).(colony.RescueRequest)
		if !colony.Rescue(request.ID) {
			return
		}
		served++
		for _, request := range colony.RescueArrivals() {
			heap.Push(q, request)
		}
	}
	colony.Report(served)
}
