package main

import (
	"context"
	"errors"
	"net"
	"strconv"
	"sync"
	"time"
)

var ErrPortPlan = errors.New("invalid local port plan")
var ErrPortDial = errors.New("dial returned no connection")

type PortOptions struct {
	Workers int
	Timeout time.Duration
	Dial    func(string, string, time.Duration) (net.Conn, error)
}
type PortStatus struct {
	Port  int
	State string
	Err   error
}

func validPortPlan(host string, ports []int, options PortOptions) net.IP {
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() || len(ports) > 64 || options.Workers < 1 || options.Workers > 16 || options.Timeout <= 0 {
		return nil
	}
	seen := map[int]bool{}
	for _, port := range ports {
		if port < 1 || port > 65535 || seen[port] {
			return nil
		}
		seen[port] = true
	}
	return ip
}
func checkDockPort(host string, port int, options PortOptions) PortStatus {
	conn, err := options.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)), options.Timeout)
	if conn != nil {
		conn.Close()
	}
	if err == nil && conn == nil {
		err = ErrPortDial
	}
	if err == nil {
		return PortStatus{Port: port, State: "open"}
	}
	state := "closed"
	var timeout net.Error
	if errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &timeout) && timeout.Timeout()) {
		state = "timeout"
	}
	return PortStatus{Port: port, State: state, Err: err}
}
func CheckDockPorts(host string, ports []int, options PortOptions) ([]PortStatus, error) {
	ip := validPortPlan(host, ports, options)
	if ip == nil {
		return nil, ErrPortPlan
	}
	if options.Dial == nil {
		options.Dial = net.DialTimeout
	}
	replies := make([]PortStatus, len(ports))
	jobs := make(chan int, len(ports))
	for i := range ports {
		jobs <- i
	}
	close(jobs)
	var wg sync.WaitGroup
	workers := options.Workers
	if workers > len(ports) {
		workers = len(ports)
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				replies[index] = checkDockPort(ip.String(), ports[index], options)
			}
		}()
	}
	wg.Wait()
	return replies, nil
}
func main() {}
