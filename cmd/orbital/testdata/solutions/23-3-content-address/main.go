package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func vaultKind(kind string) bool { return kind == "cargo" || kind == "tree" || kind == "snapshot" }
func vaultID(raw []byte) string {
	h := sha1.New()
	h.Write(raw)
	return fmt.Sprintf("%x", h.Sum(nil))
}
func vaultEnvelope(kind string, payload []byte) []byte {
	return append([]byte(fmt.Sprintf("OC23 %s %d\n", kind, len(payload))), payload...)
}
func validVaultID(id string) bool {
	if len(id) != 40 {
		return false
	}
	for _, c := range id {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func ReadVault(root, id string) (string, []byte, error) {
	if !validVaultID(id) {
		return "", nil, fmt.Errorf("invalid object ID")
	}
	packed, err := os.ReadFile(filepath.Join(root, "vault-objects", id+".z"))
	if err != nil {
		return "", nil, err
	}
	z, err := zlib.NewReader(bytes.NewReader(packed))
	if err != nil {
		return "", nil, err
	}
	raw, readErr := io.ReadAll(io.LimitReader(z, 1048641))
	closeErr := z.Close()
	if readErr != nil {
		return "", nil, readErr
	}
	if closeErr != nil {
		return "", nil, closeErr
	}
	return decodeVault(raw, id)
}
func decodeVault(raw []byte, id string) (string, []byte, error) {
	if len(raw) > 1048640 || vaultID(raw) != id {
		return "", nil, fmt.Errorf("object integrity failure")
	}
	end := bytes.IndexByte(raw, '\n')
	if end < 0 {
		return "", nil, fmt.Errorf("missing envelope")
	}
	fields := strings.Split(string(raw[:end]), " ")
	if len(fields) != 3 || !vaultKind(fields[1]) {
		return "", nil, fmt.Errorf("invalid envelope")
	}
	payload := raw[end+1:]
	if len(payload) > 1048576 || !bytes.Equal(raw, vaultEnvelope(fields[1], payload)) {
		return "", nil, fmt.Errorf("noncanonical envelope")
	}
	return fields[1], payload, nil
}
func publishVault(file string, packed []byte) error {
	f, err := os.CreateTemp(filepath.Dir(file), ".vault-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()
	if _, err := f.Write(packed); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), file)
}
func WriteVault(root, kind string, payload []byte) (string, error) {
	if !vaultKind(kind) || len(payload) > 1048576 {
		return "", fmt.Errorf("invalid kind or payload size")
	}
	raw := vaultEnvelope(kind, payload)
	id := vaultID(raw)
	oldKind, old, err := ReadVault(root, id)
	if err == nil {
		if oldKind != kind || !bytes.Equal(old, payload) {
			return "", fmt.Errorf("object collision")
		}
		return id, nil
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return "", err
	}
	var packed bytes.Buffer
	z := zlib.NewWriter(&packed)
	if _, err := z.Write(raw); err != nil {
		z.Close()
		return "", err
	}
	if err := z.Close(); err != nil {
		return "", err
	}
	dir := filepath.Join(root, "vault-objects")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	if err := publishVault(filepath.Join(dir, id+".z"), packed.Bytes()); err != nil {
		return "", err
	}
	return id, nil
}
func main() {}
