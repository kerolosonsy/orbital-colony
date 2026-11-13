package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"orbital/colony/econ"
)

type neighborOffer struct {
	id     string
	good   string
	amount int
	price  float64
}

type neighborStatus struct {
	On     bool   `json:"on"`
	Addr   string `json:"addr,omitempty"`
	Trades int    `json:"trades"`
	Drops  int    `json:"drops"`
}

type neighborReceipt struct {
	Op string `json:"op"`
	ID string `json:"id"`
	OK bool   `json:"ok"`
}

type neighborConfig struct {
	seed     int64
	high     bool
	prefix   string
	commands *os.File
	publish  func(string)
}

type colonyNeighbor struct {
	mu          sync.Mutex
	tradeMu     sync.Mutex
	listener    net.Listener
	config      neighborConfig
	status      neighborStatus
	conns       map[net.Conn]struct{}
	connections int
	next        int
	offers      map[string]neighborOffer
	paid        map[string]float64
	pending     map[string]chan neighborReceipt
	closed      chan struct{}
	once        sync.Once
	wg          sync.WaitGroup
}

func newColonyNeighbor(config neighborConfig) (*colonyNeighbor, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	n := &colonyNeighbor{listener: listener, config: config, status: neighborStatus{On: true, Addr: listener.Addr().String()}, conns: map[net.Conn]struct{}{}, next: 1, offers: map[string]neighborOffer{}, paid: map[string]float64{}, pending: map[string]chan neighborReceipt{}, closed: make(chan struct{})}
	n.wg.Add(1)
	go n.accept()
	return n, nil
}

func (n *colonyNeighbor) snapshot() neighborStatus { n.mu.Lock(); defer n.mu.Unlock(); return n.status }

func (n *colonyNeighbor) publish() {
	if n.config.publish == nil {
		return
	}
	line, _ := json.Marshal(map[string]any{"op": "systems", "systems": map[string]any{"neighbor": n.snapshot()}})
	n.config.publish(string(line))
}

func (n *colonyNeighbor) close() {
	n.once.Do(func() {
		close(n.closed)
		n.listener.Close()
		n.mu.Lock()
		n.status.On = false
		for conn := range n.conns {
			conn.Close()
		}
		n.mu.Unlock()
	})
	n.wg.Wait()
	n.publish()
}

func (n *colonyNeighbor) accept() {
	defer n.wg.Done()
	for {
		conn, err := n.listener.Accept()
		if err != nil {
			return
		}
		n.mu.Lock()
		if !n.status.On {
			n.mu.Unlock()
			conn.Close()
			return
		}
		n.connections++
		number := n.connections
		n.conns[conn] = struct{}{}
		n.wg.Add(1)
		n.mu.Unlock()
		go n.serve(conn, number)
	}
}

func (n *colonyNeighbor) serve(conn net.Conn, number int) {
	defer n.wg.Done()
	defer func() { conn.Close(); n.mu.Lock(); delete(n.conns, conn); n.mu.Unlock() }()
	if !n.disturb(conn, number) {
		return
	}
	if n.sendOffer(conn) {
		n.readClient(conn)
	}
}

func (n *colonyNeighbor) readClient(conn net.Conn) {
	reader := bufio.NewReaderSize(conn, econ.NeighborLineLimit)
	accepted := map[string]bool{}
	for {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		line, err := reader.ReadSlice('\n')
		if err != nil || len(line) > econ.NeighborLineLimit {
			return
		}
		fields := strings.Fields(string(line))
		if len(fields) != 2 {
			if !n.write(conn, "ERR - bad-command\n") {
				return
			}
			continue
		}
		if !n.reply(conn, fields, accepted) {
			return
		}
	}
}

func (n *colonyNeighbor) reply(conn net.Conn, fields []string, accepted map[string]bool) bool {
	id := fields[1]
	switch fields[0] {
	case "ACCEPT":
		n.mu.Lock()
		_, exists := n.offers[id]
		n.mu.Unlock()
		if !exists {
			return n.write(conn, "ERR "+id+" unknown-offer\n")
		}
		accepted[id] = true
		return true
	case "DECLINE":
		delete(accepted, id)
		if !n.decline(id) {
			return n.write(conn, "ERR "+id+" unknown-offer\n")
		}
		return n.sendOffer(conn)
	case "SHIP":
		reply := n.ship(id, accepted[id])
		if !n.write(conn, reply) {
			return false
		}
		if strings.HasPrefix(reply, "PAID ") {
			delete(accepted, id)
			return n.afterPayment(conn, id)
		}
		return true
	default:
		return n.write(conn, "ERR "+id+" bad-command\n")
	}
}

func (n *colonyNeighbor) afterPayment(conn net.Conn, id string) bool {
	offerNumber, _ := strconv.ParseInt(id, 10, 64)
	if n.config.high || econ.EventHash(n.config.seed, offerNumber, 32)%3 == 0 {
		n.drop()
		return false
	}
	return n.sendOffer(conn)
}

func (n *colonyNeighbor) write(conn net.Conn, line string) bool {
	conn.SetWriteDeadline(time.Now().Add(time.Second))
	_, err := io.WriteString(conn, line)
	return err == nil
}

func (n *colonyNeighbor) disturb(conn net.Conn, number int) bool {
	slots := 9
	if n.config.high {
		slots = 5
	}
	phase := (number - 1 + int(uint64(n.config.seed)%uint64(slots))) % slots
	switch phase {
	case 0:
		n.write(conn, strings.Repeat("X", econ.NeighborLineLimit*8)+"\n")
	case 1:
		if n.write(conn, "GARBAGE ? ?\n") {
			return true
		}
	case 2:
		if n.write(conn, "OFFER stalled") {
			conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			var one [1]byte
			conn.Read(one[:])
		}
	case 3:
	default:
		return true
	}
	n.drop()
	return false
}

