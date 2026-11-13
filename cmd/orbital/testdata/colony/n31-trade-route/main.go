package main

import (
	"bufio"
	"fmt"
	"net"
	"orbital/colony"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	go func() {
		for colony.Tick() < 1200 {
			time.Sleep(time.Millisecond)
		}
		for {
			conn, err := net.DialTimeout("tcp", os.Getenv("ORBITAL_NEIGHBOR_ADDR"), 100*time.Millisecond)
			if err != nil {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			trade(conn)
			conn.Close()
		}
	}()
	w, h := colony.WorldSize()
	for {
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				colony.Till()
				colony.TryHarvest()
				colony.Plant(colony.Moss)
				if x < w-1 {
					colony.Move(colony.East)
				}
			}
			for x := 1; x < w; x++ {
				colony.Move(colony.West)
			}
			if y < h-1 {
				colony.Move(colony.North)
			}
		}
		for y := 1; y < h; y++ {
			colony.Move(colony.South)
		}
	}
}

func trade(conn net.Conn) {
	reader := bufio.NewReaderSize(conn, 1024)
	for {
		conn.SetDeadline(time.Now().Add(150 * time.Millisecond))
		line, err := reader.ReadSlice('\n')
		if err != nil || len(line) > 1024 {
			return
		}
		fields := strings.Fields(string(line))
		if len(fields) != 5 || fields[0] != "OFFER" {
			continue
		}
		amount, err := strconv.Atoi(fields[3])
		if err != nil || fields[2] != "moss" || amount < 1 || amount > 10 {
			return
		}
		id := fields[1]
		if _, err := fmt.Fprintf(conn, "ACCEPT %s\n", id); err != nil {
			return
		}
		for colony.Count(colony.Moss) < amount {
			time.Sleep(time.Millisecond)
		}
		if _, err := fmt.Fprintf(conn, "SHIP %s\n", id); err != nil {
			return
		}
	}
}
