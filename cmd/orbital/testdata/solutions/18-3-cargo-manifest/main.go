package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type CargoTotals struct {
	Rows, Units int
	ByCallsign  map[string]int
}

func ReadManifest(in io.Reader) (CargoTotals, error) {
	reader := csv.NewReader(in)
	header, err := reader.Read()
	if err != nil {
		return CargoTotals{}, fmt.Errorf("header: %w", err)
	}
	if len(header) != 3 || header[0] != "id" || header[1] != "callsign" || header[2] != "cargo" {
		return CargoTotals{}, fmt.Errorf("header: expected id,callsign,cargo")
	}
	totals := CargoTotals{ByCallsign: make(map[string]int)}
	for row := 2; ; row++ {
		fields, err := reader.Read()
		if err == io.EOF {
			return totals, nil
		}
		if err != nil {
			return CargoTotals{}, fmt.Errorf("record %d: %w", row, err)
		}
		id, idErr := strconv.Atoi(fields[0])
		units, unitsErr := strconv.Atoi(fields[2])
		name := strings.TrimSpace(fields[1])
		if idErr != nil || id < 0 || unitsErr != nil || units < 0 || units > 1000000 || name == "" {
			return CargoTotals{}, fmt.Errorf("record %d: invalid cargo fields", row)
		}
		totals.Rows++
		totals.Units += units
		totals.ByCallsign[name] += units
	}
}
func main() {}
