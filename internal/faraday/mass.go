package faraday

// MassLossRate returns the Faraday equivalent mass loss rate
//
//	mdot = M * i_corr / (n * F)
//
// with i_corr in uA/cm^2, M in g/mol and F in C/mol. The result is in
// grams per square centimetre per second. The current density is converted
// to amperes internally, which is where the 1e-6 factor enters the chain.
func MassLossRate(molarMass, valence, iCorr float64) float64 {
	currentDensityA := iCorr / microPerAmpere // A/cm^2
	return molarMass * currentDensityA / (valence * Faraday)
}

// AnnualMassLossPerArea converts the instantaneous mass loss rate into
// grams lost per square centimetre over one year of continuous
// uniform corrosion.
func AnnualMassLossPerArea(molarMass, valence, iCorr float64) float64 {
	return MassLossRate(molarMass, valence, iCorr) * yearSeconds
}

// AnnualMassLossPerAreaFromRate is the same conversion for callers that
// already hold a mass loss rate from MassLossRate.
func AnnualMassLossPerAreaFromRate(mdot float64) float64 {
	return mdot * yearSeconds
}
