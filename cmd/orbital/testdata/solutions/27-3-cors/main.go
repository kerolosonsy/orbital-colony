package main

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

type Policy struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool
	MaxAge      int
}

func validOrigin(origin string) bool {
	u, err := url.Parse(origin)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil &&
		u.Scheme+"://"+u.Host == origin && strings.ToLower(origin) == origin
}

func validMethod(m string) bool {
	return m != "" && strings.Trim(m, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") == ""
}

func Shield(p Policy) (func(http.Handler) http.Handler, error) {
	if p.MaxAge < 0 || p.MaxAge > 86400 {
		return nil, errors.New("max age out of range")
	}
	for _, o := range p.Origins {
		if !validOrigin(o) {
			return nil, errors.New("invalid origin " + o)
		}
	}
	for _, m := range p.Methods {
		if !validMethod(m) {
			return nil, errors.New("invalid method " + m)
		}
	}
	for _, h := range p.Headers {
		if h == "" || strings.ContainsAny(h, " ,") {
			return nil, errors.New("invalid header " + h)
		}
	}
	origins := slices.Clone(p.Origins)
	methods := slices.Clone(p.Methods)
	headers := slices.Clone(p.Headers)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			h.Set("Referrer-Policy", "no-referrer")
			h.Set("Content-Security-Policy", "default-src 'none'")
			h.Add("Vary", "Origin")
			origin := r.Header.Get("Origin")
			wanted := r.Header.Get("Access-Control-Request-Method")
			preflight := r.Method == http.MethodOptions && wanted != ""
			if preflight {
				h.Add("Vary", "Access-Control-Request-Method")
				h.Add("Vary", "Access-Control-Request-Headers")
			}
			allowed := origin != "" && slices.Contains(origins, origin)
			if !preflight {
				if allowed {
					h.Set("Access-Control-Allow-Origin", origin)
					if p.Credentials {
						h.Set("Access-Control-Allow-Credentials", "true")
					}
				}
				next.ServeHTTP(w, r)
				return
			}
			ok := allowed && slices.Contains(methods, wanted)
			for _, name := range strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",") {
				name = strings.TrimSpace(name)
				if name != "" && !slices.ContainsFunc(headers, func(a string) bool { return strings.EqualFold(a, name) }) {
					ok = false
				}
			}
			if !ok {
				h.Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte("preflight denied\n"))
				return
			}
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			if len(headers) > 0 {
				h.Set("Access-Control-Allow-Headers", strings.Join(headers, ", "))
			}
			if p.MaxAge > 0 {
				h.Set("Access-Control-Max-Age", strconv.Itoa(p.MaxAge))
			}
			if p.Credentials {
				h.Set("Access-Control-Allow-Credentials", "true")
			}
			w.WriteHeader(http.StatusNoContent)
		})
	}, nil
}

func main() {}
