package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
)

var ErrLedgerKey = errors.New("invalid ledger key")
var ErrLedgerValue = errors.New("ledger value too large")
var ErrLedgerCorrupt = errors.New("corrupt ledger")

const (
	opSet    = 1
	opDelete = 2
	header   = 13
	maxKey   = 1024
	maxValue = 65536
)

type Ledger struct {
	f   *os.File
	mem map[string][]byte
}

func encode(op byte, key string, value []byte) []byte {
	r := make([]byte, header, header+len(key)+len(value))
	r[4] = op
	binary.BigEndian.PutUint32(r[5:9], uint32(len(key)))
	binary.BigEndian.PutUint32(r[9:13], uint32(len(value)))
	r = append(r, key...)
	r = append(r, value...)
	binary.BigEndian.PutUint32(r[0:4], crc32.ChecksumIEEE(r[4:]))
	return r
}

func OpenLedger(dir string) (*Ledger, error) {
	path := filepath.Join(dir, "ledger.wal")
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	mem := map[string][]byte{}
	off := 0
	for off < len(data) {
		rest := data[off:]
		if len(rest) < header {
			break
		}
		op := rest[4]
		klen := binary.BigEndian.Uint32(rest[5:9])
		vlen := binary.BigEndian.Uint32(rest[9:13])
		if klen == 0 || klen > maxKey || vlen > maxValue {
			return nil, fmt.Errorf("%w: bad lengths at %d", ErrLedgerCorrupt, off)
		}
		size := header + int(klen) + int(vlen)
		if len(rest) < size {
			break
		}
		if binary.BigEndian.Uint32(rest[0:4]) != crc32.ChecksumIEEE(rest[4:size]) {
			return nil, fmt.Errorf("%w: checksum at %d", ErrLedgerCorrupt, off)
		}
		key := string(rest[header : header+klen])
		switch {
		case op == opSet:
			mem[key] = append([]byte{}, rest[header+klen:size]...)
		case op == opDelete && vlen == 0:
			delete(mem, key)
		default:
			return nil, fmt.Errorf("%w: bad op at %d", ErrLedgerCorrupt, off)
		}
		off += size
	}
	if off < len(data) {
		if err := os.Truncate(path, int64(off)); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return &Ledger{f: f, mem: mem}, nil
}

func (l *Ledger) write(op byte, key string, value []byte) error {
	if key == "" || len(key) > maxKey {
		return ErrLedgerKey
	}
	if len(value) > maxValue {
		return ErrLedgerValue
	}
	_, err := l.f.Write(encode(op, key, value))
	return err
}

func (l *Ledger) Set(key string, value []byte) error {
	if err := l.write(opSet, key, value); err != nil {
		return err
	}
	l.mem[key] = append([]byte{}, value...)
	return nil
}

func (l *Ledger) Remove(key string) error {
	if err := l.write(opDelete, key, nil); err != nil {
		return err
	}
	delete(l.mem, key)
	return nil
}

func (l *Ledger) Lookup(key string) ([]byte, bool) {
	v, ok := l.mem[key]
	if !ok {
		return nil, false
	}
	return append([]byte{}, v...), true
}

func (l *Ledger) Close() error { return l.f.Close() }

func main() {}
