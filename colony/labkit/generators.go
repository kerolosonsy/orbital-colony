package labkit

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"math/rand"
	"strings"
)

// FreeAddr asks the listener to bind a free loopback port atomically.
func FreeAddr() string { return "127.0.0.1:0" }

// Logs returns n deterministic log lines ("0007 WARN drone=42") for the seed.
func Logs(seed int64, n int) string {
	rng := rand.New(rand.NewSource(seed))
	var out strings.Builder
	levels := []string{"INFO", "WARN", "ERROR"}
	for i := 0; i < n; i++ {
		fmt.Fprintf(&out, "%04d %s drone=%d\n", i, levels[rng.Intn(len(levels))], rng.Intn(100))
	}
	return out.String()
}

// CSV returns a deterministic CSV with header id,callsign,cargo and n rows,
// including quoted commas and non-ASCII names.
func CSV(seed int64, n int) string {
	rng := rand.New(rand.NewSource(seed))
	var out bytes.Buffer
	writer := csv.NewWriter(&out)
	writer.Write([]string{"id", "callsign", "cargo"})
	names := []string{"قمر", "drone, alpha", "مسبار 🚀"}
	for i := 0; i < n; i++ {
		writer.Write([]string{fmt.Sprint(i), names[rng.Intn(len(names))], fmt.Sprint(rng.Intn(1000))})
	}
	writer.Flush()
	return out.String()
}

// Text returns n deterministic space-separated words, mixing Arabic, English and emoji.
func Text(seed int64, n int) string {
	rng := rand.New(rand.NewSource(seed))
	words := []string{"قمر", "مستعمرة", "🚀", "🛰️", "طاقم", "orbit"}
	var out []string
	for i := 0; i < n; i++ {
		out = append(out, words[rng.Intn(len(words))])
	}
	return strings.Join(out, " ")
}

// Frames encodes n payloads as a big-endian uint16 length followed by bytes.
func Frames(seed int64, n int) []byte {
	rng := rand.New(rand.NewSource(seed))
	var out []byte
	for i := 0; i < n; i++ {
		payload := make([]byte, 1+rng.Intn(32))
		rng.Read(payload)
		out = binary.BigEndian.AppendUint16(out, uint16(len(payload)))
		out = append(out, payload...)
	}
	return out
}
