package main

import (
	"orbital/colony"
	"strings"
	"time"
)

func main() {
	for _, order := range colony.Orders() {
		parts := strings.Split(order, "|")
		start, err := time.Parse("2006-01-02 15:04:05", parts[0])
		duration, durationErr := time.ParseDuration(parts[1])
		if err != nil || durationErr != nil || duration < 0 || duration > 48*time.Hour {
			colony.Report("invalid")
			continue
		}
		colony.Report(start.Add(duration).Format("2006-01-02 15:04:05") + "|" + duration.String())
	}
}
