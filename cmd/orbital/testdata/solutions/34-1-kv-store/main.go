package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"
)

var ErrVaultKey = errors.New("empty vault key")
var ErrVaultTTL = errors.New("negative vault ttl")
var ErrVaultSnapshot = errors.New("invalid vault snapshot")

type vaultEntry struct {
	value   []byte
	expires time.Time
}

type Vault struct {
	now   func() time.Time
	mu    sync.RWMutex
	items map[string]vaultEntry
}

func NewVault(now func() time.Time) *Vault {
	return &Vault{now: now, items: map[string]vaultEntry{}}
}

func (e vaultEntry) alive(now time.Time) bool {
	return e.expires.IsZero() || now.Before(e.expires)
}

func (v *Vault) Put(key string, value []byte, ttl time.Duration) error {
	if key == "" {
		return ErrVaultKey
	}
	if ttl < 0 {
		return ErrVaultTTL
	}
	e := vaultEntry{value: append([]byte(nil), value...)}
	v.mu.Lock()
	defer v.mu.Unlock()
	if ttl > 0 {
		e.expires = v.now().Add(ttl)
	}
	v.items[key] = e
	return nil
}

func (v *Vault) Get(key string) ([]byte, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	e, ok := v.items[key]
	if !ok || !e.alive(v.now()) {
		return nil, false
	}
	return append([]byte{}, e.value...), true
}

func (v *Vault) Delete(key string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	e, ok := v.items[key]
	delete(v.items, key)
	return ok && e.alive(v.now())
}

func (v *Vault) Len() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	now, n := v.now(), 0
	for _, e := range v.items {
		if e.alive(now) {
			n++
		}
	}
	return n
}

type vaultRow struct {
	Key     string `json:"key"`
	Value   []byte `json:"value"`
	Expires int64  `json:"expires"`
}

func (v *Vault) Snapshot(w io.Writer) error {
	v.mu.RLock()
	now := v.now()
	rows := []vaultRow{}
	for k, e := range v.items {
		if !e.alive(now) {
			continue
		}
		r := vaultRow{Key: k, Value: e.value}
		if !e.expires.IsZero() {
			r.Expires = e.expires.UnixNano()
		}
		rows = append(rows, r)
	}
	v.mu.RUnlock()
	sort.Slice(rows, func(i, j int) bool { return rows[i].Key < rows[j].Key })
	return json.NewEncoder(w).Encode(rows)
}

func LoadVault(r io.Reader, now func() time.Time) (*Vault, error) {
	var rows []vaultRow
	if err := json.NewDecoder(r).Decode(&rows); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrVaultSnapshot, err)
	}
	v := NewVault(now)
	seen := map[string]bool{}
	for _, row := range rows {
		if row.Key == "" || seen[row.Key] {
			return nil, fmt.Errorf("%w: bad key %q", ErrVaultSnapshot, row.Key)
		}
		seen[row.Key] = true
		e := vaultEntry{value: row.Value}
		if row.Expires != 0 {
			e.expires = time.Unix(0, row.Expires)
		}
		if e.alive(now()) {
			v.items[row.Key] = e
		}
	}
	return v, nil
}

func main() {}
