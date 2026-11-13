package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/labkit"
)

func pipeNeighbor(t *testing.T) *colonyNeighbor {
	t.Helper()
	return &colonyNeighbor{config: neighborConfig{high: true, prefix: "neighbor:test"}, status: neighborStatus{On: true}, conns: map[net.Conn]struct{}{}, next: 1, offers: map[string]neighborOffer{}, paid: map[string]float64{}, pending: map[string]chan neighborReceipt{}, closed: make(chan struct{})}
}

func TestNeighborProtocolFaultsRequireOnlyBoundedLines(t *testing.T) {
	labkit.LeakCheck(t)
	for _, tc := range []struct {
		name   string
		number int
	}{{"oversized", 1}, {"garbage", 2}, {"stall", 3}, {"drop", 4}} {
		t.Run(tc.name, func(t *testing.T) {
			n := pipeNeighbor(t)
			server, client := net.Pipe()
			defer client.Close()
			finished := make(chan bool, 1)
			go func() { defer server.Close(); finished <- n.disturb(server, tc.number) }()
			client.SetReadDeadline(time.Now().Add(60 * time.Millisecond))
			reader := bufio.NewReaderSize(client, econ.NeighborLineLimit)
			line, err := reader.ReadSlice('\n')
			switch tc.name {
			case "oversized":
				if !errors.Is(err, bufio.ErrBufferFull) || len(line) != econ.NeighborLineLimit {
					t.Fatalf("unbounded read: %d %v", len(line), err)
				}
			case "garbage":
				if err != nil || string(line) != "GARBAGE ? ?\n" {
					t.Fatalf("garbage: %q %v", line, err)
				}
			case "stall":
				if timeout, ok := err.(net.Error); !ok || !timeout.Timeout() {
					t.Fatalf("stall didn't require deadline: %v", err)
				}
			case "drop":
				if !errors.Is(err, io.EOF) {
					t.Fatalf("drop: %v", err)
				}
			}
			client.Close()
			select {
			case <-finished:
			case <-time.After(time.Second):
				t.Fatal("fault goroutine leaked")
			}
		})
	}
}

func TestNeighborPaymentExactlyOnceAcrossConcurrentReconnects(t *testing.T) {
	labkit.LeakCheck(t)
	n := pipeNeighbor(t)
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	n.config.commands = writer
	var commands sync.WaitGroup
	commands.Add(1)
	payments := 0
	go func() {
		defer commands.Done()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			var request struct {
				Receipt string
				Credit  float64
				Take    map[string]int
			}
			if json.Unmarshal(scanner.Bytes(), &request) != nil {
				continue
			}
			if request.Take["moss"] < 1 || request.Credit <= 0 {
				t.Error("invalid transfer command")
			}
			payments++
			n.mu.Lock()
			ack := n.pending[request.Receipt]
			n.mu.Unlock()
			ack <- neighborReceipt{ID: request.Receipt, OK: true}
		}
	}()
	defer func() { writer.Close(); commands.Wait(); reader.Close() }()
	open := func() (net.Conn, *bufio.Reader, string) {
		server, client := net.Pipe()
		n.wg.Add(1)
		go n.serve(server, 5)
		client.SetDeadline(time.Now().Add(time.Second))
		br := bufio.NewReaderSize(client, econ.NeighborLineLimit)
		offer, err := br.ReadString('\n')
		if err != nil {
			t.Fatal(err)
		}
		fields := strings.Fields(offer)
		if len(fields) != 5 || fields[0] != "OFFER" {
			t.Fatal(offer)
		}
		return client, br, fields[1]
	}
	client, br, id := open()
	fmt.Fprintf(client, "SHIP %s\n", id)
	if line, _ := br.ReadString('\n'); line != "ERR "+id+" not-accepted\n" {
		t.Fatal(line)
	}
	fmt.Fprintf(client, "ACCEPT %s\nSHIP %s\n", id, id)
	paid, err := br.ReadString('\n')
	client.Close()
	if err != nil || !strings.HasPrefix(paid, "PAID "+id+" ") {
		t.Fatalf("payment: %q %v", paid, err)
	}
	var clients sync.WaitGroup
	for range 12 {
		clients.Add(1)
		go func() {
			defer clients.Done()
			server, client := net.Pipe()
			n.wg.Add(1)
			go n.serve(server, 5)
			defer client.Close()
			client.SetDeadline(time.Now().Add(time.Second))
			reader := bufio.NewReaderSize(client, econ.NeighborLineLimit)
			reader.ReadString('\n')
			fmt.Fprintf(client, "SHIP %s\n", id)
			line, err := reader.ReadString('\n')
			if err != nil || line != paid {
				t.Errorf("replayed payment: %q %v", line, err)
			}
		}()
	}
	clients.Wait()
	n.wg.Wait()
	writer.Close()
	commands.Wait()
	if payments != 1 || n.snapshot().Trades != 1 {
		t.Fatalf("duplicate payments: commands=%d trades=%d", payments, n.snapshot().Trades)
	}
}

func TestNeighborListenerStopsWithRunAndDoesNotLeak(t *testing.T) {
	labkit.LeakCheck(t)
	for range 3 {
		n, err := newColonyNeighbor(neighborConfig{seed: 4})
		if errors.Is(err, syscall.EPERM) {
			t.Skip("sandbox prohibits loopback listeners")
		}
		if err != nil {
			t.Fatal(err)
		}
		addr := n.snapshot().Addr
		if !strings.HasPrefix(addr, "127.0.0.1:") {
			t.Fatal(addr)
		}
		conn, err := net.DialTimeout("tcp", addr, time.Second)
		if err != nil {
			n.close()
			t.Fatal(err)
		}
		c := &colonyGame{neighbor: n}
		c.stopNeighbor()
		conn.Close()
		if n.snapshot().On {
			t.Fatal("sidecar still on")
		}
		if reopened, err := net.DialTimeout("tcp", addr, 30*time.Millisecond); err == nil {
			reopened.Close()
			t.Fatal("listener still open")
		}
	}
}

func TestNeighborDeclineAndMalformedCommands(t *testing.T) {
	labkit.LeakCheck(t)
	n := pipeNeighbor(t)
	server, client := net.Pipe()
	n.wg.Add(1)
	go n.serve(server, 5)
	defer func() { client.Close(); n.wg.Wait() }()
	client.SetDeadline(time.Now().Add(time.Second))
	reader := bufio.NewReaderSize(client, econ.NeighborLineLimit)
	offer, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	id := strings.Fields(offer)[1]
	for _, tc := range []struct{ command, want string }{
		{"garbage\n", "ERR - bad-command\n"},
		{"ACCEPT unknown\n", "ERR unknown unknown-offer\n"},
		{"WRONG " + id + "\n", "ERR " + id + " bad-command\n"},
	} {
		io.WriteString(client, tc.command)
		if reply, err := reader.ReadString('\n'); err != nil || reply != tc.want {
			t.Fatalf("%q: %q %v", tc.command, reply, err)
		}
	}
	fmt.Fprintf(client, "ACCEPT %s\nDECLINE %s\n", id, id)
	next, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(next, "OFFER 2 ") {
		t.Fatalf("next offer: %q %v", next, err)
	}
	fmt.Fprintf(client, "SHIP %s\n", id)
	if reply, err := reader.ReadString('\n'); err != nil || reply != "ERR "+id+" unknown-offer\n" {
		t.Fatalf("declined offer paid: %q %v", reply, err)
	}
}
