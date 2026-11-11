package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type FlightSummary struct {
	Info, Warn, Error int
	BadLines          []int
}

func flightLevel(line string) string {
	fields := strings.SplitN(line, " ", 3)
	if len(fields) != 3 {
		return ""
	}
	seq, err := strconv.Atoi(fields[0])
	if err != nil || seq < 0 || strings.TrimSpace(fields[2]) == "" {
		return ""
	}
	return fields[1]
}
func ReadFlight(in io.Reader) (FlightSummary, error) {
	var summary FlightSummary
	var scanner *bufio.Scanner = bufio.NewScanner(in)
	scanner.Buffer(make([]byte, 4096), (256<<10)+2)
	for n := 1; scanner.Scan(); n++ {
		line := scanner.Text()
		if len(line) > 256<<10 {
			return FlightSummary{}, fmt.Errorf("flight line too long")
		}
		switch flightLevel(line) {
		case "INFO":
			summary.Info++
		case "WARN":
			summary.Warn++
		case "ERROR":
			summary.Error++
		default:
			summary.BadLines = append(summary.BadLines, n)
		}
	}
	if err := scanner.Err(); err != nil {
		return FlightSummary{}, err
	}
	return summary, nil
}
func main() {}
