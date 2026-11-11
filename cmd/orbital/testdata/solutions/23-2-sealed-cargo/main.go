package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func cargoCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
func SealCargo(key, manifest, cargo []byte, source io.Reader) ([]byte, error) {
	aead, err := cargoCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cargo) > 65536 {
		return nil, fmt.Errorf("cargo exceeds 65536 bytes")
	}
	if source == nil {
		source = rand.Reader
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(source, nonce); err != nil {
		return nil, err
	}
	packet := append([]byte("OCV1"), nonce...)
	aad := append([]byte("OCV1"), manifest...)
	return aead.Seal(packet, nonce, cargo, aad), nil
}
func OpenCargo(key, manifest, packet []byte) ([]byte, error) {
	aead, err := cargoCipher(key)
	if err != nil {
		return nil, err
	}
	if len(packet) < 32 || len(packet) > 65568 || string(packet[:4]) != "OCV1" {
		return nil, fmt.Errorf("invalid cargo packet")
	}
	aad := append([]byte("OCV1"), manifest...)
	plain, err := aead.Open(nil, packet[4:16], packet[16:], aad)
	if err != nil {
		return nil, errors.New("authentication failed")
	}
	return plain, nil
}
func OpenCargoTo(key, manifest, packet []byte, out io.Writer) error {
	plain, err := OpenCargo(key, manifest, packet)
	if err != nil {
		return err
	}
	n, err := out.Write(plain)
	if err == nil && n != len(plain) {
		return io.ErrShortWrite
	}
	return err
}
func main() {}
