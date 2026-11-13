package main

import (
	"encoding/csv"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"orbital/colony"
)

func number(s string) (float64, bool) {
	n, err := strconv.ParseFloat(s, 64)
	return n, err == nil && n >= 0 && !math.IsNaN(n) && !math.IsInf(n, 0)
}

func main() {
	f, err := os.Open(filepath.Join(os.Getenv("ORBITAL_COLONY_DIR"), "ledger.csv"))
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	_, _ = r.Read()
	totals := map[string]float64{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || len(row) != 5 {
			continue
		}
		_, tickErr := strconv.ParseInt(row[0], 10, 64)
		_, amountOK := number(row[2])
		_, priceOK := number(row[3])
		credits, creditsOK := number(row[4])
		known := false
		for _, good := range []string{"moss", "lichen", "glowcap", "reed", "crystal", "iron", "geode", "water", "silica", "copper", "sulfur", "titanium", "helium", "uranium", "biomass", "plate", "wire", "glass", "fuel", "fertilizer", "circuit", "alloy", "battery", "coolant", "processor", "frame", "cell", "core"} {
			known = known || row[1] == good
		}
		if tickErr == nil && amountOK && priceOK && creditsOK && known {
			totals[row[1]] += credits
		}
	}
	colony.Submit("audit", totals)
}
