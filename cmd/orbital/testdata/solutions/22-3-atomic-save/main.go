package main

import (
	"io"
	"os"
	"path/filepath"
)

func SaveRelay(target string, write func(io.Writer) error) error {
	return CommitRelay(target, write, os.Rename)
}
func CommitRelay(target string, write func(io.Writer) error, rename func(string, string) error) error {
	file, err := os.CreateTemp(filepath.Dir(target), ".relay-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if err := write(file); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return rename(file.Name(), target)
}
func main() {}
