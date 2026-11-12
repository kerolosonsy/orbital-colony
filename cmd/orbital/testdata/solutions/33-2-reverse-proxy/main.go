package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type relay struct {
	routes  map[string]*url.URL
	timeout time.Duration
	proxy   *httputil.ReverseProxy
}

type routeKey struct{}

var errUpstream = errors.New("upstream failed")

func NewRelay(routes map[string]string, timeout time.Duration, transport http.RoundTripper) (http.Handler, error) {
	if len(routes) == 0 || timeout <= 0 || timeout > time.Minute || transport == nil {
		return nil, errors.New("invalid relay")
	}
	rl := &relay{routes: map[string]*url.URL{}, timeout: timeout}
	for prefix, target := range routes {
		u, err := url.Parse(target)
		if err != nil || !strings.HasPrefix(prefix, "/") || !strings.HasSuffix(prefix, "/") || len(prefix) < 3 ||
			(u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || (u.Path != "" && u.Path != "/") || u.RawQuery != "" {
			return nil, errors.New("invalid route " + prefix)
		}
		u.Path = ""
		rl.routes[prefix] = u
	}
	rl.proxy = &httputil.ReverseProxy{
		Transport: transport,
		Rewrite: func(pr *httputil.ProxyRequest) {
			prefix := pr.In.Context().Value(routeKey{}).(string)
			pr.SetURL(rl.routes[prefix])
			pr.Out.URL.Path = "/" + strings.TrimPrefix(pr.In.URL.Path, prefix)
			pr.Out.URL.RawPath = ""
			pr.Out.Header.Del("X-Colony-Role")
			pr.Out.Header.Set("X-Relay-Route", prefix)
			pr.SetXForwarded()
		},
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Del("Server")
			if resp.StatusCode >= 500 {
				return errUpstream
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			switch {
			case errors.Is(err, context.DeadlineExceeded):
				http.Error(w, "upstream timeout", http.StatusGatewayTimeout)
			case errors.Is(err, errUpstream):
				http.Error(w, "upstream failed", http.StatusBadGateway)
			default:
				http.Error(w, "upstream unavailable", http.StatusBadGateway)
			}
		},
	}
	return rl, nil
}

func (rl *relay) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	best := ""
	for prefix := range rl.routes {
		if strings.HasPrefix(r.URL.Path, prefix) && len(prefix) > len(best) {
			best = prefix
		}
	}
	if best == "" {
		http.Error(w, "no route", http.StatusNotFound)
		return
	}
	ctx, cancel := context.WithTimeout(context.WithValue(r.Context(), routeKey{}, best), rl.timeout)
	defer cancel()
	rl.proxy.ServeHTTP(w, r.WithContext(ctx))
}

func main() {}
