package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type RepairOrder struct {
	ID, Units int
	Sector    string
	Done      bool
}
type RepairBook interface {
	Queue(string, int) (int, error)
	Orders() []RepairOrder
	Complete(int) error
	Remove(int) error
}
type repairMemory struct {
	orders []RepairOrder
	issued int
}

func invalidSector(sector string) bool {
	for _, r := range sector {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
func (m *repairMemory) Queue(sector string, units int) (int, error) {
	if sector == "" || invalidSector(sector) || units < 1 || units > 100 {
		return 0, fmt.Errorf("invalid repair")
	}
	m.issued++
	m.orders = append(m.orders, RepairOrder{ID: m.issued, Sector: sector, Units: units})
	return m.issued, nil
}
func (m *repairMemory) Orders() []RepairOrder { return append([]RepairOrder(nil), m.orders...) }
func (m *repairMemory) Complete(id int) error {
	for i := range m.orders {
		if m.orders[i].ID == id {
			m.orders[i].Done = true
			return nil
		}
	}
	return fmt.Errorf("missing repair")
}
func (m *repairMemory) Remove(id int) error {
	for i := range m.orders {
		if m.orders[i].ID == id {
			m.orders = append(m.orders[:i], m.orders[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("missing repair")
}
func repairReply(args []string, book RepairBook) (string, error) {
	if len(args) == 1 && args[0] == "list" {
		var lines strings.Builder
		for _, order := range book.Orders() {
			fmt.Fprintf(&lines, "%d|%q|%d|%t\n", order.ID, order.Sector, order.Units, order.Done)
		}
		return lines.String(), nil
	}
	if len(args) == 3 && args[0] == "add" {
		units, err := strconv.Atoi(args[2])
		if err != nil {
			return "", err
		}
		id, err := book.Queue(args[1], units)
		return fmt.Sprintf("queued=%d\n", id), err
	}
	return repairIDReply(args, book)
}
func repairIDReply(args []string, book RepairBook) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("invalid command")
	}
	id, err := strconv.Atoi(args[1])
	if err != nil {
		return "", err
	}
	switch args[0] {
	case "done":
		return fmt.Sprintf("done=%d\n", id), book.Complete(id)
	case "rm":
		return fmt.Sprintf("removed=%d\n", id), book.Remove(id)
	}
	return "", fmt.Errorf("unknown command")
}
func dispatchRepair(args []string, book RepairBook, out io.Writer) int {
	reply, err := repairReply(args, book)
	code := 0
	if err != nil {
		reply = "error\n"
		code = 2
	}
	if _, err := fmt.Fprint(out, reply); err != nil {
		return 1
	}
	return code
}
func run(args []string, stdin io.Reader, stdout io.Writer) int {
	book := &repairMemory{}
	if len(args) > 0 {
		return dispatchRepair(args, book, stdout)
	}
	scanner := bufio.NewScanner(stdin)
	scanner.Buffer(make([]byte, 1024), 4098)
	code := 0
	for scanner.Scan() {
		if len(scanner.Text()) > 4096 {
			return 1
		}
		next := dispatchRepair(strings.Fields(scanner.Text()), book, stdout)
		if next == 1 {
			return 1
		}
		if next == 2 {
			code = 2
		}
	}
	if scanner.Err() != nil {
		return 1
	}
	return code
}
func main() {}
