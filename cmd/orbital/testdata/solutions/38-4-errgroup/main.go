package main

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

func SurveyAll(ctx context.Context, jobs []int, limit int, work func(context.Context, int) error) error {
	if limit <= 0 {
		return errors.New("limit must be positive")
	}
	g, groupCtx := errgroup.WithContext(ctx)
	g.SetLimit(limit)
	for _, job := range jobs {
		if groupCtx.Err() != nil {
			break
		}
		job := job
		g.Go(func() error {
			if err := groupCtx.Err(); err != nil {
				return err
			}
			return work(groupCtx, job)
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return ctx.Err()
}
func main() {}
