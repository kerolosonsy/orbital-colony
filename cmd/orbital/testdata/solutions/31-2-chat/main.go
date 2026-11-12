package main

import (
	"context"
	"errors"
	"strings"
)

var ErrRadioName = errors.New("invalid radio name")
var ErrRadioMember = errors.New("unknown radio member")
var ErrRadioCommand = errors.New("invalid radio command")

type RadioCommand struct {
	Kind, Name, Text string
	Out              chan string
	Reply            chan error
}

func radioName(name string) bool {
	if len(name) < 1 || len(name) > 16 {
		return false
	}
	for _, b := range []byte(name) {
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-') {
			return false
		}
	}
	return true
}
func radioApply(members map[string]chan string, c RadioCommand) error {
	switch c.Kind {
	case "join":
		if !radioName(c.Name) || members[c.Name] != nil || cap(c.Out) == 0 {
			return ErrRadioName
		}
		members[c.Name] = c.Out
	case "leave":
		mailbox, ok := members[c.Name]
		if !ok {
			return ErrRadioMember
		}
		delete(members, c.Name)
		close(mailbox)
	case "say":
		if members[c.Name] == nil {
			return ErrRadioMember
		}
		if len(c.Text) < 1 || len(c.Text) > 128 || strings.ContainsAny(c.Text, "\r\n") {
			return ErrRadioCommand
		}
		for name, mailbox := range members {
			if name == c.Name {
				continue
			}
			select {
			case mailbox <- c.Name + ": " + c.Text:
			default:
				delete(members, name)
				close(mailbox)
			}
		}
	default:
		return ErrRadioCommand
	}
	return nil
}
func RunRadio(ctx context.Context, commands <-chan RadioCommand) error {
	members := make(map[string]chan string)
	defer func() {
		for _, mailbox := range members {
			close(mailbox)
		}
	}()
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case c, ok := <-commands:
			if !ok {
				return nil
			}
			c.Reply <- radioApply(members, c)
		}
	}
}
func main() {}
