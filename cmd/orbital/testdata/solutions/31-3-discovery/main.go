package main

import (
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

var ErrNeighborPacket = errors.New("invalid neighbor packet")

func neighborToken(s, alphabet string) bool {
	if len(s) < 1 || len(s) > 16 {
		return false
	}
	for _, b := range []byte(s) {
		if !strings.ContainsRune(alphabet, rune(b)) {
			return false
		}
	}
	return true
}
func neighborCode(s string) bool    { return neighborToken(s, "0123456789abcdef") }
func neighborStation(s string) bool { return neighborToken(s, "abcdefghijklmnopqrstuvwxyz0123456789-") }
func OpenNeighbor(listen func(string, string) (net.PacketConn, error)) (net.PacketConn, error) {
	if listen == nil {
		listen = net.ListenPacket
	}
	return listen("udp4", "127.0.0.1:0")
}
func EncodeNeighborProbe(nonce string) ([]byte, error) {
	if !neighborCode(nonce) {
		return nil, ErrNeighborPacket
	}
	return []byte("Q1|" + nonce), nil
}
func DecodeNeighborReply(packet []byte, nonce string) (string, error) {
	parts := strings.Split(string(packet), "|")
	if len(packet) > 64 || len(parts) != 3 || parts[0] != "A1" || !neighborCode(nonce) || parts[1] != nonce || !neighborStation(parts[2]) {
		return "", ErrNeighborPacket
	}
	return parts[2], nil
}
func AnswerNeighbor(pc net.PacketConn, station string, deadline time.Time) error {
	if !neighborStation(station) {
		return ErrNeighborPacket
	}
	if err := pc.SetDeadline(deadline); err != nil {
		return err
	}
	buf := make([]byte, 65)
	n, peer, err := pc.ReadFrom(buf)
	if err != nil {
		return err
	}
	parts := strings.Split(string(buf[:n]), "|")
	if n > 64 || len(parts) != 2 || parts[0] != "Q1" || !neighborCode(parts[1]) {
		return ErrNeighborPacket
	}
	reply := []byte("A1|" + parts[1] + "|" + station)
	n, err = pc.WriteTo(reply, peer)
	if err != nil {
		return err
	}
	if n != len(reply) {
		return io.ErrShortWrite
	}
	return nil
}
func main() {}
