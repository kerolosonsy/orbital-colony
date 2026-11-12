package main

import (
	"container/heap"
	"errors"
	"time"
)

type Job struct {
	Name     string
	Every    time.Duration
	Priority int
}

type Misfire int

const (
	FireOnce Misfire = iota
	Skip
)

type Run struct {
	Name string
	At   time.Time
}

var ErrJob = errors.New("invalid job")

type entry struct {
	job     Job
	next    time.Time
	running bool
}

type dueQueue []*entry

func (q dueQueue) Len() int { return len(q) }
func (q dueQueue) Less(i, j int) bool {
	a, b := q[i], q[j]
	if !a.next.Equal(b.next) {
		return a.next.Before(b.next)
	}
	if a.job.Priority != b.job.Priority {
		return a.job.Priority > b.job.Priority
	}
	return a.job.Name < b.job.Name
}
func (q dueQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q *dueQueue) Push(x any)   { *q = append(*q, x.(*entry)) }
func (q *dueQueue) Pop() any {
	old := *q
	e := old[len(old)-1]
	*q = old[:len(old)-1]
	return e
}

type Scheduler struct {
	start  time.Time
	policy Misfire
	due    dueQueue
	jobs   map[string]*entry
}

func NewScheduler(start time.Time, policy Misfire) *Scheduler {
	return &Scheduler{start: start, policy: policy, jobs: map[string]*entry{}}
}

func (s *Scheduler) Add(j Job) error {
	if j.Name == "" || j.Every <= 0 || s.jobs[j.Name] != nil {
		return ErrJob
	}
	e := &entry{job: j, next: s.start.Add(j.Every)}
	s.jobs[j.Name] = e
	heap.Push(&s.due, e)
	return nil
}

func (s *Scheduler) Tick(now time.Time) []Run {
	runs := []Run{}
	var later []*entry
	for s.due.Len() > 0 && !s.due[0].next.After(now) {
		e := heap.Pop(&s.due).(*entry)
		missed := now.Sub(e.next) / e.job.Every
		last := e.next.Add(missed * e.job.Every)
		fire := !e.running && (s.policy == FireOnce || missed == 0)
		if fire {
			runs = append(runs, Run{Name: e.job.Name, At: e.next})
			e.running = true
		}
		e.next = last.Add(e.job.Every)
		later = append(later, e)
	}
	for _, e := range later {
		heap.Push(&s.due, e)
	}
	return runs
}

func (s *Scheduler) Done(name string) error {
	e := s.jobs[name]
	if e == nil || !e.running {
		return ErrJob
	}
	e.running = false
	return nil
}

func main() {}
