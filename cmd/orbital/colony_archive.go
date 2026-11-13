package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"orbital/colony/econ"
)

type marketReceipt struct {
	ID      string  `json:"id"`
	OK      bool    `json:"ok"`
	Credits float64 `json:"credits"`
	Error   string  `json:"error"`
}

type archiveMarketRun struct {
	mu        sync.Mutex
	writeMu   sync.Mutex
	pending   map[string]chan marketReceipt
	sequence  atomic.Uint64
	ctx       context.Context
	cancel    context.CancelFunc
	commands  *os.File
	market    *colonyMarket
	game      *colonyGame
	archiveOn bool
	dir       string
	rows      int
	journal   bool
	killAt    int64
	killed    bool
	seed      int64
	lang      string
}

func systemEnvironment(env []string) []string {
	clean := make([]string, 0, len(env))
	for _, setting := range env {
		if !strings.HasPrefix(setting, "ORBITAL_COLONY_DIR=") && !strings.HasPrefix(setting, "ORBITAL_MARKET_URL=") {
			clean = append(clean, setting)
		}
	}
	return clean
}

func (c *colonyGame) startArchiveMarket(cmd *exec.Cmd, commands *os.File, lang string) (*archiveMarketRun, error) {
	ctx, cancel := context.WithCancel(context.Background())
	systems := &archiveMarketRun{ctx: ctx, cancel: cancel, commands: commands, game: c, pending: map[string]chan marketReceipt{}, lang: lang}
	cmd.Env = systemEnvironment(cmd.Env)
	c.mu.Lock()
	archive, market := c.state.Has("sys.archive"), c.state.Has("sys.market")
	systems.seed = c.state.Seed
	systems.journal = archive && c.state.Active == "d22-atomic-journal"
	systems.killAt = econ.AtomicJournalAt(c.state.Seed, c.state.Tick)
	stubborn := c.state.Active == "h33-stubborn-client"
	c.mu.Unlock()
	if archive {
		systems.dir = filepath.Join(filepath.Dir(c.path), "colony-data")
		if err := systems.prepareArchive(); err != nil {
			cancel()
			return nil, err
		}
		systems.archiveOn, systems.rows = true, archiveRowCount(systems.dir)
		cmd.Env = append(cmd.Env, "ORBITAL_COLONY_DIR="+systems.dir)
	}
	if market {
		systems.market = &colonyMarket{entries: map[string]*marketEntry{}, seed: systems.seed, failure: 10, lang: lang, run: systems, prices: c.marketPrices}
		if stubborn {
			systems.market.failure = 30
		}
		if err := systems.market.start(); err != nil {
			cancel()
			return nil, marketStartError(lang, err)
		}
		cmd.Env = append(cmd.Env, "ORBITAL_MARKET_URL="+systems.market.url)
	}
	return systems, nil
}

func (c *colonyGame) marketPrices() map[string]float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state.Prices()
}

func journalContents(seed int64, version string) []byte {
	payload := []byte(fmt.Sprintf("version:%s\nseed:%d\n", version, seed) + strings.Repeat("drone:harvest,unload,export\n", 4096))
	checksum := sha256.Sum256(payload)
	return append(payload, []byte(fmt.Sprintf("sha256:%x\n", checksum))...)
}

func (s *archiveMarketRun) prepareArchive() error {
	dir, err := filepath.Abs(s.dir)
	if err != nil {
		return err
	}
	s.dir = dir
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return err
	}
	if !s.journal {
		return nil
	}
	if err := os.Remove(filepath.Join(s.dir, "plan.txt")); err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, version := range []string{"old", "new"} {
		if err := os.WriteFile(filepath.Join(s.dir, "plan-"+version+".txt"), journalContents(s.seed, version), 0o600); err != nil {
			return err
		}
	}
	return nil
}

func (s *archiveMarketRun) closeUnstarted(cmd *exec.Cmd) {
	if cmd.Process == nil {
		s.close()
	}
}

func (s *archiveMarketRun) close() {
	s.cancel()
	if s.market != nil {
		s.market.close()
	}
}

