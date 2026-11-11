package main

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"orbital/colony/labkit"
)

type Limits struct {
	Burst   int
	Every   time.Duration
	Idle    time.Duration
	Proxies []string
}

type bucket struct {
	credit int
	mark   time.Time // refill marker
	seen   time.Time // last request
}

type Limiter struct {
	mu      sync.Mutex
	clock   labkit.Clock
	burst   int
	every   time.Duration
	idle    time.Duration
	proxies map[string]bool
	buckets map[string]*bucket
}

func NewLimiter(clock labkit.Clock, l Limits) (*Limiter, error) {
	if clock == nil || l.Burst < 1 || l.Burst > 1000 || l.Every <= 0 || l.Every > time.Hour || l.Idle < time.Duration(l.Burst)*l.Every {
		return nil, errors.New("invalid limits")
	}
	lim := &Limiter{clock: clock, burst: l.Burst, every: l.Every, idle: l.Idle, proxies: map[string]bool{}, buckets: map[string]*bucket{}}
	for _, p := range l.Proxies {
		ip := net.ParseIP(p)
		if ip == nil {
			return nil, errors.New("invalid proxy " + p)
		}
		lim.proxies[ip.String()] = true
	}
	return lim, nil
}

// client returns the rate-limit key for r, or "" when it cannot be trusted.
func (l *Limiter) client(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	peer := net.ParseIP(host)
	if err != nil || peer == nil {
		return ""
	}
	values := r.Header.Values("X-Forwarded-For")
	if !l.proxies[peer.String()] || len(values) == 0 {
		return peer.String()
	}
	hops := strings.Split(strings.Join(values, ","), ",")
	ip := net.ParseIP(strings.TrimSpace(hops[len(hops)-1]))
	if ip == nil {
		return ""
	}
	return ip.String()
}

// take spends one permit for key; on refusal it returns the wait until the next one.
func (l *Limiter) take(key string) (int, time.Duration, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock.Now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{credit: l.burst, mark: now}
		l.buckets[key] = b
	}
	if now.After(b.mark) {
		steps := now.Sub(b.mark) / l.every
		if steps >= time.Duration(l.burst-b.credit) {
			b.credit, b.mark = l.burst, now
		} else {
			b.credit += int(steps)
			b.mark = b.mark.Add(steps * l.every)
		}
	}
	if now.After(b.seen) {
		b.seen = now
	}
	if b.credit == 0 {
		return 0, b.mark.Add(l.every).Sub(now), false
	}
	b.credit--
	return b.credit, 0, true
}

func (l *Limiter) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := l.client(r)
		if key == "" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("bad client\n"))
			return
		}
		left, wait, ok := l.take(key)
		if !ok {
			secs := (wait + time.Second - 1) / time.Second
			w.Header().Set("Retry-After", strconv.Itoa(int(secs)))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte("slow down\n"))
			return
		}
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(left))
		next.ServeHTTP(w, r)
	})
}

func (l *Limiter) Sweep() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.clock.Now()
	removed := 0
	for key, b := range l.buckets {
		if now.Sub(b.seen) >= l.idle {
			delete(l.buckets, key)
			removed++
		}
	}
	return removed
}

func (l *Limiter) Clients() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.buckets)
}

func main() {}
