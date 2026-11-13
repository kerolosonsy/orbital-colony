package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"orbital/colony/econ"
	"orbital/colony/levels"
)

type colonyRun struct {
	am *archiveMarketRun

	stopSent bool
	cmd      *exec.Cmd
	commands *os.File
	done     chan struct{}
	details  *colonyRunDetails
}

type colonyGame struct {
	ops         sync.Mutex // Serialize run, stop, buy and claim, independently of missions.
	mu          sync.Mutex
	state       *econ.State
	path        string
	lastSave    time.Time
	active      *colonyRun
	subscribers map[chan string]struct{}
}

func newColonyGame(dir string) (*colonyGame, error) {
	c := &colonyGame{path: filepath.Join(dir, "colony.json"), subscribers: make(map[chan string]struct{})}
	body, err := os.ReadFile(c.path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		c.state = econ.New(time.Now().UnixNano())
	case err != nil:
		return nil, err
	default:
		c.state = new(econ.State)
		if err := json.Unmarshal(body, c.state); err != nil {
			return nil, err
		}
		c.state.Normalize()
	}
	if err := c.persist(); err != nil {
		return nil, err
	}
	return c, nil
}

// Call with mu held after initialization so writes cannot replace a newer save.
func (c *colonyGame) persist() error {
	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(c.path), ".colony-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	err = json.NewEncoder(f).Encode(c.state)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	if closeErr != nil {
		return closeErr
	}
	if err := os.Rename(f.Name(), c.path); err != nil {
		return err
	}
	c.lastSave = time.Now()
	return nil
}

func (c *colonyGame) summary(lang string, done map[int]bool) map[string]any {
	type surfaceTile struct {
		econ.Tile
		B string `json:"b,omitempty"`
	}
	// Buildings persist separately; the map renderer uses the live tile field.
	tiles := make([]surfaceTile, len(c.state.Tiles))
	for i, tile := range c.state.Tiles {
		tiles[i].Tile = tile
	}
	for _, b := range c.state.Buildings {
		tiles[b.Y*c.state.Size+b.X].B = b.Kind
	}
	state := struct {
		*econ.State
		Tiles []surfaceTile `json:"tiles"`
	}{c.state, tiles}
	owned, total := c.state.TotalLevels()
	return map[string]any{"systems": c.archiveMarketSummary(), "running": c.active != nil, "state": state, "shop": shopForLang(c.state, lang), "tree": treeForLang(c.state, lang), "contracts": contractsForLang(c.state, lang, done), "prices": c.state.Prices(), "levels": map[string]int{"owned": owned, "total": total}, "starter": levelForLang(levels.Colony, lang).Starter}
}

// chaptersDone reads which chapters the profile finished: every mission
// with at least one star, the same rule the panel uses for lesson rewards.
func (s *server) chaptersDone(r *http.Request) map[int]bool {
	done := map[int]bool{}
	if s.profiles == nil {
		return done
	}
	save, err := readProgress(filepath.Join(s.profiles.dir(requestProfile(r)), "progress.json"))
	if err != nil {
		return done
	}
	finished := completed(save)
	for _, ch := range levels.Chapters {
		n, all := 0, true
		for _, l := range levels.All {
			if l.Chapter == ch.N {
				n++
				all = all && finished[l.ID]
			}
		}
		done[ch.N] = n > 0 && all
	}
	return done
}

// contractRequirements is checkRequirements with alternatives: "a|b" is met
// when either key is.
func contractRequirements(code string, reqs []string, lang string) []string {
	var missing []string
	for _, r := range reqs {
		var first []string
		met := false
		for _, alt := range strings.Split(r, "|") {
			m := checkRequirements(code, []string{alt}, lang)
			if len(m) == 0 {
				met = true
				break
			}
			first = append(first, m...)
		}
		if !met {
			missing = append(missing, strings.Join(first, tr(lang, " أو ", " or ")))
		}
	}
	return missing
}

func (c *colonyGame) subscribe() (chan string, func()) {
	ch := make(chan string, 512)
	c.mu.Lock()
	c.subscribers[ch] = struct{}{}
	c.mu.Unlock()
	return ch, func() { c.mu.Lock(); delete(c.subscribers, ch); c.mu.Unlock() }
}

func (c *colonyGame) broadcast(line string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for ch := range c.subscribers {
		select {
		case ch <- line:
		default:
		}
	}
}

