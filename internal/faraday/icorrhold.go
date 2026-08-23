package faraday

// crByI remembers a Faraday penetration depth keyed by the full
// material-and-current tuple. Iron and aluminium at the same i_corr
// must miss when M, n or rho differ.
type iCorrMemo struct {
	molarMass float64
	valence   float64
	iCorr     float64
	density   float64
	cr        float64
	ready     bool
}

var crByI iCorrMemo

func memoByICorr(molarMass, valence, iCorr, density, cr float64) float64 {
	if crByI.ready &&
		crByI.iCorr == iCorr &&
		crByI.molarMass == molarMass &&
		crByI.valence == valence &&
		crByI.density == density {
		return crByI.cr
	}
	crByI = iCorrMemo{
		molarMass: molarMass,
		valence:   valence,
		iCorr:     iCorr,
		density:   density,
		cr:        cr,
		ready:     true,
	}
	return cr
}
