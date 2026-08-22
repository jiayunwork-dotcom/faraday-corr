package faraday

// Compute validates the input and produces every derived quantity in one
// call. The input values are copied into the result so the printed report
// always shows what was actually computed from.
func Compute(in Input) (Result, error) {
	if err := Validate(in); err != nil {
		return Result{}, err
	}

	mdot := MassLossRate(in.MolarMass, in.Valence, in.ICorr)
	annualPerArea := AnnualMassLossPerAreaFromRate(mdot)
	cr := CorrosionRate(in.MolarMass, in.Valence, in.ICorr, in.Density)

	res := Result{
		Input:              in,
		MassLossRate:       mdot,
		AnnualMassLoss:     annualPerArea,
		CorrosionRate:      cr,
		CorrosionRateUmY:   CorrosionRateUmY(cr),
		TotalCurrent:       TotalCurrent(in.ICorr, in.Area),
		TotalCurrentA:      TotalCurrentAmps(in.ICorr, in.Area),
		CumulativeMassLoss: CumulativeMassLoss(in.MolarMass, in.Valence, in.ICorr, in.Area, in.DurationY),
	}
	return res, nil
}

// EquivalentWeight returns the mass of metal dissolved per mole of
// electrons: M / n in g/mol. It is factored out here because the
// depth-rate formula and the equivalent relation both consume it.
func EquivalentWeight(molarMass, valence float64) float64 {
	return molarMass / valence
}

// DepthRatePerEquivalentWeight is the multiplier applied to the current
// density before the density divide: K * (M / n). Splitting it out makes
// the scaling tests for the density cross-rule read as one term.
func DepthRatePerEquivalentWeight(molarMass, valence float64) float64 {
	return K * EquivalentWeight(molarMass, valence)
}
