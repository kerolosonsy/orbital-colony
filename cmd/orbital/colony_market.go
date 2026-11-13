package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"orbital/colony/econ"
)

type marketOrder struct {
	Good   string  `json:"good"`
	Amount float64 `json:"amount"`
}

type marketResponse struct {
	Key     string  `json:"key,omitempty"`
	Good    string  `json:"good,omitempty"`
	Amount  float64 `json:"amount,omitempty"`
	Credits float64 `json:"credits,omitempty"`
	Error   string  `json:"error,omitempty"`
}

type marketEntry struct {
	order    marketOrder
	ready    chan struct{}
	status   int
	slow     bool
	response marketResponse
}

type colonyMarket struct {
	mu       sync.Mutex
	entries  map[string]*marketEntry
	attempts int64
	window   time.Time
	requests int
	seed     int64
	failure  uint64
	lang     string
	url      string
	server   *http.Server
	done     chan struct{}
	run      *archiveMarketRun
	prices   func() map[string]float64
	orders   atomic.Int64
	rejected atomic.Int64
}

func (m *colonyMarket) start() error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	m.url = "http://" + listener.Addr().String()
	m.server = &http.Server{Handler: m, ReadHeaderTimeout: time.Second, ReadTimeout: 2 * time.Second, WriteTimeout: 4 * time.Second, IdleTimeout: time.Second}
	m.done = make(chan struct{})
	go func() { defer close(m.done); _ = m.server.Serve(listener) }()
	return nil
}

func (m *colonyMarket) close() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := m.server.Shutdown(ctx); err != nil {
		_ = m.server.Close()
	}
	<-m.done
}

func (m *colonyMarket) reply(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func (m *colonyMarket) refuse(w http.ResponseWriter, status int, ar, en string) {
	m.rejected.Add(1)
	m.reply(w, status, marketResponse{Error: tr(m.lang, ar, en)})
	m.run.publish()
}

func (m *colonyMarket) admission() (int, string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	if now.Sub(m.window) >= time.Second {
		m.window, m.requests = now, 0
	}
	if m.requests >= 10 {
		return 429, tr(m.lang, "طلبات كتير؛ جرّب بعد ثانية", "Too many requests; retry after one second")
	}
	m.requests++
	m.attempts++
	if econ.EventHash(m.seed, m.attempts, 24)%100 < m.failure {
		return 503, tr(m.lang, "السوق مش متاح مؤقتاً", "Market temporarily unavailable")
	}
	return 200, ""
}

func (m *colonyMarket) admit(w http.ResponseWriter) bool {
	status, message := m.admission()
	if status == 200 {
		return true
	}
	m.transient(w, status, message)
	return false
}

func (m *colonyMarket) transient(w http.ResponseWriter, status int, message string) {
	if status == 429 {
		w.Header().Set("Retry-After", "1")
	}
	m.rejected.Add(1)
	m.reply(w, status, marketResponse{Error: message})
	m.run.publish()
}

func (m *colonyMarket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/prices":
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			m.refuse(w, 405, "استخدم GET", "Use GET")
			return
		}
		if m.admit(w) {
			m.reply(w, 200, m.prices())
		}
	case "/orders":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			m.refuse(w, 405, "استخدم POST", "Use POST")
			return
		}
		m.order(w, r)
	default:
		m.refuse(w, 404, "المسار مش موجود", "Unknown endpoint")
	}
}

func (m *colonyMarket) readOrder(w http.ResponseWriter, r *http.Request) (string, marketOrder, bool) {
	key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
	var order marketOrder
	if key == "" || len(key) > 128 {
		m.refuse(w, 400, "محتاج Idempotency-Key بطول 1–128", "Idempotency-Key must contain 1–128 bytes")
		return key, order, false
	}
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&order); err != nil {
		m.refuse(w, 400, "JSON غير صالح", "Invalid JSON")
		return key, order, false
	}
	if decoder.Decode(new(any)) != io.EOF || econ.BasePrice[order.Good] <= 0 || order.Amount <= 0 || order.Amount > 1e9 {
		m.refuse(w, 400, "سلعة أو كمية غير صالحة", "Invalid good or amount")
		return key, order, false
	}
	return key, order, true
}

func (m *colonyMarket) entry(key string, order marketOrder) (*marketEntry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if entry := m.entries[key]; entry != nil {
		return entry, false
	}
	if len(m.entries) >= 4096 {
		return nil, false
	}
	entry := &marketEntry{order: order, ready: make(chan struct{})}
	m.entries[key] = entry
	return entry, true
}

func (m *colonyMarket) replay(w http.ResponseWriter, r *http.Request, entry *marketEntry) {
	select {
	case <-entry.ready:
		if entry.status == 429 {
			w.Header().Set("Retry-After", "1")
		}
		m.reply(w, entry.status, entry.response)
	case <-r.Context().Done():
	case <-m.run.ctx.Done():
		m.refuse(w, 503, "التشغيل وقف", "Run stopped")
	}
}

func (m *colonyMarket) order(w http.ResponseWriter, r *http.Request) {
	key, order, ok := m.readOrder(w, r)
	if !ok {
		return
	}
	entry, owner := m.entry(key, order)
	if entry == nil {
		w.Header().Set("Retry-After", "1")
		m.refuse(w, 429, "وصلت لحد الطلبات في التشغيل ده", "Run order limit reached")
		return
	}
	if entry.order != order {
		m.refuse(w, 409, "المفتاح مستخدم لطلب مختلف", "Key already used for a different order")
		return
	}
	if !owner {
		m.replay(w, r, entry)
		return
	}
	if status, message := m.admission(); status != 200 {
		entry.status, entry.response.Error = status, message
		m.mu.Lock()
		delete(m.entries, key)
		m.mu.Unlock()
		close(entry.ready)
		m.transient(w, status, message)
		return
	}
	m.fill(key, entry)
	if entry.slow {
		timer := time.NewTimer(600 * time.Millisecond)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-m.run.ctx.Done():
		}
	}
	m.reply(w, entry.status, entry.response)
}

func (m *colonyMarket) fill(key string, entry *marketEntry) {
	defer close(entry.ready)
	receipt, err := m.run.transact(key, entry.order)
	entry.response = marketResponse{Key: key, Good: entry.order.Good, Amount: entry.order.Amount, Credits: receipt.Credits}
	entry.status = 200
	if err != nil {
		entry.status, entry.response.Error = 503, tr(m.lang, "التشغيل وقف", "Run stopped")
	} else if !receipt.OK {
		entry.status, entry.response.Error = 409, receipt.Error
	}
	if entry.status == 200 {
		n := m.orders.Add(1)
		entry.slow = m.failure == 30 && n == int64(5+econ.EventHash(m.seed, 0, 33)%5)
	} else {
		m.rejected.Add(1)
	}
	m.run.publish()
}

func (m *colonyMarket) summary() map[string]any {
	return map[string]any{"on": true, "url": m.url, "orders": m.orders.Load(), "rejected": m.rejected.Load()}
}

func marketCommandID(n uint64) string { return "market-" + strconv.FormatUint(n, 10) }

func marketStartError(lang string, err error) error {
	return fmt.Errorf("%s: %w", tr(lang, "تعذّر تشغيل السوق", "Cannot start market"), err)
}
