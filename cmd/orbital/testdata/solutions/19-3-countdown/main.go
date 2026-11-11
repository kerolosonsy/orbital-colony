package main

import (
	"context"
	"fmt"
	"orbital/colony/labkit"
	"time"
)

func Countdown(ctx context.Context, clock labkit.Clock, seconds int, emit func(int) error) error {
	if clock == nil || seconds < 0 || emit == nil {
		return fmt.Errorf("invalid countdown")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := emit(seconds); err != nil {
		return err
	}
	if seconds == 0 {
		return nil
	}
	ticker := clock.NewTicker(time.Second)
	defer ticker.Stop()
	for seconds > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C():
			if err := ctx.Err(); err != nil {
				return err
			}
			seconds--
			if err := emit(seconds); err != nil {
				return err
			}
		}
	}
	return nil
}
func main() {}
