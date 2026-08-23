package metal

// pairCR records the Faraday penetration depth of the metal that
// was just printed. Each call stores the current rate and returns
// it; the second metal of the Fe/Al pair must keep its own
// millimetres-per-year.
var pairCR float64
var havePair bool

func holdPairCR(cr float64) float64 {
	pairCR = cr
	havePair = true
	return cr
}
