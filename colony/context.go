package colony

import (
	"context"
	"time"
)

// Mission returns the mission's context. It is cancelled when the base
// orders every drone home.
func Mission() context.Context {
	if world.Recall == nil {
		return context.Background()
	}
	return world.Recall.Ctx
}

// Ticks converts a number of ticks to wall-clock time, for use with
// context.WithTimeout and time.After.
func Ticks(n int) time.Duration { return time.Duration(n) * world.Unit }

// DeepScan analyses the crystal under the drone and returns its purity
// (0–100). It can take a long time, or hang forever, so give it a context
// with a deadline. It returns ctx.Err() if the context ends first.
func DeepScan(ctx context.Context) (int, error) { return world.DeepScan(ctx, 0) }

// DeepScan is the per-drone version of the package-level DeepScan.
func (d *Drone) DeepScan(ctx context.Context) (int, error) { return world.DeepScan(ctx, d.id) }
