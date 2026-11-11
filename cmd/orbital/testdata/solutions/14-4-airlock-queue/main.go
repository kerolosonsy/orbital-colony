package main

type Queue[T any] struct{ waiting []T }

func (q *Queue[T]) Admit(passenger T) { q.waiting = append(q.waiting, passenger) }
func (q *Queue[T]) Depart() (T, bool) {
	var zero T
	if len(q.waiting) == 0 {
		return zero, false
	}
	passenger := q.waiting[0]
	q.waiting[0] = zero
	q.waiting = q.waiting[1:]
	if len(q.waiting) == 0 {
		q.waiting = nil
	}
	return passenger, true
}
func (q *Queue[T]) Front() (T, bool) {
	if len(q.waiting) > 0 {
		return q.waiting[0], true
	}
	var zero T
	return zero, false
}
func (q *Queue[T]) Len() int { return len(q.waiting) }
func main()                  {}
