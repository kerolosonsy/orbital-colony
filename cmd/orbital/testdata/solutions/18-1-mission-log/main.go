package main

import (
	"errors"
	"strings"
)

type MissionEntry struct {
	ID     int
	Sector string
	Sealed bool
}
type MissionLog struct {
	entries []MissionEntry
	issued  int
}

func (m *MissionLog) Log(sector string) (int, error) {
	sector = strings.TrimSpace(sector)
	if sector == "" {
		return 0, errors.New("sector required")
	}
	m.issued++
	m.entries = append(m.entries, MissionEntry{ID: m.issued, Sector: sector})
	return m.issued, nil
}
func (m *MissionLog) Seal(id int) error {
	for i := range m.entries {
		if m.entries[i].ID == id {
			m.entries[i].Sealed = true
			return nil
		}
	}
	return errors.New("mission ID not found")
}
func (m *MissionLog) Erase(id int) error {
	for i := range m.entries {
		if m.entries[i].ID == id {
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return nil
		}
	}
	return errors.New("mission ID not found")
}
func (m *MissionLog) Entries() []MissionEntry {
	return append([]MissionEntry(nil), m.entries...)
}
func main() {}
