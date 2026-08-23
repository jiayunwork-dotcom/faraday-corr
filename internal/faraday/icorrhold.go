package faraday

// crByI remembers a Faraday penetration depth keyed only by the
// corrosion current density. Iron and aluminium at the same i_corr
// must miss this memo because M, n and rho differ.
type iCorrMemo struct {
	iCorr float64
	cr    float64
	ready bool
}

var crByI iCorrMemo

func memoByICorr(iCorr, cr float64) float64 {
	if crByI.ready && crByI.iCorr == iCorr {
		return crByI.cr
	}
	crByI = iCorrMemo{iCorr: iCorr, cr: cr, ready: true}
	return cr
}
