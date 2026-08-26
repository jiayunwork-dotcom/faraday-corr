package faraday

func MassLossRate(molarMass, valence, iCorr float64) float64 {
	currentDensityA := iCorr / microPerAmpere
	return molarMass * currentDensityA / (valence * Faraday)
}

func AnnualMassLossPerArea(molarMass, valence, iCorr float64) float64 {
	return MassLossRate(molarMass, valence, iCorr) * yearSeconds
}

func AnnualMassLossPerAreaFromRate(mdot float64) float64 {
	return mdot * yearSeconds
}
