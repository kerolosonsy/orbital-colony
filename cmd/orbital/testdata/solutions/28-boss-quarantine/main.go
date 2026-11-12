package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type QuarantinePod struct {
	ID        string
	Occupants int
}

func NewQuarantine(token string) *http.Server {
	var mu sync.Mutex
	pods := map[string]QuarantinePod{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", r.Header.Get("X-Request-ID"))
		fail := func(status int, code string) {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(map[string]string{"error": code})
		}
		if r.Header.Get("Authorization") != "Bearer "+token {
			fail(401, "unauthorized")
			return
		}
		if r.URL.Path != "/pods" && (!strings.HasPrefix(r.URL.Path, "/pods/") || strings.TrimPrefix(r.URL.Path, "/pods/") == "" || strings.Contains(strings.TrimPrefix(r.URL.Path, "/pods/"), "/")) {
			fail(404, "not_found")
			return
		}
		collection := r.URL.Path == "/pods"
		allow := "GET"
		if collection {
			allow = "GET, POST"
		}
		if r.Method != "GET" && !(collection && r.Method == "POST") {
			w.Header().Set("Allow", allow)
			fail(405, "method")
			return
		}
		if r.Method == "POST" && r.Header.Get("X-Role") != "medic" {
			fail(403, "forbidden")
			return
		}
		if r.Method == "POST" {
			data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024))
			if err != nil {
				fail(413, "too_large")
				return
			}
			var pod QuarantinePod
			if err = json.Unmarshal(data, &pod); err != nil || pod.ID == "" || strings.Contains(pod.ID, "/") || pod.Occupants < 0 {
				fail(400, "invalid")
				return
			}
			mu.Lock()
			defer mu.Unlock()
			if _, ok := pods[pod.ID]; ok {
				fail(409, "conflict")
				return
			}
			pods[pod.ID] = pod
			w.Header().Set("Location", "/pods/"+pod.ID)
			w.WriteHeader(201)
			json.NewEncoder(w).Encode(pod)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if !collection {
			pod, ok := pods[strings.TrimPrefix(r.URL.Path, "/pods/")]
			if !ok {
				fail(404, "not_found")
				return
			}
			json.NewEncoder(w).Encode(pod)
			return
		}
		list := make([]QuarantinePod, 0, len(pods))
		for _, pod := range pods {
			list = append(list, pod)
		}
		sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
		json.NewEncoder(w).Encode(list)
	})
	return &http.Server{Addr: "127.0.0.1:0", Handler: handler, ReadHeaderTimeout: 2 * time.Second, ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second, IdleTimeout: 30 * time.Second}
}
func main() {}
