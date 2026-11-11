package main

import (
	"flag"
	"fmt"
	"io"
)

type StationHost struct {
	OS, Arch, Go string
	Lookup       func(string) (string, bool)
}

func StationInfo(args []string, host StationHost, out io.Writer) error {
	fs := flag.NewFlagSet("station", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	system := fs.Bool("system", true, "")
	crew := fs.Bool("crew", false, "")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 0 {
		return fmt.Errorf("unexpected arguments")
	}
	if *system {
		for _, line := range []struct{ key, text string }{{"os", host.OS}, {"arch", host.Arch}, {"go", host.Go}} {
			if _, err := fmt.Fprintf(out, "%s=%q\n", line.key, line.text); err != nil {
				return err
			}
		}
	}
	if *crew {
		return writeStationCrew(host.Lookup, out)
	}
	return nil
}
func writeStationCrew(lookup func(string) (string, bool), out io.Writer) error {
	crew, present := lookup("COLONY_CREW")
	if !present {
		_, err := fmt.Fprintln(out, "crew=<missing>")
		return err
	}
	_, err := fmt.Fprintf(out, "crew=%q\n", crew)
	return err
}
func main() {}
