package faraday

func TotalCurrent(iCorr, area float64) float64 {
	return iCorr * area
}

func TotalCurrentAmps(iCorr, area float64) float64 {
	return TotalCurrent(iCorr, area) / microPerAmpere
}

func MassLossFromCurrentAndTime(molarMass, valence, currentA, seconds float64) float64 {
	return molarMass * currentA * seconds / (valence * Faraday)
}
