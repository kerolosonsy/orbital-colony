package main

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type Berth struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Slots int    `json:"slots"`
}
type berthRegistry struct {
	mu     sync.Mutex
	next   int
	berths map[int]Berth
}

func berthText(w http.ResponseWriter, r *http.Request, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	if r.Method != "HEAD" {
		_, _ = io.WriteString(w, text)
	}
}
func berthJSON(w http.ResponseWriter, status int, berth any) {
	encoded, _ := json.Marshal(berth)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(encoded)
}
func berthPath(path string) (int, int) {
	if path == "/berths" {
		return 0, 0
	}
	if !strings.HasPrefix(path, "/berths/") {
		return 0, 404
	}
	segment := strings.TrimPrefix(path, "/berths/")
	if segment == "" || strings.Contains(segment, "/") {
		return 0, 404
	}
	id, e := strconv.Atoi(segment)
	if e != nil || id <= 0 || strconv.Itoa(id) != segment {
		return 0, 400
	}
	return id, 0
}
func berthInput(w http.ResponseWriter, r *http.Request) (*Berth, bool) {
	raw, e := io.ReadAll(io.LimitReader(r.Body, 4097))
	if len(raw) > 4096 {
		berthText(w, r, 413, "too large\n")
		return nil, false
	}
	var berth *Berth
	if e != nil || json.Unmarshal(raw, &berth) != nil || berth == nil || strings.TrimSpace(berth.Label) == "" || berth.Slots < 1 || berth.Slots > 64 {
		berthText(w, r, 400, "invalid berth\n")
		return nil, false
	}
	berth.Label = strings.TrimSpace(berth.Label)
	return berth, true
}
func (s *berthRegistry) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id, status := berthPath(r.URL.Path)
	if status != 0 {
		text := "unknown berth\n"
		if status == 400 {
			text = "invalid berth\n"
		}
		berthText(w, r, status, text)
		return
	}
	allow := "GET, POST"
	if id != 0 {
		allow = "DELETE, GET, PUT"
	}
	permitted := id == 0 && (r.Method == "GET" || r.Method == "POST") || id != 0 && (r.Method == "GET" || r.Method == "PUT" || r.Method == "DELETE")
	if !permitted {
		w.Header().Set("Allow", allow)
		berthText(w, r, 405, "use route methods\n")
		return
	}
	var incoming *Berth
	if r.Method == "POST" || r.Method == "PUT" {
		var ok bool
		incoming, ok = berthInput(w, r)
		if !ok {
			return
		}
	}
	s.reply(w, r, id, incoming)
}
func (s *berthRegistry) reply(w http.ResponseWriter, r *http.Request, id int, incoming *Berth) {
	s.mu.Lock()
	if id == 0 {
		if r.Method == "GET" {
			snapshot := make([]Berth, 0, len(s.berths))
			for _, b := range s.berths {
				snapshot = append(snapshot, b)
			}
			s.mu.Unlock()
			slices.SortFunc(snapshot, func(a, b Berth) int { return a.ID - b.ID })
			berthJSON(w, 200, snapshot)
			return
		}
		s.next++
		incoming.ID = s.next
		s.berths[s.next] = *incoming
		s.mu.Unlock()
		w.Header().Set("Location", "/berths/"+strconv.Itoa(incoming.ID))
		berthJSON(w, 201, *incoming)
		return
	}
	berth, ok := s.berths[id]
	if !ok {
		s.mu.Unlock()
		berthText(w, r, 404, "unknown berth\n")
		return
	}
	switch r.Method {
	case "PUT":
		incoming.ID = id
		berth = *incoming
		s.berths[id] = berth
	case "DELETE":
		delete(s.berths, id)
	}
	s.mu.Unlock()
	if r.Method == "DELETE" {
		w.WriteHeader(204)
		return
	}
	berthJSON(w, 200, berth)
}
func NewBerthAPI() http.Handler { return &berthRegistry{berths: make(map[int]Berth)} }
func main()                     {}
