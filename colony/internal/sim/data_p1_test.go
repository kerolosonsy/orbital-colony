package sim

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"
)

func TestRescueArrivalOrderAndRejectedRequests(t *testing.T) {
	w := New(io.Discard, 1)
	w.Data.RescueWaves = [][]RescueRequest{
		{{ID: 7, Priority: 2}, {ID: 3, Priority: 2}},
		{{ID: 9, Priority: 5}},
	}
	w.Data.RescuePending = append([]RescueRequest(nil), w.Data.RescueWaves[0]...)
	arrivals := w.RescueArrivals()
	arrivals[0].Priority = 99
	if w.RescueArrivals()[0].Priority != 2 {
		t.Fatal("arrival snapshots changed pending priorities")
	}
	before := append([]RescueRequest(nil), w.Data.RescuePending...)
	if w.Rescue(7) || !w.Data.RescueFault || !reflect.DeepEqual(w.Data.RescuePending, before) || len(w.Data.Rescued) != 0 {
		t.Fatal("out-of-order rescue must fail without releasing a wave")
	}
	for _, id := range []int{3, 9, 7} {
		if !w.Rescue(id) {
			t.Fatalf("correct next request %d rejected", id)
		}
	}
	if !reflect.DeepEqual(w.Data.Rescued, []int{3, 9, 7}) || len(w.Data.RescuePending) != 0 || len(w.RescueArrivals()) != 0 {
		t.Fatal("new wave did not preempt older lower-priority requests")
	}
	if !w.Data.RescueFault || w.Rescue(7) || w.Rescue(404) {
		t.Fatal("failure must remain recorded; completed and unknown requests cannot be served")
	}
}

func TestCommandTracePreservesDeadlineAndCancellation(t *testing.T) {
	w := New(io.Discard, 1)
	parent, cancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer cancel()
	child, stop := context.WithTimeout(parent, time.Hour)
	defer stop()
	w.TraceCommand(child, "probe.exit")
	record := w.Data.CommandLog[0]
	if record.Stage != "probe.exit" || !record.HasDeadline || !record.Deadline.Equal(time.Unix(1, 0)) || !errors.Is(record.Err, context.DeadlineExceeded) {
		t.Fatalf("lost parent deadline or stop reason: %+v", record)
	}
	canceled, abort := context.WithCancel(context.Background())
	abort()
	w.TraceCommand(canceled, "command.enter")
	record = w.Data.CommandLog[1]
	if record.HasDeadline || !errors.Is(record.Err, context.Canceled) || w.Tick != 0 {
		t.Fatalf("cancellation trace changed deadline, reason, or tick: %+v", record)
	}
}
