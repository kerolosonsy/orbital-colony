package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sort"
)

var ErrDepotConfig = errors.New("invalid depot config")
var ErrChunkLost = errors.New("chunk lost")

type Pod interface {
	Store(id string, data []byte) error
	Fetch(id string) ([]byte, error)
}

type Depot struct {
	pods      []Pod
	chunkSize int
	replicas  int
	stored    map[string]bool
}

func NewDepot(pods []Pod, chunkSize, replicas int) (*Depot, error) {
	if len(pods) == 0 || chunkSize < 1 || chunkSize > 1<<20 || replicas < 1 || replicas > len(pods) {
		return nil, ErrDepotConfig
	}
	for _, p := range pods {
		if p == nil {
			return nil, ErrDepotConfig
		}
	}
	return &Depot{pods: append([]Pod{}, pods...), chunkSize: chunkSize, replicas: replicas, stored: map[string]bool{}}, nil
}

func (d *Depot) placement(sum []byte) []Pod {
	start := int(binary.BigEndian.Uint32(sum[:4]) % uint32(len(d.pods)))
	out := make([]Pod, d.replicas)
	for j := range out {
		out[j] = d.pods[(start+j)%len(d.pods)]
	}
	return out
}

func (d *Depot) Save(r io.Reader) ([]string, error) {
	ids := []string{}
	buf := make([]byte, d.chunkSize)
	for {
		n, err := io.ReadFull(r, buf)
		if err == io.EOF {
			return ids, nil
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			return ids, err
		}
		chunk := buf[:n]
		sum := sha256.Sum256(chunk)
		id := hex.EncodeToString(sum[:])
		if !d.stored[id] {
			for _, p := range d.placement(sum[:]) {
				if serr := p.Store(id, chunk); serr != nil {
					return ids, serr
				}
			}
			d.stored[id] = true
		}
		ids = append(ids, id)
		if err == io.ErrUnexpectedEOF {
			return ids, nil
		}
	}
}

func intact(id string, b []byte) bool {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]) == id
}

func decodeID(id string) ([]byte, bool) {
	if len(id) != 64 {
		return nil, false
	}
	b, err := hex.DecodeString(id)
	return b, err == nil && hex.EncodeToString(b) == id
}

// healthy returns the first intact copy among the chunk's placements.
func (d *Depot) healthy(id string) ([]byte, []Pod, error) {
	sum, ok := decodeID(id)
	if !ok {
		return nil, nil, fmt.Errorf("%w: bad id %q", ErrChunkLost, id)
	}
	pl := d.placement(sum)
	for _, p := range pl {
		if b, err := p.Fetch(id); err == nil && intact(id, b) {
			return b, pl, nil
		}
	}
	return nil, pl, fmt.Errorf("%w: %s", ErrChunkLost, id)
}

func (d *Depot) Load(ids []string, w io.Writer) error {
	for _, id := range ids {
		b, _, err := d.healthy(id)
		if err != nil {
			return err
		}
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func (d *Depot) Repair() (int, error) {
	ids := make([]string, 0, len(d.stored))
	for id := range d.stored {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	fixed := 0
	var errs []error
	for _, id := range ids {
		good, pl, err := d.healthy(id)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		for _, p := range pl {
			if b, err := p.Fetch(id); err == nil && intact(id, b) {
				continue
			}
			if err := p.Store(id, good); err != nil {
				errs = append(errs, err)
				continue
			}
			fixed++
		}
	}
	return fixed, errors.Join(errs...)
}

func main() {}
