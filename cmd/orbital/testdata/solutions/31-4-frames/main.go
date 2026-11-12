package main

import (
	"encoding/binary"
	"errors"
)

var ErrRadiogram = errors.New("invalid radiogram")
var ErrRadiogramKind = errors.New("unknown radiogram kind")

type Radiogram struct {
	Kind     byte
	Sequence uint16
	Payload  []byte
}

func ParseRadiogram(wire []byte) (Radiogram, error) {
	if len(wire) < 8 || wire[0] != 'O' || wire[1] != 'C' || wire[2] != 1 {
		return Radiogram{}, ErrRadiogram
	}
	length := int(binary.BigEndian.Uint16(wire[4:6]))
	if length > 1024 || len(wire) != 8+length {
		return Radiogram{}, ErrRadiogram
	}
	if wire[3] != 1 && wire[3] != 2 {
		return Radiogram{}, ErrRadiogramKind
	}
	return Radiogram{Kind: wire[3], Sequence: binary.BigEndian.Uint16(wire[6:8]), Payload: append([]byte{}, wire[8:]...)}, nil
}
func main() {}
