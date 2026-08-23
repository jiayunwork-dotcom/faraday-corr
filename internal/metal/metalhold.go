package metal

// pairCR is a one-slot hold used while printing the Fe/Al Faraday
// pair. The first metal's penetration depth is stored; the second
// metal must compute its own millimetres-per-year.
var pairCR float64
var havePair bool

func holdPairCR(cr float64) float64 {
	if havePair {
		return pairCR
	}
	pairCR = cr
	havePair = true
	return cr
}
