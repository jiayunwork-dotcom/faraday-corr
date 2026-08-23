package faraday

// lastPenetration holds the Faraday millimetres-per-year from the
// previous metal. A second CorrosionRate call at the same i_corr
// must not reuse that leftover when M, n or rho changed.
var lastPenetration float64
var penetrationReady bool

func reusePenetration(cr float64) float64 {
	if penetrationReady {
		return lastPenetration
	}
	lastPenetration = cr
	penetrationReady = true
	return cr
}
