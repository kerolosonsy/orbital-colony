package main

import (
	"context"
	"io"
	"net/http"
	"slices"
	"strings"
)

type Route struct {
	Method, Pattern string
	Handler         http.Handler
}
type Middleware func(http.Handler) http.Handler
type dockKey int

const (
	paramsKey dockKey = iota
	requestIDKey
)

func Param(r *http.Request, name string) string {
	params, _ := r.Context().Value(paramsKey).(map[string]string)
	return params[name]
}
func RequestID(r *http.Request) string { id, _ := r.Context().Value(requestIDKey).(string); return id }
func WithRequestID(nextID func() string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := nextID()
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), requestIDKey, id)))
		})
	}
}
func matchDock(pattern, path string) (map[string]string, bool) {
	expected, actual := strings.Split(pattern, "/"), strings.Split(path, "/")
	if len(expected) != len(actual) {
		return nil, false
	}
	params := map[string]string{}
	for i, part := range expected {
		if strings.HasPrefix(part, ":") {
			if actual[i] == "" {
				return nil, false
			}
			params[part[1:]] = actual[i]
		} else if part != actual[i] {
			return nil, false
		}
	}
	return params, true
}
func routeDock(routes []Route, w http.ResponseWriter, r *http.Request) {
	methods := []string{}
	for _, route := range routes {
		params, ok := matchDock(route.Pattern, r.URL.Path)
		if !ok {
			continue
		}
		methods = append(methods, route.Method)
		if route.Method == r.Method {
			route.Handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), paramsKey, params)))
			return
		}
	}
	status, text := 404, "no route\n"
	if len(methods) > 0 {
		slices.Sort(methods)
		methods = slices.Compact(methods)
		w.Header().Set("Allow", strings.Join(methods, ", "))
		status, text = 405, "wrong method\n"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	if r.Method != "HEAD" {
		_, _ = io.WriteString(w, text)
	}
}
func NewDockRouter(routes []Route, chain ...Middleware) http.Handler {
	snapshot := slices.Clone(routes)
	var h http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { routeDock(snapshot, w, r) })
	for i := len(chain) - 1; i >= 0; i-- {
		h = chain[i](h)
	}
	return h
}
func main() {}
