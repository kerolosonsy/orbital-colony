package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

type BeaconReply struct {
	Status int
	Err    error
}

func checkBeacon(ctx context.Context, client *http.Client, address string) BeaconReply {
	if err := ctx.Err(); err != nil {
		return BeaconReply{Err: err}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, address, nil)
	if err != nil {
		return BeaconReply{Err: err}
	}
	if req.URL.Host == "" || (req.URL.Scheme != "http" && req.URL.Scheme != "https") {
		return BeaconReply{Err: fmt.Errorf("invalid beacon address")}
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return BeaconReply{Err: err}
	}
	reply := BeaconReply{Status: resp.StatusCode}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		reply.Err = fmt.Errorf("beacon status %d", resp.StatusCode)
	}
	return reply
}
func CheckBeacons(ctx context.Context, client *http.Client, addresses []string, limit int) ([]BeaconReply, error) {
	if client == nil || limit <= 0 {
		return nil, fmt.Errorf("invalid beacon check")
	}
	replies := make([]BeaconReply, len(addresses))
	slots := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for i, address := range addresses {
		if err := ctx.Err(); err != nil {
			replies[i].Err = err
			continue
		}
		select {
		case slots <- struct{}{}:
			wg.Add(1)
			go func(i int, address string) {
				defer wg.Done()
				defer func() { <-slots }()
				replies[i] = checkBeacon(ctx, client, address)
			}(i, address)
		case <-ctx.Done():
			replies[i].Err = ctx.Err()
		}
	}
	wg.Wait()
	return replies, nil
}
func main() {}
