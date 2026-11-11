package main

import (
	"errors"
	"fmt"
	"orbital/colony"
)

func main() {
	var all error
	for i := 0; i < colony.DiagnosticCount(); i++ {
		err := colony.Diagnose(i)
		if err != nil {
			wrapped := fmt.Errorf("bay %d: %w", i, err)
			if errors.Is(wrapped, err) {
				colony.Report(colony.FaultCode(wrapped))
			}
			all = errors.Join(all, wrapped)
		}
	}
	if all != nil {
		colony.Report(all)
	}
}
