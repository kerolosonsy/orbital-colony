package main

import (
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Berth struct {
	ID    int    `json:"id"`
	Label string `json:"label"`
	Slots int    `json:"slots"`
}

type berthStore struct {
	mu     sync.Mutex
	next   int
	berths map[int]Berth
}

func noStore(c *fiber.Ctx) error {
	c.Set("Cache-Control", "no-store")
	return c.Next()
}

func berthFault(c *fiber.Ctx, err error) error {
	var fe *fiber.Error
	if errors.As(err, &fe) {
		switch fe.Code {
		case fiber.StatusNotFound:
			return c.Status(fiber.StatusNotFound).SendString("no route\n")
		case fiber.StatusMethodNotAllowed:
			return c.Status(fiber.StatusMethodNotAllowed).SendString("use route methods\n")
		}
	}
	return c.Status(fiber.StatusInternalServerError).SendString("registry fault\n")
}

func plain(c *fiber.Ctx, status int, text string) error {
	return c.Status(status).SendString(text)
}

// readBerth returns a validated berth, or the status to reply with.
func readBerth(c *fiber.Ctx) (Berth, int) {
	raw := c.Body()
	if len(raw) > 4096 {
		return Berth{}, fiber.StatusRequestEntityTooLarge
	}
	var in *Berth
	if json.Unmarshal(raw, &in) != nil || in == nil || strings.TrimSpace(in.Label) == "" || in.Slots < 1 || in.Slots > 64 {
		return Berth{}, fiber.StatusBadRequest
	}
	return Berth{Label: strings.TrimSpace(in.Label), Slots: in.Slots}, 0
}

func rejectBody(c *fiber.Ctx, status int) error {
	if status == fiber.StatusRequestEntityTooLarge {
		return plain(c, status, "too large\n")
	}
	return plain(c, status, "invalid berth\n")
}

func berthID(c *fiber.Ctx) (int, bool) {
	raw := c.Params("id")
	id, err := strconv.Atoi(raw)
	return id, err == nil && id > 0 && strconv.Itoa(id) == raw
}

func (s *berthStore) list(c *fiber.Ctx) error {
	s.mu.Lock()
	out := make([]Berth, 0, len(s.berths))
	for _, b := range s.berths {
		out = append(out, b)
	}
	s.mu.Unlock()
	slices.SortFunc(out, func(a, b Berth) int { return a.ID - b.ID })
	return c.JSON(out)
}

func (s *berthStore) create(c *fiber.Ctx) error {
	b, status := readBerth(c)
	if status != 0 {
		return rejectBody(c, status)
	}
	s.mu.Lock()
	s.next++
	b.ID = s.next
	s.berths[b.ID] = b
	s.mu.Unlock()
	c.Set("Location", "/berths/"+strconv.Itoa(b.ID))
	return c.Status(fiber.StatusCreated).JSON(b)
}

func (s *berthStore) read(c *fiber.Ctx) error {
	id, ok := berthID(c)
	if !ok {
		return plain(c, fiber.StatusBadRequest, "invalid berth\n")
	}
	s.mu.Lock()
	b, found := s.berths[id]
	s.mu.Unlock()
	if !found {
		return plain(c, fiber.StatusNotFound, "unknown berth\n")
	}
	return c.JSON(b)
}

func (s *berthStore) replace(c *fiber.Ctx) error {
	id, ok := berthID(c)
	if !ok {
		return plain(c, fiber.StatusBadRequest, "invalid berth\n")
	}
	b, status := readBerth(c)
	if status != 0 {
		return rejectBody(c, status)
	}
	b.ID = id
	s.mu.Lock()
	_, found := s.berths[id]
	if found {
		s.berths[id] = b
	}
	s.mu.Unlock()
	if !found {
		return plain(c, fiber.StatusNotFound, "unknown berth\n")
	}
	return c.JSON(b)
}

func (s *berthStore) remove(c *fiber.Ctx) error {
	id, ok := berthID(c)
	if !ok {
		return plain(c, fiber.StatusBadRequest, "invalid berth\n")
	}
	s.mu.Lock()
	_, found := s.berths[id]
	delete(s.berths, id)
	s.mu.Unlock()
	if !found {
		return plain(c, fiber.StatusNotFound, "unknown berth\n")
	}
	c.Status(fiber.StatusNoContent)
	return nil
}

func NewBerthApp() *fiber.App {
	s := &berthStore{berths: map[int]Berth{}}
	app := fiber.New(fiber.Config{StrictRouting: true, CaseSensitive: true, ErrorHandler: berthFault})
	app.Use(noStore)
	app.Get("/berths", s.list)
	app.Post("/berths", s.create)
	app.Get("/berths/:id", s.read)
	app.Put("/berths/:id", s.replace)
	app.Delete("/berths/:id", s.remove)
	return app
}

func main() {}
