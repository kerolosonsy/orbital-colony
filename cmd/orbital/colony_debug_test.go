package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"orbital/colony/levels"
)

type colonyDebugFrame struct {
	Line      int    `json:"line"`
	Func      string `json:"func"`
	Depth     int    `json:"depth"`
	Event     string `json:"event"`
	Goroutine int    `json:"goroutine"`
	Paused    bool   `json:"paused"`
}

// nextDebugFrame waits for a debug event on the colony stream.
func nextDebugFrame(t *testing.T, ch <-chan string, wantPaused bool, within time.Duration) colonyDebugFrame {
	t.Helper()
	timer := time.NewTimer(within)
	defer timer.Stop()
	for {
		select {
		case line := <-ch:
			var event struct {
				Op    string
				Frame colonyDebugFrame
			}
			if json.Unmarshal([]byte(line), &event) != nil || event.Op != "debug" {
				continue
			}
			if !wantPaused || event.Frame.Paused {
				return event.Frame
			}
		case <-timer.C:
			t.Fatalf("no debug frame (paused=%v) within %s", wantPaused, within)
		}
	}
}

// The live colony program streams the statement it runs, pauses on demand,
// steps one statement at a time and continues.
func TestColonyLiveDebuggerPausesStepsAndContinues(t *testing.T) {
	s := testColony(t)
	c := s.colony
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	body, _ := json.Marshal(map[string]string{"code": levels.Colony.Starter})
	response := colonyRequest(t, s.colonyRun, string(body))
	var run struct {
		Phase   string
		Running bool
	}
	if err := json.Unmarshal(response.Body.Bytes(), &run); err != nil || run.Phase != "run" || !run.Running {
		t.Fatalf("start: %s", response.Body)
	}
	defer c.stop()

	live := nextDebugFrame(t, ch, false, 20*time.Second)
	if live.Line <= 0 || live.Func == "" {
		t.Fatalf("live frame without position: %+v", live)
	}
	if got := colonyRequest(t, s.colonyDebugCommand, `{"cmd":"pause"}`); got.Code != 200 {
		t.Fatalf("pause: %d %s", got.Code, got.Body)
	}
	paused := nextDebugFrame(t, ch, true, 10*time.Second)
	if got := colonyRequest(t, s.colonyDebugCommand, `{"cmd":"into"}`); got.Code != 200 {
		t.Fatalf("step: %d %s", got.Code, got.Body)
	}
	stepped := nextDebugFrame(t, ch, true, 10*time.Second)
	if stepped == paused && stepped.Event == paused.Event {
		t.Fatalf("step did not move: %+v", stepped)
	}
	// while paused the program makes no progress: nothing else arrives
	quiet := time.NewTimer(600 * time.Millisecond)
	for waiting := true; waiting; {
		select {
		case line := <-ch:
			var event struct {
				Op    string
				Frame colonyDebugFrame
			}
			if json.Unmarshal([]byte(line), &event) == nil && event.Op == "debug" {
				t.Fatalf("frame arrived while paused: %s", line)
			}
		case <-quiet.C:
			waiting = false
		}
	}
	if got := colonyRequest(t, s.colonyDebugCommand, `{"cmd":"cont"}`); got.Code != 200 {
		t.Fatalf("continue: %d %s", got.Code, got.Body)
	}
	if again := nextDebugFrame(t, ch, false, 10*time.Second); again.Paused {
		t.Fatalf("still paused after continue: %+v", again)
	}
}

func TestColonyDebugRejectsBadCommands(t *testing.T) {
	s := testColony(t)
	post := func(body string) int {
		recorder := httptest.NewRecorder()
		s.colonyDebugCommand(recorder, httptest.NewRequest(http.MethodPost, "http://localhost/", strings.NewReader(body)))
		return recorder.Code
	}
	if code := post(`{"cmd":"rm -rf"}`); code != 400 {
		t.Fatalf("bad command accepted: %d", code)
	}
	if code := post(`{"cmd":"pause"}`); code != 409 {
		t.Fatalf("pause without a run: %d", code)
	}
	for key, ok := range map[string]bool{"main.go:12": true, "a/b.go:3": false, "main.go:x": false, "main.go:1 cont": false} {
		if validBreakpointKey(key) != ok {
			t.Fatalf("validBreakpointKey(%q) != %v", key, ok)
		}
	}
}
