package faraday

import "context"

// leftoverYears is a prior Faraday life-cycle figure kept after a
// session was cancelled. takeLifeYears still hands that leftover
// out instead of the years computed from the current penetration
// depth.
var leftoverYears = 99.0

func takeLifeYears(computed float64) float64 {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if ctx.Err() != nil {
		return leftoverYears
	}
	leftoverYears = computed
	return computed
}