func (n *colonyNeighbor) sendOffer(conn net.Conn) bool {
	n.mu.Lock()
	id := strconv.Itoa(n.next)
	offer, ok := n.offers[id]
	if !ok {
		hash := econ.EventHash(n.config.seed, int64(n.next), 31)
		offer = neighborOffer{id: id, good: "moss", amount: 1 + int(hash%3), price: 3 + float64(hash%5)}
		n.offers[id] = offer
	}
	n.mu.Unlock()
	return n.write(conn, fmt.Sprintf("OFFER %s %s %d %.2f\n", offer.id, offer.good, offer.amount, offer.price))
}

func (n *colonyNeighbor) ship(id string, accepted bool) string {
	n.tradeMu.Lock()
	defer n.tradeMu.Unlock()
	n.mu.Lock()
	paid, already := n.paid[id]
	offer, exists := n.offers[id]
	reconnects := max(0, n.connections-1)
	n.mu.Unlock()
	if already {
		return fmt.Sprintf("PAID %s %.2f\n", id, paid)
	}
	if !exists {
		return "ERR " + id + " unknown-offer\n"
	}
	if !accepted {
		return "ERR " + id + " not-accepted\n"
	}
	if reason := n.requestPayment(offer, reconnects); reason != "" {
		return "ERR " + id + " " + reason + "\n"
	}
	credit := offer.price * float64(offer.amount)
	n.mu.Lock()
	n.paid[id] = credit
	delete(n.offers, id)
	n.status.Trades++
	if id == strconv.Itoa(n.next) {
		n.next++
	}
	n.mu.Unlock()
	n.publish()
	return fmt.Sprintf("PAID %s %.2f\n", id, credit)
}

func (n *colonyNeighbor) requestPayment(offer neighborOffer, reconnects int) string {
	receipt := n.config.prefix + ":" + offer.id
	ack := make(chan neighborReceipt, 1)
	n.mu.Lock()
	n.pending[receipt] = ack
	n.mu.Unlock()
	defer func() { n.mu.Lock(); delete(n.pending, receipt); n.mu.Unlock() }()
	command, _ := json.Marshal(map[string]any{"credit": offer.price * float64(offer.amount), "take": map[string]int{offer.good: offer.amount}, "receipt": receipt, "neighbor": econ.NeighborProgress{Reconnects: reconnects, MaxLine: econ.NeighborLineLimit}})
	if err := n.config.commands.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		return "stopped"
	}
	if _, err := n.config.commands.Write(append(command, '\n')); err != nil {
		return "stopped"
	}
	return n.awaitPayment(ack)
}

func (n *colonyNeighbor) awaitPayment(ack <-chan neighborReceipt) string {
	select {
	case <-n.closed:
		return "stopped"
	case response := <-ack:
		if !response.OK {
			return "insufficient-goods"
		}
		return ""
	case <-time.After(2 * time.Second):
		return "timeout"
	}
}

func (c *colonyGame) prepareNeighbor(cmd *exec.Cmd, commands *os.File) error {
	c.mu.Lock()
	owned, seed, high := c.state.Has("sys.neighbor"), c.state.Seed, c.state.Active == "n31-bad-line"
	c.neighbor = nil
	c.mu.Unlock()
	if !owned {
		return nil
	}
	n, err := newColonyNeighbor(neighborConfig{seed: seed, high: high, prefix: "neighbor:" + filepath.Base(cmd.Dir), commands: commands, publish: c.broadcast})
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.neighbor = n
	c.mu.Unlock()
	cmd.Env = append(cmd.Env, "ORBITAL_NEIGHBOR_ADDR="+n.snapshot().Addr)
	n.publish()
	return nil
}

func (c *colonyGame) startNeighborProcess(cmd *exec.Cmd, commands *os.File) error {
	if err := c.prepareNeighbor(cmd, commands); err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		c.stopNeighbor()
		return err
	}
	return nil
}

func (c *colonyGame) stopNeighbor() {
	c.mu.Lock()
	n := c.neighbor
	c.mu.Unlock()
	if n != nil {
		n.close()
	}
}

func (c *colonyGame) observeNeighbor(line string) {
	var ack neighborReceipt
	if json.Unmarshal([]byte(line), &ack) != nil || ack.Op != "transfer" {
		return
	}
	c.mu.Lock()
	n := c.neighbor
	c.mu.Unlock()
	if n == nil {
		return
	}
	n.mu.Lock()
	pending := n.pending[ack.ID]
	n.mu.Unlock()
	if pending != nil {
		select {
		case pending <- ack:
		default:
		}
	}
}

func (c *colonyGame) neighborSummary(summary map[string]any) map[string]any {
	status := neighborStatus{}
	if c.neighbor != nil {
		status = c.neighbor.snapshot()
	}
	systems, ok := summary["systems"].(map[string]any)
	if !ok {
		systems = map[string]any{}
		summary["systems"] = systems
	}
	systems["neighbor"] = status
	return summary
}

func (n *colonyNeighbor) decline(id string) bool {
	n.tradeMu.Lock()
	defer n.tradeMu.Unlock()
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, exists := n.offers[id]; !exists {
		return false
	}
	delete(n.offers, id)
	if id == strconv.Itoa(n.next) {
		n.next++
	}
	return true
}

func (n *colonyNeighbor) drop() {
	n.mu.Lock()
	n.status.Drops++
	n.mu.Unlock()
	n.publish()
}
