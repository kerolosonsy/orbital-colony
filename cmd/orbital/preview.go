package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"orbital/colony/levels"
)

// Level checks use i18n's process language. A short-lived copy of this binary
// builds previews without changing the server's language or launching timers.
func init() {
	if os.Getenv("ORBITAL_PREVIEW") != "1" {
		return
	}
	seed, err := strconv.ParseInt(os.Getenv("ORBITAL_SEED"), 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	levels.Start(os.Getenv("ORBITAL_LEVEL"), os.Stdout, seed)
	os.Exit(0)
}

func previewSnapshot(ctx context.Context, level, lang string) ([]byte, error) {
	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, executable)
	configureProcess(cmd)
	cmd.Env = append(os.Environ(), "ORBITAL_PREVIEW=1", "ORBITAL_LANG="+lang,
		"ORBITAL_LEVEL="+level, "ORBITAL_SEED="+strconv.FormatInt(time.Now().UnixNano(), 10))
	var snapshot, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &snapshot, &stderr
	if err := runCommand(cmd); err != nil {
		return nil, fmt.Errorf("preview: %w: %s", err, stderr.String())
	}
	return bytes.TrimSpace(snapshot.Bytes()), nil
}
