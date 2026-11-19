package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Live colony debugger. The colony program never ends, so instead of a
// recording it runs an instrumented build that streams the statement it is
// executing (sampled) and can be paused, stepped and continued at the next
// statement — a real debugger on the live program. The plain build still
// judges compile errors; race-mode runs stay uninstrumented (the recorder's
// locks would hide races).

const colonyDebugProg = "prog-debug"

// buildColonyDebug builds the instrumented twin of a colony project into
// dir/prog-debug. It runs next to the plain build; any failure only means
// the run has no live debugger.
func (r *runner) buildColonyDebug(ctx context.Context, dir string, project *colonyProject) <-chan bool {
	done := make(chan bool, 1)
	if project.Race {
		done <- false
		return done
	}
	goFiles := map[string]string{}
	for name, body := range project.Files {
		if strings.HasSuffix(name, ".go") {
			goFiles[name] = body
		}
	}
	rewritten, recorder, err := rewriteTraceFiles(goFiles)
	if err != nil {
		done <- false
		return done
	}
	debugDir := filepath.Join(dir, "debug-src")
	files := map[string]string{}
	for name, body := range project.Files {
		files[name] = body
	}
	for name, body := range rewritten {
		files[name] = body
	}
	if err := os.MkdirAll(debugDir, 0o755); err != nil {
		done <- false
		return done
	}
	if err := r.writeModule(debugDir, colonyControlLevel()); err != nil {
		done <- false
		return done
	}
	if err := writeProjectFiles(debugDir, files); err != nil {
		done <- false
		return done
	}
	if err := os.WriteFile(filepath.Join(debugDir, recorder.name), []byte(recorder.text), 0o644); err != nil {
		done <- false
		return done
	}
	go func() {
		build := r.goCmd(ctx, debugDir, "build", "-o", filepath.Join(dir, colonyDebugProg), ".")
		out := &limitedBuffer{max: 16 << 10}
		build.Stdout, build.Stderr = out, out
		done <- runCommand(build) == nil
	}()
	return done
}

func colonyDebugAvailable(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, colonyDebugProg))
	return err == nil
}

// colonyDebug is the parent's side of the two debugger pipes. The pipes are
// recreated when the control loop restarts the child, so the struct is
// shared by the run and its replacements.
type colonyDebug struct {
	mu     sync.Mutex
	ctl    *os.File // parent writes commands; child reads fd ORBITAL_DEBUG_CTL_FD
	breaks string   // last breakpoint list, re-sent to restarted children
	paused bool

	sendMu   sync.Mutex
	lastSent time.Time
	pending  string
	flush    *time.Timer
}

// colonyDebugEnv puts the two debugger fds after the files the colony
// already passes (events fd 3, commands fd 4).
func colonyDebugEnv(extra int) []string {
	return []string{
		fmt.Sprintf("ORBITAL_DEBUG_FD=%d", 3+extra),
		fmt.Sprintf("ORBITAL_DEBUG_CTL_FD=%d", 4+extra),
	}
}

// attach adds fresh pipes to cmd (before Start). The returned function must
// run after Start: it closes the child's ends and starts reading frames.
func (d *colonyDebug) attach(c *colonyGame, cmd *exec.Cmd) (func(started bool), error) {
	framesR, framesW, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	ctlR, ctlW, err := os.Pipe()
	if err != nil {
		framesR.Close()
		framesW.Close()
		return nil, err
	}
	cmd.ExtraFiles = append(cmd.ExtraFiles, framesW, ctlR)
	return func(started bool) {
		framesW.Close()
		ctlR.Close()
		if !started {
			framesR.Close()
			ctlW.Close()
			return
		}
		d.mu.Lock()
		if d.ctl != nil {
			d.ctl.Close()
		}
		d.ctl = ctlW
		d.paused = false
		if d.breaks != "" {
			fmt.Fprintln(ctlW, d.breaks)
		}
		d.mu.Unlock()
		go d.read(c, framesR)
	}, nil
}

func (d *colonyDebug) read(c *colonyGame, frames *os.File) {
	defer frames.Close()
	scanner := bufio.NewScanner(frames)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	for scanner.Scan() {
		line := scanner.Text()
		var frame struct {
			Paused bool `json:"paused"`
		}
		if json.Unmarshal([]byte(line), &frame) != nil {
			continue
		}
		d.mu.Lock()
		d.paused = frame.Paused
		d.mu.Unlock()
		d.send(c, line, frame.Paused)
	}
}

// send forwards frames to the browser: a paused frame at once, running
// frames at most ~12 per second (always ending on the latest one).
func (d *colonyDebug) send(c *colonyGame, frame string, paused bool) {
	d.sendMu.Lock()
	defer d.sendMu.Unlock()
	if paused || time.Since(d.lastSent) >= 80*time.Millisecond {
		d.pending = ""
		d.lastSent = time.Now()
		c.broadcast(colonyDebugLine(frame))
		return
	}
	d.pending = frame
	if d.flush == nil {
		d.flush = time.AfterFunc(90*time.Millisecond, func() {
			d.sendMu.Lock()
			pending := d.pending
			d.pending, d.flush = "", nil
			if pending != "" {
				d.lastSent = time.Now()
			}
			d.sendMu.Unlock()
			if pending != "" {
				c.broadcast(colonyDebugLine(pending))
			}
		})
	}
}

func colonyDebugLine(frame string) string {
	return `{"op":"debug","frame":` + frame + `}`
}

func (d *colonyDebug) command(line string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if strings.HasPrefix(line, "breaks") {
		d.breaks = line
	}
	if d.ctl == nil {
		return os.ErrClosed
	}
	_, err := fmt.Fprintln(d.ctl, line)
	return err
}

func (d *colonyDebug) close() {
	d.mu.Lock()
	if d.ctl != nil {
		d.ctl.Close()
		d.ctl = nil
	}
	d.mu.Unlock()
}

// colonyDebugCommand: POST /api/colony/debug {"cmd": "pause"|"cont"|"into"|"over"|"out", "breaks": ["main.go:12"]}
func (s *server) colonyDebugCommand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Cmd    string   `json:"cmd"`
		Breaks []string `json:"breaks"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10)).Decode(&req); err != nil {
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	line := ""
	switch req.Cmd {
	case "pause", "cont", "into", "over", "out":
		line = req.Cmd
	case "breaks":
		keys := []string{"breaks"}
		for _, key := range req.Breaks {
			if validBreakpointKey(key) {
				keys = append(keys, key)
			}
		}
		line = strings.Join(keys, " ")
	default:
		requestError(w, r, "bad request", http.StatusBadRequest)
		return
	}
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.mu.Lock()
	debug := (*colonyDebug)(nil)
	if c.active != nil {
		debug = c.active.debug
	}
	c.mu.Unlock()
	if debug == nil {
		requestError(w, r, "no live debugger for this run", http.StatusConflict)
		return
	}
	if err := debug.command(line); err != nil {
		requestError(w, r, "debugger is not connected", http.StatusConflict)
		return
	}
	writeJSON(w, map[string]bool{"ok": true})
}

func validBreakpointKey(key string) bool {
	file, line, ok := strings.Cut(key, ":")
	if !ok || !strings.HasSuffix(file, ".go") || strings.ContainsAny(file, " \n/\\") {
		return false
	}
	for _, r := range line {
		if r < '0' || r > '9' {
			return false
		}
	}
	return line != "" && len(line) < 7
}