func (s *archiveMarketRun) transact(key string, order marketOrder) (marketReceipt, error) {
	id := marketCommandID(s.sequence.Add(1))
	reply := make(chan marketReceipt, 1)
	s.mu.Lock()
	s.pending[id] = reply
	s.mu.Unlock()
	defer func() { s.mu.Lock(); delete(s.pending, id); s.mu.Unlock() }()
	command, _ := json.Marshal(map[string]any{"id": id, "take": map[string]float64{order.Good: order.Amount}, "market": key})
	s.writeMu.Lock()
	err := s.commands.SetWriteDeadline(time.Now().Add(time.Second))
	if err == nil {
		_, err = s.commands.Write(append(command, '\n'))
	}
	s.writeMu.Unlock()
	if err != nil {
		return marketReceipt{}, err
	}
	select {
	case receipt := <-reply:
		return receipt, nil
	case <-s.ctx.Done():
		return marketReceipt{}, s.ctx.Err()
	}
}

// Called with c.mu held, including from summary and final run settlement.
func (c *colonyGame) archiveMarketSummary() map[string]any {
	own := map[string]any{}
	if c.state.Has("sys.archive") {
		own["archive"] = map[string]any{"on": false, "rows": archiveRowCount(filepath.Join(filepath.Dir(c.path), "colony-data"))}
	}
	if c.state.Has("sys.market") {
		own["market"] = map[string]any{"on": false, "orders": 0, "rejected": 0}
	}
	if c.active != nil && c.active.am != nil {
		s := c.active.am
		if s.dir != "" {
			own["archive"] = map[string]any{"on": s.archiveOn, "rows": s.rows}
		}
		if s.market != nil {
			own["market"] = s.market.summary()
		}
	}
	return own
}

func (s *archiveMarketRun) publish() {
	s.game.mu.Lock()
	own := s.game.archiveMarketSummary()
	s.game.mu.Unlock()
	if len(own) == 0 {
		return
	}
	line, _ := json.Marshal(map[string]any{"op": "systems", "systems": own})
	s.game.broadcast(string(line))
}

func (c *colonyGame) archiveMarketEvent(line string) string {
	var ev struct {
		Op      string                     `json:"op"`
		T       int64                      `json:"t"`
		Receipt *marketReceipt             `json:"receipt"`
		Systems map[string]json.RawMessage `json:"systems"`
	}
	if json.Unmarshal([]byte(line), &ev) != nil {
		return line
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active == nil || c.active.am == nil {
		return line
	}
	s := c.active.am
	if ev.Receipt != nil {
		s.mu.Lock()
		if reply := s.pending[ev.Receipt.ID]; reply != nil {
			select {
			case reply <- *ev.Receipt:
			default:
			}
		}
		s.mu.Unlock()
	}
	if s.journal && !s.killed && ev.T >= s.killAt {
		s.killed = true
		_ = killProcessTree(c.active.cmd)
	}
	if ev.Op == "systems" {
		var archive struct {
			Rows int
			On   bool
		}
		if json.Unmarshal(ev.Systems["archive"], &archive) == nil {
			s.rows, s.archiveOn = archive.Rows, archive.On
		}
		for key, status := range c.archiveMarketSummary() {
			ev.Systems[key], _ = json.Marshal(status)
		}
		encoded, _ := json.Marshal(map[string]any{"op": "systems", "systems": ev.Systems})
		return string(encoded)
	}
	return line
}

func (c *colonyGame) finishArchiveMarket(active *colonyRun) string {
	s := active.am
	if s == nil || !s.journal || c.state.Active != "d22-atomic-journal" {
		return ""
	}
	plan, err := os.ReadFile(filepath.Join(s.dir, "plan.txt"))
	correct := s.killed && err == nil && (bytes.Equal(plan, journalContents(s.seed, "old")) || bytes.Equal(plan, journalContents(s.seed, "new")))
	c.state.RecordSubmission("atomic-journal", correct)
	if correct {
		_, err = c.state.Settle(c.state.Tick)
		correct = err == nil
	}
	if !correct {
		c.state.Fail()
	}
	msg := tr(s.lang, "السجل الذرّي سليم", "Atomic journal is complete")
	if !correct {
		msg = tr(s.lang, "السجل ناقص أو البرنامج خلص قبل الاختبار", "Journal is incomplete or program exited before the test")
	}
	line, _ := json.Marshal(map[string]any{"op": "contract", "ok": correct, "msg": msg, "state": c.state, "contract": map[string]any{"id": "d22-atomic-journal"}})
	return string(line)
}

func archiveRowCount(dir string) int {
	f, err := os.Open(filepath.Join(dir, "ledger.csv"))
	if err != nil {
		return 0
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	_, _ = r.Read()
	rows := 0
	for {
		row, err := r.Read()
		if err == io.EOF || err != nil && len(row) == 0 {
			return rows
		}
		rows++
	}
}
