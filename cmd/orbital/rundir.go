package main

import "os"

// runDir makes a fresh temporary folder for one player program.
func runDir(pattern string) (string, error) {
	return os.MkdirTemp("", pattern)
}
