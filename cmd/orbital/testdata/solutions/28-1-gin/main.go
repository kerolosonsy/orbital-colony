package main

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Berth struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Slots int    `json:"slots"`
}

type berthInput struct {
	Label string `json:"label" binding:"required"`
	Slots int    `json:"slots" binding:"required,min=1,max=64"`
}

type berthFilter struct {
	MinSlots *int `form:"min_slots" binding:"omitempty,min=1,max=64"`
}

type berthStore struct {
	mu     sync.Mutex
	next   int
	berths map[int]Berth
}

func noStore(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	c.Next()
}

func readBerth(c *gin.Context) (Berth, bool) {
	raw, err := io.ReadAll(io.LimitReader(c.Request.Body, 4097))
	if err == nil && len(raw) > 4096 {
		c.String(413, "too large\n")
		return Berth{}, false
	}
	var in berthInput
	if err != nil || !json.Valid(raw) || binding.JSON.BindBody(raw, &in) != nil || strings.TrimSpace(in.Label) == "" {
		c.String(http.StatusBadRequest, "invalid berth\n")
		return Berth{}, false
	}
	return Berth{Label: strings.TrimSpace(in.Label), Slots: in.Slots}, true
}

func berthID(c *gin.Context) (int, bool) {
	raw := c.Param("id")
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 || strconv.Itoa(id) != raw {
		c.String(http.StatusBadRequest, "invalid berth\n")
		return 0, false
	}
	return id, true
}

func (s *berthStore) list(c *gin.Context) {
	var f berthFilter
	if c.ShouldBindQuery(&f) != nil {
		c.String(http.StatusBadRequest, "invalid berth\n")
		return
	}
	s.mu.Lock()
	out := make([]Berth, 0, len(s.berths))
	for _, b := range s.berths {
		if f.MinSlots == nil || b.Slots >= *f.MinSlots {
			out = append(out, b)
		}
	}
	s.mu.Unlock()
	slices.SortFunc(out, func(a, b Berth) int { return a.ID - b.ID })
	c.JSON(http.StatusOK, out)
}

func (s *berthStore) create(c *gin.Context) {
	b, ok := readBerth(c)
	if !ok {
		return
	}
	s.mu.Lock()
	s.next++
	b.ID = s.next
	s.berths[b.ID] = b
	s.mu.Unlock()
	c.Header("Location", "/berths/"+strconv.Itoa(b.ID))
	c.JSON(http.StatusCreated, b)
}

func (s *berthStore) read(c *gin.Context) {
	id, ok := berthID(c)
	if !ok {
		return
	}
	s.mu.Lock()
	b, found := s.berths[id]
	s.mu.Unlock()
	if !found {
		c.String(http.StatusNotFound, "unknown berth\n")
		return
	}
	c.JSON(http.StatusOK, b)
}

func (s *berthStore) replace(c *gin.Context) {
	id, ok := berthID(c)
	if !ok {
		return
	}
	b, ok := readBerth(c)
	if !ok {
		return
	}
	b.ID = id
	s.mu.Lock()
	_, found := s.berths[id]
	if found {
		s.berths[id] = b
	}
	s.mu.Unlock()
	if !found {
		c.String(http.StatusNotFound, "unknown berth\n")
		return
	}
	c.JSON(http.StatusOK, b)
}

func (s *berthStore) remove(c *gin.Context) {
	id, ok := berthID(c)
	if !ok {
		return
	}
	s.mu.Lock()
	_, found := s.berths[id]
	delete(s.berths, id)
	s.mu.Unlock()
	if !found {
		c.String(http.StatusNotFound, "unknown berth\n")
		return
	}
	c.Status(http.StatusNoContent)
}

func NewBerthEngine() *gin.Engine {
	s := &berthStore{berths: map[int]Berth{}}
	r := gin.New()
	r.RedirectTrailingSlash = false
	r.HandleMethodNotAllowed = true
	r.Use(noStore)
	r.NoRoute(func(c *gin.Context) { c.String(http.StatusNotFound, "no route\n") })
	r.NoMethod(func(c *gin.Context) { c.String(http.StatusMethodNotAllowed, "use route methods\n") })
	r.GET("/berths", s.list)
	r.POST("/berths", s.create)
	r.GET("/berths/:id", s.read)
	r.PUT("/berths/:id", s.replace)
	r.DELETE("/berths/:id", s.remove)
	return r
}

func main() {}
