package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"orbital/colony"
)

func main() {
	for colony.Tick() < 3600 {
		colony.Wait()
	}
	day := colony.Tick() / 1200
	report := map[string]map[string]float64{}
	for d := day - 3; d < day; d++ {
		report[strconv.Itoa(d)] = map[string]float64{}
	}
	for _, suffix := range []string{".2", ".1", ""} {
		f, err := os.Open(filepath.Join(os.Getenv("ORBITAL_COLONY_DIR"), "events.jsonl"+suffix))
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return
		}
		decoder := json.NewDecoder(f)
		for {
			var event struct {
				Op    string
				Day   int
				Goods map[string]float64
			}
			if err := decoder.Decode(&event); err != nil {
				if err != io.EOF {
					f.Close()
					return
				}
				break
			}
			if goods := report[strconv.Itoa(event.Day)]; event.Op == "harvest" && goods != nil {
				for good, amount := range event.Goods {
					goods[good] += amount
				}
			}
		}
		f.Close()
	}
	colony.Submit("shift-report", report)
}
