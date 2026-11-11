package main

import "orbital/colony"

func main() {
	for i := 0; i < colony.ReadingCount(); i++ {
		n := colony.Reading(i)
		band := "invalid"
		if n < 0 || n > 100 {
			band = "invalid"
		} else if n < 20 {
			band = "low"
		} else if n < 50 {
			band = "ready"
		} else if n < 80 {
			band = "high"
		} else {
			band = "alert"
		}
		colony.Report(band)
	}
}
