package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/segmentio/kafka-go"
)

type Telemetry struct {
	Rover   string
	Reading int
}

var ErrTelemetry = errors.New("invalid telemetry")

type Writer interface {
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Reader interface {
	FetchMessage(ctx context.Context) (kafka.Message, error)
	CommitMessages(ctx context.Context, msgs ...kafka.Message) error
}

type Handler func(ctx context.Context, t Telemetry) error

const sendAttempts = 3

func Send(ctx context.Context, w Writer, batch []Telemetry) error {
	msgs := make([]kafka.Message, 0, len(batch))
	for _, t := range batch {
		if t.Rover == "" {
			return ErrTelemetry
		}
		msgs = append(msgs, kafka.Message{Key: []byte(t.Rover), Value: []byte(strconv.Itoa(t.Reading))})
	}
	if len(msgs) == 0 {
		return nil
	}
	var err error
	for attempt := 0; attempt < sendAttempts; attempt++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if err = w.WriteMessages(ctx, msgs...); err == nil {
			return nil
		}
	}
	return fmt.Errorf("send failed after %d attempts: %w", sendAttempts, err)
}

func Consume(ctx context.Context, r Reader, h Handler, maxAttempts int) (int, error) {
	handled := 0
	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return handled, nil
			}
			return handled, err
		}
		reading, convErr := strconv.Atoi(string(msg.Value))
		if convErr == nil {
			t := Telemetry{Rover: string(msg.Key), Reading: reading}
			var hErr error
			for attempt := 0; attempt < maxAttempts; attempt++ {
				if hErr = h(ctx, t); hErr == nil || ctx.Err() != nil {
					break
				}
			}
			if hErr != nil {
				if ctx.Err() != nil {
					return handled, nil
				}
				return handled, hErr
			}
			handled++
		}
		// Finish the commit even while shutting down: the work is already done.
		if err := r.CommitMessages(context.WithoutCancel(ctx), msg); err != nil {
			return handled, err
		}
	}
}

func main() {}
