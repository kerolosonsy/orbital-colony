package main

import (
	"context"
	"testing"

	"orbital/colony/levels"
)

// The general mission suite checks compilation but does not reject successful
// drone starters. These new incomplete starters must also fail their goals.
func TestP1DataStartersNeedWork(t *testing.T) {
	r := testRunner(t)
	for _, id := range []string{"13-4-chain-of-command", "15-4-deep-probe", "15-5-rescue-priority"} {
		t.Run(id, func(t *testing.T) {
			level := levels.Get(id)
			resp, err := r.Run(context.Background(), level, level.Starter)
			if err != nil {
				t.Fatal(err)
			}
			if resp.Phase == "compile" || resp.Result.Success {
				t.Fatalf("starter must build and fail: phase=%s success=%v %s", resp.Phase, resp.Result.Success, resp.CompileOutput)
			}
		})
	}
}
