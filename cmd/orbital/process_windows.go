package main

import (
	"context"
	"os/exec"
	"strconv"
	"time"
)

func configureProcess(cmd *exec.Cmd) {
	cmd.Cancel = func() error { return killProcessTree(cmd) }
	cmd.WaitDelay = 250 * time.Millisecond
}

func killProcessTree(cmd *exec.Cmd) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "taskkill", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
}
