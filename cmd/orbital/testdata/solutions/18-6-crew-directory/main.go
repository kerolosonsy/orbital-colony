package main

import (
	"fmt"
	"slices"
	"strings"
)

type Crew struct{ ID, Callsign, Role string }
type CrewStore interface{ SaveCrew([]Crew) error }
type CrewDirectory struct{ crew []Crew }

func (d *CrewDirectory) Register(crew Crew) error {
	if strings.TrimSpace(crew.ID) == "" || strings.TrimSpace(crew.Callsign) == "" || strings.TrimSpace(crew.Role) == "" {
		return fmt.Errorf("crew fields required")
	}
	for _, known := range d.crew {
		if known.ID == crew.ID {
			return fmt.Errorf("duplicate crew ID")
		}
	}
	d.crew = append(d.crew, crew)
	return nil
}
func (d *CrewDirectory) Search(query string) []Crew {
	query = strings.ToLower(query)
	var matches []Crew
	for _, crew := range d.crew {
		if strings.Contains(strings.ToLower(crew.Callsign), query) || strings.Contains(strings.ToLower(crew.Role), query) {
			matches = append(matches, crew)
		}
	}
	slices.SortFunc(matches, func(a, b Crew) int { return strings.Compare(a.ID, b.ID) })
	return matches
}
func (d *CrewDirectory) Save(store CrewStore) error {
	if store == nil {
		return fmt.Errorf("crew store required")
	}
	return store.SaveCrew(d.Search(""))
}
func main() {}
