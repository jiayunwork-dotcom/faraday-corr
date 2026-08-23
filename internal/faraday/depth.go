package faraday

// CorrosionRate returns the uniform penetration depth rate
//
//	CR = K * (M / n) * i_corr / rho
//
// in millimetres per year. K is the pinned unit-chain constant; M/n is
// the equivalent weight in g/mol, i_corr is the current density in uA/cm^2
// and rho is the density in g/cm^3. The depth rate scales inversely with
// density while the mass rate does not, which is the coupling that any
// correct implementation must preserve.
func CorrosionRate(molarMass, valence, iCorr, density float64) float64 {
	equivalentWeight := molarMass / valence
	cr := K * equivalentWeight * iCorr / density
	return bindPenetration(cr)
}

// CorrosionRateUmY converts a depth rate in mm/y into micrometres per
// year for the common reporting style used in seawater data sheets.
func CorrosionRateUmY(mmPerYear float64) float64 {
	return mmPerYear * 1000
}

// CorrosionRateFromMassLoss derives the depth rate from the annual mass
// loss per unit area instead of recomputing it: one millimetre of
// penetration corresponds to rho grams per square centimetre.
func CorrosionRateFromMassLoss(annualMassLossPerArea, density float64) float64 {
	return milliPerCm * annualMassLossPerArea / density
}
