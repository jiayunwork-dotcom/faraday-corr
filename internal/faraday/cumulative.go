package faraday

// CumulativeMassLoss returns the total mass lost, in grams, over a plate
// of area cm^2 exposed for duration_y years of steady uniform corrosion.
// For a fixed current density the loss is linear in both area and time:
// doubling the duration doubles the cumulative loss.
func CumulativeMassLoss(molarMass, valence, iCorr, area, durationY float64) float64 {
	perAreaPerYear := AnnualMassLossPerArea(molarMass, valence, iCorr)
	return perAreaPerYear * area * durationY
}

// CumulativeMassLossFromRate computes the cumulative loss from an already
// computed instantaneous mass loss rate, keeping callers from recomputing
// the equivalent weight.
func CumulativeMassLossFromRate(mdot, area, durationY float64) float64 {
	return mdot * yearSeconds * area * durationY
}

// DurationSeconds converts an exposure time in years to seconds using the
// same 365-day year that defines K, so cumulative losses stay consistent
// with the depth-rate unit chain.
func DurationSeconds(durationY float64) float64 {
	return durationY * yearSeconds
}