func (c *colonyGame) event(line string) {
	line = c.archiveMarketEvent(line)
	c.observeRunEvent(line)
	var event struct {
		State *econ.State `json:"state"`
		T     int64       `json:"t"`
	}
	if json.Unmarshal([]byte(line), &event) == nil && event.State != nil {
		event.State.Normalize()
		c.mu.Lock()
		c.state = event.State
		if time.Since(c.lastSave) >= 2*time.Second {
			if err := c.persist(); err != nil {
				log.Printf("colony save: %v", err)
			}
		}
		c.mu.Unlock()
	}
	c.signalCleanStop(event.T)
	if !c.deferLedgerCompletion(line) {
		c.broadcast(line)
	}
}

func (c *colonyGame) readLines(reader io.Reader, op string) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 4096), 8<<20)
	for scanner.Scan() {
		line := scanner.Text()
		if op == "stderr" {
			c.detectRace(line)
		}
		if op == "" {
			c.event(line)
		} else {
			encoded, _ := json.Marshal(map[string]string{"op": op, "msg": line})
			c.broadcast(string(encoded))
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("colony %s stream: %v", op, err)
	}
}

// Caller holds ops; the exit goroutine only needs mu and can finish the save.
func (c *colonyGame) stop() {
	c.mu.Lock()
	active := c.active
	c.mu.Unlock()
	if active == nil {
		return
	}
	_ = active.cmd.Process.Signal(syscall.SIGTERM)
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()
	select {
	case <-active.done:
	case <-timer.C:
		_ = killProcessTree(active.cmd)
		<-active.done
	}
}

func (c *colonyGame) start(dir string, lang string, details *colonyRunDetails) error {
	savePath := filepath.Join(dir, "colony.json")
	c.mu.Lock()
	body, err := json.Marshal(c.state)
	c.mu.Unlock()
	if err != nil {
		return err
	}
	if err := os.WriteFile(savePath, body, 0o600); err != nil {
		return err
	}
	eventsReader, eventsWriter, err := os.Pipe()
	if err != nil {
		return err
	}
	defer eventsWriter.Close()
	commandsReader, commandsWriter, err := os.Pipe()
	if err != nil {
		eventsReader.Close()
		return err
	}
	defer commandsReader.Close()
	cmd := exec.CommandContext(context.Background(), filepath.Join(dir, "prog"))
	configureProcess(cmd)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "ORBITAL_LANG="+lang, "ORBITAL_LEVEL=colony", "ORBITAL_COLONY="+savePath, "ORBITAL_EVENTS_FD=3", "ORBITAL_CMD_FD=4")
	archiveMarket, err := c.startArchiveMarket(cmd, commandsWriter, lang)
	if err != nil {
		eventsReader.Close()
		commandsWriter.Close()
		return err
	}
	defer archiveMarket.closeUnstarted(cmd)
	cmd.ExtraFiles = []*os.File{eventsWriter, commandsReader}
	// Owned pipes let Wait reap the child without closing unread output.
	stdoutReader, stdoutWriter, err := os.Pipe()
	if err != nil {
		eventsReader.Close()
		commandsWriter.Close()
		return err
	}
	defer stdoutWriter.Close()
	stderrReader, stderrWriter, err := os.Pipe()
	if err != nil {
		eventsReader.Close()
		commandsWriter.Close()
		stdoutReader.Close()
		return err
	}
	defer stderrWriter.Close()
	cmd.Stdout, cmd.Stderr = stdoutWriter, stderrWriter
	if err := cmd.Start(); err != nil {
		eventsReader.Close()
		commandsWriter.Close()
		stdoutReader.Close()
		stderrReader.Close()
		return err
	}
	active := &colonyRun{am: archiveMarket, cmd: cmd, commands: commandsWriter, done: make(chan struct{}), details: details}
	c.mu.Lock()
	c.active = active
	c.beginRun(details)
	c.mu.Unlock()
	archiveMarket.publish()
	var readers sync.WaitGroup
	for _, stream := range []struct {
		file *os.File
		op   string
	}{{eventsReader, ""}, {stdoutReader, "stdout"}, {stderrReader, "stderr"}} {
		readers.Add(1)
		go func(f *os.File, op string) { defer readers.Done(); defer f.Close(); c.readLines(f, op) }(stream.file, stream.op)
	}
	go func() {
		waitErr := cmd.Wait()
		active.am.close()
		_ = killProcessTree(cmd) // Descendants must not keep the stream pipes open.
		readers.Wait()
		commandsWriter.Close()
		c.mu.Lock()
		archiveEvent := c.finishArchiveMarket(active)
		contractEvent := c.finishColonyEvents(lang)
		finalEvent := c.finishRun(active)
		if err := c.persist(); err != nil {
			log.Printf("colony final save: %v", err)
		}
		c.active = nil
		c.mu.Unlock()
		if archiveEvent != "" {
			c.broadcast(archiveEvent)
		}
		if contractEvent != "" {
			c.broadcast(contractEvent)
		}
		if finalEvent != "" {
			c.broadcast(finalEvent)
		}
		reason := tr(lang, "البرنامج خلص", "program exited")
		if waitErr != nil {
			reason = waitErr.Error()
		}
		stopped, _ := json.Marshal(map[string]any{"op": "stopped", "msg": reason, "code": cmd.ProcessState.ExitCode()})
		c.broadcast(string(stopped))
		os.RemoveAll(dir)
		close(active.done)
	}()
	return nil
}

