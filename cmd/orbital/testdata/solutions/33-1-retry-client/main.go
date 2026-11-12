package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Pauser waits between attempts; tests inject one that never sleeps.
type Pauser interface {
	Pause(ctx context.Context, d time.Duration) error
}

type Policy struct {
	Attempts int
	Base     time.Duration
	Max      time.Duration
}

type Retrier struct {
	next   http.RoundTripper
	p      Policy
	pause  Pauser
	jitter func(n int64) int64
}

func NewRetrier(next http.RoundTripper, p Policy, pause Pauser, jitter func(n int64) int64) (*Retrier, error) {
	if next == nil || pause == nil || jitter == nil || p.Attempts < 1 || p.Attempts > 10 || p.Base <= 0 || p.Max < p.Base {
		return nil, errors.New("invalid retry policy")
	}
	return &Retrier{next: next, p: p, pause: pause, jitter: jitter}, nil
}

// replayable reports whether sending req again cannot repeat a side effect by accident.
func replayable(req *http.Request) bool {
	switch req.Method {
	case "GET", "HEAD", "OPTIONS", "PUT", "DELETE":
	case "POST", "PATCH":
		if req.Header.Get("Idempotency-Key") == "" {
			return false
		}
	default:
		return false
	}
	return req.Body == nil || req.Body == http.NoBody || req.GetBody != nil
}

func retryStatus(code int) bool {
	return code == 429 || code == 502 || code == 503 || code == 504
}

// wait returns the delay before retry number n (1 for the first retry).
func (r *Retrier) wait(n int, resp *http.Response) time.Duration {
	if resp != nil && (resp.StatusCode == 429 || resp.StatusCode == 503) {
		if secs, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && secs >= 0 {
			return min(time.Duration(secs)*time.Second, r.p.Max)
		}
	}
	d := r.p.Base
	for i := 1; i < n && d < r.p.Max; i++ {
		d *= 2
	}
	d = min(d, r.p.Max)
	half := d / 2
	return half + time.Duration(r.jitter(int64(d-half)+1))
}

func (r *Retrier) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	tries := r.p.Attempts
	if !replayable(req) {
		tries = 1
	}
	for n := 1; ; n++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		out := req
		if n > 1 && req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, err
			}
			out = req.Clone(ctx)
			out.Body = body
		}
		resp, err := r.next.RoundTrip(out)
		if err == nil && !retryStatus(resp.StatusCode) {
			return resp, nil
		}
		if err != nil && ctx.Err() != nil {
			return nil, err
		}
		if n >= tries {
			return resp, err
		}
		d := r.wait(n, resp)
		if resp != nil {
			io.Copy(io.Discard, io.LimitReader(resp.Body, 4096))
			resp.Body.Close()
		}
		if err := r.pause.Pause(ctx, d); err != nil {
			return nil, err
		}
	}
}

func main() {}
