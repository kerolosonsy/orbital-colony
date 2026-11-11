//go:build unix

package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestCancellationStopsChildrenHoldingOutput(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "child-survived")
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", `(sleep 1; printf child > "$ORBITAL_TEST_MARKER") & wait`)
	cmd.Env = append(os.Environ(), "ORBITAL_TEST_MARKER="+marker)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	configureProcess(cmd)
	start := time.Now()
	err := runCommand(cmd)
	if err == nil || !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("expected cancellation, got %v (context %v)", err, ctx.Err())
	}
	if time.Since(start) > 900*time.Millisecond {
		t.Fatal("a child kept the output pipe open after cancellation")
	}
	// A surviving child would write after the parent has already returned.
	time.Sleep(time.Until(start.Add(1200 * time.Millisecond)))
	if _, err := os.Stat(marker); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("child survived cancellation: stat error %v", err)
	}
}
