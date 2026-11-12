package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

type slot struct {
	value []byte
	dead  bool
}

type record struct {
	key string
	slot
}

type table struct {
	seq  int
	keys []string
	vals []slot
}

type Ledger struct {
	dir    string
	f      *os.File
	mem    map[string]slot
	tables []table // oldest first
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

// parse returns the complete records and the offset after the last one.
func parse(data []byte) ([]record, int, error) {
	var out []record
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
			return nil, 0, fmt.Errorf("%w: bad lengths at %d", ErrLedgerCorrupt, off)
		}
		size := header + int(klen) + int(vlen)
		if len(rest) < size {
			break
		}
		if binary.BigEndian.Uint32(rest[0:4]) != crc32.ChecksumIEEE(rest[4:size]) {
			return nil, 0, fmt.Errorf("%w: checksum at %d", ErrLedgerCorrupt, off)
		}
		r := record{key: string(rest[header : header+klen])}
		switch {
		case op == opSet:
			r.value = append([]byte{}, rest[header+klen:size]...)
		case op == opDelete && vlen == 0:
			r.dead = true
		default:
			return nil, 0, fmt.Errorf("%w: bad op at %d", ErrLedgerCorrupt, off)
		}
		out = append(out, r)
		off += size
	}
	return out, off, nil
}

func tableName(seq int) string { return fmt.Sprintf("sst-%06d.sst", seq) }

func loadTable(path string, seq int) (table, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return table{}, err
	}
	recs, off, err := parse(data)
	if err != nil {
		return table{}, err
	}
	if off != len(data) {
		return table{}, fmt.Errorf("%w: short table %s", ErrLedgerCorrupt, path)
	}
	t := table{seq: seq}
	for i, r := range recs {
		if i > 0 && recs[i-1].key >= r.key {
			return table{}, fmt.Errorf("%w: unsorted table %s", ErrLedgerCorrupt, path)
		}
		t.keys = append(t.keys, r.key)
		t.vals = append(t.vals, r.slot)
	}
	return t, nil
}

func OpenLedger(dir string) (*Ledger, error) {
	l := &Ledger{dir: dir, mem: map[string]slot{}}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, ".tmp") {
			if err := os.Remove(filepath.Join(dir, name)); err != nil {
				return nil, err
			}
			continue
		}
		var seq int
		if !strings.HasSuffix(name, ".sst") {
			continue
		}
		if _, err := fmt.Sscanf(name, "sst-%06d.sst", &seq); err != nil {
			continue
		}
		t, err := loadTable(filepath.Join(dir, name), seq)
		if err != nil {
			return nil, err
		}
		l.tables = append(l.tables, t)
	}
	path := filepath.Join(dir, "ledger.wal")
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	recs, off, err := parse(data)
	if err != nil {
		return nil, err
	}
	for _, r := range recs {
		l.mem[r.key] = r.slot
	}
	if off < len(data) {
		if err := os.Truncate(path, int64(off)); err != nil {
			return nil, err
		}
	}
	l.f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	return l, nil
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
	l.mem[key] = slot{value: append([]byte{}, value...)}
	return nil
}

func (l *Ledger) Remove(key string) error {
	if err := l.write(opDelete, key, nil); err != nil {
		return err
	}
	l.mem[key] = slot{dead: true}
	return nil
}

func (l *Ledger) Lookup(key string) ([]byte, bool) {
	s, ok := l.mem[key]
	for i := len(l.tables) - 1; !ok && i >= 0; i-- {
		t := l.tables[i]
		if j := sort.SearchStrings(t.keys, key); j < len(t.keys) && t.keys[j] == key {
			s, ok = t.vals[j], true
		}
	}
	if !ok || s.dead {
		return nil, false
	}
	return append([]byte{}, s.value...), true
}

func (l *Ledger) nextSeq() int {
	if len(l.tables) == 0 {
		return 1
	}
	return l.tables[len(l.tables)-1].seq + 1
}

// writeTable stores the slots sorted by key atomically and returns the table.
func (l *Ledger) writeTable(slots map[string]slot) (table, error) {
	t := table{seq: l.nextSeq()}
	for k := range slots {
		t.keys = append(t.keys, k)
	}
	sort.Strings(t.keys)
	var data []byte
	for _, k := range t.keys {
		s := slots[k]
		t.vals = append(t.vals, s)
		if s.dead {
			data = append(data, encode(opDelete, k, nil)...)
		} else {
			data = append(data, encode(opSet, k, s.value)...)
		}
	}
	final := filepath.Join(l.dir, tableName(t.seq))
	tmp := final + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return table{}, err
	}
	if err := os.Rename(tmp, final); err != nil {
		return table{}, err
	}
	return t, nil
}

func (l *Ledger) Flush() error {
	if len(l.mem) == 0 {
		return nil
	}
	t, err := l.writeTable(l.mem)
	if err != nil {
		return err
	}
	l.tables = append(l.tables, t)
	if err := l.f.Truncate(0); err != nil {
		return err
	}
	l.mem = map[string]slot{}
	return nil
}

func (l *Ledger) Compact() error {
	if len(l.tables) == 0 {
		return nil
	}
	merged := map[string]slot{}
	for _, t := range l.tables {
		for i, k := range t.keys {
			merged[k] = t.vals[i]
		}
	}
	for k, s := range merged {
		if s.dead {
			delete(merged, k)
		}
	}
	t, err := l.writeTable(merged)
	if err != nil {
		return err
	}
	old := l.tables
	l.tables = []table{t}
	for _, o := range old {
		if err := os.Remove(filepath.Join(l.dir, tableName(o.seq))); err != nil {
			return err
		}
	}
	return nil
}

func (l *Ledger) Close() error { return l.f.Close() }

func main() {}
