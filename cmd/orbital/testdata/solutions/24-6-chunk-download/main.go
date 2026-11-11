package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type ChartPlan struct {
	Address     string
	Size, Chunk int64
	Workers     int
}
type ChartStage interface {
	io.WriterAt
	Commit() error
	Abort() error
}

func chartPart(ctx context.Context, client *http.Client, plan ChartPlan, start int64) ([]byte, bool, error) {
	end := start + min(plan.Chunk, plan.Size-start) - 1
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, plan.Address, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, false, err
	}
	whole := start == 0 && resp.StatusCode == http.StatusOK
	expected := end - start + 1
	if whole {
		expected = plan.Size
	} else {
		header := fmt.Sprintf("bytes %d-%d/%d", start, end, plan.Size)
		if resp.StatusCode != http.StatusPartialContent || resp.Header.Get("Content-Range") != header {
			return nil, false, fmt.Errorf("invalid chart range response")
		}
	}
	payload, err := io.ReadAll(io.LimitReader(resp.Body, expected+1))
	if err != nil {
		return nil, false, err
	}
	if int64(len(payload)) != expected {
		return nil, false, fmt.Errorf("chart length mismatch")
	}
	return payload, whole, nil
}
func chartWrite(stage ChartStage, payload []byte, start int64) error {
	n, err := stage.WriteAt(payload, start)
	if err == nil && n != len(payload) {
		return io.ErrShortWrite
	}
	return err
}
func chartRemainder(ctx context.Context, client *http.Client, plan ChartPlan, stage ChartStage) error {
	run, cancel := context.WithCancel(ctx)
	defer cancel()
	count := (plan.Size - 1) / plan.Chunk
	faults := make([]error, count)
	jobs := make(chan int64)
	var wg sync.WaitGroup
	for i := int64(0); i < min(int64(plan.Workers), count); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for part := range jobs {
				if run.Err() != nil {
					return
				}
				start := part * plan.Chunk
				payload, _, err := chartPart(run, client, plan, start)
				if err == nil {
					err = chartWrite(stage, payload, start)
				}
				if err != nil {
					faults[part-1] = err
					cancel()
					return
				}
			}
		}()
	}
feed:
	for part := int64(1); part <= count; part++ {
		if run.Err() != nil {
			break
		}
		select {
		case jobs <- part:
		case <-run.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()
	return errors.Join(ctx.Err(), errors.Join(faults...))
}
func DownloadChart(ctx context.Context, client *http.Client, plan ChartPlan, stage ChartStage) (err error) {
	req, parseErr := http.NewRequestWithContext(ctx, http.MethodGet, plan.Address, nil)
	if client == nil || stage == nil || plan.Size < 0 || plan.Size > 1048576 || plan.Chunk <= 0 || plan.Workers <= 0 || parseErr != nil {
		return fmt.Errorf("invalid chart plan")
	}
	if req.URL.Host == "" || (req.URL.Scheme != "http" && req.URL.Scheme != "https") {
		return fmt.Errorf("invalid chart address")
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, stage.Abort())
		}
	}()
	if err = ctx.Err(); err != nil {
		return err
	}
	if plan.Size == 0 {
		return stage.Commit()
	}
	run, cancel := context.WithCancel(ctx)
	defer cancel()
	payload, whole, err := chartPart(run, client, plan, 0)
	if err != nil {
		return err
	}
	if err = chartWrite(stage, payload, 0); err != nil {
		return err
	}
	if !whole {
		if err = chartRemainder(run, client, plan, stage); err != nil {
			return err
		}
	}
	if err = ctx.Err(); err != nil {
		return err
	}
	return stage.Commit()
}
func main() {}
