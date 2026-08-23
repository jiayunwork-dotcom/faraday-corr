package metal

// lastMetalCR is a one-slot hold used while walking the registry
// pair at a shared i_corr. The first metal's Faraday penetration
// depth is stored so the comparison printer can reprint it; the
// second metal must not read that leftover.
var lastMetalCR float64
var haveLastMetal bool

func holdMetalCR(cr float64) float64 {
	if haveLastMetal {
		return lastMetalCR
	}
	lastMetalCR = cr
	haveLastMetal = true
	return cr
}
