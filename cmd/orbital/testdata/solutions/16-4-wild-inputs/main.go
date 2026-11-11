package main

import (
	"errors"
	"slices"
)

func ParseCapsule(frame []byte) (byte, []byte, error) {
	if len(frame) < 3 {
		return 0, nil, errors.New("short capsule")
	}
	if frame[0] != 0xA7 {
		return 0, nil, errors.New("signature")
	}
	if frame[1] != 1 && frame[1] != 2 {
		return 0, nil, errors.New("kind")
	}
	n := int(frame[2])
	if n > 16 || len(frame) != 3+n {
		return 0, nil, errors.New("length")
	}
	return frame[1], slices.Clone(frame[3:]), nil
}
func main() {}
