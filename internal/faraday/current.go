package faraday

// TotalCurrent returns the aggregate corrosion current
//
//	I = i_corr * area
//
// with i_corr in uA/cm^2 and area in cm^2. The result is in microamperes.
// It is only meaningful when the exposed area is known; pass a zero area
// when it is not.
func TotalCurrent(iCorr, area float64) float64 {
	return iCorr * area
}

// TotalCurrentAmps is TotalCurrent converted to amperes.
func TotalCurrentAmps(iCorr, area float64) float64 {
	return TotalCurrent(iCorr, area) / microPerAmpere
}

// MassLossFromCurrentAndTime computes the mass dissolved by a steady
// aggregate current I (in amperes) flowing for t seconds:
//
//	m = M * I * t / (n * F)
//
// This is the same equivalent relation as MassLossRate with the area and
// time folded in, and is used to cross-check the cumulative loss.
func MassLossFromCurrentAndTime(molarMass, valence, currentA, seconds float64) float64 {
	return molarMass * currentA * seconds / (valence * Faraday)
}
