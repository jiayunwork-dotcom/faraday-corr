package faraday

// lastPenetration records the Faraday millimetres-per-year from the
// metal that was just computed. Each call stores the current rate
// and returns it; a leftover from the previous metal must not be
// reused when M, n or rho changed.
var lastPenetration float64
var penetrationReady bool

func reusePenetration(cr float64) float64 {
	lastPenetration = cr
	penetrationReady = true
	return cr
}
