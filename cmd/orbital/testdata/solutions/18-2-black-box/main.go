package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type BlackBox struct {
	Cycle     int      `json:"cycle"`
	Callsigns []string `json:"callsigns"`
}

func validBlackBox(box BlackBox) error {
	if box.Cycle <= 0 || box.Callsigns == nil {
		return fmt.Errorf("invalid black box schema")
	}
	for _, name := range box.Callsigns {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("invalid callsign")
		}
	}
	return nil
}
func WriteBlackBox(path string, box BlackBox) error {
	if err := validBlackBox(box); err != nil {
		return err
	}
	encoded, err := json.Marshal(box)
	if err != nil {
		return fmt.Errorf("encode %s: %w", path, err)
	}
	if err := os.WriteFile(path, encoded, 0600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
func ReadBlackBox(path string) (BlackBox, error) {
	encoded, err := os.ReadFile(path)
	if err != nil {
		return BlackBox{}, fmt.Errorf("read %s: %w", path, err)
	}
	var box BlackBox
	if err := json.Unmarshal(encoded, &box); err != nil {
		return BlackBox{}, fmt.Errorf("decode %s: %w", path, err)
	}
	if err := validBlackBox(box); err != nil {
		return BlackBox{}, fmt.Errorf("validate %s: %w", path, err)
	}
	return box, nil
}
func main() {}
