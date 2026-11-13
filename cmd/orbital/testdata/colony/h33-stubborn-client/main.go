package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "orbital/colony"
)

func sell(client *http.Client, key string) error {
	for attempt := 0; attempt < 16; attempt++ {
		req, err := http.NewRequest(http.MethodPost, os.Getenv("ORBITAL_MARKET_URL")+"/orders", strings.NewReader(`{"good":"moss","amount":1}`))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", key)
		resp, err := client.Do(req)
		pause := time.Duration(25*(1<<min(attempt, 3))) * time.Millisecond
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
			if resp.StatusCode != 429 && resp.StatusCode != 503 {
				return fmt.Errorf("order %s: %s", key, resp.Status)
			}
			if seconds, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && seconds > 0 {
				pause = time.Duration(seconds) * time.Second
			}
		}
		time.Sleep(pause)
	}
	return fmt.Errorf("order %s exhausted retries", key)
}

func main() {
	client := &http.Client{Timeout: 150 * time.Millisecond}
	defer client.CloseIdleConnections()
	jobs := make(chan string)
	failures := make(chan error, 20)
	var wg sync.WaitGroup
	for worker := 0; worker < 4; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range jobs {
				if err := sell(client, key); err != nil {
					failures <- err
				}
			}
		}()
	}
	for n := 0; n < 20; n++ {
		jobs <- fmt.Sprintf("order-%02d", n)
	}
	close(jobs)
	wg.Wait()
	close(failures)
	for err := range failures {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
