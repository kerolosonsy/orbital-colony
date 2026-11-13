package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"orbital/colony/econ"
	"orbital/colony/labkit"
)

// The pipe is the market's process boundary; sim tests exercise the real transaction.
func marketFixture(t *testing.T) (*colonyMarket, *atomic.Int64) {
	t.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	game := &colonyGame{state: econ.New(9)}
	s := &archiveMarketRun{ctx: ctx, cancel: cancel, pending: map[string]chan marketReceipt{}, commands: writer, game: game}
	m := &colonyMarket{seed: 11, entries: map[string]*marketEntry{}, run: s, prices: func() map[string]float64 { return map[string]float64{"moss": 1} }, lang: "en"}
	s.market = m
	game.active = &colonyRun{am: s}
	var calls atomic.Int64
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			var command struct {
				ID   string
				Take map[string]float64
			}
			if json.Unmarshal(scanner.Bytes(), &command) != nil {
				continue
			}
			calls.Add(1)
			receipt := marketReceipt{ID: command.ID, OK: command.Take["moss"] <= 100, Credits: command.Take["moss"] * 1.2}
			if !receipt.OK {
				receipt.Error = "Insufficient stock"
				receipt.Credits = 0
			}
			s.mu.Lock()
			if reply := s.pending[command.ID]; reply != nil {
				reply <- receipt
			}
			s.mu.Unlock()
		}
	}()
	t.Cleanup(func() { cancel(); writer.Close(); reader.Close(); <-done })
	return m, &calls
}

func marketRequest(m *colonyMarket, method, path, key, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if key != "" {
		req.Header.Set("Idempotency-Key", key)
	}
	response := httptest.NewRecorder()
	m.ServeHTTP(response, req)
	return response
}

func TestMarketEndpointStatusPathsAndErrorReplay(t *testing.T) {
	m, calls := marketFixture(t)
	for _, tc := range []struct {
		method, path, key, body string
		status                  int
	}{
		{"GET", "/prices", "", "", 200},
		{"POST", "/prices", "", "", 405},
		{"GET", "/orders", "", "", 405},
		{"POST", "/orders", "", `{"good":"moss","amount":1}`, 400},
		{"POST", "/orders", "bad-json", `{`, 400},
		{"POST", "/orders", "bad-good", `{"good":"unknown","amount":1}`, 400},
		{"POST", "/orders", "zero", `{"good":"moss","amount":0}`, 400},
		{"POST", "/orders", "negative", `{"good":"moss","amount":-1}`, 400},
		{"POST", "/orders", "overflow", `{"good":"moss","amount":1e999}`, 400},
		{"POST", "/orders", "extra", `{"good":"moss","amount":1,"x":2}`, 400},
		{"POST", "/orders", "trailing", `{"good":"moss","amount":1} {}`, 400},
		{"POST", "/orders", "stock", `{"good":"moss","amount":101}`, 409},
		{"POST", "/orders", "stock", `{"good":"moss","amount":101}`, 409},
		{"POST", "/orders", "stock", `{"good":"moss","amount":1}`, 409},
		{"GET", "/unknown", "", "", 404},
	} {
		response := marketRequest(m, tc.method, tc.path, tc.key, tc.body)
		if response.Code != tc.status {
			t.Fatalf("%s %s %s: %d %s", tc.method, tc.path, tc.key, response.Code, response.Body.String())
		}
		if response.Header().Get("Content-Type") != "application/json" || !json.Valid(response.Body.Bytes()) {
			t.Fatal("non-JSON response")
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("replayed or invalid order reached sim: %d", calls.Load())
	}
}

func TestMarketConcurrentSameKeyReturnsOneIdenticalSale(t *testing.T) {
	m, calls := marketFixture(t)
	var wg sync.WaitGroup
	responses := make(chan *httptest.ResponseRecorder, 50)
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			responses <- marketRequest(m, "POST", "/orders", "same", `{"good":"moss","amount":20}`)
		}()
	}
	wg.Wait()
	close(responses)
	body := ""
	for response := range responses {
		if response.Code != 200 {
			t.Fatal(response.Code, response.Body.String())
		}
		if body == "" {
			body = response.Body.String()
		}
		if body != response.Body.String() {
			t.Fatal("idempotency response changed")
		}
	}
	if calls.Load() != 1 || m.orders.Load() != 1 {
		t.Fatalf("duplicate transactions: %d %d", calls.Load(), m.orders.Load())
	}
	conflict := marketRequest(m, "POST", "/orders", "same", `{"good":"moss","amount":19}`)
	if conflict.Code != 409 || calls.Load() != 1 {
		t.Fatal("key conflict was sold")
	}
}

