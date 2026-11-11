package main

import "os/exec"

func runCommand(cmd *exec.Cmd) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	defer killProcessTree(cmd)
	return cmd.Wait()
}
