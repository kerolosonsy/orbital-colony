package sim

import (
	"context"
	"fmt"
	"time"
)

type FaultError struct {
	Index   int
	Code    string
	Message string
}

func (e *FaultError) Error() string { return fmt.Sprintf("%s (%d): %s", e.Code, e.Index, e.Message) }

type Question struct {
	Prompt  string
	Choices []string
}

func (w *World) CallsignList() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.Data.Callsigns...)
}
func (w *World) ReadingList() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.Data.Readings...)
}
func (w *World) OrderList() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.Data.Orders...)
}
func (w *World) Diagnose(i int) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if i < 0 || i >= len(w.Data.Faults) {
		return fmt.Errorf("%s: %d", say("رقم الفحص خارج الحدود", "diagnostic index out of range"), i)
	}
	if w.Data.Faults[i] == nil {
		return nil
	}
	fault := *w.Data.Faults[i]
	return &fault
}
func (w *World) Dock() { w.mu.Lock(); defer w.mu.Unlock(); w.Data.Docked = true }
func (w *World) Undock() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.Data.Docked {
		w.Data.Undocks++
	}
	w.Data.Docked = false
}
func (w *World) Drill() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.Data.Docked {
		panic(say("لازم تربط الدرون الأول", "dock the drone first"))
	}
	w.Data.Drills++
	if w.Data.DrillPanic {
		panic(say("عطل في المثقاب", "drill fault"))
	}
}
func (w *World) Exam() []Question {
	w.mu.Lock()
	defer w.mu.Unlock()
	questions := append([]Question(nil), w.Data.Questions...)
	for i := range questions {
		questions[i].Choices = append([]string(nil), questions[i].Choices...)
	}
	return questions
}

// Answer records a zero-based choice; invalid indices leave the answers intact.
func (w *World) Answer(i, choice int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if i < 0 || i >= len(w.Data.Questions) || choice < 0 || choice >= len(w.Data.Questions[i].Choices) {
		return false
	}
	if w.Data.Answers == nil {
		w.Data.Answers = map[int]int{}
	}
	w.Data.Answers[i] = choice
	return true
}
func (w *World) Invoice() (int, int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Data.TotalCents, w.Data.People
}

// Scalar and indexed access keeps early missions independent of slices.
func (w *World) Callsign(i int) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if i < 0 || i >= len(w.Data.Callsigns) {
		return ""
	}
	return w.Data.Callsigns[i]
}
func (w *World) CargoAmount(kind int) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if kind < 0 || kind >= len(w.Data.CargoAmounts) {
		return 0
	}
	return w.Data.CargoAmounts[kind]
}
func (w *World) ReadingCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.Data.SensorReadings)
}
func (w *World) Reading(i int) int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if i < 0 || i >= len(w.Data.SensorReadings) {
		return 0
	}
	return w.Data.SensorReadings[i]
}
func (w *World) DiagnosticCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.Data.Faults)
}
func (w *World) DrillAttempts() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Data.DrillAttempts
}
func (w *World) OrderCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.Data.Orders)
}
func (w *World) Order(i int) string {
	w.mu.Lock()
	defer w.mu.Unlock()
	if i < 0 || i >= len(w.Data.Orders) {
		return ""
	}
	return w.Data.Orders[i]
}

// CommandRecord preserves the context observed at a command boundary.
type CommandRecord struct {
	Stage       string
	Deadline    time.Time
	HasDeadline bool
	Err         error
}

func (w *World) TraceCommand(ctx context.Context, stage string) {
	deadline, ok := ctx.Deadline()
	record := CommandRecord{Stage: stage, Deadline: deadline, HasDeadline: ok, Err: ctx.Err()}
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Data.CommandLog = append(w.Data.CommandLog, record)
}

type RescueRequest struct{ ID, Priority int }

// RescueArrivals returns the wave released by the previous successful rescue.
// Call once initially and once after each rescue; snapshots do not consume it.
func (w *World) RescueArrivals() []RescueRequest {
	w.mu.Lock()
	defer w.mu.Unlock()
	wave := len(w.Data.Rescued)
	if wave >= len(w.Data.RescueWaves) {
		return nil
	}
	return append([]RescueRequest(nil), w.Data.RescueWaves[wave]...)
}
func rescueBefore(a, b RescueRequest) bool {
	return a.Priority > b.Priority || a.Priority == b.Priority && a.ID < b.ID
}

// Rescue rejects out-of-order or unknown requests without changing pending work.
func (w *World) Rescue(id int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.Data.RescuePending) == 0 {
		w.Data.RescueFault = true
		return false
	}
	best := 0
	for i, request := range w.Data.RescuePending {
		if rescueBefore(request, w.Data.RescuePending[best]) {
			best = i
		}
	}
	if w.Data.RescuePending[best].ID != id {
		w.Data.RescueFault = true
		return false
	}
	w.Data.RescuePending = append(w.Data.RescuePending[:best], w.Data.RescuePending[best+1:]...)
	w.Data.Rescued = append(w.Data.Rescued, id)
	wave := len(w.Data.Rescued)
	if wave < len(w.Data.RescueWaves) {
		w.Data.RescuePending = append(w.Data.RescuePending, w.Data.RescueWaves[wave]...)
	}
	return true
}
