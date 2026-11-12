package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"orbital/colony/labkit"
)

type Service struct {
	Prefix string
	Target string
	Public bool
}

type Config struct {
	Services []Service
	Keys     map[string]string
	Burst    int
	Every    time.Duration
	Clock    labkit.Clock
	NextID   func() string
}

type Middleware func(http.Handler) http.Handler

// Chain wraps h so the first middleware is the outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type route struct {
	prefix string
	target *url.URL
	public bool
}

type tank struct {
	credit int
	mark   time.Time
}

type gateway struct {
	c      Config
	routes []route
	keys   map[string]string
	mu     sync.Mutex
	tanks  map[string]*tank
	proxy  *httputil.ReverseProxy
}

type ctxKey int

const (
	routeCtx ctxKey = iota
	crewCtx
	idCtx
)

func NewGateway(c Config) (http.Handler, error) {
	if len(c.Services) == 0 || c.Clock == nil || c.NextID == nil || c.Burst < 1 || c.Burst > 1000 || c.Every <= 0 || c.Every > time.Hour {
		return nil, errors.New("invalid gateway")
	}
	g := &gateway{c: c, keys: map[string]string{}, tanks: map[string]*tank{}}
	seen := map[string]bool{}
	for _, s := range c.Services {
		u, err := url.Parse(s.Target)
		if err != nil || seen[s.Prefix] || !strings.HasPrefix(s.Prefix, "/") || !strings.HasSuffix(s.Prefix, "/") || len(s.Prefix) < 3 ||
			(u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" {
			return nil, errors.New("invalid service " + s.Prefix)
		}
		seen[s.Prefix] = true
		u.Path = ""
		g.routes = append(g.routes, route{prefix: s.Prefix, target: u, public: s.Public})
	}
	for k, crew := range c.Keys {
		if k == "" || crew == "" {
			return nil, errors.New("invalid key")
		}
		g.keys[k] = crew
	}
	g.proxy = &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			rt := pr.In.Context().Value(routeCtx).(route)
			pr.SetURL(rt.target)
			pr.Out.URL.Path = "/" + strings.TrimPrefix(pr.In.URL.Path, rt.prefix)
			pr.Out.URL.RawPath = ""
			pr.Out.Header.Del("X-Api-Key")
			pr.Out.Header.Del("X-Crew")
			if crew, _ := pr.In.Context().Value(crewCtx).(string); crew != "" {
				pr.Out.Header.Set("X-Crew", crew)
			}
			pr.Out.Header.Set("X-Request-ID", pr.In.Context().Value(idCtx).(string))
			pr.SetXForwarded()
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, "upstream unavailable", http.StatusBadGateway)
		},
	}
	return Chain(g.proxy, g.requestID, g.route, g.auth, g.limit), nil
}

func (g *gateway) requestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := g.c.NextID()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), idCtx, id)))
	})
}

func (g *gateway) route(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		best := -1
		for i, rt := range g.routes {
			if strings.HasPrefix(r.URL.Path, rt.prefix) && (best < 0 || len(rt.prefix) > len(g.routes[best].prefix)) {
				best = i
			}
		}
		if best < 0 {
			http.Error(w, "no route", http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), routeCtx, g.routes[best])))
	})
}

func (g *gateway) auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context().Value(routeCtx).(route).public {
			next.ServeHTTP(w, r)
			return
		}
		crew, ok := g.keys[r.Header.Get("X-Api-Key")]
		if !ok {
			http.Error(w, "bad key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), crewCtx, crew)))
	})
}

func (g *gateway) limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, _ := r.Context().Value(crewCtx).(string)
		if key == "" {
			host, _, _ := net.SplitHostPort(r.RemoteAddr)
			key = "ip:" + host
		} else {
			key = "crew:" + key
		}
		wait, ok := g.take(key)
		if !ok {
			w.Header().Set("Retry-After", strconv.Itoa(int((wait+time.Second-1)/time.Second)))
			http.Error(w, "slow down", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (g *gateway) take(key string) (time.Duration, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	now := g.c.Clock.Now()
	t, ok := g.tanks[key]
	if !ok {
		t = &tank{credit: g.c.Burst, mark: now}
		g.tanks[key] = t
	}
	if now.After(t.mark) {
		steps := now.Sub(t.mark) / g.c.Every
		if steps >= time.Duration(g.c.Burst-t.credit) {
			t.credit, t.mark = g.c.Burst, now
		} else {
			t.credit += int(steps)
			t.mark = t.mark.Add(steps * g.c.Every)
		}
	}
	if t.credit == 0 {
		return t.mark.Add(g.c.Every).Sub(now), false
	}
	t.credit--
	return 0, true
}

func main() {}
