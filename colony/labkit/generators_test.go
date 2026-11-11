package labkit

import (
	"bytes"
	"encoding/binary"
	"encoding/csv"
	"errors"
	"net"
	"strings"
	"syscall"
	"testing"
	"unicode/utf8"
)

func TestSeededGeneratorsAreRepeatableAndStructured(t *testing.T) {
	for name, generate := range map[string]func(int64, int) string{"logs": Logs, "csv": CSV, "text": Text} {
		t.Run(name, func(t *testing.T) {
			got := generate(42, 20)
			if got != generate(42, 20) || got == generate(43, 20) {
				t.Fatal("seed not respected")
			}
			if !utf8.ValidString(got) {
				t.Fatal("invalid UTF-8")
			}
		})
	}
	rows, err := csv.NewReader(strings.NewReader(CSV(42, 20))).ReadAll()
	if err != nil || len(rows) != 21 {
		t.Fatalf("rows=%d err=%v", len(rows), err)
	}
	if len(strings.Split(strings.TrimSpace(Logs(42, 20)), "\n")) != 20 {
		t.Fatal("wrong log count")
	}
	text := Text(42, 100)
	if !strings.Contains(text, "قمر") || !strings.Contains(text, "🚀") {
		t.Fatal(text)
	}
	frames := Frames(42, 20)
	if !bytes.Equal(frames, Frames(42, 20)) || bytes.Equal(frames, Frames(43, 20)) {
		t.Fatal("frame seed not respected")
	}
	count := 0
	for len(frames) > 0 {
		if len(frames) < 2 {
			t.Fatal("truncated header")
		}
		n := int(binary.BigEndian.Uint16(frames))
		frames = frames[2:]
		if n < 1 || n > len(frames) {
			t.Fatal("bad length")
		}
		frames = frames[n:]
		count++
	}
	if count != 20 {
		t.Fatal(count)
	}
	if Logs(1, 0) != "" || Text(1, 0) != "" || len(Frames(1, 0)) != 0 {
		t.Fatal("nonempty zero dataset")
	}
}
func TestFreeAddrBindsLoopback(t *testing.T) {
	listener, err := net.Listen("tcp", FreeAddr())
	if errors.Is(err, syscall.EPERM) {
		t.Skip("sandbox prohibits loopback listeners")
	}
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	addr := listener.Addr().(*net.TCPAddr)
	if !addr.IP.IsLoopback() || addr.Port == 0 {
		t.Fatal(addr)
	}
}