func TestMarketRateLimitAndSeededFailures(t *testing.T) {
	m, _ := marketFixture(t)
	for i := 0; i < 10; i++ {
		if got := marketRequest(m, "GET", "/prices", "", ""); got.Code != 200 {
			t.Fatal(got.Code)
		}
	}
	limited := marketRequest(m, "GET", "/prices", "", "")
	if limited.Code != 429 || limited.Header().Get("Retry-After") != "1" {
		t.Fatal(limited.Code, limited.Header())
	}
	m.window = time.Now().Add(-2 * time.Second)
	if got := marketRequest(m, "GET", "/prices", "", ""); got.Code != 200 {
		t.Fatal("rate window never recovered")
	}
	pattern := func(seed int64) []int {
		market, _ := marketFixture(t)
		market.seed, market.failure = seed, 30
		var statuses []int
		failures := 0
		for i := 0; i < 100; i++ {
			market.window = time.Now().Add(-2 * time.Second)
			response := marketRequest(market, "GET", "/prices", "", "")
			if response.Code == 503 {
				failures++
			} else if response.Code != 200 {
				t.Fatal(response.Code)
			}
			statuses = append(statuses, response.Code)
		}
		if failures < 15 || failures > 45 {
			t.Fatalf("30%% failure mode returned %d failures in 100 attempts", failures)
		}
		return statuses
	}
	first, second, other := fmt.Sprint(pattern(11)), fmt.Sprint(pattern(11)), fmt.Sprint(pattern(12))
	if first != second || first == other || !strings.Contains(first, "503") {
		t.Fatal("failure pattern not seeded")
	}
}

func TestMarketSlowCommittedOrderCanBeRetriedWithoutAnotherSale(t *testing.T) {
	m, calls := marketFixture(t)
	m.failure = 30
	slow := 0
	for n := 0; n < 10; n++ {
		key := fmt.Sprint(n)
		for attempt := 0; attempt < 30; attempt++ {
			m.window = time.Now().Add(-2 * time.Second)
			started := time.Now()
			response := marketRequest(m, "POST", "/orders", key, `{"good":"moss","amount":1}`)
			if response.Code == 503 {
				continue
			}
			if response.Code != 200 {
				t.Fatal(response.Code)
			}
			if time.Since(started) >= 500*time.Millisecond {
				slow++
				before := calls.Load()
				retryAt := time.Now()
				replay := marketRequest(m, "POST", "/orders", key, `{"good":"moss","amount":1}`)
				if replay.Body.String() != response.Body.String() || calls.Load() != before || time.Since(retryAt) > 200*time.Millisecond {
					t.Fatal("slow order replay was not immediate and idempotent")
				}
			}
			break
		}
	}
	if slow != 1 || calls.Load() != 10 {
		t.Fatalf("slow=%d orders=%d", slow, calls.Load())
	}
}

func TestMarketSidecarClosesListenerAndGoroutines(t *testing.T) {
	labkit.LeakCheck(t)
	requireMarketSocket(t)
	m, _ := marketFixture(t)
	if err := m.start(); err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: time.Second}
	response, err := client.Get(m.url + "/prices")
	if err != nil {
		m.run.cancel()
		m.close()
		t.Fatal(err)
	}
	response.Body.Close()
	m.run.cancel()
	m.close()
	client.CloseIdleConnections()
	conn, err := net.DialTimeout("tcp", strings.TrimPrefix(m.url, "http://"), 100*time.Millisecond)
	if err == nil {
		conn.Close()
		t.Fatal("market listener survived shutdown")
	}
}

func TestMarketCancellationReleasesPendingReceipt(t *testing.T) {
	labkit.LeakCheck(t)
	m, _ := marketFixture(t)
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	defer writer.Close()
	m.run.commands = writer
	finished := make(chan *httptest.ResponseRecorder, 1)
	go func() { finished <- marketRequest(m, "POST", "/orders", "pending", `{"good":"moss","amount":1}`) }()
	if _, err := bufio.NewReader(reader).ReadString('\n'); err != nil {
		t.Fatal(err)
	}
	m.run.cancel()
	select {
	case response := <-finished:
		if response.Code != 503 || m.orders.Load() != 0 {
			t.Fatal("pending order claimed success after shutdown")
		}
	case <-time.After(time.Second):
		t.Fatal("pending receipt leaked a handler")
	}
}
