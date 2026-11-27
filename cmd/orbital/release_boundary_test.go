package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestListenConfigOnlyAllowsLocalOrigins(t *testing.T) {
	for _, tc := range []struct {
		addr, origin string
		allowed      bool
	}{
		{"127.0.0.1:7777", "", true},
		{"localhost:7777", "http://127.0.0.1:7777", true},
		{"[::1]:7777", "http://[::1]:7777", true},
		{"0.0.0.0:8080", "", false},
		{"0.0.0.0:8080", "https://orbital.example", false},
		{"[::]:8080", "http://localhost:8080", false},
		{"127.0.0.1:7777", "https://orbital.example", false},
		{"127.0.0.1:7777", "http://localhost:7777/unexpected", false},
		{"127.0.0.1:7777", "http://localhost:7777/?unexpected", false},
	} {
		err := validateListenConfig(tc.addr, tc.origin)
		if (err == nil) != tc.allowed {
			t.Errorf("addr %q origin %q: got %v, allowed %t", tc.addr, tc.origin, err, tc.allowed)
		}
	}
}

func TestStaticAssetsExcludeTestsAndDirectories(t *testing.T) {
	handler := staticHandler(fstest.MapFS{
		"app.js":                     {Data: []byte("console.log(1)")},
		"app.test.js":                {Data: []byte("test")},
		"sandbox/harness.js":         {Data: []byte("harness")},
		"sandbox/ch01.fixtures.json": {Data: []byte("{}")},
	})
	for _, path := range []string{"/", "/sandbox/", "/app.test.js", "/sandbox/harness.js", "/sandbox/ch01.fixtures.json"} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("%s returned %d", path, rec.Code)
		}
	}
}

func TestLevelsVaryAndRespectCompressionPreference(t *testing.T) {
	s := &server{}
	req := httptest.NewRequest(http.MethodGet, "/api/levels", nil)
	req.Header.Set("Accept-Encoding", "br, gzip;q=0, *;q=0.5")
	rec := httptest.NewRecorder()
	s.levels(rec, req)
	if rec.Code != http.StatusOK || rec.Header().Get("Content-Encoding") != "" {
		t.Fatalf("expected uncompressed response, got %d %v", rec.Code, rec.Header())
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("If-None-Match", rec.Header().Get("ETag"))
	rec = httptest.NewRecorder()
	s.levels(rec, req)
	if rec.Code != http.StatusNotModified || !strings.Contains(rec.Header().Get("Vary"), "X-Orbital-Lang") ||
		!strings.Contains(rec.Header().Get("Vary"), "Accept-Encoding") {
		t.Fatalf("expected cache validators and variance, got %d %v", rec.Code, rec.Header())
	}
}
