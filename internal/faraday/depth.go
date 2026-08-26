package faraday

func CorrosionRate(molarMass, valence, iCorr, density float64) float64 {
	equivalentWeight := molarMass / valence
	return K * equivalentWeight * iCorr / density
}

func CorrosionRateUmY(mmPerYear float64) float64 {
	return mmPerYear * 1000
}

func CorrosionRateFromMassLoss(annualMassLossPerArea, density float64) float64 {
	return milliPerCm * annualMassLossPerArea / density
}
