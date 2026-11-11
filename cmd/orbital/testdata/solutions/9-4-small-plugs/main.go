package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type BeaconLookup interface{ Beacon(string) (string, bool) }

func RelayOrders(in io.Reader, out io.Writer, beacons BeaconLookup) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		tag := strings.TrimSpace(scanner.Text())
		if tag == "end" {
			_, err := fmt.Fprintln(out, "closed")
			return err
		}
		address, ok := beacons.Beacon(tag)
		if !ok {
			address = "unknown"
		}
		if _, err := fmt.Fprintln(out, address); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	_, err := fmt.Fprintln(out, "exhausted")
	return err
}
func main() {}
