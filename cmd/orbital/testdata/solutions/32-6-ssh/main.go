package main

import (
	"bytes"
	"context"
	"errors"
	"net"

	"golang.org/x/crypto/ssh"
)

var ErrDockSSHPlan = errors.New("invalid dock SSH request")

type DockSSHRequest struct {
	User    string
	HostKey ssh.PublicKey
	Auth    []ssh.AuthMethod
	Command string
	Dial    func(context.Context) (net.Conn, error)
}
type DockSSHReply struct {
	Stdout, Stderr string
	ExitCode       int
}

func dockSSHCommand(req DockSSHRequest) (string, error) {
	if len(req.User) < 1 || len(req.User) > 16 || req.HostKey == nil || len(req.Auth) == 0 || req.Dial == nil {
		return "", ErrDockSSHPlan
	}
	for _, b := range []byte(req.User) {
		if !(b >= 'a' && b <= 'z' || b >= '0' && b <= '9' || b == '-') {
			return "", ErrDockSSHPlan
		}
	}
	switch req.Command {
	case "pressure":
		return "colony pressure --brief", nil
	case "seal":
		return "colony seal --check", nil
	default:
		return "", ErrDockSSHPlan
	}
}
func dockSSHSession(client *ssh.Client, command string) (DockSSHReply, error) {
	session, err := client.NewSession()
	if err != nil {
		return DockSSHReply{}, err
	}
	defer session.Close()
	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	err = session.Run(command)
	reply := DockSSHReply{Stdout: stdout.String(), Stderr: stderr.String()}
	if err == nil {
		return reply, nil
	}
	var exit *ssh.ExitError
	if errors.As(err, &exit) {
		reply.ExitCode = exit.ExitStatus()
		return reply, err
	}
	return DockSSHReply{}, err
}
func ExecuteDockCommand(ctx context.Context, req DockSSHRequest) (DockSSHReply, error) {
	command, err := dockSSHCommand(req)
	if err != nil {
		return DockSSHReply{}, err
	}
	if ctx.Err() != nil {
		return DockSSHReply{}, ctx.Err()
	}
	raw, err := req.Dial(ctx)
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return DockSSHReply{}, err
	}
	defer raw.Close()
	stop := make(chan struct{})
	watched := make(chan struct{})
	go func() {
		defer close(watched)
		select {
		case <-ctx.Done():
			raw.Close()
		case <-stop:
		}
	}()
	defer func() { close(stop); <-watched }()
	config := &ssh.ClientConfig{User: req.User, Auth: req.Auth, HostKeyCallback: ssh.FixedHostKey(req.HostKey)}
	conn, channels, requests, err := ssh.NewClientConn(raw, "dock", config)
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		return DockSSHReply{}, err
	}
	client := ssh.NewClient(conn, channels, requests)
	defer func() { client.Close(); client.Wait() }()
	reply, err := dockSSHSession(client, command)
	if ctx.Err() != nil {
		return DockSSHReply{}, ctx.Err()
	}
	return reply, err
}
func main() {}
