package main

import (
	"bytes"
	"slices"
	"testing"
)

func FuzzCapsule(f *testing.F) {
	for _, frame := range [][]byte{nil, {}, {0xA7}, {0xA7, 1}, {0xA7, 1, 0}, {0xA7, 2, 1, 42}, {0, 1, 0}, {0xA7, 3, 0}, {0xA7, 1, 0, 99}, {0xA7, 1, 17}} {
		f.Add(frame)
	}
	f.Fuzz(func(t *testing.T, frame []byte) {
		before := slices.Clone(frame)
		valid := len(frame) >= 3 && frame[0] == 0xA7 && (frame[1] == 1 || frame[1] == 2) && int(frame[2]) <= 16 && len(frame) == 3+int(frame[2])
		kind, payload, err := ParseCapsule(frame)
		if !bytes.Equal(frame, before) {
			t.Fatal("changed input")
		}
		if !valid {
			if err == nil || kind != 0 || payload != nil {
				t.Fatal("accepted malformed frame")
			}
			return
		}
		if err != nil || kind != frame[1] || !bytes.Equal(payload, frame[3:]) {
			t.Fatal("decode contract")
		}
		if len(payload) > 0 {
			payload[0]++
			if !bytes.Equal(frame, before) {
				t.Fatal("payload aliases input")
			}
		}
	})
}
