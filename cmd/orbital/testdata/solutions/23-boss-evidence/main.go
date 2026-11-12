package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type EvidenceSample struct {
	Sensor string
	Value  int
}
type EvidenceBundle struct {
	Captured string
	Samples  []EvidenceSample
}

var ErrEvidence = errors.New("invalid evidence")

func SealEvidence(in io.Reader, dst string, key []byte, entropy io.Reader, at time.Time) error {
	reader := csv.NewReader(in)
	header, err := reader.Read()
	if err != nil {
		return err
	}
	if len(header) != 2 || header[0] != "sensor" || header[1] != "value" {
		return ErrEvidence
	}
	totals := map[string]int{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(row) != 2 || row[0] == "" {
			return ErrEvidence
		}
		value, err := strconv.Atoi(row[1])
		if err != nil || value < 0 {
			return ErrEvidence
		}
		totals[row[0]] += value
	}
	samples := make([]EvidenceSample, 0, len(totals))
	for sensor, value := range totals {
		samples = append(samples, EvidenceSample{sensor, value})
	}
	sort.Slice(samples, func(i, j int) bool { return samples[i].Sensor < samples[j].Sensor })
	plain, err := json.Marshal(EvidenceBundle{at.UTC().Format(time.RFC3339), samples})
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(entropy, nonce); err != nil {
		return err
	}
	sealed := gcm.Seal(nonce, nonce, plain, []byte("colony:evidence:v1"))
	file, err := os.CreateTemp(filepath.Dir(dst), ".evidence-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err = file.Write(sealed); err != nil {
		file.Close()
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), dst)
}
func main() {}
