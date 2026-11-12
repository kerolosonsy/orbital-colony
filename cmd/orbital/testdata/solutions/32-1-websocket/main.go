package main

import (
	"context"
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/coder/websocket"
)

type Pulse struct {
	Context context.Context
	Reply   chan error
}
type RelayOptions struct {
	Origin string
	Limit  int
	Pulses <-chan Pulse
}
type relayMessage struct {
	text []byte
	code websocket.StatusCode
	err  error
}

func readTelemetry(ctx context.Context, conn *websocket.Conn, inbox chan<- relayMessage) {
	for {
		kind, text, err := conn.Read(ctx)
		msg := relayMessage{text: text, err: err}
		if err == nil {
			switch {
			case kind != websocket.MessageText:
				msg.code = websocket.StatusUnsupportedData
			case !utf8.Valid(text):
				msg.code = websocket.StatusInvalidFramePayloadData
			case len(text) == 0:
				msg.code = websocket.StatusPolicyViolation
			}
		}
		select {
		case inbox <- msg:
		case <-ctx.Done():
			return
		}
		if err != nil || msg.code != 0 {
			return
		}
	}
}
func writeTelemetry(ctx context.Context, conn *websocket.Conn, inbox <-chan relayMessage, pulses <-chan Pulse) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pulse, ok := <-pulses:
			if !ok {
				pulses = nil
				continue
			}
			err := conn.Ping(pulse.Context)
			pulse.Reply <- err
			if err != nil {
				return err
			}
		case msg := <-inbox:
			if msg.err != nil {
				if websocket.CloseStatus(msg.err) == websocket.StatusNormalClosure {
					return nil
				}
				return msg.err
			}
			if msg.code != 0 {
				conn.Close(msg.code, "invalid telemetry")
				return errors.New("invalid telemetry")
			}
			if err := conn.Write(ctx, websocket.MessageText, append([]byte("relay:"), msg.text...)); err != nil {
				return err
			}
		}
	}
}
func RelayTelemetry(ctx context.Context, w http.ResponseWriter, r *http.Request, opts RelayOptions) error {
	if opts.Limit < 1 || opts.Limit > 4096 || opts.Origin == "" {
		http.Error(w, "invalid options", http.StatusBadRequest)
		return errors.New("invalid options")
	}
	if r.Header.Get("Origin") != opts.Origin {
		http.Error(w, "origin rejected", http.StatusForbidden)
		return errors.New("origin rejected")
	}
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}
	conn.SetReadLimit(int64(opts.Limit))
	ctx, cancel := context.WithCancel(ctx)
	inbox := make(chan relayMessage, 1)
	stopped := make(chan struct{})
	go func() { defer close(stopped); readTelemetry(ctx, conn, inbox) }()
	err = writeTelemetry(ctx, conn, inbox, opts.Pulses)
	cancel()
	conn.CloseNow()
	<-stopped
	return err
}
func main() {}
