package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var ErrAddress = errors.New("address")
var ErrRadio = errors.New("radio")
var ErrLaunchDate = errors.New("launch date")

type CrewCard struct{ Address, Radio, LaunchDate string }

var stationDomain = regexp.MustCompile(`^[A-Za-z0-9]+(\.[A-Za-z0-9]+)+$`)
var radioCode = regexp.MustCompile(`^\+[0-9]{7,12}$`)

func CheckAddress(s string) error {
	parts := strings.Split(s, "@")
	if len(parts) != 2 || parts[0] == "" || !stationDomain.MatchString(parts[1]) {
		return fmt.Errorf("crew address: %w", ErrAddress)
	}
	for _, r := range s {
		if unicode.IsSpace(r) {
			return fmt.Errorf("crew address: %w", ErrAddress)
		}
	}
	return nil
}
func CheckRadio(s string) error {
	if !radioCode.MatchString(s) {
		return fmt.Errorf("radio code: %w", ErrRadio)
	}
	return nil
}
func CheckLaunchDate(s string) error {
	t, err := time.Parse("2006/01/02", s)
	if err != nil || t.Format("2006/01/02") != s {
		return fmt.Errorf("launch date: %w", ErrLaunchDate)
	}
	return nil
}
func ValidateCrewCard(card CrewCard) error {
	return errors.Join(CheckAddress(card.Address), CheckRadio(card.Radio), CheckLaunchDate(card.LaunchDate))
}
func InvalidFields(err error) []string {
	var fields []string
	for _, entry := range []struct {
		cause error
		name  string
	}{{ErrAddress, "address"}, {ErrRadio, "radio"}, {ErrLaunchDate, "launchDate"}} {
		if errors.Is(err, entry.cause) {
			fields = append(fields, entry.name)
		}
	}
	return fields
}
func main() {}
