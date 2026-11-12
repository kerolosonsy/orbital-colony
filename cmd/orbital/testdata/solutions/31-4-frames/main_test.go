package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func FuzzRadiogram(f *testing.F) {
	for _, length := range []int{0, 1, 256, 1024} {
		wire := make([]byte, 8+length)
		copy(wire, "OC")
		wire[2] = 1
		wire[3] = 2
		binary.BigEndian.PutUint16(wire[4:6], uint16(length))
		binary.BigEndian.PutUint16(wire[6:8], 0x1234)
		f.Add(wire)
	}
	for _, wire := range [][]byte{nil, {'O'}, {'O', 'C', 1, 3, 0, 0, 0, 0}, {'O', 'C', 1, 1, 0, 1, 0, 0}, {'O', 'C', 1, 1, 4, 1, 0, 0}, {'O', 'C', 1, 1, 0, 0, 0, 0, 9}} {
		f.Add(wire)
	}
	f.Fuzz(func(t *testing.T, wire []byte) {
		before := append([]byte(nil), wire...)
		got, err := ParseRadiogram(wire)
		validShape := len(wire) >= 8 && wire[0] == 'O' && wire[1] == 'C' && wire[2] == 1
		if validShape {
			n := int(wire[4])*256 + int(wire[5])
			validShape = n <= 1024 && len(wire) == 8+n
		}
		var want error
		if !validShape {
			want = ErrRadiogram
		} else if wire[3] != 1 && wire[3] != 2 {
			want = ErrRadiogramKind
		}
		if !bytes.Equal(wire, before) {
			t.Fatal("input mutation")
		}
		if want != nil {
			if !errors.Is(err, want) || got.Kind != 0 || got.Sequence != 0 || got.Payload != nil {
				t.Fatal("malformed accepted", err)
			}
			return
		}
		if err != nil || got.Kind != wire[3] || got.Sequence != uint16(wire[6])*256+uint16(wire[7]) || !bytes.Equal(got.Payload, wire[8:]) {
			t.Fatal("decode", got, err)
		}
		if len(got.Payload) > 0 {
			got.Payload[0] ^= 0xff
			if !bytes.Equal(wire, before) {
				t.Fatal("alias")
			}
		}
	})
}