func (s *server) colonyStatus(w http.ResponseWriter, r *http.Request) {
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	writeJSON(w, c.summary(requestLang(r), s.chaptersDone(r)))
}

func (s *server) colonyRun(w http.ResponseWriter, r *http.Request) {
	s.runColonyProject(w, r)
}

func (s *server) colonyStop(w http.ResponseWriter, r *http.Request) {
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.ops.Lock()
	defer c.ops.Unlock()
	c.stop()
	s.colonyStatus(w, r)
}

func (s *server) colonyBuy(w http.ResponseWriter, r *http.Request) {
	lang := requestLang(r)
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
		requestError(w, r, "bad request", 400)
		return
	}
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.ops.Lock()
	defer c.ops.Unlock()
	c.mu.Lock()
	active := c.active
	c.mu.Unlock()
	if active != nil {
		command, _ := json.Marshal(map[string]string{"buy": req.ID})
		if err := active.commands.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
			requestError(w, r, err.Error(), 500)
			return
		}
		if _, err := active.commands.Write(append(command, '\n')); err != nil {
			requestError(w, r, err.Error(), 500)
			return
		}
		writeJSON(w, map[string]any{"ok": true, "pending": true})
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.state.Buy(req.ID)
	message := ""
	if err != nil {
		message = economyText(requestLang(r), err.Error())
	} else if err = c.persist(); err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"ok": err == nil, "error": message, "state": c.state, "shop": shopForLang(c.state, lang), "tree": treeForLang(c.state, lang)})
}

func (s *server) colonyClaim(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Chapter int `json:"chapter"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
		requestError(w, r, "bad request", 400)
		return
	}
	if req.Chapter < 1 || req.Chapter > len(levels.Chapters) {
		requestError(w, r, "unknown chapter", 400)
		return
	}
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.ops.Lock()
	defer c.ops.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active != nil {
		requestError(w, r, "Stop the colony first", http.StatusConflict)
		return
	}
	reward, err := c.state.Claim(req.Chapter)
	message := ""
	if err != nil {
		message = economyText(requestLang(r), err.Error())
	} else if err = c.persist(); err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"ok": err == nil, "error": message, "reward": reward, "state": c.state})
}

// colonyContract accepts a contract, or cancels the active one when id is
// empty. Only while the colony is stopped, like claims.
func (s *server) colonyContract(w http.ResponseWriter, r *http.Request) {
	lang := requestLang(r)
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil {
		requestError(w, r, "bad request", 400)
		return
	}
	done := s.chaptersDone(r)
	s.colonyOps.Lock()
	defer s.colonyOps.Unlock()
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	c.ops.Lock()
	defer c.ops.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active != nil {
		requestError(w, r, "Stop the colony first", http.StatusConflict)
		return
	}
	var err error
	if contract := econ.FindContract(req.ID); req.ID != "" && contract != nil && !done[contract.Chapter] {
		err = errors.New(tr(lang, fmt.Sprintf("خلّص الفصل %d الأول", contract.Chapter), fmt.Sprintf("Finish chapter %d first", contract.Chapter)))
	} else {
		err = c.state.Accept(req.ID)
	}
	message := ""
	if err != nil {
		message = economyText(lang, err.Error())
	} else if err = c.persist(); err != nil {
		requestError(w, r, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]any{"ok": err == nil, "error": message, "state": c.state, "contracts": contractsForLang(c.state, lang, done)})
}

func (s *server) colonyStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		requestError(w, r, "streaming unavailable", 500)
		return
	}
	c := s.requestGame(w, r)
	if c == nil {
		return
	}
	ch, unsubscribe := c.subscribe()
	defer unsubscribe()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case line := <-ch:
			if _, err := fmt.Fprintf(w, "data: %s\n\n", localizeStreamLine(line, requestLang(r))); err != nil {
				return
			}
			flusher.Flush()
		case <-ticker.C:
			if _, err := io.WriteString(w, ": ping\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
