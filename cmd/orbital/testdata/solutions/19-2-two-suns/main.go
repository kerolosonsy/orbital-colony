package main

import (
	"fmt"
	"time"
	_ "time/tzdata"
)

type SunReadings struct {
	Home, Remote string
	DayChanged   bool
}

func ReadSuns(instant, home, remote string) (SunReadings, error) {
	t, err := time.Parse(time.RFC3339, instant)
	if err != nil {
		return SunReadings{}, fmt.Errorf("instant: %w", err)
	}
	a, err := time.LoadLocation(home)
	if err != nil {
		return SunReadings{}, fmt.Errorf("zone %s: %w", home, err)
	}
	b, err := time.LoadLocation(remote)
	if err != nil {
		return SunReadings{}, fmt.Errorf("zone %s: %w", remote, err)
	}
	x, y := t.In(a), t.In(b)
	return SunReadings{x.Format("2006-01-02 15:04 -0700"), y.Format("2006-01-02 15:04 -0700"), x.Format("2006-01-02") != y.Format("2006-01-02")}, nil
}
func main() {}
