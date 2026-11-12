package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type unit struct {
	addr    string
	target  *url.URL
	healthy bool
}

type Pool struct {
	mu        sync.Mutex
	units     []*unit
	next      int
	transport http.RoundTripper
	proxy     *httputil.ReverseProxy
}

type pickKey struct{}

// attempt travels in the request context so the proxy callbacks know which unit is in use.
type attempt struct {
	u      *unit
	failed bool
}

var errSick = errors.New("unit answered 5xx")

func NewPool(backends []string, transport http.RoundTripper) (*Pool, error) {
	if len(backends) == 0 || len(backends) > 64 || transport == nil {
		return nil, errors.New("invalid pool")
	}
	p := &Pool{transport: transport}
	seen := map[string]bool{}
	for _, b := range backends {
		u, err := url.Parse(b)
		if err != nil || seen[b] || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" {
			return nil, errors.New("invalid backend " + b)
		}
		seen[b] = true
		u.Path = ""
		p.units = append(p.units, &unit{addr: b, target: u, healthy: true})
	}
	p.proxy = &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(pr.In.Context().Value(pickKey{}).(*attempt).u.target)
			pr.SetXForwarded()
		},
		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode >= 500 {
				return errSick
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			r.Context().Value(pickKey{}).(*attempt).failed = true
		},
	}
	return p, nil
}

// pick returns the next healthy unit in round-robin order that is not in tried.
func (p *Pool) pick(tried map[*unit]bool) *unit {
	p.mu.Lock()
	defer p.mu.Unlock()
	n := len(p.units)
	for i := 0; i < n; i++ {
		u := p.units[(p.next+i)%n]
		if u.healthy && !tried[u] {
			p.next = (p.next + i + 1) % n
			return u
		}
	}
	return nil
}

func (p *Pool) mark(u *unit, healthy bool) {
	p.mu.Lock()
	u.healthy = healthy
	p.mu.Unlock()
}

func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tried := map[*unit]bool{}
	for {
		u := p.pick(tried)
		if u == nil {
			http.Error(w, "no healthy backend", http.StatusServiceUnavailable)
			return
		}
		tried[u] = true
		a := &attempt{u: u}
		p.proxy.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), pickKey{}, a)))
		if !a.failed {
			return
		}
		p.mark(u, false)
		if r.Method != "GET" && r.Method != "HEAD" {
			http.Error(w, "backend failed", http.StatusBadGateway)
			return
		}
	}
}

// Check probes every unit's /healthz without holding the lock and returns how many are healthy.
func (p *Pool) Check(ctx context.Context) int {
	p.mu.Lock()
	units := append([]*unit(nil), p.units...)
	p.mu.Unlock()
	healthy := 0
	for _, u := range units {
		ok := false
		req, err := http.NewRequestWithContext(ctx, "GET", u.target.String()+"/healthz", nil)
		if err == nil {
			if resp, err := p.transport.RoundTrip(req); err == nil {
				ok = resp.StatusCode == 200
				resp.Body.Close()
			}
		}
		p.mark(u, ok)
		if ok {
			healthy++
		}
	}
	return healthy
}

func (p *Pool) Healthy() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	var out []string
	for _, u := range p.units {
		if u.healthy {
			out = append(out, u.addr)
		}
	}
	return out
}

func main() {}
