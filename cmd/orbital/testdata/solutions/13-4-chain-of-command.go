package main

import (
	"context"
	"orbital/colony"
)

func command(ctx context.Context) error {
	colony.TraceCommand(ctx, "command.enter")
	err := relay(ctx)
	colony.TraceCommand(ctx, "command.exit")
	return err
}
func relay(ctx context.Context) error {
	colony.TraceCommand(ctx, "relay.enter")
	err := probe(ctx)
	colony.TraceCommand(ctx, "relay.exit")
	return err
}
func probe(ctx context.Context) error {
	colony.TraceCommand(ctx, "probe.enter")
	select {
	case <-ctx.Done():
	}
	colony.TraceCommand(ctx, "probe.exit")
	return ctx.Err()
}
func main() {
	canceled, stop := context.WithCancel(context.Background())
	stop()
	first := command(canceled)
	ctx, cancel := context.WithTimeout(context.Background(), colony.Ticks(colony.Reading(0)))
	defer cancel()
	second := command(ctx)
	if first == context.Canceled && second == context.DeadlineExceeded {
		colony.Report("canceled deadline")
	}
}
