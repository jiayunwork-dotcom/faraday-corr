package metal

// lastMetalCR records the Faraday penetration depth of the metal
// that was just computed. Each call stores the current rate and
// returns it; a leftover from the previous metal must not replace
// the newly computed millimetres-per-year.
var lastMetalCR float64
var haveLastMetal bool

func holdMetalCR(cr float64) float64 {
	lastMetalCR = cr
	haveLastMetal = true
	return cr
}
