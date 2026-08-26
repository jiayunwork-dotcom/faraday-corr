package faraday

func CumulativeMassLoss(molarMass, valence, iCorr, area, durationY float64) float64 {
	perAreaPerYear := AnnualMassLossPerArea(molarMass, valence, iCorr)
	return perAreaPerYear * area * durationY
}

func CumulativeMassLossFromRate(mdot, area, durationY float64) float64 {
	return mdot * yearSeconds * area * durationY
}

func DurationSeconds(durationY float64) float64 {
	return durationY * yearSeconds
}
