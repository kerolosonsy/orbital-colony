package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

var ErrBeaconLine = errors.New("beacon line too long")

func writeBeacon(c net.Conn, line []byte) error {
	for len(line) > 0 {
		n, err := c.Write(line)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrShortWrite
		}
		line = line[n:]
	}
	return nil
}
func readBeacon(r *bufio.Reader, limit int) ([]byte, error) {
	line := make([]byte, 0)
	for {
		b, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) && len(line) > 0 {
				return nil, io.ErrUnexpectedEOF
			}
			return nil, err
		}
		if b == '\n' {
			return append(line, b), nil
		}
		if len(line) == limit {
			return nil, ErrBeaconLine
		}
		line = append(line, b)
	}
}
func EchoBeacon(c net.Conn, limit int, deadline time.Time) error {
	defer c.Close()
	if limit < 0 || limit > 4096 {
		return errors.New("invalid limit")
	}
	if err := c.SetDeadline(deadline); err != nil {
		return err
	}
	r := bufio.NewReader(c)
	for {
		line, err := readBeacon(r, limit)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err = writeBeacon(c, line); err != nil {
			return err
		}
	}
}
func main() {}
